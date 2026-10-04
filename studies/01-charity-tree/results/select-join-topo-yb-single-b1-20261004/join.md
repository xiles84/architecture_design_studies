# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 30 cells, 28 valid, engine `yugabyte`, topology `single`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-topo-yb-single-b1-20261004`
- run tag: `run/study-01/select-join-topo-yb-single-b1-20261004`
- environment: `host-zenbook-ux5406sa`
- repo commit: `dcba783b6e5aebbe8f2833c2fff03464f8508385`
- engine version: `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit`
- database engine image: `docker.io/yugabytedb/yugabyte:2025.2.6.0-b111` @ `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`
- benchmark client image: `localhost/charitybench-select:select-join-topo-yb-single-b1-20261004` @ `d8db2fd98a7bc00fa2a5561ed01373926c9d950a8d93f0ae9ba4d53be68d6e38`
- inputs digest: `845d12e5f5d79207662ff5dfd3d2a1253f0bb283a056b02348622431919582e4`
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
| yb-single | 477 | 2 |

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
| 20 (0.01%) | 61.3433 | missing | missing | missing | missing | missing | 62.6658 | 110.0308 | 65.6848 | missing |
| 2000 (1%) | 66.5045 | missing | missing | missing | missing | missing | 64.4065 | 95.6429 | 73.2589 | missing |
| 180000 (90%) | 130.1092 | missing | missing | missing | missing | missing | 136.7490 | 135.1424 | 134.4626 | missing |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | cte_joined_outer vs join_where | 61.3433 | 62.6658 | 1.0216 | false |
| 20 | cte_joined_outer_mat vs join_where | 61.3433 | 110.0308 | 1.7937 | false |
| 20 | cte_joined_outer_notmat vs join_where | 61.3433 | 65.6848 | 1.0708 | false |
| 2000 | cte_joined_outer vs join_where | 66.5045 | 64.4065 | 0.9685 | false |
| 2000 | cte_joined_outer_mat vs join_where | 66.5045 | 95.6429 | 1.4381 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 66.5045 | 73.2589 | 1.1016 | false |
| 180000 | cte_joined_outer vs join_where | 130.1092 | 136.7490 | 1.0510 | false |
| 180000 | cte_joined_outer_mat vs join_where | 130.1092 | 135.1424 | 1.0387 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 130.1092 | 134.4626 | 1.0335 | false |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_joined_outer_mat vs cte_joined_outer | 62.6658 | 110.0308 | 1.7558 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 62.6658 | 65.6848 | 1.0482 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 64.4065 | 95.6429 | 1.4850 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 64.4065 | 73.2589 | 1.1374 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 136.7490 | 135.1424 | 0.9883 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 136.7490 | 134.4626 | 0.9833 | false | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 66.6494 | 1 | `e38abba3d2c9ac0b` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 71.4869 | 1 | `7bb00fc1c1694465` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | missing | | | | |
| default | 2 | 100% | missing | | | | |
| default | 4 | 1% | 68.7220 | 4 | `be0bd93b1f12bdfd` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 62.8882 | 4 | `143ff249136971f4` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | missing | | | | |
| default | 8 | 100% | missing | | | | |
| default | 16 | 1% | 67.8303 | 16 | `a04add73c583111e` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 64.4196 | 16 | `a194a7a62cb6fbe7` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 58.6422 | 1 | `3b458f7fc231e07c` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 62.2717 | 1 | `357ecf18fa510a67` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | missing | | | | |
| mat | 2 | 100% | missing | | | | |
| mat | 4 | 1% | 62.0378 | 4 | `4456fe5fbdd04e20` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 57.0270 | 4 | `f765ed24027c99fd` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | missing | | | | |
| mat | 8 | 100% | missing | | | | |
| mat | 16 | 1% | 69.8752 | 16 | `819313793317e035` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 64.6047 | 16 | `039b64ba5faf6ab9` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 73.7302 | 1 | `b22bb96a35b56433` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 64.6614 | 1 | `b667771b185f76dc` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | missing | | | | |
| notmat | 2 | 100% | missing | | | | |
| notmat | 4 | 1% | error | 0 | `e3b0c44298fc1c14` |  | error |
| notmat | 4 | 100% | 186.7993 | 4 | `8aa8e8c486451c12` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | missing | | | | |
| notmat | 8 | 100% | missing | | | | |
| notmat | 16 | 1% | error | 0 | `e3b0c44298fc1c14` |  | error |
| notmat | 16 | 100% | 709.6863 | 16 | `a731727f3a2d0154` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 0.8799 | false |
| notmat | 1 | 1% | 1.1062 | false |
| mat | 1 | 100% | 0.8711 | false |
| notmat | 1 | 100% | 0.9045 | false |
| mat | 4 | 1% | 0.9027 | false |
| notmat | 4 | 1% | N/A | N/A |
| mat | 4 | 100% | 0.9068 | false |
| notmat | 4 | 100% | 2.9703 | false |
| mat | 16 | 1% | 1.0301 | false |
| notmat | 16 | 1% | N/A | N/A |
| mat | 16 | 100% | 1.0029 | false |
| notmat | 16 | 100% | 11.0166 | false |

