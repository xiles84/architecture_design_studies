# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 30 cells, 30 valid, engine `yugabyte`, topology `single`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-topo-yb-equal-i2-20261004`
- run tag: `run/study-01/select-join-topo-yb-equal-i2-20261004`
- environment: `host-zenbook-ux5406sa`
- repo commit: `335542f10685ca626b2f4cb2b086f0ad67e620e0`
- engine version: `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit`
- database engine image: `docker.io/yugabytedb/yugabyte:2025.2.6.0-b111` @ `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`
- benchmark client image: `localhost/charitybench-select:select-join-topo-yb-equal-i2-20261004` @ `d8db2fd98a7bc00fa2a5561ed01373926c9d950a8d93f0ae9ba4d53be68d6e38`
- inputs digest: `e7b43712745fda908d2c8da27af236914129ee8e4b887487e39ab10e6c560b65`
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
- setting `phase`: isolated
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
| yb-single | 487 | 1 |

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
| 20 (0.01%) | 54.6696 | missing | missing | missing | missing | missing | 57.8967 | 91.4008 | 55.7310 | missing |
| 2000 (1%) | 58.2504 | missing | missing | missing | missing | missing | 56.1353 | 92.6476 | 56.1235 | missing |
| 180000 (90%) | 116.6264 | missing | missing | missing | missing | missing | 124.9935 | 119.4418 | 119.8374 | missing |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | cte_joined_outer vs join_where | 54.6696 | 57.8967 | 1.0590 | false |
| 20 | cte_joined_outer_mat vs join_where | 54.6696 | 91.4008 | 1.6719 | false |
| 20 | cte_joined_outer_notmat vs join_where | 54.6696 | 55.7310 | 1.0194 | false |
| 2000 | cte_joined_outer vs join_where | 58.2504 | 56.1353 | 0.9637 | false |
| 2000 | cte_joined_outer_mat vs join_where | 58.2504 | 92.6476 | 1.5905 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 58.2504 | 56.1235 | 0.9635 | false |
| 180000 | cte_joined_outer vs join_where | 116.6264 | 124.9935 | 1.0717 | false |
| 180000 | cte_joined_outer_mat vs join_where | 116.6264 | 119.4418 | 1.0241 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 116.6264 | 119.8374 | 1.0275 | false |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_joined_outer_mat vs cte_joined_outer | 57.8967 | 91.4008 | 1.5787 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 57.8967 | 55.7310 | 0.9626 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 56.1353 | 92.6476 | 1.6504 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 56.1353 | 56.1235 | 0.9998 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 124.9935 | 119.4418 | 0.9556 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 124.9935 | 119.8374 | 0.9587 | false | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 54.5555 | 1 | `b0d52ea78f4cf146` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 54.9889 | 1 | `7ac7ad104a953cb5` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | missing | | | | |
| default | 2 | 100% | missing | | | | |
| default | 4 | 1% | 55.4522 | 4 | `349b95ccdf93ed2b` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 54.4170 | 4 | `a58a35bde09a9928` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | missing | | | | |
| default | 8 | 100% | missing | | | | |
| default | 16 | 1% | 53.8209 | 16 | `2e2c961ea8752a5a` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 57.9235 | 16 | `ef435064c6f91f09` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 59.7615 | 1 | `83cbfdcebdac3a96` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 59.5332 | 1 | `e7555e76cf040e54` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | missing | | | | |
| mat | 2 | 100% | missing | | | | |
| mat | 4 | 1% | 55.0674 | 4 | `67607ed75ece74c2` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 57.1266 | 4 | `f2d4f1ac7faaa618` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | missing | | | | |
| mat | 8 | 100% | missing | | | | |
| mat | 16 | 1% | 56.8127 | 16 | `9dc416be430b6143` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 56.8694 | 16 | `930a97be841d8927` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 57.2648 | 1 | `e34c18ba06b65fde` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 58.8432 | 1 | `27052b5af7945c85` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | missing | | | | |
| notmat | 2 | 100% | missing | | | | |
| notmat | 4 | 1% | 125.9317 | 4 | `8825df427b1c7e61` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 127.9154 | 4 | `c22f8eb061644f56` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | missing | | | | |
| notmat | 8 | 100% | missing | | | | |
| notmat | 16 | 1% | 382.2118 | 16 | `fa58cbd41a91dd29` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 388.9698 | 16 | `8c5c57246f4153b1` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 1.0954 | false |
| notmat | 1 | 1% | 1.0497 | false |
| mat | 1 | 100% | 1.0826 | false |
| notmat | 1 | 100% | 1.0701 | false |
| mat | 4 | 1% | 0.9931 | false |
| notmat | 4 | 1% | 2.2710 | false |
| mat | 4 | 100% | 1.0498 | false |
| notmat | 4 | 100% | 2.3507 | false |
| mat | 16 | 1% | 1.0556 | false |
| notmat | 16 | 1% | 7.1015 | false |
| mat | 16 | 100% | 0.9818 | false |
| notmat | 16 | 100% | 6.7152 | false |

