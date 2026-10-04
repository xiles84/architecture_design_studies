# SELECT access-path matrix (generated measurements)

_Measurements only. Verdicts for P10/P11 belong in a signed analysis._

**TL;DR:** within-run medians for OR / UNION ALL / UNION at three predicate shapes and OFFSET / keyset at five page depths, each cell verified against the independent oracle. A within-run difference is not a replicated result. Missing measurements are N/A.

| run | engine | topology | scale | page limit |
|---|---|---|---|---|
| select-order-access-yugabyte-2-20261004 | yugabyte | single | deep | 5 |

## Provenance

| Field | Value |
|---|---|
| repo commit | `ca2a38fbf9e53da00344578258138db9e01d3545` |
| benchmark image | `localhost/charitybench-select:select-order-access-yugabyte-2-20261004` |
| image digest | `8ecc060a8de698ebdf6b07ba9bd98c7663bdbb7dc6a74fa96b13a9609cd9664c` |
| run tag | `run/study-01/select-order-access-yugabyte-2-20261004` |
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

- default matrix: order seed `42759799`, 5 interleaved rounds
- expression-index arm: order seed `42754307`, 5 interleaved rounds
- default-matrix sequence (one sample of every cell per round):
  1. depth=1000;spelling=keyset
  2. shape=broad/indexable;spelling=union
  3. depth=10;spelling=keyset
  4. depth=100;spelling=keyset
  5. shape=broad/expression;spelling=union_all
  6. shape=broad/expression;spelling=union
  7. depth=10000;spelling=keyset
  8. shape=broad/cast;spelling=union_all
  9. shape=selective/cast;spelling=union
  10. shape=selective/expression;spelling=union_all
  11. shape=broad/indexable;spelling=or
  12. depth=1000;spelling=offset
  13. shape=broad/expression;spelling=or
  14. shape=broad/indexable;spelling=union_all
  15. depth=1;spelling=offset
  16. shape=broad/cast;spelling=or
  17. depth=10;spelling=offset
  18. shape=selective/indexable;spelling=union
  19. shape=broad/cast;spelling=union
  20. shape=selective/cast;spelling=union_all
  21. depth=1;spelling=keyset
  22. shape=selective/indexable;spelling=union_all
  23. depth=10000;spelling=offset
  24. shape=selective/expression;spelling=or
  25. shape=selective/indexable;spelling=or
  26. shape=selective/cast;spelling=or
  27. depth=100;spelling=offset
  28. shape=selective/expression;spelling=union
- index-arm sequence:
  1. shape=broad/expression;index=expression_index;spelling=or
  2. shape=broad/expression;index=expression_index;spelling=union_all
  3. shape=broad/expression;index=expression_index;spelling=union
  4. shape=selective/expression;index=expression_index;spelling=or
  5. shape=selective/expression;index=expression_index;spelling=union_all
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
| broad/cast | default | or | valid | 6.324200 | 12.3 | 5 | 4652 | `1dd363b32be056fc` |
| broad/cast | default | union | valid | 8.814304 | 20.8 | 5 | 4652 | `b5f6c293b9267d60` |
| broad/cast | default | union_all | valid | 8.582023 | 6.4 | 5 | 4652 | `24cabca42adfa691` |
| broad/expression | expression_index | or | valid | 5.338468 | 19.6 | 5 | 4652 | `0e731990206e1446` |
| broad/expression | expression_index | union | valid | 7.478430 | 18.8 | 5 | 4652 | `bf74591baa836f9d` |
| broad/expression | expression_index | union_all | valid | 7.154042 | 10.1 | 5 | 4652 | `a7e687668ac364de` |
| broad/expression | default | or | valid | 6.550784 | 21.2 | 5 | 4652 | `0e731990206e1446` |
| broad/expression | default | union | valid | 9.132888 | 7.9 | 5 | 4652 | `bf74591baa836f9d` |
| broad/expression | default | union_all | valid | 8.647626 | 22.0 | 5 | 4652 | `a7e687668ac364de` |
| broad/indexable | default | or | valid | 5.758767 | 10.8 | 5 | 4652 | `5191eeb70daef797` |
| broad/indexable | default | union | valid | 8.673496 | 12.4 | 5 | 4652 | `12de8f648f9520a0` |
| broad/indexable | default | union_all | valid | 8.585931 | 8.2 | 5 | 4652 | `d95e057d512f5ae5` |
| selective/cast | default | or | valid | 3.641978 | 10.8 | 5 | 251 | `57f71a4ae173cc97` |
| selective/cast | default | union | valid | 6.269404 | 12.6 | 5 | 251 | `2a31413606cd52a3` |
| selective/cast | default | union_all | valid | 6.065570 | 4.5 | 5 | 251 | `096e4d4e43a2a56e` |
| selective/expression | expression_index | or | valid | 3.300704 | 12.3 | 5 | 251 | `66a4eb5992604842` |
| selective/expression | expression_index | union | valid | 5.508483 | 21.8 | 5 | 251 | `14832cc886f611bd` |
| selective/expression | expression_index | union_all | valid | 5.334283 | 6.7 | 5 | 251 | `fee9cce2de640d1c` |
| selective/expression | default | or | valid | 3.717254 | 18.2 | 5 | 251 | `66a4eb5992604842` |
| selective/expression | default | union | valid | 6.354504 | 10.5 | 5 | 251 | `14832cc886f611bd` |
| selective/expression | default | union_all | valid | 6.242359 | 11.3 | 5 | 251 | `fee9cce2de640d1c` |
| selective/indexable | default | or | valid | 3.213270 | 10.6 | 5 | 251 | `66fcedd6207432cd` |
| selective/indexable | default | union | valid | 5.996545 | 11.5 | 5 | 251 | `40578cd947e4b93c` |
| selective/indexable | default | union_all | valid | 5.658665 | 7.1 | 5 | 251 | `8136cf54e3f9a863` |

