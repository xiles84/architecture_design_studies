# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 30 cells, 30 valid, engine `yugabyte`, topology `cluster3`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-topo-yb-cluster-bo2-20261004`
- run tag: `run/study-01/select-join-topo-yb-cluster-bo2-20261004`
- environment: `host-zenbook-ux5406sa`
- repo commit: `0da058e8bbdc36acfdeec4082f56a5b0c8068a8e`
- engine version: `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit`
- database engine image: `docker.io/yugabytedb/yugabyte:2025.2.6.0-b111` @ `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`
- benchmark client image: `localhost/charitybench-select:select-join-topo-yb-cluster-bo2-20261004` @ `79bc753ff7e57e102982a323939c43e68e9b3664887683d3a5d5bf072650f4c1`
- inputs digest: `a4978ab5ae3d94523dab787de0c1fe380ecf0358d308dcd137edddce97350f5e`
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
- setting `phase`: blended
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
| yb-n1 | 163 | 2 |
| yb-n2 | 163 | 1 |
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
| 20 (0.01%) | 62.0984 | missing | missing | missing | missing | missing | 69.3139 | 103.1969 | 64.2143 | missing |
| 2000 (1%) | 66.9318 | missing | missing | missing | missing | missing | 70.1515 | 140.0756 | 83.9865 | missing |
| 180000 (90%) | 146.4789 | missing | missing | missing | missing | missing | 165.4338 | 147.4668 | 181.3808 | missing |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | cte_joined_outer vs join_where | 62.0984 | 69.3139 | 1.1162 | false |
| 20 | cte_joined_outer_mat vs join_where | 62.0984 | 103.1969 | 1.6618 | false |
| 20 | cte_joined_outer_notmat vs join_where | 62.0984 | 64.2143 | 1.0341 | false |
| 2000 | cte_joined_outer vs join_where | 66.9318 | 70.1515 | 1.0481 | false |
| 2000 | cte_joined_outer_mat vs join_where | 66.9318 | 140.0756 | 2.0928 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 66.9318 | 83.9865 | 1.2548 | false |
| 180000 | cte_joined_outer vs join_where | 146.4789 | 165.4338 | 1.1294 | false |
| 180000 | cte_joined_outer_mat vs join_where | 146.4789 | 147.4668 | 1.0067 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 146.4789 | 181.3808 | 1.2383 | false |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_joined_outer_mat vs cte_joined_outer | 69.3139 | 103.1969 | 1.4888 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 69.3139 | 64.2143 | 0.9264 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 70.1515 | 140.0756 | 1.9968 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 70.1515 | 83.9865 | 1.1972 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 165.4338 | 147.4668 | 0.8914 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 165.4338 | 181.3808 | 1.0964 | false | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 66.4623 | 1 | `bf18f466c2f44f65` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 68.4372 | 1 | `5ed6bcadb523ddf9` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | missing | | | | |
| default | 2 | 100% | missing | | | | |
| default | 4 | 1% | 66.1701 | 4 | `694d3997efd839b8` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 69.8505 | 4 | `da82f2259c4e3d18` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | missing | | | | |
| default | 8 | 100% | missing | | | | |
| default | 16 | 1% | 67.4369 | 16 | `1ad836c656ce2509` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 64.7954 | 16 | `d640eb9df89ae894` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 67.8635 | 1 | `29e709eb0eb65c22` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 68.6306 | 1 | `bf716b38516ed624` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | missing | | | | |
| mat | 2 | 100% | missing | | | | |
| mat | 4 | 1% | 68.3145 | 4 | `fe549e3ad2b73a2b` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 64.8082 | 4 | `0e64784d6ed60137` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | missing | | | | |
| mat | 8 | 100% | missing | | | | |
| mat | 16 | 1% | 66.0063 | 16 | `dee86d09493dbd10` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 64.9486 | 16 | `aa20798ac3ba8a98` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 67.4081 | 1 | `e65ad02cb2532bcd` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 70.9761 | 1 | `1e4ca6aa895cf2b6` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | missing | | | | |
| notmat | 2 | 100% | missing | | | | |
| notmat | 4 | 1% | 209.0527 | 4 | `85d1e618bd0af6c9` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 184.0309 | 4 | `4004c917e48de42e` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | missing | | | | |
| notmat | 8 | 100% | missing | | | | |
| notmat | 16 | 1% | 756.3907 | 16 | `84f1c0e92ba72b5e` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 634.1322 | 16 | `eb5f1f5bd3a26042` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 1.0211 | false |
| notmat | 1 | 1% | 1.0142 | false |
| mat | 1 | 100% | 1.0028 | false |
| notmat | 1 | 100% | 1.0371 | false |
| mat | 4 | 1% | 1.0324 | false |
| notmat | 4 | 1% | 3.1593 | false |
| mat | 4 | 100% | 0.9278 | false |
| notmat | 4 | 100% | 2.6346 | false |
| mat | 16 | 1% | 0.9788 | false |
| notmat | 16 | 1% | 11.2163 | false |
| mat | 16 | 100% | 1.0024 | false |
| notmat | 16 | 100% | 9.7867 | false |

