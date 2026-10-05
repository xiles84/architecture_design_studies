# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 91 cells, 91 valid, engine `yugabyte`, topology `single`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-v3b-domain-yb-b2-20261005`
- run tag: `run/study-01/select-join-v3b-domain-yb-b2-20261005`
- environment: `host-zenbook-ux5406sa`
- repo commit: `7243fee074bfb16fc5ad8b6afc1a2d7549080935`
- engine version: `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit`
- database engine image: `docker.io/yugabytedb/yugabyte:2025.2.6.0-b111` @ `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`
- benchmark client image: `localhost/charitybench-select:select-join-v3b-domain-yb-b2-20261005` @ `f3f1e10c44cbcb4c7b5bc56ed95a4f071682b6f7d33b5e0a7b688c4271eeef88`
- inputs digest: `5d23235fbf9206d0c6da526df936bb9cd18736e0d6d4515beb92705a6d59389a`
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
| yb-single | 934 | 2 |

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
| 20 (0.01%) | 1.8141 | 1.9269 | 1.7353 | 1.9160 | 2.2389 | 2.2016 | 1.9081 | 87.1743 | 1.9944 | 1.9759 |
| 200 (0.1%) | 2.6196 | 2.9406 | 3.2006 | 3.3555 | 2.8621 | 3.6881 | 2.8575 | 81.9603 | 3.6449 | 2.8256 |
| 2000 (1%) | 8.9702 | 9.8402 | 9.3179 | 9.1511 | 8.8490 | 9.8422 | 8.6397 | 82.1125 | 8.2550 | 8.9554 |
| 20000 (10%) | 63.5257 | 65.1925 | 63.1170 | 64.1243 | 66.3241 | 62.0557 | 63.6403 | 85.6688 | 61.1197 | 61.9678 |
| 100000 (50%) | 86.6515 | 103.7985 | 85.8603 | 84.5817 | 94.6522 | 88.9944 | 85.5083 | 97.9048 | 86.2338 | 85.8352 |
| 180000 (90%) | 105.1388 | 105.7222 | 112.5280 | 107.8939 | 114.7388 | 107.3462 | 108.8847 | 106.4139 | 108.6003 | 106.6831 |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | join_on vs join_where | 1.8141 | 1.9269 | 1.0622 | false |
| 20 | filtered_derived vs join_where | 1.8141 | 1.7353 | 0.9566 | false |
| 20 | cte_filtered vs join_where | 1.8141 | 1.9160 | 1.0562 | false |
| 20 | cte_filtered_mat vs join_where | 1.8141 | 2.2389 | 1.2342 | false |
| 20 | cte_filtered_notmat vs join_where | 1.8141 | 2.2016 | 1.2136 | false |
| 20 | cte_joined_outer vs join_where | 1.8141 | 1.9081 | 1.0518 | false |
| 20 | cte_joined_outer_mat vs join_where | 1.8141 | 87.1743 | 48.0538 | false |
| 20 | cte_joined_outer_notmat vs join_where | 1.8141 | 1.9944 | 1.0994 | false |
| 20 | comma_join vs join_where | 1.8141 | 1.9759 | 1.0892 | false |
| 200 | join_on vs join_where | 2.6196 | 2.9406 | 1.1225 | false |
| 200 | filtered_derived vs join_where | 2.6196 | 3.2006 | 1.2218 | false |
| 200 | cte_filtered vs join_where | 2.6196 | 3.3555 | 1.2809 | false |
| 200 | cte_filtered_mat vs join_where | 2.6196 | 2.8621 | 1.0926 | false |
| 200 | cte_filtered_notmat vs join_where | 2.6196 | 3.6881 | 1.4079 | false |
| 200 | cte_joined_outer vs join_where | 2.6196 | 2.8575 | 1.0908 | false |
| 200 | cte_joined_outer_mat vs join_where | 2.6196 | 81.9603 | 31.2870 | false |
| 200 | cte_joined_outer_notmat vs join_where | 2.6196 | 3.6449 | 1.3914 | false |
| 200 | comma_join vs join_where | 2.6196 | 2.8256 | 1.0786 | false |
| 2000 | join_on vs join_where | 8.9702 | 9.8402 | 1.0970 | false |
| 2000 | filtered_derived vs join_where | 8.9702 | 9.3179 | 1.0388 | false |
| 2000 | cte_filtered vs join_where | 8.9702 | 9.1511 | 1.0202 | false |
| 2000 | cte_filtered_mat vs join_where | 8.9702 | 8.8490 | 0.9865 | false |
| 2000 | cte_filtered_notmat vs join_where | 8.9702 | 9.8422 | 1.0972 | false |
| 2000 | cte_joined_outer vs join_where | 8.9702 | 8.6397 | 0.9631 | false |
| 2000 | cte_joined_outer_mat vs join_where | 8.9702 | 82.1125 | 9.1539 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 8.9702 | 8.2550 | 0.9203 | false |
| 2000 | comma_join vs join_where | 8.9702 | 8.9554 | 0.9983 | false |
| 20000 | join_on vs join_where | 63.5257 | 65.1925 | 1.0262 | false |
| 20000 | filtered_derived vs join_where | 63.5257 | 63.1170 | 0.9936 | false |
| 20000 | cte_filtered vs join_where | 63.5257 | 64.1243 | 1.0094 | false |
| 20000 | cte_filtered_mat vs join_where | 63.5257 | 66.3241 | 1.0441 | false |
| 20000 | cte_filtered_notmat vs join_where | 63.5257 | 62.0557 | 0.9769 | false |
| 20000 | cte_joined_outer vs join_where | 63.5257 | 63.6403 | 1.0018 | false |
| 20000 | cte_joined_outer_mat vs join_where | 63.5257 | 85.6688 | 1.3486 | false |
| 20000 | cte_joined_outer_notmat vs join_where | 63.5257 | 61.1197 | 0.9621 | false |
| 20000 | comma_join vs join_where | 63.5257 | 61.9678 | 0.9755 | false |
| 100000 | join_on vs join_where | 86.6515 | 103.7985 | 1.1979 | false |
| 100000 | filtered_derived vs join_where | 86.6515 | 85.8603 | 0.9909 | false |
| 100000 | cte_filtered vs join_where | 86.6515 | 84.5817 | 0.9761 | false |
| 100000 | cte_filtered_mat vs join_where | 86.6515 | 94.6522 | 1.0923 | false |
| 100000 | cte_filtered_notmat vs join_where | 86.6515 | 88.9944 | 1.0270 | false |
| 100000 | cte_joined_outer vs join_where | 86.6515 | 85.5083 | 0.9868 | false |
| 100000 | cte_joined_outer_mat vs join_where | 86.6515 | 97.9048 | 1.1299 | false |
| 100000 | cte_joined_outer_notmat vs join_where | 86.6515 | 86.2338 | 0.9952 | false |
| 100000 | comma_join vs join_where | 86.6515 | 85.8352 | 0.9906 | false |
| 180000 | join_on vs join_where | 105.1388 | 105.7222 | 1.0055 | false |
| 180000 | filtered_derived vs join_where | 105.1388 | 112.5280 | 1.0703 | false |
| 180000 | cte_filtered vs join_where | 105.1388 | 107.8939 | 1.0262 | false |
| 180000 | cte_filtered_mat vs join_where | 105.1388 | 114.7388 | 1.0913 | false |
| 180000 | cte_filtered_notmat vs join_where | 105.1388 | 107.3462 | 1.0210 | false |
| 180000 | cte_joined_outer vs join_where | 105.1388 | 108.8847 | 1.0356 | false |
| 180000 | cte_joined_outer_mat vs join_where | 105.1388 | 106.4139 | 1.0121 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 105.1388 | 108.6003 | 1.0329 | false |
| 180000 | comma_join vs join_where | 105.1388 | 106.6831 | 1.0147 | false |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_filtered_mat vs cte_filtered | 1.9160 | 2.2389 | 1.1686 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_filtered_notmat vs cte_filtered | 1.9160 | 2.2016 | 1.1491 | false | policy=not_materialized; observed_cte_scan=false |
| 20 | cte_joined_outer_mat vs cte_joined_outer | 1.9081 | 87.1743 | 45.6863 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 1.9081 | 1.9944 | 1.0452 | false | policy=not_materialized; observed_cte_scan=false |
| 200 | cte_filtered_mat vs cte_filtered | 3.3555 | 2.8621 | 0.8530 | false | policy=materialized; observed_cte_scan=true |
| 200 | cte_filtered_notmat vs cte_filtered | 3.3555 | 3.6881 | 1.0991 | false | policy=not_materialized; observed_cte_scan=false |
| 200 | cte_joined_outer_mat vs cte_joined_outer | 2.8575 | 81.9603 | 28.6823 | false | policy=materialized; observed_cte_scan=true |
| 200 | cte_joined_outer_notmat vs cte_joined_outer | 2.8575 | 3.6449 | 1.2756 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_filtered_mat vs cte_filtered | 9.1511 | 8.8490 | 0.9670 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_filtered_notmat vs cte_filtered | 9.1511 | 9.8422 | 1.0755 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 8.6397 | 82.1125 | 9.5041 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 8.6397 | 8.2550 | 0.9555 | false | policy=not_materialized; observed_cte_scan=false |
| 20000 | cte_filtered_mat vs cte_filtered | 64.1243 | 66.3241 | 1.0343 | false | policy=materialized; observed_cte_scan=true |
| 20000 | cte_filtered_notmat vs cte_filtered | 64.1243 | 62.0557 | 0.9677 | false | policy=not_materialized; observed_cte_scan=false |
| 20000 | cte_joined_outer_mat vs cte_joined_outer | 63.6403 | 85.6688 | 1.3461 | false | policy=materialized; observed_cte_scan=true |
| 20000 | cte_joined_outer_notmat vs cte_joined_outer | 63.6403 | 61.1197 | 0.9604 | false | policy=not_materialized; observed_cte_scan=false |
| 100000 | cte_filtered_mat vs cte_filtered | 84.5817 | 94.6522 | 1.1191 | false | policy=materialized; observed_cte_scan=true |
| 100000 | cte_filtered_notmat vs cte_filtered | 84.5817 | 88.9944 | 1.0522 | false | policy=not_materialized; observed_cte_scan=false |
| 100000 | cte_joined_outer_mat vs cte_joined_outer | 85.5083 | 97.9048 | 1.1450 | false | policy=materialized; observed_cte_scan=true |
| 100000 | cte_joined_outer_notmat vs cte_joined_outer | 85.5083 | 86.2338 | 1.0085 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_filtered_mat vs cte_filtered | 107.8939 | 114.7388 | 1.0634 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_filtered_notmat vs cte_filtered | 107.8939 | 107.3462 | 0.9949 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 108.8847 | 106.4139 | 0.9773 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 108.8847 | 108.6003 | 0.9974 | false | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 6.3401 | 1 | `93f72db234a57799` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 7.8966 | 1 | `e314f2723ba6d94f` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | 8.3935 | 2 | `1205bb33df332adf` | policy=default; observed_cte_scan=true | valid |
| default | 2 | 100% | 8.5670 | 2 | `9fa4d08527ba76ff` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 1% | 8.5645 | 4 | `f39edef2d565be3a` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 8.0042 | 4 | `065876977781e296` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | 8.1899 | 8 | `cb3f118a2982bf0c` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 100% | 9.4045 | 8 | `6b9f5d2b6ca37e5c` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 1% | 9.1396 | 16 | `1c003aafb72edc4b` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 10.1888 | 16 | `06f2af63c77f8ba2` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 7.9332 | 1 | `5261998e222d1233` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 8.8313 | 1 | `68c2ec77aba1d449` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | 7.9289 | 2 | `ab7f1ca09e799217` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 100% | 8.2272 | 2 | `348d210227e46e88` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 1% | 8.5997 | 4 | `13106ecc5dd4f42c` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 8.9128 | 4 | `f2a5bbc4514a5c07` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | 8.0892 | 8 | `f976ff281c8778d0` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 100% | 9.1632 | 8 | `d08113fa3684db64` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 1% | 9.3859 | 16 | `2a66197f99386b7a` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 9.7356 | 16 | `c318c8b65b294e60` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 6.9825 | 1 | `771b3eecdf82ff8b` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 7.9434 | 1 | `11ac1ed4221d77ea` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | 12.0652 | 2 | `336e089933242dcb` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 100% | 16.2400 | 2 | `e719e1aea640df85` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 1% | 23.1679 | 4 | `8361dd9ac5aba3a7` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 28.4809 | 4 | `470cd45988f4b50a` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | 57.0332 | 8 | `3392b527d4aaeab7` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 100% | 58.0818 | 8 | `1c13f2233175fe1f` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 1% | 86.1408 | 16 | `cf3538f56a345886` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 108.3816 | 16 | `29377c98ad701594` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 1.2513 | false |
| notmat | 1 | 1% | 1.1013 | false |
| mat | 1 | 100% | 1.1184 | false |
| notmat | 1 | 100% | 1.0059 | false |
| mat | 2 | 1% | 0.9446 | false |
| notmat | 2 | 1% | 1.4374 | false |
| mat | 2 | 100% | 0.9603 | false |
| notmat | 2 | 100% | 1.8956 | false |
| mat | 4 | 1% | 1.0041 | false |
| notmat | 4 | 1% | 2.7051 | false |
| mat | 4 | 100% | 1.1135 | false |
| notmat | 4 | 100% | 3.5583 | false |
| mat | 8 | 1% | 0.9877 | false |
| notmat | 8 | 1% | 6.9638 | false |
| mat | 8 | 100% | 0.9743 | false |
| notmat | 8 | 100% | 6.1760 | false |
| mat | 16 | 1% | 1.0269 | false |
| notmat | 16 | 1% | 9.4250 | false |
| mat | 16 | 100% | 0.9555 | false |
| notmat | 16 | 100% | 10.6373 | false |

