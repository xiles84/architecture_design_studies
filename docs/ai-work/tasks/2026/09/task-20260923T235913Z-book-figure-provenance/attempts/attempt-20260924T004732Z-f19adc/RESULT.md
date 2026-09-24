# Result — figure provenance for the book build

**Task:** `task-20260923T235913Z-book-figure-provenance`
**Branch:** `repo/book-figure-provenance`
**Required tag:** `repo/book-figure-provenance-v1`
**Capability:** HIGH session executing a LOW task (model `deepseek-flash`)
**Benchmark:** none — infrastructure only

## What was delivered

1. **A settled import mechanism.** `book/FIGURES.md` records the two-variant pilot and chooses the
   generated read-only `book/assets/` export. Pilot sources are kept under `book/assets/pilot/`.
2. **A figure registry and generated manifest.** `book/assets/figures.json` (authored) and
   `book/assets/figure-manifest.json` (generated). Each manifest entry records the canonical source
   path, its git blob revision, source SHA-256, renderer image/digest/id, asset name, embedded-bytes
   SHA-256 and byte count.
3. **An export/gate script.** `book/assets/export.sh --write` renders every registered canonical
   `.puml` with the pinned renderer into `book/assets/`; `--check` fails on a missing asset, a changed
   source or a stale render.
4. **A build gate and wider source hash.** `book/build.sh` runs `export.sh --check` before compiling,
   extends `source_tree_hash_sha256` to cover the figure layer and every referenced source, and adds
   `figure_registry`, `figure_manifest`, `figure_manifest_sha256` and `figures_total` to the manifest.
5. **One registered pilot figure** (`pilot-00-overview`, from
   `studies/01-charity-tree/diagrams/00_overview.puml`) with its exported asset.
6. `CONTEXT.md` records the figure layer.

## Two-variant pilot (same canonical source, pinned Typst image)

```text
variant                                    --root book   --root <repo>
(a) direct repo-root import                    FAIL (1)      ok (0)
(b) generated book/assets export               ok (0)        ok (0)
```

Variant (a) fails with `file not found (searched at book/studies/…/00_overview.svg)`. Chosen: **(b)**.

## Gate controls (all run and recorded)

```text
export.sh --check, consistent tree                         exit 0, "verified (1 figure(s))"
mutation: edit the .puml, keep the old asset               exit 1  (manifest source hash differs)
mutation via book/build.sh on the same edit                exit 1  (fails at the figure gate, dist untouched)
missing asset                                              exit 1  ("asset missing: …")
stale render: source+manifest hash updated, asset left old exit 1  ("stale relative to its source")
export.sh --check after restoring the source               exit 0
```

## Clean-tree build check

`main.typ` compiles in the pinned image with **only `book/` mounted** and `--root book` (exit 0),
which is the mechanism-(b) clean-tree build. A full `book/build.sh` was **not** run, because it would
rewrite the released `book/dist/` PDF; the next real build writes the new manifest keys.

## Conclusion

Figures are now first-class build inputs with source, renderer and embedded-byte provenance, and the
build refuses a missing or stale figure. No measurement was performed and no database or benchmark
lock was touched.
