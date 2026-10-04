# SELECT access-path matrix (generated measurements)

_Measurements only. Verdicts for P10/P11 belong in a signed analysis._

**TL;DR:** within-run medians for OR / UNION ALL / UNION at three predicate shapes and OFFSET / keyset at five page depths, each cell verified against the independent oracle. A within-run difference is not a replicated result. Missing measurements are N/A.

| run | engine | topology | scale | page limit |
|---|---|---|---|---|
| select-access-yb-2-20261003 | yugabyte | single | deep | 5 |

## Provenance

| Field | Value |
|---|---|
| repo commit | `6ad01ef65352b0be194007dab460b297ae94a541` |
| benchmark image | `localhost/charitybench-select:select-access-yb-2-20261003` |
| image digest | `023c9fb2d709de026afbef825315bf2e29a234361a3d2a0f95507136c4fa57f5` |
| run tag | `run/study-01/select-access-yb-2-20261003` |
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
| broad/cast | default | or | valid | 6.986492 | 12.5 | 5 | 4652 | `1dd363b32be056fc` |
| broad/cast | default | union | valid | 11.844361 | 24.5 | 5 | 4652 | `b5f6c293b9267d60` |
| broad/cast | default | union_all | valid | 10.534562 | 9.7 | 5 | 4652 | `24cabca42adfa691` |
| broad/expression | expression_index | or | valid | 7.119667 | 20.3 | 5 | 4652 | `0e731990206e1446` |
| broad/expression | expression_index | union | valid | 11.705941 | 21.6 | 5 | 4652 | `bf74591baa836f9d` |
| broad/expression | expression_index | union_all | valid | 10.532490 | 9.6 | 5 | 4652 | `a7e687668ac364de` |
| broad/expression | default | or | valid | 7.530308 | 4.1 | 5 | 4652 | `0e731990206e1446` |
| broad/expression | default | union | valid | 11.142094 | 13.5 | 5 | 4652 | `bf74591baa836f9d` |
| broad/expression | default | union_all | valid | 10.005592 | 22.2 | 5 | 4652 | `a7e687668ac364de` |
| broad/indexable | default | or | valid | 8.036357 | 24.5 | 5 | 4652 | `5191eeb70daef797` |
| broad/indexable | default | union | valid | 9.568491 | 21.7 | 5 | 4652 | `12de8f648f9520a0` |
| broad/indexable | default | union_all | valid | 10.539168 | 18.2 | 5 | 4652 | `d95e057d512f5ae5` |
| selective/cast | default | or | valid | 4.038333 | 48.1 | 5 | 251 | `57f71a4ae173cc97` |
| selective/cast | default | union | valid | 6.878930 | 18.9 | 5 | 251 | `2a31413606cd52a3` |
| selective/cast | default | union_all | valid | 7.142692 | 17.9 | 5 | 251 | `096e4d4e43a2a56e` |
| selective/expression | expression_index | or | valid | 4.153415 | 10.8 | 5 | 251 | `66a4eb5992604842` |
| selective/expression | expression_index | union | valid | 7.431587 | 13.8 | 5 | 251 | `14832cc886f611bd` |
| selective/expression | expression_index | union_all | valid | 7.027629 | 7.1 | 5 | 251 | `fee9cce2de640d1c` |
| selective/expression | default | or | valid | 4.051543 | 7.7 | 5 | 251 | `66a4eb5992604842` |
| selective/expression | default | union | valid | 7.521499 | 19.2 | 5 | 251 | `14832cc886f611bd` |
| selective/expression | default | union_all | valid | 7.011259 | 15.0 | 5 | 251 | `fee9cce2de640d1c` |
| selective/indexable | default | or | valid | 3.547935 | 15.4 | 5 | 251 | `66fcedd6207432cd` |
| selective/indexable | default | union | valid | 6.531859 | 14.3 | 5 | 251 | `40578cd947e4b93c` |
| selective/indexable | default | union_all | valid | 6.273535 | 17.9 | 5 | 251 | `8136cf54e3f9a863` |

### P10 pairs

Pairs compare the **default** arm only; the labelled `expression_index` arm is a package arm and is never pooled with the query-text pair.

