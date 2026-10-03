# SELECT comparison (generated measurements)

_This file contains measurements only. Interpretation belongs in a signed analysis._

**TL;DR:** a validation run supplies correctness and plan evidence. Within-run samples do not establish a replicated performance winner. Missing measurements are N/A.

## Provenance

| Field | Value |
|---|---|
| repo commit | `b4583c1149c67c6ce0123208e3b62c8b88b231d6` |
| inputs digest (SHA-256 of result JSON bytes) | `3180cf4f671d5047281a5ab37373e14faf9da6eee3578d78d2d454fcf115ba22` |
| benchmark image | `localhost/charitybench-select:select-pg-check-20261003-final` |
| image digest | `18af3b106a1fcd39996e12ad5596eda4dd822c6c830995e63420dc879d9dfc83` |
| run tag | `run/study-01/select-pg-check-20261003-final` |
| environment | `xiles-podman-select-validation-20261003` |
| engine image digest | `7296f210ae81031ec955dbad9a67a84fe958572a2153b8d0826a647522904dc1` |
| engine version | PostgreSQL 17.11 (Debian 17.11-1.pgdg13+2) on x86_64-pc-linux-gnu, compiled by gcc (Debian 14.2.0-19) 14.2.0, 64-bit |
| dataset seed | 42 |
| dataset fingerprint | `f05b8420d983e1c8` |
| setting: client_workers | 1 sequential reader; 0 writers |
| setting: connection_pool | 2 = 1 reader/setup + 1 reserve |
| setting: distributed_rpc | N/A on PostgreSQL single-node |
| setting: effective_cache_size | 262144 |
| setting: index_set | orders_amount_idx,orders_customer_idx,primary keys |
| setting: max_connections | 200 |
| setting: prepared_plan_policy | pgx simple protocol; no prepared statements |
| setting: random_page_cost | 1.1 |
| setting: scale | tiny |
| setting: shared_buffers | 98304 |
| setting: statement_budget_ms | 120000 |
| setting: statement_timeout | 120000 |
| setting: statistics | ANALYZE after deterministic load |
| setting: topology | single |
| setting: track_io_timing | on |
| setting: whole_run_count | 1; no replicated winner |
| setting: work_mem | 32768 |

## Cells

| run | topology | statement | group | status | median ms | spread % | trials | plan structure |
|---|---|---|---|---|---|---|---|---|
| select-pg-check-20261003-final | single | s09_not_exists | sel_absence | valid | 0.115828 | 20.2 | 3 | `ee53c07e94bc3d0e` |
| select-pg-check-20261003-final | single | s10_left_anti | sel_absence | valid | 0.123656 | 43.0 | 3 | `f16e64291eebe84c` |
| select-pg-check-20261003-final | single | s12_aggregate_join | sel_aggregate | valid | 0.190127 | 34.4 | 3 | `4e6500440fe6ec01` |
| select-pg-check-20261003-final | single | s13_aggregate_preagg | sel_aggregate | valid | 0.164982 | 14.7 | 3 | `d5134ef19610103a` |
| select-pg-check-20261003-final | single | s01_join_where | sel_full_rows | valid | 0.172245 | 23.3 | 3 | `10d644054c377326` |
| select-pg-check-20261003-final | single | s02_filtered_derived | sel_full_rows | valid | 0.193453 | 16.5 | 3 | `10d644054c377326` |
| select-pg-check-20261003-final | single | s03_cte_filtered | sel_full_rows | valid | 0.184152 | 29.5 | 3 | `10d644054c377326` |
| select-pg-check-20261003-final | single | s04_join_on | sel_full_rows | valid | 0.180473 | 22.3 | 3 | `10d644054c377326` |
| select-pg-check-20261003-final | single | s05_comma_join | sel_full_rows | valid | 0.153330 | 28.6 | 3 | `10d644054c377326` |
| select-pg-check-20261003-final | single | s06_exists_parent | sel_membership | valid | 0.149768 | 8.4 | 3 | `01b008bf56170088` |
| select-pg-check-20261003-final | single | s07_in_parent | sel_membership | valid | 0.150218 | 17.6 | 3 | `01b008bf56170088` |
| select-pg-check-20261003-final | single | s08_join_distinct_parent | sel_membership | valid | 0.131152 | 22.7 | 3 | `16eb24166018a1dd` |
| select-pg-check-20261003-final | single | s11_optional_on | sel_optional_match | valid | 0.244851 | 22.3 | 3 | `f16e64291eebe84c` |
| select-pg-check-20261003-final | single | s16_page_offset | sel_pagination | valid | 0.093186 | 48.2 | 3 | `ac4509d295bf3348` |
| select-pg-check-20261003-final | single | s17_page_keyset | sel_pagination | valid | 0.119653 | 40.3 | 3 | `ac4509d295bf3348` |
| select-pg-check-20261003-final | single | s14_topn_rownumber | sel_topn | valid | 0.144106 | 22.8 | 3 | `03a28a8e98c496b2` |
| select-pg-check-20261003-final | single | s15_topn_lateral | sel_topn | valid | 0.161580 | 14.7 | 3 | `49a784cf94a2732c` |

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

