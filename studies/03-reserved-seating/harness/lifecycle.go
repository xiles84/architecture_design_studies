package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"adsplatform/core/measure"
	"adsplatform/ports"
)

// ---------------------------------------------------------------------------
// The lifecycle: the hold guarantee under compressed time.
//
// Durations are defined in HUMAN time -- a 40-minute hold, a basket of up to 45
// minutes, a payment of up to 4 minutes -- and run with one human minute lasting
// -human-minute of wall-clock time (50 ms on PostgreSQL: a hold lasts 2 s).
// Scaling every duration by the same factor keeps their ratios, and therefore the
// share of payments that straddle an expiry.
//
// Customers hold seats they chose, spend basket time, abandon a quarter of their
// holds (some explicitly, most by walking away), check the hold, pay, and
// confirm. A sweeper runs for every design. In the first event of each tier it is
// STOPPED from minute 45 to minute 85; at minute 85 every buyer is paused between
// database operations and a probe tries to hold each seat whose hold expired at
// least two minutes earlier. A lazy design must grant all of them; E1 -- whose
// expiry only takes effect when the sweeper runs -- is the control that must not.
//
// At the end, once every unconfirmed hold has expired and the sweeper has
// finished, a leak probe checks that every unsold seat can be held (INV-4).
// ---------------------------------------------------------------------------

type LifecycleResult struct {
	Trial         int     `json:"trial"`
	Tier          int     `json:"tier"`
	Events        int     `json:"events"`
	Buyers        int     `json:"buyers"`
	HumanMinuteMS float64 `json:"human_minute_ms"`

	Seats             int64  `json:"seats"`
	Sold              int64  `json:"sold"`
	Customers         int64  `json:"customers"`
	HoldsGranted      int64  `json:"holds_granted"`
	Conflicts         int64  `json:"conflicts"`
	NoBlockWaits      int64  `json:"no_block_waits"`
	AbandonedSilent   int64  `json:"abandoned_silently"`
	AbandonedExplicit int64  `json:"abandoned_explicitly"`
	ExpiredAtCheck    int64  `json:"expired_at_precheck"`
	ExpiredAtCheckout int64  `json:"expired_at_checkout,omitempty"`
	PaymentsStarted   int64  `json:"payments_started"`
	ConfirmedHolds    int64  `json:"confirmed_holds"`
	ConfirmedSeats    int64  `json:"confirmed_seats"`
	RejectedLate      int64  `json:"rejected_late"`
	RejectedBoundary  int64  `json:"rejected_boundary"`
	RejectedEarly     int64  `json:"rejected_early"`
	// RejectedEarlyTransient: the part of RejectedEarly the ledger classed transient (AM-01.1).
	RejectedEarlyTransient int64 `json:"rejected_early_transient"`
	// S1r only (AM-01.3): confirmations run a second time after a short match, and
	// how many of those sold.
	ConfirmRetries        int64 `json:"confirm_retries"`
	ConfirmRetrySuccesses int64 `json:"confirm_retry_successes"`
	// AM-03.1: reads of the sold-seat monitor the engine did not answer at once.
	// MonitorTimeouts counts events whose phase ended early because it never did.
	MonitorRetries  int64 `json:"monitor_retries"`
	MonitorTimeouts int64 `json:"monitor_timeouts"`
	EngineRetries     int64  `json:"engine_retries"`
	Errors            int64  `json:"errors"`
	FirstError        string `json:"first_error,omitempty"`

	SweptSeats           int64                `json:"swept_seats"`
	ReleaseLagHumanMin   measure.LatencyStats `json:"release_lag_human_min"`
	OutageEvents         int                  `json:"outage_events"`
	OutageProbed         int64                `json:"outage_probed"`
	OutageUnavailable    int64                `json:"outage_unavailable_seats"`
	LeakProbed           int64                `json:"leak_probed"`
	Leaked               int64                `json:"leaked_seats"`
	IdleHeldSeatHumanMin float64              `json:"idle_held_seat_human_minutes"`
	TimedOut             int                  `json:"timed_out_events"`
	WallS                float64              `json:"wall_s"`
	ConfirmedSeatsPerSec float64              `json:"confirmed_seats_per_sec"`
	HoldLatency          measure.LatencyStats `json:"hold_latency"`
	ConfirmLatency       measure.LatencyStats `json:"confirm_latency"`
	CheckoutLatency      measure.LatencyStats `json:"checkout_latency,omitempty"`

	Audit *Audit `json:"audit,omitempty"`
}

