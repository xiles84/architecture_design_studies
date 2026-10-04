# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 30 cells, 29 valid, engine `yugabyte`, topology `single`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-topo-yb-single-bl1-20261004`
- run tag: `run/study-01/select-join-topo-yb-single-bl1-20261004`
- environment: `host-zenbook-ux5406sa`
- repo commit: `622c9f5fdbd2f051e1af0a4850768a7fd073ee94`
- engine version: `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit`
- database engine image: `docker.io/yugabytedb/yugabyte:2025.2.6.0-b111` @ `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`
- benchmark client image: `localhost/charitybench-select:select-join-topo-yb-single-bl1-20261004` @ `2461bd700861e8ad769d5e61d59dfc5fbdd6f023215c0e688ba52a5a479e8c4a`
- inputs digest: `faadf2240d96a80ef37b5bd1f46c4d18c602913e1ac670c14a54edf33c1b3cc1`
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
| yb-single | 482 | 2 |

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
| 20 (0.01%) | 55.6549 | missing | missing | missing | missing | missing | 62.1824 | 91.6442 | 64.9638 | missing |
| 2000 (1%) | 56.0031 | missing | missing | missing | missing | missing | 72.1799 | 124.1417 | 87.2899 | missing |
| 180000 (90%) | 127.6305 | missing | missing | missing | missing | missing | 134.1102 | 124.1399 | 111.6823 | missing |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | cte_joined_outer vs join_where | 55.6549 | 62.1824 | 1.1173 | false |
| 20 | cte_joined_outer_mat vs join_where | 55.6549 | 91.6442 | 1.6466 | false |
| 20 | cte_joined_outer_notmat vs join_where | 55.6549 | 64.9638 | 1.1673 | false |
| 2000 | cte_joined_outer vs join_where | 56.0031 | 72.1799 | 1.2889 | false |
| 2000 | cte_joined_outer_mat vs join_where | 56.0031 | 124.1417 | 2.2167 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 56.0031 | 87.2899 | 1.5587 | false |
| 180000 | cte_joined_outer vs join_where | 127.6305 | 134.1102 | 1.0508 | false |
| 180000 | cte_joined_outer_mat vs join_where | 127.6305 | 124.1399 | 0.9727 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 127.6305 | 111.6823 | 0.8750 | false |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_joined_outer_mat vs cte_joined_outer | 62.1824 | 91.6442 | 1.4738 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 62.1824 | 64.9638 | 1.0447 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 72.1799 | 124.1417 | 1.7199 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 72.1799 | 87.2899 | 1.2093 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 134.1102 | 124.1399 | 0.9257 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 134.1102 | 111.6823 | 0.8328 | false | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 58.1108 | 1 | `69066dd344f13778` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 65.8372 | 1 | `7c8a2236a17f52c8` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | missing | | | | |
| default | 2 | 100% | missing | | | | |
| default | 4 | 1% | 66.3608 | 4 | `03e879ef863698b3` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 60.9281 | 4 | `1f3b4a14b58ec7ed` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | missing | | | | |
| default | 8 | 100% | missing | | | | |
| default | 16 | 1% | 60.9051 | 16 | `6cd8589cbfdd1537` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 67.6996 | 16 | `f5af1a2c6963fe10` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 56.6154 | 1 | `c443a32dd4947155` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 61.0683 | 1 | `5e9e6b811a86dc72` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | missing | | | | |
| mat | 2 | 100% | missing | | | | |
| mat | 4 | 1% | 68.3402 | 4 | `cc71950c00766457` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 66.5361 | 4 | `816c9e32b021df68` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | missing | | | | |
| mat | 8 | 100% | missing | | | | |
| mat | 16 | 1% | 57.5361 | 16 | `acbb24a61f8ceab2` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 62.6728 | 16 | `fc2f11924a4f925e` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 64.2556 | 1 | `43008251f698c196` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 61.7777 | 1 | `5503a316d6293de9` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | missing | | | | |
| notmat | 2 | 100% | missing | | | | |
| notmat | 4 | 1% | error | 0 | `e3b0c44298fc1c14` |  | error |
| notmat | 4 | 100% | 195.3919 | 4 | `ae1e339b042b9eda` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | missing | | | | |
| notmat | 8 | 100% | missing | | | | |
| notmat | 16 | 1% | 798.5774 | 16 | `930ff6297ae5e83a` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 705.3341 | 16 | `6aff89d1164dff81` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 0.9743 | false |
| notmat | 1 | 1% | 1.1057 | false |
| mat | 1 | 100% | 0.9276 | false |
| notmat | 1 | 100% | 0.9383 | false |
| mat | 4 | 1% | 1.0298 | false |
| notmat | 4 | 1% | N/A | N/A |
| mat | 4 | 100% | 1.0920 | false |
| notmat | 4 | 100% | 3.2069 | false |
| mat | 16 | 1% | 0.9447 | false |
| notmat | 16 | 1% | 13.1118 | false |
| mat | 16 | 100% | 0.9257 | false |
| notmat | 16 | 100% | 10.4186 | false |

