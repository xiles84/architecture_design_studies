# Task result — task-20261002T131729Z-select-access-paths

| Field | Value |
|---|---|
| Task | `task-20261002T131729Z-select-access-paths` (LOW executor, deepseek-flash / opencode CLI) |
| Status | measurements complete; submitted for HIGH review |
| Branch / worktree | `study-01/select-access-paths` / `.worktrees/select-access-paths` |
| Premeasurement commits | `68b09e6d` (harness/runner/plan), `30f3b00c` (selective+broad regimes), `6ad01ef6` (position control) |
| Measurement commit (34-cell primary runs) | `524ce328` |
| Milestone tag (on integration) | `study-01/select-v5-access-paths` |

## What was run

Seven whole runs of the bounded access-path matrix, under `ads-run-lock`, on the
fixed host (`xiles-podman-select-access-20261003`), scale `deep` (12,000 orders,
seed 42), page limit 5, five trials per cell:

| Run tag | Engine | Commit | Cells | Inputs digest |
|---|---|---|---|---|
| `run/study-01/select-access-pg-1-20261003` | postgres | `524ce328` | 34 valid | `2bb893ca5c35fe65` |
| `run/study-01/select-access-pg-2-20261003` | postgres | `524ce328` | 34 valid | `24ceae2fe2933f2a` |
| `run/study-01/select-access-pg-3-20261003` | postgres | `524ce328` | 34 valid | per run `access.json` |
| `run/study-01/select-access-pg-pos-20261003` | postgres | `6ad01ef6` | 35 valid (position control) | `8747043ad53b9991` |
| `run/study-01/select-access-yb-1-20261003` | yugabyte | `6ad01ef6` | 34 valid | `ce2139aceea0cd56` |
| `run/study-01/select-access-yb-2-20261003` | yugabyte | `6ad01ef6` | 35 valid (position control) | `bf5210d941f96803` |
| `run/study-01/select-access-yb-3-20261003` | yugabyte | `6ad01ef6` | 34 valid | `13e9fc78aae1a024` |

Every cell was verified against the independent Go oracle before timing; there
were **0 invalid, 0 unsupported, 0 timed-out and 0 errored cells**. A labelled
broad-arms-only pilot tranche (`select-access-pg-r{1,2,3}`, `select-access-yb-r1`)
was also run and committed; it is corroboration only and is not the primary set.

## Result

- **P11 — OFFSET vs keyset: confirmed**, with an engine-specific crossing
  (PostgreSQL between depth 100 and 1,000; YugabyteDB between 1,000 and 10,000).
  At depth 10,000 the keyset/OFFSET ratio is 0.20–0.23 on PG (4/4 runs) and
  0.54–0.69 on YB (3/3). Page equality under `REPEATABLE READ` held at all five
  depths in all runs.
- **P10 — OR vs UNION ALL: falsified as stated.** `OR` was never worse by a
  replicated margin; on YugabyteDB it beat both `UNION` and `UNION ALL` in all
  36 shape×run pairs (1.13–2.12×). The expected "OR blocks an index" penalty did
  not appear because the selective indexable `OR` is served by a single BitmapOr.
- **Expression-index package arm:** on PostgreSQL the selective expression `OR`
  fell from 0.845 ms to 0.265 ms (3.2×) with the index present, storage 114,688
  bytes; on YugabyteDB it did not improve (4.05 → 4.15 ms) and the size probe
  returned 0.
- **Controls:** the required `x05_union_all_overlap` control was detected in all
  six shapes in all seven runs; `x01`–`x03` are recorded `control_absent` on the
  measurement fixture.
- **Position control:** first vs last ratio 0.885 (PG) / 1.112 (YB), same plan.
- **Client:** 2 CPU / 2 GiB quota, 0 throttled periods; the read matrix is not
  client-bound.

Verdicts and reasoning are in the signed analysis
`studies/01-charity-tree/reports/analyses/20261003T-select-access--deepseek-flash--2026-10-03.md`;
measurements are in `studies/01-charity-tree/reports/20261003T-select-access-replicated.md`
and each run's `results/<run-id>/access.md`.

## Explicit limits and gaps

- Single host, single node. No YB 3-node RF=3, no equal-total-budget arm, no
  verified tablet/leader placement, no independent hardware or real network.
  These are coverage gaps, not completed comparisons.
- Latency of tiny pages from one sequential reader; no throughput, no
  concurrency, no demand generator; five trials cannot support p95/p99.
- The YugabyteDB expression-index storage/maintenance probe is unusable
  (`pg_relation_size` 0; with-index inserts faster than without, 3/3) and is
  reported as absent, not as a cost.
- The declared 3 GiB + 2 GiB budgets still exceed the 2 GiB guest.

## Deliverables

- Harness: `studies/01-charity-tree/selectbench/` (s18–s20, x05, oracles, matrix,
  reports, tests) and `studies/01-charity-tree/harness/select_access_mode.go`
  (`-cmd select-access`, `-cmd access-report`, `-position-control`).
- Runner: `studies/01-charity-tree/run-select.sh --family access`.
- Plan: `studies/01-charity-tree/protocols/select-access-paths-bounded-matrix.md`.
- Reports: per-run `results/<run-id>/access.{json,md}`; replicated
  `reports/20261003T-select-access-replicated.md`.
- Analysis: `reports/analyses/20261003T-select-access--deepseek-flash--2026-10-03.md`.

## Next

HIGH review of the Decision Log, the two-regime interpretation of the protocol's
`select-access-paths` row (an **expansion** flagged in `PROGRESS.md`), and the
material evidence; then local-main integration and the milestone tag. Book
integration belongs to the later `select-book` task; evidence synthesis belongs
to `select-synthesis`.
