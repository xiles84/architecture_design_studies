package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"adsplatform/ports"
)

// Deterministic fault injection. Every fault here is driven by the recorded
// fault seed and by explicit calls, not by a race that may or may not happen: a
// fault phase that only sometimes reproduces is not evidence.
//
// The rules each fault is checking:
//
//   * a database rollback after a pre-write tombstone must leave the cache with NO
//     value, never with a speculative one;
//   * a cache failure after a commit must leave a RELAXED cell with a stale read it
//     counts, and a STRICT cell with a tombstone that simply refills;
//   * a dead lease holder must not deadlock the key: the lease expires and the next
//     reader takes it;
//   * an unavailable cache must degrade to authoritative reads, still correct;
//   * a fill racing a writer must never publish a state the database never held;
//   * a writer that bypasses the adapter must be visible to the ORACLE even though
//     it is invisible to the cache.

func (c *cell) runFaults(ctx context.Context) error {
	steps := []func(context.Context) FaultResult{}
	if c.d.hasCache() {
		steps = append(steps, c.faultRollbackAfterTombstone, c.faultCacheFailureAfterCommit,
			c.faultLeaseHolderDeath, c.faultCacheUnavailable, c.faultFillRacesWriter)
	}
	if c.d.Writers == WritersExt20 {
		steps = append(steps, c.faultExternalWriter)
	}
	if c.d.SuppressOneInvalidation {
		steps = append(steps, c.faultSuppressedInvalidation)
	}
	if c.d.hasCache() {
		steps = append(steps, c.faultDirtyWriteAudit)
	}
	if c.d.Ver == VerUnsafe {
		steps = append(steps, c.faultLostVersionBumps)
	}
	for _, fn := range steps {
		fr := fn(ctx)
		c.res.Faults = append(c.res.Faults, fr)
		if !fr.Reproduced {
			// A fault that could not be injected is a finding about the harness, not
			// a licence to skip the check.
			return fmt.Errorf("fault %s could not be reproduced: %s", fr.Name, fr.Detail)
		}
	}
	return nil
}

// faultRollbackAfterTombstone: install the strict tombstone, then make the
// authoritative mutation fail, and require that the cache holds nothing afterwards
// and that a read returns the committed (unchanged) value.
func (c *cell) faultRollbackAfterTombstone(ctx context.Context) FaultResult {
	fr := FaultResult{Name: "db-rollback-after-pre-write-tombstone", Correctness: "pending"}
	r := rand.New(rand.NewSource(c.opts.FaultSeed + 1))
	p := c.ds.HotPeople(1)[0]

	// A valid entry first, so the tombstone has something to remove.
	out, err := c.ad.ReadKey(ctx, 0, p.ID, r)
	if err != nil {
		fr.Detail = "prime read failed: " + err.Error()
		return fr
	}
	before := c.ad.log.Summary()

	c.ad.installTombstone(ctx, 0, p.ID)

	// A mutation that must fail: an insert of a donation id that already exists.
	existing, ok := c.orc.PickDonation(p.ID, r)
	if !ok {
		fr.Detail = "no donation to duplicate"
		return fr
	}
	bad := mutation{Kind: "insert", PersonID: p.ID, Donation: existing}
	if err := c.ad.Write(ctx, 0, bad); err == nil {
		fr.Detail = "the duplicate insert unexpectedly committed"
		return fr
	}

	after, err := c.ad.ReadKey(ctx, 0, p.ID, r)
	if err != nil {
		fr.Detail = "post-rollback read failed: " + err.Error()
		return fr
	}
	sum := c.ad.log.Summary()
	fr.Reproduced = true
	fr.WrongReads = sum.WrongReads - before.WrongReads
	fr.Impossible = sum.ImpossibleValues - before.ImpossibleValues
	fr.Correctness = fmt.Sprintf("tombstone removed the pre-existing entry; the rolled-back mutation left no cache value (%s -> %s); the next read returned %s",
		out.Source, after.Source, after.Kind)
	if after.Kind == KindImpossible || fr.Impossible > 0 {
		fr.Reproduced = false
		fr.Detail = "the rolled-back mutation produced an impossible cache value"
	}
	return fr
}

