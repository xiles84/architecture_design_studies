# SELECT access-path matrix (generated measurements)

_Measurements only. Verdicts for P10/P11 belong in a signed analysis._

**TL;DR:** within-run medians for OR / UNION ALL / UNION at three predicate shapes and OFFSET / keyset at five page depths, each cell verified against the independent oracle. A within-run difference is not a replicated result. Missing measurements are N/A.

| run | engine | topology | scale | page limit |
|---|---|---|---|---|
| select-access-yb-1-20261003 | yugabyte | single | deep | 5 |

## Provenance

| Field | Value |
|---|---|
| repo commit | `6ad01ef65352b0be194007dab460b297ae94a541` |
| benchmark image | `localhost/charitybench-select:select-access-yb-1-20261003` |
| image digest | `023c9fb2d709de026afbef825315bf2e29a234361a3d2a0f95507136c4fa57f5` |
| run tag | `run/study-01/select-access-yb-1-20261003` |
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
| broad/cast | default | or | valid | 5.366463 | 15.2 | 5 | 4652 | `1dd363b32be056fc` |
| broad/cast | default | union | valid | 11.366105 | 16.4 | 5 | 4652 | `b5f6c293b9267d60` |
| broad/cast | default | union_all | valid | 11.271674 | 17.6 | 5 | 4652 | `24cabca42adfa691` |
| broad/expression | expression_index | or | valid | 7.570656 | 26.2 | 5 | 4652 | `0e731990206e1446` |
| broad/expression | expression_index | union | valid | 8.046960 | 14.2 | 5 | 4652 | `bf74591baa836f9d` |
| broad/expression | expression_index | union_all | valid | 7.561875 | 22.5 | 5 | 4652 | `a7e687668ac364de` |
| broad/expression | default | or | valid | 5.237082 | 7.6 | 5 | 4652 | `0e731990206e1446` |
| broad/expression | default | union | valid | 7.989638 | 3.7 | 5 | 4652 | `bf74591baa836f9d` |
| broad/expression | default | union_all | valid | 7.314392 | 4.6 | 5 | 4652 | `a7e687668ac364de` |
| broad/indexable | default | or | valid | 5.001992 | 7.1 | 5 | 4652 | `5191eeb70daef797` |
| broad/indexable | default | union | valid | 7.570294 | 23.9 | 5 | 4652 | `12de8f648f9520a0` |
| broad/indexable | default | union_all | valid | 7.029154 | 10.6 | 5 | 4652 | `d95e057d512f5ae5` |
| selective/cast | default | or | valid | 2.983885 | 34.4 | 5 | 251 | `57f71a4ae173cc97` |
| selective/cast | default | union | valid | 5.144890 | 17.4 | 5 | 251 | `2a31413606cd52a3` |
| selective/cast | default | union_all | valid | 5.063655 | 13.7 | 5 | 251 | `096e4d4e43a2a56e` |
| selective/expression | expression_index | or | valid | 3.323720 | 13.0 | 5 | 251 | `66a4eb5992604842` |
| selective/expression | expression_index | union | valid | 6.326547 | 20.0 | 5 | 251 | `14832cc886f611bd` |
| selective/expression | expression_index | union_all | valid | 6.220068 | 18.7 | 5 | 251 | `fee9cce2de640d1c` |
| selective/expression | default | or | valid | 3.143710 | 11.3 | 5 | 251 | `66a4eb5992604842` |
| selective/expression | default | union | valid | 4.946131 | 10.7 | 5 | 251 | `14832cc886f611bd` |
| selective/expression | default | union_all | valid | 5.245558 | 14.8 | 5 | 251 | `fee9cce2de640d1c` |
| selective/indexable | default | or | valid | 2.822762 | 26.8 | 5 | 251 | `66fcedd6207432cd` |
| selective/indexable | default | union | valid | 4.899867 | 14.4 | 5 | 251 | `40578cd947e4b93c` |
| selective/indexable | default | union_all | valid | 4.969804 | 41.5 | 5 | 251 | `8136cf54e3f9a863` |

### P10 pairs

Pairs compare the **default** arm only; the labelled `expression_index` arm is a package arm and is never pooled with the query-text pair.

