# SELECT access-path matrix (generated measurements)

_Measurements only. Verdicts for P10/P11 belong in a signed analysis._

**TL;DR:** within-run medians for OR / UNION ALL / UNION at three predicate shapes and OFFSET / keyset at five page depths, each cell verified against the independent oracle. A within-run difference is not a replicated result. Missing measurements are N/A.

| run | engine | topology | scale | page limit |
|---|---|---|---|---|
| select-order-access-postgres-3-20261004 | postgres | single | deep | 5 |

## Provenance

| Field | Value |
|---|---|
| repo commit | `ca2a38fbf9e53da00344578258138db9e01d3545` |
| benchmark image | `localhost/charitybench-select:select-order-access-postgres-3-20261004` |
| image digest | `8ecc060a8de698ebdf6b07ba9bd98c7663bdbb7dc6a74fa96b13a9609cd9664c` |
| run tag | `run/study-01/select-order-access-postgres-3-20261004` |
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

- default matrix: order seed `42930359`, 5 interleaved rounds
- expression-index arm: order seed `42688771`, 5 interleaved rounds
- default-matrix sequence (one sample of every cell per round):
  1. shape=selective/cast;spelling=or
  2. shape=selective/cast;spelling=union
  3. depth=10000;spelling=keyset
  4. shape=broad/indexable;spelling=union_all
  5. depth=10;spelling=offset
  6. shape=selective/expression;spelling=or
  7. shape=broad/expression;spelling=or
  8. shape=selective/indexable;spelling=union_all
  9. shape=broad/indexable;spelling=or
  10. shape=selective/indexable;spelling=or
  11. depth=100;spelling=offset
  12. shape=broad/cast;spelling=or
  13. depth=1000;spelling=keyset
  14. shape=broad/cast;spelling=union_all
  15. shape=broad/indexable;spelling=union
  16. depth=100;spelling=keyset
  17. depth=10;spelling=keyset
  18. shape=selective/expression;spelling=union_all
  19. depth=1;spelling=keyset
  20. shape=selective/cast;spelling=union_all
  21. shape=broad/expression;spelling=union
  22. shape=broad/cast;spelling=union
  23. depth=1000;spelling=offset
  24. depth=1;spelling=offset
  25. shape=selective/indexable;spelling=union
  26. shape=selective/expression;spelling=union
  27. depth=10000;spelling=offset
  28. shape=broad/expression;spelling=union_all
- index-arm sequence:
  1. shape=broad/expression;index=expression_index;spelling=or
  2. shape=selective/expression;index=expression_index;spelling=union_all
  3. shape=selective/expression;index=expression_index;spelling=or
  4. shape=selective/expression;index=expression_index;spelling=union
  5. shape=broad/expression;index=expression_index;spelling=union_all
  6. shape=broad/expression;index=expression_index;spelling=union

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
| broad/cast | default | or | valid | 1.324284 | 28.5 | 5 | 4652 | `a4bab1f05f5e667d` |
| broad/cast | default | union | valid | 1.976684 | 40.3 | 5 | 4652 | `6ff77b66c9bf6ee6` |
| broad/cast | default | union_all | valid | 1.793518 | 9.4 | 5 | 4652 | `3f9f3d5c5dc15dd7` |
| broad/expression | expression_index | or | valid | 1.341688 | 42.8 | 5 | 4652 | `a4bab1f05f5e667d` |
| broad/expression | expression_index | union | valid | 1.665216 | 30.0 | 5 | 4652 | `d4fbfb3aac234d58` |
| broad/expression | expression_index | union_all | valid | 1.572190 | 7.2 | 5 | 4652 | `605ee056e65421b1` |
| broad/expression | default | or | valid | 1.227995 | 28.0 | 5 | 4652 | `a4bab1f05f5e667d` |
| broad/expression | default | union | valid | 1.712651 | 26.1 | 5 | 4652 | `6ff77b66c9bf6ee6` |
| broad/expression | default | union_all | valid | 1.805940 | 30.1 | 5 | 4652 | `3f9f3d5c5dc15dd7` |
| broad/indexable | default | or | valid | 1.297230 | 36.5 | 5 | 4652 | `a4bab1f05f5e667d` |
| broad/indexable | default | union | valid | 1.788678 | 15.0 | 5 | 4652 | `ed39b768e5b063b0` |
| broad/indexable | default | union_all | valid | 1.781235 | 19.8 | 5 | 4652 | `d51936c88987c9fd` |
| selective/cast | default | or | valid | 0.878092 | 13.5 | 5 | 251 | `a4bab1f05f5e667d` |
| selective/cast | default | union | valid | 0.813842 | 17.9 | 5 | 251 | `6ff77b66c9bf6ee6` |
| selective/cast | default | union_all | valid | 0.861363 | 14.9 | 5 | 251 | `3f9f3d5c5dc15dd7` |
| selective/expression | expression_index | or | valid | 0.282496 | 22.3 | 5 | 251 | `d66da5af4c9f3384` |
| selective/expression | expression_index | union | valid | 0.332002 | 33.1 | 5 | 251 | `a5fbbd2993065c7a` |
| selective/expression | expression_index | union_all | valid | 0.422961 | 52.4 | 5 | 251 | `eeb1f1b0c0255959` |
| selective/expression | default | or | valid | 0.931551 | 13.1 | 5 | 251 | `a4bab1f05f5e667d` |
| selective/expression | default | union | valid | 0.869854 | 12.4 | 5 | 251 | `6ff77b66c9bf6ee6` |
| selective/expression | default | union_all | valid | 0.873422 | 7.4 | 5 | 251 | `3f9f3d5c5dc15dd7` |
| selective/indexable | default | or | valid | 0.374895 | 26.5 | 5 | 251 | `8bbb087a6c122dbe` |
| selective/indexable | default | union | valid | 0.428903 | 28.7 | 5 | 251 | `567de4baf4feaaf7` |
| selective/indexable | default | union_all | valid | 0.375063 | 32.0 | 5 | 251 | `e060188a3bcf6eb3` |

### P10 pairs

Pairs compare the **default** arm only; the labelled `expression_index` arm is a package arm and is never pooled with the query-text pair.

