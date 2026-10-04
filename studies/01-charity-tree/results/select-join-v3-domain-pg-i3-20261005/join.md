# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 91 cells, 91 valid, engine `postgres`, topology `single`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-v3-domain-pg-i3-20261005`
- run tag: `run/study-01/select-join-v3-domain-pg-i3-20261005`
- environment: `host-zenbook-ux5406sa`
- repo commit: `a209d4be775397a7721e736d228d5215cab8f3bb`
- engine version: `PostgreSQL 17.11 (Debian 17.11-1.pgdg13+2) on x86_64-pc-linux-gnu, compiled by gcc (Debian 14.2.0-19) 14.2.0, 64-bit`
- database engine image: `docker.io/library/postgres:17.11` @ `7296f210ae81031ec955dbad9a67a84fe958572a2153b8d0826a647522904dc1`
- benchmark client image: `localhost/charitybench-select:select-join-v3-domain-pg-i3-20261005` @ `8436e36bb66caa8178559d54771b88c063e26848d4ef1d1a56eb524b7f0ecca6`
- inputs digest: `32fce3ffd96b9bb2715c5362065426970fa1c2da2769a9e1807e1b99436ff080`
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
| 20 (0.01%) | 0.6773 | 0.6259 | 0.4551 | 0.6454 | 0.4749 | 0.4608 | 0.5682 | 46.6222 | 0.5083 | 0.5078 |
| 200 (0.1%) | 0.6812 | 0.7059 | 0.9262 | 0.7031 | 0.7242 | 1.0383 | 1.0158 | 47.2308 | 0.8536 | 0.7902 |
| 2000 (1%) | 1.8192 | 2.1357 | 1.4979 | 1.7978 | 1.8717 | 1.4847 | 1.7135 | 48.0010 | 1.5447 | 1.4899 |
| 20000 (10%) | 6.8064 | 7.8708 | 6.9476 | 6.5721 | 9.0275 | 6.9524 | 6.5760 | 51.5060 | 6.6156 | 6.6368 |
| 100000 (50%) | 44.1089 | 58.3653 | 58.2360 | 53.5641 | 57.0172 | 43.5881 | 49.6072 | 60.5675 | 44.4901 | 59.4247 |
| 180000 (90%) | 68.4025 | 70.2504 | 76.1134 | 58.2597 | 79.5338 | 49.0500 | 67.2314 | 79.3514 | 77.2493 | 71.0468 |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | join_on vs join_where | 0.6773 | 0.6259 | 0.9241 | true |
| 20 | filtered_derived vs join_where | 0.6773 | 0.4551 | 0.6720 | true |
| 20 | cte_filtered vs join_where | 0.6773 | 0.6454 | 0.9528 | true |
| 20 | cte_filtered_mat vs join_where | 0.6773 | 0.4749 | 0.7012 | false |
| 20 | cte_filtered_notmat vs join_where | 0.6773 | 0.4608 | 0.6803 | true |
| 20 | cte_joined_outer vs join_where | 0.6773 | 0.5682 | 0.8388 | true |
| 20 | cte_joined_outer_mat vs join_where | 0.6773 | 46.6222 | 68.8331 | false |
| 20 | cte_joined_outer_notmat vs join_where | 0.6773 | 0.5083 | 0.7504 | true |
| 20 | comma_join vs join_where | 0.6773 | 0.5078 | 0.7496 | true |
| 200 | join_on vs join_where | 0.6812 | 0.7059 | 1.0362 | true |
| 200 | filtered_derived vs join_where | 0.6812 | 0.9262 | 1.3596 | true |
| 200 | cte_filtered vs join_where | 0.6812 | 0.7031 | 1.0322 | true |
| 200 | cte_filtered_mat vs join_where | 0.6812 | 0.7242 | 1.0630 | false |
| 200 | cte_filtered_notmat vs join_where | 0.6812 | 1.0383 | 1.5241 | true |
| 200 | cte_joined_outer vs join_where | 0.6812 | 1.0158 | 1.4911 | true |
| 200 | cte_joined_outer_mat vs join_where | 0.6812 | 47.2308 | 69.3306 | false |
| 200 | cte_joined_outer_notmat vs join_where | 0.6812 | 0.8536 | 1.2529 | true |
| 200 | comma_join vs join_where | 0.6812 | 0.7902 | 1.1600 | true |
| 2000 | join_on vs join_where | 1.8192 | 2.1357 | 1.1740 | true |
| 2000 | filtered_derived vs join_where | 1.8192 | 1.4979 | 0.8234 | true |
| 2000 | cte_filtered vs join_where | 1.8192 | 1.7978 | 0.9883 | true |
| 2000 | cte_filtered_mat vs join_where | 1.8192 | 1.8717 | 1.0289 | false |
| 2000 | cte_filtered_notmat vs join_where | 1.8192 | 1.4847 | 0.8161 | true |
| 2000 | cte_joined_outer vs join_where | 1.8192 | 1.7135 | 0.9419 | true |
| 2000 | cte_joined_outer_mat vs join_where | 1.8192 | 48.0010 | 26.3864 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 1.8192 | 1.5447 | 0.8492 | true |
| 2000 | comma_join vs join_where | 1.8192 | 1.4899 | 0.8190 | true |
| 20000 | join_on vs join_where | 6.8064 | 7.8708 | 1.1564 | true |
| 20000 | filtered_derived vs join_where | 6.8064 | 6.9476 | 1.0208 | true |
| 20000 | cte_filtered vs join_where | 6.8064 | 6.5721 | 0.9656 | true |
| 20000 | cte_filtered_mat vs join_where | 6.8064 | 9.0275 | 1.3263 | false |
| 20000 | cte_filtered_notmat vs join_where | 6.8064 | 6.9524 | 1.0215 | true |
| 20000 | cte_joined_outer vs join_where | 6.8064 | 6.5760 | 0.9662 | true |
| 20000 | cte_joined_outer_mat vs join_where | 6.8064 | 51.5060 | 7.5673 | false |
| 20000 | cte_joined_outer_notmat vs join_where | 6.8064 | 6.6156 | 0.9720 | true |
| 20000 | comma_join vs join_where | 6.8064 | 6.6368 | 0.9751 | true |
| 100000 | join_on vs join_where | 44.1089 | 58.3653 | 1.3232 | true |
| 100000 | filtered_derived vs join_where | 44.1089 | 58.2360 | 1.3203 | true |
| 100000 | cte_filtered vs join_where | 44.1089 | 53.5641 | 1.2144 | true |
| 100000 | cte_filtered_mat vs join_where | 44.1089 | 57.0172 | 1.2926 | false |
| 100000 | cte_filtered_notmat vs join_where | 44.1089 | 43.5881 | 0.9882 | true |
| 100000 | cte_joined_outer vs join_where | 44.1089 | 49.6072 | 1.1247 | true |
| 100000 | cte_joined_outer_mat vs join_where | 44.1089 | 60.5675 | 1.3731 | false |
| 100000 | cte_joined_outer_notmat vs join_where | 44.1089 | 44.4901 | 1.0086 | true |
| 100000 | comma_join vs join_where | 44.1089 | 59.4247 | 1.3472 | true |
| 180000 | join_on vs join_where | 68.4025 | 70.2504 | 1.0270 | true |
| 180000 | filtered_derived vs join_where | 68.4025 | 76.1134 | 1.1127 | true |
| 180000 | cte_filtered vs join_where | 68.4025 | 58.2597 | 0.8517 | true |
| 180000 | cte_filtered_mat vs join_where | 68.4025 | 79.5338 | 1.1627 | false |
| 180000 | cte_filtered_notmat vs join_where | 68.4025 | 49.0500 | 0.7171 | true |
| 180000 | cte_joined_outer vs join_where | 68.4025 | 67.2314 | 0.9829 | true |
| 180000 | cte_joined_outer_mat vs join_where | 68.4025 | 79.3514 | 1.1601 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 68.4025 | 77.2493 | 1.1293 | true |
| 180000 | comma_join vs join_where | 68.4025 | 71.0468 | 1.0387 | true |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_filtered_mat vs cte_filtered | 0.6454 | 0.4749 | 0.7359 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_filtered_notmat vs cte_filtered | 0.6454 | 0.4608 | 0.7140 | true | policy=not_materialized; observed_cte_scan=false |
| 20 | cte_joined_outer_mat vs cte_joined_outer | 0.5682 | 46.6222 | 82.0584 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 0.5682 | 0.5083 | 0.8946 | true | policy=not_materialized; observed_cte_scan=false |
| 200 | cte_filtered_mat vs cte_filtered | 0.7031 | 0.7242 | 1.0299 | false | policy=materialized; observed_cte_scan=true |
| 200 | cte_filtered_notmat vs cte_filtered | 0.7031 | 1.0383 | 1.4766 | true | policy=not_materialized; observed_cte_scan=false |
| 200 | cte_joined_outer_mat vs cte_joined_outer | 1.0158 | 47.2308 | 46.4948 | false | policy=materialized; observed_cte_scan=true |
| 200 | cte_joined_outer_notmat vs cte_joined_outer | 1.0158 | 0.8536 | 0.8403 | true | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_filtered_mat vs cte_filtered | 1.7978 | 1.8717 | 1.0411 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_filtered_notmat vs cte_filtered | 1.7978 | 1.4847 | 0.8258 | true | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 1.7135 | 48.0010 | 28.0139 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 1.7135 | 1.5447 | 0.9015 | true | policy=not_materialized; observed_cte_scan=false |
| 20000 | cte_filtered_mat vs cte_filtered | 6.5721 | 9.0275 | 1.3736 | false | policy=materialized; observed_cte_scan=true |
| 20000 | cte_filtered_notmat vs cte_filtered | 6.5721 | 6.9524 | 1.0579 | true | policy=not_materialized; observed_cte_scan=false |
| 20000 | cte_joined_outer_mat vs cte_joined_outer | 6.5760 | 51.5060 | 7.8324 | false | policy=materialized; observed_cte_scan=true |
| 20000 | cte_joined_outer_notmat vs cte_joined_outer | 6.5760 | 6.6156 | 1.0060 | true | policy=not_materialized; observed_cte_scan=false |
| 100000 | cte_filtered_mat vs cte_filtered | 53.5641 | 57.0172 | 1.0645 | false | policy=materialized; observed_cte_scan=true |
| 100000 | cte_filtered_notmat vs cte_filtered | 53.5641 | 43.5881 | 0.8138 | true | policy=not_materialized; observed_cte_scan=false |
| 100000 | cte_joined_outer_mat vs cte_joined_outer | 49.6072 | 60.5675 | 1.2209 | false | policy=materialized; observed_cte_scan=true |
| 100000 | cte_joined_outer_notmat vs cte_joined_outer | 49.6072 | 44.4901 | 0.8968 | true | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_filtered_mat vs cte_filtered | 58.2597 | 79.5338 | 1.3652 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_filtered_notmat vs cte_filtered | 58.2597 | 49.0500 | 0.8419 | true | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 67.2314 | 79.3514 | 1.1803 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 67.2314 | 77.2493 | 1.1490 | true | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 0.5260 | 1 | `549a07108ab01ba9` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 2.5388 | 1 | `2ce8b85a7b60c88e` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | 2.6567 | 2 | `b307c41eb54cf452` | policy=default; observed_cte_scan=true | valid |
| default | 2 | 100% | 2.5784 | 2 | `b307c41eb54cf452` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 1% | 2.7738 | 4 | `ccfe1af393deae98` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 2.8979 | 4 | `ccfe1af393deae98` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | 2.6415 | 8 | `c8c73778097931dc` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 100% | 3.2102 | 8 | `c8c73778097931dc` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 1% | 2.9967 | 16 | `b0c356323d30c425` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 3.6697 | 16 | `b0c356323d30c425` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 2.6838 | 1 | `eaa7d6058aaa5e03` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 2.4050 | 1 | `eaa7d6058aaa5e03` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | 2.7918 | 2 | `b307c41eb54cf452` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 100% | 2.6508 | 2 | `b307c41eb54cf452` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 1% | 2.5160 | 4 | `ccfe1af393deae98` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 2.4949 | 4 | `ccfe1af393deae98` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | 2.7092 | 8 | `c8c73778097931dc` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 100% | 3.1322 | 8 | `c8c73778097931dc` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 1% | 3.3497 | 16 | `b0c356323d30c425` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 3.4001 | 16 | `b0c356323d30c425` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 0.5392 | 1 | `549a07108ab01ba9` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 2.3469 | 1 | `2ce8b85a7b60c88e` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | 0.6053 | 2 | `28cc913d71d0d72b` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 100% | 4.6252 | 2 | `c5481318e3b80c84` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 1% | 1.0423 | 4 | `9f850e53560baafb` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 7.0259 | 4 | `1fc38d7217183377` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | 1.4541 | 8 | `477f7fbaf1c627cb` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 100% | 9.1965 | 8 | `5e6834797407f304` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 1% | 2.5299 | 16 | `1c38b58d8cec08d3` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 14.9586 | 16 | `aeae71526854f82e` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 5.1020 | false |
| notmat | 1 | 1% | 1.0250 | true |
| mat | 1 | 100% | 0.9473 | false |
| notmat | 1 | 100% | 0.9244 | true |
| mat | 2 | 1% | 1.0508 | true |
| notmat | 2 | 1% | 0.2278 | false |
| mat | 2 | 100% | 1.0281 | true |
| notmat | 2 | 100% | 1.7938 | false |
| mat | 4 | 1% | 0.9070 | true |
| notmat | 4 | 1% | 0.3758 | false |
| mat | 4 | 100% | 0.8609 | true |
| notmat | 4 | 100% | 2.4245 | false |
| mat | 8 | 1% | 1.0256 | true |
| notmat | 8 | 1% | 0.5505 | false |
| mat | 8 | 100% | 0.9757 | true |
| notmat | 8 | 100% | 2.8647 | false |
| mat | 16 | 1% | 1.1178 | true |
| notmat | 16 | 1% | 0.8442 | false |
| mat | 16 | 100% | 0.9265 | true |
| notmat | 16 | 100% | 4.0763 | false |