### P10 pairs

Pairs compare the **default** arm only; the labelled `expression_index` arm is a package arm and is never pooled with the query-text pair.

| pair | A median ms | B median ms | B/A | same plan | within-run verdict |
|---|---|---|---|---|---|
| shape=selective/indexable or vs union_all | 3.213270 | 5.658665 | 1.7610 | false | within-run medians differ; whole-run replication required before a verdict |
| shape=selective/indexable or vs union | 3.213270 | 5.996545 | 1.8662 | false | within-run medians differ; whole-run replication required before a verdict |
| shape=selective/expression or vs union_all | 3.717254 | 6.242359 | 1.6793 | false | within-run medians differ; whole-run replication required before a verdict |
| shape=selective/expression or vs union | 3.717254 | 6.354504 | 1.7095 | false | within-run medians differ; whole-run replication required before a verdict |
| shape=selective/cast or vs union_all | 3.641978 | 6.065570 | 1.6655 | false | within-run medians differ; whole-run replication required before a verdict |
| shape=selective/cast or vs union | 3.641978 | 6.269404 | 1.7214 | false | within-run medians differ; whole-run replication required before a verdict |
| shape=broad/indexable or vs union_all | 5.758767 | 8.585931 | 1.4909 | false | within-run medians differ; whole-run replication required before a verdict |
| shape=broad/indexable or vs union | 5.758767 | 8.673496 | 1.5061 | false | within-run medians differ; whole-run replication required before a verdict |
| shape=broad/expression or vs union_all | 6.550784 | 8.647626 | 1.3201 | false | medians inside the two cells' own spread; no resolved difference within this run |
| shape=broad/expression or vs union | 6.550784 | 9.132888 | 1.3942 | false | within-run medians differ; whole-run replication required before a verdict |
| shape=broad/cast or vs union_all | 6.324200 | 8.582023 | 1.3570 | false | within-run medians differ; whole-run replication required before a verdict |
| shape=broad/cast or vs union | 6.324200 | 8.814304 | 1.3937 | false | within-run medians differ; whole-run replication required before a verdict |

## Pagination (P11)

| depth | spelling | status | median ms | spread % | trials | rows | plan structure |
|---|---|---|---|---|---|---|---|
| 1 | offset | valid | 4.351427 | 14.0 | 5 | 5 | `ac4509d295bf3348` |
| 1 | keyset | valid | 5.535822 | 9.9 | 5 | 5 | `166777c9921144af` |
| 10 | offset | valid | 4.345947 | 9.4 | 5 | 5 | `ac4509d295bf3348` |
| 10 | keyset | valid | 5.750417 | 10.1 | 5 | 5 | `d596ca0d42dc4da7` |
| 100 | offset | valid | 4.316243 | 16.6 | 5 | 5 | `ac4509d295bf3348` |
| 100 | keyset | valid | 5.479473 | 6.1 | 5 | 5 | `727a2962d7392ccd` |
| 1000 | offset | valid | 4.399545 | 15.8 | 5 | 5 | `ac4509d295bf3348` |
| 1000 | keyset | valid | 5.515618 | 22.9 | 5 | 5 | `441b4191cb138ed5` |
| 10000 | offset | valid | 5.657803 | 11.6 | 5 | 5 | `ac4509d295bf3348` |
| 10000 | keyset | valid | 4.091508 | 19.3 | 5 | 5 | `e1c72b3614c26105` |

Median ratio table (keyset / OFFSET; a ratio below 1 means keyset was faster in this run):

| depth | ratio (keyset/offset) | offset median ms | keyset median ms |
|---|---|---|---|
| 1 | 1.2722 | 4.351427 | 5.535822 |
| 10 | 1.3232 | 4.345947 | 5.750417 |
| 100 | 1.2695 | 4.316243 | 5.479473 |
| 1000 | 1.2537 | 4.399545 | 5.515618 |
| 10000 | 0.7232 | 5.657803 | 4.091508 |

### P11 pairs

| pair | A median ms | B median ms | B/A | same plan | within-run verdict |
|---|---|---|---|---|---|
| depth=1 offset vs keyset | 4.351427 | 5.535822 | 1.2722 | false | within-run medians differ; whole-run replication required before a verdict |
| depth=10 offset vs keyset | 4.345947 | 5.750417 | 1.3232 | false | within-run medians differ; whole-run replication required before a verdict |
| depth=100 offset vs keyset | 4.316243 | 5.479473 | 1.2695 | false | within-run medians differ; whole-run replication required before a verdict |
| depth=1000 offset vs keyset | 4.399545 | 5.515618 | 1.2537 | false | medians inside the two cells' own spread; no resolved difference within this run |
| depth=10000 offset vs keyset | 5.657803 | 4.091508 | 0.7232 | false | medians inside the two cells' own spread; no resolved difference within this run |

