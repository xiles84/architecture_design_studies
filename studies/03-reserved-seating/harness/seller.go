package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"adsplatform/core/catalog"
	"adsplatform/ports"
)

// ---------------------------------------------------------------------------
// The world: the dataset, the ledger and the id generators of one load.
// ---------------------------------------------------------------------------

type World struct {
	ds     *Dataset
	ledger *Ledger

	mu        sync.RWMutex
	published map[int64]*Event

	nextHold   atomic.Int64
	nextTicket atomic.Int64
	nextEvent  atomic.Int64
}

// NewWorld reflects a freshly loaded dataset; it is rebuilt on every reload.
func NewWorld(ds *Dataset, loadNow time.Time, guard time.Duration) *World {
	w := &World{ds: ds, ledger: NewLedger(ds, loadNow, guard), published: map[int64]*Event{}}
	// Ids allocated by the harness start far above every loaded id.
	w.nextHold.Store(1_000_000_000)
	w.nextTicket.Store(1_000_000_000)
	w.nextEvent.Store(int64(len(ds.Events)) + 1_000_000)
	return w
}

func (w *World) event(id int64) *Event {
	if e := w.ds.event(id); e != nil {
		return e
	}
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.published[id]
}

func (w *World) randomCustomer(r *rand.Rand) int64 { return r.Int63n(w.ds.Customers) + 1 }

// Block is a buyer's choice: adjacent seats in one row of one section.
type Block struct {
	Event   *Event
	Section int32
	Seats   []int32 // ascending
}

// ---------------------------------------------------------------------------
// Seller: one design's hold, checkout, sweeper, refund and publishing flows.
// ---------------------------------------------------------------------------

type Seller struct {
	db    ports.DB
	d     Design
	st    map[string]catalog.Stmt
	order []catalog.Stmt
	world *World
	// skew is each simulated application node's clock offset (E0 only).
	skew []time.Duration
}

// maxAttempts bounds retries of engine-level contention (serialization failures,
// deadlocks, lost document writes) inside one operation.
const maxAttempts = 1000

var (
	errGaveUp = errors.New("gave up after too many retries")
	// errInvariant is a design's own safety assertion tripping.
	errInvariant       = errors.New("invariant violated")
	errAttemptDeadline = errors.New("attempt deadline passed before the operation completed")
)

func NewSeller(db ports.DB, d Design, world *World, skew time.Duration) (*Seller, error) {
	stmts, err := mustStmts(d.ID, "writes.sql")
	if err != nil {
		return nil, err
	}
	reads, err := mustStmts(d.ID, "queries.sql")
	if err != nil {
		return nil, err
	}
	s := &Seller{db: db, d: d, st: catalog.Map(append(append([]catalog.Stmt(nil), stmts...), reads...)), order: stmts,
		world: world, skew: []time.Duration{0, skew}}
	for _, name := range s.required() {
		if _, ok := s.st[name]; !ok {
			return nil, fmt.Errorf("design %s: its SQL lacks %s, required by strategy %s", d.ID, name, d.Strategy)
		}
	}
	return s, nil
}

func (s *Seller) required() []string {
	req := []string{"q01_section_map", "q02_event_sections", "q03_event_available", "q04_hold_seats",
		"q05_customer_tickets", "q06_ticket_by_id", "w_insert_tickets", "w_cancel_ticket", "w_publish_event"}
	switch s.d.Strategy {
	case Guarded:
		req = append(req, "w_hold_seats", "w_check_hold", "w_confirm_seats", "w_release_hold",
			"w_expired_holds", "w_release_expired", "w_unsell_seat")
		if s.d.PaymentWindow {
			req = append(req, "w_begin_checkout")
		}
	case LockWait, LockNowait:
		req = append(req, "w_lock_seats", "w_hold_seats", "w_check_hold", "w_confirm_seats", "w_release_hold",
			"w_expired_holds", "w_release_expired", "w_unsell_seat")
	case CheckThenAct:
		req = append(req, "w_read_seats", "w_set_held", "w_check_hold", "w_confirm_seats", "w_release_hold",
			"w_expired_holds", "w_release_expired", "w_unsell_seat")
	case Cart:
		req = append(req, "w_insert_cart", "w_hold_seats", "w_check_hold", "w_confirm_cart", "w_confirm_seats",
			"w_release_cart", "w_release_hold", "w_expired_carts", "w_release_cart_seats", "w_release_carts", "w_unsell_seat")
	case Claim:
		req = append(req, "w_claim_seats", "w_check_hold", "w_confirm_seats", "w_release_hold",
			"w_expired_holds", "w_release_expired", "w_unsell_seat")
	case Document:
		req = append(req, "w_read_section", "w_write_section", "w_expired_sections")
	}
	return req
}

// appNow is simulated application node's clock. Only E0's SQL reads it.
func (s *Seller) appNow(node int) time.Time { return time.Now().Add(s.skew[node%len(s.skew)]) }

func (s *Seller) vals(node int, b Block) map[string]any {
	v := map[string]any{"app_now": s.appNow(node)}
	if b.Event != nil {
		v["event_id"] = b.Event.ID
		v["venue_id"] = b.Event.VenueID
		v["price_cents"] = b.Event.PriceCents
	}
	if b.Section != 0 {
		v["section_no"] = b.Section
	}
	if b.Seats != nil {
		v["seat_ids"] = b.Seats
	}
	return v
}

func (s *Seller) args(name string, vals map[string]any) []any {
	a, err := catalog.Bind(s.st[name].Params, vals)
	if err != nil {
		// A missing binding is a harness bug, never a runtime condition.
		panic(fmt.Sprintf("%s.%s: %v", s.d.ID, name, err))
	}
	return a
}

