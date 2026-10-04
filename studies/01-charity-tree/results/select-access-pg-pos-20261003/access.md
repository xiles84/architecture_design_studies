# SELECT access-path matrix (generated measurements)

_Measurements only. Verdicts for P10/P11 belong in a signed analysis._

**TL;DR:** within-run medians for OR / UNION ALL / UNION at three predicate shapes and OFFSET / keyset at five page depths, each cell verified against the independent oracle. A within-run difference is not a replicated result. Missing measurements are N/A.

| run | engine | topology | scale | page limit |
|---|---|---|---|---|
| select-access-pg-pos-20261003 | postgres | single | deep | 5 |

## Provenance

| Field | Value |
|---|---|
| repo commit | `6ad01ef65352b0be194007dab460b297ae94a541` |
| benchmark image | `localhost/charitybench-select:select-access-pg-pos-20261003` |
| image digest | `023c9fb2d709de026afbef825315bf2e29a234361a3d2a0f95507136c4fa57f5` |
| run tag | `run/study-01/select-access-pg-pos-20261003` |
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
| selective/indexable | `amount <= 5` | `amount >= 496` | 114 | 137 | 0 | true |
| selective/expression | `(amount + 0) <= 5` | `(amount + 0) >= 496` | 114 | 137 | 0 | true |
| selective/cast | `(amount)::bigint <= 5` | `(amount)::bigint >= 496` | 114 | 137 | 0 | true |
| broad/indexable | `amount <= 100` | `amount >= 400` | 2273 | 2379 | 0 | true |
| broad/expression | `(amount + 0) <= 100` | `(amount + 0) >= 400` | 2273 | 2379 | 0 | true |
| broad/cast | `(amount)::bigint <= 100` | `(amount)::bigint >= 400` | 2273 | 2379 | 0 | true |

## Predicate alternatives (P10)

| shape | index arm | spelling | status | median ms | spread % | trials | rows | plan structure |
|---|---|---|---|---|---|---|---|---|
| broad/cast | default | or | valid | 1.290834 | 19.5 | 5 | 4652 | `a4bab1f05f5e667d` |
| broad/cast | default | union | valid | 1.696035 | 62.1 | 5 | 4652 | `6ff77b66c9bf6ee6` |
| broad/cast | default | union_all | valid | 1.875858 | 26.7 | 5 | 4652 | `3f9f3d5c5dc15dd7` |
| broad/expression | expression_index | or | valid | 1.286213 | 50.8 | 5 | 4652 | `a4bab1f05f5e667d` |
| broad/expression | expression_index | union | valid | 1.651826 | 9.5 | 5 | 4652 | `d4fbfb3aac234d58` |
| broad/expression | expression_index | union_all | valid | 1.629259 | 12.8 | 5 | 4652 | `605ee056e65421b1` |
| broad/expression | default | or | valid | 1.355900 | 38.7 | 5 | 4652 | `a4bab1f05f5e667d` |
| broad/expression | default | union | valid | 1.633648 | 30.0 | 5 | 4652 | `6ff77b66c9bf6ee6` |
| broad/expression | default | union_all | valid | 1.925405 | 28.5 | 5 | 4652 | `3f9f3d5c5dc15dd7` |
| broad/indexable | default | or | valid | 1.834578 | 62.0 | 5 | 4652 | `a4bab1f05f5e667d` |
| broad/indexable | default | union | valid | 1.798113 | 31.1 | 5 | 4652 | `ed39b768e5b063b0` |
| broad/indexable | default | union_all | valid | 1.679676 | 10.5 | 5 | 4652 | `d51936c88987c9fd` |
| selective/cast | default | or | valid | 0.991944 | 25.4 | 5 | 251 | `a4bab1f05f5e667d` |
| selective/cast | default | union | valid | 0.851168 | 7.5 | 5 | 251 | `6ff77b66c9bf6ee6` |
| selective/cast | default | union_all | valid | 0.898037 | 12.5 | 5 | 251 | `3f9f3d5c5dc15dd7` |
| selective/expression | expression_index | or | valid | 0.264256 | 34.2 | 5 | 251 | `d66da5af4c9f3384` |
| selective/expression | expression_index | union | valid | 0.318259 | 45.4 | 5 | 251 | `a5fbbd2993065c7a` |
| selective/expression | expression_index | union_all | valid | 0.336840 | 179.0 | 5 | 251 | `eeb1f1b0c0255959` |
| selective/expression | default | or | valid | 0.863765 | 31.1 | 5 | 251 | `a4bab1f05f5e667d` |
| selective/expression | default | union | valid | 0.867208 | 12.2 | 5 | 251 | `6ff77b66c9bf6ee6` |
| selective/expression | default | union_all | valid | 0.841242 | 10.5 | 5 | 251 | `3f9f3d5c5dc15dd7` |
| selective/indexable | default | or | valid | 0.318574 | 30.9 | 5 | 251 | `8bbb087a6c122dbe` |
| selective/indexable | default | union | valid | 0.307091 | 49.6 | 5 | 251 | `567de4baf4feaaf7` |
| selective/indexable | default | union_all | valid | 0.325088 | 153.1 | 5 | 251 | `e060188a3bcf6eb3` |

### P10 pairs

Pairs compare the **default** arm only; the labelled `expression_index` arm is a package arm and is never pooled with the query-text pair.

| pair | A median ms | B median ms | B/A | same plan | within-run verdict |
|---|---|---|---|---|---|
| shape=selective/indexable or vs union_all | 0.318574 | 0.325088 | 1.0204 | false | medians inside the two cells' own spread; no resolved difference within this run |
| shape=selective/indexable or vs union | 0.318574 | 0.307091 | 0.9640 | false | medians inside the two cells' own spread; no resolved difference within this run |
| shape=selective/expression or vs union_all | 0.863765 | 0.841242 | 0.9739 | false | medians inside the two cells' own spread; no resolved difference within this run |
| shape=selective/expression or vs union | 0.863765 | 0.867208 | 1.0040 | false | medians inside the two cells' own spread; no resolved difference within this run |
| shape=selective/cast or vs union_all | 0.991944 | 0.898037 | 0.9053 | false | medians inside the two cells' own spread; no resolved difference within this run |
| shape=selective/cast or vs union | 0.991944 | 0.851168 | 0.8581 | false | medians inside the two cells' own spread; no resolved difference within this run |
| shape=broad/indexable or vs union_all | 1.834578 | 1.679676 | 0.9156 | false | medians inside the two cells' own spread; no resolved difference within this run |
| shape=broad/indexable or vs union | 1.834578 | 1.798113 | 0.9801 | false | medians inside the two cells' own spread; no resolved difference within this run |
| shape=broad/expression or vs union_all | 1.355900 | 1.925405 | 1.4200 | false | medians inside the two cells' own spread; no resolved difference within this run |
| shape=broad/expression or vs union | 1.355900 | 1.633648 | 1.2048 | false | medians inside the two cells' own spread; no resolved difference within this run |
| shape=broad/cast or vs union_all | 1.290834 | 1.875858 | 1.4532 | false | medians inside the two cells' own spread; no resolved difference within this run |
| shape=broad/cast or vs union | 1.290834 | 1.696035 | 1.3139 | false | medians inside the two cells' own spread; no resolved difference within this run |

