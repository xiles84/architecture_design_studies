# Review — task-20260924T170500Z-book-visual-correctness-a

**Reviewer:** HIGH session, model `deepseek-flash`, Deep Code CLI (effort not exposed)
**Verdict:** approve
**Disclosure:** self-review. The executor and reviewer are the same session, which is weaker than an
independent review; the owner should weigh that when judging the figures.

## What was checked

| Check | Method | Result |
|---|---|---|
| Item C wording | read `minor-variants.typ` | relaxed is a family; only bounded staleness declares Δ; measured cells unbounded-relaxed |
| Item D wording | read `major-families.typ`, `how-to-choose.typ` | same-transaction vs cross-boundary publication; write claim scoped to the measured workload |
| Item E wording | read `information-placement.typ` | "derive stores no additional copy"; no contradiction |
| Item F wording | read `concurrency-control.typ` | waiting/blocking framing; retries still possible; 1.76x scoped |
| Item G wording | read `derived-state-and-history.typ`, `major-families.typ` | derived state vs retained history; "append and fence" gone |
| Terminology | read `glossary.typ`, cache chapter | eight terms distinct; card renamed to source-generation validation |
| Figure 10 | source + rasterised page 38 | five lifelines, two panels + distinction note, no coordinator→cache, no overflow |
| Figure 11 | source + rasterised page 41 | explicit S0/v0 → S1/v1, source-generation validation, stale-fill scope stated |
| New figures | `figures.json`, manifest, pages 10/13/31/32/44 | registered, embedded once, labelled, text equivalent, no number/rate |
| No number changed | `git diff` over `book/evidence`, prose numbers | registry untouched; no measured value edited |
| Build | `export.sh --check`, `check.sh`, `build.sh` preview | 13 figures; check OK; 57 pages, 29/29 claims, 6/6 fonts, no banner, no unresolved tokens |

## Judgement on the Decision Log

1. **D1 (split the history figure) is accepted and better than the brief.** Three disconnected panels
   were reordered by the renderer and read 3-1-2; splitting fixed both order and legibility, and maps
   cleanly onto the review's separate 2.3 and 2.7. The extra figure is a page budget cost for tranche B,
   recorded.
2. **D2/D3 (two panels + `image-width`) are accepted.** A three-panel sequence overflowed the footer;
   the review explicitly allowed the distinction as a labelled note, and the width parameter is a
   reusable fix. Narrowing to 74% is the only accepted-content-preserving option.
3. **D4 (card rename) is accepted.** It is exactly what the review's terminology audit requires.
4. **D5 (no figure numbers in prose) is accepted.**

## Where this task is weak

- **One-model self-review.** No independent reader judged the figures or the corrected prose.
- **Figure 10 is small** (~5.5–6 pt labels at 74% width). Legible in the raster; a print reader may
  disagree. The fallback is to split it into two full-width figures.
- **The page count moved to 57, not the 59 first seen.** Tranche B's four figures still need the
  page-count checkpoint the brief requires.
- **`book/dist/` is untouched**, so the release artefact is still the 54-page Edition 2 build.
- The combined figure split slightly exceeds the brief's named figure list; it is recorded here and in
  `FIGURES.md`, and it changes no acceptance criterion.

## Required tag

`repo/data-architecture-book-visual-correctness-a`, created by `queue integrate` on the integrated state.
