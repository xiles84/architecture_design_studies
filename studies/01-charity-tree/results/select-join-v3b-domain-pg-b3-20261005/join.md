# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 91 cells, 91 valid, engine `postgres`, topology `single`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-v3b-domain-pg-b3-20261005`
- run tag: `run/study-01/select-join-v3b-domain-pg-b3-20261005`
- environment: `host-zenbook-ux5406sa`
- repo commit: `7243fee074bfb16fc5ad8b6afc1a2d7549080935`
- engine version: `PostgreSQL 17.11 (Debian 17.11-1.pgdg13+2) on x86_64-pc-linux-gnu, compiled by gcc (Debian 14.2.0-19) 14.2.0, 64-bit`
- database engine image: `docker.io/library/postgres:17.11` @ `7296f210ae81031ec955dbad9a67a84fe958572a2153b8d0826a647522904dc1`
- benchmark client image: `localhost/charitybench-select:select-join-v3b-domain-pg-b3-20261005` @ `f3f1e10c44cbcb4c7b5bc56ed95a4f071682b6f7d33b5e0a7b688c4271eeef88`
- inputs digest: `db176a1e371b3c5b22dd520f8087ebf675cf7129560391c68294592c11370cb6`
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
| 20 (0.01%) | 0.5992 | 0.6739 | 0.4485 | 0.4595 | 0.6517 | 0.4534 | 0.4514 | 45.6655 | 0.6821 | 0.5582 |
| 200 (0.1%) | 0.7896 | 0.7083 | 0.8574 | 0.9393 | 0.7618 | 0.9332 | 0.7179 | 45.6113 | 0.7267 | 0.8816 |
| 2000 (1%) | 2.0829 | 1.6713 | 1.7262 | 1.9683 | 2.3361 | 1.5831 | 1.7460 | 43.6970 | 2.0480 | 1.9171 |
| 20000 (10%) | 6.9500 | 6.9207 | 6.7937 | 6.8580 | 9.2961 | 6.8165 | 7.4247 | 45.3416 | 7.1753 | 6.9819 |
| 100000 (50%) | 47.5226 | 51.3712 | 48.0782 | 47.7381 | 58.4158 | 44.3648 | 49.2412 | 65.4836 | 48.8755 | 47.7892 |
| 180000 (90%) | 62.4210 | 81.3918 | 50.7758 | 78.6819 | 84.1072 | 77.7130 | 79.0973 | 90.7410 | 44.9244 | 75.9969 |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | join_on vs join_where | 0.5992 | 0.6739 | 1.1248 | true |
| 20 | filtered_derived vs join_where | 0.5992 | 0.4485 | 0.7485 | true |
| 20 | cte_filtered vs join_where | 0.5992 | 0.4595 | 0.7668 | true |
| 20 | cte_filtered_mat vs join_where | 0.5992 | 0.6517 | 1.0877 | false |
| 20 | cte_filtered_notmat vs join_where | 0.5992 | 0.4534 | 0.7567 | true |
| 20 | cte_joined_outer vs join_where | 0.5992 | 0.4514 | 0.7535 | true |
| 20 | cte_joined_outer_mat vs join_where | 0.5992 | 45.6655 | 76.2163 | false |
| 20 | cte_joined_outer_notmat vs join_where | 0.5992 | 0.6821 | 1.1384 | true |
| 20 | comma_join vs join_where | 0.5992 | 0.5582 | 0.9317 | true |
| 200 | join_on vs join_where | 0.7896 | 0.7083 | 0.8970 | true |
| 200 | filtered_derived vs join_where | 0.7896 | 0.8574 | 1.0858 | true |
| 200 | cte_filtered vs join_where | 0.7896 | 0.9393 | 1.1896 | true |
| 200 | cte_filtered_mat vs join_where | 0.7896 | 0.7618 | 0.9647 | false |
| 200 | cte_filtered_notmat vs join_where | 0.7896 | 0.9332 | 1.1817 | true |
| 200 | cte_joined_outer vs join_where | 0.7896 | 0.7179 | 0.9091 | true |
| 200 | cte_joined_outer_mat vs join_where | 0.7896 | 45.6113 | 57.7620 | false |
| 200 | cte_joined_outer_notmat vs join_where | 0.7896 | 0.7267 | 0.9203 | true |
| 200 | comma_join vs join_where | 0.7896 | 0.8816 | 1.1165 | true |
| 2000 | join_on vs join_where | 2.0829 | 1.6713 | 0.8024 | true |
| 2000 | filtered_derived vs join_where | 2.0829 | 1.7262 | 0.8287 | true |
| 2000 | cte_filtered vs join_where | 2.0829 | 1.9683 | 0.9450 | true |
| 2000 | cte_filtered_mat vs join_where | 2.0829 | 2.3361 | 1.1216 | false |
| 2000 | cte_filtered_notmat vs join_where | 2.0829 | 1.5831 | 0.7600 | true |
| 2000 | cte_joined_outer vs join_where | 2.0829 | 1.7460 | 0.8383 | true |
| 2000 | cte_joined_outer_mat vs join_where | 2.0829 | 43.6970 | 20.9793 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 2.0829 | 2.0480 | 0.9833 | true |
| 2000 | comma_join vs join_where | 2.0829 | 1.9171 | 0.9204 | true |
| 20000 | join_on vs join_where | 6.9500 | 6.9207 | 0.9958 | true |
| 20000 | filtered_derived vs join_where | 6.9500 | 6.7937 | 0.9775 | true |
| 20000 | cte_filtered vs join_where | 6.9500 | 6.8580 | 0.9868 | true |
| 20000 | cte_filtered_mat vs join_where | 6.9500 | 9.2961 | 1.3376 | false |
| 20000 | cte_filtered_notmat vs join_where | 6.9500 | 6.8165 | 0.9808 | true |
| 20000 | cte_joined_outer vs join_where | 6.9500 | 7.4247 | 1.0683 | true |
| 20000 | cte_joined_outer_mat vs join_where | 6.9500 | 45.3416 | 6.5240 | false |
| 20000 | cte_joined_outer_notmat vs join_where | 6.9500 | 7.1753 | 1.0324 | true |
| 20000 | comma_join vs join_where | 6.9500 | 6.9819 | 1.0046 | true |
| 100000 | join_on vs join_where | 47.5226 | 51.3712 | 1.0810 | true |
| 100000 | filtered_derived vs join_where | 47.5226 | 48.0782 | 1.0117 | true |
| 100000 | cte_filtered vs join_where | 47.5226 | 47.7381 | 1.0045 | true |
| 100000 | cte_filtered_mat vs join_where | 47.5226 | 58.4158 | 1.2292 | false |
| 100000 | cte_filtered_notmat vs join_where | 47.5226 | 44.3648 | 0.9335 | true |
| 100000 | cte_joined_outer vs join_where | 47.5226 | 49.2412 | 1.0362 | true |
| 100000 | cte_joined_outer_mat vs join_where | 47.5226 | 65.4836 | 1.3779 | false |
| 100000 | cte_joined_outer_notmat vs join_where | 47.5226 | 48.8755 | 1.0285 | true |
| 100000 | comma_join vs join_where | 47.5226 | 47.7892 | 1.0056 | true |
| 180000 | join_on vs join_where | 62.4210 | 81.3918 | 1.3039 | true |
| 180000 | filtered_derived vs join_where | 62.4210 | 50.7758 | 0.8134 | true |
| 180000 | cte_filtered vs join_where | 62.4210 | 78.6819 | 1.2605 | true |
| 180000 | cte_filtered_mat vs join_where | 62.4210 | 84.1072 | 1.3474 | false |
| 180000 | cte_filtered_notmat vs join_where | 62.4210 | 77.7130 | 1.2450 | true |
| 180000 | cte_joined_outer vs join_where | 62.4210 | 79.0973 | 1.2672 | true |
| 180000 | cte_joined_outer_mat vs join_where | 62.4210 | 90.7410 | 1.4537 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 62.4210 | 44.9244 | 0.7197 | true |
| 180000 | comma_join vs join_where | 62.4210 | 75.9969 | 1.2175 | true |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_filtered_mat vs cte_filtered | 0.4595 | 0.6517 | 1.4184 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_filtered_notmat vs cte_filtered | 0.4595 | 0.4534 | 0.9868 | true | policy=not_materialized; observed_cte_scan=false |
| 20 | cte_joined_outer_mat vs cte_joined_outer | 0.4514 | 45.6655 | 101.1532 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 0.4514 | 0.6821 | 1.5108 | true | policy=not_materialized; observed_cte_scan=false |
| 200 | cte_filtered_mat vs cte_filtered | 0.9393 | 0.7618 | 0.8110 | false | policy=materialized; observed_cte_scan=true |
| 200 | cte_filtered_notmat vs cte_filtered | 0.9393 | 0.9332 | 0.9934 | true | policy=not_materialized; observed_cte_scan=false |
| 200 | cte_joined_outer_mat vs cte_joined_outer | 0.7179 | 45.6113 | 63.5368 | false | policy=materialized; observed_cte_scan=true |
| 200 | cte_joined_outer_notmat vs cte_joined_outer | 0.7179 | 0.7267 | 1.0123 | true | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_filtered_mat vs cte_filtered | 1.9683 | 2.3361 | 1.1869 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_filtered_notmat vs cte_filtered | 1.9683 | 1.5831 | 0.8043 | true | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 1.7460 | 43.6970 | 25.0268 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 1.7460 | 2.0480 | 1.1730 | true | policy=not_materialized; observed_cte_scan=false |
| 20000 | cte_filtered_mat vs cte_filtered | 6.8580 | 9.2961 | 1.3555 | false | policy=materialized; observed_cte_scan=true |
| 20000 | cte_filtered_notmat vs cte_filtered | 6.8580 | 6.8165 | 0.9940 | true | policy=not_materialized; observed_cte_scan=false |
| 20000 | cte_joined_outer_mat vs cte_joined_outer | 7.4247 | 45.3416 | 6.1069 | false | policy=materialized; observed_cte_scan=true |
| 20000 | cte_joined_outer_notmat vs cte_joined_outer | 7.4247 | 7.1753 | 0.9664 | true | policy=not_materialized; observed_cte_scan=false |
| 100000 | cte_filtered_mat vs cte_filtered | 47.7381 | 58.4158 | 1.2237 | false | policy=materialized; observed_cte_scan=true |
| 100000 | cte_filtered_notmat vs cte_filtered | 47.7381 | 44.3648 | 0.9293 | true | policy=not_materialized; observed_cte_scan=false |
| 100000 | cte_joined_outer_mat vs cte_joined_outer | 49.2412 | 65.4836 | 1.3299 | false | policy=materialized; observed_cte_scan=true |
| 100000 | cte_joined_outer_notmat vs cte_joined_outer | 49.2412 | 48.8755 | 0.9926 | true | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_filtered_mat vs cte_filtered | 78.6819 | 84.1072 | 1.0690 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_filtered_notmat vs cte_filtered | 78.6819 | 77.7130 | 0.9877 | true | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 79.0973 | 90.7410 | 1.1472 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 79.0973 | 44.9244 | 0.5680 | true | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 0.4388 | 1 | `549a07108ab01ba9` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 2.5920 | 1 | `2ce8b85a7b60c88e` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | 2.6486 | 2 | `b307c41eb54cf452` | policy=default; observed_cte_scan=true | valid |
| default | 2 | 100% | 2.4816 | 2 | `b307c41eb54cf452` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 1% | 2.5308 | 4 | `ccfe1af393deae98` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 2.5929 | 4 | `ccfe1af393deae98` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | 2.6456 | 8 | `c8c73778097931dc` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 100% | 3.0538 | 8 | `c8c73778097931dc` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 1% | 3.1171 | 16 | `b0c356323d30c425` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 3.5601 | 16 | `b0c356323d30c425` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 2.4904 | 1 | `eaa7d6058aaa5e03` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 2.7651 | 1 | `eaa7d6058aaa5e03` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | 2.5532 | 2 | `b307c41eb54cf452` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 100% | 2.4868 | 2 | `b307c41eb54cf452` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 1% | 2.8890 | 4 | `ccfe1af393deae98` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 3.0520 | 4 | `ccfe1af393deae98` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | 2.8292 | 8 | `c8c73778097931dc` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 100% | 3.1845 | 8 | `c8c73778097931dc` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 1% | 3.1948 | 16 | `b0c356323d30c425` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 3.4397 | 16 | `b0c356323d30c425` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 0.6732 | 1 | `549a07108ab01ba9` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 2.4696 | 1 | `2ce8b85a7b60c88e` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | 0.6784 | 2 | `28cc913d71d0d72b` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 100% | 4.6611 | 2 | `80657236c094ba42` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 1% | 0.9517 | 4 | `9f850e53560baafb` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 7.5850 | 4 | `ce16f1b9eba58fd7` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | 1.3720 | 8 | `477f7fbaf1c627cb` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 100% | 9.8098 | 8 | `13462faf9230d225` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 1% | 2.3842 | 16 | `1c38b58d8cec08d3` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 15.4788 | 16 | `cfa5e961e551e4af` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 5.6759 | false |
| notmat | 1 | 1% | 1.5344 | true |
| mat | 1 | 100% | 1.0668 | false |
| notmat | 1 | 100% | 0.9528 | true |
| mat | 2 | 1% | 0.9640 | true |
| notmat | 2 | 1% | 0.2561 | false |
| mat | 2 | 100% | 1.0021 | true |
| notmat | 2 | 100% | 1.8783 | false |
| mat | 4 | 1% | 1.1416 | true |
| notmat | 4 | 1% | 0.3761 | false |
| mat | 4 | 100% | 1.1771 | true |
| notmat | 4 | 100% | 2.9253 | false |
| mat | 8 | 1% | 1.0694 | true |
| notmat | 8 | 1% | 0.5186 | false |
| mat | 8 | 100% | 1.0428 | true |
| notmat | 8 | 100% | 3.2123 | false |
| mat | 16 | 1% | 1.0249 | true |
| notmat | 16 | 1% | 0.7649 | false |
| mat | 16 | 100% | 0.9662 | true |
| notmat | 16 | 100% | 4.3478 | false |

