# Deferred research task: Design reproducible low/high-contention benchmark arms for current and future studies

**Origin:** Owner-requested additions to the Data Architecture Reference backlog, October 17, 2026.
**Disposition:** Proposed only; do not claim or execute until available bandwidth and formal release.

## Research / teaching scope

Reusable protocol extension for all existing and future relational studies: controlled low/high-contention factors, collision definitions, measurement and correctness gates, retry policy, failure analysis and reproducible paired design. Design only; do not run live experiments in this task.

## Scientific discipline

Begin with correctness and query semantics. Define measurable contention and predeclared controls; paired measurements require separately authorized experiments. Avoid conclusions or numerical claims from unmeasured designs. Distinguish conceptual figures from actual evidence and apply the project's signed evidence rules.

## Acceptance criteria

- Define quantitative low- and high-collision conditions via independent key-popularity distributions (uniform versus Zipf/hot set), concurrency level, transaction duration, workload mixes, and observed conflict rates.
- Specify control of data size, schema, index, hardware, batch size, client arrival rate, retry budget, correctness oracle, warmup and resource placement with equivalent settings across paired arms.
- Prescribe separate read and write throughput, p50/p95/p99 latency, tail amplification, conflict/abort/retry/deadlock/timeout rates, WAL/index churn, lock wait, CPU, memory, and fairness metrics.
- Distinguish logical collisions, physical-hot-row contention, and optimistic validation failures; measure hotspot skew and avoid inference from throughput alone.
- Provide reusable study matrix, randomized paired run order, resource-safety gates, predeclared exclusion/repetition rules, and evidence provenance. No new benchmark without separate scientific authorization.
