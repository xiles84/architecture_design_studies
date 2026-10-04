# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 91 cells, 91 valid, engine `postgres`, topology `single`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-v3-domain-pg-b3-20261005`
- run tag: `run/study-01/select-join-v3-domain-pg-b3-20261005`
- environment: `host-zenbook-ux5406sa`
- repo commit: `0f9d09212c6d52a80d8437d3b5a234f2f711ed36`
- engine version: `PostgreSQL 17.11 (Debian 17.11-1.pgdg13+2) on x86_64-pc-linux-gnu, compiled by gcc (Debian 14.2.0-19) 14.2.0, 64-bit`
- database engine image: `docker.io/library/postgres:17.11` @ `7296f210ae81031ec955dbad9a67a84fe958572a2153b8d0826a647522904dc1`
- benchmark client image: `localhost/charitybench-select:select-join-v3-domain-pg-b3-20261005` @ `8436e36bb66caa8178559d54771b88c063e26848d4ef1d1a56eb524b7f0ecca6`
- inputs digest: `c5920469c8eba10ef967d4f05126fefb46624f21a5c6434e58c35f65e32ae8c9`
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
| 20 (0.01%) | 0.4394 | 0.4860 | 0.6348 | 0.4890 | 0.6473 | 0.5371 | 0.4249 | 47.6839 | 0.6584 | 0.6212 |
| 200 (0.1%) | 0.6714 | 1.0558 | 0.7230 | 0.7156 | 0.7454 | 0.7302 | 0.9480 | 45.3663 | 0.8850 | 0.6979 |
| 2000 (1%) | 1.6659 | 1.6363 | 1.9425 | 1.8814 | 1.8377 | 1.9153 | 1.6347 | 45.7929 | 1.9699 | 1.5649 |
| 20000 (10%) | 6.7808 | 6.6044 | 6.9899 | 7.0836 | 8.9442 | 6.8254 | 6.9247 | 53.0535 | 7.3738 | 7.0641 |
| 100000 (50%) | 59.1548 | 44.0866 | 53.9677 | 52.9284 | 56.0502 | 62.7424 | 49.8375 | 59.3030 | 43.5105 | 44.1202 |
| 180000 (90%) | 70.9517 | 68.9885 | 79.6598 | 79.7020 | 88.7683 | 81.0217 | 47.5086 | 75.5793 | 84.5432 | 52.3823 |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | join_on vs join_where | 0.4394 | 0.4860 | 1.1059 | true |
| 20 | filtered_derived vs join_where | 0.4394 | 0.6348 | 1.4447 | true |
| 20 | cte_filtered vs join_where | 0.4394 | 0.4890 | 1.1128 | true |
| 20 | cte_filtered_mat vs join_where | 0.4394 | 0.6473 | 1.4731 | false |
| 20 | cte_filtered_notmat vs join_where | 0.4394 | 0.5371 | 1.2222 | true |
| 20 | cte_joined_outer vs join_where | 0.4394 | 0.4249 | 0.9669 | true |
| 20 | cte_joined_outer_mat vs join_where | 0.4394 | 47.6839 | 108.5123 | false |
| 20 | cte_joined_outer_notmat vs join_where | 0.4394 | 0.6584 | 1.4984 | true |
| 20 | comma_join vs join_where | 0.4394 | 0.6212 | 1.4136 | true |
| 200 | join_on vs join_where | 0.6714 | 1.0558 | 1.5726 | true |
| 200 | filtered_derived vs join_where | 0.6714 | 0.7230 | 1.0769 | true |
| 200 | cte_filtered vs join_where | 0.6714 | 0.7156 | 1.0659 | true |
| 200 | cte_filtered_mat vs join_where | 0.6714 | 0.7454 | 1.1104 | false |
| 200 | cte_filtered_notmat vs join_where | 0.6714 | 0.7302 | 1.0877 | true |
| 200 | cte_joined_outer vs join_where | 0.6714 | 0.9480 | 1.4121 | true |
| 200 | cte_joined_outer_mat vs join_where | 0.6714 | 45.3663 | 67.5747 | false |
| 200 | cte_joined_outer_notmat vs join_where | 0.6714 | 0.8850 | 1.3183 | true |
| 200 | comma_join vs join_where | 0.6714 | 0.6979 | 1.0396 | true |
| 2000 | join_on vs join_where | 1.6659 | 1.6363 | 0.9823 | true |
| 2000 | filtered_derived vs join_where | 1.6659 | 1.9425 | 1.1660 | true |
| 2000 | cte_filtered vs join_where | 1.6659 | 1.8814 | 1.1294 | true |
| 2000 | cte_filtered_mat vs join_where | 1.6659 | 1.8377 | 1.1032 | false |
| 2000 | cte_filtered_notmat vs join_where | 1.6659 | 1.9153 | 1.1497 | true |
| 2000 | cte_joined_outer vs join_where | 1.6659 | 1.6347 | 0.9813 | true |
| 2000 | cte_joined_outer_mat vs join_where | 1.6659 | 45.7929 | 27.4889 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 1.6659 | 1.9699 | 1.1825 | true |
| 2000 | comma_join vs join_where | 1.6659 | 1.5649 | 0.9394 | true |
| 20000 | join_on vs join_where | 6.7808 | 6.6044 | 0.9740 | true |
| 20000 | filtered_derived vs join_where | 6.7808 | 6.9899 | 1.0308 | true |
| 20000 | cte_filtered vs join_where | 6.7808 | 7.0836 | 1.0446 | true |
| 20000 | cte_filtered_mat vs join_where | 6.7808 | 8.9442 | 1.3190 | false |
| 20000 | cte_filtered_notmat vs join_where | 6.7808 | 6.8254 | 1.0066 | true |
| 20000 | cte_joined_outer vs join_where | 6.7808 | 6.9247 | 1.0212 | true |
| 20000 | cte_joined_outer_mat vs join_where | 6.7808 | 53.0535 | 7.8241 | false |
| 20000 | cte_joined_outer_notmat vs join_where | 6.7808 | 7.3738 | 1.0874 | true |
| 20000 | comma_join vs join_where | 6.7808 | 7.0641 | 1.0418 | true |
| 100000 | join_on vs join_where | 59.1548 | 44.0866 | 0.7453 | true |
| 100000 | filtered_derived vs join_where | 59.1548 | 53.9677 | 0.9123 | true |
| 100000 | cte_filtered vs join_where | 59.1548 | 52.9284 | 0.8947 | true |
| 100000 | cte_filtered_mat vs join_where | 59.1548 | 56.0502 | 0.9475 | false |
| 100000 | cte_filtered_notmat vs join_where | 59.1548 | 62.7424 | 1.0606 | true |
| 100000 | cte_joined_outer vs join_where | 59.1548 | 49.8375 | 0.8425 | true |
| 100000 | cte_joined_outer_mat vs join_where | 59.1548 | 59.3030 | 1.0025 | false |
| 100000 | cte_joined_outer_notmat vs join_where | 59.1548 | 43.5105 | 0.7355 | true |
| 100000 | comma_join vs join_where | 59.1548 | 44.1202 | 0.7458 | true |
| 180000 | join_on vs join_where | 70.9517 | 68.9885 | 0.9723 | true |
| 180000 | filtered_derived vs join_where | 70.9517 | 79.6598 | 1.1227 | true |
| 180000 | cte_filtered vs join_where | 70.9517 | 79.7020 | 1.1233 | true |
| 180000 | cte_filtered_mat vs join_where | 70.9517 | 88.7683 | 1.2511 | false |
| 180000 | cte_filtered_notmat vs join_where | 70.9517 | 81.0217 | 1.1419 | true |
| 180000 | cte_joined_outer vs join_where | 70.9517 | 47.5086 | 0.6696 | true |
| 180000 | cte_joined_outer_mat vs join_where | 70.9517 | 75.5793 | 1.0652 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 70.9517 | 84.5432 | 1.1916 | true |
| 180000 | comma_join vs join_where | 70.9517 | 52.3823 | 0.7383 | true |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_filtered_mat vs cte_filtered | 0.4890 | 0.6473 | 1.3238 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_filtered_notmat vs cte_filtered | 0.4890 | 0.5371 | 1.0983 | true | policy=not_materialized; observed_cte_scan=false |
| 20 | cte_joined_outer_mat vs cte_joined_outer | 0.4249 | 47.6839 | 112.2230 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 0.4249 | 0.6584 | 1.5496 | true | policy=not_materialized; observed_cte_scan=false |
| 200 | cte_filtered_mat vs cte_filtered | 0.7156 | 0.7454 | 1.0417 | false | policy=materialized; observed_cte_scan=true |
| 200 | cte_filtered_notmat vs cte_filtered | 0.7156 | 0.7302 | 1.0204 | true | policy=not_materialized; observed_cte_scan=false |
| 200 | cte_joined_outer_mat vs cte_joined_outer | 0.9480 | 45.3663 | 47.8546 | false | policy=materialized; observed_cte_scan=true |
| 200 | cte_joined_outer_notmat vs cte_joined_outer | 0.9480 | 0.8850 | 0.9336 | true | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_filtered_mat vs cte_filtered | 1.8814 | 1.8377 | 0.9768 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_filtered_notmat vs cte_filtered | 1.8814 | 1.9153 | 1.0180 | true | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 1.6347 | 45.7929 | 28.0125 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 1.6347 | 1.9699 | 1.2051 | true | policy=not_materialized; observed_cte_scan=false |
| 20000 | cte_filtered_mat vs cte_filtered | 7.0836 | 8.9442 | 1.2627 | false | policy=materialized; observed_cte_scan=true |
| 20000 | cte_filtered_notmat vs cte_filtered | 7.0836 | 6.8254 | 0.9636 | true | policy=not_materialized; observed_cte_scan=false |
| 20000 | cte_joined_outer_mat vs cte_joined_outer | 6.9247 | 53.0535 | 7.6615 | false | policy=materialized; observed_cte_scan=true |
| 20000 | cte_joined_outer_notmat vs cte_joined_outer | 6.9247 | 7.3738 | 1.0648 | true | policy=not_materialized; observed_cte_scan=false |
| 100000 | cte_filtered_mat vs cte_filtered | 52.9284 | 56.0502 | 1.0590 | false | policy=materialized; observed_cte_scan=true |
| 100000 | cte_filtered_notmat vs cte_filtered | 52.9284 | 62.7424 | 1.1854 | true | policy=not_materialized; observed_cte_scan=false |
| 100000 | cte_joined_outer_mat vs cte_joined_outer | 49.8375 | 59.3030 | 1.1899 | false | policy=materialized; observed_cte_scan=true |
| 100000 | cte_joined_outer_notmat vs cte_joined_outer | 49.8375 | 43.5105 | 0.8730 | true | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_filtered_mat vs cte_filtered | 79.7020 | 88.7683 | 1.1138 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_filtered_notmat vs cte_filtered | 79.7020 | 81.0217 | 1.0166 | true | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 47.5086 | 75.5793 | 1.5909 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 47.5086 | 84.5432 | 1.7795 | true | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 0.7131 | 1 | `549a07108ab01ba9` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 2.0875 | 1 | `2ce8b85a7b60c88e` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | 2.3829 | 2 | `b307c41eb54cf452` | policy=default; observed_cte_scan=true | valid |
| default | 2 | 100% | 2.6215 | 2 | `b307c41eb54cf452` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 1% | 2.8583 | 4 | `ccfe1af393deae98` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 2.9167 | 4 | `ccfe1af393deae98` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | 2.6961 | 8 | `c8c73778097931dc` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 100% | 2.8769 | 8 | `c8c73778097931dc` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 1% | 3.0774 | 16 | `b0c356323d30c425` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 3.3526 | 16 | `b0c356323d30c425` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 2.6268 | 1 | `eaa7d6058aaa5e03` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 2.6367 | 1 | `eaa7d6058aaa5e03` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | 2.3146 | 2 | `b307c41eb54cf452` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 100% | 2.6762 | 2 | `b307c41eb54cf452` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 1% | 2.4524 | 4 | `ccfe1af393deae98` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 2.5176 | 4 | `ccfe1af393deae98` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | 2.9141 | 8 | `c8c73778097931dc` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 100% | 2.8289 | 8 | `c8c73778097931dc` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 1% | 3.2741 | 16 | `b0c356323d30c425` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 3.4281 | 16 | `b0c356323d30c425` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 0.6633 | 1 | `549a07108ab01ba9` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 2.2670 | 1 | `2ce8b85a7b60c88e` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | 0.7015 | 2 | `28cc913d71d0d72b` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 100% | 4.7597 | 2 | `419a96d79d17053d` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 1% | 0.8654 | 4 | `9f850e53560baafb` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 7.4491 | 4 | `144145be26802388` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | 1.4414 | 8 | `477f7fbaf1c627cb` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 100% | 9.1435 | 8 | `2be166031907e622` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 1% | 2.4970 | 16 | `1c38b58d8cec08d3` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 14.4636 | 16 | `f4ccbc5f6dfc9fac` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 3.6837 | false |
| notmat | 1 | 1% | 0.9301 | true |
| mat | 1 | 100% | 1.2630 | false |
| notmat | 1 | 100% | 1.0860 | true |
| mat | 2 | 1% | 0.9714 | true |
| notmat | 2 | 1% | 0.2944 | false |
| mat | 2 | 100% | 1.0209 | true |
| notmat | 2 | 100% | 1.8157 | false |
| mat | 4 | 1% | 0.8580 | true |
| notmat | 4 | 1% | 0.3028 | false |
| mat | 4 | 100% | 0.8632 | true |
| notmat | 4 | 100% | 2.5540 | false |
| mat | 8 | 1% | 1.0809 | true |
| notmat | 8 | 1% | 0.5346 | false |
| mat | 8 | 100% | 0.9833 | true |
| notmat | 8 | 100% | 3.1782 | false |
| mat | 16 | 1% | 1.0639 | true |
| notmat | 16 | 1% | 0.8114 | false |
| mat | 16 | 100% | 1.0225 | true |
| notmat | 16 | 100% | 4.3142 | false |