| pair | A median ms | B median ms | B/A | same plan | within-run verdict |
|---|---|---|---|---|---|
| shape=selective/indexable or vs union_all | 2.822762 | 4.969804 | 1.7606 | false | within-run medians differ; whole-run replication required before a verdict |
| shape=selective/indexable or vs union | 2.822762 | 4.899867 | 1.7358 | false | within-run medians differ; whole-run replication required before a verdict |
| shape=selective/expression or vs union_all | 3.143710 | 5.245558 | 1.6686 | false | within-run medians differ; whole-run replication required before a verdict |
| shape=selective/expression or vs union | 3.143710 | 4.946131 | 1.5733 | false | within-run medians differ; whole-run replication required before a verdict |
| shape=selective/cast or vs union_all | 2.983885 | 5.063655 | 1.6970 | false | within-run medians differ; whole-run replication required before a verdict |
| shape=selective/cast or vs union | 2.983885 | 5.144890 | 1.7242 | false | within-run medians differ; whole-run replication required before a verdict |
| shape=broad/indexable or vs union_all | 5.001992 | 7.029154 | 1.4053 | false | within-run medians differ; whole-run replication required before a verdict |
| shape=broad/indexable or vs union | 5.001992 | 7.570294 | 1.5135 | false | within-run medians differ; whole-run replication required before a verdict |
| shape=broad/expression or vs union_all | 5.237082 | 7.314392 | 1.3967 | false | within-run medians differ; whole-run replication required before a verdict |
| shape=broad/expression or vs union | 5.237082 | 7.989638 | 1.5256 | false | within-run medians differ; whole-run replication required before a verdict |
| shape=broad/cast or vs union_all | 5.366463 | 11.271674 | 2.1004 | false | within-run medians differ; whole-run replication required before a verdict |
| shape=broad/cast or vs union | 5.366463 | 11.366105 | 2.1180 | false | within-run medians differ; whole-run replication required before a verdict |

## Pagination (P11)

| depth | spelling | status | median ms | spread % | trials | rows | plan structure |
|---|---|---|---|---|---|---|---|
| 1 | offset | valid | 6.020191 | 27.9 | 5 | 5 | `ac4509d295bf3348` |
| 1 | keyset | valid | 7.588768 | 19.8 | 5 | 5 | `166777c9921144af` |
| 10 | offset | valid | 6.419176 | 15.7 | 5 | 5 | `ac4509d295bf3348` |
| 10 | keyset | valid | 8.163653 | 12.1 | 5 | 5 | `d596ca0d42dc4da7` |
| 100 | offset | valid | 5.810455 | 40.5 | 5 | 5 | `ac4509d295bf3348` |
| 100 | keyset | valid | 6.519669 | 24.5 | 5 | 5 | `727a2962d7392ccd` |
| 1000 | offset | valid | 5.392865 | 51.4 | 5 | 5 | `ac4509d295bf3348` |
| 1000 | keyset | valid | 7.721890 | 12.8 | 5 | 5 | `441b4191cb138ed5` |
| 10000 | offset | valid | 8.699910 | 17.6 | 5 | 5 | `ac4509d295bf3348` |
| 10000 | keyset | valid | 4.685009 | 19.1 | 5 | 5 | `e1c72b3614c26105` |

Median ratio table (keyset / OFFSET; a ratio below 1 means keyset was faster in this run):

| depth | ratio (keyset/offset) | offset median ms | keyset median ms |
|---|---|---|---|
| 1 | 1.2606 | 6.020191 | 7.588768 |
| 10 | 1.2718 | 6.419176 | 8.163653 |
| 100 | 1.1221 | 5.810455 | 6.519669 |
| 1000 | 1.4319 | 5.392865 | 7.721890 |
| 10000 | 0.5385 | 8.699910 | 4.685009 |

### P11 pairs

| pair | A median ms | B median ms | B/A | same plan | within-run verdict |
|---|---|---|---|---|---|
| depth=1 offset vs keyset | 6.020191 | 7.588768 | 1.2606 | false | medians inside the two cells' own spread; no resolved difference within this run |
| depth=10 offset vs keyset | 6.419176 | 8.163653 | 1.2718 | false | medians inside the two cells' own spread; no resolved difference within this run |
| depth=100 offset vs keyset | 5.810455 | 6.519669 | 1.1221 | false | medians inside the two cells' own spread; no resolved difference within this run |
| depth=1000 offset vs keyset | 5.392865 | 7.721890 | 1.4319 | false | medians inside the two cells' own spread; no resolved difference within this run |
| depth=10000 offset vs keyset | 8.699910 | 4.685009 | 0.5385 | false | within-run medians differ; whole-run replication required before a verdict |

## Repeated-design position control

Not measured in this run.

## Expression-index arm: write and storage cost

Index `orders_amount_expr_idx` on expression `(amount + 0)`; supported=true.

| measure | value |
|---|---|
| index size (bytes) | 0 |
| insert batch with index (ms) | 15.532090 |
| insert batch without index (ms) | 21.413708 |
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

Rows: 5; serialized decoded-row JSON bytes: 51 (not wire bytes). Planning/execution: 0.082000 / 4.470000 ms.

```text
Limit  (cost=884.51..884.52 rows=5 width=8) (actual time=4.425..4.427 rows=5 loops=1)
  ->  Sort  (cost=884.51..889.29 rows=1912 width=8) (actual time=4.424..4.424 rows=5 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 25kB
        ->  Seq Scan on orders  (cost=21.08..852.75 rows=1912 width=8) (actual time=2.363..4.198 rows=2000 loops=1)
              Storage Filter: ((amount < 67) OR ((amount = 67) AND (id > 4310)) OR (amount IS NULL))
Planning Time: 0.082 ms
Execution Time: 4.470 ms
Peak Memory Usage: 92 kB

```

