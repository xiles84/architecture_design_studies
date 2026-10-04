# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 91 cells, 91 valid, engine `postgres`, topology `single`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-v3-domain-pg-i2-20261005`
- run tag: `run/study-01/select-join-v3-domain-pg-i2-20261005`
- environment: `host-zenbook-ux5406sa`
- repo commit: `a209d4be775397a7721e736d228d5215cab8f3bb`
- engine version: `PostgreSQL 17.11 (Debian 17.11-1.pgdg13+2) on x86_64-pc-linux-gnu, compiled by gcc (Debian 14.2.0-19) 14.2.0, 64-bit`
- database engine image: `docker.io/library/postgres:17.11` @ `7296f210ae81031ec955dbad9a67a84fe958572a2153b8d0826a647522904dc1`
- benchmark client image: `localhost/charitybench-select:select-join-v3-domain-pg-i2-20261005` @ `8436e36bb66caa8178559d54771b88c063e26848d4ef1d1a56eb524b7f0ecca6`
- inputs digest: `c3d49a8a5d29eafe5202c3622b6e7480b0341bfad9fea379e07510b5b359bf96`
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
| 20 (0.01%) | 0.6034 | 0.6145 | 0.6254 | 0.6259 | 0.6552 | 0.4429 | 0.6631 | 44.5245 | 0.7005 | 0.6085 |
| 200 (0.1%) | 0.7072 | 0.7188 | 0.7421 | 0.6765 | 0.7086 | 0.9289 | 0.7310 | 46.5992 | 0.6775 | 0.9260 |
| 2000 (1%) | 1.6287 | 1.8333 | 1.8257 | 1.7242 | 1.7246 | 1.4753 | 1.8919 | 46.6559 | 1.4013 | 1.8245 |
| 20000 (10%) | 6.9037 | 6.8591 | 6.7246 | 6.8664 | 8.4863 | 6.5433 | 6.7855 | 48.0432 | 6.8367 | 8.2156 |
| 100000 (50%) | 47.7222 | 49.5254 | 47.4178 | 47.3831 | 69.6435 | 45.1979 | 59.4300 | 65.7163 | 47.6143 | 46.8084 |
| 180000 (90%) | 72.0447 | 76.0941 | 77.2334 | 50.5913 | 83.6289 | 64.3750 | 47.1836 | 74.7484 | 63.5460 | 48.1376 |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | join_on vs join_where | 0.6034 | 0.6145 | 1.0184 | true |
| 20 | filtered_derived vs join_where | 0.6034 | 0.6254 | 1.0365 | true |
| 20 | cte_filtered vs join_where | 0.6034 | 0.6259 | 1.0372 | true |
| 20 | cte_filtered_mat vs join_where | 0.6034 | 0.6552 | 1.0858 | false |
| 20 | cte_filtered_notmat vs join_where | 0.6034 | 0.4429 | 0.7340 | true |
| 20 | cte_joined_outer vs join_where | 0.6034 | 0.6631 | 1.0989 | true |
| 20 | cte_joined_outer_mat vs join_where | 0.6034 | 44.5245 | 73.7869 | false |
| 20 | cte_joined_outer_notmat vs join_where | 0.6034 | 0.7005 | 1.1608 | true |
| 20 | comma_join vs join_where | 0.6034 | 0.6085 | 1.0084 | true |
| 200 | join_on vs join_where | 0.7072 | 0.7188 | 1.0165 | true |
| 200 | filtered_derived vs join_where | 0.7072 | 0.7421 | 1.0493 | true |
| 200 | cte_filtered vs join_where | 0.7072 | 0.6765 | 0.9565 | true |
| 200 | cte_filtered_mat vs join_where | 0.7072 | 0.7086 | 1.0019 | false |
| 200 | cte_filtered_notmat vs join_where | 0.7072 | 0.9289 | 1.3135 | true |
| 200 | cte_joined_outer vs join_where | 0.7072 | 0.7310 | 1.0336 | true |
| 200 | cte_joined_outer_mat vs join_where | 0.7072 | 46.5992 | 65.8924 | false |
| 200 | cte_joined_outer_notmat vs join_where | 0.7072 | 0.6775 | 0.9580 | true |
| 200 | comma_join vs join_where | 0.7072 | 0.9260 | 1.3094 | true |
| 2000 | join_on vs join_where | 1.6287 | 1.8333 | 1.1256 | true |
| 2000 | filtered_derived vs join_where | 1.6287 | 1.8257 | 1.1209 | true |
| 2000 | cte_filtered vs join_where | 1.6287 | 1.7242 | 1.0586 | true |
| 2000 | cte_filtered_mat vs join_where | 1.6287 | 1.7246 | 1.0588 | false |
| 2000 | cte_filtered_notmat vs join_where | 1.6287 | 1.4753 | 0.9058 | true |
| 2000 | cte_joined_outer vs join_where | 1.6287 | 1.8919 | 1.1616 | true |
| 2000 | cte_joined_outer_mat vs join_where | 1.6287 | 46.6559 | 28.6456 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 1.6287 | 1.4013 | 0.8604 | true |
| 2000 | comma_join vs join_where | 1.6287 | 1.8245 | 1.1202 | true |
| 20000 | join_on vs join_where | 6.9037 | 6.8591 | 0.9935 | true |
| 20000 | filtered_derived vs join_where | 6.9037 | 6.7246 | 0.9741 | true |
| 20000 | cte_filtered vs join_where | 6.9037 | 6.8664 | 0.9946 | true |
| 20000 | cte_filtered_mat vs join_where | 6.9037 | 8.4863 | 1.2292 | false |
| 20000 | cte_filtered_notmat vs join_where | 6.9037 | 6.5433 | 0.9478 | true |
| 20000 | cte_joined_outer vs join_where | 6.9037 | 6.7855 | 0.9829 | true |
| 20000 | cte_joined_outer_mat vs join_where | 6.9037 | 48.0432 | 6.9591 | false |
| 20000 | cte_joined_outer_notmat vs join_where | 6.9037 | 6.8367 | 0.9903 | true |
| 20000 | comma_join vs join_where | 6.9037 | 8.2156 | 1.1900 | true |
| 100000 | join_on vs join_where | 47.7222 | 49.5254 | 1.0378 | true |
| 100000 | filtered_derived vs join_where | 47.7222 | 47.4178 | 0.9936 | true |
| 100000 | cte_filtered vs join_where | 47.7222 | 47.3831 | 0.9929 | true |
| 100000 | cte_filtered_mat vs join_where | 47.7222 | 69.6435 | 1.4594 | false |
| 100000 | cte_filtered_notmat vs join_where | 47.7222 | 45.1979 | 0.9471 | true |
| 100000 | cte_joined_outer vs join_where | 47.7222 | 59.4300 | 1.2453 | true |
| 100000 | cte_joined_outer_mat vs join_where | 47.7222 | 65.7163 | 1.3771 | false |
| 100000 | cte_joined_outer_notmat vs join_where | 47.7222 | 47.6143 | 0.9977 | true |
| 100000 | comma_join vs join_where | 47.7222 | 46.8084 | 0.9809 | true |
| 180000 | join_on vs join_where | 72.0447 | 76.0941 | 1.0562 | true |
| 180000 | filtered_derived vs join_where | 72.0447 | 77.2334 | 1.0720 | true |
| 180000 | cte_filtered vs join_where | 72.0447 | 50.5913 | 0.7022 | true |
| 180000 | cte_filtered_mat vs join_where | 72.0447 | 83.6289 | 1.1608 | false |
| 180000 | cte_filtered_notmat vs join_where | 72.0447 | 64.3750 | 0.8935 | true |
| 180000 | cte_joined_outer vs join_where | 72.0447 | 47.1836 | 0.6549 | true |
| 180000 | cte_joined_outer_mat vs join_where | 72.0447 | 74.7484 | 1.0375 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 72.0447 | 63.5460 | 0.8820 | true |
| 180000 | comma_join vs join_where | 72.0447 | 48.1376 | 0.6682 | true |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_filtered_mat vs cte_filtered | 0.6259 | 0.6552 | 1.0469 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_filtered_notmat vs cte_filtered | 0.6259 | 0.4429 | 0.7077 | true | policy=not_materialized; observed_cte_scan=false |
| 20 | cte_joined_outer_mat vs cte_joined_outer | 0.6631 | 44.5245 | 67.1486 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 0.6631 | 0.7005 | 1.0564 | true | policy=not_materialized; observed_cte_scan=false |
| 200 | cte_filtered_mat vs cte_filtered | 0.6765 | 0.7086 | 1.0474 | false | policy=materialized; observed_cte_scan=true |
| 200 | cte_filtered_notmat vs cte_filtered | 0.6765 | 0.9289 | 1.3731 | true | policy=not_materialized; observed_cte_scan=false |
| 200 | cte_joined_outer_mat vs cte_joined_outer | 0.7310 | 46.5992 | 63.7487 | false | policy=materialized; observed_cte_scan=true |
| 200 | cte_joined_outer_notmat vs cte_joined_outer | 0.7310 | 0.6775 | 0.9268 | true | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_filtered_mat vs cte_filtered | 1.7242 | 1.7246 | 1.0002 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_filtered_notmat vs cte_filtered | 1.7242 | 1.4753 | 0.8557 | true | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 1.8919 | 46.6559 | 24.6607 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 1.8919 | 1.4013 | 0.7407 | true | policy=not_materialized; observed_cte_scan=false |
| 20000 | cte_filtered_mat vs cte_filtered | 6.8664 | 8.4863 | 1.2359 | false | policy=materialized; observed_cte_scan=true |
| 20000 | cte_filtered_notmat vs cte_filtered | 6.8664 | 6.5433 | 0.9529 | true | policy=not_materialized; observed_cte_scan=false |
| 20000 | cte_joined_outer_mat vs cte_joined_outer | 6.7855 | 48.0432 | 7.0803 | false | policy=materialized; observed_cte_scan=true |
| 20000 | cte_joined_outer_notmat vs cte_joined_outer | 6.7855 | 6.8367 | 1.0076 | true | policy=not_materialized; observed_cte_scan=false |
| 100000 | cte_filtered_mat vs cte_filtered | 47.3831 | 69.6435 | 1.4698 | false | policy=materialized; observed_cte_scan=true |
| 100000 | cte_filtered_notmat vs cte_filtered | 47.3831 | 45.1979 | 0.9539 | true | policy=not_materialized; observed_cte_scan=false |
| 100000 | cte_joined_outer_mat vs cte_joined_outer | 59.4300 | 65.7163 | 1.1058 | false | policy=materialized; observed_cte_scan=true |
| 100000 | cte_joined_outer_notmat vs cte_joined_outer | 59.4300 | 47.6143 | 0.8012 | true | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_filtered_mat vs cte_filtered | 50.5913 | 83.6289 | 1.6530 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_filtered_notmat vs cte_filtered | 50.5913 | 64.3750 | 1.2725 | true | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 47.1836 | 74.7484 | 1.5842 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 47.1836 | 63.5460 | 1.3468 | true | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 0.4212 | 1 | `549a07108ab01ba9` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 2.0413 | 1 | `2ce8b85a7b60c88e` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | 2.2995 | 2 | `b307c41eb54cf452` | policy=default; observed_cte_scan=true | valid |
| default | 2 | 100% | 2.5837 | 2 | `b307c41eb54cf452` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 1% | 2.6730 | 4 | `ccfe1af393deae98` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 2.4233 | 4 | `ccfe1af393deae98` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | 2.4455 | 8 | `c8c73778097931dc` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 100% | 2.9031 | 8 | `c8c73778097931dc` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 1% | 3.0183 | 16 | `b0c356323d30c425` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 3.5861 | 16 | `b0c356323d30c425` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 2.5969 | 1 | `eaa7d6058aaa5e03` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 2.5691 | 1 | `eaa7d6058aaa5e03` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | 2.6116 | 2 | `b307c41eb54cf452` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 100% | 2.4663 | 2 | `b307c41eb54cf452` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 1% | 2.5979 | 4 | `ccfe1af393deae98` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 2.8072 | 4 | `ccfe1af393deae98` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | 2.9256 | 8 | `c8c73778097931dc` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 100% | 2.9217 | 8 | `c8c73778097931dc` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 1% | 3.1804 | 16 | `b0c356323d30c425` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 3.3410 | 16 | `b0c356323d30c425` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 0.6853 | 1 | `549a07108ab01ba9` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 2.2056 | 1 | `2ce8b85a7b60c88e` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | 0.7773 | 2 | `28cc913d71d0d72b` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 100% | 4.7124 | 2 | `31361078a00cb8d1` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 1% | 1.1769 | 4 | `9f850e53560baafb` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 7.4035 | 4 | `eea401d25e13fd8b` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | 1.4210 | 8 | `477f7fbaf1c627cb` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 100% | 9.1495 | 8 | `74ddcd0dda4ab3fd` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 1% | 2.5300 | 16 | `1c38b58d8cec08d3` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 14.7374 | 16 | `e60f90fa1692d52c` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 6.1657 | false |
| notmat | 1 | 1% | 1.6270 | true |
| mat | 1 | 100% | 1.2586 | false |
| notmat | 1 | 100% | 1.0805 | true |
| mat | 2 | 1% | 1.1357 | true |
| notmat | 2 | 1% | 0.3380 | false |
| mat | 2 | 100% | 0.9546 | true |
| notmat | 2 | 100% | 1.8239 | false |
| mat | 4 | 1% | 0.9719 | true |
| notmat | 4 | 1% | 0.4403 | false |
| mat | 4 | 100% | 1.1584 | true |
| notmat | 4 | 100% | 3.0551 | false |
| mat | 8 | 1% | 1.1963 | true |
| notmat | 8 | 1% | 0.5811 | false |
| mat | 8 | 100% | 1.0064 | true |
| notmat | 8 | 100% | 3.1516 | false |
| mat | 16 | 1% | 1.0537 | true |
| notmat | 16 | 1% | 0.8382 | false |
| mat | 16 | 100% | 0.9316 | true |
| notmat | 16 | 100% | 4.1096 | false |

