# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 91 cells, 91 valid, engine `postgres`, topology `single`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-v3b-domain-pg-i1-20261005`
- run tag: `run/study-01/select-join-v3b-domain-pg-i1-20261005`
- environment: `host-zenbook-ux5406sa`
- repo commit: `7243fee074bfb16fc5ad8b6afc1a2d7549080935`
- engine version: `PostgreSQL 17.11 (Debian 17.11-1.pgdg13+2) on x86_64-pc-linux-gnu, compiled by gcc (Debian 14.2.0-19) 14.2.0, 64-bit`
- database engine image: `docker.io/library/postgres:17.11` @ `7296f210ae81031ec955dbad9a67a84fe958572a2153b8d0826a647522904dc1`
- benchmark client image: `localhost/charitybench-select:select-join-v3b-domain-pg-i1-20261005` @ `f3f1e10c44cbcb4c7b5bc56ed95a4f071682b6f7d33b5e0a7b688c4271eeef88`
- inputs digest: `724642ee2a9420bed8795540dfb5132831794dbb5f4c0588ace27cf0707614da`
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
| 20 (0.01%) | 0.4848 | 0.4506 | 0.6767 | 0.5432 | 0.5306 | 0.4771 | 0.4948 | 49.0254 | 0.5795 | 0.4313 |
| 200 (0.1%) | 0.7903 | 0.7477 | 1.0697 | 0.8460 | 0.8011 | 1.1406 | 0.7971 | 50.5228 | 1.0369 | 0.8444 |
| 2000 (1%) | 1.8787 | 2.2845 | 2.0170 | 2.1408 | 2.3362 | 1.8521 | 2.0924 | 50.3981 | 2.2592 | 2.1599 |
| 20000 (10%) | 7.6594 | 7.7137 | 7.2144 | 7.4998 | 9.5370 | 7.4806 | 7.6227 | 54.0679 | 7.4865 | 7.2673 |
| 100000 (50%) | 48.5971 | 64.6798 | 47.1748 | 46.3914 | 60.7805 | 66.2708 | 47.8812 | 63.2750 | 67.6797 | 55.6649 |
| 180000 (90%) | 73.1975 | 74.8171 | 52.7558 | 80.1814 | 88.7252 | 77.5869 | 80.1983 | 82.1193 | 77.1722 | 53.9519 |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | join_on vs join_where | 0.4848 | 0.4506 | 0.9295 | true |
| 20 | filtered_derived vs join_where | 0.4848 | 0.6767 | 1.3957 | true |
| 20 | cte_filtered vs join_where | 0.4848 | 0.5432 | 1.1204 | true |
| 20 | cte_filtered_mat vs join_where | 0.4848 | 0.5306 | 1.0945 | false |
| 20 | cte_filtered_notmat vs join_where | 0.4848 | 0.4771 | 0.9841 | true |
| 20 | cte_joined_outer vs join_where | 0.4848 | 0.4948 | 1.0206 | true |
| 20 | cte_joined_outer_mat vs join_where | 0.4848 | 49.0254 | 101.1198 | false |
| 20 | cte_joined_outer_notmat vs join_where | 0.4848 | 0.5795 | 1.1953 | true |
| 20 | comma_join vs join_where | 0.4848 | 0.4313 | 0.8895 | true |
| 200 | join_on vs join_where | 0.7903 | 0.7477 | 0.9460 | true |
| 200 | filtered_derived vs join_where | 0.7903 | 1.0697 | 1.3534 | true |
| 200 | cte_filtered vs join_where | 0.7903 | 0.8460 | 1.0705 | true |
| 200 | cte_filtered_mat vs join_where | 0.7903 | 0.8011 | 1.0136 | false |
| 200 | cte_filtered_notmat vs join_where | 0.7903 | 1.1406 | 1.4431 | true |
| 200 | cte_joined_outer vs join_where | 0.7903 | 0.7971 | 1.0086 | true |
| 200 | cte_joined_outer_mat vs join_where | 0.7903 | 50.5228 | 63.9266 | false |
| 200 | cte_joined_outer_notmat vs join_where | 0.7903 | 1.0369 | 1.3120 | true |
| 200 | comma_join vs join_where | 0.7903 | 0.8444 | 1.0684 | true |
| 2000 | join_on vs join_where | 1.8787 | 2.2845 | 1.2159 | true |
| 2000 | filtered_derived vs join_where | 1.8787 | 2.0170 | 1.0736 | true |
| 2000 | cte_filtered vs join_where | 1.8787 | 2.1408 | 1.1395 | true |
| 2000 | cte_filtered_mat vs join_where | 1.8787 | 2.3362 | 1.2435 | false |
| 2000 | cte_filtered_notmat vs join_where | 1.8787 | 1.8521 | 0.9858 | true |
| 2000 | cte_joined_outer vs join_where | 1.8787 | 2.0924 | 1.1137 | true |
| 2000 | cte_joined_outer_mat vs join_where | 1.8787 | 50.3981 | 26.8254 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 1.8787 | 2.2592 | 1.2025 | true |
| 2000 | comma_join vs join_where | 1.8787 | 2.1599 | 1.1497 | true |
| 20000 | join_on vs join_where | 7.6594 | 7.7137 | 1.0071 | true |
| 20000 | filtered_derived vs join_where | 7.6594 | 7.2144 | 0.9419 | true |
| 20000 | cte_filtered vs join_where | 7.6594 | 7.4998 | 0.9792 | true |
| 20000 | cte_filtered_mat vs join_where | 7.6594 | 9.5370 | 1.2451 | false |
| 20000 | cte_filtered_notmat vs join_where | 7.6594 | 7.4806 | 0.9767 | true |
| 20000 | cte_joined_outer vs join_where | 7.6594 | 7.6227 | 0.9952 | true |
| 20000 | cte_joined_outer_mat vs join_where | 7.6594 | 54.0679 | 7.0590 | false |
| 20000 | cte_joined_outer_notmat vs join_where | 7.6594 | 7.4865 | 0.9774 | true |
| 20000 | comma_join vs join_where | 7.6594 | 7.2673 | 0.9488 | true |
| 100000 | join_on vs join_where | 48.5971 | 64.6798 | 1.3309 | true |
| 100000 | filtered_derived vs join_where | 48.5971 | 47.1748 | 0.9707 | true |
| 100000 | cte_filtered vs join_where | 48.5971 | 46.3914 | 0.9546 | true |
| 100000 | cte_filtered_mat vs join_where | 48.5971 | 60.7805 | 1.2507 | false |
| 100000 | cte_filtered_notmat vs join_where | 48.5971 | 66.2708 | 1.3637 | true |
| 100000 | cte_joined_outer vs join_where | 48.5971 | 47.8812 | 0.9853 | true |
| 100000 | cte_joined_outer_mat vs join_where | 48.5971 | 63.2750 | 1.3020 | false |
| 100000 | cte_joined_outer_notmat vs join_where | 48.5971 | 67.6797 | 1.3927 | true |
| 100000 | comma_join vs join_where | 48.5971 | 55.6649 | 1.1454 | true |
| 180000 | join_on vs join_where | 73.1975 | 74.8171 | 1.0221 | true |
| 180000 | filtered_derived vs join_where | 73.1975 | 52.7558 | 0.7207 | true |
| 180000 | cte_filtered vs join_where | 73.1975 | 80.1814 | 1.0954 | true |
| 180000 | cte_filtered_mat vs join_where | 73.1975 | 88.7252 | 1.2121 | false |
| 180000 | cte_filtered_notmat vs join_where | 73.1975 | 77.5869 | 1.0600 | true |
| 180000 | cte_joined_outer vs join_where | 73.1975 | 80.1983 | 1.0956 | true |
| 180000 | cte_joined_outer_mat vs join_where | 73.1975 | 82.1193 | 1.1219 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 73.1975 | 77.1722 | 1.0543 | true |
| 180000 | comma_join vs join_where | 73.1975 | 53.9519 | 0.7371 | true |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_filtered_mat vs cte_filtered | 0.5432 | 0.5306 | 0.9769 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_filtered_notmat vs cte_filtered | 0.5432 | 0.4771 | 0.8783 | true | policy=not_materialized; observed_cte_scan=false |
| 20 | cte_joined_outer_mat vs cte_joined_outer | 0.4948 | 49.0254 | 99.0785 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 0.4948 | 0.5795 | 1.1712 | true | policy=not_materialized; observed_cte_scan=false |
| 200 | cte_filtered_mat vs cte_filtered | 0.8460 | 0.8011 | 0.9469 | false | policy=materialized; observed_cte_scan=true |
| 200 | cte_filtered_notmat vs cte_filtered | 0.8460 | 1.1406 | 1.3482 | true | policy=not_materialized; observed_cte_scan=false |
| 200 | cte_joined_outer_mat vs cte_joined_outer | 0.7971 | 50.5228 | 63.3842 | false | policy=materialized; observed_cte_scan=true |
| 200 | cte_joined_outer_notmat vs cte_joined_outer | 0.7971 | 1.0369 | 1.3009 | true | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_filtered_mat vs cte_filtered | 2.1408 | 2.3362 | 1.0913 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_filtered_notmat vs cte_filtered | 2.1408 | 1.8521 | 0.8651 | true | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 2.0924 | 50.3981 | 24.0869 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 2.0924 | 2.2592 | 1.0797 | true | policy=not_materialized; observed_cte_scan=false |
| 20000 | cte_filtered_mat vs cte_filtered | 7.4998 | 9.5370 | 1.2716 | false | policy=materialized; observed_cte_scan=true |
| 20000 | cte_filtered_notmat vs cte_filtered | 7.4998 | 7.4806 | 0.9974 | true | policy=not_materialized; observed_cte_scan=false |
| 20000 | cte_joined_outer_mat vs cte_joined_outer | 7.6227 | 54.0679 | 7.0930 | false | policy=materialized; observed_cte_scan=true |
| 20000 | cte_joined_outer_notmat vs cte_joined_outer | 7.6227 | 7.4865 | 0.9821 | true | policy=not_materialized; observed_cte_scan=false |
| 100000 | cte_filtered_mat vs cte_filtered | 46.3914 | 60.7805 | 1.3102 | false | policy=materialized; observed_cte_scan=true |
| 100000 | cte_filtered_notmat vs cte_filtered | 46.3914 | 66.2708 | 1.4285 | true | policy=not_materialized; observed_cte_scan=false |
| 100000 | cte_joined_outer_mat vs cte_joined_outer | 47.8812 | 63.2750 | 1.3215 | false | policy=materialized; observed_cte_scan=true |
| 100000 | cte_joined_outer_notmat vs cte_joined_outer | 47.8812 | 67.6797 | 1.4135 | true | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_filtered_mat vs cte_filtered | 80.1814 | 88.7252 | 1.1066 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_filtered_notmat vs cte_filtered | 80.1814 | 77.5869 | 0.9676 | true | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 80.1983 | 82.1193 | 1.0240 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 80.1983 | 77.1722 | 0.9623 | true | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 0.6139 | 1 | `549a07108ab01ba9` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 2.6169 | 1 | `2ce8b85a7b60c88e` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | 2.9484 | 2 | `b307c41eb54cf452` | policy=default; observed_cte_scan=true | valid |
| default | 2 | 100% | 2.7170 | 2 | `b307c41eb54cf452` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 1% | 3.1623 | 4 | `ccfe1af393deae98` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 3.1052 | 4 | `ccfe1af393deae98` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | 3.1344 | 8 | `c8c73778097931dc` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 100% | 3.2971 | 8 | `c8c73778097931dc` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 1% | 3.8264 | 16 | `b0c356323d30c425` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 3.8993 | 16 | `b0c356323d30c425` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 2.8372 | 1 | `eaa7d6058aaa5e03` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 2.6763 | 1 | `eaa7d6058aaa5e03` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | 3.0586 | 2 | `b307c41eb54cf452` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 100% | 2.7738 | 2 | `b307c41eb54cf452` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 1% | 2.7508 | 4 | `ccfe1af393deae98` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 2.9854 | 4 | `ccfe1af393deae98` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | 3.4181 | 8 | `c8c73778097931dc` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 100% | 3.2834 | 8 | `c8c73778097931dc` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 1% | 3.7634 | 16 | `b0c356323d30c425` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 4.2304 | 16 | `b0c356323d30c425` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 0.4855 | 1 | `549a07108ab01ba9` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 2.8192 | 1 | `2ce8b85a7b60c88e` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | 0.6380 | 2 | `28cc913d71d0d72b` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 100% | 5.3424 | 2 | `0d8b364aa09de4b1` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 1% | 0.9916 | 4 | `9f850e53560baafb` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 8.1854 | 4 | `e61aa6f65a629698` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | 1.7151 | 8 | `477f7fbaf1c627cb` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 100% | 10.7261 | 8 | `17440da042fdf60e` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 1% | 2.6808 | 16 | `1c38b58d8cec08d3` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 16.7457 | 16 | `cd83742a4c56fb2b` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 4.6216 | false |
| notmat | 1 | 1% | 0.7908 | true |
| mat | 1 | 100% | 1.0227 | false |
| notmat | 1 | 100% | 1.0773 | true |
| mat | 2 | 1% | 1.0374 | true |
| notmat | 2 | 1% | 0.2164 | false |
| mat | 2 | 100% | 1.0209 | true |
| notmat | 2 | 100% | 1.9663 | false |
| mat | 4 | 1% | 0.8699 | true |
| notmat | 4 | 1% | 0.3136 | false |
| mat | 4 | 100% | 0.9614 | true |
| notmat | 4 | 100% | 2.6361 | false |
| mat | 8 | 1% | 1.0905 | true |
| notmat | 8 | 1% | 0.5472 | false |
| mat | 8 | 100% | 0.9959 | true |
| notmat | 8 | 100% | 3.2532 | false |
| mat | 16 | 1% | 0.9835 | true |
| notmat | 16 | 1% | 0.7006 | false |
| mat | 16 | 100% | 1.0849 | true |
| notmat | 16 | 100% | 4.2945 | false |

