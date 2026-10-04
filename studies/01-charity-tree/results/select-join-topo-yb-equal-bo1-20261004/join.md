# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 30 cells, 30 valid, engine `yugabyte`, topology `single`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-topo-yb-equal-bo1-20261004`
- run tag: `run/study-01/select-join-topo-yb-equal-bo1-20261004`
- environment: `host-zenbook-ux5406sa`
- repo commit: `455b777fc02e015efdd9e54d25d3aa3343eb206f`
- engine version: `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit`
- database engine image: `docker.io/yugabytedb/yugabyte:2025.2.6.0-b111` @ `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`
- benchmark client image: `localhost/charitybench-select:select-join-topo-yb-equal-bo1-20261004` @ `79bc753ff7e57e102982a323939c43e68e9b3664887683d3a5d5bf072650f4c1`
- inputs digest: `cdf3ae2cdb1e171b126e143c23c9d4b8e2a63054f0a4ef1f73240469590d96e4`
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
| yb-single | 491 | 2 |

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
| 20 (0.01%) | 55.7718 | missing | missing | missing | missing | missing | 60.7990 | 107.1603 | 68.1762 | missing |
| 2000 (1%) | 64.3848 | missing | missing | missing | missing | missing | 64.6677 | 96.3998 | 60.9297 | missing |
| 180000 (90%) | 130.7126 | missing | missing | missing | missing | missing | 125.5740 | 130.1365 | 132.8978 | missing |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | cte_joined_outer vs join_where | 55.7718 | 60.7990 | 1.0901 | false |
| 20 | cte_joined_outer_mat vs join_where | 55.7718 | 107.1603 | 1.9214 | false |
| 20 | cte_joined_outer_notmat vs join_where | 55.7718 | 68.1762 | 1.2224 | false |
| 2000 | cte_joined_outer vs join_where | 64.3848 | 64.6677 | 1.0044 | false |
| 2000 | cte_joined_outer_mat vs join_where | 64.3848 | 96.3998 | 1.4972 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 64.3848 | 60.9297 | 0.9463 | false |
| 180000 | cte_joined_outer vs join_where | 130.7126 | 125.5740 | 0.9607 | false |
| 180000 | cte_joined_outer_mat vs join_where | 130.7126 | 130.1365 | 0.9956 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 130.7126 | 132.8978 | 1.0167 | false |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_joined_outer_mat vs cte_joined_outer | 60.7990 | 107.1603 | 1.7625 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 60.7990 | 68.1762 | 1.1213 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 64.6677 | 96.3998 | 1.4907 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 64.6677 | 60.9297 | 0.9422 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 125.5740 | 130.1365 | 1.0363 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 125.5740 | 132.8978 | 1.0583 | false | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 64.3880 | 1 | `084d37cface1ff36` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 69.2072 | 1 | `4054668c214ff07f` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | missing | | | | |
| default | 2 | 100% | missing | | | | |
| default | 4 | 1% | 62.7746 | 4 | `2586b0c680533264` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 57.7688 | 4 | `c5df8b123464f86b` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | missing | | | | |
| default | 8 | 100% | missing | | | | |
| default | 16 | 1% | 58.4504 | 16 | `269df972d7566da6` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 62.3052 | 16 | `44424af1450a368b` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 58.1759 | 1 | `104027dc2ca75704` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 61.7149 | 1 | `6b5e05911a415a2c` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | missing | | | | |
| mat | 2 | 100% | missing | | | | |
| mat | 4 | 1% | 65.0029 | 4 | `bcda715eadfedd1c` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 59.7647 | 4 | `7de06adfd27cbe48` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | missing | | | | |
| mat | 8 | 100% | missing | | | | |
| mat | 16 | 1% | 60.3551 | 16 | `2bf2aa38923bbc50` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 66.1339 | 16 | `4b150978b1ee4a44` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 56.1354 | 1 | `82f8bb3de06a4ac5` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 64.1528 | 1 | `45617973c1c85613` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | missing | | | | |
| notmat | 2 | 100% | missing | | | | |
| notmat | 4 | 1% | 145.9511 | 4 | `5d31642c484d4ad3` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 142.9396 | 4 | `8703b5aca576f404` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | missing | | | | |
| notmat | 8 | 100% | missing | | | | |
| notmat | 16 | 1% | 441.8267 | 16 | `45eea5059fd5c996` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 444.6848 | 16 | `d78346c79c8c6ac3` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 0.9035 | false |
| notmat | 1 | 1% | 0.8718 | false |
| mat | 1 | 100% | 0.8917 | false |
| notmat | 1 | 100% | 0.9270 | false |
| mat | 4 | 1% | 1.0355 | false |
| notmat | 4 | 1% | 2.3250 | false |
| mat | 4 | 100% | 1.0345 | false |
| notmat | 4 | 100% | 2.4743 | false |
| mat | 16 | 1% | 1.0326 | false |
| notmat | 16 | 1% | 7.5590 | false |
| mat | 16 | 100% | 1.0615 | false |
| notmat | 16 | 100% | 7.1372 | false |

