# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 30 cells, 30 valid, engine `yugabyte`, topology `single`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-topo-yb-equal-i1-20261004`
- run tag: `run/study-01/select-join-topo-yb-equal-i1-20261004`
- environment: `host-zenbook-ux5406sa`
- repo commit: `335542f10685ca626b2f4cb2b086f0ad67e620e0`
- engine version: `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit`
- database engine image: `docker.io/yugabytedb/yugabyte:2025.2.6.0-b111` @ `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`
- benchmark client image: `localhost/charitybench-select:select-join-topo-yb-equal-i1-20261004` @ `d8db2fd98a7bc00fa2a5561ed01373926c9d950a8d93f0ae9ba4d53be68d6e38`
- inputs digest: `0facab5446737c4c55db75ead528663c33642a667c649cb454dda69334e3aab3`
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
| 20 (0.01%) | 55.4153 | missing | missing | missing | missing | missing | 57.7261 | 85.2691 | 57.9434 | missing |
| 2000 (1%) | 56.4833 | missing | missing | missing | missing | missing | 57.7575 | 88.1047 | 56.3890 | missing |
| 180000 (90%) | 119.2619 | missing | missing | missing | missing | missing | 115.0769 | 120.7959 | 116.7473 | missing |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | cte_joined_outer vs join_where | 55.4153 | 57.7261 | 1.0417 | false |
| 20 | cte_joined_outer_mat vs join_where | 55.4153 | 85.2691 | 1.5387 | false |
| 20 | cte_joined_outer_notmat vs join_where | 55.4153 | 57.9434 | 1.0456 | false |
| 2000 | cte_joined_outer vs join_where | 56.4833 | 57.7575 | 1.0226 | false |
| 2000 | cte_joined_outer_mat vs join_where | 56.4833 | 88.1047 | 1.5598 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 56.4833 | 56.3890 | 0.9983 | false |
| 180000 | cte_joined_outer vs join_where | 119.2619 | 115.0769 | 0.9649 | false |
| 180000 | cte_joined_outer_mat vs join_where | 119.2619 | 120.7959 | 1.0129 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 119.2619 | 116.7473 | 0.9789 | false |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_joined_outer_mat vs cte_joined_outer | 57.7261 | 85.2691 | 1.4771 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 57.7261 | 57.9434 | 1.0038 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 57.7575 | 88.1047 | 1.5254 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 57.7575 | 56.3890 | 0.9763 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 115.0769 | 120.7959 | 1.0497 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 115.0769 | 116.7473 | 1.0145 | false | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 58.9152 | 1 | `4300c17bd101e461` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 57.0732 | 1 | `93911ed8c16be1eb` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | missing | | | | |
| default | 2 | 100% | missing | | | | |
| default | 4 | 1% | 56.8895 | 4 | `39db4bb22b7b3b2a` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 53.7868 | 4 | `6cc0cbd92dd0f472` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | missing | | | | |
| default | 8 | 100% | missing | | | | |
| default | 16 | 1% | 57.1678 | 16 | `c356ecd483547723` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 56.9487 | 16 | `fef505acdf7c8068` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 55.7032 | 1 | `2874af277bcff415` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 57.7261 | 1 | `1e0303ddfe4a1610` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | missing | | | | |
| mat | 2 | 100% | missing | | | | |
| mat | 4 | 1% | 57.3291 | 4 | `60e582da9a012120` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 53.7189 | 4 | `da508b195efd51bd` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | missing | | | | |
| mat | 8 | 100% | missing | | | | |
| mat | 16 | 1% | 56.0107 | 16 | `063b4de3b1b53c23` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 56.2572 | 16 | `a2f9a4144ebe238d` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 58.3360 | 1 | `e90ac86fe1de0470` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 57.9457 | 1 | `6d06dde8389461cd` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | missing | | | | |
| notmat | 2 | 100% | missing | | | | |
| notmat | 4 | 1% | 123.4679 | 4 | `217c8b3add9789c9` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 124.7187 | 4 | `6d682ba47052984a` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | missing | | | | |
| notmat | 8 | 100% | missing | | | | |
| notmat | 16 | 1% | 376.9908 | 16 | `1363c0371eefbd29` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 382.4606 | 16 | `bb826f2dcd8900e0` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 0.9455 | false |
| notmat | 1 | 1% | 0.9902 | false |
| mat | 1 | 100% | 1.0114 | false |
| notmat | 1 | 100% | 1.0153 | false |
| mat | 4 | 1% | 1.0077 | false |
| notmat | 4 | 1% | 2.1703 | false |
| mat | 4 | 100% | 0.9987 | false |
| notmat | 4 | 100% | 2.3188 | false |
| mat | 16 | 1% | 0.9798 | false |
| notmat | 16 | 1% | 6.5945 | false |
| mat | 16 | 100% | 0.9879 | false |
| notmat | 16 | 100% | 6.7159 | false |

