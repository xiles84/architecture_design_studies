# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 91 cells, 91 valid, engine `postgres`, topology `single`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-v3-domain-pg-b1-20261005`
- run tag: `run/study-01/select-join-v3-domain-pg-b1-20261005`
- environment: `host-zenbook-ux5406sa`
- repo commit: `0f9d09212c6d52a80d8437d3b5a234f2f711ed36`
- engine version: `PostgreSQL 17.11 (Debian 17.11-1.pgdg13+2) on x86_64-pc-linux-gnu, compiled by gcc (Debian 14.2.0-19) 14.2.0, 64-bit`
- database engine image: `docker.io/library/postgres:17.11` @ `7296f210ae81031ec955dbad9a67a84fe958572a2153b8d0826a647522904dc1`
- benchmark client image: `localhost/charitybench-select:select-join-v3-domain-pg-b1-20261005` @ `8436e36bb66caa8178559d54771b88c063e26848d4ef1d1a56eb524b7f0ecca6`
- inputs digest: `dfada8fb86f01cc41e7ac4d14200dcaafde4acc0f8153696134938c3a2ae005b`
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
| 20 (0.01%) | 0.4115 | 0.5793 | 0.6080 | 0.6703 | 0.4166 | 0.4508 | 0.4504 | 40.5786 | 0.5039 | 0.5229 |
| 200 (0.1%) | 0.9842 | 0.6460 | 0.7077 | 0.7423 | 0.6815 | 0.9246 | 0.9124 | 44.0453 | 0.6896 | 0.6512 |
| 2000 (1%) | 1.5709 | 1.4693 | 1.7115 | 1.7411 | 1.9761 | 1.5928 | 1.9198 | 47.6473 | 1.5498 | 1.8853 |
| 20000 (10%) | 7.0186 | 6.3797 | 6.7061 | 6.6697 | 8.7331 | 6.7058 | 6.4528 | 48.9088 | 6.7400 | 6.6498 |
| 100000 (50%) | 44.7982 | 41.7045 | 44.9644 | 43.2641 | 59.4558 | 47.3177 | 55.8208 | 61.5979 | 42.8307 | 44.9183 |
| 180000 (90%) | 44.9474 | 74.1729 | 64.4818 | 74.7432 | 68.7760 | 77.0669 | 70.6611 | 87.8842 | 45.6406 | 44.1266 |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | join_on vs join_where | 0.4115 | 0.5793 | 1.4078 | true |
| 20 | filtered_derived vs join_where | 0.4115 | 0.6080 | 1.4775 | true |
| 20 | cte_filtered vs join_where | 0.4115 | 0.6703 | 1.6289 | true |
| 20 | cte_filtered_mat vs join_where | 0.4115 | 0.4166 | 1.0123 | false |
| 20 | cte_filtered_notmat vs join_where | 0.4115 | 0.4508 | 1.0955 | true |
| 20 | cte_joined_outer vs join_where | 0.4115 | 0.4504 | 1.0944 | true |
| 20 | cte_joined_outer_mat vs join_where | 0.4115 | 40.5786 | 98.6106 | false |
| 20 | cte_joined_outer_notmat vs join_where | 0.4115 | 0.5039 | 1.2246 | true |
| 20 | comma_join vs join_where | 0.4115 | 0.5229 | 1.2707 | true |
| 200 | join_on vs join_where | 0.9842 | 0.6460 | 0.6564 | true |
| 200 | filtered_derived vs join_where | 0.9842 | 0.7077 | 0.7191 | true |
| 200 | cte_filtered vs join_where | 0.9842 | 0.7423 | 0.7542 | true |
| 200 | cte_filtered_mat vs join_where | 0.9842 | 0.6815 | 0.6925 | false |
| 200 | cte_filtered_notmat vs join_where | 0.9842 | 0.9246 | 0.9394 | true |
| 200 | cte_joined_outer vs join_where | 0.9842 | 0.9124 | 0.9271 | true |
| 200 | cte_joined_outer_mat vs join_where | 0.9842 | 44.0453 | 44.7542 | false |
| 200 | cte_joined_outer_notmat vs join_where | 0.9842 | 0.6896 | 0.7007 | true |
| 200 | comma_join vs join_where | 0.9842 | 0.6512 | 0.6617 | true |
| 2000 | join_on vs join_where | 1.5709 | 1.4693 | 0.9353 | true |
| 2000 | filtered_derived vs join_where | 1.5709 | 1.7115 | 1.0895 | true |
| 2000 | cte_filtered vs join_where | 1.5709 | 1.7411 | 1.1084 | true |
| 2000 | cte_filtered_mat vs join_where | 1.5709 | 1.9761 | 1.2580 | false |
| 2000 | cte_filtered_notmat vs join_where | 1.5709 | 1.5928 | 1.0140 | true |
| 2000 | cte_joined_outer vs join_where | 1.5709 | 1.9198 | 1.2221 | true |
| 2000 | cte_joined_outer_mat vs join_where | 1.5709 | 47.6473 | 30.3320 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 1.5709 | 1.5498 | 0.9866 | true |
| 2000 | comma_join vs join_where | 1.5709 | 1.8853 | 1.2002 | true |
| 20000 | join_on vs join_where | 7.0186 | 6.3797 | 0.9090 | true |
| 20000 | filtered_derived vs join_where | 7.0186 | 6.7061 | 0.9555 | true |
| 20000 | cte_filtered vs join_where | 7.0186 | 6.6697 | 0.9503 | true |
| 20000 | cte_filtered_mat vs join_where | 7.0186 | 8.7331 | 1.2443 | false |
| 20000 | cte_filtered_notmat vs join_where | 7.0186 | 6.7058 | 0.9554 | true |
| 20000 | cte_joined_outer vs join_where | 7.0186 | 6.4528 | 0.9194 | true |
| 20000 | cte_joined_outer_mat vs join_where | 7.0186 | 48.9088 | 6.9685 | false |
| 20000 | cte_joined_outer_notmat vs join_where | 7.0186 | 6.7400 | 0.9603 | true |
| 20000 | comma_join vs join_where | 7.0186 | 6.6498 | 0.9475 | true |
| 100000 | join_on vs join_where | 44.7982 | 41.7045 | 0.9309 | true |
| 100000 | filtered_derived vs join_where | 44.7982 | 44.9644 | 1.0037 | true |
| 100000 | cte_filtered vs join_where | 44.7982 | 43.2641 | 0.9658 | true |
| 100000 | cte_filtered_mat vs join_where | 44.7982 | 59.4558 | 1.3272 | false |
| 100000 | cte_filtered_notmat vs join_where | 44.7982 | 47.3177 | 1.0562 | true |
| 100000 | cte_joined_outer vs join_where | 44.7982 | 55.8208 | 1.2460 | true |
| 100000 | cte_joined_outer_mat vs join_where | 44.7982 | 61.5979 | 1.3750 | false |
| 100000 | cte_joined_outer_notmat vs join_where | 44.7982 | 42.8307 | 0.9561 | true |
| 100000 | comma_join vs join_where | 44.7982 | 44.9183 | 1.0027 | true |
| 180000 | join_on vs join_where | 44.9474 | 74.1729 | 1.6502 | true |
| 180000 | filtered_derived vs join_where | 44.9474 | 64.4818 | 1.4346 | true |
| 180000 | cte_filtered vs join_where | 44.9474 | 74.7432 | 1.6629 | true |
| 180000 | cte_filtered_mat vs join_where | 44.9474 | 68.7760 | 1.5301 | false |
| 180000 | cte_filtered_notmat vs join_where | 44.9474 | 77.0669 | 1.7146 | true |
| 180000 | cte_joined_outer vs join_where | 44.9474 | 70.6611 | 1.5721 | true |
| 180000 | cte_joined_outer_mat vs join_where | 44.9474 | 87.8842 | 1.9553 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 44.9474 | 45.6406 | 1.0154 | true |
| 180000 | comma_join vs join_where | 44.9474 | 44.1266 | 0.9817 | true |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_filtered_mat vs cte_filtered | 0.6703 | 0.4166 | 0.6215 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_filtered_notmat vs cte_filtered | 0.6703 | 0.4508 | 0.6725 | true | policy=not_materialized; observed_cte_scan=false |
| 20 | cte_joined_outer_mat vs cte_joined_outer | 0.4504 | 40.5786 | 90.1007 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 0.4504 | 0.5039 | 1.1189 | true | policy=not_materialized; observed_cte_scan=false |
| 200 | cte_filtered_mat vs cte_filtered | 0.7423 | 0.6815 | 0.9182 | false | policy=materialized; observed_cte_scan=true |
| 200 | cte_filtered_notmat vs cte_filtered | 0.7423 | 0.9246 | 1.2456 | true | policy=not_materialized; observed_cte_scan=false |
| 200 | cte_joined_outer_mat vs cte_joined_outer | 0.9124 | 44.0453 | 48.2731 | false | policy=materialized; observed_cte_scan=true |
| 200 | cte_joined_outer_notmat vs cte_joined_outer | 0.9124 | 0.6896 | 0.7558 | true | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_filtered_mat vs cte_filtered | 1.7411 | 1.9761 | 1.1349 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_filtered_notmat vs cte_filtered | 1.7411 | 1.5928 | 0.9148 | true | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 1.9198 | 47.6473 | 24.8192 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 1.9198 | 1.5498 | 0.8073 | true | policy=not_materialized; observed_cte_scan=false |
| 20000 | cte_filtered_mat vs cte_filtered | 6.6697 | 8.7331 | 1.3094 | false | policy=materialized; observed_cte_scan=true |
| 20000 | cte_filtered_notmat vs cte_filtered | 6.6697 | 6.7058 | 1.0054 | true | policy=not_materialized; observed_cte_scan=false |
| 20000 | cte_joined_outer_mat vs cte_joined_outer | 6.4528 | 48.9088 | 7.5795 | false | policy=materialized; observed_cte_scan=true |
| 20000 | cte_joined_outer_notmat vs cte_joined_outer | 6.4528 | 6.7400 | 1.0445 | true | policy=not_materialized; observed_cte_scan=false |
| 100000 | cte_filtered_mat vs cte_filtered | 43.2641 | 59.4558 | 1.3743 | false | policy=materialized; observed_cte_scan=true |
| 100000 | cte_filtered_notmat vs cte_filtered | 43.2641 | 47.3177 | 1.0937 | true | policy=not_materialized; observed_cte_scan=false |
| 100000 | cte_joined_outer_mat vs cte_joined_outer | 55.8208 | 61.5979 | 1.1035 | false | policy=materialized; observed_cte_scan=true |
| 100000 | cte_joined_outer_notmat vs cte_joined_outer | 55.8208 | 42.8307 | 0.7673 | true | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_filtered_mat vs cte_filtered | 74.7432 | 68.7760 | 0.9202 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_filtered_notmat vs cte_filtered | 74.7432 | 77.0669 | 1.0311 | true | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 70.6611 | 87.8842 | 1.2437 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 70.6611 | 45.6406 | 0.6459 | true | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 0.5037 | 1 | `549a07108ab01ba9` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 2.1265 | 1 | `2ce8b85a7b60c88e` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | 2.5846 | 2 | `b307c41eb54cf452` | policy=default; observed_cte_scan=true | valid |
| default | 2 | 100% | 2.6302 | 2 | `b307c41eb54cf452` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 1% | 2.7861 | 4 | `ccfe1af393deae98` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 2.8141 | 4 | `ccfe1af393deae98` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | 2.8776 | 8 | `c8c73778097931dc` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 100% | 3.0722 | 8 | `c8c73778097931dc` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 1% | 2.9816 | 16 | `b0c356323d30c425` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 3.4118 | 16 | `b0c356323d30c425` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 2.2101 | 1 | `eaa7d6058aaa5e03` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 2.3268 | 1 | `eaa7d6058aaa5e03` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | 2.3395 | 2 | `b307c41eb54cf452` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 100% | 2.7043 | 2 | `b307c41eb54cf452` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 1% | 2.5005 | 4 | `ccfe1af393deae98` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 2.6716 | 4 | `ccfe1af393deae98` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | 2.9517 | 8 | `c8c73778097931dc` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 100% | 3.2026 | 8 | `c8c73778097931dc` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 1% | 2.9838 | 16 | `b0c356323d30c425` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 3.6333 | 16 | `b0c356323d30c425` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 0.4703 | 1 | `549a07108ab01ba9` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 2.4019 | 1 | `2ce8b85a7b60c88e` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | 0.6501 | 2 | `28cc913d71d0d72b` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 100% | 4.5256 | 2 | `1454e972376b79a1` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 1% | 0.8738 | 4 | `9f850e53560baafb` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 6.7140 | 4 | `2bc0a68d5fe6d34e` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | 1.3323 | 8 | `477f7fbaf1c627cb` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 100% | 9.0259 | 8 | `338497e6d0fd64a4` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 1% | 2.3788 | 16 | `1c38b58d8cec08d3` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 14.4996 | 16 | `2cfa3d6bf19fd6bc` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 4.3878 | false |
| notmat | 1 | 1% | 0.9337 | true |
| mat | 1 | 100% | 1.0942 | false |
| notmat | 1 | 100% | 1.1295 | true |
| mat | 2 | 1% | 0.9052 | true |
| notmat | 2 | 1% | 0.2515 | false |
| mat | 2 | 100% | 1.0282 | true |
| notmat | 2 | 100% | 1.7206 | false |
| mat | 4 | 1% | 0.8975 | true |
| notmat | 4 | 1% | 0.3136 | false |
| mat | 4 | 100% | 0.9494 | true |
| notmat | 4 | 100% | 2.3858 | false |
| mat | 8 | 1% | 1.0258 | true |
| notmat | 8 | 1% | 0.4630 | false |
| mat | 8 | 100% | 1.0424 | true |
| notmat | 8 | 100% | 2.9379 | false |
| mat | 16 | 1% | 1.0008 | true |
| notmat | 16 | 1% | 0.7978 | false |
| mat | 16 | 100% | 1.0649 | true |
| notmat | 16 | 100% | 4.2498 | false |