## Pagination (P11)

| depth | spelling | status | median ms | spread % | trials | rows | plan structure |
|---|---|---|---|---|---|---|---|
| 1 | offset | valid | 0.811851 | 12.8 | 5 | 5 | `ac4509d295bf3348` |
| 1 | keyset | valid | 0.931700 | 16.1 | 5 | 5 | `ac4509d295bf3348` |
| 10 | offset | valid | 0.827168 | 14.8 | 5 | 5 | `ac4509d295bf3348` |
| 10 | keyset | valid | 0.939465 | 15.5 | 5 | 5 | `ac4509d295bf3348` |
| 100 | offset | valid | 0.902112 | 16.8 | 5 | 5 | `ac4509d295bf3348` |
| 100 | keyset | valid | 0.891731 | 7.3 | 5 | 5 | `ac4509d295bf3348` |
| 1000 | offset | valid | 1.192947 | 21.1 | 5 | 5 | `ac4509d295bf3348` |
| 1000 | keyset | valid | 0.893647 | 10.3 | 5 | 5 | `ac4509d295bf3348` |
| 10000 | offset | valid | 1.982273 | 7.0 | 5 | 5 | `ac4509d295bf3348` |
| 10000 | keyset | valid | 0.457182 | 10.7 | 5 | 5 | `dd316474b64ddc29` |

Median ratio table (keyset / OFFSET; a ratio below 1 means keyset was faster in this run):

| depth | ratio (keyset/offset) | offset median ms | keyset median ms |
|---|---|---|---|
| 1 | 1.1476 | 0.811851 | 0.931700 |
| 10 | 1.1358 | 0.827168 | 0.939465 |
| 100 | 0.9885 | 0.902112 | 0.891731 |
| 1000 | 0.7491 | 1.192947 | 0.893647 |
| 10000 | 0.2306 | 1.982273 | 0.457182 |

### P11 pairs

| pair | A median ms | B median ms | B/A | same plan | within-run verdict |
|---|---|---|---|---|---|
| depth=1 offset vs keyset | 0.811851 | 0.931700 | 1.1476 | true | medians inside the two cells' own spread; no resolved difference within this run |
| depth=10 offset vs keyset | 0.827168 | 0.939465 | 1.1358 | true | medians inside the two cells' own spread; no resolved difference within this run |
| depth=100 offset vs keyset | 0.902112 | 0.891731 | 0.9885 | true | medians inside the two cells' own spread; no resolved difference within this run |
| depth=1000 offset vs keyset | 1.192947 | 0.893647 | 0.7491 | true | medians inside the two cells' own spread; no resolved difference within this run |
| depth=10000 offset vs keyset | 1.982273 | 0.457182 | 0.2306 | false | within-run medians differ; whole-run replication required before a verdict |

## Repeated-design position control

The first matrix cell and a repeat of it at the end of the phase (`position=last`); a large ratio is a session position effect, not a spelling effect.

| pair | first median ms | last median ms | ratio (last/first) | same plan |
|---|---|---|---|---|
| position control selective/indexable or first vs last | 0.318574 | 0.281925 | 0.8850 | true |

## Expression-index arm: write and storage cost

Index `orders_amount_expr_idx` on expression `(amount + 0)`; supported=true.

| measure | value |
|---|---|
| index size (bytes) | 114688 |
| insert batch with index (ms) | 4.375290 |
| insert batch without index (ms) | 3.133509 |
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

Rows: 5; serialized decoded-row JSON bytes: 51 (not wire bytes). Planning/execution: 0.036000 / 0.330000 ms.

```text
Limit  (cost=154.12..154.13 rows=5 width=8) (actual time=0.321..0.322 rows=5 loops=1)
  Buffers: shared hit=73
  ->  Sort  (cost=154.12..158.90 rows=1912 width=8) (actual time=0.321..0.322 rows=5 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 25kB
        Buffers: shared hit=73
        ->  Bitmap Heap Scan on orders  (cost=22.66..122.36 rows=1912 width=8) (actual time=0.035..0.215 rows=2000 loops=1)
              Recheck Cond: ((amount < 67) OR (amount = 67) OR (amount IS NULL))
              Filter: ((amount < 67) OR ((amount = 67) AND (id > 4310)) OR (amount IS NULL))
              Rows Removed by Filter: 6
              Heap Blocks: exact=65
              Buffers: shared hit=73
              ->  BitmapOr  (cost=22.66..22.66 rows=1983 width=0) (actual time=0.027..0.027 rows=0 loops=1)
                    Buffers: shared hit=8
                    ->  Bitmap Index Scan on orders_amount_idx  (cost=0.00..14.38 rows=1440 width=0) (actual time=0.019..0.019 rows=1460 loops=1)
                          Index Cond: (amount < 67)
                          Buffers: shared hit=4
                    ->  Bitmap Index Scan on orders_amount_idx  (cost=0.00..1.54 rows=21 width=0) (actual time=0.001..0.001 rows=24 loops=1)
                          Index Cond: (amount = 67)
                          Buffers: shared hit=2
                    ->  Bitmap Index Scan on orders_amount_idx  (cost=0.00..5.30 rows=522 width=0) (actual time=0.006..0.006 rows=522 loops=1)
                          Index Cond: (amount IS NULL)
                          Buffers: shared hit=2
Planning Time: 0.036 ms
Execution Time: 0.330 ms

```

### s16_page_offset [depth=10000;spelling=offset]

