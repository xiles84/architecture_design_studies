package main

import "strings"

// This file is the study's scenario registry. A scenario is a CACHE scenario
// (model x backend x strategy x freshness x write regime) or a DATABASE reference
// cell measured with no cache.
//
// The registry is a product of dimensions rather than a directory per cell: the
// cache scenarios share one SQL catalogue per database model, and only the
// harness's behaviour differs between them. A directory per cell would have made
// eighteen copies of identical SQL and hidden the fact that the only thing
// changing is which cache policy brackets the same reads and writes.

// Model is which database model the scenario runs against.
type Model string

const (
	// ModelOwned: the database may be changed for cache correctness
	// (cache_version + transactional outbox + atomic version changes).
	ModelOwned Model = "owned"
	// ModelLegacy: D3 as it is, with an external adapter on top of it.
	ModelLegacy Model = "legacy"
	// ModelRef: a database-layout reference cell, measured with no cache.
	ModelRef Model = "ref"
	// ModelControl: a negative control, measured like a design and never reported
	// as a fast result.
	ModelControl Model = "ctl"
)

// VerMode is the owned model's write strategy for the version token, and the
// concurrency strategy the owner's requirement 4 asks to compare.
type VerMode string

const (
	VerOpt    VerMode = "opt"    // optimistic: bump only if the version has not moved
	VerPess   VerMode = "pess"   // pessimistic: lock the parent row first
	VerNa     VerMode = "na"     // the legacy model has no version token
	VerUnsafe VerMode = "unsafe" // negative control: unguarded read-modify-write
)

// Backend is where the cache lives.
type Backend string

const (
	BackendNone   Backend = "none"
	BackendMemory Backend = "memory"
	BackendRedis  Backend = "redis"
)

// Strategy is when the cache is written.
type Strategy string

const (
	StrategyNone    Strategy = "none"
	StrategyAside   Strategy = "aside"
	StrategyThrough Strategy = "through"
)

// Freshness is the read contract.
type Freshness string

const (
	FreshNone    Freshness = "none"
	FreshRelaxed Freshness = "relaxed"
	FreshStrict  Freshness = "strict"
)

// WritersRegime says who is allowed to write the database.
type WritersRegime string

const (
	WritersCoord WritersRegime = "coord"
	// WritersExt20: 20 % of harness writes bypass the cache adapter and commit
	// straight to the database. A strict cache cannot see them, so it must
	// validate authoritatively or bypass itself -- and the result says which.
	WritersExt20 WritersRegime = "ext20"
)

// Statement names, as constants so a typo is a compile error rather than a
// design that quietly does less than its name claims.
const (
	sPortalPerson  = "r_portal_person"
	sPortalRecent  = "r_portal_recent"
	sCharityRecent = "r_charity_recent"
	sPersonVersion = "r_person_version"

	sLoadDonations    = "w_load_donations"
	sDonationInsert   = "w_donation_insert"
	sDonationCorrect  = "w_donation_correct"
	sDonationDelete   = "w_donation_delete"
	sPersonUpdate     = "w_person_update"
	sDonationReassign = "w_donation_reassign"

	sVersionBump  = "w_version_bump"
	sOutboxInsert = "w_outbox_insert"
	sVersionCAS   = "w_version_cas"
	sLockPerson   = "w_lock_person"
	sLockPersons2 = "w_lock_persons2"

	sAuditAggregate        = "a_person_aggregate"
	sAuditDonationIDs      = "a_donation_ids"
	sAuditRecentIDs        = "a_recent_ids"
	sAuditRolldownDrift    = "a_rolldown_drift"
	sAuditVersionMonotonic = "a_version_monotonic"
	sAuditOutboxOrphans    = "a_outbox_orphans"
	sAuditRollupDrift      = "a_rollup_drift"
	sAuditEmbeddedDrift    = "a_embedded_drift"
)

