# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 30 cells, 30 valid, engine `postgres`, topology `single`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-blend-pg-blended-20261004`
- run tag: `run/study-01/select-join-blend-pg-blended-20261004`
- environment: `host-zenbook-ux5406sa`
- repo commit: `0656c5bb0e7b5e3ed379890657069ed207411f0d`
- engine version: `PostgreSQL 17.11 (Debian 17.11-1.pgdg13+2) on x86_64-pc-linux-gnu, compiled by gcc (Debian 14.2.0-19) 14.2.0, 64-bit`
- database engine image: `docker.io/library/postgres:17.11` @ `7296f210ae81031ec955dbad9a67a84fe958572a2153b8d0826a647522904dc1`
- benchmark client image: `localhost/charitybench-select:select-join-blend-pg-blended-20261004` @ `c2f822dc232766475ab08703ae4a0402982438fa77e4c9268f6c04f2452a6628`
- inputs digest: `b821e6497868a209f7858eebf7b505b5e9ae8468b1d138360c83bb738d212177`
- dataset fingerprint: `fcf26e08f377ec51`
- client CPU quota: `100000 100000`, memory: `2147483648`
- setting `client_workers`: 1 sequential reader; 0 writers
- setting `connection_pool`: 2 per endpoint = reader/setup + reserve
- setting `customers`: 2000
- setting `dataset`: study01-person-donation-charity
- setting `distributed_rpc`: N/A on PostgreSQL single-node
- setting `endpoint_mode`: single endpoint
- setting `endpoints`: pg-single
- setting `fact_count`: 200000
- setting `index_set`: child rank index and parent-key index, primary keys; no per-spelling index
- setting `phase`: blended
- setting `prepared_plan_policy`: pgx simple protocol; no prepared statements
- setting `scale`: cte-domain
- setting `selectivity_mapping`: rank <= threshold; rank is a deterministic shuffle of 1..fact_count
- setting `statement_budget_ms`: 120000
- setting `statistics`: ANALYZE after deterministic load
- setting `thresholds`: [20 2000 180000]
- setting `topology`: single
- setting `trials_per_cell`: 1
- setting `whole_run_count`: 1; no replicated winner in this file
- setting `writer_keys`: 1000
- setting `writer_target_rate_per_sec`: 10

### Endpoint distribution

A cluster arm must rotate its reader across nodes; a single seed endpoint would not be balanced cluster access.

| endpoint | operations | connections at end |
|---|---|---|
| pg-single | 367 | 2 |

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
| 20 (0.01%) | 1.0106 | missing | missing | missing | missing | missing | 0.8070 | 94.9517 | 0.9065 | missing |
| 2000 (1%) | 3.8871 | missing | missing | missing | missing | missing | 4.6823 | 106.5769 | 3.5468 | missing |
| 180000 (90%) | 107.4813 | missing | missing | missing | missing | missing | 145.8113 | 169.7241 | 113.1865 | missing |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | cte_joined_outer vs join_where | 1.0106 | 0.8070 | 0.7985 | true |
| 20 | cte_joined_outer_mat vs join_where | 1.0106 | 94.9517 | 93.9523 | false |
| 20 | cte_joined_outer_notmat vs join_where | 1.0106 | 0.9065 | 0.8970 | true |
| 2000 | cte_joined_outer vs join_where | 3.8871 | 4.6823 | 1.2046 | true |
| 2000 | cte_joined_outer_mat vs join_where | 3.8871 | 106.5769 | 27.4182 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 3.8871 | 3.5468 | 0.9125 | true |
| 180000 | cte_joined_outer vs join_where | 107.4813 | 145.8113 | 1.3566 | true |
| 180000 | cte_joined_outer_mat vs join_where | 107.4813 | 169.7241 | 1.5791 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 107.4813 | 113.1865 | 1.0531 | true |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_joined_outer_mat vs cte_joined_outer | 0.8070 | 94.9517 | 117.6639 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 0.8070 | 0.9065 | 1.1234 | true | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 4.6823 | 106.5769 | 22.7614 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 4.6823 | 3.5468 | 0.7575 | true | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 145.8113 | 169.7241 | 1.1640 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 145.8113 | 113.1865 | 0.7763 | true | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 0.9782 | 1 | `549a07108ab01ba9` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 4.7249 | 1 | `2ce8b85a7b60c88e` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | missing | | | | |
| default | 2 | 100% | missing | | | | |
| default | 4 | 1% | 5.5035 | 4 | `ccfe1af393deae98` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 5.4959 | 4 | `ccfe1af393deae98` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | missing | | | | |
| default | 8 | 100% | missing | | | | |
| default | 16 | 1% | 6.5119 | 16 | `b0c356323d30c425` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 7.4676 | 16 | `b0c356323d30c425` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 4.8477 | 1 | `eaa7d6058aaa5e03` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 4.5005 | 1 | `eaa7d6058aaa5e03` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | missing | | | | |
| mat | 2 | 100% | missing | | | | |
| mat | 4 | 1% | 4.9087 | 4 | `ccfe1af393deae98` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 5.5901 | 4 | `ccfe1af393deae98` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | missing | | | | |
| mat | 8 | 100% | missing | | | | |
| mat | 16 | 1% | 7.1201 | 16 | `b0c356323d30c425` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 8.5240 | 16 | `b0c356323d30c425` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 1.0430 | 1 | `549a07108ab01ba9` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 3.9426 | 1 | `2ce8b85a7b60c88e` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | missing | | | | |
| notmat | 2 | 100% | missing | | | | |
| notmat | 4 | 1% | 2.0376 | 4 | `9f850e53560baafb` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 20.0301 | 4 | `699a7259ed6ccf17` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | missing | | | | |
| notmat | 8 | 100% | missing | | | | |
| notmat | 16 | 1% | 7.5182 | 16 | `1c38b58d8cec08d3` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 35.4794 | 16 | `a87df3d6439cfeeb` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 4.9559 | false |
| notmat | 1 | 1% | 1.0663 | true |
| mat | 1 | 100% | 0.9525 | false |
| notmat | 1 | 100% | 0.8344 | true |
| mat | 4 | 1% | 0.8919 | true |
| notmat | 4 | 1% | 0.3702 | false |
| mat | 4 | 100% | 1.0171 | true |
| notmat | 4 | 100% | 3.6446 | false |
| mat | 16 | 1% | 1.0934 | true |
| notmat | 16 | 1% | 1.1545 | false |
| mat | 16 | 100% | 1.1415 | true |
| notmat | 16 | 100% | 4.7511 | false |

