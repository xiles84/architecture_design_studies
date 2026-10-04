package main

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"charitytree/selectbench"

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

// probeExplainOptions picks the richest EXPLAIN option list the pinned engine
// accepts: PostgreSQL VERBOSE/BUFFERS, YugabyteDB DIST/DEBUG. The chosen list is
// recorded; if the richest form is unsupported, the next is used and the
// fallback is reported as evidence rather than the absence of RPC counters being
// described as "not parsed" (review R3).
func probeExplainOptions(ctx context.Context, pool *pgxpool.Pool, engine string) (string, string) {
	candidates := []string{"ANALYZE, VERBOSE, BUFFERS, FORMAT TEXT", "ANALYZE, BUFFERS, FORMAT TEXT"}
	if engine == "yugabyte" {
		candidates = []string{"ANALYZE, DIST, DEBUG, FORMAT TEXT", "ANALYZE, DIST, FORMAT TEXT", "ANALYZE, BUFFERS, FORMAT TEXT"}
	}
	var firstErr string
	for _, c := range candidates {
		rows, err := pool.Query(ctx, "EXPLAIN ("+c+") SELECT 1")
		if err != nil {
			if firstErr == "" {
				firstErr = err.Error()
			}
			continue
		}
		for rows.Next() {
		}
		rerr := rows.Err()
		rows.Close()
		if rerr != nil {
			if firstErr == "" {
				firstErr = rerr.Error()
			}
			continue
		}
		return c, "explain options `" + c + "` supported on " + engine
	}
	return "", "no engine-specific EXPLAIN option list was supported; plain PostgreSQL options used (" + firstErr + ")"
}

func isJoinDomainScale(scale string) bool { return strings.Contains(scale, "domain") }

func selectJoinSpec(scale string) selectbench.JoinGenSpec {
	switch scale {
	case "cte-small":
		return selectbench.JoinGenSpec{Customers: 200, FactCount: 20000, NullFKEvery: 199}
	default: // cte
		return selectbench.JoinGenSpec{Customers: 2000, FactCount: 200000, NullFKEvery: 1999}
	}
}

// selectJoinDomainSpec gives the domain arm the same approved shape (2,000
// persons, 200,000 donations) so the reviewed 60-cell spelling and 30-cell reuse
// matrices are the same size, on real Study 01 entities.
func selectJoinDomainSpec(scale string) selectbench.JoinDomainGenSpec {
	switch scale {
	case "cte-domain-small":
		return selectbench.JoinDomainGenSpec{Charities: 20, Persons: 200, Donations: 20000, NullAmountEvery: 199}
	default: // cte-domain
		return selectbench.JoinDomainGenSpec{Charities: 50, Persons: 2000, Donations: 200000, NullAmountEvery: 1999}
	}
}

// joinArm is the data binding a run measures: the original synthetic
// microfixture (customers/facts) or the Study 01 domain (person/donation bound
// through the charity relationship). The catalogue is shared; only the fixture,
// the search_path schema and the provenance label differ.
type joinArm struct {
	name          string
	dataset       string
	schema        string
	rankIndex     string
	scale         string
	placement     string
	placementNote string
	indexArm      string
	indexArmNote  string
	factCount     int
	parents       int
	fingerprint   string
	oracle        selectbench.JoinFixture
	edgeOracle    selectbench.JoinFixture
	emptyOracle   selectbench.JoinFixture
	prepare       func(context.Context, selectbench.Querier) error
	prepareEdge   func(context.Context, selectbench.Querier) error
	prepareEmpty  func(context.Context, selectbench.Querier) error
}

