# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 30 cells, 30 valid, engine `yugabyte`, topology `single`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-topo-yb-single-bo2-20261004`
- run tag: `run/study-01/select-join-topo-yb-single-bo2-20261004`
- environment: `host-zenbook-ux5406sa`
- repo commit: `b8a0ff816fa169824edcf9ef89f2fd8502643df4`
- engine version: `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit`
- database engine image: `docker.io/yugabytedb/yugabyte:2025.2.6.0-b111` @ `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`
- benchmark client image: `localhost/charitybench-select:select-join-topo-yb-single-bo2-20261004` @ `79bc753ff7e57e102982a323939c43e68e9b3664887683d3a5d5bf072650f4c1`
- inputs digest: `b60f7ee32ceb5836bdcf155de7ee09b9320e8f42aca183c160e6e3b8a8477524`
- dataset fingerprint: `fcf26e08f377ec51`
- client CPU quota: `100000 100000`, memory: `2147483648`
- setting `client_workers`: 1 sequential reader; 0 writers
- setting `connection_pool`: 2 per endpoint = reader/setup + reserve
- setting `customers`: 2000
- setting `dataset`: study01-person-donation-charity
- setting `distributed_rpc`: not parsed; raw EXPLAIN retained (coverage gap)
- setting `endpoint_mode`: single endpoint
- setting `endpoints`: yb-single
- setting `fact_count`: 200000
- setting `index_set`: child rank index and parent-key index, primary keys; no per-spelling index
- setting `phase`: blended
- setting `placement`: noncolocated
- setting `prepared_plan_policy`: pgx simple protocol; no prepared statements
- setting `scale`: cte-domain
- setting `selectivity_mapping`: rank <= threshold; rank is a deterministic shuffle of 1..fact_count
- setting `statement_budget_ms`: 120000
- setting `statistics`: ANALYZE after deterministic load
- setting `thresholds`: [20 2000 180000]
- setting `topology`: single
- setting `trials_per_cell`: 5
- setting `whole_run_count`: 1; no replicated winner in this file
- setting `writer_keys`: 1000
- setting `writer_target_rate_per_sec`: 10

### Endpoint distribution

A cluster arm must rotate its reader across nodes; a single seed endpoint would not be balanced cluster access.

| endpoint | operations | connections at end |
|---|---|---|
| yb-single | 494 | 2 |

## Fixture and selectivity mapping

- arm: `study01-person-donation-charity`
- fact rows: 200000
- parent rows: 2000
- rank thresholds: [20 2000 180000]

The predicate is `rank <= threshold`, where `rank` is a deterministic shuffle of 1..200000; `rank` is a labelled proxy for the teaching selective predicate, not the same question. The `cte`/`cte-small` arm is the synthetic customers/facts microfixture; the `cte-domain` arm is Study 01's person/donation fixture with the person's charity projected and the views binding it to this catalogue.

## Per-engine feasibility evidence

- fact rows 200000, customers 2000, thresholds [20 2000 180000]
- expensive reuse token supported by the engine: yes

