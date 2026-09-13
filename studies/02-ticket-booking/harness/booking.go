package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"adsplatform/core/catalog"
	"adsplatform/ports"
)

// ---------------------------------------------------------------------------
// The world: what the harness knows independently of the database.
//
// Every successful booking and cancellation the harness observes is recorded
// per event. The audit later reconciles that ledger with the tickets that
// actually exist, which catches failures an overbooking check alone cannot:
// a booking reported to the buyer that never persisted (a lost sale), or a
// ticket that exists without any buyer being told (a phantom).
// ---------------------------------------------------------------------------

type EventRef struct {
	ID            int64
	BandID        int64
	Capacity      int
	PriceCents    int64
	Kind          string
	FirstTicketID int64
	InitialSold   int64

	booked    atomic.Int64 // successes the harness saw commit
	cancelled atomic.Int64
	// ambiguous counts commits whose outcome is unknown (connection lost during
	// COMMIT). The reconciliation allows this much disagreement, and no more.
	ambiguous atomic.Int64
}

func (e *EventRef) expectedSold() int64 {
	return e.InitialSold + e.booked.Load() - e.cancelled.Load()
}

type World struct {
	mu        sync.RWMutex
	events    map[int64]*EventRef
	ordered   []*EventRef
	customers int64
	bands     int64

	nextEvent       atomic.Int64
	nextTicket      atomic.Int64
	nextReservation atomic.Int64
}

// NewWorld reflects a freshly loaded dataset. It is rebuilt on every reload, so
// a later phase never inherits an earlier phase's ledger.
func NewWorld(ds *Dataset) *World {
	w := &World{events: map[int64]*EventRef{}, customers: ds.Customers, bands: int64(len(ds.Bands))}
	for i := range ds.Events {
		e := &ds.Events[i]
		ref := &EventRef{
			ID: e.ID, BandID: e.BandID, Capacity: e.Capacity, PriceCents: e.PriceCents,
			Kind: e.Kind, FirstTicketID: e.FirstTicketID, InitialSold: int64(e.InitialSold),
		}
		w.events[e.ID] = ref
		w.ordered = append(w.ordered, ref)
	}
	w.nextEvent.Store(int64(len(ds.Events)))
	// Ticket ids for bookings made by the harness start far above every loaded
	// seat id, so a created ticket can never collide with a pre-created one.
	w.nextTicket.Store(ds.NextTicketID + 1_000_000_000)
	w.nextReservation.Store(0)
	return w
}

func (w *World) event(id int64) *EventRef {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.events[id]
}

func (w *World) add(e *EventRef) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.events[e.ID] = e
	w.ordered = append(w.ordered, e)
}

func (w *World) snapshot() []*EventRef {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return append([]*EventRef(nil), w.ordered...)
}

func (w *World) randomCustomer(r *rand.Rand) int64 { return r.Int63n(w.customers) + 1 }

// ---------------------------------------------------------------------------
// Booker: one design's booking, cancellation, publishing and hold flows.
// ---------------------------------------------------------------------------

type Booker struct {
	db    ports.DB
	d     Design
	st    map[string]catalog.Stmt
	order []catalog.Stmt
	world *World
}

type BookResult struct {
	TicketID int64
	SoldOut  bool
	Retries  int
}

// maxAttempts bounds how long one buyer keeps losing races before the harness
// records an error. It is generous on purpose: an optimistic design near
// sell-out legitimately loses many times in a row, and giving up early would
// turn contention into a false "error" rate.
const maxAttempts = 1000

var (
	errGaveUp = errors.New("gave up after too many lost races")
	// errInvariant is a design's own safety assertion tripping: a row that was
	// locked for us was no longer in the state the lock should have preserved.
	errInvariant = errors.New("invariant violated")
)

func NewBooker(db ports.DB, d Design, world *World) (*Booker, error) {
	src, err := readSQL(d.ID, "writes.sql")
	if err != nil {
		return nil, err
	}
	stmts, err := catalog.Parse(src)
	if err != nil {
		return nil, err
	}
	b := &Booker{db: db, d: d, st: catalog.Map(stmts), order: stmts, world: world}
	for _, name := range b.required() {
		if _, ok := b.st[name]; !ok {
			return nil, fmt.Errorf("design %s: writes.sql lacks %s, required by strategy %s", d.ID, name, d.Strategy)
		}
	}
	return b, nil
}

