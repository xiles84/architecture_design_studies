# SELECT access-path matrix (generated measurements)

_Measurements only. Verdicts for P10/P11 belong in a signed analysis._

**TL;DR:** within-run medians for OR / UNION ALL / UNION at three predicate shapes and OFFSET / keyset at five page depths, each cell verified against the independent oracle. A within-run difference is not a replicated result. Missing measurements are N/A.

| run | engine | topology | scale | page limit |
|---|---|---|---|---|
| select-order-access-postgres-2-20261004 | postgres | single | deep | 5 |

## Provenance

| Field | Value |
|---|---|
| repo commit | `ca2a38fbf9e53da00344578258138db9e01d3545` |
| benchmark image | `localhost/charitybench-select:select-order-access-postgres-2-20261004` |
| image digest | `8ecc060a8de698ebdf6b07ba9bd98c7663bdbb7dc6a74fa96b13a9609cd9664c` |
| run tag | `run/study-01/select-order-access-postgres-2-20261004` |
| environment | `xiles-podman-select-order-correction-20261004` |
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

## Execution order

- default matrix: order seed `42487198`, 5 interleaved rounds
- expression-index arm: order seed `42093998`, 5 interleaved rounds
- default-matrix sequence (one sample of every cell per round):
  1. shape=broad/indexable;spelling=union
  2. shape=selective/indexable;spelling=union_all
  3. shape=broad/indexable;spelling=union_all
  4. shape=broad/cast;spelling=union_all
  5. shape=selective/expression;spelling=or
  6. shape=broad/indexable;spelling=or
  7. shape=broad/cast;spelling=or
  8. shape=broad/expression;spelling=union_all
  9. shape=selective/indexable;spelling=or
  10. shape=selective/indexable;spelling=union
  11. shape=selective/expression;spelling=union_all
  12. depth=10;spelling=keyset
  13. depth=10000;spelling=offset
  14. depth=1000;spelling=offset
  15. shape=selective/cast;spelling=or
  16. shape=broad/expression;spelling=or
  17. shape=selective/expression;spelling=union
  18. depth=10000;spelling=keyset
  19. shape=selective/cast;spelling=union_all
  20. depth=100;spelling=offset
  21. depth=100;spelling=keyset
  22. shape=broad/cast;spelling=union
  23. shape=broad/expression;spelling=union
  24. depth=1000;spelling=keyset
  25. shape=selective/cast;spelling=union
  26. depth=1;spelling=keyset
  27. depth=1;spelling=offset
  28. depth=10;spelling=offset
- index-arm sequence:
  1. shape=selective/expression;index=expression_index;spelling=or
  2. shape=broad/expression;index=expression_index;spelling=union_all
  3. shape=broad/expression;index=expression_index;spelling=union
  4. shape=broad/expression;index=expression_index;spelling=or
  5. shape=selective/expression;index=expression_index;spelling=union
  6. shape=selective/expression;index=expression_index;spelling=union_all

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
| broad/cast | default | or | valid | 1.245411 | 25.8 | 5 | 4652 | `a4bab1f05f5e667d` |
| broad/cast | default | union | valid | 1.727877 | 40.4 | 5 | 4652 | `6ff77b66c9bf6ee6` |
| broad/cast | default | union_all | valid | 1.927458 | 48.5 | 5 | 4652 | `3f9f3d5c5dc15dd7` |
| broad/expression | expression_index | or | valid | 1.202983 | 23.6 | 5 | 4652 | `a4bab1f05f5e667d` |
| broad/expression | expression_index | union | valid | 1.713714 | 56.0 | 5 | 4652 | `d4fbfb3aac234d58` |
| broad/expression | expression_index | union_all | valid | 1.553154 | 4.9 | 5 | 4652 | `605ee056e65421b1` |
| broad/expression | default | or | valid | 1.472910 | 27.9 | 5 | 4652 | `a4bab1f05f5e667d` |
| broad/expression | default | union | valid | 1.833869 | 22.0 | 5 | 4652 | `6ff77b66c9bf6ee6` |
| broad/expression | default | union_all | valid | 1.856047 | 28.4 | 5 | 4652 | `3f9f3d5c5dc15dd7` |
| broad/indexable | default | or | valid | 1.279928 | 36.9 | 5 | 4652 | `a4bab1f05f5e667d` |
| broad/indexable | default | union | valid | 1.921675 | 33.7 | 5 | 4652 | `ed39b768e5b063b0` |
| broad/indexable | default | union_all | valid | 1.799185 | 33.4 | 5 | 4652 | `d51936c88987c9fd` |
| selective/cast | default | or | valid | 0.863186 | 11.1 | 5 | 251 | `a4bab1f05f5e667d` |
| selective/cast | default | union | valid | 0.954056 | 19.6 | 5 | 251 | `6ff77b66c9bf6ee6` |
| selective/cast | default | union_all | valid | 0.856929 | 19.2 | 5 | 251 | `3f9f3d5c5dc15dd7` |
| selective/expression | expression_index | or | valid | 0.271766 | 39.3 | 5 | 251 | `d66da5af4c9f3384` |
| selective/expression | expression_index | union | valid | 0.433833 | 73.0 | 5 | 251 | `a5fbbd2993065c7a` |
| selective/expression | expression_index | union_all | valid | 0.384552 | 13.8 | 5 | 251 | `eeb1f1b0c0255959` |
| selective/expression | default | or | valid | 0.970147 | 44.1 | 5 | 251 | `a4bab1f05f5e667d` |
| selective/expression | default | union | valid | 0.923608 | 11.9 | 5 | 251 | `6ff77b66c9bf6ee6` |
| selective/expression | default | union_all | valid | 0.835096 | 8.1 | 5 | 251 | `3f9f3d5c5dc15dd7` |
| selective/indexable | default | or | valid | 0.424355 | 72.3 | 5 | 251 | `8bbb087a6c122dbe` |
| selective/indexable | default | union | valid | 0.382915 | 17.6 | 5 | 251 | `567de4baf4feaaf7` |
| selective/indexable | default | union_all | valid | 0.395387 | 22.2 | 5 | 251 | `e060188a3bcf6eb3` |

### P10 pairs

Pairs compare the **default** arm only; the labelled `expression_index` arm is a package arm and is never pooled with the query-text pair.

