# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 30 cells, 30 valid, engine `yugabyte`, topology `cluster3`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-topo-yb-cluster-i3-20261004`
- run tag: `run/study-01/select-join-topo-yb-cluster-i3-20261004`
- environment: `host-zenbook-ux5406sa`
- repo commit: `85d78f9e942acf856b48d416d0fa26de43a4af66`
- engine version: `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit`
- database engine image: `docker.io/yugabytedb/yugabyte:2025.2.6.0-b111` @ `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`
- benchmark client image: `localhost/charitybench-select:select-join-topo-yb-cluster-i3-20261004` @ `d8db2fd98a7bc00fa2a5561ed01373926c9d950a8d93f0ae9ba4d53be68d6e38`
- inputs digest: `670261f08e72281ed45e5fd18736d5eeb44cde05725dd4df5386767cdc3f6f23`
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
| 20 (0.01%) | 59.7423 | missing | missing | missing | missing | missing | 59.0641 | 137.8927 | 62.3660 | missing |
| 2000 (1%) | 71.1768 | missing | missing | missing | missing | missing | 62.8824 | 95.8888 | 61.8067 | missing |
| 180000 (90%) | 161.7484 | missing | missing | missing | missing | missing | 121.5300 | 132.0531 | 131.8830 | missing |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | cte_joined_outer vs join_where | 59.7423 | 59.0641 | 0.9886 | false |
| 20 | cte_joined_outer_mat vs join_where | 59.7423 | 137.8927 | 2.3081 | false |
| 20 | cte_joined_outer_notmat vs join_where | 59.7423 | 62.3660 | 1.0439 | false |
| 2000 | cte_joined_outer vs join_where | 71.1768 | 62.8824 | 0.8835 | false |
| 2000 | cte_joined_outer_mat vs join_where | 71.1768 | 95.8888 | 1.3472 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 71.1768 | 61.8067 | 0.8684 | false |
| 180000 | cte_joined_outer vs join_where | 161.7484 | 121.5300 | 0.7514 | false |
| 180000 | cte_joined_outer_mat vs join_where | 161.7484 | 132.0531 | 0.8164 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 161.7484 | 131.8830 | 0.8154 | false |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_joined_outer_mat vs cte_joined_outer | 59.0641 | 137.8927 | 2.3346 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 59.0641 | 62.3660 | 1.0559 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 62.8824 | 95.8888 | 1.5249 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 62.8824 | 61.8067 | 0.9829 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 121.5300 | 132.0531 | 1.0866 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 121.5300 | 131.8830 | 1.0852 | false | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 57.7788 | 1 | `ab021c1814eea664` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 63.0079 | 1 | `9c475139dcdc7a68` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | missing | | | | |
| default | 2 | 100% | missing | | | | |
| default | 4 | 1% | 58.2497 | 4 | `667c3a3262fa1745` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 63.2386 | 4 | `458f80f26d744ff3` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | missing | | | | |
| default | 8 | 100% | missing | | | | |
| default | 16 | 1% | 66.4819 | 16 | `606d8777bbd0ae0d` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 68.7383 | 16 | `79903263da17cc96` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 63.6636 | 1 | `1809f3860156895d` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 60.0949 | 1 | `f37b585a9dd46907` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | missing | | | | |
| mat | 2 | 100% | missing | | | | |
| mat | 4 | 1% | 61.2062 | 4 | `2cec5432b0250518` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 62.3945 | 4 | `de304b8e4f28daa4` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | missing | | | | |
| mat | 8 | 100% | missing | | | | |
| mat | 16 | 1% | 64.8867 | 16 | `aa24ed29655e27d5` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 58.7467 | 16 | `7c34505371366a9f` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 70.3594 | 1 | `4c1740712879f300` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 64.7306 | 1 | `b2f3e0d0367ee07d` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | missing | | | | |
| notmat | 2 | 100% | missing | | | | |
| notmat | 4 | 1% | 188.3161 | 4 | `7dc02dc1a7c9e41b` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 145.7134 | 4 | `1fe3fab0ac7dc97c` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | missing | | | | |
| notmat | 8 | 100% | missing | | | | |
| notmat | 16 | 1% | 598.5338 | 16 | `880e5a926c201ab9` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 593.5546 | 16 | `bc16de316de8de09` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 1.1019 | false |
| notmat | 1 | 1% | 1.2177 | false |
| mat | 1 | 100% | 0.9538 | false |
| notmat | 1 | 100% | 1.0273 | false |
| mat | 4 | 1% | 1.0508 | false |
| notmat | 4 | 1% | 3.2329 | false |
| mat | 4 | 100% | 0.9867 | false |
| notmat | 4 | 100% | 2.3042 | false |
| mat | 16 | 1% | 0.9760 | false |
| notmat | 16 | 1% | 9.0030 | false |
| mat | 16 | 100% | 0.8546 | false |
| notmat | 16 | 100% | 8.6350 | false |

