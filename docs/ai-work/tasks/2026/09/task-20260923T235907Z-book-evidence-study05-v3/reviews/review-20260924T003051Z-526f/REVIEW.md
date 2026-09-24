# Review — book evidence registry v3

**Task:** `task-20260923T235907Z-book-evidence-study05-v3`
**Reviewer:** HIGH session, model `deepseek-flash` (Deep Code CLI), capability HIGH, role reviewer
**Review id:** `review-20260924T0030Z-v3` (recorded by the queue)
**Verdict:** **approved for integration** — no rework required
**Required tag:** `repo/book-evidence-registry-v3`

## Decision Log reviewed first

| # | LOW decision | Verdict |
|---|---|---|
| 1 | v3 is a consolidated registry (v2 claims carried forward) rather than a delta | accepted — keeps every cited `v2-XX` id valid; a delta would have required editing chapters this task does not own |
| 2 | `v2-gap-04` retired by reference; v1/v2 untouched | accepted — `git diff repo/book-evidence-registry-v2 -- book/evidence/v2/` is empty |
| 3 | shared `validateRegistry` + `ValidateV3` over v1 ∪ v2 | accepted — v2 semantics preserved; v1 regression still 12/11/43/36 |
| 4 | new claims anchor to a package-level signed ingest analysis | accepted with a recorded limitation — `studies/05-cache-consistency/reports` is a forbidden path here; the study analyses are still cited as authoritative |
| 5 | claims cite report-visible `11k`/`8.0k`, not JSON medians | accepted — forced by the resolver, and the medians stay in the study analysis |
| 6 | locality claim uses phase-matched report rows, not the analysis's mixed triple | accepted — correct reading of the report; no disagreement with the analysis's direction |
| 7 | `book/build.sh` and `book/main.typ` moved to v3 although other tasks own them | accepted — retiring `v2-gap-04` would otherwise make the build fail `claims=22/23`; mechanical and reversible, and the ledger tells the later tasks to keep it |
| 8 | no book build run; throwaway in-container compile instead | accepted — `book/dist` is forbidden here |

No decision materially degraded correctness, requirements, architecture or the result.

## Independent checks re-run by the reviewer

```text
go vet ./...                    clean
go test ./...                   ok
validate (default v3)           evidence: 29 claims, 0 errors, 29 book inputs
validate-v2                     evidence-v2: 23 claims, 0 errors
v1-diagnostic                   12 numeric claims, 11 with a missing name, 36/43 names unresolved
v2 unchanged vs its tag         empty diff
carried-forward claims          deep-compare v2 vs v3 (except v2-gap-04) -> []
run tags cited by v3-01..v3-06  all exist
book compile (pinned image)     exit 0, no book/dist write
```

Spot checks against the acceptance criteria:

- Every new numeric claim has a structured support key whose `resolves_to` token is present verbatim
  in its primary report, and the report digest appears in `book/evidence/v3/ANALYSIS.md`; the
  validator enforces this and reports 0 errors.
- `v2-gap-04` is not reproduced and is listed in `retired_predecessor_claims`; the clause-by-clause
  disposition (superseded vs retained) is in `SUPERSESSION_LEDGER.md`.
- `book/chapters/scenarios.typ` no longer says real-TTL churn, repeated trials or cache resource
  accounting are unmeasured; it retains medium scale, placement distribution, hard expiry and
  open-loop demand as gaps. `book/README.md` names v3 as active.
- No earlier analysis, report or digest was edited; v1 and v2 are byte-identical to their tags.

## Findings (all non-blocking, already recorded by the executor)

1. `book/chapters/evidence-and-reproduction.typ` (owned by `book-figure-proofset`) still names the v2
   registry path in prose; `book/dist/build-manifest.json` is the generated manifest of the last
   pre-v3 build. Both are named in the ledger and must be updated by their owning tasks.
2. The anchor analysis for the new claims is package-level. Acceptable, but a later study-analysis
   could supersede it.
3. The new claims are single-run; hard-expiry is still unexercised — correctly retained as a gap.

## Verdict

Approved for integration. The executor's Decision Log is honest, the acceptance criteria are met,
and the one ownership-boundary decision (7) was necessary to keep the book build coherent.
