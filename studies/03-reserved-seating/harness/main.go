// Command bench is study 03's benchmark client.
//
// This file is the entry adapter of the hexagon: it parses flags, opens the
// driven adapters (pgx for the database, files for results) and hands the study's
// use cases nothing but ports. No other file in this package imports a driver.
//
// Adapted from studies/02-ticket-booking/harness/main.go at 3c0c3aa.
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

const studyID = "03-reserved-seating"

// Run is the complete record of one (design, topology) cell.
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
	// ClientDBClockOffsetMS is the database clock minus the client's (median of 20).
	ClientDBClockOffsetMS float64 `json:"client_db_clock_offset_ms"`
	ConnectionNodes       int     `json:"connection_nodes"`

	DesignID        string   `json:"design_id"`
	DesignTitle     string   `json:"design_title"`
	DesignFamily    string   `json:"design_family"`
	DesignSummary   string   `json:"design_summary"`
	Isolation       string   `json:"isolation"`
	NegativeControl string   `json:"negative_control,omitempty"`
	ControlFires    []string `json:"control_fires,omitempty"`

	Scale   string         `json:"scale"`
	Dataset map[string]any `json:"dataset"`
	Options map[string]any `json:"options"`

	Load      *LoadPhases        `json:"load,omitempty"`
	Verify    *VerifyReport      `json:"verify,omitempty"`
	LoadAudit *Audit             `json:"load_audit,omitempty"`
	PreSweep  int                `json:"pre_verify_swept_seats,omitempty"`
	Stats     *DBStats           `json:"stats,omitempty"`
	Reads     []ReadResult       `json:"reads,omitempty"`
	Writes    []WriteResult      `json:"writes,omitempty"`
	Races     []RaceResult       `json:"races,omitempty"`
	Lifecycle []LifecycleResult  `json:"lifecycle,omitempty"`
	ReloadMS  map[string]float64 `json:"reload_ms,omitempty"`
	// ClientCPU is the benchmark client's own CPU accounting per phase.
	ClientCPU map[string]cgroup.CPUStat `json:"client_cpu,omitempty"`
	Explain   map[string]string         `json:"-"`
	Error     string                    `json:"error,omitempty"`
}