Rows: 5; serialized decoded-row JSON bytes: 51 (not wire bytes). Planning/execution: 0.039000 / 2.730000 ms.

```text
Limit  (cost=1023.04..1023.06 rows=5 width=8) (actual time=2.696..2.697 rows=5 loops=1)
  Buffers: shared hit=65
  ->  Sort  (cost=998.04..1028.04 rows=12000 width=8) (actual time=2.254..2.495 rows=10005 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: quicksort  Memory: 853kB
        Buffers: shared hit=65
        ->  Seq Scan on orders  (cost=0.00..185.00 rows=12000 width=8) (actual time=0.006..0.485 rows=12000 loops=1)
              Buffers: shared hit=65
Planning Time: 0.039 ms
Execution Time: 2.730 ms

```

### s17_page_keyset [depth=1000;spelling=keyset]

Rows: 5; serialized decoded-row JSON bytes: 59 (not wire bytes). Planning/execution: 0.040000 / 1.145000 ms.

```text
Limit  (cost=450.03..450.04 rows=5 width=8) (actual time=1.137..1.138 rows=5 loops=1)
  Buffers: shared hit=65
  ->  Sort  (cost=450.03..476.38 rows=10538 width=8) (actual time=1.137..1.137 rows=5 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 25kB
        Buffers: shared hit=65
        ->  Seq Scan on orders  (cost=0.00..275.00 rows=10538 width=8) (actual time=0.005..0.583 rows=11000 loops=1)
              Filter: ((amount < 457) OR ((amount = 457) AND (id > 11254)) OR (amount IS NULL))
              Rows Removed by Filter: 1000
              Buffers: shared hit=65
Planning Time: 0.040 ms
Execution Time: 1.145 ms

```

### s16_page_offset [depth=1000;spelling=offset]

Rows: 5; serialized decoded-row JSON bytes: 59 (not wire bytes). Planning/execution: 0.020000 / 1.395000 ms.

```text
Limit  (cost=845.88..845.89 rows=5 width=8) (actual time=1.389..1.390 rows=5 loops=1)
  Buffers: shared hit=65
  ->  Sort  (cost=843.38..873.38 rows=12000 width=8) (actual time=1.350..1.371 rows=1005 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 115kB
        Buffers: shared hit=65
        ->  Seq Scan on orders  (cost=0.00..185.00 rows=12000 width=8) (actual time=0.004..0.444 rows=12000 loops=1)
              Buffers: shared hit=65
Planning Time: 0.020 ms
Execution Time: 1.395 ms

```

### s17_page_keyset [depth=100;spelling=keyset]

Rows: 5; serialized decoded-row JSON bytes: 56 (not wire bytes). Planning/execution: 0.046000 / 1.158000 ms.

```text
Limit  (cost=463.65..463.66 rows=5 width=8) (actual time=1.152..1.153 rows=5 loops=1)
  Buffers: shared hit=65
  ->  Sort  (cost=463.65..492.05 rows=11358 width=8) (actual time=1.151..1.152 rows=5 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 25kB
        Buffers: shared hit=65
        ->  Seq Scan on orders  (cost=0.00..275.00 rows=11358 width=8) (actual time=0.004..0.568 rows=11900 loops=1)
              Filter: ((amount < 497) OR ((amount = 497) AND (id > 6723)) OR (amount IS NULL))
              Rows Removed by Filter: 100
              Buffers: shared hit=65
Planning:
  Buffers: shared hit=3
Planning Time: 0.046 ms
Execution Time: 1.158 ms

```

### s16_page_offset [depth=100;spelling=offset]

Rows: 5; serialized decoded-row JSON bytes: 56 (not wire bytes). Planning/execution: 0.016000 / 1.108000 ms.

```text
Limit  (cost=648.10..648.12 rows=5 width=8) (actual time=1.103..1.104 rows=5 loops=1)
  Buffers: shared hit=65
  ->  Sort  (cost=647.85..677.85 rows=12000 width=8) (actual time=1.098..1.101 rows=105 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 31kB
        Buffers: shared hit=65
        ->  Seq Scan on orders  (cost=0.00..185.00 rows=12000 width=8) (actual time=0.003..0.463 rows=12000 loops=1)
              Buffers: shared hit=65
Planning Time: 0.016 ms
Execution Time: 1.108 ms

```

### s17_page_keyset [depth=10;spelling=keyset]

Rows: 5; serialized decoded-row JSON bytes: 56 (not wire bytes). Planning/execution: 0.045000 / 1.119000 ms.

```text
Limit  (cost=465.21..465.23 rows=5 width=8) (actual time=1.113..1.114 rows=5 loops=1)
  Buffers: shared hit=65
  ->  Sort  (cost=465.21..493.84 rows=11452 width=8) (actual time=1.113..1.113 rows=5 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 25kB
        Buffers: shared hit=65
        ->  Seq Scan on orders  (cost=0.00..275.00 rows=11452 width=8) (actual time=0.004..0.541 rows=11990 loops=1)
              Filter: ((amount < 500) OR ((amount = 500) AND (id > 6319)) OR (amount IS NULL))
              Rows Removed by Filter: 10
              Buffers: shared hit=65
Planning:
  Buffers: shared hit=3
Planning Time: 0.045 ms
Execution Time: 1.119 ms

```

### s16_page_offset [depth=10;spelling=offset]

Rows: 5; serialized decoded-row JSON bytes: 56 (not wire bytes). Planning/execution: 0.035000 / 1.118000 ms.

```text
Limit  (cost=479.44..479.45 rows=5 width=8) (actual time=1.111..1.112 rows=5 loops=1)
  Buffers: shared hit=65
  ->  Sort  (cost=479.41..509.41 rows=12000 width=8) (actual time=1.111..1.111 rows=15 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 25kB
        Buffers: shared hit=65
        ->  Seq Scan on orders  (cost=0.00..185.00 rows=12000 width=8) (actual time=0.005..0.497 rows=12000 loops=1)
              Buffers: shared hit=65
Planning Time: 0.035 ms
Execution Time: 1.118 ms

```

### s17_page_keyset [depth=1;spelling=keyset]

Rows: 5; serialized decoded-row JSON bytes: 56 (not wire bytes). Planning/execution: 0.053000 / 1.182000 ms.

