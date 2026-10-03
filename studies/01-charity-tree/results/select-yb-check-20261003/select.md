# SELECT comparison (generated measurements)

_This file contains measurements only. Interpretation belongs in a signed analysis._

**TL;DR:** a validation run supplies correctness and plan evidence. Within-run samples do not establish a replicated performance winner. Missing measurements are N/A.

## Provenance

| Field | Value |
|---|---|
| repo commit | `f37becf6793044cd31c8c35da906e5e9984ef351` |
| benchmark image | `localhost/charitybench-select:select-yb-check-20261003` |
| image digest | `923874f11103cdf30563e4edc032c74d943db7fc14a0bd8274883ef4f218b4ca` |
| run tag | `run/study-01/select-yb-check-20261003` |
| environment | `xiles-podman-select-validation-20261003` |
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
| select-yb-check-20261003 | single | s09_not_exists | sel_absence | valid | 1.481069 | 15.1 | 3 | `6b553eafd5f0122d` |
| select-yb-check-20261003 | single | s10_left_anti | sel_absence | valid | 0.859673 | 26.0 | 3 | `a2f6838ee80f8792` |
| select-yb-check-20261003 | single | s12_aggregate_join | sel_aggregate | valid | 1.079219 | 54.3 | 3 | `2a550dc1499b665c` |
| select-yb-check-20261003 | single | s13_aggregate_preagg | sel_aggregate | valid | 0.902863 | 61.2 | 3 | `120182ce300978d4` |
| select-yb-check-20261003 | single | s01_join_where | sel_full_rows | valid | 1.399511 | 26.5 | 3 | `eb7e7141c15ed0a0` |
| select-yb-check-20261003 | single | s02_filtered_derived | sel_full_rows | valid | 1.523788 | 34.7 | 3 | `8aa68d62c571b668` |
| select-yb-check-20261003 | single | s03_cte_filtered | sel_full_rows | valid | 1.067563 | 56.3 | 3 | `8aa68d62c571b668` |
| select-yb-check-20261003 | single | s04_join_on | sel_full_rows | valid | 0.891458 | 49.0 | 3 | `eb7e7141c15ed0a0` |
| select-yb-check-20261003 | single | s05_comma_join | sel_full_rows | valid | 1.106223 | 64.4 | 3 | `eb7e7141c15ed0a0` |
| select-yb-check-20261003 | single | s06_exists_parent | sel_membership | valid | 1.143520 | 16.0 | 3 | `39f4676f00c0b2cb` |
| select-yb-check-20261003 | single | s07_in_parent | sel_membership | valid | 1.278840 | 25.3 | 3 | `39f4676f00c0b2cb` |
| select-yb-check-20261003 | single | s08_join_distinct_parent | sel_membership | valid | 1.074368 | 37.9 | 3 | `28c86dff835a382c` |
| select-yb-check-20261003 | single | s11_optional_on | sel_optional_match | valid | 0.854083 | 3.6 | 3 | `a2f6838ee80f8792` |
| select-yb-check-20261003 | single | s16_page_offset | sel_pagination | valid | 0.406113 | 20.3 | 3 | `ac4509d295bf3348` |
| select-yb-check-20261003 | single | s17_page_keyset | sel_pagination | valid | 0.700543 | 70.5 | 3 | `ac4509d295bf3348` |
| select-yb-check-20261003 | single | s14_topn_rownumber | sel_topn | valid | 0.886198 | 18.8 | 3 | `5a870d51d257d840` |
| select-yb-check-20261003 | single | s15_topn_lateral | sel_topn | valid | 1.522510 | 43.3 | 3 | `9983c096f9eefa78` |

## Excluded cells

None: every cell is valid.

## Comparisons

### Same plan (no plan difference to explain)

- `s01_join_where` vs `s04_join_on` on single: same plan structure
- `s01_join_where` vs `s05_comma_join` on single: same plan structure
- `s02_filtered_derived` vs `s03_cte_filtered` on single: same plan structure
- `s04_join_on` vs `s05_comma_join` on single: same plan structure
- `s06_exists_parent` vs `s07_in_parent` on single: same plan structure
- `s16_page_offset` vs `s17_page_keyset` on single: same plan structure