## Repeated-design position control

| first median ms | last median ms | ratio (last/first) | same plan |
|---|---|---|---|
| 58.415792 | 58.761190 | 1.0059 | true |

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

- client: nr_periods 377, nr_throttled 307, throttled_usec 16877989 | quota 100000 100000
, memory 2147483648

- pg-single: nr_periods 382, nr_throttled 0, throttled_usec 0 | quota 200000 100000, memory 3221225472

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
| family=bag;stmt=cte_filtered_mat;threshold=100000;position=last | valid | true | true | 0 | 0 | 100263/0 | 0 | 0 |
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
| family=reuse;policy=default;refs=1;consumer=1% | valid | true | false | 0 | 0 | 30/0 | 1983 | 0 |
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
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | true | false | 0 | 0 | 28/0 | 1983 | 0 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | false | false | 0 | 0 | 0/0 | 0 | 0 |
| family=reuse;policy=notmat;refs=2;consumer=1% | valid | true | false | 0 | 0 | 52/0 | 1983 | 0 |
| family=reuse;policy=notmat;refs=2;consumer=100% | valid | false | false | 0 | 0 | 0/0 | 0 | 0 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | true | false | 0 | 0 | 120/0 | 1983 | 0 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | false | false | 0 | 0 | 0/0 | 0 | 0 |
| family=reuse;policy=notmat;refs=8;consumer=1% | valid | true | false | 0 | 0 | 240/0 | 1983 | 0 |
| family=reuse;policy=notmat;refs=8;consumer=100% | valid | false | false | 0 | 0 | 0/0 | 0 | 0 |

