# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 30 cells, 30 valid, engine `yugabyte`, topology `single`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-topo-yb-single-bo3-20261004`
- run tag: `run/study-01/select-join-topo-yb-single-bo3-20261004`
- environment: `host-zenbook-ux5406sa`
- repo commit: `b8a0ff816fa169824edcf9ef89f2fd8502643df4`
- engine version: `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit`
- database engine image: `docker.io/yugabytedb/yugabyte:2025.2.6.0-b111` @ `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`
- benchmark client image: `localhost/charitybench-select:select-join-topo-yb-single-bo3-20261004` @ `79bc753ff7e57e102982a323939c43e68e9b3664887683d3a5d5bf072650f4c1`
- inputs digest: `a8af0c83b2cf2e55d0f421d98b9dd05daa9f293cc251d6efe2a455e70c199c0e`
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
| yb-single | 490 | 2 |

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
| 20 (0.01%) | 59.5696 | missing | missing | missing | missing | missing | 59.7128 | 96.5852 | 55.3569 | missing |
| 2000 (1%) | 59.0996 | missing | missing | missing | missing | missing | 60.3632 | 100.2182 | 63.4369 | missing |
| 180000 (90%) | 128.5681 | missing | missing | missing | missing | missing | 125.4940 | 122.8594 | 129.8621 | missing |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | cte_joined_outer vs join_where | 59.5696 | 59.7128 | 1.0024 | false |
| 20 | cte_joined_outer_mat vs join_where | 59.5696 | 96.5852 | 1.6214 | false |
| 20 | cte_joined_outer_notmat vs join_where | 59.5696 | 55.3569 | 0.9293 | false |
| 2000 | cte_joined_outer vs join_where | 59.0996 | 60.3632 | 1.0214 | false |
| 2000 | cte_joined_outer_mat vs join_where | 59.0996 | 100.2182 | 1.6958 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 59.0996 | 63.4369 | 1.0734 | false |
| 180000 | cte_joined_outer vs join_where | 128.5681 | 125.4940 | 0.9761 | false |
| 180000 | cte_joined_outer_mat vs join_where | 128.5681 | 122.8594 | 0.9556 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 128.5681 | 129.8621 | 1.0101 | false |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_joined_outer_mat vs cte_joined_outer | 59.7128 | 96.5852 | 1.6175 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 59.7128 | 55.3569 | 0.9271 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 60.3632 | 100.2182 | 1.6603 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 60.3632 | 63.4369 | 1.0509 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 125.4940 | 122.8594 | 0.9790 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 125.4940 | 129.8621 | 1.0348 | false | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 64.2076 | 1 | `48c3360dc26089a8` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 61.3180 | 1 | `dbca7388712c192f` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | missing | | | | |
| default | 2 | 100% | missing | | | | |
| default | 4 | 1% | 56.1946 | 4 | `d9774833bab3d450` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 59.1869 | 4 | `20f11895a76aedb7` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | missing | | | | |
| default | 8 | 100% | missing | | | | |
| default | 16 | 1% | 58.1751 | 16 | `e33b692d73e204a9` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 59.5917 | 16 | `e21e5033536eec65` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 57.2049 | 1 | `14a9958e48ea12f8` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 54.8806 | 1 | `323d121527000186` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | missing | | | | |
| mat | 2 | 100% | missing | | | | |
| mat | 4 | 1% | 53.8015 | 4 | `a0e8570238aec77b` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 58.5001 | 4 | `68009261e24c39bc` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | missing | | | | |
| mat | 8 | 100% | missing | | | | |
| mat | 16 | 1% | 63.1612 | 16 | `2297f4c9516b1656` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 63.7872 | 16 | `dd552acd262869aa` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 62.8064 | 1 | `d24dbbc49d287a7c` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 62.6936 | 1 | `7c0dc957717661fa` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | missing | | | | |
| notmat | 2 | 100% | missing | | | | |
| notmat | 4 | 1% | 200.8992 | 4 | `48948ec51ed0942c` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 187.6542 | 4 | `077e53594d575969` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | missing | | | | |
| notmat | 8 | 100% | missing | | | | |
| notmat | 16 | 1% | 680.3916 | 16 | `8ecabcd8168d5868` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 712.6552 | 16 | `e4e7eb16a418e291` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 0.8909 | false |
| notmat | 1 | 1% | 0.9782 | false |
| mat | 1 | 100% | 0.8950 | false |
| notmat | 1 | 100% | 1.0224 | false |
| mat | 4 | 1% | 0.9574 | false |
| notmat | 4 | 1% | 3.5751 | false |
| mat | 4 | 100% | 0.9884 | false |
| notmat | 4 | 100% | 3.1705 | false |
| mat | 16 | 1% | 1.0857 | false |
| notmat | 16 | 1% | 11.6956 | false |
| mat | 16 | 100% | 1.0704 | false |
| notmat | 16 | 100% | 11.9590 | false |

