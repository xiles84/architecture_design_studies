# SELECT comparison (generated measurements)

_This file contains measurements only. Interpretation belongs in a signed analysis._

**TL;DR:** a validation run supplies correctness and plan evidence. Within-run samples do not establish a replicated performance winner. Missing measurements are N/A.

## Provenance

| Field | Value |
|---|---|
| repo commit | `00e94ebcf3a87bc153045a216131dfa3328a01f6` |
| inputs digest (SHA-256 of result JSON bytes) | `9e7334f6c06ab45044a716093340bbdd81a4839550a592f99a67a41cdf06e0ae` |
| benchmark image | `localhost/charitybench-select:select-mergecheck-yb-20261004` |
| image digest | `70565d24e54146717cc8e6622fc046e44e83c9966d9a1bd0d77772594ce4fda4` |
| run tag | `run/study-01/select-mergecheck-yb-20261004` |
| environment | `xiles-podman-select-aggregate-20261003` |
| engine image digest | `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35` |
| engine version | PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit |
| dataset seed | 42 |
| dataset fingerprint | `f05b8420d983e1c8` |
| setting: client_workers | 1 sequential reader; 0 writers |
| setting: connection_pool | 2 = 1 reader/setup + 1 reserve |
| setting: distributed_rpc | not parsed; raw EXPLAIN retained (coverage gap) |
| setting: effective_cache_size | 524288 |
| setting: index_set | orders_amount_idx,orders_customer_idx,primary keys |
| setting: max_connections | 300 |
| setting: prepared_plan_policy | pgx simple protocol; no prepared statements |
| setting: random_page_cost | 4 |
| setting: scale | tiny |
| setting: shared_buffers | 16384 |
| setting: statement_budget_ms | 120000 |
| setting: statement_timeout | 120000 |
| setting: statistics | ANALYZE after deterministic load |
| setting: topology | single |
| setting: track_io_timing | off |
| setting: whole_run_count | 1; no replicated winner |
| setting: work_mem | 4096 |

## Cells

| run | topology | statement | group | status | median ms | spread % | trials | plan structure |
|---|---|---|---|---|---|---|---|---|
| select-mergecheck-yb-20261004 | single | s09_not_exists | sel_absence | valid | 0.655311 | 0.0 | 1 | `4befbd47cba4cb7d` |
| select-mergecheck-yb-20261004 | single | s10_left_anti | sel_absence | valid | 0.596303 | 0.0 | 1 | `c46cbe8221360b15` |
| select-mergecheck-yb-20261004 | single | s12_aggregate_join | sel_aggregate | valid | 0.731475 | 0.0 | 1 | `e4988e033c874e00` |
| select-mergecheck-yb-20261004 | single | s13_aggregate_preagg | sel_aggregate | valid | 0.599603 | 0.0 | 1 | `1b1b6b64788bee13` |
| select-mergecheck-yb-20261004 | single | a01_aggregate_join_sel | sel_aggregate_sel | valid | 0.664743 | 0.0 | 1 | `e4988e033c874e00` |
| select-mergecheck-yb-20261004 | single | a02_aggregate_preagg_sel | sel_aggregate_sel | valid | 0.628899 | 0.0 | 1 | `1b1b6b64788bee13` |
| select-mergecheck-yb-20261004 | single | a03_aggregate_correlated_sel | sel_aggregate_sel | valid | 1.433938 | 0.0 | 1 | `a246247eaf8a7c1a` |
| select-mergecheck-yb-20261004 | single | s01_join_where | sel_full_rows | valid | 0.767637 | 0.0 | 1 | `c44ec06f089f0d1c` |
| select-mergecheck-yb-20261004 | single | s02_filtered_derived | sel_full_rows | valid | 0.567834 | 0.0 | 1 | `c44ec06f089f0d1c` |
| select-mergecheck-yb-20261004 | single | s03_cte_filtered | sel_full_rows | valid | 0.650351 | 0.0 | 1 | `c44ec06f089f0d1c` |
| select-mergecheck-yb-20261004 | single | s04_join_on | sel_full_rows | valid | 0.655241 | 0.0 | 1 | `c44ec06f089f0d1c` |
| select-mergecheck-yb-20261004 | single | s05_comma_join | sel_full_rows | valid | 0.610298 | 0.0 | 1 | `c44ec06f089f0d1c` |
| select-mergecheck-yb-20261004 | single | s06_exists_parent | sel_membership | valid | 0.628827 | 0.0 | 1 | `7829bec10abd7f80` |
| select-mergecheck-yb-20261004 | single | s07_in_parent | sel_membership | valid | 0.577840 | 0.0 | 1 | `7829bec10abd7f80` |
| select-mergecheck-yb-20261004 | single | s08_join_distinct_parent | sel_membership | valid | 0.961873 | 0.0 | 1 | `5a4de51c305ef1b0` |
| select-mergecheck-yb-20261004 | single | s11_optional_on | sel_optional_match | valid | 0.655888 | 0.0 | 1 | `c46cbe8221360b15` |
| select-mergecheck-yb-20261004 | single | s16_page_offset | sel_pagination | valid | 0.333675 | 0.0 | 1 | `ac4509d295bf3348` |
| select-mergecheck-yb-20261004 | single | s17_page_keyset | sel_pagination | valid | 0.369154 | 0.0 | 1 | `ac4509d295bf3348` |
| select-mergecheck-yb-20261004 | single | s18_or_predicate | sel_predicate_alt | valid | 0.377088 | 0.0 | 1 | `5191eeb70daef797` |
| select-mergecheck-yb-20261004 | single | s19_union_all_predicate | sel_predicate_alt | valid | 0.623446 | 0.0 | 1 | `d95e057d512f5ae5` |
| select-mergecheck-yb-20261004 | single | s20_union_predicate | sel_predicate_alt | valid | 0.555829 | 0.0 | 1 | `12de8f648f9520a0` |
| select-mergecheck-yb-20261004 | single | s14_topn_rownumber | sel_topn | valid | 0.549941 | 0.0 | 1 | `e2ad05f4f9b57fe5` |
| select-mergecheck-yb-20261004 | single | s15_topn_lateral | sel_topn | valid | 1.014916 | 0.0 | 1 | `0504ce42ca4d7a56` |
| select-mergecheck-yb-20261004 | single | a06_topn_distinct_on_sel | sel_topn_one | valid | 0.355507 | 0.0 | 1 | `520ea58629cf6ed2` |
| select-mergecheck-yb-20261004 | single | a04_topn_rownumber_sel | sel_topn_sel | valid | 0.447417 | 0.0 | 1 | `e2ad05f4f9b57fe5` |
| select-mergecheck-yb-20261004 | single | a05_topn_lateral_sel | sel_topn_sel | valid | 1.025022 | 0.0 | 1 | `0504ce42ca4d7a56` |

## Excluded cells

None: every cell is valid.

## Comparisons

### Same plan (no plan difference to explain)

- `s01_join_where` vs `s02_filtered_derived` on single: same plan structure
- `s01_join_where` vs `s03_cte_filtered` on single: same plan structure
- `s01_join_where` vs `s04_join_on` on single: same plan structure
- `s01_join_where` vs `s05_comma_join` on single: same plan structure
- `s02_filtered_derived` vs `s03_cte_filtered` on single: same plan structure
- `s02_filtered_derived` vs `s04_join_on` on single: same plan structure
- `s02_filtered_derived` vs `s05_comma_join` on single: same plan structure
- `s03_cte_filtered` vs `s04_join_on` on single: same plan structure
- `s03_cte_filtered` vs `s05_comma_join` on single: same plan structure
- `s04_join_on` vs `s05_comma_join` on single: same plan structure
- `s06_exists_parent` vs `s07_in_parent` on single: same plan structure
- `s16_page_offset` vs `s17_page_keyset` on single: same plan structure

### No resolved difference

None.

### Within-run differences (not replicated)

- `s09_not_exists` vs `s10_left_anti` on single: within-run medians 0.655311 / 0.596303 ms; whole-run replication required — no resolved difference
- `s12_aggregate_join` vs `s13_aggregate_preagg` on single: within-run medians 0.731475 / 0.599603 ms; whole-run replication required — no resolved difference
- `a01_aggregate_join_sel` vs `a02_aggregate_preagg_sel` on single: within-run medians 0.664743 / 0.628899 ms; whole-run replication required — no resolved difference
- `a01_aggregate_join_sel` vs `a03_aggregate_correlated_sel` on single: within-run medians 0.664743 / 1.433938 ms; whole-run replication required — no resolved difference
- `a02_aggregate_preagg_sel` vs `a03_aggregate_correlated_sel` on single: within-run medians 0.628899 / 1.433938 ms; whole-run replication required — no resolved difference
- `s06_exists_parent` vs `s08_join_distinct_parent` on single: within-run medians 0.628827 / 0.961873 ms; whole-run replication required — no resolved difference
- `s07_in_parent` vs `s08_join_distinct_parent` on single: within-run medians 0.577840 / 0.961873 ms; whole-run replication required — no resolved difference
- `s18_or_predicate` vs `s19_union_all_predicate` on single: within-run medians 0.377088 / 0.623446 ms; whole-run replication required — no resolved difference
- `s18_or_predicate` vs `s20_union_predicate` on single: within-run medians 0.377088 / 0.555829 ms; whole-run replication required — no resolved difference
- `s19_union_all_predicate` vs `s20_union_predicate` on single: within-run medians 0.623446 / 0.555829 ms; whole-run replication required — no resolved difference
- `s14_topn_rownumber` vs `s15_topn_lateral` on single: within-run medians 0.549941 / 1.014916 ms; whole-run replication required — no resolved difference
- `a04_topn_rownumber_sel` vs `a05_topn_lateral_sel` on single: within-run medians 0.447417 / 1.025022 ms; whole-run replication required — no resolved difference

