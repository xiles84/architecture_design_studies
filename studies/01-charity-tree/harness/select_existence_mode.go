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

// Membership / absence mode (`-cmd select-existence`).
//
// Fanout is a fixture axis: each level loads the same 10,000 child rows with a
// different children-per-parent count, then runs the same statements with the
// same predicate. The oracle gate is the same as the rest of the programme.

func runSelectExistenceMode(dsn, engine, topology, scale, out, reportOut string, seed int64, trials int, environment, runID string, statementTimeout int, positionControl bool) error {
	if dsn == "" {
		return fmt.Errorf("-cmd select-existence needs -dsn")
	}
	for _, key := range []string{"ADS_RUN_LOCK_HELD", "ADS_REPO_COMMIT", "ADS_RUN_TAG", "ADS_BENCH_IMAGE_ID", "ADS_ENGINE_IMAGE_ID"} {
		if strings.TrimSpace(os.Getenv(key)) == "" {
			return fmt.Errorf("select-existence requires %s from the locked, tagged runner", key)
		}
	}
	if environment == "" || environment == "unknown" || runID == "" || trials < 1 || statementTimeout < 1 {
		return fmt.Errorf("select-existence requires a named environment/run and positive trials/statement budget")
	}
	if scale != "fanout" {
		return fmt.Errorf("select-existence needs --scale fanout; got %q", scale)
	}
	cat, err := selectbench.Load()
	if err != nil {
		return err
	}

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
	_, _, gateErr := validateSelectFixtures(ctx, q, cat)
	if gateErr != nil {
		return fmt.Errorf("live teaching-fixture gate failed: %w", gateErr)
	}

	base := selectbench.Args{
		Pred:             selectbench.AmountGreaterThan(100),
		PredSQL:          "amount > 100",
		Trials:           trials,
		StatementTimeout: time.Duration(statementTimeout) * time.Millisecond,
	}
	controlArgs := selectbench.Args{Pred: selectbench.AmountGreaterThan(100), PredSQL: "amount > 100",
		N: 1, Trials: 1, StatementTimeout: time.Duration(statementTimeout) * time.Millisecond}

	var cells []selectbench.Cell
	var levels []selectbench.FanoutLevel
	var controls []selectbench.ControlResult
	for i, fanout := range selectbench.FanoutLevels() {
		fx := selectbench.Generate(seed, selectbench.GenSpec{
			Customers:       selectbench.TotalChildren / fanout,
			OrdersPerParent: fanout,
			NullAmountEvery: 11,
			NullFKEvery:     97,
		})
		if err := selectbench.Prepare(ctx, q, fx); err != nil {
			return fmt.Errorf("load fanout=%d fixture: %w", fanout, err)
		}
		specs := selectbench.BuildExistenceSpecs(base, fanout, trials)
		if positionControl && i == 0 {
			specs = append(specs, selectbench.ExistencePositionControlSpec(base, trials))
		}
		cells = append(cells, selectbench.RunSpecs(ctx, q, cat, fx, specs, runID, topology)...)

		levelControls := selectbench.RunControls(ctx, q, cat, fx, controlArgs)
		found := map[string]bool{}
		for _, c := range levelControls {
			if c.Statement == "x02_not_in_nullable" || c.Statement == "x07_membership_multiplier" {
				controls = append(controls, c)
				found[c.Statement] = true
				if !c.Detected {
					return fmt.Errorf("existence control %s did not fire at fanout=%d: %s", c.Statement, fanout, c.Error)
				}
			}
		}
		if len(found) != 2 {
			return fmt.Errorf("expected both existence controls at fanout=%d, got %d", fanout, len(found))
		}
		levels = append(levels, selectbench.FanoutLevel{Fanout: fanout, Parents: len(fx.Customers),
			Children: len(fx.Orders), Fingerprint: selectbench.Fingerprint(fx)})
	}

	var version string
	if err := pool.QueryRow(ctx, "SELECT version()").Scan(&version); err != nil {
		return err
	}
	settings := map[string]string{"scale": scale, "topology": topology,
		"index_set":           "orders_amount_idx,orders_customer_idx,primary keys",
		"statistics":          "ANALYZE after each deterministic load",
		"client_workers":      "1 sequential reader; 0 writers",
		"connection_pool":     "2 = 1 reader/setup + 1 reserve",
		"statement_budget_ms": fmt.Sprint(statementTimeout),
		"fanout_levels":       "1,10,100 with 10,000 child rows held constant",
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

	combined := make([]string, 0, len(levels))
	for _, l := range levels {
		combined = append(combined, fmt.Sprintf("%d:%s", l.Fanout, l.Fingerprint))
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
		DatasetFingerprint:  strings.Join(combined, ","),
		Settings:            settings,
		ClientCPUStatAfter:  readSelectCPUStat(),
		ClientCPUStatBefore: cpuBefore,
		ClientCPUQuota:      readSelectResource("/sys/fs/cgroup/cpu.max"),
		ClientMemoryLimit:   readSelectResource("/sys/fs/cgroup/memory.max"),
	}

	result := selectbench.ExistenceResult{
		RunID: runID, Engine: engine, Topology: topology, Scale: scale, Seed: seed,
		Levels: levels, Provenance: prov, Cells: cells, Controls: controls, Summary: map[string]int{},
		Notes: []string{
			"one sequential reader, no concurrent writers; latency of a drained answer, not throughput",
			"single-node only; the 3-node per-node/equal-total arms are a coverage gap on this host",
			"qualifying proportion, parent hit fraction and hot-parent skew are not swept in this bounded run",
			"YugabyteDB storage/RPC counters are not parsed; raw EXPLAIN is retained",
		},
	}
	for _, c := range cells {
		result.Summary[string(c.Status)]++
	}

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
		report, err := selectbench.RenderExistenceRun(result)
		if err != nil {
			return err
		}
		if err := os.WriteFile(reportOut, []byte(report), 0o644); err != nil {
			return err
		}
		fmt.Println("wrote", reportOut)
	}
	fmt.Printf("select-existence: %d cells — valid=%d invalid=%d unsupported=%d timeout=%d; inputs digest %s\n",
		len(cells), result.Summary["valid"], result.Summary["invalid"],
		result.Summary["unsupported"], result.Summary["timeout"], digest)
	if bad := result.Summary["invalid"] + result.Summary["error"] + result.Summary["timeout"] + result.Summary["unsupported"]; bad > 0 {
		return fmt.Errorf("select-existence has %d non-valid cells; inspect their recorded status", bad)
	}
	return nil
}

// runSelectExistenceReport renders one run from <dir>/existence.json or a
// replicated summary across run subdirectories.
func runSelectExistenceReport(dir, reportOut string) error {
	if dir == "" {
		return fmt.Errorf("-cmd existence-report needs -results")
	}
	load := func(path string) (selectbench.ExistenceResult, error) {
		var r selectbench.ExistenceResult
		data, err := os.ReadFile(path)
		if err != nil {
			return r, err
		}
		return r, json.Unmarshal(data, &r)
	}
	direct := filepath.Join(dir, "existence.json")
	var runs []selectbench.ExistenceResult
	if _, err := os.Stat(direct); err == nil {
		r, err := load(direct)
		if err != nil {
			return err
		}
		runs = append(runs, r)
	} else {
		paths, err := filepath.Glob(filepath.Join(dir, "*", "existence.json"))
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
		return fmt.Errorf("no existence.json under %s", dir)
	}
	var report string
	var err error
	if len(runs) == 1 {
		report, err = selectbench.RenderExistenceRun(runs[0])
	} else {
		report, err = selectbench.RenderExistenceReplicated(runs)
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
