# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 91 cells, 91 valid, engine `postgres`, topology `single`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-v3-domain-pg-b2-20261005`
- run tag: `run/study-01/select-join-v3-domain-pg-b2-20261005`
- environment: `host-zenbook-ux5406sa`
- repo commit: `0f9d09212c6d52a80d8437d3b5a234f2f711ed36`
- engine version: `PostgreSQL 17.11 (Debian 17.11-1.pgdg13+2) on x86_64-pc-linux-gnu, compiled by gcc (Debian 14.2.0-19) 14.2.0, 64-bit`
- database engine image: `docker.io/library/postgres:17.11` @ `7296f210ae81031ec955dbad9a67a84fe958572a2153b8d0826a647522904dc1`
- benchmark client image: `localhost/charitybench-select:select-join-v3-domain-pg-b2-20261005` @ `8436e36bb66caa8178559d54771b88c063e26848d4ef1d1a56eb524b7f0ecca6`
- inputs digest: `fb97997f4a4e90cc770f64ab502a2b3e56d0d6624632c46cf2ba2ae08e17c116`
- dataset fingerprint: `fcf26e08f377ec51`
- client CPU quota: `100000 100000`, memory: `2147483648`
- setting `client_budget`: 1 CPU 2g
- setting `client_workers`: 1 sequential reader; 0 writers
- setting `connection_pool`: 2 per endpoint = reader/setup + reserve
- setting `customers`: 2000
- setting `dataset`: study01-person-donation-charity
- setting `distributed_rpc`: N/A on PostgreSQL; shared buffers and rows removed parsed per cell where printed
- setting `endpoint_mode`: single endpoint
- setting `endpoints`: pg-single
- setting `fact_count`: 200000
- setting `index_arm`: pg-btree-rank
- setting `index_set`: child rank index and parent-key index, primary keys; no per-spelling index
- setting `node_budget`: 2 CPU 3g
- setting `phase`: blended
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
| pg-single | 921 | 2 |

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
| 20 (0.01%) | 0.4095 | 0.6040 | 0.4189 | 0.6111 | 0.3925 | 0.3996 | 0.4051 | 47.6809 | 0.4651 | 0.4139 |
| 200 (0.1%) | 1.0161 | 0.9778 | 0.8771 | 0.7084 | 0.8697 | 0.7084 | 0.7051 | 41.6973 | 0.7188 | 0.7031 |
| 2000 (1%) | 1.8286 | 1.7785 | 1.8288 | 1.7627 | 1.8512 | 1.5691 | 1.6919 | 48.8511 | 1.7899 | 1.7180 |
| 20000 (10%) | 6.7851 | 7.6238 | 6.4405 | 6.7707 | 9.0050 | 6.6817 | 6.5111 | 52.1481 | 6.8615 | 6.7956 |
| 100000 (50%) | 44.8346 | 43.9233 | 42.9626 | 40.6625 | 58.0561 | 53.6783 | 42.9374 | 64.5782 | 43.5888 | 57.6676 |
| 180000 (90%) | 59.6311 | 74.8731 | 73.9028 | 47.9628 | 71.2179 | 79.9482 | 87.3646 | 87.3942 | 75.5099 | 48.6619 |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | join_on vs join_where | 0.4095 | 0.6040 | 1.4747 | true |
| 20 | filtered_derived vs join_where | 0.4095 | 0.4189 | 1.0228 | true |
| 20 | cte_filtered vs join_where | 0.4095 | 0.6111 | 1.4922 | true |
| 20 | cte_filtered_mat vs join_where | 0.4095 | 0.3925 | 0.9584 | false |
| 20 | cte_filtered_notmat vs join_where | 0.4095 | 0.3996 | 0.9756 | true |
| 20 | cte_joined_outer vs join_where | 0.4095 | 0.4051 | 0.9892 | true |
| 20 | cte_joined_outer_mat vs join_where | 0.4095 | 47.6809 | 116.4236 | false |
| 20 | cte_joined_outer_notmat vs join_where | 0.4095 | 0.4651 | 1.1355 | true |
| 20 | comma_join vs join_where | 0.4095 | 0.4139 | 1.0106 | true |
| 200 | join_on vs join_where | 1.0161 | 0.9778 | 0.9623 | true |
| 200 | filtered_derived vs join_where | 1.0161 | 0.8771 | 0.8632 | true |
| 200 | cte_filtered vs join_where | 1.0161 | 0.7084 | 0.6972 | true |
| 200 | cte_filtered_mat vs join_where | 1.0161 | 0.8697 | 0.8559 | false |
| 200 | cte_filtered_notmat vs join_where | 1.0161 | 0.7084 | 0.6971 | true |
| 200 | cte_joined_outer vs join_where | 1.0161 | 0.7051 | 0.6939 | true |
| 200 | cte_joined_outer_mat vs join_where | 1.0161 | 41.6973 | 41.0352 | false |
| 200 | cte_joined_outer_notmat vs join_where | 1.0161 | 0.7188 | 0.7074 | true |
| 200 | comma_join vs join_where | 1.0161 | 0.7031 | 0.6920 | true |
| 2000 | join_on vs join_where | 1.8286 | 1.7785 | 0.9726 | true |
| 2000 | filtered_derived vs join_where | 1.8286 | 1.8288 | 1.0002 | true |
| 2000 | cte_filtered vs join_where | 1.8286 | 1.7627 | 0.9640 | true |
| 2000 | cte_filtered_mat vs join_where | 1.8286 | 1.8512 | 1.0124 | false |
| 2000 | cte_filtered_notmat vs join_where | 1.8286 | 1.5691 | 0.8581 | true |
| 2000 | cte_joined_outer vs join_where | 1.8286 | 1.6919 | 0.9253 | true |
| 2000 | cte_joined_outer_mat vs join_where | 1.8286 | 48.8511 | 26.7158 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 1.8286 | 1.7899 | 0.9789 | true |
| 2000 | comma_join vs join_where | 1.8286 | 1.7180 | 0.9395 | true |
| 20000 | join_on vs join_where | 6.7851 | 7.6238 | 1.1236 | true |
| 20000 | filtered_derived vs join_where | 6.7851 | 6.4405 | 0.9492 | true |
| 20000 | cte_filtered vs join_where | 6.7851 | 6.7707 | 0.9979 | true |
| 20000 | cte_filtered_mat vs join_where | 6.7851 | 9.0050 | 1.3272 | false |
| 20000 | cte_filtered_notmat vs join_where | 6.7851 | 6.6817 | 0.9848 | true |
| 20000 | cte_joined_outer vs join_where | 6.7851 | 6.5111 | 0.9596 | true |
| 20000 | cte_joined_outer_mat vs join_where | 6.7851 | 52.1481 | 7.6857 | false |
| 20000 | cte_joined_outer_notmat vs join_where | 6.7851 | 6.8615 | 1.0113 | true |
| 20000 | comma_join vs join_where | 6.7851 | 6.7956 | 1.0016 | true |
| 100000 | join_on vs join_where | 44.8346 | 43.9233 | 0.9797 | true |
| 100000 | filtered_derived vs join_where | 44.8346 | 42.9626 | 0.9582 | true |
| 100000 | cte_filtered vs join_where | 44.8346 | 40.6625 | 0.9069 | true |
| 100000 | cte_filtered_mat vs join_where | 44.8346 | 58.0561 | 1.2949 | false |
| 100000 | cte_filtered_notmat vs join_where | 44.8346 | 53.6783 | 1.1973 | true |
| 100000 | cte_joined_outer vs join_where | 44.8346 | 42.9374 | 0.9577 | true |
| 100000 | cte_joined_outer_mat vs join_where | 44.8346 | 64.5782 | 1.4404 | false |
| 100000 | cte_joined_outer_notmat vs join_where | 44.8346 | 43.5888 | 0.9722 | true |
| 100000 | comma_join vs join_where | 44.8346 | 57.6676 | 1.2862 | true |
| 180000 | join_on vs join_where | 59.6311 | 74.8731 | 1.2556 | true |
| 180000 | filtered_derived vs join_where | 59.6311 | 73.9028 | 1.2393 | true |
| 180000 | cte_filtered vs join_where | 59.6311 | 47.9628 | 0.8043 | true |
| 180000 | cte_filtered_mat vs join_where | 59.6311 | 71.2179 | 1.1943 | false |
| 180000 | cte_filtered_notmat vs join_where | 59.6311 | 79.9482 | 1.3407 | true |
| 180000 | cte_joined_outer vs join_where | 59.6311 | 87.3646 | 1.4651 | true |
| 180000 | cte_joined_outer_mat vs join_where | 59.6311 | 87.3942 | 1.4656 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 59.6311 | 75.5099 | 1.2663 | true |
| 180000 | comma_join vs join_where | 59.6311 | 48.6619 | 0.8160 | true |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_filtered_mat vs cte_filtered | 0.6111 | 0.3925 | 0.6423 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_filtered_notmat vs cte_filtered | 0.6111 | 0.3996 | 0.6538 | true | policy=not_materialized; observed_cte_scan=false |
| 20 | cte_joined_outer_mat vs cte_joined_outer | 0.4051 | 47.6809 | 117.6926 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 0.4051 | 0.4651 | 1.1479 | true | policy=not_materialized; observed_cte_scan=false |
| 200 | cte_filtered_mat vs cte_filtered | 0.7084 | 0.8697 | 1.2277 | false | policy=materialized; observed_cte_scan=true |
| 200 | cte_filtered_notmat vs cte_filtered | 0.7084 | 0.7084 | 1.0000 | true | policy=not_materialized; observed_cte_scan=false |
| 200 | cte_joined_outer_mat vs cte_joined_outer | 0.7051 | 41.6973 | 59.1404 | false | policy=materialized; observed_cte_scan=true |
| 200 | cte_joined_outer_notmat vs cte_joined_outer | 0.7051 | 0.7188 | 1.0195 | true | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_filtered_mat vs cte_filtered | 1.7627 | 1.8512 | 1.0502 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_filtered_notmat vs cte_filtered | 1.7627 | 1.5691 | 0.8901 | true | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 1.6919 | 48.8511 | 28.8734 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 1.6919 | 1.7899 | 1.0579 | true | policy=not_materialized; observed_cte_scan=false |
| 20000 | cte_filtered_mat vs cte_filtered | 6.7707 | 9.0050 | 1.3300 | false | policy=materialized; observed_cte_scan=true |
| 20000 | cte_filtered_notmat vs cte_filtered | 6.7707 | 6.6817 | 0.9869 | true | policy=not_materialized; observed_cte_scan=false |
| 20000 | cte_joined_outer_mat vs cte_joined_outer | 6.5111 | 52.1481 | 8.0092 | false | policy=materialized; observed_cte_scan=true |
| 20000 | cte_joined_outer_notmat vs cte_joined_outer | 6.5111 | 6.8615 | 1.0538 | true | policy=not_materialized; observed_cte_scan=false |
| 100000 | cte_filtered_mat vs cte_filtered | 40.6625 | 58.0561 | 1.4278 | false | policy=materialized; observed_cte_scan=true |
| 100000 | cte_filtered_notmat vs cte_filtered | 40.6625 | 53.6783 | 1.3201 | true | policy=not_materialized; observed_cte_scan=false |
| 100000 | cte_joined_outer_mat vs cte_joined_outer | 42.9374 | 64.5782 | 1.5040 | false | policy=materialized; observed_cte_scan=true |
| 100000 | cte_joined_outer_notmat vs cte_joined_outer | 42.9374 | 43.5888 | 1.0152 | true | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_filtered_mat vs cte_filtered | 47.9628 | 71.2179 | 1.4849 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_filtered_notmat vs cte_filtered | 47.9628 | 79.9482 | 1.6669 | true | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 87.3646 | 87.3942 | 1.0003 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 87.3646 | 75.5099 | 0.8643 | true | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 0.4670 | 1 | `549a07108ab01ba9` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 2.2926 | 1 | `2ce8b85a7b60c88e` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | 2.6669 | 2 | `b307c41eb54cf452` | policy=default; observed_cte_scan=true | valid |
| default | 2 | 100% | 2.6740 | 2 | `b307c41eb54cf452` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 1% | 2.3653 | 4 | `ccfe1af393deae98` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 2.6745 | 4 | `ccfe1af393deae98` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | 2.9822 | 8 | `c8c73778097931dc` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 100% | 2.8305 | 8 | `c8c73778097931dc` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 1% | 2.9521 | 16 | `b0c356323d30c425` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 3.3188 | 16 | `b0c356323d30c425` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 2.3967 | 1 | `eaa7d6058aaa5e03` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 2.4664 | 1 | `eaa7d6058aaa5e03` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | 2.6250 | 2 | `b307c41eb54cf452` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 100% | 2.4558 | 2 | `b307c41eb54cf452` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 1% | 2.4017 | 4 | `ccfe1af393deae98` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 2.8122 | 4 | `ccfe1af393deae98` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | 2.9081 | 8 | `c8c73778097931dc` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 100% | 2.7893 | 8 | `c8c73778097931dc` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 1% | 3.3597 | 16 | `b0c356323d30c425` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 3.4995 | 16 | `b0c356323d30c425` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 0.5244 | 1 | `549a07108ab01ba9` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 2.3519 | 1 | `2ce8b85a7b60c88e` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | 0.8489 | 2 | `28cc913d71d0d72b` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 100% | 4.3170 | 2 | `555b22c9849af2d9` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 1% | 0.8963 | 4 | `9f850e53560baafb` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 6.5116 | 4 | `b310752c1df5d8d2` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | 1.3628 | 8 | `477f7fbaf1c627cb` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 100% | 9.1834 | 8 | `232c329e29f3dc44` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 1% | 2.4223 | 16 | `1c38b58d8cec08d3` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 14.7379 | 16 | `e389ae3d205742e6` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 5.1321 | false |
| notmat | 1 | 1% | 1.1229 | true |
| mat | 1 | 100% | 1.0758 | false |
| notmat | 1 | 100% | 1.0258 | true |
| mat | 2 | 1% | 0.9843 | true |
| notmat | 2 | 1% | 0.3183 | false |
| mat | 2 | 100% | 0.9184 | true |
| notmat | 2 | 100% | 1.6144 | false |
| mat | 4 | 1% | 1.0154 | true |
| notmat | 4 | 1% | 0.3789 | false |
| mat | 4 | 100% | 1.0515 | true |
| notmat | 4 | 100% | 2.4347 | false |
| mat | 8 | 1% | 0.9752 | true |
| notmat | 8 | 1% | 0.4570 | false |
| mat | 8 | 100% | 0.9854 | true |
| notmat | 8 | 100% | 3.2444 | false |
| mat | 16 | 1% | 1.1381 | true |
| notmat | 16 | 1% | 0.8206 | false |
| mat | 16 | 100% | 1.0544 | true |
| notmat | 16 | 100% | 4.4407 | false |

