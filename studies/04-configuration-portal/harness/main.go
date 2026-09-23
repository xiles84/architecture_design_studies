package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	cp "configportal"

	"adsplatform/adapters/cgroup"
	"adsplatform/adapters/filestore"
	"adsplatform/adapters/pgxdb"
	"adsplatform/core/measure"
	"adsplatform/core/provenance"
	"adsplatform/ports"
)

// ClientCPUInfo records the client container's own CFS throttling across the cell.
// A throttled client puts tens of milliseconds into tails that belong to no
// design, so the counters travel with every result rather than being assumed
// absent (methodology 7).
type ClientCPUInfo struct {
	Before cgroup.CPUStat `json:"before"`
	After  cgroup.CPUStat `json:"after"`
	Delta  cgroup.CPUStat `json:"delta"`
	Seen   bool           `json:"cgroup_visible"`
}

// Study 04 — configuration portal.
//
// One cell is one (topology, design) pair on one freshly loaded database. It
// verifies, measures, and audits; it writes one JSON result and one readable
// plans file. A cell that fails its correctness gate aborts and reports no
// timing, because a design that answers the question wrongly, quickly, is worth
// nothing.

type Options struct {
	Scale      string
	Seed       int64
	Tier       int
	Fleet      int
	Regime     string
	CardMode   string
	Workers    int
	LoadConns  int
	Duration   time.Duration
	Warmup     time.Duration
	Trials     int
	Retries    int
	Batch      int
	Phases     string
	Cadence    string
	BurstSize  int
	Concurrency int
	Sample     int
}

type DatasetInfo struct {
	Seed             int64   `json:"seed"`
	ProductDefs      int     `json:"product_definitions"`
	Environments     int     `json:"environments"`
	BusinessUnits    int     `json:"business_units"`
	Installations    int     `json:"installations"`
	EntriesPerIP     int     `json:"entries_per_installed_product_tier"`
	CardinalityMode  string  `json:"cardinality_mode"`
	ValueRegime      string  `json:"value_regime"`
	KeyBytes         int     `json:"key_bytes"`
	ValueBytes       int     `json:"value_bytes"`
	TotalEntries     int     `json:"total_entries"`
	SerializedBytes  int64   `json:"total_serialized_bytes"`
	MeanSerializedB  float64 `json:"mean_serialized_bytes_per_installation"`
}

type EngineInfo struct {
	Version             string `json:"version"`
	EffectiveIsolation  string `json:"effective_isolation"`
	ServerVersion       string `json:"server_version"`
}

type ConcurrencyResult struct {
	Name              string  `json:"name"`
	Writers           int     `json:"writers"`
	Acknowledged      int64   `json:"acknowledged"`
	Retries           int64   `json:"retries"`
	Conflicts         int64   `json:"conflicts"`
	DurationS         float64 `json:"duration_s"`
	OpsPerSec         float64 `json:"ops_per_sec"`
	Latency           measure.LatencyStats `json:"latency"`
	FinalValue        int64   `json:"final_counter_value"`
	ExpectedValue     int64   `json:"expected_counter_value"`
	LostUpdates       int64   `json:"lost_updates"`
}

type CadenceResult struct {
	Regime string `json:"regime"`
	Period string `json:"period"`
	Fleet  int    `json:"fleet"`
	// OfferedRate is the calculated demand: installed products / period.
	OfferedRate float64 `json:"calculated_offered_rate_per_sec"`
	// Measured is what was actually offered to the portal in this run. For the
	// low regimes the rate is compressed and the run is labelled as such.
	Compressed bool                  `json:"rate_compressed"`
	Schedule   measure.ArrivalResult `json:"schedule"`
	Note       string                `json:"note"`
	// RollupLag is the time from a child change's commit to the design's own
	// parent aggregate reflecting it. It is present only for designs that
	// maintain an aggregate; a design with no aggregate has no staleness window
	// to measure. Timeouts are recorded separately, never folded into the stats.
	RollupLag         *measure.LatencyStats `json:"rollup_lag_ms,omitempty"`
	RollupLagTimeouts []string              `json:"rollup_lag_timeouts,omitempty"`
	RollupLagNote     string                `json:"rollup_lag_note,omitempty"`
}

