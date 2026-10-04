# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 30 cells, 30 valid, engine `yugabyte`, topology `single`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-topo-yb-equal-i3-20261004`
- run tag: `run/study-01/select-join-topo-yb-equal-i3-20261004`
- environment: `host-zenbook-ux5406sa`
- repo commit: `335542f10685ca626b2f4cb2b086f0ad67e620e0`
- engine version: `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit`
- database engine image: `docker.io/yugabytedb/yugabyte:2025.2.6.0-b111` @ `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`
- benchmark client image: `localhost/charitybench-select:select-join-topo-yb-equal-i3-20261004` @ `d8db2fd98a7bc00fa2a5561ed01373926c9d950a8d93f0ae9ba4d53be68d6e38`
- inputs digest: `82384d8f501320d635a9c7efd20dba3bb9a5c408cc343d31fbfcbcea7449d169`
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
| 20 (0.01%) | 57.6055 | missing | missing | missing | missing | missing | 58.6695 | 93.3009 | 60.1860 | missing |
| 2000 (1%) | 58.8241 | missing | missing | missing | missing | missing | 59.6361 | 101.9812 | 61.0918 | missing |
| 180000 (90%) | 121.7237 | missing | missing | missing | missing | missing | 131.2919 | 117.5039 | 128.3253 | missing |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | cte_joined_outer vs join_where | 57.6055 | 58.6695 | 1.0185 | false |
| 20 | cte_joined_outer_mat vs join_where | 57.6055 | 93.3009 | 1.6197 | false |
| 20 | cte_joined_outer_notmat vs join_where | 57.6055 | 60.1860 | 1.0448 | false |
| 2000 | cte_joined_outer vs join_where | 58.8241 | 59.6361 | 1.0138 | false |
| 2000 | cte_joined_outer_mat vs join_where | 58.8241 | 101.9812 | 1.7337 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 58.8241 | 61.0918 | 1.0386 | false |
| 180000 | cte_joined_outer vs join_where | 121.7237 | 131.2919 | 1.0786 | false |
| 180000 | cte_joined_outer_mat vs join_where | 121.7237 | 117.5039 | 0.9653 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 121.7237 | 128.3253 | 1.0542 | false |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_joined_outer_mat vs cte_joined_outer | 58.6695 | 93.3009 | 1.5903 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 58.6695 | 60.1860 | 1.0258 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 59.6361 | 101.9812 | 1.7101 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 59.6361 | 61.0918 | 1.0244 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 131.2919 | 117.5039 | 0.8950 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 131.2919 | 128.3253 | 0.9774 | false | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 60.5576 | 1 | `51d0cd00b995aaa4` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 62.9782 | 1 | `870cb6cde42ccd8c` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | missing | | | | |
| default | 2 | 100% | missing | | | | |
| default | 4 | 1% | 62.5908 | 4 | `61a80c16c3827eeb` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 57.7512 | 4 | `03aebc676846112f` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | missing | | | | |
| default | 8 | 100% | missing | | | | |
| default | 16 | 1% | 62.5562 | 16 | `fbfa91d9d372f834` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 59.7381 | 16 | `53786e8fbd59d292` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 57.0674 | 1 | `165ec7fe91e76792` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 57.7332 | 1 | `1c989baeeb638c5a` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | missing | | | | |
| mat | 2 | 100% | missing | | | | |
| mat | 4 | 1% | 57.6409 | 4 | `f93247b1cf3d3edb` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 57.3764 | 4 | `8dad9b965389390e` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | missing | | | | |
| mat | 8 | 100% | missing | | | | |
| mat | 16 | 1% | 60.7020 | 16 | `e45a3ed613cb52eb` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 59.7208 | 16 | `dffebeb20ecbb208` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 61.5403 | 1 | `44e8e67d929c0172` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 59.5980 | 1 | `7e99c0ea2bf5f956` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | missing | | | | |
| notmat | 2 | 100% | missing | | | | |
| notmat | 4 | 1% | 136.3275 | 4 | `eda3480b42d37e2f` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 140.1911 | 4 | `1dadc273273c50b5` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | missing | | | | |
| notmat | 8 | 100% | missing | | | | |
| notmat | 16 | 1% | 393.5135 | 16 | `6d2f0cf4f2a6aa8e` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 411.2624 | 16 | `1f5deaaa277c1da2` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 0.9424 | false |
| notmat | 1 | 1% | 1.0162 | false |
| mat | 1 | 100% | 0.9167 | false |
| notmat | 1 | 100% | 0.9463 | false |
| mat | 4 | 1% | 0.9209 | false |
| notmat | 4 | 1% | 2.1781 | false |
| mat | 4 | 100% | 0.9935 | false |
| notmat | 4 | 100% | 2.4275 | false |
| mat | 16 | 1% | 0.9704 | false |
| notmat | 16 | 1% | 6.2906 | false |
| mat | 16 | 100% | 0.9997 | false |
| notmat | 16 | 100% | 6.8844 | false |