type lcOne struct {
	customers, holds, conflicts, noBlock, abSilent, abExplicit atomic.Int64
	expCheck, expCheckout, payments, confHolds, confSeats      atomic.Int64
	late, boundary, early, retries, errs, swept                atomic.Int64
	confRetries, confRetrySold                                 atomic.Int64
	monRetries, monTimeouts                                    atomic.Int64
	idleMicroHuman                                             atomic.Int64
	outageProbed, outageUnavailable, leakProbed, leaked        int64
	holdL, confL, coL, lagL                                    []time.Duration
	firstErr                                                   string
	timedOut, outage                                           bool
	wall                                                       time.Duration
}

func RunLifecycle(ctx context.Context, db ports.DB, sl *Seller, d Design, ds *Dataset, s Settings, trial int) ([]LifecycleResult, error) {
	var out []LifecycleResult
	for _, tier := range ds.tiersOf("lifecycle") {
		if len(s.LifecycleTiers) > 0 && !containsInt(s.LifecycleTiers, tier) {
			continue
		}
		evs := ds.eventsOf("lifecycle", tier)
		if s.LifecycleEvents > 0 && len(evs) > s.LifecycleEvents {
			evs = evs[:s.LifecycleEvents]
		}
		res := LifecycleResult{Trial: trial, Tier: tier, Buyers: s.LifecycleBuyers,
			HumanMinuteMS: float64(s.HumanMinute.Microseconds()) / 1000}
		var holdL, confL, coL, lagL []time.Duration
		var wall time.Duration
		var ids []int64
		for i, ev := range evs {
			one, err := lifecycleEvent(ctx, sl, d, ev, s, i == 0)
			if err != nil {
				return nil, err
			}
			sold, err := soldTolerantly(ctx, db, d, ev.ID, &one.monRetries)
			if err != nil {
				// AM-03.1: the tier keeps the events it measured.
				one.monTimeouts.Add(1)
				sold = 0
			}
			ids = append(ids, ev.ID)
			res.Events++
			res.Seats += int64(ds.venue(ev.VenueID).Capacity)
			res.Sold += sold
			res.Customers += one.customers.Load()
			res.HoldsGranted += one.holds.Load()
			res.Conflicts += one.conflicts.Load()
			res.NoBlockWaits += one.noBlock.Load()
			res.AbandonedSilent += one.abSilent.Load()
			res.AbandonedExplicit += one.abExplicit.Load()
			res.ExpiredAtCheck += one.expCheck.Load()
			res.ExpiredAtCheckout += one.expCheckout.Load()
			res.PaymentsStarted += one.payments.Load()
			res.ConfirmedHolds += one.confHolds.Load()
			res.ConfirmedSeats += one.confSeats.Load()
			res.RejectedLate += one.late.Load()
			res.RejectedBoundary += one.boundary.Load()
			res.RejectedEarly += one.early.Load()
			res.ConfirmRetries += one.confRetries.Load()
			res.ConfirmRetrySuccesses += one.confRetrySold.Load()
			res.MonitorRetries += one.monRetries.Load()
			res.MonitorTimeouts += one.monTimeouts.Load()
			res.EngineRetries += one.retries.Load()
			res.Errors += one.errs.Load()
			res.SweptSeats += one.swept.Load()
			res.IdleHeldSeatHumanMin += float64(one.idleMicroHuman.Load()) / float64(time.Minute.Microseconds())
			if one.outage {
				res.OutageEvents++
				res.OutageProbed += one.outageProbed
				res.OutageUnavailable += one.outageUnavailable
			}
			res.LeakProbed += one.leakProbed
			res.Leaked += one.leaked
			if one.timedOut {
				res.TimedOut++
			}
			if res.FirstError == "" {
				res.FirstError = one.firstErr
			}
			wall += one.wall
			holdL = append(holdL, one.holdL...)
			confL = append(confL, one.confL...)
			coL = append(coL, one.coL...)
			lagL = append(lagL, one.lagL...)
		}
		res.WallS = wall.Seconds()
		if res.WallS > 0 {
			res.ConfirmedSeatsPerSec = float64(res.ConfirmedSeats) / res.WallS
		}
		res.HoldLatency = measure.Summarize(holdL)
		res.ConfirmLatency = measure.Summarize(confL)
		res.CheckoutLatency = measure.Summarize(coL)
		res.ReleaseLagHumanMin = humanLags(lagL, s)

		au, err := RunAudit(ctx, db, d, sl.world, fmt.Sprintf("lifecycle@%d#%d", tier, trial), ids)
		if err != nil {
			return nil, err
		}
		res.Audit = au
		res.RejectedEarlyTransient = au.Ledger.RejectedEarlyTransient
		status := "ok"
		if au.Violations() > 0 {
			status = "VIOLATIONS"
		}
		if res.Leaked > 0 {
			status += fmt.Sprintf(", LEAKED %d seats", res.Leaked)
		}
		fmt.Printf("    lifecycle %5d seats x%-3d %7.1f confirmed seats/s  holds=%d abandoned=%d+%d expired@check=%d rejected late/boundary/early(transient)=%d/%d/%d(%d)  confirm retries sold/run=%d/%d  outage unavailable=%d/%d  leaked=%d errors=%d  monitor retries/ended early=%d/%d  %s\n",
			tier, res.Events, res.ConfirmedSeatsPerSec, res.HoldsGranted, res.AbandonedSilent, res.AbandonedExplicit,
			res.ExpiredAtCheck, res.RejectedLate, res.RejectedBoundary, res.RejectedEarly, res.RejectedEarlyTransient,
			res.ConfirmRetrySuccesses, res.ConfirmRetries, res.OutageUnavailable,
			res.OutageProbed, res.Leaked, res.Errors, res.MonitorRetries, res.MonitorTimeouts, status)
		fmt.Printf("      audit: %s\n", au)
		if res.FirstError != "" {
			fmt.Printf("      first error: %s\n", res.FirstError)
		}
		out = append(out, res)
	}
	return out, nil
}