## Latency samples and derived throughput

_The reader is sequential, so `queries/s` is 1000/median_ms, a derived rate and not a concurrent-capacity test. min/max are the five sample extremes._

| axis | status | min ms | median ms | max ms | samples | derived queries/s |
|---|---|---|---|---|---|---|
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 124.7107 | 134.1102 | 500.8672 | 5 | 7.5 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 56.9994 | 62.1824 | 85.5740 | 5 | 16.1 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 54.6017 | 72.1799 | 100.3301 | 5 | 13.9 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 107.1795 | 124.1399 | 160.4428 | 5 | 8.1 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 86.8864 | 91.6442 | 141.4124 | 5 | 10.9 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 83.6876 | 124.1417 | 176.7833 | 5 | 8.1 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 108.4152 | 111.6823 | 180.9478 | 5 | 9.0 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 53.5926 | 64.9638 | 95.8448 | 5 | 15.4 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 55.4883 | 87.2899 | 93.2326 | 5 | 11.5 |
| family=bag;stmt=join_where;threshold=180000 | valid | 111.0937 | 127.6305 | 300.1966 | 5 | 7.8 |
| family=bag;stmt=join_where;threshold=20 | valid | 51.8790 | 55.6549 | 74.0018 | 5 | 18.0 |
| family=bag;stmt=join_where;threshold=2000 | valid | 54.9311 | 56.0031 | 78.7999 | 5 | 17.9 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 54.3404 | 60.9051 | 72.1130 | 5 | 16.4 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 56.9603 | 67.6996 | 106.3470 | 5 | 14.8 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 55.7714 | 58.1108 | 72.1557 | 5 | 17.2 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 57.3157 | 65.8372 | 109.1587 | 5 | 15.2 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 59.8119 | 66.3608 | 78.4456 | 5 | 15.1 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 58.7926 | 60.9281 | 74.0923 | 5 | 16.4 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 55.5766 | 57.5361 | 79.1655 | 5 | 17.4 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 55.0235 | 62.6728 | 106.3157 | 5 | 16.0 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 54.1520 | 56.6154 | 87.6816 | 5 | 17.7 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 55.3391 | 61.0683 | 100.4205 | 5 | 16.4 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 54.7501 | 68.3402 | 107.3391 | 5 | 14.6 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 58.0644 | 66.5361 | 109.9843 | 5 | 15.0 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 590.0703 | 798.5774 | 1128.7239 | 5 | 1.3 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 599.7508 | 705.3341 | 1238.5489 | 5 | 1.4 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 54.5356 | 64.2556 | 106.5722 | 5 | 15.6 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 55.2539 | 61.7777 | 116.1623 | 5 | 16.2 |
| family=reuse;policy=notmat;refs=4;consumer=1% | error | N/A | N/A | N/A | 0 | N/A |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 163.7464 | 195.3919 | 322.2781 | 5 | 5.1 |

## Matched blended workload