## Latency samples and derived throughput

_The reader is sequential, so `queries/s` is 1000/median_ms, a derived rate and not a concurrent-capacity test. min/max are the five sample extremes. Each sample includes any retry and backoff time the statement needed after a transient conflict; per-cell retry counts appear in the access-path table._

| axis | status | min ms | median ms | max ms | samples | derived queries/s |
|---|---|---|---|---|---|---|
| family=bag;stmt=comma_join;threshold=100000 | valid | 44.6134 | 47.7892 | 53.9490 | 5 | 20.9 |
| family=bag;stmt=comma_join;threshold=180000 | valid | 46.4863 | 75.9969 | 81.3764 | 5 | 13.2 |
| family=bag;stmt=comma_join;threshold=20 | valid | 0.5292 | 0.5582 | 0.6746 | 5 | 1791.4 |
| family=bag;stmt=comma_join;threshold=200 | valid | 0.7046 | 0.8816 | 12.6600 | 5 | 1134.3 |
| family=bag;stmt=comma_join;threshold=2000 | valid | 1.7915 | 1.9171 | 1.9769 | 5 | 521.6 |
| family=bag;stmt=comma_join;threshold=20000 | valid | 6.5054 | 6.9819 | 39.3773 | 5 | 143.2 |
| family=bag;stmt=cte_filtered;threshold=100000 | valid | 43.8284 | 47.7381 | 66.2279 | 5 | 20.9 |
| family=bag;stmt=cte_filtered;threshold=180000 | valid | 42.9950 | 78.6819 | 90.6302 | 5 | 12.7 |
| family=bag;stmt=cte_filtered;threshold=20 | valid | 0.3948 | 0.4595 | 0.4860 | 5 | 2176.5 |
| family=bag;stmt=cte_filtered;threshold=200 | valid | 0.9091 | 0.9393 | 0.9965 | 5 | 1064.6 |
| family=bag;stmt=cte_filtered;threshold=2000 | valid | 1.8215 | 1.9683 | 33.6213 | 5 | 508.0 |
| family=bag;stmt=cte_filtered;threshold=20000 | valid | 6.3669 | 6.8580 | 25.9673 | 5 | 145.8 |
| family=bag;stmt=cte_filtered_mat;threshold=100000 | valid | 56.3181 | 58.4158 | 61.2851 | 5 | 17.1 |
| family=bag;stmt=cte_filtered_mat;threshold=100000;position=last | valid | 52.4528 | 58.7612 | 59.4733 | 5 | 17.0 |
| family=bag;stmt=cte_filtered_mat;threshold=180000 | valid | 77.7401 | 84.1072 | 89.8643 | 5 | 11.9 |
| family=bag;stmt=cte_filtered_mat;threshold=20 | valid | 0.5394 | 0.6517 | 0.7387 | 5 | 1534.5 |
| family=bag;stmt=cte_filtered_mat;threshold=200 | valid | 0.7256 | 0.7618 | 0.8256 | 5 | 1312.7 |
| family=bag;stmt=cte_filtered_mat;threshold=2000 | valid | 2.1563 | 2.3361 | 2.4613 | 5 | 428.1 |
| family=bag;stmt=cte_filtered_mat;threshold=20000 | valid | 8.6496 | 9.2961 | 11.4261 | 5 | 107.6 |
| family=bag;stmt=cte_filtered_notmat;threshold=100000 | valid | 41.5588 | 44.3648 | 71.4624 | 5 | 22.5 |
| family=bag;stmt=cte_filtered_notmat;threshold=180000 | valid | 46.5794 | 77.7130 | 86.3005 | 5 | 12.9 |
| family=bag;stmt=cte_filtered_notmat;threshold=20 | valid | 0.4379 | 0.4534 | 0.6399 | 5 | 2205.5 |
| family=bag;stmt=cte_filtered_notmat;threshold=200 | valid | 0.8788 | 0.9332 | 0.9799 | 5 | 1071.6 |
| family=bag;stmt=cte_filtered_notmat;threshold=2000 | valid | 1.5692 | 1.5831 | 2.1563 | 5 | 631.7 |
| family=bag;stmt=cte_filtered_notmat;threshold=20000 | valid | 6.5994 | 6.8165 | 7.2559 | 5 | 146.7 |
| family=bag;stmt=cte_joined_outer;threshold=100000 | valid | 43.0564 | 49.2412 | 53.0021 | 5 | 20.3 |
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 67.0279 | 79.0973 | 86.1319 | 5 | 12.6 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 0.3904 | 0.4514 | 0.4889 | 5 | 2215.1 |
| family=bag;stmt=cte_joined_outer;threshold=200 | valid | 0.6841 | 0.7179 | 0.7756 | 5 | 1393.0 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 1.6489 | 1.7460 | 2.3014 | 5 | 572.7 |
| family=bag;stmt=cte_joined_outer;threshold=20000 | valid | 6.8811 | 7.4247 | 9.0709 | 5 | 134.7 |
| family=bag;stmt=cte_joined_outer_mat;threshold=100000 | valid | 61.1509 | 65.4836 | 67.3350 | 5 | 15.3 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 85.3732 | 90.7410 | 96.0097 | 5 | 11.0 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 42.7821 | 45.6655 | 49.6971 | 5 | 21.9 |
| family=bag;stmt=cte_joined_outer_mat;threshold=200 | valid | 43.6415 | 45.6113 | 52.7216 | 5 | 21.9 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 43.1028 | 43.6970 | 50.3642 | 5 | 22.9 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20000 | valid | 44.9667 | 45.3416 | 48.6147 | 5 | 22.1 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=100000 | valid | 47.5108 | 48.8755 | 52.9759 | 5 | 20.5 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 43.4518 | 44.9244 | 50.7470 | 5 | 22.3 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 0.5614 | 0.6821 | 0.7455 | 5 | 1466.1 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=200 | valid | 0.6396 | 0.7267 | 0.7694 | 5 | 1376.1 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 1.9247 | 2.0480 | 2.1007 | 5 | 488.3 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20000 | valid | 6.9365 | 7.1753 | 7.8710 | 5 | 139.4 |
| family=bag;stmt=filtered_derived;threshold=100000 | valid | 44.1575 | 48.0782 | 51.7881 | 5 | 20.8 |
| family=bag;stmt=filtered_derived;threshold=180000 | valid | 48.8794 | 50.7758 | 61.4313 | 5 | 19.7 |
| family=bag;stmt=filtered_derived;threshold=20 | valid | 0.4247 | 0.4485 | 0.5156 | 5 | 2229.8 |
| family=bag;stmt=filtered_derived;threshold=200 | valid | 0.7843 | 0.8574 | 0.8665 | 5 | 1166.3 |
| family=bag;stmt=filtered_derived;threshold=2000 | valid | 1.6517 | 1.7262 | 3.1829 | 5 | 579.3 |
| family=bag;stmt=filtered_derived;threshold=20000 | valid | 6.3867 | 6.7937 | 7.0458 | 5 | 147.2 |
| family=bag;stmt=join_on;threshold=100000 | valid | 45.3995 | 51.3712 | 57.0365 | 5 | 19.5 |
| family=bag;stmt=join_on;threshold=180000 | valid | 79.6706 | 81.3918 | 98.7522 | 5 | 12.3 |
| family=bag;stmt=join_on;threshold=20 | valid | 0.4901 | 0.6739 | 0.8972 | 5 | 1483.8 |
| family=bag;stmt=join_on;threshold=200 | valid | 0.6131 | 0.7083 | 0.9383 | 5 | 1411.9 |
| family=bag;stmt=join_on;threshold=2000 | valid | 1.6249 | 1.6713 | 1.6962 | 5 | 598.3 |
| family=bag;stmt=join_on;threshold=20000 | valid | 6.8815 | 6.9207 | 8.1233 | 5 | 144.5 |
| family=bag;stmt=join_where;threshold=100000 | valid | 42.9066 | 47.5226 | 73.7127 | 5 | 21.0 |
| family=bag;stmt=join_where;threshold=180000 | valid | 59.2612 | 62.4210 | 67.9099 | 5 | 16.0 |
| family=bag;stmt=join_where;threshold=20 | valid | 0.5764 | 0.5992 | 0.8440 | 5 | 1669.0 |
| family=bag;stmt=join_where;threshold=200 | valid | 0.7181 | 0.7896 | 0.8390 | 5 | 1266.4 |
| family=bag;stmt=join_where;threshold=2000 | valid | 2.0026 | 2.0829 | 2.4447 | 5 | 480.1 |
| family=bag;stmt=join_where;threshold=20000 | valid | 6.7599 | 6.9500 | 16.0900 | 5 | 143.9 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 2.9906 | 3.1171 | 3.4988 | 5 | 320.8 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 3.4385 | 3.5601 | 3.6648 | 5 | 280.9 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 0.3731 | 0.4388 | 0.5281 | 5 | 2279.1 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 1.9529 | 2.5920 | 2.8142 | 5 | 385.8 |
| family=reuse;policy=default;refs=2;consumer=1% | valid | 2.5542 | 2.6486 | 2.8423 | 5 | 377.6 |
| family=reuse;policy=default;refs=2;consumer=100% | valid | 2.4348 | 2.4816 | 2.6363 | 5 | 403.0 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 2.5167 | 2.5308 | 2.9074 | 5 | 395.1 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 2.5429 | 2.5929 | 2.7599 | 5 | 385.7 |
| family=reuse;policy=default;refs=8;consumer=1% | valid | 2.5751 | 2.6456 | 2.8199 | 5 | 378.0 |
| family=reuse;policy=default;refs=8;consumer=100% | valid | 3.0082 | 3.0538 | 3.4561 | 5 | 327.5 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 3.1469 | 3.1948 | 3.2974 | 5 | 313.0 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 3.2817 | 3.4397 | 31.0135 | 5 | 290.7 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 2.4225 | 2.4904 | 2.6208 | 5 | 401.5 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 2.6438 | 2.7651 | 2.9823 | 5 | 361.7 |
| family=reuse;policy=mat;refs=2;consumer=1% | valid | 2.4672 | 2.5532 | 15.2295 | 5 | 391.7 |
| family=reuse;policy=mat;refs=2;consumer=100% | valid | 2.4139 | 2.4868 | 2.9473 | 5 | 402.1 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 2.8340 | 2.8890 | 2.9804 | 5 | 346.1 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 2.8796 | 3.0520 | 3.3957 | 5 | 327.7 |
| family=reuse;policy=mat;refs=8;consumer=1% | valid | 2.5109 | 2.8292 | 3.1610 | 5 | 353.5 |
| family=reuse;policy=mat;refs=8;consumer=100% | valid | 3.0227 | 3.1845 | 3.4157 | 5 | 314.0 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 2.2612 | 2.3842 | 2.6199 | 5 | 419.4 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 15.3792 | 15.4788 | 16.7501 | 5 | 64.6 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 0.6202 | 0.6732 | 0.8357 | 5 | 1485.4 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 2.3875 | 2.4696 | 2.7594 | 5 | 404.9 |
| family=reuse;policy=notmat;refs=2;consumer=1% | valid | 0.5551 | 0.6784 | 0.7079 | 5 | 1474.0 |
| family=reuse;policy=notmat;refs=2;consumer=100% | valid | 4.4078 | 4.6611 | 4.9886 | 5 | 214.5 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 0.9125 | 0.9517 | 1.0175 | 5 | 1050.7 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 6.9796 | 7.5850 | 8.3873 | 5 | 131.8 |
| family=reuse;policy=notmat;refs=8;consumer=1% | valid | 1.2645 | 1.3720 | 2.1241 | 5 | 728.8 |
| family=reuse;policy=notmat;refs=8;consumer=100% | valid | 9.3665 | 9.8098 | 10.1623 | 5 | 101.9 |

