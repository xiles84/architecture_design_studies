# Progress — add figure provenance to the book build; settle the import mechanism

**Task:** `task-20260923T235913Z-book-figure-provenance`
**Branch:** `repo/book-figure-provenance`
**Attempt:** `attempt-20260924T004732Z-f19adc` (claim `claim-04e550978c310c3d`, epoch 1)
**Capability/role:** HIGH session executing a LOW task / executor (model `deepseek-flash`)
**Required tag:** `repo/book-figure-provenance-v1`

Infrastructure-only. No database started, no benchmark lock taken, no measurement.

## Decision Log

1. **Import mechanism = (b) the generated `book/assets/` export, chosen by pilot.** Both mechanisms
   were compiled from the same canonical source in the pinned Typst image: (a) direct repo-root import
   compiles only with `--root <repo>` and fails with `--root book`; (b) the file-relative asset import
   compiles under both roots. (b) keeps the build's `--root book` contract and the editor preview, and
   its bytes live inside `book/` so the source-tree hash covers them. Confidence: High.
2. **Registry authored, manifest generated.** `book/assets/figures.json` is the hand-authored list
   (id, canonical `.puml` source, generated asset); `book/assets/figure-manifest.json` is written by
   `book/assets/export.sh --write` and carries the source revision, source/asset SHA-256, renderer
   digest and byte count. Both live in the owned `book/assets/`. Confidence: High.
3. **One canonical editable source per figure; no book copy.** The canonical source stays a study
   `.puml`; the book holds only the generated SVG. This is the brainstorm's non-duplication rule.
   Confidence: High.
4. **The gate is a fresh pinned re-render, not just a recorded hash.** `export.sh --check` re-renders
   every source with `infra/diagram-render.sh` and compares bytes, and also compares the manifest's
   recorded source hash and the asset file. Three negative controls were run (missing asset, changed
   source, stale render) and each fails. Confidence: High.
5. **`build.sh` gates before compiling and extends the source-tree hash.** It runs
   `export.sh --check`, then hashes the Typst sources, the whole `book/assets/` layer and every
   referenced canonical source, and records `figure_registry`/`figure_manifest`/
   `figure_manifest_sha256`/`figures_total` in `dist/build-manifest.json`. Confidence: High.
6. **Pilot figure is `00_overview.puml`, not `d3_flattened_fk.puml`.** The fidelity-audit task is
   required to correct the D3 wording; using D3 as the pilot would leave the manifest stale the moment
   that task re-renders. The overview is a structural figure unlikely to need a factual correction.
   Confidence: Medium-High.
7. **No full `book/build.sh` run.** A full build rewrites the released `book/dist/` PDF, which the
   acceptance does not require; the clean-tree-build criterion was instead proven by compiling
   `main.typ` in the pinned image with only `book/` mounted (exit 0) and by the build's mutation
   failure test (exit 1). The manifest keys are exercised by the next real build. Confidence: Medium.

## What was produced

- `book/FIGURES.md` — mechanism decision, pilot table, build-gate contract, how to add a figure
- `book/assets/figures.json`, `book/assets/figure-manifest.json`, `book/assets/pilot-00-overview.svg`
- `book/assets/export.sh` — export (`--write`) and provenance gate (`--check`)
- `book/assets/pilot/direct.typ`, `book/assets/pilot/export.typ` — the two pilot variants
- `book/build.sh` — figure gate, extended source-tree hash, new manifest keys
- `CONTEXT.md` — the book figure-layer bullet

## Residual risks / limitations

- Any change to a registered canonical `.puml` requires `book/assets/export.sh --write`; the build
  will (correctly) fail until it is re-exported. The fidelity audit must re-export if it touches
  `00_overview.puml`.
- The check re-renders a whole diagrams directory per distinct source directory; with a large proof
  set this should be cached or limited, but it is exact.
- `book/dist/` was not rebuilt, so the committed `build-manifest.json` does not yet carry the new
  figure keys; the next real build will write them.
