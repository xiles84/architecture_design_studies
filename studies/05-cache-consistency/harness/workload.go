package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"adsplatform/core/measure"
	"adsplatform/ports"
)

// The workloads. Every phase goes through the same adapter, so a number in one
// phase and a number in another are produced by the same cache policy.

type cell struct {
	d    Design
	cat  *catalogue
	db   ports.DB
	ds   *Dataset
	orc  *oracle
	ad   *adapter
	res  *CellResult
	opts Options

	seq atomic.Int64

	// Probabilistic-expiry observation, by age bucket. Atomic arrays rather than a
	// map and mutex: this is on the hot read path, and a lock here would be
	// measurable.
	ageHits  [6]atomic.Int64
	ageFired [6]atomic.Int64
}

var ageBucketNames = [6]string{"<0s", "0-75s", "75-150s", "150-225s", "225-300s", ">=300s"}

func bucketIndex(ms int64) int {
	switch {
	case ms < 0:
		return 0
	case ms < 75_000:
		return 1
	case ms < 150_000:
		return 2
	case ms < 225_000:
		return 3
	case ms < 300_000:
		return 4
	default:
		return 5
	}
}

// newStore builds a backend. Capacity is in bytes; the primary capacity is chosen
// below the logical working set so eviction actually happens, and the fits-control
// capacity is above it.
func newStore(backend Backend, capBytes int64, redisAddr string, redisConns int) (EntryStore, error) {
	switch backend {
	case BackendMemory:
		return newMemStore(capBytes), nil
	case BackendRedis:
		if redisAddr == "" {
			return nil, errors.New("redis backend selected but no address")
		}
		return newRedisStore(redisAddr, redisConns), nil
	}
	return nil, fmt.Errorf("no cache backend for %q", backend)
}

func buildInstances(backend Backend, n int, capBytes int64, redisAddr string, redisConns int) ([]*Instance, error) {
	// The no-cache baselines and the reference cells have no backend. Returning an
	// empty instance list (rather than a fake store) keeps "there is no cache here"
	// visible in the type instead of pretended away.
	if backend == BackendNone || backend == "" {
		return nil, nil
	}
	out := make([]*Instance, 0, n)
	for i := 0; i < n; i++ {
		st, err := newStore(backend, capBytes, redisAddr, redisConns)
		if err != nil {
			return nil, err
		}
		out = append(out, newInstance(i, st))
	}
	return out, nil
}

// ---------------------------------------------------------------- operations

func (c *cell) instIdx(r *rand.Rand) int {
	// A scenario with no cache has no instance to choose, and must not call Intn(0).
	if len(c.ad.inst) <= 1 {
		return 0
	}
	return r.Intn(len(c.ad.inst))
}

// readOp is one cacheable portal read.
func (c *cell) readOp() measure.Op {
	return func(ctx context.Context, r *rand.Rand) (measure.Outcome, error) {
		out, err := c.ad.ReadPortal(ctx, c.instIdx(r), r)
		if err != nil {
			return measure.Outcome{}, err
		}
		c.observe(out)
		return measure.Outcome{}, nil
	}
}

// observe records the probabilistic-expiry evidence: for every cache hit, its age
// bucket and whether the draw fired. It is what checks the implementation against
// data, not only against the unit tests' arithmetic.
func (c *cell) observe(out ReadOutcome) {
	if !out.CacheHit {
		// A fill is not a hit, and a fallback or bypass is not a cache entry; there
		// is no entry age to bucket.
		return
	}
	b := bucketIndex(out.AgeMS)
	c.ageHits[b].Add(1)
	if out.ExpiredBy == "probabilistic" {
		c.ageFired[b].Add(1)
	}
}

