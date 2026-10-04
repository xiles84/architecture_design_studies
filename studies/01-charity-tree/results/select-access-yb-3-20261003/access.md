# SELECT access-path matrix (generated measurements)

_Measurements only. Verdicts for P10/P11 belong in a signed analysis._

**TL;DR:** within-run medians for OR / UNION ALL / UNION at three predicate shapes and OFFSET / keyset at five page depths, each cell verified against the independent oracle. A within-run difference is not a replicated result. Missing measurements are N/A.

| run | engine | topology | scale | page limit |
|---|---|---|---|---|
| select-access-yb-3-20261003 | yugabyte | single | deep | 5 |

## Provenance

| Field | Value |
|---|---|
| repo commit | `6ad01ef65352b0be194007dab460b297ae94a541` |
| benchmark image | `localhost/charitybench-select:select-access-yb-3-20261003` |
| image digest | `023c9fb2d709de026afbef825315bf2e29a234361a3d2a0f95507136c4fa57f5` |
| run tag | `run/study-01/select-access-yb-3-20261003` |
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
| selective/indexable | `amount <= 5` | `amount >= 496` | 114 | 137 | 0 | true |
| selective/expression | `(amount + 0) <= 5` | `(amount + 0) >= 496` | 114 | 137 | 0 | true |
| selective/cast | `(amount)::bigint <= 5` | `(amount)::bigint >= 496` | 114 | 137 | 0 | true |
| broad/indexable | `amount <= 100` | `amount >= 400` | 2273 | 2379 | 0 | true |
| broad/expression | `(amount + 0) <= 100` | `(amount + 0) >= 400` | 2273 | 2379 | 0 | true |
| broad/cast | `(amount)::bigint <= 100` | `(amount)::bigint >= 400` | 2273 | 2379 | 0 | true |

## Predicate alternatives (P10)

| shape | index arm | spelling | status | median ms | spread % | trials | rows | plan structure |
|---|---|---|---|---|---|---|---|---|
| broad/cast | default | or | valid | 16.740603 | 10.7 | 5 | 4652 | `1dd363b32be056fc` |
| broad/cast | default | union | valid | 25.919721 | 9.1 | 5 | 4652 | `b5f6c293b9267d60` |
| broad/cast | default | union_all | valid | 23.939678 | 4.1 | 5 | 4652 | `24cabca42adfa691` |
| broad/expression | expression_index | or | valid | 19.239256 | 7.5 | 5 | 4652 | `0e731990206e1446` |
| broad/expression | expression_index | union | valid | 26.107827 | 10.7 | 5 | 4652 | `bf74591baa836f9d` |
| broad/expression | expression_index | union_all | valid | 25.699853 | 17.3 | 5 | 4652 | `a7e687668ac364de` |
| broad/expression | default | or | valid | 16.901565 | 17.4 | 5 | 4652 | `0e731990206e1446` |
| broad/expression | default | union | valid | 27.355109 | 61.3 | 5 | 4652 | `bf74591baa836f9d` |
| broad/expression | default | union_all | valid | 30.641029 | 92.3 | 5 | 4652 | `a7e687668ac364de` |
| broad/indexable | default | or | valid | 19.663053 | 33.5 | 5 | 4652 | `5191eeb70daef797` |
| broad/indexable | default | union | valid | 25.104787 | 17.1 | 5 | 4652 | `12de8f648f9520a0` |
| broad/indexable | default | union_all | valid | 22.170745 | 11.8 | 5 | 4652 | `d95e057d512f5ae5` |
| selective/cast | default | or | valid | 11.706365 | 10.7 | 5 | 251 | `57f71a4ae173cc97` |
| selective/cast | default | union | valid | 17.062551 | 24.8 | 5 | 251 | `2a31413606cd52a3` |
| selective/cast | default | union_all | valid | 16.400881 | 13.4 | 5 | 251 | `096e4d4e43a2a56e` |
| selective/expression | expression_index | or | valid | 9.738460 | 7.5 | 5 | 251 | `66a4eb5992604842` |
| selective/expression | expression_index | union | valid | 16.287180 | 10.3 | 5 | 251 | `14832cc886f611bd` |
| selective/expression | expression_index | union_all | valid | 15.602001 | 12.6 | 5 | 251 | `fee9cce2de640d1c` |
| selective/expression | default | or | valid | 12.202750 | 13.4 | 5 | 251 | `66a4eb5992604842` |
| selective/expression | default | union | valid | 19.574730 | 21.0 | 5 | 251 | `14832cc886f611bd` |
| selective/expression | default | union_all | valid | 21.653088 | 16.5 | 5 | 251 | `fee9cce2de640d1c` |
| selective/indexable | default | or | valid | 8.586247 | 13.3 | 5 | 251 | `66fcedd6207432cd` |
| selective/indexable | default | union | valid | 14.989680 | 7.4 | 5 | 251 | `40578cd947e4b93c` |
| selective/indexable | default | union_all | valid | 14.813902 | 25.3 | 5 | 251 | `8136cf54e3f9a863` |

### P10 pairs

Pairs compare the **default** arm only; the labelled `expression_index` arm is a package arm and is never pooled with the query-text pair.

