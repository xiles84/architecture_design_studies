# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 30 cells, 30 valid, engine `yugabyte`, topology `single`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-topo-yb-single-i3-20261004`
- run tag: `run/study-01/select-join-topo-yb-single-i3-20261004`
- environment: `host-zenbook-ux5406sa`
- repo commit: `656489ae7d67583b0c3ba6dbb0808051e9a9cb32`
- engine version: `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit`
- database engine image: `docker.io/yugabytedb/yugabyte:2025.2.6.0-b111` @ `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`
- benchmark client image: `localhost/charitybench-select:select-join-topo-yb-single-i3-20261004` @ `d8db2fd98a7bc00fa2a5561ed01373926c9d950a8d93f0ae9ba4d53be68d6e38`
- inputs digest: `fad2afa49ce1096389e7ce1d40ec65e6bf187d688a9ca6bd78f9db316963ab80`
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
| 20 (0.01%) | 61.8880 | missing | missing | missing | missing | missing | 57.4356 | 107.0941 | 60.1146 | missing |
| 2000 (1%) | 64.8960 | missing | missing | missing | missing | missing | 65.6797 | 110.8418 | 60.5303 | missing |
| 180000 (90%) | 123.4313 | missing | missing | missing | missing | missing | 148.6556 | 126.1569 | 137.0404 | missing |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | cte_joined_outer vs join_where | 61.8880 | 57.4356 | 0.9281 | false |
| 20 | cte_joined_outer_mat vs join_where | 61.8880 | 107.0941 | 1.7304 | false |
| 20 | cte_joined_outer_notmat vs join_where | 61.8880 | 60.1146 | 0.9713 | false |
| 2000 | cte_joined_outer vs join_where | 64.8960 | 65.6797 | 1.0121 | false |
| 2000 | cte_joined_outer_mat vs join_where | 64.8960 | 110.8418 | 1.7080 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 64.8960 | 60.5303 | 0.9327 | false |
| 180000 | cte_joined_outer vs join_where | 123.4313 | 148.6556 | 1.2044 | false |
| 180000 | cte_joined_outer_mat vs join_where | 123.4313 | 126.1569 | 1.0221 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 123.4313 | 137.0404 | 1.1103 | false |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_joined_outer_mat vs cte_joined_outer | 57.4356 | 107.0941 | 1.8646 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 57.4356 | 60.1146 | 1.0466 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 65.6797 | 110.8418 | 1.6876 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 65.6797 | 60.5303 | 0.9216 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 148.6556 | 126.1569 | 0.8487 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 148.6556 | 137.0404 | 0.9219 | false | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 57.1655 | 1 | `cfddc4fdf367c745` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 68.7158 | 1 | `63ef17234ca12d72` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | missing | | | | |
| default | 2 | 100% | missing | | | | |
| default | 4 | 1% | 65.3764 | 4 | `287b39f7e71480f1` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 63.0049 | 4 | `896a6295e6c0ced9` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | missing | | | | |
| default | 8 | 100% | missing | | | | |
| default | 16 | 1% | 64.7629 | 16 | `89cc905aa1ede197` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 59.1468 | 16 | `ad79a29c13b9b3d1` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 57.9389 | 1 | `dafd8984a6ee52b5` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 57.2066 | 1 | `1f822aac1ee73b85` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | missing | | | | |
| mat | 2 | 100% | missing | | | | |
| mat | 4 | 1% | 69.4263 | 4 | `f310e87b0c7a9cd7` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 65.6005 | 4 | `b505bed047347dc6` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | missing | | | | |
| mat | 8 | 100% | missing | | | | |
| mat | 16 | 1% | 58.0999 | 16 | `0c21c40e2f92ddf7` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 62.8693 | 16 | `59f1dafbae4a390f` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 57.4687 | 1 | `748be67031c91f63` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 63.1438 | 1 | `de27309be516dabe` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | missing | | | | |
| notmat | 2 | 100% | missing | | | | |
| notmat | 4 | 1% | 199.8259 | 4 | `ee98ff923bb79cc6` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 181.3559 | 4 | `321b321745ed3e2e` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | missing | | | | |
| notmat | 8 | 100% | missing | | | | |
| notmat | 16 | 1% | 626.3086 | 16 | `8afa441e92d66c2a` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 697.1195 | 16 | `c2871f1548ff5489` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 1.0135 | false |
| notmat | 1 | 1% | 1.0053 | false |
| mat | 1 | 100% | 0.8325 | false |
| notmat | 1 | 100% | 0.9189 | false |
| mat | 4 | 1% | 1.0619 | false |
| notmat | 4 | 1% | 3.0565 | false |
| mat | 4 | 100% | 1.0412 | false |
| notmat | 4 | 100% | 2.8784 | false |
| mat | 16 | 1% | 0.8971 | false |
| notmat | 16 | 1% | 9.6708 | false |
| mat | 16 | 100% | 1.0629 | false |
| notmat | 16 | 100% | 11.7863 | false |