func (s *Seller) exec(ctx context.Context, q ports.Queryer, name string, vals map[string]any) (int64, error) {
	return q.Exec(ctx, s.st[name].SQL, s.args(name, vals)...)
}

func (s *Seller) query(ctx context.Context, q ports.Queryer, name string, vals map[string]any, scan func(ports.Rows) error) error {
	rows, err := q.Query(ctx, s.st[name].SQL, s.args(name, vals)...)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		if err := scan(rows); err != nil {
			return err
		}
	}
	return rows.Err()
}

// seatRows scans a RETURNING list of (seat_id, db_now, db_clock[, expires_at]).
func (s *Seller) seatRows(ctx context.Context, q ports.Queryer, name string, vals map[string]any, withExpiry bool) ([]Row, error) {
	var out []Row
	err := s.query(ctx, q, name, vals, func(rs ports.Rows) error {
		var r Row
		if withExpiry {
			if err := rs.Scan(&r.Seat, &r.Now, &r.Clock, &r.Expires); err != nil {
				return err
			}
		} else if err := rs.Scan(&r.Seat, &r.Now, &r.Clock); err != nil {
			return err
		}
		out = append(out, r)
		return nil
	})
	return out, err
}

type attemptDeadlineKey struct{}

// WithAttemptDeadline marks a context so that no NEW attempt starts after t. An
// attempt already running finishes, so a deadline never manufactures a commit of
// unknown outcome (study 02).
func WithAttemptDeadline(ctx context.Context, t time.Time) context.Context {
	return context.WithValue(ctx, attemptDeadlineKey{}, t)
}

// retry runs attempt until it succeeds, retrying engine contention with jittered
// exponential backoff (2 ms base, 50 ms cap, real time). attempt returns
// again=true to start over without an error (a lost document write).
func (s *Seller) retry(ctx context.Context, retries *int, attempt func() (again bool, err error)) error {
	deadline, hasDeadline := ctx.Value(attemptDeadlineKey{}).(time.Time)
	backoff := 2 * time.Millisecond
	for i := 0; i < maxAttempts; i++ {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		if hasDeadline && i > 0 && time.Now().After(deadline) {
			return errAttemptDeadline
		}
		again, err := attempt()
		if err != nil && s.db.Classify(err) != ports.ErrRetryable {
			return err
		}
		if err == nil && !again {
			return nil
		}
		// Engine contention and a lost compare-and-set are the same thing to a buyer:
		// someone else got there first. Both back off (study 02 retried at once and
		// its analysis suspected the spin understated the designs that retry most).
		*retries++
		time.Sleep(time.Duration(rand.Int63n(int64(backoff))) + backoff/2)
		backoff = min(2*backoff, 50*time.Millisecond)
	}
	return errGaveUp
}

// commit records an ambiguous outcome against the event before returning it.
func (s *Seller) commit(ctx context.Context, tx ports.Tx, ev *Event) error {
	err := tx.Commit(ctx)
	if err != nil && s.db.Classify(err) == ports.ErrAmbiguousCommit && ev != nil {
		if el := s.world.ledger.Event(ev.ID); el != nil {
			el.Ambiguous()
		}
	}
	return err
}

// ---------------------------------------------------------------------------
// Hold
// ---------------------------------------------------------------------------

type HoldResult struct {
	HoldID   int64
	Granted  bool
	Conflict bool // a seat was not free: choose again
	Retries  int
	Expires  time.Time
}

// Hold tries to take block b for customer for ttl.
func (s *Seller) Hold(ctx context.Context, node int, b Block, customer int64, ttl time.Duration) (HoldResult, error) {
	var res HoldResult
	n := len(b.Seats)
	var rows []Row
	err := s.retry(ctx, &res.Retries, func() (bool, error) {
		res.HoldID = s.world.nextHold.Add(1)
		rows = nil
		vals := s.vals(node, b)
		vals["hold_id"], vals["customer_id"], vals["hold_ms"] = res.HoldID, customer, float64(ttl.Microseconds())/1000
		if s.d.Strategy == Document {
			return s.holdDocument(ctx, b, customer, ttl, &res, &rows)
		}
		tx, err := s.db.Begin(ctx, s.d.Isolation)
		if err != nil {
			return false, err
		}
		defer tx.Rollback(ctx)
		switch s.d.Strategy {
		case Guarded:
			rows, err = s.seatRows(ctx, tx, "w_hold_seats", vals, true)
		case LockWait, LockNowait:
			free, lerr := s.lockedFree(ctx, tx, vals, n)
			if lerr != nil {
				if s.d.Strategy == LockNowait && s.lockUnavailable(lerr) {
					res.Conflict = true
					return false, nil
				}
				return false, lerr
			}
			if !free {
				res.Conflict = true
				return false, nil
			}
			rows, err = s.seatRows(ctx, tx, "w_hold_seats", vals, true)
			if err == nil && len(rows) != n {
				return false, fmt.Errorf("%w: %d of %d locked free seats taken", errInvariant, len(rows), n)
			}
		case CheckThenAct:
			var nfree int
			if err := s.query(ctx, tx, "w_read_seats", vals, func(rs ports.Rows) error {
				var seat int32
				var free bool
				if err := rs.Scan(&seat, &free); err != nil {
					return err
				}
				if free {
					nfree++
				}
				return nil
			}); err != nil {
				return false, err
			}
			if nfree != n {
				// Committed rather than rolled back: under SERIALIZABLE a
				// read-only conclusion is only valid if the transaction commits.
				if err := tx.Commit(ctx); err != nil {
					return false, err
				}
				res.Conflict = true
				return false, nil
			}
			rows, err = s.seatRows(ctx, tx, "w_set_held", vals, true)
		case Cart:
			var now, clock, exp time.Time
			if err := tx.QueryRow(ctx, s.st["w_insert_cart"].SQL, s.args("w_insert_cart", vals)...).Scan(&now, &clock, &exp); err != nil {
				return false, err
			}
			rows, err = s.seatRows(ctx, tx, "w_hold_seats", vals, false)
			for i := range rows {
				rows[i].Expires = exp
			}
		case Claim:
			rows, err = s.seatRows(ctx, tx, "w_claim_seats", vals, true)
			if s.db.Classify(err) == ports.ErrUniqueViolation {
				res.Conflict = true
				return false, nil
			}
		default:
			return false, fmt.Errorf("strategy %q not implemented", s.d.Strategy)
		}
		if err != nil {
			return false, err
		}
		if len(rows) != n {
			res.Conflict = true
			return false, nil
		}
		if err := s.commit(ctx, tx, b.Event); err != nil {
			return false, err
		}
		return false, nil
	})
	if err != nil || res.Conflict {
		return res, err
	}
	res.Granted = true
	res.Expires = rows[0].Expires
	s.world.ledger.Event(b.Event.ID).Grant(HoldInfo{ID: res.HoldID, EventID: b.Event.ID, Section: b.Section,
		Seats: b.Seats, Customer: customer}, rows)
	return res, nil
}

