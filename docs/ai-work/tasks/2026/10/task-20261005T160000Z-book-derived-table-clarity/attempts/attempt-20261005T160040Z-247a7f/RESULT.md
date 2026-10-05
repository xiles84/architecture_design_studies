# Result — derived-table clarity (Edition 3 preview 5)

**Task:** `task-20261005T160000Z-book-derived-table-clarity`. **Attempt:** `attempt-20261005T160040Z-247a7f`.
**Agent:** Claude Opus 5.5 / Claude Code desktop / HIGH / executor; effort xhigh.

## Delivered

The book now answers the owner's question.

- **Filtered derived table:** it never outperformed. It ran the same plan as `JOIN … WHERE` in every
  measured pair (`v6-sel-01`).
- **When a derived table changes the plan:** only when it groups, deduplicates, limits or ranks rows.
- **Pre-aggregation:** the measured example is `v6-sel-18`, newly registered from existing evidence.
  - It lost at 1–1,000 parents on PostgreSQL.
  - At 10,000 it dipped below 1 only inside the drift.
  - It never won on YugabyteDB.
- **Where it appears:** the Charity case (pp. 56–57), the template guide (p. 60) and theory (p. 21).

Validation and labelled self-review: `docs/book-derived-table-20261005/VALIDATION.md`. Evidence
validate: 50 claims, 0 errors. The canonical build verified 99 pages and 50/50 claims, PDF sha256
`5648b273…`.

## Decision Log

1. **Register pre-aggregation rather than only mention it.** The book may print only registered
   numbers, and the synthesis evidence already holds the six runs.
   - Prediction verdict: `unresolved`, because the 10,000-parent advantage is inside the drift.
   - Confidence: High.
2. **The textbook pre-aggregation case is marked "not measured".** It needs two child collections
   joined at once, which this fixture does not have.
   - Confidence: High.
3. **Edition bumped to preview 5.** The verified build checks the edition on the title page.
   - Confidence: High.
4. **The task's `not_before` was set 22 minutes into the future by mistake.** `task.json` is
   immutable, so the claim waited until 16:00Z. No other effect.
   - Confidence: High.
