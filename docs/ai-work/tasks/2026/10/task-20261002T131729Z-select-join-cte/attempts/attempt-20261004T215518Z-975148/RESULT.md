# Result — Study 01 SELECT join/CTE corrective round 2

**Task:** `task-20261002T131729Z-select-join-cte`
**Attempt:** `attempt-20261004T215518Z-975148`
**Capability/role:** LOW executor (model `deepseek-flash`, opencode CLI)
**Branch/worktree:** `study-01/select-join-cte` / `.worktrees/select-join-cte`
**Handoff executed:** `docs/select-query-templates/reviews/20261004-join-high-2/CORRECTIVE-HANDOFF.md`
**Measured:** yes — 37 v3b runs, one frozen client build, every cell verified.

## Delivered against the acceptance list

### M1 — the rank index is a real, checkable physical arm

- `yb-range-rank` creates `donation (rank ASC)`; `yb-default-hash-rank` and
  `pg-btree-rank` are separate labelled arms. Feature check (200k rows, evidence
  in `docs/select-query-templates/reviews/20261004-join-corrective-2/evidence/`):
  `(rank)` is hash-sharded (`Storage Filter` over every row, cost 10024);
  `(rank ASC)` gives `Index Cond` and **Storage Index Rows Scanned: 20**.
- A **physical-arm gate** checks the default join at the handoff's thresholds
  20/200/2000 before any matrix; failures record the plan and abort. The gate
  initially used the first three thresholds of the run's axis; the topology
  subset's third is 180000 (90%), where a seq scan is correct, so it wrongly
  refused to measure it (`select-join-v3-topo-single-i1-20261005` is retained as
  that evidence). Fixed in `7243fee0`; all arms were then re-run under the new
  frozen build (`v3b`).
- Result: YB index-scan plans 0 → 66–67 per domain run; `j01` at 20 rows scans
  110 storage rows (20 index rows) instead of 200,000. The corrected engine
  contrast for forced `MATERIALIZED` is PostgreSQL ~97x vs YugabyteDB ~42x at 20
  rows (v2's 1.4–1.7x was the hash-index artefact).

### M2 — colocated arm actually attempted, with placement evidence

- `infra/yb-cluster3.sh` creates `select_coloc` with `COLOCATION = true` and
  prints its endpoints; the runner passes that DSN for `--placement colocated`.
- All 6 colocated cluster3 runs (isolated ×3, blended ×3) are 30/30 valid and
  record `is_colocated = true`, one tablet per table.
- Placement evidence per cluster run: `yb_table_properties`, `yb_local_tablets`
  counts (in the result JSON) and `yb-admin list_tablets` leader-per-tablet
  (`resources/placement.txt`). Same-host containers are not used as colocation
  evidence.

### M3 — every arm factor separates compatibility classes

- Class includes environment, dataset, client/engine digests, and settings:
  topology, node budget, client budget, placement, phase, index arm, endpoint
  mode. The class description names them. Tests:
  `TestJoinClassSeparatesNodeBudget`, `TestJoinClassSeparatesIndexArm`.
- The topology report's six classes are distinct by budget, placement and phase;
  no budget pooling remains.

### M4 — truthful limits and rendered evidence

- `joinGaps` is derived from the run's own settings: a blended run records its
  writer contract, an isolated run says no writers *in this run*, a single-node
  run carries no cluster limit, and the stale `distributed_rpc` "not parsed"
  text is replaced. Test: `TestJoinGapsMatchRunSettings`.
- Per-node and client `cpu.stat` deltas (`nr_periods`, `nr_throttled`,
  `throttled_usec`) are rendered from `resources/node-throttle.json` in per-run
  and replicated reports; the replicated report lists per-run endpoint
  operations, writer offered/completed/failed, reader retries and overlap.
- Statement that sample latency includes retry and backoff, and per-cell retry
  counts, are rendered.

### Runs (all correctness-gated, one build)

| arm | runs | result |
|---|---|---|
| domain PostgreSQL `pg-btree-rank`, isolated + blended | 6 | 91/91 valid each |
| domain YugabyteDB `yb-range-rank`, isolated + blended | 6 | 91/91 valid each |
| topology single (2 CPU/3 GiB), equal-total single (6 CPU/9 GiB), cluster3, isolated + blended | 18 | 30/30 valid each |
| colocated cluster3, isolated + blended | 6 | 30/30 valid each, `is_colocated=true` |
| gate re-validation | 1 | gate 20/200/2000 = used |

The writer sustained 10.0 updates/s with 0 failures and rotated across the
cluster's endpoints (e.g. [59,59,58]); node throttling is recorded (e.g. yb-n3 15
periods / 650 ms). Both semantic controls fire in every run; no invalid, timeout
or error cell.

### Reporting and analysis

- `reports/select-join-v3b-domain-replicated-20261005.md`,
  `reports/select-join-v3b-topology-20261005.md`; the topology report re-renders
  byte-identically (digest stability proven).
- The pooled topology report and analyses v1/v2 moved to `reports/outdated/`
  with `REPLACEMENTS-select-join.md`.
- v3 signed successor analysis
  `reports/analyses/20261004T-select-join-cte-v3--deepseek-flash--2026-10-05.md`
  supersedes v2: separate P1–P4 and P12 verdicts with P5 as a control; YB per
  index arm; the hash/range pair as a physical factor; YB reuse n=1 kept
  mixed-sign/unresolved; colocated vs non-colocated and topology per budget
  compared; the rollup/rolldown/embedding non-applicability corrected to
  "out of scope for this spelling experiment", not "absent from the model".

### Gate checks

- `go vet ./...` and `go test ./...` pass in `docker.io/library/golang:1.26-bookworm`
  under Podman.
- `tools/queue/queue audit`: 91 tasks, **0 errors** (3 warnings: two pre-existing
  on an unrelated completed task, one live-ref drift on this task's claim ref).
- Tree clean; benchmark lock free; no containers running; no push; completion tag
  not created.

## Notes and remaining gaps

- Three whole runs per arm; one host; cluster nodes are containers on one
  machine, so colocation is evidenced by engine metadata, not independent
  hardware.
- The writer is a fixed maintenance load, not an optimistic/pessimistic
  concurrency comparison.
- The v3 (pre-fix-gate) runs and the v2 hash-arm runs are retained as labelled
  cross-build/physical-arm observations and are not pooled with v3b.
- No book claim is promoted; HIGH review owns promotion.
