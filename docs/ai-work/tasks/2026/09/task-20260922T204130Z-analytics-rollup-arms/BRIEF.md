# Execution brief — Run the Study 08 rollup and materialized-view arms

Implements part of `studies/08-analytics-read-models/HANDOFF.md`. Read that handoff first (and the study README).

## Acceptance criteria

- Materialized-view refresh schedules with refresh duration and interval actually crossed
- Top-N dashboard workloads measured with query gain, staleness, storage and correctness
- Refresh-boundary behaviour defined and recorded; stale-copy control fires
- Within-run comparison against the base only

## Constraints

- Take the benchmark lock; one matrix at a time; never measure in parallel.
- Correctness gates timing; record the environment, sizing and per-endpoint distribution.
- Produce a generated report (numbers only) and a signed analysis; never edit another analyst's file.
- Do not push or pull; commit and tag, and merge through the normal queue lifecycle.

NEXT MODEL: LOW
