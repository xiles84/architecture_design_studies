# SELECT comparison (generated measurements)

_This file contains measurements only. Interpretation belongs in a signed analysis._

**TL;DR:** a validation run supplies correctness and plan evidence. Within-run samples do not establish a replicated performance winner. Missing measurements are N/A.

## Provenance

| Field | Value |
|---|---|
| repo commit | `4f54d848697e50d93c528e9be37dfa06ea340258` |
| inputs digest (SHA-256 of result JSON bytes) | `6b2eabd76ca00c115ea0e6c148fb3e21e8021af9e86d8874361fff2582c7913d` |
| benchmark image | `localhost/charitybench-select:select-yb-check-20261003-final` |
| image digest | `18af3b106a1fcd39996e12ad5596eda4dd822c6c830995e63420dc879d9dfc83` |
| run tag | `run/study-01/select-yb-check-20261003-final` |
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
| select-yb-check-20261003-final | single | s09_not_exists | sel_absence | valid | 0.519004 | 21.1 | 3 | `4befbd47cba4cb7d` |
| select-yb-check-20261003-final | single | s10_left_anti | sel_absence | valid | 0.627754 | 16.0 | 3 | `c46cbe8221360b15` |
| select-yb-check-20261003-final | single | s12_aggregate_join | sel_aggregate | valid | 0.683308 | 27.3 | 3 | `e4988e033c874e00` |
| select-yb-check-20261003-final | single | s13_aggregate_preagg | sel_aggregate | valid | 0.535731 | 19.1 | 3 | `1b1b6b64788bee13` |
| select-yb-check-20261003-final | single | s01_join_where | sel_full_rows | valid | 0.546526 | 8.7 | 3 | `c44ec06f089f0d1c` |
| select-yb-check-20261003-final | single | s02_filtered_derived | sel_full_rows | valid | 0.567798 | 24.1 | 3 | `c44ec06f089f0d1c` |
| select-yb-check-20261003-final | single | s03_cte_filtered | sel_full_rows | valid | 0.835008 | 78.2 | 3 | `c44ec06f089f0d1c` |
| select-yb-check-20261003-final | single | s04_join_on | sel_full_rows | valid | 0.542680 | 20.9 | 3 | `c44ec06f089f0d1c` |
| select-yb-check-20261003-final | single | s05_comma_join | sel_full_rows | valid | 0.933422 | 61.2 | 3 | `c44ec06f089f0d1c` |
| select-yb-check-20261003-final | single | s06_exists_parent | sel_membership | valid | 0.589942 | 26.6 | 3 | `7829bec10abd7f80` |
| select-yb-check-20261003-final | single | s07_in_parent | sel_membership | valid | 0.601584 | 57.3 | 3 | `7829bec10abd7f80` |
| select-yb-check-20261003-final | single | s08_join_distinct_parent | sel_membership | valid | 0.579988 | 3.8 | 3 | `5a4de51c305ef1b0` |
| select-yb-check-20261003-final | single | s11_optional_on | sel_optional_match | valid | 0.552582 | 9.6 | 3 | `c46cbe8221360b15` |
| select-yb-check-20261003-final | single | s16_page_offset | sel_pagination | valid | 0.297046 | 28.0 | 3 | `ac4509d295bf3348` |
| select-yb-check-20261003-final | single | s17_page_keyset | sel_pagination | valid | 0.417142 | 75.8 | 3 | `ac4509d295bf3348` |
| select-yb-check-20261003-final | single | s14_topn_rownumber | sel_topn | valid | 0.381266 | 10.5 | 3 | `e2ad05f4f9b57fe5` |
| select-yb-check-20261003-final | single | s15_topn_lateral | sel_topn | valid | 1.055520 | 24.0 | 3 | `0504ce42ca4d7a56` |

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