## Repeated-design position control

| first median ms | last median ms | ratio (last/first) | same plan |
|---|---|---|---|
| 77.586930 | 77.960878 | 1.0048 | true |

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

- client: nr_periods 399, nr_throttled 314, throttled_usec 15550957 | quota 100000 100000
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
| family=bag;stmt=cte_filtered_notmat;threshold=180000;position=last | valid | false | true | 0 | 0 | 1286/0 | 20000 | 0 |
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
| family=bag;stmt=comma_join;threshold=100000 | valid | 43.6452 | 55.6649 | 63.5883 | 5 | 18.0 |
| family=bag;stmt=comma_join;threshold=180000 | valid | 48.3256 | 53.9519 | 58.0991 | 5 | 18.5 |
| family=bag;stmt=comma_join;threshold=20 | valid | 0.3569 | 0.4313 | 0.5225 | 5 | 2318.8 |
| family=bag;stmt=comma_join;threshold=200 | valid | 0.8113 | 0.8444 | 0.9813 | 5 | 1184.3 |
| family=bag;stmt=comma_join;threshold=2000 | valid | 1.9982 | 2.1599 | 2.8161 | 5 | 463.0 |
| family=bag;stmt=comma_join;threshold=20000 | valid | 7.1485 | 7.2673 | 7.9919 | 5 | 137.6 |
| family=bag;stmt=cte_filtered;threshold=100000 | valid | 43.7189 | 46.3914 | 54.5187 | 5 | 21.6 |
| family=bag;stmt=cte_filtered;threshold=180000 | valid | 47.8499 | 80.1814 | 92.9936 | 5 | 12.5 |
| family=bag;stmt=cte_filtered;threshold=20 | valid | 0.4549 | 0.5432 | 26.1384 | 5 | 1840.9 |
| family=bag;stmt=cte_filtered;threshold=200 | valid | 0.7996 | 0.8460 | 0.8649 | 5 | 1182.0 |
| family=bag;stmt=cte_filtered;threshold=2000 | valid | 1.9243 | 2.1408 | 2.7024 | 5 | 467.1 |
| family=bag;stmt=cte_filtered;threshold=20000 | valid | 7.1613 | 7.4998 | 21.2245 | 5 | 133.3 |
| family=bag;stmt=cte_filtered_mat;threshold=100000 | valid | 58.0245 | 60.7805 | 66.3816 | 5 | 16.5 |
| family=bag;stmt=cte_filtered_mat;threshold=180000 | valid | 78.9770 | 88.7252 | 92.1488 | 5 | 11.3 |
| family=bag;stmt=cte_filtered_mat;threshold=20 | valid | 0.4848 | 0.5306 | 0.7673 | 5 | 1884.5 |
| family=bag;stmt=cte_filtered_mat;threshold=200 | valid | 0.7145 | 0.8011 | 0.8384 | 5 | 1248.4 |
| family=bag;stmt=cte_filtered_mat;threshold=2000 | valid | 2.2331 | 2.3362 | 2.6831 | 5 | 428.0 |
| family=bag;stmt=cte_filtered_mat;threshold=20000 | valid | 9.2716 | 9.5370 | 43.4828 | 5 | 104.9 |
| family=bag;stmt=cte_filtered_notmat;threshold=100000 | valid | 46.1741 | 66.2708 | 72.3955 | 5 | 15.1 |
| family=bag;stmt=cte_filtered_notmat;threshold=180000 | valid | 53.5406 | 77.5869 | 87.1533 | 5 | 12.9 |
| family=bag;stmt=cte_filtered_notmat;threshold=180000;position=last | valid | 47.4997 | 77.9609 | 89.3097 | 5 | 12.8 |
| family=bag;stmt=cte_filtered_notmat;threshold=20 | valid | 0.4371 | 0.4771 | 0.5984 | 5 | 2095.9 |
| family=bag;stmt=cte_filtered_notmat;threshold=200 | valid | 0.9849 | 1.1406 | 1.2352 | 5 | 876.8 |
| family=bag;stmt=cte_filtered_notmat;threshold=2000 | valid | 1.6994 | 1.8521 | 1.9385 | 5 | 539.9 |
| family=bag;stmt=cte_filtered_notmat;threshold=20000 | valid | 7.4071 | 7.4806 | 8.1597 | 5 | 133.7 |
| family=bag;stmt=cte_joined_outer;threshold=100000 | valid | 46.3568 | 47.8812 | 54.9268 | 5 | 20.9 |
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 45.9540 | 80.1983 | 88.2784 | 5 | 12.5 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 0.4545 | 0.4948 | 0.5305 | 5 | 2021.0 |
| family=bag;stmt=cte_joined_outer;threshold=200 | valid | 0.7603 | 0.7971 | 0.9480 | 5 | 1254.6 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 2.0637 | 2.0924 | 2.1431 | 5 | 477.9 |
| family=bag;stmt=cte_joined_outer;threshold=20000 | valid | 7.2845 | 7.6227 | 7.9738 | 5 | 131.2 |
| family=bag;stmt=cte_joined_outer_mat;threshold=100000 | valid | 61.6687 | 63.2750 | 66.7146 | 5 | 15.8 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 70.5309 | 82.1193 | 87.4275 | 5 | 12.2 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 48.0131 | 49.0254 | 52.0646 | 5 | 20.4 |
| family=bag;stmt=cte_joined_outer_mat;threshold=200 | valid | 48.6973 | 50.5228 | 52.3115 | 5 | 19.8 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 48.4136 | 50.3981 | 51.5646 | 5 | 19.8 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20000 | valid | 51.5673 | 54.0679 | 58.4968 | 5 | 18.5 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=100000 | valid | 41.9229 | 67.6797 | 74.0003 | 5 | 14.8 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 67.5905 | 77.1722 | 81.5498 | 5 | 13.0 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 0.5059 | 0.5795 | 0.6517 | 5 | 1725.5 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=200 | valid | 0.9965 | 1.0369 | 1.2179 | 5 | 964.4 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 2.1852 | 2.2592 | 2.4015 | 5 | 442.6 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20000 | valid | 7.3274 | 7.4865 | 34.1180 | 5 | 133.6 |
| family=bag;stmt=filtered_derived;threshold=100000 | valid | 43.9080 | 47.1748 | 51.0970 | 5 | 21.2 |
| family=bag;stmt=filtered_derived;threshold=180000 | valid | 49.4634 | 52.7558 | 57.6167 | 5 | 19.0 |
| family=bag;stmt=filtered_derived;threshold=20 | valid | 0.5659 | 0.6767 | 0.7407 | 5 | 1477.8 |
| family=bag;stmt=filtered_derived;threshold=200 | valid | 0.9696 | 1.0697 | 1.1166 | 5 | 934.9 |
| family=bag;stmt=filtered_derived;threshold=2000 | valid | 1.9917 | 2.0170 | 2.2743 | 5 | 495.8 |
| family=bag;stmt=filtered_derived;threshold=20000 | valid | 6.9875 | 7.2144 | 8.5805 | 5 | 138.6 |
| family=bag;stmt=join_on;threshold=100000 | valid | 51.1335 | 64.6798 | 66.2111 | 5 | 15.5 |
| family=bag;stmt=join_on;threshold=180000 | valid | 48.4759 | 74.8171 | 84.1075 | 5 | 13.4 |
| family=bag;stmt=join_on;threshold=20 | valid | 0.3959 | 0.4506 | 0.4727 | 5 | 2219.1 |
| family=bag;stmt=join_on;threshold=200 | valid | 0.7065 | 0.7477 | 0.9401 | 5 | 1337.5 |
| family=bag;stmt=join_on;threshold=2000 | valid | 2.1175 | 2.2845 | 2.4310 | 5 | 437.7 |
| family=bag;stmt=join_on;threshold=20000 | valid | 7.2325 | 7.7137 | 8.0448 | 5 | 129.6 |
| family=bag;stmt=join_where;threshold=100000 | valid | 43.9880 | 48.5971 | 66.9333 | 5 | 20.6 |
| family=bag;stmt=join_where;threshold=180000 | valid | 67.8013 | 73.1975 | 78.6240 | 5 | 13.7 |
| family=bag;stmt=join_where;threshold=20 | valid | 0.4373 | 0.4848 | 0.6090 | 5 | 2062.6 |
| family=bag;stmt=join_where;threshold=200 | valid | 0.7183 | 0.7903 | 1.3484 | 5 | 1265.3 |
| family=bag;stmt=join_where;threshold=2000 | valid | 1.6939 | 1.8787 | 3.0376 | 5 | 532.3 |
| family=bag;stmt=join_where;threshold=20000 | valid | 7.5240 | 7.6594 | 45.5486 | 5 | 130.6 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 3.5556 | 3.8264 | 4.1796 | 5 | 261.3 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 3.8146 | 3.8993 | 4.1452 | 5 | 256.5 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 0.5279 | 0.6139 | 0.6256 | 5 | 1628.9 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 2.5088 | 2.6169 | 2.9763 | 5 | 382.1 |
| family=reuse;policy=default;refs=2;consumer=1% | valid | 2.8835 | 2.9484 | 3.2694 | 5 | 339.2 |
| family=reuse;policy=default;refs=2;consumer=100% | valid | 2.6519 | 2.7170 | 3.0467 | 5 | 368.0 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 3.0366 | 3.1623 | 3.4303 | 5 | 316.2 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 2.9904 | 3.1052 | 3.2825 | 5 | 322.0 |
| family=reuse;policy=default;refs=8;consumer=1% | valid | 2.9805 | 3.1344 | 3.8006 | 5 | 319.0 |
| family=reuse;policy=default;refs=8;consumer=100% | valid | 3.2016 | 3.2971 | 3.4190 | 5 | 303.3 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 3.5781 | 3.7634 | 3.9811 | 5 | 265.7 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 4.1643 | 4.2304 | 4.4679 | 5 | 236.4 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 2.7915 | 2.8372 | 3.0578 | 5 | 352.5 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 2.5309 | 2.6763 | 3.0358 | 5 | 373.6 |
| family=reuse;policy=mat;refs=2;consumer=1% | valid | 2.9218 | 3.0586 | 3.3131 | 5 | 327.0 |
| family=reuse;policy=mat;refs=2;consumer=100% | valid | 2.6243 | 2.7738 | 3.2672 | 5 | 360.5 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 2.6433 | 2.7508 | 2.8443 | 5 | 363.5 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 2.8269 | 2.9854 | 3.4044 | 5 | 335.0 |
| family=reuse;policy=mat;refs=8;consumer=1% | valid | 3.2741 | 3.4181 | 3.7266 | 5 | 292.6 |
| family=reuse;policy=mat;refs=8;consumer=100% | valid | 3.2124 | 3.2834 | 3.5222 | 5 | 304.6 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 2.5319 | 2.6808 | 2.7647 | 5 | 373.0 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 15.9592 | 16.7457 | 17.1056 | 5 | 59.7 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 0.4362 | 0.4855 | 0.5676 | 5 | 2059.8 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 2.7395 | 2.8192 | 3.0825 | 5 | 354.7 |
| family=reuse;policy=notmat;refs=2;consumer=1% | valid | 0.5808 | 0.6380 | 0.6478 | 5 | 1567.5 |
| family=reuse;policy=notmat;refs=2;consumer=100% | valid | 5.1432 | 5.3424 | 5.4726 | 5 | 187.2 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 0.8994 | 0.9916 | 1.0085 | 5 | 1008.5 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 7.9673 | 8.1854 | 8.2112 | 5 | 122.2 |
| family=reuse;policy=notmat;refs=8;consumer=1% | valid | 1.6042 | 1.7151 | 1.7457 | 5 | 583.1 |
| family=reuse;policy=notmat;refs=8;consumer=100% | valid | 10.3631 | 10.7261 | 11.2399 | 5 | 93.2 |