## Live negative controls

- `x05_aggregate_double_count`: detected=true, rows=3
- `x06_topn_rank_ties`: detected=true, rows=4
- `x01_optional_where`: detected=true, rows=3
- `x02_not_in_nullable`: detected=true, rows=0
- `x03_distinct_on_name`: detected=true, rows=1
- `x04_limit_before_join`: detected=true, rows=2
- `x05_union_all_overlap`: detected=true, rows=4

## Captured plans and result sizes

### s09_not_exists

Rows: 0; serialized decoded-row JSON bytes: 4 (not wire bytes). Planning/execution: 0.077000 / 0.476000 ms.

```text
Hash Anti Join  (cost=42.49..42.67 rows=1 width=7) (actual time=0.445..0.445 rows=0 loops=1)
  Hash Cond: (c.id = o.customer_id)
  ->  Seq Scan on customers c  (cost=21.05..21.21 rows=3 width=7) (actual time=0.215..0.216 rows=3 loops=1)
  ->  Hash  (cost=21.38..21.38 rows=5 width=4) (actual time=0.219..0.219 rows=4 loops=1)
        Buckets: 1024  Batches: 1  Memory Usage: 9kB
        ->  Seq Scan on orders o  (cost=21.05..21.38 rows=5 width=4) (actual time=0.210..0.211 rows=5 loops=1)
              Storage Filter: (amount > 100)
Planning Time: 0.077 ms
Execution Time: 0.476 ms
Peak Memory Usage: 205 kB

```

### s10_left_anti

Rows: 0; serialized decoded-row JSON bytes: 4 (not wire bytes). Planning/execution: 0.080000 / 0.440000 ms.

```text
Hash Right Join  (cost=42.30..42.70 rows=1 width=7) (actual time=0.389..0.390 rows=0 loops=1)
  Hash Cond: (o.customer_id = c.id)
  Filter: (o.id IS NULL)
  Rows Removed by Filter: 4
  ->  Seq Scan on orders o  (cost=21.05..21.40 rows=5 width=8) (actual time=0.170..0.172 rows=5 loops=1)
        Storage Filter: (amount > 100)
  ->  Hash  (cost=21.21..21.21 rows=3 width=7) (actual time=0.209..0.209 rows=3 loops=1)
        Buckets: 1024  Batches: 1  Memory Usage: 9kB
        ->  Seq Scan on customers c  (cost=21.05..21.21 rows=3 width=7) (actual time=0.200..0.202 rows=3 loops=1)
Planning Time: 0.080 ms
Execution Time: 0.440 ms
Peak Memory Usage: 191 kB

```

### s12_aggregate_join

Rows: 3; serialized decoded-row JSON bytes: 46 (not wire bytes). Planning/execution: 0.083000 / 0.439000 ms.

```text
GroupAggregate  (cost=42.72..42.89 rows=3 width=23) (actual time=0.405..0.409 rows=3 loops=1)
  Group Key: c.id
  ->  Merge Left Join  (cost=42.72..42.82 rows=6 width=15) (actual time=0.400..0.403 rows=5 loops=1)
        Merge Cond: (c.id = o.customer_id)
        ->  Sort  (cost=21.23..21.24 rows=3 width=7) (actual time=0.205..0.206 rows=3 loops=1)
              Sort Key: c.id
              Sort Method: quicksort  Memory: 25kB
              ->  Seq Scan on customers c  (cost=21.05..21.21 rows=3 width=7) (actual time=0.198..0.200 rows=3 loops=1)
        ->  Sort  (cost=21.49..21.51 rows=6 width=12) (actual time=0.192..0.192 rows=6 loops=1)
              Sort Key: o.customer_id
              Sort Method: quicksort  Memory: 25kB
              ->  Seq Scan on orders o  (cost=21.05..21.41 rows=6 width=12) (actual time=0.178..0.181 rows=6 loops=1)
Planning Time: 0.083 ms
Execution Time: 0.439 ms
Peak Memory Usage: 206 kB

```

### s13_aggregate_preagg

Rows: 3; serialized decoded-row JSON bytes: 46 (not wire bytes). Planning/execution: 0.112000 / 0.436000 ms.

```text
Hash Left Join  (cost=42.67..42.85 rows=3 width=23) (actual time=0.402..0.404 rows=3 loops=1)
  Hash Cond: (c.id = a.customer_id)
  ->  Seq Scan on customers c  (cost=21.05..21.21 rows=3 width=7) (actual time=0.142..0.143 rows=3 loops=1)
  ->  Hash  (cost=21.59..21.59 rows=3 width=20) (actual time=0.253..0.254 rows=3 loops=1)
        Buckets: 1024  Batches: 1  Memory Usage: 9kB
        ->  Subquery Scan on a  (cost=21.47..21.59 rows=3 width=20) (actual time=0.245..0.249 rows=4 loops=1)
              ->  GroupAggregate  (cost=21.47..21.56 rows=3 width=20) (actual time=0.245..0.247 rows=4 loops=1)
                    Group Key: orders.customer_id
                    ->  Sort  (cost=21.47..21.48 rows=6 width=8) (actual time=0.242..0.243 rows=6 loops=1)
                          Sort Key: orders.customer_id
                          Sort Method: quicksort  Memory: 25kB
                          ->  Seq Scan on orders  (cost=21.05..21.39 rows=6 width=8) (actual time=0.234..0.235 rows=6 loops=1)
Planning Time: 0.112 ms
Execution Time: 0.436 ms
Peak Memory Usage: 232 kB

```

### a01_aggregate_join_sel

Rows: 3; serialized decoded-row JSON bytes: 46 (not wire bytes). Planning/execution: 0.101000 / 0.679000 ms.

```text
GroupAggregate  (cost=42.72..42.89 rows=3 width=23) (actual time=0.436..0.440 rows=3 loops=1)
  Group Key: c.id
  ->  Merge Left Join  (cost=42.72..42.82 rows=6 width=15) (actual time=0.430..0.433 rows=5 loops=1)
        Merge Cond: (c.id = o.customer_id)
        ->  Sort  (cost=21.23..21.24 rows=3 width=7) (actual time=0.238..0.239 rows=3 loops=1)
              Sort Key: c.id
              Sort Method: quicksort  Memory: 25kB
              ->  Seq Scan on customers c  (cost=21.05..21.21 rows=3 width=7) (actual time=0.214..0.229 rows=3 loops=1)
        ->  Sort  (cost=21.49..21.51 rows=6 width=12) (actual time=0.189..0.189 rows=6 loops=1)
              Sort Key: o.customer_id
              Sort Method: quicksort  Memory: 25kB
              ->  Seq Scan on orders o  (cost=21.05..21.41 rows=6 width=12) (actual time=0.167..0.169 rows=6 loops=1)
Planning Time: 0.101 ms
Execution Time: 0.679 ms
Peak Memory Usage: 206 kB

```

### a02_aggregate_preagg_sel

Rows: 3; serialized decoded-row JSON bytes: 46 (not wire bytes). Planning/execution: 0.071000 / 0.424000 ms.

```text
Hash Left Join  (cost=42.67..42.85 rows=3 width=23) (actual time=0.387..0.390 rows=3 loops=1)
  Hash Cond: (c.id = a.customer_id)
  ->  Seq Scan on customers c  (cost=21.05..21.21 rows=3 width=7) (actual time=0.212..0.213 rows=3 loops=1)
  ->  Hash  (cost=21.59..21.59 rows=3 width=20) (actual time=0.166..0.166 rows=3 loops=1)
        Buckets: 1024  Batches: 1  Memory Usage: 9kB
        ->  Subquery Scan on a  (cost=21.47..21.59 rows=3 width=20) (actual time=0.159..0.162 rows=4 loops=1)
              ->  GroupAggregate  (cost=21.47..21.56 rows=3 width=20) (actual time=0.158..0.161 rows=4 loops=1)
                    Group Key: orders.customer_id
                    ->  Sort  (cost=21.47..21.48 rows=6 width=8) (actual time=0.155..0.156 rows=6 loops=1)
                          Sort Key: orders.customer_id
                          Sort Method: quicksort  Memory: 25kB
                          ->  Seq Scan on orders  (cost=21.05..21.39 rows=6 width=8) (actual time=0.147..0.148 rows=6 loops=1)
Planning Time: 0.071 ms
Execution Time: 0.424 ms
Peak Memory Usage: 232 kB

```

### a03_aggregate_correlated_sel

