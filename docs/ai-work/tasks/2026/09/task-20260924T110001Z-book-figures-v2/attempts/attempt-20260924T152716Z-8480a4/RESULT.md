# Result — figure numbering, SVG stretching and marker tracking

**Task:** `task-20260924T110001Z-book-figures-v2`
**Branch:** `book/figures-v2`
**Required tag:** `repo/book-figures-v2`
**Capability:** HIGH session executing a LOW task (model `deepseek-flash`, Deep Code CLI; effort not exposed)
**Benchmark:** none — the only container started is the pinned PlantUML renderer

## What was delivered

1. **Typst owns figure numbering.** The two book-owned PlantUML sources drew their own titles —
   `Figure 3 — the stale-fill race and the fence` and `Figure 1 — placement and ownership: where an
   answer may live` — while Typst also numbers the figure, so each page showed two numbers for one
   image. The titles now carry only the description, and the rendered page shows **one** number per
   figure: `Figure 1: Two buyers…`, `Figure 2: Placement and ownership…`, `Figure 3: When an abandoned
   hold…`, `Figure 4: The stale-fill race…`, captions only, with no in-diagram number.
2. **The stretch is fixed at the export step, not per figure.** PlantUML emits `textLength` and
   `lengthAdjust="spacing"` on every `<text>` element and Typst honours them, so each label was laid out
   at the renderer's guessed width rather than in the real font. **All five assets carried them** — the
   long-labelled stale-fill diagram was simply where it became obvious, and a fix aimed at that figure
   would have left the other four subtly wrong. `book/assets/export.sh` now strips both attributes,
   asserts the strip succeeded, and records `postprocess: strip-textlength-v1` per figure in
   `figure-manifest.json` (now schema_version 2), so changing the step invalidates every asset exactly
   as a renderer change does.
3. **Both figure rules are fatal with no waivers.** A new `figure-stretch` rule fails any asset carrying
   either attribute; the two `figure-number` waiver lines are deleted so the numbering rule applies
   everywhere. `check.sh` now reports OK with an empty waiver file.
4. **Marker tracking reduced** from 0.8pt to 0.3pt.
5. **A stale default fixed:** `check.sh` kept its own `REGISTRY_VERSION="v4"` after the book moved to v5.
   Every build was correct because `build.sh` passes the version explicitly, while a standalone
   `check.sh` run validated the frozen v4 package and reported OK — two agents could have run "the same
   check" against one tree and validated different registries.

## Verified from committed files

```text
book/assets/export.sh --check        passes; postprocess reported per asset (5 figures)
check.sh --book <absolute>           OK, with an empty waiver file
stretch attributes in any asset      0
embedded "Figure [0-9]" in any asset 0
rendered figure captions             Figures 1-4, one number each, no in-diagram number
book/build.sh --out dist/preview.pdf 51 pages, 29/29 claims indexed, 6/6 fonts embedded,
                                     evidence digest present, banner absent, no unresolved tokens
```

### Negative controls

| Control | Expected | Observed |
|---|---|---|
| A `textLength="300"` added by hand to an asset | `figure-stretch` fires | 1 failure, asset named |
| A `<text>Figure 9 — test</text>` appended to an asset | `figure-number` fires | 1 failure, asset named |
| Both reverted | checks pass again | `check: OK` |

## Where this task is weak

- **The stretch fix is verified mechanically, not visually.** I confirmed the attributes are gone, the
  assets changed, `--check` passes and each figure has one rendered number. I did not render pages to
  images and compare label widths, so "readable at 100% without zoom" rests on the attribute being the
  cause — which the reviewer's own diagnosis and the export evidence both support.
- **`figure-manifest.json` moved to schema_version 2.** Nothing reads version 1 fields defensively; a
  consumer outside this repository would need to know.
- **The remaining assets are not reviewed as designs.** This task fixes their typography, not their
  content; V2/V3/V9/V10/V12 and the lease-versus-fence figure are the pedagogy task's business.
- **The marker tracking change is a taste call** at the edge of measurement: 0.3pt still separates the
  letters of small-caps markers without the mechanical look of 0.8pt.

## Conclusion

One number per figure, owned by Typst; no stretched labels, fixed where the artefact is generated and
gated by a rule that fails rather than warns; no standing waivers left in the book; and one latent
duplicated default removed before it could make two checks disagree.