| pair | A median ms | B median ms | B/A | same plan | within-run verdict |
|---|---|---|---|---|---|
| shape=selective/indexable or vs union_all | 0.374895 | 0.375063 | 1.0004 | false | medians inside the two cells' own spread; no resolved difference within this run |
| shape=selective/indexable or vs union | 0.374895 | 0.428903 | 1.1441 | false | medians inside the two cells' own spread; no resolved difference within this run |
| shape=selective/expression or vs union_all | 0.931551 | 0.873422 | 0.9376 | false | medians inside the two cells' own spread; no resolved difference within this run |
| shape=selective/expression or vs union | 0.931551 | 0.869854 | 0.9338 | false | medians inside the two cells' own spread; no resolved difference within this run |
| shape=selective/cast or vs union_all | 0.878092 | 0.861363 | 0.9809 | false | medians inside the two cells' own spread; no resolved difference within this run |
| shape=selective/cast or vs union | 0.878092 | 0.813842 | 0.9268 | false | medians inside the two cells' own spread; no resolved difference within this run |
| shape=broad/indexable or vs union_all | 1.297230 | 1.781235 | 1.3731 | false | medians inside the two cells' own spread; no resolved difference within this run |
| shape=broad/indexable or vs union | 1.297230 | 1.788678 | 1.3788 | false | medians inside the two cells' own spread; no resolved difference within this run |
| shape=broad/expression or vs union_all | 1.227995 | 1.805940 | 1.4706 | false | medians inside the two cells' own spread; no resolved difference within this run |
| shape=broad/expression or vs union | 1.227995 | 1.712651 | 1.3947 | false | medians inside the two cells' own spread; no resolved difference within this run |
| shape=broad/cast or vs union_all | 1.324284 | 1.793518 | 1.3543 | false | medians inside the two cells' own spread; no resolved difference within this run |
| shape=broad/cast or vs union | 1.324284 | 1.976684 | 1.4926 | false | medians inside the two cells' own spread; no resolved difference within this run |

## Pagination (P11)

| depth | spelling | status | median ms | spread % | trials | rows | plan structure |
|---|---|---|---|---|---|---|---|
| 1 | offset | valid | 0.845620 | 8.5 | 5 | 5 | `ac4509d295bf3348` |
| 1 | keyset | valid | 0.903688 | 25.4 | 5 | 5 | `ac4509d295bf3348` |
| 10 | offset | valid | 0.899717 | 10.9 | 5 | 5 | `ac4509d295bf3348` |
| 10 | keyset | valid | 0.983635 | 18.8 | 5 | 5 | `ac4509d295bf3348` |
| 100 | offset | valid | 0.884186 | 13.9 | 5 | 5 | `ac4509d295bf3348` |
| 100 | keyset | valid | 0.952554 | 32.3 | 5 | 5 | `ac4509d295bf3348` |
| 1000 | offset | valid | 1.286011 | 9.4 | 5 | 5 | `ac4509d295bf3348` |
| 1000 | keyset | valid | 1.041333 | 17.9 | 5 | 5 | `ac4509d295bf3348` |
| 10000 | offset | valid | 2.074541 | 42.2 | 5 | 5 | `ac4509d295bf3348` |
| 10000 | keyset | valid | 0.469857 | 31.7 | 5 | 5 | `dd316474b64ddc29` |

Median ratio table (keyset / OFFSET; a ratio below 1 means keyset was faster in this run):

| depth | ratio (keyset/offset) | offset median ms | keyset median ms |
|---|---|---|---|
| 1 | 1.0687 | 0.845620 | 0.903688 |
| 10 | 1.0933 | 0.899717 | 0.983635 |
| 100 | 1.0773 | 0.884186 | 0.952554 |
| 1000 | 0.8097 | 1.286011 | 1.041333 |
| 10000 | 0.2265 | 2.074541 | 0.469857 |

### P11 pairs

| pair | A median ms | B median ms | B/A | same plan | within-run verdict |
|---|---|---|---|---|---|
| depth=1 offset vs keyset | 0.845620 | 0.903688 | 1.0687 | true | medians inside the two cells' own spread; no resolved difference within this run |
| depth=10 offset vs keyset | 0.899717 | 0.983635 | 1.0933 | true | medians inside the two cells' own spread; no resolved difference within this run |
| depth=100 offset vs keyset | 0.884186 | 0.952554 | 1.0773 | true | medians inside the two cells' own spread; no resolved difference within this run |
| depth=1000 offset vs keyset | 1.286011 | 1.041333 | 0.8097 | true | medians inside the two cells' own spread; no resolved difference within this run |
| depth=10000 offset vs keyset | 2.074541 | 0.469857 | 0.2265 | false | within-run medians differ; whole-run replication required before a verdict |

## Repeated-design position control

Not measured in this run.

## Expression-index arm: write and storage cost

Index `orders_amount_expr_idx` on expression `(amount + 0)`; supported=true.

| measure | value |
|---|---|
| index size (bytes) | 114688 |
| insert batch with index (ms) | 3.694945 |
| insert batch without index (ms) | 3.218416 |
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

Rows: 5; serialized decoded-row JSON bytes: 51 (not wire bytes). Planning/execution: 0.046000 / 0.321000 ms.

```text
Limit  (cost=154.12..154.13 rows=5 width=8) (actual time=0.312..0.313 rows=5 loops=1)
  Buffers: shared hit=73
  ->  Sort  (cost=154.12..158.90 rows=1912 width=8) (actual time=0.312..0.312 rows=5 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 25kB
        Buffers: shared hit=73
        ->  Bitmap Heap Scan on orders  (cost=22.66..122.36 rows=1912 width=8) (actual time=0.039..0.212 rows=2000 loops=1)
              Recheck Cond: ((amount < 67) OR (amount = 67) OR (amount IS NULL))
              Filter: ((amount < 67) OR ((amount = 67) AND (id > 4310)) OR (amount IS NULL))
              Rows Removed by Filter: 6
              Heap Blocks: exact=65
              Buffers: shared hit=73
              ->  BitmapOr  (cost=22.66..22.66 rows=1983 width=0) (actual time=0.030..0.031 rows=0 loops=1)
                    Buffers: shared hit=8
                    ->  Bitmap Index Scan on orders_amount_idx  (cost=0.00..14.38 rows=1440 width=0) (actual time=0.023..0.023 rows=1460 loops=1)
                          Index Cond: (amount < 67)
                          Buffers: shared hit=4
                    ->  Bitmap Index Scan on orders_amount_idx  (cost=0.00..1.54 rows=21 width=0) (actual time=0.002..0.002 rows=24 loops=1)
                          Index Cond: (amount = 67)
                          Buffers: shared hit=2
                    ->  Bitmap Index Scan on orders_amount_idx  (cost=0.00..5.30 rows=522 width=0) (actual time=0.006..0.006 rows=522 loops=1)
                          Index Cond: (amount IS NULL)
                          Buffers: shared hit=2
Planning Time: 0.046 ms
Execution Time: 0.321 ms

```

