# Review — T1: replicate Study 04's concurrency claims (`v2-12`, `v2-17`)

**Task:** `task-20261001T133759Z-t4-study04-replication`
**Reviewer:** HIGH, deepseek-flash via Deep Code CLI (effort not exposed)
**Reviewed at:** 2026-10-02
**Submitted result commit:** `9f6d350` on `study-04/t1-replication` (work commit `3adc5d3`, RESULT `4b0f4c1`)

## Verdict

**ACCEPT.** The two claims' fresh evidence reproduces, the re-scoped statements match what
replicated, the prediction scoring is honest, and no correction is needed. The only
integration-time work is one item the LOW left to the scratch path: the tracked book draft
is rebuilt on the branch. The verdict is recorded on the branch before the merge.

## Independent verification (HIGH)

I re-derived every figure from the raw artifacts rather than reading the analysis:

- **All 14 per-run cells reproduce exactly** from
  `studies/04-configuration-portal/results/<run-id>/pg-single/*.json`: writers, `c1`/`c2`
  ops/s, the `c2`/`c1` ratio to three decimals, `c1` conflicts, and the `x1` loss
  percentage, with each run's `gate` and `audit` passed. The executed pair order was read
  from each run's `manifest.yaml` `design_order_pg-single`, not from the analysis's labels.
- **The order control is real.** Four `c1`-first and seven `c2`-first runs now exist; the
  six added runs' manifests put `c1_optimistic_version` before `c2_pessimistic_lock`, which
  the first eight runs never did. Every figure in the analysis table matches.
- **Predictions score correctly.** P2 (1.637x, 1.573x ≥ 1.4x), P3 (3.25% lock-first and
  1.12% version-first, both < 10%), P4 (1,328 → 569 ops/s, below half), and the P1 falsifier
  firing once (3.849x) all recompute; P5's control loss 50.0% / 93.7%-98.0% recomputes.
- **The central correction holds.** `20261002T162737Z-w1` measured `c1` = 306 ops/s while
  `20261002T162400Z-w1` measured 1,286 at the same level, design and design order; the
  version-check-first runs read 0.948x and 0.905x. Client and database cgroup deltas record
  zero throttled periods. Calling it an unexplained single-run excursion is right.
- **Registry counts and bands.** `v2-12`: runs=whole_run_replications=5, anchors=5,
  band 1.56x-1.76x. `v2-17`: runs=whole_run_replications=23, anchors=23, control_coverage
  present with 1 fired. Both `run_replicated`.
- **Validators:** `tools/evidence validate-v6` 0 errors; `go test ./tools/evidence/...` ok;
  `queue audit` 0 errors (the pre-existing `t2` live-ref warnings only; this task's
  `live_ref_drift` cleared on submit).
- **Book:** the tracked draft built from the branch tree with `source_dirty: false`, evidence
  digest `8d78efca…` matching `book/evidence/v6/claims.json`, 88 pages, 6/6 fonts, 31/31
  claims, no unverified banner, no unresolved tokens.

## Decision Log rulings

| # | LOW decision | Ruling |
|---|---|---|
| 1 | Add six runs to create a genuine reversed order after the manifests showed all eight earlier runs ran `c2` first | **Accept — this is the correction that makes the result supportable.** The re-run is faithful to the committed protocol and the executed order is verified per manifest. |
| 2 | Replace the first draft's "order effect" reading with an unexplained single-run excursion | **Accept.** Independently reproduced; the alternates (order, throttling, a slower code path) are excluded by the same-order reading, the zero cgroup throttling and the unchanged SQL. Reporting rather than discarding the outlier is correct. |
| 3 | Score `v2-17`'s prediction `falsified` (P1's named falsifier fired once) while the claim is `run_replicated` | **Accept as honest.** The mismatch between prediction verdict and claim strength is the registry-wide convention question T2's review already raised, not this task's to resolve. Flagged in the RESULT; leave for a registry-level decision. |
| 4 | Run the 1-writer level with `c1`/`c2` only and report the three `x1` cells as not applicable | **Accept.** A lost update is structurally impossible at one writer, the gate is right to refuse, and the limits name the failed cells rather than hiding them. |
| 5 | Compare the retry budget only within a matched design order | **Accept.** Removes the first draft's order confound; 3.25% and 1.12% recompute. |
| 6 | Un-withhold `v2-12`'s attribution and set `v2-17`'s control coverage present | **Accept.** Both follow directly from 1, 4 and 5. |
| 7 | Counting conventions and the two spread figures | **Accept, with a negligible note.** `v2-12`'s 12.4% and `v2-17`'s 8.3% are computed over the 16-writer ratios; the exact `v2-17` spread from unrounded values is 8.35% (8.4%), a rounding-basis difference of 0.1 points in an auxiliary field. Not material; no rework. |

## Integration-time completion (not rework)

- **Tracked book draft rebuilt on the branch** at `39c32dd` (`book: rebuild on the re-scoped
  concurrency claims`). The task's prose change makes the committed draft stale, and every
  recent T1 task rebuilt it in the same change; the LOW had built only to a scratch path.
  No measurement, claim or LIMIT changed: the rebuild only re-embeds the re-scoped
  statements and corrects the tracked manifest's evidence digest.

## Notes carried forward (non-blocking)

- The environment reported an owner edit to the analysis file. The file is byte-identical to
  the committed blob (`sha256 d36042c0…`), so no owner delta is present in the branch tree;
  if an edit was intended, it did not reach the worktree and the owner should say so.
- The required tag `study-04-t1-replication-v1` must be created on the integrated `main`
  head, and the worktree removed.
- Remaining scientific gaps stated in the analysis: one fresh-load trial per cell; the 32-,
  8- and 4-writer levels still primary-only; `x1` cannot fire at one writer.
