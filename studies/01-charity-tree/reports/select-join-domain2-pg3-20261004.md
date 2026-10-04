# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 91 cells, 91 valid, engine `postgres`, topology `single`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-domain2-pg3-20261004`
- run tag: `run/study-01/select-join-domain2-pg3-20261004`
- environment: `host-zenbook-ux5406sa`
- repo commit: `445fbc1646d58074c45cbaea57a83234132adc34`
- engine version: `PostgreSQL 17.11 (Debian 17.11-1.pgdg13+2) on x86_64-pc-linux-gnu, compiled by gcc (Debian 14.2.0-19) 14.2.0, 64-bit`
- database engine image: `docker.io/library/postgres:17.11` @ `7296f210ae81031ec955dbad9a67a84fe958572a2153b8d0826a647522904dc1`
- benchmark client image: `localhost/charitybench-select:select-join-domain2-pg3-20261004` @ `79bc753ff7e57e102982a323939c43e68e9b3664887683d3a5d5bf072650f4c1`
- inputs digest: `ccd615919f584e75b6fd4640f737bcf8fc0b7c3a1bf370b6d2d483399a86c704`
- dataset fingerprint: `fcf26e08f377ec51`
- client CPU quota: `100000 100000`, memory: `2147483648`
- setting `client_workers`: 1 sequential reader; 0 writers
- setting `connection_pool`: 2 per endpoint = reader/setup + reserve
- setting `customers`: 2000
- setting `dataset`: study01-person-donation-charity
- setting `distributed_rpc`: N/A on PostgreSQL single-node
- setting `endpoint_mode`: single endpoint
- setting `endpoints`: pg-single
- setting `fact_count`: 200000
- setting `index_set`: child rank index and parent-key index, primary keys; no per-spelling index
- setting `phase`: isolated
- setting `placement`: engine default
- setting `prepared_plan_policy`: pgx simple protocol; no prepared statements
- setting `scale`: cte-domain
- setting `selectivity_mapping`: rank <= threshold; rank is a deterministic shuffle of 1..fact_count
- setting `statement_budget_ms`: 120000
- setting `statistics`: ANALYZE after deterministic load
- setting `thresholds`: [20 200 2000 20000 100000 180000]
- setting `topology`: single
- setting `trials_per_cell`: 5
- setting `whole_run_count`: 1; no replicated winner in this file
- setting `writer_keys`: 1000
- setting `writer_target_rate_per_sec`: 10

### Endpoint distribution

A cluster arm must rotate its reader across nodes; a single seed endpoint would not be balanced cluster access.