func (b *Booker) required() []string {
	common := []string{"w_cancel_ticket", "w_edit_event", "w_publish_event"}
	switch b.d.Strategy {
	case LockFirst:
		return append(common, "w_pick_seat", "w_sell_seat")
	case SkipLocked:
		req := append(common, "w_pick_seat", "w_sell_seat", "w_any_available")
		if b.d.CounterColumn != "" {
			req = append(req, "w_count_sale", "w_count_refund")
		}
		return req
	case CAS:
		return append(common, "w_candidate_from", "w_sell_seat")
	case CountCheck:
		req := append(common, "w_seats_left", "w_insert_ticket")
		if b.d.LockEvent {
			req = append(req, "w_lock_event")
		}
		return req
	case Counter:
		return append(common, "w_take_seat", "w_insert_ticket", "w_release_seat")
	case SeatUnique:
		return append(common, "w_highest_seat", "w_claim_seat")
	case Buckets:
		return append(common, "w_pick_bucket", "w_any_remaining", "w_take_from_bucket", "w_insert_ticket",
			"w_pick_refill_bucket", "w_return_to_bucket")
	case SeatPool:
		return append(common, "w_pick_slot", "w_any_slot", "w_consume_slot", "w_insert_ticket", "w_restore_slot")
	case Hold:
		return append(common, "w_take_seat", "w_insert_hold", "w_check_hold", "w_confirm_hold",
			"w_insert_ticket", "w_expired_holds", "w_release_hold", "w_release_seat")
	}
	return common
}

func (b *Booker) args(name string, vals map[string]any) []any {
	st := b.st[name]
	a, err := catalog.Bind(st.Params, vals)
	if err != nil {
		// A missing binding is a harness bug, never a runtime condition.
		panic(fmt.Sprintf("%s.%s: %v", b.d.ID, name, err))
	}
	return a
}

func (b *Booker) exec(ctx context.Context, q ports.Queryer, name string, vals map[string]any) (int64, error) {
	return q.Exec(ctx, b.st[name].SQL, b.args(name, vals)...)
}

func (b *Booker) row(ctx context.Context, q ports.Queryer, name string, vals map[string]any) ports.Row {
	return q.QueryRow(ctx, b.st[name].SQL, b.args(name, vals)...)
}

func (b *Booker) retryable(err error) bool { return b.db.Classify(err) == ports.ErrRetryable }

// commit records an ambiguous outcome against the event before returning it.
func (b *Booker) commit(ctx context.Context, tx ports.Tx, ev *EventRef) error {
	err := tx.Commit(ctx)
	if err != nil && b.db.Classify(err) == ports.ErrAmbiguousCommit {
		ev.ambiguous.Add(1)
	}
	return err
}

// Book sells one seat of ev to customer, or reports the event sold out.
func (b *Booker) Book(ctx context.Context, ev *EventRef, customer int64, r *rand.Rand) (BookResult, error) {
	var res BookResult
	var err error
	switch b.d.Strategy {
	case LockFirst, SkipLocked:
		res, err = b.bookLocked(ctx, ev, customer)
	case CAS:
		res, err = b.bookCAS(ctx, ev, customer, r)
	case CountCheck:
		res, err = b.bookCount(ctx, ev, customer)
	case Counter:
		res, err = b.bookCounter(ctx, ev, customer)
	case SeatUnique:
		res, err = b.bookSeatUnique(ctx, ev, customer)
	case Buckets:
		res, err = b.bookBuckets(ctx, ev, customer)
	case SeatPool:
		res, err = b.bookSeatPool(ctx, ev, customer)
	case Hold:
		res, err = b.bookHoldImmediate(ctx, ev, customer)
	default:
		err = fmt.Errorf("strategy %q not implemented", b.d.Strategy)
	}
	if err == nil && !res.SoldOut {
		ev.booked.Add(1)
	}
	return res, err
}

type attemptDeadlineKey struct{}