## Repeated-design position control

| first median ms | last median ms | ratio (last/first) | same plan |
|---|---|---|---|
| 0.411503 | 0.418812 | 1.0178 | true |

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

- client: nr_periods 354, nr_throttled 278, throttled_usec 14784912 | quota 100000 100000
, memory 2147483648

- pg-single: nr_periods 353, nr_throttled 0, throttled_usec 0 | quota 200000 100000, memory 3221225472

## Access path per cell

_Engine-reported from each cell's own plan: whether the rank index or a sequential scan was used, rows scanned, PostgreSQL buffers and filter removals, and reader retries of transient conflicts. A cell that scanned the whole table at a selective threshold is a different physical regime, not a faster query._

| axis | status | index scan | seq scan | storage rows scanned | index rows scanned | buffers hit/read | rows removed | retries |
|---|---|---|---|---|---|---|---|---|
| family=bag;stmt=comma_join;threshold=100000 | valid | true | true | 0 | 0 | 100265/0 | 0 | 0 |
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
| family=bag;stmt=join_where;threshold=100000 | valid | true | true | 0 | 0 | 100265/0 | 0 | 0 |
| family=bag;stmt=join_where;threshold=180000 | valid | false | true | 0 | 0 | 1286/0 | 20000 | 0 |
| family=bag;stmt=join_where;threshold=20 | valid | true | true | 0 | 0 | 84/0 | 0 | 0 |
| family=bag;stmt=join_where;threshold=200 | valid | false | true | 0 | 0 | 203/0 | 0 | 0 |
| family=bag;stmt=join_where;threshold=2000 | valid | false | true | 0 | 0 | 1033/0 | 0 | 0 |
| family=bag;stmt=join_where;threshold=20000 | valid | false | true | 0 | 0 | 1349/0 | 0 | 0 |
| family=bag;stmt=join_where;threshold=20;position=last | valid | true | true | 0 | 0 | 84/0 | 0 | 0 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | false | false | 0 | 0 | 0/0 | 0 | 0 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | false | false | 0 | 0 | 0/0 | 0 | 0 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | true | false | 0 | 0 | 24/0 | 1983 | 0 |
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
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | true | false | 0 | 0 | 416/0 | 1983 | 0 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | false | false | 0 | 0 | 0/0 | 0 | 0 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | true | false | 0 | 0 | 26/0 | 1983 | 0 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | false | false | 0 | 0 | 0/0 | 0 | 0 |
| family=reuse;policy=notmat;refs=2;consumer=1% | valid | true | false | 0 | 0 | 60/0 | 1983 | 0 |
| family=reuse;policy=notmat;refs=2;consumer=100% | valid | false | false | 0 | 0 | 0/0 | 0 | 0 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | true | false | 0 | 0 | 120/0 | 1983 | 0 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | false | false | 0 | 0 | 0/0 | 0 | 0 |
| family=reuse;policy=notmat;refs=8;consumer=1% | valid | true | false | 0 | 0 | 240/0 | 1983 | 0 |
| family=reuse;policy=notmat;refs=8;consumer=100% | valid | false | false | 0 | 0 | 0/0 | 0 | 0 |

