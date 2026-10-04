# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 91 cells, 91 valid, engine `yugabyte`, topology `single`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-v3-domain-yb-i1-20261005`
- run tag: `run/study-01/select-join-v3-domain-yb-i1-20261005`
- environment: `host-zenbook-ux5406sa`
- repo commit: `4d80c4a5b4cee2cb40c3189670cb317f74e28f78`
- engine version: `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit`
- database engine image: `docker.io/yugabytedb/yugabyte:2025.2.6.0-b111` @ `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`
- benchmark client image: `localhost/charitybench-select:select-join-v3-domain-yb-i1-20261005` @ `8436e36bb66caa8178559d54771b88c063e26848d4ef1d1a56eb524b7f0ecca6`
- inputs digest: `26978a5ec9fd2b19ac42e0e7f45cb1400fe88d85699fd2cdaecf91e11fec1feb`
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
| 20 (0.01%) | 2.1913 | 2.0977 | 1.5013 | 2.2764 | 1.5041 | 1.8373 | 1.7437 | 77.1213 | 1.5944 | 1.3876 |
| 200 (0.1%) | 3.4286 | 3.0373 | 3.3203 | 3.6375 | 2.6817 | 3.2763 | 2.7340 | 74.4573 | 3.3205 | 2.7627 |
| 2000 (1%) | 8.2696 | 8.8282 | 9.5214 | 8.2259 | 8.3444 | 7.8776 | 9.3601 | 79.3384 | 8.5146 | 8.2027 |
| 20000 (10%) | 58.1588 | 58.4553 | 59.6671 | 57.5979 | 60.8615 | 59.7119 | 58.0394 | 79.4571 | 60.9999 | 58.0668 |
| 100000 (50%) | 80.4412 | 83.3810 | 80.4251 | 80.0007 | 85.6265 | 80.0196 | 80.5928 | 84.5073 | 82.3307 | 82.3557 |
| 180000 (90%) | 98.7178 | 99.5377 | 99.0461 | 99.4646 | 101.9046 | 100.2052 | 98.0680 | 96.8588 | 99.1684 | 97.2575 |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | join_on vs join_where | 2.1913 | 2.0977 | 0.9573 | false |
| 20 | filtered_derived vs join_where | 2.1913 | 1.5013 | 0.6851 | false |
| 20 | cte_filtered vs join_where | 2.1913 | 2.2764 | 1.0388 | false |
| 20 | cte_filtered_mat vs join_where | 2.1913 | 1.5041 | 0.6864 | false |
| 20 | cte_filtered_notmat vs join_where | 2.1913 | 1.8373 | 0.8384 | false |
| 20 | cte_joined_outer vs join_where | 2.1913 | 1.7437 | 0.7957 | false |
| 20 | cte_joined_outer_mat vs join_where | 2.1913 | 77.1213 | 35.1937 | false |
| 20 | cte_joined_outer_notmat vs join_where | 2.1913 | 1.5944 | 0.7276 | false |
| 20 | comma_join vs join_where | 2.1913 | 1.3876 | 0.6332 | false |
| 200 | join_on vs join_where | 3.4286 | 3.0373 | 0.8859 | false |
| 200 | filtered_derived vs join_where | 3.4286 | 3.3203 | 0.9684 | false |
| 200 | cte_filtered vs join_where | 3.4286 | 3.6375 | 1.0609 | false |
| 200 | cte_filtered_mat vs join_where | 3.4286 | 2.6817 | 0.7822 | false |
| 200 | cte_filtered_notmat vs join_where | 3.4286 | 3.2763 | 0.9556 | false |
| 200 | cte_joined_outer vs join_where | 3.4286 | 2.7340 | 0.7974 | false |
| 200 | cte_joined_outer_mat vs join_where | 3.4286 | 74.4573 | 21.7165 | false |
| 200 | cte_joined_outer_notmat vs join_where | 3.4286 | 3.3205 | 0.9685 | false |
| 200 | comma_join vs join_where | 3.4286 | 2.7627 | 0.8058 | false |
| 2000 | join_on vs join_where | 8.2696 | 8.8282 | 1.0676 | false |
| 2000 | filtered_derived vs join_where | 8.2696 | 9.5214 | 1.1514 | false |
| 2000 | cte_filtered vs join_where | 8.2696 | 8.2259 | 0.9947 | false |
| 2000 | cte_filtered_mat vs join_where | 8.2696 | 8.3444 | 1.0090 | false |
| 2000 | cte_filtered_notmat vs join_where | 8.2696 | 7.8776 | 0.9526 | false |
| 2000 | cte_joined_outer vs join_where | 8.2696 | 9.3601 | 1.1319 | false |
| 2000 | cte_joined_outer_mat vs join_where | 8.2696 | 79.3384 | 9.5940 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 8.2696 | 8.5146 | 1.0296 | false |
| 2000 | comma_join vs join_where | 8.2696 | 8.2027 | 0.9919 | false |
| 20000 | join_on vs join_where | 58.1588 | 58.4553 | 1.0051 | false |
| 20000 | filtered_derived vs join_where | 58.1588 | 59.6671 | 1.0259 | false |
| 20000 | cte_filtered vs join_where | 58.1588 | 57.5979 | 0.9904 | false |
| 20000 | cte_filtered_mat vs join_where | 58.1588 | 60.8615 | 1.0465 | false |
| 20000 | cte_filtered_notmat vs join_where | 58.1588 | 59.7119 | 1.0267 | false |
| 20000 | cte_joined_outer vs join_where | 58.1588 | 58.0394 | 0.9979 | false |
| 20000 | cte_joined_outer_mat vs join_where | 58.1588 | 79.4571 | 1.3662 | false |
| 20000 | cte_joined_outer_notmat vs join_where | 58.1588 | 60.9999 | 1.0488 | false |
| 20000 | comma_join vs join_where | 58.1588 | 58.0668 | 0.9984 | false |
| 100000 | join_on vs join_where | 80.4412 | 83.3810 | 1.0365 | false |
| 100000 | filtered_derived vs join_where | 80.4412 | 80.4251 | 0.9998 | false |
| 100000 | cte_filtered vs join_where | 80.4412 | 80.0007 | 0.9945 | false |
| 100000 | cte_filtered_mat vs join_where | 80.4412 | 85.6265 | 1.0645 | false |
| 100000 | cte_filtered_notmat vs join_where | 80.4412 | 80.0196 | 0.9948 | false |
| 100000 | cte_joined_outer vs join_where | 80.4412 | 80.5928 | 1.0019 | false |
| 100000 | cte_joined_outer_mat vs join_where | 80.4412 | 84.5073 | 1.0505 | false |
| 100000 | cte_joined_outer_notmat vs join_where | 80.4412 | 82.3307 | 1.0235 | false |
| 100000 | comma_join vs join_where | 80.4412 | 82.3557 | 1.0238 | false |
| 180000 | join_on vs join_where | 98.7178 | 99.5377 | 1.0083 | false |
| 180000 | filtered_derived vs join_where | 98.7178 | 99.0461 | 1.0033 | false |
| 180000 | cte_filtered vs join_where | 98.7178 | 99.4646 | 1.0076 | false |
| 180000 | cte_filtered_mat vs join_where | 98.7178 | 101.9046 | 1.0323 | false |
| 180000 | cte_filtered_notmat vs join_where | 98.7178 | 100.2052 | 1.0151 | false |
| 180000 | cte_joined_outer vs join_where | 98.7178 | 98.0680 | 0.9934 | false |
| 180000 | cte_joined_outer_mat vs join_where | 98.7178 | 96.8588 | 0.9812 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 98.7178 | 99.1684 | 1.0046 | false |
| 180000 | comma_join vs join_where | 98.7178 | 97.2575 | 0.9852 | false |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_filtered_mat vs cte_filtered | 2.2764 | 1.5041 | 0.6607 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_filtered_notmat vs cte_filtered | 2.2764 | 1.8373 | 0.8071 | false | policy=not_materialized; observed_cte_scan=false |
| 20 | cte_joined_outer_mat vs cte_joined_outer | 1.7437 | 77.1213 | 44.2296 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 1.7437 | 1.5944 | 0.9144 | false | policy=not_materialized; observed_cte_scan=false |
| 200 | cte_filtered_mat vs cte_filtered | 3.6375 | 2.6817 | 0.7372 | false | policy=materialized; observed_cte_scan=true |
| 200 | cte_filtered_notmat vs cte_filtered | 3.6375 | 3.2763 | 0.9007 | false | policy=not_materialized; observed_cte_scan=false |
| 200 | cte_joined_outer_mat vs cte_joined_outer | 2.7340 | 74.4573 | 27.2340 | false | policy=materialized; observed_cte_scan=true |
| 200 | cte_joined_outer_notmat vs cte_joined_outer | 2.7340 | 3.3205 | 1.2145 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_filtered_mat vs cte_filtered | 8.2259 | 8.3444 | 1.0144 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_filtered_notmat vs cte_filtered | 8.2259 | 7.8776 | 0.9577 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 9.3601 | 79.3384 | 8.4762 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 9.3601 | 8.5146 | 0.9097 | false | policy=not_materialized; observed_cte_scan=false |
| 20000 | cte_filtered_mat vs cte_filtered | 57.5979 | 60.8615 | 1.0567 | false | policy=materialized; observed_cte_scan=true |
| 20000 | cte_filtered_notmat vs cte_filtered | 57.5979 | 59.7119 | 1.0367 | false | policy=not_materialized; observed_cte_scan=false |
| 20000 | cte_joined_outer_mat vs cte_joined_outer | 58.0394 | 79.4571 | 1.3690 | false | policy=materialized; observed_cte_scan=true |
| 20000 | cte_joined_outer_notmat vs cte_joined_outer | 58.0394 | 60.9999 | 1.0510 | false | policy=not_materialized; observed_cte_scan=false |
| 100000 | cte_filtered_mat vs cte_filtered | 80.0007 | 85.6265 | 1.0703 | false | policy=materialized; observed_cte_scan=true |
| 100000 | cte_filtered_notmat vs cte_filtered | 80.0007 | 80.0196 | 1.0002 | false | policy=not_materialized; observed_cte_scan=false |
| 100000 | cte_joined_outer_mat vs cte_joined_outer | 80.5928 | 84.5073 | 1.0486 | false | policy=materialized; observed_cte_scan=true |
| 100000 | cte_joined_outer_notmat vs cte_joined_outer | 80.5928 | 82.3307 | 1.0216 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_filtered_mat vs cte_filtered | 99.4646 | 101.9046 | 1.0245 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_filtered_notmat vs cte_filtered | 99.4646 | 100.2052 | 1.0074 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 98.0680 | 96.8588 | 0.9877 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 98.0680 | 99.1684 | 1.0112 | false | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 7.1893 | 1 | `bab144ad7db45979` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 7.4566 | 1 | `15409dd542ec76b4` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | 7.3971 | 2 | `e7bd0e71013b7d44` | policy=default; observed_cte_scan=true | valid |
| default | 2 | 100% | 8.3133 | 2 | `ae6de606dc096a49` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 1% | 8.4364 | 4 | `a1705411ecba7790` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 8.6619 | 4 | `0697b7d47cea9c02` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | 8.5018 | 8 | `8e53f03785f056d9` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 100% | 7.8533 | 8 | `39a9c9168fd46969` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 1% | 9.2401 | 16 | `deca84b3f1c0a367` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 8.8418 | 16 | `2d5afab67ab568f0` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 7.3989 | 1 | `061f8d400afc81b7` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 7.8191 | 1 | `e32ad292add0e090` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | 7.5782 | 2 | `14a9ae41f3081129` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 100% | 8.0239 | 2 | `a8f159b0283d13f5` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 1% | 7.9473 | 4 | `18d9e30e959932b3` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 9.3836 | 4 | `e8546cd38fbb2dda` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | 8.9516 | 8 | `27972cf38ca35dfa` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 100% | 8.3177 | 8 | `adc1bcb729e3762b` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 1% | 8.3627 | 16 | `5f09799019e3a341` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 9.0314 | 16 | `1e732ab5df58b7ce` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 5.9769 | 1 | `161cdd7be75e1373` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 7.8983 | 1 | `18cffd36cf4def8d` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | 10.9566 | 2 | `125b250e367b1d6a` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 100% | 14.3254 | 2 | `c5aa4c9ca2b2ac17` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 1% | 20.0876 | 4 | `1acb4d7a738023a0` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 26.1702 | 4 | `144f2084e66c2ca0` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | 52.7973 | 8 | `33818faf5a32b56f` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 100% | 51.9686 | 8 | `05c71612c62104b3` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 1% | 73.7653 | 16 | `d845f080e150dc73` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 77.5303 | 16 | `6afcd948c99c5e71` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 1.0292 | false |
| notmat | 1 | 1% | 0.8314 | false |
| mat | 1 | 100% | 1.0486 | false |
| notmat | 1 | 100% | 1.0592 | false |
| mat | 2 | 1% | 1.0245 | false |
| notmat | 2 | 1% | 1.4812 | false |
| mat | 2 | 100% | 0.9652 | false |
| notmat | 2 | 100% | 1.7232 | false |
| mat | 4 | 1% | 0.9420 | false |
| notmat | 4 | 1% | 2.3810 | false |
| mat | 4 | 100% | 1.0833 | false |
| notmat | 4 | 100% | 3.0213 | false |
| mat | 8 | 1% | 1.0529 | false |
| notmat | 8 | 1% | 6.2101 | false |
| mat | 8 | 100% | 1.0591 | false |
| notmat | 8 | 100% | 6.6174 | false |
| mat | 16 | 1% | 0.9050 | false |
| notmat | 16 | 1% | 7.9832 | false |
| mat | 16 | 100% | 1.0214 | false |
| notmat | 16 | 100% | 8.7686 | false |