// WithAttemptDeadline marks a context so that no NEW attempt starts after t.
// An attempt already running finishes normally, so a timeout never manufactures
// a commit of unknown outcome. Without this, a race's timeout stopped new buyers
// but not a buyer already retrying inside one booking: on YugabyteDB at
// SERIALIZABLE one such buyer kept a 10-seat race alive for minutes.
func WithAttemptDeadline(ctx context.Context, t time.Time) context.Context {
	return context.WithValue(ctx, attemptDeadlineKey{}, t)
}

var errAttemptDeadline = errors.New("attempt deadline passed before the booking completed")

// loop runs attempt until it reports done, retrying retryable engine errors.
// attempt returns retry=true to start over after a lost race.
func (b *Booker) loop(ctx context.Context, res *BookResult, attempt func() (retry bool, err error)) error {
	deadline, hasDeadline := ctx.Value(attemptDeadlineKey{}).(time.Time)
	for i := 0; i < maxAttempts; i++ {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		if hasDeadline && i > 0 && time.Now().After(deadline) {
			return errAttemptDeadline
		}
		retry, err := attempt()
		if err != nil {
			if b.retryable(err) {
				res.Retries++
				continue
			}
			return err
		}
		if !retry {
			return nil
		}
		res.Retries++
	}
	return errGaveUp
}

func (b *Booker) bookLocked(ctx context.Context, ev *EventRef, customer int64) (BookResult, error) {
	var res BookResult
	vals := map[string]any{"event_id": ev.ID, "customer_id": customer}
	err := b.loop(ctx, &res, func() (bool, error) {
		tx, err := b.db.Begin(ctx, b.d.Isolation)
		if err != nil {
			return false, err
		}
		defer tx.Rollback(ctx)

		var tid int64
		err = b.row(ctx, tx, "w_pick_seat", vals).Scan(&tid)
		if errors.Is(err, ports.ErrNoRows) {
			_ = tx.Rollback(ctx)
			if b.d.Strategy == LockFirst {
				// A blocking scan waits for every locked seat and re-checks it, so
				// an empty result really is sold out.
				res.SoldOut = true
				return false, nil
			}
			// SKIP LOCKED: empty may only mean "all remaining seats are locked".
			var any bool
			if err := b.row(ctx, b.db, "w_any_available", vals).Scan(&any); err != nil {
				return false, err
			}
			if !any {
				res.SoldOut = true
				return false, nil
			}
			return true, nil
		}
		if err != nil {
			return false, err
		}
		vals["ticket_id"] = tid
		n, err := b.exec(ctx, tx, "w_sell_seat", vals)
		if err != nil {
			return false, err
		}
		if n != 1 {
			return false, fmt.Errorf("%w: locked seat %d was no longer available", errInvariant, tid)
		}
		if b.d.CounterColumn != "" {
			if _, err := b.exec(ctx, tx, "w_count_sale", vals); err != nil {
				return false, err
			}
		}
		if err := b.commit(ctx, tx, ev); err != nil {
			return false, err
		}
		res.TicketID = tid
		return false, nil
	})
	return res, err
}

func (b *Booker) bookCAS(ctx context.Context, ev *EventRef, customer int64, r *rand.Rand) (BookResult, error) {
	var res BookResult
	vals := map[string]any{"event_id": ev.ID, "customer_id": customer}
	err := b.loop(ctx, &res, func() (bool, error) {
		start := 1 + r.Intn(ev.Capacity)
		var tid int64
		vals["start_seat"] = start
		err := b.row(ctx, b.db, "w_candidate_from", vals).Scan(&tid)
		if errors.Is(err, ports.ErrNoRows) && start > 1 {
			vals["start_seat"] = 1
			err = b.row(ctx, b.db, "w_candidate_from", vals).Scan(&tid)
		}
		if errors.Is(err, ports.ErrNoRows) {
			res.SoldOut = true
			return false, nil
		}
		if err != nil {
			return false, err
		}
		vals["ticket_id"] = tid
		n, err := b.exec(ctx, b.db, "w_sell_seat", vals)
		if err != nil {
			return false, err
		}
		if n == 0 {
			return true, nil // someone else sold it between our read and our write
		}
		res.TicketID = tid
		return false, nil
	})
	return res, err
}

