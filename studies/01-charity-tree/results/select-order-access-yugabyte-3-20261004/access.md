# SELECT access-path matrix (generated measurements)

_Measurements only. Verdicts for P10/P11 belong in a signed analysis._

**TL;DR:** within-run medians for OR / UNION ALL / UNION at three predicate shapes and OFFSET / keyset at five page depths, each cell verified against the independent oracle. A within-run difference is not a replicated result. Missing measurements are N/A.

| run | engine | topology | scale | page limit |
|---|---|---|---|---|
| select-order-access-yugabyte-3-20261004 | yugabyte | single | deep | 5 |

## Provenance

| Field | Value |
|---|---|
| repo commit | `ca2a38fbf9e53da00344578258138db9e01d3545` |
| benchmark image | `localhost/charitybench-select:select-order-access-yugabyte-3-20261004` |
| image digest | `8ecc060a8de698ebdf6b07ba9bd98c7663bdbb7dc6a74fa96b13a9609cd9664c` |
| run tag | `run/study-01/select-order-access-yugabyte-3-20261004` |
| environment | `xiles-podman-select-order-correction-20261004` |
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

## Execution order

- default matrix: order seed `42316638`, 5 interleaved rounds
- expression-index arm: order seed `42732718`, 5 interleaved rounds
- default-matrix sequence (one sample of every cell per round):
  1. shape=selective/indexable;spelling=or
  2. shape=broad/cast;spelling=union_all
  3. shape=broad/expression;spelling=or
  4. shape=selective/expression;spelling=or
  5. depth=1;spelling=keyset
  6. shape=broad/cast;spelling=or
  7. shape=selective/indexable;spelling=union_all
  8. shape=broad/indexable;spelling=or
  9. depth=1000;spelling=keyset
  10. shape=selective/cast;spelling=or
  11. depth=1;spelling=offset
  12. shape=broad/indexable;spelling=union_all
  13. depth=10;spelling=keyset
  14. depth=10000;spelling=offset
  15. depth=10;spelling=offset
  16. shape=selective/indexable;spelling=union
  17. shape=selective/cast;spelling=union
  18. shape=broad/indexable;spelling=union
  19. depth=100;spelling=offset
  20. depth=1000;spelling=offset
  21. shape=selective/expression;spelling=union_all
  22. shape=broad/cast;spelling=union
  23. shape=broad/expression;spelling=union
  24. shape=selective/cast;spelling=union_all
  25. shape=selective/expression;spelling=union
  26. depth=100;spelling=keyset
  27. depth=10000;spelling=keyset
  28. shape=broad/expression;spelling=union_all
- index-arm sequence:
  1. shape=broad/expression;index=expression_index;spelling=union_all
  2. shape=broad/expression;index=expression_index;spelling=union
  3. shape=selective/expression;index=expression_index;spelling=union_all
  4. shape=broad/expression;index=expression_index;spelling=or
  5. shape=selective/expression;index=expression_index;spelling=union
  6. shape=selective/expression;index=expression_index;spelling=or

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
| broad/cast | default | or | valid | 6.541074 | 12.4 | 5 | 4652 | `1dd363b32be056fc` |
| broad/cast | default | union | valid | 9.874776 | 9.3 | 5 | 4652 | `b5f6c293b9267d60` |
| broad/cast | default | union_all | valid | 8.839165 | 10.7 | 5 | 4652 | `24cabca42adfa691` |
| broad/expression | expression_index | or | valid | 4.918317 | 9.2 | 5 | 4652 | `0e731990206e1446` |
| broad/expression | expression_index | union | valid | 7.451968 | 11.3 | 5 | 4652 | `bf74591baa836f9d` |
| broad/expression | expression_index | union_all | valid | 7.321621 | 10.7 | 5 | 4652 | `a7e687668ac364de` |
| broad/expression | default | or | valid | 6.576510 | 11.8 | 5 | 4652 | `0e731990206e1446` |
| broad/expression | default | union | valid | 9.553295 | 11.5 | 5 | 4652 | `bf74591baa836f9d` |
| broad/expression | default | union_all | valid | 9.077006 | 14.2 | 5 | 4652 | `a7e687668ac364de` |
| broad/indexable | default | or | valid | 6.002137 | 17.8 | 5 | 4652 | `5191eeb70daef797` |
| broad/indexable | default | union | valid | 9.221627 | 11.0 | 5 | 4652 | `12de8f648f9520a0` |
| broad/indexable | default | union_all | valid | 8.840938 | 9.8 | 5 | 4652 | `d95e057d512f5ae5` |
| selective/cast | default | or | valid | 4.024952 | 13.9 | 5 | 251 | `57f71a4ae173cc97` |
| selective/cast | default | union | valid | 6.523992 | 9.1 | 5 | 251 | `2a31413606cd52a3` |
| selective/cast | default | union_all | valid | 6.437079 | 13.1 | 5 | 251 | `096e4d4e43a2a56e` |
| selective/expression | expression_index | or | valid | 3.194597 | 15.3 | 5 | 251 | `66a4eb5992604842` |
| selective/expression | expression_index | union | valid | 5.084925 | 22.2 | 5 | 251 | `14832cc886f611bd` |
| selective/expression | expression_index | union_all | valid | 5.232050 | 20.4 | 5 | 251 | `fee9cce2de640d1c` |
| selective/expression | default | or | valid | 3.989447 | 8.9 | 5 | 251 | `66a4eb5992604842` |
| selective/expression | default | union | valid | 6.689376 | 11.2 | 5 | 251 | `14832cc886f611bd` |
| selective/expression | default | union_all | valid | 6.783315 | 15.7 | 5 | 251 | `fee9cce2de640d1c` |
| selective/indexable | default | or | valid | 3.612128 | 6.6 | 5 | 251 | `66fcedd6207432cd` |
| selective/indexable | default | union | valid | 6.175522 | 12.3 | 5 | 251 | `40578cd947e4b93c` |
| selective/indexable | default | union_all | valid | 6.247534 | 12.2 | 5 | 251 | `8136cf54e3f9a863` |