| pair | A median ms | B median ms | B/A | same plan | within-run verdict |
|---|---|---|---|---|---|
| shape=selective/indexable or vs union_all | 0.424355 | 0.395387 | 0.9317 | false | medians inside the two cells' own spread; no resolved difference within this run |
| shape=selective/indexable or vs union | 0.424355 | 0.382915 | 0.9023 | false | medians inside the two cells' own spread; no resolved difference within this run |
| shape=selective/expression or vs union_all | 0.970147 | 0.835096 | 0.8608 | false | medians inside the two cells' own spread; no resolved difference within this run |
| shape=selective/expression or vs union | 0.970147 | 0.923608 | 0.9520 | false | medians inside the two cells' own spread; no resolved difference within this run |
| shape=selective/cast or vs union_all | 0.863186 | 0.856929 | 0.9928 | false | medians inside the two cells' own spread; no resolved difference within this run |
| shape=selective/cast or vs union | 0.863186 | 0.954056 | 1.1053 | false | medians inside the two cells' own spread; no resolved difference within this run |
| shape=broad/indexable or vs union_all | 1.279928 | 1.799185 | 1.4057 | false | medians inside the two cells' own spread; no resolved difference within this run |
| shape=broad/indexable or vs union | 1.279928 | 1.921675 | 1.5014 | false | medians inside the two cells' own spread; no resolved difference within this run |
| shape=broad/expression or vs union_all | 1.472910 | 1.856047 | 1.2601 | false | medians inside the two cells' own spread; no resolved difference within this run |
| shape=broad/expression or vs union | 1.472910 | 1.833869 | 1.2451 | false | medians inside the two cells' own spread; no resolved difference within this run |
| shape=broad/cast or vs union_all | 1.245411 | 1.927458 | 1.5476 | false | medians inside the two cells' own spread; no resolved difference within this run |
| shape=broad/cast or vs union | 1.245411 | 1.727877 | 1.3874 | false | medians inside the two cells' own spread; no resolved difference within this run |

## Pagination (P11)

| depth | spelling | status | median ms | spread % | trials | rows | plan structure |
|---|---|---|---|---|---|---|---|
| 1 | offset | valid | 0.852926 | 9.0 | 5 | 5 | `ac4509d295bf3348` |
| 1 | keyset | valid | 0.935028 | 21.8 | 5 | 5 | `ac4509d295bf3348` |
| 10 | offset | valid | 0.819740 | 10.9 | 5 | 5 | `ac4509d295bf3348` |
| 10 | keyset | valid | 0.948652 | 9.1 | 5 | 5 | `ac4509d295bf3348` |
| 100 | offset | valid | 0.855089 | 12.0 | 5 | 5 | `ac4509d295bf3348` |
| 100 | keyset | valid | 0.998355 | 42.6 | 5 | 5 | `ac4509d295bf3348` |
| 1000 | offset | valid | 1.252553 | 6.3 | 5 | 5 | `ac4509d295bf3348` |
| 1000 | keyset | valid | 1.025613 | 22.0 | 5 | 5 | `ac4509d295bf3348` |
| 10000 | offset | valid | 2.133488 | 19.3 | 5 | 5 | `ac4509d295bf3348` |
| 10000 | keyset | valid | 0.451205 | 33.1 | 5 | 5 | `dd316474b64ddc29` |

Median ratio table (keyset / OFFSET; a ratio below 1 means keyset was faster in this run):

| depth | ratio (keyset/offset) | offset median ms | keyset median ms |
|---|---|---|---|
| 1 | 1.0963 | 0.852926 | 0.935028 |
| 10 | 1.1573 | 0.819740 | 0.948652 |
| 100 | 1.1675 | 0.855089 | 0.998355 |
| 1000 | 0.8188 | 1.252553 | 1.025613 |
| 10000 | 0.2115 | 2.133488 | 0.451205 |

### P11 pairs

| pair | A median ms | B median ms | B/A | same plan | within-run verdict |
|---|---|---|---|---|---|
| depth=1 offset vs keyset | 0.852926 | 0.935028 | 1.0963 | true | medians inside the two cells' own spread; no resolved difference within this run |
| depth=10 offset vs keyset | 0.819740 | 0.948652 | 1.1573 | true | medians inside the two cells' own spread; no resolved difference within this run |
| depth=100 offset vs keyset | 0.855089 | 0.998355 | 1.1675 | true | medians inside the two cells' own spread; no resolved difference within this run |
| depth=1000 offset vs keyset | 1.252553 | 1.025613 | 0.8188 | true | medians inside the two cells' own spread; no resolved difference within this run |
| depth=10000 offset vs keyset | 2.133488 | 0.451205 | 0.2115 | false | within-run medians differ; whole-run replication required before a verdict |

## Repeated-design position control

Not measured in this run.

## Expression-index arm: write and storage cost

Index `orders_amount_expr_idx` on expression `(amount + 0)`; supported=true.

| measure | value |
|---|---|
| index size (bytes) | 114688 |
| insert batch with index (ms) | 3.666207 |
| insert batch without index (ms) | 4.070334 |
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

Rows: 5; serialized decoded-row JSON bytes: 51 (not wire bytes). Planning/execution: 0.061000 / 0.312000 ms.

```text
Limit  (cost=154.12..154.13 rows=5 width=8) (actual time=0.300..0.301 rows=5 loops=1)
  Buffers: shared hit=73
  ->  Sort  (cost=154.12..158.90 rows=1912 width=8) (actual time=0.300..0.300 rows=5 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 25kB
        Buffers: shared hit=73
        ->  Bitmap Heap Scan on orders  (cost=22.66..122.36 rows=1912 width=8) (actual time=0.036..0.200 rows=2000 loops=1)
              Recheck Cond: ((amount < 67) OR (amount = 67) OR (amount IS NULL))
              Filter: ((amount < 67) OR ((amount = 67) AND (id > 4310)) OR (amount IS NULL))
              Rows Removed by Filter: 6
              Heap Blocks: exact=65
              Buffers: shared hit=73
              ->  BitmapOr  (cost=22.66..22.66 rows=1983 width=0) (actual time=0.027..0.028 rows=0 loops=1)
                    Buffers: shared hit=8
                    ->  Bitmap Index Scan on orders_amount_idx  (cost=0.00..14.38 rows=1440 width=0) (actual time=0.021..0.021 rows=1460 loops=1)
                          Index Cond: (amount < 67)
                          Buffers: shared hit=4
                    ->  Bitmap Index Scan on orders_amount_idx  (cost=0.00..1.54 rows=21 width=0) (actual time=0.001..0.001 rows=24 loops=1)
                          Index Cond: (amount = 67)
                          Buffers: shared hit=2
                    ->  Bitmap Index Scan on orders_amount_idx  (cost=0.00..5.30 rows=522 width=0) (actual time=0.005..0.005 rows=522 loops=1)
                          Index Cond: (amount IS NULL)
                          Buffers: shared hit=2
Planning Time: 0.061 ms
Execution Time: 0.312 ms

```