## Latency samples and derived throughput

_The reader is sequential, so `queries/s` is 1000/median_ms, a derived rate and not a concurrent-capacity test. min/max are the five sample extremes._

| axis | status | min ms | median ms | max ms | samples | derived queries/s |
|---|---|---|---|---|---|---|
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 123.3805 | 125.4940 | 126.3807 | 5 | 8.0 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 51.4937 | 59.7128 | 65.2649 | 5 | 16.7 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 54.8602 | 60.3632 | 77.4556 | 5 | 16.6 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 120.1145 | 122.8594 | 129.8396 | 5 | 8.1 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 91.2233 | 96.5852 | 103.9093 | 5 | 10.4 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 90.4900 | 100.2182 | 154.4650 | 5 | 10.0 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 124.0438 | 129.8621 | 137.4045 | 5 | 7.7 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 50.0951 | 55.3569 | 64.0799 | 5 | 18.1 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 59.9833 | 63.4369 | 66.0039 | 5 | 15.8 |
| family=bag;stmt=join_where;threshold=180000 | valid | 122.8755 | 128.5681 | 139.1805 | 5 | 7.8 |
| family=bag;stmt=join_where;threshold=20 | valid | 56.7763 | 59.5696 | 67.7487 | 5 | 16.8 |
| family=bag;stmt=join_where;threshold=2000 | valid | 53.6762 | 59.0996 | 69.1515 | 5 | 16.9 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 50.0227 | 58.1751 | 64.9819 | 5 | 17.2 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 51.5942 | 59.5917 | 59.9085 | 5 | 16.8 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 53.4428 | 64.2076 | 65.4908 | 5 | 15.6 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 54.5607 | 61.3180 | 65.7382 | 5 | 16.3 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 55.2881 | 56.1946 | 60.0686 | 5 | 17.8 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 56.5062 | 59.1869 | 63.4184 | 5 | 16.9 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 51.4335 | 63.1612 | 68.2015 | 5 | 15.8 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 59.4445 | 63.7872 | 70.8843 | 5 | 15.7 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 55.7070 | 57.2049 | 62.3355 | 5 | 17.5 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 53.0941 | 54.8806 | 60.0952 | 5 | 18.2 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 52.2341 | 53.8015 | 66.5736 | 5 | 18.6 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 53.5511 | 58.5001 | 68.0192 | 5 | 17.1 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 672.6644 | 680.3916 | 961.3915 | 5 | 1.5 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 684.5457 | 712.6552 | 755.2444 | 5 | 1.4 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 58.5826 | 62.8064 | 122.7327 | 5 | 15.9 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 58.1796 | 62.6936 | 64.9026 | 5 | 16.0 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 170.8845 | 200.8992 | 405.2510 | 5 | 5.0 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 171.5575 | 187.6542 | 204.1053 | 5 | 5.3 |

## Matched blended workload

