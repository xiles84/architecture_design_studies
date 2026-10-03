# Progress — select-harness (attempt-20261003T014215Z-c6df8d)

**Status: in progress, checkpoint.** The correctness floor is built and checked in the pinned image;
the DB execution, plan capture and report layers are not written yet. Committed on
`study-01/select-harness`; queue checkpoint keeps the claim open for the next LOW session.

## Done and checked

- `studies/01-charity-tree/selectbench/` — additive Go package, no edit to the historical harness:
  - `sql/select_catalog.sql` — named statements with `group:` semantic ids, explicit projections and
    `params:`, including four deliberately-wrong controls in group `control_wrong_answer`.
  - `catalog.go` — `go:embed` loader, deterministic order, strict validation (a statement with no
    group or no SQL is a load error) and the `Comparable` guard that refuses cross-group or
    control-involving pairs.
  - `fixture.go` — the edge-case teaching fixture and a seeded generator (narrow/wide, fanout, skew,
    NULLs), plus independent Go oracles: bag, membership, absence, optional match, aggregate, top-N
    (declared total order `amount DESC NULLS LAST, id ASC`) and ordered page.
  - `eval.go` — Go mirrors of every statement, the keyset page, and the four controls.
  - `selectbench_test.go` — 8 tests, all passing.
- `studies/01-charity-tree/reports/20261003T-select-harness-check.md` — the check report.
- Verified in `golang:1.26-bookworm`: `gofmt` clean, `go vet` clean, `go build ./...` OK (existing
  harness reproducible), `go test ./selectbench/...` PASS.

## Decision Log

1. **Additive package, not a new binary and not a migration.** `selectbench` is a separate package in
   the existing `charitytree` module; the historical harness is untouched, so old Study 01 execution
   stays reproducible. Confidence high.
2. **Semantic group ids in the SQL comments, enforced by code.** The loader rejects an ungrouped
   statement and `Comparable` refuses cross-group pairs, so a faster statement answering a different
   question cannot be pooled. Confidence high.
3. **`Amount`/`CustomerID` are pointers so NULL is a real NULL**, because the NULL cases carry the
   controls. Confidence high.
4. **Declared total order `amount DESC NULLS LAST, id ASC` everywhere** (top-N and pagination), and
   keyset/offset must agree under it; the order is written down rather than assumed. Confidence high.
5. **The fixture is itself tested.** A control test that passes on data lacking the edge case proves
   nothing, so `TestFixtureExercisesEveryRequiredEdgeCase` guards NULLs, duplicate children, unmatched
   parents and all-miss sets. Confidence high.
6. **Checks ran in the pinned `golang:1.26-bookworm` image**, matching the Containerfile's build stage,
   without the benchmark lock because no database starts. Confidence high.
7. **No performance number is produced here.** The evaluators are contract checks, not timings; the
   comparison report must never present a control as a faster alternative. Confidence high.

## Update (2026-10-03, second increment)

`plan.go` + `plan_test.go` added and committed (`d9683dfc`): EXPLAIN normalization and signatures
(costs/rows/times/buffers/print-order stripped; access paths, index names and materialization kept).
13 tests pass in the pinned image; the check report was extended. Remaining item 3 below is now
**half done** — capture (`EXPLAIN (ANALYZE, BUFFERS)` per cell, stored beside the timing) is still to do.

## Update (2026-10-03, third increment)

`report.go` + `report_test.go` committed (`74009573`): `Status`/`Provenance`/`Cell` and
`RenderComparison`, the measurements-only generated report. It names invalid/unsupported/timed-out
cells and never ranks them, refuses cross-group and control pairs, classifies each comparable pair as
`same plan structure`, `no resolved difference` (difference inside the trial spread) or a resolved
difference, and **refuses a non-empty interpretation** so the generated file cannot absorb the
analysis. 6 new tests; **19 total pass in `golang:1.26-bookworm`**; check report extended. Remaining
item 4 is now mostly built; the DB execution (items 1–3) is still the gap.

## Remaining work (for the next session on this branch)

1. `-cmd select` (or an equivalent subcommand) in the harness `main.go`: bind `:pred`/`:n`/`:limit`
   parameters, execute the embedded statements against PG single-node and YB single-node, and drain
   results. This entrypoint must call `run_lock_acquire` before any database work.
2. Independent-oracle verification per cell: execute the SQL, compare against the Go oracle, abort the
   cell on a mismatch; assert each control is *seen* to fail on each relevant topology.
3. Plan capture and normalization (`EXPLAIN (ANALYZE, BUFFERS)`; strip costs/timings/rows/loops, keep
   node types, join order, index names and materialization structure) with a signature per cell.
4. Result JSON with provenance (repo commit, image digest, seed, dataset fingerprint, settings,
   statistics state, prepared-plan policy) separated from interpretation; comparison report showing
   valid/invalid/unsupported/timeout, same-plan/no-resolved-difference, trial and whole-run spread,
   generator caps and prediction ids (P1–P12 of the protocol).
5. Then submit for HIGH review of the full harness (this checkpoint is not a completed task).

## Limits

One host; the DB-backed run is not yet implemented, so the SELECT catalogue has not been executed
against any engine. `DISTINCT ON`/`LATERAL`/`MATERIALIZED` support on the pinned images remains
`unverified` until step 2 runs.
