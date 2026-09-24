---
analysis_id: book-evidence-v3-ingest--deepseek-flash--2026-09-23
run_id: 20260923T1215Z-res2-placement-yb3
environment: host-zenbook-ux5406sa
analyst: deepseek-flash
analyst_kind: ai
analyst_version: "deepseek-flash (DeepSeek), HIGH-capability, Deep Code CLI; effort setting not exposed. Executor of task-20260923T235907Z-book-evidence-study05-v3."
analyzed_at: 2026-09-23
repo_commit: edded4df1a2245764fcdbe2c327a6484fd0289f5
supersedes:
status: current
headline: Versioned ingest of the 2026-09-23 Study 05 cells into a v3 evidence package; every new number resolves to a report cell, inputs digest and producing run tag, and the stale v2-gap-04 clauses are superseded by reference.
---

# Analysis — book evidence v3 ingest of the 2026-09-23 Study 05 runs — deepseek-flash

> Deliverable of `task-20260923T235907Z-book-evidence-study05-v3`, implementing the evidence-update
> clause of `brainstorm-20260923T172022Z-cache-strategy-coverage-by-schema-control`. This is an
> ingestion and attribution analysis: it measures nothing and edits no earlier signed artifact.

## TL;DR

- The active registry moves from `book/evidence/v2/claims.json` to `book/evidence/v3/claims.json`.
  v3 carries every still-active v2 claim forward **unchanged** and adds six numeric Study 05 claims
  and one successor gap claim.
- The stale clauses of `v2-gap-04-study05-unmeasured-regimes` are superseded by reference: the new
  evidence measures real-TTL churn duration, repeated in-run trials, equal-total cache accounting and
  YugabyteDB 1-node/3-node cells. The still-open regimes are retained explicitly in
  `v3-gap-01-study05-remaining-regimes`.
- Every new numeric claim resolves a structured support key to a token in its cited report, and the
  report's `inputs digest` appears in both that report and this analysis.

## 1. Cells ingested

The 2026-09-23 Study 05 work produced three churn cells and eight resource/engine cells. Each row
below is one cited report; the digest is the report's own inputs digest and the tag is the producing
run tag.

| Run | Inputs digest | Run tag | Cells | Report |
|---|---|---|---|---|
| `20260923T1020Z-churn600-through` | `75af94d4c846e528` | `run/05-cache-consistency/20260923T1020Z-churn600-through` | 1 | `studies/05-cache-consistency/reports/20260923T1020Z-churn600-through.md` |
| `20260923T1040Z-churn600-aside` | `cf16e598ad9b32f8` | `run/05-cache-consistency/20260923T1040Z-churn600-aside` | 1 | `studies/05-cache-consistency/reports/20260923T1040Z-churn600-aside.md` |
| `20260923T1100Z-churn1800-through` | `8ad59fb36a65c17f` | `run/05-cache-consistency/20260923T1100Z-churn1800-through` | 1 | `studies/05-cache-consistency/reports/20260923T1100Z-churn1800-through.md` |
| `20260923T1148Z-res2-dbonly-pg` | `b28d8dd331c4c968` | `run/05-cache-consistency/20260923T1148Z-res2-dbonly-pg` | 1 | `studies/05-cache-consistency/reports/20260923T1148Z-res2-dbonly-pg.md` |
| `20260923T1150Z-res2-equaltotal-pg` | `e0dfc55f3c363944` | `run/05-cache-consistency/20260923T1150Z-res2-equaltotal-pg` | 1 | `studies/05-cache-consistency/reports/20260923T1150Z-res2-equaltotal-pg.md` |
| `20260923T1152Z-res2-yb1` | `6f04a8faf4a428bf` | `run/05-cache-consistency/20260923T1152Z-res2-yb1` | 2 | `studies/05-cache-consistency/reports/20260923T1152Z-res2-yb1.md` |
| `20260923T1157Z-res2-yb3` | `f2133069805cf82d` | `run/05-cache-consistency/20260923T1157Z-res2-yb3` | 4 | `studies/05-cache-consistency/reports/20260923T1157Z-res2-yb3.md` |
| `20260923T1215Z-res2-placement-yb3` | `a18817d592f767c5` | `run/05-cache-consistency/20260923T1215Z-res2-placement-yb3` | 2 | `studies/05-cache-consistency/reports/20260923T1215Z-res2-placement-yb3.md` |

