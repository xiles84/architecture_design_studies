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

// The join/filter and CTE-materialization measurement mode (`-cmd select-join`).
//
// It reuses the additive SELECT package's join catalogue, independent Go
// oracles and correctness gate. The gate runs on a small edge fixture (duplicate
// matches, NULL amounts, an orphan fact, an unmatched parent); only after it
// passes is the large deterministic fixture loaded and measured. No cell is
// timed unless its answer matched the oracle.
//
// A forced MATERIALIZED cell whose plan shows no materialization is recorded as
// unsupported rather than silently reported as a default fold. Unsupported
// cells are explicit results, not failures; invalid, timed-out and errored
// cells are never ranked.

func selectJoinSpec(scale string) selectbench.JoinGenSpec {
	switch scale {
	case "cte-small":
		return selectbench.JoinGenSpec{Customers: 200, FactCount: 20000, NullFKEvery: 199}
	default: // cte
		return selectbench.JoinGenSpec{Customers: 2000, FactCount: 200000, NullFKEvery: 1999}
	}
}

func prepareJoinFixture(ctx context.Context, q selectbench.Querier, jf selectbench.JoinFixture) error {
	if err := q.Exec(ctx, selectbench.JoinSchemaSQL()); err != nil {
		return err
	}
	for _, stmt := range selectbench.JoinLoadSQL(jf, 1000) {
		if err := q.Exec(ctx, stmt); err != nil {
			return err
		}
	}
	return q.Exec(ctx, "ANALYZE "+selectbench.JoinSchema+".customers; ANALYZE "+selectbench.JoinSchema+".facts;")
}

// runJoinGate is the correctness gate: the base spellings and the optional-match
// contract must answer the edge fixture correctly, and the WHERE-placement
// control must be seen to answer differently. A gate failure aborts before any
// timing.
func runJoinGate(ctx context.Context, q selectbench.Querier, cat *selectbench.Catalog) ([]selectbench.Cell, []selectbench.JoinControlResult, error) {
	edge := selectbench.JoinEdge()
	if err := prepareJoinFixture(ctx, q, edge); err != nil {
		return nil, nil, fmt.Errorf("gate fixture load: %w", err)
	}
	budget := 60 * time.Second
	base := selectbench.JoinArgs{Threshold: 5, Trials: 1, ValidateOnly: true, StatementTimeout: budget}
	var specs []selectbench.JoinSpec
	for _, name := range selectbench.JoinStatements {
		specs = append(specs, selectbench.JoinSpec{Statement: name, Axis: "gate=" + name, Args: base})
	}
	opt := base
	opt.Threshold = 1
	specs = append(specs, selectbench.JoinSpec{Statement: "j11_optional_on", Axis: "gate=j11_optional_on;threshold=1", Args: opt})
	cells := selectbench.RunJoinSpecs(ctx, q, cat, edge, specs, "gate", "validation")
	for _, c := range cells {
		switch c.Status {
		case selectbench.StatusValid, selectbench.StatusUnsupported:
			// Unsupported policy spellings are recorded by the feasibility
			// probe; the six defaults and the optional contract must be valid.
		default:
			return cells, nil, fmt.Errorf("gate cell %s is %s: %s", c.Statement, c.Status, c.Note)
		}
	}
	for _, name := range []string{"j01_join_where", "j02_join_on", "j03_filtered_derived",
		"j04_cte_filtered", "j05_cte_joined_outer", "j06_comma_join"} {
		found := false
		for _, c := range cells {
			if c.Statement == name && c.Status == selectbench.StatusValid {
				found = true
			}
		}
		if !found {
			return cells, nil, fmt.Errorf("gate: default spelling %s did not answer the edge fixture", name)
		}
	}
	controls := selectbench.RunJoinControls(ctx, q, cat, edge, []selectbench.JoinSpec{{
		Statement: "x01_optional_where", Axis: "control=x01_optional_where",
		Args: selectbench.JoinArgs{Threshold: 1, Trials: 1, StatementTimeout: budget},
	}})
	for _, c := range controls {
		if !c.Detected {
			return cells, controls, fmt.Errorf("gate control %s did not fire on the edge fixture: %s", c.Statement, c.Error)
		}
	}
	return cells, controls, nil
}