## Repeated-design position control

| first median ms | last median ms | ratio (last/first) | same plan |
|---|---|---|---|
| 2.405047 | 2.242638 | 0.9325 | true |

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

- client: nr_periods 366, nr_throttled 294, throttled_usec 16101931 | quota 100000 100000
, memory 2147483648

- pg-single: nr_periods 338, nr_throttled 0, throttled_usec 0 | quota 200000 100000, memory 3221225472

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
| family=reuse;policy=mat;refs=1;consumer=100%;position=last | valid | false | false | 0 | 0 | 0/0 | 0 | 0 |
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
| family=bag;stmt=comma_join;threshold=100000 | valid | 41.8031 | 59.4247 | 62.1821 | 5 | 16.8 |
| family=bag;stmt=comma_join;threshold=180000 | valid | 45.7111 | 71.0468 | 88.5764 | 5 | 14.1 |
| family=bag;stmt=comma_join;threshold=20 | valid | 0.4229 | 0.5078 | 0.5421 | 5 | 1969.5 |
| family=bag;stmt=comma_join;threshold=200 | valid | 0.7228 | 0.7902 | 0.8769 | 5 | 1265.5 |
| family=bag;stmt=comma_join;threshold=2000 | valid | 1.3926 | 1.4899 | 1.6075 | 5 | 671.2 |
| family=bag;stmt=comma_join;threshold=20000 | valid | 6.3736 | 6.6368 | 7.1368 | 5 | 150.7 |
| family=bag;stmt=cte_filtered;threshold=100000 | valid | 51.4865 | 53.5641 | 57.0319 | 5 | 18.7 |
| family=bag;stmt=cte_filtered;threshold=180000 | valid | 47.1014 | 58.2597 | 66.8286 | 5 | 17.2 |
| family=bag;stmt=cte_filtered;threshold=20 | valid | 0.5837 | 0.6454 | 0.6697 | 5 | 1549.5 |
| family=bag;stmt=cte_filtered;threshold=200 | valid | 0.6323 | 0.7031 | 1.1880 | 5 | 1422.2 |
| family=bag;stmt=cte_filtered;threshold=2000 | valid | 1.7071 | 1.7978 | 1.8969 | 5 | 556.2 |
| family=bag;stmt=cte_filtered;threshold=20000 | valid | 6.4291 | 6.5721 | 6.7055 | 5 | 152.2 |
| family=bag;stmt=cte_filtered_mat;threshold=100000 | valid | 56.6128 | 57.0172 | 59.1073 | 5 | 17.5 |
| family=bag;stmt=cte_filtered_mat;threshold=180000 | valid | 65.5078 | 79.5338 | 91.8223 | 5 | 12.6 |
| family=bag;stmt=cte_filtered_mat;threshold=20 | valid | 0.4025 | 0.4749 | 0.8772 | 5 | 2105.6 |
| family=bag;stmt=cte_filtered_mat;threshold=200 | valid | 0.6031 | 0.7242 | 0.7982 | 5 | 1380.9 |
| family=bag;stmt=cte_filtered_mat;threshold=2000 | valid | 1.8135 | 1.8717 | 2.4279 | 5 | 534.3 |
| family=bag;stmt=cte_filtered_mat;threshold=20000 | valid | 8.7783 | 9.0275 | 10.4827 | 5 | 110.8 |
| family=bag;stmt=cte_filtered_notmat;threshold=100000 | valid | 41.0108 | 43.5881 | 45.9460 | 5 | 22.9 |
| family=bag;stmt=cte_filtered_notmat;threshold=180000 | valid | 43.5978 | 49.0500 | 50.9541 | 5 | 20.4 |
| family=bag;stmt=cte_filtered_notmat;threshold=20 | valid | 0.4173 | 0.4608 | 0.5004 | 5 | 2170.3 |
| family=bag;stmt=cte_filtered_notmat;threshold=200 | valid | 0.9616 | 1.0383 | 1.2333 | 5 | 963.1 |
| family=bag;stmt=cte_filtered_notmat;threshold=2000 | valid | 1.4258 | 1.4847 | 1.6255 | 5 | 673.6 |
| family=bag;stmt=cte_filtered_notmat;threshold=20000 | valid | 6.7363 | 6.9524 | 39.5430 | 5 | 143.8 |
| family=bag;stmt=cte_joined_outer;threshold=100000 | valid | 43.1174 | 49.6072 | 64.4295 | 5 | 20.2 |
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 46.5607 | 67.2314 | 82.3494 | 5 | 14.9 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 0.4835 | 0.5682 | 0.6634 | 5 | 1760.1 |
| family=bag;stmt=cte_joined_outer;threshold=200 | valid | 0.8704 | 1.0158 | 1.1673 | 5 | 984.4 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 1.6569 | 1.7135 | 17.0356 | 5 | 583.6 |
| family=bag;stmt=cte_joined_outer;threshold=20000 | valid | 6.3143 | 6.5760 | 7.0670 | 5 | 152.1 |
| family=bag;stmt=cte_joined_outer_mat;threshold=100000 | valid | 59.5221 | 60.5675 | 67.3126 | 5 | 16.5 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 78.4601 | 79.3514 | 93.7815 | 5 | 12.6 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 45.3075 | 46.6222 | 48.4925 | 5 | 21.4 |
| family=bag;stmt=cte_joined_outer_mat;threshold=200 | valid | 41.3869 | 47.2308 | 49.2836 | 5 | 21.2 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 43.3734 | 48.0010 | 50.7353 | 5 | 20.8 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20000 | valid | 49.7248 | 51.5060 | 53.4473 | 5 | 19.4 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=100000 | valid | 44.0661 | 44.4901 | 50.1514 | 5 | 22.5 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 44.8232 | 77.2493 | 99.9417 | 5 | 12.9 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 0.4514 | 0.5083 | 0.5713 | 5 | 1967.5 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=200 | valid | 0.7656 | 0.8536 | 0.9147 | 5 | 1171.6 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 1.4101 | 1.5447 | 1.7547 | 5 | 647.4 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20000 | valid | 6.4692 | 6.6156 | 10.8305 | 5 | 151.2 |
| family=bag;stmt=filtered_derived;threshold=100000 | valid | 46.2745 | 58.2360 | 66.1644 | 5 | 17.2 |
| family=bag;stmt=filtered_derived;threshold=180000 | valid | 45.6711 | 76.1134 | 83.5714 | 5 | 13.1 |
| family=bag;stmt=filtered_derived;threshold=20 | valid | 0.4198 | 0.4551 | 0.6606 | 5 | 2197.1 |
| family=bag;stmt=filtered_derived;threshold=200 | valid | 0.9027 | 0.9262 | 1.1687 | 5 | 1079.6 |
| family=bag;stmt=filtered_derived;threshold=2000 | valid | 1.4639 | 1.4979 | 1.6328 | 5 | 667.6 |
| family=bag;stmt=filtered_derived;threshold=20000 | valid | 6.2441 | 6.9476 | 7.1546 | 5 | 143.9 |
| family=bag;stmt=join_on;threshold=100000 | valid | 43.3971 | 58.3653 | 64.0824 | 5 | 17.1 |
| family=bag;stmt=join_on;threshold=180000 | valid | 43.9253 | 70.2504 | 75.3331 | 5 | 14.2 |
| family=bag;stmt=join_on;threshold=20 | valid | 0.6069 | 0.6259 | 0.8549 | 5 | 1597.7 |
| family=bag;stmt=join_on;threshold=200 | valid | 0.5956 | 0.7059 | 1.2934 | 5 | 1416.7 |
| family=bag;stmt=join_on;threshold=2000 | valid | 2.0221 | 2.1357 | 2.2028 | 5 | 468.2 |
| family=bag;stmt=join_on;threshold=20000 | valid | 7.0632 | 7.8708 | 34.9434 | 5 | 127.1 |
| family=bag;stmt=join_where;threshold=100000 | valid | 42.8096 | 44.1089 | 55.9654 | 5 | 22.7 |
| family=bag;stmt=join_where;threshold=180000 | valid | 43.5170 | 68.4025 | 86.8644 | 5 | 14.6 |
| family=bag;stmt=join_where;threshold=20 | valid | 0.5710 | 0.6773 | 1.0821 | 5 | 1476.4 |
| family=bag;stmt=join_where;threshold=200 | valid | 0.6529 | 0.6812 | 0.7853 | 5 | 1467.9 |
| family=bag;stmt=join_where;threshold=2000 | valid | 1.6716 | 1.8192 | 2.3610 | 5 | 549.7 |
| family=bag;stmt=join_where;threshold=20000 | valid | 6.7076 | 6.8064 | 7.3738 | 5 | 146.9 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 2.9728 | 2.9967 | 3.0844 | 5 | 333.7 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 3.5509 | 3.6697 | 3.6994 | 5 | 272.5 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 0.4632 | 0.5260 | 0.5859 | 5 | 1901.1 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 2.3177 | 2.5388 | 2.7636 | 5 | 393.9 |
| family=reuse;policy=default;refs=2;consumer=1% | valid | 2.4940 | 2.6567 | 2.8705 | 5 | 376.4 |
| family=reuse;policy=default;refs=2;consumer=100% | valid | 2.4086 | 2.5784 | 2.6940 | 5 | 387.8 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 2.6830 | 2.7738 | 2.9569 | 5 | 360.5 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 2.7897 | 2.8979 | 3.1844 | 5 | 345.1 |
| family=reuse;policy=default;refs=8;consumer=1% | valid | 2.5444 | 2.6415 | 2.9153 | 5 | 378.6 |
| family=reuse;policy=default;refs=8;consumer=100% | valid | 2.9769 | 3.2102 | 3.4109 | 5 | 311.5 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 3.1654 | 3.3497 | 3.6352 | 5 | 298.5 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 3.1424 | 3.4001 | 3.6755 | 5 | 294.1 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 2.5950 | 2.6838 | 2.8086 | 5 | 372.6 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 2.1936 | 2.4050 | 2.5005 | 5 | 415.8 |
| family=reuse;policy=mat;refs=1;consumer=100%;position=last | valid | 2.1825 | 2.2426 | 2.3405 | 5 | 445.9 |
| family=reuse;policy=mat;refs=2;consumer=1% | valid | 2.6054 | 2.7918 | 2.9934 | 5 | 358.2 |
| family=reuse;policy=mat;refs=2;consumer=100% | valid | 2.5745 | 2.6508 | 3.1108 | 5 | 377.3 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 2.2412 | 2.5160 | 2.6468 | 5 | 397.5 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 2.4636 | 2.4949 | 2.6681 | 5 | 400.8 |
| family=reuse;policy=mat;refs=8;consumer=1% | valid | 2.5774 | 2.7092 | 3.3387 | 5 | 369.1 |
| family=reuse;policy=mat;refs=8;consumer=100% | valid | 2.8990 | 3.1322 | 4.0345 | 5 | 319.3 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 2.4957 | 2.5299 | 2.8893 | 5 | 395.3 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 14.1324 | 14.9586 | 15.6479 | 5 | 66.9 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 0.4364 | 0.5392 | 0.5744 | 5 | 1854.8 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 2.3273 | 2.3469 | 2.5299 | 5 | 426.1 |
| family=reuse;policy=notmat;refs=2;consumer=1% | valid | 0.5647 | 0.6053 | 0.8935 | 5 | 1652.2 |
| family=reuse;policy=notmat;refs=2;consumer=100% | valid | 4.2394 | 4.6252 | 5.3566 | 5 | 216.2 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 1.0149 | 1.0423 | 1.1570 | 5 | 959.4 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 6.7581 | 7.0259 | 7.1980 | 5 | 142.3 |
| family=reuse;policy=notmat;refs=8;consumer=1% | valid | 1.4015 | 1.4541 | 1.5261 | 5 | 687.7 |
| family=reuse;policy=notmat;refs=8;consumer=100% | valid | 8.5397 | 9.1965 | 9.6901 | 5 | 108.7 |