### s16_page_offset [depth=10000;spelling=offset]

Rows: 5; serialized decoded-row JSON bytes: 51 (not wire bytes). Planning/execution: 0.026000 / 2.530000 ms.

```text
Limit  (cost=1023.04..1023.06 rows=5 width=8) (actual time=2.524..2.525 rows=5 loops=1)
  Buffers: shared hit=65
  ->  Sort  (cost=998.04..1028.04 rows=12000 width=8) (actual time=2.116..2.342 rows=10005 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: quicksort  Memory: 853kB
        Buffers: shared hit=65
        ->  Seq Scan on orders  (cost=0.00..185.00 rows=12000 width=8) (actual time=0.004..0.474 rows=12000 loops=1)
              Buffers: shared hit=65
Planning Time: 0.026 ms
Execution Time: 2.530 ms

```

### s17_page_keyset [depth=1000;spelling=keyset]

Rows: 5; serialized decoded-row JSON bytes: 59 (not wire bytes). Planning/execution: 0.036000 / 1.082000 ms.

```text
Limit  (cost=450.03..450.04 rows=5 width=8) (actual time=1.077..1.078 rows=5 loops=1)
  Buffers: shared hit=65
  ->  Sort  (cost=450.03..476.38 rows=10538 width=8) (actual time=1.077..1.077 rows=5 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 25kB
        Buffers: shared hit=65
        ->  Seq Scan on orders  (cost=0.00..275.00 rows=10538 width=8) (actual time=0.005..0.552 rows=11000 loops=1)
              Filter: ((amount < 457) OR ((amount = 457) AND (id > 11254)) OR (amount IS NULL))
              Rows Removed by Filter: 1000
              Buffers: shared hit=65
Planning Time: 0.036 ms
Execution Time: 1.082 ms

```

### s16_page_offset [depth=1000;spelling=offset]

Rows: 5; serialized decoded-row JSON bytes: 59 (not wire bytes). Planning/execution: 0.024000 / 1.511000 ms.

```text
Limit  (cost=845.88..845.89 rows=5 width=8) (actual time=1.505..1.506 rows=5 loops=1)
  Buffers: shared hit=65
  ->  Sort  (cost=843.38..873.38 rows=12000 width=8) (actual time=1.466..1.487 rows=1005 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 115kB
        Buffers: shared hit=65
        ->  Seq Scan on orders  (cost=0.00..185.00 rows=12000 width=8) (actual time=0.005..0.497 rows=12000 loops=1)
              Buffers: shared hit=65
Planning Time: 0.024 ms
Execution Time: 1.511 ms

```

### s17_page_keyset [depth=100;spelling=keyset]

Rows: 5; serialized decoded-row JSON bytes: 56 (not wire bytes). Planning/execution: 0.030000 / 1.097000 ms.

```text
Limit  (cost=463.65..463.66 rows=5 width=8) (actual time=1.092..1.092 rows=5 loops=1)
  Buffers: shared hit=65
  ->  Sort  (cost=463.65..492.05 rows=11358 width=8) (actual time=1.091..1.092 rows=5 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 25kB
        Buffers: shared hit=65
        ->  Seq Scan on orders  (cost=0.00..275.00 rows=11358 width=8) (actual time=0.003..0.537 rows=11900 loops=1)
              Filter: ((amount < 497) OR ((amount = 497) AND (id > 6723)) OR (amount IS NULL))
              Rows Removed by Filter: 100
              Buffers: shared hit=65
Planning:
  Buffers: shared hit=3
Planning Time: 0.030 ms
Execution Time: 1.097 ms

```

### s16_page_offset [depth=100;spelling=offset]

Rows: 5; serialized decoded-row JSON bytes: 56 (not wire bytes). Planning/execution: 0.020000 / 1.139000 ms.

```text
Limit  (cost=648.10..648.12 rows=5 width=8) (actual time=1.134..1.135 rows=5 loops=1)
  Buffers: shared hit=65
  ->  Sort  (cost=647.85..677.85 rows=12000 width=8) (actual time=1.130..1.132 rows=105 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 31kB
        Buffers: shared hit=65
        ->  Seq Scan on orders  (cost=0.00..185.00 rows=12000 width=8) (actual time=0.004..0.485 rows=12000 loops=1)
              Buffers: shared hit=65
Planning Time: 0.020 ms
Execution Time: 1.139 ms

```

### s17_page_keyset [depth=10;spelling=keyset]

Rows: 5; serialized decoded-row JSON bytes: 56 (not wire bytes). Planning/execution: 0.055000 / 1.185000 ms.

```text
Limit  (cost=465.21..465.23 rows=5 width=8) (actual time=1.178..1.179 rows=5 loops=1)
  Buffers: shared hit=65
  ->  Sort  (cost=465.21..493.84 rows=11452 width=8) (actual time=1.178..1.178 rows=5 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 25kB
        Buffers: shared hit=65
        ->  Seq Scan on orders  (cost=0.00..275.00 rows=11452 width=8) (actual time=0.006..0.581 rows=11990 loops=1)
              Filter: ((amount < 500) OR ((amount = 500) AND (id > 6319)) OR (amount IS NULL))
              Rows Removed by Filter: 10
              Buffers: shared hit=65
Planning:
  Buffers: shared hit=3
Planning Time: 0.055 ms
Execution Time: 1.185 ms

```

### s16_page_offset [depth=10;spelling=offset]

Rows: 5; serialized decoded-row JSON bytes: 56 (not wire bytes). Planning/execution: 0.019000 / 1.132000 ms.

```text
Limit  (cost=479.44..479.45 rows=5 width=8) (actual time=1.126..1.127 rows=5 loops=1)
  Buffers: shared hit=65
  ->  Sort  (cost=479.41..509.41 rows=12000 width=8) (actual time=1.125..1.126 rows=15 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 25kB
        Buffers: shared hit=65
        ->  Seq Scan on orders  (cost=0.00..185.00 rows=12000 width=8) (actual time=0.004..0.498 rows=12000 loops=1)
              Buffers: shared hit=65
Planning Time: 0.019 ms
Execution Time: 1.132 ms

```