## Latency samples and derived throughput

_The reader is sequential, so `queries/s` is 1000/median_ms, a derived rate and not a concurrent-capacity test. min/max are the five sample extremes._

| axis | status | min ms | median ms | max ms | samples | derived queries/s |
|---|---|---|---|---|---|---|
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 119.6624 | 148.6556 | 187.7621 | 5 | 6.7 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 51.4643 | 57.4356 | 98.0744 | 5 | 17.4 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 59.1105 | 65.6797 | 73.3591 | 5 | 15.2 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 118.9640 | 126.1569 | 225.9655 | 5 | 7.9 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 92.5198 | 107.0941 | 152.7282 | 5 | 9.3 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 85.8156 | 110.8418 | 122.0299 | 5 | 9.0 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 123.4201 | 137.0404 | 147.1668 | 5 | 7.3 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 56.2091 | 60.1146 | 72.9730 | 5 | 16.6 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 54.6196 | 60.5303 | 71.2311 | 5 | 16.5 |
| family=bag;stmt=join_where;threshold=180000 | valid | 120.7982 | 123.4313 | 162.6114 | 5 | 8.1 |
| family=bag;stmt=join_where;threshold=20 | valid | 51.3504 | 61.8880 | 66.7669 | 5 | 16.2 |
| family=bag;stmt=join_where;threshold=2000 | valid | 63.3608 | 64.8960 | 111.5980 | 5 | 15.4 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 57.1554 | 64.7629 | 67.8215 | 5 | 15.4 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 54.4093 | 59.1468 | 70.2112 | 5 | 16.9 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 52.8891 | 57.1655 | 140.7778 | 5 | 17.5 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 60.5385 | 68.7158 | 73.3574 | 5 | 14.6 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 61.3933 | 65.3764 | 68.6177 | 5 | 15.3 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 58.2415 | 63.0049 | 75.1505 | 5 | 15.9 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 56.6788 | 58.0999 | 75.4982 | 5 | 17.2 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 57.0854 | 62.8693 | 70.8142 | 5 | 15.9 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 52.2473 | 57.9389 | 71.6584 | 5 | 17.3 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 50.4100 | 57.2066 | 73.0810 | 5 | 17.5 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 57.6837 | 69.4263 | 92.5838 | 5 | 14.4 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 52.8416 | 65.6005 | 121.4161 | 5 | 15.2 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 520.9087 | 626.3086 | 770.6471 | 5 | 1.6 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 580.5476 | 697.1195 | 824.4113 | 5 | 1.4 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 55.4288 | 57.4687 | 80.0188 | 5 | 17.4 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 56.9724 | 63.1438 | 77.5746 | 5 | 15.8 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 181.9959 | 199.8259 | 278.3412 | 5 | 5.0 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 170.9610 | 181.3559 | 231.0723 | 5 | 5.5 |

## Matched blended workload

