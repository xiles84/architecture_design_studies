# Progress — Typst book toolchain

**Task:** `task-20260922T025912Z-typst-toolchain`
**Branch:** `repo/book-typst-toolchain`
**Attempt:** `attempt-20260922T185847Z-3a1c2e` (claim `claim-7a5c4d857c8d8e1f`, epoch 1)
**Capability/role:** HIGH session executing a LOW task / executor (model `deepseek-flash`)
**Required tag:** `repo/book-typst-toolchain-v1`

Artifact-only work: no database, no benchmark, no measurement.

## Decision Log

1. **Pin the official Typst image by immutable digest, not a floating tag.**
   `ghcr.io/typst/typst@sha256:032e292249bcd378480cc7c142cfa324b63ef8aadeb88d7e7230320c4c9c422f`
   (resolves to Typst 0.15.1). The digest fixes version, bundled fonts and rendering together; the
   build manifest records it. Confidence: High.
2. **Use Typst's bundled fonts, and prove embedding in the build.** The image ships Libertinus
   Serif and DejaVu Sans Mono internally, so no host or extra pinned font package is needed; the
   build fails if any font is not embedded. Confidence: High.
3. **`book/Containerfile` is a thin build image** (typst base + poppler-utils/file) so compile and
   verification run in one pinned toolchain. Confidence: High.
4. **Chapter and concept files are structural shells only.** They carry the fixed progressive
   structure, callout vocabulary and evidence-card hooks; the prose is written by
   `book-synthesis-v1`. Inventing prose here would pre-empt the HIGH synthesis task.
   Confidence: High.
5. **The book reads exactly one registry** (`book/evidence/v2/claims.json`) via `lib/evidence.typ`;
   a missing claim id aborts the build, and v1 is not referenced. Confidence: High.
6. **Provenance is injected, never hand-written.** `build.sh` passes commit, describe, dirty state,
   build clock, Typst version/digest, evidence digest and source-tree hash through Typst
   `--input`; the manifest writes the values the run measured. Confidence: High.
7. **Verification measures the artefact.** Page count (`pdfinfo`), embedded fonts (`pdffonts`, the
   stable `yes yes yes` marker), all claim ids present in `pdftotext`, the registry digest present
   in the text, and URI-annotation count; the build exits non-zero if any check fails.
   Confidence: High. (The first two attempts got the font column and the `grep -q`/`pipefail`
   presence test wrong; both are fixed and recorded in `LESSONS_LEARNED.md`.)
8. **Record the draft as `dist/data-architecture-reference.pdf` + `dist/build-manifest.json`,
   committed.** The task asks for a generated draft and its manifest; the manifest names the exact
   source commit, so the binary is reproducible from it. Confidence: High.

## Notes

- The build clock and git describe change between builds, so two drafts are not byte-identical;
  the manifest's `built_at_utc` explains the difference, and `pdf_sha256` identifies the artefact.
- `book/build.sh` needs no benchmark lock: it starts no database and measures nothing.

## Residual risks

- The chapter prose does not exist yet; the PDF is a verified skeleton.
- Fonts depend on the pinned Typst image; a future Typst upgrade needs a new digest and a re-run of
  the build verification, not an in-place tag move.
- `source_tree_hash_sha256` covers `.typ` sources only; toolchain changes are captured by the image
  digest and the recorded build command.