### s17_page_keyset [depth=1;spelling=keyset]

Rows: 5; serialized decoded-row JSON bytes: 56 (not wire bytes). Planning/execution: 0.044000 / 1.216000 ms.

```text
Limit  (cost=465.21..465.23 rows=5 width=8) (actual time=1.209..1.210 rows=5 loops=1)
  Buffers: shared hit=65
  ->  Sort  (cost=465.21..493.84 rows=11452 width=8) (actual time=1.208..1.209 rows=5 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 25kB
        Buffers: shared hit=65
        ->  Seq Scan on orders  (cost=0.00..275.00 rows=11452 width=8) (actual time=0.004..0.596 rows=11999 loops=1)
              Filter: ((amount < 500) OR ((amount = 500) AND (id > 1554)) OR (amount IS NULL))
              Rows Removed by Filter: 1
              Buffers: shared hit=65
Planning:
  Buffers: shared hit=3
Planning Time: 0.044 ms
Execution Time: 1.216 ms

```

### s16_page_offset [depth=1;spelling=offset]

Rows: 5; serialized decoded-row JSON bytes: 56 (not wire bytes). Planning/execution: 0.022000 / 1.062000 ms.

```text
Limit  (cost=400.10..400.11 rows=5 width=8) (actual time=1.057..1.058 rows=5 loops=1)
  Buffers: shared hit=65
  ->  Sort  (cost=400.10..430.10 rows=12000 width=8) (actual time=1.057..1.057 rows=6 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 25kB
        Buffers: shared hit=65
        ->  Seq Scan on orders  (cost=0.00..185.00 rows=12000 width=8) (actual time=0.004..0.472 rows=12000 loops=1)
              Buffers: shared hit=65
Planning Time: 0.022 ms
Execution Time: 1.062 ms

```

### s18_or_predicate [shape=broad/cast;spelling=or]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.032000 / 0.879000 ms.

```text
Index Scan using orders_pkey on orders  (cost=0.29..403.89 rows=6667 width=8) (actual time=0.008..0.776 rows=4652 loops=1)
  Filter: (((amount)::bigint <= 100) OR ((amount)::bigint >= 400))
  Rows Removed by Filter: 7348
  Buffers: shared hit=99
Planning Time: 0.032 ms
Execution Time: 0.879 ms

```

### s20_union_predicate [shape=broad/cast;spelling=union]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.060000 / 1.729000 ms.

```text
Unique  (cost=968.64..1108.64 rows=8000 width=8) (actual time=1.003..1.626 rows=4652 loops=1)
  Buffers: shared hit=130
  ->  Merge Append  (cost=968.64..1068.64 rows=8000 width=8) (actual time=1.003..1.354 rows=4652 loops=1)
        Sort Key: orders.id, orders.amount
        Buffers: shared hit=130
        ->  Sort  (cost=484.32..494.32 rows=4000 width=8) (actual time=0.481..0.533 rows=2273 loops=1)
              Sort Key: orders.id, orders.amount
              Sort Method: quicksort  Memory: 150kB
              Buffers: shared hit=65
              ->  Seq Scan on orders  (cost=0.00..245.00 rows=4000 width=8) (actual time=0.007..0.363 rows=2273 loops=1)
                    Filter: ((amount)::bigint <= 100)
                    Rows Removed by Filter: 9727
                    Buffers: shared hit=65
        ->  Sort  (cost=484.32..494.32 rows=4000 width=8) (actual time=0.521..0.571 rows=2379 loops=1)
              Sort Key: orders_1.id, orders_1.amount
              Sort Method: quicksort  Memory: 152kB
              Buffers: shared hit=65
              ->  Seq Scan on orders orders_1  (cost=0.00..245.00 rows=4000 width=8) (actual time=0.002..0.382 rows=2379 loops=1)
                    Filter: ((amount)::bigint >= 400)
                    Rows Removed by Filter: 9621
                    Buffers: shared hit=65
Planning Time: 0.060 ms
Execution Time: 1.729 ms

```

### s19_union_all_predicate [shape=broad/cast;spelling=union_all]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.045000 / 1.441000 ms.

```text
Sort  (cost=1048.63..1068.63 rows=8000 width=8) (actual time=1.243..1.346 rows=4652 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 302kB
  Buffers: shared hit=130
  ->  Append  (cost=0.00..530.00 rows=8000 width=8) (actual time=0.007..0.890 rows=4652 loops=1)
        Buffers: shared hit=130
        ->  Seq Scan on orders  (cost=0.00..245.00 rows=4000 width=8) (actual time=0.007..0.367 rows=2273 loops=1)
              Filter: ((amount)::bigint <= 100)
              Rows Removed by Filter: 9727
              Buffers: shared hit=65
        ->  Seq Scan on orders orders_1  (cost=0.00..245.00 rows=4000 width=8) (actual time=0.002..0.348 rows=2379 loops=1)
              Filter: ((amount)::bigint >= 400)
              Rows Removed by Filter: 9621
              Buffers: shared hit=65
Planning Time: 0.045 ms
Execution Time: 1.441 ms

```

### s18_or_predicate [shape=broad/expression;index=expression_index;spelling=or]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.038000 / 0.900000 ms.

```text
Index Scan using orders_pkey on orders  (cost=0.29..403.89 rows=4211 width=8) (actual time=0.007..0.794 rows=4652 loops=1)
  Filter: (((amount + 0) <= 100) OR ((amount + 0) >= 400))
  Rows Removed by Filter: 7348
  Buffers: shared hit=99
Planning Time: 0.038 ms
Execution Time: 0.900 ms

```

### s20_union_predicate [shape=broad/expression;index=expression_index;spelling=union]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.064000 / 1.797000 ms.

