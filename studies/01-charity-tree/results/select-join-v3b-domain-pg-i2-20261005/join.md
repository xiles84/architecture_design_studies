# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 91 cells, 91 valid, engine `postgres`, topology `single`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-v3b-domain-pg-i2-20261005`
- run tag: `run/study-01/select-join-v3b-domain-pg-i2-20261005`
- environment: `host-zenbook-ux5406sa`
- repo commit: `7243fee074bfb16fc5ad8b6afc1a2d7549080935`
- engine version: `PostgreSQL 17.11 (Debian 17.11-1.pgdg13+2) on x86_64-pc-linux-gnu, compiled by gcc (Debian 14.2.0-19) 14.2.0, 64-bit`
- database engine image: `docker.io/library/postgres:17.11` @ `7296f210ae81031ec955dbad9a67a84fe958572a2153b8d0826a647522904dc1`
- benchmark client image: `localhost/charitybench-select:select-join-v3b-domain-pg-i2-20261005` @ `f3f1e10c44cbcb4c7b5bc56ed95a4f071682b6f7d33b5e0a7b688c4271eeef88`
- inputs digest: `093d13091d8585b37c203fd328653867313ca3d810fe35372267895672eaae8f`
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
| pg-single | 921 | 1 |

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
| 20 (0.01%) | 0.4308 | 0.7237 | 0.7151 | 0.4773 | 0.4179 | 0.4587 | 0.4778 | 48.1081 | 0.4173 | 0.5237 |
| 200 (0.1%) | 1.0708 | 0.9253 | 0.9991 | 0.7325 | 0.7822 | 1.0393 | 0.8495 | 48.7891 | 1.0311 | 0.7080 |
| 2000 (1%) | 2.1664 | 1.9979 | 1.7610 | 1.7271 | 2.1855 | 2.1203 | 2.0023 | 47.9916 | 2.2780 | 2.1973 |
| 20000 (10%) | 7.6453 | 7.5124 | 7.6786 | 7.3109 | 9.6622 | 7.5541 | 7.4931 | 53.8350 | 7.7251 | 7.5892 |
| 100000 (50%) | 49.6616 | 48.9945 | 45.2736 | 44.6471 | 60.4186 | 48.1178 | 46.0983 | 67.0791 | 46.4157 | 48.4774 |
| 180000 (90%) | 79.8609 | 77.7107 | 79.8069 | 78.3874 | 80.3723 | 75.4187 | 64.4969 | 78.6928 | 89.8642 | 72.7510 |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | join_on vs join_where | 0.4308 | 0.7237 | 1.6797 | true |
| 20 | filtered_derived vs join_where | 0.4308 | 0.7151 | 1.6597 | true |
| 20 | cte_filtered vs join_where | 0.4308 | 0.4773 | 1.1079 | true |
| 20 | cte_filtered_mat vs join_where | 0.4308 | 0.4179 | 0.9699 | false |
| 20 | cte_filtered_notmat vs join_where | 0.4308 | 0.4587 | 1.0647 | true |
| 20 | cte_joined_outer vs join_where | 0.4308 | 0.4778 | 1.1089 | true |
| 20 | cte_joined_outer_mat vs join_where | 0.4308 | 48.1081 | 111.6643 | false |
| 20 | cte_joined_outer_notmat vs join_where | 0.4308 | 0.4173 | 0.9686 | true |
| 20 | comma_join vs join_where | 0.4308 | 0.5237 | 1.2156 | true |
| 200 | join_on vs join_where | 1.0708 | 0.9253 | 0.8642 | true |
| 200 | filtered_derived vs join_where | 1.0708 | 0.9991 | 0.9331 | true |
| 200 | cte_filtered vs join_where | 1.0708 | 0.7325 | 0.6840 | true |
| 200 | cte_filtered_mat vs join_where | 1.0708 | 0.7822 | 0.7305 | false |
| 200 | cte_filtered_notmat vs join_where | 1.0708 | 1.0393 | 0.9706 | true |
| 200 | cte_joined_outer vs join_where | 1.0708 | 0.8495 | 0.7934 | true |
| 200 | cte_joined_outer_mat vs join_where | 1.0708 | 48.7891 | 45.5643 | false |
| 200 | cte_joined_outer_notmat vs join_where | 1.0708 | 1.0311 | 0.9630 | true |
| 200 | comma_join vs join_where | 1.0708 | 0.7080 | 0.6612 | true |
| 2000 | join_on vs join_where | 2.1664 | 1.9979 | 0.9222 | true |
| 2000 | filtered_derived vs join_where | 2.1664 | 1.7610 | 0.8129 | true |
| 2000 | cte_filtered vs join_where | 2.1664 | 1.7271 | 0.7972 | true |
| 2000 | cte_filtered_mat vs join_where | 2.1664 | 2.1855 | 1.0088 | false |
| 2000 | cte_filtered_notmat vs join_where | 2.1664 | 2.1203 | 0.9787 | true |
| 2000 | cte_joined_outer vs join_where | 2.1664 | 2.0023 | 0.9243 | true |
| 2000 | cte_joined_outer_mat vs join_where | 2.1664 | 47.9916 | 22.1530 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 2.1664 | 2.2780 | 1.0515 | true |
| 2000 | comma_join vs join_where | 2.1664 | 2.1973 | 1.0143 | true |
| 20000 | join_on vs join_where | 7.6453 | 7.5124 | 0.9826 | true |
| 20000 | filtered_derived vs join_where | 7.6453 | 7.6786 | 1.0044 | true |
| 20000 | cte_filtered vs join_where | 7.6453 | 7.3109 | 0.9563 | true |
| 20000 | cte_filtered_mat vs join_where | 7.6453 | 9.6622 | 1.2638 | false |
| 20000 | cte_filtered_notmat vs join_where | 7.6453 | 7.5541 | 0.9881 | true |
| 20000 | cte_joined_outer vs join_where | 7.6453 | 7.4931 | 0.9801 | true |
| 20000 | cte_joined_outer_mat vs join_where | 7.6453 | 53.8350 | 7.0415 | false |
| 20000 | cte_joined_outer_notmat vs join_where | 7.6453 | 7.7251 | 1.0104 | true |
| 20000 | comma_join vs join_where | 7.6453 | 7.5892 | 0.9927 | true |
| 100000 | join_on vs join_where | 49.6616 | 48.9945 | 0.9866 | true |
| 100000 | filtered_derived vs join_where | 49.6616 | 45.2736 | 0.9116 | true |
| 100000 | cte_filtered vs join_where | 49.6616 | 44.6471 | 0.8990 | true |
| 100000 | cte_filtered_mat vs join_where | 49.6616 | 60.4186 | 1.2166 | false |
| 100000 | cte_filtered_notmat vs join_where | 49.6616 | 48.1178 | 0.9689 | true |
| 100000 | cte_joined_outer vs join_where | 49.6616 | 46.0983 | 0.9282 | true |
| 100000 | cte_joined_outer_mat vs join_where | 49.6616 | 67.0791 | 1.3507 | false |
| 100000 | cte_joined_outer_notmat vs join_where | 49.6616 | 46.4157 | 0.9346 | true |
| 100000 | comma_join vs join_where | 49.6616 | 48.4774 | 0.9762 | true |
| 180000 | join_on vs join_where | 79.8609 | 77.7107 | 0.9731 | true |
| 180000 | filtered_derived vs join_where | 79.8609 | 79.8069 | 0.9993 | true |
| 180000 | cte_filtered vs join_where | 79.8609 | 78.3874 | 0.9815 | true |
| 180000 | cte_filtered_mat vs join_where | 79.8609 | 80.3723 | 1.0064 | false |
| 180000 | cte_filtered_notmat vs join_where | 79.8609 | 75.4187 | 0.9444 | true |
| 180000 | cte_joined_outer vs join_where | 79.8609 | 64.4969 | 0.8076 | true |
| 180000 | cte_joined_outer_mat vs join_where | 79.8609 | 78.6928 | 0.9854 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 79.8609 | 89.8642 | 1.1253 | true |
| 180000 | comma_join vs join_where | 79.8609 | 72.7510 | 0.9110 | true |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_filtered_mat vs cte_filtered | 0.4773 | 0.4179 | 0.8754 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_filtered_notmat vs cte_filtered | 0.4773 | 0.4587 | 0.9610 | true | policy=not_materialized; observed_cte_scan=false |
| 20 | cte_joined_outer_mat vs cte_joined_outer | 0.4778 | 48.1081 | 100.6938 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 0.4778 | 0.4173 | 0.8734 | true | policy=not_materialized; observed_cte_scan=false |
| 200 | cte_filtered_mat vs cte_filtered | 0.7325 | 0.7822 | 1.0679 | false | policy=materialized; observed_cte_scan=true |
| 200 | cte_filtered_notmat vs cte_filtered | 0.7325 | 1.0393 | 1.4189 | true | policy=not_materialized; observed_cte_scan=false |
| 200 | cte_joined_outer_mat vs cte_joined_outer | 0.8495 | 48.7891 | 57.4310 | false | policy=materialized; observed_cte_scan=true |
| 200 | cte_joined_outer_notmat vs cte_joined_outer | 0.8495 | 1.0311 | 1.2138 | true | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_filtered_mat vs cte_filtered | 1.7271 | 2.1855 | 1.2654 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_filtered_notmat vs cte_filtered | 1.7271 | 2.1203 | 1.2276 | true | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 2.0023 | 47.9916 | 23.9680 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 2.0023 | 2.2780 | 1.1377 | true | policy=not_materialized; observed_cte_scan=false |
| 20000 | cte_filtered_mat vs cte_filtered | 7.3109 | 9.6622 | 1.3216 | false | policy=materialized; observed_cte_scan=true |
| 20000 | cte_filtered_notmat vs cte_filtered | 7.3109 | 7.5541 | 1.0333 | true | policy=not_materialized; observed_cte_scan=false |
| 20000 | cte_joined_outer_mat vs cte_joined_outer | 7.4931 | 53.8350 | 7.1846 | false | policy=materialized; observed_cte_scan=true |
| 20000 | cte_joined_outer_notmat vs cte_joined_outer | 7.4931 | 7.7251 | 1.0310 | true | policy=not_materialized; observed_cte_scan=false |
| 100000 | cte_filtered_mat vs cte_filtered | 44.6471 | 60.4186 | 1.3532 | false | policy=materialized; observed_cte_scan=true |
| 100000 | cte_filtered_notmat vs cte_filtered | 44.6471 | 48.1178 | 1.0777 | true | policy=not_materialized; observed_cte_scan=false |
| 100000 | cte_joined_outer_mat vs cte_joined_outer | 46.0983 | 67.0791 | 1.4551 | false | policy=materialized; observed_cte_scan=true |
| 100000 | cte_joined_outer_notmat vs cte_joined_outer | 46.0983 | 46.4157 | 1.0069 | true | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_filtered_mat vs cte_filtered | 78.3874 | 80.3723 | 1.0253 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_filtered_notmat vs cte_filtered | 78.3874 | 75.4187 | 0.9621 | true | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 64.4969 | 78.6928 | 1.2201 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 64.4969 | 89.8642 | 1.3933 | true | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 0.5191 | 1 | `549a07108ab01ba9` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 2.5270 | 1 | `2ce8b85a7b60c88e` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | 2.4462 | 2 | `b307c41eb54cf452` | policy=default; observed_cte_scan=true | valid |
| default | 2 | 100% | 2.5693 | 2 | `b307c41eb54cf452` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 1% | 3.1445 | 4 | `ccfe1af393deae98` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 3.2507 | 4 | `ccfe1af393deae98` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | 2.9094 | 8 | `c8c73778097931dc` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 100% | 3.5924 | 8 | `c8c73778097931dc` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 1% | 3.4726 | 16 | `b0c356323d30c425` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 3.8954 | 16 | `b0c356323d30c425` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 2.6265 | 1 | `eaa7d6058aaa5e03` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 2.6417 | 1 | `eaa7d6058aaa5e03` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | 2.5726 | 2 | `b307c41eb54cf452` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 100% | 2.5574 | 2 | `b307c41eb54cf452` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 1% | 2.6363 | 4 | `ccfe1af393deae98` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 3.2909 | 4 | `ccfe1af393deae98` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | 3.0287 | 8 | `c8c73778097931dc` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 100% | 3.1866 | 8 | `c8c73778097931dc` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 1% | 3.6563 | 16 | `b0c356323d30c425` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 4.3386 | 16 | `b0c356323d30c425` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 0.4685 | 1 | `549a07108ab01ba9` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 2.5348 | 1 | `2ce8b85a7b60c88e` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | 0.6589 | 2 | `28cc913d71d0d72b` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 100% | 5.2653 | 2 | `7dfb959242cb3b9a` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 1% | 1.1180 | 4 | `9f850e53560baafb` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 8.2070 | 4 | `2e2ad57277b50a10` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | 1.8230 | 8 | `477f7fbaf1c627cb` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 100% | 10.3942 | 8 | `5d150ab091a05805` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 1% | 2.9461 | 16 | `1c38b58d8cec08d3` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 16.8836 | 16 | `bafd73f68501cd84` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 5.0600 | false |
| notmat | 1 | 1% | 0.9027 | true |
| mat | 1 | 100% | 1.0454 | false |
| notmat | 1 | 100% | 1.0031 | true |
| mat | 2 | 1% | 1.0517 | true |
| notmat | 2 | 1% | 0.2694 | false |
| mat | 2 | 100% | 0.9954 | true |
| notmat | 2 | 100% | 2.0493 | false |
| mat | 4 | 1% | 0.8384 | true |
| notmat | 4 | 1% | 0.3556 | false |
| mat | 4 | 100% | 1.0124 | true |
| notmat | 4 | 100% | 2.5246 | false |
| mat | 8 | 1% | 1.0410 | true |
| notmat | 8 | 1% | 0.6266 | false |
| mat | 8 | 100% | 0.8870 | true |
| notmat | 8 | 100% | 2.8934 | false |
| mat | 16 | 1% | 1.0529 | true |
| notmat | 16 | 1% | 0.8484 | false |
| mat | 16 | 100% | 1.1138 | true |
| notmat | 16 | 100% | 4.3342 | false |

