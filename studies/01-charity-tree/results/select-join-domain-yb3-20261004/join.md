# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 91 cells, 91 valid, engine `yugabyte`, topology `single`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-domain-yb3-20261004`
- run tag: `run/study-01/select-join-domain-yb3-20261004`
- environment: `host-zenbook-ux5406sa`
- repo commit: `e8b0fb1e56da6ae4b7a64f0256c588ea08e983d3`
- engine version: `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit`
- database engine image: `docker.io/yugabytedb/yugabyte:2025.2.6.0-b111` @ `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`
- benchmark client image: `localhost/charitybench-select:select-join-domain-yb3-20261004` @ `bc7ce0737685b8d0125c7291b019ef64d8dcf34c681bcc75569f23b2d083414b`
- inputs digest: `0cdb62b8b1c5504fb37fa69da1947b23fcdd09279b8af49d840f5591af8e0a30`
- dataset fingerprint: `fcf26e08f377ec51`
- client CPU quota: `100000 100000`, memory: `2147483648`
- setting `client_workers`: 1 sequential reader; 0 writers
- setting `connection_pool`: 2 = 1 reader/setup + 1 reserve
- setting `customers`: 2000
- setting `dataset`: study01-person-donation-charity
- setting `distributed_rpc`: not parsed; raw EXPLAIN retained (coverage gap)
- setting `fact_count`: 200000
- setting `index_set`: child rank index and parent-key index, primary keys; no per-spelling index
- setting `prepared_plan_policy`: pgx simple protocol; no prepared statements
- setting `scale`: cte-domain
- setting `selectivity_mapping`: rank <= threshold; rank is a deterministic shuffle of 1..fact_count
- setting `statement_budget_ms`: 120000
- setting `statistics`: ANALYZE after deterministic load
- setting `thresholds`: [20 200 2000 20000 100000 180000]
- setting `topology`: single
- setting `trials_per_cell`: 5
- setting `whole_run_count`: 1; no replicated winner in this file

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
| j01_join_where | valid | false |  |
| j02_join_on | valid | false |  |
| j03_filtered_derived | valid | false |  |
| j04_cte_filtered | valid | false | policy=default; observed_cte_scan=false |
| j04m_cte_filtered_mat | valid | false | policy=materialized; observed_cte_scan=true |
| j04n_cte_filtered_notmat | valid | false | policy=not_materialized; observed_cte_scan=false |
| j05_cte_joined_outer | valid | false | policy=default; observed_cte_scan=false |
| j05m_cte_joined_outer_mat | valid | false | policy=materialized; observed_cte_scan=true |
| j05n_cte_joined_outer_notmat | valid | false | policy=not_materialized; observed_cte_scan=false |
| j06_comma_join | valid | false |  |
| j11_optional_on | valid | false |  |
| r_default_n4_c100 | valid | false | policy=default; observed_cte_scan=true |
| r_mat_n4_c100 | valid | false | policy=materialized; observed_cte_scan=true |
| r_notmat_n4_c100 | valid | false | policy=not_materialized; observed_cte_scan=false |

## P1 — six inner-join spellings and their CTE policies (median ms)

