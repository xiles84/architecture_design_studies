package measure

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

// Outcome is what one operation did, beyond succeeding or failing.
//
// Rejected exists because some operations have a legitimate "no": booking a
// seat in a sold-out event is a correct answer, not an error, and it must be
// counted and timed separately. Folding rejections into successes would let a
// design that says "sold out" quickly look like a design that sells quickly.
type Outcome struct {
	Rejected bool
	// Retries is how many times the operation had to start again because it
	// lost a race (serialization failure, failed compare-and-set, unique
	// violation). Throughput without retries hides how the work was achieved.
	Retries int
}

// Op is one unit of work. It receives its worker's private random source so
// that keys are drawn fresh per execution without lock contention on a shared
// generator.
type Op func(ctx context.Context, r *rand.Rand) (Outcome, error)

// ErrExhausted means a destructive operation has used up the finite set of rows
// it may consume. It is a stop signal, never an error to count and loop past:
// in study 01, looping past it turned "delete" into 111 million instant
// failures and a reported 0 ops/s.
var ErrExhausted = errors.New("finite pool exhausted")

// Budget says how long, or how much, to measure.
type Budget struct {
	Workers int
	// Warmup is discarded. Either a duration or an op count may be given; an
	// op count is used for operations too expensive to run for a fixed time
	// (publishing a 100 000-seat event).
	Warmup    time.Duration
	WarmupOps int64
	// Each trial runs for Duration, or until OpsPerTrial operations have been
	// attempted, whichever comes first. A zero Duration with a positive
	// OpsPerTrial means count-based only -- the rule for destructive operations,
	// which drain their pool if measured by time.
	Duration    time.Duration
	OpsPerTrial int64
	Trials      int
	// SafetyLimit bounds a count-based trial that would otherwise never end.
	SafetyLimit time.Duration
}

// Result is the measurement of one operation.
type Result struct {
	Name       string  `json:"name"`
	Ops        int64   `json:"ops"`
	Rejected   int64   `json:"rejected,omitempty"`
	Retries    int64   `json:"retries,omitempty"`
	Errors     int64   `json:"errors"`
	FirstError string  `json:"first_error,omitempty"`
	DurationS  float64 `json:"duration_s"`
	// OpsPerSec is the MEDIAN successful ops/s across trials.
	OpsPerSec float64      `json:"ops_per_sec"`
	Latency   LatencyStats `json:"latency"`
	// RejectLatency is the distribution of the "no" answers, kept apart.
	RejectLatency *LatencyStats `json:"reject_latency,omitempty"`
	Trials        []float64     `json:"trials_ops_per_sec,omitempty"`
	SpreadPct     float64       `json:"spread_pct,omitempty"`
	OpBudget      int64         `json:"op_budget,omitempty"`
	Exhausted     bool          `json:"pool_exhausted,omitempty"`
	SafetyStopped bool          `json:"safety_stopped,omitempty"`
}

// FiniteBudget sizes a destructive operation so warmup and every trial fit
// inside HALF of the pool it consumes. Half, because the last rows of a table
// are not representative of a populated one.
func FiniteBudget(pool int64, trials int) int64 {
	if trials < 1 {
		trials = 1
	}
	b := pool / int64(2*(trials+1))
	return max(b, 1)
}

