# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 30 cells, 30 valid, engine `yugabyte`, topology `single`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-topo-yb-equal-bo2-20261004`
- run tag: `run/study-01/select-join-topo-yb-equal-bo2-20261004`
- environment: `host-zenbook-ux5406sa`
- repo commit: `455b777fc02e015efdd9e54d25d3aa3343eb206f`
- engine version: `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit`
- database engine image: `docker.io/yugabytedb/yugabyte:2025.2.6.0-b111` @ `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`
- benchmark client image: `localhost/charitybench-select:select-join-topo-yb-equal-bo2-20261004` @ `79bc753ff7e57e102982a323939c43e68e9b3664887683d3a5d5bf072650f4c1`
- inputs digest: `b98effbf8a83f961bec02d4eabc15db64ca6923d82f78d85da39c34bfa396ea4`
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
| 20 (0.01%) | 50.3190 | missing | missing | missing | missing | missing | 53.9288 | 92.4960 | 49.1918 | missing |
| 2000 (1%) | 50.2280 | missing | missing | missing | missing | missing | 51.7414 | 77.4101 | 50.2477 | missing |
| 180000 (90%) | 129.3549 | missing | missing | missing | missing | missing | 117.4256 | 113.7144 | 132.1863 | missing |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | cte_joined_outer vs join_where | 50.3190 | 53.9288 | 1.0717 | false |
| 20 | cte_joined_outer_mat vs join_where | 50.3190 | 92.4960 | 1.8382 | false |
| 20 | cte_joined_outer_notmat vs join_where | 50.3190 | 49.1918 | 0.9776 | false |
| 2000 | cte_joined_outer vs join_where | 50.2280 | 51.7414 | 1.0301 | false |
| 2000 | cte_joined_outer_mat vs join_where | 50.2280 | 77.4101 | 1.5412 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 50.2280 | 50.2477 | 1.0004 | false |
| 180000 | cte_joined_outer vs join_where | 129.3549 | 117.4256 | 0.9078 | false |
| 180000 | cte_joined_outer_mat vs join_where | 129.3549 | 113.7144 | 0.8791 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 129.3549 | 132.1863 | 1.0219 | false |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_joined_outer_mat vs cte_joined_outer | 53.9288 | 92.4960 | 1.7152 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 53.9288 | 49.1918 | 0.9122 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 51.7414 | 77.4101 | 1.4961 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 51.7414 | 50.2477 | 0.9711 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 117.4256 | 113.7144 | 0.9684 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 117.4256 | 132.1863 | 1.1257 | false | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 60.7967 | 1 | `62faf8c055f9ea3c` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 54.3752 | 1 | `3d0a338832f4a76e` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | missing | | | | |
| default | 2 | 100% | missing | | | | |
| default | 4 | 1% | 56.0266 | 4 | `c3d5aa9bc7a670db` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 49.6663 | 4 | `4572a6c78247f223` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | missing | | | | |
| default | 8 | 100% | missing | | | | |
| default | 16 | 1% | 51.0631 | 16 | `b027886f6a5adcac` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 55.3749 | 16 | `dcc33ef086423230` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 50.8713 | 1 | `d3966f66bd014bba` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 50.9506 | 1 | `01ba3e2d34c91f44` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | missing | | | | |
| mat | 2 | 100% | missing | | | | |
| mat | 4 | 1% | 54.6479 | 4 | `f2145295b617708f` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 51.1746 | 4 | `d87e4c005d34d358` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | missing | | | | |
| mat | 8 | 100% | missing | | | | |
| mat | 16 | 1% | 50.7252 | 16 | `64aad2dcfde1efe6` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 50.5845 | 16 | `b10bc23f382f98c9` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 52.6655 | 1 | `441f0fbfa0b87a4d` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 52.8947 | 1 | `11f4a5bb5ebb0c49` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | missing | | | | |
| notmat | 2 | 100% | missing | | | | |
| notmat | 4 | 1% | 105.9753 | 4 | `77fd805dc8c1f669` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 121.4790 | 4 | `36de323fbf216efa` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | missing | | | | |
| notmat | 8 | 100% | missing | | | | |
| notmat | 16 | 1% | 441.0225 | 16 | `d060fc721b606bfd` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 382.3463 | 16 | `34153325709d2320` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 0.8367 | false |
| notmat | 1 | 1% | 0.8663 | false |
| mat | 1 | 100% | 0.9370 | false |
| notmat | 1 | 100% | 0.9728 | false |
| mat | 4 | 1% | 0.9754 | false |
| notmat | 4 | 1% | 1.8915 | false |
| mat | 4 | 100% | 1.0304 | false |
| notmat | 4 | 100% | 2.4459 | false |
| mat | 16 | 1% | 0.9934 | false |
| notmat | 16 | 1% | 8.6368 | false |
| mat | 16 | 100% | 0.9135 | false |
| notmat | 16 | 100% | 6.9047 | false |

