# Progress — author the book figure proof set

**Task:** `task-20260923T235922Z-book-figure-proofset`
**Branch:** `book/figure-proofset`
**Attempt:** `attempt-20260924T012027Z-11770e` (claim `claim-a77de8b930a34624`, epoch 1)
**Capability/role:** HIGH session executing the task / analyst (model `deepseek-flash`)
**Required tag:** `repo/book-figures-v1`

Authorship and build. No database started, no benchmark lock taken, no measurement.

## Decision Log

1. **Book-owned figures get a canonical source under `book/assets/sources/`.** The non-duplication rule
   protects *study* figures; a figure the book itself owns needs a canonical editable source inside the
   book's owned paths. Study figures keep their study source and only the generated asset is embedded.
   Confidence: High.
2. **Book-only figure IDs and file names.** `fig-placement-ownership` and `fig-cache-stale-fill` are
   authored for the book; `fig-hold-arbitration` and `fig-expiry-clock` reuse study-03 sources
   unchanged (adaptation = selection, caption and text equivalent, not redraw). Confidence: High.
3. **Labels are a two-label mechanism, not decoration.** `figure-evidence(asset, form, status, caption,
   equivalent, claim)` renders the form and evidence-status labels and the adjacent text equivalent.
   An `observed result` figure names its claim and embeds no number. Confidence: High.
4. **Figure 3 redesigned after the two-reader check.** The first version showed one reader, but the
   claim it illustrates (`v2-15`) is about instances sharing a cache; a second reader in another
   instance was added so the stale hit is visible. Confidence: High.
5. **Figure 4 relabelled, not redesigned.** Its form was corrected from `state` to `sequence`; the
   study source is a sequence and the book must not hold an editable copy. Confidence: High.
6. **Both reader passes were the same model.** Recorded as a limitation in `book/FIGURES.md` and
   methodology 15 rather than presented as independent review. Confidence: High.
7. **The PDF was rebuilt.** This task owns `book/main.typ` and `book/dist`; the build wrote a new
   49-page draft with 29/29 v3 claims and the figure keys. It is labelled a draft, not a release; the
   released v1 PDF stays at `repo/data-architecture-book-v1`. Confidence: High.

## What was produced

- `book/assets/sources/{_style,fig_placement_ownership,fig_cache_stale_fill}.puml`
- `book/assets/{fig-placement-ownership,fig-hold-arbitration,fig-cache-stale-fill,fig-expiry-clock}.svg`
- `book/assets/{figures.json,figure-manifest.json}` (5 registered figures)
- `book/lib/config.typ` — `figure-evidence` helper
- Figures placed in `concepts/information-placement.typ`, `chapters/scenarios.typ`,
  `concepts/cache-consistency.typ`, `concepts/expiry-and-clock-authority.typ`
- `docs/methodology.md` — rule 15 (figure labels, text equivalent, two-reader check)
- `book/FIGURES.md` — the proof set, labels and the reader-check record
- `book/dist/` — rebuilt PDF + manifest; `CONTEXT.md` updated

## Residual risks / limitations

- The two-reader check was two passes by one model, not two independent readers.
- Figures 2 and 4 reuse study sources unchanged; a future study-source change requires a re-export and
  the build will fail until it happens.
- The rebuilt PDF is a draft build; it has not gone through a release review, and its prose predates no
  claim (29/29 v3 claims indexed).