### P10 pairs

Pairs compare the **default** arm only; the labelled `expression_index` arm is a package arm and is never pooled with the query-text pair.

| pair | A median ms | B median ms | B/A | same plan | within-run verdict |
|---|---|---|---|---|---|
| shape=selective/indexable or vs union_all | 3.612128 | 6.247534 | 1.7296 | false | within-run medians differ; whole-run replication required before a verdict |
| shape=selective/indexable or vs union | 3.612128 | 6.175522 | 1.7097 | false | within-run medians differ; whole-run replication required before a verdict |
| shape=selective/expression or vs union_all | 3.989447 | 6.783315 | 1.7003 | false | within-run medians differ; whole-run replication required before a verdict |
| shape=selective/expression or vs union | 3.989447 | 6.689376 | 1.6768 | false | within-run medians differ; whole-run replication required before a verdict |
| shape=selective/cast or vs union_all | 4.024952 | 6.437079 | 1.5993 | false | within-run medians differ; whole-run replication required before a verdict |
| shape=selective/cast or vs union | 4.024952 | 6.523992 | 1.6209 | false | within-run medians differ; whole-run replication required before a verdict |
| shape=broad/indexable or vs union_all | 6.002137 | 8.840938 | 1.4730 | false | within-run medians differ; whole-run replication required before a verdict |
| shape=broad/indexable or vs union | 6.002137 | 9.221627 | 1.5364 | false | within-run medians differ; whole-run replication required before a verdict |
| shape=broad/expression or vs union_all | 6.576510 | 9.077006 | 1.3802 | false | within-run medians differ; whole-run replication required before a verdict |
| shape=broad/expression or vs union | 6.576510 | 9.553295 | 1.4526 | false | within-run medians differ; whole-run replication required before a verdict |
| shape=broad/cast or vs union_all | 6.541074 | 8.839165 | 1.3513 | false | within-run medians differ; whole-run replication required before a verdict |
| shape=broad/cast or vs union | 6.541074 | 9.874776 | 1.5097 | false | within-run medians differ; whole-run replication required before a verdict |

## Pagination (P11)

| depth | spelling | status | median ms | spread % | trials | rows | plan structure |
|---|---|---|---|---|---|---|---|
| 1 | offset | valid | 4.481934 | 25.1 | 5 | 5 | `ac4509d295bf3348` |
| 1 | keyset | valid | 6.102109 | 11.6 | 5 | 5 | `166777c9921144af` |
| 10 | offset | valid | 4.310410 | 39.0 | 5 | 5 | `ac4509d295bf3348` |
| 10 | keyset | valid | 5.698592 | 11.4 | 5 | 5 | `d596ca0d42dc4da7` |
| 100 | offset | valid | 4.783925 | 19.8 | 5 | 5 | `ac4509d295bf3348` |
| 100 | keyset | valid | 6.481624 | 24.3 | 5 | 5 | `727a2962d7392ccd` |
| 1000 | offset | valid | 4.946153 | 32.6 | 5 | 5 | `ac4509d295bf3348` |
| 1000 | keyset | valid | 5.944456 | 26.9 | 5 | 5 | `441b4191cb138ed5` |
| 10000 | offset | valid | 6.370120 | 18.2 | 5 | 5 | `ac4509d295bf3348` |
| 10000 | keyset | valid | 4.368005 | 5.9 | 5 | 5 | `e1c72b3614c26105` |

Median ratio table (keyset / OFFSET; a ratio below 1 means keyset was faster in this run):

| depth | ratio (keyset/offset) | offset median ms | keyset median ms |
|---|---|---|---|
| 1 | 1.3615 | 4.481934 | 6.102109 |
| 10 | 1.3221 | 4.310410 | 5.698592 |
| 100 | 1.3549 | 4.783925 | 6.481624 |
| 1000 | 1.2018 | 4.946153 | 5.944456 |
| 10000 | 0.6857 | 6.370120 | 4.368005 |

### P11 pairs

| pair | A median ms | B median ms | B/A | same plan | within-run verdict |
|---|---|---|---|---|---|
| depth=1 offset vs keyset | 4.481934 | 6.102109 | 1.3615 | false | medians inside the two cells' own spread; no resolved difference within this run |
| depth=10 offset vs keyset | 4.310410 | 5.698592 | 1.3221 | false | medians inside the two cells' own spread; no resolved difference within this run |
| depth=100 offset vs keyset | 4.783925 | 6.481624 | 1.3549 | false | medians inside the two cells' own spread; no resolved difference within this run |
| depth=1000 offset vs keyset | 4.946153 | 5.944456 | 1.2018 | false | medians inside the two cells' own spread; no resolved difference within this run |
| depth=10000 offset vs keyset | 6.370120 | 4.368005 | 0.6857 | false | within-run medians differ; whole-run replication required before a verdict |

## Repeated-design position control

Not measured in this run.

## Expression-index arm: write and storage cost

Index `orders_amount_expr_idx` on expression `(amount + 0)`; supported=true.

| measure | value |
|---|---|
| index size (bytes) | 0 |
| insert batch with index (ms) | 13.909441 |
| insert batch without index (ms) | 17.997118 |
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

