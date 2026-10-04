# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 91 cells, 91 valid, engine `postgres`, topology `single`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-domain-pg1-20261004`
- run tag: `run/study-01/select-join-domain-pg1-20261004`
- environment: `host-zenbook-ux5406sa`
- repo commit: `57cd161d1ecbcc90f7ca798334e98a578046df36`
- engine version: `PostgreSQL 17.11 (Debian 17.11-1.pgdg13+2) on x86_64-pc-linux-gnu, compiled by gcc (Debian 14.2.0-19) 14.2.0, 64-bit`
- database engine image: `docker.io/library/postgres:17.11` @ `7296f210ae81031ec955dbad9a67a84fe958572a2153b8d0826a647522904dc1`
- benchmark client image: `localhost/charitybench-select:select-join-domain-pg1-20261004` @ `fa3606dc9ba1e4169c0124c33434b7aa827570eb10aa634f479901ca10b2f0be`
- inputs digest: `e285a32421c22f43f4b0dab56e20a32c7fed90c0a445dd4cecbf66de6b7cc6f3`
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
| 20 (0.01%) | 1.0654 | 0.8795 | 1.0375 | 0.9902 | 0.9213 | 0.9278 | 0.9676 | 91.2885 | 0.8834 | 1.0066 |
| 200 (0.1%) | 1.4746 | 1.5769 | 1.6059 | 1.5094 | 1.3681 | 1.7450 | 1.5839 | 90.6451 | 1.5029 | 1.5537 |
| 2000 (1%) | 3.2721 | 3.4136 | 3.5197 | 4.1071 | 3.6700 | 3.4413 | 3.4964 | 95.2411 | 3.4607 | 3.3154 |
| 20000 (10%) | 19.0600 | 15.4189 | 15.8680 | 16.9481 | 18.9753 | 16.6135 | 22.6977 | 117.0927 | 22.3343 | 19.7013 |
| 100000 (50%) | 122.3986 | 104.3220 | 120.0549 | 111.2374 | 125.4378 | 132.2922 | 127.4178 | 147.3646 | 106.3374 | 127.8640 |
| 180000 (90%) | 158.4710 | 175.9199 | 168.1826 | 125.9838 | 142.6691 | 162.8678 | 157.1349 | 212.6815 | 150.0152 | 119.9697 |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | join_on vs join_where | 1.0654 | 0.8795 | 0.8255 | true |
| 20 | filtered_derived vs join_where | 1.0654 | 1.0375 | 0.9738 | true |
| 20 | cte_filtered vs join_where | 1.0654 | 0.9902 | 0.9294 | true |
| 20 | cte_filtered_mat vs join_where | 1.0654 | 0.9213 | 0.8647 | false |
| 20 | cte_filtered_notmat vs join_where | 1.0654 | 0.9278 | 0.8709 | true |
| 20 | cte_joined_outer vs join_where | 1.0654 | 0.9676 | 0.9083 | true |
| 20 | cte_joined_outer_mat vs join_where | 1.0654 | 91.2885 | 85.6869 | false |
| 20 | cte_joined_outer_notmat vs join_where | 1.0654 | 0.8834 | 0.8292 | true |
| 20 | comma_join vs join_where | 1.0654 | 1.0066 | 0.9448 | true |
| 200 | join_on vs join_where | 1.4746 | 1.5769 | 1.0694 | true |
| 200 | filtered_derived vs join_where | 1.4746 | 1.6059 | 1.0891 | true |
| 200 | cte_filtered vs join_where | 1.4746 | 1.5094 | 1.0236 | true |
| 200 | cte_filtered_mat vs join_where | 1.4746 | 1.3681 | 0.9278 | false |
| 200 | cte_filtered_notmat vs join_where | 1.4746 | 1.7450 | 1.1834 | true |
| 200 | cte_joined_outer vs join_where | 1.4746 | 1.5839 | 1.0741 | true |
| 200 | cte_joined_outer_mat vs join_where | 1.4746 | 90.6451 | 61.4701 | false |
| 200 | cte_joined_outer_notmat vs join_where | 1.4746 | 1.5029 | 1.0192 | true |
| 200 | comma_join vs join_where | 1.4746 | 1.5537 | 1.0536 | true |
| 2000 | join_on vs join_where | 3.2721 | 3.4136 | 1.0432 | true |
| 2000 | filtered_derived vs join_where | 3.2721 | 3.5197 | 1.0757 | true |
| 2000 | cte_filtered vs join_where | 3.2721 | 4.1071 | 1.2552 | true |
| 2000 | cte_filtered_mat vs join_where | 3.2721 | 3.6700 | 1.1216 | false |
| 2000 | cte_filtered_notmat vs join_where | 3.2721 | 3.4413 | 1.0517 | true |
| 2000 | cte_joined_outer vs join_where | 3.2721 | 3.4964 | 1.0686 | true |
| 2000 | cte_joined_outer_mat vs join_where | 3.2721 | 95.2411 | 29.1067 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 3.2721 | 3.4607 | 1.0576 | true |
| 2000 | comma_join vs join_where | 3.2721 | 3.3154 | 1.0132 | true |
| 20000 | join_on vs join_where | 19.0600 | 15.4189 | 0.8090 | true |
| 20000 | filtered_derived vs join_where | 19.0600 | 15.8680 | 0.8325 | true |
| 20000 | cte_filtered vs join_where | 19.0600 | 16.9481 | 0.8892 | true |
| 20000 | cte_filtered_mat vs join_where | 19.0600 | 18.9753 | 0.9956 | false |
| 20000 | cte_filtered_notmat vs join_where | 19.0600 | 16.6135 | 0.8716 | true |
| 20000 | cte_joined_outer vs join_where | 19.0600 | 22.6977 | 1.1909 | true |
| 20000 | cte_joined_outer_mat vs join_where | 19.0600 | 117.0927 | 6.1434 | false |
| 20000 | cte_joined_outer_notmat vs join_where | 19.0600 | 22.3343 | 1.1718 | true |
| 20000 | comma_join vs join_where | 19.0600 | 19.7013 | 1.0336 | true |
| 100000 | join_on vs join_where | 122.3986 | 104.3220 | 0.8523 | true |
| 100000 | filtered_derived vs join_where | 122.3986 | 120.0549 | 0.9809 | true |
| 100000 | cte_filtered vs join_where | 122.3986 | 111.2374 | 0.9088 | true |
| 100000 | cte_filtered_mat vs join_where | 122.3986 | 125.4378 | 1.0248 | false |
| 100000 | cte_filtered_notmat vs join_where | 122.3986 | 132.2922 | 1.0808 | true |
| 100000 | cte_joined_outer vs join_where | 122.3986 | 127.4178 | 1.0410 | true |
| 100000 | cte_joined_outer_mat vs join_where | 122.3986 | 147.3646 | 1.2040 | false |
| 100000 | cte_joined_outer_notmat vs join_where | 122.3986 | 106.3374 | 0.8688 | true |
| 100000 | comma_join vs join_where | 122.3986 | 127.8640 | 1.0447 | true |
| 180000 | join_on vs join_where | 158.4710 | 175.9199 | 1.1101 | true |
| 180000 | filtered_derived vs join_where | 158.4710 | 168.1826 | 1.0613 | true |
| 180000 | cte_filtered vs join_where | 158.4710 | 125.9838 | 0.7950 | true |
| 180000 | cte_filtered_mat vs join_where | 158.4710 | 142.6691 | 0.9003 | false |
| 180000 | cte_filtered_notmat vs join_where | 158.4710 | 162.8678 | 1.0277 | true |
| 180000 | cte_joined_outer vs join_where | 158.4710 | 157.1349 | 0.9916 | true |
| 180000 | cte_joined_outer_mat vs join_where | 158.4710 | 212.6815 | 1.3421 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 158.4710 | 150.0152 | 0.9466 | true |
| 180000 | comma_join vs join_where | 158.4710 | 119.9697 | 0.7570 | true |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_filtered_mat vs cte_filtered | 0.9902 | 0.9213 | 0.9304 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_filtered_notmat vs cte_filtered | 0.9902 | 0.9278 | 0.9370 | true | policy=not_materialized; observed_cte_scan=false |
| 20 | cte_joined_outer_mat vs cte_joined_outer | 0.9676 | 91.2885 | 94.3407 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 0.9676 | 0.8834 | 0.9129 | true | policy=not_materialized; observed_cte_scan=false |
| 200 | cte_filtered_mat vs cte_filtered | 1.5094 | 1.3681 | 0.9064 | false | policy=materialized; observed_cte_scan=true |
| 200 | cte_filtered_notmat vs cte_filtered | 1.5094 | 1.7450 | 1.1561 | true | policy=not_materialized; observed_cte_scan=false |
| 200 | cte_joined_outer_mat vs cte_joined_outer | 1.5839 | 90.6451 | 57.2293 | false | policy=materialized; observed_cte_scan=true |
| 200 | cte_joined_outer_notmat vs cte_joined_outer | 1.5839 | 1.5029 | 0.9489 | true | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_filtered_mat vs cte_filtered | 4.1071 | 3.6700 | 0.8936 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_filtered_notmat vs cte_filtered | 4.1071 | 3.4413 | 0.8379 | true | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 3.4964 | 95.2411 | 27.2394 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 3.4964 | 3.4607 | 0.9898 | true | policy=not_materialized; observed_cte_scan=false |
| 20000 | cte_filtered_mat vs cte_filtered | 16.9481 | 18.9753 | 1.1196 | false | policy=materialized; observed_cte_scan=true |
| 20000 | cte_filtered_notmat vs cte_filtered | 16.9481 | 16.6135 | 0.9803 | true | policy=not_materialized; observed_cte_scan=false |
| 20000 | cte_joined_outer_mat vs cte_joined_outer | 22.6977 | 117.0927 | 5.1588 | false | policy=materialized; observed_cte_scan=true |
| 20000 | cte_joined_outer_notmat vs cte_joined_outer | 22.6977 | 22.3343 | 0.9840 | true | policy=not_materialized; observed_cte_scan=false |
| 100000 | cte_filtered_mat vs cte_filtered | 111.2374 | 125.4378 | 1.1277 | false | policy=materialized; observed_cte_scan=true |
| 100000 | cte_filtered_notmat vs cte_filtered | 111.2374 | 132.2922 | 1.1893 | true | policy=not_materialized; observed_cte_scan=false |
| 100000 | cte_joined_outer_mat vs cte_joined_outer | 127.4178 | 147.3646 | 1.1565 | false | policy=materialized; observed_cte_scan=true |
| 100000 | cte_joined_outer_notmat vs cte_joined_outer | 127.4178 | 106.3374 | 0.8346 | true | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_filtered_mat vs cte_filtered | 125.9838 | 142.6691 | 1.1324 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_filtered_notmat vs cte_filtered | 125.9838 | 162.8678 | 1.2928 | true | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 157.1349 | 212.6815 | 1.3535 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 157.1349 | 150.0152 | 0.9547 | true | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 1.0877 | 1 | `67feb88d971b83fe` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 4.7786 | 1 | `b097208aea033185` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | 4.8223 | 2 | `450416f8dc65edb3` | policy=default; observed_cte_scan=true | valid |
| default | 2 | 100% | 6.1419 | 2 | `450416f8dc65edb3` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 1% | 4.9065 | 4 | `91c3642bff5db3db` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 5.2060 | 4 | `91c3642bff5db3db` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | 5.9973 | 8 | `0c9b689bfb53ff4f` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 100% | 6.1096 | 8 | `0c9b689bfb53ff4f` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 1% | 6.5675 | 16 | `527af41dc9ea7c27` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 7.4641 | 16 | `527af41dc9ea7c27` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 4.8967 | 1 | `9384b48eeba05ee5` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 5.1171 | 1 | `9384b48eeba05ee5` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | 4.6856 | 2 | `450416f8dc65edb3` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 100% | 4.8363 | 2 | `450416f8dc65edb3` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 1% | 5.0619 | 4 | `91c3642bff5db3db` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 5.5407 | 4 | `91c3642bff5db3db` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | 5.3902 | 8 | `0c9b689bfb53ff4f` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 100% | 5.6833 | 8 | `0c9b689bfb53ff4f` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 1% | 6.4718 | 16 | `527af41dc9ea7c27` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 7.9919 | 16 | `527af41dc9ea7c27` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 1.1817 | 1 | `67feb88d971b83fe` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 4.4720 | 1 | `b097208aea033185` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | 1.6276 | 2 | `b6d08c76c5629cb2` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 100% | 11.5251 | 2 | `b5b1a9f8ee6f1d86` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 1% | 1.7496 | 4 | `4e9a2c71c9021d34` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 23.1480 | 4 | `b2dc573d0452c7ec` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | 2.9625 | 8 | `b7419850b4ea67ad` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 100% | 20.4004 | 8 | `f9a803ecffcb61fe` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 1% | 5.6004 | 16 | `a586237c6b012f4c` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 33.2215 | 16 | `daea55c209e40a05` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 4.5019 | false |
| notmat | 1 | 1% | 1.0865 | true |
| mat | 1 | 100% | 1.0708 | false |
| notmat | 1 | 100% | 0.9359 | true |
| mat | 2 | 1% | 0.9717 | true |
| notmat | 2 | 1% | 0.3375 | false |
| mat | 2 | 100% | 0.7874 | true |
| notmat | 2 | 100% | 1.8765 | false |
| mat | 4 | 1% | 1.0317 | true |
| notmat | 4 | 1% | 0.3566 | false |
| mat | 4 | 100% | 1.0643 | true |
| notmat | 4 | 100% | 4.4464 | false |
| mat | 8 | 1% | 0.8988 | true |
| notmat | 8 | 1% | 0.4940 | false |
| mat | 8 | 100% | 0.9302 | true |
| notmat | 8 | 100% | 3.3391 | false |
| mat | 16 | 1% | 0.9854 | true |
| notmat | 16 | 1% | 0.8527 | false |
| mat | 16 | 100% | 1.0707 | true |
| notmat | 16 | 100% | 4.4508 | false |

