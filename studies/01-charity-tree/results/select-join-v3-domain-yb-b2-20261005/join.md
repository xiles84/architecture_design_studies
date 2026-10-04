# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 91 cells, 91 valid, engine `yugabyte`, topology `single`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-v3-domain-yb-b2-20261005`
- run tag: `run/study-01/select-join-v3-domain-yb-b2-20261005`
- environment: `host-zenbook-ux5406sa`
- repo commit: `0f9d09212c6d52a80d8437d3b5a234f2f711ed36`
- engine version: `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit`
- database engine image: `docker.io/yugabytedb/yugabyte:2025.2.6.0-b111` @ `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`
- benchmark client image: `localhost/charitybench-select:select-join-v3-domain-yb-b2-20261005` @ `8436e36bb66caa8178559d54771b88c063e26848d4ef1d1a56eb524b7f0ecca6`
- inputs digest: `bed4e995fa3e86670108dc3fbee8833be3ff294dd8505ae8389bded608a9efdc`
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
| yb-single | 933 | 2 |

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
| 20 (0.01%) | 2.1983 | 1.7914 | 1.9585 | 2.0526 | 1.7043 | 2.3026 | 1.6921 | 82.6930 | 1.7251 | 1.7389 |
| 200 (0.1%) | 3.0922 | 3.6128 | 2.8377 | 3.7858 | 3.6037 | 3.5360 | 2.7540 | 93.5310 | 3.0461 | 2.8573 |
| 2000 (1%) | 9.2090 | 9.6649 | 9.6381 | 8.9813 | 9.5452 | 9.1439 | 9.5854 | 83.9301 | 8.5840 | 9.2026 |
| 20000 (10%) | 62.3389 | 61.9314 | 66.5514 | 64.2622 | 66.3114 | 68.5064 | 63.1333 | 86.4136 | 62.6333 | 63.7830 |
| 100000 (50%) | 87.8306 | 87.1029 | 88.3963 | 86.5522 | 88.2492 | 86.3569 | 86.6261 | 95.2325 | 87.6329 | 86.0619 |
| 180000 (90%) | 116.0717 | 107.6085 | 107.4201 | 106.3508 | 111.8531 | 103.4570 | 105.5358 | 103.8691 | 101.5432 | 108.0347 |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | join_on vs join_where | 2.1983 | 1.7914 | 0.8149 | false |
| 20 | filtered_derived vs join_where | 2.1983 | 1.9585 | 0.8909 | false |
| 20 | cte_filtered vs join_where | 2.1983 | 2.0526 | 0.9337 | false |
| 20 | cte_filtered_mat vs join_where | 2.1983 | 1.7043 | 0.7753 | false |
| 20 | cte_filtered_notmat vs join_where | 2.1983 | 2.3026 | 1.0474 | false |
| 20 | cte_joined_outer vs join_where | 2.1983 | 1.6921 | 0.7697 | false |
| 20 | cte_joined_outer_mat vs join_where | 2.1983 | 82.6930 | 37.6166 | false |
| 20 | cte_joined_outer_notmat vs join_where | 2.1983 | 1.7251 | 0.7847 | false |
| 20 | comma_join vs join_where | 2.1983 | 1.7389 | 0.7910 | false |
| 200 | join_on vs join_where | 3.0922 | 3.6128 | 1.1684 | false |
| 200 | filtered_derived vs join_where | 3.0922 | 2.8377 | 0.9177 | false |
| 200 | cte_filtered vs join_where | 3.0922 | 3.7858 | 1.2243 | false |
| 200 | cte_filtered_mat vs join_where | 3.0922 | 3.6037 | 1.1654 | false |
| 200 | cte_filtered_notmat vs join_where | 3.0922 | 3.5360 | 1.1435 | false |
| 200 | cte_joined_outer vs join_where | 3.0922 | 2.7540 | 0.8906 | false |
| 200 | cte_joined_outer_mat vs join_where | 3.0922 | 93.5310 | 30.2476 | false |
| 200 | cte_joined_outer_notmat vs join_where | 3.0922 | 3.0461 | 0.9851 | false |
| 200 | comma_join vs join_where | 3.0922 | 2.8573 | 0.9240 | false |
| 2000 | join_on vs join_where | 9.2090 | 9.6649 | 1.0495 | false |
| 2000 | filtered_derived vs join_where | 9.2090 | 9.6381 | 1.0466 | false |
| 2000 | cte_filtered vs join_where | 9.2090 | 8.9813 | 0.9753 | false |
| 2000 | cte_filtered_mat vs join_where | 9.2090 | 9.5452 | 1.0365 | false |
| 2000 | cte_filtered_notmat vs join_where | 9.2090 | 9.1439 | 0.9929 | false |
| 2000 | cte_joined_outer vs join_where | 9.2090 | 9.5854 | 1.0409 | false |
| 2000 | cte_joined_outer_mat vs join_where | 9.2090 | 83.9301 | 9.1139 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 9.2090 | 8.5840 | 0.9321 | false |
| 2000 | comma_join vs join_where | 9.2090 | 9.2026 | 0.9993 | false |
| 20000 | join_on vs join_where | 62.3389 | 61.9314 | 0.9935 | false |
| 20000 | filtered_derived vs join_where | 62.3389 | 66.5514 | 1.0676 | false |
| 20000 | cte_filtered vs join_where | 62.3389 | 64.2622 | 1.0309 | false |
| 20000 | cte_filtered_mat vs join_where | 62.3389 | 66.3114 | 1.0637 | false |
| 20000 | cte_filtered_notmat vs join_where | 62.3389 | 68.5064 | 1.0989 | false |
| 20000 | cte_joined_outer vs join_where | 62.3389 | 63.1333 | 1.0127 | false |
| 20000 | cte_joined_outer_mat vs join_where | 62.3389 | 86.4136 | 1.3862 | false |
| 20000 | cte_joined_outer_notmat vs join_where | 62.3389 | 62.6333 | 1.0047 | false |
| 20000 | comma_join vs join_where | 62.3389 | 63.7830 | 1.0232 | false |
| 100000 | join_on vs join_where | 87.8306 | 87.1029 | 0.9917 | false |
| 100000 | filtered_derived vs join_where | 87.8306 | 88.3963 | 1.0064 | false |
| 100000 | cte_filtered vs join_where | 87.8306 | 86.5522 | 0.9854 | false |
| 100000 | cte_filtered_mat vs join_where | 87.8306 | 88.2492 | 1.0048 | false |
| 100000 | cte_filtered_notmat vs join_where | 87.8306 | 86.3569 | 0.9832 | false |
| 100000 | cte_joined_outer vs join_where | 87.8306 | 86.6261 | 0.9863 | false |
| 100000 | cte_joined_outer_mat vs join_where | 87.8306 | 95.2325 | 1.0843 | false |
| 100000 | cte_joined_outer_notmat vs join_where | 87.8306 | 87.6329 | 0.9977 | false |
| 100000 | comma_join vs join_where | 87.8306 | 86.0619 | 0.9799 | false |
| 180000 | join_on vs join_where | 116.0717 | 107.6085 | 0.9271 | false |
| 180000 | filtered_derived vs join_where | 116.0717 | 107.4201 | 0.9255 | false |
| 180000 | cte_filtered vs join_where | 116.0717 | 106.3508 | 0.9163 | false |
| 180000 | cte_filtered_mat vs join_where | 116.0717 | 111.8531 | 0.9637 | false |
| 180000 | cte_filtered_notmat vs join_where | 116.0717 | 103.4570 | 0.8913 | false |
| 180000 | cte_joined_outer vs join_where | 116.0717 | 105.5358 | 0.9092 | false |
| 180000 | cte_joined_outer_mat vs join_where | 116.0717 | 103.8691 | 0.8949 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 116.0717 | 101.5432 | 0.8748 | false |
| 180000 | comma_join vs join_where | 116.0717 | 108.0347 | 0.9308 | false |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_filtered_mat vs cte_filtered | 2.0526 | 1.7043 | 0.8303 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_filtered_notmat vs cte_filtered | 2.0526 | 2.3026 | 1.1218 | false | policy=not_materialized; observed_cte_scan=false |
| 20 | cte_joined_outer_mat vs cte_joined_outer | 1.6921 | 82.6930 | 48.8713 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 1.6921 | 1.7251 | 1.0195 | false | policy=not_materialized; observed_cte_scan=false |
| 200 | cte_filtered_mat vs cte_filtered | 3.7858 | 3.6037 | 0.9519 | false | policy=materialized; observed_cte_scan=true |
| 200 | cte_filtered_notmat vs cte_filtered | 3.7858 | 3.5360 | 0.9340 | false | policy=not_materialized; observed_cte_scan=false |
| 200 | cte_joined_outer_mat vs cte_joined_outer | 2.7540 | 93.5310 | 33.9617 | false | policy=materialized; observed_cte_scan=true |
| 200 | cte_joined_outer_notmat vs cte_joined_outer | 2.7540 | 3.0461 | 1.1060 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_filtered_mat vs cte_filtered | 8.9813 | 9.5452 | 1.0628 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_filtered_notmat vs cte_filtered | 8.9813 | 9.1439 | 1.0181 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 9.5854 | 83.9301 | 8.7561 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 9.5854 | 8.5840 | 0.8955 | false | policy=not_materialized; observed_cte_scan=false |
| 20000 | cte_filtered_mat vs cte_filtered | 64.2622 | 66.3114 | 1.0319 | false | policy=materialized; observed_cte_scan=true |
| 20000 | cte_filtered_notmat vs cte_filtered | 64.2622 | 68.5064 | 1.0660 | false | policy=not_materialized; observed_cte_scan=false |
| 20000 | cte_joined_outer_mat vs cte_joined_outer | 63.1333 | 86.4136 | 1.3687 | false | policy=materialized; observed_cte_scan=true |
| 20000 | cte_joined_outer_notmat vs cte_joined_outer | 63.1333 | 62.6333 | 0.9921 | false | policy=not_materialized; observed_cte_scan=false |
| 100000 | cte_filtered_mat vs cte_filtered | 86.5522 | 88.2492 | 1.0196 | false | policy=materialized; observed_cte_scan=true |
| 100000 | cte_filtered_notmat vs cte_filtered | 86.5522 | 86.3569 | 0.9977 | false | policy=not_materialized; observed_cte_scan=false |
| 100000 | cte_joined_outer_mat vs cte_joined_outer | 86.6261 | 95.2325 | 1.0994 | false | policy=materialized; observed_cte_scan=true |
| 100000 | cte_joined_outer_notmat vs cte_joined_outer | 86.6261 | 87.6329 | 1.0116 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_filtered_mat vs cte_filtered | 106.3508 | 111.8531 | 1.0517 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_filtered_notmat vs cte_filtered | 106.3508 | 103.4570 | 0.9728 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 105.5358 | 103.8691 | 0.9842 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 105.5358 | 101.5432 | 0.9622 | false | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 6.3009 | 1 | `8d398a2cee5943c3` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 8.7621 | 1 | `f04625e62d5faada` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | 9.3312 | 2 | `e744104fdb4622d5` | policy=default; observed_cte_scan=true | valid |
| default | 2 | 100% | 8.2106 | 2 | `32a2b792810dc2a8` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 1% | 8.5648 | 4 | `81fa286e97be8d25` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 8.3781 | 4 | `9c9aa704da2fb1b4` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | 9.3836 | 8 | `1c66778d54405c53` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 100% | 9.4281 | 8 | `6c2b01729f08759e` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 1% | 8.7374 | 16 | `6c20417679f5fbfd` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 9.8167 | 16 | `5f5d10ac883293de` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 9.2904 | 1 | `6689b2d0fc8606f8` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 8.9986 | 1 | `b26644247da5b41b` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | 7.7799 | 2 | `8b6750db6eaaa211` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 100% | 8.4157 | 2 | `26ed7e06a8edf4db` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 1% | 8.8143 | 4 | `5117e7368dc54c32` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 9.7157 | 4 | `5776bbb330b77a18` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | 8.5764 | 8 | `9517f0c04a534fd9` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 100% | 9.5841 | 8 | `0c41888a3b10e411` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 1% | 10.4571 | 16 | `def8116385a54d81` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 10.9405 | 16 | `7a87518a258f56ad` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 7.2889 | 1 | `e52181004a68fc5d` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 8.7757 | 1 | `40f1af0c733c861f` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | 13.2200 | 2 | `c614e4eae8838fcf` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 100% | 15.2284 | 2 | `8d48a37065cf22ac` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 1% | 23.9589 | 4 | `0c0b0c14bddb1fe0` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 29.0846 | 4 | `7d07c3c9c394d739` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | 58.4491 | 8 | `8f5e74f1794de76b` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 100% | 59.8194 | 8 | `af027cc73ee78117` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 1% | 74.4537 | 16 | `e9cf25e328e4e9fa` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 114.2530 | 16 | `ce7d445a99b77c8f` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 1.4745 | false |
| notmat | 1 | 1% | 1.1568 | false |
| mat | 1 | 100% | 1.0270 | false |
| notmat | 1 | 100% | 1.0016 | false |
| mat | 2 | 1% | 0.8338 | false |
| notmat | 2 | 1% | 1.4167 | false |
| mat | 2 | 100% | 1.0250 | false |
| notmat | 2 | 100% | 1.8547 | false |
| mat | 4 | 1% | 1.0291 | false |
| notmat | 4 | 1% | 2.7974 | false |
| mat | 4 | 100% | 1.1597 | false |
| notmat | 4 | 100% | 3.4715 | false |
| mat | 8 | 1% | 0.9140 | false |
| notmat | 8 | 1% | 6.2288 | false |
| mat | 8 | 100% | 1.0165 | false |
| notmat | 8 | 100% | 6.3448 | false |
| mat | 16 | 1% | 1.1968 | false |
| notmat | 16 | 1% | 8.5213 | false |
| mat | 16 | 100% | 1.1145 | false |
| notmat | 16 | 100% | 11.6387 | false |

