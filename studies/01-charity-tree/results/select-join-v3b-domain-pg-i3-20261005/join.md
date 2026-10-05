# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 91 cells, 91 valid, engine `postgres`, topology `single`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-v3b-domain-pg-i3-20261005`
- run tag: `run/study-01/select-join-v3b-domain-pg-i3-20261005`
- environment: `host-zenbook-ux5406sa`
- repo commit: `7243fee074bfb16fc5ad8b6afc1a2d7549080935`
- engine version: `PostgreSQL 17.11 (Debian 17.11-1.pgdg13+2) on x86_64-pc-linux-gnu, compiled by gcc (Debian 14.2.0-19) 14.2.0, 64-bit`
- database engine image: `docker.io/library/postgres:17.11` @ `7296f210ae81031ec955dbad9a67a84fe958572a2153b8d0826a647522904dc1`
- benchmark client image: `localhost/charitybench-select:select-join-v3b-domain-pg-i3-20261005` @ `f3f1e10c44cbcb4c7b5bc56ed95a4f071682b6f7d33b5e0a7b688c4271eeef88`
- inputs digest: `6a2111f1912b1cd00d196e68d818df6cf08c060976d579f884c1b6e033110d3d`
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
| 20 (0.01%) | 0.7061 | 0.4750 | 0.5387 | 0.6643 | 0.7395 | 0.4287 | 0.4971 | 48.4540 | 0.7524 | 0.7213 |
| 200 (0.1%) | 0.8913 | 0.9292 | 0.8377 | 0.7709 | 0.7268 | 1.0479 | 0.9989 | 51.7686 | 1.0361 | 0.7560 |
| 2000 (1%) | 1.8410 | 1.7797 | 1.8511 | 2.0126 | 2.0906 | 2.0417 | 1.7816 | 47.5223 | 1.8431 | 1.9456 |
| 20000 (10%) | 7.6768 | 7.3328 | 7.1140 | 7.6167 | 9.5962 | 7.0524 | 7.1662 | 55.8272 | 7.1995 | 8.3874 |
| 100000 (50%) | 52.2632 | 48.0515 | 54.3888 | 57.8186 | 63.8025 | 50.1601 | 46.5513 | 66.9646 | 46.8125 | 50.1389 |
| 180000 (90%) | 70.0986 | 66.8251 | 88.3348 | 50.2553 | 80.8848 | 49.2250 | 73.4084 | 76.6270 | 49.7390 | 75.3407 |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | join_on vs join_where | 0.7061 | 0.4750 | 0.6728 | true |
| 20 | filtered_derived vs join_where | 0.7061 | 0.5387 | 0.7630 | true |
| 20 | cte_filtered vs join_where | 0.7061 | 0.6643 | 0.9409 | true |
| 20 | cte_filtered_mat vs join_where | 0.7061 | 0.7395 | 1.0474 | false |
| 20 | cte_filtered_notmat vs join_where | 0.7061 | 0.4287 | 0.6071 | true |
| 20 | cte_joined_outer vs join_where | 0.7061 | 0.4971 | 0.7040 | true |
| 20 | cte_joined_outer_mat vs join_where | 0.7061 | 48.4540 | 68.6250 | false |
| 20 | cte_joined_outer_notmat vs join_where | 0.7061 | 0.7524 | 1.0657 | true |
| 20 | comma_join vs join_where | 0.7061 | 0.7213 | 1.0215 | true |
| 200 | join_on vs join_where | 0.8913 | 0.9292 | 1.0425 | true |
| 200 | filtered_derived vs join_where | 0.8913 | 0.8377 | 0.9399 | true |
| 200 | cte_filtered vs join_where | 0.8913 | 0.7709 | 0.8649 | true |
| 200 | cte_filtered_mat vs join_where | 0.8913 | 0.7268 | 0.8154 | false |
| 200 | cte_filtered_notmat vs join_where | 0.8913 | 1.0479 | 1.1757 | true |
| 200 | cte_joined_outer vs join_where | 0.8913 | 0.9989 | 1.1207 | true |
| 200 | cte_joined_outer_mat vs join_where | 0.8913 | 51.7686 | 58.0797 | false |
| 200 | cte_joined_outer_notmat vs join_where | 0.8913 | 1.0361 | 1.1624 | true |
| 200 | comma_join vs join_where | 0.8913 | 0.7560 | 0.8482 | true |
| 2000 | join_on vs join_where | 1.8410 | 1.7797 | 0.9667 | true |
| 2000 | filtered_derived vs join_where | 1.8410 | 1.8511 | 1.0055 | true |
| 2000 | cte_filtered vs join_where | 1.8410 | 2.0126 | 1.0932 | true |
| 2000 | cte_filtered_mat vs join_where | 1.8410 | 2.0906 | 1.1356 | false |
| 2000 | cte_filtered_notmat vs join_where | 1.8410 | 2.0417 | 1.1090 | true |
| 2000 | cte_joined_outer vs join_where | 1.8410 | 1.7816 | 0.9678 | true |
| 2000 | cte_joined_outer_mat vs join_where | 1.8410 | 47.5223 | 25.8137 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 1.8410 | 1.8431 | 1.0011 | true |
| 2000 | comma_join vs join_where | 1.8410 | 1.9456 | 1.0568 | true |
| 20000 | join_on vs join_where | 7.6768 | 7.3328 | 0.9552 | true |
| 20000 | filtered_derived vs join_where | 7.6768 | 7.1140 | 0.9267 | true |
| 20000 | cte_filtered vs join_where | 7.6768 | 7.6167 | 0.9922 | true |
| 20000 | cte_filtered_mat vs join_where | 7.6768 | 9.5962 | 1.2500 | false |
| 20000 | cte_filtered_notmat vs join_where | 7.6768 | 7.0524 | 0.9187 | true |
| 20000 | cte_joined_outer vs join_where | 7.6768 | 7.1662 | 0.9335 | true |
| 20000 | cte_joined_outer_mat vs join_where | 7.6768 | 55.8272 | 7.2722 | false |
| 20000 | cte_joined_outer_notmat vs join_where | 7.6768 | 7.1995 | 0.9378 | true |
| 20000 | comma_join vs join_where | 7.6768 | 8.3874 | 1.0926 | true |
| 100000 | join_on vs join_where | 52.2632 | 48.0515 | 0.9194 | true |
| 100000 | filtered_derived vs join_where | 52.2632 | 54.3888 | 1.0407 | true |
| 100000 | cte_filtered vs join_where | 52.2632 | 57.8186 | 1.1063 | true |
| 100000 | cte_filtered_mat vs join_where | 52.2632 | 63.8025 | 1.2208 | false |
| 100000 | cte_filtered_notmat vs join_where | 52.2632 | 50.1601 | 0.9598 | true |
| 100000 | cte_joined_outer vs join_where | 52.2632 | 46.5513 | 0.8907 | true |
| 100000 | cte_joined_outer_mat vs join_where | 52.2632 | 66.9646 | 1.2813 | false |
| 100000 | cte_joined_outer_notmat vs join_where | 52.2632 | 46.8125 | 0.8957 | true |
| 100000 | comma_join vs join_where | 52.2632 | 50.1389 | 0.9594 | true |
| 180000 | join_on vs join_where | 70.0986 | 66.8251 | 0.9533 | true |
| 180000 | filtered_derived vs join_where | 70.0986 | 88.3348 | 1.2602 | true |
| 180000 | cte_filtered vs join_where | 70.0986 | 50.2553 | 0.7169 | true |
| 180000 | cte_filtered_mat vs join_where | 70.0986 | 80.8848 | 1.1539 | false |
| 180000 | cte_filtered_notmat vs join_where | 70.0986 | 49.2250 | 0.7022 | true |
| 180000 | cte_joined_outer vs join_where | 70.0986 | 73.4084 | 1.0472 | true |
| 180000 | cte_joined_outer_mat vs join_where | 70.0986 | 76.6270 | 1.0931 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 70.0986 | 49.7390 | 0.7096 | true |
| 180000 | comma_join vs join_where | 70.0986 | 75.3407 | 1.0748 | true |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_filtered_mat vs cte_filtered | 0.6643 | 0.7395 | 1.1132 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_filtered_notmat vs cte_filtered | 0.6643 | 0.4287 | 0.6452 | true | policy=not_materialized; observed_cte_scan=false |
| 20 | cte_joined_outer_mat vs cte_joined_outer | 0.4971 | 48.4540 | 97.4732 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 0.4971 | 0.7524 | 1.5136 | true | policy=not_materialized; observed_cte_scan=false |
| 200 | cte_filtered_mat vs cte_filtered | 0.7709 | 0.7268 | 0.9428 | false | policy=materialized; observed_cte_scan=true |
| 200 | cte_filtered_notmat vs cte_filtered | 0.7709 | 1.0479 | 1.3594 | true | policy=not_materialized; observed_cte_scan=false |
| 200 | cte_joined_outer_mat vs cte_joined_outer | 0.9989 | 51.7686 | 51.8235 | false | policy=materialized; observed_cte_scan=true |
| 200 | cte_joined_outer_notmat vs cte_joined_outer | 0.9989 | 1.0361 | 1.0372 | true | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_filtered_mat vs cte_filtered | 2.0126 | 2.0906 | 1.0388 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_filtered_notmat vs cte_filtered | 2.0126 | 2.0417 | 1.0145 | true | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 1.7816 | 47.5223 | 26.6738 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 1.7816 | 1.8431 | 1.0345 | true | policy=not_materialized; observed_cte_scan=false |
| 20000 | cte_filtered_mat vs cte_filtered | 7.6167 | 9.5962 | 1.2599 | false | policy=materialized; observed_cte_scan=true |
| 20000 | cte_filtered_notmat vs cte_filtered | 7.6167 | 7.0524 | 0.9259 | true | policy=not_materialized; observed_cte_scan=false |
| 20000 | cte_joined_outer_mat vs cte_joined_outer | 7.1662 | 55.8272 | 7.7903 | false | policy=materialized; observed_cte_scan=true |
| 20000 | cte_joined_outer_notmat vs cte_joined_outer | 7.1662 | 7.1995 | 1.0046 | true | policy=not_materialized; observed_cte_scan=false |
| 100000 | cte_filtered_mat vs cte_filtered | 57.8186 | 63.8025 | 1.1035 | false | policy=materialized; observed_cte_scan=true |
| 100000 | cte_filtered_notmat vs cte_filtered | 57.8186 | 50.1601 | 0.8675 | true | policy=not_materialized; observed_cte_scan=false |
| 100000 | cte_joined_outer_mat vs cte_joined_outer | 46.5513 | 66.9646 | 1.4385 | false | policy=materialized; observed_cte_scan=true |
| 100000 | cte_joined_outer_notmat vs cte_joined_outer | 46.5513 | 46.8125 | 1.0056 | true | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_filtered_mat vs cte_filtered | 50.2553 | 80.8848 | 1.6095 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_filtered_notmat vs cte_filtered | 50.2553 | 49.2250 | 0.9795 | true | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 73.4084 | 76.6270 | 1.0438 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 73.4084 | 49.7390 | 0.6776 | true | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 0.7328 | 1 | `549a07108ab01ba9` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 2.7212 | 1 | `2ce8b85a7b60c88e` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | 2.4471 | 2 | `b307c41eb54cf452` | policy=default; observed_cte_scan=true | valid |
| default | 2 | 100% | 3.0303 | 2 | `b307c41eb54cf452` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 1% | 2.7572 | 4 | `ccfe1af393deae98` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 3.0613 | 4 | `ccfe1af393deae98` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | 2.8666 | 8 | `c8c73778097931dc` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 100% | 3.2170 | 8 | `c8c73778097931dc` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 1% | 3.8915 | 16 | `b0c356323d30c425` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 4.0729 | 16 | `b0c356323d30c425` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 2.5651 | 1 | `eaa7d6058aaa5e03` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 2.7046 | 1 | `eaa7d6058aaa5e03` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | 2.7639 | 2 | `b307c41eb54cf452` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 100% | 2.9788 | 2 | `b307c41eb54cf452` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 1% | 2.7532 | 4 | `ccfe1af393deae98` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 3.0078 | 4 | `ccfe1af393deae98` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | 2.9525 | 8 | `c8c73778097931dc` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 100% | 3.2280 | 8 | `c8c73778097931dc` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 1% | 3.8925 | 16 | `b0c356323d30c425` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 4.0885 | 16 | `b0c356323d30c425` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 0.4815 | 1 | `549a07108ab01ba9` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 2.6732 | 1 | `2ce8b85a7b60c88e` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | 0.8109 | 2 | `28cc913d71d0d72b` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 100% | 5.0722 | 2 | `53ffbc7320ed83a8` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 1% | 0.9861 | 4 | `9f850e53560baafb` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 7.5602 | 4 | `1888423fc6a4024a` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | 1.4900 | 8 | `477f7fbaf1c627cb` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 100% | 10.5705 | 8 | `e2f06cbda0e4dce8` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 1% | 2.6584 | 16 | `1c38b58d8cec08d3` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 16.1016 | 16 | `f632f422e7d95269` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 3.5003 | false |
| notmat | 1 | 1% | 0.6570 | true |
| mat | 1 | 100% | 0.9939 | false |
| notmat | 1 | 100% | 0.9824 | true |
| mat | 2 | 1% | 1.1294 | true |
| notmat | 2 | 1% | 0.3313 | false |
| mat | 2 | 100% | 0.9830 | true |
| notmat | 2 | 100% | 1.6738 | false |
| mat | 4 | 1% | 0.9986 | true |
| notmat | 4 | 1% | 0.3577 | false |
| mat | 4 | 100% | 0.9825 | true |
| notmat | 4 | 100% | 2.4696 | false |
| mat | 8 | 1% | 1.0300 | true |
| notmat | 8 | 1% | 0.5198 | false |
| mat | 8 | 100% | 1.0034 | true |
| notmat | 8 | 100% | 3.2858 | false |
| mat | 16 | 1% | 1.0003 | true |
| notmat | 16 | 1% | 0.6831 | false |
| mat | 16 | 100% | 1.0038 | true |
| notmat | 16 | 100% | 3.9533 | false |

