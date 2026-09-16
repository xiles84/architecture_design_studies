// Command bench is study 02's benchmark client.
//
// This file is the entry adapter of the hexagon: it parses flags, opens the
// driven adapters (pgx for the database, files for results) and hands the study's
// use cases nothing but ports. No other file in this package imports a driver.
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"

	"adsplatform/adapters/cgroup"
	"adsplatform/adapters/filestore"
	"adsplatform/adapters/pgxdb"
	"adsplatform/core/catalog"
	"adsplatform/core/inspect"
	"adsplatform/core/provenance"
	"adsplatform/ports"
)

const studyID = "02-ticket-booking"

// Run is the complete record of one (design, topology) cell. Everything needed
// to interpret the numbers travels with them -- including the repository
// version, so a result read a year from now can be traced to the exact code and
// SQL that produced it.
type Run struct {
	Study       string                 `json:"study"`
	RunID       string                 `json:"run_id"`
	Repo        provenance.RepoVersion `json:"repo"`
	StartedAt   time.Time              `json:"started_at"`
	FinishedAt  time.Time              `json:"finished_at"`
	Environment string                 `json:"environment"`
	Engine      string                 `json:"engine"`
	Topology    string                 `json:"topology"`
	EngineInfo  inspect.EngineInfo     `json:"engine_info"`

	DesignID        string `json:"design_id"`
	DesignTitle     string `json:"design_title"`
	DesignFamily    string `json:"design_family"`
	DesignSummary   string `json:"design_summary"`
	Isolation       string `json:"isolation"`
	NegativeControl string `json:"negative_control,omitempty"`

	Scale   string         `json:"scale"`
	Dataset map[string]any `json:"dataset"`
	Options map[string]any `json:"options"`

	Load        *LoadPhases             `json:"load,omitempty"`
	Verify      *VerifyReport           `json:"verify,omitempty"`
	LoadAudit   *Audit                  `json:"load_audit,omitempty"`
	Stats       *DBStats                `json:"stats,omitempty"`
	Reads       []ReadResult            `json:"reads,omitempty"`
	Reports     []ReadResult            `json:"reports,omitempty"`
	ReportCov   map[string]ReportStatus `json:"report_coverage,omitempty"`
	LedgerAudit *LedgerAudit            `json:"ledger_audit,omitempty"`
	Writes      []WriteResult           `json:"writes,omitempty"`
	Races       []RaceResult            `json:"races,omitempty"`
	Holds       []HoldResult            `json:"holds,omitempty"`
	ReloadMS    map[string]float64      `json:"reload_ms,omitempty"`
	// ClientCPU is the benchmark client's own CPU accounting per phase. Throttled
	// periods mean the client, not the database, may have set the tail latency.
	ClientCPU map[string]cgroup.CPUStat `json:"client_cpu,omitempty"`
	Explain   map[string]string         `json:"-"`
	Error     string                    `json:"error,omitempty"`
}

