# SELECT access-path matrix (generated measurements)

_Measurements only. Verdicts for P10/P11 belong in a signed analysis._

**TL;DR:** within-run medians for OR / UNION ALL / UNION at three predicate shapes and OFFSET / keyset at five page depths, each cell verified against the independent oracle. A within-run difference is not a replicated result. Missing measurements are N/A.

| run | engine | topology | scale | page limit |
|---|---|---|---|---|
| select-access-yb-r1-20261003 | yugabyte | single | deep | 5 |

## Provenance

| Field | Value |
|---|---|
| repo commit | `2dae4b5acf56d2a70e2b8e556f365d16079e478f` |
| benchmark image | `localhost/charitybench-select:select-access-yb-r1-20261003` |
| image digest | `528af40e70150ef1c049d260c18d1c1db1474b3a4e495eea15d8e7b851f5eb5c` |
| run tag | `run/study-01/select-access-yb-r1-20261003` |
| environment | `xiles-podman-select-access-20261003` |
| engine image digest | `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35` |
| engine version | PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit |
| dataset seed | 42 |
| dataset fingerprint | `88c40592216c0ce9` |
| setting: client_workers | 1 sequential reader; 0 writers |
| setting: connection_pool | 2 = 1 reader/setup + 1 reserve |
| setting: distributed_rpc | not parsed; raw EXPLAIN retained (coverage gap) |
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
| cast | default | or | valid | 5.912337 | 13.5 | 5 | 4652 | `1dd363b32be056fc` |
| cast | default | union | valid | 8.378739 | 6.9 | 5 | 4652 | `b5f6c293b9267d60` |
| cast | default | union_all | valid | 7.869411 | 18.1 | 5 | 4652 | `24cabca42adfa691` |
| expression | expression_index | or | valid | 6.417091 | 23.1 | 5 | 4652 | `0e731990206e1446` |
| expression | expression_index | union | valid | 8.730833 | 29.1 | 5 | 4652 | `bf74591baa836f9d` |
| expression | expression_index | union_all | valid | 9.030609 | 29.6 | 5 | 4652 | `a7e687668ac364de` |
| expression | default | or | valid | 6.281501 | 13.7 | 5 | 4652 | `0e731990206e1446` |
| expression | default | union | valid | 8.601029 | 24.3 | 5 | 4652 | `bf74591baa836f9d` |
| expression | default | union_all | valid | 8.106574 | 4.9 | 5 | 4652 | `a7e687668ac364de` |
| indexable | default | or | valid | 5.357923 | 18.1 | 5 | 4652 | `5191eeb70daef797` |
| indexable | default | union | valid | 8.117927 | 4.2 | 5 | 4652 | `12de8f648f9520a0` |
| indexable | default | union_all | valid | 8.102317 | 14.6 | 5 | 4652 | `d95e057d512f5ae5` |

### P10 pairs

Pairs compare the **default** arm only; the labelled `expression_index` arm is a package arm and is never pooled with the query-text pair.

| pair | A median ms | B median ms | B/A | same plan | within-run verdict |
|---|---|---|---|---|---|
| shape=indexable or vs union_all | 5.357923 | 8.102317 | 1.5122 | false | within-run medians differ; whole-run replication required before a verdict |
| shape=indexable or vs union | 5.357923 | 8.117927 | 1.5151 | false | within-run medians differ; whole-run replication required before a verdict |
| shape=expression or vs union_all | 6.281501 | 8.106574 | 1.2905 | false | within-run medians differ; whole-run replication required before a verdict |
| shape=expression or vs union | 6.281501 | 8.601029 | 1.3693 | false | medians inside the two cells' own spread; no resolved difference within this run |
| shape=cast or vs union_all | 5.912337 | 7.869411 | 1.3310 | false | within-run medians differ; whole-run replication required before a verdict |
| shape=cast or vs union | 5.912337 | 8.378739 | 1.4172 | false | within-run medians differ; whole-run replication required before a verdict |

## Pagination (P11)

