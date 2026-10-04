# SELECT access-path matrix (generated measurements)

_Measurements only. Verdicts for P10/P11 belong in a signed analysis._

**TL;DR:** within-run medians for OR / UNION ALL / UNION at three predicate shapes and OFFSET / keyset at five page depths, each cell verified against the independent oracle. A within-run difference is not a replicated result. Missing measurements are N/A.

| run | engine | topology | scale | page limit |
|---|---|---|---|---|
| select-access-pg-r1-20261003 | postgres | single | deep | 5 |

## Provenance

| Field | Value |
|---|---|
| repo commit | `68b09e6d3bbbcf03deb53bd489e0059a978b192f` |
| benchmark image | `localhost/charitybench-select:select-access-pg-r1-20261003` |
| image digest | `371240b0ab41987c5579724320a2cd947c9a62cdb4a5374ad9480ef86ab52329` |
| run tag | `run/study-01/select-access-pg-r1-20261003` |
| environment | `xiles-podman-select-access-20261003` |
| engine image digest | `7296f210ae81031ec955dbad9a67a84fe958572a2153b8d0826a647522904dc1` |
| engine version | PostgreSQL 17.11 (Debian 17.11-1.pgdg13+2) on x86_64-pc-linux-gnu, compiled by gcc (Debian 14.2.0-19) 14.2.0, 64-bit |
| dataset seed | 42 |
| dataset fingerprint | `88c40592216c0ce9` |
| setting: client_workers | 1 sequential reader; 0 writers |
| setting: connection_pool | 2 = 1 reader/setup + 1 reserve |
| setting: distributed_rpc | N/A on PostgreSQL single-node |
| setting: index_set | orders_amount_idx(amount),orders_customer_idx(customer_id),primary keys; expression arm adds orders_amount_expr_idx((amount + 0)) |
| setting: page_limit | 5 |
| setting: prepared_plan_policy | pgx simple protocol; no prepared statements |
| setting: scale | deep |
| setting: snapshot_isolation | repeatable read |
| setting: statement_budget_ms | 120000 |
| setting: statistics | ANALYZE after deterministic load |
| setting: topology | single |
| setting: trials_per_cell | 5 |

## Disjunction disjointness proof

UNION ALL is a competing spelling of OR only while the two arms cannot both match one row. The count is computed by the independent Go oracle from the deterministic fixture.

| shape | arm A | arm B | A matches | B matches | overlapping rows | disjoint |
|---|---|---|---|---|---|---|
| indexable | `amount <= 100` | `amount >= 400` | 2273 | 2379 | 0 | true |
| expression | `(amount + 0) <= 100` | `(amount + 0) >= 400` | 2273 | 2379 | 0 | true |
| cast | `(amount)::bigint <= 100` | `(amount)::bigint >= 400` | 2273 | 2379 | 0 | true |

## Predicate alternatives (P10)

| shape | index arm | spelling | status | median ms | spread % | trials | rows | plan structure |
|---|---|---|---|---|---|---|---|---|
| cast | default | or | valid | 3.511739 | 69.3 | 5 | 4652 | `a4bab1f05f5e667d` |
| cast | default | union | valid | 4.413927 | 26.3 | 5 | 4652 | `6ff77b66c9bf6ee6` |
| cast | default | union_all | valid | 5.842239 | 28.6 | 5 | 4652 | `3f9f3d5c5dc15dd7` |
| expression | expression_index | or | valid | 3.608903 | 45.6 | 5 | 4652 | `a4bab1f05f5e667d` |
| expression | expression_index | union | valid | 4.502599 | 1.3 | 5 | 4652 | `d4fbfb3aac234d58` |
| expression | expression_index | union_all | valid | 4.951228 | 31.1 | 5 | 4652 | `605ee056e65421b1` |
| expression | default | or | valid | 3.697289 | 16.9 | 5 | 4652 | `a4bab1f05f5e667d` |
| expression | default | union | valid | 6.097718 | 35.1 | 5 | 4652 | `6ff77b66c9bf6ee6` |
| expression | default | union_all | valid | 6.725520 | 70.9 | 5 | 4652 | `3f9f3d5c5dc15dd7` |
| indexable | default | or | valid | 3.941529 | 78.3 | 5 | 4652 | `a4bab1f05f5e667d` |
| indexable | default | union | valid | 4.885129 | 34.7 | 5 | 4652 | `ed39b768e5b063b0` |
| indexable | default | union_all | valid | 4.678244 | 14.3 | 5 | 4652 | `d51936c88987c9fd` |

