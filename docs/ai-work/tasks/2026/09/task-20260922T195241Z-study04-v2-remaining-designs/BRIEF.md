# Execution brief — Run the Study 04 v2 remaining designs, equal-total arm and YugabyteDB coverage

Implements part of `studies/04-configuration-portal/HANDOFF-AMENDMENT-01-V2.md`
(tag `study-04/v2-handoff`). Read that amendment and the v1 `HANDOFF.md` first.

## Acceptance criteria

- Build and run the eight unbuilt v1 designs, or record each as a coverage gap with the reason
- Separately labelled equal-total-resource arm alongside the per-node baseline
- YugabyteDB 1-node RF=1 and 3-node RF=3 for the cardinality and rollup designs, reported per engine
- Placement designs y1/y2 executed with the actual tablet/leader distribution recorded; no placement record means gap

## Constraints

- Take the benchmark lock; one matrix at a time; never measure in parallel.
- Correctness gates timing: a design that fails verification aborts its cell.
- Produce a generated report (numbers only) and a signed analysis; never edit another analyst's file.
- Do not push or pull; commit and tag, and merge through the normal queue lifecycle.

NEXT MODEL: LOW
