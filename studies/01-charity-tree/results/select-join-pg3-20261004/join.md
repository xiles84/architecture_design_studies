# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 91 cells, 91 valid, engine `postgres`, topology `single`, scale `cte`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-pg3-20261004`
- run tag: `run/study-01/select-join-pg3-20261004`
- environment: `host-zenbook-ux5406sa`
- repo commit: `7e148174b5211225f6b4f243fc02e0de417a0586`
- engine version: `PostgreSQL 17.11 (Debian 17.11-1.pgdg13+2) on x86_64-pc-linux-gnu, compiled by gcc (Debian 14.2.0-19) 14.2.0, 64-bit`
- engine image: `localhost/charitybench-select:select-join-pg3-20261004` @ `7296f210ae81031ec955dbad9a67a84fe958572a2153b8d0826a647522904dc1`
- inputs digest: `05b2d7a99e38c77c84baf413a769f531179682921a2419bf73da89153b2d00ac`
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
| 20 (0.01%) | 0.5173 | 0.4731 | 0.5433 | 0.5586 | 0.6535 | 0.5953 | 0.5549 | 47.3537 | 0.6638 | 0.5099 |
| 200 (0.1%) | 0.9129 | 0.7535 | 1.1272 | 0.8483 | 0.8686 | 0.9135 | 0.7998 | 45.4371 | 1.1907 | 0.9139 |
| 2000 (1%) | 1.9756 | 1.9229 | 1.9151 | 2.1902 | 2.2407 | 2.1657 | 1.8898 | 51.8563 | 1.9320 | 2.2905 |
| 20000 (10%) | 7.6702 | 9.1660 | 7.6269 | 7.4979 | 7.8886 | 7.0544 | 7.3033 | 52.3460 | 8.6984 | 7.3755 |
| 100000 (50%) | 54.1921 | 53.4579 | 49.7021 | 60.2255 | 62.9446 | 58.0338 | 58.9167 | 72.6632 | 61.9921 | 60.0191 |
| 180000 (90%) | 68.7031 | 62.4330 | 67.8667 | 60.0565 | 71.2601 | 51.6417 | 66.3503 | 84.1167 | 58.6311 | 74.2699 |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | join_on vs join_where | 0.5173 | 0.4731 | 0.9146 | true |
| 20 | filtered_derived vs join_where | 0.5173 | 0.5433 | 1.0503 | true |
| 20 | cte_filtered vs join_where | 0.5173 | 0.5586 | 1.0799 | true |
| 20 | cte_filtered_mat vs join_where | 0.5173 | 0.6535 | 1.2634 | false |
| 20 | cte_filtered_notmat vs join_where | 0.5173 | 0.5953 | 1.1509 | true |
| 20 | cte_joined_outer vs join_where | 0.5173 | 0.5549 | 1.0727 | true |
| 20 | cte_joined_outer_mat vs join_where | 0.5173 | 47.3537 | 91.5453 | false |
| 20 | cte_joined_outer_notmat vs join_where | 0.5173 | 0.6638 | 1.2832 | true |
| 20 | comma_join vs join_where | 0.5173 | 0.5099 | 0.9857 | true |
| 200 | join_on vs join_where | 0.9129 | 0.7535 | 0.8254 | true |
| 200 | filtered_derived vs join_where | 0.9129 | 1.1272 | 1.2348 | true |
| 200 | cte_filtered vs join_where | 0.9129 | 0.8483 | 0.9293 | true |
| 200 | cte_filtered_mat vs join_where | 0.9129 | 0.8686 | 0.9515 | false |
| 200 | cte_filtered_notmat vs join_where | 0.9129 | 0.9135 | 1.0007 | true |
| 200 | cte_joined_outer vs join_where | 0.9129 | 0.7998 | 0.8762 | true |
| 200 | cte_joined_outer_mat vs join_where | 0.9129 | 45.4371 | 49.7734 | false |
| 200 | cte_joined_outer_notmat vs join_where | 0.9129 | 1.1907 | 1.3044 | true |
| 200 | comma_join vs join_where | 0.9129 | 0.9139 | 1.0011 | true |
| 2000 | join_on vs join_where | 1.9756 | 1.9229 | 0.9733 | true |
| 2000 | filtered_derived vs join_where | 1.9756 | 1.9151 | 0.9694 | true |
| 2000 | cte_filtered vs join_where | 1.9756 | 2.1902 | 1.1086 | true |
| 2000 | cte_filtered_mat vs join_where | 1.9756 | 2.2407 | 1.1342 | false |
| 2000 | cte_filtered_notmat vs join_where | 1.9756 | 2.1657 | 1.0962 | true |
| 2000 | cte_joined_outer vs join_where | 1.9756 | 1.8898 | 0.9565 | true |
| 2000 | cte_joined_outer_mat vs join_where | 1.9756 | 51.8563 | 26.2481 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 1.9756 | 1.9320 | 0.9779 | true |
| 2000 | comma_join vs join_where | 1.9756 | 2.2905 | 1.1594 | true |
| 20000 | join_on vs join_where | 7.6702 | 9.1660 | 1.1950 | true |
| 20000 | filtered_derived vs join_where | 7.6702 | 7.6269 | 0.9944 | true |
| 20000 | cte_filtered vs join_where | 7.6702 | 7.4979 | 0.9775 | true |
| 20000 | cte_filtered_mat vs join_where | 7.6702 | 7.8886 | 1.0285 | false |
| 20000 | cte_filtered_notmat vs join_where | 7.6702 | 7.0544 | 0.9197 | true |
| 20000 | cte_joined_outer vs join_where | 7.6702 | 7.3033 | 0.9522 | true |
| 20000 | cte_joined_outer_mat vs join_where | 7.6702 | 52.3460 | 6.8246 | false |
| 20000 | cte_joined_outer_notmat vs join_where | 7.6702 | 8.6984 | 1.1340 | true |
| 20000 | comma_join vs join_where | 7.6702 | 7.3755 | 0.9616 | true |
| 100000 | join_on vs join_where | 54.1921 | 53.4579 | 0.9865 | true |
| 100000 | filtered_derived vs join_where | 54.1921 | 49.7021 | 0.9171 | true |
| 100000 | cte_filtered vs join_where | 54.1921 | 60.2255 | 1.1113 | true |
| 100000 | cte_filtered_mat vs join_where | 54.1921 | 62.9446 | 1.1615 | false |
| 100000 | cte_filtered_notmat vs join_where | 54.1921 | 58.0338 | 1.0709 | true |
| 100000 | cte_joined_outer vs join_where | 54.1921 | 58.9167 | 1.0872 | true |
| 100000 | cte_joined_outer_mat vs join_where | 54.1921 | 72.6632 | 1.3408 | false |
| 100000 | cte_joined_outer_notmat vs join_where | 54.1921 | 61.9921 | 1.1439 | true |
| 100000 | comma_join vs join_where | 54.1921 | 60.0191 | 1.1075 | true |
| 180000 | join_on vs join_where | 68.7031 | 62.4330 | 0.9087 | true |
| 180000 | filtered_derived vs join_where | 68.7031 | 67.8667 | 0.9878 | true |
| 180000 | cte_filtered vs join_where | 68.7031 | 60.0565 | 0.8741 | true |
| 180000 | cte_filtered_mat vs join_where | 68.7031 | 71.2601 | 1.0372 | false |
| 180000 | cte_filtered_notmat vs join_where | 68.7031 | 51.6417 | 0.7517 | true |
| 180000 | cte_joined_outer vs join_where | 68.7031 | 66.3503 | 0.9658 | true |
| 180000 | cte_joined_outer_mat vs join_where | 68.7031 | 84.1167 | 1.2244 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 68.7031 | 58.6311 | 0.8534 | true |
| 180000 | comma_join vs join_where | 68.7031 | 74.2699 | 1.0810 | true |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_filtered_mat vs cte_filtered | 0.5586 | 0.6535 | 1.1700 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_filtered_notmat vs cte_filtered | 0.5586 | 0.5953 | 1.0657 | true | policy=not_materialized; observed_cte_scan=false |
| 20 | cte_joined_outer_mat vs cte_joined_outer | 0.5549 | 47.3537 | 85.3413 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 0.5549 | 0.6638 | 1.1962 | true | policy=not_materialized; observed_cte_scan=false |
| 200 | cte_filtered_mat vs cte_filtered | 0.8483 | 0.8686 | 1.0239 | false | policy=materialized; observed_cte_scan=true |
| 200 | cte_filtered_notmat vs cte_filtered | 0.8483 | 0.9135 | 1.0768 | true | policy=not_materialized; observed_cte_scan=false |
| 200 | cte_joined_outer_mat vs cte_joined_outer | 0.7998 | 45.4371 | 56.8073 | false | policy=materialized; observed_cte_scan=true |
| 200 | cte_joined_outer_notmat vs cte_joined_outer | 0.7998 | 1.1907 | 1.4887 | true | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_filtered_mat vs cte_filtered | 2.1902 | 2.2407 | 1.0230 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_filtered_notmat vs cte_filtered | 2.1902 | 2.1657 | 0.9888 | true | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 1.8898 | 51.8563 | 27.4406 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 1.8898 | 1.9320 | 1.0223 | true | policy=not_materialized; observed_cte_scan=false |
| 20000 | cte_filtered_mat vs cte_filtered | 7.4979 | 7.8886 | 1.0521 | false | policy=materialized; observed_cte_scan=true |
| 20000 | cte_filtered_notmat vs cte_filtered | 7.4979 | 7.0544 | 0.9408 | true | policy=not_materialized; observed_cte_scan=false |
| 20000 | cte_joined_outer_mat vs cte_joined_outer | 7.3033 | 52.3460 | 7.1675 | false | policy=materialized; observed_cte_scan=true |
| 20000 | cte_joined_outer_notmat vs cte_joined_outer | 7.3033 | 8.6984 | 1.1910 | true | policy=not_materialized; observed_cte_scan=false |
| 100000 | cte_filtered_mat vs cte_filtered | 60.2255 | 62.9446 | 1.0451 | false | policy=materialized; observed_cte_scan=true |
| 100000 | cte_filtered_notmat vs cte_filtered | 60.2255 | 58.0338 | 0.9636 | true | policy=not_materialized; observed_cte_scan=false |
| 100000 | cte_joined_outer_mat vs cte_joined_outer | 58.9167 | 72.6632 | 1.2333 | false | policy=materialized; observed_cte_scan=true |
| 100000 | cte_joined_outer_notmat vs cte_joined_outer | 58.9167 | 61.9921 | 1.0522 | true | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_filtered_mat vs cte_filtered | 60.0565 | 71.2601 | 1.1866 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_filtered_notmat vs cte_filtered | 60.0565 | 51.6417 | 0.8599 | true | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 66.3503 | 84.1167 | 1.2678 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 66.3503 | 58.6311 | 0.8837 | true | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 0.7223 | 1 | `2beb12c2700fe6cd` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 2.6865 | 1 | `ae8e946d87d174f1` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | 3.0633 | 2 | `aef241ae98ffed1c` | policy=default; observed_cte_scan=true | valid |
| default | 2 | 100% | 2.9653 | 2 | `aef241ae98ffed1c` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 1% | 3.0217 | 4 | `4052ca75f43959fb` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 3.3216 | 4 | `4052ca75f43959fb` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | 3.4752 | 8 | `836170c2201c675d` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 100% | 3.9767 | 8 | `836170c2201c675d` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 1% | 4.0233 | 16 | `9e8e5819d256402e` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 4.6956 | 16 | `9e8e5819d256402e` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 3.2872 | 1 | `3b31b8c96a2778e9` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 3.0881 | 1 | `3b31b8c96a2778e9` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | 3.8939 | 2 | `aef241ae98ffed1c` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 100% | 3.1538 | 2 | `aef241ae98ffed1c` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 1% | 3.1626 | 4 | `4052ca75f43959fb` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 4.4975 | 4 | `4052ca75f43959fb` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | 3.3360 | 8 | `836170c2201c675d` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 100% | 3.5456 | 8 | `836170c2201c675d` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 1% | 3.8049 | 16 | `9e8e5819d256402e` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 4.4203 | 16 | `9e8e5819d256402e` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 0.6402 | 1 | `2beb12c2700fe6cd` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 2.8040 | 1 | `ae8e946d87d174f1` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | 0.7837 | 2 | `b5a88f6403cacc71` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 100% | 5.3479 | 2 | `3dec56afda6778f7` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 1% | 1.3533 | 4 | `2b6ccce2a72df800` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 8.9609 | 4 | `9bc6945fc91b283b` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | 1.8828 | 8 | `caf493c80d09c2fe` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 100% | 14.2132 | 8 | `3f71bddec758345a` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 1% | 2.9672 | 16 | `7cb8b15715a86c9f` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 19.2252 | 16 | `a1e0601200b9912d` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 4.5509 | false |
| notmat | 1 | 1% | 0.8863 | true |
| mat | 1 | 100% | 1.1495 | false |
| notmat | 1 | 100% | 1.0438 | true |
| mat | 2 | 1% | 1.2711 | true |
| notmat | 2 | 1% | 0.2558 | false |
| mat | 2 | 100% | 1.0635 | true |
| notmat | 2 | 100% | 1.8035 | false |
| mat | 4 | 1% | 1.0466 | true |
| notmat | 4 | 1% | 0.4479 | false |
| mat | 4 | 100% | 1.3540 | true |
| notmat | 4 | 100% | 2.6977 | false |
| mat | 8 | 1% | 0.9600 | true |
| notmat | 8 | 1% | 0.5418 | false |
| mat | 8 | 100% | 0.8916 | true |
| notmat | 8 | 100% | 3.5741 | false |
| mat | 16 | 1% | 0.9457 | true |
| notmat | 16 | 1% | 0.7375 | false |
| mat | 16 | 100% | 0.9414 | true |
| notmat | 16 | 100% | 4.0943 | false |