## Repeated-design position control

| first median ms | last median ms | ratio (last/first) | same plan |
|---|---|---|---|
| 7.933197 | 7.456521 | 0.9399 | false |

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

- client: nr_periods 605, nr_throttled 251, throttled_usec 8872000 | quota 100000 100000
, memory 2147483648

- yb-single: nr_periods 672, nr_throttled 72, throttled_usec 2544553 | quota 200000 100000, memory 3221225472

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
| family=bag;stmt=cte_joined_outer_mat;threshold=100000 | valid | false | true | 202050 | 0 | 0/0 | 100000 | 1 |
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
| family=reuse;policy=mat;refs=1;consumer=1%;position=last | valid | true | false | 4000 | 2000 | 0/0 | 1983 | 0 |
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
| family=bag;stmt=comma_join;threshold=100000 | valid | 84.7633 | 85.8352 | 90.6317 | 5 | 11.7 |
| family=bag;stmt=comma_join;threshold=180000 | valid | 103.0613 | 106.6831 | 110.0636 | 5 | 9.4 |
| family=bag;stmt=comma_join;threshold=20 | valid | 1.6855 | 1.9759 | 2.0552 | 5 | 506.1 |
| family=bag;stmt=comma_join;threshold=200 | valid | 2.7714 | 2.8256 | 3.2010 | 5 | 353.9 |
| family=bag;stmt=comma_join;threshold=2000 | valid | 8.7244 | 8.9554 | 9.7895 | 5 | 111.7 |
| family=bag;stmt=comma_join;threshold=20000 | valid | 60.0886 | 61.9678 | 64.2018 | 5 | 16.1 |
| family=bag;stmt=cte_filtered;threshold=100000 | valid | 81.4864 | 84.5817 | 89.0418 | 5 | 11.8 |
| family=bag;stmt=cte_filtered;threshold=180000 | valid | 106.2319 | 107.8939 | 111.2061 | 5 | 9.3 |
| family=bag;stmt=cte_filtered;threshold=20 | valid | 1.8286 | 1.9160 | 2.1520 | 5 | 521.9 |
| family=bag;stmt=cte_filtered;threshold=200 | valid | 3.0834 | 3.3555 | 4.1314 | 5 | 298.0 |
| family=bag;stmt=cte_filtered;threshold=2000 | valid | 8.5081 | 9.1511 | 9.6179 | 5 | 109.3 |
| family=bag;stmt=cte_filtered;threshold=20000 | valid | 62.5714 | 64.1243 | 87.5854 | 5 | 15.6 |
| family=bag;stmt=cte_filtered_mat;threshold=100000 | valid | 88.8738 | 94.6522 | 96.9963 | 5 | 10.6 |
| family=bag;stmt=cte_filtered_mat;threshold=180000 | valid | 111.3140 | 114.7388 | 119.9404 | 5 | 8.7 |
| family=bag;stmt=cte_filtered_mat;threshold=20 | valid | 2.1491 | 2.2389 | 2.3104 | 5 | 446.6 |
| family=bag;stmt=cte_filtered_mat;threshold=200 | valid | 2.8159 | 2.8621 | 2.9586 | 5 | 349.4 |
| family=bag;stmt=cte_filtered_mat;threshold=2000 | valid | 8.3913 | 8.8490 | 9.5938 | 5 | 113.0 |
| family=bag;stmt=cte_filtered_mat;threshold=20000 | valid | 64.5654 | 66.3241 | 66.7431 | 5 | 15.1 |
| family=bag;stmt=cte_filtered_notmat;threshold=100000 | valid | 84.3300 | 88.9944 | 91.6596 | 5 | 11.2 |
| family=bag;stmt=cte_filtered_notmat;threshold=180000 | valid | 103.0038 | 107.3462 | 112.7032 | 5 | 9.3 |
| family=bag;stmt=cte_filtered_notmat;threshold=20 | valid | 2.0581 | 2.2016 | 2.2769 | 5 | 454.2 |
| family=bag;stmt=cte_filtered_notmat;threshold=200 | valid | 3.5339 | 3.6881 | 3.8808 | 5 | 271.1 |
| family=bag;stmt=cte_filtered_notmat;threshold=2000 | valid | 8.8853 | 9.8422 | 10.0065 | 5 | 101.6 |
| family=bag;stmt=cte_filtered_notmat;threshold=20000 | valid | 60.9359 | 62.0557 | 63.1397 | 5 | 16.1 |
| family=bag;stmt=cte_joined_outer;threshold=100000 | valid | 83.6432 | 85.5083 | 93.4909 | 5 | 11.7 |
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 104.4436 | 108.8847 | 109.6653 | 5 | 9.2 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 1.7910 | 1.9081 | 2.4205 | 5 | 524.1 |
| family=bag;stmt=cte_joined_outer;threshold=200 | valid | 2.6584 | 2.8575 | 2.9361 | 5 | 350.0 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 8.5528 | 8.6397 | 9.3337 | 5 | 115.7 |
| family=bag;stmt=cte_joined_outer;threshold=20000 | valid | 60.4780 | 63.6403 | 66.0410 | 5 | 15.7 |
| family=bag;stmt=cte_joined_outer_mat;threshold=100000 | valid | 92.0320 | 97.9048 | 179.7601 | 5 | 10.2 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 103.6654 | 106.4139 | 110.0470 | 5 | 9.4 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 82.1990 | 87.1743 | 98.3301 | 5 | 11.5 |
| family=bag;stmt=cte_joined_outer_mat;threshold=200 | valid | 77.0635 | 81.9603 | 85.7602 | 5 | 12.2 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 80.5686 | 82.1125 | 87.8390 | 5 | 12.2 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20000 | valid | 85.0211 | 85.6688 | 88.3706 | 5 | 11.7 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=100000 | valid | 82.4922 | 86.2338 | 91.1959 | 5 | 11.6 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 102.4324 | 108.6003 | 109.3736 | 5 | 9.2 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 1.7962 | 1.9944 | 2.2115 | 5 | 501.4 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=200 | valid | 3.1842 | 3.6449 | 3.8310 | 5 | 274.4 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 8.1277 | 8.2550 | 9.1196 | 5 | 121.1 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20000 | valid | 59.9062 | 61.1197 | 62.1028 | 5 | 16.4 |
| family=bag;stmt=filtered_derived;threshold=100000 | valid | 85.1024 | 85.8603 | 89.9854 | 5 | 11.6 |
| family=bag;stmt=filtered_derived;threshold=180000 | valid | 104.1163 | 112.5280 | 113.6305 | 5 | 8.9 |
| family=bag;stmt=filtered_derived;threshold=20 | valid | 1.6038 | 1.7353 | 1.8472 | 5 | 576.3 |
| family=bag;stmt=filtered_derived;threshold=200 | valid | 2.8142 | 3.2006 | 3.7054 | 5 | 312.4 |
| family=bag;stmt=filtered_derived;threshold=2000 | valid | 8.7577 | 9.3179 | 11.2769 | 5 | 107.3 |
| family=bag;stmt=filtered_derived;threshold=20000 | valid | 62.9366 | 63.1170 | 63.8047 | 5 | 15.8 |
| family=bag;stmt=join_on;threshold=100000 | valid | 89.2214 | 103.7985 | 116.0311 | 5 | 9.6 |
| family=bag;stmt=join_on;threshold=180000 | valid | 103.8650 | 105.7222 | 111.5068 | 5 | 9.5 |
| family=bag;stmt=join_on;threshold=20 | valid | 1.6586 | 1.9269 | 2.2529 | 5 | 519.0 |
| family=bag;stmt=join_on;threshold=200 | valid | 2.7761 | 2.9406 | 3.2718 | 5 | 340.1 |
| family=bag;stmt=join_on;threshold=2000 | valid | 9.2023 | 9.8402 | 10.2419 | 5 | 101.6 |
| family=bag;stmt=join_on;threshold=20000 | valid | 62.4068 | 65.1925 | 66.5710 | 5 | 15.3 |
| family=bag;stmt=join_where;threshold=100000 | valid | 83.1076 | 86.6515 | 93.6373 | 5 | 11.5 |
| family=bag;stmt=join_where;threshold=180000 | valid | 104.4007 | 105.1388 | 113.7498 | 5 | 9.5 |
| family=bag;stmt=join_where;threshold=20 | valid | 1.6242 | 1.8141 | 1.9928 | 5 | 551.2 |
| family=bag;stmt=join_where;threshold=200 | valid | 2.3375 | 2.6196 | 2.6932 | 5 | 381.7 |
| family=bag;stmt=join_where;threshold=2000 | valid | 8.0802 | 8.9702 | 10.0088 | 5 | 111.5 |
| family=bag;stmt=join_where;threshold=20000 | valid | 62.0344 | 63.5257 | 67.0433 | 5 | 15.7 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 8.8708 | 9.1396 | 9.7443 | 5 | 109.4 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 9.8677 | 10.1888 | 11.4873 | 5 | 98.1 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 6.1073 | 6.3401 | 6.9955 | 5 | 157.7 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 7.6988 | 7.8966 | 8.1972 | 5 | 126.6 |
| family=reuse;policy=default;refs=2;consumer=1% | valid | 7.9201 | 8.3935 | 8.9865 | 5 | 119.1 |
| family=reuse;policy=default;refs=2;consumer=100% | valid | 8.3441 | 8.5670 | 8.9760 | 5 | 116.7 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 8.2411 | 8.5645 | 8.8943 | 5 | 116.8 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 7.8231 | 8.0042 | 8.3492 | 5 | 124.9 |
| family=reuse;policy=default;refs=8;consumer=1% | valid | 8.0924 | 8.1899 | 8.3465 | 5 | 122.1 |
| family=reuse;policy=default;refs=8;consumer=100% | valid | 8.7657 | 9.4045 | 9.6923 | 5 | 106.3 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 8.8663 | 9.3859 | 9.7516 | 5 | 106.5 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 9.0785 | 9.7356 | 10.6317 | 5 | 102.7 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 7.7766 | 7.9332 | 8.6546 | 5 | 126.1 |
| family=reuse;policy=mat;refs=1;consumer=1%;position=last | valid | 6.9320 | 7.4565 | 7.7459 | 5 | 134.1 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 8.1075 | 8.8313 | 9.1325 | 5 | 113.2 |
| family=reuse;policy=mat;refs=2;consumer=1% | valid | 7.7873 | 7.9289 | 8.6596 | 5 | 126.1 |
| family=reuse;policy=mat;refs=2;consumer=100% | valid | 7.9345 | 8.2272 | 10.4654 | 5 | 121.5 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 8.3128 | 8.5997 | 9.1563 | 5 | 116.3 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 8.4044 | 8.9128 | 36.0261 | 5 | 112.2 |
| family=reuse;policy=mat;refs=8;consumer=1% | valid | 7.8076 | 8.0892 | 8.2700 | 5 | 123.6 |
| family=reuse;policy=mat;refs=8;consumer=100% | valid | 8.9315 | 9.1632 | 9.3916 | 5 | 109.1 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 74.2591 | 86.1408 | 94.4510 | 5 | 11.6 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 84.7369 | 108.3816 | 123.7329 | 5 | 9.2 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 6.2907 | 6.9825 | 7.5237 | 5 | 143.2 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 7.4201 | 7.9434 | 7.9988 | 5 | 125.9 |
| family=reuse;policy=notmat;refs=2;consumer=1% | valid | 11.5660 | 12.0652 | 12.7542 | 5 | 82.9 |
| family=reuse;policy=notmat;refs=2;consumer=100% | valid | 14.3283 | 16.2400 | 17.6219 | 5 | 61.6 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 21.7536 | 23.1679 | 24.5001 | 5 | 43.2 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 27.7003 | 28.4809 | 30.7821 | 5 | 35.1 |
| family=reuse;policy=notmat;refs=8;consumer=1% | valid | 54.2300 | 57.0332 | 98.2339 | 5 | 17.5 |
| family=reuse;policy=notmat;refs=8;consumer=100% | valid | 55.1428 | 58.0818 | 66.9861 | 5 | 17.2 |

