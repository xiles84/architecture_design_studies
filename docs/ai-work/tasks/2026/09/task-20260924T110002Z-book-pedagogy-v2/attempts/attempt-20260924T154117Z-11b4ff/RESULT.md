# Result — reader-facing framework: part opener, glossary, visual tranche 1

**Task:** `task-20260924T110002Z-book-pedagogy-v2`
**Branch:** `book/pedagogy-v2`
**Required tag:** `repo/data-architecture-book-pedagogy-v2`
**Capability:** HIGH session executing a LOW task (model `deepseek-flash`, Deep Code CLI; effort not exposed)
**Benchmark:** none — the only container started is the pinned PlantUML renderer

## What was delivered

**The part opener (R23).** The Concepts divider was a heading over an almost blank page. It is now
"Part II — Concepts" with a two-sentence lead naming the order the decisions bind in, and a pointer to
the glossary.

**The glossary (R24).** `book/glossary.typ` is new and is the single definition site for 26 canonical
terms: claim and the strength values; direct evidence, mechanism evidence, analogy; the four
`gap_kind` values; partially superseded and retired; run, replication, trial and the in-run repeat
distinction; `read_score`; `wrong_reads`; the strict contract and the relaxed family with bounded
staleness as one member; resource framing; colocation versus verified placement; rollup, rolldown and
embedding; source versus copy correctness; invalidation versus publication fence; confound; and
family/variant/scenario. Where a term is also a registry field (`gap_kind`, `trials`,
`resource_framing`), the field is named, so the prose vocabulary and the schema cannot drift apart
silently.

**Visual tranche 1 (V1, V7, V11).** Three additions, each with a form and evidence-status label and an
adjacent text equivalent:

- **V1**, the decision path, in Chapter 2 — seven steps from answerability, through arbitration, to
  stating the evidence level.
- **V7**, the freshness timeline, in Chapter 15 — where the strict contract binds, and the note that
  every relaxed contract is a named weakening of it.
- **V11**, the benchmark-reading checklist, as a **table** rather than a figure: eight questions, each
  pointing at where the answer lives on the card and what to do when it is missing.

Neither figure asserts a measured value of its own: one orders the book's recommendations, the other
draws the contract the claims are judged against. Both are book-owned sources under
`book/assets/sources/`, registered in `figures.json` and rendered by the pinned exporter.

**Also fixed in the same file:** Chapter 17's closed-loop tail was called "a floor, not an SLO";
closed-loop load can *understate* fixed-rate behaviour, which is defensible, and "floor" implied a
guaranteed lower bound.

## Verified from committed files

```text
book/build.sh --out dist/preview.pdf   56 pages, 29/29 claims indexed, 6/6 fonts embedded,
                                       evidence digest present, banner absent, no unresolved tokens
book/assets/export.sh --check          passes; seven figures, postprocess reported per asset
check.sh --book <absolute>             OK
rendered figures                       Figure 1 the decision path … Figure 6 the stale-fill race,
                                       one number each, no in-diagram number
new content rendered                   "Part II — Concepts", Glossary, "Before quoting a number",
                                       schema limitation ×12, publication fence ×5, in-run repeat ×2
fragment lines                         39 total, 3 on the worst page — the glossary adds no
                                       collapsed column
```

## Where this task is weak

- **Tranche 2 is not done.** The review's V2, V3, V4–V6, V9, V10 and V12 remain unimplemented; the task
  brief deliberately scoped this first tranche and required the page count to be reviewed before
  starting the rest. V8 is the mechanism cards, which Task H already delivered.
- **The glossary duplicates definitions that still appear in the chapters.** It is the canonical site,
  not a replacement: a later pass could shorten the in-place sentences and point here.
- **The figures' text equivalents are prose paragraphs**, which is the existing convention; they are not
  structured alternatives (no table of the same content).
- **Visual QA was textual.** I confirmed the figures render, are registered, have captions and appear in
  order; I did not inspect the rendered images for label overlap after the attribute strip, which is the
  one thing the figures task could not verify mechanically either.
- **The part opener's page is now a full page**, which changes the book's page count and the positioning
  of the contents; the release task will record the final count.

## Required tag

`repo/data-architecture-book-pedagogy-v2`, created by `queue integrate` on the integrated state.
