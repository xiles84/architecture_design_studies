# Review — task-20260924T123503Z-book-review2-delta

**Reviewer:** HIGH session, model `deepseek-flash`, Deep Code CLI (effort not exposed)
**Verdict:** approve
**Disclosure:** self-review, against the committed sources and the rebuilt PDF.

## What was checked

| Check | Method | Result |
|---|---|---|
| Capsules apply to repeats only | rendered text | 34 capsules; the registry index still renders full cards |
| The first citation is still full | read two chapters in order | the first `registry-card` occurrence renders the card |
| Every limit survives in a capsule | source and page | `limits.join("; ")`, not `at(0)` |
| The counter cannot drift | architecture | Typst state in document order; no build-generated index exists |
| The dead input is gone | grep for `citation_index` | none |
| Lifecycle rendering | rendered text | Moved to / Still open in / unchanged from, all present |
| The statement stays verbatim | diff the registry text against the card's rendering path | the card renders `c.statement` unchanged |
| The legend is derived | source and page | reads `gap_kind` values from the active claims; panics on an undescribed kind |
| The legend order is deterministic | rendered text | schema limitation first, then coverage, instrument, unstable measurement |
| The ladder no longer overclaims | source | no "fastest read", no "only step that adds a consistency contract" |
| The new figure is proper | registry, manifest, page | source, entry, Figure 6 caption, text equivalent, no number of its own |
| No regressions | `book/build.sh`, `export.sh --check`, `check.sh` | 54 pages, 29/29 claims, 6/6 fonts, digest present, no banner, no unresolved tokens; exports verify; no waivers |

## Judgement

1. **Using Typst state rather than a build-time index is the decision that matters here.** The scan-based
   version would have worked on the day it was written and would have gone wrong silently the first time
   someone reordered an include. The page's own order cannot disagree with the page.
2. **Carrying all limits in the capsule is the difference between a summary and a smaller card.** The
   reviewer's sketch was compact because their example claim had one limit; several of this corpus's
   claims have three, and a capsule that shows one would quietly break the book's central promise.
3. **The generated legend closes the same class of defect as the retired gap claims**, one layer up: prose
   restating data. It now fails the build instead of going stale.
4. **The lifecycle split removes a real ambiguity.** A reader who saw "partially superseded by v4-gap-01"
   above a statement that still called placement absent could not tell what was open; the card now says.
5. **The lease-versus-fence figure earns its place** because the chapter argues a distinction a reader
   cannot easily see — it shows the expiry and the rejection side by side.

## Where this task is weak

- **The capsule loses the confound register**, which the full card carries. Following the pointer is now
  required for a reader who meets a confounded number at a repeat citation.
- **"Moved to successor claims" renders bare ids**, with no strength or title.
- **Tranche 2 of the visuals is still outstanding** — V2, V3, V4–V6, V9, V10 and V12 — and the page-count
  effect of six more figures has not been assessed.
- **The two-page saving is modest.** A capsule that condensed the limits would save more, at the cost of
  rewriting evidence text, which this task deliberately does not do.
- **Visual QA was textual** once more: rendered text and captions, not images.

## Required tag

`repo/data-architecture-book-review2-delta`, created by `queue integrate` on the integrated state.