type Result struct {
	Study       string `json:"study"`
	Design      string `json:"design"`
	DesignShort string `json:"design_short"`
	DesignTitle string `json:"design_title"`
	Family      string `json:"family"`
	Summary     string `json:"summary"`
	Pair        string `json:"controlled_pair,omitempty"`
	Negative    bool   `json:"negative_control,omitempty"`

	Topology    string `json:"topology"`
	Engine      string `json:"engine"`
	RunID       string `json:"run_id"`
	Environment string `json:"environment"`
	RepoCommit  string `json:"repo_commit"`
	RepoDesc    string `json:"repo_describe"`
	RepoDirty   bool   `json:"repo_dirty"`
	RunTag      string `json:"run_tag,omitempty"`

	StartedAt  string `json:"started_at"`
	FinishedAt string `json:"finished_at"`

	Options    Options        `json:"options"`
	EngineInfo EngineInfo     `json:"engine_info"`
	Dataset    DatasetInfo    `json:"dataset"`
	Load       LoadInfo       `json:"load"`
	Gate       GateInfo       `json:"gate"`
	ClientCPU  *ClientCPUInfo `json:"client_cpu,omitempty"`
	Reads      []measure.Result `json:"reads,omitempty"`
	Writes     []measure.Result `json:"writes,omitempty"`
	Concurrency []ConcurrencyResult `json:"concurrency,omitempty"`
	Cadence    []CadenceResult `json:"cadence,omitempty"`
	Audits     []AuditInfo    `json:"audits,omitempty"`
	Storage    map[string]int64 `json:"storage_bytes,omitempty"`

	Error string `json:"error,omitempty"`
}

