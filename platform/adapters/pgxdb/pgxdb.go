// Package pgxdb adapts pgx to ports.DB for every PostgreSQL-protocol engine the
// studies run: PostgreSQL itself and YugabyteDB's YSQL layer.
//
// It is the only package in the repository that imports a database driver.
package pgxdb

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

	"adsplatform/ports"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	DSN                string
	MaxConns           int
	StatementTimeoutMS int
	ApplicationName    string
}

type DB struct{ pool *pgxpool.Pool }

var _ ports.DB = (*DB)(nil)

func Open(ctx context.Context, c Config) (*DB, error) {
	cfg, err := pgxpool.ParseConfig(c.DSN)
	if err != nil {
		return nil, err
	}
	cfg.MaxConns = int32(max(c.MaxConns, 2))
	cfg.MinConns = 1
	if c.StatementTimeoutMS > 0 {
		cfg.ConnConfig.RuntimeParams["statement_timeout"] = fmt.Sprint(c.StatementTimeoutMS)
	}
	if c.ApplicationName != "" {
		cfg.ConnConfig.RuntimeParams["application_name"] = c.ApplicationName
	}
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}
	return &DB{pool: pool}, nil
}

func (d *DB) Close() { d.pool.Close() }

func (d *DB) Exec(ctx context.Context, sql string, args ...any) (int64, error) {
	tag, err := d.pool.Exec(ctx, sql, args...)
	return tag.RowsAffected(), err
}

func (d *DB) QueryRow(ctx context.Context, sql string, args ...any) ports.Row {
	return row{d.pool.QueryRow(ctx, sql, args...)}
}

func (d *DB) Query(ctx context.Context, sql string, args ...any) (ports.Rows, error) {
	return d.pool.Query(ctx, sql, args...)
}

func (d *DB) Begin(ctx context.Context, iso ports.Isolation) (ports.Tx, error) {
	tx, err := d.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: isoLevel(iso)})
	if err != nil {
		return nil, err
	}
	return &txn{tx: tx}, nil
}

// CopyFrom uses the COPY protocol. Callers that want parallel streams split the
// row range themselves; see ParallelCopy.
func (d *DB) CopyFrom(ctx context.Context, table string, cols []string, n int, rowFn func(i int) ([]any, error)) (int64, error) {
	i := 0
	return d.pool.CopyFrom(ctx, pgx.Identifier{table}, cols, pgx.CopyFromFunc(func() ([]any, error) {
		if i >= n {
			return nil, nil
		}
		v, err := rowFn(i)
		i++
		return v, err
	}))
}

// Classify maps SQLSTATEs onto the classes the core acts on.
func (d *DB) Classify(err error) ports.ErrorClass { return Classify(err) }

func Classify(err error) ports.ErrorClass {
	if err == nil {
		return ports.ErrOther
	}
	var ae *ambiguousCommit
	if errors.As(err, &ae) {
		return ports.ErrAmbiguousCommit
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "40001", "40P01":
			return ports.ErrRetryable
		case "23505":
			return ports.ErrUniqueViolation
		case "23514":
			return ports.ErrCheckViolation
		case "57P01", "57P02", "57P03", "08006", "08003":
			return ports.ErrServerTerminated
		case "55P03":
			// lock_not_available: FOR UPDATE NOWAIT met a locked row.
			return ports.ErrLockNotAvailable
		}
		// YugabyteDB reports some transaction conflicts with a generic code and
		// a descriptive message. They are contention, not failure.
		msg := strings.ToLower(pgErr.Message)
		if strings.Contains(msg, "restart read") || strings.Contains(msg, "try again") ||
			strings.Contains(msg, "conflicts with") || strings.Contains(msg, "transaction aborted") ||
			strings.Contains(msg, "could not serialize") {
			return ports.ErrRetryable
		}
		return ports.ErrOther
	}
	msg := err.Error()
	if strings.Contains(msg, "conn closed") || strings.Contains(msg, "unexpected EOF") ||
		strings.Contains(msg, "connection reset") {
		return ports.ErrServerTerminated
	}
	return ports.ErrOther
}

func isoLevel(i ports.Isolation) pgx.TxIsoLevel {
	switch i {
	case ports.RepeatableRead:
		return pgx.RepeatableRead
	case ports.Serializable:
		return pgx.Serializable
	default:
		return pgx.ReadCommitted
	}
}

type row struct{ r pgx.Row }

func (r row) Scan(dest ...any) error {
	err := r.r.Scan(dest...)
	if errors.Is(err, pgx.ErrNoRows) {
		return ports.ErrNoRows
	}
	return err
}

type txn struct {
	tx   pgx.Tx
	once sync.Once
}

func (t *txn) Exec(ctx context.Context, sql string, args ...any) (int64, error) {
	tag, err := t.tx.Exec(ctx, sql, args...)
	return tag.RowsAffected(), err
}

func (t *txn) QueryRow(ctx context.Context, sql string, args ...any) ports.Row {
	return row{t.tx.QueryRow(ctx, sql, args...)}
}

func (t *txn) Query(ctx context.Context, sql string, args ...any) (ports.Rows, error) {
	return t.tx.Query(ctx, sql, args...)
}

// Commit distinguishes a commit the server refused (a real, classifiable error)
// from a commit whose outcome is unknown because the connection broke. The
// second kind must be counted separately by any audit that reconciles
// "operations the client saw succeed" against "rows that exist".
func (t *txn) Commit(ctx context.Context) error {
	err := t.tx.Commit(ctx)
	if err == nil {
		return nil
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return err
	}
	if pgconn.SafeToRetry(err) {
		return err
	}
	return &ambiguousCommit{err}
}

func (t *txn) Rollback(ctx context.Context) error {
	err := t.tx.Rollback(ctx)
	if errors.Is(err, pgx.ErrTxClosed) {
		return nil
	}
	return err
}

type ambiguousCommit struct{ err error }

func (a *ambiguousCommit) Error() string { return "commit outcome unknown: " + a.err.Error() }
func (a *ambiguousCommit) Unwrap() error { return a.err }
