package main

import (
	"context"
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"time"

	"adsplatform/core/catalog"
	"adsplatform/ports"
)

// ---------------------------------------------------------------------------
// Correctness gate.
//
// Before any timing, every design must answer the six read questions exactly as
// computed in Go from the generated dataset, on a load that contains sold seats,
// live holds AND expired holds that nobody has cleaned up -- the state that tells
// lazy expiry from a design that forgot about it. A failure aborts the cell: a
// design that answers wrongly, quickly, is worth nothing.
//
// Expiry makes the truth time-dependent. Loaded live holds expire 10-40 minutes
// after the load and expired ones ended up to an hour before it, so the answer is
// exact as long as verification starts at least five minutes before the earliest
// live hold expires; the gate checks that first rather than assuming it.
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
	// MarginMS is how far the earliest live loaded hold was from expiring when
	// verification started.
	MarginMS float64 `json:"margin_ms"`
}

func (v *VerifyReport) add(c Check) {
	if c.OK {
		v.Passed++
	} else {
		v.Failed++
	}
	if !c.OK || len(v.Checks) < 40 {
		v.Checks = append(v.Checks, c)
	}
}

const gateMargin = 5 * time.Minute

// columns is each read's result width, identical in every design.
var columns = map[string]int{"q01_section_map": 2, "q02_event_sections": 2, "q03_event_available": 1,
	"q04_hold_seats": 2, "q05_customer_tickets": 3, "q06_ticket_by_id": 4,
	// Operational reports (REPORTS.md v2, AM-02).
	"r01_holds_expiring_soon": 4, "r02_seat_status_lookup": 4, "r03_section_sales_window": 3,
	"r04_event_recent_confirmations": 4, "r05_customers_last_purchase_window": 2}

// Verify runs the gate. E0's statements read the application clock with no skew.
func Verify(ctx context.Context, db ports.DB, d Design, ds *Dataset, loadNow time.Time) (*VerifyReport, error) {
	stmts, err := mustStmts(d.ID, "queries.sql")
	if err != nil {
		return nil, err
	}
	q := catalog.Map(stmts)
	vr := &VerifyReport{}

	var now time.Time
	if err := db.QueryRow(ctx, "SELECT now()").Scan(&now); err != nil {
		return nil, err
	}
	earliest := time.Time{}
	for _, h := range ds.Holds {
		if h.live() && (earliest.IsZero() || h.expiresAt(loadNow).Before(earliest)) {
			earliest = h.expiresAt(loadNow)
		}
	}
	if !earliest.IsZero() {
		vr.MarginMS = float64(earliest.Sub(now).Microseconds()) / 1000
		if earliest.Sub(now) < gateMargin {
			return nil, fmt.Errorf("load took too long: the earliest loaded live hold expires in %s (< %s), so the gate's truth would be stale",
				earliest.Sub(now).Round(time.Second), gateMargin)
		}
	}

	run := func(name string, key string, args []any, want []string) {
		vals := map[string]any{"app_now": time.Now()}
		st := q[name]
		// Arguments are positional in the order the catalogue names them; E0's
		// statements add app_now last.
		bound := args
		if len(st.Params) > len(args) {
			bound = append(append([]any(nil), args...), vals["app_now"])
		}
		got, err := collect(ctx, db, st.SQL, bound, columns[name])
		c := Check{Query: name, Key: key}
		if err != nil {
			c.Got = "error: " + err.Error()
		} else if strings.Join(got, ";") == strings.Join(want, ";") {
			c.OK = true
		} else {
			c.Expect, c.Got = abbreviate(want), abbreviate(got)
		}
		vr.add(c)
	}

	for _, tier := range ds.tiersOf("catalogue") {
		ev := mostSold(ds, tier)
		v := ds.venue(ev.VenueID)
		st := ds.State[ev.ID]
		for _, sec := range []int32{sectionMost(v, st, stSold), sectionMost(v, st, stLive)} {
			run("q01_section_map", fmt.Sprintf("event %d section %d", ev.ID, sec), []any{ev.ID, sec}, truthMap(v, st, sec))
		}
		run("q02_event_sections", fmt.Sprintf("event %d", ev.ID), []any{ev.ID}, truthSections(v, st))
		run("q03_event_available", fmt.Sprintf("event %d", ev.ID), []any{ev.ID}, []string{fmt.Sprint(truthAvailable(st))})
	}
	for _, tier := range ds.tiersOf("race") {
		ev := ds.eventsOf("race", tier)[0]
		v := ds.venue(ev.VenueID)
		free := make([]uint8, v.Capacity)
		run("q01_section_map", fmt.Sprintf("race event %d section 1", ev.ID), []any{ev.ID, int32(1)}, truthMap(v, free, 1))
		run("q02_event_sections", fmt.Sprintf("race event %d", ev.ID), []any{ev.ID}, truthSections(v, free))
		run("q03_event_available", fmt.Sprintf("race event %d", ev.ID), []any{ev.ID}, []string{fmt.Sprint(v.Capacity)})
	}

	var biggestLive, anExpired *LoadedHold
	for i := range ds.Holds {
		h := &ds.Holds[i]
		if h.live() && (biggestLive == nil || len(h.Seats) > len(biggestLive.Seats)) {
			biggestLive = h
		}
		if !h.live() && anExpired == nil {
			anExpired = h
		}
	}
	if biggestLive != nil {
		var want []string
		for _, s := range biggestLive.Seats {
			want = append(want, fmt.Sprintf("%d@%d", s, biggestLive.expiresAt(loadNow).UnixMilli()))
		}
		run("q04_hold_seats", fmt.Sprintf("live hold %d", biggestLive.ID), []any{biggestLive.EventID, biggestLive.Section, biggestLive.ID}, want)
	}
	if anExpired != nil {
		run("q04_hold_seats", fmt.Sprintf("expired hold %d", anExpired.ID), []any{anExpired.EventID, anExpired.Section, anExpired.ID}, nil)
	}

	byCustomer := map[int64][]*Ticket{}
	for i := range ds.Sold {
		t := &ds.Sold[i]
		byCustomer[t.CustomerID] = append(byCustomer[t.CustomerID], t)
	}
	var busiest int64
	for c, ts := range byCustomer {
		if busiest == 0 || len(ts) > len(byCustomer[busiest]) || (len(ts) == len(byCustomer[busiest]) && c < busiest) {
			busiest = c
		}
	}
	if busiest != 0 {
		ts := byCustomer[busiest]
		sort.Slice(ts, func(i, j int) bool { return ts[i].ID < ts[j].ID })
		var want []string
		for _, t := range ts {
			want = append(want, fmt.Sprintf("%d|%d|%d", t.ID, t.EventID, t.SeatID))
		}
		run("q05_customer_tickets", fmt.Sprintf("customer %d", busiest), []any{busiest}, want)
	}
	r := rand.New(rand.NewSource(ds.Seed + 7))
	for i := 0; i < 20 && len(ds.Sold) > 0; i++ {
		t := ds.Sold[r.Intn(len(ds.Sold))]
		run("q06_ticket_by_id", fmt.Sprintf("ticket %d", t.ID), []any{t.ID},
			[]string{fmt.Sprintf("%d|%d|%d|%d", t.ID, t.EventID, t.SeatID, t.CustomerID)})
	}

	verifyReports(ctx, db, d, ds, q, loadNow, vr)
	return vr, nil
}

