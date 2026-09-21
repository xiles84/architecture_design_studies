package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"sync/atomic"
	"time"

	"adsplatform/ports"
)

// The cache adapter: the policy that sits between the portal's reads and writes and
// the database.
//
// Three invariants are enforced here and nowhere else, because every scenario shares
// this code and a scenario that quietly implemented a different policy would make
// the comparison meaningless:
//
//  1. NOTHING IS PUBLISHED BEFORE THE DATABASE COMMITS. The relaxed paths commit
//     first and then touch the cache; the strict paths tombstone, commit, and only
//     then acknowledge. There is no code path that writes a cache value derived from
//     tentative application state.
//  2. EVERY PUBLISH IS FENCED. A publish is refused if the stored entry is newer,
//     so an expired lease holder cannot overwrite a newer committed version.
//  3. A FILL READS ONE SNAPSHOT. The portal's two statements run in one REPEATABLE
//     READ transaction, so the published content is a state the database actually
//     held. Reading them separately could stitch two states into one that never
//     existed -- an impossible cache value produced by the harness rather than by a
//     design, which is exactly the kind of false finding the gate exists to avoid.

type strictPolicy string

const (
	// PolicyInvalidationOnly: the writer's pre-commit fence is a sufficient proof of
	// freshness, because every writer goes through the same cache AND every reader
	// shares that cache's fence. The fence is what makes it sufficient: a plain
	// deletion would leave a reader that had already begun its fill free to
	// republish the superseded state it read.
	PolicyInvalidationOnly strictPolicy = "invalidation-only"
	// PolicyVersionValidate: a local cache in a multi-instance deployment cannot
	// rely on its own invalidation, so a strict reader validates the entry's
	// version against the authoritative token before serving it.
	PolicyVersionValidate strictPolicy = "version-validation"
	// PolicyAuthoritative: freshness cannot be proven from the cache at all
	// (no token, and another writer the adapter cannot see), so the strict reader
	// reads the database and the result records the cache as bypassed.
	PolicyAuthoritative strictPolicy = "authoritative-read"
)

// strictPolicyFor decides how a strict read proves freshness. The decision is a
// recorded property of the scenario, not an implementation detail: it is the answer
// to "what did strict cost here, and why".
func strictPolicyFor(d Design, instances int) strictPolicy {
	if d.Writers == WritersExt20 {
		// A writer that bypasses the adapter cannot be seen by it. On the legacy
		// model there is no token to check either, so no cache-side proof exists.
		return PolicyAuthoritative
	}
	switch d.Backend {
	case BackendRedis:
		// One shared store: a pre-commit tombstone is visible to every instance.
		return PolicyInvalidationOnly
	case BackendMemory:
		if instances > 1 {
			if d.usesVersion() {
				return PolicyVersionValidate
			}
			// Legacy + separate LRUs: instance B cannot see that instance A
			// invalidated, and there is no token to validate. Bypass, and say so.
			return PolicyAuthoritative
		}
		return PolicyInvalidationOnly
	}
	return PolicyInvalidationOnly
}

// mutation is one logical write the workload asks the adapter to perform.
type mutation struct {
	Kind        string // insert | correct | delete | person_update | reassign
	PersonID    int64
	NewPersonID int64
	DonationID  int64
	AmountCents int64
	FullName    string
	Email       string
	Donation    Donation
}

// keys are the cache keys this mutation affects. Reassignment affects two, which is
// the property that makes it the interesting mutation.
func (m mutation) keys() []int64 {
	if m.Kind == "reassign" {
		return []int64{m.PersonID, m.NewPersonID}
	}
	return []int64{m.PersonID}
}

type adapter struct {
	d    Design
	cat  *catalogue
	db   ports.DB
	ds   *Dataset
	orc  *oracle
	opts Options
	inst []*Instance
	log  *readLog

	policy     strictPolicy
	leaseTTL   time.Duration
	leaseCalib string

	seqGen atomic.Int64

	// phaseName is the phase currently running, used only to attribute a residual
	// stale read to what produced it.
	phaseName string

	// counters the report needs
	fills           atomic.Int64
	fillErrors      atomic.Int64
	publishes       atomic.Int64
	publishFenced   atomic.Int64
	publishFailed   atomic.Int64
	invalidations   atomic.Int64
	tombstonePre    atomic.Int64
	validationQ     atomic.Int64
	bypassReads     atomic.Int64
	fallbackReads   atomic.Int64
	leasesAcquired  atomic.Int64
	leasesContend   atomic.Int64
	leaseTimeout    atomic.Int64
	duplicateFills  atomic.Int64
	writeCount      atomic.Int64
	writeErrors     atomic.Int64
	writeConflicts  atomic.Int64
	ambiguousWrites atomic.Int64
	externalWrites  atomic.Int64
	extBumps        atomic.Int64
	strictWrong     atomic.Int64
	impossible      atomic.Int64
	// unrecordedConfirmed counts reads whose value was committed but not yet on the
	// ledger when the read finished. They are classified ahead, never impossible.
	unrecordedConfirmed atomic.Int64

	// fault injection, driven by the faults phase. Each is a one-shot or a switch,
	// set before the phase and cleared after it, so a fault cannot leak into the
	// next scenario in the same cell.
	faultSuppressInvalidation atomic.Bool
	faultCacheDown            atomic.Bool
	faultPublishFail          atomic.Bool
	faultSkipNextInvalidation atomic.Bool
	suppressedInvalidations   atomic.Int64

	// fillLatency samples the database fill path, which calibrates the lease.
	fillLatency []float64
	fillMu      chanMutex
}