// Run drives b.Workers workers against op: an unrecorded warmup, then b.Trials
// recorded trials. The loop is closed -- each worker waits for its own response
// -- so deep tails are a floor, not an SLO figure (coordinated omission).
func Run(ctx context.Context, b Budget, name string, op Op) Result {
	res := Result{Name: name, OpBudget: b.OpsPerTrial}
	var exhausted atomic.Bool
	if b.Workers < 1 {
		b.Workers = 1
	}

	type phaseOut struct {
		ok, rej       []time.Duration
		ops, rejected int64
		retries, errs int64
		firstErr      string
		safety        bool
	}

	phase := func(d time.Duration, opLimit int64, record bool) phaseOut {
		if d <= 0 && opLimit <= 0 {
			return phaseOut{}
		}
		var deadline time.Time
		switch {
		case d > 0:
			deadline = time.Now().Add(d)
		case b.SafetyLimit > 0:
			deadline = time.Now().Add(b.SafetyLimit)
		default:
			deadline = time.Now().Add(30 * time.Minute)
		}
		var (
			wg                           sync.WaitGroup
			mu                           sync.Mutex
			out                          phaseOut
			claimed, ops, rejected, errs atomic.Int64
			retries                      atomic.Int64
			errOnce                      sync.Once
		)
		for w := 0; w < b.Workers; w++ {
			wg.Add(1)
			go func(seed int64) {
				defer wg.Done()
				r := rand.New(rand.NewSource(seed))
				var okL, rejL []time.Duration
				for time.Now().Before(deadline) {
					// Claim before executing, so a count budget is exact
					// even with many workers racing for the last op.
					if opLimit > 0 && claimed.Add(1) > opLimit {
						break
					}
					if ctx.Err() != nil {
						break
					}
					t0 := time.Now()
					o, err := op(ctx, r)
					el := time.Since(t0)
					retries.Add(int64(o.Retries))
					if errors.Is(err, ErrExhausted) {
						exhausted.Store(true)
						break
					}
					if err != nil {
						errs.Add(1)
						errOnce.Do(func() { mu.Lock(); out.firstErr = err.Error(); mu.Unlock() })
						continue
					}
					if o.Rejected {
						rejected.Add(1)
						if record {
							rejL = append(rejL, el)
						}
						continue
					}
					ops.Add(1)
					if record {
						okL = append(okL, el)
					}
				}
				if record {
					mu.Lock()
					out.ok = append(out.ok, okL...)
					out.rej = append(out.rej, rejL...)
					mu.Unlock()
				}
			}(time.Now().UnixNano() + int64(w)*7919)
		}
		wg.Wait()
		out.ops, out.rejected, out.errs, out.retries = ops.Load(), rejected.Load(), errs.Load(), retries.Load()
		out.safety = d <= 0 && opLimit > 0 && claimed.Load() < opLimit && !exhausted.Load() && time.Now().After(deadline)
		return out
	}

	phase(b.Warmup, b.WarmupOps, false)

	trials := max(b.Trials, 1)
	var okAll, rejAll []time.Duration
	for t := 0; t < trials; t++ {
		start := time.Now()
		p := phase(b.Duration, b.OpsPerTrial, true)
		el := time.Since(start).Seconds()
		res.Ops += p.ops
		res.Rejected += p.rejected
		res.Errors += p.errs
		res.Retries += p.retries
		res.DurationS += el
		res.SafetyStopped = res.SafetyStopped || p.safety
		if res.FirstError == "" {
			res.FirstError = p.firstErr
		}
		if el > 0 {
			res.Trials = append(res.Trials, float64(p.ops)/el)
		}
		okAll = append(okAll, p.ok...)
		rejAll = append(rejAll, p.rej...)
		if exhausted.Load() {
			break
		}
	}
	res.OpsPerSec = Median(res.Trials)
	res.SpreadPct = SpreadPct(res.Trials)
	res.Exhausted = exhausted.Load()
	res.Latency = Summarize(okAll)
	if len(rejAll) > 0 {
		rl := Summarize(rejAll)
		res.RejectLatency = &rl
	}
	return res
}

// SpreadNote flags, at the point of measurement, a cell whose own trials
// disagree by more than a fifth.
func SpreadNote(r Result) string {
	if len(r.Trials) < 2 {
		return ""
	}
	flag := ""
	if r.SpreadPct > 20 {
		flag = "  <-- noisy"
	}
	return fmt.Sprintf("  [%d trials, spread %.0f%%]%s", len(r.Trials), r.SpreadPct, flag)
}