## Repeated-design position control

| first median ms | last median ms | ratio (last/first) | same plan |
|---|---|---|---|
| 48.661921 | 53.999749 | 1.1097 | true |

## Physical arm

- `charity_pkey`: `CREATE UNIQUE INDEX charity_pkey ON selectjoin_domain.charity USING btree (charity_id)`
- `donation_person_idx`: `CREATE INDEX donation_person_idx ON selectjoin_domain.donation USING btree (person_id)`
- `donation_pkey`: `CREATE UNIQUE INDEX donation_pkey ON selectjoin_domain.donation USING btree (donation_id)`
- `donation_rank_idx`: `CREATE INDEX donation_rank_idx ON selectjoin_domain.donation USING btree (rank)`
- `person_pkey`: `CREATE UNIQUE INDEX person_pkey ON selectjoin_domain.person USING btree (person_id)`

- index arm `pg-btree-rank`; placement `engine default`; node budget `2 CPU 3g`; client budget `1 CPU 2g`

Physical-arm gate (default join must use the rank index at its three most selective thresholds):

- threshold 20: used
- threshold 200: used
- threshold 2000: used

## Throttling (cgroup `cpu.stat` deltas)

- client: nr_periods 374, nr_throttled 301, throttled_usec 17031223 | quota 100000 100000
, memory 2147483648

