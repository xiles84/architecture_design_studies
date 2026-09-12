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

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ---------------------------------------------------------------------------
// Latency accounting
// ---------------------------------------------------------------------------

// LatencyStats is reported in milliseconds. Percentiles are computed from every
// recorded sample rather than from a bucketed histogram: at these op counts the
// memory is trivial and exact tails matter more than the memory does. Tail
// latency is usually where a design decision actually shows up -- means hide it.
type LatencyStats struct {
	Count  int64   `json:"count"`
	MinMS  float64 `json:"min_ms"`
	MeanMS float64 `json:"mean_ms"`
	P50MS  float64 `json:"p50_ms"`
	P90MS  float64 `json:"p90_ms"`
	P95MS  float64 `json:"p95_ms"`
	P99MS  float64 `json:"p99_ms"`
	// Deep tails are omitted rather than guessed when there are too few samples
	// to support them -- see minSamplesFor* below. A zero here means "not enough
	// data", not "zero milliseconds", and the report renders it as a dash.
	P999MS  float64 `json:"p999_ms,omitempty"`
	P9999MS float64 `json:"p9999_ms,omitempty"`
	MaxMS   float64 `json:"max_ms"`
}

// A percentile is only meaningful when enough samples sit beyond it. These
// thresholds require at least ten samples above the quantile, which is the
// minimum at which the number describes a population rather than one unlucky
// request.
//
// This matters concretely here: within a single run, cell sample counts span
// 363 to 392 000. At 392 000 samples p99.9 is backed by ~392 observations and is
// a real measurement; at 363 samples it would be the top 0.36 of one sample --
// literally the maximum, relabelled as a percentile and looking far more
// authoritative than it is. The slow queries that produce small sample counts
// are exactly the ones a reader would most want a tail figure for, which is why
// this has to be enforced rather than left to judgement.
const (
	minSamplesForP999  = 10_000
	minSamplesForP9999 = 100_000
)

func summarize(samples []time.Duration) LatencyStats {
	if len(samples) == 0 {
		return LatencyStats{}
	}
	sort.Slice(samples, func(i, j int) bool { return samples[i] < samples[j] })
	ms := func(d time.Duration) float64 { return float64(d.Nanoseconds()) / 1e6 }
	pct := func(p float64) float64 {
		idx := int(p * float64(len(samples)-1))
		return ms(samples[idx])
	}
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
	if len(samples) >= minSamplesForP999 {
		st.P999MS = pct(0.999)
	}
	if len(samples) >= minSamplesForP9999 {
		st.P9999MS = pct(0.9999)
	}
	return st
}

// ---------------------------------------------------------------------------
// Parameter binding
// ---------------------------------------------------------------------------

// Binder turns the "-- params:" names in a catalogue into concrete values, so
// one driver can execute statements whose signatures differ between designs.
type Binder struct {
	Charities int64
	People    int64
	Donations int64
	// personCharity[i] is the charity of person id i+1, needed to write the
	// denormalised key in D3/D4/D5/D7.
	personCharity []int64
	nextDonation  atomic.Int64
}

func NewBinder(ds *Dataset) *Binder {
	b := &Binder{
		Charities:     int64(len(ds.Charities)),
		People:        int64(len(ds.People)),
		Donations:     int64(len(ds.Donations)),
		personCharity: make([]int64, len(ds.People)),
	}
	for i := range ds.People {
		b.personCharity[i] = ds.People[i].CharityID
	}
	b.nextDonation.Store(int64(len(ds.Donations)) + 1)
	return b
}

func bind(params []string, vals map[string]any) ([]any, error) {
	out := make([]any, 0, len(params))
	for _, p := range params {
		v, ok := vals[p]
		if !ok {
			return nil, fmt.Errorf("no value bound for parameter %q", p)
		}
		out = append(out, v)
	}
	return out, nil
}

// readVals draws a fresh random key set for one read execution. Fresh keys per
// execution are essential: repeating one key would measure the buffer cache, not
// the access path.
func (b *Binder) readVals(r *rand.Rand) map[string]any {
	return map[string]any{
		"charity_id":  r.Int63n(b.Charities) + 1,
		"person_id":   r.Int63n(b.People) + 1,
		"donation_id": r.Int63n(b.Donations) + 1,
	}
}

