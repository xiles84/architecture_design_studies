package main

import (
	"context"
	"errors"
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
// The sell-out race.
//
// A crowd of buyers arrives at an unsold event at the same instant and keeps
// buying until the event answers "sold out". This is the hot drop, and it is
// where overbooking strategies actually differ: with demand spread across the
// catalogue, almost any strategy is fine.
//
// Each buyer goroutine is a queue of customers: it books one ticket, then the
// next customer in its queue tries, until a customer is told "sold out" and the
// queue goes home. Buyers start behind a gate so they really do arrive together.
//
// Two correctness outcomes are measured beside the speed:
//
//	overbooked   the event ended with more tickets than seats
//	underbooked  every buyer went home told "sold out", yet seats are unsold --
//	             the failure mode of an over-cautious design, invisible to any
//	             overbooking check
//
// Churn mode adds refunds: a fraction of successful buyers cancel straight away.
// A refund after the crowd has left would legitimately stay unsold, so after
// the crowd is done one last buyer sweeps the event until told "sold out"; only
// seats still unsold after THAT are underbooked.
//
// An organiser edits the event page throughout, so designs that put booking
// contention on the event row can be seen making unrelated writes wait.
// ---------------------------------------------------------------------------

type RaceResult struct {
	Mode   string `json:"mode"` // race | churn
	Trial  int    `json:"trial"`
	Tier   int    `json:"tier"`
	Events int    `json:"events"`
	// EventsPlanned is how many events of this tier existed; Events < EventsPlanned
	// means the tier's time budget ran out and the rest were not raced.
	EventsPlanned int `json:"events_planned"`
	Buyers        int `json:"buyers"`

	Seats      int64  `json:"seats"`
	Sold       int64  `json:"sold"`
	Successes  int64  `json:"successes"`
	Rejections int64  `json:"rejections"`
	Cancels    int64  `json:"cancels,omitempty"`
	SweepSold  int64  `json:"sweep_sold,omitempty"`
	Retries    int64  `json:"retries"`
	Errors     int64  `json:"errors"`
	FirstError string `json:"first_error,omitempty"`

	// SellOutS sums, over the tier's events, the time from the gate opening to
	// the last successful sale. SoldPerSec is measured over it. WallS adds the
	// time spent telling the remaining crowd "sold out"; for a 10-seat event with
	// 32 buyers that tail is most of the wall time, and dividing sales by it
	// would measure the rejection path, not selling.
	SellOutS   float64 `json:"sell_out_s"`
	WallS      float64 `json:"wall_s"`
	SoldPerSec float64 `json:"sold_per_sec"`
	// AttemptsPerSuccess = (successes + retries) / successes.
	AttemptsPerSuccess float64              `json:"attempts_per_success"`
	Latency            measure.LatencyStats `json:"latency"`
	RejectLatency      measure.LatencyStats `json:"reject_latency"`
	EditorLatency      measure.LatencyStats `json:"editor_latency"`
	EditorErrors       int64                `json:"editor_errors,omitempty"`

	TimedOut         int              `json:"timed_out_events"`
	TimedOutSoldPct  float64          `json:"timed_out_sold_pct,omitempty"`
	Overbooked       int              `json:"overbooked_events"`
	OverbookedSeats  int64            `json:"overbooked_seats"`
	Underbooked      int              `json:"underbooked_events"`
	UnderbookedSeats int64            `json:"underbooked_seats"`
	Examples         []EventViolation `json:"examples,omitempty"`
	Audit            *Audit           `json:"audit,omitempty"`
	LedgerAudit      *LedgerAudit     `json:"ledger_audit,omitempty"`
	// ReportChecks: r01/r05 against the harness's own counters, set once per
	// churn trial (on the last tier's result) after every tier has run on
	// that trial's world -- never for plain "race" (EH-02 AM-03.5).
	ReportChecks []Check `json:"report_checks,omitempty"`

	// Filled only when a report combines several trials: per-trial sold/s and
	// their spread. Never written by the harness.
	TrialsSoldPerSec []float64 `json:"-"`
	SpreadPct        float64   `json:"-"`
}

type raceOne struct {
	ok, rej, edit                      []time.Duration
	successes, rejections, cancels     int64
	retries, errors, editErrors, sweep int64
	firstErr                           string
	timedOut                           bool
	wall, sellOut                      time.Duration
}

func RunRaces(ctx context.Context, db ports.DB, bk *Booker, d Design, ds *Dataset, world *World, mode string, s Settings) ([]RaceResult, error) {
	auditStmts, err := mustStmts(d.ID, "audit.sql")
	if err != nil {
		return nil, err
	}
	soldSQL := catalog.Map(auditStmts)["a_event_sold"].SQL
	churn := 0
	kind := "race"
	if mode == "churn" {
		churn, kind = s.ChurnPct, "churn"
	}

	var out []RaceResult
	for _, tier := range ds.tiersOf(kind) {
		if len(s.RaceTiers) > 0 && !containsInt(s.RaceTiers, tier) {
			continue
		}
		res := RaceResult{Mode: mode, Tier: tier, Buyers: s.RaceBuyers}
		var ok, rej, edit []time.Duration
		var wall, sellOut time.Duration
		var timedOutSold, timedOutSeats int64

		tierStart := time.Now()
		planned := ds.eventsOf(kind, tier)
		res.EventsPlanned = len(planned)
		for _, e := range planned {
			// A tier of 100 small events on a slow configuration can otherwise
			// run for an hour. Events not started within the budget are skipped
			// and reported as such; the ones that ran are complete.
			if s.RaceTierBudget > 0 && time.Since(tierStart) > s.RaceTierBudget {
				break
			}
			ev := world.event(e.ID)
			one := raceEvent(ctx, db, bk, world, ev, s, churn)
			sold, err := eventSold(ctx, db, soldSQL, ev.ID)
			if err != nil {
				return nil, fmt.Errorf("race %s event %d: final count: %w", mode, ev.ID, err)
			}
			res.Events++
			res.Seats += int64(ev.Capacity)
			res.Sold += sold
			res.Successes += one.successes
			res.Rejections += one.rejections
			res.Cancels += one.cancels
			res.SweepSold += one.sweep
			res.Retries += one.retries
			res.Errors += one.errors
			res.EditorErrors += one.editErrors
			if res.FirstError == "" {
				res.FirstError = one.firstErr
			}
			wall += one.wall
			sellOut += one.sellOut
			ok = append(ok, one.ok...)
			rej = append(rej, one.rej...)
			edit = append(edit, one.edit...)

			switch {
			case sold > int64(ev.Capacity):
				res.Overbooked++
				res.OverbookedSeats += sold - int64(ev.Capacity)
				if len(res.Examples) < auditExamples {
					res.Examples = append(res.Examples, EventViolation{EventID: ev.ID, Capacity: int64(ev.Capacity), Actual: sold, Note: "overbooked"})
				}
			case one.timedOut:
				res.TimedOut++
				timedOutSold += sold
				timedOutSeats += int64(ev.Capacity)
			case sold < int64(ev.Capacity):
				res.Underbooked++
				res.UnderbookedSeats += int64(ev.Capacity) - sold
				if len(res.Examples) < auditExamples {
					res.Examples = append(res.Examples, EventViolation{EventID: ev.ID, Capacity: int64(ev.Capacity), Actual: sold,
						Note: "every buyer was told sold out with seats unsold"})
				}
			}
		}
		res.WallS = wall.Seconds()
		res.SellOutS = sellOut.Seconds()
		if res.SellOutS > 0 {
			res.SoldPerSec = float64(res.Successes) / res.SellOutS
		}
		if res.Successes > 0 {
			res.AttemptsPerSuccess = float64(res.Successes+res.Retries) / float64(res.Successes)
		}
		if timedOutSeats > 0 {
			res.TimedOutSoldPct = 100 * float64(timedOutSold) / float64(timedOutSeats)
		}
		res.Latency = measure.Summarize(ok)
		res.RejectLatency = measure.Summarize(rej)
		res.EditorLatency = measure.Summarize(edit)

		au, err := RunAudit(ctx, db, d, world, fmt.Sprintf("%s@%d", mode, tier))
		if err != nil {
			return nil, err
		}
		res.Audit = au

		// EH-02 AM-03.2: the ledger reconciliation audit runs after every
		// writing phase, and the sell-out race is X1's headline experiment --
		// the one place its extra insert is meant to be felt. A no-op for
		// every design without a ledger (RunLedgerAudit checks d.Ledger).
		la, err := RunLedgerAudit(ctx, db, d, fmt.Sprintf("%s@%d", mode, tier))
		if err != nil {
			return nil, err
		}
		res.LedgerAudit = la

		status := "ok"
		if res.Overbooked > 0 {
			status = fmt.Sprintf("OVERBOOKED %d events (+%d seats)", res.Overbooked, res.OverbookedSeats)
		} else if res.Underbooked > 0 {
			status = fmt.Sprintf("UNDERBOOKED %d events (%d seats)", res.Underbooked, res.UnderbookedSeats)
		}
		if res.TimedOut > 0 {
			status += fmt.Sprintf(", %d timed out (%.0f%% sold)", res.TimedOut, res.TimedOutSoldPct)
		}
		fmt.Printf("    %-6s %6d seats x%-4d %9.1f sold/s  p50=%7.3fms p99=%8.3fms  attempts/sale=%.2f errors=%d  editor p99=%.2fms  %s\n",
			mode, tier, res.Events, res.SoldPerSec, res.Latency.P50MS, res.Latency.P99MS, res.AttemptsPerSuccess,
			res.Errors, res.EditorLatency.P99MS, status)
		fmt.Printf("      audit: %s\n", au)
		if la != nil {
			fmt.Printf("      ledger audit: %s\n", la)
		}
		if res.FirstError != "" {
			fmt.Printf("      first error: %s\n", res.FirstError)
		}
		out = append(out, res)
	}
	return out, nil
}

func raceEvent(ctx context.Context, db ports.DB, bk *Booker, world *World, ev *EventRef, s Settings, churnPct int) raceOne {
	var (
		one                          raceOne
		mu                           sync.Mutex
		wg                           sync.WaitGroup
		gate                         = make(chan struct{})
		deadline                     time.Time
		stopped                      atomic.Bool
		timedOut                     atomic.Bool
		lastSale                     atomic.Int64 // nanoseconds after the gate opened
		start                        time.Time
		succ, rejc, canc, retr, errs atomic.Int64
		errOnce                      sync.Once
	)
	for w := 0; w < s.RaceBuyers; w++ {
		wg.Add(1)
		go func(seed int64) {
			defer wg.Done()
			r := rand.New(rand.NewSource(seed))
			<-gate
			var okL, rejL []time.Duration
			consecutiveErrors := 0
			for {
				// The deadline stops NEW attempts only. In-flight bookings finish
				// with the parent context, so a timeout never produces a commit of
				// unknown outcome that the reconciliation would have to excuse.
				if time.Now().After(deadline) {
					timedOut.Store(true)
					break
				}
				t0 := time.Now()
				br, err := bk.Book(WithAttemptDeadline(ctx, deadline), ev, world.randomCustomer(r), r)
				el := time.Since(t0)
				retr.Add(int64(br.Retries))
				if errors.Is(err, errAttemptDeadline) {
					timedOut.Store(true)
					break
				}
				if err != nil {
					errs.Add(1)
					errOnce.Do(func() { mu.Lock(); one.firstErr = err.Error(); mu.Unlock() })
					if consecutiveErrors++; consecutiveErrors > 100 || ctx.Err() != nil {
						break
					}
					continue
				}
				consecutiveErrors = 0
				if br.SoldOut {
					rejc.Add(1)
					rejL = append(rejL, el)
					break
				}
				succ.Add(1)
				okL = append(okL, el)
				for at := int64(time.Since(start)); ; {
					prev := lastSale.Load()
					if at <= prev || lastSale.CompareAndSwap(prev, at) {
						break
					}
				}
				if churnPct > 0 && r.Intn(100) < churnPct {
					if okc, rt, err := bk.Cancel(ctx, br.TicketID); err == nil && okc {
						canc.Add(1)
						retr.Add(int64(rt))
					} else if err != nil {
						errs.Add(1)
					}
				}
			}
			mu.Lock()
			one.ok = append(one.ok, okL...)
			one.rej = append(one.rej, rejL...)
			mu.Unlock()
		}(time.Now().UnixNano() + int64(w)*104729)
	}

	// The organiser, editing the page at a fixed interval while the crowd buys.
	editDone := make(chan struct{})
	go func() {
		defer close(editDone)
		r := rand.New(rand.NewSource(ev.ID))
		<-gate
		if s.EditorInterval <= 0 {
			return
		}
		tick := time.NewTicker(s.EditorInterval)
		defer tick.Stop()
		for !stopped.Load() {
			t0 := time.Now()
			err := bk.Edit(ctx, ev, r)
			el := time.Since(t0)
			mu.Lock()
			if err != nil {
				one.editErrors++
			} else {
				one.edit = append(one.edit, el)
			}
			mu.Unlock()
			<-tick.C
		}
	}()

	start = time.Now()
	deadline = start.Add(s.RaceTimeout)
	close(gate)
	wg.Wait()
	one.wall = time.Since(start)
	one.sellOut = time.Duration(lastSale.Load())
	if one.sellOut == 0 {
		one.sellOut = one.wall
	}
	stopped.Store(true)
	<-editDone

	one.successes, one.rejections, one.cancels = succ.Load(), rejc.Load(), canc.Load()
	one.retries, one.errors = retr.Load(), errs.Load()
	one.timedOut = timedOut.Load()

	if churnPct > 0 && !one.timedOut {
		// One last customer buys until told sold out, reclaiming seats refunded
		// after the crowd went home.
		r := rand.New(rand.NewSource(ev.ID + 1))
		sweepCtx := WithAttemptDeadline(ctx, time.Now().Add(s.RaceTimeout))
		for i := 0; i < ev.Capacity+1; i++ {
			br, err := bk.Book(sweepCtx, ev, world.randomCustomer(r), r)
			if errors.Is(err, errAttemptDeadline) {
				one.timedOut = true
				break
			}
			if err != nil {
				one.errors++
				continue
			}
			if br.SoldOut {
				break
			}
			one.sweep++
		}
	}
	return one
}

func containsInt(xs []int, v int) bool {
	for _, x := range xs {
		if x == v {
			return true
		}
	}
	return false
}