### P10 pairs

Pairs compare the **default** arm only; the labelled `expression_index` arm is a package arm and is never pooled with the query-text pair.

| pair | A median ms | B median ms | B/A | same plan | within-run verdict |
|---|---|---|---|---|---|
| shape=indexable or vs union_all | 3.941529 | 4.678244 | 1.1869 | false | medians inside the two cells' own spread; no resolved difference within this run |
| shape=indexable or vs union | 3.941529 | 4.885129 | 1.2394 | false | medians inside the two cells' own spread; no resolved difference within this run |
| shape=expression or vs union_all | 3.697289 | 6.725520 | 1.8190 | false | medians inside the two cells' own spread; no resolved difference within this run |
| shape=expression or vs union | 3.697289 | 6.097718 | 1.6492 | false | within-run medians differ; whole-run replication required before a verdict |
| shape=cast or vs union_all | 3.511739 | 5.842239 | 1.6636 | false | medians inside the two cells' own spread; no resolved difference within this run |
| shape=cast or vs union | 3.511739 | 4.413927 | 1.2569 | false | medians inside the two cells' own spread; no resolved difference within this run |

## Pagination (P11)

| depth | spelling | status | median ms | spread % | trials | rows | plan structure |
|---|---|---|---|---|---|---|---|
| 1 | offset | valid | 2.918058 | 125.5 | 5 | 5 | `ac4509d295bf3348` |
| 1 | keyset | valid | 2.864689 | 17.1 | 5 | 5 | `ac4509d295bf3348` |
| 10 | offset | valid | 2.954553 | 50.1 | 5 | 5 | `ac4509d295bf3348` |
| 10 | keyset | valid | 2.888451 | 17.0 | 5 | 5 | `ac4509d295bf3348` |
| 100 | offset | valid | 2.923831 | 25.8 | 5 | 5 | `ac4509d295bf3348` |
| 100 | keyset | valid | 2.852522 | 13.8 | 5 | 5 | `ac4509d295bf3348` |
| 1000 | offset | valid | 3.837056 | 153.3 | 5 | 5 | `ac4509d295bf3348` |
| 1000 | keyset | valid | 3.277525 | 76.6 | 5 | 5 | `ac4509d295bf3348` |
| 10000 | offset | valid | 6.078742 | 15.6 | 5 | 5 | `ac4509d295bf3348` |
| 10000 | keyset | valid | 1.199865 | 17.5 | 5 | 5 | `dd316474b64ddc29` |

Median ratio table (keyset / OFFSET; a ratio below 1 means keyset was faster in this run):

| depth | ratio (keyset/offset) | offset median ms | keyset median ms |
|---|---|---|---|
| 1 | 0.9817 | 2.918058 | 2.864689 |
| 10 | 0.9776 | 2.954553 | 2.888451 |
| 100 | 0.9756 | 2.923831 | 2.852522 |
| 1000 | 0.8542 | 3.837056 | 3.277525 |
| 10000 | 0.1974 | 6.078742 | 1.199865 |

### P11 pairs

| pair | A median ms | B median ms | B/A | same plan | within-run verdict |
|---|---|---|---|---|---|
| depth=1 offset vs keyset | 2.918058 | 2.864689 | 0.9817 | true | medians inside the two cells' own spread; no resolved difference within this run |
| depth=10 offset vs keyset | 2.954553 | 2.888451 | 0.9776 | true | medians inside the two cells' own spread; no resolved difference within this run |
| depth=100 offset vs keyset | 2.923831 | 2.852522 | 0.9756 | true | medians inside the two cells' own spread; no resolved difference within this run |
| depth=1000 offset vs keyset | 3.837056 | 3.277525 | 0.8542 | true | medians inside the two cells' own spread; no resolved difference within this run |
| depth=10000 offset vs keyset | 6.078742 | 1.199865 | 0.1974 | false | within-run medians differ; whole-run replication required before a verdict |

## Expression-index arm: write and storage cost

Index `orders_amount_expr_idx` on expression `(amount + 0)`; supported=true.

| measure | value |
|---|---|
| index size (bytes) | 114688 |
| insert batch with index (ms) | 10.585041 |
| insert batch without index (ms) | 8.687693 |
| batch rows | 500 |