## Repeated-design position control

| first median ms | last median ms | ratio (last/first) | same plan |
|---|---|---|---|
| 0.477766 | 0.414905 | 0.8684 | true |

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

- client: nr_periods 383, nr_throttled 307, throttled_usec 17621363 | quota 100000 100000
, memory 2147483648

- pg-single: nr_periods 376, nr_throttled 0, throttled_usec 0 | quota 200000 100000, memory 3221225472

## Access path per cell

_Engine-reported from each cell's own plan: whether the rank index or a sequential scan was used, rows scanned, PostgreSQL buffers and filter removals, and reader retries of transient conflicts. A cell that scanned the whole table at a selective threshold is a different physical regime, not a faster query._

| axis | status | index scan | seq scan | storage rows scanned | index rows scanned | buffers hit/read | rows removed | retries |
|---|---|---|---|---|---|---|---|---|
| family=bag;stmt=comma_join;threshold=100000 | valid | true | true | 0 | 0 | 100263/0 | 0 | 0 |
| family=bag;stmt=comma_join;threshold=180000 | valid | false | true | 0 | 0 | 1286/0 | 20000 | 0 |
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
| family=bag;stmt=cte_joined_outer;threshold=20;position=last | valid | true | true | 0 | 0 | 84/0 | 0 | 0 |
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
| family=reuse;policy=default;refs=1;consumer=1% | valid | true | false | 0 | 0 | 21/0 | 1983 | 0 |
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
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | true | false | 0 | 0 | 336/0 | 1983 | 0 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | false | false | 0 | 0 | 0/0 | 0 | 0 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | true | false | 0 | 0 | 21/0 | 1983 | 0 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | false | false | 0 | 0 | 0/0 | 0 | 0 |
| family=reuse;policy=notmat;refs=2;consumer=1% | valid | true | false | 0 | 0 | 42/0 | 1983 | 0 |
| family=reuse;policy=notmat;refs=2;consumer=100% | valid | false | false | 0 | 0 | 0/0 | 0 | 0 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | true | false | 0 | 0 | 84/0 | 1983 | 0 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | false | false | 0 | 0 | 0/0 | 0 | 0 |
| family=reuse;policy=notmat;refs=8;consumer=1% | valid | true | false | 0 | 0 | 168/0 | 1983 | 0 |
| family=reuse;policy=notmat;refs=8;consumer=100% | valid | false | false | 0 | 0 | 0/0 | 0 | 0 |