## Repeated-design position control

| first median ms | last median ms | ratio (last/first) | same plan |
|---|---|---|---|
| 4.836267 | 4.981896 | 1.0301 | true |

## Latency samples and derived throughput

_The reader is sequential, so `queries/s` is 1000/median_ms, a derived rate and not a concurrent-capacity test. min/max are the five sample extremes._

| axis | status | min ms | median ms | max ms | samples | derived queries/s |
|---|---|---|---|---|---|---|
| family=bag;stmt=comma_join;threshold=100000 | valid | 123.3866 | 127.8640 | 135.9354 | 5 | 7.8 |
| family=bag;stmt=comma_join;threshold=180000 | valid | 94.9313 | 119.9697 | 129.0196 | 5 | 8.3 |
| family=bag;stmt=comma_join;threshold=20 | valid | 0.8685 | 1.0066 | 1.0743 | 5 | 993.5 |
| family=bag;stmt=comma_join;threshold=200 | valid | 1.3657 | 1.5537 | 1.6795 | 5 | 643.6 |
| family=bag;stmt=comma_join;threshold=2000 | valid | 2.9989 | 3.3154 | 4.0821 | 5 | 301.6 |
| family=bag;stmt=comma_join;threshold=20000 | valid | 13.0617 | 19.7013 | 36.3771 | 5 | 50.8 |
| family=bag;stmt=cte_filtered;threshold=100000 | valid | 102.8133 | 111.2374 | 129.2585 | 5 | 9.0 |
| family=bag;stmt=cte_filtered;threshold=180000 | valid | 109.2999 | 125.9838 | 140.5868 | 5 | 7.9 |
| family=bag;stmt=cte_filtered;threshold=20 | valid | 0.8768 | 0.9902 | 1.2546 | 5 | 1009.9 |
| family=bag;stmt=cte_filtered;threshold=200 | valid | 1.4012 | 1.5094 | 1.5706 | 5 | 662.5 |
| family=bag;stmt=cte_filtered;threshold=2000 | valid | 3.3534 | 4.1071 | 6.6777 | 5 | 243.5 |
| family=bag;stmt=cte_filtered;threshold=20000 | valid | 13.2636 | 16.9481 | 20.8519 | 5 | 59.0 |
| family=bag;stmt=cte_filtered_mat;threshold=100000 | valid | 97.8624 | 125.4378 | 127.4602 | 5 | 8.0 |
| family=bag;stmt=cte_filtered_mat;threshold=180000 | valid | 119.3648 | 142.6691 | 167.0740 | 5 | 7.0 |
| family=bag;stmt=cte_filtered_mat;threshold=20 | valid | 0.8432 | 0.9213 | 2.8981 | 5 | 1085.4 |
| family=bag;stmt=cte_filtered_mat;threshold=200 | valid | 1.2818 | 1.3681 | 1.8759 | 5 | 730.9 |
| family=bag;stmt=cte_filtered_mat;threshold=2000 | valid | 3.3444 | 3.6700 | 13.6529 | 5 | 272.5 |
| family=bag;stmt=cte_filtered_mat;threshold=20000 | valid | 15.6305 | 18.9753 | 33.6764 | 5 | 52.7 |
| family=bag;stmt=cte_filtered_notmat;threshold=100000 | valid | 116.5500 | 132.2922 | 156.9998 | 5 | 7.6 |
| family=bag;stmt=cte_filtered_notmat;threshold=180000 | valid | 130.5682 | 162.8678 | 185.8836 | 5 | 6.1 |
| family=bag;stmt=cte_filtered_notmat;threshold=20 | valid | 0.8829 | 0.9278 | 1.3352 | 5 | 1077.8 |
| family=bag;stmt=cte_filtered_notmat;threshold=200 | valid | 1.5014 | 1.7450 | 2.4589 | 5 | 573.1 |
| family=bag;stmt=cte_filtered_notmat;threshold=2000 | valid | 3.1649 | 3.4413 | 10.6356 | 5 | 290.6 |
| family=bag;stmt=cte_filtered_notmat;threshold=20000 | valid | 16.2372 | 16.6135 | 18.8133 | 5 | 60.2 |
| family=bag;stmt=cte_joined_outer;threshold=100000 | valid | 105.5650 | 127.4178 | 144.8843 | 5 | 7.8 |
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 126.4980 | 157.1349 | 186.7965 | 5 | 6.4 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 0.8534 | 0.9676 | 2.4759 | 5 | 1033.4 |
| family=bag;stmt=cte_joined_outer;threshold=200 | valid | 1.4962 | 1.5839 | 1.7770 | 5 | 631.4 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 3.0670 | 3.4964 | 6.6330 | 5 | 286.0 |
| family=bag;stmt=cte_joined_outer;threshold=20000 | valid | 16.5092 | 22.6977 | 36.8247 | 5 | 44.1 |
| family=bag;stmt=cte_joined_outer_mat;threshold=100000 | valid | 114.3013 | 147.3646 | 156.5618 | 5 | 6.8 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 198.8655 | 212.6815 | 232.4632 | 5 | 4.7 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 84.1251 | 91.2885 | 95.5564 | 5 | 11.0 |
| family=bag;stmt=cte_joined_outer_mat;threshold=200 | valid | 87.8607 | 90.6451 | 104.2553 | 5 | 11.0 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 86.1941 | 95.2411 | 103.9314 | 5 | 10.5 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20000 | valid | 71.3570 | 117.0927 | 125.2656 | 5 | 8.5 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=100000 | valid | 105.4699 | 106.3374 | 132.1562 | 5 | 9.4 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 103.8149 | 150.0152 | 167.1704 | 5 | 6.7 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 0.8408 | 0.8834 | 1.9418 | 5 | 1132.0 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=200 | valid | 1.4621 | 1.5029 | 17.6872 | 5 | 665.4 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 3.3277 | 3.4607 | 6.8010 | 5 | 289.0 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20000 | valid | 17.2359 | 22.3343 | 34.0151 | 5 | 44.8 |
| family=bag;stmt=filtered_derived;threshold=100000 | valid | 103.6891 | 120.0549 | 130.8479 | 5 | 8.3 |
| family=bag;stmt=filtered_derived;threshold=180000 | valid | 121.3185 | 168.1826 | 182.5573 | 5 | 5.9 |
| family=bag;stmt=filtered_derived;threshold=20 | valid | 0.8472 | 1.0375 | 28.3193 | 5 | 963.9 |
| family=bag;stmt=filtered_derived;threshold=200 | valid | 1.5762 | 1.6059 | 1.9400 | 5 | 622.7 |
| family=bag;stmt=filtered_derived;threshold=2000 | valid | 2.9388 | 3.5197 | 6.5199 | 5 | 284.1 |
| family=bag;stmt=filtered_derived;threshold=20000 | valid | 15.8218 | 15.8680 | 24.0943 | 5 | 63.0 |
| family=bag;stmt=join_on;threshold=100000 | valid | 53.7134 | 104.3220 | 138.1172 | 5 | 9.6 |
| family=bag;stmt=join_on;threshold=180000 | valid | 167.2067 | 175.9199 | 207.7571 | 5 | 5.7 |
| family=bag;stmt=join_on;threshold=20 | valid | 0.8610 | 0.8795 | 1.2039 | 5 | 1137.1 |
| family=bag;stmt=join_on;threshold=200 | valid | 1.4709 | 1.5769 | 1.9477 | 5 | 634.1 |
| family=bag;stmt=join_on;threshold=2000 | valid | 3.0575 | 3.4136 | 4.4084 | 5 | 292.9 |
| family=bag;stmt=join_on;threshold=20000 | valid | 13.8406 | 15.4189 | 22.3865 | 5 | 64.9 |
| family=bag;stmt=join_where;threshold=100000 | valid | 101.2471 | 122.3986 | 160.3515 | 5 | 8.2 |
| family=bag;stmt=join_where;threshold=180000 | valid | 117.9888 | 158.4710 | 166.3505 | 5 | 6.3 |
| family=bag;stmt=join_where;threshold=20 | valid | 0.9736 | 1.0654 | 1.2373 | 5 | 938.6 |
| family=bag;stmt=join_where;threshold=200 | valid | 1.3948 | 1.4746 | 2.4881 | 5 | 678.1 |
| family=bag;stmt=join_where;threshold=2000 | valid | 3.1656 | 3.2721 | 3.4560 | 5 | 305.6 |
| family=bag;stmt=join_where;threshold=20000 | valid | 14.6885 | 19.0600 | 22.1371 | 5 | 52.5 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 6.0669 | 6.5675 | 8.4543 | 5 | 152.3 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 6.6597 | 7.4641 | 8.3713 | 5 | 134.0 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 0.9166 | 1.0877 | 1.7122 | 5 | 919.4 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 4.4584 | 4.7786 | 5.7275 | 5 | 209.3 |
| family=reuse;policy=default;refs=2;consumer=1% | valid | 4.7208 | 4.8223 | 10.9034 | 5 | 207.4 |
| family=reuse;policy=default;refs=2;consumer=100% | valid | 4.4507 | 6.1419 | 8.4020 | 5 | 162.8 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 4.8064 | 4.9065 | 5.7313 | 5 | 203.8 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 4.8962 | 5.2060 | 7.6252 | 5 | 192.1 |
| family=reuse;policy=default;refs=8;consumer=1% | valid | 4.9704 | 5.9973 | 6.0263 | 5 | 166.7 |
| family=reuse;policy=default;refs=8;consumer=100% | valid | 5.9218 | 6.1096 | 6.6994 | 5 | 163.7 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 6.1515 | 6.4718 | 6.8266 | 5 | 154.5 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 6.9223 | 7.9919 | 8.3871 | 5 | 125.1 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 4.3739 | 4.8967 | 5.0003 | 5 | 204.2 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 4.6573 | 5.1171 | 5.4243 | 5 | 195.4 |
| family=reuse;policy=mat;refs=2;consumer=1% | valid | 3.3720 | 4.6856 | 4.7424 | 5 | 213.4 |
| family=reuse;policy=mat;refs=2;consumer=100% | valid | 4.7042 | 4.8363 | 5.9637 | 5 | 206.8 |
| family=reuse;policy=mat;refs=2;consumer=100%;position=last | valid | 4.8244 | 4.9819 | 5.7042 | 5 | 200.7 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 4.7236 | 5.0619 | 6.7324 | 5 | 197.6 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 5.0567 | 5.5407 | 7.6530 | 5 | 180.5 |
| family=reuse;policy=mat;refs=8;consumer=1% | valid | 5.2435 | 5.3902 | 5.4304 | 5 | 185.5 |
| family=reuse;policy=mat;refs=8;consumer=100% | valid | 5.4298 | 5.6833 | 6.2250 | 5 | 176.0 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 4.8227 | 5.6004 | 8.3071 | 5 | 178.6 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 31.9130 | 33.2215 | 49.9895 | 5 | 30.1 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 0.9237 | 1.1817 | 2.5664 | 5 | 846.2 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 4.1997 | 4.4720 | 4.6922 | 5 | 223.6 |
| family=reuse;policy=notmat;refs=2;consumer=1% | valid | 1.2413 | 1.6276 | 2.0457 | 5 | 614.4 |
| family=reuse;policy=notmat;refs=2;consumer=100% | valid | 8.1851 | 11.5251 | 20.4123 | 5 | 86.8 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 1.2881 | 1.7496 | 2.1463 | 5 | 571.6 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 14.3121 | 23.1480 | 26.4190 | 5 | 43.2 |
| family=reuse;policy=notmat;refs=8;consumer=1% | valid | 2.7019 | 2.9625 | 3.0016 | 5 | 337.6 |
| family=reuse;policy=notmat;refs=8;consumer=100% | valid | 18.0397 | 20.4004 | 22.9434 | 5 | 49.0 |