```text
Limit  (cost=465.21..465.23 rows=5 width=8) (actual time=1.172..1.174 rows=5 loops=1)
  Buffers: shared hit=65
  ->  Sort  (cost=465.21..493.84 rows=11452 width=8) (actual time=1.172..1.172 rows=5 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 25kB
        Buffers: shared hit=65
        ->  Seq Scan on orders  (cost=0.00..275.00 rows=11452 width=8) (actual time=0.005..0.573 rows=11999 loops=1)
              Filter: ((amount < 500) OR ((amount = 500) AND (id > 1554)) OR (amount IS NULL))
              Rows Removed by Filter: 1
              Buffers: shared hit=65
Planning:
  Buffers: shared hit=3
Planning Time: 0.053 ms
Execution Time: 1.182 ms

```

### s16_page_offset [depth=1;spelling=offset]

Rows: 5; serialized decoded-row JSON bytes: 56 (not wire bytes). Planning/execution: 0.023000 / 1.113000 ms.

```text
Limit  (cost=400.10..400.11 rows=5 width=8) (actual time=1.106..1.107 rows=5 loops=1)
  Buffers: shared hit=65
  ->  Sort  (cost=400.10..430.10 rows=12000 width=8) (actual time=1.105..1.105 rows=6 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 25kB
        Buffers: shared hit=65
        ->  Seq Scan on orders  (cost=0.00..185.00 rows=12000 width=8) (actual time=0.004..0.477 rows=12000 loops=1)
              Buffers: shared hit=65
Planning Time: 0.023 ms
Execution Time: 1.113 ms

```

### s18_or_predicate [shape=broad/cast;spelling=or]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.031000 / 0.817000 ms.

```text
Index Scan using orders_pkey on orders  (cost=0.29..403.89 rows=6667 width=8) (actual time=0.009..0.721 rows=4652 loops=1)
  Filter: (((amount)::bigint <= 100) OR ((amount)::bigint >= 400))
  Rows Removed by Filter: 7348
  Buffers: shared hit=99
Planning Time: 0.031 ms
Execution Time: 0.817 ms

```

### s20_union_predicate [shape=broad/cast;spelling=union]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.051000 / 1.768000 ms.

```text
Unique  (cost=968.64..1108.64 rows=8000 width=8) (actual time=1.056..1.664 rows=4652 loops=1)
  Buffers: shared hit=130
  ->  Merge Append  (cost=968.64..1068.64 rows=8000 width=8) (actual time=1.055..1.401 rows=4652 loops=1)
        Sort Key: orders.id, orders.amount
        Buffers: shared hit=130
        ->  Sort  (cost=484.32..494.32 rows=4000 width=8) (actual time=0.567..0.622 rows=2273 loops=1)
              Sort Key: orders.id, orders.amount
              Sort Method: quicksort  Memory: 150kB
              Buffers: shared hit=65
              ->  Seq Scan on orders  (cost=0.00..245.00 rows=4000 width=8) (actual time=0.010..0.424 rows=2273 loops=1)
                    Filter: ((amount)::bigint <= 100)
                    Rows Removed by Filter: 9727
                    Buffers: shared hit=65
        ->  Sort  (cost=484.32..494.32 rows=4000 width=8) (actual time=0.487..0.536 rows=2379 loops=1)
              Sort Key: orders_1.id, orders_1.amount
              Sort Method: quicksort  Memory: 152kB
              Buffers: shared hit=65
              ->  Seq Scan on orders orders_1  (cost=0.00..245.00 rows=4000 width=8) (actual time=0.008..0.368 rows=2379 loops=1)
                    Filter: ((amount)::bigint >= 400)
                    Rows Removed by Filter: 9621
                    Buffers: shared hit=65
Planning Time: 0.051 ms
Execution Time: 1.768 ms

```

### s19_union_all_predicate [shape=broad/cast;spelling=union_all]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.047000 / 1.537000 ms.

```text
Sort  (cost=1048.63..1068.63 rows=8000 width=8) (actual time=1.320..1.433 rows=4652 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 302kB
  Buffers: shared hit=130
  ->  Append  (cost=0.00..530.00 rows=8000 width=8) (actual time=0.008..0.940 rows=4652 loops=1)
        Buffers: shared hit=130
        ->  Seq Scan on orders  (cost=0.00..245.00 rows=4000 width=8) (actual time=0.008..0.398 rows=2273 loops=1)
              Filter: ((amount)::bigint <= 100)
              Rows Removed by Filter: 9727
              Buffers: shared hit=65
        ->  Seq Scan on orders orders_1  (cost=0.00..245.00 rows=4000 width=8) (actual time=0.002..0.358 rows=2379 loops=1)
              Filter: ((amount)::bigint >= 400)
              Rows Removed by Filter: 9621
              Buffers: shared hit=65
Planning Time: 0.047 ms
Execution Time: 1.537 ms

```

### s18_or_predicate [shape=broad/expression;index=expression_index;spelling=or]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.032000 / 0.839000 ms.

```text
Index Scan using orders_pkey on orders  (cost=0.29..403.89 rows=4211 width=8) (actual time=0.008..0.738 rows=4652 loops=1)
  Filter: (((amount + 0) <= 100) OR ((amount + 0) >= 400))
  Rows Removed by Filter: 7348
  Buffers: shared hit=99
Planning Time: 0.032 ms
Execution Time: 0.839 ms

```

### s20_union_predicate [shape=broad/expression;index=expression_index;spelling=union]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.063000 / 1.697000 ms.

```text
Unique  (cost=492.69..574.31 rows=4664 width=8) (actual time=0.953..1.593 rows=4652 loops=1)
  Buffers: shared hit=3934
  ->  Merge Append  (cost=492.69..550.99 rows=4664 width=8) (actual time=0.952..1.310 rows=4652 loops=1)
        Sort Key: orders.id, orders.amount
        Buffers: shared hit=3934
        ->  Sort  (cost=241.10..246.77 rows=2267 width=8) (actual time=0.482..0.534 rows=2273 loops=1)
              Sort Key: orders.id, orders.amount
              Sort Method: quicksort  Memory: 150kB
              Buffers: shared hit=1921
              ->  Index Scan using orders_amount_expr_idx on orders  (cost=0.29..114.75 rows=2267 width=8) (actual time=0.007..0.273 rows=2273 loops=1)
                    Index Cond: ((amount + 0) <= 100)
                    Buffers: shared hit=1921
        ->  Sort  (cost=251.58..257.57 rows=2397 width=8) (actual time=0.469..0.522 rows=2379 loops=1)
              Sort Key: orders_1.id, orders_1.amount
              Sort Method: quicksort  Memory: 152kB
              Buffers: shared hit=2013
              ->  Index Scan using orders_amount_expr_idx on orders orders_1  (cost=0.29..117.03 rows=2397 width=8) (actual time=0.003..0.255 rows=2379 loops=1)
                    Index Cond: ((amount + 0) >= 400)
                    Buffers: shared hit=2013
Planning Time: 0.063 ms
Execution Time: 1.697 ms

```

