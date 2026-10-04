package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Balanced cluster access for the topology arms (review R1 / handoff C4).
//
// A cluster arm must not be served by one seed node: a single endpoint stands
// in for balanced cluster access only by accident. rrQuerier holds one pool per
// endpoint and rotates successive statements across them, counting operations
// and reporting connections per endpoint so the result can show the distribution
// was real. With one endpoint (single-node arms) it is a thin wrapper and the
// distribution is trivially that endpoint.

type rrQuerier struct {
	pools       []*pgxpool.Pool
	endpoints   []string
	next        int
	ops         []int
	explainOpts string

	// writer, when set, lets the reader count how many of its operations ran
	// while at least one maintenance update committed. A phase whose reader
	// operations never overlapped a write is not a concurrent measurement.
	writer  *maintenanceWriter
	samples int
	overlap int
}

// resetWriterOverlap scopes the overlap counters to the measured cell phase, so
// gate and feasibility operations do not inflate the denominator.
func (q *rrQuerier) resetWriterOverlap() { q.samples, q.overlap = 0, 0 }

func (q *rrQuerier) writerOverlap() (samples, overlap int) { return q.samples, q.overlap }

func (q *rrQuerier) beginOp() int64 {
	if q.writer == nil {
		return 0
	}
	return q.writer.completedCount()
}

func (q *rrQuerier) endOp(before int64) {
	if q.writer == nil {
		return
	}
	q.samples++
	if q.writer.completedCount() > before {
		q.overlap++
	}
}

func (q *rrQuerier) pick() int {
	i := q.next % len(q.pools)
	q.ops[i]++
	q.next++
	return i
}

func (q *rrQuerier) Exec(ctx context.Context, sql string) error {
	_, err := q.pools[q.pick()].Exec(ctx, sql)
	return classifyEngineError(err)
}

func (q *rrQuerier) Query(ctx context.Context, sqlText string) ([][]any, error) {
	before := q.beginOp()
	rows, err := q.pools[q.pick()].Query(ctx, sqlText)
	if err != nil {
		q.endOp(before)
		return nil, classifyEngineError(err)
	}
	defer rows.Close()
	var out [][]any
	for rows.Next() {
		vals, err := rows.Values()
		if err != nil {
			q.endOp(before)
			return nil, err
		}
		out = append(out, vals)
	}
	q.endOp(before)
	return out, rows.Err()
}

func (q *rrQuerier) Explain(ctx context.Context, sqlText string) (string, error) {
	opt := q.explainOpts
	if opt == "" {
		opt = "ANALYZE, BUFFERS, FORMAT TEXT"
	}
	before := q.beginOp()
	rows, err := q.pools[q.pick()].Query(ctx, "EXPLAIN ("+opt+") "+sqlText)
	if err != nil {
		q.endOp(before)
		return "", classifyEngineError(err)
	}
	defer rows.Close()
	var b strings.Builder
	for rows.Next() {
		var line string
		if err := rows.Scan(&line); err != nil {
			q.endOp(before)
			return "", err
		}
		b.WriteString(line)
		b.WriteString("\n")
	}
	q.endOp(before)
	return b.String(), rows.Err()
}

// conns reports the open connection count per endpoint at the end of the run.
func (q *rrQuerier) conns() []int {
	out := make([]int, len(q.pools))
	for i, p := range q.pools {
		out[i] = int(p.Stat().TotalConns())
	}
	return out
}

// connectEndpoints parses a comma-separated DSN list and builds one pool per
// endpoint with identical session settings. Two connections per endpoint cover
// the sequential reader, endpoint setup and a control; the count is recorded in
// the run's provenance so the pool sizing is not implicit.
func connectEndpoints(ctx context.Context, dsnList, schema string, statementTimeout int) ([]*pgxpool.Pool, []string, error) {
	var pools []*pgxpool.Pool
	var names []string
	for _, d := range strings.Split(dsnList, ",") {
		d = strings.TrimSpace(d)
		if d == "" {
			continue
		}
		cfg, err := pgxpool.ParseConfig(d)
		if err != nil {
			closePools(pools)
			return nil, nil, fmt.Errorf("dsn %q: %w", d, err)
		}
		cfg.ConnConfig.RuntimeParams["search_path"] = schema + ",public"
		cfg.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol
		cfg.ConnConfig.RuntimeParams["statement_timeout"] = fmt.Sprint(statementTimeout)
		cfg.MaxConns = 2
		p, err := pgxpool.NewWithConfig(ctx, cfg)
		if err != nil {
			closePools(pools)
			return nil, nil, fmt.Errorf("connect %s: %w", cfg.ConnConfig.Host, err)
		}
		pools = append(pools, p)
		names = append(names, cfg.ConnConfig.Host)
	}
	if len(pools) == 0 {
		return nil, nil, fmt.Errorf("no database endpoints in %q", dsnList)
	}
	return pools, names, nil
}

func closePools(pools []*pgxpool.Pool) {
	for _, p := range pools {
		p.Close()
	}
}
