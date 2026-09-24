# Result — text and terminology pass

**Task:** `task-20260924T110000Z-book-text-pass-v2`
**Branch:** `book/text-pass-v2`
**Required tag:** `repo/data-architecture-book-text-pass-v2`
**Capability:** HIGH session executing a LOW task (model `deepseek-flash`, Deep Code CLI; effort not exposed)
**Benchmark:** none — artifact-only

## What was delivered

**Taxonomy and structure (R05, R06).** Page 3 described an evidence system the book no longer has: it
claimed every family, variant, topology and scenario section follows one structure, and reduced the
whole gap vocabulary to "a gap callout is missing evidence". The page now defines *direct*,
*mechanism*, *analogy*, the four `gap_kind` values (schema limitation, coverage gap, instrument gap,
unstable measurement), *partially superseded* and *retired*, and scopes the structure claim to
major-family sections, saying why variant, scenario and concept chapters differ.

**Claims scoped to their evidence (R08–R12).**

- The topology chapter scopes its consensus sentence to the measured YugabyteDB RF=3 topology, naming
  asynchronous replication as the case it does not describe, and calls the seat-map result a layout
  *package* comparison rather than "one topology change".
- Normalized facts: normalization decides where a fact lives; it does not by itself prevent
  overbooking or double-selling, which needs a declared constraint and a concurrency strategy.
- Derived state: same-database rollups are updated in one transaction, and "publication" is reserved
  for a copy crossing the authoritative boundary.
- Append-only history claims better *answerability*, not better correctness, and names concurrency,
  ordering, duplicates, invalid events and atomicity as things it does not settle.
- "writes are cheap" is scoped to the measured cells rather than stated as a law.

**Vocabulary (R19–R22, R25, R29).**

- `read score` is defined at first use: a blended within-run composite, not the throughput of any one
  query, not comparable across runs.
- An *in-run repeat* is defined as another measurement window inside one benchmark process, explicitly
  not an independent run replication, so "Trials: 1" beside it cannot mislead.
- Numeric ranges in the generated cards are typeset with en dashes by `typeset-ranges`, so prose and
  cards agree while the registry text stays verbatim ASCII.

**A build defect found and fixed.** `book/assets/export.sh` recorded `source_revision` as
`git rev-parse HEAD:<source>`. Run with `--write` while a source change was uncommitted it stored the
*previous* blob, and the next `--check` — after the commit existed — computed a different value and
failed the build on a figure nobody had touched. It now records the committed blob only when that blob
describes the file's content, and `uncommitted` otherwise.

## Verified from committed files

```text
book/build.sh --out dist/preview.pdf       51 pages, 29/29 claims indexed, 6/6 fonts embedded,
                                           evidence digest present, banner absent, no unresolved tokens
book/assets/export.sh --write then --check both succeed and agree (the source_revision fix)
check.sh --book <absolute>                 OK
page-3 legend as rendered                  scope, direct/mechanism/analogy, the four gap kinds,
                                           partially superseded and retired all present
card ranges                                4–8x ×7, 33–35% ×12, 2.2–3.5x ×9, 4.5–11x ×6;
                                           ASCII "4-8x" and "33-35%" count 0 in card text
```

## Where this task is weak

- **Two acceptance criteria are open-ended by nature**: "family summaries separate the mechanistic
  trade-off from what the corpus observed" and "forward references carry a one-sentence definition".
  I fixed the instances I could locate and judged; a reader with a different threshold would find more.
- **Harness terminology is defined at first use in the chapters, not yet centralised.** The glossary
  belongs to the pedagogy task, which owns that structure; this pass added the definitions that a
  reader meets before any glossary would help.
- **The range fix is a display transform.** A reader extracting the PDF and comparing text to the
  registry will see an en dash where the JSON has a hyphen. That is deliberate and documented, and no
  digit changes.
- **`read score`'s formula is still not shown** — the claim is a composite and the book does not have a
  published weighting. I defined what it is and what it is not, and did not invent a formula.
- **Visual QA was textual**: fragment lines and the rendered text layer, not image inspection.

## Required tag

`repo/data-architecture-book-text-pass-v2`, created by `queue integrate` on the integrated state.