## Repeated-design position control

Not measured in this run.

## Expression-index arm: write and storage cost

Index `orders_amount_expr_idx` on expression `(amount + 0)`; supported=true.

| measure | value |
|---|---|
| index size (bytes) | 0 |
| insert batch with index (ms) | 13.616289 |
| insert batch without index (ms) | 16.602277 |
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

Rows: 5; serialized decoded-row JSON bytes: 51 (not wire bytes). Planning/execution: 0.080000 / 4.188000 ms.

```text
Limit  (cost=884.51..884.52 rows=5 width=8) (actual time=4.153..4.155 rows=5 loops=1)
  ->  Sort  (cost=884.51..889.29 rows=1912 width=8) (actual time=4.152..4.153 rows=5 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 25kB
        ->  Seq Scan on orders  (cost=21.08..852.75 rows=1912 width=8) (actual time=1.967..3.882 rows=2000 loops=1)
              Storage Filter: ((amount < 67) OR ((amount = 67) AND (id > 4310)) OR (amount IS NULL))
Planning Time: 0.080 ms
Execution Time: 4.188 ms
Peak Memory Usage: 92 kB

```

### s16_page_offset [depth=10000;spelling=offset]

Rows: 5; serialized decoded-row JSON bytes: 51 (not wire bytes). Planning/execution: 0.047000 / 7.060000 ms.

```text
Limit  (cost=1879.09..1879.10 rows=5 width=8) (actual time=7.020..7.022 rows=5 loops=1)
  ->  Sort  (cost=1854.09..1884.09 rows=12000 width=8) (actual time=5.800..6.141 rows=10005 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: quicksort  Memory: 947kB
        ->  Seq Scan on orders  (cost=21.08..1041.05 rows=12000 width=8) (actual time=0.497..2.555 rows=12000 loops=1)
Planning Time: 0.047 ms
Execution Time: 7.060 ms
Peak Memory Usage: 1516 kB

```

### s17_page_keyset [depth=1000;spelling=keyset]

Rows: 5; serialized decoded-row JSON bytes: 59 (not wire bytes). Planning/execution: 0.089000 / 5.853000 ms.

```text
Limit  (cost=1274.79..1274.81 rows=5 width=8) (actual time=5.812..5.815 rows=5 loops=1)
  ->  Sort  (cost=1274.79..1301.14 rows=10538 width=8) (actual time=5.811..5.812 rows=5 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 25kB
        ->  Seq Scan on orders  (cost=21.08..1099.76 rows=10538 width=8) (actual time=0.717..4.319 rows=11000 loops=1)
              Storage Filter: ((amount < 457) OR ((amount = 457) AND (id > 11254)) OR (amount IS NULL))
Planning Time: 0.089 ms
Execution Time: 5.853 ms
Peak Memory Usage: 92 kB

```

### s16_page_offset [depth=1000;spelling=offset]

Rows: 5; serialized decoded-row JSON bytes: 59 (not wire bytes). Planning/execution: 0.065000 / 4.247000 ms.

```text
Limit  (cost=1701.93..1701.94 rows=5 width=8) (actual time=4.041..4.042 rows=5 loops=1)
  ->  Sort  (cost=1699.43..1729.43 rows=12000 width=8) (actual time=3.918..3.949 rows=1005 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 130kB
        ->  Seq Scan on orders  (cost=21.08..1041.05 rows=12000 width=8) (actual time=0.516..1.928 rows=12000 loops=1)
Planning Time: 0.065 ms
Execution Time: 4.247 ms
Peak Memory Usage: 276 kB

```

### s17_page_keyset [depth=100;spelling=keyset]

Rows: 5; serialized decoded-row JSON bytes: 56 (not wire bytes). Planning/execution: 0.075000 / 6.429000 ms.

```text
Limit  (cost=1315.00..1315.01 rows=5 width=8) (actual time=6.386..6.388 rows=5 loops=1)
  ->  Sort  (cost=1315.00..1343.41 rows=11366 width=8) (actual time=6.384..6.385 rows=5 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 25kB
        ->  Seq Scan on orders  (cost=21.08..1126.21 rows=11366 width=8) (actual time=0.696..4.856 rows=11900 loops=1)
              Storage Filter: ((amount < 497) OR ((amount = 497) AND (id > 6723)) OR (amount IS NULL))
Planning Time: 0.075 ms
Execution Time: 6.429 ms
Peak Memory Usage: 92 kB

```

### s16_page_offset [depth=100;spelling=offset]

Rows: 5; serialized decoded-row JSON bytes: 56 (not wire bytes). Planning/execution: 0.055000 / 3.585000 ms.

