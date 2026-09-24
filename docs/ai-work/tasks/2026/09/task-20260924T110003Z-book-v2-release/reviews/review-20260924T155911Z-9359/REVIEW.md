# Review — task-20260924T110003Z-book-v2-release

**Reviewer:** HIGH session, model `deepseek-flash`, Deep Code CLI (effort not exposed)
**Verdict:** approve
**Disclosure:** self-review. The checks below were run against the committed PDF, its manifest and the
installed tooling, not from memory.

## What was checked

| Check | Method | Result |
|---|---|---|
| The file matches its manifest | `sha256sum` against `pdf_sha256` | equal |
| The manifest describes this PDF | edition, registry, figures, release flag | "Edition 2 (draft)", `book/evidence/v5/claims.json`, 8 figures, `release: true` |
| The edition is really on the page | `edition_in_pdf` gate and an independent text search | true; "Edition 2" ×54, "Edition 1" ×0 |
| The artefact carries every claim | `claims_indexed` | 29/29 |
| No preview artefact would pass as a release | banner gate | absent |
| No template token reached the page | unresolved-token gate | false |
| Fonts are embedded | font gate | 6/6 |
| The book's own checks pass | `check.sh`, `export.sh --check` | OK, no waivers; eight figures verified |
| The tracked tree is what was built | `git status` after commit | clean; the PDF is committed |

## Judgement

1. **This is the artefact the reviewer asked for, and the first one that is not a preview.** Every previous
   external verdict in this project was taken on a hand-built draft; the release path now produces a PDF
   whose pages 1–2 carry real commit, describe, dirty state, clock, image digest, evidence registry digest
   and source-tree hash, and whose manifest says which registry and which figures it was built from.
2. **Both blockers found here were in the manifest, not the page.** That is worth naming: the project has
   spent three tasks making the *page* self-checking, and the file that a reader trusts to describe the
   page had no check at all. It has one now, and it fires if the edition never renders.
3. **The `source_revision` fix removes a class rather than an instance.** The field failed twice in
   opposite directions — first recording the previous commit's blob, then recording `uncommitted` before
   its own commit. Excluding a repository-state field from a content comparison is the general fix.
4. **Replacing the Edition 1 PDF rather than adding a second file is right.** `dist/` is the current
   artefact; the old state is in history and under two tags, and a folder holding both would invite a
   reader to quote the wrong one.

## Where this task is weak

- **`pdf_sha256` describes one build, not a reproducible one.** The build clock is an input, so a rebuild
  produces a different file with identical inputs otherwise. The manifest says so, and the source-tree and
  evidence digests are the stable identifiers.
- **`source_dirty: true` in this manifest** is honest but unattractive: the release commit is the one that
  carries the PDF, so the tree was dirty when it was compiled. A later release from a clean tree would
  report false; a truly clean record would need a two-commit dance.
- **The `(draft)` label stays** until an external reader signs off, which is deliberate and means this is
  not a "final" edition.
- **Tranche 2 of the visuals remains unimplemented**, and six further figures would change the page count
  again — a decision for whoever reads the release.
- **The reviewer's verdict is necessarily absent.** This is self-review of the artefact that exists to
  receive an external one.

## Required tag

`repo/data-architecture-book-v2`, created by `queue integrate` on the integrated state.