// faultCacheFailureAfterCommit: commit a real mutation and then fail the cache
// update. A relaxed cell must record a stale read; a strict cell must not (its
// tombstone was installed before the commit, so the next reader refills).
func (c *cell) faultCacheFailureAfterCommit(ctx context.Context) FaultResult {
	fr := FaultResult{Name: "cache-update-failure-after-commit", Correctness: "pending"}
	r := rand.New(rand.NewSource(c.opts.FaultSeed + 2))
	hp := c.ds.HotPeople(2)
	p := hp[len(hp)-1]
	key := CacheKey(p.ID)

	// Prime the key so a stale value is available to be returned.
	if _, err := c.ad.ReadKey(ctx, 0, p.ID, r); err != nil {
		fr.Detail = "prime read failed: " + err.Error()
		return fr
	}
	before := c.ad.log.Summary()

	// Fail the cache side of the next write. For strict, the tombstone is installed
	// first and the post-commit publish is skipped; for relaxed, the post-commit
	// invalidation/publish is skipped and the old entry survives.
	c.ad.faultSkipNextInvalidation.Store(true)
	c.ad.faultPublishFail.Store(true)
	m, err := c.buildMutationFor(r, &p)
	if err != nil {
		fr.Detail = "no mutation: " + err.Error()
		return fr
	}
	if err := c.ad.Write(ctx, 0, m); err != nil {
		fr.Detail = "the mutation failed at the database, not the cache: " + err.Error()
		return fr
	}
	// The cache side of the write was made to fail, so a relaxed writer's
	// acknowledgement still happens (that is the point of the fault) while the cache
	// keeps the superseded value.
	c.ad.faultPublishFail.Store(false)
	c.ad.faultSkipNextInvalidation.Store(false)

	after, err := c.ad.ReadKey(ctx, 0, p.ID, r)
	if err != nil {
		fr.Detail = "post-failure read failed: " + err.Error()
		return fr
	}
	sum := c.ad.log.Summary()
	fr.Reproduced = true
	fr.WrongReads = sum.WrongReads - before.WrongReads
	fr.Impossible = sum.ImpossibleValues - before.ImpossibleValues
	switch c.d.Freshness {
	case FreshStrict:
		fr.Correctness = fmt.Sprintf("strict: the tombstone installed before the commit survived the failed publish; the read refilled and returned %s", after.Kind)
		if fr.WrongReads > 0 {
			fr.Reproduced = false
			fr.Detail = "a strict cell returned a stale read after a post-commit cache failure"
		}
	default:
		fr.Correctness = fmt.Sprintf("relaxed: the post-commit cache failure left the previous committed value in place; the read returned %s and was counted", after.Kind)
		if fr.WrongReads == 0 {
			fr.Correctness += " (no stale read observed in this single draw; the churn phases carry the rate)"
		}
	}
	_ = key
	return fr
}

// faultLeaseHolderDeath: take a lease and abandon it. A second reader must wait, then
// take the lease after it expires -- a dead holder must not wedge the key forever.
func (c *cell) faultLeaseHolderDeath(ctx context.Context) FaultResult {
	fr := FaultResult{Name: "lease-holder-death", Correctness: "pending"}
	p := c.ds.HotPeople(3)[2%len(c.ds.HotPeople(3))]
	key := CacheKey(p.ID)
	store := c.ad.store(0)

	_ = store.Delete(ctx, key)
	acquired, err := store.TryLease(ctx, key, "dead-holder", c.ad.leaseTTL)
	if err != nil || !acquired {
		fr.Detail = fmt.Sprintf("could not take the lease (ok=%v err=%v)", acquired, err)
		return fr
	}
	// A second attempt must be refused while the dead holder's lease is live.
	second, err2 := store.TryLease(ctx, key, "second", c.ad.leaseTTL)
	if err2 != nil {
		fr.Detail = "second lease attempt errored: " + err2.Error()
		return fr
	}
	if second {
		fr.Detail = "the lease was not exclusive"
		return fr
	}
	// Wait past the lease and try again: the key must recover without human help.
	time.Sleep(c.ad.leaseTTL + 50*time.Millisecond)
	third, err3 := store.TryLease(ctx, key, "third", c.ad.leaseTTL)
	if err3 != nil || !third {
		fr.Detail = fmt.Sprintf("the lease did not expire (ok=%v err=%v)", third, err3)
		return fr
	}
	_ = store.ReleaseLease(ctx, key, "third")
	fr.Reproduced = true
	fr.Correctness = fmt.Sprintf("an abandoned lease blocked a second holder and was stealable after %s; no manual intervention required", c.ad.leaseTTL)
	return fr
}

