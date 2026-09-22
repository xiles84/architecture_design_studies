# Result — book release review v1

**Task:** `task-20260922T025912Z-book-release-review-v1`
**Branch:** `repo/book-release-review-v1`
**Required tag:** `repo/data-architecture-book-v1`
**Capability:** HIGH session (reviewer + integrator) / executor (model `deepseek-flash`)

## Release candidate (final, clean tree)

```text
artefact:               book/dist/data-architecture-reference.pdf
pdf sha256:             58d61df0600a46e9b019aa83067475add4e3a2b986345c6041deeea3c8516cb1
pages:                  36
fonts embedded:         4/4
claims indexed:         23/23
evidence digest in pdf: true
uri annotations:        1
source commit:          cc5ad3ff95232a52c558b4d5e40973233e214ae2 (tree dirty=false)
typst image digest:     sha256:032e292249bcd378480cc7c142cfa324b63ef8aadeb88d7e7230320c4c9c422f
evidence registry:      book/evidence/v2/claims.json (23 claims, 0 errors)
```

## Acceptance checks

| Check | Result |
|---|---|
| Evidence registry valid in Podman | `evidence: 23 claims, 0 errors, 23 book inputs` |
| Every number resolves to a claim | 54 distinct numeric tokens; 0 not found in the v2 registry |
| Every claim cited in prose | 23/23 claim ids appear in chapter/concept prose (after correction 1) |
| Direct-evidence / analogy boundary | direct callouts on measured results; the supermarket/auction/e-commerce transfers and the native-document store transfer are labelled analogy |
| Cross-study comparisons | none; the comparisons chapter states the refusal and lists all six confounds |
| Topology qualifications | shared cores, throttled nodes and no real network stated in topologies, every family boundary and the provenance page |
| Failed-cell handling | the Study 02 failed cells are described in the family/variant text and the run-level notes; no failed cell is a winner |
| Superseded material | the book reads v2 only; v1 is never cited |
| Beginner path | `how-to-choose` gives the answerability-first decision path and a worked decision |
| Expert path | concept chapters plus evidence-and-reproduction give mechanism and reproduction |
| PDF integrity | committed blob sha256 equals the manifest hash; fonts embedded; registry digest present |
| Visual inspection | title, provenance and a content spread rendered and inspected; three defects found and fixed |

## Defects fixed in this task

1. Coverage gaps `v2-gap-03` / `v2-gap-04` cited in prose.
2. Title-page provenance no longer renders a literal "n" and is not justified.
3. The blank front-matter page removed (37 → 36 pages).

## Release statement

The v1 book, `book/dist/data-architecture-reference.pdf`, is accepted as the first released artefact
of the data-architecture reference. It is evidence-bound: every figure comes from one of the 23
active v2 claims, with confounds, gaps and analogies on the page. It is not a final edition: the
evidence is single-run and single-host, and the coverage gaps cited in the scenarios chapter are
open work. No measurement was performed during this task.