func newAdapter(d Design, cat *catalogue, db ports.DB, ds *Dataset, orc *oracle, opts Options, insts []*Instance, log *readLog) *adapter {
	a := &adapter{
		d: d, cat: cat, db: db, ds: ds, orc: orc, opts: opts, inst: insts, log: log,
		policy:     strictPolicyFor(d, len(insts)),
		leaseTTL:   500 * time.Millisecond,
		leaseCalib: "not yet calibrated",
		fillMu:     newChanMutex(),
	}
	if !d.hasCache() {
		a.policy = PolicyAuthoritative
		a.leaseCalib = "no cache in this scenario"
	}
	return a
}

func (a *adapter) store(instIdx int) EntryStore { return a.inst[instIdx].Store }

func (a *adapter) nextSeq() int64 { return a.seqGen.Add(1) }

// ---------------------------------------------------------------- snapshot read

// readVersion reads the authoritative token from an open transaction.
func (a *adapter) readVersion(ctx context.Context, q ports.Queryer, personID int64) (int64, error) {
	args, err := a.cat.args(sPersonVersion, map[string]any{"person_id": personID})
	if err != nil {
		return 0, err
	}
	var v int64
	if err := q.QueryRow(ctx, a.cat.stmt(sPersonVersion).SQL, args...).Scan(&v); err != nil {
		return 0, err
	}
	return v, nil
}

// readPortalSnapshot reads the portal view in ONE repeatable-read transaction and
// returns the content, the authoritative version (owned model) and how long the
// fill took.
func (a *adapter) readPortalSnapshot(ctx context.Context, personID int64) (PortalContent, int64, time.Duration, error) {
	start := time.Now()
	// One implementation of the portal read, shared with the correctness gate and
	// the audit (readPortalSQL in verify.go). A second implementation could drift
	// from the measured one and then certify a payload the harness never reads.
	content, version, err := readPortalSQL(ctx, a.db, a.cat, a.d.usesVersion(), personID)
	if err != nil {
		return PortalContent{}, 0, 0, err
	}
	return content, version, time.Since(start), nil
}

// scanPortalPerson reads the one-row identity-and-aggregate block.
func scanPortalPerson(rows ports.Rows, personID int64) (PortalContent, error) {
	defer rows.Close()
	var c PortalContent
	var joined time.Time
	if !rows.Next() {
		return c, fmt.Errorf("portal view for person %d returned no row", personID)
	}
	if err := rows.Scan(&c.PersonID, &c.FullName, &c.Email, &joined, &c.CharityID,
		&c.CharityName, &c.CharityCountry, &c.DonationCount, &c.DonationTotalCents); err != nil {
		return c, err
	}
	// One canonical timestamp rendering, so a payload built by the oracle and one
	// built from the database hash identically.
	c.JoinedAt = joined.UTC().Format(time.RFC3339Nano)
	return c, rows.Err()
}

func scanPortalRecent(rows ports.Rows) ([]PortalDonation, error) {
	defer rows.Close()
	out := []PortalDonation{}
	for rows.Next() {
		var d PortalDonation
		var at time.Time
		var note *string
		if err := rows.Scan(&d.ID, &d.AmountCents, &d.Currency, &at, &note); err != nil {
			// The embedded design returns `note` as text; the row designs return a
			// typed column. Both scan into *string with pgx, and a NULL becomes nil.
			return nil, err
		}
		d.DonatedAt = at.UTC().Format(time.RFC3339Nano)
		if note != nil {
			d.Note = *note
		}
		out = append(out, d)
	}
	return out, rows.Err()
}

// ---------------------------------------------------------------- reads

// ReadPortal serves one portal read for logical instance instIdx and reports what
// it returned, from where, and how that compares with the oracle's requirement.
func (a *adapter) ReadPortal(ctx context.Context, instIdx int, r *rand.Rand) (ReadOutcome, error) {
	person := a.ds.PickPerson(r)
	return a.ReadKey(ctx, instIdx, person.ID, r)
}