## Matched blended workload

- phase: `blended`; target offered rate: 10 updates/s
- offered 349, committed 349, failed 0 over 34.803 s
- offered rate 10.028/s, completed rate 10.028/s
- writer operations per endpoint: [349]
- reader operations overlapping a committed write: 114 of 637
- reader retries of transient conflicts: 0
- per key: 349 keys, min 1, max 1 per key
- ledger: 349 keys reconciled, 0 mismatch
- update latency: avg 1.617 ms, min 0.867, max 17.387 over 349

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42283460`, 5 interleaved rounds
  1. family=bag;stmt=cte_filtered_mat;threshold=100000
  2. family=bag;stmt=join_on;threshold=20
  3. family=bag;stmt=filtered_derived;threshold=20000
  4. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  5. family=bag;stmt=cte_joined_outer_notmat;threshold=100000
  6. family=bag;stmt=cte_joined_outer_mat;threshold=100000
  7. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  8. family=bag;stmt=comma_join;threshold=20
  9. family=bag;stmt=filtered_derived;threshold=200
  10. family=bag;stmt=cte_filtered_notmat;threshold=20000
  11. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  12. family=reuse;policy=notmat;refs=4;consumer=100%
  13. family=bag;stmt=cte_filtered_mat;threshold=200
  14. family=bag;stmt=cte_joined_outer;threshold=200
  15. family=bag;stmt=filtered_derived;threshold=2000
  16. family=bag;stmt=cte_joined_outer;threshold=20
  17. family=bag;stmt=join_on;threshold=2000
  18. family=bag;stmt=join_where;threshold=180000
  19. family=bag;stmt=join_on;threshold=180000
  20. family=reuse;policy=mat;refs=1;consumer=100%
  21. family=reuse;policy=mat;refs=2;consumer=100%
  22. family=bag;stmt=join_where;threshold=100000
  23. family=bag;stmt=cte_filtered_mat;threshold=2000
  24. family=reuse;policy=notmat;refs=16;consumer=100%
  25. family=reuse;policy=mat;refs=8;consumer=100%
  26. family=bag;stmt=cte_joined_outer;threshold=100000
  27. family=reuse;policy=mat;refs=4;consumer=100%
  28. family=reuse;policy=notmat;refs=2;consumer=1%
  29. family=bag;stmt=join_on;threshold=20000
  30. family=bag;stmt=cte_filtered;threshold=2000
  31. family=reuse;policy=mat;refs=8;consumer=1%
  32. family=reuse;policy=default;refs=1;consumer=100%
  33. family=bag;stmt=filtered_derived;threshold=100000
  34. family=bag;stmt=cte_filtered;threshold=180000
  35. family=bag;stmt=cte_filtered;threshold=100000
  36. family=reuse;policy=mat;refs=4;consumer=1%
  37. family=bag;stmt=comma_join;threshold=180000
  38. family=reuse;policy=notmat;refs=1;consumer=1%
  39. family=reuse;policy=mat;refs=1;consumer=1%
  40. family=reuse;policy=notmat;refs=2;consumer=100%
  41. family=bag;stmt=join_on;threshold=100000
  42. family=bag;stmt=cte_joined_outer;threshold=180000
  43. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  44. family=bag;stmt=comma_join;threshold=20000
  45. family=reuse;policy=notmat;refs=4;consumer=1%
  46. family=bag;stmt=join_where;threshold=200
  47. family=bag;stmt=cte_joined_outer;threshold=2000
  48. family=reuse;policy=mat;refs=2;consumer=1%
  49. family=bag;stmt=cte_filtered_notmat;threshold=20
  50. family=reuse;policy=notmat;refs=8;consumer=1%
  51. family=bag;stmt=join_on;threshold=200
  52. family=reuse;policy=default;refs=16;consumer=1%
  53. family=bag;stmt=cte_joined_outer_notmat;threshold=20000
  54. family=reuse;policy=default;refs=2;consumer=1%
  55. family=bag;stmt=cte_filtered_mat;threshold=20000
  56. family=reuse;policy=notmat;refs=16;consumer=1%
  57. family=reuse;policy=default;refs=8;consumer=100%
  58. family=bag;stmt=comma_join;threshold=200
  59. family=reuse;policy=default;refs=4;consumer=1%
  60. family=bag;stmt=cte_joined_outer_notmat;threshold=200
  61. family=bag;stmt=join_where;threshold=20000
  62. family=bag;stmt=comma_join;threshold=2000
  63. family=bag;stmt=filtered_derived;threshold=20
  64. family=reuse;policy=default;refs=16;consumer=100%
  65. family=reuse;policy=default;refs=4;consumer=100%
  66. family=reuse;policy=default;refs=1;consumer=1%
  67. family=reuse;policy=notmat;refs=8;consumer=100%
  68. family=bag;stmt=cte_joined_outer_mat;threshold=200
  69. family=bag;stmt=filtered_derived;threshold=180000
  70. family=bag;stmt=cte_filtered_mat;threshold=180000
  71. family=bag;stmt=cte_joined_outer_mat;threshold=20000
  72. family=reuse;policy=notmat;refs=1;consumer=100%
  73. family=bag;stmt=cte_joined_outer_mat;threshold=20
  74. family=bag;stmt=join_where;threshold=2000
  75. family=reuse;policy=default;refs=8;consumer=1%
  76. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  77. family=bag;stmt=cte_filtered_notmat;threshold=100000
  78. family=bag;stmt=cte_filtered;threshold=200
  79. family=reuse;policy=default;refs=2;consumer=100%
  80. family=reuse;policy=mat;refs=16;consumer=100%
  81. family=bag;stmt=cte_filtered_notmat;threshold=2000
  82. family=bag;stmt=cte_filtered;threshold=20000
  83. family=bag;stmt=cte_filtered_mat;threshold=20
  84. family=bag;stmt=comma_join;threshold=100000
  85. family=bag;stmt=cte_joined_outer;threshold=20000
  86. family=bag;stmt=cte_filtered_notmat;threshold=200
  87. family=bag;stmt=cte_filtered_notmat;threshold=180000
  88. family=bag;stmt=join_where;threshold=20
  89. family=reuse;policy=mat;refs=16;consumer=1%
  90. family=bag;stmt=cte_filtered;threshold=20
  91. family=bag;stmt=cte_filtered_mat;threshold=100000;position=last

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
| family=bag;stmt=cte_filtered_mat;threshold=100000;position=last | j04m_cte_filtered_mat | valid | `4b4ce4c83ae7fcc8` | policy=materialized; observed_cte_scan=true |
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
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `cfa5e961e551e4af` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `549a07108ab01ba9` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `2ce8b85a7b60c88e` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=2;consumer=1% | r_notmat_n2_c1 | valid | `28cc913d71d0d72b` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=2;consumer=100% | r_notmat_n2_c100 | valid | `80657236c094ba42` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `9f850e53560baafb` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `ce16f1b9eba58fd7` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=8;consumer=1% | r_notmat_n8_c1 | valid | `477f7fbaf1c627cb` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=8;consumer=100% | r_notmat_n8_c100 | valid | `13462faf9230d225` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, VERBOSE, BUFFERS, FORMAT TEXT` supported on postgres
- one host: this run used node budget 2 CPU 3g and client budget 1 CPU 2g, so budget comparisons are separate classes here, not independent hardware
- one client process with one fixed maintenance writer at 10 updates/s; readers are sequential; the reported sample latency includes retry and backoff time
- physical rank-index arm pg-btree-rank; an engine whose planner ignores the index is a different physical design, not a faster one
- distributed RPC counters are N/A on single-node PostgreSQL; shared buffers and rows removed are parsed per cell where printed
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