// collect renders every row as "a|b|c", with timestamps as epoch milliseconds
// (rounded: L2 rebuilds them from milliseconds through a float division), and a
// (seat, expires) pair as "seat@ms".
func collect(ctx context.Context, db ports.DB, sql string, args []any, n int) ([]string, error) {
	rows, err := db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []string
	for rows.Next() {
		var a, b, c, d any
		if err := rows.Scan([]any{&a, &b, &c, &d}[:n]...); err != nil {
			return nil, err
		}
		parts := []string{}
		for _, v := range []any{a, b, c, d}[:n] {
			parts = append(parts, render(v))
		}
		if n == 2 {
			if t, ok := b.(time.Time); ok {
				out = append(out, fmt.Sprintf("%s@%d", render(a), t.Round(time.Millisecond).UnixMilli()))
				continue
			}
		}
		out = append(out, strings.Join(parts, "|"))
	}
	return out, rows.Err()
}

func render(v any) string {
	switch x := v.(type) {
	case nil:
		return "NULL"
	case time.Time:
		return fmt.Sprint(x.Round(time.Millisecond).UnixMilli())
	case int16, int32, int64, int:
		return fmt.Sprint(x)
	default:
		return fmt.Sprint(x)
	}
}

func abbreviate(xs []string) string {
	s := strings.Join(xs, ";")
	if len(s) > 400 {
		return s[:400] + fmt.Sprintf("... (%d rows)", len(xs))
	}
	if s == "" {
		return "(no rows)"
	}
	return s
}

func mostSold(ds *Dataset, tier int) *Event {
	var best *Event
	bestN := -1
	for _, e := range ds.eventsOf("catalogue", tier) {
		n := 0
		for _, st := range ds.State[e.ID] {
			if st == stSold {
				n++
			}
		}
		if n > bestN {
			best, bestN = e, n
		}
	}
	return best
}

func sectionMost(v *Venue, st []uint8, want uint8) int32 {
	best, bestN := int32(1), -1
	for _, sec := range v.Sections {
		n := 0
		for s := sec.FirstSeat; s < sec.FirstSeat+int32(sec.Capacity); s++ {
			if st[s-1] == want {
				n++
			}
		}
		if n > bestN {
			best, bestN = sec.No, n
		}
	}
	return best
}

func truthMap(v *Venue, st []uint8, secNo int32) []string {
	sec := v.section(secNo)
	var out []string
	for s := sec.FirstSeat; s < sec.FirstSeat+int32(sec.Capacity); s++ {
		switch st[s-1] {
		case stLive:
			out = append(out, fmt.Sprintf("%d|1", s))
		case stSold:
			out = append(out, fmt.Sprintf("%d|2", s))
		}
	}
	return out
}

func truthSections(v *Venue, st []uint8) []string {
	var out []string
	for _, sec := range v.Sections {
		n := 0
		for s := sec.FirstSeat; s < sec.FirstSeat+int32(sec.Capacity); s++ {
			if st[s-1] == stFree || st[s-1] == stExpired {
				n++
			}
		}
		out = append(out, fmt.Sprintf("%d|%d", sec.No, n))
	}
	return out
}

func truthAvailable(st []uint8) int {
	n := 0
	for _, x := range st {
		if x == stFree || x == stExpired {
			n++
		}
	}
	return n
}