```text
Limit  (cost=1504.15..1504.16 rows=5 width=8) (actual time=3.550..3.551 rows=5 loops=1)
  ->  Sort  (cost=1503.90..1533.90 rows=12000 width=8) (actual time=3.535..3.539 rows=105 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 33kB
        ->  Seq Scan on orders  (cost=21.08..1041.05 rows=12000 width=8) (actual time=0.474..2.103 rows=12000 loops=1)
Planning Time: 0.055 ms
Execution Time: 3.585 ms
Peak Memory Usage: 128 kB

```

### s17_page_keyset [depth=10;spelling=keyset]

Rows: 5; serialized decoded-row JSON bytes: 56 (not wire bytes). Planning/execution: 0.105000 / 6.147000 ms.

```text
Limit  (cost=1317.54..1317.55 rows=5 width=8) (actual time=6.097..6.099 rows=5 loops=1)
  ->  Sort  (cost=1317.54..1346.22 rows=11471 width=8) (actual time=6.096..6.096 rows=5 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 25kB
        ->  Seq Scan on orders  (cost=21.08..1127.01 rows=11471 width=8) (actual time=0.835..4.584 rows=11990 loops=1)
              Storage Filter: ((amount < 500) OR ((amount = 500) AND (id > 6319)) OR (amount IS NULL))
Planning Time: 0.105 ms
Execution Time: 6.147 ms
Peak Memory Usage: 92 kB

```

### s16_page_offset [depth=10;spelling=offset]

Rows: 5; serialized decoded-row JSON bytes: 56 (not wire bytes). Planning/execution: 0.057000 / 4.855000 ms.

```text
Limit  (cost=1335.49..1335.50 rows=5 width=8) (actual time=4.818..4.819 rows=5 loops=1)
  ->  Sort  (cost=1335.46..1365.46 rows=12000 width=8) (actual time=4.815..4.816 rows=15 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 26kB
        ->  Seq Scan on orders  (cost=21.08..1041.05 rows=12000 width=8) (actual time=0.566..3.074 rows=12000 loops=1)
Planning Time: 0.057 ms
Execution Time: 4.855 ms
Peak Memory Usage: 92 kB

```

### s17_page_keyset [depth=1;spelling=keyset]

Rows: 5; serialized decoded-row JSON bytes: 56 (not wire bytes). Planning/execution: 0.103000 / 5.741000 ms.

```text
Limit  (cost=1317.57..1317.58 rows=5 width=8) (actual time=5.696..5.698 rows=5 loops=1)
  ->  Sort  (cost=1317.57..1346.25 rows=11472 width=8) (actual time=5.695..5.696 rows=5 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 25kB
        ->  Seq Scan on orders  (cost=21.08..1127.02 rows=11472 width=8) (actual time=0.711..4.282 rows=11999 loops=1)
              Storage Filter: ((amount < 500) OR ((amount = 500) AND (id > 1554)) OR (amount IS NULL))
Planning Time: 0.103 ms
Execution Time: 5.741 ms
Peak Memory Usage: 92 kB

```

### s16_page_offset [depth=1;spelling=offset]

Rows: 5; serialized decoded-row JSON bytes: 56 (not wire bytes). Planning/execution: 0.089000 / 4.323000 ms.

```text
Limit  (cost=1256.15..1256.16 rows=5 width=8) (actual time=4.288..4.289 rows=5 loops=1)
  ->  Sort  (cost=1256.15..1286.15 rows=12000 width=8) (actual time=4.287..4.287 rows=6 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 25kB
        ->  Seq Scan on orders  (cost=21.08..1041.05 rows=12000 width=8) (actual time=0.596..2.872 rows=12000 loops=1)
Planning Time: 0.089 ms
Execution Time: 4.323 ms
Peak Memory Usage: 92 kB

```

### s18_or_predicate [shape=broad/cast;spelling=or]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.064000 / 5.530000 ms.

```text
Sort  (cost=1443.14..1459.81 rows=6667 width=8) (actual time=4.925..5.093 rows=4652 loops=1)
  Sort Key: id
  Sort Method: quicksort  Memory: 411kB
  ->  Seq Scan on orders  (cost=21.08..1019.70 rows=6667 width=8) (actual time=1.215..4.020 rows=4652 loops=1)
        Storage Filter: (((amount)::bigint <= 100) OR ((amount)::bigint >= 400))
Planning Time: 0.064 ms
Execution Time: 5.530 ms
Peak Memory Usage: 564 kB

```

### s20_union_predicate [shape=broad/cast;spelling=union]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.085000 / 9.250000 ms.

```text
Sort  (cost=2516.53..2536.53 rows=8000 width=8) (actual time=8.484..8.728 rows=4652 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 411kB
  ->  HashAggregate  (cost=1917.90..1997.90 rows=8000 width=8) (actual time=7.203..7.565 rows=4652 loops=1)
        Group Key: orders.id, orders.amount
        Batches: 1  Memory Usage: 657kB
        ->  Append  (cost=21.08..1877.90 rows=8000 width=8) (actual time=1.570..6.286 rows=4652 loops=1)
              ->  Seq Scan on orders  (cost=21.08..878.95 rows=4000 width=8) (actual time=1.569..2.990 rows=2273 loops=1)
                    Storage Filter: ((amount)::bigint <= 100)
              ->  Seq Scan on orders orders_1  (cost=21.08..878.95 rows=4000 width=8) (actual time=1.393..2.687 rows=2379 loops=1)
                    Storage Filter: ((amount)::bigint >= 400)
Planning Time: 0.085 ms
Execution Time: 9.250 ms
Peak Memory Usage: 1252 kB

```