## Repeated-design position control

| first median ms | last median ms | ratio (last/first) | same plan |
|---|---|---|---|
| 2.198310 | 1.637466 | 0.7449 | false |

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

- client: nr_periods 630, nr_throttled 246, throttled_usec 10237110 | quota 100000 100000
, memory 2147483648

- yb-single: nr_periods 694, nr_throttled 69, throttled_usec 2758335 | quota 200000 100000, memory 3221225472

## Access path per cell

_Engine-reported from each cell's own plan: whether the rank index or a sequential scan was used, rows scanned, PostgreSQL buffers and filter removals, and reader retries of transient conflicts. A cell that scanned the whole table at a selective threshold is a different physical regime, not a faster query._

| axis | status | index scan | seq scan | storage rows scanned | index rows scanned | buffers hit/read | rows removed | retries |
|---|---|---|---|---|---|---|---|---|
| family=bag;stmt=comma_join;threshold=100000 | valid | false | true | 202050 | 0 | 0/0 | 0 | 0 |
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
| family=bag;stmt=join_where;threshold=20;position=last | valid | true | true | 110 | 20 | 0/0 | 0 | 0 |
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
| family=bag;stmt=comma_join;threshold=100000 | valid | 84.7477 | 86.0619 | 87.7577 | 5 | 11.6 |
| family=bag;stmt=comma_join;threshold=180000 | valid | 105.0384 | 108.0347 | 110.7538 | 5 | 9.3 |
| family=bag;stmt=comma_join;threshold=20 | valid | 1.7024 | 1.7389 | 1.8067 | 5 | 575.1 |
| family=bag;stmt=comma_join;threshold=200 | valid | 2.5278 | 2.8573 | 24.1305 | 5 | 350.0 |
| family=bag;stmt=comma_join;threshold=2000 | valid | 8.7996 | 9.2026 | 9.5481 | 5 | 108.7 |
| family=bag;stmt=comma_join;threshold=20000 | valid | 63.5349 | 63.7830 | 64.8591 | 5 | 15.7 |
| family=bag;stmt=cte_filtered;threshold=100000 | valid | 84.5000 | 86.5522 | 95.5378 | 5 | 11.6 |
| family=bag;stmt=cte_filtered;threshold=180000 | valid | 103.2464 | 106.3508 | 108.4290 | 5 | 9.4 |
| family=bag;stmt=cte_filtered;threshold=20 | valid | 1.6951 | 2.0526 | 2.2231 | 5 | 487.2 |
| family=bag;stmt=cte_filtered;threshold=200 | valid | 3.5805 | 3.7858 | 4.1215 | 5 | 264.1 |
| family=bag;stmt=cte_filtered;threshold=2000 | valid | 8.7352 | 8.9813 | 9.1242 | 5 | 111.3 |
| family=bag;stmt=cte_filtered;threshold=20000 | valid | 63.0790 | 64.2622 | 74.1387 | 5 | 15.6 |
| family=bag;stmt=cte_filtered_mat;threshold=100000 | valid | 87.2752 | 88.2492 | 94.2738 | 5 | 11.3 |
| family=bag;stmt=cte_filtered_mat;threshold=180000 | valid | 106.8806 | 111.8531 | 114.7182 | 5 | 8.9 |
| family=bag;stmt=cte_filtered_mat;threshold=20 | valid | 1.5031 | 1.7043 | 1.9242 | 5 | 586.7 |
| family=bag;stmt=cte_filtered_mat;threshold=200 | valid | 3.3310 | 3.6037 | 3.9597 | 5 | 277.5 |
| family=bag;stmt=cte_filtered_mat;threshold=2000 | valid | 9.0448 | 9.5452 | 11.7044 | 5 | 104.8 |
| family=bag;stmt=cte_filtered_mat;threshold=20000 | valid | 64.2494 | 66.3114 | 68.0161 | 5 | 15.1 |
| family=bag;stmt=cte_filtered_notmat;threshold=100000 | valid | 83.9545 | 86.3569 | 92.4289 | 5 | 11.6 |
| family=bag;stmt=cte_filtered_notmat;threshold=180000 | valid | 102.2164 | 103.4570 | 107.6835 | 5 | 9.7 |
| family=bag;stmt=cte_filtered_notmat;threshold=20 | valid | 1.9377 | 2.3026 | 2.4237 | 5 | 434.3 |
| family=bag;stmt=cte_filtered_notmat;threshold=200 | valid | 3.4030 | 3.5360 | 3.9508 | 5 | 282.8 |
| family=bag;stmt=cte_filtered_notmat;threshold=2000 | valid | 8.4660 | 9.1439 | 9.8009 | 5 | 109.4 |
| family=bag;stmt=cte_filtered_notmat;threshold=20000 | valid | 61.2387 | 68.5064 | 69.8656 | 5 | 14.6 |
| family=bag;stmt=cte_joined_outer;threshold=100000 | valid | 85.6300 | 86.6261 | 90.8822 | 5 | 11.5 |
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 101.8228 | 105.5358 | 115.5715 | 5 | 9.5 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 1.5587 | 1.6921 | 1.9873 | 5 | 591.0 |
| family=bag;stmt=cte_joined_outer;threshold=200 | valid | 2.5596 | 2.7540 | 3.0559 | 5 | 363.1 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 9.3211 | 9.5854 | 10.3940 | 5 | 104.3 |
| family=bag;stmt=cte_joined_outer;threshold=20000 | valid | 60.1634 | 63.1333 | 63.7636 | 5 | 15.8 |
| family=bag;stmt=cte_joined_outer_mat;threshold=100000 | valid | 85.5570 | 95.2325 | 98.1320 | 5 | 10.5 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 98.0327 | 103.8691 | 108.5873 | 5 | 9.6 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 78.5456 | 82.6930 | 87.5176 | 5 | 12.1 |
| family=bag;stmt=cte_joined_outer_mat;threshold=200 | valid | 85.0712 | 93.5310 | 97.5931 | 5 | 10.7 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 78.0425 | 83.9301 | 88.1169 | 5 | 11.9 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20000 | valid | 81.3899 | 86.4136 | 88.8089 | 5 | 11.6 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=100000 | valid | 86.6054 | 87.6329 | 92.5211 | 5 | 11.4 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 100.4521 | 101.5432 | 105.1164 | 5 | 9.8 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 1.5915 | 1.7251 | 1.9381 | 5 | 579.7 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=200 | valid | 2.6858 | 3.0461 | 3.0854 | 5 | 328.3 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 8.0914 | 8.5840 | 9.2056 | 5 | 116.5 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20000 | valid | 61.0585 | 62.6333 | 67.9157 | 5 | 16.0 |
| family=bag;stmt=filtered_derived;threshold=100000 | valid | 84.5187 | 88.3963 | 91.8505 | 5 | 11.3 |
| family=bag;stmt=filtered_derived;threshold=180000 | valid | 100.8397 | 107.4201 | 109.0011 | 5 | 9.3 |
| family=bag;stmt=filtered_derived;threshold=20 | valid | 1.8196 | 1.9585 | 2.0394 | 5 | 510.6 |
| family=bag;stmt=filtered_derived;threshold=200 | valid | 2.8185 | 2.8377 | 2.9123 | 5 | 352.4 |
| family=bag;stmt=filtered_derived;threshold=2000 | valid | 9.3658 | 9.6381 | 10.1020 | 5 | 103.8 |
| family=bag;stmt=filtered_derived;threshold=20000 | valid | 62.7842 | 66.5514 | 99.5966 | 5 | 15.0 |
| family=bag;stmt=join_on;threshold=100000 | valid | 82.7084 | 87.1029 | 90.1046 | 5 | 11.5 |
| family=bag;stmt=join_on;threshold=180000 | valid | 102.7970 | 107.6085 | 120.8752 | 5 | 9.3 |
| family=bag;stmt=join_on;threshold=20 | valid | 1.7417 | 1.7914 | 1.9559 | 5 | 558.2 |
| family=bag;stmt=join_on;threshold=200 | valid | 3.1495 | 3.6128 | 3.8976 | 5 | 276.8 |
| family=bag;stmt=join_on;threshold=2000 | valid | 8.6778 | 9.6649 | 9.7211 | 5 | 103.5 |
| family=bag;stmt=join_on;threshold=20000 | valid | 59.8556 | 61.9314 | 64.7971 | 5 | 16.1 |
| family=bag;stmt=join_where;threshold=100000 | valid | 84.5643 | 87.8306 | 90.8759 | 5 | 11.4 |
| family=bag;stmt=join_where;threshold=180000 | valid | 107.4098 | 116.0717 | 137.1786 | 5 | 8.6 |
| family=bag;stmt=join_where;threshold=20 | valid | 2.0284 | 2.1983 | 2.7514 | 5 | 454.9 |
| family=bag;stmt=join_where;threshold=200 | valid | 2.8339 | 3.0922 | 3.3734 | 5 | 323.4 |
| family=bag;stmt=join_where;threshold=2000 | valid | 8.6950 | 9.2090 | 10.3576 | 5 | 108.6 |
| family=bag;stmt=join_where;threshold=20000 | valid | 61.9013 | 62.3389 | 63.4288 | 5 | 16.0 |
| family=bag;stmt=join_where;threshold=20;position=last | valid | 1.4234 | 1.6375 | 2.6010 | 5 | 610.7 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 8.5703 | 8.7374 | 9.5823 | 5 | 114.5 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 9.4580 | 9.8167 | 10.3109 | 5 | 101.9 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 6.0496 | 6.3009 | 6.7405 | 5 | 158.7 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 8.3624 | 8.7621 | 9.0158 | 5 | 114.1 |
| family=reuse;policy=default;refs=2;consumer=1% | valid | 8.9579 | 9.3312 | 9.6931 | 5 | 107.2 |
| family=reuse;policy=default;refs=2;consumer=100% | valid | 7.5128 | 8.2106 | 8.9583 | 5 | 121.8 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 8.3218 | 8.5648 | 9.8098 | 5 | 116.8 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 8.0133 | 8.3781 | 8.9301 | 5 | 119.4 |
| family=reuse;policy=default;refs=8;consumer=1% | valid | 9.0827 | 9.3836 | 9.4572 | 5 | 106.6 |
| family=reuse;policy=default;refs=8;consumer=100% | valid | 8.8955 | 9.4281 | 9.8901 | 5 | 106.1 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 9.4500 | 10.4571 | 10.6281 | 5 | 95.6 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 10.3454 | 10.9405 | 11.6848 | 5 | 91.4 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 8.8464 | 9.2904 | 9.9840 | 5 | 107.6 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 8.8060 | 8.9986 | 9.2676 | 5 | 111.1 |
| family=reuse;policy=mat;refs=2;consumer=1% | valid | 7.5712 | 7.7799 | 8.3529 | 5 | 128.5 |
| family=reuse;policy=mat;refs=2;consumer=100% | valid | 8.3199 | 8.4157 | 8.7719 | 5 | 118.8 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 7.7079 | 8.8143 | 9.4670 | 5 | 113.5 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 9.4228 | 9.7157 | 10.0282 | 5 | 102.9 |
| family=reuse;policy=mat;refs=8;consumer=1% | valid | 8.0261 | 8.5764 | 8.7335 | 5 | 116.6 |
| family=reuse;policy=mat;refs=8;consumer=100% | valid | 9.1374 | 9.5841 | 10.6117 | 5 | 104.3 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 73.9143 | 74.4537 | 117.6855 | 5 | 13.4 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 111.0390 | 114.2530 | 126.5011 | 5 | 8.8 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 7.0992 | 7.2889 | 7.5394 | 5 | 137.2 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 8.5180 | 8.7757 | 9.5174 | 5 | 114.0 |
| family=reuse;policy=notmat;refs=2;consumer=1% | valid | 11.9367 | 13.2200 | 13.8253 | 5 | 75.6 |
| family=reuse;policy=notmat;refs=2;consumer=100% | valid | 14.9387 | 15.2284 | 16.5043 | 5 | 65.7 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 20.8532 | 23.9589 | 29.3448 | 5 | 41.7 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 27.6240 | 29.0846 | 31.8812 | 5 | 34.4 |
| family=reuse;policy=notmat;refs=8;consumer=1% | valid | 57.5858 | 58.4491 | 63.6480 | 5 | 17.1 |
| family=reuse;policy=notmat;refs=8;consumer=100% | valid | 57.7249 | 59.8194 | 73.2150 | 5 | 16.7 |

