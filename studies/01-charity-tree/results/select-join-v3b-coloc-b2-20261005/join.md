# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 30 cells, 30 valid, engine `yugabyte`, topology `cluster3`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-v3b-coloc-b2-20261005`
- run tag: `run/study-01/select-join-v3b-coloc-b2-20261005`
- environment: `host-zenbook-ux5406sa`
- repo commit: `7243fee074bfb16fc5ad8b6afc1a2d7549080935`
- engine version: `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit`
- database engine image: `docker.io/yugabytedb/yugabyte:2025.2.6.0-b111` @ `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`
- benchmark client image: `localhost/charitybench-select:select-join-v3b-coloc-b2-20261005` @ `f3f1e10c44cbcb4c7b5bc56ed95a4f071682b6f7d33b5e0a7b688c4271eeef88`
- inputs digest: `0581e4a967498f54a9448d5aabf854e944ba0c2020a8b1ab5153d3aa21b02256`
- dataset fingerprint: `fcf26e08f377ec51`
- client CPU quota: `100000 100000`, memory: `2147483648`
- setting `client_budget`: 1 CPU 2g
- setting `client_workers`: 1 sequential reader; 0 writers
- setting `connection_pool`: 2 per endpoint = reader/setup + reserve
- setting `customers`: 2000
- setting `dataset`: study01-person-donation-charity
- setting `distributed_rpc`: DIST/DEBUG storage counters parsed per cell (Storage Rows Scanned, Storage Index Rows Scanned); RPC-level counters not aggregated
- setting `endpoint_mode`: round-robin across 3 endpoints
- setting `endpoints`: yb-n1,yb-n2,yb-n3
- setting `fact_count`: 200000
- setting `index_arm`: yb-range-rank
- setting `index_set`: child rank index and parent-key index, primary keys; no per-spelling index
- setting `node_budget`: 2 CPU 3g
- setting `phase`: blended
- setting `placement`: colocated
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
| yb-n1 | 168 | 2 |
| yb-n2 | 168 | 2 |
| yb-n3 | 168 | 2 |

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
| 20 (0.01%) | 2.4009 | missing | missing | missing | missing | missing | 2.7007 | 78.8698 | 2.5863 | missing |
| 2000 (1%) | 8.0471 | missing | missing | missing | missing | missing | 7.2531 | 126.6278 | 7.2220 | missing |
| 180000 (90%) | 149.1025 | missing | missing | missing | missing | missing | 144.9383 | 147.7394 | 154.6765 | missing |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | cte_joined_outer vs join_where | 2.4009 | 2.7007 | 1.1249 | false |
| 20 | cte_joined_outer_mat vs join_where | 2.4009 | 78.8698 | 32.8499 | false |
| 20 | cte_joined_outer_notmat vs join_where | 2.4009 | 2.5863 | 1.0772 | false |
| 2000 | cte_joined_outer vs join_where | 8.0471 | 7.2531 | 0.9013 | false |
| 2000 | cte_joined_outer_mat vs join_where | 8.0471 | 126.6278 | 15.7359 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 8.0471 | 7.2220 | 0.8975 | false |
| 180000 | cte_joined_outer vs join_where | 149.1025 | 144.9383 | 0.9721 | false |
| 180000 | cte_joined_outer_mat vs join_where | 149.1025 | 147.7394 | 0.9909 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 149.1025 | 154.6765 | 1.0374 | false |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_joined_outer_mat vs cte_joined_outer | 2.7007 | 78.8698 | 29.2035 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 2.7007 | 2.5863 | 0.9576 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 7.2531 | 126.6278 | 17.4585 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 7.2531 | 7.2220 | 0.9957 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 144.9383 | 147.7394 | 1.0193 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 144.9383 | 154.6765 | 1.0672 | false | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 2.0189 | 1 | `f33eba260e73e3e7` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 6.7427 | 1 | `47b26c94f5788a95` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | missing | | | | |
| default | 2 | 100% | missing | | | | |
| default | 4 | 1% | 7.8086 | 4 | `398b1753006fcdd9` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 6.8956 | 4 | `a1b229931121f3a9` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | missing | | | | |
| default | 8 | 100% | missing | | | | |
| default | 16 | 1% | 8.7312 | 16 | `ea89b99b7aa0b0ed` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 8.2241 | 16 | `4637e66906d7bdc2` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 6.4261 | 1 | `5c3c50fde0c2ecf5` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 6.4124 | 1 | `b207778182468121` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | missing | | | | |
| mat | 2 | 100% | missing | | | | |
| mat | 4 | 1% | 8.1476 | 4 | `64ab28e00dc29470` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 8.2863 | 4 | `ccf583256c6b8417` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | missing | | | | |
| mat | 8 | 100% | missing | | | | |
| mat | 16 | 1% | 8.1164 | 16 | `7e9427cb2ab7343d` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 9.2145 | 16 | `ac30f7a5623a7891` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 1.8658 | 1 | `6082cddbe2aa5941` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 7.8148 | 1 | `96872545720613b3` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | missing | | | | |
| notmat | 2 | 100% | missing | | | | |
| notmat | 4 | 1% | 4.9077 | 4 | `b9a8dc94ea97e3c7` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 27.0590 | 4 | `67ca833d54f17631` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | missing | | | | |
| notmat | 8 | 100% | missing | | | | |
| notmat | 16 | 1% | 16.9673 | 16 | `575a74664debf8de` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 96.9797 | 16 | `016172c88794412a` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 3.1830 | false |
| notmat | 1 | 1% | 0.9242 | false |
| mat | 1 | 100% | 0.9510 | false |
| notmat | 1 | 100% | 1.1590 | false |
| mat | 4 | 1% | 1.0434 | false |
| notmat | 4 | 1% | 0.6285 | false |
| mat | 4 | 100% | 1.2017 | false |
| notmat | 4 | 100% | 3.9241 | false |
| mat | 16 | 1% | 0.9296 | false |
| notmat | 16 | 1% | 1.9433 | false |
| mat | 16 | 100% | 1.1204 | false |
| notmat | 16 | 100% | 11.7921 | false |

