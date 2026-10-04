# SELECT access-path matrix (generated measurements)

_Measurements only. Verdicts for P10/P11 belong in a signed analysis._

**TL;DR:** within-run medians for OR / UNION ALL / UNION at three predicate shapes and OFFSET / keyset at five page depths, each cell verified against the independent oracle. A within-run difference is not a replicated result. Missing measurements are N/A.

| run | engine | topology | scale | page limit |
|---|---|---|---|---|
| select-order-access-yugabyte-1-20261004 | yugabyte | single | deep | 5 |

## Provenance

| Field | Value |
|---|---|
| repo commit | `ca2a38fbf9e53da00344578258138db9e01d3545` |
| benchmark image | `localhost/charitybench-select:select-order-access-yugabyte-1-20261004` |
| image digest | `8ecc060a8de698ebdf6b07ba9bd98c7663bdbb7dc6a74fa96b13a9609cd9664c` |
| run tag | `run/study-01/select-order-access-yugabyte-1-20261004` |
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

- default matrix: order seed `42118484`, 5 interleaved rounds
- expression-index arm: order seed `42309600`, 5 interleaved rounds
- default-matrix sequence (one sample of every cell per round):
  1. depth=1;spelling=offset
  2. shape=selective/cast;spelling=union
  3. shape=broad/expression;spelling=or
  4. shape=selective/indexable;spelling=or
  5. shape=broad/cast;spelling=or
  6. depth=1;spelling=keyset
  7. shape=broad/expression;spelling=union
  8. depth=100;spelling=offset
  9. shape=broad/indexable;spelling=or
  10. shape=selective/expression;spelling=or
  11. depth=1000;spelling=offset
  12. shape=broad/cast;spelling=union
  13. depth=10000;spelling=keyset
  14. shape=selective/expression;spelling=union
  15. shape=selective/cast;spelling=or
  16. shape=selective/expression;spelling=union_all
  17. shape=selective/cast;spelling=union_all
  18. depth=100;spelling=keyset
  19. shape=broad/cast;spelling=union_all
  20. shape=selective/indexable;spelling=union
  21. depth=1000;spelling=keyset
  22. shape=broad/expression;spelling=union_all
  23. depth=10;spelling=keyset
  24. depth=10;spelling=offset
  25. shape=broad/indexable;spelling=union
  26. depth=10000;spelling=offset
  27. shape=broad/indexable;spelling=union_all
  28. shape=selective/indexable;spelling=union_all
  29. depth=1;spelling=offset;position=last
- index-arm sequence:
  1. shape=broad/expression;index=expression_index;spelling=union
  2. shape=selective/expression;index=expression_index;spelling=union_all
  3. shape=broad/expression;index=expression_index;spelling=union_all
  4. shape=selective/expression;index=expression_index;spelling=or
  5. shape=broad/expression;index=expression_index;spelling=or
  6. shape=selective/expression;index=expression_index;spelling=union

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
| broad/cast | default | or | valid | 6.191972 | 12.1 | 5 | 4652 | `1dd363b32be056fc` |
| broad/cast | default | union | valid | 10.073319 | 5.4 | 5 | 4652 | `b5f6c293b9267d60` |
| broad/cast | default | union_all | valid | 9.171276 | 9.3 | 5 | 4652 | `24cabca42adfa691` |
| broad/expression | expression_index | or | valid | 4.811472 | 8.8 | 5 | 4652 | `0e731990206e1446` |
| broad/expression | expression_index | union | valid | 7.166819 | 8.5 | 5 | 4652 | `bf74591baa836f9d` |
| broad/expression | expression_index | union_all | valid | 6.639805 | 15.9 | 5 | 4652 | `a7e687668ac364de` |
| broad/expression | default | or | valid | 6.476607 | 8.1 | 5 | 4652 | `0e731990206e1446` |
| broad/expression | default | union | valid | 9.538274 | 5.7 | 5 | 4652 | `bf74591baa836f9d` |
| broad/expression | default | union_all | valid | 9.331759 | 17.5 | 5 | 4652 | `a7e687668ac364de` |
| broad/indexable | default | or | valid | 6.172420 | 22.5 | 5 | 4652 | `5191eeb70daef797` |
| broad/indexable | default | union | valid | 9.360888 | 14.0 | 5 | 4652 | `12de8f648f9520a0` |
| broad/indexable | default | union_all | valid | 8.929360 | 17.4 | 5 | 4652 | `d95e057d512f5ae5` |
| selective/cast | default | or | valid | 3.976737 | 10.4 | 5 | 251 | `57f71a4ae173cc97` |
| selective/cast | default | union | valid | 6.602031 | 9.2 | 5 | 251 | `2a31413606cd52a3` |
| selective/cast | default | union_all | valid | 6.503088 | 12.9 | 5 | 251 | `096e4d4e43a2a56e` |
| selective/expression | expression_index | or | valid | 3.094914 | 22.6 | 5 | 251 | `66a4eb5992604842` |
| selective/expression | expression_index | union | valid | 4.846780 | 27.8 | 5 | 251 | `14832cc886f611bd` |
| selective/expression | expression_index | union_all | valid | 4.969379 | 18.3 | 5 | 251 | `fee9cce2de640d1c` |
| selective/expression | default | or | valid | 4.125288 | 22.0 | 5 | 251 | `66a4eb5992604842` |
| selective/expression | default | union | valid | 6.553698 | 5.6 | 5 | 251 | `14832cc886f611bd` |
| selective/expression | default | union_all | valid | 6.374217 | 12.2 | 5 | 251 | `fee9cce2de640d1c` |
| selective/indexable | default | or | valid | 3.618328 | 9.5 | 5 | 251 | `66fcedd6207432cd` |
| selective/indexable | default | union | valid | 6.067404 | 14.2 | 5 | 251 | `40578cd947e4b93c` |
| selective/indexable | default | union_all | valid | 6.173958 | 4.3 | 5 | 251 | `8136cf54e3f9a863` |