## Latency samples and derived throughput

_The reader is sequential, so `queries/s` is 1000/median_ms, a derived rate and not a concurrent-capacity test. min/max are the five sample extremes._

| axis | status | min ms | median ms | max ms | samples | derived queries/s |
|---|---|---|---|---|---|---|
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 123.5750 | 125.5740 | 133.5012 | 5 | 8.0 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 52.1196 | 60.7990 | 68.1581 | 5 | 16.4 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 56.1900 | 64.6677 | 69.9356 | 5 | 15.5 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 128.4741 | 130.1365 | 139.2376 | 5 | 7.7 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 93.2900 | 107.1603 | 112.3024 | 5 | 9.3 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 91.7562 | 96.3998 | 122.6152 | 5 | 10.4 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 116.7946 | 132.8978 | 137.9029 | 5 | 7.5 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 53.4965 | 68.1762 | 74.2840 | 5 | 14.7 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 54.2067 | 60.9297 | 64.0312 | 5 | 16.4 |
| family=bag;stmt=join_where;threshold=180000 | valid | 128.0648 | 130.7126 | 148.2096 | 5 | 7.7 |
| family=bag;stmt=join_where;threshold=20 | valid | 52.0400 | 55.7718 | 67.0758 | 5 | 17.9 |
| family=bag;stmt=join_where;threshold=2000 | valid | 57.6885 | 64.3848 | 69.7310 | 5 | 15.5 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 57.5224 | 58.4504 | 67.1283 | 5 | 17.1 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 52.6888 | 62.3052 | 70.7476 | 5 | 16.1 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 56.4874 | 64.3880 | 71.6221 | 5 | 15.5 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 55.5352 | 69.2072 | 71.9260 | 5 | 14.4 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 56.0061 | 62.7746 | 67.1200 | 5 | 15.9 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 49.6540 | 57.7688 | 58.5686 | 5 | 17.3 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 55.8806 | 60.3551 | 64.1469 | 5 | 16.6 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 56.8763 | 66.1339 | 68.5815 | 5 | 15.1 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 47.6084 | 58.1759 | 63.0534 | 5 | 17.2 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 60.6254 | 61.7149 | 71.4640 | 5 | 16.2 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 51.1312 | 65.0029 | 66.5345 | 5 | 15.4 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 57.2457 | 59.7647 | 62.4809 | 5 | 16.7 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 435.5079 | 441.8267 | 594.2187 | 5 | 2.3 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 433.4950 | 444.6848 | 459.7109 | 5 | 2.2 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 52.6322 | 56.1354 | 68.9310 | 5 | 17.8 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 62.2578 | 64.1528 | 71.9616 | 5 | 15.6 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 134.8188 | 145.9511 | 152.9314 | 5 | 6.9 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 134.6330 | 142.9396 | 157.8021 | 5 | 7.0 |

## Matched blended workload

