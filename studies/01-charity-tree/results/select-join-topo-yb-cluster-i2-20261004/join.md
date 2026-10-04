# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 30 cells, 30 valid, engine `yugabyte`, topology `cluster3`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-topo-yb-cluster-i2-20261004`
- run tag: `run/study-01/select-join-topo-yb-cluster-i2-20261004`
- environment: `host-zenbook-ux5406sa`
- repo commit: `85d78f9e942acf856b48d416d0fa26de43a4af66`
- engine version: `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit`
- database engine image: `docker.io/yugabytedb/yugabyte:2025.2.6.0-b111` @ `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`
- benchmark client image: `localhost/charitybench-select:select-join-topo-yb-cluster-i2-20261004` @ `d8db2fd98a7bc00fa2a5561ed01373926c9d950a8d93f0ae9ba4d53be68d6e38`
- inputs digest: `c1c021644575ba9b05ac523ad5b300d257d273f27f0ef586b873b4046d1ca82d`
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
| 20 (0.01%) | 61.2836 | missing | missing | missing | missing | missing | 63.7881 | 92.8607 | 58.3723 | missing |
| 2000 (1%) | 58.8678 | missing | missing | missing | missing | missing | 58.2179 | 91.6260 | 62.4607 | missing |
| 180000 (90%) | 159.0656 | missing | missing | missing | missing | missing | 157.5145 | 116.2442 | 163.0366 | missing |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | cte_joined_outer vs join_where | 61.2836 | 63.7881 | 1.0409 | false |
| 20 | cte_joined_outer_mat vs join_where | 61.2836 | 92.8607 | 1.5153 | false |
| 20 | cte_joined_outer_notmat vs join_where | 61.2836 | 58.3723 | 0.9525 | false |
| 2000 | cte_joined_outer vs join_where | 58.8678 | 58.2179 | 0.9890 | false |
| 2000 | cte_joined_outer_mat vs join_where | 58.8678 | 91.6260 | 1.5565 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 58.8678 | 62.4607 | 1.0610 | false |
| 180000 | cte_joined_outer vs join_where | 159.0656 | 157.5145 | 0.9902 | false |
| 180000 | cte_joined_outer_mat vs join_where | 159.0656 | 116.2442 | 0.7308 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 159.0656 | 163.0366 | 1.0250 | false |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_joined_outer_mat vs cte_joined_outer | 63.7881 | 92.8607 | 1.4558 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 63.7881 | 58.3723 | 0.9151 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 58.2179 | 91.6260 | 1.5738 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 58.2179 | 62.4607 | 1.0729 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 157.5145 | 116.2442 | 0.7380 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 157.5145 | 163.0366 | 1.0351 | false | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 64.1190 | 1 | `952046204d72c8f5` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 60.6322 | 1 | `2f2b4e59ef7b2954` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | missing | | | | |
| default | 2 | 100% | missing | | | | |
| default | 4 | 1% | 61.0438 | 4 | `1645debb31238e9e` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 59.9284 | 4 | `ceff64c3ef6b3050` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | missing | | | | |
| default | 8 | 100% | missing | | | | |
| default | 16 | 1% | 61.0842 | 16 | `b811fe5594504cf7` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 66.6020 | 16 | `c8f19f925c28ee30` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 61.6748 | 1 | `267de0ef63af580d` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 60.1288 | 1 | `2af47f0aa9be0245` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | missing | | | | |
| mat | 2 | 100% | missing | | | | |
| mat | 4 | 1% | 59.8558 | 4 | `3ff18f5a0012f577` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 61.4999 | 4 | `f9ef0bb71e9ff8f4` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | missing | | | | |
| mat | 8 | 100% | missing | | | | |
| mat | 16 | 1% | 60.6726 | 16 | `bc3df43c43539a0b` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 64.7898 | 16 | `a02fd336841ae218` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 59.4324 | 1 | `c2de57a1b99a6f0a` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 62.8535 | 1 | `c75422c7f8f00495` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | missing | | | | |
| notmat | 2 | 100% | missing | | | | |
| notmat | 4 | 1% | 140.4409 | 4 | `c2098e0edac97b26` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 177.7283 | 4 | `0d0bee58a700bba5` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | missing | | | | |
| notmat | 8 | 100% | missing | | | | |
| notmat | 16 | 1% | 607.9929 | 16 | `15e91002e24a33c4` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 615.0791 | 16 | `c286d553a54198f1` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 0.9619 | false |
| notmat | 1 | 1% | 0.9269 | false |
| mat | 1 | 100% | 0.9917 | false |
| notmat | 1 | 100% | 1.0366 | false |
| mat | 4 | 1% | 0.9805 | false |
| notmat | 4 | 1% | 2.3007 | false |
| mat | 4 | 100% | 1.0262 | false |
| notmat | 4 | 100% | 2.9657 | false |
| mat | 16 | 1% | 0.9933 | false |
| notmat | 16 | 1% | 9.9534 | false |
| mat | 16 | 100% | 0.9728 | false |
| notmat | 16 | 100% | 9.2351 | false |

