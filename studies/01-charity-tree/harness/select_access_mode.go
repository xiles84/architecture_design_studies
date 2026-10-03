package main

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"time"

	"charitytree/selectbench"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// The access-path measurement mode (`-cmd select-access`).
//
// It reuses the additive SELECT schema, deterministic fixture, independent
// oracles and correctness gate from select_mode.go, and runs the bounded
// access-path matrix the protocol fixes for task select-access-paths:
//
//	P10  OR vs overlap-safe UNION ALL at three predicate shapes, plus the
//	     deliberately-wrong overlapping control per shape;
//	P11  OFFSET vs keyset at five page depths, each cell verified against the
//	     contract oracle;
//	     a fixed-snapshot page-equality check inside one transaction;
//	     a labelled expression-index arm with its storage and insert cost.
//
// Only the harness touches the database; the matrix, oracles and reports live
// in the selectbench package and are unit-tested without an engine.

func runSelectAccessMode(dsn, engine, topology, scale, out, reportOut string, seed int64, trials, pageLimit int, environment, runID string, statementTimeout int) error {
	if dsn == "" {
		return fmt.Errorf("-cmd select-access needs -dsn")
	}
	for _, key := range []string{"ADS_RUN_LOCK_HELD", "ADS_REPO_COMMIT", "ADS_RUN_TAG", "ADS_BENCH_IMAGE_ID", "ADS_ENGINE_IMAGE_ID"} {
		if strings.TrimSpace(os.Getenv(key)) == "" {
			return fmt.Errorf("SELECT access requires %s from the locked, tagged runner", key)
		}
	}
	if environment == "" || environment == "unknown" || runID == "" || trials < 1 || pageLimit < 1 || statementTimeout < 1 {
		return fmt.Errorf("SELECT access requires a named environment/run and positive trials/limit/statement budget")
	}
	cat, err := selectbench.Load()
	if err != nil {
		return err
	}
	spec := selectSpec(scale)
	fx := selectbench.Generate(seed, spec)

	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Minute)
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

	result := selectbench.AccessResult{
		RunID: runID, Engine: engine, Topology: topology, Scale: scale, Seed: seed,
		Limit: pageLimit, Shapes: shapeNames(selectbench.DefaultPredicateShapes()),
		Depths: selectbench.DefaultPageDepths(), Summary: map[string]int{},
	}

	// Gate 1: the reviewed live fixture/control gate over the whole catalogue on
	// the teaching fixture (includes the overlap control's unit-level shape).
	validation, teachingControls, gateErr := validateSelectFixtures(ctx, q, cat)
	if gateErr != nil {
		return writeAccessOutputs(result, validation, gateErr, out, reportOut)
	}
	_ = teachingControls

	if err := selectbench.Prepare(ctx, q, fx); err != nil {
		return fmt.Errorf("load fixture: %w", err)
	}

	base := selectbench.Args{
		Pred: selectbench.AmountGreaterThan(100), PredSQL: "amount > 100",
		Limit: pageLimit, Trials: trials, StatementTimeout: time.Duration(statementTimeout) * time.Millisecond,
	}
	shapes := selectbench.DefaultPredicateShapes()
	depths := selectbench.DefaultPageDepths()

	// The default arm must contain no added index: drop the expression index if
	// a previous run left one behind.
	if err := q.Exec(ctx, selectbench.DropExpressionIndexSQL()); err != nil {
		result.Notes = append(result.Notes, "could not drop a pre-existing expression index: "+err.Error())
	}

	// Overlap proof per shape (disjointness is what makes UNION ALL comparable).
	for _, shape := range shapes {
		rec := selectbench.OverlapRecord{
			Shape: shape.Name, ASQL: shape.ASQL, BSQL: shape.BSQL,
			ArmAMatches: countPred(fx, shape.A), ArmBMatches: countPred(fx, shape.B),
			OverlapCount: selectbench.OverlapCount(fx, shape.A, shape.B),
		}
		rec.Disjoint = rec.OverlapCount == 0
		result.Overlap = append(result.Overlap, rec)
		if !rec.Disjoint {
			return fmt.Errorf("shape %s: arms overlap on %d rows; UNION ALL is not a comparable spelling", shape.Name, rec.OverlapCount)
		}
		if rec.ArmAMatches == 0 {
			return fmt.Errorf("shape %s: arm A matches no row; the overlap control cannot fire", shape.Name)
		}
	}

	// Default arm cells.
	specs := selectbench.BuildAccessSpecs(fx, base, shapes, depths, pageLimit, trials)
	result.Cells = append(result.Cells, selectbench.RunSpecs(ctx, q, cat, fx, specs, runID, topology)...)

	// The overlap control must fire in every predicate shape on the measured
	// data: that is this family's required negative control.
	for _, shape := range shapes {
		ctrlArgs := base
		ctrlArgs.PredA, ctrlArgs.PredB = shape.A, shape.B
		ctrlArgs.PredASQL, ctrlArgs.PredBSQL = shape.ASQL, shape.BSQL
		controls := selectbench.RunControls(ctx, q, cat, fx, ctrlArgs)
		for _, c := range controls {
			result.Controls = append(result.Controls, selectbench.AccessControl{Shape: shape.Name, ControlResult: c})
		}
		var overlap *selectbench.ControlResult
		for i := range controls {
			if controls[i].Statement == "x05_union_all_overlap" {
				overlap = &controls[i]
			}
		}
		if overlap == nil || !overlap.Detected {
			note := "not run"
			if overlap != nil {
				note = overlap.Error
			}
			return fmt.Errorf("shape %s: the overlapping-UNION-ALL control did not fire: %s", shape.Name, note)
		}
	}

	// Fixed-snapshot page equality.
	result.Snapshot = checkSnapshot(ctx, pool, cat, fx, depths, pageLimit)

	// Labelled expression-index arm: the query text is unchanged, so any gain
	// belongs to the index, not to the spelling.
	if err := q.Exec(ctx, selectbench.ExpressionIndexSQL()); err != nil {
		result.Index = selectbench.IndexCost{
			IndexName: selectbench.ExpressionIndexName, Expression: "(amount + 0)",
			Supported: false, Note: "expression index not created: " + err.Error(),
		}
		result.Notes = append(result.Notes, "expression-index arm skipped: "+err.Error()+" (explicit coverage gap)")
	} else {
		if err := q.Exec(ctx, "ANALYZE "+selectbench.TeachingSchema+".orders"); err != nil {
			result.Notes = append(result.Notes, "ANALYZE after creating the expression index failed: "+err.Error())
		}
		for _, exprShape := range shapes {
			if !strings.HasSuffix(exprShape.Name, "/expression") {
				continue
			}
			indexSpecs := shapeSpecsForIndex(base, exprShape, trials, "expression_index")
			result.Cells = append(result.Cells, selectbench.RunSpecs(ctx, q, cat, fx, indexSpecs, runID, topology)...)
		}
		result.Index = measureIndexCost(ctx, pool, fx, 500)
	}

	version := ""
	_ = pool.QueryRow(ctx, "SELECT version()").Scan(&version)

	result.Provenance = selectbench.Provenance{
		RepoCommit:          envOr("ADS_REPO_COMMIT", "unknown"),
		EngineImage:         envOr("ADS_BENCH_IMAGE", "unknown"),
		ImageDigest:         os.Getenv("ADS_BENCH_IMAGE_ID"),
		EngineImageDigest:   os.Getenv("ADS_ENGINE_IMAGE_ID"),
		EngineVersion:       version,
		Environment:         environment,
		RunTag:              os.Getenv("ADS_RUN_TAG"),
		Seed:                seed,
		DatasetFingerprint:  selectbench.Fingerprint(fx),
		Settings:            accessSettings(scale, engine, topology, pageLimit, trials, statementTimeout, result.Snapshot),
		ClientCPUStatBefore: cpuBefore,
		ClientCPUStatAfter:  readSelectCPUStat(),
		ClientCPUQuota:      readSelectResource("/sys/fs/cgroup/cpu.max"),
		ClientMemoryLimit:   readSelectResource("/sys/fs/cgroup/memory.max"),
	}
	for _, c := range result.Cells {
		result.Summary[string(c.Status)]++
	}
	result.Notes = append(result.Notes, accessGaps(engine, topology, result)...)
	return writeAccessOutputs(result, validation, nil, out, reportOut)
}

