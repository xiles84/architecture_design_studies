# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 30 cells, 30 valid, engine `yugabyte`, topology `cluster3`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-v4-topo-tablet-b3-20261005`
- run tag: `run/study-01/select-join-v4-topo-tablet-b3-20261005`
- environment: `host-zenbook-ux5406sa`
- repo commit: `ba821e7b686ac3adb1cda56dd45f15c8ae699ef8`
- engine version: `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit`
- database engine image: `docker.io/yugabytedb/yugabyte:2025.2.6.0-b111` @ `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`
- benchmark client image: `localhost/charitybench-select:select-join-v4-topo-tablet-b3-20261005` @ `b6ea3c007cf459112eaa8e9e85381a7dc25a7710ff07648fde3048f747df6ab0`
- inputs digest: `00fd931d2a48a57a27f502788a1ea816a3a69d7365f1016c6e129d4a8e0d9a6b`
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
- setting `placement`: noncolocated-single-tablet
- setting `prepared_plan_policy`: pgx simple protocol; no prepared statements
- setting `scale`: cte-domain
- setting `selectivity_mapping`: rank <= threshold; rank is a deterministic shuffle of 1..fact_count
- setting `shards_per_tserver`: ysql_num_shards_per_tserver=2
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
| 20 (0.01%) | 2.9713 | missing | missing | missing | missing | missing | 2.8027 | 128.7220 | 2.8511 | missing |
| 2000 (1%) | 11.7788 | missing | missing | missing | missing | missing | 9.5828 | 122.2600 | 11.6188 | missing |
| 180000 (90%) | 116.4572 | missing | missing | missing | missing | missing | 113.5581 | 105.9493 | 117.0366 | missing |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | cte_joined_outer vs join_where | 2.9713 | 2.8027 | 0.9432 | false |
| 20 | cte_joined_outer_mat vs join_where | 2.9713 | 128.7220 | 43.3218 | false |
| 20 | cte_joined_outer_notmat vs join_where | 2.9713 | 2.8511 | 0.9595 | false |
| 2000 | cte_joined_outer vs join_where | 11.7788 | 9.5828 | 0.8136 | false |
| 2000 | cte_joined_outer_mat vs join_where | 11.7788 | 122.2600 | 10.3796 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 11.7788 | 11.6188 | 0.9864 | false |
| 180000 | cte_joined_outer vs join_where | 116.4572 | 113.5581 | 0.9751 | false |
| 180000 | cte_joined_outer_mat vs join_where | 116.4572 | 105.9493 | 0.9098 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 116.4572 | 117.0366 | 1.0050 | false |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_joined_outer_mat vs cte_joined_outer | 2.8027 | 128.7220 | 45.9284 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 2.8027 | 2.8511 | 1.0173 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 9.5828 | 122.2600 | 12.7583 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 9.5828 | 11.6188 | 1.2125 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 113.5581 | 105.9493 | 0.9330 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 113.5581 | 117.0366 | 1.0306 | false | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 6.6946 | 1 | `815f9487c8765ca7` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 9.9197 | 1 | `0c0e81d35231e05c` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | missing | | | | |
| default | 2 | 100% | missing | | | | |
| default | 4 | 1% | 9.8833 | 4 | `08aa2969542c945c` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 10.0376 | 4 | `c3ea09c400cbc001` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | missing | | | | |
| default | 8 | 100% | missing | | | | |
| default | 16 | 1% | 10.0611 | 16 | `76da6aa991de1378` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 10.1862 | 16 | `f445a60a4f6e369a` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 9.5204 | 1 | `10ab0bef109d98d8` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 9.3698 | 1 | `dce5f0f725534eb4` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | missing | | | | |
| mat | 2 | 100% | missing | | | | |
| mat | 4 | 1% | 11.3920 | 4 | `4e4400e633b384a8` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 10.6557 | 4 | `8dbc5b50970963e2` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | missing | | | | |
| mat | 8 | 100% | missing | | | | |
| mat | 16 | 1% | 10.6292 | 16 | `feeb5c3bf295e741` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 12.0248 | 16 | `1b46058ceb1e7a03` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 8.3972 | 1 | `24529a0f778d4615` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 9.5145 | 1 | `f66c8ea61882a505` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | missing | | | | |
| notmat | 2 | 100% | missing | | | | |
| notmat | 4 | 1% | 27.5815 | 4 | `c8e0c949fd7ef0f9` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 31.0676 | 4 | `7282307b234c4a51` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | missing | | | | |
| notmat | 8 | 100% | missing | | | | |
| notmat | 16 | 1% | 85.9427 | 16 | `98a1977e30388e3d` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 88.2381 | 16 | `e0701941a11f5173` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 1.4221 | false |
| notmat | 1 | 1% | 1.2543 | false |
| mat | 1 | 100% | 0.9446 | false |
| notmat | 1 | 100% | 0.9592 | false |
| mat | 4 | 1% | 1.1527 | false |
| notmat | 4 | 1% | 2.7907 | false |
| mat | 4 | 100% | 1.0616 | false |
| notmat | 4 | 100% | 3.0951 | false |
| mat | 16 | 1% | 1.0565 | false |
| notmat | 16 | 1% | 8.5421 | false |
| mat | 16 | 100% | 1.1805 | false |
| notmat | 16 | 100% | 8.6625 | false |

