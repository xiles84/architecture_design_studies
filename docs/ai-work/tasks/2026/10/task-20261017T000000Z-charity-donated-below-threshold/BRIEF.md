# Deferred research task: Compare approaches for lifetime charity donations below X including non-donors

**Origin:** Owner-requested additions to the Data Architecture Reference backlog, October 17, 2026.
**Disposition:** Proposed only; do not claim or execute until available bandwidth and formal release.

## Research / teaching scope

Clarify that 'donated less than X' means lifetime SUM(amount) < X by default, not 'at least one donation < X'. Include never-donors by preserving outer person rows. Example HAVING COALESCE(SUM(d.amount),0) < :x using LEFT JOIN. Compare row-level cached totals, separate per-person total table, live aggregation, asynchronous derived view and hybrid cache; attribute maintenance costs to writes. Examine refunds and currencies, and test low/high collision, read/write mix, skew and thresholds.

## Scientific discipline

Begin with correctness and query semantics. Define measurable contention and predeclared controls; paired measurements require separately authorized experiments. Avoid conclusions or numerical claims from unmeasured designs. Distinguish conceptual figures from actual evidence and apply the project's signed evidence rules.

## Acceptance criteria

- Define the metric as total eligible lifetime donated amount strictly less than X, including zero-donation people; make X, currency/precision, refunds, cancellation and negative totals explicit.
- Compare LEFT JOIN + GROUP BY + HAVING COALESCE(SUM(amount),0) < X, pre-aggregated subqueries, correlated aggregate, and persisted total on person or separate per-person aggregate table/materialized view.
- Consider index strategies for donation(person_id, amount), aggregation coverage, thresholds and population selectivity; compare EXPLAIN plans and data size sensitivity.
- Measure correct read and write throughput (first donation, additional donations, refunds/deletes), contention, atomic aggregate maintenance and staleness trade-offs under low/high collision.
- Propose two-transaction anomaly schedules where concurrent child writes break naive derived totals; demand invariant/correctness tests and signed evidence before claims.