## Execution order

- order seed `42196000`, 5 interleaved rounds
  1. family=reuse;policy=mat;refs=2;consumer=100%
  2. family=bag;stmt=cte_filtered_notmat;threshold=20
  3. family=reuse;policy=mat;refs=1;consumer=1%
  4. family=bag;stmt=cte_filtered_mat;threshold=200
  5. family=reuse;policy=notmat;refs=16;consumer=1%
  6. family=reuse;policy=mat;refs=16;consumer=1%
  7. family=bag;stmt=cte_filtered_notmat;threshold=100000
  8. family=bag;stmt=filtered_derived;threshold=180000
  9. family=bag;stmt=cte_joined_outer_mat;threshold=200
  10. family=bag;stmt=comma_join;threshold=180000
  11. family=bag;stmt=cte_filtered_mat;threshold=180000
  12. family=reuse;policy=default;refs=8;consumer=1%
  13. family=reuse;policy=default;refs=2;consumer=100%
  14. family=reuse;policy=default;refs=8;consumer=100%
  15. family=reuse;policy=notmat;refs=2;consumer=1%
  16. family=reuse;policy=default;refs=4;consumer=100%
  17. family=reuse;policy=mat;refs=16;consumer=100%
  18. family=bag;stmt=cte_filtered_notmat;threshold=2000
  19. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  20. family=bag;stmt=cte_filtered_notmat;threshold=20000
  21. family=bag;stmt=cte_joined_outer;threshold=200
  22. family=bag;stmt=comma_join;threshold=2000
  23. family=bag;stmt=join_where;threshold=20000
  24. family=bag;stmt=join_on;threshold=2000
  25. family=bag;stmt=filtered_derived;threshold=20
  26. family=bag;stmt=cte_filtered;threshold=200
  27. family=bag;stmt=cte_joined_outer_notmat;threshold=200
  28. family=bag;stmt=cte_filtered;threshold=2000
  29. family=reuse;policy=default;refs=4;consumer=1%
  30. family=bag;stmt=comma_join;threshold=100000
  31. family=reuse;policy=notmat;refs=1;consumer=1%
  32. family=bag;stmt=join_where;threshold=180000
  33. family=bag;stmt=filtered_derived;threshold=200
  34. family=bag;stmt=comma_join;threshold=20
  35. family=reuse;policy=default;refs=1;consumer=100%
  36. family=bag;stmt=join_where;threshold=100000
  37. family=bag;stmt=cte_joined_outer_mat;threshold=100000
  38. family=bag;stmt=join_where;threshold=2000
  39. family=reuse;policy=mat;refs=4;consumer=1%
  40. family=bag;stmt=cte_filtered_mat;threshold=20
  41. family=bag;stmt=join_where;threshold=200
  42. family=reuse;policy=notmat;refs=4;consumer=100%
  43. family=bag;stmt=filtered_derived;threshold=2000
  44. family=bag;stmt=cte_joined_outer_notmat;threshold=100000
  45. family=bag;stmt=cte_joined_outer;threshold=180000
  46. family=bag;stmt=cte_joined_outer;threshold=100000
  47. family=bag;stmt=filtered_derived;threshold=20000
  48. family=bag;stmt=join_on;threshold=200
  49. family=bag;stmt=cte_filtered;threshold=20000
  50. family=bag;stmt=join_on;threshold=180000
  51. family=reuse;policy=notmat;refs=16;consumer=100%
  52. family=bag;stmt=cte_filtered;threshold=180000
  53. family=bag;stmt=cte_filtered_notmat;threshold=180000
  54. family=reuse;policy=mat;refs=8;consumer=100%
  55. family=bag;stmt=join_on;threshold=20000
  56. family=reuse;policy=default;refs=1;consumer=1%
  57. family=reuse;policy=mat;refs=4;consumer=100%
  58. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  59. family=bag;stmt=comma_join;threshold=200
  60. family=reuse;policy=mat;refs=2;consumer=1%
  61. family=reuse;policy=notmat;refs=2;consumer=100%
  62. family=reuse;policy=notmat;refs=4;consumer=1%
  63. family=bag;stmt=cte_joined_outer_mat;threshold=20000
  64. family=bag;stmt=join_on;threshold=100000
  65. family=bag;stmt=cte_filtered_mat;threshold=100000
  66. family=bag;stmt=join_where;threshold=20
  67. family=reuse;policy=default;refs=16;consumer=100%
  68. family=bag;stmt=filtered_derived;threshold=100000
  69. family=bag;stmt=cte_filtered_mat;threshold=2000
  70. family=bag;stmt=cte_joined_outer_notmat;threshold=20000
  71. family=bag;stmt=join_on;threshold=20
  72. family=bag;stmt=cte_filtered;threshold=20
  73. family=bag;stmt=comma_join;threshold=20000
  74. family=reuse;policy=mat;refs=8;consumer=1%
  75. family=reuse;policy=default;refs=16;consumer=1%
  76. family=bag;stmt=cte_joined_outer_mat;threshold=20
  77. family=bag;stmt=cte_filtered_notmat;threshold=200
  78. family=reuse;policy=notmat;refs=1;consumer=100%
  79. family=bag;stmt=cte_joined_outer;threshold=20000
  80. family=bag;stmt=cte_joined_outer;threshold=20
  81. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  82. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  83. family=bag;stmt=cte_filtered;threshold=100000
  84. family=reuse;policy=notmat;refs=8;consumer=100%
  85. family=bag;stmt=cte_filtered_mat;threshold=20000
  86. family=reuse;policy=mat;refs=1;consumer=100%
  87. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  88. family=bag;stmt=cte_joined_outer;threshold=2000
  89. family=reuse;policy=notmat;refs=8;consumer=1%
  90. family=reuse;policy=default;refs=2;consumer=1%
  91. family=reuse;policy=mat;refs=2;consumer=100%;position=last

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
| family=reuse;policy=mat;refs=2;consumer=100%;position=last | r_mat_n2_c100 | valid | `450416f8dc65edb3` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=1% | r_mat_n4_c1 | valid | `91c3642bff5db3db` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=100% | r_mat_n4_c100 | valid | `91c3642bff5db3db` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=8;consumer=1% | r_mat_n8_c1 | valid | `0c9b689bfb53ff4f` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=8;consumer=100% | r_mat_n8_c100 | valid | `0c9b689bfb53ff4f` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=notmat;refs=16;consumer=1% | r_notmat_n16_c1 | valid | `a586237c6b012f4c` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `daea55c209e40a05` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `67feb88d971b83fe` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `b097208aea033185` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=2;consumer=1% | r_notmat_n2_c1 | valid | `b6d08c76c5629cb2` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=2;consumer=100% | r_notmat_n2_c100 | valid | `b5b1a9f8ee6f1d86` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `4e9a2c71c9021d34` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `b2dc573d0452c7ec` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=8;consumer=1% | r_notmat_n8_c1 | valid | `b7419850b4ea67ad` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=8;consumer=100% | r_notmat_n8_c100 | valid | `f9a803ecffcb61fe` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, VERBOSE, BUFFERS, FORMAT TEXT` supported on postgres
- one host, single client: no independent hardware, real-network or equal-total-budget comparison
- no concurrent writers: the isolated read contract only; a matched blended read/write phase is a recorded gap
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