func (a *adapter) ReadKey(ctx context.Context, instIdx int, keyID int64, r *rand.Rand) (ReadOutcome, error) {
	key := CacheKey(keyID)
	reqHash, reqSeq := a.orc.Required(keyID)
	out := ReadOutcome{KeyID: keyID, RequiredHash: reqHash, RequiredSeq: reqSeq}

	// No cache at all: the baseline. The read still reports itself so the oracle's
	// account of a correct design can be seen to work.
	if !a.d.hasCache() {
		content, _, _, err := a.readPortalSnapshot(ctx, keyID)
		if err != nil {
			return out, err
		}
		out.Source = SrcDatabase
		return a.finish(ctx, keyID, content, out, reqHash, reqSeq), nil
	}

	// A strict scenario whose freshness cannot be proven from the cache reads the
	// database every time, and says so. This is the honest outcome for a legacy
	// cache under external writers, and for a local cache in a multi-instance
	// legacy deployment.
	if a.policy == PolicyAuthoritative {
		a.bypassReads.Add(1)
		content, _, _, err := a.readPortalSnapshot(ctx, keyID)
		if err != nil {
			return out, err
		}
		out.Source = SrcBypass
		out.Cause = "strict freshness unsupported by the cache alone"
		return a.finish(ctx, keyID, content, out, reqHash, reqSeq), nil
	}

	if a.faultCacheDown.Load() {
		// Cache outage: the read must still be correct. It falls back to the
		// authoritative database and is counted.
		a.bypassReads.Add(1)
		content, _, _, err := a.readPortalSnapshot(ctx, keyID)
		if err != nil {
			return out, err
		}
		out.Source = SrcBypass
		out.Cause = "cache outage"
		return a.finish(ctx, keyID, content, out, reqHash, reqSeq), nil
	}

	store := a.store(instIdx % len(a.inst))
	now := nowMS()

	// --- hit path
	if e, hit, err := store.Get(ctx, key); err == nil && hit {
		out.CacheHit = true
		age := now - e.CreatedMS
		expiredBy := ""
		switch {
		case e.HardExpired(now):
			expiredBy = "hard"
		case ShouldExpireEarly(e.TimeRemainingS(now), hardTTLSeconds, r.Float64()):
			// Probabilistic early expiration: the entry is treated as expired
			// BEFORE it is served, and refreshed through the lease path.
			expiredBy = "probabilistic"
		}
		if expiredBy == "" {
			switch a.policy {
			case PolicyVersionValidate:
				v, verr := a.readVersion(ctx, a.db, keyID)
				a.validationQ.Add(1)
				if verr == nil && v == e.Version {
					out.AgeMS = age
					out.ClaimedVersion = e.Version
					out.Source = SrcValidated
					return a.finish(ctx, keyID, e.Content, out, reqHash, reqSeq), nil
				}
				// The local entry is behind the authoritative version: it must not
				// be served. Fall through to the fill path.
			default:
				out.AgeMS = age
				out.ClaimedVersion = e.Version
				out.Source = SrcHit
				return a.finish(ctx, keyID, e.Content, out, reqHash, reqSeq), nil
			}
		} else {
			out.ExpiredBy = expiredBy
		}
	} else if err != nil {
		// A backend error is not a stale read: the read falls back and is counted.
		a.bypassReads.Add(1)
		content, _, _, ferr := a.readPortalSnapshot(ctx, keyID)
		if ferr != nil {
			return out, ferr
		}
		out.Source = SrcBypass
		out.Cause = "cache outage"
		return a.finish(ctx, keyID, content, out, reqHash, reqSeq), nil
	}

	// --- fill path, through the per-key lease
	token := fmt.Sprintf("%d-%d", keyID, a.nextSeq())
	acquired, lerr := store.TryLease(ctx, key, token, a.leaseTTL)
	if lerr != nil {
		a.bypassReads.Add(1)
		content, _, _, ferr := a.readPortalSnapshot(ctx, keyID)
		if ferr != nil {
			return out, ferr
		}
		out.Source = SrcBypass
		out.Cause = "cache outage"
		return a.finish(ctx, keyID, content, out, reqHash, reqSeq), nil
	}
	if acquired {
		a.leasesAcquired.Add(1)
		defer func() { _ = store.ReleaseLease(context.Background(), key, token) }()

		// Double-check after acquiring: a previous holder may have published while
		// this reader waited for the lease.
		if e2, hit2, err2 := store.Get(ctx, key); err2 == nil && hit2 && !e2.HardExpired(nowMS()) {
			if a.policy != PolicyVersionValidate || a.versionMatches(ctx, keyID, e2) {
				out.CacheHit = true
				out.Source = SrcHit
				out.ClaimedVersion = e2.Version
				return a.finish(ctx, keyID, e2.Content, out, reqHash, reqSeq), nil
			}
		}
		// The fence is captured BEFORE the database snapshot. If the key is
		// invalidated while this fill is in flight, the fence moves and the publish
		// below is refused: the content read here is a committed but superseded
		// state, and caching it is exactly how cache-aside produces a stale read
		// that no "publish after commit" rule prevents.
		fence, ferr := store.FenceOf(ctx, key)
		if ferr != nil {
			a.fillErrors.Add(1)
			return out, ferr
		}
		content, version, dur, err := a.readPortalSnapshot(ctx, keyID)
		if err != nil {
			a.fillErrors.Add(1)
			return out, err
		}
		a.noteFillLatency(dur)
		a.fills.Add(1)
		e := newEntry(content, version, a.nextSeq(), nowMS())
		e.Fence = fence
		ok, perr := store.Put(ctx, key, e, hardTTL)
		if perr != nil {
			a.publishFailed.Add(1)
		} else if ok {
			a.publishes.Add(1)
		} else {
			a.publishFenced.Add(1)
		}
		out.Source = SrcFill
		out.ClaimedVersion = version
		return a.finish(ctx, keyID, content, out, reqHash, reqSeq), nil
	}

	// Contended: wait a bounded, jittered time, re-check, then fall back to the
	// database rather than waiting on someone else's lease forever.
	a.leasesContend.Add(1)
	wait, werr := waitForLease(ctx, a.leaseTTL, r)
	a.log.NoteLeaseWait(float64(wait.Milliseconds()))
	if werr != nil {
		a.leaseTimeout.Add(1)
	}
	if e3, hit3, err3 := store.Get(ctx, key); err3 == nil && hit3 && !e3.HardExpired(nowMS()) {
		if a.policy != PolicyVersionValidate || a.versionMatches(ctx, keyID, e3) {
			out.CacheHit = true
			out.Source = SrcHit
			out.ClaimedVersion = e3.Version
			return a.finish(ctx, keyID, e3.Content, out, reqHash, reqSeq), nil
		}
	}
	a.fallbackReads.Add(1)
	content, _, _, err := a.readPortalSnapshot(ctx, keyID)
	if err != nil {
		return out, err
	}
	out.Source = SrcFallback
	return a.finish(ctx, keyID, content, out, reqHash, reqSeq), nil
}

