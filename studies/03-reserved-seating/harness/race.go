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
// The race: a hot drop with seat choice.
//
// A crowd of buyers arrives at an unsold event at once. Each buyer is a queue of
// customers: a customer reads the seat map, chooses the best block it can see,
// holds it and confirms; on a conflict it re-reads and chooses again. Holds use
// the REAL 40-minute TTL.
//
// One customer in ten who gets a hold does not confirm straight away: they keep
// the hold while the rest of the crowd hunts for seats around them, and pay only
// once the crowd has gone. That is the owner's scenario taken literally -- choose
// seats, keep them through a sell-out storm, pay at the end -- and every one of
// those confirmations must succeed (INV-3).
//
// After the crowd and the deferred confirmations, one last customer buys single
// seats until nothing is available. Only then are seats left unsold counted as
// under-selling (INV-8), and a leak probe checks that any unsold seat without a
// valid hold can be held (INV-4).
// ---------------------------------------------------------------------------

type RaceResult struct {
	Trial         int `json:"trial"`
	Tier          int `json:"tier"`
	Events        int `json:"events"`
	EventsPlanned int `json:"events_planned"`
	Buyers        int `json:"buyers"`

	Seats           int64  `json:"seats"`
	Sold            int64  `json:"sold"`
	Customers       int64  `json:"customers"`
	HoldsGranted    int64  `json:"holds_granted"`
	SeatsGranted    int64  `json:"seats_granted"`
	HoldAttempts    int64  `json:"hold_attempts"`
	Conflicts       int64  `json:"conflicts"`
	MapReads        int64  `json:"map_reads"`
	SectionReads    int64  `json:"section_reads"`
	EngineRetries   int64  `json:"engine_retries"`
	GaveUp          int64  `json:"gave_up"`
	NoBlock         int64  `json:"no_block"`
	Errors          int64  `json:"errors"`
	FirstError      string `json:"first_error,omitempty"`
	ImmediateFailed int64  `json:"immediate_confirms_refused"`
	// S1r only (AM-01.3): confirmations run a second time after a short match, and
	// how many of those sold.
	ConfirmRetries        int64 `json:"confirm_retries"`
	ConfirmRetrySuccesses int64 `json:"confirm_retry_successes"`

	DeferredHolds     int64 `json:"deferred_holds"`
	DeferredConfirmed int64 `json:"deferred_confirmed"`
	DeferredRefused   int64 `json:"deferred_refused"`

	SweepSold      int64 `json:"final_sweep_sold"`
	UnderSoldSeats int64 `json:"undersold_seats"`
	UnderSold      int   `json:"undersold_events"`
	LeakProbed     int64 `json:"leak_probed"`
	Leaked         int64 `json:"leaked_seats"`

	// AllocateS sums, over the tier's events, the time from the gate to the last
	// hold granted; SeatsPerSec is seats granted over it.
	AllocateS        float64              `json:"allocate_s"`
	WallS            float64              `json:"wall_s"`
	SeatsPerSec      float64              `json:"seats_per_sec"`
	ConflictsPerHold float64              `json:"conflicts_per_hold"`
	MapReadsPerHold  float64              `json:"map_reads_per_hold"`
	HoldLatency      measure.LatencyStats `json:"hold_latency"`
	ConfirmLatency   measure.LatencyStats `json:"confirm_latency"`
	DeferredLatency  measure.LatencyStats `json:"deferred_confirm_latency"`
	CustomerLatency  measure.LatencyStats `json:"customer_latency"`
	TimedOut         int                  `json:"timed_out_events"`
	TimedOutSoldPct  float64              `json:"timed_out_sold_pct,omitempty"`

	Audit *Audit `json:"audit,omitempty"`

	// Filled only when a report combines trials. Never written by the harness.
	TrialsSeatsPerSec []float64 `json:"-"`
	SpreadPct         float64   `json:"-"`
}

