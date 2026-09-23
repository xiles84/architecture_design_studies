package pgxdb

import (
	"context"
	"fmt"
	"sync/atomic"

	"adsplatform/ports"
)

// spreadPool spreads work over one pool per database node, round robin.
//
// The runner hands a multi-node topology a comma-separated DSN list. A single
// pgxpool cannot parse that list, and connecting every client to one node's SQL
// layer is what Study 02's analysis named as the reason a 3-node comparison said
// little about the cluster: the other nodes did no work. Opening one pool per
// endpoint and spreading operations over them makes the measurement reflect the
// whole cluster. A transaction stays on the pool it began on, so its statements
// keep the engine's own routing.
//
// This lives in the adapter rather than in a study because it is the same
// fidelity problem for every multi-node study; a study that needs it passes the
// endpoint list and depends only on the port.
type spreadPool struct {
	pools []*DB
	next  atomic.Uint64
}

var _ ports.DB = (*spreadPool)(nil)

// OpenSpread opens one pool per endpoint and returns a DB that spreads work over
// them. With a single endpoint it returns that pool unchanged, so single-node
// behaviour is not altered. If any endpoint fails to open, the ones already open
// are closed before the error is returned.
func OpenSpread(ctx context.Context, cfgs []Config) (ports.DB, error) {
	if len(cfgs) == 0 {
		return nil, fmt.Errorf("pgxdb: OpenSpread needs at least one endpoint")
	}
	pools := make([]*DB, 0, len(cfgs))
	for i, c := range cfgs {
		p, err := Open(ctx, c)
		if err != nil {
			for _, q := range pools {
				q.Close()
			}
			return nil, fmt.Errorf("endpoint %d: %w", i, err)
		}
		pools = append(pools, p)
	}
	if len(pools) == 1 {
		return pools[0], nil
	}
	return &spreadPool{pools: pools}, nil
}

// Endpoints reports how many endpoints the spread DB was opened with.
func (s *spreadPool) Endpoints() int { return len(s.pools) }

func (s *spreadPool) pick() *DB { return s.pools[s.next.Add(1)%uint64(len(s.pools))] }

func (s *spreadPool) Exec(ctx context.Context, sql string, args ...any) (int64, error) {
	return s.pick().Exec(ctx, sql, args...)
}

func (s *spreadPool) QueryRow(ctx context.Context, sql string, args ...any) ports.Row {
	return s.pick().QueryRow(ctx, sql, args...)
}

func (s *spreadPool) Query(ctx context.Context, sql string, args ...any) (ports.Rows, error) {
	return s.pick().Query(ctx, sql, args...)
}

func (s *spreadPool) Begin(ctx context.Context, iso ports.Isolation) (ports.Tx, error) {
	return s.pick().Begin(ctx, iso)
}

func (s *spreadPool) CopyFrom(ctx context.Context, table string, cols []string, n int, row func(i int) ([]any, error)) (int64, error) {
	return s.pick().CopyFrom(ctx, table, cols, n, row)
}

// Classify asks any pool: the error classes come from the one engine, so they are
// the same whichever node produced the error.
func (s *spreadPool) Classify(err error) ports.ErrorClass { return s.pools[0].Classify(err) }

func (s *spreadPool) Close() {
	for _, p := range s.pools {
		p.Close()
	}
}
