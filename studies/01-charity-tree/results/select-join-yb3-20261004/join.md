# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 91 cells, 91 valid, engine `yugabyte`, topology `single`, scale `cte`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-yb3-20261004`
- run tag: `run/study-01/select-join-yb3-20261004`
- environment: `host-zenbook-ux5406sa`
- repo commit: `7e148174b5211225f6b4f243fc02e0de417a0586`
- engine version: `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit`
- engine image: `localhost/charitybench-select:select-join-yb3-20261004` @ `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`
- inputs digest: `2d8394b9cc2a9a0f649367486f9d994fcfc1538e3e7d151bc75a1e8febe82e46`
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
| 20 (0.01%) | 50.2806 | 55.5802 | 54.3175 | 48.9501 | 47.0736 | 50.4773 | 50.3063 | 76.8114 | 48.8650 | 51.8888 |
| 200 (0.1%) | 48.5288 | 51.3271 | 49.4475 | 48.8216 | 48.6444 | 49.4430 | 50.1562 | 76.8135 | 51.3568 | 51.5427 |
| 2000 (1%) | 49.5256 | 50.8117 | 51.0314 | 50.4312 | 48.3419 | 49.0466 | 51.4149 | 75.9808 | 51.4716 | 50.2710 |
| 20000 (10%) | 62.0317 | 60.1219 | 61.3544 | 60.1950 | 61.9860 | 61.4589 | 62.0711 | 81.2126 | 61.7382 | 62.1412 |
| 100000 (50%) | 83.3312 | 82.6664 | 82.6168 | 82.5497 | 87.8927 | 83.2714 | 81.3736 | 88.2493 | 83.5110 | 82.9680 |
| 180000 (90%) | 105.8799 | 102.6127 | 102.5922 | 102.7894 | 108.7353 | 104.2315 | 103.9512 | 98.1551 | 105.3809 | 104.8518 |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | join_on vs join_where | 50.2806 | 55.5802 | 1.1054 | true |
| 20 | filtered_derived vs join_where | 50.2806 | 54.3175 | 1.0803 | true |
| 20 | cte_filtered vs join_where | 50.2806 | 48.9501 | 0.9735 | true |
| 20 | cte_filtered_mat vs join_where | 50.2806 | 47.0736 | 0.9362 | false |
| 20 | cte_filtered_notmat vs join_where | 50.2806 | 50.4773 | 1.0039 | true |
| 20 | cte_joined_outer vs join_where | 50.2806 | 50.3063 | 1.0005 | true |
| 20 | cte_joined_outer_mat vs join_where | 50.2806 | 76.8114 | 1.5277 | false |
| 20 | cte_joined_outer_notmat vs join_where | 50.2806 | 48.8650 | 0.9718 | true |
| 20 | comma_join vs join_where | 50.2806 | 51.8888 | 1.0320 | true |
| 200 | join_on vs join_where | 48.5288 | 51.3271 | 1.0577 | true |
| 200 | filtered_derived vs join_where | 48.5288 | 49.4475 | 1.0189 | true |
| 200 | cte_filtered vs join_where | 48.5288 | 48.8216 | 1.0060 | true |
| 200 | cte_filtered_mat vs join_where | 48.5288 | 48.6444 | 1.0024 | false |
| 200 | cte_filtered_notmat vs join_where | 48.5288 | 49.4430 | 1.0188 | true |
| 200 | cte_joined_outer vs join_where | 48.5288 | 50.1562 | 1.0335 | true |
| 200 | cte_joined_outer_mat vs join_where | 48.5288 | 76.8135 | 1.5828 | false |
| 200 | cte_joined_outer_notmat vs join_where | 48.5288 | 51.3568 | 1.0583 | true |
| 200 | comma_join vs join_where | 48.5288 | 51.5427 | 1.0621 | true |
| 2000 | join_on vs join_where | 49.5256 | 50.8117 | 1.0260 | true |
| 2000 | filtered_derived vs join_where | 49.5256 | 51.0314 | 1.0304 | true |
| 2000 | cte_filtered vs join_where | 49.5256 | 50.4312 | 1.0183 | true |
| 2000 | cte_filtered_mat vs join_where | 49.5256 | 48.3419 | 0.9761 | false |
| 2000 | cte_filtered_notmat vs join_where | 49.5256 | 49.0466 | 0.9903 | true |
| 2000 | cte_joined_outer vs join_where | 49.5256 | 51.4149 | 1.0381 | true |
| 2000 | cte_joined_outer_mat vs join_where | 49.5256 | 75.9808 | 1.5342 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 49.5256 | 51.4716 | 1.0393 | true |
| 2000 | comma_join vs join_where | 49.5256 | 50.2710 | 1.0150 | true |
| 20000 | join_on vs join_where | 62.0317 | 60.1219 | 0.9692 | true |
| 20000 | filtered_derived vs join_where | 62.0317 | 61.3544 | 0.9891 | true |
| 20000 | cte_filtered vs join_where | 62.0317 | 60.1950 | 0.9704 | true |
| 20000 | cte_filtered_mat vs join_where | 62.0317 | 61.9860 | 0.9993 | false |
| 20000 | cte_filtered_notmat vs join_where | 62.0317 | 61.4589 | 0.9908 | true |
| 20000 | cte_joined_outer vs join_where | 62.0317 | 62.0711 | 1.0006 | true |
| 20000 | cte_joined_outer_mat vs join_where | 62.0317 | 81.2126 | 1.3092 | false |
| 20000 | cte_joined_outer_notmat vs join_where | 62.0317 | 61.7382 | 0.9953 | true |
| 20000 | comma_join vs join_where | 62.0317 | 62.1412 | 1.0018 | true |
| 100000 | join_on vs join_where | 83.3312 | 82.6664 | 0.9920 | true |
| 100000 | filtered_derived vs join_where | 83.3312 | 82.6168 | 0.9914 | true |
| 100000 | cte_filtered vs join_where | 83.3312 | 82.5497 | 0.9906 | true |
| 100000 | cte_filtered_mat vs join_where | 83.3312 | 87.8927 | 1.0547 | false |
| 100000 | cte_filtered_notmat vs join_where | 83.3312 | 83.2714 | 0.9993 | true |
| 100000 | cte_joined_outer vs join_where | 83.3312 | 81.3736 | 0.9765 | true |
| 100000 | cte_joined_outer_mat vs join_where | 83.3312 | 88.2493 | 1.0590 | false |
| 100000 | cte_joined_outer_notmat vs join_where | 83.3312 | 83.5110 | 1.0022 | true |
| 100000 | comma_join vs join_where | 83.3312 | 82.9680 | 0.9956 | true |
| 180000 | join_on vs join_where | 105.8799 | 102.6127 | 0.9691 | true |
| 180000 | filtered_derived vs join_where | 105.8799 | 102.5922 | 0.9689 | true |
| 180000 | cte_filtered vs join_where | 105.8799 | 102.7894 | 0.9708 | true |
| 180000 | cte_filtered_mat vs join_where | 105.8799 | 108.7353 | 1.0270 | false |
| 180000 | cte_filtered_notmat vs join_where | 105.8799 | 104.2315 | 0.9844 | true |
| 180000 | cte_joined_outer vs join_where | 105.8799 | 103.9512 | 0.9818 | true |
| 180000 | cte_joined_outer_mat vs join_where | 105.8799 | 98.1551 | 0.9270 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 105.8799 | 105.3809 | 0.9953 | true |
| 180000 | comma_join vs join_where | 105.8799 | 104.8518 | 0.9903 | true |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_filtered_mat vs cte_filtered | 48.9501 | 47.0736 | 0.9617 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_filtered_notmat vs cte_filtered | 48.9501 | 50.4773 | 1.0312 | true | policy=not_materialized; observed_cte_scan=false |
| 20 | cte_joined_outer_mat vs cte_joined_outer | 50.3063 | 76.8114 | 1.5269 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 50.3063 | 48.8650 | 0.9713 | true | policy=not_materialized; observed_cte_scan=false |
| 200 | cte_filtered_mat vs cte_filtered | 48.8216 | 48.6444 | 0.9964 | false | policy=materialized; observed_cte_scan=true |
| 200 | cte_filtered_notmat vs cte_filtered | 48.8216 | 49.4430 | 1.0127 | true | policy=not_materialized; observed_cte_scan=false |
| 200 | cte_joined_outer_mat vs cte_joined_outer | 50.1562 | 76.8135 | 1.5315 | false | policy=materialized; observed_cte_scan=true |
| 200 | cte_joined_outer_notmat vs cte_joined_outer | 50.1562 | 51.3568 | 1.0239 | true | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_filtered_mat vs cte_filtered | 50.4312 | 48.3419 | 0.9586 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_filtered_notmat vs cte_filtered | 50.4312 | 49.0466 | 0.9725 | true | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 51.4149 | 75.9808 | 1.4778 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 51.4149 | 51.4716 | 1.0011 | true | policy=not_materialized; observed_cte_scan=false |
| 20000 | cte_filtered_mat vs cte_filtered | 60.1950 | 61.9860 | 1.0298 | false | policy=materialized; observed_cte_scan=true |
| 20000 | cte_filtered_notmat vs cte_filtered | 60.1950 | 61.4589 | 1.0210 | true | policy=not_materialized; observed_cte_scan=false |
| 20000 | cte_joined_outer_mat vs cte_joined_outer | 62.0711 | 81.2126 | 1.3084 | false | policy=materialized; observed_cte_scan=true |
| 20000 | cte_joined_outer_notmat vs cte_joined_outer | 62.0711 | 61.7382 | 0.9946 | true | policy=not_materialized; observed_cte_scan=false |
| 100000 | cte_filtered_mat vs cte_filtered | 82.5497 | 87.8927 | 1.0647 | false | policy=materialized; observed_cte_scan=true |
| 100000 | cte_filtered_notmat vs cte_filtered | 82.5497 | 83.2714 | 1.0087 | true | policy=not_materialized; observed_cte_scan=false |
| 100000 | cte_joined_outer_mat vs cte_joined_outer | 81.3736 | 88.2493 | 1.0845 | false | policy=materialized; observed_cte_scan=true |
| 100000 | cte_joined_outer_notmat vs cte_joined_outer | 81.3736 | 83.5110 | 1.0263 | true | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_filtered_mat vs cte_filtered | 102.7894 | 108.7353 | 1.0578 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_filtered_notmat vs cte_filtered | 102.7894 | 104.2315 | 1.0140 | true | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 103.9512 | 98.1551 | 0.9442 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 103.9512 | 105.3809 | 1.0138 | true | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 54.5585 | 1 | `bf773e2e40365851` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 52.2319 | 1 | `ff7b5ad359874144` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | 52.5271 | 2 | `4a04183093a6d236` | policy=default; observed_cte_scan=true | valid |
| default | 2 | 100% | 48.3360 | 2 | `4a04183093a6d236` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 1% | 51.7899 | 4 | `f9f0719a4d742bb9` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 51.9626 | 4 | `f9f0719a4d742bb9` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | 49.7483 | 8 | `6c203246b5dd2c10` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 100% | 50.9801 | 8 | `6c203246b5dd2c10` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 1% | 53.9796 | 16 | `cc27dd85a1f4ec9b` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 52.8587 | 16 | `cc27dd85a1f4ec9b` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 50.0198 | 1 | `3b3f73bb7eb53e1f` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 50.8898 | 1 | `3b3f73bb7eb53e1f` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | 51.3200 | 2 | `4a04183093a6d236` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 100% | 49.3620 | 2 | `4a04183093a6d236` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 1% | 51.0538 | 4 | `f9f0719a4d742bb9` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 51.9624 | 4 | `f9f0719a4d742bb9` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | 50.6843 | 8 | `6c203246b5dd2c10` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 100% | 49.7646 | 8 | `6c203246b5dd2c10` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 1% | 50.0097 | 16 | `cc27dd85a1f4ec9b` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 53.1911 | 16 | `cc27dd85a1f4ec9b` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 51.9662 | 1 | `bf773e2e40365851` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 53.0314 | 1 | `ff7b5ad359874144` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | 88.3428 | 2 | `c3b4d26e971e009a` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 100% | 87.3859 | 2 | `6feeba9fc23b7461` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 1% | 127.3600 | 4 | `3bbc29b17b9460f9` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 124.4023 | 4 | `211f12963d9e68cf` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | 234.9841 | 8 | `4dc778394bad8168` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 100% | 269.7885 | 8 | `f8a70902ada8c64d` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 1% | 496.9699 | 16 | `dedad8f633c2e2dd` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 507.6440 | 16 | `538c32ce5f76c888` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 0.9168 | false |
| notmat | 1 | 1% | 0.9525 | true |
| mat | 1 | 100% | 0.9743 | false |
| notmat | 1 | 100% | 1.0153 | true |
| mat | 2 | 1% | 0.9770 | true |
| notmat | 2 | 1% | 1.6819 | false |
| mat | 2 | 100% | 1.0212 | true |
| notmat | 2 | 100% | 1.8079 | false |
| mat | 4 | 1% | 0.9858 | true |
| notmat | 4 | 1% | 2.4592 | false |
| mat | 4 | 100% | 1.0000 | true |
| notmat | 4 | 100% | 2.3941 | false |
| mat | 8 | 1% | 1.0188 | true |
| notmat | 8 | 1% | 4.7235 | false |
| mat | 8 | 100% | 0.9762 | true |
| notmat | 8 | 100% | 5.2920 | false |
| mat | 16 | 1% | 0.9265 | true |
| notmat | 16 | 1% | 9.2066 | false |
| mat | 16 | 100% | 1.0063 | true |
| notmat | 16 | 100% | 9.6038 | false |