type raceOne struct {
	customers, holds, seats, attempts, conflicts, mapReads, secReads atomic.Int64
	retries, gaveUp, noBlock, errs, immFailed                        atomic.Int64
	deferredConfirmed, deferredRefused                               atomic.Int64
	confRetries, confRetrySold                                       atomic.Int64
	holdL, confL, defL, custL                                        []time.Duration
	deferred                                                         []deferredHold
	firstErr                                                         string
	timedOut                                                         bool
	sweepSold, underSold, leakProbed, leaked                         int64
	wall, allocate                                                   time.Duration
}

// noteConfirm counts S1r's second attempts (zero for every other design).
func (one *raceOne) noteConfirm(cr ConfirmResult) {
	if cr.ShortRetried {
		one.confRetries.Add(1)
	}
	if cr.RetrySold {
		one.confRetrySold.Add(1)
	}
}

type deferredHold struct {
	node     int
	block    Block
	hold     int64
	customer int64
}

func RunRaces(ctx context.Context, db ports.DB, sl *Seller, d Design, ds *Dataset, s Settings, trial int) ([]RaceResult, error) {
	var out []RaceResult
	for _, tier := range ds.tiersOf("race") {
		if len(s.RaceTiers) > 0 && !containsInt(s.RaceTiers, tier) {
			continue
		}
		res := RaceResult{Trial: trial, Tier: tier, Buyers: s.RaceBuyers}
		var holdL, confL, defL, custL []time.Duration
		var wall, allocate time.Duration
		var tOutSold, tOutSeats int64
		planned := ds.eventsOf("race", tier)
		res.EventsPlanned = len(planned)
		var ids []int64
		tierStart := time.Now()
		for _, ev := range planned {
			if s.RaceTierBudget > 0 && time.Since(tierStart) > s.RaceTierBudget {
				break
			}
			one := raceEvent(ctx, sl, ev, s)
			sold, err := eventSold(ctx, db, d, ev.ID)
			if err != nil {
				return nil, fmt.Errorf("race event %d: final count: %w", ev.ID, err)
			}
			ids = append(ids, ev.ID)
			capacity := int64(ds.venue(ev.VenueID).Capacity)
			res.Events++
			res.Seats += capacity
			res.Sold += sold
			res.Customers += one.customers.Load()
			res.HoldsGranted += one.holds.Load()
			res.SeatsGranted += one.seats.Load()
			res.HoldAttempts += one.attempts.Load()
			res.Conflicts += one.conflicts.Load()
			res.MapReads += one.mapReads.Load()
			res.SectionReads += one.secReads.Load()
			res.EngineRetries += one.retries.Load()
			res.GaveUp += one.gaveUp.Load()
			res.NoBlock += one.noBlock.Load()
			res.Errors += one.errs.Load()
			res.ImmediateFailed += one.immFailed.Load()
			res.ConfirmRetries += one.confRetries.Load()
			res.ConfirmRetrySuccesses += one.confRetrySold.Load()
			res.DeferredHolds += int64(len(one.deferred))
			res.DeferredConfirmed += one.deferredConfirmed.Load()
			res.DeferredRefused += one.deferredRefused.Load()
			res.SweepSold += one.sweepSold
			res.LeakProbed += one.leakProbed
			res.Leaked += one.leaked
			if res.FirstError == "" {
				res.FirstError = one.firstErr
			}
			if one.timedOut {
				res.TimedOut++
				tOutSold += sold
				tOutSeats += capacity
			} else if one.underSold > 0 {
				res.UnderSold++
				res.UnderSoldSeats += one.underSold
			}
			wall += one.wall
			allocate += one.allocate
			holdL = append(holdL, one.holdL...)
			confL = append(confL, one.confL...)
			defL = append(defL, one.defL...)
			custL = append(custL, one.custL...)
		}
		res.WallS, res.AllocateS = wall.Seconds(), allocate.Seconds()
		if res.AllocateS > 0 {
			res.SeatsPerSec = float64(res.SeatsGranted) / res.AllocateS
		}
		if res.HoldsGranted > 0 {
			res.ConflictsPerHold = float64(res.Conflicts) / float64(res.HoldsGranted)
			res.MapReadsPerHold = float64(res.MapReads) / float64(res.HoldsGranted)
		}
		if tOutSeats > 0 {
			res.TimedOutSoldPct = 100 * float64(tOutSold) / float64(tOutSeats)
		}
		res.HoldLatency = measure.Summarize(holdL)
		res.ConfirmLatency = measure.Summarize(confL)
		res.DeferredLatency = measure.Summarize(defL)
		res.CustomerLatency = measure.Summarize(custL)

		au, err := RunAudit(ctx, db, d, sl.world, fmt.Sprintf("race@%d#%d", tier, trial), ids)
		if err != nil {
			return nil, err
		}
		res.Audit = au
		status := "ok"
		if au.Violations() > 0 {
			status = "VIOLATIONS"
		}
		if res.UnderSold > 0 {
			status += fmt.Sprintf(", UNDERSOLD %d events (%d seats)", res.UnderSold, res.UnderSoldSeats)
		}
		if res.Leaked > 0 {
			status += fmt.Sprintf(", LEAKED %d seats", res.Leaked)
		}
		if res.TimedOut > 0 {
			status += fmt.Sprintf(", %d timed out (%.0f%% sold)", res.TimedOut, res.TimedOutSoldPct)
		}
		fmt.Printf("    race %6d seats x%-4d %9.1f seats/s  conflicts/hold=%.2f reads/hold=%.2f  hold p50=%.2fms p99=%.2fms  deferred %d/%d ok  confirm retries sold/run=%d/%d  errors=%d  %s\n",
			tier, res.Events, res.SeatsPerSec, res.ConflictsPerHold, res.MapReadsPerHold, res.HoldLatency.P50MS,
			res.HoldLatency.P99MS, res.DeferredConfirmed, res.DeferredHolds, res.ConfirmRetrySuccesses, res.ConfirmRetries, res.Errors, status)
		fmt.Printf("      audit: %s\n", au)
		if res.FirstError != "" {
			fmt.Printf("      first error: %s\n", res.FirstError)
		}
		out = append(out, res)
	}
	return out, nil
}