- shape=selective/indexable `x05_aggregate_double_count`: detected=true, rows=200
- shape=selective/indexable `x06_topn_rank_ties`: detected=true, rows=210
- shape=selective/indexable `x01_optional_where`: detected=false, rows=9175 — x01_optional_where was NOT detected: its answer matched the sel_optional_match oracle
- shape=selective/indexable `x02_not_in_nullable`: detected=false, rows=0 — x02_not_in_nullable was NOT detected: its answer matched the sel_absence oracle
- shape=selective/indexable `x03_distinct_on_name`: detected=false, rows=200 — x03_distinct_on_name was NOT detected: its answer matched the sel_membership oracle
- shape=selective/indexable `x04_limit_before_join`: detected=true, rows=5
- shape=selective/indexable `x05_union_all_overlap`: detected=true, rows=228
- shape=selective/indexable `x07_membership_multiplier`: detected=true, rows=9175
- shape=selective/expression `x05_aggregate_double_count`: detected=true, rows=200
- shape=selective/expression `x06_topn_rank_ties`: detected=true, rows=210
- shape=selective/expression `x01_optional_where`: detected=false, rows=9175 — x01_optional_where was NOT detected: its answer matched the sel_optional_match oracle
- shape=selective/expression `x02_not_in_nullable`: detected=false, rows=0 — x02_not_in_nullable was NOT detected: its answer matched the sel_absence oracle
- shape=selective/expression `x03_distinct_on_name`: detected=false, rows=200 — x03_distinct_on_name was NOT detected: its answer matched the sel_membership oracle
- shape=selective/expression `x04_limit_before_join`: detected=true, rows=5
- shape=selective/expression `x05_union_all_overlap`: detected=true, rows=228
- shape=selective/expression `x07_membership_multiplier`: detected=true, rows=9175
- shape=selective/cast `x05_aggregate_double_count`: detected=true, rows=200
- shape=selective/cast `x06_topn_rank_ties`: detected=true, rows=210
- shape=selective/cast `x01_optional_where`: detected=false, rows=9175 — x01_optional_where was NOT detected: its answer matched the sel_optional_match oracle
- shape=selective/cast `x02_not_in_nullable`: detected=false, rows=0 — x02_not_in_nullable was NOT detected: its answer matched the sel_absence oracle
- shape=selective/cast `x03_distinct_on_name`: detected=false, rows=200 — x03_distinct_on_name was NOT detected: its answer matched the sel_membership oracle
- shape=selective/cast `x04_limit_before_join`: detected=true, rows=5
- shape=selective/cast `x05_union_all_overlap`: detected=true, rows=228
- shape=selective/cast `x07_membership_multiplier`: detected=true, rows=9175
- shape=broad/indexable `x05_aggregate_double_count`: detected=true, rows=200
- shape=broad/indexable `x06_topn_rank_ties`: detected=true, rows=210
- shape=broad/indexable `x01_optional_where`: detected=false, rows=9175 — x01_optional_where was NOT detected: its answer matched the sel_optional_match oracle
- shape=broad/indexable `x02_not_in_nullable`: detected=false, rows=0 — x02_not_in_nullable was NOT detected: its answer matched the sel_absence oracle
- shape=broad/indexable `x03_distinct_on_name`: detected=false, rows=200 — x03_distinct_on_name was NOT detected: its answer matched the sel_membership oracle
- shape=broad/indexable `x04_limit_before_join`: detected=true, rows=5
- shape=broad/indexable `x05_union_all_overlap`: detected=true, rows=4546
- shape=broad/indexable `x07_membership_multiplier`: detected=true, rows=9175
- shape=broad/expression `x05_aggregate_double_count`: detected=true, rows=200
- shape=broad/expression `x06_topn_rank_ties`: detected=true, rows=210
- shape=broad/expression `x01_optional_where`: detected=false, rows=9175 — x01_optional_where was NOT detected: its answer matched the sel_optional_match oracle
- shape=broad/expression `x02_not_in_nullable`: detected=false, rows=0 — x02_not_in_nullable was NOT detected: its answer matched the sel_absence oracle
- shape=broad/expression `x03_distinct_on_name`: detected=false, rows=200 — x03_distinct_on_name was NOT detected: its answer matched the sel_membership oracle
- shape=broad/expression `x04_limit_before_join`: detected=true, rows=5
- shape=broad/expression `x05_union_all_overlap`: detected=true, rows=4546
- shape=broad/expression `x07_membership_multiplier`: detected=true, rows=9175
- shape=broad/cast `x05_aggregate_double_count`: detected=true, rows=200
- shape=broad/cast `x06_topn_rank_ties`: detected=true, rows=210
- shape=broad/cast `x01_optional_where`: detected=false, rows=9175 — x01_optional_where was NOT detected: its answer matched the sel_optional_match oracle
- shape=broad/cast `x02_not_in_nullable`: detected=false, rows=0 — x02_not_in_nullable was NOT detected: its answer matched the sel_absence oracle
- shape=broad/cast `x03_distinct_on_name`: detected=false, rows=200 — x03_distinct_on_name was NOT detected: its answer matched the sel_membership oracle
- shape=broad/cast `x04_limit_before_join`: detected=true, rows=5
- shape=broad/cast `x05_union_all_overlap`: detected=true, rows=4546
- shape=broad/cast `x07_membership_multiplier`: detected=true, rows=9175

## Excluded cells

None: every cell is valid.

## Captured plans and result sizes

### s17_page_keyset [depth=10000;spelling=keyset]

Rows: 5; serialized decoded-row JSON bytes: 51 (not wire bytes). Planning/execution: 0.083000 / 4.058000 ms.

```text
Limit  (cost=884.51..884.52 rows=5 width=8) (actual time=4.023..4.025 rows=5 loops=1)
  ->  Sort  (cost=884.51..889.29 rows=1912 width=8) (actual time=4.021..4.022 rows=5 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 25kB
        ->  Seq Scan on orders  (cost=21.08..852.75 rows=1912 width=8) (actual time=2.056..3.763 rows=2000 loops=1)
              Storage Filter: ((amount < 67) OR ((amount = 67) AND (id > 4310)) OR (amount IS NULL))
Planning Time: 0.083 ms
Execution Time: 4.058 ms
Peak Memory Usage: 92 kB

```

### s16_page_offset [depth=10000;spelling=offset]

Rows: 5; serialized decoded-row JSON bytes: 51 (not wire bytes). Planning/execution: 0.050000 / 7.552000 ms.