func (s *Seller) lockedFree(ctx context.Context, tx ports.Tx, vals map[string]any, n int) (bool, error) {
	count, free := 0, true
	err := s.query(ctx, tx, "w_lock_seats", vals, func(rs ports.Rows) error {
		var seat int32
		var f bool
		if err := rs.Scan(&seat, &f); err != nil {
			return err
		}
		count++
		free = free && f
		return nil
	})
	return err == nil && free && count == n, err
}

// lockUnavailable: NOWAIT met a locked seat. Engine probe P2 decides whether an
// engine reports it as 55P03 or as a serialization-class conflict; from
// w_lock_seats either means "being taken -- choose again".
func (s *Seller) lockUnavailable(err error) bool {
	c := s.db.Classify(err)
	return c == ports.ErrLockNotAvailable || c == ports.ErrRetryable
}

// ---------------------------------------------------------------------------
// L2: the section document
// ---------------------------------------------------------------------------

type sectionDoc struct {
	version int64
	claims  Claims
	now     time.Time
	clock   time.Time
}

func (s *Seller) readSection(ctx context.Context, q ports.Queryer, ev *Event, section int32) (sectionDoc, error) {
	var doc sectionDoc
	var text string
	err := q.QueryRow(ctx, s.st["w_read_section"].SQL, s.args("w_read_section",
		map[string]any{"event_id": ev.ID, "section_no": section})...).Scan(&doc.version, &text, &doc.now, &doc.clock)
	if err != nil {
		return doc, fmt.Errorf("w_read_section %d/%d: %w", ev.ID, section, err)
	}
	doc.claims, err = DecodeClaims(text)
	return doc, err
}

// writeSection is the compare-and-set. ok=false: someone wrote the section first.
func (s *Seller) writeSection(ctx context.Context, q ports.Queryer, ev *Event, section int32, doc sectionDoc) (ok bool, now, clock time.Time, err error) {
	text, next := doc.claims.Encode()
	err = q.QueryRow(ctx, s.st["w_write_section"].SQL, s.args("w_write_section", map[string]any{
		"claims": text, "next_expiry": next, "event_id": ev.ID, "section_no": section, "version": doc.version,
	})...).Scan(&now, &clock)
	if errors.Is(err, ports.ErrNoRows) {
		return false, now, clock, nil
	}
	return err == nil, now, clock, err
}

func (s *Seller) holdDocument(ctx context.Context, b Block, customer int64, ttl time.Duration, res *HoldResult, rows *[]Row) (bool, error) {
	doc, err := s.readSection(ctx, s.db, b.Event, b.Section)
	if err != nil {
		return false, err
	}
	nowMS := doc.now.UnixMilli()
	for _, seat := range b.Seats {
		if c, ok := doc.claims[seat]; ok && (c.State == 2 || c.ExpiresMS > nowMS) {
			res.Conflict = true
			return false, nil
		}
	}
	exp := doc.now.Add(ttl).Truncate(time.Millisecond)
	for _, seat := range b.Seats {
		doc.claims[seat] = ClaimVal{State: 1, Hold: res.HoldID, Customer: customer, ExpiresMS: exp.UnixMilli()}
	}
	ok, _, clock, err := s.writeSection(ctx, s.db, b.Event, b.Section, doc)
	if err != nil {
		return false, err
	}
	if !ok {
		return true, nil // lost the compare-and-set: read the section again
	}
	*rows = (*rows)[:0]
	for _, seat := range b.Seats {
		// The decision was taken on the read's clock; the write's clock orders it.
		*rows = append(*rows, Row{Seat: seat, Now: doc.now, Clock: clock, Expires: exp})
	}
	return false, nil
}

// ---------------------------------------------------------------------------
// Pre-payment check, checkout, confirmation
// ---------------------------------------------------------------------------

// CheckHold is the friendly pre-payment check.
func (s *Seller) CheckHold(ctx context.Context, node int, b Block, holdID int64) (bool, error) {
	if s.d.Strategy == Document {
		doc, err := s.readSection(ctx, s.db, b.Event, b.Section)
		if err != nil {
			return false, err
		}
		return docHoldValid(doc, b.Seats, holdID), nil
	}
	vals := s.vals(node, b)
	vals["hold_id"] = holdID
	var valid bool
	err := s.db.QueryRow(ctx, s.st["w_check_hold"].SQL, s.args("w_check_hold", vals)...).Scan(&valid)
	return valid, err
}