func (b *Booker) insertVals(ev *EventRef, customer int64) map[string]any {
	return map[string]any{
		"ticket_id": b.world.nextTicket.Add(1), "event_id": ev.ID, "seat_no": nil,
		"customer_id": customer, "price_cents": ev.PriceCents, "reservation_id": nil,
	}
}

func (b *Booker) bookCount(ctx context.Context, ev *EventRef, customer int64) (BookResult, error) {
	var res BookResult
	err := b.loop(ctx, &res, func() (bool, error) {
		vals := b.insertVals(ev, customer)
		tx, err := b.db.Begin(ctx, b.d.Isolation)
		if err != nil {
			return false, err
		}
		defer tx.Rollback(ctx)
		if b.d.LockEvent {
			var id int64
			if err := b.row(ctx, tx, "w_lock_event", vals).Scan(&id); err != nil {
				return false, err
			}
		}
		var left int64
		if err := b.row(ctx, tx, "w_seats_left", vals).Scan(&left); err != nil {
			return false, err
		}
		if left <= 0 {
			// Committed rather than rolled back: under SERIALIZABLE even a
			// read-only conclusion is only valid if the transaction commits.
			if err := tx.Commit(ctx); err != nil {
				return false, err
			}
			res.SoldOut = true
			return false, nil
		}
		if _, err := b.exec(ctx, tx, "w_insert_ticket", vals); err != nil {
			return false, err
		}
		if err := b.commit(ctx, tx, ev); err != nil {
			return false, err
		}
		res.TicketID = vals["ticket_id"].(int64)
		return false, nil
	})
	return res, err
}

func (b *Booker) bookCounter(ctx context.Context, ev *EventRef, customer int64) (BookResult, error) {
	var res BookResult
	err := b.loop(ctx, &res, func() (bool, error) {
		vals := b.insertVals(ev, customer)
		tx, err := b.db.Begin(ctx, b.d.Isolation)
		if err != nil {
			return false, err
		}
		defer tx.Rollback(ctx)
		var n int64
		err = b.row(ctx, tx, "w_take_seat", vals).Scan(&n)
		if errors.Is(err, ports.ErrNoRows) {
			res.SoldOut = true
			return false, nil
		}
		if err != nil {
			return false, err
		}
		if _, err := b.exec(ctx, tx, "w_insert_ticket", vals); err != nil {
			return false, err
		}
		if err := b.commit(ctx, tx, ev); err != nil {
			return false, err
		}
		res.TicketID = vals["ticket_id"].(int64)
		return false, nil
	})
	return res, err
}

func (b *Booker) bookSeatUnique(ctx context.Context, ev *EventRef, customer int64) (BookResult, error) {
	var res BookResult
	err := b.loop(ctx, &res, func() (bool, error) {
		vals := b.insertVals(ev, customer)
		var highest int
		if err := b.row(ctx, b.db, "w_highest_seat", vals).Scan(&highest); err != nil {
			return false, err
		}
		if highest >= ev.Capacity {
			res.SoldOut = true
			return false, nil
		}
		vals["seat_no"] = highest + 1
		n, err := b.exec(ctx, b.db, "w_claim_seat", vals)
		switch b.db.Classify(err) {
		case ports.ErrUniqueViolation:
			return true, nil
		case ports.ErrCheckViolation:
			// Seat past the venue: the constraint, not the code, said sold out.
			res.SoldOut = true
			return false, nil
		}
		if err != nil {
			return false, err
		}
		if n == 0 {
			return true, nil // ON CONFLICT DO NOTHING: another buyer took that seat number
		}
		res.TicketID = vals["ticket_id"].(int64)
		return false, nil
	})
	return res, err
}