```text
Unique  (cost=492.69..574.31 rows=4664 width=8) (actual time=1.008..1.690 rows=4652 loops=1)
  Buffers: shared hit=3934
  ->  Merge Append  (cost=492.69..550.99 rows=4664 width=8) (actual time=1.008..1.401 rows=4652 loops=1)
        Sort Key: orders.id, orders.amount
        Buffers: shared hit=3934
        ->  Sort  (cost=241.10..246.77 rows=2267 width=8) (actual time=0.516..0.573 rows=2273 loops=1)
              Sort Key: orders.id, orders.amount
              Sort Method: quicksort  Memory: 150kB
              Buffers: shared hit=1921
              ->  Index Scan using orders_amount_expr_idx on orders  (cost=0.29..114.75 rows=2267 width=8) (actual time=0.009..0.295 rows=2273 loops=1)
                    Index Cond: ((amount + 0) <= 100)
                    Buffers: shared hit=1921
        ->  Sort  (cost=251.58..257.57 rows=2397 width=8) (actual time=0.492..0.548 rows=2379 loops=1)
              Sort Key: orders_1.id, orders_1.amount
              Sort Method: quicksort  Memory: 152kB
              Buffers: shared hit=2013
              ->  Index Scan using orders_amount_expr_idx on orders orders_1  (cost=0.29..117.03 rows=2397 width=8) (actual time=0.003..0.268 rows=2379 loops=1)
                    Index Cond: ((amount + 0) >= 400)
                    Buffers: shared hit=2013
Planning Time: 0.064 ms
Execution Time: 1.797 ms

```

### s19_union_all_predicate [shape=broad/expression;index=expression_index;spelling=union_all]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.056000 / 1.444000 ms.

```text
Sort  (cost=539.31..550.97 rows=4664 width=8) (actual time=1.219..1.331 rows=4652 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 302kB
  Buffers: shared hit=3934
  ->  Append  (cost=0.29..255.10 rows=4664 width=8) (actual time=0.009..0.743 rows=4652 loops=1)
        Buffers: shared hit=3934
        ->  Index Scan using orders_amount_expr_idx on orders  (cost=0.29..114.75 rows=2267 width=8) (actual time=0.008..0.283 rows=2273 loops=1)
              Index Cond: ((amount + 0) <= 100)
              Buffers: shared hit=1921
        ->  Index Scan using orders_amount_expr_idx on orders orders_1  (cost=0.29..117.03 rows=2397 width=8) (actual time=0.003..0.269 rows=2379 loops=1)
              Index Cond: ((amount + 0) >= 400)
              Buffers: shared hit=2013
Planning Time: 0.056 ms
Execution Time: 1.444 ms

```

### s18_or_predicate [shape=broad/expression;spelling=or]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.032000 / 0.937000 ms.

```text
Index Scan using orders_pkey on orders  (cost=0.29..403.89 rows=6667 width=8) (actual time=0.008..0.827 rows=4652 loops=1)
  Filter: (((amount + 0) <= 100) OR ((amount + 0) >= 400))
  Rows Removed by Filter: 7348
  Buffers: shared hit=99
Planning Time: 0.032 ms
Execution Time: 0.937 ms

```

### s20_union_predicate [shape=broad/expression;spelling=union]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.046000 / 1.653000 ms.

```text
Unique  (cost=968.64..1108.64 rows=8000 width=8) (actual time=0.933..1.554 rows=4652 loops=1)
  Buffers: shared hit=130
  ->  Merge Append  (cost=968.64..1068.64 rows=8000 width=8) (actual time=0.932..1.280 rows=4652 loops=1)
        Sort Key: orders.id, orders.amount
        Buffers: shared hit=130
        ->  Sort  (cost=484.32..494.32 rows=4000 width=8) (actual time=0.465..0.516 rows=2273 loops=1)
              Sort Key: orders.id, orders.amount
              Sort Method: quicksort  Memory: 150kB
              Buffers: shared hit=65
              ->  Seq Scan on orders  (cost=0.00..245.00 rows=4000 width=8) (actual time=0.006..0.346 rows=2273 loops=1)
                    Filter: ((amount + 0) <= 100)
                    Rows Removed by Filter: 9727
                    Buffers: shared hit=65
        ->  Sort  (cost=484.32..494.32 rows=4000 width=8) (actual time=0.466..0.517 rows=2379 loops=1)
              Sort Key: orders_1.id, orders_1.amount
              Sort Method: quicksort  Memory: 152kB
              Buffers: shared hit=65
              ->  Seq Scan on orders orders_1  (cost=0.00..245.00 rows=4000 width=8) (actual time=0.002..0.336 rows=2379 loops=1)
                    Filter: ((amount + 0) >= 400)
                    Rows Removed by Filter: 9621
                    Buffers: shared hit=65
Planning Time: 0.046 ms
Execution Time: 1.653 ms

```

### s19_union_all_predicate [shape=broad/expression;spelling=union_all]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.059000 / 1.489000 ms.

```text
Sort  (cost=1048.63..1068.63 rows=8000 width=8) (actual time=1.278..1.377 rows=4652 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 302kB
  Buffers: shared hit=130
  ->  Append  (cost=0.00..530.00 rows=8000 width=8) (actual time=0.009..0.906 rows=4652 loops=1)
        Buffers: shared hit=130
        ->  Seq Scan on orders  (cost=0.00..245.00 rows=4000 width=8) (actual time=0.008..0.381 rows=2273 loops=1)
              Filter: ((amount + 0) <= 100)
              Rows Removed by Filter: 9727
              Buffers: shared hit=65
        ->  Seq Scan on orders orders_1  (cost=0.00..245.00 rows=4000 width=8) (actual time=0.002..0.346 rows=2379 loops=1)
              Filter: ((amount + 0) >= 400)
              Rows Removed by Filter: 9621
              Buffers: shared hit=65
Planning Time: 0.059 ms
Execution Time: 1.489 ms

```

### s18_or_predicate [shape=broad/indexable;spelling=or]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.039000 / 0.822000 ms.

```text
Index Scan using orders_pkey on orders  (cost=0.29..343.89 rows=4211 width=8) (actual time=0.009..0.719 rows=4652 loops=1)
  Filter: ((amount <= 100) OR (amount >= 400))
  Rows Removed by Filter: 7348
  Buffers: shared hit=99
Planning Time: 0.039 ms
Execution Time: 0.822 ms

```

### s20_union_predicate [shape=broad/indexable;spelling=union]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.063000 / 1.838000 ms.