- `s09_not_exists` vs `s10_left_anti` on single: medians differ by 0.008 ms, inside the 20.2%/43.0% spread — no resolved difference
- `s12_aggregate_join` vs `s13_aggregate_preagg` on single: medians differ by 0.025 ms, inside the 34.4%/14.7% spread — no resolved difference
- `s06_exists_parent` vs `s08_join_distinct_parent` on single: medians differ by 0.019 ms, inside the 8.4%/22.7% spread — no resolved difference
- `s07_in_parent` vs `s08_join_distinct_parent` on single: medians differ by 0.019 ms, inside the 17.6%/22.7% spread — no resolved difference
- `s14_topn_rownumber` vs `s15_topn_lateral` on single: medians differ by 0.017 ms, inside the 22.8%/14.7% spread — no resolved difference

### Within-run differences (not replicated)

None.

## Live negative controls

- `x01_optional_where`: detected=true, rows=3
- `x02_not_in_nullable`: detected=true, rows=0
- `x03_distinct_on_name`: detected=true, rows=1
- `x04_limit_before_join`: detected=true, rows=2

## Captured plans and result sizes

### s09_not_exists

Rows: 0; serialized decoded-row JSON bytes: 4 (not wire bytes). Planning/execution: 0.032000 / 0.010000 ms.

```text
Hash Anti Join  (cost=1.14..2.19 rows=1 width=7) (actual time=0.005..0.006 rows=0 loops=1)
  Hash Cond: (c.id = o.customer_id)
  Buffers: shared hit=2
  ->  Seq Scan on customers c  (cost=0.00..1.03 rows=3 width=7) (actual time=0.001..0.001 rows=3 loops=1)
        Buffers: shared hit=1
  ->  Hash  (cost=1.07..1.07 rows=5 width=4) (actual time=0.003..0.003 rows=4 loops=1)
        Buckets: 1024  Batches: 1  Memory Usage: 9kB
        Buffers: shared hit=1
        ->  Seq Scan on orders o  (cost=0.00..1.07 rows=5 width=4) (actual time=0.001..0.002 rows=5 loops=1)
              Filter: (amount > 100)
              Rows Removed by Filter: 1
              Buffers: shared hit=1
Planning:
  Buffers: shared hit=2
Planning Time: 0.032 ms
Execution Time: 0.010 ms

```

### s10_left_anti

Rows: 0; serialized decoded-row JSON bytes: 4 (not wire bytes). Planning/execution: 0.039000 / 0.012000 ms.

```text
Hash Right Join  (cost=1.07..2.17 rows=1 width=7) (actual time=0.007..0.008 rows=0 loops=1)
  Hash Cond: (o.customer_id = c.id)
  Filter: (o.id IS NULL)
  Rows Removed by Filter: 4
  Buffers: shared hit=2
  ->  Seq Scan on orders o  (cost=0.00..1.07 rows=5 width=8) (actual time=0.001..0.002 rows=5 loops=1)
        Filter: (amount > 100)
        Rows Removed by Filter: 1
        Buffers: shared hit=1
  ->  Hash  (cost=1.03..1.03 rows=3 width=7) (actual time=0.003..0.003 rows=3 loops=1)
        Buckets: 1024  Batches: 1  Memory Usage: 9kB
        Buffers: shared hit=1
        ->  Seq Scan on customers c  (cost=0.00..1.03 rows=3 width=7) (actual time=0.002..0.002 rows=3 loops=1)
              Buffers: shared hit=1
Planning:
  Buffers: shared hit=2
Planning Time: 0.039 ms
Execution Time: 0.012 ms

```

### s12_aggregate_join

Rows: 3; serialized decoded-row JSON bytes: 46 (not wire bytes). Planning/execution: 0.049000 / 0.035000 ms.