func (b *Booker) bookBuckets(ctx context.Context, ev *EventRef, customer int64) (BookResult, error) {
	var res BookResult
	err := b.loop(ctx, &res, func() (bool, error) {
		vals := b.insertVals(ev, customer)
		tx, err := b.db.Begin(ctx, b.d.Isolation)
		if err != nil {
			return false, err
		}
		defer tx.Rollback(ctx)
		var bucket int
		err = b.row(ctx, tx, "w_pick_bucket", vals).Scan(&bucket)
		if errors.Is(err, ports.ErrNoRows) {
			_ = tx.Rollback(ctx)
			var any bool
			if err := b.row(ctx, b.db, "w_any_remaining", vals).Scan(&any); err != nil {
				return false, err
			}
			if !any {
				res.SoldOut = true
				return false, nil
			}
			return true, nil
		}
		if err != nil {
			return false, err
		}
		vals["bucket_no"] = bucket
		n, err := b.exec(ctx, tx, "w_take_from_bucket", vals)
		if err != nil {
			return false, err
		}
		if n != 1 {
			return false, fmt.Errorf("%w: locked bucket %d of event %d had no seat", errInvariant, bucket, ev.ID)
		}
		if _, err := b.exec(ctx, tx, "w_insert_ticket", vals); err != nil {
			return false, err
		}
		if err := b.commit(ctx, tx, ev); err != nil {
			return false, err
		}
		res.TicketID = vals["ticket_id"].(int64)
		return false, nil
	})
	return res, err
}

func (b *Booker) bookSeatPool(ctx context.Context, ev *EventRef, customer int64) (BookResult, error) {
	var res BookResult
	err := b.loop(ctx, &res, func() (bool, error) {
		vals := b.insertVals(ev, customer)
		tx, err := b.db.Begin(ctx, b.d.Isolation)
		if err != nil {
			return false, err
		}
		defer tx.Rollback(ctx)
		var seat int
		err = b.row(ctx, tx, "w_pick_slot", vals).Scan(&seat)
		if errors.Is(err, ports.ErrNoRows) {
			_ = tx.Rollback(ctx)
			var any bool
			if err := b.row(ctx, b.db, "w_any_slot", vals).Scan(&any); err != nil {
				return false, err
			}
			if !any {
				res.SoldOut = true
				return false, nil
			}
			return true, nil
		}
		if err != nil {
			return false, err
		}
		vals["seat_no"] = seat
		n, err := b.exec(ctx, tx, "w_consume_slot", vals)
		if err != nil {
			return false, err
		}
		if n != 1 {
			return false, fmt.Errorf("%w: locked slot %d/%d was already consumed", errInvariant, ev.ID, seat)
		}
		if _, err := b.exec(ctx, tx, "w_insert_ticket", vals); err != nil {
			return false, err
		}
		if err := b.commit(ctx, tx, ev); err != nil {
			return false, err
		}
		res.TicketID = vals["ticket_id"].(int64)
		return false, nil
	})
	return res, err
}

// ---------------------------------------------------------------------------
// Holds
// ---------------------------------------------------------------------------

// Hold takes a seat into a basket for holdMS milliseconds.
func (b *Booker) Hold(ctx context.Context, ev *EventRef, customer int64, holdMS int) (resID int64, soldOut bool, retries int, err error) {
	var res BookResult
	err = b.loop(ctx, &res, func() (bool, error) {
		id := b.world.nextReservation.Add(1)
		vals := map[string]any{"event_id": ev.ID, "reservation_id": id, "customer_id": customer, "hold_ms": holdMS}
		tx, err := b.db.Begin(ctx, b.d.Isolation)
		if err != nil {
			return false, err
		}
		defer tx.Rollback(ctx)
		var n int64
		err = b.row(ctx, tx, "w_take_seat", vals).Scan(&n)
		if errors.Is(err, ports.ErrNoRows) {
			res.SoldOut = true
			return false, nil
		}
		if err != nil {
			return false, err
		}
		if _, err := b.exec(ctx, tx, "w_insert_hold", vals); err != nil {
			return false, err
		}
		if err := tx.Commit(ctx); err != nil {
			return false, err
		}
		resID = id
		return false, nil
	})
	return resID, res.SoldOut, res.Retries, err
}

// CheckHold is the friendly pre-payment check. It is identical in H0 and H1.
func (b *Booker) CheckHold(ctx context.Context, resID int64) (bool, error) {
	var valid bool
	err := b.row(ctx, b.db, "w_check_hold", map[string]any{"reservation_id": resID}).Scan(&valid)
	if errors.Is(err, ports.ErrNoRows) {
		return false, nil
	}
	return valid, err
}