func main() {
	var (
		dsn         = flag.String("dsn", os.Getenv("BENCH_DSN"), "PostgreSQL-protocol connection string; a comma-separated list spreads connections over nodes")
		engine      = flag.String("engine", "postgres", "postgres | yugabyte")
		topology    = flag.String("topology", "pg-single", "pg-single | yb-single | yb-cluster3")
		designID    = flag.String("design", "", "design id (see -cmd list)")
		scale       = flag.String("scale", "small", "tiny | small")
		seed        = flag.Int64("seed", 42, "dataset seed; changing it invalidates comparisons")
		tiers       = flag.String("tiers", "", "restrict the dataset to these event sizes")
		workers     = flag.Int("conns", 8, "concurrent workers for reads and isolated writes")
		streams     = flag.Int("load-conns", 4, "concurrent COPY streams during load")
		duration    = flag.Duration("duration", 5*time.Second, "measured duration per read / write trial")
		warmup      = flag.Duration("warmup", 2*time.Second, "discarded warmup per read / write")
		trials      = flag.Int("trials", 1, "trials per read / write measurement; throughput is the median")
		writeOps    = flag.String("write-ops", "hold,release,cancel,publish", "isolated write ops, or empty to skip")
		phases      = flag.String("phases", "verify,explain,read,write,race,lifecycle", "phases for -cmd full")
		holdTTL     = flag.Duration("hold-ttl", 40*time.Minute, "hold TTL in the race and isolated writes (real time)")
		maxConf     = flag.Int("max-conflicts", 20, "conflicts after which a customer gives up")
		raceBuyers  = flag.Int("race-buyers", 32, "concurrent buyers per event in the race")
		raceTimeout = flag.Duration("race-timeout", 60*time.Second, "stop starting new holds in one event's race after this long")
		raceTrials  = flag.Int("race-trials", 1, "repeat the race, each on a fresh load")
		tierBudget  = flag.Duration("race-tier-budget", 3*time.Minute, "stop starting new events of a tier after this long (0 = none)")
		raceTiers   = flag.String("race-tiers", "", "restrict races to these event sizes")
		deferPct    = flag.Int("race-defer-pct", 10, "percent of granted holds kept until the crowd has gone, then confirmed")
		humanMinute = flag.Duration("human-minute", 50*time.Millisecond, "wall-clock length of one human minute in the lifecycle")
		lcBuyers    = flag.Int("lifecycle-buyers", 32, "concurrent buyers per lifecycle event")
		lcTiers     = flag.String("lifecycle-tiers", "10,100,1000", "event sizes in the lifecycle")
		lcEvents    = flag.Int("lifecycle-events-per-tier", 6, "lifecycle events per tier")
		lcTrials    = flag.Int("lifecycle-trials", 1, "repeat the lifecycle, each on a fresh load")
		lcTTL       = flag.Duration("lc-ttl", 40*time.Minute, "lifecycle hold TTL (human time)")
		lcAbandon   = flag.Int("lc-abandon-pct", 25, "percent of holds abandoned")
		lcExplicit  = flag.Int("lc-abandon-explicit-pct", 40, "percent of abandoned holds released explicitly")
		lcWindow    = flag.Duration("lc-payment-window", 10*time.Minute, "K1 payment window (human time)")
		lcGuard     = flag.Duration("lc-guard", 2*time.Minute, "guard margin G: rejections this far or more before expiry are violations (human time)")
		lcSweep     = flag.Duration("lc-sweep-every", time.Minute, "sweeper interval (human time)")
		lcOutFrom   = flag.Duration("lc-outage-from", 45*time.Minute, "sweeper outage start, first event of each tier (human time)")
		lcOutTo     = flag.Duration("lc-outage-to", 85*time.Minute, "sweeper outage end and probe (human time; 0 disables)")
		lcSkew      = flag.Duration("lc-e0-skew", 10*time.Minute, "E0: clock skew of the fast application node (human time; real time in the race)")
		lcTimeout   = flag.Duration("lc-timeout", 240*time.Minute, "lifecycle event timeout (human time)")
		probeCap    = flag.Int("probe-cap", 200, "free seats probed per event by the leak probe, besides every expired-hold seat")
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
			fmt.Printf("%-34s %-12s %-4s %s\n", d.ID, d.Family, d.Isolation, d.Summary)
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
	if d.YBOnly && *engine != "yugabyte" {
		fatal(fmt.Errorf("design %s is YugabyteDB-only", d.ID))
	}
	s := Settings{
		Workers: *workers, Duration: *duration, Warmup: *warmup, Trials: *trials, Streams: *streams,
		Tiers: parseInts(*tiers), WriteOps: splitList(*writeOps), HoldTTL: *holdTTL, MaxConflicts: *maxConf,
		RaceBuyers: *raceBuyers, RaceTimeout: *raceTimeout, RaceTrials: *raceTrials, RaceTierBudget: *tierBudget,
		RaceTiers: parseInts(*raceTiers), RaceDeferPct: *deferPct,
		HumanMinute: *humanMinute, LifecycleBuyers: *lcBuyers, LifecycleTiers: parseInts(*lcTiers),
		LifecycleEvents: *lcEvents, LifecycleTrials: *lcTrials, LcTTL: *lcTTL, LcAbandonPct: *lcAbandon,
		LcAbandonExplicit: *lcExplicit, LcPaymentWindow: *lcWindow, LcGuard: *lcGuard, LcSweepEvery: *lcSweep,
		LcOutageFrom: *lcOutFrom, LcOutageTo: *lcOutTo, LcE0Skew: *lcSkew, LcTimeout: *lcTimeout,
		LcSweepBatch: 200, LcProbeCapPerEvent: *probeCap,
	}
	if *runID == "" {
		*runID = time.Now().UTC().Format("20060102T150405Z")
	}
	dsns := splitList(*dsn)
	run := &Run{
		Study: studyID, RunID: *runID, StartedAt: time.Now().UTC(), Environment: *envName,
		Repo:   provenance.RepoVersion{Commit: *repoCommit, Describe: *repoDesc, Dirty: *repoDirty, Tag: *runTag},
		Engine: *engine, Topology: *topology, ConnectionNodes: len(dsns),
		DesignID: d.ID, DesignTitle: d.Title, DesignFamily: d.Family, DesignSummary: d.Summary,
		Isolation: string(d.Isolation), NegativeControl: d.NegativeControl, ControlFires: d.ControlFires, Scale: *scale,
		Options: map[string]any{
			"conns": *workers, "load_conns": *streams, "duration": duration.String(), "warmup": warmup.String(),
			"trials": *trials, "seed": *seed, "tiers": *tiers, "write_ops": *writeOps, "phases": *phases,
			"hold_ttl": holdTTL.String(), "max_conflicts": *maxConf,
			"race_buyers": *raceBuyers, "race_timeout": raceTimeout.String(), "race_trials": *raceTrials,
			"race_tier_budget": tierBudget.String(), "race_tiers": *raceTiers, "race_defer_pct": *deferPct,
			"human_minute": humanMinute.String(), "lifecycle_buyers": *lcBuyers, "lifecycle_tiers": *lcTiers,
			"lifecycle_events_per_tier": *lcEvents, "lifecycle_trials": *lcTrials, "lc_ttl": lcTTL.String(),
			"lc_abandon_pct": *lcAbandon, "lc_abandon_explicit_pct": *lcExplicit, "lc_payment_window": lcWindow.String(),
			"lc_guard": lcGuard.String(), "lc_sweep_every": lcSweep.String(), "lc_outage_from": lcOutFrom.String(),
			"lc_outage_to": lcOutTo.String(), "lc_e0_skew": lcSkew.String(), "lc_timeout": lcTimeout.String(),
			"probe_cap": *probeCap, "statement_timeout_ms": *stmtTimeout, "connection_nodes": len(dsns),
		},
	}

	ctx, stopSignals := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stopSignals()
	var pools []ports.DB
	per := max(*workers, *raceBuyers, *lcBuyers)*2 + *streams + 8
	for _, one := range dsns {
		p, err := pgxdb.Open(ctx, pgxdb.Config{DSN: one, MaxConns: per, StatementTimeoutMS: *stmtTimeout, ApplicationName: "reservedseating-bench"})
		if err != nil {
			fatal(err)
		}
		pools = append(pools, p)
	}
	db := newSpreadDB(pools)
	defer db.Close()

	if err := execute(ctx, run, db, d, *cmd, *seed, splitList(*phases), s, *explainOut); err != nil {
		run.Error = err.Error()
		fmt.Fprintln(os.Stderr, "ERROR:", err)
	}
	if ctx.Err() != nil {
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
	run.ClientDBClockOffsetMS = clockOffset(ctx, db)
	fmt.Printf("  engine: %s\n", firstLine(run.EngineInfo.Version))
	if run.EngineInfo.EffectiveIsolation != "" {
		fmt.Printf("  effective isolation for READ COMMITTED requests: %s\n", run.EngineInfo.EffectiveIsolation)
	}
	fmt.Printf("  database clock minus client clock: %.3f ms\n", run.ClientDBClockOffsetMS)

	ds, err := Generate(run.Scale, seed, s.Tiers)
	if err != nil {
		return err
	}
	run.Dataset = ds.Summary()
	fmt.Printf("  dataset: %d venues, %d events, %d tickets and %d holds at load\n", len(ds.Venues), len(ds.Events), len(ds.Sold), len(ds.Holds))

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
	// The lifecycle runs on compressed human time; the race and the writes hold for a
	// real 40 minutes. Guard margin and E0's skew follow the phase's clock.
	guardFor := func(phase string) time.Duration {
		if phase == "lifecycle" {
			return s.human(s.LcGuard)
		}
		return s.LcGuard
	}
	skewFor := func(phase string) time.Duration {
		if !d.AppClock {
			return 0
		}
		if phase == "lifecycle" {
			return s.human(s.LcE0Skew)
		}
		return s.LcE0Skew
	}
	run.ReloadMS = map[string]float64{}
	var world *World
	var sl *Seller
	reload := func(why, phase string) error {
		t := time.Now()
		ph, err := Load(ctx, db, d, ds, s.Streams)
		if err != nil {
			return fmt.Errorf("load before %s: %w", why, err)
		}
		if run.Load == nil {
			run.Load = ph
		} else {
			run.ReloadMS[why] = msSince(t)
		}
		world = NewWorld(ds, ph.LoadNow, guardFor(phase))
		sl, err = NewSeller(db, d, world, skewFor(phase))
		if err != nil {
			return err
		}
		if d.ExpiryOnSweeper {
			// E1's steady state: its sweeper has released the expired holds found at load.
			swept := 0
			for {
				r, err := sl.Sweep(ctx, 0, 1000)
				if err != nil {
					return fmt.Errorf("pre-verification sweep: %w", err)
				}
				if r.Released == 0 {
					break
				}
				swept += r.Released
			}
			if run.PreSweep == 0 {
				run.PreSweep = swept
			}
		}
		return nil
	}

	fmt.Printf("  loading %s ...\n", d.ID)
	if err := reload("initial", "verify"); err != nil {
		return err
	}
	fmt.Printf("  loaded in %.1fs (copy %.1fs, index %.1fs)\n", run.Load.TotalMS/1000, run.Load.CopyMS/1000, run.Load.IndexMS/1000)

	if has("verify") {
		fmt.Println("  verifying answers against ground truth ...")
		vr, err := Verify(ctx, db, d, ds, run.Load.LoadNow)
		if err != nil {
			return err
		}
		run.Verify = vr
		for _, c := range vr.Checks {
			if !c.OK {
				fmt.Printf("    FAIL %-22s %s: expected %s, got %s\n", c.Query, c.Key, c.Expect, c.Got)
			}
		}
		au, err := RunAudit(ctx, db, d, world, "load", nil)
		if err != nil {
			return err
		}
		run.LoadAudit = au
		fmt.Printf("    %d/%d checks passed (margin %.0fs); audit on load: %s\n", vr.Passed, vr.Passed+vr.Failed, vr.MarginMS/1000, au)
		if vr.Failed > 0 || au.Violations() > 0 {
			return fmt.Errorf("design %s failed the correctness gate; refusing to report timings", d.ID)
		}
		if st, err := CollectStats(ctx, db, run.Engine); err == nil {
			run.Stats = st
		}
	}
	if cmd == "verify" {
		return nil
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
		rs, err := BenchmarkReads(ctx, db, d, ds, run.Load.LoadNow, s)
		stop()
		if err != nil {
			return err
		}
		run.Reads = rs
	}

	// Every writing phase gets its own fresh load; the initial load is still
	// pristine for the first one (explain rolled back, reads did not write).
	fresh := true
	nextLoad := func(why, phase string) error {
		if fresh {
			fresh = false
			world.ledger.setGuard(guardFor(phase))
			var err error
			sl, err = NewSeller(db, d, world, skewFor(phase))
			return err
		}
		return reload(why, phase)
	}

	if has("write") {
		fmt.Println("  write benchmark (fresh load per op):")
		for _, op := range s.WriteOps {
			if err := nextLoad(op, "write"); err != nil {
				return err
			}
			var results []WriteResult
			stop := probe(run, "write:"+op)
			switch op {
			case "hold":
				results = append(results, BenchmarkHoldSpread(ctx, sl, s))
			case "release":
				results = append(results, BenchmarkRelease(ctx, sl, s))
			case "cancel":
				results = append(results, BenchmarkCancel(ctx, sl, s))
			case "publish":
				for _, v := range ds.Venues {
					results = append(results, BenchmarkPublish(ctx, sl, v, s))
				}
			default:
				return fmt.Errorf("unknown write op %q", op)
			}
			stop()
			au, err := RunAudit(ctx, db, d, world, "write:"+op, nil)
			if err != nil {
				return err
			}
			fmt.Printf("      audit after %s: %s\n", op, au)
			for i := range results {
				results[i].Audit = au
			}
			run.Writes = append(run.Writes, results...)
		}
	}

	if has("race") {
		for trial := 1; trial <= max(s.RaceTrials, 1); trial++ {
			if err := nextLoad(fmt.Sprintf("race#%d", trial), "race"); err != nil {
				return err
			}
			fmt.Printf("  race trial %d/%d (%d buyers per event, fresh load):\n", trial, max(s.RaceTrials, 1), s.RaceBuyers)
			stop := probe(run, fmt.Sprintf("race#%d", trial))
			rr, err := RunRaces(ctx, db, sl, d, ds, s, trial)
			stop()
			if err != nil {
				return err
			}
			run.Races = append(run.Races, rr...)
		}
	}

	if has("lifecycle") {
		for trial := 1; trial <= max(s.LifecycleTrials, 1); trial++ {
			if err := nextLoad(fmt.Sprintf("lifecycle#%d", trial), "lifecycle"); err != nil {
				return err
			}
			fmt.Printf("  lifecycle trial %d/%d (%d buyers per event, one human minute = %s, fresh load):\n",
				trial, max(s.LifecycleTrials, 1), s.LifecycleBuyers, s.HumanMinute)
			stop := probe(run, fmt.Sprintf("lifecycle#%d", trial))
			lr, err := RunLifecycle(ctx, db, sl, d, ds, s, trial)
			stop()
			if err != nil {
				return err
			}
			run.Lifecycle = append(run.Lifecycle, lr...)
		}
	}
	return nil
}

// clockOffset: the database clock minus the client's, median of 20 samples.
func clockOffset(ctx context.Context, db ports.DB) float64 {
	var offs []float64
	for i := 0; i < 20; i++ {
		t0 := time.Now()
		var n time.Time
		if err := db.QueryRow(ctx, "SELECT clock_timestamp()").Scan(&n); err != nil {
			return 0
		}
		mid := t0.Add(time.Since(t0) / 2)
		offs = append(offs, float64(n.Sub(mid).Microseconds())/1000)
	}
	sort.Float64s(offs)
	return offs[len(offs)/2]
}

// probe records how much the client container was CPU-throttled during a phase.
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
	fmt.Fprintf(&sb, "keys:        the largest catalogue event, one of its live holds (or its section with the most\n")
	fmt.Fprintf(&sb, "             sold seats), a sold ticket. Write statements run inside transactions that are rolled back.\n\n")
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