func main() {
	var (
		dsn         = flag.String("dsn", os.Getenv("BENCH_DSN"), "PostgreSQL-protocol connection string")
		engine      = flag.String("engine", "postgres", "postgres | yugabyte")
		topology    = flag.String("topology", "pg-single", "topology id, e.g. pg-single | yb-single | yb-cluster3")
		designID    = flag.String("design", "", "design id (see -cmd list)")
		scale       = flag.String("scale", "small", "tiny | small")
		seed        = flag.Int64("seed", 42, "dataset seed; changing it invalidates comparisons")
		tiers       = flag.String("tiers", "", "restrict the dataset to these event sizes, e.g. 10,100,1000")
		workers     = flag.Int("conns", 8, "concurrent workers for reads and isolated writes")
		streams     = flag.Int("load-conns", 4, "concurrent COPY streams during load")
		duration    = flag.Duration("duration", 5*time.Second, "measured duration per read query / write op trial")
		warmup      = flag.Duration("warmup", 2*time.Second, "discarded warmup per read query / write op")
		trials      = flag.Int("trials", 1, "trials per measurement; reported throughput is the median")
		writeOps    = flag.String("write-ops", "book,cancel,publish", "isolated write ops, or empty to skip")
		phases      = flag.String("phases", "verify,explain,read,write,race,churn,holds", "phases for -cmd full")
		raceBuyers  = flag.Int("race-buyers", 32, "concurrent buyers per event in the sell-out race")
		raceTimeout = flag.Duration("race-timeout", 60*time.Second, "give up on one event's race after this long")
		raceTrials  = flag.Int("race-trials", 1, "repeat the race and churn phases, each on a fresh load")
		tierBudget  = flag.Duration("race-tier-budget", 3*time.Minute, "stop starting new events of a tier after this long (0 = none)")
		raceTiers   = flag.String("race-tiers", "", "restrict races to these event sizes")
		churnPct    = flag.Int("churn-pct", 10, "churn race: percent of successful buyers who cancel at once")
		editorEvery = flag.Duration("editor-interval", 20*time.Millisecond, "organiser edit interval during races (0 = none)")
		holdBuyers  = flag.Int("hold-buyers", 32, "concurrent buyers per event in the holds experiment")
		holdTTL     = flag.Duration("hold-ttl", 250*time.Millisecond, "hold expiry")
		holdThink   = flag.Duration("hold-think-max", 200*time.Millisecond, "max basket time before paying")
		holdPayMin  = flag.Duration("hold-pay-min", 20*time.Millisecond, "min payment time")
		holdPayMax  = flag.Duration("hold-pay-max", 100*time.Millisecond, "max payment time")
		holdAbandon = flag.Int("hold-abandon-pct", 25, "percent of holds abandoned in the basket")
		holdSweep   = flag.Duration("hold-sweep-every", 50*time.Millisecond, "sweeper interval")
		holdScale   = flag.Float64("hold-time-scale", 1, "multiply hold TTL, basket, payment and sweeper timings by this factor")
		holdEvents  = flag.Int("hold-events-per-tier", 10, "events per tier in the holds experiment")
		holdTiers   = flag.String("hold-tiers", "10,100,1000", "event sizes in the holds experiment")
		stmtTimeout = flag.Int("stmt-timeout-ms", 300000, "server-side statement_timeout")
		out         = flag.String("out", "", "path to write the JSON result")
		explainOut  = flag.String("explain-out", "", "path to write the readable plans")
		envName     = flag.String("environment", os.Getenv("BENCH_ENVIRONMENT"), "environment id from docs/environments")
		runID       = flag.String("run-id", "", "run identifier")
		repoCommit  = flag.String("repo-commit", "", "git commit the run was produced from (set by the runner)")
		repoDesc    = flag.String("repo-describe", "", "git describe --tags --always --dirty")
		repoDirty   = flag.Bool("repo-dirty", false, "the working tree had uncommitted changes")
		runTag      = flag.String("run-tag", "", "git tag created for this run")
		cmd         = flag.String("cmd", "full", "full | verify | list | report | digest")
		resultsDir  = flag.String("results", "", "report/digest: results directory of one run")
		reportOut   = flag.String("report-out", "", "report: markdown file to write")
	)
	flag.Parse()

	switch *cmd {
	case "list":
		for _, d := range designs {
			fmt.Printf("%-28s %-26s %-6s %s\n", d.ID, d.Family, d.Isolation, d.Summary)
		}
		return
	case "digest":
		dg, n, err := provenance.Digest(os.DirFS(*resultsDir))
		if err != nil {
			fatal(err)
		}
		fmt.Println(dg, n)
		return
	case "report":
		if *resultsDir == "" || *reportOut == "" {
			fatal(fmt.Errorf("-cmd report needs -results and -report-out"))
		}
		if err := WriteReport(*resultsDir, *reportOut); err != nil {
			fatal(err)
		}
		fmt.Println("wrote", *reportOut)
		return
	}

	if *dsn == "" {
		fatal(fmt.Errorf("-dsn is required"))
	}
	d, err := designByID(*designID)
	if err != nil {
		fatal(err)
	}
	// One factor scales every hold timing together. The late-payment probability
	// depends only on their ratios, so scaling keeps the experiment the same
	// experiment on an engine whose hold itself takes longer than an unscaled
	// TTL (YugabyteDB here: at scale 1 nearly every hold expired before payment).
	for _, p := range []*time.Duration{holdTTL, holdThink, holdPayMin, holdPayMax, holdSweep} {
		*p = time.Duration(float64(*p) * *holdScale)
	}
	s := Settings{
		Workers: *workers, Duration: *duration, Warmup: *warmup, Trials: *trials, Streams: *streams,
		Tiers: parseInts(*tiers), WriteOps: splitList(*writeOps), RaceTiers: parseInts(*raceTiers),
		RaceBuyers: *raceBuyers, RaceTimeout: *raceTimeout, RaceTrials: *raceTrials, RaceTierBudget: *tierBudget, ChurnPct: *churnPct, EditorInterval: *editorEvery,
		HoldBuyers: *holdBuyers, HoldTTL: *holdTTL, HoldThinkMax: *holdThink, HoldPayMin: *holdPayMin,
		HoldPayMax: *holdPayMax, HoldAbandonPct: *holdAbandon, HoldSweepEvery: *holdSweep,
		HoldEventsPerTier: *holdEvents, HoldTiers: parseInts(*holdTiers),
	}
	if *runID == "" {
		*runID = time.Now().UTC().Format("20060102T150405Z")
	}
	run := &Run{
		Study: studyID, RunID: *runID, StartedAt: time.Now().UTC(), Environment: *envName,
		Repo:   provenance.RepoVersion{Commit: *repoCommit, Describe: *repoDesc, Dirty: *repoDirty, Tag: *runTag},
		Engine: *engine, Topology: *topology,
		DesignID: d.ID, DesignTitle: d.Title, DesignFamily: d.Family, DesignSummary: d.Summary,
		Isolation: string(d.Isolation), NegativeControl: d.NegativeControl, Scale: *scale,
		Options: map[string]any{
			"conns": *workers, "load_conns": *streams, "duration": duration.String(), "warmup": warmup.String(),
			"trials": *trials, "seed": *seed, "tiers": *tiers, "write_ops": *writeOps, "phases": *phases,
			"race_buyers": *raceBuyers, "race_timeout": raceTimeout.String(), "race_trials": *raceTrials, "race_tier_budget": tierBudget.String(), "churn_pct": *churnPct,
			"editor_interval": editorEvery.String(), "hold_buyers": *holdBuyers, "hold_ttl": holdTTL.String(),
			"hold_think_max": holdThink.String(), "hold_pay": holdPayMin.String() + "-" + holdPayMax.String(),
			"hold_abandon_pct": *holdAbandon, "hold_sweep_every": holdSweep.String(), "hold_time_scale": *holdScale,
			"hold_events_per_tier": *holdEvents, "hold_tiers": *holdTiers, "statement_timeout_ms": *stmtTimeout,
		},
	}

	// A cell stopped from outside (podman stop, Ctrl+C) still writes its result
	// file, with the interruption recorded as its error: a cell that vanishes
	// without a trace looks like a cell that was never run.
	ctx, stopSignals := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stopSignals()
	// Pool sized over every concurrent user of it -- workers, race buyers, the
	// editor, the sweeper, the monitor -- so the pool is never the bottleneck.
	db, err := pgxdb.Open(ctx, pgxdb.Config{
		DSN: *dsn, MaxConns: max(*workers, *raceBuyers, *holdBuyers) + *streams + 8,
		StatementTimeoutMS: *stmtTimeout, ApplicationName: "ticketbooking-bench",
	})
	if err != nil {
		fatal(err)
	}
	defer db.Close()

	if err := execute(ctx, run, db, d, *cmd, *seed, splitList(*phases), s, *explainOut); err != nil {
		run.Error = err.Error()
		fmt.Fprintln(os.Stderr, "ERROR:", err)
	}
	if ctx.Err() != nil {
		// Phases that swallow per-operation errors may have returned "normally"
		// after the signal; their numbers are partial and must not look complete.
		run.Error = strings.TrimPrefix(run.Error+"; interrupted by signal: results are partial", "; ")
	}
	run.FinishedAt = time.Now().UTC()
	if *out != "" {
		if err := filestore.WriteJSON(*out, run); err != nil {
			fatal(err)
		}
		fmt.Println("  wrote", *out)
	}
	if run.Error != "" {
		os.Exit(1)
	}
}

