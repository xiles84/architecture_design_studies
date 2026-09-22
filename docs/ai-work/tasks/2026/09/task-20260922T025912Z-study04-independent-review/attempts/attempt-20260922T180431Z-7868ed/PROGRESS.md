# Progress — Study 04 independent HIGH review

**Task:** `task-20260922T025912Z-study04-independent-review`
**Branch:** `repo/study04-independent-review`
**Attempt:** `attempt-20260922T180431Z-7868ed` (claim `claim-13fd21bccfee8d8d`, epoch 1)
**Capability/role:** HIGH / reviewer (model `deepseek-flash`, tool/effort not exposed)
**Required tag:** `repo/study04-independent-review`

Artifact-only review. No database started, no measurement run.

## Decision Log

1. **Re-derive from the result JSON, not the generated report.** Every accepted/narrowed figure was
   extracted from `results/20260921T1215Z-small/pg-single/<design>.json` and the design order from
   the manifest. Confidence: High.
2. **Keep the ~1.7x read floor but separate the artefact from its explanation.** The three
   byte-identical reads really span 1.59-1.66x, so the floor stands. The proposed cell-order cause
   does not survive a check against the recorded order (slowest identical-SQL cell is 8th of 10; the
   last two cells beat cells 6-8), so the review records the mechanism as a hypothesis and requires
   the repeated-design control. Confidence: High (order is in the manifest).
3. **`claim-10` is narrowed, not rejected.** Its second and third clauses (93.7% control loss,
   6.85x trigger cost) are re-derived exactly; only the 1.76x clause needs the upper-bound caveat
   the author's own discussion already states. Confidence: High.
4. **No new task for the gaps.** `task-20260922T025912Z-study04-v2-protocol` already owns them and
   depends on this review; the review adds three required controls to that protocol's scope rather
   than creating a duplicate. Confidence: High.
5. **`claims.json` v1 was not edited.** Verdicts live in a new file,
   `book/evidence/reviews/20260922-study04-independent-review.json`, matching the phase-3b review's
   pattern. Confidence: High.
6. **Report regeneration is allowed and explained.** The report carries a `Generated <timestamp>`
   line, so regeneration changes it plus one analyses-index row; verified that nothing else differs.
   Confidence: High.

## What was produced

- `studies/04-configuration-portal/reports/analyses/20260921T1215Z-small--deepseek-flash--2026-09-22.md`
- `book/evidence/reviews/20260922-study04-independent-review.json`
- Regenerated report (index row + generated-at line only)
- `CONTEXT.md` Study 04 section and open-gap 3 updated

## Residual risks / limitations

- One trial per measurement in the source run; only large effects are accepted.
- No independent SQL inspection of every design was performed beyond the statement excerpts quoted
  in the discussion; the review spot-checked the concurrency and rollup paths it judges.
- The v2 protocol additions are requirements recorded here, not yet implemented; they land with
  `task-20260922T025912Z-study04-v2-protocol`.