- phase: `blended`; target offered rate: 10 updates/s
- offered 359, committed 359, failed 0 over 35.849 s
- offered rate 10.014/s, completed rate 10.014/s
- reader operations overlapping a committed write: 164 of 213
- reader retries of transient conflicts: 3
- per key: 359 keys, min 1, max 1 per key
- ledger: 359 keys reconciled, 0 mismatch
- update latency: avg 15.547 ms, min 1.584, max 69.141 over 359

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42851537`, 5 interleaved rounds
  1. family=reuse;policy=default;refs=1;consumer=1%
  2. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  3. family=reuse;policy=default;refs=4;consumer=100%
  4. family=bag;stmt=cte_joined_outer;threshold=20
  5. family=reuse;policy=notmat;refs=1;consumer=100%
  6. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  7. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  8. family=bag;stmt=cte_joined_outer_mat;threshold=20
  9. family=reuse;policy=notmat;refs=16;consumer=1%
  10. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  11. family=reuse;policy=default;refs=16;consumer=100%
  12. family=bag;stmt=join_where;threshold=20
  13. family=bag;stmt=cte_joined_outer;threshold=180000
  14. family=bag;stmt=join_where;threshold=2000
  15. family=reuse;policy=mat;refs=16;consumer=100%
  16. family=reuse;policy=notmat;refs=1;consumer=1%
  17. family=reuse;policy=notmat;refs=4;consumer=100%
  18. family=reuse;policy=mat;refs=16;consumer=1%
  19. family=reuse;policy=default;refs=16;consumer=1%
  20. family=reuse;policy=mat;refs=1;consumer=100%
  21. family=reuse;policy=notmat;refs=4;consumer=1%
  22. family=bag;stmt=join_where;threshold=180000
  23. family=reuse;policy=default;refs=4;consumer=1%
  24. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  25. family=reuse;policy=mat;refs=1;consumer=1%
  26. family=reuse;policy=notmat;refs=16;consumer=100%
  27. family=bag;stmt=cte_joined_outer;threshold=2000
  28. family=reuse;policy=mat;refs=4;consumer=1%
  29. family=reuse;policy=mat;refs=4;consumer=100%
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
| family=bag;stmt=cte_joined_outer;threshold=180000 | j05_cte_joined_outer | valid | `8ccb3d72ecd4e0ac` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20 | j05_cte_joined_outer | valid | `09316f370509139d` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=2000 | j05_cte_joined_outer | valid | `829808dc15af6a36` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | j05m_cte_joined_outer_mat | valid | `a074028f40cd3800` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | j05m_cte_joined_outer_mat | valid | `37071db9a1293962` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | j05m_cte_joined_outer_mat | valid | `40acf0b02b4f4858` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | j05n_cte_joined_outer_notmat | valid | `dcdcf0e8eb09c894` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | j05n_cte_joined_outer_notmat | valid | `8b71d37e2d00380d` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | j05n_cte_joined_outer_notmat | valid | `c0d1e43bb0d848be` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=join_where;threshold=180000 | j01_join_where | valid | `00a2be6b2e3c6271` |  |
| family=bag;stmt=join_where;threshold=20 | j01_join_where | valid | `9f9cd650148849f8` |  |
| family=bag;stmt=join_where;threshold=2000 | j01_join_where | valid | `6cc7ac877a74e2a0` |  |
| family=reuse;policy=default;refs=16;consumer=1% | r_default_n16_c1 | valid | `e33b692d73e204a9` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=16;consumer=100% | r_default_n16_c100 | valid | `e21e5033536eec65` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=1;consumer=1% | r_default_n1_c1 | valid | `48c3360dc26089a8` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=1;consumer=100% | r_default_n1_c100 | valid | `dbca7388712c192f` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=4;consumer=1% | r_default_n4_c1 | valid | `d9774833bab3d450` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=100% | r_default_n4_c100 | valid | `20f11895a76aedb7` | policy=default; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=1% | r_mat_n16_c1 | valid | `2297f4c9516b1656` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=100% | r_mat_n16_c100 | valid | `dd552acd262869aa` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=1% | r_mat_n1_c1 | valid | `14a9958e48ea12f8` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=100% | r_mat_n1_c100 | valid | `323d121527000186` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=1% | r_mat_n4_c1 | valid | `a0e8570238aec77b` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=100% | r_mat_n4_c100 | valid | `68009261e24c39bc` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=notmat;refs=16;consumer=1% | r_notmat_n16_c1 | valid | `8ecabcd8168d5868` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `e4e7eb16a418e291` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `d24dbbc49d287a7c` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `7c0dc957717661fa` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `48948ec51ed0942c` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `077e53594d575969` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, DIST, DEBUG, FORMAT TEXT` supported on yugabyte
- one host, single client: no independent hardware, real-network or equal-total-budget comparison
- no concurrent writers: the isolated read contract only; a matched blended read/write phase is a recorded gap
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