## Fixed-snapshot page equality

Isolation: `repeatable read`; all depths equal: true.

| depth | offset rows | keyset rows | equal |
|---|---|---|---|
| 1 | 5 | 5 | true |
| 10 | 5 | 5 | true |
| 100 | 5 | 5 | true |
| 1000 | 5 | 5 | true |
| 10000 | 5 | 5 | true |

## Live negative controls (per predicate shape)

`x05_union_all_overlap` is this family's required control and must be detected in every shape. `x01`-`x04` belong to other contracts; they are gated live on the teaching fixture, and here they are recorded as `control_absent` when the measurement fixture cannot exercise them (for example `x01` needs an unmatched parent, which a fixture where every parent has orders does not contain). An absent control is not evidence for the cells it cannot police.

- shape=indexable `x01_optional_where`: detected=false, rows=9175 — x01_optional_where was NOT detected: its answer matched the sel_optional_match oracle
- shape=indexable `x02_not_in_nullable`: detected=false, rows=0 — x02_not_in_nullable was NOT detected: its answer matched the sel_absence oracle
- shape=indexable `x03_distinct_on_name`: detected=false, rows=200 — x03_distinct_on_name was NOT detected: its answer matched the sel_membership oracle
- shape=indexable `x04_limit_before_join`: detected=true, rows=5
- shape=indexable `x05_union_all_overlap`: detected=true, rows=4546
- shape=expression `x01_optional_where`: detected=false, rows=9175 — x01_optional_where was NOT detected: its answer matched the sel_optional_match oracle
- shape=expression `x02_not_in_nullable`: detected=false, rows=0 — x02_not_in_nullable was NOT detected: its answer matched the sel_absence oracle
- shape=expression `x03_distinct_on_name`: detected=false, rows=200 — x03_distinct_on_name was NOT detected: its answer matched the sel_membership oracle
- shape=expression `x04_limit_before_join`: detected=true, rows=5
- shape=expression `x05_union_all_overlap`: detected=true, rows=4546
- shape=cast `x01_optional_where`: detected=false, rows=9175 — x01_optional_where was NOT detected: its answer matched the sel_optional_match oracle
- shape=cast `x02_not_in_nullable`: detected=false, rows=0 — x02_not_in_nullable was NOT detected: its answer matched the sel_absence oracle
- shape=cast `x03_distinct_on_name`: detected=false, rows=200 — x03_distinct_on_name was NOT detected: its answer matched the sel_membership oracle
- shape=cast `x04_limit_before_join`: detected=true, rows=5
- shape=cast `x05_union_all_overlap`: detected=true, rows=4546

## Excluded cells

None: every cell is valid.

## Captured plans and result sizes

### s17_page_keyset [depth=10000;spelling=keyset]

Rows: 5; serialized decoded-row JSON bytes: 51 (not wire bytes). Planning/execution: 0.214000 / 0.909000 ms.

```text
Limit  (cost=154.12..154.13 rows=5 width=8) (actual time=0.871..0.874 rows=5 loops=1)
  Buffers: shared hit=73
  ->  Sort  (cost=154.12..158.90 rows=1912 width=8) (actual time=0.870..0.872 rows=5 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 25kB
        Buffers: shared hit=73
        ->  Bitmap Heap Scan on orders  (cost=22.66..122.36 rows=1912 width=8) (actual time=0.112..0.600 rows=2000 loops=1)
              Recheck Cond: ((amount < 67) OR (amount = 67) OR (amount IS NULL))
              Filter: ((amount < 67) OR ((amount = 67) AND (id > 4310)) OR (amount IS NULL))
              Rows Removed by Filter: 6
              Heap Blocks: exact=65
              Buffers: shared hit=73
              ->  BitmapOr  (cost=22.66..22.66 rows=1983 width=0) (actual time=0.088..0.090 rows=0 loops=1)
                    Buffers: shared hit=8
                    ->  Bitmap Index Scan on orders_amount_idx  (cost=0.00..14.38 rows=1440 width=0) (actual time=0.069..0.069 rows=1460 loops=1)
                          Index Cond: (amount < 67)
                          Buffers: shared hit=4
                    ->  Bitmap Index Scan on orders_amount_idx  (cost=0.00..1.54 rows=21 width=0) (actual time=0.004..0.004 rows=24 loops=1)
                          Index Cond: (amount = 67)
                          Buffers: shared hit=2
                    ->  Bitmap Index Scan on orders_amount_idx  (cost=0.00..5.30 rows=522 width=0) (actual time=0.015..0.015 rows=522 loops=1)
                          Index Cond: (amount IS NULL)
                          Buffers: shared hit=2
Planning Time: 0.214 ms
Execution Time: 0.909 ms

```