### s16_page_offset [depth=10000;spelling=offset]

Rows: 5; serialized decoded-row JSON bytes: 51 (not wire bytes). Planning/execution: 0.066000 / 8.632000 ms.

```text
Limit  (cost=1879.09..1879.10 rows=5 width=8) (actual time=8.574..8.576 rows=5 loops=1)
  ->  Sort  (cost=1854.09..1884.09 rows=12000 width=8) (actual time=7.449..7.795 rows=10005 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: quicksort  Memory: 947kB
        ->  Seq Scan on orders  (cost=21.08..1041.05 rows=12000 width=8) (actual time=1.007..4.108 rows=12000 loops=1)
Planning Time: 0.066 ms
Execution Time: 8.632 ms
Peak Memory Usage: 1516 kB

```

### s17_page_keyset [depth=1000;spelling=keyset]

Rows: 5; serialized decoded-row JSON bytes: 59 (not wire bytes). Planning/execution: 0.088000 / 8.000000 ms.

```text
Limit  (cost=1274.79..1274.81 rows=5 width=8) (actual time=7.957..7.959 rows=5 loops=1)
  ->  Sort  (cost=1274.79..1301.14 rows=10538 width=8) (actual time=7.956..7.956 rows=5 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 25kB
        ->  Seq Scan on orders  (cost=21.08..1099.76 rows=10538 width=8) (actual time=0.793..6.700 rows=11000 loops=1)
              Storage Filter: ((amount < 457) OR ((amount = 457) AND (id > 11254)) OR (amount IS NULL))
Planning Time: 0.088 ms
Execution Time: 8.000 ms
Peak Memory Usage: 92 kB

```

### s16_page_offset [depth=1000;spelling=offset]

Rows: 5; serialized decoded-row JSON bytes: 59 (not wire bytes). Planning/execution: 0.056000 / 6.152000 ms.

```text
Limit  (cost=1701.93..1701.94 rows=5 width=8) (actual time=6.112..6.114 rows=5 loops=1)
  ->  Sort  (cost=1699.43..1729.43 rows=12000 width=8) (actual time=6.000..6.029 rows=1005 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 130kB
        ->  Seq Scan on orders  (cost=21.08..1041.05 rows=12000 width=8) (actual time=0.605..3.799 rows=12000 loops=1)
Planning Time: 0.056 ms
Execution Time: 6.152 ms
Peak Memory Usage: 276 kB

```

### s17_page_keyset [depth=100;spelling=keyset]

Rows: 5; serialized decoded-row JSON bytes: 56 (not wire bytes). Planning/execution: 0.085000 / 6.724000 ms.

```text
Limit  (cost=1315.00..1315.01 rows=5 width=8) (actual time=6.681..6.682 rows=5 loops=1)
  ->  Sort  (cost=1315.00..1343.41 rows=11366 width=8) (actual time=6.679..6.680 rows=5 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 25kB
        ->  Seq Scan on orders  (cost=21.08..1126.21 rows=11366 width=8) (actual time=0.769..5.403 rows=11900 loops=1)
              Storage Filter: ((amount < 497) OR ((amount = 497) AND (id > 6723)) OR (amount IS NULL))
Planning Time: 0.085 ms
Execution Time: 6.724 ms
Peak Memory Usage: 92 kB

```

### s16_page_offset [depth=100;spelling=offset]

Rows: 5; serialized decoded-row JSON bytes: 56 (not wire bytes). Planning/execution: 0.058000 / 8.415000 ms.

```text
Limit  (cost=1504.15..1504.16 rows=5 width=8) (actual time=8.372..8.374 rows=5 loops=1)
  ->  Sort  (cost=1503.90..1533.90 rows=12000 width=8) (actual time=8.360..8.363 rows=105 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 33kB
        ->  Seq Scan on orders  (cost=21.08..1041.05 rows=12000 width=8) (actual time=0.856..6.826 rows=12000 loops=1)
Planning Time: 0.058 ms
Execution Time: 8.415 ms
Peak Memory Usage: 128 kB

```

### s17_page_keyset [depth=10;spelling=keyset]

Rows: 5; serialized decoded-row JSON bytes: 56 (not wire bytes). Planning/execution: 0.088000 / 7.255000 ms.

```text
Limit  (cost=1317.54..1317.55 rows=5 width=8) (actual time=7.215..7.217 rows=5 loops=1)
  ->  Sort  (cost=1317.54..1346.22 rows=11471 width=8) (actual time=7.213..7.214 rows=5 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 25kB
        ->  Seq Scan on orders  (cost=21.08..1127.01 rows=11471 width=8) (actual time=0.741..5.919 rows=11990 loops=1)
              Storage Filter: ((amount < 500) OR ((amount = 500) AND (id > 6319)) OR (amount IS NULL))
Planning Time: 0.088 ms
Execution Time: 7.255 ms
Peak Memory Usage: 92 kB

```

