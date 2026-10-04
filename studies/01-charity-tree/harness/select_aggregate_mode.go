package main

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"charitytree/selectbench"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Aggregate / top-N sweep mode (`-cmd select-aggregate`).
//
// It reuses the selectbench package's oracles, controls, plan signatures and
// report renderers; here it only supplies the pgx-backed Querier, loads the
// deterministic fixture, runs the parent-level matrix and writes the
// measurements. The correctness gate is the same as the rest of the programme:
// a cell whose answer differs from the oracle is invalid and is never timed.

func runSelectAggregateMode(dsn, engine, topology, scale, out, reportOut string, seed int64, trials int, environment, runID string, statementTimeout int, positionControl bool) error {
	if dsn == "" {
		return fmt.Errorf("-cmd select-aggregate needs -dsn")
	}
	for _, key := range []string{"ADS_RUN_LOCK_HELD", "ADS_REPO_COMMIT", "ADS_RUN_TAG", "ADS_BENCH_IMAGE_ID", "ADS_ENGINE_IMAGE_ID"} {
		if strings.TrimSpace(os.Getenv(key)) == "" {
			return fmt.Errorf("select-aggregate requires %s from the locked, tagged runner", key)
		}
	}
	if environment == "" || environment == "unknown" || runID == "" || trials < 1 || statementTimeout < 1 {
		return fmt.Errorf("select-aggregate requires a named environment/run and positive trials/statement budget")
	}
	if scale != "sweep" {
		return fmt.Errorf("select-aggregate needs --scale sweep (10,000 parents); got %q", scale)
	}
	cat, err := selectbench.Load()
	if err != nil {
		return err
	}

	// One deterministic fixture; GuaranteedTie makes the top of the smallest
	// multi-child parent a tie so the x06 control always has a tie to keep.
	fx := selectbench.GuaranteedTie(selectbench.Generate(seed, selectSpec(scale)))

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Hour)
	defer cancel()

	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return fmt.Errorf("dsn: %w", err)
	}
	cfg.ConnConfig.RuntimeParams["search_path"] = selectbench.TeachingSchema + ",public"
	cfg.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol
	cfg.ConnConfig.RuntimeParams["statement_timeout"] = fmt.Sprint(statementTimeout)
	cfg.MaxConns = 2
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return fmt.Errorf("connect: %w", err)
	}
	defer pool.Close()

	q := pgxQuerier{pool: pool}
	cpuBefore := readSelectCPUStat()
	validation, _, gateErr := validateSelectFixtures(ctx, q, cat)
	if gateErr != nil {
		return fmt.Errorf("live teaching-fixture gate failed: %w", gateErr)
	}
	if err := selectbench.Prepare(ctx, q, fx); err != nil {
		return fmt.Errorf("load fixture: %w", err)
	}

	base := selectbench.Args{
		Trials:           trials,
		StatementTimeout: time.Duration(statementTimeout) * time.Millisecond,
	}
	var cells []selectbench.Cell
	if gateErr == nil {
		specs := selectbench.BuildAggregateSpecs(base, selectbench.ParentLevels(), trials)
		if positionControl {
			specs = append(specs, selectbench.AggregatePositionControlSpec(base, trials))
		}
		cells = selectbench.RunSpecs(ctx, q, cat, fx, specs, runID, topology)
	}

	// Controls for this family run at the widest level, N=1. RunControls also
	// runs the other catalogue controls, which need the predicate bound, so the
	// same predicate the teaching gate uses is supplied here.
	controlArgs := selectbench.Args{Pred: selectbench.AmountGreaterThan(100), PredSQL: "amount > 100",
		N: 1, ParentMax: 0, Trials: 1,
		StatementTimeout: time.Duration(statementTimeout) * time.Millisecond}
	controls := selectbench.RunControls(ctx, q, cat, fx, controlArgs)
	var familyControls []selectbench.ControlResult
	for _, c := range controls {
		if c.Statement == "x05_aggregate_double_count" || c.Statement == "x06_topn_rank_ties" {
			familyControls = append(familyControls, c)
		}
	}
	if len(familyControls) != 2 {
		return fmt.Errorf("expected the two aggregate controls, got %d", len(familyControls))
	}
	for _, c := range familyControls {
		if !c.Detected {
			return fmt.Errorf("aggregate control %s did not fire: %s", c.Statement, c.Error)
		}
	}

	var version string
	if err := pool.QueryRow(ctx, "SELECT version()").Scan(&version); err != nil {
		return err
	}
	settings := map[string]string{"scale": scale, "topology": topology,
		"index_set":           "orders_amount_idx,orders_customer_idx,primary keys",
		"statistics":          "ANALYZE after deterministic load",
		"client_workers":      "1 sequential reader; 0 writers",
		"connection_pool":     "2 = 1 reader/setup + 1 reserve",
		"statement_budget_ms": fmt.Sprint(statementTimeout),
		"levels":              "1,10,100,1000,10000 parents (the protocol's parent-selection levels plus 0.01% so a crossing has >=5 levels)",
		"topn_n":              "1 (DISTINCT ON is N=1-only)",
		"distributed_rpc":     "N/A on PostgreSQL single-node",
	}
	if engine != "postgres" || topology != "single" {
		settings["distributed_rpc"] = "not parsed; raw EXPLAIN retained (coverage gap)"
	}
	settingRows, err := pool.Query(ctx, "SELECT name, setting FROM pg_settings WHERE name IN ('shared_buffers','work_mem','effective_cache_size','max_connections','random_page_cost','track_io_timing','statement_timeout')")
	if err != nil {
		return err
	}
	for settingRows.Next() {
		var key, value string
		if err := settingRows.Scan(&key, &value); err != nil {
			settingRows.Close()
			return err
		}
		settings[key] = value
	}
	settingRows.Close()
	if err := settingRows.Err(); err != nil {
		return err
	}

	prov := selectbench.Provenance{
		RepoCommit:          envOr("ADS_REPO_COMMIT", "unknown"),
		EngineImage:         envOr("ADS_BENCH_IMAGE", "unknown"),
		ImageDigest:         os.Getenv("ADS_BENCH_IMAGE_ID"),
		EngineImageDigest:   os.Getenv("ADS_ENGINE_IMAGE_ID"),
		EngineVersion:       version,
		Environment:         environment,
		RunTag:              os.Getenv("ADS_RUN_TAG"),
		Seed:                seed,
		DatasetFingerprint:  selectbench.Fingerprint(fx),
		Settings:            settings,
		ClientCPUStatAfter:  readSelectCPUStat(),
		ClientCPUStatBefore: cpuBefore,
		ClientCPUQuota:      readSelectResource("/sys/fs/cgroup/cpu.max"),
		ClientMemoryLimit:   readSelectResource("/sys/fs/cgroup/memory.max"),
	}

	result := selectbench.AggResult{
		RunID: runID, Engine: engine, Topology: topology, Scale: scale, Seed: seed,
		Levels: selectbench.ParentLevels(), Provenance: prov, Cells: cells,
		Controls: familyControls, Summary: map[string]int{},
		Notes: []string{
			"one sequential reader, no concurrent writers; latency of a drained answer, not throughput",
			"single-node only; the 3-node per-node/equal-total arms are a coverage gap on this host",
			"YugabyteDB storage/RPC counters are not parsed; raw EXPLAIN is retained",
		},
	}
	for _, c := range cells {
		result.Summary[string(c.Status)]++
	}
	_ = validation

	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	digest := fmt.Sprintf("%x", sha256.Sum256(data))
	if out != "" {
		if err := os.WriteFile(out, data, 0o644); err != nil {
			return err
		}
		fmt.Println("wrote", out)
	}
	if reportOut != "" {
		report, err := selectbench.RenderAggRun(result)
		if err != nil {
			return err
		}
		if err := os.WriteFile(reportOut, []byte(report), 0o644); err != nil {
			return err
		}
		fmt.Println("wrote", reportOut)
	}
	fmt.Printf("select-aggregate: %d cells — valid=%d invalid=%d unsupported=%d timeout=%d; inputs digest %s\n",
		len(cells), result.Summary["valid"], result.Summary["invalid"],
		result.Summary["unsupported"], result.Summary["timeout"], digest)
	if bad := result.Summary["invalid"] + result.Summary["error"] + result.Summary["timeout"] + result.Summary["unsupported"]; bad > 0 {
		return fmt.Errorf("select-aggregate has %d non-valid cells; inspect their recorded status", bad)
	}
	return nil
}

