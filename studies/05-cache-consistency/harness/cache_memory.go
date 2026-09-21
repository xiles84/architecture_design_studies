package main

import (
	"container/list"
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"
)

// The in-process memory backend: a byte-bounded, concurrency-safe EXACT LRU.
//
// Exact, not approximate: an eviction here is the least-recently-used entry, with
// no sampling and no approximation error. Redis's allkeys-lru is approximate, and
// the study says so wherever the two are compared -- an eviction difference between
// the backends is partly a policy difference, and hiding that would turn an
// implementation detail into a finding.
//
// Capacity is charged in bytes rather than in keys. A key count would let a cache
// look bounded while the payloads grew past the container's memory limit, and the
// study's whole resource accounting (add-cache vs equal-total) depends on the
// cache's resident bytes being a number, not an assumption.

type memItem struct {
	key   string
	entry *Entry
	size  int64
}

type memLease struct {
	token     string
	expiresMS int64
}

type memStore struct {
	mu       sync.Mutex
	capBytes int64
	used     int64
	items    map[string]*list.Element
	lru      *list.List
	leases   map[string]memLease
	now      func() int64

	evictions   int64
	puts        int64
	putFenced   int64
	putTooLarge int64
	gets        int64
	hits        int64
	misses      int64
	hardExpired int64
	deletes     int64
	leaseOK     int64
	leaseBusy   int64
	leaseLost   int64
	flushCount  int64
	closed      bool
}

func newMemStore(capBytes int64) *memStore {
	return &memStore{
		capBytes: capBytes,
		items:    map[string]*list.Element{},
		lru:      list.New(),
		leases:   map[string]memLease{},
		now:      func() int64 { return time.Now().UnixNano() / int64(time.Millisecond) },
	}
}

func (m *memStore) Kind() string { return string(BackendMemory) }

func (m *memStore) Get(_ context.Context, key string) (*Entry, bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.gets++
	if m.closed {
		m.misses++
		return nil, false, fmt.Errorf("memory cache closed")
	}
	el, ok := m.items[key]
	if !ok {
		m.misses++
		return nil, false, nil
	}
	it := el.Value.(*memItem)
	if it.entry.HardExpired(m.now()) {
		// Hard expiry is enforced by the store as well as by the adapter: a value
		// past its deadline must never be served, whatever the caller believes.
		m.removeElement(el)
		m.hardExpired++
		m.misses++
		return nil, false, nil
	}
	m.lru.MoveToFront(el)
	m.hits++
	return it.entry, true, nil
}

func (m *memStore) Put(_ context.Context, key string, e *Entry, _ time.Duration) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.puts++
	if m.closed {
		return false, fmt.Errorf("memory cache closed")
	}
	if el, ok := m.items[key]; ok {
		cur := el.Value.(*memItem)
		// Version fencing: a newer committed version must not be replaced by an
		// older one, which is what an expired lease holder would otherwise do.
		if cur.entry.Seq > e.Seq {
			m.putFenced++
			return false, nil
		}
		m.removeElement(el)
	}
	size := e.Size(key)
	if size > m.capBytes {
		// The entry alone does not fit. Publishing it would break the byte bound
		// silently; refusing it is the honest behaviour and is counted.
		m.putTooLarge++
		return false, nil
	}
	el := m.lru.PushFront(&memItem{key: key, entry: e, size: size})
	m.items[key] = el
	m.used += size
	for m.used > m.capBytes && m.lru.Len() > 0 {
		back := m.lru.Back()
		if back == nil {
			break
		}
		m.removeElement(back)
		m.evictions++
	}
	return true, nil
}

func (m *memStore) Delete(_ context.Context, key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.deletes++
	if el, ok := m.items[key]; ok {
		m.removeElement(el)
	}
	return nil
}

func (m *memStore) Flush(_ context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.items = map[string]*list.Element{}
	m.lru.Init()
	m.used = 0
	m.leases = map[string]memLease{}
	m.flushCount++
	return nil
}

func (m *memStore) TryLease(_ context.Context, key, token string, ttl time.Duration) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.closed {
		return false, fmt.Errorf("memory cache closed")
	}
	now := m.now()
	if l, ok := m.leases[key]; ok {
		if l.expiresMS > now {
			m.leaseBusy++
			return false, nil
		}
		// An expired lease is stealable. Counting the steal is how a lease-holder
		// death becomes a measured event rather than an invisible one.
		m.leaseLost++
	}
	m.leases[key] = memLease{token: token, expiresMS: now + ttl.Milliseconds()}
	m.leaseOK++
	m.sweepLeasesLocked(now)
	return true, nil
}

func (m *memStore) ReleaseLease(_ context.Context, key, token string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if l, ok := m.leases[key]; ok && l.token == token {
		delete(m.leases, key)
	}
	return nil
}

// sweepLeasesLocked bounds the lease map. Without it a long churn run accumulates
// one dead lease entry per distinct key forever, which is a leak dressed up as
// bookkeeping.
func (m *memStore) sweepLeasesLocked(now int64) {
	const sweepAt = 4096
	if len(m.leases) < sweepAt {
		return
	}
	for k, l := range m.leases {
		if l.expiresMS <= now {
			delete(m.leases, k)
		}
	}
}

func (m *memStore) removeElement(el *list.Element) {
	it := el.Value.(*memItem)
	m.lru.Remove(el)
	delete(m.items, it.key)
	m.used -= it.size
}

func (m *memStore) Close(_ context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.closed = true
	return nil
}

// FailNext makes the store unavailable for the fault-injection phases, which is how
// a cache outage is reproduced deterministically instead of by stopping a container
// and hoping the timing is right.
func (m *memStore) FailNext(on bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.closed = on
}

// SetAvailable turns the store on or off for the fault-injection phases. It is how a
// cache outage is reproduced deterministically instead of by stopping a process and
// hoping the timing is right.
func (m *memStore) SetAvailable(on bool) {
	m.mu.Lock()
	m.closed = !on
	m.mu.Unlock()
}

func (m *memStore) Stats() map[string]int64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	return map[string]int64{
		"resident_bytes":       m.used,
		"capacity_bytes":       m.capBytes,
		"items":                int64(m.lru.Len()),
		"evictions":            m.evictions,
		"puts":                 m.puts,
		"put_fenced":           m.putFenced,
		"put_too_large":        m.putTooLarge,
		"gets":                 m.gets,
		"hits":                 m.hits,
		"misses":               m.misses,
		"hard_expired":         m.hardExpired,
		"deletes":              m.deletes,
		"lease_acquired":       m.leaseOK,
		"lease_contended":      m.leaseBusy,
		"lease_expired_stolen": m.leaseLost,
		"flushes":              m.flushCount,
		"active_leases":        int64(len(m.leases)),
	}
}

func (m *memStore) Info(_ context.Context) map[string]string {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	m.mu.Lock()
	defer m.mu.Unlock()
	return map[string]string{
		"backend":           "memory",
		"lru_policy":        "exact (container/list, byte-bounded)",
		"capacity_bytes":    fmt.Sprint(m.capBytes),
		"resident_bytes":    fmt.Sprint(m.used),
		"go_heap_alloc":     fmt.Sprint(ms.HeapAlloc),
		"go_total_alloc":    fmt.Sprint(ms.TotalAlloc),
		"go_num_gc":         fmt.Sprint(ms.NumGC),
		"go_pause_total_ns": fmt.Sprint(ms.PauseTotalNs),
	}
}
