# Bounded measurement plan — aggregate and per-parent top-N (`select-aggregate-topn`)

LOW worker plan for the reviewed protocol's task-5 row
(`protocols/select-query-templates.md`, P8/P9). It fixes the matrix, the
controls, the levels and the replication **before** any database run.

## 1. Question and predictions

- **P8** — parent aggregate: post-join aggregate vs pre-aggregation vs
  correlated scalar. Predicted: pre-aggregation avoids the join multiplier;
  correlated pays per parent; correlated wins at small parent sets and loses at
  large.
- **P9** — per-parent top-N: LATERAL/correlated indexed probe vs `ROW_NUMBER` vs
  `DISTINCT ON` (N=1). Predicted: the indexed probe wins at small parent sets
  and loses at large; a bounded crossing exists vs the number of parents.

Falsifiers (from the protocol): no reversal across ≥5 parent levels, or totals
that disagree with the oracle (P8); no crossing across ≥5 levels, or a tie
behaviour that differs from the declared contract (P9).

## 2. Data, levels and controls

- One deterministic fixture, scale `sweep`: 10,000 customers, 2 orders each,
  seed 42; `GuaranteedTie` puts a tie at the top of the smallest multi-child
  parent.
- **Parent levels {1, 10, 100, 1000, 10000}** — the protocol's
  parent-selection {0.1%, 1%, 10%, 100%} plus 0.01%, so a crossing claim has the
  five levels the protocol requires. On a 10,000-parent fixture these are the
  same counts as the top-N parent-count axis {10, 100, 1k, 10k} plus 1.
- Cell = (family, spelling, level); all spellings in a family share the level,
  the arguments and the fixture, so each pair isolates the spelling.
- **Controls (must fire):** `x05_aggregate_double_count` (joins the child
  relation to itself, multiplying totals — the brief's wrong-grain
  counterexample) and `x06_topn_rank_ties` (RANK keeps boundary ties, a separate
  contract). Both are run live on the teaching fixture by the gate and on the
  measurement fixture by the run; a non-firing control fails the run.
- Top-N is compared **at N=1 only**, because `DISTINCT ON` is an N=1-only
  spelling; the report says so on every row.

## 3. Arms and evidence

- PostgreSQL 17.11 single-node: ≥3 whole runs, 5 trials per cell, one run also
  carrying a repeated-design position control.
- YugabyteDB 2025.2.6 single-node: 3 whole runs.
- Per cell: oracle-verified answer, client median and spread, `EXPLAIN ANALYZE`
  structure signature, planning/execution server time, decoded-row count.
- Reports carry provenance, per-run cells, per-level P8/P9 ratio tables,
  controls, plan signatures and the recorded limits. Verdicts go in a signed
  analysis beside the report, never in the generated report.

## 4. Coverage, gaps and isolates

Covered: P8 across five parent-selection levels; P9 across five parent-count
levels; both engines single-node; wrong-grain and tie controls live.

Explicit gaps: no 3-node per-node or equal-total arm, no verified placement, no
concurrent demand, no YuanybyteDB storage/RPC counters, no N>1 comparison for
`DISTINCT ON`, and no fanout/skew axis (a separate task). Persistent
rollups/materialized read models are out of scope (Study 01/04/08 evidence).

## 5. Replication and controls for order

Each whole run executes the matrix in a fixed order interleaved by family; at
least one run per engine repeats the first cell last. Across-run medians and
spreads are reported; a within-run difference inside the cells' own spread is
`no resolved difference`, and a difference is `within-run only` until whole runs
agree on its sign.