- pg-single: nr_periods 377, nr_throttled 0, throttled_usec 0 | quota 200000 100000, memory 3221225472

## Access path per cell

_Engine-reported from each cell's own plan: whether the rank index or a sequential scan was used, rows scanned, PostgreSQL buffers and filter removals, and reader retries of transient conflicts. A cell that scanned the whole table at a selective threshold is a different physical regime, not a faster query._

| axis | status | index scan | seq scan | storage rows scanned | index rows scanned | buffers hit/read | rows removed | retries |
|---|---|---|---|---|---|---|---|---|
| family=bag;stmt=comma_join;threshold=100000 | valid | true | true | 0 | 0 | 100263/0 | 0 | 0 |
| family=bag;stmt=comma_join;threshold=180000 | valid | false | true | 0 | 0 | 1286/0 | 20000 | 0 |
| family=bag;stmt=comma_join;threshold=180000;position=last | valid | false | true | 0 | 0 | 1288/0 | 20000 | 0 |
| family=bag;stmt=comma_join;threshold=20 | valid | true | true | 0 | 0 | 84/0 | 0 | 0 |
| family=bag;stmt=comma_join;threshold=200 | valid | false | true | 0 | 0 | 203/0 | 0 | 0 |
| family=bag;stmt=comma_join;threshold=2000 | valid | false | true | 0 | 0 | 1033/0 | 0 | 0 |
| family=bag;stmt=comma_join;threshold=20000 | valid | false | true | 0 | 0 | 1349/0 | 0 | 0 |
| family=bag;stmt=cte_filtered;threshold=100000 | valid | true | true | 0 | 0 | 100263/0 | 0 | 0 |
| family=bag;stmt=cte_filtered;threshold=180000 | valid | false | true | 0 | 0 | 1286/0 | 20000 | 0 |
| family=bag;stmt=cte_filtered;threshold=20 | valid | true | true | 0 | 0 | 84/0 | 0 | 0 |
| family=bag;stmt=cte_filtered;threshold=200 | valid | false | true | 0 | 0 | 203/0 | 0 | 0 |
| family=bag;stmt=cte_filtered;threshold=2000 | valid | false | true | 0 | 0 | 1033/0 | 0 | 0 |
| family=bag;stmt=cte_filtered;threshold=20000 | valid | false | true | 0 | 0 | 1349/0 | 0 | 0 |
| family=bag;stmt=cte_filtered_mat;threshold=100000 | valid | true | true | 0 | 0 | 100263/0 | 0 | 0 |
| family=bag;stmt=cte_filtered_mat;threshold=180000 | valid | false | true | 0 | 0 | 1286/0 | 20000 | 0 |
| family=bag;stmt=cte_filtered_mat;threshold=20 | valid | true | true | 0 | 0 | 84/0 | 0 | 0 |
| family=bag;stmt=cte_filtered_mat;threshold=200 | valid | false | true | 0 | 0 | 203/0 | 0 | 0 |
| family=bag;stmt=cte_filtered_mat;threshold=2000 | valid | false | true | 0 | 0 | 1033/0 | 0 | 0 |
| family=bag;stmt=cte_filtered_mat;threshold=20000 | valid | false | true | 0 | 0 | 1349/0 | 0 | 0 |
| family=bag;stmt=cte_filtered_notmat;threshold=100000 | valid | true | true | 0 | 0 | 100263/0 | 0 | 0 |
| family=bag;stmt=cte_filtered_notmat;threshold=180000 | valid | false | true | 0 | 0 | 1286/0 | 20000 | 0 |
| family=bag;stmt=cte_filtered_notmat;threshold=20 | valid | true | true | 0 | 0 | 84/0 | 0 | 0 |
| family=bag;stmt=cte_filtered_notmat;threshold=200 | valid | false | true | 0 | 0 | 203/0 | 0 | 0 |
| family=bag;stmt=cte_filtered_notmat;threshold=2000 | valid | false | true | 0 | 0 | 1033/0 | 0 | 0 |
| family=bag;stmt=cte_filtered_notmat;threshold=20000 | valid | false | true | 0 | 0 | 1349/0 | 0 | 0 |
| family=bag;stmt=cte_joined_outer;threshold=100000 | valid | true | true | 0 | 0 | 100263/0 | 0 | 0 |
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | false | true | 0 | 0 | 1286/0 | 20000 | 0 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | true | true | 0 | 0 | 84/0 | 0 | 0 |
| family=bag;stmt=cte_joined_outer;threshold=200 | valid | false | true | 0 | 0 | 203/0 | 0 | 0 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | false | true | 0 | 0 | 1033/0 | 0 | 0 |
| family=bag;stmt=cte_joined_outer;threshold=20000 | valid | false | true | 0 | 0 | 1349/0 | 0 | 0 |
| family=bag;stmt=cte_joined_outer_mat;threshold=100000 | valid | false | true | 0 | 0 | 1286/0 | 100000 | 0 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | false | true | 0 | 0 | 1286/0 | 20000 | 0 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | false | true | 0 | 0 | 1286/0 | 199980 | 0 |
| family=bag;stmt=cte_joined_outer_mat;threshold=200 | valid | false | true | 0 | 0 | 1286/0 | 199800 | 0 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | false | true | 0 | 0 | 1286/0 | 198000 | 0 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20000 | valid | false | true | 0 | 0 | 1286/0 | 180000 | 0 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=100000 | valid | true | true | 0 | 0 | 100263/0 | 0 | 0 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | false | true | 0 | 0 | 1286/0 | 20000 | 0 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | true | true | 0 | 0 | 84/0 | 0 | 0 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=200 | valid | false | true | 0 | 0 | 203/0 | 0 | 0 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | false | true | 0 | 0 | 1033/0 | 0 | 0 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20000 | valid | false | true | 0 | 0 | 1349/0 | 0 | 0 |
| family=bag;stmt=filtered_derived;threshold=100000 | valid | true | true | 0 | 0 | 100263/0 | 0 | 0 |
| family=bag;stmt=filtered_derived;threshold=180000 | valid | false | true | 0 | 0 | 1286/0 | 20000 | 0 |
| family=bag;stmt=filtered_derived;threshold=20 | valid | true | true | 0 | 0 | 84/0 | 0 | 0 |
| family=bag;stmt=filtered_derived;threshold=200 | valid | false | true | 0 | 0 | 203/0 | 0 | 0 |
| family=bag;stmt=filtered_derived;threshold=2000 | valid | false | true | 0 | 0 | 1033/0 | 0 | 0 |
| family=bag;stmt=filtered_derived;threshold=20000 | valid | false | true | 0 | 0 | 1349/0 | 0 | 0 |
| family=bag;stmt=join_on;threshold=100000 | valid | true | true | 0 | 0 | 100263/0 | 0 | 0 |
| family=bag;stmt=join_on;threshold=180000 | valid | false | true | 0 | 0 | 1286/0 | 20000 | 0 |
| family=bag;stmt=join_on;threshold=20 | valid | true | true | 0 | 0 | 84/0 | 0 | 0 |
| family=bag;stmt=join_on;threshold=200 | valid | false | true | 0 | 0 | 203/0 | 0 | 0 |
| family=bag;stmt=join_on;threshold=2000 | valid | false | true | 0 | 0 | 1033/0 | 0 | 0 |
| family=bag;stmt=join_on;threshold=20000 | valid | false | true | 0 | 0 | 1349/0 | 0 | 0 |
| family=bag;stmt=join_where;threshold=100000 | valid | true | true | 0 | 0 | 100263/0 | 0 | 0 |
| family=bag;stmt=join_where;threshold=180000 | valid | false | true | 0 | 0 | 1286/0 | 20000 | 0 |
| family=bag;stmt=join_where;threshold=20 | valid | true | true | 0 | 0 | 84/0 | 0 | 0 |
| family=bag;stmt=join_where;threshold=200 | valid | false | true | 0 | 0 | 203/0 | 0 | 0 |
| family=bag;stmt=join_where;threshold=2000 | valid | false | true | 0 | 0 | 1033/0 | 0 | 0 |
| family=bag;stmt=join_where;threshold=20000 | valid | false | true | 0 | 0 | 1349/0 | 0 | 0 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | false | false | 0 | 0 | 0/0 | 0 | 0 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | false | false | 0 | 0 | 0/0 | 0 | 0 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | true | false | 0 | 0 | 26/0 | 1983 | 0 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | false | false | 0 | 0 | 0/0 | 0 | 0 |
| family=reuse;policy=default;refs=2;consumer=1% | valid | false | false | 0 | 0 | 0/0 | 0 | 0 |
| family=reuse;policy=default;refs=2;consumer=100% | valid | false | false | 0 | 0 | 0/0 | 0 | 0 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | false | false | 0 | 0 | 0/0 | 0 | 0 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | false | false | 0 | 0 | 0/0 | 0 | 0 |
| family=reuse;policy=default;refs=8;consumer=1% | valid | false | false | 0 | 0 | 0/0 | 0 | 0 |
| family=reuse;policy=default;refs=8;consumer=100% | valid | false | false | 0 | 0 | 0/0 | 0 | 0 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | false | false | 0 | 0 | 0/0 | 0 | 0 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | false | false | 0 | 0 | 0/0 | 0 | 0 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | false | false | 0 | 0 | 0/0 | 0 | 0 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | false | false | 0 | 0 | 0/0 | 0 | 0 |
| family=reuse;policy=mat;refs=2;consumer=1% | valid | false | false | 0 | 0 | 0/0 | 0 | 0 |
| family=reuse;policy=mat;refs=2;consumer=100% | valid | false | false | 0 | 0 | 0/0 | 0 | 0 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | false | false | 0 | 0 | 0/0 | 0 | 0 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | false | false | 0 | 0 | 0/0 | 0 | 0 |
| family=reuse;policy=mat;refs=8;consumer=1% | valid | false | false | 0 | 0 | 0/0 | 0 | 0 |
| family=reuse;policy=mat;refs=8;consumer=100% | valid | false | false | 0 | 0 | 0/0 | 0 | 0 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | true | false | 0 | 0 | 480/0 | 1983 | 0 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | false | false | 0 | 0 | 0/0 | 0 | 0 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | true | false | 0 | 0 | 30/0 | 1983 | 0 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | false | false | 0 | 0 | 0/0 | 0 | 0 |
| family=reuse;policy=notmat;refs=2;consumer=1% | valid | true | false | 0 | 0 | 64/0 | 1983 | 0 |
| family=reuse;policy=notmat;refs=2;consumer=100% | valid | false | false | 0 | 0 | 0/0 | 0 | 0 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | true | false | 0 | 0 | 112/0 | 1983 | 0 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | false | false | 0 | 0 | 0/0 | 0 | 0 |
| family=reuse;policy=notmat;refs=8;consumer=1% | valid | true | false | 0 | 0 | 256/0 | 1983 | 0 |
| family=reuse;policy=notmat;refs=8;consumer=100% | valid | false | false | 0 | 0 | 0/0 | 0 | 0 |