// ---------------------------------------------------------------------------
// Read benchmark
// ---------------------------------------------------------------------------

type QueryResult struct {
	Query      string  `json:"query"`
	Ops        int64   `json:"ops"`
	Errors     int64   `json:"errors"`
	FirstError string  `json:"first_error,omitempty"`
	DurationS  float64 `json:"duration_s"`
	// OpsPerSec is the MEDIAN across trials, not the mean. On a laptop a single
	// trial can be wrecked by an autovacuum pass, a checkpoint or the OS
	// scheduler moving a container onto an efficiency core; the median ignores
	// one bad trial, where a mean would absorb it and quietly report it as a
	// design difference.
	OpsPerSec float64      `json:"ops_per_sec"`
	Latency   LatencyStats `json:"latency"`

	// Per-trial throughput, and the spread across them. Spread is the honest
	// error bar: a 2x difference between designs means nothing if either design
	// varies by 2x between its own trials.
	Trials    []float64 `json:"trials_ops_per_sec,omitempty"`
	SpreadPct float64   `json:"spread_pct,omitempty"`

	// Destructive operations are measured by count, not by duration (see
	// finiteBudget). OpBudget records that count; Exhausted records that some
	// trial ran out of rows before its budget, which makes that cell suspect.
	OpBudget  int64 `json:"op_budget,omitempty"`
	Exhausted bool  `json:"pool_exhausted,omitempty"`
}