## Latency samples and derived throughput

_The reader is sequential, so `queries/s` is 1000/median_ms, a derived rate and not a concurrent-capacity test. min/max are the five sample extremes._

| axis | status | min ms | median ms | max ms | samples | derived queries/s |
|---|---|---|---|---|---|---|
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 148.5640 | 165.4338 | 224.9161 | 5 | 6.0 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 53.8370 | 69.3139 | 70.5417 | 5 | 14.4 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 58.5758 | 70.1515 | 86.5648 | 5 | 14.3 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 141.3504 | 147.4668 | 178.4269 | 5 | 6.8 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 82.0885 | 103.1969 | 106.7640 | 5 | 9.7 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 105.2707 | 140.0756 | 183.0653 | 5 | 7.1 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 170.8229 | 181.3808 | 215.8870 | 5 | 5.5 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 63.1893 | 64.2143 | 76.6184 | 5 | 15.6 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 59.5934 | 83.9865 | 97.8594 | 5 | 11.9 |
| family=bag;stmt=join_where;threshold=180000 | valid | 124.1691 | 146.4789 | 201.2320 | 5 | 6.8 |
| family=bag;stmt=join_where;threshold=20 | valid | 59.6306 | 62.0984 | 67.7456 | 5 | 16.1 |
| family=bag;stmt=join_where;threshold=2000 | valid | 61.5742 | 66.9318 | 278.4879 | 5 | 14.9 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 52.3746 | 67.4369 | 71.8330 | 5 | 14.8 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 63.2043 | 64.7954 | 68.7269 | 5 | 15.4 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 63.5642 | 66.4623 | 69.2560 | 5 | 15.0 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 63.8709 | 68.4372 | 70.0230 | 5 | 14.6 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 61.2350 | 66.1701 | 72.5009 | 5 | 15.1 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 52.3880 | 69.8505 | 72.4551 | 5 | 14.3 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 51.4445 | 66.0063 | 76.6926 | 5 | 15.2 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 50.1981 | 64.9486 | 66.2662 | 5 | 15.4 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 64.1144 | 67.8635 | 68.5055 | 5 | 14.7 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 52.8710 | 68.6306 | 73.8074 | 5 | 14.6 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 63.1836 | 68.3145 | 72.0413 | 5 | 14.6 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 63.7965 | 64.8082 | 70.9513 | 5 | 15.4 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 619.8683 | 756.3907 | 836.7400 | 5 | 1.3 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 490.9418 | 634.1322 | 655.7851 | 5 | 1.6 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 53.3259 | 67.4081 | 68.1065 | 5 | 14.8 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 68.0935 | 70.9761 | 74.3481 | 5 | 14.1 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 154.0399 | 209.0527 | 359.9639 | 5 | 4.8 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 134.2851 | 184.0309 | 251.6547 | 5 | 5.4 |

## Matched blended workload