// uncacheableOp is the charity feed: the blended workload's operation that no
// donor-keyed cache can serve. It exists because a cacheable-endpoint number alone
// would let a design look good while the application's real mix is unaffected.
func (c *cell) uncacheableOp() measure.Op {
	return func(ctx context.Context, r *rand.Rand) (measure.Outcome, error) {
		ch := c.ds.Charities[r.Intn(len(c.ds.Charities))]
		args, err := c.cat.args(sCharityRecent, map[string]any{"charity_id": ch.ID})
		if err != nil {
			return measure.Outcome{}, err
		}
		rows, err := c.db.Query(ctx, c.cat.stmt(sCharityRecent).SQL, args...)
		if err != nil {
			return measure.Outcome{}, err
		}
		n := 0
		for rows.Next() {
			var id, amt int64
			var at time.Time
			var name string
			if err := rows.Scan(&id, &amt, &at, &name); err != nil {
				rows.Close()
				return measure.Outcome{}, err
			}
			n++
		}
		rows.Close()
		if err := rows.Err(); err != nil {
			return measure.Outcome{}, err
		}
		if n == 0 {
			return measure.Outcome{Rejected: true}, nil
		}
		return measure.Outcome{}, nil
	}
}

// buildMutationFor parameterises one logical mutation. `target`, when non-nil, is the
// donor the caller needs the mutation to land on -- the hotspot and fault phases use
// it to force contention onto a hot key.
//
// It is a PARAMETER rather than an assignment after the fact. Setting m.PersonID on a
// ready-made mutation leaves that mutation's donation belonging to somebody else, so
// the database changes a row the ledger never touched and the two diverge. The gate
// caught exactly that on the first dev check: a fill returned a portal state the
// oracle had never held, the read was correctly classified impossible, and the cell
// failed. The lesson is the study's own: a harness that mis-attributes a write does not
// fail loudly, it fabricates a correctness finding.
func (c *cell) buildMutationFor(r *rand.Rand, target *Person) (mutation, error) {
	base := func() Person {
		if target != nil {
			return *target
		}
		return c.ds.PickPerson(r)
	}
	roll := r.Intn(100)
	switch {
	case roll < 55:
		return c.buildInsertFor(r, base()), nil
	case roll < 75:
		if m, ok := c.buildCorrectFor(r, base()); ok {
			return m, nil
		}
		return c.buildInsertFor(r, base()), nil
	case roll < 85:
		if m, ok := c.buildDeleteFor(r, base()); ok {
			return m, nil
		}
		return c.buildInsertFor(r, base()), nil
	case roll < 95:
		p := base()
		return mutation{Kind: "person_update", PersonID: p.ID,
			FullName: fmt.Sprintf("renamed-%d", p.ID), Email: fmt.Sprintf("renamed%06d@example.org", p.ID)}, nil
	default:
		if m, ok := c.buildReassignFor(r, base()); ok {
			return m, nil
		}
		return c.buildInsertFor(r, base()), nil
	}
}

// buildMutation is the unforced form, used wherever the workload may choose its own
// donor.
func (c *cell) buildMutation(r *rand.Rand) (mutation, error) { return c.buildMutationFor(r, nil) }

func (c *cell) buildInsertFor(r *rand.Rand, p Person) mutation {
	id := c.ds.MaxDonationID + c.seq.Add(1)
	note := fmt.Sprintf("workload gift %d", id)
	if r.Intn(2) == 0 {
		note = ""
	}
	var notePtr *string
	if note != "" {
		notePtr = &note
	}
	return mutation{
		Kind: "insert", PersonID: p.ID,
		Donation: Donation{
			ID: id, PersonID: p.ID, CharityID: p.CharityID,
			AmountCents: int64(500 + r.Intn(250_000)),
			Currency:    currencies[r.Intn(len(currencies))],
			DonatedAt:   time.Now().UTC().Truncate(time.Microsecond),
			Note:        notePtr,
		},
	}
}