| endpoint | operations | connections at end |
|---|---|---|
| pg-single | 914 | 1 |

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
| 20 (0.01%) | 0.6540 | 0.4178 | 0.4883 | 0.4654 | 0.4717 | 0.4807 | 0.6866 | 55.9168 | 0.5703 | 0.6343 |
| 200 (0.1%) | 0.7302 | 1.0121 | 0.8932 | 0.7962 | 0.6940 | 0.7118 | 1.1010 | 54.9834 | 0.9680 | 0.6717 |
| 2000 (1%) | 1.7240 | 1.7962 | 1.6968 | 1.6022 | 1.7417 | 1.5883 | 1.5808 | 46.0628 | 1.9358 | 1.8350 |
| 20000 (10%) | 7.2611 | 6.9573 | 7.7755 | 7.4579 | 10.9048 | 8.7815 | 7.3442 | 54.1713 | 7.5793 | 7.2235 |
| 100000 (50%) | 48.3034 | 48.6207 | 48.0486 | 45.9197 | 59.9158 | 54.0796 | 46.0135 | 65.1432 | 49.6292 | 47.9625 |
| 180000 (90%) | 77.9629 | 78.3122 | 72.0006 | 79.3914 | 71.9129 | 76.9339 | 75.2346 | 87.9106 | 86.7068 | 58.2125 |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | join_on vs join_where | 0.6540 | 0.4178 | 0.6388 | true |
| 20 | filtered_derived vs join_where | 0.6540 | 0.4883 | 0.7466 | true |
| 20 | cte_filtered vs join_where | 0.6540 | 0.4654 | 0.7116 | true |
| 20 | cte_filtered_mat vs join_where | 0.6540 | 0.4717 | 0.7212 | false |
| 20 | cte_filtered_notmat vs join_where | 0.6540 | 0.4807 | 0.7350 | true |
| 20 | cte_joined_outer vs join_where | 0.6540 | 0.6866 | 1.0498 | true |
| 20 | cte_joined_outer_mat vs join_where | 0.6540 | 55.9168 | 85.4967 | false |
| 20 | cte_joined_outer_notmat vs join_where | 0.6540 | 0.5703 | 0.8720 | true |
| 20 | comma_join vs join_where | 0.6540 | 0.6343 | 0.9699 | true |
| 200 | join_on vs join_where | 0.7302 | 1.0121 | 1.3861 | true |
| 200 | filtered_derived vs join_where | 0.7302 | 0.8932 | 1.2233 | true |
| 200 | cte_filtered vs join_where | 0.7302 | 0.7962 | 1.0904 | true |
| 200 | cte_filtered_mat vs join_where | 0.7302 | 0.6940 | 0.9504 | false |
| 200 | cte_filtered_notmat vs join_where | 0.7302 | 0.7118 | 0.9748 | true |
| 200 | cte_joined_outer vs join_where | 0.7302 | 1.1010 | 1.5077 | true |
| 200 | cte_joined_outer_mat vs join_where | 0.7302 | 54.9834 | 75.2983 | false |
| 200 | cte_joined_outer_notmat vs join_where | 0.7302 | 0.9680 | 1.3256 | true |
| 200 | comma_join vs join_where | 0.7302 | 0.6717 | 0.9198 | true |
| 2000 | join_on vs join_where | 1.7240 | 1.7962 | 1.0419 | true |
| 2000 | filtered_derived vs join_where | 1.7240 | 1.6968 | 0.9842 | true |
| 2000 | cte_filtered vs join_where | 1.7240 | 1.6022 | 0.9293 | true |
| 2000 | cte_filtered_mat vs join_where | 1.7240 | 1.7417 | 1.0102 | false |
| 2000 | cte_filtered_notmat vs join_where | 1.7240 | 1.5883 | 0.9213 | true |
| 2000 | cte_joined_outer vs join_where | 1.7240 | 1.5808 | 0.9169 | true |
| 2000 | cte_joined_outer_mat vs join_where | 1.7240 | 46.0628 | 26.7179 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 1.7240 | 1.9358 | 1.1228 | true |
| 2000 | comma_join vs join_where | 1.7240 | 1.8350 | 1.0644 | true |
| 20000 | join_on vs join_where | 7.2611 | 6.9573 | 0.9582 | true |
| 20000 | filtered_derived vs join_where | 7.2611 | 7.7755 | 1.0708 | true |
| 20000 | cte_filtered vs join_where | 7.2611 | 7.4579 | 1.0271 | true |
| 20000 | cte_filtered_mat vs join_where | 7.2611 | 10.9048 | 1.5018 | false |
| 20000 | cte_filtered_notmat vs join_where | 7.2611 | 8.7815 | 1.2094 | true |
| 20000 | cte_joined_outer vs join_where | 7.2611 | 7.3442 | 1.0114 | true |
| 20000 | cte_joined_outer_mat vs join_where | 7.2611 | 54.1713 | 7.4604 | false |
| 20000 | cte_joined_outer_notmat vs join_where | 7.2611 | 7.5793 | 1.0438 | true |
| 20000 | comma_join vs join_where | 7.2611 | 7.2235 | 0.9948 | true |
| 100000 | join_on vs join_where | 48.3034 | 48.6207 | 1.0066 | true |
| 100000 | filtered_derived vs join_where | 48.3034 | 48.0486 | 0.9947 | true |
| 100000 | cte_filtered vs join_where | 48.3034 | 45.9197 | 0.9507 | true |
| 100000 | cte_filtered_mat vs join_where | 48.3034 | 59.9158 | 1.2404 | false |
| 100000 | cte_filtered_notmat vs join_where | 48.3034 | 54.0796 | 1.1196 | true |
| 100000 | cte_joined_outer vs join_where | 48.3034 | 46.0135 | 0.9526 | true |
| 100000 | cte_joined_outer_mat vs join_where | 48.3034 | 65.1432 | 1.3486 | false |
| 100000 | cte_joined_outer_notmat vs join_where | 48.3034 | 49.6292 | 1.0274 | true |
| 100000 | comma_join vs join_where | 48.3034 | 47.9625 | 0.9929 | true |
| 180000 | join_on vs join_where | 77.9629 | 78.3122 | 1.0045 | true |
| 180000 | filtered_derived vs join_where | 77.9629 | 72.0006 | 0.9235 | true |
| 180000 | cte_filtered vs join_where | 77.9629 | 79.3914 | 1.0183 | true |
| 180000 | cte_filtered_mat vs join_where | 77.9629 | 71.9129 | 0.9224 | false |
| 180000 | cte_filtered_notmat vs join_where | 77.9629 | 76.9339 | 0.9868 | true |
| 180000 | cte_joined_outer vs join_where | 77.9629 | 75.2346 | 0.9650 | true |
| 180000 | cte_joined_outer_mat vs join_where | 77.9629 | 87.9106 | 1.1276 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 77.9629 | 86.7068 | 1.1122 | true |
| 180000 | comma_join vs join_where | 77.9629 | 58.2125 | 0.7467 | true |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_filtered_mat vs cte_filtered | 0.4654 | 0.4717 | 1.0135 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_filtered_notmat vs cte_filtered | 0.4654 | 0.4807 | 1.0329 | true | policy=not_materialized; observed_cte_scan=false |
| 20 | cte_joined_outer_mat vs cte_joined_outer | 0.6866 | 55.9168 | 81.4411 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 0.6866 | 0.5703 | 0.8306 | true | policy=not_materialized; observed_cte_scan=false |
| 200 | cte_filtered_mat vs cte_filtered | 0.7962 | 0.6940 | 0.8717 | false | policy=materialized; observed_cte_scan=true |
| 200 | cte_filtered_notmat vs cte_filtered | 0.7962 | 0.7118 | 0.8940 | true | policy=not_materialized; observed_cte_scan=false |
| 200 | cte_joined_outer_mat vs cte_joined_outer | 1.1010 | 54.9834 | 49.9410 | false | policy=materialized; observed_cte_scan=true |
| 200 | cte_joined_outer_notmat vs cte_joined_outer | 1.1010 | 0.9680 | 0.8792 | true | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_filtered_mat vs cte_filtered | 1.6022 | 1.7417 | 1.0871 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_filtered_notmat vs cte_filtered | 1.6022 | 1.5883 | 0.9913 | true | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 1.5808 | 46.0628 | 29.1387 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 1.5808 | 1.9358 | 1.2246 | true | policy=not_materialized; observed_cte_scan=false |
| 20000 | cte_filtered_mat vs cte_filtered | 7.4579 | 10.9048 | 1.4622 | false | policy=materialized; observed_cte_scan=true |
| 20000 | cte_filtered_notmat vs cte_filtered | 7.4579 | 8.7815 | 1.1775 | true | policy=not_materialized; observed_cte_scan=false |
| 20000 | cte_joined_outer_mat vs cte_joined_outer | 7.3442 | 54.1713 | 7.3761 | false | policy=materialized; observed_cte_scan=true |
| 20000 | cte_joined_outer_notmat vs cte_joined_outer | 7.3442 | 7.5793 | 1.0320 | true | policy=not_materialized; observed_cte_scan=false |
| 100000 | cte_filtered_mat vs cte_filtered | 45.9197 | 59.9158 | 1.3048 | false | policy=materialized; observed_cte_scan=true |
| 100000 | cte_filtered_notmat vs cte_filtered | 45.9197 | 54.0796 | 1.1777 | true | policy=not_materialized; observed_cte_scan=false |
| 100000 | cte_joined_outer_mat vs cte_joined_outer | 46.0135 | 65.1432 | 1.4157 | false | policy=materialized; observed_cte_scan=true |
| 100000 | cte_joined_outer_notmat vs cte_joined_outer | 46.0135 | 49.6292 | 1.0786 | true | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_filtered_mat vs cte_filtered | 79.3914 | 71.9129 | 0.9058 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_filtered_notmat vs cte_filtered | 79.3914 | 76.9339 | 0.9690 | true | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 75.2346 | 87.9106 | 1.1685 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 75.2346 | 86.7068 | 1.1525 | true | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 0.4957 | 1 | `549a07108ab01ba9` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 2.5516 | 1 | `2ce8b85a7b60c88e` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | 2.2806 | 2 | `b307c41eb54cf452` | policy=default; observed_cte_scan=true | valid |
| default | 2 | 100% | 2.3475 | 2 | `b307c41eb54cf452` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 1% | 2.9179 | 4 | `ccfe1af393deae98` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 3.0049 | 4 | `ccfe1af393deae98` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | 2.9520 | 8 | `c8c73778097931dc` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 100% | 3.2185 | 8 | `c8c73778097931dc` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 1% | 3.2959 | 16 | `b0c356323d30c425` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 3.4677 | 16 | `b0c356323d30c425` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 2.6219 | 1 | `eaa7d6058aaa5e03` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 2.6474 | 1 | `eaa7d6058aaa5e03` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | 2.5742 | 2 | `b307c41eb54cf452` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 100% | 2.6991 | 2 | `b307c41eb54cf452` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 1% | 2.5027 | 4 | `ccfe1af393deae98` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 2.9640 | 4 | `ccfe1af393deae98` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | 3.1591 | 8 | `c8c73778097931dc` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 100% | 3.1196 | 8 | `c8c73778097931dc` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 1% | 3.3769 | 16 | `b0c356323d30c425` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 3.7104 | 16 | `b0c356323d30c425` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 0.4840 | 1 | `549a07108ab01ba9` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 2.3495 | 1 | `2ce8b85a7b60c88e` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | 0.6854 | 2 | `28cc913d71d0d72b` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 100% | 4.3344 | 2 | `58596ee2bf1338ac` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 1% | 0.9101 | 4 | `9f850e53560baafb` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 7.2314 | 4 | `5bb21449fb3f5342` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | 1.3829 | 8 | `477f7fbaf1c627cb` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 100% | 8.9714 | 8 | `1353860480f2b2fb` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 1% | 2.8243 | 16 | `1c38b58d8cec08d3` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 16.2864 | 16 | `9779fc7cd40b6621` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 5.2892 | false |
| notmat | 1 | 1% | 0.9764 | true |
| mat | 1 | 100% | 1.0376 | false |
| notmat | 1 | 100% | 0.9208 | true |
| mat | 2 | 1% | 1.1287 | true |
| notmat | 2 | 1% | 0.3005 | false |
| mat | 2 | 100% | 1.1498 | true |
| notmat | 2 | 100% | 1.8464 | false |
| mat | 4 | 1% | 0.8577 | true |
| notmat | 4 | 1% | 0.3119 | false |
| mat | 4 | 100% | 0.9864 | true |
| notmat | 4 | 100% | 2.4066 | false |
| mat | 8 | 1% | 1.0702 | true |
| notmat | 8 | 1% | 0.4685 | false |
| mat | 8 | 100% | 0.9693 | true |
| notmat | 8 | 100% | 2.7874 | false |
| mat | 16 | 1% | 1.0246 | true |
| notmat | 16 | 1% | 0.8569 | false |
| mat | 16 | 100% | 1.0700 | true |
| notmat | 16 | 100% | 4.6966 | false |