func raceEvent(ctx context.Context, sl *Seller, ev *Event, s Settings) *raceOne {
	one := &raceOne{}
	var (
		mu       sync.Mutex
		wg       sync.WaitGroup
		gate     = make(chan struct{})
		start    time.Time
		deadline time.Time
		lastHold atomic.Int64
		timedOut atomic.Bool
		errOnce  sync.Once
	)
	recordErr := func(err error) {
		one.errs.Add(1)
		errOnce.Do(func() { mu.Lock(); one.firstErr = err.Error(); mu.Unlock() })
	}
	for w := 0; w < s.RaceBuyers; w++ {
		wg.Add(1)
		go func(w int, seed int64) {
			defer wg.Done()
			r := rand.New(rand.NewSource(seed))
			buyer := &Buyer{seller: sl, node: w % 2, r: r, maxConflicts: s.MaxConflicts}
			<-gate
			var hl, cl, ul []time.Duration
			var def []deferredHold
			consecutiveErrors := 0
			for {
				if time.Now().After(deadline) {
					timedOut.Store(true)
					break
				}
				cust := sl.world.randomCustomer(r)
				n := partySize(r)
				t0 := time.Now()
				actx := WithAttemptDeadline(ctx, deadline)
				a, err := buyer.Acquire(actx, ev, cust, n, s.HoldTTL)
				one.customers.Add(1)
				one.attempts.Add(int64(a.Conflicts))
				one.conflicts.Add(int64(a.Conflicts))
				one.mapReads.Add(int64(a.MapReads))
				one.secReads.Add(int64(a.SectionReads))
				one.retries.Add(int64(a.EngineRetries))
				if errors.Is(err, errAttemptDeadline) {
					timedOut.Store(true)
					break
				}
				if err != nil {
					recordErr(err)
					if consecutiveErrors++; consecutiveErrors > 100 || ctx.Err() != nil {
						break
					}
					continue
				}
				consecutiveErrors = 0
				if a.GaveUp {
					one.gaveUp.Add(1)
					continue
				}
				if a.NoBlock {
					one.noBlock.Add(1)
					// A party that finds no block leaves; the queue stops only
					// when nothing at all is available.
					left, err := sl.Available(ctx, buyer.node, ev)
					if err != nil {
						recordErr(err)
						continue
					}
					if left == 0 || n == 1 {
						break
					}
					continue
				}
				one.attempts.Add(1)
				one.holds.Add(1)
				one.seats.Add(int64(n))
				hl = append(hl, a.HoldLatency)
				for at := int64(time.Since(start)); ; {
					prev := lastHold.Load()
					if at <= prev || lastHold.CompareAndSwap(prev, at) {
						break
					}
				}
				if r.Intn(100) < s.RaceDeferPct {
					def = append(def, deferredHold{node: buyer.node, block: a.Block, hold: a.Hold.HoldID, customer: cust})
					ul = append(ul, time.Since(t0))
					continue
				}
				t1 := time.Now()
				cr, err := sl.Confirm(ctx, buyer.node, a.Block, a.Hold.HoldID, cust)
				one.retries.Add(int64(cr.Retries))
				one.noteConfirm(cr)
				if err != nil {
					recordErr(err)
					continue
				}
				if !cr.Sold {
					one.immFailed.Add(1)
					continue
				}
				cl = append(cl, time.Since(t1))
				ul = append(ul, time.Since(t0))
			}
			mu.Lock()
			one.holdL = append(one.holdL, hl...)
			one.confL = append(one.confL, cl...)
			one.custL = append(one.custL, ul...)
			one.deferred = append(one.deferred, def...)
			mu.Unlock()
		}(w, ev.ID*7919+int64(w)*104729)
	}
	start = time.Now()
	deadline = start.Add(s.RaceTimeout)
	fmt.Printf("      race event %d: gate %s, deadline %s\n", ev.ID, start.UTC().Format(time.RFC3339Nano), deadline.UTC().Format(time.RFC3339Nano))
	close(gate)
	wg.Wait()
	one.allocate = time.Duration(lastHold.Load())
	if one.allocate == 0 {
		one.allocate = time.Since(start)
	}
	one.timedOut = timedOut.Load()

	// The deferred customers pay now, after the crowd -- concurrently, as they would.
	work := make(chan deferredHold)
	var dwg sync.WaitGroup
	for w := 0; w < s.RaceBuyers; w++ {
		dwg.Add(1)
		go func() {
			defer dwg.Done()
			var dl []time.Duration
			for dh := range work {
				t0 := time.Now()
				dctx := WithAttemptDeadline(ctx, time.Now().Add(10*time.Second+s.RaceTimeout))
				cr, err := sl.Confirm(dctx, dh.node, dh.block, dh.hold, dh.customer)
				one.noteConfirm(cr)
				if err != nil {
					recordErr(err)
					continue
				}
				if cr.Sold {
					one.deferredConfirmed.Add(1)
					dl = append(dl, time.Since(t0))
				} else {
					one.deferredRefused.Add(1)
				}
			}
			mu.Lock()
			one.defL = append(one.defL, dl...)
			mu.Unlock()
		}()
	}
	for _, dh := range one.deferred {
		work <- dh
	}
	close(work)
	dwg.Wait()
	one.wall = time.Since(start)

	if !one.timedOut {
		finishRaceEvent(ctx, sl, ev, s, one)
	}
	return one
}

