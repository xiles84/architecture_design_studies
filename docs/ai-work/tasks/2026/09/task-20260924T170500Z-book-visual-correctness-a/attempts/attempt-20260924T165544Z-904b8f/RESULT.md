# Result — Book Edition 2 revision A

**Task:** `task-20260924T170500Z-book-visual-correctness-a`
**Attempt:** `attempt-20260924T165544Z-904b8f`
**Worker:** model `deepseek-flash`, tool Deep Code CLI, effort unknown, capability HIGH, role executor.
**Commit:** `743a4a0` (branch `book/visual-correctness-a`).
**Tag required:** `repo/data-architecture-book-visual-correctness-a`.

## What was delivered

Technical-correctness items A–G of the owner review, the cache-terminology audit, the redraw of the
lease-vs-fence and stale-fill figures, and four new tranche-A structure figures. No measured number,
claim id, trial count, evidence classification, confound or gap changed; `book/evidence/` is untouched.

**Prose (items A–G).**
- `minor-variants.typ`: relaxed freshness is a family of weakenings; only bounded staleness declares a
  maximum Δ and the measured cells declared none (item C).
- `major-families.typ` rollup Mechanism: same database → one transaction; cross-boundary copy → publication
  with delivery, idempotency and reconciliation. `how-to-choose.typ` scopes the application-rollup write
  claim to the measured configuration-portal workload (item D).
- `information-placement.typ`: derive stores no additional copy; the remaining rungs add stored derived
  structure or copy (item E).
- `concurrency-control.typ`: pessimistic contention is waiting/blocking, latency may rise or fall, and
  retries can still be required; 1.76x stays scoped (item F).
- `derived-state-and-history.typ`: derived state separated from retained history; "the only structures"
  removed in favour of the current-state-only fundamental claim. `major-families.typ` drops "writes append
  and fence" (item G).

**Terminology (item 3).** `glossary.typ` now defines invalid, fill lease, fencing token,
source-generation validation, cache CAS, hit validation, transactional outbox and TTL/early expiry
separately; the cache chapter's "Publication fence / CAS" card and matrix row were renamed to name
**source-generation validation**, and the chapter intro states the distinction from a sink-side fencing
token and from cache CAS.

**Figures (items A, B, 2.1, 2.3+2.7, 2.8).**
- `fig_lease_vs_fence` redrawn with five lifelines and two panels plus a distinction note; the coordinator
  no longer appears to publish to the cache and the "the source generation rejects it" note is gone.
- `fig_cache_stale_fill` redrawn with explicit `S0/v0 → S1/v1` provenance, the mechanism named
  source-generation validation, an authoritative-read path for R2, and scoping to the stale-fill race.
- New `fig_families_er`, `fig_rollup_vs_rolldown`, `fig_derived_vs_history`, `fig_history_vs_current`,
  `fig_topology_physical_boundary`, all labelled form × evidence status with a text equivalent.
- Shared `_style.puml` gained authoritative / derived / external stereotypes for one visual language.

## Verification

- `book/assets/export.sh --write` then `--check`: 13 registered figures, pinned renderer digest
  `sha256:9b9ee6af54a8…`, `strip-textlength-v1` applied.
- `book/evidence/check.sh`: OK (no orphan figure, no embedded number, no stretch, registry path derived).
- `book/build.sh --out dist/preview.pdf --manifest dist/preview-manifest.json`: **57 pages**, 29/29 claims
  indexed, 6/6 fonts embedded, evidence digest present, no unverified banner, no unresolved tokens.
  Both preview files were deleted before committing.
- Rasterised pages 10, 13, 31, 32, 38, 41, 44 and inspected them: no clipping, overlap or broken arrows.
  Figure 10 (lease-vs-fence) is rendered at 74% width via the new `image-width` parameter because at full
  width the tall sequence overflowed into the running footer; with that width it fits with its caption,
  text equivalent and the following card.

## Decision Log

| # | Decision | Reason | Confidence | Affects |
|---|---|---|---|---|
| D1 | Split the brief's single `fig_history_vs_derived` into `fig_derived_vs_history` + `fig_history_vs_current` | Three disconnected panels were reordered by the renderer and read 3-1-2; splitting fixes both reading order and legibility, and matches the review's separate items 2.3 and 2.7 | High | Tranche B page budget |
| D2 | Lease-vs-fence is two panels plus a distinction note, not three panels | A third panel made the sequence taller than the page; it overflowed the footer. The review allowed the distinction as a labelled note | High | — |
| D3 | Added `image-width` to `figure-evidence` and render figure 10 at 74% | The tall sequence still exceeded the page at full width; narrowing is the only fix that keeps all accepted content on one page | High | Any future tall figure |
| D4 | Renamed the cache mechanism card to "Publication fence — source-generation validation" | The card's mechanism is a source check; a sink-side fencing token and cache CAS are different and must not share its name | High | Cache chapter, glossary |
| D5 | Figure numbering was not asserted in prose | Typst owns the numbers and they shifted when figures were inserted; captions and `FIGURES.md` name figures by id | High | — |

## Where this is weak

- **Visual review is one model's.** The rasterised-page inspection is not an independent read, and no
  reader test measured comprehension gain.
- **Figure 10 is small.** At 74% of text width its labels are around 5.5–6 pt; legible in the raster but
  a print reader may disagree. Splitting it into two full-width figures is the fallback if HIGH says so.
- **The combined figure was split, so the book is 57 pages, not 59.** Tranche B's four figures will add
  more; the page-count checkpoint belongs to that task.
- **`book/dist/` was not touched**, so the tracked PDF is still the 54-page Edition 2 build until the
  release task runs.
