# Bounded measurement plan — join/filter spellings and CTE materialization (`select-join-cte`)

LOW worker plan for the reviewed protocol's task-3 row
(`protocols/select-query-templates.md`, P1–P4). Committed before any run. This
record is a checkpoint: the matrix below is the reviewed one and has not yet been
executed.

## 1. Question and predictions

- **P1** six inner-join spellings pairwise: no resolved difference in most
  regimes; a difference appears only where a CTE materializes or an index cannot
  be used. Falsifier: a pair differs beyond spread, same direction, ≥3 runs at
  ≥3 trials, with an identified plan difference.
- **P2** derived table vs plain join: no resolved difference.
- **P3** default CTE vs `MATERIALIZED` under a selective one-use predicate:
  `MATERIALIZED` is slower; reverses under high reuse.
- **P4** `MATERIALIZED` vs `NOT MATERIALIZED` under repeated reuse:
  `MATERIALIZED` is faster above a reuse threshold; needs ≥5 reuse levels to
  bound a crossing.

## 2. Matrix (reviewed) and the two recalibrations it needs

Reviewed bounded matrix: 6 spellings × 6 selectivity levels
{0.01%, 0.1%, 1%, 10%, 50%, 90%} × {default, MATERIALIZED, NOT MATERIALIZED} at
reuse ∈ {1, 8}.

Two points must be calibrated before measurement, and both extend the reviewed
matrix (flagged for HIGH like the earlier families' expansions):

1. **Reuse must reach ≥5 levels.** P4's own falsifier is "no crossing across ≥5
   reuse levels"; reuse {1,8} cannot bound it. Proposed levels {1,2,4,8,16}.
2. **Materialization variants only exist for the CTE spellings.** `MATERIALIZED`
   is not valid on a plain join or a derived table, so those spellings contribute
   one cell per level, not three.

Proposed cell count per run, before calibration of the per-cell statement count:

- 3 non-CTE spellings (plain join, filtered derived table, comma join) × 6 levels
  × reuse 1 = 18 cells
- 3 CTE spellings (filtered input, joined rows, and a 6th "CTE of joined rows
  filtered outside") × 6 levels × 3 modes × reuse {1,8} = 108 cells
- plus the P4 reuse ladder at the expensive-reuse contract: 3 CTE spellings × 3
  reuse levels beyond {1,8} × 3 modes = 27 cells

≈ 150–216 cells per whole run. Every cell is 5 trials, so a whole run is on the
order of 1,000 drained queries; this is the largest SELECT matrix and the reason
this task is a separate execution.

## 3. Fixture feasibility (must be calibrated before measuring)

The existing teaching fixture is far too small for a 0.01% selectivity probe on a
fact table. The protocol requires "small indexed probes and broad scans", so the
fixture must have a fact table large enough that 0.01% selects at least one row
and 90% is a broad scan, with an index on the filtered column. Proposed: a
`cte` scale of ~200,000 fact rows with a selective indexed column, loaded once
per run; row width and fanout levels are recorded. Feasibility (load time on the
2 GiB YugabyteDB guest, and whether 0.01% is distinguishable from zero rows) must
be measured and the levels adjusted before the real matrix, and any adjustment
recorded in the Decision Log.

## 4. Controls and semantic contracts

- The reviewed control is `X-stale`; the harness already carries
  `x03_distinct_on_name` and the LEFT-JOIN predicate-placement examples as
  distinct semantic contracts. LEFT-JOIN predicate placement must be kept as a
  different answer, never a faster competitor, and every invalid substitution
  must be detected.
- A CTE-materialization control must ensure a `MATERIALIZED` CTE is actually
  materialized on the pinned engines (plan evidence), or the cell is
  `unsupported`, not silently counted as a default fold.

## 5. Arms and replication

- PostgreSQL 17.11 single-node and YugabyteDB 2025.2.6 single-node; ≥3 whole runs
  each, 5 trials, one position control per engine; a labelled
  broad-arm/feasibility pilot before the matrix.
- Reports per run and replicated per engine with the six-spelling table per
  selectivity level, the P3/P4 tables, plan signatures and cost counters.

## 6. Explicit gaps

No 3-node/equal-total/placement arm, no real network, no concurrent demand, no
write/mutation arm, no row-width sweep beyond what the fixture records, and the
>10-level feasibility calibration may remove the 0.01% level if it is not
distinguishable from zero rows.