## Latency samples and derived throughput

_The reader is sequential, so `queries/s` is 1000/median_ms, a derived rate and not a concurrent-capacity test. min/max are the five sample extremes._

| axis | status | min ms | median ms | max ms | samples | derived queries/s |
|---|---|---|---|---|---|---|
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 145.8113 | 145.8113 | 145.8113 | 1 | 6.9 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 0.8070 | 0.8070 | 0.8070 | 1 | 1239.2 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 4.6823 | 4.6823 | 4.6823 | 1 | 213.6 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 169.7241 | 169.7241 | 169.7241 | 1 | 5.9 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 94.9517 | 94.9517 | 94.9517 | 1 | 10.5 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 106.5769 | 106.5769 | 106.5769 | 1 | 9.4 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 113.1865 | 113.1865 | 113.1865 | 1 | 8.8 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 0.9065 | 0.9065 | 0.9065 | 1 | 1103.1 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 3.5468 | 3.5468 | 3.5468 | 1 | 281.9 |
| family=bag;stmt=join_where;threshold=180000 | valid | 107.4813 | 107.4813 | 107.4813 | 1 | 9.3 |
| family=bag;stmt=join_where;threshold=20 | valid | 1.0106 | 1.0106 | 1.0106 | 1 | 989.5 |
| family=bag;stmt=join_where;threshold=2000 | valid | 3.8871 | 3.8871 | 3.8871 | 1 | 257.3 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 6.5119 | 6.5119 | 6.5119 | 1 | 153.6 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 7.4676 | 7.4676 | 7.4676 | 1 | 133.9 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 0.9782 | 0.9782 | 0.9782 | 1 | 1022.3 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 4.7249 | 4.7249 | 4.7249 | 1 | 211.6 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 5.5035 | 5.5035 | 5.5035 | 1 | 181.7 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 5.4959 | 5.4959 | 5.4959 | 1 | 182.0 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 7.1201 | 7.1201 | 7.1201 | 1 | 140.4 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 8.5240 | 8.5240 | 8.5240 | 1 | 117.3 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 4.8477 | 4.8477 | 4.8477 | 1 | 206.3 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 4.5005 | 4.5005 | 4.5005 | 1 | 222.2 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 4.9087 | 4.9087 | 4.9087 | 1 | 203.7 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 5.5901 | 5.5901 | 5.5901 | 1 | 178.9 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 7.5182 | 7.5182 | 7.5182 | 1 | 133.0 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 35.4794 | 35.4794 | 35.4794 | 1 | 28.2 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 1.0430 | 1.0430 | 1.0430 | 1 | 958.7 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 3.9426 | 3.9426 | 3.9426 | 1 | 253.6 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 2.0376 | 2.0376 | 2.0376 | 1 | 490.8 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 20.0301 | 20.0301 | 20.0301 | 1 | 49.9 |

## Matched blended workload