| depth | spelling | status | median ms | spread % | trials | rows | plan structure |
|---|---|---|---|---|---|---|---|
| 1 | offset | valid | 3.972189 | 6.6 | 5 | 5 | `ac4509d295bf3348` |
| 1 | keyset | valid | 5.216671 | 7.9 | 5 | 5 | `166777c9921144af` |
| 10 | offset | valid | 4.181904 | 22.4 | 5 | 5 | `ac4509d295bf3348` |
| 10 | keyset | valid | 5.391089 | 8.5 | 5 | 5 | `d596ca0d42dc4da7` |
| 100 | offset | valid | 4.369848 | 12.0 | 5 | 5 | `ac4509d295bf3348` |
| 100 | keyset | valid | 5.850156 | 17.7 | 5 | 5 | `727a2962d7392ccd` |
| 1000 | offset | valid | 5.134938 | 11.8 | 5 | 5 | `ac4509d295bf3348` |
| 1000 | keyset | valid | 5.912185 | 11.6 | 5 | 5 | `441b4191cb138ed5` |
| 10000 | offset | valid | 6.335135 | 22.7 | 5 | 5 | `ac4509d295bf3348` |
| 10000 | keyset | valid | 3.767875 | 7.2 | 5 | 5 | `e1c72b3614c26105` |

Median ratio table (keyset / OFFSET; a ratio below 1 means keyset was faster in this run):

| depth | ratio (keyset/offset) | offset median ms | keyset median ms |
|---|---|---|---|
| 1 | 1.3133 | 3.972189 | 5.216671 |
| 10 | 1.2891 | 4.181904 | 5.391089 |
| 100 | 1.3388 | 4.369848 | 5.850156 |
| 1000 | 1.1514 | 5.134938 | 5.912185 |
| 10000 | 0.5948 | 6.335135 | 3.767875 |

### P11 pairs

| pair | A median ms | B median ms | B/A | same plan | within-run verdict |
|---|---|---|---|---|---|
| depth=1 offset vs keyset | 3.972189 | 5.216671 | 1.3133 | false | within-run medians differ; whole-run replication required before a verdict |
| depth=10 offset vs keyset | 4.181904 | 5.391089 | 1.2891 | false | medians inside the two cells' own spread; no resolved difference within this run |
| depth=100 offset vs keyset | 4.369848 | 5.850156 | 1.3388 | false | within-run medians differ; whole-run replication required before a verdict |
| depth=1000 offset vs keyset | 5.134938 | 5.912185 | 1.1514 | false | medians inside the two cells' own spread; no resolved difference within this run |
| depth=10000 offset vs keyset | 6.335135 | 3.767875 | 0.5948 | false | within-run medians differ; whole-run replication required before a verdict |

## Expression-index arm: write and storage cost

Index `orders_amount_expr_idx` on expression `(amount + 0)`; supported=true.

| measure | value |
|---|---|
| index size (bytes) | 0 |
| insert batch with index (ms) | 17.744512 |
| insert batch without index (ms) | 20.975606 |
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

Rows: 5; serialized decoded-row JSON bytes: 51 (not wire bytes). Planning/execution: 0.080000 / 3.624000 ms.

```text
Limit  (cost=884.51..884.52 rows=5 width=8) (actual time=3.587..3.588 rows=5 loops=1)
  ->  Sort  (cost=884.51..889.29 rows=1912 width=8) (actual time=3.585..3.586 rows=5 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 25kB
        ->  Seq Scan on orders  (cost=21.08..852.75 rows=1912 width=8) (actual time=1.817..3.378 rows=2000 loops=1)
              Storage Filter: ((amount < 67) OR ((amount = 67) AND (id > 4310)) OR (amount IS NULL))
Planning Time: 0.080 ms
Execution Time: 3.624 ms
Peak Memory Usage: 92 kB

```

### s16_page_offset [depth=10000;spelling=offset]

Rows: 5; serialized decoded-row JSON bytes: 51 (not wire bytes). Planning/execution: 0.066000 / 7.124000 ms.

```text
Limit  (cost=1879.09..1879.10 rows=5 width=8) (actual time=7.070..7.071 rows=5 loops=1)
  ->  Sort  (cost=1854.09..1884.09 rows=12000 width=8) (actual time=5.951..6.323 rows=10005 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: quicksort  Memory: 947kB
        ->  Seq Scan on orders  (cost=21.08..1041.05 rows=12000 width=8) (actual time=0.740..3.380 rows=12000 loops=1)
Planning Time: 0.066 ms
Execution Time: 7.124 ms
Peak Memory Usage: 1516 kB

```

### s17_page_keyset [depth=1000;spelling=keyset]

Rows: 5; serialized decoded-row JSON bytes: 59 (not wire bytes). Planning/execution: 0.073000 / 5.653000 ms.