The two signed study analyses this ingest cites are
`studies/05-cache-consistency/reports/analyses/20260923T1100Z-churn--deepseek-flash--2026-09-23.md`
and
`studies/05-cache-consistency/reports/analyses/20260923T1215Z-resources-engines--deepseek-flash--2026-09-23.md`.
They are not edited; this file is an independent ingest and may be cited as the evidence anchor's
analysis for the new claims.

## 2. What the new claims state, and what they refuse to state

- **Churn (v3-01, v3-02).** Duration and wrong-read counts are claimed. The 1801 s cell crossed 6.0
  real 300 s boundaries with 0 wrong reads in 2,879,077 reads; the two 601 s cells crossed 2.0 each
  with 0 wrong reads. The hard-expiry path fired zero times in every cell, so "an entry reached its
  hard TTL" is **not** claimed — it is retained in `v3-gap-01`.
- **Equal-total (v3-03).** The arm is claimed as a *label and separation*, not as a throughput
  ratio: the same design and topology read 11k ops/s under equal-total framing and 8.0k under
  db-only, within the run's spread. The 2.2–3.5x add-cache claim (`v2-14`) is untouched and keeps
  `conf-04`.
- **Engines (v3-04, v3-05).** Each engine's strict/relaxed app-99 pair is reported as one run each.
  The directions differ (1-node relaxed faster, 3-node strict faster) and neither is attributed;
  this replaces "no YugabyteDB cell for Study 05" with evidence, not with a ranking.
- **Locality (v3-06).** The colocated/non-colocated pair is measured (882 vs 673 warm, 852 vs 520
  cacheable-99) but the tablet/leader distribution was not produced, so the claim is scoped to
  engine key locality on one host and never to verified colocation.
- **Remaining gaps (v3-gap-01).** Hard TTL expiry unexercised, medium scale beyond the 3,000-person
  harness maximum, absent placement distribution and unrun open-loop demand are carried forward
  explicitly.

## 3. Supersession, not editing

`v2-gap-04-study05-unmeasured-regimes` is **not** reproduced in v3 and is not edited in v2. Its
retirement and the clause-by-clause mapping are recorded in
[`SUPERSESSION_LEDGER.md`](SUPERSESSION_LEDGER.md). `book/chapters/scenarios.typ` and
`book/README.md` are corrected by reference to the new claims.

## 4. Where this ingest is weak

1. **It re-authors nothing.** The v3 file copies the v2 claims byte-for-byte in content; a v2 claim
   whose limits have since been partly closed elsewhere is carried with those limits unchanged. The
   only v2 claim whose content is not carried forward is `v2-gap-04`.
2. **The new claims inherit their runs' single-trial, single-host, closed-loop character.** In
   particular the churn cells are single sustained measurements (no repeat), and the placement
   distribution is absent.
3. **The anchor analysis for the new claims is this package-level ingest**, not a study analysis.
   The two signed study analyses remain the authoritative discussion of their runs; this file maps
   their digests into the registry.
4. **The book is not rebuilt here.** `book/lib/evidence.typ`, `book/build.sh` and `book/main.typ` all
   name v3 now, so the build's claims-index check and the rendered index agree; but `book/dist/` is
   untouched and no build was run, so the current PDF still records the pre-v2-gap-04 retirement
   state. Rebuilding the PDF belongs to `task-20260923T235919Z-book-cache-decision-map` and the two
   figure tasks.