### s16_page_offset [depth=10000;spelling=offset]

Rows: 5; serialized decoded-row JSON bytes: 51 (not wire bytes). Planning/execution: 0.103000 / 7.779000 ms.

```text
Limit  (cost=1023.04..1023.06 rows=5 width=8) (actual time=7.748..7.752 rows=5 loops=1)
  Buffers: shared hit=65
  ->  Sort  (cost=998.04..1028.04 rows=12000 width=8) (actual time=6.523..7.211 rows=10005 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: quicksort  Memory: 853kB
        Buffers: shared hit=65
        ->  Seq Scan on orders  (cost=0.00..185.00 rows=12000 width=8) (actual time=0.020..1.576 rows=12000 loops=1)
              Buffers: shared hit=65
Planning Time: 0.103 ms
Execution Time: 7.779 ms

```

### s17_page_keyset [depth=1000;spelling=keyset]

Rows: 5; serialized decoded-row JSON bytes: 59 (not wire bytes). Planning/execution: 0.191000 / 5.670000 ms.

```text
Limit  (cost=450.03..450.04 rows=5 width=8) (actual time=5.632..5.637 rows=5 loops=1)
  Buffers: shared hit=65
  ->  Sort  (cost=450.03..476.38 rows=10538 width=8) (actual time=5.628..5.631 rows=5 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 25kB
        Buffers: shared hit=65
        ->  Seq Scan on orders  (cost=0.00..275.00 rows=10538 width=8) (actual time=0.026..3.098 rows=11000 loops=1)
              Filter: ((amount < 457) OR ((amount = 457) AND (id > 11254)) OR (amount IS NULL))
              Rows Removed by Filter: 1000
              Buffers: shared hit=65
Planning Time: 0.191 ms
Execution Time: 5.670 ms

```

### s16_page_offset [depth=1000;spelling=offset]

Rows: 5; serialized decoded-row JSON bytes: 59 (not wire bytes). Planning/execution: 0.086000 / 4.535000 ms.

```text
Limit  (cost=845.88..845.89 rows=5 width=8) (actual time=4.511..4.514 rows=5 loops=1)
  Buffers: shared hit=65
  ->  Sort  (cost=843.38..873.38 rows=12000 width=8) (actual time=4.395..4.458 rows=1005 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 115kB
        Buffers: shared hit=65
        ->  Seq Scan on orders  (cost=0.00..185.00 rows=12000 width=8) (actual time=0.017..1.463 rows=12000 loops=1)
              Buffers: shared hit=65
Planning Time: 0.086 ms
Execution Time: 4.535 ms

```

### s17_page_keyset [depth=100;spelling=keyset]

Rows: 5; serialized decoded-row JSON bytes: 56 (not wire bytes). Planning/execution: 0.188000 / 3.390000 ms.

```text
Limit  (cost=463.65..463.66 rows=5 width=8) (actual time=3.367..3.370 rows=5 loops=1)
  Buffers: shared hit=65
  ->  Sort  (cost=463.65..492.05 rows=11358 width=8) (actual time=3.365..3.367 rows=5 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 25kB
        Buffers: shared hit=65
        ->  Seq Scan on orders  (cost=0.00..275.00 rows=11358 width=8) (actual time=0.016..1.709 rows=11900 loops=1)
              Filter: ((amount < 497) OR ((amount = 497) AND (id > 6723)) OR (amount IS NULL))
              Rows Removed by Filter: 100
              Buffers: shared hit=65
Planning:
  Buffers: shared hit=3
Planning Time: 0.188 ms
Execution Time: 3.390 ms

```

### s16_page_offset [depth=100;spelling=offset]

Rows: 5; serialized decoded-row JSON bytes: 56 (not wire bytes). Planning/execution: 0.091000 / 3.373000 ms.

