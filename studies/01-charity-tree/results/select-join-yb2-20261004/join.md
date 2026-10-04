# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 91 cells, 91 valid, engine `yugabyte`, topology `single`, scale `cte`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-yb2-20261004`
- run tag: `run/study-01/select-join-yb2-20261004`
- environment: `host-zenbook-ux5406sa`
- repo commit: `7e148174b5211225f6b4f243fc02e0de417a0586`
- engine version: `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit`
- engine image: `localhost/charitybench-select:select-join-yb2-20261004` @ `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`
- inputs digest: `25e5f7bf01de5149bed501f089eebe900d7b462c2927eedaee302436b664d198`
- dataset fingerprint: `2e06ecb3915eb290`
- client CPU quota: `200000 100000`, memory: `2147483648`
- setting `client_workers`: 1 sequential reader; 0 writers
- setting `connection_pool`: 2 = 1 reader/setup + 1 reserve
- setting `customers`: 2000
- setting `distributed_rpc`: not parsed; raw EXPLAIN retained (coverage gap)
- setting `fact_count`: 200000
- setting `index_set`: facts_rank_idx(rank),facts_customer_idx(customer_id),primary keys; no per-spelling index
- setting `prepared_plan_policy`: pgx simple protocol; no prepared statements
- setting `scale`: cte
- setting `selectivity_mapping`: rank <= threshold; rank is a deterministic shuffle of 1..fact_count
- setting `statement_budget_ms`: 120000
- setting `statistics`: ANALYZE after deterministic load
- setting `thresholds`: [20 200 2000 20000 100000 180000]
- setting `topology`: single
- setting `trials_per_cell`: 5
- setting `whole_run_count`: 1; no replicated winner in this file

## Fixture and selectivity mapping

- fact rows: 200000
- customers: 2000
- rank thresholds: [20 200 2000 20000 100000 180000]

The predicate is `rank <= threshold`, where `rank` is a deterministic shuffle of 1..200000; `rank` is only a labelled proxy for the teaching amount predicate.

## Per-engine feasibility evidence

- fact rows 200000, customers 2000, thresholds [20 200 2000 20000 100000 180000]
- expensive reuse token supported by the engine: yes

| threshold | rows with rank <= threshold |
|---|---|
| 100000 | 100000 |
| 180000 | 180000 |
| 20 | 20 |
| 200 | 200 |
| 2000 | 2000 |
| 20000 | 20000 |

| statement | support | rank index observed | forced policy evidence |
|---|---|---|---|
| j01_join_where | valid | false |  |
| j02_join_on | valid | false |  |
| j03_filtered_derived | valid | false |  |
| j04_cte_filtered | valid | false | policy=default; observed_cte_scan=false |
| j04m_cte_filtered_mat | valid | false | policy=materialized; observed_cte_scan=true |
| j04n_cte_filtered_notmat | valid | false | policy=not_materialized; observed_cte_scan=false |
| j05_cte_joined_outer | valid | false | policy=default; observed_cte_scan=false |
| j05m_cte_joined_outer_mat | valid | false | policy=materialized; observed_cte_scan=true |
| j05n_cte_joined_outer_notmat | valid | false | policy=not_materialized; observed_cte_scan=false |
| j06_comma_join | valid | false |  |
| j11_optional_on | valid | false |  |
| r_default_n4_c100 | valid | false | policy=default; observed_cte_scan=true |
| r_mat_n4_c100 | valid | false | policy=materialized; observed_cte_scan=true |
| r_notmat_n4_c100 | valid | false | policy=not_materialized; observed_cte_scan=false |

## P1 — six inner-join spellings and their CTE policies (median ms)