| pair | A median ms | B median ms | B/A | same plan | within-run verdict |
|---|---|---|---|---|---|
| shape=selective/indexable or vs union_all | 3.547935 | 6.273535 | 1.7682 | false | within-run medians differ; whole-run replication required before a verdict |
| shape=selective/indexable or vs union | 3.547935 | 6.531859 | 1.8410 | false | within-run medians differ; whole-run replication required before a verdict |
| shape=selective/expression or vs union_all | 4.051543 | 7.011259 | 1.7305 | false | within-run medians differ; whole-run replication required before a verdict |
| shape=selective/expression or vs union | 4.051543 | 7.521499 | 1.8565 | false | within-run medians differ; whole-run replication required before a verdict |
| shape=selective/cast or vs union_all | 4.038333 | 7.142692 | 1.7687 | false | within-run medians differ; whole-run replication required before a verdict |
| shape=selective/cast or vs union | 4.038333 | 6.878930 | 1.7034 | false | within-run medians differ; whole-run replication required before a verdict |
| shape=broad/indexable or vs union_all | 8.036357 | 10.539168 | 1.3114 | false | medians inside the two cells' own spread; no resolved difference within this run |
| shape=broad/indexable or vs union | 8.036357 | 9.568491 | 1.1907 | false | medians inside the two cells' own spread; no resolved difference within this run |
| shape=broad/expression or vs union_all | 7.530308 | 10.005592 | 1.3287 | false | within-run medians differ; whole-run replication required before a verdict |
| shape=broad/expression or vs union | 7.530308 | 11.142094 | 1.4796 | false | within-run medians differ; whole-run replication required before a verdict |
| shape=broad/cast or vs union_all | 6.986492 | 10.534562 | 1.5078 | false | within-run medians differ; whole-run replication required before a verdict |
| shape=broad/cast or vs union | 6.986492 | 11.844361 | 1.6953 | false | within-run medians differ; whole-run replication required before a verdict |

## Pagination (P11)

| depth | spelling | status | median ms | spread % | trials | rows | plan structure |
|---|---|---|---|---|---|---|---|
| 1 | offset | valid | 6.739808 | 16.5 | 5 | 5 | `ac4509d295bf3348` |
| 1 | keyset | valid | 7.936387 | 17.7 | 5 | 5 | `166777c9921144af` |
| 10 | offset | valid | 6.560888 | 23.3 | 5 | 5 | `ac4509d295bf3348` |
| 10 | keyset | valid | 8.254318 | 17.5 | 5 | 5 | `d596ca0d42dc4da7` |
| 100 | offset | valid | 6.360430 | 30.4 | 5 | 5 | `ac4509d295bf3348` |
| 100 | keyset | valid | 7.006257 | 8.0 | 5 | 5 | `727a2962d7392ccd` |
| 1000 | offset | valid | 5.836768 | 37.2 | 5 | 5 | `ac4509d295bf3348` |
| 1000 | keyset | valid | 6.822253 | 7.1 | 5 | 5 | `441b4191cb138ed5` |
| 10000 | offset | valid | 7.418655 | 10.6 | 5 | 5 | `ac4509d295bf3348` |
| 10000 | keyset | valid | 4.998019 | 14.7 | 5 | 5 | `e1c72b3614c26105` |

Median ratio table (keyset / OFFSET; a ratio below 1 means keyset was faster in this run):

| depth | ratio (keyset/offset) | offset median ms | keyset median ms |
|---|---|---|---|
| 1 | 1.1775 | 6.739808 | 7.936387 |
| 10 | 1.2581 | 6.560888 | 8.254318 |
| 100 | 1.1015 | 6.360430 | 7.006257 |
| 1000 | 1.1688 | 5.836768 | 6.822253 |
| 10000 | 0.6737 | 7.418655 | 4.998019 |

### P11 pairs

| pair | A median ms | B median ms | B/A | same plan | within-run verdict |
|---|---|---|---|---|---|
| depth=1 offset vs keyset | 6.739808 | 7.936387 | 1.1775 | false | medians inside the two cells' own spread; no resolved difference within this run |
| depth=10 offset vs keyset | 6.560888 | 8.254318 | 1.2581 | false | medians inside the two cells' own spread; no resolved difference within this run |
| depth=100 offset vs keyset | 6.360430 | 7.006257 | 1.1015 | false | medians inside the two cells' own spread; no resolved difference within this run |
| depth=1000 offset vs keyset | 5.836768 | 6.822253 | 1.1688 | false | medians inside the two cells' own spread; no resolved difference within this run |
| depth=10000 offset vs keyset | 7.418655 | 4.998019 | 0.6737 | false | within-run medians differ; whole-run replication required before a verdict |

## Repeated-design position control

The first matrix cell and a repeat of it at the end of the phase (`position=last`); a large ratio is a session position effect, not a spelling effect.

| pair | first median ms | last median ms | ratio (last/first) | same plan |
|---|---|---|---|---|
| position control selective/indexable or first vs last | 3.547935 | 3.944609 | 1.1118 | true |

## Expression-index arm: write and storage cost

Index `orders_amount_expr_idx` on expression `(amount + 0)`; supported=true.

| measure | value |
|---|---|
| index size (bytes) | 0 |
| insert batch with index (ms) | 22.869442 |
| insert batch without index (ms) | 31.237521 |
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

Rows: 5; serialized decoded-row JSON bytes: 51 (not wire bytes). Planning/execution: 0.116000 / 5.051000 ms.

```text
Limit  (cost=884.51..884.52 rows=5 width=8) (actual time=5.003..5.005 rows=5 loops=1)
  ->  Sort  (cost=884.51..889.29 rows=1912 width=8) (actual time=5.001..5.002 rows=5 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 25kB
        ->  Seq Scan on orders  (cost=21.08..852.75 rows=1912 width=8) (actual time=2.201..4.738 rows=2000 loops=1)
              Storage Filter: ((amount < 67) OR ((amount = 67) AND (id > 4310)) OR (amount IS NULL))
Planning Time: 0.116 ms
Execution Time: 5.051 ms
Peak Memory Usage: 92 kB

```

