package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"adsplatform/core/catalog"
	"adsplatform/core/measure"
	"adsplatform/ports"
)

// ---------------------------------------------------------------------------
// Experiment: holds with expiry (H0, H1 only).
//
// Buyers take a hold, spend time in the basket, sometimes abandon it, and
// otherwise pay and confirm. A sweeper releases expired holds. Timings are set
// so that a minority of payments finish after their hold has expired -- the
// real-world case of a slow card authorisation -- and a released seat is
// immediately wanted by buyers who were told "sold out" moments before.
//
// Timing distribution (defaults): basket time U[0, 200 ms], payment U[20, 100 ms],
// hold TTL 250 ms. About 8% of confirmations arrive after expiry; the sweeper
// runs every 50 ms, so some of those find their seat already resold.
//
// An event ends when its tickets reach capacity (watched by a monitor) or at the
// timeout. Buyers already mid-payment are allowed to finish -- which is exactly
// when H0's late confirmations land.
// ---------------------------------------------------------------------------

type HoldResult struct {
	Tier   int `json:"tier"`
	Events int `json:"events"`
	Buyers int `json:"buyers"`

	Seats           int64   `json:"seats"`
	Sold            int64   `json:"sold"`
	Holds           int64   `json:"holds"`
	HoldsRejected   int64   `json:"holds_rejected"`
	Abandoned       int64   `json:"abandoned"`
	ExpiredAtCheck  int64   `json:"expired_at_precheck"`
	Confirmed       int64   `json:"confirmed"`
	LateRejected    int64   `json:"late_confirms_rejected"`
	SweptReleased   int64   `json:"swept_released"`
	Retries         int64   `json:"retries"`
	Errors          int64   `json:"errors"`
	FirstError      string  `json:"first_error,omitempty"`
	WallS           float64 `json:"wall_s"`
	ConfirmedPerSec float64 `json:"confirmed_per_sec"`

	HoldLatency    measure.LatencyStats `json:"hold_latency"`
	ConfirmLatency measure.LatencyStats `json:"confirm_latency"`

	TimedOut        int              `json:"timed_out_events"`
	Overbooked      int              `json:"overbooked_events"`
	OverbookedSeats int64            `json:"overbooked_seats"`
	Examples        []EventViolation `json:"examples,omitempty"`
	Audit           *Audit           `json:"audit,omitempty"`
}

func RunHolds(ctx context.Context, db ports.DB, bk *Booker, d Design, ds *Dataset, world *World, s Settings) ([]HoldResult, error) {
	auditStmts, err := mustStmts(d.ID, "audit.sql")
	if err != nil {
		return nil, err
	}
	soldSQL := catalog.Map(auditStmts)["a_event_sold"].SQL

	var out []HoldResult
	for _, tier := range ds.tiersOf("race") {
		if len(s.HoldTiers) > 0 && !containsInt(s.HoldTiers, tier) {
			continue
		}
		res := HoldResult{Tier: tier, Buyers: s.HoldBuyers}
		var holdL, confL []time.Duration
		var wall time.Duration
		evs := ds.eventsOf("race", tier)
		if s.HoldEventsPerTier > 0 && len(evs) > s.HoldEventsPerTier {
			evs = evs[:s.HoldEventsPerTier]
		}
		for _, e := range evs {
			ev := world.event(e.ID)
			one, err := holdEvent(ctx, db, bk, world, ev, soldSQL, s)
			if err != nil {
				return nil, err
			}
			sold, err := eventSold(ctx, db, soldSQL, ev.ID)
			if err != nil {
				return nil, err
			}
			res.Events++
			res.Seats += int64(ev.Capacity)
			res.Sold += sold
			res.Holds += one.holds
			res.HoldsRejected += one.holdsRejected
			res.Abandoned += one.abandoned
			res.ExpiredAtCheck += one.expiredAtCheck
			res.Confirmed += one.confirmed
			res.LateRejected += one.lateRejected
			res.SweptReleased += one.swept
			res.Retries += one.retries
			res.Errors += one.errors
			if res.FirstError == "" {
				res.FirstError = one.firstErr
			}
			if one.timedOut {
				res.TimedOut++
			}
			wall += one.wall
			holdL = append(holdL, one.holdL...)
			confL = append(confL, one.confL...)
			if sold > int64(ev.Capacity) {
				res.Overbooked++
				res.OverbookedSeats += sold - int64(ev.Capacity)
				if len(res.Examples) < auditExamples {
					res.Examples = append(res.Examples, EventViolation{EventID: ev.ID, Capacity: int64(ev.Capacity), Actual: sold, Note: "overbooked"})
				}
			}
		}
		res.WallS = wall.Seconds()
		if res.WallS > 0 {
			res.ConfirmedPerSec = float64(res.Confirmed) / res.WallS
		}
		res.HoldLatency = measure.Summarize(holdL)
		res.ConfirmLatency = measure.Summarize(confL)

		// Release what is left so the counter audit compares settled state:
		// every hold still 'held' after the sweeper stopped has expired by now.
		time.Sleep(s.HoldTTL + 50*time.Millisecond)
		for {
			n, err := bk.Sweep(ctx, 1000)
			if err != nil {
				return nil, fmt.Errorf("final sweep: %w", err)
			}
			res.SweptReleased += int64(n)
			if n == 0 {
				break
			}
		}
		au, err := RunAudit(ctx, db, d, world, fmt.Sprintf("holds@%d", tier))
		if err != nil {
			return nil, err
		}
		res.Audit = au

		status := "ok"
		if res.Overbooked > 0 {
			status = fmt.Sprintf("OVERBOOKED %d events (+%d seats)", res.Overbooked, res.OverbookedSeats)
		}
		fmt.Printf("    holds  %6d seats x%-3d %8.1f confirmed/s  holds=%d abandoned=%d late-rejected=%d swept=%d errors=%d  %s\n",
			tier, res.Events, res.ConfirmedPerSec, res.Holds, res.Abandoned, res.LateRejected, res.SweptReleased, res.Errors, status)
		fmt.Printf("      audit: %s\n", au)
		if res.FirstError != "" {
			fmt.Printf("      first error: %s\n", res.FirstError)
		}
		out = append(out, res)
	}
	return out, nil
}