```text
HashAggregate  (cost=2.21..2.24 rows=3 width=23) (actual time=0.024..0.025 rows=3 loops=1)
  Group Key: c.id
  Batches: 1  Memory Usage: 24kB
  Buffers: shared hit=2
  ->  Hash Right Join  (cost=1.07..2.17 rows=6 width=15) (actual time=0.019..0.021 rows=5 loops=1)
        Hash Cond: (o.customer_id = c.id)
        Buffers: shared hit=2
        ->  Seq Scan on orders o  (cost=0.00..1.06 rows=6 width=12) (actual time=0.013..0.014 rows=6 loops=1)
              Buffers: shared hit=1
        ->  Hash  (cost=1.03..1.03 rows=3 width=7) (actual time=0.003..0.004 rows=3 loops=1)
              Buckets: 1024  Batches: 1  Memory Usage: 9kB
              Buffers: shared hit=1
              ->  Seq Scan on customers c  (cost=0.00..1.03 rows=3 width=7) (actual time=0.002..0.002 rows=3 loops=1)
                    Buffers: shared hit=1
Planning:
  Buffers: shared hit=2
Planning Time: 0.049 ms
Execution Time: 0.035 ms

```

### s13_aggregate_preagg

Rows: 3; serialized decoded-row JSON bytes: 46 (not wire bytes). Planning/execution: 0.037000 / 0.020000 ms.

```text
Hash Left Join  (cost=1.20..2.25 rows=3 width=23) (actual time=0.009..0.011 rows=3 loops=1)
  Hash Cond: (c.id = a.customer_id)
  Buffers: shared hit=2
  ->  Seq Scan on customers c  (cost=0.00..1.03 rows=3 width=7) (actual time=0.002..0.002 rows=3 loops=1)
        Buffers: shared hit=1
  ->  Hash  (cost=1.17..1.17 rows=3 width=20) (actual time=0.006..0.006 rows=3 loops=1)
        Buckets: 1024  Batches: 1  Memory Usage: 9kB
        Buffers: shared hit=1
        ->  Subquery Scan on a  (cost=1.10..1.17 rows=3 width=20) (actual time=0.004..0.005 rows=4 loops=1)
              Buffers: shared hit=1
              ->  HashAggregate  (cost=1.10..1.14 rows=3 width=20) (actual time=0.004..0.005 rows=4 loops=1)
                    Group Key: orders.customer_id
                    Batches: 1  Memory Usage: 24kB
                    Buffers: shared hit=1
                    ->  Seq Scan on orders  (cost=0.00..1.06 rows=6 width=8) (actual time=0.001..0.001 rows=6 loops=1)
                          Buffers: shared hit=1
Planning Time: 0.037 ms
Execution Time: 0.020 ms

```

### s01_join_where

Rows: 4; serialized decoded-row JSON bytes: 73 (not wire bytes). Planning/execution: 0.056000 / 0.021000 ms.

```text
Hash Join  (cost=1.07..2.17 rows=5 width=15) (actual time=0.013..0.015 rows=4 loops=1)
  Hash Cond: (o.customer_id = c.id)
  Buffers: shared hit=2
  ->  Seq Scan on orders o  (cost=0.00..1.07 rows=5 width=12) (actual time=0.003..0.004 rows=5 loops=1)
        Filter: (amount > 100)
        Rows Removed by Filter: 1
        Buffers: shared hit=1
  ->  Hash  (cost=1.03..1.03 rows=3 width=7) (actual time=0.007..0.007 rows=3 loops=1)
        Buckets: 1024  Batches: 1  Memory Usage: 9kB
        Buffers: shared hit=1
        ->  Seq Scan on customers c  (cost=0.00..1.03 rows=3 width=7) (actual time=0.001..0.002 rows=3 loops=1)
              Buffers: shared hit=1
Planning:
  Buffers: shared hit=2
Planning Time: 0.056 ms
Execution Time: 0.021 ms

```

### s02_filtered_derived

Rows: 4; serialized decoded-row JSON bytes: 73 (not wire bytes). Planning/execution: 0.046000 / 0.017000 ms.

