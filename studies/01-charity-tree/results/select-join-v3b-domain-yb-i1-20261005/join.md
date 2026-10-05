# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 91 cells, 91 valid, engine `yugabyte`, topology `single`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-v3b-domain-yb-i1-20261005`
- run tag: `run/study-01/select-join-v3b-domain-yb-i1-20261005`
- environment: `host-zenbook-ux5406sa`
- repo commit: `7243fee074bfb16fc5ad8b6afc1a2d7549080935`
- engine version: `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit`
- database engine image: `docker.io/yugabytedb/yugabyte:2025.2.6.0-b111` @ `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`
- benchmark client image: `localhost/charitybench-select:select-join-v3b-domain-yb-i1-20261005` @ `f3f1e10c44cbcb4c7b5bc56ed95a4f071682b6f7d33b5e0a7b688c4271eeef88`
- inputs digest: `aed6b10e8425f3871462de7a93c38b198680edb2c1ef539e7a94d85c90c8ffa6`
- dataset fingerprint: `fcf26e08f377ec51`
- client CPU quota: `100000 100000`, memory: `2147483648`
- setting `client_budget`: 1 CPU 2g
- setting `client_workers`: 1 sequential reader; 0 writers
- setting `connection_pool`: 2 per endpoint = reader/setup + reserve
- setting `customers`: 2000
- setting `dataset`: study01-person-donation-charity
- setting `distributed_rpc`: DIST/DEBUG storage counters parsed per cell (Storage Rows Scanned, Storage Index Rows Scanned); RPC-level counters not aggregated
- setting `endpoint_mode`: single endpoint
- setting `endpoints`: yb-single
- setting `fact_count`: 200000
- setting `index_arm`: yb-range-rank
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
| yb-single | 933 | 1 |

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
| 20 (0.01%) | 2.1515 | 1.9798 | 1.9032 | 1.7762 | 1.8862 | 2.4434 | 2.0400 | 88.4830 | 1.9852 | 1.9191 |
| 200 (0.1%) | 3.1810 | 3.1592 | 3.6986 | 3.6180 | 3.6131 | 3.2571 | 4.2330 | 86.6851 | 3.7862 | 3.2192 |
| 2000 (1%) | 9.6082 | 8.4765 | 9.3428 | 9.0673 | 10.3554 | 9.5377 | 9.6825 | 85.8304 | 10.0079 | 10.4146 |
| 20000 (10%) | 66.7985 | 62.5451 | 67.4144 | 68.0607 | 67.1924 | 68.1455 | 65.8666 | 90.7192 | 64.8759 | 67.2274 |
| 100000 (50%) | 86.5406 | 89.8909 | 89.8128 | 95.5874 | 96.8301 | 91.5077 | 91.6823 | 100.7127 | 91.1699 | 91.7894 |
| 180000 (90%) | 108.4982 | 107.3997 | 118.3836 | 108.9847 | 121.2151 | 111.1850 | 111.0178 | 113.7440 | 115.3694 | 129.6052 |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | join_on vs join_where | 2.1515 | 1.9798 | 0.9202 | false |
| 20 | filtered_derived vs join_where | 2.1515 | 1.9032 | 0.8846 | false |
| 20 | cte_filtered vs join_where | 2.1515 | 1.7762 | 0.8256 | false |
| 20 | cte_filtered_mat vs join_where | 2.1515 | 1.8862 | 0.8767 | false |
| 20 | cte_filtered_notmat vs join_where | 2.1515 | 2.4434 | 1.1357 | false |
| 20 | cte_joined_outer vs join_where | 2.1515 | 2.0400 | 0.9482 | false |
| 20 | cte_joined_outer_mat vs join_where | 2.1515 | 88.4830 | 41.1269 | false |
| 20 | cte_joined_outer_notmat vs join_where | 2.1515 | 1.9852 | 0.9227 | false |
| 20 | comma_join vs join_where | 2.1515 | 1.9191 | 0.8920 | false |
| 200 | join_on vs join_where | 3.1810 | 3.1592 | 0.9931 | false |
| 200 | filtered_derived vs join_where | 3.1810 | 3.6986 | 1.1627 | false |
| 200 | cte_filtered vs join_where | 3.1810 | 3.6180 | 1.1374 | false |
| 200 | cte_filtered_mat vs join_where | 3.1810 | 3.6131 | 1.1359 | false |
| 200 | cte_filtered_notmat vs join_where | 3.1810 | 3.2571 | 1.0239 | false |
| 200 | cte_joined_outer vs join_where | 3.1810 | 4.2330 | 1.3307 | false |
| 200 | cte_joined_outer_mat vs join_where | 3.1810 | 86.6851 | 27.2512 | false |
| 200 | cte_joined_outer_notmat vs join_where | 3.1810 | 3.7862 | 1.1903 | false |
| 200 | comma_join vs join_where | 3.1810 | 3.2192 | 1.0120 | false |
| 2000 | join_on vs join_where | 9.6082 | 8.4765 | 0.8822 | false |
| 2000 | filtered_derived vs join_where | 9.6082 | 9.3428 | 0.9724 | false |
| 2000 | cte_filtered vs join_where | 9.6082 | 9.0673 | 0.9437 | false |
| 2000 | cte_filtered_mat vs join_where | 9.6082 | 10.3554 | 1.0778 | false |
| 2000 | cte_filtered_notmat vs join_where | 9.6082 | 9.5377 | 0.9927 | false |
| 2000 | cte_joined_outer vs join_where | 9.6082 | 9.6825 | 1.0077 | false |
| 2000 | cte_joined_outer_mat vs join_where | 9.6082 | 85.8304 | 8.9331 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 9.6082 | 10.0079 | 1.0416 | false |
| 2000 | comma_join vs join_where | 9.6082 | 10.4146 | 1.0839 | false |
| 20000 | join_on vs join_where | 66.7985 | 62.5451 | 0.9363 | false |
| 20000 | filtered_derived vs join_where | 66.7985 | 67.4144 | 1.0092 | false |
| 20000 | cte_filtered vs join_where | 66.7985 | 68.0607 | 1.0189 | false |
| 20000 | cte_filtered_mat vs join_where | 66.7985 | 67.1924 | 1.0059 | false |
| 20000 | cte_filtered_notmat vs join_where | 66.7985 | 68.1455 | 1.0202 | false |
| 20000 | cte_joined_outer vs join_where | 66.7985 | 65.8666 | 0.9860 | false |
| 20000 | cte_joined_outer_mat vs join_where | 66.7985 | 90.7192 | 1.3581 | false |
| 20000 | cte_joined_outer_notmat vs join_where | 66.7985 | 64.8759 | 0.9712 | false |
| 20000 | comma_join vs join_where | 66.7985 | 67.2274 | 1.0064 | false |
| 100000 | join_on vs join_where | 86.5406 | 89.8909 | 1.0387 | false |
| 100000 | filtered_derived vs join_where | 86.5406 | 89.8128 | 1.0378 | false |
| 100000 | cte_filtered vs join_where | 86.5406 | 95.5874 | 1.1045 | false |
| 100000 | cte_filtered_mat vs join_where | 86.5406 | 96.8301 | 1.1189 | false |
| 100000 | cte_filtered_notmat vs join_where | 86.5406 | 91.5077 | 1.0574 | false |
| 100000 | cte_joined_outer vs join_where | 86.5406 | 91.6823 | 1.0594 | false |
| 100000 | cte_joined_outer_mat vs join_where | 86.5406 | 100.7127 | 1.1638 | false |
| 100000 | cte_joined_outer_notmat vs join_where | 86.5406 | 91.1699 | 1.0535 | false |
| 100000 | comma_join vs join_where | 86.5406 | 91.7894 | 1.0607 | false |
| 180000 | join_on vs join_where | 108.4982 | 107.3997 | 0.9899 | false |
| 180000 | filtered_derived vs join_where | 108.4982 | 118.3836 | 1.0911 | false |
| 180000 | cte_filtered vs join_where | 108.4982 | 108.9847 | 1.0045 | false |
| 180000 | cte_filtered_mat vs join_where | 108.4982 | 121.2151 | 1.1172 | false |
| 180000 | cte_filtered_notmat vs join_where | 108.4982 | 111.1850 | 1.0248 | false |
| 180000 | cte_joined_outer vs join_where | 108.4982 | 111.0178 | 1.0232 | false |
| 180000 | cte_joined_outer_mat vs join_where | 108.4982 | 113.7440 | 1.0483 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 108.4982 | 115.3694 | 1.0633 | false |
| 180000 | comma_join vs join_where | 108.4982 | 129.6052 | 1.1945 | false |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_filtered_mat vs cte_filtered | 1.7762 | 1.8862 | 1.0620 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_filtered_notmat vs cte_filtered | 1.7762 | 2.4434 | 1.3757 | false | policy=not_materialized; observed_cte_scan=false |
| 20 | cte_joined_outer_mat vs cte_joined_outer | 2.0400 | 88.4830 | 43.3747 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 2.0400 | 1.9852 | 0.9731 | false | policy=not_materialized; observed_cte_scan=false |
| 200 | cte_filtered_mat vs cte_filtered | 3.6180 | 3.6131 | 0.9987 | false | policy=materialized; observed_cte_scan=true |
| 200 | cte_filtered_notmat vs cte_filtered | 3.6180 | 3.2571 | 0.9003 | false | policy=not_materialized; observed_cte_scan=false |
| 200 | cte_joined_outer_mat vs cte_joined_outer | 4.2330 | 86.6851 | 20.4786 | false | policy=materialized; observed_cte_scan=true |
| 200 | cte_joined_outer_notmat vs cte_joined_outer | 4.2330 | 3.7862 | 0.8945 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_filtered_mat vs cte_filtered | 9.0673 | 10.3554 | 1.1421 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_filtered_notmat vs cte_filtered | 9.0673 | 9.5377 | 1.0519 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 9.6825 | 85.8304 | 8.8645 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 9.6825 | 10.0079 | 1.0336 | false | policy=not_materialized; observed_cte_scan=false |
| 20000 | cte_filtered_mat vs cte_filtered | 68.0607 | 67.1924 | 0.9872 | false | policy=materialized; observed_cte_scan=true |
| 20000 | cte_filtered_notmat vs cte_filtered | 68.0607 | 68.1455 | 1.0012 | false | policy=not_materialized; observed_cte_scan=false |
| 20000 | cte_joined_outer_mat vs cte_joined_outer | 65.8666 | 90.7192 | 1.3773 | false | policy=materialized; observed_cte_scan=true |
| 20000 | cte_joined_outer_notmat vs cte_joined_outer | 65.8666 | 64.8759 | 0.9850 | false | policy=not_materialized; observed_cte_scan=false |
| 100000 | cte_filtered_mat vs cte_filtered | 95.5874 | 96.8301 | 1.0130 | false | policy=materialized; observed_cte_scan=true |
| 100000 | cte_filtered_notmat vs cte_filtered | 95.5874 | 91.5077 | 0.9573 | false | policy=not_materialized; observed_cte_scan=false |
| 100000 | cte_joined_outer_mat vs cte_joined_outer | 91.6823 | 100.7127 | 1.0985 | false | policy=materialized; observed_cte_scan=true |
| 100000 | cte_joined_outer_notmat vs cte_joined_outer | 91.6823 | 91.1699 | 0.9944 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_filtered_mat vs cte_filtered | 108.9847 | 121.2151 | 1.1122 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_filtered_notmat vs cte_filtered | 108.9847 | 111.1850 | 1.0202 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 111.0178 | 113.7440 | 1.0246 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 111.0178 | 115.3694 | 1.0392 | false | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 7.2801 | 1 | `1e464b3f33f72b04` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 8.7633 | 1 | `fbac5a7a22297229` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | 8.1555 | 2 | `bc0b0617e1357ec5` | policy=default; observed_cte_scan=true | valid |
| default | 2 | 100% | 9.2403 | 2 | `4e205dbd55a42666` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 1% | 8.0565 | 4 | `3865366089004bbd` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 9.3743 | 4 | `2bbfe6ff213a1d88` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | 8.7641 | 8 | `62fbb2d6d3c66d6c` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 100% | 9.0791 | 8 | `18d8ad5a2031c6cf` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 1% | 9.3403 | 16 | `f69657e41cd7f916` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 10.1317 | 16 | `349c18ad78591646` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 8.5311 | 1 | `0bfec48ca42bab72` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 8.8117 | 1 | `143eafb5c72914c7` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | 8.7456 | 2 | `3d6e7ab8a9d2d9ac` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 100% | 8.5401 | 2 | `61a52babce16bb39` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 1% | 8.3248 | 4 | `4085881595386429` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 8.4238 | 4 | `998578789d48b8c2` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | 9.2116 | 8 | `6039c78bbdcd81f8` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 100% | 8.8693 | 8 | `d21dcc95a3bf78ae` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 1% | 9.4584 | 16 | `d9f5c886922ac5bd` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 10.1585 | 16 | `0ad8c4af6c5cd5fd` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 6.9244 | 1 | `ff49c8e7207cee76` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 8.5240 | 1 | `3ae4ae6412eafec1` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | 13.1376 | 2 | `4eecbcb8dd2cac99` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 100% | 15.8540 | 2 | `7746a0ede4756129` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 1% | 24.3854 | 4 | `3ecf93bd2f911acc` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 31.0500 | 4 | `08d8852b8756a6ca` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | 56.9874 | 8 | `eaec6c400e1aae9b` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 100% | 59.7160 | 8 | `2d4b80ab1946b7fc` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 1% | 106.8276 | 16 | `d4ddc9f769d31f2d` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 91.3246 | 16 | `5142d6d94c0d564d` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 1.1718 | false |
| notmat | 1 | 1% | 0.9511 | false |
| mat | 1 | 100% | 1.0055 | false |
| notmat | 1 | 100% | 0.9727 | false |
| mat | 2 | 1% | 1.0724 | false |
| notmat | 2 | 1% | 1.6109 | false |
| mat | 2 | 100% | 0.9242 | false |
| notmat | 2 | 100% | 1.7158 | false |
| mat | 4 | 1% | 1.0333 | false |
| notmat | 4 | 1% | 3.0268 | false |
| mat | 4 | 100% | 0.8986 | false |
| notmat | 4 | 100% | 3.3122 | false |
| mat | 8 | 1% | 1.0511 | false |
| notmat | 8 | 1% | 6.5024 | false |
| mat | 8 | 100% | 0.9769 | false |
| notmat | 8 | 100% | 6.5773 | false |
| mat | 16 | 1% | 1.0126 | false |
| notmat | 16 | 1% | 11.4373 | false |
| mat | 16 | 100% | 1.0026 | false |
| notmat | 16 | 100% | 9.0138 | false |