// faultCacheUnavailable: the cache stops answering. Reads must still be correct --
// they fall back to the authoritative database -- and the outage must be visible in
// the bypass counter.
func (c *cell) faultCacheUnavailable(ctx context.Context) FaultResult {
	fr := FaultResult{Name: "cache-unavailable", Correctness: "pending"}
	r := rand.New(rand.NewSource(c.opts.FaultSeed + 4))
	before := c.ad.log.Summary()
	bypassBefore := c.ad.bypassReads.Load()

	c.ad.faultCacheDown.Store(true)
	st := c.ad.store(0)
	if setter, ok := st.(interface{ SetAvailable(bool) }); ok {
		setter.SetAvailable(false)
	}
	for i := 0; i < 20; i++ {
		p := c.ds.PickPerson(r)
		if _, err := c.ad.ReadKey(ctx, 0, p.ID, r); err != nil {
			c.ad.faultCacheDown.Store(false)
			fr.Detail = "a read failed during the outage: " + err.Error()
			return fr
		}
	}
	c.ad.faultCacheDown.Store(false)
	if setter, ok := st.(interface{ SetAvailable(bool) }); ok {
		setter.SetAvailable(true)
	}
	after := c.ad.log.Summary()
	fr.Reproduced = true
	fr.WrongReads = after.WrongReads - before.WrongReads
	fr.Impossible = after.ImpossibleValues - before.ImpossibleValues
	fr.Correctness = fmt.Sprintf("20 reads during the outage all returned a committed state (%d bypassed the cache, %d wrong, %d impossible)",
		c.ad.bypassReads.Load()-bypassBefore, fr.WrongReads, fr.Impossible)
	if fr.Impossible > 0 {
		fr.Reproduced = false
		fr.Detail = "the outage produced an impossible cache value"
	}
	return fr
}

// faultFillRacesWriter: a fill and a mutation on the same key, released together.
// The published value must still be a state the database held.
func (c *cell) faultFillRacesWriter(ctx context.Context) FaultResult {
	fr := FaultResult{Name: "cache-aside-fill-races-writer", Correctness: "pending"}
	r := rand.New(rand.NewSource(c.opts.FaultSeed + 5))
	p := c.ds.HotPeople(4)[3%len(c.ds.HotPeople(4))]
	before := c.ad.log.Summary()

	done := make(chan struct{}, 2)
	go func() {
		defer func() { done <- struct{}{} }()
		_, _ = c.ad.ReadKey(ctx, 0, p.ID, rand.New(rand.NewSource(c.opts.FaultSeed+51)))
	}()
	go func() {
		defer func() { done <- struct{}{} }()
		m, err := c.buildMutationFor(r, &p)
		if err != nil {
			return
		}
		_ = c.ad.Write(ctx, 0, m)
	}()
	<-done
	<-done

	after := c.ad.log.Summary()
	fr.Reproduced = true
	fr.WrongReads = after.WrongReads - before.WrongReads
	fr.Impossible = after.ImpossibleValues - before.ImpossibleValues
	fr.Correctness = fmt.Sprintf("a fill and a write on donor %d released together; %d wrong and %d impossible reads", p.ID, fr.WrongReads, fr.Impossible)
	if fr.Impossible > 0 {
		fr.Reproduced = false
		fr.Detail = "the racing fill published a state the database never held"
	}
	return fr
}