### s16_page_offset [depth=10000;spelling=offset]

Rows: 5; serialized decoded-row JSON bytes: 51 (not wire bytes). Planning/execution: 0.029000 / 2.613000 ms.

```text
Limit  (cost=1023.04..1023.06 rows=5 width=8) (actual time=2.606..2.607 rows=5 loops=1)
  Buffers: shared hit=65
  ->  Sort  (cost=998.04..1028.04 rows=12000 width=8) (actual time=2.181..2.417 rows=10005 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: quicksort  Memory: 853kB
        Buffers: shared hit=65
        ->  Seq Scan on orders  (cost=0.00..185.00 rows=12000 width=8) (actual time=0.005..0.467 rows=12000 loops=1)
              Buffers: shared hit=65
Planning Time: 0.029 ms
Execution Time: 2.613 ms

```

### s17_page_keyset [depth=1000;spelling=keyset]

Rows: 5; serialized decoded-row JSON bytes: 59 (not wire bytes). Planning/execution: 0.040000 / 1.170000 ms.

```text
Limit  (cost=450.03..450.04 rows=5 width=8) (actual time=1.164..1.164 rows=5 loops=1)
  Buffers: shared hit=65
  ->  Sort  (cost=450.03..476.38 rows=10538 width=8) (actual time=1.163..1.163 rows=5 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 25kB
        Buffers: shared hit=65
        ->  Seq Scan on orders  (cost=0.00..275.00 rows=10538 width=8) (actual time=0.005..0.574 rows=11000 loops=1)
              Filter: ((amount < 457) OR ((amount = 457) AND (id > 11254)) OR (amount IS NULL))
              Rows Removed by Filter: 1000
              Buffers: shared hit=65
Planning Time: 0.040 ms
Execution Time: 1.170 ms

```

### s16_page_offset [depth=1000;spelling=offset]

Rows: 5; serialized decoded-row JSON bytes: 59 (not wire bytes). Planning/execution: 0.028000 / 1.571000 ms.

```text
Limit  (cost=845.88..845.89 rows=5 width=8) (actual time=1.564..1.565 rows=5 loops=1)
  Buffers: shared hit=65
  ->  Sort  (cost=843.38..873.38 rows=12000 width=8) (actual time=1.520..1.544 rows=1005 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 115kB
        Buffers: shared hit=65
        ->  Seq Scan on orders  (cost=0.00..185.00 rows=12000 width=8) (actual time=0.005..0.515 rows=12000 loops=1)
              Buffers: shared hit=65
Planning Time: 0.028 ms
Execution Time: 1.571 ms

```

### s17_page_keyset [depth=100;spelling=keyset]

Rows: 5; serialized decoded-row JSON bytes: 56 (not wire bytes). Planning/execution: 0.054000 / 1.278000 ms.

```text
Limit  (cost=463.65..463.66 rows=5 width=8) (actual time=1.269..1.270 rows=5 loops=1)
  Buffers: shared hit=65
  ->  Sort  (cost=463.65..492.05 rows=11358 width=8) (actual time=1.268..1.269 rows=5 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 25kB
        Buffers: shared hit=65
        ->  Seq Scan on orders  (cost=0.00..275.00 rows=11358 width=8) (actual time=0.005..0.583 rows=11900 loops=1)
              Filter: ((amount < 497) OR ((amount = 497) AND (id > 6723)) OR (amount IS NULL))
              Rows Removed by Filter: 100
              Buffers: shared hit=65
Planning:
  Buffers: shared hit=3
Planning Time: 0.054 ms
Execution Time: 1.278 ms

```

### s16_page_offset [depth=100;spelling=offset]

Rows: 5; serialized decoded-row JSON bytes: 56 (not wire bytes). Planning/execution: 0.032000 / 1.067000 ms.

```text
Limit  (cost=648.10..648.12 rows=5 width=8) (actual time=1.061..1.062 rows=5 loops=1)
  Buffers: shared hit=65
  ->  Sort  (cost=647.85..677.85 rows=12000 width=8) (actual time=1.057..1.059 rows=105 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 31kB
        Buffers: shared hit=65
        ->  Seq Scan on orders  (cost=0.00..185.00 rows=12000 width=8) (actual time=0.006..0.462 rows=12000 loops=1)
              Buffers: shared hit=65
Planning Time: 0.032 ms
Execution Time: 1.067 ms

```

### s17_page_keyset [depth=10;spelling=keyset]

Rows: 5; serialized decoded-row JSON bytes: 56 (not wire bytes). Planning/execution: 0.036000 / 1.115000 ms.

```text
Limit  (cost=465.21..465.23 rows=5 width=8) (actual time=1.110..1.111 rows=5 loops=1)
  Buffers: shared hit=65
  ->  Sort  (cost=465.21..493.84 rows=11452 width=8) (actual time=1.110..1.110 rows=5 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 25kB
        Buffers: shared hit=65
        ->  Seq Scan on orders  (cost=0.00..275.00 rows=11452 width=8) (actual time=0.003..0.540 rows=11990 loops=1)
              Filter: ((amount < 500) OR ((amount = 500) AND (id > 6319)) OR (amount IS NULL))
              Rows Removed by Filter: 10
              Buffers: shared hit=65
Planning:
  Buffers: shared hit=3
Planning Time: 0.036 ms
Execution Time: 1.115 ms

```

### s16_page_offset [depth=10;spelling=offset]

Rows: 5; serialized decoded-row JSON bytes: 56 (not wire bytes). Planning/execution: 0.023000 / 1.119000 ms.

```text
Limit  (cost=479.44..479.45 rows=5 width=8) (actual time=1.113..1.114 rows=5 loops=1)
  Buffers: shared hit=65
  ->  Sort  (cost=479.41..509.41 rows=12000 width=8) (actual time=1.112..1.113 rows=15 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 25kB
        Buffers: shared hit=65
        ->  Seq Scan on orders  (cost=0.00..185.00 rows=12000 width=8) (actual time=0.004..0.492 rows=12000 loops=1)
              Buffers: shared hit=65
Planning Time: 0.023 ms
Execution Time: 1.119 ms

```

