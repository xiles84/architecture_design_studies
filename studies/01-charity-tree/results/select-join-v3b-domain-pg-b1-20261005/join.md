# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 91 cells, 91 valid, engine `postgres`, topology `single`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-v3b-domain-pg-b1-20261005`
- run tag: `run/study-01/select-join-v3b-domain-pg-b1-20261005`
- environment: `host-zenbook-ux5406sa`
- repo commit: `7243fee074bfb16fc5ad8b6afc1a2d7549080935`
- engine version: `PostgreSQL 17.11 (Debian 17.11-1.pgdg13+2) on x86_64-pc-linux-gnu, compiled by gcc (Debian 14.2.0-19) 14.2.0, 64-bit`
- database engine image: `docker.io/library/postgres:17.11` @ `7296f210ae81031ec955dbad9a67a84fe958572a2153b8d0826a647522904dc1`
- benchmark client image: `localhost/charitybench-select:select-join-v3b-domain-pg-b1-20261005` @ `f3f1e10c44cbcb4c7b5bc56ed95a4f071682b6f7d33b5e0a7b688c4271eeef88`
- inputs digest: `4ad32c4a1ab1d6ca5968ee0a7bb75f9580e86aa1e6ec5df69ce60d8ec0dd3f8c`
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
| 20 (0.01%) | 0.4642 | 0.4160 | 0.5129 | 0.5001 | 0.6091 | 0.5153 | 0.4750 | 47.0695 | 0.5622 | 0.5194 |
| 200 (0.1%) | 0.7512 | 0.8461 | 0.7617 | 0.9974 | 0.7521 | 1.0120 | 0.9248 | 47.4744 | 0.7677 | 0.7180 |
| 2000 (1%) | 1.8015 | 1.9423 | 1.7229 | 2.0456 | 2.0084 | 1.8422 | 1.8930 | 43.7699 | 1.9513 | 1.6872 |
| 20000 (10%) | 7.6169 | 7.1561 | 7.2271 | 7.2616 | 9.4334 | 7.0651 | 7.2137 | 52.9223 | 7.5818 | 7.5551 |
| 100000 (50%) | 61.1514 | 46.3995 | 45.9740 | 50.4568 | 61.8963 | 49.3899 | 45.6567 | 63.2256 | 46.6147 | 54.8663 |
| 180000 (90%) | 85.1346 | 64.7109 | 48.5609 | 50.1551 | 69.6185 | 82.0364 | 64.5399 | 93.0220 | 79.4112 | 78.4841 |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | join_on vs join_where | 0.4642 | 0.4160 | 0.8963 | true |
| 20 | filtered_derived vs join_where | 0.4642 | 0.5129 | 1.1049 | true |
| 20 | cte_filtered vs join_where | 0.4642 | 0.5001 | 1.0773 | true |
| 20 | cte_filtered_mat vs join_where | 0.4642 | 0.6091 | 1.3121 | false |
| 20 | cte_filtered_notmat vs join_where | 0.4642 | 0.5153 | 1.1101 | true |
| 20 | cte_joined_outer vs join_where | 0.4642 | 0.4750 | 1.0233 | true |
| 20 | cte_joined_outer_mat vs join_where | 0.4642 | 47.0695 | 101.4022 | false |
| 20 | cte_joined_outer_notmat vs join_where | 0.4642 | 0.5622 | 1.2111 | true |
| 20 | comma_join vs join_where | 0.4642 | 0.5194 | 1.1189 | true |
| 200 | join_on vs join_where | 0.7512 | 0.8461 | 1.1263 | true |
| 200 | filtered_derived vs join_where | 0.7512 | 0.7617 | 1.0139 | true |
| 200 | cte_filtered vs join_where | 0.7512 | 0.9974 | 1.3277 | true |
| 200 | cte_filtered_mat vs join_where | 0.7512 | 0.7521 | 1.0012 | false |
| 200 | cte_filtered_notmat vs join_where | 0.7512 | 1.0120 | 1.3471 | true |
| 200 | cte_joined_outer vs join_where | 0.7512 | 0.9248 | 1.2310 | true |
| 200 | cte_joined_outer_mat vs join_where | 0.7512 | 47.4744 | 63.1960 | false |
| 200 | cte_joined_outer_notmat vs join_where | 0.7512 | 0.7677 | 1.0220 | true |
| 200 | comma_join vs join_where | 0.7512 | 0.7180 | 0.9558 | true |
| 2000 | join_on vs join_where | 1.8015 | 1.9423 | 1.0781 | true |
| 2000 | filtered_derived vs join_where | 1.8015 | 1.7229 | 0.9563 | true |
| 2000 | cte_filtered vs join_where | 1.8015 | 2.0456 | 1.1355 | true |
| 2000 | cte_filtered_mat vs join_where | 1.8015 | 2.0084 | 1.1148 | false |
| 2000 | cte_filtered_notmat vs join_where | 1.8015 | 1.8422 | 1.0226 | true |
| 2000 | cte_joined_outer vs join_where | 1.8015 | 1.8930 | 1.0508 | true |
| 2000 | cte_joined_outer_mat vs join_where | 1.8015 | 43.7699 | 24.2961 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 1.8015 | 1.9513 | 1.0831 | true |
| 2000 | comma_join vs join_where | 1.8015 | 1.6872 | 0.9365 | true |
| 20000 | join_on vs join_where | 7.6169 | 7.1561 | 0.9395 | true |
| 20000 | filtered_derived vs join_where | 7.6169 | 7.2271 | 0.9488 | true |
| 20000 | cte_filtered vs join_where | 7.6169 | 7.2616 | 0.9534 | true |
| 20000 | cte_filtered_mat vs join_where | 7.6169 | 9.4334 | 1.2385 | false |
| 20000 | cte_filtered_notmat vs join_where | 7.6169 | 7.0651 | 0.9276 | true |
| 20000 | cte_joined_outer vs join_where | 7.6169 | 7.2137 | 0.9471 | true |
| 20000 | cte_joined_outer_mat vs join_where | 7.6169 | 52.9223 | 6.9480 | false |
| 20000 | cte_joined_outer_notmat vs join_where | 7.6169 | 7.5818 | 0.9954 | true |
| 20000 | comma_join vs join_where | 7.6169 | 7.5551 | 0.9919 | true |
| 100000 | join_on vs join_where | 61.1514 | 46.3995 | 0.7588 | true |
| 100000 | filtered_derived vs join_where | 61.1514 | 45.9740 | 0.7518 | true |
| 100000 | cte_filtered vs join_where | 61.1514 | 50.4568 | 0.8251 | true |
| 100000 | cte_filtered_mat vs join_where | 61.1514 | 61.8963 | 1.0122 | false |
| 100000 | cte_filtered_notmat vs join_where | 61.1514 | 49.3899 | 0.8077 | true |
| 100000 | cte_joined_outer vs join_where | 61.1514 | 45.6567 | 0.7466 | true |
| 100000 | cte_joined_outer_mat vs join_where | 61.1514 | 63.2256 | 1.0339 | false |
| 100000 | cte_joined_outer_notmat vs join_where | 61.1514 | 46.6147 | 0.7623 | true |
| 100000 | comma_join vs join_where | 61.1514 | 54.8663 | 0.8972 | true |
| 180000 | join_on vs join_where | 85.1346 | 64.7109 | 0.7601 | true |
| 180000 | filtered_derived vs join_where | 85.1346 | 48.5609 | 0.5704 | true |
| 180000 | cte_filtered vs join_where | 85.1346 | 50.1551 | 0.5891 | true |
| 180000 | cte_filtered_mat vs join_where | 85.1346 | 69.6185 | 0.8177 | false |
| 180000 | cte_filtered_notmat vs join_where | 85.1346 | 82.0364 | 0.9636 | true |
| 180000 | cte_joined_outer vs join_where | 85.1346 | 64.5399 | 0.7581 | true |
| 180000 | cte_joined_outer_mat vs join_where | 85.1346 | 93.0220 | 1.0926 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 85.1346 | 79.4112 | 0.9328 | true |
| 180000 | comma_join vs join_where | 85.1346 | 78.4841 | 0.9219 | true |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_filtered_mat vs cte_filtered | 0.5001 | 0.6091 | 1.2179 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_filtered_notmat vs cte_filtered | 0.5001 | 0.5153 | 1.0304 | true | policy=not_materialized; observed_cte_scan=false |
| 20 | cte_joined_outer_mat vs cte_joined_outer | 0.4750 | 47.0695 | 99.0916 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 0.4750 | 0.5622 | 1.1835 | true | policy=not_materialized; observed_cte_scan=false |
| 200 | cte_filtered_mat vs cte_filtered | 0.9974 | 0.7521 | 0.7541 | false | policy=materialized; observed_cte_scan=true |
| 200 | cte_filtered_notmat vs cte_filtered | 0.9974 | 1.0120 | 1.0146 | true | policy=not_materialized; observed_cte_scan=false |
| 200 | cte_joined_outer_mat vs cte_joined_outer | 0.9248 | 47.4744 | 51.3356 | false | policy=materialized; observed_cte_scan=true |
| 200 | cte_joined_outer_notmat vs cte_joined_outer | 0.9248 | 0.7677 | 0.8302 | true | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_filtered_mat vs cte_filtered | 2.0456 | 2.0084 | 0.9818 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_filtered_notmat vs cte_filtered | 2.0456 | 1.8422 | 0.9006 | true | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 1.8930 | 43.7699 | 23.1225 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 1.8930 | 1.9513 | 1.0308 | true | policy=not_materialized; observed_cte_scan=false |
| 20000 | cte_filtered_mat vs cte_filtered | 7.2616 | 9.4334 | 1.2991 | false | policy=materialized; observed_cte_scan=true |
| 20000 | cte_filtered_notmat vs cte_filtered | 7.2616 | 7.0651 | 0.9729 | true | policy=not_materialized; observed_cte_scan=false |
| 20000 | cte_joined_outer_mat vs cte_joined_outer | 7.2137 | 52.9223 | 7.3363 | false | policy=materialized; observed_cte_scan=true |
| 20000 | cte_joined_outer_notmat vs cte_joined_outer | 7.2137 | 7.5818 | 1.0510 | true | policy=not_materialized; observed_cte_scan=false |
| 100000 | cte_filtered_mat vs cte_filtered | 50.4568 | 61.8963 | 1.2267 | false | policy=materialized; observed_cte_scan=true |
| 100000 | cte_filtered_notmat vs cte_filtered | 50.4568 | 49.3899 | 0.9789 | true | policy=not_materialized; observed_cte_scan=false |
| 100000 | cte_joined_outer_mat vs cte_joined_outer | 45.6567 | 63.2256 | 1.3848 | false | policy=materialized; observed_cte_scan=true |
| 100000 | cte_joined_outer_notmat vs cte_joined_outer | 45.6567 | 46.6147 | 1.0210 | true | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_filtered_mat vs cte_filtered | 50.1551 | 69.6185 | 1.3881 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_filtered_notmat vs cte_filtered | 50.1551 | 82.0364 | 1.6357 | true | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 64.5399 | 93.0220 | 1.4413 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 64.5399 | 79.4112 | 1.2304 | true | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 0.6132 | 1 | `549a07108ab01ba9` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 2.3676 | 1 | `2ce8b85a7b60c88e` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | 2.7435 | 2 | `b307c41eb54cf452` | policy=default; observed_cte_scan=true | valid |
| default | 2 | 100% | 2.6144 | 2 | `b307c41eb54cf452` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 1% | 2.7322 | 4 | `ccfe1af393deae98` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 2.7126 | 4 | `ccfe1af393deae98` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | 2.4437 | 8 | `c8c73778097931dc` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 100% | 3.2566 | 8 | `c8c73778097931dc` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 1% | 3.0162 | 16 | `b0c356323d30c425` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 3.7032 | 16 | `b0c356323d30c425` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 2.6243 | 1 | `eaa7d6058aaa5e03` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 2.5801 | 1 | `eaa7d6058aaa5e03` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | 2.6620 | 2 | `b307c41eb54cf452` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 100% | 2.7935 | 2 | `b307c41eb54cf452` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 1% | 2.5024 | 4 | `ccfe1af393deae98` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 2.8295 | 4 | `ccfe1af393deae98` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | 3.0563 | 8 | `c8c73778097931dc` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 100% | 3.1890 | 8 | `c8c73778097931dc` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 1% | 3.3480 | 16 | `b0c356323d30c425` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 3.4463 | 16 | `b0c356323d30c425` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 0.4883 | 1 | `549a07108ab01ba9` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 2.5503 | 1 | `2ce8b85a7b60c88e` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | 0.7900 | 2 | `28cc913d71d0d72b` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 100% | 4.8823 | 2 | `164db6b1595310eb` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 1% | 1.0573 | 4 | `9f850e53560baafb` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 7.5347 | 4 | `96aa76348edaa095` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | 1.4932 | 8 | `477f7fbaf1c627cb` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 100% | 9.3082 | 8 | `ab83914f43e32200` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 1% | 2.2999 | 16 | `1c38b58d8cec08d3` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 15.8456 | 16 | `957d40747f6f7eb1` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 4.2799 | false |
| notmat | 1 | 1% | 0.7964 | true |
| mat | 1 | 100% | 1.0898 | false |
| notmat | 1 | 100% | 1.0772 | true |
| mat | 2 | 1% | 0.9703 | true |
| notmat | 2 | 1% | 0.2880 | false |
| mat | 2 | 100% | 1.0685 | true |
| notmat | 2 | 100% | 1.8675 | false |
| mat | 4 | 1% | 0.9159 | true |
| notmat | 4 | 1% | 0.3870 | false |
| mat | 4 | 100% | 1.0431 | true |
| notmat | 4 | 100% | 2.7777 | false |
| mat | 8 | 1% | 1.2507 | true |
| notmat | 8 | 1% | 0.6110 | false |
| mat | 8 | 100% | 0.9793 | true |
| notmat | 8 | 100% | 2.8583 | false |
| mat | 16 | 1% | 1.1100 | true |
| notmat | 16 | 1% | 0.7625 | false |
| mat | 16 | 100% | 0.9306 | true |
| notmat | 16 | 100% | 4.2788 | false |