func main() {
	var (
		dsn        = flag.String("dsn", os.Getenv("BENCH_DSN"), "PostgreSQL-protocol connection string; a comma-separated list spreads connections over nodes")
		engine     = flag.String("engine", "postgres", "postgres | yugabyte")
		topology   = flag.String("topology", "pg-single", "pg-single | yb-single | yb-cluster3")
		designID   = flag.String("design", "", "design id (see -cmd list)")
		scale      = flag.String("scale", "small", "tiny | small")
		seed       = flag.Int64("seed", 42, "dataset seed; changing it invalidates comparisons")
		tier       = flag.Int("entries-per-ip", 0, "configuration entries per installed product (0 = scale default)")
		fleet      = flag.Int("fleet", 0, "installed products (0 = scale default)")
		regime     = flag.String("value-regime", "small", "small | large")
		cardMode   = flag.String("cardinality-mode", "exact", "exact | constant-entries | constant-bytes | fixed-fleet | skewed")
		workers    = flag.Int("conns", 8, "concurrent workers for reads")
		loadConns  = flag.Int("load-conns", 4, "concurrent COPY streams during load")
		duration   = flag.Duration("duration", 5*time.Second, "measured duration per read trial")
		warmup     = flag.Duration("warmup", 2*time.Second, "discarded warmup per read")
		trials     = flag.Int("trials", 1, "trials per read measurement; throughput is the median")
		retries    = flag.Int("retries", 8, "bounded retries for guarded publications")
		batch      = flag.Int("batch", 8, "changes per atomic batch publication")
		phases     = flag.String("phases", "verify,explain,read,write,contention,cadence", "phases for -cmd full")
		cadence    = flag.String("cadence", "none", "none | daily | hourly | minutely | secondly")
		burst      = flag.Int("burst-size", 1, "installed products publishing together in a synchronized burst")
		concurrency = flag.Int("concurrency", 16, "concurrent writers on the hot key in the contention phase")
		sample     = flag.Int("sample", 12, "installations verified and audited per cell (0 = all)")
		stmtTO     = flag.Int("stmt-timeout-ms", 300000, "server-side statement_timeout")
		out        = flag.String("out", "", "path to write the JSON result")
		explainOut = flag.String("explain-out", "", "path to write the readable plans")
		envName    = flag.String("environment", os.Getenv("BENCH_ENVIRONMENT"), "environment id from docs/environments")
		runID      = flag.String("run-id", "", "run identifier")
		repoCommit = flag.String("repo-commit", "", "git commit the run was produced from (set by the runner)")
		repoDesc   = flag.String("repo-describe", "", "git describe --tags --always --dirty")
		repoDirty  = flag.Bool("repo-dirty", false, "the working tree had uncommitted changes")
		runTag     = flag.String("run-tag", "", "git tag created for this run")
		cmd        = flag.String("cmd", "full", "full | verify | list | report | digest")
		resultsDir = flag.String("results", "", "report/digest: results directory of one run")
		reportOut  = flag.String("report-out", "", "report: markdown file to write")
	)
	flag.Parse()

	switch *cmd {
	case "list":
		for _, d := range designs {
			fmt.Printf("%-24s %-6s %-12s %-12s pair=%s\n", d.ID, d.Short, d.Family, d.Kind, orNone(d.Pair))
		}
		return
	case "report":
		if err := writeReport(*resultsDir, *reportOut); err != nil {
			fatal(err)
		}
		return
	case "digest":
		// The digest is a content hash over every result file: it is what a later
		// analyst quotes to say which exact data an analysis was written about.
		dg, files, err := provenance.Digest(os.DirFS(*resultsDir))
		if err != nil {
			fatal(err)
		}
		fmt.Printf("%s (%d files)\n", dg, files)
		return
	}

	d, ok := designByID(*designID)
	if !ok {
		fatal(fmt.Errorf("unknown design %q (try -cmd list)", *designID))
	}
	if d.YBOnly && *engine != "yugabyte" {
		fatal(fmt.Errorf("design %s is YugabyteDB only", d.ID))
	}

	opts := Options{
		Scale: *scale, Seed: *seed, Tier: *tier, Fleet: *fleet, Regime: *regime,
		CardMode: *cardMode, Workers: *workers, LoadConns: *loadConns,
		Duration: *duration, Warmup: *warmup, Trials: *trials, Retries: *retries,
		Batch: *batch, Phases: *phases, Cadence: *cadence, BurstSize: *burst,
		Concurrency: *concurrency, Sample: *sample,
	}
	applyScale(&opts)

	res := Result{
		Study: "04-configuration-portal", Design: d.ID, DesignShort: d.Short,
		DesignTitle: d.Title, Family: d.Family, Summary: d.Summary, Pair: d.Pair,
		Negative: d.NegativeControl, Topology: *topology, Engine: *engine,
		RunID: *runID, Environment: *envName, RepoCommit: *repoCommit,
		RepoDesc: *repoDesc, RepoDirty: *repoDirty, RunTag: *runTag,
		StartedAt: time.Now().UTC().Format(time.RFC3339),
		Options:   opts,
	}
	if err := runCell(context.Background(), &res, d, *dsn, *explainOut, *stmtTO); err != nil {
		res.Error = err.Error()
	}
	res.FinishedAt = time.Now().UTC().Format(time.RFC3339)

	if *out != "" {
		if err := filestore.WriteJSON(*out, res); err != nil {
			fatal(err)
		}
	}
	if res.Error != "" {
		fmt.Fprintf(os.Stderr, "[%s/%s] FAILED: %s\n", *topology, d.ID, res.Error)
		os.Exit(1)
	}
}

func orNone(s string) string {
	if s == "" {
		return "-"
	}
	return s
}

// applyScale turns a named scale into concrete sizes. `tiny` exists so that every
// design, invariant and negative control can be exercised in minutes; it is a dev
// check and its numbers are never reported.
func applyScale(o *Options) {
	switch o.Scale {
	case "tiny":
		if o.Tier == 0 {
			o.Tier = 8
		}
		if o.Fleet == 0 {
			o.Fleet = 6
		}
		if o.Duration == 5*time.Second {
			o.Duration = 400 * time.Millisecond
		}
		if o.Warmup == 2*time.Second {
			o.Warmup = 100 * time.Millisecond
		}
		if o.Sample == 12 {
			o.Sample = 0
		}
	case "small":
		if o.Tier == 0 {
			o.Tier = 60
		}
		if o.Fleet == 0 {
			o.Fleet = 60
		}
	}
}