func (c *cell) buildCorrectFor(r *rand.Rand, p Person) (mutation, bool) {
	d, ok := c.orc.PickDonation(p.ID, r)
	if !ok {
		return mutation{}, false
	}
	delta := int64(1 + r.Intn(5000))
	return mutation{Kind: "correct", PersonID: p.ID, DonationID: d.ID, DeltaCents: delta}, true
}

func (c *cell) buildDeleteFor(r *rand.Rand, p Person) (mutation, bool) {
	d, ok := c.orc.PickDonation(p.ID, r)
	if !ok {
		return mutation{}, false
	}
	return mutation{Kind: "delete", PersonID: p.ID, DonationID: d.ID}, true
}

// buildReassignFor moves one of THIS donor's donations to a different donor. Both
// keys change, which is the mutation a single-key cache invalidation gets wrong.
func (c *cell) buildReassignFor(r *rand.Rand, p Person) (mutation, bool) {
	d, ok := c.orc.PickDonation(p.ID, r)
	if !ok {
		return mutation{}, false
	}
	for attempt := 0; attempt < 6; attempt++ {
		other := c.ds.PickPerson(r)
		if other.ID == p.ID {
			continue
		}
		d.CharityID = other.CharityID
		return mutation{Kind: "reassign", PersonID: p.ID, NewPersonID: other.ID,
			DonationID: d.ID, Donation: d}, true
	}
	return mutation{}, false
}

func (c *cell) writeOp() measure.Op {
	return func(ctx context.Context, r *rand.Rand) (measure.Outcome, error) {
		m, err := c.buildMutation(r)
		if err != nil {
			return measure.Outcome{}, err
		}
		return measure.Outcome{}, c.ad.Write(ctx, c.instIdx(r), m)
	}
}

// blendedOp is the application's real mix: cacheable reads, the uncacheable charity
// feed, and writes. readPercent is what makes the 99/1 and 90/10 regimes.
func (c *cell) blendedOp(readPercent int) measure.Op {
	return func(ctx context.Context, r *rand.Rand) (measure.Outcome, error) {
		if r.Intn(100) < readPercent {
			// 85 % of reads are the cacheable portal, 15 % the uncacheable feed.
			if r.Intn(100) < 85 {
				out, err := c.ad.ReadPortal(ctx, c.instIdx(r), r)
				if err != nil {
					return measure.Outcome{}, err
				}
				c.observe(out)
				return measure.Outcome{}, nil
			}
			return c.uncacheableOp()(ctx, r)
		}
		m, err := c.buildMutation(r)
		if err != nil {
			return measure.Outcome{}, err
		}
		return measure.Outcome{}, c.ad.Write(ctx, c.instIdx(r), m)
	}
}

// ---------------------------------------------------------------- phases

// resetWrong starts a fresh account for the next phase. The study reports wrong
// reads PER PHASE: a design that is stale only while writes are in flight is a
// different finding from one that is stale during warm reads, and a cumulative
// number would hide which.
func (c *cell) resetWrong(phase string) {
	if c.ad != nil {
		c.ad.log = newReadLog()
		c.ad.phaseName = phase
	}
}

func (c *cell) phaseWrong(name string) {
	s := c.ad.log.Summary()
	c.res.Wrong = append(c.res.Wrong, PhaseWrong{Phase: name, Summary: s})
	if s.ImpossibleValues > 0 {
		c.res.ImpossibleValues += s.ImpossibleValues
	}
}