## Latency samples and derived throughput

_The reader is sequential, so `queries/s` is 1000/median_ms, a derived rate and not a concurrent-capacity test. min/max are the five sample extremes. Each sample includes any retry and backoff time the statement needed after a transient conflict; per-cell retry counts appear in the access-path table._

| axis | status | min ms | median ms | max ms | samples | derived queries/s |
|---|---|---|---|---|---|---|
| family=bag;stmt=comma_join;threshold=100000 | valid | 41.0847 | 57.6676 | 77.6259 | 5 | 17.3 |
| family=bag;stmt=comma_join;threshold=180000 | valid | 45.2125 | 48.6619 | 52.0028 | 5 | 20.5 |
| family=bag;stmt=comma_join;threshold=180000;position=last | valid | 44.0239 | 53.9997 | 86.8967 | 5 | 18.5 |
| family=bag;stmt=comma_join;threshold=20 | valid | 0.3519 | 0.4139 | 0.4659 | 5 | 2416.1 |
| family=bag;stmt=comma_join;threshold=200 | valid | 0.6704 | 0.7031 | 1.0684 | 5 | 1422.2 |
| family=bag;stmt=comma_join;threshold=2000 | valid | 1.5618 | 1.7180 | 1.9262 | 5 | 582.1 |
| family=bag;stmt=comma_join;threshold=20000 | valid | 6.3922 | 6.7956 | 7.4831 | 5 | 147.2 |
| family=bag;stmt=cte_filtered;threshold=100000 | valid | 39.5666 | 40.6625 | 45.6924 | 5 | 24.6 |
| family=bag;stmt=cte_filtered;threshold=180000 | valid | 45.7661 | 47.9628 | 50.0339 | 5 | 20.8 |
| family=bag;stmt=cte_filtered;threshold=20 | valid | 0.5671 | 0.6111 | 0.7936 | 5 | 1636.3 |
| family=bag;stmt=cte_filtered;threshold=200 | valid | 0.6770 | 0.7084 | 0.7894 | 5 | 1411.6 |
| family=bag;stmt=cte_filtered;threshold=2000 | valid | 1.6307 | 1.7627 | 2.4020 | 5 | 567.3 |
| family=bag;stmt=cte_filtered;threshold=20000 | valid | 6.4106 | 6.7707 | 7.8432 | 5 | 147.7 |
| family=bag;stmt=cte_filtered_mat;threshold=100000 | valid | 56.5851 | 58.0561 | 59.1598 | 5 | 17.2 |
| family=bag;stmt=cte_filtered_mat;threshold=180000 | valid | 60.6078 | 71.2179 | 79.9559 | 5 | 14.0 |
| family=bag;stmt=cte_filtered_mat;threshold=20 | valid | 0.3473 | 0.3925 | 0.4702 | 5 | 2547.6 |
| family=bag;stmt=cte_filtered_mat;threshold=200 | valid | 0.7921 | 0.8697 | 1.0268 | 5 | 1149.8 |
| family=bag;stmt=cte_filtered_mat;threshold=2000 | valid | 1.7077 | 1.8512 | 1.9427 | 5 | 540.2 |
| family=bag;stmt=cte_filtered_mat;threshold=20000 | valid | 8.2361 | 9.0050 | 40.5114 | 5 | 111.0 |
| family=bag;stmt=cte_filtered_notmat;threshold=100000 | valid | 41.7921 | 53.6783 | 64.8539 | 5 | 18.6 |
| family=bag;stmt=cte_filtered_notmat;threshold=180000 | valid | 45.5570 | 79.9482 | 92.9617 | 5 | 12.5 |
| family=bag;stmt=cte_filtered_notmat;threshold=20 | valid | 0.3936 | 0.3996 | 0.4242 | 5 | 2502.7 |
| family=bag;stmt=cte_filtered_notmat;threshold=200 | valid | 0.6775 | 0.7084 | 29.8199 | 5 | 1411.6 |
| family=bag;stmt=cte_filtered_notmat;threshold=2000 | valid | 1.5266 | 1.5691 | 2.0999 | 5 | 637.3 |
| family=bag;stmt=cte_filtered_notmat;threshold=20000 | valid | 6.3095 | 6.6817 | 6.8158 | 5 | 149.7 |
| family=bag;stmt=cte_joined_outer;threshold=100000 | valid | 40.7872 | 42.9374 | 59.3952 | 5 | 23.3 |
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 74.1437 | 87.3646 | 91.8399 | 5 | 11.4 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 0.3753 | 0.4051 | 0.4275 | 5 | 2468.3 |
| family=bag;stmt=cte_joined_outer;threshold=200 | valid | 0.6193 | 0.7051 | 28.7968 | 5 | 1418.3 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 1.6522 | 1.6919 | 1.9961 | 5 | 591.0 |
| family=bag;stmt=cte_joined_outer;threshold=20000 | valid | 6.4134 | 6.5111 | 6.8985 | 5 | 153.6 |
| family=bag;stmt=cte_joined_outer_mat;threshold=100000 | valid | 56.0676 | 64.5782 | 75.1135 | 5 | 15.5 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 71.9982 | 87.3942 | 97.2767 | 5 | 11.4 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 41.3921 | 47.6809 | 48.4566 | 5 | 21.0 |
| family=bag;stmt=cte_joined_outer_mat;threshold=200 | valid | 40.0911 | 41.6973 | 49.4877 | 5 | 24.0 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 42.9124 | 48.8511 | 49.1291 | 5 | 20.5 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20000 | valid | 50.7150 | 52.1481 | 53.5374 | 5 | 19.2 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=100000 | valid | 42.7732 | 43.5888 | 50.9369 | 5 | 22.9 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 43.5227 | 75.5099 | 89.1625 | 5 | 13.2 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 0.4037 | 0.4651 | 0.6055 | 5 | 2150.3 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=200 | valid | 0.6822 | 0.7188 | 0.8061 | 5 | 1391.2 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 1.6697 | 1.7899 | 1.8458 | 5 | 558.7 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20000 | valid | 6.6208 | 6.8615 | 7.1488 | 5 | 145.7 |
| family=bag;stmt=filtered_derived;threshold=100000 | valid | 40.3688 | 42.9626 | 46.4995 | 5 | 23.3 |
| family=bag;stmt=filtered_derived;threshold=180000 | valid | 47.3767 | 73.9028 | 95.6383 | 5 | 13.5 |
| family=bag;stmt=filtered_derived;threshold=20 | valid | 0.3901 | 0.4189 | 0.4658 | 5 | 2387.3 |
| family=bag;stmt=filtered_derived;threshold=200 | valid | 0.8047 | 0.8771 | 1.2295 | 5 | 1140.1 |
| family=bag;stmt=filtered_derived;threshold=2000 | valid | 1.7544 | 1.8288 | 1.9084 | 5 | 546.8 |
| family=bag;stmt=filtered_derived;threshold=20000 | valid | 6.3129 | 6.4405 | 6.7317 | 5 | 155.3 |
| family=bag;stmt=join_on;threshold=100000 | valid | 42.4470 | 43.9233 | 66.0942 | 5 | 22.8 |
| family=bag;stmt=join_on;threshold=180000 | valid | 67.4030 | 74.8731 | 82.6633 | 5 | 13.4 |
| family=bag;stmt=join_on;threshold=20 | valid | 0.5738 | 0.6040 | 0.7347 | 5 | 1655.8 |
| family=bag;stmt=join_on;threshold=200 | valid | 0.9074 | 0.9778 | 1.1067 | 5 | 1022.7 |
| family=bag;stmt=join_on;threshold=2000 | valid | 1.7096 | 1.7785 | 2.6712 | 5 | 562.3 |
| family=bag;stmt=join_on;threshold=20000 | valid | 6.7470 | 7.6238 | 9.6829 | 5 | 131.2 |
| family=bag;stmt=join_where;threshold=100000 | valid | 40.7388 | 44.8346 | 53.6539 | 5 | 22.3 |
| family=bag;stmt=join_where;threshold=180000 | valid | 46.3909 | 59.6311 | 84.1494 | 5 | 16.8 |
| family=bag;stmt=join_where;threshold=20 | valid | 0.3980 | 0.4095 | 0.5083 | 5 | 2441.7 |
| family=bag;stmt=join_where;threshold=200 | valid | 0.8790 | 1.0161 | 1.0813 | 5 | 984.1 |
| family=bag;stmt=join_where;threshold=2000 | valid | 1.6648 | 1.8286 | 2.0952 | 5 | 546.9 |
| family=bag;stmt=join_where;threshold=20000 | valid | 6.5584 | 6.7851 | 7.1320 | 5 | 147.4 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 2.8790 | 2.9521 | 3.2337 | 5 | 338.7 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 3.2199 | 3.3188 | 3.7790 | 5 | 301.3 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 0.4421 | 0.4670 | 0.5248 | 5 | 2141.3 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 2.0784 | 2.2926 | 2.4257 | 5 | 436.2 |
| family=reuse;policy=default;refs=2;consumer=1% | valid | 2.5462 | 2.6669 | 2.8134 | 5 | 375.0 |
| family=reuse;policy=default;refs=2;consumer=100% | valid | 2.5484 | 2.6740 | 2.8808 | 5 | 374.0 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 2.2979 | 2.3653 | 2.4766 | 5 | 422.8 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 2.6251 | 2.6745 | 3.0771 | 5 | 373.9 |
| family=reuse;policy=default;refs=8;consumer=1% | valid | 2.7936 | 2.9822 | 3.9540 | 5 | 335.3 |
| family=reuse;policy=default;refs=8;consumer=100% | valid | 2.6767 | 2.8305 | 2.8957 | 5 | 353.3 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 3.1887 | 3.3597 | 3.4555 | 5 | 297.6 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 3.3910 | 3.4995 | 3.5793 | 5 | 285.8 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 2.3272 | 2.3967 | 2.4697 | 5 | 417.2 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 2.3595 | 2.4664 | 2.5492 | 5 | 405.5 |
| family=reuse;policy=mat;refs=2;consumer=1% | valid | 2.5473 | 2.6250 | 3.0867 | 5 | 381.0 |
| family=reuse;policy=mat;refs=2;consumer=100% | valid | 2.3633 | 2.4558 | 2.4598 | 5 | 407.2 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 2.2619 | 2.4017 | 2.5818 | 5 | 416.4 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 2.7492 | 2.8122 | 2.9807 | 5 | 355.6 |
| family=reuse;policy=mat;refs=8;consumer=1% | valid | 2.8333 | 2.9081 | 2.9474 | 5 | 343.9 |
| family=reuse;policy=mat;refs=8;consumer=100% | valid | 2.6940 | 2.7893 | 3.0615 | 5 | 358.5 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 2.2153 | 2.4223 | 2.6767 | 5 | 412.8 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 14.2182 | 14.7379 | 15.1827 | 5 | 67.9 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 0.4469 | 0.5244 | 0.6162 | 5 | 1906.9 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 2.2443 | 2.3519 | 2.4178 | 5 | 425.2 |
| family=reuse;policy=notmat;refs=2;consumer=1% | valid | 0.7934 | 0.8489 | 0.8713 | 5 | 1178.1 |
| family=reuse;policy=notmat;refs=2;consumer=100% | valid | 3.9986 | 4.3170 | 4.7527 | 5 | 231.6 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 0.8175 | 0.8963 | 0.9323 | 5 | 1115.7 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 6.3410 | 6.5116 | 7.3228 | 5 | 153.6 |
| family=reuse;policy=notmat;refs=8;consumer=1% | valid | 1.3150 | 1.3628 | 1.3793 | 5 | 733.8 |
| family=reuse;policy=notmat;refs=8;consumer=100% | valid | 8.6154 | 9.1834 | 10.4971 | 5 | 108.9 |

