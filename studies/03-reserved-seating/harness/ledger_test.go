package main

import (
	"testing"
	"time"
)

// The ledger decides whether a design broke an invariant, so it is tested on its
// own, with no database: every case the handoff names (§5.4), plus the lock-wait
// ordering that motivated recording clock_timestamp() beside now().

var t0 = time.Date(2026, 9, 14, 12, 0, 0, 0, time.UTC)

func at(d time.Duration) time.Time { return t0.Add(d) }

func newTestLedger() *EventLedger {
	v := buildVenue(1, 10)
	return &EventLedger{ID: 7, venue: v, guard: 2 * time.Minute, seats: map[int32][]seatEvent{},
		holds: map[int64]*HoldInfo{}, tickets: map[int64]int32{}}
}

func grant(el *EventLedger, hold int64, seat int32, now, clock, exp time.Time) {
	el.Grant(HoldInfo{ID: hold, EventID: el.ID, Section: 1, Seats: []int32{seat}, Customer: hold},
		[]Row{{Seat: seat, Now: now, Clock: clock, Expires: exp}})
}

func TestLegitimateStealAfterExpiry(t *testing.T) {
	el := newTestLedger()
	grant(el, 1, 3, at(0), at(0), at(40*time.Minute))
	grant(el, 2, 3, at(40*time.Minute), at(40*time.Minute+time.Millisecond), at(80*time.Minute))
	v, final := el.Evaluate()
	if v.Count() != 0 {
		t.Fatalf("a steal at expiry is legitimate, got %+v", v)
	}
	if final[3].owner != 2 {
		t.Fatalf("final owner = %d, want 2", final[3].owner)
	}
}

func TestTheftBeforeExpiry(t *testing.T) {
	el := newTestLedger()
	grant(el, 1, 3, at(0), at(0), at(40*time.Minute))
	grant(el, 2, 3, at(39*time.Minute), at(39*time.Minute), at(79*time.Minute))
	v, _ := el.Evaluate()
	if v.Thefts != 1 {
		t.Fatalf("thefts = %d, want 1 (%+v)", v.Thefts, v)
	}
}

func TestReleaseThenRehold(t *testing.T) {
	el := newTestLedger()
	grant(el, 1, 3, at(0), at(0), at(40*time.Minute))
	el.Release(1, []Row{{Seat: 3, Now: at(5 * time.Minute), Clock: at(5 * time.Minute)}})
	grant(el, 2, 3, at(6*time.Minute), at(6*time.Minute), at(46*time.Minute))
	if v, _ := el.Evaluate(); v.Count() != 0 {
		t.Fatalf("a hold after a release is legitimate, got %+v", v)
	}
}

// A grant that waited for a release's row lock carries a now() older than the
// release; its clock_timestamp() is later. Judged by now() alone it would look
// like theft.
func TestGrantThatWaitedForARelease(t *testing.T) {
	el := newTestLedger()
	grant(el, 1, 3, at(0), at(0), at(40*time.Minute))
	el.Release(1, []Row{{Seat: 3, Now: at(5*time.Minute + time.Millisecond), Clock: at(5*time.Minute + 2*time.Millisecond)}})
	grant(el, 2, 3, at(5*time.Minute), at(5*time.Minute+3*time.Millisecond), at(45*time.Minute))
	if v, _ := el.Evaluate(); v.Count() != 0 {
		t.Fatalf("the grant was ordered after the release by its clock, got %+v", v)
	}
}

func TestConfirmedThenHeld(t *testing.T) {
	el := newTestLedger()
	grant(el, 1, 3, at(0), at(0), at(40*time.Minute))
	el.Sale(1, []Row{{Seat: 3, Now: at(10 * time.Minute), Clock: at(10 * time.Minute)}}, map[int32]int64{3: 100})
	grant(el, 2, 3, at(50*time.Minute), at(50*time.Minute), at(90*time.Minute))
	v, _ := el.Evaluate()
	if v.SoldSeatGrants != 1 {
		t.Fatalf("grants on sold seats = %d, want 1 (%+v)", v.SoldSeatGrants, v)
	}
}