## Latency samples and derived throughput

_The reader is sequential, so `queries/s` is 1000/median_ms, a derived rate and not a concurrent-capacity test. min/max are the five sample extremes._

| axis | status | min ms | median ms | max ms | samples | derived queries/s |
|---|---|---|---|---|---|---|
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 123.2597 | 136.7490 | 165.0253 | 5 | 7.3 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 56.0099 | 62.6658 | 68.0892 | 5 | 16.0 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 54.6004 | 64.4065 | 69.3886 | 5 | 15.5 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 132.1565 | 135.1424 | 156.1595 | 5 | 7.4 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 97.5845 | 110.0308 | 131.6069 | 5 | 9.1 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 88.7677 | 95.6429 | 117.5157 | 5 | 10.5 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 127.3516 | 134.4626 | 140.8769 | 5 | 7.4 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 60.9443 | 65.6848 | 69.5699 | 5 | 15.2 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 57.5987 | 73.2589 | 75.0338 | 5 | 13.7 |
| family=bag;stmt=join_where;threshold=180000 | valid | 128.0087 | 130.1092 | 160.7129 | 5 | 7.7 |
| family=bag;stmt=join_where;threshold=20 | valid | 59.7027 | 61.3433 | 65.1409 | 5 | 16.3 |
| family=bag;stmt=join_where;threshold=2000 | valid | 62.6292 | 66.5045 | 79.3759 | 5 | 15.0 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 61.5184 | 67.8303 | 71.7809 | 5 | 14.7 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 63.7189 | 64.4196 | 67.9428 | 5 | 15.5 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 60.1502 | 66.6494 | 76.7236 | 5 | 15.0 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 62.9406 | 71.4869 | 82.2089 | 5 | 14.0 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 64.6847 | 68.7220 | 71.0949 | 5 | 14.6 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 62.1608 | 62.8882 | 67.8957 | 5 | 15.9 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 64.6015 | 69.8752 | 105.8845 | 5 | 14.3 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 58.7750 | 64.6047 | 69.5054 | 5 | 15.5 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 54.0709 | 58.6422 | 73.8502 | 5 | 17.1 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 58.2588 | 62.2717 | 66.7590 | 5 | 16.1 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 57.6249 | 62.0378 | 85.7961 | 5 | 16.1 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 55.7483 | 57.0270 | 64.1030 | 5 | 17.5 |
| family=reuse;policy=notmat;refs=16;consumer=1% | error | N/A | N/A | N/A | 0 | N/A |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 683.5254 | 709.6863 | 737.0626 | 5 | 1.4 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 59.7183 | 73.7302 | 98.4119 | 5 | 13.6 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 58.5049 | 64.6614 | 77.3383 | 5 | 15.5 |
| family=reuse;policy=notmat;refs=4;consumer=1% | error | N/A | N/A | N/A | 0 | N/A |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 175.4668 | 186.7993 | 243.3340 | 5 | 5.4 |

## Matched blended workload

