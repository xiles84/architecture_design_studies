# Progress — version the Study 05 v2 evidence (v3 package)

**Task:** `task-20260923T235907Z-book-evidence-study05-v3`
**Branch:** `repo/book-evidence-study05-v3`
**Attempt:** `attempt-20260924T001226Z-8816b7` (claim `claim-9db53283c12214e2`, epoch 1)
**Capability/role:** HIGH session executing a LOW task / executor (model `deepseek-flash`)
**Required tag:** `repo/book-evidence-registry-v3`

Artifact-only work over committed evidence. No database started, no benchmark run, no new design, no
measurement; the benchmark lock was not taken. A throwaway in-container Typst compile was run to
confirm the book sources still compile against v3 (no `book/dist` write).

## Decision Log

1. **v3 is a consolidated registry, not a delta addendum.** `book/evidence/v3/claims.json` carries
   all 22 still-active v2 claims forward byte-for-byte in content and adds `v3-01`…`v3-06` plus
   `v3-gap-01`. Chosen over a v2+v3 union so the book keeps one active file and every existing
   `v2-XX` citation id stays valid without re-citing every chapter (most chapters are not owned by
   this task). Confidence: High.
2. **`v2-gap-04` is retired by reference; v2 and v1 are untouched.** `git diff
   repo/book-evidence-registry-v2 -- book/evidence/v2/` is empty, and the clause-by-clause
   disposition is the new `SUPERSESSION_LEDGER.md`. Confidence: High.
3. **The validator gained a v3 predecessor set rather than a v2 edit.** `ValidateV2` was refactored
   to a shared `validateRegistry` with an allowed-predecessor set; `ValidateV3` closes supersessions
   and retirements over v1 ∪ v2. `validate` now defaults to v3; `validate-v2`/`v1-diagnostic` keep
   the old behaviour and the v1 12/11/43/36 regression. Confidence: High.
4. **New claims anchor to a package-level signed ingest analysis** (`book/evidence/v3/ANALYSIS.md`)
   rather than to a study analysis. Reason: the study analyses each record one inputs digest, and
   `studies/05-cache-consistency/reports` is a forbidden path for this task, so a package-level
   analysis that lists all eight digests is the only place every new anchor's digest resolves. The
   two signed study analyses are still cited as the authoritative discussion in the ledger.
   Confidence: Medium (a reviewer may prefer per-report study analyses; the digests and cells are
   identical either way).
5. **Equal-total/churn claims cite report-visible values.** The report renders throughput as `11k`
   and `8.0k`; the run JSON's medians are 10,837 and 8,027. Because the resolver demands a token
   that appears verbatim in the cited report, the claims use `11k`/`8.0k`; the precise medians stay
   in the signed analysis. Confidence: High.
6. **The locality claim uses phase-matched report rows** (warm 882 vs 673; cacheable-99 852 vs 520)
   instead of the signed analysis's "882/852/558 vs 673/557/387" triple, which mixes `_99` and
   `_90` phases across the two layouts. This is a citation choice, not a disagreement with the
   analysis's direction. Confidence: High.
7. **`book/build.sh` and `book/main.typ` were moved to v3** although they sit in later tasks' owned
   paths (`book-figure-provenance`, `book-figure-proofset`). Required, not cosmetic: retiring
   `v2-gap-04` makes the build's v2 claims-index check fail (`claims=22/23`) after `evidence.typ`
   loads v3. The change is mechanical, local and reversible; the ledger tells the later tasks to
   preserve it. Confidence: High.
8. **No book build was run.** `book/dist` is a forbidden path for this task and `build.sh` writes
   it. A read-only compile of `main.typ` in the pinned image to a throwaway output path exited 0,
   proving the v3 registry, `evidence.typ` index and `scenarios.typ` cards compile. Confidence: High.

## What was produced

- `book/evidence/v3/{claims.json,schema.json,confounds.json,coverage.json,SUPERSESSION_LEDGER.md,ANALYSIS.md}`
- `book/evidence/README.md` (v3 active, v2/v1 frozen)
- `tools/evidence/internal/evidence/v3.go`, `v3_test.go`; `v2.go` refactored to the shared resolver
- `tools/evidence/main.go` default + `validate-v3`; `tools/evidence/README.md`
- `book/lib/evidence.typ`, `book/chapters/scenarios.typ`, `book/README.md`, `book/build.sh`,
  `book/main.typ`, `CONTEXT.md`

## Residual risks / limitations

- The new claims inherit single-trial, single-host, closed-loop runs; the churn cells have no repeat
  and the hard-expiry counter is 0.
- The placement distribution remains unevidenced; `v3-06` is scoped to engine key locality.
- `book/chapters/evidence-and-reproduction.typ` (owned by `book-figure-proofset`) and the generated
  `book/dist/build-manifest.json` still name v2; both are recorded in the ledger.
- Because v3 copies v2 content, a future change to a frozen v2 file would not propagate; that is the
  intended versioning model.
