# Task result — task-20261002T131729Z-select-existence

| Field | Value |
|---|---|
| Task | `task-20261002T131729Z-select-existence` (LOW executor, deepseek-flash / opencode CLI) |
| Status | measurements complete; submitted for HIGH review |
| Branch / worktree | `study-01/select-existence` / `.worktrees/select-existence` |
| Premeasurement commit | `78fdd9c9` (plan/env/Decision Log) |
| Measurement commit | `805b5b85` (harness shape fix, all six runs) |
| Milestone tag (on integration) | `study-01/select-v3-existence` |

## What was run

A fanout sweep holding 10,000 child rows constant with 1, 10 and 100 children
per parent (10,000 / 1,000 / 100 parents), seed 42, five trials per cell, one
sequential reader, under `ads-run-lock`, on
`xiles-podman-select-existence-20261004`:

| Run tag | Engine | Cells | Position control |
|---|---|---|---|
| `run/study-01/select-exist-pg-1b-20261004` | postgres | 21 valid | no |
| `run/study-01/select-exist-pg-2-20261004` | postgres | 22 valid | yes |
| `run/study-01/select-exist-pg-3-20261004` | postgres | 21 valid | no |
| `run/study-01/select-exist-yb-1-20261004` | yugabyte | 21 valid | no |
| `run/study-01/select-exist-yb-2-20261004` | yugabyte | 22 valid | yes |
| `run/study-01/select-exist-yb-3-20261004` | yugabyte | 21 valid | no |

All cells were verified against the independent Go oracle before timing; there
were **0 invalid, 0 unsupported, 0 timed-out and 0 errored cells**. Both
controls fired at every fanout level in every run: x02 nullable `NOT IN`
returned 0 rows (18/18) and x07 membership multiplier returned 7,240–7,241 rows
(18/18). A first PostgreSQL attempt (`select-exist-pg-1-20261004`) failed its
gate because the harness built the raw fixture; it is retained as a
harness-defect record and its tag is not reused.

## Result

- **P6 confirmed on PostgreSQL, refuted on YugabyteDB.** PG `JOIN + DISTINCT`
  over `EXISTS`: 1.44–1.47× at fanout 1, 1.30–1.57× at fanout 10, **3.61–4.40×**
  at fanout 100, stable 3/3 runs. YB: 0.92–1.06× at every level, no trend.
- **`IN` and `EXISTS` are interchangeable** on both engines (PG 0.86–1.13, YB
  0.99–1.04, mixed signs).
- **P7's NULL trap confirmed:** nullable `NOT IN` returned the empty set at every
  level (x02 fired 18/18 with 0 rows); the NOT-NULL `NOT IN` spelling matched the
  absence oracle. On PG, `NOT EXISTS` also beats `LEFT JOIN ... IS NULL` and
  `NOT IN` by 2–3× at fanout 100 (0.22–0.34 ms vs 0.63–0.72 ms, 3/3 runs); on YB
  the three absence spellings are within ~10%.
- **Position control:** 0.871 (PG) and 0.995 (YB), same plan.
- **Client:** 2 CPU/2 GiB quota, not throttled.

Verdicts and reasoning are in the signed analysis
`studies/01-charity-tree/reports/analyses/20261004T-select-exist--deepseek-flash--2026-10-04.md`;
measurements are in `studies/01-charity-tree/reports/20261004T-select-exist-replicated.md`
and each run's `results/<run-id>/existence.md`.

## Explicit limits and gaps

- Fanout and parent count move together (constant child volume), so the growing
  `JOIN + DISTINCT` penalty is confounded with the shrinking parent set; three
  levels bound the direction but locate no crossing.
- No selectivity, hit-fraction or skew axis; no 3-node/equal-total/placement arm;
  no index redesign arm; five samples; latency only. YugabyteDB comparisons sit
  near its per-statement floor and storage/RPC counters are not parsed.

## Deliverables

- Harness: `studies/01-charity-tree/harness/select_existence_mode.go` (`-cmd
  select-existence`, `-cmd existence-report`, `-position-control`), scale `fanout`.
- Catalogue/oracles: `selectbench/sql/select_existence_catalog.sql` (e12, e13,
  x07), `selectbench/existence.go`, `selectbench/existencereport.go`,
  `ExistenceShape`, `MembershipWithMultiplicity`.
- Runner: `studies/01-charity-tree/run-select.sh --family existence`.
- Plan: `studies/01-charity-tree/protocols/select-existence-bounded-matrix.md`.
- Reports: per-run `results/<run-id>/existence.{json,md}`; replicated
  `reports/20261004T-select-exist-replicated.md`.
- Analysis: `reports/analyses/20261004T-select-exist--deepseek-flash--2026-10-04.md`.

## Next

HIGH review of the Decision Log and the material evidence, then local-main
integration and the milestone tag `study-01/select-v3-existence`.