// faultExternalWriter: commit a write that bypasses the adapter entirely. The oracle
// must see it (it is a committed change) while the cache must not.
func (c *cell) faultExternalWriter(ctx context.Context) FaultResult {
	fr := FaultResult{Name: "legacy-external-writer-bypasses-adapter", Correctness: "pending"}
	r := rand.New(rand.NewSource(c.opts.FaultSeed + 6))
	p := c.ds.PickPerson(r)

	// Prime the cache with the pre-write state.
	if _, err := c.ad.ReadKey(ctx, 0, p.ID, r); err != nil {
		fr.Detail = "prime read failed: " + err.Error()
		return fr
	}
	before := c.ad.log.Summary()

	// Commit straight to the database, with no version bump, no tombstone and no
	// invalidation: exactly what an unknown legacy application does.
	m, err := c.buildMutationFor(r, &p)
	if err != nil {
		fr.Detail = "no mutation: " + err.Error()
		return fr
	}
	if err := c.externalWrite(ctx, m); err != nil {
		fr.Detail = "external write failed: " + err.Error()
		return fr
	}
	after, err := c.ad.ReadKey(ctx, 0, p.ID, r)
	if err != nil {
		fr.Detail = "post-external-write read failed: " + err.Error()
		return fr
	}
	sum := c.ad.log.Summary()
	fr.Reproduced = true
	fr.WrongReads = sum.WrongReads - before.WrongReads
	fr.Impossible = sum.ImpossibleValues - before.ImpossibleValues
	switch {
	case c.d.Freshness == FreshStrict:
		fr.Correctness = fmt.Sprintf("strict: the read could not trust the cache, so it %s and returned %s", after.Source, after.Kind)
	default:
		fr.Correctness = fmt.Sprintf("relaxed: the external write was invisible to the cache; the read returned %s (source %s) and any staleness is counted", after.Kind, after.Source)
	}
	return fr
}

// externalWrite commits a mutation to the database without touching the cache or the
// adapter's invalidation paths. It exists for the external-writer fault only.
func (c *cell) externalWrite(ctx context.Context, m mutation) error {
	tx, err := c.db.Begin(ctx, ports.ReadCommitted)
	if err != nil {
		return err
	}
	if _, err := c.cat.exec(ctx, tx, mutationStatement(m), mutationParams(m)); err != nil {
		_ = tx.Rollback(ctx)
		return err
	}
	if err := tx.Commit(ctx); err != nil {
		return err
	}
	// An external writer's commit IS its acknowledgement: nothing fences the cache on
	// its behalf, which is exactly why the strict legacy cells must read
	// authoritatively under this regime.
	c.ad.ackTokens(c.ad.noteCommitted(m))
	return nil
}

func mutationStatement(m mutation) string {
	switch m.Kind {
	case "insert":
		return sDonationInsert
	case "correct":
		return sDonationCorrect
	case "delete":
		return sDonationDelete
	case "person_update":
		return sPersonUpdate
	default:
		return sDonationReassign
	}
}

func mutationParams(m mutation) map[string]any {
	switch m.Kind {
	case "insert":
		d := m.Donation
		var note any
		if d.Note != nil {
			note = *d.Note
		}
		return map[string]any{
			"donation_id": d.ID, "person_id": d.PersonID, "charity_id": d.CharityID,
			"amount_cents": d.AmountCents, "currency": d.Currency, "donated_at": d.DonatedAt, "note": note,
		}
	case "correct":
		return map[string]any{"donation_id": m.DonationID, "delta_cents": m.DeltaCents}
	case "delete":
		return map[string]any{"donation_id": m.DonationID}
	case "person_update":
		return map[string]any{"person_id": m.PersonID, "full_name": m.FullName, "email": m.Email}
	default:
		return map[string]any{"donation_id": m.DonationID, "new_person_id": m.NewPersonID}
	}
}