## Latency samples and derived throughput

_The reader is sequential, so `queries/s` is 1000/median_ms, a derived rate and not a concurrent-capacity test. min/max are the five sample extremes._

| axis | status | min ms | median ms | max ms | samples | derived queries/s |
|---|---|---|---|---|---|---|
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 120.0605 | 121.5300 | 184.0101 | 5 | 8.2 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 53.1006 | 59.0641 | 69.2108 | 5 | 16.9 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 52.8905 | 62.8824 | 72.2339 | 5 | 15.9 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 124.2524 | 132.0531 | 145.5032 | 5 | 7.6 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 124.1440 | 137.8927 | 206.1857 | 5 | 7.3 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 91.9232 | 95.8888 | 122.8999 | 5 | 10.4 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 120.5954 | 131.8830 | 286.6628 | 5 | 7.6 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 56.3760 | 62.3660 | 66.6201 | 5 | 16.0 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 57.5836 | 61.8067 | 66.9382 | 5 | 16.2 |
| family=bag;stmt=join_where;threshold=180000 | valid | 145.4563 | 161.7484 | 175.9341 | 5 | 6.2 |
| family=bag;stmt=join_where;threshold=20 | valid | 56.5290 | 59.7423 | 62.9998 | 5 | 16.7 |
| family=bag;stmt=join_where;threshold=2000 | valid | 54.7198 | 71.1768 | 90.3598 | 5 | 14.0 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 56.2860 | 66.4819 | 76.4254 | 5 | 15.0 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 53.4254 | 68.7383 | 74.7320 | 5 | 14.5 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 49.4657 | 57.7788 | 68.4395 | 5 | 17.3 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 55.8639 | 63.0079 | 63.6284 | 5 | 15.9 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 48.9533 | 58.2497 | 64.2493 | 5 | 17.2 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 57.8944 | 63.2386 | 76.9176 | 5 | 15.8 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 55.6829 | 64.8867 | 72.5699 | 5 | 15.4 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 56.6199 | 58.7467 | 73.5244 | 5 | 17.0 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 60.2896 | 63.6636 | 66.8985 | 5 | 15.7 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 56.3142 | 60.0949 | 71.2015 | 5 | 16.6 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 54.2145 | 61.2062 | 71.5546 | 5 | 16.3 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 59.2737 | 62.3945 | 68.9847 | 5 | 16.0 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 573.2427 | 598.5338 | 684.2067 | 5 | 1.7 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 568.2068 | 593.5546 | 656.9210 | 5 | 1.7 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 59.2013 | 70.3594 | 72.5509 | 5 | 14.2 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 57.1105 | 64.7306 | 70.6341 | 5 | 15.4 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 157.3170 | 188.3161 | 202.4960 | 5 | 5.3 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 135.3623 | 145.7134 | 179.0060 | 5 | 6.9 |

## Matched blended workload

