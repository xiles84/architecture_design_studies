# Result — Book Edition 2 revision B

**Task:** `task-20260924T170501Z-book-visual-tranche-b`
**Attempt:** `attempt-20260924T171649Z-1666c8`
**Worker:** model `deepseek-flash`, tool Deep Code CLI, effort unknown, capability HIGH, role executor.
**Commit:** `16be6f6` (branch `book/visual-tranche-b`).
**Tag required:** `repo/data-architecture-book-visual-tranche-b`.

## What was delivered

The four remaining accepted visuals from the owner review, each a new book-owned PlantUML source
registered in `figures.json`, embedded once with a form × evidence-status label and a text equivalent:

| id | content | placed in |
|---|---|---|
| `fig-embedded-vs-normalized` | normalized child rows vs a bounded collection inside the parent, with the whole-aggregate vs one-child trade-off | `chapters/major-families.typ`, Embedded documents |
| `fig-cache-read-copy` | client → application → {authoritative source database, non-authoritative cache}, annotated with source vs copy correctness | `concepts/cache-consistency.typ`, before the races |
| `fig-index-vs-hotrow` | the narrow index path vs several writers serialising on one contended row | `concepts/indexes-and-hot-rows.typ` |
| `fig-optimistic-vs-pessimistic` | parallel optimistic (version-check) and pessimistic (lock) sequences, with the invariant recheck named as what rejects an invalidated operation | `concepts/concurrency-control.typ` |

No measured number, claim id, trial count, evidence classification or gap changed; `book/evidence/` is
untouched.

## Verification

- `book/assets/export.sh --write` then `--check`: **17 registered figures**, pinned renderer
  `sha256:9b9ee6af54a8…`, `strip-textlength-v1` applied.
- `book/evidence/check.sh`: OK.
- `book/build.sh --out dist/preview.pdf --manifest dist/preview-manifest.json`: **61 pages**, 29/29 claims
  indexed, 6/6 fonts embedded, evidence digest present, no banner, no unresolved tokens. Both preview
  files deleted before committing.
- Rasterised and inspected the four figure pages (14, 30, 32, 37): no clipping, overlap or broken arrows.
  `fig-optimistic-vs-pessimistic` is rendered at 72% width via `image-width` because the tall sequence
  would otherwise overflow the footer.

## Page-count checkpoint

Revision A left the book at 57 pages. This tranche took it to **61 pages** (+4). The two wide structure
figures are read at the top of their concepts and the two-page embeddings family gained one figure; no
figure needed a split, so no accepted item was dropped. The growth is within the plan's expectation and
is recorded for the release task.

## Decision Log

| # | Decision | Reason | Confidence | Affects |
|---|---|---|---|---|
| D6 | `fig-index-vs-hotrow` panels stacked vertically, not side by side | The side-by-side version placed the hot-row panel below-right of the index panel, a diagonal that read poorly; a forced vertical stack is narrow and legible | High | — |
| D7 | `fig-optimistic-vs-pessimistic` rendered at 72% width | Two full sequence panels at full width exceed the page height | High | Any future tall sequence |
| D8 | The optimistic/pessimistic figure is generic (row at version 5), not tied to a study design | The concept chapter is reusable; a study-specific story would duplicate the scenario figure and over-scope the section | High | — |

## Where this is weak

- **One-model self-review.** As in revision A, no independent reader judged the four figures.
- **Two figures are rendered below full width** (`fig-optimistic-vs-pessimistic` at 72%; the revision-A
  lease-vs-fence at 74%). Legible in the raster, but a print reader may disagree.
- **Page growth is real.** 54 → 61 pages across the revision; the release task should confirm the PDF is
  still coherent.
- **`book/dist/` is untouched**; the tracked PDF remains the 54-page Edition 2 build until release.
