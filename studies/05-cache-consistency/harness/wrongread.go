package main

import (
	"sort"
	"sync"
)

// Exact wrong-read accounting.
//
// The rule this file exists to respect: a timed cache hit is NEVER validated by a
// hidden synchronous full database read, because that would destroy the performance
// question. Instead every read reports what it returned and what the oracle says it
// required, and the classification below separates the four outcomes the study's
// vocabulary names.

const (
	// KindFresh: the value is the committed state the read required.
	KindFresh = "fresh"
	// KindStale: the value is an EARLIER committed state than the read required.
	// This is the violation the strict contract forbids.
	KindStale = "stale"
	// KindAhead: the value is a LATER committed state than the read required,
	// which happens legitimately when a write commits during the read. Not a
	// violation; counted separately so it is never quietly folded into "fresh".
	KindAhead = "ahead"
	// KindImpossible: the value corresponds to no committed state this key ever
	// had. Always a failure, under every freshness policy.
	KindImpossible = "impossible"
)

// Read sources. Which path served a read is part of the result, because "a cache
// hit" and "a cache that fell back to the database" are different products.
const (
	SrcDatabase  = "database"      // no cache in this scenario
	SrcHit       = "hit"           // served from the cache
	SrcValidated = "hit-validated" // served from the cache after an authoritative version check
	SrcFill      = "fill"          // this reader took the lease and filled
	SrcFallback  = "fallback"      // the lease did not resolve; authoritative read
	SrcBypass    = "bypass"        // strict freshness required an authoritative read
	SrcError     = "error"         // the cache errored
)

// ReadOutcome is one read's account of itself.
type ReadOutcome struct {
	KeyID          int64
	RequiredHash   string
	RequiredSeq    int64
	ReturnedHash   string
	ReturnedSeq    int64
	Behind         int64
	Kind           string
	Source         string
	AgeMS          int64
	Overlapped     bool
	ClaimedVersion int64
	// The next four are diagnostic and are filled only on a stale cache HIT: they say
	// whether the surviving entry was published at the fence that is current now (a
	// publication-ordering bug) or outlived a fence that should have removed it (a
	// deletion or store-identity bug).
	EntryFence   int64
	CurrentFence int64
	EntrySeq     int64
	CurrentSeq   int64
	// Cause names the evidenced cause when this read was wrong. The adapter sets
	// it from the phase it is in and from what failed; the log only groups it.
	Cause    string
	CacheHit bool
	// StaleMS is how old the served committed state was when it was served. It is
	// the staleness duration the protocol asks to report.
	StaleMS float64
	// ExpiredBy is "hard", "probabilistic" or "" for a hit that was refreshed.
	ExpiredBy string
}

// staleSampleLimit bounds the diagnostic sample so a badly wrong cell cannot write a
// huge result file.
const staleSampleLimit = 12

type readLog struct {
	mu sync.Mutex

	total, fresh, wrong, ahead, impossible, overlapped int64
	hits                                               int64
	bySource                                           map[string]int64
	byCause                                            map[string]int64
	expiredHard, expiredProb                           int64
	staleMS                                            []float64
	keysWrong                                          map[int64]int64
	maxBehind                                          int64
	streak                                             map[int64]int64
	maxStreak                                          int64
	pendingAck                                         map[int64]float64
	ttfMS                                              []float64
	staleSamples                                       []ReadOutcome
	leaseWaitMS                                        []float64
	// maxStaleSamples bounds the detail retained. Counts are exact regardless;
	// only the distribution is bounded, and the bound is large enough that no run
	// in this study reaches it.
	maxStaleSamples int
}

func newReadLog() *readLog {
	return &readLog{
		bySource:        map[string]int64{},
		byCause:         map[string]int64{},
		keysWrong:       map[int64]int64{},
		streak:          map[int64]int64{},
		pendingAck:      map[int64]float64{},
		maxStaleSamples: 200_000,
	}
}

// NoteAck records that a write to this key was acknowledged, starting the
// time-to-freshness clock. The most recent acknowledgement wins, so the reported
// time is "how long until this key was fresh again after the last write".
func (l *readLog) NoteAck(keyID int64) {
	l.mu.Lock()
	l.pendingAck[keyID] = float64(nowMS())
	l.mu.Unlock()
}