// joinFeasibility records the per-engine facts the handoff requires before the
// matrix runs: exact threshold counts, statement support, materialization
// evidence, rank-index use, expensive-token support and a full drain at the
// widest threshold.
func joinFeasibility(ctx context.Context, pool *pgxpool.Pool, q selectbench.Querier, cat *selectbench.Catalog,
	jf selectbench.JoinFixture, thresholds []int, budget time.Duration) *selectbench.JoinFeasibility {
	f := &selectbench.JoinFeasibility{
		FactCount:        len(jf.Facts),
		Customers:        len(jf.Customers),
		Thresholds:       thresholds,
		ThresholdCounts:  map[string]int{},
		StatementSupport: map[string]string{},
		PolicyEvidence:   map[string]string{},
		RankIndexUsed:    map[string]bool{},
	}
	for _, t := range thresholds {
		if n, err := scalarInt64(ctx, pool, fmt.Sprintf("SELECT count(*) FROM %s.facts WHERE rank <= %d", selectbench.JoinSchema, t)); err == nil {
			f.ThresholdCounts[fmt.Sprint(t)] = int(n)
		} else {
			f.ThresholdCounts[fmt.Sprint(t)] = -1
			f.Notes = append(f.Notes, fmt.Sprintf("threshold %d count failed: %v", t, err))
		}
	}
	var tok int
	if err := pool.QueryRow(ctx, "SELECT ascii(substr(md5(repeat('x'::text, 32)), 1, 1))").Scan(&tok); err != nil {
		f.ReuseExprError = err.Error()
	} else {
		f.ReuseExprSupported = true
	}

	probeThreshold := thresholds[2] // 1%
	probeArgs := selectbench.JoinArgs{Threshold: probeThreshold, Refs: 4, ConsumerMax: len(jf.Facts), ValidateOnly: true, StatementTimeout: budget}
	names := append(append([]string{}, selectbench.JoinStatements...), "j11_optional_on")
	for _, name := range names {
		specs := []selectbench.JoinSpec{{Statement: name, Axis: "probe=" + name, Args: probeArgs}}
		cells := selectbench.RunJoinSpecs(ctx, q, cat, jf, specs, "feasibility", "validation")
		c := cells[0]
		f.StatementSupport[name] = string(c.Status)
		f.PolicyEvidence[name] = c.PolicyNote
		f.RankIndexUsed[name] = strings.Contains(c.Plan, "facts_rank_idx")
		if c.Status != selectbench.StatusValid && c.Status != selectbench.StatusUnsupported {
			f.Notes = append(f.Notes, fmt.Sprintf("probe %s is %s: %s", name, c.Status, c.Note))
		}
	}
	for _, policy := range []string{"default", "mat", "notmat"} {
		name := selectbench.ReuseStatementName(policy, 4, 100)
		specs := []selectbench.JoinSpec{{Statement: name, Axis: "probe=" + name, Args: probeArgs}}
		cells := selectbench.RunJoinSpecs(ctx, q, cat, jf, specs, "feasibility", "validation")
		c := cells[0]
		f.StatementSupport[name] = string(c.Status)
		f.PolicyEvidence[name] = c.PolicyNote
	}

	// Full drain at the widest threshold: proves the client can consume the
	// broad-scan answer, not just the selective one.
	maxArgs := selectbench.JoinArgs{Threshold: thresholds[len(thresholds)-1], Trials: 1, StatementTimeout: budget}
	cells := selectbench.RunJoinSpecs(ctx, q, cat, jf, []selectbench.JoinSpec{
		{Statement: "j01_join_where", Axis: "probe=max_threshold", Args: maxArgs},
	}, "feasibility", "validation")
	if cells[0].Status == selectbench.StatusValid {
		f.MaxThresholdDrained = cells[0].ResultRows
	} else {
		f.Notes = append(f.Notes, fmt.Sprintf("max-threshold drain cell is %s: %s", cells[0].Status, cells[0].Note))
	}
	return f
}

