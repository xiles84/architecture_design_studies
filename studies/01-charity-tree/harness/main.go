package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Run is the complete record of one (design, engine, topology) measurement.
// Everything needed to interpret or reproduce the numbers travels with them:
// the engine banner, the dataset fingerprint, the benchmark settings and the
// plans. A results file that needs external context to read is a results file
// that will be misread later.
type Run struct {
	RunID       string            `json:"run_id"`
	StartedAt   time.Time         `json:"started_at"`
	FinishedAt  time.Time         `json:"finished_at"`
	Environment string            `json:"environment"`
	Engine      string            `json:"engine"`
	Topology    string            `json:"topology"`
	EngineVer   string            `json:"engine_version"`
	DesignID    string            `json:"design_id"`
	DesignTitle string            `json:"design_title"`
	DesignNote  string            `json:"design_summary"`
	Scale       string            `json:"scale"`
	Dataset     map[string]any    `json:"dataset"`
	Options     map[string]any    `json:"options"`
	Load        *LoadPhases       `json:"load,omitempty"`
	Verify      *VerifyReport     `json:"verify,omitempty"`
	Stats       *DBStats          `json:"stats,omitempty"`
	Reads       []QueryResult     `json:"reads,omitempty"`
	Writes      []WriteResult     `json:"writes,omitempty"`
	Audit       *RollupAudit      `json:"rollup_audit,omitempty"`
	Mixed       []*MixedResult    `json:"mixed,omitempty"`
	Explain     map[string]string `json:"explain,omitempty"`
	Error       string            `json:"error,omitempty"`
}