### No resolved difference

- `s12_aggregate_join` vs `s13_aggregate_preagg` on single: medians differ by 0.176 ms, inside the 54.3%/61.2% spread — no resolved difference
- `s01_join_where` vs `s02_filtered_derived` on single: medians differ by 0.124 ms, inside the 26.5%/34.7% spread — no resolved difference
- `s01_join_where` vs `s03_cte_filtered` on single: medians differ by 0.332 ms, inside the 26.5%/56.3% spread — no resolved difference
- `s02_filtered_derived` vs `s04_join_on` on single: medians differ by 0.632 ms, inside the 34.7%/49.0% spread — no resolved difference
- `s02_filtered_derived` vs `s05_comma_join` on single: medians differ by 0.418 ms, inside the 34.7%/64.4% spread — no resolved difference
- `s03_cte_filtered` vs `s04_join_on` on single: medians differ by 0.176 ms, inside the 56.3%/49.0% spread — no resolved difference
- `s03_cte_filtered` vs `s05_comma_join` on single: medians differ by 0.039 ms, inside the 56.3%/64.4% spread — no resolved difference
- `s06_exists_parent` vs `s08_join_distinct_parent` on single: medians differ by 0.069 ms, inside the 16.0%/37.9% spread — no resolved difference
- `s07_in_parent` vs `s08_join_distinct_parent` on single: medians differ by 0.204 ms, inside the 25.3%/37.9% spread — no resolved difference

### Within-run differences (not replicated)

- `s09_not_exists` vs `s10_left_anti` on single: within-run medians 1.481069 / 0.859673 ms; whole-run replication required — no resolved difference
- `s14_topn_rownumber` vs `s15_topn_lateral` on single: within-run medians 0.886198 / 1.522510 ms; whole-run replication required — no resolved difference

## Live negative controls

- `x01_optional_where`: detected=true, rows=3
- `x02_not_in_nullable`: detected=true, rows=0
- `x03_distinct_on_name`: detected=true, rows=1
- `x04_limit_before_join`: detected=true, rows=2

## Captured plans and result sizes

### s09_not_exists

Rows: 0; serialized decoded-row JSON bytes: 4 (not wire bytes). Planning/execution: 0.090000 / 0.620000 ms.

```text
Hash Anti Join  (cost=42.49..42.67 rows=1 width=7) (actual time=0.582..0.583 rows=0 loops=1)
  Hash Cond: (c.id = o.customer_id)
  ->  Seq Scan on customers c  (cost=21.05..21.21 rows=3 width=7) (actual time=0.237..0.238 rows=3 loops=1)
  ->  Hash  (cost=21.38..21.38 rows=5 width=4) (actual time=0.333..0.333 rows=4 loops=1)
        Buckets: 1024  Batches: 1  Memory Usage: 9kB
        ->  Seq Scan on orders o  (cost=21.05..21.38 rows=5 width=4) (actual time=0.319..0.322 rows=5 loops=1)
              Storage Filter: (amount > 100)
Planning Time: 0.090 ms
Execution Time: 0.620 ms
Peak Memory Usage: 205 kB

```

### s10_left_anti

Rows: 0; serialized decoded-row JSON bytes: 4 (not wire bytes). Planning/execution: 0.199000 / 0.975000 ms.

```text
Hash Right Join  (cost=42.30..42.70 rows=1 width=7) (actual time=0.927..0.929 rows=0 loops=1)
  Hash Cond: (o.customer_id = c.id)
  Filter: (o.id IS NULL)
  Rows Removed by Filter: 4
  ->  Seq Scan on orders o  (cost=21.05..21.40 rows=5 width=8) (actual time=0.318..0.321 rows=5 loops=1)
        Storage Filter: (amount > 100)
  ->  Hash  (cost=21.21..21.21 rows=3 width=7) (actual time=0.593..0.594 rows=3 loops=1)
        Buckets: 1024  Batches: 1  Memory Usage: 9kB
        ->  Seq Scan on customers c  (cost=21.05..21.21 rows=3 width=7) (actual time=0.579..0.582 rows=3 loops=1)
Planning Time: 0.199 ms
Execution Time: 0.975 ms
Peak Memory Usage: 191 kB

```

