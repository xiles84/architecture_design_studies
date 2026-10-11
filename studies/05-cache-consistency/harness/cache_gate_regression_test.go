package main

import (
	"context"
	"sync"
	"testing"
)

// Concurrent writers of diagnostic samples must not race or lose counts.
func TestReadLogConcurrentCacheGateAccounting(t *testing.T) {
	const workers = 24
	const iterations = 200

	l := newReadLog()
	var wg sync.WaitGroup

	for worker := 0; worker < workers; worker++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			for i := 0; i < iterations; i++ {
				kind := KindImpossible
				if i%2 == 0 {
					kind = KindStale
				}

				l.Record(ReadOutcome{
					KeyID:    int64(id + 1),
					Kind:     kind,
					Source:   SrcHit,
					CacheHit: true,
					StaleMS:  10,
				})
			}
		}(worker)
	}

	wg.Wait()

	s := l.Summary()
	expected := int64(workers * iterations / 2)

	if s.TotalReads != 2*expected {
		t.Fatalf("total=%d, expected=%d", s.TotalReads, 2*expected)
	}
	if s.WrongReads != expected {
		t.Fatalf("stale=%d, expected=%d", s.WrongReads, expected)
	}
	if s.ImpossibleValues != expected {
		t.Fatalf("impossible=%d, expected=%d",
			s.ImpossibleValues, expected)
	}
	if len(s.StaleSamples) != staleSampleLimit {
		t.Fatalf("stale samples=%d, expected=%d",
			len(s.StaleSamples), staleSampleLimit)
	}
}

// Multiple diagnostic snapshots must not increase the cell-wide counter.
func TestPhaseSnapshotsDoNotDoubleCountImpossible(t *testing.T) {
	l := newReadLog()
	a := &adapter{log: l}
	c := &cell{
		ad:  a,
		res: &CellResult{},
	}

	l.Record(ReadOutcome{
		KeyID: 1, Kind: KindImpossible,
		Source: SrcHit, CacheHit: true,
	})

	c.phaseWrong("cold_fill")
	c.phaseWrong("warm_read_only")

	if len(c.res.Wrong) != 2 {
		t.Fatalf("expected two diagnostic snapshots")
	}
	if c.res.ImpossibleValues != 0 {
		t.Fatalf("phase summaries must not accumulate cell total")
	}

	a.impossible.Store(1)

	if got := a.impossible.Load(); got != 1 {
		t.Fatalf("atomic impossible count=%d, expected=1", got)
	}
}

// A fabricated cache HIT must be classified as impossible and
// increment the adapter's authoritative counter.
func TestImpossibleCacheHitIsCounted(t *testing.T) {
	ds := BuildDataset(3, "tiny")
	orc := newOracle(ds)
	key := ds.People[0].ID

	requiredHash, requiredSeq := orc.Required(key)

	a := &adapter{
		orc: orc,
		log: newReadLog(),
	}

	fabricated := PortalContent{
		PersonID: key,
		FullName: "never-committed-cache-state",
	}

	out := a.finish(
		context.Background(),
		key,
		fabricated,
		ReadOutcome{Source: SrcHit, CacheHit: true},
		requiredHash,
		requiredSeq,
	)

	if out.Kind != KindImpossible {
		t.Fatalf("classification=%q, expected impossible", out.Kind)
	}
	if a.impossible.Load() != 1 {
		t.Fatalf("adapter lost impossible cache-hit count")
	}
	if a.log.Summary().ImpossibleValues != 1 {
		t.Fatalf("diagnostic log lost impossible read")
	}
}
