# Review — task-20260924T110002Z-book-pedagogy-v2

**Reviewer:** HIGH session, model `deepseek-flash`, Deep Code CLI (effort not exposed)
**Verdict:** approve
**Disclosure:** self-review, against the committed sources, the exporter and the rebuilt PDF.

## What was checked

| Check | Method | Result |
|---|---|---|
| The divider is a real opener | rendered page | "Part II — Concepts", lead and glossary pointer present |
| The glossary is the definition site | rendered text | all 26 terms present, including retired and partially superseded |
| Figure V1 added properly | registry + manifest + page | source, asset, manifest entry, Figure 1 caption, text equivalent |
| Figure V7 added properly | registry + manifest + page | source, asset, manifest entry, Figure 5 caption, text equivalent |
| V11 is a table, not an image | source | a `#table` in the reading chapter |
| No figure asserts a number | read both sources | neither contains a rate; both are labelled conceptual illustrations |
| No embedded figure numbers | the `figure-number` rule with no waivers | OK |
| Every registered figure is embedded | the `orphan-figure` rule | OK, seven figures |
| The exports are reproducible | `export.sh --check` | passes, postprocess reported for all seven |
| No regressions | `book/build.sh` | 56 pages, 29/29 claims, 6/6 fonts, digest present, no banner, no unresolved tokens |
| No new typographic defect | fragment-line count | 39 with 3 on the worst page |

## Judgement

1. **The glossary is written against the registry's vocabulary, not the prose's**, which is why
   `retired` and `partially superseded` have definitions at all — no chapter had ever named them. Naming
   the schema field for each shared term is the detail that keeps it from drifting.
2. **Both figures were kept free of numbers.** A decision tree that quotes a throughput would be a claim
   without a card; the text equivalent says where the numbers live instead.
3. **Making V11 a table is the right call and the reviewer's own suggestion.** A checklist that carries
   text is readable, searchable and cannot be letter-spaced by a renderer.
4. **Scoping the closed-loop sentence in the same file** avoids leaving a claim the reviewer flagged in
   round 3 sitting beside a checklist that tells readers not to trust closed-loop tails.

## Where this task is weak

- **Tranche 2 is deliberately not done**: V2, V3, V4–V6, V9, V10 and V12 remain. The brief required the
  page count to be reviewed first, and the book is now 56 pages from 51 — a second tranche of six
  figures would move that materially, which is a judgement for the release, not for this task.
- **The glossary coexists with the in-chapter definitions** rather than replacing them, so a term is
  still explained twice. That is deliberate (pointers beat removals) but it is duplication.
- **Figure quality was not inspected as an image.** I verified registration, captions, order and the
  absence of embedded numbers; label overlap after the attribute strip is still unverified visually.
- **The part opener's exact wording is mine**, not the reviewer's, and it makes a claim about the order
  of the chapters that a future editor could disagree with.

## Required tag

`repo/data-architecture-book-pedagogy-v2`, created by `queue integrate` on the integrated state.