func runCell(ctx context.Context, res *Result, d Design, dsn, explainPath string, stmtTO int) error {
	openCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	db, err := pgxdb.Open(openCtx, pgxdb.Config{
		DSN: dsn, MaxConns: 48, StatementTimeoutMS: stmtTO,
		ApplicationName: "study04-" + d.Short,
	})
	if err != nil {
		return fmt.Errorf("connect: %w", err)
	}
	defer db.Close()

	res.EngineInfo = probeEngine(ctx, db, res.Engine)

	cat, err := loadCatalogue(cp.SQL, d)
	if err != nil {
		return err
	}
	ds := BuildDataset(res.Options.Seed, res.Options.Fleet, res.Options.Tier,
		res.Options.Regime, res.Options.CardMode)
	res.Dataset = describeDataset(ds)

	// The client's own throttling is captured across the whole cell. The client
	// runs in a container under a CFS quota on the same cores as the database,
	// and a throttled client manufactures tail latency that belongs to no design.
	cpuBefore, seenBefore := cgroup.Read()

	// One ledger for the whole cell. A per-phase ledger reconciles the database
	// against a state that no longer exists: the first full dev check audited the
	// contention phase against a fresh dataset and reported every installation as
	// missing keys the write phase had legitimately changed.
	led := newLedger(ds)
	for _, p := range strings.Split(res.Options.Phases, ",") {
		if err := runPhase(ctx, res, d, cat, ds, db, led, strings.TrimSpace(p), explainPath); err != nil {
			return fmt.Errorf("phase %s: %w", p, err)
		}
	}

	if cpuAfter, seenAfter := cgroup.Read(); seenBefore && seenAfter {
		res.ClientCPU = &ClientCPUInfo{
			Before: cpuBefore, After: cpuAfter, Delta: cpuAfter.Sub(cpuBefore), Seen: true,
		}
	}
	return nil
}

func describeDataset(ds *Dataset) DatasetInfo {
	di := DatasetInfo{
		Seed: ds.Seed, ProductDefs: len(ds.Definitions), Environments: len(ds.Environments),
		BusinessUnits: len(ds.BusinessUnits), Installations: len(ds.Installations),
		EntriesPerIP: ds.Tier, CardinalityMode: ds.CardMode, ValueRegime: ds.Regime,
		KeyBytes: ds.KeyBytes, ValueBytes: ds.ValueBytes,
		TotalEntries: ds.TotalEntries, SerializedBytes: ds.SerializedBytes,
	}
	if len(ds.Installations) > 0 {
		di.MeanSerializedB = float64(ds.SerializedBytes) / float64(len(ds.Installations))
	}
	return di
}

func runPhase(ctx context.Context, res *Result, d Design, cat *catalogue, ds *Dataset, db ports.DB, led *ledger, phase, explainPath string) error {
	switch phase {
	case "":
		return nil
	case "verify":
		DropSchema(ctx, db)
		if err := ApplySchema(ctx, db, cat); err != nil {
			return err
		}
		li, err := BulkLoad(ctx, db, cat, ds)
		if err != nil {
			return err
		}
		res.Load = li
		res.Storage = relationSizes(ctx, db, res.Engine)
		g := verifyReads(ctx, db, cat, ds, res.Options.Sample)
		res.Gate = g
		if !g.Passed {
			return fmt.Errorf("correctness gate failed: %s", strings.Join(g.Failures, "; "))
		}
		return nil

	case "explain":
		if explainPath == "" {
			return nil
		}
		return capturePlans(ctx, db, cat, ds, res.Engine, explainPath)

	case "read":
		c := &cell{design: d, cat: cat, db: db, ds: ds, led: led, opts: res.Options}
		for _, name := range baseReads {
			r := measure.Run(ctx, measure.Budget{
				Workers: res.Options.Workers, Warmup: res.Options.Warmup,
				Duration: res.Options.Duration, Trials: res.Options.Trials,
			}, name, c.readOp(name))
			res.Reads = append(res.Reads, r)
		}
		return nil

	case "write":
		ops := []string{"w01_modify_key", "w02_add_key", "w03_delete_key",
			"w04_publish_batch", "w05_replace_all", "w06_update_metadata"}
		if d.Kind == Doc {
			ops = []string{"wd1_write_doc", "w06_update_metadata"}
		}
		c := &cell{design: d, cat: cat, db: db, ds: ds, led: led, opts: res.Options}
		for _, name := range ops {
			// One worker: the isolated write measurement answers "what does one
			// update cost". Concurrency is a different question, asked separately
			// with controlled contention, and mixing them would let either one
			// explain the other's number.
			r := measure.Run(ctx, measure.Budget{
				Workers: 1, Warmup: res.Options.Warmup, Duration: res.Options.Duration,
				Trials: res.Options.Trials,
			}, name, c.writeOp(name))
			res.Writes = append(res.Writes, r)
		}
		a := audit(ctx, db, cat, ds, led, res.Options.Sample)
		res.Audits = append(res.Audits, a)
		// A negative control that breaks an invariant has done its job. Failing
		// the cell there would make the control look like a broken design, and
		// the control firing is exactly what licenses the other designs'
		// clean audits.
		if !a.Passed && !d.NegativeControl {
			return fmt.Errorf("audit after writes failed: %s", strings.Join(a.Failures, "; "))
		}
		return nil

	case "contention":
		// A design with no arbitration has nothing to contend with -- except the
		// drift control, whose race is on the derived rollup rather than on a key
		// and which therefore must still be run here.
		if !d.hasContention() {
			return nil
		}
		return runContention(ctx, res, d, cat, ds, db, led)

	case "cadence":
		if res.Options.Cadence == "none" {
			return nil
		}
		return runCadence(ctx, res, d, cat, ds, db, led)

	case "audit":
		a := audit(ctx, db, cat, ds, led, 0)
		res.Audits = append(res.Audits, a)
		return nil
	}
	return fmt.Errorf("unknown phase %q", phase)
}