### s12_aggregate_join

Rows: 3; serialized decoded-row JSON bytes: 46 (not wire bytes). Planning/execution: 0.099000 / 0.464000 ms.

```text
GroupAggregate  (cost=42.72..42.89 rows=3 width=23) (actual time=0.415..0.420 rows=3 loops=1)
  Group Key: c.id
  ->  Merge Left Join  (cost=42.72..42.82 rows=6 width=15) (actual time=0.410..0.413 rows=5 loops=1)
        Merge Cond: (c.id = o.customer_id)
        ->  Sort  (cost=21.23..21.24 rows=3 width=7) (actual time=0.228..0.229 rows=3 loops=1)
              Sort Key: c.id
              Sort Method: quicksort  Memory: 25kB
              ->  Seq Scan on customers c  (cost=21.05..21.21 rows=3 width=7) (actual time=0.218..0.220 rows=3 loops=1)
        ->  Sort  (cost=21.49..21.51 rows=6 width=12) (actual time=0.179..0.179 rows=6 loops=1)
              Sort Key: o.customer_id
              Sort Method: quicksort  Memory: 25kB
              ->  Seq Scan on orders o  (cost=21.05..21.41 rows=6 width=12) (actual time=0.166..0.168 rows=6 loops=1)
Planning Time: 0.099 ms
Execution Time: 0.464 ms
Peak Memory Usage: 206 kB

```

### s13_aggregate_preagg

Rows: 3; serialized decoded-row JSON bytes: 46 (not wire bytes). Planning/execution: 0.138000 / 0.577000 ms.

```text
Hash Left Join  (cost=42.67..42.85 rows=3 width=23) (actual time=0.525..0.528 rows=3 loops=1)
  Hash Cond: (c.id = a.customer_id)
  ->  Seq Scan on customers c  (cost=21.05..21.21 rows=3 width=7) (actual time=0.260..0.261 rows=3 loops=1)
  ->  Hash  (cost=21.59..21.59 rows=3 width=20) (actual time=0.253..0.254 rows=3 loops=1)
        Buckets: 1024  Batches: 1  Memory Usage: 9kB
        ->  Subquery Scan on a  (cost=21.47..21.59 rows=3 width=20) (actual time=0.242..0.246 rows=4 loops=1)
              ->  GroupAggregate  (cost=21.47..21.56 rows=3 width=20) (actual time=0.241..0.244 rows=4 loops=1)
                    Group Key: orders.customer_id
                    ->  Sort  (cost=21.47..21.48 rows=6 width=8) (actual time=0.236..0.237 rows=6 loops=1)
                          Sort Key: orders.customer_id
                          Sort Method: quicksort  Memory: 25kB
                          ->  Seq Scan on orders  (cost=21.05..21.39 rows=6 width=8) (actual time=0.223..0.225 rows=6 loops=1)
Planning Time: 0.138 ms
Execution Time: 0.577 ms
Peak Memory Usage: 232 kB

```

### s01_join_where

Rows: 4; serialized decoded-row JSON bytes: 73 (not wire bytes). Planning/execution: 0.143000 / 0.576000 ms.

```text
Hash Join  (cost=42.30..42.72 rows=5 width=15) (actual time=0.532..0.536 rows=4 loops=1)
  Hash Cond: (o.customer_id = c.id)
  ->  Seq Scan on orders o  (cost=21.05..21.42 rows=5 width=12) (actual time=0.329..0.330 rows=5 loops=1)
        Storage Filter: (amount > 100)
  ->  Hash  (cost=21.21..21.21 rows=3 width=7) (actual time=0.192..0.192 rows=3 loops=1)
        Buckets: 1024  Batches: 1  Memory Usage: 9kB
        ->  Seq Scan on customers c  (cost=21.05..21.21 rows=3 width=7) (actual time=0.183..0.185 rows=3 loops=1)
Planning Time: 0.143 ms
Execution Time: 0.576 ms
Peak Memory Usage: 205 kB

```

### s02_filtered_derived

Rows: 4; serialized decoded-row JSON bytes: 73 (not wire bytes). Planning/execution: 0.138000 / 0.635000 ms.

