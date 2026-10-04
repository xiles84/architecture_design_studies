# Progress — Study 01 SELECT corrective round 2 (LOW `deepseek-flash`)

**Task:** `task-20261002T131729Z-select-join-cte`
**Attempt:** `attempt-20261004T215518Z-975148`
**Claim:** `claim-03114d688bd5fe53`... see queue `guard` (recovered claim, real clock)
**Worktree:** `.worktrees/select-join-cte`; branch `study-01/select-join-cte`
**Supersedes this file's predecessor** `attempt-20261004T200000Z-c30935/PROGRESS.md`
(that attempt's own file was a byte copy of the previous one; its consequential
decisions are logged here retroactively, per corrective handoff 2 item 8).

This attempt executes
`docs/select-query-templates/reviews/20261004-join-high-2/CORRECTIVE-HANDOFF.md`
(amends update 1 and the original handoff; everything they require still applies).

## Round-2 work items

- [~] **M1 range index arm.** Implemented `yb-range-rank` (`donation (rank ASC)`),
  `IndexArmYBDefaultHash`/`IndexArmPGBtree`, `-index-arm`, `index_arm` setting.
  Feature-checked on 200k rows and in a feasibility run:
  `select-join-v3-feas-yb-range-20261005` → indexdef `(rank ASC)`, physical-arm
  gate `20/200/2000 = used`. Evidence:
  `docs/select-query-templates/reviews/20261004-join-corrective-2/evidence/yb-range-index-check.txt`.
- [~] **M1 physical-arm gate + per-cell access evidence.** `JoinFeasibility`
  gained `default_index_gate`/`default_index_plans`; the matrix aborts if the
  default join does not use the rank index. Cells carry index/seq scan, YB
  Storage Rows Scanned, index rows scanned, PG buffers, rows removed and retries.
- [~] **M2 colocated arm.** `infra/yb-cluster3.sh` gained `colocated-db` /
  `colocated-endpoints` (`CREATE DATABASE … WITH COLOCATION = true`); the runner
  captures `yb-admin list_tablets` leaders; the harness records `pg_indexes.indexdef`
  and `yb_table_properties` + `yb_local_tablets`.
- [~] **M3 class includes arm factors.** `node_budget`, `client_budget`,
  `index_arm` added; class description names topology, budgets, placement, phase,
  index arm and endpoint mode; tests assert budget and index-arm separation.
- [~] **M4 truthful limits and throttling.** `joinGaps` is derived from the run's
  settings (no "no concurrent writers" on blended; no cluster limit on single);
  `distributed_rpc` text fixed; per-node and client `cpu.stat` deltas rendered in
  per-run and replicated reports; replicated report lists endpoint ops, writer
  counts, retries and overlap.
- [~] **Cluster writer rotation.** The writer rotates across endpoints and records
  `endpoint_ops`.
- [ ] **Matrix runs** (both engines, both phases, 3 whole runs; topology subset on
  the range arm; colocated cluster3) — see below.
- [ ] **Corrected topology report**, pooled one moved to `reports/outdated/`,
  digest-stability proof.
- [ ] **v3 signed analysis** superseding v2; v1 and v2 moved to `reports/outdated/`.

## Retroactive Decision Log — decisions of attempt `c30935` (previously unlogged)

1. **Reader retry policy for transient conflicts.** Context: the first blended
   YB run aborted with 2/30 cells `SQLSTATE 40001`. Decision: re-execute the
   statement up to 10 times with a short backoff, record the count on the cell,
   and verify the oracle after every retry; `EXPLAIN ANALYZE` is retried by the
   same policy because it executes the statement. Reason: a transient conflict is
   not a wrong answer, and aborting makes the blended arm unmeasurable on YB.
   Confidence: High. Alternative considered: report the failure rate and do not
   retry (rejected: three of thirty cells would be untimed for a load artifact).
   Affects: all blended arms; the report states that sample latency includes
   retry and backoff.
2. **Domain tables bound through views.** Context: the catalogue SQL must be
   byte-identical between the synthetic and domain arms. Decision: real
   `charity`/`person`/`donation` tables plus `customers`/`facts` views. Reason:
   the measured SQL and the oracle contract stay unchanged; only the physical
   fixture changes. Confidence: High. Affects: every domain run.
3. **Domain matrix re-run.** Decision: re-run the 3+3 matrix from one frozen
   commit (`select-join-domain2-*`) because the earlier six runs spanned two
   client builds and could not be pooled. Confidence: High.
4. **Writer endpoint = one endpoint.** Decision at the time: bind the writer to
   one pool. Reason: the task had one writer and a single-endpoint contract.
   Correction: review M-gap and handoff item 5 require rotation; now implemented
   and recorded as `endpoint_ops`. Confidence in the correction: High.
5. **Placement conclusion.** Decision: record `colocation = true` on the default
   database as an engine rejection (SQLSTATE 0A000) and label the colocated arm a
   gap. Correction: review M2 — that was untested, not unsupported; the colocated
   *database* arm is now attempted. Confidence: High.

## Decision Log — this attempt (round 2)

1. **Range index as a physical arm, not a syntax effect.** `(rank ASC)` is a
   different physical design (range-sharded) from `(rank)` (hash-sharded); both
   are labelled and never pooled. Confidence: High; feature-checked.
2. **Gate as a precondition.** A run whose default join does not use the rank
   index writes its plans and stops before measuring. Confidence: High; the
   handout requires it and it is the guard against the M1 defect recurring.
3. **`colocated-db` creation in the topology script** rather than in the
   harness: database creation is infra, and a run must not silently create
   clusters. Confidence: High.
4. **Re-render the per-run report in the runner** after writing
   `node-throttle.json`, because the "after" counters exist only once the run has
   ended; the re-render hashes the same join.json bytes, so the digest is stable.
   Confidence: High.
5. **No `--now` on live queue commands** (process finding from review 2); the
   claim was taken with the real clock. Confidence: High.

## Status

Code frozen at `b108fafe` (plus evidence commit `65f79e04`). `go vet` and
`go test` pass. The matrix is starting; each run uses a new id under
`select-join-v3-…-20261005` and an annotated producing tag.