func docHoldValid(doc sectionDoc, seats []int32, holdID int64) bool {
	nowMS := doc.now.UnixMilli()
	for _, seat := range seats {
		c, ok := doc.claims[seat]
		if !ok || c.State != 1 || c.Hold != holdID || c.ExpiresMS <= nowMS {
			return false
		}
	}
	return true
}

// txnNow is the harness's own instrumentation: the confirming transaction's
// now(), so a refused confirmation has a timestamp too. Same in every design.
func txnNow(ctx context.Context, tx ports.Tx) (time.Time, error) {
	var t time.Time
	err := tx.QueryRow(ctx, "SELECT now()").Scan(&t)
	return t, err
}

type ConfirmResult struct {
	Sold    bool
	Class   string // rejection class when not sold
	Retries int
	// ShortRetried: S1r ran the confirmation a second time after a short match.
	// RetrySold: and the second attempt sold.
	ShortRetried, RetrySold bool
}

// earlyAt says whether the ledger would call a refusal of the hold at txn now early
// (at least G before its expiry). Only then does the harness spend statements on
// diagnostics, so late and boundary refusals take exactly the design's path.
func (s *Seller) earlyAt(b Block, holdID int64, tnow time.Time) bool {
	h, ok := s.world.ledger.Event(b.Event.ID).Hold(holdID)
	return ok && h.Expires.Sub(tnow) >= s.world.ledger.guard
}

// refusal is what the in-transaction diagnostics of an early refusal found (AM-01.2).
type refusal struct {
	ran       bool
	matched   int    // rows the refusing statement matched; -1: a one-row statement returned none
	inTx      string // the re-read and the re-issue, inside the refusing transaction
	transient bool   // the re-issue matched every seat the first execution did not
}

// diagnoseInTx runs inside the refusing transaction, before its rollback: the hold's
// rows re-read, then the identical statement with identical bindings. The refusal
// stands whatever the re-issue matches -- the transaction is rolled back as always --
// it only tells a transient miss from a persistent one. Harness instrumentation, the
// same for every design with a guarded statement (ER-01, AM-01.2).
func (s *Seller) diagnoseInTx(ctx context.Context, tx ports.Tx, b Block, holdID int64, name string, vals map[string]any, matched int, withExpiry bool) refusal {
	r := refusal{ran: true, matched: matched, inTx: s.seatStates(ctx, tx, b, holdID)}
	if matched < 0 {
		var now, clock time.Time
		err := tx.QueryRow(ctx, s.st[name].SQL, s.args(name, vals)...).Scan(&now, &clock)
		switch {
		case err == nil:
			r.transient = true
			r.inTx += "; the same statement issued again in this transaction returned its row"
		case errors.Is(err, ports.ErrNoRows):
			r.inTx += "; the same statement issued again in this transaction returned no row"
		default:
			r.inTx += "; the same statement issued again in this transaction failed: " + err.Error()
		}
		return r
	}
	again, err := s.seatRows(ctx, tx, name, vals, withExpiry)
	if err != nil {
		r.inTx += "; the same statement issued again in this transaction failed: " + err.Error()
		return r
	}
	n := len(b.Seats)
	r.transient = matched+len(again) == n
	r.inTx += fmt.Sprintf("; the same statement issued again in this transaction matched %d more (%d of %d in all)",
		len(again), matched+len(again), n)
	return r
}

// BeginCheckout is K1's payment window. Rejected = the hold expired before paying.
func (s *Seller) BeginCheckout(ctx context.Context, node int, b Block, holdID int64, window time.Duration) (ConfirmResult, error) {
	var res ConfirmResult
	var rows []Row
	var tnow time.Time
	var ref refusal
	err := s.retry(ctx, &res.Retries, func() (bool, error) {
		rows, ref = nil, refusal{}
		vals := s.vals(node, b)
		vals["hold_id"], vals["payment_window_ms"] = holdID, float64(window.Microseconds())/1000
		tx, err := s.db.Begin(ctx, s.d.Isolation)
		if err != nil {
			return false, err
		}
		defer tx.Rollback(ctx)
		if tnow, err = txnNow(ctx, tx); err != nil {
			return false, err
		}
		if rows, err = s.seatRows(ctx, tx, "w_begin_checkout", vals, true); err != nil {
			return false, err
		}
		if len(rows) != len(b.Seats) {
			if s.earlyAt(b, holdID, tnow) {
				ref = s.diagnoseInTx(ctx, tx, b, holdID, "w_begin_checkout", vals, len(rows), true)
			}
			rows = nil
			return false, nil
		}
		return false, s.commit(ctx, tx, b.Event)
	})
	if err != nil {
		return res, err
	}
	el := s.world.ledger.Event(b.Event.ID)
	if rows == nil {
		res.Class = el.Reject(holdID, tnow, s.diagnoseRefusal(ctx, node, b, holdID, tnow, ref), ref.transient)
		return res, nil
	}
	el.Extend(holdID, rows)
	res.Sold = true
	return res, nil
}

