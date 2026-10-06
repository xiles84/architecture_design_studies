# Edition 3 preview 7 — SQL before every plan (validation and labelled self-review)

Claude Opus 5.5 / Claude Code desktop / HIGH; effort xhigh; 2026-10-06.
Task `task-20261005T194500Z-book-plan-appendix-sql`. No database or benchmark was started.

## Owner feedback

"show the query that was used before the text of the explain analyze. the current text is not easy
to identify the query behind the explain analyze."

## What changed in the book

- **Each of the 45 plans in *Query plans behind the SELECT claims* (pp. 75–92) is preceded by:**
  - its statement name and bound parameters;
  - the bound SQL, in a box with a left rule and syntax highlighting;
  - then a "Plan" caption with the cell axis.
- **Statements in node-tree checks also show their SQL.** This covers the five join spellings
  checked against `JOIN ... WHERE`, the default CTE policy, IN and the broad expression shape. A
  compared plan the appendix does not quote gets its SQL under "Compared plan, not quoted here". SQL
  already printed is referenced ("Same SQL as in the … exhibit above") instead of repeated.
- **The intro explains** where the SQL comes from and the constant switches such as
  `WHERE 1 = 0 OR c.id <= 1`.
- **The join section explains** that `customers` and `facts` are views over person, charity and
  donation, which is why the plans name those tables.

## How the SQL is made trustworthy

- **Rebuilt with the code that ran.** Result cells store the statement name and the bound
  parameters (`k=v,...`) but not the text. `studies/01-charity-tree/selectbench/cmd/appendixsql`
  (new, DB-free) uses the harness's own `selectbench.Load`/`LoadJoin` and `selectbench.Bind`.
  - `LoadJoin` also generates the CTE-reuse statements, which exist only in Go.
  - The exporter inverts the `k=v` list and proves the inversion by rebuilding the string.
  - It writes 62 statements to `book/appendix/query-plans.sql.json`, listing per catalogue the
    source files that determine the text.
- **`tools/evidence plans` refuses the SQL when:**
  - its statement or parameters differ from the cell;
  - no entry exists for a quoted cell, or an entry is used by no exhibit;
  - any listed catalogue source is not byte-identical (`git show <commit>:<path>`, line endings
    normalized) at the exhibit's producing commit.
  - All five join and five teaching sources matched at `7243fee0`, `ba821e7b` and `ca2a38fb`.
- **New tests:**
  - a stale-parameter control;
  - a missing-SQL control;
  - a source absent at the producing commit (`joincatalog.go` at `ca2a38fb`), refused;
  - an unchanged source accepted.

## Checks

- `go vet ./...`, `gofmt -l .`, `go test -count=1 ./...` (tools/evidence): ok.
  `go vet ./selectbench/cmd/...` (Study 01): ok.
- `go run . plans --repo ../..`: current (25 exhibits, 45 plans).
- `go run . validate --repo ../..`: 50 claims, 0 errors.
- `book/evidence/check.sh`: OK.
- Canonical `book/build.sh` in WSL from clean commit `c54b8f62`, verified ([log](build.log)):
  - 117 pages, 4/4 fonts, 50/50 claims, edition "Edition 3 preview 7";
  - PDF sha256 `4ddc80d85c53e47fb59daf2d2aa435c0ccf6992ed00927854dfe31c3d30a4ef9`.
- Pages 75–92 checked: 45 plan SQL captions. [`proof/`](proof/) keeps pages 75, 76, 78, 83 and 90.
  - Fixed during inspection: a SQL box split across a page break (now unbreakable).
  - Fixed during inspection: two queries ended a page with their plan on the next page (the
    plan's SQL box is now sticky). A build of `254f3a01` showed both defects.

## Self-review

Same session as the author; labelled.

- The SQL is the harness's binding, not authored text.
- The intro's claim that the planner folds the constant switches is visible in the quoted plans,
  e.g. `Index Cond: (id <= 1)` under `WHERE 1 = 0 OR c.id <= 1`.
- The existing reading notes stay accurate against the now-visible SQL; for example, `v6-sel-09`'s
  "two scalar subqueries" are the SUM and COUNT subselects.

Approve.
