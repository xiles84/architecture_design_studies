# Deferred research task: Teach concurrency anomalies, isolation, and control strategies with progressive examples

**Origin:** Owner-requested additions to the Data Architecture Reference backlog, October 17, 2026.
**Disposition:** Proposed only; do not claim or execute until available bandwidth and formal release.

## Research / teaching scope

Build a didactic concurrency field guide inspired by the teaching style of *Designing Data-Intensive Applications*, without copying text. Start with an invariant and two actors and build toward anomalies and remedies. Include SQL before-and-after and a decision map (optimistic vs pessimistic vs constraints / serializable). Keep correctness and throughput separate; compare assumptions, failure modes, and portability across PostgreSQL and YugabyteDB. Propose figures/timelines and end-of-section exercises. Do not fabricate measurements.

## Scientific discipline

Begin with correctness and query semantics. Define measurable contention and predeclared controls; paired measurements require separately authorized experiments. Avoid conclusions or numerical claims from unmeasured designs. Distinguish conceptual figures from actual evidence and apply the project's signed evidence rules.

## Acceptance criteria

- Introduce transaction vocabulary via paired timelines: dirty read, dirty write, lost update, non-repeatable read, read skew, phantom read, write skew, deadlock, starvation; distinguish phenomena from guarantees.
- Explain read committed, repeatable read/snapshot isolation, serializable, MVCC, optimistic validation/versioned CAS and pessimistic locking; call out PostgreSQL and YugabyteDB-specific behavior and SQLSTATE/retry boundaries.
- For each problem, show invariant, a small reproducible two-transaction schedule, how to detect it, two alternative solutions, performance/operational trade-offs, and one reader exercise with answer.
- Discuss first-order lock ordering, isolation scope, transient conflicts, idempotent retries, constraint-enforced invariants, and why dirty writes or reads are not interchangeable with lost updates.
- All diagrams say conceptual or measured; do not invent evidence or label a theoretical guarantee as experimental fact.