## Latency samples and derived throughput

_The reader is sequential, so `queries/s` is 1000/median_ms, a derived rate and not a concurrent-capacity test. min/max are the five sample extremes._

| axis | status | min ms | median ms | max ms | samples | derived queries/s |
|---|---|---|---|---|---|---|
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 109.5514 | 117.4256 | 126.9095 | 5 | 8.5 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 53.0086 | 53.9288 | 76.4957 | 5 | 18.5 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 49.7033 | 51.7414 | 60.3163 | 5 | 19.3 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 106.9880 | 113.7144 | 115.0102 | 5 | 8.8 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 72.7297 | 92.4960 | 103.2910 | 5 | 10.8 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 75.3151 | 77.4101 | 93.0208 | 5 | 12.9 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 115.3494 | 132.1863 | 155.8511 | 5 | 7.6 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 48.4832 | 49.1918 | 60.6038 | 5 | 20.3 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 48.9894 | 50.2477 | 60.5331 | 5 | 19.9 |
| family=bag;stmt=join_where;threshold=180000 | valid | 119.6155 | 129.3549 | 153.5304 | 5 | 7.7 |
| family=bag;stmt=join_where;threshold=20 | valid | 49.0546 | 50.3190 | 62.4413 | 5 | 19.9 |
| family=bag;stmt=join_where;threshold=2000 | valid | 47.7684 | 50.2280 | 62.3433 | 5 | 19.9 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 50.3927 | 51.0631 | 63.6440 | 5 | 19.6 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 53.1364 | 55.3749 | 76.9963 | 5 | 18.1 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 53.3693 | 60.7967 | 68.6374 | 5 | 16.4 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 51.3848 | 54.3752 | 67.1975 | 5 | 18.4 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 54.2018 | 56.0266 | 71.8133 | 5 | 17.8 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 47.9746 | 49.6663 | 61.7117 | 5 | 20.1 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 49.6613 | 50.7252 | 68.7899 | 5 | 19.7 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 48.7621 | 50.5845 | 65.0665 | 5 | 19.8 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 48.8659 | 50.8713 | 63.7157 | 5 | 19.7 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 47.8798 | 50.9506 | 65.3460 | 5 | 19.6 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 49.4289 | 54.6479 | 67.9841 | 5 | 18.3 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 49.9855 | 51.1746 | 63.2683 | 5 | 19.5 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 356.0622 | 441.0225 | 522.9900 | 5 | 2.3 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 376.2778 | 382.3463 | 461.0048 | 5 | 2.6 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 49.6592 | 52.6655 | 67.2806 | 5 | 19.0 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 50.7126 | 52.8947 | 65.7460 | 5 | 18.9 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 105.1213 | 105.9753 | 260.8279 | 5 | 9.4 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 104.4472 | 121.4790 | 150.4365 | 5 | 8.2 |

## Matched blended workload