func shapeNames(shapes []selectbench.PredShape) []string {
	out := make([]string, 0, len(shapes))
	for _, s := range shapes {
		out = append(out, s.Name)
	}
	return out
}

func countPred(fx selectbench.Fixture, pred selectbench.Pred) int {
	if pred == nil {
		return 0
	}
	n := 0
	for _, o := range fx.Orders {
		if pred(o) {
			n++
		}
	}
	return n
}

// shapeSpecsForIndex exposes the index-arm spec builder via BuildAccessSpecs on
// a single selected shape.
func shapeSpecsForIndex(base selectbench.Args, shape selectbench.PredShape, trials int, arm string) []selectbench.CellSpec {
	return selectbench.IndexArmSpecs(base, shape, trials, arm)
}

func accessSettings(scale, engine, topology string, pageLimit, trials, timeout int, snap selectbench.SnapshotCheck) map[string]string {
	s := map[string]string{
		"scale": scale, "topology": topology, "page_limit": fmt.Sprint(pageLimit),
		"trials_per_cell": fmt.Sprint(trials), "statement_budget_ms": fmt.Sprint(timeout),
		"index_set": "orders_amount_idx(amount),orders_customer_idx(customer_id),primary keys;" +
			" expression arm adds " + selectbench.ExpressionIndexName + "((amount + 0))",
		"statistics": "ANALYZE after deterministic load", "prepared_plan_policy": "pgx simple protocol; no prepared statements",
		"client_workers": "1 sequential reader; 0 writers", "connection_pool": "2 = 1 reader/setup + 1 reserve",
		"snapshot_isolation": snap.Isolation,
	}
	if engine != "postgres" || topology != "single" {
		s["distributed_rpc"] = "not parsed; raw EXPLAIN retained (coverage gap)"
	} else {
		s["distributed_rpc"] = "N/A on PostgreSQL single-node"
	}
	return s
}

