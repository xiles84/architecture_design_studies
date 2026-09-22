# Execution brief — Run the Study 09 read and storage axes

Implements part of `studies/09-hierarchy/HANDOFF.md`. Read that handoff first (and the study README).

## Acceptance criteria

- Every read operation measured on balanced and skewed trees with cross-representation equality
- Storage (rows, indexes, total) reported per representation at three tree sizes
- Within-run comparison only; per-node and equal-total arms labelled
- No failed cell becomes a winner

## Constraints

- Take the benchmark lock; one matrix at a time; never measure in parallel.
- Correctness gates timing; record the environment, sizing and per-endpoint distribution.
- Produce a generated report (numbers only) and a signed analysis; never edit another analyst's file.
- Do not push or pull; commit and tag, and merge through the normal queue lifecycle.

NEXT MODEL: LOW
