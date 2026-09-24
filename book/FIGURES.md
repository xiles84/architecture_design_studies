# Figures — import mechanism and provenance

The book's figures are **generated artefacts**. One canonical editable source per figure lives in
the study that owns it; the book holds only a generated, read-only SVG asset beside its Typst
sources. Nothing in `book/` is an independent copy of a figure.

- **Registry (authored):** [`assets/figures.json`](assets/figures.json) — one entry per figure:
  `id`, canonical `source` (a study `.puml`), generated `asset`, and a `status`.
- **Manifest (generated):** [`assets/figure-manifest.json`](assets/figure-manifest.json) — per figure:
  the source path, its recorded git blob revision, its SHA-256, the renderer image/digest/id, the
  asset name, the asset SHA-256 (the embedded bytes) and the byte count.
- **Exporter/gate:** [`assets/export.sh`](assets/export.sh) — `--write` renders every source with the
  pinned renderer and refreshes the assets and manifest; `--check` fails on a missing asset, a source
  that changed without a re-export, or an asset that no longer matches a fresh pinned render.

## The import mechanism was chosen by pilot, not preference

Two mechanisms were executed against the same canonical source
(`studies/01-charity-tree/diagrams/00_overview.puml`, a structure figure) in the pinned Typst image:

| Mechanism | Editor preview | Clean-tree build | Manifest coverage |
|---|---|---|---|
| **(a) direct repo-root import** — `#image("/studies/…/rendered/00_overview.svg")` | only when the editor roots at the **repository** (a leading-slash path is root-relative) | compiles with `--root <repo>` (exit 0) but **fails** with `--root book`: `file not found (searched at book/studies/…/00_overview.svg)` (exit 1) | requires hashing files outside `book/`; the build would have to mount the whole repository |
| **(b) generated `book/assets/` export** — `#image("../pilot-00-overview.svg")` | works under **either** root, because a path without a leading slash resolves against the file | compiles with `--root book` (exit 0) **and** `--root <repo>` (exit 0) | the asset, registry and manifest are all inside `book/`, so the source-tree hash covers them naturally |

**Chosen: (b) the generated `book/assets/` export.** It is the only mechanism that keeps the build's
`--root book` contract (and therefore the existing container mount and the editor preview) while
letting the manifest cover the embedded bytes. Pilot sources are kept under `assets/pilot/` so the
comparison stays reproducible.

## The build gate

`book/build.sh` runs `assets/export.sh --check` **before** compiling. The build therefore fails when:

- a figure's asset file is missing;
- a canonical `.puml` changed without `assets/export.sh --write` (the manifest's recorded source hash
  no longer matches);
- the asset no longer matches a fresh render with the pinned renderer (so an asset cannot be
  hand-edited or left behind after a source change).

`source_tree_hash_sha256` in `dist/build-manifest.json` now covers the Typst sources, everything under
`book/assets/` and every referenced canonical figure source, so a changed figure cannot leave the hash
unchanged. `dist/build-manifest.json` also records `figure_registry`, `figure_manifest`,
`figure_manifest_sha256` and `figures_total`.

## Adding or changing a figure

1. Author or edit the canonical `.puml` in its study; render it with the pinned renderer
   (`studies/<study>/diagrams/render.sh`).
2. Add or keep its entry in `assets/figures.json`.
3. Run `book/assets/export.sh --write`; commit the source, the asset and the refreshed manifest.
4. Reference the asset from a book source with a file-relative path (`../pilot-00-overview.svg` from a
   file in `assets/`).

A canonical source change **always** requires step 3. The renderer is pinned once in
[`../infra/versions.env`](../infra/versions.env) (`PLANTUML_IMAGE`, digest
`sha256:9b9ee6af…`, PlantUML 1.2026.8); see [`../infra/diagram-render.sh`](../infra/diagram-render.sh).