```text
Limit  (cost=648.10..648.12 rows=5 width=8) (actual time=3.352..3.355 rows=5 loops=1)
  Buffers: shared hit=65
  ->  Sort  (cost=647.85..677.85 rows=12000 width=8) (actual time=3.339..3.346 rows=105 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 31kB
        Buffers: shared hit=65
        ->  Seq Scan on orders  (cost=0.00..185.00 rows=12000 width=8) (actual time=0.016..1.473 rows=12000 loops=1)
              Buffers: shared hit=65
Planning Time: 0.091 ms
Execution Time: 3.373 ms

```

### s17_page_keyset [depth=10;spelling=keyset]

Rows: 5; serialized decoded-row JSON bytes: 56 (not wire bytes). Planning/execution: 0.359000 / 3.487000 ms.

```text
Limit  (cost=465.21..465.23 rows=5 width=8) (actual time=3.460..3.462 rows=5 loops=1)
  Buffers: shared hit=65
  ->  Sort  (cost=465.21..493.84 rows=11452 width=8) (actual time=3.458..3.459 rows=5 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 25kB
        Buffers: shared hit=65
        ->  Seq Scan on orders  (cost=0.00..275.00 rows=11452 width=8) (actual time=0.019..1.728 rows=11990 loops=1)
              Filter: ((amount < 500) OR ((amount = 500) AND (id > 6319)) OR (amount IS NULL))
              Rows Removed by Filter: 10
              Buffers: shared hit=65
Planning:
  Buffers: shared hit=3
Planning Time: 0.359 ms
Execution Time: 3.487 ms

```

### s16_page_offset [depth=10;spelling=offset]

Rows: 5; serialized decoded-row JSON bytes: 56 (not wire bytes). Planning/execution: 0.079000 / 3.353000 ms.

```text
Limit  (cost=479.44..479.45 rows=5 width=8) (actual time=3.333..3.335 rows=5 loops=1)
  Buffers: shared hit=65
  ->  Sort  (cost=479.41..509.41 rows=12000 width=8) (actual time=3.330..3.331 rows=15 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 25kB
        Buffers: shared hit=65
        ->  Seq Scan on orders  (cost=0.00..185.00 rows=12000 width=8) (actual time=0.014..1.542 rows=12000 loops=1)
              Buffers: shared hit=65
Planning Time: 0.079 ms
Execution Time: 3.353 ms

```

### s17_page_keyset [depth=1;spelling=keyset]

Rows: 5; serialized decoded-row JSON bytes: 56 (not wire bytes). Planning/execution: 0.496000 / 3.902000 ms.

```text
Limit  (cost=465.21..465.23 rows=5 width=8) (actual time=3.794..3.797 rows=5 loops=1)
  Buffers: shared hit=65
  ->  Sort  (cost=465.21..493.84 rows=11452 width=8) (actual time=3.792..3.793 rows=5 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 25kB
        Buffers: shared hit=65
        ->  Seq Scan on orders  (cost=0.00..275.00 rows=11452 width=8) (actual time=0.023..1.952 rows=11999 loops=1)
              Filter: ((amount < 500) OR ((amount = 500) AND (id > 1554)) OR (amount IS NULL))
              Rows Removed by Filter: 1
              Buffers: shared hit=65
Planning:
  Buffers: shared hit=3
Planning Time: 0.496 ms
Execution Time: 3.902 ms

```

### s16_page_offset [depth=1;spelling=offset]

Rows: 5; serialized decoded-row JSON bytes: 56 (not wire bytes). Planning/execution: 0.111000 / 4.455000 ms.

```text
Limit  (cost=400.10..400.11 rows=5 width=8) (actual time=4.428..4.430 rows=5 loops=1)
  Buffers: shared hit=65
  ->  Sort  (cost=400.10..430.10 rows=12000 width=8) (actual time=4.425..4.426 rows=6 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 25kB
        Buffers: shared hit=65
        ->  Seq Scan on orders  (cost=0.00..185.00 rows=12000 width=8) (actual time=0.020..2.443 rows=12000 loops=1)
              Buffers: shared hit=65
Planning Time: 0.111 ms
Execution Time: 4.455 ms

```

### s18_or_predicate [shape=cast;spelling=or]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.080000 / 2.582000 ms.

```text
Index Scan using orders_pkey on orders  (cost=0.29..403.89 rows=6667 width=8) (actual time=0.021..2.301 rows=4652 loops=1)
  Filter: (((amount)::bigint <= 100) OR ((amount)::bigint >= 400))
  Rows Removed by Filter: 7348
  Buffers: shared hit=99
Planning Time: 0.080 ms
Execution Time: 2.582 ms

```