### s17_page_keyset [depth=1;spelling=keyset]

Rows: 5; serialized decoded-row JSON bytes: 56 (not wire bytes). Planning/execution: 0.048000 / 1.197000 ms.

```text
Limit  (cost=465.21..465.23 rows=5 width=8) (actual time=1.191..1.191 rows=5 loops=1)
  Buffers: shared hit=65
  ->  Sort  (cost=465.21..493.84 rows=11452 width=8) (actual time=1.190..1.191 rows=5 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 25kB
        Buffers: shared hit=65
        ->  Seq Scan on orders  (cost=0.00..275.00 rows=11452 width=8) (actual time=0.004..0.552 rows=11999 loops=1)
              Filter: ((amount < 500) OR ((amount = 500) AND (id > 1554)) OR (amount IS NULL))
              Rows Removed by Filter: 1
              Buffers: shared hit=65
Planning:
  Buffers: shared hit=3
Planning Time: 0.048 ms
Execution Time: 1.197 ms

```

### s16_page_offset [depth=1;spelling=offset]

Rows: 5; serialized decoded-row JSON bytes: 56 (not wire bytes). Planning/execution: 0.020000 / 1.040000 ms.

```text
Limit  (cost=400.10..400.11 rows=5 width=8) (actual time=1.035..1.036 rows=5 loops=1)
  Buffers: shared hit=65
  ->  Sort  (cost=400.10..430.10 rows=12000 width=8) (actual time=1.035..1.035 rows=6 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 25kB
        Buffers: shared hit=65
        ->  Seq Scan on orders  (cost=0.00..185.00 rows=12000 width=8) (actual time=0.004..0.457 rows=12000 loops=1)
              Buffers: shared hit=65
Planning Time: 0.020 ms
Execution Time: 1.040 ms

```

### s18_or_predicate [shape=broad/cast;spelling=or]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.027000 / 0.872000 ms.

```text
Index Scan using orders_pkey on orders  (cost=0.29..403.89 rows=6667 width=8) (actual time=0.008..0.769 rows=4652 loops=1)
  Filter: (((amount)::bigint <= 100) OR ((amount)::bigint >= 400))
  Rows Removed by Filter: 7348
  Buffers: shared hit=99
Planning Time: 0.027 ms
Execution Time: 0.872 ms

```

### s20_union_predicate [shape=broad/cast;spelling=union]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.055000 / 1.664000 ms.

```text
Unique  (cost=968.64..1108.64 rows=8000 width=8) (actual time=0.922..1.562 rows=4652 loops=1)
  Buffers: shared hit=130
  ->  Merge Append  (cost=968.64..1068.64 rows=8000 width=8) (actual time=0.922..1.264 rows=4652 loops=1)
        Sort Key: orders.id, orders.amount
        Buffers: shared hit=130
        ->  Sort  (cost=484.32..494.32 rows=4000 width=8) (actual time=0.466..0.516 rows=2273 loops=1)
              Sort Key: orders.id, orders.amount
              Sort Method: quicksort  Memory: 150kB
              Buffers: shared hit=65
              ->  Seq Scan on orders  (cost=0.00..245.00 rows=4000 width=8) (actual time=0.007..0.351 rows=2273 loops=1)
                    Filter: ((amount)::bigint <= 100)
                    Rows Removed by Filter: 9727
                    Buffers: shared hit=65
        ->  Sort  (cost=484.32..494.32 rows=4000 width=8) (actual time=0.455..0.508 rows=2379 loops=1)
              Sort Key: orders_1.id, orders_1.amount
              Sort Method: quicksort  Memory: 152kB
              Buffers: shared hit=65
              ->  Seq Scan on orders orders_1  (cost=0.00..245.00 rows=4000 width=8) (actual time=0.002..0.339 rows=2379 loops=1)
                    Filter: ((amount)::bigint >= 400)
                    Rows Removed by Filter: 9621
                    Buffers: shared hit=65
Planning Time: 0.055 ms
Execution Time: 1.664 ms

```

### s19_union_all_predicate [shape=broad/cast;spelling=union_all]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.043000 / 1.441000 ms.

```text
Sort  (cost=1048.63..1068.63 rows=8000 width=8) (actual time=1.244..1.343 rows=4652 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 302kB
  Buffers: shared hit=130
  ->  Append  (cost=0.00..530.00 rows=8000 width=8) (actual time=0.008..0.887 rows=4652 loops=1)
        Buffers: shared hit=130
        ->  Seq Scan on orders  (cost=0.00..245.00 rows=4000 width=8) (actual time=0.008..0.381 rows=2273 loops=1)
              Filter: ((amount)::bigint <= 100)
              Rows Removed by Filter: 9727
              Buffers: shared hit=65
        ->  Seq Scan on orders orders_1  (cost=0.00..245.00 rows=4000 width=8) (actual time=0.002..0.328 rows=2379 loops=1)
              Filter: ((amount)::bigint >= 400)
              Rows Removed by Filter: 9621
              Buffers: shared hit=65
Planning Time: 0.043 ms
Execution Time: 1.441 ms

```

### s18_or_predicate [shape=broad/expression;index=expression_index;spelling=or]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.034000 / 0.853000 ms.

```text
Index Scan using orders_pkey on orders  (cost=0.29..403.89 rows=4211 width=8) (actual time=0.008..0.754 rows=4652 loops=1)
  Filter: (((amount + 0) <= 100) OR ((amount + 0) >= 400))
  Rows Removed by Filter: 7348
  Buffers: shared hit=99
Planning Time: 0.034 ms
Execution Time: 0.853 ms

```

### s20_union_predicate [shape=broad/expression;index=expression_index;spelling=union]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.062000 / 2.012000 ms.