- phase: `blended`; target offered rate: 10 updates/s
- offered 324, committed 324, failed 0 over 32.334 s
- offered rate 10.020/s, completed rate 10.020/s
- reader operations overlapping a committed write: 155 of 200
- per key: 324 keys, min 1, max 1 per key
- ledger: 324 keys reconciled, 0 mismatch
- update latency: avg 10.600 ms, min 1.655, max 94.724 over 324

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42251436`, 5 interleaved rounds
  1. family=bag;stmt=join_where;threshold=2000
  2. family=reuse;policy=notmat;refs=16;consumer=100%
  3. family=reuse;policy=default;refs=4;consumer=1%
  4. family=reuse;policy=default;refs=1;consumer=1%
  5. family=reuse;policy=notmat;refs=4;consumer=1%
  6. family=reuse;policy=mat;refs=1;consumer=1%
  7. family=reuse;policy=default;refs=4;consumer=100%
  8. family=reuse;policy=mat;refs=4;consumer=1%
  9. family=reuse;policy=mat;refs=1;consumer=100%
  10. family=bag;stmt=cte_joined_outer;threshold=20
  11. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  12. family=reuse;policy=notmat;refs=4;consumer=100%
  13. family=reuse;policy=default;refs=16;consumer=1%
  14. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  15. family=reuse;policy=default;refs=1;consumer=100%
  16. family=bag;stmt=cte_joined_outer;threshold=180000
  17. family=reuse;policy=mat;refs=16;consumer=1%
  18. family=reuse;policy=notmat;refs=1;consumer=1%
  19. family=bag;stmt=cte_joined_outer_mat;threshold=20
  20. family=bag;stmt=cte_joined_outer;threshold=2000
  21. family=bag;stmt=join_where;threshold=20
  22. family=reuse;policy=default;refs=16;consumer=100%
  23. family=reuse;policy=notmat;refs=1;consumer=100%
  24. family=reuse;policy=notmat;refs=16;consumer=1%
  25. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  26. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  27. family=bag;stmt=join_where;threshold=180000
  28. family=reuse;policy=mat;refs=4;consumer=100%
  29. family=reuse;policy=mat;refs=16;consumer=100%
  30. family=bag;stmt=cte_joined_outer_mat;threshold=180000

## Live negative controls

| control | pretended contract | detected | rows | note |
|---|---|---|---|---|
| x01_optional_where | sel_join_optional | true | 1 |  |
| x02_join_limit_before | sel_join_bag | true | 7 |  |

## Excluded cells

- `family=reuse;policy=notmat;refs=4;consumer=1%` (r_notmat_n4_c1) — **error** — ERROR: Restart read required (SQLSTATE 40001)
- `family=reuse;policy=notmat;refs=16;consumer=1%` (r_notmat_n16_c1) — **error** — ERROR: Restart read required (SQLSTATE 40001)

A cell that is invalid, unsupported or timed out is never ranked.

## Plan signatures

| axis | statement | status | structure signature | policy note |
|---|---|---|---|---|
| family=bag;stmt=cte_joined_outer;threshold=180000 | j05_cte_joined_outer | valid | `6df68f5308c8f0b4` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20 | j05_cte_joined_outer | valid | `7bb0f82e1ffb1f3b` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=2000 | j05_cte_joined_outer | valid | `6cd4b4af5914659c` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | j05m_cte_joined_outer_mat | valid | `611f06667c99e28b` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | j05m_cte_joined_outer_mat | valid | `ea5471bdbd47f2cb` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | j05m_cte_joined_outer_mat | valid | `456d9abf0ef0581b` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | j05n_cte_joined_outer_notmat | valid | `b77f099457d76985` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | j05n_cte_joined_outer_notmat | valid | `2d03b34eeaa42479` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | j05n_cte_joined_outer_notmat | valid | `57eed0f080b59537` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=join_where;threshold=180000 | j01_join_where | valid | `4c19191a908be55d` |  |
| family=bag;stmt=join_where;threshold=20 | j01_join_where | valid | `60dd82583c6bf1c4` |  |
| family=bag;stmt=join_where;threshold=2000 | j01_join_where | valid | `b839ae18398a3fc5` |  |
| family=reuse;policy=default;refs=16;consumer=1% | r_default_n16_c1 | valid | `a04add73c583111e` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=16;consumer=100% | r_default_n16_c100 | valid | `a194a7a62cb6fbe7` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=1;consumer=1% | r_default_n1_c1 | valid | `e38abba3d2c9ac0b` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=1;consumer=100% | r_default_n1_c100 | valid | `7bb00fc1c1694465` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=4;consumer=1% | r_default_n4_c1 | valid | `be0bd93b1f12bdfd` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=100% | r_default_n4_c100 | valid | `143ff249136971f4` | policy=default; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=1% | r_mat_n16_c1 | valid | `819313793317e035` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=100% | r_mat_n16_c100 | valid | `039b64ba5faf6ab9` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=1% | r_mat_n1_c1 | valid | `3b458f7fc231e07c` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=100% | r_mat_n1_c100 | valid | `357ecf18fa510a67` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=1% | r_mat_n4_c1 | valid | `4456fe5fbdd04e20` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=100% | r_mat_n4_c100 | valid | `f765ed24027c99fd` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=notmat;refs=16;consumer=1% | r_notmat_n16_c1 | error | `e3b0c44298fc1c14` |  |
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `a731727f3a2d0154` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `b22bb96a35b56433` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `b667771b185f76dc` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | error | `e3b0c44298fc1c14` |  |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `8aa8e8c486451c12` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, DIST, DEBUG, FORMAT TEXT` supported on yugabyte
- one host, single client: no independent hardware, real-network or equal-total-budget comparison
- no concurrent writers: the isolated read contract only; a matched blended read/write phase is a recorded gap
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