```text
Limit  (cost=1879.09..1879.10 rows=5 width=8) (actual time=7.481..7.484 rows=5 loops=1)
  ->  Sort  (cost=1854.09..1884.09 rows=12000 width=8) (actual time=6.116..6.521 rows=10005 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: quicksort  Memory: 947kB
        ->  Seq Scan on orders  (cost=21.08..1041.05 rows=12000 width=8) (actual time=0.518..2.665 rows=12000 loops=1)
Planning Time: 0.050 ms
Execution Time: 7.552 ms
Peak Memory Usage: 1516 kB

```

### s17_page_keyset [depth=1000;spelling=keyset]

Rows: 5; serialized decoded-row JSON bytes: 59 (not wire bytes). Planning/execution: 0.088000 / 5.762000 ms.

```text
Limit  (cost=1274.79..1274.81 rows=5 width=8) (actual time=5.718..5.720 rows=5 loops=1)
  ->  Sort  (cost=1274.79..1301.14 rows=10538 width=8) (actual time=5.716..5.717 rows=5 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 25kB
        ->  Seq Scan on orders  (cost=21.08..1099.76 rows=10538 width=8) (actual time=0.751..4.297 rows=11000 loops=1)
              Storage Filter: ((amount < 457) OR ((amount = 457) AND (id > 11254)) OR (amount IS NULL))
Planning Time: 0.088 ms
Execution Time: 5.762 ms
Peak Memory Usage: 92 kB

```

### s16_page_offset [depth=1000;spelling=offset]

Rows: 5; serialized decoded-row JSON bytes: 59 (not wire bytes). Planning/execution: 0.094000 / 4.753000 ms.

```text
Limit  (cost=1701.93..1701.94 rows=5 width=8) (actual time=4.693..4.695 rows=5 loops=1)
  ->  Sort  (cost=1699.43..1729.43 rows=12000 width=8) (actual time=4.565..4.598 rows=1005 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 130kB
        ->  Seq Scan on orders  (cost=21.08..1041.05 rows=12000 width=8) (actual time=0.682..2.567 rows=12000 loops=1)
Planning Time: 0.094 ms
Execution Time: 4.753 ms
Peak Memory Usage: 276 kB

```

### s17_page_keyset [depth=100;spelling=keyset]

Rows: 5; serialized decoded-row JSON bytes: 56 (not wire bytes). Planning/execution: 0.095000 / 5.388000 ms.

```text
Limit  (cost=1315.00..1315.01 rows=5 width=8) (actual time=5.354..5.355 rows=5 loops=1)
  ->  Sort  (cost=1315.00..1343.41 rows=11366 width=8) (actual time=5.353..5.353 rows=5 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 25kB
        ->  Seq Scan on orders  (cost=21.08..1126.21 rows=11366 width=8) (actual time=0.646..3.907 rows=11900 loops=1)
              Storage Filter: ((amount < 497) OR ((amount = 497) AND (id > 6723)) OR (amount IS NULL))
Planning Time: 0.095 ms
Execution Time: 5.388 ms
Peak Memory Usage: 92 kB

```

### s16_page_offset [depth=100;spelling=offset]

Rows: 5; serialized decoded-row JSON bytes: 56 (not wire bytes). Planning/execution: 0.050000 / 4.092000 ms.

```text
Limit  (cost=1504.15..1504.16 rows=5 width=8) (actual time=4.061..4.062 rows=5 loops=1)
  ->  Sort  (cost=1503.90..1533.90 rows=12000 width=8) (actual time=4.047..4.051 rows=105 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 33kB
        ->  Seq Scan on orders  (cost=21.08..1041.05 rows=12000 width=8) (actual time=0.509..2.462 rows=12000 loops=1)
Planning Time: 0.050 ms
Execution Time: 4.092 ms
Peak Memory Usage: 128 kB

```

### s17_page_keyset [depth=10;spelling=keyset]

Rows: 5; serialized decoded-row JSON bytes: 56 (not wire bytes). Planning/execution: 0.065000 / 5.888000 ms.

```text
Limit  (cost=1317.54..1317.55 rows=5 width=8) (actual time=5.851..5.853 rows=5 loops=1)
  ->  Sort  (cost=1317.54..1346.22 rows=11471 width=8) (actual time=5.850..5.850 rows=5 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 25kB
        ->  Seq Scan on orders  (cost=21.08..1127.01 rows=11471 width=8) (actual time=0.762..4.297 rows=11990 loops=1)
              Storage Filter: ((amount < 500) OR ((amount = 500) AND (id > 6319)) OR (amount IS NULL))
Planning Time: 0.065 ms
Execution Time: 5.888 ms
Peak Memory Usage: 92 kB

```

### s16_page_offset [depth=10;spelling=offset]

Rows: 5; serialized decoded-row JSON bytes: 56 (not wire bytes). Planning/execution: 0.060000 / 5.307000 ms.

```text
Limit  (cost=1335.49..1335.50 rows=5 width=8) (actual time=5.266..5.268 rows=5 loops=1)
  ->  Sort  (cost=1335.46..1365.46 rows=12000 width=8) (actual time=5.263..5.264 rows=15 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 26kB
        ->  Seq Scan on orders  (cost=21.08..1041.05 rows=12000 width=8) (actual time=0.556..3.539 rows=12000 loops=1)
Planning Time: 0.060 ms
Execution Time: 5.307 ms
Peak Memory Usage: 92 kB

```

### s17_page_keyset [depth=1;spelling=keyset]

Rows: 5; serialized decoded-row JSON bytes: 56 (not wire bytes). Planning/execution: 0.114000 / 6.219000 ms.