## Latency samples and derived throughput

_The reader is sequential, so `queries/s` is 1000/median_ms, a derived rate and not a concurrent-capacity test. min/max are the five sample extremes. Each sample includes any retry and backoff time the statement needed after a transient conflict; per-cell retry counts appear in the access-path table._

| axis | status | min ms | median ms | max ms | samples | derived queries/s |
|---|---|---|---|---|---|---|
| family=bag;stmt=comma_join;threshold=100000 | valid | 42.2266 | 44.9183 | 52.3433 | 5 | 22.3 |
| family=bag;stmt=comma_join;threshold=180000 | valid | 43.9076 | 44.1266 | 83.6710 | 5 | 22.7 |
| family=bag;stmt=comma_join;threshold=20 | valid | 0.4961 | 0.5229 | 0.5959 | 5 | 1912.4 |
| family=bag;stmt=comma_join;threshold=200 | valid | 0.6275 | 0.6512 | 0.6901 | 5 | 1535.7 |
| family=bag;stmt=comma_join;threshold=2000 | valid | 1.7965 | 1.8853 | 2.0533 | 5 | 530.4 |
| family=bag;stmt=comma_join;threshold=20000 | valid | 6.5128 | 6.6498 | 6.6942 | 5 | 150.4 |
| family=bag;stmt=cte_filtered;threshold=100000 | valid | 39.0927 | 43.2641 | 44.8008 | 5 | 23.1 |
| family=bag;stmt=cte_filtered;threshold=180000 | valid | 46.3884 | 74.7432 | 77.7143 | 5 | 13.4 |
| family=bag;stmt=cte_filtered;threshold=20 | valid | 0.5644 | 0.6703 | 0.7159 | 5 | 1491.9 |
| family=bag;stmt=cte_filtered;threshold=200 | valid | 0.6470 | 0.7423 | 0.8485 | 5 | 1347.2 |
| family=bag;stmt=cte_filtered;threshold=2000 | valid | 1.6654 | 1.7411 | 1.8194 | 5 | 574.3 |
| family=bag;stmt=cte_filtered;threshold=20000 | valid | 6.3946 | 6.6697 | 36.6273 | 5 | 149.9 |
| family=bag;stmt=cte_filtered_mat;threshold=100000 | valid | 54.4688 | 59.4558 | 70.3235 | 5 | 16.8 |
| family=bag;stmt=cte_filtered_mat;threshold=180000 | valid | 63.5488 | 68.7760 | 80.1887 | 5 | 14.5 |
| family=bag;stmt=cte_filtered_mat;threshold=20 | valid | 0.3193 | 0.4166 | 0.4730 | 5 | 2400.6 |
| family=bag;stmt=cte_filtered_mat;threshold=200 | valid | 0.5617 | 0.6815 | 0.7077 | 5 | 1467.3 |
| family=bag;stmt=cte_filtered_mat;threshold=2000 | valid | 1.8964 | 1.9761 | 38.0879 | 5 | 506.1 |
| family=bag;stmt=cte_filtered_mat;threshold=20000 | valid | 8.0546 | 8.7331 | 34.8245 | 5 | 114.5 |
| family=bag;stmt=cte_filtered_notmat;threshold=100000 | valid | 38.3251 | 47.3177 | 60.6785 | 5 | 21.1 |
| family=bag;stmt=cte_filtered_notmat;threshold=180000 | valid | 43.1659 | 77.0669 | 90.8676 | 5 | 13.0 |
| family=bag;stmt=cte_filtered_notmat;threshold=20 | valid | 0.3814 | 0.4508 | 0.4779 | 5 | 2218.3 |
| family=bag;stmt=cte_filtered_notmat;threshold=200 | valid | 0.8983 | 0.9246 | 0.9854 | 5 | 1081.6 |
| family=bag;stmt=cte_filtered_notmat;threshold=2000 | valid | 1.4713 | 1.5928 | 1.7393 | 5 | 627.8 |
| family=bag;stmt=cte_filtered_notmat;threshold=20000 | valid | 6.4808 | 6.7058 | 7.9385 | 5 | 149.1 |
| family=bag;stmt=cte_joined_outer;threshold=100000 | valid | 49.6244 | 55.8208 | 60.7228 | 5 | 17.9 |
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 43.7254 | 70.6611 | 74.1292 | 5 | 14.2 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 0.4247 | 0.4504 | 0.5091 | 5 | 2220.4 |
| family=bag;stmt=cte_joined_outer;threshold=200 | valid | 0.7759 | 0.9124 | 1.0437 | 5 | 1096.0 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 1.8163 | 1.9198 | 2.5380 | 5 | 520.9 |
| family=bag;stmt=cte_joined_outer;threshold=20000 | valid | 6.4410 | 6.4528 | 7.0271 | 5 | 155.0 |
| family=bag;stmt=cte_joined_outer_mat;threshold=100000 | valid | 59.0224 | 61.5979 | 62.3909 | 5 | 16.2 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 78.3927 | 87.8842 | 95.1124 | 5 | 11.4 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 39.3852 | 40.5786 | 48.5925 | 5 | 24.6 |
| family=bag;stmt=cte_joined_outer_mat;threshold=200 | valid | 43.0763 | 44.0453 | 44.4957 | 5 | 22.7 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 44.6898 | 47.6473 | 50.6179 | 5 | 21.0 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20000 | valid | 45.9869 | 48.9088 | 53.6323 | 5 | 20.4 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=100000 | valid | 38.5825 | 42.8307 | 61.5343 | 5 | 23.3 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 41.9680 | 45.6406 | 50.5756 | 5 | 21.9 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 0.4462 | 0.5039 | 0.5838 | 5 | 1984.5 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=200 | valid | 0.6624 | 0.6896 | 0.8510 | 5 | 1450.2 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 1.4538 | 1.5498 | 1.8964 | 5 | 645.2 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20000 | valid | 6.4152 | 6.7400 | 6.9996 | 5 | 148.4 |
| family=bag;stmt=filtered_derived;threshold=100000 | valid | 41.0955 | 44.9644 | 64.8908 | 5 | 22.2 |
| family=bag;stmt=filtered_derived;threshold=180000 | valid | 59.7847 | 64.4818 | 91.2923 | 5 | 15.5 |
| family=bag;stmt=filtered_derived;threshold=20 | valid | 0.5590 | 0.6080 | 0.8152 | 5 | 1644.8 |
| family=bag;stmt=filtered_derived;threshold=200 | valid | 0.6452 | 0.7077 | 1.2740 | 5 | 1413.0 |
| family=bag;stmt=filtered_derived;threshold=2000 | valid | 1.5535 | 1.7115 | 1.9109 | 5 | 584.3 |
| family=bag;stmt=filtered_derived;threshold=20000 | valid | 6.2359 | 6.7061 | 43.7660 | 5 | 149.1 |
| family=bag;stmt=join_on;threshold=100000 | valid | 39.7688 | 41.7045 | 60.3498 | 5 | 24.0 |
| family=bag;stmt=join_on;threshold=180000 | valid | 72.0909 | 74.1729 | 77.7034 | 5 | 13.5 |
| family=bag;stmt=join_on;threshold=20 | valid | 0.4875 | 0.5793 | 0.6817 | 5 | 1726.1 |
| family=bag;stmt=join_on;threshold=200 | valid | 0.6241 | 0.6460 | 0.6818 | 5 | 1548.1 |
| family=bag;stmt=join_on;threshold=2000 | valid | 1.3960 | 1.4693 | 1.5251 | 5 | 680.6 |
| family=bag;stmt=join_on;threshold=20000 | valid | 6.0913 | 6.3797 | 8.7627 | 5 | 156.7 |
| family=bag;stmt=join_where;threshold=100000 | valid | 41.3739 | 44.7982 | 62.0759 | 5 | 22.3 |
| family=bag;stmt=join_where;threshold=180000 | valid | 43.7395 | 44.9474 | 86.1218 | 5 | 22.2 |
| family=bag;stmt=join_where;threshold=20 | valid | 0.3835 | 0.4115 | 0.5025 | 5 | 2430.1 |
| family=bag;stmt=join_where;threshold=200 | valid | 0.8622 | 0.9842 | 42.6789 | 5 | 1016.1 |
| family=bag;stmt=join_where;threshold=2000 | valid | 1.4455 | 1.5709 | 1.6655 | 5 | 636.6 |
| family=bag;stmt=join_where;threshold=20000 | valid | 6.5348 | 7.0186 | 33.4481 | 5 | 142.5 |
| family=bag;stmt=join_where;threshold=20;position=last | valid | 0.3530 | 0.4188 | 0.5150 | 5 | 2387.7 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 2.8429 | 2.9816 | 3.1478 | 5 | 335.4 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 3.2459 | 3.4118 | 3.7417 | 5 | 293.1 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 0.4086 | 0.5037 | 0.5512 | 5 | 1985.3 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 2.0275 | 2.1265 | 2.8348 | 5 | 470.3 |
| family=reuse;policy=default;refs=2;consumer=1% | valid | 2.3738 | 2.5846 | 3.1013 | 5 | 386.9 |
| family=reuse;policy=default;refs=2;consumer=100% | valid | 2.5373 | 2.6302 | 3.3724 | 5 | 380.2 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 2.6868 | 2.7861 | 2.8848 | 5 | 358.9 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 2.6316 | 2.8141 | 4.2531 | 5 | 355.4 |
| family=reuse;policy=default;refs=8;consumer=1% | valid | 2.7242 | 2.8776 | 3.3502 | 5 | 347.5 |
| family=reuse;policy=default;refs=8;consumer=100% | valid | 2.9403 | 3.0722 | 3.1745 | 5 | 325.5 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 2.9112 | 2.9838 | 3.3880 | 5 | 335.1 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 3.4644 | 3.6333 | 3.7736 | 5 | 275.2 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 2.1872 | 2.2101 | 2.3732 | 5 | 452.5 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 2.2286 | 2.3268 | 2.4434 | 5 | 429.8 |
| family=reuse;policy=mat;refs=2;consumer=1% | valid | 2.1523 | 2.3395 | 2.4979 | 5 | 427.4 |
| family=reuse;policy=mat;refs=2;consumer=100% | valid | 2.5556 | 2.7043 | 2.7900 | 5 | 369.8 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 2.3636 | 2.5005 | 2.7295 | 5 | 399.9 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 2.5826 | 2.6716 | 2.7883 | 5 | 374.3 |
| family=reuse;policy=mat;refs=8;consumer=1% | valid | 2.7542 | 2.9517 | 3.3985 | 5 | 338.8 |
| family=reuse;policy=mat;refs=8;consumer=100% | valid | 3.0331 | 3.2026 | 3.3535 | 5 | 312.2 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 2.2396 | 2.3788 | 2.4799 | 5 | 420.4 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 13.7773 | 14.4996 | 16.9133 | 5 | 69.0 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 0.3946 | 0.4703 | 0.5559 | 5 | 2126.2 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 2.3067 | 2.4019 | 2.5814 | 5 | 416.3 |
| family=reuse;policy=notmat;refs=2;consumer=1% | valid | 0.5823 | 0.6501 | 1.2644 | 5 | 1538.2 |
| family=reuse;policy=notmat;refs=2;consumer=100% | valid | 4.3756 | 4.5256 | 4.6972 | 5 | 221.0 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 0.7979 | 0.8738 | 0.9811 | 5 | 1144.5 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 6.6405 | 6.7140 | 7.3040 | 5 | 148.9 |
| family=reuse;policy=notmat;refs=8;consumer=1% | valid | 1.2969 | 1.3323 | 24.6276 | 5 | 750.6 |
| family=reuse;policy=notmat;refs=8;consumer=100% | valid | 8.6540 | 9.0259 | 10.2159 | 5 | 110.8 |

