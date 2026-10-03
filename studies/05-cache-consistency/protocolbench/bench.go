package main

import (
	"context"
	"embed"
	"fmt"
	"math/rand"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"adsplatform/core/catalog"
	"adsplatform/ports"
)

//go:embed sql/*.sql
var sqlFiles embed.FS

type counters struct{ snapshots, validations, relayQueries, deliveries atomic.Int64 }
type bench struct {
	db                       ports.DB
	cache                    *redisStore
	sql                      map[string]catalog.Stmt
	ledger                   *ledger
	protocol, strategy, kind string
	keys                     int
	locks                    []sync.Mutex
	stats                    *counters
	relayCancel              context.CancelFunc
	relayDone                chan struct{}
	relayMu                  sync.Mutex
	paused                   bool
	deliveryMu               sync.Mutex
	delivered                map[int64]int64
	relayError               error
	skipBarrier              bool
	disableValidation        bool
}

func newBench(db ports.DB, r *redisStore, protocol, strategy, kind string, keys int) (*bench, error) {
	src, e := sqlFiles.ReadFile("sql/queries.sql")
	if e != nil {
		return nil, e
	}
	ss, e := catalog.Parse(string(src))
	if e != nil {
		return nil, e
	}
	return &bench{db: db, cache: r, sql: catalog.Map(ss), protocol: protocol, strategy: strategy, kind: kind, keys: keys, locks: make([]sync.Mutex, keys+1)}, nil
}
func (b *bench) q(n string) string { return b.sql[n].SQL }
func (b *bench) stopRelay() {
	if b.relayCancel != nil {
		b.relayCancel()
		<-b.relayDone
		b.relayCancel = nil
	}
}
func (b *bench) reset(ctx context.Context) error {
	b.stopRelay()
	src, e := sqlFiles.ReadFile("sql/schema.sql")
	if e != nil {
		return e
	}
	for _, s := range catalog.SplitDDL(string(src)) {
		if _, e = b.db.Exec(ctx, s); e != nil {
			return e
		}
	}
	if _, e = b.db.Exec(ctx, "INSERT INTO protocol_generation SELECT i,1 FROM generate_series(1,$1::int) i", b.keys); e != nil {
		return e
	}
	if _, e = b.db.Exec(ctx, "INSERT INTO protocol_object SELECT i,1,i*1000+1,false FROM generate_series(1,$1::int) i", b.keys); e != nil {
		return e
	}
	if _, e = b.cache.do(ctx, "FLUSHDB"); e != nil {
		return e
	}
	b.ledger = newLedger(b.keys)
	b.stats = &counters{}
	b.delivered = map[int64]int64{}
	b.relayError = nil
	b.paused = false
	b.skipBarrier = false
	b.disableValidation = false
	if b.protocol == "outbox" {
		rctx, cancel := context.WithCancel(ctx)
		b.relayCancel = cancel
		b.relayDone = make(chan struct{})
		go b.relay(rctx)
	}
	return nil
}
func (b *bench) snapshot(ctx context.Context, k int64) (state, error) {
	var s state
	b.stats.snapshots.Add(1)
	e := b.db.QueryRow(ctx, b.q("snapshot"), k).Scan(&s.Generation, &s.Value, &s.Deleted)
	return s, e
}
func (b *bench) read(ctx context.Context, k int64) (state, error) {
	required := b.ledger.required(k)
	s, hit, e := b.cache.get(ctx, k)
	if e != nil {
		return s, e
	}
	if hit && b.protocol == "version" && !b.disableValidation {
		var g int64
		b.stats.validations.Add(1)
		e = b.db.QueryRow(ctx, b.q("version"), k).Scan(&g)
		if e != nil {
			return s, e
		}
		hit = g == s.Generation
	}
	if !hit {
		f, e := b.cache.fence(ctx, k)
		if e != nil {
			return s, e
		}
		s, e = b.snapshot(ctx, k)
		if e != nil {
			return s, e
		}
		ttl := 300 * time.Second
		if s.Deleted && b.kind == "tombstone" {
			ttl = 100 * time.Millisecond
		}
		if _, e = b.cache.put(ctx, k, s, f, ttl); e != nil {
			return s, e
		}
	}
	class, e := b.ledger.classify(ctx, k, s, required)
	if e != nil {
		return s, e
	}
	if class != "committed" {
		return s, fmt.Errorf("%s key=%d generation=%d required=%d", class, k, s.Generation, required)
	}
	return s, nil
}