```text
Hash Join  (cost=1.07..2.17 rows=5 width=15) (actual time=0.010..0.012 rows=4 loops=1)
  Hash Cond: (orders.customer_id = c.id)
  Buffers: shared hit=2
  ->  Seq Scan on orders  (cost=0.00..1.07 rows=5 width=12) (actual time=0.002..0.003 rows=5 loops=1)
        Filter: (amount > 100)
        Rows Removed by Filter: 1
        Buffers: shared hit=1
  ->  Hash  (cost=1.03..1.03 rows=3 width=7) (actual time=0.006..0.007 rows=3 loops=1)
        Buckets: 1024  Batches: 1  Memory Usage: 9kB
        Buffers: shared hit=1
        ->  Seq Scan on customers c  (cost=0.00..1.03 rows=3 width=7) (actual time=0.001..0.001 rows=3 loops=1)
              Buffers: shared hit=1
Planning:
  Buffers: shared hit=2
Planning Time: 0.046 ms
Execution Time: 0.017 ms

```

### s03_cte_filtered

Rows: 4; serialized decoded-row JSON bytes: 73 (not wire bytes). Planning/execution: 0.068000 / 0.015000 ms.

```text
Hash Join  (cost=1.07..2.17 rows=5 width=15) (actual time=0.007..0.009 rows=4 loops=1)
  Hash Cond: (orders.customer_id = c.id)
  Buffers: shared hit=2
  ->  Seq Scan on orders  (cost=0.00..1.07 rows=5 width=12) (actual time=0.002..0.003 rows=5 loops=1)
        Filter: (amount > 100)
        Rows Removed by Filter: 1
        Buffers: shared hit=1
  ->  Hash  (cost=1.03..1.03 rows=3 width=7) (actual time=0.003..0.003 rows=3 loops=1)
        Buckets: 1024  Batches: 1  Memory Usage: 9kB
        Buffers: shared hit=1
        ->  Seq Scan on customers c  (cost=0.00..1.03 rows=3 width=7) (actual time=0.001..0.002 rows=3 loops=1)
              Buffers: shared hit=1
Planning:
  Buffers: shared hit=2
Planning Time: 0.068 ms
Execution Time: 0.015 ms

```

### s04_join_on

Rows: 4; serialized decoded-row JSON bytes: 73 (not wire bytes). Planning/execution: 0.050000 / 0.016000 ms.

```text
Hash Join  (cost=1.07..2.17 rows=5 width=15) (actual time=0.008..0.009 rows=4 loops=1)
  Hash Cond: (o.customer_id = c.id)
  Buffers: shared hit=2
  ->  Seq Scan on orders o  (cost=0.00..1.07 rows=5 width=12) (actual time=0.003..0.004 rows=5 loops=1)
        Filter: (amount > 100)
        Rows Removed by Filter: 1
        Buffers: shared hit=1
  ->  Hash  (cost=1.03..1.03 rows=3 width=7) (actual time=0.003..0.003 rows=3 loops=1)
        Buckets: 1024  Batches: 1  Memory Usage: 9kB
        Buffers: shared hit=1
        ->  Seq Scan on customers c  (cost=0.00..1.03 rows=3 width=7) (actual time=0.001..0.002 rows=3 loops=1)
              Buffers: shared hit=1
Planning:
  Buffers: shared hit=2
Planning Time: 0.050 ms
Execution Time: 0.016 ms

```

### s05_comma_join

Rows: 4; serialized decoded-row JSON bytes: 73 (not wire bytes). Planning/execution: 0.047000 / 0.011000 ms.

```text
Hash Join  (cost=1.07..2.17 rows=5 width=15) (actual time=0.006..0.007 rows=4 loops=1)
  Hash Cond: (o.customer_id = c.id)
  Buffers: shared hit=2
  ->  Seq Scan on orders o  (cost=0.00..1.07 rows=5 width=12) (actual time=0.002..0.003 rows=5 loops=1)
        Filter: (amount > 100)
        Rows Removed by Filter: 1
        Buffers: shared hit=1
  ->  Hash  (cost=1.03..1.03 rows=3 width=7) (actual time=0.002..0.002 rows=3 loops=1)
        Buckets: 1024  Batches: 1  Memory Usage: 9kB
        Buffers: shared hit=1
        ->  Seq Scan on customers c  (cost=0.00..1.03 rows=3 width=7) (actual time=0.001..0.001 rows=3 loops=1)
              Buffers: shared hit=1
Planning:
  Buffers: shared hit=2
Planning Time: 0.047 ms
Execution Time: 0.011 ms

```

### s06_exists_parent

Rows: 3; serialized decoded-row JSON bytes: 28 (not wire bytes). Planning/execution: 0.056000 / 0.015000 ms.