### s19_union_all_predicate [shape=broad/cast;spelling=union_all]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.101000 / 8.553000 ms.

```text
Sort  (cost=2396.53..2416.53 rows=8000 width=8) (actual time=7.893..8.090 rows=4652 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 411kB
  ->  Append  (cost=21.08..1877.90 rows=8000 width=8) (actual time=1.822..6.942 rows=4652 loops=1)
        ->  Seq Scan on orders  (cost=21.08..878.95 rows=4000 width=8) (actual time=1.821..3.365 rows=2273 loops=1)
              Storage Filter: ((amount)::bigint <= 100)
        ->  Seq Scan on orders orders_1  (cost=21.08..878.95 rows=4000 width=8) (actual time=1.416..2.999 rows=2379 loops=1)
              Storage Filter: ((amount)::bigint >= 400)
Planning Time: 0.101 ms
Execution Time: 8.553 ms
Peak Memory Usage: 625 kB

```

### s18_or_predicate [shape=broad/expression;index=expression_index;spelling=or]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.055000 / 4.367000 ms.

```text
Sort  (cost=1214.19..1224.72 rows=4211 width=8) (actual time=3.914..4.036 rows=4652 loops=1)
  Sort Key: id
  Sort Method: quicksort  Memory: 411kB
  ->  Seq Scan on orders  (cost=21.08..960.69 rows=4211 width=8) (actual time=0.853..3.187 rows=4652 loops=1)
        Storage Filter: (((amount + 0) <= 100) OR ((amount + 0) >= 400))
Planning Time: 0.055 ms
Execution Time: 4.367 ms
Peak Memory Usage: 548 kB

```

### s20_union_predicate [shape=broad/expression;index=expression_index;spelling=union]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.067000 / 7.032000 ms.

```text
Sort  (cost=2116.31..2127.97 rows=4664 width=8) (actual time=6.563..6.676 rows=4652 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 411kB
  ->  HashAggregate  (cost=1785.46..1832.10 rows=4664 width=8) (actual time=5.600..5.874 rows=4652 loops=1)
        Group Key: orders.id, orders.amount
        Batches: 1  Memory Usage: 465kB
        ->  Append  (cost=21.08..1762.14 rows=4664 width=8) (actual time=1.215..4.936 rows=4652 loops=1)
              ->  Seq Scan on orders  (cost=21.08..845.59 rows=2267 width=8) (actual time=1.214..2.449 rows=2273 loops=1)
                    Storage Filter: ((amount + 0) <= 100)
              ->  Seq Scan on orders orders_1  (cost=21.08..846.59 rows=2397 width=8) (actual time=1.047..2.034 rows=2379 loops=1)
                    Storage Filter: ((amount + 0) >= 400)
Planning Time: 0.067 ms
Execution Time: 7.032 ms
Peak Memory Usage: 1060 kB

```

### s19_union_all_predicate [shape=broad/expression;index=expression_index;spelling=union_all]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.108000 / 6.389000 ms.

```text
Sort  (cost=2046.35..2058.01 rows=4664 width=8) (actual time=5.943..6.060 rows=4652 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 411kB
  ->  Append  (cost=21.08..1762.14 rows=4664 width=8) (actual time=1.365..5.221 rows=4652 loops=1)
        ->  Seq Scan on orders  (cost=21.08..845.59 rows=2267 width=8) (actual time=1.364..2.494 rows=2273 loops=1)
              Storage Filter: ((amount + 0) <= 100)
        ->  Seq Scan on orders orders_1  (cost=21.08..846.59 rows=2397 width=8) (actual time=1.078..2.295 rows=2379 loops=1)
              Storage Filter: ((amount + 0) >= 400)
Planning Time: 0.108 ms
Execution Time: 6.389 ms
Peak Memory Usage: 625 kB

```

### s18_or_predicate [shape=broad/expression;spelling=or]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.061000 / 6.139000 ms.

```text
Sort  (cost=1443.14..1459.81 rows=6667 width=8) (actual time=5.396..5.648 rows=4652 loops=1)
  Sort Key: id
  Sort Method: quicksort  Memory: 411kB
  ->  Seq Scan on orders  (cost=21.08..1019.70 rows=6667 width=8) (actual time=1.222..4.385 rows=4652 loops=1)
        Storage Filter: (((amount + 0) <= 100) OR ((amount + 0) >= 400))
Planning Time: 0.061 ms
Execution Time: 6.139 ms
Peak Memory Usage: 548 kB

```

### s20_union_predicate [shape=broad/expression;spelling=union]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.078000 / 9.951000 ms.

