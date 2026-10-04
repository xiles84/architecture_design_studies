# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 30 cells, 30 valid, engine `yugabyte`, topology `single`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-topo-yb-single-i2-20261004`
- run tag: `run/study-01/select-join-topo-yb-single-i2-20261004`
- environment: `host-zenbook-ux5406sa`
- repo commit: `656489ae7d67583b0c3ba6dbb0808051e9a9cb32`
- engine version: `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit`
- database engine image: `docker.io/yugabytedb/yugabyte:2025.2.6.0-b111` @ `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`
- benchmark client image: `localhost/charitybench-select:select-join-topo-yb-single-i2-20261004` @ `d8db2fd98a7bc00fa2a5561ed01373926c9d950a8d93f0ae9ba4d53be68d6e38`
- inputs digest: `60192ec76d86536bed2998ae40d62d5a2f96635f703d02ab86a17d4ed3e70308`
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
| 20 (0.01%) | 58.9060 | missing | missing | missing | missing | missing | 58.2810 | 102.3405 | 57.5154 | missing |
| 2000 (1%) | 56.8327 | missing | missing | missing | missing | missing | 53.2324 | 88.2535 | 54.8947 | missing |
| 180000 (90%) | 117.9654 | missing | missing | missing | missing | missing | 114.2169 | 111.5070 | 112.3960 | missing |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | cte_joined_outer vs join_where | 58.9060 | 58.2810 | 0.9894 | false |
| 20 | cte_joined_outer_mat vs join_where | 58.9060 | 102.3405 | 1.7374 | false |
| 20 | cte_joined_outer_notmat vs join_where | 58.9060 | 57.5154 | 0.9764 | false |
| 2000 | cte_joined_outer vs join_where | 56.8327 | 53.2324 | 0.9366 | false |
| 2000 | cte_joined_outer_mat vs join_where | 56.8327 | 88.2535 | 1.5529 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 56.8327 | 54.8947 | 0.9659 | false |
| 180000 | cte_joined_outer vs join_where | 117.9654 | 114.2169 | 0.9682 | false |
| 180000 | cte_joined_outer_mat vs join_where | 117.9654 | 111.5070 | 0.9453 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 117.9654 | 112.3960 | 0.9528 | false |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_joined_outer_mat vs cte_joined_outer | 58.2810 | 102.3405 | 1.7560 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 58.2810 | 57.5154 | 0.9869 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 53.2324 | 88.2535 | 1.6579 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 53.2324 | 54.8947 | 1.0312 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 114.2169 | 111.5070 | 0.9763 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 114.2169 | 112.3960 | 0.9841 | false | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 58.4089 | 1 | `9c1b3b3128e73bff` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 59.2367 | 1 | `e1dd4ba22ebf8119` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | missing | | | | |
| default | 2 | 100% | missing | | | | |
| default | 4 | 1% | 56.8430 | 4 | `29dcec76d2411783` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 55.3589 | 4 | `24bdcc8b9c1bc9c9` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | missing | | | | |
| default | 8 | 100% | missing | | | | |
| default | 16 | 1% | 56.8409 | 16 | `883eabbf6932f4f1` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 58.3907 | 16 | `6e6ed7dfb336a674` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 54.7696 | 1 | `1f743cecb2f1f0a6` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 58.2328 | 1 | `be107f6397d551b9` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | missing | | | | |
| mat | 2 | 100% | missing | | | | |
| mat | 4 | 1% | 55.9978 | 4 | `25c46eca56b92a9a` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 57.7456 | 4 | `47bc68e8d34fbd94` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | missing | | | | |
| mat | 8 | 100% | missing | | | | |
| mat | 16 | 1% | 57.3420 | 16 | `3a5cde1136d5063e` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 57.9552 | 16 | `590ac95534dfddcb` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 56.2393 | 1 | `3bf77e70fd286dd8` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 57.5775 | 1 | `56f9f7faf6140aff` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | missing | | | | |
| notmat | 2 | 100% | missing | | | | |
| notmat | 4 | 1% | 150.9267 | 4 | `90ae6f638e7cd7b8` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 141.8750 | 4 | `0654fbf3ddd5ccd3` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | missing | | | | |
| notmat | 8 | 100% | missing | | | | |
| notmat | 16 | 1% | 588.9804 | 16 | `cf39891d205579ea` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 586.5704 | 16 | `d035fa763726b085` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 0.9377 | false |
| notmat | 1 | 1% | 0.9629 | false |
| mat | 1 | 100% | 0.9831 | false |
| notmat | 1 | 100% | 0.9720 | false |
| mat | 4 | 1% | 0.9851 | false |
| notmat | 4 | 1% | 2.6551 | false |
| mat | 4 | 100% | 1.0431 | false |
| notmat | 4 | 100% | 2.5628 | false |
| mat | 16 | 1% | 1.0088 | false |
| notmat | 16 | 1% | 10.3619 | false |
| mat | 16 | 100% | 0.9925 | false |
| notmat | 16 | 100% | 10.0456 | false |