| threshold | rows with rank <= threshold |
|---|---|
| 180000 | 180000 |
| 20 | 20 |
| 2000 | 2000 |

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
| 20 (0.01%) | 57.2546 | missing | missing | missing | missing | missing | 57.8894 | 94.9424 | 56.4210 | missing |
| 2000 (1%) | 64.0663 | missing | missing | missing | missing | missing | 59.0295 | 92.3406 | 59.0918 | missing |
| 180000 (90%) | 140.8515 | missing | missing | missing | missing | missing | 125.5266 | 124.6410 | 158.4339 | missing |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | cte_joined_outer vs join_where | 57.2546 | 57.8894 | 1.0111 | false |
| 20 | cte_joined_outer_mat vs join_where | 57.2546 | 94.9424 | 1.6582 | false |
| 20 | cte_joined_outer_notmat vs join_where | 57.2546 | 56.4210 | 0.9854 | false |
| 2000 | cte_joined_outer vs join_where | 64.0663 | 59.0295 | 0.9214 | false |
| 2000 | cte_joined_outer_mat vs join_where | 64.0663 | 92.3406 | 1.4413 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 64.0663 | 59.0918 | 0.9224 | false |
| 180000 | cte_joined_outer vs join_where | 140.8515 | 125.5266 | 0.8912 | false |
| 180000 | cte_joined_outer_mat vs join_where | 140.8515 | 124.6410 | 0.8849 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 140.8515 | 158.4339 | 1.1248 | false |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_joined_outer_mat vs cte_joined_outer | 57.8894 | 94.9424 | 1.6401 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 57.8894 | 56.4210 | 0.9746 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 59.0295 | 92.3406 | 1.5643 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 59.0295 | 59.0918 | 1.0011 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 125.5266 | 124.6410 | 0.9929 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 125.5266 | 158.4339 | 1.2622 | false | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 58.3898 | 1 | `1802cc33f4b6cf11` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 60.0259 | 1 | `b3e01ce83a94bf40` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | missing | | | | |
| default | 2 | 100% | missing | | | | |
| default | 4 | 1% | 60.5417 | 4 | `500c4554626756d6` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 57.5371 | 4 | `d24d6356c4e46ff9` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | missing | | | | |
| default | 8 | 100% | missing | | | | |
| default | 16 | 1% | 56.5063 | 16 | `089f54087b3593d7` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 58.9286 | 16 | `c7c42e3f447b58d9` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 58.5865 | 1 | `777a3bb1a087c774` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 55.6929 | 1 | `f6424de37922b8d2` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | missing | | | | |
| mat | 2 | 100% | missing | | | | |
| mat | 4 | 1% | 55.2996 | 4 | `9ca8d343f1185ae2` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 59.5685 | 4 | `47c776a1a0737f09` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | missing | | | | |
| mat | 8 | 100% | missing | | | | |
| mat | 16 | 1% | 60.0821 | 16 | `d6f05ff7a76062cf` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 57.9146 | 16 | `9676ac467e0207b2` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 56.3947 | 1 | `364adac21747c656` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 64.4059 | 1 | `d061b5a744840281` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | missing | | | | |
| notmat | 2 | 100% | missing | | | | |
| notmat | 4 | 1% | 189.9891 | 4 | `ff0d8a1fe1c6d930` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 169.2588 | 4 | `c035267c9a7e3cf6` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | missing | | | | |
| notmat | 8 | 100% | missing | | | | |
| notmat | 16 | 1% | 843.0020 | 16 | `5f1a6ede3ef7cfd1` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 653.0132 | 16 | `077799500ca528ee` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 1.0034 | false |
| notmat | 1 | 1% | 0.9658 | false |
| mat | 1 | 100% | 0.9278 | false |
| notmat | 1 | 100% | 1.0730 | false |
| mat | 4 | 1% | 0.9134 | false |
| notmat | 4 | 1% | 3.1382 | false |
| mat | 4 | 100% | 1.0353 | false |
| notmat | 4 | 100% | 2.9417 | false |
| mat | 16 | 1% | 1.0633 | false |
| notmat | 16 | 1% | 14.9187 | false |
| mat | 16 | 100% | 0.9828 | false |
| notmat | 16 | 100% | 11.0814 | false |

## Latency samples and derived throughput

_The reader is sequential, so `queries/s` is 1000/median_ms, a derived rate and not a concurrent-capacity test. min/max are the five sample extremes._

| axis | status | min ms | median ms | max ms | samples | derived queries/s |
|---|---|---|---|---|---|---|
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 112.3566 | 125.5266 | 135.4642 | 5 | 8.0 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 54.0213 | 57.8894 | 63.4993 | 5 | 17.3 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 57.2411 | 59.0295 | 63.5602 | 5 | 16.9 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 118.0100 | 124.6410 | 140.9334 | 5 | 8.0 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 85.6629 | 94.9424 | 99.9431 | 5 | 10.5 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 90.4185 | 92.3406 | 112.7302 | 5 | 10.8 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 125.8390 | 158.4339 | 167.6135 | 5 | 6.3 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 54.3545 | 56.4210 | 67.4089 | 5 | 17.7 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 53.4713 | 59.0918 | 61.5758 | 5 | 16.9 |
| family=bag;stmt=join_where;threshold=180000 | valid | 123.7762 | 140.8515 | 162.7798 | 5 | 7.1 |
| family=bag;stmt=join_where;threshold=20 | valid | 55.3910 | 57.2546 | 59.2947 | 5 | 17.5 |
| family=bag;stmt=join_where;threshold=2000 | valid | 54.5167 | 64.0663 | 68.7488 | 5 | 15.6 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 55.0970 | 56.5063 | 62.5975 | 5 | 17.7 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 54.5371 | 58.9286 | 61.1039 | 5 | 17.0 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 50.3548 | 58.3898 | 62.0638 | 5 | 17.1 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 58.0839 | 60.0259 | 60.4722 | 5 | 16.7 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 53.7787 | 60.5417 | 102.0868 | 5 | 16.5 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 53.1432 | 57.5371 | 63.1215 | 5 | 17.4 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 56.0151 | 60.0821 | 63.8946 | 5 | 16.6 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 55.9375 | 57.9146 | 64.2912 | 5 | 17.3 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 55.6122 | 58.5865 | 63.9063 | 5 | 17.1 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 53.6480 | 55.6929 | 61.6062 | 5 | 18.0 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 53.0282 | 55.2996 | 59.8260 | 5 | 18.1 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 56.8273 | 59.5685 | 64.5488 | 5 | 16.8 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 609.6680 | 843.0020 | 915.2231 | 5 | 1.2 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 647.9726 | 653.0132 | 684.9767 | 5 | 1.5 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 53.2885 | 56.3947 | 59.8850 | 5 | 17.7 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 57.3438 | 64.4059 | 72.8139 | 5 | 15.5 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 162.3191 | 189.9891 | 395.8618 | 5 | 5.3 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 158.8751 | 169.2588 | 180.9938 | 5 | 5.9 |

## Matched blended workload

