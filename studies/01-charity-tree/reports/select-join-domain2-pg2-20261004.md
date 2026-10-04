# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 91 cells, 91 valid, engine `postgres`, topology `single`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-domain2-pg2-20261004`
- run tag: `run/study-01/select-join-domain2-pg2-20261004`
- environment: `host-zenbook-ux5406sa`
- repo commit: `445fbc1646d58074c45cbaea57a83234132adc34`
- engine version: `PostgreSQL 17.11 (Debian 17.11-1.pgdg13+2) on x86_64-pc-linux-gnu, compiled by gcc (Debian 14.2.0-19) 14.2.0, 64-bit`
- database engine image: `docker.io/library/postgres:17.11` @ `7296f210ae81031ec955dbad9a67a84fe958572a2153b8d0826a647522904dc1`
- benchmark client image: `localhost/charitybench-select:select-join-domain2-pg2-20261004` @ `79bc753ff7e57e102982a323939c43e68e9b3664887683d3a5d5bf072650f4c1`
- inputs digest: `d4a7ba24e4ff2986f6b0640e3a4e72ab1b07969591ea625b3228c29d0aa6edbe`
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
| 20 (0.01%) | 0.4097 | 0.4641 | 0.6065 | 0.5249 | 0.4508 | 0.4473 | 0.5652 | 49.4461 | 0.5000 | 0.4164 |
| 200 (0.1%) | 0.7297 | 0.8040 | 0.9941 | 0.9581 | 0.6885 | 0.9287 | 0.8177 | 51.2854 | 0.8054 | 0.7678 |
| 2000 (1%) | 1.7043 | 1.5673 | 2.0227 | 1.8291 | 1.8953 | 1.5587 | 1.6324 | 51.3152 | 1.5677 | 1.9993 |
| 20000 (10%) | 7.5717 | 7.6304 | 7.1127 | 7.1967 | 9.2963 | 7.7338 | 8.2149 | 55.4266 | 7.0063 | 7.9642 |
| 100000 (50%) | 62.2255 | 55.9779 | 54.7401 | 45.3867 | 61.4459 | 47.1518 | 61.2868 | 68.0637 | 47.7809 | 50.9131 |
| 180000 (90%) | 71.1818 | 77.2965 | 84.9338 | 82.7155 | 80.3233 | 73.6976 | 81.2149 | 84.8802 | 61.8462 | 73.7781 |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | join_on vs join_where | 0.4097 | 0.4641 | 1.1328 | true |
| 20 | filtered_derived vs join_where | 0.4097 | 0.6065 | 1.4804 | true |
| 20 | cte_filtered vs join_where | 0.4097 | 0.5249 | 1.2811 | true |
| 20 | cte_filtered_mat vs join_where | 0.4097 | 0.4508 | 1.1003 | false |
| 20 | cte_filtered_notmat vs join_where | 0.4097 | 0.4473 | 1.0918 | true |
| 20 | cte_joined_outer vs join_where | 0.4097 | 0.5652 | 1.3795 | true |
| 20 | cte_joined_outer_mat vs join_where | 0.4097 | 49.4461 | 120.6871 | false |
| 20 | cte_joined_outer_notmat vs join_where | 0.4097 | 0.5000 | 1.2204 | true |
| 20 | comma_join vs join_where | 0.4097 | 0.4164 | 1.0164 | true |
| 200 | join_on vs join_where | 0.7297 | 0.8040 | 1.1018 | true |
| 200 | filtered_derived vs join_where | 0.7297 | 0.9941 | 1.3624 | true |
| 200 | cte_filtered vs join_where | 0.7297 | 0.9581 | 1.3131 | true |
| 200 | cte_filtered_mat vs join_where | 0.7297 | 0.6885 | 0.9436 | false |
| 200 | cte_filtered_notmat vs join_where | 0.7297 | 0.9287 | 1.2728 | true |
| 200 | cte_joined_outer vs join_where | 0.7297 | 0.8177 | 1.1206 | true |
| 200 | cte_joined_outer_mat vs join_where | 0.7297 | 51.2854 | 70.2849 | false |
| 200 | cte_joined_outer_notmat vs join_where | 0.7297 | 0.8054 | 1.1038 | true |
| 200 | comma_join vs join_where | 0.7297 | 0.7678 | 1.0522 | true |
| 2000 | join_on vs join_where | 1.7043 | 1.5673 | 0.9196 | true |
| 2000 | filtered_derived vs join_where | 1.7043 | 2.0227 | 1.1868 | true |
| 2000 | cte_filtered vs join_where | 1.7043 | 1.8291 | 1.0732 | true |
| 2000 | cte_filtered_mat vs join_where | 1.7043 | 1.8953 | 1.1120 | false |
| 2000 | cte_filtered_notmat vs join_where | 1.7043 | 1.5587 | 0.9145 | true |
| 2000 | cte_joined_outer vs join_where | 1.7043 | 1.6324 | 0.9578 | true |
| 2000 | cte_joined_outer_mat vs join_where | 1.7043 | 51.3152 | 30.1085 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 1.7043 | 1.5677 | 0.9198 | true |
| 2000 | comma_join vs join_where | 1.7043 | 1.9993 | 1.1731 | true |
| 20000 | join_on vs join_where | 7.5717 | 7.6304 | 1.0078 | true |
| 20000 | filtered_derived vs join_where | 7.5717 | 7.1127 | 0.9394 | true |
| 20000 | cte_filtered vs join_where | 7.5717 | 7.1967 | 0.9505 | true |
| 20000 | cte_filtered_mat vs join_where | 7.5717 | 9.2963 | 1.2278 | false |
| 20000 | cte_filtered_notmat vs join_where | 7.5717 | 7.7338 | 1.0214 | true |
| 20000 | cte_joined_outer vs join_where | 7.5717 | 8.2149 | 1.0850 | true |
| 20000 | cte_joined_outer_mat vs join_where | 7.5717 | 55.4266 | 7.3202 | false |
| 20000 | cte_joined_outer_notmat vs join_where | 7.5717 | 7.0063 | 0.9253 | true |
| 20000 | comma_join vs join_where | 7.5717 | 7.9642 | 1.0518 | true |
| 100000 | join_on vs join_where | 62.2255 | 55.9779 | 0.8996 | true |
| 100000 | filtered_derived vs join_where | 62.2255 | 54.7401 | 0.8797 | true |
| 100000 | cte_filtered vs join_where | 62.2255 | 45.3867 | 0.7294 | true |
| 100000 | cte_filtered_mat vs join_where | 62.2255 | 61.4459 | 0.9875 | false |
| 100000 | cte_filtered_notmat vs join_where | 62.2255 | 47.1518 | 0.7578 | true |
| 100000 | cte_joined_outer vs join_where | 62.2255 | 61.2868 | 0.9849 | true |
| 100000 | cte_joined_outer_mat vs join_where | 62.2255 | 68.0637 | 1.0938 | false |
| 100000 | cte_joined_outer_notmat vs join_where | 62.2255 | 47.7809 | 0.7679 | true |
| 100000 | comma_join vs join_where | 62.2255 | 50.9131 | 0.8182 | true |
| 180000 | join_on vs join_where | 71.1818 | 77.2965 | 1.0859 | true |
| 180000 | filtered_derived vs join_where | 71.1818 | 84.9338 | 1.1932 | true |
| 180000 | cte_filtered vs join_where | 71.1818 | 82.7155 | 1.1620 | true |
| 180000 | cte_filtered_mat vs join_where | 71.1818 | 80.3233 | 1.1284 | false |
| 180000 | cte_filtered_notmat vs join_where | 71.1818 | 73.6976 | 1.0353 | true |
| 180000 | cte_joined_outer vs join_where | 71.1818 | 81.2149 | 1.1410 | true |
| 180000 | cte_joined_outer_mat vs join_where | 71.1818 | 84.8802 | 1.1924 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 71.1818 | 61.8462 | 0.8688 | true |
| 180000 | comma_join vs join_where | 71.1818 | 73.7781 | 1.0365 | true |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_filtered_mat vs cte_filtered | 0.5249 | 0.4508 | 0.8588 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_filtered_notmat vs cte_filtered | 0.5249 | 0.4473 | 0.8522 | true | policy=not_materialized; observed_cte_scan=false |
| 20 | cte_joined_outer_mat vs cte_joined_outer | 0.5652 | 49.4461 | 87.4873 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 0.5652 | 0.5000 | 0.8847 | true | policy=not_materialized; observed_cte_scan=false |
| 200 | cte_filtered_mat vs cte_filtered | 0.9581 | 0.6885 | 0.7186 | false | policy=materialized; observed_cte_scan=true |
| 200 | cte_filtered_notmat vs cte_filtered | 0.9581 | 0.9287 | 0.9693 | true | policy=not_materialized; observed_cte_scan=false |
| 200 | cte_joined_outer_mat vs cte_joined_outer | 0.8177 | 51.2854 | 62.7208 | false | policy=materialized; observed_cte_scan=true |
| 200 | cte_joined_outer_notmat vs cte_joined_outer | 0.8177 | 0.8054 | 0.9850 | true | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_filtered_mat vs cte_filtered | 1.8291 | 1.8953 | 1.0362 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_filtered_notmat vs cte_filtered | 1.8291 | 1.5587 | 0.8522 | true | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 1.6324 | 51.3152 | 31.4352 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 1.6324 | 1.5677 | 0.9604 | true | policy=not_materialized; observed_cte_scan=false |
| 20000 | cte_filtered_mat vs cte_filtered | 7.1967 | 9.2963 | 1.2918 | false | policy=materialized; observed_cte_scan=true |
| 20000 | cte_filtered_notmat vs cte_filtered | 7.1967 | 7.7338 | 1.0746 | true | policy=not_materialized; observed_cte_scan=false |
| 20000 | cte_joined_outer_mat vs cte_joined_outer | 8.2149 | 55.4266 | 6.7471 | false | policy=materialized; observed_cte_scan=true |
| 20000 | cte_joined_outer_notmat vs cte_joined_outer | 8.2149 | 7.0063 | 0.8529 | true | policy=not_materialized; observed_cte_scan=false |
| 100000 | cte_filtered_mat vs cte_filtered | 45.3867 | 61.4459 | 1.3538 | false | policy=materialized; observed_cte_scan=true |
| 100000 | cte_filtered_notmat vs cte_filtered | 45.3867 | 47.1518 | 1.0389 | true | policy=not_materialized; observed_cte_scan=false |
| 100000 | cte_joined_outer_mat vs cte_joined_outer | 61.2868 | 68.0637 | 1.1106 | false | policy=materialized; observed_cte_scan=true |
| 100000 | cte_joined_outer_notmat vs cte_joined_outer | 61.2868 | 47.7809 | 0.7796 | true | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_filtered_mat vs cte_filtered | 82.7155 | 80.3233 | 0.9711 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_filtered_notmat vs cte_filtered | 82.7155 | 73.6976 | 0.8910 | true | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 81.2149 | 84.8802 | 1.0451 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 81.2149 | 61.8462 | 0.7615 | true | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 0.4873 | 1 | `549a07108ab01ba9` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 2.3190 | 1 | `2ce8b85a7b60c88e` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | 2.8397 | 2 | `b307c41eb54cf452` | policy=default; observed_cte_scan=true | valid |
| default | 2 | 100% | 2.8106 | 2 | `b307c41eb54cf452` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 1% | 2.7852 | 4 | `ccfe1af393deae98` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 2.7997 | 4 | `ccfe1af393deae98` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | 2.6959 | 8 | `c8c73778097931dc` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 100% | 2.8873 | 8 | `c8c73778097931dc` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 1% | 3.4334 | 16 | `b0c356323d30c425` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 3.8454 | 16 | `b0c356323d30c425` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 2.5048 | 1 | `eaa7d6058aaa5e03` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 2.2675 | 1 | `eaa7d6058aaa5e03` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | 2.7040 | 2 | `b307c41eb54cf452` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 100% | 2.8119 | 2 | `b307c41eb54cf452` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 1% | 2.4500 | 4 | `ccfe1af393deae98` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 3.0233 | 4 | `ccfe1af393deae98` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | 2.7836 | 8 | `c8c73778097931dc` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 100% | 3.3625 | 8 | `c8c73778097931dc` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 1% | 2.9772 | 16 | `b0c356323d30c425` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 3.4599 | 16 | `b0c356323d30c425` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 0.6413 | 1 | `549a07108ab01ba9` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 2.3894 | 1 | `2ce8b85a7b60c88e` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | 0.8477 | 2 | `28cc913d71d0d72b` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 100% | 4.7580 | 2 | `c0d217dc0afc8e9b` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 1% | 1.1600 | 4 | `9f850e53560baafb` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 6.9940 | 4 | `5c23b4d1c6deb269` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | 1.5805 | 8 | `477f7fbaf1c627cb` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 100% | 9.9142 | 8 | `880756aa372bb547` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 1% | 2.7899 | 16 | `1c38b58d8cec08d3` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 15.7597 | 16 | `46c798816d58114e` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 5.1399 | false |
| notmat | 1 | 1% | 1.3161 | true |
| mat | 1 | 100% | 0.9778 | false |
| notmat | 1 | 100% | 1.0303 | true |
| mat | 2 | 1% | 0.9522 | true |
| notmat | 2 | 1% | 0.2985 | false |
| mat | 2 | 100% | 1.0004 | true |
| notmat | 2 | 100% | 1.6929 | false |
| mat | 4 | 1% | 0.8796 | true |
| notmat | 4 | 1% | 0.4165 | false |
| mat | 4 | 100% | 1.0799 | true |
| notmat | 4 | 100% | 2.4981 | false |
| mat | 8 | 1% | 1.0325 | true |
| notmat | 8 | 1% | 0.5863 | false |
| mat | 8 | 100% | 1.1646 | true |
| notmat | 8 | 100% | 3.4337 | false |
| mat | 16 | 1% | 0.8671 | true |
| notmat | 16 | 1% | 0.8126 | false |
| mat | 16 | 100% | 0.8998 | true |
| notmat | 16 | 100% | 4.0984 | false |