- phase: `blended`; target offered rate: 10 updates/s
- offered 309, committed 309, failed 0 over 30.808 s
- offered rate 10.030/s, completed rate 10.030/s
- reader operations overlapping a committed write: 160 of 214
- reader retries of transient conflicts: 4
- per key: 309 keys, min 1, max 1 per key
- ledger: 309 keys reconciled, 0 mismatch
- update latency: avg 3.625 ms, min 1.687, max 63.806 over 309

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42139479`, 5 interleaved rounds
  1. family=bag;stmt=join_where;threshold=2000
  2. family=reuse;policy=mat;refs=4;consumer=1%
  3. family=reuse;policy=mat;refs=1;consumer=1%
  4. family=bag;stmt=cte_joined_outer;threshold=180000
  5. family=reuse;policy=mat;refs=16;consumer=1%
  6. family=reuse;policy=default;refs=4;consumer=100%
  7. family=reuse;policy=mat;refs=16;consumer=100%
  8. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  9. family=reuse;policy=notmat;refs=4;consumer=100%
  10. family=bag;stmt=cte_joined_outer;threshold=2000
  11. family=reuse;policy=default;refs=1;consumer=100%
  12. family=reuse;policy=notmat;refs=1;consumer=100%
  13. family=reuse;policy=default;refs=1;consumer=1%
  14. family=bag;stmt=join_where;threshold=20
  15. family=reuse;policy=default;refs=4;consumer=1%
  16. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  17. family=reuse;policy=default;refs=16;consumer=1%
  18. family=bag;stmt=cte_joined_outer;threshold=20
  19. family=reuse;policy=mat;refs=1;consumer=100%
  20. family=reuse;policy=default;refs=16;consumer=100%
  21. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  22. family=reuse;policy=notmat;refs=4;consumer=1%
  23. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  24. family=bag;stmt=cte_joined_outer_mat;threshold=20
  25. family=reuse;policy=notmat;refs=16;consumer=1%
  26. family=bag;stmt=join_where;threshold=180000
  27. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  28. family=reuse;policy=mat;refs=4;consumer=100%
  29. family=reuse;policy=notmat;refs=1;consumer=1%
  30. family=reuse;policy=notmat;refs=16;consumer=100%

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
| family=bag;stmt=cte_joined_outer;threshold=180000 | j05_cte_joined_outer | valid | `555b709754de86b6` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20 | j05_cte_joined_outer | valid | `05970441897cc50b` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=2000 | j05_cte_joined_outer | valid | `69621a0fdb22de54` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | j05m_cte_joined_outer_mat | valid | `ffe93a68f4e029a8` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | j05m_cte_joined_outer_mat | valid | `b9c436e0b5359acf` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | j05m_cte_joined_outer_mat | valid | `286a24f01a683a05` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | j05n_cte_joined_outer_notmat | valid | `9be34e6bacbaaff8` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | j05n_cte_joined_outer_notmat | valid | `fd3a4be996c5a3e0` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | j05n_cte_joined_outer_notmat | valid | `dbe8e610a3a583ed` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=join_where;threshold=180000 | j01_join_where | valid | `2a802f194fdce6ff` |  |
| family=bag;stmt=join_where;threshold=20 | j01_join_where | valid | `4ec951b302a52c15` |  |
| family=bag;stmt=join_where;threshold=2000 | j01_join_where | valid | `9b7c27647cbd613a` |  |
| family=reuse;policy=default;refs=16;consumer=1% | r_default_n16_c1 | valid | `269df972d7566da6` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=16;consumer=100% | r_default_n16_c100 | valid | `44424af1450a368b` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=1;consumer=1% | r_default_n1_c1 | valid | `084d37cface1ff36` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=1;consumer=100% | r_default_n1_c100 | valid | `4054668c214ff07f` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=4;consumer=1% | r_default_n4_c1 | valid | `2586b0c680533264` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=100% | r_default_n4_c100 | valid | `c5df8b123464f86b` | policy=default; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=1% | r_mat_n16_c1 | valid | `2bf2aa38923bbc50` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=100% | r_mat_n16_c100 | valid | `4b150978b1ee4a44` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=1% | r_mat_n1_c1 | valid | `104027dc2ca75704` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=100% | r_mat_n1_c100 | valid | `6b5e05911a415a2c` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=1% | r_mat_n4_c1 | valid | `bcda715eadfedd1c` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=100% | r_mat_n4_c100 | valid | `7de06adfd27cbe48` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=notmat;refs=16;consumer=1% | r_notmat_n16_c1 | valid | `45eea5059fd5c996` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `d78346c79c8c6ac3` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `82f8bb3de06a4ac5` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `45617973c1c85613` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `5d31642c484d4ad3` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `8703b5aca576f404` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, DIST, DEBUG, FORMAT TEXT` supported on yugabyte
- one host, single client: no independent hardware, real-network or equal-total-budget comparison
- no concurrent writers: the isolated read contract only; a matched blended read/write phase is a recorded gap
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