// runContention puts every writer on one key of one installed product. That is
// the only scenario in which a lost update is possible at all, so it is the
// scenario the negative control must fire in (methodology 5a).
func runContention(ctx context.Context, res *Result, d Design, cat *catalogue, ds *Dataset, db ports.DB, led *ledger) error {
	ip := ds.Installations[0]
	key := ds.Entries[ip.ID][0].Key
	c := &cell{design: d, cat: cat, db: db, ds: ds, led: led, opts: res.Options,
		concIP: ip.ID, concKey: key}

	// The phase prepares its own hot spot. An earlier phase in the same cell may
	// have replaced this installation's whole key set (w05_replace_all does), and
	// a race on a key that no longer exists measures the error path rather than
	// the arbitration. This is a routine, recorded choice: the counter starts at
	// zero and the ledger is told so.
	if err := c.publish(ctx, ip.ID, func(q ports.Queryer) error {
		_, e := c.cat.exec(ctx, q, sW02, map[string]any{
			"installed_product_id": ip.ID, "key": key, "value": "0"})
		return e
	}); err != nil {
		return fmt.Errorf("prepare contention hot spot: %w", err)
	}
	c.led.set(ip.ID, key, "0")

	start := time.Now()
	r := measure.Run(ctx, measure.Budget{
		Workers: res.Options.Concurrency, Warmup: res.Options.Warmup,
		Duration: res.Options.Duration, Trials: res.Options.Trials, SafetyLimit: 2 * time.Minute,
	}, "contention:"+string(d.Conc), c.incrementOp())
	elapsed := time.Since(start).Seconds()

	// The expected counter is where it started plus every increment the client saw
	// acknowledged. A guard that works makes the database agree; a guard that does
	// not leaves the database behind (INV-3).
	var final int64
	vals, _ := cat.args(sAuditKeys, map[string]any{"installed_product_id": ip.ID})
	rows, err := db.Query(ctx, cat.stmt(sAuditKeys).SQL, vals...)
	if err != nil {
		return err
	}
	kvs, err := scanPairs(rows)
	if err != nil {
		return err
	}
	for _, kv := range kvs {
		if kv.Key == key {
			final = int64(atoi(kv.Value))
		}
	}
	// The counter was prepared at zero above, and each acknowledged increment
	// adds one.
	want := led.incAcked.Load()

	cr := ConcurrencyResult{
		Name: r.Name, Writers: res.Options.Concurrency, Acknowledged: led.acked.Load(),
		Retries: led.retries.Load(), Conflicts: led.conflicts.Load(), DurationS: elapsed,
		OpsPerSec: r.OpsPerSec, Latency: r.Latency,
		FinalValue: final, ExpectedValue: want,
	}
	if want > final {
		cr.LostUpdates = want - final
	}
	res.Concurrency = append(res.Concurrency, cr)

	a := audit(ctx, db, cat, ds, led, res.Options.Sample)
	res.Audits = append(res.Audits, a)

	// A negative control that does not fire is a finding about the workload, not
	// a licence to soften the invariant. It is recorded and the cell fails, so a
	// control that cannot fire is never reported as a control that did not need
	// to.
	if d.NegativeControl {
		fired, why := controlFired(d, cr, a)
		if !fired {
			return fmt.Errorf("negative control %s did not fire: %s", d.ID, why)
		}
	}
	if !d.NegativeControl && !a.Passed {
		return fmt.Errorf("audit after contention failed: %s", strings.Join(a.Failures, "; "))
	}
	return nil
}

