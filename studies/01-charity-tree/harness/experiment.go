package main

import (
	"context"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// These controls belong to the follow-up experiment, not to the legacy survey.
// A trial is an independently loaded process, not another timing window on a
// database that faster designs have already grown by different amounts.
type ExperimentOptions struct {
	Mode              string  `json:"mode"`
	Condition         string  `json:"condition"`
	Trial             int     `json:"trial"`
	Order             int     `json:"order"`
	Preparation       string  `json:"preparation"`
	HistoryMultiplier int     `json:"history_multiplier"`
	ReadWorkers       int     `json:"read_workers"`
	WriteWorkers      int     `json:"write_workers"`
	ArrivalRate       float64 `json:"arrival_rate"`
	Queue             int     `json:"queue_capacity"`
	HotProbability    float64 `json:"hot_probability"`
	HotDonors         int     `json:"hot_donors"`
	Cycles            int     `json:"cycles"`
	Batch             int     `json:"batch"`
	DBMemoryBytes     int64   `json:"db_memory_limit_bytes"`
	Targeted          bool    `json:"targeted_reads"`
}

func experimentFlags() *ExperimentOptions {
	o := &ExperimentOptions{}
	flag.StringVar(&o.Mode, "experiment", "reads", "reads | writes | arrival | growth | verify")
	flag.StringVar(&o.Condition, "condition", "default", "named experimental condition")
	flag.IntVar(&o.Trial, "trial-id", 1, "independently loaded trial number")
	flag.IntVar(&o.Order, "order", 1, "execution position within this trial")
	flag.StringVar(&o.Preparation, "preparation", "analyze", "analyze | vacuum (PostgreSQL only)")
	flag.IntVar(&o.HistoryMultiplier, "history-multiplier", 1, "multiply each generated donor history length")
	flag.IntVar(&o.ReadWorkers, "read-workers", 4, "closed-loop readers, held fixed across offered write rates")
	flag.IntVar(&o.WriteWorkers, "write-workers", 8, "workers draining scheduled write arrivals")
	flag.Float64Var(&o.ArrivalRate, "arrival-rate", 500, "scheduled inserts per second, independent of replies")
	flag.IntVar(&o.Queue, "arrival-queue", 256, "bounded waiting requests; overflow is recorded, never hidden")
	flag.Float64Var(&o.HotProbability, "hot-probability", 0, "probability of directing an insert to the hot donor set")
	flag.IntVar(&o.HotDonors, "hot-donors", 0, "hot donors from the largest charity (0 = every donor there)")
	flag.IntVar(&o.Cycles, "cycles", 3, "deterministic grow/correct/delete cycles")
	flag.IntVar(&o.Batch, "batch", 1000, "inserts per growth cycle; half are deleted after correction")
	flag.Int64Var(&o.DBMemoryBytes, "db-memory-bytes", 3221225472, "live database memory limit supplied by runner")
	flag.BoolVar(&o.Targeted, "targeted", false, "also benchmark q08/q12 for largest and smallest charity")
	return o
}

type ExperimentResult struct {
	Settings           ExperimentOptions            `json:"settings"`
	Keys               map[string]int64             `json:"charity_keys"`
	Plans              map[string]string            `json:"plans"`
	Arrival            *ArrivalResult               `json:"arrival,omitempty"`
	Growth             []GrowthPhase                `json:"growth,omitempty"`
	Snapshots          map[string]map[string]string `json:"snapshots"`
	EffectiveIsolation string                       `json:"effective_isolation"`
}

func charityExtremes(ds *Dataset) (int64, int64) {
	counts := make([]int, len(ds.Charities))
	for _, d := range ds.Donations {
		counts[d.CharityID-1]++
	}
	lo, hi := 0, 0
	for i, n := range counts {
		if n < counts[lo] {
			lo = i
		}
		if n > counts[hi] {
			hi = i
		}
	}
	return int64(hi + 1), int64(lo + 1)
}

func configureHot(b *Binder, ds *Dataset, o ExperimentOptions) {
	hi, _ := charityExtremes(ds)
	for _, p := range ds.People {
		if p.CharityID == hi {
			b.HotPeople = append(b.HotPeople, p.ID)
		}
		if o.HotDonors > 0 && len(b.HotPeople) >= o.HotDonors {
			break
		}
	}
	b.HotProbability = o.HotProbability
}

func gate(ctx context.Context, pool *pgxpool.Pool, d Design, ds *Dataset) (*VerifyReport, error) {
	v, err := Verify(ctx, pool, d, ds)
	if err != nil {
		return v, err
	}
	hi, lo := charityExtremes(ds)
	for _, id := range []int64{hi, lo} {
		if err := verifyTarget(ctx, pool, d, ds, id); err != nil {
			v.add(Check{Query: fmt.Sprintf("q08/q12@charity%d", id), Detail: err.Error()})
		} else {
			v.add(Check{Query: fmt.Sprintf("q08/q12@charity%d", id), OK: true})
		}
	}
	if v.Failed > 0 {
		return v, fmt.Errorf("correctness gate failed: %d checks; no valid timings", v.Failed)
	}
	fmt.Printf("  correctness: %d/%d checks passed\n", v.Passed, v.Passed+v.Failed)
	return v, nil
}

// The feed gate checks every returned row and the LIMIT boundary, including
// ties, rather than accepting the right row count with the wrong donations.
func verifyTarget(ctx context.Context, pool *pgxpool.Pool, d Design, ds *Dataset, charity int64) error {
	src, err := readSQL(d.ID, "queries.sql")
	if err != nil {
		return err
	}
	ss, err := ParseCatalog(src)
	if err != nil {
		return err
	}
	q := StmtMap(ss)
	var candidates []Donation
	var sum int64
	byID := map[int64]Donation{}
	for _, v := range ds.Donations {
		if v.CharityID == charity {
			candidates = append(candidates, v)
			sum += v.AmountCents
			byID[v.ID] = v
		}
	}
	var got any
	if err := pool.QueryRow(ctx, q["q08_total_donated_charity"].SQL, charity).Scan(&got); err != nil {
		return err
	}
	if toI64(got) != sum {
		return fmt.Errorf("charity %d sum: got %v want %d", charity, got, sum)
	}
	sort.Slice(candidates, func(i, j int) bool { return candidates[i].DonatedAt.After(candidates[j].DonatedAt) })
	n := len(candidates)
	if n > 50 {
		n = 50
	}
	rows, err := pool.Query(ctx, q["q12_charity_recent_feed"].SQL, charity)
	if err != nil {
		return err
	}
	defer rows.Close()
	seen := map[int64]bool{}
	var previous time.Time
	for rows.Next() {
		var id, amount int64
		var at time.Time
		var name string
		if err := rows.Scan(&id, &amount, &at, &name); err != nil {
			return err
		}
		expected, ok := byID[id]
		if !ok || seen[id] || expected.AmountCents != amount || !expected.DonatedAt.Equal(at) || name != ds.People[expected.PersonID-1].FullName || (!previous.IsZero() && at.After(previous)) || n == 0 || at.Before(candidates[n-1].DonatedAt) {
			return fmt.Errorf("charity %d feed incorrect at donation %d", charity, id)
		}
		seen[id] = true
		previous = at
	}
	if err := rows.Err(); err != nil {
		return err
	}
	if len(seen) != n {
		return fmt.Errorf("charity %d feed length %d want %d", charity, len(seen), n)
	}
	return nil
}

func snapshot(ctx context.Context, pool *pgxpool.Pool, engine string) map[string]string {
	m := map[string]string{}
	for _, name := range []string{"cpu.stat", "memory.current", "memory.peak", "memory.max", "memory.stat", "io.stat"} {
		b, err := os.ReadFile("/sys/fs/cgroup/" + name)
		if err != nil {
			m["client_"+name] = "unavailable: " + err.Error()
		} else {
			m["client_"+name] = strings.TrimSpace(string(b))
		}
	}
	if engine == "postgres" {
		var read, hit int64
		if err := pool.QueryRow(ctx, "SELECT blks_read, blks_hit FROM pg_stat_database WHERE datname=current_database()").Scan(&read, &hit); err != nil {
			m["db_blocks"] = err.Error()
		} else {
			m["db_blks_read"] = fmt.Sprint(read)
			m["db_blks_hit"] = fmt.Sprint(hit)
		}
	}
	return m
}

func executeExperiment(ctx context.Context, pool *pgxpool.Pool, run *Run, d Design, ds *Dataset, opts BenchOpts, writeOps, queries, planPath string, loadConns int) error {
	o := opts.Experiment
	if o.Trial < 1 || o.HistoryMultiplier < 1 || o.HotProbability < 0 || o.HotProbability > 1 || o.HotDonors < 0 {
		return fmt.Errorf("invalid experiment controls")
	}
	if o.Preparation != "analyze" && o.Preparation != "vacuum" {
		return fmt.Errorf("invalid preparation %q", o.Preparation)
	}
	if o.Preparation == "vacuum" && run.Engine != "postgres" {
		return fmt.Errorf("vacuum preparation is PostgreSQL-only")
	}
	hi, lo := charityExtremes(ds)
	x := &ExperimentResult{Settings: o, Keys: map[string]int64{"largest": hi, "smallest": lo}, Plans: map[string]string{}, Snapshots: map[string]map[string]string{}}
	run.Experiment = x
	isoSQL := "SHOW transaction_isolation"
	if run.Engine == "yugabyte" {
		isoSQL = "SHOW yb_effective_transaction_isolation_level"
	}
	if err := pool.QueryRow(ctx, isoSQL).Scan(&x.EffectiveIsolation); err != nil {
		return fmt.Errorf("effective isolation: %w", err)
	}
	var err error
	run.Load, err = SetupAndLoadRetry(ctx, pool, d, ds, run.Engine, loadConns)
	if err != nil {
		return err
	}
	run.Verify, err = gate(ctx, pool, d, ds)
	if err != nil {
		return err
	}
	if o.Preparation == "vacuum" {
		for _, table := range []string{"charity", "person", "donation"} {
			if d.Embedded && table == "donation" {
				continue
			}
			if _, err := pool.Exec(ctx, "VACUUM (ANALYZE) "+table); err != nil {
				return err
			}
		}
	}
	run.Stats, err = CollectStats(ctx, pool, d, run.Engine)
	if err != nil {
		return err
	}
	x.Snapshots["before"] = snapshot(ctx, pool, run.Engine)
	defer func() { x.Snapshots["after"] = snapshot(ctx, pool, run.Engine) }()
	// Full-catalogue plans are the floor, even for a write-only experiment.
	run.Explain, err = ExplainAll(ctx, pool, d, ds, run.Engine)
	if err != nil {
		return err
	}
	for name, plan := range run.Explain {
		if strings.HasPrefix(plan, "ERROR:") {
			return fmt.Errorf("plan %s: %s", name, plan)
		}
	}
	if planPath != "" {
		if err := writePlans(planPath, run, run.Explain); err != nil {
			return err
		}
	}
	b := NewBinder(ds)
	configureHot(b, ds, o)
	opts.Trials = 1 // repeats are separate container invocations with fresh loads
	switch o.Mode {
	case "verify":
		return nil
	case "reads":
		var only []string
		if queries != "" {
			only = strings.Split(queries, ",")
		}
		if err := captureTargetPlans(ctx, pool, run, d, "before", planPath); err != nil {
			return err
		}
		run.Reads, err = BenchmarkReads(ctx, pool, d, b, opts, only)
		if err != nil {
			return err
		}
		if o.Targeted {
			for _, label := range []string{"largest", "smallest"} {
				b.FixedCharity = x.Keys[label]
				rs, err := BenchmarkReads(ctx, pool, d, b, opts, []string{"q08_total_donated_charity", "q12_charity_recent_feed"})
				if err != nil {
					return err
				}
				for _, r := range rs {
					r.Query += "@" + label
					run.Reads = append(run.Reads, r)
				}
			}
		}
		return captureTargetPlans(ctx, pool, run, d, "after", planPath)
	case "writes":
		ops := strings.Split(writeOps, ",")
		if len(ops) != 1 || ops[0] == "" {
			return fmt.Errorf("each experiment write cell must specify exactly one operation")
		}
		run.Writes, err = BenchmarkWrites(ctx, pool, d, b, opts, ops)
		if err != nil {
			return err
		}
		if d.Rollups || d.SumOnly || d.RecentCache {
			run.Audit, err = AuditRollups(ctx, pool, d)
			if err != nil {
				return err
			}
			if len(run.Writes) > 0 {
				run.Writes[0].Audit = run.Audit
			}
		}
	case "arrival":
		x.Arrival, err = benchmarkArrival(ctx, pool, d, ds, b, opts)
		if err != nil {
			return err
		}
		if d.Rollups || d.SumOnly || d.RecentCache {
			run.Audit, err = AuditRollups(ctx, pool, d)
			if err != nil {
				return err
			}
		}
	case "growth":
		x.Growth, err = benchmarkGrowth(ctx, pool, run, d, ds, opts)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown experiment %q", o.Mode)
	}
	return nil
}

func captureTargetPlans(ctx context.Context, pool *pgxpool.Pool, run *Run, d Design, phase, path string) error {
	src, err := readSQL(d.ID, "queries.sql")
	if err != nil {
		return err
	}
	ss, err := ParseCatalog(src)
	if err != nil {
		return err
	}
	var dump strings.Builder
	fmt.Fprintf(&dump, "# Targeted plans — %s\nrun: %s; commit: %s; preparation: %s\n", phase, run.RunID, run.Provenance["repo_commit"], run.Experiment.Settings.Preparation)
	for _, label := range []string{"largest", "smallest"} {
		for _, st := range ss {
			if st.Name != "q08_total_donated_charity" && st.Name != "q12_charity_recent_feed" {
				continue
			}
			var plan string
			for i := 0; i < 2; i++ {
				rows, err := pool.Query(ctx, explainPrefix(run.Engine)+" "+st.SQL, run.Experiment.Keys[label])
				if err != nil {
					return err
				}
				var sb strings.Builder
				for rows.Next() {
					var s string
					if err := rows.Scan(&s); err != nil {
						rows.Close()
						return err
					}
					sb.WriteString(s + "\n")
				}
				rows.Close()
				if err := rows.Err(); err != nil {
					return err
				}
				plan = sb.String()
			}
			key := phase + "/" + label + "/" + st.Name
			run.Experiment.Plans[key] = plan
			fmt.Fprintf(&dump, "\n## %s; charity_id=%d\n\n%s\n\n%s\n", key, run.Experiment.Keys[label], st.SQL, plan)
		}
	}
	if path != "" {
		return os.WriteFile(strings.TrimSuffix(path, filepath.Ext(path))+"-"+phase+".txt", []byte(dump.String()), 0644)
	}
	return nil
}

// Deterministic request seeds preserve the offered parent distribution across
// different service rates. They are deliberately independent of wall-clock time.
func requestRandom(sequence int64) *rand.Rand { return rand.New(rand.NewSource(42 + sequence*7919)) }
