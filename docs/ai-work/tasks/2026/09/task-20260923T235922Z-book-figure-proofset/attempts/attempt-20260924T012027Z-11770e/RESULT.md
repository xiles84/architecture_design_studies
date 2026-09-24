# Result — book figure proof set

**Task:** `task-20260923T235922Z-book-figure-proofset`
**Branch:** `book/figure-proofset`
**Required tag:** `repo/book-figures-v1`
**Capability:** HIGH session executing the task / analyst (model `deepseek-flash`)
**Benchmark:** none — authorship and build only

## What was delivered

**Four book figures**, each with two independent labels (form; evidence status) and an adjacent text
equivalent:

| # | Figure | Source | Form | Evidence status | Placed in |
|---|---|---|---|---|---|
| 1 | placement and ownership | `book/assets/sources/fig_placement_ownership.puml` | structure | conceptual illustration | `concepts/information-placement.typ` |
| 2 | two buyers, same marked seats | `studies/03-reserved-seating/diagrams/s_arbitration.puml` (reuse) | sequence | implemented design contract (labelled negative control) | `chapters/scenarios.typ` |
| 3 | the stale-fill race and the fence | `book/assets/sources/fig_cache_stale_fill.puml` | sequence | observed result — illustrates `v2-15`, embeds no rate | `concepts/cache-consistency.typ` |
| 4 | expiry and clock authority | `studies/03-reserved-seating/diagrams/e_expiry.puml` (reuse) | sequence | implemented design contract (labelled negative control) | `concepts/expiry-and-clock-authority.typ` |

Also: the `figure-evidence` helper (`book/lib/config.typ`), methodology rule 15, and `book/FIGURES.md`
updated with the proof set and the reader-check record.

## Two-reader check (recorded)

```text
reader 1 (no caption): what does the picture say?
reader 2 (sceptic):    does it match the claim/SQL? does it imply a performance number?
independence:          both passes were the same model — weaker than two independent readers
```

| Figure | Outcome |
|---|---|
| 1 placement/ownership | kept — five placements and the ownership note are readable without the caption |
| 2 hold arbitration | kept — shows the contract plus the labelled `S0` control; the text equivalent states it carries no measured value |
| 3 stale fill | **redesigned** — the first version showed one reader; claim `v2-15` is about instances sharing a cache, so a second instance's reader was added to hit the stale fill |
| 4 expiry | **relabelled** — form corrected from `state` to `sequence`; the study source is a sequence and the book does not hold an editable copy |

No figure was dropped. Figures 1–3 were visually inspected as rendered PNGs.

## Validation

```text
book/assets/export.sh --check                  5 figures, assets and manifest verified
book/assets/export.sh --write                  after sources were committed, so revisions are real
book/build.sh                                  exit 0
  pages                                        49
  fonts embedded                               5/5
  claims indexed                               29/29 (v3)
  evidence digest in PDF                       true (814f40949d162f56…)
  figure_registry / figure_manifest / sha      recorded in dist/build-manifest.json
  figures_total                                5 (4 book figures + pilot)
  source_dirty                                 false
pdf sha256                                     d6d3c0299f96c9a74f2ea8f52517562ae1993409e2e4dc1eb56e67e589ac9ada
```

Acceptance criteria:

- Three to four labelled figures including a book-level structure figure, a Study 02/03 sequence reuse,
  a new Study 05 stale-fill sequence and an expiry figure — met (four).
- Two independent labels plus an adjacent text equivalent — met.
- An observed-result figure cites an active claim and embeds no performance number — met (figure 3,
  `v2-15`, no number in the SVG).
- Two-reader check recorded with keep/redesign/drop — met (both passes one model, stated).
- Methodology figure-evidence rule added — met (rule 15).
- PDF rebuilds with every figure source and embedded asset in the manifest — met.

## Limitations

- The reader check is not independent (one model, two passes).
- Figures 2 and 4 reuse study sources; a later study-source change requires a re-export and fails the
  build until then.
- The rebuilt PDF is a draft; it has not been release-reviewed.
