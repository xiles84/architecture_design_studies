# SELECT access-path matrix (generated measurements)

_Measurements only. Verdicts for P10/P11 belong in a signed analysis._

**TL;DR:** within-run medians for OR / UNION ALL / UNION at three predicate shapes and OFFSET / keyset at five page depths, each cell verified against the independent oracle. A within-run difference is not a replicated result. Missing measurements are N/A.

| run | engine | topology | scale | page limit |
|---|---|---|---|---|
| select-access-pg-3-20261003 | postgres | single | deep | 5 |

## Provenance

| Field | Value |
|---|---|
| repo commit | `524ce328afa3e738f095b71214051240bcc6f1f0` |
| benchmark image | `localhost/charitybench-select:select-access-pg-3-20261003` |
| image digest | `287420e008898d30b9764fc4a59182fad7670675139c797225182929bfa3622f` |
| run tag | `run/study-01/select-access-pg-3-20261003` |
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
| broad/cast | default | or | valid | 1.332954 | 32.7 | 5 | 4652 | `a4bab1f05f5e667d` |
| broad/cast | default | union | valid | 1.720254 | 19.5 | 5 | 4652 | `6ff77b66c9bf6ee6` |
| broad/cast | default | union_all | valid | 1.848675 | 40.3 | 5 | 4652 | `3f9f3d5c5dc15dd7` |
| broad/expression | expression_index | or | valid | 1.344989 | 44.5 | 5 | 4652 | `a4bab1f05f5e667d` |
| broad/expression | expression_index | union | valid | 1.673449 | 66.4 | 5 | 4652 | `d4fbfb3aac234d58` |
| broad/expression | expression_index | union_all | valid | 1.880469 | 18.6 | 5 | 4652 | `605ee056e65421b1` |
| broad/expression | default | or | valid | 1.277451 | 10.8 | 5 | 4652 | `a4bab1f05f5e667d` |
| broad/expression | default | union | valid | 1.747708 | 12.2 | 5 | 4652 | `6ff77b66c9bf6ee6` |
| broad/expression | default | union_all | valid | 1.759378 | 5.8 | 5 | 4652 | `3f9f3d5c5dc15dd7` |
| broad/indexable | default | or | valid | 1.152580 | 16.1 | 5 | 4652 | `a4bab1f05f5e667d` |
| broad/indexable | default | union | valid | 1.812832 | 42.8 | 5 | 4652 | `ed39b768e5b063b0` |
| broad/indexable | default | union_all | valid | 1.848096 | 25.0 | 5 | 4652 | `d51936c88987c9fd` |
| selective/cast | default | or | valid | 0.849832 | 12.5 | 5 | 251 | `a4bab1f05f5e667d` |
| selective/cast | default | union | valid | 0.809649 | 36.0 | 5 | 251 | `6ff77b66c9bf6ee6` |
| selective/cast | default | union_all | valid | 0.854721 | 24.7 | 5 | 251 | `3f9f3d5c5dc15dd7` |
| selective/expression | expression_index | or | valid | 0.250803 | 40.8 | 5 | 251 | `d66da5af4c9f3384` |
| selective/expression | expression_index | union | valid | 0.283988 | 20.3 | 5 | 251 | `a5fbbd2993065c7a` |
| selective/expression | expression_index | union_all | valid | 0.314781 | 96.8 | 5 | 251 | `eeb1f1b0c0255959` |
| selective/expression | default | or | valid | 0.895906 | 12.5 | 5 | 251 | `a4bab1f05f5e667d` |
| selective/expression | default | union | valid | 0.828612 | 20.8 | 5 | 251 | `6ff77b66c9bf6ee6` |
| selective/expression | default | union_all | valid | 0.858703 | 20.6 | 5 | 251 | `3f9f3d5c5dc15dd7` |
| selective/indexable | default | or | valid | 0.338102 | 29.4 | 5 | 251 | `8bbb087a6c122dbe` |
| selective/indexable | default | union | valid | 0.369664 | 90.6 | 5 | 251 | `567de4baf4feaaf7` |
| selective/indexable | default | union_all | valid | 0.464899 | 78.9 | 5 | 251 | `e060188a3bcf6eb3` |

### P10 pairs

Pairs compare the **default** arm only; the labelled `expression_index` arm is a package arm and is never pooled with the query-text pair.

| pair | A median ms | B median ms | B/A | same plan | within-run verdict |
|---|---|---|---|---|---|
| shape=selective/indexable or vs union_all | 0.338102 | 0.464899 | 1.3750 | false | medians inside the two cells' own spread; no resolved difference within this run |
| shape=selective/indexable or vs union | 0.338102 | 0.369664 | 1.0934 | false | medians inside the two cells' own spread; no resolved difference within this run |
| shape=selective/expression or vs union_all | 0.895906 | 0.858703 | 0.9585 | false | medians inside the two cells' own spread; no resolved difference within this run |
| shape=selective/expression or vs union | 0.895906 | 0.828612 | 0.9249 | false | medians inside the two cells' own spread; no resolved difference within this run |
| shape=selective/cast or vs union_all | 0.849832 | 0.854721 | 1.0058 | false | medians inside the two cells' own spread; no resolved difference within this run |
| shape=selective/cast or vs union | 0.849832 | 0.809649 | 0.9527 | false | medians inside the two cells' own spread; no resolved difference within this run |
| shape=broad/indexable or vs union_all | 1.152580 | 1.848096 | 1.6034 | false | within-run medians differ; whole-run replication required before a verdict |
| shape=broad/indexable or vs union | 1.152580 | 1.812832 | 1.5728 | false | medians inside the two cells' own spread; no resolved difference within this run |
| shape=broad/expression or vs union_all | 1.277451 | 1.759378 | 1.3773 | false | within-run medians differ; whole-run replication required before a verdict |
| shape=broad/expression or vs union | 1.277451 | 1.747708 | 1.3681 | false | within-run medians differ; whole-run replication required before a verdict |
| shape=broad/cast or vs union_all | 1.332954 | 1.848675 | 1.3869 | false | medians inside the two cells' own spread; no resolved difference within this run |
| shape=broad/cast or vs union | 1.332954 | 1.720254 | 1.2906 | false | medians inside the two cells' own spread; no resolved difference within this run |

