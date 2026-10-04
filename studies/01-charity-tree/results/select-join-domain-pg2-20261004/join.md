# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 91 cells, 91 valid, engine `postgres`, topology `single`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-domain-pg2-20261004`
- run tag: `run/study-01/select-join-domain-pg2-20261004`
- environment: `host-zenbook-ux5406sa`
- repo commit: `ea77bc5c677aa9582a2646b597ea0211eac9d0e0`
- engine version: `PostgreSQL 17.11 (Debian 17.11-1.pgdg13+2) on x86_64-pc-linux-gnu, compiled by gcc (Debian 14.2.0-19) 14.2.0, 64-bit`
- database engine image: `docker.io/library/postgres:17.11` @ `7296f210ae81031ec955dbad9a67a84fe958572a2153b8d0826a647522904dc1`
- benchmark client image: `localhost/charitybench-select:select-join-domain-pg2-20261004` @ `fa3606dc9ba1e4169c0124c33434b7aa827570eb10aa634f479901ca10b2f0be`
- inputs digest: `d6de501895a3181190de1a1825c5965c62e9b99f7dde692b1a10c467890d3bc2`
- dataset fingerprint: `fcf26e08f377ec51`
- client CPU quota: `100000 100000`, memory: `2147483648`
- setting `client_workers`: 1 sequential reader; 0 writers
- setting `connection_pool`: 2 = 1 reader/setup + 1 reserve
- setting `customers`: 2000
- setting `dataset`: study01-person-donation-charity
- setting `distributed_rpc`: N/A on PostgreSQL single-node
- setting `fact_count`: 200000
- setting `index_set`: child rank index and parent-key index, primary keys; no per-spelling index
- setting `prepared_plan_policy`: pgx simple protocol; no prepared statements
- setting `scale`: cte-domain
- setting `selectivity_mapping`: rank <= threshold; rank is a deterministic shuffle of 1..fact_count
- setting `statement_budget_ms`: 120000
- setting `statistics`: ANALYZE after deterministic load
- setting `thresholds`: [20 200 2000 20000 100000 180000]
- setting `topology`: single
- setting `trials_per_cell`: 5
- setting `whole_run_count`: 1; no replicated winner in this file

## Fixture and selectivity mapping

- arm: `study01-person-donation-charity`
- fact rows: 200000
- parent rows: 2000
- rank thresholds: [20 200 2000 20000 100000 180000]

