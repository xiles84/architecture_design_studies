# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 30 cells, 30 valid, engine `yugabyte`, topology `single`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-topo-yb-single-i1-20261004`
- run tag: `run/study-01/select-join-topo-yb-single-i1-20261004`
- environment: `host-zenbook-ux5406sa`
- repo commit: `656489ae7d67583b0c3ba6dbb0808051e9a9cb32`
- engine version: `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit`
- database engine image: `docker.io/yugabytedb/yugabyte:2025.2.6.0-b111` @ `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`
- benchmark client image: `localhost/charitybench-select:select-join-topo-yb-single-i1-20261004` @ `d8db2fd98a7bc00fa2a5561ed01373926c9d950a8d93f0ae9ba4d53be68d6e38`
- inputs digest: `703f2b33238eef2a48030da5cc54b7e81cbbbfdfa6177b3dbb94b90f86a6e6db`
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
| 20 (0.01%) | 56.7329 | missing | missing | missing | missing | missing | 54.5759 | 84.7686 | 54.7661 | missing |
| 2000 (1%) | 57.8108 | missing | missing | missing | missing | missing | 57.1971 | 85.3482 | 58.5522 | missing |
| 180000 (90%) | 119.9078 | missing | missing | missing | missing | missing | 115.3082 | 115.5422 | 111.8797 | missing |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | cte_joined_outer vs join_where | 56.7329 | 54.5759 | 0.9620 | false |
| 20 | cte_joined_outer_mat vs join_where | 56.7329 | 84.7686 | 1.4942 | false |
| 20 | cte_joined_outer_notmat vs join_where | 56.7329 | 54.7661 | 0.9653 | false |
| 2000 | cte_joined_outer vs join_where | 57.8108 | 57.1971 | 0.9894 | false |
| 2000 | cte_joined_outer_mat vs join_where | 57.8108 | 85.3482 | 1.4763 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 57.8108 | 58.5522 | 1.0128 | false |
| 180000 | cte_joined_outer vs join_where | 119.9078 | 115.3082 | 0.9616 | false |
| 180000 | cte_joined_outer_mat vs join_where | 119.9078 | 115.5422 | 0.9636 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 119.9078 | 111.8797 | 0.9330 | false |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_joined_outer_mat vs cte_joined_outer | 54.5759 | 84.7686 | 1.5532 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 54.5759 | 54.7661 | 1.0035 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 57.1971 | 85.3482 | 1.4922 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 57.1971 | 58.5522 | 1.0237 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 115.3082 | 115.5422 | 1.0020 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 115.3082 | 111.8797 | 0.9703 | false | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 60.3669 | 1 | `95a949c4fa06fcb6` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 56.9828 | 1 | `1b3ea281dc6417ad` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | missing | | | | |
| default | 2 | 100% | missing | | | | |
| default | 4 | 1% | 54.9054 | 4 | `8eca22aeced63efd` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 55.5775 | 4 | `27b8cc72dd629518` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | missing | | | | |
| default | 8 | 100% | missing | | | | |
| default | 16 | 1% | 57.1951 | 16 | `700b4c5ce6676793` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 57.3638 | 16 | `607ecb87d7d729b5` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 53.3344 | 1 | `88ee5a1cade66632` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 53.3539 | 1 | `addd6c6e0e963daf` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | missing | | | | |
| mat | 2 | 100% | missing | | | | |
| mat | 4 | 1% | 56.3732 | 4 | `2c6fea0459bc2186` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 53.5907 | 4 | `041d130310da2037` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | missing | | | | |
| mat | 8 | 100% | missing | | | | |
| mat | 16 | 1% | 56.9096 | 16 | `859c5c71265aac3b` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 54.1710 | 16 | `e0ca1773b0f25dd2` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 55.0814 | 1 | `05b2f54b3de34b99` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 61.7294 | 1 | `8e7b0f2b6410988a` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | missing | | | | |
| notmat | 2 | 100% | missing | | | | |
| notmat | 4 | 1% | 160.9036 | 4 | `87447c128f03ee5c` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 179.9540 | 4 | `d063bfcf476f7a0f` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | missing | | | | |
| notmat | 8 | 100% | missing | | | | |
| notmat | 16 | 1% | 559.5926 | 16 | `5f8e00ca79b9e1d0` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 577.3888 | 16 | `924cb6cdeeffa35d` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 0.8835 | false |
| notmat | 1 | 1% | 0.9124 | false |
| mat | 1 | 100% | 0.9363 | false |
| notmat | 1 | 100% | 1.0833 | false |
| mat | 4 | 1% | 1.0267 | false |
| notmat | 4 | 1% | 2.9306 | false |
| mat | 4 | 100% | 0.9643 | false |
| notmat | 4 | 100% | 3.2379 | false |
| mat | 16 | 1% | 0.9950 | false |
| notmat | 16 | 1% | 9.7839 | false |
| mat | 16 | 100% | 0.9443 | false |
| notmat | 16 | 100% | 10.0654 | false |

