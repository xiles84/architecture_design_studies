package main

import (
	"fmt"
	"sort"
	"sync"
	"time"
)

// ---------------------------------------------------------------------------
// The ledger: what the harness knows independently of the database.
//
// Every grant, extension, release, sale and refund the harness sees commit is
// recorded per seat, with the database clocks the statement itself returned:
//
//	now   -- now(), the transaction start: the clock a design's expiry logic uses
//	clock -- clock_timestamp(), when the row was written: the order of a release
//	         and the grant that waited for its lock
//
// Replaying a seat's events in clock order decides, exactly and after the fact:
//
//	theft (INV-2)        a grant while another hold was valid (grant now < its
//	                     expiry, and no release before) or while the seat was sold
//	sale without owner   a sale by a hold that no longer had the seat (INV-1 if the
//	                     seat was already sold: a double sale)
//	late sale (INV-6)    a sale whose transaction started at or after its expiry
//
// Rejected confirmations are classified when they happen (INV-3), by the margin
// between the hold's expiry and the confirming transaction's now().
//
// The ledger never asks the database what happened; the audit compares the two.
// ---------------------------------------------------------------------------

type evKind uint8

const (
	evGrant evKind = iota + 1
	evExtend
	evRelease
	evSale
	evRefund
)

func (k evKind) String() string {
	return [...]string{"", "grant", "extend", "release", "sale", "refund"}[k]
}

type seatEvent struct {
	kind    evKind
	hold    int64
	now     time.Time
	clock   time.Time
	expires time.Time // grant, extend
	ticket  int64     // sale, refund
}

// HoldInfo is one hold as the harness saw it granted.
type HoldInfo struct {
	ID       int64
	EventID  int64
	Section  int32
	Seats    []int32
	Customer int64
	GrantNow time.Time
	Expires  time.Time // current: K1 may extend it
	Loaded   bool
}

// Row is one row of a RETURNING list: the seat and the database clocks.
type Row struct {
	Seat    int32
	Now     time.Time
	Clock   time.Time
	Expires time.Time
}

// Rejection classes of a confirmation (or K1 checkout start) that did not succeed.
const (
	RejLate     = "late"     // m <= 0: expected under strict expiry
	RejBoundary = "boundary" // 0 < m < G: a reported cost, not a violation
	RejEarly    = "early"    // m >= G: the hold was not honored (INV-3 violation)
)

type Example struct {
	Kind    string `json:"kind"`
	EventID int64  `json:"event_id"`
	SeatID  int32  `json:"seat_id,omitempty"`
	HoldID  int64  `json:"hold_id,omitempty"`
	Other   int64  `json:"other_hold_id,omitempty"`
	Detail  string `json:"detail,omitempty"`
}

const maxExamples = 5

// EventLedger is the record of one event. Safe for concurrent use.
type EventLedger struct {
	mu      sync.Mutex
	ID      int64
	venue   *Venue
	guard   time.Duration
	seats   map[int32][]seatEvent
	holds   map[int64]*HoldInfo
	tickets map[int64]int32 // ticket id -> seat, tickets the harness saw created
	// ambiguous counts commits of unknown outcome touching this event.
	ambiguous int64

	rejLate, rejBoundary, rejEarly int64
	examples                       []Example
}

type Ledger struct {
	mu     sync.RWMutex
	events map[int64]*EventLedger
	guard  time.Duration
}

// NewLedger reflects a freshly loaded dataset. It is rebuilt on every reload, so a
// later phase never inherits an earlier phase's ledger.
func NewLedger(ds *Dataset, loadNow time.Time, guard time.Duration) *Ledger {
	l := &Ledger{events: map[int64]*EventLedger{}, guard: guard}
	for i := range ds.Events {
		e := &ds.Events[i]
		l.events[e.ID] = &EventLedger{ID: e.ID, venue: ds.venue(e.VenueID), guard: guard,
			seats: map[int32][]seatEvent{}, holds: map[int64]*HoldInfo{}, tickets: map[int64]int32{}}
	}
	old := loadNow.Add(-24 * 365 * time.Hour)
	for _, t := range ds.Sold {
		el := l.events[t.EventID]
		el.seats[t.SeatID] = append(el.seats[t.SeatID], seatEvent{kind: evSale, now: old, clock: old, ticket: t.ID})
		el.tickets[t.ID] = t.SeatID
	}
	for _, h := range ds.Holds {
		el := l.events[h.EventID]
		created, exp := h.createdAt(loadNow), h.expiresAt(loadNow)
		el.holds[h.ID] = &HoldInfo{ID: h.ID, EventID: h.EventID, Section: h.Section, Seats: h.Seats,
			Customer: h.CustomerID, GrantNow: created, Expires: exp, Loaded: true}
		for _, s := range h.Seats {
			el.seats[s] = append(el.seats[s], seatEvent{kind: evGrant, hold: h.ID, now: created, clock: created, expires: exp})
		}
	}
	return l
}