```text
Limit  (cost=1274.79..1274.81 rows=5 width=8) (actual time=5.617..5.619 rows=5 loops=1)
  ->  Sort  (cost=1274.79..1301.14 rows=10538 width=8) (actual time=5.616..5.617 rows=5 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 25kB
        ->  Seq Scan on orders  (cost=21.08..1099.76 rows=10538 width=8) (actual time=0.554..4.513 rows=11000 loops=1)
              Storage Filter: ((amount < 457) OR ((amount = 457) AND (id > 11254)) OR (amount IS NULL))
Planning Time: 0.073 ms
Execution Time: 5.653 ms
Peak Memory Usage: 92 kB

```

### s16_page_offset [depth=1000;spelling=offset]

Rows: 5; serialized decoded-row JSON bytes: 59 (not wire bytes). Planning/execution: 0.052000 / 3.971000 ms.

```text
Limit  (cost=1701.93..1701.94 rows=5 width=8) (actual time=3.900..3.902 rows=5 loops=1)
  ->  Sort  (cost=1699.43..1729.43 rows=12000 width=8) (actual time=3.804..3.828 rows=1005 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 130kB
        ->  Seq Scan on orders  (cost=21.08..1041.05 rows=12000 width=8) (actual time=0.452..2.158 rows=12000 loops=1)
Planning Time: 0.052 ms
Execution Time: 3.971 ms
Peak Memory Usage: 276 kB

```

### s17_page_keyset [depth=100;spelling=keyset]

Rows: 5; serialized decoded-row JSON bytes: 56 (not wire bytes). Planning/execution: 0.074000 / 4.852000 ms.

```text
Limit  (cost=1315.00..1315.01 rows=5 width=8) (actual time=4.819..4.820 rows=5 loops=1)
  ->  Sort  (cost=1315.00..1343.41 rows=11366 width=8) (actual time=4.818..4.818 rows=5 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 25kB
        ->  Seq Scan on orders  (cost=21.08..1126.21 rows=11366 width=8) (actual time=0.553..3.694 rows=11900 loops=1)
              Storage Filter: ((amount < 497) OR ((amount = 497) AND (id > 6723)) OR (amount IS NULL))
Planning Time: 0.074 ms
Execution Time: 4.852 ms
Peak Memory Usage: 92 kB

```

### s16_page_offset [depth=100;spelling=offset]

Rows: 5; serialized decoded-row JSON bytes: 56 (not wire bytes). Planning/execution: 0.053000 / 4.256000 ms.

```text
Limit  (cost=1504.15..1504.16 rows=5 width=8) (actual time=4.216..4.218 rows=5 loops=1)
  ->  Sort  (cost=1503.90..1533.90 rows=12000 width=8) (actual time=4.205..4.208 rows=105 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 33kB
        ->  Seq Scan on orders  (cost=21.08..1041.05 rows=12000 width=8) (actual time=0.559..2.886 rows=12000 loops=1)
Planning Time: 0.053 ms
Execution Time: 4.256 ms
Peak Memory Usage: 128 kB

```

### s17_page_keyset [depth=10;spelling=keyset]

Rows: 5; serialized decoded-row JSON bytes: 56 (not wire bytes). Planning/execution: 0.069000 / 4.893000 ms.

```text
Limit  (cost=1317.54..1317.55 rows=5 width=8) (actual time=4.860..4.862 rows=5 loops=1)
  ->  Sort  (cost=1317.54..1346.22 rows=11471 width=8) (actual time=4.859..4.860 rows=5 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 25kB
        ->  Seq Scan on orders  (cost=21.08..1127.01 rows=11471 width=8) (actual time=0.617..3.729 rows=11990 loops=1)
              Storage Filter: ((amount < 500) OR ((amount = 500) AND (id > 6319)) OR (amount IS NULL))
Planning Time: 0.069 ms
Execution Time: 4.893 ms
Peak Memory Usage: 92 kB

```

### s16_page_offset [depth=10;spelling=offset]

Rows: 5; serialized decoded-row JSON bytes: 56 (not wire bytes). Planning/execution: 0.052000 / 3.715000 ms.

```text
Limit  (cost=1335.49..1335.50 rows=5 width=8) (actual time=3.682..3.683 rows=5 loops=1)
  ->  Sort  (cost=1335.46..1365.46 rows=12000 width=8) (actual time=3.680..3.680 rows=15 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 26kB
        ->  Seq Scan on orders  (cost=21.08..1041.05 rows=12000 width=8) (actual time=0.563..2.525 rows=12000 loops=1)
Planning Time: 0.052 ms
Execution Time: 3.715 ms
Peak Memory Usage: 92 kB

```

