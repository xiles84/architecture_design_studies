package main

import (
	"bufio"
	"bytes"
	"context"
	"math/rand"
	"strings"
	"testing"
	"time"
)

// The unit tests decide whether a number is a measurement or a confident mistake.
// Each one pins a rule the study's conclusions rest on.

// ---------------------------------------------------------------- canonical hash

func TestCanonicalHashIgnoresSliceOrderButNotMembership(t *testing.T) {
	base := PortalContent{PersonID: 1, FullName: "a", DonationCount: 2, Recent: []PortalDonation{
		{ID: 2, AmountCents: 200, Currency: "USD", DonatedAt: "2026-01-02T00:00:00Z"},
		{ID: 1, AmountCents: 100, Currency: "USD", DonatedAt: "2026-01-01T00:00:00Z"},
	}}
	reordered := base
	reordered.Recent = []PortalDonation{base.Recent[1], base.Recent[0]}
	if base.ContentHash() != reordered.ContentHash() {
		t.Fatal("canonical hash depends on slice order; a design that returns the same donations in a different order would be called wrong")
	}
	changed := base
	changed.Recent = []PortalDonation{base.Recent[1]}
	if base.ContentHash() == changed.ContentHash() {
		t.Fatal("canonical hash ignores slice membership")
	}
	var zero PortalContent
	if zero.ContentHash() == "" {
		t.Fatal("empty content must still hash")
	}
}

// ---------------------------------------------------------------- LRU

func TestMemoryLRUEvictsLeastRecentlyUsedAndRespectsByteBound(t *testing.T) {
	ctx := context.Background()
	now := nowMS()
	mk := func(id int64, v int64) *Entry {
		c := PortalContent{PersonID: id, DonationTotalCents: v}
		return newEntry(c, 0, id, now)
	}
	// Measure one entry's real size first, then size the store to fit exactly two:
	// a capacity chosen by guessing would make this a test of the guess.
	probe := newMemStore(1 << 20)
	if ok, _ := probe.Put(ctx, "a", mk(1, 1), hardTTL); !ok {
		t.Fatal("probe put a failed")
	}
	one := probe.Stats()["resident_bytes"]
	if one <= 0 {
		t.Fatal("the store reported no resident bytes for a stored entry")
	}
	s := newMemStore(one * 2)
	if ok, _ := s.Put(ctx, "a", mk(1, 1), hardTTL); !ok {
		t.Fatal("put a failed")
	}
	if ok, _ := s.Put(ctx, "b", mk(2, 2), hardTTL); !ok {
		t.Fatal("put b failed")
	}
	// Touch a so b becomes the least recently used.
	if _, hit, _ := s.Get(ctx, "a"); !hit {
		t.Fatal("a should be present")
	}
	if ok, _ := s.Put(ctx, "c", mk(3, 3), hardTTL); !ok {
		t.Fatal("put c failed")
	}
	stats := s.Stats()
	if stats["resident_bytes"] > one*2 {
		t.Fatalf("byte bound exceeded: %d > %d", stats["resident_bytes"], one*2)
	}
	if stats["evictions"] == 0 {
		t.Fatal("no eviction occurred although the bound was exceeded")
	}
	if _, hit, _ := s.Get(ctx, "b"); hit {
		t.Fatal("b should have been evicted as least recently used")
	}
	if _, hit, _ := s.Get(ctx, "a"); !hit {
		t.Fatal("a should have survived: it was used more recently than b")
	}
}

func TestMemoryStoreRefusesAnEntryLargerThanCapacity(t *testing.T) {
	s := newMemStore(64)
	ctx := context.Background()
	e := newEntry(PortalContent{PersonID: 1, FullName: strings.Repeat("x", 500)}, 0, 1, nowMS())
	ok, err := s.Put(ctx, "big", e, 0)
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Fatal("an entry larger than the capacity was accepted, which breaks the byte bound silently")
	}
	if s.Stats()["put_too_large"] != 1 {
		t.Fatal("the refusal was not counted")
	}
}

// ---------------------------------------------------------------- expiry