## Latency samples and derived throughput

_The reader is sequential, so `queries/s` is 1000/median_ms, a derived rate and not a concurrent-capacity test. min/max are the five sample extremes. Each sample includes any retry and backoff time the statement needed after a transient conflict; per-cell retry counts appear in the access-path table._

| axis | status | min ms | median ms | max ms | samples | derived queries/s |
|---|---|---|---|---|---|---|
| family=bag;stmt=comma_join;threshold=100000 | valid | 46.9938 | 48.4774 | 55.2583 | 5 | 20.6 |
| family=bag;stmt=comma_join;threshold=180000 | valid | 61.8378 | 72.7510 | 83.3384 | 5 | 13.7 |
| family=bag;stmt=comma_join;threshold=20 | valid | 0.3702 | 0.5237 | 0.6521 | 5 | 1909.5 |
| family=bag;stmt=comma_join;threshold=200 | valid | 0.6787 | 0.7080 | 1.1944 | 5 | 1412.5 |
| family=bag;stmt=comma_join;threshold=2000 | valid | 2.0833 | 2.1973 | 2.9078 | 5 | 455.1 |
| family=bag;stmt=comma_join;threshold=20000 | valid | 7.3152 | 7.5892 | 14.6444 | 5 | 131.8 |
| family=bag;stmt=cte_filtered;threshold=100000 | valid | 43.5390 | 44.6471 | 54.9190 | 5 | 22.4 |
| family=bag;stmt=cte_filtered;threshold=180000 | valid | 43.4470 | 78.3874 | 86.1072 | 5 | 12.8 |
| family=bag;stmt=cte_filtered;threshold=20 | valid | 0.4251 | 0.4773 | 0.5093 | 5 | 2095.0 |
| family=bag;stmt=cte_filtered;threshold=200 | valid | 0.6335 | 0.7325 | 1.2390 | 5 | 1365.3 |
| family=bag;stmt=cte_filtered;threshold=2000 | valid | 1.6105 | 1.7271 | 1.9136 | 5 | 579.0 |
| family=bag;stmt=cte_filtered;threshold=20000 | valid | 7.2567 | 7.3109 | 7.4609 | 5 | 136.8 |
| family=bag;stmt=cte_filtered_mat;threshold=100000 | valid | 59.6436 | 60.4186 | 68.0205 | 5 | 16.6 |
| family=bag;stmt=cte_filtered_mat;threshold=180000 | valid | 67.7919 | 80.3723 | 96.9650 | 5 | 12.4 |
| family=bag;stmt=cte_filtered_mat;threshold=20 | valid | 0.4028 | 0.4179 | 0.4751 | 5 | 2393.1 |
| family=bag;stmt=cte_filtered_mat;threshold=200 | valid | 0.6640 | 0.7822 | 0.8269 | 5 | 1278.4 |
| family=bag;stmt=cte_filtered_mat;threshold=2000 | valid | 2.1660 | 2.1855 | 2.3368 | 5 | 457.6 |
| family=bag;stmt=cte_filtered_mat;threshold=20000 | valid | 9.4227 | 9.6622 | 44.7759 | 5 | 103.5 |
| family=bag;stmt=cte_filtered_notmat;threshold=100000 | valid | 42.6800 | 48.1178 | 65.6194 | 5 | 20.8 |
| family=bag;stmt=cte_filtered_notmat;threshold=180000 | valid | 67.9800 | 75.4187 | 84.6860 | 5 | 13.3 |
| family=bag;stmt=cte_filtered_notmat;threshold=20 | valid | 0.4502 | 0.4587 | 0.6509 | 5 | 2180.1 |
| family=bag;stmt=cte_filtered_notmat;threshold=200 | valid | 0.9787 | 1.0393 | 1.3323 | 5 | 962.2 |
| family=bag;stmt=cte_filtered_notmat;threshold=2000 | valid | 2.0766 | 2.1203 | 2.2442 | 5 | 471.6 |
| family=bag;stmt=cte_filtered_notmat;threshold=20000 | valid | 7.3639 | 7.5541 | 8.5005 | 5 | 132.4 |
| family=bag;stmt=cte_joined_outer;threshold=100000 | valid | 42.9535 | 46.0983 | 72.7040 | 5 | 21.7 |
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 48.0396 | 64.4969 | 73.8306 | 5 | 15.5 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 0.3029 | 0.4778 | 0.5287 | 5 | 2093.1 |
| family=bag;stmt=cte_joined_outer;threshold=200 | valid | 0.7127 | 0.8495 | 0.9668 | 5 | 1177.1 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 1.9226 | 2.0023 | 2.1026 | 5 | 499.4 |
| family=bag;stmt=cte_joined_outer;threshold=20000 | valid | 7.0903 | 7.4931 | 8.7270 | 5 | 133.5 |
| family=bag;stmt=cte_joined_outer;threshold=20;position=last | valid | 0.3155 | 0.4149 | 0.5871 | 5 | 2410.2 |
| family=bag;stmt=cte_joined_outer_mat;threshold=100000 | valid | 59.9539 | 67.0791 | 71.5180 | 5 | 14.9 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 75.3026 | 78.6928 | 115.2418 | 5 | 12.7 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 47.8469 | 48.1081 | 50.8071 | 5 | 20.8 |
| family=bag;stmt=cte_joined_outer_mat;threshold=200 | valid | 46.7552 | 48.7891 | 51.1116 | 5 | 20.5 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 47.0174 | 47.9916 | 50.1080 | 5 | 20.8 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20000 | valid | 52.7561 | 53.8350 | 55.0978 | 5 | 18.6 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=100000 | valid | 45.0967 | 46.4157 | 56.2578 | 5 | 21.5 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 82.1577 | 89.8642 | 98.5045 | 5 | 11.1 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 0.3861 | 0.4173 | 0.4792 | 5 | 2396.3 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=200 | valid | 0.9604 | 1.0311 | 1.1382 | 5 | 969.8 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 2.0765 | 2.2780 | 2.5172 | 5 | 439.0 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20000 | valid | 7.4879 | 7.7251 | 8.8862 | 5 | 129.4 |
| family=bag;stmt=filtered_derived;threshold=100000 | valid | 42.1504 | 45.2736 | 50.2058 | 5 | 22.1 |
| family=bag;stmt=filtered_derived;threshold=180000 | valid | 47.1764 | 79.8069 | 95.7735 | 5 | 12.5 |
| family=bag;stmt=filtered_derived;threshold=20 | valid | 0.6523 | 0.7151 | 0.9289 | 5 | 1398.5 |
| family=bag;stmt=filtered_derived;threshold=200 | valid | 0.9289 | 0.9991 | 1.1108 | 5 | 1000.9 |
| family=bag;stmt=filtered_derived;threshold=2000 | valid | 1.6716 | 1.7610 | 1.9613 | 5 | 567.9 |
| family=bag;stmt=filtered_derived;threshold=20000 | valid | 7.2520 | 7.6786 | 8.0576 | 5 | 130.2 |
| family=bag;stmt=join_on;threshold=100000 | valid | 42.2145 | 48.9945 | 54.5269 | 5 | 20.4 |
| family=bag;stmt=join_on;threshold=180000 | valid | 48.0823 | 77.7107 | 85.4048 | 5 | 12.9 |
| family=bag;stmt=join_on;threshold=20 | valid | 0.7031 | 0.7237 | 0.7932 | 5 | 1381.8 |
| family=bag;stmt=join_on;threshold=200 | valid | 0.8745 | 0.9253 | 1.2458 | 5 | 1080.7 |
| family=bag;stmt=join_on;threshold=2000 | valid | 1.8574 | 1.9979 | 28.6453 | 5 | 500.5 |
| family=bag;stmt=join_on;threshold=20000 | valid | 7.0371 | 7.5124 | 8.5655 | 5 | 133.1 |
| family=bag;stmt=join_where;threshold=100000 | valid | 42.7571 | 49.6616 | 68.7678 | 5 | 20.1 |
| family=bag;stmt=join_where;threshold=180000 | valid | 47.6952 | 79.8609 | 92.6744 | 5 | 12.5 |
| family=bag;stmt=join_where;threshold=20 | valid | 0.3713 | 0.4308 | 0.5009 | 5 | 2321.1 |
| family=bag;stmt=join_where;threshold=200 | valid | 0.9684 | 1.0708 | 1.1408 | 5 | 933.9 |
| family=bag;stmt=join_where;threshold=2000 | valid | 2.1069 | 2.1664 | 2.2865 | 5 | 461.6 |
| family=bag;stmt=join_where;threshold=20000 | valid | 7.3075 | 7.6453 | 7.9434 | 5 | 130.8 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 3.3982 | 3.4726 | 3.5596 | 5 | 288.0 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 3.7548 | 3.8954 | 4.0835 | 5 | 256.7 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 0.4300 | 0.5191 | 0.6364 | 5 | 1926.5 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 2.4211 | 2.5270 | 2.5636 | 5 | 395.7 |
| family=reuse;policy=default;refs=2;consumer=1% | valid | 2.3943 | 2.4462 | 2.7241 | 5 | 408.8 |
| family=reuse;policy=default;refs=2;consumer=100% | valid | 2.3923 | 2.5693 | 2.8563 | 5 | 389.2 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 3.0433 | 3.1445 | 3.3555 | 5 | 318.0 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 2.9801 | 3.2507 | 3.4366 | 5 | 307.6 |
| family=reuse;policy=default;refs=8;consumer=1% | valid | 2.7180 | 2.9094 | 3.0396 | 5 | 343.7 |
| family=reuse;policy=default;refs=8;consumer=100% | valid | 3.3965 | 3.5924 | 3.8065 | 5 | 278.4 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 3.5812 | 3.6563 | 3.7920 | 5 | 273.5 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 4.1286 | 4.3386 | 4.6350 | 5 | 230.5 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 2.4253 | 2.6265 | 2.7245 | 5 | 380.7 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 2.4782 | 2.6417 | 2.7377 | 5 | 378.5 |
| family=reuse;policy=mat;refs=2;consumer=1% | valid | 2.4517 | 2.5726 | 3.1113 | 5 | 388.7 |
| family=reuse;policy=mat;refs=2;consumer=100% | valid | 2.5058 | 2.5574 | 2.7693 | 5 | 391.0 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 2.5505 | 2.6363 | 2.7023 | 5 | 379.3 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 3.1220 | 3.2909 | 3.7994 | 5 | 303.9 |
| family=reuse;policy=mat;refs=8;consumer=1% | valid | 2.9465 | 3.0287 | 3.0974 | 5 | 330.2 |
| family=reuse;policy=mat;refs=8;consumer=100% | valid | 3.0144 | 3.1866 | 3.2635 | 5 | 313.8 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 2.8451 | 2.9461 | 3.2007 | 5 | 339.4 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 16.1984 | 16.8836 | 18.0556 | 5 | 59.2 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 0.4545 | 0.4685 | 0.5147 | 5 | 2134.3 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 2.2288 | 2.5348 | 2.7165 | 5 | 394.5 |
| family=reuse;policy=notmat;refs=2;consumer=1% | valid | 0.6202 | 0.6589 | 0.7191 | 5 | 1517.6 |
| family=reuse;policy=notmat;refs=2;consumer=100% | valid | 4.9155 | 5.2653 | 5.3929 | 5 | 189.9 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 1.0796 | 1.1180 | 1.3011 | 5 | 894.4 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 8.0739 | 8.2070 | 8.5494 | 5 | 121.8 |
| family=reuse;policy=notmat;refs=8;consumer=1% | valid | 1.6247 | 1.8230 | 2.0277 | 5 | 548.5 |
| family=reuse;policy=notmat;refs=8;consumer=100% | valid | 10.2188 | 10.3942 | 11.7710 | 5 | 96.2 |

