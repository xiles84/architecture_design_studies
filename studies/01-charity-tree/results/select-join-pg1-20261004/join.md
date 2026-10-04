# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 91 cells, 91 valid, engine `postgres`, topology `single`, scale `cte`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-pg1-20261004`
- run tag: `run/study-01/select-join-pg1-20261004`
- environment: `host-zenbook-ux5406sa`
- repo commit: `7e148174b5211225f6b4f243fc02e0de417a0586`
- engine version: `PostgreSQL 17.11 (Debian 17.11-1.pgdg13+2) on x86_64-pc-linux-gnu, compiled by gcc (Debian 14.2.0-19) 14.2.0, 64-bit`
- engine image: `localhost/charitybench-select:select-join-pg1-20261004` @ `7296f210ae81031ec955dbad9a67a84fe958572a2153b8d0826a647522904dc1`
- inputs digest: `070df793e1354a33d7e9363877d83f9b77fc14647b9818a0631096c314f2cfac`
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
| 20 (0.01%) | 0.4177 | 0.3605 | 0.6337 | 0.4700 | 0.3920 | 0.4427 | 0.5048 | 36.9053 | 0.3844 | 0.4979 |
| 200 (0.1%) | 0.6714 | 0.7784 | 0.7846 | 0.6489 | 0.5953 | 0.8553 | 0.5446 | 40.2530 | 0.7009 | 0.5782 |
| 2000 (1%) | 1.4119 | 1.3183 | 1.7095 | 1.7570 | 1.5628 | 1.6232 | 1.4107 | 37.5496 | 1.5546 | 1.3679 |
| 20000 (10%) | 6.1398 | 5.5920 | 5.7220 | 5.9134 | 6.5758 | 5.8440 | 6.1904 | 42.0132 | 5.6211 | 6.2381 |
| 100000 (50%) | 38.7161 | 34.9793 | 37.1299 | 37.8795 | 41.9282 | 38.4509 | 37.6850 | 53.5020 | 36.8440 | 37.6825 |
| 180000 (90%) | 46.7475 | 45.3185 | 44.1086 | 48.0588 | 48.8781 | 45.6820 | 45.0621 | 63.7415 | 42.5202 | 45.5845 |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | join_on vs join_where | 0.4177 | 0.3605 | 0.8631 | true |
| 20 | filtered_derived vs join_where | 0.4177 | 0.6337 | 1.5172 | true |
| 20 | cte_filtered vs join_where | 0.4177 | 0.4700 | 1.1253 | true |
| 20 | cte_filtered_mat vs join_where | 0.4177 | 0.3920 | 0.9384 | false |
| 20 | cte_filtered_notmat vs join_where | 0.4177 | 0.4427 | 1.0600 | true |
| 20 | cte_joined_outer vs join_where | 0.4177 | 0.5048 | 1.2087 | true |
| 20 | cte_joined_outer_mat vs join_where | 0.4177 | 36.9053 | 88.3604 | false |
| 20 | cte_joined_outer_notmat vs join_where | 0.4177 | 0.3844 | 0.9203 | true |
| 20 | comma_join vs join_where | 0.4177 | 0.4979 | 1.1921 | true |
| 200 | join_on vs join_where | 0.6714 | 0.7784 | 1.1593 | true |
| 200 | filtered_derived vs join_where | 0.6714 | 0.7846 | 1.1685 | true |
| 200 | cte_filtered vs join_where | 0.6714 | 0.6489 | 0.9664 | true |
| 200 | cte_filtered_mat vs join_where | 0.6714 | 0.5953 | 0.8866 | false |
| 200 | cte_filtered_notmat vs join_where | 0.6714 | 0.8553 | 1.2739 | true |
| 200 | cte_joined_outer vs join_where | 0.6714 | 0.5446 | 0.8111 | true |
| 200 | cte_joined_outer_mat vs join_where | 0.6714 | 40.2530 | 59.9526 | false |
| 200 | cte_joined_outer_notmat vs join_where | 0.6714 | 0.7009 | 1.0440 | true |
| 200 | comma_join vs join_where | 0.6714 | 0.5782 | 0.8612 | true |
| 2000 | join_on vs join_where | 1.4119 | 1.3183 | 0.9337 | true |
| 2000 | filtered_derived vs join_where | 1.4119 | 1.7095 | 1.2107 | true |
| 2000 | cte_filtered vs join_where | 1.4119 | 1.7570 | 1.2444 | true |
| 2000 | cte_filtered_mat vs join_where | 1.4119 | 1.5628 | 1.1068 | false |
| 2000 | cte_filtered_notmat vs join_where | 1.4119 | 1.6232 | 1.1497 | true |
| 2000 | cte_joined_outer vs join_where | 1.4119 | 1.4107 | 0.9991 | true |
| 2000 | cte_joined_outer_mat vs join_where | 1.4119 | 37.5496 | 26.5944 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 1.4119 | 1.5546 | 1.1010 | true |
| 2000 | comma_join vs join_where | 1.4119 | 1.3679 | 0.9688 | true |
| 20000 | join_on vs join_where | 6.1398 | 5.5920 | 0.9108 | true |
| 20000 | filtered_derived vs join_where | 6.1398 | 5.7220 | 0.9320 | true |
| 20000 | cte_filtered vs join_where | 6.1398 | 5.9134 | 0.9631 | true |
| 20000 | cte_filtered_mat vs join_where | 6.1398 | 6.5758 | 1.0710 | false |
| 20000 | cte_filtered_notmat vs join_where | 6.1398 | 5.8440 | 0.9518 | true |
| 20000 | cte_joined_outer vs join_where | 6.1398 | 6.1904 | 1.0082 | true |
| 20000 | cte_joined_outer_mat vs join_where | 6.1398 | 42.0132 | 6.8428 | false |
| 20000 | cte_joined_outer_notmat vs join_where | 6.1398 | 5.6211 | 0.9155 | true |
| 20000 | comma_join vs join_where | 6.1398 | 6.2381 | 1.0160 | true |
| 100000 | join_on vs join_where | 38.7161 | 34.9793 | 0.9035 | true |
| 100000 | filtered_derived vs join_where | 38.7161 | 37.1299 | 0.9590 | true |
| 100000 | cte_filtered vs join_where | 38.7161 | 37.8795 | 0.9784 | true |
| 100000 | cte_filtered_mat vs join_where | 38.7161 | 41.9282 | 1.0830 | false |
| 100000 | cte_filtered_notmat vs join_where | 38.7161 | 38.4509 | 0.9932 | true |
| 100000 | cte_joined_outer vs join_where | 38.7161 | 37.6850 | 0.9734 | true |
| 100000 | cte_joined_outer_mat vs join_where | 38.7161 | 53.5020 | 1.3819 | false |
| 100000 | cte_joined_outer_notmat vs join_where | 38.7161 | 36.8440 | 0.9516 | true |
| 100000 | comma_join vs join_where | 38.7161 | 37.6825 | 0.9733 | true |
| 180000 | join_on vs join_where | 46.7475 | 45.3185 | 0.9694 | true |
| 180000 | filtered_derived vs join_where | 46.7475 | 44.1086 | 0.9435 | true |
| 180000 | cte_filtered vs join_where | 46.7475 | 48.0588 | 1.0281 | true |
| 180000 | cte_filtered_mat vs join_where | 46.7475 | 48.8781 | 1.0456 | false |
| 180000 | cte_filtered_notmat vs join_where | 46.7475 | 45.6820 | 0.9772 | true |
| 180000 | cte_joined_outer vs join_where | 46.7475 | 45.0621 | 0.9639 | true |
| 180000 | cte_joined_outer_mat vs join_where | 46.7475 | 63.7415 | 1.3635 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 46.7475 | 42.5202 | 0.9096 | true |
| 180000 | comma_join vs join_where | 46.7475 | 45.5845 | 0.9751 | true |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_filtered_mat vs cte_filtered | 0.4700 | 0.3920 | 0.8340 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_filtered_notmat vs cte_filtered | 0.4700 | 0.4427 | 0.9420 | true | policy=not_materialized; observed_cte_scan=false |
| 20 | cte_joined_outer_mat vs cte_joined_outer | 0.5048 | 36.9053 | 73.1056 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 0.5048 | 0.3844 | 0.7615 | true | policy=not_materialized; observed_cte_scan=false |
| 200 | cte_filtered_mat vs cte_filtered | 0.6489 | 0.5953 | 0.9174 | false | policy=materialized; observed_cte_scan=true |
| 200 | cte_filtered_notmat vs cte_filtered | 0.6489 | 0.8553 | 1.3182 | true | policy=not_materialized; observed_cte_scan=false |
| 200 | cte_joined_outer_mat vs cte_joined_outer | 0.5446 | 40.2530 | 73.9123 | false | policy=materialized; observed_cte_scan=true |
| 200 | cte_joined_outer_notmat vs cte_joined_outer | 0.5446 | 0.7009 | 1.2871 | true | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_filtered_mat vs cte_filtered | 1.7570 | 1.5628 | 0.8895 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_filtered_notmat vs cte_filtered | 1.7570 | 1.6232 | 0.9239 | true | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 1.4107 | 37.5496 | 26.6177 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 1.4107 | 1.5546 | 1.1020 | true | policy=not_materialized; observed_cte_scan=false |
| 20000 | cte_filtered_mat vs cte_filtered | 5.9134 | 6.5758 | 1.1120 | false | policy=materialized; observed_cte_scan=true |
| 20000 | cte_filtered_notmat vs cte_filtered | 5.9134 | 5.8440 | 0.9883 | true | policy=not_materialized; observed_cte_scan=false |
| 20000 | cte_joined_outer_mat vs cte_joined_outer | 6.1904 | 42.0132 | 6.7868 | false | policy=materialized; observed_cte_scan=true |
| 20000 | cte_joined_outer_notmat vs cte_joined_outer | 6.1904 | 5.6211 | 0.9080 | true | policy=not_materialized; observed_cte_scan=false |
| 100000 | cte_filtered_mat vs cte_filtered | 37.8795 | 41.9282 | 1.1069 | false | policy=materialized; observed_cte_scan=true |
| 100000 | cte_filtered_notmat vs cte_filtered | 37.8795 | 38.4509 | 1.0151 | true | policy=not_materialized; observed_cte_scan=false |
| 100000 | cte_joined_outer_mat vs cte_joined_outer | 37.6850 | 53.5020 | 1.4197 | false | policy=materialized; observed_cte_scan=true |
| 100000 | cte_joined_outer_notmat vs cte_joined_outer | 37.6850 | 36.8440 | 0.9777 | true | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_filtered_mat vs cte_filtered | 48.0588 | 48.8781 | 1.0170 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_filtered_notmat vs cte_filtered | 48.0588 | 45.6820 | 0.9505 | true | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 45.0621 | 63.7415 | 1.4145 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 45.0621 | 42.5202 | 0.9436 | true | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 0.4576 | 1 | `2beb12c2700fe6cd` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 2.2153 | 1 | `ae8e946d87d174f1` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | 2.3980 | 2 | `aef241ae98ffed1c` | policy=default; observed_cte_scan=true | valid |
| default | 2 | 100% | 2.3958 | 2 | `aef241ae98ffed1c` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 1% | 2.5367 | 4 | `4052ca75f43959fb` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 2.4807 | 4 | `4052ca75f43959fb` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | 2.8126 | 8 | `836170c2201c675d` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 100% | 2.9185 | 8 | `836170c2201c675d` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 1% | 3.0468 | 16 | `9e8e5819d256402e` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 3.3187 | 16 | `9e8e5819d256402e` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 2.1583 | 1 | `3b31b8c96a2778e9` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 2.4804 | 1 | `3b31b8c96a2778e9` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | 2.5836 | 2 | `aef241ae98ffed1c` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 100% | 2.6858 | 2 | `aef241ae98ffed1c` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 1% | 2.6090 | 4 | `4052ca75f43959fb` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 2.4531 | 4 | `4052ca75f43959fb` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | 2.4965 | 8 | `836170c2201c675d` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 100% | 3.2429 | 8 | `836170c2201c675d` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 1% | 3.2646 | 16 | `9e8e5819d256402e` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 3.3909 | 16 | `9e8e5819d256402e` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 0.4459 | 1 | `2beb12c2700fe6cd` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 2.3161 | 1 | `ae8e946d87d174f1` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | 0.7016 | 2 | `b5a88f6403cacc71` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 100% | 4.9255 | 2 | `97c514398ae7345d` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 1% | 1.0949 | 4 | `2b6ccce2a72df800` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 7.0942 | 4 | `9bc6945fc91b283b` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | 1.3449 | 8 | `caf493c80d09c2fe` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 100% | 8.9133 | 8 | `b09c2c0768541f93` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 1% | 2.2689 | 16 | `7cb8b15715a86c9f` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 14.2454 | 16 | `a1e0601200b9912d` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 4.7170 | false |
| notmat | 1 | 1% | 0.9745 | true |
| mat | 1 | 100% | 1.1197 | false |
| notmat | 1 | 100% | 1.0455 | true |
| mat | 2 | 1% | 1.0774 | true |
| notmat | 2 | 1% | 0.2926 | false |
| mat | 2 | 100% | 1.1210 | true |
| notmat | 2 | 100% | 2.0559 | false |
| mat | 4 | 1% | 1.0285 | true |
| notmat | 4 | 1% | 0.4316 | false |
| mat | 4 | 100% | 0.9889 | true |
| notmat | 4 | 100% | 2.8597 | false |
| mat | 8 | 1% | 0.8876 | true |
| notmat | 8 | 1% | 0.4782 | false |
| mat | 8 | 100% | 1.1112 | true |
| notmat | 8 | 100% | 3.0541 | false |
| mat | 16 | 1% | 1.0715 | true |
| notmat | 16 | 1% | 0.7447 | false |
| mat | 16 | 100% | 1.0218 | true |
| notmat | 16 | 100% | 4.2925 | false |