func main() {
	var (
		dsn          = flag.String("dsn", envOr("BENCH_DSN", ""), "PostgreSQL-protocol connection string")
		engine       = flag.String("engine", envOr("BENCH_ENGINE", "postgres"), "postgres | yugabyte")
		topology     = flag.String("topology", envOr("BENCH_TOPOLOGY", "single"), "single | cluster3")
		designID     = flag.String("design", "", "design id, e.g. d3_flattened_fk")
		scale        = flag.String("scale", "small", "tiny | small | medium | large")
		seed         = flag.Int64("seed", 42, "dataset seed; changing it invalidates comparisons")
		conns        = flag.Int("conns", 8, "concurrent benchmark connections")
		loadConns    = flag.Int("load-conns", 4, "concurrent COPY streams during load")
		duration     = flag.Duration("duration", 10*time.Second, "measured duration per query/op")
		warmup       = flag.Duration("warmup", 2*time.Second, "unrecorded warmup per query/op")
		maxOps       = flag.Int64("max-ops", 0, "stop a query early after this many ops (0 = no cap)")
		trials       = flag.Int("trials", 1, "repeat each measured phase N times; reported throughput is the median")
		mixSplits    = flag.String("mix", "8:0,6:2,4:4", "mixed workload read:write worker splits (experiment D)")
		readMixName  = flag.String("read-mix", "dashboard", "mixed workload read mix: dashboard | portal")
		maxPerPerson = flag.Int("max-per-person", 0, "cap donations per donor (0 = unbounded); total donations are held constant")
		charities    = flag.Int("charities", 0, "override charity count (0 = scale default of 10); people and donations are held constant")
		strategy     = flag.String("strategy", "blind", "D5 concurrency control: blind | forupdate | optimistic")
		isolation    = flag.String("isolation", "rc", "transaction isolation: rc | rr | ser")
		writeOps     = flag.String("write-ops", "insert,update,delete", "comma-separated write ops, or empty to skip")
		queries      = flag.String("queries", "", "comma-separated query ids to restrict reads to")
		out          = flag.String("out", "", "path to write the JSON result")
		explainTo    = flag.String("explain-out", "", "path to write the human-readable plan dump")
		envName      = flag.String("environment", envOr("BENCH_ENVIRONMENT", "unknown"), "environment id from docs/environments")
		runID        = flag.String("run-id", "", "run identifier; defaults to a timestamp")
		cmd          = flag.String("cmd", "full", "full | load | verify | read | write | mixed | explain | stats | report | list")
		resultsDir   = flag.String("results", "", "report: directory of result JSON files to summarise")
		reportOut    = flag.String("report-out", "", "report: markdown file to write")
		compareDir   = flag.String("compare", "", "report-compare: the alternative-regime results directory")
		labelA       = flag.String("label-a", "", "report-compare: name for the -results regime")
		labelB       = flag.String("label-b", "", "report-compare: name for the -compare regime")
		stmtTimeout  = flag.Int("stmt-timeout-ms", 120000, "server-side statement_timeout")
	)
	flag.Parse()

	switch *cmd {
	case "digest":
		dg, n, err := ResultsDigest(*resultsDir)
		if err != nil {
			fatal(err)
		}
		fmt.Println(dg, n)
		return
	case "report", "report-mixed", "report-concurrency", "report-compare":
		if *resultsDir == "" || *reportOut == "" {
			fatal(fmt.Errorf("-cmd %s needs -results <dir> and -report-out <file>", *cmd))
		}
		var err error
		switch *cmd {
		case "report":
			err = writeReport(*resultsDir, *reportOut)
		case "report-mixed":
			err = writeMixedReport(*resultsDir, *reportOut)
		case "report-concurrency":
			err = writeConcurrencyReport(*resultsDir, *reportOut)
		case "report-compare":
			err = writeCompareReport(*resultsDir, *compareDir, *labelA, *labelB, *reportOut)
		}
		if err != nil {
			fatal(err)
		}
		fmt.Println("wrote", *reportOut)
		return
	}
	if *cmd == "list" {
		for _, d := range designs {
			fmt.Printf("%-24s %-20s engines=%-22s %s\n", d.ID, d.Title, strings.Join(d.Engines, ","), d.Summary)
		}
		return
	}
	if *dsn == "" {
		fatal(fmt.Errorf("-dsn is required"))
	}
	if *runID == "" {
		*runID = time.Now().UTC().Format("20060102T150405Z")
	}

	d, err := designByID(*designID)
	if err != nil {
		fatal(err)
	}
	if !d.supports(*engine) {
		fatal(fmt.Errorf("design %s does not support engine %s (supported: %s)",
			d.ID, *engine, strings.Join(d.Engines, ",")))
	}

	ctx := context.Background()
	run := &Run{
		RunID:       *runID,
		StartedAt:   time.Now().UTC(),
		Environment: *envName,
		Engine:      *engine,
		Topology:    *topology,
		DesignID:    d.ID,
		DesignTitle: d.Title,
		DesignNote:  d.Summary,
		Scale:       *scale,
		Options: map[string]any{
			"conns": *conns, "load_conns": *loadConns,
			"duration": duration.String(), "warmup": warmup.String(),
			"max_ops": *maxOps, "strategy": *strategy, "isolation": *isolation,
			"seed": *seed, "statement_timeout_ms": *stmtTimeout,
			"max_per_person": *maxPerPerson, "read_mix": *readMixName, "charities": *charities,
		},
	}

	opts := BenchOpts{
		Conns: *conns, Duration: *duration, Warmup: *warmup, MaxOps: *maxOps,
		Strategy: *strategy, Isolation: *isolation,
		StmtTimeoutMS: *stmtTimeout, Trials: *trials,
		MaxPerPerson: *maxPerPerson, ReadMix: *readMixName, Charities: *charities,
	}
	run.Options["trials"] = *trials

	run.Options["mix"] = *mixSplits

	if err := execute(ctx, run, d, *dsn, *cmd, *scale, *seed, *loadConns,
		opts, *writeOps, *queries, *explainTo, *mixSplits); err != nil {
		run.Error = err.Error()
		fmt.Fprintln(os.Stderr, "ERROR:", err)
	}
	run.FinishedAt = time.Now().UTC()

	if *out != "" {
		if err := writeJSON(*out, run); err != nil {
			fatal(err)
		}
		fmt.Println("  wrote", *out)
	}
	if run.Error != "" {
		os.Exit(1)
	}
}