## Latency samples and derived throughput

_The reader is sequential, so `queries/s` is 1000/median_ms, a derived rate and not a concurrent-capacity test. min/max are the five sample extremes._

| axis | status | min ms | median ms | max ms | samples | derived queries/s |
|---|---|---|---|---|---|---|
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 120.0888 | 131.2919 | 141.7544 | 5 | 7.6 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 58.1775 | 58.6695 | 66.0297 | 5 | 17.0 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 56.5239 | 59.6361 | 63.1446 | 5 | 16.8 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 114.2465 | 117.5039 | 169.0208 | 5 | 8.5 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 87.5303 | 93.3009 | 101.8051 | 5 | 10.7 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 82.0932 | 101.9812 | 102.1259 | 5 | 9.8 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 116.0758 | 128.3253 | 134.4010 | 5 | 7.8 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 54.8102 | 60.1860 | 64.1806 | 5 | 16.6 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 59.8975 | 61.0918 | 65.6388 | 5 | 16.4 |
| family=bag;stmt=join_where;threshold=180000 | valid | 115.6425 | 121.7237 | 139.7411 | 5 | 8.2 |
| family=bag;stmt=join_where;threshold=20 | valid | 54.9858 | 57.6055 | 65.6624 | 5 | 17.4 |
| family=bag;stmt=join_where;threshold=2000 | valid | 58.1827 | 58.8241 | 64.4363 | 5 | 17.0 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 59.9094 | 62.5562 | 65.5443 | 5 | 16.0 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 59.3111 | 59.7381 | 62.0335 | 5 | 16.7 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 56.6547 | 60.5576 | 71.9759 | 5 | 16.5 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 56.2842 | 62.9782 | 69.3493 | 5 | 15.9 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 58.0861 | 62.5908 | 64.2965 | 5 | 16.0 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 54.3502 | 57.7512 | 62.0592 | 5 | 17.3 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 56.5879 | 60.7020 | 70.5635 | 5 | 16.5 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 57.6421 | 59.7208 | 65.6823 | 5 | 16.7 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 54.8855 | 57.0674 | 61.8291 | 5 | 17.5 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 54.4473 | 57.7332 | 59.4917 | 5 | 17.3 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 56.0256 | 57.6409 | 64.6120 | 5 | 17.3 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 54.2824 | 57.3764 | 61.5522 | 5 | 17.4 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 376.5504 | 393.5135 | 405.7439 | 5 | 2.5 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 397.1573 | 411.2624 | 416.2518 | 5 | 2.4 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 57.7453 | 61.5403 | 66.3578 | 5 | 16.2 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 56.7261 | 59.5980 | 60.3911 | 5 | 16.8 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 125.0468 | 136.3275 | 143.4675 | 5 | 7.3 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 135.2320 | 140.1911 | 143.1420 | 5 | 7.1 |

## Matched blended workload