```text
Hash Join  (cost=42.30..42.72 rows=5 width=15) (actual time=0.592..0.596 rows=4 loops=1)
  Hash Cond: (orders.customer_id = c.id)
  ->  Seq Scan on orders  (cost=21.05..21.42 rows=5 width=12) (actual time=0.310..0.312 rows=5 loops=1)
        Storage Filter: (amount > 100)
  ->  Hash  (cost=21.21..21.21 rows=3 width=7) (actual time=0.272..0.273 rows=3 loops=1)
        Buckets: 1024  Batches: 1  Memory Usage: 9kB
        ->  Seq Scan on customers c  (cost=21.05..21.21 rows=3 width=7) (actual time=0.260..0.263 rows=3 loops=1)
Planning Time: 0.138 ms
Execution Time: 0.635 ms
Peak Memory Usage: 205 kB

```

### s03_cte_filtered

Rows: 4; serialized decoded-row JSON bytes: 73 (not wire bytes). Planning/execution: 0.139000 / 0.537000 ms.

```text
Hash Join  (cost=42.30..42.72 rows=5 width=15) (actual time=0.495..0.500 rows=4 loops=1)
  Hash Cond: (orders.customer_id = c.id)
  ->  Seq Scan on orders  (cost=21.05..21.42 rows=5 width=12) (actual time=0.268..0.270 rows=5 loops=1)
        Storage Filter: (amount > 100)
  ->  Hash  (cost=21.21..21.21 rows=3 width=7) (actual time=0.218..0.218 rows=3 loops=1)
        Buckets: 1024  Batches: 1  Memory Usage: 9kB
        ->  Seq Scan on customers c  (cost=21.05..21.21 rows=3 width=7) (actual time=0.206..0.208 rows=3 loops=1)
Planning Time: 0.139 ms
Execution Time: 0.537 ms
Peak Memory Usage: 205 kB

```

### s04_join_on

Rows: 4; serialized decoded-row JSON bytes: 73 (not wire bytes). Planning/execution: 0.147000 / 0.991000 ms.

```text
Hash Join  (cost=42.30..42.72 rows=5 width=15) (actual time=0.663..0.666 rows=4 loops=1)
  Hash Cond: (o.customer_id = c.id)
  ->  Seq Scan on orders o  (cost=21.05..21.42 rows=5 width=12) (actual time=0.425..0.426 rows=5 loops=1)
        Storage Filter: (amount > 100)
  ->  Hash  (cost=21.21..21.21 rows=3 width=7) (actual time=0.223..0.223 rows=3 loops=1)
        Buckets: 1024  Batches: 1  Memory Usage: 9kB
        ->  Seq Scan on customers c  (cost=21.05..21.21 rows=3 width=7) (actual time=0.215..0.217 rows=3 loops=1)
Planning Time: 0.147 ms
Execution Time: 0.991 ms
Peak Memory Usage: 205 kB

```

### s05_comma_join

Rows: 4; serialized decoded-row JSON bytes: 73 (not wire bytes). Planning/execution: 0.172000 / 0.537000 ms.

```text
Hash Join  (cost=42.30..42.72 rows=5 width=15) (actual time=0.489..0.493 rows=4 loops=1)
  Hash Cond: (o.customer_id = c.id)
  ->  Seq Scan on orders o  (cost=21.05..21.42 rows=5 width=12) (actual time=0.275..0.277 rows=5 loops=1)
        Storage Filter: (amount > 100)
  ->  Hash  (cost=21.21..21.21 rows=3 width=7) (actual time=0.201..0.201 rows=3 loops=1)
        Buckets: 1024  Batches: 1  Memory Usage: 9kB
        ->  Seq Scan on customers c  (cost=21.05..21.21 rows=3 width=7) (actual time=0.191..0.193 rows=3 loops=1)
Planning Time: 0.172 ms
Execution Time: 0.537 ms
Peak Memory Usage: 205 kB

```

### s06_exists_parent

Rows: 3; serialized decoded-row JSON bytes: 28 (not wire bytes). Planning/execution: 0.163000 / 0.900000 ms.