## Repeated-design position control

| first median ms | last median ms | ratio (last/first) | same plan |
|---|---|---|---|
| 2.977199 | 3.117931 | 1.0473 | true |

## Latency samples and derived throughput

_The reader is sequential, so `queries/s` is 1000/median_ms, a derived rate and not a concurrent-capacity test. min/max are the five sample extremes._

| axis | status | min ms | median ms | max ms | samples | derived queries/s |
|---|---|---|---|---|---|---|
| family=bag;stmt=comma_join;threshold=100000 | valid | 46.2070 | 50.9131 | 68.5563 | 5 | 19.6 |
| family=bag;stmt=comma_join;threshold=180000 | valid | 51.7546 | 73.7781 | 88.1060 | 5 | 13.6 |
| family=bag;stmt=comma_join;threshold=20 | valid | 0.3993 | 0.4164 | 1.2667 | 5 | 2401.4 |
| family=bag;stmt=comma_join;threshold=200 | valid | 0.7099 | 0.7678 | 1.7660 | 5 | 1302.5 |
| family=bag;stmt=comma_join;threshold=2000 | valid | 1.9261 | 1.9993 | 2.5968 | 5 | 500.2 |
| family=bag;stmt=comma_join;threshold=20000 | valid | 6.5585 | 7.9642 | 8.8978 | 5 | 125.6 |
| family=bag;stmt=cte_filtered;threshold=100000 | valid | 44.2125 | 45.3867 | 52.6520 | 5 | 22.0 |
| family=bag;stmt=cte_filtered;threshold=180000 | valid | 47.5130 | 82.7155 | 91.8056 | 5 | 12.1 |
| family=bag;stmt=cte_filtered;threshold=20 | valid | 0.4567 | 0.5249 | 0.6010 | 5 | 1905.2 |
| family=bag;stmt=cte_filtered;threshold=200 | valid | 0.8924 | 0.9581 | 1.0000 | 5 | 1043.7 |
| family=bag;stmt=cte_filtered;threshold=2000 | valid | 1.4700 | 1.8291 | 2.1347 | 5 | 546.7 |
| family=bag;stmt=cte_filtered;threshold=20000 | valid | 6.4563 | 7.1967 | 9.3344 | 5 | 139.0 |
| family=bag;stmt=cte_filtered_mat;threshold=100000 | valid | 59.3838 | 61.4459 | 72.3610 | 5 | 16.3 |
| family=bag;stmt=cte_filtered_mat;threshold=180000 | valid | 67.1722 | 80.3233 | 91.3873 | 5 | 12.4 |
| family=bag;stmt=cte_filtered_mat;threshold=20 | valid | 0.3928 | 0.4508 | 0.5258 | 5 | 2218.4 |
| family=bag;stmt=cte_filtered_mat;threshold=200 | valid | 0.6615 | 0.6885 | 1.0552 | 5 | 1452.4 |
| family=bag;stmt=cte_filtered_mat;threshold=2000 | valid | 1.5562 | 1.8953 | 39.6440 | 5 | 527.6 |
| family=bag;stmt=cte_filtered_mat;threshold=20000 | valid | 8.3229 | 9.2963 | 9.6243 | 5 | 107.6 |
| family=bag;stmt=cte_filtered_notmat;threshold=100000 | valid | 44.4592 | 47.1518 | 57.3894 | 5 | 21.2 |
| family=bag;stmt=cte_filtered_notmat;threshold=180000 | valid | 47.6788 | 73.6976 | 78.4490 | 5 | 13.6 |
| family=bag;stmt=cte_filtered_notmat;threshold=20 | valid | 0.4135 | 0.4473 | 0.5891 | 5 | 2235.6 |
| family=bag;stmt=cte_filtered_notmat;threshold=200 | valid | 0.8743 | 0.9287 | 1.4797 | 5 | 1076.8 |
| family=bag;stmt=cte_filtered_notmat;threshold=2000 | valid | 1.4759 | 1.5587 | 1.6631 | 5 | 641.6 |
| family=bag;stmt=cte_filtered_notmat;threshold=20000 | valid | 6.4639 | 7.7338 | 48.3131 | 5 | 129.3 |
| family=bag;stmt=cte_joined_outer;threshold=100000 | valid | 53.7916 | 61.2868 | 74.2144 | 5 | 16.3 |
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 77.6482 | 81.2149 | 88.8158 | 5 | 12.3 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 0.4854 | 0.5652 | 20.1520 | 5 | 1769.3 |
| family=bag;stmt=cte_joined_outer;threshold=200 | valid | 0.7276 | 0.8177 | 0.8411 | 5 | 1223.0 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 1.5619 | 1.6324 | 3.2773 | 5 | 612.6 |
| family=bag;stmt=cte_joined_outer;threshold=20000 | valid | 6.4991 | 8.2149 | 41.6199 | 5 | 121.7 |
| family=bag;stmt=cte_joined_outer_mat;threshold=100000 | valid | 61.5153 | 68.0637 | 90.2017 | 5 | 14.7 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 78.9112 | 84.8802 | 93.5733 | 5 | 11.8 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 47.2176 | 49.4461 | 55.6074 | 5 | 20.2 |
| family=bag;stmt=cte_joined_outer_mat;threshold=200 | valid | 44.3059 | 51.2854 | 58.2724 | 5 | 19.5 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 48.1020 | 51.3152 | 55.4391 | 5 | 19.5 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20000 | valid | 50.0235 | 55.4266 | 64.2515 | 5 | 18.0 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=100000 | valid | 46.3470 | 47.7809 | 54.3979 | 5 | 20.9 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 45.2773 | 61.8462 | 148.1527 | 5 | 16.2 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 0.3787 | 0.5000 | 1.0143 | 5 | 2000.0 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=200 | valid | 0.7708 | 0.8054 | 0.8489 | 5 | 1241.6 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 1.3688 | 1.5677 | 1.9911 | 5 | 637.9 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20000 | valid | 6.5149 | 7.0063 | 33.7618 | 5 | 142.7 |
| family=bag;stmt=filtered_derived;threshold=100000 | valid | 44.3076 | 54.7401 | 70.1241 | 5 | 18.3 |
| family=bag;stmt=filtered_derived;threshold=180000 | valid | 42.1700 | 84.9338 | 93.0764 | 5 | 11.8 |
| family=bag;stmt=filtered_derived;threshold=20 | valid | 0.5611 | 0.6065 | 1.3630 | 5 | 1648.7 |
| family=bag;stmt=filtered_derived;threshold=200 | valid | 0.9657 | 0.9941 | 1.0815 | 5 | 1005.9 |
| family=bag;stmt=filtered_derived;threshold=2000 | valid | 1.7706 | 2.0227 | 2.5248 | 5 | 494.4 |
| family=bag;stmt=filtered_derived;threshold=20000 | valid | 6.3740 | 7.1127 | 30.0847 | 5 | 140.6 |
| family=bag;stmt=join_on;threshold=100000 | valid | 46.6845 | 55.9779 | 79.3859 | 5 | 17.9 |
| family=bag;stmt=join_on;threshold=180000 | valid | 50.7479 | 77.2965 | 92.3748 | 5 | 12.9 |
| family=bag;stmt=join_on;threshold=20 | valid | 0.3698 | 0.4641 | 0.6345 | 5 | 2154.6 |
| family=bag;stmt=join_on;threshold=200 | valid | 0.7015 | 0.8040 | 1.2200 | 5 | 1243.9 |
| family=bag;stmt=join_on;threshold=2000 | valid | 1.4223 | 1.5673 | 1.7192 | 5 | 638.0 |
| family=bag;stmt=join_on;threshold=20000 | valid | 6.8167 | 7.6304 | 9.0802 | 5 | 131.1 |
| family=bag;stmt=join_where;threshold=100000 | valid | 46.7761 | 62.2255 | 70.5485 | 5 | 16.1 |
| family=bag;stmt=join_where;threshold=180000 | valid | 46.3008 | 71.1818 | 93.9226 | 5 | 14.0 |
| family=bag;stmt=join_where;threshold=20 | valid | 0.3928 | 0.4097 | 0.5060 | 5 | 2440.8 |
| family=bag;stmt=join_where;threshold=200 | valid | 0.6908 | 0.7297 | 1.0065 | 5 | 1370.5 |
| family=bag;stmt=join_where;threshold=2000 | valid | 1.6614 | 1.7043 | 2.0002 | 5 | 586.7 |
| family=bag;stmt=join_where;threshold=20000 | valid | 7.2603 | 7.5717 | 10.2252 | 5 | 132.1 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 3.2204 | 3.4334 | 3.6687 | 5 | 291.3 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 3.5296 | 3.8454 | 4.6759 | 5 | 260.1 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 0.4508 | 0.4873 | 0.6697 | 5 | 2052.1 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 2.1968 | 2.3190 | 2.6452 | 5 | 431.2 |
| family=reuse;policy=default;refs=2;consumer=1% | valid | 2.5754 | 2.8397 | 2.9639 | 5 | 352.2 |
| family=reuse;policy=default;refs=2;consumer=100% | valid | 2.6365 | 2.8106 | 3.1402 | 5 | 355.8 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 2.6706 | 2.7852 | 2.9772 | 5 | 359.0 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 2.7365 | 2.7997 | 3.0993 | 5 | 357.2 |
| family=reuse;policy=default;refs=8;consumer=1% | valid | 2.5776 | 2.6959 | 3.9699 | 5 | 370.9 |
| family=reuse;policy=default;refs=8;consumer=100% | valid | 2.8066 | 2.8873 | 3.2748 | 5 | 346.3 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 2.8747 | 2.9772 | 3.6084 | 5 | 335.9 |
| family=reuse;policy=mat;refs=16;consumer=1%;position=last | valid | 2.8510 | 3.1179 | 3.3878 | 5 | 320.7 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 3.1952 | 3.4599 | 3.7944 | 5 | 289.0 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 2.2795 | 2.5048 | 2.8616 | 5 | 399.2 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 2.1936 | 2.2675 | 2.8029 | 5 | 441.0 |
| family=reuse;policy=mat;refs=2;consumer=1% | valid | 2.5275 | 2.7040 | 3.0461 | 5 | 369.8 |
| family=reuse;policy=mat;refs=2;consumer=100% | valid | 2.4144 | 2.8119 | 3.4775 | 5 | 355.6 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 2.4042 | 2.4500 | 2.7547 | 5 | 408.2 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 2.7901 | 3.0233 | 3.6433 | 5 | 330.8 |
| family=reuse;policy=mat;refs=8;consumer=1% | valid | 2.5417 | 2.7836 | 2.8251 | 5 | 359.2 |
| family=reuse;policy=mat;refs=8;consumer=100% | valid | 3.0355 | 3.3625 | 3.4344 | 5 | 297.4 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 2.5143 | 2.7899 | 3.8601 | 5 | 358.4 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 14.7529 | 15.7597 | 22.3214 | 5 | 63.5 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 0.6374 | 0.6413 | 0.7795 | 5 | 1559.2 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 2.2477 | 2.3894 | 2.9258 | 5 | 418.5 |
| family=reuse;policy=notmat;refs=2;consumer=1% | valid | 0.7550 | 0.8477 | 0.9931 | 5 | 1179.6 |
| family=reuse;policy=notmat;refs=2;consumer=100% | valid | 4.4109 | 4.7580 | 6.9135 | 5 | 210.2 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 1.0654 | 1.1600 | 1.4326 | 5 | 862.0 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 6.7383 | 6.9940 | 7.4838 | 5 | 143.0 |
| family=reuse;policy=notmat;refs=8;consumer=1% | valid | 1.4773 | 1.5805 | 1.8252 | 5 | 632.7 |
| family=reuse;policy=notmat;refs=8;consumer=100% | valid | 8.7866 | 9.9142 | 12.5430 | 5 | 100.9 |