The predicate is `rank <= threshold`, where `rank` is a deterministic shuffle of 1..200000; `rank` is a labelled proxy for the teaching selective predicate, not the same question. The `cte`/`cte-small` arm is the synthetic customers/facts microfixture; the `cte-domain` arm is Study 01's person/donation fixture with the person's charity projected and the views binding it to this catalogue.

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
| 20 (0.01%) | 0.8339 | 0.9895 | 0.6565 | 0.8937 | 1.0511 | 1.1268 | 0.6325 | 96.0967 | 0.9156 | 0.7429 |
| 200 (0.1%) | 1.1568 | 1.0699 | 1.3723 | 1.1353 | 1.4749 | 1.2630 | 1.1862 | 69.8653 | 1.3588 | 1.4749 |
| 2000 (1%) | 3.2849 | 2.2548 | 3.0841 | 3.5722 | 4.6169 | 3.0845 | 2.6834 | 71.4844 | 3.4040 | 3.5948 |
| 20000 (10%) | 17.1928 | 15.6357 | 15.7120 | 16.9754 | 14.4377 | 15.7976 | 14.6243 | 91.0315 | 16.7984 | 21.3029 |
| 100000 (50%) | 85.7644 | 90.2344 | 82.7441 | 81.8470 | 80.1030 | 79.9130 | 73.7816 | 151.0271 | 84.1068 | 85.7783 |
| 180000 (90%) | 98.1927 | 98.4095 | 101.6153 | 173.9902 | 139.4390 | 125.3316 | 130.8156 | 149.2009 | 92.3036 | 130.0475 |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | join_on vs join_where | 0.8339 | 0.9895 | 1.1866 | true |
| 20 | filtered_derived vs join_where | 0.8339 | 0.6565 | 0.7873 | true |
| 20 | cte_filtered vs join_where | 0.8339 | 0.8937 | 1.0717 | true |
| 20 | cte_filtered_mat vs join_where | 0.8339 | 1.0511 | 1.2605 | false |
| 20 | cte_filtered_notmat vs join_where | 0.8339 | 1.1268 | 1.3512 | true |
| 20 | cte_joined_outer vs join_where | 0.8339 | 0.6325 | 0.7584 | true |
| 20 | cte_joined_outer_mat vs join_where | 0.8339 | 96.0967 | 115.2330 | false |
| 20 | cte_joined_outer_notmat vs join_where | 0.8339 | 0.9156 | 1.0979 | true |
| 20 | comma_join vs join_where | 0.8339 | 0.7429 | 0.8908 | true |
| 200 | join_on vs join_where | 1.1568 | 1.0699 | 0.9249 | true |
| 200 | filtered_derived vs join_where | 1.1568 | 1.3723 | 1.1863 | true |
| 200 | cte_filtered vs join_where | 1.1568 | 1.1353 | 0.9814 | true |
| 200 | cte_filtered_mat vs join_where | 1.1568 | 1.4749 | 1.2750 | false |
| 200 | cte_filtered_notmat vs join_where | 1.1568 | 1.2630 | 1.0919 | true |
| 200 | cte_joined_outer vs join_where | 1.1568 | 1.1862 | 1.0254 | true |
| 200 | cte_joined_outer_mat vs join_where | 1.1568 | 69.8653 | 60.3971 | false |
| 200 | cte_joined_outer_notmat vs join_where | 1.1568 | 1.3588 | 1.1747 | true |
| 200 | comma_join vs join_where | 1.1568 | 1.4749 | 1.2750 | true |
| 2000 | join_on vs join_where | 3.2849 | 2.2548 | 0.6864 | true |
| 2000 | filtered_derived vs join_where | 3.2849 | 3.0841 | 0.9389 | true |
| 2000 | cte_filtered vs join_where | 3.2849 | 3.5722 | 1.0875 | true |
| 2000 | cte_filtered_mat vs join_where | 3.2849 | 4.6169 | 1.4055 | false |
| 2000 | cte_filtered_notmat vs join_where | 3.2849 | 3.0845 | 0.9390 | true |
| 2000 | cte_joined_outer vs join_where | 3.2849 | 2.6834 | 0.8169 | true |
| 2000 | cte_joined_outer_mat vs join_where | 3.2849 | 71.4844 | 21.7614 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 3.2849 | 3.4040 | 1.0362 | true |
| 2000 | comma_join vs join_where | 3.2849 | 3.5948 | 1.0943 | true |
| 20000 | join_on vs join_where | 17.1928 | 15.6357 | 0.9094 | true |
| 20000 | filtered_derived vs join_where | 17.1928 | 15.7120 | 0.9139 | true |
| 20000 | cte_filtered vs join_where | 17.1928 | 16.9754 | 0.9874 | true |
| 20000 | cte_filtered_mat vs join_where | 17.1928 | 14.4377 | 0.8398 | false |
| 20000 | cte_filtered_notmat vs join_where | 17.1928 | 15.7976 | 0.9189 | true |
| 20000 | cte_joined_outer vs join_where | 17.1928 | 14.6243 | 0.8506 | true |
| 20000 | cte_joined_outer_mat vs join_where | 17.1928 | 91.0315 | 5.2947 | false |
| 20000 | cte_joined_outer_notmat vs join_where | 17.1928 | 16.7984 | 0.9771 | true |
| 20000 | comma_join vs join_where | 17.1928 | 21.3029 | 1.2391 | true |
| 100000 | join_on vs join_where | 85.7644 | 90.2344 | 1.0521 | true |
| 100000 | filtered_derived vs join_where | 85.7644 | 82.7441 | 0.9648 | true |
| 100000 | cte_filtered vs join_where | 85.7644 | 81.8470 | 0.9543 | true |
| 100000 | cte_filtered_mat vs join_where | 85.7644 | 80.1030 | 0.9340 | false |
| 100000 | cte_filtered_notmat vs join_where | 85.7644 | 79.9130 | 0.9318 | true |
| 100000 | cte_joined_outer vs join_where | 85.7644 | 73.7816 | 0.8603 | true |
| 100000 | cte_joined_outer_mat vs join_where | 85.7644 | 151.0271 | 1.7610 | false |
| 100000 | cte_joined_outer_notmat vs join_where | 85.7644 | 84.1068 | 0.9807 | true |
| 100000 | comma_join vs join_where | 85.7644 | 85.7783 | 1.0002 | true |
| 180000 | join_on vs join_where | 98.1927 | 98.4095 | 1.0022 | true |
| 180000 | filtered_derived vs join_where | 98.1927 | 101.6153 | 1.0349 | true |
| 180000 | cte_filtered vs join_where | 98.1927 | 173.9902 | 1.7719 | true |
| 180000 | cte_filtered_mat vs join_where | 98.1927 | 139.4390 | 1.4201 | false |
| 180000 | cte_filtered_notmat vs join_where | 98.1927 | 125.3316 | 1.2764 | true |
| 180000 | cte_joined_outer vs join_where | 98.1927 | 130.8156 | 1.3322 | true |
| 180000 | cte_joined_outer_mat vs join_where | 98.1927 | 149.2009 | 1.5195 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 98.1927 | 92.3036 | 0.9400 | true |
| 180000 | comma_join vs join_where | 98.1927 | 130.0475 | 1.3244 | true |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_filtered_mat vs cte_filtered | 0.8937 | 1.0511 | 1.1761 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_filtered_notmat vs cte_filtered | 0.8937 | 1.1268 | 1.2608 | true | policy=not_materialized; observed_cte_scan=false |
| 20 | cte_joined_outer_mat vs cte_joined_outer | 0.6325 | 96.0967 | 151.9410 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 0.6325 | 0.9156 | 1.4476 | true | policy=not_materialized; observed_cte_scan=false |
| 200 | cte_filtered_mat vs cte_filtered | 1.1353 | 1.4749 | 1.2991 | false | policy=materialized; observed_cte_scan=true |
| 200 | cte_filtered_notmat vs cte_filtered | 1.1353 | 1.2630 | 1.1125 | true | policy=not_materialized; observed_cte_scan=false |
| 200 | cte_joined_outer_mat vs cte_joined_outer | 1.1862 | 69.8653 | 58.8999 | false | policy=materialized; observed_cte_scan=true |
| 200 | cte_joined_outer_notmat vs cte_joined_outer | 1.1862 | 1.3588 | 1.1456 | true | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_filtered_mat vs cte_filtered | 3.5722 | 4.6169 | 1.2924 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_filtered_notmat vs cte_filtered | 3.5722 | 3.0845 | 0.8635 | true | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 2.6834 | 71.4844 | 26.6394 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 2.6834 | 3.4040 | 1.2685 | true | policy=not_materialized; observed_cte_scan=false |
| 20000 | cte_filtered_mat vs cte_filtered | 16.9754 | 14.4377 | 0.8505 | false | policy=materialized; observed_cte_scan=true |
| 20000 | cte_filtered_notmat vs cte_filtered | 16.9754 | 15.7976 | 0.9306 | true | policy=not_materialized; observed_cte_scan=false |
| 20000 | cte_joined_outer_mat vs cte_joined_outer | 14.6243 | 91.0315 | 6.2247 | false | policy=materialized; observed_cte_scan=true |
| 20000 | cte_joined_outer_notmat vs cte_joined_outer | 14.6243 | 16.7984 | 1.1487 | true | policy=not_materialized; observed_cte_scan=false |
| 100000 | cte_filtered_mat vs cte_filtered | 81.8470 | 80.1030 | 0.9787 | false | policy=materialized; observed_cte_scan=true |
| 100000 | cte_filtered_notmat vs cte_filtered | 81.8470 | 79.9130 | 0.9764 | true | policy=not_materialized; observed_cte_scan=false |
| 100000 | cte_joined_outer_mat vs cte_joined_outer | 73.7816 | 151.0271 | 2.0469 | false | policy=materialized; observed_cte_scan=true |
| 100000 | cte_joined_outer_notmat vs cte_joined_outer | 73.7816 | 84.1068 | 1.1399 | true | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_filtered_mat vs cte_filtered | 173.9902 | 139.4390 | 0.8014 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_filtered_notmat vs cte_filtered | 173.9902 | 125.3316 | 0.7203 | true | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 130.8156 | 149.2009 | 1.1405 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 130.8156 | 92.3036 | 0.7056 | true | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 0.9778 | 1 | `67feb88d971b83fe` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 4.6296 | 1 | `b097208aea033185` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | 3.4488 | 2 | `450416f8dc65edb3` | policy=default; observed_cte_scan=true | valid |
| default | 2 | 100% | 4.2866 | 2 | `450416f8dc65edb3` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 1% | 3.8070 | 4 | `91c3642bff5db3db` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 3.9164 | 4 | `91c3642bff5db3db` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | 5.8029 | 8 | `0c9b689bfb53ff4f` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 100% | 5.4107 | 8 | `0c9b689bfb53ff4f` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 1% | 4.4871 | 16 | `527af41dc9ea7c27` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 5.3225 | 16 | `527af41dc9ea7c27` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 4.6471 | 1 | `9384b48eeba05ee5` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 3.5360 | 1 | `9384b48eeba05ee5` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | 3.7076 | 2 | `450416f8dc65edb3` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 100% | 3.7604 | 2 | `450416f8dc65edb3` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 1% | 5.3971 | 4 | `91c3642bff5db3db` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 4.4396 | 4 | `91c3642bff5db3db` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | 3.9185 | 8 | `0c9b689bfb53ff4f` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 100% | 4.2207 | 8 | `0c9b689bfb53ff4f` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 1% | 4.6714 | 16 | `527af41dc9ea7c27` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 7.0619 | 16 | `527af41dc9ea7c27` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 0.6936 | 1 | `67feb88d971b83fe` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 3.0746 | 1 | `b097208aea033185` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | 1.0894 | 2 | `b6d08c76c5629cb2` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 100% | 11.5578 | 2 | `cd29fa7a7d29b636` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 1% | 1.3025 | 4 | `4e9a2c71c9021d34` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 13.4820 | 4 | `3527c6c56c9bb190` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | 3.0246 | 8 | `b7419850b4ea67ad` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 100% | 16.9787 | 8 | `eb0f827e9bbc0d1f` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 1% | 3.7988 | 16 | `a586237c6b012f4c` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 27.5992 | 16 | `1754eea8ba4cc2a2` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 4.7527 | false |
| notmat | 1 | 1% | 0.7094 | true |
| mat | 1 | 100% | 0.7638 | false |
| notmat | 1 | 100% | 0.6641 | true |
| mat | 2 | 1% | 1.0750 | true |
| notmat | 2 | 1% | 0.3159 | false |
| mat | 2 | 100% | 0.8772 | true |
| notmat | 2 | 100% | 2.6963 | false |
| mat | 4 | 1% | 1.4177 | true |
| notmat | 4 | 1% | 0.3421 | false |
| mat | 4 | 100% | 1.1336 | true |
| notmat | 4 | 100% | 3.4425 | false |
| mat | 8 | 1% | 0.6753 | true |
| notmat | 8 | 1% | 0.5212 | false |
| mat | 8 | 100% | 0.7801 | true |
| notmat | 8 | 100% | 3.1380 | false |
| mat | 16 | 1% | 1.0411 | true |
| notmat | 16 | 1% | 0.8466 | false |
| mat | 16 | 100% | 1.3268 | true |
| notmat | 16 | 100% | 5.1853 | false |

