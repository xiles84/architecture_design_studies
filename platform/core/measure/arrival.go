package measure

import (
	"context"
	"errors"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

// This file adds the *open-loop* driver. Run() in loop.go is closed-loop: each
// worker waits for its own response before issuing the next request, so when the
// server stalls the load generator stalls with it and the stall is recorded once
// instead of being charged to every request that would have arrived during it --
// coordinated omission.
//
// A cadence experiment cannot use that. The question there is "the portal is
// offered N updates per second; what does it deliver?", and the offered rate is
// a property of the fleet, not of the server's speed. So arrivals are scheduled
// on a clock and the *offered* stream does not slow down when the database does.
// What the generator could not deliver is dropped and counted, never silently
// absorbed: a run whose client saturated must say so, because otherwise a
// generator limit is indistinguishable from a database limit
// (LESSONS_LEARNED, "Zero operation errors does not mean offered demand was
// served").

// Schedule describes an open-loop arrival process.
type Schedule struct {
	// Rate is the offered arrival rate in operations per second. It is what the
	// cadence calculation produces: installed products / update period. See
	// CadenceRate.
	//
	// A daily cadence on a small fleet is a handful of operations per hour, and
	// no run may wait for them. The protocol compresses such a rate and labels
	// the result a compressed-time validation, exactly as study 03 compressed a
	// forty-minute hold into a lifecycle it could measure; a compressed run is
	// never reported as the temporal result it stands in for.
	Rate float64
	// Duration is the offered window.
	Duration time.Duration
	// Count, when positive, caps how many operations are offered. It exists for
	// the low-cadence regimes: "once per day" on a small fleet is a handful of
	// operations, and waiting a day for them would measure the wall clock rather
	// than the portal.
	Count int64
	// Workers bounds how many operations may be in flight.
	Workers int
	// Queue is the arrival buffer. When it is full an arrival is dropped and
	// counted. Sizing it generously would hide saturation; sizing it at zero
	// would measure the scheduler.
	Queue int
	// Burst > 1 makes the schedule synchronized: Burst operations are offered at
	// the same instant, every Burst/Rate seconds. A fleet that publishes on the
	// minute is a burst, not a smooth trickle, and the two are different
	// workloads.
	Burst int
	// Jitter spreads each arrival uniformly within its slot. Without it a fixed
	// rate degenerates into a metronome, and a periodic workload can align with
	// the engine's own background work (a checkpoint, an autovacuum pass) in a
	// way that repeats identically every run.
	Jitter bool
	// Seed makes the jitter reproducible.
	Seed int64
}

// ArrivalResult separates what was offered from what was delivered. The gap is
// the finding, not a footnote.
type ArrivalResult struct {
	Name          string  `json:"name"`
	OfferedRate   float64 `json:"offered_rate_per_sec"`
	DeliveredRate float64 `json:"delivered_rate_per_sec"`
	Offered       int64   `json:"offered"`
	Started       int64   `json:"started"`
	Completed     int64   `json:"completed"`
	Rejected      int64   `json:"rejected"`
	Errors        int64   `json:"errors"`
	Retries       int64   `json:"retries"`
	// Dropped counts arrivals the generator could not hand to a worker before
	// the buffer filled. These are offered work that was never attempted.
	Dropped       int64 `json:"dropped"`
	QueueDepthMax int64 `json:"queue_depth_max"`
	QueueBound    int   `json:"queue_bound"`
	// SchedulingLag is how late each operation started relative to the instant it
	// was scheduled for. It is the client's queueing delay, and it is what makes
	// "the client saturated" a measurement rather than a guess.
	SchedulingLagMS LatencyStats `json:"scheduling_lag_ms"`
	Latency         LatencyStats `json:"latency"`
	DurationS       float64      `json:"duration_s"`
	// ClientSaturated is set when the generator could not deliver its own
	// schedule: arrivals were dropped, or the tail scheduling lag exceeded one
	// arrival slot. A saturated client invalidates a claim about the server.
	ClientSaturated  bool   `json:"client_saturated"`
	SaturationReason string `json:"saturation_reason,omitempty"`
	FirstError       string `json:"first_error,omitempty"`
}

// CadenceRate is the global offered update rate the study protocol defines:
//
//	installed products / update period
//
// It is a calculation. A "once per day" cadence on a 500-product fleet is
// 0.0058 updates per second, and the honest way to test it is a small fixed
// count landing on a pre-aged history -- not a day of wall-clock time. Any
// capacity figure derived from this rate is a *projection* and must be labelled
// as one wherever it is reported.
func CadenceRate(installedProducts int, period time.Duration) float64 {
	if period <= 0 || installedProducts <= 0 {
		return 0
	}
	return float64(installedProducts) / period.Seconds()
}

// RunOpenLoop offers operations on a schedule and measures what the database
// delivered. Unlike Run, the offered stream is not slowed by a stalled server.
func RunOpenLoop(ctx context.Context, s Schedule, name string, op Op) ArrivalResult {
	res := ArrivalResult{Name: name, OfferedRate: s.Rate}
	if s.Workers < 1 {
		s.Workers = 1
	}
	if s.Queue < 1 {
		s.Queue = s.Workers * 4
	}
	if s.Burst < 1 {
		s.Burst = 1
	}
	if s.Rate <= 0 {
		s.Rate = 1
	}
	if s.Duration <= 0 {
		if s.Count > 0 {
			// A count-bounded schedule needs a window that can actually hold its
			// count. Deriving one from the rate stops "offer these N operations"
			// from silently truncating at a default second and reporting a
			// shortfall that belongs to the caller, not the database.
			s.Duration = time.Duration(float64(s.Count)/s.Rate*float64(time.Second)) + time.Second
		} else {
			s.Duration = time.Second
		}
	}
	// One arrival instant every Burst/Rate seconds, so the offered rate is Rate
	// whether the fleet trickles or bursts.
	slot := time.Duration(float64(time.Second) * float64(s.Burst) / s.Rate)
	if slot <= 0 {
		slot = time.Nanosecond
	}
	res.QueueBound = s.Queue

	type job struct {
		scheduled time.Time
		r         *rand.Rand
	}

	jobs := make(chan job, s.Queue)
	var (
		offered, started, completed, rejected, errs, retries, dropped atomic.Int64
		queueMax                                                       atomic.Int64
		mu                                                             sync.Mutex
		lags, lat                                                      []time.Duration
		firstErr                                                       string
	)

	start := time.Now()
	deadline := start.Add(s.Duration)

	schedDone := make(chan struct{})
	go func() {
		defer close(schedDone)
		defer close(jobs)
		r := rand.New(rand.NewSource(s.Seed))
		for i := int64(0); ; i++ {
			if ctx.Err() != nil {
				return
			}
			if s.Count > 0 && offered.Load() >= s.Count {
				return
			}
			at := start.Add(time.Duration(float64(i) * float64(slot)))
			// The window is half-open: an instant exactly at the deadline is
			// not offered. Including it made a burst schedule offer one extra
			// burst -- 30 operations where 50/s over 400 ms is 20.
			if !at.Before(deadline) {
				return
			}
			if s.Jitter {
				at = at.Add(time.Duration(r.Int63n(int64(slot) + 1)))
			}
			if d := time.Until(at); d > 0 {
				t := time.NewTimer(d)
				select {
				case <-t.C:
				case <-ctx.Done():
					t.Stop()
					return
				}
			}
			for b := 0; b < s.Burst; b++ {
				if s.Count > 0 && offered.Load() >= s.Count {
					return
				}
				n := offered.Load()
				j := job{
					scheduled: at,
					// A private source per operation, derived from the schedule
					// seed, so keys are drawn fresh per execution and the run
					// replays exactly.
					r: rand.New(rand.NewSource(s.Seed + n*7919)),
				}
				offered.Add(1)
				select {
				case jobs <- j:
				default:
					// Offered, never attempted. Counting it here is what keeps a
					// generator limit from being read as a database limit.
					dropped.Add(1)
				}
				if q := int64(len(jobs)); q > queueMax.Load() {
					queueMax.Store(q)
				}
			}
		}
	}()

	var wg sync.WaitGroup
	for w := 0; w < s.Workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := range jobs {
				started.Add(1)
				lag := time.Since(j.scheduled)
				if lag < 0 {
					lag = 0
				}
				mu.Lock()
				lags = append(lags, lag)
				mu.Unlock()

				t0 := time.Now()
				o, err := op(ctx, j.r)
				el := time.Since(t0)
				retries.Add(int64(o.Retries))
				if errors.Is(err, ErrExhausted) {
					return
				}
				if err != nil {
					errs.Add(1)
					mu.Lock()
					if firstErr == "" {
						firstErr = err.Error()
					}
					mu.Unlock()
					continue
				}
				if o.Rejected {
					rejected.Add(1)
					continue
				}
				completed.Add(1)
				mu.Lock()
				lat = append(lat, el)
				mu.Unlock()
			}
		}()
	}

	<-schedDone
	wg.Wait()

	res.DurationS = time.Since(start).Seconds()
	res.Offered = offered.Load()
	res.Started = started.Load()
	res.Completed = completed.Load()
	res.Rejected = rejected.Load()
	res.Errors = errs.Load()
	res.Retries = retries.Load()
	res.Dropped = dropped.Load()
	res.QueueDepthMax = queueMax.Load()
	res.FirstError = firstErr
	if res.DurationS > 0 {
		res.DeliveredRate = float64(res.Completed) / res.DurationS
	}

	mu.Lock()
	res.SchedulingLagMS = Summarize(lags)
	res.Latency = Summarize(lat)
	mu.Unlock()

	// Saturation is decided by fixed rules, so it cannot be argued away later.
	switch {
	case res.Dropped > 0:
		res.ClientSaturated = true
		res.SaturationReason = "arrivals dropped: the buffer filled before a worker could take them"
	case res.SchedulingLagMS.P99MS > float64(slot.Nanoseconds())/1e6:
		res.ClientSaturated = true
		res.SaturationReason = "scheduling lag exceeded one arrival slot at p99"
	}
	return res
}
