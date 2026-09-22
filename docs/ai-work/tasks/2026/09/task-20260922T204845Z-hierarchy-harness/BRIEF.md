# Execution brief — Build the Study 09 hierarchy harness and correctness gate

Implements part of `studies/09-hierarchy/HANDOFF.md`. Read that handoff first (and the study README).

## Acceptance criteria

- Deterministic tree generator at 1k/100k/1M nodes, balanced and skewed, seed 42
- Loaders and schemas for adjacency list, materialized path, closure table, nested sets and bounded embedding
- Cross-representation answer checker plus cycle, orphan and depth-bound controls that are seen to fire
- Same statement names and result shapes across representations

## Constraints

- Take the benchmark lock; one matrix at a time; never measure in parallel.
- Correctness gates timing; record the environment, sizing and per-endpoint distribution.
- Produce a generated report (numbers only) and a signed analysis; never edit another analyst's file.
- Do not push or pull; commit and tag, and merge through the normal queue lifecycle.

NEXT MODEL: LOW
