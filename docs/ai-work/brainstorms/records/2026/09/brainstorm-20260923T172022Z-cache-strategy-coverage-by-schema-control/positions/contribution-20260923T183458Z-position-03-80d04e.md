# Independent position

| Field | Value |
|---|---|
| Brainstorm | `brainstorm-20260923T172022Z-cache-strategy-coverage-by-schema-control` |
| Contribution | `contribution-20260923T183458Z-position-03-80d04e` |
| Slot | `position-03` |
| Actor | model=unknown tool=Codex effort=unknown session=01a0cf34-bce0-7162-843a-e5ec65412925-cache-third capability=HIGH role=analyst |
| Capability input | HIGH |
| Submitted | `2026-09-23T18:34:58Z` |
| Confidence | medium |

## Summary

Map cache choices by source schema and writer control, then cache topology. Versions order known publications but do not alone prove currentness; Study 05 supports scoped warm-read and multi-instance findings.

## Disagreements

Whether versions and tombstones can avoid shared coordination, and how newer unregistered v2 evidence should enter the book.

## Contribution

# Independent HIGH position — cache strategy coverage by schema control

**Brainstorm:** `brainstorm-20260923T172022Z-cache-strategy-coverage-by-schema-control`  
**Slot:** `position-03`  
**Author:** GPT-6-family Codex desktop, HIGH session; exact model variant and effort were not exposed. The live claim stores `model: unknown` for that reason. This is a separate session of the same model family as Codex position-02, not a third distinct model.  
**Confidence:** high on the logical guarantee boundaries; medium on the proposed book organization; low on comparative costs of tactics Study 05 has not isolated.

## Evidence inspected

I read the question and its listed packet: `book/concepts/cache-consistency.typ`, `book/chapters/{major-families,minor-variants,scenarios}.typ`, `book/evidence/v2/claims.json`, Study 05 `README.md`, `HANDOFF-AMENDMENT-01-V2.md`, the four current signed analyses and the partially superseded fence discussion under Study 05 `reports/`, `sql/README.md`, `docs/methodology.md`, and `CONTEXT.md`. I also read the repository brainstorm workflow and schema. I did not read either earlier brainstorm position. No new measurement or task was created.

## Recommendation

Teach caches as **read copies with an explicit acknowledgement-relative freshness contract**, then map strategy on two independent axes: (1) who can change the source and whether the schema can carry a transactional version/outbox, and (2) whether the cache is private to a process or shared. Cache-aside and write-through are *fill/publication paths*, not consistency guarantees. A 300-second TTL, a fill lease, or a soft delete does not turn either path into strict freshness.

The book's six-family taxonomy can remain. Expand the external-read-copies chapter and cache-consistency concept into a decision table, with each row labelled `measured`, `mechanism demonstrated`, or `proposed/unmeasured`. Keep the v1 book's numbers tied to the 23 active v2 registry claims. The 2026-09-23 Study 05 analyses are newer evidence but are not yet active v2 claims; a later book edition should reconcile and validate them before quoting new numbers or repeating a now-stale blanket gap statement.

| Source/writers | Process-local cache | Shared cache | Freshness boundary and cost |
|---|---|---|---|
| **Frozen schema; all writers through one adapter** | Cache-aside or post-commit write-through can use the adapter's per-key fence and bounded fill lease. A single process has one coordination domain. | Cache-aside or write-through needs an atomic *shared* fence/conditional publication, and every instance must participate. Private leases cannot coordinate the shared key. | Study 05 demonstrates the fence mechanism and a one-process strict contract under its workload. A crash between DB commit and cache coordination still needs an explicit fail-closed/recovery policy before claiming strictness through failures. |
| **Frozen schema; writers may bypass adapter** | Local invalidation misses external commits. | Shared invalidation also misses them. | Strict reads require an authoritative read or a complete, timely change source with a validated delivery/acknowledgement protocol; otherwise state a weaker freshness bound. Study 05's external-writer arm falls back toward authoritative reads and loses the cache gain. CDC/log capture is a candidate, not measured here. |
| **Owned schema; all writers obey protocol** | A version stored with the row allows a cache hit to be checked against DB authority; this may cost a DB lookup per hit. | Transactional version/outbox can give durable change identity, replay and cross-instance invalidation. Cache publication still needs atomic compare/fence and an ordered consumer. | The owned model implements version bump and outbox, but the measured legacy/owned contrast is a *package*, not the isolated price or benefit of a version column. An async outbox alone does not make a post-ack cache hit current. |
| **Owned schema; writers bypass the protocol** | A version or outbox omitted by any writer is no longer a valid completeness proof. | The same problem crosses instances. | Ownership means the protocol can be installed on all sanctioned paths; it is not a guarantee until coverage, migration, bulk writes, and failure paths are audited. |