## Repeated-design position control

| first median ms | last median ms | ratio (last/first) | same plan |
|---|---|---|---|
| 3.018348 | 2.656370 | 0.8801 | true |

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

- client: nr_periods 369, nr_throttled 290, throttled_usec 13859747 | quota 100000 100000
, memory 2147483648

- pg-single: nr_periods 363, nr_throttled 0, throttled_usec 0 | quota 200000 100000, memory 3221225472

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
| family=reuse;policy=default;refs=16;consumer=1%;position=last | valid | false | false | 0 | 0 | 0/0 | 0 | 0 |
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
| family=bag;stmt=comma_join;threshold=100000 | valid | 44.2612 | 46.8084 | 68.7264 | 5 | 21.4 |
| family=bag;stmt=comma_join;threshold=180000 | valid | 43.6869 | 48.1376 | 76.8200 | 5 | 20.8 |
| family=bag;stmt=comma_join;threshold=20 | valid | 0.5732 | 0.6085 | 0.6764 | 5 | 1643.4 |
| family=bag;stmt=comma_join;threshold=200 | valid | 0.8628 | 0.9260 | 1.0273 | 5 | 1079.9 |
| family=bag;stmt=comma_join;threshold=2000 | valid | 1.7703 | 1.8245 | 23.6559 | 5 | 548.1 |
| family=bag;stmt=comma_join;threshold=20000 | valid | 6.6226 | 8.2156 | 35.1460 | 5 | 121.7 |
| family=bag;stmt=cte_filtered;threshold=100000 | valid | 46.1189 | 47.3831 | 53.1032 | 5 | 21.1 |
| family=bag;stmt=cte_filtered;threshold=180000 | valid | 49.4651 | 50.5913 | 56.3066 | 5 | 19.8 |
| family=bag;stmt=cte_filtered;threshold=20 | valid | 0.5749 | 0.6259 | 0.7287 | 5 | 1597.8 |
| family=bag;stmt=cte_filtered;threshold=200 | valid | 0.6375 | 0.6765 | 0.8991 | 5 | 1478.3 |
| family=bag;stmt=cte_filtered;threshold=2000 | valid | 1.5584 | 1.7242 | 1.9689 | 5 | 580.0 |
| family=bag;stmt=cte_filtered;threshold=20000 | valid | 6.7219 | 6.8664 | 8.1164 | 5 | 145.6 |
| family=bag;stmt=cte_filtered_mat;threshold=100000 | valid | 57.0525 | 69.6435 | 75.5319 | 5 | 14.4 |
| family=bag;stmt=cte_filtered_mat;threshold=180000 | valid | 66.0741 | 83.6289 | 95.2330 | 5 | 12.0 |
| family=bag;stmt=cte_filtered_mat;threshold=20 | valid | 0.5116 | 0.6552 | 0.9547 | 5 | 1526.2 |
| family=bag;stmt=cte_filtered_mat;threshold=200 | valid | 0.6763 | 0.7086 | 0.8046 | 5 | 1411.3 |
| family=bag;stmt=cte_filtered_mat;threshold=2000 | valid | 1.6514 | 1.7246 | 2.4101 | 5 | 579.9 |
| family=bag;stmt=cte_filtered_mat;threshold=20000 | valid | 8.2959 | 8.4863 | 9.1841 | 5 | 117.8 |
| family=bag;stmt=cte_filtered_notmat;threshold=100000 | valid | 43.3785 | 45.1979 | 63.4015 | 5 | 22.1 |
| family=bag;stmt=cte_filtered_notmat;threshold=180000 | valid | 44.8802 | 64.3750 | 94.5165 | 5 | 15.5 |
| family=bag;stmt=cte_filtered_notmat;threshold=20 | valid | 0.3679 | 0.4429 | 0.6800 | 5 | 2257.7 |
| family=bag;stmt=cte_filtered_notmat;threshold=200 | valid | 0.8634 | 0.9289 | 1.0589 | 5 | 1076.6 |
| family=bag;stmt=cte_filtered_notmat;threshold=2000 | valid | 1.4196 | 1.4753 | 1.5942 | 5 | 677.8 |
| family=bag;stmt=cte_filtered_notmat;threshold=20000 | valid | 6.3770 | 6.5433 | 7.9928 | 5 | 152.8 |
| family=bag;stmt=cte_joined_outer;threshold=100000 | valid | 44.2720 | 59.4300 | 66.3882 | 5 | 16.8 |
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 46.2891 | 47.1836 | 48.4301 | 5 | 21.2 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 0.6108 | 0.6631 | 0.8008 | 5 | 1508.1 |
| family=bag;stmt=cte_joined_outer;threshold=200 | valid | 0.7008 | 0.7310 | 0.8895 | 5 | 1368.0 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 1.8182 | 1.8919 | 1.9919 | 5 | 528.6 |
| family=bag;stmt=cte_joined_outer;threshold=20000 | valid | 6.4712 | 6.7855 | 8.8696 | 5 | 147.4 |
| family=bag;stmt=cte_joined_outer_mat;threshold=100000 | valid | 59.7406 | 65.7163 | 66.4912 | 5 | 15.2 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 74.3435 | 74.7484 | 79.9279 | 5 | 13.4 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 40.6995 | 44.5245 | 51.6353 | 5 | 22.5 |
| family=bag;stmt=cte_joined_outer_mat;threshold=200 | valid | 41.6870 | 46.5992 | 51.3750 | 5 | 21.5 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 43.0419 | 46.6559 | 51.3936 | 5 | 21.4 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20000 | valid | 43.7076 | 48.0432 | 52.0591 | 5 | 20.8 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=100000 | valid | 42.8262 | 47.6143 | 70.6763 | 5 | 21.0 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 46.6444 | 63.5460 | 82.0100 | 5 | 15.7 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 0.6173 | 0.7005 | 0.7790 | 5 | 1427.6 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=200 | valid | 0.5799 | 0.6775 | 0.7429 | 5 | 1476.1 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 1.3350 | 1.4013 | 1.9746 | 5 | 713.6 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20000 | valid | 6.4719 | 6.8367 | 50.6399 | 5 | 146.3 |
| family=bag;stmt=filtered_derived;threshold=100000 | valid | 42.8181 | 47.4178 | 66.4071 | 5 | 21.1 |
| family=bag;stmt=filtered_derived;threshold=180000 | valid | 44.2394 | 77.2334 | 84.8752 | 5 | 12.9 |
| family=bag;stmt=filtered_derived;threshold=20 | valid | 0.6121 | 0.6254 | 0.6955 | 5 | 1598.9 |
| family=bag;stmt=filtered_derived;threshold=200 | valid | 0.6636 | 0.7421 | 0.7729 | 5 | 1347.6 |
| family=bag;stmt=filtered_derived;threshold=2000 | valid | 1.7251 | 1.8257 | 1.8907 | 5 | 547.7 |
| family=bag;stmt=filtered_derived;threshold=20000 | valid | 6.4764 | 6.7246 | 21.4879 | 5 | 148.7 |
| family=bag;stmt=join_on;threshold=100000 | valid | 46.9721 | 49.5254 | 64.2661 | 5 | 20.2 |
| family=bag;stmt=join_on;threshold=180000 | valid | 44.0411 | 76.0941 | 94.3189 | 5 | 13.1 |
| family=bag;stmt=join_on;threshold=20 | valid | 0.5349 | 0.6145 | 0.6752 | 5 | 1627.3 |
| family=bag;stmt=join_on;threshold=200 | valid | 0.6708 | 0.7188 | 1.3379 | 5 | 1391.1 |
| family=bag;stmt=join_on;threshold=2000 | valid | 1.7775 | 1.8333 | 2.3126 | 5 | 545.5 |
| family=bag;stmt=join_on;threshold=20000 | valid | 6.4001 | 6.8591 | 38.4126 | 5 | 145.8 |
| family=bag;stmt=join_where;threshold=100000 | valid | 45.6317 | 47.7222 | 53.7393 | 5 | 21.0 |
| family=bag;stmt=join_where;threshold=180000 | valid | 68.7858 | 72.0447 | 96.1866 | 5 | 13.9 |
| family=bag;stmt=join_where;threshold=20 | valid | 0.5570 | 0.6034 | 0.6965 | 5 | 1657.2 |
| family=bag;stmt=join_where;threshold=200 | valid | 0.6358 | 0.7072 | 0.7390 | 5 | 1414.0 |
| family=bag;stmt=join_where;threshold=2000 | valid | 1.5912 | 1.6287 | 1.7510 | 5 | 614.0 |
| family=bag;stmt=join_where;threshold=20000 | valid | 6.8022 | 6.9037 | 7.9382 | 5 | 144.9 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 2.8220 | 3.0183 | 3.1158 | 5 | 331.3 |
| family=reuse;policy=default;refs=16;consumer=1%;position=last | valid | 2.5897 | 2.6564 | 2.7805 | 5 | 376.5 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 3.4042 | 3.5861 | 4.2665 | 5 | 278.9 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 0.3958 | 0.4212 | 28.6241 | 5 | 2374.2 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 1.9442 | 2.0413 | 2.2919 | 5 | 489.9 |
| family=reuse;policy=default;refs=2;consumer=1% | valid | 2.1834 | 2.2995 | 2.7401 | 5 | 434.9 |
| family=reuse;policy=default;refs=2;consumer=100% | valid | 2.5408 | 2.5837 | 2.6218 | 5 | 387.0 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 2.6428 | 2.6730 | 3.4191 | 5 | 374.1 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 2.3505 | 2.4233 | 2.7130 | 5 | 412.7 |
| family=reuse;policy=default;refs=8;consumer=1% | valid | 2.2995 | 2.4455 | 2.6583 | 5 | 408.9 |
| family=reuse;policy=default;refs=8;consumer=100% | valid | 2.8370 | 2.9031 | 2.9445 | 5 | 344.5 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 3.1552 | 3.1804 | 3.4913 | 5 | 314.4 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 3.1617 | 3.3410 | 3.7237 | 5 | 299.3 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 2.3877 | 2.5969 | 2.7594 | 5 | 385.1 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 2.3681 | 2.5691 | 2.6640 | 5 | 389.2 |
| family=reuse;policy=mat;refs=2;consumer=1% | valid | 2.5629 | 2.6116 | 2.8134 | 5 | 382.9 |
| family=reuse;policy=mat;refs=2;consumer=100% | valid | 2.3592 | 2.4663 | 2.4752 | 5 | 405.5 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 2.3765 | 2.5979 | 2.8061 | 5 | 384.9 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 2.7048 | 2.8072 | 2.8860 | 5 | 356.2 |
| family=reuse;policy=mat;refs=8;consumer=1% | valid | 2.8282 | 2.9256 | 3.3507 | 5 | 341.8 |
| family=reuse;policy=mat;refs=8;consumer=100% | valid | 2.8674 | 2.9217 | 3.0154 | 5 | 342.3 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 2.4960 | 2.5300 | 2.7959 | 5 | 395.3 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 14.3600 | 14.7374 | 15.7396 | 5 | 67.9 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 0.5901 | 0.6853 | 0.8336 | 5 | 1459.3 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 2.1605 | 2.2056 | 2.3474 | 5 | 453.4 |
| family=reuse;policy=notmat;refs=2;consumer=1% | valid | 0.7478 | 0.7773 | 0.7904 | 5 | 1286.5 |
| family=reuse;policy=notmat;refs=2;consumer=100% | valid | 4.4928 | 4.7124 | 6.9125 | 5 | 212.2 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 1.1510 | 1.1769 | 1.3310 | 5 | 849.7 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 6.8040 | 7.4035 | 7.9588 | 5 | 135.1 |
| family=reuse;policy=notmat;refs=8;consumer=1% | valid | 1.3889 | 1.4210 | 1.5187 | 5 | 703.7 |
| family=reuse;policy=notmat;refs=8;consumer=100% | valid | 8.8064 | 9.1495 | 9.5679 | 5 | 109.3 |

