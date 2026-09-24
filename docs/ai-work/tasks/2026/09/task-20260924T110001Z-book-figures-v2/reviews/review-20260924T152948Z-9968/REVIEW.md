# Review — task-20260924T110001Z-book-figures-v2

**Reviewer:** HIGH session, model `deepseek-flash`, Deep Code CLI (effort not exposed)
**Verdict:** approve
**Disclosure:** self-review, against the committed assets, the re-run exporter and the rebuilt PDF.

## What was checked

| Check | Method | Result |
|---|---|---|
| Typst alone owns the number | rendered text | Figures 1–4, caption only; `Figure 3 —` and `Figure 1 —` appear zero times |
| The sources carry no number | grep the two `.puml` titles | description only |
| No asset is stretched | grep all five SVGs for `textLength`/`lengthAdjust` | zero |
| The fix survives a re-export | `export.sh --check` | passes, reporting `strip-textlength-v1` per asset |
| The step is recorded | `figure-manifest.json` | `postprocess` per figure, `schema_version: 2` |
| Both rules are fatal | `check.sh` with an empty waiver file | OK; controls below fire when violated |
| `figure-stretch` control | a hand-added `textLength="300"` | 1 failure naming the asset |
| `figure-number` control | a hand-appended `Figure 9` | 1 failure naming the asset |
| No regressions | `book/build.sh` | 51 pages, 29/29 claims, 6/6 fonts, digest present, no banner, no unresolved tokens |

## Judgement

1. **Fixing the export rather than the figure is the whole value of this task.** The reviewer saw one
   letter-spaced diagram; all five assets carried the attributes. A per-figure fix would have satisfied
   the review and left the book still relying on the renderer's guessed metrics.
2. **Recording `postprocess` in the manifest is what makes it a build input.** A future change to the
   strip now invalidates every asset the same way a renderer change does, and the check mode re-renders
   and re-strips before comparing, so a hand-edited asset cannot pass.
3. **Removing the waivers instead of fixing the rule's scope is correct.** The book now has no standing
   exceptions in `check-waivers.txt`, which was the point of making debt visible: debt paid.
4. **The stale `check.sh` default was worth catching here.** It is tangential to figures, and it is the
   kind of defect that never fails a build — it fails a *person* who runs the check by hand and believes
   the result. One line, fixed and recorded.

## Where this task is weak

- **Verification is mechanical, not visual.** I proved the attributes are gone and each figure renders
  one number; I did not compare label widths in rendered images, so the claim "readable at 100%" rests
  on the causal story rather than on measurement.
- **The other figures' content was not reviewed.** Typography only; V2/V3/V9/V10/V12 and the
  lease-versus-fence figure belong to the pedagogy task.
- **`schema_version: 2` of the manifest is a breaking shape change** for any consumer outside this
  repository, and nothing here reads version 1 defensively.
- **The tracking value (0.3pt) is a judgement**, not a measurement, and a different reader might prefer
  0 or 0.5.

## Required tag

`repo/book-figures-v2`, created by `queue integrate` on the integrated state.