### P10 pairs

Pairs compare the **default** arm only; the labelled `expression_index` arm is a package arm and is never pooled with the query-text pair.

| pair | A median ms | B median ms | B/A | same plan | within-run verdict |
|---|---|---|---|---|---|
| shape=selective/indexable or vs union_all | 3.618328 | 6.173958 | 1.7063 | false | within-run medians differ; whole-run replication required before a verdict |
| shape=selective/indexable or vs union | 3.618328 | 6.067404 | 1.6769 | false | within-run medians differ; whole-run replication required before a verdict |
| shape=selective/expression or vs union_all | 4.125288 | 6.374217 | 1.5452 | false | within-run medians differ; whole-run replication required before a verdict |
| shape=selective/expression or vs union | 4.125288 | 6.553698 | 1.5887 | false | within-run medians differ; whole-run replication required before a verdict |
| shape=selective/cast or vs union_all | 3.976737 | 6.503088 | 1.6353 | false | within-run medians differ; whole-run replication required before a verdict |
| shape=selective/cast or vs union | 3.976737 | 6.602031 | 1.6602 | false | within-run medians differ; whole-run replication required before a verdict |
| shape=broad/indexable or vs union_all | 6.172420 | 8.929360 | 1.4467 | false | within-run medians differ; whole-run replication required before a verdict |
| shape=broad/indexable or vs union | 6.172420 | 9.360888 | 1.5166 | false | within-run medians differ; whole-run replication required before a verdict |
| shape=broad/expression or vs union_all | 6.476607 | 9.331759 | 1.4408 | false | within-run medians differ; whole-run replication required before a verdict |
| shape=broad/expression or vs union | 6.476607 | 9.538274 | 1.4727 | false | within-run medians differ; whole-run replication required before a verdict |
| shape=broad/cast or vs union_all | 6.191972 | 9.171276 | 1.4812 | false | within-run medians differ; whole-run replication required before a verdict |
| shape=broad/cast or vs union | 6.191972 | 10.073319 | 1.6268 | false | within-run medians differ; whole-run replication required before a verdict |

## Pagination (P11)

| depth | spelling | status | median ms | spread % | trials | rows | plan structure |
|---|---|---|---|---|---|---|---|
| 1 | offset | valid | 4.433638 | 38.2 | 5 | 5 | `ac4509d295bf3348` |
| 1 | keyset | valid | 5.941060 | 34.7 | 5 | 5 | `166777c9921144af` |
| 10 | offset | valid | 4.695505 | 28.8 | 5 | 5 | `ac4509d295bf3348` |
| 10 | keyset | valid | 5.961666 | 41.6 | 5 | 5 | `d596ca0d42dc4da7` |
| 100 | offset | valid | 4.552814 | 13.1 | 5 | 5 | `ac4509d295bf3348` |
| 100 | keyset | valid | 5.943401 | 24.7 | 5 | 5 | `727a2962d7392ccd` |
| 1000 | offset | valid | 4.544106 | 19.6 | 5 | 5 | `ac4509d295bf3348` |
| 1000 | keyset | valid | 5.657415 | 23.3 | 5 | 5 | `441b4191cb138ed5` |
| 10000 | offset | valid | 6.098029 | 8.8 | 5 | 5 | `ac4509d295bf3348` |
| 10000 | keyset | valid | 4.348607 | 8.2 | 5 | 5 | `e1c72b3614c26105` |

Median ratio table (keyset / OFFSET; a ratio below 1 means keyset was faster in this run):

| depth | ratio (keyset/offset) | offset median ms | keyset median ms |
|---|---|---|---|
| 1 | 1.3400 | 4.433638 | 5.941060 |
| 10 | 1.2697 | 4.695505 | 5.961666 |
| 100 | 1.3054 | 4.552814 | 5.943401 |
| 1000 | 1.2450 | 4.544106 | 5.657415 |
| 10000 | 0.7131 | 6.098029 | 4.348607 |

### P11 pairs

| pair | A median ms | B median ms | B/A | same plan | within-run verdict |
|---|---|---|---|---|---|
| depth=1 offset vs keyset | 4.433638 | 5.941060 | 1.3400 | false | medians inside the two cells' own spread; no resolved difference within this run |
| depth=10 offset vs keyset | 4.695505 | 5.961666 | 1.2697 | false | medians inside the two cells' own spread; no resolved difference within this run |
| depth=100 offset vs keyset | 4.552814 | 5.943401 | 1.3054 | false | medians inside the two cells' own spread; no resolved difference within this run |
| depth=1000 offset vs keyset | 4.544106 | 5.657415 | 1.2450 | false | medians inside the two cells' own spread; no resolved difference within this run |
| depth=10000 offset vs keyset | 6.098029 | 4.348607 | 0.7131 | false | within-run medians differ; whole-run replication required before a verdict |