func TestRefundThenRehold(t *testing.T) {
	el := newTestLedger()
	grant(el, 1, 3, at(0), at(0), at(40*time.Minute))
	el.Sale(1, []Row{{Seat: 3, Now: at(10 * time.Minute), Clock: at(10 * time.Minute)}}, map[int32]int64{3: 100})
	el.Refund(100, 3, at(20*time.Minute), at(20*time.Minute))
	grant(el, 2, 3, at(21*time.Minute), at(21*time.Minute), at(61*time.Minute))
	if v, _ := el.Evaluate(); v.Count() != 0 {
		t.Fatalf("a hold after a refund is legitimate, got %+v", v)
	}
	if len(el.Tickets()) != 0 {
		t.Fatalf("a refunded ticket must leave the ledger's ticket set")
	}
}

func TestPaymentWindowExtension(t *testing.T) {
	el := newTestLedger()
	grant(el, 1, 3, at(0), at(0), at(40*time.Minute))
	el.Extend(1, []Row{{Seat: 3, Now: at(39 * time.Minute), Clock: at(39 * time.Minute), Expires: at(49 * time.Minute)}})
	el.Sale(1, []Row{{Seat: 3, Now: at(42 * time.Minute), Clock: at(42 * time.Minute)}}, map[int32]int64{3: 100})
	if v, _ := el.Evaluate(); v.Count() != 0 {
		t.Fatalf("a sale inside the payment window is legitimate, got %+v", v)
	}
	// A steal between the old and the new expiry is theft.
	el2 := newTestLedger()
	grant(el2, 1, 3, at(0), at(0), at(40*time.Minute))
	el2.Extend(1, []Row{{Seat: 3, Now: at(39 * time.Minute), Clock: at(39 * time.Minute), Expires: at(49 * time.Minute)}})
	grant(el2, 2, 3, at(41*time.Minute), at(41*time.Minute), at(81*time.Minute))
	if v, _ := el2.Evaluate(); v.Thefts != 1 {
		t.Fatalf("a steal inside the payment window is theft, got %+v", v)
	}
}

func TestLateSale(t *testing.T) {
	el := newTestLedger()
	grant(el, 1, 3, at(0), at(0), at(40*time.Minute))
	el.Sale(1, []Row{{Seat: 3, Now: at(41 * time.Minute), Clock: at(41 * time.Minute)}}, map[int32]int64{3: 100})
	v, _ := el.Evaluate()
	if v.LateSales != 1 {
		t.Fatalf("late sales = %d, want 1 (%+v)", v.LateSales, v)
	}
}

// K0's failure: a naive confirmation sells a seat another hold took after expiry.
func TestNaiveSaleOverAnotherHold(t *testing.T) {
	el := newTestLedger()
	grant(el, 1, 3, at(0), at(0), at(40*time.Minute))
	grant(el, 2, 3, at(41*time.Minute), at(41*time.Minute), at(81*time.Minute))
	el.Sale(1, []Row{{Seat: 3, Now: at(42 * time.Minute), Clock: at(42 * time.Minute)}}, map[int32]int64{3: 100})
	el.Sale(2, []Row{{Seat: 3, Now: at(43 * time.Minute), Clock: at(43 * time.Minute)}}, map[int32]int64{3: 101})
	v, _ := el.Evaluate()
	if v.LateSales != 1 || v.SalesWithoutHold != 1 || v.DoubleSales != 1 {
		t.Fatalf("want 1 late sale, 1 sale without the hold, 1 double sale; got %+v", v)
	}
}

