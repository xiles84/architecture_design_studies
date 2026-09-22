# Progress — Study 05 independent HIGH review

**Task:** `task-20260922T025912Z-study05-independent-review`
**Branch:** `repo/study05-independent-review`
**Attempt:** `attempt-20260922T182234Z-b46fb1` (claim `claim-ee404a0092cd7e60`, epoch 1)
**Capability/role:** HIGH / reviewer (model `deepseek-flash`, tool/effort not exposed)
**Required tag:** `repo/study05-independent-review`

Artifact-only review. No database started, no measurement run.

## Decision Log

1. **Re-derive every verdict figure from the result JSON.** `warm_reads`, `wrong_reads` and
   `instances` were extracted per scenario and compared with the report, analysis and discussion.
   This is what caught the identity swap. Confidence: High.
2. **The identity swap is recorded as a material correction, not applied to the signed analysis.**
   The analysis is signed and immutable; the book/derived-claim fix belongs to the evidence
   correction task, and the reviewer's own signed analysis states the corrected mapping. Confidence:
   High (report line and JSON agree).
3. **`claim-11` and `claim-12` are accepted with explicit scope, not broadened.** The throughput
   gain holds only for warm cacheable reads on one small node in `db-only` framing; the strict
   result is a read-throughput non-finding, not "strict is free". Confidence: High.
4. **Redis cache/authoritative-store boundary made a hard scope requirement.** No cell measures
   Redis as a system of record; the book must not present these numbers as a database comparison.
   Confidence: High.
5. **No new task for the gaps.** `task-20260922T025912Z-study05-v2-protocol` owns them and depends
   on this review; four requirements are added to its scope. Confidence: High.
6. **`claims.json` v1 not edited.** Verdicts go to a new file, matching the phase-3b and Study 04
   review pattern. Confidence: High.
7. **The partially-superseded discussion is cited only for mechanism.** Its retired residual-failure
   sections are not used. Confidence: High.

## What was produced

- `studies/05-cache-consistency/reports/analyses/20260921-cache-consistency-allgreen--deepseek-flash--2026-09-22.md`
- `book/evidence/reviews/20260922-study05-independent-review.json`
- Report regenerated (index row + generated-at line only)
- `CONTEXT.md` Study 05 section and open-gap 4 updated; one `LESSONS_LEARNED.md` entry on
  cross-cell identity swaps

## Residual risks / limitations

- One trial per measurement in the source run; only the 87%-vs-0% separation is safely outside the
  ~20% noise floor.
- The "86% in a second run" figure in the all-green analysis was not re-derived here; the correction
  task should check it against its own run before quoting.
- The three-instance generality (second machine, repeated trials) is unproven and is left to the v2
  protocol.