## Matched blended workload

- phase: `isolated`; target offered rate: 10 updates/s
- offered 0, committed 0, failed 0 over 0.000 s
- offered rate 0.000/s, completed rate 0.000/s
- per key: no keys touched
- ledger: 0 keys reconciled, 0 mismatch

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42678819`, 5 interleaved rounds
  1. family=reuse;policy=mat;refs=16;consumer=1%
  2. family=bag;stmt=cte_joined_outer;threshold=2000
  3. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  4. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  5. family=reuse;policy=notmat;refs=2;consumer=1%
  6. family=bag;stmt=cte_joined_outer;threshold=100000
  7. family=reuse;policy=notmat;refs=8;consumer=1%
  8. family=bag;stmt=cte_joined_outer_notmat;threshold=200
  9. family=bag;stmt=cte_filtered_notmat;threshold=20000
  10. family=bag;stmt=cte_filtered_notmat;threshold=200
  11. family=reuse;policy=mat;refs=2;consumer=100%
  12. family=reuse;policy=mat;refs=8;consumer=1%
  13. family=bag;stmt=join_on;threshold=180000
  14. family=reuse;policy=default;refs=4;consumer=1%
  15. family=bag;stmt=cte_filtered_notmat;threshold=2000
  16. family=bag;stmt=cte_filtered_notmat;threshold=20
  17. family=reuse;policy=default;refs=1;consumer=100%
  18. family=bag;stmt=comma_join;threshold=100000
  19. family=bag;stmt=cte_joined_outer_mat;threshold=200
  20. family=bag;stmt=join_on;threshold=20000
  21. family=bag;stmt=cte_joined_outer_mat;threshold=20000
  22. family=bag;stmt=cte_filtered_notmat;threshold=100000
  23. family=reuse;policy=notmat;refs=4;consumer=1%
  24. family=bag;stmt=cte_filtered_mat;threshold=100000
  25. family=bag;stmt=comma_join;threshold=20000
  26. family=bag;stmt=cte_joined_outer;threshold=20000
  27. family=bag;stmt=cte_joined_outer;threshold=180000
  28. family=reuse;policy=mat;refs=8;consumer=100%
  29. family=bag;stmt=join_where;threshold=180000
  30. family=bag;stmt=filtered_derived;threshold=180000
  31. family=bag;stmt=cte_filtered;threshold=200
  32. family=bag;stmt=cte_filtered;threshold=20000
  33. family=bag;stmt=filtered_derived;threshold=20
  34. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  35. family=bag;stmt=cte_joined_outer_notmat;threshold=100000
  36. family=bag;stmt=cte_filtered;threshold=180000
  37. family=bag;stmt=comma_join;threshold=2000
  38. family=reuse;policy=default;refs=8;consumer=100%
  39. family=bag;stmt=cte_filtered;threshold=20
  40. family=bag;stmt=cte_filtered_mat;threshold=20
  41. family=bag;stmt=cte_joined_outer_mat;threshold=20
  42. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  43. family=reuse;policy=mat;refs=4;consumer=100%
  44. family=bag;stmt=cte_filtered;threshold=2000
  45. family=bag;stmt=cte_joined_outer_notmat;threshold=20000
  46. family=reuse;policy=notmat;refs=16;consumer=1%
  47. family=bag;stmt=comma_join;threshold=180000
  48. family=reuse;policy=notmat;refs=16;consumer=100%
  49. family=bag;stmt=cte_joined_outer;threshold=20
  50. family=reuse;policy=notmat;refs=1;consumer=100%
  51. family=bag;stmt=filtered_derived;threshold=100000
  52. family=reuse;policy=notmat;refs=1;consumer=1%
  53. family=bag;stmt=join_on;threshold=20
  54. family=reuse;policy=notmat;refs=2;consumer=100%
  55. family=bag;stmt=cte_filtered_mat;threshold=200
  56. family=reuse;policy=mat;refs=1;consumer=1%
  57. family=bag;stmt=join_where;threshold=20
  58. family=bag;stmt=cte_joined_outer_mat;threshold=100000
  59. family=bag;stmt=filtered_derived;threshold=2000
  60. family=reuse;policy=mat;refs=1;consumer=100%
  61. family=bag;stmt=cte_filtered_notmat;threshold=180000
  62. family=reuse;policy=default;refs=2;consumer=1%
  63. family=bag;stmt=cte_filtered_mat;threshold=2000
  64. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  65. family=reuse;policy=default;refs=8;consumer=1%
  66. family=bag;stmt=comma_join;threshold=20
  67. family=bag;stmt=cte_filtered_mat;threshold=180000
  68. family=reuse;policy=default;refs=16;consumer=100%
  69. family=reuse;policy=default;refs=1;consumer=1%
  70. family=bag;stmt=join_where;threshold=100000
  71. family=bag;stmt=filtered_derived;threshold=200
  72. family=bag;stmt=join_on;threshold=200
  73. family=bag;stmt=join_where;threshold=2000
  74. family=bag;stmt=cte_filtered_mat;threshold=20000
  75. family=reuse;policy=default;refs=4;consumer=100%
  76. family=bag;stmt=join_where;threshold=200
  77. family=bag;stmt=filtered_derived;threshold=20000
  78. family=reuse;policy=mat;refs=2;consumer=1%
  79. family=reuse;policy=mat;refs=4;consumer=1%
  80. family=reuse;policy=mat;refs=16;consumer=100%
  81. family=bag;stmt=comma_join;threshold=200
  82. family=reuse;policy=notmat;refs=4;consumer=100%
  83. family=bag;stmt=cte_joined_outer;threshold=200
  84. family=reuse;policy=notmat;refs=8;consumer=100%
  85. family=bag;stmt=join_where;threshold=20000
  86. family=bag;stmt=cte_filtered;threshold=100000
  87. family=reuse;policy=default;refs=16;consumer=1%
  88. family=bag;stmt=join_on;threshold=2000
  89. family=bag;stmt=join_on;threshold=100000
  90. family=reuse;policy=default;refs=2;consumer=100%
  91. family=reuse;policy=mat;refs=16;consumer=1%;position=last

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
| family=reuse;policy=mat;refs=16;consumer=1%;position=last | r_mat_n16_c1 | valid | `b0c356323d30c425` | policy=materialized; observed_cte_scan=true |
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
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `46c798816d58114e` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `549a07108ab01ba9` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `2ce8b85a7b60c88e` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=2;consumer=1% | r_notmat_n2_c1 | valid | `28cc913d71d0d72b` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=2;consumer=100% | r_notmat_n2_c100 | valid | `c0d217dc0afc8e9b` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `9f850e53560baafb` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `5c23b4d1c6deb269` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=8;consumer=1% | r_notmat_n8_c1 | valid | `477f7fbaf1c627cb` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=8;consumer=100% | r_notmat_n8_c100 | valid | `880756aa372bb547` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, VERBOSE, BUFFERS, FORMAT TEXT` supported on postgres
- one host, single client: no independent hardware, real-network or equal-total-budget comparison
- no concurrent writers: the isolated read contract only; a matched blended read/write phase is a recorded gap
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
