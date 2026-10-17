# Deferred research task: Benchmark parent row lock versus optimistic CAS update for OneToMany invariants

**Origin:** Owner-requested additions to the Data Architecture Reference backlog, October 17, 2026.
**Disposition:** Proposed only; do not claim or execute until available bandwidth and formal release.

## Research / teaching scope

FOCUS: inserting new children in a OneToMany relation, the hard case for optimistic CAS because the child row does not exist yet. Scenario: parent has finite capacity and transactions concurrently attempt to INSERT new child rows. Arm A: lock parent SELECT FOR UPDATE, read active children/capacity, insert child in one transaction. Arm B: read parent version and child predicate, perform successful UPDATE parent SET lock_version = lock_version + 1 WHERE version = expected, then insert child within the SAME transaction; on failed CAS abort/retry the WHOLE transaction. Arm C: atomic conditional parent UPDATE that reserves capacity then INSERT child in the same transaction. Test first child, same-parent concurrent siblings, nearly full parent, duplicate child ID/business key, child-first versus CAS-first ordering, INSERT errors after successful parent CAS, rollback and idempotent retry. Test competing writers that do not update parent version as a correctness-negative control, demonstrating why a bare parent version CAS cannot defend an uncoordinated child predicate. Contrast DB UNIQUE/exclusion, slot allocation, SERIALIZABLE and locks where equivalent invariant is expressible. Compare physical parent lock versus parent write amplification, MVCC/WAL, foreign-key side effects, lock queues versus abort retries and correctness across LOW/HIGH collision (hot keys and uniformly distributed keys), capturing child insert throughput, latency and failure rates. Distinguish design/theory from measured PostgreSQL/YugabyteDB evidence; do not claim a universal winner.

## Scientific discipline

Begin with correctness and query semantics. Define measurable contention and predeclared controls; paired measurements require separately authorized experiments. Avoid conclusions or numerical claims from unmeasured designs. Distinguish conceptual figures from actual evidence and apply the project's signed evidence rules.

## Acceptance criteria

- Make INSERT INTO child the PRIMARY use case, not an incidental variation: test creating the first child, inserting concurrent siblings for the same parent, capacity-bound inserts, duplicate logical child creation, and failed/rolled-back inserts.
- Choose precise shared invariants spanning a parent and newly created children (e.g. max active children, reservation capacity, or unique child allocation); ensure all arms enforce exactly the SAME invariant atomically even when two children do not exist at the initial read.
- Compare (A) SELECT parent ... FOR UPDATE then check capacity/child predicate then INSERT child, (B) read version/predicate then CAS UPDATE parent version/counter and INSERT child in the SAME transaction, and (C) atomic conditional UPDATE parent SET reserved=reserved+1 WHERE reserved<capacity RETURNING ... followed by INSERT child.
- Include child-first then parent-CAS vs parent-CAS-first then child INSERT; explain rollback, visibility, FK/unique-index interactions, deadlocks, conflict detection and whether wasted insert work or parent writes are incurred.
- Explain explicitly why a CAS on a non-existent child row cannot protect concurrent child INSERTs, when a parent CAS acts as a serialization point, and why optimistic CAS loses correctness if any competing child writer bypasses the agreed protocol.
- Compare additional valid alternatives: unique/exclusion constraints or slot-allocation tables where expressible, SERIALIZABLE transactions, parent advisory locks, and database-enforced constraints; never assume one mechanism alone protects arbitrary COUNT(child) predicates.
- Distinguish SELECT FOR UPDATE locking from UPDATE: parent tuple/WAL cost, version changes, indexes, bloat, FK lock footprints, round trips, retries, aborted child inserts and contention on the hot parent row; distinguish PostgreSQL and YugabyteDB behavior.
- Benchmark both LOW and HIGH collision, parent-key distribution (uniform/hotspot), children per parent, 1 vs N concurrent inserts, success/failure ratios and read/write mix. Capture correctness anomalies, inserted-child throughput, p50/p95/p99 latency, lock wait, abort/retry, WAL and CPU.
- Provide two-transaction SQL schedules for same-parent concurrent inserts, near-capacity races and duplicate insert races, plus transactional rollback/retry/idempotency tests; include book-ready explanations and evidence provenance without inventing results.
