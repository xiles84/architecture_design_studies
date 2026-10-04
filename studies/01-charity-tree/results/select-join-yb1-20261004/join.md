# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 91 cells, 91 valid, engine `yugabyte`, topology `single`, scale `cte`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-yb1-20261004`
- run tag: `run/study-01/select-join-yb1-20261004`
- environment: `host-zenbook-ux5406sa`
- repo commit: `7e148174b5211225f6b4f243fc02e0de417a0586`
- engine version: `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit`
- engine image: `localhost/charitybench-select:select-join-yb1-20261004` @ `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`
- inputs digest: `f650acde1ebb5541080360bca3e467863611e00322d3f580cefdd36ba90e0281`
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
| 20 (0.01%) | 54.1855 | 51.3924 | 51.9190 | 53.4507 | 52.6065 | 55.8474 | 55.8498 | 87.0800 | 56.1114 | 56.1530 |
| 200 (0.1%) | 53.6080 | 52.2678 | 58.7998 | 52.4404 | 53.3720 | 53.3432 | 56.8219 | 88.6169 | 52.0684 | 52.0503 |
| 2000 (1%) | 54.0418 | 58.4189 | 54.4900 | 52.9426 | 54.3065 | 53.7409 | 60.4125 | 87.4509 | 57.2107 | 51.7476 |
| 20000 (10%) | 69.6583 | 63.5956 | 65.5495 | 66.8677 | 70.4617 | 66.4306 | 66.2433 | 86.8973 | 68.3329 | 68.5227 |
| 100000 (50%) | 89.8912 | 86.5841 | 93.6500 | 90.2188 | 97.3519 | 86.9815 | 90.1083 | 95.3078 | 90.2313 | 89.8597 |
| 180000 (90%) | 121.1290 | 111.9913 | 110.6885 | 107.3187 | 114.6110 | 109.1626 | 102.6674 | 104.4287 | 108.8500 | 111.7094 |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | join_on vs join_where | 54.1855 | 51.3924 | 0.9485 | true |
| 20 | filtered_derived vs join_where | 54.1855 | 51.9190 | 0.9582 | true |
| 20 | cte_filtered vs join_where | 54.1855 | 53.4507 | 0.9864 | true |
| 20 | cte_filtered_mat vs join_where | 54.1855 | 52.6065 | 0.9709 | false |
| 20 | cte_filtered_notmat vs join_where | 54.1855 | 55.8474 | 1.0307 | true |
| 20 | cte_joined_outer vs join_where | 54.1855 | 55.8498 | 1.0307 | true |
| 20 | cte_joined_outer_mat vs join_where | 54.1855 | 87.0800 | 1.6071 | false |
| 20 | cte_joined_outer_notmat vs join_where | 54.1855 | 56.1114 | 1.0355 | true |
| 20 | comma_join vs join_where | 54.1855 | 56.1530 | 1.0363 | true |
| 200 | join_on vs join_where | 53.6080 | 52.2678 | 0.9750 | true |
| 200 | filtered_derived vs join_where | 53.6080 | 58.7998 | 1.0968 | true |
| 200 | cte_filtered vs join_where | 53.6080 | 52.4404 | 0.9782 | true |
| 200 | cte_filtered_mat vs join_where | 53.6080 | 53.3720 | 0.9956 | false |
| 200 | cte_filtered_notmat vs join_where | 53.6080 | 53.3432 | 0.9951 | true |
| 200 | cte_joined_outer vs join_where | 53.6080 | 56.8219 | 1.0600 | true |
| 200 | cte_joined_outer_mat vs join_where | 53.6080 | 88.6169 | 1.6531 | false |
| 200 | cte_joined_outer_notmat vs join_where | 53.6080 | 52.0684 | 0.9713 | true |
| 200 | comma_join vs join_where | 53.6080 | 52.0503 | 0.9709 | true |
| 2000 | join_on vs join_where | 54.0418 | 58.4189 | 1.0810 | true |
| 2000 | filtered_derived vs join_where | 54.0418 | 54.4900 | 1.0083 | true |
| 2000 | cte_filtered vs join_where | 54.0418 | 52.9426 | 0.9797 | true |
| 2000 | cte_filtered_mat vs join_where | 54.0418 | 54.3065 | 1.0049 | false |
| 2000 | cte_filtered_notmat vs join_where | 54.0418 | 53.7409 | 0.9944 | true |
| 2000 | cte_joined_outer vs join_where | 54.0418 | 60.4125 | 1.1179 | true |
| 2000 | cte_joined_outer_mat vs join_where | 54.0418 | 87.4509 | 1.6182 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 54.0418 | 57.2107 | 1.0586 | true |
| 2000 | comma_join vs join_where | 54.0418 | 51.7476 | 0.9575 | true |
| 20000 | join_on vs join_where | 69.6583 | 63.5956 | 0.9130 | true |
| 20000 | filtered_derived vs join_where | 69.6583 | 65.5495 | 0.9410 | true |
| 20000 | cte_filtered vs join_where | 69.6583 | 66.8677 | 0.9599 | true |
| 20000 | cte_filtered_mat vs join_where | 69.6583 | 70.4617 | 1.0115 | false |
| 20000 | cte_filtered_notmat vs join_where | 69.6583 | 66.4306 | 0.9537 | true |
| 20000 | cte_joined_outer vs join_where | 69.6583 | 66.2433 | 0.9510 | true |
| 20000 | cte_joined_outer_mat vs join_where | 69.6583 | 86.8973 | 1.2475 | false |
| 20000 | cte_joined_outer_notmat vs join_where | 69.6583 | 68.3329 | 0.9810 | true |
| 20000 | comma_join vs join_where | 69.6583 | 68.5227 | 0.9837 | true |
| 100000 | join_on vs join_where | 89.8912 | 86.5841 | 0.9632 | true |
| 100000 | filtered_derived vs join_where | 89.8912 | 93.6500 | 1.0418 | true |
| 100000 | cte_filtered vs join_where | 89.8912 | 90.2188 | 1.0036 | true |
| 100000 | cte_filtered_mat vs join_where | 89.8912 | 97.3519 | 1.0830 | false |
| 100000 | cte_filtered_notmat vs join_where | 89.8912 | 86.9815 | 0.9676 | true |
| 100000 | cte_joined_outer vs join_where | 89.8912 | 90.1083 | 1.0024 | true |
| 100000 | cte_joined_outer_mat vs join_where | 89.8912 | 95.3078 | 1.0603 | false |
| 100000 | cte_joined_outer_notmat vs join_where | 89.8912 | 90.2313 | 1.0038 | true |
| 100000 | comma_join vs join_where | 89.8912 | 89.8597 | 0.9996 | true |
| 180000 | join_on vs join_where | 121.1290 | 111.9913 | 0.9246 | true |
| 180000 | filtered_derived vs join_where | 121.1290 | 110.6885 | 0.9138 | true |
| 180000 | cte_filtered vs join_where | 121.1290 | 107.3187 | 0.8860 | true |
| 180000 | cte_filtered_mat vs join_where | 121.1290 | 114.6110 | 0.9462 | false |
| 180000 | cte_filtered_notmat vs join_where | 121.1290 | 109.1626 | 0.9012 | true |
| 180000 | cte_joined_outer vs join_where | 121.1290 | 102.6674 | 0.8476 | true |
| 180000 | cte_joined_outer_mat vs join_where | 121.1290 | 104.4287 | 0.8621 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 121.1290 | 108.8500 | 0.8986 | true |
| 180000 | comma_join vs join_where | 121.1290 | 111.7094 | 0.9222 | true |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_filtered_mat vs cte_filtered | 53.4507 | 52.6065 | 0.9842 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_filtered_notmat vs cte_filtered | 53.4507 | 55.8474 | 1.0448 | true | policy=not_materialized; observed_cte_scan=false |
| 20 | cte_joined_outer_mat vs cte_joined_outer | 55.8498 | 87.0800 | 1.5592 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 55.8498 | 56.1114 | 1.0047 | true | policy=not_materialized; observed_cte_scan=false |
| 200 | cte_filtered_mat vs cte_filtered | 52.4404 | 53.3720 | 1.0178 | false | policy=materialized; observed_cte_scan=true |
| 200 | cte_filtered_notmat vs cte_filtered | 52.4404 | 53.3432 | 1.0172 | true | policy=not_materialized; observed_cte_scan=false |
| 200 | cte_joined_outer_mat vs cte_joined_outer | 56.8219 | 88.6169 | 1.5596 | false | policy=materialized; observed_cte_scan=true |
| 200 | cte_joined_outer_notmat vs cte_joined_outer | 56.8219 | 52.0684 | 0.9163 | true | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_filtered_mat vs cte_filtered | 52.9426 | 54.3065 | 1.0258 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_filtered_notmat vs cte_filtered | 52.9426 | 53.7409 | 1.0151 | true | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 60.4125 | 87.4509 | 1.4476 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 60.4125 | 57.2107 | 0.9470 | true | policy=not_materialized; observed_cte_scan=false |
| 20000 | cte_filtered_mat vs cte_filtered | 66.8677 | 70.4617 | 1.0537 | false | policy=materialized; observed_cte_scan=true |
| 20000 | cte_filtered_notmat vs cte_filtered | 66.8677 | 66.4306 | 0.9935 | true | policy=not_materialized; observed_cte_scan=false |
| 20000 | cte_joined_outer_mat vs cte_joined_outer | 66.2433 | 86.8973 | 1.3118 | false | policy=materialized; observed_cte_scan=true |
| 20000 | cte_joined_outer_notmat vs cte_joined_outer | 66.2433 | 68.3329 | 1.0315 | true | policy=not_materialized; observed_cte_scan=false |
| 100000 | cte_filtered_mat vs cte_filtered | 90.2188 | 97.3519 | 1.0791 | false | policy=materialized; observed_cte_scan=true |
| 100000 | cte_filtered_notmat vs cte_filtered | 90.2188 | 86.9815 | 0.9641 | true | policy=not_materialized; observed_cte_scan=false |
| 100000 | cte_joined_outer_mat vs cte_joined_outer | 90.1083 | 95.3078 | 1.0577 | false | policy=materialized; observed_cte_scan=true |
| 100000 | cte_joined_outer_notmat vs cte_joined_outer | 90.1083 | 90.2313 | 1.0014 | true | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_filtered_mat vs cte_filtered | 107.3187 | 114.6110 | 1.0680 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_filtered_notmat vs cte_filtered | 107.3187 | 109.1626 | 1.0172 | true | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 102.6674 | 104.4287 | 1.0172 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 102.6674 | 108.8500 | 1.0602 | true | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 56.3565 | 1 | `bf773e2e40365851` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 55.8016 | 1 | `ff7b5ad359874144` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | 51.8952 | 2 | `4a04183093a6d236` | policy=default; observed_cte_scan=true | valid |
| default | 2 | 100% | 56.1597 | 2 | `4a04183093a6d236` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 1% | 54.1220 | 4 | `f9f0719a4d742bb9` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 52.9323 | 4 | `f9f0719a4d742bb9` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | 55.4008 | 8 | `6c203246b5dd2c10` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 100% | 57.1514 | 8 | `6c203246b5dd2c10` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 1% | 54.8308 | 16 | `cc27dd85a1f4ec9b` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 52.4722 | 16 | `cc27dd85a1f4ec9b` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 51.3329 | 1 | `3b3f73bb7eb53e1f` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 55.8041 | 1 | `3b3f73bb7eb53e1f` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | 55.3564 | 2 | `4a04183093a6d236` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 100% | 55.0306 | 2 | `4a04183093a6d236` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 1% | 51.5940 | 4 | `f9f0719a4d742bb9` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 59.0504 | 4 | `f9f0719a4d742bb9` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | 56.6268 | 8 | `6c203246b5dd2c10` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 100% | 51.4913 | 8 | `6c203246b5dd2c10` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 1% | 54.3485 | 16 | `cc27dd85a1f4ec9b` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 56.4226 | 16 | `cc27dd85a1f4ec9b` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 56.9631 | 1 | `bf773e2e40365851` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 56.8319 | 1 | `ff7b5ad359874144` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | 94.5241 | 2 | `c3b4d26e971e009a` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 100% | 97.2606 | 2 | `6feeba9fc23b7461` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 1% | 157.2088 | 4 | `3bbc29b17b9460f9` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 182.2373 | 4 | `211f12963d9e68cf` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | 283.7898 | 8 | `4dc778394bad8168` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 100% | 265.0696 | 8 | `f8a70902ada8c64d` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 1% | 490.0560 | 16 | `dedad8f633c2e2dd` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 585.2708 | 16 | `538c32ce5f76c888` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 0.9109 | false |
| notmat | 1 | 1% | 1.0108 | true |
| mat | 1 | 100% | 1.0000 | false |
| notmat | 1 | 100% | 1.0185 | true |
| mat | 2 | 1% | 1.0667 | true |
| notmat | 2 | 1% | 1.8214 | false |
| mat | 2 | 100% | 0.9799 | true |
| notmat | 2 | 100% | 1.7319 | false |
| mat | 4 | 1% | 0.9533 | true |
| notmat | 4 | 1% | 2.9047 | false |
| mat | 4 | 100% | 1.1156 | true |
| notmat | 4 | 100% | 3.4428 | false |
| mat | 8 | 1% | 1.0221 | true |
| notmat | 8 | 1% | 5.1225 | false |
| mat | 8 | 100% | 0.9010 | true |
| notmat | 8 | 100% | 4.6380 | false |
| mat | 16 | 1% | 0.9912 | true |
| notmat | 16 | 1% | 8.9376 | false |
| mat | 16 | 100% | 1.0753 | true |
| notmat | 16 | 100% | 11.1539 | false |

