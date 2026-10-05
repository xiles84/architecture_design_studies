# Book: say plainly when a derived table does (not) outperform other SELECT forms

**Task:** `task-20261005T160000Z-book-derived-table-clarity` · **capability:** HIGH · **benchmark:** none

## Owner feedback (2026-10-05)

"The book does not make clear to me when (or if there is none) situation when filtered derived table queries outperform other forms."

## What the evidence already says

- The filtered derived table ran the same plan as `JOIN ... WHERE` in every measured pair (0.01%-90% selectivity, PostgreSQL and YugabyteDB, isolated and blended): `v6-sel-01`. It never outperformed.
- A derived table changes the plan only when it holds work the planner cannot merge. The measured case is pre-aggregation (`a02_aggregate_preagg_sel`) in the synthesis evidence tables: slower than join-then-GROUP-BY for 1-1,000 parents on PostgreSQL, below 1 only at 10,000 parents inside the 29% drift, and never faster on YugabyteDB. It is not yet a registered claim.

## Work

Register the pre-aggregation result from `studies/01-charity-tree/reports/20261005T-select-synthesis-evidence.md` (no new measurement), state both facts in the Charity case and guide, add the number-free explanation to theory, rebuild and inspect.

NEXT MODEL: HIGH