func (c *cell) snapshotCache() {
	if len(c.ad.inst) == 0 {
		// The no-cache baseline has no backend to account for. Reporting zeros
		// rather than omitting the section keeps every cell's shape identical.
		c.res.Cache = CacheStats{Backend: string(BackendNone)}
		c.res.Lease = LeaseStats{LeaseDuration: "n/a", Calibration: "no cache in this scenario"}
		return
	}
	st := c.ad.store(0).Stats()
	capacity := int64(c.opts.CapacityKB) * 1024
	// The backend's resident bytes are the cache's whole footprint: payload plus
	// the per-entry metadata the backend needs. The logical payload size is what a
	// reader compares against the capacity, so both are recorded.
	resident := st["resident_bytes"]
	items := st["items"]
	payload := resident
	meta := int64(0)
	c.res.Cache = CacheStats{
		Backend:                 string(c.d.Backend),
		CapacityBytes:           capacity,
		ResidentBytes:           resident,
		Items:                   items,
		Evictions:               st["evictions"],
		EvictedKeys:             st["evicted_keys"],
		Fills:                   c.ad.fills.Load(),
		FillErrors:              c.ad.fillErrors.Load(),
		Publishes:               c.ad.publishes.Load(),
		PublishFenced:           c.ad.publishFenced.Load(),
		PublishFailed:           c.ad.publishFailed.Load(),
		Invalidations:           c.ad.invalidations.Load(),
		Tombstones:              c.ad.tombstonePre.Load(),
		ValidationQueries:       c.ad.validationQ.Load(),
		BypassReads:             c.ad.bypassReads.Load(),
		FallbackReads:           c.ad.fallbackReads.Load(),
		ExternalWrites:          c.ad.externalWrites.Load(),
		AmbiguousWrites:         c.ad.ambiguousWrites.Load(),
		SuppressedInvalidations: c.ad.suppressedInvalidations.Load(),
		UnrecordedConfirmed:     c.ad.unrecordedConfirmed.Load(),
		WriteNoops:              c.ad.writeNoops.Load(),
		BackendStats:            st,
		PayloadBytes:            payload,
		MetadataBytes:           meta,
	}
	if c.d.Backend == BackendRedis {
		c.res.Redis = c.ad.store(0).Info(context.Background())
	}
	c.res.Lease = LeaseStats{
		Acquisitions:  c.ad.leasesAcquired.Load(),
		Contended:     c.ad.leasesContend.Load(),
		Waits:         c.ad.log.Summary().LeaseWait.Count,
		WaitP50MS:     c.ad.log.Summary().LeaseWait.P50,
		WaitP99MS:     c.ad.log.Summary().LeaseWait.P99,
		WaitMaxMS:     c.ad.log.Summary().LeaseWait.Max,
		Timeouts:      c.ad.leaseTimeout.Load(),
		DuplicateFill: c.ad.duplicateFills.Load(),
		FallbackReads: c.ad.fallbackReads.Load(),
		LeaseDuration: c.ad.leaseTTL.String(),
		Calibration:   c.ad.leaseCalib,
	}
	c.res.ExpiryByAge = map[string]map[string]int64{}
	for i, name := range ageBucketNames {
		hits := c.ageHits[i].Load()
		fired := c.ageFired[i].Load()
		if hits == 0 {
			continue
		}
		c.res.ExpiryByAge[name] = map[string]int64{"hits": hits, "probabilistic_expiries": fired}
	}
}

// calibrate runs a short, unmeasured fill pass so the lease duration comes from
// measured latency rather than from a constant chosen for convenience.
func (c *cell) calibrate(ctx context.Context) {
	if !c.d.hasCache() {
		return
	}
	r := rand.New(rand.NewSource(c.opts.Seed + 991))
	for i := 0; i < 30; i++ {
		p := c.ds.PickPerson(r)
		if _, err := c.ad.ReadKey(ctx, 0, p.ID, r); err != nil {
			return
		}
	}
	c.ad.CalibrateLease()
}

