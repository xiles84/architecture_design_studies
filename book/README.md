# book — the assembled reference

This directory is the reader layer: modular Typst sources, the pinned build image, and the
generated PDF plus its manifest. The chapter and concept files carry the prose; every number in them
is rendered from the active evidence package under `book/evidence/`, which is validated by
[`tools/evidence`](../tools/evidence/README.md).

## Layout

```text
book/
  STYLE.md                 the didactic contract: structure, chapter template, writing rules
  main.typ                 title page, edition/provenance page, contents, the Edition 3 parts,
                           the Edition 2 divider, glossary, evidence index
  front/*.typ              Edition 3 front matter (how to use this book, with the badge legend)
  lib/config.typ           the visual system, the callouts, and the Edition 3 components
                           (concept-opener, fit-table, cost-table, case-box, remember,
                           checklist-table, part-opener) — each validating its input
  lib/evidence.typ         loads the active evidence package, renders the confidence cards, and
                           renders inline `claim-badge` citations that link to the card's page
  cases/*.typ              Edition 3 study profiles (one page per study with claims)
  chapters/*.typ           7 Edition 2 chapters (how to choose, families, variants, comparisons,
                           topologies, scenarios, evidence and reproduction) — awaiting conversion
  concepts/*.typ           Edition 3 concept chapters (rollup) plus the 8 Edition 2 concepts
                           (placement, indexes, concurrency, derived state, expiry, cache
                           consistency, distributed placement, reading benchmarks)
  evidence/                the frozen v1-v5 packages, the active v6 package, and check.sh
  assets/                  generated figures, their canonical .puml sources and the manifest
  Containerfile            pinned Typst image (by immutable manifest digest)
  build.sh                 builds and verifies the PDF, writes dist/build-manifest.json
  dist/                    generated PDF and build manifest
```

## Building

```bash
book/build.sh
book/build.sh --out dist/preview.pdf --manifest dist/preview-manifest.json   # scratch build
book/evidence/check.sh                                                       # source rules only
```

## Editions

`book-edition` in `lib/config.typ` names the edition the sources build. **Edition 3 preview 1** is the
concept-first edition (owner decision 2026-10-01, `STYLE.md`): each chapter teaches one concept and the
studies show when it fits. It is being built part by part — a part opener names the chapters it will
hold, a converted chapter appears in its part, and the Edition 2 chapters follow a divider that says
they await conversion. The first converted part is part I's chapter 4 (Rollup) plus the Study 04
profile; `front/how-to-use.typ` explains the shape to the reader.

The owner authorized integration and a new local release on 2026-10-04. The tracked
`dist/data-architecture-reference.pdf` now contains the approved pilot slice and the retained
Edition 2 chapters. This is a preview release: the remaining concept-first conversion is unfinished.
Its immutable milestone is `repo/data-architecture-book-v3-preview-1`; the earlier Edition 2
release remains available through its existing tags. The original pilot PDF and manifest are
preserved under `docs/book-release-20261004/archive/edition3-pilot/`.

Nothing runs on the host: the script builds the pinned image, compiles `main.typ` to
`dist/data-architecture-reference.pdf`, verifies the result, and writes `dist/build-manifest.json`.
The manifest records the source commit, describe string and dirty state, the build clock, the Typst
version and image digest, the active registry version and its digest, a hash over every `.typ`
source, the figure manifest, the PDF sha256, the page count, the embedded-font count, how many
claims render, and whether the unverified-build banner appeared.

Rebuilding is serialised by a `ads-book-build-lock` podman volume, because every worktree shares the
`localhost/ads-book:1` image tag. That lock is separate from the benchmark lock and never touches it.

## Previewing

The sources compile under **any** Typst root: every file path is resolved against the file that
names it, not against the root. That matters for editors whose language server (tinymist, typst-lsp)
roots the project at the repository rather than at `book/`. A preview compiled without the build
script's `--input` values renders an explicit *Unverified development build* banner instead of a
provenance page full of placeholders.

## Rules the sources follow

- One visual system: chapters import `lib/config.typ`; they do not set their own fonts or margins.
- One evidence source: `lib/evidence.typ` derives its path from the single `registry_version` build
  input. No other file may name a version — `book/evidence/check.sh` fails the build if one does,
  because Edition 1 printed a v2 path in a chapter while the cover said v3.
- No number without a claim: a claim id that is unknown, retired or superseded aborts the build.
- Edition 3 components validate their input: a key idea over 25 words, a fit table without
  3-4 rows per side, a verdict outside {fits, does not fit, depends}, a checklist with other
  than 8 rows, an unknown evidence status, more than 3 remember bullets or a row longer than
  its word limit panics with the rule and the offending value.
- A gap states what sort of gap it is (`gap_kind`), so a schema limitation is not badged as a regime
  that was simply not run.
- Gaps and analogies are labelled on the page, never smoothed into the prose.
- Typst owns figure numbering; no asset may draw its own "Figure N".
- Fonts are Typst's bundled set, pinned by the image digest; no host font installation is required.