- phase: `blended`; target offered rate: 10 updates/s
- offered 379, committed 379, failed 0 over 37.857 s
- offered rate 10.011/s, completed rate 10.011/s
- reader operations overlapping a committed write: 173 of 211
- reader retries of transient conflicts: 1
- per key: 379 keys, min 1, max 1 per key
- ledger: 379 keys reconciled, 0 mismatch
- update latency: avg 7.749 ms, min 2.796, max 38.116 over 379

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42492698`, 5 interleaved rounds
  1. family=reuse;policy=mat;refs=4;consumer=1%
  2. family=reuse;policy=notmat;refs=1;consumer=100%
  3. family=reuse;policy=default;refs=4;consumer=1%
  4. family=reuse;policy=mat;refs=1;consumer=1%
  5. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  6. family=reuse;policy=mat;refs=4;consumer=100%
  7. family=bag;stmt=join_where;threshold=20
  8. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  9. family=reuse;policy=default;refs=1;consumer=1%
  10. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  11. family=reuse;policy=default;refs=1;consumer=100%
  12. family=bag;stmt=join_where;threshold=180000
  13. family=reuse;policy=default;refs=16;consumer=100%
  14. family=bag;stmt=join_where;threshold=2000
  15. family=reuse;policy=notmat;refs=4;consumer=1%
  16. family=bag;stmt=cte_joined_outer;threshold=2000
  17. family=reuse;policy=notmat;refs=16;consumer=1%
  18. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  19. family=reuse;policy=notmat;refs=4;consumer=100%
  20. family=bag;stmt=cte_joined_outer;threshold=20
  21. family=reuse;policy=notmat;refs=1;consumer=1%
  22. family=reuse;policy=mat;refs=1;consumer=100%
  23. family=reuse;policy=default;refs=4;consumer=100%
  24. family=reuse;policy=mat;refs=16;consumer=1%
  25. family=bag;stmt=cte_joined_outer;threshold=180000
  26. family=reuse;policy=default;refs=16;consumer=1%
  27. family=reuse;policy=mat;refs=16;consumer=100%
  28. family=reuse;policy=notmat;refs=16;consumer=100%
  29. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
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
| family=bag;stmt=cte_joined_outer;threshold=180000 | j05_cte_joined_outer | valid | `8afab9ab3a752dd4` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20 | j05_cte_joined_outer | valid | `b158cda01019e180` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=2000 | j05_cte_joined_outer | valid | `b1884493de0e2a28` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | j05m_cte_joined_outer_mat | valid | `a187e2803f22600b` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | j05m_cte_joined_outer_mat | valid | `897ff3c51cc33b92` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | j05m_cte_joined_outer_mat | valid | `3eab71b6a6bc4efb` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | j05n_cte_joined_outer_notmat | valid | `a165b22aeb89c036` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | j05n_cte_joined_outer_notmat | valid | `d44d8906fb0598ed` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | j05n_cte_joined_outer_notmat | valid | `72c9dfcab408d629` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=join_where;threshold=180000 | j01_join_where | valid | `6afcc62e1ec770e5` |  |
| family=bag;stmt=join_where;threshold=20 | j01_join_where | valid | `85752500dd362e63` |  |
| family=bag;stmt=join_where;threshold=2000 | j01_join_where | valid | `ec58d5be517ed55b` |  |
| family=reuse;policy=default;refs=16;consumer=1% | r_default_n16_c1 | valid | `1ad836c656ce2509` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=16;consumer=100% | r_default_n16_c100 | valid | `d640eb9df89ae894` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=1;consumer=1% | r_default_n1_c1 | valid | `bf18f466c2f44f65` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=1;consumer=100% | r_default_n1_c100 | valid | `5ed6bcadb523ddf9` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=4;consumer=1% | r_default_n4_c1 | valid | `694d3997efd839b8` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=100% | r_default_n4_c100 | valid | `da82f2259c4e3d18` | policy=default; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=1% | r_mat_n16_c1 | valid | `dee86d09493dbd10` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=100% | r_mat_n16_c100 | valid | `aa20798ac3ba8a98` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=1% | r_mat_n1_c1 | valid | `29e709eb0eb65c22` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=100% | r_mat_n1_c100 | valid | `bf716b38516ed624` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=1% | r_mat_n4_c1 | valid | `fe549e3ad2b73a2b` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=100% | r_mat_n4_c100 | valid | `0e64784d6ed60137` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=notmat;refs=16;consumer=1% | r_notmat_n16_c1 | valid | `84f1c0e92ba72b5e` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `eb5f1f5bd3a26042` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `e65ad02cb2532bcd` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `1e4ca6aa895cf2b6` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `85d1e618bd0af6c9` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `4004c917e48de42e` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, DIST, DEBUG, FORMAT TEXT` supported on yugabyte
- one host, single client: no independent hardware, real-network or equal-total-budget comparison
- no concurrent writers: the isolated read contract only; a matched blended read/write phase is a recorded gap
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