Rows: 3; serialized decoded-row JSON bytes: 46 (not wire bytes). Planning/execution: 0.065000 / 1.214000 ms.

```text
Seq Scan on customers c  (cost=21.05..149.52 rows=3 width=23) (actual time=0.556..1.184 rows=3 loops=1)
  SubPlan 1
    ->  Finalize Aggregate  (cost=21.37..21.38 rows=1 width=8) (actual time=0.162..0.162 rows=1 loops=3)
          ->  Seq Scan on orders o  (cost=21.05..21.37 rows=2 width=4) (actual time=0.158..0.159 rows=1 loops=3)
                Storage Filter: (customer_id = c.id)
                Partial Aggregate: true
  SubPlan 2
    ->  Finalize Aggregate  (cost=21.38..21.39 rows=1 width=8) (actual time=0.162..0.162 rows=1 loops=3)
          ->  Seq Scan on orders o_1  (cost=21.05..21.37 rows=2 width=0) (actual time=0.159..0.159 rows=1 loops=3)
                Storage Filter: (customer_id = c.id)
                Partial Aggregate: true
Planning Time: 0.065 ms
Execution Time: 1.214 ms
Peak Memory Usage: 128 kB

```

### s01_join_where

Rows: 4; serialized decoded-row JSON bytes: 73 (not wire bytes). Planning/execution: 0.089000 / 0.376000 ms.

```text
Hash Join  (cost=42.30..42.72 rows=5 width=15) (actual time=0.346..0.349 rows=4 loops=1)
  Hash Cond: (o.customer_id = c.id)
  ->  Seq Scan on orders o  (cost=21.05..21.42 rows=5 width=12) (actual time=0.182..0.183 rows=5 loops=1)
        Storage Filter: (amount > 100)
  ->  Hash  (cost=21.21..21.21 rows=3 width=7) (actual time=0.156..0.156 rows=3 loops=1)
        Buckets: 1024  Batches: 1  Memory Usage: 9kB
        ->  Seq Scan on customers c  (cost=21.05..21.21 rows=3 width=7) (actual time=0.149..0.151 rows=3 loops=1)
Planning Time: 0.089 ms
Execution Time: 0.376 ms
Peak Memory Usage: 205 kB

```

### s02_filtered_derived

Rows: 4; serialized decoded-row JSON bytes: 73 (not wire bytes). Planning/execution: 0.087000 / 0.364000 ms.

```text
Hash Join  (cost=42.30..42.72 rows=5 width=15) (actual time=0.338..0.341 rows=4 loops=1)
  Hash Cond: (orders.customer_id = c.id)
  ->  Seq Scan on orders  (cost=21.05..21.42 rows=5 width=12) (actual time=0.168..0.169 rows=5 loops=1)
        Storage Filter: (amount > 100)
  ->  Hash  (cost=21.21..21.21 rows=3 width=7) (actual time=0.164..0.164 rows=3 loops=1)
        Buckets: 1024  Batches: 1  Memory Usage: 9kB
        ->  Seq Scan on customers c  (cost=21.05..21.21 rows=3 width=7) (actual time=0.157..0.158 rows=3 loops=1)
Planning Time: 0.087 ms
Execution Time: 0.364 ms
Peak Memory Usage: 205 kB

```

### s03_cte_filtered

Rows: 4; serialized decoded-row JSON bytes: 73 (not wire bytes). Planning/execution: 0.096000 / 0.412000 ms.

```text
Hash Join  (cost=42.30..42.72 rows=5 width=15) (actual time=0.380..0.383 rows=4 loops=1)
  Hash Cond: (orders.customer_id = c.id)
  ->  Seq Scan on orders  (cost=21.05..21.42 rows=5 width=12) (actual time=0.199..0.200 rows=5 loops=1)
        Storage Filter: (amount > 100)
  ->  Hash  (cost=21.21..21.21 rows=3 width=7) (actual time=0.172..0.173 rows=3 loops=1)
        Buckets: 1024  Batches: 1  Memory Usage: 9kB
        ->  Seq Scan on customers c  (cost=21.05..21.21 rows=3 width=7) (actual time=0.163..0.164 rows=3 loops=1)
Planning Time: 0.096 ms
Execution Time: 0.412 ms
Peak Memory Usage: 205 kB

```

### s04_join_on

Rows: 4; serialized decoded-row JSON bytes: 73 (not wire bytes). Planning/execution: 0.082000 / 0.402000 ms.

```text
Hash Join  (cost=42.30..42.72 rows=5 width=15) (actual time=0.371..0.375 rows=4 loops=1)
  Hash Cond: (o.customer_id = c.id)
  ->  Seq Scan on orders o  (cost=21.05..21.42 rows=5 width=12) (actual time=0.201..0.203 rows=5 loops=1)
        Storage Filter: (amount > 100)
  ->  Hash  (cost=21.21..21.21 rows=3 width=7) (actual time=0.161..0.161 rows=3 loops=1)
        Buckets: 1024  Batches: 1  Memory Usage: 9kB
        ->  Seq Scan on customers c  (cost=21.05..21.21 rows=3 width=7) (actual time=0.151..0.152 rows=3 loops=1)
Planning Time: 0.082 ms
Execution Time: 0.402 ms
Peak Memory Usage: 205 kB

```

### s05_comma_join

Rows: 4; serialized decoded-row JSON bytes: 73 (not wire bytes). Planning/execution: 0.078000 / 0.390000 ms.

```text
Hash Join  (cost=42.30..42.72 rows=5 width=15) (actual time=0.346..0.349 rows=4 loops=1)
  Hash Cond: (o.customer_id = c.id)
  ->  Seq Scan on orders o  (cost=21.05..21.42 rows=5 width=12) (actual time=0.190..0.191 rows=5 loops=1)
        Storage Filter: (amount > 100)
  ->  Hash  (cost=21.21..21.21 rows=3 width=7) (actual time=0.148..0.148 rows=3 loops=1)
        Buckets: 1024  Batches: 1  Memory Usage: 9kB
        ->  Seq Scan on customers c  (cost=21.05..21.21 rows=3 width=7) (actual time=0.141..0.143 rows=3 loops=1)
Planning Time: 0.078 ms
Execution Time: 0.390 ms
Peak Memory Usage: 205 kB

```

### s06_exists_parent

Rows: 3; serialized decoded-row JSON bytes: 28 (not wire bytes). Planning/execution: 0.084000 / 0.428000 ms.

```text
Hash Semi Join  (cost=42.49..42.70 rows=3 width=7) (actual time=0.399..0.402 rows=3 loops=1)
  Hash Cond: (c.id = o.customer_id)
  ->  Seq Scan on customers c  (cost=21.05..21.21 rows=3 width=7) (actual time=0.194..0.195 rows=3 loops=1)
  ->  Hash  (cost=21.38..21.38 rows=5 width=4) (actual time=0.197..0.197 rows=4 loops=1)
        Buckets: 1024  Batches: 1  Memory Usage: 9kB
        ->  Seq Scan on orders o  (cost=21.05..21.38 rows=5 width=4) (actual time=0.186..0.188 rows=5 loops=1)
              Storage Filter: (amount > 100)
Planning Time: 0.084 ms
Execution Time: 0.428 ms
Peak Memory Usage: 205 kB

```

### s07_in_parent

Rows: 3; serialized decoded-row JSON bytes: 28 (not wire bytes). Planning/execution: 0.088000 / 0.376000 ms.

```text
Hash Semi Join  (cost=42.49..42.70 rows=3 width=7) (actual time=0.349..0.351 rows=3 loops=1)
  Hash Cond: (c.id = o.customer_id)
  ->  Seq Scan on customers c  (cost=21.05..21.21 rows=3 width=7) (actual time=0.170..0.170 rows=3 loops=1)
  ->  Hash  (cost=21.38..21.38 rows=5 width=4) (actual time=0.174..0.174 rows=4 loops=1)
        Buckets: 1024  Batches: 1  Memory Usage: 9kB
        ->  Seq Scan on orders o  (cost=21.05..21.38 rows=5 width=4) (actual time=0.167..0.168 rows=5 loops=1)
              Storage Filter: (amount > 100)
Planning Time: 0.088 ms
Execution Time: 0.376 ms
Peak Memory Usage: 205 kB

```

### s08_join_distinct_parent

Rows: 3; serialized decoded-row JSON bytes: 28 (not wire bytes). Planning/execution: 0.099000 / 0.437000 ms.

