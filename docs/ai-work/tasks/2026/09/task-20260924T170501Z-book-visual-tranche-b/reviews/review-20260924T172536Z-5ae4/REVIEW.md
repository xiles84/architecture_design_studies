# Review — task-20260924T170501Z-book-visual-tranche-b

**Reviewer:** HIGH session, model `deepseek-flash`, Deep Code CLI (effort not exposed)
**Verdict:** approve
**Disclosure:** self-review by the same session that executed both tranches; weaker than an independent
review.

## What was checked

| Check | Method | Result |
|---|---|---|
| Four figures registered | `figures.json` + manifest | 17 figures; each new id present with a source and asset |
| Embedded once | `check.sh` orphan-figure rule | OK |
| Form × evidence label | read the four `figure-evidence` calls | structure ×3, sequence ×1; all conceptual illustration |
| Text equivalent | read the four calls | present and specific |
| No embedded number or rate | grep the sources | none |
| Placement | read the four chapter call sites | embedded family, cache intro, index concept, concurrency concept |
| Fits the page | rasterised pages 14, 30, 32, 37 | no clipping or overlap; the sequence uses `image-width: 72%` |
| Build | `export.sh --check`, `check.sh`, `build.sh` preview | 61 pages, 29/29 claims, 6/6 fonts, no banner, no unresolved tokens |
| No number changed | git diff over `book/evidence` and prose numbers | registry untouched |

## Judgement on the Decision Log

1. **D6 (vertical index/hot-row stack) is accepted.** The first render placed the hot-row panel below-right
   of the index panel; the vertical stack reads as two ordered panels and is narrower.
2. **D7 (72% sequence width) is accepted** for the same reason as revision A's 74%.
3. **D8 (generic concurrency example) is accepted.** The concept chapter is reusable and a study story
   would duplicate the scenario figure.

## Where this task is weak

- **One-model self-review**, as in revision A.
- **Two figures render below full width**; legible in the raster but not independently confirmed in print.
- **The book grew 54 → 61 pages** across the revision; the release task must confirm the artefact is
  coherent and record the final page count.
- **`book/dist/` is untouched**, so the tracked PDF is still the 54-page Edition 2 build.

## Required tag

`repo/data-architecture-book-visual-tranche-b`, created by `queue integrate`.