func accessGaps(engine, topology string, result selectbench.AccessResult) []string {
	var notes []string
	notes = append(notes, "no concurrent writers: client-decode JSON bytes are not wire bytes; tails are under-sampled at the protocol trial count")
	notes = append(notes, "one host: no independent hardware, real-network, 3-node placement or equal-total-budget comparison is established")
	if topology == "single" {
		notes = append(notes, "single-node only; the YB 3-node RF=3 per-node/equal-total arms are a coverage gap on this host")
	}
	return notes
}

// --- index cost ------------------------------------------------------------

func measureIndexCost(ctx context.Context, pool *pgxpool.Pool, fx selectbench.Fixture, rows int) selectbench.IndexCost {
	q := pgxQuerier{pool: pool}
	ic := selectbench.IndexCost{IndexName: selectbench.ExpressionIndexName, Expression: "(amount + 0)", Supported: true, Rows: rows}
	startID := 2_000_000_000
	batch := indexCostBatch(fx, startID, rows)

	// Without the expression index.
	if err := q.Exec(ctx, selectbench.DropExpressionIndexSQL()); err != nil {
		ic.Note = "drop index: " + err.Error()
	}
	without, err := timedExec(ctx, q, batch)
	if err != nil {
		ic.Note = "insert without index: " + err.Error()
		return ic
	}
	ic.WithoutIndexMS = &without
	_ = q.Exec(ctx, deleteBatchSQL(startID, rows))

	// With the expression index.
	if err := q.Exec(ctx, selectbench.ExpressionIndexSQL()); err != nil {
		ic.Supported = false
		ic.Note = "create expression index: " + err.Error()
		return ic
	}
	if err := q.Exec(ctx, "ANALYZE "+selectbench.TeachingSchema+".orders"); err != nil {
		ic.Note = "analyze: " + err.Error()
	}
	with, err := timedExec(ctx, q, batch)
	if err != nil {
		ic.Note = "insert with index: " + err.Error()
		return ic
	}
	ic.WithIndexMS = &with
	if size, err := scalarInt64(ctx, pool, "SELECT pg_relation_size('"+selectbench.TeachingSchema+"."+selectbench.ExpressionIndexName+"')"); err == nil {
		ic.SizeBytes = &size
	} else if ic.Note == "" {
		ic.Note = "pg_relation_size unavailable: " + err.Error()
	}
	_ = q.Exec(ctx, deleteBatchSQL(startID, rows))
	return ic
}

