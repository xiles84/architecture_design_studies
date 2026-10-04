# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 30 cells, 30 valid, engine `postgres`, topology `single`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-blend-pg-isolated-20261004`
- run tag: `run/study-01/select-join-blend-pg-isolated-20261004`
- environment: `host-zenbook-ux5406sa`
- repo commit: `0656c5bb0e7b5e3ed379890657069ed207411f0d`
- engine version: `PostgreSQL 17.11 (Debian 17.11-1.pgdg13+2) on x86_64-pc-linux-gnu, compiled by gcc (Debian 14.2.0-19) 14.2.0, 64-bit`
- database engine image: `docker.io/library/postgres:17.11` @ `7296f210ae81031ec955dbad9a67a84fe958572a2153b8d0826a647522904dc1`
- benchmark client image: `localhost/charitybench-select:select-join-blend-pg-isolated-20261004` @ `c2f822dc232766475ab08703ae4a0402982438fa77e4c9268f6c04f2452a6628`
- inputs digest: `da57ecc6e4ec8de63d8d90542b60b31007da7506297a0bbe88dc9a94b9b2021a`
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
- setting `phase`: isolated
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
| pg-single | 367 | 1 |

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
| 20 (0.01%) | 0.9667 | missing | missing | missing | missing | missing | 1.0818 | 126.6456 | 1.0544 | missing |
| 2000 (1%) | 3.6713 | missing | missing | missing | missing | missing | 7.6706 | 90.8731 | 4.0379 | missing |
| 180000 (90%) | 123.5987 | missing | missing | missing | missing | missing | 125.4244 | 182.8792 | 108.8432 | missing |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | cte_joined_outer vs join_where | 0.9667 | 1.0818 | 1.1191 | true |
| 20 | cte_joined_outer_mat vs join_where | 0.9667 | 126.6456 | 131.0070 | false |
| 20 | cte_joined_outer_notmat vs join_where | 0.9667 | 1.0544 | 1.0907 | true |
| 2000 | cte_joined_outer vs join_where | 3.6713 | 7.6706 | 2.0893 | true |
| 2000 | cte_joined_outer_mat vs join_where | 3.6713 | 90.8731 | 24.7523 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 3.6713 | 4.0379 | 1.0998 | true |
| 180000 | cte_joined_outer vs join_where | 123.5987 | 125.4244 | 1.0148 | true |
| 180000 | cte_joined_outer_mat vs join_where | 123.5987 | 182.8792 | 1.4796 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 123.5987 | 108.8432 | 0.8806 | true |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_joined_outer_mat vs cte_joined_outer | 1.0818 | 126.6456 | 117.0663 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 1.0818 | 1.0544 | 0.9746 | true | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 7.6706 | 90.8731 | 11.8469 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 7.6706 | 4.0379 | 0.5264 | true | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 125.4244 | 182.8792 | 1.4581 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 125.4244 | 108.8432 | 0.8678 | true | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 1.0487 | 1 | `549a07108ab01ba9` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 9.3268 | 1 | `2ce8b85a7b60c88e` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | missing | | | | |
| default | 2 | 100% | missing | | | | |
| default | 4 | 1% | 14.5888 | 4 | `ccfe1af393deae98` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 3.4174 | 4 | `ccfe1af393deae98` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | missing | | | | |
| default | 8 | 100% | missing | | | | |
| default | 16 | 1% | 6.5988 | 16 | `b0c356323d30c425` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 9.5244 | 16 | `b0c356323d30c425` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 3.7331 | 1 | `eaa7d6058aaa5e03` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 5.6353 | 1 | `eaa7d6058aaa5e03` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | missing | | | | |
| mat | 2 | 100% | missing | | | | |
| mat | 4 | 1% | 3.3799 | 4 | `ccfe1af393deae98` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 6.6161 | 4 | `ccfe1af393deae98` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | missing | | | | |
| mat | 8 | 100% | missing | | | | |
| mat | 16 | 1% | 6.3531 | 16 | `b0c356323d30c425` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 8.2586 | 16 | `b0c356323d30c425` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 1.5052 | 1 | `549a07108ab01ba9` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 4.5689 | 1 | `2ce8b85a7b60c88e` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | missing | | | | |
| notmat | 2 | 100% | missing | | | | |
| notmat | 4 | 1% | 2.8690 | 4 | `9f850e53560baafb` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 22.1289 | 4 | `d556f90b36f5b191` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | missing | | | | |
| notmat | 8 | 100% | missing | | | | |
| notmat | 16 | 1% | 9.0988 | 16 | `1c38b58d8cec08d3` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 34.8603 | 16 | `f4d3bb967cdb4e35` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 3.5597 | false |
| notmat | 1 | 1% | 1.4353 | true |
| mat | 1 | 100% | 0.6042 | false |
| notmat | 1 | 100% | 0.4899 | true |
| mat | 4 | 1% | 0.2317 | true |
| notmat | 4 | 1% | 0.1967 | false |
| mat | 4 | 100% | 1.9360 | true |
| notmat | 4 | 100% | 6.4753 | false |
| mat | 16 | 1% | 0.9628 | true |
| notmat | 16 | 1% | 1.3789 | false |
| mat | 16 | 100% | 0.8671 | true |
| notmat | 16 | 100% | 3.6601 | false |

## Latency samples and derived throughput

