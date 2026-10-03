package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"charitytree/selectbench"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// SELECT experiment mode (`-cmd select`).
//
// This is the only file that touches a database for the SELECT families. The
// cell loop, oracles, controls, plan signatures, report and verification all
// live in the selectbench package and are unit-tested without an engine; here we
// only supply a pgx-backed selectbench.Querier, load the deterministic fixture,
// run the cells and write the measurements. The runner enforces the correctness
// gate: a cell whose answer differs from the oracle is marked invalid and is
// never timed.
//
// The existing Study 01 modes are untouched: this mode writes to its own schema
// (selectbench) and its own output file, so historical results stay reproducible.

// pgxQuerier adapts a pgx pool to selectbench.Querier.
type pgxQuerier struct{ pool *pgxpool.Pool }

func classifyEngineError(err error) error {
	if err == nil {
		return nil
	}
	msg := strings.ToLower(err.Error())
	for _, marker := range []string{"does not exist", "syntax error", "not supported", "unrecognized", "unsupported"} {
		if strings.Contains(msg, marker) {
			return fmt.Errorf("%w: %v", selectbench.ErrUnsupported, err)
		}
	}
	return err
}

func (q pgxQuerier) Exec(ctx context.Context, sql string) error {
	_, err := q.pool.Exec(ctx, sql)
	return classifyEngineError(err)
}

func (q pgxQuerier) Explain(ctx context.Context, sql string) (string, error) {
	rows, err := q.pool.Query(ctx, "EXPLAIN (ANALYZE, BUFFERS, FORMAT TEXT) "+sql)
	if err != nil {
		return "", classifyEngineError(err)
	}
	defer rows.Close()
	var b strings.Builder
	for rows.Next() {
		var line string
		if err := rows.Scan(&line); err != nil {
			return "", err
		}
		b.WriteString(line)
		b.WriteString("\n")
	}
	return b.String(), rows.Err()
}

func (q pgxQuerier) Query(ctx context.Context, sql string) ([][]any, error) {
	rows, err := q.pool.Query(ctx, sql)
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

// selectSpec maps the harness scale onto a deterministic dataset shape. The
// fixture is generated in Go and loaded into the engine, so the oracle and the
// query always see the same rows.
func selectSpec(scale string) selectbench.GenSpec {
	switch scale {
	case "tiny":
		return selectbench.GenSpec{Customers: 3, OrdersPerParent: 2, NullAmountEvery: 4, NullFKEvery: 7}
	case "medium":
		return selectbench.GenSpec{Customers: 200, OrdersPerParent: 10, NullAmountEvery: 11, NullFKEvery: 97, Wide: true}
	case "large":
		return selectbench.GenSpec{Customers: 2000, OrdersPerParent: 20, NullAmountEvery: 13, NullFKEvery: 199,
			HotCustomerShare: 0.5, Wide: true}
	default: // small
		return selectbench.GenSpec{Customers: 50, OrdersPerParent: 5, NullAmountEvery: 7, NullFKEvery: 23}
	}
}

type selectRun struct {
	RunID      string                 `json:"run_id"`
	Engine     string                 `json:"engine"`
	Topology   string                 `json:"topology"`
	Scale      string                 `json:"scale"`
	Seed       int64                  `json:"seed"`
	Fixture    string                 `json:"fixture_fingerprint"`
	Provenance selectbench.Provenance `json:"provenance"`
	Cells      []selectbench.Cell     `json:"cells"`
	Summary    map[string]int         `json:"summary"`
	Error      string                 `json:"error,omitempty"`
}

func runSelectMode(dsn, engine, topology, scale, queries, out, reportOut string, seed int64) error {
	if dsn == "" {
		return fmt.Errorf("-cmd select needs -dsn")
	}
	cat, err := selectbench.Load()
	if err != nil {
		return err
	}

	spec := selectSpec(scale)
	fx := selectbench.Generate(seed, spec)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return fmt.Errorf("dsn: %w", err)
	}
	// The catalogue uses bare table names; resolve them in the SELECT schema.
	cfg.ConnConfig.RuntimeParams["search_path"] = selectbench.TeachingSchema + ",public"
	cfg.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return fmt.Errorf("connect: %w", err)
	}
	defer pool.Close()

	q := pgxQuerier{pool: pool}
	if err := selectbench.Prepare(ctx, q, fx); err != nil {
		return fmt.Errorf("load fixture: %w", err)
	}

	var names []string
	if strings.TrimSpace(queries) != "" {
		for _, n := range strings.Split(queries, ",") {
			if n = strings.TrimSpace(n); n != "" {
				names = append(names, n)
			}
		}
	}

	a := selectbench.Args{
		Pred:    selectbench.AmountGreaterThan(100),
		PredSQL: "amount > 100",
		N:       3,
		Limit:   2,
		Offset:  0,
	}
	runID := time.Now().UTC().Format("20060102T150405Z")
	cells := selectbench.RunCells(ctx, q, cat, fx, a, names, runID, topology)

	prov := selectbench.Provenance{
		RepoCommit:         envOr("ADS_REPO_COMMIT", "unknown"),
		EngineImage:        engine,
		Seed:               seed,
		DatasetFingerprint: selectbench.Fingerprint(fx),
		Settings:           map[string]string{"scale": scale, "topology": topology},
	}

	result := selectRun{
		RunID: runID, Engine: engine, Topology: topology, Scale: scale, Seed: seed,
		Fixture: selectbench.Fingerprint(fx), Provenance: prov, Cells: cells,
		Summary: map[string]int{},
	}
	for _, c := range cells {
		result.Summary[string(c.Status)]++
	}

	if out != "" {
		data, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return err
		}
		if err := os.WriteFile(out, append(data, '\n'), 0o644); err != nil {
			return err
		}
		fmt.Println("wrote", out)
	}

	if reportOut != "" {
		report, err := selectbench.RenderComparison(selectbench.ReportInput{
			Provenance: prov,
			Cells:      cells,
			Predictions: []string{"P1-P5 join/CTE spellings", "P6-P7 membership/absence",
				"P8-P9 aggregates/top-N", "P10-P11 predicates/pagination"},
		})
		if err != nil {
			return err
		}
		if err := os.WriteFile(reportOut, []byte(report), 0o644); err != nil {
			return err
		}
		fmt.Println("wrote", reportOut)
	}

	fmt.Printf("select: %d cells — valid=%d invalid=%d unsupported=%d timeout=%d\n",
		len(cells), result.Summary["valid"], result.Summary["invalid"],
		result.Summary["unsupported"], result.Summary["timeout"])
	if result.Summary["invalid"] > 0 {
		return fmt.Errorf("%d cell(s) failed the correctness gate", result.Summary["invalid"])
	}
	return nil
}
