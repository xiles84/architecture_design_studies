# Book: show the executed SQL before each EXPLAIN ANALYZE plan

**Task:** `task-20261005T194500Z-book-plan-appendix-sql` · **capability:** HIGH · **benchmark:** none

## Owner feedback (2026-10-05)

"show the query that was used before the text of the explain analyze. the current text is not easy to identify the query behind the explain analyze."

## Decisions

- Result cells store the statement name and the bound parameters (`k=v,...`) but not the SQL text. The SQL is rebuilt with the harness's own code: `selectbench.Load`/`LoadJoin` (including the generated CTE-reuse statements) and `selectbench.Bind`, in a DB-free exporter under `studies/01-charity-tree/selectbench/cmd/`. Nothing is retyped.
- The catalogue sources (`selectbench/sql/**`, `catalog.go`, `bind.go`, `joincatalog.go`, `joinfixture.go`) are byte-identical at HEAD and at the three producing commits (`ca2a38fb`, `7243fee0`, `ba821e7b`); `tools/evidence plans` re-checks this with git for every exhibit and refuses on a difference.
- The SQL is shown before each plan, and for every statement in a node-tree check. The domain join SQL reads the `customers`/`facts` views over person, charity and donation; the appendix says so.

NEXT MODEL: HIGH