### s19_union_all_predicate [shape=broad/expression;index=expression_index;spelling=union_all]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.047000 / 1.324000 ms.

```text
Sort  (cost=539.31..550.97 rows=4664 width=8) (actual time=1.122..1.228 rows=4652 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 302kB
  Buffers: shared hit=3934
  ->  Append  (cost=0.29..255.10 rows=4664 width=8) (actual time=0.009..0.690 rows=4652 loops=1)
        Buffers: shared hit=3934
        ->  Index Scan using orders_amount_expr_idx on orders  (cost=0.29..114.75 rows=2267 width=8) (actual time=0.008..0.269 rows=2273 loops=1)
              Index Cond: ((amount + 0) <= 100)
              Buffers: shared hit=1921
        ->  Index Scan using orders_amount_expr_idx on orders orders_1  (cost=0.29..117.03 rows=2397 width=8) (actual time=0.003..0.251 rows=2379 loops=1)
              Index Cond: ((amount + 0) >= 400)
              Buffers: shared hit=2013
Planning Time: 0.047 ms
Execution Time: 1.324 ms

```

### s18_or_predicate [shape=broad/expression;spelling=or]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.038000 / 0.952000 ms.

```text
Index Scan using orders_pkey on orders  (cost=0.29..403.89 rows=6667 width=8) (actual time=0.010..0.842 rows=4652 loops=1)
  Filter: (((amount + 0) <= 100) OR ((amount + 0) >= 400))
  Rows Removed by Filter: 7348
  Buffers: shared hit=99
Planning Time: 0.038 ms
Execution Time: 0.952 ms

```

### s20_union_predicate [shape=broad/expression;spelling=union]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.050000 / 1.718000 ms.

```text
Unique  (cost=968.64..1108.64 rows=8000 width=8) (actual time=0.978..1.615 rows=4652 loops=1)
  Buffers: shared hit=130
  ->  Merge Append  (cost=968.64..1068.64 rows=8000 width=8) (actual time=0.978..1.351 rows=4652 loops=1)
        Sort Key: orders.id, orders.amount
        Buffers: shared hit=130
        ->  Sort  (cost=484.32..494.32 rows=4000 width=8) (actual time=0.513..0.577 rows=2273 loops=1)
              Sort Key: orders.id, orders.amount
              Sort Method: quicksort  Memory: 150kB
              Buffers: shared hit=65
              ->  Seq Scan on orders  (cost=0.00..245.00 rows=4000 width=8) (actual time=0.006..0.389 rows=2273 loops=1)
                    Filter: ((amount + 0) <= 100)
                    Rows Removed by Filter: 9727
                    Buffers: shared hit=65
        ->  Sort  (cost=484.32..494.32 rows=4000 width=8) (actual time=0.464..0.515 rows=2379 loops=1)
              Sort Key: orders_1.id, orders_1.amount
              Sort Method: quicksort  Memory: 152kB
              Buffers: shared hit=65
              ->  Seq Scan on orders orders_1  (cost=0.00..245.00 rows=4000 width=8) (actual time=0.003..0.343 rows=2379 loops=1)
                    Filter: ((amount + 0) >= 400)
                    Rows Removed by Filter: 9621
                    Buffers: shared hit=65
Planning Time: 0.050 ms
Execution Time: 1.718 ms

```

### s19_union_all_predicate [shape=broad/expression;spelling=union_all]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.045000 / 1.567000 ms.

```text
Sort  (cost=1048.63..1068.63 rows=8000 width=8) (actual time=1.343..1.458 rows=4652 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 302kB
  Buffers: shared hit=130
  ->  Append  (cost=0.00..530.00 rows=8000 width=8) (actual time=0.008..0.928 rows=4652 loops=1)
        Buffers: shared hit=130
        ->  Seq Scan on orders  (cost=0.00..245.00 rows=4000 width=8) (actual time=0.007..0.373 rows=2273 loops=1)
              Filter: ((amount + 0) <= 100)
              Rows Removed by Filter: 9727
              Buffers: shared hit=65
        ->  Seq Scan on orders orders_1  (cost=0.00..245.00 rows=4000 width=8) (actual time=0.003..0.363 rows=2379 loops=1)
              Filter: ((amount + 0) >= 400)
              Rows Removed by Filter: 9621
              Buffers: shared hit=65
Planning Time: 0.045 ms
Execution Time: 1.567 ms

```

### s18_or_predicate [shape=broad/indexable;spelling=or]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.052000 / 0.843000 ms.

```text
Index Scan using orders_pkey on orders  (cost=0.29..343.89 rows=4211 width=8) (actual time=0.010..0.736 rows=4652 loops=1)
  Filter: ((amount <= 100) OR (amount >= 400))
  Rows Removed by Filter: 7348
  Buffers: shared hit=99
Planning Time: 0.052 ms
Execution Time: 0.843 ms

```

### s20_union_predicate [shape=broad/indexable;spelling=union]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.053000 / 1.869000 ms.

```text
Unique  (cost=495.99..577.61 rows=4664 width=8) (actual time=1.052..1.753 rows=4652 loops=1)
  Buffers: shared hit=3936
  ->  Merge Append  (cost=495.99..554.29 rows=4664 width=8) (actual time=1.052..1.428 rows=4652 loops=1)
        Sort Key: orders.id, orders.amount
        Buffers: shared hit=3936
        ->  Sort  (cost=242.20..247.87 rows=2267 width=8) (actual time=0.558..0.620 rows=2273 loops=1)
              Sort Key: orders.id, orders.amount
              Sort Method: quicksort  Memory: 150kB
              Buffers: shared hit=1922
              ->  Index Scan using orders_amount_idx on orders  (cost=0.29..115.85 rows=2267 width=8) (actual time=0.007..0.336 rows=2273 loops=1)
                    Index Cond: (amount <= 100)
                    Buffers: shared hit=1922
        ->  Sort  (cost=253.78..259.77 rows=2397 width=8) (actual time=0.493..0.555 rows=2379 loops=1)
              Sort Key: orders_1.id, orders_1.amount
              Sort Method: quicksort  Memory: 152kB
              Buffers: shared hit=2014
              ->  Index Scan using orders_amount_idx on orders orders_1  (cost=0.29..119.23 rows=2397 width=8) (actual time=0.004..0.262 rows=2379 loops=1)
                    Index Cond: (amount >= 400)
                    Buffers: shared hit=2014
Planning Time: 0.053 ms
Execution Time: 1.869 ms

```