| pair | A median ms | B median ms | B/A | same plan | within-run verdict |
|---|---|---|---|---|---|
| shape=selective/indexable or vs union_all | 8.586247 | 14.813902 | 1.7253 | false | within-run medians differ; whole-run replication required before a verdict |
| shape=selective/indexable or vs union | 8.586247 | 14.989680 | 1.7458 | false | within-run medians differ; whole-run replication required before a verdict |
| shape=selective/expression or vs union_all | 12.202750 | 21.653088 | 1.7744 | false | within-run medians differ; whole-run replication required before a verdict |
| shape=selective/expression or vs union | 12.202750 | 19.574730 | 1.6041 | false | within-run medians differ; whole-run replication required before a verdict |
| shape=selective/cast or vs union_all | 11.706365 | 16.400881 | 1.4010 | false | within-run medians differ; whole-run replication required before a verdict |
| shape=selective/cast or vs union | 11.706365 | 17.062551 | 1.4575 | false | within-run medians differ; whole-run replication required before a verdict |
| shape=broad/indexable or vs union_all | 19.663053 | 22.170745 | 1.1275 | false | medians inside the two cells' own spread; no resolved difference within this run |
| shape=broad/indexable or vs union | 19.663053 | 25.104787 | 1.2767 | false | medians inside the two cells' own spread; no resolved difference within this run |
| shape=broad/expression or vs union_all | 16.901565 | 30.641029 | 1.8129 | false | medians inside the two cells' own spread; no resolved difference within this run |
| shape=broad/expression or vs union | 16.901565 | 27.355109 | 1.6185 | false | medians inside the two cells' own spread; no resolved difference within this run |
| shape=broad/cast or vs union_all | 16.740603 | 23.939678 | 1.4300 | false | within-run medians differ; whole-run replication required before a verdict |
| shape=broad/cast or vs union | 16.740603 | 25.919721 | 1.5483 | false | within-run medians differ; whole-run replication required before a verdict |

## Pagination (P11)

| depth | spelling | status | median ms | spread % | trials | rows | plan structure |
|---|---|---|---|---|---|---|---|
| 1 | offset | valid | 14.194939 | 12.2 | 5 | 5 | `ac4509d295bf3348` |
| 1 | keyset | valid | 16.705487 | 20.4 | 5 | 5 | `166777c9921144af` |
| 10 | offset | valid | 13.942363 | 9.6 | 5 | 5 | `ac4509d295bf3348` |
| 10 | keyset | valid | 16.486891 | 10.8 | 5 | 5 | `d596ca0d42dc4da7` |
| 100 | offset | valid | 14.058823 | 28.1 | 5 | 5 | `ac4509d295bf3348` |
| 100 | keyset | valid | 17.572849 | 4.1 | 5 | 5 | `727a2962d7392ccd` |
| 1000 | offset | valid | 12.319564 | 14.1 | 5 | 5 | `ac4509d295bf3348` |
| 1000 | keyset | valid | 16.402403 | 14.9 | 5 | 5 | `441b4191cb138ed5` |
| 10000 | offset | valid | 16.887356 | 9.0 | 5 | 5 | `ac4509d295bf3348` |
| 10000 | keyset | valid | 11.682140 | 5.9 | 5 | 5 | `e1c72b3614c26105` |

Median ratio table (keyset / OFFSET; a ratio below 1 means keyset was faster in this run):

| depth | ratio (keyset/offset) | offset median ms | keyset median ms |
|---|---|---|---|
| 1 | 1.1769 | 14.194939 | 16.705487 |
| 10 | 1.1825 | 13.942363 | 16.486891 |
| 100 | 1.2500 | 14.058823 | 17.572849 |
| 1000 | 1.3314 | 12.319564 | 16.402403 |
| 10000 | 0.6918 | 16.887356 | 11.682140 |

### P11 pairs

| pair | A median ms | B median ms | B/A | same plan | within-run verdict |
|---|---|---|---|---|---|
| depth=1 offset vs keyset | 14.194939 | 16.705487 | 1.1769 | false | medians inside the two cells' own spread; no resolved difference within this run |
| depth=10 offset vs keyset | 13.942363 | 16.486891 | 1.1825 | false | medians inside the two cells' own spread; no resolved difference within this run |
| depth=100 offset vs keyset | 14.058823 | 17.572849 | 1.2500 | false | medians inside the two cells' own spread; no resolved difference within this run |
| depth=1000 offset vs keyset | 12.319564 | 16.402403 | 1.3314 | false | within-run medians differ; whole-run replication required before a verdict |
| depth=10000 offset vs keyset | 16.887356 | 11.682140 | 0.6918 | false | within-run medians differ; whole-run replication required before a verdict |

## Repeated-design position control

Not measured in this run.

## Expression-index arm: write and storage cost

Index `orders_amount_expr_idx` on expression `(amount + 0)`; supported=true.

| measure | value |
|---|---|
| index size (bytes) | 0 |
| insert batch with index (ms) | 23.529761 |
| insert batch without index (ms) | 62.053034 |
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

