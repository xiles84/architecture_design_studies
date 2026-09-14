package main

import (
	"context"
	"fmt"
	"sort"
	"time"

	"adsplatform/core/catalog"
	"adsplatform/ports"
)

// ---------------------------------------------------------------------------
// The audit.
//
// Run after every phase that writes, on quiescent data (the phase has stopped and
// its sweeper finished), on the state that phase produced. It combines two
// independent views:
//
//	the database's   every ticket, the inventory's own view of sold seats, the
//	                 holds valid now, rows that do not fit the venue
//	the ledger's     what the harness saw commit, replayed seat by seat
//
// and reports, with examples rather than tallies alone:
//
//	INV-1  a seat with more than one ticket; a double sale in the ledger
//	INV-2  theft; a grant on a sold seat
//	INV-3  early rejections (classified when they happened)
//	INV-5  a live hold the database has only part of, or a valid hold the
//	       harness never saw granted
//	INV-6  late sales
//	INV-7  tickets that exist without a sale the harness saw, or sales the
//	       harness saw that left no ticket (allowing ambiguous commits)
//	drift  sold inventory that disagrees with the tickets; rows outside the venue
// ---------------------------------------------------------------------------

type Audit struct {
	Phase   string `json:"phase"`
	Events  int    `json:"events_checked"`
	Tickets int64  `json:"tickets"`

	Ledger Violations `json:"ledger"`

	DuplicateSeats     int64 `json:"duplicate_seats"`
	InventoryDrift     int64 `json:"inventory_drift"`
	InvalidSeats       int64 `json:"invalid_seats"`
	PartialHolds       int64 `json:"partial_or_unknown_holds"`
	TicketMismatches   int64 `json:"ticket_mismatch_events"`
	AmbiguousCommits   int64 `json:"ambiguous_commits,omitempty"`
	HoldsSkippedAtEdge int64 `json:"holds_skipped_at_expiry_edge,omitempty"`

	Examples []Example `json:"examples,omitempty"`
}

// Violations is the number that gates a correct design's phase.
func (a *Audit) Violations() int64 {
	if a == nil {
		return 0
	}
	return a.Ledger.Count() + a.DuplicateSeats + a.InventoryDrift + a.InvalidSeats + a.PartialHolds + a.TicketMismatches
}

func (a *Audit) String() string {
	if a.Violations() == 0 {
		return fmt.Sprintf("consistent (%d events, %d tickets; rejections late=%d boundary=%d)",
			a.Events, a.Tickets, a.Ledger.RejectedLate, a.Ledger.RejectedBoundary)
	}
	return fmt.Sprintf("VIOLATIONS — thefts %d, grants on sold seats %d, double sales %d, sales without the hold %d, late sales %d, early rejections %d, duplicate seats %d, inventory drift %d, invalid seats %d, partial holds %d, ticket mismatches %d",
		a.Ledger.Thefts, a.Ledger.SoldSeatGrants, a.Ledger.DoubleSales, a.Ledger.SalesWithoutHold, a.Ledger.LateSales,
		a.Ledger.RejectedEarly, a.DuplicateSeats, a.InventoryDrift, a.InvalidSeats, a.PartialHolds, a.TicketMismatches)
}

func (a *Audit) example(x Example) {
	if len(a.Examples) < 3*maxExamples {
		a.Examples = append(a.Examples, x)
	}
}

type seatKey struct {
	event int64
	seat  int32
}