// Design is one scenario.
type Design struct {
	ID      string
	Short   string
	Group   string
	Title   string
	Summary string
	// Risk names the correctness risk this scenario carries, for the report.
	Risk string

	// OwnDir holds schema.sql, indexes.sql and (where present) triggers.sql.
	OwnDir string
	// SQLDir holds queries.sql, writes.sql and audit.sql. It usually equals
	// OwnDir; the placement pair points it at the flattened catalogue so those two
	// cells differ by exactly one decision.
	SQLDir string

	Model Model
	// ModelUsed is the database model a control borrows its SQL catalogue from.
	// The controls reuse a model's SQL but are not that model.
	ModelUsed Model
	Ver       VerMode
	Backend   Backend
	Strategy  Strategy
	Freshness Freshness
	Writers   WritersRegime

	// SuppressOneInvalidation makes the writer skip exactly one post-commit
	// invalidation. It is the stale-read negative control.
	SuppressOneInvalidation bool
	NegativeControl         bool

	// YBOnly marks SQL that uses YugabyteDB-only syntax.
	YBOnly bool

	// Pair names the scenario this one differs from by exactly one decision.
	Pair string

	// Needs lists statements that must exist in the scenario's catalogue.
	Needs []string
}

// baseNeeds is the statement set every scenario drives.
func baseNeeds() []string {
	return []string{
		sPortalPerson, sPortalRecent, sCharityRecent,
		sLoadDonations, sDonationInsert, sDonationCorrect, sDonationDelete,
		sPersonUpdate, sDonationReassign,
		sAuditAggregate, sAuditDonationIDs, sAuditRecentIDs,
	}
}

var ownedExtra = []string{
	sPersonVersion, sVersionBump, sOutboxInsert, sVersionCAS, sLockPerson, sLockPersons2,
	sAuditVersionMonotonic, sAuditOutboxOrphans,
}

func withNeeds(extra ...string) []string {
	out := append([]string{}, baseNeeds()...)
	return append(out, extra...)
}

// scenarioShort is the compact label the report tables use. It keeps every
// dimension of the scenario id, because a table row whose label hides which
// freshness contract it ran under is exactly the kind of artefact this study
// exists to avoid.
func scenarioShort(m Model, v VerMode, b Backend, s Strategy, f Freshness, w WritersRegime) string {
	bk := map[Backend]string{BackendNone: "db", BackendMemory: "mem", BackendRedis: "red"}[b]
	st := map[Strategy]string{StrategyNone: "", StrategyAside: "asd", StrategyThrough: "wth"}[s]
	fr := map[Freshness]string{FreshNone: "", FreshRelaxed: "rel", FreshStrict: "str"}[f]
	out := string(m[:1]) + ":" + string(v) + ":" + bk
	if st != "" {
		out += ":" + st
	}
	if fr != "" {
		out += ":" + fr
	}
	if w == WritersExt20 {
		out += ":ext"
	}
	return strings.Trim(out, ":")
}

func scenarioID(m Model, v VerMode, b Backend, s Strategy, f Freshness, w WritersRegime) string {
	return strings.Join([]string{string(m), string(v), string(b), string(s), string(f), string(w)}, "-")
}