- phase: `blended`; target offered rate: 10 updates/s
- offered 366, committed 366, failed 0 over 36.515 s
- offered rate 10.023/s, completed rate 10.023/s
- reader operations overlapping a committed write: 161 of 217
- reader retries of transient conflicts: 7
- per key: 366 keys, min 1, max 1 per key
- ledger: 366 keys reconciled, 0 mismatch
- update latency: avg 3.372 ms, min 1.558, max 56.516 over 366

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42367464`, 5 interleaved rounds
  1. family=reuse;policy=notmat;refs=1;consumer=100%
  2. family=reuse;policy=notmat;refs=1;consumer=1%
  3. family=reuse;policy=notmat;refs=16;consumer=1%
  4. family=bag;stmt=join_where;threshold=180000
  5. family=bag;stmt=cte_joined_outer;threshold=20
  6. family=bag;stmt=join_where;threshold=2000
  7. family=reuse;policy=notmat;refs=4;consumer=1%
  8. family=reuse;policy=default;refs=4;consumer=1%
  9. family=bag;stmt=cte_joined_outer;threshold=180000
  10. family=reuse;policy=mat;refs=4;consumer=100%
  11. family=reuse;policy=notmat;refs=4;consumer=100%
  12. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  13. family=reuse;policy=default;refs=16;consumer=100%
  14. family=reuse;policy=notmat;refs=16;consumer=100%
  15. family=reuse;policy=mat;refs=4;consumer=1%
  16. family=reuse;policy=default;refs=1;consumer=1%
  17. family=reuse;policy=mat;refs=1;consumer=1%
  18. family=bag;stmt=join_where;threshold=20
  19. family=reuse;policy=default;refs=4;consumer=100%
  20. family=reuse;policy=mat;refs=16;consumer=100%
  21. family=reuse;policy=mat;refs=16;consumer=1%
  22. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  23. family=reuse;policy=default;refs=16;consumer=1%
  24. family=reuse;policy=mat;refs=1;consumer=100%
  25. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  26. family=reuse;policy=default;refs=1;consumer=100%
  27. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  28. family=bag;stmt=cte_joined_outer_mat;threshold=20
  29. family=bag;stmt=cte_joined_outer;threshold=2000
  30. family=bag;stmt=cte_joined_outer_notmat;threshold=20

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
| family=bag;stmt=cte_joined_outer;threshold=180000 | j05_cte_joined_outer | valid | `6ff691c91f51090f` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20 | j05_cte_joined_outer | valid | `d3185dd2dacd4e38` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=2000 | j05_cte_joined_outer | valid | `e01e1e4f5ab89d5d` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | j05m_cte_joined_outer_mat | valid | `e53ba341b5a22f73` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | j05m_cte_joined_outer_mat | valid | `bc96e9fdebbdcf25` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | j05m_cte_joined_outer_mat | valid | `8aa8d08317e2cda6` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | j05n_cte_joined_outer_notmat | valid | `03b32664024accd1` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | j05n_cte_joined_outer_notmat | valid | `aae41dd9a96eb3e3` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | j05n_cte_joined_outer_notmat | valid | `0a48d16fd8ce8903` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=join_where;threshold=180000 | j01_join_where | valid | `a8abe92fcce4f2e6` |  |
| family=bag;stmt=join_where;threshold=20 | j01_join_where | valid | `5b1476de7bee563f` |  |
| family=bag;stmt=join_where;threshold=2000 | j01_join_where | valid | `dbc9b8bf896000f7` |  |
| family=reuse;policy=default;refs=16;consumer=1% | r_default_n16_c1 | valid | `089f54087b3593d7` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=16;consumer=100% | r_default_n16_c100 | valid | `c7c42e3f447b58d9` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=1;consumer=1% | r_default_n1_c1 | valid | `1802cc33f4b6cf11` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=1;consumer=100% | r_default_n1_c100 | valid | `b3e01ce83a94bf40` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=4;consumer=1% | r_default_n4_c1 | valid | `500c4554626756d6` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=100% | r_default_n4_c100 | valid | `d24d6356c4e46ff9` | policy=default; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=1% | r_mat_n16_c1 | valid | `d6f05ff7a76062cf` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=100% | r_mat_n16_c100 | valid | `9676ac467e0207b2` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=1% | r_mat_n1_c1 | valid | `777a3bb1a087c774` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=100% | r_mat_n1_c100 | valid | `f6424de37922b8d2` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=1% | r_mat_n4_c1 | valid | `9ca8d343f1185ae2` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=100% | r_mat_n4_c100 | valid | `47c776a1a0737f09` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=notmat;refs=16;consumer=1% | r_notmat_n16_c1 | valid | `5f1a6ede3ef7cfd1` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `077799500ca528ee` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `364adac21747c656` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `d061b5a744840281` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `ff0d8a1fe1c6d930` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `c035267c9a7e3cf6` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, DIST, DEBUG, FORMAT TEXT` supported on yugabyte
- one host, single client: no independent hardware, real-network or equal-total-budget comparison
- no concurrent writers: the isolated read contract only; a matched blended read/write phase is a recorded gap
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