- shape=selective/indexable `x01_optional_where`: detected=false, rows=9175 — x01_optional_where was NOT detected: its answer matched the sel_optional_match oracle
- shape=selective/indexable `x02_not_in_nullable`: detected=false, rows=0 — x02_not_in_nullable was NOT detected: its answer matched the sel_absence oracle
- shape=selective/indexable `x03_distinct_on_name`: detected=false, rows=200 — x03_distinct_on_name was NOT detected: its answer matched the sel_membership oracle
- shape=selective/indexable `x04_limit_before_join`: detected=true, rows=5
- shape=selective/indexable `x05_union_all_overlap`: detected=true, rows=228
- shape=selective/expression `x01_optional_where`: detected=false, rows=9175 — x01_optional_where was NOT detected: its answer matched the sel_optional_match oracle
- shape=selective/expression `x02_not_in_nullable`: detected=false, rows=0 — x02_not_in_nullable was NOT detected: its answer matched the sel_absence oracle
- shape=selective/expression `x03_distinct_on_name`: detected=false, rows=200 — x03_distinct_on_name was NOT detected: its answer matched the sel_membership oracle
- shape=selective/expression `x04_limit_before_join`: detected=true, rows=5
- shape=selective/expression `x05_union_all_overlap`: detected=true, rows=228
- shape=selective/cast `x01_optional_where`: detected=false, rows=9175 — x01_optional_where was NOT detected: its answer matched the sel_optional_match oracle
- shape=selective/cast `x02_not_in_nullable`: detected=false, rows=0 — x02_not_in_nullable was NOT detected: its answer matched the sel_absence oracle
- shape=selective/cast `x03_distinct_on_name`: detected=false, rows=200 — x03_distinct_on_name was NOT detected: its answer matched the sel_membership oracle
- shape=selective/cast `x04_limit_before_join`: detected=true, rows=5
- shape=selective/cast `x05_union_all_overlap`: detected=true, rows=228
- shape=broad/indexable `x01_optional_where`: detected=false, rows=9175 — x01_optional_where was NOT detected: its answer matched the sel_optional_match oracle
- shape=broad/indexable `x02_not_in_nullable`: detected=false, rows=0 — x02_not_in_nullable was NOT detected: its answer matched the sel_absence oracle
- shape=broad/indexable `x03_distinct_on_name`: detected=false, rows=200 — x03_distinct_on_name was NOT detected: its answer matched the sel_membership oracle
- shape=broad/indexable `x04_limit_before_join`: detected=true, rows=5
- shape=broad/indexable `x05_union_all_overlap`: detected=true, rows=4546
- shape=broad/expression `x01_optional_where`: detected=false, rows=9175 — x01_optional_where was NOT detected: its answer matched the sel_optional_match oracle
- shape=broad/expression `x02_not_in_nullable`: detected=false, rows=0 — x02_not_in_nullable was NOT detected: its answer matched the sel_absence oracle
- shape=broad/expression `x03_distinct_on_name`: detected=false, rows=200 — x03_distinct_on_name was NOT detected: its answer matched the sel_membership oracle
- shape=broad/expression `x04_limit_before_join`: detected=true, rows=5
- shape=broad/expression `x05_union_all_overlap`: detected=true, rows=4546
- shape=broad/cast `x01_optional_where`: detected=false, rows=9175 — x01_optional_where was NOT detected: its answer matched the sel_optional_match oracle
- shape=broad/cast `x02_not_in_nullable`: detected=false, rows=0 — x02_not_in_nullable was NOT detected: its answer matched the sel_absence oracle
- shape=broad/cast `x03_distinct_on_name`: detected=false, rows=200 — x03_distinct_on_name was NOT detected: its answer matched the sel_membership oracle
- shape=broad/cast `x04_limit_before_join`: detected=true, rows=5
- shape=broad/cast `x05_union_all_overlap`: detected=true, rows=4546

## Excluded cells

None: every cell is valid.

## Captured plans and result sizes

### s17_page_keyset [depth=10000;spelling=keyset]

Rows: 5; serialized decoded-row JSON bytes: 51 (not wire bytes). Planning/execution: 0.154000 / 10.590000 ms.

```text
Limit  (cost=884.51..884.52 rows=5 width=8) (actual time=10.496..10.500 rows=5 loops=1)
  ->  Sort  (cost=884.51..889.29 rows=1912 width=8) (actual time=10.492..10.494 rows=5 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 25kB
        ->  Seq Scan on orders  (cost=21.08..852.75 rows=1912 width=8) (actual time=4.940..9.837 rows=2000 loops=1)
              Storage Filter: ((amount < 67) OR ((amount = 67) AND (id > 4310)) OR (amount IS NULL))
Planning Time: 0.154 ms
Execution Time: 10.590 ms
Peak Memory Usage: 92 kB

```

### s16_page_offset [depth=10000;spelling=offset]

Rows: 5; serialized decoded-row JSON bytes: 51 (not wire bytes). Planning/execution: 0.134000 / 20.922000 ms.

```text
Limit  (cost=1879.09..1879.10 rows=5 width=8) (actual time=20.800..20.807 rows=5 loops=1)
  ->  Sort  (cost=1854.09..1884.09 rows=12000 width=8) (actual time=16.996..18.330 rows=10005 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: quicksort  Memory: 947kB
        ->  Seq Scan on orders  (cost=21.08..1041.05 rows=12000 width=8) (actual time=1.457..8.689 rows=12000 loops=1)
Planning Time: 0.134 ms
Execution Time: 20.922 ms
Peak Memory Usage: 1516 kB

```

### s17_page_keyset [depth=1000;spelling=keyset]

Rows: 5; serialized decoded-row JSON bytes: 59 (not wire bytes). Planning/execution: 0.160000 / 14.290000 ms.

```text
Limit  (cost=1274.79..1274.81 rows=5 width=8) (actual time=14.210..14.214 rows=5 loops=1)
  ->  Sort  (cost=1274.79..1301.14 rows=10538 width=8) (actual time=14.207..14.209 rows=5 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 25kB
        ->  Seq Scan on orders  (cost=21.08..1099.76 rows=10538 width=8) (actual time=1.818..11.139 rows=11000 loops=1)
              Storage Filter: ((amount < 457) OR ((amount = 457) AND (id > 11254)) OR (amount IS NULL))
Planning Time: 0.160 ms
Execution Time: 14.290 ms
Peak Memory Usage: 92 kB

```

### s16_page_offset [depth=1000;spelling=offset]