### s20_union_predicate [shape=cast;spelling=union]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.144000 / 4.721000 ms.

```text
Unique  (cost=968.64..1108.64 rows=8000 width=8) (actual time=2.689..4.417 rows=4652 loops=1)
  Buffers: shared hit=130
  ->  Merge Append  (cost=968.64..1068.64 rows=8000 width=8) (actual time=2.688..3.674 rows=4652 loops=1)
        Sort Key: orders.id, orders.amount
        Buffers: shared hit=130
        ->  Sort  (cost=484.32..494.32 rows=4000 width=8) (actual time=1.404..1.542 rows=2273 loops=1)
              Sort Key: orders.id, orders.amount
              Sort Method: quicksort  Memory: 150kB
              Buffers: shared hit=65
              ->  Seq Scan on orders  (cost=0.00..245.00 rows=4000 width=8) (actual time=0.021..1.066 rows=2273 loops=1)
                    Filter: ((amount)::bigint <= 100)
                    Rows Removed by Filter: 9727
                    Buffers: shared hit=65
        ->  Sort  (cost=484.32..494.32 rows=4000 width=8) (actual time=1.282..1.425 rows=2379 loops=1)
              Sort Key: orders_1.id, orders_1.amount
              Sort Method: quicksort  Memory: 152kB
              Buffers: shared hit=65
              ->  Seq Scan on orders orders_1  (cost=0.00..245.00 rows=4000 width=8) (actual time=0.008..0.948 rows=2379 loops=1)
                    Filter: ((amount)::bigint >= 400)
                    Rows Removed by Filter: 9621
                    Buffers: shared hit=65
Planning Time: 0.144 ms
Execution Time: 4.721 ms

```

### s19_union_all_predicate [shape=cast;spelling=union_all]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.119000 / 4.173000 ms.

```text
Sort  (cost=1048.63..1068.63 rows=8000 width=8) (actual time=3.585..3.882 rows=4652 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 302kB
  Buffers: shared hit=130
  ->  Append  (cost=0.00..530.00 rows=8000 width=8) (actual time=0.019..2.580 rows=4652 loops=1)
        Buffers: shared hit=130
        ->  Seq Scan on orders  (cost=0.00..245.00 rows=4000 width=8) (actual time=0.018..1.047 rows=2273 loops=1)
              Filter: ((amount)::bigint <= 100)
              Rows Removed by Filter: 9727
              Buffers: shared hit=65
        ->  Seq Scan on orders orders_1  (cost=0.00..245.00 rows=4000 width=8) (actual time=0.008..1.006 rows=2379 loops=1)
              Filter: ((amount)::bigint >= 400)
              Rows Removed by Filter: 9621
              Buffers: shared hit=65
Planning Time: 0.119 ms
Execution Time: 4.173 ms

```

### s18_or_predicate [shape=expression;index=expression_index;spelling=or]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.115000 / 2.687000 ms.

```text
Index Scan using orders_pkey on orders  (cost=0.29..403.89 rows=4211 width=8) (actual time=0.028..2.359 rows=4652 loops=1)
  Filter: (((amount + 0) <= 100) OR ((amount + 0) >= 400))
  Rows Removed by Filter: 7348
  Buffers: shared hit=99
Planning Time: 0.115 ms
Execution Time: 2.687 ms

```

### s20_union_predicate [shape=expression;index=expression_index;spelling=union]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.170000 / 4.916000 ms.