func execute(ctx context.Context, run *Run, d Design, dsn, cmd, scale string,
	seed int64, loadConns int, opts BenchOpts, writeOps, queries, explainTo, mixSplits string) error {

	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return err
	}
	// Leave headroom over the worker count: the pool must never be the bottleneck
	// we accidentally measure.
	cfg.MaxConns = int32(opts.Conns + loadConns + 4)
	cfg.MinConns = 1
	cfg.ConnConfig.RuntimeParams["statement_timeout"] = fmt.Sprint(opts.StmtTimeoutMS)
	cfg.ConnConfig.RuntimeParams["application_name"] = "charitytree-bench"

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return err
	}
	defer pool.Close()

	if err := waitReady(ctx, pool, 120*time.Second); err != nil {
		return err
	}
	run.EngineVer = EngineVersion(ctx, pool)
	fmt.Printf("  engine: %s\n", firstLine(run.EngineVer))

	fmt.Printf("  generating dataset scale=%s seed=%d ...\n", scale, seed)
	ds, err := GenerateProfile(scale, seed, Profile{MaxPerPerson: opts.MaxPerPerson, Charities: opts.Charities})
	if err != nil {
		return err
	}
	run.Dataset = ds.Summary()
	fmt.Printf("  dataset: %d charities, %d people, %d donations\n",
		len(ds.Charities), len(ds.People), len(ds.Donations))

	binder := NewBinder(ds)

	// verify must load its own data: it compares against the pristine generated
	// dataset, and any previous run's write benchmark will have mutated whatever
	// is already in the database.
	// "mixed" loads once per split instead, so that a later split never inherits
	// the bloat an earlier one created.
	needLoad := cmd == "full" || cmd == "load" || cmd == "verify"
	if needLoad {
		fmt.Printf("  loading %s ...\n", d.ID)
		ph, err := SetupAndLoadRetry(ctx, pool, d, ds, run.Engine, loadConns)
		if err != nil {
			return err
		}
		run.Load = ph
		fmt.Printf("  loaded in %.1fs (copy %.1fs, rollup %.1fs, index %.1fs)\n",
			ph.TotalMS/1000, ph.CopyMS/1000, ph.RollupMS/1000, ph.IndexMS/1000)
	}

	// Correctness before speed. Verification runs against freshly loaded data and
	// before any write benchmark mutates it, and a failure aborts the run: a
	// design that answers the question wrongly has no meaningful throughput.
	if cmd == "full" || cmd == "verify" {
		fmt.Println("  verifying answers against ground truth ...")
		vr, err := Verify(ctx, pool, d, ds)
		if err != nil {
			return err
		}
		run.Verify = vr
		for _, c := range vr.Checks {
			if !c.OK {
				fmt.Printf("    FAIL %-30s expected %s, got %s %s\n", c.Query, c.Expect, c.Got, c.Detail)
			}
		}
		fmt.Printf("    %d/%d checks passed\n", vr.Passed, vr.Passed+vr.Failed)
		if vr.Failed > 0 {
			return fmt.Errorf("%d correctness checks failed for design %s; refusing to report timings", vr.Failed, d.ID)
		}
	}

	if cmd == "full" || cmd == "stats" {
		st, err := CollectStats(ctx, pool, d, run.Engine)
		if err != nil {
			return err
		}
		run.Stats = st
	}

	if cmd == "full" || cmd == "explain" {
		fmt.Println("  capturing plans ...")
		pl, err := ExplainAll(ctx, pool, d, ds, run.Engine)
		if err != nil {
			return err
		}
		run.Explain = pl
		if explainTo != "" {
			if err := writePlans(explainTo, run, pl); err != nil {
				return err
			}
			fmt.Println("  wrote", explainTo)
		}
	}

	if cmd == "full" || cmd == "read" {
		fmt.Println("  read benchmark:")
		var only []string
		if queries != "" {
			only = strings.Split(queries, ",")
		}
		rs, err := BenchmarkReads(ctx, pool, d, binder, opts, only)
		if err != nil {
			return err
		}
		run.Reads = rs
	}

	// Experiment D: readers and writers sharing the database.
	//
	// The first split is conventionally "N:0" -- all readers, no writers. That is
	// the baseline every other split is measured against, and taking it here
	// rather than from the isolated read benchmark matters: it uses the same
	// weighted mix, the same data and the same process, so the only thing that
	// changes between splits is the presence of writers.
	if cmd == "mixed" {
		for _, sp := range strings.Split(mixSplits, ",") {
			rw, ww, err := parseSplit(sp)
			if err != nil {
				return err
			}
			fmt.Printf("  mixed workload %s (reload first) ...\n", sp)
			if _, err := SetupAndLoadRetry(ctx, pool, d, ds, run.Engine, loadConns); err != nil {
				return err
			}
			res, err := BenchmarkMixed(ctx, pool, d, binder, opts, rw, ww)
			if err != nil {
				return err
			}
			if ww > 0 && (d.Rollups || d.RecentCache) {
				au, err := AuditRollups(ctx, pool, d)
				if err != nil {
					return err
				}
				res.Audit = au
				bad := au.PersonMismatches + au.CharityMismatches + au.CacheMismatches
				if bad == 0 {
					fmt.Printf("      audit: consistent\n")
				} else {
					fmt.Printf("      audit: INCONSISTENT — %d rollup and %d cache rows disagree\n",
						au.PersonMismatches+au.CharityMismatches, au.CacheMismatches)
				}
			}
			run.Mixed = append(run.Mixed, res)
		}
	}

	if (cmd == "full" || cmd == "write") && strings.TrimSpace(writeOps) != "" {
		fmt.Println("  write benchmark:")
		// Each write operation gets a freshly loaded database.
		//
		// Running insert, then update, then delete against the same table makes
		// every later phase inherit the earlier phases' growth, dead tuples and
		// index bloat -- by DIFFERENT amounts per design, because each design's
		// earlier phases ran at different speeds. This is not hypothetical: the
		// survey measured D6 deleting 28 donations/s after insert and update had run,
		// and the same delete on a fresh load measured ~830/s -- a 30x distortion
		// in a headline number, caused entirely by phase order. The reload happens
		// before every op except the first (which follows the initial load).
		ops := strings.Split(writeOps, ",")
		for i, op := range ops {
			if i > 0 && cmd == "full" {
				if _, err := SetupAndLoadRetry(ctx, pool, d, ds, run.Engine, loadConns); err != nil {
					return fmt.Errorf("reload before %s: %w", op, err)
				}
			}
			ws, err := BenchmarkWrites(ctx, pool, d, binder, opts, []string{op})
			if err != nil {
				return err
			}

			// Throughput means nothing if the writes lost data, so each isolated
			// op is audited on the state IT produced. Auditing only at the end
			// would check nothing but the last op, now that ops no longer share
			// a table.
			if (d.Rollups || d.RecentCache) && len(ws) > 0 && ws[0].Skipped == "" {
				au, err := AuditRollups(ctx, pool, d)
				if err != nil {
					return err
				}
				ws[0].Audit = au
				reportAudit(op, d, au)
				// run.Audit keeps the worst result, so a single inconsistent op
				// cannot be hidden behind a later consistent one.
				if run.Audit == nil || auditBad(au) > auditBad(run.Audit) {
					run.Audit = au
				}
			}
			run.Writes = append(run.Writes, ws...)
		}
	}
	return nil
}