```text
Limit  (cost=1317.57..1317.58 rows=5 width=8) (actual time=6.133..6.135 rows=5 loops=1)
  ->  Sort  (cost=1317.57..1346.25 rows=11472 width=8) (actual time=6.132..6.132 rows=5 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 25kB
        ->  Seq Scan on orders  (cost=21.08..1127.02 rows=11472 width=8) (actual time=0.963..4.502 rows=11999 loops=1)
              Storage Filter: ((amount < 500) OR ((amount = 500) AND (id > 1554)) OR (amount IS NULL))
Planning Time: 0.114 ms
Execution Time: 6.219 ms
Peak Memory Usage: 92 kB

```

### s16_page_offset [depth=1;spelling=offset]

Rows: 5; serialized decoded-row JSON bytes: 56 (not wire bytes). Planning/execution: 0.095000 / 4.879000 ms.

```text
Limit  (cost=1256.15..1256.16 rows=5 width=8) (actual time=4.841..4.842 rows=5 loops=1)
  ->  Sort  (cost=1256.15..1286.15 rows=12000 width=8) (actual time=4.839..4.840 rows=6 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 25kB
        ->  Seq Scan on orders  (cost=21.08..1041.05 rows=12000 width=8) (actual time=0.713..3.185 rows=12000 loops=1)
Planning Time: 0.095 ms
Execution Time: 4.879 ms
Peak Memory Usage: 92 kB

```

### s18_or_predicate [shape=broad/cast;spelling=or]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.085000 / 6.077000 ms.

```text
Sort  (cost=1443.14..1459.81 rows=6667 width=8) (actual time=5.371..5.593 rows=4652 loops=1)
  Sort Key: id
  Sort Method: quicksort  Memory: 411kB
  ->  Seq Scan on orders  (cost=21.08..1019.70 rows=6667 width=8) (actual time=1.237..4.383 rows=4652 loops=1)
        Storage Filter: (((amount)::bigint <= 100) OR ((amount)::bigint >= 400))
Planning Time: 0.085 ms
Execution Time: 6.077 ms
Peak Memory Usage: 548 kB

```

### s20_union_predicate [shape=broad/cast;spelling=union]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.076000 / 9.689000 ms.

```text
Sort  (cost=2516.53..2536.53 rows=8000 width=8) (actual time=9.018..9.180 rows=4652 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 411kB
  ->  HashAggregate  (cost=1917.90..1997.90 rows=8000 width=8) (actual time=7.674..8.049 rows=4652 loops=1)
        Group Key: orders.id, orders.amount
        Batches: 1  Memory Usage: 657kB
        ->  Append  (cost=21.08..1877.90 rows=8000 width=8) (actual time=1.773..6.727 rows=4652 loops=1)
              ->  Seq Scan on orders  (cost=21.08..878.95 rows=4000 width=8) (actual time=1.772..3.300 rows=2273 loops=1)
                    Storage Filter: ((amount)::bigint <= 100)
              ->  Seq Scan on orders orders_1  (cost=21.08..878.95 rows=4000 width=8) (actual time=1.482..2.803 rows=2379 loops=1)
                    Storage Filter: ((amount)::bigint >= 400)
Planning Time: 0.076 ms
Execution Time: 9.689 ms
Peak Memory Usage: 1252 kB

```

### s19_union_all_predicate [shape=broad/cast;spelling=union_all]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.104000 / 8.971000 ms.

```text
Sort  (cost=2396.53..2416.53 rows=8000 width=8) (actual time=8.280..8.480 rows=4652 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 411kB
  ->  Append  (cost=21.08..1877.90 rows=8000 width=8) (actual time=2.058..7.229 rows=4652 loops=1)
        ->  Seq Scan on orders  (cost=21.08..878.95 rows=4000 width=8) (actual time=2.057..3.592 rows=2273 loops=1)
              Storage Filter: ((amount)::bigint <= 100)
        ->  Seq Scan on orders orders_1  (cost=21.08..878.95 rows=4000 width=8) (actual time=1.536..2.982 rows=2379 loops=1)
              Storage Filter: ((amount)::bigint >= 400)
Planning Time: 0.104 ms
Execution Time: 8.971 ms
Peak Memory Usage: 625 kB

```

### s18_or_predicate [shape=broad/expression;index=expression_index;spelling=or]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.057000 / 4.925000 ms.

```text
Sort  (cost=1214.19..1224.72 rows=4211 width=8) (actual time=4.440..4.599 rows=4652 loops=1)
  Sort Key: id
  Sort Method: quicksort  Memory: 411kB
  ->  Seq Scan on orders  (cost=21.08..960.69 rows=4211 width=8) (actual time=0.932..2.930 rows=4652 loops=1)
        Storage Filter: (((amount + 0) <= 100) OR ((amount + 0) >= 400))
Planning Time: 0.057 ms
Execution Time: 4.925 ms
Peak Memory Usage: 548 kB

```

### s20_union_predicate [shape=broad/expression;index=expression_index;spelling=union]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.077000 / 7.373000 ms.

```text
Sort  (cost=2116.31..2127.97 rows=4664 width=8) (actual time=6.876..7.000 rows=4652 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 411kB
  ->  HashAggregate  (cost=1785.46..1832.10 rows=4664 width=8) (actual time=6.007..6.226 rows=4652 loops=1)
        Group Key: orders.id, orders.amount
        Batches: 1  Memory Usage: 465kB
        ->  Append  (cost=21.08..1762.14 rows=4664 width=8) (actual time=1.338..5.314 rows=4652 loops=1)
              ->  Seq Scan on orders  (cost=21.08..845.59 rows=2267 width=8) (actual time=1.334..2.600 rows=2273 loops=1)
                    Storage Filter: ((amount + 0) <= 100)
              ->  Seq Scan on orders orders_1  (cost=21.08..846.59 rows=2397 width=8) (actual time=1.170..2.257 rows=2379 loops=1)
                    Storage Filter: ((amount + 0) >= 400)
Planning Time: 0.077 ms
Execution Time: 7.373 ms
Peak Memory Usage: 1060 kB

```