- `s09_not_exists` vs `s10_left_anti` on single: medians differ by 0.109 ms, inside the 21.1%/16.0% spread — no resolved difference
- `s12_aggregate_join` vs `s13_aggregate_preagg` on single: medians differ by 0.148 ms, inside the 27.3%/19.1% spread — no resolved difference
- `s06_exists_parent` vs `s08_join_distinct_parent` on single: medians differ by 0.010 ms, inside the 26.6%/3.8% spread — no resolved difference
- `s07_in_parent` vs `s08_join_distinct_parent` on single: medians differ by 0.022 ms, inside the 57.3%/3.8% spread — no resolved difference

### Within-run differences (not replicated)

- `s14_topn_rownumber` vs `s15_topn_lateral` on single: within-run medians 0.381266 / 1.055520 ms; whole-run replication required — no resolved difference

## Live negative controls

- `x01_optional_where`: detected=true, rows=3
- `x02_not_in_nullable`: detected=true, rows=0
- `x03_distinct_on_name`: detected=true, rows=1
- `x04_limit_before_join`: detected=true, rows=2

## Captured plans and result sizes

### s09_not_exists

Rows: 0; serialized decoded-row JSON bytes: 4 (not wire bytes). Planning/execution: 0.058000 / 0.391000 ms.

```text
Hash Anti Join  (cost=42.49..42.67 rows=1 width=7) (actual time=0.365..0.366 rows=0 loops=1)
  Hash Cond: (c.id = o.customer_id)
  ->  Seq Scan on customers c  (cost=21.05..21.21 rows=3 width=7) (actual time=0.164..0.165 rows=3 loops=1)
  ->  Hash  (cost=21.38..21.38 rows=5 width=4) (actual time=0.186..0.186 rows=4 loops=1)
        Buckets: 1024  Batches: 1  Memory Usage: 9kB
        ->  Seq Scan on orders o  (cost=21.05..21.38 rows=5 width=4) (actual time=0.178..0.179 rows=5 loops=1)
              Storage Filter: (amount > 100)
Planning Time: 0.058 ms
Execution Time: 0.391 ms
Peak Memory Usage: 205 kB

```

### s10_left_anti

Rows: 0; serialized decoded-row JSON bytes: 4 (not wire bytes). Planning/execution: 0.079000 / 0.529000 ms.

```text
Hash Right Join  (cost=42.30..42.70 rows=1 width=7) (actual time=0.495..0.496 rows=0 loops=1)
  Hash Cond: (o.customer_id = c.id)
  Filter: (o.id IS NULL)
  Rows Removed by Filter: 4
  ->  Seq Scan on orders o  (cost=21.05..21.40 rows=5 width=8) (actual time=0.220..0.221 rows=5 loops=1)
        Storage Filter: (amount > 100)
  ->  Hash  (cost=21.21..21.21 rows=3 width=7) (actual time=0.259..0.259 rows=3 loops=1)
        Buckets: 1024  Batches: 1  Memory Usage: 9kB
        ->  Seq Scan on customers c  (cost=21.05..21.21 rows=3 width=7) (actual time=0.248..0.249 rows=3 loops=1)
Planning Time: 0.079 ms
Execution Time: 0.529 ms
Peak Memory Usage: 191 kB

```

### s12_aggregate_join

Rows: 3; serialized decoded-row JSON bytes: 46 (not wire bytes). Planning/execution: 0.073000 / 0.302000 ms.

```text
GroupAggregate  (cost=42.72..42.89 rows=3 width=23) (actual time=0.270..0.273 rows=3 loops=1)
  Group Key: c.id
  ->  Merge Left Join  (cost=42.72..42.82 rows=6 width=15) (actual time=0.267..0.269 rows=5 loops=1)
        Merge Cond: (c.id = o.customer_id)
        ->  Sort  (cost=21.23..21.24 rows=3 width=7) (actual time=0.149..0.149 rows=3 loops=1)
              Sort Key: c.id
              Sort Method: quicksort  Memory: 25kB
              ->  Seq Scan on customers c  (cost=21.05..21.21 rows=3 width=7) (actual time=0.141..0.143 rows=3 loops=1)
        ->  Sort  (cost=21.49..21.51 rows=6 width=12) (actual time=0.116..0.116 rows=6 loops=1)
              Sort Key: o.customer_id
              Sort Method: quicksort  Memory: 25kB
              ->  Seq Scan on orders o  (cost=21.05..21.41 rows=6 width=12) (actual time=0.107..0.108 rows=6 loops=1)
Planning Time: 0.073 ms
Execution Time: 0.302 ms
Peak Memory Usage: 206 kB

```

