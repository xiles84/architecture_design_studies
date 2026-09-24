# Review — task-20260924T123502Z-book-cache-chapter-v3

**Reviewer:** HIGH session, model `deepseek-flash`, Deep Code CLI (effort not exposed)
**Verdict:** approve
**Disclosure:** self-review. Every check below was run against the committed chapter and the rebuilt PDF.

## What was checked

| Check | Method | Result |
|---|---|---|
| Every retracted phrasing is gone | grep the source for the seven sentences | none remain |
| The lock is two cards | read the source-group block | "Pessimistic lock + invariant recheck" and "Conditional update / compare-and-set" are separate, with no combined singleton guarantee |
| The lease and TTL are scoped | read both cards and the matrix row | lease "while a valid lease is held"; TTL "eligible to serve", "no hard upper bound on refill traffic" |
| The outbox is conditional | read the card | states the asynchronous-publication condition and the two synchronous alternatives |
| Observation is its own dimension | rendered matrix | `Mechanism \| Mutation \| Observation \| Freshness \| Load \| Recovery`, outbox under Observation |
| The taxonomy renders | rendered text | "Bounded staleness (≤ Δ)" present; the cells are labelled unbounded-relaxed |
| §15.7 no longer claims causality | read the paragraph | observation + mechanism + an explicit "not as an ablation" |
| No new typographic defect | orphan-hyphen fragment lines | 29 with 3 on the worst page; the six-column matrix adds none |
| The book still builds clean | `book/build.sh` | 51 pages, 29/29 claims, 6/6 fonts, digest present, no banner, no unresolved tokens |

## Judgement

1. **The split is the substantive fix.** The old card gave one guarantee to two mechanisms that provide
   it only together. Naming the lock's contribution as serialisation, and rejection as the recheck's,
   is what the code in the corpus actually does — `c2_pessimistic_lock` versus `c1_optimistic_version`
   is precisely the pair the study measured.
2. **Weakening the lease and TTL claims makes them stronger, not weaker.** "Suppressed while a valid
   lease is held" is defensible and teachable; the previous "one filler at a time" was false under
   lease expiry, and the chapter's own fence discussion is why.
3. **The matrix column is the right repair.** Putting "durable observation" under mutation correctness
   was a category error; the outbox does not make the mutation correct, it makes the mutation's
   *notification* durable. The six dimensions now map onto the cards one to one.
4. **The relaxed taxonomy resolves a contradiction this project created.** Bounded staleness is now one
   member of a family, and the measured cells are honestly labelled unbounded-relaxed, which is what
   §15.7's `wrong_reads` paragraph already required.
5. **The two "does not solve" scopes are the highest-value edits.** The invalidation card now says
   in-flight reads and fills are outside its guarantee, and the fence card says the before/after writer
   protocol is the one that was measured — exactly the distinction the reviewer drew in §3.2 and §3.3.

## Where this task is weak

- **The chapter is one page longer** (50 → 51). The release task will confirm the final edition count,
  and the added taxonomy list is the cost of removing the contradiction.
- **The cards remain terse by design.** A reader who reads only the Guarantee lines still meets
  confident sentences; the mitigation is the "Does not solve" line, not a structural guarantee.
- **No reviewer has seen this revision.** The work implements their round-3 list exactly, and the
  release artefact is their next look; I did not add corrections beyond the list, so if the list itself
  had a gap, this task inherits it.
- **Visual QA was textual.** Fragment-line counting and reading the rendered header, not image
  inspection of each page.

## Required tag

`repo/data-architecture-book-cache-chapter-v3`, created by `queue integrate` on the integrated state.