The first two rows are useful even when schema control is absent. They also make a necessary distinction: **strictness is still possible with a frozen schema by reading the database on each relevant read**; what becomes unavailable without complete write observation is a guaranteed *cache-hit* shortcut. Claiming that legacy strictness is impossible would overstate the evidence.

## What the evidence supports today

- The active registry's `v2-14` supports a **within-run add-cache result**: 2.2–3.5× warm cacheable-read throughput and p99 moving from 15.8–27.1 ms to 0.8–2.0 ms, at 800 donors, one PostgreSQL host, closed loop and `db-only` resource framing. Redis is a cache, never an authoritative key-value database in this study. It does not establish a universal cache preference or equal-total efficiency.
- `v2-15` supports a **specific multi-instance failure mechanism**: the three-instance, shared-Redis, write-through, relaxed cells returned 87.24% wrong reads for legacy and 87.81% for owned, while the single-instance relaxed cells observed zero. The independent review corrected a swap of those model labels in the original analysis. The rate is one trial on one host, not a forecast for another service. The owned version did not cure this path by itself.
- `v2-16` supports a narrow negative result: strict/relaxed read-throughput differences of −14.2% to +12.3% lay inside this run's approximate 20% noise floor. It does **not** mean strictness is free. Strict writers use more cache fences, and costs/failure behavior belong on the write and coordination path.
- The fence discussion gives the decisive race: a reader snapshots committed state S0, a writer commits S1 and deletes the key, then the reader republishes S0 after deletion. `publish after commit` and simple invalidation do not block this. A per-key generation captured **before** the snapshot, checked atomically at publication, plus the writer's before/after-commit fences prevents that stale refill for coordinated writes. The all-green analysis reports refused stale publications and zero strict stale-after-ack reads in its measured cells; both negative controls fired and the independent ledger found no impossible values. These are source/test results for that protocol, not a proof for every crash or writer topology.
- The newer signed churn analysis closes only the v1 claim that *wall time never passed* 300 seconds: measured strict owned cells ran through 2 and 6 TTL intervals with zero wrong reads. All observed expiry was probabilistic early expiry; **hard expiry fired zero times**, refresh counts were absent, and the burst and requested larger scales were not run. The newer resource/engine analysis labels an equal-total arm, runs YugabyteDB strict/relaxed and intended placement pairs with three endpoints, but the tablet/leader distribution is still missing and all nodes share one host. The pair directions are single-run results. Thus the v1 book's blanket list of unmeasured regimes has become partly stale, while hard-expiry, scale, burst, placement proof, open-loop demand, repeated multi-instance behavior and failure coverage remain gaps. The active registry still describes the v1 scope; do not silently promote the newer analyses into its claim cards.
- Study 05 includes both cache-aside and write-through and both process-local and Redis backends, but the book itself classifies aside versus through as **not a settled controlled variant recommendation**. The matrix is coverage of implementations, not evidence that one strategy is generally faster or safer. Memory and Redis also differ in capacity policy (exact versus approximate LRU) and network/failure domain.

## Which "avoid locks" claim is defensible?

First name the lock. A **database lock** serializes a conflicting business mutation; a **fill lease** suppresses a cache stampede; an **atomic cache fence/CAS** orders publication against invalidation. These solve different races. Owned row versions can support optimistic conditional updates and may replace an *application-held pessimistic lock* for some mutations, at the price of retries. They do not remove the database's internal arbitration or prove the business invariant under contention. Study 05 deliberately has both optimistic and pessimistic owned write paths; its one-trial pair is not a general lock-cost ranking.

