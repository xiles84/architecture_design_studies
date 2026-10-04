# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 30 cells, 30 valid, engine `yugabyte`, topology `single`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-topo-yb-equal-bo3-20261004`
- run tag: `run/study-01/select-join-topo-yb-equal-bo3-20261004`
- environment: `host-zenbook-ux5406sa`
- repo commit: `455b777fc02e015efdd9e54d25d3aa3343eb206f`
- engine version: `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit`
- database engine image: `docker.io/yugabytedb/yugabyte:2025.2.6.0-b111` @ `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`
- benchmark client image: `localhost/charitybench-select:select-join-topo-yb-equal-bo3-20261004` @ `79bc753ff7e57e102982a323939c43e68e9b3664887683d3a5d5bf072650f4c1`
- inputs digest: `66adf072b41a282ea590391e71adb8fc8a6fb5734b9e83c818274fbc83282583`
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
| yb-single | 488 | 2 |

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
| 20 (0.01%) | 60.4591 | missing | missing | missing | missing | missing | 61.7175 | 97.6875 | 61.2431 | missing |
| 2000 (1%) | 63.9741 | missing | missing | missing | missing | missing | 63.4538 | 95.3780 | 65.8494 | missing |
| 180000 (90%) | 130.7395 | missing | missing | missing | missing | missing | 144.3187 | 121.1784 | 125.5814 | missing |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | cte_joined_outer vs join_where | 60.4591 | 61.7175 | 1.0208 | false |
| 20 | cte_joined_outer_mat vs join_where | 60.4591 | 97.6875 | 1.6158 | false |
| 20 | cte_joined_outer_notmat vs join_where | 60.4591 | 61.2431 | 1.0130 | false |
| 2000 | cte_joined_outer vs join_where | 63.9741 | 63.4538 | 0.9919 | false |
| 2000 | cte_joined_outer_mat vs join_where | 63.9741 | 95.3780 | 1.4909 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 63.9741 | 65.8494 | 1.0293 | false |
| 180000 | cte_joined_outer vs join_where | 130.7395 | 144.3187 | 1.1039 | false |
| 180000 | cte_joined_outer_mat vs join_where | 130.7395 | 121.1784 | 0.9269 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 130.7395 | 125.5814 | 0.9605 | false |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_joined_outer_mat vs cte_joined_outer | 61.7175 | 97.6875 | 1.5828 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 61.7175 | 61.2431 | 0.9923 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 63.4538 | 95.3780 | 1.5031 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 63.4538 | 65.8494 | 1.0378 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 144.3187 | 121.1784 | 0.8397 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 144.3187 | 125.5814 | 0.8702 | false | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 65.0221 | 1 | `22acc8d2e9fad541` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 72.9628 | 1 | `2d2cf4e37e5b66fd` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | missing | | | | |
| default | 2 | 100% | missing | | | | |
| default | 4 | 1% | 63.2297 | 4 | `2954d61583de3d30` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 67.6621 | 4 | `8e57ba4aa4086247` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | missing | | | | |
| default | 8 | 100% | missing | | | | |
| default | 16 | 1% | 64.7180 | 16 | `c091ecefc63e2b86` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 64.6120 | 16 | `43b51b44a30c0a9f` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 71.6489 | 1 | `21387fb75a526c7a` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 61.1923 | 1 | `a74dfddd2f611fa0` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | missing | | | | |
| mat | 2 | 100% | missing | | | | |
| mat | 4 | 1% | 72.4158 | 4 | `e382053167e2e320` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 61.8597 | 4 | `7c37943718c26f19` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | missing | | | | |
| mat | 8 | 100% | missing | | | | |
| mat | 16 | 1% | 61.2280 | 16 | `37c1a9232ea8e0d6` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 64.2402 | 16 | `043b44d95f764e42` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 65.2811 | 1 | `822370d1fb919eaf` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 67.5904 | 1 | `4020363533148bbb` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | missing | | | | |
| notmat | 2 | 100% | missing | | | | |
| notmat | 4 | 1% | 144.6326 | 4 | `b445c8fe927aac21` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 139.1589 | 4 | `931d1c8071b4c11a` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | missing | | | | |
| notmat | 8 | 100% | missing | | | | |
| notmat | 16 | 1% | 514.7835 | 16 | `eaea7d8565c084f3` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 462.0861 | 16 | `5cbbed4287a50268` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 1.1019 | false |
| notmat | 1 | 1% | 1.0040 | false |
| mat | 1 | 100% | 0.8387 | false |
| notmat | 1 | 100% | 0.9264 | false |
| mat | 4 | 1% | 1.1453 | false |
| notmat | 4 | 1% | 2.2874 | false |
| mat | 4 | 100% | 0.9142 | false |
| notmat | 4 | 100% | 2.0567 | false |
| mat | 16 | 1% | 0.9461 | false |
| notmat | 16 | 1% | 7.9543 | false |
| mat | 16 | 100% | 0.9942 | false |
| notmat | 16 | 100% | 7.1517 | false |

