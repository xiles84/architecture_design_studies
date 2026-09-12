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

	"github.com/jackc/pgx/v5/pgxpool"
)

// ---------------------------------------------------------------------------
// Experiment D — reads and writes running at the same time.
//
// Readers and writers share the database for the whole measured window, drawing
// from the weighted mixes in bench.go. Each query and each write operation is
// timed separately, so the result keeps the attribution that isolated
// benchmarking gives while exposing the interference that isolated benchmarking
// hides.
//
// The number that matters is not the throughput in this run -- it is the RATIO
// against the isolated run of the same design. That ratio is the answer to
// "does this design's advantage survive a real workload?", and it is the only
// question a person choosing a schema actually needs answered.
// ---------------------------------------------------------------------------

type MixedResult struct {
	ReadWorkers  int `json:"read_workers"`
	WriteWorkers int `json:"write_workers"`

	// Throughput aggregated across the whole mix, which is what an operator sees.
	ReadOpsPerSec  float64 `json:"read_ops_per_sec"`
	WriteOpsPerSec float64 `json:"write_ops_per_sec"`
	DurationS      float64 `json:"duration_s"`
	Errors         int64   `json:"errors"`
	Retries        int64   `json:"retries"`

	// Per-operation detail, so a degraded blend can be traced to the operation
	// that degraded. Without this the experiment would only say "slower".
	PerQuery    map[string]LatencyStats `json:"per_query_latency"`
	PerQueryOps map[string]float64      `json:"per_query_ops_per_sec"`
	PerWrite    map[string]LatencyStats `json:"per_write_latency"`
	PerWriteOps map[string]float64      `json:"per_write_ops_per_sec"`

	ReadMix  map[string]int `json:"read_mix_weights"`
	WriteMix map[string]int `json:"write_mix_weights"`

	// Did the consolidated aggregates and embedded caches survive a workload
	// where reads and writes ran together? Correctness under interference is a
	// different question from correctness under a pure write burst.
	Audit *RollupAudit `json:"rollup_audit,omitempty"`
}

// collector accumulates per-operation timings from many goroutines.
type collector struct {
	mu      sync.Mutex
	samples map[string][]time.Duration
	counts  map[string]int64
}

func newCollector() *collector {
	return &collector{samples: map[string][]time.Duration{}, counts: map[string]int64{}}
}

func (c *collector) add(name string, d time.Duration) {
	c.mu.Lock()
	c.samples[name] = append(c.samples[name], d)
	c.counts[name]++
	c.mu.Unlock()
}

func (c *collector) total() int64 {
	var n int64
	for _, v := range c.counts {
		n += v
	}
	return n
}