## Repeated-design position control

| first median ms | last median ms | ratio (last/first) | same plan |
|---|---|---|---|
| 82.355710 | 79.508536 | 0.9654 | false |

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

- client: nr_periods 579, nr_throttled 244, throttled_usec 8254885 | quota 100000 100000
, memory 2147483648

- yb-single: nr_periods 651, nr_throttled 62, throttled_usec 1785998 | quota 200000 100000, memory 3221225472

## Access path per cell

_Engine-reported from each cell's own plan: whether the rank index or a sequential scan was used, rows scanned, PostgreSQL buffers and filter removals, and reader retries of transient conflicts. A cell that scanned the whole table at a selective threshold is a different physical regime, not a faster query._

| axis | status | index scan | seq scan | storage rows scanned | index rows scanned | buffers hit/read | rows removed | retries |
|---|---|---|---|---|---|---|---|---|
| family=bag;stmt=comma_join;threshold=100000 | valid | false | true | 202050 | 0 | 0/0 | 0 | 0 |
| family=bag;stmt=comma_join;threshold=100000;position=last | valid | false | true | 202050 | 0 | 0/0 | 0 | 0 |
| family=bag;stmt=comma_join;threshold=180000 | valid | false | true | 202050 | 0 | 0/0 | 0 | 0 |
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
| family=bag;stmt=comma_join;threshold=100000 | valid | 77.7712 | 82.3557 | 84.1556 | 5 | 12.1 |
| family=bag;stmt=comma_join;threshold=100000;position=last | valid | 77.4420 | 79.5085 | 82.0802 | 5 | 12.6 |
| family=bag;stmt=comma_join;threshold=180000 | valid | 96.4223 | 97.2575 | 100.5431 | 5 | 10.3 |
| family=bag;stmt=comma_join;threshold=20 | valid | 1.3174 | 1.3876 | 1.6580 | 5 | 720.7 |
| family=bag;stmt=comma_join;threshold=200 | valid | 2.6370 | 2.7627 | 2.9974 | 5 | 362.0 |
| family=bag;stmt=comma_join;threshold=2000 | valid | 7.5620 | 8.2027 | 9.1625 | 5 | 121.9 |
| family=bag;stmt=comma_join;threshold=20000 | valid | 55.3388 | 58.0668 | 59.4592 | 5 | 17.2 |
| family=bag;stmt=cte_filtered;threshold=100000 | valid | 77.1169 | 80.0007 | 83.7769 | 5 | 12.5 |
| family=bag;stmt=cte_filtered;threshold=180000 | valid | 96.3800 | 99.4646 | 107.5767 | 5 | 10.1 |
| family=bag;stmt=cte_filtered;threshold=20 | valid | 2.0135 | 2.2764 | 2.4144 | 5 | 439.3 |
| family=bag;stmt=cte_filtered;threshold=200 | valid | 3.4303 | 3.6375 | 3.8571 | 5 | 274.9 |
| family=bag;stmt=cte_filtered;threshold=2000 | valid | 7.7881 | 8.2259 | 8.8903 | 5 | 121.6 |
| family=bag;stmt=cte_filtered;threshold=20000 | valid | 56.1249 | 57.5979 | 58.7430 | 5 | 17.4 |
| family=bag;stmt=cte_filtered_mat;threshold=100000 | valid | 80.5607 | 85.6265 | 87.1776 | 5 | 11.7 |
| family=bag;stmt=cte_filtered_mat;threshold=180000 | valid | 100.2432 | 101.9046 | 114.0429 | 5 | 9.8 |
| family=bag;stmt=cte_filtered_mat;threshold=20 | valid | 1.1999 | 1.5041 | 1.5077 | 5 | 664.9 |
| family=bag;stmt=cte_filtered_mat;threshold=200 | valid | 2.6394 | 2.6817 | 3.0839 | 5 | 372.9 |
| family=bag;stmt=cte_filtered_mat;threshold=2000 | valid | 7.9541 | 8.3444 | 8.7238 | 5 | 119.8 |
| family=bag;stmt=cte_filtered_mat;threshold=20000 | valid | 57.7059 | 60.8615 | 65.3137 | 5 | 16.4 |
| family=bag;stmt=cte_filtered_notmat;threshold=100000 | valid | 79.1824 | 80.0196 | 81.5454 | 5 | 12.5 |
| family=bag;stmt=cte_filtered_notmat;threshold=180000 | valid | 96.2724 | 100.2052 | 102.2770 | 5 | 10.0 |
| family=bag;stmt=cte_filtered_notmat;threshold=20 | valid | 1.7830 | 1.8373 | 2.0162 | 5 | 544.3 |
| family=bag;stmt=cte_filtered_notmat;threshold=200 | valid | 2.9976 | 3.2763 | 3.5699 | 5 | 305.2 |
| family=bag;stmt=cte_filtered_notmat;threshold=2000 | valid | 7.5074 | 7.8776 | 8.9748 | 5 | 126.9 |
| family=bag;stmt=cte_filtered_notmat;threshold=20000 | valid | 55.8241 | 59.7119 | 63.4561 | 5 | 16.7 |
| family=bag;stmt=cte_joined_outer;threshold=100000 | valid | 78.4636 | 80.5928 | 83.1494 | 5 | 12.4 |
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 94.3332 | 98.0680 | 103.2862 | 5 | 10.2 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 1.5587 | 1.7437 | 1.8043 | 5 | 573.5 |
| family=bag;stmt=cte_joined_outer;threshold=200 | valid | 2.6421 | 2.7340 | 2.7700 | 5 | 365.8 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 8.9535 | 9.3601 | 9.4403 | 5 | 106.8 |
| family=bag;stmt=cte_joined_outer;threshold=20000 | valid | 57.8350 | 58.0394 | 59.4758 | 5 | 17.2 |
| family=bag;stmt=cte_joined_outer_mat;threshold=100000 | valid | 82.5597 | 84.5073 | 91.5078 | 5 | 11.8 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 93.5537 | 96.8588 | 101.7060 | 5 | 10.3 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 69.7410 | 77.1213 | 80.5527 | 5 | 13.0 |
| family=bag;stmt=cte_joined_outer_mat;threshold=200 | valid | 70.8216 | 74.4573 | 76.0498 | 5 | 13.4 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 74.7168 | 79.3384 | 80.9385 | 5 | 12.6 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20000 | valid | 76.5788 | 79.4571 | 83.7176 | 5 | 12.6 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=100000 | valid | 77.9596 | 82.3307 | 85.5808 | 5 | 12.1 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 97.0026 | 99.1684 | 100.0482 | 5 | 10.1 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 1.5515 | 1.5944 | 1.9530 | 5 | 627.2 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=200 | valid | 3.2542 | 3.3205 | 3.4711 | 5 | 301.2 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 8.1144 | 8.5146 | 8.8491 | 5 | 117.4 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20000 | valid | 58.2593 | 60.9999 | 95.7857 | 5 | 16.4 |
| family=bag;stmt=filtered_derived;threshold=100000 | valid | 77.6551 | 80.4251 | 86.2047 | 5 | 12.4 |
| family=bag;stmt=filtered_derived;threshold=180000 | valid | 94.8460 | 99.0461 | 102.0143 | 5 | 10.1 |
| family=bag;stmt=filtered_derived;threshold=20 | valid | 1.3728 | 1.5013 | 1.5869 | 5 | 666.1 |
| family=bag;stmt=filtered_derived;threshold=200 | valid | 3.1858 | 3.3203 | 3.4655 | 5 | 301.2 |
| family=bag;stmt=filtered_derived;threshold=2000 | valid | 8.0305 | 9.5214 | 9.9262 | 5 | 105.0 |
| family=bag;stmt=filtered_derived;threshold=20000 | valid | 56.6911 | 59.6671 | 60.3419 | 5 | 16.8 |
| family=bag;stmt=join_on;threshold=100000 | valid | 79.5408 | 83.3810 | 87.1300 | 5 | 12.0 |
| family=bag;stmt=join_on;threshold=180000 | valid | 98.2234 | 99.5377 | 103.2284 | 5 | 10.0 |
| family=bag;stmt=join_on;threshold=20 | valid | 1.8638 | 2.0977 | 2.3493 | 5 | 476.7 |
| family=bag;stmt=join_on;threshold=200 | valid | 2.6667 | 3.0373 | 3.2103 | 5 | 329.2 |
| family=bag;stmt=join_on;threshold=2000 | valid | 8.4897 | 8.8282 | 9.0065 | 5 | 113.3 |
| family=bag;stmt=join_on;threshold=20000 | valid | 56.0077 | 58.4553 | 58.9396 | 5 | 17.1 |
| family=bag;stmt=join_where;threshold=100000 | valid | 78.6620 | 80.4412 | 80.8682 | 5 | 12.4 |
| family=bag;stmt=join_where;threshold=180000 | valid | 94.5899 | 98.7178 | 104.2045 | 5 | 10.1 |
| family=bag;stmt=join_where;threshold=20 | valid | 2.1814 | 2.1913 | 2.5230 | 5 | 456.3 |
| family=bag;stmt=join_where;threshold=200 | valid | 3.1067 | 3.4286 | 3.8611 | 5 | 291.7 |
| family=bag;stmt=join_where;threshold=2000 | valid | 8.0249 | 8.2696 | 9.1952 | 5 | 120.9 |
| family=bag;stmt=join_where;threshold=20000 | valid | 56.3156 | 58.1588 | 60.3905 | 5 | 17.2 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 8.9442 | 9.2401 | 9.4830 | 5 | 108.2 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 8.5836 | 8.8418 | 9.3797 | 5 | 113.1 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 6.9174 | 7.1893 | 7.3799 | 5 | 139.1 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 7.2804 | 7.4566 | 8.1741 | 5 | 134.1 |
| family=reuse;policy=default;refs=2;consumer=1% | valid | 7.1284 | 7.3971 | 8.3727 | 5 | 135.2 |
| family=reuse;policy=default;refs=2;consumer=100% | valid | 8.1445 | 8.3133 | 9.5019 | 5 | 120.3 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 8.2282 | 8.4364 | 9.3958 | 5 | 118.5 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 8.3035 | 8.6619 | 9.2764 | 5 | 115.4 |
| family=reuse;policy=default;refs=8;consumer=1% | valid | 8.2349 | 8.5018 | 8.5536 | 5 | 117.6 |
| family=reuse;policy=default;refs=8;consumer=100% | valid | 7.6705 | 7.8533 | 8.2710 | 5 | 127.3 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 8.1042 | 8.3627 | 8.6760 | 5 | 119.6 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 8.5430 | 9.0314 | 9.9535 | 5 | 110.7 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 7.2036 | 7.3989 | 7.7026 | 5 | 135.2 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 7.5985 | 7.8191 | 8.1626 | 5 | 127.9 |
| family=reuse;policy=mat;refs=2;consumer=1% | valid | 7.1016 | 7.5782 | 7.9403 | 5 | 132.0 |
| family=reuse;policy=mat;refs=2;consumer=100% | valid | 7.8156 | 8.0239 | 9.3303 | 5 | 124.6 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 7.5075 | 7.9473 | 8.2678 | 5 | 125.8 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 8.5674 | 9.3836 | 9.7518 | 5 | 106.6 |
| family=reuse;policy=mat;refs=8;consumer=1% | valid | 8.4280 | 8.9516 | 9.8216 | 5 | 111.7 |
| family=reuse;policy=mat;refs=8;consumer=100% | valid | 7.6097 | 8.3177 | 8.8870 | 5 | 120.2 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 70.6303 | 73.7653 | 93.3399 | 5 | 13.6 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 76.3424 | 77.5303 | 119.1744 | 5 | 12.9 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 5.8892 | 5.9769 | 6.8790 | 5 | 167.3 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 7.6508 | 7.8983 | 8.1182 | 5 | 126.6 |
| family=reuse;policy=notmat;refs=2;consumer=1% | valid | 10.6100 | 10.9566 | 11.4707 | 5 | 91.3 |
| family=reuse;policy=notmat;refs=2;consumer=100% | valid | 14.2903 | 14.3254 | 15.5293 | 5 | 69.8 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 18.6790 | 20.0876 | 22.5198 | 5 | 49.8 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 24.7639 | 26.1702 | 28.2016 | 5 | 38.2 |
| family=reuse;policy=notmat;refs=8;consumer=1% | valid | 49.5949 | 52.7973 | 60.2099 | 5 | 18.9 |
| family=reuse;policy=notmat;refs=8;consumer=100% | valid | 51.6867 | 51.9686 | 54.1966 | 5 | 19.2 |