// runCadence offers publications at the fleet's calculated rate. The rate is a
// calculation (installed products / period); a low rate is compressed so the run
// finishes, and the result says so rather than pretending to have waited.
func runCadence(ctx context.Context, res *Result, d Design, cat *catalogue, ds *Dataset, db ports.DB, led *ledger) error {
	periods := map[string]time.Duration{
		"daily": 24 * time.Hour, "hourly": time.Hour,
		"minutely": time.Minute, "secondly": time.Second,
	}
	period, ok := periods[res.Options.Cadence]
	if !ok {
		return fmt.Errorf("unknown cadence %q", res.Options.Cadence)
	}
	rate := measure.CadenceRate(len(ds.Installations), period)

	// Compression: a rate this harness cannot sustain is scaled up, and the run
	// is labelled. Nothing here is reported as the temporal result it stands in
	// for.
	const maxRate = 500.0
	compressed := false
	offered := rate
	if offered < 1 {
		offered = 1
		compressed = true
	}
	if offered > maxRate {
		offered = maxRate
		compressed = true
	}

	c := &cell{design: d, cat: cat, db: db, ds: ds, led: led, opts: res.Options}
	ar := measure.RunOpenLoop(ctx, measure.Schedule{
		Rate: offered, Duration: res.Options.Duration, Workers: res.Options.Workers,
		Queue: res.Options.Workers * 8, Burst: res.Options.BurstSize,
		Jitter: res.Options.BurstSize <= 1, Seed: res.Options.Seed,
	}, "cadence:"+res.Options.Cadence, c.writeOp("w01_modify_key"))

	note := "measured directly"
	if compressed {
		note = "rate compressed for the run duration; the offered rate is a calculation and this is a compressed-time validation, not a temporal result"
	}
	cr := CadenceResult{
		Regime: res.Options.Cadence, Period: period.String(), Fleet: len(ds.Installations),
		OfferedRate: rate, Compressed: compressed, Schedule: ar, Note: note,
	}
	// A maintained aggregate gets its staleness window measured; a design with no
	// aggregate has none to measure, and saying so is better than a missing field.
	if d.Rollup == RollupTrigger || d.Rollup == RollupApp {
		stats, timeouts, err := c.stalenessProbe(ctx, res.Options.Sample)
		if err != nil {
			return fmt.Errorf("rollup staleness probe: %w", err)
		}
		cr.RollupLag = &stats
		cr.RollupLagTimeouts = timeouts
		cr.RollupLagNote = fmt.Sprintf(
			"time from a child change's commit to this design's own aggregate reflecting it, %d sample(s); trigger and application maintenance are both inside the publishing transaction",
			stats.Count)
	}
	res.Cadence = append(res.Cadence, cr)
	return nil
}

