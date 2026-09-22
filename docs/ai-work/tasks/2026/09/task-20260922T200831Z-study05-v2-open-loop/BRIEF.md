# Execution brief — Run the Study 05 v2 open-loop demand phase

Implements part of `studies/05-cache-consistency/HANDOFF-AMENDMENT-01-V2.md`
(tag `study-05/v4-handoff`). Read that amendment and the v1 `HANDOFF.md` first.

## Acceptance criteria

- Open-loop arrival at two fixed rates with the platform arrival driver; no closed-loop backpressure
- Report achieved throughput, non-attempted operations and latency percentiles per offered rate
- State the offered rate and whether the server saturated; closed-loop numbers are not reused as open-loop

## Constraints

- Take the benchmark lock; one matrix at a time; never measure in parallel.
- Correctness gates timing; report wrong-read rate and throughput from the same cell.
- Produce a generated report and a signed analysis; never edit another analyst's file.
- Do not push or pull; commit and tag, and merge through the normal queue lifecycle.

NEXT MODEL: LOW