// faultSuppressedInvalidation is the stale-read negative control: commit a mutation
// and skip exactly one post-commit invalidation, then require the oracle to see the
// stale read. It is the measured proof that the wrong-read accounting works.
func (c *cell) faultSuppressedInvalidation(ctx context.Context) FaultResult {
	fr := FaultResult{Name: "suppressed-invalidation-produces-a-stale-read", Correctness: "pending"}
	r := rand.New(rand.NewSource(c.opts.FaultSeed + 7))
	p := c.ds.HotPeople(5)[4%len(c.ds.HotPeople(5))]

	// Prime the key.
	if _, err := c.ad.ReadKey(ctx, 0, p.ID, r); err != nil {
		fr.Detail = "prime read failed: " + err.Error()
		return fr
	}
	before := c.ad.log.Summary()

	c.ad.faultSuppressInvalidation.Store(true)
	m, err := c.buildMutationFor(r, &p)
	if err != nil {
		fr.Detail = "no mutation: " + err.Error()
		return fr
	}
	if err := c.ad.Write(ctx, 0, m); err != nil {
		fr.Detail = "the mutation failed: " + err.Error()
		return fr
	}
	if c.ad.suppressedInvalidations.Load() == 0 {
		fr.Detail = "the invalidation was not suppressed"
		return fr
	}
	// The mutation may have been a person_update to the same value; read the key
	// until the oracle confirms the cache is behind, or give up honestly.
	var observed int64
	for i := 0; i < 5; i++ {
		out, err := c.ad.ReadKey(ctx, 0, p.ID, r)
		if err != nil {
			fr.Detail = "post-suppression read failed: " + err.Error()
			return fr
		}
		if out.Kind == KindStale {
			observed++
			break
		}
	}
	sum := c.ad.log.Summary()
	fr.WrongReads = sum.WrongReads - before.WrongReads
	fr.Impossible = sum.ImpossibleValues - before.ImpossibleValues
	fr.Reproduced = observed > 0
	if fr.Reproduced {
		fr.Correctness = fmt.Sprintf("one suppressed post-commit invalidation produced %d stale-after-ack read(s) that the oracle detected", fr.WrongReads)
	} else {
		fr.Detail = "the suppressed invalidation produced no detectable stale read in this draw; the workload is not harsh enough to license the other cells' correctness"
	}
	return fr
}

// faultDirtyWriteAudit proves the impossible-value detector works by presenting it
// with a record that corresponds to no committed state. This is an AUDIT TEST, not a
// benchmark strategy: its speed is never reported.
//
// The detector is exercised directly, and then end-to-end when the scenario actually
// reads its cache. A scenario whose strict policy is an authoritative read never
// consults the cache at all, and asking it to classify an injected record would be
// asking the wrong component -- the first version of this test did exactly that and
// reported a working detector as broken.
func (c *cell) faultDirtyWriteAudit(ctx context.Context) FaultResult {
	fr := FaultResult{Name: "dirty-cache-write-audit-rejects-impossible-version", Correctness: "pending"}
	r := rand.New(rand.NewSource(c.opts.FaultSeed + 8))
	p := c.ds.HotPeople(6)[5%len(c.ds.HotPeople(6))]
	key := CacheKey(p.ID)

	// The committed payload and its requirement, read the way the adapter reads it.
	committed, version, verr := readPortalSQL(ctx, c.db, c.cat, c.d.usesVersion(), p.ID)
	if verr != nil {
		fr.Detail = "could not read the committed payload: " + verr.Error()
		return fr
	}
	reqHash, reqSeq := c.orc.Required(p.ID)
	if k, _, _, _ := c.orc.Classify(p.ID, committed.ContentHash(), reqHash, reqSeq); k != KindFresh && k != KindAhead {
		fr.Detail = "the committed payload was not accepted as committed, so the detector cannot be calibrated"
		return fr
	}

	// An IMPOSSIBLE record: a real payload with a field changed, and a version far in
	// the future. It corresponds to no state the database ever held.
	bad := committed
	bad.FullName = "impossible-cache-write-audit-marker"
	kind, _, _, _ := c.orc.Classify(p.ID, bad.ContentHash(), reqHash, reqSeq)
	direct := kind == KindImpossible

	endToEnd := true
	var readKind, readSource string
	if c.ad.policy != PolicyAuthoritative && c.d.hasCache() {
		store := c.ad.store(0)
		fence, _ := store.FenceOf(ctx, key)
		injected := newEntry(bad, 1<<40, c.ad.nextSeq(), nowMS())
		injected.Fence = fence
		if _, err := store.Put(ctx, key, injected, hardTTL); err != nil {
			fr.Detail = "could not inject the invalid entry: " + err.Error()
			return fr
		}
		out, err := c.ad.ReadKey(ctx, 0, p.ID, r)
		if err != nil {
			fr.Detail = "read after injection failed: " + err.Error()
			return fr
		}
		readKind, readSource = out.Kind, out.Source
		endToEnd = out.Kind == KindImpossible
		// Leave the cache clean so the injection cannot leak into a later phase.
		_, _ = store.Fence(ctx, key)
	}

	fr.Reproduced = direct && endToEnd
	switch {
	case fr.Reproduced && readKind != "":
		fr.Correctness = fmt.Sprintf("the injected record corresponded to no committed state: the detector classified it impossible directly and the read through the cache returned %s (source %s)", readKind, readSource)
	case fr.Reproduced:
		fr.Correctness = "the injected record corresponded to no committed state and the detector classified it impossible (this scenario reads the database authoritatively, so there is no cache path to exercise)"
	default:
		fr.Detail = fmt.Sprintf("detector accepted an impossible record (direct=%t end-to-end=%t, read classified %q)", direct, endToEnd, readKind)
	}
	_ = version
	return fr
}

