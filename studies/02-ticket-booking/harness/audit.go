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

// ---------------------------------------------------------------------------
// Ledger reconciliation audit (X1 only, REPORTS.md section 4 / AM-02.3).
//
// Reconciliation against what the harness saw commit is already covered by
// RunAudit's existing ledger check (ticket.status='sold' counts against
// world.expectedSold()); this audit ties sale_event to that SAME ticket
// state, per seat, so the composition of the two covers both halves of the
// requirement without duplicating the harness-observation machinery.
// ---------------------------------------------------------------------------

type LedgerMismatch struct {
	EventID    int64 `json:"event_id"`
	SeatNo     int64 `json:"seat_no"`
	LedgerNet  int64 `json:"ledger_net"`
	ActualSold int64 `json:"actual_sold"`
}

// LedgerAttributionIssue is one row of a_ledger_attribution: a seat whose
// ledger history does not read as sold/cancelled/sold/... with each refund
// naming the buyer it reverses (EH-02 AM-03.1 -- the net-count check above
// cannot see a wrong customer, a missing predecessor, or two sales in a row).
type LedgerAttributionIssue struct {
	EventID int64  `json:"event_id"`
	SeatNo  int64  `json:"seat_no"`
	Problem string `json:"problem"`
}

type LedgerAudit struct {
	Phase                 string                   `json:"phase"`
	SeatsChecked          int64                    `json:"seats_checked"`
	Mismatches            int64                    `json:"mismatches"`
	MismatchExamples      []LedgerMismatch         `json:"mismatch_examples,omitempty"`
	AttributionMismatches int64                    `json:"attribution_mismatches"`
	AttributionExamples   []LedgerAttributionIssue `json:"attribution_examples,omitempty"`
}

// Violations is the count that gates correctness: both the net-count check
// and the attribution check must be zero for the ledger to be trusted.
func (a *LedgerAudit) Violations() int64 {
	if a == nil {
		return 0
	}
	return a.Mismatches + a.AttributionMismatches
}

func (a *LedgerAudit) String() string {
	if a == nil {
		return "consistent (0 seats checked)"
	}
	if a.Violations() == 0 {
		return fmt.Sprintf("consistent (%d seats checked)", a.SeatsChecked)
	}
	return fmt.Sprintf("INCONSISTENT — %d net mismatches, %d attribution problems (of %d seats checked)",
		a.Mismatches, a.AttributionMismatches, a.SeatsChecked)
}

func RunLedgerAudit(ctx context.Context, db ports.DB, d Design, phase string) (*LedgerAudit, error) {
	if !d.Ledger {
		return nil, nil
	}
	stmts, err := mustStmts(d.ID, "audit.sql")
	if err != nil {
		return nil, err
	}
	q := catalog.Map(stmts)
	a := &LedgerAudit{Phase: phase}

	rows, err := db.Query(ctx, q["a_ledger_mismatches"].SQL)
	if err != nil {
		return nil, fmt.Errorf("a_ledger_mismatches: %w", err)
	}
	for rows.Next() {
		var m LedgerMismatch
		var seat int32
		if err := rows.Scan(&m.EventID, &seat, &m.LedgerNet, &m.ActualSold); err != nil {
			rows.Close()
			return nil, err
		}
		m.SeatNo = int64(seat)
		a.Mismatches++
		if len(a.MismatchExamples) < auditExamples {
			a.MismatchExamples = append(a.MismatchExamples, m)
		}
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return nil, err
	}
	rows.Close()

	attrRows, err := db.Query(ctx, q["a_ledger_attribution"].SQL)
	if err != nil {
		return nil, fmt.Errorf("a_ledger_attribution: %w", err)
	}
	for attrRows.Next() {
		var iss LedgerAttributionIssue
		var seat int32
		if err := attrRows.Scan(&iss.EventID, &seat, &iss.Problem); err != nil {
			attrRows.Close()
			return nil, err
		}
		iss.SeatNo = int64(seat)
		a.AttributionMismatches++
		if len(a.AttributionExamples) < auditExamples {
			a.AttributionExamples = append(a.AttributionExamples, iss)
		}
	}
	if err := attrRows.Err(); err != nil {
		attrRows.Close()
		return nil, err
	}
	attrRows.Close()

	// SeatsChecked: every ticket row (the reconciliation query's left side).
	if err := db.QueryRow(ctx, "SELECT COUNT(*) FROM ticket").Scan(&a.SeatsChecked); err != nil {
		return nil, fmt.Errorf("ticket count: %w", err)
	}
	return a, nil
}

// ---------------------------------------------------------------------------
// Ledger fault injection -- TESTING ONLY (EH-02 AM-03.3). Proves the audit
// above actually fires, by corrupting a real ledger row directly in the
// database rather than through any design statement. The fault SQL lives
// here, never in a design's own catalogue, so it is invisible to the
// SQL-binding test and to EXPLAIN. Callers validate that this only runs with
// -cmd verify on a Ledger design (main.go); RunLedgerAudit itself already
// no-ops for a non-Ledger design.
// ---------------------------------------------------------------------------

// injectLedgerFault corrupts the ledger row for the first sold seat of the
// dataset's busiest catalogue event (chooseKeys' own choice of key, so the
// fault lands on a seat the correctness gate already exercises).
//
//   - "drop-sale" deletes that seat's sold ledger row: the seat is still
//     sold, but the ledger no longer says so -- a_ledger_mismatches must
//     report exactly one net mismatch.
//   - "wrong-customer" changes that row's customer_id: the seat's sale is
//     still recorded, but to the wrong buyer -- a_ledger_attribution must
//     report exactly one "live sale disagrees with its newest ledger row".
func injectLedgerFault(ctx context.Context, db ports.DB, ds *Dataset, fault string) error {
	k := chooseKeys(ds)
	if k.soldTicket == nil {
		return fmt.Errorf("dataset has no sold ticket to corrupt")
	}
	switch fault {
	case "drop-sale":
		n, err := db.Exec(ctx, "DELETE FROM sale_event WHERE event_id = $1 AND seat_no = $2 AND kind = 'sold'",
			k.soldTicket.EventID, k.soldTicket.SeatNo)
		if err != nil {
			return err
		}
		if n == 0 {
			return fmt.Errorf("no sold ledger row for event %d seat %d", k.soldTicket.EventID, k.soldTicket.SeatNo)
		}
	case "wrong-customer":
		n, err := db.Exec(ctx, "UPDATE sale_event SET customer_id = customer_id + 1 WHERE event_id = $1 AND seat_no = $2 AND kind = 'sold'",
			k.soldTicket.EventID, k.soldTicket.SeatNo)
		if err != nil {
			return err
		}
		if n == 0 {
			return fmt.Errorf("no sold ledger row for event %d seat %d", k.soldTicket.EventID, k.soldTicket.SeatNo)
		}
	default:
		return fmt.Errorf("unknown ledger fault %q", fault)
	}
	return nil
}