## Repeated-design position control

| first median ms | last median ms | ratio (last/first) | same plan |
|---|---|---|---|
| 129.605220 | 108.967709 | 0.8408 | false |

## Physical arm

- `charity_pkey`: `CREATE UNIQUE INDEX charity_pkey ON selectjoin_domain.charity USING lsm (charity_id HASH)`
- `donation_person_idx`: `CREATE INDEX donation_person_idx ON selectjoin_domain.donation USING lsm (person_id HASH)`
- `donation_pkey`: `CREATE UNIQUE INDEX donation_pkey ON selectjoin_domain.donation USING lsm (donation_id HASH)`
- `donation_rank_idx`: `CREATE INDEX donation_rank_idx ON selectjoin_domain.donation USING lsm (rank ASC)`
- `person_pkey`: `CREATE UNIQUE INDEX person_pkey ON selectjoin_domain.person USING lsm (person_id HASH)`

- index arm `yb-range-rank`; placement `engine default`; node budget `2 CPU 3g`; client budget `1 CPU 2g`

| table | tablets | hash key columns | colocated | tablegroup | colocation id |
|---|---|---|---|---|---|
| selectjoin_domain.charity | 1 | 1 | false | unknown | unknown |
| selectjoin_domain.person | 1 | 1 | false | unknown | unknown |
| selectjoin_domain.donation | 1 | 1 | false | unknown | unknown |
| selectjoin_domain.donation_pkey | 0 | 0 | false | unknown | unknown |
| selectjoin_domain.donation_rank_idx | 1 | 0 | false | unknown | unknown |
| selectjoin_domain.donation_person_idx | 1 | 1 | false | unknown | unknown |