// RunAudit audits the events in ids (every event when ids is empty).
func RunAudit(ctx context.Context, db ports.DB, d Design, w *World, phase string, ids []int64) (*Audit, error) {
	stmts, err := mustStmts(d.ID, "audit.sql")
	if err != nil {
		return nil, err
	}
	q := catalog.Map(stmts)
	a := &Audit{Phase: phase}
	scope := map[int64]bool{}
	for _, id := range ids {
		scope[id] = true
	}
	in := func(ev int64) bool { return len(scope) == 0 || scope[ev] }

	var t0 time.Time
	if err := db.QueryRow(ctx, "SELECT now()").Scan(&t0); err != nil {
		return nil, err
	}
	tickets := map[seatKey][]int64{}
	ticketsByEvent := map[int64]map[int64]bool{}
	if err := scanAll(ctx, db, q["a_tickets"].SQL, func(rs ports.Rows) error {
		var ev, tid, cust, hold int64
		var seat int32
		if err := rs.Scan(&ev, &seat, &tid, &cust, &hold); err != nil {
			return err
		}
		if !in(ev) {
			return nil
		}
		k := seatKey{ev, seat}
		tickets[k] = append(tickets[k], tid)
		if ticketsByEvent[ev] == nil {
			ticketsByEvent[ev] = map[int64]bool{}
		}
		ticketsByEvent[ev][tid] = true
		a.Tickets++
		return nil
	}); err != nil {
		return nil, fmt.Errorf("a_tickets: %w", err)
	}
	invSold := map[seatKey]bool{}
	if err := scanAll(ctx, db, q["a_inventory_sold"].SQL, func(rs ports.Rows) error {
		var ev, hold int64
		var seat int32
		if err := rs.Scan(&ev, &seat, &hold); err != nil {
			return err
		}
		if in(ev) {
			invSold[seatKey{ev, seat}] = true
		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("a_inventory_sold: %w", err)
	}
	dbHolds := map[seatKey]int64{}
	dbExp := map[seatKey]time.Time{}
	if err := scanAll(ctx, db, q["a_valid_holds"].SQL, func(rs ports.Rows) error {
		var hold, ev int64
		var seat int32
		var exp time.Time
		if err := rs.Scan(&hold, &ev, &seat, &exp); err != nil {
			return err
		}
		if in(ev) {
			dbHolds[seatKey{ev, seat}] = hold
			dbExp[seatKey{ev, seat}] = exp
		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("a_valid_holds: %w", err)
	}
	var t1 time.Time
	if err := db.QueryRow(ctx, "SELECT now()").Scan(&t1); err != nil {
		return nil, err
	}
	for _, name := range []string{"a_inventory_drift", "a_invalid_seats"} {
		if err := scanAll(ctx, db, q[name].SQL, func(rs ports.Rows) error {
			var ev, x int64
			var y any
			if name == "a_inventory_drift" {
				if err := rs.Scan(&ev, &x, &y); err != nil {
					return err
				}
			} else if err := rs.Scan(&ev, &x); err != nil {
				return err
			}
			if !in(ev) {
				return nil
			}
			if name == "a_inventory_drift" {
				a.InventoryDrift++
				a.example(Example{Kind: "inventory rows do not fit the venue", EventID: ev, Detail: fmt.Sprintf("stored %d actual %v", x, y)})
			} else {
				a.InvalidSeats++
				a.example(Example{Kind: "ticket for a seat outside the venue", EventID: ev, SeatID: int32(x)})
			}
			return nil
		}); err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
	}

	// Database-side checks.
	for k, ts := range tickets {
		if len(ts) > 1 {
			a.DuplicateSeats++
			a.example(Example{Kind: "seat with more than one ticket", EventID: k.event, SeatID: k.seat, Detail: fmt.Sprint(ts)})
		}
		if !invSold[k] {
			a.InventoryDrift++
			a.example(Example{Kind: "ticket without a sold seat in the inventory", EventID: k.event, SeatID: k.seat})
		}
	}
	for k := range invSold {
		if len(tickets[k]) == 0 {
			a.InventoryDrift++
			a.example(Example{Kind: "sold seat in the inventory without a ticket", EventID: k.event, SeatID: k.seat})
		}
	}

	// Ledger-side checks, per event.
	var events []*EventLedger
	if len(scope) == 0 {
		events = w.ledger.all()
	} else {
		for id := range scope {
			if el := w.ledger.Event(id); el != nil {
				events = append(events, el)
			}
		}
		sort.Slice(events, func(i, j int) bool { return events[i].ID < events[j].ID })
	}
	edge := time.Second
	for _, el := range events {
		a.Events++
		v, final := el.Evaluate()
		a.Ledger.add(v)
		a.AmbiguousCommits += v.Ambiguous

		// INV-7: ticket ids.
		want := el.Tickets()
		got := ticketsByEvent[el.ID]
		diff := 0
		for tid := range want {
			if !got[tid] {
				diff++
			}
		}
		for tid := range got {
			if _, ok := want[tid]; !ok {
				diff++
			}
		}
		if int64(diff) > v.Ambiguous {
			a.TicketMismatches++
			a.example(Example{Kind: "tickets differ from the sales the harness saw", EventID: el.ID,
				Detail: fmt.Sprintf("%d ticket ids differ (ledger %d, database %d)", diff, len(want), len(got))})
		}

		// INV-5: holds valid at the audit, away from the expiry edge.
		live := liveHoldsAt(final, t1.Add(edge))
		for seat, hold := range live {
			k := seatKey{el.ID, seat}
			if dbHolds[k] != hold {
				a.PartialHolds++
				a.example(Example{Kind: "live hold missing from the database", EventID: el.ID, SeatID: seat, HoldID: hold, Other: dbHolds[k]})
			}
		}
		for k, hold := range dbHolds {
			if k.event != el.ID {
				continue
			}
			if exp := dbExp[k]; exp.Before(t1.Add(edge)) || exp.Before(t0.Add(edge)) {
				a.HoldsSkippedAtEdge++
				continue
			}
			if live[k.seat] != hold {
				a.PartialHolds++
				a.example(Example{Kind: "valid hold the harness never saw granted", EventID: el.ID, SeatID: k.seat, HoldID: hold, Other: live[k.seat]})
			}
		}
	}
	for _, x := range a.Ledger.Examples {
		a.example(x)
	}
	return a, nil
}

func scanAll(ctx context.Context, db ports.DB, sql string, f func(ports.Rows) error) error {
	rows, err := db.Query(ctx, sql)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		if err := f(rows); err != nil {
			return err
		}
	}
	return rows.Err()
}

// eventSold reads one event's sold seats through the design's own audit SQL.
func eventSold(ctx context.Context, db ports.DB, d Design, ev int64) (int64, error) {
	stmts, err := mustStmts(d.ID, "audit.sql")
	if err != nil {
		return 0, err
	}
	var n int64
	err = db.QueryRow(ctx, catalog.Map(stmts)["a_event_sold"].SQL, ev).Scan(&n)
	return n, err
}
