package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"charitytree/selectbench"

	"github.com/jackc/pgx/v5/pgxpool"
)

// The blended phase's single fixed writer (review R2 / handoff C5).
//
// It increments an independent maintenance counter on a fixed, deterministic set
// of existing child rows at an offered rate. The counter is not in the catalogue
// projection, not in the rank/index set and not in the MD5 token, so the read
// answers and their oracles are unchanged by the concurrent update. The writer
// is a fixed maintenance load, not an optimistic/pessimistic comparison.

type maintenanceWriter struct {
	pools  []*pgxpool.Pool
	schema string
	table  string
	idCol  string
	keys   []int
	rate   int

	mu        sync.Mutex
	offered   int64
	completed int64
	failed    int64
	latency   []float64
	perKey    map[int]int
	intended  map[int]int
	ops       []int
	lastErr   string

	stop chan struct{}
	wg   sync.WaitGroup
}

func newMaintenanceWriter(pools []*pgxpool.Pool, schema, table, idCol string, keys []int, rate int) *maintenanceWriter {
	return &maintenanceWriter{
		pools: pools, schema: schema, table: table, idCol: idCol, keys: keys, rate: rate,
		perKey: map[int]int{}, intended: map[int]int{}, ops: make([]int, len(pools)),
		stop: make(chan struct{}),
	}
}

// start launches the writer. It must be called before phase warmup so a cell
// cannot be at the start of the phase with no writer running.
func (w *maintenanceWriter) start(ctx context.Context) {
	w.wg.Add(1)
	go func() {
		defer w.wg.Done()
		interval := time.Second / time.Duration(w.rate)
		start := time.Now()
		for i := 0; ; i++ {
			target := start.Add(time.Duration(i) * interval)
			if d := time.Until(target); d > 0 {
				timer := time.NewTimer(d)
				select {
				case <-timer.C:
				case <-w.stop:
					timer.Stop()
					return
				}
			}
			select {
			case <-w.stop:
				return
			default:
			}
			key := w.keys[i%len(w.keys)]
			w.mu.Lock()
			w.offered++
			w.mu.Unlock()
			// Rotate the writer across the run's endpoints exactly as the reader
			// does, so a cluster arm's write path is not silently single-ended
			// (review M-gap: the writer used one endpoint).
			ep := i % len(w.pools)
			t0 := time.Now()
			_, err := w.pools[ep].Exec(ctx,
				fmt.Sprintf("UPDATE %s.%s SET maintenance_updates = maintenance_updates + 1 WHERE %s = $1", w.schema, w.table, w.idCol), key)
			lat := float64(time.Since(t0).Microseconds()) / 1000
			w.mu.Lock()
			w.ops[ep]++
			if err != nil {
				w.failed++
				if w.lastErr == "" {
					w.lastErr = err.Error()
				}
			} else {
				w.completed++
				w.latency = append(w.latency, lat)
				w.perKey[key]++
				w.intended[key]++
			}
			w.mu.Unlock()
		}
	}()
}

func (w *maintenanceWriter) stopAndWait() {
	select {
	case <-w.stop:
	default:
		close(w.stop)
	}
	w.wg.Wait()
}

func (w *maintenanceWriter) completedCount() int64 {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.completed
}

// stats returns the writer's account for the phase. wallMS is the measured phase
// wall time; offered and completed rates are per second over that window.
func (w *maintenanceWriter) stats(phase string, wall float64) *selectbench.WriterStats {
	w.mu.Lock()
	defer w.mu.Unlock()
	s := &selectbench.WriterStats{
		Phase: phase, TargetRate: w.rate,
		Offered: w.offered, Completed: w.completed, Failed: w.failed,
		WallSeconds: wall, LatencyMS: append([]float64(nil), w.latency...),
		PerKey: map[string]int{}, EndpointOps: append([]int(nil), w.ops...), Error: w.lastErr,
	}
	for k, v := range w.perKey {
		s.PerKey[fmt.Sprint(k)] = v
	}
	if wall > 0 {
		s.OfferedRate = float64(w.offered) / wall
		s.CompletedRate = float64(w.completed) / wall
	}
	return s
}

// ledger reconciles the writer's intended per-key increments with the database
// counter values. Only the writer touches those rows in a phase, so any
// mismatch is a real lost/duplicated update, not noise.
func (w *maintenanceWriter) ledger(ctx context.Context) (keys, mismatch int, detail string) {
	w.mu.Lock()
	want := make(map[int]int, len(w.intended))
	maxKey := 0
	for k, v := range w.intended {
		want[k] = v
		if k > maxKey {
			maxKey = k
		}
	}
	w.mu.Unlock()
	if len(want) == 0 {
		return 0, 0, "writer committed no updates in this phase"
	}
	rows, err := w.pools[0].Query(ctx,
		fmt.Sprintf("SELECT %s, maintenance_updates FROM %s.%s WHERE %s <= $1", w.idCol, w.schema, w.table, w.idCol), maxKey)
	if err != nil {
		return 0, 0, "ledger read failed: " + err.Error()
	}
	defer rows.Close()
	for rows.Next() {
		var id, got int
		if err := rows.Scan(&id, &got); err != nil {
			return 0, 0, "ledger scan failed: " + err.Error()
		}
		if w, ok := want[id]; ok {
			keys++
			if got != w {
				mismatch++
				if detail == "" {
					detail = fmt.Sprintf("first mismatch key %d: db=%d writer=%d", id, got, w)
				}
			}
		}
	}
	return keys, mismatch, detail
}

// writerPlan names the real table and key column a blended phase updates for an
// arm. The domain arm updates the donation table itself, not the binding view.
func writerPlan(arm joinArm) (table, idCol string) {
	if arm.dataset == "study01-person-donation-charity" {
		return "donation", "donation_id"
	}
	return "facts", "id"
}