// median of a copy, so the caller's slice ordering survives.
func median(xs []float64) float64 {
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

// spreadPct is (max-min)/median, the width of the trial band relative to the
// reported number. Anything above roughly 20% means this cell should not be used
// to argue a small difference.
func spreadPct(xs []float64) float64 {
	if len(xs) < 2 {
		return 0
	}
	lo, hi := xs[0], xs[0]
	for _, v := range xs {
		if v < lo {
			lo = v
		}
		if v > hi {
			hi = v
		}
	}
	m := median(xs)
	if m <= 0 {
		return 0
	}
	return (hi - lo) / m * 100
}

type BenchOpts struct {
	Conns         int
	Duration      time.Duration
	Warmup        time.Duration
	MaxOps        int64
	Strategy      string // D5 only: blind | forupdate | optimistic
	Isolation     string // rc | rr | ser
	StmtTimeoutMS int
	// Trials repeats the measured phase; the reported throughput is the median.
	Trials int
	// MaxPerPerson caps donor history in the generated dataset (0 = unbounded).
	MaxPerPerson int
	// ReadMix selects the mixed-workload query weights: dashboard or portal.
	ReadMix string
	// Charities overrides the number of charities (0 = scale default).
	Charities int
}

// spreadNote flags a cell whose own trials disagree with each other by more than
// a fifth. Such a cell cannot support an argument about a small difference
// between designs, and saying so at the point of measurement is better than
// leaving a reader to discover it in the JSON.
func spreadNote(r QueryResult) string {
	if len(r.Trials) < 2 {
		return ""
	}
	flag := ""
	if r.SpreadPct > 20 {
		flag = "  <-- noisy"
	}
	return fmt.Sprintf("  [%d trials, spread %.0f%%]%s", len(r.Trials), r.SpreadPct, flag)
}

func isoLevel(s string) pgx.TxIsoLevel {
	switch s {
	case "rr":
		return pgx.RepeatableRead
	case "ser":
		return pgx.Serializable
	default:
		return pgx.ReadCommitted
	}
}

// runLoop drives `conns` workers against `op` for the configured duration, first
// through an unrecorded warmup so that plan caching, connection establishment and
// buffer warming do not land in the reported numbers.
func runLoop(ctx context.Context, o BenchOpts, name string, op func(ctx context.Context, r *rand.Rand) error) QueryResult {
	res := QueryResult{Query: name, OpBudget: o.MaxOps}
	var exhausted atomic.Bool

	warm := func(d time.Duration, record bool) ([]time.Duration, int64, int64, string) {
		if d <= 0 {
			return nil, 0, 0, ""
		}
		deadline := time.Now().Add(d)
		var (
			wg       sync.WaitGroup
			mu       sync.Mutex
			all      []time.Duration
			ops      atomic.Int64
			errCount atomic.Int64
			firstErr string
			errOnce  sync.Once
		)
		for w := 0; w < o.Conns; w++ {
			wg.Add(1)
			go func(seed int64) {
				defer wg.Done()
				r := rand.New(rand.NewSource(seed))
				var local []time.Duration
				for time.Now().Before(deadline) {
					if o.MaxOps > 0 && ops.Load() >= o.MaxOps {
						break
					}
					t0 := time.Now()
					err := op(ctx, r)
					el := time.Since(t0)
					if errors.Is(err, errExhausted) {
						// The finite pool this operation consumes is used up.
						// Stop this worker; do NOT count an error and do NOT loop.
						// Spinning here is what once turned "delete" into 111
						// million instant failures and a reported 0 ops/s.
						exhausted.Store(true)
						break
					}
					if err != nil {
						errCount.Add(1)
						errOnce.Do(func() { firstErr = err.Error() })
						if ctx.Err() != nil {
							return
						}
						continue
					}
					ops.Add(1)
					if record {
						local = append(local, el)
					}
				}
				if record && len(local) > 0 {
					mu.Lock()
					all = append(all, local...)
					mu.Unlock()
				}
			}(time.Now().UnixNano() + int64(w)*7919)
		}
		wg.Wait()
		return all, ops.Load(), errCount.Load(), firstErr
	}

	warm(o.Warmup, false)

	trials := o.Trials
	if trials < 1 {
		trials = 1
	}

	var allSamples []time.Duration
	for t := 0; t < trials; t++ {
		start := time.Now()
		samples, ops, errs, firstErr := warm(o.Duration, true)
		elapsed := time.Since(start).Seconds()

		res.Ops += ops
		res.Errors += errs
		res.DurationS += elapsed
		if res.FirstError == "" {
			res.FirstError = firstErr
		}
		if elapsed > 0 {
			res.Trials = append(res.Trials, float64(ops)/elapsed)
		}
		// Latency samples are pooled across trials: percentiles over the whole
		// population describe the workload better than a percentile of one trial.
		allSamples = append(allSamples, samples...)
	}

	res.OpsPerSec = median(res.Trials)
	res.SpreadPct = spreadPct(res.Trials)
	res.Exhausted = exhausted.Load()
	res.Latency = summarize(allSamples)
	return res
}

// BenchmarkReads runs every query in the design's catalogue in isolation.
//
// Isolated rather than mixed on purpose: a mixed workload gives one number that
// hides which query moved, and the point of this study is attributing cost to a
// specific question. A weighted mixed run is a separate experiment.
func BenchmarkReads(ctx context.Context, pool *pgxpool.Pool, d Design, b *Binder, o BenchOpts, only []string) ([]QueryResult, error) {
	src, err := readSQL(d.ID, "queries.sql")
	if err != nil {
		return nil, err
	}
	stmts, err := ParseCatalog(src)
	if err != nil {
		return nil, err
	}
	keep := map[string]bool{}
	for _, q := range only {
		keep[q] = true
	}

	var out []QueryResult
	for _, st := range stmts {
		if len(keep) > 0 && !keep[st.Name] {
			continue
		}
		st := st
		sql := st.SQL
		r := runLoop(ctx, o, st.Name, func(ctx context.Context, rnd *rand.Rand) error {
			args, err := bind(st.Params, b.readVals(rnd))
			if err != nil {
				return err
			}
			rows, err := pool.Query(ctx, sql, args...)
			if err != nil {
				return err
			}
			defer rows.Close()
			for rows.Next() {
				// Drain so the server actually materialises every row; without
				// this a LIMIT-less aggregate could be measured unfairly.
			}
			return rows.Err()
		})
		fmt.Printf("    %-30s %10.1f ops/s  p50=%7.3fms p99=%8.3fms  errors=%d%s\n",
			r.Query, r.OpsPerSec, r.Latency.P50MS, r.Latency.P99MS, r.Errors, spreadNote(r))
		out = append(out, r)
	}
	return out, nil
}

// ---------------------------------------------------------------------------
// Write benchmark
// ---------------------------------------------------------------------------

type WriteResult struct {
	QueryResult
	Op        string `json:"op"`
	Strategy  string `json:"strategy,omitempty"`
	Isolation string `json:"isolation,omitempty"`
	Retries   int64  `json:"retries"`
	// Skipped explains why a cell has no number, so the gap is a documented
	// decision rather than something a reader has to guess at.
	Skipped string `json:"skipped,omitempty"`
	// Audit is the correctness check run on the state this op produced.
	Audit *RollupAudit `json:"rollup_audit,omitempty"`
}

var errNoop = errors.New("statement affected no rows")

// errExhausted means a destructive operation has used up the finite set of rows
// it is allowed to consume (donation ids for delete, person ids for erasure).
// It is not a failure of the design and must never be counted as one: it means
// "stop", and runLoop treats it that way.
var errExhausted = errors.New("finite id pool exhausted")

// finiteBudget sizes a destructive operation so warmup and every trial fit
// inside HALF of the pool it consumes.
//
// Duration-based measurement is wrong for these operations. A delete can only
// happen once per row, so at a few thousand deletes per second the pool drains in
// seconds -- the first repeated-trials pass used 3 s of warmup plus 5 x 8 s of
// trials against 108 000 donations, ran out partway through, and every later
// trial measured nothing but instant "no such row" errors. Each trial therefore
// gets a fixed number of operations instead, and its throughput is that count
// over the time it actually took.
//
// Half the pool, not all of it, because deleting the last rows of a table is not
// representative of deleting rows from a populated one.
func finiteBudget(pool int64, trials int) int64 {
	if trials < 1 {
		trials = 1
	}
	b := pool / int64(2*(trials+1))
	if b < 1 {
		b = 1
	}
	return b
}

// ---------------------------------------------------------------------------
// Workload mix (experiment D)
//
// Isolated benchmarks answer "what does this design do to this query". They
// cannot answer "does that advantage survive a real workload", because several
// costs only exist when reads and writes overlap in time:
//
//   * MVCC bloat -- D6 rewrites a whole person row per donation, and those dead
//     versions pile up WHILE readers are scanning. Measured alone, the read
//     benchmark runs against a clean table and never meets them.
//   * hot-row contention -- every write touches D4's charity rollup row, so a
//     reader of that row walks version chains it never sees in isolation.
//   * vacuum and compaction -- autovacuum in PostgreSQL, RocksDB compaction in
//     YugabyteDB. Both are triggered by writes and both steal CPU from reads.
//   * buffer-cache and index-page competition.
//
// Every one of those penalises the designs that buy read speed with redundancy
// -- which is to say, the designs this study is mostly about. Isolation
// therefore biases the results in favour of exactly the designs under scrutiny,
// and that is why this experiment exists rather than being optional.
//
// The mix below is a donor portal plus a dashboard: mostly per-donor reads, a
// steady feed, and occasional aggregate panels. Weights sum to 100.
// ---------------------------------------------------------------------------

type weighted struct {
	Name   string
	Weight int
}

var readMix = []weighted{
	{"q09_person_recent_donations", 30}, // the donor's own history page
	{"q10_donation_by_id", 20},          // receipts, deep links
	{"q12_charity_recent_feed", 15},     // "recent activity" panel
	{"q11_person_donation_count", 10},   // badge on the donor header
	{"q02_last_donation_charity", 8},
	{"q06_person_first_last", 5},
	{"q08_total_donated_charity", 5},
	{"q03_top_donor_charity", 4}, // leaderboard panel
	{"q07_total_donated_global", 3},
}

// portalMix is the regime in which the embedded design's biggest weakness does not
// apply. D6 collapses on every question that crosses donors -- totals, leaderboards,
// feeds, global recency -- because it has to unnest every array in the table. A
// donor-facing portal never asks those questions: each request is about ONE
// donor. If D6 wins anywhere, it wins here, and the study should say so with a
// measurement rather than leave "embedding is bad" as an unqualified conclusion.
//
// Point lookups by donation id are excluded too: in a portal the donor reaches a
// donation by browsing their own history, not by global id -- and id lookup is the
// other access path D6 structurally lacks.
var portalMix = []weighted{
	{"q09_person_recent_donations", 60},
	{"q06_person_first_last", 20},
	{"q11_person_donation_count", 20},
}

func readMixFor(name string) ([]weighted, error) {
	switch name {
	case "", "dashboard":
		return readMix, nil
	case "portal":
		return portalMix, nil
	}
	return nil, fmt.Errorf("unknown read mix %q (want dashboard or portal)", name)
}

// Writes are dominated by new donations; corrections and refunds are rare but
// disproportionately expensive in several designs, so they must not be zero.
//
// update_person is included because it is the one common write that originates
// at the PARENT rather than propagating up from a child, and the designs differ
// enormously in what they have made a person row into.
//
// delete_person is deliberately EXCLUDED. Erasure removes donors that concurrent
// inserts are still choosing at random, which would produce a stream of foreign
// key violations that say nothing about any design -- the error rate would be an
// artefact of the harness picking a deleted id. It is measured in the isolated
// write benchmark instead, where it can have the id space to itself.
var writeMix = []weighted{
	{"insert", 85},
	{"update", 7},
	{"update_person", 5},
	{"delete", 3},
}

func pickWeighted(ws []weighted, r *rand.Rand) string {
	total := 0
	for _, w := range ws {
		total += w.Weight
	}
	x := r.Intn(total)
	for _, w := range ws {
		if x < w.Weight {
			return w.Name
		}
		x -= w.Weight
	}
	return ws[len(ws)-1].Name
}

// BenchmarkWrites measures the three mutations that matter for a tree schema:
// appending a child, correcting one in place, and removing one. Insert is the
// hot path in production; update and delete are where rollup and embedded
// designs reveal costs that inserts alone never show.
func BenchmarkWrites(ctx context.Context, pool *pgxpool.Pool, d Design, b *Binder, o BenchOpts, ops []string) ([]WriteResult, error) {
	src, err := readSQL(d.ID, "writes.sql")
	if err != nil {
		return nil, err
	}
	stmts, err := ParseCatalog(src)
	if err != nil {
		return nil, err
	}
	w := StmtMap(stmts)
	var retries atomic.Int64
	var out []WriteResult

	for _, op := range ops {
		retries.Store(0)
		var fn func(context.Context, *rand.Rand) error

		// D5 implements application-maintained rollups on the INSERT path only.
		// Correcting or removing a donation would require recomputing the MIN and
		// MAX, which cannot be derived from a delta -- the work D4's trigger does
		// in donation_rollup_apply().
		//
		// Measuring those paths anyway would produce a number that looks
		// comparable to D4 and is not: D5 would simply be doing less work, and
		// leaving the aggregates wrong afterwards. Refusing to produce the number
		// is better than publishing one that invites a false conclusion. D4 vs D5
		// is therefore compared on the insert path, which is where the
		// concurrency question lives in any case.
		// Only the operations that CHANGE an aggregate are affected. Editing a
		// donor's email touches no rollup, so D5 runs it exactly like every other
		// design and the comparison stands.
		rollupSensitive := map[string]bool{"update": true, "delete": true, "delete_person": true}
		if d.AppRollup && rollupSensitive[op] {
			reason := "D5 maintains its rollups only on the insert path, so this " +
				"operation would do strictly less work than D4's trigger and leave " +
				"the aggregates stale; the number would not be comparable"
			fmt.Printf("    %-30s skipped — %s\n", "write:"+op, reason)
			out = append(out, WriteResult{
				QueryResult: QueryResult{Query: op}, Op: op,
				Strategy: o.Strategy, Isolation: o.Isolation, Skipped: reason,
			})
			continue
		}

		var cursor atomic.Int64
		fn, err = b.writeFn(pool, d, w, o, op, &retries, &cursor)
		if err != nil {
			return nil, err
		}

		// Destructive operations are measured by a fixed count per trial, sized
		// to fit warmup plus every trial inside half the rows they consume.
		oo := o
		switch op {
		case "delete":
			oo.MaxOps = finiteBudget(b.Donations, o.Trials)
		case "delete_person":
			oo.MaxOps = finiteBudget(b.People, o.Trials)
		}

		r := runLoop(ctx, oo, op, fn)
		if r.Exhausted {
			fmt.Printf("    WARNING %s ran out of rows before its budget; treat this cell as suspect\n", op)
		}
		wr := WriteResult{QueryResult: r, Op: op, Retries: retries.Load()}
		if d.AppRollup {
			wr.Strategy = o.Strategy
			wr.Isolation = o.Isolation
		}
		fmt.Printf("    %-30s %10.1f ops/s  p50=%7.3fms p99=%8.3fms  errors=%d retries=%d%s\n",
			"write:"+op, r.OpsPerSec, r.Latency.P50MS, r.Latency.P99MS, r.Errors, wr.Retries, spreadNote(r))
		out = append(out, wr)
	}
	return out, nil
}

// writeFn builds the driver for one write operation. Factored out of
// BenchmarkWrites so that the mixed workload can issue the same operations,
// against the same statements, without a second implementation that could drift
// away from this one.
//
// cursor is supplied by the caller because deletes must advance monotonically
// across the whole run: sharing one cursor guarantees two workers never target
// the same row and never double-delete, which would otherwise be recorded as a
// suspiciously fast success.
func (b *Binder) writeFn(pool *pgxpool.Pool, d Design, w map[string]Stmt, o BenchOpts,
	op string, retries *atomic.Int64, cursor *atomic.Int64) (func(context.Context, *rand.Rand) error, error) {

	switch op {
	case "insert":
		return b.insertFn(pool, d, w, o, retries), nil

	case "update":
		st := w["w_update_amount"]
		return func(ctx context.Context, r *rand.Rand) error {
			vals := map[string]any{
				"amount_cents": int64(500 + r.Intn(500000)),
				"donation_id":  r.Int63n(b.Donations) + 1,
			}
			args, err := bind(st.Params, vals)
			if err != nil {
				return err
			}
			_, err = pool.Exec(ctx, st.SQL, args...)
			return err
		}, nil

	case "delete":
		st := w["w_delete_donation"]
		return func(ctx context.Context, r *rand.Rand) error {
			id := cursor.Add(1)
			if id > b.Donations {
				return errExhausted
			}
			args, err := bind(st.Params, map[string]any{"donation_id": id})
			if err != nil {
				return err
			}
			tag, err := pool.Exec(ctx, st.SQL, args...)
			if err != nil {
				return err
			}
			if tag.RowsAffected() == 0 {
				return errNoop
			}
			return nil
		}, nil

	case "update_person":
		// A donor edits their profile: a write that ORIGINATES at the parent
		// rather than propagating up from a child. Identical SQL in every design;
		// only the shape of the row it lands on differs.
		st := w["w_update_person_profile"]
		return func(ctx context.Context, r *rand.Rand) error {
			pid := r.Int63n(b.People) + 1
			args, err := bind(st.Params, map[string]any{
				"email":     fmt.Sprintf("donor%d+%d@example.org", pid, r.Intn(1_000_000)),
				"person_id": pid,
			})
			if err != nil {
				return err
			}
			_, err = pool.Exec(ctx, st.SQL, args...)
			return err
		}, nil

	case "delete_person":
		// Erasure. Walks the person id space monotonically so no donor is deleted
		// twice -- a second delete would affect no rows and be recorded as a
		// suspiciously fast success.
		st := w["w_delete_person_cascade"]
		return func(ctx context.Context, r *rand.Rand) error {
			id := cursor.Add(1)
			if id > b.People {
				return errExhausted
			}
			args, err := bind(st.Params, map[string]any{"person_id": id})
			if err != nil {
				return err
			}
			tag, err := pool.Exec(ctx, st.SQL, args...)
			if err != nil {
				return err
			}
			if tag.RowsAffected() == 0 {
				return errNoop
			}
			return nil
		}, nil
	}
	return nil, fmt.Errorf("unknown write op %q", op)
}

// insertFn builds the append path for a design. This is where the designs stop
// looking alike: one statement for D1-D4, a multi-statement transaction for D5,
// and a whole-document rewrite for D6.
func (b *Binder) insertFn(pool *pgxpool.Pool, d Design, w map[string]Stmt, o BenchOpts, retries *atomic.Int64) func(context.Context, *rand.Rand) error {
	ins := w["w_insert_donation"]

	newDonation := func(r *rand.Rand) map[string]any {
		pid := r.Int63n(b.People) + 1
		note := noteTemplates[r.Intn(len(noteTemplates))]
		return map[string]any{
			"donation_id":  b.nextDonation.Add(1),
			"person_id":    pid,
			"charity_id":   b.personCharity[pid-1],
			"amount_cents": int64(500 + r.Intn(500000)),
			"currency":     "USD",
			"donated_at":   time.Now().UTC(),
			"note":         note,
		}
	}

	if !d.AppRollup {
		// D1-D4 and D6: a single statement. For D4 the trigger fans it out
		// server-side; for D6 the statement rewrites the person document.
		return func(ctx context.Context, r *rand.Rand) error {
			args, err := bind(ins.Params, newDonation(r))
			if err != nil {
				return err
			}
			_, err = pool.Exec(ctx, ins.SQL, args...)
			return err
		}
	}

	// D5: the application performs the same consolidation the trigger would,
	// under a selectable concurrency-control strategy.
	rp, rc := w["w_rollup_person_ins"], w["w_rollup_charity_ins"]
	rpCAS, rcCAS := w["w_rollup_person_ins_cas"], w["w_rollup_charity_ins_cas"]
	lockP, lockC := w["w_lock_person"], w["w_lock_charity"]
	readP, readC := w["w_read_person_version"], w["w_read_charity_version"]

	return func(ctx context.Context, r *rand.Rand) error {
		const maxAttempts = 50
		for attempt := 0; attempt < maxAttempts; attempt++ {
			vals := newDonation(r)
			err := func() error {
				tx, err := pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: isoLevel(o.Isolation)})
				if err != nil {
					return err
				}
				defer tx.Rollback(ctx)

				if o.Strategy == "forupdate" {
					// Deterministic lock order (person before charity) removes
					// the deadlock this pattern otherwise invites.
					for _, l := range []struct {
						st  Stmt
						key string
					}{{lockP, "person_id"}, {lockC, "charity_id"}} {
						a, err := bind(l.st.Params, vals)
						if err != nil {
							return err
						}
						var v int64
						if err := tx.QueryRow(ctx, l.st.SQL, a...).Scan(&v); err != nil {
							return err
						}
					}
				}

				var pv, cv int64
				if o.Strategy == "optimistic" {
					a, _ := bind(readP.Params, vals)
					if err := tx.QueryRow(ctx, readP.SQL, a...).Scan(&pv); err != nil {
						return err
					}
					a, _ = bind(readC.Params, vals)
					if err := tx.QueryRow(ctx, readC.SQL, a...).Scan(&cv); err != nil {
						return err
					}
					vals["version"] = pv
				}

				a, err := bind(ins.Params, vals)
				if err != nil {
					return err
				}
				if _, err := tx.Exec(ctx, ins.SQL, a...); err != nil {
					return err
				}

				pStmt, cStmt := rp, rc
				if o.Strategy == "optimistic" {
					pStmt, cStmt = rpCAS, rcCAS
				}

				vals["version"] = pv
				a, err = bind(pStmt.Params, vals)
				if err != nil {
					return err
				}
				tag, err := tx.Exec(ctx, pStmt.SQL, a...)
				if err != nil {
					return err
				}
				if o.Strategy == "optimistic" && tag.RowsAffected() == 0 {
					return errNoop
				}

				vals["version"] = cv
				a, err = bind(cStmt.Params, vals)
				if err != nil {
					return err
				}
				tag, err = tx.Exec(ctx, cStmt.SQL, a...)
				if err != nil {
					return err
				}
				if o.Strategy == "optimistic" && tag.RowsAffected() == 0 {
					return errNoop
				}
				return tx.Commit(ctx)
			}()

			if err == nil {
				return nil
			}
			if isRetryable(err) && attempt < maxAttempts-1 {
				retries.Add(1)
				continue
			}
			return err
		}
		return fmt.Errorf("gave up after %d attempts", maxAttempts)
	}
}

// isRetryable covers the two ways this workload legitimately fails under
// contention: an optimistic version check that lost, and a serialization
// failure raised by the engine. Both mean "try again", and counting them is
// half the point of experiment C.
func isRetryable(err error) bool {
	if errors.Is(err, errNoop) {
		return true
	}
	var pgErr interface{ SQLState() string }
	if errors.As(err, &pgErr) {
		switch pgErr.SQLState() {
		case "40001", "40P01": // serialization_failure, deadlock_detected
			return true
		}
	}
	return false
}