## Pagination (P11)

| depth | spelling | status | median ms | spread % | trials | rows | plan structure |
|---|---|---|---|---|---|---|---|
| 1 | offset | valid | 0.819929 | 17.6 | 5 | 5 | `ac4509d295bf3348` |
| 1 | keyset | valid | 0.986581 | 11.1 | 5 | 5 | `ac4509d295bf3348` |
| 10 | offset | valid | 0.909782 | 16.3 | 5 | 5 | `ac4509d295bf3348` |
| 10 | keyset | valid | 1.066094 | 8.3 | 5 | 5 | `ac4509d295bf3348` |
| 100 | offset | valid | 0.931667 | 8.6 | 5 | 5 | `ac4509d295bf3348` |
| 100 | keyset | valid | 0.903991 | 28.6 | 5 | 5 | `ac4509d295bf3348` |
| 1000 | offset | valid | 1.270472 | 11.5 | 5 | 5 | `ac4509d295bf3348` |
| 1000 | keyset | valid | 0.930056 | 14.0 | 5 | 5 | `ac4509d295bf3348` |
| 10000 | offset | valid | 2.080571 | 8.4 | 5 | 5 | `ac4509d295bf3348` |
| 10000 | keyset | valid | 0.486106 | 47.3 | 5 | 5 | `dd316474b64ddc29` |

Median ratio table (keyset / OFFSET; a ratio below 1 means keyset was faster in this run):

| depth | ratio (keyset/offset) | offset median ms | keyset median ms |
|---|---|---|---|
| 1 | 1.2033 | 0.819929 | 0.986581 |
| 10 | 1.1718 | 0.909782 | 1.066094 |
| 100 | 0.9703 | 0.931667 | 0.903991 |
| 1000 | 0.7321 | 1.270472 | 0.930056 |
| 10000 | 0.2336 | 2.080571 | 0.486106 |

### P11 pairs

| pair | A median ms | B median ms | B/A | same plan | within-run verdict |
|---|---|---|---|---|---|
| depth=1 offset vs keyset | 0.819929 | 0.986581 | 1.2033 | true | medians inside the two cells' own spread; no resolved difference within this run |
| depth=10 offset vs keyset | 0.909782 | 1.066094 | 1.1718 | true | medians inside the two cells' own spread; no resolved difference within this run |
| depth=100 offset vs keyset | 0.931667 | 0.903991 | 0.9703 | true | medians inside the two cells' own spread; no resolved difference within this run |
| depth=1000 offset vs keyset | 1.270472 | 0.930056 | 0.7321 | true | within-run medians differ; whole-run replication required before a verdict |
| depth=10000 offset vs keyset | 2.080571 | 0.486106 | 0.2336 | false | within-run medians differ; whole-run replication required before a verdict |

## Expression-index arm: write and storage cost

Index `orders_amount_expr_idx` on expression `(amount + 0)`; supported=true.

| measure | value |
|---|---|
| index size (bytes) | 114688 |
| insert batch with index (ms) | 3.983595 |
| insert batch without index (ms) | 4.229128 |
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

Rows: 5; serialized decoded-row JSON bytes: 51 (not wire bytes). Planning/execution: 0.044000 / 0.325000 ms.

```text
Limit  (cost=154.12..154.13 rows=5 width=8) (actual time=0.315..0.316 rows=5 loops=1)
  Buffers: shared hit=73
  ->  Sort  (cost=154.12..158.90 rows=1912 width=8) (actual time=0.314..0.315 rows=5 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 25kB
        Buffers: shared hit=73
        ->  Bitmap Heap Scan on orders  (cost=22.66..122.36 rows=1912 width=8) (actual time=0.036..0.214 rows=2000 loops=1)
              Recheck Cond: ((amount < 67) OR (amount = 67) OR (amount IS NULL))
              Filter: ((amount < 67) OR ((amount = 67) AND (id > 4310)) OR (amount IS NULL))
              Rows Removed by Filter: 6
              Heap Blocks: exact=65
              Buffers: shared hit=73
              ->  BitmapOr  (cost=22.66..22.66 rows=1983 width=0) (actual time=0.028..0.028 rows=0 loops=1)
                    Buffers: shared hit=8
                    ->  Bitmap Index Scan on orders_amount_idx  (cost=0.00..14.38 rows=1440 width=0) (actual time=0.021..0.021 rows=1460 loops=1)
                          Index Cond: (amount < 67)
                          Buffers: shared hit=4
                    ->  Bitmap Index Scan on orders_amount_idx  (cost=0.00..1.54 rows=21 width=0) (actual time=0.001..0.001 rows=24 loops=1)
                          Index Cond: (amount = 67)
                          Buffers: shared hit=2
                    ->  Bitmap Index Scan on orders_amount_idx  (cost=0.00..5.30 rows=522 width=0) (actual time=0.006..0.006 rows=522 loops=1)
                          Index Cond: (amount IS NULL)
                          Buffers: shared hit=2
Planning Time: 0.044 ms
Execution Time: 0.325 ms

```

### s16_page_offset [depth=10000;spelling=offset]

Rows: 5; serialized decoded-row JSON bytes: 51 (not wire bytes). Planning/execution: 0.030000 / 2.952000 ms.

```text
Limit  (cost=1023.04..1023.06 rows=5 width=8) (actual time=2.941..2.942 rows=5 loops=1)
  Buffers: shared hit=65
  ->  Sort  (cost=998.04..1028.04 rows=12000 width=8) (actual time=2.492..2.742 rows=10005 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: quicksort  Memory: 853kB
        Buffers: shared hit=65
        ->  Seq Scan on orders  (cost=0.00..185.00 rows=12000 width=8) (actual time=0.006..0.534 rows=12000 loops=1)
              Buffers: shared hit=65
Planning Time: 0.030 ms
Execution Time: 2.952 ms

```