```text
Unique  (cost=492.69..574.31 rows=4664 width=8) (actual time=1.018..1.907 rows=4652 loops=1)
  Buffers: shared hit=3934
  ->  Merge Append  (cost=492.69..550.99 rows=4664 width=8) (actual time=1.017..1.636 rows=4652 loops=1)
        Sort Key: orders.id, orders.amount
        Buffers: shared hit=3934
        ->  Sort  (cost=241.10..246.77 rows=2267 width=8) (actual time=0.526..0.579 rows=2273 loops=1)
              Sort Key: orders.id, orders.amount
              Sort Method: quicksort  Memory: 150kB
              Buffers: shared hit=1921
              ->  Index Scan using orders_amount_expr_idx on orders  (cost=0.29..114.75 rows=2267 width=8) (actual time=0.008..0.307 rows=2273 loops=1)
                    Index Cond: ((amount + 0) <= 100)
                    Buffers: shared hit=1921
        ->  Sort  (cost=251.58..257.57 rows=2397 width=8) (actual time=0.490..0.545 rows=2379 loops=1)
              Sort Key: orders_1.id, orders_1.amount
              Sort Method: quicksort  Memory: 152kB
              Buffers: shared hit=2013
              ->  Index Scan using orders_amount_expr_idx on orders orders_1  (cost=0.29..117.03 rows=2397 width=8) (actual time=0.003..0.263 rows=2379 loops=1)
                    Index Cond: ((amount + 0) >= 400)
                    Buffers: shared hit=2013
Planning Time: 0.062 ms
Execution Time: 2.012 ms

```

### s19_union_all_predicate [shape=broad/expression;index=expression_index;spelling=union_all]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.050000 / 1.369000 ms.

```text
Sort  (cost=539.31..550.97 rows=4664 width=8) (actual time=1.154..1.263 rows=4652 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 302kB
  Buffers: shared hit=3934
  ->  Append  (cost=0.29..255.10 rows=4664 width=8) (actual time=0.008..0.711 rows=4652 loops=1)
        Buffers: shared hit=3934
        ->  Index Scan using orders_amount_expr_idx on orders  (cost=0.29..114.75 rows=2267 width=8) (actual time=0.007..0.288 rows=2273 loops=1)
              Index Cond: ((amount + 0) <= 100)
              Buffers: shared hit=1921
        ->  Index Scan using orders_amount_expr_idx on orders orders_1  (cost=0.29..117.03 rows=2397 width=8) (actual time=0.004..0.251 rows=2379 loops=1)
              Index Cond: ((amount + 0) >= 400)
              Buffers: shared hit=2013
Planning Time: 0.050 ms
Execution Time: 1.369 ms

```

### s18_or_predicate [shape=broad/expression;spelling=or]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.044000 / 0.843000 ms.

```text
Index Scan using orders_pkey on orders  (cost=0.29..403.89 rows=6667 width=8) (actual time=0.011..0.741 rows=4652 loops=1)
  Filter: (((amount + 0) <= 100) OR ((amount + 0) >= 400))
  Rows Removed by Filter: 7348
  Buffers: shared hit=99
Planning Time: 0.044 ms
Execution Time: 0.843 ms

```

### s20_union_predicate [shape=broad/expression;spelling=union]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.054000 / 1.743000 ms.

```text
Unique  (cost=968.64..1108.64 rows=8000 width=8) (actual time=0.955..1.619 rows=4652 loops=1)
  Buffers: shared hit=130
  ->  Merge Append  (cost=968.64..1068.64 rows=8000 width=8) (actual time=0.955..1.342 rows=4652 loops=1)
        Sort Key: orders.id, orders.amount
        Buffers: shared hit=130
        ->  Sort  (cost=484.32..494.32 rows=4000 width=8) (actual time=0.490..0.547 rows=2273 loops=1)
              Sort Key: orders.id, orders.amount
              Sort Method: quicksort  Memory: 150kB
              Buffers: shared hit=65
              ->  Seq Scan on orders  (cost=0.00..245.00 rows=4000 width=8) (actual time=0.007..0.366 rows=2273 loops=1)
                    Filter: ((amount + 0) <= 100)
                    Rows Removed by Filter: 9727
                    Buffers: shared hit=65
        ->  Sort  (cost=484.32..494.32 rows=4000 width=8) (actual time=0.464..0.518 rows=2379 loops=1)
              Sort Key: orders_1.id, orders_1.amount
              Sort Method: quicksort  Memory: 152kB
              Buffers: shared hit=65
              ->  Seq Scan on orders orders_1  (cost=0.00..245.00 rows=4000 width=8) (actual time=0.002..0.342 rows=2379 loops=1)
                    Filter: ((amount + 0) >= 400)
                    Rows Removed by Filter: 9621
                    Buffers: shared hit=65
Planning Time: 0.054 ms
Execution Time: 1.743 ms

```

### s19_union_all_predicate [shape=broad/expression;spelling=union_all]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.053000 / 1.768000 ms.

```text
Sort  (cost=1048.63..1068.63 rows=8000 width=8) (actual time=1.577..1.672 rows=4652 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 302kB
  Buffers: shared hit=130
  ->  Append  (cost=0.00..530.00 rows=8000 width=8) (actual time=0.009..1.206 rows=4652 loops=1)
        Buffers: shared hit=130
        ->  Seq Scan on orders  (cost=0.00..245.00 rows=4000 width=8) (actual time=0.009..0.674 rows=2273 loops=1)
              Filter: ((amount + 0) <= 100)
              Rows Removed by Filter: 9727
              Buffers: shared hit=65
        ->  Seq Scan on orders orders_1  (cost=0.00..245.00 rows=4000 width=8) (actual time=0.007..0.358 rows=2379 loops=1)
              Filter: ((amount + 0) >= 400)
              Rows Removed by Filter: 9621
              Buffers: shared hit=65
Planning Time: 0.053 ms
Execution Time: 1.768 ms

```

### s18_or_predicate [shape=broad/indexable;spelling=or]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.040000 / 0.826000 ms.

```text
Index Scan using orders_pkey on orders  (cost=0.29..343.89 rows=4211 width=8) (actual time=0.009..0.722 rows=4652 loops=1)
  Filter: ((amount <= 100) OR (amount >= 400))
  Rows Removed by Filter: 7348
  Buffers: shared hit=99
Planning Time: 0.040 ms
Execution Time: 0.826 ms

```

### s20_union_predicate [shape=broad/indexable;spelling=union]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.068000 / 1.683000 ms.

