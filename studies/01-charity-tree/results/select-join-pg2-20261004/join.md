# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 91 cells, 91 valid, engine `postgres`, topology `single`, scale `cte`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-pg2-20261004`
- run tag: `run/study-01/select-join-pg2-20261004`
- environment: `host-zenbook-ux5406sa`
- repo commit: `7e148174b5211225f6b4f243fc02e0de417a0586`
- engine version: `PostgreSQL 17.11 (Debian 17.11-1.pgdg13+2) on x86_64-pc-linux-gnu, compiled by gcc (Debian 14.2.0-19) 14.2.0, 64-bit`
- engine image: `localhost/charitybench-select:select-join-pg2-20261004` @ `7296f210ae81031ec955dbad9a67a84fe958572a2153b8d0826a647522904dc1`
- inputs digest: `112577024534a9e91932f78d3ba8b9e530e0caea0aeefdf851737e578be7bab6`
- dataset fingerprint: `2e06ecb3915eb290`
- client CPU quota: `200000 100000`, memory: `2147483648`
- setting `client_workers`: 1 sequential reader; 0 writers
- setting `connection_pool`: 2 = 1 reader/setup + 1 reserve
- setting `customers`: 2000
- setting `distributed_rpc`: N/A on PostgreSQL single-node
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
| j01_join_where | valid | true |  |
| j02_join_on | valid | true |  |
| j03_filtered_derived | valid | true |  |
| j04_cte_filtered | valid | true | policy=default; observed_cte_scan=false |
| j04m_cte_filtered_mat | valid | true | policy=materialized; observed_cte_scan=true |
| j04n_cte_filtered_notmat | valid | true | policy=not_materialized; observed_cte_scan=false |
| j05_cte_joined_outer | valid | true | policy=default; observed_cte_scan=false |
| j05m_cte_joined_outer_mat | valid | false | policy=materialized; observed_cte_scan=true |
| j05n_cte_joined_outer_notmat | valid | true | policy=not_materialized; observed_cte_scan=false |
| j06_comma_join | valid | true |  |
| j11_optional_on | valid | true |  |
| r_default_n4_c100 | valid | false | policy=default; observed_cte_scan=true |
| r_mat_n4_c100 | valid | false | policy=materialized; observed_cte_scan=true |
| r_notmat_n4_c100 | valid | false | policy=not_materialized; observed_cte_scan=false |

## P1 — six inner-join spellings and their CTE policies (median ms)