### s19_union_all_predicate [shape=broad/expression;index=expression_index;spelling=union_all]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.105000 / 6.519000 ms.

```text
Sort  (cost=2046.35..2058.01 rows=4664 width=8) (actual time=6.055..6.187 rows=4652 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 411kB
  ->  Append  (cost=21.08..1762.14 rows=4664 width=8) (actual time=1.331..5.287 rows=4652 loops=1)
        ->  Seq Scan on orders  (cost=21.08..845.59 rows=2267 width=8) (actual time=1.330..2.423 rows=2273 loops=1)
              Storage Filter: ((amount + 0) <= 100)
        ->  Seq Scan on orders orders_1  (cost=21.08..846.59 rows=2397 width=8) (actual time=1.069..2.387 rows=2379 loops=1)
              Storage Filter: ((amount + 0) >= 400)
Planning Time: 0.105 ms
Execution Time: 6.519 ms
Peak Memory Usage: 625 kB

```

### s18_or_predicate [shape=broad/expression;spelling=or]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.094000 / 6.378000 ms.

```text
Sort  (cost=1443.14..1459.81 rows=6667 width=8) (actual time=5.674..5.893 rows=4652 loops=1)
  Sort Key: id
  Sort Method: quicksort  Memory: 411kB
  ->  Seq Scan on orders  (cost=21.08..1019.70 rows=6667 width=8) (actual time=1.384..4.590 rows=4652 loops=1)
        Storage Filter: (((amount + 0) <= 100) OR ((amount + 0) >= 400))
Planning Time: 0.094 ms
Execution Time: 6.378 ms
Peak Memory Usage: 548 kB

```

### s20_union_predicate [shape=broad/expression;spelling=union]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.118000 / 9.775000 ms.

```text
Sort  (cost=2516.53..2536.53 rows=8000 width=8) (actual time=9.106..9.270 rows=4652 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 411kB
  ->  HashAggregate  (cost=1917.90..1997.90 rows=8000 width=8) (actual time=7.670..8.185 rows=4652 loops=1)
        Group Key: orders.id, orders.amount
        Batches: 1  Memory Usage: 657kB
        ->  Append  (cost=21.08..1877.90 rows=8000 width=8) (actual time=2.020..6.691 rows=4652 loops=1)
              ->  Seq Scan on orders  (cost=21.08..878.95 rows=4000 width=8) (actual time=2.017..3.303 rows=2273 loops=1)
                    Storage Filter: ((amount + 0) <= 100)
              ->  Seq Scan on orders orders_1  (cost=21.08..878.95 rows=4000 width=8) (actual time=1.387..2.800 rows=2379 loops=1)
                    Storage Filter: ((amount + 0) >= 400)
Planning Time: 0.118 ms
Execution Time: 9.775 ms
Peak Memory Usage: 1252 kB

```

### s19_union_all_predicate [shape=broad/expression;spelling=union_all]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.087000 / 8.421000 ms.

```text
Sort  (cost=2396.53..2416.53 rows=8000 width=8) (actual time=7.742..7.930 rows=4652 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 411kB
  ->  Append  (cost=21.08..1877.90 rows=8000 width=8) (actual time=1.658..6.786 rows=4652 loops=1)
        ->  Seq Scan on orders  (cost=21.08..878.95 rows=4000 width=8) (actual time=1.656..3.232 rows=2273 loops=1)
              Storage Filter: ((amount + 0) <= 100)
        ->  Seq Scan on orders orders_1  (cost=21.08..878.95 rows=4000 width=8) (actual time=1.487..2.977 rows=2379 loops=1)
              Storage Filter: ((amount + 0) >= 400)
Planning Time: 0.087 ms
Execution Time: 8.421 ms
Peak Memory Usage: 625 kB

```

### s18_or_predicate [shape=broad/indexable;spelling=or]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.080000 / 5.673000 ms.

```text
Sort  (cost=1154.19..1164.72 rows=4211 width=8) (actual time=4.970..5.167 rows=4652 loops=1)
  Sort Key: id
  Sort Method: quicksort  Memory: 411kB
  ->  Seq Scan on orders  (cost=21.08..900.69 rows=4211 width=8) (actual time=1.211..3.879 rows=4652 loops=1)
        Storage Filter: ((amount <= 100) OR (amount >= 400))
Planning Time: 0.080 ms
Execution Time: 5.673 ms
Peak Memory Usage: 548 kB

```

### s20_union_predicate [shape=broad/indexable;spelling=union]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.097000 / 9.148000 ms.

```text
Sort  (cost=2056.31..2067.97 rows=4664 width=8) (actual time=8.370..8.562 rows=4652 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 411kB
  ->  HashAggregate  (cost=1725.46..1772.10 rows=4664 width=8) (actual time=7.055..7.411 rows=4652 loops=1)
        Group Key: orders.id, orders.amount
        Batches: 1  Memory Usage: 465kB
        ->  Append  (cost=21.08..1702.14 rows=4664 width=8) (actual time=1.626..6.183 rows=4652 loops=1)
              ->  Seq Scan on orders  (cost=21.08..815.59 rows=2267 width=8) (actual time=1.624..2.990 rows=2273 loops=1)
                    Storage Filter: (amount <= 100)
              ->  Seq Scan on orders orders_1  (cost=21.08..816.59 rows=2397 width=8) (actual time=1.345..2.613 rows=2379 loops=1)
                    Storage Filter: (amount >= 400)
Planning Time: 0.097 ms
Execution Time: 9.148 ms
Peak Memory Usage: 1060 kB

```

### s19_union_all_predicate [shape=broad/indexable;spelling=union_all]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.088000 / 7.973000 ms.