## Repeated-design position control

Not measured in this run.

## Expression-index arm: write and storage cost

Index `orders_amount_expr_idx` on expression `(amount + 0)`; supported=true.

| measure | value |
|---|---|
| index size (bytes) | 0 |
| insert batch with index (ms) | 14.482221 |
| insert batch without index (ms) | 19.701623 |
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

Rows: 5; serialized decoded-row JSON bytes: 51 (not wire bytes). Planning/execution: 0.094000 / 4.553000 ms.

```text
Limit  (cost=884.51..884.52 rows=5 width=8) (actual time=4.516..4.517 rows=5 loops=1)
  ->  Sort  (cost=884.51..889.29 rows=1912 width=8) (actual time=4.514..4.515 rows=5 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 25kB
        ->  Seq Scan on orders  (cost=21.08..852.75 rows=1912 width=8) (actual time=2.191..4.234 rows=2000 loops=1)
              Storage Filter: ((amount < 67) OR ((amount = 67) AND (id > 4310)) OR (amount IS NULL))
Planning Time: 0.094 ms
Execution Time: 4.553 ms
Peak Memory Usage: 92 kB

```

### s16_page_offset [depth=10000;spelling=offset]

Rows: 5; serialized decoded-row JSON bytes: 51 (not wire bytes). Planning/execution: 0.050000 / 7.246000 ms.

```text
Limit  (cost=1879.09..1879.10 rows=5 width=8) (actual time=7.202..7.205 rows=5 loops=1)
  ->  Sort  (cost=1854.09..1884.09 rows=12000 width=8) (actual time=5.852..6.254 rows=10005 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: quicksort  Memory: 947kB
        ->  Seq Scan on orders  (cost=21.08..1041.05 rows=12000 width=8) (actual time=0.475..2.720 rows=12000 loops=1)
Planning Time: 0.050 ms
Execution Time: 7.246 ms
Peak Memory Usage: 1516 kB

```

### s17_page_keyset [depth=1000;spelling=keyset]

Rows: 5; serialized decoded-row JSON bytes: 59 (not wire bytes). Planning/execution: 0.086000 / 5.763000 ms.

```text
Limit  (cost=1274.79..1274.81 rows=5 width=8) (actual time=5.728..5.730 rows=5 loops=1)
  ->  Sort  (cost=1274.79..1301.14 rows=10538 width=8) (actual time=5.727..5.728 rows=5 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 25kB
        ->  Seq Scan on orders  (cost=21.08..1099.76 rows=10538 width=8) (actual time=0.662..4.409 rows=11000 loops=1)
              Storage Filter: ((amount < 457) OR ((amount = 457) AND (id > 11254)) OR (amount IS NULL))
Planning Time: 0.086 ms
Execution Time: 5.763 ms
Peak Memory Usage: 92 kB

```

### s16_page_offset [depth=1000;spelling=offset]

Rows: 5; serialized decoded-row JSON bytes: 59 (not wire bytes). Planning/execution: 0.060000 / 4.759000 ms.

```text
Limit  (cost=1701.93..1701.94 rows=5 width=8) (actual time=4.722..4.724 rows=5 loops=1)
  ->  Sort  (cost=1699.43..1729.43 rows=12000 width=8) (actual time=4.597..4.629 rows=1005 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 130kB
        ->  Seq Scan on orders  (cost=21.08..1041.05 rows=12000 width=8) (actual time=0.543..2.471 rows=12000 loops=1)
Planning Time: 0.060 ms
Execution Time: 4.759 ms
Peak Memory Usage: 276 kB

```

### s17_page_keyset [depth=100;spelling=keyset]

Rows: 5; serialized decoded-row JSON bytes: 56 (not wire bytes). Planning/execution: 0.070000 / 5.444000 ms.

```text
Limit  (cost=1315.00..1315.01 rows=5 width=8) (actual time=5.409..5.412 rows=5 loops=1)
  ->  Sort  (cost=1315.00..1343.41 rows=11366 width=8) (actual time=5.408..5.409 rows=5 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 25kB
        ->  Seq Scan on orders  (cost=21.08..1126.21 rows=11366 width=8) (actual time=0.607..3.901 rows=11900 loops=1)
              Storage Filter: ((amount < 497) OR ((amount = 497) AND (id > 6723)) OR (amount IS NULL))
Planning Time: 0.070 ms
Execution Time: 5.444 ms
Peak Memory Usage: 92 kB

```

### s16_page_offset [depth=100;spelling=offset]

Rows: 5; serialized decoded-row JSON bytes: 56 (not wire bytes). Planning/execution: 0.048000 / 4.309000 ms.

```text
Limit  (cost=1504.15..1504.16 rows=5 width=8) (actual time=4.274..4.276 rows=5 loops=1)
  ->  Sort  (cost=1503.90..1533.90 rows=12000 width=8) (actual time=4.260..4.264 rows=105 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 33kB
        ->  Seq Scan on orders  (cost=21.08..1041.05 rows=12000 width=8) (actual time=0.541..2.673 rows=12000 loops=1)
Planning Time: 0.048 ms
Execution Time: 4.309 ms
Peak Memory Usage: 128 kB

```