func TestRejectionClasses(t *testing.T) {
	el := newTestLedger()
	grant(el, 1, 3, at(0), at(0), at(40*time.Minute))
	exp := at(40 * time.Minute)
	cases := []struct {
		txn  time.Time
		want string
	}{
		{exp.Add(time.Millisecond), RejLate},
		{exp, RejLate},
		{exp.Add(-time.Millisecond), RejBoundary},
		{exp.Add(-2*time.Minute + time.Millisecond), RejBoundary},
		{exp.Add(-2 * time.Minute), RejEarly},
	}
	for _, c := range cases {
		if got := el.Reject(1, c.txn, "", false, true); got != c.want {
			t.Errorf("margin %s: class %s, want %s", exp.Sub(c.txn), got, c.want)
		}
	}
	v, _ := el.Evaluate()
	if v.RejectedLate != 2 || v.RejectedBoundary != 2 || v.RejectedEarly != 1 || v.Count() != 1 {
		t.Fatalf("tallies %+v", v)
	}
}

func TestPartialHoldIsVisibleAsLiveSeats(t *testing.T) {
	el := newTestLedger()
	el.Grant(HoldInfo{ID: 1, EventID: el.ID, Section: 1, Seats: []int32{3, 4}, Customer: 1},
		[]Row{{Seat: 3, Now: at(0), Clock: at(0), Expires: at(40 * time.Minute)}, {Seat: 4, Now: at(0), Clock: at(0), Expires: at(40 * time.Minute)}})
	_, final := el.Evaluate()
	live := liveHoldsAt(final, at(time.Minute))
	if live[3] != 1 || live[4] != 1 || len(live) != 2 {
		t.Fatalf("the audit compares every seat of a live hold, got %v", live)
	}
}

func TestAmbiguousCommitsAreCounted(t *testing.T) {
	el := newTestLedger()
	el.Ambiguous()
	el.Ambiguous()
	if v, _ := el.Evaluate(); v.Ambiguous != 2 {
		t.Fatalf("ambiguous = %d, want 2", v.Ambiguous)
	}
}

func TestLoadedStateSeedsTheLedger(t *testing.T) {
	ds, err := Generate("tiny", 42, nil)
	if err != nil {
		t.Fatal(err)
	}
	l := NewLedger(ds, t0, 2*time.Minute)
	v := l.EvaluateAll(nil)
	if v.Count() != 0 {
		t.Fatalf("the loaded state must be consistent, got %+v", v)
	}
	var live int
	for _, el := range l.all() {
		_, final := el.Evaluate()
		live += len(liveHoldsAt(final, t0))
	}
	want := 0
	for _, h := range ds.Holds {
		if h.live() {
			want += len(h.Seats)
		}
	}
	if live != want {
		t.Fatalf("live held seats at load: %d, want %d", live, want)
	}
}

// AM-01: a transient early rejection is still an early rejection (and a violation);
// the transient flag means nothing for late and boundary rejections.
func TestTransientEarlyRejectionCountsInBothTotals(t *testing.T) {
	el := newTestLedger()
	grant(el, 1, 3, at(0), at(0), at(40*time.Minute))
	exp := at(40 * time.Minute)
	if got := el.Reject(1, exp.Add(-10*time.Minute), "", true, true); got != RejEarly {
		t.Fatalf("class %s, want %s", got, RejEarly)
	}
	el.Reject(1, exp.Add(-10*time.Minute), "", false, true)
	v, _ := el.Evaluate()
	if v.RejectedEarly != 2 || v.RejectedEarlyTransient != 1 || v.Count() != 2 {
		t.Fatalf("tallies %+v", v)
	}
}

func TestTransientFlagIgnoredForLateAndBoundary(t *testing.T) {
	el := newTestLedger()
	grant(el, 1, 3, at(0), at(0), at(40*time.Minute))
	exp := at(40 * time.Minute)
	el.Reject(1, exp.Add(time.Second), "", true, true)
	el.Reject(1, exp.Add(-time.Second), "", true, true)
	v, _ := el.Evaluate()
	if v.RejectedLate != 1 || v.RejectedBoundary != 1 || v.RejectedEarly != 0 || v.RejectedEarlyTransient != 0 || v.Count() != 0 {
		t.Fatalf("tallies %+v", v)
	}
}