### s16_page_offset [depth=10000;spelling=offset]

Rows: 5; serialized decoded-row JSON bytes: 51 (not wire bytes). Planning/execution: 0.067000 / 8.485000 ms.

```text
Limit  (cost=1879.09..1879.10 rows=5 width=8) (actual time=8.417..8.419 rows=5 loops=1)
  ->  Sort  (cost=1854.09..1884.09 rows=12000 width=8) (actual time=7.165..7.491 rows=10005 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: quicksort  Memory: 947kB
        ->  Seq Scan on orders  (cost=21.08..1041.05 rows=12000 width=8) (actual time=0.744..3.578 rows=12000 loops=1)
Planning Time: 0.067 ms
Execution Time: 8.485 ms
Peak Memory Usage: 1516 kB

```

### s17_page_keyset [depth=1000;spelling=keyset]

Rows: 5; serialized decoded-row JSON bytes: 59 (not wire bytes). Planning/execution: 0.101000 / 6.666000 ms.

```text
Limit  (cost=1274.79..1274.81 rows=5 width=8) (actual time=6.622..6.624 rows=5 loops=1)
  ->  Sort  (cost=1274.79..1301.14 rows=10538 width=8) (actual time=6.620..6.621 rows=5 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 25kB
        ->  Seq Scan on orders  (cost=21.08..1099.76 rows=10538 width=8) (actual time=1.174..5.288 rows=11000 loops=1)
              Storage Filter: ((amount < 457) OR ((amount = 457) AND (id > 11254)) OR (amount IS NULL))
Planning Time: 0.101 ms
Execution Time: 6.666 ms
Peak Memory Usage: 92 kB

```

### s16_page_offset [depth=1000;spelling=offset]

Rows: 5; serialized decoded-row JSON bytes: 59 (not wire bytes). Planning/execution: 0.060000 / 5.653000 ms.

```text
Limit  (cost=1701.93..1701.94 rows=5 width=8) (actual time=5.601..5.602 rows=5 loops=1)
  ->  Sort  (cost=1699.43..1729.43 rows=12000 width=8) (actual time=5.481..5.512 rows=1005 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 130kB
        ->  Seq Scan on orders  (cost=21.08..1041.05 rows=12000 width=8) (actual time=0.527..3.467 rows=12000 loops=1)
Planning Time: 0.060 ms
Execution Time: 5.653 ms
Peak Memory Usage: 276 kB

```

### s17_page_keyset [depth=100;spelling=keyset]

Rows: 5; serialized decoded-row JSON bytes: 56 (not wire bytes). Planning/execution: 0.099000 / 7.214000 ms.

```text
Limit  (cost=1315.00..1315.01 rows=5 width=8) (actual time=7.166..7.168 rows=5 loops=1)
  ->  Sort  (cost=1315.00..1343.41 rows=11366 width=8) (actual time=7.165..7.165 rows=5 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 25kB
        ->  Seq Scan on orders  (cost=21.08..1126.21 rows=11366 width=8) (actual time=0.849..5.736 rows=11900 loops=1)
              Storage Filter: ((amount < 497) OR ((amount = 497) AND (id > 6723)) OR (amount IS NULL))
Planning Time: 0.099 ms
Execution Time: 7.214 ms
Peak Memory Usage: 92 kB

```

### s16_page_offset [depth=100;spelling=offset]

Rows: 5; serialized decoded-row JSON bytes: 56 (not wire bytes). Planning/execution: 0.068000 / 5.636000 ms.

```text
Limit  (cost=1504.15..1504.16 rows=5 width=8) (actual time=5.593..5.595 rows=5 loops=1)
  ->  Sort  (cost=1503.90..1533.90 rows=12000 width=8) (actual time=5.579..5.583 rows=105 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 33kB
        ->  Seq Scan on orders  (cost=21.08..1041.05 rows=12000 width=8) (actual time=0.960..4.069 rows=12000 loops=1)
Planning Time: 0.068 ms
Execution Time: 5.636 ms
Peak Memory Usage: 128 kB

```

### s17_page_keyset [depth=10;spelling=keyset]

Rows: 5; serialized decoded-row JSON bytes: 56 (not wire bytes). Planning/execution: 0.095000 / 8.119000 ms.

```text
Limit  (cost=1317.54..1317.55 rows=5 width=8) (actual time=8.072..8.074 rows=5 loops=1)
  ->  Sort  (cost=1317.54..1346.22 rows=11471 width=8) (actual time=8.070..8.071 rows=5 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 25kB
        ->  Seq Scan on orders  (cost=21.08..1127.01 rows=11471 width=8) (actual time=0.886..6.363 rows=11990 loops=1)
              Storage Filter: ((amount < 500) OR ((amount = 500) AND (id > 6319)) OR (amount IS NULL))
Planning Time: 0.095 ms
Execution Time: 8.119 ms
Peak Memory Usage: 92 kB

```

### s16_page_offset [depth=10;spelling=offset]