## Repeated-design position control

| first median ms | last median ms | ratio (last/first) | same plan |
|---|---|---|---|
| 98.155142 | 99.160352 | 1.0102 | true |

## Execution order

- order seed `42439153`, 5 interleaved rounds
  1. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  2. family=bag;stmt=join_where;threshold=2000
  3. family=bag;stmt=cte_joined_outer;threshold=2000
  4. family=reuse;policy=notmat;refs=2;consumer=100%
  5. family=bag;stmt=comma_join;threshold=2000
  6. family=bag;stmt=cte_filtered_notmat;threshold=200
  7. family=bag;stmt=cte_filtered;threshold=20
  8. family=bag;stmt=filtered_derived;threshold=180000
  9. family=reuse;policy=notmat;refs=8;consumer=1%
  10. family=bag;stmt=comma_join;threshold=100000
  11. family=bag;stmt=join_on;threshold=2000
  12. family=reuse;policy=default;refs=16;consumer=1%
  13. family=reuse;policy=notmat;refs=8;consumer=100%
  14. family=bag;stmt=cte_joined_outer;threshold=100000
  15. family=reuse;policy=mat;refs=1;consumer=100%
  16. family=bag;stmt=filtered_derived;threshold=200
  17. family=bag;stmt=cte_filtered_mat;threshold=200
  18. family=reuse;policy=notmat;refs=16;consumer=100%
  19. family=bag;stmt=comma_join;threshold=200
  20. family=bag;stmt=join_where;threshold=20
  21. family=reuse;policy=mat;refs=16;consumer=100%
  22. family=bag;stmt=cte_joined_outer_notmat;threshold=100000
  23. family=bag;stmt=filtered_derived;threshold=100000
  24. family=reuse;policy=mat;refs=2;consumer=1%
  25. family=reuse;policy=mat;refs=16;consumer=1%
  26. family=bag;stmt=comma_join;threshold=20000
  27. family=bag;stmt=cte_filtered_mat;threshold=20000
  28. family=reuse;policy=default;refs=4;consumer=100%
  29. family=reuse;policy=mat;refs=2;consumer=100%
  30. family=reuse;policy=notmat;refs=1;consumer=1%
  31. family=reuse;policy=mat;refs=8;consumer=1%
  32. family=bag;stmt=cte_joined_outer_notmat;threshold=200
  33. family=bag;stmt=cte_filtered;threshold=200
  34. family=bag;stmt=join_where;threshold=200
  35. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  36. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  37. family=bag;stmt=cte_filtered_mat;threshold=20
  38. family=reuse;policy=default;refs=8;consumer=1%
  39. family=bag;stmt=cte_filtered;threshold=100000
  40. family=bag;stmt=cte_filtered;threshold=180000
  41. family=reuse;policy=default;refs=1;consumer=100%
  42. family=bag;stmt=cte_filtered;threshold=2000
  43. family=bag;stmt=join_on;threshold=20000
  44. family=bag;stmt=cte_filtered;threshold=20000
  45. family=reuse;policy=notmat;refs=1;consumer=100%
  46. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  47. family=bag;stmt=join_where;threshold=20000
  48. family=bag;stmt=cte_joined_outer_mat;threshold=200
  49. family=reuse;policy=notmat;refs=16;consumer=1%
  50. family=bag;stmt=cte_joined_outer_notmat;threshold=20000
  51. family=bag;stmt=cte_joined_outer;threshold=20000
  52. family=bag;stmt=cte_joined_outer;threshold=200
  53. family=bag;stmt=cte_filtered_notmat;threshold=2000
  54. family=reuse;policy=notmat;refs=4;consumer=1%
  55. family=bag;stmt=comma_join;threshold=20
  56. family=bag;stmt=cte_joined_outer;threshold=180000
  57. family=bag;stmt=cte_filtered_notmat;threshold=100000
  58. family=reuse;policy=mat;refs=4;consumer=1%
  59. family=bag;stmt=filtered_derived;threshold=20000
  60. family=bag;stmt=cte_joined_outer_mat;threshold=20
  61. family=bag;stmt=join_on;threshold=20
  62. family=bag;stmt=cte_filtered_mat;threshold=2000
  63. family=bag;stmt=join_on;threshold=180000
  64. family=reuse;policy=default;refs=16;consumer=100%
  65. family=bag;stmt=join_where;threshold=100000
  66. family=bag;stmt=comma_join;threshold=180000
  67. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  68. family=reuse;policy=default;refs=1;consumer=1%
  69. family=reuse;policy=notmat;refs=2;consumer=1%
  70. family=reuse;policy=notmat;refs=4;consumer=100%
  71. family=bag;stmt=filtered_derived;threshold=20
  72. family=reuse;policy=default;refs=2;consumer=1%
  73. family=bag;stmt=cte_joined_outer_mat;threshold=20000
  74. family=bag;stmt=join_where;threshold=180000
  75. family=reuse;policy=mat;refs=4;consumer=100%
  76. family=bag;stmt=join_on;threshold=100000
  77. family=bag;stmt=cte_joined_outer_mat;threshold=100000
  78. family=bag;stmt=cte_filtered_notmat;threshold=180000
  79. family=bag;stmt=cte_filtered_mat;threshold=100000
  80. family=bag;stmt=join_on;threshold=200
  81. family=reuse;policy=default;refs=8;consumer=100%
  82. family=bag;stmt=cte_filtered_notmat;threshold=20000
  83. family=reuse;policy=default;refs=4;consumer=1%
  84. family=reuse;policy=default;refs=2;consumer=100%
  85. family=bag;stmt=cte_filtered_mat;threshold=180000
  86. family=bag;stmt=filtered_derived;threshold=2000
  87. family=reuse;policy=mat;refs=8;consumer=100%
  88. family=reuse;policy=mat;refs=1;consumer=1%
  89. family=bag;stmt=cte_joined_outer;threshold=20
  90. family=bag;stmt=cte_filtered_notmat;threshold=20
  91. family=bag;stmt=cte_joined_outer_mat;threshold=180000;position=last

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
| family=bag;stmt=cte_joined_outer_mat;threshold=180000;position=last | j05m_cte_joined_outer_mat | valid | `9260d16a40769f90` | policy=materialized; observed_cte_scan=true |
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
