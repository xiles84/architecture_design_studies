# Deferred research task: Compare query and write trade-offs for finding charity people who never donated

**Origin:** Owner-requested additions to the Data Architecture Reference backlog, October 17, 2026.
**Disposition:** Proposed only; do not claim or execute until available bandwidth and formal release.

## Research / teaching scope

Primary question: list persons for whom no eligible donation exists. Compare SELECT p.* FROM person p WHERE NOT EXISTS (SELECT 1 FROM donation d WHERE d.person_id=p.id) with LEFT JOIN ... ON and d.person_id IS NULL. Investigate B-tree index donation(person_id), index-only plans where supported, anti-join plans and selectivity; add maintained metadata on person (has_donated/first or last donation date), separate donor marker and aggregate read model. Include the consistency cost of correcting the flag on deletes/refunds. Analyze how inserts to a child table impact hotspot parents and indexes, including contention matrix.

## Scientific discipline

Begin with correctness and query semantics. Define measurable contention and predeclared controls; paired measurements require separately authorized experiments. Avoid conclusions or numerical claims from unmeasured designs. Distinguish conceptual figures from actual evidence and apply the project's signed evidence rules.

## Acceptance criteria

- Define never donated using the committed donation lifecycle and nullable/soft-deleted/refunded donation semantics; never confuse no qualifying donation with amount = 0.
- Compare indexed NOT EXISTS anti-semi-join, LEFT JOIN IS NULL anti-join, aggregate HAVING, and schema alternatives (person.has_donated, first/last donation timestamp, separate donor marker/summary).
- For every arm document index design (especially donation.person_id; consider covering and partial indexes), EXPLAIN plans, scan cardinality and correctness on skewed populations.
- Include read-only, write-heavy, mixed, first-donation race, donation deletion/rollback and low/high-contention workloads; evaluate read and write throughput, latency, amplification and maintaining derived state atomically.
- Provide benchmark protocol and book-ready SQL walkthroughs, but defer results until approved runs and verified evidence.