### s17_page_keyset [depth=10;spelling=keyset]

Rows: 5; serialized decoded-row JSON bytes: 56 (not wire bytes). Planning/execution: 0.095000 / 5.642000 ms.

```text
Limit  (cost=1317.54..1317.55 rows=5 width=8) (actual time=5.606..5.608 rows=5 loops=1)
  ->  Sort  (cost=1317.54..1346.22 rows=11471 width=8) (actual time=5.605..5.605 rows=5 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 25kB
        ->  Seq Scan on orders  (cost=21.08..1127.01 rows=11471 width=8) (actual time=0.661..4.117 rows=11990 loops=1)
              Storage Filter: ((amount < 500) OR ((amount = 500) AND (id > 6319)) OR (amount IS NULL))
Planning Time: 0.095 ms
Execution Time: 5.642 ms
Peak Memory Usage: 92 kB

```

### s16_page_offset [depth=10;spelling=offset]

Rows: 5; serialized decoded-row JSON bytes: 56 (not wire bytes). Planning/execution: 0.059000 / 4.674000 ms.

```text
Limit  (cost=1335.49..1335.50 rows=5 width=8) (actual time=4.641..4.642 rows=5 loops=1)
  ->  Sort  (cost=1335.46..1365.46 rows=12000 width=8) (actual time=4.638..4.639 rows=15 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 26kB
        ->  Seq Scan on orders  (cost=21.08..1041.05 rows=12000 width=8) (actual time=0.477..3.131 rows=12000 loops=1)
Planning Time: 0.059 ms
Execution Time: 4.674 ms
Peak Memory Usage: 92 kB

```

### s17_page_keyset [depth=1;spelling=keyset]

Rows: 5; serialized decoded-row JSON bytes: 56 (not wire bytes). Planning/execution: 0.087000 / 6.574000 ms.

```text
Limit  (cost=1317.57..1317.58 rows=5 width=8) (actual time=6.533..6.535 rows=5 loops=1)
  ->  Sort  (cost=1317.57..1346.25 rows=11472 width=8) (actual time=6.532..6.532 rows=5 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 25kB
        ->  Seq Scan on orders  (cost=21.08..1127.02 rows=11472 width=8) (actual time=0.757..4.902 rows=11999 loops=1)
              Storage Filter: ((amount < 500) OR ((amount = 500) AND (id > 1554)) OR (amount IS NULL))
Planning Time: 0.087 ms
Execution Time: 6.574 ms
Peak Memory Usage: 92 kB

```

### s16_page_offset [depth=1;spelling=offset]

Rows: 5; serialized decoded-row JSON bytes: 56 (not wire bytes). Planning/execution: 0.051000 / 4.754000 ms.

```text
Limit  (cost=1256.15..1256.16 rows=5 width=8) (actual time=4.718..4.720 rows=5 loops=1)
  ->  Sort  (cost=1256.15..1286.15 rows=12000 width=8) (actual time=4.717..4.717 rows=6 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 25kB
        ->  Seq Scan on orders  (cost=21.08..1041.05 rows=12000 width=8) (actual time=0.545..3.077 rows=12000 loops=1)
Planning Time: 0.051 ms
Execution Time: 4.754 ms
Peak Memory Usage: 92 kB

```

### s16_page_offset [depth=1;spelling=offset;position=last]

Rows: 5; serialized decoded-row JSON bytes: 56 (not wire bytes). Planning/execution: 0.046000 / 4.025000 ms.

```text
Limit  (cost=1256.15..1256.16 rows=5 width=8) (actual time=3.994..3.996 rows=5 loops=1)
  ->  Sort  (cost=1256.15..1286.15 rows=12000 width=8) (actual time=3.993..3.994 rows=6 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 25kB
        ->  Seq Scan on orders  (cost=21.08..1041.05 rows=12000 width=8) (actual time=0.508..2.597 rows=12000 loops=1)
Planning Time: 0.046 ms
Execution Time: 4.025 ms
Peak Memory Usage: 92 kB

```

### s18_or_predicate [shape=broad/cast;spelling=or]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.057000 / 6.607000 ms.

```text
Sort  (cost=1443.14..1459.81 rows=6667 width=8) (actual time=5.908..6.107 rows=4652 loops=1)
  Sort Key: id
  Sort Method: quicksort  Memory: 411kB
  ->  Seq Scan on orders  (cost=21.08..1019.70 rows=6667 width=8) (actual time=1.166..4.870 rows=4652 loops=1)
        Storage Filter: (((amount)::bigint <= 100) OR ((amount)::bigint >= 400))
Planning Time: 0.057 ms
Execution Time: 6.607 ms
Peak Memory Usage: 548 kB

```

### s20_union_predicate [shape=broad/cast;spelling=union]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.088000 / 10.425000 ms.

