# Result — cache chapter: cards instead of a collapsed table, and the terminology fixes

**Task:** `task-20260924T105000Z-book-cache-chapter-v2`
**Branch:** `book/cache-chapter-v2`
**Required tag:** `repo/data-architecture-book-cache-chapter-v2`
**Capability:** HIGH session executing a LOW task (model `deepseek-flash`, Deep Code CLI; effort not exposed)
**Benchmark:** none — artifact-only; no database, no measurement

## What was delivered

1. **The three-column mechanism table is gone.** `book/concepts/cache-consistency.typ` now presents
   one card per mechanism via a new `mechanism-card` component in `book/lib/config.typ`, each carrying
   Race closed / Mechanism / Guarantee / **Does not solve** / Evidence. The "does not solve" line is
   mandatory precisely because five of these can be present and the contract still be broken.
2. **The mechanisms are grouped by layer (R16).** *Source correctness* — database mutation lock and
   conditional update — arbitrates the mutation. *Copy correctness* — invalidation, publication
   fence/CAS, source version validated at the hit boundary, fill lease, TTL/early expiry,
   transactional outbox — orders what leaves the transaction.
3. **Every evidence-status class the table carried is preserved.** The fill lease is still "newer
   signed but unregistered evidence"; the outbox is still "proposed / unmeasured" and visibly so; the
   lock, fence, TTL and equal-total items still name their registered claims.
4. **R13.** The chapter no longer says "invalidation is a fence". Invalidation clears an entry; it
   cannot stop an in-flight reader publishing an older value; the publication fence is the extra
   check. The stale-fill race the chapter already draws is named as the counter-example.
5. **R15.** The strict-after-acknowledgement contract is stated formally: with `W1`, `ack(W1)`, `R1`
   and `begin(R1)`, `begin(R1) > ack(W1)` implies `R1` must not return a version older than `W1`,
   unless a later write supersedes it.
6. **R14.** The harness metric is named and defined: `wrong_reads` counts reads that violated the
   strict comparator over committed-but-older values — not dirty, torn, uncommitted or impossible
   ones — and the bounded-window question is flagged as unanswerable without a bounded contract.
7. **R17/R18.** The outbox card states that a durable notification path is a prerequisite for
   strictness, not a guarantee of freshness, and a purpose matrix gives, per mechanism, whether it
   addresses source correctness, freshness proof, load control or recovery.
8. **`CONTEXT.md` and `LESSONS_LEARNED.md`** updated, including the measuring technique below.

## The defect, measured rather than described

Orphan-hyphen fragment lines in the extracted page text — the signature of the collapsed column,
where Typst hyphenates because the column is narrower than a word:

| | Edition 1 (committed draft) | After this change |
|---|---|---|
| Orphan-hyphen lines | 99 | 30 |
| Worst page | 33, with 17 | 34, with 3 |
| Book pages | 49 | 50 (the cache chapter itself loses three) |

The remaining 30 are ordinary line-end hyphenation distributed through the book. The specific
fragments the review quoted — `se - ri - alises`, `asyn - chro - nous`, `oth - ers`, `suf - fi - cient` —
no longer occur.

## Validation

```text
book/build.sh --out dist/preview.pdf     50 pages, 29/29 claims indexed, 6/6 fonts embedded,
                                         evidence digest present, unverified banner absent — verified
book/evidence/check.sh                   OK (2 waived figure-number defects, owned by book/figures-v2)
rendered card text                       Race closed / Mechanism / Guarantee / Does not solve / Evidence
                                         present for every mechanism, under the two group headings
fragment-line measurement                before 99 (worst 17) → after 30 (worst 3)
```

## Conclusion

The book's most conceptual section is now readable at normal size, and the chapter states what each
mechanism does not solve rather than implying that naming a mechanism closes the problem. No number,
registry claim, report or analysis was touched; the only files changed are the chapter and the shared
visual system.