// Confirm turns a valid hold into tickets.
func (s *Seller) Confirm(ctx context.Context, node int, b Block, holdID, customer int64) (ConfirmResult, error) {
	var res ConfirmResult
	n := len(b.Seats)
	ticketIDs := make([]int64, n)
	ticketOf := map[int32]int64{}
	for i, seat := range b.Seats {
		ticketIDs[i] = s.world.nextTicket.Add(1)
		ticketOf[seat] = ticketIDs[i]
	}
	var rows []Row
	var tnow time.Time
	var rejected bool
	var ref refusal
	for attempt := 1; ; attempt++ {
		// S1r (AM-01.3): a first confirmation that matches too few seats is rolled
		// back and run once more, at once, in a new transaction. It runs no
		// diagnostics; the second attempt is exactly S1's.
		retryShort := s.d.RetryShortConfirm && attempt == 1
		err := s.retry(ctx, &res.Retries, func() (bool, error) {
			rows, rejected, ref = nil, false, refusal{}
			vals := s.vals(node, b)
			vals["hold_id"], vals["customer_id"], vals["ticket_ids"] = holdID, customer, ticketIDs
			if s.d.Strategy == Document {
				return s.confirmDocument(ctx, b, holdID, customer, vals, &rows, &tnow, &rejected)
			}
			tx, err := s.db.Begin(ctx, s.d.Isolation)
			if err != nil {
				return false, err
			}
			defer tx.Rollback(ctx)
			if tnow, err = txnNow(ctx, tx); err != nil {
				return false, err
			}
			if s.d.Strategy == Cart {
				var now, clock time.Time
				err := tx.QueryRow(ctx, s.st["w_confirm_cart"].SQL, s.args("w_confirm_cart", vals)...).Scan(&now, &clock)
				if errors.Is(err, ports.ErrNoRows) {
					rejected = true
					if s.earlyAt(b, holdID, tnow) {
						ref = s.diagnoseInTx(ctx, tx, b, holdID, "w_confirm_cart", vals, -1, false)
					}
					return false, nil
				}
				if err != nil {
					return false, err
				}
			}
			if rows, err = s.seatRows(ctx, tx, "w_confirm_seats", vals, false); err != nil {
				return false, err
			}
			if len(rows) != n {
				rejected = true
				if !retryShort && s.earlyAt(b, holdID, tnow) {
					ref = s.diagnoseInTx(ctx, tx, b, holdID, "w_confirm_seats", vals, len(rows), false)
				}
				return false, nil
			}
			if _, err := s.exec(ctx, tx, "w_insert_tickets", vals); err != nil {
				if s.db.Classify(err) == ports.ErrUniqueViolation {
					return false, fmt.Errorf("%w: a confirmed seat already had a ticket: %v", errInvariant, err)
				}
				return false, err
			}
			return false, s.commit(ctx, tx, b.Event)
		})
		if err != nil {
			return res, err
		}
		if rejected && retryShort {
			res.ShortRetried = true
			continue
		}
		break
	}
	el := s.world.ledger.Event(b.Event.ID)
	if rejected {
		res.Class = el.Reject(holdID, tnow, s.diagnoseRefusal(ctx, node, b, holdID, tnow, ref), ref.transient)
		return res, nil
	}
	el.Sale(holdID, rows, ticketOf)
	res.Sold = true
	res.RetrySold = res.ShortRetried
	return res, nil
}

func (s *Seller) confirmDocument(ctx context.Context, b Block, holdID, customer int64, vals map[string]any, rows *[]Row, tnow *time.Time, rejected *bool) (bool, error) {
	doc, err := s.readSection(ctx, s.db, b.Event, b.Section)
	if err != nil {
		return false, err
	}
	*tnow = doc.now
	if !docHoldValid(doc, b.Seats, holdID) {
		*rejected = true
		return false, nil
	}
	for _, seat := range b.Seats {
		doc.claims[seat] = ClaimVal{State: 2, Hold: holdID, Customer: customer}
	}
	tx, err := s.db.Begin(ctx, s.d.Isolation)
	if err != nil {
		return false, err
	}
	defer tx.Rollback(ctx)
	ok, _, clock, err := s.writeSection(ctx, tx, b.Event, b.Section, doc)
	if err != nil {
		return false, err
	}
	if !ok {
		return true, nil // lost the compare-and-set: read again
	}
	if _, err := s.exec(ctx, tx, "w_insert_tickets", vals); err != nil {
		return false, err
	}
	if err := s.commit(ctx, tx, b.Event); err != nil {
		return false, err
	}
	*rows = (*rows)[:0]
	for _, seat := range b.Seats {
		*rows = append(*rows, Row{Seat: seat, Now: doc.now, Clock: clock})
	}
	return false, nil
}

// ---------------------------------------------------------------------------
// Release and sweep
// ---------------------------------------------------------------------------

// Release gives back whatever seats of the hold it still has.
func (s *Seller) Release(ctx context.Context, node int, b Block, holdID int64) (int, error) {
	var rows []Row
	var retries int
	err := s.retry(ctx, &retries, func() (bool, error) {
		rows = nil
		if s.d.Strategy == Document {
			doc, err := s.readSection(ctx, s.db, b.Event, b.Section)
			if err != nil {
				return false, err
			}
			var mine []int32
			for _, seat := range b.Seats {
				if c, ok := doc.claims[seat]; ok && c.State == 1 && c.Hold == holdID {
					delete(doc.claims, seat)
					mine = append(mine, seat)
				}
			}
			if len(mine) == 0 {
				return false, nil
			}
			ok, _, clock, err := s.writeSection(ctx, s.db, b.Event, b.Section, doc)
			if err != nil || !ok {
				return !ok && err == nil, err
			}
			for _, seat := range mine {
				rows = append(rows, Row{Seat: seat, Now: doc.now, Clock: clock})
			}
			return false, nil
		}
		vals := s.vals(node, b)
		vals["hold_id"] = holdID
		tx, err := s.db.Begin(ctx, ports.ReadCommitted)
		if err != nil {
			return false, err
		}
		defer tx.Rollback(ctx)
		if s.d.Strategy == Cart {
			if _, err := s.exec(ctx, tx, "w_release_cart", vals); err != nil {
				return false, err
			}
		}
		if rows, err = s.seatRows(ctx, tx, "w_release_hold", vals, false); err != nil {
			return false, err
		}
		return false, s.commit(ctx, tx, b.Event)
	})
	if err != nil {
		return 0, err
	}
	if len(rows) > 0 {
		s.world.ledger.Event(b.Event.ID).Release(holdID, rows)
	}
	return len(rows), nil
}