```text
Unique  (cost=495.99..577.61 rows=4664 width=8) (actual time=0.966..1.580 rows=4652 loops=1)
  Buffers: shared hit=3936
  ->  Merge Append  (cost=495.99..554.29 rows=4664 width=8) (actual time=0.966..1.323 rows=4652 loops=1)
        Sort Key: orders.id, orders.amount
        Buffers: shared hit=3936
        ->  Sort  (cost=242.20..247.87 rows=2267 width=8) (actual time=0.516..0.566 rows=2273 loops=1)
              Sort Key: orders.id, orders.amount
              Sort Method: quicksort  Memory: 150kB
              Buffers: shared hit=1922
              ->  Index Scan using orders_amount_idx on orders  (cost=0.29..115.85 rows=2267 width=8) (actual time=0.010..0.303 rows=2273 loops=1)
                    Index Cond: (amount <= 100)
                    Buffers: shared hit=1922
        ->  Sort  (cost=253.78..259.77 rows=2397 width=8) (actual time=0.448..0.500 rows=2379 loops=1)
              Sort Key: orders_1.id, orders_1.amount
              Sort Method: quicksort  Memory: 152kB
              Buffers: shared hit=2014
              ->  Index Scan using orders_amount_idx on orders orders_1  (cost=0.29..119.23 rows=2397 width=8) (actual time=0.004..0.243 rows=2379 loops=1)
                    Index Cond: (amount >= 400)
                    Buffers: shared hit=2014
Planning Time: 0.068 ms
Execution Time: 1.683 ms

```

### s19_union_all_predicate [shape=broad/indexable;spelling=union_all]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.056000 / 1.448000 ms.

```text
Sort  (cost=542.61..554.27 rows=4664 width=8) (actual time=1.235..1.341 rows=4652 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 302kB
  Buffers: shared hit=3936
  ->  Append  (cost=0.29..258.40 rows=4664 width=8) (actual time=0.009..0.784 rows=4652 loops=1)
        Buffers: shared hit=3936
        ->  Index Scan using orders_amount_idx on orders  (cost=0.29..115.85 rows=2267 width=8) (actual time=0.009..0.323 rows=2273 loops=1)
              Index Cond: (amount <= 100)
              Buffers: shared hit=1922
        ->  Index Scan using orders_amount_idx on orders orders_1  (cost=0.29..119.23 rows=2397 width=8) (actual time=0.004..0.258 rows=2379 loops=1)
              Index Cond: (amount >= 400)
              Buffers: shared hit=2014
Planning Time: 0.056 ms
Execution Time: 1.448 ms

```

### s18_or_predicate [shape=selective/cast;spelling=or]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.023000 / 0.624000 ms.

```text
Index Scan using orders_pkey on orders  (cost=0.29..403.89 rows=6667 width=8) (actual time=0.009..0.614 rows=251 loops=1)
  Filter: (((amount)::bigint <= 5) OR ((amount)::bigint >= 496))
  Rows Removed by Filter: 11749
  Buffers: shared hit=99
Planning Time: 0.023 ms
Execution Time: 0.624 ms

```

### s20_union_predicate [shape=selective/cast;spelling=union]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.043000 / 0.647000 ms.

```text
Unique  (cost=968.64..1108.64 rows=8000 width=8) (actual time=0.599..0.634 rows=251 loops=1)
  Buffers: shared hit=130
  ->  Merge Append  (cost=968.64..1068.64 rows=8000 width=8) (actual time=0.599..0.618 rows=251 loops=1)
        Sort Key: orders.id, orders.amount
        Buffers: shared hit=130
        ->  Sort  (cost=484.32..494.32 rows=4000 width=8) (actual time=0.330..0.333 rows=114 loops=1)
              Sort Key: orders.id, orders.amount
              Sort Method: quicksort  Memory: 27kB
              Buffers: shared hit=65
              ->  Seq Scan on orders  (cost=0.00..245.00 rows=4000 width=8) (actual time=0.006..0.322 rows=114 loops=1)
                    Filter: ((amount)::bigint <= 5)
                    Rows Removed by Filter: 11886
                    Buffers: shared hit=65
        ->  Sort  (cost=484.32..494.32 rows=4000 width=8) (actual time=0.268..0.271 rows=137 loops=1)
              Sort Key: orders_1.id, orders_1.amount
              Sort Method: quicksort  Memory: 28kB
              Buffers: shared hit=65
              ->  Seq Scan on orders orders_1  (cost=0.00..245.00 rows=4000 width=8) (actual time=0.003..0.260 rows=137 loops=1)
                    Filter: ((amount)::bigint >= 496)
                    Rows Removed by Filter: 11863
                    Buffers: shared hit=65
Planning Time: 0.043 ms
Execution Time: 0.647 ms

```

### s19_union_all_predicate [shape=selective/cast;spelling=union_all]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.033000 / 0.609000 ms.

```text
Sort  (cost=1048.63..1068.63 rows=8000 width=8) (actual time=0.591..0.597 rows=251 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 30kB
  Buffers: shared hit=130
  ->  Append  (cost=0.00..530.00 rows=8000 width=8) (actual time=0.006..0.569 rows=251 loops=1)
        Buffers: shared hit=130
        ->  Seq Scan on orders  (cost=0.00..245.00 rows=4000 width=8) (actual time=0.005..0.282 rows=114 loops=1)
              Filter: ((amount)::bigint <= 5)
              Rows Removed by Filter: 11886
              Buffers: shared hit=65
        ->  Seq Scan on orders orders_1  (cost=0.00..245.00 rows=4000 width=8) (actual time=0.003..0.276 rows=137 loops=1)
              Filter: ((amount)::bigint >= 496)
              Rows Removed by Filter: 11863
              Buffers: shared hit=65
Planning Time: 0.033 ms
Execution Time: 0.609 ms

```

### s18_or_predicate [shape=selective/expression;index=expression_index;spelling=or]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.039000 / 0.090000 ms.

```text
Sort  (cost=86.11..86.78 rows=268 width=8) (actual time=0.064..0.079 rows=251 loops=1)
  Sort Key: id
  Sort Method: quicksort  Memory: 30kB
  Buffers: shared hit=67
  ->  Bitmap Heap Scan on orders  (cost=4.92..75.30 rows=268 width=8) (actual time=0.011..0.050 rows=251 loops=1)
        Recheck Cond: (((amount + 0) <= 5) OR ((amount + 0) >= 496))
        Heap Blocks: exact=63
        Buffers: shared hit=67
        ->  BitmapOr  (cost=4.92..4.92 rows=269 width=0) (actual time=0.005..0.006 rows=0 loops=1)
              Buffers: shared hit=4
              ->  Bitmap Index Scan on orders_amount_expr_idx  (cost=0.00..2.19 rows=107 width=0) (actual time=0.003..0.003 rows=114 loops=1)
                    Index Cond: ((amount + 0) <= 5)
                    Buffers: shared hit=2
              ->  Bitmap Index Scan on orders_amount_expr_idx  (cost=0.00..2.60 rows=162 width=0) (actual time=0.002..0.002 rows=137 loops=1)
                    Index Cond: ((amount + 0) >= 496)
                    Buffers: shared hit=2
Planning:
  Buffers: shared hit=6
Planning Time: 0.039 ms
Execution Time: 0.090 ms

```