### s17_page_keyset [depth=1000;spelling=keyset]

Rows: 5; serialized decoded-row JSON bytes: 59 (not wire bytes). Planning/execution: 0.057000 / 1.165000 ms.

```text
Limit  (cost=450.03..450.04 rows=5 width=8) (actual time=1.157..1.157 rows=5 loops=1)
  Buffers: shared hit=65
  ->  Sort  (cost=450.03..476.38 rows=10538 width=8) (actual time=1.156..1.156 rows=5 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 25kB
        Buffers: shared hit=65
        ->  Seq Scan on orders  (cost=0.00..275.00 rows=10538 width=8) (actual time=0.006..0.594 rows=11000 loops=1)
              Filter: ((amount < 457) OR ((amount = 457) AND (id > 11254)) OR (amount IS NULL))
              Rows Removed by Filter: 1000
              Buffers: shared hit=65
Planning Time: 0.057 ms
Execution Time: 1.165 ms

```

### s16_page_offset [depth=1000;spelling=offset]

Rows: 5; serialized decoded-row JSON bytes: 59 (not wire bytes). Planning/execution: 0.022000 / 1.565000 ms.

```text
Limit  (cost=845.88..845.89 rows=5 width=8) (actual time=1.556..1.557 rows=5 loops=1)
  Buffers: shared hit=65
  ->  Sort  (cost=843.38..873.38 rows=12000 width=8) (actual time=1.514..1.537 rows=1005 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 115kB
        Buffers: shared hit=65
        ->  Seq Scan on orders  (cost=0.00..185.00 rows=12000 width=8) (actual time=0.004..0.524 rows=12000 loops=1)
              Buffers: shared hit=65
Planning Time: 0.022 ms
Execution Time: 1.565 ms

```

### s17_page_keyset [depth=100;spelling=keyset]

Rows: 5; serialized decoded-row JSON bytes: 56 (not wire bytes). Planning/execution: 0.049000 / 1.152000 ms.

```text
Limit  (cost=463.65..463.66 rows=5 width=8) (actual time=1.145..1.146 rows=5 loops=1)
  Buffers: shared hit=65
  ->  Sort  (cost=463.65..492.05 rows=11358 width=8) (actual time=1.144..1.145 rows=5 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 25kB
        Buffers: shared hit=65
        ->  Seq Scan on orders  (cost=0.00..275.00 rows=11358 width=8) (actual time=0.005..0.567 rows=11900 loops=1)
              Filter: ((amount < 497) OR ((amount = 497) AND (id > 6723)) OR (amount IS NULL))
              Rows Removed by Filter: 100
              Buffers: shared hit=65
Planning:
  Buffers: shared hit=3
Planning Time: 0.049 ms
Execution Time: 1.152 ms

```

### s16_page_offset [depth=100;spelling=offset]

Rows: 5; serialized decoded-row JSON bytes: 56 (not wire bytes). Planning/execution: 0.047000 / 1.200000 ms.

```text
Limit  (cost=648.10..648.12 rows=5 width=8) (actual time=1.189..1.190 rows=5 loops=1)
  Buffers: shared hit=65
  ->  Sort  (cost=647.85..677.85 rows=12000 width=8) (actual time=1.184..1.187 rows=105 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 31kB
        Buffers: shared hit=65
        ->  Seq Scan on orders  (cost=0.00..185.00 rows=12000 width=8) (actual time=0.008..0.517 rows=12000 loops=1)
              Buffers: shared hit=65
Planning Time: 0.047 ms
Execution Time: 1.200 ms

```

### s17_page_keyset [depth=10;spelling=keyset]

Rows: 5; serialized decoded-row JSON bytes: 56 (not wire bytes). Planning/execution: 0.125000 / 1.142000 ms.

```text
Limit  (cost=465.21..465.23 rows=5 width=8) (actual time=1.134..1.134 rows=5 loops=1)
  Buffers: shared hit=65
  ->  Sort  (cost=465.21..493.84 rows=11452 width=8) (actual time=1.133..1.133 rows=5 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 25kB
        Buffers: shared hit=65
        ->  Seq Scan on orders  (cost=0.00..275.00 rows=11452 width=8) (actual time=0.007..0.555 rows=11990 loops=1)
              Filter: ((amount < 500) OR ((amount = 500) AND (id > 6319)) OR (amount IS NULL))
              Rows Removed by Filter: 10
              Buffers: shared hit=65
Planning:
  Buffers: shared hit=3
Planning Time: 0.125 ms
Execution Time: 1.142 ms

```

### s16_page_offset [depth=10;spelling=offset]

Rows: 5; serialized decoded-row JSON bytes: 56 (not wire bytes). Planning/execution: 0.030000 / 1.082000 ms.

```text
Limit  (cost=479.44..479.45 rows=5 width=8) (actual time=1.075..1.076 rows=5 loops=1)
  Buffers: shared hit=65
  ->  Sort  (cost=479.41..509.41 rows=12000 width=8) (actual time=1.074..1.075 rows=15 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 25kB
        Buffers: shared hit=65
        ->  Seq Scan on orders  (cost=0.00..185.00 rows=12000 width=8) (actual time=0.005..0.478 rows=12000 loops=1)
              Buffers: shared hit=65
Planning Time: 0.030 ms
Execution Time: 1.082 ms

```

### s17_page_keyset [depth=1;spelling=keyset]

Rows: 5; serialized decoded-row JSON bytes: 56 (not wire bytes). Planning/execution: 0.054000 / 1.277000 ms.

```text
Limit  (cost=465.21..465.23 rows=5 width=8) (actual time=1.262..1.264 rows=5 loops=1)
  Buffers: shared hit=65
  ->  Sort  (cost=465.21..493.84 rows=11452 width=8) (actual time=1.262..1.262 rows=5 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 25kB
        Buffers: shared hit=65
        ->  Seq Scan on orders  (cost=0.00..275.00 rows=11452 width=8) (actual time=0.005..0.641 rows=11999 loops=1)
              Filter: ((amount < 500) OR ((amount = 500) AND (id > 1554)) OR (amount IS NULL))
              Rows Removed by Filter: 1
              Buffers: shared hit=65
Planning:
  Buffers: shared hit=3
Planning Time: 0.054 ms
Execution Time: 1.277 ms

```

