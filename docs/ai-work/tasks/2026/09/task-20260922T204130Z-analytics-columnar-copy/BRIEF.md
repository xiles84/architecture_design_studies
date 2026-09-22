# Execution brief — Run the Study 08 analytical (columnar) copy arm

Implements part of `studies/08-analytics-read-models/HANDOFF.md`. Read that handoff first (and the study README).

## Acceptance criteria

- ClickHouse copy loaded from the base at a recorded cadence; loader cadence and lag reported
- Top-N dashboard workloads measured with query gain, staleness, storage and correctness
- Per-node and equal-total resource arms labelled and separate
- Within-run comparison against the base only

## Constraints

- Take the benchmark lock; one matrix at a time; never measure in parallel.
- Correctness gates timing; record the environment, sizing and per-endpoint distribution.
- Produce a generated report (numbers only) and a signed analysis; never edit another analyst's file.
- Do not push or pull; commit and tag, and merge through the normal queue lifecycle.

NEXT MODEL: LOW