- phase: `blended`; target offered rate: 10 updates/s
- offered 308, committed 308, failed 0 over 30.776 s
- offered rate 10.008/s, completed rate 10.008/s
- reader operations overlapping a committed write: 154 of 213
- reader retries of transient conflicts: 3
- per key: 308 keys, min 1, max 1 per key
- ledger: 308 keys reconciled, 0 mismatch
- update latency: avg 3.672 ms, min 1.779, max 56.075 over 308

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42908404`, 5 interleaved rounds
  1. family=reuse;policy=mat;refs=16;consumer=100%
  2. family=reuse;policy=notmat;refs=4;consumer=100%
  3. family=reuse;policy=mat;refs=1;consumer=1%
  4. family=reuse;policy=notmat;refs=1;consumer=1%
  5. family=reuse;policy=mat;refs=16;consumer=1%
  6. family=reuse;policy=mat;refs=4;consumer=100%
  7. family=reuse;policy=default;refs=16;consumer=1%
  8. family=bag;stmt=cte_joined_outer;threshold=180000
  9. family=bag;stmt=cte_joined_outer_mat;threshold=20
  10. family=reuse;policy=notmat;refs=16;consumer=1%
  11. family=reuse;policy=default;refs=4;consumer=1%
  12. family=reuse;policy=default;refs=1;consumer=1%
  13. family=reuse;policy=notmat;refs=16;consumer=100%
  14. family=bag;stmt=join_where;threshold=180000
  15. family=reuse;policy=default;refs=16;consumer=100%
  16. family=bag;stmt=cte_joined_outer;threshold=20
  17. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  18. family=reuse;policy=mat;refs=4;consumer=1%
  19. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  20. family=bag;stmt=cte_joined_outer;threshold=2000
  21. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  22. family=reuse;policy=notmat;refs=1;consumer=100%
  23. family=reuse;policy=default;refs=1;consumer=100%
  24. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  25. family=reuse;policy=notmat;refs=4;consumer=1%
  26. family=reuse;policy=mat;refs=1;consumer=100%
  27. family=bag;stmt=join_where;threshold=20
  28. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  29. family=reuse;policy=default;refs=4;consumer=100%
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
| family=bag;stmt=cte_joined_outer;threshold=180000 | j05_cte_joined_outer | valid | `2122d4eba97ddb51` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20 | j05_cte_joined_outer | valid | `56d3c8de557bb0b6` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=2000 | j05_cte_joined_outer | valid | `e7f5b2ec95a5a013` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | j05m_cte_joined_outer_mat | valid | `7b938c04446e8e39` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | j05m_cte_joined_outer_mat | valid | `fdee9ea5d17a2d4e` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | j05m_cte_joined_outer_mat | valid | `921c086268f499a3` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | j05n_cte_joined_outer_notmat | valid | `0351bda8a8ac1f69` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | j05n_cte_joined_outer_notmat | valid | `65e7b07bc65ed158` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | j05n_cte_joined_outer_notmat | valid | `9e5d7d66a8198fe5` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=join_where;threshold=180000 | j01_join_where | valid | `196f0df181e68977` |  |
| family=bag;stmt=join_where;threshold=20 | j01_join_where | valid | `07f369542bb060ca` |  |
| family=bag;stmt=join_where;threshold=2000 | j01_join_where | valid | `b2f33a9b7cc6113b` |  |
| family=reuse;policy=default;refs=16;consumer=1% | r_default_n16_c1 | valid | `b027886f6a5adcac` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=16;consumer=100% | r_default_n16_c100 | valid | `dcc33ef086423230` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=1;consumer=1% | r_default_n1_c1 | valid | `62faf8c055f9ea3c` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=1;consumer=100% | r_default_n1_c100 | valid | `3d0a338832f4a76e` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=4;consumer=1% | r_default_n4_c1 | valid | `c3d5aa9bc7a670db` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=100% | r_default_n4_c100 | valid | `4572a6c78247f223` | policy=default; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=1% | r_mat_n16_c1 | valid | `64aad2dcfde1efe6` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=100% | r_mat_n16_c100 | valid | `b10bc23f382f98c9` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=1% | r_mat_n1_c1 | valid | `d3966f66bd014bba` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=100% | r_mat_n1_c100 | valid | `01ba3e2d34c91f44` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=1% | r_mat_n4_c1 | valid | `f2145295b617708f` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=100% | r_mat_n4_c100 | valid | `d87e4c005d34d358` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=notmat;refs=16;consumer=1% | r_notmat_n16_c1 | valid | `d060fc721b606bfd` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `34153325709d2320` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `441f0fbfa0b87a4d` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `11f4a5bb5ebb0c49` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `77fd805dc8c1f669` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `36de323fbf216efa` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, DIST, DEBUG, FORMAT TEXT` supported on yugabyte
- one host, single client: no independent hardware, real-network or equal-total-budget comparison
- no concurrent writers: the isolated read contract only; a matched blended read/write phase is a recorded gap
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