## Repeated-design position control

| first median ms | last median ms | ratio (last/first) | same plan |
|---|---|---|---|
| 53.053546 | 51.722769 | 0.9749 | true |

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

- client: nr_periods 366, nr_throttled 294, throttled_usec 16906802 | quota 100000 100000
, memory 2147483648

- pg-single: nr_periods 366, nr_throttled 0, throttled_usec 0 | quota 200000 100000, memory 3221225472

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
| family=bag;stmt=cte_filtered_mat;threshold=100000 | valid | true | true | 0 | 0 | 100265/0 | 0 | 0 |
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
| family=bag;stmt=cte_joined_outer_mat;threshold=20000;position=last | valid | false | true | 0 | 0 | 1288/0 | 180000 | 0 |
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
| family=reuse;policy=default;refs=1;consumer=1% | valid | true | false | 0 | 0 | 32/0 | 1983 | 0 |
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
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | true | false | 0 | 0 | 30/0 | 1983 | 0 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | false | false | 0 | 0 | 0/0 | 0 | 0 |
| family=reuse;policy=notmat;refs=2;consumer=1% | valid | true | false | 0 | 0 | 60/0 | 1983 | 0 |
| family=reuse;policy=notmat;refs=2;consumer=100% | valid | false | false | 0 | 0 | 0/0 | 0 | 0 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | true | false | 0 | 0 | 96/0 | 1983 | 0 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | false | false | 0 | 0 | 0/0 | 0 | 0 |
| family=reuse;policy=notmat;refs=8;consumer=1% | valid | true | false | 0 | 0 | 256/0 | 1983 | 0 |
| family=reuse;policy=notmat;refs=8;consumer=100% | valid | false | false | 0 | 0 | 0/0 | 0 | 0 |

