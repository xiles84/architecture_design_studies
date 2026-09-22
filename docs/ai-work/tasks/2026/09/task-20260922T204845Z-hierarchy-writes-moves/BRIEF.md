# Execution brief — Run the Study 09 insert, delete and move axes

Implements part of `studies/09-hierarchy/HANDOFF.md`. Read that handoff first (and the study README).

## Acceptance criteria

- Move-subtree cost reported as a function of subtree size and depth, not one number
- Insert/delete leaf and subtree, batched and single, measured
- Cycle and orphan controls fire; move results separated by representation
- Within-run comparison only

## Constraints

- Take the benchmark lock; one matrix at a time; never measure in parallel.
- Correctness gates timing; record the environment, sizing and per-endpoint distribution.
- Produce a generated report (numbers only) and a signed analysis; never edit another analyst's file.
- Do not push or pull; commit and tag, and merge through the normal queue lifecycle.

NEXT MODEL: LOW
