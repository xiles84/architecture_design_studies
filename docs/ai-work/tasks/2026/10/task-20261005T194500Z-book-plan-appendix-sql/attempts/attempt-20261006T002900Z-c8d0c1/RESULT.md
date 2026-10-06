# Result — SQL before every plan (Edition 3 preview 7)

**Task:** `task-20261005T194500Z-book-plan-appendix-sql`. **Attempt:** `attempt-20261006T002900Z-c8d0c1`.
**Agent:** Claude Opus 5.5 / Claude Code desktop / HIGH / executor; effort xhigh.

## Delivered

- **Every plan is preceded by the SQL that produced it.** All 45 plans in the query-plan appendix
  (pp. 75–92) now show the bound SQL, the statement name and the parameters. The statements in
  node-tree checks show theirs too, once each.
- **The SQL comes from the code that ran it.** A new DB-free exporter,
  `selectbench/cmd/appendixsql`, rebuilds it with the harness's catalogue and binder.
- **`tools/evidence plans` checks it.** The SQL must match each cell, and every catalogue source
  must be byte-identical at the producing commit.

Validation and labelled self-review: `docs/book-plan-appendix-sql-20261006/VALIDATION.md`. The
canonical build verified 117 pages and 50/50 claims, PDF sha256 `4ddc80d8…`.

## Decision Log

1. **Rebuild the SQL with the study's own package, not a copy in `tools/evidence`.** The CTE-reuse
   statements exist only in Go, and a copied binder could drift.
   - Cost: a new `cmd` in Study 01; the study's harness code is unchanged.
   - Confidence: High.
2. **Prove the rebuild applies to past runs with a git byte-identity check, not an assumption.** If
   a later catalogue change breaks it, the generator fails, and the fix is to rebuild from the
   producing commit.
   - Confidence: High.
3. **Show SQL for checked statements once and reference it after.** This saves about two pages of
   repetition between the PostgreSQL and YugabyteDB exhibits.
   - Confidence: High.
4. **Keep the theory's 9pt SQL blocks unchanged.** The appendix uses its own compact box.
   - Confidence: High.
5. **Recommendation for future harness work: store the bound SQL in each result cell.** It is
   recorded in `LESSONS_LEARNED.md`; the harness was not changed here.
   - Confidence: High.