## Physical arm

- `charity_pkey`: `CREATE UNIQUE INDEX charity_pkey ON selectjoin_domain.charity USING lsm (charity_id ASC)`
- `donation_person_idx`: `CREATE INDEX donation_person_idx ON selectjoin_domain.donation USING lsm (person_id ASC)`
- `donation_pkey`: `CREATE UNIQUE INDEX donation_pkey ON selectjoin_domain.donation USING lsm (donation_id ASC)`
- `donation_rank_idx`: `CREATE INDEX donation_rank_idx ON selectjoin_domain.donation USING lsm (rank ASC)`
- `person_pkey`: `CREATE UNIQUE INDEX person_pkey ON selectjoin_domain.person USING lsm (person_id ASC)`

- index arm `yb-range-rank`; placement `colocated`; node budget `2 CPU 3g`; client budget `1 CPU 2g`

| table | tablets | hash key columns | colocated | tablegroup | colocation id |
|---|---|---|---|---|---|
| selectjoin_domain.charity | 1 | 0 | true | 16388 | 3898818602 |
| selectjoin_domain.person | 1 | 0 | true | 16388 | 1670309505 |
| selectjoin_domain.donation | 1 | 0 | true | 16388 | 760246173 |
| selectjoin_domain.donation_pkey | 0 | 0 | false | unknown | unknown |
| selectjoin_domain.donation_rank_idx | 1 | 0 | true | 16388 | 2670493207 |
| selectjoin_domain.donation_person_idx | 1 | 0 | true | 16388 | 1023403763 |

Physical-arm gate (default join must use the rank index at its three most selective thresholds):

- threshold 20: used
- threshold 2000: used