## Latency samples and derived throughput

_The reader is sequential, so `queries/s` is 1000/median_ms, a derived rate and not a concurrent-capacity test. min/max are the five sample extremes._

| axis | status | min ms | median ms | max ms | samples | derived queries/s |
|---|---|---|---|---|---|---|
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 112.3864 | 115.0769 | 129.0200 | 5 | 8.7 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 54.9658 | 57.7261 | 60.0730 | 5 | 17.3 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 55.5181 | 57.7575 | 58.4926 | 5 | 17.3 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 109.4299 | 120.7959 | 122.8606 | 5 | 8.3 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 81.6448 | 85.2691 | 101.9207 | 5 | 11.7 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 83.2018 | 88.1047 | 98.4894 | 5 | 11.4 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 111.7603 | 116.7473 | 139.6827 | 5 | 8.6 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 54.6655 | 57.9434 | 66.1686 | 5 | 17.3 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 53.7915 | 56.3890 | 62.3529 | 5 | 17.7 |
| family=bag;stmt=join_where;threshold=180000 | valid | 112.8504 | 119.2619 | 120.0822 | 5 | 8.4 |
| family=bag;stmt=join_where;threshold=20 | valid | 53.7863 | 55.4153 | 57.7596 | 5 | 18.0 |
| family=bag;stmt=join_where;threshold=2000 | valid | 54.8051 | 56.4833 | 60.5152 | 5 | 17.7 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 51.7155 | 57.1678 | 64.4444 | 5 | 17.5 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 55.7873 | 56.9487 | 64.8504 | 5 | 17.6 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 54.0182 | 58.9152 | 62.0529 | 5 | 17.0 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 55.2651 | 57.0732 | 61.7902 | 5 | 17.5 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 54.0756 | 56.8895 | 58.9007 | 5 | 17.6 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 49.3665 | 53.7868 | 57.4548 | 5 | 18.6 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 55.4460 | 56.0107 | 56.8849 | 5 | 17.9 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 53.7680 | 56.2572 | 59.6467 | 5 | 17.8 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 53.0520 | 55.7032 | 57.6006 | 5 | 18.0 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 53.3666 | 57.7261 | 59.7695 | 5 | 17.3 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 55.0563 | 57.3291 | 62.7102 | 5 | 17.4 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 51.4646 | 53.7189 | 58.9877 | 5 | 18.6 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 366.9918 | 376.9908 | 410.8095 | 5 | 2.7 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 376.3942 | 382.4606 | 387.5783 | 5 | 2.6 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 57.1765 | 58.3360 | 67.9587 | 5 | 17.1 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 55.2456 | 57.9457 | 62.4508 | 5 | 17.3 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 121.7159 | 123.4679 | 127.0250 | 5 | 8.1 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 123.4758 | 124.7187 | 136.1180 | 5 | 8.0 |

## Matched blended workload

