# book — the assembled reference

This directory is the reader layer: modular Typst sources, the pinned build image, and the
generated PDF plus its manifest. It contains **no scientific prose of its own** — the chapter and
concept files are structural shells whose prose is written by the HIGH `book-synthesis-v1` task
from `book/evidence/v3/claims.json`.

## Layout

```text
book/
  main.typ                 title page, edition/provenance page, contents, includes, evidence index
  lib/config.typ           the visual system: page geometry, type scale, callouts, evidence cards
  lib/evidence.typ         loads book/evidence/v3/claims.json and renders the confidence cards
  chapters/*.typ           7 chapter shells (how to choose, families, variants, comparisons,
                           topologies, scenarios, evidence and reproduction)
  concepts/*.typ           8 concept shells (placement, indexes, concurrency, derived state,
                           expiry, cache consistency, placement, reading benchmarks)
  evidence/                the v1 and v2 registries (frozen) and the active v3 package
  Containerfile            pinned Typst image (by immutable manifest digest)
  build.sh                 builds and verifies the PDF, writes dist/build-manifest.json
  dist/                    generated PDF and build manifest
```

## Building

```bash
book/build.sh
```

Nothing runs on the host: the script builds the pinned image, compiles `main.typ` to
`dist/data-architecture-reference.pdf`, verifies the result, and writes `dist/build-manifest.json`.
The manifest records the source commit, describe string and dirty state, the build clock, the Typst
version and image digest, the evidence-registry digest, a hash over every `.typ` source, the PDF
sha256, the page count, the embedded-font count, the number of claims indexed and whether the
evidence digest appears in the rendered text.

## Previewing

The sources compile under **any** Typst root: every file path is resolved against the file that
names it, not against the root. That matters for editors whose language server (tinymist, typst-lsp)
roots the project at the repository rather than at `book/` — with a root-relative path the preview
failed immediately with `file not found (searched at <repo>/evidence/v3/claims.json)`. If your
editor still shows that error, set its root to this directory (`book/`), but it should not be
necessary. `book/build.sh` continues to pass `--root book`, so the built artefact is unaffected.

The build **fails** if the PDF is implausibly short, if any font is not embedded, or if the
evidence digest does not appear on the page. A green build means the registry compiled into the
artefact, not merely that Typst exited zero.

## Rules the sources follow

- One visual system: chapters import `lib/config.typ`; they do not set their own fonts or margins.
- One evidence source: `lib/evidence.typ` reads only `evidence/v3/claims.json`; v1 and v2 are
  frozen historical reference.
- No number without a claim: a claim id that does not exist in the registry aborts the build.
- Gaps and analogies are labelled on the page, never smoothed into the prose.
- Fonts are Typst's bundled set, pinned by the image digest; no host font installation is required.
