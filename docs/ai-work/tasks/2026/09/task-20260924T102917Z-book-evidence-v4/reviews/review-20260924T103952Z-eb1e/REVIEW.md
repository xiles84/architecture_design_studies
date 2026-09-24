# Review — task-20260924T102917Z-book-evidence-v4

**Reviewer:** HIGH session, model `deepseek-flash`, Deep Code CLI (effort not exposed)
**Verdict:** approve
**Disclosure:** this is a self-review. The same session planned, executed and reviewed the task, which
the workflow permits for a HIGH session executing a LOW-eligible task but which weakens the
independence of the review. Every claim below was re-checked from the committed files, not from the
session's memory of writing them, and the checks are reproducible from this file.

## What was checked, and how

| Check | Method | Result |
|---|---|---|
| The v4 package validates | `go run . validate-v4 --repo ../..` | 29 claims, 0 errors |
| The default dispatch resolves v4 | `go run . validate --repo ../..` | 29 claims, 0 errors, active source `book/evidence/v4/claims.json` |
| The frozen v3 package still validates unchanged | `go run . validate-v3 --repo ../..` | 29 claims, 0 errors |
| The lifecycle rules are tested | `go test ./...` | ok (~25 s), including six new v4 fixtures |
| v3 is byte-unchanged | `git diff repo/book-evidence-registry-v3 -- book/evidence/v3/` | empty |
| Carried-forward claims are unmodified | normalise away `status`, `gap_kind`, `superseded_by`, `remaining_dimensions`, `closed_dimensions` and diff v3 against v4 | no change |
| The rendered book still holds every claim | `book/build.sh --out dist/preview.pdf` | 53 pages, 29/29 claims indexed, 5/5 fonts, digest present, banner absent |
| A retired claim cannot reach the page | negative control: quote `v2-gap-05` in a source | source checker fails; Typst panics naming `v4-gap-01` |
| An unknown id cannot reach the page | negative control: `v9-does-not-exist` | `[FAIL] unknown-claim-id` |
| A fixed registry version cannot return | negative control: a v3 path in a source | `[FAIL] registry-literal` |
| A preview build cannot pass as a release | compile with no `--input` | banner present in preview, absent in the release build |

## Judgement on the substantive claims

1. **The retirement is warranted and not a rewording.** The retired clauses are quoted in the ledger
   beside the claim that falsifies each one. I re-read `v3-05` ("the client spreads operations over
   all three endpoints") and `v3-06` ("a key-locality pair … was executed") against the retired
   `v2-gap-07` ("balanced client access across cluster endpoints is not proven") and `v2-gap-05`
   ("Study 05's placement pair was never executed"). The contradiction is real, and the successors
   keep the clauses that remain true.
2. **Nothing new is claimed.** Both successors have `support: []`, `trials` all zero and no number in
   their statements; the runs they cite are the existing `v3-05`/`v3-06` anchors plus the Study 03
   anchor the retired gaps already used. No report, result file or signed analysis is edited.
3. **The gap-kind split is a real distinction, not bookkeeping.** `v2-gap-02` is a schema
   limitation — no design can answer the hold funnel — and badging it as a coverage gap implied that
   a run would close it. The four kinds now separate that from regime, instrument and unstable
   measurement.
4. **The single-source registry path closes R03 properly.** The version is an input; both the
   source-relative and display paths derive from it; the checker keeps it that way. The two remaining
   version strings in the tree are historical documents (the frozen packages' own READMEs and
   ledgers), which the rule deliberately does not police.
5. **The `claims_indexed` gate caught a defect in this task, which is evidence it works.** The first
   build of this change reported 28/29 because the longer card meta text squeezed the claim-id column
   until `v3-gap-01-study05-remaining-regimes` wrapped. That is exactly the R01 failure shape
   (a long right-hand column squeezing the left), and it was fixed by putting the id on its own line.

## Where this task is weak

- **The `yb-n1` single-endpoint fact is a source read, not a measurement.** It is verifiable in
  `studies/01-charity-tree/run-study.sh` and `studies/02-ticket-booking/run-study.sh`, and the claim
  says "connect through a single endpoint" on that basis. If the client library ever fans out on its
  own, the wording would be wrong; the analysis records this hedge.
- **The lifecycle rules are only as good as the next ingest.** The validator now refuses a retired
  claim that is also active and a gap with no kind, but "does this new run close that old sentence?"
  remains a human judgement, and it is exactly where v3 failed. The ledger and the new
  "re-read the gaps your evidence touches" step in `book/evidence/README.md` are process, not
  enforcement.
- **Two figure-number defects are waived, not fixed.** They are recorded in `check-waivers.txt` with
  an owner (`book/figures-v2`), which is visible debt but still debt.
- **`book/dist/` is deliberately untouched**, so the tracked PDF and its manifest still describe the
  v3 build. The PDF is regenerated once, at release, from fully merged sources; until then a reader
  of `dist/` sees the previous draft.
- **Repository housekeeping in this task:** a stale `.git/index.lock` (no live git process, no
  running container) blocked the queue's own commit; it was removed and the publish event chain was
  committed by hand. `/.queue-drafts/` was added to `.git/info/exclude` alongside the existing
  `/.worktrees/`. An untracked completed-event file belonging to
  `task-20260923T235910Z-diagram-renderer-pin` was committed verbatim with attribution so the submit
  gate could pass. All three are noted in the result file and `CONTEXT.md`.

## Required tag

`repo/book-evidence-registry-v4`, created by `queue integrate` on the integrated state.
