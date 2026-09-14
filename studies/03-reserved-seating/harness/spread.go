package main

import (
	"context"
	"sync/atomic"

	"adsplatform/ports"
)

// spreadDB spreads work over one pool per database node, round robin. Study 02's
// clients all connected to one YugabyteDB node, whose SQL layer did every query
// and lock -- which its analysis named as the reason 1-node vs 3-node said little.
// A transaction stays on the pool it began on. Depends only on the port.
type spreadDB struct {
	pools []ports.DB
	next  atomic.Uint64
}

var _ ports.DB = (*spreadDB)(nil)

func newSpreadDB(pools []ports.DB) ports.DB {
	if len(pools) == 1 {
		return pools[0]
	}
	return &spreadDB{pools: pools}
}

func (s *spreadDB) pick() ports.DB { return s.pools[s.next.Add(1)%uint64(len(s.pools))] }

func (s *spreadDB) Exec(ctx context.Context, sql string, args ...any) (int64, error) {
	return s.pick().Exec(ctx, sql, args...)
}

func (s *spreadDB) QueryRow(ctx context.Context, sql string, args ...any) ports.Row {
	return s.pick().QueryRow(ctx, sql, args...)
}

func (s *spreadDB) Query(ctx context.Context, sql string, args ...any) (ports.Rows, error) {
	return s.pick().Query(ctx, sql, args...)
}

func (s *spreadDB) Begin(ctx context.Context, iso ports.Isolation) (ports.Tx, error) {
	return s.pick().Begin(ctx, iso)
}

func (s *spreadDB) CopyFrom(ctx context.Context, table string, cols []string, n int, row func(i int) ([]any, error)) (int64, error) {
	return s.pick().CopyFrom(ctx, table, cols, n, row)
}

func (s *spreadDB) Classify(err error) ports.ErrorClass { return s.pools[0].Classify(err) }

func (s *spreadDB) Close() {
	for _, p := range s.pools {
		p.Close()
	}
}