### s16_page_offset [depth=1;spelling=offset]

Rows: 5; serialized decoded-row JSON bytes: 56 (not wire bytes). Planning/execution: 0.019000 / 1.031000 ms.

```text
Limit  (cost=400.10..400.11 rows=5 width=8) (actual time=1.026..1.027 rows=5 loops=1)
  Buffers: shared hit=65
  ->  Sort  (cost=400.10..430.10 rows=12000 width=8) (actual time=1.025..1.026 rows=6 loops=1)
        Sort Key: amount DESC NULLS LAST, id
        Sort Method: top-N heapsort  Memory: 25kB
        Buffers: shared hit=65
        ->  Seq Scan on orders  (cost=0.00..185.00 rows=12000 width=8) (actual time=0.004..0.458 rows=12000 loops=1)
              Buffers: shared hit=65
Planning Time: 0.019 ms
Execution Time: 1.031 ms

```

### s18_or_predicate [shape=broad/cast;spelling=or]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.050000 / 0.920000 ms.

```text
Index Scan using orders_pkey on orders  (cost=0.29..403.89 rows=6667 width=8) (actual time=0.011..0.813 rows=4652 loops=1)
  Filter: (((amount)::bigint <= 100) OR ((amount)::bigint >= 400))
  Rows Removed by Filter: 7348
  Buffers: shared hit=99
Planning Time: 0.050 ms
Execution Time: 0.920 ms

```

### s20_union_predicate [shape=broad/cast;spelling=union]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.070000 / 1.800000 ms.

```text
Unique  (cost=968.64..1108.64 rows=8000 width=8) (actual time=1.051..1.690 rows=4652 loops=1)
  Buffers: shared hit=130
  ->  Merge Append  (cost=968.64..1068.64 rows=8000 width=8) (actual time=1.051..1.410 rows=4652 loops=1)
        Sort Key: orders.id, orders.amount
        Buffers: shared hit=130
        ->  Sort  (cost=484.32..494.32 rows=4000 width=8) (actual time=0.545..0.598 rows=2273 loops=1)
              Sort Key: orders.id, orders.amount
              Sort Method: quicksort  Memory: 150kB
              Buffers: shared hit=65
              ->  Seq Scan on orders  (cost=0.00..245.00 rows=4000 width=8) (actual time=0.008..0.405 rows=2273 loops=1)
                    Filter: ((amount)::bigint <= 100)
                    Rows Removed by Filter: 9727
                    Buffers: shared hit=65
        ->  Sort  (cost=484.32..494.32 rows=4000 width=8) (actual time=0.505..0.556 rows=2379 loops=1)
              Sort Key: orders_1.id, orders_1.amount
              Sort Method: quicksort  Memory: 152kB
              Buffers: shared hit=65
              ->  Seq Scan on orders orders_1  (cost=0.00..245.00 rows=4000 width=8) (actual time=0.003..0.377 rows=2379 loops=1)
                    Filter: ((amount)::bigint >= 400)
                    Rows Removed by Filter: 9621
                    Buffers: shared hit=65
Planning Time: 0.070 ms
Execution Time: 1.800 ms

```

### s19_union_all_predicate [shape=broad/cast;spelling=union_all]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.039000 / 1.474000 ms.

```text
Sort  (cost=1048.63..1068.63 rows=8000 width=8) (actual time=1.240..1.370 rows=4652 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 302kB
  Buffers: shared hit=130
  ->  Append  (cost=0.00..530.00 rows=8000 width=8) (actual time=0.006..0.881 rows=4652 loops=1)
        Buffers: shared hit=130
        ->  Seq Scan on orders  (cost=0.00..245.00 rows=4000 width=8) (actual time=0.005..0.347 rows=2273 loops=1)
              Filter: ((amount)::bigint <= 100)
              Rows Removed by Filter: 9727
              Buffers: shared hit=65
        ->  Seq Scan on orders orders_1  (cost=0.00..245.00 rows=4000 width=8) (actual time=0.003..0.356 rows=2379 loops=1)
              Filter: ((amount)::bigint >= 400)
              Rows Removed by Filter: 9621
              Buffers: shared hit=65
Planning Time: 0.039 ms
Execution Time: 1.474 ms

```

### s18_or_predicate [shape=broad/expression;index=expression_index;spelling=or]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.029000 / 0.858000 ms.

```text
Index Scan using orders_pkey on orders  (cost=0.29..403.89 rows=4211 width=8) (actual time=0.006..0.755 rows=4652 loops=1)
  Filter: (((amount + 0) <= 100) OR ((amount + 0) >= 400))
  Rows Removed by Filter: 7348
  Buffers: shared hit=99
Planning Time: 0.029 ms
Execution Time: 0.858 ms

```

### s20_union_predicate [shape=broad/expression;index=expression_index;spelling=union]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.081000 / 1.821000 ms.

```text
Unique  (cost=492.69..574.31 rows=4664 width=8) (actual time=0.983..1.687 rows=4652 loops=1)
  Buffers: shared hit=3934
  ->  Merge Append  (cost=492.69..550.99 rows=4664 width=8) (actual time=0.983..1.375 rows=4652 loops=1)
        Sort Key: orders.id, orders.amount
        Buffers: shared hit=3934
        ->  Sort  (cost=241.10..246.77 rows=2267 width=8) (actual time=0.523..0.580 rows=2273 loops=1)
              Sort Key: orders.id, orders.amount
              Sort Method: quicksort  Memory: 150kB
              Buffers: shared hit=1921
              ->  Index Scan using orders_amount_expr_idx on orders  (cost=0.29..114.75 rows=2267 width=8) (actual time=0.011..0.312 rows=2273 loops=1)
                    Index Cond: ((amount + 0) <= 100)
                    Buffers: shared hit=1921
        ->  Sort  (cost=251.58..257.57 rows=2397 width=8) (actual time=0.459..0.518 rows=2379 loops=1)
              Sort Key: orders_1.id, orders_1.amount
              Sort Method: quicksort  Memory: 152kB
              Buffers: shared hit=2013
              ->  Index Scan using orders_amount_expr_idx on orders orders_1  (cost=0.29..117.03 rows=2397 width=8) (actual time=0.003..0.249 rows=2379 loops=1)
                    Index Cond: ((amount + 0) >= 400)
                    Buffers: shared hit=2013
Planning Time: 0.081 ms
Execution Time: 1.821 ms

```

