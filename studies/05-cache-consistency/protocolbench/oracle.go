package main

import (
	"context"
	"fmt"
	"sync"
)

type state struct {
	Generation int64 `json:"generation"`
	Value      int64 `json:"value"`
	Deleted    bool  `json:"deleted"`
}

func expected(key, g int64, deleted bool) state {
	v := key*1000 + g
	if deleted {
		v = 0
	}
	return state{g, v, deleted}
}

type pending struct {
	value     state
	done      chan struct{}
	committed bool
}
type ledger struct {
	mu      sync.Mutex
	history map[int64]map[int64]state
	pending map[int64]map[int64]*pending
	ack     map[int64]int64
	current map[int64]state
}

func newLedger(n int) *ledger {
	l := &ledger{history: map[int64]map[int64]state{}, pending: map[int64]map[int64]*pending{}, ack: map[int64]int64{}, current: map[int64]state{}}
	for k := int64(1); k <= int64(n); k++ {
		l.history[k] = map[int64]state{1: expected(k, 1, false)}
		l.current[k] = expected(k, 1, false)
		l.pending[k] = map[int64]*pending{}
		l.ack[k] = 1
	}
	return l
}
func (l *ledger) required(k int64) int64 { l.mu.Lock(); defer l.mu.Unlock(); return l.ack[k] }
func (l *ledger) latest(k int64) state   { l.mu.Lock(); defer l.mu.Unlock(); return l.current[k] }
func (l *ledger) begin(k int64, s state) *pending {
	l.mu.Lock()
	defer l.mu.Unlock()
	p := &pending{value: s, done: make(chan struct{})}
	l.pending[k][s.Generation] = p
	return p
}
func (l *ledger) finish(k int64, p *pending, committed bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	p.committed = committed
	if committed {
		l.history[k][p.value.Generation] = p.value
		l.current[k] = p.value
	}
	close(p.done)
}
func (l *ledger) acknowledge(k int64, g int64) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if g > l.ack[k] {
		l.ack[k] = g
	}
}

// An ahead database snapshot waits for the independently recorded commit outcome;
// a merely intended or rolled-back state is never admitted into committed history.
func (l *ledger) classify(ctx context.Context, k int64, s state, required int64) (string, error) {
	l.mu.Lock()
	known, ok := l.history[k][s.Generation]
	p := l.pending[k][s.Generation]
	l.mu.Unlock()
	if !ok && p != nil {
		select {
		case <-p.done:
		case <-ctx.Done():
			return "", ctx.Err()
		}
		l.mu.Lock()
		known, ok = l.history[k][s.Generation]
		l.mu.Unlock()
	}
	if !ok || known != s {
		return "impossible", nil
	}
	if s.Generation < required {
		return "stale-after-ack", nil
	}
	return "committed", nil
}
func requireClass(ctx context.Context, l *ledger, k int64, s state, want string) error {
	got, err := l.classify(ctx, k, s, l.required(k))
	if err != nil {
		return err
	}
	if got != want {
		return fmt.Errorf("oracle: want %s got %s state=%+v", want, got, s)
	}
	return nil
}
