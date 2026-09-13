// Package inspect captures what the engine did, as opposed to how long it took:
// plans, versions, effective settings. It talks to the database only through
// the port, so it is shared by every study.
package inspect

import (
	"context"
	"fmt"
	"strings"
	"time"

	"adsplatform/core/catalog"
	"adsplatform/ports"
)

// ExplainPrefix picks the richest EXPLAIN the engine understands.
//
//	PostgreSQL: BUFFERS attributes work to block reads -- wall-clock alone
//	            confuses "less work" with "warmer cache".
//	YugabyteDB: DIST reports RPC counts, the cost signal that survives the move
//	            from a laptop's loopback to a real network.
func ExplainPrefix(engine string) string {
	if engine == "yugabyte" {
		return "EXPLAIN (ANALYZE, DIST, COSTS OFF, TIMING OFF)"
	}
	return "EXPLAIN (ANALYZE, BUFFERS, VERBOSE, COSTS ON)"
}

// Explain captures a plan for every statement, executed twice so the kept plan
// reflects a warm cache rather than first touch.
//
// When rollback is true each statement runs inside its own transaction that is
// rolled back, which is how write statements get real ANALYZE output without
// mutating the dataset the benchmark is about to measure. A statement that
// cannot be explained with the representative values (because it depends on a
// row an earlier step would have produced) records its error as its plan: a
// missing plan should be visible, not silently skipped.
func Explain(ctx context.Context, db ports.DB, engine string, stmts []catalog.Stmt, vals map[string]any, rollback bool) map[string]string {
	out := map[string]string{}
	for _, st := range stmts {
		args, err := catalog.Bind(st.Params, vals)
		if err != nil {
			out[st.Name] = "NOT CAPTURED: " + err.Error()
			continue
		}
		q := ExplainPrefix(engine) + " " + st.SQL
		var plan string
		for attempt := 0; attempt < 2; attempt++ {
			plan = explainOnce(ctx, db, q, args, rollback)
		}
		out[st.Name] = strings.TrimRight(plan, "\n")
	}
	return out
}

// ExplainSequence captures plans for statements that depend on each other (a
// parent insert, then its children) by explaining them in order inside ONE
// transaction that is rolled back at the end. Each is executed once, since a
// second execution would collide with the first.
func ExplainSequence(ctx context.Context, db ports.DB, engine string, stmts []catalog.Stmt, vals map[string]any) map[string]string {
	out := map[string]string{}
	tx, err := db.Begin(ctx, ports.ReadCommitted)
	if err != nil {
		for _, st := range stmts {
			out[st.Name] = "ERROR: " + err.Error()
		}
		return out
	}
	defer tx.Rollback(ctx)
	failed := ""
	for _, st := range stmts {
		if failed != "" {
			out[st.Name] = "NOT CAPTURED: an earlier statement in the sequence failed (" + failed + ")"
			continue
		}
		args, err := catalog.Bind(st.Params, vals)
		if err != nil {
			out[st.Name] = "NOT CAPTURED: " + err.Error()
			continue
		}
		plan := queryLines(ctx, tx, ExplainPrefix(engine)+" "+st.SQL, args)
		if strings.HasPrefix(plan, "ERROR:") {
			failed = st.Name
		}
		out[st.Name] = strings.TrimRight(plan, "\n")
	}
	return out
}

func explainOnce(ctx context.Context, db ports.DB, q string, args []any, rollback bool) string {
	var qr ports.Queryer = db
	if rollback {
		tx, err := db.Begin(ctx, ports.ReadCommitted)
		if err != nil {
			return "ERROR: " + err.Error()
		}
		defer tx.Rollback(ctx)
		qr = tx
	}
	return queryLines(ctx, qr, q, args)
}

func queryLines(ctx context.Context, qr ports.Queryer, q string, args []any) string {
	rows, err := qr.Query(ctx, q, args...)
	if err != nil {
		return "ERROR: " + err.Error()
	}
	defer rows.Close()
	var sb strings.Builder
	for rows.Next() {
		var line string
		if err := rows.Scan(&line); err != nil {
			return "ERROR: " + err.Error()
		}
		sb.WriteString(line + "\n")
	}
	if err := rows.Err(); err != nil {
		return "ERROR: " + err.Error()
	}
	return sb.String()
}

// EngineInfo records exactly what was measured. "PostgreSQL 17" is not a
// result; the version banner is. EffectiveIsolation matters on YugabyteDB,
// where depending on version and server flags READ COMMITTED may run as Snapshot
// Isolation -- a study that believes it measured RC must be able to prove it.
type EngineInfo struct {
	Version            string `json:"version"`
	DefaultIsolation   string `json:"default_isolation"`
	EffectiveIsolation string `json:"effective_isolation,omitempty"`
}

func Describe(ctx context.Context, db ports.DB, engine string) EngineInfo {
	var info EngineInfo
	if err := db.QueryRow(ctx, "SELECT version()").Scan(&info.Version); err != nil {
		info.Version = "unknown: " + err.Error()
	}
	info.Version = strings.TrimSpace(info.Version)
	_ = db.QueryRow(ctx, "SHOW default_transaction_isolation").Scan(&info.DefaultIsolation)
	if engine == "yugabyte" {
		// Asked inside a READ COMMITTED transaction, because that is the
		// question: when a design requests RC, what does it get?
		tx, err := db.Begin(ctx, ports.ReadCommitted)
		if err == nil {
			if err := tx.QueryRow(ctx, "SHOW yb_effective_transaction_isolation_level").Scan(&info.EffectiveIsolation); err != nil {
				info.EffectiveIsolation = "unknown: " + err.Error()
			}
			_ = tx.Rollback(ctx)
		}
	}
	return info
}

// RetryOnTermination repeats an idempotent setup step when -- and only when --
// the server killed the session. On the 3-node YugabyteDB cluster, repeated
// drop/create/COPY/backfill cycles occasionally get a session terminated
// (57P01) mid-DDL; that says nothing about a design and once cost a whole cell.
// A schema error, a failed verification or a timeout is a real result and is
// returned immediately.
func RetryOnTermination[T any](ctx context.Context, db ports.DB, attempts int, fn func() (T, error)) (T, error) {
	var zero T
	var last error
	for i := 1; i <= attempts; i++ {
		v, err := fn()
		if err == nil {
			if i > 1 {
				fmt.Printf("  succeeded on attempt %d after: %v\n", i, last)
			}
			return v, nil
		}
		last = err
		if db.Classify(err) != ports.ErrServerTerminated {
			return zero, err
		}
		fmt.Printf("  attempt %d/%d lost its connection (%v); retrying\n", i, attempts, err)
		select {
		case <-ctx.Done():
			return zero, ctx.Err()
		case <-time.After(time.Duration(i*10) * time.Second):
		}
	}
	return zero, fmt.Errorf("failed after %d attempts: %w", attempts, last)
}