Rows: 5; serialized decoded-row JSON bytes: 56 (not wire bytes). Planning/execution: 0.062000 / 5.769000 ms.

```text
Limit  (cost=1335.49..1335.50 rows=5 width=8) (actual time=5.726..5.728 rows=5 loops=1)
  ->  Sort  (cost=1335.46..1365.46 rows=12000 width=8) (actual time=5.723..5.724 rows=15 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 26kB
        ->  Seq Scan on orders  (cost=21.08..1041.05 rows=12000 width=8) (actual time=0.791..4.181 rows=12000 loops=1)
Planning Time: 0.062 ms
Execution Time: 5.769 ms
Peak Memory Usage: 92 kB

```

### s17_page_keyset [depth=1;spelling=keyset]

Rows: 5; serialized decoded-row JSON bytes: 56 (not wire bytes). Planning/execution: 0.159000 / 6.908000 ms.

```text
Limit  (cost=1317.57..1317.58 rows=5 width=8) (actual time=6.865..6.866 rows=5 loops=1)
  ->  Sort  (cost=1317.57..1346.25 rows=11472 width=8) (actual time=6.863..6.864 rows=5 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 25kB
        ->  Seq Scan on orders  (cost=21.08..1127.02 rows=11472 width=8) (actual time=0.753..5.408 rows=11999 loops=1)
              Storage Filter: ((amount < 500) OR ((amount = 500) AND (id > 1554)) OR (amount IS NULL))
Planning Time: 0.159 ms
Execution Time: 6.908 ms
Peak Memory Usage: 92 kB

```

### s16_page_offset [depth=1;spelling=offset]

Rows: 5; serialized decoded-row JSON bytes: 56 (not wire bytes). Planning/execution: 0.062000 / 5.724000 ms.

```text
Limit  (cost=1256.15..1256.16 rows=5 width=8) (actual time=5.678..5.681 rows=5 loops=1)
  ->  Sort  (cost=1256.15..1286.15 rows=12000 width=8) (actual time=5.676..5.677 rows=6 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 25kB
        ->  Seq Scan on orders  (cost=21.08..1041.05 rows=12000 width=8) (actual time=0.686..4.219 rows=12000 loops=1)
Planning Time: 0.062 ms
Execution Time: 5.724 ms
Peak Memory Usage: 92 kB

```

### s18_or_predicate [shape=broad/cast;spelling=or]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.067000 / 6.479000 ms.

```text
Sort  (cost=1443.14..1459.81 rows=6667 width=8) (actual time=5.853..6.032 rows=4652 loops=1)
  Sort Key: id
  Sort Method: quicksort  Memory: 411kB
  ->  Seq Scan on orders  (cost=21.08..1019.70 rows=6667 width=8) (actual time=1.396..4.854 rows=4652 loops=1)
        Storage Filter: (((amount)::bigint <= 100) OR ((amount)::bigint >= 400))
Planning Time: 0.067 ms
Execution Time: 6.479 ms
Peak Memory Usage: 548 kB

```

### s20_union_predicate [shape=broad/cast;spelling=union]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.145000 / 10.197000 ms.

```text
Sort  (cost=2516.53..2536.53 rows=8000 width=8) (actual time=9.545..9.704 rows=4652 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 411kB
  ->  HashAggregate  (cost=1917.90..1997.90 rows=8000 width=8) (actual time=8.211..8.627 rows=4652 loops=1)
        Group Key: orders.id, orders.amount
        Batches: 1  Memory Usage: 657kB
        ->  Append  (cost=21.08..1877.90 rows=8000 width=8) (actual time=2.079..7.228 rows=4652 loops=1)
              ->  Seq Scan on orders  (cost=21.08..878.95 rows=4000 width=8) (actual time=2.077..3.511 rows=2273 loops=1)
                    Storage Filter: ((amount)::bigint <= 100)
              ->  Seq Scan on orders orders_1  (cost=21.08..878.95 rows=4000 width=8) (actual time=1.471..3.103 rows=2379 loops=1)
                    Storage Filter: ((amount)::bigint >= 400)
Planning Time: 0.145 ms
Execution Time: 10.197 ms
Peak Memory Usage: 1252 kB

```

### s19_union_all_predicate [shape=broad/cast;spelling=union_all]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.098000 / 10.027000 ms.

```text
Sort  (cost=2396.53..2416.53 rows=8000 width=8) (actual time=8.995..9.267 rows=4652 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 411kB
  ->  Append  (cost=21.08..1877.90 rows=8000 width=8) (actual time=1.707..7.618 rows=4652 loops=1)
        ->  Seq Scan on orders  (cost=21.08..878.95 rows=4000 width=8) (actual time=1.705..3.681 rows=2273 loops=1)
              Storage Filter: ((amount)::bigint <= 100)
        ->  Seq Scan on orders orders_1  (cost=21.08..878.95 rows=4000 width=8) (actual time=1.835..3.345 rows=2379 loops=1)
              Storage Filter: ((amount)::bigint >= 400)
Planning Time: 0.098 ms
Execution Time: 10.027 ms
Peak Memory Usage: 625 kB

```