Rows: 5; serialized decoded-row JSON bytes: 59 (not wire bytes). Planning/execution: 0.143000 / 12.923000 ms.

```text
Limit  (cost=1701.93..1701.94 rows=5 width=8) (actual time=12.831..12.835 rows=5 loops=1)
  ->  Sort  (cost=1699.43..1729.43 rows=12000 width=8) (actual time=12.568..12.638 rows=1005 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 130kB
        ->  Seq Scan on orders  (cost=21.08..1041.05 rows=12000 width=8) (actual time=1.688..7.639 rows=12000 loops=1)
Planning Time: 0.143 ms
Execution Time: 12.923 ms
Peak Memory Usage: 276 kB

```

### s17_page_keyset [depth=100;spelling=keyset]

Rows: 5; serialized decoded-row JSON bytes: 56 (not wire bytes). Planning/execution: 0.180000 / 15.563000 ms.

```text
Limit  (cost=1315.00..1315.01 rows=5 width=8) (actual time=15.481..15.485 rows=5 loops=1)
  ->  Sort  (cost=1315.00..1343.41 rows=11366 width=8) (actual time=15.478..15.480 rows=5 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 25kB
        ->  Seq Scan on orders  (cost=21.08..1126.21 rows=11366 width=8) (actual time=1.976..11.874 rows=11900 loops=1)
              Storage Filter: ((amount < 497) OR ((amount = 497) AND (id > 6723)) OR (amount IS NULL))
Planning Time: 0.180 ms
Execution Time: 15.563 ms
Peak Memory Usage: 92 kB

```

### s16_page_offset [depth=100;spelling=offset]

Rows: 5; serialized decoded-row JSON bytes: 56 (not wire bytes). Planning/execution: 0.104000 / 11.431000 ms.

```text
Limit  (cost=1504.15..1504.16 rows=5 width=8) (actual time=11.321..11.325 rows=5 loops=1)
  ->  Sort  (cost=1503.90..1533.90 rows=12000 width=8) (actual time=11.292..11.300 rows=105 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 33kB
        ->  Seq Scan on orders  (cost=21.08..1041.05 rows=12000 width=8) (actual time=1.216..7.468 rows=12000 loops=1)
Planning Time: 0.104 ms
Execution Time: 11.431 ms
Peak Memory Usage: 128 kB

```

### s17_page_keyset [depth=10;spelling=keyset]

Rows: 5; serialized decoded-row JSON bytes: 56 (not wire bytes). Planning/execution: 0.170000 / 15.964000 ms.

```text
Limit  (cost=1317.54..1317.55 rows=5 width=8) (actual time=15.870..15.874 rows=5 loops=1)
  ->  Sort  (cost=1317.54..1346.22 rows=11471 width=8) (actual time=15.867..15.868 rows=5 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 25kB
        ->  Seq Scan on orders  (cost=21.08..1127.01 rows=11471 width=8) (actual time=1.953..12.398 rows=11990 loops=1)
              Storage Filter: ((amount < 500) OR ((amount = 500) AND (id > 6319)) OR (amount IS NULL))
Planning Time: 0.170 ms
Execution Time: 15.964 ms
Peak Memory Usage: 92 kB

```

### s16_page_offset [depth=10;spelling=offset]

Rows: 5; serialized decoded-row JSON bytes: 56 (not wire bytes). Planning/execution: 0.194000 / 13.357000 ms.

```text
Limit  (cost=1335.49..1335.50 rows=5 width=8) (actual time=13.257..13.261 rows=5 loops=1)
  ->  Sort  (cost=1335.46..1365.46 rows=12000 width=8) (actual time=13.250..13.253 rows=15 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 26kB
        ->  Seq Scan on orders  (cost=21.08..1041.05 rows=12000 width=8) (actual time=1.722..9.491 rows=12000 loops=1)
Planning Time: 0.194 ms
Execution Time: 13.357 ms
Peak Memory Usage: 92 kB

```

### s17_page_keyset [depth=1;spelling=keyset]

Rows: 5; serialized decoded-row JSON bytes: 56 (not wire bytes). Planning/execution: 0.157000 / 17.043000 ms.

```text
Limit  (cost=1317.57..1317.58 rows=5 width=8) (actual time=16.945..16.949 rows=5 loops=1)
  ->  Sort  (cost=1317.57..1346.25 rows=11472 width=8) (actual time=16.942..16.943 rows=5 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 25kB
        ->  Seq Scan on orders  (cost=21.08..1127.02 rows=11472 width=8) (actual time=1.719..13.524 rows=11999 loops=1)
              Storage Filter: ((amount < 500) OR ((amount = 500) AND (id > 1554)) OR (amount IS NULL))
Planning Time: 0.157 ms
Execution Time: 17.043 ms
Peak Memory Usage: 92 kB

```

### s16_page_offset [depth=1;spelling=offset]

Rows: 5; serialized decoded-row JSON bytes: 56 (not wire bytes). Planning/execution: 0.168000 / 14.723000 ms.

```text
Limit  (cost=1256.15..1256.16 rows=5 width=8) (actual time=14.539..14.543 rows=5 loops=1)
  ->  Sort  (cost=1256.15..1286.15 rows=12000 width=8) (actual time=14.536..14.538 rows=6 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 25kB
        ->  Seq Scan on orders  (cost=21.08..1041.05 rows=12000 width=8) (actual time=2.004..10.946 rows=12000 loops=1)
Planning Time: 0.168 ms
Execution Time: 14.723 ms
Peak Memory Usage: 92 kB

```

### s18_or_predicate [shape=broad/cast;spelling=or]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.139000 / 14.872000 ms.