## Latency samples and derived throughput

_The reader is sequential, so `queries/s` is 1000/median_ms, a derived rate and not a concurrent-capacity test. min/max are the five sample extremes._

| axis | status | min ms | median ms | max ms | samples | derived queries/s |
|---|---|---|---|---|---|---|
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 112.0044 | 114.2169 | 116.5648 | 5 | 8.8 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 51.9363 | 58.2810 | 89.1776 | 5 | 17.2 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 48.8595 | 53.2324 | 60.0097 | 5 | 18.8 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 105.7155 | 111.5070 | 115.4972 | 5 | 9.0 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 95.7200 | 102.3405 | 108.2553 | 5 | 9.8 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 82.5820 | 88.2535 | 92.0595 | 5 | 11.3 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 109.4420 | 112.3960 | 118.8621 | 5 | 8.9 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 48.6110 | 57.5154 | 60.8098 | 5 | 17.4 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 53.3506 | 54.8947 | 58.3428 | 5 | 18.2 |
| family=bag;stmt=join_where;threshold=180000 | valid | 112.2073 | 117.9654 | 121.7773 | 5 | 8.5 |
| family=bag;stmt=join_where;threshold=20 | valid | 52.2293 | 58.9060 | 61.3257 | 5 | 17.0 |
| family=bag;stmt=join_where;threshold=2000 | valid | 53.5118 | 56.8327 | 65.3571 | 5 | 17.6 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 53.6785 | 56.8409 | 62.8365 | 5 | 17.6 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 55.6454 | 58.3907 | 61.4211 | 5 | 17.1 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 58.1232 | 58.4089 | 59.0411 | 5 | 17.1 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 57.5175 | 59.2367 | 61.4086 | 5 | 16.9 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 53.0692 | 56.8430 | 66.0475 | 5 | 17.6 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 55.1170 | 55.3589 | 59.1823 | 5 | 18.1 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 57.0657 | 57.3420 | 59.6798 | 5 | 17.4 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 54.7162 | 57.9552 | 60.7831 | 5 | 17.3 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 52.7848 | 54.7696 | 58.2924 | 5 | 18.3 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 53.8098 | 58.2328 | 58.9346 | 5 | 17.2 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 54.3473 | 55.9978 | 56.4242 | 5 | 17.9 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 52.6719 | 57.7456 | 58.0117 | 5 | 17.3 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 586.2938 | 588.9804 | 596.8184 | 5 | 1.7 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 552.6979 | 586.5704 | 622.1084 | 5 | 1.7 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 52.6547 | 56.2393 | 62.4129 | 5 | 17.8 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 56.9686 | 57.5775 | 61.3569 | 5 | 17.4 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 138.2118 | 150.9267 | 156.7416 | 5 | 6.6 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 134.5992 | 141.8750 | 176.2361 | 5 | 7.0 |

## Matched blended workload