## Latency samples and derived throughput

_The reader is sequential, so `queries/s` is 1000/median_ms, a derived rate and not a concurrent-capacity test. min/max are the five sample extremes._

| axis | status | min ms | median ms | max ms | samples | derived queries/s |
|---|---|---|---|---|---|---|
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 114.4330 | 115.3082 | 117.3172 | 5 | 8.7 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 54.0829 | 54.5759 | 56.7780 | 5 | 18.3 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 55.6406 | 57.1971 | 75.4124 | 5 | 17.5 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 105.5467 | 115.5422 | 143.5445 | 5 | 8.7 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 83.1423 | 84.7686 | 94.3441 | 5 | 11.8 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 82.9552 | 85.3482 | 88.2876 | 5 | 11.7 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 108.6284 | 111.8797 | 123.2964 | 5 | 8.9 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 52.2647 | 54.7661 | 97.9733 | 5 | 18.3 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 52.4016 | 58.5522 | 62.3139 | 5 | 17.1 |
| family=bag;stmt=join_where;threshold=180000 | valid | 108.2402 | 119.9078 | 135.0832 | 5 | 8.3 |
| family=bag;stmt=join_where;threshold=20 | valid | 56.1672 | 56.7329 | 59.4671 | 5 | 17.6 |
| family=bag;stmt=join_where;threshold=2000 | valid | 52.1503 | 57.8108 | 59.2708 | 5 | 17.3 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 54.6119 | 57.1951 | 67.6262 | 5 | 17.5 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 54.3827 | 57.3638 | 69.9581 | 5 | 17.4 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 54.5676 | 60.3669 | 62.0646 | 5 | 16.6 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 53.8011 | 56.9828 | 61.4589 | 5 | 17.5 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 53.2391 | 54.9054 | 68.2415 | 5 | 18.2 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 52.0608 | 55.5775 | 87.4568 | 5 | 18.0 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 51.8311 | 56.9096 | 64.1810 | 5 | 17.6 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 52.8203 | 54.1710 | 76.5053 | 5 | 18.5 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 50.9219 | 53.3344 | 56.0932 | 5 | 18.7 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 53.0709 | 53.3539 | 69.7710 | 5 | 18.7 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 53.8233 | 56.3732 | 58.7858 | 5 | 17.7 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 53.2987 | 53.5907 | 55.6503 | 5 | 18.7 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 545.2323 | 559.5926 | 714.9737 | 5 | 1.8 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 529.1783 | 577.3888 | 753.2869 | 5 | 1.7 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 52.5223 | 55.0814 | 63.3808 | 5 | 18.2 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 55.0350 | 61.7294 | 65.3789 | 5 | 16.2 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 126.7658 | 160.9036 | 176.6728 | 5 | 6.2 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 150.3001 | 179.9540 | 191.4865 | 5 | 5.6 |

## Matched blended workload