func TestHardExpiryIsExact(t *testing.T) {
	s := newMemStore(1 << 20)
	now := int64(1_000_000)
	s.now = func() int64 { return now }
	ctx := context.Background()

	e := newEntry(PortalContent{PersonID: 1}, 0, 1, now)
	if _, hit, _ := s.Get(ctx, "k"); hit {
		t.Fatal("empty store returned a hit")
	}
	if _, err := s.Put(ctx, "k", e, hardTTL); err != nil {
		t.Fatal(err)
	}
	if e.ExpiresMS-now != 300_000 {
		t.Fatalf("hard expiry is %d ms after creation, expected exactly 300000", e.ExpiresMS-now)
	}
	now += 299_999
	if _, hit, _ := s.Get(ctx, "k"); !hit {
		t.Fatal("entry expired before its deadline")
	}
	now += 1
	if _, hit, _ := s.Get(ctx, "k"); hit {
		t.Fatal("entry was served AT its hard-expiry deadline; expiry must be exact")
	}
	if s.Stats()["hard_expired"] != 1 {
		t.Fatal("the hard expiry was not counted")
	}
}

// TestProbabilisticEarlyExpiry pins the protocol's formula:
//
//	p = clamp(1 - remaining/300s, 0, 1)
//
// at the five ages the protocol names, with a deterministic draw.
func TestProbabilisticEarlyExpiry(t *testing.T) {
	cases := []struct {
		ageS     float64
		wantProb float64
	}{
		{0, 0.0},
		{75, 0.25},
		{150, 0.5},
		{225, 0.75},
		{300, 1.0},
	}
	for _, c := range cases {
		remaining := hardTTLSeconds - c.ageS
		// A draw just below p must fire; a draw just above must not. This pins the
		// threshold itself, not just its direction.
		below := c.wantProb - 0.001
		above := c.wantProb + 0.001
		if below < 0 {
			below = 0
		}
		if below < c.wantProb && !ShouldExpireEarly(remaining, hardTTLSeconds, below) {
			t.Fatalf("age %vs: a draw of %.4f should fire (p=%.2f)", c.ageS, below, c.wantProb)
		}
		if c.wantProb < 1 && ShouldExpireEarly(remaining, hardTTLSeconds, above) {
			t.Fatalf("age %vs: a draw of %.4f should not fire (p=%.2f)", c.ageS, above, c.wantProb)
		}
	}
	// The deterministic draw from a recorded source must reproduce exactly, so a
	// run can be replayed from its seed.
	r1 := rand.New(rand.NewSource(7))
	r2 := rand.New(rand.NewSource(7))
	for i := 0; i < 100; i++ {
		a := ShouldExpireEarly(150, hardTTLSeconds, r1.Float64())
		b := ShouldExpireEarly(150, hardTTLSeconds, r2.Float64())
		if a != b {
			t.Fatal("the probabilistic expiry draw is not reproducible from its seed")
		}
	}
	// At the start of an entry's life nothing should ever fire, and past the
	// deadline everything must.
	if ShouldExpireEarly(hardTTLSeconds, hardTTLSeconds, 0.0) {
		t.Fatal("a brand-new entry fired the early-expiry draw")
	}
	if !ShouldExpireEarly(0, hardTTLSeconds, 0.999999) {
		t.Fatal("an entry past its deadline did not fire")
	}
}

func TestAgeBucketsCoverTheFivePinnedAges(t *testing.T) {
	cases := map[int64]string{
		0: "0-75s", 75_000: "75-150s", 150_000: "150-225s", 225_000: "225-300s", 300_000: ">=300s",
	}
	for age, want := range cases {
		if got := ageBucketNames[bucketIndex(age)]; got != want {
			t.Fatalf("age %d ms fell into bucket %q, expected %q", age, got, want)
		}
	}
}

// ---------------------------------------------------------------- leases