// finishRaceEvent: the last customer sells what is left, one seat at a time; then
// unsold seats are under-sold, and any of them that cannot be held has leaked.
func finishRaceEvent(ctx context.Context, sl *Seller, ev *Event, s Settings, one *raceOne) {
	r := rand.New(rand.NewSource(ev.ID + 1))
	buyer := &Buyer{seller: sl, node: 0, r: r, maxConflicts: s.MaxConflicts}
	sweepCtx := WithAttemptDeadline(ctx, time.Now().Add(s.RaceTimeout))
	capacity := sl.world.ds.venue(ev.VenueID).Capacity
	for i := 0; i <= capacity; i++ {
		cust := sl.world.randomCustomer(r)
		a, err := buyer.Acquire(sweepCtx, ev, cust, 1, s.HoldTTL)
		if errors.Is(err, errAttemptDeadline) {
			one.timedOut = true
			return
		}
		if err != nil {
			one.errs.Add(1)
			continue
		}
		if a.NoBlock {
			break
		}
		if !a.Granted {
			continue
		}
		cr, err := sl.Confirm(ctx, 0, a.Block, a.Hold.HoldID, cust)
		one.noteConfirm(cr)
		if err != nil {
			one.errs.Add(1)
			continue
		}
		if cr.Sold {
			one.sweepSold++
		}
	}
	left, err := sl.Available(ctx, 0, ev)
	if err != nil {
		one.errs.Add(1)
		return
	}
	one.underSold = left
	if left > 0 {
		one.leakProbed, one.leaked = probeLeaks(ctx, sl, ev, s.LcProbeCapPerEvent, s.HoldTTL, nil)
	}
}