| threshold | join_where | join_on | filtered_derived | cte_filtered | cte_filtered_mat | cte_filtered_notmat | cte_joined_outer | cte_joined_outer_mat | cte_joined_outer_notmat | comma_join |
|---|---|---|---|---|---|---|---|---|---|---|
| 20 (0.01%) | 56.9529 | 58.2015 | 56.2280 | 55.3352 | 55.1612 | 56.1700 | 57.1734 | 86.1778 | 54.7682 | 56.2012 |
| 200 (0.1%) | 56.2774 | 55.9788 | 55.3489 | 58.5707 | 55.2458 | 56.9247 | 58.7040 | 89.9664 | 57.9382 | 59.1745 |
| 2000 (1%) | 57.6083 | 56.5501 | 57.0162 | 57.1560 | 57.4772 | 58.5656 | 56.2968 | 91.8593 | 57.1141 | 59.6796 |
| 20000 (10%) | 73.9714 | 70.6003 | 71.1113 | 70.3850 | 71.3000 | 70.5793 | 68.1864 | 94.1146 | 70.8516 | 69.9753 |
| 100000 (50%) | 89.9328 | 100.1651 | 94.9401 | 93.8807 | 101.3366 | 94.8852 | 96.6752 | 101.4762 | 93.4475 | 94.6713 |
| 180000 (90%) | 119.3936 | 125.0222 | 117.8714 | 113.9327 | 130.7327 | 118.1213 | 122.3637 | 108.2137 | 116.5785 | 114.8336 |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | join_on vs join_where | 56.9529 | 58.2015 | 1.0219 | true |
| 20 | filtered_derived vs join_where | 56.9529 | 56.2280 | 0.9873 | true |
| 20 | cte_filtered vs join_where | 56.9529 | 55.3352 | 0.9716 | true |
| 20 | cte_filtered_mat vs join_where | 56.9529 | 55.1612 | 0.9685 | false |
| 20 | cte_filtered_notmat vs join_where | 56.9529 | 56.1700 | 0.9863 | true |
| 20 | cte_joined_outer vs join_where | 56.9529 | 57.1734 | 1.0039 | true |
| 20 | cte_joined_outer_mat vs join_where | 56.9529 | 86.1778 | 1.5131 | false |
| 20 | cte_joined_outer_notmat vs join_where | 56.9529 | 54.7682 | 0.9616 | true |
| 20 | comma_join vs join_where | 56.9529 | 56.2012 | 0.9868 | true |
| 200 | join_on vs join_where | 56.2774 | 55.9788 | 0.9947 | true |
| 200 | filtered_derived vs join_where | 56.2774 | 55.3489 | 0.9835 | true |
| 200 | cte_filtered vs join_where | 56.2774 | 58.5707 | 1.0408 | true |
| 200 | cte_filtered_mat vs join_where | 56.2774 | 55.2458 | 0.9817 | false |
| 200 | cte_filtered_notmat vs join_where | 56.2774 | 56.9247 | 1.0115 | true |
| 200 | cte_joined_outer vs join_where | 56.2774 | 58.7040 | 1.0431 | true |
| 200 | cte_joined_outer_mat vs join_where | 56.2774 | 89.9664 | 1.5986 | false |
| 200 | cte_joined_outer_notmat vs join_where | 56.2774 | 57.9382 | 1.0295 | true |
| 200 | comma_join vs join_where | 56.2774 | 59.1745 | 1.0515 | true |
| 2000 | join_on vs join_where | 57.6083 | 56.5501 | 0.9816 | true |
| 2000 | filtered_derived vs join_where | 57.6083 | 57.0162 | 0.9897 | true |
| 2000 | cte_filtered vs join_where | 57.6083 | 57.1560 | 0.9921 | true |
| 2000 | cte_filtered_mat vs join_where | 57.6083 | 57.4772 | 0.9977 | false |
| 2000 | cte_filtered_notmat vs join_where | 57.6083 | 58.5656 | 1.0166 | true |
| 2000 | cte_joined_outer vs join_where | 57.6083 | 56.2968 | 0.9772 | true |
| 2000 | cte_joined_outer_mat vs join_where | 57.6083 | 91.8593 | 1.5946 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 57.6083 | 57.1141 | 0.9914 | true |
| 2000 | comma_join vs join_where | 57.6083 | 59.6796 | 1.0360 | true |
| 20000 | join_on vs join_where | 73.9714 | 70.6003 | 0.9544 | true |
| 20000 | filtered_derived vs join_where | 73.9714 | 71.1113 | 0.9613 | true |
| 20000 | cte_filtered vs join_where | 73.9714 | 70.3850 | 0.9515 | true |
| 20000 | cte_filtered_mat vs join_where | 73.9714 | 71.3000 | 0.9639 | false |
| 20000 | cte_filtered_notmat vs join_where | 73.9714 | 70.5793 | 0.9541 | true |
| 20000 | cte_joined_outer vs join_where | 73.9714 | 68.1864 | 0.9218 | true |
| 20000 | cte_joined_outer_mat vs join_where | 73.9714 | 94.1146 | 1.2723 | false |
| 20000 | cte_joined_outer_notmat vs join_where | 73.9714 | 70.8516 | 0.9578 | true |
| 20000 | comma_join vs join_where | 73.9714 | 69.9753 | 0.9460 | true |
| 100000 | join_on vs join_where | 89.9328 | 100.1651 | 1.1138 | true |
| 100000 | filtered_derived vs join_where | 89.9328 | 94.9401 | 1.0557 | true |
| 100000 | cte_filtered vs join_where | 89.9328 | 93.8807 | 1.0439 | true |
| 100000 | cte_filtered_mat vs join_where | 89.9328 | 101.3366 | 1.1268 | false |
| 100000 | cte_filtered_notmat vs join_where | 89.9328 | 94.8852 | 1.0551 | true |
| 100000 | cte_joined_outer vs join_where | 89.9328 | 96.6752 | 1.0750 | true |
| 100000 | cte_joined_outer_mat vs join_where | 89.9328 | 101.4762 | 1.1284 | false |
| 100000 | cte_joined_outer_notmat vs join_where | 89.9328 | 93.4475 | 1.0391 | true |
| 100000 | comma_join vs join_where | 89.9328 | 94.6713 | 1.0527 | true |
| 180000 | join_on vs join_where | 119.3936 | 125.0222 | 1.0471 | true |
| 180000 | filtered_derived vs join_where | 119.3936 | 117.8714 | 0.9873 | true |
| 180000 | cte_filtered vs join_where | 119.3936 | 113.9327 | 0.9543 | true |
| 180000 | cte_filtered_mat vs join_where | 119.3936 | 130.7327 | 1.0950 | false |
| 180000 | cte_filtered_notmat vs join_where | 119.3936 | 118.1213 | 0.9893 | true |
| 180000 | cte_joined_outer vs join_where | 119.3936 | 122.3637 | 1.0249 | true |
| 180000 | cte_joined_outer_mat vs join_where | 119.3936 | 108.2137 | 0.9064 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 119.3936 | 116.5785 | 0.9764 | true |
| 180000 | comma_join vs join_where | 119.3936 | 114.8336 | 0.9618 | true |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_filtered_mat vs cte_filtered | 55.3352 | 55.1612 | 0.9969 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_filtered_notmat vs cte_filtered | 55.3352 | 56.1700 | 1.0151 | true | policy=not_materialized; observed_cte_scan=false |
| 20 | cte_joined_outer_mat vs cte_joined_outer | 57.1734 | 86.1778 | 1.5073 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 57.1734 | 54.7682 | 0.9579 | true | policy=not_materialized; observed_cte_scan=false |
| 200 | cte_filtered_mat vs cte_filtered | 58.5707 | 55.2458 | 0.9432 | false | policy=materialized; observed_cte_scan=true |
| 200 | cte_filtered_notmat vs cte_filtered | 58.5707 | 56.9247 | 0.9719 | true | policy=not_materialized; observed_cte_scan=false |
| 200 | cte_joined_outer_mat vs cte_joined_outer | 58.7040 | 89.9664 | 1.5325 | false | policy=materialized; observed_cte_scan=true |
| 200 | cte_joined_outer_notmat vs cte_joined_outer | 58.7040 | 57.9382 | 0.9870 | true | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_filtered_mat vs cte_filtered | 57.1560 | 57.4772 | 1.0056 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_filtered_notmat vs cte_filtered | 57.1560 | 58.5656 | 1.0247 | true | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 56.2968 | 91.8593 | 1.6317 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 56.2968 | 57.1141 | 1.0145 | true | policy=not_materialized; observed_cte_scan=false |
| 20000 | cte_filtered_mat vs cte_filtered | 70.3850 | 71.3000 | 1.0130 | false | policy=materialized; observed_cte_scan=true |
| 20000 | cte_filtered_notmat vs cte_filtered | 70.3850 | 70.5793 | 1.0028 | true | policy=not_materialized; observed_cte_scan=false |
| 20000 | cte_joined_outer_mat vs cte_joined_outer | 68.1864 | 94.1146 | 1.3803 | false | policy=materialized; observed_cte_scan=true |
| 20000 | cte_joined_outer_notmat vs cte_joined_outer | 68.1864 | 70.8516 | 1.0391 | true | policy=not_materialized; observed_cte_scan=false |
| 100000 | cte_filtered_mat vs cte_filtered | 93.8807 | 101.3366 | 1.0794 | false | policy=materialized; observed_cte_scan=true |
| 100000 | cte_filtered_notmat vs cte_filtered | 93.8807 | 94.8852 | 1.0107 | true | policy=not_materialized; observed_cte_scan=false |
| 100000 | cte_joined_outer_mat vs cte_joined_outer | 96.6752 | 101.4762 | 1.0497 | false | policy=materialized; observed_cte_scan=true |
| 100000 | cte_joined_outer_notmat vs cte_joined_outer | 96.6752 | 93.4475 | 0.9666 | true | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_filtered_mat vs cte_filtered | 113.9327 | 130.7327 | 1.1475 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_filtered_notmat vs cte_filtered | 113.9327 | 118.1213 | 1.0368 | true | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 122.3637 | 108.2137 | 0.8844 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 122.3637 | 116.5785 | 0.9527 | true | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 60.8007 | 1 | `bf773e2e40365851` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 59.6535 | 1 | `ff7b5ad359874144` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | 56.9339 | 2 | `4a04183093a6d236` | policy=default; observed_cte_scan=true | valid |
| default | 2 | 100% | 57.8339 | 2 | `4a04183093a6d236` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 1% | 60.2149 | 4 | `f9f0719a4d742bb9` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 56.8766 | 4 | `f9f0719a4d742bb9` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | 56.3602 | 8 | `6c203246b5dd2c10` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 100% | 55.3011 | 8 | `6c203246b5dd2c10` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 1% | 58.1748 | 16 | `cc27dd85a1f4ec9b` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 59.5489 | 16 | `cc27dd85a1f4ec9b` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 55.9581 | 1 | `3b3f73bb7eb53e1f` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 57.8953 | 1 | `3b3f73bb7eb53e1f` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | 56.3527 | 2 | `4a04183093a6d236` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 100% | 56.6137 | 2 | `4a04183093a6d236` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 1% | 55.3893 | 4 | `f9f0719a4d742bb9` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 57.5499 | 4 | `f9f0719a4d742bb9` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | 58.1511 | 8 | `6c203246b5dd2c10` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 100% | 58.1307 | 8 | `6c203246b5dd2c10` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 1% | 59.2150 | 16 | `cc27dd85a1f4ec9b` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 58.3002 | 16 | `cc27dd85a1f4ec9b` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 58.3730 | 1 | `bf773e2e40365851` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 57.6320 | 1 | `ff7b5ad359874144` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | 101.3628 | 2 | `c3b4d26e971e009a` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 100% | 96.0695 | 2 | `6feeba9fc23b7461` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 1% | 150.6242 | 4 | `3bbc29b17b9460f9` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 153.2996 | 4 | `211f12963d9e68cf` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | 290.6342 | 8 | `4dc778394bad8168` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 100% | 295.1429 | 8 | `f8a70902ada8c64d` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 1% | 550.9748 | 16 | `dedad8f633c2e2dd` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 589.9626 | 16 | `538c32ce5f76c888` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 0.9204 | false |
| notmat | 1 | 1% | 0.9601 | true |
| mat | 1 | 100% | 0.9705 | false |
| notmat | 1 | 100% | 0.9661 | true |
| mat | 2 | 1% | 0.9898 | true |
| notmat | 2 | 1% | 1.7804 | false |
| mat | 2 | 100% | 0.9789 | true |
| notmat | 2 | 100% | 1.6611 | false |
| mat | 4 | 1% | 0.9199 | true |
| notmat | 4 | 1% | 2.5014 | false |
| mat | 4 | 100% | 1.0118 | true |
| notmat | 4 | 100% | 2.6953 | false |
| mat | 8 | 1% | 1.0318 | true |
| notmat | 8 | 1% | 5.1567 | false |
| mat | 8 | 100% | 1.0512 | true |
| notmat | 8 | 100% | 5.3370 | false |
| mat | 16 | 1% | 1.0179 | true |
| notmat | 16 | 1% | 9.4710 | false |
| mat | 16 | 100% | 0.9790 | true |
| notmat | 16 | 100% | 9.9072 | false |

