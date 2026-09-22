# Review — Typst book toolchain and shell

**Task:** `task-20260922T025912Z-typst-toolchain`
**Reviewer:** HIGH session, model `deepseek-flash`, role `reviewer`
**Reviewed:** `attempt-20260922T185847Z-3a1c2e`, submitted at `f98d4a1`.
**Decision:** **approve for integration.** No corrective handoff.

## Decision Log reviewed first

| # | Decision | Verdict |
|---|---|---|
| 1 | Pin the Typst image by immutable digest | Accept; `book/Containerfile` and the manifest both carry `sha256:032e…`, resolving to Typst 0.15.1. |
| 2 | Use Typst's bundled fonts and prove embedding | Accept; independently re-read `pdffonts`: 3 embedded of 3. |
| 3 | Thin container image (typst + poppler) | Accept; compile and verification share the pinned toolchain. |
| 4 | Chapter/concept files are structural shells only | Accept; no scientific prose was invented, so the HIGH synthesis task is not pre-empted. |
| 5 | Exactly one registry, missing id aborts the build | Accept; `lib/evidence.typ` reads `v2/claims.json` only and panics on an unknown id. |
| 6 | Provenance injected, never hand-set | Accept; values flow from `build.sh` through Typst `--input`; the manifest matches the build output. |
| 7 | Verification measures the artefact | Accept; the two defects found during the task (pdffonts column, `grep -q`/`pipefail`) are fixed and recorded. |
| 8 | Commit the draft and manifest | Accept; with the follow-up fix that stores the PDF byte-for-byte. |

All eight are within the brief and reversible; none affects measurement or another study.

## Independent verification of the committed artefact

The reviewer did **not** re-run `build.sh` (a rebuild changes the build clock and therefore the PDF
hash, which would invalidate the committed manifest). Instead it verified the committed bytes:

```text
manifest pdf_sha256 == git cat-file HEAD:book/dist/data-architecture-reference.pdf | sha256sum
   da433f01436f33238a381d2db2c61f66d9ad887cd4ad2536e4d83a8c1a7379ef  (match)
source_commit 9b04419 is an ancestor of HEAD; manifest source_dirty=false
pdfinfo:     30 pages
pdffonts:    3 embedded of 3 (Libertinus Serif regular/bold, DejaVu Sans Mono)
pdftotext:   registry digest present; 23 of 23 v2 claim ids present
```

The binary-exemption fix (`.gitattributes` `book/dist/**/*.pdf -text`) is what makes the first
check pass; before it, git normalised the PDF and the stored blob did not match the manifest. That
defect was found and corrected inside the task, and recorded in `LESSONS_LEARNED.md`.

## Result reviewed against the brief

- **Exact `book/` structure**: present (`main.typ`, `lib/`, 7 chapters, 8 concepts, `evidence/`,
  `Containerfile`, `README.md`, `dist/`).
- **Pin Typst by digest and fonts**: done; fonts are the image's bundled set, pinned by the digest,
  and embedding is verified.
- **Compile only in Podman; modular `.typ`; docs Markdown; state JSON**: yes.
- **Visual system, title/edition/provenance, TOC, callouts, evidence labels, source links, headers,
  footers, tables, accessible code blocks**: all present; the evidence index renders from the
  registry and each card carries a source pointer.
- **Chapter shells and navigation without invented prose**: confirmed — chapters carry the
  progressive structure and claim hooks only.
- **Draft PDF and manifest from a clean commit; verify pages, fonts, links, evidence digest, hash;
  actual clock/version outputs**: yes, `source_dirty=false` and `built_at_utc` is the real build
  clock.

## Conditions on the outcome

The draft is a verified skeleton, not a finished book; `book-synthesis-v1` owns the prose. A future
Typst upgrade must add a new digest and re-run the build verification. No measurement was performed.

## Integration

Approve and integrate into local `main`; tag the integrated state `repo/book-typst-toolchain-v1`.
