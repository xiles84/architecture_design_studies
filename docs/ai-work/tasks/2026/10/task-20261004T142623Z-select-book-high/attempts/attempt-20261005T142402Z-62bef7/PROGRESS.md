# Progress — HIGH SELECT book publication

**Task:** `task-20261004T142623Z-select-book-high`. **Attempt:** `attempt-20261005T142402Z-62bef7`,
claim `claim-dc61831e9fef4e72` (epoch 1, real clock).
**Agent:** Claude Opus 5.5 (`claude-opus-5-5`) / Claude Code desktop / HIGH / executor; effort xhigh.
**Worktree:** `.worktrees/select-book-high`; branch `book/select-query-templates-high`.

Inputs:

- the completed synthesis (`study-01/select-v6-analysis`) and its book handoff
  `docs/select-query-templates/BOOK-HANDOFF-20261005.md`;
- the 18 `v6-sel-*` claims;
- `book/STYLE.md`;
- the join handoff and the order-correction qualification.

No database or benchmark was started.

## Decision Log

1. **All SELECT results go into the Charity case**, as variations after the recency variation.
   They are not a new case chapter: STYLE maps SELECT theory to Study 01 and the brief says
   "existing Charity donations case". The existing label `case-charity-join-plan` is kept for
   inbound links.
   - Confidence: High.
2. **Fixture honesty.** The join, CTE and placement variations name Study 01's 200,000-donation
   table. Membership, aggregate and pagination name the generated SELECT fixture standing in for
   donors and donations (10,000 or 12,000 rows), as its claims say.
   - Confidence: High.
3. **Numbers are copied from claim statements.** No ratio is recomputed and no summary range is
   merged across claims. Example: YugabyteDB membership is printed as "1.10–1.16x at fanout 1 and
   0.95–1.01x above it", as the claim states it.
   - Confidence: High.
4. **The template guide sits in the case** (`case-charity-select-guide`), not in theory, because it
   carries observed boundaries. Its claim ids are short raw text; the full claims are linked from
   the badges above it.
   - Confidence: Medium-High.
5. **Theory stays number-free.** Each SELECT strategy's "In practice" link now points at its
   measured variation, with `measured: true`. The pagination example gains a number-free
   access-path note (OR, UNION ALL overlap, expression index). The answer-contract strategy points
   at the guide.
   - Confidence: High.
6. **Badge label fix.** `repeated_controlled` rendered as "repeated within one run". That is false
   for the four SELECT claims at this strength (3–12 whole runs), and no earlier claim uses it. The
   label is now "repeated — no prior prediction", and the reading guide's legend gains the row.
   - Confidence: High.
   - Alternative: re-grade the claims, which would overstate them.
7. **Stale placement wording corrected.** The cache case's gap callouts and the glossary entry for
   colocation no longer say placement was never verified anywhere; they point at the Charity
   placement variation. This matches the narrowed `v4-gap-01`/`v4-gap-02`.
   - Confidence: High.
8. **Edition bumped** to "Edition 3 preview 4" with a changelog entry. No figure was added: tables
   and the existing theory pictures carry the variations. FIGURES.md is unchanged.
   - Confidence: Medium-High.
9. **The stale placement figure was regenerated, not captioned around.** Its source said
   tablet/leader placement was "UNVERIFIED in this corpus". The PlantUML source now names the
   verified Charity runs, and `book/assets/export.sh --write` re-rendered it with the pinned
   renderer. Only that asset changed.
   - Confidence: High.
   - Alternative: a caption caveat, which would leave a false figure label.
10. **Layout fixes after the first rendered review:** badges after a manual line break, an
    unbreakable expression, and an unjustified guide table with unbreakable ids. The canonical build
    was then rerun from a clean committed tree (`bf03e408`). Build 1's log, text and proofs are kept.
    - Confidence: High.

## Status

The final canonical build verified: 98 pages, 49/49 claims, 4/4 fonts. The reader audit passes
34/34, and the changed pages were inspected. Submitted for labelled HIGH self-review.