func (a *adapter) versionMatches(ctx context.Context, keyID int64, e *Entry) bool {
	v, err := a.readVersion(ctx, a.db, keyID)
	a.validationQ.Add(1)
	return err == nil && v == e.Version
}

// finish classifies the read against the oracle and records it. The strict
// violation and impossible-value counters are incremented HERE, at the moment of
// the read, so a cell cannot pass by aggregating its own evidence incorrectly.
func (a *adapter) finish(ctx context.Context, keyID int64, content PortalContent, out ReadOutcome, reqHash string, reqSeq int64) ReadOutcome {
	h := content.ContentHash()
	out.ReturnedHash = h
	kind, rseq, behind, at := a.orc.Classify(keyID, h, reqHash, reqSeq)
	if kind == KindImpossible {
		// An UNRECORDED state is not an impossible one. The oracle learns about a
		// committed state after the commit returns, so a reader racing that commit can
		// hold a state the ledger has not recorded yet. One conditional read settles
		// it, and it runs only when a hash is otherwise unknown -- never on a hit, so
		// it cannot distort the measurement it protects.
		if db, _, derr := readPortalSQL(ctx, a.db, a.cat, a.d.usesVersion(), keyID); derr == nil && db.ContentHash() == h {
			kind = KindAhead
			a.unrecordedConfirmed.Add(1)
		}
	}
	out.Kind = kind
	out.ReturnedSeq = rseq
	out.Behind = behind
	if a.orc.CurrentSeq(keyID) > reqSeq {
		out.Overlapped = true
	}
	if kind == KindStale && at > 0 {
		out.StaleMS = float64(nowMS() - at)
	}
	if kind == KindStale {
		out.Cause = a.staleCause(out)
	}
	if kind == KindStale && a.d.Freshness == FreshStrict {
		a.strictWrong.Add(1)
	}
	if kind == KindImpossible {
		a.impossible.Add(1)
	}
	a.log.Record(out)
	return out
}

// staleCause names the phase and the path a stale read came from, so a residual
// violation can be attributed instead of guessed at.
func (a *adapter) staleCause(out ReadOutcome) string {
	phase := a.phaseName
	if phase == "" {
		phase = "unattributed"
	}
	where := map[string]string{
		SrcHit:       "cache hit served a superseded entry",
		SrcValidated: "cache hit served an entry that passed validation",
		SrcFill:      "the refilling reader published before the write committed",
		SrcFallback:  "lease fallback read a superseded snapshot",
		SrcBypass:    "authoritative read returned a superseded snapshot",
		SrcDatabase:  "database read returned a superseded snapshot",
	}[out.Source]
	if where == "" {
		where = "unknown path " + out.Source
	}
	return phase + ": " + where
}

func (a *adapter) noteFillLatency(d time.Duration) {
	a.fillMu.Lock()
	defer a.fillMu.Unlock()
	if len(a.fillLatency) < 5000 {
		a.fillLatency = append(a.fillLatency, float64(d.Microseconds())/1000.0)
	}
}

