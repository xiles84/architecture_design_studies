// Package ports defines the boundary between the study-independent core and
// the outside world.
//
// The core (platform/core) and every study's use cases talk only to these
// interfaces. Concrete technology lives in platform/adapters: today that is pgx
// for PostgreSQL-protocol engines and plain files for results. Keeping the
// boundary this thin is deliberate. A benchmark harness that reaches straight
// into a driver ends up measuring the driver's idioms as much as the design
// under test, and cannot be pointed at a second engine without a rewrite.
package ports

import (
	"context"
	"errors"
)

// Isolation is a transaction isolation level, named the way the studies' flags
// name it. The adapter maps it onto whatever the engine calls it.
type Isolation string

const (
	ReadCommitted  Isolation = "rc"
	RepeatableRead Isolation = "rr"
	Serializable   Isolation = "ser"
)

// Row is a single-row result. Scan returns ErrNoRows when the query matched
// nothing, whichever driver sits behind it.
type Row interface {
	Scan(dest ...any) error
}

// Rows is a multi-row result. Callers must drain it: a query whose rows are
// never fetched has not necessarily been executed (LESSONS_LEARNED).
type Rows interface {
	Next() bool
	Scan(dest ...any) error
	Err() error
	Close()
}

// Queryer is what both a pool and a transaction can do.
type Queryer interface {
	// Exec returns the number of rows the statement affected.
	Exec(ctx context.Context, sql string, args ...any) (int64, error)
	QueryRow(ctx context.Context, sql string, args ...any) Row
	Query(ctx context.Context, sql string, args ...any) (Rows, error)
}

// Tx is an open transaction.
type Tx interface {
	Queryer
	Commit(ctx context.Context) error
	// Rollback is safe to call after Commit; it is then a no-op.
	Rollback(ctx context.Context) error
}

// DB is a connection pool to one database under test.
type DB interface {
	Queryer
	Begin(ctx context.Context, iso Isolation) (Tx, error)
	// CopyFrom bulk-loads n rows produced by row(i). Engines without a bulk
	// protocol may fall back to batched inserts; the load phase is timed but is
	// never the thing a study compares.
	CopyFrom(ctx context.Context, table string, cols []string, n int, row func(i int) ([]any, error)) (int64, error)
	// Classify tells the core what an error means without the core having to
	// know which driver produced it.
	Classify(err error) ErrorClass
	Close()
}

// ErrorClass is the part of an error the core is allowed to act on.
type ErrorClass int

const (
	// ErrOther is a real failure: count it, report it, do not retry it.
	ErrOther ErrorClass = iota
	// ErrRetryable is a legitimate consequence of contention -- a serialization
	// failure, a deadlock victim, a read restart. Retrying is correct, and
	// counting the retries is part of the measurement.
	ErrRetryable
	// ErrUniqueViolation means a uniqueness constraint arbitrated a race.
	ErrUniqueViolation
	// ErrCheckViolation means a CHECK constraint rejected the row.
	ErrCheckViolation
	// ErrServerTerminated means the server killed the session. A load is
	// idempotent and may be retried on this; a measurement may not.
	ErrServerTerminated
	// ErrAmbiguousCommit means the connection failed during COMMIT, so the
	// transaction may or may not have committed. An audit must allow for it.
	ErrAmbiguousCommit
)

// ErrNoRows is returned by Row.Scan when the query matched nothing.
var ErrNoRows = errors.New("no rows in result set")
