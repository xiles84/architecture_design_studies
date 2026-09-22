# Result — book synthesis v1 (draft)

**Task:** `task-20260922T025912Z-book-synthesis-v1`
**Branch:** `repo/book-synthesis-v1`
**Required tag:** `repo/data-architecture-book-v1-draft`
**Capability:** HIGH session (author) / executor (model `deepseek-flash`)
**Benchmark:** none — authorship and a documentation build only

## What was delivered

The v1 prose for the data-architecture reference, written into the existing `book/` Typst shell:

- **`chapters/major-families.typ`** — the six fixed families (normalized facts, rolldown/duplicated
  data, rollups/materialized aggregates, embedded documents, append-only history/snapshots,
  external read copies/caches), each with quick choice, thrives/perishes, scenarios, mechanism,
  read/write/storage/correctness/operations costs, direct evidence, boundaries and reproduction.
- **`chapters/how-to-choose.typ`** — the answerability-first decision path, arbitration, derived-state
  maintenance and the evidence-level rule.
- **`chapters/minor-variants.typ`** — the measured controlled pairs and the variants that are named
  but not controlled.
- **`chapters/cross-family-comparisons.typ`** — what a controlled pair requires, the comparisons the
  book refuses to make, and the six confounds published in the chapter.
- **`chapters/topologies.typ`**, **`chapters/scenarios.typ`** — the topology axis and the domain
  scenarios, with transfers labelled analogy.
- **`concepts/*.typ`** — eight mechanism chapters: information placement, indexes and hot rows,
  concurrency control, derived state and history, expiry and clock authority, cache consistency,
  distributed placement, reading benchmarks honestly.

Every figure is drawn from one of the 23 active `book/evidence/v2/claims.json` claims and rendered
through its evidence card; confounds and coverage gaps appear beside the numbers, and no
cross-study numeric ranking, native-family transfer or capacity claim is made.

## Verified build (pinned image, clean tree)

```text
pages:                  37
fonts embedded:         4/4   (Libertinus Serif regular/bold/italic, DejaVu Sans Mono)
claims indexed:         23/23 (every v2 claim id appears in the rendered text)
evidence digest in pdf: true
uri annotations:        1
pdf sha256:             8490e4d89e25876a974956583648852a38dc26cf9996ecaa427fad54ffaad59a
source commit:          6b2bc53e58a15317589b27716d6e2f95d796c2b3 (tree dirty=false, built 2026-09-22T19:15:09Z)
typst image digest:     sha256:032e292249bcd378480cc7c142cfa324b63ef8aadeb88d7e7230320c4c9c422f
```

`book/build.sh` fails the build if the PDF is implausibly short, any font is not embedded, the
registry digest is missing from the text, or a claim id is missing.

## Limitations

- This is a v1 draft: it selects and explains committed claims; it re-derives no number from result
  JSON and runs no measurement.
- Every family result is single-run and single-host; the environment caveat is repeated wherever a
  conclusion is drawn.
- The scenario set is deliberately short; additional domains belong in a later edition as labelled
  analogies.
- The draft PDF is submitted for independent HIGH review (`book-release-review-v1`); it is not the
  release artefact.