```text
Unique  (cost=495.99..577.61 rows=4664 width=8) (actual time=1.028..1.728 rows=4652 loops=1)
  Buffers: shared hit=3936
  ->  Merge Append  (cost=495.99..554.29 rows=4664 width=8) (actual time=1.028..1.449 rows=4652 loops=1)
        Sort Key: orders.id, orders.amount
        Buffers: shared hit=3936
        ->  Sort  (cost=242.20..247.87 rows=2267 width=8) (actual time=0.542..0.601 rows=2273 loops=1)
              Sort Key: orders.id, orders.amount
              Sort Method: quicksort  Memory: 150kB
              Buffers: shared hit=1922
              ->  Index Scan using orders_amount_idx on orders  (cost=0.29..115.85 rows=2267 width=8) (actual time=0.010..0.302 rows=2273 loops=1)
                    Index Cond: (amount <= 100)
                    Buffers: shared hit=1922
        ->  Sort  (cost=253.78..259.77 rows=2397 width=8) (actual time=0.485..0.543 rows=2379 loops=1)
              Sort Key: orders_1.id, orders_1.amount
              Sort Method: quicksort  Memory: 152kB
              Buffers: shared hit=2014
              ->  Index Scan using orders_amount_idx on orders orders_1  (cost=0.29..119.23 rows=2397 width=8) (actual time=0.004..0.264 rows=2379 loops=1)
                    Index Cond: (amount >= 400)
                    Buffers: shared hit=2014
Planning Time: 0.063 ms
Execution Time: 1.838 ms

```

### s19_union_all_predicate [shape=broad/indexable;spelling=union_all]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.050000 / 1.485000 ms.

```text
Sort  (cost=542.61..554.27 rows=4664 width=8) (actual time=1.263..1.379 rows=4652 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 302kB
  Buffers: shared hit=3936
  ->  Append  (cost=0.29..258.40 rows=4664 width=8) (actual time=0.008..0.803 rows=4652 loops=1)
        Buffers: shared hit=3936
        ->  Index Scan using orders_amount_idx on orders  (cost=0.29..115.85 rows=2267 width=8) (actual time=0.008..0.341 rows=2273 loops=1)
              Index Cond: (amount <= 100)
              Buffers: shared hit=1922
        ->  Index Scan using orders_amount_idx on orders orders_1  (cost=0.29..119.23 rows=2397 width=8) (actual time=0.005..0.263 rows=2379 loops=1)
              Index Cond: (amount >= 400)
              Buffers: shared hit=2014
Planning Time: 0.050 ms
Execution Time: 1.485 ms

```

### s18_or_predicate [shape=selective/cast;spelling=or]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.026000 / 0.685000 ms.

```text
Index Scan using orders_pkey on orders  (cost=0.29..403.89 rows=6667 width=8) (actual time=0.009..0.672 rows=251 loops=1)
  Filter: (((amount)::bigint <= 5) OR ((amount)::bigint >= 496))
  Rows Removed by Filter: 11749
  Buffers: shared hit=99
Planning Time: 0.026 ms
Execution Time: 0.685 ms

```

### s20_union_predicate [shape=selective/cast;spelling=union]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.042000 / 0.663000 ms.

```text
Unique  (cost=968.64..1108.64 rows=8000 width=8) (actual time=0.611..0.649 rows=251 loops=1)
  Buffers: shared hit=130
  ->  Merge Append  (cost=968.64..1068.64 rows=8000 width=8) (actual time=0.611..0.632 rows=251 loops=1)
        Sort Key: orders.id, orders.amount
        Buffers: shared hit=130
        ->  Sort  (cost=484.32..494.32 rows=4000 width=8) (actual time=0.329..0.331 rows=114 loops=1)
              Sort Key: orders.id, orders.amount
              Sort Method: quicksort  Memory: 27kB
              Buffers: shared hit=65
              ->  Seq Scan on orders  (cost=0.00..245.00 rows=4000 width=8) (actual time=0.009..0.321 rows=114 loops=1)
                    Filter: ((amount)::bigint <= 5)
                    Rows Removed by Filter: 11886
                    Buffers: shared hit=65
        ->  Sort  (cost=484.32..494.32 rows=4000 width=8) (actual time=0.282..0.285 rows=137 loops=1)
              Sort Key: orders_1.id, orders_1.amount
              Sort Method: quicksort  Memory: 28kB
              Buffers: shared hit=65
              ->  Seq Scan on orders orders_1  (cost=0.00..245.00 rows=4000 width=8) (actual time=0.003..0.274 rows=137 loops=1)
                    Filter: ((amount)::bigint >= 496)
                    Rows Removed by Filter: 11863
                    Buffers: shared hit=65
Planning Time: 0.042 ms
Execution Time: 0.663 ms

```

### s19_union_all_predicate [shape=selective/cast;spelling=union_all]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.035000 / 0.589000 ms.

```text
Sort  (cost=1048.63..1068.63 rows=8000 width=8) (actual time=0.572..0.577 rows=251 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 30kB
  Buffers: shared hit=130
  ->  Append  (cost=0.00..530.00 rows=8000 width=8) (actual time=0.006..0.551 rows=251 loops=1)
        Buffers: shared hit=130
        ->  Seq Scan on orders  (cost=0.00..245.00 rows=4000 width=8) (actual time=0.006..0.282 rows=114 loops=1)
              Filter: ((amount)::bigint <= 5)
              Rows Removed by Filter: 11886
              Buffers: shared hit=65
        ->  Seq Scan on orders orders_1  (cost=0.00..245.00 rows=4000 width=8) (actual time=0.003..0.259 rows=137 loops=1)
              Filter: ((amount)::bigint >= 496)
              Rows Removed by Filter: 11863
              Buffers: shared hit=65
Planning Time: 0.035 ms
Execution Time: 0.589 ms

```

### s18_or_predicate [shape=selective/expression;index=expression_index;spelling=or]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.037000 / 0.076000 ms.

