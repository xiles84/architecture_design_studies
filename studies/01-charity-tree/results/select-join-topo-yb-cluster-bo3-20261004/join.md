# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 30 cells, 30 valid, engine `yugabyte`, topology `cluster3`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-topo-yb-cluster-bo3-20261004`
- run tag: `run/study-01/select-join-topo-yb-cluster-bo3-20261004`
- environment: `host-zenbook-ux5406sa`
- repo commit: `0da058e8bbdc36acfdeec4082f56a5b0c8068a8e`
- engine version: `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit`
- database engine image: `docker.io/yugabytedb/yugabyte:2025.2.6.0-b111` @ `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`
- benchmark client image: `localhost/charitybench-select:select-join-topo-yb-cluster-bo3-20261004` @ `79bc753ff7e57e102982a323939c43e68e9b3664887683d3a5d5bf072650f4c1`
- inputs digest: `d040c5a573f18ee2d8bd92137d2bee0ba0f36ea85717be975ee6b4e8784cd8ec`
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
| 20 (0.01%) | 59.7879 | missing | missing | missing | missing | missing | 60.6358 | 164.2476 | 64.7563 | missing |
| 2000 (1%) | 60.6704 | missing | missing | missing | missing | missing | 60.8721 | 133.6519 | 69.1991 | missing |
| 180000 (90%) | 181.1888 | missing | missing | missing | missing | missing | 164.4606 | 153.3465 | 130.1848 | missing |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | cte_joined_outer vs join_where | 59.7879 | 60.6358 | 1.0142 | false |
| 20 | cte_joined_outer_mat vs join_where | 59.7879 | 164.2476 | 2.7472 | false |
| 20 | cte_joined_outer_notmat vs join_where | 59.7879 | 64.7563 | 1.0831 | false |
| 2000 | cte_joined_outer vs join_where | 60.6704 | 60.8721 | 1.0033 | false |
| 2000 | cte_joined_outer_mat vs join_where | 60.6704 | 133.6519 | 2.2029 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 60.6704 | 69.1991 | 1.1406 | false |
| 180000 | cte_joined_outer vs join_where | 181.1888 | 164.4606 | 0.9077 | false |
| 180000 | cte_joined_outer_mat vs join_where | 181.1888 | 153.3465 | 0.8463 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 181.1888 | 130.1848 | 0.7185 | false |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_joined_outer_mat vs cte_joined_outer | 60.6358 | 164.2476 | 2.7088 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 60.6358 | 64.7563 | 1.0680 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 60.8721 | 133.6519 | 2.1956 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 60.8721 | 69.1991 | 1.1368 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 164.4606 | 153.3465 | 0.9324 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 164.4606 | 130.1848 | 0.7916 | false | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 59.2939 | 1 | `1cfe29374af0ce07` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 60.3320 | 1 | `4e49a6f1fb5f1650` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | missing | | | | |
| default | 2 | 100% | missing | | | | |
| default | 4 | 1% | 63.5566 | 4 | `3d8f19eb74b69d6e` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 62.7718 | 4 | `9c2da56bf9fd0c93` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | missing | | | | |
| default | 8 | 100% | missing | | | | |
| default | 16 | 1% | 57.9646 | 16 | `c1db7bbbe9fa1bee` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 60.7695 | 16 | `c58d59d93e4c0bca` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 63.2637 | 1 | `6e644ccb68174ee8` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 59.0598 | 1 | `869711974e988524` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | missing | | | | |
| mat | 2 | 100% | missing | | | | |
| mat | 4 | 1% | 57.2060 | 4 | `dc7eec524b8e534a` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 65.1011 | 4 | `1a5d6fa1d35a1ec5` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | missing | | | | |
| mat | 8 | 100% | missing | | | | |
| mat | 16 | 1% | 61.7703 | 16 | `ba4d83867bf84a87` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 60.4983 | 16 | `c24501b02bc690b0` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 59.8556 | 1 | `65e12e06ac02f1eb` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 58.3541 | 1 | `fc7efef577c3942a` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | missing | | | | |
| notmat | 2 | 100% | missing | | | | |
| notmat | 4 | 1% | 322.6704 | 4 | `488934f2763e4e99` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 148.9859 | 4 | `65f6a2c19dc432f2` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | missing | | | | |
| notmat | 8 | 100% | missing | | | | |
| notmat | 16 | 1% | 916.4125 | 16 | `2c0030907235bfb0` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 605.5264 | 16 | `b49a286371bdc7f0` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 1.0669 | false |
| notmat | 1 | 1% | 1.0095 | false |
| mat | 1 | 100% | 0.9789 | false |
| notmat | 1 | 100% | 0.9672 | false |
| mat | 4 | 1% | 0.9001 | false |
| notmat | 4 | 1% | 5.0769 | false |
| mat | 4 | 100% | 1.0371 | false |
| notmat | 4 | 100% | 2.3735 | false |
| mat | 16 | 1% | 1.0657 | false |
| notmat | 16 | 1% | 15.8099 | false |
| mat | 16 | 100% | 0.9955 | false |
| notmat | 16 | 100% | 9.9643 | false |