| threshold | join_where | join_on | filtered_derived | cte_filtered | cte_filtered_mat | cte_filtered_notmat | cte_joined_outer | cte_joined_outer_mat | cte_joined_outer_notmat | comma_join |
|---|---|---|---|---|---|---|---|---|---|---|
| 20 (0.01%) | 0.3890 | 0.3628 | 0.5796 | 0.4103 | 0.6016 | 0.5527 | 0.3993 | 37.0778 | 0.3957 | 0.5505 |
| 200 (0.1%) | 0.8207 | 0.8235 | 0.6290 | 0.5609 | 0.5746 | 0.8413 | 0.8494 | 40.4844 | 0.7032 | 0.6202 |
| 2000 (1%) | 1.3901 | 1.6320 | 1.5771 | 1.4668 | 1.5003 | 1.3624 | 1.2700 | 38.9063 | 1.6029 | 1.6504 |
| 20000 (10%) | 6.0667 | 5.8614 | 5.9304 | 5.5818 | 6.1904 | 5.8829 | 5.8073 | 42.7557 | 6.1446 | 5.5945 |
| 100000 (50%) | 37.6218 | 37.8893 | 39.8349 | 39.8509 | 44.9175 | 37.4683 | 36.4659 | 53.1524 | 36.7364 | 37.0255 |
| 180000 (90%) | 47.6480 | 40.9529 | 42.5508 | 49.7117 | 48.4859 | 44.6485 | 47.2707 | 64.7961 | 46.9925 | 46.1841 |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | join_on vs join_where | 0.3890 | 0.3628 | 0.9326 | true |
| 20 | filtered_derived vs join_where | 0.3890 | 0.5796 | 1.4900 | true |
| 20 | cte_filtered vs join_where | 0.3890 | 0.4103 | 1.0547 | true |
| 20 | cte_filtered_mat vs join_where | 0.3890 | 0.6016 | 1.5465 | false |
| 20 | cte_filtered_notmat vs join_where | 0.3890 | 0.5527 | 1.4209 | true |
| 20 | cte_joined_outer vs join_where | 0.3890 | 0.3993 | 1.0265 | true |
| 20 | cte_joined_outer_mat vs join_where | 0.3890 | 37.0778 | 95.3136 | false |
| 20 | cte_joined_outer_notmat vs join_where | 0.3890 | 0.3957 | 1.0173 | true |
| 20 | comma_join vs join_where | 0.3890 | 0.5505 | 1.4151 | true |
| 200 | join_on vs join_where | 0.8207 | 0.8235 | 1.0034 | true |
| 200 | filtered_derived vs join_where | 0.8207 | 0.6290 | 0.7664 | true |
| 200 | cte_filtered vs join_where | 0.8207 | 0.5609 | 0.6834 | true |
| 200 | cte_filtered_mat vs join_where | 0.8207 | 0.5746 | 0.7001 | false |
| 200 | cte_filtered_notmat vs join_where | 0.8207 | 0.8413 | 1.0251 | true |
| 200 | cte_joined_outer vs join_where | 0.8207 | 0.8494 | 1.0349 | true |
| 200 | cte_joined_outer_mat vs join_where | 0.8207 | 40.4844 | 49.3263 | false |
| 200 | cte_joined_outer_notmat vs join_where | 0.8207 | 0.7032 | 0.8567 | true |
| 200 | comma_join vs join_where | 0.8207 | 0.6202 | 0.7556 | true |
| 2000 | join_on vs join_where | 1.3901 | 1.6320 | 1.1740 | true |
| 2000 | filtered_derived vs join_where | 1.3901 | 1.5771 | 1.1345 | true |
| 2000 | cte_filtered vs join_where | 1.3901 | 1.4668 | 1.0552 | true |
| 2000 | cte_filtered_mat vs join_where | 1.3901 | 1.5003 | 1.0793 | false |
| 2000 | cte_filtered_notmat vs join_where | 1.3901 | 1.3624 | 0.9801 | true |
| 2000 | cte_joined_outer vs join_where | 1.3901 | 1.2700 | 0.9137 | true |
| 2000 | cte_joined_outer_mat vs join_where | 1.3901 | 38.9063 | 27.9888 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 1.3901 | 1.6029 | 1.1531 | true |
| 2000 | comma_join vs join_where | 1.3901 | 1.6504 | 1.1873 | true |
| 20000 | join_on vs join_where | 6.0667 | 5.8614 | 0.9662 | true |
| 20000 | filtered_derived vs join_where | 6.0667 | 5.9304 | 0.9775 | true |
| 20000 | cte_filtered vs join_where | 6.0667 | 5.5818 | 0.9201 | true |
| 20000 | cte_filtered_mat vs join_where | 6.0667 | 6.1904 | 1.0204 | false |
| 20000 | cte_filtered_notmat vs join_where | 6.0667 | 5.8829 | 0.9697 | true |
| 20000 | cte_joined_outer vs join_where | 6.0667 | 5.8073 | 0.9573 | true |
| 20000 | cte_joined_outer_mat vs join_where | 6.0667 | 42.7557 | 7.0476 | false |
| 20000 | cte_joined_outer_notmat vs join_where | 6.0667 | 6.1446 | 1.0129 | true |
| 20000 | comma_join vs join_where | 6.0667 | 5.5945 | 0.9222 | true |
| 100000 | join_on vs join_where | 37.6218 | 37.8893 | 1.0071 | true |
| 100000 | filtered_derived vs join_where | 37.6218 | 39.8349 | 1.0588 | true |
| 100000 | cte_filtered vs join_where | 37.6218 | 39.8509 | 1.0592 | true |
| 100000 | cte_filtered_mat vs join_where | 37.6218 | 44.9175 | 1.1939 | false |
| 100000 | cte_filtered_notmat vs join_where | 37.6218 | 37.4683 | 0.9959 | true |
| 100000 | cte_joined_outer vs join_where | 37.6218 | 36.4659 | 0.9693 | true |
| 100000 | cte_joined_outer_mat vs join_where | 37.6218 | 53.1524 | 1.4128 | false |
| 100000 | cte_joined_outer_notmat vs join_where | 37.6218 | 36.7364 | 0.9765 | true |
| 100000 | comma_join vs join_where | 37.6218 | 37.0255 | 0.9841 | true |
| 180000 | join_on vs join_where | 47.6480 | 40.9529 | 0.8595 | true |
| 180000 | filtered_derived vs join_where | 47.6480 | 42.5508 | 0.8930 | true |
| 180000 | cte_filtered vs join_where | 47.6480 | 49.7117 | 1.0433 | true |
| 180000 | cte_filtered_mat vs join_where | 47.6480 | 48.4859 | 1.0176 | false |
| 180000 | cte_filtered_notmat vs join_where | 47.6480 | 44.6485 | 0.9370 | true |
| 180000 | cte_joined_outer vs join_where | 47.6480 | 47.2707 | 0.9921 | true |
| 180000 | cte_joined_outer_mat vs join_where | 47.6480 | 64.7961 | 1.3599 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 47.6480 | 46.9925 | 0.9862 | true |
| 180000 | comma_join vs join_where | 47.6480 | 46.1841 | 0.9693 | true |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_filtered_mat vs cte_filtered | 0.4103 | 0.6016 | 1.4662 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_filtered_notmat vs cte_filtered | 0.4103 | 0.5527 | 1.3472 | true | policy=not_materialized; observed_cte_scan=false |
| 20 | cte_joined_outer_mat vs cte_joined_outer | 0.3993 | 37.0778 | 92.8513 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 0.3993 | 0.3957 | 0.9910 | true | policy=not_materialized; observed_cte_scan=false |
| 200 | cte_filtered_mat vs cte_filtered | 0.5609 | 0.5746 | 1.0245 | false | policy=materialized; observed_cte_scan=true |
| 200 | cte_filtered_notmat vs cte_filtered | 0.5609 | 0.8413 | 1.5000 | true | policy=not_materialized; observed_cte_scan=false |
| 200 | cte_joined_outer_mat vs cte_joined_outer | 0.8494 | 40.4844 | 47.6644 | false | policy=materialized; observed_cte_scan=true |
| 200 | cte_joined_outer_notmat vs cte_joined_outer | 0.8494 | 0.7032 | 0.8279 | true | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_filtered_mat vs cte_filtered | 1.4668 | 1.5003 | 1.0228 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_filtered_notmat vs cte_filtered | 1.4668 | 1.3624 | 0.9288 | true | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 1.2700 | 38.9063 | 30.6338 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 1.2700 | 1.6029 | 1.2621 | true | policy=not_materialized; observed_cte_scan=false |
| 20000 | cte_filtered_mat vs cte_filtered | 5.5818 | 6.1904 | 1.1090 | false | policy=materialized; observed_cte_scan=true |
| 20000 | cte_filtered_notmat vs cte_filtered | 5.5818 | 5.8829 | 1.0539 | true | policy=not_materialized; observed_cte_scan=false |
| 20000 | cte_joined_outer_mat vs cte_joined_outer | 5.8073 | 42.7557 | 7.3624 | false | policy=materialized; observed_cte_scan=true |
| 20000 | cte_joined_outer_notmat vs cte_joined_outer | 5.8073 | 6.1446 | 1.0581 | true | policy=not_materialized; observed_cte_scan=false |
| 100000 | cte_filtered_mat vs cte_filtered | 39.8509 | 44.9175 | 1.1271 | false | policy=materialized; observed_cte_scan=true |
| 100000 | cte_filtered_notmat vs cte_filtered | 39.8509 | 37.4683 | 0.9402 | true | policy=not_materialized; observed_cte_scan=false |
| 100000 | cte_joined_outer_mat vs cte_joined_outer | 36.4659 | 53.1524 | 1.4576 | false | policy=materialized; observed_cte_scan=true |
| 100000 | cte_joined_outer_notmat vs cte_joined_outer | 36.4659 | 36.7364 | 1.0074 | true | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_filtered_mat vs cte_filtered | 49.7117 | 48.4859 | 0.9753 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_filtered_notmat vs cte_filtered | 49.7117 | 44.6485 | 0.8981 | true | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 47.2707 | 64.7961 | 1.3707 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 47.2707 | 46.9925 | 0.9941 | true | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 0.6543 | 1 | `2beb12c2700fe6cd` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 1.9196 | 1 | `ae8e946d87d174f1` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | 2.1791 | 2 | `aef241ae98ffed1c` | policy=default; observed_cte_scan=true | valid |
| default | 2 | 100% | 2.3629 | 2 | `aef241ae98ffed1c` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 1% | 2.3796 | 4 | `4052ca75f43959fb` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 2.4233 | 4 | `4052ca75f43959fb` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | 2.7135 | 8 | `836170c2201c675d` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 100% | 3.1735 | 8 | `836170c2201c675d` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 1% | 3.1462 | 16 | `9e8e5819d256402e` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 3.1885 | 16 | `9e8e5819d256402e` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 2.2570 | 1 | `3b31b8c96a2778e9` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 2.4236 | 1 | `3b31b8c96a2778e9` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | 2.6064 | 2 | `aef241ae98ffed1c` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 100% | 2.2794 | 2 | `aef241ae98ffed1c` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 1% | 2.2673 | 4 | `4052ca75f43959fb` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 2.7523 | 4 | `4052ca75f43959fb` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | 2.5571 | 8 | `836170c2201c675d` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 100% | 2.9563 | 8 | `836170c2201c675d` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 1% | 3.2338 | 16 | `9e8e5819d256402e` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 3.6705 | 16 | `9e8e5819d256402e` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 0.4827 | 1 | `2beb12c2700fe6cd` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 2.3477 | 1 | `ae8e946d87d174f1` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | 0.5664 | 2 | `b5a88f6403cacc71` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 100% | 4.8402 | 2 | `3dec56afda6778f7` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 1% | 0.8563 | 4 | `2b6ccce2a72df800` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 7.0597 | 4 | `9bc6945fc91b283b` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | 1.3824 | 8 | `caf493c80d09c2fe` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 100% | 9.4583 | 8 | `b09c2c0768541f93` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 1% | 2.5444 | 16 | `7cb8b15715a86c9f` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 14.6995 | 16 | `c2a32a53c5d580bc` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 3.4495 | false |
| notmat | 1 | 1% | 0.7377 | true |
| mat | 1 | 100% | 1.2625 | false |
| notmat | 1 | 100% | 1.2230 | true |
| mat | 2 | 1% | 1.1961 | true |
| notmat | 2 | 1% | 0.2599 | false |
| mat | 2 | 100% | 0.9647 | true |
| notmat | 2 | 100% | 2.0484 | false |
| mat | 4 | 1% | 0.9528 | true |
| notmat | 4 | 1% | 0.3599 | false |
| mat | 4 | 100% | 1.1357 | true |
| notmat | 4 | 100% | 2.9132 | false |
| mat | 8 | 1% | 0.9424 | true |
| notmat | 8 | 1% | 0.5094 | false |
| mat | 8 | 100% | 0.9316 | true |
| notmat | 8 | 100% | 2.9804 | false |
| mat | 16 | 1% | 1.0279 | true |
| notmat | 16 | 1% | 0.8087 | false |
| mat | 16 | 100% | 1.1512 | true |
| notmat | 16 | 100% | 4.6102 | false |