### s13_aggregate_preagg

Rows: 3; serialized decoded-row JSON bytes: 46 (not wire bytes). Planning/execution: 0.070000 / 0.415000 ms.

```text
Hash Left Join  (cost=42.67..42.85 rows=3 width=23) (actual time=0.368..0.370 rows=3 loops=1)
  Hash Cond: (c.id = a.customer_id)
  ->  Seq Scan on customers c  (cost=21.05..21.21 rows=3 width=7) (actual time=0.183..0.184 rows=3 loops=1)
  ->  Hash  (cost=21.59..21.59 rows=3 width=20) (actual time=0.177..0.178 rows=3 loops=1)
        Buckets: 1024  Batches: 1  Memory Usage: 9kB
        ->  Subquery Scan on a  (cost=21.47..21.59 rows=3 width=20) (actual time=0.169..0.172 rows=4 loops=1)
              ->  GroupAggregate  (cost=21.47..21.56 rows=3 width=20) (actual time=0.168..0.171 rows=4 loops=1)
                    Group Key: orders.customer_id
                    ->  Sort  (cost=21.47..21.48 rows=6 width=8) (actual time=0.165..0.166 rows=6 loops=1)
                          Sort Key: orders.customer_id
                          Sort Method: quicksort  Memory: 25kB
                          ->  Seq Scan on orders  (cost=21.05..21.39 rows=6 width=8) (actual time=0.156..0.158 rows=6 loops=1)
Planning Time: 0.070 ms
Execution Time: 0.415 ms
Peak Memory Usage: 232 kB

```

### s01_join_where

Rows: 4; serialized decoded-row JSON bytes: 73 (not wire bytes). Planning/execution: 0.112000 / 0.405000 ms.

```text
Hash Join  (cost=42.30..42.72 rows=5 width=15) (actual time=0.373..0.375 rows=4 loops=1)
  Hash Cond: (o.customer_id = c.id)
  ->  Seq Scan on orders o  (cost=21.05..21.42 rows=5 width=12) (actual time=0.210..0.211 rows=5 loops=1)
        Storage Filter: (amount > 100)
  ->  Hash  (cost=21.21..21.21 rows=3 width=7) (actual time=0.153..0.153 rows=3 loops=1)
        Buckets: 1024  Batches: 1  Memory Usage: 9kB
        ->  Seq Scan on customers c  (cost=21.05..21.21 rows=3 width=7) (actual time=0.146..0.147 rows=3 loops=1)
Planning Time: 0.112 ms
Execution Time: 0.405 ms
Peak Memory Usage: 205 kB

```

### s02_filtered_derived

Rows: 4; serialized decoded-row JSON bytes: 73 (not wire bytes). Planning/execution: 0.089000 / 0.360000 ms.

```text
Hash Join  (cost=42.30..42.72 rows=5 width=15) (actual time=0.331..0.334 rows=4 loops=1)
  Hash Cond: (orders.customer_id = c.id)
  ->  Seq Scan on orders  (cost=21.05..21.42 rows=5 width=12) (actual time=0.171..0.172 rows=5 loops=1)
        Storage Filter: (amount > 100)
  ->  Hash  (cost=21.21..21.21 rows=3 width=7) (actual time=0.151..0.151 rows=3 loops=1)
        Buckets: 1024  Batches: 1  Memory Usage: 9kB
        ->  Seq Scan on customers c  (cost=21.05..21.21 rows=3 width=7) (actual time=0.135..0.137 rows=3 loops=1)
Planning Time: 0.089 ms
Execution Time: 0.360 ms
Peak Memory Usage: 205 kB

```

