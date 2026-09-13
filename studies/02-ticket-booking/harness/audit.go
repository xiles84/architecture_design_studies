package main

import (
	"context"
	"fmt"
	"sort"

	"adsplatform/core/catalog"
	"adsplatform/ports"
)

// ---------------------------------------------------------------------------
// The overbooking audit.
//
// Run after every phase that writes, on the state that phase produced (study 01
// learned that an audit placed after a later phase can be healed into passing).
// It recomputes, from each design's own tables, four independent things:
//
//	overbooked     events holding more sold tickets than seats
//	duplicates     a seat number sold twice
//	derived drift  a counter, bucket set, slot pool or seat frontier that
//	               disagrees with the tickets it summarises
//	reconciliation the harness's own ledger (loaded + bookings it saw commit
//	               - cancellations it saw commit) against the tickets that exist:
//	               a mismatch is a sale the buyer was told about that did not
//	               persist, or a ticket nobody was told about
//
// It records examples, not just counts: in study 01 the examples were what
// turned "the cache drifted" into two distinct, fixable bugs.
// ---------------------------------------------------------------------------

type EventViolation struct {
	EventID  int64  `json:"event_id"`
	Capacity int64  `json:"capacity,omitempty"`
	Stored   int64  `json:"stored,omitempty"`
	Actual   int64  `json:"actual"`
	Expected int64  `json:"expected,omitempty"`
	Note     string `json:"note,omitempty"`
}

type Audit struct {
	Phase         string `json:"phase"`
	EventsChecked int64  `json:"events_checked"`
	SeatsSold     int64  `json:"seats_sold"`

	OverbookedEvents int64            `json:"overbooked_events"`
	OverbookedSeats  int64            `json:"overbooked_seats"`
	Overbooked       []EventViolation `json:"overbooked_examples,omitempty"`

	DuplicateSeats int64            `json:"duplicate_seats"`
	Duplicates     []EventViolation `json:"duplicate_examples,omitempty"`

	DriftEvents int64            `json:"derived_drift_events"`
	Drift       []EventViolation `json:"derived_drift_examples,omitempty"`

	// LedgerMismatches counts events whose tickets differ from what the harness
	// saw committed by more than the ambiguous commits recorded for them.
	LedgerMismatches int64            `json:"ledger_mismatch_events"`
	Ledger           []EventViolation `json:"ledger_examples,omitempty"`
	AmbiguousCommits int64            `json:"ambiguous_commits,omitempty"`
}

// Violations is the number that gates timing. Derived drift counts: a counter
// that disagrees with its tickets will oversell or undersell on its next use.
func (a *Audit) Violations() int64 {
	if a == nil {
		return 0
	}
	return a.OverbookedEvents + a.DuplicateSeats + a.DriftEvents + a.LedgerMismatches
}

func (a *Audit) String() string {
	if a.Violations() == 0 {
		return fmt.Sprintf("consistent (%d events, %d seats sold)", a.EventsChecked, a.SeatsSold)
	}
	return fmt.Sprintf("VIOLATIONS — %d overbooked events (+%d seats), %d duplicate seats, %d derived-state drifts, %d ledger mismatches",
		a.OverbookedEvents, a.OverbookedSeats, a.DuplicateSeats, a.DriftEvents, a.LedgerMismatches)
}

const auditExamples = 5

func RunAudit(ctx context.Context, db ports.DB, d Design, world *World, phase string) (*Audit, error) {
	stmts, err := mustStmts(d.ID, "audit.sql")
	if err != nil {
		return nil, err
	}
	q := catalog.Map(stmts)
	a := &Audit{Phase: phase}

	rows, err := db.Query(ctx, q["a_sold_per_event"].SQL)
	if err != nil {
		return nil, fmt.Errorf("a_sold_per_event: %w", err)
	}
	type evRow struct{ id, capacity, sold int64 }
	var evs []evRow
	for rows.Next() {
		var r evRow
		var capacity int32
		if err := rows.Scan(&r.id, &capacity, &r.sold); err != nil {
			rows.Close()
			return nil, err
		}
		r.capacity = int64(capacity)
		evs = append(evs, r)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return nil, err
	}
	sort.Slice(evs, func(i, j int) bool { return evs[i].id < evs[j].id })

	for _, r := range evs {
		a.EventsChecked++
		a.SeatsSold += r.sold
		if r.sold > r.capacity {
			a.OverbookedEvents++
			a.OverbookedSeats += r.sold - r.capacity
			if len(a.Overbooked) < auditExamples {
				a.Overbooked = append(a.Overbooked, EventViolation{EventID: r.id, Capacity: r.capacity, Actual: r.sold})
			}
		}
		if world == nil {
			continue
		}
		ref := world.event(r.id)
		if ref == nil {
			a.LedgerMismatches++
			if len(a.Ledger) < auditExamples {
				a.Ledger = append(a.Ledger, EventViolation{EventID: r.id, Actual: r.sold, Note: "event unknown to the harness"})
			}
			continue
		}
		amb := ref.ambiguous.Load()
		a.AmbiguousCommits += amb
		exp := ref.expectedSold()
		if diff := r.sold - exp; diff > amb || -diff > amb {
			a.LedgerMismatches++
			if len(a.Ledger) < auditExamples {
				note := "tickets exist that no buyer was told about"
				if diff < 0 {
					note = "buyers were told they had tickets that do not exist"
				}
				a.Ledger = append(a.Ledger, EventViolation{EventID: r.id, Capacity: r.capacity, Actual: r.sold, Expected: exp, Note: note})
			}
		}
	}

	if err := collectViolations(ctx, db, q["a_duplicate_seats"].SQL, func(ev, seat, n int64) {
		a.DuplicateSeats++
		if len(a.Duplicates) < auditExamples {
			a.Duplicates = append(a.Duplicates, EventViolation{EventID: ev, Stored: seat, Actual: n, Note: "seat_no in stored, tickets in actual"})
		}
	}); err != nil {
		return nil, fmt.Errorf("a_duplicate_seats: %w", err)
	}
	if err := collectViolations(ctx, db, q["a_derived_drift"].SQL, func(ev, stored, actual int64) {
		a.DriftEvents++
		if len(a.Drift) < auditExamples {
			a.Drift = append(a.Drift, EventViolation{EventID: ev, Stored: stored, Actual: actual})
		}
	}); err != nil {
		return nil, fmt.Errorf("a_derived_drift: %w", err)
	}
	return a, nil
}

func collectViolations(ctx context.Context, db ports.DB, sql string, fn func(a, b, c int64)) error {
	rows, err := db.Query(ctx, sql)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var x, y, z int64
		if err := rows.Scan(&x, &y, &z); err != nil {
			return err
		}
		fn(x, y, z)
	}
	return rows.Err()
}

// eventSold reads one event's sold count through the design's own audit SQL.
func eventSold(ctx context.Context, db ports.DB, auditSQL string, eventID int64) (int64, error) {
	var n int64
	err := db.QueryRow(ctx, auditSQL, eventID).Scan(&n)
	return n, err
}