// humanLags converts wall-clock lags into human minutes, in the LatencyStats
// fields (whose "ms" names then mean human minutes).
// monitorBudget bounds how long the lifecycle's sold-seat monitor retries a read
// that the engine did not answer (AM-03.1). It is generous enough to outlast the
// statement timeouts a saturated node produces, and short enough that a database
// which has genuinely stopped answering is reported rather than waited on.
const monitorBudget = 30 * time.Second

// soldTolerantly reads the event's sold count, retrying a failed read within the
// budget. Retries are counted so a reader can see when the engine stopped
// answering; the returned error means it never did.
func soldTolerantly(ctx context.Context, db ports.DB, d Design, ev int64, retries *atomic.Int64) (int64, error) {
	sold, err := eventSold(ctx, db, d, ev)
	if err == nil {
		return sold, nil
	}
	deadline := time.Now().Add(monitorBudget)
	backoff := 100 * time.Millisecond
	for err != nil && time.Now().Before(deadline) && ctx.Err() == nil {
		retries.Add(1)
		time.Sleep(backoff)
		backoff = min(2*backoff, 2*time.Second)
		sold, err = eventSold(ctx, db, d, ev)
	}
	return sold, err
}

func humanLags(lags []time.Duration, s Settings) measure.LatencyStats {
	scaled := make([]time.Duration, len(lags))
	for i, l := range lags {
		// One human minute becomes one millisecond, so the *_ms fields read as minutes.
		scaled[i] = time.Duration(float64(l) / float64(s.HumanMinute) * float64(time.Millisecond))
	}
	return measure.Summarize(scaled)
}

// basketTime: 70% U[1,10] min, 20% U[10,30] min, 10% U[30,45] min (human).
func basketTime(r *rand.Rand) time.Duration {
	u := func(lo, hi time.Duration) time.Duration { return lo + time.Duration(r.Int63n(int64(hi-lo))) }
	x := r.Intn(100)
	switch {
	case x < 70:
		return u(time.Minute, 10*time.Minute)
	case x < 90:
		return u(10*time.Minute, 30*time.Minute)
	default:
		return u(30*time.Minute, 45*time.Minute)
	}
}

func paymentTime(r *rand.Rand) time.Duration {
	return 30*time.Second + time.Duration(r.Int63n(int64(4*time.Minute-30*time.Second)))
}