## Latency samples and derived throughput

_The reader is sequential, so `queries/s` is 1000/median_ms, a derived rate and not a concurrent-capacity test. min/max are the five sample extremes._

| axis | status | min ms | median ms | max ms | samples | derived queries/s |
|---|---|---|---|---|---|---|
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 153.0076 | 157.5145 | 207.9381 | 5 | 6.3 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 58.2100 | 63.7881 | 68.5649 | 5 | 15.7 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 56.5559 | 58.2179 | 65.6350 | 5 | 17.2 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 111.4909 | 116.2442 | 140.7393 | 5 | 8.6 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 86.5773 | 92.8607 | 95.7319 | 5 | 10.8 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 89.3522 | 91.6260 | 100.8012 | 5 | 10.9 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 145.9287 | 163.0366 | 181.3205 | 5 | 6.1 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 55.9220 | 58.3723 | 70.4047 | 5 | 17.1 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 59.3346 | 62.4607 | 66.5124 | 5 | 16.0 |
| family=bag;stmt=join_where;threshold=180000 | valid | 150.5681 | 159.0656 | 192.4554 | 5 | 6.3 |
| family=bag;stmt=join_where;threshold=20 | valid | 55.2601 | 61.2836 | 62.4931 | 5 | 16.3 |
| family=bag;stmt=join_where;threshold=2000 | valid | 58.6470 | 58.8678 | 68.7080 | 5 | 17.0 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 59.2569 | 61.0842 | 65.5080 | 5 | 16.4 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 59.9697 | 66.6020 | 81.4181 | 5 | 15.0 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 54.6696 | 64.1190 | 69.4348 | 5 | 15.6 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 58.0314 | 60.6322 | 67.3457 | 5 | 16.5 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 60.2604 | 61.0438 | 66.5347 | 5 | 16.4 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 57.4200 | 59.9284 | 66.1619 | 5 | 16.7 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 58.6695 | 60.6726 | 64.8796 | 5 | 16.5 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 60.8387 | 64.7898 | 67.5717 | 5 | 15.4 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 55.4056 | 61.6748 | 63.9786 | 5 | 16.2 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 57.4572 | 60.1288 | 69.8902 | 5 | 16.6 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 57.8721 | 59.8558 | 66.4628 | 5 | 16.7 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 60.3499 | 61.4999 | 63.7524 | 5 | 16.3 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 572.2283 | 607.9929 | 695.0231 | 5 | 1.6 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 547.6257 | 615.0791 | 670.8115 | 5 | 1.6 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 58.5090 | 59.4324 | 66.7004 | 5 | 16.8 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 58.7723 | 62.8535 | 65.4446 | 5 | 15.9 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 133.3123 | 140.4409 | 169.3184 | 5 | 7.1 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 137.9112 | 177.7283 | 216.1151 | 5 | 5.6 |

## Matched blended workload