## Latency samples and derived throughput

_The reader is sequential, so `queries/s` is 1000/median_ms, a derived rate and not a concurrent-capacity test. min/max are the five sample extremes._

| axis | status | min ms | median ms | max ms | samples | derived queries/s |
|---|---|---|---|---|---|---|
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 150.0791 | 164.4606 | 373.2006 | 5 | 6.1 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 56.3621 | 60.6358 | 64.4859 | 5 | 16.5 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 56.3556 | 60.8721 | 61.7028 | 5 | 16.4 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 145.2837 | 153.3465 | 176.2612 | 5 | 6.5 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 129.0927 | 164.2476 | 176.4965 | 5 | 6.1 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 125.4421 | 133.6519 | 197.3478 | 5 | 7.5 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 120.7698 | 130.1848 | 142.6553 | 5 | 7.7 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 54.2278 | 64.7563 | 115.9671 | 5 | 15.4 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 57.1963 | 69.1991 | 146.6330 | 5 | 14.5 |
| family=bag;stmt=join_where;threshold=180000 | valid | 149.5564 | 181.1888 | 225.2351 | 5 | 5.5 |
| family=bag;stmt=join_where;threshold=20 | valid | 51.7096 | 59.7879 | 64.1789 | 5 | 16.7 |
| family=bag;stmt=join_where;threshold=2000 | valid | 54.9680 | 60.6704 | 66.1449 | 5 | 16.5 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 54.5279 | 57.9646 | 66.3869 | 5 | 17.3 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 57.4409 | 60.7695 | 69.2260 | 5 | 16.5 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 55.3972 | 59.2939 | 61.0414 | 5 | 16.9 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 54.5140 | 60.3320 | 63.9584 | 5 | 16.6 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 57.0301 | 63.5566 | 121.9153 | 5 | 15.7 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 56.2060 | 62.7718 | 101.9120 | 5 | 15.9 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 53.7519 | 61.7703 | 116.8378 | 5 | 16.2 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 57.2610 | 60.4983 | 65.7846 | 5 | 16.5 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 58.8184 | 63.2637 | 73.9538 | 5 | 15.8 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 54.0640 | 59.0598 | 68.9287 | 5 | 16.9 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 53.5013 | 57.2060 | 67.1449 | 5 | 17.5 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 55.9694 | 65.1011 | 72.1259 | 5 | 15.4 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 626.0325 | 916.4125 | 1070.8895 | 5 | 1.1 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 594.5214 | 605.5264 | 630.3152 | 5 | 1.7 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 52.8137 | 59.8556 | 63.3875 | 5 | 16.7 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 57.1034 | 58.3541 | 72.6643 | 5 | 17.1 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 265.5488 | 322.6704 | 327.2539 | 5 | 3.1 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 135.6939 | 148.9859 | 173.7332 | 5 | 6.7 |

## Matched blended workload