## Repeated-design position control

| first median ms | last median ms | ratio (last/first) | same plan |
|---|---|---|---|
| 0.778366 | 0.623130 | 0.8006 | true |

## Execution order

- order seed `42528501`, 5 interleaved rounds
  1. family=bag;stmt=join_on;threshold=200
  2. family=reuse;policy=notmat;refs=8;consumer=1%
  3. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  4. family=reuse;policy=default;refs=2;consumer=1%
  5. family=bag;stmt=cte_filtered;threshold=180000
  6. family=reuse;policy=notmat;refs=2;consumer=100%
  7. family=bag;stmt=cte_filtered_mat;threshold=20000
  8. family=bag;stmt=cte_filtered_notmat;threshold=2000
  9. family=bag;stmt=cte_filtered;threshold=20000
  10. family=bag;stmt=cte_filtered_notmat;threshold=20000
  11. family=bag;stmt=filtered_derived;threshold=200
  12. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  13. family=bag;stmt=cte_joined_outer;threshold=200
  14. family=bag;stmt=filtered_derived;threshold=100000
  15. family=reuse;policy=notmat;refs=16;consumer=100%
  16. family=bag;stmt=cte_filtered_mat;threshold=2000
  17. family=bag;stmt=join_on;threshold=20000
  18. family=bag;stmt=comma_join;threshold=20
  19. family=reuse;policy=default;refs=2;consumer=100%
  20. family=reuse;policy=notmat;refs=8;consumer=100%
  21. family=bag;stmt=cte_joined_outer_notmat;threshold=200
  22. family=reuse;policy=default;refs=1;consumer=100%
  23. family=bag;stmt=cte_filtered_mat;threshold=180000
  24. family=bag;stmt=cte_filtered_notmat;threshold=200
  25. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  26. family=reuse;policy=default;refs=4;consumer=1%
  27. family=bag;stmt=cte_joined_outer_mat;threshold=20
  28. family=bag;stmt=cte_filtered_mat;threshold=100000
  29. family=bag;stmt=join_where;threshold=100000
  30. family=bag;stmt=cte_joined_outer_mat;threshold=20000
  31. family=reuse;policy=mat;refs=1;consumer=100%
  32. family=bag;stmt=comma_join;threshold=20000
  33. family=bag;stmt=join_where;threshold=180000
  34. family=reuse;policy=mat;refs=2;consumer=100%
  35. family=bag;stmt=cte_filtered_mat;threshold=20
  36. family=bag;stmt=join_on;threshold=180000
  37. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  38. family=bag;stmt=filtered_derived;threshold=20
  39. family=reuse;policy=mat;refs=16;consumer=100%
  40. family=bag;stmt=cte_joined_outer;threshold=2000
  41. family=bag;stmt=filtered_derived;threshold=180000
  42. family=reuse;policy=default;refs=8;consumer=1%
  43. family=bag;stmt=cte_joined_outer_notmat;threshold=100000
  44. family=bag;stmt=cte_filtered_notmat;threshold=100000
  45. family=bag;stmt=comma_join;threshold=100000
  46. family=reuse;policy=mat;refs=8;consumer=100%
  47. family=reuse;policy=mat;refs=4;consumer=100%
  48. family=bag;stmt=join_on;threshold=2000
  49. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  50. family=bag;stmt=join_where;threshold=20000
  51. family=bag;stmt=cte_filtered;threshold=100000
  52. family=bag;stmt=cte_joined_outer;threshold=100000
  53. family=bag;stmt=cte_filtered;threshold=2000
  54. family=bag;stmt=cte_filtered;threshold=20
  55. family=bag;stmt=comma_join;threshold=180000
  56. family=reuse;policy=notmat;refs=4;consumer=1%
  57. family=bag;stmt=cte_filtered;threshold=200
  58. family=bag;stmt=cte_joined_outer_mat;threshold=200
  59. family=reuse;policy=mat;refs=16;consumer=1%
  60. family=reuse;policy=mat;refs=1;consumer=1%
  61. family=bag;stmt=cte_joined_outer_notmat;threshold=20000
  62. family=bag;stmt=cte_joined_outer;threshold=20
  63. family=reuse;policy=notmat;refs=2;consumer=1%
  64. family=bag;stmt=join_where;threshold=20
  65. family=reuse;policy=notmat;refs=4;consumer=100%
  66. family=reuse;policy=default;refs=16;consumer=1%
  67. family=reuse;policy=default;refs=1;consumer=1%
  68. family=reuse;policy=notmat;refs=1;consumer=1%
  69. family=bag;stmt=comma_join;threshold=2000
  70. family=bag;stmt=cte_joined_outer;threshold=180000
  71. family=reuse;policy=mat;refs=2;consumer=1%
  72. family=reuse;policy=default;refs=16;consumer=100%
  73. family=bag;stmt=cte_filtered_mat;threshold=200
  74. family=bag;stmt=cte_joined_outer_mat;threshold=100000
  75. family=reuse;policy=notmat;refs=1;consumer=100%
  76. family=reuse;policy=mat;refs=8;consumer=1%
  77. family=bag;stmt=join_where;threshold=200
  78. family=reuse;policy=default;refs=4;consumer=100%
  79. family=reuse;policy=notmat;refs=16;consumer=1%
  80. family=bag;stmt=join_where;threshold=2000
  81. family=bag;stmt=comma_join;threshold=200
  82. family=bag;stmt=cte_filtered_notmat;threshold=180000
  83. family=reuse;policy=mat;refs=4;consumer=1%
  84. family=bag;stmt=join_on;threshold=100000
  85. family=bag;stmt=filtered_derived;threshold=2000
  86. family=bag;stmt=cte_filtered_notmat;threshold=20
  87. family=reuse;policy=default;refs=8;consumer=100%
  88. family=bag;stmt=join_on;threshold=20
  89. family=bag;stmt=filtered_derived;threshold=20000
  90. family=bag;stmt=cte_joined_outer;threshold=20000
  91. family=bag;stmt=join_on;threshold=200;position=last

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
| family=bag;stmt=join_on;threshold=200;position=last | j02_join_on | valid | `e9b5c491684a318e` |  |
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
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `a1e0601200b9912d` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `2beb12c2700fe6cd` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `ae8e946d87d174f1` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=2;consumer=1% | r_notmat_n2_c1 | valid | `b5a88f6403cacc71` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=2;consumer=100% | r_notmat_n2_c100 | valid | `97c514398ae7345d` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `2b6ccce2a72df800` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `9bc6945fc91b283b` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=8;consumer=1% | r_notmat_n8_c1 | valid | `caf493c80d09c2fe` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=8;consumer=100% | r_notmat_n8_c100 | valid | `b09c2c0768541f93` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- one host, single client: no independent hardware, real-network or equal-total-budget comparison
- no concurrent writers: the isolated read contract only; a matched blended read/write phase is a recorded gap
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
