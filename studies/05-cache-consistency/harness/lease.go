package main

import (
	"context"
	"errors"
	"math/rand"
	"time"
)

// The per-key fill lease, and the calibration that decides how long it lasts.
//
// The lease exists to stop a stampede: without it, a cold or soft-expired key is
// loaded once per concurrent reader, and the database sees the burst rather than
// the cache. The lease is NOT a lock on the data -- writers do not take it -- it is
// a claim on the RIGHT TO FILL, so the correctness of a value still rests on
// publishing only committed state and on version fencing.

var errLeaseWaitTimeout = errors.New("fill lease did not resolve within the wait bound")

// LeaseStats is what the report needs about leases. Every field is a question the
// protocol asks: did the lease prevent duplicate fills, at what waiting cost, and
// what happened when a holder died or a waiter gave up.
type LeaseStats struct {
	Acquisitions  int64   `json:"acquisitions"`
	Contended     int64   `json:"contended"`
	Waits         int64   `json:"waits"`
	WaitP50MS     float64 `json:"wait_p50_ms"`
	WaitP99MS     float64 `json:"wait_p99_ms"`
	WaitMaxMS     float64 `json:"wait_max_ms"`
	Timeouts      int64   `json:"timeouts"`
	Lost          int64   `json:"lost"`
	DuplicateFill int64   `json:"duplicate_fills"`
	FallbackReads int64   `json:"fallback_reads"`
	// FillsPerBurst is the measured database-load count per synchronized miss
	// burst divided by the number of keys in the burst. 1.0 means the lease worked:
	// one load per key per burst, however many readers arrived at once.
	FillsPerBurst float64 `json:"fills_per_key_per_burst"`
	LeaseDuration string  `json:"lease_duration"`
	Calibration   string  `json:"calibration"`
}

// msStat is a small percentile summary for durations in milliseconds.
type msStat struct {
	Count int64   `json:"count"`
	P50   float64 `json:"p50_ms"`
	P99   float64 `json:"p99_ms"`
	Max   float64 `json:"max_ms"`
}

// calibrateLeaseDuration turns the measured database fill latency into a lease
// duration. The rule is fixed before any measurement runs: four times the p99 of
// the fill, clamped to [100 ms, 2 s].
//
// The reasoning, recorded because the protocol requires it: a lease shorter than a
// fill expires while its holder is still loading, a second reader steals it, and the
// miss burst produces exactly the duplicate loads the lease was meant to prevent.
// Four times the p99 leaves room for the tail without holding a dead key for
// seconds. An arbitrary large constant would pass a stampede test while hiding a
// cache that stalls for two seconds whenever a key expires.
func calibrateLeaseDuration(p99Fill time.Duration) (time.Duration, string) {
	d := 4 * p99Fill
	reason := "4 x measured p99 fill latency"
	if d < 100*time.Millisecond {
		d = 100 * time.Millisecond
		reason = "4 x measured p99 fill latency, raised to the 100 ms floor"
	}
	if d > 2*time.Second {
		d = 2 * time.Second
		reason = "4 x measured p99 fill latency, capped at 2 s"
	}
	reason += " (p99 fill = " + p99Fill.String() + ")"
	return d, reason
}

// jitteredBackoff is a bounded exponential backoff with jitter. Jitter matters:
// without it every waiter that lost the lease retries on the same schedule and the
// retries themselves become a burst.
func jitteredBackoff(r *rand.Rand, attempt int) time.Duration {
	base := 2 * time.Millisecond
	for i := 0; i < attempt && base < 200*time.Millisecond; i++ {
		base *= 2
	}
	if base > 200*time.Millisecond {
		base = 200 * time.Millisecond
	}
	if base <= time.Millisecond {
		return base
	}
	jitter := time.Duration(r.Int63n(int64(base/2) + 1))
	return base/2 + jitter
}

// waitForLease bounds how long a waiter will wait before falling back to an
// authoritative database read. The bound is derived from the lease duration, so it
// scales with the same calibration rather than being an unrelated constant.
func waitForLease(ctx context.Context, leaseTTL time.Duration, r *rand.Rand) (time.Duration, error) {
	deadline := time.Now().Add(leaseTTL)
	start := time.Now()
	for attempt := 0; ; attempt++ {
		if ctx.Err() != nil {
			return time.Since(start), ctx.Err()
		}
		if time.Now().After(deadline) {
			return time.Since(start), errLeaseWaitTimeout
		}
		select {
		case <-ctx.Done():
			return time.Since(start), ctx.Err()
		case <-time.After(jitteredBackoff(r, attempt)):
		}
	}
}