## Latency samples and derived throughput

_The reader is sequential, so `queries/s` is 1000/median_ms, a derived rate and not a concurrent-capacity test. min/max are the five sample extremes._

| axis | status | min ms | median ms | max ms | samples | derived queries/s |
|---|---|---|---|---|---|---|
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 142.1707 | 144.3187 | 149.0501 | 5 | 6.9 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 57.5148 | 61.7175 | 64.8667 | 5 | 16.2 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 60.5117 | 63.4538 | 64.2236 | 5 | 15.8 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 108.6411 | 121.1784 | 199.1424 | 5 | 8.3 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 92.2411 | 97.6875 | 99.7981 | 5 | 10.2 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 91.9011 | 95.3780 | 113.7879 | 5 | 10.5 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 118.5552 | 125.5814 | 131.3664 | 5 | 8.0 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 59.0966 | 61.2431 | 66.0688 | 5 | 16.3 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 63.3185 | 65.8494 | 71.4524 | 5 | 15.2 |
| family=bag;stmt=join_where;threshold=180000 | valid | 126.6084 | 130.7395 | 139.9304 | 5 | 7.6 |
| family=bag;stmt=join_where;threshold=20 | valid | 59.4448 | 60.4591 | 64.4733 | 5 | 16.5 |
| family=bag;stmt=join_where;threshold=2000 | valid | 61.9123 | 63.9741 | 94.6410 | 5 | 15.6 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 60.8171 | 64.7180 | 68.3556 | 5 | 15.5 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 60.9822 | 64.6120 | 65.8375 | 5 | 15.5 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 62.6392 | 65.0221 | 67.2025 | 5 | 15.4 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 70.2093 | 72.9628 | 77.3859 | 5 | 13.7 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 60.3863 | 63.2297 | 64.7319 | 5 | 15.8 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 65.8982 | 67.6621 | 69.2062 | 5 | 14.8 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 58.7808 | 61.2280 | 63.9431 | 5 | 16.3 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 61.6865 | 64.2402 | 69.6156 | 5 | 15.6 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 69.3414 | 71.6489 | 75.8221 | 5 | 14.0 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 57.4092 | 61.1923 | 62.8778 | 5 | 16.3 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 68.8798 | 72.4158 | 75.4653 | 5 | 13.8 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 59.7094 | 61.8597 | 68.3743 | 5 | 16.2 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 419.9891 | 514.7835 | 572.2596 | 5 | 1.9 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 457.0255 | 462.0861 | 476.3233 | 5 | 2.2 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 59.5976 | 65.2811 | 67.5905 | 5 | 15.3 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 65.0966 | 67.5904 | 70.4388 | 5 | 14.8 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 139.2237 | 144.6326 | 272.6527 | 5 | 6.9 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 132.0363 | 139.1589 | 149.6716 | 5 | 7.2 |

## Matched blended workload