### s17_page_keyset [depth=1;spelling=keyset]

Rows: 5; serialized decoded-row JSON bytes: 56 (not wire bytes). Planning/execution: 0.072000 / 4.916000 ms.

```text
Limit  (cost=1317.57..1317.58 rows=5 width=8) (actual time=4.882..4.883 rows=5 loops=1)
  ->  Sort  (cost=1317.57..1346.25 rows=11472 width=8) (actual time=4.881..4.881 rows=5 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 25kB
        ->  Seq Scan on orders  (cost=21.08..1127.02 rows=11472 width=8) (actual time=0.510..3.721 rows=11999 loops=1)
              Storage Filter: ((amount < 500) OR ((amount = 500) AND (id > 1554)) OR (amount IS NULL))
Planning Time: 0.072 ms
Execution Time: 4.916 ms
Peak Memory Usage: 124 kB

```

### s16_page_offset [depth=1;spelling=offset]

Rows: 5; serialized decoded-row JSON bytes: 56 (not wire bytes). Planning/execution: 0.051000 / 3.750000 ms.

```text
Limit  (cost=1256.15..1256.16 rows=5 width=8) (actual time=3.718..3.720 rows=5 loops=1)
  ->  Sort  (cost=1256.15..1286.15 rows=12000 width=8) (actual time=3.717..3.718 rows=6 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 25kB
        ->  Seq Scan on orders  (cost=21.08..1041.05 rows=12000 width=8) (actual time=0.536..2.509 rows=12000 loops=1)
Planning Time: 0.051 ms
Execution Time: 3.750 ms
Peak Memory Usage: 92 kB

```

### s18_or_predicate [shape=cast;spelling=or]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.059000 / 5.030000 ms.

```text
Sort  (cost=1443.14..1459.81 rows=6667 width=8) (actual time=4.520..4.673 rows=4652 loops=1)
  Sort Key: id
  Sort Method: quicksort  Memory: 411kB
  ->  Seq Scan on orders  (cost=21.08..1019.70 rows=6667 width=8) (actual time=1.005..3.714 rows=4652 loops=1)
        Storage Filter: (((amount)::bigint <= 100) OR ((amount)::bigint >= 400))
Planning Time: 0.059 ms
Execution Time: 5.030 ms
Peak Memory Usage: 548 kB

```

### s20_union_predicate [shape=cast;spelling=union]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.072000 / 8.157000 ms.

```text
Sort  (cost=2516.53..2536.53 rows=8000 width=8) (actual time=7.609..7.753 rows=4652 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 411kB
  ->  HashAggregate  (cost=1917.90..1997.90 rows=8000 width=8) (actual time=6.580..6.874 rows=4652 loops=1)
        Group Key: orders.id, orders.amount
        Batches: 1  Memory Usage: 657kB
        ->  Append  (cost=21.08..1877.90 rows=8000 width=8) (actual time=1.515..5.824 rows=4652 loops=1)
              ->  Seq Scan on orders  (cost=21.08..878.95 rows=4000 width=8) (actual time=1.514..2.990 rows=2273 loops=1)
                    Storage Filter: ((amount)::bigint <= 100)
              ->  Seq Scan on orders orders_1  (cost=21.08..878.95 rows=4000 width=8) (actual time=1.212..2.385 rows=2379 loops=1)
                    Storage Filter: ((amount)::bigint >= 400)
Planning Time: 0.072 ms
Execution Time: 8.157 ms
Peak Memory Usage: 1252 kB

```

### s19_union_all_predicate [shape=cast;spelling=union_all]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.089000 / 7.622000 ms.

```text
Sort  (cost=2396.53..2416.53 rows=8000 width=8) (actual time=7.092..7.253 rows=4652 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 411kB
  ->  Append  (cost=21.08..1877.90 rows=8000 width=8) (actual time=1.649..6.318 rows=4652 loops=1)
        ->  Seq Scan on orders  (cost=21.08..878.95 rows=4000 width=8) (actual time=1.647..3.040 rows=2273 loops=1)
              Storage Filter: ((amount)::bigint <= 100)
        ->  Seq Scan on orders orders_1  (cost=21.08..878.95 rows=4000 width=8) (actual time=1.217..2.823 rows=2379 loops=1)
              Storage Filter: ((amount)::bigint >= 400)
Planning Time: 0.089 ms
Execution Time: 7.622 ms
Peak Memory Usage: 625 kB

```

