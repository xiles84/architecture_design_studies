# Book handoff — SELECT case results (from the HIGH synthesis)

For `task-20261004T142623Z-select-book-high` (HIGH). Prepared by the synthesis task
`task-20261002T131729Z-select-synthesis`. Claude Opus 5.5 / Claude Code desktop, HIGH,
effort xhigh, 2026-10-05.

**Inputs:**

- signed synthesis analysis
  `studies/01-charity-tree/reports/analyses/20261005T-select-synthesis--claude-opus-5-5--2026-10-05.md`;
- generated evidence tables `studies/01-charity-tree/reports/20261005T-select-synthesis-evidence.md`;
- registered claims `v6-sel-01` … `v6-sel-17` and `v6-sel-gap-01` in
  `book/evidence/v6/claims.json`.

Every number below is copied from those claims. Print nothing that is not in a claim.

## Decisions already made (do not reopen)

- **Teach answer contracts before speed.**
  - Optional match: `WHERE` on the child is a different answer.
  - Membership: a plain join multiplies parents.
  - Absence: nullable `NOT IN` returns nothing.
  - Top-N: `ROW_NUMBER` needs a total-order tiebreak.

  These are control results, never timing wins.
- **No universal ranking.** Every winner is printed with its engine, its regime (selectivity,
  references, fanout, parents, depth) and the claim id.
- **YugabyteDB access/existence numbers are scan-only.** Those runs used hash-sharded default
  indexes, and no plan used an index (`v6-sel-gap-01`). Print them only with that label, or
  leave them out. The join family's YugabyteDB numbers use a range index and need no such label.
- **Placement is a separate decision.** Print `v6-sel-16`/`v6-sel-17` in the placement or case
  section with "containers on one laptop" and "verified tablets and leaders". Never present them as
  spelling advice.
- **Cite the qualification.** Where the v4 join analysis is cited, cite the review-4
  qualification `20261005T-select-join-review-4--claude-opus-5-5` with it.
- **Teach `NOT IN` with three-valued logic:** `FALSE` or `UNKNOWN`, never `TRUE` once the
  subquery holds a NULL. GPT-6's correction applies.
- **Expectations and observations stay apart.** Theory keeps the expectations; cases keep the
  observations. The customers/orders fixture is a teaching analogue with no timing. Map it to
  person/donation (the domain views `customers`/`facts` over `person`/`donation`).

## Case variations to add (Charity donations case)

| Variation | Show | Fits | Dies / avoid | Claims |
|---|---|---|---|---|
| Six join/filter forms | Same plan in all 120 pairs, 0.01–90% selectivity, both engines | Write the form that reads best | No form is faster; do not rank them | `v6-sel-01` |
| CTE fence on a selective one-use filter | Forced `MATERIALIZED` with the filter outside: 84–101x (PG), 34–49x (YB) at 0.01%, ~1x at 90%. Filter inside: ≤1.61x (PG) | Default policy, or put the filter inside the CTE | Forcing `MATERIALIZED` with an outside selective filter | `v6-sel-02` |
| Reused expensive CTE | `NOT MATERIALIZED`/`MATERIALIZED` crosses at 1–2 references (PG 100%, YB both). PG 1% consumer never crosses within 1–16 | Engine default (inline 1, materialize ≥2) | Forcing `NOT MATERIALIZED` at ≥2 references, except on PG with a selective consumer (1.2–4.8x faster) | `v6-sel-03`, `v6-sel-04` |
| Optional match | `LEFT JOIN … ON` keeps every donor; `WHERE` drops unmatched ones (control fired 12/12) | Child predicate in `ON` | Child predicate in `WHERE` | `v6-sel-01` (controls) |

## Other SELECT family cases