## Repeated-design position control

| first median ms | last median ms | ratio (last/first) | same plan |
|---|---|---|---|
| 9.433355 | 9.394800 | 0.9959 | false |

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

- client: nr_periods 364, nr_throttled 287, throttled_usec 16046843 | quota 100000 100000
, memory 2147483648

- pg-single: nr_periods 369, nr_throttled 0, throttled_usec 0 | quota 200000 100000, memory 3221225472

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
| family=bag;stmt=cte_filtered_mat;threshold=20000;position=last | valid | false | true | 0 | 0 | 1351/0 | 0 | 0 |
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
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | true | false | 0 | 0 | 384/0 | 1983 | 0 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | false | false | 0 | 0 | 0/0 | 0 | 0 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | true | false | 0 | 0 | 26/0 | 1983 | 0 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | false | false | 0 | 0 | 0/0 | 0 | 0 |
| family=reuse;policy=notmat;refs=2;consumer=1% | valid | true | false | 0 | 0 | 64/0 | 1983 | 0 |
| family=reuse;policy=notmat;refs=2;consumer=100% | valid | false | false | 0 | 0 | 0/0 | 0 | 0 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | true | false | 0 | 0 | 104/0 | 1983 | 0 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | false | false | 0 | 0 | 0/0 | 0 | 0 |
| family=reuse;policy=notmat;refs=8;consumer=1% | valid | true | false | 0 | 0 | 176/0 | 1983 | 0 |
| family=reuse;policy=notmat;refs=8;consumer=100% | valid | false | false | 0 | 0 | 0/0 | 0 | 0 |