```text
Hash Semi Join  (cost=42.49..42.70 rows=3 width=7) (actual time=0.853..0.858 rows=3 loops=1)
  Hash Cond: (c.id = o.customer_id)
  ->  Seq Scan on customers c  (cost=21.05..21.21 rows=3 width=7) (actual time=0.419..0.420 rows=3 loops=1)
  ->  Hash  (cost=21.38..21.38 rows=5 width=4) (actual time=0.419..0.419 rows=4 loops=1)
        Buckets: 1024  Batches: 1  Memory Usage: 9kB
        ->  Seq Scan on orders o  (cost=21.05..21.38 rows=5 width=4) (actual time=0.404..0.407 rows=5 loops=1)
              Storage Filter: (amount > 100)
Planning Time: 0.163 ms
Execution Time: 0.900 ms
Peak Memory Usage: 205 kB

```

### s07_in_parent

Rows: 3; serialized decoded-row JSON bytes: 28 (not wire bytes). Planning/execution: 0.244000 / 0.776000 ms.

```text
Hash Semi Join  (cost=42.49..42.70 rows=3 width=7) (actual time=0.728..0.731 rows=3 loops=1)
  Hash Cond: (c.id = o.customer_id)
  ->  Seq Scan on customers c  (cost=21.05..21.21 rows=3 width=7) (actual time=0.298..0.299 rows=3 loops=1)
  ->  Hash  (cost=21.38..21.38 rows=5 width=4) (actual time=0.417..0.417 rows=4 loops=1)
        Buckets: 1024  Batches: 1  Memory Usage: 9kB
        ->  Seq Scan on orders o  (cost=21.05..21.38 rows=5 width=4) (actual time=0.406..0.409 rows=5 loops=1)
              Storage Filter: (amount > 100)
Planning Time: 0.244 ms
Execution Time: 0.776 ms
Peak Memory Usage: 205 kB

```

### s08_join_distinct_parent

Rows: 3; serialized decoded-row JSON bytes: 28 (not wire bytes). Planning/execution: 0.137000 / 0.635000 ms.

```text
Unique  (cost=42.81..42.84 rows=3 width=7) (actual time=0.589..0.593 rows=3 loops=1)
  ->  Sort  (cost=42.81..42.82 rows=5 width=7) (actual time=0.589..0.590 rows=4 loops=1)
        Sort Key: c.id, c.name
        Sort Method: quicksort  Memory: 25kB
        ->  Merge Join  (cost=42.67..42.75 rows=5 width=7) (actual time=0.577..0.581 rows=4 loops=1)
              Merge Cond: (c.id = o.customer_id)
              ->  Sort  (cost=21.23..21.24 rows=3 width=7) (actual time=0.334..0.335 rows=3 loops=1)
                    Sort Key: c.id
                    Sort Method: quicksort  Memory: 25kB
                    ->  Seq Scan on customers c  (cost=21.05..21.21 rows=3 width=7) (actual time=0.323..0.326 rows=3 loops=1)
              ->  Sort  (cost=21.44..21.45 rows=5 width=4) (actual time=0.239..0.240 rows=5 loops=1)
                    Sort Key: o.customer_id
                    Sort Method: quicksort  Memory: 25kB
                    ->  Seq Scan on orders o  (cost=21.05..21.38 rows=5 width=4) (actual time=0.224..0.226 rows=5 loops=1)
                          Storage Filter: (amount > 100)
Planning Time: 0.137 ms
Execution Time: 0.635 ms
Peak Memory Usage: 242 kB

```

### s11_optional_on

Rows: 4; serialized decoded-row JSON bytes: 73 (not wire bytes). Planning/execution: 0.122000 / 0.607000 ms.

```text
Hash Right Join  (cost=42.30..42.72 rows=5 width=15) (actual time=0.560..0.568 rows=4 loops=1)
  Hash Cond: (o.customer_id = c.id)
  ->  Seq Scan on orders o  (cost=21.05..21.42 rows=5 width=12) (actual time=0.247..0.249 rows=5 loops=1)
        Storage Filter: (amount > 100)
  ->  Hash  (cost=21.21..21.21 rows=3 width=7) (actual time=0.306..0.306 rows=3 loops=1)
        Buckets: 1024  Batches: 1  Memory Usage: 9kB
        ->  Seq Scan on customers c  (cost=21.05..21.21 rows=3 width=7) (actual time=0.296..0.298 rows=3 loops=1)
Planning Time: 0.122 ms
Execution Time: 0.607 ms
Peak Memory Usage: 191 kB

```