## Repeated-design position control

| first median ms | last median ms | ratio (last/first) | same plan |
|---|---|---|---|
| 40.484395 | 36.286378 | 0.8963 | true |

## Execution order

- order seed `42643254`, 5 interleaved rounds
  1. family=bag;stmt=cte_joined_outer_mat;threshold=200
  2. family=reuse;policy=mat;refs=2;consumer=1%
  3. family=bag;stmt=join_where;threshold=2000
  4. family=bag;stmt=cte_filtered;threshold=180000
  5. family=bag;stmt=filtered_derived;threshold=20
  6. family=reuse;policy=notmat;refs=1;consumer=100%
  7. family=bag;stmt=cte_filtered_notmat;threshold=100000
  8. family=bag;stmt=join_on;threshold=200
  9. family=bag;stmt=cte_filtered;threshold=2000
  10. family=reuse;policy=mat;refs=8;consumer=1%
  11. family=reuse;policy=mat;refs=4;consumer=1%
  12. family=bag;stmt=cte_filtered_mat;threshold=200
  13. family=bag;stmt=cte_filtered_mat;threshold=180000
  14. family=bag;stmt=join_where;threshold=200
  15. family=reuse;policy=mat;refs=8;consumer=100%
  16. family=bag;stmt=cte_joined_outer_mat;threshold=20
  17. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  18. family=bag;stmt=cte_joined_outer;threshold=200
  19. family=bag;stmt=cte_joined_outer;threshold=100000
  20. family=bag;stmt=cte_filtered_notmat;threshold=180000
  21. family=reuse;policy=mat;refs=16;consumer=1%
  22. family=bag;stmt=join_where;threshold=20
  23. family=reuse;policy=mat;refs=2;consumer=100%
  24. family=reuse;policy=default;refs=16;consumer=100%
  25. family=bag;stmt=cte_filtered_notmat;threshold=20000
  26. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  27. family=reuse;policy=notmat;refs=8;consumer=1%
  28. family=bag;stmt=cte_filtered;threshold=20
  29. family=bag;stmt=cte_joined_outer;threshold=180000
  30. family=bag;stmt=cte_filtered_mat;threshold=20
  31. family=bag;stmt=filtered_derived;threshold=2000
  32. family=bag;stmt=comma_join;threshold=20000
  33. family=bag;stmt=filtered_derived;threshold=20000
  34. family=bag;stmt=cte_filtered;threshold=100000
  35. family=bag;stmt=cte_joined_outer_mat;threshold=20000
  36. family=bag;stmt=cte_filtered_notmat;threshold=20
  37. family=bag;stmt=cte_filtered_mat;threshold=100000
  38. family=bag;stmt=cte_filtered_notmat;threshold=200
  39. family=bag;stmt=filtered_derived;threshold=200
  40. family=reuse;policy=default;refs=8;consumer=1%
  41. family=reuse;policy=default;refs=1;consumer=100%
  42. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  43. family=bag;stmt=cte_joined_outer_mat;threshold=100000
  44. family=bag;stmt=join_on;threshold=20000
  45. family=bag;stmt=comma_join;threshold=20
  46. family=reuse;policy=default;refs=16;consumer=1%
  47. family=bag;stmt=comma_join;threshold=200
  48. family=reuse;policy=default;refs=2;consumer=100%
  49. family=reuse;policy=notmat;refs=2;consumer=1%
  50. family=reuse;policy=default;refs=2;consumer=1%
  51. family=bag;stmt=cte_filtered;threshold=200
  52. family=reuse;policy=notmat;refs=16;consumer=100%
  53. family=bag;stmt=join_where;threshold=180000
  54. family=reuse;policy=default;refs=8;consumer=100%
  55. family=bag;stmt=cte_joined_outer;threshold=20
  56. family=bag;stmt=cte_filtered_mat;threshold=2000
  57. family=reuse;policy=notmat;refs=4;consumer=1%
  58. family=bag;stmt=cte_filtered;threshold=20000
  59. family=bag;stmt=comma_join;threshold=180000
  60. family=bag;stmt=filtered_derived;threshold=180000
  61. family=reuse;policy=notmat;refs=2;consumer=100%
  62. family=bag;stmt=join_on;threshold=100000
  63. family=reuse;policy=default;refs=1;consumer=1%
  64. family=bag;stmt=cte_joined_outer_notmat;threshold=200
  65. family=bag;stmt=cte_joined_outer_notmat;threshold=20000
  66. family=bag;stmt=join_where;threshold=100000
  67. family=bag;stmt=cte_filtered_mat;threshold=20000
  68. family=reuse;policy=notmat;refs=8;consumer=100%
  69. family=reuse;policy=mat;refs=1;consumer=100%
  70. family=bag;stmt=cte_filtered_notmat;threshold=2000
  71. family=reuse;policy=notmat;refs=4;consumer=100%
  72. family=reuse;policy=notmat;refs=16;consumer=1%
  73. family=bag;stmt=cte_joined_outer;threshold=20000
  74. family=reuse;policy=mat;refs=16;consumer=100%
  75. family=reuse;policy=default;refs=4;consumer=100%
  76. family=bag;stmt=filtered_derived;threshold=100000
  77. family=bag;stmt=join_where;threshold=20000
  78. family=bag;stmt=cte_joined_outer_notmat;threshold=100000
  79. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  80. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  81. family=bag;stmt=join_on;threshold=2000
  82. family=reuse;policy=notmat;refs=1;consumer=1%
  83. family=bag;stmt=join_on;threshold=180000
  84. family=reuse;policy=mat;refs=4;consumer=100%
  85. family=reuse;policy=default;refs=4;consumer=1%
  86. family=bag;stmt=comma_join;threshold=100000
  87. family=bag;stmt=comma_join;threshold=2000
  88. family=reuse;policy=mat;refs=1;consumer=1%
  89. family=bag;stmt=cte_joined_outer;threshold=2000
  90. family=bag;stmt=join_on;threshold=20
  91. family=bag;stmt=cte_joined_outer_mat;threshold=200;position=last

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
| family=bag;stmt=comma_join;threshold=100000 | j06_comma_join | valid | `f59fe6b12975ad43` |  |
| family=bag;stmt=comma_join;threshold=180000 | j06_comma_join | valid | `052dca246c43e506` |  |
| family=bag;stmt=comma_join;threshold=20 | j06_comma_join | valid | `b6e06754f69a798c` |  |
| family=bag;stmt=comma_join;threshold=200 | j06_comma_join | valid | `e9b5c491684a318e` |  |
| family=bag;stmt=comma_join;threshold=2000 | j06_comma_join | valid | `72883cfe65e9dc9a` |  |
| family=bag;stmt=comma_join;threshold=20000 | j06_comma_join | valid | `3f780ee12ad1d16b` |  |
| family=bag;stmt=cte_filtered;threshold=100000 | j04_cte_filtered | valid | `f59fe6b12975ad43` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=180000 | j04_cte_filtered | valid | `052dca246c43e506` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=20 | j04_cte_filtered | valid | `b6e06754f69a798c` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=200 | j04_cte_filtered | valid | `e9b5c491684a318e` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=2000 | j04_cte_filtered | valid | `72883cfe65e9dc9a` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=20000 | j04_cte_filtered | valid | `3f780ee12ad1d16b` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_mat;threshold=100000 | j04m_cte_filtered_mat | valid | `2ed36a8082a03f8d` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=180000 | j04m_cte_filtered_mat | valid | `31669afcb0a16b04` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=20 | j04m_cte_filtered_mat | valid | `fef2a456dbb4dbce` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=200 | j04m_cte_filtered_mat | valid | `cdfe2f017b46298c` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=2000 | j04m_cte_filtered_mat | valid | `9f56a67dd55d1250` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=20000 | j04m_cte_filtered_mat | valid | `62bd67b707a60ca7` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_notmat;threshold=100000 | j04n_cte_filtered_notmat | valid | `f59fe6b12975ad43` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=180000 | j04n_cte_filtered_notmat | valid | `052dca246c43e506` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=20 | j04n_cte_filtered_notmat | valid | `b6e06754f69a798c` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=200 | j04n_cte_filtered_notmat | valid | `e9b5c491684a318e` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=2000 | j04n_cte_filtered_notmat | valid | `72883cfe65e9dc9a` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=20000 | j04n_cte_filtered_notmat | valid | `3f780ee12ad1d16b` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=100000 | j05_cte_joined_outer | valid | `f59fe6b12975ad43` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=180000 | j05_cte_joined_outer | valid | `052dca246c43e506` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20 | j05_cte_joined_outer | valid | `b6e06754f69a798c` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=200 | j05_cte_joined_outer | valid | `e9b5c491684a318e` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=2000 | j05_cte_joined_outer | valid | `72883cfe65e9dc9a` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20000 | j05_cte_joined_outer | valid | `3f780ee12ad1d16b` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_mat;threshold=100000 | j05m_cte_joined_outer_mat | valid | `9260d16a40769f90` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | j05m_cte_joined_outer_mat | valid | `9260d16a40769f90` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | j05m_cte_joined_outer_mat | valid | `9260d16a40769f90` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=200 | j05m_cte_joined_outer_mat | valid | `9260d16a40769f90` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | j05m_cte_joined_outer_mat | valid | `9260d16a40769f90` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20000 | j05m_cte_joined_outer_mat | valid | `9260d16a40769f90` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=200;position=last | j05m_cte_joined_outer_mat | valid | `9260d16a40769f90` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_notmat;threshold=100000 | j05n_cte_joined_outer_notmat | valid | `f59fe6b12975ad43` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | j05n_cte_joined_outer_notmat | valid | `052dca246c43e506` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | j05n_cte_joined_outer_notmat | valid | `b6e06754f69a798c` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=200 | j05n_cte_joined_outer_notmat | valid | `e9b5c491684a318e` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | j05n_cte_joined_outer_notmat | valid | `72883cfe65e9dc9a` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20000 | j05n_cte_joined_outer_notmat | valid | `3f780ee12ad1d16b` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=filtered_derived;threshold=100000 | j03_filtered_derived | valid | `f59fe6b12975ad43` |  |
| family=bag;stmt=filtered_derived;threshold=180000 | j03_filtered_derived | valid | `052dca246c43e506` |  |
| family=bag;stmt=filtered_derived;threshold=20 | j03_filtered_derived | valid | `b6e06754f69a798c` |  |
| family=bag;stmt=filtered_derived;threshold=200 | j03_filtered_derived | valid | `e9b5c491684a318e` |  |
| family=bag;stmt=filtered_derived;threshold=2000 | j03_filtered_derived | valid | `72883cfe65e9dc9a` |  |
| family=bag;stmt=filtered_derived;threshold=20000 | j03_filtered_derived | valid | `3f780ee12ad1d16b` |  |
| family=bag;stmt=join_on;threshold=100000 | j02_join_on | valid | `f59fe6b12975ad43` |  |
| family=bag;stmt=join_on;threshold=180000 | j02_join_on | valid | `052dca246c43e506` |  |
| family=bag;stmt=join_on;threshold=20 | j02_join_on | valid | `b6e06754f69a798c` |  |
| family=bag;stmt=join_on;threshold=200 | j02_join_on | valid | `e9b5c491684a318e` |  |
| family=bag;stmt=join_on;threshold=2000 | j02_join_on | valid | `72883cfe65e9dc9a` |  |
| family=bag;stmt=join_on;threshold=20000 | j02_join_on | valid | `3f780ee12ad1d16b` |  |
| family=bag;stmt=join_where;threshold=100000 | j01_join_where | valid | `f59fe6b12975ad43` |  |
| family=bag;stmt=join_where;threshold=180000 | j01_join_where | valid | `052dca246c43e506` |  |
| family=bag;stmt=join_where;threshold=20 | j01_join_where | valid | `b6e06754f69a798c` |  |
| family=bag;stmt=join_where;threshold=200 | j01_join_where | valid | `e9b5c491684a318e` |  |
| family=bag;stmt=join_where;threshold=2000 | j01_join_where | valid | `72883cfe65e9dc9a` |  |
| family=bag;stmt=join_where;threshold=20000 | j01_join_where | valid | `3f780ee12ad1d16b` |  |
| family=reuse;policy=default;refs=16;consumer=1% | r_default_n16_c1 | valid | `9e8e5819d256402e` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=16;consumer=100% | r_default_n16_c100 | valid | `9e8e5819d256402e` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=1;consumer=1% | r_default_n1_c1 | valid | `2beb12c2700fe6cd` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=1;consumer=100% | r_default_n1_c100 | valid | `ae8e946d87d174f1` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=2;consumer=1% | r_default_n2_c1 | valid | `aef241ae98ffed1c` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=2;consumer=100% | r_default_n2_c100 | valid | `aef241ae98ffed1c` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=1% | r_default_n4_c1 | valid | `4052ca75f43959fb` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=100% | r_default_n4_c100 | valid | `4052ca75f43959fb` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=8;consumer=1% | r_default_n8_c1 | valid | `836170c2201c675d` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=8;consumer=100% | r_default_n8_c100 | valid | `836170c2201c675d` | policy=default; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=1% | r_mat_n16_c1 | valid | `9e8e5819d256402e` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=100% | r_mat_n16_c100 | valid | `9e8e5819d256402e` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=1% | r_mat_n1_c1 | valid | `3b31b8c96a2778e9` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=100% | r_mat_n1_c100 | valid | `3b31b8c96a2778e9` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=2;consumer=1% | r_mat_n2_c1 | valid | `aef241ae98ffed1c` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=2;consumer=100% | r_mat_n2_c100 | valid | `aef241ae98ffed1c` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=1% | r_mat_n4_c1 | valid | `4052ca75f43959fb` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=100% | r_mat_n4_c100 | valid | `4052ca75f43959fb` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=8;consumer=1% | r_mat_n8_c1 | valid | `836170c2201c675d` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=8;consumer=100% | r_mat_n8_c100 | valid | `836170c2201c675d` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=notmat;refs=16;consumer=1% | r_notmat_n16_c1 | valid | `7cb8b15715a86c9f` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `c2a32a53c5d580bc` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `2beb12c2700fe6cd` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `ae8e946d87d174f1` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=2;consumer=1% | r_notmat_n2_c1 | valid | `b5a88f6403cacc71` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=2;consumer=100% | r_notmat_n2_c100 | valid | `3dec56afda6778f7` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `2b6ccce2a72df800` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `9bc6945fc91b283b` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=8;consumer=1% | r_notmat_n8_c1 | valid | `caf493c80d09c2fe` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=8;consumer=100% | r_notmat_n8_c100 | valid | `b09c2c0768541f93` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- one host, single client: no independent hardware, real-network or equal-total-budget comparison
- no concurrent writers: the isolated read contract only; a matched blended read/write phase is a recorded gap
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