type holdOne struct {
	holds, holdsRejected, abandoned, expiredAtCheck int64
	confirmed, lateRejected, swept, retries, errors int64
	firstErr                                        string
	timedOut                                        bool
	wall                                            time.Duration
	holdL, confL                                    []time.Duration
}

func holdEvent(ctx context.Context, db ports.DB, bk *Booker, world *World, ev *EventRef, soldSQL string, s Settings) (holdOne, error) {
	var (
		one                                   holdOne
		mu                                    sync.Mutex
		wg                                    sync.WaitGroup
		stop                                  atomic.Bool
		holds, holdsRej, abandoned, expired   atomic.Int64
		confirmed, late, swept, retries, errs atomic.Int64
		errOnce                               sync.Once
	)
	recordErr := func(err error) {
		errs.Add(1)
		errOnce.Do(func() { mu.Lock(); one.firstErr = err.Error(); mu.Unlock() })
	}
	uniform := func(r *rand.Rand, lo, hi time.Duration) time.Duration {
		if hi <= lo {
			return lo
		}
		return lo + time.Duration(r.Int63n(int64(hi-lo)))
	}
	ttlMS := int(s.HoldTTL / time.Millisecond)
	start := time.Now()
	deadline := start.Add(s.RaceTimeout)
	// Retries inside a hold or confirm stop shortly after the event's deadline,
	// so one pathological transaction cannot keep the experiment alive.
	ctx = WithAttemptDeadline(ctx, deadline.Add(10*time.Second))

	for w := 0; w < s.HoldBuyers; w++ {
		wg.Add(1)
		go func(seed int64) {
			defer wg.Done()
			r := rand.New(rand.NewSource(seed))
			var hl, cl []time.Duration
			for !stop.Load() {
				cust := world.randomCustomer(r)
				t0 := time.Now()
				resID, soldOut, rt, err := bk.Hold(ctx, ev, cust, ttlMS)
				retries.Add(int64(rt))
				if err != nil {
					recordErr(err)
					continue
				}
				if soldOut {
					// Sold out *for now*: holds expire, so a buyer who waits a moment
					// may still get in. That is exactly who takes a released seat.
					holdsRej.Add(1)
					time.Sleep(10 * time.Millisecond)
					continue
				}
				hl = append(hl, time.Since(t0))
				holds.Add(1)
				time.Sleep(uniform(r, 0, s.HoldThinkMax))
				if r.Intn(100) < s.HoldAbandonPct {
					abandoned.Add(1) // walks away; the sweeper will release it
					continue
				}
				valid, err := bk.CheckHold(ctx, resID)
				if err != nil {
					recordErr(err)
					continue
				}
				if !valid {
					expired.Add(1)
					continue
				}
				time.Sleep(uniform(r, s.HoldPayMin, s.HoldPayMax)) // card authorisation
				t1 := time.Now()
				_, lost, rt, err := bk.Confirm(ctx, ev, resID, cust)
				retries.Add(int64(rt))
				if err != nil {
					recordErr(err)
					continue
				}
				cl = append(cl, time.Since(t1))
				if lost {
					late.Add(1) // H1: the hold had expired; refund, no ticket
				} else {
					confirmed.Add(1)
				}
			}
			mu.Lock()
			one.holdL = append(one.holdL, hl...)
			one.confL = append(one.confL, cl...)
			mu.Unlock()
		}(time.Now().UnixNano() + int64(w)*15485863)
	}

	sweepDone := make(chan struct{})
	var sweepStop atomic.Bool
	go func() {
		defer close(sweepDone)
		for !sweepStop.Load() {
			n, err := bk.Sweep(ctx, 200)
			if err != nil {
				recordErr(fmt.Errorf("sweeper: %w", err))
			}
			swept.Add(int64(n))
			time.Sleep(s.HoldSweepEvery)
		}
	}()

	// Monitor: the event is done when its tickets reach capacity.
	for {
		time.Sleep(50 * time.Millisecond)
		sold, err := eventSold(ctx, db, soldSQL, ev.ID)
		if err != nil {
			stop.Store(true)
			wg.Wait()
			sweepStop.Store(true)
			<-sweepDone
			return one, fmt.Errorf("hold monitor event %d: %w", ev.ID, err)
		}
		if sold >= int64(ev.Capacity) {
			break
		}
		if time.Now().After(deadline) {
			one.timedOut = true
			break
		}
	}
	stop.Store(true)
	wg.Wait() // buyers mid-payment finish here
	one.wall = time.Since(start)
	sweepStop.Store(true)
	<-sweepDone

	one.holds, one.holdsRejected, one.abandoned = holds.Load(), holdsRej.Load(), abandoned.Load()
	one.expiredAtCheck, one.confirmed, one.lateRejected = expired.Load(), confirmed.Load(), late.Load()
	one.swept, one.retries, one.errors = swept.Load(), retries.Load(), errs.Load()
	return one, nil
}