| threshold | join_where | join_on | filtered_derived | cte_filtered | cte_filtered_mat | cte_filtered_notmat | cte_joined_outer | cte_joined_outer_mat | cte_joined_outer_notmat | comma_join |
|---|---|---|---|---|---|---|---|---|---|---|
| 20 (0.01%) | 87.1451 | 110.1639 | 93.5598 | 93.0575 | 103.9383 | 73.5282 | 85.2483 | 214.1780 | 73.8782 | 75.0968 |
| 200 (0.1%) | 103.9895 | 83.4269 | 113.0063 | 110.6672 | 78.2634 | 109.2058 | 95.4511 | 136.4335 | 79.0972 | 72.1468 |
| 2000 (1%) | 106.0996 | 80.4195 | 121.5485 | 74.6743 | 82.7469 | 74.9934 | 84.1131 | 144.2880 | 74.7956 | 99.5767 |
| 20000 (10%) | 101.1035 | 98.6984 | 105.5479 | 124.0584 | 89.2912 | 101.3577 | 94.0174 | 226.3498 | 120.7276 | 115.4632 |
| 100000 (50%) | 174.8875 | 145.1905 | 177.5349 | 139.6369 | 150.6452 | 156.0797 | 171.9248 | 172.9103 | 128.2294 | 188.7832 |
| 180000 (90%) | 225.4778 | 183.7657 | 243.6141 | 232.7039 | 205.2236 | 209.7476 | 168.7611 | 229.4704 | 184.6368 | 174.5075 |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | join_on vs join_where | 87.1451 | 110.1639 | 1.2641 | false |
| 20 | filtered_derived vs join_where | 87.1451 | 93.5598 | 1.0736 | false |
| 20 | cte_filtered vs join_where | 87.1451 | 93.0575 | 1.0678 | false |
| 20 | cte_filtered_mat vs join_where | 87.1451 | 103.9383 | 1.1927 | false |
| 20 | cte_filtered_notmat vs join_where | 87.1451 | 73.5282 | 0.8437 | false |
| 20 | cte_joined_outer vs join_where | 87.1451 | 85.2483 | 0.9782 | false |
| 20 | cte_joined_outer_mat vs join_where | 87.1451 | 214.1780 | 2.4577 | false |
| 20 | cte_joined_outer_notmat vs join_where | 87.1451 | 73.8782 | 0.8478 | false |
| 20 | comma_join vs join_where | 87.1451 | 75.0968 | 0.8617 | false |
| 200 | join_on vs join_where | 103.9895 | 83.4269 | 0.8023 | false |
| 200 | filtered_derived vs join_where | 103.9895 | 113.0063 | 1.0867 | false |
| 200 | cte_filtered vs join_where | 103.9895 | 110.6672 | 1.0642 | false |
| 200 | cte_filtered_mat vs join_where | 103.9895 | 78.2634 | 0.7526 | false |
| 200 | cte_filtered_notmat vs join_where | 103.9895 | 109.2058 | 1.0502 | false |
| 200 | cte_joined_outer vs join_where | 103.9895 | 95.4511 | 0.9179 | false |
| 200 | cte_joined_outer_mat vs join_where | 103.9895 | 136.4335 | 1.3120 | false |
| 200 | cte_joined_outer_notmat vs join_where | 103.9895 | 79.0972 | 0.7606 | false |
| 200 | comma_join vs join_where | 103.9895 | 72.1468 | 0.6938 | false |
| 2000 | join_on vs join_where | 106.0996 | 80.4195 | 0.7580 | false |
| 2000 | filtered_derived vs join_where | 106.0996 | 121.5485 | 1.1456 | false |
| 2000 | cte_filtered vs join_where | 106.0996 | 74.6743 | 0.7038 | false |
| 2000 | cte_filtered_mat vs join_where | 106.0996 | 82.7469 | 0.7799 | false |
| 2000 | cte_filtered_notmat vs join_where | 106.0996 | 74.9934 | 0.7068 | false |
| 2000 | cte_joined_outer vs join_where | 106.0996 | 84.1131 | 0.7928 | false |
| 2000 | cte_joined_outer_mat vs join_where | 106.0996 | 144.2880 | 1.3599 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 106.0996 | 74.7956 | 0.7050 | false |
| 2000 | comma_join vs join_where | 106.0996 | 99.5767 | 0.9385 | false |
| 20000 | join_on vs join_where | 101.1035 | 98.6984 | 0.9762 | false |
| 20000 | filtered_derived vs join_where | 101.1035 | 105.5479 | 1.0440 | false |
| 20000 | cte_filtered vs join_where | 101.1035 | 124.0584 | 1.2270 | false |
| 20000 | cte_filtered_mat vs join_where | 101.1035 | 89.2912 | 0.8832 | false |
| 20000 | cte_filtered_notmat vs join_where | 101.1035 | 101.3577 | 1.0025 | false |
| 20000 | cte_joined_outer vs join_where | 101.1035 | 94.0174 | 0.9299 | false |
| 20000 | cte_joined_outer_mat vs join_where | 101.1035 | 226.3498 | 2.2388 | false |
| 20000 | cte_joined_outer_notmat vs join_where | 101.1035 | 120.7276 | 1.1941 | false |
| 20000 | comma_join vs join_where | 101.1035 | 115.4632 | 1.1420 | false |
| 100000 | join_on vs join_where | 174.8875 | 145.1905 | 0.8302 | false |
| 100000 | filtered_derived vs join_where | 174.8875 | 177.5349 | 1.0151 | false |
| 100000 | cte_filtered vs join_where | 174.8875 | 139.6369 | 0.7984 | false |
| 100000 | cte_filtered_mat vs join_where | 174.8875 | 150.6452 | 0.8614 | false |
| 100000 | cte_filtered_notmat vs join_where | 174.8875 | 156.0797 | 0.8925 | false |
| 100000 | cte_joined_outer vs join_where | 174.8875 | 171.9248 | 0.9831 | false |
| 100000 | cte_joined_outer_mat vs join_where | 174.8875 | 172.9103 | 0.9887 | false |
| 100000 | cte_joined_outer_notmat vs join_where | 174.8875 | 128.2294 | 0.7332 | false |
| 100000 | comma_join vs join_where | 174.8875 | 188.7832 | 1.0795 | false |
| 180000 | join_on vs join_where | 225.4778 | 183.7657 | 0.8150 | false |
| 180000 | filtered_derived vs join_where | 225.4778 | 243.6141 | 1.0804 | false |
| 180000 | cte_filtered vs join_where | 225.4778 | 232.7039 | 1.0320 | false |
| 180000 | cte_filtered_mat vs join_where | 225.4778 | 205.2236 | 0.9102 | false |
| 180000 | cte_filtered_notmat vs join_where | 225.4778 | 209.7476 | 0.9302 | false |
| 180000 | cte_joined_outer vs join_where | 225.4778 | 168.7611 | 0.7485 | false |
| 180000 | cte_joined_outer_mat vs join_where | 225.4778 | 229.4704 | 1.0177 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 225.4778 | 184.6368 | 0.8189 | false |
| 180000 | comma_join vs join_where | 225.4778 | 174.5075 | 0.7739 | false |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_filtered_mat vs cte_filtered | 93.0575 | 103.9383 | 1.1169 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_filtered_notmat vs cte_filtered | 93.0575 | 73.5282 | 0.7901 | false | policy=not_materialized; observed_cte_scan=false |
| 20 | cte_joined_outer_mat vs cte_joined_outer | 85.2483 | 214.1780 | 2.5124 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 85.2483 | 73.8782 | 0.8666 | false | policy=not_materialized; observed_cte_scan=false |
| 200 | cte_filtered_mat vs cte_filtered | 110.6672 | 78.2634 | 0.7072 | false | policy=materialized; observed_cte_scan=true |
| 200 | cte_filtered_notmat vs cte_filtered | 110.6672 | 109.2058 | 0.9868 | false | policy=not_materialized; observed_cte_scan=false |
| 200 | cte_joined_outer_mat vs cte_joined_outer | 95.4511 | 136.4335 | 1.4294 | false | policy=materialized; observed_cte_scan=true |
| 200 | cte_joined_outer_notmat vs cte_joined_outer | 95.4511 | 79.0972 | 0.8287 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_filtered_mat vs cte_filtered | 74.6743 | 82.7469 | 1.1081 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_filtered_notmat vs cte_filtered | 74.6743 | 74.9934 | 1.0043 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 84.1131 | 144.2880 | 1.7154 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 84.1131 | 74.7956 | 0.8892 | false | policy=not_materialized; observed_cte_scan=false |
| 20000 | cte_filtered_mat vs cte_filtered | 124.0584 | 89.2912 | 0.7198 | false | policy=materialized; observed_cte_scan=true |
| 20000 | cte_filtered_notmat vs cte_filtered | 124.0584 | 101.3577 | 0.8170 | false | policy=not_materialized; observed_cte_scan=false |
| 20000 | cte_joined_outer_mat vs cte_joined_outer | 94.0174 | 226.3498 | 2.4075 | false | policy=materialized; observed_cte_scan=true |
| 20000 | cte_joined_outer_notmat vs cte_joined_outer | 94.0174 | 120.7276 | 1.2841 | false | policy=not_materialized; observed_cte_scan=false |
| 100000 | cte_filtered_mat vs cte_filtered | 139.6369 | 150.6452 | 1.0788 | false | policy=materialized; observed_cte_scan=true |
| 100000 | cte_filtered_notmat vs cte_filtered | 139.6369 | 156.0797 | 1.1178 | false | policy=not_materialized; observed_cte_scan=false |
| 100000 | cte_joined_outer_mat vs cte_joined_outer | 171.9248 | 172.9103 | 1.0057 | false | policy=materialized; observed_cte_scan=true |
| 100000 | cte_joined_outer_notmat vs cte_joined_outer | 171.9248 | 128.2294 | 0.7458 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_filtered_mat vs cte_filtered | 232.7039 | 205.2236 | 0.8819 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_filtered_notmat vs cte_filtered | 232.7039 | 209.7476 | 0.9013 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 168.7611 | 229.4704 | 1.3597 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 168.7611 | 184.6368 | 1.0941 | false | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 92.4769 | 1 | `70fa7780328362a2` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 85.9608 | 1 | `433eb554db8690b9` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | 91.0910 | 2 | `666cdb8c3f3a5a86` | policy=default; observed_cte_scan=true | valid |
| default | 2 | 100% | 76.7118 | 2 | `96933d0bf8432dfb` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 1% | 88.9721 | 4 | `7e0bf55a42ae0dfa` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 81.1366 | 4 | `a4f852d3e5dc0cb1` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | 72.9649 | 8 | `bf9f5223db607fb2` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 100% | 70.5471 | 8 | `ff67b98e8d61965a` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 1% | 73.8573 | 16 | `e9949a76d592b697` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 93.9867 | 16 | `77b5f8d3514276a6` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 87.8327 | 1 | `bb52eb04f09e65bc` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 102.5196 | 1 | `12c0a0a668dcc12a` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | 94.1085 | 2 | `3283f00b5edf0012` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 100% | 100.8830 | 2 | `1a8ef14e201df112` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 1% | 77.3296 | 4 | `cfe3fb80362a9c8a` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 89.5087 | 4 | `7aa8c8c2a4b3205e` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | 90.8158 | 8 | `8a92857e216ae7f5` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 100% | 76.8639 | 8 | `eb735dea667440af` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 1% | 82.0594 | 16 | `1732500e975be145` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 78.3027 | 16 | `a711475916bd2835` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 77.4532 | 1 | `4388b1e94492f68b` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 97.9613 | 1 | `5420d77bf93dd781` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | 204.0947 | 2 | `cf6b13768743afa8` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 100% | 219.6598 | 2 | `6cfb5c00796dc850` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 1% | 359.5680 | 4 | `3888327d375b83b7` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 285.6017 | 4 | `ff9aec201bee741d` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | 621.2279 | 8 | `e58023205d5af8b3` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 100% | 635.6290 | 8 | `7f8ee940110a4016` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 1% | 1052.2198 | 16 | `b6c1e248463c5b00` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 1160.0005 | 16 | `69df98b52c51f084` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 0.9498 | false |
| notmat | 1 | 1% | 0.8375 | false |
| mat | 1 | 100% | 1.1926 | false |
| notmat | 1 | 100% | 1.1396 | false |
| mat | 2 | 1% | 1.0331 | false |
| notmat | 2 | 1% | 2.2406 | false |
| mat | 2 | 100% | 1.3151 | false |
| notmat | 2 | 100% | 2.8634 | false |
| mat | 4 | 1% | 0.8691 | false |
| notmat | 4 | 1% | 4.0414 | false |
| mat | 4 | 100% | 1.1032 | false |
| notmat | 4 | 100% | 3.5200 | false |
| mat | 8 | 1% | 1.2447 | false |
| notmat | 8 | 1% | 8.5141 | false |
| mat | 8 | 100% | 1.0895 | false |
| notmat | 8 | 100% | 9.0100 | false |
| mat | 16 | 1% | 1.1111 | false |
| notmat | 16 | 1% | 14.2467 | false |
| mat | 16 | 100% | 0.8331 | false |
| notmat | 16 | 100% | 12.3422 | false |

