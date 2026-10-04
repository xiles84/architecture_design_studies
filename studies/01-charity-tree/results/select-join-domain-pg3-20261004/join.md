# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 91 cells, 91 valid, engine `postgres`, topology `single`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-domain-pg3-20261004`
- run tag: `run/study-01/select-join-domain-pg3-20261004`
- environment: `host-zenbook-ux5406sa`
- repo commit: `ea77bc5c677aa9582a2646b597ea0211eac9d0e0`
- engine version: `PostgreSQL 17.11 (Debian 17.11-1.pgdg13+2) on x86_64-pc-linux-gnu, compiled by gcc (Debian 14.2.0-19) 14.2.0, 64-bit`
- database engine image: `docker.io/library/postgres:17.11` @ `7296f210ae81031ec955dbad9a67a84fe958572a2153b8d0826a647522904dc1`
- benchmark client image: `localhost/charitybench-select:select-join-domain-pg3-20261004` @ `fa3606dc9ba1e4169c0124c33434b7aa827570eb10aa634f479901ca10b2f0be`
- inputs digest: `8af025b10b490d7d3713d6b1a6e91dd092afa6480c83e59c42a1e1a2dfa22aba`
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
| 20 (0.01%) | 1.1177 | 1.0379 | 0.9249 | 1.0551 | 1.1285 | 1.0677 | 1.1123 | 91.0875 | 1.0277 | 1.0890 |
| 200 (0.1%) | 1.6067 | 1.7515 | 1.5025 | 1.6146 | 1.4828 | 1.4606 | 1.4381 | 71.3702 | 1.7451 | 1.4856 |
| 2000 (1%) | 3.5932 | 3.3856 | 3.2125 | 3.6275 | 4.9725 | 5.0778 | 3.3905 | 91.6398 | 5.4642 | 6.2662 |
| 20000 (10%) | 24.8175 | 13.9540 | 10.1589 | 22.0181 | 33.6297 | 20.0279 | 34.5062 | 120.0458 | 14.2443 | 18.0951 |
| 100000 (50%) | 111.3829 | 81.9913 | 112.1246 | 112.8479 | 119.3190 | 130.0051 | 135.8569 | 166.4652 | 82.8897 | 86.0306 |
| 180000 (90%) | 103.4024 | 139.1107 | 108.4891 | 137.4895 | 124.3221 | 156.8097 | 149.6988 | 196.5463 | 166.3514 | 174.9010 |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | join_on vs join_where | 1.1177 | 1.0379 | 0.9286 | true |
| 20 | filtered_derived vs join_where | 1.1177 | 0.9249 | 0.8275 | true |
| 20 | cte_filtered vs join_where | 1.1177 | 1.0551 | 0.9440 | true |
| 20 | cte_filtered_mat vs join_where | 1.1177 | 1.1285 | 1.0097 | false |
| 20 | cte_filtered_notmat vs join_where | 1.1177 | 1.0677 | 0.9553 | true |
| 20 | cte_joined_outer vs join_where | 1.1177 | 1.1123 | 0.9952 | true |
| 20 | cte_joined_outer_mat vs join_where | 1.1177 | 91.0875 | 81.4974 | false |
| 20 | cte_joined_outer_notmat vs join_where | 1.1177 | 1.0277 | 0.9195 | true |
| 20 | comma_join vs join_where | 1.1177 | 1.0890 | 0.9743 | true |
| 200 | join_on vs join_where | 1.6067 | 1.7515 | 1.0901 | true |
| 200 | filtered_derived vs join_where | 1.6067 | 1.5025 | 0.9351 | true |
| 200 | cte_filtered vs join_where | 1.6067 | 1.6146 | 1.0049 | true |
| 200 | cte_filtered_mat vs join_where | 1.6067 | 1.4828 | 0.9229 | false |
| 200 | cte_filtered_notmat vs join_where | 1.6067 | 1.4606 | 0.9091 | true |
| 200 | cte_joined_outer vs join_where | 1.6067 | 1.4381 | 0.8951 | true |
| 200 | cte_joined_outer_mat vs join_where | 1.6067 | 71.3702 | 44.4190 | false |
| 200 | cte_joined_outer_notmat vs join_where | 1.6067 | 1.7451 | 1.0861 | true |
| 200 | comma_join vs join_where | 1.6067 | 1.4856 | 0.9246 | true |
| 2000 | join_on vs join_where | 3.5932 | 3.3856 | 0.9422 | true |
| 2000 | filtered_derived vs join_where | 3.5932 | 3.2125 | 0.8940 | true |
| 2000 | cte_filtered vs join_where | 3.5932 | 3.6275 | 1.0095 | true |
| 2000 | cte_filtered_mat vs join_where | 3.5932 | 4.9725 | 1.3839 | false |
| 2000 | cte_filtered_notmat vs join_where | 3.5932 | 5.0778 | 1.4132 | true |
| 2000 | cte_joined_outer vs join_where | 3.5932 | 3.3905 | 0.9436 | true |
| 2000 | cte_joined_outer_mat vs join_where | 3.5932 | 91.6398 | 25.5034 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 3.5932 | 5.4642 | 1.5207 | true |
| 2000 | comma_join vs join_where | 3.5932 | 6.2662 | 1.7439 | true |
| 20000 | join_on vs join_where | 24.8175 | 13.9540 | 0.5623 | true |
| 20000 | filtered_derived vs join_where | 24.8175 | 10.1589 | 0.4093 | true |
| 20000 | cte_filtered vs join_where | 24.8175 | 22.0181 | 0.8872 | true |
| 20000 | cte_filtered_mat vs join_where | 24.8175 | 33.6297 | 1.3551 | false |
| 20000 | cte_filtered_notmat vs join_where | 24.8175 | 20.0279 | 0.8070 | true |
| 20000 | cte_joined_outer vs join_where | 24.8175 | 34.5062 | 1.3904 | true |
| 20000 | cte_joined_outer_mat vs join_where | 24.8175 | 120.0458 | 4.8371 | false |
| 20000 | cte_joined_outer_notmat vs join_where | 24.8175 | 14.2443 | 0.5740 | true |
| 20000 | comma_join vs join_where | 24.8175 | 18.0951 | 0.7291 | true |
| 100000 | join_on vs join_where | 111.3829 | 81.9913 | 0.7361 | true |
| 100000 | filtered_derived vs join_where | 111.3829 | 112.1246 | 1.0067 | true |
| 100000 | cte_filtered vs join_where | 111.3829 | 112.8479 | 1.0132 | true |
| 100000 | cte_filtered_mat vs join_where | 111.3829 | 119.3190 | 1.0713 | false |
| 100000 | cte_filtered_notmat vs join_where | 111.3829 | 130.0051 | 1.1672 | true |
| 100000 | cte_joined_outer vs join_where | 111.3829 | 135.8569 | 1.2197 | true |
| 100000 | cte_joined_outer_mat vs join_where | 111.3829 | 166.4652 | 1.4945 | false |
| 100000 | cte_joined_outer_notmat vs join_where | 111.3829 | 82.8897 | 0.7442 | true |
| 100000 | comma_join vs join_where | 111.3829 | 86.0306 | 0.7724 | true |
| 180000 | join_on vs join_where | 103.4024 | 139.1107 | 1.3453 | true |
| 180000 | filtered_derived vs join_where | 103.4024 | 108.4891 | 1.0492 | true |
| 180000 | cte_filtered vs join_where | 103.4024 | 137.4895 | 1.3297 | true |
| 180000 | cte_filtered_mat vs join_where | 103.4024 | 124.3221 | 1.2023 | false |
| 180000 | cte_filtered_notmat vs join_where | 103.4024 | 156.8097 | 1.5165 | true |
| 180000 | cte_joined_outer vs join_where | 103.4024 | 149.6988 | 1.4477 | true |
| 180000 | cte_joined_outer_mat vs join_where | 103.4024 | 196.5463 | 1.9008 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 103.4024 | 166.3514 | 1.6088 | true |
| 180000 | comma_join vs join_where | 103.4024 | 174.9010 | 1.6915 | true |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_filtered_mat vs cte_filtered | 1.0551 | 1.1285 | 1.0695 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_filtered_notmat vs cte_filtered | 1.0551 | 1.0677 | 1.0119 | true | policy=not_materialized; observed_cte_scan=false |
| 20 | cte_joined_outer_mat vs cte_joined_outer | 1.1123 | 91.0875 | 81.8944 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 1.1123 | 1.0277 | 0.9240 | true | policy=not_materialized; observed_cte_scan=false |
| 200 | cte_filtered_mat vs cte_filtered | 1.6146 | 1.4828 | 0.9184 | false | policy=materialized; observed_cte_scan=true |
| 200 | cte_filtered_notmat vs cte_filtered | 1.6146 | 1.4606 | 0.9046 | true | policy=not_materialized; observed_cte_scan=false |
| 200 | cte_joined_outer_mat vs cte_joined_outer | 1.4381 | 71.3702 | 49.6269 | false | policy=materialized; observed_cte_scan=true |
| 200 | cte_joined_outer_notmat vs cte_joined_outer | 1.4381 | 1.7451 | 1.2134 | true | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_filtered_mat vs cte_filtered | 3.6275 | 4.9725 | 1.3708 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_filtered_notmat vs cte_filtered | 3.6275 | 5.0778 | 1.3998 | true | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 3.3905 | 91.6398 | 27.0284 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 3.3905 | 5.4642 | 1.6116 | true | policy=not_materialized; observed_cte_scan=false |
| 20000 | cte_filtered_mat vs cte_filtered | 22.0181 | 33.6297 | 1.5274 | false | policy=materialized; observed_cte_scan=true |
| 20000 | cte_filtered_notmat vs cte_filtered | 22.0181 | 20.0279 | 0.9096 | true | policy=not_materialized; observed_cte_scan=false |
| 20000 | cte_joined_outer_mat vs cte_joined_outer | 34.5062 | 120.0458 | 3.4790 | false | policy=materialized; observed_cte_scan=true |
| 20000 | cte_joined_outer_notmat vs cte_joined_outer | 34.5062 | 14.2443 | 0.4128 | true | policy=not_materialized; observed_cte_scan=false |
| 100000 | cte_filtered_mat vs cte_filtered | 112.8479 | 119.3190 | 1.0573 | false | policy=materialized; observed_cte_scan=true |
| 100000 | cte_filtered_notmat vs cte_filtered | 112.8479 | 130.0051 | 1.1520 | true | policy=not_materialized; observed_cte_scan=false |
| 100000 | cte_joined_outer_mat vs cte_joined_outer | 135.8569 | 166.4652 | 1.2253 | false | policy=materialized; observed_cte_scan=true |
| 100000 | cte_joined_outer_notmat vs cte_joined_outer | 135.8569 | 82.8897 | 0.6101 | true | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_filtered_mat vs cte_filtered | 137.4895 | 124.3221 | 0.9042 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_filtered_notmat vs cte_filtered | 137.4895 | 156.8097 | 1.1405 | true | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 149.6988 | 196.5463 | 1.3129 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 149.6988 | 166.3514 | 1.1112 | true | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 1.1366 | 1 | `67feb88d971b83fe` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 4.4159 | 1 | `b097208aea033185` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | 4.5282 | 2 | `450416f8dc65edb3` | policy=default; observed_cte_scan=true | valid |
| default | 2 | 100% | 4.6779 | 2 | `450416f8dc65edb3` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 1% | 5.6812 | 4 | `91c3642bff5db3db` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 4.2442 | 4 | `91c3642bff5db3db` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | 6.1257 | 8 | `0c9b689bfb53ff4f` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 100% | 6.3234 | 8 | `0c9b689bfb53ff4f` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 1% | 6.0684 | 16 | `527af41dc9ea7c27` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 6.5124 | 16 | `527af41dc9ea7c27` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 5.0182 | 1 | `9384b48eeba05ee5` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 5.0806 | 1 | `9384b48eeba05ee5` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | 5.5009 | 2 | `450416f8dc65edb3` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 100% | 5.1618 | 2 | `450416f8dc65edb3` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 1% | 5.6445 | 4 | `91c3642bff5db3db` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 3.6059 | 4 | `91c3642bff5db3db` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | 5.7651 | 8 | `0c9b689bfb53ff4f` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 100% | 4.3625 | 8 | `0c9b689bfb53ff4f` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 1% | 6.8189 | 16 | `527af41dc9ea7c27` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 7.4622 | 16 | `527af41dc9ea7c27` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 1.0737 | 1 | `67feb88d971b83fe` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 4.5309 | 1 | `b097208aea033185` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | 1.4273 | 2 | `b6d08c76c5629cb2` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 100% | 15.5487 | 2 | `9c54ea41f4c79ade` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 1% | 1.9974 | 4 | `4e9a2c71c9021d34` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 19.4619 | 4 | `9b7035b7de60e1df` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | 2.9976 | 8 | `b7419850b4ea67ad` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 100% | 23.5775 | 8 | `f60391fbf46e4cb3` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 1% | 5.1315 | 16 | `a586237c6b012f4c` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 53.1498 | 16 | `b2256a0d4351b9bd` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 4.4152 | false |
| notmat | 1 | 1% | 0.9447 | true |
| mat | 1 | 100% | 1.1505 | false |
| notmat | 1 | 100% | 1.0260 | true |
| mat | 2 | 1% | 1.2148 | true |
| notmat | 2 | 1% | 0.3152 | false |
| mat | 2 | 100% | 1.1035 | true |
| notmat | 2 | 100% | 3.3239 | false |
| mat | 4 | 1% | 0.9935 | true |
| notmat | 4 | 1% | 0.3516 | false |
| mat | 4 | 100% | 0.8496 | true |
| notmat | 4 | 100% | 4.5856 | false |
| mat | 8 | 1% | 0.9411 | true |
| notmat | 8 | 1% | 0.4893 | false |
| mat | 8 | 100% | 0.6899 | true |
| notmat | 8 | 100% | 3.7286 | false |
| mat | 16 | 1% | 1.1237 | true |
| notmat | 16 | 1% | 0.8456 | false |
| mat | 16 | 100% | 1.1458 | true |
| notmat | 16 | 100% | 8.1613 | false |