Physical-arm gate (default join must use the rank index at its three most selective thresholds):

- threshold 20: used
- threshold 200: used
- threshold 2000: used

## Throttling (cgroup `cpu.stat` deltas)

- client: nr_periods 647, nr_throttled 254, throttled_usec 9634765 | quota 100000 100000
, memory 2147483648

- yb-single: nr_periods 714, nr_throttled 80, throttled_usec 3524146 | quota 200000 100000, memory 3221225472

## Access path per cell

_Engine-reported from each cell's own plan: whether the rank index or a sequential scan was used, rows scanned, PostgreSQL buffers and filter removals, and reader retries of transient conflicts. A cell that scanned the whole table at a selective threshold is a different physical regime, not a faster query._

| axis | status | index scan | seq scan | storage rows scanned | index rows scanned | buffers hit/read | rows removed | retries |
|---|---|---|---|---|---|---|---|---|
| family=bag;stmt=comma_join;threshold=100000 | valid | false | true | 202050 | 0 | 0/0 | 0 | 0 |
| family=bag;stmt=comma_join;threshold=180000 | valid | false | true | 202050 | 0 | 0/0 | 0 | 0 |
| family=bag;stmt=comma_join;threshold=180000;position=last | valid | false | true | 202050 | 0 | 0/0 | 0 | 0 |
| family=bag;stmt=comma_join;threshold=20 | valid | true | true | 110 | 20 | 0/0 | 0 | 0 |
| family=bag;stmt=comma_join;threshold=200 | valid | true | true | 2450 | 200 | 0/0 | 0 | 0 |
| family=bag;stmt=comma_join;threshold=2000 | valid | true | true | 6050 | 2000 | 0/0 | 0 | 0 |
| family=bag;stmt=comma_join;threshold=20000 | valid | true | true | 42050 | 20000 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_filtered;threshold=100000 | valid | false | true | 202050 | 0 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_filtered;threshold=180000 | valid | false | true | 202050 | 0 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_filtered;threshold=20 | valid | true | true | 110 | 20 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_filtered;threshold=200 | valid | true | true | 2450 | 200 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_filtered;threshold=2000 | valid | true | true | 6050 | 2000 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_filtered;threshold=20000 | valid | true | true | 42050 | 20000 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_filtered_mat;threshold=100000 | valid | false | true | 202050 | 0 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_filtered_mat;threshold=180000 | valid | false | true | 202050 | 0 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_filtered_mat;threshold=20 | valid | true | true | 110 | 20 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_filtered_mat;threshold=200 | valid | true | true | 2450 | 200 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_filtered_mat;threshold=2000 | valid | true | true | 6050 | 2000 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_filtered_mat;threshold=20000 | valid | true | true | 42050 | 20000 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_filtered_notmat;threshold=100000 | valid | false | true | 202050 | 0 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_filtered_notmat;threshold=180000 | valid | false | true | 202050 | 0 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_filtered_notmat;threshold=20 | valid | true | true | 110 | 20 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_filtered_notmat;threshold=200 | valid | true | true | 2450 | 200 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_filtered_notmat;threshold=2000 | valid | true | true | 6050 | 2000 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_filtered_notmat;threshold=20000 | valid | true | true | 42050 | 20000 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_joined_outer;threshold=100000 | valid | false | true | 202050 | 0 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | false | true | 202050 | 0 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | true | true | 110 | 20 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_joined_outer;threshold=200 | valid | true | true | 2450 | 200 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | true | true | 6050 | 2000 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_joined_outer;threshold=20000 | valid | true | true | 42050 | 20000 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_joined_outer_mat;threshold=100000 | valid | false | true | 202050 | 0 | 0/0 | 100000 | 0 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | false | true | 202050 | 0 | 0/0 | 20000 | 0 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | false | true | 202050 | 0 | 0/0 | 199980 | 0 |
| family=bag;stmt=cte_joined_outer_mat;threshold=200 | valid | false | true | 202050 | 0 | 0/0 | 199800 | 0 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | false | true | 202050 | 0 | 0/0 | 198000 | 0 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20000 | valid | false | true | 202050 | 0 | 0/0 | 180000 | 0 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=100000 | valid | false | true | 202050 | 0 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | false | true | 202050 | 0 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | true | true | 110 | 20 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=200 | valid | true | true | 2450 | 200 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | true | true | 6050 | 2000 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20000 | valid | true | true | 42050 | 20000 | 0/0 | 0 | 0 |
| family=bag;stmt=filtered_derived;threshold=100000 | valid | false | true | 202050 | 0 | 0/0 | 0 | 0 |
| family=bag;stmt=filtered_derived;threshold=180000 | valid | false | true | 202050 | 0 | 0/0 | 0 | 0 |
| family=bag;stmt=filtered_derived;threshold=20 | valid | true | true | 110 | 20 | 0/0 | 0 | 0 |
| family=bag;stmt=filtered_derived;threshold=200 | valid | true | true | 2450 | 200 | 0/0 | 0 | 0 |
| family=bag;stmt=filtered_derived;threshold=2000 | valid | true | true | 6050 | 2000 | 0/0 | 0 | 0 |
| family=bag;stmt=filtered_derived;threshold=20000 | valid | true | true | 42050 | 20000 | 0/0 | 0 | 0 |
| family=bag;stmt=join_on;threshold=100000 | valid | false | true | 202050 | 0 | 0/0 | 0 | 0 |
| family=bag;stmt=join_on;threshold=180000 | valid | false | true | 202050 | 0 | 0/0 | 0 | 0 |
| family=bag;stmt=join_on;threshold=20 | valid | true | true | 110 | 20 | 0/0 | 0 | 0 |
| family=bag;stmt=join_on;threshold=200 | valid | true | true | 2450 | 200 | 0/0 | 0 | 0 |
| family=bag;stmt=join_on;threshold=2000 | valid | true | true | 6050 | 2000 | 0/0 | 0 | 0 |
| family=bag;stmt=join_on;threshold=20000 | valid | true | true | 42050 | 20000 | 0/0 | 0 | 0 |
| family=bag;stmt=join_where;threshold=100000 | valid | false | true | 202050 | 0 | 0/0 | 0 | 0 |
| family=bag;stmt=join_where;threshold=180000 | valid | false | true | 202050 | 0 | 0/0 | 0 | 0 |
| family=bag;stmt=join_where;threshold=20 | valid | true | true | 110 | 20 | 0/0 | 0 | 0 |
| family=bag;stmt=join_where;threshold=200 | valid | true | true | 2450 | 200 | 0/0 | 0 | 0 |
| family=bag;stmt=join_where;threshold=2000 | valid | true | true | 6050 | 2000 | 0/0 | 0 | 0 |
| family=bag;stmt=join_where;threshold=20000 | valid | true | true | 42050 | 20000 | 0/0 | 0 | 0 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | true | false | 4000 | 2000 | 0/0 | 1983 | 0 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | true | false | 4000 | 2000 | 0/0 | 0 | 0 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | true | false | 4000 | 2000 | 0/0 | 0 | 0 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | true | false | 4000 | 2000 | 0/0 | 0 | 0 |
| family=reuse;policy=default;refs=2;consumer=1% | valid | true | false | 4000 | 2000 | 0/0 | 1983 | 0 |
| family=reuse;policy=default;refs=2;consumer=100% | valid | true | false | 4000 | 2000 | 0/0 | 0 | 0 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | true | false | 4000 | 2000 | 0/0 | 1983 | 0 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | true | false | 4000 | 2000 | 0/0 | 0 | 0 |
| family=reuse;policy=default;refs=8;consumer=1% | valid | true | false | 4000 | 2000 | 0/0 | 1983 | 0 |
| family=reuse;policy=default;refs=8;consumer=100% | valid | true | false | 4000 | 2000 | 0/0 | 0 | 0 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | true | false | 4000 | 2000 | 0/0 | 1983 | 0 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | true | false | 4000 | 2000 | 0/0 | 0 | 0 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | true | false | 4000 | 2000 | 0/0 | 1983 | 0 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | true | false | 4000 | 2000 | 0/0 | 0 | 0 |
| family=reuse;policy=mat;refs=2;consumer=1% | valid | true | false | 4000 | 2000 | 0/0 | 1983 | 0 |
| family=reuse;policy=mat;refs=2;consumer=100% | valid | true | false | 4000 | 2000 | 0/0 | 0 | 0 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | true | false | 4000 | 2000 | 0/0 | 1983 | 0 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | true | false | 4000 | 2000 | 0/0 | 0 | 0 |
| family=reuse;policy=mat;refs=8;consumer=1% | valid | true | false | 4000 | 2000 | 0/0 | 1983 | 0 |
| family=reuse;policy=mat;refs=8;consumer=100% | valid | true | false | 4000 | 2000 | 0/0 | 0 | 0 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | true | false | 64000 | 2000 | 0/0 | 0 | 0 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | true | false | 64000 | 2000 | 0/0 | 0 | 0 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | true | false | 4000 | 2000 | 0/0 | 0 | 0 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | true | false | 4000 | 2000 | 0/0 | 0 | 0 |
| family=reuse;policy=notmat;refs=2;consumer=1% | valid | true | false | 8000 | 2000 | 0/0 | 0 | 0 |
| family=reuse;policy=notmat;refs=2;consumer=100% | valid | true | false | 8000 | 2000 | 0/0 | 0 | 0 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | true | false | 16000 | 2000 | 0/0 | 0 | 0 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | true | false | 16000 | 2000 | 0/0 | 0 | 0 |
| family=reuse;policy=notmat;refs=8;consumer=1% | valid | true | false | 32000 | 2000 | 0/0 | 0 | 0 |
| family=reuse;policy=notmat;refs=8;consumer=100% | valid | true | false | 32000 | 2000 | 0/0 | 0 | 0 |