### s19_union_all_predicate [shape=broad/indexable;spelling=union_all]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.076000 / 1.412000 ms.

```text
Sort  (cost=542.61..554.27 rows=4664 width=8) (actual time=1.207..1.311 rows=4652 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 302kB
  Buffers: shared hit=3936
  ->  Append  (cost=0.29..258.40 rows=4664 width=8) (actual time=0.010..0.771 rows=4652 loops=1)
        Buffers: shared hit=3936
        ->  Index Scan using orders_amount_idx on orders  (cost=0.29..115.85 rows=2267 width=8) (actual time=0.009..0.323 rows=2273 loops=1)
              Index Cond: (amount <= 100)
              Buffers: shared hit=1922
        ->  Index Scan using orders_amount_idx on orders orders_1  (cost=0.29..119.23 rows=2397 width=8) (actual time=0.005..0.258 rows=2379 loops=1)
              Index Cond: (amount >= 400)
              Buffers: shared hit=2014
Planning Time: 0.076 ms
Execution Time: 1.412 ms

```

### s18_or_predicate [shape=selective/cast;spelling=or]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.029000 / 0.772000 ms.

```text
Index Scan using orders_pkey on orders  (cost=0.29..403.89 rows=6667 width=8) (actual time=0.008..0.758 rows=251 loops=1)
  Filter: (((amount)::bigint <= 5) OR ((amount)::bigint >= 496))
  Rows Removed by Filter: 11749
  Buffers: shared hit=99
Planning Time: 0.029 ms
Execution Time: 0.772 ms

```

### s20_union_predicate [shape=selective/cast;spelling=union]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.042000 / 0.659000 ms.

```text
Unique  (cost=968.64..1108.64 rows=8000 width=8) (actual time=0.589..0.644 rows=251 loops=1)
  Buffers: shared hit=130
  ->  Merge Append  (cost=968.64..1068.64 rows=8000 width=8) (actual time=0.588..0.626 rows=251 loops=1)
        Sort Key: orders.id, orders.amount
        Buffers: shared hit=130
        ->  Sort  (cost=484.32..494.32 rows=4000 width=8) (actual time=0.298..0.301 rows=114 loops=1)
              Sort Key: orders.id, orders.amount
              Sort Method: quicksort  Memory: 27kB
              Buffers: shared hit=65
              ->  Seq Scan on orders  (cost=0.00..245.00 rows=4000 width=8) (actual time=0.006..0.290 rows=114 loops=1)
                    Filter: ((amount)::bigint <= 5)
                    Rows Removed by Filter: 11886
                    Buffers: shared hit=65
        ->  Sort  (cost=484.32..494.32 rows=4000 width=8) (actual time=0.290..0.310 rows=137 loops=1)
              Sort Key: orders_1.id, orders_1.amount
              Sort Method: quicksort  Memory: 28kB
              Buffers: shared hit=65
              ->  Seq Scan on orders orders_1  (cost=0.00..245.00 rows=4000 width=8) (actual time=0.004..0.282 rows=137 loops=1)
                    Filter: ((amount)::bigint >= 496)
                    Rows Removed by Filter: 11863
                    Buffers: shared hit=65
Planning Time: 0.042 ms
Execution Time: 0.659 ms

```

### s19_union_all_predicate [shape=selective/cast;spelling=union_all]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.037000 / 0.662000 ms.

```text
Sort  (cost=1048.63..1068.63 rows=8000 width=8) (actual time=0.643..0.649 rows=251 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 30kB
  Buffers: shared hit=130
  ->  Append  (cost=0.00..530.00 rows=8000 width=8) (actual time=0.007..0.620 rows=251 loops=1)
        Buffers: shared hit=130
        ->  Seq Scan on orders  (cost=0.00..245.00 rows=4000 width=8) (actual time=0.006..0.321 rows=114 loops=1)
              Filter: ((amount)::bigint <= 5)
              Rows Removed by Filter: 11886
              Buffers: shared hit=65
        ->  Seq Scan on orders orders_1  (cost=0.00..245.00 rows=4000 width=8) (actual time=0.004..0.288 rows=137 loops=1)
              Filter: ((amount)::bigint >= 496)
              Rows Removed by Filter: 11863
              Buffers: shared hit=65
Planning Time: 0.037 ms
Execution Time: 0.662 ms

```

### s18_or_predicate [shape=selective/expression;index=expression_index;spelling=or]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.041000 / 0.080000 ms.

```text
Sort  (cost=86.11..86.78 rows=268 width=8) (actual time=0.062..0.068 rows=251 loops=1)
  Sort Key: id
  Sort Method: quicksort  Memory: 30kB
  Buffers: shared hit=67
  ->  Bitmap Heap Scan on orders  (cost=4.92..75.30 rows=268 width=8) (actual time=0.013..0.048 rows=251 loops=1)
        Recheck Cond: (((amount + 0) <= 5) OR ((amount + 0) >= 496))
        Heap Blocks: exact=63
        Buffers: shared hit=67
        ->  BitmapOr  (cost=4.92..4.92 rows=269 width=0) (actual time=0.005..0.005 rows=0 loops=1)
              Buffers: shared hit=4
              ->  Bitmap Index Scan on orders_amount_expr_idx  (cost=0.00..2.19 rows=107 width=0) (actual time=0.003..0.003 rows=114 loops=1)
                    Index Cond: ((amount + 0) <= 5)
                    Buffers: shared hit=2
              ->  Bitmap Index Scan on orders_amount_expr_idx  (cost=0.00..2.60 rows=162 width=0) (actual time=0.002..0.002 rows=137 loops=1)
                    Index Cond: ((amount + 0) >= 496)
                    Buffers: shared hit=2
Planning:
  Buffers: shared hit=6
Planning Time: 0.041 ms
Execution Time: 0.080 ms

```