## Latency samples and derived throughput

_The reader is sequential, so `queries/s` is 1000/median_ms, a derived rate and not a concurrent-capacity test. min/max are the five sample extremes._

| axis | status | min ms | median ms | max ms | samples | derived queries/s |
|---|---|---|---|---|---|---|
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 122.0047 | 124.9935 | 131.6242 | 5 | 8.0 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 55.6521 | 57.8967 | 61.1849 | 5 | 17.3 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 54.7546 | 56.1353 | 56.4307 | 5 | 17.8 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 115.6572 | 119.4418 | 128.9973 | 5 | 8.4 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 89.1663 | 91.4008 | 98.2967 | 5 | 10.9 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 87.0693 | 92.6476 | 102.3741 | 5 | 10.8 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 114.6706 | 119.8374 | 125.3922 | 5 | 8.3 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 52.5723 | 55.7310 | 58.2213 | 5 | 17.9 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 54.3545 | 56.1235 | 61.2074 | 5 | 17.8 |
| family=bag;stmt=join_where;threshold=180000 | valid | 113.6301 | 116.6264 | 124.8158 | 5 | 8.6 |
| family=bag;stmt=join_where;threshold=20 | valid | 51.7619 | 54.6696 | 60.7571 | 5 | 18.3 |
| family=bag;stmt=join_where;threshold=2000 | valid | 57.4505 | 58.2504 | 60.6655 | 5 | 17.2 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 53.1112 | 53.8209 | 57.4800 | 5 | 18.6 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 56.4618 | 57.9235 | 61.1597 | 5 | 17.3 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 52.2046 | 54.5555 | 58.3483 | 5 | 18.3 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 52.2378 | 54.9889 | 57.5479 | 5 | 18.2 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 52.6503 | 55.4522 | 57.6871 | 5 | 18.0 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 52.7646 | 54.4170 | 58.4860 | 5 | 18.4 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 55.3378 | 56.8127 | 57.5914 | 5 | 17.6 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 54.3918 | 56.8694 | 61.7672 | 5 | 17.6 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 55.0247 | 59.7615 | 62.6167 | 5 | 16.7 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 56.4963 | 59.5332 | 64.5274 | 5 | 16.8 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 53.9955 | 55.0674 | 62.0998 | 5 | 18.2 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 54.9936 | 57.1266 | 62.5907 | 5 | 17.5 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 369.0131 | 382.2118 | 395.3233 | 5 | 2.6 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 380.0343 | 388.9698 | 399.1309 | 5 | 2.6 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 55.4519 | 57.2648 | 59.4510 | 5 | 17.5 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 57.3367 | 58.8432 | 63.3220 | 5 | 17.0 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 121.8954 | 125.9317 | 132.1375 | 5 | 7.9 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 124.3400 | 127.9154 | 128.9819 | 5 | 7.8 |

## Matched blended workload