func lifecycleEvent(ctx context.Context, sl *Seller, d Design, ev *Event, s Settings, withOutage bool) (*lcOne, error) {
	one := &lcOne{outage: withOutage && s.LcOutageTo > s.LcOutageFrom}
	capacity := int64(sl.world.ds.venue(ev.VenueID).Capacity)
	var (
		mu       sync.Mutex
		gate     sync.RWMutex // buyers read-lock around database operations; the probe write-locks
		wg       sync.WaitGroup
		stop     atomic.Bool
		sweeping atomic.Bool
		// probeDone ends the outage. The outage ends when the probe has run, not at a
		// wall-clock minute: the first dev check had the sweeper resume on its own tick
		// and release every expired seat before the probe could look at them.
		probeDone atomic.Bool
		errOnce   sync.Once
	)
	recordErr := func(err error) {
		one.errs.Add(1)
		errOnce.Do(func() { mu.Lock(); one.firstErr = err.Error(); mu.Unlock() })
	}
	h := s.human
	start := time.Now()
	deadline := start.Add(h(s.LcTimeout))
	ctx = WithAttemptDeadline(ctx, deadline.Add(10*time.Second))
	sleep := func(dur time.Duration) {
		t := time.NewTimer(dur)
		defer t.Stop()
		select {
		case <-t.C:
		case <-ctx.Done():
		}
	}

	for w := 0; w < s.LifecycleBuyers; w++ {
		wg.Add(1)
		go func(w int, seed int64) {
			defer wg.Done()
			r := rand.New(rand.NewSource(seed))
			buyer := &Buyer{seller: sl, node: w % 2, r: r, maxConflicts: s.MaxConflicts, gate: &gate}
			var hl, cl, col []time.Duration
			for !stop.Load() {
				cust := sl.world.randomCustomer(r)
				n := partySize(r)
				one.customers.Add(1)
				a, err := buyer.Acquire(ctx, ev, cust, n, h(s.LcTTL))
				one.conflicts.Add(int64(a.Conflicts))
				one.retries.Add(int64(a.EngineRetries))
				if err != nil {
					if !errors.Is(err, errAttemptDeadline) {
						recordErr(err)
					}
					sleep(h(time.Minute))
					continue
				}
				if !a.Granted {
					one.noBlock.Add(1)
					sleep(h(time.Minute)) // sold out for now: holds expire, try again soon
					continue
				}
				one.holds.Add(1)
				hl = append(hl, a.HoldLatency)
				blk, hold := a.Block, a.Hold.HoldID
				basket := basketTime(r)
				sleep(h(basket))
				if r.Intn(100) < s.LcAbandonPct {
					if r.Intn(100) < s.LcAbandonExplicit {
						one.abExplicit.Add(1)
						one.idleMicroHuman.Add(int64(n) * basket.Microseconds())
						var err error
						buyer.do(func() { _, err = sl.Release(ctx, buyer.node, blk, hold) })
						if err != nil {
							recordErr(err)
						}
					} else {
						one.abSilent.Add(1)
						one.idleMicroHuman.Add(int64(n) * s.LcTTL.Microseconds())
					}
					continue
				}
				var valid bool
				buyer.do(func() { valid, err = sl.CheckHold(ctx, buyer.node, blk, hold) })
				if err != nil {
					recordErr(err)
					continue
				}
				if !valid {
					one.expCheck.Add(1)
					continue
				}
				if d.PaymentWindow {
					t0 := time.Now()
					var cr ConfirmResult
					buyer.do(func() { cr, err = sl.BeginCheckout(ctx, buyer.node, blk, hold, h(s.LcPaymentWindow)) })
					if err != nil {
						recordErr(err)
						continue
					}
					if !cr.Sold {
						one.expCheckout.Add(1)
						countClass(one, cr.Class)
						continue
					}
					col = append(col, time.Since(t0))
				}
				one.payments.Add(1)
				sleep(h(paymentTime(r)))
				t1 := time.Now()
				var cr ConfirmResult
				buyer.do(func() { cr, err = sl.Confirm(ctx, buyer.node, blk, hold, cust) })
				one.retries.Add(int64(cr.Retries))
				if cr.ShortRetried {
					one.confRetries.Add(1)
				}
				if cr.RetrySold {
					one.confRetrySold.Add(1)
				}
				if err != nil {
					recordErr(err)
					continue
				}
				if cr.Sold {
					one.confHolds.Add(1)
					one.confSeats.Add(int64(n))
					cl = append(cl, time.Since(t1))
				} else {
					countClass(one, cr.Class)
				}
			}
			mu.Lock()
			one.holdL = append(one.holdL, hl...)
			one.confL = append(one.confL, cl...)
			one.coL = append(one.coL, col...)
			mu.Unlock()
		}(w, ev.ID*6007+int64(w)*15485863)
	}

	sweepNode := 0
	if d.AppClock {
		sweepNode = 1 // E0: the sweeper runs on the node whose clock runs ahead
	}
	sweepOnce := func() int {
		var res SweepResult
		var err error
		gate.RLock()
		res, err = sl.Sweep(ctx, sweepNode, s.LcSweepBatch)
		gate.RUnlock()
		if err != nil {
			recordErr(fmt.Errorf("sweeper: %w", err))
			return 0
		}
		one.swept.Add(int64(res.Released))
		mu.Lock()
		one.lagL = append(one.lagL, res.Lags...)
		mu.Unlock()
		return res.Released
	}
	sweepStop := make(chan struct{})
	sweepDone := make(chan struct{})
	sweeping.Store(true)
	go func() {
		defer close(sweepDone)
		tick := time.NewTicker(h(s.LcSweepEvery))
		defer tick.Stop()
		for {
			select {
			case <-sweepStop:
				return
			case <-tick.C:
			}
			if one.outage {
				el := time.Since(start)
				if el >= h(s.LcOutageFrom) && !probeDone.Load() {
					continue
				}
			}
			sweepOnce()
		}
	}()

	// Monitor: sold out, timeout, and the outage probe.
	probed := !one.outage
	for {
		time.Sleep(max(50*time.Millisecond, s.HumanMinute))
		if !probed && time.Since(start) >= h(s.LcOutageTo) {
			probed = true
			gate.Lock() // every buyer is now between database operations
			one.outageProbed, one.outageUnavailable = outageProbe(ctx, sl, ev, s)
			probeDone.Store(true)
			gate.Unlock()
		}
		sold, err := soldTolerantly(ctx, sl.db, d, ev.ID, &one.monRetries)
		if err != nil {
			// AM-03.1: the monitor watches the workload; it is not part of any design.
			// A design whose statements are retried must not lose its phase because an
			// observer query was not. The event ends here and is recorded as such; the
			// tier and the cell continue.
			one.monTimeouts.Add(1)
			recordErr(fmt.Errorf("lifecycle monitor event %d, phase of this event ended early: %w", ev.ID, err))
			break
		}
		if sold >= capacity && probed {
			break
		}
		if time.Now().After(deadline) {
			one.timedOut = true
			break
		}
	}
	stop.Store(true)
	wg.Wait()
	one.wall = time.Since(start)
	close(sweepStop)
	<-sweepDone

	// Settle: every unconfirmed hold expires, the sweeper finishes, then probe.
	el := sl.world.ledger.Event(ev.ID)
	_, final := el.Evaluate()
	var latest time.Time
	for _, f := range final {
		if !f.sold && f.owner != 0 && f.expires.After(latest) {
			latest = f.expires
		}
	}
	var dbNow time.Time
	if err := sl.db.QueryRow(ctx, "SELECT now()").Scan(&dbNow); err != nil {
		return nil, err
	}
	if wait := latest.Sub(dbNow); wait > 0 {
		sleep(wait + 50*time.Millisecond)
	}
	for sweepOnce() > 0 {
	}
	one.leakProbed, one.leaked = probeLeaks(ctx, sl, ev, s.LcProbeCapPerEvent, h(time.Minute), nil)
	return one, nil
}