```text
Sort  (cost=2516.53..2536.53 rows=8000 width=8) (actual time=9.620..9.827 rows=4652 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 411kB
  ->  HashAggregate  (cost=1917.90..1997.90 rows=8000 width=8) (actual time=7.944..8.317 rows=4652 loops=1)
        Group Key: orders.id, orders.amount
        Batches: 1  Memory Usage: 657kB
        ->  Append  (cost=21.08..1877.90 rows=8000 width=8) (actual time=1.738..6.988 rows=4652 loops=1)
              ->  Seq Scan on orders  (cost=21.08..878.95 rows=4000 width=8) (actual time=1.737..3.214 rows=2273 loops=1)
                    Storage Filter: ((amount)::bigint <= 100)
              ->  Seq Scan on orders orders_1  (cost=21.08..878.95 rows=4000 width=8) (actual time=1.548..3.154 rows=2379 loops=1)
                    Storage Filter: ((amount)::bigint >= 400)
Planning Time: 0.088 ms
Execution Time: 10.425 ms
Peak Memory Usage: 1260 kB

```

### s19_union_all_predicate [shape=broad/cast;spelling=union_all]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.093000 / 8.563000 ms.

```text
Sort  (cost=2396.53..2416.53 rows=8000 width=8) (actual time=7.892..8.075 rows=4652 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 411kB
  ->  Append  (cost=21.08..1877.90 rows=8000 width=8) (actual time=1.663..6.820 rows=4652 loops=1)
        ->  Seq Scan on orders  (cost=21.08..878.95 rows=4000 width=8) (actual time=1.662..3.209 rows=2273 loops=1)
              Storage Filter: ((amount)::bigint <= 100)
        ->  Seq Scan on orders orders_1  (cost=21.08..878.95 rows=4000 width=8) (actual time=1.542..2.956 rows=2379 loops=1)
              Storage Filter: ((amount)::bigint >= 400)
Planning Time: 0.093 ms
Execution Time: 8.563 ms
Peak Memory Usage: 625 kB

```

### s18_or_predicate [shape=broad/expression;index=expression_index;spelling=or]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.090000 / 4.931000 ms.

```text
Sort  (cost=1214.19..1224.72 rows=4211 width=8) (actual time=4.446..4.595 rows=4652 loops=1)
  Sort Key: id
  Sort Method: quicksort  Memory: 411kB
  ->  Seq Scan on orders  (cost=21.08..960.69 rows=4211 width=8) (actual time=1.161..3.658 rows=4652 loops=1)
        Storage Filter: (((amount + 0) <= 100) OR ((amount + 0) >= 400))
Planning Time: 0.090 ms
Execution Time: 4.931 ms
Peak Memory Usage: 548 kB

```

### s20_union_predicate [shape=broad/expression;index=expression_index;spelling=union]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.081000 / 7.383000 ms.

```text
Sort  (cost=2116.31..2127.97 rows=4664 width=8) (actual time=6.880..7.000 rows=4652 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 411kB
  ->  HashAggregate  (cost=1785.46..1832.10 rows=4664 width=8) (actual time=5.947..6.181 rows=4652 loops=1)
        Group Key: orders.id, orders.amount
        Batches: 1  Memory Usage: 465kB
        ->  Append  (cost=21.08..1762.14 rows=4664 width=8) (actual time=1.495..5.238 rows=4652 loops=1)
              ->  Seq Scan on orders  (cost=21.08..845.59 rows=2267 width=8) (actual time=1.494..2.664 rows=2273 loops=1)
                    Storage Filter: ((amount + 0) <= 100)
              ->  Seq Scan on orders orders_1  (cost=21.08..846.59 rows=2397 width=8) (actual time=1.053..2.127 rows=2379 loops=1)
                    Storage Filter: ((amount + 0) >= 400)
Planning Time: 0.081 ms
Execution Time: 7.383 ms
Peak Memory Usage: 1060 kB

```

### s19_union_all_predicate [shape=broad/expression;index=expression_index;spelling=union_all]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.110000 / 7.162000 ms.

```text
Sort  (cost=2046.35..2058.01 rows=4664 width=8) (actual time=6.487..6.707 rows=4652 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 411kB
  ->  Append  (cost=21.08..1762.14 rows=4664 width=8) (actual time=1.221..5.517 rows=4652 loops=1)
        ->  Seq Scan on orders  (cost=21.08..845.59 rows=2267 width=8) (actual time=1.219..2.437 rows=2273 loops=1)
              Storage Filter: ((amount + 0) <= 100)
        ->  Seq Scan on orders orders_1  (cost=21.08..846.59 rows=2397 width=8) (actual time=1.122..2.510 rows=2379 loops=1)
              Storage Filter: ((amount + 0) >= 400)
Planning Time: 0.110 ms
Execution Time: 7.162 ms
Peak Memory Usage: 625 kB

```

### s18_or_predicate [shape=broad/expression;spelling=or]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.061000 / 6.557000 ms.

```text
Sort  (cost=1443.14..1459.81 rows=6667 width=8) (actual time=5.886..6.073 rows=4652 loops=1)
  Sort Key: id
  Sort Method: quicksort  Memory: 411kB
  ->  Seq Scan on orders  (cost=21.08..1019.70 rows=6667 width=8) (actual time=1.173..4.736 rows=4652 loops=1)
        Storage Filter: (((amount + 0) <= 100) OR ((amount + 0) >= 400))
Planning Time: 0.061 ms
Execution Time: 6.557 ms
Peak Memory Usage: 548 kB

```

### s20_union_predicate [shape=broad/expression;spelling=union]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.084000 / 10.232000 ms.