// runSelectAggregateReport renders one run from <dir>/aggregate.json, or a
// replicated summary when <dir> contains run subdirectories.
func runSelectAggregateReport(dir, reportOut string) error {
	if dir == "" {
		return fmt.Errorf("-cmd aggregate-report needs -results")
	}
	load := func(path string) (selectbench.AggResult, error) {
		var r selectbench.AggResult
		data, err := os.ReadFile(path)
		if err != nil {
			return r, err
		}
		return r, json.Unmarshal(data, &r)
	}
	direct := filepath.Join(dir, "aggregate.json")
	var runs []selectbench.AggResult
	if _, err := os.Stat(direct); err == nil {
		r, err := load(direct)
		if err != nil {
			return err
		}
		runs = append(runs, r)
	} else {
		paths, err := filepath.Glob(filepath.Join(dir, "*", "aggregate.json"))
		if err != nil {
			return err
		}
		sort.Strings(paths)
		for _, p := range paths {
			r, err := load(p)
			if err != nil {
				return err
			}
			runs = append(runs, r)
		}
	}
	if len(runs) == 0 {
		return fmt.Errorf("no aggregate.json under %s", dir)
	}
	var report string
	var err error
	if len(runs) == 1 {
		report, err = selectbench.RenderAggRun(runs[0])
	} else {
		report, err = selectbench.RenderAggReplicated(runs)
	}
	if err != nil {
		return err
	}
	if reportOut == "" {
		fmt.Print(report)
		return nil
	}
	return os.WriteFile(reportOut, []byte(report), 0o644)
}