## Latency samples and derived throughput

_The reader is sequential, so `queries/s` is 1000/median_ms, a derived rate and not a concurrent-capacity test. min/max are the five sample extremes. Each sample includes any retry and backoff time the statement needed after a transient conflict; per-cell retry counts appear in the access-path table._

| axis | status | min ms | median ms | max ms | samples | derived queries/s |
|---|---|---|---|---|---|---|
| family=bag;stmt=comma_join;threshold=100000 | valid | 47.0140 | 54.8663 | 63.2076 | 5 | 18.2 |
| family=bag;stmt=comma_join;threshold=180000 | valid | 68.4158 | 78.4841 | 80.4692 | 5 | 12.7 |
| family=bag;stmt=comma_join;threshold=20 | valid | 0.4479 | 0.5194 | 0.5376 | 5 | 1925.4 |
| family=bag;stmt=comma_join;threshold=200 | valid | 0.6872 | 0.7180 | 0.7722 | 5 | 1392.8 |
| family=bag;stmt=comma_join;threshold=2000 | valid | 1.6121 | 1.6872 | 2.2191 | 5 | 592.7 |
| family=bag;stmt=comma_join;threshold=20000 | valid | 6.5418 | 7.5551 | 8.5896 | 5 | 132.4 |
| family=bag;stmt=cte_filtered;threshold=100000 | valid | 42.2224 | 50.4568 | 70.3370 | 5 | 19.8 |
| family=bag;stmt=cte_filtered;threshold=180000 | valid | 48.4641 | 50.1551 | 90.6545 | 5 | 19.9 |
| family=bag;stmt=cte_filtered;threshold=20 | valid | 0.4127 | 0.5001 | 0.6147 | 5 | 1999.7 |
| family=bag;stmt=cte_filtered;threshold=200 | valid | 0.8343 | 0.9974 | 1.0260 | 5 | 1002.6 |
| family=bag;stmt=cte_filtered;threshold=2000 | valid | 1.8647 | 2.0456 | 2.6905 | 5 | 488.9 |
| family=bag;stmt=cte_filtered;threshold=20000 | valid | 6.9747 | 7.2616 | 9.3657 | 5 | 137.7 |
| family=bag;stmt=cte_filtered_mat;threshold=100000 | valid | 61.0293 | 61.8963 | 63.2174 | 5 | 16.2 |
| family=bag;stmt=cte_filtered_mat;threshold=180000 | valid | 64.2604 | 69.6185 | 95.9669 | 5 | 14.4 |
| family=bag;stmt=cte_filtered_mat;threshold=20 | valid | 0.5564 | 0.6091 | 0.9004 | 5 | 1641.9 |
| family=bag;stmt=cte_filtered_mat;threshold=200 | valid | 0.6876 | 0.7521 | 1.9329 | 5 | 1329.5 |
| family=bag;stmt=cte_filtered_mat;threshold=2000 | valid | 1.7122 | 2.0084 | 2.4693 | 5 | 497.9 |
| family=bag;stmt=cte_filtered_mat;threshold=20000 | valid | 9.1019 | 9.4334 | 9.8674 | 5 | 106.0 |
| family=bag;stmt=cte_filtered_mat;threshold=20000;position=last | valid | 8.9230 | 9.3948 | 10.3481 | 5 | 106.4 |
| family=bag;stmt=cte_filtered_notmat;threshold=100000 | valid | 42.5658 | 49.3899 | 64.0379 | 5 | 20.2 |
| family=bag;stmt=cte_filtered_notmat;threshold=180000 | valid | 43.4766 | 82.0364 | 90.3376 | 5 | 12.2 |
| family=bag;stmt=cte_filtered_notmat;threshold=20 | valid | 0.3980 | 0.5153 | 0.7775 | 5 | 1940.7 |
| family=bag;stmt=cte_filtered_notmat;threshold=200 | valid | 0.9427 | 1.0120 | 1.0245 | 5 | 988.2 |
| family=bag;stmt=cte_filtered_notmat;threshold=2000 | valid | 1.7560 | 1.8422 | 1.9304 | 5 | 542.8 |
| family=bag;stmt=cte_filtered_notmat;threshold=20000 | valid | 6.5289 | 7.0651 | 7.5256 | 5 | 141.5 |
| family=bag;stmt=cte_joined_outer;threshold=100000 | valid | 42.8945 | 45.6567 | 51.1494 | 5 | 21.9 |
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 58.8677 | 64.5399 | 72.5403 | 5 | 15.5 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 0.4420 | 0.4750 | 0.5258 | 5 | 2105.2 |
| family=bag;stmt=cte_joined_outer;threshold=200 | valid | 0.8936 | 0.9248 | 1.0065 | 5 | 1081.3 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 1.8184 | 1.8930 | 2.1213 | 5 | 528.3 |
| family=bag;stmt=cte_joined_outer;threshold=20000 | valid | 6.7940 | 7.2137 | 7.3047 | 5 | 138.6 |
| family=bag;stmt=cte_joined_outer_mat;threshold=100000 | valid | 59.6730 | 63.2256 | 65.0788 | 5 | 15.8 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 83.8005 | 93.0220 | 94.2637 | 5 | 10.8 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 45.1865 | 47.0695 | 49.7085 | 5 | 21.2 |
| family=bag;stmt=cte_joined_outer_mat;threshold=200 | valid | 42.6116 | 47.4744 | 54.7391 | 5 | 21.1 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 40.7670 | 43.7699 | 44.0077 | 5 | 22.8 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20000 | valid | 47.6480 | 52.9223 | 54.4257 | 5 | 18.9 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=100000 | valid | 45.5263 | 46.6147 | 65.8085 | 5 | 21.5 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 44.2776 | 79.4112 | 86.4963 | 5 | 12.6 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 0.4856 | 0.5622 | 0.7738 | 5 | 1778.8 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=200 | valid | 0.6981 | 0.7677 | 0.8872 | 5 | 1302.6 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 1.8406 | 1.9513 | 2.2051 | 5 | 512.5 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20000 | valid | 6.6922 | 7.5818 | 9.6909 | 5 | 131.9 |
| family=bag;stmt=filtered_derived;threshold=100000 | valid | 45.0639 | 45.9740 | 47.7011 | 5 | 21.8 |
| family=bag;stmt=filtered_derived;threshold=180000 | valid | 48.2626 | 48.5609 | 97.5586 | 5 | 20.6 |
| family=bag;stmt=filtered_derived;threshold=20 | valid | 0.4955 | 0.5129 | 0.5613 | 5 | 1949.9 |
| family=bag;stmt=filtered_derived;threshold=200 | valid | 0.6820 | 0.7617 | 0.8334 | 5 | 1312.9 |
| family=bag;stmt=filtered_derived;threshold=2000 | valid | 1.5279 | 1.7229 | 2.0544 | 5 | 580.4 |
| family=bag;stmt=filtered_derived;threshold=20000 | valid | 6.7203 | 7.2271 | 7.4801 | 5 | 138.4 |
| family=bag;stmt=join_on;threshold=100000 | valid | 42.9079 | 46.3995 | 64.8418 | 5 | 21.6 |
| family=bag;stmt=join_on;threshold=180000 | valid | 45.4337 | 64.7109 | 80.5530 | 5 | 15.5 |
| family=bag;stmt=join_on;threshold=20 | valid | 0.3880 | 0.4160 | 0.4357 | 5 | 2403.6 |
| family=bag;stmt=join_on;threshold=200 | valid | 0.8252 | 0.8461 | 0.9543 | 5 | 1181.9 |
| family=bag;stmt=join_on;threshold=2000 | valid | 1.8253 | 1.9423 | 2.7088 | 5 | 514.9 |
| family=bag;stmt=join_on;threshold=20000 | valid | 6.9364 | 7.1561 | 35.1758 | 5 | 139.7 |
| family=bag;stmt=join_where;threshold=100000 | valid | 40.8640 | 61.1514 | 76.1517 | 5 | 16.4 |
| family=bag;stmt=join_where;threshold=180000 | valid | 68.9285 | 85.1346 | 89.9839 | 5 | 11.7 |
| family=bag;stmt=join_where;threshold=20 | valid | 0.3956 | 0.4642 | 0.5967 | 5 | 2154.3 |
| family=bag;stmt=join_where;threshold=200 | valid | 0.7009 | 0.7512 | 0.8494 | 5 | 1331.2 |
| family=bag;stmt=join_where;threshold=2000 | valid | 1.6355 | 1.8015 | 1.9276 | 5 | 555.1 |
| family=bag;stmt=join_where;threshold=20000 | valid | 6.8121 | 7.6169 | 8.5355 | 5 | 131.3 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 2.8717 | 3.0162 | 3.2170 | 5 | 331.5 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 3.3815 | 3.7032 | 4.1970 | 5 | 270.0 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 0.6056 | 0.6132 | 0.7012 | 5 | 1630.9 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 2.1942 | 2.3676 | 2.4570 | 5 | 422.4 |
| family=reuse;policy=default;refs=2;consumer=1% | valid | 2.6314 | 2.7435 | 2.8978 | 5 | 364.5 |
| family=reuse;policy=default;refs=2;consumer=100% | valid | 2.4084 | 2.6144 | 2.7035 | 5 | 382.5 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 2.5760 | 2.7322 | 2.9821 | 5 | 366.0 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 2.6414 | 2.7126 | 3.7339 | 5 | 368.6 |
| family=reuse;policy=default;refs=8;consumer=1% | valid | 2.3939 | 2.4437 | 3.0122 | 5 | 409.2 |
| family=reuse;policy=default;refs=8;consumer=100% | valid | 3.1619 | 3.2566 | 4.6141 | 5 | 307.1 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 3.2517 | 3.3480 | 4.2060 | 5 | 298.7 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 3.2061 | 3.4463 | 3.6106 | 5 | 290.2 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 2.5593 | 2.6243 | 2.7012 | 5 | 381.1 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 2.1059 | 2.5801 | 3.1625 | 5 | 387.6 |
| family=reuse;policy=mat;refs=2;consumer=1% | valid | 2.5637 | 2.6620 | 3.4726 | 5 | 375.7 |
| family=reuse;policy=mat;refs=2;consumer=100% | valid | 2.5132 | 2.7935 | 17.0592 | 5 | 358.0 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 2.4495 | 2.5024 | 4.0323 | 5 | 399.6 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 2.5552 | 2.8295 | 3.0595 | 5 | 353.4 |
| family=reuse;policy=mat;refs=8;consumer=1% | valid | 2.8447 | 3.0563 | 3.8416 | 5 | 327.2 |
| family=reuse;policy=mat;refs=8;consumer=100% | valid | 3.1186 | 3.1890 | 3.8723 | 5 | 313.6 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 2.2262 | 2.2999 | 2.5041 | 5 | 434.8 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 14.8126 | 15.8456 | 16.3138 | 5 | 63.1 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 0.4064 | 0.4883 | 0.6866 | 5 | 2047.8 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 2.4664 | 2.5503 | 3.4351 | 5 | 392.1 |
| family=reuse;policy=notmat;refs=2;consumer=1% | valid | 0.7174 | 0.7900 | 0.8596 | 5 | 1265.8 |
| family=reuse;policy=notmat;refs=2;consumer=100% | valid | 4.5216 | 4.8823 | 5.4747 | 5 | 204.8 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 1.0290 | 1.0573 | 1.1305 | 5 | 945.8 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 7.3780 | 7.5347 | 7.7602 | 5 | 132.7 |
| family=reuse;policy=notmat;refs=8;consumer=1% | valid | 1.4413 | 1.4932 | 1.8875 | 5 | 669.7 |
| family=reuse;policy=notmat;refs=8;consumer=100% | valid | 8.7395 | 9.3082 | 10.9141 | 5 | 107.4 |