```text
Sort  (cost=2516.53..2536.53 rows=8000 width=8) (actual time=9.518..9.694 rows=4652 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 411kB
  ->  HashAggregate  (cost=1917.90..1997.90 rows=8000 width=8) (actual time=8.065..8.465 rows=4652 loops=1)
        Group Key: orders.id, orders.amount
        Batches: 1  Memory Usage: 657kB
        ->  Append  (cost=21.08..1877.90 rows=8000 width=8) (actual time=1.849..7.169 rows=4652 loops=1)
              ->  Seq Scan on orders  (cost=21.08..878.95 rows=4000 width=8) (actual time=1.847..3.313 rows=2273 loops=1)
                    Storage Filter: ((amount + 0) <= 100)
              ->  Seq Scan on orders orders_1  (cost=21.08..878.95 rows=4000 width=8) (actual time=1.642..3.262 rows=2379 loops=1)
                    Storage Filter: ((amount + 0) >= 400)
Planning Time: 0.084 ms
Execution Time: 10.232 ms
Peak Memory Usage: 1252 kB

```

### s19_union_all_predicate [shape=broad/expression;spelling=union_all]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.111000 / 8.774000 ms.

```text
Sort  (cost=2396.53..2416.53 rows=8000 width=8) (actual time=8.115..8.293 rows=4652 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 411kB
  ->  Append  (cost=21.08..1877.90 rows=8000 width=8) (actual time=1.780..7.097 rows=4652 loops=1)
        ->  Seq Scan on orders  (cost=21.08..878.95 rows=4000 width=8) (actual time=1.778..3.439 rows=2273 loops=1)
              Storage Filter: ((amount + 0) <= 100)
        ->  Seq Scan on orders orders_1  (cost=21.08..878.95 rows=4000 width=8) (actual time=1.503..3.067 rows=2379 loops=1)
              Storage Filter: ((amount + 0) >= 400)
Planning Time: 0.111 ms
Execution Time: 8.774 ms
Peak Memory Usage: 625 kB

```

### s18_or_predicate [shape=broad/indexable;spelling=or]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.069000 / 6.110000 ms.

```text
Sort  (cost=1154.19..1164.72 rows=4211 width=8) (actual time=5.416..5.620 rows=4652 loops=1)
  Sort Key: id
  Sort Method: quicksort  Memory: 411kB
  ->  Seq Scan on orders  (cost=21.08..900.69 rows=4211 width=8) (actual time=1.299..4.396 rows=4652 loops=1)
        Storage Filter: ((amount <= 100) OR (amount >= 400))
Planning Time: 0.069 ms
Execution Time: 6.110 ms
Peak Memory Usage: 548 kB

```

### s20_union_predicate [shape=broad/indexable;spelling=union]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.095000 / 9.887000 ms.

```text
Sort  (cost=2056.31..2067.97 rows=4664 width=8) (actual time=8.842..9.017 rows=4652 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 411kB
  ->  HashAggregate  (cost=1725.46..1772.10 rows=4664 width=8) (actual time=7.541..7.883 rows=4652 loops=1)
        Group Key: orders.id, orders.amount
        Batches: 1  Memory Usage: 465kB
        ->  Append  (cost=21.08..1702.14 rows=4664 width=8) (actual time=1.750..6.636 rows=4652 loops=1)
              ->  Seq Scan on orders  (cost=21.08..815.59 rows=2267 width=8) (actual time=1.749..3.175 rows=2273 loops=1)
                    Storage Filter: (amount <= 100)
              ->  Seq Scan on orders orders_1  (cost=21.08..816.59 rows=2397 width=8) (actual time=1.469..2.863 rows=2379 loops=1)
                    Storage Filter: (amount >= 400)
Planning Time: 0.095 ms
Execution Time: 9.887 ms
Peak Memory Usage: 1060 kB

```

### s19_union_all_predicate [shape=broad/indexable;spelling=union_all]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.096000 / 8.203000 ms.

```text
Sort  (cost=1986.35..1998.01 rows=4664 width=8) (actual time=7.541..7.729 rows=4652 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 411kB
  ->  Append  (cost=21.08..1702.14 rows=4664 width=8) (actual time=1.647..6.557 rows=4652 loops=1)
        ->  Seq Scan on orders  (cost=21.08..815.59 rows=2267 width=8) (actual time=1.642..3.028 rows=2273 loops=1)
              Storage Filter: (amount <= 100)
        ->  Seq Scan on orders orders_1  (cost=21.08..816.59 rows=2397 width=8) (actual time=1.450..2.939 rows=2379 loops=1)
              Storage Filter: (amount >= 400)
Planning Time: 0.096 ms
Execution Time: 8.203 ms
Peak Memory Usage: 625 kB

```

### s18_or_predicate [shape=selective/cast;spelling=or]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.052000 / 3.781000 ms.

```text
Sort  (cost=1443.14..1459.81 rows=6667 width=8) (actual time=3.721..3.730 rows=251 loops=1)
  Sort Key: id
  Sort Method: quicksort  Memory: 36kB
  ->  Seq Scan on orders  (cost=21.08..1019.70 rows=6667 width=8) (actual time=3.641..3.662 rows=251 loops=1)
        Storage Filter: (((amount)::bigint <= 5) OR ((amount)::bigint >= 496))
Planning Time: 0.052 ms
Execution Time: 3.781 ms
Peak Memory Usage: 125 kB

```