func auditBad(a *RollupAudit) int64 {
	if a == nil {
		return 0
	}
	return a.PersonMismatches + a.CharityMismatches + a.CacheMismatches
}

func reportAudit(op string, d Design, au *RollupAudit) {
	if d.Rollups {
		if au.PersonMismatches+au.CharityMismatches == 0 {
			fmt.Printf("    audit after %s: rollups consistent (%d person + %d charity rows)\n",
				op, au.PersonRows, au.CharityRows)
		} else {
			fmt.Printf("    audit after %s: rollups INCONSISTENT — %d person and %d charity rows "+
				"disagree with their donations; drift %d cents\n",
				op, au.PersonMismatches, au.CharityMismatches, au.DriftCents)
		}
	}
	if d.RecentCache {
		if au.CacheMismatches == 0 {
			fmt.Printf("    audit after %s: embedded cache consistent (%d rows)\n", op, au.CacheChecked)
		} else {
			fmt.Printf("    audit after %s: embedded cache STALE — %d of %d person rows wrong\n",
				op, au.CacheMismatches, au.CacheChecked)
		}
	}
}

// parseSplit reads a "readers:writers" worker split, e.g. "6:2".
func parseSplit(s string) (int, int, error) {
	var r, w int
	if _, err := fmt.Sscanf(strings.TrimSpace(s), "%d:%d", &r, &w); err != nil {
		return 0, 0, fmt.Errorf("bad mix split %q, want READERS:WRITERS (e.g. 6:2): %w", s, err)
	}
	if r < 0 || w < 0 || r+w == 0 {
		return 0, 0, fmt.Errorf("bad mix split %q: needs at least one worker", s)
	}
	return r, w, nil
}