### s18_or_predicate [shape=broad/expression;index=expression_index;spelling=or]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.079000 / 7.813000 ms.

```text
Sort  (cost=1214.19..1224.72 rows=4211 width=8) (actual time=7.146..7.355 rows=4652 loops=1)
  Sort Key: id
  Sort Method: quicksort  Memory: 411kB
  ->  Seq Scan on orders  (cost=21.08..960.69 rows=4211 width=8) (actual time=1.273..6.211 rows=4652 loops=1)
        Storage Filter: (((amount + 0) <= 100) OR ((amount + 0) >= 400))
Planning Time: 0.079 ms
Execution Time: 7.813 ms
Peak Memory Usage: 548 kB

```

### s20_union_predicate [shape=broad/expression;index=expression_index;spelling=union]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.130000 / 11.713000 ms.

```text
Sort  (cost=2116.31..2127.97 rows=4664 width=8) (actual time=10.950..11.137 rows=4652 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 411kB
  ->  HashAggregate  (cost=1785.46..1832.10 rows=4664 width=8) (actual time=9.678..10.024 rows=4652 loops=1)
        Group Key: orders.id, orders.amount
        Batches: 1  Memory Usage: 465kB
        ->  Append  (cost=21.08..1762.14 rows=4664 width=8) (actual time=2.140..8.798 rows=4652 loops=1)
              ->  Seq Scan on orders  (cost=21.08..845.59 rows=2267 width=8) (actual time=2.139..3.496 rows=2273 loops=1)
                    Storage Filter: ((amount + 0) <= 100)
              ->  Seq Scan on orders orders_1  (cost=21.08..846.59 rows=2397 width=8) (actual time=3.092..4.725 rows=2379 loops=1)
                    Storage Filter: ((amount + 0) >= 400)
Planning Time: 0.130 ms
Execution Time: 11.713 ms
Peak Memory Usage: 1060 kB

```

### s19_union_all_predicate [shape=broad/expression;index=expression_index;spelling=union_all]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.111000 / 10.550000 ms.

```text
Sort  (cost=2046.35..2058.01 rows=4664 width=8) (actual time=9.864..10.085 rows=4652 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 411kB
  ->  Append  (cost=21.08..1762.14 rows=4664 width=8) (actual time=2.026..8.915 rows=4652 loops=1)
        ->  Seq Scan on orders  (cost=21.08..845.59 rows=2267 width=8) (actual time=2.024..3.536 rows=2273 loops=1)
              Storage Filter: ((amount + 0) <= 100)
        ->  Seq Scan on orders orders_1  (cost=21.08..846.59 rows=2397 width=8) (actual time=1.476..4.734 rows=2379 loops=1)
              Storage Filter: ((amount + 0) >= 400)
Planning Time: 0.111 ms
Execution Time: 10.550 ms
Peak Memory Usage: 625 kB

```

### s18_or_predicate [shape=broad/expression;spelling=or]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.074000 / 6.577000 ms.

```text
Sort  (cost=1443.14..1459.81 rows=6667 width=8) (actual time=5.863..6.052 rows=4652 loops=1)
  Sort Key: id
  Sort Method: quicksort  Memory: 411kB
  ->  Seq Scan on orders  (cost=21.08..1019.70 rows=6667 width=8) (actual time=1.387..4.880 rows=4652 loops=1)
        Storage Filter: (((amount + 0) <= 100) OR ((amount + 0) >= 400))
Planning Time: 0.074 ms
Execution Time: 6.577 ms
Peak Memory Usage: 548 kB

```

### s20_union_predicate [shape=broad/expression;spelling=union]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.087000 / 10.494000 ms.

```text
Sort  (cost=2516.53..2536.53 rows=8000 width=8) (actual time=9.711..9.942 rows=4652 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 411kB
  ->  HashAggregate  (cost=1917.90..1997.90 rows=8000 width=8) (actual time=8.447..8.824 rows=4652 loops=1)
        Group Key: orders.id, orders.amount
        Batches: 1  Memory Usage: 657kB
        ->  Append  (cost=21.08..1877.90 rows=8000 width=8) (actual time=2.003..7.504 rows=4652 loops=1)
              ->  Seq Scan on orders  (cost=21.08..878.95 rows=4000 width=8) (actual time=2.002..3.599 rows=2273 loops=1)
                    Storage Filter: ((amount + 0) <= 100)
              ->  Seq Scan on orders orders_1  (cost=21.08..878.95 rows=4000 width=8) (actual time=1.657..3.343 rows=2379 loops=1)
                    Storage Filter: ((amount + 0) >= 400)
Planning Time: 0.087 ms
Execution Time: 10.494 ms
Peak Memory Usage: 1252 kB

```

### s19_union_all_predicate [shape=broad/expression;spelling=union_all]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.099000 / 8.747000 ms.