## Throttling (cgroup `cpu.stat` deltas)

- client: nr_periods 294, nr_throttled 90, throttled_usec 4295606 | quota 100000 100000
, memory 2147483648

- yb-n1: nr_periods 550, nr_throttled 5, throttled_usec 572002 | quota 200000 100000, memory 3221225472
- yb-n2: nr_periods 544, nr_throttled 12, throttled_usec 1377765 | quota 200000 100000, memory 3221225472
- yb-n3: nr_periods 538, nr_throttled 0, throttled_usec 0 | quota 200000 100000, memory 3221225472

## Access path per cell

_Engine-reported from each cell's own plan: whether the rank index or a sequential scan was used, rows scanned, PostgreSQL buffers and filter removals, and reader retries of transient conflicts. A cell that scanned the whole table at a selective threshold is a different physical regime, not a faster query._

| axis | status | index scan | seq scan | storage rows scanned | index rows scanned | buffers hit/read | rows removed | retries |
|---|---|---|---|---|---|---|---|---|
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | true | true | 202050 | 0 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | true | false | 110 | 20 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | true | false | 6050 | 2000 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | true | true | 202050 | 0 | 0/0 | 20000 | 0 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | true | true | 202050 | 0 | 0/0 | 199980 | 0 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | true | true | 202050 | 0 | 0/0 | 198000 | 0 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | true | true | 202050 | 0 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | true | false | 110 | 20 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | true | false | 6050 | 2000 | 0/0 | 0 | 0 |
| family=bag;stmt=join_where;threshold=180000 | valid | true | true | 202050 | 0 | 0/0 | 0 | 0 |
| family=bag;stmt=join_where;threshold=20 | valid | true | false | 110 | 20 | 0/0 | 0 | 0 |
| family=bag;stmt=join_where;threshold=2000 | valid | true | false | 6050 | 2000 | 0/0 | 0 | 0 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | true | false | 4000 | 2000 | 0/0 | 1983 | 0 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | true | false | 4000 | 2000 | 0/0 | 0 | 0 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | true | false | 2000 | 0 | 0/0 | 0 | 0 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | true | false | 4000 | 2000 | 0/0 | 0 | 0 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | true | false | 4000 | 2000 | 0/0 | 1983 | 0 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | true | false | 4000 | 2000 | 0/0 | 0 | 0 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | true | false | 4000 | 2000 | 0/0 | 1983 | 0 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | true | false | 4000 | 2000 | 0/0 | 0 | 0 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | true | false | 4000 | 2000 | 0/0 | 1983 | 0 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | true | false | 4000 | 2000 | 0/0 | 0 | 0 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | true | false | 4000 | 2000 | 0/0 | 1983 | 0 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | true | false | 4000 | 2000 | 0/0 | 0 | 0 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | true | false | 32000 | 0 | 0/0 | 0 | 0 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | true | false | 64000 | 2000 | 0/0 | 0 | 0 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | true | false | 2000 | 0 | 0/0 | 0 | 0 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | true | false | 4000 | 2000 | 0/0 | 0 | 0 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | true | false | 8000 | 0 | 0/0 | 0 | 0 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | true | false | 16000 | 2000 | 0/0 | 0 | 0 |

## Latency samples and derived throughput

_The reader is sequential, so `queries/s` is 1000/median_ms, a derived rate and not a concurrent-capacity test. min/max are the five sample extremes. Each sample includes any retry and backoff time the statement needed after a transient conflict; per-cell retry counts appear in the access-path table._