// CalibrateLease sets the lease duration from the measured fill latency. It is
// called once, after a short calibration pass, before any measured phase.
func (a *adapter) CalibrateLease() {
	a.fillMu.Lock()
	samples := append([]float64(nil), a.fillLatency...)
	a.fillMu.Unlock()
	if len(samples) == 0 {
		a.leaseCalib = "no fill samples; lease left at the 500 ms default (NOT calibrated)"
		return
	}
	p99 := statFromSamples(samples).P99
	d, why := calibrateLeaseDuration(time.Duration(p99 * float64(time.Millisecond)))
	a.leaseTTL = d
	a.leaseCalib = why
}

// ---------------------------------------------------------------- writes

// Write applies one mutation. The ordering is the study's correctness floor and is
// not negotiable:
//
//	strict    tombstone -> COMMIT -> (republish for through) -> acknowledge
//	relaxed   COMMIT -> best-effort invalidate or republish -> acknowledge
//
// There is no path that touches the cache before the commit except installing a
// tombstone, and a tombstone contains no data, so no speculative or rolled-back
// value can ever be published.
func (a *adapter) Write(ctx context.Context, instIdx int, m mutation) error {
	n := a.writeCount.Add(1)
	ext := a.d.Writers == WritersExt20 && n%5 == 0
	strict := a.d.Freshness == FreshStrict
	keys := m.keys()

	// 1. Strict scenario: fence the cache BEFORE the authoritative mutation. The
	//    tombstone carries no data; it only makes readers fill again.
	if strict {
		for _, k := range keys {
			a.installTombstone(ctx, instIdx, k)
		}
	}

	// 2. The authoritative mutation, in one transaction with its version bump and
	//    its outbox event.
	retries, err := a.applyDB(ctx, m)
	if err != nil {
		a.writeErrors.Add(1)
		a.reconcileAmbiguous(ctx, m)
		return err
	}
	a.noteCommitted(m)

	// 3. An external writer commits and then does NOT touch the cache. It is the
	//    regime the strict legacy cells have to survive, so it must not be modelled
	//    as if the adapter had seen it.
	if ext {
		a.externalWrites.Add(1)
		return nil
	}

	// A scenario with no cache has no instance to bump, and must never be given a
	// pretend one: the baseline is the database alone.
	for _, k := range keys {
		if len(a.inst) > 0 {
			a.inst[instIdx%len(a.inst)].BumpGen(k)
		}
	}

	switch {
	case strict && a.d.Strategy == StrategyAside:
		// Fence a SECOND time, now that the mutation has committed. The pre-commit
		// fence invalidates the cache for the duration of the mutation; this one
		// bounds the fence token itself. Without it, a fill that read its fence in the
		// window between the pre-commit fence and the commit would read the
		// pre-mutation state and be allowed to publish it, because the fence it
		// captured was already the new one. That window is real under load: the first
		// version of this code recorded ~14 000 stale-after-ack reads in a strict
		// cell that fenced only before the commit.
		for _, k := range keys {
			a.installTombstone(ctx, instIdx, k)
		}
		// Leave the cache invalidated and only then acknowledge, per the contract.
		for _, k := range keys {
			a.log.NoteAck(k)
		}
		_ = retries
		return nil

	case strict && a.d.Strategy == StrategyThrough:
		// The same second fence, then republish the committed representation.
		for _, k := range keys {
			a.installTombstone(ctx, instIdx, k)
		}
		// The committed representation is the ONLY thing republished. A failed
		// republish leaves the tombstone installed in step 1 -- a correct, if
		// slower, state. Never the old value.
		for _, k := range keys {
			if a.publishFail() {
				a.publishFailed.Add(1)
				continue
			}
			if err := a.republish(ctx, instIdx, k); err != nil {
				a.publishFailed.Add(1)
			}
		}
		for _, k := range keys {
			a.log.NoteAck(k)
		}
		return nil

	case a.d.Freshness == FreshRelaxed && a.d.Strategy == StrategyAside:
		// Commit first, then best-effort invalidation. Stale committed reads are
		// permitted and every one of them is counted against this cell.
		for _, k := range keys {
			if a.suppressInvalidation() {
				continue
			}
			a.invalidate(ctx, instIdx, k)
		}
		for _, k := range keys {
			a.log.NoteAck(k)
		}
		return nil

	case a.d.Freshness == FreshRelaxed && a.d.Strategy == StrategyThrough:
		for _, k := range keys {
			if a.publishFail() {
				a.publishFailed.Add(1)
				continue
			}
			if err := a.republish(ctx, instIdx, k); err != nil {
				a.publishFailed.Add(1)
			}
		}
		for _, k := range keys {
			a.log.NoteAck(k)
		}
		return nil
	}

	// A no-cache scenario still records the acknowledgement so the oracle's
	// time-to-freshness clock runs; there is no cache to update.
	for _, k := range keys {
		a.log.NoteAck(k)
	}
	_ = retries
	return nil
}