### s20_union_predicate [shape=selective/expression;index=expression_index;spelling=union]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.059000 / 0.124000 ms.

```text
Unique  (cost=142.73..147.44 rows=269 width=8) (actual time=0.073..0.110 rows=251 loops=1)
  Buffers: shared hit=114
  ->  Merge Append  (cost=142.73..146.09 rows=269 width=8) (actual time=0.073..0.093 rows=251 loops=1)
        Sort Key: orders.id, orders.amount
        Buffers: shared hit=114
        ->  Sort  (cost=66.70..66.97 rows=107 width=8) (actual time=0.037..0.040 rows=114 loops=1)
              Sort Key: orders.id, orders.amount
              Sort Method: quicksort  Memory: 27kB
              Buffers: shared hit=56
              ->  Bitmap Heap Scan on orders  (cost=2.21..63.10 rows=107 width=8) (actual time=0.009..0.030 rows=114 loops=1)
                    Recheck Cond: ((amount + 0) <= 5)
                    Heap Blocks: exact=54
                    Buffers: shared hit=56
                    ->  Bitmap Index Scan on orders_amount_expr_idx  (cost=0.00..2.19 rows=107 width=0) (actual time=0.003..0.003 rows=114 loops=1)
                          Index Cond: ((amount + 0) <= 5)
                          Buffers: shared hit=2
        ->  Sort  (cost=76.02..76.42 rows=162 width=8) (actual time=0.035..0.038 rows=137 loops=1)
              Sort Key: orders_1.id, orders_1.amount
              Sort Method: quicksort  Memory: 28kB
              Buffers: shared hit=58
              ->  Bitmap Heap Scan on orders orders_1  (cost=2.64..70.07 rows=162 width=8) (actual time=0.007..0.027 rows=137 loops=1)
                    Recheck Cond: ((amount + 0) >= 496)
                    Heap Blocks: exact=56
                    Buffers: shared hit=58
                    ->  Bitmap Index Scan on orders_amount_expr_idx  (cost=0.00..2.60 rows=162 width=0) (actual time=0.003..0.003 rows=137 loops=1)
                          Index Cond: ((amount + 0) >= 496)
                          Buffers: shared hit=2
Planning:
  Buffers: shared hit=6
Planning Time: 0.059 ms
Execution Time: 0.124 ms

```

### s19_union_all_predicate [shape=selective/expression;index=expression_index;spelling=union_all]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.043000 / 0.112000 ms.

```text
Sort  (cost=145.37..146.04 rows=269 width=8) (actual time=0.093..0.100 rows=251 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 30kB
  Buffers: shared hit=114
  ->  Append  (cost=2.21..134.51 rows=269 width=8) (actual time=0.008..0.071 rows=251 loops=1)
        Buffers: shared hit=114
        ->  Bitmap Heap Scan on orders  (cost=2.21..63.10 rows=107 width=8) (actual time=0.008..0.030 rows=114 loops=1)
              Recheck Cond: ((amount + 0) <= 5)
              Heap Blocks: exact=54
              Buffers: shared hit=56
              ->  Bitmap Index Scan on orders_amount_expr_idx  (cost=0.00..2.19 rows=107 width=0) (actual time=0.003..0.003 rows=114 loops=1)
                    Index Cond: ((amount + 0) <= 5)
                    Buffers: shared hit=2
        ->  Bitmap Heap Scan on orders orders_1  (cost=2.64..70.07 rows=162 width=8) (actual time=0.007..0.029 rows=137 loops=1)
              Recheck Cond: ((amount + 0) >= 496)
              Heap Blocks: exact=56
              Buffers: shared hit=58
              ->  Bitmap Index Scan on orders_amount_expr_idx  (cost=0.00..2.60 rows=162 width=0) (actual time=0.003..0.003 rows=137 loops=1)
                    Index Cond: ((amount + 0) >= 496)
                    Buffers: shared hit=2
Planning:
  Buffers: shared hit=6
Planning Time: 0.043 ms
Execution Time: 0.112 ms

```

### s18_or_predicate [shape=selective/expression;spelling=or]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.023000 / 0.666000 ms.

```text
Index Scan using orders_pkey on orders  (cost=0.29..403.89 rows=6667 width=8) (actual time=0.008..0.656 rows=251 loops=1)
  Filter: (((amount + 0) <= 5) OR ((amount + 0) >= 496))
  Rows Removed by Filter: 11749
  Buffers: shared hit=99
Planning Time: 0.023 ms
Execution Time: 0.666 ms

```

### s20_union_predicate [shape=selective/expression;spelling=union]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.041000 / 0.600000 ms.

```text
Unique  (cost=968.64..1108.64 rows=8000 width=8) (actual time=0.552..0.587 rows=251 loops=1)
  Buffers: shared hit=130
  ->  Merge Append  (cost=968.64..1068.64 rows=8000 width=8) (actual time=0.552..0.571 rows=251 loops=1)
        Sort Key: orders.id, orders.amount
        Buffers: shared hit=130
        ->  Sort  (cost=484.32..494.32 rows=4000 width=8) (actual time=0.295..0.297 rows=114 loops=1)
              Sort Key: orders.id, orders.amount
              Sort Method: quicksort  Memory: 27kB
              Buffers: shared hit=65
              ->  Seq Scan on orders  (cost=0.00..245.00 rows=4000 width=8) (actual time=0.005..0.287 rows=114 loops=1)
                    Filter: ((amount + 0) <= 5)
                    Rows Removed by Filter: 11886
                    Buffers: shared hit=65
        ->  Sort  (cost=484.32..494.32 rows=4000 width=8) (actual time=0.257..0.260 rows=137 loops=1)
              Sort Key: orders_1.id, orders_1.amount
              Sort Method: quicksort  Memory: 28kB
              Buffers: shared hit=65
              ->  Seq Scan on orders orders_1  (cost=0.00..245.00 rows=4000 width=8) (actual time=0.004..0.249 rows=137 loops=1)
                    Filter: ((amount + 0) >= 496)
                    Rows Removed by Filter: 11863
                    Buffers: shared hit=65
Planning Time: 0.041 ms
Execution Time: 0.600 ms

```

### s19_union_all_predicate [shape=selective/expression;spelling=union_all]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.032000 / 0.627000 ms.