func execute(ctx context.Context, run *Run, db ports.DB, d Design, cmd string, seed int64, phases []string, s Settings, explainOut string) error {
	if err := waitReady(ctx, db, 120*time.Second); err != nil {
		return err
	}
	run.EngineInfo = inspect.Describe(ctx, db, run.Engine)
	fmt.Printf("  engine: %s\n", firstLine(run.EngineInfo.Version))
	if run.EngineInfo.EffectiveIsolation != "" {
		fmt.Printf("  effective isolation for READ COMMITTED requests: %s\n", run.EngineInfo.EffectiveIsolation)
	}

	ds, err := Generate(run.Scale, seed, s.Tiers)
	if err != nil {
		return err
	}
	run.Dataset = ds.Summary()
	fmt.Printf("  dataset: %d bands, %d events, %d seats, %d sold at load\n",
		len(ds.Bands), len(ds.Events), ds.NextTicketID-1, len(ds.Sold))

	has := func(p string) bool {
		if cmd == "verify" {
			return p == "verify"
		}
		for _, x := range phases {
			if x == p {
				return true
			}
		}
		return false
	}
	run.ReloadMS = map[string]float64{}
	reload := func(why string) (*World, error) {
		t := time.Now()
		ph, err := Load(ctx, db, d, ds, s.Streams)
		if err != nil {
			return nil, fmt.Errorf("load before %s: %w", why, err)
		}
		if run.Load == nil {
			run.Load = ph
		} else {
			run.ReloadMS[why] = msSince(t)
		}
		return NewWorld(ds), nil
	}

	fmt.Printf("  loading %s ...\n", d.ID)
	world, err := reload("initial")
	if err != nil {
		return err
	}
	fmt.Printf("  loaded in %.1fs (copy %.1fs, index %.1fs)\n", run.Load.TotalMS/1000, run.Load.CopyMS/1000, run.Load.IndexMS/1000)

	// Correctness before speed.
	if has("verify") {
		fmt.Println("  verifying answers against ground truth ...")
		vr, err := Verify(ctx, db, d, ds)
		if err != nil {
			return err
		}
		run.Verify = vr
		for _, c := range vr.Checks {
			if !c.OK {
				fmt.Printf("    FAIL %-26s %s: expected %s, got %s\n", c.Query, c.Key, c.Expect, c.Got)
			}
		}
		au, err := RunAudit(ctx, db, d, world, "load")
		if err != nil {
			return err
		}
		run.LoadAudit = au
		fmt.Printf("    %d/%d checks passed; audit on load: %s\n", vr.Passed, vr.Passed+vr.Failed, au)
		la, err := RunLedgerAudit(ctx, db, d, "load")
		if err != nil {
			return err
		}
		if la != nil {
			run.LedgerAudit = la
			fmt.Printf("    ledger audit on load: %s\n", la)
		}
		if vr.Failed > 0 || au.Violations() > 0 || (la != nil && la.Mismatches > 0) {
			return fmt.Errorf("design %s failed the correctness gate; refusing to report timings", d.ID)
		}
		if st, err := CollectStats(ctx, db, run.Engine); err == nil {
			run.Stats = st
		}
	}
	if cmd == "verify" {
		return nil
	}

	bk, err := NewBooker(db, d, world)
	if err != nil {
		return err
	}

	if has("explain") {
		fmt.Println("  capturing plans ...")
		pl, err := ExplainAll(ctx, db, d, ds, world, run.Engine)
		if err != nil {
			return err
		}
		run.Explain = pl
		if explainOut != "" {
			if err := writePlans(explainOut, run, d, pl); err != nil {
				return err
			}
		}
	}

	if has("read") {
		fmt.Println("  read benchmark:")
		stop := probe(run, "read")
		rs, err := BenchmarkReads(ctx, db, d, ds, s)
		stop()
		if err != nil {
			return err
		}
		run.Reads = rs

		run.ReportCov = ReportCoverage(d)
		fmt.Println("  operational reports:")
		stop = probe(run, "reports")
		reps, err := BenchmarkReports(ctx, db, d, ds, s)
		stop()
		if err != nil {
			return err
		}
		run.Reports = reps
	}

	// Every isolated write op gets its own fresh load, audited on the state it
	// produced. The initial load is still pristine for the first one.
	fresh := true
	nextWorld := func(why string) error {
		if fresh {
			fresh = false
			return nil
		}
		w, err := reload(why)
		if err != nil {
			return err
		}
		world = w
		bk, err = NewBooker(db, d, world)
		return err
	}
	// Reads and plans did not write; only EXPLAIN ANALYZE inside rolled-back
	// transactions touched the tables. The first write op may use this load.

	if has("write") {
		fmt.Println("  write benchmark (fresh load per op):")
		for _, op := range s.WriteOps {
			if err := nextWorld(op); err != nil {
				return err
			}
			var results []WriteResult
			stop := probe(run, "write:"+op)
			switch op {
			case "book":
				results = append(results, BenchmarkBookSpread(ctx, bk, world, s))
			case "cancel":
				results = append(results, BenchmarkCancel(ctx, bk, ds, s))
			case "publish":
				// Smallest first, so the large tiers' bulk inserts cannot bloat
				// the table the small tiers are measured against.
				for _, tier := range allTiers {
					if len(s.Tiers) > 0 && !containsInt(s.Tiers, tier) {
						continue
					}
					results = append(results, BenchmarkPublish(ctx, bk, world, tier, s))
				}
			default:
				return fmt.Errorf("unknown write op %q", op)
			}
			stop()
			au, err := RunAudit(ctx, db, d, world, "write:"+op)
			if err != nil {
				return err
			}
			fmt.Printf("      audit after %s: %s\n", op, au)
			la, err := RunLedgerAudit(ctx, db, d, "write:"+op)
			if err != nil {
				return err
			}
			if la != nil {
				fmt.Printf("      ledger audit after %s: %s\n", op, la)
			}
			for i := range results {
				results[i].Audit = au
				results[i].LedgerAudit = la
			}
			run.Writes = append(run.Writes, results...)
		}
	}

	// Races repeat on a fresh load per trial. The large tiers are one event each,
	// so without repetition a 100 000-seat race is a single sample.
	for trial := 1; trial <= max(s.RaceTrials, 1); trial++ {
		for _, mode := range []string{"race", "churn"} {
			if !has(mode) {
				continue
			}
			if err := nextWorld(fmt.Sprintf("%s#%d", mode, trial)); err != nil {
				return err
			}
			fmt.Printf("  %s trial %d/%d (%d buyers per event, fresh load):\n", mode, trial, max(s.RaceTrials, 1), s.RaceBuyers)
			stop := probe(run, fmt.Sprintf("%s#%d", mode, trial))
			rr, err := RunRaces(ctx, db, bk, d, ds, world, mode, s)
			stop()
			if err != nil {
				return err
			}
			for i := range rr {
				rr[i].Trial = trial
			}
			run.Races = append(run.Races, rr...)
		}
	}

	if has("holds") && d.Holds {
		if err := nextWorld("holds"); err != nil {
			return err
		}
		fmt.Printf("  holds (%d buyers per event, TTL %s, fresh load):\n", s.HoldBuyers, s.HoldTTL)
		stop := probe(run, "holds")
		hr, err := RunHolds(ctx, db, bk, d, ds, world, s)
		stop()
		if err != nil {
			return err
		}
		run.Holds = hr
	}
	return nil
}