func indexCostBatch(fx selectbench.Fixture, startID, rows int) string {
	var b strings.Builder
	b.WriteString("INSERT INTO " + selectbench.TeachingSchema + ".orders (id, customer_id, amount, payload) VALUES ")
	n := len(fx.Customers)
	if n == 0 {
		n = 1
	}
	for k := 0; k < rows; k++ {
		if k > 0 {
			b.WriteString(",")
		}
		cust := (k % n) + 1
		fmt.Fprintf(&b, "(%d,%d,%d,NULL)", startID+k, cust, (k%500)+1)
	}
	return b.String()
}

func deleteBatchSQL(startID, rows int) string {
	return fmt.Sprintf("DELETE FROM %s.orders WHERE id >= %d AND id < %d", selectbench.TeachingSchema, startID, startID+rows)
}

func timedExec(ctx context.Context, q selectbench.Querier, sql string) (float64, error) {
	start := time.Now()
	if err := q.Exec(ctx, sql); err != nil {
		return 0, err
	}
	return float64(time.Since(start).Nanoseconds()) / 1e6, nil
}

func scalarInt64(ctx context.Context, pool *pgxpool.Pool, sql string) (int64, error) {
	var v int64
	if err := pool.QueryRow(ctx, sql).Scan(&v); err != nil {
		return 0, err
	}
	return v, nil
}

// --- fixed snapshot --------------------------------------------------------

// checkSnapshot runs the OFFSET and keyset pages for every depth inside ONE
// transaction and requires the two pages to contain the same rows. With no
// concurrent writers the fixture is fixed by construction; the transaction makes
// the "same snapshot" part explicit rather than assumed.
func checkSnapshot(ctx context.Context, pool *pgxpool.Pool, cat *selectbench.Catalog, fx selectbench.Fixture, depths []int, limit int) selectbench.SnapshotCheck {
	sc := selectbench.SnapshotCheck{AllEqual: true}
	tx, err := pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.RepeatableRead})
	if err != nil {
		tx, err = pool.Begin(ctx)
		if err != nil {
			sc.Error = "begin: " + err.Error()
			sc.AllEqual = false
			return sc
		}
		sc.Isolation = "read committed (REPEATABLE READ rejected by the engine; fixed only by the absence of writers)"
	} else {
		sc.Isolation = "repeatable read"
	}
	defer tx.Rollback(ctx)

	offStmt := cat.Statements["s16_page_offset"]
	keyStmt := cat.Statements["s17_page_keyset"]
	for _, d := range depths {
		offArgs := selectbench.Args{Limit: limit, Offset: d}
		keyArgs := selectbench.Args{Limit: limit, Last: snapshotCursor(fx, d-1)}
		offRows, err := queryTx(ctx, tx, offStmt, offArgs)
		if err != nil {
			sc.Error = fmt.Sprintf("depth %d offset: %v", d, err)
			sc.AllEqual = false
			break
		}
		keyRows, err := queryTx(ctx, tx, keyStmt, keyArgs)
		if err != nil {
			sc.Error = fmt.Sprintf("depth %d keyset: %v", d, err)
			sc.AllEqual = false
			break
		}
		offCanon, err1 := selectbench.Canonicalize("sel_pagination", offRows)
		keyCanon, err2 := selectbench.Canonicalize("sel_pagination", keyRows)
		// Canonical values hold *int amounts, so compare by value, not by the
		// pointers' printed addresses.
		equal := err1 == nil && err2 == nil && reflect.DeepEqual(offCanon, keyCanon)
		if !equal {
			sc.AllEqual = false
		}
		sc.Depths = append(sc.Depths, selectbench.SnapshotAt{Depth: d, Equal: equal, OffsetN: len(offRows), KeysetN: len(keyRows)})
	}
	return sc
}