## Matched blended workload

- phase: `blended`; target offered rate: 10 updates/s
- offered 344, committed 344, failed 0 over 34.398 s
- offered rate 10.000/s, completed rate 10.000/s
- writer operations per endpoint: [344]
- reader operations overlapping a committed write: 129 of 637
- reader retries of transient conflicts: 0
- per key: 344 keys, min 1, max 1 per key
- ledger: 344 keys reconciled, 0 mismatch
- update latency: avg 2.117 ms, min 0.787, max 52.418 over 344

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42907739`, 5 interleaved rounds
  1. family=bag;stmt=comma_join;threshold=180000
  2. family=reuse;policy=default;refs=2;consumer=100%
  3. family=bag;stmt=cte_filtered_notmat;threshold=2000
  4. family=bag;stmt=cte_filtered_mat;threshold=20000
  5. family=bag;stmt=cte_filtered_notmat;threshold=20000
  6. family=bag;stmt=cte_joined_outer_notmat;threshold=20000
  7. family=bag;stmt=cte_filtered_mat;threshold=180000
  8. family=bag;stmt=join_on;threshold=200
  9. family=reuse;policy=default;refs=4;consumer=100%
  10. family=bag;stmt=cte_joined_outer_mat;threshold=200
  11. family=reuse;policy=notmat;refs=1;consumer=100%
  12. family=bag;stmt=comma_join;threshold=2000
  13. family=bag;stmt=cte_filtered_mat;threshold=100000
  14. family=bag;stmt=cte_filtered_notmat;threshold=100000
  15. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  16. family=bag;stmt=comma_join;threshold=100000
  17. family=bag;stmt=join_on;threshold=20
  18. family=reuse;policy=mat;refs=1;consumer=1%
  19. family=bag;stmt=cte_joined_outer_mat;threshold=20000
  20. family=reuse;policy=mat;refs=8;consumer=1%
  21. family=reuse;policy=default;refs=1;consumer=1%
  22. family=reuse;policy=mat;refs=4;consumer=1%
  23. family=bag;stmt=cte_filtered_notmat;threshold=20
  24. family=reuse;policy=mat;refs=8;consumer=100%
  25. family=bag;stmt=cte_filtered_mat;threshold=20
  26. family=bag;stmt=cte_filtered;threshold=180000
  27. family=reuse;policy=mat;refs=4;consumer=100%
  28. family=bag;stmt=filtered_derived;threshold=20000
  29. family=reuse;policy=default;refs=1;consumer=100%
  30. family=bag;stmt=join_on;threshold=180000
  31. family=reuse;policy=mat;refs=16;consumer=1%
  32. family=bag;stmt=join_where;threshold=100000
  33. family=bag;stmt=join_on;threshold=100000
  34. family=bag;stmt=cte_joined_outer_mat;threshold=100000
  35. family=reuse;policy=default;refs=8;consumer=1%
  36. family=reuse;policy=notmat;refs=4;consumer=1%
  37. family=bag;stmt=cte_joined_outer;threshold=100000
  38. family=bag;stmt=cte_filtered;threshold=20
  39. family=bag;stmt=cte_joined_outer_notmat;threshold=200
  40. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  41. family=bag;stmt=filtered_derived;threshold=180000
  42. family=bag;stmt=cte_joined_outer;threshold=180000
  43. family=bag;stmt=cte_filtered_mat;threshold=200
  44. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  45. family=bag;stmt=comma_join;threshold=20
  46. family=reuse;policy=notmat;refs=16;consumer=1%
  47. family=bag;stmt=comma_join;threshold=20000
  48. family=bag;stmt=join_where;threshold=2000
  49. family=bag;stmt=cte_joined_outer;threshold=200
  50. family=bag;stmt=cte_joined_outer;threshold=20000
  51. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  52. family=reuse;policy=default;refs=8;consumer=100%
  53. family=bag;stmt=join_where;threshold=180000
  54. family=reuse;policy=mat;refs=2;consumer=1%
  55. family=reuse;policy=default;refs=16;consumer=100%
  56. family=bag;stmt=cte_filtered_notmat;threshold=200
  57. family=reuse;policy=notmat;refs=1;consumer=1%
  58. family=bag;stmt=filtered_derived;threshold=20
  59. family=bag;stmt=cte_filtered;threshold=200
  60. family=bag;stmt=cte_joined_outer;threshold=20
  61. family=reuse;policy=mat;refs=2;consumer=100%
  62. family=bag;stmt=cte_filtered;threshold=20000
  63. family=bag;stmt=cte_joined_outer_notmat;threshold=100000
  64. family=bag;stmt=filtered_derived;threshold=2000
  65. family=bag;stmt=join_where;threshold=20
  66. family=bag;stmt=cte_filtered_notmat;threshold=180000
  67. family=bag;stmt=join_where;threshold=200
  68. family=bag;stmt=cte_joined_outer;threshold=2000
  69. family=bag;stmt=cte_filtered;threshold=100000
  70. family=reuse;policy=default;refs=2;consumer=1%
  71. family=bag;stmt=filtered_derived;threshold=100000
  72. family=reuse;policy=notmat;refs=2;consumer=1%
  73. family=bag;stmt=comma_join;threshold=200
  74. family=reuse;policy=notmat;refs=8;consumer=100%
  75. family=reuse;policy=notmat;refs=4;consumer=100%
  76. family=bag;stmt=cte_filtered;threshold=2000
  77. family=bag;stmt=join_where;threshold=20000
  78. family=bag;stmt=join_on;threshold=2000
  79. family=bag;stmt=join_on;threshold=20000
  80. family=bag;stmt=filtered_derived;threshold=200
  81. family=bag;stmt=cte_joined_outer_mat;threshold=20
  82. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  83. family=reuse;policy=mat;refs=1;consumer=100%
  84. family=reuse;policy=notmat;refs=2;consumer=100%
  85. family=reuse;policy=default;refs=16;consumer=1%
  86. family=reuse;policy=notmat;refs=8;consumer=1%
  87. family=reuse;policy=default;refs=4;consumer=1%
  88. family=reuse;policy=notmat;refs=16;consumer=100%
  89. family=reuse;policy=mat;refs=16;consumer=100%
  90. family=bag;stmt=cte_filtered_mat;threshold=2000
  91. family=bag;stmt=comma_join;threshold=180000;position=last

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
| family=bag;stmt=comma_join;threshold=180000;position=last | j06_comma_join | valid | `18e11c5d2e685630` |  |
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
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `e389ae3d205742e6` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `549a07108ab01ba9` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `2ce8b85a7b60c88e` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=2;consumer=1% | r_notmat_n2_c1 | valid | `28cc913d71d0d72b` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=2;consumer=100% | r_notmat_n2_c100 | valid | `555b22c9849af2d9` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `9f850e53560baafb` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `b310752c1df5d8d2` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=8;consumer=1% | r_notmat_n8_c1 | valid | `477f7fbaf1c627cb` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=8;consumer=100% | r_notmat_n8_c100 | valid | `232c329e29f3dc44` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, VERBOSE, BUFFERS, FORMAT TEXT` supported on postgres
- one host: this run used node budget 2 CPU 3g and client budget 1 CPU 2g, so budget comparisons are separate classes here, not independent hardware
- one client process with one fixed maintenance writer at 10 updates/s; readers are sequential; the reported sample latency includes retry and backoff time
- physical rank-index arm pg-btree-rank; an engine whose planner ignores the index is a different physical design, not a faster one
- distributed RPC counters are N/A on single-node PostgreSQL; shared buffers and rows removed are parsed per cell where printed
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