```text
Sort  (cost=2396.53..2416.53 rows=8000 width=8) (actual time=8.052..8.226 rows=4652 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 411kB
  ->  Append  (cost=21.08..1877.90 rows=8000 width=8) (actual time=1.807..7.125 rows=4652 loops=1)
        ->  Seq Scan on orders  (cost=21.08..878.95 rows=4000 width=8) (actual time=1.806..3.579 rows=2273 loops=1)
              Storage Filter: ((amount + 0) <= 100)
        ->  Seq Scan on orders orders_1  (cost=21.08..878.95 rows=4000 width=8) (actual time=1.497..2.972 rows=2379 loops=1)
              Storage Filter: ((amount + 0) >= 400)
Planning Time: 0.099 ms
Execution Time: 8.747 ms
Peak Memory Usage: 625 kB

```

### s18_or_predicate [shape=broad/indexable;spelling=or]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.104000 / 5.605000 ms.

```text
Sort  (cost=1154.19..1164.72 rows=4211 width=8) (actual time=4.978..5.145 rows=4652 loops=1)
  Sort Key: id
  Sort Method: quicksort  Memory: 411kB
  ->  Seq Scan on orders  (cost=21.08..900.69 rows=4211 width=8) (actual time=1.312..3.814 rows=4652 loops=1)
        Storage Filter: ((amount <= 100) OR (amount >= 400))
Planning Time: 0.104 ms
Execution Time: 5.605 ms
Peak Memory Usage: 548 kB

```

### s20_union_predicate [shape=broad/indexable;spelling=union]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.103000 / 9.720000 ms.

```text
Sort  (cost=2056.31..2067.97 rows=4664 width=8) (actual time=9.063..9.221 rows=4652 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 411kB
  ->  HashAggregate  (cost=1725.46..1772.10 rows=4664 width=8) (actual time=7.841..8.159 rows=4652 loops=1)
        Group Key: orders.id, orders.amount
        Batches: 1  Memory Usage: 465kB
        ->  Append  (cost=21.08..1702.14 rows=4664 width=8) (actual time=2.060..6.721 rows=4652 loops=1)
              ->  Seq Scan on orders  (cost=21.08..815.59 rows=2267 width=8) (actual time=2.057..3.252 rows=2273 loops=1)
                    Storage Filter: (amount <= 100)
              ->  Seq Scan on orders orders_1  (cost=21.08..816.59 rows=2397 width=8) (actual time=1.469..2.814 rows=2379 loops=1)
                    Storage Filter: (amount >= 400)
Planning Time: 0.103 ms
Execution Time: 9.720 ms
Peak Memory Usage: 1060 kB

```

### s19_union_all_predicate [shape=broad/indexable;spelling=union_all]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.132000 / 9.896000 ms.

```text
Sort  (cost=1986.35..1998.01 rows=4664 width=8) (actual time=9.209..9.416 rows=4652 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 411kB
  ->  Append  (cost=21.08..1702.14 rows=4664 width=8) (actual time=2.809..8.279 rows=4652 loops=1)
        ->  Seq Scan on orders  (cost=21.08..815.59 rows=2267 width=8) (actual time=2.807..4.436 rows=2273 loops=1)
              Storage Filter: (amount <= 100)
        ->  Seq Scan on orders orders_1  (cost=21.08..816.59 rows=2397 width=8) (actual time=1.517..3.281 rows=2379 loops=1)
              Storage Filter: (amount >= 400)
Planning Time: 0.132 ms
Execution Time: 9.896 ms
Peak Memory Usage: 625 kB

```

### s18_or_predicate [shape=selective/cast;spelling=or]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.042000 / 3.495000 ms.

```text
Sort  (cost=1443.14..1459.81 rows=6667 width=8) (actual time=3.429..3.437 rows=251 loops=1)
  Sort Key: id
  Sort Method: quicksort  Memory: 36kB
  ->  Seq Scan on orders  (cost=21.08..1019.70 rows=6667 width=8) (actual time=3.334..3.354 rows=251 loops=1)
        Storage Filter: (((amount)::bigint <= 5) OR ((amount)::bigint >= 496))
Planning Time: 0.042 ms
Execution Time: 3.495 ms
Peak Memory Usage: 125 kB

```

### s20_union_predicate [shape=selective/cast;spelling=union]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.081000 / 6.085000 ms.

```text
Sort  (cost=2516.53..2536.53 rows=8000 width=8) (actual time=5.992..6.000 rows=251 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 36kB
  ->  HashAggregate  (cost=1917.90..1997.90 rows=8000 width=8) (actual time=5.906..5.943 rows=251 loops=1)
        Group Key: orders.id, orders.amount
        Batches: 1  Memory Usage: 417kB
        ->  Append  (cost=21.08..1877.90 rows=8000 width=8) (actual time=2.968..5.849 rows=251 loops=1)
              ->  Seq Scan on orders  (cost=21.08..878.95 rows=4000 width=8) (actual time=2.966..2.976 rows=114 loops=1)
                    Storage Filter: ((amount)::bigint <= 5)
              ->  Seq Scan on orders orders_1  (cost=21.08..878.95 rows=4000 width=8) (actual time=2.831..2.842 rows=137 loops=1)
                    Storage Filter: ((amount)::bigint >= 496)
Planning Time: 0.081 ms
Execution Time: 6.085 ms
Peak Memory Usage: 600 kB

```