// capturePlans writes EXPLAIN output for every read and write statement. PostgreSQL
// gets ANALYZE + BUFFERS; YugabyteDB gets ANALYZE + DIST, whose RPC counts are the
// portable distributed-cost signal on a single machine.
func capturePlans(ctx context.Context, db ports.DB, cat *catalogue, ds *Dataset, engine, path string) error {
	ip := ds.Installations[0]
	key := ds.Entries[ip.ID][0].Key
	vals := map[string]any{
		"installed_product_id": ip.ID, "key": key, "known_revision": int64(0),
		"scope_kind": "product_definition", "scope_id": int64(ip.ProductDefID),
	}
	var b strings.Builder
	fmt.Fprintf(&b, "Read and write plans — design %s, engine %s\n", cat.design.ID, engine)
	fmt.Fprintf(&b, "installations=%d entries_per_installation=%d value_regime=%s\n",
		len(ds.Installations), ds.Tier, ds.Regime)
	for _, name := range cat.order {
		s := cat.stmt(name)
		if strings.HasPrefix(name, "w_") || strings.HasPrefix(name, "a_") {
			continue
		}
		a, err := cat.args(name, vals)
		if err != nil {
			// A statement whose signature this probe does not know is skipped
			// rather than guessed at.
			continue
		}
		prefix := "EXPLAIN (ANALYZE, BUFFERS, VERBOSE) "
		if engine == "yugabyte" {
			prefix = "EXPLAIN (ANALYZE, DIST) "
		}
		fmt.Fprintf(&b, "\n=== %s\n", name)
		fmt.Fprintf(&b, "-- declared parameters: %s\n", strings.Join(s.Params, ", "))
		rows, err := db.Query(ctx, prefix+s.SQL, a...)
		if err != nil {
			fmt.Fprintf(&b, "-- EXPLAIN failed: %v\n", err)
			continue
		}
		for rows.Next() {
			var line string
			if err := rows.Scan(&line); err != nil {
				break
			}
			b.WriteString(line)
			b.WriteString("\n")
		}
		rows.Close()
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(b.String()), 0o644)
}

// relationSizes records storage. On YugabyteDB pg_total_relation_size is
// meaningless, so the study does not report it there rather than reporting a
// number that would be read as bytes on disk.
func relationSizes(ctx context.Context, db ports.DB, engine string) map[string]int64 {
	if engine == "yugabyte" {
		return nil
	}
	out := map[string]int64{}
	for _, t := range []string{"product_definition", "environment", "business_unit",
		"installed_product", "config_entry", "config_document", "load_entry"} {
		var n int64
		err := db.QueryRow(ctx,
			"SELECT CASE WHEN to_regclass($1) IS NULL THEN 0 ELSE pg_total_relation_size(to_regclass($1)) END",
			t).Scan(&n)
		if err == nil && n > 0 {
			out[t] = n
		}
	}
	return out
}

func probeEngine(ctx context.Context, db ports.DB, engine string) EngineInfo {
	ei := EngineInfo{Version: engine}
	if err := db.QueryRow(ctx, "SELECT version()").Scan(&ei.ServerVersion); err != nil {
		ei.ServerVersion = "unknown"
	}
	var iso string
	if err := db.QueryRow(ctx, "SHOW transaction_isolation").Scan(&iso); err != nil {
		iso = "unknown"
	}
	ei.EffectiveIsolation = iso
	return ei
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, "fatal:", err)
	os.Exit(1)
}

// nowMS is the harness's only clock. Milliseconds, because every duration that
// reaches a result is either an elapsed budget or a percentile.
func nowMS() float64 { return float64(time.Now().UnixNano()) / 1e6 }

// controlFired decides, for one negative control, whether its invariant actually
// broke under the contention workload. A control that did not fire is a finding
// about the workload -- and it means the other designs' clean audits are not
// evidence of anything -- so it fails the cell rather than being recorded as a
// quiet success (methodology 5a).
func controlFired(d Design, cr ConcurrencyResult, a AuditInfo) (bool, string) {
	if d.Conc == ConcUnsafe {
		if cr.LostUpdates > 0 {
			return true, ""
		}
		return false, fmt.Sprintf("no lost update: %d writers, %d acknowledged, counter %d (expected %d)",
			cr.Writers, cr.Acknowledged, cr.FinalValue, cr.ExpectedValue)
	}
	// The drift control breaks derived data, so it is its own rollup audit that
	// must catch it, not the counter.
	for name, n := range a.DesignChecks {
		if n > 0 {
			return true, ""
		}
		_ = name
	}
	if !a.Passed {
		return true, ""
	}
	return false, "the stored rollup still agreed with its entries after a concurrent publication"
}