// isExternalWrite decides, per write, whether this mutation belongs to the external
// fraction. It uses the adapter's own deterministic source so a run is reproducible
// from its seed.
func (a *adapter) isExternalWrite() bool {
	n := a.writeCount.Load()
	// 20 %: every fifth write, by count rather than by a random draw, so the
	// fraction is exactly 20 % at any run length and the regime is reproducible.
	return n%5 == 0
}

func (a *adapter) suppressInvalidation() bool {
	// A fault phase arms faultSkipNextInvalidation explicitly, so any relaxed cell
	// can be made to skip exactly one post-commit invalidation. The negative
	// control's own flag gates the second arm, so an ordinary cell cannot turn
	// itself into the control by accident.
	if a.faultSkipNextInvalidation.CompareAndSwap(true, false) {
		a.suppressedInvalidations.Add(1)
		return true
	}
	if !a.d.SuppressOneInvalidation {
		return false
	}
	if a.faultSuppressInvalidation.CompareAndSwap(true, false) {
		a.suppressedInvalidations.Add(1)
		return true
	}
	return false
}

func (a *adapter) publishFail() bool { return a.faultPublishFail.Load() }

// installTombstone removes any cached value for this key BEFORE the authoritative
// mutation starts. It is the strict writer's fence: a tombstone contains no data, so
// installing it early cannot publish anything speculative, and a reader that arrives
// during the mutation misses and fills from the database.
//
// With the memory backend and several logical instances, this clears only the
// WRITING instance's LRU -- which is the honest model of independent in-process
// caches, and the reason a multi-instance local cache needs version validation or an
// authoritative read to be strict at all.
func (a *adapter) installTombstone(ctx context.Context, instIdx int, keyID int64) {
	key := CacheKey(keyID)
	store := a.inst[instIdx%len(a.inst)].Store
	if a.d.Backend == BackendRedis {
		store = a.inst[0].Store
	}
	// Fence, not a bare delete: the deletion makes readers fill again, and the fence
	// makes any fill that already started refuse to publish what it read.
	_, _ = store.Fence(ctx, key)
	a.tombstonePre.Add(1)
}

func (a *adapter) invalidate(ctx context.Context, instIdx int, keyID int64) {
	if a.faultCacheDown.Load() {
		a.publishFailed.Add(1)
		return
	}
	key := CacheKey(keyID)
	store := a.inst[instIdx%len(a.inst)].Store
	if a.d.Backend == BackendRedis {
		store = a.inst[0].Store
	}
	// The relaxed writer's post-commit invalidation is the same operation as the
	// strict writer's pre-commit tombstone. What makes a scenario relaxed is WHEN it
	// runs and that it may be skipped or fail -- not a weaker fence.
	_, _ = store.Fence(ctx, key)
	a.invalidations.Add(1)
}

// reconcileAmbiguous handles a mutation whose outcome is unknown. The database may
// or may not hold it, so the cache must not be trusted either way: every affected
// key is invalidated and a later reader reconciles from the database. The mutation is
// NOT applied to the oracle, because an unacknowledged write is not a promise.
func (a *adapter) reconcileAmbiguous(ctx context.Context, m mutation) {
	if !a.d.hasCache() {
		return
	}
	for _, k := range m.keys() {
		a.installTombstone(ctx, int(m.PersonID)%len(a.inst), k)
	}
	if a.d.usesVersion() {
		a.ambiguousWrites.Add(1)
	}
}

// noteCommitted applies the mutation to the oracle, which is the only record of
// what the client was promised.
func (a *adapter) noteCommitted(m mutation) {
	switch m.Kind {
	case "insert":
		a.orc.ApplyInsert(m.PersonID, m.Donation)
	case "correct":
		a.orc.ApplyCorrect(m.PersonID, m.DonationID, m.AmountCents)
	case "delete":
		a.orc.ApplyDelete(m.PersonID, m.DonationID)
	case "person_update":
		a.orc.ApplyPersonUpdate(m.PersonID, m.FullName, m.Email)
	case "reassign":
		a.orc.ApplyReassign(m.DonationID, m.PersonID, m.NewPersonID, m.Donation.CharityID)
	}
}

// republish fills the committed view and publishes it. It runs ONLY after the
// mutation's transaction committed.
func (a *adapter) republish(ctx context.Context, instIdx int, keyID int64) error {
	store := a.store(instIdx)
	fence, err := store.FenceOf(ctx, CacheKey(keyID))
	if err != nil {
		return err
	}
	content, version, dur, err := a.readPortalSnapshot(ctx, keyID)
	if err != nil {
		return err
	}
	a.noteFillLatency(dur)
	a.fills.Add(1)
	e := newEntry(content, version, a.nextSeq(), nowMS())
	e.Fence = fence
	ok, err := store.Put(ctx, CacheKey(keyID), e, hardTTL)
	if err != nil {
		return err
	}
	if ok {
		a.publishes.Add(1)
	} else {
		a.publishFenced.Add(1)
	}
	return nil
}