// BenchmarkMixed runs readers and writers concurrently for one measured window.
func BenchmarkMixed(ctx context.Context, pool *pgxpool.Pool, d Design, b *Binder,
	o BenchOpts, readWorkers, writeWorkers int) (*MixedResult, error) {

	qsrc, err := readSQL(d.ID, "queries.sql")
	if err != nil {
		return nil, err
	}
	qstmts, err := ParseCatalog(qsrc)
	if err != nil {
		return nil, err
	}
	queries := StmtMap(qstmts)

	wsrc, err := readSQL(d.ID, "writes.sql")
	if err != nil {
		return nil, err
	}
	wstmts, err := ParseCatalog(wsrc)
	if err != nil {
		return nil, err
	}
	w := StmtMap(wstmts)

	// D5 maintains its rollups only on the insert path (see BenchmarkWrites), so
	// including corrections and refunds for it would quietly measure a design
	// that leaves its aggregates stale. Its mix is inserts only, and the result
	// records that rather than hiding it.
	rmix, err := readMixFor(o.ReadMix)
	if err != nil {
		return nil, err
	}

	wmix := writeMix
	if d.AppRollup {
		wmix = []weighted{{"insert", 100}}
	}

	var retries, errCount atomic.Int64
	var cursor atomic.Int64

	writeFns := map[string]func(context.Context, *rand.Rand) error{}
	for _, m := range wmix {
		fn, err := b.writeFn(pool, d, w, o, m.Name, &retries, &cursor)
		if err != nil {
			return nil, err
		}
		writeFns[m.Name] = fn
	}

	readOne := func(ctx context.Context, r *rand.Rand, name string) error {
		st, ok := queries[name]
		if !ok {
			return fmt.Errorf("query %s missing from catalogue", name)
		}
		args, err := bind(st.Params, b.readVals(r))
		if err != nil {
			return err
		}
		rows, err := pool.Query(ctx, st.SQL, args...)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
		}
		return rows.Err()
	}

	run := func(d time.Duration, record bool) (*collector, *collector, float64) {
		reads, writes := newCollector(), newCollector()
		if d <= 0 {
			return reads, writes, 0
		}
		deadline := time.Now().Add(d)
		var wg sync.WaitGroup
		start := time.Now()

		for i := 0; i < readWorkers; i++ {
			wg.Add(1)
			go func(seed int64) {
				defer wg.Done()
				r := rand.New(rand.NewSource(seed))
				for time.Now().Before(deadline) {
					name := pickWeighted(rmix, r)
					t0 := time.Now()
					err := readOne(ctx, r, name)
					el := time.Since(t0)
					if err != nil {
						errCount.Add(1)
						if ctx.Err() != nil {
							return
						}
						continue
					}
					if record {
						reads.add(name, el)
					}
				}
			}(time.Now().UnixNano() + int64(i)*7919)
		}

		for i := 0; i < writeWorkers; i++ {
			wg.Add(1)
			go func(seed int64) {
				defer wg.Done()
				r := rand.New(rand.NewSource(seed))
				for time.Now().Before(deadline) {
					name := pickWeighted(wmix, r)
					t0 := time.Now()
					err := writeFns[name](ctx, r)
					el := time.Since(t0)
					if errors.Is(err, errExhausted) {
						// Deletes are 3% of the write mix and cannot drain the
						// pool in one window at these rates, but if they ever do,
						// skip the op rather than counting a design error.
						continue
					}
					if err != nil {
						errCount.Add(1)
						if ctx.Err() != nil {
							return
						}
						continue
					}
					if record {
						writes.add(name, el)
					}
				}
			}(time.Now().UnixNano() + int64(i)*104729)
		}

		wg.Wait()
		return reads, writes, time.Since(start).Seconds()
	}

	run(o.Warmup, false)
	errCount.Store(0)
	retries.Store(0)

	reads, writes, elapsed := run(o.Duration, true)

	res := &MixedResult{
		ReadWorkers: readWorkers, WriteWorkers: writeWorkers,
		DurationS:   elapsed,
		Errors:      errCount.Load(),
		Retries:     retries.Load(),
		PerQuery:    map[string]LatencyStats{},
		PerQueryOps: map[string]float64{},
		PerWrite:    map[string]LatencyStats{},
		PerWriteOps: map[string]float64{},
		ReadMix:     map[string]int{},
		WriteMix:    map[string]int{},
	}
	for _, m := range rmix {
		res.ReadMix[m.Name] = m.Weight
	}
	for _, m := range wmix {
		res.WriteMix[m.Name] = m.Weight
	}
	if elapsed > 0 {
		res.ReadOpsPerSec = float64(reads.total()) / elapsed
		res.WriteOpsPerSec = float64(writes.total()) / elapsed
	}
	for name, s := range reads.samples {
		res.PerQuery[name] = summarize(s)
		if elapsed > 0 {
			res.PerQueryOps[name] = float64(len(s)) / elapsed
		}
	}
	for name, s := range writes.samples {
		res.PerWrite[name] = summarize(s)
		if elapsed > 0 {
			res.PerWriteOps[name] = float64(len(s)) / elapsed
		}
	}

	fmt.Printf("    %d readers + %d writers: %.0f reads/s, %.0f writes/s, errors=%d retries=%d\n",
		readWorkers, writeWorkers, res.ReadOpsPerSec, res.WriteOpsPerSec, res.Errors, res.Retries)
	names := make([]string, 0, len(res.PerQueryOps))
	for n := range res.PerQueryOps {
		names = append(names, n)
	}
	sort.Strings(names)
	for _, n := range names {
		l := res.PerQuery[n]
		fmt.Printf("      %-30s %9.1f ops/s  p50=%7.3fms p99=%8.3fms\n",
			n, res.PerQueryOps[n], l.P50MS, l.P99MS)
	}
	return res, nil
}