- phase: `isolated`; target offered rate: 10 updates/s
- offered 0, committed 0, failed 0 over 0.000 s
- offered rate 0.000/s, completed rate 0.000/s
- per key: no keys touched
- ledger: 0 keys reconciled, 0 mismatch

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42130175`, 5 interleaved rounds
  1. family=reuse;policy=mat;refs=1;consumer=1%
  2. family=bag;stmt=cte_joined_outer;threshold=2000
  3. family=reuse;policy=notmat;refs=16;consumer=1%
  4. family=bag;stmt=cte_joined_outer;threshold=20
  5. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  6. family=reuse;policy=default;refs=16;consumer=100%
  7. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  8. family=reuse;policy=mat;refs=16;consumer=100%
  9. family=bag;stmt=join_where;threshold=20
  10. family=reuse;policy=mat;refs=4;consumer=1%
  11. family=reuse;policy=default;refs=4;consumer=100%
  12. family=reuse;policy=default;refs=1;consumer=1%
  13. family=reuse;policy=mat;refs=16;consumer=1%
  14. family=reuse;policy=default;refs=1;consumer=100%
  15. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  16. family=reuse;policy=notmat;refs=1;consumer=1%
  17. family=bag;stmt=join_where;threshold=180000
  18. family=reuse;policy=mat;refs=4;consumer=100%
  19. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  20. family=reuse;policy=notmat;refs=16;consumer=100%
  21. family=bag;stmt=cte_joined_outer;threshold=180000
  22. family=reuse;policy=notmat;refs=1;consumer=100%
  23. family=reuse;policy=mat;refs=1;consumer=100%
  24. family=reuse;policy=notmat;refs=4;consumer=1%
  25. family=reuse;policy=notmat;refs=4;consumer=100%
  26. family=reuse;policy=default;refs=4;consumer=1%
  27. family=bag;stmt=cte_joined_outer_mat;threshold=20
  28. family=reuse;policy=default;refs=16;consumer=1%
  29. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  30. family=bag;stmt=join_where;threshold=2000

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
| family=bag;stmt=cte_joined_outer;threshold=180000 | j05_cte_joined_outer | valid | `d37c12318cec9570` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20 | j05_cte_joined_outer | valid | `1963bacba7f71e91` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=2000 | j05_cte_joined_outer | valid | `9fc8c95b9bdbd482` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | j05m_cte_joined_outer_mat | valid | `89bf2bdf4475e89d` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | j05m_cte_joined_outer_mat | valid | `aec681a14ad951eb` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | j05m_cte_joined_outer_mat | valid | `a9811440bc2d374d` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | j05n_cte_joined_outer_notmat | valid | `234161870ea3b34c` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | j05n_cte_joined_outer_notmat | valid | `1d72690fce2049ca` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | j05n_cte_joined_outer_notmat | valid | `7818b83b0a96bcf5` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=join_where;threshold=180000 | j01_join_where | valid | `488d7ea79629fa97` |  |
| family=bag;stmt=join_where;threshold=20 | j01_join_where | valid | `1e2a0aa383b52218` |  |
| family=bag;stmt=join_where;threshold=2000 | j01_join_where | valid | `6b2f277f4ef05845` |  |
| family=reuse;policy=default;refs=16;consumer=1% | r_default_n16_c1 | valid | `fbfa91d9d372f834` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=16;consumer=100% | r_default_n16_c100 | valid | `53786e8fbd59d292` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=1;consumer=1% | r_default_n1_c1 | valid | `51d0cd00b995aaa4` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=1;consumer=100% | r_default_n1_c100 | valid | `870cb6cde42ccd8c` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=4;consumer=1% | r_default_n4_c1 | valid | `61a80c16c3827eeb` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=100% | r_default_n4_c100 | valid | `03aebc676846112f` | policy=default; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=1% | r_mat_n16_c1 | valid | `e45a3ed613cb52eb` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=100% | r_mat_n16_c100 | valid | `dffebeb20ecbb208` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=1% | r_mat_n1_c1 | valid | `165ec7fe91e76792` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=100% | r_mat_n1_c100 | valid | `1c989baeeb638c5a` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=1% | r_mat_n4_c1 | valid | `f93247b1cf3d3edb` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=100% | r_mat_n4_c100 | valid | `8dad9b965389390e` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=notmat;refs=16;consumer=1% | r_notmat_n16_c1 | valid | `6d2f0cf4f2a6aa8e` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `1f5deaaa277c1da2` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `44e8e67d929c0172` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `7e99c0ea2bf5f956` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `eda3480b42d37e2f` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `1dadc273273c50b5` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, DIST, DEBUG, FORMAT TEXT` supported on yugabyte
- one host, single client: no independent hardware, real-network or equal-total-budget comparison
- no concurrent writers: the isolated read contract only; a matched blended read/write phase is a recorded gap
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
