# HIGH review — book evidence registry and validators

**Review id:** `review-20260922T141600Z-high`
**Reviewer:** HIGH session, work role `reviewer`. Model/effort not exposed.
**Task:** `task-20260922T025912Z-evidence-registry`, attempt `attempt-20260922T135849Z-578160`,
submitted `event-20260922T140534Z-submitted` (result commit `0b46a05`).
**Verdict:** **accepted for integration.**

## Acceptance criteria

| Criterion | Verdict | Evidence |
|---|---|---|
| `claims.json` validates against a committed schema | Pass | `book/evidence/claims.schema.json`; the validator evaluates the same subset in-tree and the committed registry reports 0 errors. |
| Every initial claim records strength, scope, limits, run, digest, tag, environment, topology, report and analysis where applicable | Pass | All 16 claims carry strength/trials/limits and a provenance block; run tag present for the Study 05 claims and null where a run was untagged. |
| Numeric claims cannot cite failed or invalid cells as winners | Pass | `validate` rejects a winner whose cell status is not `valid` and cross-checks `failed_cells`; fixture `failed-cell-winner.json`. |
| Superseded analyses and stale digests are detected | Pass | Digest must appear in both report and analysis (`stale-digest.json`); supersession must be symmetric and explicit (`superseded-not-marked.json`); `--list-inputs` excludes superseded claims. |
| Cross-study numeric comparisons require an explicit comparability record | Pass | `comparability: cross-study` without a resolving record fails (`cross-study-without-record.json`); no cross-study claim is seeded yet. |
| Scenario transfers are labelled direct evidence or analogy | Pass | Schema requires `transfer`; analogy strength must carry `transfer: analogy`; a numeric claim cannot be an analogy. |
| No new interpretation introduced by the extractor | Pass | `extract` copies analysis front-matter verbatim; each statement is the analysis headline/sentence; gaps are stated as gaps. |

## Checks performed

- `go test ./...` green: the committed registry validates cleanly, all five negative fixtures fail
  for the intended reason, and the mini-schema rejects an out-of-enum value.
- `go run . validate --repo ../..` → **16 claims, 0 errors, 16 book inputs**.
- Every referenced report, analysis, environment page, results directory and run tag resolves;
  each digest appears in both its report and its signed analysis.
- Branch diff is limited to `book/evidence/`, `tools/evidence/` and this task's archive.

## Notes and residual risks

- The seed is the headline claim per current analysis plus four gaps, not every sentence in every
  report; later chapters must add claims rather than quote unregistered numbers. This is stated in
  the review of the result and is the intended workflow.
- `provenance.cells` are curator-recorded design ids; the validator enforces statuses and the
  winner rule but does not re-derive them from result JSON. A future task may cross-check against
  `results/*/*.json`.
- Review performed in the same session that produced the work; the checks are mechanical and
  re-runnable.

## Integration instruction

Fast-forward local `main` to this branch and create the annotated tag
`repo/book-evidence-registry-v1`. Do not push or pull.