// probe records how much the client container was CPU-throttled during a phase.
// The returned function closes the measurement.
func probe(run *Run, phase string) func() {
	before, ok := cgroup.Read()
	if !ok {
		return func() {}
	}
	return func() {
		after, ok := cgroup.Read()
		if !ok {
			return
		}
		if run.ClientCPU == nil {
			run.ClientCPU = map[string]cgroup.CPUStat{}
		}
		run.ClientCPU[phase] = after.Sub(before)
	}
}

func writePlans(path string, run *Run, d Design, plans map[string]string) error {
	var sb strings.Builder
	fmt.Fprintf(&sb, "# Query plans — study %s\n\n", studyID)
	fmt.Fprintf(&sb, "run_id:      %s\nrepo:        %s (dirty=%v)\nenvironment: %s\nengine:      %s (%s)\nversion:     %s\n",
		run.RunID, run.Repo.Describe, run.Repo.Dirty, run.Environment, run.Engine, run.Topology, firstLine(run.EngineInfo.Version))
	fmt.Fprintf(&sb, "design:      %s (%s)\ncaptured:    %s\n", d.ID, d.Title, inspect.ExplainPrefix(run.Engine))
	fmt.Fprintf(&sb, "keys:        the largest partially sold catalogue event, its busiest customer and one of its sold\n")
	fmt.Fprintf(&sb, "             tickets. Write statements run inside transactions that are rolled back.\n\n")
	for _, file := range []string{"queries.sql", "writes.sql"} {
		stmts, err := mustStmts(d.ID, file)
		if err != nil {
			return err
		}
		for _, st := range stmts {
			writePlan(&sb, st, plans[st.Name])
		}
	}
	return filestore.WriteText(path, sb.String())
}