### s18_or_predicate [shape=expression;index=expression_index;spelling=or]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.071000 / 5.257000 ms.

```text
Sort  (cost=1214.19..1224.72 rows=4211 width=8) (actual time=4.768..4.904 rows=4652 loops=1)
  Sort Key: id
  Sort Method: quicksort  Memory: 411kB
  ->  Seq Scan on orders  (cost=21.08..960.69 rows=4211 width=8) (actual time=1.163..3.986 rows=4652 loops=1)
        Storage Filter: (((amount + 0) <= 100) OR ((amount + 0) >= 400))
Planning Time: 0.071 ms
Execution Time: 5.257 ms
Peak Memory Usage: 545 kB

```

### s20_union_predicate [shape=expression;index=expression_index;spelling=union]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.086000 / 10.304000 ms.

```text
Sort  (cost=2116.31..2127.97 rows=4664 width=8) (actual time=9.791..9.920 rows=4652 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 411kB
  ->  HashAggregate  (cost=1785.46..1832.10 rows=4664 width=8) (actual time=8.675..8.961 rows=4652 loops=1)
        Group Key: orders.id, orders.amount
        Batches: 1  Memory Usage: 465kB
        ->  Append  (cost=21.08..1762.14 rows=4664 width=8) (actual time=1.646..7.952 rows=4652 loops=1)
              ->  Seq Scan on orders  (cost=21.08..845.59 rows=2267 width=8) (actual time=1.645..3.716 rows=2273 loops=1)
                    Storage Filter: ((amount + 0) <= 100)
              ->  Seq Scan on orders orders_1  (cost=21.08..846.59 rows=2397 width=8) (actual time=1.306..3.774 rows=2379 loops=1)
                    Storage Filter: ((amount + 0) >= 400)
Planning Time: 0.086 ms
Execution Time: 10.304 ms
Peak Memory Usage: 1060 kB

```

### s19_union_all_predicate [shape=expression;index=expression_index;spelling=union_all]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.139000 / 8.545000 ms.

```text
Sort  (cost=2046.35..2058.01 rows=4664 width=8) (actual time=7.986..8.173 rows=4652 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 411kB
  ->  Append  (cost=21.08..1762.14 rows=4664 width=8) (actual time=1.772..7.138 rows=4652 loops=1)
        ->  Seq Scan on orders  (cost=21.08..845.59 rows=2267 width=8) (actual time=1.770..3.138 rows=2273 loops=1)
              Storage Filter: ((amount + 0) <= 100)
        ->  Seq Scan on orders orders_1  (cost=21.08..846.59 rows=2397 width=8) (actual time=2.034..3.522 rows=2379 loops=1)
              Storage Filter: ((amount + 0) >= 400)
Planning Time: 0.139 ms
Execution Time: 8.545 ms
Peak Memory Usage: 625 kB

```

### s18_or_predicate [shape=expression;spelling=or]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.060000 / 6.249000 ms.

```text
Sort  (cost=1443.14..1459.81 rows=6667 width=8) (actual time=5.685..5.860 rows=4652 loops=1)
  Sort Key: id
  Sort Method: quicksort  Memory: 411kB
  ->  Seq Scan on orders  (cost=21.08..1019.70 rows=6667 width=8) (actual time=1.706..4.794 rows=4652 loops=1)
        Storage Filter: (((amount + 0) <= 100) OR ((amount + 0) >= 400))
Planning Time: 0.060 ms
Execution Time: 6.249 ms
Peak Memory Usage: 548 kB

```

### s20_union_predicate [shape=expression;spelling=union]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.073000 / 8.326000 ms.

```text
Sort  (cost=2516.53..2536.53 rows=8000 width=8) (actual time=7.785..7.921 rows=4652 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 411kB
  ->  HashAggregate  (cost=1917.90..1997.90 rows=8000 width=8) (actual time=6.724..7.040 rows=4652 loops=1)
        Group Key: orders.id, orders.amount
        Batches: 1  Memory Usage: 657kB
        ->  Append  (cost=21.08..1877.90 rows=8000 width=8) (actual time=1.557..5.963 rows=4652 loops=1)
              ->  Seq Scan on orders  (cost=21.08..878.95 rows=4000 width=8) (actual time=1.555..2.899 rows=2273 loops=1)
                    Storage Filter: ((amount + 0) <= 100)
              ->  Seq Scan on orders orders_1  (cost=21.08..878.95 rows=4000 width=8) (actual time=1.345..2.609 rows=2379 loops=1)
                    Storage Filter: ((amount + 0) >= 400)
Planning Time: 0.073 ms
Execution Time: 8.326 ms
Peak Memory Usage: 1252 kB

```

