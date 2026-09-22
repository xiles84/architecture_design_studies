# Progress — book synthesis v1

**Task:** `task-20260922T025912Z-book-synthesis-v1`
**Branch:** `repo/book-synthesis-v1`
**Attempt:** `attempt-20260922T191804Z-49e273` (claim `claim-06f5b0ecca2083ab`, epoch 1)
**Capability/role:** HIGH session (author) / executor (model `deepseek-flash`)
**Required tag:** `repo/data-architecture-book-v1-draft`

Artifact-only authorship. No database started, no benchmark run, no measurement.

## Decision Log

1. **Write only from the 23 active v2 claims.** Every figure in the prose appears in a v2 claim's
   statement or support `expected` value; each section cites its claim ids through `registry-card`,
   so the page cannot drift from the data. v1 is not read. Confidence: High.
2. **Keep the fixed six families and the progressive structure.** Each family section carries the
   quick choice, thrives/perishes, scenarios, mechanism, costs, direct evidence, boundaries and
   reproduction sub-sections required by the brief. Confidence: High.
3. **Mechanisms live in the concept chapters; families point to them.** This avoids repeating the
   index/hot-row, concurrency, derived-state, expiry, cache-fence and placement explanations.
   Confidence: High.
4. **Confounds, gaps and analogies are on the page, not in footnotes.** `conf-01`…`conf-06` are
   listed in the comparisons chapter; every family states its coverage gap; the supermarket/
   auction/e-commerce transfers are labelled analogy. Confidence: High.
5. **Do not present JSONB as a native document-store comparison or Redis as an authoritative
   key-value comparison.** The embedding family is labelled an analogy to a native document
   database; the cache family states Redis was measured only as a cache. Confidence: High.
6. **Do not turn laptop topologies into capacity guidance.** The topologies and every run boundary
   repeat the shared-core, no-real-network caveat. Confidence: High.
7. **Numbers are kept with their caveats in the same breath.** Where a result is single-run or
   unsettled, the sentence that quotes it also states the limitation (the 128-buyer scope, the
   order-sensitivity, the upper-bound status), so a reader cannot take the figure without the
   caveat. Confidence: High.
8. **The draft is generated, not hand-set.** `book/build.sh` builds and verifies it in the pinned
   image; the build clock and image digest are injected. Confidence: High.

## Residual risks

- This is a v1 draft: it selects claims, it does not re-derive any number from result JSON.
- The scenario transfer list is short; more domains could be added, each as a labelled analogy.
- The book's relationship to each claim is by id; a future task could render a claim-to-page index
  automatically.
