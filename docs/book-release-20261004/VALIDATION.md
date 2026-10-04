# Edition 3 preview 1 — release validation

HIGH, GPT-6 via Codex desktop; selected effort unknown; 2026-10-04.
Self-review of the release is disclosed. Independent review of LOW's correction
is recorded separately in `REVIEW-ORDER.md` and its signed companion.

## Result

| Check | Observed result |
|---|---|
| Canonical clean build | `book/build.sh`, Podman Typst 0.15.1, source `3d43d9720fb45c29fa08249e7bdee8fb669c1fba`, dirty=false |
| PDF | 98 A4 pages; SHA-256 `fc20df64b07549e44c9b287915815049507e649ff5d67cb4a4dcff9cd13db603` |
| Source/figure checks | pass; all 17 pinned figures reproduced; no compiler warnings |
| Semantic evidence resolver | v6: 31 claims, zero errors, Podman Go 1.26 |
| Embedded fonts / indexed claims | 6/6 and 31/31 |
| Provenance | edition and evidence digest present; no unverified banner or unresolved tokens |
| SELECT source checks | containerized vet, unit tests and static build passed; interleaving/late-mismatch and actual RunTag regression tests included |
| Correction inputs | 18 tagged whole runs, 516 valid cells, 2,580 samples; required controls fired; three distinct orders per family/engine |
| Preservation | 120 files, 1,816,351 bytes; every SHA-256 matched both preserved files and committed Git blobs |
| Queue before integration | 88 tasks; zero errors; two existing T2 coordination warnings plus this new claim's pre-integration live-ref warning |

The manifest under `book/dist/` records the complete producing image, evidence,
source and artifact identities. Producing run tags and original signed analyses
were neither edited nor moved. No database was started for release checks.

## Rendered proof

Inspected seven numbered contact sheets covering all 98 pages. Inspected the
reading guide, Rollup diagram, tables, case boxes, checklist, Study 04 profile,
legacy divider and the two shortened legacy passages at larger resolution.

The combined layout split the Rollup cost table. Grouping its heading and five
rows fixed that defect. The final PDF was rendered again; only pages 1, 2, 4, 8
and 9 changed (provenance, contents and the repaired table). All other page PNGs
are byte-identical to the already inspected render. Re-inspected all changed
pages. Representative final proofs are in `proof/`.

Reading guide: page 3. Rollup: pages 8–10. Study 04: page 14. The checklist stays
together; the cost table stays together; diagrams/captions remain legible and
no clipping or overlaps were observed.

## Scope and limits

The owner authorized integration of the reviewed pilot and a preview release.
The remaining Edition 2 chapters are explicitly retained pending conversion;
their long evidence cards are not claimed to meet the final Edition 3 template.
The current edit shortens the Rollup opener, its diagram equivalent and two
dense legacy passages; `book/STYLE.md` now explicitly requires useful visuals
for new or complex concepts. Full conversion remains separate unfinished work.

SELECT correction data are consolidated but no new SELECT timing claim enters
the frozen book evidence package. The signed HIGH qualification keeps noisy or
conflicting aggregate comparisons unresolved. Join/CTE yielded a plan without
running its matrix, so its branch and worktree remain available to resume.

## Cleanup protocol

The recorded pre-release inventory contains 63 branch tips and 18 worktrees,
including main and this release task. Every tip in that inventory is reachable
from the combined release branch. Before any removal, recheck reachability from
main, task status, the benchmark lock and exact source/archive hashes.

Only captured loose/ignored files may be removed after their archive is on main.
Then remove clean completed worktrees through Git, without forcing removal of
dirty trees. Delete only branches reachable from main using `git branch -d`.
Keep the unfinished join/CTE worktree. Final operations are recorded separately
in `CLEANUP-RESULT.json`; new work from a concurrent session is left alone.

The broad restore initially proposed for temporary line-ending normalization
was rejected by automatic approval review because it could discard unrelated
edits. Read-only checks proved zero content differences; a narrowly scoped
restoration of the temporary SELECT manifest line endings succeeded. Portable
LF checkout rules fixed the build without rewriting any scientific evidence.

Release milestone: `repo/data-architecture-book-v3-preview-1`. Local only;
remote synchronization remains with the owner.

## Cleanup observed after integration

`CLEANUP-RESULT.json` and the three `REMOVED-*.json` receipts record 15 completed
worktrees, 60 merged branches and 120 archived originals removed. Every original
branch tip remains reachable from main; the unfinished join/CTE branch and ready
worktree remain. No worktree removal was forced and no tag was deleted or moved.

The post-integration queue audit covers 88 tasks: zero errors and two older
completed-T2 coordination warnings. The release live-ref warning is resolved.
The pilot and SELECT order correction both have completed queue receipts and tags.

After release queue completion, remove its clean checkout and merged branch.
The annotated `repo/book-release-cleanup-20261004` tag records the final verified
counts including that disposal; the current JSON records the observed counts above.
