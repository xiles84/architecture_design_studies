# Figures — import mechanism and provenance

The book's figures are **generated artefacts**. One canonical editable source per figure lives either in
the study that owns it (study figures) or under `book/assets/sources/` (figures the book itself owns);
the book holds only generated, read-only SVG assets beside its Typst sources. Nothing in `book/` is an
independent copy of a **study** figure.

- **Registry (authored):** [`assets/figures.json`](assets/figures.json) — one entry per figure:
  `id`, canonical `source` (a study `.puml` or a `book/assets/sources/*.puml`), generated `asset`, and a
  `status`.
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

## The proof set

Every book figure carries two independent labels — its **form** (`structure` / `sequence` / `state`)
and its **evidence status** (`conceptual illustration` / `implemented design contract` / `negative
control` / `observed result`) — plus an adjacent text equivalent in the prose. The rule is
[methodology 15](../docs/methodology.md).

| Figure | Source | Form | Evidence status | Placed in |
|---|---|---|---|---|
| 1 — placement and ownership | `assets/sources/fig_placement_ownership.puml` (book) | structure | conceptual illustration | `concepts/information-placement.typ` |
| 2 — two buyers, same marked seats | `studies/03-reserved-seating/diagrams/s_arbitration.puml` (study reuse) | sequence | implemented design contract (contains a labelled negative control) | `chapters/scenarios.typ` |
| 3 — the stale-fill race and the fence | `assets/sources/fig_cache_stale_fill.puml` (book) | sequence | observed result (illustrates `v2-15`; embeds no rate) | `concepts/cache-consistency.typ` |
| 4 — when an abandoned hold returns its seat | `studies/03-reserved-seating/diagrams/e_expiry.puml` (study reuse) | sequence | implemented design contract (contains a labelled negative control) | `concepts/expiry-and-clock-authority.typ` |

`pilot-00-overview` is the import-mechanism pilot and remains registered but is not a book figure.

**Two-reader check (2026-09-23).** Two reading passes were run: a first reader (what does the picture
say without the caption?) and a second reader (does it match the claim/SQL, and does it imply a
performance number?). Both passes were by the same model, which is weaker than two independent readers.
Outcomes: figure 1 kept; figure 2 kept; **figure 3 redesigned** (the first version showed one reader,
but claim `v2-15` is about instances sharing a cache, so a second reader in another instance was added
to show the stale hit); figure 4 **relabelled**, not redesigned (its form was corrected from `state`
to `sequence` — the study source is a sequence and the book does not hold an editable copy). No figure
was dropped.



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

1. Author or edit the canonical `.puml`: a study figure in its study (render with
   `studies/<study>/diagrams/render.sh`), a book figure under `assets/sources/` (with `_style.puml`
   beside it).
2. Add or keep its entry in `assets/figures.json`.
3. Run `book/assets/export.sh --write`; commit the source, the asset and the refreshed manifest.
4. Reference the asset from a book source with a file-relative path (`../pilot-00-overview.svg` from a
   file in `assets/`).

A canonical source change **always** requires step 3. The renderer is pinned once in
[`../infra/versions.env`](../infra/versions.env) (`PLANTUML_IMAGE`, digest
`sha256:9b9ee6af…`, PlantUML 1.2026.8); see [`../infra/diagram-render.sh`](../infra/diagram-render.sh).