### s19_union_all_predicate [shape=broad/expression;index=expression_index;spelling=union_all]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.041000 / 1.471000 ms.

```text
Sort  (cost=539.31..550.97 rows=4664 width=8) (actual time=1.251..1.364 rows=4652 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 302kB
  Buffers: shared hit=3934
  ->  Append  (cost=0.29..255.10 rows=4664 width=8) (actual time=0.006..0.720 rows=4652 loops=1)
        Buffers: shared hit=3934
        ->  Index Scan using orders_amount_expr_idx on orders  (cost=0.29..114.75 rows=2267 width=8) (actual time=0.005..0.287 rows=2273 loops=1)
              Index Cond: ((amount + 0) <= 100)
              Buffers: shared hit=1921
        ->  Index Scan using orders_amount_expr_idx on orders orders_1  (cost=0.29..117.03 rows=2397 width=8) (actual time=0.002..0.252 rows=2379 loops=1)
              Index Cond: ((amount + 0) >= 400)
              Buffers: shared hit=2013
Planning Time: 0.041 ms
Execution Time: 1.471 ms

```

### s18_or_predicate [shape=broad/expression;spelling=or]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.034000 / 0.903000 ms.

```text
Index Scan using orders_pkey on orders  (cost=0.29..403.89 rows=6667 width=8) (actual time=0.009..0.795 rows=4652 loops=1)
  Filter: (((amount + 0) <= 100) OR ((amount + 0) >= 400))
  Rows Removed by Filter: 7348
  Buffers: shared hit=99
Planning Time: 0.034 ms
Execution Time: 0.903 ms

```

### s20_union_predicate [shape=broad/expression;spelling=union]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.043000 / 1.964000 ms.

```text
Unique  (cost=968.64..1108.64 rows=8000 width=8) (actual time=1.017..1.769 rows=4652 loops=1)
  Buffers: shared hit=130
  ->  Merge Append  (cost=968.64..1068.64 rows=8000 width=8) (actual time=1.017..1.461 rows=4652 loops=1)
        Sort Key: orders.id, orders.amount
        Buffers: shared hit=130
        ->  Sort  (cost=484.32..494.32 rows=4000 width=8) (actual time=0.510..0.568 rows=2273 loops=1)
              Sort Key: orders.id, orders.amount
              Sort Method: quicksort  Memory: 150kB
              Buffers: shared hit=65
              ->  Seq Scan on orders  (cost=0.00..245.00 rows=4000 width=8) (actual time=0.005..0.387 rows=2273 loops=1)
                    Filter: ((amount + 0) <= 100)
                    Rows Removed by Filter: 9727
                    Buffers: shared hit=65
        ->  Sort  (cost=484.32..494.32 rows=4000 width=8) (actual time=0.506..0.567 rows=2379 loops=1)
              Sort Key: orders_1.id, orders_1.amount
              Sort Method: quicksort  Memory: 152kB
              Buffers: shared hit=65
              ->  Seq Scan on orders orders_1  (cost=0.00..245.00 rows=4000 width=8) (actual time=0.004..0.375 rows=2379 loops=1)
                    Filter: ((amount + 0) >= 400)
                    Rows Removed by Filter: 9621
                    Buffers: shared hit=65
Planning Time: 0.043 ms
Execution Time: 1.964 ms

```

### s19_union_all_predicate [shape=broad/expression;spelling=union_all]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.078000 / 1.460000 ms.

```text
Sort  (cost=1048.63..1068.63 rows=8000 width=8) (actual time=1.256..1.362 rows=4652 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 302kB
  Buffers: shared hit=130
  ->  Append  (cost=0.00..530.00 rows=8000 width=8) (actual time=0.013..0.896 rows=4652 loops=1)
        Buffers: shared hit=130
        ->  Seq Scan on orders  (cost=0.00..245.00 rows=4000 width=8) (actual time=0.013..0.373 rows=2273 loops=1)
              Filter: ((amount + 0) <= 100)
              Rows Removed by Filter: 9727
              Buffers: shared hit=65
        ->  Seq Scan on orders orders_1  (cost=0.00..245.00 rows=4000 width=8) (actual time=0.003..0.340 rows=2379 loops=1)
              Filter: ((amount + 0) >= 400)
              Rows Removed by Filter: 9621
              Buffers: shared hit=65
Planning Time: 0.078 ms
Execution Time: 1.460 ms

```

### s18_or_predicate [shape=broad/indexable;spelling=or]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.053000 / 0.771000 ms.

```text
Index Scan using orders_pkey on orders  (cost=0.29..343.89 rows=4211 width=8) (actual time=0.022..0.675 rows=4652 loops=1)
  Filter: ((amount <= 100) OR (amount >= 400))
  Rows Removed by Filter: 7348
  Buffers: shared hit=99
Planning Time: 0.053 ms
Execution Time: 0.771 ms

```

### s20_union_predicate [shape=broad/indexable;spelling=union]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.083000 / 1.755000 ms.