### s16_page_offset [depth=10;spelling=offset]

Rows: 5; serialized decoded-row JSON bytes: 56 (not wire bytes). Planning/execution: 0.073000 / 6.721000 ms.

```text
Limit  (cost=1335.49..1335.50 rows=5 width=8) (actual time=6.682..6.684 rows=5 loops=1)
  ->  Sort  (cost=1335.46..1365.46 rows=12000 width=8) (actual time=6.679..6.680 rows=15 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 26kB
        ->  Seq Scan on orders  (cost=21.08..1041.05 rows=12000 width=8) (actual time=0.889..5.388 rows=12000 loops=1)
Planning Time: 0.073 ms
Execution Time: 6.721 ms
Peak Memory Usage: 92 kB

```

### s17_page_keyset [depth=1;spelling=keyset]

Rows: 5; serialized decoded-row JSON bytes: 56 (not wire bytes). Planning/execution: 0.090000 / 6.760000 ms.

```text
Limit  (cost=1317.57..1317.58 rows=5 width=8) (actual time=6.715..6.717 rows=5 loops=1)
  ->  Sort  (cost=1317.57..1346.25 rows=11472 width=8) (actual time=6.714..6.715 rows=5 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 25kB
        ->  Seq Scan on orders  (cost=21.08..1127.02 rows=11472 width=8) (actual time=0.987..5.415 rows=11999 loops=1)
              Storage Filter: ((amount < 500) OR ((amount = 500) AND (id > 1554)) OR (amount IS NULL))
Planning Time: 0.090 ms
Execution Time: 6.760 ms
Peak Memory Usage: 92 kB

```

### s16_page_offset [depth=1;spelling=offset]

Rows: 5; serialized decoded-row JSON bytes: 56 (not wire bytes). Planning/execution: 0.080000 / 5.261000 ms.

```text
Limit  (cost=1256.15..1256.16 rows=5 width=8) (actual time=5.205..5.206 rows=5 loops=1)
  ->  Sort  (cost=1256.15..1286.15 rows=12000 width=8) (actual time=5.203..5.203 rows=6 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 25kB
        ->  Seq Scan on orders  (cost=21.08..1041.05 rows=12000 width=8) (actual time=0.905..3.897 rows=12000 loops=1)
Planning Time: 0.080 ms
Execution Time: 5.261 ms
Peak Memory Usage: 92 kB

```

### s18_or_predicate [shape=broad/cast;spelling=or]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.061000 / 4.735000 ms.

```text
Sort  (cost=1443.14..1459.81 rows=6667 width=8) (actual time=4.271..4.400 rows=4652 loops=1)
  Sort Key: id
  Sort Method: quicksort  Memory: 411kB
  ->  Seq Scan on orders  (cost=21.08..1019.70 rows=6667 width=8) (actual time=0.980..3.540 rows=4652 loops=1)
        Storage Filter: (((amount)::bigint <= 100) OR ((amount)::bigint >= 400))
Planning Time: 0.061 ms
Execution Time: 4.735 ms
Peak Memory Usage: 548 kB

```

### s20_union_predicate [shape=broad/cast;spelling=union]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.086000 / 11.159000 ms.

```text
Sort  (cost=2516.53..2536.53 rows=8000 width=8) (actual time=10.358..10.567 rows=4652 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 411kB
  ->  HashAggregate  (cost=1917.90..1997.90 rows=8000 width=8) (actual time=8.989..9.410 rows=4652 loops=1)
        Group Key: orders.id, orders.amount
        Batches: 1  Memory Usage: 657kB
        ->  Append  (cost=21.08..1877.90 rows=8000 width=8) (actual time=1.907..7.887 rows=4652 loops=1)
              ->  Seq Scan on orders  (cost=21.08..878.95 rows=4000 width=8) (actual time=1.905..3.794 rows=2273 loops=1)
                    Storage Filter: ((amount)::bigint <= 100)
              ->  Seq Scan on orders orders_1  (cost=21.08..878.95 rows=4000 width=8) (actual time=1.669..3.500 rows=2379 loops=1)
                    Storage Filter: ((amount)::bigint >= 400)
Planning Time: 0.086 ms
Execution Time: 11.159 ms
Peak Memory Usage: 1252 kB

```

### s19_union_all_predicate [shape=broad/cast;spelling=union_all]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.108000 / 9.704000 ms.

```text
Sort  (cost=2396.53..2416.53 rows=8000 width=8) (actual time=9.091..9.296 rows=4652 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 411kB
  ->  Append  (cost=21.08..1877.90 rows=8000 width=8) (actual time=1.967..8.231 rows=4652 loops=1)
        ->  Seq Scan on orders  (cost=21.08..878.95 rows=4000 width=8) (actual time=1.965..4.208 rows=2273 loops=1)
              Storage Filter: ((amount)::bigint <= 100)
        ->  Seq Scan on orders orders_1  (cost=21.08..878.95 rows=4000 width=8) (actual time=1.676..3.500 rows=2379 loops=1)
              Storage Filter: ((amount)::bigint >= 400)
Planning Time: 0.108 ms
Execution Time: 9.704 ms
Peak Memory Usage: 625 kB

```