- phase: `isolated`; target offered rate: 10 updates/s
- offered 0, committed 0, failed 0 over 0.000 s
- offered rate 0.000/s, completed rate 0.000/s
- per key: no keys touched
- ledger: 0 keys reconciled, 0 mismatch

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42229993`, 5 interleaved rounds
  1. family=bag;stmt=cte_joined_outer_mat;threshold=20
  2. family=reuse;policy=default;refs=1;consumer=1%
  3. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  4. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  5. family=reuse;policy=notmat;refs=4;consumer=100%
  6. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  7. family=reuse;policy=mat;refs=1;consumer=1%
  8. family=reuse;policy=notmat;refs=1;consumer=1%
  9. family=reuse;policy=mat;refs=4;consumer=1%
  10. family=reuse;policy=default;refs=1;consumer=100%
  11. family=bag;stmt=join_where;threshold=2000
  12. family=bag;stmt=join_where;threshold=20
  13. family=reuse;policy=mat;refs=16;consumer=1%
  14. family=reuse;policy=notmat;refs=1;consumer=100%
  15. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  16. family=bag;stmt=cte_joined_outer;threshold=180000
  17. family=reuse;policy=mat;refs=4;consumer=100%
  18. family=bag;stmt=cte_joined_outer;threshold=20
  19. family=reuse;policy=notmat;refs=4;consumer=1%
  20. family=reuse;policy=mat;refs=16;consumer=100%
  21. family=reuse;policy=notmat;refs=16;consumer=1%
  22. family=reuse;policy=default;refs=16;consumer=100%
  23. family=reuse;policy=mat;refs=1;consumer=100%
  24. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  25. family=bag;stmt=cte_joined_outer;threshold=2000
  26. family=bag;stmt=join_where;threshold=180000
  27. family=reuse;policy=notmat;refs=16;consumer=100%
  28. family=reuse;policy=default;refs=4;consumer=100%
  29. family=reuse;policy=default;refs=4;consumer=1%
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
| family=bag;stmt=cte_joined_outer;threshold=180000 | j05_cte_joined_outer | valid | `befb7be37350fb98` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20 | j05_cte_joined_outer | valid | `0f8e66011fa37d32` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=2000 | j05_cte_joined_outer | valid | `524af04e3c7a5c13` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | j05m_cte_joined_outer_mat | valid | `d7cb797f306fdd6c` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | j05m_cte_joined_outer_mat | valid | `e60fe5aea707a998` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | j05m_cte_joined_outer_mat | valid | `8d0228d45791be34` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | j05n_cte_joined_outer_notmat | valid | `e644fe24616278ed` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | j05n_cte_joined_outer_notmat | valid | `23c5ae9f7a98be1a` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | j05n_cte_joined_outer_notmat | valid | `545404b2a16956e6` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=join_where;threshold=180000 | j01_join_where | valid | `5646db58cfda33b8` |  |
| family=bag;stmt=join_where;threshold=20 | j01_join_where | valid | `351c6751034e55c4` |  |
| family=bag;stmt=join_where;threshold=2000 | j01_join_where | valid | `10194ccf6f04d761` |  |
| family=reuse;policy=default;refs=16;consumer=1% | r_default_n16_c1 | valid | `700b4c5ce6676793` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=16;consumer=100% | r_default_n16_c100 | valid | `607ecb87d7d729b5` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=1;consumer=1% | r_default_n1_c1 | valid | `95a949c4fa06fcb6` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=1;consumer=100% | r_default_n1_c100 | valid | `1b3ea281dc6417ad` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=4;consumer=1% | r_default_n4_c1 | valid | `8eca22aeced63efd` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=100% | r_default_n4_c100 | valid | `27b8cc72dd629518` | policy=default; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=1% | r_mat_n16_c1 | valid | `859c5c71265aac3b` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=100% | r_mat_n16_c100 | valid | `e0ca1773b0f25dd2` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=1% | r_mat_n1_c1 | valid | `88ee5a1cade66632` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=100% | r_mat_n1_c100 | valid | `addd6c6e0e963daf` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=1% | r_mat_n4_c1 | valid | `2c6fea0459bc2186` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=100% | r_mat_n4_c100 | valid | `041d130310da2037` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=notmat;refs=16;consumer=1% | r_notmat_n16_c1 | valid | `5f8e00ca79b9e1d0` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `924cb6cdeeffa35d` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `05b2f54b3de34b99` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `8e7b0f2b6410988a` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `87447c128f03ee5c` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `d063bfcf476f7a0f` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, DIST, DEBUG, FORMAT TEXT` supported on yugabyte
- one host, single client: no independent hardware, real-network or equal-total-budget comparison
- no concurrent writers: the isolated read contract only; a matched blended read/write phase is a recorded gap
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