```text
Sort  (cost=86.11..86.78 rows=268 width=8) (actual time=0.059..0.066 rows=251 loops=1)
  Sort Key: id
  Sort Method: quicksort  Memory: 30kB
  Buffers: shared hit=67
  ->  Bitmap Heap Scan on orders  (cost=4.92..75.30 rows=268 width=8) (actual time=0.012..0.045 rows=251 loops=1)
        Recheck Cond: (((amount + 0) <= 5) OR ((amount + 0) >= 496))
        Heap Blocks: exact=63
        Buffers: shared hit=67
        ->  BitmapOr  (cost=4.92..4.92 rows=269 width=0) (actual time=0.006..0.006 rows=0 loops=1)
              Buffers: shared hit=4
              ->  Bitmap Index Scan on orders_amount_expr_idx  (cost=0.00..2.19 rows=107 width=0) (actual time=0.004..0.004 rows=114 loops=1)
                    Index Cond: ((amount + 0) <= 5)
                    Buffers: shared hit=2
              ->  Bitmap Index Scan on orders_amount_expr_idx  (cost=0.00..2.60 rows=162 width=0) (actual time=0.002..0.002 rows=137 loops=1)
                    Index Cond: ((amount + 0) >= 496)
                    Buffers: shared hit=2
Planning:
  Buffers: shared hit=6
Planning Time: 0.037 ms
Execution Time: 0.076 ms

```

### s20_union_predicate [shape=selective/expression;index=expression_index;spelling=union]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.054000 / 0.121000 ms.

```text
Unique  (cost=142.73..147.44 rows=269 width=8) (actual time=0.070..0.108 rows=251 loops=1)
  Buffers: shared hit=114
  ->  Merge Append  (cost=142.73..146.09 rows=269 width=8) (actual time=0.070..0.091 rows=251 loops=1)
        Sort Key: orders.id, orders.amount
        Buffers: shared hit=114
        ->  Sort  (cost=66.70..66.97 rows=107 width=8) (actual time=0.035..0.038 rows=114 loops=1)
              Sort Key: orders.id, orders.amount
              Sort Method: quicksort  Memory: 27kB
              Buffers: shared hit=56
              ->  Bitmap Heap Scan on orders  (cost=2.21..63.10 rows=107 width=8) (actual time=0.008..0.028 rows=114 loops=1)
                    Recheck Cond: ((amount + 0) <= 5)
                    Heap Blocks: exact=54
                    Buffers: shared hit=56
                    ->  Bitmap Index Scan on orders_amount_expr_idx  (cost=0.00..2.19 rows=107 width=0) (actual time=0.003..0.003 rows=114 loops=1)
                          Index Cond: ((amount + 0) <= 5)
                          Buffers: shared hit=2
        ->  Sort  (cost=76.02..76.42 rows=162 width=8) (actual time=0.034..0.037 rows=137 loops=1)
              Sort Key: orders_1.id, orders_1.amount
              Sort Method: quicksort  Memory: 28kB
              Buffers: shared hit=58
              ->  Bitmap Heap Scan on orders orders_1  (cost=2.64..70.07 rows=162 width=8) (actual time=0.007..0.026 rows=137 loops=1)
                    Recheck Cond: ((amount + 0) >= 496)
                    Heap Blocks: exact=56
                    Buffers: shared hit=58
                    ->  Bitmap Index Scan on orders_amount_expr_idx  (cost=0.00..2.60 rows=162 width=0) (actual time=0.003..0.003 rows=137 loops=1)
                          Index Cond: ((amount + 0) >= 496)
                          Buffers: shared hit=2
Planning:
  Buffers: shared hit=6
Planning Time: 0.054 ms
Execution Time: 0.121 ms

```

### s19_union_all_predicate [shape=selective/expression;index=expression_index;spelling=union_all]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.070000 / 0.149000 ms.

```text
Sort  (cost=145.37..146.04 rows=269 width=8) (actual time=0.125..0.131 rows=251 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 30kB
  Buffers: shared hit=114
  ->  Append  (cost=2.21..134.51 rows=269 width=8) (actual time=0.013..0.100 rows=251 loops=1)
        Buffers: shared hit=114
        ->  Bitmap Heap Scan on orders  (cost=2.21..63.10 rows=107 width=8) (actual time=0.012..0.050 rows=114 loops=1)
              Recheck Cond: ((amount + 0) <= 5)
              Heap Blocks: exact=54
              Buffers: shared hit=56
              ->  Bitmap Index Scan on orders_amount_expr_idx  (cost=0.00..2.19 rows=107 width=0) (actual time=0.005..0.005 rows=114 loops=1)
                    Index Cond: ((amount + 0) <= 5)
                    Buffers: shared hit=2
        ->  Bitmap Heap Scan on orders orders_1  (cost=2.64..70.07 rows=162 width=8) (actual time=0.011..0.038 rows=137 loops=1)
              Recheck Cond: ((amount + 0) >= 496)
              Heap Blocks: exact=56
              Buffers: shared hit=58
              ->  Bitmap Index Scan on orders_amount_expr_idx  (cost=0.00..2.60 rows=162 width=0) (actual time=0.007..0.007 rows=137 loops=1)
                    Index Cond: ((amount + 0) >= 496)
                    Buffers: shared hit=2
Planning:
  Buffers: shared hit=6
Planning Time: 0.070 ms
Execution Time: 0.149 ms

```

### s18_or_predicate [shape=selective/expression;spelling=or]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.027000 / 0.667000 ms.

```text
Index Scan using orders_pkey on orders  (cost=0.29..403.89 rows=6667 width=8) (actual time=0.009..0.655 rows=251 loops=1)
  Filter: (((amount + 0) <= 5) OR ((amount + 0) >= 496))
  Rows Removed by Filter: 11749
  Buffers: shared hit=99
Planning Time: 0.027 ms
Execution Time: 0.667 ms

```

### s20_union_predicate [shape=selective/expression;spelling=union]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.034000 / 0.627000 ms.

```text
Unique  (cost=968.64..1108.64 rows=8000 width=8) (actual time=0.582..0.615 rows=251 loops=1)
  Buffers: shared hit=130
  ->  Merge Append  (cost=968.64..1068.64 rows=8000 width=8) (actual time=0.582..0.600 rows=251 loops=1)
        Sort Key: orders.id, orders.amount
        Buffers: shared hit=130
        ->  Sort  (cost=484.32..494.32 rows=4000 width=8) (actual time=0.319..0.322 rows=114 loops=1)
              Sort Key: orders.id, orders.amount
              Sort Method: quicksort  Memory: 27kB
              Buffers: shared hit=65
              ->  Seq Scan on orders  (cost=0.00..245.00 rows=4000 width=8) (actual time=0.005..0.312 rows=114 loops=1)
                    Filter: ((amount + 0) <= 5)
                    Rows Removed by Filter: 11886
                    Buffers: shared hit=65
        ->  Sort  (cost=484.32..494.32 rows=4000 width=8) (actual time=0.262..0.265 rows=137 loops=1)
              Sort Key: orders_1.id, orders_1.amount
              Sort Method: quicksort  Memory: 28kB
              Buffers: shared hit=65
              ->  Seq Scan on orders orders_1  (cost=0.00..245.00 rows=4000 width=8) (actual time=0.003..0.255 rows=137 loops=1)
                    Filter: ((amount + 0) >= 496)
                    Rows Removed by Filter: 11863
                    Buffers: shared hit=65
Planning Time: 0.034 ms
Execution Time: 0.627 ms

```