```text
Unique  (cost=42.81..42.84 rows=3 width=7) (actual time=0.402..0.404 rows=3 loops=1)
  ->  Sort  (cost=42.81..42.82 rows=5 width=7) (actual time=0.401..0.402 rows=4 loops=1)
        Sort Key: c.id, c.name
        Sort Method: quicksort  Memory: 25kB
        ->  Merge Join  (cost=42.67..42.75 rows=5 width=7) (actual time=0.392..0.395 rows=4 loops=1)
              Merge Cond: (c.id = o.customer_id)
              ->  Sort  (cost=21.23..21.24 rows=3 width=7) (actual time=0.206..0.206 rows=3 loops=1)
                    Sort Key: c.id
                    Sort Method: quicksort  Memory: 25kB
                    ->  Seq Scan on customers c  (cost=21.05..21.21 rows=3 width=7) (actual time=0.198..0.200 rows=3 loops=1)
              ->  Sort  (cost=21.44..21.45 rows=5 width=4) (actual time=0.184..0.184 rows=5 loops=1)
                    Sort Key: o.customer_id
                    Sort Method: quicksort  Memory: 25kB
                    ->  Seq Scan on orders o  (cost=21.05..21.38 rows=5 width=4) (actual time=0.173..0.175 rows=5 loops=1)
                          Storage Filter: (amount > 100)
Planning Time: 0.099 ms
Execution Time: 0.437 ms
Peak Memory Usage: 242 kB

```

### s11_optional_on

Rows: 4; serialized decoded-row JSON bytes: 73 (not wire bytes). Planning/execution: 0.064000 / 0.442000 ms.

```text
Hash Right Join  (cost=42.30..42.72 rows=5 width=15) (actual time=0.412..0.418 rows=4 loops=1)
  Hash Cond: (o.customer_id = c.id)
  ->  Seq Scan on orders o  (cost=21.05..21.42 rows=5 width=12) (actual time=0.222..0.224 rows=5 loops=1)
        Storage Filter: (amount > 100)
  ->  Hash  (cost=21.21..21.21 rows=3 width=7) (actual time=0.185..0.186 rows=3 loops=1)
        Buckets: 1024  Batches: 1  Memory Usage: 9kB
        ->  Seq Scan on customers c  (cost=21.05..21.21 rows=3 width=7) (actual time=0.179..0.180 rows=3 loops=1)
Planning Time: 0.064 ms
Execution Time: 0.442 ms
Peak Memory Usage: 191 kB

```

### s16_page_offset

Rows: 2; serialized decoded-row JSON bytes: 23 (not wire bytes). Planning/execution: 0.021000 / 0.233000 ms.

```text
Limit  (cost=21.45..21.46 rows=2 width=8) (actual time=0.212..0.213 rows=2 loops=1)
  ->  Sort  (cost=21.45..21.47 rows=6 width=8) (actual time=0.211..0.212 rows=2 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 25kB
        ->  Seq Scan on orders  (cost=21.05..21.39 rows=6 width=8) (actual time=0.204..0.206 rows=6 loops=1)
Planning Time: 0.021 ms
Execution Time: 0.233 ms
Peak Memory Usage: 92 kB

```

### s17_page_keyset

Rows: 2; serialized decoded-row JSON bytes: 23 (not wire bytes). Planning/execution: 0.028000 / 0.232000 ms.

```text
Limit  (cost=21.45..21.46 rows=2 width=8) (actual time=0.206..0.207 rows=2 loops=1)
  ->  Sort  (cost=21.45..21.47 rows=6 width=8) (actual time=0.205..0.205 rows=2 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 25kB
        ->  Seq Scan on orders  (cost=21.05..21.39 rows=6 width=8) (actual time=0.194..0.196 rows=6 loops=1)
Planning Time: 0.028 ms
Execution Time: 0.232 ms
Peak Memory Usage: 92 kB

```

### s18_or_predicate

Rows: 2; serialized decoded-row JSON bytes: 23 (not wire bytes). Planning/execution: 0.028000 / 0.322000 ms.

```text
Sort  (cost=21.42..21.43 rows=3 width=8) (actual time=0.301..0.301 rows=2 loops=1)
  Sort Key: id
  Sort Method: quicksort  Memory: 25kB
  ->  Seq Scan on orders  (cost=21.05..21.40 rows=3 width=8) (actual time=0.291..0.293 rows=2 loops=1)
        Storage Filter: ((amount <= 100) OR (amount >= 400))
Planning Time: 0.028 ms
Execution Time: 0.322 ms
Peak Memory Usage: 100 kB

```

### s19_union_all_predicate

Rows: 2; serialized decoded-row JSON bytes: 23 (not wire bytes). Planning/execution: 0.050000 / 0.412000 ms.

```text
Sort  (cost=42.85..42.86 rows=4 width=8) (actual time=0.385..0.387 rows=2 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 25kB
  ->  Append  (cost=21.05..42.81 rows=4 width=8) (actual time=0.377..0.380 rows=2 loops=1)
        ->  Seq Scan on orders  (cost=21.05..21.37 rows=1 width=8) (actual time=0.192..0.192 rows=0 loops=1)
              Storage Filter: (amount <= 100)
        ->  Seq Scan on orders orders_1  (cost=21.05..21.38 rows=3 width=8) (actual time=0.182..0.184 rows=2 loops=1)
              Storage Filter: (amount >= 400)
Planning Time: 0.050 ms
Execution Time: 0.412 ms
Peak Memory Usage: 164 kB

```

### s20_union_predicate

Rows: 2; serialized decoded-row JSON bytes: 23 (not wire bytes). Planning/execution: 0.040000 / 0.388000 ms.

```text
Sort  (cost=42.91..42.92 rows=4 width=8) (actual time=0.353..0.354 rows=2 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 25kB
  ->  HashAggregate  (cost=42.83..42.87 rows=4 width=8) (actual time=0.347..0.348 rows=2 loops=1)
        Group Key: orders.id, orders.amount
        Batches: 1  Memory Usage: 24kB
        ->  Append  (cost=21.05..42.81 rows=4 width=8) (actual time=0.341..0.344 rows=2 loops=1)
              ->  Seq Scan on orders  (cost=21.05..21.37 rows=1 width=8) (actual time=0.173..0.173 rows=0 loops=1)
                    Storage Filter: (amount <= 100)
              ->  Seq Scan on orders orders_1  (cost=21.05..21.38 rows=3 width=8) (actual time=0.164..0.166 rows=2 loops=1)
                    Storage Filter: (amount >= 400)
Planning Time: 0.040 ms
Execution Time: 0.388 ms
Peak Memory Usage: 164 kB

```

### s14_topn_rownumber

Rows: 5; serialized decoded-row JSON bytes: 67 (not wire bytes). Planning/execution: 0.045000 / 0.246000 ms.

```text
Subquery Scan on t  (cost=21.46..21.62 rows=5 width=12) (actual time=0.212..0.220 rows=5 loops=1)
  ->  WindowAgg  (cost=21.46..21.57 rows=5 width=52) (actual time=0.212..0.219 rows=5 loops=1)
        Run Condition: (row_number() OVER (?) <= 3)
        ->  Sort  (cost=21.46..21.47 rows=5 width=12) (actual time=0.198..0.198 rows=5 loops=1)
              Sort Key: o.customer_id, o.amount DESC NULLS LAST, o.id
              Sort Method: quicksort  Memory: 25kB
              ->  Seq Scan on orders o  (cost=21.05..21.40 rows=5 width=12) (actual time=0.190..0.191 rows=5 loops=1)
                    Storage Filter: (customer_id IS NOT NULL)
Planning Time: 0.045 ms
Execution Time: 0.246 ms
Peak Memory Usage: 198 kB

```

### s15_topn_lateral

Rows: 5; serialized decoded-row JSON bytes: 67 (not wire bytes). Planning/execution: 0.219000 / 1.083000 ms.

```text
Sort  (cost=85.55..85.56 rows=6 width=12) (actual time=1.040..1.040 rows=5 loops=1)
  Sort Key: c.id, o_1.amount DESC NULLS LAST, o_1.id
  Sort Method: quicksort  Memory: 25kB
  ->  Nested Loop  (cost=42.43..85.47 rows=6 width=12) (actual time=0.471..1.033 rows=5 loops=1)
        ->  Seq Scan on customers c  (cost=21.05..21.18 rows=3 width=4) (actual time=0.258..0.259 rows=3 loops=1)
        ->  Limit  (cost=21.39..21.39 rows=2 width=8) (actual time=0.255..0.256 rows=2 loops=3)
              ->  Sort  (cost=21.39..21.39 rows=2 width=8) (actual time=0.254..0.254 rows=2 loops=3)
                    Sort Key: o_1.amount DESC NULLS LAST, o_1.id
                    Sort Method: quicksort  Memory: 25kB
                    ->  Seq Scan on orders o_1  (cost=21.05..21.38 rows=2 width=8) (actual time=0.243..0.245 rows=2 loops=3)
                          Storage Filter: (customer_id = c.id)
Planning Time: 0.219 ms
Execution Time: 1.083 ms
Peak Memory Usage: 212 kB

```

### a06_topn_distinct_on_sel

Rows: 3; serialized decoded-row JSON bytes: 40 (not wire bytes). Planning/execution: 0.028000 / 0.209000 ms.

```text
Unique  (cost=21.46..21.49 rows=3 width=12) (actual time=0.187..0.189 rows=3 loops=1)
  ->  Sort  (cost=21.46..21.47 rows=5 width=12) (actual time=0.187..0.187 rows=5 loops=1)
        Sort Key: customer_id, amount DESC NULLS LAST, id
        Sort Method: quicksort  Memory: 25kB
        ->  Seq Scan on orders o  (cost=21.05..21.40 rows=5 width=12) (actual time=0.179..0.180 rows=5 loops=1)
              Storage Filter: (customer_id IS NOT NULL)
Planning Time: 0.028 ms
Execution Time: 0.209 ms
Peak Memory Usage: 100 kB

```

