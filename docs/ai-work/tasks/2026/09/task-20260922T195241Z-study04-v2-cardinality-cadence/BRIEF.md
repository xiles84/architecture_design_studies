# Execution brief — Run the Study 04 v2 cardinality and cadence phases

Implements part of `studies/04-configuration-portal/HANDOFF-AMENDMENT-01-V2.md`
(tag `study-04/v2-handoff`). Read that amendment and the v1 `HANDOFF.md` first.

## Acceptance criteria

- Cardinality tiers 1, 10, 30, 60, 120, 500 entries for n0_rows_unindexed, n1_rows_indexed, n3_rollup_trigger, n4_rollup_app
- Cadence regimes quiet, steady, burst at the 120-entry tier; staleness and rollup lag reported for trigger vs application maintenance
- >=3 independent trials per cell, or a recorded randomized order with the limitation stated
- State the refresh/expiry actually crossed in each regime; a regime that crossed none is a gap, not a null result

## Constraints

- Take the benchmark lock; one matrix at a time; never measure in parallel.
- Correctness gates timing: a design that fails verification aborts its cell.
- Produce a generated report (numbers only) and a signed analysis; never edit another analyst's file.
- Do not push or pull; commit and tag, and merge through the normal queue lifecycle.

NEXT MODEL: LOW