- phase: `isolated`; target offered rate: 10 updates/s
- offered 0, committed 0, failed 0 over 0.000 s
- offered rate 0.000/s, completed rate 0.000/s
- per key: no keys touched
- ledger: 0 keys reconciled, 0 mismatch

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42759185`, 5 interleaved rounds
  1. family=reuse;policy=default;refs=1;consumer=1%
  2. family=reuse;policy=notmat;refs=4;consumer=1%
  3. family=reuse;policy=default;refs=4;consumer=1%
  4. family=reuse;policy=default;refs=4;consumer=100%
  5. family=reuse;policy=notmat;refs=1;consumer=100%
  6. family=bag;stmt=cte_joined_outer;threshold=20
  7. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  8. family=reuse;policy=default;refs=1;consumer=100%
  9. family=reuse;policy=notmat;refs=4;consumer=100%
  10. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  11. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  12. family=reuse;policy=default;refs=16;consumer=100%
  13. family=bag;stmt=cte_joined_outer_mat;threshold=20
  14. family=reuse;policy=mat;refs=1;consumer=1%
  15. family=bag;stmt=join_where;threshold=20
  16. family=bag;stmt=join_where;threshold=180000
  17. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  18. family=reuse;policy=notmat;refs=16;consumer=100%
  19. family=reuse;policy=default;refs=16;consumer=1%
  20. family=reuse;policy=mat;refs=16;consumer=100%
  21. family=reuse;policy=notmat;refs=16;consumer=1%
  22. family=bag;stmt=join_where;threshold=2000
  23. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  24. family=reuse;policy=mat;refs=4;consumer=100%
  25. family=reuse;policy=mat;refs=1;consumer=100%
  26. family=reuse;policy=mat;refs=16;consumer=1%
  27. family=bag;stmt=cte_joined_outer;threshold=2000
  28. family=reuse;policy=notmat;refs=1;consumer=1%
  29. family=bag;stmt=cte_joined_outer;threshold=180000
  30. family=reuse;policy=mat;refs=4;consumer=1%

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
| family=bag;stmt=cte_joined_outer;threshold=180000 | j05_cte_joined_outer | valid | `9146c6f05151b017` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20 | j05_cte_joined_outer | valid | `5e36e5b771c9ab04` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=2000 | j05_cte_joined_outer | valid | `f580da6fd26bf3e3` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | j05m_cte_joined_outer_mat | valid | `9f682dae93f9a683` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | j05m_cte_joined_outer_mat | valid | `f3d3be3053854a66` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | j05m_cte_joined_outer_mat | valid | `ff77036567f3492f` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | j05n_cte_joined_outer_notmat | valid | `2693f335838b0b9f` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | j05n_cte_joined_outer_notmat | valid | `970eb94b09f016d0` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | j05n_cte_joined_outer_notmat | valid | `7ce7238a076d6c00` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=join_where;threshold=180000 | j01_join_where | valid | `a357c0e488bbd362` |  |
| family=bag;stmt=join_where;threshold=20 | j01_join_where | valid | `e9d47b7932de2f6e` |  |
| family=bag;stmt=join_where;threshold=2000 | j01_join_where | valid | `a6ff797f54218758` |  |
| family=reuse;policy=default;refs=16;consumer=1% | r_default_n16_c1 | valid | `606d8777bbd0ae0d` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=16;consumer=100% | r_default_n16_c100 | valid | `79903263da17cc96` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=1;consumer=1% | r_default_n1_c1 | valid | `ab021c1814eea664` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=1;consumer=100% | r_default_n1_c100 | valid | `9c475139dcdc7a68` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=4;consumer=1% | r_default_n4_c1 | valid | `667c3a3262fa1745` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=100% | r_default_n4_c100 | valid | `458f80f26d744ff3` | policy=default; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=1% | r_mat_n16_c1 | valid | `aa24ed29655e27d5` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=100% | r_mat_n16_c100 | valid | `7c34505371366a9f` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=1% | r_mat_n1_c1 | valid | `1809f3860156895d` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=100% | r_mat_n1_c100 | valid | `f37b585a9dd46907` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=1% | r_mat_n4_c1 | valid | `2cec5432b0250518` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=100% | r_mat_n4_c100 | valid | `de304b8e4f28daa4` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=notmat;refs=16;consumer=1% | r_notmat_n16_c1 | valid | `880e5a926c201ab9` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `bc16de316de8de09` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `4c1740712879f300` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `b2f3e0d0367ee07d` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `7dc02dc1a7c9e41b` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `1fe3fab0ac7dc97c` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, DIST, DEBUG, FORMAT TEXT` supported on yugabyte
- one host, single client: no independent hardware, real-network or equal-total-budget comparison
- no concurrent writers: the isolated read contract only; a matched blended read/write phase is a recorded gap
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