// SweepResult is one sweeper batch.
type SweepResult struct {
	Released int
	// Lags are release time minus expiry, per released seat.
	Lags []time.Duration
}

type released struct {
	event int64
	hold  int64
	row   Row
	exp   time.Time
}

// Sweep releases up to batch expired holds.
func (s *Seller) Sweep(ctx context.Context, node, batch int) (SweepResult, error) {
	var out []released
	var retries int
	err := s.retry(ctx, &retries, func() (bool, error) {
		out = nil
		tx, err := s.db.Begin(ctx, ports.ReadCommitted)
		if err != nil {
			return false, err
		}
		defer tx.Rollback(ctx)
		switch s.d.Strategy {
		case Cart:
			out, err = s.sweepCarts(ctx, tx, node, batch)
		case Document:
			out, err = s.sweepDocuments(ctx, tx, batch)
		default:
			out, err = s.sweepRows(ctx, tx, node, batch)
		}
		if err != nil {
			return false, err
		}
		return false, tx.Commit(ctx)
	})
	var res SweepResult
	if err != nil {
		return res, err
	}
	type key struct{ event, hold int64 }
	byHold := map[key][]Row{}
	for _, r := range out {
		byHold[key{r.event, r.hold}] = append(byHold[key{r.event, r.hold}], r.row)
		// Release lag measures the sweeper against holds granted in this phase. Holds
		// loaded already expired (the dataset's 1%) expired before the phase began and
		// would report the load's age instead (tens of thousands of human minutes).
		if el := s.world.ledger.Event(r.event); el != nil {
			if h, ok := el.Hold(r.hold); ok && h.Loaded {
				continue
			}
		}
		res.Lags = append(res.Lags, r.row.Now.Sub(r.exp))
	}
	for k, rows := range byHold {
		if el := s.world.ledger.Event(k.event); el != nil {
			el.Release(k.hold, rows)
		}
	}
	res.Released = len(out)
	return res, nil
}

func (s *Seller) sweepRows(ctx context.Context, tx ports.Tx, node, batch int) ([]released, error) {
	type exp struct {
		event   int64
		section int32
		seat    int32
		hold    int64
		expires time.Time
	}
	var cand []exp
	vals := map[string]any{"batch": batch, "app_now": s.appNow(node)}
	if err := s.query(ctx, tx, "w_expired_holds", vals, func(rs ports.Rows) error {
		var e exp
		var hold *int64
		if err := rs.Scan(&e.event, &e.section, &e.seat, &hold, &e.expires); err != nil {
			return err
		}
		if hold != nil {
			e.hold = *hold
		}
		cand = append(cand, e)
		return nil
	}); err != nil {
		return nil, err
	}
	type grp struct {
		event   int64
		section int32
	}
	groups := map[grp][]exp{}
	for _, c := range cand {
		groups[grp{c.event, c.section}] = append(groups[grp{c.event, c.section}], c)
	}
	var out []released
	for g, cs := range groups {
		seats := make([]int32, 0, len(cs))
		info := map[int32]exp{}
		for _, c := range cs {
			seats = append(seats, c.seat)
			info[c.seat] = c
		}
		sort.Slice(seats, func(i, j int) bool { return seats[i] < seats[j] })
		v := map[string]any{"event_id": g.event, "section_no": g.section, "seat_ids": seats, "app_now": s.appNow(node)}
		rows, err := s.seatRows(ctx, tx, "w_release_expired", v, false)
		if err != nil {
			return nil, err
		}
		for _, r := range rows {
			c := info[r.Seat]
			out = append(out, released{event: g.event, hold: c.hold, row: r, exp: c.expires})
		}
	}
	return out, nil
}

func (s *Seller) sweepCarts(ctx context.Context, tx ports.Tx, node, batch int) ([]released, error) {
	type cart struct {
		hold, event int64
		expires     time.Time
	}
	var carts []cart
	if err := s.query(ctx, tx, "w_expired_carts", map[string]any{"batch": batch}, func(rs ports.Rows) error {
		var c cart
		if err := rs.Scan(&c.hold, &c.event, &c.expires); err != nil {
			return err
		}
		carts = append(carts, c)
		return nil
	}); err != nil {
		return nil, err
	}
	if len(carts) == 0 {
		return nil, nil
	}
	ids := make([]int64, len(carts))
	for i, c := range carts {
		ids[i] = c.hold
	}
	// Which cart a released seat belonged to: the ledger knows each cart's seats.
	owner := map[int64]map[int32]cart{}
	for _, c := range carts {
		if el := s.world.ledger.Event(c.event); el != nil {
			if h, ok := el.Hold(c.hold); ok {
				if owner[c.event] == nil {
					owner[c.event] = map[int32]cart{}
				}
				for _, seat := range h.Seats {
					owner[c.event][seat] = c
				}
			}
		}
	}
	var out []released
	if err := s.query(ctx, tx, "w_release_cart_seats", map[string]any{"hold_ids": ids}, func(rs ports.Rows) error {
		var ev int64
		var r Row
		if err := rs.Scan(&ev, &r.Seat, &r.Now, &r.Clock); err != nil {
			return err
		}
		c := owner[ev][r.Seat]
		out = append(out, released{event: ev, hold: c.hold, row: r, exp: c.expires})
		return nil
	}); err != nil {
		return nil, err
	}
	if _, err := s.exec(ctx, tx, "w_release_carts", map[string]any{"hold_ids": ids}); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *Seller) sweepDocuments(ctx context.Context, tx ports.Tx, batch int) ([]released, error) {
	type sec struct {
		event   int64
		section int32
	}
	var secs []sec
	if err := s.query(ctx, tx, "w_expired_sections", map[string]any{"batch": batch}, func(rs ports.Rows) error {
		var x sec
		if err := rs.Scan(&x.event, &x.section); err != nil {
			return err
		}
		secs = append(secs, x)
		return nil
	}); err != nil {
		return nil, err
	}
	var out []released
	for _, x := range secs {
		ev := s.world.event(x.event)
		doc, err := s.readSection(ctx, tx, ev, x.section)
		if err != nil {
			return nil, err
		}
		nowMS := doc.now.UnixMilli()
		var gone []released
		for seat, c := range doc.claims {
			if c.State == 1 && c.ExpiresMS <= nowMS {
				delete(doc.claims, seat)
				gone = append(gone, released{event: x.event, hold: c.Hold, row: Row{Seat: seat, Now: doc.now},
					exp: time.UnixMilli(c.ExpiresMS)})
			}
		}
		ok, _, clock, err := s.writeSection(ctx, tx, ev, x.section, doc)
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, fmt.Errorf("%w: a section locked by the sweeper changed version", errInvariant)
		}
		for i := range gone {
			gone[i].row.Clock = clock
		}
		out = append(out, gone...)
	}
	return out, nil
}