## Repeated-design position control

| first median ms | last median ms | ratio (last/first) | same plan |
|---|---|---|---|
| 85.764382 | 83.713543 | 0.9761 | true |

## Latency samples and derived throughput

_The reader is sequential, so `queries/s` is 1000/median_ms, a derived rate and not a concurrent-capacity test. min/max are the five sample extremes._

| axis | status | min ms | median ms | max ms | samples | derived queries/s |
|---|---|---|---|---|---|---|
| family=bag;stmt=comma_join;threshold=100000 | valid | 69.0883 | 85.7783 | 115.3295 | 5 | 11.7 |
| family=bag;stmt=comma_join;threshold=180000 | valid | 94.7375 | 130.0475 | 165.0708 | 5 | 7.7 |
| family=bag;stmt=comma_join;threshold=20 | valid | 0.5580 | 0.7429 | 1.1206 | 5 | 1346.1 |
| family=bag;stmt=comma_join;threshold=200 | valid | 0.8604 | 1.4749 | 1.5297 | 5 | 678.0 |
| family=bag;stmt=comma_join;threshold=2000 | valid | 2.0419 | 3.5948 | 4.5423 | 5 | 278.2 |
| family=bag;stmt=comma_join;threshold=20000 | valid | 9.8638 | 21.3029 | 49.0359 | 5 | 46.9 |
| family=bag;stmt=cte_filtered;threshold=100000 | valid | 64.3403 | 81.8470 | 86.3183 | 5 | 12.2 |
| family=bag;stmt=cte_filtered;threshold=180000 | valid | 91.7438 | 173.9902 | 187.4689 | 5 | 5.7 |
| family=bag;stmt=cte_filtered;threshold=20 | valid | 0.6277 | 0.8937 | 29.4533 | 5 | 1118.9 |
| family=bag;stmt=cte_filtered;threshold=200 | valid | 0.9698 | 1.1353 | 2.0312 | 5 | 880.8 |
| family=bag;stmt=cte_filtered;threshold=2000 | valid | 2.1312 | 3.5722 | 4.1845 | 5 | 279.9 |
| family=bag;stmt=cte_filtered;threshold=20000 | valid | 9.1134 | 16.9754 | 45.5138 | 5 | 58.9 |
| family=bag;stmt=cte_filtered_mat;threshold=100000 | valid | 65.8791 | 80.1030 | 111.8985 | 5 | 12.5 |
| family=bag;stmt=cte_filtered_mat;threshold=180000 | valid | 91.7417 | 139.4390 | 173.5681 | 5 | 7.2 |
| family=bag;stmt=cte_filtered_mat;threshold=20 | valid | 0.6896 | 1.0511 | 1.7834 | 5 | 951.3 |
| family=bag;stmt=cte_filtered_mat;threshold=200 | valid | 0.9041 | 1.4749 | 1.6428 | 5 | 678.0 |
| family=bag;stmt=cte_filtered_mat;threshold=2000 | valid | 2.2698 | 4.6169 | 5.2052 | 5 | 216.6 |
| family=bag;stmt=cte_filtered_mat;threshold=20000 | valid | 10.3245 | 14.4377 | 31.7953 | 5 | 69.3 |
| family=bag;stmt=cte_filtered_notmat;threshold=100000 | valid | 77.9424 | 79.9130 | 103.6390 | 5 | 12.5 |
| family=bag;stmt=cte_filtered_notmat;threshold=180000 | valid | 110.4493 | 125.3316 | 199.6787 | 5 | 8.0 |
| family=bag;stmt=cte_filtered_notmat;threshold=20 | valid | 0.7186 | 1.1268 | 1.3706 | 5 | 887.5 |
| family=bag;stmt=cte_filtered_notmat;threshold=200 | valid | 0.9520 | 1.2630 | 1.7632 | 5 | 791.7 |
| family=bag;stmt=cte_filtered_notmat;threshold=2000 | valid | 2.1130 | 3.0845 | 5.6591 | 5 | 324.2 |
| family=bag;stmt=cte_filtered_notmat;threshold=20000 | valid | 13.8444 | 15.7976 | 18.3850 | 5 | 63.3 |
| family=bag;stmt=cte_joined_outer;threshold=100000 | valid | 66.7867 | 73.7816 | 90.1703 | 5 | 13.6 |
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 104.6234 | 130.8156 | 238.1579 | 5 | 7.6 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 0.6003 | 0.6325 | 0.8800 | 5 | 1581.1 |
| family=bag;stmt=cte_joined_outer;threshold=200 | valid | 0.9540 | 1.1862 | 1.9342 | 5 | 843.0 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 2.1338 | 2.6834 | 4.1642 | 5 | 372.7 |
| family=bag;stmt=cte_joined_outer;threshold=20000 | valid | 10.6923 | 14.6243 | 23.6898 | 5 | 68.4 |
| family=bag;stmt=cte_joined_outer_mat;threshold=100000 | valid | 97.4513 | 151.0271 | 187.8415 | 5 | 6.6 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 110.0158 | 149.2009 | 202.2770 | 5 | 6.7 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 61.7461 | 96.0967 | 106.4004 | 5 | 10.4 |
| family=bag;stmt=cte_joined_outer_mat;threshold=200 | valid | 65.1559 | 69.8653 | 97.5031 | 5 | 14.3 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 61.5097 | 71.4844 | 103.8998 | 5 | 14.0 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20000 | valid | 76.1267 | 91.0315 | 114.9590 | 5 | 11.0 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=100000 | valid | 60.9998 | 84.1068 | 162.3011 | 5 | 11.9 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 88.4492 | 92.3036 | 108.9286 | 5 | 10.8 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 0.5566 | 0.9156 | 1.0215 | 5 | 1092.2 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=200 | valid | 0.8942 | 1.3588 | 1.3886 | 5 | 735.9 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 2.0884 | 3.4040 | 3.6306 | 5 | 293.8 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20000 | valid | 13.1821 | 16.7984 | 22.6156 | 5 | 59.5 |
| family=bag;stmt=filtered_derived;threshold=100000 | valid | 76.4131 | 82.7441 | 151.9572 | 5 | 12.1 |
| family=bag;stmt=filtered_derived;threshold=180000 | valid | 93.4936 | 101.6153 | 136.2286 | 5 | 9.8 |
| family=bag;stmt=filtered_derived;threshold=20 | valid | 0.5169 | 0.6565 | 0.9146 | 5 | 1523.1 |
| family=bag;stmt=filtered_derived;threshold=200 | valid | 0.9861 | 1.3723 | 1.7194 | 5 | 728.7 |
| family=bag;stmt=filtered_derived;threshold=2000 | valid | 2.2221 | 3.0841 | 3.5586 | 5 | 324.2 |
| family=bag;stmt=filtered_derived;threshold=20000 | valid | 10.2192 | 15.7120 | 36.4056 | 5 | 63.6 |
| family=bag;stmt=join_on;threshold=100000 | valid | 64.3891 | 90.2344 | 124.5505 | 5 | 11.1 |
| family=bag;stmt=join_on;threshold=180000 | valid | 84.1467 | 98.4095 | 104.0676 | 5 | 10.2 |
| family=bag;stmt=join_on;threshold=20 | valid | 0.8013 | 0.9895 | 1.3119 | 5 | 1010.6 |
| family=bag;stmt=join_on;threshold=200 | valid | 0.8209 | 1.0699 | 1.4457 | 5 | 934.7 |
| family=bag;stmt=join_on;threshold=2000 | valid | 1.9126 | 2.2548 | 4.7458 | 5 | 443.5 |
| family=bag;stmt=join_on;threshold=20000 | valid | 11.3282 | 15.6357 | 42.2879 | 5 | 64.0 |
| family=bag;stmt=join_where;threshold=100000 | valid | 76.4283 | 85.7644 | 164.8350 | 5 | 11.7 |
| family=bag;stmt=join_where;threshold=100000;position=last | valid | 73.7832 | 83.7135 | 110.0279 | 5 | 11.9 |
| family=bag;stmt=join_where;threshold=180000 | valid | 84.8729 | 98.1927 | 175.1900 | 5 | 10.2 |
| family=bag;stmt=join_where;threshold=20 | valid | 0.5323 | 0.8339 | 0.8743 | 5 | 1199.1 |
| family=bag;stmt=join_where;threshold=200 | valid | 1.0845 | 1.1568 | 1.8490 | 5 | 864.5 |
| family=bag;stmt=join_where;threshold=2000 | valid | 2.1440 | 3.2849 | 7.5977 | 5 | 304.4 |
| family=bag;stmt=join_where;threshold=20000 | valid | 12.0600 | 17.1928 | 25.3309 | 5 | 58.2 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 4.0276 | 4.4871 | 5.7746 | 5 | 222.9 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 4.9948 | 5.3225 | 5.5388 | 5 | 187.9 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 0.6934 | 0.9778 | 1.1828 | 5 | 1022.7 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 2.8544 | 4.6296 | 5.1521 | 5 | 216.0 |
| family=reuse;policy=default;refs=2;consumer=1% | valid | 3.1864 | 3.4488 | 6.0869 | 5 | 290.0 |
| family=reuse;policy=default;refs=2;consumer=100% | valid | 3.2714 | 4.2866 | 5.3071 | 5 | 233.3 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 3.5402 | 3.8070 | 5.1292 | 5 | 262.7 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 3.3929 | 3.9164 | 5.0218 | 5 | 255.3 |
| family=reuse;policy=default;refs=8;consumer=1% | valid | 3.4534 | 5.8029 | 21.0145 | 5 | 172.3 |
| family=reuse;policy=default;refs=8;consumer=100% | valid | 4.5605 | 5.4107 | 10.6767 | 5 | 184.8 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 4.2511 | 4.6714 | 6.1617 | 5 | 214.1 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 4.6204 | 7.0619 | 8.0310 | 5 | 141.6 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 3.3550 | 4.6471 | 5.9395 | 5 | 215.2 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 3.2599 | 3.5360 | 5.7003 | 5 | 282.8 |
| family=reuse;policy=mat;refs=2;consumer=1% | valid | 2.9335 | 3.7076 | 4.7798 | 5 | 269.7 |
| family=reuse;policy=mat;refs=2;consumer=100% | valid | 3.2606 | 3.7604 | 5.2953 | 5 | 265.9 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 3.9667 | 5.3971 | 6.2665 | 5 | 185.3 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 3.7325 | 4.4396 | 9.3655 | 5 | 225.2 |
| family=reuse;policy=mat;refs=8;consumer=1% | valid | 3.7477 | 3.9185 | 6.2704 | 5 | 255.2 |
| family=reuse;policy=mat;refs=8;consumer=100% | valid | 4.1148 | 4.2207 | 5.2795 | 5 | 236.9 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 3.5482 | 3.7988 | 8.7075 | 5 | 263.2 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 25.5247 | 27.5992 | 34.0545 | 5 | 36.2 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 0.6576 | 0.6936 | 0.8569 | 5 | 1441.7 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 2.6814 | 3.0746 | 8.5066 | 5 | 325.2 |
| family=reuse;policy=notmat;refs=2;consumer=1% | valid | 1.0380 | 1.0894 | 1.9218 | 5 | 917.9 |
| family=reuse;policy=notmat;refs=2;consumer=100% | valid | 9.4785 | 11.5578 | 22.1443 | 5 | 86.5 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 1.2178 | 1.3025 | 1.8939 | 5 | 767.7 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 11.4368 | 13.4820 | 15.2580 | 5 | 74.2 |
| family=reuse;policy=notmat;refs=8;consumer=1% | valid | 2.1121 | 3.0246 | 3.4289 | 5 | 330.6 |
| family=reuse;policy=notmat;refs=8;consumer=100% | valid | 16.1455 | 16.9787 | 27.0016 | 5 | 58.9 |