```text
Sort  (cost=2516.53..2536.53 rows=8000 width=8) (actual time=9.173..9.342 rows=4652 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 411kB
  ->  HashAggregate  (cost=1917.90..1997.90 rows=8000 width=8) (actual time=7.893..8.255 rows=4652 loops=1)
        Group Key: orders.id, orders.amount
        Batches: 1  Memory Usage: 657kB
        ->  Append  (cost=21.08..1877.90 rows=8000 width=8) (actual time=1.673..6.885 rows=4652 loops=1)
              ->  Seq Scan on orders  (cost=21.08..878.95 rows=4000 width=8) (actual time=1.672..3.375 rows=2273 loops=1)
                    Storage Filter: ((amount + 0) <= 100)
              ->  Seq Scan on orders orders_1  (cost=21.08..878.95 rows=4000 width=8) (actual time=1.503..2.904 rows=2379 loops=1)
                    Storage Filter: ((amount + 0) >= 400)
Planning Time: 0.078 ms
Execution Time: 9.951 ms
Peak Memory Usage: 1252 kB

```

### s19_union_all_predicate [shape=broad/expression;spelling=union_all]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.104000 / 8.725000 ms.

```text
Sort  (cost=2396.53..2416.53 rows=8000 width=8) (actual time=8.049..8.248 rows=4652 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 411kB
  ->  Append  (cost=21.08..1877.90 rows=8000 width=8) (actual time=1.747..7.023 rows=4652 loops=1)
        ->  Seq Scan on orders  (cost=21.08..878.95 rows=4000 width=8) (actual time=1.744..3.346 rows=2273 loops=1)
              Storage Filter: ((amount + 0) <= 100)
        ->  Seq Scan on orders orders_1  (cost=21.08..878.95 rows=4000 width=8) (actual time=1.480..3.062 rows=2379 loops=1)
              Storage Filter: ((amount + 0) >= 400)
Planning Time: 0.104 ms
Execution Time: 8.725 ms
Peak Memory Usage: 625 kB

```

### s18_or_predicate [shape=broad/indexable;spelling=or]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.059000 / 5.668000 ms.

```text
Sort  (cost=1154.19..1164.72 rows=4211 width=8) (actual time=4.979..5.159 rows=4652 loops=1)
  Sort Key: id
  Sort Method: quicksort  Memory: 411kB
  ->  Seq Scan on orders  (cost=21.08..900.69 rows=4211 width=8) (actual time=1.294..3.934 rows=4652 loops=1)
        Storage Filter: ((amount <= 100) OR (amount >= 400))
Planning Time: 0.059 ms
Execution Time: 5.668 ms
Peak Memory Usage: 548 kB

```

### s20_union_predicate [shape=broad/indexable;spelling=union]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.090000 / 9.961000 ms.

```text
Sort  (cost=2056.31..2067.97 rows=4664 width=8) (actual time=9.109..9.340 rows=4652 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 411kB
  ->  HashAggregate  (cost=1725.46..1772.10 rows=4664 width=8) (actual time=7.858..8.182 rows=4652 loops=1)
        Group Key: orders.id, orders.amount
        Batches: 1  Memory Usage: 465kB
        ->  Append  (cost=21.08..1702.14 rows=4664 width=8) (actual time=1.727..6.904 rows=4652 loops=1)
              ->  Seq Scan on orders  (cost=21.08..815.59 rows=2267 width=8) (actual time=1.726..3.243 rows=2273 loops=1)
                    Storage Filter: (amount <= 100)
              ->  Seq Scan on orders orders_1  (cost=21.08..816.59 rows=2397 width=8) (actual time=1.614..3.037 rows=2379 loops=1)
                    Storage Filter: (amount >= 400)
Planning Time: 0.090 ms
Execution Time: 9.961 ms
Peak Memory Usage: 1060 kB

```

### s19_union_all_predicate [shape=broad/indexable;spelling=union_all]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.095000 / 8.057000 ms.

```text
Sort  (cost=1986.35..1998.01 rows=4664 width=8) (actual time=7.416..7.590 rows=4652 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 411kB
  ->  Append  (cost=21.08..1702.14 rows=4664 width=8) (actual time=1.506..6.445 rows=4652 loops=1)
        ->  Seq Scan on orders  (cost=21.08..815.59 rows=2267 width=8) (actual time=1.505..2.736 rows=2273 loops=1)
              Storage Filter: (amount <= 100)
        ->  Seq Scan on orders orders_1  (cost=21.08..816.59 rows=2397 width=8) (actual time=1.439..2.889 rows=2379 loops=1)
              Storage Filter: (amount >= 400)
Planning Time: 0.095 ms
Execution Time: 8.057 ms
Peak Memory Usage: 625 kB

```

### s18_or_predicate [shape=selective/cast;spelling=or]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.047000 / 3.589000 ms.

```text
Sort  (cost=1443.14..1459.81 rows=6667 width=8) (actual time=3.525..3.533 rows=251 loops=1)
  Sort Key: id
  Sort Method: quicksort  Memory: 36kB
  ->  Seq Scan on orders  (cost=21.08..1019.70 rows=6667 width=8) (actual time=3.438..3.468 rows=251 loops=1)
        Storage Filter: (((amount)::bigint <= 5) OR ((amount)::bigint >= 496))
Planning Time: 0.047 ms
Execution Time: 3.589 ms
Peak Memory Usage: 125 kB

```

### s20_union_predicate [shape=selective/cast;spelling=union]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.082000 / 6.530000 ms.

