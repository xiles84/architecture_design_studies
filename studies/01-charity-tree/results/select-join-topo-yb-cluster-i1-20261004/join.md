# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 30 cells, 30 valid, engine `yugabyte`, topology `cluster3`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-topo-yb-cluster-i1-20261004`
- run tag: `run/study-01/select-join-topo-yb-cluster-i1-20261004`
- environment: `host-zenbook-ux5406sa`
- repo commit: `85d78f9e942acf856b48d416d0fa26de43a4af66`
- engine version: `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit`
- database engine image: `docker.io/yugabytedb/yugabyte:2025.2.6.0-b111` @ `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`
- benchmark client image: `localhost/charitybench-select:select-join-topo-yb-cluster-i1-20261004` @ `d8db2fd98a7bc00fa2a5561ed01373926c9d950a8d93f0ae9ba4d53be68d6e38`
- inputs digest: `d8b0ad3fe902c1427281c60a89e03fc24f19f18e9b7364bd95455de2d60bb6dd`
- dataset fingerprint: `fcf26e08f377ec51`
- client CPU quota: `100000 100000`, memory: `2147483648`
- setting `client_workers`: 1 sequential reader; 0 writers
- setting `connection_pool`: 2 per endpoint = reader/setup + reserve
- setting `customers`: 2000
- setting `dataset`: study01-person-donation-charity
- setting `distributed_rpc`: not parsed; raw EXPLAIN retained (coverage gap)
- setting `endpoint_mode`: round-robin across 3 endpoints
- setting `endpoints`: yb-n1,yb-n2,yb-n3
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
- setting `topology`: cluster3
- setting `trials_per_cell`: 5
- setting `whole_run_count`: 1; no replicated winner in this file
- setting `writer_keys`: 1000
- setting `writer_target_rate_per_sec`: 10

### Endpoint distribution

A cluster arm must rotate its reader across nodes; a single seed endpoint would not be balanced cluster access.

| endpoint | operations | connections at end |
|---|---|---|
| yb-n1 | 163 | 1 |
| yb-n2 | 162 | 1 |
| yb-n3 | 162 | 1 |

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
| 20 (0.01%) | 58.1714 | missing | missing | missing | missing | missing | 61.0815 | 136.1561 | 58.0963 | missing |
| 2000 (1%) | 61.6067 | missing | missing | missing | missing | missing | 56.9925 | 92.5430 | 60.1670 | missing |
| 180000 (90%) | 156.3926 | missing | missing | missing | missing | missing | 153.6320 | 170.6905 | 123.1130 | missing |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | cte_joined_outer vs join_where | 58.1714 | 61.0815 | 1.0500 | false |
| 20 | cte_joined_outer_mat vs join_where | 58.1714 | 136.1561 | 2.3406 | false |
| 20 | cte_joined_outer_notmat vs join_where | 58.1714 | 58.0963 | 0.9987 | false |
| 2000 | cte_joined_outer vs join_where | 61.6067 | 56.9925 | 0.9251 | false |
| 2000 | cte_joined_outer_mat vs join_where | 61.6067 | 92.5430 | 1.5022 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 61.6067 | 60.1670 | 0.9766 | false |
| 180000 | cte_joined_outer vs join_where | 156.3926 | 153.6320 | 0.9823 | false |
| 180000 | cte_joined_outer_mat vs join_where | 156.3926 | 170.6905 | 1.0914 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 156.3926 | 123.1130 | 0.7872 | false |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_joined_outer_mat vs cte_joined_outer | 61.0815 | 136.1561 | 2.2291 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 61.0815 | 58.0963 | 0.9511 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 56.9925 | 92.5430 | 1.6238 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 56.9925 | 60.1670 | 1.0557 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 153.6320 | 170.6905 | 1.1110 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 153.6320 | 123.1130 | 0.8014 | false | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 57.6754 | 1 | `0fe61bd60760ddbe` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 60.2633 | 1 | `a9f11740cb3799ba` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | missing | | | | |
| default | 2 | 100% | missing | | | | |
| default | 4 | 1% | 59.7394 | 4 | `a9a68c9e1d4e02ec` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 56.4671 | 4 | `00e14f665e7e7690` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | missing | | | | |
| default | 8 | 100% | missing | | | | |
| default | 16 | 1% | 63.1134 | 16 | `9f85923e855531e2` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 59.3168 | 16 | `868900b79a97fb1f` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 56.7674 | 1 | `41dd493c78db0764` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 60.1879 | 1 | `157aac690e9feebc` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | missing | | | | |
| mat | 2 | 100% | missing | | | | |
| mat | 4 | 1% | 57.4937 | 4 | `c1142467bb3c21bd` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 59.9555 | 4 | `427347c3517cbba1` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | missing | | | | |
| mat | 8 | 100% | missing | | | | |
| mat | 16 | 1% | 58.5864 | 16 | `0081e443800a097b` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 58.1184 | 16 | `f591a7fc222cef27` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 59.1565 | 1 | `ee9763807c346827` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 59.3696 | 1 | `0fa0b20bb88348ad` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | missing | | | | |
| notmat | 2 | 100% | missing | | | | |
| notmat | 4 | 1% | 137.5559 | 4 | `fc0de162da31ae31` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 153.7558 | 4 | `5ec7a93137311c24` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | missing | | | | |
| notmat | 8 | 100% | missing | | | | |
| notmat | 16 | 1% | 533.9507 | 16 | `85856d61e18d286a` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 554.3707 | 16 | `e8cd1bd9572b64e5` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 0.9843 | false |
| notmat | 1 | 1% | 1.0257 | false |
| mat | 1 | 100% | 0.9987 | false |
| notmat | 1 | 100% | 0.9852 | false |
| mat | 4 | 1% | 0.9624 | false |
| notmat | 4 | 1% | 2.3026 | false |
| mat | 4 | 100% | 1.0618 | false |
| notmat | 4 | 100% | 2.7229 | false |
| mat | 16 | 1% | 0.9283 | false |
| notmat | 16 | 1% | 8.4602 | false |
| mat | 16 | 100% | 0.9798 | false |
| notmat | 16 | 100% | 9.3459 | false |