## Matched blended workload

- phase: `blended`; target offered rate: 10 updates/s
- offered 324, committed 324, failed 0 over 32.326 s
- offered rate 10.023/s, completed rate 10.023/s
- writer operations per endpoint: [324]
- reader operations overlapping a committed write: 101 of 637
- reader retries of transient conflicts: 0
- per key: 324 keys, min 1, max 1 per key
- ledger: 324 keys reconciled, 0 mismatch
- update latency: avg 1.814 ms, min 0.817, max 21.331 over 324

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42907624`, 5 interleaved rounds
  1. family=bag;stmt=join_where;threshold=20
  2. family=bag;stmt=join_on;threshold=100000
  3. family=reuse;policy=mat;refs=2;consumer=100%
  4. family=bag;stmt=join_where;threshold=20000
  5. family=bag;stmt=cte_filtered_mat;threshold=2000
  6. family=reuse;policy=default;refs=16;consumer=100%
  7. family=bag;stmt=cte_filtered;threshold=180000
  8. family=bag;stmt=cte_filtered;threshold=20
  9. family=reuse;policy=default;refs=1;consumer=1%
  10. family=reuse;policy=default;refs=8;consumer=1%
  11. family=bag;stmt=comma_join;threshold=200
  12. family=bag;stmt=comma_join;threshold=20000
  13. family=bag;stmt=filtered_derived;threshold=100000
  14. family=bag;stmt=cte_joined_outer;threshold=2000
  15. family=reuse;policy=mat;refs=1;consumer=100%
  16. family=bag;stmt=join_on;threshold=200
  17. family=reuse;policy=mat;refs=2;consumer=1%
  18. family=bag;stmt=join_on;threshold=2000
  19. family=bag;stmt=cte_joined_outer;threshold=180000
  20. family=reuse;policy=mat;refs=16;consumer=100%
  21. family=bag;stmt=cte_filtered;threshold=200
  22. family=bag;stmt=cte_filtered_notmat;threshold=20
  23. family=bag;stmt=cte_filtered_notmat;threshold=20000
  24. family=bag;stmt=join_on;threshold=180000
  25. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  26. family=reuse;policy=mat;refs=8;consumer=1%
  27. family=reuse;policy=notmat;refs=1;consumer=1%
  28. family=reuse;policy=mat;refs=1;consumer=1%
  29. family=bag;stmt=cte_joined_outer;threshold=20000
  30. family=bag;stmt=cte_filtered;threshold=2000
  31. family=bag;stmt=filtered_derived;threshold=200
  32. family=bag;stmt=cte_joined_outer;threshold=100000
  33. family=reuse;policy=default;refs=4;consumer=1%
  34. family=bag;stmt=cte_filtered_notmat;threshold=2000
  35. family=reuse;policy=notmat;refs=16;consumer=1%
  36. family=bag;stmt=cte_joined_outer_notmat;threshold=200
  37. family=bag;stmt=cte_joined_outer_mat;threshold=20000
  38. family=bag;stmt=join_where;threshold=180000
  39. family=bag;stmt=filtered_derived;threshold=20
  40. family=bag;stmt=filtered_derived;threshold=20000
  41. family=bag;stmt=cte_filtered_notmat;threshold=180000
  42. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  43. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  44. family=bag;stmt=join_where;threshold=200
  45. family=reuse;policy=notmat;refs=16;consumer=100%
  46. family=bag;stmt=filtered_derived;threshold=2000
  47. family=reuse;policy=default;refs=16;consumer=1%
  48. family=bag;stmt=cte_filtered_mat;threshold=180000
  49. family=reuse;policy=notmat;refs=8;consumer=100%
  50. family=bag;stmt=comma_join;threshold=100000
  51. family=reuse;policy=notmat;refs=1;consumer=100%
  52. family=bag;stmt=cte_joined_outer;threshold=20
  53. family=bag;stmt=cte_joined_outer_mat;threshold=20
  54. family=reuse;policy=default;refs=2;consumer=1%
  55. family=bag;stmt=cte_joined_outer_notmat;threshold=20000
  56. family=reuse;policy=notmat;refs=4;consumer=100%
  57. family=bag;stmt=comma_join;threshold=20
  58. family=reuse;policy=notmat;refs=4;consumer=1%
  59. family=bag;stmt=cte_filtered_mat;threshold=20
  60. family=bag;stmt=comma_join;threshold=180000
  61. family=bag;stmt=comma_join;threshold=2000
  62. family=bag;stmt=join_on;threshold=20000
  63. family=bag;stmt=join_on;threshold=20
  64. family=bag;stmt=cte_joined_outer_mat;threshold=100000
  65. family=reuse;policy=default;refs=2;consumer=100%
  66. family=reuse;policy=notmat;refs=8;consumer=1%
  67. family=reuse;policy=mat;refs=4;consumer=1%
  68. family=bag;stmt=join_where;threshold=2000
  69. family=bag;stmt=cte_filtered_mat;threshold=200
  70. family=reuse;policy=mat;refs=16;consumer=1%
  71. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  72. family=bag;stmt=filtered_derived;threshold=180000
  73. family=reuse;policy=default;refs=8;consumer=100%
  74. family=reuse;policy=default;refs=1;consumer=100%
  75. family=bag;stmt=cte_filtered_mat;threshold=20000
  76. family=bag;stmt=cte_filtered;threshold=20000
  77. family=reuse;policy=default;refs=4;consumer=100%
  78. family=reuse;policy=notmat;refs=2;consumer=1%
  79. family=bag;stmt=join_where;threshold=100000
  80. family=reuse;policy=notmat;refs=2;consumer=100%
  81. family=bag;stmt=cte_joined_outer_mat;threshold=200
  82. family=reuse;policy=mat;refs=4;consumer=100%
  83. family=bag;stmt=cte_filtered;threshold=100000
  84. family=reuse;policy=mat;refs=8;consumer=100%
  85. family=bag;stmt=cte_filtered_notmat;threshold=100000
  86. family=bag;stmt=cte_filtered_mat;threshold=100000
  87. family=bag;stmt=cte_joined_outer;threshold=200
  88. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  89. family=bag;stmt=cte_joined_outer_notmat;threshold=100000
  90. family=bag;stmt=cte_filtered_notmat;threshold=200
  91. family=bag;stmt=join_where;threshold=20;position=last

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
| family=bag;stmt=join_where;threshold=20;position=last | j01_join_where | valid | `61819ba7d0f5f952` |  |
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
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `2cfa3d6bf19fd6bc` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `549a07108ab01ba9` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `2ce8b85a7b60c88e` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=2;consumer=1% | r_notmat_n2_c1 | valid | `28cc913d71d0d72b` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=2;consumer=100% | r_notmat_n2_c100 | valid | `1454e972376b79a1` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `9f850e53560baafb` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `2bc0a68d5fe6d34e` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=8;consumer=1% | r_notmat_n8_c1 | valid | `477f7fbaf1c627cb` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=8;consumer=100% | r_notmat_n8_c100 | valid | `338497e6d0fd64a4` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, VERBOSE, BUFFERS, FORMAT TEXT` supported on postgres
- one host: this run used node budget 2 CPU 3g and client budget 1 CPU 2g, so budget comparisons are separate classes here, not independent hardware
- one client process with one fixed maintenance writer at 10 updates/s; readers are sequential; the reported sample latency includes retry and backoff time
- physical rank-index arm pg-btree-rank; an engine whose planner ignores the index is a different physical design, not a faster one
- distributed RPC counters are N/A on single-node PostgreSQL; shared buffers and rows removed are parsed per cell where printed
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