## Repeated-design position control

| first median ms | last median ms | ratio (last/first) | same plan |
|---|---|---|---|
| 81.991315 | 72.533000 | 0.8846 | true |

## Latency samples and derived throughput

_The reader is sequential, so `queries/s` is 1000/median_ms, a derived rate and not a concurrent-capacity test. min/max are the five sample extremes._

| axis | status | min ms | median ms | max ms | samples | derived queries/s |
|---|---|---|---|---|---|---|
| family=bag;stmt=comma_join;threshold=100000 | valid | 72.8386 | 86.0306 | 121.3326 | 5 | 11.6 |
| family=bag;stmt=comma_join;threshold=180000 | valid | 125.4210 | 174.9010 | 225.6030 | 5 | 5.7 |
| family=bag;stmt=comma_join;threshold=20 | valid | 0.9505 | 1.0890 | 1.2222 | 5 | 918.3 |
| family=bag;stmt=comma_join;threshold=200 | valid | 1.1114 | 1.4856 | 1.8463 | 5 | 673.1 |
| family=bag;stmt=comma_join;threshold=2000 | valid | 3.2548 | 6.2662 | 11.3660 | 5 | 159.6 |
| family=bag;stmt=comma_join;threshold=20000 | valid | 13.1232 | 18.0951 | 42.2770 | 5 | 55.3 |
| family=bag;stmt=cte_filtered;threshold=100000 | valid | 100.0913 | 112.8479 | 147.1873 | 5 | 8.9 |
| family=bag;stmt=cte_filtered;threshold=180000 | valid | 78.1436 | 137.4895 | 161.2238 | 5 | 7.3 |
| family=bag;stmt=cte_filtered;threshold=20 | valid | 0.8944 | 1.0551 | 1.1772 | 5 | 947.7 |
| family=bag;stmt=cte_filtered;threshold=200 | valid | 0.9043 | 1.6146 | 2.1856 | 5 | 619.3 |
| family=bag;stmt=cte_filtered;threshold=2000 | valid | 2.5615 | 3.6275 | 4.1287 | 5 | 275.7 |
| family=bag;stmt=cte_filtered;threshold=20000 | valid | 18.1460 | 22.0181 | 55.1332 | 5 | 45.4 |
| family=bag;stmt=cte_filtered_mat;threshold=100000 | valid | 76.6899 | 119.3190 | 141.4544 | 5 | 8.4 |
| family=bag;stmt=cte_filtered_mat;threshold=180000 | valid | 75.0237 | 124.3221 | 141.4808 | 5 | 8.0 |
| family=bag;stmt=cte_filtered_mat;threshold=20 | valid | 0.7049 | 1.1285 | 2.1020 | 5 | 886.1 |
| family=bag;stmt=cte_filtered_mat;threshold=200 | valid | 0.7913 | 1.4828 | 1.8976 | 5 | 674.4 |
| family=bag;stmt=cte_filtered_mat;threshold=2000 | valid | 3.0651 | 4.9725 | 11.6827 | 5 | 201.1 |
| family=bag;stmt=cte_filtered_mat;threshold=20000 | valid | 24.8856 | 33.6297 | 82.1633 | 5 | 29.7 |
| family=bag;stmt=cte_filtered_notmat;threshold=100000 | valid | 66.0608 | 130.0051 | 181.4892 | 5 | 7.7 |
| family=bag;stmt=cte_filtered_notmat;threshold=180000 | valid | 136.0005 | 156.8097 | 186.8654 | 5 | 6.4 |
| family=bag;stmt=cte_filtered_notmat;threshold=20 | valid | 0.8964 | 1.0677 | 14.4127 | 5 | 936.6 |
| family=bag;stmt=cte_filtered_notmat;threshold=200 | valid | 1.0662 | 1.4606 | 1.7660 | 5 | 684.6 |
| family=bag;stmt=cte_filtered_notmat;threshold=2000 | valid | 1.9909 | 5.0778 | 23.7343 | 5 | 196.9 |
| family=bag;stmt=cte_filtered_notmat;threshold=20000 | valid | 15.9758 | 20.0279 | 25.6070 | 5 | 49.9 |
| family=bag;stmt=cte_joined_outer;threshold=100000 | valid | 121.6379 | 135.8569 | 168.1772 | 5 | 7.4 |
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 100.8929 | 149.6988 | 170.6871 | 5 | 6.7 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 0.7923 | 1.1123 | 1.2000 | 5 | 899.1 |
| family=bag;stmt=cte_joined_outer;threshold=200 | valid | 1.1466 | 1.4381 | 3.3705 | 5 | 695.3 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 2.2095 | 3.3905 | 4.8941 | 5 | 294.9 |
| family=bag;stmt=cte_joined_outer;threshold=20000 | valid | 13.6886 | 34.5062 | 39.2731 | 5 | 29.0 |
| family=bag;stmt=cte_joined_outer_mat;threshold=100000 | valid | 100.7760 | 166.4652 | 179.8689 | 5 | 6.0 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 158.2006 | 196.5463 | 214.4608 | 5 | 5.1 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 65.3721 | 91.0875 | 103.7676 | 5 | 11.0 |
| family=bag;stmt=cte_joined_outer_mat;threshold=200 | valid | 52.2862 | 71.3702 | 98.2347 | 5 | 14.0 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 66.0980 | 91.6398 | 109.2455 | 5 | 10.9 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20000 | valid | 102.1571 | 120.0458 | 125.6684 | 5 | 8.3 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=100000 | valid | 72.6043 | 82.8897 | 130.8703 | 5 | 12.1 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 129.4766 | 166.3514 | 180.3314 | 5 | 6.0 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 0.8841 | 1.0277 | 1.2267 | 5 | 973.1 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=200 | valid | 1.5411 | 1.7451 | 2.1492 | 5 | 573.0 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 3.0352 | 5.4642 | 15.4837 | 5 | 183.0 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20000 | valid | 9.5950 | 14.2443 | 26.7163 | 5 | 70.2 |
| family=bag;stmt=filtered_derived;threshold=100000 | valid | 83.7044 | 112.1246 | 150.9153 | 5 | 8.9 |
| family=bag;stmt=filtered_derived;threshold=180000 | valid | 65.1741 | 108.4891 | 168.3544 | 5 | 9.2 |
| family=bag;stmt=filtered_derived;threshold=20 | valid | 0.6792 | 0.9249 | 1.0394 | 5 | 1081.2 |
| family=bag;stmt=filtered_derived;threshold=200 | valid | 1.3992 | 1.5025 | 1.7366 | 5 | 665.6 |
| family=bag;stmt=filtered_derived;threshold=2000 | valid | 2.1550 | 3.2125 | 3.6892 | 5 | 311.3 |
| family=bag;stmt=filtered_derived;threshold=20000 | valid | 6.4353 | 10.1589 | 18.4923 | 5 | 98.4 |
| family=bag;stmt=join_on;threshold=100000 | valid | 66.6522 | 81.9913 | 152.9574 | 5 | 12.2 |
| family=bag;stmt=join_on;threshold=100000;position=last | valid | 64.3870 | 72.5330 | 127.0302 | 5 | 13.8 |
| family=bag;stmt=join_on;threshold=180000 | valid | 106.9674 | 139.1107 | 161.4690 | 5 | 7.2 |
| family=bag;stmt=join_on;threshold=20 | valid | 0.8694 | 1.0379 | 65.3720 | 5 | 963.5 |
| family=bag;stmt=join_on;threshold=200 | valid | 1.4770 | 1.7515 | 3.1512 | 5 | 570.9 |
| family=bag;stmt=join_on;threshold=2000 | valid | 3.1980 | 3.3856 | 15.6650 | 5 | 295.4 |
| family=bag;stmt=join_on;threshold=20000 | valid | 9.9301 | 13.9540 | 21.0154 | 5 | 71.7 |
| family=bag;stmt=join_where;threshold=100000 | valid | 80.2037 | 111.3829 | 172.8450 | 5 | 9.0 |
| family=bag;stmt=join_where;threshold=180000 | valid | 90.8215 | 103.4024 | 161.7417 | 5 | 9.7 |
| family=bag;stmt=join_where;threshold=20 | valid | 1.0006 | 1.1177 | 1.7746 | 5 | 894.7 |
| family=bag;stmt=join_where;threshold=200 | valid | 1.5267 | 1.6067 | 1.8533 | 5 | 622.4 |
| family=bag;stmt=join_where;threshold=2000 | valid | 3.4068 | 3.5932 | 4.1907 | 5 | 278.3 |
| family=bag;stmt=join_where;threshold=20000 | valid | 14.7167 | 24.8175 | 43.2478 | 5 | 40.3 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 5.1136 | 6.0684 | 6.5096 | 5 | 164.8 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 5.2690 | 6.5124 | 7.3996 | 5 | 153.6 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 0.8967 | 1.1366 | 1.1378 | 5 | 879.8 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 3.0157 | 4.4159 | 5.0492 | 5 | 226.5 |
| family=reuse;policy=default;refs=2;consumer=1% | valid | 3.3724 | 4.5282 | 7.5642 | 5 | 220.8 |
| family=reuse;policy=default;refs=2;consumer=100% | valid | 3.2443 | 4.6779 | 9.6032 | 5 | 213.8 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 4.6891 | 5.6812 | 7.9604 | 5 | 176.0 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 2.8848 | 4.2442 | 5.9120 | 5 | 235.6 |
| family=reuse;policy=default;refs=8;consumer=1% | valid | 5.3299 | 6.1257 | 8.3603 | 5 | 163.2 |
| family=reuse;policy=default;refs=8;consumer=100% | valid | 5.8687 | 6.3234 | 6.9055 | 5 | 158.1 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 6.4857 | 6.8189 | 9.7687 | 5 | 146.7 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 4.7518 | 7.4622 | 8.9676 | 5 | 134.0 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 4.2561 | 5.0182 | 5.4900 | 5 | 199.3 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 3.0959 | 5.0806 | 5.9027 | 5 | 196.8 |
| family=reuse;policy=mat;refs=2;consumer=1% | valid | 3.4381 | 5.5009 | 7.7786 | 5 | 181.8 |
| family=reuse;policy=mat;refs=2;consumer=100% | valid | 4.6536 | 5.1618 | 6.1096 | 5 | 193.7 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 4.6739 | 5.6445 | 8.6535 | 5 | 177.2 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 2.9165 | 3.6059 | 6.1080 | 5 | 277.3 |
| family=reuse;policy=mat;refs=8;consumer=1% | valid | 5.4567 | 5.7651 | 11.6001 | 5 | 173.5 |
| family=reuse;policy=mat;refs=8;consumer=100% | valid | 4.1933 | 4.3625 | 6.4156 | 5 | 229.2 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 3.6402 | 5.1315 | 6.0098 | 5 | 194.9 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 35.1342 | 53.1498 | 62.9689 | 5 | 18.8 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 0.8594 | 1.0737 | 1.1445 | 5 | 931.3 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 4.2040 | 4.5309 | 5.3899 | 5 | 220.7 |
| family=reuse;policy=notmat;refs=2;consumer=1% | valid | 1.0283 | 1.4273 | 2.1683 | 5 | 700.6 |
| family=reuse;policy=notmat;refs=2;consumer=100% | valid | 9.2804 | 15.5487 | 21.7590 | 5 | 64.3 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 1.1379 | 1.9974 | 3.6578 | 5 | 500.7 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 14.3644 | 19.4619 | 30.7318 | 5 | 51.4 |
| family=reuse;policy=notmat;refs=8;consumer=1% | valid | 2.2408 | 2.9976 | 5.1545 | 5 | 333.6 |
| family=reuse;policy=notmat;refs=8;consumer=100% | valid | 21.0091 | 23.5775 | 24.9520 | 5 | 42.4 |