```text
Unique  (cost=495.99..577.61 rows=4664 width=8) (actual time=1.029..1.653 rows=4652 loops=1)
  Buffers: shared hit=3936
  ->  Merge Append  (cost=495.99..554.29 rows=4664 width=8) (actual time=1.029..1.394 rows=4652 loops=1)
        Sort Key: orders.id, orders.amount
        Buffers: shared hit=3936
        ->  Sort  (cost=242.20..247.87 rows=2267 width=8) (actual time=0.520..0.575 rows=2273 loops=1)
              Sort Key: orders.id, orders.amount
              Sort Method: quicksort  Memory: 150kB
              Buffers: shared hit=1922
              ->  Index Scan using orders_amount_idx on orders  (cost=0.29..115.85 rows=2267 width=8) (actual time=0.011..0.293 rows=2273 loops=1)
                    Index Cond: (amount <= 100)
                    Buffers: shared hit=1922
        ->  Sort  (cost=253.78..259.77 rows=2397 width=8) (actual time=0.508..0.560 rows=2379 loops=1)
              Sort Key: orders_1.id, orders_1.amount
              Sort Method: quicksort  Memory: 152kB
              Buffers: shared hit=2014
              ->  Index Scan using orders_amount_idx on orders orders_1  (cost=0.29..119.23 rows=2397 width=8) (actual time=0.005..0.277 rows=2379 loops=1)
                    Index Cond: (amount >= 400)
                    Buffers: shared hit=2014
Planning Time: 0.083 ms
Execution Time: 1.755 ms

```

### s19_union_all_predicate [shape=broad/indexable;spelling=union_all]

Rows: 4652; serialized decoded-row JSON bytes: 49919 (not wire bytes). Planning/execution: 0.048000 / 1.384000 ms.

```text
Sort  (cost=542.61..554.27 rows=4664 width=8) (actual time=1.187..1.289 rows=4652 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 302kB
  Buffers: shared hit=3936
  ->  Append  (cost=0.29..258.40 rows=4664 width=8) (actual time=0.007..0.738 rows=4652 loops=1)
        Buffers: shared hit=3936
        ->  Index Scan using orders_amount_idx on orders  (cost=0.29..115.85 rows=2267 width=8) (actual time=0.007..0.283 rows=2273 loops=1)
              Index Cond: (amount <= 100)
              Buffers: shared hit=1922
        ->  Index Scan using orders_amount_idx on orders orders_1  (cost=0.29..119.23 rows=2397 width=8) (actual time=0.003..0.270 rows=2379 loops=1)
              Index Cond: (amount >= 400)
              Buffers: shared hit=2014
Planning Time: 0.048 ms
Execution Time: 1.384 ms

```

### s18_or_predicate [shape=selective/cast;spelling=or]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.048000 / 0.602000 ms.

```text
Index Scan using orders_pkey on orders  (cost=0.29..403.89 rows=6667 width=8) (actual time=0.008..0.591 rows=251 loops=1)
  Filter: (((amount)::bigint <= 5) OR ((amount)::bigint >= 496))
  Rows Removed by Filter: 11749
  Buffers: shared hit=99
Planning Time: 0.048 ms
Execution Time: 0.602 ms

```

### s20_union_predicate [shape=selective/cast;spelling=union]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.049000 / 0.714000 ms.

```text
Unique  (cost=968.64..1108.64 rows=8000 width=8) (actual time=0.634..0.672 rows=251 loops=1)
  Buffers: shared hit=130
  ->  Merge Append  (cost=968.64..1068.64 rows=8000 width=8) (actual time=0.634..0.655 rows=251 loops=1)
        Sort Key: orders.id, orders.amount
        Buffers: shared hit=130
        ->  Sort  (cost=484.32..494.32 rows=4000 width=8) (actual time=0.317..0.320 rows=114 loops=1)
              Sort Key: orders.id, orders.amount
              Sort Method: quicksort  Memory: 27kB
              Buffers: shared hit=65
              ->  Seq Scan on orders  (cost=0.00..245.00 rows=4000 width=8) (actual time=0.008..0.307 rows=114 loops=1)
                    Filter: ((amount)::bigint <= 5)
                    Rows Removed by Filter: 11886
                    Buffers: shared hit=65
        ->  Sort  (cost=484.32..494.32 rows=4000 width=8) (actual time=0.316..0.320 rows=137 loops=1)
              Sort Key: orders_1.id, orders_1.amount
              Sort Method: quicksort  Memory: 28kB
              Buffers: shared hit=65
              ->  Seq Scan on orders orders_1  (cost=0.00..245.00 rows=4000 width=8) (actual time=0.004..0.308 rows=137 loops=1)
                    Filter: ((amount)::bigint >= 496)
                    Rows Removed by Filter: 11863
                    Buffers: shared hit=65
Planning Time: 0.049 ms
Execution Time: 0.714 ms

```

### s19_union_all_predicate [shape=selective/cast;spelling=union_all]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.029000 / 0.615000 ms.

```text
Sort  (cost=1048.63..1068.63 rows=8000 width=8) (actual time=0.597..0.603 rows=251 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 30kB
  Buffers: shared hit=130
  ->  Append  (cost=0.00..530.00 rows=8000 width=8) (actual time=0.007..0.577 rows=251 loops=1)
        Buffers: shared hit=130
        ->  Seq Scan on orders  (cost=0.00..245.00 rows=4000 width=8) (actual time=0.007..0.284 rows=114 loops=1)
              Filter: ((amount)::bigint <= 5)
              Rows Removed by Filter: 11886
              Buffers: shared hit=65
        ->  Seq Scan on orders orders_1  (cost=0.00..245.00 rows=4000 width=8) (actual time=0.004..0.283 rows=137 loops=1)
              Filter: ((amount)::bigint >= 496)
              Rows Removed by Filter: 11863
              Buffers: shared hit=65
Planning Time: 0.029 ms
Execution Time: 0.615 ms

```

### s18_or_predicate [shape=selective/expression;index=expression_index;spelling=or]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.042000 / 0.076000 ms.

```text
Sort  (cost=86.11..86.78 rows=268 width=8) (actual time=0.058..0.064 rows=251 loops=1)
  Sort Key: id
  Sort Method: quicksort  Memory: 30kB
  Buffers: shared hit=67
  ->  Bitmap Heap Scan on orders  (cost=4.92..75.30 rows=268 width=8) (actual time=0.012..0.044 rows=251 loops=1)
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
Planning Time: 0.042 ms
Execution Time: 0.076 ms

```

### s20_union_predicate [shape=selective/expression;index=expression_index;spelling=union]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.064000 / 0.145000 ms.