## Repeated-design position control

| first median ms | last median ms | ratio (last/first) | same plan |
|---|---|---|---|
| 0.428654 | 0.608068 | 1.4186 | true |

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

- client: nr_periods 404, nr_throttled 325, throttled_usec 15767589 | quota 100000 100000
, memory 2147483648

- pg-single: nr_periods 379, nr_throttled 0, throttled_usec 0 | quota 200000 100000, memory 3221225472

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
| family=bag;stmt=cte_filtered_notmat;threshold=20;position=last | valid | true | true | 0 | 0 | 84/0 | 0 | 0 |
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
| family=bag;stmt=comma_join;threshold=100000 | valid | 45.6424 | 50.1389 | 55.6037 | 5 | 19.9 |
| family=bag;stmt=comma_join;threshold=180000 | valid | 44.5477 | 75.3407 | 83.8305 | 5 | 13.3 |
| family=bag;stmt=comma_join;threshold=20 | valid | 0.6209 | 0.7213 | 0.7957 | 5 | 1386.5 |
| family=bag;stmt=comma_join;threshold=200 | valid | 0.7010 | 0.7560 | 1.5121 | 5 | 1322.7 |
| family=bag;stmt=comma_join;threshold=2000 | valid | 1.7023 | 1.9456 | 2.3160 | 5 | 514.0 |
| family=bag;stmt=comma_join;threshold=20000 | valid | 7.2689 | 8.3874 | 14.6389 | 5 | 119.2 |
| family=bag;stmt=cte_filtered;threshold=100000 | valid | 44.5707 | 57.8186 | 66.7759 | 5 | 17.3 |
| family=bag;stmt=cte_filtered;threshold=180000 | valid | 43.1123 | 50.2553 | 92.4202 | 5 | 19.9 |
| family=bag;stmt=cte_filtered;threshold=20 | valid | 0.6268 | 0.6643 | 0.9589 | 5 | 1505.2 |
| family=bag;stmt=cte_filtered;threshold=200 | valid | 0.7018 | 0.7709 | 0.9531 | 5 | 1297.2 |
| family=bag;stmt=cte_filtered;threshold=2000 | valid | 1.9036 | 2.0126 | 2.3904 | 5 | 496.9 |
| family=bag;stmt=cte_filtered;threshold=20000 | valid | 7.5448 | 7.6167 | 24.1733 | 5 | 131.3 |
| family=bag;stmt=cte_filtered_mat;threshold=100000 | valid | 58.7417 | 63.8025 | 65.4325 | 5 | 15.7 |
| family=bag;stmt=cte_filtered_mat;threshold=180000 | valid | 71.0025 | 80.8848 | 223.0808 | 5 | 12.4 |
| family=bag;stmt=cte_filtered_mat;threshold=20 | valid | 0.6457 | 0.7395 | 1.0630 | 5 | 1352.2 |
| family=bag;stmt=cte_filtered_mat;threshold=200 | valid | 0.6907 | 0.7268 | 1.4097 | 5 | 1375.9 |
| family=bag;stmt=cte_filtered_mat;threshold=2000 | valid | 1.8933 | 2.0906 | 17.0888 | 5 | 478.3 |
| family=bag;stmt=cte_filtered_mat;threshold=20000 | valid | 9.3340 | 9.5962 | 34.7997 | 5 | 104.2 |
| family=bag;stmt=cte_filtered_notmat;threshold=100000 | valid | 44.3772 | 50.1601 | 75.1361 | 5 | 19.9 |
| family=bag;stmt=cte_filtered_notmat;threshold=180000 | valid | 47.2959 | 49.2250 | 123.3307 | 5 | 20.3 |
| family=bag;stmt=cte_filtered_notmat;threshold=20 | valid | 0.3874 | 0.4287 | 2.4429 | 5 | 2332.9 |
| family=bag;stmt=cte_filtered_notmat;threshold=200 | valid | 0.9504 | 1.0479 | 1.4460 | 5 | 954.3 |
| family=bag;stmt=cte_filtered_notmat;threshold=2000 | valid | 1.8705 | 2.0417 | 3.4174 | 5 | 489.8 |
| family=bag;stmt=cte_filtered_notmat;threshold=20000 | valid | 6.7361 | 7.0524 | 10.5837 | 5 | 141.8 |
| family=bag;stmt=cte_filtered_notmat;threshold=20;position=last | valid | 0.3386 | 0.6081 | 0.9172 | 5 | 1644.6 |
| family=bag;stmt=cte_joined_outer;threshold=100000 | valid | 45.0208 | 46.5513 | 50.0076 | 5 | 21.5 |
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 68.3701 | 73.4084 | 92.7591 | 5 | 13.6 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 0.3930 | 0.4971 | 0.6520 | 5 | 2011.7 |
| family=bag;stmt=cte_joined_outer;threshold=200 | valid | 0.9588 | 0.9989 | 1.0657 | 5 | 1001.1 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 1.6851 | 1.7816 | 2.5860 | 5 | 561.3 |
| family=bag;stmt=cte_joined_outer;threshold=20000 | valid | 6.8992 | 7.1662 | 28.2570 | 5 | 139.5 |
| family=bag;stmt=cte_joined_outer_mat;threshold=100000 | valid | 63.2554 | 66.9646 | 105.1200 | 5 | 14.9 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 72.4784 | 76.6270 | 133.4934 | 5 | 13.1 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 47.6531 | 48.4540 | 51.1833 | 5 | 20.6 |
| family=bag;stmt=cte_joined_outer_mat;threshold=200 | valid | 48.4304 | 51.7686 | 86.8588 | 5 | 19.3 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 46.5000 | 47.5223 | 60.5952 | 5 | 21.0 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20000 | valid | 53.3324 | 55.8272 | 94.3057 | 5 | 17.9 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=100000 | valid | 44.0491 | 46.8125 | 62.5849 | 5 | 21.4 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 47.6984 | 49.7390 | 85.3110 | 5 | 20.1 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 0.6655 | 0.7524 | 0.8880 | 5 | 1329.0 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=200 | valid | 0.9156 | 1.0361 | 1.4870 | 5 | 965.2 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 1.7169 | 1.8431 | 3.1447 | 5 | 542.6 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20000 | valid | 6.8088 | 7.1995 | 18.4108 | 5 | 138.9 |
| family=bag;stmt=filtered_derived;threshold=100000 | valid | 46.2075 | 54.3888 | 132.2322 | 5 | 18.4 |
| family=bag;stmt=filtered_derived;threshold=180000 | valid | 72.6521 | 88.3348 | 96.6954 | 5 | 11.3 |
| family=bag;stmt=filtered_derived;threshold=20 | valid | 0.4792 | 0.5387 | 4.9754 | 5 | 1856.3 |
| family=bag;stmt=filtered_derived;threshold=200 | valid | 0.7169 | 0.8377 | 33.2841 | 5 | 1193.7 |
| family=bag;stmt=filtered_derived;threshold=2000 | valid | 1.6650 | 1.8511 | 3.8439 | 5 | 540.2 |
| family=bag;stmt=filtered_derived;threshold=20000 | valid | 6.7297 | 7.1140 | 21.2081 | 5 | 140.6 |
| family=bag;stmt=join_on;threshold=100000 | valid | 47.0274 | 48.0515 | 154.5636 | 5 | 20.8 |
| family=bag;stmt=join_on;threshold=180000 | valid | 54.3435 | 66.8251 | 78.2977 | 5 | 15.0 |
| family=bag;stmt=join_on;threshold=20 | valid | 0.4113 | 0.4750 | 0.7606 | 5 | 2105.0 |
| family=bag;stmt=join_on;threshold=200 | valid | 0.9153 | 0.9292 | 1.1634 | 5 | 1076.2 |
| family=bag;stmt=join_on;threshold=2000 | valid | 1.7007 | 1.7797 | 1.9431 | 5 | 561.9 |
| family=bag;stmt=join_on;threshold=20000 | valid | 6.9034 | 7.3328 | 7.7726 | 5 | 136.4 |
| family=bag;stmt=join_where;threshold=100000 | valid | 43.7777 | 52.2632 | 73.9065 | 5 | 19.1 |
| family=bag;stmt=join_where;threshold=180000 | valid | 44.9155 | 70.0986 | 78.2170 | 5 | 14.3 |
| family=bag;stmt=join_where;threshold=20 | valid | 0.6035 | 0.7061 | 4.4278 | 5 | 1416.3 |
| family=bag;stmt=join_where;threshold=200 | valid | 0.8258 | 0.8913 | 1.1019 | 5 | 1121.9 |
| family=bag;stmt=join_where;threshold=2000 | valid | 1.7346 | 1.8410 | 3.6060 | 5 | 543.2 |
| family=bag;stmt=join_where;threshold=20000 | valid | 6.8516 | 7.6768 | 22.2774 | 5 | 130.3 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 3.6200 | 3.8915 | 4.6291 | 5 | 257.0 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 3.9329 | 4.0729 | 4.1470 | 5 | 245.5 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 0.5790 | 0.7328 | 0.7621 | 5 | 1364.6 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 2.6192 | 2.7212 | 3.1250 | 5 | 367.5 |
| family=reuse;policy=default;refs=2;consumer=1% | valid | 2.2600 | 2.4471 | 3.1859 | 5 | 408.6 |
| family=reuse;policy=default;refs=2;consumer=100% | valid | 2.8551 | 3.0303 | 3.7869 | 5 | 330.0 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 2.5468 | 2.7572 | 2.9358 | 5 | 362.7 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 2.7089 | 3.0613 | 4.4499 | 5 | 326.7 |
| family=reuse;policy=default;refs=8;consumer=1% | valid | 2.8330 | 2.8666 | 2.8764 | 5 | 348.8 |
| family=reuse;policy=default;refs=8;consumer=100% | valid | 3.0309 | 3.2170 | 3.2675 | 5 | 310.8 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 3.5506 | 3.8925 | 4.8796 | 5 | 256.9 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 3.7839 | 4.0885 | 9.2431 | 5 | 244.6 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 2.5176 | 2.5651 | 3.7778 | 5 | 389.8 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 2.6020 | 2.7046 | 3.5495 | 5 | 369.7 |
| family=reuse;policy=mat;refs=2;consumer=1% | valid | 2.2766 | 2.7639 | 3.9540 | 5 | 361.8 |
| family=reuse;policy=mat;refs=2;consumer=100% | valid | 2.8460 | 2.9788 | 3.2270 | 5 | 335.7 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 2.4537 | 2.7532 | 4.0124 | 5 | 363.2 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 2.8222 | 3.0078 | 4.0435 | 5 | 332.5 |
| family=reuse;policy=mat;refs=8;consumer=1% | valid | 2.6814 | 2.9525 | 3.1955 | 5 | 338.7 |
| family=reuse;policy=mat;refs=8;consumer=100% | valid | 2.9869 | 3.2280 | 4.6192 | 5 | 309.8 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 2.3713 | 2.6584 | 3.8345 | 5 | 376.2 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 15.7784 | 16.1016 | 49.8549 | 5 | 62.1 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 0.4663 | 0.4815 | 9.0963 | 5 | 2077.1 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 2.3669 | 2.6732 | 3.4174 | 5 | 374.1 |
| family=reuse;policy=notmat;refs=2;consumer=1% | valid | 0.7652 | 0.8109 | 1.1135 | 5 | 1233.3 |
| family=reuse;policy=notmat;refs=2;consumer=100% | valid | 4.8794 | 5.0722 | 6.7326 | 5 | 197.2 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 0.8968 | 0.9861 | 1.3862 | 5 | 1014.1 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 7.2988 | 7.5602 | 8.8711 | 5 | 132.3 |
| family=reuse;policy=notmat;refs=8;consumer=1% | valid | 1.3370 | 1.4900 | 2.0615 | 5 | 671.1 |
| family=reuse;policy=notmat;refs=8;consumer=100% | valid | 9.3293 | 10.5705 | 30.6656 | 5 | 94.6 |

