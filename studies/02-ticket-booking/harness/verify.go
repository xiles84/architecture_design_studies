package main

import (
	"context"
	"fmt"
	"math/rand"
	"sort"
	"strings"

	"adsplatform/core/catalog"
	"adsplatform/ports"
)

// ---------------------------------------------------------------------------
// Correctness gate.
//
// Before any timing, every design must answer the five read questions exactly as
// computed in Go from the generated dataset, and its tables must pass the
// overbooking audit on quiescent data. A failure aborts the cell: a design that
// answers wrongly, quickly, is worth nothing.
//
// Keys are chosen to make the checks mean something: for every size tier, the
// most-sold catalogue event and an unsold race event; the band with the most
// events; the customer with the most tickets; and a random sample of tickets.
// An event with nothing sold and one with everything sold both pass trivially
// in a broken design, so neither is the only thing checked.
// ---------------------------------------------------------------------------

type Check struct {
	Query  string `json:"query"`
	Key    string `json:"key"`
	OK     bool   `json:"ok"`
	Expect string `json:"expect,omitempty"`
	Got    string `json:"got,omitempty"`
}

type VerifyReport struct {
	Passed int     `json:"passed"`
	Failed int     `json:"failed"`
	Checks []Check `json:"checks"`
}

func (v *VerifyReport) add(c Check) {
	if c.OK {
		v.Passed++
	} else {
		v.Failed++
	}
	// Passing checks are kept only as a count beyond the first few, so a result
	// file stays readable; every failure is kept.
	if !c.OK || len(v.Checks) < 40 {
		v.Checks = append(v.Checks, c)
	}
}

type keySet struct {
	tierEvents []*Event
	band       int64
	otherBand  int64
	customer   int64
	tickets    []*Ticket
	busyEvent  *Event // largest catalogue event, for plans
	soldTicket *Ticket
}

func chooseKeys(ds *Dataset) keySet {
	var k keySet
	for _, tier := range allTiers {
		var best *Event
		for _, e := range ds.eventsOf("catalogue", tier) {
			if best == nil || e.InitialSold > best.InitialSold {
				best = e
			}
		}
		if best != nil {
			k.tierEvents = append(k.tierEvents, best)
			if k.busyEvent == nil || best.Capacity > k.busyEvent.Capacity {
				k.busyEvent = best
			}
		}
		if race := ds.eventsOf("race", tier); len(race) > 0 {
			k.tierEvents = append(k.tierEvents, race[0])
		}
	}
	perBand := map[int64]int{}
	for _, e := range ds.Events {
		perBand[e.BandID]++
	}
	for id, n := range perBand {
		if k.band == 0 || n > perBand[k.band] || (n == perBand[k.band] && id < k.band) {
			k.band = id
		}
	}
	k.otherBand = int64(len(ds.Bands)/2) + 1
	perCustomer := map[int64]int{}
	for _, t := range ds.Sold {
		perCustomer[t.CustomerID]++
	}
	for id, n := range perCustomer {
		if k.customer == 0 || n > perCustomer[k.customer] || (n == perCustomer[k.customer] && id < k.customer) {
			k.customer = id
		}
	}
	r := rand.New(rand.NewSource(ds.Seed + 7))
	for i := 0; i < 25 && len(ds.Sold) > 0; i++ {
		k.tickets = append(k.tickets, &ds.Sold[r.Intn(len(ds.Sold))])
	}
	if k.busyEvent != nil {
		for i := range ds.Sold {
			if ds.Sold[i].EventID == k.busyEvent.ID {
				k.soldTicket = &ds.Sold[i]
				break
			}
		}
	}
	return k
}