## Latency samples and derived throughput

_The reader is sequential, so `queries/s` is 1000/median_ms, a derived rate and not a concurrent-capacity test. min/max are the five sample extremes. Each sample includes any retry and backoff time the statement needed after a transient conflict; per-cell retry counts appear in the access-path table._

| axis | status | min ms | median ms | max ms | samples | derived queries/s |
|---|---|---|---|---|---|---|
| family=bag;stmt=comma_join;threshold=100000 | valid | 43.3794 | 44.1202 | 44.9458 | 5 | 22.7 |
| family=bag;stmt=comma_join;threshold=180000 | valid | 45.6545 | 52.3823 | 71.2890 | 5 | 19.1 |
| family=bag;stmt=comma_join;threshold=20 | valid | 0.5661 | 0.6212 | 0.7766 | 5 | 1609.8 |
| family=bag;stmt=comma_join;threshold=200 | valid | 0.6434 | 0.6979 | 0.9127 | 5 | 1432.9 |
| family=bag;stmt=comma_join;threshold=2000 | valid | 1.5390 | 1.5649 | 1.8079 | 5 | 639.0 |
| family=bag;stmt=comma_join;threshold=20000 | valid | 6.5614 | 7.0641 | 8.4650 | 5 | 141.6 |
| family=bag;stmt=cte_filtered;threshold=100000 | valid | 40.5394 | 52.9284 | 65.7476 | 5 | 18.9 |
| family=bag;stmt=cte_filtered;threshold=180000 | valid | 71.2114 | 79.7020 | 85.1716 | 5 | 12.5 |
| family=bag;stmt=cte_filtered;threshold=20 | valid | 0.4538 | 0.4890 | 0.5560 | 5 | 2045.0 |
| family=bag;stmt=cte_filtered;threshold=200 | valid | 0.6214 | 0.7156 | 0.7800 | 5 | 1397.4 |
| family=bag;stmt=cte_filtered;threshold=2000 | valid | 1.7907 | 1.8814 | 1.9053 | 5 | 531.5 |
| family=bag;stmt=cte_filtered;threshold=20000 | valid | 6.6633 | 7.0836 | 7.5031 | 5 | 141.2 |
| family=bag;stmt=cte_filtered_mat;threshold=100000 | valid | 53.1814 | 56.0502 | 60.5319 | 5 | 17.8 |
| family=bag;stmt=cte_filtered_mat;threshold=180000 | valid | 84.6759 | 88.7683 | 92.3029 | 5 | 11.3 |
| family=bag;stmt=cte_filtered_mat;threshold=20 | valid | 0.5920 | 0.6473 | 0.6657 | 5 | 1544.8 |
| family=bag;stmt=cte_filtered_mat;threshold=200 | valid | 0.7092 | 0.7454 | 0.7961 | 5 | 1341.5 |
| family=bag;stmt=cte_filtered_mat;threshold=2000 | valid | 1.7049 | 1.8377 | 2.9929 | 5 | 544.1 |
| family=bag;stmt=cte_filtered_mat;threshold=20000 | valid | 8.3892 | 8.9442 | 9.2495 | 5 | 111.8 |
| family=bag;stmt=cte_filtered_notmat;threshold=100000 | valid | 37.7100 | 62.7424 | 68.5052 | 5 | 15.9 |
| family=bag;stmt=cte_filtered_notmat;threshold=180000 | valid | 41.6372 | 81.0217 | 93.3260 | 5 | 12.3 |
| family=bag;stmt=cte_filtered_notmat;threshold=20 | valid | 0.4304 | 0.5371 | 0.8437 | 5 | 1861.9 |
| family=bag;stmt=cte_filtered_notmat;threshold=200 | valid | 0.6797 | 0.7302 | 0.7825 | 5 | 1369.5 |
| family=bag;stmt=cte_filtered_notmat;threshold=2000 | valid | 1.7966 | 1.9153 | 2.0668 | 5 | 522.1 |
| family=bag;stmt=cte_filtered_notmat;threshold=20000 | valid | 6.4970 | 6.8254 | 23.0002 | 5 | 146.5 |
| family=bag;stmt=cte_joined_outer;threshold=100000 | valid | 44.5273 | 49.8375 | 60.9934 | 5 | 20.1 |
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 44.0053 | 47.5086 | 91.3643 | 5 | 21.0 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 0.4155 | 0.4249 | 0.4688 | 5 | 2353.5 |
| family=bag;stmt=cte_joined_outer;threshold=200 | valid | 0.8469 | 0.9480 | 1.3783 | 5 | 1054.8 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 1.6252 | 1.6347 | 1.7538 | 5 | 611.7 |
| family=bag;stmt=cte_joined_outer;threshold=20000 | valid | 6.4642 | 6.9247 | 7.4849 | 5 | 144.4 |
| family=bag;stmt=cte_joined_outer_mat;threshold=100000 | valid | 57.5751 | 59.3030 | 67.8610 | 5 | 16.9 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 71.0364 | 75.5793 | 119.6301 | 5 | 13.2 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 40.0393 | 47.6839 | 49.6014 | 5 | 21.0 |
| family=bag;stmt=cte_joined_outer_mat;threshold=200 | valid | 43.1215 | 45.3663 | 50.5639 | 5 | 22.0 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 44.7735 | 45.7929 | 50.7093 | 5 | 21.8 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20000 | valid | 47.8252 | 53.0535 | 58.0539 | 5 | 18.8 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20000;position=last | valid | 49.8888 | 51.7228 | 55.6417 | 5 | 19.3 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=100000 | valid | 41.6001 | 43.5105 | 53.2475 | 5 | 23.0 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 69.9330 | 84.5432 | 96.3464 | 5 | 11.8 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 0.6220 | 0.6584 | 0.8703 | 5 | 1518.8 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=200 | valid | 0.8414 | 0.8850 | 0.9153 | 5 | 1129.9 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 1.8156 | 1.9699 | 2.2705 | 5 | 507.6 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20000 | valid | 6.9683 | 7.3738 | 36.9394 | 5 | 135.6 |
| family=bag;stmt=filtered_derived;threshold=100000 | valid | 40.1005 | 53.9677 | 69.2133 | 5 | 18.5 |
| family=bag;stmt=filtered_derived;threshold=180000 | valid | 77.7398 | 79.6598 | 81.9523 | 5 | 12.6 |
| family=bag;stmt=filtered_derived;threshold=20 | valid | 0.6049 | 0.6348 | 0.7366 | 5 | 1575.2 |
| family=bag;stmt=filtered_derived;threshold=200 | valid | 0.6117 | 0.7230 | 1.3868 | 5 | 1383.2 |
| family=bag;stmt=filtered_derived;threshold=2000 | valid | 1.8698 | 1.9425 | 2.2029 | 5 | 514.8 |
| family=bag;stmt=filtered_derived;threshold=20000 | valid | 6.7121 | 6.9899 | 9.1508 | 5 | 143.1 |
| family=bag;stmt=join_on;threshold=100000 | valid | 38.4806 | 44.0866 | 68.2960 | 5 | 22.7 |
| family=bag;stmt=join_on;threshold=180000 | valid | 44.6785 | 68.9885 | 90.0813 | 5 | 14.5 |
| family=bag;stmt=join_on;threshold=20 | valid | 0.3721 | 0.4860 | 0.5394 | 5 | 2057.7 |
| family=bag;stmt=join_on;threshold=200 | valid | 0.8575 | 1.0558 | 1.3484 | 5 | 947.2 |
| family=bag;stmt=join_on;threshold=2000 | valid | 1.5505 | 1.6363 | 1.8937 | 5 | 611.1 |
| family=bag;stmt=join_on;threshold=20000 | valid | 6.3569 | 6.6044 | 7.6145 | 5 | 151.4 |
| family=bag;stmt=join_where;threshold=100000 | valid | 39.1447 | 59.1548 | 67.4908 | 5 | 16.9 |
| family=bag;stmt=join_where;threshold=180000 | valid | 63.1303 | 70.9517 | 72.1999 | 5 | 14.1 |
| family=bag;stmt=join_where;threshold=20 | valid | 0.3839 | 0.4394 | 0.5742 | 5 | 2275.7 |
| family=bag;stmt=join_where;threshold=200 | valid | 0.6216 | 0.6714 | 0.7282 | 5 | 1489.5 |
| family=bag;stmt=join_where;threshold=2000 | valid | 1.5324 | 1.6659 | 47.3503 | 5 | 600.3 |
| family=bag;stmt=join_where;threshold=20000 | valid | 6.5206 | 6.7808 | 6.9060 | 5 | 147.5 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 3.0440 | 3.0774 | 3.7515 | 5 | 325.0 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 3.0542 | 3.3526 | 3.5671 | 5 | 298.3 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 0.6093 | 0.7131 | 0.7336 | 5 | 1402.3 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 2.0373 | 2.0875 | 2.9020 | 5 | 479.0 |
| family=reuse;policy=default;refs=2;consumer=1% | valid | 2.2908 | 2.3829 | 2.6440 | 5 | 419.7 |
| family=reuse;policy=default;refs=2;consumer=100% | valid | 2.5683 | 2.6215 | 3.1959 | 5 | 381.5 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 2.7186 | 2.8583 | 2.8769 | 5 | 349.9 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 2.7310 | 2.9167 | 3.5671 | 5 | 342.9 |
| family=reuse;policy=default;refs=8;consumer=1% | valid | 2.6191 | 2.6961 | 2.8444 | 5 | 370.9 |
| family=reuse;policy=default;refs=8;consumer=100% | valid | 2.7589 | 2.8769 | 2.9169 | 5 | 347.6 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 3.1117 | 3.2741 | 3.4179 | 5 | 305.4 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 3.3463 | 3.4281 | 3.6855 | 5 | 291.7 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 2.3747 | 2.6268 | 2.6610 | 5 | 380.7 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 2.4940 | 2.6367 | 2.7448 | 5 | 379.3 |
| family=reuse;policy=mat;refs=2;consumer=1% | valid | 2.3057 | 2.3146 | 2.5287 | 5 | 432.0 |
| family=reuse;policy=mat;refs=2;consumer=100% | valid | 2.6074 | 2.6762 | 2.7804 | 5 | 373.7 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 2.2972 | 2.4524 | 2.5942 | 5 | 407.8 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 2.5024 | 2.5176 | 2.6774 | 5 | 397.2 |
| family=reuse;policy=mat;refs=8;consumer=1% | valid | 2.8164 | 2.9141 | 3.0921 | 5 | 343.2 |
| family=reuse;policy=mat;refs=8;consumer=100% | valid | 2.7737 | 2.8289 | 3.4479 | 5 | 353.5 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 2.3663 | 2.4970 | 2.6198 | 5 | 400.5 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 13.7337 | 14.4636 | 14.7527 | 5 | 69.1 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 0.5834 | 0.6633 | 0.7905 | 5 | 1507.7 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 2.0914 | 2.2670 | 2.5412 | 5 | 441.1 |
| family=reuse;policy=notmat;refs=2;consumer=1% | valid | 0.6957 | 0.7015 | 0.8019 | 5 | 1425.5 |
| family=reuse;policy=notmat;refs=2;consumer=100% | valid | 4.3705 | 4.7597 | 4.9380 | 5 | 210.1 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 0.8298 | 0.8654 | 0.8840 | 5 | 1155.6 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 6.9068 | 7.4491 | 8.0455 | 5 | 134.2 |
| family=reuse;policy=notmat;refs=8;consumer=1% | valid | 1.3734 | 1.4414 | 1.5345 | 5 | 693.8 |
| family=reuse;policy=notmat;refs=8;consumer=100% | valid | 8.4767 | 9.1435 | 10.0119 | 5 | 109.4 |