// matrixScenarios expands the core matrix: two owned no-cache baselines, one
// legacy no-cache baseline, the cached block, the cached pessimistic cell, and the
// two legacy external-writer cells.
func matrixScenarios() []Design {
	var out []Design

	build := func(m Model, v VerMode, b Backend, s Strategy, f Freshness, w WritersRegime, title, summary, risk string, neg bool) Design {
		d := Design{
			ID: scenarioID(m, v, b, s, f, w), Short: scenarioShort(m, v, b, s, f, w),
			Group: string(m), Model: m, Ver: v, Backend: b, Strategy: s, Freshness: f, Writers: w,
			Title: title, Summary: summary, Risk: risk, NegativeControl: neg,
			OwnDir: string(m) + "/d3_" + string(m), Needs: withNeeds(),
		}
		if m == ModelOwned {
			d.Needs = withNeeds(ownedExtra...)
		}
		return d
	}

	// --- no-cache baselines: the version/outbox overhead must be visible before
	// any cached comparison is made.
	out = append(out, build(ModelOwned, VerOpt, BackendNone, StrategyNone, FreshNone, WritersCoord,
		"owned, optimistic, no cache",
		"Every read goes to the database. The mutation's transaction bumps cache_version with an optimistic compare-and-set and appends an outbox event. This is the floor the owned cached cells are compared against: it already pays for the machinery the cache will use.",
		"the version bump can conflict under contention and be retried (bounded)", false))
	out = append(out, build(ModelOwned, VerPess, BackendNone, StrategyNone, FreshNone, WritersCoord,
		"owned, pessimistic, no cache",
		"Byte-identical read and mutation SQL to the optimistic baseline; the only difference is that the transaction locks the parent row before touching anything, so competitors wait instead of being rejected.",
		"waiters hold a connection; a long transaction serialises the whole donor", false))
	out = append(out, build(ModelLegacy, VerNa, BackendNone, StrategyNone, FreshNone, WritersCoord,
		"legacy, no cache (D3 flattened-FK reference)",
		"The legacy model with no cache at all: Study 01's D3 schema and queries, remeasured in this run and environment. This cell is also the flattened-FK reference of the database-reference set, because the legacy model IS D3 -- running the same SQL twice would produce two cells of the same code.",
		"none specific", false))

	// --- the cached block.
	for _, m := range []Model{ModelOwned, ModelLegacy} {
		for _, b := range []Backend{BackendMemory, BackendRedis} {
			for _, s := range []Strategy{StrategyAside, StrategyThrough} {
				for _, f := range []Freshness{FreshRelaxed, FreshStrict} {
					v := VerOpt
					if m == ModelLegacy {
						v = VerNa
					}
					title := string(m) + ", " + string(b) + ", " + string(s) + ", " + string(f)
					out = append(out, build(m, v, b, s, f, WritersCoord, title, cacheSummary(m, b, s, f), cacheRisk(s, f), false))
				}
			}
		}
	}

	// --- one cached pessimistic cell, so the concurrency pair is measured with a
	// cache in front of it and not only on the no-cache baseline.
	out = append(out, build(ModelOwned, VerPess, BackendRedis, StrategyAside, FreshStrict, WritersCoord,
		"owned, pessimistic, redis, cache-aside, strict",
		"The strict owned cache-aside cell with the pessimistic writer instead of the optimistic one. Same reads, same cache policy, same contract; only who arbitrates the parent row changes.",
		"waiters hold a connection on the write path", false))

	// --- the legacy external-writer regime. A strict cache cannot see a writer
	// that bypasses it, so these two cells are where that limit is measured.
	out = append(out, build(ModelLegacy, VerNa, BackendRedis, StrategyAside, FreshRelaxed, WritersExt20,
		"legacy, redis, cache-aside, relaxed, 20% external writers",
		"One in five harness writes commits straight to the database without passing through the cache adapter. A relaxed cache accepts the staleness this creates; the wrong-read count is the price of the throughput beside it.",
		"external writes are invisible to the adapter; every resulting stale read must be counted", false))
	out = append(out, build(ModelLegacy, VerNa, BackendRedis, StrategyAside, FreshStrict, WritersExt20,
		"legacy, redis, cache-aside, strict, 20% external writers",
		"Relaxed is not acceptable here, so the strict reader must prove freshness itself. The legacy model offers no token, so the only sound proof is an authoritative database read: the cell reports the cache as bypassed rather than pretending the cache is strict.",
		"strict freshness is unsupported by the cache alone; the result must say so", false))

	return out
}

