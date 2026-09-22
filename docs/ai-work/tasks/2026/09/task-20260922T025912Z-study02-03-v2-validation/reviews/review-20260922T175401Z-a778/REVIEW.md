# Review — Study 02/03 v2 HIGH validation

**Task:** `task-20260922T025912Z-study02-03-v2-validation`
**Reviewer:** HIGH session, model `deepseek-flash`, work role `reviewer`
**Reviewed:** `attempt-20260922T173839Z-679172` (Decision Log in `PROGRESS.md`, result in
`RESULT.md`), submitted at commit `0b7dfd8`, branch `repo/study02-03-v2-high-validation`.
**Decision:** **approve for integration.** No corrective handoff needed.

## Decision Log reviewed first

| # | Decision | Verdict |
|---|---|---|
| 1 | Validate against run JSON, not the generated report | Accept. It is the decision that produced the two findings below; the report alone would have repeated the 3b-2 headline. |
| 2 | One cross-run analysis per study, indexed via `inputs:` | Accept. It is efficient, and the author proved the four reports regenerate byte-identically from the recorded images before regenerating, so no published measurement text changed. |
| 3 | Do not edit `claims.json` (v1); write verdicts to a new review file | Accept. Correctly treats v1 as a protected artifact and defers the fix to the correction task. |
| 4 | Record the `claim-07`/`claim-05` defects instead of repairing them in place | Accept. They are derivation defects in signed artifacts; the correction task owns the atomic re-derivation. |
| 5 | Leave the 32-buyer PostgreSQL sign flip unresolved | Accept. Buyer count and design order moved together; resolving it from this data is impossible and the review says so. |
| 6 | No new measurement required | Accept. The task may not start a database; the named follow-ups are optional and cheap. |
| 7 | Retire `claim-gap-04` | Accept. The claim's whole content is the now-false statement that validation is open. |

All seven are local, reversible and within the brief. None changes the scientific question,
acceptance criteria or a signed conclusion.

## Result reviewed

- The two signed analyses are present, carry honest frontmatter (digest, commit, `inputs:`) and
  name their weaknesses. The Study 02 file covers all three phase-3b runs; the Study 03 file
  covers its run and correctly marks `claim-08` out of scope.
- The claim-verdict JSON is well-formed and its verdicts are backed by re-derived JSON figures:
  `claim-06` accepted (1.06-2.35x PG, 1.68-2.21x YB, refund 2.1-2.5x), `claim-09` accepted
  (r06 unanswerable everywhere; 3 102 vs 679 ops/s; r02 = 26.4 ops/s), `claim-07` narrowed
  (1.7-5.7x), `claim-05` narrowed (drop "8x load"; keep +33-35% storage).
- Spot-rechecks by the reviewer: Study 02 3b-2 manifest lists `yb-cluster3/x1_cas_ledger` as
  failed and its JSON holds 18 ledger audits, 2-trial medians — a partial cell, not a gate
  failure, as the review says. 3b-3 X1 cells hold 26 ledger audits each, all zero; report checks
  3/6 per cell, zero failures. Study 03 `l3_section_sharded` `r02` = 26.4 ops/s and `l2` 100k
  seat map = 3.10k ops/s, matching the analysis. The `claim-05` load contradiction is real:
  same loader, 6410 ms vs 1062/1106 ms, empty harness diff.
- `CONTEXT.md` and `LESSONS_LEARNED.md` updates are accurate and small; the two lessons state
  the evidence rather than the narrative.
- `tools/evidence validate` runs clean (`16 claims, 0 errors`); the new review file is additive
  and does not touch the registry.

## Conditions on the outcome

The validation closes the open question and authorizes no measurement. The correcting task must
honour the three required corrections in the verdict file; until it does, `claim-05`, `claim-07`
and `claim-gap-04` must not be used as book inputs.

## Integration

Approve and integrate into local `main`; tag the integrated state
`repo/study02-03-v2-high-validation`.