### s18_or_predicate [shape=broad/expression;index=expression_index;spelling=or]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.113000 / 9.679000 ms.

```text
Sort  (cost=1214.19..1224.72 rows=4211 width=8) (actual time=9.064..9.253 rows=4652 loops=1)
  Sort Key: id
  Sort Method: quicksort  Memory: 411kB
  ->  Seq Scan on orders  (cost=21.08..960.69 rows=4211 width=8) (actual time=2.852..7.825 rows=4652 loops=1)
        Storage Filter: (((amount + 0) <= 100) OR ((amount + 0) >= 400))
Planning Time: 0.113 ms
Execution Time: 9.679 ms
Peak Memory Usage: 548 kB

```

### s20_union_predicate [shape=broad/expression;index=expression_index;spelling=union]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.088000 / 7.871000 ms.

```text
Sort  (cost=2116.31..2127.97 rows=4664 width=8) (actual time=7.405..7.520 rows=4652 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 411kB
  ->  HashAggregate  (cost=1785.46..1832.10 rows=4664 width=8) (actual time=6.524..6.744 rows=4652 loops=1)
        Group Key: orders.id, orders.amount
        Batches: 1  Memory Usage: 465kB
        ->  Append  (cost=21.08..1762.14 rows=4664 width=8) (actual time=1.296..5.738 rows=4652 loops=1)
              ->  Seq Scan on orders  (cost=21.08..845.59 rows=2267 width=8) (actual time=1.292..2.666 rows=2273 loops=1)
                    Storage Filter: ((amount + 0) <= 100)
              ->  Seq Scan on orders orders_1  (cost=21.08..846.59 rows=2397 width=8) (actual time=1.299..2.570 rows=2379 loops=1)
                    Storage Filter: ((amount + 0) >= 400)
Planning Time: 0.088 ms
Execution Time: 7.871 ms
Peak Memory Usage: 1060 kB

```

### s19_union_all_predicate [shape=broad/expression;index=expression_index;spelling=union_all]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.124000 / 9.429000 ms.

```text
Sort  (cost=2046.35..2058.01 rows=4664 width=8) (actual time=8.748..8.965 rows=4652 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 411kB
  ->  Append  (cost=21.08..1762.14 rows=4664 width=8) (actual time=1.762..7.839 rows=4652 loops=1)
        ->  Seq Scan on orders  (cost=21.08..845.59 rows=2267 width=8) (actual time=1.760..3.264 rows=2273 loops=1)
              Storage Filter: ((amount + 0) <= 100)
        ->  Seq Scan on orders orders_1  (cost=21.08..846.59 rows=2397 width=8) (actual time=1.430..4.027 rows=2379 loops=1)
              Storage Filter: ((amount + 0) >= 400)
Planning Time: 0.124 ms
Execution Time: 9.429 ms
Peak Memory Usage: 625 kB

```

### s18_or_predicate [shape=broad/expression;spelling=or]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.057000 / 4.358000 ms.

```text
Sort  (cost=1443.14..1459.81 rows=6667 width=8) (actual time=3.845..3.980 rows=4652 loops=1)
  Sort Key: id
  Sort Method: quicksort  Memory: 411kB
  ->  Seq Scan on orders  (cost=21.08..1019.70 rows=6667 width=8) (actual time=0.869..3.093 rows=4652 loops=1)
        Storage Filter: (((amount + 0) <= 100) OR ((amount + 0) >= 400))
Planning Time: 0.057 ms
Execution Time: 4.358 ms
Peak Memory Usage: 548 kB

```

### s20_union_predicate [shape=broad/expression;spelling=union]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.074000 / 7.958000 ms.

```text
Sort  (cost=2516.53..2536.53 rows=8000 width=8) (actual time=7.469..7.591 rows=4652 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 411kB
  ->  HashAggregate  (cost=1917.90..1997.90 rows=8000 width=8) (actual time=6.449..6.739 rows=4652 loops=1)
        Group Key: orders.id, orders.amount
        Batches: 1  Memory Usage: 657kB
        ->  Append  (cost=21.08..1877.90 rows=8000 width=8) (actual time=1.406..5.700 rows=4652 loops=1)
              ->  Seq Scan on orders  (cost=21.08..878.95 rows=4000 width=8) (actual time=1.405..2.860 rows=2273 loops=1)
                    Storage Filter: ((amount + 0) <= 100)
              ->  Seq Scan on orders orders_1  (cost=21.08..878.95 rows=4000 width=8) (actual time=1.296..2.358 rows=2379 loops=1)
                    Storage Filter: ((amount + 0) >= 400)
Planning Time: 0.074 ms
Execution Time: 7.958 ms
Peak Memory Usage: 1252 kB

```

### s19_union_all_predicate [shape=broad/expression;spelling=union_all]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.124000 / 6.867000 ms.