- phase: `isolated`; target offered rate: 10 updates/s
- offered 0, committed 0, failed 0 over 0.000 s
- offered rate 0.000/s, completed rate 0.000/s
- per key: no keys touched
- ledger: 0 keys reconciled, 0 mismatch

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42446758`, 5 interleaved rounds
  1. family=reuse;policy=notmat;refs=1;consumer=1%
  2. family=reuse;policy=mat;refs=16;consumer=1%
  3. family=reuse;policy=default;refs=1;consumer=1%
  4. family=reuse;policy=default;refs=16;consumer=100%
  5. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  6. family=reuse;policy=default;refs=4;consumer=100%
  7. family=reuse;policy=notmat;refs=16;consumer=100%
  8. family=bag;stmt=cte_joined_outer;threshold=20
  9. family=bag;stmt=cte_joined_outer;threshold=180000
  10. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  11. family=bag;stmt=cte_joined_outer;threshold=2000
  12. family=bag;stmt=join_where;threshold=180000
  13. family=reuse;policy=default;refs=1;consumer=100%
  14. family=bag;stmt=join_where;threshold=20
  15. family=bag;stmt=cte_joined_outer_mat;threshold=20
  16. family=reuse;policy=default;refs=4;consumer=1%
  17. family=reuse;policy=mat;refs=16;consumer=100%
  18. family=reuse;policy=notmat;refs=16;consumer=1%
  19. family=reuse;policy=mat;refs=1;consumer=1%
  20. family=reuse;policy=mat;refs=1;consumer=100%
  21. family=reuse;policy=mat;refs=4;consumer=1%
  22. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  23. family=reuse;policy=default;refs=16;consumer=1%
  24. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  25. family=reuse;policy=notmat;refs=4;consumer=100%
  26. family=reuse;policy=notmat;refs=1;consumer=100%
  27. family=reuse;policy=notmat;refs=4;consumer=1%
  28. family=bag;stmt=join_where;threshold=2000
  29. family=reuse;policy=mat;refs=4;consumer=100%
  30. family=bag;stmt=cte_joined_outer_notmat;threshold=180000

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
| family=bag;stmt=cte_joined_outer;threshold=180000 | j05_cte_joined_outer | valid | `270f4b5ec89de8e1` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20 | j05_cte_joined_outer | valid | `5ce3b97592611adb` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=2000 | j05_cte_joined_outer | valid | `8e44065a360ac1ff` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | j05m_cte_joined_outer_mat | valid | `f512f2563f8448bb` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | j05m_cte_joined_outer_mat | valid | `1f23a54547928cca` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | j05m_cte_joined_outer_mat | valid | `f104f2ebcfe283c5` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | j05n_cte_joined_outer_notmat | valid | `2ef9e04d01fd303e` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | j05n_cte_joined_outer_notmat | valid | `9278a4047a0cc0a1` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | j05n_cte_joined_outer_notmat | valid | `babb9d0cdad9469c` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=join_where;threshold=180000 | j01_join_where | valid | `98a3ea05ac8580aa` |  |
| family=bag;stmt=join_where;threshold=20 | j01_join_where | valid | `6e92196a7b326ecd` |  |
| family=bag;stmt=join_where;threshold=2000 | j01_join_where | valid | `9f7d585132699983` |  |
| family=reuse;policy=default;refs=16;consumer=1% | r_default_n16_c1 | valid | `2e2c961ea8752a5a` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=16;consumer=100% | r_default_n16_c100 | valid | `ef435064c6f91f09` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=1;consumer=1% | r_default_n1_c1 | valid | `b0d52ea78f4cf146` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=1;consumer=100% | r_default_n1_c100 | valid | `7ac7ad104a953cb5` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=4;consumer=1% | r_default_n4_c1 | valid | `349b95ccdf93ed2b` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=100% | r_default_n4_c100 | valid | `a58a35bde09a9928` | policy=default; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=1% | r_mat_n16_c1 | valid | `9dc416be430b6143` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=100% | r_mat_n16_c100 | valid | `930a97be841d8927` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=1% | r_mat_n1_c1 | valid | `83cbfdcebdac3a96` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=100% | r_mat_n1_c100 | valid | `e7555e76cf040e54` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=1% | r_mat_n4_c1 | valid | `67607ed75ece74c2` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=100% | r_mat_n4_c100 | valid | `f2d4f1ac7faaa618` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=notmat;refs=16;consumer=1% | r_notmat_n16_c1 | valid | `fa58cbd41a91dd29` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `8c5c57246f4153b1` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `e34c18ba06b65fde` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `27052b5af7945c85` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `8825df427b1c7e61` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `c22f8eb061644f56` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, DIST, DEBUG, FORMAT TEXT` supported on yugabyte
- one host, single client: no independent hardware, real-network or equal-total-budget comparison
- no concurrent writers: the isolated read contract only; a matched blended read/write phase is a recorded gap
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
