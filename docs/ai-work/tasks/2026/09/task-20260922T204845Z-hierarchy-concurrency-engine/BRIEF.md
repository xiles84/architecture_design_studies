# Execution brief — Run the Study 09 concurrency and second-engine arms

Implements part of `studies/09-hierarchy/HANDOFF.md`. Read that handoff first (and the study README).

## Acceptance criteria

- Contended concurrent moves of one subtree and a read-during-move, with a lost-update control
- PostgreSQL and YugabyteDB both run the same correctness gate and operation set
- Engines reported separately; no cross-engine ranking pooled with the representation comparison
- Per-node and equal-total arms labelled and separate

## Constraints

- Take the benchmark lock; one matrix at a time; never measure in parallel.
- Correctness gates timing; record the environment, sizing and per-endpoint distribution.
- Produce a generated report (numbers only) and a signed analysis; never edit another analyst's file.
- Do not push or pull; commit and tag, and merge through the normal queue lifecycle.

NEXT MODEL: LOW