func runSelectJoinMode(dsn, engine, topology, scale, out, reportOut string, seed int64, trials int, environment, runID string, statementTimeout int, positionControl, feasibilityOnly bool) error {
	if dsn == "" {
		return fmt.Errorf("-cmd select-join needs -dsn")
	}
	for _, key := range []string{"ADS_RUN_LOCK_HELD", "ADS_REPO_COMMIT", "ADS_RUN_TAG", "ADS_BENCH_IMAGE_ID", "ADS_ENGINE_IMAGE_ID"} {
		if strings.TrimSpace(os.Getenv(key)) == "" {
			return fmt.Errorf("SELECT join requires %s from the locked, tagged runner", key)
		}
	}
	if environment == "" || environment == "unknown" || runID == "" || trials < 1 || statementTimeout < 1 {
		return fmt.Errorf("SELECT join requires a named environment/run and positive trials/statement budget")
	}
	cat, err := selectbench.LoadJoin()
	if err != nil {
		return err
	}
	spec := selectJoinSpec(scale)
	jf := selectbench.GenerateJoin(seed, spec)
	thresholds := selectbench.JoinThresholds(spec.FactCount)
	budget := time.Duration(statementTimeout) * time.Millisecond

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Minute)
	defer cancel()

	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return fmt.Errorf("dsn: %w", err)
	}
	cfg.ConnConfig.RuntimeParams["search_path"] = selectbench.JoinSchema + ",public"
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

	result := selectbench.JoinResult{
		RunID: runID, Engine: engine, Topology: topology, Scale: scale,
		FactCount: spec.FactCount, Customers: spec.Customers, Seed: seed,
		Thresholds: thresholds, Summary: map[string]int{},
	}

	gateCells, gateControls, gateErr := runJoinGate(ctx, q, cat)
	_ = gateCells
	if gateErr != nil {
		return writeJoinOutputs(result, gateErr, false, out, reportOut)
	}

	if err := prepareJoinFixture(ctx, q, jf); err != nil {
		return fmt.Errorf("load fixture: %w", err)
	}
	result.Controls = append(result.Controls, gateControls...)

	// The limit-before-join control is proved on the large fixture (the edge
	// fixture is too small for a limit of 7 to truncate anything).
	bigControls := selectbench.RunJoinControls(ctx, q, cat, jf, []selectbench.JoinSpec{{
		Statement: "x02_join_limit_before", Axis: "control=x02_join_limit_before",
		Args: selectbench.JoinArgs{Threshold: int(float64(spec.FactCount) * 10 / 100), Trials: 1, StatementTimeout: budget},
	}})
	for _, c := range bigControls {
		if !c.Detected {
			return writeJoinOutputs(result, fmt.Errorf("limit-before-join control did not fire on the large fixture: %s", c.Error), feasibilityOnly, out, reportOut)
		}
	}
	result.Controls = append(result.Controls, bigControls...)

	result.Feasibility = joinFeasibility(ctx, pool, q, cat, jf, thresholds, budget)

	base := selectbench.JoinArgs{Trials: trials, StatementTimeout: budget}
	if !feasibilityOnly {
		specs := selectbench.BuildJoinSpecs(base, spec.FactCount, trials)
		plan := selectbench.PlanOrder(selectbench.OrderSeedFor(seed, runID), len(specs))
		first := specs[plan.Order[0]]
		result.Cells = selectbench.RunJoinSpecsInterleaved(ctx, q, cat, jf, specs, plan, runID, topology)
		result.Order = selectbench.RunOrder{Seed: plan.Seed, Rounds: trials, Sequence: plan.JoinSequence(specs)}
		if positionControl {
			pc := first
			pc.Axis += selectbench.PositionSuffix
			result.Cells = append(result.Cells, selectbench.RunJoinSpecs(ctx, q, cat, jf, []selectbench.JoinSpec{pc}, runID, topology)...)
			result.Order.Sequence = append(result.Order.Sequence, pc.Axis)
		}
	}

	version := ""
	_ = pool.QueryRow(ctx, "SELECT version()").Scan(&version)
	result.Provenance = joinProvenance(version, environment, scale, spec, thresholds, engine, topology, seed, trials, statementTimeout, cpuBefore, selectbench.JoinFingerprint(jf))
	result.Notes = append(result.Notes, joinGaps(engine, topology)...)
	return writeJoinOutputs(result, nil, feasibilityOnly, out, reportOut)
}

func joinProvenance(version, environment, scale string, spec selectbench.JoinGenSpec, thresholds []int, engine, topology string, seed int64, trials, statementTimeout int, cpuBefore, fingerprint string) selectbench.Provenance {
	rpc := "not parsed; raw EXPLAIN retained (coverage gap)"
	if engine == "postgres" && topology == "single" {
		rpc = "N/A on PostgreSQL single-node"
	}
	settings := map[string]string{
		"scale": scale, "topology": topology, "fact_count": fmt.Sprint(spec.FactCount),
		"customers": fmt.Sprint(spec.Customers), "thresholds": fmt.Sprint(thresholds),
		"selectivity_mapping":  "rank <= threshold; rank is a deterministic shuffle of 1..fact_count",
		"index_set":            "facts_rank_idx(rank),facts_customer_idx(customer_id),primary keys; no per-spelling index",
		"statistics":           "ANALYZE after deterministic load",
		"prepared_plan_policy": "pgx simple protocol; no prepared statements",
		"client_workers":       "1 sequential reader; 0 writers", "connection_pool": "2 = 1 reader/setup + 1 reserve",
		"trials_per_cell": fmt.Sprint(trials), "statement_budget_ms": fmt.Sprint(statementTimeout),
		"whole_run_count": "1; no replicated winner in this file", "distributed_rpc": rpc,
	}
	return selectbench.Provenance{
		RepoCommit:          envOr("ADS_REPO_COMMIT", "unknown"),
		EngineImage:         envOr("ADS_BENCH_IMAGE", "unknown"),
		ImageDigest:         os.Getenv("ADS_BENCH_IMAGE_ID"),
		EngineImageDigest:   os.Getenv("ADS_ENGINE_IMAGE_ID"),
		EngineVersion:       version,
		Environment:         environment,
		RunTag:              os.Getenv("ADS_RUN_TAG"),
		Seed:                seed,
		DatasetFingerprint:  fingerprint,
		Settings:            settings,
		ClientCPUStatBefore: cpuBefore,
		ClientCPUStatAfter:  readSelectCPUStat(),
		ClientCPUQuota:      readSelectResource("/sys/fs/cgroup/cpu.max"),
		ClientMemoryLimit:   readSelectResource("/sys/fs/cgroup/memory.max"),
	}
}