## Matched blended workload

- phase: `isolated`; target offered rate: 10 updates/s
- offered 0, committed 0, failed 0 over 0.000 s
- offered rate 0.000/s, completed rate 0.000/s
- per key: no keys touched
- ledger: 0 keys reconciled, 0 mismatch

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42405278`, 5 interleaved rounds
  1. family=reuse;policy=default;refs=16;consumer=1%
  2. family=reuse;policy=default;refs=1;consumer=1%
  3. family=reuse;policy=default;refs=2;consumer=1%
  4. family=bag;stmt=join_where;threshold=100000
  5. family=reuse;policy=notmat;refs=1;consumer=1%
  6. family=bag;stmt=cte_joined_outer;threshold=20000
  7. family=reuse;policy=notmat;refs=2;consumer=1%
  8. family=bag;stmt=cte_filtered_mat;threshold=20000
  9. family=bag;stmt=comma_join;threshold=2000
  10. family=bag;stmt=comma_join;threshold=100000
  11. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  12. family=reuse;policy=mat;refs=4;consumer=1%
  13. family=reuse;policy=default;refs=1;consumer=100%
  14. family=bag;stmt=join_on;threshold=200
  15. family=bag;stmt=filtered_derived;threshold=20000
  16. family=reuse;policy=mat;refs=2;consumer=1%
  17. family=bag;stmt=cte_filtered_mat;threshold=2000
  18. family=reuse;policy=default;refs=4;consumer=100%
  19. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  20. family=reuse;policy=mat;refs=16;consumer=100%
  21. family=bag;stmt=cte_joined_outer_mat;threshold=20
  22. family=bag;stmt=cte_joined_outer;threshold=2000
  23. family=bag;stmt=cte_filtered_notmat;threshold=20
  24. family=bag;stmt=cte_filtered;threshold=20000
  25. family=bag;stmt=cte_joined_outer_mat;threshold=100000
  26. family=reuse;policy=mat;refs=8;consumer=1%
  27. family=bag;stmt=cte_joined_outer;threshold=200
  28. family=bag;stmt=join_where;threshold=2000
  29. family=bag;stmt=cte_filtered;threshold=200
  30. family=bag;stmt=cte_filtered_notmat;threshold=100000
  31. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  32. family=reuse;policy=mat;refs=1;consumer=100%
  33. family=bag;stmt=cte_filtered_notmat;threshold=2000
  34. family=bag;stmt=cte_joined_outer;threshold=180000
  35. family=bag;stmt=join_where;threshold=20000
  36. family=bag;stmt=cte_filtered_notmat;threshold=200
  37. family=reuse;policy=mat;refs=2;consumer=100%
  38. family=bag;stmt=join_where;threshold=200
  39. family=reuse;policy=default;refs=16;consumer=100%
  40. family=bag;stmt=comma_join;threshold=20000
  41. family=bag;stmt=filtered_derived;threshold=2000
  42. family=bag;stmt=cte_filtered;threshold=100000
  43. family=bag;stmt=comma_join;threshold=200
  44. family=reuse;policy=default;refs=2;consumer=100%
  45. family=bag;stmt=join_on;threshold=180000
  46. family=reuse;policy=notmat;refs=4;consumer=1%
  47. family=bag;stmt=filtered_derived;threshold=180000
  48. family=reuse;policy=mat;refs=16;consumer=1%
  49. family=reuse;policy=default;refs=8;consumer=1%
  50. family=reuse;policy=notmat;refs=16;consumer=1%
  51. family=bag;stmt=cte_joined_outer_notmat;threshold=100000
  52. family=bag;stmt=cte_joined_outer_mat;threshold=200
  53. family=bag;stmt=join_on;threshold=20
  54. family=bag;stmt=join_on;threshold=2000
  55. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  56. family=bag;stmt=filtered_derived;threshold=20
  57. family=bag;stmt=cte_filtered_mat;threshold=100000
  58. family=bag;stmt=cte_filtered;threshold=20
  59. family=reuse;policy=default;refs=8;consumer=100%
  60. family=bag;stmt=cte_joined_outer_notmat;threshold=20000
  61. family=bag;stmt=join_on;threshold=100000
  62. family=reuse;policy=mat;refs=4;consumer=100%
  63. family=bag;stmt=cte_filtered_mat;threshold=180000
  64. family=bag;stmt=cte_joined_outer;threshold=20
  65. family=reuse;policy=mat;refs=1;consumer=1%
  66. family=bag;stmt=join_on;threshold=20000
  67. family=bag;stmt=cte_filtered_mat;threshold=20
  68. family=reuse;policy=notmat;refs=4;consumer=100%
  69. family=bag;stmt=cte_filtered;threshold=2000
  70. family=reuse;policy=notmat;refs=16;consumer=100%
  71. family=bag;stmt=cte_filtered;threshold=180000
  72. family=bag;stmt=filtered_derived;threshold=100000
  73. family=bag;stmt=join_where;threshold=180000
  74. family=bag;stmt=cte_joined_outer;threshold=100000
  75. family=reuse;policy=default;refs=4;consumer=1%
  76. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  77. family=bag;stmt=cte_filtered_notmat;threshold=180000
  78. family=bag;stmt=cte_filtered_notmat;threshold=20000
  79. family=bag;stmt=comma_join;threshold=20
  80. family=reuse;policy=notmat;refs=1;consumer=100%
  81. family=bag;stmt=filtered_derived;threshold=200
  82. family=bag;stmt=cte_joined_outer_mat;threshold=20000
  83. family=reuse;policy=notmat;refs=2;consumer=100%
  84. family=reuse;policy=notmat;refs=8;consumer=100%
  85. family=bag;stmt=comma_join;threshold=180000
  86. family=bag;stmt=join_where;threshold=20
  87. family=bag;stmt=cte_filtered_mat;threshold=200
  88. family=reuse;policy=mat;refs=8;consumer=100%
  89. family=bag;stmt=cte_joined_outer_notmat;threshold=200
  90. family=reuse;policy=notmat;refs=8;consumer=1%
  91. family=reuse;policy=default;refs=16;consumer=1%;position=last

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
| family=reuse;policy=default;refs=16;consumer=1%;position=last | r_default_n16_c1 | valid | `b0c356323d30c425` | policy=default; observed_cte_scan=true |
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
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `e60f90fa1692d52c` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `549a07108ab01ba9` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `2ce8b85a7b60c88e` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=2;consumer=1% | r_notmat_n2_c1 | valid | `28cc913d71d0d72b` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=2;consumer=100% | r_notmat_n2_c100 | valid | `31361078a00cb8d1` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `9f850e53560baafb` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `eea401d25e13fd8b` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=8;consumer=1% | r_notmat_n8_c1 | valid | `477f7fbaf1c627cb` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=8;consumer=100% | r_notmat_n8_c100 | valid | `74ddcd0dda4ab3fd` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, VERBOSE, BUFFERS, FORMAT TEXT` supported on postgres
- one host: this run used node budget 2 CPU 3g and client budget 1 CPU 2g, so budget comparisons are separate classes here, not independent hardware
- isolated phase: no concurrent writers in this run; the blended read/write phase is a separate labelled run
- physical rank-index arm pg-btree-rank; an engine whose planner ignores the index is a different physical design, not a faster one
- distributed RPC counters are N/A on single-node PostgreSQL; shared buffers and rows removed are parsed per cell where printed
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