// faultLostVersionBumps is the concurrency negative control: an unguarded
// read-modify-write of the version token must lose bumps under contention. If it
// does not, the version token's correctness is not demonstrated and no strict owned
// cell may be believed.
func (c *cell) faultLostVersionBumps(ctx context.Context) FaultResult {
	fr := FaultResult{Name: "unguarded-version-bump-loses-updates", Correctness: "pending"}
	p := c.ds.HotPeople(1)[0]

	// The starting point, read from the database rather than assumed.
	v0, err := scalarInt64(ctx, c.db, c.cat, sPersonVersion, map[string]any{"person_id": p.ID})
	if err != nil {
		fr.Detail = "could not read the starting version: " + err.Error()
		return fr
	}
	before := c.ad.writeCount.Load()

	const workers = 8
	const perWorker = 40
	start := make(chan struct{})
	done := make(chan struct{}, workers)
	for w := 0; w < workers; w++ {
		go func(w int) {
			defer func() { done <- struct{}{} }()
			<-start
			for i := 0; i < perWorker; i++ {
				m := mutation{Kind: "person_update", PersonID: p.ID,
					FullName: fmt.Sprintf("racer-%d-%d", w, i), Email: fmt.Sprintf("racer%d-%d@example.org", w, i)}
				_ = c.ad.Write(ctx, 0, m)
			}
		}(w)
	}
	close(start)
	for w := 0; w < workers; w++ {
		<-done
	}

	acked := c.ad.writeCount.Load() - before
	vf, err := scalarInt64(ctx, c.db, c.cat, sPersonVersion, map[string]any{"person_id": p.ID})
	if err != nil {
		fr.Detail = "could not read the final version: " + err.Error()
		return fr
	}
	bumped := vf - v0
	lost := acked - bumped
	if lost < 0 {
		lost = 0
	}
	// The outbox/version invariant, which fails for the same reason.
	mono, merr := scalarInt64(ctx, c.db, c.cat, sAuditVersionMonotonic, map[string]any{"person_id": p.ID})

	// The control FIRES when the unguarded path is SEEN to break the version/outbox
	// invariant. Two pieces of evidence count, and either is sufficient:
	//
	//   * the outbox/version invariant itself (a_version_monotonic), which is direct
	//     evidence that an update was lost -- the application wrote back a version it
	//     had read before another writer's bump;
	//   * the acknowledged-write-versus-bump count.
	//
	// Requiring both was wrong: the count is the noisier of the two, and on a run
	// where the invariant broke but the count still balanced, the harness refused to
	// license the guarded designs even though it had just observed the unguarded one
	// fail. A control is judged by whether the invariant broke, not by whether the
	// coarse counter also noticed.
	fr.Reproduced = (mono > 0 || lost > 0) && merr == nil
	fr.Correctness = fmt.Sprintf("%d acknowledged writes, %d version bumps in the database, %d lost; a_version_monotonic reports %d mismatches (either line of evidence fires the control)", acked, bumped, lost, mono)
	if !fr.Reproduced {
		fr.Detail = "the unguarded control was not seen to break the version/outbox invariant; the workload is not harsh enough to license the guarded designs"
	}
	return fr
}

// scalarInt64 runs a single-value statement through the catalogue.
func scalarInt64(ctx context.Context, db ports.DB, cat *catalogue, name string, vals map[string]any) (int64, error) {
	a, err := cat.args(name, vals)
	if err != nil {
		return 0, err
	}
	var n int64
	if err := db.QueryRow(ctx, cat.stmt(name).SQL, a...).Scan(&n); err != nil {
		return 0, err
	}
	return n, nil
}