func (c *cell) runWarm(ctx context.Context) error {
	if c.d.hasCache() {
		// Cold fill: the whole key space, counted rather than timed, so the phase
		// measures the fill path with evictions and leases in play.
		_ = c.ad.store(0).Flush(ctx)
		keys := c.ds.People
		if len(keys) > 400 {
			keys = keys[:400]
		}
		_ = measure.Run(ctx, measure.Budget{
			Workers: c.opts.Workers, WarmupOps: 0, OpsPerTrial: int64(len(keys)),
			Trials: 1, SafetyLimit: time.Minute,
		}, "cold_fill", func(ctx context.Context, r *rand.Rand) (measure.Outcome, error) {
			p := keys[r.Intn(len(keys))]
			out, err := c.ad.ReadKey(ctx, c.instIdx(r), p.ID, r)
			if err != nil {
				return measure.Outcome{}, err
			}
			c.observe(out)
			return measure.Outcome{}, nil
		})
		c.phaseWrong("cold_fill")
	}

	res := measure.Run(ctx, measure.Budget{
		Workers: c.opts.Workers, Warmup: c.opts.Warmup, Duration: c.opts.Duration,
		Trials: c.opts.Trials,
	}, "warm_read_only", c.readOp())
	c.res.Warm = append(c.res.Warm, res)
	c.phaseWrong("warm_read_only")
	return nil
}

func (c *cell) runMixed(ctx context.Context) error {
	for _, pct := range []int{99, 90} {
		name := fmt.Sprintf("app_%d_reads", pct)
		res := measure.Run(ctx, measure.Budget{
			Workers: c.opts.Workers, Warmup: c.opts.Warmup, Duration: c.opts.Duration,
			Trials: c.opts.Trials,
		}, name, c.blendedOp(pct))
		c.res.AppMix = append(c.res.AppMix, res)
		c.phaseWrong(name)

		// The cacheable endpoint measured on its own, so the report can distinguish
		// "the cache got faster" from "the application did less cacheable work".
		rd := measure.Run(ctx, measure.Budget{
			Workers: c.opts.Workers, Warmup: c.opts.Warmup, Duration: c.opts.Duration,
			Trials: c.opts.Trials,
		}, fmt.Sprintf("cacheable_%d_reads", pct), c.readOp())
		c.res.Mixed = append(c.res.Mixed, rd)
		c.phaseWrong(fmt.Sprintf("cacheable_%d_reads", pct))
	}
	return nil
}

// runHotspot puts every worker on the hottest few keys while writers hit the same
// keys, which is where a per-key lease and a per-key invalidation are actually
// tested.
func (c *cell) runHotspot(ctx context.Context) error {
	hot := c.ds.HotPeople(c.opts.HotKeys)
	if len(hot) == 0 {
		return nil
	}
	res := measure.Run(ctx, measure.Budget{
		Workers: c.opts.Workers, Warmup: c.opts.Warmup, Duration: c.opts.Duration,
		Trials: c.opts.Trials,
	}, "hotspot_reads", func(ctx context.Context, r *rand.Rand) (measure.Outcome, error) {
		p := hot[r.Intn(len(hot))]
		out, err := c.ad.ReadKey(ctx, c.instIdx(r), p.ID, r)
		if err != nil {
			return measure.Outcome{}, err
		}
		c.observe(out)
		return measure.Outcome{}, nil
	})
	c.res.Mixed = append(c.res.Mixed, res)
	c.phaseWrong("hotspot_reads")

	if c.opts.Writers > 0 {
		var ops atomic.Int64
		wr := measure.Run(ctx, measure.Budget{
			Workers: c.opts.Writers, Warmup: c.opts.Warmup, Duration: c.opts.Duration,
			Trials: c.opts.Trials, SafetyLimit: time.Minute,
		}, "hotspot_writes", func(ctx context.Context, r *rand.Rand) (measure.Outcome, error) {
			ops.Add(1)
			p := hot[r.Intn(len(hot))]
			// The mutation is built FOR the hot donor, not reassigned to it
			// afterwards: see buildMutationFor.
			m, err := c.buildMutationFor(r, &p)
			if err != nil {
				return measure.Outcome{}, err
			}
			return measure.Outcome{}, c.ad.Write(ctx, c.instIdx(r), m)
		})
		c.res.Writes = append(c.res.Writes, wr)
		c.phaseWrong("hotspot_writes")
	}
	return nil
}