### s03_cte_filtered

Rows: 4; serialized decoded-row JSON bytes: 73 (not wire bytes). Planning/execution: 0.099000 / 0.354000 ms.

```text
Hash Join  (cost=42.30..42.72 rows=5 width=15) (actual time=0.325..0.327 rows=4 loops=1)
  Hash Cond: (orders.customer_id = c.id)
  ->  Seq Scan on orders  (cost=21.05..21.42 rows=5 width=12) (actual time=0.171..0.172 rows=5 loops=1)
        Storage Filter: (amount > 100)
  ->  Hash  (cost=21.21..21.21 rows=3 width=7) (actual time=0.146..0.146 rows=3 loops=1)
        Buckets: 1024  Batches: 1  Memory Usage: 9kB
        ->  Seq Scan on customers c  (cost=21.05..21.21 rows=3 width=7) (actual time=0.130..0.131 rows=3 loops=1)
Planning Time: 0.099 ms
Execution Time: 0.354 ms
Peak Memory Usage: 205 kB

```

### s04_join_on

Rows: 4; serialized decoded-row JSON bytes: 73 (not wire bytes). Planning/execution: 0.091000 / 0.702000 ms.

```text
Hash Join  (cost=42.30..42.72 rows=5 width=15) (actual time=0.357..0.360 rows=4 loops=1)
  Hash Cond: (o.customer_id = c.id)
  ->  Seq Scan on orders o  (cost=21.05..21.42 rows=5 width=12) (actual time=0.176..0.177 rows=5 loops=1)
        Storage Filter: (amount > 100)
  ->  Hash  (cost=21.21..21.21 rows=3 width=7) (actual time=0.172..0.173 rows=3 loops=1)
        Buckets: 1024  Batches: 1  Memory Usage: 9kB
        ->  Seq Scan on customers c  (cost=21.05..21.21 rows=3 width=7) (actual time=0.166..0.167 rows=3 loops=1)
Planning Time: 0.091 ms
Execution Time: 0.702 ms
Peak Memory Usage: 205 kB

```

### s05_comma_join

Rows: 4; serialized decoded-row JSON bytes: 73 (not wire bytes). Planning/execution: 0.084000 / 0.335000 ms.

```text
Hash Join  (cost=42.30..42.72 rows=5 width=15) (actual time=0.305..0.307 rows=4 loops=1)
  Hash Cond: (o.customer_id = c.id)
  ->  Seq Scan on orders o  (cost=21.05..21.42 rows=5 width=12) (actual time=0.160..0.161 rows=5 loops=1)
        Storage Filter: (amount > 100)
  ->  Hash  (cost=21.21..21.21 rows=3 width=7) (actual time=0.136..0.136 rows=3 loops=1)
        Buckets: 1024  Batches: 1  Memory Usage: 9kB
        ->  Seq Scan on customers c  (cost=21.05..21.21 rows=3 width=7) (actual time=0.128..0.130 rows=3 loops=1)
Planning Time: 0.084 ms
Execution Time: 0.335 ms
Peak Memory Usage: 205 kB

```

### s06_exists_parent

Rows: 3; serialized decoded-row JSON bytes: 28 (not wire bytes). Planning/execution: 0.102000 / 0.413000 ms.

```text
Hash Semi Join  (cost=42.49..42.70 rows=3 width=7) (actual time=0.382..0.384 rows=3 loops=1)
  Hash Cond: (c.id = o.customer_id)
  ->  Seq Scan on customers c  (cost=21.05..21.21 rows=3 width=7) (actual time=0.207..0.208 rows=3 loops=1)
  ->  Hash  (cost=21.38..21.38 rows=5 width=4) (actual time=0.166..0.167 rows=4 loops=1)
        Buckets: 1024  Batches: 1  Memory Usage: 9kB
        ->  Seq Scan on orders o  (cost=21.05..21.38 rows=5 width=4) (actual time=0.159..0.161 rows=5 loops=1)
              Storage Filter: (amount > 100)
Planning Time: 0.102 ms
Execution Time: 0.413 ms
Peak Memory Usage: 205 kB

```