func countClass(one *lcOne, class string) {
	switch class {
	case RejLate:
		one.late.Add(1)
	case RejBoundary:
		one.boundary.Add(1)
	default:
		one.early.Add(1)
	}
}

// outageProbe: with every buyer paused and the sweeper stopped, try to hold each
// seat whose latest hold expired at least two human minutes ago and was neither
// released nor confirmed. A lazy design grants every one.
func outageProbe(ctx context.Context, sl *Seller, ev *Event, s Settings) (probed, unavailable int64) {
	var now time.Time
	if err := sl.db.QueryRow(ctx, "SELECT now()").Scan(&now); err != nil {
		return 0, 0
	}
	_, final := sl.world.ledger.Event(ev.ID).Evaluate()
	var seats []int32
	for seat, f := range final {
		if !f.sold && f.owner != 0 && f.expires.Before(now.Add(-s.human(2*time.Minute))) {
			seats = append(seats, seat)
		}
	}
	sort.Slice(seats, func(i, j int) bool { return seats[i] < seats[j] })
	if len(seats) > s.LcProbeCapPerEvent {
		seats = seats[:s.LcProbeCapPerEvent]
	}
	if len(seats) == 0 {
		return 0, 0
	}
	return probeLeaks(ctx, sl, ev, s.LcProbeCapPerEvent, s.human(time.Minute), seats)
}