func (l *readLog) Record(o ReadOutcome) {
	if o.Kind == KindStale && len(l.staleSamples) < staleSampleLimit {
		l.staleSamples = append(l.staleSamples, o)
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	l.total++
	l.bySource[o.Source]++
	if o.CacheHit {
		l.hits++
	}
	if o.Overlapped {
		l.overlapped++
	}
	switch o.ExpiredBy {
	case "hard":
		l.expiredHard++
	case "probabilistic":
		l.expiredProb++
	}
	switch o.Kind {
	case KindFresh:
		l.fresh++
		l.streak[o.KeyID] = 0
		if t, ok := l.pendingAck[o.KeyID]; ok {
			l.ttfMS = append(l.ttfMS, float64(nowMS())-t)
			delete(l.pendingAck, o.KeyID)
		}
	case KindStale:
		l.wrong++
		cause := o.Cause
		if cause == "" {
			cause = "other evidenced cause"
		}
		l.byCause[cause]++
		if len(l.staleMS) < l.maxStaleSamples {
			l.staleMS = append(l.staleMS, o.StaleMS)
		}
		l.keysWrong[o.KeyID]++
		l.streak[o.KeyID]++
		if l.streak[o.KeyID] > l.maxStreak {
			l.maxStreak = l.streak[o.KeyID]
		}
	case KindAhead:
		l.ahead++
		l.streak[o.KeyID] = 0
	case KindImpossible:
		l.impossible++
		l.keysWrong[o.KeyID]++
	}
	if o.Behind > l.maxBehind {
		l.maxBehind = o.Behind
	}
}

func (l *readLog) NoteLeaseWait(ms float64) {
	l.mu.Lock()
	l.leaseWaitMS = append(l.leaseWaitMS, ms)
	l.mu.Unlock()
}

// WrongReadSummary is what the report prints. Every field the protocol asks for is
// here, including the two that are easy to omit and that change the reading:
// wrong reads as a share of cache HITS (not only of all reads), and the number of
// concurrent/ambiguous reads kept separate from both.
type WrongReadSummary struct {
	TotalReads           int64            `json:"total_reads"`
	FreshReads           int64            `json:"fresh_reads"`
	WrongReads           int64            `json:"wrong_stale_reads"`
	AheadReads           int64            `json:"ahead_reads"`
	CacheHits            int64            `json:"cache_hits"`
	WrongPctOfAllReads   float64          `json:"wrong_pct_of_all_reads"`
	WrongPctOfCacheHits  float64          `json:"wrong_pct_of_cache_hits"`
	UniqueKeysAffected   int              `json:"unique_keys_affected"`
	MaxVersionsBehind    int64            `json:"max_versions_behind"`
	StaleDuration        msStat           `json:"staleness_duration_ms"`
	ConsecutiveWrongMax  int64            `json:"consecutive_wrong_read_streak_max"`
	TimeToFreshness      msStat           `json:"time_to_freshness_after_ack_ms"`
	ConcurrentAmbiguous  int64            `json:"concurrent_ambiguous_reads"`
	ImpossibleValues     int64            `json:"impossible_cache_values"`
	ExpiredHard          int64            `json:"expired_hard"`
	ExpiredProbabilistic int64            `json:"expired_probabilistic"`
	BySource             map[string]int64 `json:"reads_by_source"`
	ByCause              map[string]int64 `json:"wrong_reads_by_cause"`
	// StaleSamples keeps the first few stale reads in full, with the surviving
	// entry's fence and sequence beside the current ones. The count says a rule
	// failed; only these say which rule.
	StaleSamples []ReadOutcome `json:"stale_read_samples,omitempty"`
	LeaseWait    msStat        `json:"lease_wait_ms"`
}

func (l *readLog) Summary() WrongReadSummary {
	l.mu.Lock()
	defer l.mu.Unlock()
	s := WrongReadSummary{
		TotalReads:           l.total,
		FreshReads:           l.fresh,
		WrongReads:           l.wrong,
		AheadReads:           l.ahead,
		CacheHits:            l.hits,
		UniqueKeysAffected:   len(l.keysWrong),
		MaxVersionsBehind:    l.maxBehind,
		ConsecutiveWrongMax:  l.maxStreak,
		ConcurrentAmbiguous:  l.overlapped,
		ImpossibleValues:     l.impossible,
		ExpiredHard:          l.expiredHard,
		ExpiredProbabilistic: l.expiredProb,
		BySource:             copyCounts(l.bySource),
		ByCause:              copyCounts(l.byCause),
		StaleSamples:         append([]ReadOutcome(nil), l.staleSamples...),
		StaleDuration:        statFromSamples(l.staleMS),
		TimeToFreshness:      statFromSamples(l.ttfMS),
		LeaseWait:            statFromSamples(l.leaseWaitMS),
	}
	if l.total > 0 {
		s.WrongPctOfAllReads = 100 * float64(l.wrong) / float64(l.total)
	}
	if l.hits > 0 {
		s.WrongPctOfCacheHits = 100 * float64(l.wrong) / float64(l.hits)
	}
	return s
}

// statFromSamples computes p50/p99/max over a copy. Percentiles are computed from
// every recorded sample rather than from a histogram, and `max` is the honest tail
// figure where the sample count does not support a deeper percentile
// (methodology 7).
func statFromSamples(xs []float64) msStat {
	if len(xs) == 0 {
		return msStat{}
	}
	c := append([]float64(nil), xs...)
	sort.Float64s(c)
	at := func(p float64) float64 { return c[int(p*float64(len(c)-1))] }
	return msStat{Count: int64(len(c)), P50: at(0.50), P99: at(0.99), Max: c[len(c)-1]}
}

// copyCounts returns a snapshot. Returning the live map would make every phase's
// summary share one growing map, and the report would show the cell's final totals
// on every line while the counters beside them differed -- a contradiction no reader
// could resolve.
func copyCounts(m map[string]int64) map[string]int64 {
	out := make(map[string]int64, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}