## Physical arm

- `charity_pkey`: `CREATE UNIQUE INDEX charity_pkey ON selectjoin_domain.charity USING lsm (charity_id HASH)`
- `donation_person_idx`: `CREATE INDEX donation_person_idx ON selectjoin_domain.donation USING lsm (person_id HASH)`
- `donation_pkey`: `CREATE UNIQUE INDEX donation_pkey ON selectjoin_domain.donation USING lsm (donation_id HASH)`
- `donation_rank_idx`: `CREATE INDEX donation_rank_idx ON selectjoin_domain.donation USING lsm (rank ASC)`
- `person_pkey`: `CREATE UNIQUE INDEX person_pkey ON selectjoin_domain.person USING lsm (person_id HASH)`

- index arm `yb-range-rank`; placement `noncolocated-single-tablet`; node budget `2 CPU 3g`; client budget `1 CPU 2g`

| table | tablets | hash key columns | colocated | tablegroup | colocation id |
|---|---|---|---|---|---|
| selectjoin_domain.charity | 1 | 1 | false | unknown | unknown |
| selectjoin_domain.person | 1 | 1 | false | unknown | unknown |
| selectjoin_domain.donation | 1 | 1 | false | unknown | unknown |
| selectjoin_domain.donation_pkey | 0 | 0 | false | unknown | unknown |
| selectjoin_domain.donation_rank_idx | 1 | 0 | false | unknown | unknown |
| selectjoin_domain.donation_person_idx | 1 | 1 | false | unknown | unknown |

Physical-arm gate (default join must use the rank index at its three most selective thresholds):

- threshold 20: used
- threshold 2000: used

Tablet leaders per relation (placement evidence; the runner read these through yb-admin):

| relation | leader nodes |
|---|---|
| charity | yb-n2:9100 |
| person | yb-n2:9100 |
| donation | yb-n2:9100 |
| donation_rank_idx | yb-n2:9100 |
| donation_person_idx | yb-n2:9100 |

Effective sharding flag (from the engine): `ysql_num_shards_per_tserver=2`. A relation with fewer tablets than this setting was created without an explicit split; the per-table `num_tablets` above is the authority.

## Throttling (cgroup `cpu.stat` deltas)

- client: nr_periods 321, nr_throttled 119, throttled_usec 4324233 | quota 100000 100000
, memory 2147483648