```text
Sort  (cost=1986.35..1998.01 rows=4664 width=8) (actual time=7.294..7.489 rows=4652 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 411kB
  ->  Append  (cost=21.08..1702.14 rows=4664 width=8) (actual time=1.622..6.234 rows=4652 loops=1)
        ->  Seq Scan on orders  (cost=21.08..815.59 rows=2267 width=8) (actual time=1.621..2.803 rows=2273 loops=1)
              Storage Filter: (amount <= 100)
        ->  Seq Scan on orders orders_1  (cost=21.08..816.59 rows=2397 width=8) (actual time=1.357..2.807 rows=2379 loops=1)
              Storage Filter: (amount >= 400)
Planning Time: 0.088 ms
Execution Time: 7.973 ms
Peak Memory Usage: 625 kB

```

### s18_or_predicate [shape=selective/cast;spelling=or]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.064000 / 4.054000 ms.

```text
Sort  (cost=1443.14..1459.81 rows=6667 width=8) (actual time=3.986..3.995 rows=251 loops=1)
  Sort Key: id
  Sort Method: quicksort  Memory: 36kB
  ->  Seq Scan on orders  (cost=21.08..1019.70 rows=6667 width=8) (actual time=3.892..3.915 rows=251 loops=1)
        Storage Filter: (((amount)::bigint <= 5) OR ((amount)::bigint >= 496))
Planning Time: 0.064 ms
Execution Time: 4.054 ms
Peak Memory Usage: 125 kB

```

### s20_union_predicate [shape=selective/cast;spelling=union]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.081000 / 6.371000 ms.

```text
Sort  (cost=2516.53..2536.53 rows=8000 width=8) (actual time=6.251..6.274 rows=251 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 36kB
  ->  HashAggregate  (cost=1917.90..1997.90 rows=8000 width=8) (actual time=6.116..6.191 rows=251 loops=1)
        Group Key: orders.id, orders.amount
        Batches: 1  Memory Usage: 417kB
        ->  Append  (cost=21.08..1877.90 rows=8000 width=8) (actual time=3.157..6.027 rows=251 loops=1)
              ->  Seq Scan on orders  (cost=21.08..878.95 rows=4000 width=8) (actual time=3.154..3.165 rows=114 loops=1)
                    Storage Filter: ((amount)::bigint <= 5)
              ->  Seq Scan on orders orders_1  (cost=21.08..878.95 rows=4000 width=8) (actual time=2.809..2.822 rows=137 loops=1)
                    Storage Filter: ((amount)::bigint >= 496)
Planning Time: 0.081 ms
Execution Time: 6.371 ms
Peak Memory Usage: 600 kB

```

### s19_union_all_predicate [shape=selective/cast;spelling=union_all]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.101000 / 6.176000 ms.

```text
Sort  (cost=2396.53..2416.53 rows=8000 width=8) (actual time=6.096..6.107 rows=251 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 36kB
  ->  Append  (cost=21.08..1877.90 rows=8000 width=8) (actual time=3.016..6.010 rows=251 loops=1)
        ->  Seq Scan on orders  (cost=21.08..878.95 rows=4000 width=8) (actual time=3.015..3.025 rows=114 loops=1)
              Storage Filter: ((amount)::bigint <= 5)
        ->  Seq Scan on orders orders_1  (cost=21.08..878.95 rows=4000 width=8) (actual time=2.936..2.949 rows=137 loops=1)
              Storage Filter: ((amount)::bigint >= 496)
Planning Time: 0.101 ms
Execution Time: 6.176 ms
Peak Memory Usage: 202 kB

```

### s18_or_predicate [shape=selective/expression;index=expression_index;spelling=or]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.040000 / 2.843000 ms.

```text
Sort  (cost=860.56..861.21 rows=262 width=8) (actual time=2.794..2.801 rows=251 loops=1)
  Sort Key: id
  Sort Method: quicksort  Memory: 36kB
  ->  Seq Scan on orders  (cost=21.08..850.03 rows=262 width=8) (actual time=2.723..2.741 rows=251 loops=1)
        Storage Filter: (((amount + 0) <= 5) OR ((amount + 0) >= 496))
Planning Time: 0.040 ms
Execution Time: 2.843 ms
Peak Memory Usage: 125 kB

```

### s20_union_predicate [shape=selective/expression;index=expression_index;spelling=union]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.063000 / 4.614000 ms.

```text
Sort  (cost=1596.53..1597.19 rows=263 width=8) (actual time=4.556..4.563 rows=251 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 36kB
  ->  HashAggregate  (cost=1583.33..1585.96 rows=263 width=8) (actual time=4.498..4.514 rows=251 loops=1)
        Group Key: orders.id, orders.amount
        Batches: 1  Memory Usage: 45kB
        ->  Append  (cost=21.08..1582.02 rows=263 width=8) (actual time=2.269..4.448 rows=251 loops=1)
              ->  Seq Scan on orders  (cost=21.08..788.85 rows=107 width=8) (actual time=2.265..2.273 rows=114 loops=1)
                    Storage Filter: ((amount + 0) <= 5)
              ->  Seq Scan on orders orders_1  (cost=21.08..789.22 rows=156 width=8) (actual time=2.134..2.144 rows=137 loops=1)
                    Storage Filter: ((amount + 0) >= 496)
Planning Time: 0.063 ms
Execution Time: 4.614 ms
Peak Memory Usage: 224 kB

```

### s19_union_all_predicate [shape=selective/expression;index=expression_index;spelling=union_all]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.086000 / 4.859000 ms.

```text
Sort  (cost=1592.59..1593.25 rows=263 width=8) (actual time=4.801..4.809 rows=251 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 36kB
  ->  Append  (cost=21.08..1582.02 rows=263 width=8) (actual time=2.539..4.752 rows=251 loops=1)
        ->  Seq Scan on orders  (cost=21.08..788.85 rows=107 width=8) (actual time=2.539..2.547 rows=114 loops=1)
              Storage Filter: ((amount + 0) <= 5)
        ->  Seq Scan on orders orders_1  (cost=21.08..789.22 rows=156 width=8) (actual time=2.167..2.176 rows=137 loops=1)
              Storage Filter: ((amount + 0) >= 496)
Planning Time: 0.086 ms
Execution Time: 4.859 ms
Peak Memory Usage: 200 kB

```