## Execution order

- order seed `42729011`, 5 interleaved rounds
  1. family=bag;stmt=join_where;threshold=100000
  2. family=bag;stmt=cte_joined_outer_mat;threshold=200
  3. family=reuse;policy=mat;refs=4;consumer=100%
  4. family=reuse;policy=notmat;refs=16;consumer=1%
  5. family=bag;stmt=cte_joined_outer;threshold=2000
  6. family=bag;stmt=cte_joined_outer;threshold=200
  7. family=bag;stmt=cte_joined_outer_notmat;threshold=200
  8. family=bag;stmt=filtered_derived;threshold=20
  9. family=bag;stmt=join_where;threshold=20
  10. family=bag;stmt=join_on;threshold=200
  11. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  12. family=bag;stmt=join_where;threshold=20000
  13. family=bag;stmt=cte_filtered_mat;threshold=20000
  14. family=bag;stmt=cte_filtered_notmat;threshold=20000
  15. family=bag;stmt=cte_joined_outer_mat;threshold=100000
  16. family=bag;stmt=comma_join;threshold=20000
  17. family=reuse;policy=notmat;refs=8;consumer=1%
  18. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  19. family=reuse;policy=default;refs=8;consumer=1%
  20. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  21. family=reuse;policy=mat;refs=16;consumer=100%
  22. family=bag;stmt=cte_joined_outer_notmat;threshold=20000
  23. family=bag;stmt=cte_joined_outer_mat;threshold=20
  24. family=bag;stmt=cte_filtered_mat;threshold=200
  25. family=reuse;policy=default;refs=1;consumer=100%
  26. family=bag;stmt=comma_join;threshold=200
  27. family=bag;stmt=cte_filtered_mat;threshold=180000
  28. family=bag;stmt=cte_joined_outer;threshold=180000
  29. family=reuse;policy=mat;refs=4;consumer=1%
  30. family=reuse;policy=mat;refs=1;consumer=100%
  31. family=reuse;policy=default;refs=8;consumer=100%
  32. family=reuse;policy=notmat;refs=2;consumer=100%
  33. family=bag;stmt=cte_filtered_notmat;threshold=200
  34. family=bag;stmt=comma_join;threshold=20
  35. family=reuse;policy=default;refs=2;consumer=1%
  36. family=reuse;policy=notmat;refs=1;consumer=100%
  37. family=bag;stmt=cte_filtered_mat;threshold=100000
  38. family=bag;stmt=join_on;threshold=100000
  39. family=reuse;policy=default;refs=4;consumer=1%
  40. family=bag;stmt=filtered_derived;threshold=180000
  41. family=bag;stmt=filtered_derived;threshold=20000
  42. family=bag;stmt=cte_filtered_notmat;threshold=180000
  43. family=reuse;policy=default;refs=16;consumer=100%
  44. family=bag;stmt=cte_joined_outer;threshold=20
  45. family=reuse;policy=mat;refs=8;consumer=100%
  46. family=bag;stmt=join_on;threshold=180000
  47. family=bag;stmt=cte_joined_outer;threshold=20000
  48. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  49. family=reuse;policy=default;refs=1;consumer=1%
  50. family=bag;stmt=cte_filtered_notmat;threshold=2000
  51. family=bag;stmt=comma_join;threshold=100000
  52. family=reuse;policy=notmat;refs=16;consumer=100%
  53. family=reuse;policy=default;refs=4;consumer=100%
  54. family=reuse;policy=mat;refs=16;consumer=1%
  55. family=reuse;policy=notmat;refs=4;consumer=1%
  56. family=bag;stmt=cte_filtered;threshold=2000
  57. family=reuse;policy=notmat;refs=4;consumer=100%
  58. family=bag;stmt=filtered_derived;threshold=200
  59. family=bag;stmt=filtered_derived;threshold=2000
  60. family=bag;stmt=cte_filtered_notmat;threshold=100000
  61. family=bag;stmt=join_on;threshold=20000
  62. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  63. family=bag;stmt=cte_filtered;threshold=200
  64. family=bag;stmt=cte_joined_outer_notmat;threshold=100000
  65. family=bag;stmt=cte_filtered;threshold=100000
  66. family=bag;stmt=join_where;threshold=200
  67. family=reuse;policy=default;refs=16;consumer=1%
  68. family=reuse;policy=notmat;refs=8;consumer=100%
  69. family=reuse;policy=mat;refs=8;consumer=1%
  70. family=reuse;policy=notmat;refs=1;consumer=1%
  71. family=bag;stmt=cte_joined_outer;threshold=100000
  72. family=reuse;policy=notmat;refs=2;consumer=1%
  73. family=bag;stmt=join_where;threshold=180000
  74. family=bag;stmt=cte_filtered;threshold=180000
  75. family=bag;stmt=cte_filtered_notmat;threshold=20
  76. family=bag;stmt=filtered_derived;threshold=100000
  77. family=bag;stmt=join_on;threshold=20
  78. family=reuse;policy=mat;refs=2;consumer=100%
  79. family=bag;stmt=comma_join;threshold=2000
  80. family=bag;stmt=cte_filtered_mat;threshold=20
  81. family=bag;stmt=cte_filtered;threshold=20000
  82. family=bag;stmt=cte_filtered_mat;threshold=2000
  83. family=bag;stmt=join_where;threshold=2000
  84. family=bag;stmt=join_on;threshold=2000
  85. family=bag;stmt=cte_filtered;threshold=20
  86. family=bag;stmt=cte_joined_outer_mat;threshold=20000
  87. family=reuse;policy=default;refs=2;consumer=100%
  88. family=reuse;policy=mat;refs=2;consumer=1%
  89. family=reuse;policy=mat;refs=1;consumer=1%
  90. family=bag;stmt=comma_join;threshold=180000
  91. family=bag;stmt=join_where;threshold=100000;position=last

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
| family=bag;stmt=comma_join;threshold=100000 | j06_comma_join | valid | `6f0091d8475bff3c` |  |
| family=bag;stmt=comma_join;threshold=180000 | j06_comma_join | valid | `116d58c44d68b0e2` |  |
| family=bag;stmt=comma_join;threshold=20 | j06_comma_join | valid | `9192a1cb0d423ae7` |  |
| family=bag;stmt=comma_join;threshold=200 | j06_comma_join | valid | `da1656fbf706c609` |  |
| family=bag;stmt=comma_join;threshold=2000 | j06_comma_join | valid | `9174168e6180c2f9` |  |
| family=bag;stmt=comma_join;threshold=20000 | j06_comma_join | valid | `d43dd14675befada` |  |
| family=bag;stmt=cte_filtered;threshold=100000 | j04_cte_filtered | valid | `6f0091d8475bff3c` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=180000 | j04_cte_filtered | valid | `116d58c44d68b0e2` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=20 | j04_cte_filtered | valid | `9192a1cb0d423ae7` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=200 | j04_cte_filtered | valid | `da1656fbf706c609` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=2000 | j04_cte_filtered | valid | `9174168e6180c2f9` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=20000 | j04_cte_filtered | valid | `d43dd14675befada` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_mat;threshold=100000 | j04m_cte_filtered_mat | valid | `4b4ce4c83ae7fcc8` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=180000 | j04m_cte_filtered_mat | valid | `6a10a6fd7e521978` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=20 | j04m_cte_filtered_mat | valid | `d4ca4a9d1c0fd3e3` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=200 | j04m_cte_filtered_mat | valid | `c1c51eb45bf14a85` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=2000 | j04m_cte_filtered_mat | valid | `d5f7ae93b0e0ee0b` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=20000 | j04m_cte_filtered_mat | valid | `a5a9631585a1d3ee` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_notmat;threshold=100000 | j04n_cte_filtered_notmat | valid | `6f0091d8475bff3c` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=180000 | j04n_cte_filtered_notmat | valid | `116d58c44d68b0e2` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=20 | j04n_cte_filtered_notmat | valid | `9192a1cb0d423ae7` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=200 | j04n_cte_filtered_notmat | valid | `da1656fbf706c609` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=2000 | j04n_cte_filtered_notmat | valid | `9174168e6180c2f9` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=20000 | j04n_cte_filtered_notmat | valid | `d43dd14675befada` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=100000 | j05_cte_joined_outer | valid | `6f0091d8475bff3c` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=180000 | j05_cte_joined_outer | valid | `116d58c44d68b0e2` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20 | j05_cte_joined_outer | valid | `9192a1cb0d423ae7` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=200 | j05_cte_joined_outer | valid | `da1656fbf706c609` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=2000 | j05_cte_joined_outer | valid | `9174168e6180c2f9` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20000 | j05_cte_joined_outer | valid | `d43dd14675befada` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_mat;threshold=100000 | j05m_cte_joined_outer_mat | valid | `177a6e26b8c3a785` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | j05m_cte_joined_outer_mat | valid | `177a6e26b8c3a785` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | j05m_cte_joined_outer_mat | valid | `177a6e26b8c3a785` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=200 | j05m_cte_joined_outer_mat | valid | `177a6e26b8c3a785` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | j05m_cte_joined_outer_mat | valid | `177a6e26b8c3a785` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20000 | j05m_cte_joined_outer_mat | valid | `177a6e26b8c3a785` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_notmat;threshold=100000 | j05n_cte_joined_outer_notmat | valid | `6f0091d8475bff3c` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | j05n_cte_joined_outer_notmat | valid | `116d58c44d68b0e2` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | j05n_cte_joined_outer_notmat | valid | `9192a1cb0d423ae7` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=200 | j05n_cte_joined_outer_notmat | valid | `da1656fbf706c609` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | j05n_cte_joined_outer_notmat | valid | `9174168e6180c2f9` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20000 | j05n_cte_joined_outer_notmat | valid | `d43dd14675befada` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=filtered_derived;threshold=100000 | j03_filtered_derived | valid | `6f0091d8475bff3c` |  |
| family=bag;stmt=filtered_derived;threshold=180000 | j03_filtered_derived | valid | `116d58c44d68b0e2` |  |
| family=bag;stmt=filtered_derived;threshold=20 | j03_filtered_derived | valid | `9192a1cb0d423ae7` |  |
| family=bag;stmt=filtered_derived;threshold=200 | j03_filtered_derived | valid | `da1656fbf706c609` |  |
| family=bag;stmt=filtered_derived;threshold=2000 | j03_filtered_derived | valid | `9174168e6180c2f9` |  |
| family=bag;stmt=filtered_derived;threshold=20000 | j03_filtered_derived | valid | `d43dd14675befada` |  |
| family=bag;stmt=join_on;threshold=100000 | j02_join_on | valid | `6f0091d8475bff3c` |  |
| family=bag;stmt=join_on;threshold=180000 | j02_join_on | valid | `116d58c44d68b0e2` |  |
| family=bag;stmt=join_on;threshold=20 | j02_join_on | valid | `9192a1cb0d423ae7` |  |
| family=bag;stmt=join_on;threshold=200 | j02_join_on | valid | `da1656fbf706c609` |  |
| family=bag;stmt=join_on;threshold=2000 | j02_join_on | valid | `9174168e6180c2f9` |  |
| family=bag;stmt=join_on;threshold=20000 | j02_join_on | valid | `d43dd14675befada` |  |
| family=bag;stmt=join_where;threshold=100000 | j01_join_where | valid | `6f0091d8475bff3c` |  |
| family=bag;stmt=join_where;threshold=100000;position=last | j01_join_where | valid | `6f0091d8475bff3c` |  |
| family=bag;stmt=join_where;threshold=180000 | j01_join_where | valid | `116d58c44d68b0e2` |  |
| family=bag;stmt=join_where;threshold=20 | j01_join_where | valid | `9192a1cb0d423ae7` |  |
| family=bag;stmt=join_where;threshold=200 | j01_join_where | valid | `da1656fbf706c609` |  |
| family=bag;stmt=join_where;threshold=2000 | j01_join_where | valid | `9174168e6180c2f9` |  |
| family=bag;stmt=join_where;threshold=20000 | j01_join_where | valid | `d43dd14675befada` |  |
| family=reuse;policy=default;refs=16;consumer=1% | r_default_n16_c1 | valid | `527af41dc9ea7c27` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=16;consumer=100% | r_default_n16_c100 | valid | `527af41dc9ea7c27` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=1;consumer=1% | r_default_n1_c1 | valid | `67feb88d971b83fe` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=1;consumer=100% | r_default_n1_c100 | valid | `b097208aea033185` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=2;consumer=1% | r_default_n2_c1 | valid | `450416f8dc65edb3` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=2;consumer=100% | r_default_n2_c100 | valid | `450416f8dc65edb3` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=1% | r_default_n4_c1 | valid | `91c3642bff5db3db` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=100% | r_default_n4_c100 | valid | `91c3642bff5db3db` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=8;consumer=1% | r_default_n8_c1 | valid | `0c9b689bfb53ff4f` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=8;consumer=100% | r_default_n8_c100 | valid | `0c9b689bfb53ff4f` | policy=default; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=1% | r_mat_n16_c1 | valid | `527af41dc9ea7c27` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=100% | r_mat_n16_c100 | valid | `527af41dc9ea7c27` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=1% | r_mat_n1_c1 | valid | `9384b48eeba05ee5` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=100% | r_mat_n1_c100 | valid | `9384b48eeba05ee5` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=2;consumer=1% | r_mat_n2_c1 | valid | `450416f8dc65edb3` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=2;consumer=100% | r_mat_n2_c100 | valid | `450416f8dc65edb3` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=1% | r_mat_n4_c1 | valid | `91c3642bff5db3db` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=100% | r_mat_n4_c100 | valid | `91c3642bff5db3db` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=8;consumer=1% | r_mat_n8_c1 | valid | `0c9b689bfb53ff4f` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=8;consumer=100% | r_mat_n8_c100 | valid | `0c9b689bfb53ff4f` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=notmat;refs=16;consumer=1% | r_notmat_n16_c1 | valid | `a586237c6b012f4c` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `1754eea8ba4cc2a2` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `67feb88d971b83fe` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `b097208aea033185` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=2;consumer=1% | r_notmat_n2_c1 | valid | `b6d08c76c5629cb2` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=2;consumer=100% | r_notmat_n2_c100 | valid | `cd29fa7a7d29b636` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `4e9a2c71c9021d34` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `3527c6c56c9bb190` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=8;consumer=1% | r_notmat_n8_c1 | valid | `b7419850b4ea67ad` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=8;consumer=100% | r_notmat_n8_c100 | valid | `eb0f827e9bbc0d1f` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, VERBOSE, BUFFERS, FORMAT TEXT` supported on postgres
- one host, single client: no independent hardware, real-network or equal-total-budget comparison
- no concurrent writers: the isolated read contract only; a matched blended read/write phase is a recorded gap
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