## Latency samples and derived throughput

_The reader is sequential, so `queries/s` is 1000/median_ms, a derived rate and not a concurrent-capacity test. min/max are the five sample extremes. Each sample includes any retry and backoff time the statement needed after a transient conflict; per-cell retry counts appear in the access-path table._

| axis | status | min ms | median ms | max ms | samples | derived queries/s |
|---|---|---|---|---|---|---|
| family=bag;stmt=comma_join;threshold=100000 | valid | 90.0341 | 91.7894 | 95.0134 | 5 | 10.9 |
| family=bag;stmt=comma_join;threshold=180000 | valid | 116.5140 | 129.6052 | 145.3516 | 5 | 7.7 |
| family=bag;stmt=comma_join;threshold=180000;position=last | valid | 105.7962 | 108.9677 | 115.1161 | 5 | 9.2 |
| family=bag;stmt=comma_join;threshold=20 | valid | 1.7254 | 1.9191 | 2.3491 | 5 | 521.1 |
| family=bag;stmt=comma_join;threshold=200 | valid | 2.9338 | 3.2192 | 3.8786 | 5 | 310.6 |
| family=bag;stmt=comma_join;threshold=2000 | valid | 9.3956 | 10.4146 | 10.9946 | 5 | 96.0 |
| family=bag;stmt=comma_join;threshold=20000 | valid | 65.6225 | 67.2274 | 69.8070 | 5 | 14.9 |
| family=bag;stmt=cte_filtered;threshold=100000 | valid | 88.8456 | 95.5874 | 112.7286 | 5 | 10.5 |
| family=bag;stmt=cte_filtered;threshold=180000 | valid | 102.1540 | 108.9847 | 113.8875 | 5 | 9.2 |
| family=bag;stmt=cte_filtered;threshold=20 | valid | 1.6524 | 1.7762 | 1.9326 | 5 | 563.0 |
| family=bag;stmt=cte_filtered;threshold=200 | valid | 3.0858 | 3.6180 | 6.4280 | 5 | 276.4 |
| family=bag;stmt=cte_filtered;threshold=2000 | valid | 8.8917 | 9.0673 | 9.6712 | 5 | 110.3 |
| family=bag;stmt=cte_filtered;threshold=20000 | valid | 64.0105 | 68.0607 | 74.1242 | 5 | 14.7 |
| family=bag;stmt=cte_filtered_mat;threshold=100000 | valid | 89.5223 | 96.8301 | 102.4832 | 5 | 10.3 |
| family=bag;stmt=cte_filtered_mat;threshold=180000 | valid | 115.6127 | 121.2151 | 148.2920 | 5 | 8.2 |
| family=bag;stmt=cte_filtered_mat;threshold=20 | valid | 1.8352 | 1.8862 | 2.2762 | 5 | 530.2 |
| family=bag;stmt=cte_filtered_mat;threshold=200 | valid | 3.4463 | 3.6131 | 4.7123 | 5 | 276.8 |
| family=bag;stmt=cte_filtered_mat;threshold=2000 | valid | 9.4210 | 10.3554 | 11.2167 | 5 | 96.6 |
| family=bag;stmt=cte_filtered_mat;threshold=20000 | valid | 64.2284 | 67.1924 | 68.8834 | 5 | 14.9 |
| family=bag;stmt=cte_filtered_notmat;threshold=100000 | valid | 84.3694 | 91.5077 | 97.8536 | 5 | 10.9 |
| family=bag;stmt=cte_filtered_notmat;threshold=180000 | valid | 108.1701 | 111.1850 | 124.5009 | 5 | 9.0 |
| family=bag;stmt=cte_filtered_notmat;threshold=20 | valid | 2.1991 | 2.4434 | 2.6767 | 5 | 409.3 |
| family=bag;stmt=cte_filtered_notmat;threshold=200 | valid | 2.8730 | 3.2571 | 3.6284 | 5 | 307.0 |
| family=bag;stmt=cte_filtered_notmat;threshold=2000 | valid | 8.9851 | 9.5377 | 10.5975 | 5 | 104.8 |
| family=bag;stmt=cte_filtered_notmat;threshold=20000 | valid | 64.7566 | 68.1455 | 69.5271 | 5 | 14.7 |
| family=bag;stmt=cte_joined_outer;threshold=100000 | valid | 88.5647 | 91.6823 | 98.5374 | 5 | 10.9 |
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 105.9813 | 111.0178 | 115.1532 | 5 | 9.0 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 1.8116 | 2.0400 | 2.4096 | 5 | 490.2 |
| family=bag;stmt=cte_joined_outer;threshold=200 | valid | 3.5322 | 4.2330 | 4.3215 | 5 | 236.2 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 9.3328 | 9.6825 | 9.9069 | 5 | 103.3 |
| family=bag;stmt=cte_joined_outer;threshold=20000 | valid | 62.5668 | 65.8666 | 82.8601 | 5 | 15.2 |
| family=bag;stmt=cte_joined_outer_mat;threshold=100000 | valid | 89.8137 | 100.7127 | 107.5103 | 5 | 9.9 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 106.6292 | 113.7440 | 116.3122 | 5 | 8.8 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 81.3760 | 88.4830 | 92.3386 | 5 | 11.3 |
| family=bag;stmt=cte_joined_outer_mat;threshold=200 | valid | 82.8741 | 86.6851 | 91.9429 | 5 | 11.5 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 81.2698 | 85.8304 | 92.5383 | 5 | 11.7 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20000 | valid | 85.6120 | 90.7192 | 95.0719 | 5 | 11.0 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=100000 | valid | 84.1101 | 91.1699 | 93.0840 | 5 | 11.0 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 110.9669 | 115.3694 | 148.2352 | 5 | 8.7 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 1.7495 | 1.9852 | 2.4285 | 5 | 503.7 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=200 | valid | 3.4874 | 3.7862 | 3.9417 | 5 | 264.1 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 9.1527 | 10.0079 | 10.3020 | 5 | 99.9 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20000 | valid | 60.7738 | 64.8759 | 67.8351 | 5 | 15.4 |
| family=bag;stmt=filtered_derived;threshold=100000 | valid | 89.1180 | 89.8128 | 92.3775 | 5 | 11.1 |
| family=bag;stmt=filtered_derived;threshold=180000 | valid | 109.7725 | 118.3836 | 129.8024 | 5 | 8.4 |
| family=bag;stmt=filtered_derived;threshold=20 | valid | 1.8241 | 1.9032 | 2.4830 | 5 | 525.4 |
| family=bag;stmt=filtered_derived;threshold=200 | valid | 3.4709 | 3.6986 | 4.5477 | 5 | 270.4 |
| family=bag;stmt=filtered_derived;threshold=2000 | valid | 8.3828 | 9.3428 | 9.7750 | 5 | 107.0 |
| family=bag;stmt=filtered_derived;threshold=20000 | valid | 64.8428 | 67.4144 | 68.2803 | 5 | 14.8 |
| family=bag;stmt=join_on;threshold=100000 | valid | 83.9656 | 89.8909 | 92.3096 | 5 | 11.1 |
| family=bag;stmt=join_on;threshold=180000 | valid | 104.1885 | 107.3997 | 113.3772 | 5 | 9.3 |
| family=bag;stmt=join_on;threshold=20 | valid | 1.7050 | 1.9798 | 2.5057 | 5 | 505.1 |
| family=bag;stmt=join_on;threshold=200 | valid | 3.0329 | 3.1592 | 3.7948 | 5 | 316.5 |
| family=bag;stmt=join_on;threshold=2000 | valid | 8.3290 | 8.4765 | 10.2176 | 5 | 118.0 |
| family=bag;stmt=join_on;threshold=20000 | valid | 62.2193 | 62.5451 | 74.9313 | 5 | 16.0 |
| family=bag;stmt=join_where;threshold=100000 | valid | 83.6083 | 86.5406 | 92.7615 | 5 | 11.6 |
| family=bag;stmt=join_where;threshold=180000 | valid | 102.0053 | 108.4982 | 115.9880 | 5 | 9.2 |
| family=bag;stmt=join_where;threshold=20 | valid | 1.7524 | 2.1515 | 2.5089 | 5 | 464.8 |
| family=bag;stmt=join_where;threshold=200 | valid | 2.3110 | 3.1810 | 3.8979 | 5 | 314.4 |
| family=bag;stmt=join_where;threshold=2000 | valid | 9.3066 | 9.6082 | 10.4800 | 5 | 104.1 |
| family=bag;stmt=join_where;threshold=20000 | valid | 63.3664 | 66.7985 | 69.8839 | 5 | 15.0 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 8.8325 | 9.3403 | 9.7822 | 5 | 107.1 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 9.9122 | 10.1317 | 10.9425 | 5 | 98.7 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 6.1796 | 7.2801 | 8.1874 | 5 | 137.4 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 8.5488 | 8.7633 | 9.4130 | 5 | 114.1 |
| family=reuse;policy=default;refs=2;consumer=1% | valid | 7.6540 | 8.1555 | 8.8456 | 5 | 122.6 |
| family=reuse;policy=default;refs=2;consumer=100% | valid | 8.5343 | 9.2403 | 9.5505 | 5 | 108.2 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 7.5639 | 8.0565 | 8.3698 | 5 | 124.1 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 8.8476 | 9.3743 | 10.0652 | 5 | 106.7 |
| family=reuse;policy=default;refs=8;consumer=1% | valid | 8.3252 | 8.7641 | 9.1940 | 5 | 114.1 |
| family=reuse;policy=default;refs=8;consumer=100% | valid | 8.7426 | 9.0791 | 9.2416 | 5 | 110.1 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 8.9017 | 9.4584 | 9.8793 | 5 | 105.7 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 9.9048 | 10.1585 | 10.7815 | 5 | 98.4 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 8.0439 | 8.5311 | 9.9824 | 5 | 117.2 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 8.4377 | 8.8117 | 43.9753 | 5 | 113.5 |
| family=reuse;policy=mat;refs=2;consumer=1% | valid | 8.4152 | 8.7456 | 9.5141 | 5 | 114.3 |
| family=reuse;policy=mat;refs=2;consumer=100% | valid | 8.0377 | 8.5401 | 9.4108 | 5 | 117.1 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 7.9585 | 8.3248 | 9.4126 | 5 | 120.1 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 8.0755 | 8.4238 | 9.0656 | 5 | 118.7 |
| family=reuse;policy=mat;refs=8;consumer=1% | valid | 9.0576 | 9.2116 | 9.9108 | 5 | 108.6 |
| family=reuse;policy=mat;refs=8;consumer=100% | valid | 8.4224 | 8.8693 | 10.2133 | 5 | 112.7 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 79.2058 | 106.8276 | 148.5089 | 5 | 9.4 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 84.5355 | 91.3246 | 126.8360 | 5 | 10.9 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 6.2311 | 6.9244 | 7.5476 | 5 | 144.4 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 8.4340 | 8.5240 | 9.7278 | 5 | 117.3 |
| family=reuse;policy=notmat;refs=2;consumer=1% | valid | 12.4354 | 13.1376 | 14.5360 | 5 | 76.1 |
| family=reuse;policy=notmat;refs=2;consumer=100% | valid | 14.7506 | 15.8540 | 16.5253 | 5 | 63.1 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 22.7781 | 24.3854 | 25.6378 | 5 | 41.0 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 30.1765 | 31.0500 | 31.9242 | 5 | 32.2 |
| family=reuse;policy=notmat;refs=8;consumer=1% | valid | 54.6466 | 56.9874 | 67.8909 | 5 | 17.5 |
| family=reuse;policy=notmat;refs=8;consumer=100% | valid | 54.9964 | 59.7160 | 86.0307 | 5 | 16.7 |