## Repeated-design position control

| first median ms | last median ms | ratio (last/first) | same plan |
|---|---|---|---|
| 71.299964 | 96.039652 | 1.3470 | true |

## Execution order

- order seed `42569480`, 5 interleaved rounds
  1. family=bag;stmt=cte_filtered_mat;threshold=20000
  2. family=reuse;policy=default;refs=16;consumer=100%
  3. family=bag;stmt=cte_filtered_notmat;threshold=200
  4. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  5. family=reuse;policy=notmat;refs=1;consumer=100%
  6. family=bag;stmt=cte_joined_outer;threshold=2000
  7. family=bag;stmt=join_on;threshold=100000
  8. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  9. family=reuse;policy=notmat;refs=8;consumer=100%
  10. family=bag;stmt=comma_join;threshold=2000
  11. family=bag;stmt=cte_filtered_notmat;threshold=100000
  12. family=bag;stmt=cte_filtered_mat;threshold=2000
  13. family=bag;stmt=comma_join;threshold=200
  14. family=bag;stmt=filtered_derived;threshold=180000
  15. family=reuse;policy=default;refs=1;consumer=1%
  16. family=bag;stmt=filtered_derived;threshold=20000
  17. family=reuse;policy=default;refs=1;consumer=100%
  18. family=bag;stmt=comma_join;threshold=20000
  19. family=bag;stmt=filtered_derived;threshold=200
  20. family=reuse;policy=default;refs=8;consumer=100%
  21. family=bag;stmt=cte_joined_outer_mat;threshold=20
  22. family=reuse;policy=mat;refs=1;consumer=100%
  23. family=bag;stmt=cte_filtered;threshold=20
  24. family=reuse;policy=mat;refs=1;consumer=1%
  25. family=reuse;policy=default;refs=4;consumer=100%
  26. family=reuse;policy=notmat;refs=8;consumer=1%
  27. family=bag;stmt=join_on;threshold=180000
  28. family=bag;stmt=cte_filtered;threshold=100000
  29. family=bag;stmt=cte_filtered_mat;threshold=200
  30. family=bag;stmt=comma_join;threshold=20
  31. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  32. family=bag;stmt=join_on;threshold=200
  33. family=reuse;policy=mat;refs=2;consumer=1%
  34. family=bag;stmt=cte_joined_outer_notmat;threshold=20000
  35. family=bag;stmt=cte_filtered_mat;threshold=100000
  36. family=bag;stmt=cte_joined_outer_mat;threshold=200
  37. family=reuse;policy=mat;refs=16;consumer=100%
  38. family=reuse;policy=mat;refs=8;consumer=100%
  39. family=bag;stmt=join_where;threshold=180000
  40. family=bag;stmt=cte_joined_outer_mat;threshold=100000
  41. family=bag;stmt=cte_filtered;threshold=180000
  42. family=bag;stmt=cte_joined_outer;threshold=200
  43. family=bag;stmt=join_where;threshold=100000
  44. family=reuse;policy=notmat;refs=2;consumer=100%
  45. family=reuse;policy=notmat;refs=1;consumer=1%
  46. family=bag;stmt=cte_filtered_mat;threshold=20
  47. family=reuse;policy=mat;refs=4;consumer=1%
  48. family=reuse;policy=default;refs=8;consumer=1%
  49. family=bag;stmt=cte_joined_outer;threshold=20000
  50. family=reuse;policy=notmat;refs=16;consumer=1%
  51. family=bag;stmt=cte_filtered_notmat;threshold=2000
  52. family=bag;stmt=cte_joined_outer;threshold=20
  53. family=bag;stmt=cte_joined_outer_mat;threshold=20000
  54. family=bag;stmt=join_on;threshold=2000
  55. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  56. family=bag;stmt=filtered_derived;threshold=100000
  57. family=reuse;policy=mat;refs=2;consumer=100%
  58. family=bag;stmt=cte_joined_outer_notmat;threshold=100000
  59. family=bag;stmt=join_where;threshold=2000
  60. family=reuse;policy=mat;refs=16;consumer=1%
  61. family=reuse;policy=default;refs=2;consumer=1%
  62. family=reuse;policy=mat;refs=8;consumer=1%
  63. family=bag;stmt=comma_join;threshold=100000
  64. family=bag;stmt=cte_joined_outer;threshold=100000
  65. family=bag;stmt=join_where;threshold=20000
  66. family=bag;stmt=cte_joined_outer;threshold=180000
  67. family=bag;stmt=cte_filtered_notmat;threshold=20000
  68. family=reuse;policy=notmat;refs=4;consumer=100%
  69. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  70. family=reuse;policy=notmat;refs=2;consumer=1%
  71. family=bag;stmt=cte_filtered_notmat;threshold=180000
  72. family=bag;stmt=join_where;threshold=200
  73. family=reuse;policy=default;refs=16;consumer=1%
  74. family=reuse;policy=mat;refs=4;consumer=100%
  75. family=bag;stmt=cte_filtered;threshold=2000
  76. family=reuse;policy=notmat;refs=4;consumer=1%
  77. family=bag;stmt=cte_filtered_mat;threshold=180000
  78. family=bag;stmt=filtered_derived;threshold=20
  79. family=bag;stmt=comma_join;threshold=180000
  80. family=bag;stmt=cte_filtered;threshold=20000
  81. family=bag;stmt=cte_filtered_notmat;threshold=20
  82. family=reuse;policy=notmat;refs=16;consumer=100%
  83. family=reuse;policy=default;refs=4;consumer=1%
  84. family=reuse;policy=default;refs=2;consumer=100%
  85. family=bag;stmt=cte_joined_outer_notmat;threshold=200
  86. family=bag;stmt=filtered_derived;threshold=2000
  87. family=bag;stmt=join_where;threshold=20
  88. family=bag;stmt=join_on;threshold=20
  89. family=bag;stmt=cte_filtered;threshold=200
  90. family=bag;stmt=join_on;threshold=20000
  91. family=bag;stmt=cte_filtered_mat;threshold=20000;position=last