```text
Hash Semi Join  (cost=1.14..2.22 rows=3 width=7) (actual time=0.008..0.009 rows=3 loops=1)
  Hash Cond: (c.id = o.customer_id)
  Buffers: shared hit=2
  ->  Seq Scan on customers c  (cost=0.00..1.03 rows=3 width=7) (actual time=0.002..0.003 rows=3 loops=1)
        Buffers: shared hit=1
  ->  Hash  (cost=1.07..1.07 rows=5 width=4) (actual time=0.004..0.004 rows=4 loops=1)
        Buckets: 1024  Batches: 1  Memory Usage: 9kB
        Buffers: shared hit=1
        ->  Seq Scan on orders o  (cost=0.00..1.07 rows=5 width=4) (actual time=0.002..0.002 rows=5 loops=1)
              Filter: (amount > 100)
              Rows Removed by Filter: 1
              Buffers: shared hit=1
Planning:
  Buffers: shared hit=2
Planning Time: 0.056 ms
Execution Time: 0.015 ms

```

### s07_in_parent

Rows: 3; serialized decoded-row JSON bytes: 28 (not wire bytes). Planning/execution: 0.045000 / 0.012000 ms.

```text
Hash Semi Join  (cost=1.14..2.22 rows=3 width=7) (actual time=0.006..0.007 rows=3 loops=1)
  Hash Cond: (c.id = o.customer_id)
  Buffers: shared hit=2
  ->  Seq Scan on customers c  (cost=0.00..1.03 rows=3 width=7) (actual time=0.002..0.002 rows=3 loops=1)
        Buffers: shared hit=1
  ->  Hash  (cost=1.07..1.07 rows=5 width=4) (actual time=0.003..0.003 rows=4 loops=1)
        Buckets: 1024  Batches: 1  Memory Usage: 9kB
        Buffers: shared hit=1
        ->  Seq Scan on orders o  (cost=0.00..1.07 rows=5 width=4) (actual time=0.001..0.002 rows=5 loops=1)
              Filter: (amount > 100)
              Rows Removed by Filter: 1
              Buffers: shared hit=1
Planning:
  Buffers: shared hit=2
Planning Time: 0.045 ms
Execution Time: 0.012 ms

```

### s08_join_distinct_parent

Rows: 3; serialized decoded-row JSON bytes: 28 (not wire bytes). Planning/execution: 0.042000 / 0.018000 ms.

```text
HashAggregate  (cost=2.19..2.22 rows=3 width=7) (actual time=0.009..0.010 rows=3 loops=1)
  Group Key: c.id, c.name
  Batches: 1  Memory Usage: 24kB
  Buffers: shared hit=2
  ->  Hash Join  (cost=1.07..2.17 rows=5 width=7) (actual time=0.006..0.007 rows=4 loops=1)
        Hash Cond: (o.customer_id = c.id)
        Buffers: shared hit=2
        ->  Seq Scan on orders o  (cost=0.00..1.07 rows=5 width=4) (actual time=0.002..0.003 rows=5 loops=1)
              Filter: (amount > 100)
              Rows Removed by Filter: 1
              Buffers: shared hit=1
        ->  Hash  (cost=1.03..1.03 rows=3 width=7) (actual time=0.002..0.002 rows=3 loops=1)
              Buckets: 1024  Batches: 1  Memory Usage: 9kB
              Buffers: shared hit=1
              ->  Seq Scan on customers c  (cost=0.00..1.03 rows=3 width=7) (actual time=0.001..0.001 rows=3 loops=1)
                    Buffers: shared hit=1
Planning:
  Buffers: shared hit=2
Planning Time: 0.042 ms
Execution Time: 0.018 ms

```

### s11_optional_on

Rows: 4; serialized decoded-row JSON bytes: 73 (not wire bytes). Planning/execution: 0.061000 / 0.015000 ms.

```text
Hash Right Join  (cost=1.07..2.17 rows=5 width=15) (actual time=0.007..0.009 rows=4 loops=1)
  Hash Cond: (o.customer_id = c.id)
  Buffers: shared hit=2
  ->  Seq Scan on orders o  (cost=0.00..1.07 rows=5 width=12) (actual time=0.002..0.003 rows=5 loops=1)
        Filter: (amount > 100)
        Rows Removed by Filter: 1
        Buffers: shared hit=1
  ->  Hash  (cost=1.03..1.03 rows=3 width=7) (actual time=0.003..0.003 rows=3 loops=1)
        Buckets: 1024  Batches: 1  Memory Usage: 9kB
        Buffers: shared hit=1
        ->  Seq Scan on customers c  (cost=0.00..1.03 rows=3 width=7) (actual time=0.002..0.002 rows=3 loops=1)
              Buffers: shared hit=1
Planning:
  Buffers: shared hit=2
Planning Time: 0.061 ms
Execution Time: 0.015 ms

```