- yb-n1: nr_periods 599, nr_throttled 5, throttled_usec 496966 | quota 200000 100000, memory 3221225472
- yb-n2: nr_periods 592, nr_throttled 17, throttled_usec 1563387 | quota 200000 100000, memory 3221225472
- yb-n3: nr_periods 586, nr_throttled 0, throttled_usec 0 | quota 200000 100000, memory 3221225472

## Access path per cell

_Engine-reported from each cell's own plan. The two scan flags are **whole-plan** booleans: they are true when any node of the plan is an index scan / sequential scan, including small parent-table scans, so they are not the rank access path. `storage rows scanned` is the engine's own count (YugabyteDB); PostgreSQL does not print one and it is shown as n/a. A cell that scanned the whole child table at a selective threshold is a different physical regime, not a faster query._

| axis | status | whole-plan index scan | whole-plan seq scan | storage rows scanned | index rows scanned | buffers hit/read | rows removed | retries |
|---|---|---|---|---|---|---|---|---|
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | false | true | 202050 | 0 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | true | true | 110 | 20 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | true | true | 6050 | 2000 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | false | true | 202050 | 0 | 0/0 | 20000 | 0 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | false | true | 202050 | 0 | 0/0 | 199980 | 0 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | false | true | 202050 | 0 | 0/0 | 198000 | 0 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | false | true | 202050 | 0 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | true | true | 110 | 20 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | true | true | 6050 | 2000 | 0/0 | 0 | 0 |
| family=bag;stmt=join_where;threshold=180000 | valid | false | true | 202050 | 0 | 0/0 | 0 | 0 |
| family=bag;stmt=join_where;threshold=20 | valid | true | true | 110 | 20 | 0/0 | 0 | 0 |
| family=bag;stmt=join_where;threshold=2000 | valid | true | true | 6050 | 2000 | 0/0 | 0 | 0 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | true | false | 4000 | 2000 | 0/0 | 1983 | 0 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | true | false | 4000 | 2000 | 0/0 | 0 | 0 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | true | false | 4000 | 2000 | 0/0 | 0 | 0 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | true | false | 4000 | 2000 | 0/0 | 0 | 0 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | true | false | 4000 | 2000 | 0/0 | 1983 | 0 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | true | false | 4000 | 2000 | 0/0 | 0 | 0 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | true | false | 4000 | 2000 | 0/0 | 1983 | 0 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | true | false | 4000 | 2000 | 0/0 | 0 | 0 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | true | false | 4000 | 2000 | 0/0 | 1983 | 0 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | true | false | 4000 | 2000 | 0/0 | 0 | 0 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | true | false | 4000 | 2000 | 0/0 | 1983 | 0 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | true | false | 4000 | 2000 | 0/0 | 0 | 0 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | true | false | 64000 | 2000 | 0/0 | 0 | 0 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | true | false | 64000 | 2000 | 0/0 | 0 | 0 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | true | false | 4000 | 2000 | 0/0 | 0 | 0 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | true | false | 4000 | 2000 | 0/0 | 0 | 0 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | true | false | 16000 | 2000 | 0/0 | 0 | 0 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | true | false | 16000 | 2000 | 0/0 | 0 | 0 |

## Latency samples and derived throughput

_The reader is sequential, so `queries/s` is 1000/median_ms, a derived rate and not a concurrent-capacity test. min/max are the five sample extremes. Each sample includes any retry and backoff time the statement needed after a transient conflict; per-cell retry counts appear in the access-path table._

