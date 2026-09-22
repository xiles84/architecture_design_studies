# Execution brief — Run the Study 04 v2 instrument and contention controls

Implements part of `studies/04-configuration-portal/HANDOFF-AMENDMENT-01-V2.md`
(tag `study-04/v2-handoff`). Read that amendment and the v1 `HANDOFF.md` first.

## Acceptance criteria

- Control A: n1_rows_indexed first and last, >=3 repeats, fresh database per cell or recorded randomized order; report the identical-SQL read floor as a measured range
- Control B: c1_optimistic_version with retries disabled under 16 writers and one key; only then may 1.76x be restated as an upper bound
- Control C: c1_optimistic_version vs c2_pessimistic_lock at 1, 4, 16 and 64 writers; report throughput and abort/retry rate per point
- Correctness gate passes for every cell; negative controls fired; one matrix at a time under the benchmark lock

## Constraints

- Take the benchmark lock; one matrix at a time; never measure in parallel.
- Correctness gates timing: a design that fails verification aborts its cell.
- Produce a generated report (numbers only) and a signed analysis; never edit another analyst's file.
- Do not push or pull; commit and tag, and merge through the normal queue lifecycle.

NEXT MODEL: LOW
