package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// The Redis backend, and the minimal RESP2 client it speaks.
//
// The client is written here rather than pulled in as a dependency, for the same
// reason the study's SQL is embedded: what the cache does must be readable beside
// the result it produced. RESP2 is small and the commands this study needs are a
// handful -- GET, SET with NX PX, DEL, EVAL and INFO -- and every one of them is
// exercised by a dev check against a real Redis before any measurement runs.
//
// The store implements the SAME EntryStore interface as the memory backend, so the
// lease, fencing and expiry semantics above this line are identical and the
// memory-vs-Redis comparison is about the backend, not about two policies.

// ---------------------------------------------------------------- RESP2 client

type respConn struct {
	c net.Conn
	r *bufio.Reader
	w *bufio.Writer
}

func dialRESP(addr string, timeout time.Duration) (*respConn, error) {
	c, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		return nil, err
	}
	_ = c.SetDeadline(time.Now().Add(timeout))
	return &respConn{c: c, r: bufio.NewReader(c), w: bufio.NewWriter(c)}, nil
}

func (rc *respConn) Close() error { return rc.c.Close() }

// do sends one command as a RESP array of bulk strings and returns the decoded
// reply: string for simple strings, int64 for integers, []byte for bulk strings,
// nil for a null bulk/array, and []any for arrays. A RESP error reply becomes a Go
// error carrying the server's message.
func (rc *respConn) do(args ...string) (any, error) {
	_ = rc.c.SetDeadline(time.Now().Add(30 * time.Second))
	if _, err := fmt.Fprintf(rc.w, "*%d\r\n", len(args)); err != nil {
		return nil, err
	}
	for _, a := range args {
		if _, err := fmt.Fprintf(rc.w, "$%d\r\n%s\r\n", len(a), a); err != nil {
			return nil, err
		}
	}
	if err := rc.w.Flush(); err != nil {
		return nil, err
	}
	return readReply(rc.r)
}

func readReply(r *bufio.Reader) (any, error) {
	line, err := r.ReadString('\n')
	if err != nil {
		return nil, err
	}
	line = strings.TrimRight(line, "\r\n")
	if line == "" {
		return nil, fmt.Errorf("redis: empty reply line")
	}
	switch line[0] {
	case '+':
		return line[1:], nil
	case '-':
		return nil, fmt.Errorf("redis: %s", line[1:])
	case ':':
		n, err := strconv.ParseInt(line[1:], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("redis: bad integer %q", line)
		}
		return n, nil
	case '$':
		n, err := strconv.Atoi(line[1:])
		if err != nil {
			return nil, fmt.Errorf("redis: bad bulk length %q", line)
		}
		if n < 0 {
			return nil, nil
		}
		buf := make([]byte, n+2)
		if _, err := io.ReadFull(r, buf); err != nil {
			return nil, err
		}
		return buf[:n], nil
	case '*':
		n, err := strconv.Atoi(line[1:])
		if err != nil {
			return nil, fmt.Errorf("redis: bad array length %q", line)
		}
		if n < 0 {
			return nil, nil
		}
		out := make([]any, 0, n)
		for i := 0; i < n; i++ {
			v, err := readReply(r)
			if err != nil {
				return nil, err
			}
			out = append(out, v)
		}
		return out, nil
	}
	return nil, fmt.Errorf("redis: unknown reply type %q", line)
}

// ---------------------------------------------------------------- the store

const (
	// luaPublish refuses a publication whose fence has moved (KEYS[2] is the key's
	// fence) and one whose publication order is older than the stored entry's.
	luaPublish = `local f = redis.call('GET', KEYS[2])
if (tonumber(f) or 0) ~= tonumber(ARGV[3]) then
  return 0
end
local cur = redis.call('GET', KEYS[1])
if cur then
  local ok, c = pcall(cjson.decode, cur)
  if ok and type(c) == 'table' and c['q'] and tonumber(c['q']) > tonumber(ARGV[2]) then
    return 0
  end
end
redis.call('SET', KEYS[1], ARGV[1], 'PX', ARGV[4])
return 1`

	luaRelease = `if redis.call('GET', KEYS[1]) == ARGV[1] then
  return redis.call('DEL', KEYS[1])
else
  return 0
end`
)

type redisStore struct {
	addr    string
	timeout time.Duration
	pool    chan *respConn
	sem     chan struct{}

	hits        atomic.Int64
	misses      atomic.Int64
	puts        atomic.Int64
	putFenced   atomic.Int64
	deletes     atomic.Int64
	errors      atomic.Int64
	leaseOK     atomic.Int64
	leaseBusy   atomic.Int64
	leaseLost   atomic.Int64
	flushes     atomic.Int64
	corruptVals atomic.Int64
	closed      atomic.Bool

	mu       sync.Mutex
	lastInfo map[string]string
}