## Matched blended workload

- phase: `isolated`; target offered rate: 10 updates/s
- offered 0, committed 0, failed 0 over 0.000 s
- offered rate 0.000/s, completed rate 0.000/s
- per key: no keys touched
- ledger: 0 keys reconciled, 0 mismatch

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42208848`, 5 interleaved rounds
  1. family=bag;stmt=cte_joined_outer;threshold=20
  2. family=reuse;policy=mat;refs=1;consumer=100%
  3. family=reuse;policy=mat;refs=4;consumer=1%
  4. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  5. family=bag;stmt=comma_join;threshold=200
  6. family=bag;stmt=filtered_derived;threshold=2000
  7. family=reuse;policy=default;refs=2;consumer=100%
  8. family=bag;stmt=cte_filtered;threshold=100000
  9. family=bag;stmt=cte_filtered_notmat;threshold=200
  10. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  11. family=bag;stmt=join_on;threshold=200
  12. family=bag;stmt=cte_joined_outer_mat;threshold=200
  13. family=bag;stmt=filtered_derived;threshold=100000
  14. family=bag;stmt=join_on;threshold=100000
  15. family=reuse;policy=mat;refs=16;consumer=100%
  16. family=bag;stmt=cte_filtered_notmat;threshold=100000
  17. family=bag;stmt=join_where;threshold=2000
  18. family=reuse;policy=notmat;refs=1;consumer=100%
  19. family=bag;stmt=join_where;threshold=100000
  20. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  21. family=reuse;policy=default;refs=1;consumer=100%
  22. family=reuse;policy=mat;refs=2;consumer=100%
  23. family=reuse;policy=notmat;refs=1;consumer=1%
  24. family=reuse;policy=notmat;refs=2;consumer=1%
  25. family=bag;stmt=cte_filtered;threshold=20000
  26. family=bag;stmt=cte_filtered_notmat;threshold=2000
  27. family=reuse;policy=default;refs=1;consumer=1%
  28. family=bag;stmt=cte_joined_outer;threshold=180000
  29. family=bag;stmt=filtered_derived;threshold=20000
  30. family=bag;stmt=cte_filtered_mat;threshold=180000
  31. family=reuse;policy=notmat;refs=16;consumer=1%
  32. family=bag;stmt=cte_joined_outer;threshold=20000
  33. family=bag;stmt=filtered_derived;threshold=200
  34. family=bag;stmt=cte_filtered_notmat;threshold=180000
  35. family=bag;stmt=cte_joined_outer_mat;threshold=100000
  36. family=bag;stmt=comma_join;threshold=2000
  37. family=bag;stmt=cte_joined_outer;threshold=200
  38. family=reuse;policy=default;refs=16;consumer=1%
  39. family=bag;stmt=join_where;threshold=180000
  40. family=reuse;policy=notmat;refs=8;consumer=100%
  41. family=bag;stmt=cte_filtered;threshold=180000
  42. family=bag;stmt=join_where;threshold=20000
  43. family=reuse;policy=notmat;refs=16;consumer=100%
  44. family=bag;stmt=comma_join;threshold=20000
  45. family=reuse;policy=default;refs=4;consumer=100%
  46. family=reuse;policy=notmat;refs=4;consumer=100%
  47. family=bag;stmt=cte_filtered_notmat;threshold=20000
  48. family=reuse;policy=default;refs=8;consumer=100%
  49. family=reuse;policy=default;refs=16;consumer=100%
  50. family=bag;stmt=cte_filtered_notmat;threshold=20
  51. family=bag;stmt=cte_joined_outer_mat;threshold=20000
  52. family=bag;stmt=filtered_derived;threshold=180000
  53. family=reuse;policy=notmat;refs=2;consumer=100%
  54. family=reuse;policy=mat;refs=8;consumer=1%
  55. family=bag;stmt=cte_filtered;threshold=2000
  56. family=bag;stmt=cte_filtered;threshold=200
  57. family=bag;stmt=cte_filtered;threshold=20
  58. family=bag;stmt=cte_joined_outer_mat;threshold=20
  59. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  60. family=bag;stmt=join_on;threshold=20
  61. family=bag;stmt=join_where;threshold=20
  62. family=bag;stmt=cte_joined_outer;threshold=2000
  63. family=bag;stmt=cte_joined_outer;threshold=100000
  64. family=bag;stmt=cte_joined_outer_notmat;threshold=100000
  65. family=reuse;policy=mat;refs=4;consumer=100%
  66. family=reuse;policy=mat;refs=1;consumer=1%
  67. family=reuse;policy=default;refs=2;consumer=1%
  68. family=bag;stmt=join_on;threshold=20000
  69. family=bag;stmt=cte_joined_outer_notmat;threshold=200
  70. family=bag;stmt=comma_join;threshold=20
  71. family=bag;stmt=comma_join;threshold=100000
  72. family=bag;stmt=cte_joined_outer_notmat;threshold=20000
  73. family=reuse;policy=default;refs=4;consumer=1%
  74. family=reuse;policy=mat;refs=2;consumer=1%
  75. family=bag;stmt=cte_filtered_mat;threshold=100000
  76. family=bag;stmt=join_where;threshold=200
  77. family=bag;stmt=join_on;threshold=2000
  78. family=bag;stmt=join_on;threshold=180000
  79. family=bag;stmt=filtered_derived;threshold=20
  80. family=reuse;policy=mat;refs=16;consumer=1%
  81. family=bag;stmt=cte_filtered_mat;threshold=200
  82. family=bag;stmt=cte_filtered_mat;threshold=20000
  83. family=reuse;policy=notmat;refs=4;consumer=1%
  84. family=bag;stmt=comma_join;threshold=180000
  85. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  86. family=reuse;policy=notmat;refs=8;consumer=1%
  87. family=bag;stmt=cte_filtered_mat;threshold=2000
  88. family=reuse;policy=mat;refs=8;consumer=100%
  89. family=reuse;policy=default;refs=8;consumer=1%
  90. family=bag;stmt=cte_filtered_mat;threshold=20
  91. family=bag;stmt=cte_joined_outer;threshold=20;position=last

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
| family=bag;stmt=cte_joined_outer;threshold=20;position=last | j05_cte_joined_outer | valid | `61819ba7d0f5f952` | policy=default; observed_cte_scan=false |
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
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `bafd73f68501cd84` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `549a07108ab01ba9` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `2ce8b85a7b60c88e` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=2;consumer=1% | r_notmat_n2_c1 | valid | `28cc913d71d0d72b` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=2;consumer=100% | r_notmat_n2_c100 | valid | `7dfb959242cb3b9a` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `9f850e53560baafb` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `2e2ad57277b50a10` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=8;consumer=1% | r_notmat_n8_c1 | valid | `477f7fbaf1c627cb` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=8;consumer=100% | r_notmat_n8_c100 | valid | `5d150ab091a05805` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, VERBOSE, BUFFERS, FORMAT TEXT` supported on postgres
- one host: this run used node budget 2 CPU 3g and client budget 1 CPU 2g, so budget comparisons are separate classes here, not independent hardware
- isolated phase: no concurrent writers in this run; the blended read/write phase is a separate labelled run
- physical rank-index arm pg-btree-rank; an engine whose planner ignores the index is a different physical design, not a faster one
- distributed RPC counters are N/A on single-node PostgreSQL; shared buffers and rows removed are parsed per cell where printed
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