## Matched blended workload

- phase: `blended`; target offered rate: 10 updates/s
- offered 523, committed 523, failed 0 over 52.257 s
- offered rate 10.008/s, completed rate 10.008/s
- writer operations per endpoint: [523]
- reader operations overlapping a committed write: 234 of 638
- reader retries of transient conflicts: 1
- per key: 523 keys, min 1, max 1 per key
- ledger: 523 keys reconciled, 0 mismatch
- update latency: avg 2.702 ms, min 1.444, max 152.503 over 523

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42500627`, 5 interleaved rounds
  1. family=reuse;policy=mat;refs=1;consumer=1%
  2. family=reuse;policy=notmat;refs=1;consumer=100%
  3. family=bag;stmt=cte_joined_outer;threshold=100000
  4. family=bag;stmt=cte_filtered_notmat;threshold=100000
  5. family=bag;stmt=cte_joined_outer_mat;threshold=20000
  6. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  7. family=bag;stmt=cte_joined_outer_notmat;threshold=100000
  8. family=bag;stmt=cte_filtered_notmat;threshold=2000
  9. family=reuse;policy=notmat;refs=1;consumer=1%
  10. family=bag;stmt=comma_join;threshold=180000
  11. family=reuse;policy=default;refs=2;consumer=100%
  12. family=bag;stmt=cte_filtered_mat;threshold=2000
  13. family=bag;stmt=cte_joined_outer;threshold=20
  14. family=bag;stmt=comma_join;threshold=200
  15. family=reuse;policy=default;refs=1;consumer=1%
  16. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  17. family=bag;stmt=cte_joined_outer_mat;threshold=200
  18. family=reuse;policy=notmat;refs=2;consumer=100%
  19. family=bag;stmt=join_where;threshold=2000
  20. family=reuse;policy=mat;refs=1;consumer=100%
  21. family=reuse;policy=notmat;refs=16;consumer=100%
  22. family=bag;stmt=cte_filtered;threshold=20
  23. family=bag;stmt=cte_filtered;threshold=20000
  24. family=bag;stmt=filtered_derived;threshold=2000
  25. family=reuse;policy=default;refs=1;consumer=100%
  26. family=bag;stmt=cte_filtered;threshold=2000
  27. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  28. family=bag;stmt=comma_join;threshold=2000
  29. family=reuse;policy=mat;refs=16;consumer=100%
  30. family=bag;stmt=cte_joined_outer;threshold=2000
  31. family=reuse;policy=mat;refs=2;consumer=100%
  32. family=bag;stmt=join_where;threshold=20000
  33. family=bag;stmt=cte_filtered_mat;threshold=20
  34. family=bag;stmt=cte_filtered_notmat;threshold=180000
  35. family=bag;stmt=join_on;threshold=2000
  36. family=bag;stmt=cte_joined_outer;threshold=180000
  37. family=bag;stmt=cte_joined_outer;threshold=20000
  38. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  39. family=bag;stmt=cte_filtered;threshold=200
  40. family=bag;stmt=filtered_derived;threshold=20000
  41. family=reuse;policy=default;refs=16;consumer=1%
  42. family=bag;stmt=cte_filtered_mat;threshold=200
  43. family=bag;stmt=filtered_derived;threshold=180000
  44. family=bag;stmt=cte_filtered;threshold=100000
  45. family=bag;stmt=cte_joined_outer_notmat;threshold=200
  46. family=bag;stmt=comma_join;threshold=20
  47. family=bag;stmt=cte_filtered_notmat;threshold=20000
  48. family=bag;stmt=cte_filtered_mat;threshold=20000
  49. family=bag;stmt=cte_filtered_notmat;threshold=200
  50. family=reuse;policy=notmat;refs=4;consumer=1%
  51. family=reuse;policy=mat;refs=8;consumer=1%
  52. family=bag;stmt=filtered_derived;threshold=20
  53. family=reuse;policy=mat;refs=2;consumer=1%
  54. family=reuse;policy=notmat;refs=4;consumer=100%
  55. family=bag;stmt=cte_joined_outer;threshold=200
  56. family=bag;stmt=join_where;threshold=200
  57. family=reuse;policy=default;refs=8;consumer=1%
  58. family=reuse;policy=notmat;refs=8;consumer=100%
  59. family=bag;stmt=join_on;threshold=100000
  60. family=bag;stmt=cte_filtered_notmat;threshold=20
  61. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  62. family=reuse;policy=notmat;refs=2;consumer=1%
  63. family=bag;stmt=join_where;threshold=100000
  64. family=reuse;policy=notmat;refs=8;consumer=1%
  65. family=reuse;policy=mat;refs=4;consumer=100%
  66. family=bag;stmt=filtered_derived;threshold=200
  67. family=bag;stmt=cte_filtered;threshold=180000
  68. family=reuse;policy=default;refs=4;consumer=1%
  69. family=reuse;policy=default;refs=4;consumer=100%
  70. family=reuse;policy=notmat;refs=16;consumer=1%
  71. family=bag;stmt=cte_joined_outer_mat;threshold=20
  72. family=reuse;policy=default;refs=2;consumer=1%
  73. family=bag;stmt=cte_filtered_mat;threshold=100000
  74. family=bag;stmt=cte_joined_outer_mat;threshold=100000
  75. family=reuse;policy=default;refs=16;consumer=100%
  76. family=bag;stmt=join_on;threshold=20
  77. family=bag;stmt=join_on;threshold=200
  78. family=bag;stmt=comma_join;threshold=20000
  79. family=reuse;policy=default;refs=8;consumer=100%
  80. family=bag;stmt=join_on;threshold=180000
  81. family=reuse;policy=mat;refs=16;consumer=1%
  82. family=bag;stmt=join_where;threshold=180000
  83. family=bag;stmt=filtered_derived;threshold=100000
  84. family=reuse;policy=mat;refs=4;consumer=1%
  85. family=bag;stmt=cte_joined_outer_notmat;threshold=20000
  86. family=bag;stmt=comma_join;threshold=100000
  87. family=bag;stmt=cte_filtered_mat;threshold=180000
  88. family=bag;stmt=join_on;threshold=20000
  89. family=reuse;policy=mat;refs=8;consumer=100%
  90. family=bag;stmt=join_where;threshold=20
  91. family=reuse;policy=mat;refs=1;consumer=1%;position=last

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
| family=bag;stmt=comma_join;threshold=100000 | j06_comma_join | valid | `0fed8684e42b4f27` |  |
| family=bag;stmt=comma_join;threshold=180000 | j06_comma_join | valid | `c08084c9ec84cae3` |  |
| family=bag;stmt=comma_join;threshold=20 | j06_comma_join | valid | `5b09e32327f88cb6` |  |
| family=bag;stmt=comma_join;threshold=200 | j06_comma_join | valid | `628e7275c2fe8133` |  |
| family=bag;stmt=comma_join;threshold=2000 | j06_comma_join | valid | `380efbcc485bc934` |  |
| family=bag;stmt=comma_join;threshold=20000 | j06_comma_join | valid | `f2f0d982d41df47c` |  |
| family=bag;stmt=cte_filtered;threshold=100000 | j04_cte_filtered | valid | `38ba7961b76cbab2` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=180000 | j04_cte_filtered | valid | `829e577c9526ea04` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=20 | j04_cte_filtered | valid | `c943f3e4fd44504a` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=200 | j04_cte_filtered | valid | `ed1a1afb30ce86ab` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=2000 | j04_cte_filtered | valid | `08cacdf96c5e26e3` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=20000 | j04_cte_filtered | valid | `059b9badb523dd1a` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_mat;threshold=100000 | j04m_cte_filtered_mat | valid | `e184d8047fd5efb9` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=180000 | j04m_cte_filtered_mat | valid | `684547e569093d79` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=20 | j04m_cte_filtered_mat | valid | `820064ba1f842410` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=200 | j04m_cte_filtered_mat | valid | `47e889fd2b4afa1c` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=2000 | j04m_cte_filtered_mat | valid | `110131b60ebe19ac` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=20000 | j04m_cte_filtered_mat | valid | `4074546221a99136` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_notmat;threshold=100000 | j04n_cte_filtered_notmat | valid | `a15184eab35ee62c` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=180000 | j04n_cte_filtered_notmat | valid | `2a2075847e9570a3` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=20 | j04n_cte_filtered_notmat | valid | `4c61e9544c6151ac` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=200 | j04n_cte_filtered_notmat | valid | `ef21c53d6d1b8717` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=2000 | j04n_cte_filtered_notmat | valid | `7cd27431c05dfddc` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=20000 | j04n_cte_filtered_notmat | valid | `35709e286f844788` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=100000 | j05_cte_joined_outer | valid | `7a098e3ce1e579c3` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=180000 | j05_cte_joined_outer | valid | `84e8f79714a0c566` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20 | j05_cte_joined_outer | valid | `7c6f2acff40250f5` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=200 | j05_cte_joined_outer | valid | `afcdf705e6d2d167` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=2000 | j05_cte_joined_outer | valid | `0e7436552ccc9ee8` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20000 | j05_cte_joined_outer | valid | `595c922ccec1684a` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_mat;threshold=100000 | j05m_cte_joined_outer_mat | valid | `12112d02e42acb80` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | j05m_cte_joined_outer_mat | valid | `0a88ac71fa0ba8da` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | j05m_cte_joined_outer_mat | valid | `f1ecdfb14ce979f9` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=200 | j05m_cte_joined_outer_mat | valid | `e610035adc678802` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | j05m_cte_joined_outer_mat | valid | `8fa2e4f0afe8fc6f` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20000 | j05m_cte_joined_outer_mat | valid | `1135856bf6415c92` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_notmat;threshold=100000 | j05n_cte_joined_outer_notmat | valid | `ce854c9d2a847077` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | j05n_cte_joined_outer_notmat | valid | `50237692781b8c29` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | j05n_cte_joined_outer_notmat | valid | `83cec787606d459e` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=200 | j05n_cte_joined_outer_notmat | valid | `020031651cc3a6cb` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | j05n_cte_joined_outer_notmat | valid | `cd4b77eee8377a87` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20000 | j05n_cte_joined_outer_notmat | valid | `ac658ceb7c3db0ca` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=filtered_derived;threshold=100000 | j03_filtered_derived | valid | `c16404c61ae8b7c8` |  |
| family=bag;stmt=filtered_derived;threshold=180000 | j03_filtered_derived | valid | `ca968901fba2ca5b` |  |
| family=bag;stmt=filtered_derived;threshold=20 | j03_filtered_derived | valid | `26db2b7e4bcea694` |  |
| family=bag;stmt=filtered_derived;threshold=200 | j03_filtered_derived | valid | `21e0f0970aa55e11` |  |
| family=bag;stmt=filtered_derived;threshold=2000 | j03_filtered_derived | valid | `86b74a930345b977` |  |
| family=bag;stmt=filtered_derived;threshold=20000 | j03_filtered_derived | valid | `81c69250c867e0fc` |  |
| family=bag;stmt=join_on;threshold=100000 | j02_join_on | valid | `97e0b5a9a702288d` |  |
| family=bag;stmt=join_on;threshold=180000 | j02_join_on | valid | `50fb73382a42200a` |  |
| family=bag;stmt=join_on;threshold=20 | j02_join_on | valid | `c873335283f38f92` |  |
| family=bag;stmt=join_on;threshold=200 | j02_join_on | valid | `a5893c301ce35c0e` |  |
| family=bag;stmt=join_on;threshold=2000 | j02_join_on | valid | `6cb2c2a3bb7f3871` |  |
| family=bag;stmt=join_on;threshold=20000 | j02_join_on | valid | `3cd5140eda6f972f` |  |
| family=bag;stmt=join_where;threshold=100000 | j01_join_where | valid | `79fb5bd161704f49` |  |
| family=bag;stmt=join_where;threshold=180000 | j01_join_where | valid | `fbdc2e878ef55e9f` |  |
| family=bag;stmt=join_where;threshold=20 | j01_join_where | valid | `a72e6e2647c49a2a` |  |
| family=bag;stmt=join_where;threshold=200 | j01_join_where | valid | `8d992812ba259e98` |  |
| family=bag;stmt=join_where;threshold=2000 | j01_join_where | valid | `a03ff066ca7d0f36` |  |
| family=bag;stmt=join_where;threshold=20000 | j01_join_where | valid | `8add9d00c6de1fb5` |  |
| family=reuse;policy=default;refs=16;consumer=1% | r_default_n16_c1 | valid | `1c003aafb72edc4b` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=16;consumer=100% | r_default_n16_c100 | valid | `06f2af63c77f8ba2` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=1;consumer=1% | r_default_n1_c1 | valid | `93f72db234a57799` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=1;consumer=100% | r_default_n1_c100 | valid | `e314f2723ba6d94f` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=2;consumer=1% | r_default_n2_c1 | valid | `1205bb33df332adf` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=2;consumer=100% | r_default_n2_c100 | valid | `9fa4d08527ba76ff` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=1% | r_default_n4_c1 | valid | `f39edef2d565be3a` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=100% | r_default_n4_c100 | valid | `065876977781e296` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=8;consumer=1% | r_default_n8_c1 | valid | `cb3f118a2982bf0c` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=8;consumer=100% | r_default_n8_c100 | valid | `6b9f5d2b6ca37e5c` | policy=default; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=1% | r_mat_n16_c1 | valid | `2a66197f99386b7a` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=100% | r_mat_n16_c100 | valid | `c318c8b65b294e60` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=1% | r_mat_n1_c1 | valid | `5261998e222d1233` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=1%;position=last | r_mat_n1_c1 | valid | `af1e2b0d7bb259b3` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=100% | r_mat_n1_c100 | valid | `68c2ec77aba1d449` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=2;consumer=1% | r_mat_n2_c1 | valid | `ab7f1ca09e799217` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=2;consumer=100% | r_mat_n2_c100 | valid | `348d210227e46e88` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=1% | r_mat_n4_c1 | valid | `13106ecc5dd4f42c` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=100% | r_mat_n4_c100 | valid | `f2a5bbc4514a5c07` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=8;consumer=1% | r_mat_n8_c1 | valid | `f976ff281c8778d0` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=8;consumer=100% | r_mat_n8_c100 | valid | `d08113fa3684db64` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=notmat;refs=16;consumer=1% | r_notmat_n16_c1 | valid | `cf3538f56a345886` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `29377c98ad701594` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `771b3eecdf82ff8b` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `11ac1ed4221d77ea` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=2;consumer=1% | r_notmat_n2_c1 | valid | `336e089933242dcb` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=2;consumer=100% | r_notmat_n2_c100 | valid | `e719e1aea640df85` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `8361dd9ac5aba3a7` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `470cd45988f4b50a` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=8;consumer=1% | r_notmat_n8_c1 | valid | `3392b527d4aaeab7` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=8;consumer=100% | r_notmat_n8_c100 | valid | `1c13f2233175fe1f` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, DIST, DEBUG, FORMAT TEXT` supported on yugabyte
- one host: this run used node budget 2 CPU 3g and client budget 1 CPU 2g, so budget comparisons are separate classes here, not independent hardware
- one client process with one fixed maintenance writer at 10 updates/s; readers are sequential; the reported sample latency includes retry and backoff time
- physical rank-index arm yb-range-rank; an engine whose planner ignores the index is a different physical design, not a faster one
- YugabyteDB DIST/DEBUG storage counters (rows and index rows scanned) are parsed per cell; RPC-level counters are not aggregated
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
