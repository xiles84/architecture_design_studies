# Execution brief — Run the Study 05 v2 churn and scale phases

Implements part of `studies/05-cache-consistency/HANDOFF-AMENDMENT-01-V2.md`
(tag `study-05/v4-handoff`). Read that amendment and the v1 `HANDOFF.md` first.

## Acceptance criteria

- Churn cells cross the 300 s TTL (600 s, 1800 s, 5x burst for 120 s) and record refreshes and expiries actually fired
- Medium scale at 8k, 80k and 800k donors with the working-set-to-RAM ratio and Redis eviction reported
- >=3 trials per cell, or a recorded randomized order with the limitation stated
- Wrong-read rate and throughput reported from the same cell; a regime that crossed no TTL is a gap

## Constraints

- Take the benchmark lock; one matrix at a time; never measure in parallel.
- Correctness gates timing; report wrong-read rate and throughput from the same cell.
- Produce a generated report and a signed analysis; never edit another analyst's file.
- Do not push or pull; commit and tag, and merge through the normal queue lifecycle.

NEXT MODEL: LOW