func snapshotCursor(fx selectbench.Fixture, index int) *selectbench.PageRow {
	return selectbench.CursorAt(fx, index)
}

func queryTx(ctx context.Context, tx pgx.Tx, st selectbench.Statement, a selectbench.Args) ([][]any, error) {
	bound, _, err := selectbench.Bind(st.SQL, a.ParamMap(st))
	if err != nil {
		return nil, err
	}
	rows, err := tx.Query(ctx, bound)
	if err != nil {
		return nil, classifyEngineError(err)
	}
	defer rows.Close()
	var out [][]any
	for rows.Next() {
		vals, err := rows.Values()
		if err != nil {
			return nil, err
		}
		out = append(out, vals)
	}
	return out, rows.Err()
}

// --- outputs ---------------------------------------------------------------

func writeAccessOutputs(result selectbench.AccessResult, validation []selectValidation, gateErr error, out, reportOut string) error {
	if gateErr != nil {
		result.Error = gateErr.Error()
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
		report, err := selectbench.RenderAccessRun(result)
		if err != nil {
			return err
		}
		if err := os.WriteFile(reportOut, []byte(report), 0o644); err != nil {
			return err
		}
		fmt.Println("wrote", reportOut)
	}
	_ = validation
	fmt.Printf("access: %d cells — valid=%d invalid=%d unsupported=%d timeout=%d error=%d; inputs digest %s\n",
		len(result.Cells), result.Summary["valid"], result.Summary["invalid"],
		result.Summary["unsupported"], result.Summary["timeout"], result.Summary["error"], inputsDigest)
	if result.Error != "" {
		return fmt.Errorf("%s", result.Error)
	}
	if result.Summary["invalid"]+result.Summary["error"]+result.Summary["timeout"]+result.Summary["unsupported"] > 0 {
		return fmt.Errorf("SELECT access has non-valid cells; inspect their recorded status")
	}
	return nil
}

// runSelectAccessReport writes an access-path measurements report. Given a
// directory that itself holds `access.json` it re-renders that single run
// (useful for regenerating a report from a committed result); given a directory
// of run subdirectories it writes the cross-whole-run report.
func runSelectAccessReport(resultsDir, reportOut string) error {
	if resultsDir == "" || reportOut == "" {
		return fmt.Errorf("-cmd access-report needs -results <dir> and -report-out <file>")
	}
	if fi, err := os.Stat(filepath.Join(resultsDir, "access.json")); err == nil && !fi.IsDir() {
		data, err := os.ReadFile(filepath.Join(resultsDir, "access.json"))
		if err != nil {
			return err
		}
		var r selectbench.AccessResult
		if err := json.Unmarshal(data, &r); err != nil {
			return err
		}
		report, err := selectbench.RenderAccessRun(r)
		if err != nil {
			return err
		}
		return os.WriteFile(reportOut, []byte(report), 0o644)
	}
	paths, err := filepath.Glob(filepath.Join(resultsDir, "*", "access.json"))
	if err != nil {
		return err
	}
	if len(paths) == 0 {
		return fmt.Errorf("no access.json at or under %s", resultsDir)
	}
	sort.Strings(paths)
	var runs []selectbench.AccessResult
	for _, p := range paths {
		data, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		var r selectbench.AccessResult
		if err := json.Unmarshal(data, &r); err != nil {
			return fmt.Errorf("%s: %w", p, err)
		}
		runs = append(runs, r)
	}
	report, err := selectbench.RenderAccessReplicated(runs)
	if err != nil {
		return err
	}
	return os.WriteFile(reportOut, []byte(report), 0o644)
}