```text
Unique  (cost=142.73..147.44 rows=269 width=8) (actual time=0.082..0.129 rows=251 loops=1)
  Buffers: shared hit=114
  ->  Merge Append  (cost=142.73..146.09 rows=269 width=8) (actual time=0.082..0.110 rows=251 loops=1)
        Sort Key: orders.id, orders.amount
        Buffers: shared hit=114
        ->  Sort  (cost=66.70..66.97 rows=107 width=8) (actual time=0.042..0.045 rows=114 loops=1)
              Sort Key: orders.id, orders.amount
              Sort Method: quicksort  Memory: 27kB
              Buffers: shared hit=56
              ->  Bitmap Heap Scan on orders  (cost=2.21..63.10 rows=107 width=8) (actual time=0.010..0.033 rows=114 loops=1)
                    Recheck Cond: ((amount + 0) <= 5)
                    Heap Blocks: exact=54
                    Buffers: shared hit=56
                    ->  Bitmap Index Scan on orders_amount_expr_idx  (cost=0.00..2.19 rows=107 width=0) (actual time=0.004..0.004 rows=114 loops=1)
                          Index Cond: ((amount + 0) <= 5)
                          Buffers: shared hit=2
        ->  Sort  (cost=76.02..76.42 rows=162 width=8) (actual time=0.039..0.043 rows=137 loops=1)
              Sort Key: orders_1.id, orders_1.amount
              Sort Method: quicksort  Memory: 28kB
              Buffers: shared hit=58
              ->  Bitmap Heap Scan on orders orders_1  (cost=2.64..70.07 rows=162 width=8) (actual time=0.008..0.031 rows=137 loops=1)
                    Recheck Cond: ((amount + 0) >= 496)
                    Heap Blocks: exact=56
                    Buffers: shared hit=58
                    ->  Bitmap Index Scan on orders_amount_expr_idx  (cost=0.00..2.60 rows=162 width=0) (actual time=0.004..0.004 rows=137 loops=1)
                          Index Cond: ((amount + 0) >= 496)
                          Buffers: shared hit=2
Planning:
  Buffers: shared hit=6
Planning Time: 0.064 ms
Execution Time: 0.145 ms

```

### s19_union_all_predicate [shape=selective/expression;index=expression_index;spelling=union_all]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.045000 / 0.107000 ms.

```text
Sort  (cost=145.37..146.04 rows=269 width=8) (actual time=0.088..0.095 rows=251 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 30kB
  Buffers: shared hit=114
  ->  Append  (cost=2.21..134.51 rows=269 width=8) (actual time=0.008..0.066 rows=251 loops=1)
        Buffers: shared hit=114
        ->  Bitmap Heap Scan on orders  (cost=2.21..63.10 rows=107 width=8) (actual time=0.008..0.029 rows=114 loops=1)
              Recheck Cond: ((amount + 0) <= 5)
              Heap Blocks: exact=54
              Buffers: shared hit=56
              ->  Bitmap Index Scan on orders_amount_expr_idx  (cost=0.00..2.19 rows=107 width=0) (actual time=0.003..0.003 rows=114 loops=1)
                    Index Cond: ((amount + 0) <= 5)
                    Buffers: shared hit=2
        ->  Bitmap Heap Scan on orders orders_1  (cost=2.64..70.07 rows=162 width=8) (actual time=0.007..0.027 rows=137 loops=1)
              Recheck Cond: ((amount + 0) >= 496)
              Heap Blocks: exact=56
              Buffers: shared hit=58
              ->  Bitmap Index Scan on orders_amount_expr_idx  (cost=0.00..2.60 rows=162 width=0) (actual time=0.003..0.003 rows=137 loops=1)
                    Index Cond: ((amount + 0) >= 496)
                    Buffers: shared hit=2
Planning:
  Buffers: shared hit=6
Planning Time: 0.045 ms
Execution Time: 0.107 ms

```

### s18_or_predicate [shape=selective/expression;spelling=or]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.021000 / 0.729000 ms.

```text
Index Scan using orders_pkey on orders  (cost=0.29..403.89 rows=6667 width=8) (actual time=0.008..0.719 rows=251 loops=1)
  Filter: (((amount + 0) <= 5) OR ((amount + 0) >= 496))
  Rows Removed by Filter: 11749
  Buffers: shared hit=99
Planning Time: 0.021 ms
Execution Time: 0.729 ms

```

### s20_union_predicate [shape=selective/expression;spelling=union]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.038000 / 0.661000 ms.

```text
Unique  (cost=968.64..1108.64 rows=8000 width=8) (actual time=0.609..0.647 rows=251 loops=1)
  Buffers: shared hit=130
  ->  Merge Append  (cost=968.64..1068.64 rows=8000 width=8) (actual time=0.608..0.629 rows=251 loops=1)
        Sort Key: orders.id, orders.amount
        Buffers: shared hit=130
        ->  Sort  (cost=484.32..494.32 rows=4000 width=8) (actual time=0.294..0.297 rows=114 loops=1)
              Sort Key: orders.id, orders.amount
              Sort Method: quicksort  Memory: 27kB
              Buffers: shared hit=65
              ->  Seq Scan on orders  (cost=0.00..245.00 rows=4000 width=8) (actual time=0.005..0.286 rows=114 loops=1)
                    Filter: ((amount + 0) <= 5)
                    Rows Removed by Filter: 11886
                    Buffers: shared hit=65
        ->  Sort  (cost=484.32..494.32 rows=4000 width=8) (actual time=0.314..0.317 rows=137 loops=1)
              Sort Key: orders_1.id, orders_1.amount
              Sort Method: quicksort  Memory: 28kB
              Buffers: shared hit=65
              ->  Seq Scan on orders orders_1  (cost=0.00..245.00 rows=4000 width=8) (actual time=0.004..0.306 rows=137 loops=1)
                    Filter: ((amount + 0) >= 496)
                    Rows Removed by Filter: 11863
                    Buffers: shared hit=65
Planning Time: 0.038 ms
Execution Time: 0.661 ms

```