```text
Unique  (cost=492.69..574.31 rows=4664 width=8) (actual time=2.803..4.624 rows=4652 loops=1)
  Buffers: shared hit=3934
  ->  Merge Append  (cost=492.69..550.99 rows=4664 width=8) (actual time=2.802..3.828 rows=4652 loops=1)
        Sort Key: orders.id, orders.amount
        Buffers: shared hit=3934
        ->  Sort  (cost=241.10..246.77 rows=2267 width=8) (actual time=1.389..1.537 rows=2273 loops=1)
              Sort Key: orders.id, orders.amount
              Sort Method: quicksort  Memory: 150kB
              Buffers: shared hit=1921
              ->  Index Scan using orders_amount_expr_idx on orders  (cost=0.29..114.75 rows=2267 width=8) (actual time=0.041..0.820 rows=2273 loops=1)
                    Index Cond: ((amount + 0) <= 100)
                    Buffers: shared hit=1921
        ->  Sort  (cost=251.58..257.57 rows=2397 width=8) (actual time=1.411..1.583 rows=2379 loops=1)
              Sort Key: orders_1.id, orders_1.amount
              Sort Method: quicksort  Memory: 152kB
              Buffers: shared hit=2013
              ->  Index Scan using orders_amount_expr_idx on orders orders_1  (cost=0.29..117.03 rows=2397 width=8) (actual time=0.008..0.797 rows=2379 loops=1)
                    Index Cond: ((amount + 0) >= 400)
                    Buffers: shared hit=2013
Planning Time: 0.170 ms
Execution Time: 4.916 ms

```

### s19_union_all_predicate [shape=expression;index=expression_index;spelling=union_all]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.122000 / 3.794000 ms.

```text
Sort  (cost=539.31..550.97 rows=4664 width=8) (actual time=3.228..3.519 rows=4652 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 302kB
  Buffers: shared hit=3934
  ->  Append  (cost=0.29..255.10 rows=4664 width=8) (actual time=0.015..1.987 rows=4652 loops=1)
        Buffers: shared hit=3934
        ->  Index Scan using orders_amount_expr_idx on orders  (cost=0.29..114.75 rows=2267 width=8) (actual time=0.014..0.773 rows=2273 loops=1)
              Index Cond: ((amount + 0) <= 100)
              Buffers: shared hit=1921
        ->  Index Scan using orders_amount_expr_idx on orders orders_1  (cost=0.29..117.03 rows=2397 width=8) (actual time=0.007..0.702 rows=2379 loops=1)
              Index Cond: ((amount + 0) >= 400)
              Buffers: shared hit=2013
Planning Time: 0.122 ms
Execution Time: 3.794 ms

```

### s18_or_predicate [shape=expression;spelling=or]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.092000 / 2.647000 ms.

```text
Index Scan using orders_pkey on orders  (cost=0.29..403.89 rows=6667 width=8) (actual time=0.029..2.355 rows=4652 loops=1)
  Filter: (((amount + 0) <= 100) OR ((amount + 0) >= 400))
  Rows Removed by Filter: 7348
  Buffers: shared hit=99
Planning Time: 0.092 ms
Execution Time: 2.647 ms

```

### s20_union_predicate [shape=expression;spelling=union]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.138000 / 5.061000 ms.

```text
Unique  (cost=968.64..1108.64 rows=8000 width=8) (actual time=2.879..4.765 rows=4652 loops=1)
  Buffers: shared hit=130
  ->  Merge Append  (cost=968.64..1068.64 rows=8000 width=8) (actual time=2.878..3.916 rows=4652 loops=1)
        Sort Key: orders.id, orders.amount
        Buffers: shared hit=130
        ->  Sort  (cost=484.32..494.32 rows=4000 width=8) (actual time=1.475..1.619 rows=2273 loops=1)
              Sort Key: orders.id, orders.amount
              Sort Method: quicksort  Memory: 150kB
              Buffers: shared hit=65
              ->  Seq Scan on orders  (cost=0.00..245.00 rows=4000 width=8) (actual time=0.017..1.111 rows=2273 loops=1)
                    Filter: ((amount + 0) <= 100)
                    Rows Removed by Filter: 9727
                    Buffers: shared hit=65
        ->  Sort  (cost=484.32..494.32 rows=4000 width=8) (actual time=1.398..1.567 rows=2379 loops=1)
              Sort Key: orders_1.id, orders_1.amount
              Sort Method: quicksort  Memory: 152kB
              Buffers: shared hit=65
              ->  Seq Scan on orders orders_1  (cost=0.00..245.00 rows=4000 width=8) (actual time=0.015..1.050 rows=2379 loops=1)
                    Filter: ((amount + 0) >= 400)
                    Rows Removed by Filter: 9621
                    Buffers: shared hit=65
Planning Time: 0.138 ms
Execution Time: 5.061 ms

```

### s19_union_all_predicate [shape=expression;spelling=union_all]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.161000 / 4.192000 ms.