## Matched blended workload

- phase: `isolated`; target offered rate: 10 updates/s
- offered 0, committed 0, failed 0 over 0.000 s
- offered rate 0.000/s, completed rate 0.000/s
- per key: no keys touched
- ledger: 0 keys reconciled, 0 mismatch

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42177627`, 5 interleaved rounds
  1. family=reuse;policy=mat;refs=1;consumer=100%
  2. family=bag;stmt=cte_filtered;threshold=200
  3. family=bag;stmt=comma_join;threshold=180000
  4. family=bag;stmt=cte_joined_outer_mat;threshold=20000
  5. family=bag;stmt=filtered_derived;threshold=20000
  6. family=reuse;policy=default;refs=2;consumer=1%
  7. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  8. family=bag;stmt=cte_filtered_notmat;threshold=100000
  9. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  10. family=bag;stmt=cte_filtered;threshold=20
  11. family=bag;stmt=cte_filtered_notmat;threshold=20000
  12. family=bag;stmt=cte_joined_outer;threshold=200
  13. family=bag;stmt=cte_joined_outer_mat;threshold=20
  14. family=reuse;policy=default;refs=16;consumer=100%
  15. family=bag;stmt=cte_filtered;threshold=20000
  16. family=bag;stmt=cte_joined_outer_notmat;threshold=100000
  17. family=reuse;policy=mat;refs=16;consumer=1%
  18. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  19. family=bag;stmt=filtered_derived;threshold=180000
  20. family=reuse;policy=default;refs=1;consumer=100%
  21. family=bag;stmt=comma_join;threshold=100000
  22. family=bag;stmt=filtered_derived;threshold=100000
  23. family=bag;stmt=join_on;threshold=20000
  24. family=reuse;policy=mat;refs=8;consumer=100%
  25. family=bag;stmt=cte_filtered_mat;threshold=20
  26. family=bag;stmt=join_where;threshold=100000
  27. family=bag;stmt=cte_filtered_notmat;threshold=200
  28. family=bag;stmt=join_where;threshold=180000
  29. family=reuse;policy=mat;refs=1;consumer=1%
  30. family=reuse;policy=mat;refs=8;consumer=1%
  31. family=reuse;policy=mat;refs=4;consumer=1%
  32. family=bag;stmt=cte_filtered_notmat;threshold=2000
  33. family=reuse;policy=default;refs=16;consumer=1%
  34. family=bag;stmt=cte_filtered_mat;threshold=20000
  35. family=bag;stmt=filtered_derived;threshold=200
  36. family=reuse;policy=notmat;refs=1;consumer=1%
  37. family=bag;stmt=join_where;threshold=2000
  38. family=bag;stmt=cte_filtered_notmat;threshold=20
  39. family=bag;stmt=cte_filtered_mat;threshold=200
  40. family=bag;stmt=join_on;threshold=200
  41. family=bag;stmt=cte_joined_outer;threshold=2000
  42. family=bag;stmt=filtered_derived;threshold=2000
  43. family=reuse;policy=notmat;refs=2;consumer=100%
  44. family=reuse;policy=mat;refs=16;consumer=100%
  45. family=bag;stmt=cte_joined_outer;threshold=20000
  46. family=bag;stmt=cte_filtered;threshold=2000
  47. family=bag;stmt=comma_join;threshold=2000
  48. family=bag;stmt=cte_filtered_mat;threshold=100000
  49. family=reuse;policy=default;refs=4;consumer=100%
  50. family=bag;stmt=cte_joined_outer_mat;threshold=200
  51. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  52. family=reuse;policy=notmat;refs=1;consumer=100%
  53. family=bag;stmt=filtered_derived;threshold=20
  54. family=bag;stmt=cte_filtered_notmat;threshold=180000
  55. family=bag;stmt=cte_filtered_mat;threshold=180000
  56. family=reuse;policy=mat;refs=2;consumer=1%
  57. family=bag;stmt=cte_filtered_mat;threshold=2000
  58. family=bag;stmt=join_where;threshold=200
  59. family=reuse;policy=default;refs=8;consumer=1%
  60. family=bag;stmt=join_on;threshold=180000
  61. family=reuse;policy=notmat;refs=16;consumer=1%
  62. family=bag;stmt=join_on;threshold=100000
  63. family=bag;stmt=join_on;threshold=2000
  64. family=reuse;policy=notmat;refs=2;consumer=1%
  65. family=reuse;policy=notmat;refs=8;consumer=100%
  66. family=bag;stmt=comma_join;threshold=200
  67. family=bag;stmt=cte_filtered;threshold=100000
  68. family=reuse;policy=notmat;refs=16;consumer=100%
  69. family=bag;stmt=cte_joined_outer_notmat;threshold=200
  70. family=bag;stmt=comma_join;threshold=20000
  71. family=reuse;policy=notmat;refs=8;consumer=1%
  72. family=reuse;policy=default;refs=2;consumer=100%
  73. family=bag;stmt=cte_filtered;threshold=180000
  74. family=bag;stmt=join_where;threshold=20
  75. family=bag;stmt=cte_joined_outer_mat;threshold=100000
  76. family=reuse;policy=notmat;refs=4;consumer=1%
  77. family=bag;stmt=cte_joined_outer;threshold=20
  78. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  79. family=reuse;policy=default;refs=4;consumer=1%
  80. family=reuse;policy=mat;refs=4;consumer=100%
  81. family=bag;stmt=join_where;threshold=20000
  82. family=bag;stmt=cte_joined_outer;threshold=100000
  83. family=bag;stmt=join_on;threshold=20
  84. family=reuse;policy=mat;refs=2;consumer=100%
  85. family=bag;stmt=cte_joined_outer_notmat;threshold=20000
  86. family=reuse;policy=notmat;refs=4;consumer=100%
  87. family=reuse;policy=default;refs=1;consumer=1%
  88. family=bag;stmt=cte_joined_outer;threshold=180000
  89. family=reuse;policy=default;refs=8;consumer=100%
  90. family=bag;stmt=comma_join;threshold=20
  91. family=reuse;policy=mat;refs=1;consumer=100%;position=last

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
| family=reuse;policy=mat;refs=16;consumer=100% | r_mat_n16_c100 | valid | `b0c356323d30c425` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=1% | r_mat_n1_c1 | valid | `eaa7d6058aaa5e03` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=100% | r_mat_n1_c100 | valid | `eaa7d6058aaa5e03` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=100%;position=last | r_mat_n1_c100 | valid | `eaa7d6058aaa5e03` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=2;consumer=1% | r_mat_n2_c1 | valid | `b307c41eb54cf452` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=2;consumer=100% | r_mat_n2_c100 | valid | `b307c41eb54cf452` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=1% | r_mat_n4_c1 | valid | `ccfe1af393deae98` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=100% | r_mat_n4_c100 | valid | `ccfe1af393deae98` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=8;consumer=1% | r_mat_n8_c1 | valid | `c8c73778097931dc` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=8;consumer=100% | r_mat_n8_c100 | valid | `c8c73778097931dc` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=notmat;refs=16;consumer=1% | r_notmat_n16_c1 | valid | `1c38b58d8cec08d3` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `aeae71526854f82e` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `549a07108ab01ba9` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `2ce8b85a7b60c88e` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=2;consumer=1% | r_notmat_n2_c1 | valid | `28cc913d71d0d72b` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=2;consumer=100% | r_notmat_n2_c100 | valid | `c5481318e3b80c84` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `9f850e53560baafb` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `1fc38d7217183377` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=8;consumer=1% | r_notmat_n8_c1 | valid | `477f7fbaf1c627cb` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=8;consumer=100% | r_notmat_n8_c100 | valid | `5e6834797407f304` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, VERBOSE, BUFFERS, FORMAT TEXT` supported on postgres
- one host: this run used node budget 2 CPU 3g and client budget 1 CPU 2g, so budget comparisons are separate classes here, not independent hardware
- isolated phase: no concurrent writers in this run; the blended read/write phase is a separate labelled run
- physical rank-index arm pg-btree-rank; an engine whose planner ignores the index is a different physical design, not a faster one
- distributed RPC counters are N/A on single-node PostgreSQL; shared buffers and rows removed are parsed per cell where printed
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