### a04_topn_rownumber_sel

Rows: 5; serialized decoded-row JSON bytes: 67 (not wire bytes). Planning/execution: 0.055000 / 0.302000 ms.

```text
Subquery Scan on t  (cost=21.46..21.62 rows=5 width=12) (actual time=0.253..0.271 rows=5 loops=1)
  ->  WindowAgg  (cost=21.46..21.57 rows=5 width=20) (actual time=0.252..0.268 rows=5 loops=1)
        Run Condition: (row_number() OVER (?) <= 3)
        ->  Sort  (cost=21.46..21.47 rows=5 width=12) (actual time=0.240..0.240 rows=5 loops=1)
              Sort Key: o.customer_id, o.amount DESC NULLS LAST, o.id
              Sort Method: quicksort  Memory: 25kB
              ->  Seq Scan on orders o  (cost=21.05..21.40 rows=5 width=12) (actual time=0.228..0.230 rows=5 loops=1)
                    Storage Filter: (customer_id IS NOT NULL)
Planning Time: 0.055 ms
Execution Time: 0.302 ms
Peak Memory Usage: 222 kB

```

### a05_topn_lateral_sel

Rows: 5; serialized decoded-row JSON bytes: 67 (not wire bytes). Planning/execution: 0.075000 / 0.833000 ms.

```text
Sort  (cost=85.55..85.56 rows=6 width=12) (actual time=0.804..0.805 rows=5 loops=1)
  Sort Key: c.id, o_1.amount DESC NULLS LAST, o_1.id
  Sort Method: quicksort  Memory: 25kB
  ->  Nested Loop  (cost=42.43..85.47 rows=6 width=12) (actual time=0.428..0.797 rows=5 loops=1)
        ->  Seq Scan on customers c  (cost=21.05..21.18 rows=3 width=4) (actual time=0.223..0.225 rows=3 loops=1)
        ->  Limit  (cost=21.39..21.39 rows=2 width=8) (actual time=0.187..0.188 rows=2 loops=3)
              ->  Sort  (cost=21.39..21.39 rows=2 width=8) (actual time=0.186..0.186 rows=2 loops=3)
                    Sort Key: o_1.amount DESC NULLS LAST, o_1.id
                    Sort Method: quicksort  Memory: 25kB
                    ->  Seq Scan on orders o_1  (cost=21.05..21.38 rows=2 width=8) (actual time=0.174..0.175 rows=2 loops=3)
                          Storage Filter: (customer_id = c.id)
Planning Time: 0.075 ms
Execution Time: 0.833 ms
Peak Memory Usage: 212 kB

```

### Not comparable (never ranked)