## Matched blended workload

- phase: `isolated`; target offered rate: 10 updates/s
- offered 0, committed 0, failed 0 over 0.000 s
- offered rate 0.000/s, completed rate 0.000/s
- per key: no keys touched
- ledger: 0 keys reconciled, 0 mismatch

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42119389`, 5 interleaved rounds
  1. family=bag;stmt=comma_join;threshold=180000
  2. family=bag;stmt=cte_joined_outer_mat;threshold=200
  3. family=bag;stmt=cte_filtered;threshold=200
  4. family=bag;stmt=cte_filtered;threshold=180000
  5. family=bag;stmt=cte_filtered_notmat;threshold=20
  6. family=reuse;policy=notmat;refs=4;consumer=100%
  7. family=bag;stmt=cte_joined_outer_mat;threshold=20000
  8. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  9. family=reuse;policy=default;refs=1;consumer=100%
  10. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  11. family=bag;stmt=cte_joined_outer_mat;threshold=20
  12. family=bag;stmt=join_where;threshold=20
  13. family=reuse;policy=notmat;refs=16;consumer=1%
  14. family=bag;stmt=cte_joined_outer;threshold=20
  15. family=bag;stmt=cte_joined_outer;threshold=2000
  16. family=bag;stmt=comma_join;threshold=200
  17. family=bag;stmt=cte_joined_outer;threshold=20000
  18. family=bag;stmt=filtered_derived;threshold=180000
  19. family=bag;stmt=join_on;threshold=20000
  20. family=bag;stmt=comma_join;threshold=2000
  21. family=bag;stmt=filtered_derived;threshold=100000
  22. family=reuse;policy=mat;refs=2;consumer=1%
  23. family=reuse;policy=default;refs=4;consumer=1%
  24. family=bag;stmt=join_on;threshold=2000
  25. family=bag;stmt=join_where;threshold=100000
  26. family=bag;stmt=filtered_derived;threshold=20000
  27. family=reuse;policy=notmat;refs=1;consumer=100%
  28. family=bag;stmt=cte_filtered;threshold=20
  29. family=reuse;policy=default;refs=8;consumer=1%
  30. family=reuse;policy=mat;refs=1;consumer=1%
  31. family=bag;stmt=join_where;threshold=2000
  32. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  33. family=bag;stmt=cte_filtered_notmat;threshold=2000
  34. family=bag;stmt=join_on;threshold=100000
  35. family=reuse;policy=notmat;refs=8;consumer=1%
  36. family=bag;stmt=join_where;threshold=20000
  37. family=reuse;policy=notmat;refs=2;consumer=1%
  38. family=reuse;policy=default;refs=8;consumer=100%
  39. family=bag;stmt=join_on;threshold=200
  40. family=bag;stmt=cte_filtered_mat;threshold=20
  41. family=bag;stmt=comma_join;threshold=100000
  42. family=bag;stmt=cte_joined_outer;threshold=100000
  43. family=reuse;policy=notmat;refs=8;consumer=100%
  44. family=reuse;policy=default;refs=4;consumer=100%
  45. family=reuse;policy=default;refs=1;consumer=1%
  46. family=bag;stmt=filtered_derived;threshold=2000
  47. family=reuse;policy=mat;refs=8;consumer=1%
  48. family=reuse;policy=default;refs=2;consumer=1%
  49. family=reuse;policy=mat;refs=8;consumer=100%
  50. family=reuse;policy=default;refs=16;consumer=1%
  51. family=reuse;policy=notmat;refs=2;consumer=100%
  52. family=bag;stmt=cte_filtered_notmat;threshold=200
  53. family=bag;stmt=cte_filtered_mat;threshold=100000
  54. family=bag;stmt=cte_joined_outer;threshold=200
  55. family=bag;stmt=cte_filtered_notmat;threshold=180000
  56. family=bag;stmt=cte_joined_outer_notmat;threshold=20000
  57. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  58. family=bag;stmt=cte_filtered_mat;threshold=20000
  59. family=bag;stmt=cte_filtered;threshold=20000
  60. family=reuse;policy=default;refs=2;consumer=100%
  61. family=bag;stmt=comma_join;threshold=20
  62. family=reuse;policy=mat;refs=2;consumer=100%
  63. family=bag;stmt=filtered_derived;threshold=20
  64. family=reuse;policy=notmat;refs=4;consumer=1%
  65. family=reuse;policy=notmat;refs=1;consumer=1%
  66. family=bag;stmt=cte_filtered;threshold=2000
  67. family=bag;stmt=cte_filtered_mat;threshold=180000
  68. family=bag;stmt=join_where;threshold=180000
  69. family=bag;stmt=cte_joined_outer_mat;threshold=100000
  70. family=bag;stmt=cte_filtered_mat;threshold=200
  71. family=reuse;policy=mat;refs=4;consumer=1%
  72. family=bag;stmt=cte_filtered_notmat;threshold=20000
  73. family=bag;stmt=comma_join;threshold=20000
  74. family=bag;stmt=cte_joined_outer_notmat;threshold=200
  75. family=bag;stmt=join_where;threshold=200
  76. family=reuse;policy=notmat;refs=16;consumer=100%
  77. family=bag;stmt=cte_filtered_mat;threshold=2000
  78. family=bag;stmt=join_on;threshold=20
  79. family=reuse;policy=mat;refs=16;consumer=1%
  80. family=bag;stmt=cte_filtered;threshold=100000
  81. family=bag;stmt=filtered_derived;threshold=200
  82. family=reuse;policy=mat;refs=4;consumer=100%
  83. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  84. family=bag;stmt=join_on;threshold=180000
  85. family=reuse;policy=mat;refs=16;consumer=100%
  86. family=bag;stmt=cte_joined_outer_notmat;threshold=100000
  87. family=bag;stmt=cte_filtered_notmat;threshold=100000
  88. family=reuse;policy=mat;refs=1;consumer=100%
  89. family=bag;stmt=cte_joined_outer;threshold=180000
  90. family=reuse;policy=default;refs=16;consumer=100%
  91. family=bag;stmt=comma_join;threshold=180000;position=last

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
| family=bag;stmt=comma_join;threshold=100000 | j06_comma_join | valid | `ee3611feaa6ec8b1` |  |
| family=bag;stmt=comma_join;threshold=180000 | j06_comma_join | valid | `e56972ccf1afb5a8` |  |
| family=bag;stmt=comma_join;threshold=180000;position=last | j06_comma_join | valid | `0e4dd666a43e3c3d` |  |
| family=bag;stmt=comma_join;threshold=20 | j06_comma_join | valid | `c584b1730f2f1c6d` |  |
| family=bag;stmt=comma_join;threshold=200 | j06_comma_join | valid | `6f9c30f003b18775` |  |
| family=bag;stmt=comma_join;threshold=2000 | j06_comma_join | valid | `3913059b7fb90da9` |  |
| family=bag;stmt=comma_join;threshold=20000 | j06_comma_join | valid | `659ad906d6b7e1b5` |  |
| family=bag;stmt=cte_filtered;threshold=100000 | j04_cte_filtered | valid | `8241121c1ef109c4` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=180000 | j04_cte_filtered | valid | `67925af93299dce1` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=20 | j04_cte_filtered | valid | `0a475f20cfbc0f01` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=200 | j04_cte_filtered | valid | `e6e6b7999e16ab02` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=2000 | j04_cte_filtered | valid | `42a514445d6c73d5` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=20000 | j04_cte_filtered | valid | `78e77696338170ec` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_mat;threshold=100000 | j04m_cte_filtered_mat | valid | `13ea7bc5ba18960e` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=180000 | j04m_cte_filtered_mat | valid | `06ffcc20a5cb6928` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=20 | j04m_cte_filtered_mat | valid | `50fe15cc77e839a3` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=200 | j04m_cte_filtered_mat | valid | `d2d5c3df74a77bd2` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=2000 | j04m_cte_filtered_mat | valid | `cafad85e51b969aa` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=20000 | j04m_cte_filtered_mat | valid | `8a3dff4b2dc15bf6` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_notmat;threshold=100000 | j04n_cte_filtered_notmat | valid | `e911965e013003b0` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=180000 | j04n_cte_filtered_notmat | valid | `c5fad1ac3b52ef1c` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=20 | j04n_cte_filtered_notmat | valid | `d3d7cf948c064daf` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=200 | j04n_cte_filtered_notmat | valid | `4918a243ce7acab8` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=2000 | j04n_cte_filtered_notmat | valid | `efcfcd059313f2f4` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=20000 | j04n_cte_filtered_notmat | valid | `ad109ce3b79b1a7c` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=100000 | j05_cte_joined_outer | valid | `29a0bcefdf6d3ec4` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=180000 | j05_cte_joined_outer | valid | `42f39348d08c34f5` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20 | j05_cte_joined_outer | valid | `3eb4c5b0bda40919` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=200 | j05_cte_joined_outer | valid | `c92817736907848b` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=2000 | j05_cte_joined_outer | valid | `039710944b4a0475` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20000 | j05_cte_joined_outer | valid | `e795fa296374b3aa` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_mat;threshold=100000 | j05m_cte_joined_outer_mat | valid | `a673f4dfd13bf067` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | j05m_cte_joined_outer_mat | valid | `918bf029599b2477` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | j05m_cte_joined_outer_mat | valid | `9263eac7fdd6c067` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=200 | j05m_cte_joined_outer_mat | valid | `e4d7221a0bd746d3` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | j05m_cte_joined_outer_mat | valid | `1a0dae077d3e97c6` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20000 | j05m_cte_joined_outer_mat | valid | `ac85b9b0a2af1ae3` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_notmat;threshold=100000 | j05n_cte_joined_outer_notmat | valid | `976cde6d0f17e1f5` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | j05n_cte_joined_outer_notmat | valid | `fefe2a7140c2ca44` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | j05n_cte_joined_outer_notmat | valid | `6771b8a44573d9d9` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=200 | j05n_cte_joined_outer_notmat | valid | `188a50c3a14e121d` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | j05n_cte_joined_outer_notmat | valid | `65d73e40c85f5251` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20000 | j05n_cte_joined_outer_notmat | valid | `d49edafb0d95dc91` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=filtered_derived;threshold=100000 | j03_filtered_derived | valid | `4a9443a4e5447c6b` |  |
| family=bag;stmt=filtered_derived;threshold=180000 | j03_filtered_derived | valid | `7dee6f5587cac023` |  |
| family=bag;stmt=filtered_derived;threshold=20 | j03_filtered_derived | valid | `3d9371af1f139502` |  |
| family=bag;stmt=filtered_derived;threshold=200 | j03_filtered_derived | valid | `5a3b685ecee028fd` |  |
| family=bag;stmt=filtered_derived;threshold=2000 | j03_filtered_derived | valid | `2cb38da86758e4f6` |  |
| family=bag;stmt=filtered_derived;threshold=20000 | j03_filtered_derived | valid | `bd420602cfb24b47` |  |
| family=bag;stmt=join_on;threshold=100000 | j02_join_on | valid | `13cf89052da9eba4` |  |
| family=bag;stmt=join_on;threshold=180000 | j02_join_on | valid | `c67d69da2067f7e7` |  |
| family=bag;stmt=join_on;threshold=20 | j02_join_on | valid | `9ab3a53b0e18a392` |  |
| family=bag;stmt=join_on;threshold=200 | j02_join_on | valid | `298e2429a99e6424` |  |
| family=bag;stmt=join_on;threshold=2000 | j02_join_on | valid | `2570db5a7c57004f` |  |
| family=bag;stmt=join_on;threshold=20000 | j02_join_on | valid | `06d0b4adcef330fa` |  |
| family=bag;stmt=join_where;threshold=100000 | j01_join_where | valid | `f93a1694494097bf` |  |
| family=bag;stmt=join_where;threshold=180000 | j01_join_where | valid | `d6e86336279de0f4` |  |
| family=bag;stmt=join_where;threshold=20 | j01_join_where | valid | `828773bb6094e7a5` |  |
| family=bag;stmt=join_where;threshold=200 | j01_join_where | valid | `6b373ace58e66229` |  |
| family=bag;stmt=join_where;threshold=2000 | j01_join_where | valid | `81bb41178acfb814` |  |
| family=bag;stmt=join_where;threshold=20000 | j01_join_where | valid | `4d10a4b63cb856bb` |  |
| family=reuse;policy=default;refs=16;consumer=1% | r_default_n16_c1 | valid | `f69657e41cd7f916` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=16;consumer=100% | r_default_n16_c100 | valid | `349c18ad78591646` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=1;consumer=1% | r_default_n1_c1 | valid | `1e464b3f33f72b04` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=1;consumer=100% | r_default_n1_c100 | valid | `fbac5a7a22297229` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=2;consumer=1% | r_default_n2_c1 | valid | `bc0b0617e1357ec5` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=2;consumer=100% | r_default_n2_c100 | valid | `4e205dbd55a42666` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=1% | r_default_n4_c1 | valid | `3865366089004bbd` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=100% | r_default_n4_c100 | valid | `2bbfe6ff213a1d88` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=8;consumer=1% | r_default_n8_c1 | valid | `62fbb2d6d3c66d6c` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=8;consumer=100% | r_default_n8_c100 | valid | `18d8ad5a2031c6cf` | policy=default; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=1% | r_mat_n16_c1 | valid | `d9f5c886922ac5bd` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=100% | r_mat_n16_c100 | valid | `0ad8c4af6c5cd5fd` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=1% | r_mat_n1_c1 | valid | `0bfec48ca42bab72` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=100% | r_mat_n1_c100 | valid | `143eafb5c72914c7` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=2;consumer=1% | r_mat_n2_c1 | valid | `3d6e7ab8a9d2d9ac` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=2;consumer=100% | r_mat_n2_c100 | valid | `61a52babce16bb39` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=1% | r_mat_n4_c1 | valid | `4085881595386429` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=100% | r_mat_n4_c100 | valid | `998578789d48b8c2` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=8;consumer=1% | r_mat_n8_c1 | valid | `6039c78bbdcd81f8` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=8;consumer=100% | r_mat_n8_c100 | valid | `d21dcc95a3bf78ae` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=notmat;refs=16;consumer=1% | r_notmat_n16_c1 | valid | `d4ddc9f769d31f2d` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `5142d6d94c0d564d` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `ff49c8e7207cee76` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `3ae4ae6412eafec1` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=2;consumer=1% | r_notmat_n2_c1 | valid | `4eecbcb8dd2cac99` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=2;consumer=100% | r_notmat_n2_c100 | valid | `7746a0ede4756129` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `3ecf93bd2f911acc` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `08d8852b8756a6ca` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=8;consumer=1% | r_notmat_n8_c1 | valid | `eaec6c400e1aae9b` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=8;consumer=100% | r_notmat_n8_c100 | valid | `2d4b80ab1946b7fc` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, DIST, DEBUG, FORMAT TEXT` supported on yugabyte
- one host: this run used node budget 2 CPU 3g and client budget 1 CPU 2g, so budget comparisons are separate classes here, not independent hardware
- isolated phase: no concurrent writers in this run; the blended read/write phase is a separate labelled run
- physical rank-index arm yb-range-rank; an engine whose planner ignores the index is a different physical design, not a faster one
- YugabyteDB DIST/DEBUG storage counters (rows and index rows scanned) are parsed per cell; RPC-level counters are not aggregated
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
