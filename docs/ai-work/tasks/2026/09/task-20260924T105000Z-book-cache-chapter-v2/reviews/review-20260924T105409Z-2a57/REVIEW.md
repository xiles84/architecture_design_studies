# Review — task-20260924T105000Z-book-cache-chapter-v2

**Reviewer:** HIGH session, model `deepseek-flash`, Deep Code CLI (effort not exposed)
**Verdict:** approve
**Disclosure:** self-review, as with the v4 registry task. The checks below were run against the
committed sources and the rebuilt PDF, not from memory.

## What was checked

| Check | Method | Result |
|---|---|---|
| The collapsed table is gone | read the chapter; grep for `#table(` in it | one table remains, the purpose matrix, with fractional columns and no prose cell |
| Every mechanism survived | diff the mechanism list against the table's rows | 7 mechanisms: mutation lock, invalidation, publication fence/CAS, source version at hit boundary, fill lease, TTL/early expiry, outbox (the table merged the lock and the conditional update into one row; splitting them is the R16 layering point) |
| Evidence statuses preserved | read each card's Evidence line | active registered claim ×4, newer signed but unregistered ×1, mechanism demonstrated ×1, proposed/unmeasured ×1 — the four classes are intact |
| The page is readable | count orphan-hyphen lines in the extracted text, before vs after | 99 → 30; worst page 17 → 3; the quoted fragments no longer occur |
| The book still builds clean | `book/build.sh --out dist/preview.pdf` | 50 pages, 29/29 claims indexed, 6/6 fonts embedded, digest present, banner absent |
| No registry or evidence file touched | `git diff --stat` on the branch | only `book/concepts/cache-consistency.typ`, `book/lib/config.typ`, `CONTEXT.md`, `LESSONS_LEARNED.md` |
| The source checker still passes | `book/evidence/check.sh` | OK, with the two known figure-number waivers |

## Judgement

1. **The redesign does what R01 asked and more.** The review's acceptance criteria were: no column
   narrower than ~35% of text width, no word broken into repeated 1–3 character fragments, each
   mechanism understandable without reading another row, page breaks not splitting a heading from its
   explanation. Cards satisfy all four structurally, and the fragment measurement confirms it
   numerically rather than by assertion.
2. **The "Does not solve" line is the substantive improvement, not decoration.** The table's evidence
   column answered "how strong is the evidence"; nothing answered "what remains broken if I use
   this". Two mechanisms (fill lease, TTL) are now explicitly load controls that prove nothing about
   freshness, which is the confusion the review's R18 identified.
3. **The grouping resolves a real conflation.** Putting a database row lock in the same list as a
   cache publication fence implied they operate at one layer. They do not: one arbitrates the
   mutation, the other orders a copy leaving the transaction.
4. **The R13 rewrite is correct and is backed by the chapter's own figure.** Invalidation clears an
   entry; the stale-fill race shows an older value being published afterwards. Saying "invalidation is
   a fence" described a mechanism that cannot close the race the same chapter draws.

## Where this task is weak

- **The page-count and fragment numbers are from one build on one toolchain.** They compare like with
  like (same extraction method, same pinned image) but they are not a typographic guarantee for every
  font or paper size.
- **Visual QA was textual.** I measured fragment lines from the extracted text and read the rendered
  card structure in the text layer; I did not render every page to an image and inspect it at 100%,
  which the plan lists as the ideal check. The measurement is a good proxy for this specific defect
  (a column too narrow to hold a word) but not a substitute for looking.
- **The remaining 30 orphan hyphens are unexplained individually.** They are ordinary line-end
  hyphenation spread across the book rather than clustered, which is why I treated them as normal,
  but I did not confirm each.
- **`book/dist/` is still the v3 build.** As with the previous task, the tracked PDF is regenerated
  once at release from merged sources.
- **Tooling friction worth recording:** in this WSL + Windows-podman setup the queue CLI's own git
  commits repeatedly leave `.git/index.lock` behind, which blocked the publish/claim of this task until
  the lock was removed by hand. It is noted here because a LOW executor hitting it will otherwise
  diagnose a phantom "another git process".

## Required tag

`repo/data-architecture-book-cache-chapter-v2`, created by `queue integrate` on the integrated state.
