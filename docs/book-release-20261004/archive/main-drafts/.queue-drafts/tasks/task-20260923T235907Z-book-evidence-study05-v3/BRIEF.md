# Execution brief — Version the Study 05 v2 evidence and correct the stale coverage gap

Read `docs/ai-work/brainstorms/records/2026/09/brainstorm-20260923T172022Z-cache-strategy-coverage-by-schema-control/CONCLUSION.md`
first, then `book/evidence/README.md`, `studies/05-cache-consistency/HANDOFF-AMENDMENT-01-V2.md`, and the
2026-09-23 signed analyses under `studies/05-cache-consistency/reports/analyses/`.

## Goal

The active evidence registry (`book/evidence/v2/`) was written before the 2026-09-23 Study 05 runs. Its
`v2-gap-04` and the book sentence that repeats it still call real-TTL churn, repeated trials and cache
resource accounting unmeasured, but those runs are committed. Produce the next versioned evidence package
that ingests them and supersedes the stale statements — without editing any earlier signed artifact.

## Decisions already made (do not re-litigate)

- The new package is a **new version**; v2 stays frozen historical reference. Never edit a v1/v2 claim or digest.
- Supersede by reference, exactly as the v1→v2 correction did.
- The newer runs are **mechanism and scoped-duration evidence**, not a general cache result: the retained
  limits are hard expiry unexercised (`expired_hard = 0`), medium scale beyond the harness maximum
  (3,000 people), placement distribution absent, single host, single run per arm, open loop not run.
- Every number must resolve to a report cell and a run tag. A number that cannot be resolved is a gap, not a claim.

## Acceptance criteria

- New versioned evidence package; each new numeric claim resolves to a report cell, inputs digest and producing tag.
- Stale parts of `v2-gap-04` superseded by reference; the open gaps above retained explicitly.
- `tools/evidence validate` reports 0 errors with the new version as default.
- `book/chapters/scenarios.typ` and `book/README.md` corrected by supersession (the active-source pointer updated).
- All v2 material remains readable; no signed analysis, report or digest edited.

## Constraints

- No measurement: this task runs nothing, takes no benchmark lock.
- Do not edit another analyst's report or analysis; write your own supersession ledger and signed analysis.
- Do not push or pull. Commit explicit paths, tag `repo/book-evidence-registry-v3`, merge through the queue lifecycle.
- Record a Decision Log entry for any judgement about which v2 statement is superseded versus retained.

## Escalate only if

- a numeric claim cannot be resolved to a report cell and the report itself is ambiguous (do not invent a value);
- superseding a v2 statement would silently change a released book number (that is a material defect, not a local choice).

NEXT MODEL: LOW
