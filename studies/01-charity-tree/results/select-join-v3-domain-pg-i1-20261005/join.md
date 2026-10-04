# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 91 cells, 91 valid, engine `postgres`, topology `single`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-v3-domain-pg-i1-20261005`
- run tag: `run/study-01/select-join-v3-domain-pg-i1-20261005`
- environment: `host-zenbook-ux5406sa`
- repo commit: `a209d4be775397a7721e736d228d5215cab8f3bb`
- engine version: `PostgreSQL 17.11 (Debian 17.11-1.pgdg13+2) on x86_64-pc-linux-gnu, compiled by gcc (Debian 14.2.0-19) 14.2.0, 64-bit`
- database engine image: `docker.io/library/postgres:17.11` @ `7296f210ae81031ec955dbad9a67a84fe958572a2153b8d0826a647522904dc1`
- benchmark client image: `localhost/charitybench-select:select-join-v3-domain-pg-i1-20261005` @ `8436e36bb66caa8178559d54771b88c063e26848d4ef1d1a56eb524b7f0ecca6`
- inputs digest: `5f5433eefa6ebc670bffc8458178151b210467ddd9d90faadc91e74a6cbbc3d2`
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
| 20 (0.01%) | 0.4520 | 0.5058 | 0.6753 | 0.6697 | 0.6529 | 0.4735 | 0.4569 | 51.4304 | 0.5120 | 0.6137 |
| 200 (0.1%) | 0.8891 | 1.0189 | 0.8907 | 0.7424 | 0.8258 | 0.8451 | 0.7964 | 49.5737 | 0.7399 | 0.7610 |
| 2000 (1%) | 2.6115 | 1.4608 | 1.9061 | 1.6372 | 1.9423 | 1.9220 | 2.2694 | 47.1928 | 1.6927 | 1.7917 |
| 20000 (10%) | 7.6728 | 7.5633 | 6.9022 | 7.4681 | 9.4570 | 6.9693 | 7.3746 | 54.9229 | 7.4877 | 6.9428 |
| 100000 (50%) | 52.2563 | 51.9661 | 50.1330 | 56.9868 | 66.3234 | 50.8794 | 62.5575 | 64.1920 | 48.6364 | 45.8468 |
| 180000 (90%) | 69.8956 | 47.7103 | 48.9875 | 83.7385 | 74.5401 | 81.6412 | 80.3031 | 80.0276 | 77.2632 | 68.8424 |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | join_on vs join_where | 0.4520 | 0.5058 | 1.1190 | true |
| 20 | filtered_derived vs join_where | 0.4520 | 0.6753 | 1.4941 | true |
| 20 | cte_filtered vs join_where | 0.4520 | 0.6697 | 1.4817 | true |
| 20 | cte_filtered_mat vs join_where | 0.4520 | 0.6529 | 1.4445 | false |
| 20 | cte_filtered_notmat vs join_where | 0.4520 | 0.4735 | 1.0477 | true |
| 20 | cte_joined_outer vs join_where | 0.4520 | 0.4569 | 1.0109 | true |
| 20 | cte_joined_outer_mat vs join_where | 0.4520 | 51.4304 | 113.7859 | false |
| 20 | cte_joined_outer_notmat vs join_where | 0.4520 | 0.5120 | 1.1327 | true |
| 20 | comma_join vs join_where | 0.4520 | 0.6137 | 1.3577 | true |
| 200 | join_on vs join_where | 0.8891 | 1.0189 | 1.1459 | true |
| 200 | filtered_derived vs join_where | 0.8891 | 0.8907 | 1.0018 | true |
| 200 | cte_filtered vs join_where | 0.8891 | 0.7424 | 0.8350 | true |
| 200 | cte_filtered_mat vs join_where | 0.8891 | 0.8258 | 0.9288 | false |
| 200 | cte_filtered_notmat vs join_where | 0.8891 | 0.8451 | 0.9505 | true |
| 200 | cte_joined_outer vs join_where | 0.8891 | 0.7964 | 0.8956 | true |
| 200 | cte_joined_outer_mat vs join_where | 0.8891 | 49.5737 | 55.7541 | false |
| 200 | cte_joined_outer_notmat vs join_where | 0.8891 | 0.7399 | 0.8321 | true |
| 200 | comma_join vs join_where | 0.8891 | 0.7610 | 0.8558 | true |
| 2000 | join_on vs join_where | 2.6115 | 1.4608 | 0.5594 | true |
| 2000 | filtered_derived vs join_where | 2.6115 | 1.9061 | 0.7299 | true |
| 2000 | cte_filtered vs join_where | 2.6115 | 1.6372 | 0.6269 | true |
| 2000 | cte_filtered_mat vs join_where | 2.6115 | 1.9423 | 0.7437 | false |
| 2000 | cte_filtered_notmat vs join_where | 2.6115 | 1.9220 | 0.7360 | true |
| 2000 | cte_joined_outer vs join_where | 2.6115 | 2.2694 | 0.8690 | true |
| 2000 | cte_joined_outer_mat vs join_where | 2.6115 | 47.1928 | 18.0709 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 2.6115 | 1.6927 | 0.6482 | true |
| 2000 | comma_join vs join_where | 2.6115 | 1.7917 | 0.6861 | true |
| 20000 | join_on vs join_where | 7.6728 | 7.5633 | 0.9857 | true |
| 20000 | filtered_derived vs join_where | 7.6728 | 6.9022 | 0.8996 | true |
| 20000 | cte_filtered vs join_where | 7.6728 | 7.4681 | 0.9733 | true |
| 20000 | cte_filtered_mat vs join_where | 7.6728 | 9.4570 | 1.2325 | false |
| 20000 | cte_filtered_notmat vs join_where | 7.6728 | 6.9693 | 0.9083 | true |
| 20000 | cte_joined_outer vs join_where | 7.6728 | 7.3746 | 0.9611 | true |
| 20000 | cte_joined_outer_mat vs join_where | 7.6728 | 54.9229 | 7.1582 | false |
| 20000 | cte_joined_outer_notmat vs join_where | 7.6728 | 7.4877 | 0.9759 | true |
| 20000 | comma_join vs join_where | 7.6728 | 6.9428 | 0.9049 | true |
| 100000 | join_on vs join_where | 52.2563 | 51.9661 | 0.9944 | true |
| 100000 | filtered_derived vs join_where | 52.2563 | 50.1330 | 0.9594 | true |
| 100000 | cte_filtered vs join_where | 52.2563 | 56.9868 | 1.0905 | true |
| 100000 | cte_filtered_mat vs join_where | 52.2563 | 66.3234 | 1.2692 | false |
| 100000 | cte_filtered_notmat vs join_where | 52.2563 | 50.8794 | 0.9737 | true |
| 100000 | cte_joined_outer vs join_where | 52.2563 | 62.5575 | 1.1971 | true |
| 100000 | cte_joined_outer_mat vs join_where | 52.2563 | 64.1920 | 1.2284 | false |
| 100000 | cte_joined_outer_notmat vs join_where | 52.2563 | 48.6364 | 0.9307 | true |
| 100000 | comma_join vs join_where | 52.2563 | 45.8468 | 0.8773 | true |
| 180000 | join_on vs join_where | 69.8956 | 47.7103 | 0.6826 | true |
| 180000 | filtered_derived vs join_where | 69.8956 | 48.9875 | 0.7009 | true |
| 180000 | cte_filtered vs join_where | 69.8956 | 83.7385 | 1.1981 | true |
| 180000 | cte_filtered_mat vs join_where | 69.8956 | 74.5401 | 1.0664 | false |
| 180000 | cte_filtered_notmat vs join_where | 69.8956 | 81.6412 | 1.1680 | true |
| 180000 | cte_joined_outer vs join_where | 69.8956 | 80.3031 | 1.1489 | true |
| 180000 | cte_joined_outer_mat vs join_where | 69.8956 | 80.0276 | 1.1450 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 69.8956 | 77.2632 | 1.1054 | true |
| 180000 | comma_join vs join_where | 69.8956 | 68.8424 | 0.9849 | true |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_filtered_mat vs cte_filtered | 0.6697 | 0.6529 | 0.9749 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_filtered_notmat vs cte_filtered | 0.6697 | 0.4735 | 0.7071 | true | policy=not_materialized; observed_cte_scan=false |
| 20 | cte_joined_outer_mat vs cte_joined_outer | 0.4569 | 51.4304 | 112.5552 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 0.4569 | 0.5120 | 1.1204 | true | policy=not_materialized; observed_cte_scan=false |
| 200 | cte_filtered_mat vs cte_filtered | 0.7424 | 0.8258 | 1.1124 | false | policy=materialized; observed_cte_scan=true |
| 200 | cte_filtered_notmat vs cte_filtered | 0.7424 | 0.8451 | 1.1384 | true | policy=not_materialized; observed_cte_scan=false |
| 200 | cte_joined_outer_mat vs cte_joined_outer | 0.7964 | 49.5737 | 62.2503 | false | policy=materialized; observed_cte_scan=true |
| 200 | cte_joined_outer_notmat vs cte_joined_outer | 0.7964 | 0.7399 | 0.9291 | true | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_filtered_mat vs cte_filtered | 1.6372 | 1.9423 | 1.1864 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_filtered_notmat vs cte_filtered | 1.6372 | 1.9220 | 1.1740 | true | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 2.2694 | 47.1928 | 20.7949 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 2.2694 | 1.6927 | 0.7459 | true | policy=not_materialized; observed_cte_scan=false |
| 20000 | cte_filtered_mat vs cte_filtered | 7.4681 | 9.4570 | 1.2663 | false | policy=materialized; observed_cte_scan=true |
| 20000 | cte_filtered_notmat vs cte_filtered | 7.4681 | 6.9693 | 0.9332 | true | policy=not_materialized; observed_cte_scan=false |
| 20000 | cte_joined_outer_mat vs cte_joined_outer | 7.3746 | 54.9229 | 7.4476 | false | policy=materialized; observed_cte_scan=true |
| 20000 | cte_joined_outer_notmat vs cte_joined_outer | 7.3746 | 7.4877 | 1.0153 | true | policy=not_materialized; observed_cte_scan=false |
| 100000 | cte_filtered_mat vs cte_filtered | 56.9868 | 66.3234 | 1.1638 | false | policy=materialized; observed_cte_scan=true |
| 100000 | cte_filtered_notmat vs cte_filtered | 56.9868 | 50.8794 | 0.8928 | true | policy=not_materialized; observed_cte_scan=false |
| 100000 | cte_joined_outer_mat vs cte_joined_outer | 62.5575 | 64.1920 | 1.0261 | false | policy=materialized; observed_cte_scan=true |
| 100000 | cte_joined_outer_notmat vs cte_joined_outer | 62.5575 | 48.6364 | 0.7775 | true | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_filtered_mat vs cte_filtered | 83.7385 | 74.5401 | 0.8902 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_filtered_notmat vs cte_filtered | 83.7385 | 81.6412 | 0.9750 | true | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 80.3031 | 80.0276 | 0.9966 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 80.3031 | 77.2632 | 0.9621 | true | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 0.4675 | 1 | `549a07108ab01ba9` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 2.4483 | 1 | `2ce8b85a7b60c88e` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | 2.1693 | 2 | `b307c41eb54cf452` | policy=default; observed_cte_scan=true | valid |
| default | 2 | 100% | 2.3762 | 2 | `b307c41eb54cf452` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 1% | 2.3183 | 4 | `ccfe1af393deae98` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 2.9495 | 4 | `ccfe1af393deae98` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | 3.0337 | 8 | `c8c73778097931dc` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 100% | 2.8433 | 8 | `c8c73778097931dc` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 1% | 3.0917 | 16 | `b0c356323d30c425` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 3.6390 | 16 | `b0c356323d30c425` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 2.5203 | 1 | `eaa7d6058aaa5e03` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 2.4870 | 1 | `eaa7d6058aaa5e03` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | 2.5440 | 2 | `b307c41eb54cf452` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 100% | 2.4921 | 2 | `b307c41eb54cf452` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 1% | 2.7703 | 4 | `ccfe1af393deae98` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 2.8452 | 4 | `ccfe1af393deae98` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | 2.7965 | 8 | `c8c73778097931dc` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 100% | 2.9734 | 8 | `c8c73778097931dc` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 1% | 3.4366 | 16 | `b0c356323d30c425` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 3.7431 | 16 | `b0c356323d30c425` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 0.6558 | 1 | `549a07108ab01ba9` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 2.6986 | 1 | `2ce8b85a7b60c88e` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | 0.8398 | 2 | `28cc913d71d0d72b` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 100% | 4.6751 | 2 | `fc4d3d1577720fdb` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 1% | 0.8894 | 4 | `9f850e53560baafb` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 6.9953 | 4 | `9608e3d297c0b2d8` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | 1.4123 | 8 | `477f7fbaf1c627cb` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 100% | 9.7695 | 8 | `2f8f777dad660f79` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 1% | 2.2139 | 16 | `1c38b58d8cec08d3` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 15.6528 | 16 | `596a11d1cf99201c` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 5.3909 | false |
| notmat | 1 | 1% | 1.4027 | true |
| mat | 1 | 100% | 1.0158 | false |
| notmat | 1 | 100% | 1.1022 | true |
| mat | 2 | 1% | 1.1727 | true |
| notmat | 2 | 1% | 0.3871 | false |
| mat | 2 | 100% | 1.0488 | true |
| notmat | 2 | 100% | 1.9674 | false |
| mat | 4 | 1% | 1.1950 | true |
| notmat | 4 | 1% | 0.3836 | false |
| mat | 4 | 100% | 0.9646 | true |
| notmat | 4 | 100% | 2.3717 | false |
| mat | 8 | 1% | 0.9218 | true |
| notmat | 8 | 1% | 0.4656 | false |
| mat | 8 | 100% | 1.0458 | true |
| notmat | 8 | 100% | 3.4359 | false |
| mat | 16 | 1% | 1.1116 | true |
| notmat | 16 | 1% | 0.7161 | false |
| mat | 16 | 100% | 1.0286 | true |
| notmat | 16 | 100% | 4.3015 | false |