| axis | status | min ms | median ms | max ms | samples | derived queries/s |
|---|---|---|---|---|---|---|
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 138.6325 | 144.9383 | 156.6671 | 5 | 6.9 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 2.5393 | 2.7007 | 8.1741 | 5 | 370.3 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 7.0426 | 7.2531 | 9.4135 | 5 | 137.9 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 135.7102 | 147.7394 | 151.9061 | 5 | 6.8 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 77.3542 | 78.8698 | 87.4769 | 5 | 12.7 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 116.6782 | 126.6278 | 165.3192 | 5 | 7.9 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 146.1842 | 154.6765 | 163.7896 | 5 | 6.5 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 2.4468 | 2.5863 | 4.2310 | 5 | 386.7 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 6.6633 | 7.2220 | 9.4960 | 5 | 138.5 |
| family=bag;stmt=join_where;threshold=180000 | valid | 142.0205 | 149.1025 | 157.9821 | 5 | 6.7 |
| family=bag;stmt=join_where;threshold=20 | valid | 2.3303 | 2.4009 | 5.1257 | 5 | 416.5 |
| family=bag;stmt=join_where;threshold=2000 | valid | 7.5370 | 8.0471 | 8.6569 | 5 | 124.3 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 8.5315 | 8.7312 | 9.4071 | 5 | 114.5 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 8.0695 | 8.2241 | 10.4333 | 5 | 121.6 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 1.7140 | 2.0189 | 2.4393 | 5 | 495.3 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 6.5301 | 6.7427 | 6.9698 | 5 | 148.3 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 6.7842 | 7.8086 | 8.8420 | 5 | 128.1 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 6.5846 | 6.8956 | 9.5064 | 5 | 145.0 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 7.5116 | 8.1164 | 11.5244 | 5 | 123.2 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 8.7157 | 9.2145 | 11.7875 | 5 | 108.5 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 5.9385 | 6.4261 | 7.2024 | 5 | 155.6 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 6.2609 | 6.4124 | 9.6826 | 5 | 155.9 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 7.6131 | 8.1476 | 9.1332 | 5 | 122.7 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 7.5333 | 8.2863 | 11.5238 | 5 | 120.7 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 16.0994 | 16.9673 | 24.4081 | 5 | 58.9 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 85.3633 | 96.9797 | 203.5793 | 5 | 10.3 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 1.6589 | 1.8658 | 3.7623 | 5 | 536.0 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 7.4317 | 7.8148 | 8.4344 | 5 | 128.0 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 4.6631 | 4.9077 | 6.4188 | 5 | 203.8 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 24.2368 | 27.0590 | 33.1131 | 5 | 37.0 |

## Matched blended workload

