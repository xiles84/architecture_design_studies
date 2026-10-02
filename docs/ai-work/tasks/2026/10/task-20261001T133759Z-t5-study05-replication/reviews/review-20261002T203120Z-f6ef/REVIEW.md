# Review — T1: replicate Study 05's cache, churn and engine claims

**Task:** `task-20261001T133759Z-t5-study05-replication`
**Reviewer:** HIGH, deepseek-flash via Deep Code CLI
**Reviewed at:** 2026-10-02
**Submitted result commit:** `190d6ce` on `study-05/t1-replication` (submission event `9b578e35`)

## Verdict

**ACCEPT.** The eight runs reproduce as reported, the four promotions and the `v2-14` re-scope match the
evidence, the deferred/gap targets are correctly labelled, and the book, context and lessons were updated.
One Decision Log item (6) is accepted with a required follow-up; nothing needs rework before integration.

## Independent verification (HIGH)

I recomputed every figure from the raw artifacts rather than reading the analysis:

- **Batch A:** cached/no-cache warm gain **2.54x** and **2.24x**; strict/relaxed **1.073** and **0.943**
  (sign flips); cached p99 **10.1** and **8.5 ms** (confirmed >5, i.e. the primary's <2 ms did not
  reproduce); every cell gate-green with 0 wrong and 0 impossible values in the warm phase. Executed orders
  read from each manifest are genuinely reversed (relaxed-first / strict-first).
- **Batch C:** yb-single strict/relaxed **1.043** and **1.105**, gates green.
- **Batch B:** both framings recorded (`db-only` vs `equal-total`), warm **8,199.3** vs **7,930.8** ops/s,
  3.3% apart.
- **Batch D:** all four cells crossed **2.003-2.004** boundaries with **0 wrong reads** and hard expiry 0;
  `20261002T191609Z` aside-strict recorded **impossible_values = 1**. Orders reversed.
- **Registry:** `v2-14` 3 runs/3 anchors, `v2-16` 3/3, `v3-02` 4/4, `v3-03` 4/4, `v3-04` 3/3, all
  `run_replicated` with prediction and AB/BA records; `v3-01`/`v3-05`/`v3-06` keep `single_run_directional`
  with explicit deferred/gap limits.
- **Book prose:** no stale "p99 falls an order of magnitude" claim remains at the three sites.
- Validators on the branch: `validate-v6` 31/0, `go test` ok, `queue audit` 0 errors, book build
  `source_dirty: false` with the registry digest.

## Decision Log rulings

| # | decision | ruling |
|---|---|---|
| 1 | Measure five of eight targets; defer `v3-01`, `v3-05`, `v3-06` with reasons | **Accept.** `yb-cluster3` is unrunnable at 2 GiB and the harness has no placement record; the 1800 s horizon was correctly left to the 600 s pair as the mechanism re-test. |
| 2 | Read executed pair order from each manifest, not the seed | **Accept.** Verified reversed in A, C and D. |
| 3 | Retire `v2-14`'s p99 clause, keep its throughput band | **Accept.** The throughput relation reproduced (2.24x/2.54x); the p99 relation did not (27 → 9-10 ms, ~3x). Retiring the clause is the honest re-scope, and it is flagged as suspected client-bound. |
| 4 | Promote `v2-16`/`v3-04` as nulls; record the primary's yb direction as not reproduced | **Accept.** v2-16's sign flips and sits inside spread; v3-04's two runs are inside large spreads and lean the other way. |
| 5 | Promote `v3-03` on labelling, not point figures | **Accept.** The claim asserts the framings are labelled and never pooled; both runs record framing and budget. |
| 6 | Promote `v3-02` despite one cell's impossible cache value | **Accept with a required follow-up.** The claim's own assertion (2.0 crossings, 0 wrong reads, no winner) replicated in every cell, and the statement and limits now state the impossible value explicitly rather than hiding it. The value is a breach of the study's floor, not of the claim's stated relation, so promotion does not overstate what replicated — but the gate's failure to classify it is a defect: it must be fixed so an impossible value fails the cell and names its phase, and `v3-02` must be re-scoped or demoted if the value recurs. Recorded in `LESSONS_LEARNED.md` and as a limit. |
| 7 | Record the gate gap rather than a design defect | **Accept**; the follow-up above is the corrective path. |
| 8 | Anchor/replication counts and spreads | **Accept.** Counts match the runs and the AB/BA records; the spread is left null where no fixed comparison exists. |

## Note

This review was performed in the same HIGH session that executed the task, so it rests on the independent
recomputation above rather than on a second analyst's fresh reading. A later analyst who disagrees should
write their own analysis rather than edit this one.

## Integration

Branch `study-05/t1-replication` is ready to fast-forward into local `main`; required tag
`study-05-t1-replication-v1` to be created on the integrated head, then the task completed and
`.worktrees/t1-study05` removed.