### s19_union_all_predicate [shape=selective/cast;spelling=union_all]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.114000 / 6.457000 ms.

```text
Sort  (cost=2396.53..2416.53 rows=8000 width=8) (actual time=6.333..6.343 rows=251 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 36kB
  ->  Append  (cost=21.08..1877.90 rows=8000 width=8) (actual time=2.861..6.267 rows=251 loops=1)
        ->  Seq Scan on orders  (cost=21.08..878.95 rows=4000 width=8) (actual time=2.860..2.870 rows=114 loops=1)
              Storage Filter: ((amount)::bigint <= 5)
        ->  Seq Scan on orders orders_1  (cost=21.08..878.95 rows=4000 width=8) (actual time=3.348..3.361 rows=137 loops=1)
              Storage Filter: ((amount)::bigint >= 496)
Planning Time: 0.114 ms
Execution Time: 6.457 ms
Peak Memory Usage: 202 kB

```

### s18_or_predicate [shape=selective/expression;index=expression_index;spelling=or]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.081000 / 4.058000 ms.

```text
Sort  (cost=860.56..861.21 rows=262 width=8) (actual time=3.986..3.994 rows=251 loops=1)
  Sort Key: id
  Sort Method: quicksort  Memory: 36kB
  ->  Seq Scan on orders  (cost=21.08..850.03 rows=262 width=8) (actual time=3.891..3.915 rows=251 loops=1)
        Storage Filter: (((amount + 0) <= 5) OR ((amount + 0) >= 496))
Planning Time: 0.081 ms
Execution Time: 4.058 ms
Peak Memory Usage: 125 kB

```

### s20_union_predicate [shape=selective/expression;index=expression_index;spelling=union]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.125000 / 6.537000 ms.

```text
Sort  (cost=1596.53..1597.19 rows=263 width=8) (actual time=6.429..6.439 rows=251 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 36kB
  ->  HashAggregate  (cost=1583.33..1585.96 rows=263 width=8) (actual time=6.348..6.370 rows=251 loops=1)
        Group Key: orders.id, orders.amount
        Batches: 1  Memory Usage: 45kB
        ->  Append  (cost=21.08..1582.02 rows=263 width=8) (actual time=3.135..6.281 rows=251 loops=1)
              ->  Seq Scan on orders  (cost=21.08..788.85 rows=107 width=8) (actual time=3.134..3.145 rows=114 loops=1)
                    Storage Filter: ((amount + 0) <= 5)
              ->  Seq Scan on orders orders_1  (cost=21.08..789.22 rows=156 width=8) (actual time=3.083..3.097 rows=137 loops=1)
                    Storage Filter: ((amount + 0) >= 496)
Planning Time: 0.125 ms
Execution Time: 6.537 ms
Peak Memory Usage: 224 kB

```

### s19_union_all_predicate [shape=selective/expression;index=expression_index;spelling=union_all]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.111000 / 6.453000 ms.

```text
Sort  (cost=1592.59..1593.25 rows=263 width=8) (actual time=6.378..6.387 rows=251 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 36kB
  ->  Append  (cost=21.08..1582.02 rows=263 width=8) (actual time=3.027..6.314 rows=251 loops=1)
        ->  Seq Scan on orders  (cost=21.08..788.85 rows=107 width=8) (actual time=3.026..3.037 rows=114 loops=1)
              Storage Filter: ((amount + 0) <= 5)
        ->  Seq Scan on orders orders_1  (cost=21.08..789.22 rows=156 width=8) (actual time=3.229..3.242 rows=137 loops=1)
              Storage Filter: ((amount + 0) >= 496)
Planning Time: 0.111 ms
Execution Time: 6.453 ms
Peak Memory Usage: 202 kB

```

### s18_or_predicate [shape=selective/expression;spelling=or]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.065000 / 3.684000 ms.

```text
Sort  (cost=1443.14..1459.81 rows=6667 width=8) (actual time=3.616..3.625 rows=251 loops=1)
  Sort Key: id
  Sort Method: quicksort  Memory: 36kB
  ->  Seq Scan on orders  (cost=21.08..1019.70 rows=6667 width=8) (actual time=3.531..3.552 rows=251 loops=1)
        Storage Filter: (((amount + 0) <= 5) OR ((amount + 0) >= 496))
Planning Time: 0.065 ms
Execution Time: 3.684 ms
Peak Memory Usage: 125 kB

```

### s20_union_predicate [shape=selective/expression;spelling=union]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.083000 / 6.054000 ms.

```text
Sort  (cost=2516.53..2536.53 rows=8000 width=8) (actual time=5.917..5.927 rows=251 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 36kB
  ->  HashAggregate  (cost=1917.90..1997.90 rows=8000 width=8) (actual time=5.814..5.860 rows=251 loops=1)
        Group Key: orders.id, orders.amount
        Batches: 1  Memory Usage: 417kB
        ->  Append  (cost=21.08..1877.90 rows=8000 width=8) (actual time=2.932..5.749 rows=251 loops=1)
              ->  Seq Scan on orders  (cost=21.08..878.95 rows=4000 width=8) (actual time=2.930..2.941 rows=114 loops=1)
                    Storage Filter: ((amount + 0) <= 5)
              ->  Seq Scan on orders orders_1  (cost=21.08..878.95 rows=4000 width=8) (actual time=2.759..2.771 rows=137 loops=1)
                    Storage Filter: ((amount + 0) >= 496)
Planning Time: 0.083 ms
Execution Time: 6.054 ms
Peak Memory Usage: 600 kB

```