- phase: `blended`; target offered rate: 10 updates/s
- offered 362, committed 362, failed 0 over 36.185 s
- offered rate 10.004/s, completed rate 10.004/s
- reader operations overlapping a committed write: 155 of 205
- reader retries of transient conflicts: 0
- per key: 362 keys, min 1, max 1 per key
- ledger: 362 keys reconciled, 0 mismatch
- update latency: avg 3.723 ms, min 1.528, max 67.948 over 362

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42848202`, 5 interleaved rounds
  1. family=reuse;policy=default;refs=4;consumer=100%
  2. family=bag;stmt=cte_joined_outer_mat;threshold=20
  3. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  4. family=bag;stmt=cte_joined_outer;threshold=20
  5. family=reuse;policy=mat;refs=16;consumer=1%
  6. family=bag;stmt=join_where;threshold=2000
  7. family=reuse;policy=mat;refs=1;consumer=1%
  8. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  9. family=reuse;policy=default;refs=16;consumer=1%
  10. family=reuse;policy=default;refs=1;consumer=1%
  11. family=bag;stmt=join_where;threshold=20
  12. family=reuse;policy=notmat;refs=16;consumer=1%
  13. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  14. family=reuse;policy=mat;refs=1;consumer=100%
  15. family=reuse;policy=notmat;refs=1;consumer=1%
  16. family=bag;stmt=join_where;threshold=180000
  17. family=reuse;policy=mat;refs=4;consumer=100%
  18. family=reuse;policy=notmat;refs=4;consumer=100%
  19. family=bag;stmt=cte_joined_outer;threshold=180000
  20. family=reuse;policy=default;refs=16;consumer=100%
  21. family=reuse;policy=default;refs=1;consumer=100%
  22. family=reuse;policy=mat;refs=4;consumer=1%
  23. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  24. family=reuse;policy=notmat;refs=1;consumer=100%
  25. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  26. family=reuse;policy=mat;refs=16;consumer=100%
  27. family=bag;stmt=cte_joined_outer;threshold=2000
  28. family=reuse;policy=notmat;refs=16;consumer=100%
  29. family=reuse;policy=notmat;refs=4;consumer=1%
  30. family=reuse;policy=default;refs=4;consumer=1%

## Live negative controls

| control | pretended contract | detected | rows | note |
|---|---|---|---|---|
| x01_optional_where | sel_join_optional | true | 1 |  |
| x02_join_limit_before | sel_join_bag | true | 7 |  |

## Excluded cells

- `family=reuse;policy=notmat;refs=4;consumer=1%` (r_notmat_n4_c1) — **error** — ERROR: Restart read required (SQLSTATE 40001)

A cell that is invalid, unsupported or timed out is never ranked.

## Plan signatures

| axis | statement | status | structure signature | policy note |
|---|---|---|---|---|
| family=bag;stmt=cte_joined_outer;threshold=180000 | j05_cte_joined_outer | valid | `294f9b71ccbf565c` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20 | j05_cte_joined_outer | valid | `4da1062c27bc2230` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=2000 | j05_cte_joined_outer | valid | `534b6497b0c464d1` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | j05m_cte_joined_outer_mat | valid | `64c84490df1b23a7` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | j05m_cte_joined_outer_mat | valid | `3144c8d593ca894e` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | j05m_cte_joined_outer_mat | valid | `570d1c569f41b7e1` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | j05n_cte_joined_outer_notmat | valid | `f49195e860c66663` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | j05n_cte_joined_outer_notmat | valid | `99656770c702461d` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | j05n_cte_joined_outer_notmat | valid | `85c0254d43960f08` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=join_where;threshold=180000 | j01_join_where | valid | `37b0f33e1532f8ab` |  |
| family=bag;stmt=join_where;threshold=20 | j01_join_where | valid | `a532b6bc99f08fd1` |  |
| family=bag;stmt=join_where;threshold=2000 | j01_join_where | valid | `b9533cf3e3547979` |  |
| family=reuse;policy=default;refs=16;consumer=1% | r_default_n16_c1 | valid | `6cd8589cbfdd1537` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=16;consumer=100% | r_default_n16_c100 | valid | `f5af1a2c6963fe10` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=1;consumer=1% | r_default_n1_c1 | valid | `69066dd344f13778` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=1;consumer=100% | r_default_n1_c100 | valid | `7c8a2236a17f52c8` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=4;consumer=1% | r_default_n4_c1 | valid | `03e879ef863698b3` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=100% | r_default_n4_c100 | valid | `1f3b4a14b58ec7ed` | policy=default; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=1% | r_mat_n16_c1 | valid | `acbb24a61f8ceab2` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=100% | r_mat_n16_c100 | valid | `fc2f11924a4f925e` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=1% | r_mat_n1_c1 | valid | `c443a32dd4947155` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=100% | r_mat_n1_c100 | valid | `5e9e6b811a86dc72` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=1% | r_mat_n4_c1 | valid | `cc71950c00766457` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=100% | r_mat_n4_c100 | valid | `816c9e32b021df68` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=notmat;refs=16;consumer=1% | r_notmat_n16_c1 | valid | `930ff6297ae5e83a` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `6aff89d1164dff81` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `43008251f698c196` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `5503a316d6293de9` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | error | `e3b0c44298fc1c14` |  |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `ae1e339b042b9eda` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, DIST, DEBUG, FORMAT TEXT` supported on yugabyte
- one host, single client: no independent hardware, real-network or equal-total-budget comparison
- no concurrent writers: the isolated read contract only; a matched blended read/write phase is a recorded gap
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
