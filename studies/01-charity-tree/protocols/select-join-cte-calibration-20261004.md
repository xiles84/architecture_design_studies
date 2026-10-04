# Join/CTE matrix calibration and premeasurement commitment (2026-10-04)

LOW executor: model `deepseek-flash`, tool opencode/CLI, effort unknown, session
capability LOW, role executor. Task `task-20261002T131729Z-select-join-cte`,
branch `study-01/select-join-cte`, milestone tag `study-01/select-v2-join-cte`.

This file is **new** and does not edit the previous LOW worker's committed
bounded-matrix plan (`select-join-cte-bounded-matrix.md`), which remains the
record of what was proposed and yielded. It commits the calibrated matrix and
fixture **before any measurement**, as the handoff requires.

## Fixture calibration

The handoff fixes 200,000 fact rows with deterministic shuffled ranks covering
1..200,000 and a fixed parent set. The implemented fixture is exactly that:

- schema `selectjoin`: `customers(id, name)` and
  `facts(id, customer_id, rank, amount, payload)`;
- `rank` is a seeded Fisher-Yates permutation of 1..N, so `rank <= threshold`
  selects an exact count; the predicate is a labelled proxy for the teaching
  amount predicate, not the same question;
- fixed index arm for every spelling pair: `facts_rank_idx(rank)`,
  `facts_customer_idx(customer_id)` and the primary keys. No per-spelling index;
- thresholds for 200,000 rows: **20 / 200 / 2,000 / 20,000 / 100,000 /
  180,000** (= 0.01% / 0.1% / 1% / 10% / 50% / 90%).

Feasibility may reject this size; if so the rejection and the new size are
recorded in the attempt Decision Log and here, and the axes are not silently
changed.

## Calibrated matrix

**Selectivity matrix** — 10 bag spellings x 6 thresholds = 60 cells:
`j01_join_where`, `j02_join_on`, `j03_filtered_derived`, `j04_cte_filtered`,
`j04m_cte_filtered_mat`, `j04n_cte_filtered_notmat`, `j05_cte_joined_outer`,
`j05m_cte_joined_outer_mat`, `j05n_cte_joined_outer_notmat`, `j06_comma_join`.

**Reuse matrix** — 3 policies x 5 reference counts x 2 consumer selections = 30
cells, at a fixed 1% selective CTE and a per-consumer filter outside the common
CTE: policies `default` / `MATERIALIZED` / `NOT MATERIALIZED`, references
**1, 2, 4, 8, 16**, consumer selection **1% / 100%**. The CTE consumes the
side-effect-free token `ascii(substr(md5(repeat(id::text, 32)), 1, 1))`; the
token's support is verified per engine during feasibility.

Whole-run size: 90 cells x 5 trials = 450 drained queries plus the gate,
feasibility probes and controls.

## Controls

- `x01_optional_where` (child predicate in WHERE) is detected against the
  optional-match oracle on the **edge fixture** (duplicate matches, NULL amount,
  an orphan fact and an unmatched parent). On the large fixture at 10% almost
  every parent has a qualifying child, so this control is not run there.
- `x02_join_limit_before` (fixed limit 7) is detected against the bag oracle on
  the **large fixture**, where it truncates.

## Predictions under test

- **P1** six inner-join spellings: no resolved difference in most regimes; a
  difference only where a CTE materializes or an index cannot be used.
- **P2** filtered derived table vs plain join: no resolved difference.
- **P3** default CTE vs `MATERIALIZED` under a selective one-use predicate:
  `MATERIALIZED` slower; reverses under high reuse.
- **P4** `MATERIALIZED` vs `NOT MATERIALIZED` under repeated reuse:
  `MATERIALIZED` faster above a reuse threshold, bounded by the five levels.
- **P12** plan-equivalent variants and no-resolved-difference outcomes are
  retained, not discarded.

## Engines, replication and arms

- PostgreSQL 17 single-node and pinned YugabyteDB single-node: >=3 whole runs
  each, 5 trials per cell, a seeded per-run cell permutation with interleaved
  samples, and the actual first cell repeated at the end (position control).
- The handoff's 3-node YB RF=3 subset and the colocated/non-colocated placement
  pair are attempted only after the single-node work. If the fixed 2 GiB podman
  machine cannot carry them, they are recorded as explicit coverage gaps with
  the failing resource evidence, never presented as a completed comparison.
  Same-host containers are not data-colocation evidence.