## Matched blended workload

- phase: `blended`; target offered rate: 10 updates/s
- offered 539, committed 539, failed 0 over 53.894 s
- offered rate 10.001/s, completed rate 10.001/s
- writer operations per endpoint: [539]
- reader operations overlapping a committed write: 226 of 637
- reader retries of transient conflicts: 0
- per key: 539 keys, min 1, max 1 per key
- ledger: 539 keys reconciled, 0 mismatch
- update latency: avg 3.755 ms, min 1.453, max 197.722 over 539

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42077781`, 5 interleaved rounds
  1. family=bag;stmt=join_where;threshold=20
  2. family=bag;stmt=join_where;threshold=180000
  3. family=reuse;policy=mat;refs=16;consumer=100%
  4. family=bag;stmt=join_where;threshold=2000
  5. family=reuse;policy=default;refs=16;consumer=100%
  6. family=bag;stmt=comma_join;threshold=20
  7. family=bag;stmt=join_where;threshold=100000
  8. family=reuse;policy=default;refs=2;consumer=1%
  9. family=reuse;policy=notmat;refs=4;consumer=100%
  10. family=reuse;policy=default;refs=2;consumer=100%
  11. family=reuse;policy=mat;refs=4;consumer=1%
  12. family=bag;stmt=cte_joined_outer_mat;threshold=100000
  13. family=bag;stmt=cte_filtered_mat;threshold=180000
  14. family=reuse;policy=mat;refs=1;consumer=1%
  15. family=bag;stmt=filtered_derived;threshold=20
  16. family=bag;stmt=comma_join;threshold=100000
  17. family=bag;stmt=comma_join;threshold=20000
  18. family=bag;stmt=cte_joined_outer;threshold=100000
  19. family=bag;stmt=cte_joined_outer;threshold=2000
  20. family=reuse;policy=default;refs=1;consumer=1%
  21. family=bag;stmt=join_on;threshold=20000
  22. family=reuse;policy=mat;refs=4;consumer=100%
  23. family=bag;stmt=cte_joined_outer;threshold=20000
  24. family=bag;stmt=cte_filtered_mat;threshold=20000
  25. family=reuse;policy=default;refs=1;consumer=100%
  26. family=bag;stmt=cte_filtered_notmat;threshold=180000
  27. family=bag;stmt=cte_joined_outer_notmat;threshold=20000
  28. family=bag;stmt=join_on;threshold=200
  29. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  30. family=reuse;policy=notmat;refs=1;consumer=1%
  31. family=reuse;policy=mat;refs=2;consumer=1%
  32. family=bag;stmt=join_where;threshold=20000
  33. family=reuse;policy=notmat;refs=2;consumer=1%
  34. family=bag;stmt=cte_joined_outer_mat;threshold=20000
  35. family=reuse;policy=notmat;refs=1;consumer=100%
  36. family=bag;stmt=join_where;threshold=200
  37. family=bag;stmt=cte_filtered_mat;threshold=20
  38. family=bag;stmt=cte_filtered_mat;threshold=100000
  39. family=reuse;policy=mat;refs=1;consumer=100%
  40. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  41. family=bag;stmt=filtered_derived;threshold=180000
  42. family=reuse;policy=notmat;refs=8;consumer=1%
  43. family=bag;stmt=cte_joined_outer_mat;threshold=200
  44. family=bag;stmt=join_on;threshold=180000
  45. family=bag;stmt=cte_filtered;threshold=20000
  46. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  47. family=reuse;policy=notmat;refs=16;consumer=100%
  48. family=bag;stmt=filtered_derived;threshold=100000
  49. family=reuse;policy=mat;refs=16;consumer=1%
  50. family=bag;stmt=comma_join;threshold=180000
  51. family=bag;stmt=cte_filtered;threshold=180000
  52. family=bag;stmt=filtered_derived;threshold=2000
  53. family=bag;stmt=join_on;threshold=20
  54. family=reuse;policy=notmat;refs=2;consumer=100%
  55. family=reuse;policy=default;refs=16;consumer=1%
  56. family=bag;stmt=cte_filtered_notmat;threshold=100000
  57. family=bag;stmt=cte_filtered_notmat;threshold=20000
  58. family=bag;stmt=cte_joined_outer;threshold=180000
  59. family=bag;stmt=cte_filtered;threshold=200
  60. family=reuse;policy=mat;refs=2;consumer=100%
  61. family=reuse;policy=default;refs=4;consumer=100%
  62. family=reuse;policy=notmat;refs=4;consumer=1%
  63. family=bag;stmt=join_on;threshold=100000
  64. family=bag;stmt=cte_filtered_mat;threshold=200
  65. family=bag;stmt=cte_joined_outer;threshold=200
  66. family=reuse;policy=mat;refs=8;consumer=100%
  67. family=bag;stmt=cte_filtered;threshold=100000
  68. family=reuse;policy=default;refs=8;consumer=1%
  69. family=bag;stmt=cte_joined_outer_notmat;threshold=200
  70. family=bag;stmt=comma_join;threshold=200
  71. family=bag;stmt=cte_joined_outer;threshold=20
  72. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  73. family=bag;stmt=cte_filtered_mat;threshold=2000
  74. family=reuse;policy=default;refs=8;consumer=100%
  75. family=bag;stmt=cte_filtered_notmat;threshold=2000
  76. family=reuse;policy=default;refs=4;consumer=1%
  77. family=reuse;policy=mat;refs=8;consumer=1%
  78. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  79. family=bag;stmt=filtered_derived;threshold=200
  80. family=bag;stmt=comma_join;threshold=2000
  81. family=bag;stmt=cte_filtered;threshold=2000
  82. family=reuse;policy=notmat;refs=16;consumer=1%
  83. family=bag;stmt=cte_filtered_notmat;threshold=20
  84. family=bag;stmt=filtered_derived;threshold=20000
  85. family=bag;stmt=cte_filtered_notmat;threshold=200
  86. family=bag;stmt=cte_joined_outer_mat;threshold=20
  87. family=bag;stmt=cte_filtered;threshold=20
  88. family=bag;stmt=join_on;threshold=2000
  89. family=bag;stmt=cte_joined_outer_notmat;threshold=100000
  90. family=reuse;policy=notmat;refs=8;consumer=100%
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
| family=bag;stmt=comma_join;threshold=100000 | j06_comma_join | valid | `29806faf9d993eb2` |  |
| family=bag;stmt=comma_join;threshold=180000 | j06_comma_join | valid | `d9b3d8fee62c9aad` |  |
| family=bag;stmt=comma_join;threshold=20 | j06_comma_join | valid | `b99bc45005d3297d` |  |
| family=bag;stmt=comma_join;threshold=200 | j06_comma_join | valid | `e379120c5ca8ffdc` |  |
| family=bag;stmt=comma_join;threshold=2000 | j06_comma_join | valid | `3258038f84f8215e` |  |
| family=bag;stmt=comma_join;threshold=20000 | j06_comma_join | valid | `fb5cee7375f28c9d` |  |
| family=bag;stmt=cte_filtered;threshold=100000 | j04_cte_filtered | valid | `be8c0e97bc054f96` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=180000 | j04_cte_filtered | valid | `2e9fa3a0d760655b` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=20 | j04_cte_filtered | valid | `0168a2bdfaacfc42` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=200 | j04_cte_filtered | valid | `9026a869f2acc0d3` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=2000 | j04_cte_filtered | valid | `845a4dbcd5e53092` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=20000 | j04_cte_filtered | valid | `ab667545e5133f96` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_mat;threshold=100000 | j04m_cte_filtered_mat | valid | `27e92063dbecbe54` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=180000 | j04m_cte_filtered_mat | valid | `617a9802c7f1c586` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=20 | j04m_cte_filtered_mat | valid | `e1bc3d86e3ffc9f0` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=200 | j04m_cte_filtered_mat | valid | `f01d6920f4926202` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=2000 | j04m_cte_filtered_mat | valid | `06f653f990b9f0c1` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=20000 | j04m_cte_filtered_mat | valid | `41e2dfa70a0e6e96` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_notmat;threshold=100000 | j04n_cte_filtered_notmat | valid | `991d79d44ee2cfe5` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=180000 | j04n_cte_filtered_notmat | valid | `93d01c12ca2b4749` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=20 | j04n_cte_filtered_notmat | valid | `34a5c5cd9fe4b218` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=200 | j04n_cte_filtered_notmat | valid | `c397b4c6f35ff935` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=2000 | j04n_cte_filtered_notmat | valid | `19aa7ed15cd01850` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=20000 | j04n_cte_filtered_notmat | valid | `ccd7a394ee4b2cf2` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=100000 | j05_cte_joined_outer | valid | `b150d886203d7fc2` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=180000 | j05_cte_joined_outer | valid | `85105525bfd8facc` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20 | j05_cte_joined_outer | valid | `b462873d5aa390ab` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=200 | j05_cte_joined_outer | valid | `f95766154cf966e9` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=2000 | j05_cte_joined_outer | valid | `ca84a3f604cbe2ce` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20000 | j05_cte_joined_outer | valid | `66802665d99e9803` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_mat;threshold=100000 | j05m_cte_joined_outer_mat | valid | `b2e0c7f66ea3fd6d` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | j05m_cte_joined_outer_mat | valid | `76cbf35b49020c8d` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | j05m_cte_joined_outer_mat | valid | `a76609b1e9d9c4d8` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=200 | j05m_cte_joined_outer_mat | valid | `337de54331066248` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | j05m_cte_joined_outer_mat | valid | `b2ff89af10218be5` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20000 | j05m_cte_joined_outer_mat | valid | `aa2aa0966afafb13` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_notmat;threshold=100000 | j05n_cte_joined_outer_notmat | valid | `34e85e348eac0387` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | j05n_cte_joined_outer_notmat | valid | `feadbee73df310bc` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | j05n_cte_joined_outer_notmat | valid | `2e2f87197e47e59b` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=200 | j05n_cte_joined_outer_notmat | valid | `ca43bde8295194e5` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | j05n_cte_joined_outer_notmat | valid | `0aa025d8aeda0e9f` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20000 | j05n_cte_joined_outer_notmat | valid | `c381ace7a4e96e2d` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=filtered_derived;threshold=100000 | j03_filtered_derived | valid | `002032492dbd5dd9` |  |
| family=bag;stmt=filtered_derived;threshold=180000 | j03_filtered_derived | valid | `720387462908170c` |  |
| family=bag;stmt=filtered_derived;threshold=20 | j03_filtered_derived | valid | `cad8a00773bbb5f4` |  |
| family=bag;stmt=filtered_derived;threshold=200 | j03_filtered_derived | valid | `75653f29898006e1` |  |
| family=bag;stmt=filtered_derived;threshold=2000 | j03_filtered_derived | valid | `a235599a69cf5851` |  |
| family=bag;stmt=filtered_derived;threshold=20000 | j03_filtered_derived | valid | `da3a56cf3f9a0a71` |  |
| family=bag;stmt=join_on;threshold=100000 | j02_join_on | valid | `4b4e9bbc3e94b0e6` |  |
| family=bag;stmt=join_on;threshold=180000 | j02_join_on | valid | `0c9175683da8f7a1` |  |
| family=bag;stmt=join_on;threshold=20 | j02_join_on | valid | `a913236c397db0c0` |  |
| family=bag;stmt=join_on;threshold=200 | j02_join_on | valid | `27f6027ea5521619` |  |
| family=bag;stmt=join_on;threshold=2000 | j02_join_on | valid | `35aaa45287e0e5f0` |  |
| family=bag;stmt=join_on;threshold=20000 | j02_join_on | valid | `cdc0a8ad69dabd43` |  |
| family=bag;stmt=join_where;threshold=100000 | j01_join_where | valid | `8eb2e4298b9a21df` |  |
| family=bag;stmt=join_where;threshold=180000 | j01_join_where | valid | `042ef9e50c9ba842` |  |
| family=bag;stmt=join_where;threshold=20 | j01_join_where | valid | `c9b7a53a89375bce` |  |
| family=bag;stmt=join_where;threshold=200 | j01_join_where | valid | `04e4ec716105e868` |  |
| family=bag;stmt=join_where;threshold=2000 | j01_join_where | valid | `7235f6dddf24be1e` |  |
| family=bag;stmt=join_where;threshold=20000 | j01_join_where | valid | `4c398f9d21de74eb` |  |
| family=bag;stmt=join_where;threshold=20;position=last | j01_join_where | valid | `42086d2757bbec1a` |  |
| family=reuse;policy=default;refs=16;consumer=1% | r_default_n16_c1 | valid | `6c20417679f5fbfd` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=16;consumer=100% | r_default_n16_c100 | valid | `5f5d10ac883293de` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=1;consumer=1% | r_default_n1_c1 | valid | `8d398a2cee5943c3` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=1;consumer=100% | r_default_n1_c100 | valid | `f04625e62d5faada` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=2;consumer=1% | r_default_n2_c1 | valid | `e744104fdb4622d5` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=2;consumer=100% | r_default_n2_c100 | valid | `32a2b792810dc2a8` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=1% | r_default_n4_c1 | valid | `81fa286e97be8d25` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=100% | r_default_n4_c100 | valid | `9c9aa704da2fb1b4` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=8;consumer=1% | r_default_n8_c1 | valid | `1c66778d54405c53` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=8;consumer=100% | r_default_n8_c100 | valid | `6c2b01729f08759e` | policy=default; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=1% | r_mat_n16_c1 | valid | `def8116385a54d81` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=100% | r_mat_n16_c100 | valid | `7a87518a258f56ad` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=1% | r_mat_n1_c1 | valid | `6689b2d0fc8606f8` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=100% | r_mat_n1_c100 | valid | `b26644247da5b41b` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=2;consumer=1% | r_mat_n2_c1 | valid | `8b6750db6eaaa211` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=2;consumer=100% | r_mat_n2_c100 | valid | `26ed7e06a8edf4db` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=1% | r_mat_n4_c1 | valid | `5117e7368dc54c32` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=100% | r_mat_n4_c100 | valid | `5776bbb330b77a18` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=8;consumer=1% | r_mat_n8_c1 | valid | `9517f0c04a534fd9` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=8;consumer=100% | r_mat_n8_c100 | valid | `0c41888a3b10e411` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=notmat;refs=16;consumer=1% | r_notmat_n16_c1 | valid | `e9cf25e328e4e9fa` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `ce7d445a99b77c8f` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `e52181004a68fc5d` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `40f1af0c733c861f` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=2;consumer=1% | r_notmat_n2_c1 | valid | `c614e4eae8838fcf` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=2;consumer=100% | r_notmat_n2_c100 | valid | `8d48a37065cf22ac` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `0c0b0c14bddb1fe0` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `7d07c3c9c394d739` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=8;consumer=1% | r_notmat_n8_c1 | valid | `8f5e74f1794de76b` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=8;consumer=100% | r_notmat_n8_c100 | valid | `af027cc73ee78117` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, DIST, DEBUG, FORMAT TEXT` supported on yugabyte
- one host: this run used node budget 2 CPU 3g and client budget 1 CPU 2g, so budget comparisons are separate classes here, not independent hardware
- one client process with one fixed maintenance writer at 10 updates/s; readers are sequential; the reported sample latency includes retry and backoff time
- physical rank-index arm yb-range-rank; an engine whose planner ignores the index is a different physical design, not a faster one
- YugabyteDB DIST/DEBUG storage counters (rows and index rows scanned) are parsed per cell; RPC-level counters are not aggregated
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