## Repeated-design position control

| first median ms | last median ms | ratio (last/first) | same plan |
|---|---|---|---|
| 1.012149 | 0.700468 | 0.6921 | true |

## Latency samples and derived throughput

_The reader is sequential, so `queries/s` is 1000/median_ms, a derived rate and not a concurrent-capacity test. min/max are the five sample extremes._

| axis | status | min ms | median ms | max ms | samples | derived queries/s |
|---|---|---|---|---|---|---|
| family=bag;stmt=comma_join;threshold=100000 | valid | 41.5582 | 47.9625 | 67.0993 | 5 | 20.8 |
| family=bag;stmt=comma_join;threshold=180000 | valid | 54.7828 | 58.2125 | 67.2842 | 5 | 17.2 |
| family=bag;stmt=comma_join;threshold=20 | valid | 0.5429 | 0.6343 | 1.1752 | 5 | 1576.5 |
| family=bag;stmt=comma_join;threshold=200 | valid | 0.6522 | 0.6717 | 0.8634 | 5 | 1488.8 |
| family=bag;stmt=comma_join;threshold=2000 | valid | 1.8036 | 1.8350 | 1.8861 | 5 | 544.9 |
| family=bag;stmt=comma_join;threshold=20000 | valid | 6.9731 | 7.2235 | 37.8248 | 5 | 138.4 |
| family=bag;stmt=cte_filtered;threshold=100000 | valid | 38.8672 | 45.9197 | 46.7787 | 5 | 21.8 |
| family=bag;stmt=cte_filtered;threshold=180000 | valid | 48.6597 | 79.3914 | 94.3483 | 5 | 12.6 |
| family=bag;stmt=cte_filtered;threshold=20 | valid | 0.3815 | 0.4654 | 0.5326 | 5 | 2148.7 |
| family=bag;stmt=cte_filtered;threshold=200 | valid | 0.7098 | 0.7962 | 0.8363 | 5 | 1256.0 |
| family=bag;stmt=cte_filtered;threshold=2000 | valid | 1.4592 | 1.6022 | 1.6264 | 5 | 624.1 |
| family=bag;stmt=cte_filtered;threshold=20000 | valid | 6.5092 | 7.4579 | 8.2224 | 5 | 134.1 |
| family=bag;stmt=cte_filtered_mat;threshold=100000 | valid | 50.8881 | 59.9158 | 61.3274 | 5 | 16.7 |
| family=bag;stmt=cte_filtered_mat;threshold=180000 | valid | 67.4267 | 71.9129 | 78.3674 | 5 | 13.9 |
| family=bag;stmt=cte_filtered_mat;threshold=20 | valid | 0.3720 | 0.4717 | 0.4910 | 5 | 2120.1 |
| family=bag;stmt=cte_filtered_mat;threshold=200 | valid | 0.6500 | 0.6940 | 36.1158 | 5 | 1440.9 |
| family=bag;stmt=cte_filtered_mat;threshold=2000 | valid | 1.6402 | 1.7417 | 1.8642 | 5 | 574.2 |
| family=bag;stmt=cte_filtered_mat;threshold=20000 | valid | 9.5184 | 10.9048 | 32.4911 | 5 | 91.7 |
| family=bag;stmt=cte_filtered_notmat;threshold=100000 | valid | 46.6981 | 54.0796 | 60.0383 | 5 | 18.5 |
| family=bag;stmt=cte_filtered_notmat;threshold=180000 | valid | 44.3083 | 76.9339 | 94.4849 | 5 | 13.0 |
| family=bag;stmt=cte_filtered_notmat;threshold=20 | valid | 0.4305 | 0.4807 | 0.6635 | 5 | 2080.3 |
| family=bag;stmt=cte_filtered_notmat;threshold=200 | valid | 0.6609 | 0.7118 | 0.7809 | 5 | 1404.8 |
| family=bag;stmt=cte_filtered_notmat;threshold=2000 | valid | 1.5559 | 1.5883 | 1.7871 | 5 | 629.6 |
| family=bag;stmt=cte_filtered_notmat;threshold=20000 | valid | 6.5833 | 8.7815 | 26.0033 | 5 | 113.9 |
| family=bag;stmt=cte_joined_outer;threshold=100000 | valid | 44.0160 | 46.0135 | 73.8857 | 5 | 21.7 |
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 67.0115 | 75.2346 | 89.8269 | 5 | 13.3 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 0.6129 | 0.6866 | 0.8152 | 5 | 1456.5 |
| family=bag;stmt=cte_joined_outer;threshold=200 | valid | 0.9621 | 1.1010 | 1.4210 | 5 | 908.3 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 1.5345 | 1.5808 | 1.8378 | 5 | 632.6 |
| family=bag;stmt=cte_joined_outer;threshold=20000 | valid | 6.9132 | 7.3442 | 8.3704 | 5 | 136.2 |
| family=bag;stmt=cte_joined_outer_mat;threshold=100000 | valid | 60.7129 | 65.1432 | 68.1764 | 5 | 15.4 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 82.8451 | 87.9106 | 91.8754 | 5 | 11.4 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 44.1972 | 55.9168 | 60.3749 | 5 | 17.9 |
| family=bag;stmt=cte_joined_outer_mat;threshold=200 | valid | 43.3713 | 54.9834 | 57.6821 | 5 | 18.2 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 43.8909 | 46.0628 | 60.8578 | 5 | 21.7 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20000 | valid | 46.8936 | 54.1713 | 64.4744 | 5 | 18.5 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=100000 | valid | 44.2197 | 49.6292 | 53.7234 | 5 | 20.1 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 79.4379 | 86.7068 | 93.4509 | 5 | 11.5 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 0.5423 | 0.5703 | 0.5791 | 5 | 1753.5 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=200 | valid | 0.9361 | 0.9680 | 1.0148 | 5 | 1033.1 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 1.8543 | 1.9358 | 2.1244 | 5 | 516.6 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20000 | valid | 6.9881 | 7.5793 | 11.9943 | 5 | 131.9 |
| family=bag;stmt=filtered_derived;threshold=100000 | valid | 43.8370 | 48.0486 | 62.0364 | 5 | 20.8 |
| family=bag;stmt=filtered_derived;threshold=180000 | valid | 57.3512 | 72.0006 | 88.5371 | 5 | 13.9 |
| family=bag;stmt=filtered_derived;threshold=20 | valid | 0.4243 | 0.4883 | 0.5490 | 5 | 2047.9 |
| family=bag;stmt=filtered_derived;threshold=200 | valid | 0.8165 | 0.8932 | 1.1115 | 5 | 1119.5 |
| family=bag;stmt=filtered_derived;threshold=2000 | valid | 1.4597 | 1.6968 | 1.7830 | 5 | 589.3 |
| family=bag;stmt=filtered_derived;threshold=20000 | valid | 6.6285 | 7.7755 | 10.6313 | 5 | 128.6 |
| family=bag;stmt=join_on;threshold=100000 | valid | 38.4825 | 48.6207 | 69.7445 | 5 | 20.6 |
| family=bag;stmt=join_on;threshold=180000 | valid | 50.8843 | 78.3122 | 78.5909 | 5 | 12.8 |
| family=bag;stmt=join_on;threshold=20 | valid | 0.3984 | 0.4178 | 0.5419 | 5 | 2393.7 |
| family=bag;stmt=join_on;threshold=200 | valid | 0.9465 | 1.0121 | 1.0929 | 5 | 988.0 |
| family=bag;stmt=join_on;threshold=2000 | valid | 1.7669 | 1.7962 | 2.7784 | 5 | 556.7 |
| family=bag;stmt=join_on;threshold=20000 | valid | 6.5836 | 6.9573 | 7.9621 | 5 | 143.7 |
| family=bag;stmt=join_on;threshold=200;position=last | valid | 0.6220 | 0.7005 | 0.7759 | 5 | 1427.6 |
| family=bag;stmt=join_where;threshold=100000 | valid | 41.7719 | 48.3034 | 49.8237 | 5 | 20.7 |
| family=bag;stmt=join_where;threshold=180000 | valid | 45.0144 | 77.9629 | 87.2551 | 5 | 12.8 |
| family=bag;stmt=join_where;threshold=20 | valid | 0.5490 | 0.6540 | 0.8124 | 5 | 1529.0 |
| family=bag;stmt=join_where;threshold=200 | valid | 0.6740 | 0.7302 | 0.7947 | 5 | 1369.5 |
| family=bag;stmt=join_where;threshold=2000 | valid | 1.6687 | 1.7240 | 1.8084 | 5 | 580.0 |
| family=bag;stmt=join_where;threshold=20000 | valid | 6.3288 | 7.2611 | 7.7298 | 5 | 137.7 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 3.1883 | 3.2959 | 3.3297 | 5 | 303.4 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 3.4091 | 3.4677 | 3.7410 | 5 | 288.4 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 0.4725 | 0.4957 | 0.7551 | 5 | 2017.3 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 2.5017 | 2.5516 | 3.1894 | 5 | 391.9 |
| family=reuse;policy=default;refs=2;consumer=1% | valid | 2.2296 | 2.2806 | 2.3154 | 5 | 438.5 |
| family=reuse;policy=default;refs=2;consumer=100% | valid | 2.2206 | 2.3475 | 3.9764 | 5 | 426.0 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 2.6867 | 2.9179 | 3.2969 | 5 | 342.7 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 2.7393 | 3.0049 | 3.7010 | 5 | 332.8 |
| family=reuse;policy=default;refs=8;consumer=1% | valid | 2.8137 | 2.9520 | 3.0983 | 5 | 338.7 |
| family=reuse;policy=default;refs=8;consumer=100% | valid | 3.0688 | 3.2185 | 3.4931 | 5 | 310.7 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 3.2038 | 3.3769 | 3.4619 | 5 | 296.1 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 3.5266 | 3.7104 | 3.9563 | 5 | 269.5 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 2.5345 | 2.6219 | 2.8755 | 5 | 381.4 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 2.4771 | 2.6474 | 3.0820 | 5 | 377.7 |
| family=reuse;policy=mat;refs=2;consumer=1% | valid | 2.4805 | 2.5742 | 2.7411 | 5 | 388.5 |
| family=reuse;policy=mat;refs=2;consumer=100% | valid | 2.6139 | 2.6991 | 2.7176 | 5 | 370.5 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 2.4611 | 2.5027 | 2.7543 | 5 | 399.6 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 2.7439 | 2.9640 | 3.0432 | 5 | 337.4 |
| family=reuse;policy=mat;refs=8;consumer=1% | valid | 2.9075 | 3.1591 | 3.5452 | 5 | 316.5 |
| family=reuse;policy=mat;refs=8;consumer=100% | valid | 3.0274 | 3.1196 | 3.6207 | 5 | 320.6 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 2.6980 | 2.8243 | 3.9119 | 5 | 354.1 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 14.1382 | 16.2864 | 17.1263 | 5 | 61.4 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 0.4148 | 0.4840 | 0.6714 | 5 | 2066.0 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 2.1938 | 2.3495 | 2.4231 | 5 | 425.6 |
| family=reuse;policy=notmat;refs=2;consumer=1% | valid | 0.6351 | 0.6854 | 0.7833 | 5 | 1459.0 |
| family=reuse;policy=notmat;refs=2;consumer=100% | valid | 4.1595 | 4.3344 | 6.0683 | 5 | 230.7 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 0.7837 | 0.9101 | 0.9580 | 5 | 1098.8 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 6.8263 | 7.2314 | 7.8672 | 5 | 138.3 |
| family=reuse;policy=notmat;refs=8;consumer=1% | valid | 1.3264 | 1.3829 | 1.4181 | 5 | 723.1 |
| family=reuse;policy=notmat;refs=8;consumer=100% | valid | 8.8722 | 8.9714 | 14.4431 | 5 | 111.5 |