### s19_union_all_predicate [shape=expression;spelling=union_all]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.090000 / 8.218000 ms.

```text
Sort  (cost=2396.53..2416.53 rows=8000 width=8) (actual time=7.709..7.860 rows=4652 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 411kB
  ->  Append  (cost=21.08..1877.90 rows=8000 width=8) (actual time=1.683..6.945 rows=4652 loops=1)
        ->  Seq Scan on orders  (cost=21.08..878.95 rows=4000 width=8) (actual time=1.682..3.516 rows=2273 loops=1)
              Storage Filter: ((amount + 0) <= 100)
        ->  Seq Scan on orders orders_1  (cost=21.08..878.95 rows=4000 width=8) (actual time=1.599..2.961 rows=2379 loops=1)
              Storage Filter: ((amount + 0) >= 400)
Planning Time: 0.090 ms
Execution Time: 8.218 ms
Peak Memory Usage: 625 kB

```

### s18_or_predicate [shape=indexable;spelling=or]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.065000 / 4.789000 ms.

```text
Sort  (cost=1154.19..1164.72 rows=4211 width=8) (actual time=4.251..4.402 rows=4652 loops=1)
  Sort Key: id
  Sort Method: quicksort  Memory: 411kB
  ->  Seq Scan on orders  (cost=21.08..900.69 rows=4211 width=8) (actual time=1.061..3.471 rows=4652 loops=1)
        Storage Filter: ((amount <= 100) OR (amount >= 400))
Planning Time: 0.065 ms
Execution Time: 4.789 ms
Peak Memory Usage: 548 kB

```

### s20_union_predicate [shape=indexable;spelling=union]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.081000 / 7.871000 ms.

```text
Sort  (cost=2056.31..2067.97 rows=4664 width=8) (actual time=7.325..7.466 rows=4652 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 411kB
  ->  HashAggregate  (cost=1725.46..1772.10 rows=4664 width=8) (actual time=6.290..6.549 rows=4652 loops=1)
        Group Key: orders.id, orders.amount
        Batches: 1  Memory Usage: 465kB
        ->  Append  (cost=21.08..1702.14 rows=4664 width=8) (actual time=1.560..5.557 rows=4652 loops=1)
              ->  Seq Scan on orders  (cost=21.08..815.59 rows=2267 width=8) (actual time=1.559..2.645 rows=2273 loops=1)
                    Storage Filter: (amount <= 100)
              ->  Seq Scan on orders orders_1  (cost=21.08..816.59 rows=2397 width=8) (actual time=1.128..2.429 rows=2379 loops=1)
                    Storage Filter: (amount >= 400)
Planning Time: 0.081 ms
Execution Time: 7.871 ms
Peak Memory Usage: 1060 kB

```

### s19_union_all_predicate [shape=indexable;spelling=union_all]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.117000 / 7.174000 ms.

```text
Sort  (cost=1986.35..1998.01 rows=4664 width=8) (actual time=6.619..6.785 rows=4652 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 411kB
  ->  Append  (cost=21.08..1702.14 rows=4664 width=8) (actual time=1.545..5.853 rows=4652 loops=1)
        ->  Seq Scan on orders  (cost=21.08..815.59 rows=2267 width=8) (actual time=1.543..2.777 rows=2273 loops=1)
              Storage Filter: (amount <= 100)
        ->  Seq Scan on orders orders_1  (cost=21.08..816.59 rows=2397 width=8) (actual time=1.285..2.623 rows=2379 loops=1)
              Storage Filter: (amount >= 400)
Planning Time: 0.117 ms
Execution Time: 7.174 ms
Peak Memory Usage: 625 kB

```

## Limits and gaps recorded by the runner

- no concurrent writers: client-decode JSON bytes are not wire bytes; tails are under-sampled at the protocol trial count
- one host: no independent hardware, real-network, 3-node placement or equal-total-budget comparison is established
- single-node only; the YB 3-node RF=3 per-node/equal-total arms are a coverage gap on this host