- phase: `isolated`; target offered rate: 10 updates/s
- offered 0, committed 0, failed 0 over 0.000 s
- offered rate 0.000/s, completed rate 0.000/s
- per key: no keys touched
- ledger: 0 keys reconciled, 0 mismatch

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42499994`, 5 interleaved rounds
  1. family=reuse;policy=default;refs=1;consumer=1%
  2. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  3. family=bag;stmt=join_where;threshold=20
  4. family=reuse;policy=default;refs=16;consumer=100%
  5. family=bag;stmt=cte_joined_outer;threshold=2000
  6. family=reuse;policy=notmat;refs=1;consumer=100%
  7. family=reuse;policy=notmat;refs=16;consumer=100%
  8. family=bag;stmt=join_where;threshold=2000
  9. family=reuse;policy=default;refs=16;consumer=1%
  10. family=reuse;policy=notmat;refs=1;consumer=1%
  11. family=reuse;policy=notmat;refs=4;consumer=1%
  12. family=reuse;policy=notmat;refs=16;consumer=1%
  13. family=bag;stmt=cte_joined_outer;threshold=20
  14. family=reuse;policy=default;refs=4;consumer=100%
  15. family=reuse;policy=mat;refs=16;consumer=1%
  16. family=reuse;policy=default;refs=4;consumer=1%
  17. family=reuse;policy=mat;refs=16;consumer=100%
  18. family=bag;stmt=join_where;threshold=180000
  19. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  20. family=reuse;policy=default;refs=1;consumer=100%
  21. family=reuse;policy=mat;refs=4;consumer=1%
  22. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  23. family=bag;stmt=cte_joined_outer;threshold=180000
  24. family=reuse;policy=mat;refs=4;consumer=100%
  25. family=reuse;policy=mat;refs=1;consumer=1%
  26. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  27. family=reuse;policy=mat;refs=1;consumer=100%
  28. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  29. family=reuse;policy=notmat;refs=4;consumer=100%
  30. family=bag;stmt=cte_joined_outer_mat;threshold=20

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
| family=bag;stmt=cte_joined_outer;threshold=180000 | j05_cte_joined_outer | valid | `98fc155d2be7ecfa` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20 | j05_cte_joined_outer | valid | `c4caf3d38532ca81` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=2000 | j05_cte_joined_outer | valid | `2d4bed75b4e4b408` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | j05m_cte_joined_outer_mat | valid | `3db368ed9173304a` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | j05m_cte_joined_outer_mat | valid | `313dedf3296c551e` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | j05m_cte_joined_outer_mat | valid | `4050eef18b4e6683` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | j05n_cte_joined_outer_notmat | valid | `56d9a31ae3e57473` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | j05n_cte_joined_outer_notmat | valid | `df8f3d78d25ebfb8` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | j05n_cte_joined_outer_notmat | valid | `c771887c2a442a6d` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=join_where;threshold=180000 | j01_join_where | valid | `2ca2d4eddc04c02c` |  |
| family=bag;stmt=join_where;threshold=20 | j01_join_where | valid | `d39744c3b7b014f6` |  |
| family=bag;stmt=join_where;threshold=2000 | j01_join_where | valid | `034f15a04d2a997e` |  |
| family=reuse;policy=default;refs=16;consumer=1% | r_default_n16_c1 | valid | `883eabbf6932f4f1` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=16;consumer=100% | r_default_n16_c100 | valid | `6e6ed7dfb336a674` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=1;consumer=1% | r_default_n1_c1 | valid | `9c1b3b3128e73bff` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=1;consumer=100% | r_default_n1_c100 | valid | `e1dd4ba22ebf8119` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=4;consumer=1% | r_default_n4_c1 | valid | `29dcec76d2411783` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=100% | r_default_n4_c100 | valid | `24bdcc8b9c1bc9c9` | policy=default; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=1% | r_mat_n16_c1 | valid | `3a5cde1136d5063e` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=100% | r_mat_n16_c100 | valid | `590ac95534dfddcb` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=1% | r_mat_n1_c1 | valid | `1f743cecb2f1f0a6` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=100% | r_mat_n1_c100 | valid | `be107f6397d551b9` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=1% | r_mat_n4_c1 | valid | `25c46eca56b92a9a` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=100% | r_mat_n4_c100 | valid | `47bc68e8d34fbd94` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=notmat;refs=16;consumer=1% | r_notmat_n16_c1 | valid | `cf39891d205579ea` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `d035fa763726b085` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `3bf77e70fd286dd8` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `56f9f7faf6140aff` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `90ae6f638e7cd7b8` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `0654fbf3ddd5ccd3` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, DIST, DEBUG, FORMAT TEXT` supported on yugabyte
- one host, single client: no independent hardware, real-network or equal-total-budget comparison
- no concurrent writers: the isolated read contract only; a matched blended read/write phase is a recorded gap
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
