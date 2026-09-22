# Result — Typst book toolchain and shell

**Task:** `task-20260922T025912Z-typst-toolchain`
**Branch:** `repo/book-typst-toolchain`
**Required tag:** `repo/book-typst-toolchain-v1`
**Capability:** HIGH session executing a LOW task (model `deepseek-flash`)
**Benchmark:** none — documentation build only

## What was delivered

1. **The `book/` structure from the goal**: `main.typ`, `lib/config.typ`, `lib/evidence.typ`,
   7 chapter shells under `chapters/`, 8 concept shells under `concepts/`, `Containerfile`,
   `build.sh`, `README.md`, and `dist/` with the generated PDF and manifest.
2. **Visual system** (`lib/config.typ`): A4 geometry, running header/footer, a fixed type scale,
   `direct`/`analogy`/`gap` callouts, family/variant/topology/scenario markers, an evidence-card
   component, accessible code blocks (monospace, no hyphenation, labelled surface) and the
   progressive-structure helper.
3. **Title, edition and provenance pages** plus a table of contents and an automatically rendered
   evidence index. Every provenance field is injected via Typst `--input` from the build script —
   nothing is hand-set.
4. **One evidence source**: `lib/evidence.typ` loads `book/evidence/v2/claims.json`; a claim id that
   does not exist aborts the build. The index renders all 23 active claims with strength, trials,
   confounds, limits and a source pointer.
5. **Pinned toolchain**: `book/Containerfile` from
   `ghcr.io/typst/typst@sha256:032e292249bcd378480cc7c142cfa324b63ef8aadeb88d7e7230320c4c9c422f`
   (Typst 0.15.1) plus poppler for verification. Typst's bundled Libertinus Serif and DejaVu Sans
   Mono are embedded in the PDF; no host font is needed.
6. **`book/build.sh`**: builds the image, compiles, verifies and writes `dist/build-manifest.json`
   recording source commit, describe, dirty state, build clock, Typst version/digest, evidence
   registry digest, source-tree hash, PDF sha256, pages, fonts, claims indexed and URI annotations.
7. **A generated draft** `book/dist/data-architecture-reference.pdf` from a clean commit, committed
   with its manifest. **No scientific prose was invented**; chapters and concepts are structural
   shells for `book-synthesis-v1`.

## Verified build (all checks in the pinned image, from a clean `dirty=false` tree)

```text
pages:                  30
fonts embedded:         3/3   (Libertinus Serif regular/bold, DejaVu Sans Mono)
claims indexed:         23/23 (every v2 claim id appears in the rendered text)
evidence digest in pdf: true  (registry sha256 present on the page)
uri annotations:        1
pdf sha256:             da433f01436f33238a381d2db2c61f66d9ad887cd4ad2536e4d83a8c1a7379ef
source commit:          9b044195316b62973520dd8e5d620e8f1141b2e4 (tree dirty=false, built 2026-09-22T19:12:08Z)
typst image digest:     sha256:032e292249bcd378480cc7c142cfa324b63ef8aadeb88d7e7230320c4c9c422f
```

`book/build.sh` exits non-zero if the PDF is implausibly short, any font is not embedded, the
evidence digest is not in the text, or a registry claim id is missing.

## Notes and limitations

- The PDF is a verified skeleton: chapter prose is the next task's deliverable.
- Two build-verification defects were found and fixed during the task (a `pdffonts` column assumption
  and a `grep -q`/`pipefail` presence test that reported false for a true match); both are recorded in
  `LESSONS_LEARNED.md` so the failure modes do not recur.
- No database or benchmark ran; the benchmark lock was untouched.