```text
Sort  (cost=1443.14..1459.81 rows=6667 width=8) (actual time=13.341..13.768 rows=4652 loops=1)
  Sort Key: id
  Sort Method: quicksort  Memory: 411kB
  ->  Seq Scan on orders  (cost=21.08..1019.70 rows=6667 width=8) (actual time=2.943..10.862 rows=4652 loops=1)
        Storage Filter: (((amount)::bigint <= 100) OR ((amount)::bigint >= 400))
Planning Time: 0.139 ms
Execution Time: 14.872 ms
Peak Memory Usage: 548 kB

```

### s20_union_predicate [shape=broad/cast;spelling=union]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.323000 / 24.902000 ms.

```text
Sort  (cost=2516.53..2536.53 rows=8000 width=8) (actual time=23.229..23.598 rows=4652 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 411kB
  ->  HashAggregate  (cost=1917.90..1997.90 rows=8000 width=8) (actual time=20.274..21.156 rows=4652 loops=1)
        Group Key: orders.id, orders.amount
        Batches: 1  Memory Usage: 657kB
        ->  Append  (cost=21.08..1877.90 rows=8000 width=8) (actual time=5.320..17.949 rows=4652 loops=1)
              ->  Seq Scan on orders  (cost=21.08..878.95 rows=4000 width=8) (actual time=5.315..9.130 rows=2273 loops=1)
                    Storage Filter: ((amount)::bigint <= 100)
              ->  Seq Scan on orders orders_1  (cost=21.08..878.95 rows=4000 width=8) (actual time=3.894..7.530 rows=2379 loops=1)
                    Storage Filter: ((amount)::bigint >= 400)
Planning Time: 0.323 ms
Execution Time: 24.902 ms
Peak Memory Usage: 1252 kB

```

### s19_union_all_predicate [shape=broad/cast;spelling=union_all]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.200000 / 23.776000 ms.

```text
Sort  (cost=2396.53..2416.53 rows=8000 width=8) (actual time=21.306..22.434 rows=4652 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 411kB
  ->  Append  (cost=21.08..1877.90 rows=8000 width=8) (actual time=5.465..18.902 rows=4652 loops=1)
        ->  Seq Scan on orders  (cost=21.08..878.95 rows=4000 width=8) (actual time=5.458..9.656 rows=2273 loops=1)
              Storage Filter: ((amount)::bigint <= 100)
        ->  Seq Scan on orders orders_1  (cost=21.08..878.95 rows=4000 width=8) (actual time=3.909..7.802 rows=2379 loops=1)
              Storage Filter: ((amount)::bigint >= 400)
Planning Time: 0.200 ms
Execution Time: 23.776 ms
Peak Memory Usage: 625 kB

```

### s18_or_predicate [shape=broad/expression;index=expression_index;spelling=or]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.162000 / 16.164000 ms.

```text
Sort  (cost=1214.19..1224.72 rows=4211 width=8) (actual time=14.355..14.943 rows=4652 loops=1)
  Sort Key: id
  Sort Method: quicksort  Memory: 411kB
  ->  Seq Scan on orders  (cost=21.08..960.69 rows=4211 width=8) (actual time=2.878..11.570 rows=4652 loops=1)
        Storage Filter: (((amount + 0) <= 100) OR ((amount + 0) >= 400))
Planning Time: 0.162 ms
Execution Time: 16.164 ms
Peak Memory Usage: 548 kB

```

### s20_union_predicate [shape=broad/expression;index=expression_index;spelling=union]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.184000 / 25.454000 ms.

```text
Sort  (cost=2116.31..2127.97 rows=4664 width=8) (actual time=23.709..24.161 rows=4652 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 411kB
  ->  HashAggregate  (cost=1785.46..1832.10 rows=4664 width=8) (actual time=20.533..21.346 rows=4652 loops=1)
        Group Key: orders.id, orders.amount
        Batches: 1  Memory Usage: 465kB
        ->  Append  (cost=21.08..1762.14 rows=4664 width=8) (actual time=4.306..18.088 rows=4652 loops=1)
              ->  Seq Scan on orders  (cost=21.08..845.59 rows=2267 width=8) (actual time=4.303..8.171 rows=2273 loops=1)
                    Storage Filter: ((amount + 0) <= 100)
              ->  Seq Scan on orders orders_1  (cost=21.08..846.59 rows=2397 width=8) (actual time=3.901..8.378 rows=2379 loops=1)
                    Storage Filter: ((amount + 0) >= 400)
Planning Time: 0.184 ms
Execution Time: 25.454 ms
Peak Memory Usage: 1060 kB

```

### s19_union_all_predicate [shape=broad/expression;index=expression_index;spelling=union_all]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.223000 / 23.696000 ms.

```text
Sort  (cost=2046.35..2058.01 rows=4664 width=8) (actual time=21.715..22.283 rows=4652 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 411kB
  ->  Append  (cost=21.08..1762.14 rows=4664 width=8) (actual time=4.581..19.179 rows=4652 loops=1)
        ->  Seq Scan on orders  (cost=21.08..845.59 rows=2267 width=8) (actual time=4.577..9.125 rows=2273 loops=1)
              Storage Filter: ((amount + 0) <= 100)
        ->  Seq Scan on orders orders_1  (cost=21.08..846.59 rows=2397 width=8) (actual time=4.128..8.478 rows=2379 loops=1)
              Storage Filter: ((amount + 0) >= 400)
Planning Time: 0.223 ms
Execution Time: 23.696 ms
Peak Memory Usage: 625 kB

```