### s19_union_all_predicate [shape=selective/expression;spelling=union_all]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.107000 / 5.902000 ms.

```text
Sort  (cost=2396.53..2416.53 rows=8000 width=8) (actual time=5.821..5.832 rows=251 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 36kB
  ->  Append  (cost=21.08..1877.90 rows=8000 width=8) (actual time=2.922..5.751 rows=251 loops=1)
        ->  Seq Scan on orders  (cost=21.08..878.95 rows=4000 width=8) (actual time=2.919..2.931 rows=114 loops=1)
              Storage Filter: ((amount + 0) <= 5)
        ->  Seq Scan on orders orders_1  (cost=21.08..878.95 rows=4000 width=8) (actual time=2.768..2.781 rows=137 loops=1)
              Storage Filter: ((amount + 0) >= 496)
Planning Time: 0.107 ms
Execution Time: 5.902 ms
Peak Memory Usage: 202 kB

```

### s18_or_predicate [shape=selective/indexable;spelling=or]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.063000 / 3.199000 ms.

```text
Sort  (cost=800.56..801.21 rows=262 width=8) (actual time=3.133..3.141 rows=251 loops=1)
  Sort Key: id
  Sort Method: quicksort  Memory: 36kB
  ->  Seq Scan on orders  (cost=21.08..790.03 rows=262 width=8) (actual time=3.046..3.067 rows=251 loops=1)
        Storage Filter: ((amount <= 5) OR (amount >= 496))
Planning Time: 0.063 ms
Execution Time: 3.199 ms
Peak Memory Usage: 125 kB

```

### s18_or_predicate [shape=selective/indexable;spelling=or;position=last]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.074000 / 3.506000 ms.

```text
Sort  (cost=800.56..801.21 rows=262 width=8) (actual time=3.424..3.447 rows=251 loops=1)
  Sort Key: id
  Sort Method: quicksort  Memory: 36kB
  ->  Seq Scan on orders  (cost=21.08..790.03 rows=262 width=8) (actual time=3.329..3.352 rows=251 loops=1)
        Storage Filter: ((amount <= 5) OR (amount >= 496))
Planning Time: 0.074 ms
Execution Time: 3.506 ms
Peak Memory Usage: 125 kB

```

### s20_union_predicate [shape=selective/indexable;spelling=union]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.116000 / 5.818000 ms.

```text
Sort  (cost=1536.53..1537.19 rows=263 width=8) (actual time=5.731..5.741 rows=251 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 36kB
  ->  HashAggregate  (cost=1523.33..1525.96 rows=263 width=8) (actual time=5.654..5.674 rows=251 loops=1)
        Group Key: orders.id, orders.amount
        Batches: 1  Memory Usage: 45kB
        ->  Append  (cost=21.08..1522.02 rows=263 width=8) (actual time=2.882..5.588 rows=251 loops=1)
              ->  Seq Scan on orders  (cost=21.08..758.85 rows=107 width=8) (actual time=2.880..2.891 rows=114 loops=1)
                    Storage Filter: (amount <= 5)
              ->  Seq Scan on orders orders_1  (cost=21.08..759.22 rows=156 width=8) (actual time=2.650..2.662 rows=137 loops=1)
                    Storage Filter: (amount >= 496)
Planning Time: 0.116 ms
Execution Time: 5.818 ms
Peak Memory Usage: 224 kB

```

### s19_union_all_predicate [shape=selective/indexable;spelling=union_all]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.102000 / 5.772000 ms.

```text
Sort  (cost=1532.59..1533.25 rows=263 width=8) (actual time=5.699..5.708 rows=251 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 36kB
  ->  Append  (cost=21.08..1522.02 rows=263 width=8) (actual time=2.773..5.625 rows=251 loops=1)
        ->  Seq Scan on orders  (cost=21.08..758.85 rows=107 width=8) (actual time=2.772..2.782 rows=114 loops=1)
              Storage Filter: (amount <= 5)
        ->  Seq Scan on orders orders_1  (cost=21.08..759.22 rows=156 width=8) (actual time=2.798..2.809 rows=137 loops=1)
              Storage Filter: (amount >= 496)
Planning Time: 0.102 ms
Execution Time: 5.772 ms
Peak Memory Usage: 202 kB

```

## Limits and gaps recorded by the runner

- no concurrent writers: client-decode JSON bytes are not wire bytes; tails are under-sampled at the protocol trial count
- one host: no independent hardware, real-network, 3-node placement or equal-total-budget comparison is established
- single-node only; the YB 3-node RF=3 per-node/equal-total arms are a coverage gap on this host