## Matched blended workload

- phase: `isolated`; target offered rate: 10 updates/s
- offered 0, committed 0, failed 0 over 0.000 s
- offered rate 0.000/s, completed rate 0.000/s
- per key: no keys touched
- ledger: 0 keys reconciled, 0 mismatch

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42248983`, 5 interleaved rounds
  1. family=bag;stmt=comma_join;threshold=100000
  2. family=bag;stmt=cte_filtered_notmat;threshold=20000
  3. family=reuse;policy=mat;refs=4;consumer=100%
  4. family=bag;stmt=cte_filtered_mat;threshold=20000
  5. family=bag;stmt=filtered_derived;threshold=20000
  6. family=bag;stmt=cte_joined_outer_notmat;threshold=200
  7. family=bag;stmt=filtered_derived;threshold=20
  8. family=reuse;policy=mat;refs=1;consumer=1%
  9. family=bag;stmt=join_on;threshold=180000
  10. family=bag;stmt=cte_joined_outer;threshold=2000
  11. family=reuse;policy=default;refs=2;consumer=1%
  12. family=reuse;policy=notmat;refs=16;consumer=1%
  13. family=bag;stmt=join_where;threshold=20000
  14. family=reuse;policy=default;refs=2;consumer=100%
  15. family=bag;stmt=filtered_derived;threshold=180000
  16. family=bag;stmt=cte_filtered;threshold=20
  17. family=bag;stmt=comma_join;threshold=20
  18. family=reuse;policy=notmat;refs=4;consumer=1%
  19. family=bag;stmt=cte_filtered_notmat;threshold=2000
  20. family=reuse;policy=notmat;refs=16;consumer=100%
  21. family=bag;stmt=cte_joined_outer_notmat;threshold=20000
  22. family=bag;stmt=join_where;threshold=200
  23. family=reuse;policy=mat;refs=4;consumer=1%
  24. family=bag;stmt=cte_filtered_notmat;threshold=100000
  25. family=bag;stmt=cte_joined_outer_mat;threshold=20000
  26. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  27. family=bag;stmt=cte_joined_outer;threshold=100000
  28. family=reuse;policy=notmat;refs=2;consumer=100%
  29. family=bag;stmt=comma_join;threshold=20000
  30. family=bag;stmt=join_on;threshold=100000
  31. family=bag;stmt=join_where;threshold=20
  32. family=reuse;policy=notmat;refs=1;consumer=100%
  33. family=bag;stmt=cte_filtered_mat;threshold=2000
  34. family=reuse;policy=mat;refs=2;consumer=1%
  35. family=bag;stmt=cte_filtered_mat;threshold=100000
  36. family=bag;stmt=filtered_derived;threshold=200
  37. family=reuse;policy=mat;refs=16;consumer=100%
  38. family=reuse;policy=notmat;refs=2;consumer=1%
  39. family=bag;stmt=cte_filtered_mat;threshold=180000
  40. family=reuse;policy=default;refs=16;consumer=1%
  41. family=bag;stmt=cte_filtered_mat;threshold=200
  42. family=reuse;policy=mat;refs=1;consumer=100%
  43. family=reuse;policy=mat;refs=8;consumer=100%
  44. family=bag;stmt=join_where;threshold=100000
  45. family=reuse;policy=mat;refs=8;consumer=1%
  46. family=bag;stmt=comma_join;threshold=200
  47. family=bag;stmt=comma_join;threshold=2000
  48. family=reuse;policy=notmat;refs=8;consumer=1%
  49. family=bag;stmt=cte_joined_outer;threshold=20000
  50. family=reuse;policy=default;refs=8;consumer=1%
  51. family=bag;stmt=filtered_derived;threshold=100000
  52. family=bag;stmt=join_where;threshold=180000
  53. family=bag;stmt=cte_filtered;threshold=100000
  54. family=bag;stmt=cte_filtered;threshold=20000
  55. family=reuse;policy=default;refs=4;consumer=100%
  56. family=bag;stmt=cte_joined_outer;threshold=20
  57. family=reuse;policy=default;refs=16;consumer=100%
  58. family=bag;stmt=join_on;threshold=200
  59. family=bag;stmt=cte_filtered_notmat;threshold=180000
  60. family=reuse;policy=default;refs=1;consumer=1%
  61. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  62. family=bag;stmt=join_where;threshold=2000
  63. family=bag;stmt=cte_joined_outer;threshold=180000
  64. family=bag;stmt=cte_filtered_notmat;threshold=200
  65. family=bag;stmt=cte_joined_outer_mat;threshold=100000
  66. family=bag;stmt=cte_joined_outer_mat;threshold=200
  67. family=bag;stmt=filtered_derived;threshold=2000
  68. family=reuse;policy=default;refs=1;consumer=100%
  69. family=bag;stmt=cte_filtered;threshold=180000
  70. family=bag;stmt=cte_joined_outer_mat;threshold=20
  71. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  72. family=reuse;policy=mat;refs=16;consumer=1%
  73. family=reuse;policy=notmat;refs=4;consumer=100%
  74. family=bag;stmt=comma_join;threshold=180000
  75. family=bag;stmt=cte_filtered;threshold=200
  76. family=reuse;policy=notmat;refs=1;consumer=1%
  77. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  78. family=bag;stmt=cte_filtered_notmat;threshold=20
  79. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  80. family=reuse;policy=default;refs=4;consumer=1%
  81. family=bag;stmt=cte_joined_outer;threshold=200
  82. family=reuse;policy=default;refs=8;consumer=100%
  83. family=bag;stmt=join_on;threshold=20000
  84. family=bag;stmt=join_on;threshold=20
  85. family=bag;stmt=cte_filtered_mat;threshold=20
  86. family=reuse;policy=mat;refs=2;consumer=100%
  87. family=bag;stmt=cte_joined_outer_notmat;threshold=100000
  88. family=reuse;policy=notmat;refs=8;consumer=100%
  89. family=bag;stmt=join_on;threshold=2000
  90. family=bag;stmt=cte_filtered;threshold=2000
  91. family=bag;stmt=comma_join;threshold=100000;position=last

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
| family=bag;stmt=comma_join;threshold=100000 | j06_comma_join | valid | `fd197942c4364573` |  |
| family=bag;stmt=comma_join;threshold=100000;position=last | j06_comma_join | valid | `d619a76e89955504` |  |
| family=bag;stmt=comma_join;threshold=180000 | j06_comma_join | valid | `2b755a7d3ebe97ae` |  |
| family=bag;stmt=comma_join;threshold=20 | j06_comma_join | valid | `73d1a0a01ef7632f` |  |
| family=bag;stmt=comma_join;threshold=200 | j06_comma_join | valid | `46339b3bb8ffc745` |  |
| family=bag;stmt=comma_join;threshold=2000 | j06_comma_join | valid | `44d8dee2eb458654` |  |
| family=bag;stmt=comma_join;threshold=20000 | j06_comma_join | valid | `28892841e91500ba` |  |
| family=bag;stmt=cte_filtered;threshold=100000 | j04_cte_filtered | valid | `6ced9f89f0047ce6` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=180000 | j04_cte_filtered | valid | `bcb33a7c2db081ce` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=20 | j04_cte_filtered | valid | `0605e8d4267e5afa` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=200 | j04_cte_filtered | valid | `8c48278c8cc9cdb6` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=2000 | j04_cte_filtered | valid | `3514c49c4ce790f6` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=20000 | j04_cte_filtered | valid | `ba659807589f4421` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_mat;threshold=100000 | j04m_cte_filtered_mat | valid | `48b94d1dfdaf650a` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=180000 | j04m_cte_filtered_mat | valid | `c1296c36599f301f` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=20 | j04m_cte_filtered_mat | valid | `760e8f567f7fa84a` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=200 | j04m_cte_filtered_mat | valid | `f8551775d9aa6004` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=2000 | j04m_cte_filtered_mat | valid | `3a0782ecc88b93aa` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=20000 | j04m_cte_filtered_mat | valid | `3019e1c3300fb589` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_notmat;threshold=100000 | j04n_cte_filtered_notmat | valid | `cd66c3c7a1274b8c` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=180000 | j04n_cte_filtered_notmat | valid | `00f5307aa896958e` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=20 | j04n_cte_filtered_notmat | valid | `f253a05f8693ea53` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=200 | j04n_cte_filtered_notmat | valid | `49f72a23249984b8` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=2000 | j04n_cte_filtered_notmat | valid | `b291fcd836b3115d` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=20000 | j04n_cte_filtered_notmat | valid | `dc6bf0a7c9649482` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=100000 | j05_cte_joined_outer | valid | `86dc57806a138e12` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=180000 | j05_cte_joined_outer | valid | `9e6ecab55eb66a4c` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20 | j05_cte_joined_outer | valid | `6c6b1d9063b61d1f` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=200 | j05_cte_joined_outer | valid | `840f524d3ce00c2f` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=2000 | j05_cte_joined_outer | valid | `6535ff81ec774814` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20000 | j05_cte_joined_outer | valid | `062bbd70b53b9814` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_mat;threshold=100000 | j05m_cte_joined_outer_mat | valid | `5c4c5700209362e1` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | j05m_cte_joined_outer_mat | valid | `66a5a92cff6db2a3` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | j05m_cte_joined_outer_mat | valid | `f02538feca15dd47` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=200 | j05m_cte_joined_outer_mat | valid | `3cf4397c79bc360a` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | j05m_cte_joined_outer_mat | valid | `f8631425ac0df7a1` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20000 | j05m_cte_joined_outer_mat | valid | `7ec05cce0f288c46` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_notmat;threshold=100000 | j05n_cte_joined_outer_notmat | valid | `1f08194cf00ae0c0` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | j05n_cte_joined_outer_notmat | valid | `81fe21139d1bf98f` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | j05n_cte_joined_outer_notmat | valid | `ab880461e2c8c168` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=200 | j05n_cte_joined_outer_notmat | valid | `246bc23b92b9e746` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | j05n_cte_joined_outer_notmat | valid | `061d33535c4854a8` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20000 | j05n_cte_joined_outer_notmat | valid | `88ecc4ffe08f7b7a` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=filtered_derived;threshold=100000 | j03_filtered_derived | valid | `43b3ee1b80f44e6d` |  |
| family=bag;stmt=filtered_derived;threshold=180000 | j03_filtered_derived | valid | `e34b3c4353b5578b` |  |
| family=bag;stmt=filtered_derived;threshold=20 | j03_filtered_derived | valid | `794fd63bf733a84d` |  |
| family=bag;stmt=filtered_derived;threshold=200 | j03_filtered_derived | valid | `d0bc7a943b424570` |  |
| family=bag;stmt=filtered_derived;threshold=2000 | j03_filtered_derived | valid | `baa3ff5998567839` |  |
| family=bag;stmt=filtered_derived;threshold=20000 | j03_filtered_derived | valid | `2eee2e8f9567d5b2` |  |
| family=bag;stmt=join_on;threshold=100000 | j02_join_on | valid | `050a68b9ab719898` |  |
| family=bag;stmt=join_on;threshold=180000 | j02_join_on | valid | `34d4e7f4a7d25e46` |  |
| family=bag;stmt=join_on;threshold=20 | j02_join_on | valid | `6f2a3512f223841a` |  |
| family=bag;stmt=join_on;threshold=200 | j02_join_on | valid | `ef79f544a0242eb4` |  |
| family=bag;stmt=join_on;threshold=2000 | j02_join_on | valid | `3d6a57a58a342425` |  |
| family=bag;stmt=join_on;threshold=20000 | j02_join_on | valid | `18022eeeb297e798` |  |
| family=bag;stmt=join_where;threshold=100000 | j01_join_where | valid | `1ff9dfdac4a5fcbd` |  |
| family=bag;stmt=join_where;threshold=180000 | j01_join_where | valid | `9a0079e253636f80` |  |
| family=bag;stmt=join_where;threshold=20 | j01_join_where | valid | `5511436f92b8f8a6` |  |
| family=bag;stmt=join_where;threshold=200 | j01_join_where | valid | `5840c190b7ba3186` |  |
| family=bag;stmt=join_where;threshold=2000 | j01_join_where | valid | `d1a4bd89a9e494f4` |  |
| family=bag;stmt=join_where;threshold=20000 | j01_join_where | valid | `c7571197fb4c45cf` |  |
| family=reuse;policy=default;refs=16;consumer=1% | r_default_n16_c1 | valid | `deca84b3f1c0a367` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=16;consumer=100% | r_default_n16_c100 | valid | `2d5afab67ab568f0` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=1;consumer=1% | r_default_n1_c1 | valid | `bab144ad7db45979` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=1;consumer=100% | r_default_n1_c100 | valid | `15409dd542ec76b4` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=2;consumer=1% | r_default_n2_c1 | valid | `e7bd0e71013b7d44` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=2;consumer=100% | r_default_n2_c100 | valid | `ae6de606dc096a49` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=1% | r_default_n4_c1 | valid | `a1705411ecba7790` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=100% | r_default_n4_c100 | valid | `0697b7d47cea9c02` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=8;consumer=1% | r_default_n8_c1 | valid | `8e53f03785f056d9` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=8;consumer=100% | r_default_n8_c100 | valid | `39a9c9168fd46969` | policy=default; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=1% | r_mat_n16_c1 | valid | `5f09799019e3a341` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=100% | r_mat_n16_c100 | valid | `1e732ab5df58b7ce` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=1% | r_mat_n1_c1 | valid | `061f8d400afc81b7` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=100% | r_mat_n1_c100 | valid | `e32ad292add0e090` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=2;consumer=1% | r_mat_n2_c1 | valid | `14a9ae41f3081129` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=2;consumer=100% | r_mat_n2_c100 | valid | `a8f159b0283d13f5` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=1% | r_mat_n4_c1 | valid | `18d9e30e959932b3` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=100% | r_mat_n4_c100 | valid | `e8546cd38fbb2dda` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=8;consumer=1% | r_mat_n8_c1 | valid | `27972cf38ca35dfa` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=8;consumer=100% | r_mat_n8_c100 | valid | `adc1bcb729e3762b` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=notmat;refs=16;consumer=1% | r_notmat_n16_c1 | valid | `d845f080e150dc73` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `6afcd948c99c5e71` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `161cdd7be75e1373` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `18cffd36cf4def8d` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=2;consumer=1% | r_notmat_n2_c1 | valid | `125b250e367b1d6a` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=2;consumer=100% | r_notmat_n2_c100 | valid | `c5aa4c9ca2b2ac17` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `1acb4d7a738023a0` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `144f2084e66c2ca0` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=8;consumer=1% | r_notmat_n8_c1 | valid | `33818faf5a32b56f` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=8;consumer=100% | r_notmat_n8_c100 | valid | `05c71612c62104b3` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, DIST, DEBUG, FORMAT TEXT` supported on yugabyte
- one host: this run used node budget 2 CPU 3g and client budget 1 CPU 2g, so budget comparisons are separate classes here, not independent hardware
- isolated phase: no concurrent writers in this run; the blended read/write phase is a separate labelled run
- physical rank-index arm yb-range-rank; an engine whose planner ignores the index is a different physical design, not a faster one
- YugabyteDB DIST/DEBUG storage counters (rows and index rows scanned) are parsed per cell; RPC-level counters are not aggregated
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