### s16_page_offset

Rows: 2; serialized decoded-row JSON bytes: 23 (not wire bytes). Planning/execution: 0.029000 / 0.259000 ms.

```text
Limit  (cost=21.45..21.46 rows=2 width=8) (actual time=0.224..0.225 rows=2 loops=1)
  ->  Sort  (cost=21.45..21.47 rows=6 width=8) (actual time=0.222..0.223 rows=2 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 25kB
        ->  Seq Scan on orders  (cost=21.05..21.39 rows=6 width=8) (actual time=0.208..0.211 rows=6 loops=1)
Planning Time: 0.029 ms
Execution Time: 0.259 ms
Peak Memory Usage: 92 kB

```

### s17_page_keyset

Rows: 2; serialized decoded-row JSON bytes: 23 (not wire bytes). Planning/execution: 0.025000 / 0.218000 ms.

```text
Limit  (cost=21.45..21.46 rows=2 width=8) (actual time=0.189..0.190 rows=2 loops=1)
  ->  Sort  (cost=21.45..21.47 rows=6 width=8) (actual time=0.188..0.188 rows=2 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 25kB
        ->  Seq Scan on orders  (cost=21.05..21.39 rows=6 width=8) (actual time=0.175..0.178 rows=6 loops=1)
Planning Time: 0.025 ms
Execution Time: 0.218 ms
Peak Memory Usage: 92 kB

```

### s14_topn_rownumber

Rows: 5; serialized decoded-row JSON bytes: 67 (not wire bytes). Planning/execution: 0.070000 / 0.453000 ms.

```text
Subquery Scan on t  (cost=21.46..21.62 rows=5 width=12) (actual time=0.396..0.409 rows=5 loops=1)
  ->  WindowAgg  (cost=21.46..21.57 rows=5 width=52) (actual time=0.395..0.406 rows=5 loops=1)
        Run Condition: (row_number() OVER (?) <= 3)
        ->  Sort  (cost=21.46..21.47 rows=5 width=12) (actual time=0.383..0.384 rows=5 loops=1)
              Sort Key: o.customer_id, o.amount DESC NULLS LAST, o.id
              Sort Method: quicksort  Memory: 25kB
              ->  Seq Scan on orders o  (cost=21.05..21.40 rows=5 width=12) (actual time=0.368..0.371 rows=5 loops=1)
                    Storage Filter: (customer_id IS NOT NULL)
Planning Time: 0.070 ms
Execution Time: 0.453 ms
Peak Memory Usage: 198 kB

```

### s15_topn_lateral

Rows: 5; serialized decoded-row JSON bytes: 67 (not wire bytes). Planning/execution: 0.104000 / 0.889000 ms.

```text
Sort  (cost=85.55..85.56 rows=6 width=12) (actual time=0.845..0.846 rows=5 loops=1)
  Sort Key: c.id, o_1.amount DESC NULLS LAST, o_1.id
  Sort Method: quicksort  Memory: 25kB
  ->  Nested Loop  (cost=42.43..85.47 rows=6 width=12) (actual time=0.437..0.836 rows=5 loops=1)
        ->  Seq Scan on customers c  (cost=21.05..21.18 rows=3 width=4) (actual time=0.225..0.227 rows=3 loops=1)
        ->  Limit  (cost=21.39..21.39 rows=2 width=8) (actual time=0.200..0.201 rows=2 loops=3)
              ->  Sort  (cost=21.39..21.39 rows=2 width=8) (actual time=0.199..0.199 rows=2 loops=3)
                    Sort Key: o_1.amount DESC NULLS LAST, o_1.id
                    Sort Method: quicksort  Memory: 25kB
                    ->  Seq Scan on orders o_1  (cost=21.05..21.38 rows=2 width=8) (actual time=0.189..0.190 rows=2 loops=3)
                          Storage Filter: (customer_id = c.id)
Planning Time: 0.104 ms
Execution Time: 0.889 ms
Peak Memory Usage: 212 kB

```