func joinGaps(engine, topology string) []string {
	return []string{
		"one host, single client: no independent hardware, real-network or equal-total-budget comparison",
		"no concurrent writers: the isolated read contract only; a matched blended read/write phase is a recorded gap",
		"a flat cell reports a client or engine ceiling; it is not by itself a constant query cost",
	}
}

func writeJoinOutputs(result selectbench.JoinResult, runErr error, feasibilityOnly bool, out, reportOut string) error {
	if runErr != nil {
		result.Error = runErr.Error()
	}
	for _, c := range result.Cells {
		result.Summary[string(c.Status)]++
	}
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	inputsDigest := fmt.Sprintf("%x", sha256.Sum256(data))
	if out != "" {
		if err := os.WriteFile(out, data, 0o644); err != nil {
			return err
		}
		fmt.Println("wrote", out)
	}
	if reportOut != "" {
		report, err := selectbench.RenderJoinRun(result, inputsDigest)
		if err != nil {
			return err
		}
		if err := os.WriteFile(reportOut, []byte(report), 0o644); err != nil {
			return err
		}
		fmt.Println("wrote", reportOut)
	}
	fmt.Printf("join: %d cells — valid=%d invalid=%d unsupported=%d timeout=%d error=%d; inputs digest %s\n",
		len(result.Cells), result.Summary["valid"], result.Summary["invalid"],
		result.Summary["unsupported"], result.Summary["timeout"], result.Summary["error"], inputsDigest)
	if result.Error != "" {
		return fmt.Errorf("%s", result.Error)
	}
	fatal := result.Summary["invalid"] + result.Summary["timeout"] + result.Summary["error"]
	if fatal > 0 && !feasibilityOnly {
		return fmt.Errorf("SELECT join has %d invalid/timeout/error cells; inspect their recorded status", fatal)
	}
	return nil
}

// runSelectJoinReport re-renders one run from a directory holding join.json, or
// writes the cross-whole-run report for a directory of run subdirectories.
func runSelectJoinReport(resultsDir, reportOut string) error {
	if resultsDir == "" || reportOut == "" {
		return fmt.Errorf("-cmd join-report needs -results <dir> and -report-out <file>")
	}
	if fi, err := os.Stat(filepath.Join(resultsDir, "join.json")); err == nil && !fi.IsDir() {
		data, err := os.ReadFile(filepath.Join(resultsDir, "join.json"))
		if err != nil {
			return err
		}
		var r selectbench.JoinResult
		if err := json.Unmarshal(data, &r); err != nil {
			return err
		}
		digest := fmt.Sprintf("%x", sha256.Sum256(append(data, '\n')))
		report, err := selectbench.RenderJoinRun(r, digest)
		if err != nil {
			return err
		}
		return os.WriteFile(reportOut, []byte(report), 0o644)
	}
	paths, err := filepath.Glob(filepath.Join(resultsDir, "*", "join.json"))
	if err != nil {
		return err
	}
	if len(paths) == 0 {
		return fmt.Errorf("no join.json at or under %s", resultsDir)
	}
	sort.Strings(paths)
	var runs []selectbench.JoinResult
	for _, p := range paths {
		data, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		var r selectbench.JoinResult
		if err := json.Unmarshal(data, &r); err != nil {
			return fmt.Errorf("%s: %w", p, err)
		}
		// A feasibility-only run measures probes, not the matrix; it has no
		// cells to replicate and must not appear as a whole run.
		if len(r.Cells) == 0 {
			continue
		}
		runs = append(runs, r)
	}
	if len(runs) == 0 {
		return fmt.Errorf("no whole join runs (with cells) at or under %s", resultsDir)
	}
	report, err := selectbench.RenderJoinReplicated(runs)
	if err != nil {
		return err
	}
	return os.WriteFile(reportOut, []byte(report), 0o644)
}