```text
Sort  (cost=2396.53..2416.53 rows=8000 width=8) (actual time=6.364..6.520 rows=4652 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 411kB
  ->  Append  (cost=21.08..1877.90 rows=8000 width=8) (actual time=1.288..5.607 rows=4652 loops=1)
        ->  Seq Scan on orders  (cost=21.08..878.95 rows=4000 width=8) (actual time=1.287..2.546 rows=2273 loops=1)
              Storage Filter: ((amount + 0) <= 100)
        ->  Seq Scan on orders orders_1  (cost=21.08..878.95 rows=4000 width=8) (actual time=1.262..2.595 rows=2379 loops=1)
              Storage Filter: ((amount + 0) >= 400)
Planning Time: 0.124 ms
Execution Time: 6.867 ms
Peak Memory Usage: 625 kB

```

### s18_or_predicate [shape=broad/indexable;spelling=or]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.064000 / 4.178000 ms.

```text
Sort  (cost=1154.19..1164.72 rows=4211 width=8) (actual time=3.678..3.814 rows=4652 loops=1)
  Sort Key: id
  Sort Method: quicksort  Memory: 411kB
  ->  Seq Scan on orders  (cost=21.08..900.69 rows=4211 width=8) (actual time=0.995..2.892 rows=4652 loops=1)
        Storage Filter: ((amount <= 100) OR (amount >= 400))
Planning Time: 0.064 ms
Execution Time: 4.178 ms
Peak Memory Usage: 548 kB

```

### s20_union_predicate [shape=broad/indexable;spelling=union]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.075000 / 7.249000 ms.

```text
Sort  (cost=2056.31..2067.97 rows=4664 width=8) (actual time=6.751..6.873 rows=4652 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 411kB
  ->  HashAggregate  (cost=1725.46..1772.10 rows=4664 width=8) (actual time=5.798..6.043 rows=4652 loops=1)
        Group Key: orders.id, orders.amount
        Batches: 1  Memory Usage: 465kB
        ->  Append  (cost=21.08..1702.14 rows=4664 width=8) (actual time=1.277..5.065 rows=4652 loops=1)
              ->  Seq Scan on orders  (cost=21.08..815.59 rows=2267 width=8) (actual time=1.276..2.528 rows=2273 loops=1)
                    Storage Filter: (amount <= 100)
              ->  Seq Scan on orders orders_1  (cost=21.08..816.59 rows=2397 width=8) (actual time=1.086..2.066 rows=2379 loops=1)
                    Storage Filter: (amount >= 400)
Planning Time: 0.075 ms
Execution Time: 7.249 ms
Peak Memory Usage: 1060 kB

```

### s19_union_all_predicate [shape=broad/indexable;spelling=union_all]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.089000 / 6.457000 ms.

```text
Sort  (cost=1986.35..1998.01 rows=4664 width=8) (actual time=5.973..6.115 rows=4652 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 411kB
  ->  Append  (cost=21.08..1702.14 rows=4664 width=8) (actual time=1.348..5.167 rows=4652 loops=1)
        ->  Seq Scan on orders  (cost=21.08..815.59 rows=2267 width=8) (actual time=1.347..2.487 rows=2273 loops=1)
              Storage Filter: (amount <= 100)
        ->  Seq Scan on orders orders_1  (cost=21.08..816.59 rows=2397 width=8) (actual time=1.094..2.208 rows=2379 loops=1)
              Storage Filter: (amount >= 400)
Planning Time: 0.089 ms
Execution Time: 6.457 ms
Peak Memory Usage: 625 kB

```

### s18_or_predicate [shape=selective/cast;spelling=or]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.034000 / 2.990000 ms.

```text
Sort  (cost=1443.14..1459.81 rows=6667 width=8) (actual time=2.937..2.944 rows=251 loops=1)
  Sort Key: id
  Sort Method: quicksort  Memory: 36kB
  ->  Seq Scan on orders  (cost=21.08..1019.70 rows=6667 width=8) (actual time=2.872..2.889 rows=251 loops=1)
        Storage Filter: (((amount)::bigint <= 5) OR ((amount)::bigint >= 496))
Planning Time: 0.034 ms
Execution Time: 2.990 ms
Peak Memory Usage: 125 kB

```

### s20_union_predicate [shape=selective/cast;spelling=union]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.056000 / 4.956000 ms.

```text
Sort  (cost=2516.53..2536.53 rows=8000 width=8) (actual time=4.882..4.889 rows=251 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 36kB
  ->  HashAggregate  (cost=1917.90..1997.90 rows=8000 width=8) (actual time=4.792..4.833 rows=251 loops=1)
        Group Key: orders.id, orders.amount
        Batches: 1  Memory Usage: 417kB
        ->  Append  (cost=21.08..1877.90 rows=8000 width=8) (actual time=2.372..4.715 rows=251 loops=1)
              ->  Seq Scan on orders  (cost=21.08..878.95 rows=4000 width=8) (actual time=2.371..2.379 rows=114 loops=1)
                    Storage Filter: ((amount)::bigint <= 5)
              ->  Seq Scan on orders orders_1  (cost=21.08..878.95 rows=4000 width=8) (actual time=2.299..2.309 rows=137 loops=1)
                    Storage Filter: ((amount)::bigint >= 496)
Planning Time: 0.056 ms
Execution Time: 4.956 ms
Peak Memory Usage: 600 kB

```

