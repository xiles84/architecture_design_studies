# Result — cache chapter: card corrections and the relaxed taxonomy

**Task:** `task-20260924T123502Z-book-cache-chapter-v3`
**Branch:** `book/cache-chapter-v3`
**Required tag:** `repo/data-architecture-book-cache-chapter-v3`
**Capability:** HIGH session executing a LOW task (model `deepseek-flash`, Deep Code CLI; effort not exposed)
**Benchmark:** none — artifact-only; no database, no measurement

## What was delivered

**Cards** (all six defects the reviewer listed, four of them introduced by this project's own card rewrite):

1. The lock card is **split**. *Pessimistic lock + invariant recheck* says the lock supplies
   serialisation and that rejecting the obsolete operation is the recheck's contribution.
   *Conditional update / compare-and-set* puts the expectation into the write, so a zero-row update
   reports the loser. The combined "at most one of the two writers commits" guarantee is gone — it was
   true only of the two mechanisms together.
2. The **fill lease** suppresses duplicate fills *while a valid lease is held*, and says duplicate work
   can reappear if the lease expires or ownership is lost. "Load coordination, not freshness" is kept.
3. The **TTL** bounds *eligibility to serve*, not refill load; early expiry spreads work and is not a
   hard cap. The matrix cell matches.
4. The **outbox** is conditional, not universal: when coherence depends on asynchronous publication it
   makes the committed change durably observable, and strictness can also come from synchronous source
   validation or an authoritative read.
5. The **source-version** card states its prerequisites — every relevant mutation advances the token
   atomically, including deletes and external writers, and the validating read is authoritative enough
   for the declared contract.
6. **Invalidation** is scoped to lookups after it is applied, with in-flight reads and fills explicitly
   outside it; the **fence** card scopes its before/after writer protocol to the one that was measured.

**Structure:**

- `"Each one closes exactly one race"` and `"only one of them is a performance control"` are removed.
  Two mechanisms are load controls and the TTL closes no race at all; the heading is now "one primary
  failure mode each".
- The purpose matrix gained a **Change observation** column, so the outbox is no longer filed under
  mutation correctness. Headers are short so no word breaks in a narrow column, and a sentence defines
  each dimension.
- **"Relaxed" is a family**, not a definition: strict-after-acknowledgement, session/read-your-writes,
  bounded staleness (≤ Δ), eventual freshness, best effort. The measured cells are labelled
  **unbounded-relaxed**, which removes the contradiction with the `wrong_reads` explanation.
- §15.7's "stayed correct only because" became an observation plus the mechanism those cells are
  consistent with, and states that no unvalidated owned arm was run, so it is not an ablation.
- The `v2-16` prose keeps the read result and labels the write-path fencing as mechanism, noting the
  claim quantifies no write-path cost.

## Verified from committed files

```text
retracted phrasings in the source        none remain (one race, only one performance control,
                                         bounds residency, stops being served, one filler at a time,
                                         at most one of the two writers, stayed correct only because)
book/build.sh --out dist/preview.pdf     51 pages, 29/29 claims indexed, 6/6 fonts embedded,
                                         evidence digest present, banner absent, no unresolved tokens
matrix header as rendered                Mechanism | Mutation | Observation | Freshness | Load | Recovery
orphan-hyphen fragment lines             29 total, 3 on the worst page — the six-column matrix
                                         introduced no collapsed column
```

## Judgement on the reviewer's findings

All six card findings, the two summarising sentences, the matrix classification, the relaxed
contradiction, §15.7 and the `v2-16` prose are implemented as asked. I have no disagreement with any of
them: the split lock/CAS point is the sharpest of the six, because the original card's guarantee was
true only of the combination and the card presented it as the lock's.

## Where this task is weak

- **The card format still makes summaries look authoritative.** The "Does not solve" line is the
  mitigation, not a fix: a reader who skims only the Guarantee lines still meets confident sentences.
- **The taxonomy list adds a page.** The chapter grew from 50 to 51 pages; the release task will confirm
  the final count.
- **Visual QA was textual.** I measured fragment lines and read the rendered header and matrix in the
  text layer; I did not inspect every page as an image.
- **The reviewer has not seen this revision.** Their round-3 list is the specification I worked to, and
  the release artefact will be their next look — I did not invent additional corrections beyond it.

## Required tag

`repo/data-architecture-book-cache-chapter-v3`, created by `queue integrate` on the integrated state.
