# Result — Study 08 analytics/read-model protocol

**Task:** `task-20260922T025912Z-analytics-read-model-protocol`
**Capability/role:** HIGH planner (model `deepseek-flash`)
**Required tag:** `repo/analytics-read-model-handoff-v1`
**Measured:** nothing.

## Delivered

- `studies/08-analytics-read-models/HANDOFF.md` and `README.md`: a decision-complete protocol
  comparing base OLTP, rollup/materialized views, an in-engine search copy and a ClickHouse
  analytical copy for top-N/dashboard workloads, measuring query gain, refresh/maintenance,
  staleness, storage and correctness with negative controls.
- Four separately claimable execution tasks published as `proposed`:
  `task-20260922T204130Z-analytics-harness`, `…-rollup-arms`, `…-search-copy`, `…-columnar-copy`.
- `CONTEXT.md` open gap 8 records completion and the new tasks.

## Acceptance criteria check

| Criterion | Result |
|---|---|
| Compares base OLTP, rollups/materialized views, search copies and analytical copies for top-N/dashboard workloads | met (handoff §2–3) |
| Measures refresh/maintenance, staleness, storage, correctness and query gains | met (§4) |
| Creates execution tasks but does not measure | met (four published) |

## Notes

The book's rollup claims come from an in-engine materialized/trigger design; this study is what will
let the book speak about read-model *copies* generally.