// Confirm turns a hold into a ticket. lost=true means the confirmation matched
// no hold (H1 only: the hold expired or was released) and no ticket was issued.
// The flow is identical for H0 and H1; the difference is entirely in the SQL of
// w_confirm_hold.
func (b *Booker) Confirm(ctx context.Context, ev *EventRef, resID, customer int64) (ticketID int64, lost bool, retries int, err error) {
	var res BookResult
	err = b.loop(ctx, &res, func() (bool, error) {
		vals := b.insertVals(ev, customer)
		vals["reservation_id"] = resID
		tx, err := b.db.Begin(ctx, b.d.Isolation)
		if err != nil {
			return false, err
		}
		defer tx.Rollback(ctx)
		n, err := b.exec(ctx, tx, "w_confirm_hold", vals)
		if err != nil {
			return false, err
		}
		if n == 0 {
			lost = true
			return false, nil
		}
		if _, err := b.exec(ctx, tx, "w_insert_ticket", vals); err != nil {
			return false, err
		}
		if err := b.commit(ctx, tx, ev); err != nil {
			return false, err
		}
		ticketID = vals["ticket_id"].(int64)
		return false, nil
	})
	if err == nil && !lost {
		ev.booked.Add(1)
	}
	return ticketID, lost, res.Retries, err
}

// bookHoldImmediate is a purchase with no basket time: hold, then confirm at
// once. Used wherever every design must answer the same "book a seat" request.
func (b *Booker) bookHoldImmediate(ctx context.Context, ev *EventRef, customer int64) (BookResult, error) {
	var res BookResult
	resID, soldOut, retries, err := b.Hold(ctx, ev, customer, 60_000)
	res.Retries += retries
	if err != nil || soldOut {
		res.SoldOut = soldOut
		return res, err
	}
	tid, lost, retries, err := b.Confirm(ctx, ev, resID, customer)
	res.Retries += retries
	if err != nil {
		return res, err
	}
	if lost {
		return res, fmt.Errorf("%w: a 60 s hold was lost before its immediate confirmation", errInvariant)
	}
	// Confirm already counted the booking; Book must not count it twice.
	ev.booked.Add(-1)
	res.TicketID = tid
	return res, nil
}

// Sweep releases up to batch expired holds and returns how many seats it gave back.
func (b *Booker) Sweep(ctx context.Context, batch int) (int, error) {
	tx, err := b.db.Begin(ctx, ports.ReadCommitted)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)
	rows, err := tx.Query(ctx, b.st["w_expired_holds"].SQL, b.args("w_expired_holds", map[string]any{"batch": batch})...)
	if err != nil {
		return 0, err
	}
	type hold struct{ id, event int64 }
	var expired []hold
	for rows.Next() {
		var h hold
		if err := rows.Scan(&h.id, &h.event); err != nil {
			rows.Close()
			return 0, err
		}
		expired = append(expired, h)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return 0, err
	}
	released := 0
	for _, h := range expired {
		n, err := b.exec(ctx, tx, "w_release_hold", map[string]any{"reservation_id": h.id})
		if err != nil {
			return 0, err
		}
		if n == 1 {
			if _, err := b.exec(ctx, tx, "w_release_seat", map[string]any{"event_id": h.event}); err != nil {
				return 0, err
			}
			released++
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return 0, err
	}
	return released, nil
}

// ---------------------------------------------------------------------------
// Cancellation, publishing, editing
// ---------------------------------------------------------------------------