// ---------------------------------------------------------------------------
// Refund and publishing
// ---------------------------------------------------------------------------

// Cancel refunds a ticket. ok=false: it was already refunded.
func (s *Seller) Cancel(ctx context.Context, ticketID int64) (bool, error) {
	var (
		ok           bool
		ev           int64
		section      int32
		seat         int32
		now, clock   time.Time
		retries      int
		unsoldAnyway bool
	)
	err := s.retry(ctx, &retries, func() (bool, error) {
		ok, unsoldAnyway = false, false
		tx, err := s.db.Begin(ctx, ports.ReadCommitted)
		if err != nil {
			return false, err
		}
		defer tx.Rollback(ctx)
		err = tx.QueryRow(ctx, s.st["w_cancel_ticket"].SQL, s.args("w_cancel_ticket", map[string]any{"ticket_id": ticketID})...).
			Scan(&ev, &section, &seat, &now, &clock)
		if errors.Is(err, ports.ErrNoRows) {
			return false, nil
		}
		if err != nil {
			return false, err
		}
		e := s.world.event(ev)
		if s.d.Strategy == Document {
			doc, err := s.readSection(ctx, tx, e, section)
			if err != nil {
				return false, err
			}
			if c, has := doc.claims[seat]; has && c.State == 2 {
				delete(doc.claims, seat)
			} else {
				unsoldAnyway = true
			}
			wok, _, _, err := s.writeSection(ctx, tx, e, section, doc)
			if err != nil {
				return false, err
			}
			if !wok {
				return true, nil
			}
		} else {
			n, err := s.exec(ctx, tx, "w_unsell_seat", map[string]any{"event_id": ev, "section_no": section, "seat_id": seat})
			if err != nil {
				return false, err
			}
			unsoldAnyway = n != 1
		}
		if err := s.commit(ctx, tx, e); err != nil {
			return false, err
		}
		ok = true
		return false, nil
	})
	if err != nil || !ok {
		return false, err
	}
	if unsoldAnyway {
		return true, fmt.Errorf("%w: refunded ticket %d's seat %d/%d was not sold in the inventory", errInvariant, ticketID, ev, seat)
	}
	s.world.ledger.Event(ev).Refund(ticketID, seat, now, clock)
	return true, nil
}

// Publish creates an event at a venue, with its inventory, in one transaction:
// every w_publish_* statement of the design in file order. That cost is the
// pre-creation question in its purest form.
func (s *Seller) Publish(ctx context.Context, bandID int64, v *Venue, r *rand.Rand) (*Event, error) {
	e := &Event{ID: s.world.nextEvent.Add(1), BandID: bandID, VenueID: v.ID, Tier: v.Tier,
		Name: "published show", StartsAt: showEpoch.Add(time.Duration(r.Intn(365*24)) * time.Hour),
		PriceCents: int64(2000 + 500*r.Intn(30)), Description: "Just announced.", Kind: "published"}
	vals := map[string]any{"event_id": e.ID, "band_id": bandID, "venue_id": v.ID, "name": e.Name,
		"starts_at": e.StartsAt, "price_cents": e.PriceCents, "description": e.Description}
	var retries int
	err := s.retry(ctx, &retries, func() (bool, error) {
		tx, err := s.db.Begin(ctx, ports.ReadCommitted)
		if err != nil {
			return false, err
		}
		defer tx.Rollback(ctx)
		for _, st := range s.order {
			if !strings.HasPrefix(st.Name, "w_publish_") {
				continue
			}
			if _, err := s.exec(ctx, tx, st.Name, vals); err != nil {
				return false, fmt.Errorf("%s: %w", st.Name, err)
			}
		}
		return false, tx.Commit(ctx)
	})
	if err != nil {
		return nil, err
	}
	s.world.mu.Lock()
	s.world.published[e.ID] = e
	s.world.mu.Unlock()
	s.world.ledger.AddEvent(e.ID, v)
	return e, nil
}

// ---------------------------------------------------------------------------
// Reads buyers make
// ---------------------------------------------------------------------------