## Matched blended workload

- phase: `isolated`; target offered rate: 10 updates/s
- offered 0, committed 0, failed 0 over 0.000 s
- offered rate 0.000/s, completed rate 0.000/s
- per key: no keys touched
- ledger: 0 keys reconciled, 0 mismatch

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42173002`, 5 interleaved rounds
  1. family=bag;stmt=join_on;threshold=200
  2. family=bag;stmt=join_on;threshold=20
  3. family=bag;stmt=comma_join;threshold=2000
  4. family=reuse;policy=notmat;refs=8;consumer=1%
  5. family=bag;stmt=cte_filtered;threshold=20
  6. family=bag;stmt=join_where;threshold=2000
  7. family=reuse;policy=notmat;refs=1;consumer=1%
  8. family=reuse;policy=default;refs=16;consumer=100%
  9. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  10. family=bag;stmt=join_where;threshold=100000
  11. family=reuse;policy=mat;refs=2;consumer=100%
  12. family=bag;stmt=cte_joined_outer;threshold=180000
  13. family=reuse;policy=notmat;refs=16;consumer=1%
  14. family=bag;stmt=cte_filtered_notmat;threshold=20
  15. family=reuse;policy=notmat;refs=16;consumer=100%
  16. family=reuse;policy=notmat;refs=8;consumer=100%
  17. family=bag;stmt=comma_join;threshold=180000
  18. family=reuse;policy=mat;refs=4;consumer=100%
  19. family=bag;stmt=cte_filtered_notmat;threshold=2000
  20. family=bag;stmt=cte_filtered_mat;threshold=200
  21. family=bag;stmt=filtered_derived;threshold=100000
  22. family=reuse;policy=mat;refs=8;consumer=1%
  23. family=reuse;policy=default;refs=2;consumer=1%
  24. family=bag;stmt=cte_joined_outer_notmat;threshold=20000
  25. family=reuse;policy=default;refs=16;consumer=1%
  26. family=bag;stmt=filtered_derived;threshold=180000
  27. family=bag;stmt=cte_joined_outer;threshold=200
  28. family=reuse;policy=mat;refs=4;consumer=1%
  29. family=bag;stmt=cte_filtered_notmat;threshold=200
  30. family=bag;stmt=filtered_derived;threshold=20
  31. family=bag;stmt=join_on;threshold=180000
  32. family=reuse;policy=default;refs=8;consumer=100%
  33. family=bag;stmt=filtered_derived;threshold=2000
  34. family=bag;stmt=cte_joined_outer_mat;threshold=20000
  35. family=bag;stmt=comma_join;threshold=20
  36. family=reuse;policy=notmat;refs=1;consumer=100%
  37. family=bag;stmt=join_where;threshold=200
  38. family=reuse;policy=notmat;refs=2;consumer=1%
  39. family=bag;stmt=join_where;threshold=20000
  40. family=reuse;policy=mat;refs=2;consumer=1%
  41. family=reuse;policy=notmat;refs=4;consumer=1%
  42. family=bag;stmt=cte_filtered;threshold=200
  43. family=bag;stmt=cte_joined_outer;threshold=20000
  44. family=reuse;policy=mat;refs=1;consumer=100%
  45. family=bag;stmt=filtered_derived;threshold=20000
  46. family=reuse;policy=default;refs=4;consumer=1%
  47. family=bag;stmt=comma_join;threshold=200
  48. family=bag;stmt=cte_filtered_notmat;threshold=20000
  49. family=reuse;policy=mat;refs=16;consumer=1%
  50. family=reuse;policy=notmat;refs=2;consumer=100%
  51. family=bag;stmt=cte_filtered_notmat;threshold=100000
  52. family=bag;stmt=cte_joined_outer;threshold=20
  53. family=bag;stmt=cte_filtered;threshold=180000
  54. family=reuse;policy=mat;refs=8;consumer=100%
  55. family=bag;stmt=cte_filtered_mat;threshold=2000
  56. family=bag;stmt=join_where;threshold=180000
  57. family=bag;stmt=join_on;threshold=100000
  58. family=reuse;policy=default;refs=1;consumer=100%
  59. family=bag;stmt=cte_filtered_notmat;threshold=180000
  60. family=reuse;policy=notmat;refs=4;consumer=100%
  61. family=bag;stmt=cte_joined_outer_mat;threshold=20
  62. family=bag;stmt=cte_filtered_mat;threshold=180000
  63. family=reuse;policy=default;refs=8;consumer=1%
  64. family=reuse;policy=default;refs=2;consumer=100%
  65. family=bag;stmt=cte_joined_outer;threshold=2000
  66. family=bag;stmt=cte_filtered_mat;threshold=100000
  67. family=reuse;policy=mat;refs=1;consumer=1%
  68. family=bag;stmt=comma_join;threshold=100000
  69. family=bag;stmt=cte_joined_outer_mat;threshold=100000
  70. family=bag;stmt=join_where;threshold=20
  71. family=bag;stmt=cte_joined_outer;threshold=100000
  72. family=bag;stmt=comma_join;threshold=20000
  73. family=bag;stmt=cte_joined_outer_mat;threshold=200
  74. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  75. family=bag;stmt=cte_joined_outer_notmat;threshold=100000
  76. family=bag;stmt=join_on;threshold=20000
  77. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  78. family=bag;stmt=cte_filtered;threshold=2000
  79. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  80. family=bag;stmt=cte_joined_outer_notmat;threshold=200
  81. family=bag;stmt=cte_filtered;threshold=20000
  82. family=bag;stmt=filtered_derived;threshold=200
  83. family=reuse;policy=mat;refs=16;consumer=100%
  84. family=bag;stmt=cte_filtered_mat;threshold=20
  85. family=reuse;policy=default;refs=1;consumer=1%
  86. family=bag;stmt=cte_filtered_mat;threshold=20000
  87. family=bag;stmt=join_on;threshold=2000
  88. family=bag;stmt=cte_filtered;threshold=100000
  89. family=reuse;policy=default;refs=4;consumer=100%
  90. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
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
| family=bag;stmt=comma_join;threshold=100000 | j06_comma_join | valid | `4bb1561e58681728` |  |
| family=bag;stmt=comma_join;threshold=180000 | j06_comma_join | valid | `18e11c5d2e685630` |  |
| family=bag;stmt=comma_join;threshold=20 | j06_comma_join | valid | `61819ba7d0f5f952` |  |
| family=bag;stmt=comma_join;threshold=200 | j06_comma_join | valid | `41f23de3437828d5` |  |
| family=bag;stmt=comma_join;threshold=2000 | j06_comma_join | valid | `5573434208edb047` |  |
| family=bag;stmt=comma_join;threshold=20000 | j06_comma_join | valid | `09a028d269057a3e` |  |
| family=bag;stmt=cte_filtered;threshold=100000 | j04_cte_filtered | valid | `4bb1561e58681728` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=180000 | j04_cte_filtered | valid | `18e11c5d2e685630` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=20 | j04_cte_filtered | valid | `61819ba7d0f5f952` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=200 | j04_cte_filtered | valid | `41f23de3437828d5` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=2000 | j04_cte_filtered | valid | `5573434208edb047` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=20000 | j04_cte_filtered | valid | `09a028d269057a3e` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_mat;threshold=100000 | j04m_cte_filtered_mat | valid | `4b4ce4c83ae7fcc8` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=180000 | j04m_cte_filtered_mat | valid | `6a10a6fd7e521978` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=20 | j04m_cte_filtered_mat | valid | `d4ca4a9d1c0fd3e3` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=200 | j04m_cte_filtered_mat | valid | `a398c2b3383a151d` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=2000 | j04m_cte_filtered_mat | valid | `d4df6a0ef2f05082` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=20000 | j04m_cte_filtered_mat | valid | `5db776896b6f593f` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_notmat;threshold=100000 | j04n_cte_filtered_notmat | valid | `4bb1561e58681728` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=180000 | j04n_cte_filtered_notmat | valid | `18e11c5d2e685630` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=20 | j04n_cte_filtered_notmat | valid | `61819ba7d0f5f952` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=200 | j04n_cte_filtered_notmat | valid | `41f23de3437828d5` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=2000 | j04n_cte_filtered_notmat | valid | `5573434208edb047` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=20000 | j04n_cte_filtered_notmat | valid | `09a028d269057a3e` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=100000 | j05_cte_joined_outer | valid | `4bb1561e58681728` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=180000 | j05_cte_joined_outer | valid | `18e11c5d2e685630` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20 | j05_cte_joined_outer | valid | `61819ba7d0f5f952` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=200 | j05_cte_joined_outer | valid | `41f23de3437828d5` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=2000 | j05_cte_joined_outer | valid | `5573434208edb047` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20000 | j05_cte_joined_outer | valid | `09a028d269057a3e` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_mat;threshold=100000 | j05m_cte_joined_outer_mat | valid | `7f87bace4cf1da4c` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | j05m_cte_joined_outer_mat | valid | `7f87bace4cf1da4c` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | j05m_cte_joined_outer_mat | valid | `7f87bace4cf1da4c` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=200 | j05m_cte_joined_outer_mat | valid | `7f87bace4cf1da4c` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | j05m_cte_joined_outer_mat | valid | `7f87bace4cf1da4c` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20000 | j05m_cte_joined_outer_mat | valid | `7f87bace4cf1da4c` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_notmat;threshold=100000 | j05n_cte_joined_outer_notmat | valid | `4bb1561e58681728` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | j05n_cte_joined_outer_notmat | valid | `18e11c5d2e685630` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | j05n_cte_joined_outer_notmat | valid | `61819ba7d0f5f952` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=200 | j05n_cte_joined_outer_notmat | valid | `41f23de3437828d5` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | j05n_cte_joined_outer_notmat | valid | `5573434208edb047` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20000 | j05n_cte_joined_outer_notmat | valid | `09a028d269057a3e` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=filtered_derived;threshold=100000 | j03_filtered_derived | valid | `4bb1561e58681728` |  |
| family=bag;stmt=filtered_derived;threshold=180000 | j03_filtered_derived | valid | `18e11c5d2e685630` |  |
| family=bag;stmt=filtered_derived;threshold=20 | j03_filtered_derived | valid | `61819ba7d0f5f952` |  |
| family=bag;stmt=filtered_derived;threshold=200 | j03_filtered_derived | valid | `41f23de3437828d5` |  |
| family=bag;stmt=filtered_derived;threshold=2000 | j03_filtered_derived | valid | `5573434208edb047` |  |
| family=bag;stmt=filtered_derived;threshold=20000 | j03_filtered_derived | valid | `09a028d269057a3e` |  |
| family=bag;stmt=join_on;threshold=100000 | j02_join_on | valid | `4bb1561e58681728` |  |
| family=bag;stmt=join_on;threshold=180000 | j02_join_on | valid | `18e11c5d2e685630` |  |
| family=bag;stmt=join_on;threshold=20 | j02_join_on | valid | `61819ba7d0f5f952` |  |
| family=bag;stmt=join_on;threshold=200 | j02_join_on | valid | `41f23de3437828d5` |  |
| family=bag;stmt=join_on;threshold=2000 | j02_join_on | valid | `5573434208edb047` |  |
| family=bag;stmt=join_on;threshold=20000 | j02_join_on | valid | `09a028d269057a3e` |  |
| family=bag;stmt=join_on;threshold=200;position=last | j02_join_on | valid | `41f23de3437828d5` |  |
| family=bag;stmt=join_where;threshold=100000 | j01_join_where | valid | `4bb1561e58681728` |  |
| family=bag;stmt=join_where;threshold=180000 | j01_join_where | valid | `18e11c5d2e685630` |  |
| family=bag;stmt=join_where;threshold=20 | j01_join_where | valid | `61819ba7d0f5f952` |  |
| family=bag;stmt=join_where;threshold=200 | j01_join_where | valid | `41f23de3437828d5` |  |
| family=bag;stmt=join_where;threshold=2000 | j01_join_where | valid | `5573434208edb047` |  |
| family=bag;stmt=join_where;threshold=20000 | j01_join_where | valid | `09a028d269057a3e` |  |
| family=reuse;policy=default;refs=16;consumer=1% | r_default_n16_c1 | valid | `b0c356323d30c425` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=16;consumer=100% | r_default_n16_c100 | valid | `b0c356323d30c425` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=1;consumer=1% | r_default_n1_c1 | valid | `549a07108ab01ba9` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=1;consumer=100% | r_default_n1_c100 | valid | `2ce8b85a7b60c88e` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=2;consumer=1% | r_default_n2_c1 | valid | `b307c41eb54cf452` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=2;consumer=100% | r_default_n2_c100 | valid | `b307c41eb54cf452` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=1% | r_default_n4_c1 | valid | `ccfe1af393deae98` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=100% | r_default_n4_c100 | valid | `ccfe1af393deae98` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=8;consumer=1% | r_default_n8_c1 | valid | `c8c73778097931dc` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=8;consumer=100% | r_default_n8_c100 | valid | `c8c73778097931dc` | policy=default; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=1% | r_mat_n16_c1 | valid | `b0c356323d30c425` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=100% | r_mat_n16_c100 | valid | `b0c356323d30c425` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=1% | r_mat_n1_c1 | valid | `eaa7d6058aaa5e03` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=100% | r_mat_n1_c100 | valid | `eaa7d6058aaa5e03` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=2;consumer=1% | r_mat_n2_c1 | valid | `b307c41eb54cf452` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=2;consumer=100% | r_mat_n2_c100 | valid | `b307c41eb54cf452` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=1% | r_mat_n4_c1 | valid | `ccfe1af393deae98` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=100% | r_mat_n4_c100 | valid | `ccfe1af393deae98` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=8;consumer=1% | r_mat_n8_c1 | valid | `c8c73778097931dc` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=8;consumer=100% | r_mat_n8_c100 | valid | `c8c73778097931dc` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=notmat;refs=16;consumer=1% | r_notmat_n16_c1 | valid | `1c38b58d8cec08d3` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `9779fc7cd40b6621` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `549a07108ab01ba9` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `2ce8b85a7b60c88e` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=2;consumer=1% | r_notmat_n2_c1 | valid | `28cc913d71d0d72b` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=2;consumer=100% | r_notmat_n2_c100 | valid | `58596ee2bf1338ac` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `9f850e53560baafb` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `5bb21449fb3f5342` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=8;consumer=1% | r_notmat_n8_c1 | valid | `477f7fbaf1c627cb` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=8;consumer=100% | r_notmat_n8_c100 | valid | `1353860480f2b2fb` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, VERBOSE, BUFFERS, FORMAT TEXT` supported on postgres
- one host, single client: no independent hardware, real-network or equal-total-budget comparison
- no concurrent writers: the isolated read contract only; a matched blended read/write phase is a recorded gap
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