type writeTimes struct{ DBMS, AckMS, LockMS float64 }

func (b *bench) write(ctx context.Context, k int64, deleted bool) (writeTimes, error) {
	var ts writeTimes
	start := time.Now()
	b.locks[k].Lock()
	defer b.locks[k].Unlock()
	ts.LockMS = float64(time.Since(start).Microseconds()) / 1000
	old := b.ledger.latest(k)
	next := expected(k, old.Generation+1, deleted)
	if b.protocol == "sync" {
		if e := b.cache.invalidate(ctx, k); e != nil {
			return ts, e
		}
	}
	p := b.ledger.begin(k, next)
	committed := false
	defer func() {
		if !committed {
			b.ledger.finish(k, p, false)
		}
	}()
	dbStart := time.Now()
	tx, e := b.db.Begin(ctx, ports.ReadCommitted)
	if e != nil {
		return ts, e
	}
	defer tx.Rollback(context.Background())
	var g int64
	if b.strategy == "lock" {
		if e = tx.QueryRow(ctx, b.q("lock"), k).Scan(&g); e != nil {
			return ts, e
		}
		if g != old.Generation {
			return ts, fmt.Errorf("ledger/source generation mismatch")
		}
		e = tx.QueryRow(ctx, b.q("bump"), k).Scan(&g)
	} else {
		e = tx.QueryRow(ctx, b.q("cas"), k, old.Generation).Scan(&g)
	}
	if e != nil {
		return ts, e
	}
	if g != next.Generation {
		return ts, fmt.Errorf("generation advance mismatch")
	}
	if deleted && b.kind != "soft" {
		_, e = tx.Exec(ctx, b.q("delete"), k)
	} else {
		_, e = tx.Exec(ctx, b.q("upsert"), k, g, next.Value, deleted)
	}
	if e != nil {
		return ts, e
	}
	if _, e = tx.Exec(ctx, b.q("event"), k, g); e != nil {
		return ts, e
	}
	if e = tx.Commit(ctx); e != nil {
		return ts, e
	}
	ts.DBMS = float64(time.Since(dbStart).Microseconds()) / 1000
	commitTime := time.Now()
	b.ledger.finish(k, p, true)
	committed = true
	if b.protocol == "sync" {
		if e = b.cache.invalidate(ctx, k); e != nil {
			return ts, e
		}
	}
	if b.protocol == "outbox" && !b.skipBarrier {
		for {
			b.deliveryMu.Lock()
			done := b.delivered[k] >= g
			re := b.relayError
			b.deliveryMu.Unlock()
			if re != nil {
				return ts, re
			}
			if done {
				break
			}
			select {
			case <-ctx.Done():
				return ts, ctx.Err()
			case <-time.After(200 * time.Microsecond):
			}
		}
	}
	b.ledger.acknowledge(k, g)
	ts.AckMS = float64(time.Since(commitTime).Microseconds()) / 1000
	return ts, nil
}
func (b *bench) pauseRelay(paused bool) { b.relayMu.Lock(); b.paused = paused; b.relayMu.Unlock() }
func (b *bench) relay(ctx context.Context) {
	defer close(b.relayDone)
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		b.relayMu.Lock()
		if !b.paused {
			if e := b.deliver(ctx); e != nil && ctx.Err() == nil {
				b.deliveryMu.Lock()
				b.relayError = e
				b.deliveryMu.Unlock()
				b.relayMu.Unlock()
				return
			}
		}
		b.relayMu.Unlock()
		select {
		case <-ctx.Done():
			return
		case <-time.After(2 * time.Millisecond):
		}
	}
}

type event struct{ seq, key, gen int64 }