- phase: `isolated`; target offered rate: 10 updates/s
- offered 0, committed 0, failed 0 over 0.000 s
- offered rate 0.000/s, completed rate 0.000/s
- per key: no keys touched
- ledger: 0 keys reconciled, 0 mismatch

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42303909`, 5 interleaved rounds
  1. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  2. family=reuse;policy=notmat;refs=4;consumer=100%
  3. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  4. family=reuse;policy=notmat;refs=16;consumer=1%
  5. family=reuse;policy=notmat;refs=1;consumer=1%
  6. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  7. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  8. family=bag;stmt=cte_joined_outer_mat;threshold=20
  9. family=reuse;policy=mat;refs=1;consumer=100%
  10. family=reuse;policy=notmat;refs=1;consumer=100%
  11. family=reuse;policy=default;refs=4;consumer=1%
  12. family=bag;stmt=cte_joined_outer;threshold=20
  13. family=bag;stmt=join_where;threshold=180000
  14. family=bag;stmt=join_where;threshold=20
  15. family=reuse;policy=notmat;refs=4;consumer=1%
  16. family=reuse;policy=mat;refs=1;consumer=1%
  17. family=reuse;policy=mat;refs=4;consumer=100%
  18. family=reuse;policy=mat;refs=16;consumer=100%
  19. family=reuse;policy=default;refs=4;consumer=100%
  20. family=reuse;policy=notmat;refs=16;consumer=100%
  21. family=bag;stmt=cte_joined_outer;threshold=2000
  22. family=reuse;policy=default;refs=1;consumer=1%
  23. family=reuse;policy=default;refs=16;consumer=100%
  24. family=reuse;policy=mat;refs=4;consumer=1%
  25. family=bag;stmt=cte_joined_outer;threshold=180000
  26. family=bag;stmt=join_where;threshold=2000
  27. family=reuse;policy=default;refs=1;consumer=100%
  28. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  29. family=reuse;policy=mat;refs=16;consumer=1%
  30. family=reuse;policy=default;refs=16;consumer=1%

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
| family=bag;stmt=cte_joined_outer;threshold=180000 | j05_cte_joined_outer | valid | `e7a79f1a81da4f6b` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20 | j05_cte_joined_outer | valid | `d0ab9ff7735ce3f9` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=2000 | j05_cte_joined_outer | valid | `df70b8c79dddca2e` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | j05m_cte_joined_outer_mat | valid | `fb6e0d9a3c27433b` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | j05m_cte_joined_outer_mat | valid | `b6614c6df78aba47` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | j05m_cte_joined_outer_mat | valid | `4d3726e0343ab6a8` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | j05n_cte_joined_outer_notmat | valid | `b15dfbc2c72e0d1c` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | j05n_cte_joined_outer_notmat | valid | `bebd917720fae306` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | j05n_cte_joined_outer_notmat | valid | `85fbf8a26060199f` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=join_where;threshold=180000 | j01_join_where | valid | `19a11a8f318f1782` |  |
| family=bag;stmt=join_where;threshold=20 | j01_join_where | valid | `a8a7a9559a14dc26` |  |
| family=bag;stmt=join_where;threshold=2000 | j01_join_where | valid | `ebb8ff5e290d01d5` |  |
| family=reuse;policy=default;refs=16;consumer=1% | r_default_n16_c1 | valid | `c356ecd483547723` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=16;consumer=100% | r_default_n16_c100 | valid | `fef505acdf7c8068` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=1;consumer=1% | r_default_n1_c1 | valid | `4300c17bd101e461` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=1;consumer=100% | r_default_n1_c100 | valid | `93911ed8c16be1eb` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=4;consumer=1% | r_default_n4_c1 | valid | `39db4bb22b7b3b2a` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=100% | r_default_n4_c100 | valid | `6cc0cbd92dd0f472` | policy=default; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=1% | r_mat_n16_c1 | valid | `063b4de3b1b53c23` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=100% | r_mat_n16_c100 | valid | `a2f9a4144ebe238d` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=1% | r_mat_n1_c1 | valid | `2874af277bcff415` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=100% | r_mat_n1_c100 | valid | `1e0303ddfe4a1610` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=1% | r_mat_n4_c1 | valid | `60e582da9a012120` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=100% | r_mat_n4_c100 | valid | `da508b195efd51bd` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=notmat;refs=16;consumer=1% | r_notmat_n16_c1 | valid | `1363c0371eefbd29` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `bb826f2dcd8900e0` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `e90ac86fe1de0470` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `6d06dde8389461cd` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `217c8b3add9789c9` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `6d682ba47052984a` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, DIST, DEBUG, FORMAT TEXT` supported on yugabyte
- one host, single client: no independent hardware, real-network or equal-total-budget comparison
- no concurrent writers: the isolated read contract only; a matched blended read/write phase is a recorded gap
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