### s20_union_predicate [shape=selective/expression;index=expression_index;spelling=union]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.052000 / 0.122000 ms.

```text
Unique  (cost=142.73..147.44 rows=269 width=8) (actual time=0.072..0.109 rows=251 loops=1)
  Buffers: shared hit=114
  ->  Merge Append  (cost=142.73..146.09 rows=269 width=8) (actual time=0.071..0.092 rows=251 loops=1)
        Sort Key: orders.id, orders.amount
        Buffers: shared hit=114
        ->  Sort  (cost=66.70..66.97 rows=107 width=8) (actual time=0.034..0.037 rows=114 loops=1)
              Sort Key: orders.id, orders.amount
              Sort Method: quicksort  Memory: 27kB
              Buffers: shared hit=56
              ->  Bitmap Heap Scan on orders  (cost=2.21..63.10 rows=107 width=8) (actual time=0.008..0.027 rows=114 loops=1)
                    Recheck Cond: ((amount + 0) <= 5)
                    Heap Blocks: exact=54
                    Buffers: shared hit=56
                    ->  Bitmap Index Scan on orders_amount_expr_idx  (cost=0.00..2.19 rows=107 width=0) (actual time=0.003..0.003 rows=114 loops=1)
                          Index Cond: ((amount + 0) <= 5)
                          Buffers: shared hit=2
        ->  Sort  (cost=76.02..76.42 rows=162 width=8) (actual time=0.037..0.041 rows=137 loops=1)
              Sort Key: orders_1.id, orders_1.amount
              Sort Method: quicksort  Memory: 28kB
              Buffers: shared hit=58
              ->  Bitmap Heap Scan on orders orders_1  (cost=2.64..70.07 rows=162 width=8) (actual time=0.007..0.028 rows=137 loops=1)
                    Recheck Cond: ((amount + 0) >= 496)
                    Heap Blocks: exact=56
                    Buffers: shared hit=58
                    ->  Bitmap Index Scan on orders_amount_expr_idx  (cost=0.00..2.60 rows=162 width=0) (actual time=0.004..0.004 rows=137 loops=1)
                          Index Cond: ((amount + 0) >= 496)
                          Buffers: shared hit=2
Planning:
  Buffers: shared hit=6
Planning Time: 0.052 ms
Execution Time: 0.122 ms

```

### s19_union_all_predicate [shape=selective/expression;index=expression_index;spelling=union_all]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.049000 / 0.133000 ms.

```text
Sort  (cost=145.37..146.04 rows=269 width=8) (actual time=0.113..0.120 rows=251 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 30kB
  Buffers: shared hit=114
  ->  Append  (cost=2.21..134.51 rows=269 width=8) (actual time=0.008..0.092 rows=251 loops=1)
        Buffers: shared hit=114
        ->  Bitmap Heap Scan on orders  (cost=2.21..63.10 rows=107 width=8) (actual time=0.008..0.031 rows=114 loops=1)
              Recheck Cond: ((amount + 0) <= 5)
              Heap Blocks: exact=54
              Buffers: shared hit=56
              ->  Bitmap Index Scan on orders_amount_expr_idx  (cost=0.00..2.19 rows=107 width=0) (actual time=0.003..0.003 rows=114 loops=1)
                    Index Cond: ((amount + 0) <= 5)
                    Buffers: shared hit=2
        ->  Bitmap Heap Scan on orders orders_1  (cost=2.64..70.07 rows=162 width=8) (actual time=0.030..0.050 rows=137 loops=1)
              Recheck Cond: ((amount + 0) >= 496)
              Heap Blocks: exact=56
              Buffers: shared hit=58
              ->  Bitmap Index Scan on orders_amount_expr_idx  (cost=0.00..2.60 rows=162 width=0) (actual time=0.025..0.025 rows=137 loops=1)
                    Index Cond: ((amount + 0) >= 496)
                    Buffers: shared hit=2
Planning:
  Buffers: shared hit=6
Planning Time: 0.049 ms
Execution Time: 0.133 ms

```

### s18_or_predicate [shape=selective/expression;spelling=or]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.022000 / 0.702000 ms.

```text
Index Scan using orders_pkey on orders  (cost=0.29..403.89 rows=6667 width=8) (actual time=0.008..0.691 rows=251 loops=1)
  Filter: (((amount + 0) <= 5) OR ((amount + 0) >= 496))
  Rows Removed by Filter: 11749
  Buffers: shared hit=99
Planning Time: 0.022 ms
Execution Time: 0.702 ms

```

### s20_union_predicate [shape=selective/expression;spelling=union]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.037000 / 0.655000 ms.

```text
Unique  (cost=968.64..1108.64 rows=8000 width=8) (actual time=0.600..0.641 rows=251 loops=1)
  Buffers: shared hit=130
  ->  Merge Append  (cost=968.64..1068.64 rows=8000 width=8) (actual time=0.599..0.621 rows=251 loops=1)
        Sort Key: orders.id, orders.amount
        Buffers: shared hit=130
        ->  Sort  (cost=484.32..494.32 rows=4000 width=8) (actual time=0.290..0.294 rows=114 loops=1)
              Sort Key: orders.id, orders.amount
              Sort Method: quicksort  Memory: 27kB
              Buffers: shared hit=65
              ->  Seq Scan on orders  (cost=0.00..245.00 rows=4000 width=8) (actual time=0.005..0.283 rows=114 loops=1)
                    Filter: ((amount + 0) <= 5)
                    Rows Removed by Filter: 11886
                    Buffers: shared hit=65
        ->  Sort  (cost=484.32..494.32 rows=4000 width=8) (actual time=0.308..0.311 rows=137 loops=1)
              Sort Key: orders_1.id, orders_1.amount
              Sort Method: quicksort  Memory: 28kB
              Buffers: shared hit=65
              ->  Seq Scan on orders orders_1  (cost=0.00..245.00 rows=4000 width=8) (actual time=0.004..0.299 rows=137 loops=1)
                    Filter: ((amount + 0) >= 496)
                    Rows Removed by Filter: 11863
                    Buffers: shared hit=65
Planning Time: 0.037 ms
Execution Time: 0.655 ms

```

### s19_union_all_predicate [shape=selective/expression;spelling=union_all]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.028000 / 0.570000 ms.