func (b *bench) deliver(ctx context.Context) error {
	tx, e := b.db.Begin(ctx, ports.ReadCommitted)
	if e != nil {
		return e
	}
	defer tx.Rollback(context.Background())
	b.stats.relayQueries.Add(1)
	rows, e := tx.Query(ctx, b.q("pending"))
	if e != nil {
		return e
	}
	var events []event
	for rows.Next() {
		var ev event
		if e = rows.Scan(&ev.seq, &ev.key, &ev.gen); e != nil {
			rows.Close()
			return e
		}
		events = append(events, ev)
	}
	e = rows.Err()
	rows.Close()
	if e != nil {
		return e
	}
	for _, ev := range events {
		if e = b.cache.invalidate(ctx, ev.key); e != nil {
			return e
		}
		if _, e = tx.Exec(ctx, b.q("delivered"), ev.seq); e != nil {
			return e
		}
	}
	if e = tx.Commit(ctx); e != nil {
		return e
	}
	b.deliveryMu.Lock()
	for _, ev := range events {
		if ev.gen > b.delivered[ev.key] {
			b.delivered[ev.key] = ev.gen
		}
		b.stats.deliveries.Add(1)
	}
	b.deliveryMu.Unlock()
	return nil
}
func (b *bench) audit(ctx context.Context) error {
	for k := int64(1); k <= int64(b.keys); k++ {
		s, e := b.snapshot(ctx, k)
		if e != nil {
			return e
		}
		if s != b.ledger.latest(k) {
			return fmt.Errorf("source/ledger disagreement at key %d", k)
		}
		if _, e = b.read(ctx, k); e != nil {
			return e
		}
	}
	return nil
}

type control struct {
	Name            string  `json:"name"`
	Passed          bool    `json:"passed"`
	Detail          string  `json:"detail"`
	BlockedMS       float64 `json:"blocked_ms,omitempty"`
	RecoveryMS      float64 `json:"recovery_ms,omitempty"`
	StaleDurationMS float64 `json:"stale_duration_ms,omitempty"`
}