func Verify(ctx context.Context, db ports.DB, d Design, ds *Dataset) (*VerifyReport, error) {
	stmts, err := mustStmts(d.ID, "queries.sql")
	if err != nil {
		return nil, err
	}
	q := catalog.Map(stmts)
	k := chooseKeys(ds)
	v := &VerifyReport{}

	// q01: seats left, per tier, busy and unsold events.
	for _, e := range k.tierEvents {
		var got int64
		err := db.QueryRow(ctx, q["q01_event_availability"].SQL, e.ID).Scan(&got)
		want := int64(e.Capacity - e.InitialSold)
		v.add(Check{Query: "q01_event_availability", Key: fmt.Sprintf("event %d (%s, %d seats)", e.ID, e.Kind, e.Capacity),
			OK: err == nil && got == want, Expect: fmt.Sprint(want), Got: gotOrErr(got, err)})
	}

	// q02 and q05: per band.
	for _, band := range []int64{k.band, k.otherBand} {
		var evs []*Event
		for i := range ds.Events {
			if ds.Events[i].BandID == band {
				evs = append(evs, &ds.Events[i])
			}
		}
		sort.Slice(evs, func(i, j int) bool {
			if !evs[i].StartsAt.Equal(evs[j].StartsAt) {
				return evs[i].StartsAt.Before(evs[j].StartsAt)
			}
			return evs[i].ID < evs[j].ID
		})
		var want []string
		var sold int64
		for _, e := range evs {
			want = append(want, fmt.Sprintf("%d:%d:%d", e.ID, e.Capacity, e.Capacity-e.InitialSold))
			sold += int64(e.InitialSold)
		}
		got, err := scanStrings(ctx, db, q["q02_band_events"].SQL, band, 3)
		v.add(Check{Query: "q02_band_events", Key: fmt.Sprintf("band %d (%d events)", band, len(evs)),
			OK:     err == nil && strings.Join(got, ",") == strings.Join(want, ","),
			Expect: abbreviate(want), Got: abbreviateErr(got, err)})

		var gotSold int64
		err = db.QueryRow(ctx, q["q05_band_tickets_sold"].SQL, band).Scan(&gotSold)
		v.add(Check{Query: "q05_band_tickets_sold", Key: fmt.Sprintf("band %d", band),
			OK: err == nil && gotSold == sold, Expect: fmt.Sprint(sold), Got: gotOrErr(gotSold, err)})
	}

	// q03: a customer's tickets.
	var want []string
	for _, t := range ds.Sold {
		if t.CustomerID == k.customer {
			want = append(want, fmt.Sprintf("%d:%d:%d", t.ID, t.EventID, t.SeatNo))
		}
	}
	sort.Slice(want, func(i, j int) bool { return idOf(want[i]) < idOf(want[j]) })
	got, err := scanStrings(ctx, db, q["q03_customer_tickets"].SQL, k.customer, 3)
	v.add(Check{Query: "q03_customer_tickets", Key: fmt.Sprintf("customer %d (%d tickets)", k.customer, len(want)),
		OK: err == nil && strings.Join(got, ",") == strings.Join(want, ","), Expect: abbreviate(want), Got: abbreviateErr(got, err)})

	// q04: ticket by id.
	for _, t := range k.tickets {
		var id, ev, cust int64
		var seat int32
		err := db.QueryRow(ctx, q["q04_ticket_by_id"].SQL, t.ID).Scan(&id, &ev, &seat, &cust)
		want := fmt.Sprintf("%d:%d:%d:%d", t.ID, t.EventID, t.SeatNo, t.CustomerID)
		gotS := fmt.Sprintf("%d:%d:%d:%d", id, ev, seat, cust)
		if err != nil {
			gotS = "error: " + err.Error()
		}
		v.add(Check{Query: "q04_ticket_by_id", Key: fmt.Sprintf("ticket %d", t.ID), OK: err == nil && gotS == want, Expect: want, Got: gotS})
	}
	return v, nil
}

func idOf(s string) int64 {
	var id int64
	fmt.Sscanf(s, "%d:", &id)
	return id
}

// scanStrings renders each row's integer columns as "a:b:c", draining every row.
func scanStrings(ctx context.Context, db ports.DB, sql string, arg any, cols int) ([]string, error) {
	rows, err := db.Query(ctx, sql, arg)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []string
	for rows.Next() {
		vals := make([]*int64, cols)
		dest := make([]any, cols)
		for i := range vals {
			dest[i] = &vals[i]
		}
		if err := rows.Scan(dest...); err != nil {
			return nil, err
		}
		parts := make([]string, cols)
		for i, p := range vals {
			if p == nil {
				parts[i] = "null"
			} else {
				parts[i] = fmt.Sprint(*p)
			}
		}
		out = append(out, strings.Join(parts, ":"))
	}
	return out, rows.Err()
}

func gotOrErr(v int64, err error) string {
	if err != nil {
		return "error: " + err.Error()
	}
	return fmt.Sprint(v)
}

func abbreviate(xs []string) string {
	s := strings.Join(xs, ",")
	if len(s) > 160 {
		return fmt.Sprintf("%s... (%d rows)", s[:157], len(xs))
	}
	return s
}

func abbreviateErr(xs []string, err error) string {
	if err != nil {
		return "error: " + err.Error()
	}
	return abbreviate(xs)
}
