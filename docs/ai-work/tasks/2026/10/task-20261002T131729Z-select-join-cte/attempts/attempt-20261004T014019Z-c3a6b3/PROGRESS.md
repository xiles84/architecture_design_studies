# Progress — task-20261002T131729Z-select-join-cte

Attempt `attempt-20261004T014019Z-c3a6b3`, branch `study-01/select-join-cte`,
LOW executor (deepseek-flash / opencode CLI). Milestone tag:
`study-01/select-v2-join-cte`.

**Status: recon and premeasurement plan only; no run has been started.** This
checkpoint was committed deliberately: the join/CTE matrix is the largest of the
four SELECT families (~150–216 cells per whole run) and needs two matrix
recalibrations and a fixture-feasibility calibration before measurement. The
worktree and branch are left ready for the next LOW session to resume.

## Recon findings

- The shared runner is already on `main` (integrated with the access-paths and
  aggregate/top-N families): `CellSpec`/`RunSpecs`/`Axis`, eight live controls on
  the teaching fixture. This task can build on it directly.
- Existing join/filter statements: `s01_join_where`, `s02_filtered_derived`,
  `s03_cte_filtered`, `s04_join_on`, `s05_comma_join` (five of the six spellings).
  The sixth (CTE of joined rows, filtered outside) and the
  `MATERIALIZED`/`NOT MATERIALIZED` variants still have to be added, as does a
  reuse-N query form.
- The teaching fixture is far too small for a 0.01% selectivity probe: a new
  fact-table scale is required before the matrix can run.

## Decision Log

1. **Stop at a committed plan rather than start the matrix.** The matrix is
   ~150–216 cells per run at 5 trials (~1,000 queries/run), needs a new large
   fact fixture, and needs feasibility calibration of the 0.01% level. Starting
   it with insufficient remaining budget would leave a half-built harness.
   Confidence High. Alternative: run a reduced matrix now (would not satisfy the
   accepted matrix and would produce a report that cannot answer P1–P4).
2. **Reuse levels must be extended to ≥5** ({1,2,4,8,16}) or P4 cannot bound a
   crossing; the reviewed matrix's reuse {1,8} has two levels. Flagged for HIGH as
   an expansion, exactly as the earlier families' regime/level expansions were.
3. **Materialization variants apply only to CTE spellings.** `MATERIALIZED` is
   not valid syntax on a plain join or a derived table; those spellings
   contribute one cell per level. Flagged for HIGH.
4. **A CTE-materialization control must verify the CTE is really materialized**
   (plan evidence) or mark the cell `unsupported`, so a default fold is never
   silently counted as a `MATERIALIZED` run.
5. **Proposed fixture: ~200,000 fact rows** with an indexed selective column,
   loaded once per whole run; the 0.01% level must be shown to select at least
   one row or it is dropped and recorded.

## Plan for the resuming session

1. Recover/claim this task and its worktree; read this file and the bounded
   matrix.
2. Add the sixth spelling, the three materialization variants for each CTE
   spelling, and the reuse-N query form; add the fixture scale and calibrate the
   0.01% level and load feasibility.
3. Add the materialization-evidence control and the LEFT-JOIN contract controls.
4. Pilot one level; then run ≥3 whole runs per engine with 5 trials and a
   position control; write the reports and the signed analysis with P1–P4
   verdicts; RESULT.md; submit for review.