## Repeated-design position control

| first median ms | last median ms | ratio (last/first) | same plan |
|---|---|---|---|
| 184.636751 | 220.671625 | 1.1952 | false |

## Latency samples and derived throughput

_The reader is sequential, so `queries/s` is 1000/median_ms, a derived rate and not a concurrent-capacity test. min/max are the five sample extremes._

| axis | status | min ms | median ms | max ms | samples | derived queries/s |
|---|---|---|---|---|---|---|
| family=bag;stmt=comma_join;threshold=100000 | valid | 143.5775 | 188.7832 | 200.9888 | 5 | 5.3 |
| family=bag;stmt=comma_join;threshold=180000 | valid | 159.2727 | 174.5075 | 260.4963 | 5 | 5.7 |
| family=bag;stmt=comma_join;threshold=20 | valid | 63.4873 | 75.0968 | 98.8650 | 5 | 13.3 |
| family=bag;stmt=comma_join;threshold=200 | valid | 71.4049 | 72.1468 | 101.8517 | 5 | 13.9 |
| family=bag;stmt=comma_join;threshold=2000 | valid | 76.3240 | 99.5767 | 132.3638 | 5 | 10.0 |
| family=bag;stmt=comma_join;threshold=20000 | valid | 90.0365 | 115.4632 | 134.4770 | 5 | 8.7 |
| family=bag;stmt=cte_filtered;threshold=100000 | valid | 120.3006 | 139.6369 | 236.7630 | 5 | 7.2 |
| family=bag;stmt=cte_filtered;threshold=180000 | valid | 162.9244 | 232.7039 | 259.2880 | 5 | 4.3 |
| family=bag;stmt=cte_filtered;threshold=20 | valid | 71.6697 | 93.0575 | 110.2753 | 5 | 10.7 |
| family=bag;stmt=cte_filtered;threshold=200 | valid | 70.8615 | 110.6672 | 118.0410 | 5 | 9.0 |
| family=bag;stmt=cte_filtered;threshold=2000 | valid | 66.4730 | 74.6743 | 98.8190 | 5 | 13.4 |
| family=bag;stmt=cte_filtered;threshold=20000 | valid | 89.0103 | 124.0584 | 159.6175 | 5 | 8.1 |
| family=bag;stmt=cte_filtered_mat;threshold=100000 | valid | 121.3560 | 150.6452 | 191.3084 | 5 | 6.6 |
| family=bag;stmt=cte_filtered_mat;threshold=180000 | valid | 155.7079 | 205.2236 | 261.3214 | 5 | 4.9 |
| family=bag;stmt=cte_filtered_mat;threshold=20 | valid | 75.8154 | 103.9383 | 132.2972 | 5 | 9.6 |
| family=bag;stmt=cte_filtered_mat;threshold=200 | valid | 65.8288 | 78.2634 | 97.7805 | 5 | 12.8 |
| family=bag;stmt=cte_filtered_mat;threshold=2000 | valid | 72.2803 | 82.7469 | 121.8181 | 5 | 12.1 |
| family=bag;stmt=cte_filtered_mat;threshold=20000 | valid | 79.9848 | 89.2912 | 135.4109 | 5 | 11.2 |
| family=bag;stmt=cte_filtered_notmat;threshold=100000 | valid | 121.9224 | 156.0797 | 197.2959 | 5 | 6.4 |
| family=bag;stmt=cte_filtered_notmat;threshold=180000 | valid | 163.2852 | 209.7476 | 300.2524 | 5 | 4.8 |
| family=bag;stmt=cte_filtered_notmat;threshold=20 | valid | 68.7032 | 73.5282 | 102.4163 | 5 | 13.6 |
| family=bag;stmt=cte_filtered_notmat;threshold=200 | valid | 90.0275 | 109.2058 | 125.7668 | 5 | 9.2 |
| family=bag;stmt=cte_filtered_notmat;threshold=2000 | valid | 68.0764 | 74.9934 | 100.8611 | 5 | 13.3 |
| family=bag;stmt=cte_filtered_notmat;threshold=20000 | valid | 86.1449 | 101.3577 | 125.3741 | 5 | 9.9 |
| family=bag;stmt=cte_joined_outer;threshold=100000 | valid | 124.1064 | 171.9248 | 227.6319 | 5 | 5.8 |
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 163.6813 | 168.7611 | 224.8246 | 5 | 5.9 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 71.3418 | 85.2483 | 103.6901 | 5 | 11.7 |
| family=bag;stmt=cte_joined_outer;threshold=200 | valid | 73.5611 | 95.4511 | 131.2638 | 5 | 10.5 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 73.2372 | 84.1131 | 118.2943 | 5 | 11.9 |
| family=bag;stmt=cte_joined_outer;threshold=20000 | valid | 82.3293 | 94.0174 | 126.7837 | 5 | 10.6 |
| family=bag;stmt=cte_joined_outer_mat;threshold=100000 | valid | 122.0386 | 172.9103 | 234.7764 | 5 | 5.8 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 129.9621 | 229.4704 | 287.6574 | 5 | 4.4 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 116.9008 | 214.1780 | 244.3066 | 5 | 4.7 |
| family=bag;stmt=cte_joined_outer_mat;threshold=200 | valid | 109.9610 | 136.4335 | 208.9593 | 5 | 7.3 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 103.3746 | 144.2880 | 203.4999 | 5 | 6.9 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20000 | valid | 135.3789 | 226.3498 | 266.6322 | 5 | 4.4 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=100000 | valid | 117.6387 | 128.2294 | 197.8817 | 5 | 7.8 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 148.7049 | 184.6368 | 288.9512 | 5 | 5.4 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000;position=last | valid | 218.6182 | 220.6716 | 386.8064 | 5 | 4.5 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 68.5515 | 73.8782 | 118.9590 | 5 | 13.5 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=200 | valid | 68.2955 | 79.0972 | 111.3629 | 5 | 12.6 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 70.1814 | 74.7956 | 88.1819 | 5 | 13.4 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20000 | valid | 92.6970 | 120.7276 | 130.5851 | 5 | 8.3 |
| family=bag;stmt=filtered_derived;threshold=100000 | valid | 122.5975 | 177.5349 | 195.1923 | 5 | 5.6 |
| family=bag;stmt=filtered_derived;threshold=180000 | valid | 191.5185 | 243.6141 | 261.2045 | 5 | 4.1 |
| family=bag;stmt=filtered_derived;threshold=20 | valid | 77.6780 | 93.5598 | 102.0028 | 5 | 10.7 |
| family=bag;stmt=filtered_derived;threshold=200 | valid | 74.2698 | 113.0063 | 129.5529 | 5 | 8.8 |
| family=bag;stmt=filtered_derived;threshold=2000 | valid | 81.2586 | 121.5485 | 135.9960 | 5 | 8.2 |
| family=bag;stmt=filtered_derived;threshold=20000 | valid | 85.3835 | 105.5479 | 126.4462 | 5 | 9.5 |
| family=bag;stmt=join_on;threshold=100000 | valid | 111.0307 | 145.1905 | 256.0550 | 5 | 6.9 |
| family=bag;stmt=join_on;threshold=180000 | valid | 154.9164 | 183.7657 | 219.1742 | 5 | 5.4 |
| family=bag;stmt=join_on;threshold=20 | valid | 71.6769 | 110.1639 | 123.7760 | 5 | 9.1 |
| family=bag;stmt=join_on;threshold=200 | valid | 71.5399 | 83.4269 | 101.4177 | 5 | 12.0 |
| family=bag;stmt=join_on;threshold=2000 | valid | 69.7071 | 80.4195 | 110.8995 | 5 | 12.4 |
| family=bag;stmt=join_on;threshold=20000 | valid | 86.5853 | 98.6984 | 124.7305 | 5 | 10.1 |
| family=bag;stmt=join_where;threshold=100000 | valid | 125.7620 | 174.8875 | 197.7037 | 5 | 5.7 |
| family=bag;stmt=join_where;threshold=180000 | valid | 149.4693 | 225.4778 | 280.7417 | 5 | 4.4 |
| family=bag;stmt=join_where;threshold=20 | valid | 66.2338 | 87.1451 | 107.2095 | 5 | 11.5 |
| family=bag;stmt=join_where;threshold=200 | valid | 70.3212 | 103.9895 | 130.2247 | 5 | 9.6 |
| family=bag;stmt=join_where;threshold=2000 | valid | 76.8903 | 106.0996 | 138.7202 | 5 | 9.4 |
| family=bag;stmt=join_where;threshold=20000 | valid | 90.2295 | 101.1035 | 122.9608 | 5 | 9.9 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 69.4866 | 73.8573 | 106.6756 | 5 | 13.5 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 71.5513 | 93.9867 | 97.8497 | 5 | 10.6 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 70.3628 | 92.4769 | 102.8499 | 5 | 10.8 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 74.7271 | 85.9608 | 110.1768 | 5 | 11.6 |
| family=reuse;policy=default;refs=2;consumer=1% | valid | 78.2468 | 91.0910 | 132.1154 | 5 | 11.0 |
| family=reuse;policy=default;refs=2;consumer=100% | valid | 67.3567 | 76.7118 | 109.1070 | 5 | 13.0 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 73.3677 | 88.9721 | 98.6294 | 5 | 11.2 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 70.6944 | 81.1366 | 97.7936 | 5 | 12.3 |
| family=reuse;policy=default;refs=8;consumer=1% | valid | 70.7200 | 72.9649 | 82.4688 | 5 | 13.7 |
| family=reuse;policy=default;refs=8;consumer=100% | valid | 69.4868 | 70.5471 | 104.1173 | 5 | 14.2 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 69.5153 | 82.0594 | 107.8151 | 5 | 12.2 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 73.0100 | 78.3027 | 109.6774 | 5 | 12.8 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 78.4137 | 87.8327 | 113.5336 | 5 | 11.4 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 84.4420 | 102.5196 | 150.8040 | 5 | 9.8 |
| family=reuse;policy=mat;refs=2;consumer=1% | valid | 66.2634 | 94.1085 | 112.0365 | 5 | 10.6 |
| family=reuse;policy=mat;refs=2;consumer=100% | valid | 73.4105 | 100.8830 | 115.6991 | 5 | 9.9 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 65.2712 | 77.3296 | 100.8856 | 5 | 12.9 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 71.9200 | 89.5087 | 121.0299 | 5 | 11.2 |
| family=reuse;policy=mat;refs=8;consumer=1% | valid | 73.9229 | 90.8158 | 113.4584 | 5 | 11.0 |
| family=reuse;policy=mat;refs=8;consumer=100% | valid | 74.0099 | 76.8639 | 109.4005 | 5 | 13.0 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 956.1567 | 1052.2198 | 1181.9696 | 5 | 1.0 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 1016.0862 | 1160.0005 | 1358.3052 | 5 | 0.9 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 68.6997 | 77.4532 | 102.8394 | 5 | 12.9 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 73.8164 | 97.9613 | 126.3362 | 5 | 10.2 |
| family=reuse;policy=notmat;refs=2;consumer=1% | valid | 156.2173 | 204.0947 | 270.2503 | 5 | 4.9 |
| family=reuse;policy=notmat;refs=2;consumer=100% | valid | 139.5277 | 219.6598 | 295.0942 | 5 | 4.6 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 254.7947 | 359.5680 | 366.3627 | 5 | 2.8 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 187.9038 | 285.6017 | 367.0692 | 5 | 3.5 |
| family=reuse;policy=notmat;refs=8;consumer=1% | valid | 461.2738 | 621.2279 | 685.3946 | 5 | 1.6 |
| family=reuse;policy=notmat;refs=8;consumer=100% | valid | 555.0041 | 635.6290 | 803.0339 | 5 | 1.6 |