### s07_in_parent

Rows: 3; serialized decoded-row JSON bytes: 28 (not wire bytes). Planning/execution: 0.091000 / 0.324000 ms.

```text
Hash Semi Join  (cost=42.49..42.70 rows=3 width=7) (actual time=0.295..0.297 rows=3 loops=1)
  Hash Cond: (c.id = o.customer_id)
  ->  Seq Scan on customers c  (cost=21.05..21.21 rows=3 width=7) (actual time=0.155..0.156 rows=3 loops=1)
  ->  Hash  (cost=21.38..21.38 rows=5 width=4) (actual time=0.133..0.133 rows=4 loops=1)
        Buckets: 1024  Batches: 1  Memory Usage: 9kB
        ->  Seq Scan on orders o  (cost=21.05..21.38 rows=5 width=4) (actual time=0.128..0.129 rows=5 loops=1)
              Storage Filter: (amount > 100)
Planning Time: 0.091 ms
Execution Time: 0.324 ms
Peak Memory Usage: 205 kB

```

### s08_join_distinct_parent

Rows: 3; serialized decoded-row JSON bytes: 28 (not wire bytes). Planning/execution: 0.099000 / 0.362000 ms.

```text
Unique  (cost=42.81..42.84 rows=3 width=7) (actual time=0.331..0.333 rows=3 loops=1)
  ->  Sort  (cost=42.81..42.82 rows=5 width=7) (actual time=0.330..0.331 rows=4 loops=1)
        Sort Key: c.id, c.name
        Sort Method: quicksort  Memory: 25kB
        ->  Merge Join  (cost=42.67..42.75 rows=5 width=7) (actual time=0.322..0.325 rows=4 loops=1)
              Merge Cond: (c.id = o.customer_id)
              ->  Sort  (cost=21.23..21.24 rows=3 width=7) (actual time=0.173..0.174 rows=3 loops=1)
                    Sort Key: c.id
                    Sort Method: quicksort  Memory: 25kB
                    ->  Seq Scan on customers c  (cost=21.05..21.21 rows=3 width=7) (actual time=0.166..0.167 rows=3 loops=1)
              ->  Sort  (cost=21.44..21.45 rows=5 width=4) (actual time=0.146..0.146 rows=5 loops=1)
                    Sort Key: o.customer_id
                    Sort Method: quicksort  Memory: 25kB
                    ->  Seq Scan on orders o  (cost=21.05..21.38 rows=5 width=4) (actual time=0.137..0.138 rows=5 loops=1)
                          Storage Filter: (amount > 100)
Planning Time: 0.099 ms
Execution Time: 0.362 ms
Peak Memory Usage: 242 kB

```

### s11_optional_on

Rows: 4; serialized decoded-row JSON bytes: 73 (not wire bytes). Planning/execution: 0.059000 / 0.381000 ms.

```text
Hash Right Join  (cost=42.30..42.72 rows=5 width=15) (actual time=0.348..0.355 rows=4 loops=1)
  Hash Cond: (o.customer_id = c.id)
  ->  Seq Scan on orders o  (cost=21.05..21.42 rows=5 width=12) (actual time=0.171..0.173 rows=5 loops=1)
        Storage Filter: (amount > 100)
  ->  Hash  (cost=21.21..21.21 rows=3 width=7) (actual time=0.171..0.171 rows=3 loops=1)
        Buckets: 1024  Batches: 1  Memory Usage: 9kB
        ->  Seq Scan on customers c  (cost=21.05..21.21 rows=3 width=7) (actual time=0.157..0.159 rows=3 loops=1)
Planning Time: 0.059 ms
Execution Time: 0.381 ms
Peak Memory Usage: 191 kB

```