func (b *bench) controls(ctx context.Context) ([]control, error) {
	var out []control
	check := func(name string, ok bool, detail string) error {
		out = append(out, control{Name: name, Passed: ok, Detail: detail})
		if !ok {
			return fmt.Errorf("control did not fire: %s: %s", name, detail)
		}
		return nil
	}
	if e := b.reset(ctx); e != nil {
		return out, e
	}
	old, e := b.read(ctx, 1)
	if e != nil {
		return out, e
	}
	f, e := b.cache.fence(ctx, 1)
	if e != nil {
		return out, e
	}
	if _, e = b.write(ctx, 1, false); e != nil {
		return out, e
	}
	accepted, e := b.cache.put(ctx, 1, old, f, 300*time.Second)
	if e != nil {
		return out, e
	}
	s, re := b.read(ctx, 1)
	// Version validation may admit an old committed fill, but must reject it at serve time.
	if e = check("paused-old-fill", re == nil && s.Generation == 2, fmt.Sprintf("old_publish=%t served_generation=%d", accepted, s.Generation)); e != nil {
		return out, e
	}
	if e = b.reset(ctx); e != nil {
		return out, e
	}
	forged := expected(1, 1, false)
	forged.Value++
	if _, e = b.cache.do(ctx, "SET", key(1), encode(forged), "PX", "300000"); e != nil {
		return out, e
	}
	_, re = b.read(ctx, 1)
	if e = check("uncommitted-value", re != nil && stringsHas(re.Error(), "impossible"), fmt.Sprint(re)); e != nil {
		return out, e
	}
	if e = b.reset(ctx); e != nil {
		return out, e
	}
	old, e = b.read(ctx, 1)
	if e != nil {
		return out, e
	}
	if _, e = b.write(ctx, 1, false); e != nil {
		return out, e
	}
	if _, e = b.cache.do(ctx, "SET", key(1), encode(old), "PX", "300000"); e != nil {
		return out, e
	}
	b.disableValidation = true
	_, re = b.read(ctx, 1)
	if e = check("stale-read-with-protocol-disabled", re != nil && stringsHas(re.Error(), "stale-after-ack"), fmt.Sprint(re)); e != nil {
		return out, e
	}
	if e = b.reset(ctx); e != nil {
		return out, e
	}
	if _, e = b.db.Exec(ctx, "UPDATE protocol_generation SET generation=2 WHERE key_id=1"); e != nil {
		return out, e
	}
	if _, e = b.db.Exec(ctx, "UPDATE protocol_generation SET generation=2 WHERE key_id=1"); e != nil {
		return out, e
	}
	var actual int64
	if e = b.db.QueryRow(ctx, b.q("version"), int64(1)).Scan(&actual); e != nil {
		return out, e
	}
	if e = check("unguarded-lost-update", actual != 3, fmt.Sprintf("two intended increments: expected=3 actual=%d", actual)); e != nil {
		return out, e
	}
	if b.protocol == "outbox" {
		if e = b.reset(ctx); e != nil {
			return out, e
		}
		if _, e = b.read(ctx, 1); e != nil {
			return out, e
		}
		b.pauseRelay(true)
		done := make(chan error, 1)
		go func() { _, we := b.write(ctx, 1, false); done <- we }()
		limit := time.Now().Add(2 * time.Second)
		for b.ledger.latest(1).Generation < 2 && time.Now().Before(limit) {
			time.Sleep(time.Millisecond)
		}
		if b.ledger.latest(1).Generation != 2 {
			b.pauseRelay(false)
			return out, fmt.Errorf("relay barrier control never reached source commit")
		}
		blockStart := time.Now()
		blocked := false
		select {
		case e = <-done:
			if e != nil {
				return out, e
			}
		case <-time.After(50 * time.Millisecond):
			blocked = true
		}
		recoveryStart := time.Now()
		b.pauseRelay(false)
		if blocked {
			select {
			case e = <-done:
			case <-time.After(2 * time.Second):
				e = fmt.Errorf("relay failed to recover")
			}
			if e != nil {
				return out, e
			}
		}
		s, re = b.read(ctx, 1)
		cr := control{Name: "relay-acknowledgement-barrier-and-recovery", Passed: blocked && re == nil && s.Generation == 2, Detail: fmt.Sprintf("served generation=%d", s.Generation), BlockedMS: float64(recoveryStart.Sub(blockStart).Microseconds()) / 1000, RecoveryMS: float64(time.Since(recoveryStart).Microseconds()) / 1000}
		out = append(out, cr)
		if !cr.Passed {
			return out, fmt.Errorf("relay barrier did not hold")
		}
		if e = b.reset(ctx); e != nil {
			return out, e
		}
		if _, e = b.read(ctx, 1); e != nil {
			return out, e
		}
		b.pauseRelay(true)
		b.skipBarrier = true
		if _, e = b.write(ctx, 1, false); e != nil {
			return out, e
		}
		staleStart := time.Now()
		_, re = b.read(ctx, 1)
		if e = check("outbox-without-ack-barrier", re != nil && stringsHas(re.Error(), "stale-after-ack"), fmt.Sprint(re)); e != nil {
			return out, e
		}
		time.Sleep(30 * time.Millisecond)
		b.pauseRelay(false)
		limit = time.Now().Add(2 * time.Second)
		for time.Now().Before(limit) {
			_, re = b.read(ctx, 1)
			if re == nil {
				break
			}
			time.Sleep(time.Millisecond)
		}
		if re != nil {
			return out, re
		}
		out[len(out)-1].StaleDurationMS = float64(time.Since(staleStart).Microseconds()) / 1000
	}
	return out, b.reset(ctx)
}
func stringsHas(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
func (b *bench) deletionControls(ctx context.Context) ([]control, error) {
	var out []control
	add := func(n string, ok bool, d string) error {
		out = append(out, control{Name: n, Passed: ok, Detail: d})
		if !ok {
			return fmt.Errorf("deletion gate failed: %s", n)
		}
		return nil
	}
	if e := b.reset(ctx); e != nil {
		return out, e
	}
	old, e := b.snapshot(ctx, 1)
	if e != nil {
		return out, e
	}
	f, e := b.cache.fence(ctx, 1)
	if e != nil {
		return out, e
	}
	if _, e = b.write(ctx, 1, true); e != nil {
		return out, e
	}
	accepted, e := b.cache.put(ctx, 1, old, f, 300*time.Second)
	if e != nil {
		return out, e
	}
	s, re := b.read(ctx, 1)
	if e = add("old-fill-after-delete", !accepted && re == nil && s.Deleted, "durable generation and fences rejected old fill"); e != nil {
		return out, e
	}
	if b.kind == "tombstone" {
		limit := time.Now().Add(2 * time.Second)
		expired := false
		for time.Now().Before(limit) {
			_, hit, ce := b.cache.get(ctx, 1)
			if ce != nil {
				return out, ce
			}
			if !hit {
				expired = true
				break
			}
			time.Sleep(5 * time.Millisecond)
		}
		if e = add("absence-marker-expired", expired, "GET observed the 100-ms marker absent"); e != nil {
			return out, e
		}
	}
	accepted, e = b.cache.put(ctx, 1, old, f, 300*time.Second)
	if e != nil {
		return out, e
	}
	if e = add("old-fill-after-marker-window", !accepted, "fencing metadata remains present after the marker window"); e != nil {
		return out, e
	}
	if _, e = b.write(ctx, 1, false); e != nil {
		return out, e
	}
	accepted, e = b.cache.put(ctx, 1, old, f, 300*time.Second)
	if e != nil {
		return out, e
	}
	s, re = b.read(ctx, 1)
	if e = add("old-fill-after-recreate", !accepted && re == nil && s.Generation == 3 && !s.Deleted, fmt.Sprintf("served generation=%d", s.Generation)); e != nil {
		return out, e
	}
	if e = b.reset(ctx); e != nil {
		return out, e
	}
	old, e = b.snapshot(ctx, 1)
	if e != nil {
		return out, e
	}
	if _, e = b.write(ctx, 1, true); e != nil {
		return out, e
	}
	if _, e = b.cache.do(ctx, "PEXPIRE", fenceKey(1), "30"); e != nil {
		return out, e
	}
	limit := time.Now().Add(time.Second)
	expired := false
	for time.Now().Before(limit) {
		ff, ce := b.cache.fence(ctx, 1)
		if ce != nil {
			return out, ce
		}
		if ff == 0 {
			expired = true
			break
		}
		time.Sleep(5 * time.Millisecond)
	}
	accepted, e = b.cache.put(ctx, 1, old, 0, 300*time.Second)
	if e != nil {
		return out, e
	}
	_, re = b.read(ctx, 1)
	if e = add("expiring-fence-resurrection-control", expired && accepted && re != nil && stringsHas(re.Error(), "stale-after-ack"), fmt.Sprint(re)); e != nil {
		return out, e
	}
	return out, b.reset(ctx)
}

type trial struct {
	Phase                                                                                string  `json:"phase"`
	Number                                                                               int     `json:"trial"`
	ElapsedS                                                                             float64 `json:"elapsed_s"`
	Reads, Writes                                                                        int64
	ReadOps, WriteOps                                                                    float64
	ReadP50MS, ReadP99MS, WriteP50MS, WriteP99MS, DBP50MS, AckP50MS, AckP99MS, LockP99MS float64
	SnapshotQueries, ValidationQueries, RelayQueries, DeliveredEvents                    int64
	SourceQueriesPerRead                                                                 float64
	WrongReads, ImpossibleReads, Errors                                                  int64
}
type localStats struct {
	reads, writes, errors, wrong, impossible  int64
	readLat, writeLat, dbLat, ackLat, lockLat []float64
	err                                       error
}

func percentile(xs []float64, q float64) float64 {
	if len(xs) == 0 {
		return 0
	}
	sort.Float64s(xs)
	i := int(float64(len(xs)-1) * q)
	return xs[i]
}
func (b *bench) measure(ctx context.Context, phase string, number, readers, writers int, duration time.Duration, seed int64) (trial, error) {
	if e := b.reset(ctx); e != nil {
		return trial{}, e
	}
	for k := int64(1); k <= int64(b.keys); k++ {
		if _, e := b.read(ctx, k); e != nil {
			return trial{}, e
		}
	}
	s0, v0, r0, d0 := b.stats.snapshots.Load(), b.stats.validations.Load(), b.stats.relayQueries.Load(), b.stats.deliveries.Load()
	nread, nwrite := readers, writers
	if phase == "read" {
		nwrite = 0
	}
	if phase == "write" || phase == "delete-recreate" {
		nread = 0
	}
	timer, stop := context.WithTimeout(ctx, duration)
	defer stop()
	start := time.Now()
	ch := make(chan localStats, nread+nwrite)
	for worker := 0; worker < nread+nwrite; worker++ {
		go func(w int) {
			var l localStats
			r := rand.New(rand.NewSource(seed + int64(w)*7919))
			for timer.Err() == nil {
				k := int64(r.Intn(b.keys) + 1)
				t := time.Now()
				if w < nread {
					_, e := b.read(ctx, k)
					if e != nil {
						l.err = e
						break
					}
					l.reads++
					l.readLat = append(l.readLat, float64(time.Since(t).Microseconds())/1000)
				} else {
					ops := []bool{false}
					if phase == "delete-recreate" {
						ops = []bool{true, false}
					}
					for _, deleted := range ops {
						ws := time.Now()
						wt, e := b.write(ctx, k, deleted)
						if e != nil {
							l.err = e
							break
						}
						l.writes++
						l.writeLat = append(l.writeLat, float64(time.Since(ws).Microseconds())/1000)
						l.dbLat = append(l.dbLat, wt.DBMS)
						l.ackLat = append(l.ackLat, wt.AckMS)
						l.lockLat = append(l.lockLat, wt.LockMS)
					}
					if l.err != nil {
						break
					}
				}
			}
			if l.err != nil {
				l.errors = 1
				if stringsHas(l.err.Error(), "stale-after-ack") {
					l.wrong = 1
				}
				if stringsHas(l.err.Error(), "impossible") {
					l.impossible = 1
				}
			}
			ch <- l
		}(worker)
	}
	var all localStats
	for i := 0; i < nread+nwrite; i++ {
		l := <-ch
		all.reads += l.reads
		all.writes += l.writes
		all.errors += l.errors
		all.wrong += l.wrong
		all.impossible += l.impossible
		all.readLat = append(all.readLat, l.readLat...)
		all.writeLat = append(all.writeLat, l.writeLat...)
		all.dbLat = append(all.dbLat, l.dbLat...)
		all.ackLat = append(all.ackLat, l.ackLat...)
		all.lockLat = append(all.lockLat, l.lockLat...)
		if l.err != nil {
			all.err = l.err
		}
	}
	elapsed := time.Since(start).Seconds()
	t := trial{Phase: phase, Number: number, ElapsedS: elapsed, Reads: all.reads, Writes: all.writes, ReadOps: float64(all.reads) / elapsed, WriteOps: float64(all.writes) / elapsed, ReadP50MS: percentile(all.readLat, .5), ReadP99MS: percentile(all.readLat, .99), WriteP50MS: percentile(all.writeLat, .5), WriteP99MS: percentile(all.writeLat, .99), DBP50MS: percentile(all.dbLat, .5), AckP50MS: percentile(all.ackLat, .5), AckP99MS: percentile(all.ackLat, .99), LockP99MS: percentile(all.lockLat, .99), SnapshotQueries: b.stats.snapshots.Load() - s0, ValidationQueries: b.stats.validations.Load() - v0, RelayQueries: b.stats.relayQueries.Load() - r0, DeliveredEvents: b.stats.deliveries.Load() - d0}
	if t.Reads > 0 {
		t.SourceQueriesPerRead = float64(t.SnapshotQueries+t.ValidationQueries) / float64(t.Reads)
	}
	if all.err != nil {
		t.Errors = all.errors
		t.WrongReads = all.wrong
		t.ImpossibleReads = all.impossible
		return t, all.err
	}
	if e := b.audit(ctx); e != nil {
		t.Errors = 1
		if stringsHas(e.Error(), "stale-after-ack") {
			t.WrongReads = 1
		}
		if stringsHas(e.Error(), "impossible") {
			t.ImpossibleReads = 1
		}
		return t, e
	}
	return t, nil
}