## Repeated-design position control

| first median ms | last median ms | ratio (last/first) | same plan |
|---|---|---|---|
| 7.672757 | 9.045127 | 1.1789 | true |

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

- client: nr_periods 379, nr_throttled 306, throttled_usec 15255077 | quota 100000 100000
, memory 2147483648

- pg-single: nr_periods 368, nr_throttled 0, throttled_usec 0 | quota 200000 100000, memory 3221225472

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
| family=bag;stmt=join_where;threshold=20000;position=last | valid | false | true | 0 | 0 | 1349/0 | 0 | 0 |
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
| family=bag;stmt=comma_join;threshold=100000 | valid | 44.6875 | 45.8468 | 54.5541 | 5 | 21.8 |
| family=bag;stmt=comma_join;threshold=180000 | valid | 47.0511 | 68.8424 | 129.2178 | 5 | 14.5 |
| family=bag;stmt=comma_join;threshold=20 | valid | 0.5870 | 0.6137 | 0.7114 | 5 | 1629.6 |
| family=bag;stmt=comma_join;threshold=200 | valid | 0.6588 | 0.7610 | 0.9187 | 5 | 1314.1 |
| family=bag;stmt=comma_join;threshold=2000 | valid | 1.4891 | 1.7917 | 2.2198 | 5 | 558.1 |
| family=bag;stmt=comma_join;threshold=20000 | valid | 6.5602 | 6.9428 | 29.9032 | 5 | 144.0 |
| family=bag;stmt=cte_filtered;threshold=100000 | valid | 45.3858 | 56.9868 | 63.4946 | 5 | 17.5 |
| family=bag;stmt=cte_filtered;threshold=180000 | valid | 75.5095 | 83.7385 | 89.8175 | 5 | 11.9 |
| family=bag;stmt=cte_filtered;threshold=20 | valid | 0.6093 | 0.6697 | 1.0207 | 5 | 1493.2 |
| family=bag;stmt=cte_filtered;threshold=200 | valid | 0.7223 | 0.7424 | 0.7890 | 5 | 1347.0 |
| family=bag;stmt=cte_filtered;threshold=2000 | valid | 1.5659 | 1.6372 | 1.8302 | 5 | 610.8 |
| family=bag;stmt=cte_filtered;threshold=20000 | valid | 6.6269 | 7.4681 | 10.3755 | 5 | 133.9 |
| family=bag;stmt=cte_filtered_mat;threshold=100000 | valid | 55.9670 | 66.3234 | 67.7005 | 5 | 15.1 |
| family=bag;stmt=cte_filtered_mat;threshold=180000 | valid | 64.2198 | 74.5401 | 90.2399 | 5 | 13.4 |
| family=bag;stmt=cte_filtered_mat;threshold=20 | valid | 0.4539 | 0.6529 | 0.7896 | 5 | 1531.7 |
| family=bag;stmt=cte_filtered_mat;threshold=200 | valid | 0.6858 | 0.8258 | 0.9794 | 5 | 1210.9 |
| family=bag;stmt=cte_filtered_mat;threshold=2000 | valid | 1.6528 | 1.9423 | 2.2718 | 5 | 514.9 |
| family=bag;stmt=cte_filtered_mat;threshold=20000 | valid | 8.6424 | 9.4570 | 35.6513 | 5 | 105.7 |
| family=bag;stmt=cte_filtered_notmat;threshold=100000 | valid | 40.0337 | 50.8794 | 68.6364 | 5 | 19.7 |
| family=bag;stmt=cte_filtered_notmat;threshold=180000 | valid | 44.1458 | 81.6412 | 82.7131 | 5 | 12.2 |
| family=bag;stmt=cte_filtered_notmat;threshold=20 | valid | 0.4434 | 0.4735 | 0.8921 | 5 | 2111.7 |
| family=bag;stmt=cte_filtered_notmat;threshold=200 | valid | 0.7120 | 0.8451 | 0.9914 | 5 | 1183.2 |
| family=bag;stmt=cte_filtered_notmat;threshold=2000 | valid | 1.8369 | 1.9220 | 2.2851 | 5 | 520.3 |
| family=bag;stmt=cte_filtered_notmat;threshold=20000 | valid | 6.8485 | 6.9693 | 7.4970 | 5 | 143.5 |
| family=bag;stmt=cte_joined_outer;threshold=100000 | valid | 56.3679 | 62.5575 | 73.6930 | 5 | 16.0 |
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 46.9436 | 80.3031 | 99.5804 | 5 | 12.5 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 0.4150 | 0.4569 | 0.6045 | 5 | 2188.5 |
| family=bag;stmt=cte_joined_outer;threshold=200 | valid | 0.6645 | 0.7964 | 0.9287 | 5 | 1255.7 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 1.8399 | 2.2694 | 2.6233 | 5 | 440.6 |
| family=bag;stmt=cte_joined_outer;threshold=20000 | valid | 6.6901 | 7.3746 | 11.0726 | 5 | 135.6 |
| family=bag;stmt=cte_joined_outer_mat;threshold=100000 | valid | 61.0915 | 64.1920 | 66.4915 | 5 | 15.6 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 72.4017 | 80.0276 | 99.9394 | 5 | 12.5 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 48.3299 | 51.4304 | 53.3119 | 5 | 19.4 |
| family=bag;stmt=cte_joined_outer_mat;threshold=200 | valid | 48.9839 | 49.5737 | 53.6404 | 5 | 20.2 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 45.3699 | 47.1928 | 52.9497 | 5 | 21.2 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20000 | valid | 49.1412 | 54.9229 | 60.8075 | 5 | 18.2 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=100000 | valid | 42.7936 | 48.6364 | 61.6433 | 5 | 20.6 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 46.1624 | 77.2632 | 86.9012 | 5 | 12.9 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 0.4864 | 0.5120 | 0.8007 | 5 | 1953.3 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=200 | valid | 0.7235 | 0.7399 | 35.0586 | 5 | 1351.5 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 1.6133 | 1.6927 | 2.3187 | 5 | 590.8 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20000 | valid | 7.2018 | 7.4877 | 8.5716 | 5 | 133.6 |
| family=bag;stmt=filtered_derived;threshold=100000 | valid | 48.5392 | 50.1330 | 64.0693 | 5 | 19.9 |
| family=bag;stmt=filtered_derived;threshold=180000 | valid | 44.4486 | 48.9875 | 51.3459 | 5 | 20.4 |
| family=bag;stmt=filtered_derived;threshold=20 | valid | 0.5818 | 0.6753 | 0.6796 | 5 | 1480.7 |
| family=bag;stmt=filtered_derived;threshold=200 | valid | 0.7985 | 0.8907 | 1.3862 | 5 | 1122.7 |
| family=bag;stmt=filtered_derived;threshold=2000 | valid | 1.8923 | 1.9061 | 2.6397 | 5 | 524.6 |
| family=bag;stmt=filtered_derived;threshold=20000 | valid | 6.5612 | 6.9022 | 31.9001 | 5 | 144.9 |
| family=bag;stmt=join_on;threshold=100000 | valid | 43.2192 | 51.9661 | 58.7803 | 5 | 19.2 |
| family=bag;stmt=join_on;threshold=180000 | valid | 44.0564 | 47.7103 | 49.8056 | 5 | 21.0 |
| family=bag;stmt=join_on;threshold=20 | valid | 0.4309 | 0.5058 | 0.6454 | 5 | 1977.1 |
| family=bag;stmt=join_on;threshold=200 | valid | 0.8738 | 1.0189 | 12.9268 | 5 | 981.5 |
| family=bag;stmt=join_on;threshold=2000 | valid | 1.3367 | 1.4608 | 1.6516 | 5 | 684.5 |
| family=bag;stmt=join_on;threshold=20000 | valid | 6.4614 | 7.5633 | 16.4175 | 5 | 132.2 |
| family=bag;stmt=join_where;threshold=100000 | valid | 41.8664 | 52.2563 | 54.3415 | 5 | 19.1 |
| family=bag;stmt=join_where;threshold=180000 | valid | 51.3316 | 69.8956 | 78.5772 | 5 | 14.3 |
| family=bag;stmt=join_where;threshold=20 | valid | 0.3764 | 0.4520 | 0.4935 | 5 | 2212.4 |
| family=bag;stmt=join_where;threshold=200 | valid | 0.6876 | 0.8891 | 18.2996 | 5 | 1124.7 |
| family=bag;stmt=join_where;threshold=2000 | valid | 1.6680 | 2.6115 | 2.8649 | 5 | 382.9 |
| family=bag;stmt=join_where;threshold=20000 | valid | 6.6655 | 7.6728 | 29.7161 | 5 | 130.3 |
| family=bag;stmt=join_where;threshold=20000;position=last | valid | 6.5311 | 9.0451 | 37.4924 | 5 | 110.6 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 2.8690 | 3.0917 | 3.4179 | 5 | 323.4 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 3.5640 | 3.6390 | 3.7593 | 5 | 274.8 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 0.3983 | 0.4675 | 0.4763 | 5 | 2139.0 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 2.2616 | 2.4483 | 2.6588 | 5 | 408.4 |
| family=reuse;policy=default;refs=2;consumer=1% | valid | 2.0284 | 2.1693 | 2.5262 | 5 | 461.0 |
| family=reuse;policy=default;refs=2;consumer=100% | valid | 2.2837 | 2.3762 | 2.6138 | 5 | 420.8 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 2.1811 | 2.3183 | 2.7267 | 5 | 431.4 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 2.7538 | 2.9495 | 3.0765 | 5 | 339.0 |
| family=reuse;policy=default;refs=8;consumer=1% | valid | 2.7843 | 3.0337 | 3.1410 | 5 | 329.6 |
| family=reuse;policy=default;refs=8;consumer=100% | valid | 2.7505 | 2.8433 | 3.0701 | 5 | 351.7 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 3.3048 | 3.4366 | 3.4450 | 5 | 291.0 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 3.4622 | 3.7431 | 4.4649 | 5 | 267.2 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 2.2381 | 2.5203 | 3.5255 | 5 | 396.8 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 2.3215 | 2.4870 | 3.1461 | 5 | 402.1 |
| family=reuse;policy=mat;refs=2;consumer=1% | valid | 2.2561 | 2.5440 | 3.2908 | 5 | 393.1 |
| family=reuse;policy=mat;refs=2;consumer=100% | valid | 2.3594 | 2.4921 | 2.5776 | 5 | 401.3 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 2.6206 | 2.7703 | 2.8823 | 5 | 361.0 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 2.6756 | 2.8452 | 3.5856 | 5 | 351.5 |
| family=reuse;policy=mat;refs=8;consumer=1% | valid | 2.7734 | 2.7965 | 2.9748 | 5 | 357.6 |
| family=reuse;policy=mat;refs=8;consumer=100% | valid | 2.6184 | 2.9734 | 3.6377 | 5 | 336.3 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 2.1539 | 2.2139 | 2.3209 | 5 | 451.7 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 14.8995 | 15.6528 | 17.0191 | 5 | 63.9 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 0.6061 | 0.6558 | 0.9897 | 5 | 1524.9 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 2.4549 | 2.6986 | 3.1081 | 5 | 370.6 |
| family=reuse;policy=notmat;refs=2;consumer=1% | valid | 0.8351 | 0.8398 | 0.8540 | 5 | 1190.8 |
| family=reuse;policy=notmat;refs=2;consumer=100% | valid | 4.2848 | 4.6751 | 6.3312 | 5 | 213.9 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 0.8538 | 0.8894 | 0.9305 | 5 | 1124.4 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 6.6042 | 6.9953 | 8.3993 | 5 | 143.0 |
| family=reuse;policy=notmat;refs=8;consumer=1% | valid | 1.3230 | 1.4123 | 2.4997 | 5 | 708.0 |
| family=reuse;policy=notmat;refs=8;consumer=100% | valid | 8.6556 | 9.7695 | 9.8934 | 5 | 102.4 |