### Not comparable (never ranked)

- - `s09_not_exists` vs `s12_aggregate_join`: not comparable — different semantic groups answer different questions
- - `s09_not_exists` vs `s13_aggregate_preagg`: not comparable — different semantic groups answer different questions
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
- - `s09_not_exists` vs `s14_topn_rownumber`: not comparable — different semantic groups answer different questions
- - `s09_not_exists` vs `s15_topn_lateral`: not comparable — different semantic groups answer different questions
- - `s10_left_anti` vs `s12_aggregate_join`: not comparable — different semantic groups answer different questions
- - `s10_left_anti` vs `s13_aggregate_preagg`: not comparable — different semantic groups answer different questions
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
- - `s10_left_anti` vs `s14_topn_rownumber`: not comparable — different semantic groups answer different questions
- - `s10_left_anti` vs `s15_topn_lateral`: not comparable — different semantic groups answer different questions
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
- - `s12_aggregate_join` vs `s14_topn_rownumber`: not comparable — different semantic groups answer different questions
- - `s12_aggregate_join` vs `s15_topn_lateral`: not comparable — different semantic groups answer different questions
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
- - `s13_aggregate_preagg` vs `s14_topn_rownumber`: not comparable — different semantic groups answer different questions
- - `s13_aggregate_preagg` vs `s15_topn_lateral`: not comparable — different semantic groups answer different questions
- - `s01_join_where` vs `s06_exists_parent`: not comparable — different semantic groups answer different questions
- - `s01_join_where` vs `s07_in_parent`: not comparable — different semantic groups answer different questions
- - `s01_join_where` vs `s08_join_distinct_parent`: not comparable — different semantic groups answer different questions
- - `s01_join_where` vs `s11_optional_on`: not comparable — different semantic groups answer different questions
- - `s01_join_where` vs `s16_page_offset`: not comparable — different semantic groups answer different questions
- - `s01_join_where` vs `s17_page_keyset`: not comparable — different semantic groups answer different questions
- - `s01_join_where` vs `s14_topn_rownumber`: not comparable — different semantic groups answer different questions
- - `s01_join_where` vs `s15_topn_lateral`: not comparable — different semantic groups answer different questions
- - `s02_filtered_derived` vs `s06_exists_parent`: not comparable — different semantic groups answer different questions
- - `s02_filtered_derived` vs `s07_in_parent`: not comparable — different semantic groups answer different questions
- - `s02_filtered_derived` vs `s08_join_distinct_parent`: not comparable — different semantic groups answer different questions
- - `s02_filtered_derived` vs `s11_optional_on`: not comparable — different semantic groups answer different questions
- - `s02_filtered_derived` vs `s16_page_offset`: not comparable — different semantic groups answer different questions
- - `s02_filtered_derived` vs `s17_page_keyset`: not comparable — different semantic groups answer different questions
- - `s02_filtered_derived` vs `s14_topn_rownumber`: not comparable — different semantic groups answer different questions
- - `s02_filtered_derived` vs `s15_topn_lateral`: not comparable — different semantic groups answer different questions
- - `s03_cte_filtered` vs `s06_exists_parent`: not comparable — different semantic groups answer different questions
- - `s03_cte_filtered` vs `s07_in_parent`: not comparable — different semantic groups answer different questions
- - `s03_cte_filtered` vs `s08_join_distinct_parent`: not comparable — different semantic groups answer different questions
- - `s03_cte_filtered` vs `s11_optional_on`: not comparable — different semantic groups answer different questions
- - `s03_cte_filtered` vs `s16_page_offset`: not comparable — different semantic groups answer different questions
- - `s03_cte_filtered` vs `s17_page_keyset`: not comparable — different semantic groups answer different questions
- - `s03_cte_filtered` vs `s14_topn_rownumber`: not comparable — different semantic groups answer different questions
- - `s03_cte_filtered` vs `s15_topn_lateral`: not comparable — different semantic groups answer different questions
- - `s04_join_on` vs `s06_exists_parent`: not comparable — different semantic groups answer different questions
- - `s04_join_on` vs `s07_in_parent`: not comparable — different semantic groups answer different questions
- - `s04_join_on` vs `s08_join_distinct_parent`: not comparable — different semantic groups answer different questions
- - `s04_join_on` vs `s11_optional_on`: not comparable — different semantic groups answer different questions
- - `s04_join_on` vs `s16_page_offset`: not comparable — different semantic groups answer different questions
- - `s04_join_on` vs `s17_page_keyset`: not comparable — different semantic groups answer different questions
- - `s04_join_on` vs `s14_topn_rownumber`: not comparable — different semantic groups answer different questions
- - `s04_join_on` vs `s15_topn_lateral`: not comparable — different semantic groups answer different questions
- - `s05_comma_join` vs `s06_exists_parent`: not comparable — different semantic groups answer different questions
- - `s05_comma_join` vs `s07_in_parent`: not comparable — different semantic groups answer different questions
- - `s05_comma_join` vs `s08_join_distinct_parent`: not comparable — different semantic groups answer different questions
- - `s05_comma_join` vs `s11_optional_on`: not comparable — different semantic groups answer different questions
- - `s05_comma_join` vs `s16_page_offset`: not comparable — different semantic groups answer different questions
- - `s05_comma_join` vs `s17_page_keyset`: not comparable — different semantic groups answer different questions
- - `s05_comma_join` vs `s14_topn_rownumber`: not comparable — different semantic groups answer different questions
- - `s05_comma_join` vs `s15_topn_lateral`: not comparable — different semantic groups answer different questions
- - `s06_exists_parent` vs `s11_optional_on`: not comparable — different semantic groups answer different questions
- - `s06_exists_parent` vs `s16_page_offset`: not comparable — different semantic groups answer different questions
- - `s06_exists_parent` vs `s17_page_keyset`: not comparable — different semantic groups answer different questions
- - `s06_exists_parent` vs `s14_topn_rownumber`: not comparable — different semantic groups answer different questions
- - `s06_exists_parent` vs `s15_topn_lateral`: not comparable — different semantic groups answer different questions
- - `s07_in_parent` vs `s11_optional_on`: not comparable — different semantic groups answer different questions
- - `s07_in_parent` vs `s16_page_offset`: not comparable — different semantic groups answer different questions
- - `s07_in_parent` vs `s17_page_keyset`: not comparable — different semantic groups answer different questions
- - `s07_in_parent` vs `s14_topn_rownumber`: not comparable — different semantic groups answer different questions
- - `s07_in_parent` vs `s15_topn_lateral`: not comparable — different semantic groups answer different questions
- - `s08_join_distinct_parent` vs `s11_optional_on`: not comparable — different semantic groups answer different questions
- - `s08_join_distinct_parent` vs `s16_page_offset`: not comparable — different semantic groups answer different questions
- - `s08_join_distinct_parent` vs `s17_page_keyset`: not comparable — different semantic groups answer different questions
- - `s08_join_distinct_parent` vs `s14_topn_rownumber`: not comparable — different semantic groups answer different questions
- - `s08_join_distinct_parent` vs `s15_topn_lateral`: not comparable — different semantic groups answer different questions
- - `s11_optional_on` vs `s16_page_offset`: not comparable — different semantic groups answer different questions
- - `s11_optional_on` vs `s17_page_keyset`: not comparable — different semantic groups answer different questions
- - `s11_optional_on` vs `s14_topn_rownumber`: not comparable — different semantic groups answer different questions
- - `s11_optional_on` vs `s15_topn_lateral`: not comparable — different semantic groups answer different questions
- - `s16_page_offset` vs `s14_topn_rownumber`: not comparable — different semantic groups answer different questions
- - `s16_page_offset` vs `s15_topn_lateral`: not comparable — different semantic groups answer different questions
- - `s17_page_keyset` vs `s14_topn_rownumber`: not comparable — different semantic groups answer different questions
- - `s17_page_keyset` vs `s15_topn_lateral`: not comparable — different semantic groups answer different questions

## Predictions this run informs

- P1-P5 join/CTE spellings
- P6-P7 membership/absence
- P8-P9 aggregates/top-N
- P10-P11 predicates/pagination

Verdicts are assigned in the signed analysis, not here.