func TestMemoryLeaseIsExclusiveAndStealableAfterExpiry(t *testing.T) {
	s := newMemStore(1 << 20)
	now := int64(0)
	s.now = func() int64 { return now }
	ctx := context.Background()

	ok, err := s.TryLease(ctx, "k", "t1", 100*time.Millisecond)
	if err != nil || !ok {
		t.Fatalf("first lease failed: ok=%v err=%v", ok, err)
	}
	if ok2, _ := s.TryLease(ctx, "k", "t2", 100*time.Millisecond); ok2 {
		t.Fatal("a second holder acquired a lease that was already held")
	}
	// A release with the wrong token must not remove the holder's lease.
	_ = s.ReleaseLease(ctx, "k", "t2")
	if ok3, _ := s.TryLease(ctx, "k", "t3", 100*time.Millisecond); ok3 {
		t.Fatal("a wrong-token release freed the lease")
	}
	// After expiry, a new holder must be able to take it.
	now += 101
	if ok4, _ := s.TryLease(ctx, "k", "t4", 100*time.Millisecond); !ok4 {
		t.Fatal("an expired lease was not stealable: a dead holder would wedge the key forever")
	}
	if s.Stats()["lease_expired_stolen"] == 0 {
		t.Fatal("the stolen lease was not counted")
	}
	// The correct token releases it.
	_ = s.ReleaseLease(ctx, "k", "t4")
	if ok5, _ := s.TryLease(ctx, "k", "t5", 100*time.Millisecond); !ok5 {
		t.Fatal("a compare-and-delete release did not free the lease")
	}
}

func TestVersionFencingRefusesAnOlderPublication(t *testing.T) {
	s := newMemStore(1 << 20)
	ctx := context.Background()
	now := nowMS()
	newer := newEntry(PortalContent{PersonID: 1, FullName: "new"}, 5, 5, now)
	if ok, _ := s.Put(ctx, "k", newer, 0); !ok {
		t.Fatal("newer publication failed")
	}
	older := newEntry(PortalContent{PersonID: 1, FullName: "old"}, 3, 3, now)
	ok, err := s.Put(ctx, "k", older, 0)
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Fatal("an older publication overwrote a newer committed version")
	}
	if s.Stats()["put_fenced"] != 1 {
		t.Fatal("the fenced publication was not counted")
	}
	if got, hit, _ := s.Get(ctx, "k"); !hit || got.Hash != newer.Hash {
		t.Fatal("the stored value changed after a fenced publication")
	}
}

// ---------------------------------------------------------------- oracle

func TestOracleClassifiesFreshStaleAheadAndImpossible(t *testing.T) {
	ds := BuildDataset(3, "tiny")
	orc := newOracle(ds)
	p := ds.People[0]
	key := p.ID

	hash0, seq0 := orc.Required(key)
	if seq0 != 0 {
		t.Fatalf("a freshly loaded key starts at sequence %d, expected 0", seq0)
	}
	if k, _, _, _ := orc.Classify(key, hash0, hash0, seq0); k != KindFresh {
		t.Fatalf("the loaded state was classified %q, expected fresh", k)
	}

	// One committed insert: a new state, sequence 1.
	d, _ := orc.PickDonation(key, rand.New(rand.NewSource(1)))
	d.ID = ds.MaxDonationID + 1
	d.PersonID = key
	d.CharityID = p.CharityID
	orc.ApplyInsert(key, d)
	hash1, seq1 := orc.Required(key)
	if seq1 != 1 {
		t.Fatalf("after one acknowledged write the sequence is %d, expected 1", seq1)
	}

	// A read that requires hash1 and gets hash0 is STALE, one version behind.
	if k, _, behind, _ := orc.Classify(key, hash0, hash1, seq1); k != KindStale || behind != 1 {
		t.Fatalf("a read of the older state was classified %q behind %d, expected stale/1", k, behind)
	}
	// A read that required hash0 and got hash1 is AHEAD, not wrong.
	if k, _, _, _ := orc.Classify(key, hash1, hash0, seq0); k != KindAhead {
		t.Fatalf("a read that saw a later committed state was classified %q, expected ahead", k)
	}
	// A hash the key never had is IMPOSSIBLE.
	if k, _, _, _ := orc.Classify(key, "deadbeef", hash1, seq1); k != KindImpossible {
		t.Fatalf("an uncommitted payload was classified %q, expected impossible", k)
	}
	// The oracle's view and its expected content agree.
	want, ok := orc.ExpectedContent(key)
	if !ok || want.ContentHash() != hash1 {
		t.Fatal("ExpectedContent does not reproduce the oracle's current hash")
	}
}