func cacheSummary(m Model, b Backend, s Strategy, f Freshness) string {
	var sb strings.Builder
	switch b {
	case BackendMemory:
		sb.WriteString("A byte-bounded exact LRU inside the application process. Included in the budget rather than treated as free. ")
	case BackendRedis:
		sb.WriteString("A shared Redis with a recorded maxmemory and allkeys-lru policy (approximate LRU, and the report says so). ")
	}
	switch s {
	case StrategyAside:
		sb.WriteString("Cache-aside: a miss takes a per-key fill lease, double-checks, loads a committed view and publishes it. ")
	case StrategyThrough:
		sb.WriteString("Write-through: the new value is never published before the database commit; the writer tombstones, commits, then republishes only the committed representation. ")
	}
	switch f {
	case FreshRelaxed:
		sb.WriteString("A reader may receive an older committed value than the latest acknowledged write, and every such read is counted.")
	case FreshStrict:
		sb.WriteString("A read that begins after a write to that key was acknowledged must not return an older value; a violation fails the cell.")
	}
	if m == ModelOwned {
		sb.WriteString(" The owned model supplies a monotonic cache_version to fence and validate against.")
	} else {
		sb.WriteString(" The legacy model supplies no token, so strict freshness rests on invalidation alone or on an authoritative read.")
	}
	return sb.String()
}

func cacheRisk(s Strategy, f Freshness) string {
	if f == FreshStrict {
		return "MUST record zero stale-after-ack reads and zero impossible cache values; a single wrong read invalidates the cell for performance conclusions"
	}
	return "may return stale committed data; throughput is only reportable with the wrong-read count and rate beside it"
}

// referenceScenarios are the database-layout cells, all measured with no cache.
func referenceScenarios() []Design {
	return []Design{
		{
			ID: "ref-normalized-indexed", Short: "ref-norm", Group: "reference", Model: ModelRef,
			Title:   "normalized, indexed",
			Summary: "Study 01's D2 shape: no copied grandparent key on the child, so the charity feed pays the join and an insert maintains one fewer index.",
			Risk:    "none specific",
			OwnDir:  "reference/normalized_indexed", SQLDir: "reference/normalized_indexed",
			Pair:  "legacy-na-none-none-none-coord",
			Needs: withNeeds(),
		},
		{
			ID: "ref-rollup-trigger", Short: "ref-roll", Group: "reference", Model: ModelRef,
			Title:   "parent rollup maintained by trigger",
			Summary: "The portal's two aggregates become stored person columns maintained by a child trigger that recomputes rather than increments. The portal read stops aggregating; every mutation starts paying.",
			Risk:    "the stored aggregate can drift from the donations it summarises",
			OwnDir:  "reference/rollup_trigger", SQLDir: "reference/rollup_trigger",
			Pair:  "legacy-na-none-none-none-coord",
			Needs: withNeeds(sAuditRollupDrift),
		},
		{
			ID: "ref-embedded-locked", Short: "ref-emb", Group: "reference", Model: ModelRef,
			Title:   "bounded embedding, concurrency-correct (D10)",
			Summary: "The newest 20 donations become a bounded JSONB slice on the parent, maintained by Study 01's D10 trigger: merge-and-sort on insert, parent row lock before a rebuild. Study 01 measured the naive D9 form corrupting donor caches under overlap; only the corrected form is a design here.",
			Risk:    "the embedded slice can drift or be reordered; a_embedded_drift and the Go oracle both check it",
			OwnDir:  "reference/embedded_locked", SQLDir: "reference/embedded_locked",
			Pair:  "legacy-na-none-none-none-coord",
			Needs: withNeeds(sAuditEmbeddedDrift),
		},
		{
			ID: "ref-y1-colocated", Short: "y1-coloc", Group: "reference", Model: ModelRef,
			Title:   "child keyed by person (YugabyteDB)",
			Summary: "The flattened design with the child's primary key chosen for placement: every donation of one donor lives in one tablet. Point lookups by donation id become a two-hop distributed index lookup.",
			Risk:    "donation_id is no longer unique by construction of the key; a unique index enforces it",
			OwnDir:  "reference/y1_colocated", SQLDir: "legacy/d3_legacy",
			Pair: "ref-y2-noncolocated", YBOnly: true,
			Needs: withNeeds(sAuditRolldownDrift),
		},
		{
			ID: "ref-y2-noncolocated", Short: "y2-noncol", Group: "reference", Model: ModelRef,
			Title:   "child keyed by identity (YugabyteDB)",
			Summary: "The same flattened design with the child's natural primary key, so a donor's donations hash across every tablet in the cluster. The pair isolates physical placement with columns, operations, indexes and budget held constant.",
			Risk:    "none specific; this is the reference shape",
			OwnDir:  "reference/y2_noncolocated", SQLDir: "legacy/d3_legacy",
			Pair: "ref-y1-colocated", YBOnly: true,
			Needs: withNeeds(sAuditRolldownDrift),
		},
	}
}