// SectionMap returns the non-available seats of a section (1 held, 2 sold).
func (s *Seller) SectionMap(ctx context.Context, node int, ev *Event, section int32) (map[int32]uint8, error) {
	out := map[int32]uint8{}
	err := s.query(ctx, s.db, "q01_section_map", s.vals(node, Block{Event: ev, Section: section}), func(rs ports.Rows) error {
		var seat int32
		var state int16
		if err := rs.Scan(&seat, &state); err != nil {
			return err
		}
		out[seat] = uint8(state)
		return nil
	})
	return out, err
}

type SectionAvail struct {
	No        int32
	Available int64
}

// Sections returns seats available per section.
func (s *Seller) Sections(ctx context.Context, node int, ev *Event) ([]SectionAvail, error) {
	var out []SectionAvail
	err := s.query(ctx, s.db, "q02_event_sections", s.vals(node, Block{Event: ev}), func(rs ports.Rows) error {
		var a SectionAvail
		if err := rs.Scan(&a.No, &a.Available); err != nil {
			return err
		}
		out = append(out, a)
		return nil
	})
	return out, err
}

// Available returns the event's available seats.
func (s *Seller) Available(ctx context.Context, node int, ev *Event) (int64, error) {
	var n int64
	err := s.db.QueryRow(ctx, s.st["q03_event_available"].SQL, s.args("q03_event_available", s.vals(node, Block{Event: ev}))...).Scan(&n)
	return n, err
}

// diagnoseRefusal is instrumentation for a refusal the ledger will call early (the
// hold had at least G left), run after the rollback: what the refusing transaction
// found (ref), and how many seats of the hold the database reports as validly held
// right afterwards. Together they tell an engine that silently matched too few rows
// from a hold that something the harness did not record had changed. Empty when not
// early.
func (s *Seller) diagnoseRefusal(ctx context.Context, node int, b Block, holdID int64, tnow time.Time, ref refusal) string {
	h, ok := s.world.ledger.Event(b.Event.ID).Hold(holdID)
	if !ok || h.Expires.Sub(tnow) < s.world.ledger.guard {
		return ""
	}
	first, inTx := "n/a", "no diagnostics (not a guarded statement)"
	if ref.ran {
		inTx = ref.inTx
		first = fmt.Sprintf("%d of %d seats", ref.matched, len(b.Seats))
		if ref.matched < 0 {
			first = "no row"
		}
	}
	var still int
	vals := s.vals(node, b)
	vals["hold_id"] = holdID
	err := s.query(ctx, s.db, "q04_hold_seats", vals, func(rs ports.Rows) error {
		var seat int32
		var exp time.Time
		if err := rs.Scan(&seat, &exp); err != nil {
			return err
		}
		still++
		return nil
	})
	if err != nil {
		return fmt.Sprintf("; diagnosis failed: %v", err)
	}
	after := s.seatStates(ctx, s.db, b, holdID)
	return fmt.Sprintf("; the refusing statement matched %s (seat ids %v); right after, the database showed %d of the hold's seats validly held by it; inside the refusing transaction: %s; re-read after: %s; grant now() %s, confirm now() %s, client %s",
		first, b.Seats, still, inTx, after, h.GrantNow.Format(time.RFC3339Nano), tnow.Format(time.RFC3339Nano), time.Now().UTC().Format(time.RFC3339Nano))
}

// seatStates is diagnostic SQL (harness instrumentation, not design SQL): the exact
// seats a statement was given, as the database shows them to a new statement on q --
// state, hold, whether that hold is valid now -- and the backend that answered.
func (s *Seller) seatStates(ctx context.Context, q ports.Queryer, b Block, holdID int64) string {
	var sql string
	switch s.d.Layout {
	case SeatRows:
		sql = `SELECT seat_id, status, COALESCE(hold_id, 0), COALESCE(hold_expires_at > now(), false), now(), pg_backend_pid()
			FROM event_seat WHERE event_id = $1 AND seat_id = ANY ($2::INT[]) ORDER BY seat_id`
	case CartRows:
		sql = `SELECT s.seat_id, s.status, COALESCE(s.hold_id, 0),
			COALESCE((SELECT h.status = 'held' AND h.expires_at > now() FROM hold h WHERE h.hold_id = s.hold_id), false),
			now(), pg_backend_pid()
			FROM event_seat s WHERE s.event_id = $1 AND s.seat_id = ANY ($2::INT[]) ORDER BY s.seat_id`
	case ClaimRows:
		sql = `SELECT seat_id, CASE WHEN expires_at = 'infinity' THEN 'sold' ELSE 'claimed' END, COALESCE(hold_id, 0),
			expires_at > now() AND expires_at <> 'infinity', now(), pg_backend_pid()
			FROM seat_claim WHERE event_id = $1 AND seat_id = ANY ($2::INT[]) ORDER BY seat_id`
	default:
		return "n/a (no seat rows)"
	}
	rows, err := q.Query(ctx, sql, b.Event.ID, b.Seats)
	if err != nil {
		return "error: " + err.Error()
	}
	defer rows.Close()
	var parts []string
	var at time.Time
	var pid int32
	for rows.Next() {
		var seat int32
		var status string
		var hold int64
		var valid bool
		if err := rows.Scan(&seat, &status, &hold, &valid, &at, &pid); err != nil {
			return "error: " + err.Error()
		}
		mine := "other hold"
		if hold == holdID {
			mine = "this hold"
		}
		parts = append(parts, fmt.Sprintf("%d %s %s valid=%v", seat, status, mine, valid))
	}
	return fmt.Sprintf("[%s] (%d rows for %d seats) at now() %s, backend pid %d", strings.Join(parts, "; "), len(parts), len(b.Seats), at.Format(time.RFC3339Nano), pid)
}