### s18_or_predicate [shape=broad/expression;spelling=or]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.134000 / 17.881000 ms.

```text
Sort  (cost=1443.14..1459.81 rows=6667 width=8) (actual time=16.341..16.816 rows=4652 loops=1)
  Sort Key: id
  Sort Method: quicksort  Memory: 411kB
  ->  Seq Scan on orders  (cost=21.08..1019.70 rows=6667 width=8) (actual time=3.722..14.094 rows=4652 loops=1)
        Storage Filter: (((amount + 0) <= 100) OR ((amount + 0) >= 400))
Planning Time: 0.134 ms
Execution Time: 17.881 ms
Peak Memory Usage: 548 kB

```

### s20_union_predicate [shape=broad/expression;spelling=union]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.210000 / 30.315000 ms.

```text
Sort  (cost=2516.53..2536.53 rows=8000 width=8) (actual time=26.070..27.571 rows=4652 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 411kB
  ->  HashAggregate  (cost=1917.90..1997.90 rows=8000 width=8) (actual time=22.083..23.433 rows=4652 loops=1)
        Group Key: orders.id, orders.amount
        Batches: 1  Memory Usage: 657kB
        ->  Append  (cost=21.08..1877.90 rows=8000 width=8) (actual time=4.229..18.762 rows=4652 loops=1)
              ->  Seq Scan on orders  (cost=21.08..878.95 rows=4000 width=8) (actual time=4.226..8.356 rows=2273 loops=1)
                    Storage Filter: ((amount + 0) <= 100)
              ->  Seq Scan on orders orders_1  (cost=21.08..878.95 rows=4000 width=8) (actual time=4.480..8.853 rows=2379 loops=1)
                    Storage Filter: ((amount + 0) >= 400)
Planning Time: 0.210 ms
Execution Time: 30.315 ms
Peak Memory Usage: 1252 kB

```

### s19_union_all_predicate [shape=broad/expression;spelling=union_all]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.193000 / 24.956000 ms.

```text
Sort  (cost=2396.53..2416.53 rows=8000 width=8) (actual time=22.993..23.576 rows=4652 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 411kB
  ->  Append  (cost=21.08..1877.90 rows=8000 width=8) (actual time=4.598..16.561 rows=4652 loops=1)
        ->  Seq Scan on orders  (cost=21.08..878.95 rows=4000 width=8) (actual time=4.594..8.788 rows=2273 loops=1)
              Storage Filter: ((amount + 0) <= 100)
        ->  Seq Scan on orders orders_1  (cost=21.08..878.95 rows=4000 width=8) (actual time=4.968..6.257 rows=2379 loops=1)
              Storage Filter: ((amount + 0) >= 400)
Planning Time: 0.193 ms
Execution Time: 24.956 ms
Peak Memory Usage: 625 kB

```

### s18_or_predicate [shape=broad/indexable;spelling=or]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.139000 / 15.078000 ms.

```text
Sort  (cost=1154.19..1164.72 rows=4211 width=8) (actual time=11.969..12.765 rows=4652 loops=1)
  Sort Key: id
  Sort Method: quicksort  Memory: 411kB
  ->  Seq Scan on orders  (cost=21.08..900.69 rows=4211 width=8) (actual time=3.027..9.241 rows=4652 loops=1)
        Storage Filter: ((amount <= 100) OR (amount >= 400))
Planning Time: 0.139 ms
Execution Time: 15.078 ms
Peak Memory Usage: 548 kB

```

### s20_union_predicate [shape=broad/indexable;spelling=union]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.160000 / 22.530000 ms.

```text
Sort  (cost=2056.31..2067.97 rows=4664 width=8) (actual time=20.737..21.192 rows=4652 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 411kB
  ->  HashAggregate  (cost=1725.46..1772.10 rows=4664 width=8) (actual time=17.797..18.599 rows=4652 loops=1)
        Group Key: orders.id, orders.amount
        Batches: 1  Memory Usage: 465kB
        ->  Append  (cost=21.08..1702.14 rows=4664 width=8) (actual time=4.147..15.639 rows=4652 loops=1)
              ->  Seq Scan on orders  (cost=21.08..815.59 rows=2267 width=8) (actual time=4.143..6.991 rows=2273 loops=1)
                    Storage Filter: (amount <= 100)
              ->  Seq Scan on orders orders_1  (cost=21.08..816.59 rows=2397 width=8) (actual time=3.477..7.190 rows=2379 loops=1)
                    Storage Filter: (amount >= 400)
Planning Time: 0.160 ms
Execution Time: 22.530 ms
Peak Memory Usage: 1060 kB

```

### s19_union_all_predicate [shape=broad/indexable;spelling=union_all]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.231000 / 20.055000 ms.

```text
Sort  (cost=1986.35..1998.01 rows=4664 width=8) (actual time=18.339..18.924 rows=4652 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 411kB
  ->  Append  (cost=21.08..1702.14 rows=4664 width=8) (actual time=3.919..16.096 rows=4652 loops=1)
        ->  Seq Scan on orders  (cost=21.08..815.59 rows=2267 width=8) (actual time=3.916..7.946 rows=2273 loops=1)
              Storage Filter: (amount <= 100)
        ->  Seq Scan on orders orders_1  (cost=21.08..816.59 rows=2397 width=8) (actual time=3.209..6.501 rows=2379 loops=1)
              Storage Filter: (amount >= 400)
Planning Time: 0.231 ms
Execution Time: 20.055 ms
Peak Memory Usage: 625 kB

```

### s18_or_predicate [shape=selective/cast;spelling=or]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.139000 / 9.537000 ms.