### s16_page_offset

Rows: 2; serialized decoded-row JSON bytes: 23 (not wire bytes). Planning/execution: 0.021000 / 0.170000 ms.

```text
Limit  (cost=21.45..21.46 rows=2 width=8) (actual time=0.151..0.152 rows=2 loops=1)
  ->  Sort  (cost=21.45..21.47 rows=6 width=8) (actual time=0.150..0.151 rows=2 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 25kB
        ->  Seq Scan on orders  (cost=21.05..21.39 rows=6 width=8) (actual time=0.144..0.146 rows=6 loops=1)
Planning Time: 0.021 ms
Execution Time: 0.170 ms
Peak Memory Usage: 92 kB

```

### s17_page_keyset

Rows: 2; serialized decoded-row JSON bytes: 23 (not wire bytes). Planning/execution: 0.026000 / 0.197000 ms.

```text
Limit  (cost=21.45..21.46 rows=2 width=8) (actual time=0.172..0.173 rows=2 loops=1)
  ->  Sort  (cost=21.45..21.47 rows=6 width=8) (actual time=0.171..0.171 rows=2 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 25kB
        ->  Seq Scan on orders  (cost=21.05..21.39 rows=6 width=8) (actual time=0.160..0.163 rows=6 loops=1)
Planning Time: 0.026 ms
Execution Time: 0.197 ms
Peak Memory Usage: 92 kB

```

### s14_topn_rownumber

Rows: 5; serialized decoded-row JSON bytes: 67 (not wire bytes). Planning/execution: 0.055000 / 0.218000 ms.

```text
Subquery Scan on t  (cost=21.46..21.62 rows=5 width=12) (actual time=0.186..0.193 rows=5 loops=1)
  ->  WindowAgg  (cost=21.46..21.57 rows=5 width=52) (actual time=0.186..0.192 rows=5 loops=1)
        Run Condition: (row_number() OVER (?) <= 3)
        ->  Sort  (cost=21.46..21.47 rows=5 width=12) (actual time=0.179..0.179 rows=5 loops=1)
              Sort Key: o.customer_id, o.amount DESC NULLS LAST, o.id
              Sort Method: quicksort  Memory: 25kB
              ->  Seq Scan on orders o  (cost=21.05..21.40 rows=5 width=12) (actual time=0.170..0.172 rows=5 loops=1)
                    Storage Filter: (customer_id IS NOT NULL)
Planning Time: 0.055 ms
Execution Time: 0.218 ms
Peak Memory Usage: 198 kB

```

### s15_topn_lateral

Rows: 5; serialized decoded-row JSON bytes: 67 (not wire bytes). Planning/execution: 0.087000 / 0.873000 ms.

```text
Sort  (cost=85.55..85.56 rows=6 width=12) (actual time=0.842..0.843 rows=5 loops=1)
  Sort Key: c.id, o_1.amount DESC NULLS LAST, o_1.id
  Sort Method: quicksort  Memory: 25kB
  ->  Nested Loop  (cost=42.43..85.47 rows=6 width=12) (actual time=0.497..0.836 rows=5 loops=1)
        ->  Seq Scan on customers c  (cost=21.05..21.18 rows=3 width=4) (actual time=0.230..0.232 rows=3 loops=1)
        ->  Limit  (cost=21.39..21.39 rows=2 width=8) (actual time=0.198..0.199 rows=2 loops=3)
              ->  Sort  (cost=21.39..21.39 rows=2 width=8) (actual time=0.197..0.197 rows=2 loops=3)
                    Sort Key: o_1.amount DESC NULLS LAST, o_1.id
                    Sort Method: quicksort  Memory: 25kB
                    ->  Seq Scan on orders o_1  (cost=21.05..21.38 rows=2 width=8) (actual time=0.186..0.187 rows=2 loops=3)
                          Storage Filter: (customer_id = c.id)
Planning Time: 0.087 ms
Execution Time: 0.873 ms
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