```text
Sort  (cost=1048.63..1068.63 rows=8000 width=8) (actual time=0.554..0.560 rows=251 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 30kB
  Buffers: shared hit=130
  ->  Append  (cost=0.00..530.00 rows=8000 width=8) (actual time=0.005..0.536 rows=251 loops=1)
        Buffers: shared hit=130
        ->  Seq Scan on orders  (cost=0.00..245.00 rows=4000 width=8) (actual time=0.004..0.271 rows=114 loops=1)
              Filter: ((amount + 0) <= 5)
              Rows Removed by Filter: 11886
              Buffers: shared hit=65
        ->  Seq Scan on orders orders_1  (cost=0.00..245.00 rows=4000 width=8) (actual time=0.003..0.254 rows=137 loops=1)
              Filter: ((amount + 0) >= 496)
              Rows Removed by Filter: 11863
              Buffers: shared hit=65
Planning Time: 0.028 ms
Execution Time: 0.570 ms

```

### s18_or_predicate [shape=selective/indexable;spelling=or]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.039000 / 0.088000 ms.

```text
Sort  (cost=84.77..85.44 rows=268 width=8) (actual time=0.068..0.075 rows=251 loops=1)
  Sort Key: id
  Sort Method: quicksort  Memory: 30kB
  Buffers: shared hit=67
  ->  Bitmap Heap Scan on orders  (cost=4.92..73.96 rows=268 width=8) (actual time=0.014..0.052 rows=251 loops=1)
        Recheck Cond: ((amount <= 5) OR (amount >= 496))
        Heap Blocks: exact=63
        Buffers: shared hit=67
        ->  BitmapOr  (cost=4.92..4.92 rows=269 width=0) (actual time=0.007..0.008 rows=0 loops=1)
              Buffers: shared hit=4
              ->  Bitmap Index Scan on orders_amount_idx  (cost=0.00..2.19 rows=107 width=0) (actual time=0.004..0.004 rows=114 loops=1)
                    Index Cond: (amount <= 5)
                    Buffers: shared hit=2
              ->  Bitmap Index Scan on orders_amount_idx  (cost=0.00..2.60 rows=162 width=0) (actual time=0.003..0.003 rows=137 loops=1)
                    Index Cond: (amount >= 496)
                    Buffers: shared hit=2
Planning:
  Buffers: shared hit=6
Planning Time: 0.039 ms
Execution Time: 0.088 ms

```

### s18_or_predicate [shape=selective/indexable;spelling=or;position=last]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.043000 / 0.090000 ms.

```text
Sort  (cost=84.77..85.44 rows=268 width=8) (actual time=0.069..0.076 rows=251 loops=1)
  Sort Key: id
  Sort Method: quicksort  Memory: 30kB
  Buffers: shared hit=67
  ->  Bitmap Heap Scan on orders  (cost=4.92..73.96 rows=268 width=8) (actual time=0.014..0.053 rows=251 loops=1)
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
Planning Time: 0.043 ms
Execution Time: 0.090 ms

```

### s20_union_predicate [shape=selective/indexable;spelling=union]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.050000 / 0.127000 ms.

```text
Unique  (cost=142.06..146.77 rows=269 width=8) (actual time=0.074..0.111 rows=251 loops=1)
  Buffers: shared hit=114
  ->  Merge Append  (cost=142.06..145.42 rows=269 width=8) (actual time=0.074..0.095 rows=251 loops=1)
        Sort Key: orders.id, orders.amount
        Buffers: shared hit=114
        ->  Sort  (cost=66.44..66.70 rows=107 width=8) (actual time=0.038..0.041 rows=114 loops=1)
              Sort Key: orders.id, orders.amount
              Sort Method: quicksort  Memory: 27kB
              Buffers: shared hit=56
              ->  Bitmap Heap Scan on orders  (cost=2.21..62.83 rows=107 width=8) (actual time=0.008..0.030 rows=114 loops=1)
                    Recheck Cond: (amount <= 5)
                    Heap Blocks: exact=54
                    Buffers: shared hit=56
                    ->  Bitmap Index Scan on orders_amount_idx  (cost=0.00..2.19 rows=107 width=0) (actual time=0.003..0.003 rows=114 loops=1)
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
                    ->  Bitmap Index Scan on orders_amount_idx  (cost=0.00..2.60 rows=162 width=0) (actual time=0.003..0.003 rows=137 loops=1)
                          Index Cond: (amount >= 496)
                          Buffers: shared hit=2
Planning:
  Buffers: shared hit=6
Planning Time: 0.050 ms
Execution Time: 0.127 ms

```

### s19_union_all_predicate [shape=selective/indexable;spelling=union_all]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.049000 / 0.109000 ms.

```text
Sort  (cost=144.70..145.37 rows=269 width=8) (actual time=0.090..0.097 rows=251 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 30kB
  Buffers: shared hit=114
  ->  Append  (cost=2.21..133.84 rows=269 width=8) (actual time=0.009..0.069 rows=251 loops=1)
        Buffers: shared hit=114
        ->  Bitmap Heap Scan on orders  (cost=2.21..62.83 rows=107 width=8) (actual time=0.009..0.030 rows=114 loops=1)
              Recheck Cond: (amount <= 5)
              Heap Blocks: exact=54
              Buffers: shared hit=56
              ->  Bitmap Index Scan on orders_amount_idx  (cost=0.00..2.19 rows=107 width=0) (actual time=0.004..0.004 rows=114 loops=1)
                    Index Cond: (amount <= 5)
                    Buffers: shared hit=2
        ->  Bitmap Heap Scan on orders orders_1  (cost=2.64..69.67 rows=162 width=8) (actual time=0.007..0.028 rows=137 loops=1)
              Recheck Cond: (amount >= 496)
              Heap Blocks: exact=56
              Buffers: shared hit=58
              ->  Bitmap Index Scan on orders_amount_idx  (cost=0.00..2.60 rows=162 width=0) (actual time=0.004..0.004 rows=137 loops=1)
                    Index Cond: (amount >= 496)
                    Buffers: shared hit=2
Planning:
  Buffers: shared hit=6
Planning Time: 0.049 ms
Execution Time: 0.109 ms

```

## Limits and gaps recorded by the runner

- no concurrent writers: client-decode JSON bytes are not wire bytes; tails are under-sampled at the protocol trial count
- one host: no independent hardware, real-network, 3-node placement or equal-total-budget comparison is established
- single-node only; the YB 3-node RF=3 per-node/equal-total arms are a coverage gap on this host