## Matched blended workload

- phase: `isolated`; target offered rate: 10 updates/s
- offered 0, committed 0, failed 0 over 0.000 s
- offered rate 0.000/s, completed rate 0.000/s
- per key: no keys touched
- ledger: 0 keys reconciled, 0 mismatch

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42261981`, 5 interleaved rounds
  1. family=bag;stmt=cte_filtered_notmat;threshold=20
  2. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  3. family=reuse;policy=default;refs=16;consumer=1%
  4. family=bag;stmt=cte_joined_outer_mat;threshold=100000
  5. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  6. family=bag;stmt=cte_filtered;threshold=20000
  7. family=reuse;policy=default;refs=1;consumer=100%
  8. family=bag;stmt=join_where;threshold=100000
  9. family=bag;stmt=cte_joined_outer_notmat;threshold=100000
  10. family=bag;stmt=comma_join;threshold=20
  11. family=reuse;policy=default;refs=16;consumer=100%
  12. family=bag;stmt=join_where;threshold=180000
  13. family=bag;stmt=join_on;threshold=20000
  14. family=reuse;policy=mat;refs=2;consumer=100%
  15. family=reuse;policy=default;refs=4;consumer=1%
  16. family=reuse;policy=notmat;refs=4;consumer=100%
  17. family=bag;stmt=cte_joined_outer;threshold=20000
  18. family=bag;stmt=comma_join;threshold=100000
  19. family=bag;stmt=cte_filtered_mat;threshold=100000
  20. family=bag;stmt=cte_joined_outer;threshold=200
  21. family=bag;stmt=join_on;threshold=2000
  22. family=reuse;policy=mat;refs=8;consumer=1%
  23. family=bag;stmt=cte_filtered;threshold=200
  24. family=reuse;policy=default;refs=8;consumer=100%
  25. family=reuse;policy=default;refs=8;consumer=1%
  26. family=bag;stmt=join_on;threshold=180000
  27. family=bag;stmt=cte_joined_outer_notmat;threshold=200
  28. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  29. family=bag;stmt=cte_filtered;threshold=180000
  30. family=bag;stmt=cte_joined_outer_mat;threshold=20
  31. family=bag;stmt=cte_joined_outer;threshold=100000
  32. family=bag;stmt=filtered_derived;threshold=180000
  33. family=reuse;policy=default;refs=1;consumer=1%
  34. family=bag;stmt=comma_join;threshold=180000
  35. family=bag;stmt=cte_joined_outer_notmat;threshold=20000
  36. family=reuse;policy=notmat;refs=2;consumer=1%
  37. family=bag;stmt=join_on;threshold=20
  38. family=bag;stmt=cte_filtered_mat;threshold=180000
  39. family=bag;stmt=cte_filtered_mat;threshold=20
  40. family=bag;stmt=join_where;threshold=2000
  41. family=reuse;policy=notmat;refs=16;consumer=100%
  42. family=reuse;policy=mat;refs=16;consumer=100%
  43. family=bag;stmt=comma_join;threshold=20000
  44. family=bag;stmt=cte_filtered_notmat;threshold=2000
  45. family=reuse;policy=notmat;refs=1;consumer=1%
  46. family=reuse;policy=notmat;refs=8;consumer=100%
  47. family=bag;stmt=filtered_derived;threshold=20000
  48. family=bag;stmt=join_on;threshold=100000
  49. family=bag;stmt=cte_filtered_notmat;threshold=200
  50. family=reuse;policy=mat;refs=2;consumer=1%
  51. family=reuse;policy=mat;refs=8;consumer=100%
  52. family=reuse;policy=notmat;refs=16;consumer=1%
  53. family=reuse;policy=mat;refs=4;consumer=1%
  54. family=bag;stmt=cte_filtered_mat;threshold=2000
  55. family=reuse;policy=notmat;refs=4;consumer=1%
  56. family=reuse;policy=mat;refs=1;consumer=1%
  57. family=bag;stmt=cte_filtered_notmat;threshold=20000
  58. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  59. family=bag;stmt=cte_filtered;threshold=100000
  60. family=reuse;policy=notmat;refs=1;consumer=100%
  61. family=bag;stmt=filtered_derived;threshold=200
  62. family=bag;stmt=cte_filtered_mat;threshold=200
  63. family=reuse;policy=mat;refs=1;consumer=100%
  64. family=bag;stmt=cte_filtered_mat;threshold=20000
  65. family=bag;stmt=join_on;threshold=200
  66. family=bag;stmt=cte_filtered_notmat;threshold=100000
  67. family=bag;stmt=cte_filtered;threshold=2000
  68. family=reuse;policy=default;refs=2;consumer=1%
  69. family=bag;stmt=cte_joined_outer;threshold=20
  70. family=bag;stmt=cte_joined_outer;threshold=2000
  71. family=bag;stmt=cte_joined_outer;threshold=180000
  72. family=reuse;policy=notmat;refs=2;consumer=100%
  73. family=bag;stmt=comma_join;threshold=2000
  74. family=bag;stmt=join_where;threshold=20000
  75. family=bag;stmt=cte_filtered;threshold=20
  76. family=reuse;policy=mat;refs=4;consumer=100%
  77. family=bag;stmt=filtered_derived;threshold=2000
  78. family=bag;stmt=filtered_derived;threshold=20
  79. family=reuse;policy=default;refs=4;consumer=100%
  80. family=bag;stmt=filtered_derived;threshold=100000
  81. family=reuse;policy=default;refs=2;consumer=100%
  82. family=reuse;policy=notmat;refs=8;consumer=1%
  83. family=bag;stmt=cte_joined_outer_mat;threshold=200
  84. family=bag;stmt=cte_filtered_notmat;threshold=180000
  85. family=bag;stmt=join_where;threshold=20
  86. family=bag;stmt=join_where;threshold=200
  87. family=bag;stmt=cte_joined_outer_mat;threshold=20000
  88. family=reuse;policy=mat;refs=16;consumer=1%
  89. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  90. family=bag;stmt=comma_join;threshold=200
  91. family=bag;stmt=cte_filtered_notmat;threshold=20;position=last

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
| family=bag;stmt=cte_filtered_notmat;threshold=20;position=last | j04n_cte_filtered_notmat | valid | `61819ba7d0f5f952` | policy=not_materialized; observed_cte_scan=false |
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
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `f632f422e7d95269` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `549a07108ab01ba9` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `2ce8b85a7b60c88e` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=2;consumer=1% | r_notmat_n2_c1 | valid | `28cc913d71d0d72b` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=2;consumer=100% | r_notmat_n2_c100 | valid | `53ffbc7320ed83a8` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `9f850e53560baafb` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `1888423fc6a4024a` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=8;consumer=1% | r_notmat_n8_c1 | valid | `477f7fbaf1c627cb` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=8;consumer=100% | r_notmat_n8_c100 | valid | `e2f06cbda0e4dce8` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, VERBOSE, BUFFERS, FORMAT TEXT` supported on postgres
- one host: this run used node budget 2 CPU 3g and client budget 1 CPU 2g, so budget comparisons are separate classes here, not independent hardware
- isolated phase: no concurrent writers in this run; the blended read/write phase is a separate labelled run
- physical rank-index arm pg-btree-rank; an engine whose planner ignores the index is a different physical design, not a faster one
- distributed RPC counters are N/A on single-node PostgreSQL; shared buffers and rows removed are parsed per cell where printed
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