// applyDB runs the mutation's database transaction, including the owned model's
// version machinery and concurrency strategy.
func (a *adapter) applyDB(ctx context.Context, m mutation) (int, error) {
	model := a.d.effectiveModel()
	if model != ModelOwned {
		tx, err := a.db.Begin(ctx, ports.ReadCommitted)
		if err != nil {
			return 0, err
		}
		if err := a.execMutation(ctx, tx, m); err != nil {
			_ = tx.Rollback(ctx)
			return 0, err
		}
		if err := tx.Commit(ctx); err != nil {
			return 0, err
		}
		return 0, nil
	}

	// Owned model.
	if a.d.Ver == VerUnsafe {
		// The negative control: read the version in the application and write back
		// a value derived from it, with no guard. Two concurrent writers can both
		// be acknowledged while one bump is lost.
		return a.applyUnsafe(ctx, m)
	}
	if a.d.Ver == VerPess {
		tx, err := a.db.Begin(ctx, ports.ReadCommitted)
		if err != nil {
			return 0, err
		}
		if err := a.lockForMutation(ctx, tx, m); err != nil {
			_ = tx.Rollback(ctx)
			return 0, err
		}
		if err := a.execMutation(ctx, tx, m); err != nil {
			_ = tx.Rollback(ctx)
			return 0, err
		}
		if err := a.bumpAndLog(ctx, tx, m, false, 0); err != nil {
			_ = tx.Rollback(ctx)
			return 0, err
		}
		if err := tx.Commit(ctx); err != nil {
			return 0, err
		}
		return 0, nil
	}

	// Optimistic: bounded retries around a compare-and-set on the version.
	for attempt := 0; attempt <= a.opts.Retries; attempt++ {
		tx, err := a.db.Begin(ctx, ports.ReadCommitted)
		if err != nil {
			return attempt, err
		}
		v1, err := a.readVersion(ctx, tx, m.PersonID)
		if err != nil {
			_ = tx.Rollback(ctx)
			return attempt, err
		}
		v2 := int64(0)
		if m.Kind == "reassign" {
			if v2, err = a.readVersion(ctx, tx, m.NewPersonID); err != nil {
				_ = tx.Rollback(ctx)
				return attempt, err
			}
		}
		if err := a.execMutation(ctx, tx, m); err != nil {
			_ = tx.Rollback(ctx)
			return attempt, err
		}
		n, err := a.exec(ctx, tx, sVersionCAS, map[string]any{"person_id": m.PersonID, "expected_version": v1})
		if err != nil {
			_ = tx.Rollback(ctx)
			return attempt, err
		}
		if n == 0 {
			_ = tx.Rollback(ctx)
			a.conflict()
			continue
		}
		if m.Kind == "reassign" {
			n2, err := a.exec(ctx, tx, sVersionCAS, map[string]any{"person_id": m.NewPersonID, "expected_version": v2})
			if err != nil {
				_ = tx.Rollback(ctx)
				return attempt, err
			}
			if n2 == 0 {
				_ = tx.Rollback(ctx)
				a.conflict()
				continue
			}
		}
		// The CAS is the bump, so the outbox row is appended without another bump.
		if err := a.outbox(ctx, tx, m); err != nil {
			_ = tx.Rollback(ctx)
			return attempt, err
		}
		if err := tx.Commit(ctx); err != nil {
			return attempt, err
		}
		return attempt, nil
	}
	return a.opts.Retries, fmt.Errorf("optimistic version CAS did not succeed after %d attempts", a.opts.Retries)
}

func (a *adapter) conflict() { a.writeConflicts.Add(1) }

func (a *adapter) applyUnsafe(ctx context.Context, m mutation) (int, error) {
	for attempt := 0; attempt <= a.opts.Retries; attempt++ {
		tx, err := a.db.Begin(ctx, ports.ReadCommitted)
		if err != nil {
			return attempt, err
		}
		v, err := a.readVersion(ctx, tx, m.PersonID)
		if err != nil {
			_ = tx.Rollback(ctx)
			return attempt, err
		}
		if err := a.execMutation(ctx, tx, m); err != nil {
			_ = tx.Rollback(ctx)
			return attempt, err
		}
		// Application-computed bump: this is where the lost update happens.
		if _, err := a.exec(ctx, tx, "w_version_set", map[string]any{"person_id": m.PersonID, "new_version": v + 1}); err != nil {
			_ = tx.Rollback(ctx)
			return attempt, err
		}
		if err := a.outbox(ctx, tx, m); err != nil {
			_ = tx.Rollback(ctx)
			return attempt, err
		}
		if err := tx.Commit(ctx); err != nil {
			return attempt, err
		}
		return attempt, nil
	}
	return a.opts.Retries, errors.New("unsafe control exhausted retries")
}