### s20_union_predicate [shape=selective/cast;spelling=union]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.073000 / 6.931000 ms.

```text
Sort  (cost=2516.53..2536.53 rows=8000 width=8) (actual time=6.829..6.839 rows=251 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 36kB
  ->  HashAggregate  (cost=1917.90..1997.90 rows=8000 width=8) (actual time=6.719..6.769 rows=251 loops=1)
        Group Key: orders.id, orders.amount
        Batches: 1  Memory Usage: 417kB
        ->  Append  (cost=21.08..1877.90 rows=8000 width=8) (actual time=3.352..6.651 rows=251 loops=1)
              ->  Seq Scan on orders  (cost=21.08..878.95 rows=4000 width=8) (actual time=3.351..3.363 rows=114 loops=1)
                    Storage Filter: ((amount)::bigint <= 5)
              ->  Seq Scan on orders orders_1  (cost=21.08..878.95 rows=4000 width=8) (actual time=3.236..3.249 rows=137 loops=1)
                    Storage Filter: ((amount)::bigint >= 496)
Planning Time: 0.073 ms
Execution Time: 6.931 ms
Peak Memory Usage: 600 kB

```

### s19_union_all_predicate [shape=selective/cast;spelling=union_all]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.112000 / 6.647000 ms.

```text
Sort  (cost=2396.53..2416.53 rows=8000 width=8) (actual time=6.568..6.578 rows=251 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 36kB
  ->  Append  (cost=21.08..1877.90 rows=8000 width=8) (actual time=3.430..6.500 rows=251 loops=1)
        ->  Seq Scan on orders  (cost=21.08..878.95 rows=4000 width=8) (actual time=3.428..3.439 rows=114 loops=1)
              Storage Filter: ((amount)::bigint <= 5)
        ->  Seq Scan on orders orders_1  (cost=21.08..878.95 rows=4000 width=8) (actual time=3.015..3.026 rows=137 loops=1)
              Storage Filter: ((amount)::bigint >= 496)
Planning Time: 0.112 ms
Execution Time: 6.647 ms
Peak Memory Usage: 202 kB

```

### s18_or_predicate [shape=selective/expression;index=expression_index;spelling=or]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.037000 / 2.584000 ms.

```text
Sort  (cost=860.56..861.21 rows=262 width=8) (actual time=2.541..2.547 rows=251 loops=1)
  Sort Key: id
  Sort Method: quicksort  Memory: 36kB
  ->  Seq Scan on orders  (cost=21.08..850.03 rows=262 width=8) (actual time=2.486..2.500 rows=251 loops=1)
        Storage Filter: (((amount + 0) <= 5) OR ((amount + 0) >= 496))
Planning Time: 0.037 ms
Execution Time: 2.584 ms
Peak Memory Usage: 125 kB

```

### s20_union_predicate [shape=selective/expression;index=expression_index;spelling=union]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.056000 / 4.985000 ms.

```text
Sort  (cost=1596.53..1597.19 rows=263 width=8) (actual time=4.927..4.934 rows=251 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 36kB
  ->  HashAggregate  (cost=1583.33..1585.96 rows=263 width=8) (actual time=4.873..4.889 rows=251 loops=1)
        Group Key: orders.id, orders.amount
        Batches: 1  Memory Usage: 45kB
        ->  Append  (cost=21.08..1582.02 rows=263 width=8) (actual time=2.424..4.792 rows=251 loops=1)
              ->  Seq Scan on orders  (cost=21.08..788.85 rows=107 width=8) (actual time=2.422..2.448 rows=114 loops=1)
                    Storage Filter: ((amount + 0) <= 5)
              ->  Seq Scan on orders orders_1  (cost=21.08..789.22 rows=156 width=8) (actual time=2.280..2.290 rows=137 loops=1)
                    Storage Filter: ((amount + 0) >= 496)
Planning Time: 0.056 ms
Execution Time: 4.985 ms
Peak Memory Usage: 222 kB

```

### s19_union_all_predicate [shape=selective/expression;index=expression_index;spelling=union_all]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.068000 / 4.565000 ms.

```text
Sort  (cost=1592.59..1593.25 rows=263 width=8) (actual time=4.510..4.518 rows=251 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 36kB
  ->  Append  (cost=21.08..1582.02 rows=263 width=8) (actual time=2.265..4.466 rows=251 loops=1)
        ->  Seq Scan on orders  (cost=21.08..788.85 rows=107 width=8) (actual time=2.264..2.272 rows=114 loops=1)
              Storage Filter: ((amount + 0) <= 5)
        ->  Seq Scan on orders orders_1  (cost=21.08..789.22 rows=156 width=8) (actual time=2.160..2.169 rows=137 loops=1)
              Storage Filter: ((amount + 0) >= 496)
Planning Time: 0.068 ms
Execution Time: 4.565 ms
Peak Memory Usage: 202 kB

```

### s18_or_predicate [shape=selective/expression;spelling=or]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.044000 / 3.728000 ms.

```text
Sort  (cost=1443.14..1459.81 rows=6667 width=8) (actual time=3.667..3.676 rows=251 loops=1)
  Sort Key: id
  Sort Method: quicksort  Memory: 36kB
  ->  Seq Scan on orders  (cost=21.08..1019.70 rows=6667 width=8) (actual time=3.590..3.610 rows=251 loops=1)
        Storage Filter: (((amount + 0) <= 5) OR ((amount + 0) >= 496))
Planning Time: 0.044 ms
Execution Time: 3.728 ms
Peak Memory Usage: 125 kB

```