// Cancel refunds a sold ticket. ok=false means the ticket was not sold (already
// cancelled), which is not an error.
func (b *Booker) Cancel(ctx context.Context, ticketID int64) (ok bool, retries int, err error) {
	var res BookResult
	var ev *EventRef
	err = b.loop(ctx, &res, func() (bool, error) {
		vals := map[string]any{"ticket_id": ticketID}
		tx, err := b.db.Begin(ctx, ports.ReadCommitted)
		if err != nil {
			return false, err
		}
		defer tx.Rollback(ctx)
		var eventID int64
		var seat *int32
		err = b.row(ctx, tx, "w_cancel_ticket", vals).Scan(&eventID, &seat)
		if errors.Is(err, ports.ErrNoRows) {
			return false, nil
		}
		if err != nil {
			return false, err
		}
		ev = b.world.event(eventID)
		vals["event_id"] = eventID
		switch {
		case b.d.Strategy == SkipLocked && b.d.CounterColumn != "":
			if _, err := b.exec(ctx, tx, "w_count_refund", vals); err != nil {
				return false, err
			}
		case b.d.Strategy == Counter || b.d.Strategy == Hold:
			if _, err := b.exec(ctx, tx, "w_release_seat", vals); err != nil {
				return false, err
			}
		case b.d.Strategy == Buckets:
			var bucket int
			err := b.row(ctx, tx, "w_pick_refill_bucket", vals).Scan(&bucket)
			if errors.Is(err, ports.ErrNoRows) {
				return true, nil // every bucket with room is locked right now
			}
			if err != nil {
				return false, err
			}
			vals["bucket_no"] = bucket
			n, err := b.exec(ctx, tx, "w_return_to_bucket", vals)
			if err != nil {
				return false, err
			}
			if n != 1 {
				return false, fmt.Errorf("%w: locked refill bucket %d was full", errInvariant, bucket)
			}
		case b.d.Strategy == SeatPool:
			if seat == nil {
				return false, fmt.Errorf("%w: pooled ticket %d has no seat", errInvariant, ticketID)
			}
			vals["seat_no"] = *seat
			if _, err := b.exec(ctx, tx, "w_restore_slot", vals); err != nil {
				return false, err
			}
		}
		if err := tx.Commit(ctx); err != nil {
			if ev != nil && b.db.Classify(err) == ports.ErrAmbiguousCommit {
				ev.ambiguous.Add(1)
			}
			return false, err
		}
		ok = true
		return false, nil
	})
	if err == nil && ok && ev != nil {
		ev.cancelled.Add(1)
	}
	return ok, res.Retries, err
}

// Publish creates an event with its inventory, in one transaction, by running
// every w_publish_* statement of the design in file order. What that costs is
// the pre-creation question in its purest form: one row, or one row per seat.
func (b *Booker) Publish(ctx context.Context, bandID int64, capacity int, r *rand.Rand) (*EventRef, error) {
	ev := &EventRef{
		ID: b.world.nextEvent.Add(1), BandID: bandID, Capacity: capacity,
		PriceCents: int64(2000 + 500*r.Intn(30)), Kind: "published",
	}
	ev.FirstTicketID = b.world.nextTicket.Add(int64(capacity)) - int64(capacity) + 1
	vals := map[string]any{
		"event_id": ev.ID, "band_id": bandID, "name": fmt.Sprintf("published show %d", ev.ID),
		"venue": venueFor(capacity), "starts_at": showEpoch.Add(time.Duration(r.Intn(365*24)) * time.Hour),
		"capacity": capacity, "price_cents": ev.PriceCents, "description": "Just announced.",
		"first_ticket_id": ev.FirstTicketID,
	}
	var res BookResult
	err := b.loop(ctx, &res, func() (bool, error) {
		tx, err := b.db.Begin(ctx, ports.ReadCommitted)
		if err != nil {
			return false, err
		}
		defer tx.Rollback(ctx)
		for _, st := range b.order {
			if !strings.HasPrefix(st.Name, "w_publish_") {
				continue
			}
			if _, err := b.exec(ctx, tx, st.Name, vals); err != nil {
				return false, fmt.Errorf("%s: %w", st.Name, err)
			}
		}
		return false, tx.Commit(ctx)
	})
	if err != nil {
		return nil, err
	}
	b.world.add(ev)
	return ev, nil
}

// Edit is the organiser changing the event page mid-sale.
func (b *Booker) Edit(ctx context.Context, ev *EventRef, r *rand.Rand) error {
	_, err := b.exec(ctx, b.db, "w_edit_event", map[string]any{
		"event_id": ev.ID, "description": fmt.Sprintf("Updated set list, revision %d.", r.Intn(1_000_000)),
	})
	return err
}