```text
Sort  (cost=1048.63..1068.63 rows=8000 width=8) (actual time=0.610..0.616 rows=251 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 30kB
  Buffers: shared hit=130
  ->  Append  (cost=0.00..530.00 rows=8000 width=8) (actual time=0.005..0.589 rows=251 loops=1)
        Buffers: shared hit=130
        ->  Seq Scan on orders  (cost=0.00..245.00 rows=4000 width=8) (actual time=0.005..0.319 rows=114 loops=1)
              Filter: ((amount + 0) <= 5)
              Rows Removed by Filter: 11886
              Buffers: shared hit=65
        ->  Seq Scan on orders orders_1  (cost=0.00..245.00 rows=4000 width=8) (actual time=0.004..0.259 rows=137 loops=1)
              Filter: ((amount + 0) >= 496)
              Rows Removed by Filter: 11863
              Buffers: shared hit=65
Planning Time: 0.032 ms
Execution Time: 0.627 ms

```

### s18_or_predicate [shape=selective/indexable;spelling=or]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.040000 / 0.081000 ms.

```text
Sort  (cost=84.77..85.44 rows=268 width=8) (actual time=0.062..0.069 rows=251 loops=1)
  Sort Key: id
  Sort Method: quicksort  Memory: 30kB
  Buffers: shared hit=67
  ->  Bitmap Heap Scan on orders  (cost=4.92..73.96 rows=268 width=8) (actual time=0.013..0.048 rows=251 loops=1)
        Recheck Cond: ((amount <= 5) OR (amount >= 496))
        Heap Blocks: exact=63
        Buffers: shared hit=67
        ->  BitmapOr  (cost=4.92..4.92 rows=269 width=0) (actual time=0.007..0.007 rows=0 loops=1)
              Buffers: shared hit=4
              ->  Bitmap Index Scan on orders_amount_idx  (cost=0.00..2.19 rows=107 width=0) (actual time=0.004..0.004 rows=114 loops=1)
                    Index Cond: (amount <= 5)
                    Buffers: shared hit=2
              ->  Bitmap Index Scan on orders_amount_idx  (cost=0.00..2.60 rows=162 width=0) (actual time=0.003..0.003 rows=137 loops=1)
                    Index Cond: (amount >= 496)
                    Buffers: shared hit=2
Planning:
  Buffers: shared hit=6
Planning Time: 0.040 ms
Execution Time: 0.081 ms

```

### s20_union_predicate [shape=selective/indexable;spelling=union]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.058000 / 0.158000 ms.

```text
Unique  (cost=142.06..146.77 rows=269 width=8) (actual time=0.100..0.142 rows=251 loops=1)
  Buffers: shared hit=114
  ->  Merge Append  (cost=142.06..145.42 rows=269 width=8) (actual time=0.100..0.123 rows=251 loops=1)
        Sort Key: orders.id, orders.amount
        Buffers: shared hit=114
        ->  Sort  (cost=66.44..66.70 rows=107 width=8) (actual time=0.058..0.061 rows=114 loops=1)
              Sort Key: orders.id, orders.amount
              Sort Method: quicksort  Memory: 27kB
              Buffers: shared hit=56
              ->  Bitmap Heap Scan on orders  (cost=2.21..62.83 rows=107 width=8) (actual time=0.010..0.035 rows=114 loops=1)
                    Recheck Cond: (amount <= 5)
                    Heap Blocks: exact=54
                    Buffers: shared hit=56
                    ->  Bitmap Index Scan on orders_amount_idx  (cost=0.00..2.19 rows=107 width=0) (actual time=0.004..0.004 rows=114 loops=1)
                          Index Cond: (amount <= 5)
                          Buffers: shared hit=2
        ->  Sort  (cost=75.61..76.02 rows=162 width=8) (actual time=0.041..0.045 rows=137 loops=1)
              Sort Key: orders_1.id, orders_1.amount
              Sort Method: quicksort  Memory: 28kB
              Buffers: shared hit=58
              ->  Bitmap Heap Scan on orders orders_1  (cost=2.64..69.67 rows=162 width=8) (actual time=0.009..0.032 rows=137 loops=1)
                    Recheck Cond: (amount >= 496)
                    Heap Blocks: exact=56
                    Buffers: shared hit=58
                    ->  Bitmap Index Scan on orders_amount_idx  (cost=0.00..2.60 rows=162 width=0) (actual time=0.004..0.004 rows=137 loops=1)
                          Index Cond: (amount >= 496)
                          Buffers: shared hit=2
Planning:
  Buffers: shared hit=6
Planning Time: 0.058 ms
Execution Time: 0.158 ms

```

### s19_union_all_predicate [shape=selective/indexable;spelling=union_all]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.053000 / 0.138000 ms.

```text
Sort  (cost=144.70..145.37 rows=269 width=8) (actual time=0.118..0.124 rows=251 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 30kB
  Buffers: shared hit=114
  ->  Append  (cost=2.21..133.84 rows=269 width=8) (actual time=0.011..0.093 rows=251 loops=1)
        Buffers: shared hit=114
        ->  Bitmap Heap Scan on orders  (cost=2.21..62.83 rows=107 width=8) (actual time=0.011..0.052 rows=114 loops=1)
              Recheck Cond: (amount <= 5)
              Heap Blocks: exact=54
              Buffers: shared hit=56
              ->  Bitmap Index Scan on orders_amount_idx  (cost=0.00..2.19 rows=107 width=0) (actual time=0.004..0.004 rows=114 loops=1)
                    Index Cond: (amount <= 5)
                    Buffers: shared hit=2
        ->  Bitmap Heap Scan on orders orders_1  (cost=2.64..69.67 rows=162 width=8) (actual time=0.008..0.030 rows=137 loops=1)
              Recheck Cond: (amount >= 496)
              Heap Blocks: exact=56
              Buffers: shared hit=58
              ->  Bitmap Index Scan on orders_amount_idx  (cost=0.00..2.60 rows=162 width=0) (actual time=0.004..0.004 rows=137 loops=1)
                    Index Cond: (amount >= 496)
                    Buffers: shared hit=2
Planning:
  Buffers: shared hit=6
Planning Time: 0.053 ms
Execution Time: 0.138 ms

```

## Limits and gaps recorded by the runner

- no concurrent writers: client-decode JSON bytes are not wire bytes; tails are under-sampled at the protocol trial count
- one host: no independent hardware, real-network, 3-node placement or equal-total-budget comparison is established
- single-node only; the YB 3-node RF=3 per-node/equal-total arms are a coverage gap on this host