```text
Sort  (cost=2516.53..2536.53 rows=8000 width=8) (actual time=6.426..6.436 rows=251 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 36kB
  ->  HashAggregate  (cost=1917.90..1997.90 rows=8000 width=8) (actual time=6.321..6.370 rows=251 loops=1)
        Group Key: orders.id, orders.amount
        Batches: 1  Memory Usage: 417kB
        ->  Append  (cost=21.08..1877.90 rows=8000 width=8) (actual time=3.301..6.250 rows=251 loops=1)
              ->  Seq Scan on orders  (cost=21.08..878.95 rows=4000 width=8) (actual time=3.300..3.312 rows=114 loops=1)
                    Storage Filter: ((amount)::bigint <= 5)
              ->  Seq Scan on orders orders_1  (cost=21.08..878.95 rows=4000 width=8) (actual time=2.879..2.900 rows=137 loops=1)
                    Storage Filter: ((amount)::bigint >= 496)
Planning Time: 0.082 ms
Execution Time: 6.530 ms
Peak Memory Usage: 600 kB

```

### s19_union_all_predicate [shape=selective/cast;spelling=union_all]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.100000 / 6.459000 ms.

```text
Sort  (cost=2396.53..2416.53 rows=8000 width=8) (actual time=6.352..6.361 rows=251 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 36kB
  ->  Append  (cost=21.08..1877.90 rows=8000 width=8) (actual time=3.032..6.284 rows=251 loops=1)
        ->  Seq Scan on orders  (cost=21.08..878.95 rows=4000 width=8) (actual time=3.030..3.052 rows=114 loops=1)
              Storage Filter: ((amount)::bigint <= 5)
        ->  Seq Scan on orders orders_1  (cost=21.08..878.95 rows=4000 width=8) (actual time=3.183..3.195 rows=137 loops=1)
              Storage Filter: ((amount)::bigint >= 496)
Planning Time: 0.100 ms
Execution Time: 6.459 ms
Peak Memory Usage: 202 kB

```

### s18_or_predicate [shape=selective/expression;index=expression_index;spelling=or]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.055000 / 2.794000 ms.

```text
Sort  (cost=860.56..861.21 rows=262 width=8) (actual time=2.741..2.747 rows=251 loops=1)
  Sort Key: id
  Sort Method: quicksort  Memory: 36kB
  ->  Seq Scan on orders  (cost=21.08..850.03 rows=262 width=8) (actual time=2.672..2.689 rows=251 loops=1)
        Storage Filter: (((amount + 0) <= 5) OR ((amount + 0) >= 496))
Planning Time: 0.055 ms
Execution Time: 2.794 ms
Peak Memory Usage: 125 kB

```

### s20_union_predicate [shape=selective/expression;index=expression_index;spelling=union]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.057000 / 5.098000 ms.

```text
Sort  (cost=1596.53..1597.19 rows=263 width=8) (actual time=5.035..5.043 rows=251 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 36kB
  ->  HashAggregate  (cost=1583.33..1585.96 rows=263 width=8) (actual time=4.970..4.987 rows=251 loops=1)
        Group Key: orders.id, orders.amount
        Batches: 1  Memory Usage: 45kB
        ->  Append  (cost=21.08..1582.02 rows=263 width=8) (actual time=2.381..4.907 rows=251 loops=1)
              ->  Seq Scan on orders  (cost=21.08..788.85 rows=107 width=8) (actual time=2.380..2.389 rows=114 loops=1)
                    Storage Filter: ((amount + 0) <= 5)
              ->  Seq Scan on orders orders_1  (cost=21.08..789.22 rows=156 width=8) (actual time=2.478..2.489 rows=137 loops=1)
                    Storage Filter: ((amount + 0) >= 496)
Planning Time: 0.057 ms
Execution Time: 5.098 ms
Peak Memory Usage: 224 kB

```

### s19_union_all_predicate [shape=selective/expression;index=expression_index;spelling=union_all]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.095000 / 4.914000 ms.

```text
Sort  (cost=1592.59..1593.25 rows=263 width=8) (actual time=4.857..4.865 rows=251 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 36kB
  ->  Append  (cost=21.08..1582.02 rows=263 width=8) (actual time=2.486..4.810 rows=251 loops=1)
        ->  Seq Scan on orders  (cost=21.08..788.85 rows=107 width=8) (actual time=2.485..2.493 rows=114 loops=1)
              Storage Filter: ((amount + 0) <= 5)
        ->  Seq Scan on orders orders_1  (cost=21.08..789.22 rows=156 width=8) (actual time=2.281..2.290 rows=137 loops=1)
              Storage Filter: ((amount + 0) >= 496)
Planning Time: 0.095 ms
Execution Time: 4.914 ms
Peak Memory Usage: 200 kB

```

### s18_or_predicate [shape=selective/expression;spelling=or]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.045000 / 3.660000 ms.

```text
Sort  (cost=1443.14..1459.81 rows=6667 width=8) (actual time=3.600..3.609 rows=251 loops=1)
  Sort Key: id
  Sort Method: quicksort  Memory: 36kB
  ->  Seq Scan on orders  (cost=21.08..1019.70 rows=6667 width=8) (actual time=3.457..3.485 rows=251 loops=1)
        Storage Filter: (((amount + 0) <= 5) OR ((amount + 0) >= 496))
Planning Time: 0.045 ms
Execution Time: 3.660 ms
Peak Memory Usage: 125 kB

```