| axis | status | min ms | median ms | max ms | samples | derived queries/s |
|---|---|---|---|---|---|---|
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 111.8479 | 113.5581 | 124.8886 | 5 | 8.8 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 2.5648 | 2.8027 | 3.0486 | 5 | 356.8 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 9.1310 | 9.5828 | 9.8304 | 5 | 104.4 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 103.4207 | 105.9493 | 128.6843 | 5 | 9.4 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 127.4073 | 128.7220 | 144.6758 | 5 | 7.8 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 119.8102 | 122.2600 | 134.2384 | 5 | 8.2 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 113.2218 | 117.0366 | 121.1096 | 5 | 8.5 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 2.5035 | 2.8511 | 3.1671 | 5 | 350.7 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 10.8060 | 11.6188 | 11.7222 | 5 | 86.1 |
| family=bag;stmt=join_where;threshold=180000 | valid | 109.3070 | 116.4572 | 127.2911 | 5 | 8.6 |
| family=bag;stmt=join_where;threshold=20 | valid | 2.6268 | 2.9713 | 3.5500 | 5 | 336.6 |
| family=bag;stmt=join_where;threshold=2000 | valid | 11.5654 | 11.7788 | 12.5090 | 5 | 84.9 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 9.5459 | 10.0611 | 11.2458 | 5 | 99.4 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 9.6229 | 10.1862 | 10.9003 | 5 | 98.2 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 6.4485 | 6.6946 | 7.5247 | 5 | 149.4 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 9.4356 | 9.9197 | 10.2322 | 5 | 100.8 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 9.6160 | 9.8833 | 10.3867 | 5 | 101.2 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 9.7565 | 10.0376 | 11.9807 | 5 | 99.6 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 10.1766 | 10.6292 | 11.0674 | 5 | 94.1 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 11.5339 | 12.0248 | 12.4827 | 5 | 83.2 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 9.0747 | 9.5204 | 9.6612 | 5 | 105.0 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 9.2495 | 9.3698 | 9.5340 | 5 | 106.7 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 10.1301 | 11.3920 | 12.1500 | 5 | 87.8 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 9.9783 | 10.6557 | 10.9110 | 5 | 93.8 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 80.2334 | 85.9427 | 94.2412 | 5 | 11.6 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 84.8267 | 88.2381 | 91.5079 | 5 | 11.3 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 7.4861 | 8.3972 | 8.9877 | 5 | 119.1 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 9.4743 | 9.5145 | 10.4711 | 5 | 105.1 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 26.4028 | 27.5815 | 29.6156 | 5 | 36.3 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 30.0074 | 31.0676 | 31.7134 | 5 | 32.2 |

## Matched blended workload