func TestOracleReassignmentMovesBothKeys(t *testing.T) {
	ds := BuildDataset(4, "tiny")
	orc := newOracle(ds)
	a, b := ds.People[0], ds.People[1]
	d, ok := orc.PickDonation(a.ID, rand.New(rand.NewSource(2)))
	if !ok {
		t.Fatal("no donation to move")
	}
	beforeA, seqA := orc.Required(a.ID)
	beforeB, seqB := orc.Required(b.ID)
	d.CharityID = b.CharityID
	orc.ApplyReassign(d.ID, a.ID, b.ID, b.CharityID)

	afterA, seqA2 := orc.Required(a.ID)
	afterB, seqB2 := orc.Required(b.ID)
	if seqA2 != seqA+1 || seqB2 != seqB+1 {
		t.Fatalf("reassignment bumped %d and %d; both keys must advance", seqA2-seqA, seqB2-seqB)
	}
	if afterA == beforeA || afterB == beforeB {
		t.Fatal("reassignment did not change both ports of the portal view")
	}
	ids := orc.ExpectedDonationIDs(a.ID)
	for _, id := range ids {
		if id == d.ID {
			t.Fatal("the moved donation is still on the source donor")
		}
	}
}

// ---------------------------------------------------------------- wrong-read log

func TestReadLogSeparatesStaleFromConcurrentAndImpossible(t *testing.T) {
	l := newReadLog()
	l.Record(ReadOutcome{KeyID: 1, Kind: KindFresh, Source: SrcHit, CacheHit: true})
	l.Record(ReadOutcome{KeyID: 1, Kind: KindStale, Source: SrcHit, CacheHit: true, Behind: 2, StaleMS: 40, Cause: "invalidation race", Overlapped: true})
	l.Record(ReadOutcome{KeyID: 1, Kind: KindStale, Source: SrcHit, CacheHit: true, Behind: 1, StaleMS: 10, Cause: "invalidation race"})
	l.Record(ReadOutcome{KeyID: 2, Kind: KindAhead, Source: SrcFill})
	l.Record(ReadOutcome{KeyID: 3, Kind: KindImpossible, Source: SrcHit, CacheHit: true})
	s := l.Summary()
	if s.TotalReads != 5 || s.FreshReads != 1 || s.WrongReads != 2 || s.AheadReads != 1 || s.ImpossibleValues != 1 {
		t.Fatalf("unexpected counts: %+v", s)
	}
	if s.ConcurrentAmbiguous != 1 {
		t.Fatalf("concurrent/ambiguous reads were not kept separate: %d", s.ConcurrentAmbiguous)
	}
	if s.UniqueKeysAffected != 2 {
		t.Fatalf("unique keys affected = %d, expected 2", s.UniqueKeysAffected)
	}
	if s.MaxVersionsBehind != 2 {
		t.Fatalf("max versions behind = %d, expected 2", s.MaxVersionsBehind)
	}
	if s.ConsecutiveWrongMax != 2 {
		t.Fatalf("longest consecutive wrong-read streak = %d, expected 2", s.ConsecutiveWrongMax)
	}
	if s.WrongPctOfCacheHits <= 0 || s.WrongPctOfAllReads <= 0 {
		t.Fatal("wrong-read percentages were not computed")
	}
	if s.ByCause["invalidation race"] != 2 {
		t.Fatalf("wrong reads were not grouped by cause: %v", s.ByCause)
	}
}

// ---------------------------------------------------------------- RESP parsing