## Matched blended workload

- phase: `isolated`; target offered rate: 10 updates/s
- offered 0, committed 0, failed 0 over 0.000 s
- offered rate 0.000/s, completed rate 0.000/s
- per key: no keys touched
- ledger: 0 keys reconciled, 0 mismatch

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42702115`, 5 interleaved rounds
  1. family=bag;stmt=cte_filtered_notmat;threshold=180000
  2. family=bag;stmt=cte_filtered_mat;threshold=180000
  3. family=bag;stmt=comma_join;threshold=20000
  4. family=reuse;policy=default;refs=2;consumer=1%
  5. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  6. family=bag;stmt=cte_joined_outer_notmat;threshold=100000
  7. family=bag;stmt=join_on;threshold=2000
  8. family=bag;stmt=join_on;threshold=180000
  9. family=bag;stmt=join_where;threshold=20000
  10. family=reuse;policy=mat;refs=16;consumer=100%
  11. family=bag;stmt=cte_filtered;threshold=20
  12. family=bag;stmt=cte_joined_outer_mat;threshold=100000
  13. family=bag;stmt=cte_joined_outer;threshold=180000
  14. family=bag;stmt=join_where;threshold=100000
  15. family=reuse;policy=notmat;refs=1;consumer=100%
  16. family=bag;stmt=cte_joined_outer;threshold=20
  17. family=reuse;policy=mat;refs=1;consumer=100%
  18. family=reuse;policy=notmat;refs=8;consumer=100%
  19. family=bag;stmt=cte_filtered;threshold=2000
  20. family=bag;stmt=filtered_derived;threshold=20000
  21. family=reuse;policy=default;refs=16;consumer=1%
  22. family=bag;stmt=cte_filtered_mat;threshold=200
  23. family=bag;stmt=filtered_derived;threshold=100000
  24. family=bag;stmt=cte_filtered_notmat;threshold=200
  25. family=bag;stmt=cte_filtered;threshold=20000
  26. family=bag;stmt=join_on;threshold=20000
  27. family=reuse;policy=mat;refs=16;consumer=1%
  28. family=bag;stmt=cte_joined_outer;threshold=200
  29. family=reuse;policy=default;refs=16;consumer=100%
  30. family=bag;stmt=join_where;threshold=200
  31. family=reuse;policy=mat;refs=4;consumer=100%
  32. family=bag;stmt=cte_joined_outer_mat;threshold=20
  33. family=bag;stmt=comma_join;threshold=2000
  34. family=bag;stmt=join_on;threshold=200
  35. family=bag;stmt=join_where;threshold=2000
  36. family=bag;stmt=cte_joined_outer;threshold=20000
  37. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  38. family=bag;stmt=cte_joined_outer;threshold=2000
  39. family=bag;stmt=cte_filtered_mat;threshold=100000
  40. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  41. family=bag;stmt=cte_joined_outer_mat;threshold=20000
  42. family=reuse;policy=notmat;refs=8;consumer=1%
  43. family=reuse;policy=notmat;refs=4;consumer=1%
  44. family=bag;stmt=cte_filtered_mat;threshold=2000
  45. family=reuse;policy=mat;refs=2;consumer=100%
  46. family=bag;stmt=join_on;threshold=20
  47. family=reuse;policy=mat;refs=8;consumer=100%
  48. family=bag;stmt=comma_join;threshold=20
  49. family=bag;stmt=cte_joined_outer_mat;threshold=200
  50. family=reuse;policy=default;refs=1;consumer=1%
  51. family=bag;stmt=cte_filtered_mat;threshold=20
  52. family=bag;stmt=comma_join;threshold=200
  53. family=reuse;policy=notmat;refs=1;consumer=1%
  54. family=bag;stmt=filtered_derived;threshold=2000
  55. family=bag;stmt=filtered_derived;threshold=180000
  56. family=reuse;policy=notmat;refs=2;consumer=100%
  57. family=bag;stmt=cte_filtered_notmat;threshold=20000
  58. family=bag;stmt=filtered_derived;threshold=20
  59. family=bag;stmt=cte_joined_outer;threshold=100000
  60. family=bag;stmt=filtered_derived;threshold=200
  61. family=bag;stmt=comma_join;threshold=100000
  62. family=reuse;policy=mat;refs=8;consumer=1%
  63. family=bag;stmt=join_where;threshold=20
  64. family=reuse;policy=default;refs=1;consumer=100%
  65. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  66. family=bag;stmt=cte_joined_outer_notmat;threshold=200
  67. family=bag;stmt=cte_filtered_notmat;threshold=20
  68. family=bag;stmt=join_on;threshold=100000
  69. family=bag;stmt=cte_filtered;threshold=180000
  70. family=reuse;policy=notmat;refs=16;consumer=100%
  71. family=reuse;policy=default;refs=8;consumer=1%
  72. family=reuse;policy=default;refs=2;consumer=100%
  73. family=reuse;policy=notmat;refs=2;consumer=1%
  74. family=bag;stmt=comma_join;threshold=180000
  75. family=reuse;policy=mat;refs=2;consumer=1%
  76. family=bag;stmt=cte_filtered;threshold=200
  77. family=reuse;policy=notmat;refs=16;consumer=1%
  78. family=reuse;policy=default;refs=8;consumer=100%
  79. family=bag;stmt=cte_filtered_notmat;threshold=2000
  80. family=reuse;policy=mat;refs=4;consumer=1%
  81. family=reuse;policy=notmat;refs=4;consumer=100%
  82. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  83. family=reuse;policy=default;refs=4;consumer=100%
  84. family=bag;stmt=cte_filtered;threshold=100000
  85. family=bag;stmt=cte_filtered_notmat;threshold=100000
  86. family=bag;stmt=join_where;threshold=180000
  87. family=reuse;policy=default;refs=4;consumer=1%
  88. family=bag;stmt=cte_filtered_mat;threshold=20000
  89. family=bag;stmt=cte_joined_outer_notmat;threshold=20000
  90. family=reuse;policy=mat;refs=1;consumer=1%
  91. family=bag;stmt=cte_filtered_notmat;threshold=180000;position=last

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
| family=bag;stmt=cte_filtered_notmat;threshold=180000;position=last | j04n_cte_filtered_notmat | valid | `18e11c5d2e685630` | policy=not_materialized; observed_cte_scan=false |
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
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `cd83742a4c56fb2b` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `549a07108ab01ba9` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `2ce8b85a7b60c88e` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=2;consumer=1% | r_notmat_n2_c1 | valid | `28cc913d71d0d72b` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=2;consumer=100% | r_notmat_n2_c100 | valid | `0d8b364aa09de4b1` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `9f850e53560baafb` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `e61aa6f65a629698` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=8;consumer=1% | r_notmat_n8_c1 | valid | `477f7fbaf1c627cb` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=8;consumer=100% | r_notmat_n8_c100 | valid | `17440da042fdf60e` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, VERBOSE, BUFFERS, FORMAT TEXT` supported on postgres
- one host: this run used node budget 2 CPU 3g and client budget 1 CPU 2g, so budget comparisons are separate classes here, not independent hardware
- isolated phase: no concurrent writers in this run; the blended read/write phase is a separate labelled run
- physical rank-index arm pg-btree-rank; an engine whose planner ignores the index is a different physical design, not a faster one
- distributed RPC counters are N/A on single-node PostgreSQL; shared buffers and rows removed are parsed per cell where printed
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
