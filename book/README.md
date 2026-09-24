# book — the assembled reference

This directory is the reader layer: modular Typst sources, the pinned build image, and the
generated PDF plus its manifest. The chapter and concept files carry the prose; every number in them
is rendered from the active evidence package under `book/evidence/`, which is validated by
[`tools/evidence`](../tools/evidence/README.md).

## Layout

```text
book/
  main.typ                 title page, edition/provenance page, contents, includes, evidence index
  lib/config.typ           the visual system: page geometry, type scale, callouts
  lib/evidence.typ         loads the active evidence package and renders the confidence cards
  chapters/*.typ           7 chapters (how to choose, families, variants, comparisons,
                           topologies, scenarios, evidence and reproduction)
  concepts/*.typ           8 concepts (placement, indexes, concurrency, derived state,
                           expiry, cache consistency, distributed placement, reading benchmarks)
  evidence/                the frozen v1-v4 packages, the active v5 package, and check.sh
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
- A gap states what sort of gap it is (`gap_kind`), so a schema limitation is not badged as a regime
  that was simply not run.
- Gaps and analogies are labelled on the page, never smoothed into the prose.
- Typst owns figure numbering; no asset may draw its own "Figure N".
- Fonts are Typst's bundled set, pinned by the image digest; no host font installation is required.