### s16_page_offset

Rows: 2; serialized decoded-row JSON bytes: 23 (not wire bytes). Planning/execution: 0.011000 / 0.008000 ms.

```text
Limit  (cost=1.12..1.12 rows=2 width=8) (actual time=0.004..0.004 rows=2 loops=1)
  Buffers: shared hit=1
  ->  Sort  (cost=1.12..1.14 rows=6 width=8) (actual time=0.004..0.004 rows=2 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 25kB
        Buffers: shared hit=1
        ->  Seq Scan on orders  (cost=0.00..1.06 rows=6 width=8) (actual time=0.002..0.002 rows=6 loops=1)
              Buffers: shared hit=1
Planning Time: 0.011 ms
Execution Time: 0.008 ms

```

### s17_page_keyset

Rows: 2; serialized decoded-row JSON bytes: 23 (not wire bytes). Planning/execution: 0.025000 / 0.010000 ms.

```text
Limit  (cost=1.12..1.12 rows=2 width=8) (actual time=0.006..0.006 rows=2 loops=1)
  Buffers: shared hit=1
  ->  Sort  (cost=1.12..1.14 rows=6 width=8) (actual time=0.005..0.006 rows=2 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 25kB
        Buffers: shared hit=1
        ->  Seq Scan on orders  (cost=0.00..1.06 rows=6 width=8) (actual time=0.003..0.003 rows=6 loops=1)
              Buffers: shared hit=1
Planning Time: 0.025 ms
Execution Time: 0.010 ms

```

### s14_topn_rownumber

Rows: 5; serialized decoded-row JSON bytes: 67 (not wire bytes). Planning/execution: 0.036000 / 0.019000 ms.

```text
Subquery Scan on t  (cost=1.14..1.28 rows=5 width=12) (actual time=0.008..0.012 rows=5 loops=1)
  Buffers: shared hit=1
  ->  WindowAgg  (cost=1.14..1.23 rows=5 width=52) (actual time=0.008..0.011 rows=5 loops=1)
        Run Condition: (row_number() OVER (?) <= 3)
        Buffers: shared hit=1
        ->  Sort  (cost=1.12..1.13 rows=5 width=12) (actual time=0.007..0.007 rows=5 loops=1)
              Sort Key: o.customer_id, o.amount DESC NULLS LAST, o.id
              Sort Method: quicksort  Memory: 25kB
              Buffers: shared hit=1
              ->  Seq Scan on orders o  (cost=0.00..1.06 rows=5 width=12) (actual time=0.003..0.004 rows=5 loops=1)
                    Filter: (customer_id IS NOT NULL)
                    Rows Removed by Filter: 1
                    Buffers: shared hit=1
Planning Time: 0.036 ms
Execution Time: 0.019 ms

```

### s15_topn_lateral

Rows: 5; serialized decoded-row JSON bytes: 67 (not wire bytes). Planning/execution: 0.035000 / 0.019000 ms.

```text
Sort  (cost=4.44..4.45 rows=6 width=12) (actual time=0.013..0.014 rows=5 loops=1)
  Sort Key: c.id, o.amount DESC NULLS LAST, o.id
  Sort Method: quicksort  Memory: 25kB
  Buffers: shared hit=4
  ->  Nested Loop  (cost=1.08..4.36 rows=6 width=12) (actual time=0.006..0.012 rows=5 loops=1)
        Buffers: shared hit=4
        ->  Seq Scan on customers c  (cost=0.00..1.03 rows=3 width=4) (actual time=0.002..0.002 rows=3 loops=1)
              Buffers: shared hit=1
        ->  Limit  (cost=1.08..1.09 rows=2 width=8) (actual time=0.002..0.002 rows=2 loops=3)
              Buffers: shared hit=3
              ->  Sort  (cost=1.08..1.09 rows=2 width=8) (actual time=0.002..0.002 rows=2 loops=3)
                    Sort Key: o.amount DESC NULLS LAST, o.id
                    Sort Method: quicksort  Memory: 25kB
                    Buffers: shared hit=3
                    ->  Seq Scan on orders o  (cost=0.00..1.07 rows=2 width=8) (actual time=0.001..0.001 rows=2 loops=3)
                          Filter: (customer_id = c.id)
                          Rows Removed by Filter: 4
                          Buffers: shared hit=3
Planning Time: 0.035 ms
Execution Time: 0.019 ms

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