_The reader is sequential, so `queries/s` is 1000/median_ms, a derived rate and not a concurrent-capacity test. min/max are the five sample extremes._

| axis | status | min ms | median ms | max ms | samples | derived queries/s |
|---|---|---|---|---|---|---|
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 125.4244 | 125.4244 | 125.4244 | 1 | 8.0 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 1.0818 | 1.0818 | 1.0818 | 1 | 924.4 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 7.6706 | 7.6706 | 7.6706 | 1 | 130.4 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 182.8792 | 182.8792 | 182.8792 | 1 | 5.5 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 126.6456 | 126.6456 | 126.6456 | 1 | 7.9 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 90.8731 | 90.8731 | 90.8731 | 1 | 11.0 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 108.8432 | 108.8432 | 108.8432 | 1 | 9.2 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 1.0544 | 1.0544 | 1.0544 | 1 | 948.4 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 4.0379 | 4.0379 | 4.0379 | 1 | 247.7 |
| family=bag;stmt=join_where;threshold=180000 | valid | 123.5987 | 123.5987 | 123.5987 | 1 | 8.1 |
| family=bag;stmt=join_where;threshold=20 | valid | 0.9667 | 0.9667 | 0.9667 | 1 | 1034.4 |
| family=bag;stmt=join_where;threshold=2000 | valid | 3.6713 | 3.6713 | 3.6713 | 1 | 272.4 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 6.5988 | 6.5988 | 6.5988 | 1 | 151.5 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 9.5244 | 9.5244 | 9.5244 | 1 | 105.0 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 1.0487 | 1.0487 | 1.0487 | 1 | 953.5 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 9.3268 | 9.3268 | 9.3268 | 1 | 107.2 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 14.5888 | 14.5888 | 14.5888 | 1 | 68.5 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 3.4174 | 3.4174 | 3.4174 | 1 | 292.6 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 6.3531 | 6.3531 | 6.3531 | 1 | 157.4 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 8.2586 | 8.2586 | 8.2586 | 1 | 121.1 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 3.7331 | 3.7331 | 3.7331 | 1 | 267.9 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 5.6353 | 5.6353 | 5.6353 | 1 | 177.5 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 3.3799 | 3.3799 | 3.3799 | 1 | 295.9 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 6.6161 | 6.6161 | 6.6161 | 1 | 151.1 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 9.0988 | 9.0988 | 9.0988 | 1 | 109.9 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 34.8603 | 34.8603 | 34.8603 | 1 | 28.7 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 1.5052 | 1.5052 | 1.5052 | 1 | 664.4 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 4.5689 | 4.5689 | 4.5689 | 1 | 218.9 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 2.8690 | 2.8690 | 2.8690 | 1 | 348.6 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 22.1289 | 22.1289 | 22.1289 | 1 | 45.2 |

## Matched blended workload

- phase: `isolated`; target offered rate: 10 updates/s
- offered 0, committed 0, failed 0 over 0.000 s
- offered rate 0.000/s, completed rate 0.000/s
- per key: no keys touched
- ledger: 0 keys reconciled, 0 mismatch

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42450882`, 1 interleaved rounds
  1. family=reuse;policy=default;refs=4;consumer=1%
  2. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  3. family=bag;stmt=cte_joined_outer;threshold=20
  4. family=bag;stmt=join_where;threshold=180000
  5. family=reuse;policy=mat;refs=16;consumer=1%
  6. family=reuse;policy=notmat;refs=16;consumer=100%
  7. family=reuse;policy=default;refs=16;consumer=1%
  8. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  9. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  10. family=reuse;policy=default;refs=1;consumer=100%
  11. family=reuse;policy=mat;refs=16;consumer=100%
  12. family=bag;stmt=cte_joined_outer;threshold=2000
  13. family=reuse;policy=notmat;refs=4;consumer=1%
  14. family=reuse;policy=notmat;refs=16;consumer=1%
  15. family=reuse;policy=notmat;refs=1;consumer=100%
  16. family=reuse;policy=default;refs=16;consumer=100%
  17. family=bag;stmt=join_where;threshold=2000
  18. family=bag;stmt=cte_joined_outer_mat;threshold=20
  19. family=reuse;policy=notmat;refs=4;consumer=100%
  20. family=reuse;policy=mat;refs=1;consumer=100%
  21. family=reuse;policy=mat;refs=4;consumer=100%
  22. family=reuse;policy=mat;refs=1;consumer=1%
  23. family=reuse;policy=default;refs=4;consumer=100%
  24. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  25. family=reuse;policy=default;refs=1;consumer=1%
  26. family=reuse;policy=mat;refs=4;consumer=1%
  27. family=bag;stmt=cte_joined_outer;threshold=180000
  28. family=reuse;policy=notmat;refs=1;consumer=1%
  29. family=bag;stmt=join_where;threshold=20
  30. family=bag;stmt=cte_joined_outer_notmat;threshold=2000

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
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `f4d3bb967cdb4e35` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `549a07108ab01ba9` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `2ce8b85a7b60c88e` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `9f850e53560baafb` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `d556f90b36f5b191` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, VERBOSE, BUFFERS, FORMAT TEXT` supported on postgres
- one host, single client: no independent hardware, real-network or equal-total-budget comparison
- no concurrent writers: the isolated read contract only; a matched blended read/write phase is a recorded gap
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
