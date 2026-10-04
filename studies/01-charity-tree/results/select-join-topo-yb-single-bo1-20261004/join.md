# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 30 cells, 30 valid, engine `yugabyte`, topology `single`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-topo-yb-single-bo1-20261004`
- run tag: `run/study-01/select-join-topo-yb-single-bo1-20261004`
- environment: `host-zenbook-ux5406sa`
- repo commit: `b8a0ff816fa169824edcf9ef89f2fd8502643df4`
- engine version: `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit`
- database engine image: `docker.io/yugabytedb/yugabyte:2025.2.6.0-b111` @ `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`
- benchmark client image: `localhost/charitybench-select:select-join-topo-yb-single-bo1-20261004` @ `79bc753ff7e57e102982a323939c43e68e9b3664887683d3a5d5bf072650f4c1`
- inputs digest: `f85692bc31bdead23b8e23b7e61bbcc7c8c4e63c9f3cc35d003808708e679a14`
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
| yb-single | 489 | 2 |

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
| 20 (0.01%) | 62.6396 | missing | missing | missing | missing | missing | 57.5677 | 108.0833 | 60.9876 | missing |
| 2000 (1%) | 59.6312 | missing | missing | missing | missing | missing | 61.7730 | 104.8964 | 60.0741 | missing |
| 180000 (90%) | 119.6240 | missing | missing | missing | missing | missing | 125.0343 | 125.7757 | 136.3596 | missing |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | cte_joined_outer vs join_where | 62.6396 | 57.5677 | 0.9190 | false |
| 20 | cte_joined_outer_mat vs join_where | 62.6396 | 108.0833 | 1.7255 | false |
| 20 | cte_joined_outer_notmat vs join_where | 62.6396 | 60.9876 | 0.9736 | false |
| 2000 | cte_joined_outer vs join_where | 59.6312 | 61.7730 | 1.0359 | false |
| 2000 | cte_joined_outer_mat vs join_where | 59.6312 | 104.8964 | 1.7591 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 59.6312 | 60.0741 | 1.0074 | false |
| 180000 | cte_joined_outer vs join_where | 119.6240 | 125.0343 | 1.0452 | false |
| 180000 | cte_joined_outer_mat vs join_where | 119.6240 | 125.7757 | 1.0514 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 119.6240 | 136.3596 | 1.1399 | false |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_joined_outer_mat vs cte_joined_outer | 57.5677 | 108.0833 | 1.8775 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 57.5677 | 60.9876 | 1.0594 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 61.7730 | 104.8964 | 1.6981 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 61.7730 | 60.0741 | 0.9725 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 125.0343 | 125.7757 | 1.0059 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 125.0343 | 136.3596 | 1.0906 | false | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 54.4408 | 1 | `06489355976acb98` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 61.9992 | 1 | `404620459ba78fe7` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | missing | | | | |
| default | 2 | 100% | missing | | | | |
| default | 4 | 1% | 56.8826 | 4 | `3ab270084b2fd518` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 57.4907 | 4 | `ede86c1d8934114d` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | missing | | | | |
| default | 8 | 100% | missing | | | | |
| default | 16 | 1% | 62.2506 | 16 | `2999200d1b3ef299` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 65.7374 | 16 | `94da61e5b1f38c3f` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 56.9979 | 1 | `f6df43823a7cba9c` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 54.4799 | 1 | `8da50e32df37a3c5` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | missing | | | | |
| mat | 2 | 100% | missing | | | | |
| mat | 4 | 1% | 59.9407 | 4 | `029f6540a5749725` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 57.1308 | 4 | `569b13100c66be40` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | missing | | | | |
| mat | 8 | 100% | missing | | | | |
| mat | 16 | 1% | 60.6356 | 16 | `172097f86edf1a59` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 58.7067 | 16 | `3e8cc414eeda9844` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 59.6673 | 1 | `32f937d2e31938ff` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 66.9390 | 1 | `5cd0bc040852a1c0` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | missing | | | | |
| notmat | 2 | 100% | missing | | | | |
| notmat | 4 | 1% | 329.3815 | 4 | `c92a2c212ae1b018` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 175.1936 | 4 | `e5f72d4298b36c41` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | missing | | | | |
| notmat | 8 | 100% | missing | | | | |
| notmat | 16 | 1% | 786.3588 | 16 | `2d4008794c6a2b14` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 672.4271 | 16 | `677fbd08e7aff1e8` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 1.0470 | false |
| notmat | 1 | 1% | 1.0960 | false |
| mat | 1 | 100% | 0.8787 | false |
| notmat | 1 | 100% | 1.0797 | false |
| mat | 4 | 1% | 1.0538 | false |
| notmat | 4 | 1% | 5.7906 | false |
| mat | 4 | 100% | 0.9937 | false |
| notmat | 4 | 100% | 3.0473 | false |
| mat | 16 | 1% | 0.9741 | false |
| notmat | 16 | 1% | 12.6321 | false |
| mat | 16 | 100% | 0.8930 | false |
| notmat | 16 | 100% | 10.2290 | false |