| Family | Show | Fits | Dies / avoid | Claims |
|---|---|---|---|---|
| Membership | PG `JOIN + DISTINCT` 1.40–1.69x (fanout 1), 3.64–4.23x (fanout 100) of `EXISTS`; YB ≤1.16x; `IN` ≈ `EXISTS` | `EXISTS` or `IN` | `JOIN + DISTINCT` (PG); `DISTINCT` on a non-key | `v6-sel-05`, `v6-sel-06` |
| Absence | Nullable `NOT IN` returned 0 rows in 18/18 checks; PG `NOT EXISTS` 2.7–3.3x faster at fanout 100 | `NOT EXISTS` | Nullable `NOT IN`, always | `v6-sel-07`, `v6-sel-08` |
| Parent aggregates | PG correlated 0.47–0.90x at 1 parent, 2.49–2.70x at 10,000. YB correlated never wins (398–541x, 6.7–20 s at 10,000) | Join then `GROUP BY`; correlated only for a handful of parents on PG | Correlated scalar on YB beyond one parent; summing across two child joins (double count, 6/6) | `v6-sel-09`, `v6-sel-10` |
| Top 1 per parent | `LATERAL` wins only at 1 parent (YB 0.54–0.70x) and loses to 257–289x at 10,000; PG `DISTINCT ON` 0.61–0.64x at 10,000 | `ROW_NUMBER` with a tiebreak; `DISTINCT ON` (PG, N=1) at scale | `LATERAL` beyond a few parents; `ROW_NUMBER` without a tiebreak | `v6-sel-11`, `v6-sel-12` |
| Disjunctions | `UNION ALL`/`UNION` never beat `OR` (PG broad 1.23–1.55x slower) | `OR` | `UNION ALL` over overlapping arms (duplicates; control 36/36) | `v6-sel-13` |
| Expression predicate | Matching expression index: PG selective 0.28–0.33x | Index the expression you filter on | An expression predicate with no matching index | `v6-sel-14` |
| Pagination | Keyset 1.06–1.29x at depth ≤100, 0.21–0.24x at 10,000 (PG); same page at every depth | `OFFSET` for shallow pages; keyset for deep pages | Keyset at shallow pages; claiming flat keyset cost without the composite index | `v6-sel-15` |
| Placement (YB 3-node) | Split 0.54–0.73x at 2,000–180,000 rows, point read 1.22x; colocated 1.21–1.54x slower at 180,000 rows | Split for broad reads; colocate for point and medium reads | Reading container locality as placement | `v6-sel-16`, `v6-sel-17` |

## Compact selection guide (for the guide page)

Use the synthesis analysis's "Selection matrix" table as the source. One row per required answer.
Keep its "Evidence" column as claim ids.

## Code to show (reviewed; copy from the catalogues, do not retype)

- `studies/01-charity-tree/selectbench/sql/join/select_join_catalog.sql`:
  - `j01_join_where`;
  - `j05_cte_joined_outer` and `j05m_cte_joined_outer_mat` (the fence);
  - `j04_cte_filtered` (filter inside);
  - `j11_optional_on` against `x01_optional_where`.
- `studies/01-charity-tree/selectbench/joincatalog.go` `buildReuseSQL`: the reuse ladder shape, a
  CTE consumed n times with `UNION ALL`.
- `studies/01-charity-tree/selectbench/sql/select_catalog.sql`:
  - `s06_exists_parent`, `s08_join_distinct_parent`;
  - `s09_not_exists`, `x02_not_in_nullable`;
  - `s16_page_offset`, `s17_page_keyset`;
  - `s18_or_predicate`, `s19_union_all_predicate`.
- `studies/01-charity-tree/selectbench/sql/select_aggregate_catalog.sql`:
  - `a01_aggregate_join_sel`, `a03_aggregate_correlated_sel`;
  - `a04_topn_rownumber_sel`, `a05_topn_lateral_sel`, `a06_topn_distinct_on_sel`;
  - `x05_aggregate_double_count`.

## Limits every page must carry

- One laptop; containers share eight cores; no real network. Latency of one sequential reader,
  five samples per cell, three runs per arm.
- YugabyteDB access/existence plans are scan-only (`v6-sel-gap-01`). The three non-join families
  have no cluster arm.
- Crossings are bracketed between two levels, not located.
- Position drift is large in the aggregate family (29% PG, 74% YB). The PG pre-aggregation
  advantage at 10,000 parents (0.77–0.80x) is inside it, so it is not a registered claim and
  must not be printed as a finding.

## Escalate only if

A claim needed for a planned variation does not resolve, or the build's evidence gate fails on a
`v6-sel-*` claim. Hand a targeted correction to LOW; never re-measure inside the book task.

NEXT MODEL: HIGH