func waitReady(ctx context.Context, pool *pgxpool.Pool, limit time.Duration) error {
	deadline := time.Now().Add(limit)
	var last error
	for time.Now().Before(deadline) {
		c, cancel := context.WithTimeout(ctx, 5*time.Second)
		_, err := pool.Exec(c, "SELECT 1")
		cancel()
		if err == nil {
			return nil
		}
		last = err
		time.Sleep(2 * time.Second)
	}
	return fmt.Errorf("database not ready after %s: %w", limit, last)
}

func writeJSON(path string, v any) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(b, '\n'), 0o644)
}

// writePlans emits the plan dump as readable text. JSON is for machines to
// aggregate; a plan is for a person to read, and burying it in escaped strings
// makes the single most informative artefact of the run unusable.
func writePlans(path string, run *Run, plans map[string]string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "# Query plans\n\n")
	fmt.Fprintf(&sb, "run_id:      %s\n", run.RunID)
	fmt.Fprintf(&sb, "environment: %s\n", run.Environment)
	fmt.Fprintf(&sb, "engine:      %s (%s)\n", run.Engine, run.Topology)
	fmt.Fprintf(&sb, "version:     %s\n", firstLine(run.EngineVer))
	fmt.Fprintf(&sb, "design:      %s (%s)\n", run.DesignID, run.DesignTitle)
	fmt.Fprintf(&sb, "scale:       %s -- %v people, %v donations\n",
		run.Scale, run.Dataset["people"], run.Dataset["donations"])
	fmt.Fprintf(&sb, "captured:    %s\n", explainPrefix(run.Engine))
	fmt.Fprintf(&sb, "keys used:   charity_id=1 (the largest charity, 30%% of donors); person_id and\n")
	fmt.Fprintf(&sb, "             donation_id belong to that charity's busiest donor, so these plans\n")
	fmt.Fprintf(&sb, "             show real work rather than a donor who happens to have three rows\n\n")

	src, err := readSQL(run.DesignID, "queries.sql")
	if err != nil {
		return err
	}
	stmts, err := ParseCatalog(src)
	if err != nil {
		return err
	}
	for _, st := range stmts {
		fmt.Fprintf(&sb, "%s\n## %s\n%s\n\n", strings.Repeat("=", 78), st.Name, strings.Repeat("=", 78))
		if st.Doc != "" {
			for _, l := range strings.Split(st.Doc, "\n") {
				fmt.Fprintf(&sb, "# %s\n", l)
			}
			sb.WriteString("\n")
		}
		fmt.Fprintf(&sb, "%s\n\n%s\n\n", strings.TrimSpace(st.SQL), plans[st.Name])
	}
	return os.WriteFile(path, []byte(sb.String()), 0o644)
}

func firstLine(s string) string {
	if i := strings.IndexByte(s, '\n'); i >= 0 {
		return s[:i]
	}
	return s
}

func envOr(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, "fatal:", err)
	os.Exit(1)
}