## Live negative controls

| control | pretended contract | detected | rows | note |
|---|---|---|---|---|
| x01_optional_where | sel_join_optional | true | 1 |  |
| x02_join_limit_before | sel_join_bag | true | 7 |  |

## Excluded cells

None: every cell is valid.

## Plan signatures

| axis | statement | status | structure signature | policy note |
|---|---|---|---|---|
| family=bag;stmt=comma_join;threshold=100000 | j06_comma_join | valid | `346cb9e7edfc6298` |  |
| family=bag;stmt=comma_join;threshold=180000 | j06_comma_join | valid | `c393d9dbc6d22445` |  |
| family=bag;stmt=comma_join;threshold=20 | j06_comma_join | valid | `3aa28fab69c4084c` |  |
| family=bag;stmt=comma_join;threshold=200 | j06_comma_join | valid | `62e720d49063502a` |  |
| family=bag;stmt=comma_join;threshold=2000 | j06_comma_join | valid | `24b87f8e6484e319` |  |
| family=bag;stmt=comma_join;threshold=20000 | j06_comma_join | valid | `dbf5856d383eb7fb` |  |
| family=bag;stmt=cte_filtered;threshold=100000 | j04_cte_filtered | valid | `346cb9e7edfc6298` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=180000 | j04_cte_filtered | valid | `c393d9dbc6d22445` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=20 | j04_cte_filtered | valid | `3aa28fab69c4084c` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=200 | j04_cte_filtered | valid | `62e720d49063502a` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=2000 | j04_cte_filtered | valid | `24b87f8e6484e319` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=20000 | j04_cte_filtered | valid | `dbf5856d383eb7fb` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_mat;threshold=100000 | j04m_cte_filtered_mat | valid | `aadbbae70868c951` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=180000 | j04m_cte_filtered_mat | valid | `a17d828e8e6039bd` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=20 | j04m_cte_filtered_mat | valid | `64c934f5018e9c24` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=200 | j04m_cte_filtered_mat | valid | `445b3a9a7668a3e3` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=2000 | j04m_cte_filtered_mat | valid | `f302eb3fe01479ee` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=20000 | j04m_cte_filtered_mat | valid | `0890f84d077d4624` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=20000;position=last | j04m_cte_filtered_mat | valid | `0890f84d077d4624` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_notmat;threshold=100000 | j04n_cte_filtered_notmat | valid | `346cb9e7edfc6298` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=180000 | j04n_cte_filtered_notmat | valid | `c393d9dbc6d22445` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=20 | j04n_cte_filtered_notmat | valid | `3aa28fab69c4084c` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=200 | j04n_cte_filtered_notmat | valid | `62e720d49063502a` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=2000 | j04n_cte_filtered_notmat | valid | `24b87f8e6484e319` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=20000 | j04n_cte_filtered_notmat | valid | `dbf5856d383eb7fb` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=100000 | j05_cte_joined_outer | valid | `346cb9e7edfc6298` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=180000 | j05_cte_joined_outer | valid | `c393d9dbc6d22445` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20 | j05_cte_joined_outer | valid | `3aa28fab69c4084c` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=200 | j05_cte_joined_outer | valid | `62e720d49063502a` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=2000 | j05_cte_joined_outer | valid | `24b87f8e6484e319` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20000 | j05_cte_joined_outer | valid | `dbf5856d383eb7fb` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_mat;threshold=100000 | j05m_cte_joined_outer_mat | valid | `9260d16a40769f90` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | j05m_cte_joined_outer_mat | valid | `9260d16a40769f90` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | j05m_cte_joined_outer_mat | valid | `9260d16a40769f90` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=200 | j05m_cte_joined_outer_mat | valid | `9260d16a40769f90` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | j05m_cte_joined_outer_mat | valid | `9260d16a40769f90` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20000 | j05m_cte_joined_outer_mat | valid | `9260d16a40769f90` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_notmat;threshold=100000 | j05n_cte_joined_outer_notmat | valid | `346cb9e7edfc6298` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | j05n_cte_joined_outer_notmat | valid | `c393d9dbc6d22445` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | j05n_cte_joined_outer_notmat | valid | `3aa28fab69c4084c` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=200 | j05n_cte_joined_outer_notmat | valid | `62e720d49063502a` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | j05n_cte_joined_outer_notmat | valid | `24b87f8e6484e319` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20000 | j05n_cte_joined_outer_notmat | valid | `dbf5856d383eb7fb` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=filtered_derived;threshold=100000 | j03_filtered_derived | valid | `346cb9e7edfc6298` |  |
| family=bag;stmt=filtered_derived;threshold=180000 | j03_filtered_derived | valid | `c393d9dbc6d22445` |  |
| family=bag;stmt=filtered_derived;threshold=20 | j03_filtered_derived | valid | `3aa28fab69c4084c` |  |
| family=bag;stmt=filtered_derived;threshold=200 | j03_filtered_derived | valid | `62e720d49063502a` |  |
| family=bag;stmt=filtered_derived;threshold=2000 | j03_filtered_derived | valid | `24b87f8e6484e319` |  |
| family=bag;stmt=filtered_derived;threshold=20000 | j03_filtered_derived | valid | `dbf5856d383eb7fb` |  |
| family=bag;stmt=join_on;threshold=100000 | j02_join_on | valid | `346cb9e7edfc6298` |  |
| family=bag;stmt=join_on;threshold=180000 | j02_join_on | valid | `c393d9dbc6d22445` |  |
| family=bag;stmt=join_on;threshold=20 | j02_join_on | valid | `3aa28fab69c4084c` |  |
| family=bag;stmt=join_on;threshold=200 | j02_join_on | valid | `62e720d49063502a` |  |
| family=bag;stmt=join_on;threshold=2000 | j02_join_on | valid | `24b87f8e6484e319` |  |
| family=bag;stmt=join_on;threshold=20000 | j02_join_on | valid | `dbf5856d383eb7fb` |  |
| family=bag;stmt=join_where;threshold=100000 | j01_join_where | valid | `346cb9e7edfc6298` |  |
| family=bag;stmt=join_where;threshold=180000 | j01_join_where | valid | `c393d9dbc6d22445` |  |
| family=bag;stmt=join_where;threshold=20 | j01_join_where | valid | `3aa28fab69c4084c` |  |
| family=bag;stmt=join_where;threshold=200 | j01_join_where | valid | `62e720d49063502a` |  |
| family=bag;stmt=join_where;threshold=2000 | j01_join_where | valid | `24b87f8e6484e319` |  |
| family=bag;stmt=join_where;threshold=20000 | j01_join_where | valid | `dbf5856d383eb7fb` |  |
| family=reuse;policy=default;refs=16;consumer=1% | r_default_n16_c1 | valid | `cc27dd85a1f4ec9b` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=16;consumer=100% | r_default_n16_c100 | valid | `cc27dd85a1f4ec9b` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=1;consumer=1% | r_default_n1_c1 | valid | `bf773e2e40365851` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=1;consumer=100% | r_default_n1_c100 | valid | `ff7b5ad359874144` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=2;consumer=1% | r_default_n2_c1 | valid | `4a04183093a6d236` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=2;consumer=100% | r_default_n2_c100 | valid | `4a04183093a6d236` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=1% | r_default_n4_c1 | valid | `f9f0719a4d742bb9` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=100% | r_default_n4_c100 | valid | `f9f0719a4d742bb9` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=8;consumer=1% | r_default_n8_c1 | valid | `6c203246b5dd2c10` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=8;consumer=100% | r_default_n8_c100 | valid | `6c203246b5dd2c10` | policy=default; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=1% | r_mat_n16_c1 | valid | `cc27dd85a1f4ec9b` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=100% | r_mat_n16_c100 | valid | `cc27dd85a1f4ec9b` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=1% | r_mat_n1_c1 | valid | `3b3f73bb7eb53e1f` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=100% | r_mat_n1_c100 | valid | `3b3f73bb7eb53e1f` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=2;consumer=1% | r_mat_n2_c1 | valid | `4a04183093a6d236` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=2;consumer=100% | r_mat_n2_c100 | valid | `4a04183093a6d236` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=1% | r_mat_n4_c1 | valid | `f9f0719a4d742bb9` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=100% | r_mat_n4_c100 | valid | `f9f0719a4d742bb9` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=8;consumer=1% | r_mat_n8_c1 | valid | `6c203246b5dd2c10` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=8;consumer=100% | r_mat_n8_c100 | valid | `6c203246b5dd2c10` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=notmat;refs=16;consumer=1% | r_notmat_n16_c1 | valid | `dedad8f633c2e2dd` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `538c32ce5f76c888` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `bf773e2e40365851` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `ff7b5ad359874144` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=2;consumer=1% | r_notmat_n2_c1 | valid | `c3b4d26e971e009a` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=2;consumer=100% | r_notmat_n2_c100 | valid | `6feeba9fc23b7461` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `3bbc29b17b9460f9` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `211f12963d9e68cf` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=8;consumer=1% | r_notmat_n8_c1 | valid | `4dc778394bad8168` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=8;consumer=100% | r_notmat_n8_c100 | valid | `f8a70902ada8c64d` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- one host, single client: no independent hardware, real-network or equal-total-budget comparison
- no concurrent writers: the isolated read contract only; a matched blended read/write phase is a recorded gap
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