- phase: `blended`; target offered rate: 10 updates/s
- offered 393, committed 393, failed 0 over 39.218 s
- offered rate 10.021/s, completed rate 10.021/s
- reader operations overlapping a committed write: 156 of 210
- reader retries of transient conflicts: 0
- per key: 393 keys, min 1, max 1 per key
- ledger: 393 keys reconciled, 0 mismatch
- update latency: avg 5.679 ms, min 2.542, max 161.641 over 393

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42912691`, 5 interleaved rounds
  1. family=reuse;policy=default;refs=4;consumer=100%
  2. family=reuse;policy=mat;refs=16;consumer=1%
  3. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  4. family=bag;stmt=cte_joined_outer;threshold=180000
  5. family=reuse;policy=notmat;refs=16;consumer=1%
  6. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  7. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  8. family=reuse;policy=mat;refs=1;consumer=1%
  9. family=reuse;policy=mat;refs=4;consumer=100%
  10. family=bag;stmt=join_where;threshold=180000
  11. family=reuse;policy=mat;refs=16;consumer=100%
  12. family=reuse;policy=default;refs=16;consumer=1%
  13. family=reuse;policy=mat;refs=1;consumer=100%
  14. family=bag;stmt=join_where;threshold=20
  15. family=bag;stmt=join_where;threshold=2000
  16. family=reuse;policy=notmat;refs=1;consumer=100%
  17. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  18. family=reuse;policy=default;refs=16;consumer=100%
  19. family=bag;stmt=cte_joined_outer_mat;threshold=20
  20. family=reuse;policy=default;refs=4;consumer=1%
  21. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  22. family=reuse;policy=notmat;refs=4;consumer=100%
  23. family=bag;stmt=cte_joined_outer;threshold=20
  24. family=reuse;policy=notmat;refs=4;consumer=1%
  25. family=reuse;policy=mat;refs=4;consumer=1%
  26. family=reuse;policy=default;refs=1;consumer=100%
  27. family=reuse;policy=notmat;refs=1;consumer=1%
  28. family=reuse;policy=default;refs=1;consumer=1%
  29. family=bag;stmt=cte_joined_outer;threshold=2000
  30. family=reuse;policy=notmat;refs=16;consumer=100%

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
| family=bag;stmt=cte_joined_outer;threshold=180000 | j05_cte_joined_outer | valid | `2a88a8a99f278ba8` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20 | j05_cte_joined_outer | valid | `9dbd3e26cfce4146` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=2000 | j05_cte_joined_outer | valid | `e2fdbbd3d726c862` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | j05m_cte_joined_outer_mat | valid | `1c9d6b6051b3433f` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | j05m_cte_joined_outer_mat | valid | `d1f4b86f690e8ccd` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | j05m_cte_joined_outer_mat | valid | `2e3dce8cbd38235d` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | j05n_cte_joined_outer_notmat | valid | `1722325aa825bf47` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | j05n_cte_joined_outer_notmat | valid | `06e0ca7772fc4728` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | j05n_cte_joined_outer_notmat | valid | `b1bb5759319da139` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=join_where;threshold=180000 | j01_join_where | valid | `ead7c1cdd28927fb` |  |
| family=bag;stmt=join_where;threshold=20 | j01_join_where | valid | `3904603d993f82d2` |  |
| family=bag;stmt=join_where;threshold=2000 | j01_join_where | valid | `c97fb9fcfc588f62` |  |
| family=reuse;policy=default;refs=16;consumer=1% | r_default_n16_c1 | valid | `c1db7bbbe9fa1bee` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=16;consumer=100% | r_default_n16_c100 | valid | `c58d59d93e4c0bca` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=1;consumer=1% | r_default_n1_c1 | valid | `1cfe29374af0ce07` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=1;consumer=100% | r_default_n1_c100 | valid | `4e49a6f1fb5f1650` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=4;consumer=1% | r_default_n4_c1 | valid | `3d8f19eb74b69d6e` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=100% | r_default_n4_c100 | valid | `9c2da56bf9fd0c93` | policy=default; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=1% | r_mat_n16_c1 | valid | `ba4d83867bf84a87` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=100% | r_mat_n16_c100 | valid | `c24501b02bc690b0` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=1% | r_mat_n1_c1 | valid | `6e644ccb68174ee8` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=100% | r_mat_n1_c100 | valid | `869711974e988524` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=1% | r_mat_n4_c1 | valid | `dc7eec524b8e534a` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=100% | r_mat_n4_c100 | valid | `1a5d6fa1d35a1ec5` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=notmat;refs=16;consumer=1% | r_notmat_n16_c1 | valid | `2c0030907235bfb0` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `b49a286371bdc7f0` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `65e12e06ac02f1eb` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `fc7efef577c3942a` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `488934f2763e4e99` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `65f6a2c19dc432f2` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, DIST, DEBUG, FORMAT TEXT` supported on yugabyte
- one host, single client: no independent hardware, real-network or equal-total-budget comparison
- no concurrent writers: the isolated read contract only; a matched blended read/write phase is a recorded gap
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
