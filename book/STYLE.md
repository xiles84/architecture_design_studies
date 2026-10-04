# Book style — theory, then applied study cases

Owner clarification, 2026-10-04. This is the single didactic contract.

## Reader hierarchy

| Part | Topic / case | Subdivision |
|---|---|---|
| I — Theory | Database schemas | Normalized rows · indexes · rolldown · rollup · embedding |
| I — Theory | Database SELECTs | Answer contracts · joins/filter placement · CTE policies · membership · aggregation/top-N · pagination |
| I — Theory | Concurrent writes and time | Hot rows · optimistic · pessimistic · history · expiry |
| I — Theory | Caches | Cache-aside/freshness · invalidation · leases/fencing |
| I — Theory | Queues | Asynchronous handoff · duplicate delivery/idempotency |
| I — Theory | Cluster placement | Colocation and distributed access |
| II — Study cases | Charity · ticket booking · reserved seating · configuration · cache consistency | Variations mapped to theory, thrive/die situations and summary numbers |
| Appendix | Reading results · reproduction · glossary · full evidence cards | Audit detail kept outside teaching chapters |

The contents exposes all three levels: part → topic/case → strategy/variation.
No empty promised chapter or duplicate older edition appears in the reader path.
Historical sources and immutable evidence remain preserved in the repository.

## Strategy lesson

Target 2–3 pages, with the Rollup pilot's depth and a repeated sequence:

1. Key idea, at most 25 words; one concrete problem.
2. A mechanism picture with an adjacent short text equivalent.
3. Three to five short steps.
4. Side-by-side fit and failure conditions.
5. Five cost rows: reads, writes, storage, correctness and operations.
6. A link to the matching case, or an explicit unmeasured coverage entry.
7. The eight-question checklist from methodology 21, with honest evidence status.

Theory teaches mechanisms and expected boundaries. It does not contain benchmark numbers.
A plausible failure condition is not labelled as an observed crossover.

## Applied case

Target 2–4 pages per business case; useful variation sections may extend a case.

- Establish the business, the invariant and a small model diagram.
- Map each variant to a named theory strategy using a link.
- Show the situation in which it thrives or dies.
- Distinguish **observed**, **inferred** and **not measured**.
- Put compact summary numbers beside the relevant claim badge.
- State the material limit and keep missing coverage visible.
- Give the study path; full run and input provenance lives on the linked card.

A summary table uses two or three generous proportional columns. Do not compress long
sentences into narrow identifier columns. Case boxes are optional when a table explains better.

## Short writing and visuals

1. One idea per paragraph: at most three sentences and 60 words.
2. Prefer labelled steps, small tables and short callouts for dense explanations.
3. Aim for sentences under 25 words.
4. Explain a thing once; link to the explanation elsewhere.
5. New or complex concepts need a diagram showing ownership, sequence or failure boundary.
6. Figures carry form, evidence status and a brief text equivalent (methodology 15).
7. Native Typst diagrams are conceptual illustrations, with their source included in the build hash.
8. No reader exercises. Checklists are a reading aid.

## Evidence

- Use inline `claim-badge(id)` citations in teaching chapters.
- Every active registered claim remains cited and has a full card in the appendix.
- Copy summary values from cited claim statements or limits; do not calculate new ratios.
- Do not quote points the claim says did not reproduce.
- No numeric ranking across studies or environments.
- Queues have no measured case. New SELECT results are not promoted without registry review.
- Do not edit signed analyses, reports, producing tags or frozen evidence packages.
- Owner clarification, 2026-10-04: LOW performs the join/CTE study; HIGH reviews and adds its results to the book.
- Preserve all original registered figure provenance and canonical build gates.

## Release check

Build from committed clean sources through `book/build.sh`; inspect rendered pages,
contents and cross-links. Record the source commit, PDF hash, evidence coverage and
remaining scientific gaps. A coherent reading structure does not close measurement gaps.