```text
Sort  (cost=1443.14..1459.81 rows=6667 width=8) (actual time=9.381..9.402 rows=251 loops=1)
  Sort Key: id
  Sort Method: quicksort  Memory: 36kB
  ->  Seq Scan on orders  (cost=21.08..1019.70 rows=6667 width=8) (actual time=9.104..9.175 rows=251 loops=1)
        Storage Filter: (((amount)::bigint <= 5) OR ((amount)::bigint >= 496))
Planning Time: 0.139 ms
Execution Time: 9.537 ms
Peak Memory Usage: 125 kB

```

### s20_union_predicate [shape=selective/cast;spelling=union]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.162000 / 14.943000 ms.

```text
Sort  (cost=2516.53..2536.53 rows=8000 width=8) (actual time=14.738..14.757 rows=251 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 36kB
  ->  HashAggregate  (cost=1917.90..1997.90 rows=8000 width=8) (actual time=14.531..14.626 rows=251 loops=1)
        Group Key: orders.id, orders.amount
        Batches: 1  Memory Usage: 417kB
        ->  Append  (cost=21.08..1877.90 rows=8000 width=8) (actual time=6.689..14.288 rows=251 loops=1)
              ->  Seq Scan on orders  (cost=21.08..878.95 rows=4000 width=8) (actual time=6.685..6.712 rows=114 loops=1)
                    Storage Filter: ((amount)::bigint <= 5)
              ->  Seq Scan on orders orders_1  (cost=21.08..878.95 rows=4000 width=8) (actual time=7.448..7.485 rows=137 loops=1)
                    Storage Filter: ((amount)::bigint >= 496)
Planning Time: 0.162 ms
Execution Time: 14.943 ms
Peak Memory Usage: 600 kB

```

### s19_union_all_predicate [shape=selective/cast;spelling=union_all]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.201000 / 17.548000 ms.

```text
Sort  (cost=2396.53..2416.53 rows=8000 width=8) (actual time=17.297..17.314 rows=251 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 36kB
  ->  Append  (cost=21.08..1877.90 rows=8000 width=8) (actual time=7.985..17.145 rows=251 loops=1)
        ->  Seq Scan on orders  (cost=21.08..878.95 rows=4000 width=8) (actual time=7.982..8.007 rows=114 loops=1)
              Storage Filter: ((amount)::bigint <= 5)
        ->  Seq Scan on orders orders_1  (cost=21.08..878.95 rows=4000 width=8) (actual time=9.020..9.045 rows=137 loops=1)
              Storage Filter: ((amount)::bigint >= 496)
Planning Time: 0.201 ms
Execution Time: 17.548 ms
Peak Memory Usage: 202 kB

```

### s18_or_predicate [shape=selective/expression;index=expression_index;spelling=or]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.160000 / 9.551000 ms.

```text
Sort  (cost=860.56..861.21 rows=262 width=8) (actual time=9.365..9.388 rows=251 loops=1)
  Sort Key: id
  Sort Method: quicksort  Memory: 36kB
  ->  Seq Scan on orders  (cost=21.08..850.03 rows=262 width=8) (actual time=9.141..9.197 rows=251 loops=1)
        Storage Filter: (((amount + 0) <= 5) OR ((amount + 0) >= 496))
Planning Time: 0.160 ms
Execution Time: 9.551 ms
Peak Memory Usage: 125 kB

```

### s20_union_predicate [shape=selective/expression;index=expression_index;spelling=union]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.172000 / 14.707000 ms.

```text
Sort  (cost=1596.53..1597.19 rows=263 width=8) (actual time=14.528..14.550 rows=251 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 36kB
  ->  HashAggregate  (cost=1583.33..1585.96 rows=263 width=8) (actual time=14.332..14.380 rows=251 loops=1)
        Group Key: orders.id, orders.amount
        Batches: 1  Memory Usage: 45kB
        ->  Append  (cost=21.08..1582.02 rows=263 width=8) (actual time=7.056..14.181 rows=251 loops=1)
              ->  Seq Scan on orders  (cost=21.08..788.85 rows=107 width=8) (actual time=7.053..7.081 rows=114 loops=1)
                    Storage Filter: ((amount + 0) <= 5)
              ->  Seq Scan on orders orders_1  (cost=21.08..789.22 rows=156 width=8) (actual time=6.978..7.011 rows=137 loops=1)
                    Storage Filter: ((amount + 0) >= 496)
Planning Time: 0.172 ms
Execution Time: 14.707 ms
Peak Memory Usage: 224 kB

```

### s19_union_all_predicate [shape=selective/expression;index=expression_index;spelling=union_all]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.199000 / 14.427000 ms.

```text
Sort  (cost=1592.59..1593.25 rows=263 width=8) (actual time=14.248..14.296 rows=251 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 36kB
  ->  Append  (cost=21.08..1582.02 rows=263 width=8) (actual time=7.132..14.091 rows=251 loops=1)
        ->  Seq Scan on orders  (cost=21.08..788.85 rows=107 width=8) (actual time=7.128..7.153 rows=114 loops=1)
              Storage Filter: ((amount + 0) <= 5)
        ->  Seq Scan on orders orders_1  (cost=21.08..789.22 rows=156 width=8) (actual time=6.834..6.862 rows=137 loops=1)
              Storage Filter: ((amount + 0) >= 496)
Planning Time: 0.199 ms
Execution Time: 14.427 ms
Peak Memory Usage: 202 kB

```

### s18_or_predicate [shape=selective/expression;spelling=or]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.163000 / 10.181000 ms.