- phase: `isolated`; target offered rate: 10 updates/s
- offered 0, committed 0, failed 0 over 0.000 s
- offered rate 0.000/s, completed rate 0.000/s
- per key: no keys touched
- ledger: 0 keys reconciled, 0 mismatch

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42368371`, 5 interleaved rounds
  1. family=reuse;policy=default;refs=4;consumer=1%
  2. family=reuse;policy=mat;refs=16;consumer=100%
  3. family=reuse;policy=default;refs=16;consumer=1%
  4. family=reuse;policy=default;refs=1;consumer=100%
  5. family=reuse;policy=default;refs=16;consumer=100%
  6. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  7. family=bag;stmt=join_where;threshold=20
  8. family=reuse;policy=notmat;refs=4;consumer=100%
  9. family=bag;stmt=cte_joined_outer;threshold=2000
  10. family=bag;stmt=cte_joined_outer_mat;threshold=20
  11. family=bag;stmt=join_where;threshold=2000
  12. family=bag;stmt=cte_joined_outer;threshold=180000
  13. family=reuse;policy=notmat;refs=16;consumer=1%
  14. family=reuse;policy=mat;refs=4;consumer=1%
  15. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  16. family=reuse;policy=default;refs=4;consumer=100%
  17. family=reuse;policy=mat;refs=1;consumer=1%
  18. family=bag;stmt=join_where;threshold=180000
  19. family=bag;stmt=cte_joined_outer;threshold=20
  20. family=reuse;policy=mat;refs=4;consumer=100%
  21. family=reuse;policy=default;refs=1;consumer=1%
  22. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  23. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  24. family=reuse;policy=mat;refs=16;consumer=1%
  25. family=reuse;policy=notmat;refs=4;consumer=1%
  26. family=reuse;policy=notmat;refs=1;consumer=1%
  27. family=reuse;policy=mat;refs=1;consumer=100%
  28. family=reuse;policy=notmat;refs=1;consumer=100%
  29. family=bag;stmt=cte_joined_outer_mat;threshold=2000
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
| family=bag;stmt=cte_joined_outer;threshold=180000 | j05_cte_joined_outer | valid | `139be6c822909f6d` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20 | j05_cte_joined_outer | valid | `b4d9236317c15e95` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=2000 | j05_cte_joined_outer | valid | `7df9d220c0c49988` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | j05m_cte_joined_outer_mat | valid | `fe78346e31a3ebd3` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | j05m_cte_joined_outer_mat | valid | `a72ea6726cad70d4` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | j05m_cte_joined_outer_mat | valid | `8e24f2a46f93d128` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | j05n_cte_joined_outer_notmat | valid | `19c6f452160e6d95` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | j05n_cte_joined_outer_notmat | valid | `95f0e01b35efe809` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | j05n_cte_joined_outer_notmat | valid | `671c762da73bf662` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=join_where;threshold=180000 | j01_join_where | valid | `77bf17dcd77e517e` |  |
| family=bag;stmt=join_where;threshold=20 | j01_join_where | valid | `dda91d40172a2eee` |  |
| family=bag;stmt=join_where;threshold=2000 | j01_join_where | valid | `e8c1389b97bbb180` |  |
| family=reuse;policy=default;refs=16;consumer=1% | r_default_n16_c1 | valid | `89cc905aa1ede197` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=16;consumer=100% | r_default_n16_c100 | valid | `ad79a29c13b9b3d1` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=1;consumer=1% | r_default_n1_c1 | valid | `cfddc4fdf367c745` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=1;consumer=100% | r_default_n1_c100 | valid | `63ef17234ca12d72` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=4;consumer=1% | r_default_n4_c1 | valid | `287b39f7e71480f1` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=100% | r_default_n4_c100 | valid | `896a6295e6c0ced9` | policy=default; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=1% | r_mat_n16_c1 | valid | `0c21c40e2f92ddf7` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=100% | r_mat_n16_c100 | valid | `59f1dafbae4a390f` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=1% | r_mat_n1_c1 | valid | `dafd8984a6ee52b5` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=100% | r_mat_n1_c100 | valid | `1f822aac1ee73b85` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=1% | r_mat_n4_c1 | valid | `f310e87b0c7a9cd7` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=100% | r_mat_n4_c100 | valid | `b505bed047347dc6` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=notmat;refs=16;consumer=1% | r_notmat_n16_c1 | valid | `8afa441e92d66c2a` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `c2871f1548ff5489` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `748be67031c91f63` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `de27309be516dabe` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `ee98ff923bb79cc6` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `321b321745ed3e2e` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, DIST, DEBUG, FORMAT TEXT` supported on yugabyte
- one host, single client: no independent hardware, real-network or equal-total-budget comparison
- no concurrent writers: the isolated read contract only; a matched blended read/write phase is a recorded gap
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