func newRedisStore(addr string, maxConns int) *redisStore {
	if maxConns < 2 {
		maxConns = 2
	}
	return &redisStore{
		addr:    addr,
		timeout: 10 * time.Second,
		pool:    make(chan *respConn, maxConns),
		sem:     make(chan struct{}, maxConns),
	}
}

func (r *redisStore) Kind() string { return string(BackendRedis) }

func (r *redisStore) withConn(fn func(*respConn) error) error {
	if r.closed.Load() {
		return fmt.Errorf("redis: store closed")
	}
	// Try an idle connection first; otherwise take a slot and dial a new one.
	select {
	case rc := <-r.pool:
		if err := fn(rc); err != nil {
			rc.Close()
			r.errors.Add(1)
			return err
		}
		select {
		case r.pool <- rc:
		default:
			rc.Close()
		}
		return nil
	default:
	}
	select {
	case r.sem <- struct{}{}:
	default:
		// All slots busy: wait for one. A bounded wait, so a wedged server becomes
		// an error rather than a hang.
		select {
		case r.sem <- struct{}{}:
		case <-time.After(10 * time.Second):
			return fmt.Errorf("redis: no connection slot within 10s")
		}
	}
	defer func() { <-r.sem }()
	rc, err := dialRESP(r.addr, r.timeout)
	if err != nil {
		r.errors.Add(1)
		return err
	}
	if err := fn(rc); err != nil {
		rc.Close()
		r.errors.Add(1)
		return err
	}
	select {
	case r.pool <- rc:
	default:
		rc.Close()
	}
	return nil
}

func (r *redisStore) Get(_ context.Context, key string) (*Entry, bool, error) {
	var e *Entry
	err := r.withConn(func(rc *respConn) error {
		reply, err := rc.do("GET", key)
		if err != nil {
			return err
		}
		b, _ := reply.([]byte)
		if b == nil {
			return nil
		}
		var v Entry
		if err := json.Unmarshal(b, &v); err != nil {
			// A value that cannot be decoded corresponds to no committed state.
			// It is dropped and counted rather than served.
			r.corruptVals.Add(1)
			_, _ = rc.do("DEL", key)
			return nil
		}
		e = &v
		return nil
	})
	if err != nil {
		r.misses.Add(1)
		return nil, false, err
	}
	if e == nil {
		r.misses.Add(1)
		return nil, false, nil
	}
	if e.HardExpired(nowMS()) {
		// Redis's own PX should have removed it; this is the belt-and-braces check
		// so that memory and Redis cannot disagree about hard expiry.
		_ = r.Delete(context.Background(), key)
		r.misses.Add(1)
		return nil, false, nil
	}
	r.hits.Add(1)
	return e, true, nil
}

func (r *redisStore) Put(_ context.Context, key string, e *Entry, ttl time.Duration) (bool, error) {
	b, err := json.Marshal(e)
	if err != nil {
		return false, err
	}
	// The stored PX is the HARD expiry only. Probabilistic early expiration is the
	// adapter's decision, so it must not be expressed by a shorter Redis TTL: a
	// shorter TTL would silently change the semantics being compared.
	px := ttl.Milliseconds()
	if px < 1 {
		px = 1
	}
	var ok bool
	err = r.withConn(func(rc *respConn) error {
		reply, err := rc.do("EVAL", luaPublish, "2", key, fenceKey(key), string(b),
			strconv.FormatInt(e.Seq, 10), strconv.FormatInt(e.Fence, 10), strconv.FormatInt(px, 10))
		if err != nil {
			return err
		}
		n, _ := reply.(int64)
		ok = n == 1
		return nil
	})
	if err != nil {
		return false, err
	}
	r.puts.Add(1)
	if !ok {
		r.putFenced.Add(1)
	}
	return ok, nil
}

func (r *redisStore) Delete(_ context.Context, key string) error {
	r.deletes.Add(1)
	return r.withConn(func(rc *respConn) error {
		_, err := rc.do("DEL", key)
		return err
	})
}

func (r *redisStore) Flush(_ context.Context) error {
	r.flushes.Add(1)
	return r.withConn(func(rc *respConn) error {
		_, err := rc.do("FLUSHDB")
		return err
	})
}

func (r *redisStore) TryLease(_ context.Context, key, token string, ttl time.Duration) (bool, error) {
	var ok bool
	err := r.withConn(func(rc *respConn) error {
		reply, err := rc.do("SET", leaseKey(key), token, "NX", "PX", strconv.FormatInt(ttl.Milliseconds(), 10))
		if err != nil {
			return err
		}
		// "+OK" on success, nil when the key existed.
		ok = reply != nil
		return nil
	})
	if err != nil {
		return false, err
	}
	if ok {
		r.leaseOK.Add(1)
	} else {
		r.leaseBusy.Add(1)
	}
	return ok, nil
}