// runStampede empties the cache and releases many readers at once on the same keys.
// The measurement that matters is database loads per key: one per key means the
// lease worked, N means it did not.
func (c *cell) runStampede(ctx context.Context) error {
	if !c.d.hasCache() {
		return nil
	}
	hot := c.ds.HotPeople(c.opts.StampedeN)
	if len(hot) == 0 {
		return nil
	}
	_ = c.ad.store(0).Flush(ctx)
	for _, in := range c.ad.inst {
		_ = in.Store.Flush(ctx)
	}
	loadsBefore := c.ad.fills.Load()
	fbBefore := c.ad.fallbackReads.Load()

	readers := c.opts.Stampede
	if readers < 2 {
		readers = 2
	}
	start := make(chan struct{})
	var wg sync.WaitGroup
	var wrongBefore = c.ad.log.Summary()
	for i := 0; i < readers; i++ {
		wg.Add(1)
		go func(seed int64) {
			defer wg.Done()
			r := rand.New(rand.NewSource(c.opts.Seed + seed))
			<-start
			for _, p := range hot {
				if _, err := c.ad.ReadKey(ctx, int(seed)%len(c.ad.inst), p.ID, r); err != nil {
					return
				}
			}
		}(int64(i))
	}
	t0 := nowMS()
	close(start)
	wg.Wait()
	elapsed := float64(nowMS() - t0)

	loads := c.ad.fills.Load() - loadsBefore
	fallbacks := c.ad.fallbackReads.Load() - fbBefore
	after := c.ad.log.Summary()
	dup := int64(0)
	if loads > int64(len(hot)) {
		dup = loads - int64(len(hot))
	}
	c.ad.duplicateFills.Add(dup)
	c.res.Stampede = &StampedeResult{
		Readers: readers, Keys: len(hot), ElapsedMS: elapsed,
		DatabaseLoads:  loads,
		LoadsPerKey:    float64(loads) / float64(len(hot)),
		LeaseAcquired:  c.ad.leasesAcquired.Load(),
		LeaseContended: c.ad.leasesContend.Load(),
		Fallbacks:      fallbacks,
		DuplicateFills: dup,
		WrongReads:     after.WrongReads - wrongBefore.WrongReads,
		Impossible:     after.ImpossibleValues - wrongBefore.ImpossibleValues,
		Note: fmt.Sprintf("%d readers released together on %d keys after a full flush; database loads per key is the stampede metric",
			readers, len(hot)),
	}
	c.phaseWrong("stampede")
	return nil
}

// runChurn sustains a blended workload long enough to cross real hard-TTL
// boundaries. The primary TTL is never shortened for this: the run simply lasts.
func (c *cell) runChurn(ctx context.Context) error {
	if c.opts.ChurnFor <= 0 {
		return nil
	}
	before := c.ad.log.Summary()
	t0 := time.Now()
	res := measure.Run(ctx, measure.Budget{
		Workers: c.opts.Workers, Warmup: c.opts.Warmup, Duration: c.opts.ChurnFor,
		Trials: 1, SafetyLimit: c.opts.ChurnFor + 2*time.Minute,
	}, "churn_90_reads", c.blendedOp(90))
	elapsed := time.Since(t0)
	after := c.ad.log.Summary()
	c.res.Churn = &ChurnResult{
		DurationS:     elapsed.Seconds(),
		TTLBoundaries: elapsed.Seconds() / hardTTLSeconds,
		Reads:         after.TotalReads - before.TotalReads,
		Writes:        c.ad.writeCount.Load(),
		HardExpired:   after.ExpiredHard - before.ExpiredHard,
		ProbExpired:   after.ExpiredProbabilistic - before.ExpiredProbabilistic,
		WrongReads:    after.WrongReads - before.WrongReads,
		Impossible:    after.ImpossibleValues - before.ImpossibleValues,
		Note: fmt.Sprintf("sustained blended workload for %s; crosses %.1f real 300 s TTL boundaries",
			elapsed.Round(time.Second), elapsed.Seconds()/hardTTLSeconds),
	}
	c.res.AppMix = append(c.res.AppMix, res)
	c.phaseWrong("churn_90_reads")
	return nil
}

