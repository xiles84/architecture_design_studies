# Task result — task-20261004T100822Z-select-order-correction

| Field | Value |
|---|---|
| Task | `task-20261004T100822Z-select-order-correction` (LOW executor, deepseek-flash / opencode CLI) |
| Status | measurements complete; submitted for HIGH review |
| Branch / worktree | `study-01/select-order-correction` / `.worktrees/select-order-correction` |
| Source | `docs/select-evidence-review-20261004/LOW-CORRECTION-HANDOFF.md` |
| Code commit | `ca2a38fb` (order machinery) + the run-id/report fixes |
| Milestone tag (on integration) | `study-01/select-v7-order-controls` |

## What was corrected

Each whole run now permutes its cells with a seed derived from the dataset seed
and the run id, verifies/plans every cell once, and then takes one sample of every
cell per round (interleaved). The seed, round count and actual axis sequence are
recorded in the run JSON and rendered in the report. The position control repeats
the **actual first cell** of that run's randomized order; the access expression
index is a separate labelled phase with its own order; existence keeps its
per-fanout phases. Correctness is unchanged: every timed query is drained and
oracle-verified, and all fixture and control gates fired in every run.

## Runs (18 whole runs, 5 trials, serial under ads-run-lock)

| family | PostgreSQL | YugabyteDB | cells |
|---|---|---|---|
| access | `select-order-access-postgres-1/2/3-20261004` | `...-yugabyte-1/2/3-...` | 34–35 valid each |
| aggregate | `select-order-aggregate-postgres-1/2/3-...` | `...-yugabyte-1/2/3-...` | 30–31 valid each |
| existence | `select-order-existence-postgres-1/2/3-...` | `...-yugabyte-1/2/3-...` | 21–22 valid each |

All cells valid (0 invalid/unsupported/timeout/error). Run 1 per engine/family
carries the position control. Each family/engine has 3 distinct recorded
sequences.

## Result

- **Position drift is now measured**: access 1.5–2.4%, existence 13–22%,
  aggregate 29–74% first-to-last of the actual first cell. Only effects larger
  than the drift are promoted.
- **Access P11 holds** with a stable sign at every depth: PostgreSQL crossing
  between 100 and 1,000 (depth 100 now 3/3, keyset slower by 6–17%),
  YugabyteDB between 1,000 and 10,000.
- **Existence P6/P7 hold** under randomized order: PG `JOIN + DISTINCT` over
  `EXISTS` is 1.40–1.69× at fanout 1, 1.38–1.44× at fanout 10 and 3.64–4.23× at
  fanout 100 (3/3); YB within noise; nullable `NOT IN` empty in every run.
- **Aggregate P8/P9 partly unresolved**: `DISTINCT ON` vs `ROW_NUMBER` is
  mixed-sign at PG level 10 and 100 and at YB level 100; stable only at the
  levels reported in the signed analysis. The earlier "`DISTINCT ON` fastest at
  scale on both engines" is narrowed accordingly.

## Deliverables

- Order machinery in `selectbench/runner.go` (`PlanOrder`, `OrderSeedFor`,
  `RunSpecsInterleaved`, `RunOrder`) with `order_test.go`.
- Recorded order rendered in all three per-run reports; replicated ratio rows
  carry run ids; `-cmd manifest` emits per-file and corpus digests.
- New runs under `results/select-order-*`; reports
  `reports/20261004T-select-order-{access,agg,exist}-replicated.md`; manifests
  `reports/20261004T-select-order-manifest-*.md`; signed analysis
  `reports/analyses/20261004T-select-order-correction--deepseek-flash--2026-10-04.md`.
- Corrected environment page
  `docs/environments/xiles-podman-select-order-correction-20261004.md`.

## Explicit gaps

Single-node only (no cluster per-node/equal-total arm despite now-recorded live
resources); one sequential reader; five samples; three pagination levels bound
but do not locate the P11 crossing; aggregate drift is the largest and the family
most in need of more rounds/runs.

## Next

HIGH review of the correction before numeric promotion and local-main
integration; the original fixed-order runs and analyses are retained and not
edited.