// probeLeaks tries to hold, one at a time, every seat the ledger shows unsold and
// without a valid hold (at most cap of them, chosen at random, plus every seat
// whose last hold has expired). A refusal means the design keeps a seat nobody
// can have. Granted probe holds are released at once. extra restricts the probe
// to those seats when non-nil.
func probeLeaks(ctx context.Context, sl *Seller, ev *Event, cap int, ttl time.Duration, extra []int32) (probed, leaked int64) {
	el := sl.world.ledger.Event(ev.ID)
	v := sl.world.ds.venue(ev.VenueID)
	var now time.Time
	if err := sl.db.QueryRow(ctx, "SELECT now()").Scan(&now); err != nil {
		return 0, 0
	}
	_, final := el.Evaluate()
	var expired, free []int32
	if extra != nil {
		expired = extra
	} else {
		for _, seat := range v.Seats {
			f, ok := final[seat.ID]
			switch {
			case !ok:
				free = append(free, seat.ID)
			case f.sold:
			case f.owner != 0 && !f.expires.After(now):
				expired = append(expired, seat.ID)
			}
		}
		r := rand.New(rand.NewSource(ev.ID + 17))
		r.Shuffle(len(free), func(i, j int) { free[i], free[j] = free[j], free[i] })
		if len(free) > cap {
			free = free[:cap]
		}
	}
	targets := append(expired, free...)
	sort.Slice(targets, func(i, j int) bool { return targets[i] < targets[j] })
	const probeCustomer = int64(-1)
	for _, seat := range targets {
		b := Block{Event: ev, Section: v.seat(seat).Section, Seats: []int32{seat}}
		hr, err := sl.Hold(ctx, 0, b, probeCustomer, ttl)
		if err != nil {
			continue
		}
		probed++
		if !hr.Granted {
			leaked++
			el.mu.Lock()
			el.example(Example{Kind: "leaked seat: unsold, no valid hold, and a hold was refused", EventID: ev.ID, SeatID: seat})
			el.mu.Unlock()
			continue
		}
		_, _ = sl.Release(ctx, 0, b, hr.HoldID)
	}
	return probed, leaked
}

func containsInt(xs []int, v int) bool {
	for _, x := range xs {
		if x == v {
			return true
		}
	}
	return false
}
