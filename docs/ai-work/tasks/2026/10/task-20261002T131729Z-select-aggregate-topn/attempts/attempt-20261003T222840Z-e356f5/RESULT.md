# Task result — task-20261002T131729Z-select-aggregate-topn

| Field | Value |
|---|---|
| Task | `task-20261002T131729Z-select-aggregate-topn` (LOW executor, deepseek-flash / opencode CLI) |
| Status | measurements complete; submitted for HIGH review |
| Branch / worktree | `study-01/select-aggregate-topn` / `.worktrees/select-aggregate-topn` |
| Premeasurement commits | `00f2aa3d` (harness/runner), `ad048e0c` (plan/env/Decision Log) |
| Measurement commits | `08a3c553` (pg-1b), `33b76901` (pg-2/3, yb-1/2/3) |
| Milestone tag (on integration) | `study-01/select-v4-aggregate-topn` |

## What was run

A 10,000-parent deterministic fixture (2 orders per parent, seed 42, `sweep`
scale), parent-prefix levels {1, 10, 100, 1000, 10000}, five trials per cell,
one sequential reader, under `ads-run-lock`, on
`xiles-podman-select-aggregate-20261003`:

| Run tag | Engine | Cells | Position control |
|---|---|---|---|
| `run/study-01/select-agg-pg-1b-20261003` | postgres | 30 valid | no |
| `run/study-01/select-agg-pg-2-20261003` | postgres | 30 valid | no |
| `run/study-01/select-agg-pg-3-20261003` | postgres | 31 valid | yes |
| `run/study-01/select-agg-yb-1-20261003` | yugabyte | 30 valid | no |
| `run/study-01/select-agg-yb-2-20261003` | yugabyte | 31 valid | yes |
| `run/study-01/select-agg-yb-3-20261003` | yugabyte | 30 valid | no |

All cells were verified against the independent Go oracle before timing; there
were **0 invalid, 0 unsupported, 0 timed-out and 0 errored cells**, and both
family controls fired in every run (`x05` 10,000 rows; `x06` 10,021 rows). A
first PostgreSQL attempt (`select-agg-pg-1-20261003`, commit `08a3c553`) panicked
in the harness before measuring — its tag and log are retained as a failed
harness run, not a measurement.

## Result

- **P8 confirmed on PostgreSQL, refuted on YugabyteDB.** PG: correlated wins
  below ~10 parents (0.70–0.98× post-join) and loses from 100 (crossing between
  10 and 100, 3/3 runs); pre-aggregation wins only at 10,000 parents (crossing
  between 1,000 and 10,000, 3/3 runs). YB: pre-aggregation never wins
  (1.13–3.26×), and the correlated shape collapses from 3.99 ms at one parent to
  **8,460 ms** at 10,000 (489–507× post-join).
- **P9 partially confirmed; the protocol's favourite is wrong at scale.** LATERAL
  has no resolved small-parent advantage on PG, loses from 100 parents, and on YB
  wins only at one parent before losing 272–278× at 10,000. `DISTINCT ON` is the
  fastest top-N spelling from 100 parents up on both engines (PG 0.63–0.73×
  ROW_NUMBER; YB 0.83–1.02×).
- **Controls:** `x05_aggregate_double_count` and `x06_topn_rank_ties` fired in
  6/6 runs; the wrong-grain control would otherwise pass as a fast aggregate.
- **Position control:** first/last ratio 1.278 (PG) and 0.951 (YB), same plan.
- **Client:** 2 CPU/2 GiB quota, 0 throttled periods.

Verdicts and reasoning are in the signed analysis
`studies/01-charity-tree/reports/analyses/20261003T-select-agg--deepseek-flash--2026-10-03.md`;
measurements are in `studies/01-charity-tree/reports/20261003T-select-agg-replicated.md`
and each run's `results/<run-id>/aggregate.md`.

## Explicit limits and gaps

- Crossings are brackets (10–100, 1,000–10,000), not located; the one-parent
  level is sub-millisecond and fixed-overhead dominated.
- Single host, single node: no YB 3-node RF=3 arm, no equal-total budget, no
  verified placement, no real network. YB absolute times include single-node
  YSQL proxy cost; storage/RPC counters are not parsed.
- Fixed two children per parent, no fanout/skew axis; N=1 only for the top-N
  family (`DISTINCT ON` cannot express N>1); no real second one-to-many relation
  for the join-multiplier contract.
- Five samples per cell cannot support p95/p99; latency only, no throughput.

## Deliverables

- Harness: `studies/01-charity-tree/harness/select_aggregate_mode.go` (`-cmd
  select-aggregate`, `-cmd aggregate-report`, `-position-control`), scale `sweep`.
- Catalogue/oracles: `selectbench/sql/select_aggregate_catalog.sql` (a01–a06,
  x05, x06), `selectbench/aggregate.go`, `selectbench/aggregatereport.go`,
  fixture/oracle additions.
- Runner: `studies/01-charity-tree/run-select.sh --family aggregate`.
- Plan: `studies/01-charity-tree/protocols/select-aggregate-topn-bounded-matrix.md`.
- Reports: per-run `results/<run-id>/aggregate.{json,md}`; replicated
  `reports/20261003T-select-agg-replicated.md`.
- Analysis: `reports/analyses/20261003T-select-agg--deepseek-flash--2026-10-03.md`.

## Next

HIGH review of the Decision Log (notably the added 0.01% level and the
`RunCells`→`RunSpecs` refactor shared with the unmerged access-paths branch) and
the material evidence, then local-main integration and the milestone tag.