## Latency samples and derived throughput

_The reader is sequential, so `queries/s` is 1000/median_ms, a derived rate and not a concurrent-capacity test. min/max are the five sample extremes._

| axis | status | min ms | median ms | max ms | samples | derived queries/s |
|---|---|---|---|---|---|---|
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 146.6092 | 153.6320 | 253.9958 | 5 | 6.5 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 56.7309 | 61.0815 | 69.0046 | 5 | 16.4 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 54.4866 | 56.9925 | 61.0522 | 5 | 17.5 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 145.9890 | 170.6905 | 173.2675 | 5 | 5.9 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 125.0354 | 136.1561 | 158.8012 | 5 | 7.3 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 84.3367 | 92.5430 | 94.6426 | 5 | 10.8 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 108.4931 | 123.1130 | 133.4648 | 5 | 8.1 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 57.0138 | 58.0963 | 64.9113 | 5 | 17.2 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 55.6558 | 60.1670 | 66.8262 | 5 | 16.6 |
| family=bag;stmt=join_where;threshold=180000 | valid | 145.4577 | 156.3926 | 176.5973 | 5 | 6.4 |
| family=bag;stmt=join_where;threshold=20 | valid | 56.6184 | 58.1714 | 71.6083 | 5 | 17.2 |
| family=bag;stmt=join_where;threshold=2000 | valid | 57.0249 | 61.6067 | 73.0558 | 5 | 16.2 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 57.1520 | 63.1134 | 70.6554 | 5 | 15.8 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 55.6304 | 59.3168 | 70.7142 | 5 | 16.9 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 56.6082 | 57.6754 | 65.3493 | 5 | 17.3 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 58.5278 | 60.2633 | 62.7524 | 5 | 16.6 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 55.5205 | 59.7394 | 61.9664 | 5 | 16.7 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 55.1389 | 56.4671 | 60.6538 | 5 | 17.7 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 55.9234 | 58.5864 | 66.5020 | 5 | 17.1 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 56.7344 | 58.1184 | 69.9459 | 5 | 17.2 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 54.5988 | 56.7674 | 71.8139 | 5 | 17.6 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 56.6408 | 60.1879 | 76.3129 | 5 | 16.6 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 56.7864 | 57.4937 | 74.0293 | 5 | 17.4 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 54.3723 | 59.9555 | 60.5570 | 5 | 16.7 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 518.4398 | 533.9507 | 561.2395 | 5 | 1.9 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 536.3479 | 554.3707 | 621.6924 | 5 | 1.8 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 55.5362 | 59.1565 | 76.8855 | 5 | 16.9 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 58.4892 | 59.3696 | 77.8845 | 5 | 16.8 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 122.2511 | 137.5559 | 171.2778 | 5 | 7.3 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 135.4208 | 153.7558 | 164.7332 | 5 | 6.5 |

## Matched blended workload