## Repeated-design position control

| first median ms | last median ms | ratio (last/first) | same plan |
|---|---|---|---|
| 0.653533 | 0.395802 | 0.6056 | true |

## Execution order

- order seed `42237135`, 5 interleaved rounds
  1. family=bag;stmt=cte_filtered_mat;threshold=20
  2. family=bag;stmt=join_on;threshold=2000
  3. family=reuse;policy=default;refs=8;consumer=1%
  4. family=reuse;policy=notmat;refs=8;consumer=100%
  5. family=bag;stmt=filtered_derived;threshold=100000
  6. family=bag;stmt=cte_joined_outer;threshold=20000
  7. family=reuse;policy=notmat;refs=2;consumer=100%
  8. family=bag;stmt=join_where;threshold=20
  9. family=bag;stmt=join_on;threshold=100000
  10. family=bag;stmt=filtered_derived;threshold=180000
  11. family=bag;stmt=cte_filtered_mat;threshold=20000
  12. family=bag;stmt=comma_join;threshold=20000
  13. family=reuse;policy=default;refs=2;consumer=1%
  14. family=reuse;policy=default;refs=4;consumer=100%
  15. family=bag;stmt=cte_filtered;threshold=200
  16. family=bag;stmt=cte_joined_outer;threshold=200
  17. family=bag;stmt=comma_join;threshold=180000
  18. family=bag;stmt=cte_joined_outer_mat;threshold=20000
  19. family=reuse;policy=mat;refs=4;consumer=1%
  20. family=bag;stmt=join_where;threshold=100000
  21. family=bag;stmt=cte_joined_outer_mat;threshold=200
  22. family=bag;stmt=cte_joined_outer;threshold=2000
  23. family=bag;stmt=filtered_derived;threshold=20000
  24. family=reuse;policy=default;refs=1;consumer=1%
  25. family=bag;stmt=cte_joined_outer_notmat;threshold=100000
  26. family=bag;stmt=filtered_derived;threshold=200
  27. family=bag;stmt=cte_joined_outer;threshold=180000
  28. family=reuse;policy=mat;refs=4;consumer=100%
  29. family=bag;stmt=cte_joined_outer_mat;threshold=20
  30. family=bag;stmt=cte_filtered;threshold=180000
  31. family=bag;stmt=cte_filtered_notmat;threshold=20
  32. family=reuse;policy=notmat;refs=2;consumer=1%
  33. family=reuse;policy=default;refs=4;consumer=1%
  34. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  35. family=reuse;policy=mat;refs=2;consumer=100%
  36. family=bag;stmt=cte_filtered;threshold=20000
  37. family=bag;stmt=cte_joined_outer_notmat;threshold=200
  38. family=bag;stmt=cte_filtered_notmat;threshold=180000
  39. family=reuse;policy=default;refs=8;consumer=100%
  40. family=reuse;policy=mat;refs=8;consumer=100%
  41. family=bag;stmt=cte_filtered;threshold=100000
  42. family=bag;stmt=cte_filtered_notmat;threshold=20000
  43. family=bag;stmt=cte_filtered_notmat;threshold=2000
  44. family=bag;stmt=join_where;threshold=20000
  45. family=reuse;policy=default;refs=16;consumer=100%
  46. family=bag;stmt=comma_join;threshold=20
  47. family=bag;stmt=comma_join;threshold=100000
  48. family=bag;stmt=comma_join;threshold=2000
  49. family=reuse;policy=notmat;refs=1;consumer=1%
  50. family=reuse;policy=notmat;refs=1;consumer=100%
  51. family=bag;stmt=join_on;threshold=180000
  52. family=reuse;policy=notmat;refs=4;consumer=1%
  53. family=reuse;policy=mat;refs=1;consumer=100%
  54. family=bag;stmt=filtered_derived;threshold=20
  55. family=bag;stmt=cte_filtered_mat;threshold=2000
  56. family=bag;stmt=cte_joined_outer_mat;threshold=100000
  57. family=reuse;policy=mat;refs=2;consumer=1%
  58. family=bag;stmt=cte_filtered;threshold=20
  59. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  60. family=bag;stmt=join_where;threshold=2000
  61. family=reuse;policy=default;refs=2;consumer=100%
  62. family=bag;stmt=cte_filtered_notmat;threshold=200
  63. family=reuse;policy=default;refs=16;consumer=1%
  64. family=reuse;policy=mat;refs=16;consumer=100%
  65. family=bag;stmt=cte_joined_outer_notmat;threshold=20000
  66. family=bag;stmt=comma_join;threshold=200
  67. family=reuse;policy=notmat;refs=4;consumer=100%
  68. family=bag;stmt=cte_joined_outer;threshold=20
  69. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  70. family=bag;stmt=join_where;threshold=200
  71. family=reuse;policy=notmat;refs=16;consumer=1%
  72. family=bag;stmt=cte_filtered_mat;threshold=180000
  73. family=reuse;policy=mat;refs=1;consumer=1%
  74. family=bag;stmt=cte_filtered_notmat;threshold=100000
  75. family=reuse;policy=notmat;refs=16;consumer=100%
  76. family=bag;stmt=cte_filtered_mat;threshold=100000
  77. family=bag;stmt=cte_joined_outer;threshold=100000
  78. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  79. family=bag;stmt=join_on;threshold=200
  80. family=reuse;policy=notmat;refs=8;consumer=1%
  81. family=bag;stmt=join_on;threshold=20
  82. family=reuse;policy=default;refs=1;consumer=100%
  83. family=bag;stmt=join_where;threshold=180000
  84. family=bag;stmt=cte_filtered;threshold=2000
  85. family=bag;stmt=filtered_derived;threshold=2000
  86. family=bag;stmt=cte_filtered_mat;threshold=200
  87. family=reuse;policy=mat;refs=16;consumer=1%
  88. family=reuse;policy=mat;refs=8;consumer=1%
  89. family=bag;stmt=join_on;threshold=20000
  90. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  91. family=bag;stmt=cte_filtered_mat;threshold=20;position=last

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
| family=bag;stmt=cte_filtered_mat;threshold=20;position=last | j04m_cte_filtered_mat | valid | `fef2a456dbb4dbce` | policy=materialized; observed_cte_scan=true |
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
| family=reuse;policy=notmat;refs=2;consumer=100% | r_notmat_n2_c100 | valid | `3dec56afda6778f7` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `2b6ccce2a72df800` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `9bc6945fc91b283b` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=8;consumer=1% | r_notmat_n8_c1 | valid | `caf493c80d09c2fe` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=8;consumer=100% | r_notmat_n8_c100 | valid | `3f71bddec758345a` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- one host, single client: no independent hardware, real-network or equal-total-budget comparison
- no concurrent writers: the isolated read contract only; a matched blended read/write phase is a recorded gap
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