## Execution order

- order seed `42304840`, 5 interleaved rounds
  1. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  2. family=bag;stmt=cte_filtered;threshold=180000
  3. family=bag;stmt=cte_joined_outer;threshold=2000
  4. family=bag;stmt=cte_filtered_mat;threshold=200
  5. family=bag;stmt=cte_filtered_mat;threshold=20000
  6. family=bag;stmt=filtered_derived;threshold=20000
  7. family=reuse;policy=default;refs=8;consumer=1%
  8. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  9. family=bag;stmt=join_on;threshold=180000
  10. family=reuse;policy=mat;refs=2;consumer=100%
  11. family=bag;stmt=cte_filtered;threshold=20000
  12. family=bag;stmt=cte_filtered_mat;threshold=20
  13. family=bag;stmt=join_where;threshold=100000
  14. family=reuse;policy=mat;refs=4;consumer=100%
  15. family=bag;stmt=join_on;threshold=100000
  16. family=reuse;policy=notmat;refs=4;consumer=1%
  17. family=reuse;policy=default;refs=2;consumer=1%
  18. family=bag;stmt=join_where;threshold=20000
  19. family=reuse;policy=notmat;refs=4;consumer=100%
  20. family=reuse;policy=mat;refs=1;consumer=100%
  21. family=bag;stmt=comma_join;threshold=2000
  22. family=reuse;policy=notmat;refs=16;consumer=100%
  23. family=bag;stmt=comma_join;threshold=20000
  24. family=bag;stmt=cte_filtered_notmat;threshold=20
  25. family=bag;stmt=join_where;threshold=20
  26. family=reuse;policy=default;refs=16;consumer=1%
  27. family=bag;stmt=cte_filtered_notmat;threshold=2000
  28. family=reuse;policy=mat;refs=4;consumer=1%
  29. family=bag;stmt=cte_joined_outer_notmat;threshold=100000
  30. family=reuse;policy=default;refs=4;consumer=1%
  31. family=reuse;policy=mat;refs=1;consumer=1%
  32. family=bag;stmt=cte_joined_outer_mat;threshold=200
  33. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  34. family=bag;stmt=cte_joined_outer_mat;threshold=100000
  35. family=reuse;policy=default;refs=2;consumer=100%
  36. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  37. family=reuse;policy=notmat;refs=8;consumer=100%
  38. family=reuse;policy=notmat;refs=8;consumer=1%
  39. family=reuse;policy=default;refs=4;consumer=100%
  40. family=bag;stmt=cte_filtered;threshold=100000
  41. family=bag;stmt=join_where;threshold=2000
  42. family=bag;stmt=cte_filtered;threshold=200
  43. family=reuse;policy=mat;refs=8;consumer=100%
  44. family=bag;stmt=comma_join;threshold=200
  45. family=bag;stmt=comma_join;threshold=20
  46. family=reuse;policy=default;refs=1;consumer=1%
  47. family=bag;stmt=cte_filtered_notmat;threshold=20000
  48. family=bag;stmt=cte_filtered_notmat;threshold=180000
  49. family=bag;stmt=cte_joined_outer;threshold=20
  50. family=bag;stmt=join_on;threshold=200
  51. family=bag;stmt=cte_joined_outer;threshold=20000
  52. family=bag;stmt=cte_filtered_notmat;threshold=100000
  53. family=reuse;policy=notmat;refs=1;consumer=1%
  54. family=bag;stmt=cte_filtered;threshold=2000
  55. family=bag;stmt=join_where;threshold=200
  56. family=reuse;policy=mat;refs=2;consumer=1%
  57. family=bag;stmt=cte_joined_outer;threshold=180000
  58. family=bag;stmt=cte_filtered_mat;threshold=180000
  59. family=bag;stmt=filtered_derived;threshold=2000
  60. family=bag;stmt=join_on;threshold=20
  61. family=bag;stmt=cte_filtered;threshold=20
  62. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  63. family=bag;stmt=cte_joined_outer_notmat;threshold=20000
  64. family=bag;stmt=filtered_derived;threshold=20
  65. family=bag;stmt=comma_join;threshold=100000
  66. family=bag;stmt=cte_joined_outer_mat;threshold=20
  67. family=bag;stmt=filtered_derived;threshold=100000
  68. family=reuse;policy=notmat;refs=1;consumer=100%
  69. family=bag;stmt=cte_joined_outer;threshold=200
  70. family=reuse;policy=default;refs=16;consumer=100%
  71. family=bag;stmt=filtered_derived;threshold=180000
  72. family=bag;stmt=cte_filtered_notmat;threshold=200
  73. family=bag;stmt=filtered_derived;threshold=200
  74. family=reuse;policy=notmat;refs=2;consumer=1%
  75. family=bag;stmt=cte_joined_outer;threshold=100000
  76. family=bag;stmt=cte_joined_outer_mat;threshold=20000
  77. family=reuse;policy=mat;refs=8;consumer=1%
  78. family=reuse;policy=mat;refs=16;consumer=1%
  79. family=bag;stmt=join_on;threshold=2000
  80. family=bag;stmt=join_where;threshold=180000
  81. family=reuse;policy=default;refs=1;consumer=100%
  82. family=reuse;policy=mat;refs=16;consumer=100%
  83. family=bag;stmt=cte_joined_outer_notmat;threshold=200
  84. family=reuse;policy=notmat;refs=2;consumer=100%
  85. family=bag;stmt=cte_filtered_mat;threshold=100000
  86. family=reuse;policy=notmat;refs=16;consumer=1%
  87. family=bag;stmt=join_on;threshold=20000
  88. family=bag;stmt=comma_join;threshold=180000
  89. family=bag;stmt=cte_filtered_mat;threshold=2000
  90. family=reuse;policy=default;refs=8;consumer=100%
  91. family=bag;stmt=cte_joined_outer_notmat;threshold=180000;position=last

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
| family=bag;stmt=comma_join;threshold=100000 | j06_comma_join | valid | `7f5e173da6e3176e` |  |
| family=bag;stmt=comma_join;threshold=180000 | j06_comma_join | valid | `758e482b2564df8f` |  |
| family=bag;stmt=comma_join;threshold=20 | j06_comma_join | valid | `20022d6cc43ee384` |  |
| family=bag;stmt=comma_join;threshold=200 | j06_comma_join | valid | `49fe175df2fec984` |  |
| family=bag;stmt=comma_join;threshold=2000 | j06_comma_join | valid | `e671f56c6e14637c` |  |
| family=bag;stmt=comma_join;threshold=20000 | j06_comma_join | valid | `d406b110dcf20a71` |  |
| family=bag;stmt=cte_filtered;threshold=100000 | j04_cte_filtered | valid | `c1bdf1a37b13c5ea` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=180000 | j04_cte_filtered | valid | `269a93712efb2046` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=20 | j04_cte_filtered | valid | `411f4fecd9c07342` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=200 | j04_cte_filtered | valid | `e5ef39e40c3efcf7` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=2000 | j04_cte_filtered | valid | `c4d3d7452e00aaeb` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=20000 | j04_cte_filtered | valid | `d3ff3f3bcec55cd0` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_mat;threshold=100000 | j04m_cte_filtered_mat | valid | `615e54a520cba74c` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=180000 | j04m_cte_filtered_mat | valid | `b691fd64a18a6c94` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=20 | j04m_cte_filtered_mat | valid | `34dfd062fead2f60` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=200 | j04m_cte_filtered_mat | valid | `df7e538ba26f899c` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=2000 | j04m_cte_filtered_mat | valid | `31f0fd1988835561` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=20000 | j04m_cte_filtered_mat | valid | `9b9a979b703d1d97` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_notmat;threshold=100000 | j04n_cte_filtered_notmat | valid | `d3bccf86663c508c` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=180000 | j04n_cte_filtered_notmat | valid | `d71da201abadc1bd` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=20 | j04n_cte_filtered_notmat | valid | `8b488e63bd5da547` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=200 | j04n_cte_filtered_notmat | valid | `2826fe5125c83cb6` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=2000 | j04n_cte_filtered_notmat | valid | `21b43a4cbb8bde77` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=20000 | j04n_cte_filtered_notmat | valid | `173f1b70d604336a` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=100000 | j05_cte_joined_outer | valid | `6a1bed6c1d520320` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=180000 | j05_cte_joined_outer | valid | `37cd1ca3e7d56c3a` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20 | j05_cte_joined_outer | valid | `13698fd12c6a58cf` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=200 | j05_cte_joined_outer | valid | `bfed798066d0017d` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=2000 | j05_cte_joined_outer | valid | `bc037a3f0c76168c` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20000 | j05_cte_joined_outer | valid | `889aaf70b36d1130` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_mat;threshold=100000 | j05m_cte_joined_outer_mat | valid | `c7e8493ff802070c` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | j05m_cte_joined_outer_mat | valid | `a92f0aa1f2d110f5` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | j05m_cte_joined_outer_mat | valid | `7e9805024af445f8` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=200 | j05m_cte_joined_outer_mat | valid | `13e1be8831a71333` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | j05m_cte_joined_outer_mat | valid | `2530a42421df920a` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20000 | j05m_cte_joined_outer_mat | valid | `a041e57b7882bdb4` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_notmat;threshold=100000 | j05n_cte_joined_outer_notmat | valid | `3584888d564591e7` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | j05n_cte_joined_outer_notmat | valid | `cd8d71d813f33225` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000;position=last | j05n_cte_joined_outer_notmat | valid | `70f633d21acce040` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | j05n_cte_joined_outer_notmat | valid | `48badd672d0c958a` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=200 | j05n_cte_joined_outer_notmat | valid | `d4d0ebd7ce145857` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | j05n_cte_joined_outer_notmat | valid | `8e8cc1598befb814` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20000 | j05n_cte_joined_outer_notmat | valid | `2a69594ccb932dc7` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=filtered_derived;threshold=100000 | j03_filtered_derived | valid | `c021269890a29a7f` |  |
| family=bag;stmt=filtered_derived;threshold=180000 | j03_filtered_derived | valid | `951262930ac8a856` |  |
| family=bag;stmt=filtered_derived;threshold=20 | j03_filtered_derived | valid | `5d9d96682e50c6e6` |  |
| family=bag;stmt=filtered_derived;threshold=200 | j03_filtered_derived | valid | `e243657c0b82a0be` |  |
| family=bag;stmt=filtered_derived;threshold=2000 | j03_filtered_derived | valid | `bc388ab11c492a69` |  |
| family=bag;stmt=filtered_derived;threshold=20000 | j03_filtered_derived | valid | `160af61f6233eb40` |  |
| family=bag;stmt=join_on;threshold=100000 | j02_join_on | valid | `590f83aee8e9f5a1` |  |
| family=bag;stmt=join_on;threshold=180000 | j02_join_on | valid | `a33ad2edd53b5874` |  |
| family=bag;stmt=join_on;threshold=20 | j02_join_on | valid | `1999e5b468ac7c6c` |  |
| family=bag;stmt=join_on;threshold=200 | j02_join_on | valid | `b0654b4f7f9027ac` |  |
| family=bag;stmt=join_on;threshold=2000 | j02_join_on | valid | `5d95140a6a11629c` |  |
| family=bag;stmt=join_on;threshold=20000 | j02_join_on | valid | `210161fac1a56428` |  |
| family=bag;stmt=join_where;threshold=100000 | j01_join_where | valid | `89a66c1b6e891494` |  |
| family=bag;stmt=join_where;threshold=180000 | j01_join_where | valid | `4110afa8bd4a428b` |  |
| family=bag;stmt=join_where;threshold=20 | j01_join_where | valid | `200973dab9bfdee0` |  |
| family=bag;stmt=join_where;threshold=200 | j01_join_where | valid | `db1b1b74fa90854a` |  |
| family=bag;stmt=join_where;threshold=2000 | j01_join_where | valid | `9eecd60ceace2193` |  |
| family=bag;stmt=join_where;threshold=20000 | j01_join_where | valid | `353d04b0ea0cba61` |  |
| family=reuse;policy=default;refs=16;consumer=1% | r_default_n16_c1 | valid | `e9949a76d592b697` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=16;consumer=100% | r_default_n16_c100 | valid | `77b5f8d3514276a6` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=1;consumer=1% | r_default_n1_c1 | valid | `70fa7780328362a2` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=1;consumer=100% | r_default_n1_c100 | valid | `433eb554db8690b9` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=2;consumer=1% | r_default_n2_c1 | valid | `666cdb8c3f3a5a86` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=2;consumer=100% | r_default_n2_c100 | valid | `96933d0bf8432dfb` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=1% | r_default_n4_c1 | valid | `7e0bf55a42ae0dfa` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=100% | r_default_n4_c100 | valid | `a4f852d3e5dc0cb1` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=8;consumer=1% | r_default_n8_c1 | valid | `bf9f5223db607fb2` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=8;consumer=100% | r_default_n8_c100 | valid | `ff67b98e8d61965a` | policy=default; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=1% | r_mat_n16_c1 | valid | `1732500e975be145` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=100% | r_mat_n16_c100 | valid | `a711475916bd2835` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=1% | r_mat_n1_c1 | valid | `bb52eb04f09e65bc` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=100% | r_mat_n1_c100 | valid | `12c0a0a668dcc12a` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=2;consumer=1% | r_mat_n2_c1 | valid | `3283f00b5edf0012` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=2;consumer=100% | r_mat_n2_c100 | valid | `1a8ef14e201df112` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=1% | r_mat_n4_c1 | valid | `cfe3fb80362a9c8a` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=100% | r_mat_n4_c100 | valid | `7aa8c8c2a4b3205e` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=8;consumer=1% | r_mat_n8_c1 | valid | `8a92857e216ae7f5` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=8;consumer=100% | r_mat_n8_c100 | valid | `eb735dea667440af` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=notmat;refs=16;consumer=1% | r_notmat_n16_c1 | valid | `b6c1e248463c5b00` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `69df98b52c51f084` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `4388b1e94492f68b` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `5420d77bf93dd781` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=2;consumer=1% | r_notmat_n2_c1 | valid | `cf6b13768743afa8` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=2;consumer=100% | r_notmat_n2_c100 | valid | `6cfb5c00796dc850` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `3888327d375b83b7` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `ff9aec201bee741d` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=8;consumer=1% | r_notmat_n8_c1 | valid | `e58023205d5af8b3` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=8;consumer=100% | r_notmat_n8_c100 | valid | `7f8ee940110a4016` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, DIST, DEBUG, FORMAT TEXT` supported on yugabyte
- one host, single client: no independent hardware, real-network or equal-total-budget comparison
- no concurrent writers: the isolated read contract only; a matched blended read/write phase is a recorded gap
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
