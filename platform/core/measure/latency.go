// Package measure is the study-independent measurement core: latency
// summaries, trial statistics and the closed-loop driver.
//
// Everything here was first written inside study 01's harness and moved here
// unchanged in behaviour. Each rule encodes a lesson that cost a wrong number
// once; see LESSONS_LEARNED.md for the incidents.
package measure

import (
	"sort"
	"time"
)

// LatencyStats is reported in milliseconds. Percentiles are computed from every
// recorded sample rather than from a bucketed histogram: at these op counts the
// memory is trivial and exact tails matter more than the memory does.
type LatencyStats struct {
	Count  int64   `json:"count"`
	MinMS  float64 `json:"min_ms"`
	MeanMS float64 `json:"mean_ms"`
	P50MS  float64 `json:"p50_ms"`
	P90MS  float64 `json:"p90_ms"`
	P95MS  float64 `json:"p95_ms"`
	P99MS  float64 `json:"p99_ms"`
	// Deep tails are omitted rather than guessed when there are too few samples
	// to support them. A zero means "not enough data", never "zero milliseconds".
	P999MS  float64 `json:"p999_ms,omitempty"`
	P9999MS float64 `json:"p9999_ms,omitempty"`
	MaxMS   float64 `json:"max_ms"`
}

// A percentile is only meaningful when at least ten samples sit beyond it.
// Study 01 saw cells from 363 to 392 000 samples in one run; at 363 a "p99.9"
// is the maximum relabelled, and those low-count cells are exactly the slow
// ones a reader most wants a tail figure for.
const (
	MinSamplesForP999  = 10_000
	MinSamplesForP9999 = 100_000
)

// Summarize sorts samples in place and returns their distribution.
func Summarize(samples []time.Duration) LatencyStats {
	if len(samples) == 0 {
		return LatencyStats{}
	}
	sort.Slice(samples, func(i, j int) bool { return samples[i] < samples[j] })
	ms := func(d time.Duration) float64 { return float64(d.Nanoseconds()) / 1e6 }
	pct := func(p float64) float64 { return ms(samples[int(p*float64(len(samples)-1))]) }
	var sum time.Duration
	for _, s := range samples {
		sum += s
	}
	st := LatencyStats{
		Count:  int64(len(samples)),
		MinMS:  ms(samples[0]),
		MeanMS: ms(sum) / float64(len(samples)),
		P50MS:  pct(0.50),
		P90MS:  pct(0.90),
		P95MS:  pct(0.95),
		P99MS:  pct(0.99),
		MaxMS:  ms(samples[len(samples)-1]),
	}
	if len(samples) >= MinSamplesForP999 {
		st.P999MS = pct(0.999)
	}
	if len(samples) >= MinSamplesForP9999 {
		st.P9999MS = pct(0.9999)
	}
	return st
}

// Median of a copy. The median, not the mean, is the reported throughput across
// trials: it ignores one trial wrecked by an autovacuum pass or a scheduler
// migration onto an efficiency core, where a mean would launder it into a
// finding.
func Median(xs []float64) float64 {
	if len(xs) == 0 {
		return 0
	}
	c := append([]float64(nil), xs...)
	sort.Float64s(c)
	n := len(c)
	if n%2 == 1 {
		return c[n/2]
	}
	return (c[n/2-1] + c[n/2]) / 2
}

// SpreadPct is (max-min)/median: the honest error bar of a repeated
// measurement. Above ~20% a cell cannot support an argument about a small
// difference between designs.
func SpreadPct(xs []float64) float64 {
	if len(xs) < 2 {
		return 0
	}
	lo, hi := xs[0], xs[0]
	for _, v := range xs {
		lo = min(lo, v)
		hi = max(hi, v)
	}
	m := Median(xs)
	if m <= 0 {
		return 0
	}
	return (hi - lo) / m * 100
}