### s19_union_all_predicate [shape=selective/cast;spelling=union_all]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.072000 / 4.634000 ms.

```text
Sort  (cost=2396.53..2416.53 rows=8000 width=8) (actual time=4.552..4.560 rows=251 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 36kB
  ->  Append  (cost=21.08..1877.90 rows=8000 width=8) (actual time=2.373..4.506 rows=251 loops=1)
        ->  Seq Scan on orders  (cost=21.08..878.95 rows=4000 width=8) (actual time=2.372..2.380 rows=114 loops=1)
              Storage Filter: ((amount)::bigint <= 5)
        ->  Seq Scan on orders orders_1  (cost=21.08..878.95 rows=4000 width=8) (actual time=2.091..2.100 rows=137 loops=1)
              Storage Filter: ((amount)::bigint >= 496)
Planning Time: 0.072 ms
Execution Time: 4.634 ms
Peak Memory Usage: 202 kB

```

### s18_or_predicate [shape=selective/expression;index=expression_index;spelling=or]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.048000 / 3.242000 ms.

```text
Sort  (cost=860.56..861.21 rows=262 width=8) (actual time=3.182..3.190 rows=251 loops=1)
  Sort Key: id
  Sort Method: quicksort  Memory: 36kB
  ->  Seq Scan on orders  (cost=21.08..850.03 rows=262 width=8) (actual time=3.090..3.108 rows=251 loops=1)
        Storage Filter: (((amount + 0) <= 5) OR ((amount + 0) >= 496))
Planning Time: 0.048 ms
Execution Time: 3.242 ms
Peak Memory Usage: 125 kB

```

### s20_union_predicate [shape=selective/expression;index=expression_index;spelling=union]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.143000 / 6.262000 ms.

```text
Sort  (cost=1596.53..1597.19 rows=263 width=8) (actual time=6.176..6.184 rows=251 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 36kB
  ->  HashAggregate  (cost=1583.33..1585.96 rows=263 width=8) (actual time=6.114..6.130 rows=251 loops=1)
        Group Key: orders.id, orders.amount
        Batches: 1  Memory Usage: 45kB
        ->  Append  (cost=21.08..1582.02 rows=263 width=8) (actual time=3.188..6.055 rows=251 loops=1)
              ->  Seq Scan on orders  (cost=21.08..788.85 rows=107 width=8) (actual time=3.186..3.197 rows=114 loops=1)
                    Storage Filter: ((amount + 0) <= 5)
              ->  Seq Scan on orders orders_1  (cost=21.08..789.22 rows=156 width=8) (actual time=2.814..2.825 rows=137 loops=1)
                    Storage Filter: ((amount + 0) >= 496)
Planning Time: 0.143 ms
Execution Time: 6.262 ms
Peak Memory Usage: 224 kB

```

### s19_union_all_predicate [shape=selective/expression;index=expression_index;spelling=union_all]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.097000 / 5.409000 ms.

```text
Sort  (cost=1592.59..1593.25 rows=263 width=8) (actual time=5.338..5.346 rows=251 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 36kB
  ->  Append  (cost=21.08..1582.02 rows=263 width=8) (actual time=2.677..5.277 rows=251 loops=1)
        ->  Seq Scan on orders  (cost=21.08..788.85 rows=107 width=8) (actual time=2.675..2.685 rows=114 loops=1)
              Storage Filter: ((amount + 0) <= 5)
        ->  Seq Scan on orders orders_1  (cost=21.08..789.22 rows=156 width=8) (actual time=2.550..2.561 rows=137 loops=1)
              Storage Filter: ((amount + 0) >= 496)
Planning Time: 0.097 ms
Execution Time: 5.409 ms
Peak Memory Usage: 202 kB

```

### s18_or_predicate [shape=selective/expression;spelling=or]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.040000 / 3.180000 ms.

```text
Sort  (cost=1443.14..1459.81 rows=6667 width=8) (actual time=3.123..3.131 rows=251 loops=1)
  Sort Key: id
  Sort Method: quicksort  Memory: 36kB
  ->  Seq Scan on orders  (cost=21.08..1019.70 rows=6667 width=8) (actual time=3.052..3.071 rows=251 loops=1)
        Storage Filter: (((amount + 0) <= 5) OR ((amount + 0) >= 496))
Planning Time: 0.040 ms
Execution Time: 3.180 ms
Peak Memory Usage: 125 kB

```

### s20_union_predicate [shape=selective/expression;spelling=union]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.068000 / 4.619000 ms.