func (r *redisStore) ReleaseLease(_ context.Context, key, token string) error {
	return r.withConn(func(rc *respConn) error {
		_, err := rc.do("EVAL", luaRelease, "1", leaseKey(key), token)
		return err
	})
}

// leaseKey and fenceKey namespace the lease and the fence away from the value
// keyspace, so a value can never collide with either.
func leaseKey(key string) string { return key + ":fill-lease" }
func fenceKey(key string) string { return key + ":fence" }

// Fence advances the key's invalidation fence and removes any value, in the store
// itself so that every instance of a multi-instance deployment shares it. This is
// what makes a shared cache's strict freshness supportable across instances: the
// tombstone is not just a deletion, it is a statement that anything read before it
// is no longer publishable.
func (r *redisStore) Fence(_ context.Context, key string) (int64, error) {
	var out int64
	err := r.withConn(func(rc *respConn) error {
		reply, err := rc.do("INCR", fenceKey(key))
		if err != nil {
			return err
		}
		out, _ = reply.(int64)
		if _, err := rc.do("DEL", key); err != nil {
			return err
		}
		return nil
	})
	return out, err
}

func (r *redisStore) FenceOf(_ context.Context, key string) (int64, error) {
	var out int64
	err := r.withConn(func(rc *respConn) error {
		reply, err := rc.do("GET", fenceKey(key))
		if err != nil {
			return err
		}
		if b, ok := reply.([]byte); ok && b != nil {
			n, perr := strconv.ParseInt(string(b), 10, 64)
			if perr == nil {
				out = n
			}
		}
		return nil
	})
	return out, err
}

func (r *redisStore) Close(_ context.Context) error {
	r.closed.Store(true)
	for {
		select {
		case rc := <-r.pool:
			rc.Close()
		default:
			return nil
		}
	}
}

// Ping proves the cache under test answers. A scenario that needs a cache must fail
// loudly when there is none: the first dev check of this study ran every redis
// scenario against a cache that was never started, and because a cache outage
// degrades to authoritative reads, every cell "passed" its phases while measuring
// nothing but the database. A cache that is absent must never look like a cache that
// is merely cold.
func (r *redisStore) Ping(ctx context.Context) error {
	return r.withConn(func(rc *respConn) error {
		v, err := rc.do("PING")
		if err != nil {
			return err
		}
		if s, _ := v.(string); s != "PONG" {
			return fmt.Errorf("redis: unexpected PING reply %v", v)
		}
		return nil
	})
}

// SetAvailable turns the store on or off for the fault-injection phases.
func (r *redisStore) SetAvailable(on bool) { r.closed.Store(!on) }

func (r *redisStore) Stats() map[string]int64 {
	return map[string]int64{
		"hits":                 r.hits.Load(),
		"misses":               r.misses.Load(),
		"puts":                 r.puts.Load(),
		"put_fenced":           r.putFenced.Load(),
		"deletes":              r.deletes.Load(),
		"errors":               r.errors.Load(),
		"lease_acquired":       r.leaseOK.Load(),
		"lease_contended":      r.leaseBusy.Load(),
		"lease_expired_stolen": r.leaseLost.Load(),
		"flushes":              r.flushes.Load(),
		"corrupt_values":       r.corruptVals.Load(),
	}
}

// Info reads Redis's own account of itself. The fields are the ones the study's
// metrics list names: maxmemory and its policy, keyspace hits and misses, expired
// and evicted keys, command counts, memory, CPU and network bytes.
func (r *redisStore) Info(ctx context.Context) map[string]string {
	out := map[string]string{}
	err := r.withConn(func(rc *respConn) error {
		reply, err := rc.do("INFO")
		if err != nil {
			return err
		}
		b, _ := reply.([]byte)
		if b == nil {
			return fmt.Errorf("redis: INFO returned no body")
		}
		for _, line := range strings.Split(string(b), "\n") {
			line = strings.TrimRight(line, "\r")
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			k, v, ok := strings.Cut(line, ":")
			if !ok {
				continue
			}
			switch k {
			case "used_memory", "used_memory_peak", "maxmemory", "maxmemory_policy",
				"maxmemory_samples", "keyspace_hits", "keyspace_misses", "expired_keys",
				"evicted_keys", "total_commands_processed", "total_net_input_bytes",
				"total_net_output_bytes", "used_cpu_sys", "used_cpu_user",
				"instantaneous_ops_per_sec", "mem_fragmentation_ratio", "redis_version":
				out[k] = v
			}
			if strings.HasPrefix(k, "cmdstat_get") || strings.HasPrefix(k, "cmdstat_set") ||
				strings.HasPrefix(k, "cmdstat_eval") || strings.HasPrefix(k, "cmdstat_del") {
				out[k] = v
			}
		}
		return nil
	})
	if err != nil {
		out["error"] = err.Error()
	}
	r.mu.Lock()
	r.lastInfo = out
	r.mu.Unlock()
	_ = ctx
	return out
}
