# Progress — book release review v1

**Task:** `task-20260922T025912Z-book-release-review-v1`
**Branch:** `repo/book-release-review-v1`
**Attempt:** `attempt-20260922T192910Z-5bce5d` (claim `claim-7502c7fe87cb83e8`, epoch 1)
**Capability/role:** HIGH session (reviewer + integrator) / executor (model `deepseek-flash`)
**Required tag:** `repo/data-architecture-book-v1`

Artifact-only review and release. No database started, no measurement.

## What the review did

1. Read the synthesis Decision Log first (`book-synthesis-v1` attempt
   `attempt-20260922T191804Z-49e273`) and judged each decision.
2. Re-ran the evidence registry validation in Podman on the active v2 source.
3. Mechanically audited the book's numbers against the registry and the claim ids against the
   prose.
4. Visually inspected rendered pages (title, provenance, how-to-choose) and found three defects,
   corrected below.
5. Rebuilt and re-verified the final PDF; produced the release candidate on a clean tree.

## Corrections made (targeted, material to the reader)

1. **Uncited coverage gaps.** `v2-gap-03` and `v2-gap-04` appeared only in the evidence index; the
   scenarios chapter now cites both and explains how each scenario family stops. This closes the
   synthesis review's first recommendation.
2. **Title-page provenance rendered a literal "n"** after each sentence (a `\n` in Typst markup is
   not a newline) and the long hash stretched under justification. Replaced with explicit Typst line
   breaks and no justification.
3. **A blank page sat between the title and the provenance page.** The level-1 heading show rule's
   `pagebreak(weak: true)` fired on the provenance page's leading heading. The provenance title is
   now styled text (front matter, not an outline entry), and the blank page is gone: 37 → 36 pages.

## Decision Log

1. **Correct in this task rather than spawn a correction task.** The three defects are small,
   reader-visible formatting/navigation faults, not scientific or architectural problems; they are
   within a release review's remit and cheaper to fix now than to re-open the draft. Confidence: High.
2. **Keep the release candidate's manifest on a clean tree.** The `dist` cycle is commit → build →
   commit so `source_dirty=false` is truthful. Confidence: High.
3. **Do not re-derive numbers from result JSON.** That is the standing limitation recorded in the
   evidence registry; the release checks that every number resolves to a claim, not that the claim
   re-computes from the raw run. Confidence: High.

## Residual risks

- The book remains a v1 draft with single-run evidence and a short scenario set; the release tag
  marks the v1 artefact, not a claim of completeness.
- The two Study 04/05 protocol tasks that would close the cited gaps are still proposed.