```text
Sort  (cost=2516.53..2536.53 rows=8000 width=8) (actual time=4.488..4.494 rows=251 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 36kB
  ->  HashAggregate  (cost=1917.90..1997.90 rows=8000 width=8) (actual time=4.393..4.425 rows=251 loops=1)
        Group Key: orders.id, orders.amount
        Batches: 1  Memory Usage: 417kB
        ->  Append  (cost=21.08..1877.90 rows=8000 width=8) (actual time=2.177..4.353 rows=251 loops=1)
              ->  Seq Scan on orders  (cost=21.08..878.95 rows=4000 width=8) (actual time=2.176..2.184 rows=114 loops=1)
                    Storage Filter: ((amount + 0) <= 5)
              ->  Seq Scan on orders orders_1  (cost=21.08..878.95 rows=4000 width=8) (actual time=2.136..2.144 rows=137 loops=1)
                    Storage Filter: ((amount + 0) >= 496)
Planning Time: 0.068 ms
Execution Time: 4.619 ms
Peak Memory Usage: 600 kB

```

### s19_union_all_predicate [shape=selective/expression;spelling=union_all]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.094000 / 4.857000 ms.

```text
Sort  (cost=2396.53..2416.53 rows=8000 width=8) (actual time=4.743..4.751 rows=251 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 36kB
  ->  Append  (cost=21.08..1877.90 rows=8000 width=8) (actual time=2.320..4.685 rows=251 loops=1)
        ->  Seq Scan on orders  (cost=21.08..878.95 rows=4000 width=8) (actual time=2.319..2.328 rows=114 loops=1)
              Storage Filter: ((amount + 0) <= 5)
        ->  Seq Scan on orders orders_1  (cost=21.08..878.95 rows=4000 width=8) (actual time=2.301..2.329 rows=137 loops=1)
              Storage Filter: ((amount + 0) >= 496)
Planning Time: 0.094 ms
Execution Time: 4.857 ms
Peak Memory Usage: 202 kB

```

### s18_or_predicate [shape=selective/indexable;spelling=or]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.065000 / 2.563000 ms.

```text
Sort  (cost=800.56..801.21 rows=262 width=8) (actual time=2.512..2.519 rows=251 loops=1)
  Sort Key: id
  Sort Method: quicksort  Memory: 36kB
  ->  Seq Scan on orders  (cost=21.08..790.03 rows=262 width=8) (actual time=2.414..2.431 rows=251 loops=1)
        Storage Filter: ((amount <= 5) OR (amount >= 496))
Planning Time: 0.065 ms
Execution Time: 2.563 ms
Peak Memory Usage: 125 kB

```

### s20_union_predicate [shape=selective/indexable;spelling=union]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.070000 / 4.662000 ms.

```text
Sort  (cost=1536.53..1537.19 rows=263 width=8) (actual time=4.600..4.607 rows=251 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 36kB
  ->  HashAggregate  (cost=1523.33..1525.96 rows=263 width=8) (actual time=4.542..4.559 rows=251 loops=1)
        Group Key: orders.id, orders.amount
        Batches: 1  Memory Usage: 45kB
        ->  Append  (cost=21.08..1522.02 rows=263 width=8) (actual time=2.192..4.307 rows=251 loops=1)
              ->  Seq Scan on orders  (cost=21.08..758.85 rows=107 width=8) (actual time=2.188..2.248 rows=114 loops=1)
                    Storage Filter: (amount <= 5)
              ->  Seq Scan on orders orders_1  (cost=21.08..759.22 rows=156 width=8) (actual time=1.944..1.954 rows=137 loops=1)
                    Storage Filter: (amount >= 496)
Planning Time: 0.070 ms
Execution Time: 4.662 ms
Peak Memory Usage: 224 kB

```

### s19_union_all_predicate [shape=selective/indexable;spelling=union_all]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.079000 / 4.470000 ms.

```text
Sort  (cost=1532.59..1533.25 rows=263 width=8) (actual time=4.386..4.393 rows=251 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 36kB
  ->  Append  (cost=21.08..1522.02 rows=263 width=8) (actual time=2.185..4.321 rows=251 loops=1)
        ->  Seq Scan on orders  (cost=21.08..758.85 rows=107 width=8) (actual time=2.182..2.191 rows=114 loops=1)
              Storage Filter: (amount <= 5)
        ->  Seq Scan on orders orders_1  (cost=21.08..759.22 rows=156 width=8) (actual time=2.092..2.102 rows=137 loops=1)
              Storage Filter: (amount >= 496)
Planning Time: 0.079 ms
Execution Time: 4.470 ms
Peak Memory Usage: 202 kB

```

## Limits and gaps recorded by the runner

- no concurrent writers: client-decode JSON bytes are not wire bytes; tails are under-sampled at the protocol trial count
- one host: no independent hardware, real-network, 3-node placement or equal-total-budget comparison is established
- single-node only; the YB 3-node RF=3 per-node/equal-total arms are a coverage gap on this host