## Matched blended workload

- phase: `isolated`; target offered rate: 10 updates/s
- offered 0, committed 0, failed 0 over 0.000 s
- offered rate 0.000/s, completed rate 0.000/s
- per key: no keys touched
- ledger: 0 keys reconciled, 0 mismatch

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42364181`, 5 interleaved rounds
  1. family=bag;stmt=join_where;threshold=20000
  2. family=bag;stmt=join_on;threshold=200
  3. family=bag;stmt=cte_joined_outer_mat;threshold=20000
  4. family=bag;stmt=cte_joined_outer_mat;threshold=20
  5. family=bag;stmt=cte_filtered_mat;threshold=180000
  6. family=reuse;policy=default;refs=1;consumer=100%
  7. family=bag;stmt=cte_filtered;threshold=200
  8. family=bag;stmt=cte_filtered;threshold=180000
  9. family=bag;stmt=cte_filtered_notmat;threshold=100000
  10. family=reuse;policy=notmat;refs=2;consumer=1%
  11. family=reuse;policy=mat;refs=8;consumer=1%
  12. family=bag;stmt=join_on;threshold=20000
  13. family=bag;stmt=cte_filtered_notmat;threshold=2000
  14. family=reuse;policy=mat;refs=1;consumer=100%
  15. family=reuse;policy=mat;refs=1;consumer=1%
  16. family=bag;stmt=cte_filtered_notmat;threshold=20
  17. family=bag;stmt=comma_join;threshold=2000
  18. family=bag;stmt=comma_join;threshold=180000
  19. family=bag;stmt=filtered_derived;threshold=20
  20. family=bag;stmt=cte_filtered_notmat;threshold=200
  21. family=bag;stmt=cte_joined_outer;threshold=100000
  22. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  23. family=reuse;policy=notmat;refs=16;consumer=100%
  24. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  25. family=bag;stmt=join_where;threshold=180000
  26. family=bag;stmt=cte_filtered_notmat;threshold=20000
  27. family=bag;stmt=cte_filtered_mat;threshold=100000
  28. family=bag;stmt=cte_filtered;threshold=20
  29. family=bag;stmt=cte_joined_outer;threshold=20
  30. family=bag;stmt=cte_joined_outer_notmat;threshold=100000
  31. family=reuse;policy=mat;refs=16;consumer=100%
  32. family=bag;stmt=cte_joined_outer_mat;threshold=100000
  33. family=bag;stmt=cte_joined_outer;threshold=2000
  34. family=bag;stmt=cte_filtered_notmat;threshold=180000
  35. family=reuse;policy=mat;refs=16;consumer=1%
  36. family=reuse;policy=default;refs=8;consumer=100%
  37. family=bag;stmt=cte_filtered_mat;threshold=200
  38. family=reuse;policy=notmat;refs=8;consumer=1%
  39. family=bag;stmt=filtered_derived;threshold=100000
  40. family=bag;stmt=filtered_derived;threshold=2000
  41. family=reuse;policy=default;refs=2;consumer=100%
  42. family=reuse;policy=default;refs=1;consumer=1%
  43. family=bag;stmt=cte_joined_outer_notmat;threshold=20000
  44. family=bag;stmt=cte_filtered;threshold=100000
  45. family=reuse;policy=default;refs=8;consumer=1%
  46. family=bag;stmt=filtered_derived;threshold=20000
  47. family=bag;stmt=join_where;threshold=2000
  48. family=bag;stmt=cte_filtered_mat;threshold=2000
  49. family=bag;stmt=comma_join;threshold=20000
  50. family=bag;stmt=comma_join;threshold=20
  51. family=bag;stmt=cte_filtered;threshold=20000
  52. family=reuse;policy=notmat;refs=1;consumer=1%
  53. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  54. family=bag;stmt=join_on;threshold=100000
  55. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  56. family=bag;stmt=filtered_derived;threshold=200
  57. family=bag;stmt=filtered_derived;threshold=180000
  58. family=reuse;policy=mat;refs=4;consumer=1%
  59. family=bag;stmt=join_on;threshold=20
  60. family=bag;stmt=join_where;threshold=20
  61. family=reuse;policy=mat;refs=2;consumer=100%
  62. family=bag;stmt=cte_joined_outer;threshold=200
  63. family=reuse;policy=notmat;refs=4;consumer=1%
  64. family=bag;stmt=cte_joined_outer_mat;threshold=200
  65. family=bag;stmt=join_on;threshold=180000
  66. family=bag;stmt=cte_joined_outer;threshold=180000
  67. family=reuse;policy=default;refs=4;consumer=100%
  68. family=bag;stmt=cte_joined_outer_notmat;threshold=200
  69. family=reuse;policy=notmat;refs=2;consumer=100%
  70. family=bag;stmt=cte_filtered_mat;threshold=20000
  71. family=reuse;policy=notmat;refs=1;consumer=100%
  72. family=bag;stmt=join_where;threshold=100000
  73. family=reuse;policy=default;refs=16;consumer=100%
  74. family=reuse;policy=notmat;refs=16;consumer=1%
  75. family=bag;stmt=cte_filtered;threshold=2000
  76. family=reuse;policy=default;refs=4;consumer=1%
  77. family=bag;stmt=join_on;threshold=2000
  78. family=reuse;policy=default;refs=2;consumer=1%
  79. family=reuse;policy=notmat;refs=4;consumer=100%
  80. family=reuse;policy=mat;refs=2;consumer=1%
  81. family=reuse;policy=default;refs=16;consumer=1%
  82. family=bag;stmt=comma_join;threshold=200
  83. family=bag;stmt=comma_join;threshold=100000
  84. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  85. family=reuse;policy=notmat;refs=8;consumer=100%
  86. family=bag;stmt=cte_filtered_mat;threshold=20
  87. family=bag;stmt=cte_joined_outer;threshold=20000
  88. family=reuse;policy=mat;refs=4;consumer=100%
  89. family=reuse;policy=mat;refs=8;consumer=100%
  90. family=bag;stmt=join_where;threshold=200
  91. family=bag;stmt=join_where;threshold=20000;position=last

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
| family=bag;stmt=join_where;threshold=20000;position=last | j01_join_where | valid | `09a028d269057a3e` |  |
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
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `596a11d1cf99201c` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `549a07108ab01ba9` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `2ce8b85a7b60c88e` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=2;consumer=1% | r_notmat_n2_c1 | valid | `28cc913d71d0d72b` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=2;consumer=100% | r_notmat_n2_c100 | valid | `fc4d3d1577720fdb` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `9f850e53560baafb` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `9608e3d297c0b2d8` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=8;consumer=1% | r_notmat_n8_c1 | valid | `477f7fbaf1c627cb` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=8;consumer=100% | r_notmat_n8_c100 | valid | `2f8f777dad660f79` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, VERBOSE, BUFFERS, FORMAT TEXT` supported on postgres
- one host: this run used node budget 2 CPU 3g and client budget 1 CPU 2g, so budget comparisons are separate classes here, not independent hardware
- isolated phase: no concurrent writers in this run; the blended read/write phase is a separate labelled run
- physical rank-index arm pg-btree-rank; an engine whose planner ignores the index is a different physical design, not a faster one
- distributed RPC counters are N/A on single-node PostgreSQL; shared buffers and rows removed are parsed per cell where printed
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