// lockForMutation is the pessimistic strategy: take the parent row lock BEFORE
// reading or writing anything. Reassignment locks TWO parents in a deterministic
// order, so two concurrent reassignments in opposite directions cannot deadlock.
func (a *adapter) lockForMutation(ctx context.Context, tx ports.Tx, m mutation) error {
	if m.Kind == "reassign" {
		a1, a2 := m.PersonID, m.NewPersonID
		if a1 > a2 {
			a1, a2 = a2, a1
		}
		args, err := a.cat.args(sLockPersons2, map[string]any{"person_a": a1, "person_b": a2})
		if err != nil {
			return err
		}
		rows, err := tx.Query(ctx, a.cat.stmt(sLockPersons2).SQL, args...)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			var id, v int64
			if err := rows.Scan(&id, &v); err != nil {
				return err
			}
		}
		return rows.Err()
	}
	args, err := a.cat.args(sLockPerson, map[string]any{"person_id": m.PersonID})
	if err != nil {
		return err
	}
	var v int64
	return tx.QueryRow(ctx, a.cat.stmt(sLockPerson).SQL, args...).Scan(&v)
}

func (a *adapter) execMutation(ctx context.Context, q ports.Queryer, m mutation) error {
	switch m.Kind {
	case "insert":
		d := m.Donation
		var note any
		if d.Note != nil {
			note = *d.Note
		}
		_, err := a.exec(ctx, q, sDonationInsert, map[string]any{
			"donation_id": d.ID, "person_id": d.PersonID, "charity_id": d.CharityID,
			"amount_cents": d.AmountCents, "currency": d.Currency, "donated_at": d.DonatedAt, "note": note,
		})
		return err
	case "correct":
		_, err := a.exec(ctx, q, sDonationCorrect, map[string]any{"donation_id": m.DonationID, "amount_cents": m.AmountCents})
		return err
	case "delete":
		_, err := a.exec(ctx, q, sDonationDelete, map[string]any{"donation_id": m.DonationID})
		return err
	case "person_update":
		_, err := a.exec(ctx, q, sPersonUpdate, map[string]any{"person_id": m.PersonID, "full_name": m.FullName, "email": m.Email})
		return err
	case "reassign":
		_, err := a.exec(ctx, q, sDonationReassign, map[string]any{"donation_id": m.DonationID, "new_person_id": m.NewPersonID})
		return err
	}
	return fmt.Errorf("unknown mutation kind %q", m.Kind)
}

// bumpAndLog performs the plain (non-CAS) version bump and the outbox append inside
// the mutation's transaction. This is the non-optimistic path; the optimistic path
// uses the CAS as its bump.
func (a *adapter) bumpAndLog(ctx context.Context, tx ports.Tx, m mutation, bumpFirst bool, expected int64) error {
	if _, err := a.exec(ctx, tx, sVersionBump, map[string]any{"person_id": m.PersonID}); err != nil {
		return err
	}
	if m.Kind == "reassign" {
		if _, err := a.exec(ctx, tx, sVersionBump, map[string]any{"person_id": m.NewPersonID}); err != nil {
			return err
		}
	}
	return a.outbox(ctx, tx, m)
}

// outbox appends one event per affected key, in the mutation's own transaction. If
// this fails the transaction rolls back and the mutation did not happen, which is
// what makes "the version moved" and "the log recorded why" the same event.
func (a *adapter) outbox(ctx context.Context, tx ports.Tx, m mutation) error {
	donationID := any(nil)
	if m.Kind != "person_update" {
		donationID = m.DonationID
	}
	if _, err := a.exec(ctx, tx, sOutboxInsert, map[string]any{
		"person_id": m.PersonID, "kind": m.Kind, "donation_id": donationID,
	}); err != nil {
		return err
	}
	if m.Kind == "reassign" {
		if _, err := a.exec(ctx, tx, sOutboxInsert, map[string]any{
			"person_id": m.NewPersonID, "kind": m.Kind, "donation_id": donationID,
		}); err != nil {
			return err
		}
	}
	return nil
}

func (a *adapter) exec(ctx context.Context, q ports.Queryer, name string, vals map[string]any) (int64, error) {
	args, err := a.cat.args(name, vals)
	if err != nil {
		return 0, err
	}
	n, err := q.Exec(ctx, a.cat.stmt(name).SQL, args...)
	if err != nil {
		return n, fmt.Errorf("%s: %w", name, err)
	}
	return n, nil
}

// writeConflictsCount counts optimistic CAS failures. It is on the adapter, not on
// the oracle, because a conflict is a property of the write strategy.
func (a *adapter) writeConflictsCount() int64 { return a.writeConflicts.Load() }