- phase: `isolated`; target offered rate: 10 updates/s
- offered 0, committed 0, failed 0 over 0.000 s
- offered rate 0.000/s, completed rate 0.000/s
- per key: no keys touched
- ledger: 0 keys reconciled, 0 mismatch

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42275112`, 5 interleaved rounds
  1. family=reuse;policy=notmat;refs=16;consumer=1%
  2. family=bag;stmt=cte_joined_outer;threshold=20
  3. family=reuse;policy=notmat;refs=16;consumer=100%
  4. family=reuse;policy=default;refs=16;consumer=100%
  5. family=reuse;policy=mat;refs=1;consumer=1%
  6. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  7. family=reuse;policy=mat;refs=16;consumer=100%
  8. family=reuse;policy=default;refs=16;consumer=1%
  9. family=reuse;policy=mat;refs=16;consumer=1%
  10. family=reuse;policy=mat;refs=4;consumer=1%
  11. family=bag;stmt=join_where;threshold=2000
  12. family=bag;stmt=cte_joined_outer;threshold=180000
  13. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  14. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  15. family=reuse;policy=notmat;refs=4;consumer=1%
  16. family=reuse;policy=default;refs=1;consumer=100%
  17. family=reuse;policy=notmat;refs=1;consumer=100%
  18. family=bag;stmt=join_where;threshold=20
  19. family=reuse;policy=default;refs=4;consumer=100%
  20. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  21. family=reuse;policy=default;refs=4;consumer=1%
  22. family=reuse;policy=notmat;refs=4;consumer=100%
  23. family=reuse;policy=notmat;refs=1;consumer=1%
  24. family=reuse;policy=mat;refs=1;consumer=100%
  25. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  26. family=bag;stmt=cte_joined_outer_mat;threshold=20
  27. family=bag;stmt=join_where;threshold=180000
  28. family=bag;stmt=cte_joined_outer;threshold=2000
  29. family=reuse;policy=mat;refs=4;consumer=100%
  30. family=reuse;policy=default;refs=1;consumer=1%

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
| family=bag;stmt=cte_joined_outer;threshold=180000 | j05_cte_joined_outer | valid | `fde8d47376799e72` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20 | j05_cte_joined_outer | valid | `c0fa8a68264891b0` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=2000 | j05_cte_joined_outer | valid | `02a50abc46d6fa6b` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | j05m_cte_joined_outer_mat | valid | `9cb26bcb91c760a5` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | j05m_cte_joined_outer_mat | valid | `c07be71a760540c3` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | j05m_cte_joined_outer_mat | valid | `443c06b824f4dddf` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | j05n_cte_joined_outer_notmat | valid | `fb772b6566a87221` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | j05n_cte_joined_outer_notmat | valid | `22c263a382b49e66` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | j05n_cte_joined_outer_notmat | valid | `bb34f108f2692515` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=join_where;threshold=180000 | j01_join_where | valid | `120d9420a4c0dbf8` |  |
| family=bag;stmt=join_where;threshold=20 | j01_join_where | valid | `8b7bbb0f9c877108` |  |
| family=bag;stmt=join_where;threshold=2000 | j01_join_where | valid | `df154847f4f3fd1c` |  |
| family=reuse;policy=default;refs=16;consumer=1% | r_default_n16_c1 | valid | `b811fe5594504cf7` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=16;consumer=100% | r_default_n16_c100 | valid | `c8f19f925c28ee30` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=1;consumer=1% | r_default_n1_c1 | valid | `952046204d72c8f5` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=1;consumer=100% | r_default_n1_c100 | valid | `2f2b4e59ef7b2954` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=4;consumer=1% | r_default_n4_c1 | valid | `1645debb31238e9e` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=100% | r_default_n4_c100 | valid | `ceff64c3ef6b3050` | policy=default; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=1% | r_mat_n16_c1 | valid | `bc3df43c43539a0b` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=100% | r_mat_n16_c100 | valid | `a02fd336841ae218` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=1% | r_mat_n1_c1 | valid | `267de0ef63af580d` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=100% | r_mat_n1_c100 | valid | `2af47f0aa9be0245` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=1% | r_mat_n4_c1 | valid | `3ff18f5a0012f577` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=100% | r_mat_n4_c100 | valid | `f9ef0bb71e9ff8f4` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=notmat;refs=16;consumer=1% | r_notmat_n16_c1 | valid | `15e91002e24a33c4` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `c286d553a54198f1` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `c2de57a1b99a6f0a` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `c75422c7f8f00495` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `c2098e0edac97b26` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `0d0bee58a700bba5` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, DIST, DEBUG, FORMAT TEXT` supported on yugabyte
- one host, single client: no independent hardware, real-network or equal-total-budget comparison
- no concurrent writers: the isolated read contract only; a matched blended read/write phase is a recorded gap
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