- phase: `isolated`; target offered rate: 10 updates/s
- offered 0, committed 0, failed 0 over 0.000 s
- offered rate 0.000/s, completed rate 0.000/s
- per key: no keys touched
- ledger: 0 keys reconciled, 0 mismatch

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42531419`, 5 interleaved rounds
  1. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  2. family=reuse;policy=notmat;refs=4;consumer=1%
  3. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  4. family=reuse;policy=default;refs=4;consumer=1%
  5. family=reuse;policy=default;refs=1;consumer=100%
  6. family=reuse;policy=notmat;refs=16;consumer=1%
  7. family=reuse;policy=default;refs=16;consumer=1%
  8. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  9. family=reuse;policy=default;refs=1;consumer=1%
  10. family=reuse;policy=mat;refs=16;consumer=100%
  11. family=bag;stmt=cte_joined_outer;threshold=180000
  12. family=reuse;policy=notmat;refs=1;consumer=100%
  13. family=bag;stmt=join_where;threshold=20
  14. family=reuse;policy=mat;refs=1;consumer=1%
  15. family=bag;stmt=join_where;threshold=2000
  16. family=reuse;policy=notmat;refs=1;consumer=1%
  17. family=bag;stmt=cte_joined_outer;threshold=20
  18. family=reuse;policy=mat;refs=16;consumer=1%
  19. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  20. family=reuse;policy=notmat;refs=4;consumer=100%
  21. family=reuse;policy=notmat;refs=16;consumer=100%
  22. family=reuse;policy=mat;refs=1;consumer=100%
  23. family=reuse;policy=mat;refs=4;consumer=1%
  24. family=bag;stmt=cte_joined_outer_mat;threshold=20
  25. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  26. family=bag;stmt=join_where;threshold=180000
  27. family=reuse;policy=default;refs=16;consumer=100%
  28. family=bag;stmt=cte_joined_outer;threshold=2000
  29. family=reuse;policy=default;refs=4;consumer=100%
  30. family=reuse;policy=mat;refs=4;consumer=100%

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
| family=bag;stmt=cte_joined_outer;threshold=180000 | j05_cte_joined_outer | valid | `ef2513d343356cd8` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20 | j05_cte_joined_outer | valid | `de5009602a5f6bea` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=2000 | j05_cte_joined_outer | valid | `46eb5016ecc79969` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | j05m_cte_joined_outer_mat | valid | `df33d024652ca0ef` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | j05m_cte_joined_outer_mat | valid | `b91afbb3d877cb0b` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | j05m_cte_joined_outer_mat | valid | `ec98fc6a9554bcd8` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | j05n_cte_joined_outer_notmat | valid | `17b8c753250b7e37` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | j05n_cte_joined_outer_notmat | valid | `0d7258898f35446c` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | j05n_cte_joined_outer_notmat | valid | `8b5a357ed11054d5` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=join_where;threshold=180000 | j01_join_where | valid | `34b4545097e37cde` |  |
| family=bag;stmt=join_where;threshold=20 | j01_join_where | valid | `7ea5e44ffcbeaf35` |  |
| family=bag;stmt=join_where;threshold=2000 | j01_join_where | valid | `45bb53fbccf80efd` |  |
| family=reuse;policy=default;refs=16;consumer=1% | r_default_n16_c1 | valid | `9f85923e855531e2` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=16;consumer=100% | r_default_n16_c100 | valid | `868900b79a97fb1f` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=1;consumer=1% | r_default_n1_c1 | valid | `0fe61bd60760ddbe` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=1;consumer=100% | r_default_n1_c100 | valid | `a9f11740cb3799ba` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=4;consumer=1% | r_default_n4_c1 | valid | `a9a68c9e1d4e02ec` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=100% | r_default_n4_c100 | valid | `00e14f665e7e7690` | policy=default; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=1% | r_mat_n16_c1 | valid | `0081e443800a097b` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=100% | r_mat_n16_c100 | valid | `f591a7fc222cef27` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=1% | r_mat_n1_c1 | valid | `41dd493c78db0764` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=100% | r_mat_n1_c100 | valid | `157aac690e9feebc` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=1% | r_mat_n4_c1 | valid | `c1142467bb3c21bd` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=100% | r_mat_n4_c100 | valid | `427347c3517cbba1` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=notmat;refs=16;consumer=1% | r_notmat_n16_c1 | valid | `85856d61e18d286a` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `e8cd1bd9572b64e5` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `ee9763807c346827` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `0fa0b20bb88348ad` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `fc0de162da31ae31` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `5ec7a93137311c24` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, DIST, DEBUG, FORMAT TEXT` supported on yugabyte
- one host, single client: no independent hardware, real-network or equal-total-budget comparison
- no concurrent writers: the isolated read contract only; a matched blended read/write phase is a recorded gap
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