## Latency samples and derived throughput

_The reader is sequential, so `queries/s` is 1000/median_ms, a derived rate and not a concurrent-capacity test. min/max are the five sample extremes._

| axis | status | min ms | median ms | max ms | samples | derived queries/s |
|---|---|---|---|---|---|---|
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 108.8311 | 125.0343 | 130.2587 | 5 | 8.0 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 46.8576 | 57.5677 | 62.1325 | 5 | 17.4 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 57.7229 | 61.7730 | 65.4778 | 5 | 16.2 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 112.2815 | 125.7757 | 149.6623 | 5 | 8.0 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 89.2622 | 108.0833 | 131.9329 | 5 | 9.3 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 92.5577 | 104.8964 | 123.5366 | 5 | 9.5 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 110.0207 | 136.3596 | 138.8805 | 5 | 7.3 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 57.8680 | 60.9876 | 71.2542 | 5 | 16.4 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 49.2044 | 60.0741 | 64.9295 | 5 | 16.6 |
| family=bag;stmt=join_where;threshold=180000 | valid | 108.9602 | 119.6240 | 132.8001 | 5 | 8.4 |
| family=bag;stmt=join_where;threshold=20 | valid | 49.6898 | 62.6396 | 66.9136 | 5 | 16.0 |
| family=bag;stmt=join_where;threshold=2000 | valid | 54.9825 | 59.6312 | 104.4711 | 5 | 16.8 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 56.1440 | 62.2506 | 121.6891 | 5 | 16.1 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 55.2129 | 65.7374 | 109.9562 | 5 | 15.2 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 47.7752 | 54.4408 | 58.8880 | 5 | 18.4 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 55.2200 | 61.9992 | 62.9817 | 5 | 16.1 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 51.9681 | 56.8826 | 73.6417 | 5 | 17.6 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 53.9916 | 57.4907 | 62.5478 | 5 | 17.4 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 56.9249 | 60.6356 | 71.3293 | 5 | 16.5 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 54.2671 | 58.7067 | 78.4096 | 5 | 17.0 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 53.8699 | 56.9979 | 83.9982 | 5 | 17.5 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 52.4619 | 54.4799 | 63.2256 | 5 | 18.4 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 56.9000 | 59.9407 | 66.2395 | 5 | 16.7 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 50.0719 | 57.1308 | 63.3761 | 5 | 17.5 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 666.1053 | 786.3588 | 1023.2894 | 5 | 1.3 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 581.5066 | 672.4271 | 726.5584 | 5 | 1.5 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 53.7459 | 59.6673 | 63.3324 | 5 | 16.8 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 54.9897 | 66.9390 | 126.5398 | 5 | 14.9 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 160.8315 | 329.3815 | 496.0976 | 5 | 3.0 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 120.0643 | 175.1936 | 225.2271 | 5 | 5.7 |

## Matched blended workload