## Matched blended workload

- phase: `blended`; target offered rate: 10 updates/s
- offered 335, committed 335, failed 0 over 33.473 s
- offered rate 10.008/s, completed rate 10.008/s
- writer operations per endpoint: [335]
- reader operations overlapping a committed write: 93 of 637
- reader retries of transient conflicts: 0
- per key: 335 keys, min 1, max 1 per key
- ledger: 335 keys reconciled, 0 mismatch
- update latency: avg 1.564 ms, min 0.909, max 10.899 over 335

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42746858`, 5 interleaved rounds
  1. family=bag;stmt=cte_filtered_mat;threshold=20000
  2. family=reuse;policy=notmat;refs=8;consumer=1%
  3. family=bag;stmt=cte_joined_outer_notmat;threshold=100000
  4. family=reuse;policy=notmat;refs=1;consumer=100%
  5. family=bag;stmt=join_where;threshold=2000
  6. family=bag;stmt=cte_filtered_notmat;threshold=20000
  7. family=bag;stmt=join_on;threshold=2000
  8. family=bag;stmt=comma_join;threshold=100000
  9. family=reuse;policy=mat;refs=1;consumer=1%
  10. family=reuse;policy=default;refs=16;consumer=1%
  11. family=reuse;policy=notmat;refs=16;consumer=1%
  12. family=bag;stmt=comma_join;threshold=2000
  13. family=bag;stmt=comma_join;threshold=200
  14. family=bag;stmt=cte_joined_outer;threshold=180000
  15. family=reuse;policy=default;refs=2;consumer=1%
  16. family=reuse;policy=default;refs=8;consumer=1%
  17. family=reuse;policy=mat;refs=1;consumer=100%
  18. family=bag;stmt=cte_joined_outer_mat;threshold=100000
  19. family=bag;stmt=cte_filtered_notmat;threshold=200
  20. family=bag;stmt=join_where;threshold=20
  21. family=reuse;policy=mat;refs=2;consumer=100%
  22. family=bag;stmt=join_on;threshold=180000
  23. family=bag;stmt=cte_filtered;threshold=200
  24. family=reuse;policy=notmat;refs=1;consumer=1%
  25. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  26. family=reuse;policy=mat;refs=8;consumer=100%
  27. family=bag;stmt=comma_join;threshold=20
  28. family=reuse;policy=notmat;refs=8;consumer=100%
  29. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  30. family=reuse;policy=default;refs=4;consumer=1%
  31. family=bag;stmt=join_on;threshold=20
  32. family=reuse;policy=default;refs=2;consumer=100%
  33. family=reuse;policy=mat;refs=4;consumer=1%
  34. family=bag;stmt=cte_filtered_mat;threshold=2000
  35. family=bag;stmt=filtered_derived;threshold=2000
  36. family=bag;stmt=cte_joined_outer;threshold=100000
  37. family=bag;stmt=join_where;threshold=180000
  38. family=reuse;policy=notmat;refs=4;consumer=1%
  39. family=bag;stmt=cte_filtered_notmat;threshold=100000
  40. family=bag;stmt=cte_filtered;threshold=100000
  41. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  42. family=reuse;policy=notmat;refs=2;consumer=100%
  43. family=bag;stmt=join_where;threshold=100000
  44. family=bag;stmt=cte_joined_outer;threshold=20000
  45. family=bag;stmt=cte_joined_outer_mat;threshold=20000
  46. family=bag;stmt=filtered_derived;threshold=100000
  47. family=bag;stmt=comma_join;threshold=20000
  48. family=bag;stmt=cte_filtered_mat;threshold=20
  49. family=reuse;policy=mat;refs=2;consumer=1%
  50. family=bag;stmt=cte_joined_outer_notmat;threshold=200
  51. family=bag;stmt=comma_join;threshold=180000
  52. family=bag;stmt=cte_filtered_notmat;threshold=180000
  53. family=reuse;policy=mat;refs=8;consumer=1%
  54. family=bag;stmt=join_on;threshold=100000
  55. family=reuse;policy=notmat;refs=4;consumer=100%
  56. family=bag;stmt=cte_joined_outer_mat;threshold=20
  57. family=bag;stmt=join_on;threshold=200
  58. family=reuse;policy=default;refs=4;consumer=100%
  59. family=reuse;policy=mat;refs=16;consumer=100%
  60. family=bag;stmt=cte_filtered;threshold=20000
  61. family=bag;stmt=cte_joined_outer_notmat;threshold=20000
  62. family=reuse;policy=default;refs=1;consumer=1%
  63. family=reuse;policy=default;refs=1;consumer=100%
  64. family=bag;stmt=cte_filtered_mat;threshold=200
  65. family=bag;stmt=cte_joined_outer_mat;threshold=200
  66. family=reuse;policy=notmat;refs=16;consumer=100%
  67. family=bag;stmt=cte_filtered_notmat;threshold=2000
  68. family=bag;stmt=cte_filtered_notmat;threshold=20
  69. family=bag;stmt=cte_filtered;threshold=180000
  70. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  71. family=reuse;policy=mat;refs=4;consumer=100%
  72. family=bag;stmt=cte_filtered;threshold=20
  73. family=bag;stmt=join_where;threshold=200
  74. family=bag;stmt=join_where;threshold=20000
  75. family=reuse;policy=notmat;refs=2;consumer=1%
  76. family=bag;stmt=cte_joined_outer;threshold=2000
  77. family=bag;stmt=cte_joined_outer;threshold=20
  78. family=bag;stmt=filtered_derived;threshold=20000
  79. family=bag;stmt=cte_filtered;threshold=2000
  80. family=reuse;policy=default;refs=16;consumer=100%
  81. family=bag;stmt=cte_filtered_mat;threshold=100000
  82. family=reuse;policy=default;refs=8;consumer=100%
  83. family=bag;stmt=filtered_derived;threshold=200
  84. family=bag;stmt=join_on;threshold=20000
  85. family=bag;stmt=filtered_derived;threshold=180000
  86. family=bag;stmt=cte_joined_outer;threshold=200
  87. family=bag;stmt=cte_filtered_mat;threshold=180000
  88. family=reuse;policy=mat;refs=16;consumer=1%
  89. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  90. family=bag;stmt=filtered_derived;threshold=20
  91. family=bag;stmt=cte_filtered_mat;threshold=20000;position=last

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
| family=bag;stmt=cte_filtered_mat;threshold=20000;position=last | j04m_cte_filtered_mat | valid | `28a581de662adc77` | policy=materialized; observed_cte_scan=true |
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
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `957d40747f6f7eb1` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `549a07108ab01ba9` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `2ce8b85a7b60c88e` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=2;consumer=1% | r_notmat_n2_c1 | valid | `28cc913d71d0d72b` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=2;consumer=100% | r_notmat_n2_c100 | valid | `164db6b1595310eb` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `9f850e53560baafb` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `96aa76348edaa095` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=8;consumer=1% | r_notmat_n8_c1 | valid | `477f7fbaf1c627cb` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=8;consumer=100% | r_notmat_n8_c100 | valid | `ab83914f43e32200` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, VERBOSE, BUFFERS, FORMAT TEXT` supported on postgres
- one host: this run used node budget 2 CPU 3g and client budget 1 CPU 2g, so budget comparisons are separate classes here, not independent hardware
- one client process with one fixed maintenance writer at 10 updates/s; readers are sequential; the reported sample latency includes retry and backoff time
- physical rank-index arm pg-btree-rank; an engine whose planner ignores the index is a different physical design, not a faster one
- distributed RPC counters are N/A on single-node PostgreSQL; shared buffers and rows removed are parsed per cell where printed
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