- phase: `blended`; target offered rate: 10 updates/s
- offered 162, committed 162, failed 0 over 16.170 s
- offered rate 10.019/s, completed rate 10.019/s
- writer operations per endpoint: [54 54 54]
- reader operations overlapping a committed write: 72 of 210
- reader retries of transient conflicts: 0
- per key: 162 keys, min 1, max 1 per key
- ledger: 162 keys reconciled, 0 mismatch
- update latency: avg 6.167 ms, min 2.353, max 155.229 over 162

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42045222`, 5 interleaved rounds
  1. family=reuse;policy=default;refs=4;consumer=100%
  2. family=bag;stmt=cte_joined_outer;threshold=180000
  3. family=reuse;policy=mat;refs=4;consumer=100%
  4. family=reuse;policy=notmat;refs=1;consumer=1%
  5. family=reuse;policy=default;refs=16;consumer=100%
  6. family=bag;stmt=join_where;threshold=20
  7. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  8. family=bag;stmt=cte_joined_outer;threshold=2000
  9. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  10. family=reuse;policy=mat;refs=1;consumer=1%
  11. family=reuse;policy=notmat;refs=4;consumer=100%
  12. family=reuse;policy=default;refs=4;consumer=1%
  13. family=reuse;policy=mat;refs=1;consumer=100%
  14. family=reuse;policy=default;refs=1;consumer=1%
  15. family=bag;stmt=cte_joined_outer;threshold=20
  16. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  17. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  18. family=bag;stmt=join_where;threshold=2000
  19. family=reuse;policy=notmat;refs=4;consumer=1%
  20. family=reuse;policy=notmat;refs=16;consumer=100%
  21. family=bag;stmt=cte_joined_outer_mat;threshold=20
  22. family=reuse;policy=notmat;refs=1;consumer=100%
  23. family=bag;stmt=join_where;threshold=180000
  24. family=reuse;policy=mat;refs=4;consumer=1%
  25. family=reuse;policy=notmat;refs=16;consumer=1%
  26. family=reuse;policy=default;refs=16;consumer=1%
  27. family=reuse;policy=default;refs=1;consumer=100%
  28. family=reuse;policy=mat;refs=16;consumer=1%
  29. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  30. family=reuse;policy=mat;refs=16;consumer=100%

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
| family=bag;stmt=cte_joined_outer;threshold=180000 | j05_cte_joined_outer | valid | `c90378db939241c0` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20 | j05_cte_joined_outer | valid | `e60abcee2bbfdb9f` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=2000 | j05_cte_joined_outer | valid | `3daeedecfb425cbe` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | j05m_cte_joined_outer_mat | valid | `9e8ba25bde3dde22` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | j05m_cte_joined_outer_mat | valid | `baf894b17610be98` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | j05m_cte_joined_outer_mat | valid | `0fc76c1c70f5684f` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | j05n_cte_joined_outer_notmat | valid | `d160df921eb79dc5` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | j05n_cte_joined_outer_notmat | valid | `da20263fe4f6d022` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | j05n_cte_joined_outer_notmat | valid | `f79a2f5d89afe57a` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=join_where;threshold=180000 | j01_join_where | valid | `2f628fb441f70e1f` |  |
| family=bag;stmt=join_where;threshold=20 | j01_join_where | valid | `978f88e05955bb0c` |  |
| family=bag;stmt=join_where;threshold=2000 | j01_join_where | valid | `9cbc9e158540c999` |  |
| family=reuse;policy=default;refs=16;consumer=1% | r_default_n16_c1 | valid | `76da6aa991de1378` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=16;consumer=100% | r_default_n16_c100 | valid | `f445a60a4f6e369a` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=1;consumer=1% | r_default_n1_c1 | valid | `815f9487c8765ca7` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=1;consumer=100% | r_default_n1_c100 | valid | `0c0e81d35231e05c` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=4;consumer=1% | r_default_n4_c1 | valid | `08aa2969542c945c` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=100% | r_default_n4_c100 | valid | `c3ea09c400cbc001` | policy=default; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=1% | r_mat_n16_c1 | valid | `feeb5c3bf295e741` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=100% | r_mat_n16_c100 | valid | `1b46058ceb1e7a03` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=1% | r_mat_n1_c1 | valid | `10ab0bef109d98d8` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=100% | r_mat_n1_c100 | valid | `dce5f0f725534eb4` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=1% | r_mat_n4_c1 | valid | `4e4400e633b384a8` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=100% | r_mat_n4_c100 | valid | `8dbc5b50970963e2` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=notmat;refs=16;consumer=1% | r_notmat_n16_c1 | valid | `98a1977e30388e3d` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `e0701941a11f5173` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `24529a0f778d4615` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `f66c8ea61882a505` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `c8e0c949fd7ef0f9` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `7282307b234c4a51` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, DIST, DEBUG, FORMAT TEXT` supported on yugabyte
- one host: this run used node budget 2 CPU 3g and client budget 1 CPU 2g, so budget comparisons are separate classes here, not independent hardware
- one client process with one fixed maintenance writer at 10 updates/s; readers are sequential; the reported sample latency includes retry and backoff time
- cluster nodes are containers on the same host, so these are not network or independent-hardware measurements; container colocation is not data-colocation evidence
- readers round-robin across 3 endpoints; the writer's endpoint policy is recorded separately
- physical rank-index arm yb-range-rank; an engine whose planner ignores the index is a different physical design, not a faster one
- YugabyteDB DIST/DEBUG storage counters (rows and index rows scanned) are parsed per cell; RPC-level counters are not aggregated
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