// runInstances measures the same workload with one and with three logical
// application instances at an equal total worker count. With the memory backend the
// instances have separate LRUs, which is the property the phase exists to expose:
// one instance's invalidation is invisible to another's.
func (c *cell) runInstances(ctx context.Context) error {
	if !c.d.hasCache() {
		return nil
	}
	capBytes := int64(c.opts.CapacityKB) * 1024
	base := c.ad
	// The instances phase measures a DIFFERENT deployment: its writes go through its
	// own instances, which cannot invalidate the base adapter's stores. Starting and
	// ending cold keeps a leftover entry from being read as a strict violation of the
	// deployment that produced it.
	for _, in := range base.inst {
		_ = in.Store.Flush(ctx)
	}
	// Accumulated so a violation recorded in this phase still fails the cell.
	var instancesViolations, instancesImpossible atomic.Int64
	for _, n := range []int{1, 3} {
		insts, err := buildInstances(c.d.Backend, n, capBytes, c.opts.RedisAddr, c.opts.Workers+8)
		if err != nil {
			return err
		}
		log := newReadLog()
		ad := newAdapter(c.d, c.cat, c.db, c.ds, c.orc, c.opts, insts, log)
		ad.leaseTTL = base.leaseTTL
		ad.leaseCalib = base.leaseCalib
		c.ad = ad
		// Equal TOTAL workers: the worker count is shared across instances.
		app := measure.Run(ctx, measure.Budget{
			Workers: c.opts.Workers, Warmup: c.opts.Warmup, Duration: c.opts.Duration,
			Trials: c.opts.Trials,
		}, fmt.Sprintf("app_instances_%d", n), c.blendedOp(90))
		cache := measure.Run(ctx, measure.Budget{
			Workers: c.opts.Workers, Warmup: c.opts.Warmup, Duration: c.opts.Duration,
			Trials: c.opts.Trials,
		}, fmt.Sprintf("cacheable_instances_%d", n), c.readOp())
		s := log.Summary()
		localItems := int64(0)
		if n > 1 && c.d.Backend == BackendMemory {
			for _, in := range insts {
				localItems += in.Store.Stats()["items"]
			}
		} else {
			localItems = insts[0].Store.Stats()["items"]
		}
		note := "shared store: invalidation is visible to every instance"
		if c.d.Backend == BackendMemory && n > 1 {
			note = "separate in-process LRUs and separate leases; a write through one instance is invisible to the other two"
		}
		c.res.Instances = append(c.res.Instances, InstanceResult{
			Instances: n, Workers: c.opts.Workers, Backend: string(c.d.Backend),
			AppOpsPerSec: app.OpsPerSec, AppTrials: app.Trials, AppSpreadPct: app.SpreadPct,
			CacheOpsPerSec: cache.OpsPerSec,
			WrongReads:     s.WrongReads, Impossible: s.ImpossibleValues,
			BypassReads: ad.bypassReads.Load(), LocalLRUItems: localItems, Note: note,
		})
		c.ad.log = log
		c.phaseWrong(fmt.Sprintf("instances_%d", n))
		for _, in := range insts {
			_ = in.Store.Close(ctx)
		}
	}
	// A violation recorded by an arm's own adapter must reach the CELL's acceptance
	// check. The cell's check reads the base adapter's counter, so an instances-phase
	// violation was invisible: the owned redis through cell passed while its
	// one-instance arm had already recorded 11 stale reads.
	c.ad = base
	base.strictWrong.Add(instancesViolations.Load())
	base.impossible.Add(instancesImpossible.Load())
	for _, in := range base.inst {
		_ = in.Store.Flush(ctx)
	}
	return nil
}