- phase: `blended`; target offered rate: 10 updates/s
- offered 66, committed 66, failed 0 over 6.599 s
- offered rate 10.002/s, completed rate 10.002/s
- reader operations overlapping a committed write: 23 of 90
- per key: 66 keys, min 1, max 1 per key
- ledger: 66 keys reconciled, 0 mismatch
- update latency: avg 6.009 ms, min 1.555, max 31.341 over 66

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42608475`, 1 interleaved rounds
  1. family=reuse;policy=default;refs=16;consumer=1%
  2. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  3. family=reuse;policy=default;refs=4;consumer=100%
  4. family=bag;stmt=join_where;threshold=180000
  5. family=bag;stmt=cte_joined_outer_mat;threshold=20
  6. family=bag;stmt=join_where;threshold=20
  7. family=bag;stmt=join_where;threshold=2000
  8. family=reuse;policy=mat;refs=1;consumer=1%
  9. family=reuse;policy=notmat;refs=4;consumer=100%
  10. family=reuse;policy=default;refs=1;consumer=1%
  11. family=reuse;policy=default;refs=16;consumer=100%
  12. family=reuse;policy=notmat;refs=16;consumer=1%
  13. family=reuse;policy=mat;refs=16;consumer=1%
  14. family=bag;stmt=cte_joined_outer;threshold=180000
  15. family=reuse;policy=default;refs=1;consumer=100%
  16. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  17. family=reuse;policy=notmat;refs=16;consumer=100%
  18. family=reuse;policy=notmat;refs=1;consumer=100%
  19. family=reuse;policy=mat;refs=4;consumer=1%
  20. family=bag;stmt=cte_joined_outer;threshold=20
  21. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  22. family=reuse;policy=mat;refs=1;consumer=100%
  23. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  24. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  25. family=reuse;policy=notmat;refs=1;consumer=1%
  26. family=reuse;policy=mat;refs=4;consumer=100%
  27. family=bag;stmt=cte_joined_outer;threshold=2000
  28. family=reuse;policy=mat;refs=16;consumer=100%
  29. family=reuse;policy=default;refs=4;consumer=1%
  30. family=reuse;policy=notmat;refs=4;consumer=1%

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
| family=bag;stmt=cte_joined_outer;threshold=180000 | j05_cte_joined_outer | valid | `18e11c5d2e685630` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20 | j05_cte_joined_outer | valid | `61819ba7d0f5f952` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=2000 | j05_cte_joined_outer | valid | `5573434208edb047` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | j05m_cte_joined_outer_mat | valid | `7f87bace4cf1da4c` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | j05m_cte_joined_outer_mat | valid | `7f87bace4cf1da4c` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | j05m_cte_joined_outer_mat | valid | `7f87bace4cf1da4c` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | j05n_cte_joined_outer_notmat | valid | `18e11c5d2e685630` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | j05n_cte_joined_outer_notmat | valid | `61819ba7d0f5f952` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | j05n_cte_joined_outer_notmat | valid | `5573434208edb047` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=join_where;threshold=180000 | j01_join_where | valid | `18e11c5d2e685630` |  |
| family=bag;stmt=join_where;threshold=20 | j01_join_where | valid | `61819ba7d0f5f952` |  |
| family=bag;stmt=join_where;threshold=2000 | j01_join_where | valid | `5573434208edb047` |  |
| family=reuse;policy=default;refs=16;consumer=1% | r_default_n16_c1 | valid | `b0c356323d30c425` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=16;consumer=100% | r_default_n16_c100 | valid | `b0c356323d30c425` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=1;consumer=1% | r_default_n1_c1 | valid | `549a07108ab01ba9` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=1;consumer=100% | r_default_n1_c100 | valid | `2ce8b85a7b60c88e` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=4;consumer=1% | r_default_n4_c1 | valid | `ccfe1af393deae98` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=100% | r_default_n4_c100 | valid | `ccfe1af393deae98` | policy=default; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=1% | r_mat_n16_c1 | valid | `b0c356323d30c425` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=100% | r_mat_n16_c100 | valid | `b0c356323d30c425` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=1% | r_mat_n1_c1 | valid | `eaa7d6058aaa5e03` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=100% | r_mat_n1_c100 | valid | `eaa7d6058aaa5e03` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=1% | r_mat_n4_c1 | valid | `ccfe1af393deae98` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=100% | r_mat_n4_c100 | valid | `ccfe1af393deae98` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=notmat;refs=16;consumer=1% | r_notmat_n16_c1 | valid | `1c38b58d8cec08d3` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `a87df3d6439cfeeb` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `549a07108ab01ba9` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `2ce8b85a7b60c88e` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `9f850e53560baafb` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `699a7259ed6ccf17` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, VERBOSE, BUFFERS, FORMAT TEXT` supported on postgres
- one host, single client: no independent hardware, real-network or equal-total-budget comparison
- no concurrent writers: the isolated read contract only; a matched blended read/write phase is a recorded gap
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
