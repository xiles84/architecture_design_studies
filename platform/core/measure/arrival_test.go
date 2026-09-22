package measure

import (
	"context"
	"math"
	"math/rand"
	"testing"
	"time"
)

// CadenceRate is the study's definition of cadence: installed products divided
// by the update period. The four regimes in the handoff must fall out of it,
// because the whole "do not wait a day" argument rests on this number.
func TestCadenceRateMatchesTheHandoffTable(t *testing.T) {
	const fleet = 500
	cases := []struct {
		name   string
		period time.Duration
		want   float64
	}{
		{"daily", 24 * time.Hour, 0.005787},
		{"hourly", time.Hour, 0.138889},
		{"minutely", time.Minute, 8.333333},
		{"secondly", time.Second, 500},
	}
	for _, c := range cases {
		got := CadenceRate(fleet, c.period)
		if math.Abs(got-c.want) > 0.0001 {
			t.Errorf("%s: CadenceRate(%d, %s) = %.6f, want %.6f", c.name, fleet, c.period, got, c.want)
		}
	}
}

func TestCadenceRateRejectsNonsense(t *testing.T) {
	if got := CadenceRate(0, time.Second); got != 0 {
		t.Errorf("no installed products should offer nothing, got %v", got)
	}
	if got := CadenceRate(500, 0); got != 0 {
		t.Errorf("a zero period has no rate, got %v", got)
	}
}

// With an instantaneous operation the generator must deliver essentially
// everything it offered. If this fails, every cadence number is measuring the
// harness rather than the portal.
func TestOpenLoopDeliversWhatItOffers(t *testing.T) {
	res := RunOpenLoop(context.Background(), Schedule{
		Rate: 200, Duration: 500 * time.Millisecond, Workers: 4, Queue: 64, Seed: 1,
	}, "instant", func(ctx context.Context, r *rand.Rand) (Outcome, error) {
		return Outcome{}, nil
	})
	if res.Dropped != 0 {
		t.Errorf("an instant op dropped %d arrivals; the harness is the bottleneck", res.Dropped)
	}
	if res.ClientSaturated {
		t.Errorf("an instant op reported a saturated client: %s", res.SaturationReason)
	}
	if res.Completed != res.Offered {
		t.Errorf("completed %d of %d offered", res.Completed, res.Offered)
	}
	// ~200/s over 0.5s, allowing for the first slot landing at t=0.
	if res.Offered < 80 || res.Offered > 105 {
		t.Errorf("offered %d operations at 200/s over 500ms; expected roughly 100", res.Offered)
	}
}

// The whole point of the open-loop driver: when the server is slower than the
// offered rate, the generator must report that it could not deliver, rather than
// quietly slowing down with the server.
func TestOpenLoopReportsClientSaturationInsteadOfSlowingDown(t *testing.T) {
	res := RunOpenLoop(context.Background(), Schedule{
		Rate: 2000, Duration: 300 * time.Millisecond, Workers: 1, Queue: 2, Seed: 2,
	}, "slow", func(ctx context.Context, r *rand.Rand) (Outcome, error) {
		time.Sleep(20 * time.Millisecond)
		return Outcome{}, nil
	})
	if res.Offered <= res.Completed {
		t.Errorf("expected offered (%d) to exceed completed (%d) when the server is slower than the rate",
			res.Offered, res.Completed)
	}
	if res.Dropped == 0 {
		t.Errorf("a full buffer must drop and count arrivals, not block the scheduler")
	}
	if !res.ClientSaturated {
		t.Errorf("client saturation must be reported; reason was %q", res.SaturationReason)
	}
}

// A synchronized burst offers Burst operations at one instant every Burst/Rate
// seconds, so the offered *rate* is unchanged while the arrival *shape* is not.
func TestOpenLoopBurstKeepsTheRateAndChangesTheShape(t *testing.T) {
	res := RunOpenLoop(context.Background(), Schedule{
		Rate: 50, Duration: 400 * time.Millisecond, Workers: 8, Queue: 512, Burst: 10, Seed: 3,
	}, "burst", func(ctx context.Context, r *rand.Rand) (Outcome, error) {
		return Outcome{}, nil
	})
	if res.Offered == 0 {
		t.Fatalf("a burst schedule offered nothing")
	}
	if res.Offered%10 != 0 {
		t.Errorf("offered %d with a burst of 10: bursts must be whole", res.Offered)
	}
	if res.Offered < 15 || res.Offered > 25 {
		t.Errorf("offered %d at 50/s over 400ms in bursts of 10; expected roughly 20", res.Offered)
	}
}

// A cadence run has to replay from its seed, or two runs of "the same" workload
// cannot be compared.
func TestOpenLoopIsReproducible(t *testing.T) {
	run := func() int64 {
		return RunOpenLoop(context.Background(), Schedule{
			Rate: 500, Duration: 200 * time.Millisecond, Workers: 4, Queue: 256,
			Jitter: true, Seed: 99,
		}, "jitter", func(ctx context.Context, r *rand.Rand) (Outcome, error) {
			return Outcome{}, nil
		}).Offered
	}
	if a, b := run(), run(); a != b {
		t.Errorf("same schedule and seed offered %d then %d operations", a, b)
	}
}

// A count bound is how a low cadence is validated without waiting out its
// period: the run stops after a fixed number of operations, inside a window
// deliberately left far longer.
func TestOpenLoopHonoursACountBound(t *testing.T) {
	start := time.Now()
	res := RunOpenLoop(context.Background(), Schedule{
		Rate: 1000, Duration: 5 * time.Second, Count: 7, Workers: 2, Queue: 16, Seed: 4,
	}, "compressed-daily", func(ctx context.Context, r *rand.Rand) (Outcome, error) {
		return Outcome{}, nil
	})
	if res.Offered != 7 || res.Completed != 7 {
		t.Errorf("count-bounded schedule offered %d completed %d, want 7 and 7", res.Offered, res.Completed)
	}
	if elapsed := time.Since(start); elapsed > time.Second {
		t.Errorf("a 7-operation bound ran for %s inside a 5s window; it must stop at the count", elapsed)
	}
}

// A count bound with no window derives one from the rate. Otherwise it would
// truncate at a default second and report a shortfall that is the caller's
// mistake as though it were the database's behaviour.
func TestCountBoundDerivesAWindowRatherThanTruncating(t *testing.T) {
	res := RunOpenLoop(context.Background(), Schedule{
		Rate: 2000, Count: 40, Workers: 4, Queue: 64, Seed: 5,
	}, "derived-window", func(ctx context.Context, r *rand.Rand) (Outcome, error) {
		return Outcome{}, nil
	})
	if res.Offered != 40 {
		t.Errorf("offered %d of a 40-operation bound with no explicit window", res.Offered)
	}
}