```text
Sort  (cost=1443.14..1459.81 rows=6667 width=8) (actual time=10.022..10.043 rows=251 loops=1)
  Sort Key: id
  Sort Method: quicksort  Memory: 36kB
  ->  Seq Scan on orders  (cost=21.08..1019.70 rows=6667 width=8) (actual time=9.708..9.798 rows=251 loops=1)
        Storage Filter: (((amount + 0) <= 5) OR ((amount + 0) >= 496))
Planning Time: 0.163 ms
Execution Time: 10.181 ms
Peak Memory Usage: 125 kB

```

### s20_union_predicate [shape=selective/expression;spelling=union]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.397000 / 20.761000 ms.

```text
Sort  (cost=2516.53..2536.53 rows=8000 width=8) (actual time=20.482..20.506 rows=251 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 36kB
  ->  HashAggregate  (cost=1917.90..1997.90 rows=8000 width=8) (actual time=20.041..20.212 rows=251 loops=1)
        Group Key: orders.id, orders.amount
        Batches: 1  Memory Usage: 417kB
        ->  Append  (cost=21.08..1877.90 rows=8000 width=8) (actual time=10.028..19.810 rows=251 loops=1)
              ->  Seq Scan on orders  (cost=21.08..878.95 rows=4000 width=8) (actual time=10.023..10.061 rows=114 loops=1)
                    Storage Filter: ((amount + 0) <= 5)
              ->  Seq Scan on orders orders_1  (cost=21.08..878.95 rows=4000 width=8) (actual time=9.592..9.633 rows=137 loops=1)
                    Storage Filter: ((amount + 0) >= 496)
Planning Time: 0.397 ms
Execution Time: 20.761 ms
Peak Memory Usage: 600 kB

```

### s19_union_all_predicate [shape=selective/expression;spelling=union_all]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.218000 / 19.595000 ms.

```text
Sort  (cost=2396.53..2416.53 rows=8000 width=8) (actual time=19.408..19.433 rows=251 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 36kB
  ->  Append  (cost=21.08..1877.90 rows=8000 width=8) (actual time=9.506..19.232 rows=251 loops=1)
        ->  Seq Scan on orders  (cost=21.08..878.95 rows=4000 width=8) (actual time=9.502..9.532 rows=114 loops=1)
              Storage Filter: ((amount + 0) <= 5)
        ->  Seq Scan on orders orders_1  (cost=21.08..878.95 rows=4000 width=8) (actual time=9.558..9.597 rows=137 loops=1)
              Storage Filter: ((amount + 0) >= 496)
Planning Time: 0.218 ms
Execution Time: 19.595 ms
Peak Memory Usage: 202 kB

```

### s18_or_predicate [shape=selective/indexable;spelling=or]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.239000 / 8.821000 ms.

```text
Sort  (cost=800.56..801.21 rows=262 width=8) (actual time=8.668..8.689 rows=251 loops=1)
  Sort Key: id
  Sort Method: quicksort  Memory: 36kB
  ->  Seq Scan on orders  (cost=21.08..790.03 rows=262 width=8) (actual time=8.458..8.509 rows=251 loops=1)
        Storage Filter: ((amount <= 5) OR (amount >= 496))
Planning Time: 0.239 ms
Execution Time: 8.821 ms
Peak Memory Usage: 125 kB

```

### s20_union_predicate [shape=selective/indexable;spelling=union]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.289000 / 16.528000 ms.

```text
Sort  (cost=1536.53..1537.19 rows=263 width=8) (actual time=16.310..16.339 rows=251 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 36kB
  ->  HashAggregate  (cost=1523.33..1525.96 rows=263 width=8) (actual time=15.988..16.049 rows=251 loops=1)
        Group Key: orders.id, orders.amount
        Batches: 1  Memory Usage: 45kB
        ->  Append  (cost=21.08..1522.02 rows=263 width=8) (actual time=7.870..15.660 rows=251 loops=1)
              ->  Seq Scan on orders  (cost=21.08..758.85 rows=107 width=8) (actual time=7.864..7.892 rows=114 loops=1)
                    Storage Filter: (amount <= 5)
              ->  Seq Scan on orders orders_1  (cost=21.08..759.22 rows=156 width=8) (actual time=7.648..7.682 rows=137 loops=1)
                    Storage Filter: (amount >= 496)
Planning Time: 0.289 ms
Execution Time: 16.528 ms
Peak Memory Usage: 224 kB

```

### s19_union_all_predicate [shape=selective/indexable;spelling=union_all]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.188000 / 13.999000 ms.

```text
Sort  (cost=1532.59..1533.25 rows=263 width=8) (actual time=13.813..13.835 rows=251 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 36kB
  ->  Append  (cost=21.08..1522.02 rows=263 width=8) (actual time=6.690..13.678 rows=251 loops=1)
        ->  Seq Scan on orders  (cost=21.08..758.85 rows=107 width=8) (actual time=6.687..6.713 rows=114 loops=1)
              Storage Filter: (amount <= 5)
        ->  Seq Scan on orders orders_1  (cost=21.08..759.22 rows=156 width=8) (actual time=6.864..6.890 rows=137 loops=1)
              Storage Filter: (amount >= 496)
Planning Time: 0.188 ms
Execution Time: 13.999 ms
Peak Memory Usage: 202 kB

```

## Limits and gaps recorded by the runner

- no concurrent writers: client-decode JSON bytes are not wire bytes; tails are under-sampled at the protocol trial count
- one host: no independent hardware, real-network, 3-node placement or equal-total-budget comparison is established
- single-node only; the YB 3-node RF=3 per-node/equal-total arms are a coverage gap on this host