### s19_union_all_predicate [shape=selective/expression;spelling=union_all]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.027000 / 0.634000 ms.

```text
Sort  (cost=1048.63..1068.63 rows=8000 width=8) (actual time=0.617..0.623 rows=251 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 30kB
  Buffers: shared hit=130
  ->  Append  (cost=0.00..530.00 rows=8000 width=8) (actual time=0.005..0.588 rows=251 loops=1)
        Buffers: shared hit=130
        ->  Seq Scan on orders  (cost=0.00..245.00 rows=4000 width=8) (actual time=0.005..0.272 rows=114 loops=1)
              Filter: ((amount + 0) <= 5)
              Rows Removed by Filter: 11886
              Buffers: shared hit=65
        ->  Seq Scan on orders orders_1  (cost=0.00..245.00 rows=4000 width=8) (actual time=0.003..0.306 rows=137 loops=1)
              Filter: ((amount + 0) >= 496)
              Rows Removed by Filter: 11863
              Buffers: shared hit=65
Planning Time: 0.027 ms
Execution Time: 0.634 ms

```

### s18_or_predicate [shape=selective/indexable;spelling=or]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.042000 / 0.100000 ms.

```text
Sort  (cost=84.77..85.44 rows=268 width=8) (actual time=0.078..0.086 rows=251 loops=1)
  Sort Key: id
  Sort Method: quicksort  Memory: 30kB
  Buffers: shared hit=67
  ->  Bitmap Heap Scan on orders  (cost=4.92..73.96 rows=268 width=8) (actual time=0.015..0.060 rows=251 loops=1)
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
Planning Time: 0.042 ms
Execution Time: 0.100 ms

```

### s20_union_predicate [shape=selective/indexable;spelling=union]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.054000 / 0.124000 ms.

```text
Unique  (cost=142.06..146.77 rows=269 width=8) (actual time=0.074..0.110 rows=251 loops=1)
  Buffers: shared hit=114
  ->  Merge Append  (cost=142.06..145.42 rows=269 width=8) (actual time=0.073..0.094 rows=251 loops=1)
        Sort Key: orders.id, orders.amount
        Buffers: shared hit=114
        ->  Sort  (cost=66.44..66.70 rows=107 width=8) (actual time=0.038..0.041 rows=114 loops=1)
              Sort Key: orders.id, orders.amount
              Sort Method: quicksort  Memory: 27kB
              Buffers: shared hit=56
              ->  Bitmap Heap Scan on orders  (cost=2.21..62.83 rows=107 width=8) (actual time=0.009..0.031 rows=114 loops=1)
                    Recheck Cond: (amount <= 5)
                    Heap Blocks: exact=54
                    Buffers: shared hit=56
                    ->  Bitmap Index Scan on orders_amount_idx  (cost=0.00..2.19 rows=107 width=0) (actual time=0.004..0.004 rows=114 loops=1)
                          Index Cond: (amount <= 5)
                          Buffers: shared hit=2
        ->  Sort  (cost=75.61..76.02 rows=162 width=8) (actual time=0.035..0.038 rows=137 loops=1)
              Sort Key: orders_1.id, orders_1.amount
              Sort Method: quicksort  Memory: 28kB
              Buffers: shared hit=58
              ->  Bitmap Heap Scan on orders orders_1  (cost=2.64..69.67 rows=162 width=8) (actual time=0.007..0.027 rows=137 loops=1)
                    Recheck Cond: (amount >= 496)
                    Heap Blocks: exact=56
                    Buffers: shared hit=58
                    ->  Bitmap Index Scan on orders_amount_idx  (cost=0.00..2.60 rows=162 width=0) (actual time=0.003..0.003 rows=137 loops=1)
                          Index Cond: (amount >= 496)
                          Buffers: shared hit=2
Planning:
  Buffers: shared hit=6
Planning Time: 0.054 ms
Execution Time: 0.124 ms

```

### s19_union_all_predicate [shape=selective/indexable;spelling=union_all]

Rows: 251; serialized decoded-row JSON bytes: 2609 (not wire bytes). Planning/execution: 0.061000 / 0.122000 ms.

```text
Sort  (cost=144.70..145.37 rows=269 width=8) (actual time=0.101..0.108 rows=251 loops=1)
  Sort Key: orders.id
  Sort Method: quicksort  Memory: 30kB
  Buffers: shared hit=114
  ->  Append  (cost=2.21..133.84 rows=269 width=8) (actual time=0.012..0.078 rows=251 loops=1)
        Buffers: shared hit=114
        ->  Bitmap Heap Scan on orders  (cost=2.21..62.83 rows=107 width=8) (actual time=0.012..0.038 rows=114 loops=1)
              Recheck Cond: (amount <= 5)
              Heap Blocks: exact=54
              Buffers: shared hit=56
              ->  Bitmap Index Scan on orders_amount_idx  (cost=0.00..2.19 rows=107 width=0) (actual time=0.004..0.005 rows=114 loops=1)
                    Index Cond: (amount <= 5)
                    Buffers: shared hit=2
        ->  Bitmap Heap Scan on orders orders_1  (cost=2.64..69.67 rows=162 width=8) (actual time=0.008..0.029 rows=137 loops=1)
              Recheck Cond: (amount >= 496)
              Heap Blocks: exact=56
              Buffers: shared hit=58
              ->  Bitmap Index Scan on orders_amount_idx  (cost=0.00..2.60 rows=162 width=0) (actual time=0.004..0.004 rows=137 loops=1)
                    Index Cond: (amount >= 496)
                    Buffers: shared hit=2
Planning:
  Buffers: shared hit=6
Planning Time: 0.061 ms
Execution Time: 0.122 ms

```

## Limits and gaps recorded by the runner

- no concurrent writers: client-decode JSON bytes are not wire bytes; tails are under-sampled at the protocol trial count
- one host: no independent hardware, real-network, 3-node placement or equal-total-budget comparison is established
- single-node only; the YB 3-node RF=3 per-node/equal-total arms are a coverage gap on this host