## Execution order

- order seed `42860634`, 5 interleaved rounds
  1. family=bag;stmt=join_on;threshold=100000
  2. family=bag;stmt=comma_join;threshold=200
  3. family=bag;stmt=cte_filtered_notmat;threshold=2000
  4. family=reuse;policy=mat;refs=8;consumer=100%
  5. family=bag;stmt=cte_joined_outer_notmat;threshold=100000
  6. family=bag;stmt=join_where;threshold=20
  7. family=reuse;policy=notmat;refs=4;consumer=1%
  8. family=bag;stmt=cte_filtered;threshold=200
  9. family=bag;stmt=cte_joined_outer;threshold=180000
  10. family=reuse;policy=notmat;refs=4;consumer=100%
  11. family=reuse;policy=mat;refs=4;consumer=1%
  12. family=bag;stmt=comma_join;threshold=2000
  13. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  14. family=reuse;policy=notmat;refs=1;consumer=1%
  15. family=bag;stmt=cte_joined_outer_notmat;threshold=200
  16. family=reuse;policy=notmat;refs=16;consumer=100%
  17. family=bag;stmt=join_on;threshold=180000
  18. family=reuse;policy=mat;refs=16;consumer=1%
  19. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  20. family=bag;stmt=filtered_derived;threshold=200
  21. family=reuse;policy=default;refs=4;consumer=1%
  22. family=bag;stmt=cte_filtered_notmat;threshold=20
  23. family=reuse;policy=mat;refs=8;consumer=1%
  24. family=reuse;policy=mat;refs=2;consumer=100%
  25. family=bag;stmt=cte_filtered;threshold=100000
  26. family=bag;stmt=join_where;threshold=2000
  27. family=bag;stmt=cte_filtered;threshold=20
  28. family=bag;stmt=join_where;threshold=20000
  29. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  30. family=bag;stmt=filtered_derived;threshold=100000
  31. family=reuse;policy=mat;refs=1;consumer=1%
  32. family=bag;stmt=cte_filtered_mat;threshold=2000
  33. family=bag;stmt=cte_joined_outer;threshold=100000
  34. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  35. family=bag;stmt=join_on;threshold=200
  36. family=reuse;policy=default;refs=8;consumer=1%
  37. family=bag;stmt=cte_filtered;threshold=20000
  38. family=reuse;policy=default;refs=8;consumer=100%
  39. family=bag;stmt=cte_filtered_notmat;threshold=180000
  40. family=bag;stmt=cte_filtered_mat;threshold=20000
  41. family=bag;stmt=comma_join;threshold=20
  42. family=bag;stmt=join_on;threshold=20
  43. family=bag;stmt=join_where;threshold=200
  44. family=bag;stmt=comma_join;threshold=180000
  45. family=reuse;policy=mat;refs=16;consumer=100%
  46. family=bag;stmt=cte_filtered_mat;threshold=20
  47. family=bag;stmt=cte_joined_outer;threshold=2000
  48. family=bag;stmt=cte_joined_outer;threshold=200
  49. family=reuse;policy=mat;refs=1;consumer=100%
  50. family=bag;stmt=cte_filtered_notmat;threshold=200
  51. family=bag;stmt=cte_joined_outer_mat;threshold=100000
  52. family=bag;stmt=cte_filtered_mat;threshold=100000
  53. family=reuse;policy=default;refs=4;consumer=100%
  54. family=bag;stmt=filtered_derived;threshold=20000
  55. family=bag;stmt=cte_joined_outer_mat;threshold=200
  56. family=reuse;policy=mat;refs=4;consumer=100%
  57. family=bag;stmt=cte_filtered_mat;threshold=200
  58. family=reuse;policy=notmat;refs=16;consumer=1%
  59. family=bag;stmt=filtered_derived;threshold=180000
  60. family=reuse;policy=default;refs=1;consumer=100%
  61. family=reuse;policy=default;refs=2;consumer=100%
  62. family=bag;stmt=join_where;threshold=180000
  63. family=bag;stmt=comma_join;threshold=100000
  64. family=reuse;policy=default;refs=16;consumer=1%
  65. family=bag;stmt=filtered_derived;threshold=20
  66. family=reuse;policy=notmat;refs=8;consumer=100%
  67. family=bag;stmt=cte_filtered;threshold=180000
  68. family=bag;stmt=cte_filtered;threshold=2000
  69. family=bag;stmt=join_on;threshold=2000
  70. family=bag;stmt=cte_filtered_notmat;threshold=100000
  71. family=bag;stmt=cte_joined_outer;threshold=20
  72. family=bag;stmt=cte_filtered_notmat;threshold=20000
  73. family=bag;stmt=comma_join;threshold=20000
  74. family=reuse;policy=notmat;refs=2;consumer=1%
  75. family=bag;stmt=join_where;threshold=100000
  76. family=bag;stmt=cte_joined_outer;threshold=20000
  77. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  78. family=reuse;policy=notmat;refs=1;consumer=100%
  79. family=reuse;policy=mat;refs=2;consumer=1%
  80. family=bag;stmt=filtered_derived;threshold=2000
  81. family=reuse;policy=notmat;refs=2;consumer=100%
  82. family=bag;stmt=cte_filtered_mat;threshold=180000
  83. family=bag;stmt=cte_joined_outer_mat;threshold=20000
  84. family=reuse;policy=default;refs=1;consumer=1%
  85. family=bag;stmt=cte_joined_outer_mat;threshold=20
  86. family=reuse;policy=notmat;refs=8;consumer=1%
  87. family=reuse;policy=default;refs=2;consumer=1%
  88. family=bag;stmt=join_on;threshold=20000
  89. family=bag;stmt=cte_joined_outer_notmat;threshold=20000
  90. family=reuse;policy=default;refs=16;consumer=100%
  91. family=bag;stmt=join_on;threshold=100000;position=last

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
| family=bag;stmt=join_on;threshold=100000;position=last | j02_join_on | valid | `6f0091d8475bff3c` |  |
| family=bag;stmt=join_on;threshold=180000 | j02_join_on | valid | `116d58c44d68b0e2` |  |
| family=bag;stmt=join_on;threshold=20 | j02_join_on | valid | `9192a1cb0d423ae7` |  |
| family=bag;stmt=join_on;threshold=200 | j02_join_on | valid | `da1656fbf706c609` |  |
| family=bag;stmt=join_on;threshold=2000 | j02_join_on | valid | `9174168e6180c2f9` |  |
| family=bag;stmt=join_on;threshold=20000 | j02_join_on | valid | `d43dd14675befada` |  |
| family=bag;stmt=join_where;threshold=100000 | j01_join_where | valid | `6f0091d8475bff3c` |  |
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
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `b2256a0d4351b9bd` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `67feb88d971b83fe` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `b097208aea033185` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=2;consumer=1% | r_notmat_n2_c1 | valid | `b6d08c76c5629cb2` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=2;consumer=100% | r_notmat_n2_c100 | valid | `9c54ea41f4c79ade` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `4e9a2c71c9021d34` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `9b7035b7de60e1df` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=8;consumer=1% | r_notmat_n8_c1 | valid | `b7419850b4ea67ad` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=8;consumer=100% | r_notmat_n8_c100 | valid | `f60391fbf46e4cb3` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, VERBOSE, BUFFERS, FORMAT TEXT` supported on postgres
- one host, single client: no independent hardware, real-network or equal-total-budget comparison
- no concurrent writers: the isolated read contract only; a matched blended read/write phase is a recorded gap
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