```text
Sort  (cost=1048.63..1068.63 rows=8000 width=8) (actual time=3.627..3.910 rows=4652 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 302kB
  Buffers: shared hit=130
  ->  Append  (cost=0.00..530.00 rows=8000 width=8) (actual time=0.023..2.573 rows=4652 loops=1)
        Buffers: shared hit=130
        ->  Seq Scan on orders  (cost=0.00..245.00 rows=4000 width=8) (actual time=0.022..1.056 rows=2273 loops=1)
              Filter: ((amount + 0) <= 100)
              Rows Removed by Filter: 9727
              Buffers: shared hit=65
        ->  Seq Scan on orders orders_1  (cost=0.00..245.00 rows=4000 width=8) (actual time=0.008..1.007 rows=2379 loops=1)
              Filter: ((amount + 0) >= 400)
              Rows Removed by Filter: 9621
              Buffers: shared hit=65
Planning Time: 0.161 ms
Execution Time: 4.192 ms

```

### s18_or_predicate [shape=indexable;spelling=or]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.134000 / 2.480000 ms.

```text
Index Scan using orders_pkey on orders  (cost=0.29..343.89 rows=4211 width=8) (actual time=0.027..2.211 rows=4652 loops=1)
  Filter: ((amount <= 100) OR (amount >= 400))
  Rows Removed by Filter: 7348
  Buffers: shared hit=99
Planning Time: 0.134 ms
Execution Time: 2.480 ms

```

### s20_union_predicate [shape=indexable;spelling=union]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.164000 / 6.776000 ms.

```text
Unique  (cost=495.99..577.61 rows=4664 width=8) (actual time=4.624..6.483 rows=4652 loops=1)
  Buffers: shared hit=3936
  ->  Merge Append  (cost=495.99..554.29 rows=4664 width=8) (actual time=4.619..5.708 rows=4652 loops=1)
        Sort Key: orders.id, orders.amount
        Buffers: shared hit=3936
        ->  Sort  (cost=242.20..247.87 rows=2267 width=8) (actual time=1.303..1.475 rows=2273 loops=1)
              Sort Key: orders.id, orders.amount
              Sort Method: quicksort  Memory: 150kB
              Buffers: shared hit=1922
              ->  Index Scan using orders_amount_idx on orders  (cost=0.29..115.85 rows=2267 width=8) (actual time=0.019..0.767 rows=2273 loops=1)
                    Index Cond: (amount <= 100)
                    Buffers: shared hit=1922
        ->  Sort  (cost=253.78..259.77 rows=2397 width=8) (actual time=3.307..3.474 rows=2379 loops=1)
              Sort Key: orders_1.id, orders_1.amount
              Sort Method: quicksort  Memory: 152kB
              Buffers: shared hit=2014
              ->  Index Scan using orders_amount_idx on orders orders_1  (cost=0.29..119.23 rows=2397 width=8) (actual time=0.012..1.413 rows=2379 loops=1)
                    Index Cond: (amount >= 400)
                    Buffers: shared hit=2014
Planning Time: 0.164 ms
Execution Time: 6.776 ms

```

### s19_union_all_predicate [shape=indexable;spelling=union_all]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.140000 / 5.232000 ms.

```text
Sort  (cost=542.61..554.27 rows=4664 width=8) (actual time=4.647..4.962 rows=4652 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 302kB
  Buffers: shared hit=3936
  ->  Append  (cost=0.29..258.40 rows=4664 width=8) (actual time=0.026..3.033 rows=4652 loops=1)
        Buffers: shared hit=3936
        ->  Index Scan using orders_amount_idx on orders  (cost=0.29..115.85 rows=2267 width=8) (actual time=0.025..1.747 rows=2273 loops=1)
              Index Cond: (amount <= 100)
              Buffers: shared hit=1922
        ->  Index Scan using orders_amount_idx on orders orders_1  (cost=0.29..119.23 rows=2397 width=8) (actual time=0.015..0.767 rows=2379 loops=1)
              Index Cond: (amount >= 400)
              Buffers: shared hit=2014
Planning Time: 0.140 ms
Execution Time: 5.232 ms

```

## Limits and gaps recorded by the runner

- no concurrent writers: client-decode JSON bytes are not wire bytes; tails are under-sampled at the protocol trial count
- one host: no independent hardware, real-network, 3-node placement or equal-total-budget comparison is established
- single-node only; the YB 3-node RF=3 per-node/equal-total arms are a coverage gap on this host