## Repeated-design position control

| first median ms | last median ms | ratio (last/first) | same plan |
|---|---|---|---|
| 68.332942 | 58.825523 | 0.8609 | true |

## Execution order

- order seed `42180155`, 5 interleaved rounds
  1. family=bag;stmt=cte_joined_outer_notmat;threshold=20000
  2. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  3. family=bag;stmt=cte_joined_outer_notmat;threshold=100000
  4. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  5. family=reuse;policy=default;refs=1;consumer=100%
  6. family=bag;stmt=cte_joined_outer;threshold=100000
  7. family=bag;stmt=join_on;threshold=20
  8. family=reuse;policy=mat;refs=2;consumer=100%
  9. family=reuse;policy=default;refs=4;consumer=100%
  10. family=bag;stmt=cte_filtered_notmat;threshold=200
  11. family=bag;stmt=cte_joined_outer;threshold=180000
  12. family=reuse;policy=mat;refs=1;consumer=1%
  13. family=bag;stmt=comma_join;threshold=200
  14. family=bag;stmt=join_on;threshold=100000
  15. family=bag;stmt=cte_filtered_mat;threshold=20
  16. family=bag;stmt=join_on;threshold=200
  17. family=bag;stmt=cte_filtered_mat;threshold=180000
  18. family=reuse;policy=notmat;refs=16;consumer=1%
  19. family=reuse;policy=mat;refs=4;consumer=100%
  20. family=bag;stmt=cte_filtered_mat;threshold=2000
  21. family=bag;stmt=cte_filtered_notmat;threshold=20000
  22. family=bag;stmt=filtered_derived;threshold=20
  23. family=bag;stmt=cte_filtered_notmat;threshold=180000
  24. family=bag;stmt=cte_filtered_notmat;threshold=100000
  25. family=reuse;policy=mat;refs=16;consumer=1%
  26. family=bag;stmt=filtered_derived;threshold=180000
  27. family=bag;stmt=filtered_derived;threshold=2000
  28. family=bag;stmt=cte_filtered;threshold=100000
  29. family=bag;stmt=cte_filtered_notmat;threshold=2000
  30. family=bag;stmt=comma_join;threshold=180000
  31. family=bag;stmt=join_where;threshold=100000
  32. family=reuse;policy=default;refs=8;consumer=1%
  33. family=reuse;policy=notmat;refs=8;consumer=1%
  34. family=reuse;policy=notmat;refs=4;consumer=100%
  35. family=reuse;policy=default;refs=2;consumer=100%
  36. family=reuse;policy=mat;refs=1;consumer=100%
  37. family=reuse;policy=default;refs=8;consumer=100%
  38. family=reuse;policy=notmat;refs=2;consumer=1%
  39. family=bag;stmt=cte_joined_outer;threshold=200
  40. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  41. family=bag;stmt=cte_joined_outer;threshold=20
  42. family=reuse;policy=mat;refs=2;consumer=1%
  43. family=bag;stmt=cte_joined_outer_mat;threshold=200
  44. family=bag;stmt=cte_filtered_mat;threshold=20000
  45. family=reuse;policy=notmat;refs=1;consumer=1%
  46. family=reuse;policy=default;refs=16;consumer=1%
  47. family=bag;stmt=cte_joined_outer;threshold=20000
  48. family=bag;stmt=join_where;threshold=20
  49. family=reuse;policy=notmat;refs=1;consumer=100%
  50. family=bag;stmt=join_where;threshold=2000
  51. family=bag;stmt=cte_filtered_mat;threshold=200
  52. family=bag;stmt=comma_join;threshold=100000
  53. family=bag;stmt=cte_filtered;threshold=180000
  54. family=bag;stmt=cte_joined_outer_mat;threshold=20
  55. family=bag;stmt=cte_filtered;threshold=20
  56. family=bag;stmt=cte_filtered;threshold=20000
  57. family=bag;stmt=comma_join;threshold=20000
  58. family=bag;stmt=filtered_derived;threshold=20000
  59. family=reuse;policy=mat;refs=4;consumer=1%
  60. family=reuse;policy=notmat;refs=2;consumer=100%
  61. family=bag;stmt=join_on;threshold=180000
  62. family=bag;stmt=cte_filtered_mat;threshold=100000
  63. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  64. family=bag;stmt=cte_joined_outer_notmat;threshold=200
  65. family=reuse;policy=mat;refs=8;consumer=100%
  66. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  67. family=bag;stmt=comma_join;threshold=2000
  68. family=reuse;policy=notmat;refs=8;consumer=100%
  69. family=bag;stmt=cte_joined_outer;threshold=2000
  70. family=bag;stmt=cte_filtered;threshold=2000
  71. family=reuse;policy=default;refs=1;consumer=1%
  72. family=reuse;policy=default;refs=2;consumer=1%
  73. family=bag;stmt=cte_joined_outer_mat;threshold=20000
  74. family=reuse;policy=default;refs=16;consumer=100%
  75. family=bag;stmt=cte_filtered;threshold=200
  76. family=bag;stmt=cte_joined_outer_mat;threshold=100000
  77. family=bag;stmt=join_on;threshold=20000
  78. family=reuse;policy=mat;refs=8;consumer=1%
  79. family=bag;stmt=join_where;threshold=200
  80. family=reuse;policy=notmat;refs=16;consumer=100%
  81. family=bag;stmt=comma_join;threshold=20
  82. family=bag;stmt=cte_filtered_notmat;threshold=20
  83. family=reuse;policy=default;refs=4;consumer=1%
  84. family=bag;stmt=join_where;threshold=180000
  85. family=bag;stmt=join_on;threshold=2000
  86. family=bag;stmt=filtered_derived;threshold=200
  87. family=bag;stmt=join_where;threshold=20000
  88. family=bag;stmt=filtered_derived;threshold=100000
  89. family=reuse;policy=notmat;refs=4;consumer=1%
  90. family=reuse;policy=mat;refs=16;consumer=100%
  91. family=bag;stmt=cte_joined_outer_notmat;threshold=20000;position=last

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
| family=bag;stmt=cte_joined_outer_notmat;threshold=20000;position=last | j05n_cte_joined_outer_notmat | valid | `dbf5856d383eb7fb` | policy=not_materialized; observed_cte_scan=false |
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