// setGuard changes the guard margin G, for a load reused by a phase on another clock.
func (l *Ledger) setGuard(g time.Duration) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.guard = g
	for _, el := range l.events {
		el.mu.Lock()
		el.guard = g
		el.mu.Unlock()
	}
}

func (l *Ledger) Event(id int64) *EventLedger {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.events[id]
}

// AddEvent registers an event published during a phase.
func (l *Ledger) AddEvent(id int64, v *Venue) *EventLedger {
	l.mu.Lock()
	defer l.mu.Unlock()
	el := &EventLedger{ID: id, venue: v, guard: l.guard, seats: map[int32][]seatEvent{},
		holds: map[int64]*HoldInfo{}, tickets: map[int64]int32{}}
	l.events[id] = el
	return el
}

func (l *Ledger) all() []*EventLedger {
	l.mu.RLock()
	defer l.mu.RUnlock()
	out := make([]*EventLedger, 0, len(l.events))
	for _, el := range l.events {
		out = append(out, el)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}

func (el *EventLedger) example(x Example) {
	if len(el.examples) < maxExamples {
		el.examples = append(el.examples, x)
	}
}

// Grant records a hold the harness saw commit, with one returned row per seat.
func (el *EventLedger) Grant(h HoldInfo, rows []Row) {
	el.mu.Lock()
	defer el.mu.Unlock()
	hc := h
	hc.Seats = append([]int32(nil), h.Seats...)
	if len(rows) > 0 {
		hc.GrantNow = rows[0].Now
		if !rows[0].Expires.IsZero() {
			hc.Expires = rows[0].Expires
		}
	}
	el.holds[h.ID] = &hc
	for _, r := range rows {
		exp := r.Expires
		if exp.IsZero() {
			exp = hc.Expires
		}
		el.seats[r.Seat] = append(el.seats[r.Seat], seatEvent{kind: evGrant, hold: h.ID, now: r.Now, clock: r.Clock, expires: exp})
	}
}

// Extend records K1's payment window.
func (el *EventLedger) Extend(holdID int64, rows []Row) {
	el.mu.Lock()
	defer el.mu.Unlock()
	if h := el.holds[holdID]; h != nil && len(rows) > 0 {
		h.Expires = rows[0].Expires
	}
	for _, r := range rows {
		el.seats[r.Seat] = append(el.seats[r.Seat], seatEvent{kind: evExtend, hold: holdID, now: r.Now, clock: r.Clock, expires: r.Expires})
	}
}

// Release records seats the database said it released for a hold (the buyer's
// release or the sweeper's).
func (el *EventLedger) Release(holdID int64, rows []Row) {
	el.mu.Lock()
	defer el.mu.Unlock()
	for _, r := range rows {
		el.seats[r.Seat] = append(el.seats[r.Seat], seatEvent{kind: evRelease, hold: holdID, now: r.Now, clock: r.Clock})
	}
}

// Sale records a confirmation that committed. RETURNING order is not the order
// the seats were sent in, so tickets are matched by seat.
func (el *EventLedger) Sale(holdID int64, rows []Row, ticketOf map[int32]int64) {
	el.mu.Lock()
	defer el.mu.Unlock()
	for _, r := range rows {
		tid := ticketOf[r.Seat]
		el.seats[r.Seat] = append(el.seats[r.Seat], seatEvent{kind: evSale, hold: holdID, now: r.Now, clock: r.Clock, ticket: tid})
		el.tickets[tid] = r.Seat
	}
}

// Refund records a cancelled ticket.
func (el *EventLedger) Refund(ticketID int64, seat int32, now, clock time.Time) {
	el.mu.Lock()
	defer el.mu.Unlock()
	el.seats[seat] = append(el.seats[seat], seatEvent{kind: evRefund, now: now, clock: clock, ticket: ticketID})
	delete(el.tickets, ticketID)
}

// Ambiguous records a commit whose outcome is unknown.
func (el *EventLedger) Ambiguous() {
	el.mu.Lock()
	el.ambiguous++
	el.mu.Unlock()
}

// Reject classifies a confirmation (or checkout start) the design refused.
func (el *EventLedger) Reject(holdID int64, txnNow time.Time) string {
	el.mu.Lock()
	defer el.mu.Unlock()
	h := el.holds[holdID]
	if h == nil {
		el.rejEarly++
		el.example(Example{Kind: "rejection of an unknown hold", EventID: el.ID, HoldID: holdID})
		return RejEarly
	}
	m := h.Expires.Sub(txnNow)
	switch {
	case m <= 0:
		el.rejLate++
		return RejLate
	case m < el.guard:
		el.rejBoundary++
		return RejBoundary
	default:
		el.rejEarly++
		el.example(Example{Kind: "early rejection", EventID: el.ID, HoldID: holdID,
			Detail: fmt.Sprintf("refused %s before the hold's expiry", m.Round(time.Millisecond))})
		return RejEarly
	}
}

// Hold returns a copy of a hold the ledger knows.
func (el *EventLedger) Hold(id int64) (HoldInfo, bool) {
	el.mu.Lock()
	defer el.mu.Unlock()
	h := el.holds[id]
	if h == nil {
		return HoldInfo{}, false
	}
	c := *h
	c.Seats = append([]int32(nil), h.Seats...)
	return c, true
}

// Violations is the ledger's verdict over a set of events.
type Violations struct {
	Thefts             int64     `json:"thefts"`
	SoldSeatGrants     int64     `json:"grants_on_sold_seats"`
	DoubleSales        int64     `json:"double_sales"`
	SalesWithoutHold   int64     `json:"sales_without_the_hold"`
	LateSales          int64     `json:"late_sales"`
	RejectedLate       int64     `json:"rejected_late"`
	RejectedBoundary   int64     `json:"rejected_boundary"`
	RejectedEarly      int64     `json:"rejected_early"`
	Ambiguous          int64     `json:"ambiguous_commits,omitempty"`
	Examples           []Example `json:"examples,omitempty"`
	examplesByCategory map[string]int
}

// Count is the number that makes a correct design's phase fail: every class
// except late and boundary rejections, which are expected costs.
func (v *Violations) Count() int64 {
	return v.Thefts + v.SoldSeatGrants + v.DoubleSales + v.SalesWithoutHold + v.LateSales + v.RejectedEarly
}

func (v *Violations) add(o Violations) {
	v.Thefts += o.Thefts
	v.SoldSeatGrants += o.SoldSeatGrants
	v.DoubleSales += o.DoubleSales
	v.SalesWithoutHold += o.SalesWithoutHold
	v.LateSales += o.LateSales
	v.RejectedLate += o.RejectedLate
	v.RejectedBoundary += o.RejectedBoundary
	v.RejectedEarly += o.RejectedEarly
	v.Ambiguous += o.Ambiguous
	for _, x := range o.Examples {
		if v.examplesByCategory == nil {
			v.examplesByCategory = map[string]int{}
		}
		// Keep a few examples of every kind rather than the first five of one.
		if v.examplesByCategory[x.Kind] < 2 && len(v.Examples) < 3*maxExamples {
			v.examplesByCategory[x.Kind]++
			v.Examples = append(v.Examples, x)
		}
	}
}

// seatFinal is a seat's state after replaying its events.
type seatFinal struct {
	owner   int64 // hold id holding or having bought the seat; 0 = none
	sold    bool
	expires time.Time
}

// Evaluate replays every seat of the event. It is exact given the recorded
// clocks, and it does not depend on the order the harness recorded events in.
func (el *EventLedger) Evaluate() (Violations, map[int32]seatFinal) {
	el.mu.Lock()
	defer el.mu.Unlock()
	v := Violations{RejectedLate: el.rejLate, RejectedBoundary: el.rejBoundary, RejectedEarly: el.rejEarly,
		Ambiguous: el.ambiguous}
	ex := append([]Example(nil), el.examples...)
	final := make(map[int32]seatFinal, len(el.seats))

	seats := make([]int32, 0, len(el.seats))
	for s := range el.seats {
		seats = append(seats, s)
	}
	sort.Slice(seats, func(i, j int) bool { return seats[i] < seats[j] })

	for _, seat := range seats {
		evs := append([]seatEvent(nil), el.seats[seat]...)
		sort.SliceStable(evs, func(i, j int) bool { return evs[i].clock.Before(evs[j].clock) })
		var (
			owner    int64
			ownerExp time.Time
			released bool // the owner released the seat (or never had one)
			sold     bool
		)
		released = true
		note := func(kind string, hold, other int64, detail string) {
			if len(ex) < 4*maxExamples {
				ex = append(ex, Example{Kind: kind, EventID: el.ID, SeatID: seat, HoldID: hold, Other: other, Detail: detail})
			}
		}
		for _, e := range evs {
			switch e.kind {
			case evGrant:
				switch {
				case sold:
					v.SoldSeatGrants++
					note("grant on a sold seat", e.hold, owner, "")
				case !released && e.now.Before(ownerExp):
					v.Thefts++
					note("theft", e.hold, owner, fmt.Sprintf("granted %s before the previous hold expired",
						ownerExp.Sub(e.now).Round(time.Millisecond)))
				}
				owner, ownerExp, released, sold = e.hold, e.expires, false, false
			case evExtend:
				if owner == e.hold {
					ownerExp = e.expires
				}
			case evRelease:
				if owner == e.hold && !sold {
					owner, released = 0, true
				}
			case evSale:
				if e.hold != 0 {
					// The hold's current expiry: a K1 extension always precedes its sale.
					if h := el.holds[e.hold]; h != nil && !e.now.Before(h.Expires) {
						v.LateSales++
						note("late sale", e.hold, 0, fmt.Sprintf("sold %s after the hold expired", e.now.Sub(h.Expires).Round(time.Millisecond)))
					}
				}
				switch {
				case sold:
					v.DoubleSales++
					note("double sale", e.hold, owner, "")
				case e.hold != 0 && (owner != e.hold || released):
					v.SalesWithoutHold++
					note("sale without the hold", e.hold, owner, "")
				}
				owner, sold, released = e.hold, true, false
			case evRefund:
				sold, released, owner = false, true, 0
			}
		}
		if owner != 0 || sold {
			final[seat] = seatFinal{owner: owner, sold: sold, expires: ownerExp}
		}
	}
	for _, x := range ex {
		v.add(Violations{Examples: []Example{x}})
	}
	return v, final
}

// Tickets returns the ticket ids the harness saw created and not refunded.
func (el *EventLedger) Tickets() map[int64]int32 {
	el.mu.Lock()
	defer el.mu.Unlock()
	out := make(map[int64]int32, len(el.tickets))
	for k, s := range el.tickets {
		out[k] = s
	}
	return out
}

// LiveHolds returns, for every seat whose final owner is an unsold hold valid at
// t, that hold -- what a_valid_holds must return at quiescence.
func liveHoldsAt(final map[int32]seatFinal, t time.Time) map[int32]int64 {
	out := map[int32]int64{}
	for s, f := range final {
		if f.owner != 0 && !f.sold && f.expires.After(t) {
			out[s] = f.owner
		}
	}
	return out
}

// EvaluateAll sums the verdict over events (all of them when ids is empty).
func (l *Ledger) EvaluateAll(ids []int64) Violations {
	var total Violations
	if len(ids) == 0 {
		for _, el := range l.all() {
			v, _ := el.Evaluate()
			total.add(v)
		}
		return total
	}
	for _, id := range ids {
		if el := l.Event(id); el != nil {
			v, _ := el.Evaluate()
			total.add(v)
		}
	}
	return total
}