func TestRESPDecoding(t *testing.T) {
	cases := []struct {
		in   string
		want any
	}{
		{"+OK\r\n", "OK"},
		{":42\r\n", int64(42)},
		{"$3\r\nabc\r\n", []byte("abc")},
		{"$-1\r\n", nil},
		{"*2\r\n$1\r\na\r\n:1\r\n", []any{[]byte("a"), int64(1)}},
		{"*0\r\n", []any{}},
	}
	for _, c := range cases {
		got, err := readReply(bufio.NewReader(bytes.NewBufferString(c.in)))
		if err != nil {
			t.Fatalf("decoding %q: %v", c.in, err)
		}
		if !deepEqualAny(got, c.want) {
			t.Fatalf("decoding %q gave %#v, expected %#v", c.in, got, c.want)
		}
	}
	if _, err := readReply(bufio.NewReader(bytes.NewBufferString("-ERR nope\r\n"))); err == nil {
		t.Fatal("a RESP error reply was not turned into a Go error")
	}
}

func deepEqualAny(a, b any) bool {
	switch x := a.(type) {
	case []byte:
		y, ok := b.([]byte)
		return ok && bytes.Equal(x, y)
	case []any:
		y, ok := b.([]any)
		if !ok || len(x) != len(y) {
			return false
		}
		for i := range x {
			if !deepEqualAny(x[i], y[i]) {
				return false
			}
		}
		return true
	default:
		return a == b
	}
}

// ---------------------------------------------------------------- lease calibration

func TestLeaseCalibrationIsBoundedAndExplained(t *testing.T) {
	d, why := calibrateLeaseDuration(5 * time.Millisecond)
	if d != 100*time.Millisecond {
		t.Fatalf("a tiny p99 gave a %s lease, expected the 100 ms floor", d)
	}
	if !strings.Contains(why, "floor") {
		t.Fatalf("the calibration reasoning was not recorded: %q", why)
	}
	d2, _ := calibrateLeaseDuration(2 * time.Second)
	if d2 != 2*time.Second {
		t.Fatalf("a huge p99 gave %s, expected the 2 s cap", d2)
	}
	d3, _ := calibrateLeaseDuration(50 * time.Millisecond)
	if d3 != 200*time.Millisecond {
		t.Fatalf("4 x 50 ms should be 200 ms, got %s", d3)
	}
}

// TestFenceRefusesAFillThatStartedBeforeAnInvalidation is the study's central
// correctness rule at the store level: invalidation is not just a deletion. A reader
// whose fill began before the invalidation holds a committed but superseded state,
// and publishing it afterwards is how cache-aside produces a stale read that no
// "publish only after commit" rule prevents.
func TestFenceRefusesAFillThatStartedBeforeAnInvalidation(t *testing.T) {
	s := newMemStore(1 << 20)
	ctx := context.Background()
	f0, err := s.FenceOf(ctx, "k")
	if err != nil {
		t.Fatal(err)
	}
	if f0 != 0 {
		t.Fatalf("a never-invalidated key reports fence %d, expected 0", f0)
	}
	pre := newEntry(PortalContent{PersonID: 1, FullName: "pre"}, 0, 1, nowMS())
	pre.Fence = f0
	if ok, _ := s.Put(ctx, "k", pre, hardTTL); !ok {
		t.Fatal("the initial publish was refused")
	}

	// The write path invalidates while that fill is notionally still in flight.
	f1, err := s.Fence(ctx, "k")
	if err != nil {
		t.Fatal(err)
	}
	if f1 != f0+1 {
		t.Fatalf("the fence moved to %d, expected %d", f1, f0+1)
	}
	if _, hit, _ := s.Get(ctx, "k"); hit {
		t.Fatal("the invalidated value is still being served")
	}
	if ok, _ := s.Put(ctx, "k", pre, hardTTL); ok {
		t.Fatal("a fill that began before the invalidation republished its superseded content")
	}
	if s.Stats()["put_fenced"] == 0 {
		t.Fatal("the refused publication was not counted as fenced")
	}

	// A fill that began AFTER the invalidation publishes normally.
	post := newEntry(PortalContent{PersonID: 1, FullName: "post"}, 0, 2, nowMS())
	post.Fence = f1
	if ok, _ := s.Put(ctx, "k", post, hardTTL); !ok {
		t.Fatal("a fill that began after the invalidation was refused")
	}
}