### s20_union_predicate [shape=selective/expression;spelling=union]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.099000 / 6.550000 ms.

```text
Sort  (cost=2516.53..2536.53 rows=8000 width=8) (actual time=6.445..6.454 rows=251 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 36kB
  ->  HashAggregate  (cost=1917.90..1997.90 rows=8000 width=8) (actual time=6.332..6.393 rows=251 loops=1)
        Group Key: orders.id, orders.amount
        Batches: 1  Memory Usage: 417kB
        ->  Append  (cost=21.08..1877.90 rows=8000 width=8) (actual time=3.306..6.252 rows=251 loops=1)
              ->  Seq Scan on orders  (cost=21.08..878.95 rows=4000 width=8) (actual time=3.303..3.313 rows=114 loops=1)
                    Storage Filter: ((amount + 0) <= 5)
              ->  Seq Scan on orders orders_1  (cost=21.08..878.95 rows=4000 width=8) (actual time=2.887..2.900 rows=137 loops=1)
                    Storage Filter: ((amount + 0) >= 496)
Planning Time: 0.099 ms
Execution Time: 6.550 ms
Peak Memory Usage: 600 kB

```

### s19_union_all_predicate [shape=selective/expression;spelling=union_all]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.092000 / 6.380000 ms.

```text
Sort  (cost=2396.53..2416.53 rows=8000 width=8) (actual time=6.299..6.311 rows=251 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 36kB
  ->  Append  (cost=21.08..1877.90 rows=8000 width=8) (actual time=3.198..6.216 rows=251 loops=1)
        ->  Seq Scan on orders  (cost=21.08..878.95 rows=4000 width=8) (actual time=3.197..3.207 rows=114 loops=1)
              Storage Filter: ((amount + 0) <= 5)
        ->  Seq Scan on orders orders_1  (cost=21.08..878.95 rows=4000 width=8) (actual time=2.955..2.970 rows=137 loops=1)
              Storage Filter: ((amount + 0) >= 496)
Planning Time: 0.092 ms
Execution Time: 6.380 ms
Peak Memory Usage: 202 kB

```

### s18_or_predicate [shape=selective/indexable;spelling=or]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.067000 / 3.131000 ms.

```text
Sort  (cost=800.56..801.21 rows=262 width=8) (actual time=3.074..3.081 rows=251 loops=1)
  Sort Key: id
  Sort Method: quicksort  Memory: 36kB
  ->  Seq Scan on orders  (cost=21.08..790.03 rows=262 width=8) (actual time=3.001..3.019 rows=251 loops=1)
        Storage Filter: ((amount <= 5) OR (amount >= 496))
Planning Time: 0.067 ms
Execution Time: 3.131 ms
Peak Memory Usage: 125 kB

```

### s20_union_predicate [shape=selective/indexable;spelling=union]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.075000 / 5.843000 ms.

```text
Sort  (cost=1536.53..1537.19 rows=263 width=8) (actual time=5.772..5.781 rows=251 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 36kB
  ->  HashAggregate  (cost=1523.33..1525.96 rows=263 width=8) (actual time=5.703..5.722 rows=251 loops=1)
        Group Key: orders.id, orders.amount
        Batches: 1  Memory Usage: 45kB
        ->  Append  (cost=21.08..1522.02 rows=263 width=8) (actual time=2.726..5.650 rows=251 loops=1)
              ->  Seq Scan on orders  (cost=21.08..758.85 rows=107 width=8) (actual time=2.725..2.738 rows=114 loops=1)
                    Storage Filter: (amount <= 5)
              ->  Seq Scan on orders orders_1  (cost=21.08..759.22 rows=156 width=8) (actual time=2.852..2.863 rows=137 loops=1)
                    Storage Filter: (amount >= 496)
Planning Time: 0.075 ms
Execution Time: 5.843 ms
Peak Memory Usage: 224 kB

```

### s19_union_all_predicate [shape=selective/indexable;spelling=union_all]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.090000 / 5.897000 ms.

```text
Sort  (cost=1532.59..1533.25 rows=263 width=8) (actual time=5.823..5.833 rows=251 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 36kB
  ->  Append  (cost=21.08..1522.02 rows=263 width=8) (actual time=2.897..5.757 rows=251 loops=1)
        ->  Seq Scan on orders  (cost=21.08..758.85 rows=107 width=8) (actual time=2.894..2.903 rows=114 loops=1)
              Storage Filter: (amount <= 5)
        ->  Seq Scan on orders orders_1  (cost=21.08..759.22 rows=156 width=8) (actual time=2.803..2.816 rows=137 loops=1)
              Storage Filter: (amount >= 496)
Planning Time: 0.090 ms
Execution Time: 5.897 ms
Peak Memory Usage: 210 kB

```

## Limits and gaps recorded by the runner

- no concurrent writers: client-decode JSON bytes are not wire bytes; tails are under-sampled at the protocol trial count
- one host: no independent hardware, real-network, 3-node placement or equal-total-budget comparison is established
- single-node only; the YB 3-node RF=3 per-node/equal-total arms are a coverage gap on this host