## Matched blended workload

- phase: `blended`; target offered rate: 10 updates/s
- offered 335, committed 335, failed 0 over 33.460 s
- offered rate 10.012/s, completed rate 10.012/s
- writer operations per endpoint: [335]
- reader operations overlapping a committed write: 117 of 637
- reader retries of transient conflicts: 0
- per key: 335 keys, min 1, max 1 per key
- ledger: 335 keys reconciled, 0 mismatch
- update latency: avg 2.873 ms, min 0.824, max 27.697 over 335

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42135390`, 5 interleaved rounds
  1. family=bag;stmt=cte_joined_outer_mat;threshold=20000
  2. family=reuse;policy=mat;refs=1;consumer=100%
  3. family=reuse;policy=default;refs=2;consumer=1%
  4. family=bag;stmt=comma_join;threshold=200
  5. family=reuse;policy=mat;refs=4;consumer=100%
  6. family=bag;stmt=comma_join;threshold=180000
  7. family=bag;stmt=join_on;threshold=180000
  8. family=bag;stmt=cte_filtered_notmat;threshold=180000
  9. family=bag;stmt=comma_join;threshold=20
  10. family=bag;stmt=cte_filtered;threshold=180000
  11. family=bag;stmt=cte_filtered_mat;threshold=20
  12. family=bag;stmt=cte_joined_outer;threshold=2000
  13. family=reuse;policy=notmat;refs=16;consumer=1%
  14. family=bag;stmt=cte_filtered_notmat;threshold=200
  15. family=reuse;policy=notmat;refs=4;consumer=1%
  16. family=reuse;policy=mat;refs=2;consumer=1%
  17. family=bag;stmt=join_where;threshold=200
  18. family=bag;stmt=comma_join;threshold=20000
  19. family=bag;stmt=cte_filtered_notmat;threshold=20000
  20. family=reuse;policy=mat;refs=2;consumer=100%
  21. family=bag;stmt=cte_filtered_mat;threshold=100000
  22. family=bag;stmt=filtered_derived;threshold=100000
  23. family=reuse;policy=mat;refs=8;consumer=1%
  24. family=bag;stmt=cte_joined_outer_mat;threshold=20
  25. family=bag;stmt=filtered_derived;threshold=2000
  26. family=bag;stmt=cte_joined_outer;threshold=180000
  27. family=bag;stmt=cte_joined_outer;threshold=20000
  28. family=bag;stmt=cte_joined_outer_notmat;threshold=200
  29. family=bag;stmt=join_on;threshold=20
  30. family=bag;stmt=join_where;threshold=20000
  31. family=reuse;policy=mat;refs=1;consumer=1%
  32. family=reuse;policy=notmat;refs=2;consumer=100%
  33. family=reuse;policy=mat;refs=16;consumer=100%
  34. family=bag;stmt=comma_join;threshold=100000
  35. family=bag;stmt=join_on;threshold=200
  36. family=bag;stmt=filtered_derived;threshold=200
  37. family=bag;stmt=cte_filtered_notmat;threshold=100000
  38. family=reuse;policy=default;refs=4;consumer=100%
  39. family=reuse;policy=default;refs=1;consumer=100%
  40. family=bag;stmt=filtered_derived;threshold=180000
  41. family=reuse;policy=default;refs=4;consumer=1%
  42. family=bag;stmt=join_on;threshold=2000
  43. family=bag;stmt=cte_filtered_mat;threshold=180000
  44. family=bag;stmt=cte_filtered_notmat;threshold=2000
  45. family=reuse;policy=default;refs=8;consumer=100%
  46. family=bag;stmt=cte_joined_outer_notmat;threshold=100000
  47. family=reuse;policy=notmat;refs=4;consumer=100%
  48. family=reuse;policy=notmat;refs=16;consumer=100%
  49. family=bag;stmt=join_where;threshold=180000
  50. family=bag;stmt=filtered_derived;threshold=20
  51. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  52. family=reuse;policy=notmat;refs=2;consumer=1%
  53. family=bag;stmt=cte_joined_outer_mat;threshold=200
  54. family=bag;stmt=cte_joined_outer_mat;threshold=100000
  55. family=bag;stmt=cte_joined_outer_notmat;threshold=20000
  56. family=reuse;policy=default;refs=2;consumer=100%
  57. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  58. family=reuse;policy=notmat;refs=1;consumer=1%
  59. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  60. family=bag;stmt=join_where;threshold=100000
  61. family=reuse;policy=mat;refs=16;consumer=1%
  62. family=bag;stmt=cte_filtered;threshold=20
  63. family=reuse;policy=default;refs=16;consumer=1%
  64. family=bag;stmt=cte_joined_outer;threshold=20
  65. family=reuse;policy=mat;refs=8;consumer=100%
  66. family=bag;stmt=cte_filtered_mat;threshold=2000
  67. family=reuse;policy=mat;refs=4;consumer=1%
  68. family=bag;stmt=cte_joined_outer;threshold=100000
  69. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  70. family=bag;stmt=join_where;threshold=20
  71. family=bag;stmt=cte_filtered_notmat;threshold=20
  72. family=bag;stmt=cte_filtered_mat;threshold=200
  73. family=bag;stmt=join_where;threshold=2000
  74. family=reuse;policy=default;refs=16;consumer=100%
  75. family=bag;stmt=cte_filtered;threshold=100000
  76. family=reuse;policy=default;refs=1;consumer=1%
  77. family=reuse;policy=default;refs=8;consumer=1%
  78. family=bag;stmt=comma_join;threshold=2000
  79. family=bag;stmt=cte_filtered_mat;threshold=20000
  80. family=bag;stmt=cte_filtered;threshold=2000
  81. family=reuse;policy=notmat;refs=8;consumer=100%
  82. family=bag;stmt=filtered_derived;threshold=20000
  83. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  84. family=bag;stmt=cte_filtered;threshold=200
  85. family=reuse;policy=notmat;refs=1;consumer=100%
  86. family=reuse;policy=notmat;refs=8;consumer=1%
  87. family=bag;stmt=join_on;threshold=20000
  88. family=bag;stmt=cte_joined_outer;threshold=200
  89. family=bag;stmt=cte_filtered;threshold=20000
  90. family=bag;stmt=join_on;threshold=100000
  91. family=bag;stmt=cte_joined_outer_mat;threshold=20000;position=last

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
| family=bag;stmt=cte_joined_outer_mat;threshold=20000;position=last | j05m_cte_joined_outer_mat | valid | `7f87bace4cf1da4c` | policy=materialized; observed_cte_scan=true |
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
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `f4ccbc5f6dfc9fac` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `549a07108ab01ba9` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `2ce8b85a7b60c88e` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=2;consumer=1% | r_notmat_n2_c1 | valid | `28cc913d71d0d72b` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=2;consumer=100% | r_notmat_n2_c100 | valid | `419a96d79d17053d` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `9f850e53560baafb` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `144145be26802388` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=8;consumer=1% | r_notmat_n8_c1 | valid | `477f7fbaf1c627cb` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=8;consumer=100% | r_notmat_n8_c100 | valid | `2be166031907e622` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, VERBOSE, BUFFERS, FORMAT TEXT` supported on postgres
- one host: this run used node budget 2 CPU 3g and client budget 1 CPU 2g, so budget comparisons are separate classes here, not independent hardware
- one client process with one fixed maintenance writer at 10 updates/s; readers are sequential; the reported sample latency includes retry and backoff time
- physical rank-index arm pg-btree-rank; an engine whose planner ignores the index is a different physical design, not a faster one
- distributed RPC counters are N/A on single-node PostgreSQL; shared buffers and rows removed are parsed per cell where printed
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