### s20_union_predicate [shape=selective/expression;spelling=union]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.077000 / 6.984000 ms.

```text
Sort  (cost=2516.53..2536.53 rows=8000 width=8) (actual time=6.870..6.880 rows=251 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 36kB
  ->  HashAggregate  (cost=1917.90..1997.90 rows=8000 width=8) (actual time=6.768..6.814 rows=251 loops=1)
        Group Key: orders.id, orders.amount
        Batches: 1  Memory Usage: 417kB
        ->  Append  (cost=21.08..1877.90 rows=8000 width=8) (actual time=3.281..6.705 rows=251 loops=1)
              ->  Seq Scan on orders  (cost=21.08..878.95 rows=4000 width=8) (actual time=3.280..3.290 rows=114 loops=1)
                    Storage Filter: ((amount + 0) <= 5)
              ->  Seq Scan on orders orders_1  (cost=21.08..878.95 rows=4000 width=8) (actual time=3.366..3.379 rows=137 loops=1)
                    Storage Filter: ((amount + 0) >= 496)
Planning Time: 0.077 ms
Execution Time: 6.984 ms
Peak Memory Usage: 600 kB

```

### s19_union_all_predicate [shape=selective/expression;spelling=union_all]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.075000 / 6.529000 ms.

```text
Sort  (cost=2396.53..2416.53 rows=8000 width=8) (actual time=6.432..6.443 rows=251 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 36kB
  ->  Append  (cost=21.08..1877.90 rows=8000 width=8) (actual time=3.162..6.363 rows=251 loops=1)
        ->  Seq Scan on orders  (cost=21.08..878.95 rows=4000 width=8) (actual time=3.161..3.172 rows=114 loops=1)
              Storage Filter: ((amount + 0) <= 5)
        ->  Seq Scan on orders orders_1  (cost=21.08..878.95 rows=4000 width=8) (actual time=3.138..3.151 rows=137 loops=1)
              Storage Filter: ((amount + 0) >= 496)
Planning Time: 0.075 ms
Execution Time: 6.529 ms
Peak Memory Usage: 202 kB

```

### s18_or_predicate [shape=selective/indexable;spelling=or]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.057000 / 3.704000 ms.

```text
Sort  (cost=800.56..801.21 rows=262 width=8) (actual time=3.625..3.636 rows=251 loops=1)
  Sort Key: id
  Sort Method: quicksort  Memory: 36kB
  ->  Seq Scan on orders  (cost=21.08..790.03 rows=262 width=8) (actual time=3.465..3.544 rows=251 loops=1)
        Storage Filter: ((amount <= 5) OR (amount >= 496))
Planning Time: 0.057 ms
Execution Time: 3.704 ms
Peak Memory Usage: 157 kB

```

### s20_union_predicate [shape=selective/indexable;spelling=union]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.081000 / 6.057000 ms.

```text
Sort  (cost=1536.53..1537.19 rows=263 width=8) (actual time=5.965..5.976 rows=251 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 36kB
  ->  HashAggregate  (cost=1523.33..1525.96 rows=263 width=8) (actual time=5.882..5.904 rows=251 loops=1)
        Group Key: orders.id, orders.amount
        Batches: 1  Memory Usage: 45kB
        ->  Append  (cost=21.08..1522.02 rows=263 width=8) (actual time=2.968..5.821 rows=251 loops=1)
              ->  Seq Scan on orders  (cost=21.08..758.85 rows=107 width=8) (actual time=2.967..2.978 rows=114 loops=1)
                    Storage Filter: (amount <= 5)
              ->  Seq Scan on orders orders_1  (cost=21.08..759.22 rows=156 width=8) (actual time=2.792..2.806 rows=137 loops=1)
                    Storage Filter: (amount >= 496)
Planning Time: 0.081 ms
Execution Time: 6.057 ms
Peak Memory Usage: 224 kB

```

### s19_union_all_predicate [shape=selective/indexable;spelling=union_all]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.099000 / 5.931000 ms.

```text
Sort  (cost=1532.59..1533.25 rows=263 width=8) (actual time=5.855..5.869 rows=251 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 36kB
  ->  Append  (cost=21.08..1522.02 rows=263 width=8) (actual time=3.087..5.800 rows=251 loops=1)
        ->  Seq Scan on orders  (cost=21.08..758.85 rows=107 width=8) (actual time=3.085..3.101 rows=114 loops=1)
              Storage Filter: (amount <= 5)
        ->  Seq Scan on orders orders_1  (cost=21.08..759.22 rows=156 width=8) (actual time=2.645..2.655 rows=137 loops=1)
              Storage Filter: (amount >= 496)
Planning Time: 0.099 ms
Execution Time: 5.931 ms
Peak Memory Usage: 202 kB

```

## Limits and gaps recorded by the runner

- no concurrent writers: client-decode JSON bytes are not wire bytes; tails are under-sampled at the protocol trial count
- one host: no independent hardware, real-network, 3-node placement or equal-total-budget comparison is established
- single-node only; the YB 3-node RF=3 per-node/equal-total arms are a coverage gap on this host