- phase: `blended`; target offered rate: 10 updates/s
- offered 140, committed 140, failed 0 over 13.956 s
- offered rate 10.031/s, completed rate 10.031/s
- writer operations per endpoint: [47 47 46]
- reader operations overlapping a committed write: 66 of 210
- reader retries of transient conflicts: 0
- per key: 140 keys, min 1, max 1 per key
- ledger: 140 keys reconciled, 0 mismatch
- update latency: avg 5.822 ms, min 2.522, max 93.535 over 140

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42844991`, 5 interleaved rounds
  1. family=reuse;policy=notmat;refs=4;consumer=100%
  2. family=reuse;policy=default;refs=4;consumer=100%
  3. family=reuse;policy=notmat;refs=16;consumer=1%
  4. family=reuse;policy=default;refs=1;consumer=1%
  5. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  6. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  7. family=reuse;policy=mat;refs=16;consumer=100%
  8. family=reuse;policy=mat;refs=1;consumer=100%
  9. family=reuse;policy=notmat;refs=4;consumer=1%
  10. family=bag;stmt=join_where;threshold=20
  11. family=reuse;policy=default;refs=16;consumer=100%
  12. family=reuse;policy=default;refs=4;consumer=1%
  13. family=reuse;policy=notmat;refs=16;consumer=100%
  14. family=reuse;policy=notmat;refs=1;consumer=1%
  15. family=reuse;policy=default;refs=16;consumer=1%
  16. family=bag;stmt=join_where;threshold=2000
  17. family=reuse;policy=mat;refs=16;consumer=1%
  18. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  19. family=bag;stmt=cte_joined_outer;threshold=180000
  20. family=bag;stmt=cte_joined_outer;threshold=2000
  21. family=bag;stmt=join_where;threshold=180000
  22. family=reuse;policy=mat;refs=4;consumer=1%
  23. family=reuse;policy=default;refs=1;consumer=100%
  24. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  25. family=reuse;policy=notmat;refs=1;consumer=100%
  26. family=reuse;policy=mat;refs=1;consumer=1%
  27. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  28. family=bag;stmt=cte_joined_outer;threshold=20
  29. family=bag;stmt=cte_joined_outer_mat;threshold=20
  30. family=reuse;policy=mat;refs=4;consumer=100%

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
| family=bag;stmt=cte_joined_outer;threshold=180000 | j05_cte_joined_outer | valid | `98921f8c6891904a` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20 | j05_cte_joined_outer | valid | `0fefec5de6888f93` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=2000 | j05_cte_joined_outer | valid | `89ae13100e87336e` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | j05m_cte_joined_outer_mat | valid | `3bc1eba064b1053d` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | j05m_cte_joined_outer_mat | valid | `b18a7f0e183c34ba` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | j05m_cte_joined_outer_mat | valid | `721d1bf770f428f4` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | j05n_cte_joined_outer_notmat | valid | `6c167287b9d39f21` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | j05n_cte_joined_outer_notmat | valid | `e427431df521ed4c` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | j05n_cte_joined_outer_notmat | valid | `28b17fb272e53708` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=join_where;threshold=180000 | j01_join_where | valid | `5db48b506f9d3c23` |  |
| family=bag;stmt=join_where;threshold=20 | j01_join_where | valid | `055bcc90834483e3` |  |
| family=bag;stmt=join_where;threshold=2000 | j01_join_where | valid | `ab8f2d10d1ef8d92` |  |
| family=reuse;policy=default;refs=16;consumer=1% | r_default_n16_c1 | valid | `ea89b99b7aa0b0ed` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=16;consumer=100% | r_default_n16_c100 | valid | `4637e66906d7bdc2` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=1;consumer=1% | r_default_n1_c1 | valid | `f33eba260e73e3e7` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=1;consumer=100% | r_default_n1_c100 | valid | `47b26c94f5788a95` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=4;consumer=1% | r_default_n4_c1 | valid | `398b1753006fcdd9` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=100% | r_default_n4_c100 | valid | `a1b229931121f3a9` | policy=default; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=1% | r_mat_n16_c1 | valid | `7e9427cb2ab7343d` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=100% | r_mat_n16_c100 | valid | `ac30f7a5623a7891` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=1% | r_mat_n1_c1 | valid | `5c3c50fde0c2ecf5` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=100% | r_mat_n1_c100 | valid | `b207778182468121` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=1% | r_mat_n4_c1 | valid | `64ab28e00dc29470` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=100% | r_mat_n4_c100 | valid | `ccf583256c6b8417` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=notmat;refs=16;consumer=1% | r_notmat_n16_c1 | valid | `575a74664debf8de` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `016172c88794412a` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `6082cddbe2aa5941` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `96872545720613b3` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `b9a8dc94ea97e3c7` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `67ca833d54f17631` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, DIST, DEBUG, FORMAT TEXT` supported on yugabyte
- one host: this run used node budget 2 CPU 3g and client budget 1 CPU 2g, so budget comparisons are separate classes here, not independent hardware
- one client process with one fixed maintenance writer at 10 updates/s; readers are sequential; the reported sample latency includes retry and backoff time
- cluster nodes are containers on the same host, so these are not network or independent-hardware measurements; container colocation is not data-colocation evidence
- readers round-robin across 3 endpoints; the writer's endpoint policy is recorded separately
- physical rank-index arm yb-range-rank; an engine whose planner ignores the index is a different physical design, not a faster one
- YugabyteDB DIST/DEBUG storage counters (rows and index rows scanned) are parsed per cell; RPC-level counters are not aggregated
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
