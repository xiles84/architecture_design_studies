# Execution brief — Build the Study 08 read-model harness, pins and correctness gate

Implements part of `studies/08-analytics-read-models/HANDOFF.md`. Read that handoff first (and the study README).

## Acceptance criteria

- ClickHouse Containerfile and versions.env pin with recorded digest; a moved tag blocks the run
- Deterministic loader for the base and each copy, with duplicate/lost-row controls
- Refresh, staleness and storage instrumentation shared by every strategy
- Correctness gate returns the base answer under the stated freshness contract; controls fire

## Constraints

- Take the benchmark lock; one matrix at a time; never measure in parallel.
- Correctness gates timing; record the environment, sizing and per-endpoint distribution.
- Produce a generated report (numbers only) and a signed analysis; never edit another analyst's file.
- Do not push or pull; commit and tag, and merge through the normal queue lifecycle.

NEXT MODEL: LOW