// controlScenarios are measured like designs and never reported as fast results.
func controlScenarios() []Design {
	stale := Design{
		ID: "ctl-stale-invalidation", Short: "ctl-stale", Group: "control", Model: ModelControl,
		Title:     "one suppressed post-commit invalidation (negative control)",
		Summary:   "A relaxed cache-aside cell whose writer commits the mutation and then skips exactly one invalidation. The next read of that key must return the older committed value, and the oracle must see it. Its speed is never shown as a valid result.",
		Risk:      "MUST produce at least one stale-after-ack read; if it does not, no design's freshness is demonstrated",
		ModelUsed: ModelLegacy, Ver: VerNa, Backend: BackendMemory, Strategy: StrategyAside,
		Freshness: FreshRelaxed, Writers: WritersCoord,
		NegativeControl: true, SuppressOneInvalidation: true,
		OwnDir: "legacy/d3_legacy", SQLDir: "legacy/d3_legacy",
		Needs: withNeeds(),
	}
	rmw := Design{
		ID: "ctl-unchecked-rmw", Short: "ctl-rmw", Group: "control", Model: ModelControl,
		Title:     "unguarded read-modify-write on the version token (negative control)",
		Summary:   "The owned no-cache write path with the compare-and-set replaced by a plain read-then-write. Two concurrent writers can both be acknowledged while one bump is lost, so the version token stops describing the committed state -- which is the assumption every strict owned cell rests on.",
		Risk:      "MUST lose version bumps under the hotspot workload; if it does not, the version token's correctness is not demonstrated",
		ModelUsed: ModelOwned, Ver: VerUnsafe, Backend: BackendNone, Strategy: StrategyNone,
		Freshness: FreshNone, Writers: WritersCoord,
		NegativeControl: true,
		OwnDir:          "owned/d3_owned", SQLDir: "owned/d3_owned",
		Needs: withNeeds(ownedExtra...),
	}
	return []Design{stale, rmw}
}

var _ = Design{} // keep the type visible to readers of this file

var designs = func() []Design {
	out := append([]Design{}, matrixScenarios()...)
	out = append(out, referenceScenarios()...)
	out = append(out, controlScenarios()...)
	return out
}()

var designIndex = func() map[string]Design {
	m := make(map[string]Design, len(designs))
	for _, d := range designs {
		if _, dup := m[d.ID]; dup {
			panic("duplicate scenario id " + d.ID)
		}
		m[d.ID] = d
	}
	return m
}()

func designByID(id string) (Design, bool) {
	d, ok := designIndex[id]
	return d, ok
}

// hasCache reports whether the scenario has a cache in front of the database.
func (d Design) hasCache() bool { return d.Backend != BackendNone }

// effectiveModel is the database model the SQL catalogue comes from. The controls
// reuse a model's catalogue but are not that model.
func (d Design) effectiveModel() Model {
	if d.Model == ModelControl {
		return d.ModelUsed
	}
	return d.Model
}

func (d Design) usesVersion() bool { return d.effectiveModel() == ModelOwned }