A monotonic version and `put-if-newer` can stop an older filler from overwriting a newer **known** cache value. It cannot guarantee that the known cache value is the latest database value. After DB commit v8 but before any cache event arrives, a cache still holding v7 can pass every local version comparison and serve a stale hit. An outbox makes a committed change record durable, but an asynchronous consumer still has delivery lag and can replay events out of order. To promise no stale-after-ack hits, the writer must wait until the relevant invalidation/publication is effective, or the read must validate against authority (or a proven equivalent barrier); failures need a fail-closed policy. These mechanisms have different DB-read, write-latency and operational costs and were not compared in Study 05 as isolated pairs.

Likewise, distinguish **soft delete in the database** from a **cache tombstone**. A DB soft delete preserves a row and may preserve a monotonic version; a cache tombstone tells cache readers/fillers that a key is deleted. Neither alone prevents resurrection. Example: R reads live v7; W commits deleted v8 and invalidates; R publishes v7 afterward. A cache-side generation/CAS must reject R. If v8's tombstone expires or is evicted before R publishes, a plain version comparison can lose its reference state. Key reuse or version reset causes an ABA risk. A defensible lock-avoidance design requires a monotonic generation across delete/recreate, atomic comparisons, retained/authoritative delete knowledge for the maximum in-flight fill window, and all writers participating. Its exact retention, outage and cost behavior is **not measured** here.

A cache fill lease can still be useful even with versions: it limits duplicate DB work when a hot key expires. Study 05 measured one load per key in an eight-key stampede, but a lease must have a unique token, bounded lifetime/wait and compare-token release; an expired holder still needs a publication fence. Calling a versioned design "lock-free" without distinguishing these roles would mislead readers.

## Alternatives and anticipated disagreements

1. **Treat owned version/outbox as the preferred default and omit the legacy path.** This yields a simpler narrative but would erase the frozen-schema constraint and the external-writer limit that Study 05 was designed to expose. Keep owned as a conditional upgrade, not a universal prerequisite.
2. **Lead with cache-aside versus write-through.** They are familiar names and both are implemented, but they do not decide whether an old fill can republish, whether another process sees invalidation, or whether a writer can bypass the service. Lead with contract and writer coverage; discuss the two write/fill paths under each cell of the map.
3. **Use short TTL as the legacy safety answer.** TTL can bound staleness only under explicit clock, expiry and reachability assumptions; it cannot meet zero stale-after-ack and the hard-expiry mechanism was not exercised in the new churn runs. Label this as a relaxed contract with an observable wrong-read budget, not strictness.
4. **Say versions and soft deletes eliminate locks.** They can reduce particular application mutex uses when every operation is versioned and cache operations are atomic. The stale-hit window, delete resurrection, expired lease holder, cache outage, and business write contention remain. Replace that slogan with the named guarantee and the conditions above.
5. **Promote all 2026-09-23 results to the book immediately.** They close parts of the old gap, but the active registry's claim cards still point at the v1 digest. The book's numeric rule requires a versioned, semantically checked evidence update before new figures appear in prose.

## Assumptions and evidence that could falsify this position

I assume the desired strict contract is Study 05's: a read starting after an acknowledged write to the same key must not return older committed data; overlapping reads are classified separately; impossible/uncommitted cache values are forbidden in every arm. I assume the cache is a read copy, and all tabulated DB models keep PostgreSQL/YugabyteDB authoritative. A weaker bounded-staleness product contract would change the recommended tactic, not the race analysis.

The strongest falsification tests for a future governed study would be: (a) pause an aside filler between DB snapshot and cache publish while a writer commits/invalidate, with fence/CAS disabled and enabled; (b) crash a coordinated writer after commit but before cache fence/ack, restart and read from another process; (c) delay, duplicate and reorder owned outbox events, then read after writer ack; (d) delete and recreate the same key after tombstone expiry while an old filler or expired lease holder publishes; (e) let an external legacy writer mutate a hot key without going through the adapter and compare cache-hit answers with an independent DB/operation oracle; and (f) compare version-check-on-hit, synchronous shared fence, and async outbox under equal-total resources, offered load, cache outage and actual write latency. Each must report DB reads per logical read, wrong/impossible reads, stale duration, writes/retries, contention, refresh amplification, bytes and operational dependencies, with negative controls that fire. A design that keeps strictness through (b)–(e) while retaining cache-hit shortcuts without complete write observation would directly challenge my guarantee boundary.

**Position:** preserve the family taxonomy, make the schema/writer-coverage by cache-topology map the organizing table, and qualify every lock-avoidance statement by the exact race it solves. No tasks or measurements follow from this contribution alone.