### s19_union_all_predicate [shape=selective/expression;spelling=union_all]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.028000 / 0.602000 ms.

```text
Sort  (cost=1048.63..1068.63 rows=8000 width=8) (actual time=0.586..0.592 rows=251 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 30kB
  Buffers: shared hit=130
  ->  Append  (cost=0.00..530.00 rows=8000 width=8) (actual time=0.005..0.567 rows=251 loops=1)
        Buffers: shared hit=130
        ->  Seq Scan on orders  (cost=0.00..245.00 rows=4000 width=8) (actual time=0.004..0.287 rows=114 loops=1)
              Filter: ((amount + 0) <= 5)
              Rows Removed by Filter: 11886
              Buffers: shared hit=65
        ->  Seq Scan on orders orders_1  (cost=0.00..245.00 rows=4000 width=8) (actual time=0.003..0.270 rows=137 loops=1)
              Filter: ((amount + 0) >= 496)
              Rows Removed by Filter: 11863
              Buffers: shared hit=65
Planning Time: 0.028 ms
Execution Time: 0.602 ms

```

### s18_or_predicate [shape=selective/indexable;spelling=or]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.041000 / 0.092000 ms.

```text
Sort  (cost=84.77..85.44 rows=268 width=8) (actual time=0.072..0.079 rows=251 loops=1)
  Sort Key: id
  Sort Method: quicksort  Memory: 30kB
  Buffers: shared hit=67
  ->  Bitmap Heap Scan on orders  (cost=4.92..73.96 rows=268 width=8) (actual time=0.015..0.055 rows=251 loops=1)
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
Planning Time: 0.041 ms
Execution Time: 0.092 ms

```

### s20_union_predicate [shape=selective/indexable;spelling=union]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.147000 / 0.130000 ms.

```text
Unique  (cost=142.06..146.77 rows=269 width=8) (actual time=0.077..0.114 rows=251 loops=1)
  Buffers: shared hit=114
  ->  Merge Append  (cost=142.06..145.42 rows=269 width=8) (actual time=0.076..0.098 rows=251 loops=1)
        Sort Key: orders.id, orders.amount
        Buffers: shared hit=114
        ->  Sort  (cost=66.44..66.70 rows=107 width=8) (actual time=0.040..0.043 rows=114 loops=1)
              Sort Key: orders.id, orders.amount
              Sort Method: quicksort  Memory: 27kB
              Buffers: shared hit=56
              ->  Bitmap Heap Scan on orders  (cost=2.21..62.83 rows=107 width=8) (actual time=0.010..0.033 rows=114 loops=1)
                    Recheck Cond: (amount <= 5)
                    Heap Blocks: exact=54
                    Buffers: shared hit=56
                    ->  Bitmap Index Scan on orders_amount_idx  (cost=0.00..2.19 rows=107 width=0) (actual time=0.003..0.004 rows=114 loops=1)
                          Index Cond: (amount <= 5)
                          Buffers: shared hit=2
        ->  Sort  (cost=75.61..76.02 rows=162 width=8) (actual time=0.036..0.039 rows=137 loops=1)
              Sort Key: orders_1.id, orders_1.amount
              Sort Method: quicksort  Memory: 28kB
              Buffers: shared hit=58
              ->  Bitmap Heap Scan on orders orders_1  (cost=2.64..69.67 rows=162 width=8) (actual time=0.007..0.028 rows=137 loops=1)
                    Recheck Cond: (amount >= 496)
                    Heap Blocks: exact=56
                    Buffers: shared hit=58
                    ->  Bitmap Index Scan on orders_amount_idx  (cost=0.00..2.60 rows=162 width=0) (actual time=0.004..0.004 rows=137 loops=1)
                          Index Cond: (amount >= 496)
                          Buffers: shared hit=2
Planning:
  Buffers: shared hit=6
Planning Time: 0.147 ms
Execution Time: 0.130 ms

```

### s19_union_all_predicate [shape=selective/indexable;spelling=union_all]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.057000 / 0.110000 ms.

```text
Sort  (cost=144.70..145.37 rows=269 width=8) (actual time=0.091..0.097 rows=251 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 30kB
  Buffers: shared hit=114
  ->  Append  (cost=2.21..133.84 rows=269 width=8) (actual time=0.009..0.070 rows=251 loops=1)
        Buffers: shared hit=114
        ->  Bitmap Heap Scan on orders  (cost=2.21..62.83 rows=107 width=8) (actual time=0.009..0.031 rows=114 loops=1)
              Recheck Cond: (amount <= 5)
              Heap Blocks: exact=54
              Buffers: shared hit=56
              ->  Bitmap Index Scan on orders_amount_idx  (cost=0.00..2.19 rows=107 width=0) (actual time=0.003..0.003 rows=114 loops=1)
                    Index Cond: (amount <= 5)
                    Buffers: shared hit=2
        ->  Bitmap Heap Scan on orders orders_1  (cost=2.64..69.67 rows=162 width=8) (actual time=0.007..0.029 rows=137 loops=1)
              Recheck Cond: (amount >= 496)
              Heap Blocks: exact=56
              Buffers: shared hit=58
              ->  Bitmap Index Scan on orders_amount_idx  (cost=0.00..2.60 rows=162 width=0) (actual time=0.004..0.004 rows=137 loops=1)
                    Index Cond: (amount >= 496)
                    Buffers: shared hit=2
Planning:
  Buffers: shared hit=6
Planning Time: 0.057 ms
Execution Time: 0.110 ms

```

## Limits and gaps recorded by the runner

- no concurrent writers: client-decode JSON bytes are not wire bytes; tails are under-sampled at the protocol trial count
- one host: no independent hardware, real-network, 3-node placement or equal-total-budget comparison is established
- single-node only; the YB 3-node RF=3 per-node/equal-total arms are a coverage gap on this host