- phase: `blended`; target offered rate: 10 updates/s
- offered 351, committed 351, failed 0 over 35.020 s
- offered rate 10.023/s, completed rate 10.023/s
- reader operations overlapping a committed write: 161 of 212
- reader retries of transient conflicts: 2
- per key: 351 keys, min 1, max 1 per key
- ledger: 351 keys reconciled, 0 mismatch
- update latency: avg 14.782 ms, min 1.579, max 68.722 over 351

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42623771`, 5 interleaved rounds
  1. family=reuse;policy=default;refs=1;consumer=1%
  2. family=bag;stmt=cte_joined_outer;threshold=180000
  3. family=reuse;policy=notmat;refs=4;consumer=100%
  4. family=reuse;policy=notmat;refs=1;consumer=100%
  5. family=reuse;policy=mat;refs=4;consumer=100%
  6. family=reuse;policy=default;refs=4;consumer=100%
  7. family=bag;stmt=join_where;threshold=180000
  8. family=bag;stmt=cte_joined_outer;threshold=20
  9. family=reuse;policy=notmat;refs=16;consumer=100%
  10. family=reuse;policy=mat;refs=16;consumer=100%
  11. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  12. family=reuse;policy=mat;refs=1;consumer=1%
  13. family=reuse;policy=default;refs=16;consumer=100%
  14. family=reuse;policy=default;refs=16;consumer=1%
  15. family=bag;stmt=join_where;threshold=2000
  16. family=reuse;policy=mat;refs=16;consumer=1%
  17. family=reuse;policy=mat;refs=1;consumer=100%
  18. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  19. family=reuse;policy=mat;refs=4;consumer=1%
  20. family=reuse;policy=default;refs=4;consumer=1%
  21. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  22. family=reuse;policy=notmat;refs=4;consumer=1%
  23. family=reuse;policy=notmat;refs=16;consumer=1%
  24. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  25. family=bag;stmt=cte_joined_outer;threshold=2000
  26. family=bag;stmt=join_where;threshold=20
  27. family=bag;stmt=cte_joined_outer_mat;threshold=20
  28. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  29. family=reuse;policy=notmat;refs=1;consumer=1%
  30. family=reuse;policy=default;refs=1;consumer=100%

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
| family=bag;stmt=cte_joined_outer;threshold=180000 | j05_cte_joined_outer | valid | `dceb128ce3730e0f` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20 | j05_cte_joined_outer | valid | `ab4a806492dbecbf` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=2000 | j05_cte_joined_outer | valid | `7be0cba8c32218e5` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | j05m_cte_joined_outer_mat | valid | `a5ec68246decfd44` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | j05m_cte_joined_outer_mat | valid | `0bf4b6c61b15f743` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | j05m_cte_joined_outer_mat | valid | `8a00a0533b861d8d` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | j05n_cte_joined_outer_notmat | valid | `10080e0c7bc14343` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | j05n_cte_joined_outer_notmat | valid | `67e3b5fd437bd9fd` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | j05n_cte_joined_outer_notmat | valid | `483c62a6bca3ca0e` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=join_where;threshold=180000 | j01_join_where | valid | `8be2f4498324f29b` |  |
| family=bag;stmt=join_where;threshold=20 | j01_join_where | valid | `2bf194084f31e9e3` |  |
| family=bag;stmt=join_where;threshold=2000 | j01_join_where | valid | `4e8a9668bafce1c4` |  |
| family=reuse;policy=default;refs=16;consumer=1% | r_default_n16_c1 | valid | `2999200d1b3ef299` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=16;consumer=100% | r_default_n16_c100 | valid | `94da61e5b1f38c3f` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=1;consumer=1% | r_default_n1_c1 | valid | `06489355976acb98` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=1;consumer=100% | r_default_n1_c100 | valid | `404620459ba78fe7` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=4;consumer=1% | r_default_n4_c1 | valid | `3ab270084b2fd518` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=100% | r_default_n4_c100 | valid | `ede86c1d8934114d` | policy=default; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=1% | r_mat_n16_c1 | valid | `172097f86edf1a59` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=100% | r_mat_n16_c100 | valid | `3e8cc414eeda9844` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=1% | r_mat_n1_c1 | valid | `f6df43823a7cba9c` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=100% | r_mat_n1_c100 | valid | `8da50e32df37a3c5` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=1% | r_mat_n4_c1 | valid | `029f6540a5749725` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=100% | r_mat_n4_c100 | valid | `569b13100c66be40` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=notmat;refs=16;consumer=1% | r_notmat_n16_c1 | valid | `2d4008794c6a2b14` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `677fbd08e7aff1e8` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `32f937d2e31938ff` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `5cd0bc040852a1c0` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `c92a2c212ae1b018` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `e5f72d4298b36c41` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, DIST, DEBUG, FORMAT TEXT` supported on yugabyte
- one host, single client: no independent hardware, real-network or equal-total-budget comparison
- no concurrent writers: the isolated read contract only; a matched blended read/write phase is a recorded gap
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