func writePlan(sb *strings.Builder, st catalog.Stmt, plan string) {
	fmt.Fprintf(sb, "%s\n## %s\n%s\n\n", strings.Repeat("=", 78), st.Name, strings.Repeat("=", 78))
	for _, l := range strings.Split(st.Doc, "\n") {
		if l != "" {
			fmt.Fprintf(sb, "# %s\n", l)
		}
	}
	fmt.Fprintf(sb, "\n%s\n\n%s\n\n", strings.TrimSpace(st.SQL), plan)
}

func waitReady(ctx context.Context, db ports.DB, limit time.Duration) error {
	deadline := time.Now().Add(limit)
	var last error
	for time.Now().Before(deadline) {
		c, cancel := context.WithTimeout(ctx, 5*time.Second)
		_, err := db.Exec(c, "SELECT 1")
		cancel()
		if err == nil {
			return nil
		}
		last = err
		time.Sleep(2 * time.Second)
	}
	return fmt.Errorf("database not ready after %s: %w", limit, last)
}

func parseInts(s string) []int {
	var out []int
	for _, p := range splitList(s) {
		if v, err := strconv.Atoi(p); err == nil {
			out = append(out, v)
		}
	}
	sort.Ints(out)
	return out
}

func splitList(s string) []string {
	var out []string
	for _, p := range strings.Split(s, ",") {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}

func firstLine(s string) string {
	if i := strings.IndexByte(s, '\n'); i >= 0 {
		return s[:i]
	}
	return s
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, "fatal:", err)
	os.Exit(1)
}
