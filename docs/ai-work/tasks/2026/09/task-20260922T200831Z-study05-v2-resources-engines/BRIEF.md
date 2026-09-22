# Execution brief — Run the Study 05 v2 equal-total accounting and YugabyteDB/placement phases

Implements part of `studies/05-cache-consistency/HANDOFF-AMENDMENT-01-V2.md`
(tag `study-05/v4-handoff`). Read that amendment and the v1 `HANDOFF.md` first.

## Acceptance criteria

- Separately labelled equal-total arm charging a share to the Redis container; never pooled with db-only
- Strict/relaxed pair repeated on YugabyteDB 1-node RF=1 and 3-node RF=3, reported per engine
- Colocated vs non-colocated pair executed with the tablet/leader distribution recorded; no record means gap
- Correctness gate passes; negative controls fired; one matrix at a time under the benchmark lock

## Constraints

- Take the benchmark lock; one matrix at a time; never measure in parallel.
- Correctness gates timing; report wrong-read rate and throughput from the same cell.
- Produce a generated report and a signed analysis; never edit another analyst's file.
- Do not push or pull; commit and tag, and merge through the normal queue lifecycle.

NEXT MODEL: LOW