func selectJoinArm(seed int64, scale, engine, placement, indexArm string) joinArm {
	if isJoinDomainScale(scale) {
		dspec := selectJoinDomainSpec(scale)
		df := selectbench.GenerateJoinDomain(seed, dspec)
		edge := selectbench.JoinDomainEdge()
		empty := selectbench.JoinDomainFixture{
			Charities: []selectbench.Charity{{ID: 1, Name: "Empty"}},
			Persons: []selectbench.Person{
				{ID: 1, CharityID: 1, FullName: "Ana"},
				{ID: 2, CharityID: 1, FullName: "Ben"},
			},
		}
		// Colocation is a YugabyteDB table property. Requesting it on any other
		// engine is recorded as not applicable, never silently applied.
		effPlacement, note := "", ""
		if placement != "" {
			if engine == "yugabyte" {
				effPlacement = placement
			} else {
				note = "placement `" + placement + "` is not applicable on " + engine + "; engine default used"
			}
		}
		// The rank-index arm is a physical design decision, not a syntax one
		// (review M1). An explicit arm the engine cannot express is recorded.
		effIndexArm := selectbench.JoinDomainDefaultIndexArm(engine, indexArm)
		indexNote := ""
		if indexArm != "" && !selectbench.JoinDomainIndexArmApplies(engine, indexArm) {
			indexNote = "index arm `" + indexArm + "` is not applicable on " + engine + "; " + effIndexArm + " used"
		}
		return joinArm{
			name: "study01-domain", dataset: "study01-person-donation-charity",
			schema: selectbench.JoinDomainSchema, rankIndex: "donation_rank_idx", scale: scale,
			placement: effPlacement, placementNote: note,
			indexArm: effIndexArm, indexArmNote: indexNote,
			factCount: dspec.Donations, parents: dspec.Persons,
			fingerprint: selectbench.JoinDomainFingerprint(df),
			oracle:      df.AsJoinFixture(),
			edgeOracle:  edge.AsJoinFixture(),
			emptyOracle: empty.AsJoinFixture(),
			prepare: func(ctx context.Context, q selectbench.Querier) error {
				return prepareJoinDomain(ctx, q, df, effPlacement, effIndexArm)
			},
			prepareEdge: func(ctx context.Context, q selectbench.Querier) error {
				return prepareJoinDomain(ctx, q, edge, effPlacement, effIndexArm)
			},
			prepareEmpty: func(ctx context.Context, q selectbench.Querier) error {
				return prepareJoinDomain(ctx, q, empty, effPlacement, effIndexArm)
			},
		}
	}
	spec := selectJoinSpec(scale)
	jf := selectbench.GenerateJoin(seed, spec)
	empty := selectbench.JoinFixture{Customers: []selectbench.Customer{{ID: 1, Name: "Ana"}, {ID: 2, Name: "Ben"}}}
	note := ""
	if placement != "" {
		note = "placement `" + placement + "` is not applicable to the synthetic microfixture arm; engine default used"
	}
	indexNote := ""
	if indexArm != "" {
		indexNote = "index arm `" + indexArm + "` is not applicable to the synthetic microfixture arm; " + selectbench.IndexArmPGBtree + "/engine default used"
	}
	return joinArm{
		name: "synthetic-microfixture", dataset: "synthetic-customers-facts",
		schema: selectbench.JoinSchema, rankIndex: "facts_rank_idx", scale: scale,
		placementNote: note, indexArmNote: indexNote,
		factCount: spec.FactCount, parents: spec.Customers,
		fingerprint: selectbench.JoinFingerprint(jf),
		oracle:      jf,
		edgeOracle:  selectbench.JoinEdge(),
		emptyOracle: empty,
		prepare:     func(ctx context.Context, q selectbench.Querier) error { return prepareJoinFixture(ctx, q, jf) },
		prepareEdge: func(ctx context.Context, q selectbench.Querier) error {
			return prepareJoinFixture(ctx, q, selectbench.JoinEdge())
		},
		prepareEmpty: func(ctx context.Context, q selectbench.Querier) error { return prepareJoinFixture(ctx, q, empty) },
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

// prepareJoinDomain loads the Study 01 domain tables and their catalogue-binding
// views. The views are the only place person/donation are mapped onto the shared
// projection, so the measured SQL never changes between arms.
func prepareJoinDomain(ctx context.Context, q selectbench.Querier, df selectbench.JoinDomainFixture, placement, indexArm string) error {
	if err := q.Exec(ctx, selectbench.JoinDomainSchemaSQLOptions(placement, indexArm)); err != nil {
		return err
	}
	for _, stmt := range selectbench.JoinDomainLoadSQL(df, 1000) {
		if err := q.Exec(ctx, stmt); err != nil {
			return err
		}
	}
	return q.Exec(ctx, selectbench.JoinDomainAnalyzeSQL())
}

// runJoinGate is the correctness gate: the base spellings and the optional-match
// contract must answer the edge fixture correctly, and the WHERE-placement
// control must be seen to answer differently. A gate failure aborts before any
// timing.
func runJoinGate(ctx context.Context, q selectbench.Querier, cat *selectbench.Catalog, arm joinArm) ([]selectbench.Cell, []selectbench.JoinControlResult, error) {
	edge := arm.edgeOracle
	if err := arm.prepareEdge(ctx, q); err != nil {
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
	specs = append(specs,
		selectbench.JoinSpec{Statement: "j11_optional_on", Axis: "gate=j11_optional_on;threshold=1", Args: opt},
		selectbench.JoinSpec{Statement: "j11_optional_derived", Axis: "gate=j11_optional_derived;threshold=1", Args: opt},
	)
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
	// Empty input: the bag is empty and the optional contract preserves every
	// parent with NULLs. This is a live DB check, not only the Go oracle, and it
	// covers both the ON and filtered-derived optional spellings.
	empty := arm.emptyOracle
	if err := arm.prepareEmpty(ctx, q); err != nil {
		return cells, controls, fmt.Errorf("gate empty-fixture load: %w", err)
	}
	emptyCells := selectbench.RunJoinSpecs(ctx, q, cat, empty, []selectbench.JoinSpec{
		{Statement: "j01_join_where", Axis: "gate=empty;j01_join_where", Args: base},
		{Statement: "j11_optional_on", Axis: "gate=empty;j11_optional_on", Args: base},
		{Statement: "j11_optional_derived", Axis: "gate=empty;j11_optional_derived", Args: base},
	}, "gate", "validation")
	for _, c := range emptyCells {
		if c.Status != selectbench.StatusValid {
			return cells, controls, fmt.Errorf("empty-input gate cell %s is %s: %s", c.Statement, c.Status, c.Note)
		}
	}
	return cells, controls, nil
}

// joinFeasibility records the per-engine facts the handoff requires before the
// matrix runs: exact threshold counts, statement support, materialization
// evidence, rank-index use, expensive-token support and a full drain at the
// widest threshold.
func joinFeasibility(ctx context.Context, pool *pgxpool.Pool, q selectbench.Querier, cat *selectbench.Catalog,
	jf selectbench.JoinFixture, thresholds []int, budget time.Duration, schema, rankIndex string) *selectbench.JoinFeasibility {
	f := &selectbench.JoinFeasibility{
		FactCount:         len(jf.Facts),
		Customers:         len(jf.Customers),
		Thresholds:        thresholds,
		ThresholdCounts:   map[string]int{},
		StatementSupport:  map[string]string{},
		PolicyEvidence:    map[string]string{},
		RankIndexUsed:     map[string]bool{},
		DefaultIndexGate:  map[string]string{},
		DefaultIndexPlans: map[string]string{},
	}
	for _, t := range thresholds {
		if n, err := scalarInt64(ctx, pool, fmt.Sprintf("SELECT count(*) FROM %s.facts WHERE rank <= %d", schema, t)); err == nil {
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
		// Volatility evidence: the token must be a per-row expression, not a
		// constant the planner can hoist. Distinct ids must yield distinct
		// tokens, and the reuse plan must show the expression.
		var t1, t2 int
		e1 := pool.QueryRow(ctx, "SELECT ascii(substr(md5(repeat(1::text, 32)), 1, 1))").Scan(&t1)
		e2 := pool.QueryRow(ctx, "SELECT ascii(substr(md5(repeat(2::text, 32)), 1, 1))").Scan(&t2)
		if e1 == nil && e2 == nil {
			f.TokenEvidence = fmt.Sprintf("token(id=1)=%d token(id=2)=%d distinct=%v", t1, t2, t1 != t2)
		} else if e1 != nil {
			f.TokenEvidence = "token evaluation failed: " + e1.Error()
		}
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
		f.RankIndexUsed[name] = strings.Contains(c.Plan, rankIndex) || strings.Contains(c.Plan, "rank_idx")
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
		if strings.Contains(strings.ToLower(c.Plan), "md5") && strings.Contains(strings.ToLower(c.Plan), "repeat") {
			f.TokenEvidence += fmt.Sprintf("; %s plan carries the token expression", name)
		}
	}

	// Physical-arm gate (review M1): the default join at the named selective
	// thresholds must use the rank index, or the two engines are not running a
	// matched index arm. The gate records the plan and the engine's estimates
	// rather than proceeding silently to a matrix that would compare an index
	// path with a full scan.
	// The gate checks the thresholds the handoff names (20, 200 and 2000), not
	// the first three of whatever axis the run uses. The topology subset's third
	// threshold is 180000 (90%), where a sequential scan is the correct plan and
	// the index must not be required; gating on it wrongly refused to measure the
	// subset (run select-join-v3-topo-single-i1-20261005 recorded exactly that).
	for _, t := range gateThresholdsFor(thresholds) {
		args := selectbench.JoinArgs{Threshold: t, Refs: 4, ConsumerMax: len(jf.Facts), ValidateOnly: true, StatementTimeout: budget}
		specs := []selectbench.JoinSpec{{Statement: "j01_join_where", Axis: fmt.Sprintf("indexgate=%d", t), Args: args}}
		cells := selectbench.RunJoinSpecs(ctx, q, cat, jf, specs, "feasibility", "validation")
		c := cells[0]
		f.DefaultIndexPlans[fmt.Sprint(t)] = c.Plan
		switch {
		case c.Status != selectbench.StatusValid:
			f.DefaultIndexGate[fmt.Sprint(t)] = "error: " + string(c.Status) + ": " + c.Note
		case strings.Contains(c.Plan, rankIndex):
			f.DefaultIndexGate[fmt.Sprint(t)] = "used"
		default:
			f.DefaultIndexGate[fmt.Sprint(t)] = "NOT USED"
			f.Notes = append(f.Notes, fmt.Sprintf("physical-arm gate: j01 at threshold %d did not name %s", t, rankIndex))
		}
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

func runSelectJoinMode(dsn, engine, topology, scale, out, reportOut string, seed int64, trials int, environment, runID, phase, placement, indexArm string, statementTimeout int, positionControl, feasibilityOnly, topoSubset bool) error {
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
	if phase == "" {
		phase = "isolated"
	}
	if phase != "isolated" && phase != "blended" {
		return fmt.Errorf("phase must be isolated or blended, got %q", phase)
	}
	cat, err := selectbench.LoadJoin()
	if err != nil {
		return err
	}
	// The arm is chosen by the scale name: `cte`/`cte-small` measure the
	// synthetic microfixture, `cte-domain`/`cte-domain-small` measure the Study 01
	// person/donation fixture bound through the charity views.
	arm := selectJoinArm(seed, scale, engine, placement, indexArm)
	thresholds := selectbench.JoinThresholds(arm.factCount)
	if topoSubset {
		thresholds = selectbench.JoinTopologyThresholds(arm.factCount)
	}
	budget := time.Duration(statementTimeout) * time.Millisecond

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Minute)
	defer cancel()

	// One pool per endpoint: a cluster arm rotates its reader across every node
	// rather than silently using the seed endpoint alone.
	pools, endpoints, err := connectEndpoints(ctx, dsn, arm.schema, statementTimeout)
	if err != nil {
		return err
	}
	defer closePools(pools)

	q := &rrQuerier{pools: pools, endpoints: endpoints, ops: make([]int, len(pools))}
	explainOpts, explainEvidence := probeExplainOptions(ctx, pools[0], engine)
	q.explainOpts = explainOpts
	cpuBefore := readSelectCPUStat()

	result := selectbench.JoinResult{
		RunID: runID, Engine: engine, Topology: topology, Scale: scale,
		Arm: arm.dataset, FactCount: arm.factCount, Customers: arm.parents, Seed: seed,
		Thresholds: thresholds, Summary: map[string]int{},
		ExplainOptions: explainOpts,
	}
	result.Notes = append(result.Notes, explainEvidence)

	gateCells, gateControls, gateErr := runJoinGate(ctx, q, cat, arm)
	_ = gateCells
	if gateErr != nil {
		return writeJoinOutputs(result, gateErr, false, out, reportOut)
	}

	if err := arm.prepare(ctx, q); err != nil {
		return fmt.Errorf("load fixture: %w", err)
	}
	// Physical-arm evidence: the engine's own index definitions, and on
	// YugabyteDB the colocation/tablet properties of each domain table. These
	// make the arm checkable instead of asserted (review M1/M2).
	result.IndexDefs = queryIndexDefs(ctx, q, arm.schema)
	if engine == "yugabyte" && arm.dataset == "study01-person-donation-charity" {
		result.Placement = queryPlacements(ctx, q, arm.schema,
			[]string{"charity", "person", "donation", "donation_pkey", "donation_rank_idx", "donation_person_idx"})
	}
	result.Controls = append(result.Controls, gateControls...)

	// The limit-before-join control is proved on the large fixture (the edge
	// fixture is too small for a limit of 7 to truncate anything).
	bigControls := selectbench.RunJoinControls(ctx, q, cat, arm.oracle, []selectbench.JoinSpec{{
		Statement: "x02_join_limit_before", Axis: "control=x02_join_limit_before",
		Args: selectbench.JoinArgs{Threshold: int(float64(arm.factCount) * 10 / 100), Trials: 1, StatementTimeout: budget},
	}})
	for _, c := range bigControls {
		if !c.Detected {
			return writeJoinOutputs(result, fmt.Errorf("limit-before-join control did not fire on the large fixture: %s", c.Error), feasibilityOnly, out, reportOut)
		}
	}
	result.Controls = append(result.Controls, bigControls...)

	result.Feasibility = joinFeasibility(ctx, pools[0], q, cat, arm.oracle, thresholds, budget, arm.schema, arm.rankIndex)

	// The physical-arm gate is a precondition, not an observation: if the
	// default join does not use the rank index, the "matched index arm" claim is
	// false and the matrix would compare an index path with a full scan. Record
	// the plans and stop before measuring.
	if gateFailed(result.Feasibility) {
		result.Notes = append(result.Notes, "physical-arm gate failed for the default join at its most selective thresholds; matrix not measured (see feasibility.default_index_gate and default_index_plans)")
		return writeJoinOutputs(result, fmt.Errorf("physical-arm gate failed: default join does not use %s", arm.rankIndex), feasibilityOnly, out, reportOut)
	}

	base := selectbench.JoinArgs{Trials: trials, StatementTimeout: budget}
	// Blended phase: the fixed writer starts before the measured cells so no
	// cell runs concurrently with a writer that has not started, and the reader
	// overlap counters are scoped to exactly this window.
	result.Phase = phase
	var writer *maintenanceWriter
	if phase == "blended" {
		table, idCol := writerPlan(arm)
		writer = newMaintenanceWriter(pools, arm.schema, table, idCol,
			selectbench.MaintenanceKeys(arm.factCount, maintenanceKeyCount), selectbench.WriterTargetRate)
		q.writer = writer
		writer.start(ctx)
	}
	q.resetWriterOverlap()
	phaseStart := time.Now()
	if !feasibilityOnly {
		var specs []selectbench.JoinSpec
		if topoSubset {
			specs = selectbench.BuildJoinTopologySpecs(base, arm.factCount, trials)
		} else {
			specs = selectbench.BuildJoinSpecs(base, arm.factCount, trials)
		}
		plan := selectbench.PlanOrder(selectbench.OrderSeedFor(seed, runID), len(specs))
		first := specs[plan.Order[0]]
		result.Cells = selectbench.RunJoinSpecsInterleaved(ctx, q, cat, arm.oracle, specs, plan, runID, topology)
		result.Order = selectbench.RunOrder{Seed: plan.Seed, Rounds: trials, Sequence: plan.JoinSequence(specs)}
		if positionControl {
			pc := first
			pc.Axis += selectbench.PositionSuffix
			result.Cells = append(result.Cells, selectbench.RunJoinSpecs(ctx, q, cat, arm.oracle, []selectbench.JoinSpec{pc}, runID, topology)...)
			result.Order.Sequence = append(result.Order.Sequence, pc.Axis)
		}
	}
	phaseWall := time.Since(phaseStart).Seconds()
	samples, overlap := q.writerOverlap()
	result.WriterSamples, result.WriterOverlapSamples = samples, overlap
	if writer != nil {
		writer.stopAndWait()
		st := writer.stats(phase, phaseWall)
		k, mm, detail := writer.ledger(ctx)
		st.LedgerKeys, st.LedgerMismatch, st.LedgerDetail = k, mm, detail
		result.Writer = st
	} else {
		result.Writer = &selectbench.WriterStats{Phase: phase, TargetRate: selectbench.WriterTargetRate}
	}

	version := ""
	_ = pools[0].QueryRow(ctx, "SELECT version()").Scan(&version)
	result.Endpoints = endpoints
	result.EndpointOps = q.ops
	result.EndpointConns = q.conns()
	result.Provenance = joinProvenance(version, environment, scale, arm.dataset, arm.factCount, arm.parents, thresholds, engine, topology, seed, trials, statementTimeout, cpuBefore, arm.fingerprint, endpoints)
	result.Provenance.Settings["phase"] = phase
	result.Provenance.Settings["writer_target_rate_per_sec"] = fmt.Sprint(selectbench.WriterTargetRate)
	result.Provenance.Settings["writer_keys"] = fmt.Sprint(maintenanceKeyCount)
	placementLabel := arm.placement
	if placementLabel == "" {
		placementLabel = "engine default"
	}
	result.Provenance.Settings["placement"] = placementLabel
	if arm.placementNote != "" {
		result.Notes = append(result.Notes, arm.placementNote)
	}
	// The physical index arm and the node/client budgets are arm factors: two
	// runs that differ in any of them are not the same experiment, so they enter
	// the compatibility class (review M3).
	result.Provenance.Settings["index_arm"] = arm.indexArm
	if arm.indexArmNote != "" {
		result.Notes = append(result.Notes, arm.indexArmNote)
	}
	result.Provenance.Settings["node_budget"] = envOr("ADS_NODE_BUDGET", "unknown")
	result.Provenance.Settings["client_budget"] = envOr("ADS_CLIENT_BUDGET", "unknown")
	result.Notes = append(result.Notes, joinGaps(engine, topology, phase, len(endpoints), result.Provenance.Settings["node_budget"], result.Provenance.Settings["client_budget"], arm.indexArm)...)
	return writeJoinOutputs(result, nil, feasibilityOnly, out, reportOut)
}

// maintenanceKeyCount is the size of the fixed writer key set.
const maintenanceKeyCount = 1000

func joinProvenance(version, environment, scale, dataset string, factCount, parents int, thresholds []int, engine, topology string, seed int64, trials, statementTimeout int, cpuBefore, fingerprint string, endpoints []string) selectbench.Provenance {
	rpc := "DIST/DEBUG storage counters parsed per cell (Storage Rows Scanned, Storage Index Rows Scanned); RPC-level counters not aggregated"
	if engine == "postgres" {
		rpc = "N/A on PostgreSQL; shared buffers and rows removed parsed per cell where printed"
	}
	mode := "single endpoint"
	if len(endpoints) > 1 {
		mode = "round-robin across " + fmt.Sprint(len(endpoints)) + " endpoints"
	}
	settings := map[string]string{
		"scale": scale, "dataset": dataset, "topology": topology, "fact_count": fmt.Sprint(factCount),
		"customers": fmt.Sprint(parents), "thresholds": fmt.Sprint(thresholds),
		"selectivity_mapping":  "rank <= threshold; rank is a deterministic shuffle of 1..fact_count",
		"index_set":            "child rank index and parent-key index, primary keys; no per-spelling index",
		"statistics":           "ANALYZE after deterministic load",
		"prepared_plan_policy": "pgx simple protocol; no prepared statements",
		"client_workers":       "1 sequential reader; 0 writers", "connection_pool": "2 per endpoint = reader/setup + reserve",
		"endpoints": strings.Join(endpoints, ","), "endpoint_mode": mode,
		"trials_per_cell": fmt.Sprint(trials), "statement_budget_ms": fmt.Sprint(statementTimeout),
		"whole_run_count": "1; no replicated winner in this file", "distributed_rpc": rpc,
	}
	return selectbench.Provenance{
		RepoCommit: envOr("ADS_REPO_COMMIT", "unknown"),
		// The client and the database are different images. Keep them
		// separate: labelling the client image `EngineImage` beside a database
		// digest makes the provenance self-contradictory (review R4).
		ClientImage:         envOr("ADS_BENCH_IMAGE", "unknown"),
		ClientImageDigest:   os.Getenv("ADS_BENCH_IMAGE_ID"),
		ImageDigest:         os.Getenv("ADS_BENCH_IMAGE_ID"),
		EngineImage:         envOr("ADS_ENGINE_IMAGE", "unknown"),
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

// joinGaps derives a run's recorded limits from its own settings. Writing a
// fixed list produced limits that contradicted the run: every blended report
// still said "no concurrent writers", every equal-total or cluster run said
// "no equal-total-budget comparison", and a captured YugabyteDB DIST plan still
// said the RPC counters were unparsed (review M4).
func joinGaps(engine, topology, phase string, endpointCount int, nodeBudget, clientBudget, indexArm string) []string {
	gaps := []string{
		fmt.Sprintf("one host: this run used node budget %s and client budget %s, so budget comparisons are separate classes here, not independent hardware", nodeBudget, clientBudget),
	}
	if phase == "blended" {
		gaps = append(gaps, fmt.Sprintf("one client process with one fixed maintenance writer at %d updates/s; readers are sequential; the reported sample latency includes retry and backoff time", selectbench.WriterTargetRate))
	} else {
		gaps = append(gaps, "isolated phase: no concurrent writers in this run; the blended read/write phase is a separate labelled run")
	}
	if topology != "single" {
		gaps = append(gaps, "cluster nodes are containers on the same host, so these are not network or independent-hardware measurements; container colocation is not data-colocation evidence")
	}
	if endpointCount > 1 {
		gaps = append(gaps, fmt.Sprintf("readers round-robin across %d endpoints; the writer's endpoint policy is recorded separately", endpointCount))
	}
	if indexArm == "" {
		indexArm = "engine default"
	}
	gaps = append(gaps, "physical rank-index arm "+indexArm+"; an engine whose planner ignores the index is a different physical design, not a faster one")
	if engine == "yugabyte" {
		gaps = append(gaps, "YugabyteDB DIST/DEBUG storage counters (rows and index rows scanned) are parsed per cell; RPC-level counters are not aggregated")
	} else {
		gaps = append(gaps, "distributed RPC counters are N/A on single-node PostgreSQL; shared buffers and rows removed are parsed per cell where printed")
	}
	gaps = append(gaps, "a flat cell reports a client or engine ceiling; it is not by itself a constant query cost")
	return gaps
}

// joinInputsDigest is the single definition of a run's inputs digest: the
// SHA-256 of the exact result-file bytes (which end in one newline). The writer
// and the re-renderer must both call it, or a report re-rendered from committed
// bytes advertises a digest the run never had.
func joinInputsDigest(data []byte) string {
	return fmt.Sprintf("%x", sha256.Sum256(data))
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
	inputsDigest := joinInputsDigest(data)
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

func armFilterLabel(armFilter string) string {
	if armFilter == "" {
		return ""
	}
	return ", arm=" + armFilter
}

// runSelectJoinReport re-renders one run from a directory holding join.json, or
// writes the cross-whole-run report for a directory of run subdirectories.
func runSelectJoinReport(resultsDir, reportOut, armFilter string) error {
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
		// Hash the exact file bytes. The writer appended the single trailing
		// newline before hashing, so adding another here (the original defect)
		// advertised a digest the run never had.
		digest := joinInputsDigest(data)
		r.NodeThrottle = loadNodeThrottle(resultsDir)
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
		// cells to replicate and must not appear as a whole run. When an arm
		// filter is given, runs from the other arm are excluded rather than
		// pooled; the report's compatibility classes then only ever see one arm.
		if len(r.Cells) == 0 {
			continue
		}
		if armFilter != "" && r.Arm != armFilter {
			continue
		}
		r.NodeThrottle = loadNodeThrottle(filepath.Dir(p))
		runs = append(runs, r)
	}
	if len(runs) == 0 {
		return fmt.Errorf("no whole join runs (with cells%s) at or under %s", armFilterLabel(armFilter), resultsDir)
	}
	report, err := selectbench.RenderJoinReplicated(runs)
	if err != nil {
		return err
	}
	return os.WriteFile(reportOut, []byte(report), 0o644)
}

// loadNodeThrottle reads the per-node throttling deltas the runner computed from
// the before/after cpu.stat files. An old run without the file renders "not
// recorded" rather than a false zero.
func loadNodeThrottle(runDir string) []selectbench.NodeThrottleStat {
	data, err := os.ReadFile(filepath.Join(runDir, "resources", "node-throttle.json"))
	if err != nil {
		return nil
	}
	var out []selectbench.NodeThrottleStat
	if err := json.Unmarshal(data, &out); err != nil {
		return nil
	}
	return out
}

// gateFailed reports whether the physical-arm gate found the default join not
// using the rank index at any gated threshold.
func gateFailed(f *selectbench.JoinFeasibility) bool {
	if f == nil {
		return false
	}
	for _, v := range f.DefaultIndexGate {
		if v == "NOT USED" {
			return true
		}
	}
	return false
}

// queryIndexDefs reads the engine's own definition of every index in a schema.
func queryIndexDefs(ctx context.Context, q selectbench.Querier, schema string) []selectbench.IndexDef {
	rows, err := q.Query(ctx, "SELECT indexname, indexdef FROM pg_indexes WHERE schemaname = '"+schema+"' ORDER BY indexname")
	if err != nil {
		return []selectbench.IndexDef{{Name: "pg_indexes", Def: "unavailable: " + err.Error()}}
	}
	var out []selectbench.IndexDef
	for _, r := range rows {
		if len(r) >= 2 {
			out = append(out, selectbench.IndexDef{Name: fmt.Sprint(r[0]), Def: fmt.Sprint(r[1])})
		}
	}
	return out
}

// queryPlacements reads YugabyteDB's colocation and tablet properties for the
// named relations, plus the local tablet count. On an engine without these
// functions the query fails and the list stays empty, which the report renders
// as not recorded rather than as a success.
func queryPlacements(ctx context.Context, q selectbench.Querier, schema string, tables []string) []selectbench.TablePlacement {
	var out []selectbench.TablePlacement
	for _, t := range tables {
		p := selectbench.TablePlacement{Table: schema + "." + t}
		rows, err := q.Query(ctx, fmt.Sprintf(
			"SELECT num_tablets, num_hash_key_columns, is_colocated, COALESCE(tablegroup_oid::text,''), COALESCE(colocation_id::text,'') FROM yb_table_properties('%s.%s'::regclass)",
			schema, t))
		if err == nil && len(rows) > 0 && len(rows[0]) >= 5 {
			p.NumTablets = anyInt(rows[0][0])
			p.NumHashKeyCols = anyInt(rows[0][1])
			p.IsColocated = anyBool(rows[0][2])
			p.TablegroupOID = fmt.Sprint(rows[0][3])
			p.ColocationID = fmt.Sprint(rows[0][4])
		}
		if crows, cerr := q.Query(ctx, "SELECT count(*) FROM yb_local_tablets WHERE table_name = '"+t+"'"); cerr == nil && len(crows) > 0 && len(crows[0]) > 0 {
			p.TabletCountLocal = anyInt(crows[0][0])
		}
		out = append(out, p)
	}
	return out
}

func anyInt(v any) int {
	switch n := v.(type) {
	case int:
		return n
	case int32:
		return int(n)
	case int64:
		return int(n)
	case float64:
		return int(n)
	default:
		i, _ := strconv.Atoi(fmt.Sprint(v))
		return i
	}
}

// containsInt reports whether xs contains v.
func containsInt(xs []int, v int) bool {
	for _, x := range xs {
		if x == v {
			return true
		}
	}
	return false
}

// gateThresholdsFor returns the thresholds the physical-arm gate checks: the
// handoff's 20/200/2000 when the run includes them, otherwise the two most
// selective thresholds the run has. The widest thresholds are excluded because
// a sequential scan there is the correct plan, not a missing index.
func gateThresholdsFor(thresholds []int) []int {
	var out []int
	for _, t := range []int{20, 200, 2000} {
		if containsInt(thresholds, t) {
			out = append(out, t)
		}
	}
	if len(out) > 0 {
		return out
	}
	if len(thresholds) > 2 {
		return thresholds[:2]
	}
	return thresholds
}

func anyBool(v any) bool {
	switch b := v.(type) {
	case bool:
		return b
	case string:
		return b == "t" || b == "true"
	default:
		s := fmt.Sprint(v)
		return s == "true" || s == "t"
	}
}