### s18_or_predicate [shape=selective/expression;spelling=or]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.052000 / 4.203000 ms.

```text
Sort  (cost=1443.14..1459.81 rows=6667 width=8) (actual time=4.131..4.140 rows=251 loops=1)
  Sort Key: id
  Sort Method: quicksort  Memory: 36kB
  ->  Seq Scan on orders  (cost=21.08..1019.70 rows=6667 width=8) (actual time=4.041..4.063 rows=251 loops=1)
        Storage Filter: (((amount + 0) <= 5) OR ((amount + 0) >= 496))
Planning Time: 0.052 ms
Execution Time: 4.203 ms
Peak Memory Usage: 125 kB

```

### s20_union_predicate [shape=selective/expression;spelling=union]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.080000 / 7.030000 ms.

```text
Sort  (cost=2516.53..2536.53 rows=8000 width=8) (actual time=6.888..6.899 rows=251 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 36kB
  ->  HashAggregate  (cost=1917.90..1997.90 rows=8000 width=8) (actual time=6.774..6.825 rows=251 loops=1)
        Group Key: orders.id, orders.amount
        Batches: 1  Memory Usage: 417kB
        ->  Append  (cost=21.08..1877.90 rows=8000 width=8) (actual time=3.514..6.700 rows=251 loops=1)
              ->  Seq Scan on orders  (cost=21.08..878.95 rows=4000 width=8) (actual time=3.513..3.525 rows=114 loops=1)
                    Storage Filter: ((amount + 0) <= 5)
              ->  Seq Scan on orders orders_1  (cost=21.08..878.95 rows=4000 width=8) (actual time=3.122..3.136 rows=137 loops=1)
                    Storage Filter: ((amount + 0) >= 496)
Planning Time: 0.080 ms
Execution Time: 7.030 ms
Peak Memory Usage: 600 kB

```

### s19_union_all_predicate [shape=selective/expression;spelling=union_all]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.087000 / 6.029000 ms.

```text
Sort  (cost=2396.53..2416.53 rows=8000 width=8) (actual time=5.959..5.969 rows=251 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 36kB
  ->  Append  (cost=21.08..1877.90 rows=8000 width=8) (actual time=2.945..5.883 rows=251 loops=1)
        ->  Seq Scan on orders  (cost=21.08..878.95 rows=4000 width=8) (actual time=2.944..2.954 rows=114 loops=1)
              Storage Filter: ((amount + 0) <= 5)
        ->  Seq Scan on orders orders_1  (cost=21.08..878.95 rows=4000 width=8) (actual time=2.882..2.894 rows=137 loops=1)
              Storage Filter: ((amount + 0) >= 496)
Planning Time: 0.087 ms
Execution Time: 6.029 ms
Peak Memory Usage: 202 kB

```

### s18_or_predicate [shape=selective/indexable;spelling=or]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.056000 / 3.554000 ms.

```text
Sort  (cost=800.56..801.21 rows=262 width=8) (actual time=3.484..3.494 rows=251 loops=1)
  Sort Key: id
  Sort Method: quicksort  Memory: 36kB
  ->  Seq Scan on orders  (cost=21.08..790.03 rows=262 width=8) (actual time=3.375..3.397 rows=251 loops=1)
        Storage Filter: ((amount <= 5) OR (amount >= 496))
Planning Time: 0.056 ms
Execution Time: 3.554 ms
Peak Memory Usage: 125 kB

```

### s20_union_predicate [shape=selective/indexable;spelling=union]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.095000 / 6.258000 ms.

```text
Sort  (cost=1536.53..1537.19 rows=263 width=8) (actual time=6.170..6.181 rows=251 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 36kB
  ->  HashAggregate  (cost=1523.33..1525.96 rows=263 width=8) (actual time=6.087..6.109 rows=251 loops=1)
        Group Key: orders.id, orders.amount
        Batches: 1  Memory Usage: 45kB
        ->  Append  (cost=21.08..1522.02 rows=263 width=8) (actual time=2.945..6.024 rows=251 loops=1)
              ->  Seq Scan on orders  (cost=21.08..758.85 rows=107 width=8) (actual time=2.944..2.954 rows=114 loops=1)
                    Storage Filter: (amount <= 5)
              ->  Seq Scan on orders orders_1  (cost=21.08..759.22 rows=156 width=8) (actual time=3.019..3.033 rows=137 loops=1)
                    Storage Filter: (amount >= 496)
Planning Time: 0.095 ms
Execution Time: 6.258 ms
Peak Memory Usage: 224 kB

```

### s19_union_all_predicate [shape=selective/indexable;spelling=union_all]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.082000 / 5.973000 ms.

```text
Sort  (cost=1532.59..1533.25 rows=263 width=8) (actual time=5.901..5.911 rows=251 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 36kB
  ->  Append  (cost=21.08..1522.02 rows=263 width=8) (actual time=2.967..5.830 rows=251 loops=1)
        ->  Seq Scan on orders  (cost=21.08..758.85 rows=107 width=8) (actual time=2.966..2.977 rows=114 loops=1)
              Storage Filter: (amount <= 5)
        ->  Seq Scan on orders orders_1  (cost=21.08..759.22 rows=156 width=8) (actual time=2.803..2.817 rows=137 loops=1)
              Storage Filter: (amount >= 496)
Planning Time: 0.082 ms
Execution Time: 5.973 ms
Peak Memory Usage: 202 kB

```

## Limits and gaps recorded by the runner

- no concurrent writers: client-decode JSON bytes are not wire bytes; tails are under-sampled at the protocol trial count
- one host: no independent hardware, real-network, 3-node placement or equal-total-budget comparison is established
- single-node only; the YB 3-node RF=3 per-node/equal-total arms are a coverage gap on this host

