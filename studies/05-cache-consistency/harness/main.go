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

	cp "cacheconsistency"

	"adsplatform/adapters/cgroup"
	"adsplatform/adapters/filestore"
	"adsplatform/adapters/pgxdb"
	"adsplatform/core/catalog"
	"adsplatform/core/inspect"
	"adsplatform/core/provenance"
	"adsplatform/ports"
)

// Study 05 -- external cache throughput and consistency.
//
// One cell is one scenario on one freshly loaded database. It verifies, measures,
// injects its faults, audits, and writes one JSON result plus a readable plans file.
// A cell that fails its correctness gate aborts and reports no timing, because a
// design that answers the question wrongly, quickly, is worth nothing.

func main() {
	var (
		dsn          = flag.String("dsn", os.Getenv("BENCH_DSN"), "PostgreSQL-protocol connection string; a comma-separated list spreads connections over nodes")
		engine       = flag.String("engine", "postgres", "postgres | yugabyte")
		topology     = flag.String("topology", "pg-single", "pg-single | yb-single | yb-cluster3")
		scenario     = flag.String("scenario", "", "scenario id (see -cmd list)")
		scale        = flag.String("scale", "small", "tiny | small | medium")
		seed         = flag.Int64("seed", 42, "dataset seed; changing it invalidates comparisons")
		faultSeed    = flag.Int64("fault-seed", 4242, "fault-injection schedule seed")
		conns        = flag.Int("conns", 8, "concurrent reader workers (per instance group)")
		writeConns   = flag.Int("write-conns", 2, "writer workers for the hotspot phase")
		duration     = flag.Duration("duration", 0, "measured duration per phase (0 = scale default)")
		warmup       = flag.Duration("warmup", 0, "discarded warmup per phase (0 = scale default)")
		trials       = flag.Int("trials", 1, "trials per measurement; throughput is the median")
		retries      = flag.Int("retries", 8, "bounded retries for the optimistic version CAS")
		phases       = flag.String("phases", "verify,explain,calibrate,warm,mixed,hotspot,stampede,instances,churn,faults,audit", "phases for -cmd full")
		instances    = flag.Int("instances", 1, "logical application instances")
		churnFor     = flag.Duration("churn-duration", 0, "sustained churn duration (0 = skip; the 300 s TTL is never shortened)")
		stampede     = flag.Int("stampede-readers", 16, "readers released together in the stampede phase")
		stampedeN    = flag.Int("stampede-keys", 8, "keys in the stampede burst")
		hotKeys      = flag.Int("hot-keys", 4, "keys the hotspot phase targets")
		capacityKB   = flag.Int("cache-capacity-kib", 0, "cache capacity in KiB (0 = derived from the working set)")
		fits         = flag.Bool("cache-fits", false, "size the cache so the whole working set fits (the control condition)")
		redisAddr    = flag.String("redis-addr", os.Getenv("BENCH_REDIS_ADDR"), "host:port of the cache under test")
		redisMB      = flag.Int("redis-maxmemory-mb", 0, "recorded Redis maxmemory, for the report only")
		frame        = flag.String("resource-frame", "db-only", "db-only | add-cache | equal-total")
		sample       = flag.Int("sample", 12, "donors verified and audited per cell (0 = all)")
		stmtTO       = flag.Int("stmt-timeout-ms", 300000, "server-side statement_timeout")
		out          = flag.String("out", "", "path to write the JSON result")
		explainOut   = flag.String("explain-out", "", "path to write the readable plans")
		envName      = flag.String("environment", os.Getenv("BENCH_ENVIRONMENT"), "environment id from docs/environments")
		runID        = flag.String("run-id", "", "run identifier")
		repoCommit   = flag.String("repo-commit", "", "git commit the run was produced from (set by the runner)")
		repoDesc     = flag.String("repo-describe", "", "git describe --tags --always --dirty")
		repoDirty    = flag.Bool("repo-dirty", false, "the working tree had uncommitted changes")
		runTag       = flag.String("run-tag", "", "git tag created for this run")
		benchImage   = flag.String("bench-image", "", "benchmark image name recorded in the result")
		benchImageID = flag.String("bench-image-id", "", "benchmark image id recorded in the result")
		cmd          = flag.String("cmd", "full", "full | verify | list | report | digest | probe")
		resultsDir   = flag.String("results", "", "report/digest: results directory of one run")
		reportOut    = flag.String("report-out", "", "report: markdown file to write")
	)
	flag.Parse()

	switch *cmd {
	case "list":
		for _, d := range designs {
			fmt.Printf("%-52s %-14s %-10s pair=%s\n", d.ID, d.Short, d.Group, orNone(d.Pair))
		}
		return
	case "report":
		if err := writeReport(*resultsDir, *reportOut); err != nil {
			fatal(err)
		}
		return
	case "digest":
		dg, files, err := provenance.Digest(os.DirFS(*resultsDir))
		if err != nil {
			fatal(err)
		}
		fmt.Printf("%s (%d files)\n", dg, files)
		return
	case "probe":
		res, err := probeRedisCapabilities(context.Background(), *redisAddr, *out)
		if err != nil {
			fatal(err)
		}
		_ = res
		return
	}

	d, ok := designByID(*scenario)
	if !ok {
		fatal(fmt.Errorf("unknown scenario %q (try -cmd list)", *scenario))
	}
	if d.YBOnly && *engine != "yugabyte" {
		fatal(fmt.Errorf("scenario %s is YugabyteDB only", d.ID))
	}

	opts := Options{
		Scale: *scale, Seed: *seed, FaultSeed: *faultSeed, Workers: *conns,
		Writers: *writeConns, Duration: *duration, Warmup: *warmup, Trials: *trials,
		Retries: *retries, Phases: *phases, Instances: *instances, ChurnFor: *churnFor,
		Stampede: *stampede, StampedeN: *stampedeN, HotKeys: *hotKeys,
		CapacityKB: *capacityKB, CacheFit: *fits, RedisAddr: *redisAddr,
		RedisMaxMB: *redisMB, ResourceFrame: *frame,
	}
	applyScale(&opts)

	res := CellResult{
		Study: "05-cache-consistency", Scenario: d.ID, ScenarioShort: d.Short,
		Group: d.Group, Title: d.Title, Summary: d.Summary, Risk: d.Risk, Pair: d.Pair,
		Negative: d.NegativeControl,
		ModelDim: string(d.effectiveModel()), VersionDim: string(d.Ver),
		BackendDim: string(d.Backend), StrategyDim: string(d.Strategy),
		FreshnessDim: string(d.Freshness), WritersDim: string(d.Writers),
		Topology: *topology, Engine: *engine, RunID: *runID, Environment: *envName,
		RepoCommit: *repoCommit, RepoDesc: *repoDesc, RepoDirty: *repoDirty, RunTag: *runTag,
		BenchImage: *benchImage, BenchImageID: *benchImageID,
		StartedAt: time.Now().UTC().Format(time.RFC3339),
		Options:   opts,
	}

	err := runCell(context.Background(), &res, d, *dsn, *explainOut, *stmtTO, *sample, *duration == 0)
	if err != nil {
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

// applyScale turns a named scale into concrete sizes and defaults. `tiny` exists so
// every scenario, invariant and control can be exercised in minutes; its numbers are
// never reported.
func applyScale(o *Options) {
	switch o.Scale {
	case "tiny":
		if o.Duration == 0 {
			o.Duration = 400 * time.Millisecond
		}
		if o.Warmup == 0 {
			o.Warmup = 100 * time.Millisecond
		}
	case "medium":
		if o.Duration == 0 {
			o.Duration = 3 * time.Second
		}
		if o.Warmup == 0 {
			o.Warmup = time.Second
		}
	default:
		if o.Duration == 0 {
			o.Duration = 2 * time.Second
		}
		if o.Warmup == 0 {
			o.Warmup = 500 * time.Millisecond
		}
	}
}

// capacityFor derives the cache's byte capacity. The primary condition is
// deliberately BELOW the logical working set so eviction occurs; the fits control is
// above it. Both numbers are recorded, because "capacity" is only meaningful beside
// the working set it is compared with.
func capacityFor(o *Options, ds *Dataset) int64 {
	workingSet := int64(len(ds.People)) * 2500
	if o.CacheFit {
		return workingSet * 2
	}
	if o.CapacityKB > 0 {
		return int64(o.CapacityKB) * 1024
	}
	c := workingSet / 2
	if c > 256*1024 {
		c = 256 * 1024
	}
	if c < 16*1024 {
		c = 16 * 1024
	}
	return c
}

func runCell(ctx context.Context, res *CellResult, d Design, dsn, explainPath string, stmtTO, sample int, durationWasDefault bool) error {
	openCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	db, err := pgxdb.Open(openCtx, pgxdb.Config{
		DSN: dsn, MaxConns: 48, StatementTimeoutMS: stmtTO,
		ApplicationName: "study05-" + d.Short,
	})
	if err != nil {
		return fmt.Errorf("connect: %w", err)
	}
	defer db.Close()

	res.EngineInfo = probeEngine(ctx, db, res.Engine)
	res.StrictPolicy = string(strictPolicyFor(d, res.Options.Instances))

	cat, err := loadCatalogue(cp.SQL, d)
	if err != nil {
		return err
	}
	ds := BuildDataset(res.Options.Seed, res.Options.Scale)
	res.Dataset = describeDataset(ds)
	res.Options.CapacityKB = int(capacityFor(&res.Options, ds) / 1024)

	// The client's own throttling is captured across the whole cell. The client runs
	// in a container under a CFS quota on the same cores as the database, and a
	// throttled client manufactures tail latency that belongs to no scenario.
	cpuBefore, seenBefore := cgroup.Read()

	orc := newOracle(ds)
	capBytes := capacityFor(&res.Options, ds)
	insts, err := buildInstances(d.Backend, res.Options.Instances, capBytes, res.Options.RedisAddr, res.Options.Workers+8)
	if err != nil {
		return err
	}
	defer func() {
		for _, in := range insts {
			_ = in.Store.Close(ctx)
		}
	}()
	// A scenario that needs a cache must not run against a cache that is not there.
	// A cache outage degrades to authoritative reads by design, so without this check
	// an unreachable broker would look like a scenario that simply never hit.
	if d.Backend == BackendRedis {
		if p, ok := insts[0].Store.(interface{ Ping(context.Context) error }); ok {
			if err := p.Ping(ctx); err != nil {
				return fmt.Errorf("the cache under test does not answer: %w", err)
			}
		} else {
			return fmt.Errorf("the redis backend does not implement Ping; refusing to measure against an unverified cache")
		}
	}

	log := newReadLog()
	ad := newAdapter(d, cat, db, ds, orc, res.Options, insts, log)
	c := &cell{d: d, cat: cat, db: db, ds: ds, orc: orc, ad: ad, res: res, opts: res.Options}

	for _, p := range strings.Split(res.Options.Phases, ",") {
		if err := c.runPhase(ctx, strings.TrimSpace(p), explainPath, sample); err != nil {
			return fmt.Errorf("phase %s: %w", p, err)
		}
	}

	c.snapshotCache()
	res.StrictViolations = ad.strictWrong.Load()
	// res.ImpossibleValues was accumulated from the PER-PHASE summaries, which is the
	// only place a read is counted once. Adding the adapter's running counter here as
	// well double-counted every occurrence.

	// The acceptance rules. A strict cell with a wrong read is a correctness
	// failure and cannot support a performance conclusion; an impossible cache value
	// invalidates the cell under every freshness policy.
	if d.Freshness == FreshStrict && !d.NegativeControl && res.StrictViolations > 0 {
		return fmt.Errorf("STRICT CONTRACT VIOLATED: %d stale-after-ack read(s); the cell cannot support a performance conclusion", res.StrictViolations)
	}
	if res.ImpossibleValues > 0 {
		return fmt.Errorf("IMPOSSIBLE CACHE VALUE: %d read(s) returned a state the database never held; the cell is invalid", res.ImpossibleValues)
	}
	// A negative control that did not fire is a finding about the workload, not a
	// licence to soften the invariant.
	if d.NegativeControl {
		if err := negativeControlFired(res, d); err != nil {
			return err
		}
	}

	if cpuAfter, seenAfter := cgroup.Read(); seenBefore && seenAfter {
		res.ClientCPU = &ClientCPUInfo{Before: cpuBefore, After: cpuAfter, Delta: cpuAfter.Sub(cpuBefore), Seen: true}
	}
	_ = durationWasDefault
	return nil
}

func (c *cell) runPhase(ctx context.Context, phase, explainPath string, sample int) error {
	switch phase {
	case "":
		return nil

	case "verify":
		DropSchema(ctx, c.db)
		li, err := BulkLoad(ctx, c.db, c.cat, c.ds)
		if err != nil {
			return err
		}
		c.res.Load = li
		c.res.Storage = relationSizes(ctx, c.db, c.res.Engine)
		// A cell starts from a COLD cache: eviction counts, hit rates and expiry
		// observations are per-cell, and a previous cell's entries would make them
		// a measurement of the previous cell.
		if c.d.hasCache() {
			for _, in := range c.ad.inst {
				_ = in.Store.Flush(ctx)
			}
		}
		g := verifyReads(ctx, c.db, c.cat, c.ds, c.orc, sample)
		c.res.Gate = g
		if !g.Passed {
			return fmt.Errorf("correctness gate failed: %s", strings.Join(g.Failures, "; "))
		}
		if g.Checks == 0 {
			return fmt.Errorf("the correctness gate ran no checks")
		}
		return nil

	case "explain":
		if explainPath == "" {
			return nil
		}
		return capturePlans(ctx, c.db, c.cat, c.ds, c.res.Engine, explainPath)

	case "calibrate":
		c.calibrate(ctx)
		return nil

	case "warm":
		c.resetWrong()
		return c.runWarm(ctx)

	case "mixed":
		c.resetWrong()
		return c.runMixed(ctx)

	case "hotspot":
		c.resetWrong()
		return c.runHotspot(ctx)

	case "stampede":
		c.resetWrong()
		return c.runStampede(ctx)

	case "instances":
		c.resetWrong()
		return c.runInstances(ctx)

	case "churn":
		c.resetWrong()
		return c.runChurn(ctx)

	case "faults":
		c.resetWrong()
		return c.runFaults(ctx)

	case "audit":
		a := auditCell(ctx, c.db, c.cat, c.ds, c.orc, sample)
		c.res.Audits = append(c.res.Audits, a)
		if !a.Passed && !c.d.NegativeControl {
			return fmt.Errorf("ledger replay failed: %s", strings.Join(a.Failures, "; "))
		}
		return nil
	}
	return fmt.Errorf("unknown phase %q", phase)
}

func describeDataset(ds *Dataset) DatasetInfo {
	di := DatasetInfo{
		Seed: ds.Seed, Charities: len(ds.Charities), People: len(ds.People),
		Donations: ds.TotalDonations, SerializedBytes: ds.SerializedBytes,
	}
	di.WorkingSetBytes = int64(len(ds.People)) * 2500
	if len(ds.People) > 0 {
		di.MeanPayloadB = float64(di.WorkingSetBytes) / float64(len(ds.People))
	}
	return di
}

func probeEngine(ctx context.Context, db ports.DB, engine string) EngineInfo {
	ei := EngineInfo{Version: engine}
	if err := db.QueryRow(ctx, "SELECT version()").Scan(&ei.ServerVersion); err != nil {
		ei.ServerVersion = "unknown: " + err.Error()
	}
	var iso string
	if err := db.QueryRow(ctx, "SHOW transaction_isolation").Scan(&iso); err != nil {
		iso = "unknown"
	}
	ei.EffectiveIsolation = iso
	if engine == "yugabyte" {
		// Asked inside a READ COMMITTED transaction, because that is the question:
		// when a design requests RC, what does it get?
		tx, err := db.Begin(ctx, ports.ReadCommitted)
		if err == nil {
			var eff string
			if err := tx.QueryRow(ctx, "SHOW yb_effective_transaction_isolation_level").Scan(&eff); err == nil {
				ei.EffectiveIsolation = eff
			}
			_ = tx.Rollback(ctx)
		}
	}
	return ei
}

// capturePlans writes EXPLAIN output for the statements a reader would want to see.
// A cache HIT has no SQL plan: the report says the database operation was saved
// rather than inventing a plan for it.
func capturePlans(ctx context.Context, db ports.DB, cat *catalogue, ds *Dataset, engine, path string) error {
	p := ds.People[0]
	vals := map[string]any{
		"person_id": p.ID, "charity_id": p.CharityID,
	}
	wanted := []string{sPortalPerson, sPortalRecent, sCharityRecent, sPersonVersion}
	var b strings.Builder
	fmt.Fprintf(&b, "Plans — scenario %s, engine %s\n", cat.design.ID, engine)
	fmt.Fprintf(&b, "people=%d scale=%s\n", len(ds.People), ds.Scale)
	fmt.Fprintf(&b, "A cache hit executes no SQL. These are the statements a miss, a fill and a strict validation run.\n")
	for _, name := range wanted {
		if !cat.has(name) {
			fmt.Fprintf(&b, "\n=== %s\n-- not part of this scenario's catalogue\n", name)
			continue
		}
		st := cat.stmt(name)
		args, err := catalog.Bind(st.Params, vals)
		if err != nil {
			fmt.Fprintf(&b, "\n=== %s\n-- not captured: %v\n", name, err)
			continue
		}
		fmt.Fprintf(&b, "\n=== %s  (params: %s)\n", name, strings.Join(st.Params, ", "))
		plan := inspect.ExplainPrefix(engine) + " " + st.SQL
		rows, err := db.Query(ctx, plan, args...)
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

// negativeControlFired reports whether a control actually broke its invariant. A
// control that did not fire means the other cells' clean audits are NOT evidence of
// their correctness.
func negativeControlFired(res *CellResult, d Design) error {
	for _, f := range res.Faults {
		if f.Reproduced {
			return nil
		}
	}
	switch d.Ver {
	case VerUnsafe:
		return fmt.Errorf("negative control %s did not fire: the unguarded version bump lost no updates", d.ID)
	default:
		return fmt.Errorf("negative control %s did not fire: the suppressed invalidation produced no detectable stale read", d.ID)
	}
}

// probeRedisCapabilities is the phase A dev check: every Redis primitive the cache
// policy depends on is exercised against the real server before any measurement.
// The result is written to out (or stdout) as JSON.
func probeRedisCapabilities(ctx context.Context, addr, out string) (map[string]any, error) {
	if addr == "" {
		return nil, fmt.Errorf("probe: -redis-addr is required")
	}
	st := newRedisStore(addr, 4)
	defer st.Close(ctx)
	res := map[string]any{"addr": addr}
	conn, err := dialRESP(addr, 5*time.Second)
	if err != nil {
		return nil, fmt.Errorf("dial: %w", err)
	}
	defer conn.Close()
	if v, err := conn.do("PING"); err != nil {
		return nil, fmt.Errorf("PING: %w", err)
	} else {
		res["ping"] = v
	}
	if info := st.Info(ctx); len(info) > 0 {
		res["info"] = info
	}

	// Lease primitives.
	const key = "probe:lease"
	_ = st.Delete(ctx, key)
	ok1, err := st.TryLease(ctx, key, "tok-1", time.Second)
	if err != nil {
		return nil, fmt.Errorf("first SET NX PX: %w", err)
	}
	ok2, err := st.TryLease(ctx, key, "tok-2", time.Second)
	if err != nil {
		return nil, fmt.Errorf("second SET NX PX: %w", err)
	}
	res["lease_acquired_first"] = ok1
	res["lease_refused_second"] = !ok2
	// A release with the WRONG token must fail, and with the right token succeed:
	// that is the compare-and-delete the lease depends on.
	_ = st.ReleaseLease(ctx, key, "tok-2")
	stillHeld, _, _ := st.Get(ctx, key)
	res["value_key_untouched_by_lease"] = stillHeld == nil
	if err := st.ReleaseLease(ctx, key, "tok-1"); err != nil {
		return nil, fmt.Errorf("release: %w", err)
	}
	reok, err := st.TryLease(ctx, key, "tok-3", time.Second)
	if err != nil {
		return nil, fmt.Errorf("re-acquire: %w", err)
	}
	res["lease_released_by_compare_and_delete"] = reok
	_ = st.ReleaseLease(ctx, key, "tok-3")

	// Version fencing on publish.
	vkey := "probe:value"
	_ = st.Delete(ctx, vkey)
	c1 := PortalContent{PersonID: 1, FullName: "a"}
	e5 := newEntry(c1, 5, 5, nowMS())
	if _, err := st.Put(ctx, vkey, e5, time.Minute); err != nil {
		return nil, fmt.Errorf("publish seq 5: %w", err)
	}
	e3 := newEntry(c1, 3, 3, nowMS())
	fenced, err := st.Put(ctx, vkey, e3, time.Minute)
	if err != nil {
		return nil, fmt.Errorf("publish seq 3: %w", err)
	}
	res["older_publication_refused"] = !fenced
	e6 := newEntry(PortalContent{PersonID: 1, FullName: "b"}, 6, 6, nowMS())
	newer, err := st.Put(ctx, vkey, e6, time.Minute)
	if err != nil {
		return nil, fmt.Errorf("publish seq 6: %w", err)
	}
	res["newer_publication_accepted"] = newer
	if got, hit, _ := st.Get(ctx, vkey); hit {
		res["round_trip_hash_matches"] = got.Hash == e6.Hash
	} else {
		res["round_trip_hash_matches"] = false
	}
	_ = st.Delete(ctx, vkey)

	// Invalidation fence: a fill that began before an invalidation must not be able
	// to publish what it read, and one that began after it must publish normally.
	fkey := "probe:fence"
	_ = st.Delete(ctx, fkey)
	f0, _ := st.FenceOf(ctx, fkey)
	pre := newEntry(PortalContent{PersonID: 2, FullName: "pre"}, 1, 1, nowMS())
	pre.Fence = f0
	if _, err := st.Put(ctx, fkey, pre, time.Minute); err != nil {
		return nil, fmt.Errorf("fence probe publish: %w", err)
	}
	if _, err := st.Fence(ctx, fkey); err != nil {
		return nil, fmt.Errorf("fence probe fence: %w", err)
	}
	refused, err := st.Put(ctx, fkey, pre, time.Minute)
	if err != nil {
		return nil, fmt.Errorf("fence probe republish: %w", err)
	}
	res["fill_started_before_invalidation_is_refused"] = !refused
	f1, err := st.FenceOf(ctx, fkey)
	if err != nil {
		return nil, fmt.Errorf("fence probe fenceof: %w", err)
	}
	post := newEntry(PortalContent{PersonID: 2, FullName: "post"}, 2, 2, nowMS())
	post.Fence = f1
	accepted, err := st.Put(ctx, fkey, post, time.Minute)
	if err != nil {
		return nil, fmt.Errorf("fence probe post publish: %w", err)
	}
	res["fill_after_invalidation_is_accepted"] = accepted
	_ = st.Delete(ctx, fkey)
	res["stats"] = st.Stats()

	// Hard expiry: a 1 ms TTL must be gone immediately afterwards.
	ekey := "probe:expiry"
	if _, err := st.Put(ctx, ekey, e6, time.Millisecond); err != nil {
		return nil, fmt.Errorf("publish short-ttl: %w", err)
	}
	time.Sleep(30 * time.Millisecond)
	_, hit, _ := st.Get(ctx, ekey)
	res["hard_expiry_enforced_by_server"] = !hit

	b, _ := json.MarshalIndent(res, "", "  ")
	if out != "" {
		if err := os.WriteFile(out, append(b, '\n'), 0o644); err != nil {
			return res, err
		}
	} else {
		fmt.Println(string(b))
	}
	// A probe that did not observe every required behaviour is a failure, so it
	// cannot be quietly skimmed.
	for _, k := range []string{"lease_acquired_first", "lease_refused_second", "lease_released_by_compare_and_delete",
		"older_publication_refused", "newer_publication_accepted", "round_trip_hash_matches",
		"hard_expiry_enforced_by_server", "fill_started_before_invalidation_is_refused",
		"fill_after_invalidation_is_accepted"} {
		if v, ok := res[k].(bool); !ok || !v {
			return res, fmt.Errorf("probe: %s was not satisfied (%v)", k, res[k])
		}
	}
	return res, nil
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, "fatal:", err)
	os.Exit(1)
}