- - `s09_not_exists` vs `s12_aggregate_join`: not comparable — different semantic groups answer different questions
- - `s09_not_exists` vs `s13_aggregate_preagg`: not comparable — different semantic groups answer different questions
- - `s09_not_exists` vs `a01_aggregate_join_sel`: not comparable — different semantic groups answer different questions
- - `s09_not_exists` vs `a02_aggregate_preagg_sel`: not comparable — different semantic groups answer different questions
- - `s09_not_exists` vs `a03_aggregate_correlated_sel`: not comparable — different semantic groups answer different questions
- - `s09_not_exists` vs `s01_join_where`: not comparable — different semantic groups answer different questions
- - `s09_not_exists` vs `s02_filtered_derived`: not comparable — different semantic groups answer different questions
- - `s09_not_exists` vs `s03_cte_filtered`: not comparable — different semantic groups answer different questions
- - `s09_not_exists` vs `s04_join_on`: not comparable — different semantic groups answer different questions
- - `s09_not_exists` vs `s05_comma_join`: not comparable — different semantic groups answer different questions
- - `s09_not_exists` vs `s06_exists_parent`: not comparable — different semantic groups answer different questions
- - `s09_not_exists` vs `s07_in_parent`: not comparable — different semantic groups answer different questions
- - `s09_not_exists` vs `s08_join_distinct_parent`: not comparable — different semantic groups answer different questions
- - `s09_not_exists` vs `s11_optional_on`: not comparable — different semantic groups answer different questions
- - `s09_not_exists` vs `s16_page_offset`: not comparable — different semantic groups answer different questions
- - `s09_not_exists` vs `s17_page_keyset`: not comparable — different semantic groups answer different questions
- - `s09_not_exists` vs `s18_or_predicate`: not comparable — different semantic groups answer different questions
- - `s09_not_exists` vs `s19_union_all_predicate`: not comparable — different semantic groups answer different questions
- - `s09_not_exists` vs `s20_union_predicate`: not comparable — different semantic groups answer different questions
- - `s09_not_exists` vs `s14_topn_rownumber`: not comparable — different semantic groups answer different questions
- - `s09_not_exists` vs `s15_topn_lateral`: not comparable — different semantic groups answer different questions
- - `s09_not_exists` vs `a06_topn_distinct_on_sel`: not comparable — different semantic groups answer different questions
- - `s09_not_exists` vs `a04_topn_rownumber_sel`: not comparable — different semantic groups answer different questions
- - `s09_not_exists` vs `a05_topn_lateral_sel`: not comparable — different semantic groups answer different questions
- - `s10_left_anti` vs `s12_aggregate_join`: not comparable — different semantic groups answer different questions
- - `s10_left_anti` vs `s13_aggregate_preagg`: not comparable — different semantic groups answer different questions
- - `s10_left_anti` vs `a01_aggregate_join_sel`: not comparable — different semantic groups answer different questions
- - `s10_left_anti` vs `a02_aggregate_preagg_sel`: not comparable — different semantic groups answer different questions
- - `s10_left_anti` vs `a03_aggregate_correlated_sel`: not comparable — different semantic groups answer different questions
- - `s10_left_anti` vs `s01_join_where`: not comparable — different semantic groups answer different questions
- - `s10_left_anti` vs `s02_filtered_derived`: not comparable — different semantic groups answer different questions
- - `s10_left_anti` vs `s03_cte_filtered`: not comparable — different semantic groups answer different questions
- - `s10_left_anti` vs `s04_join_on`: not comparable — different semantic groups answer different questions
- - `s10_left_anti` vs `s05_comma_join`: not comparable — different semantic groups answer different questions
- - `s10_left_anti` vs `s06_exists_parent`: not comparable — different semantic groups answer different questions
- - `s10_left_anti` vs `s07_in_parent`: not comparable — different semantic groups answer different questions
- - `s10_left_anti` vs `s08_join_distinct_parent`: not comparable — different semantic groups answer different questions
- - `s10_left_anti` vs `s11_optional_on`: not comparable — different semantic groups answer different questions
- - `s10_left_anti` vs `s16_page_offset`: not comparable — different semantic groups answer different questions
- - `s10_left_anti` vs `s17_page_keyset`: not comparable — different semantic groups answer different questions
- - `s10_left_anti` vs `s18_or_predicate`: not comparable — different semantic groups answer different questions
- - `s10_left_anti` vs `s19_union_all_predicate`: not comparable — different semantic groups answer different questions
- - `s10_left_anti` vs `s20_union_predicate`: not comparable — different semantic groups answer different questions
- - `s10_left_anti` vs `s14_topn_rownumber`: not comparable — different semantic groups answer different questions
- - `s10_left_anti` vs `s15_topn_lateral`: not comparable — different semantic groups answer different questions
- - `s10_left_anti` vs `a06_topn_distinct_on_sel`: not comparable — different semantic groups answer different questions
- - `s10_left_anti` vs `a04_topn_rownumber_sel`: not comparable — different semantic groups answer different questions
- - `s10_left_anti` vs `a05_topn_lateral_sel`: not comparable — different semantic groups answer different questions
- - `s12_aggregate_join` vs `a01_aggregate_join_sel`: not comparable — different semantic groups answer different questions
- - `s12_aggregate_join` vs `a02_aggregate_preagg_sel`: not comparable — different semantic groups answer different questions
- - `s12_aggregate_join` vs `a03_aggregate_correlated_sel`: not comparable — different semantic groups answer different questions
- - `s12_aggregate_join` vs `s01_join_where`: not comparable — different semantic groups answer different questions
- - `s12_aggregate_join` vs `s02_filtered_derived`: not comparable — different semantic groups answer different questions
- - `s12_aggregate_join` vs `s03_cte_filtered`: not comparable — different semantic groups answer different questions
- - `s12_aggregate_join` vs `s04_join_on`: not comparable — different semantic groups answer different questions
- - `s12_aggregate_join` vs `s05_comma_join`: not comparable — different semantic groups answer different questions
- - `s12_aggregate_join` vs `s06_exists_parent`: not comparable — different semantic groups answer different questions
- - `s12_aggregate_join` vs `s07_in_parent`: not comparable — different semantic groups answer different questions
- - `s12_aggregate_join` vs `s08_join_distinct_parent`: not comparable — different semantic groups answer different questions
- - `s12_aggregate_join` vs `s11_optional_on`: not comparable — different semantic groups answer different questions
- - `s12_aggregate_join` vs `s16_page_offset`: not comparable — different semantic groups answer different questions
- - `s12_aggregate_join` vs `s17_page_keyset`: not comparable — different semantic groups answer different questions
- - `s12_aggregate_join` vs `s18_or_predicate`: not comparable — different semantic groups answer different questions
- - `s12_aggregate_join` vs `s19_union_all_predicate`: not comparable — different semantic groups answer different questions
- - `s12_aggregate_join` vs `s20_union_predicate`: not comparable — different semantic groups answer different questions
- - `s12_aggregate_join` vs `s14_topn_rownumber`: not comparable — different semantic groups answer different questions
- - `s12_aggregate_join` vs `s15_topn_lateral`: not comparable — different semantic groups answer different questions
- - `s12_aggregate_join` vs `a06_topn_distinct_on_sel`: not comparable — different semantic groups answer different questions
- - `s12_aggregate_join` vs `a04_topn_rownumber_sel`: not comparable — different semantic groups answer different questions
- - `s12_aggregate_join` vs `a05_topn_lateral_sel`: not comparable — different semantic groups answer different questions
- - `s13_aggregate_preagg` vs `a01_aggregate_join_sel`: not comparable — different semantic groups answer different questions
- - `s13_aggregate_preagg` vs `a02_aggregate_preagg_sel`: not comparable — different semantic groups answer different questions
- - `s13_aggregate_preagg` vs `a03_aggregate_correlated_sel`: not comparable — different semantic groups answer different questions
- - `s13_aggregate_preagg` vs `s01_join_where`: not comparable — different semantic groups answer different questions
- - `s13_aggregate_preagg` vs `s02_filtered_derived`: not comparable — different semantic groups answer different questions
- - `s13_aggregate_preagg` vs `s03_cte_filtered`: not comparable — different semantic groups answer different questions
- - `s13_aggregate_preagg` vs `s04_join_on`: not comparable — different semantic groups answer different questions
- - `s13_aggregate_preagg` vs `s05_comma_join`: not comparable — different semantic groups answer different questions
- - `s13_aggregate_preagg` vs `s06_exists_parent`: not comparable — different semantic groups answer different questions
- - `s13_aggregate_preagg` vs `s07_in_parent`: not comparable — different semantic groups answer different questions
- - `s13_aggregate_preagg` vs `s08_join_distinct_parent`: not comparable — different semantic groups answer different questions
- - `s13_aggregate_preagg` vs `s11_optional_on`: not comparable — different semantic groups answer different questions
- - `s13_aggregate_preagg` vs `s16_page_offset`: not comparable — different semantic groups answer different questions
- - `s13_aggregate_preagg` vs `s17_page_keyset`: not comparable — different semantic groups answer different questions
- - `s13_aggregate_preagg` vs `s18_or_predicate`: not comparable — different semantic groups answer different questions
- - `s13_aggregate_preagg` vs `s19_union_all_predicate`: not comparable — different semantic groups answer different questions
- - `s13_aggregate_preagg` vs `s20_union_predicate`: not comparable — different semantic groups answer different questions
- - `s13_aggregate_preagg` vs `s14_topn_rownumber`: not comparable — different semantic groups answer different questions
- - `s13_aggregate_preagg` vs `s15_topn_lateral`: not comparable — different semantic groups answer different questions
- - `s13_aggregate_preagg` vs `a06_topn_distinct_on_sel`: not comparable — different semantic groups answer different questions
- - `s13_aggregate_preagg` vs `a04_topn_rownumber_sel`: not comparable — different semantic groups answer different questions
- - `s13_aggregate_preagg` vs `a05_topn_lateral_sel`: not comparable — different semantic groups answer different questions
- - `a01_aggregate_join_sel` vs `s01_join_where`: not comparable — different semantic groups answer different questions
- - `a01_aggregate_join_sel` vs `s02_filtered_derived`: not comparable — different semantic groups answer different questions
- - `a01_aggregate_join_sel` vs `s03_cte_filtered`: not comparable — different semantic groups answer different questions
- - `a01_aggregate_join_sel` vs `s04_join_on`: not comparable — different semantic groups answer different questions
- - `a01_aggregate_join_sel` vs `s05_comma_join`: not comparable — different semantic groups answer different questions
- - `a01_aggregate_join_sel` vs `s06_exists_parent`: not comparable — different semantic groups answer different questions
- - `a01_aggregate_join_sel` vs `s07_in_parent`: not comparable — different semantic groups answer different questions
- - `a01_aggregate_join_sel` vs `s08_join_distinct_parent`: not comparable — different semantic groups answer different questions
- - `a01_aggregate_join_sel` vs `s11_optional_on`: not comparable — different semantic groups answer different questions
- - `a01_aggregate_join_sel` vs `s16_page_offset`: not comparable — different semantic groups answer different questions
- - `a01_aggregate_join_sel` vs `s17_page_keyset`: not comparable — different semantic groups answer different questions
- - `a01_aggregate_join_sel` vs `s18_or_predicate`: not comparable — different semantic groups answer different questions
- - `a01_aggregate_join_sel` vs `s19_union_all_predicate`: not comparable — different semantic groups answer different questions
- - `a01_aggregate_join_sel` vs `s20_union_predicate`: not comparable — different semantic groups answer different questions
- - `a01_aggregate_join_sel` vs `s14_topn_rownumber`: not comparable — different semantic groups answer different questions
- - `a01_aggregate_join_sel` vs `s15_topn_lateral`: not comparable — different semantic groups answer different questions
- - `a01_aggregate_join_sel` vs `a06_topn_distinct_on_sel`: not comparable — different semantic groups answer different questions
- - `a01_aggregate_join_sel` vs `a04_topn_rownumber_sel`: not comparable — different semantic groups answer different questions
- - `a01_aggregate_join_sel` vs `a05_topn_lateral_sel`: not comparable — different semantic groups answer different questions
- - `a02_aggregate_preagg_sel` vs `s01_join_where`: not comparable — different semantic groups answer different questions
- - `a02_aggregate_preagg_sel` vs `s02_filtered_derived`: not comparable — different semantic groups answer different questions
- - `a02_aggregate_preagg_sel` vs `s03_cte_filtered`: not comparable — different semantic groups answer different questions
- - `a02_aggregate_preagg_sel` vs `s04_join_on`: not comparable — different semantic groups answer different questions
- - `a02_aggregate_preagg_sel` vs `s05_comma_join`: not comparable — different semantic groups answer different questions
- - `a02_aggregate_preagg_sel` vs `s06_exists_parent`: not comparable — different semantic groups answer different questions
- - `a02_aggregate_preagg_sel` vs `s07_in_parent`: not comparable — different semantic groups answer different questions
- - `a02_aggregate_preagg_sel` vs `s08_join_distinct_parent`: not comparable — different semantic groups answer different questions
- - `a02_aggregate_preagg_sel` vs `s11_optional_on`: not comparable — different semantic groups answer different questions
- - `a02_aggregate_preagg_sel` vs `s16_page_offset`: not comparable — different semantic groups answer different questions
- - `a02_aggregate_preagg_sel` vs `s17_page_keyset`: not comparable — different semantic groups answer different questions
- - `a02_aggregate_preagg_sel` vs `s18_or_predicate`: not comparable — different semantic groups answer different questions
- - `a02_aggregate_preagg_sel` vs `s19_union_all_predicate`: not comparable — different semantic groups answer different questions
- - `a02_aggregate_preagg_sel` vs `s20_union_predicate`: not comparable — different semantic groups answer different questions
- - `a02_aggregate_preagg_sel` vs `s14_topn_rownumber`: not comparable — different semantic groups answer different questions
- - `a02_aggregate_preagg_sel` vs `s15_topn_lateral`: not comparable — different semantic groups answer different questions
- - `a02_aggregate_preagg_sel` vs `a06_topn_distinct_on_sel`: not comparable — different semantic groups answer different questions
- - `a02_aggregate_preagg_sel` vs `a04_topn_rownumber_sel`: not comparable — different semantic groups answer different questions
- - `a02_aggregate_preagg_sel` vs `a05_topn_lateral_sel`: not comparable — different semantic groups answer different questions
- - `a03_aggregate_correlated_sel` vs `s01_join_where`: not comparable — different semantic groups answer different questions
- - `a03_aggregate_correlated_sel` vs `s02_filtered_derived`: not comparable — different semantic groups answer different questions
- - `a03_aggregate_correlated_sel` vs `s03_cte_filtered`: not comparable — different semantic groups answer different questions
- - `a03_aggregate_correlated_sel` vs `s04_join_on`: not comparable — different semantic groups answer different questions
- - `a03_aggregate_correlated_sel` vs `s05_comma_join`: not comparable — different semantic groups answer different questions
- - `a03_aggregate_correlated_sel` vs `s06_exists_parent`: not comparable — different semantic groups answer different questions
- - `a03_aggregate_correlated_sel` vs `s07_in_parent`: not comparable — different semantic groups answer different questions
- - `a03_aggregate_correlated_sel` vs `s08_join_distinct_parent`: not comparable — different semantic groups answer different questions
- - `a03_aggregate_correlated_sel` vs `s11_optional_on`: not comparable — different semantic groups answer different questions
- - `a03_aggregate_correlated_sel` vs `s16_page_offset`: not comparable — different semantic groups answer different questions
- - `a03_aggregate_correlated_sel` vs `s17_page_keyset`: not comparable — different semantic groups answer different questions
- - `a03_aggregate_correlated_sel` vs `s18_or_predicate`: not comparable — different semantic groups answer different questions
- - `a03_aggregate_correlated_sel` vs `s19_union_all_predicate`: not comparable — different semantic groups answer different questions
- - `a03_aggregate_correlated_sel` vs `s20_union_predicate`: not comparable — different semantic groups answer different questions
- - `a03_aggregate_correlated_sel` vs `s14_topn_rownumber`: not comparable — different semantic groups answer different questions
- - `a03_aggregate_correlated_sel` vs `s15_topn_lateral`: not comparable — different semantic groups answer different questions
- - `a03_aggregate_correlated_sel` vs `a06_topn_distinct_on_sel`: not comparable — different semantic groups answer different questions
- - `a03_aggregate_correlated_sel` vs `a04_topn_rownumber_sel`: not comparable — different semantic groups answer different questions
- - `a03_aggregate_correlated_sel` vs `a05_topn_lateral_sel`: not comparable — different semantic groups answer different questions
- - `s01_join_where` vs `s06_exists_parent`: not comparable — different semantic groups answer different questions
- - `s01_join_where` vs `s07_in_parent`: not comparable — different semantic groups answer different questions
- - `s01_join_where` vs `s08_join_distinct_parent`: not comparable — different semantic groups answer different questions
- - `s01_join_where` vs `s11_optional_on`: not comparable — different semantic groups answer different questions
- - `s01_join_where` vs `s16_page_offset`: not comparable — different semantic groups answer different questions
- - `s01_join_where` vs `s17_page_keyset`: not comparable — different semantic groups answer different questions
- - `s01_join_where` vs `s18_or_predicate`: not comparable — different semantic groups answer different questions
- - `s01_join_where` vs `s19_union_all_predicate`: not comparable — different semantic groups answer different questions
- - `s01_join_where` vs `s20_union_predicate`: not comparable — different semantic groups answer different questions
- - `s01_join_where` vs `s14_topn_rownumber`: not comparable — different semantic groups answer different questions
- - `s01_join_where` vs `s15_topn_lateral`: not comparable — different semantic groups answer different questions
- - `s01_join_where` vs `a06_topn_distinct_on_sel`: not comparable — different semantic groups answer different questions
- - `s01_join_where` vs `a04_topn_rownumber_sel`: not comparable — different semantic groups answer different questions
- - `s01_join_where` vs `a05_topn_lateral_sel`: not comparable — different semantic groups answer different questions
- - `s02_filtered_derived` vs `s06_exists_parent`: not comparable — different semantic groups answer different questions
- - `s02_filtered_derived` vs `s07_in_parent`: not comparable — different semantic groups answer different questions
- - `s02_filtered_derived` vs `s08_join_distinct_parent`: not comparable — different semantic groups answer different questions
- - `s02_filtered_derived` vs `s11_optional_on`: not comparable — different semantic groups answer different questions
- - `s02_filtered_derived` vs `s16_page_offset`: not comparable — different semantic groups answer different questions
- - `s02_filtered_derived` vs `s17_page_keyset`: not comparable — different semantic groups answer different questions
- - `s02_filtered_derived` vs `s18_or_predicate`: not comparable — different semantic groups answer different questions
- - `s02_filtered_derived` vs `s19_union_all_predicate`: not comparable — different semantic groups answer different questions
- - `s02_filtered_derived` vs `s20_union_predicate`: not comparable — different semantic groups answer different questions
- - `s02_filtered_derived` vs `s14_topn_rownumber`: not comparable — different semantic groups answer different questions
- - `s02_filtered_derived` vs `s15_topn_lateral`: not comparable — different semantic groups answer different questions
- - `s02_filtered_derived` vs `a06_topn_distinct_on_sel`: not comparable — different semantic groups answer different questions
- - `s02_filtered_derived` vs `a04_topn_rownumber_sel`: not comparable — different semantic groups answer different questions
- - `s02_filtered_derived` vs `a05_topn_lateral_sel`: not comparable — different semantic groups answer different questions
- - `s03_cte_filtered` vs `s06_exists_parent`: not comparable — different semantic groups answer different questions
- - `s03_cte_filtered` vs `s07_in_parent`: not comparable — different semantic groups answer different questions
- - `s03_cte_filtered` vs `s08_join_distinct_parent`: not comparable — different semantic groups answer different questions
- - `s03_cte_filtered` vs `s11_optional_on`: not comparable — different semantic groups answer different questions
- - `s03_cte_filtered` vs `s16_page_offset`: not comparable — different semantic groups answer different questions
- - `s03_cte_filtered` vs `s17_page_keyset`: not comparable — different semantic groups answer different questions
- - `s03_cte_filtered` vs `s18_or_predicate`: not comparable — different semantic groups answer different questions
- - `s03_cte_filtered` vs `s19_union_all_predicate`: not comparable — different semantic groups answer different questions
- - `s03_cte_filtered` vs `s20_union_predicate`: not comparable — different semantic groups answer different questions
- - `s03_cte_filtered` vs `s14_topn_rownumber`: not comparable — different semantic groups answer different questions
- - `s03_cte_filtered` vs `s15_topn_lateral`: not comparable — different semantic groups answer different questions
- - `s03_cte_filtered` vs `a06_topn_distinct_on_sel`: not comparable — different semantic groups answer different questions
- - `s03_cte_filtered` vs `a04_topn_rownumber_sel`: not comparable — different semantic groups answer different questions
- - `s03_cte_filtered` vs `a05_topn_lateral_sel`: not comparable — different semantic groups answer different questions
- - `s04_join_on` vs `s06_exists_parent`: not comparable — different semantic groups answer different questions
- - `s04_join_on` vs `s07_in_parent`: not comparable — different semantic groups answer different questions
- - `s04_join_on` vs `s08_join_distinct_parent`: not comparable — different semantic groups answer different questions
- - `s04_join_on` vs `s11_optional_on`: not comparable — different semantic groups answer different questions
- - `s04_join_on` vs `s16_page_offset`: not comparable — different semantic groups answer different questions
- - `s04_join_on` vs `s17_page_keyset`: not comparable — different semantic groups answer different questions
- - `s04_join_on` vs `s18_or_predicate`: not comparable — different semantic groups answer different questions
- - `s04_join_on` vs `s19_union_all_predicate`: not comparable — different semantic groups answer different questions
- - `s04_join_on` vs `s20_union_predicate`: not comparable — different semantic groups answer different questions
- - `s04_join_on` vs `s14_topn_rownumber`: not comparable — different semantic groups answer different questions
- - `s04_join_on` vs `s15_topn_lateral`: not comparable — different semantic groups answer different questions
- - `s04_join_on` vs `a06_topn_distinct_on_sel`: not comparable — different semantic groups answer different questions
- - `s04_join_on` vs `a04_topn_rownumber_sel`: not comparable — different semantic groups answer different questions
- - `s04_join_on` vs `a05_topn_lateral_sel`: not comparable — different semantic groups answer different questions
- - `s05_comma_join` vs `s06_exists_parent`: not comparable — different semantic groups answer different questions
- - `s05_comma_join` vs `s07_in_parent`: not comparable — different semantic groups answer different questions
- - `s05_comma_join` vs `s08_join_distinct_parent`: not comparable — different semantic groups answer different questions
- - `s05_comma_join` vs `s11_optional_on`: not comparable — different semantic groups answer different questions
- - `s05_comma_join` vs `s16_page_offset`: not comparable — different semantic groups answer different questions
- - `s05_comma_join` vs `s17_page_keyset`: not comparable — different semantic groups answer different questions
- - `s05_comma_join` vs `s18_or_predicate`: not comparable — different semantic groups answer different questions
- - `s05_comma_join` vs `s19_union_all_predicate`: not comparable — different semantic groups answer different questions
- - `s05_comma_join` vs `s20_union_predicate`: not comparable — different semantic groups answer different questions
- - `s05_comma_join` vs `s14_topn_rownumber`: not comparable — different semantic groups answer different questions
- - `s05_comma_join` vs `s15_topn_lateral`: not comparable — different semantic groups answer different questions
- - `s05_comma_join` vs `a06_topn_distinct_on_sel`: not comparable — different semantic groups answer different questions
- - `s05_comma_join` vs `a04_topn_rownumber_sel`: not comparable — different semantic groups answer different questions
- - `s05_comma_join` vs `a05_topn_lateral_sel`: not comparable — different semantic groups answer different questions
- - `s06_exists_parent` vs `s11_optional_on`: not comparable — different semantic groups answer different questions
- - `s06_exists_parent` vs `s16_page_offset`: not comparable — different semantic groups answer different questions
- - `s06_exists_parent` vs `s17_page_keyset`: not comparable — different semantic groups answer different questions
- - `s06_exists_parent` vs `s18_or_predicate`: not comparable — different semantic groups answer different questions
- - `s06_exists_parent` vs `s19_union_all_predicate`: not comparable — different semantic groups answer different questions
- - `s06_exists_parent` vs `s20_union_predicate`: not comparable — different semantic groups answer different questions
- - `s06_exists_parent` vs `s14_topn_rownumber`: not comparable — different semantic groups answer different questions
- - `s06_exists_parent` vs `s15_topn_lateral`: not comparable — different semantic groups answer different questions
- - `s06_exists_parent` vs `a06_topn_distinct_on_sel`: not comparable — different semantic groups answer different questions
- - `s06_exists_parent` vs `a04_topn_rownumber_sel`: not comparable — different semantic groups answer different questions
- - `s06_exists_parent` vs `a05_topn_lateral_sel`: not comparable — different semantic groups answer different questions
- - `s07_in_parent` vs `s11_optional_on`: not comparable — different semantic groups answer different questions
- - `s07_in_parent` vs `s16_page_offset`: not comparable — different semantic groups answer different questions
- - `s07_in_parent` vs `s17_page_keyset`: not comparable — different semantic groups answer different questions
- - `s07_in_parent` vs `s18_or_predicate`: not comparable — different semantic groups answer different questions
- - `s07_in_parent` vs `s19_union_all_predicate`: not comparable — different semantic groups answer different questions
- - `s07_in_parent` vs `s20_union_predicate`: not comparable — different semantic groups answer different questions
- - `s07_in_parent` vs `s14_topn_rownumber`: not comparable — different semantic groups answer different questions
- - `s07_in_parent` vs `s15_topn_lateral`: not comparable — different semantic groups answer different questions
- - `s07_in_parent` vs `a06_topn_distinct_on_sel`: not comparable — different semantic groups answer different questions
- - `s07_in_parent` vs `a04_topn_rownumber_sel`: not comparable — different semantic groups answer different questions
- - `s07_in_parent` vs `a05_topn_lateral_sel`: not comparable — different semantic groups answer different questions
- - `s08_join_distinct_parent` vs `s11_optional_on`: not comparable — different semantic groups answer different questions
- - `s08_join_distinct_parent` vs `s16_page_offset`: not comparable — different semantic groups answer different questions
- - `s08_join_distinct_parent` vs `s17_page_keyset`: not comparable — different semantic groups answer different questions
- - `s08_join_distinct_parent` vs `s18_or_predicate`: not comparable — different semantic groups answer different questions
- - `s08_join_distinct_parent` vs `s19_union_all_predicate`: not comparable — different semantic groups answer different questions
- - `s08_join_distinct_parent` vs `s20_union_predicate`: not comparable — different semantic groups answer different questions
- - `s08_join_distinct_parent` vs `s14_topn_rownumber`: not comparable — different semantic groups answer different questions
- - `s08_join_distinct_parent` vs `s15_topn_lateral`: not comparable — different semantic groups answer different questions
- - `s08_join_distinct_parent` vs `a06_topn_distinct_on_sel`: not comparable — different semantic groups answer different questions
- - `s08_join_distinct_parent` vs `a04_topn_rownumber_sel`: not comparable — different semantic groups answer different questions
- - `s08_join_distinct_parent` vs `a05_topn_lateral_sel`: not comparable — different semantic groups answer different questions
- - `s11_optional_on` vs `s16_page_offset`: not comparable — different semantic groups answer different questions
- - `s11_optional_on` vs `s17_page_keyset`: not comparable — different semantic groups answer different questions
- - `s11_optional_on` vs `s18_or_predicate`: not comparable — different semantic groups answer different questions
- - `s11_optional_on` vs `s19_union_all_predicate`: not comparable — different semantic groups answer different questions
- - `s11_optional_on` vs `s20_union_predicate`: not comparable — different semantic groups answer different questions
- - `s11_optional_on` vs `s14_topn_rownumber`: not comparable — different semantic groups answer different questions
- - `s11_optional_on` vs `s15_topn_lateral`: not comparable — different semantic groups answer different questions
- - `s11_optional_on` vs `a06_topn_distinct_on_sel`: not comparable — different semantic groups answer different questions
- - `s11_optional_on` vs `a04_topn_rownumber_sel`: not comparable — different semantic groups answer different questions
- - `s11_optional_on` vs `a05_topn_lateral_sel`: not comparable — different semantic groups answer different questions
- - `s16_page_offset` vs `s18_or_predicate`: not comparable — different semantic groups answer different questions
- - `s16_page_offset` vs `s19_union_all_predicate`: not comparable — different semantic groups answer different questions
- - `s16_page_offset` vs `s20_union_predicate`: not comparable — different semantic groups answer different questions
- - `s16_page_offset` vs `s14_topn_rownumber`: not comparable — different semantic groups answer different questions
- - `s16_page_offset` vs `s15_topn_lateral`: not comparable — different semantic groups answer different questions
- - `s16_page_offset` vs `a06_topn_distinct_on_sel`: not comparable — different semantic groups answer different questions
- - `s16_page_offset` vs `a04_topn_rownumber_sel`: not comparable — different semantic groups answer different questions
- - `s16_page_offset` vs `a05_topn_lateral_sel`: not comparable — different semantic groups answer different questions
- - `s17_page_keyset` vs `s18_or_predicate`: not comparable — different semantic groups answer different questions
- - `s17_page_keyset` vs `s19_union_all_predicate`: not comparable — different semantic groups answer different questions
- - `s17_page_keyset` vs `s20_union_predicate`: not comparable — different semantic groups answer different questions
- - `s17_page_keyset` vs `s14_topn_rownumber`: not comparable — different semantic groups answer different questions
- - `s17_page_keyset` vs `s15_topn_lateral`: not comparable — different semantic groups answer different questions
- - `s17_page_keyset` vs `a06_topn_distinct_on_sel`: not comparable — different semantic groups answer different questions
- - `s17_page_keyset` vs `a04_topn_rownumber_sel`: not comparable — different semantic groups answer different questions
- - `s17_page_keyset` vs `a05_topn_lateral_sel`: not comparable — different semantic groups answer different questions
- - `s18_or_predicate` vs `s14_topn_rownumber`: not comparable — different semantic groups answer different questions
- - `s18_or_predicate` vs `s15_topn_lateral`: not comparable — different semantic groups answer different questions
- - `s18_or_predicate` vs `a06_topn_distinct_on_sel`: not comparable — different semantic groups answer different questions
- - `s18_or_predicate` vs `a04_topn_rownumber_sel`: not comparable — different semantic groups answer different questions
- - `s18_or_predicate` vs `a05_topn_lateral_sel`: not comparable — different semantic groups answer different questions
- - `s19_union_all_predicate` vs `s14_topn_rownumber`: not comparable — different semantic groups answer different questions
- - `s19_union_all_predicate` vs `s15_topn_lateral`: not comparable — different semantic groups answer different questions
- - `s19_union_all_predicate` vs `a06_topn_distinct_on_sel`: not comparable — different semantic groups answer different questions
- - `s19_union_all_predicate` vs `a04_topn_rownumber_sel`: not comparable — different semantic groups answer different questions
- - `s19_union_all_predicate` vs `a05_topn_lateral_sel`: not comparable — different semantic groups answer different questions
- - `s20_union_predicate` vs `s14_topn_rownumber`: not comparable — different semantic groups answer different questions
- - `s20_union_predicate` vs `s15_topn_lateral`: not comparable — different semantic groups answer different questions
- - `s20_union_predicate` vs `a06_topn_distinct_on_sel`: not comparable — different semantic groups answer different questions
- - `s20_union_predicate` vs `a04_topn_rownumber_sel`: not comparable — different semantic groups answer different questions
- - `s20_union_predicate` vs `a05_topn_lateral_sel`: not comparable — different semantic groups answer different questions
- - `s14_topn_rownumber` vs `a06_topn_distinct_on_sel`: not comparable — different semantic groups answer different questions
- - `s14_topn_rownumber` vs `a04_topn_rownumber_sel`: not comparable — different semantic groups answer different questions
- - `s14_topn_rownumber` vs `a05_topn_lateral_sel`: not comparable — different semantic groups answer different questions
- - `s15_topn_lateral` vs `a06_topn_distinct_on_sel`: not comparable — different semantic groups answer different questions
- - `s15_topn_lateral` vs `a04_topn_rownumber_sel`: not comparable — different semantic groups answer different questions
- - `s15_topn_lateral` vs `a05_topn_lateral_sel`: not comparable — different semantic groups answer different questions
- - `a06_topn_distinct_on_sel` vs `a04_topn_rownumber_sel`: not comparable — different semantic groups answer different questions
- - `a06_topn_distinct_on_sel` vs `a05_topn_lateral_sel`: not comparable — different semantic groups answer different questions

## Predictions this run informs

- P1-P5 join/CTE spellings
- P6-P7 membership/absence
- P8-P9 aggregates/top-N
- P10-P11 predicates/pagination

Verdicts are assigned in the signed analysis, not here.