- phase: `blended`; target offered rate: 10 updates/s
- offered 326, committed 326, failed 0 over 32.523 s
- offered rate 10.024/s, completed rate 10.024/s
- reader operations overlapping a committed write: 158 of 211
- reader retries of transient conflicts: 1
- per key: 326 keys, min 1, max 1 per key
- ledger: 326 keys reconciled, 0 mismatch
- update latency: avg 3.740 ms, min 1.822, max 74.126 over 326

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42488893`, 5 interleaved rounds
  1. family=reuse;policy=mat;refs=4;consumer=1%
  2. family=reuse;policy=default;refs=1;consumer=100%
  3. family=reuse;policy=notmat;refs=16;consumer=100%
  4. family=reuse;policy=mat;refs=1;consumer=1%
  5. family=bag;stmt=cte_joined_outer;threshold=180000
  6. family=reuse;policy=default;refs=4;consumer=100%
  7. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  8. family=bag;stmt=join_where;threshold=180000
  9. family=reuse;policy=notmat;refs=1;consumer=100%
  10. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  11. family=reuse;policy=mat;refs=4;consumer=100%
  12. family=bag;stmt=cte_joined_outer_mat;threshold=20
  13. family=reuse;policy=default;refs=16;consumer=100%
  14. family=bag;stmt=join_where;threshold=20
  15. family=reuse;policy=mat;refs=16;consumer=100%
  16. family=reuse;policy=default;refs=16;consumer=1%
  17. family=reuse;policy=mat;refs=16;consumer=1%
  18. family=bag;stmt=cte_joined_outer;threshold=2000
  19. family=bag;stmt=join_where;threshold=2000
  20. family=reuse;policy=mat;refs=1;consumer=100%
  21. family=reuse;policy=notmat;refs=4;consumer=1%
  22. family=reuse;policy=default;refs=1;consumer=1%
  23. family=reuse;policy=default;refs=4;consumer=1%
  24. family=reuse;policy=notmat;refs=1;consumer=1%
  25. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  26. family=bag;stmt=cte_joined_outer;threshold=20
  27. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  28. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  29. family=reuse;policy=notmat;refs=4;consumer=100%
  30. family=reuse;policy=notmat;refs=16;consumer=1%

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
| family=bag;stmt=cte_joined_outer;threshold=180000 | j05_cte_joined_outer | valid | `8fc9c0d331061895` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20 | j05_cte_joined_outer | valid | `963655742e81033d` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=2000 | j05_cte_joined_outer | valid | `3416b01d647f8833` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | j05m_cte_joined_outer_mat | valid | `19617b3a83701fd2` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | j05m_cte_joined_outer_mat | valid | `dcb1376f31c60d9a` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | j05m_cte_joined_outer_mat | valid | `c3d8131eb0f295b6` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | j05n_cte_joined_outer_notmat | valid | `1a796322150d21f2` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | j05n_cte_joined_outer_notmat | valid | `4dde4083f09b1713` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | j05n_cte_joined_outer_notmat | valid | `c8f9e1b455e713d8` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=join_where;threshold=180000 | j01_join_where | valid | `7e67e667c1048767` |  |
| family=bag;stmt=join_where;threshold=20 | j01_join_where | valid | `f19cc4d8df1eee9b` |  |
| family=bag;stmt=join_where;threshold=2000 | j01_join_where | valid | `4d85596740e11648` |  |
| family=reuse;policy=default;refs=16;consumer=1% | r_default_n16_c1 | valid | `c091ecefc63e2b86` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=16;consumer=100% | r_default_n16_c100 | valid | `43b51b44a30c0a9f` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=1;consumer=1% | r_default_n1_c1 | valid | `22acc8d2e9fad541` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=1;consumer=100% | r_default_n1_c100 | valid | `2d2cf4e37e5b66fd` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=4;consumer=1% | r_default_n4_c1 | valid | `2954d61583de3d30` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=100% | r_default_n4_c100 | valid | `8e57ba4aa4086247` | policy=default; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=1% | r_mat_n16_c1 | valid | `37c1a9232ea8e0d6` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=100% | r_mat_n16_c100 | valid | `043b44d95f764e42` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=1% | r_mat_n1_c1 | valid | `21387fb75a526c7a` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=100% | r_mat_n1_c100 | valid | `a74dfddd2f611fa0` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=1% | r_mat_n4_c1 | valid | `e382053167e2e320` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=100% | r_mat_n4_c100 | valid | `7c37943718c26f19` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=notmat;refs=16;consumer=1% | r_notmat_n16_c1 | valid | `eaea7d8565c084f3` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `5cbbed4287a50268` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `822370d1fb919eaf` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `4020363533148bbb` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `b445c8fe927aac21` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `931d1c8071b4c11a` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, DIST, DEBUG, FORMAT TEXT` supported on yugabyte
- one host, single client: no independent hardware, real-network or equal-total-budget comparison
- no concurrent writers: the isolated read contract only; a matched blended read/write phase is a recorded gap
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
