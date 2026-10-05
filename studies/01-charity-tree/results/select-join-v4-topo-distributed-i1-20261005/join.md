# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 30 cells, 30 valid, engine `yugabyte`, topology `cluster3`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-v4-topo-distributed-i1-20261005`
- run tag: `run/study-01/select-join-v4-topo-distributed-i1-20261005`
- environment: `host-zenbook-ux5406sa`
- repo commit: `ba821e7b686ac3adb1cda56dd45f15c8ae699ef8`
- engine version: `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit`
- database engine image: `docker.io/yugabytedb/yugabyte:2025.2.6.0-b111` @ `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`
- benchmark client image: `localhost/charitybench-select:select-join-v4-topo-distributed-i1-20261005` @ `b6ea3c007cf459112eaa8e9e85381a7dc25a7710ff07648fde3048f747df6ab0`
- inputs digest: `1e89ee032810798eedb2bc9d723a7bd8386f9717146a338ac038608541db9805`
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
- setting `phase`: isolated
- setting `placement`: distributed
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
| yb-n1 | 168 | 1 |
| yb-n2 | 168 | 1 |
| yb-n3 | 168 | 1 |

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
| 20 (0.01%) | 3.3142 | missing | missing | missing | missing | missing | 2.8150 | 63.0618 | 2.6038 | missing |
| 2000 (1%) | 7.8363 | missing | missing | missing | missing | missing | 7.6579 | 65.0933 | 8.0165 | missing |
| 180000 (90%) | 75.4595 | missing | missing | missing | missing | missing | 78.6368 | 101.6240 | 81.1422 | missing |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | cte_joined_outer vs join_where | 3.3142 | 2.8150 | 0.8494 | false |
| 20 | cte_joined_outer_mat vs join_where | 3.3142 | 63.0618 | 19.0278 | false |
| 20 | cte_joined_outer_notmat vs join_where | 3.3142 | 2.6038 | 0.7856 | false |
| 2000 | cte_joined_outer vs join_where | 7.8363 | 7.6579 | 0.9772 | false |
| 2000 | cte_joined_outer_mat vs join_where | 7.8363 | 65.0933 | 8.3067 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 7.8363 | 8.0165 | 1.0230 | false |
| 180000 | cte_joined_outer vs join_where | 75.4595 | 78.6368 | 1.0421 | false |
| 180000 | cte_joined_outer_mat vs join_where | 75.4595 | 101.6240 | 1.3467 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 75.4595 | 81.1422 | 1.0753 | false |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_joined_outer_mat vs cte_joined_outer | 2.8150 | 63.0618 | 22.4022 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 2.8150 | 2.6038 | 0.9250 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 7.6579 | 65.0933 | 8.5001 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 7.6579 | 8.0165 | 1.0468 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 78.6368 | 101.6240 | 1.2923 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 78.6368 | 81.1422 | 1.0319 | false | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 5.4924 | 1 | `757ac310ff55667c` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 7.0748 | 1 | `f14adfed90d54f59` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | missing | | | | |
| default | 2 | 100% | missing | | | | |
| default | 4 | 1% | 6.7126 | 4 | `71926aea07aa689a` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 7.0435 | 4 | `80859c63e975d4c1` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | missing | | | | |
| default | 8 | 100% | missing | | | | |
| default | 16 | 1% | 7.7865 | 16 | `0574d513250608af` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 7.9741 | 16 | `aa650fb653aed351` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 6.9168 | 1 | `b006357d1cb08d14` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 6.6983 | 1 | `b2309a4b0dd4bc9b` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | missing | | | | |
| mat | 2 | 100% | missing | | | | |
| mat | 4 | 1% | 7.0160 | 4 | `019ee96dff167386` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 6.8674 | 4 | `b1bab7382a051b0b` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | missing | | | | |
| mat | 8 | 100% | missing | | | | |
| mat | 16 | 1% | 7.5942 | 16 | `27276d7ea0694d34` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 8.3630 | 16 | `d33a595f67255747` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 4.8997 | 1 | `69d471faf9b88e63` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 6.4239 | 1 | `64d7122c959b1544` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | missing | | | | |
| notmat | 2 | 100% | missing | | | | |
| notmat | 4 | 1% | 17.3076 | 4 | `a8c9b77f5b8104d0` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 23.4287 | 4 | `c8ecb3257ada41cb` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | missing | | | | |
| notmat | 8 | 100% | missing | | | | |
| notmat | 16 | 1% | 71.1265 | 16 | `4c3a79bcc7e31f31` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 89.1630 | 16 | `add7113932f09219` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 1.2594 | false |
| notmat | 1 | 1% | 0.8921 | false |
| mat | 1 | 100% | 0.9468 | false |
| notmat | 1 | 100% | 0.9080 | false |
| mat | 4 | 1% | 1.0452 | false |
| notmat | 4 | 1% | 2.5784 | false |
| mat | 4 | 100% | 0.9750 | false |
| notmat | 4 | 100% | 3.3263 | false |
| mat | 16 | 1% | 0.9753 | false |
| notmat | 16 | 1% | 9.1346 | false |
| mat | 16 | 100% | 1.0488 | false |
| notmat | 16 | 100% | 11.1815 | false |

## Physical arm

- `charity_pkey`: `CREATE UNIQUE INDEX charity_pkey ON selectjoin_domain.charity USING lsm (charity_id HASH)`
- `donation_person_idx`: `CREATE INDEX donation_person_idx ON selectjoin_domain.donation USING lsm (person_id HASH)`
- `donation_pkey`: `CREATE UNIQUE INDEX donation_pkey ON selectjoin_domain.donation USING lsm (donation_id HASH)`
- `donation_rank_idx`: `CREATE INDEX donation_rank_idx ON selectjoin_domain.donation USING lsm (rank ASC)`
- `person_pkey`: `CREATE UNIQUE INDEX person_pkey ON selectjoin_domain.person USING lsm (person_id HASH)`

- index arm `yb-range-rank`; placement `distributed`; node budget `2 CPU 3g`; client budget `1 CPU 2g`

| table | tablets | hash key columns | colocated | tablegroup | colocation id |
|---|---|---|---|---|---|
| selectjoin_domain.charity | 3 | 1 | false | unknown | unknown |
| selectjoin_domain.person | 3 | 1 | false | unknown | unknown |
| selectjoin_domain.donation | 3 | 1 | false | unknown | unknown |
| selectjoin_domain.donation_pkey | 0 | 0 | false | unknown | unknown |
| selectjoin_domain.donation_rank_idx | 3 | 0 | false | unknown | unknown |
| selectjoin_domain.donation_person_idx | 3 | 1 | false | unknown | unknown |

Physical-arm gate (default join must use the rank index at its three most selective thresholds):

- threshold 20: used
- threshold 2000: used

Tablet leaders per relation (placement evidence; the runner read these through yb-admin):

| relation | leader nodes |
|---|---|
| charity | yb-n1:9100, yb-n2:9100, yb-n3:9100 |
| person | yb-n1:9100, yb-n2:9100, yb-n3:9100 |
| donation | yb-n1:9100, yb-n2:9100, yb-n3:9100 |
| donation_rank_idx | yb-n1:9100, yb-n2:9100, yb-n3:9100 |
| donation_person_idx | yb-n1:9100, yb-n2:9100, yb-n3:9100 |

Effective sharding flag (from the engine): `ysql_num_shards_per_tserver=2`. A relation with fewer tablets than this setting was created without an explicit split; the per-table `num_tablets` above is the authority.

## Throttling (cgroup `cpu.stat` deltas)

- client: nr_periods 277, nr_throttled 128, throttled_usec 4130668 | quota 100000 100000
, memory 2147483648

- yb-n1: nr_periods 594, nr_throttled 10, throttled_usec 895807 | quota 200000 100000, memory 3221225472
- yb-n2: nr_periods 587, nr_throttled 4, throttled_usec 664926 | quota 200000 100000, memory 3221225472
- yb-n3: nr_periods 585, nr_throttled 1, throttled_usec 2385 | quota 200000 100000, memory 3221225472

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
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 78.2705 | 78.6368 | 83.3416 | 5 | 12.7 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 2.6131 | 2.8150 | 3.0819 | 5 | 355.2 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 6.9075 | 7.6579 | 8.3391 | 5 | 130.6 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 96.5671 | 101.6240 | 104.9613 | 5 | 9.8 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 62.4351 | 63.0618 | 65.0666 | 5 | 15.9 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 59.0162 | 65.0933 | 66.7447 | 5 | 15.4 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 74.7519 | 81.1422 | 90.5430 | 5 | 12.3 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 2.5710 | 2.6038 | 2.7758 | 5 | 384.1 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 7.7530 | 8.0165 | 12.7213 | 5 | 124.7 |
| family=bag;stmt=join_where;threshold=180000 | valid | 70.7538 | 75.4595 | 123.8268 | 5 | 13.3 |
| family=bag;stmt=join_where;threshold=20 | valid | 2.7006 | 3.3142 | 4.3248 | 5 | 301.7 |
| family=bag;stmt=join_where;threshold=2000 | valid | 7.2381 | 7.8363 | 9.5627 | 5 | 127.6 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 7.2582 | 7.7865 | 8.2113 | 5 | 128.4 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 7.7663 | 7.9741 | 8.9218 | 5 | 125.4 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 4.8680 | 5.4924 | 5.8803 | 5 | 182.1 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 6.3238 | 7.0748 | 7.2559 | 5 | 141.3 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 6.2390 | 6.7126 | 7.3620 | 5 | 149.0 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 6.4897 | 7.0435 | 7.6383 | 5 | 142.0 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 7.1912 | 7.5942 | 8.2007 | 5 | 131.7 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 8.1200 | 8.3630 | 8.7601 | 5 | 119.6 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 6.3470 | 6.9168 | 7.5800 | 5 | 144.6 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 6.2985 | 6.6983 | 8.6921 | 5 | 149.3 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 6.5888 | 7.0160 | 7.9504 | 5 | 142.5 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 6.6996 | 6.8674 | 7.7615 | 5 | 145.6 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 68.8700 | 71.1265 | 71.8697 | 5 | 14.1 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 83.9356 | 89.1630 | 94.9656 | 5 | 11.2 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 4.6825 | 4.8997 | 5.4664 | 5 | 204.1 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 6.3233 | 6.4239 | 6.7854 | 5 | 155.7 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 15.7009 | 17.3076 | 18.9194 | 5 | 57.8 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 22.6900 | 23.4287 | 26.2472 | 5 | 42.7 |

## Matched blended workload

- phase: `isolated`; target offered rate: 10 updates/s
- offered 0, committed 0, failed 0 over 0.000 s
- offered rate 0.000/s, completed rate 0.000/s
- per key: no keys touched
- ledger: 0 keys reconciled, 0 mismatch

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42363330`, 5 interleaved rounds
  1. family=reuse;policy=notmat;refs=16;consumer=100%
  2. family=bag;stmt=join_where;threshold=2000
  3. family=bag;stmt=join_where;threshold=20
  4. family=reuse;policy=mat;refs=1;consumer=100%
  5. family=reuse;policy=mat;refs=1;consumer=1%
  6. family=bag;stmt=cte_joined_outer;threshold=180000
  7. family=bag;stmt=cte_joined_outer_mat;threshold=20
  8. family=reuse;policy=notmat;refs=4;consumer=100%
  9. family=reuse;policy=mat;refs=16;consumer=100%
  10. family=reuse;policy=mat;refs=16;consumer=1%
  11. family=reuse;policy=notmat;refs=4;consumer=1%
  12. family=bag;stmt=cte_joined_outer;threshold=2000
  13. family=reuse;policy=default;refs=16;consumer=100%
  14. family=bag;stmt=cte_joined_outer;threshold=20
  15. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  16. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  17. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  18. family=reuse;policy=default;refs=1;consumer=100%
  19. family=reuse;policy=notmat;refs=16;consumer=1%
  20. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  21. family=reuse;policy=mat;refs=4;consumer=1%
  22. family=reuse;policy=default;refs=4;consumer=1%
  23. family=reuse;policy=mat;refs=4;consumer=100%
  24. family=reuse;policy=notmat;refs=1;consumer=1%
  25. family=reuse;policy=default;refs=4;consumer=100%
  26. family=reuse;policy=default;refs=1;consumer=1%
  27. family=bag;stmt=join_where;threshold=180000
  28. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  29. family=reuse;policy=notmat;refs=1;consumer=100%
  30. family=reuse;policy=default;refs=16;consumer=1%

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
| family=bag;stmt=cte_joined_outer;threshold=180000 | j05_cte_joined_outer | valid | `7854d1e97718c574` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20 | j05_cte_joined_outer | valid | `1e582e1e73296622` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=2000 | j05_cte_joined_outer | valid | `54ab90bb9e7560c2` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | j05m_cte_joined_outer_mat | valid | `7e12ac33bb9706fa` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | j05m_cte_joined_outer_mat | valid | `2eab7e4e1f03f259` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | j05m_cte_joined_outer_mat | valid | `2f8ddd8fdefecd61` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | j05n_cte_joined_outer_notmat | valid | `95971879da10b0a1` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | j05n_cte_joined_outer_notmat | valid | `8c3bcbb15b44ef23` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | j05n_cte_joined_outer_notmat | valid | `42b9a227ad8b92c2` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=join_where;threshold=180000 | j01_join_where | valid | `c20af64250a55249` |  |
| family=bag;stmt=join_where;threshold=20 | j01_join_where | valid | `9ada89e18e52b3b1` |  |
| family=bag;stmt=join_where;threshold=2000 | j01_join_where | valid | `afa50a9364dc1a70` |  |
| family=reuse;policy=default;refs=16;consumer=1% | r_default_n16_c1 | valid | `0574d513250608af` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=16;consumer=100% | r_default_n16_c100 | valid | `aa650fb653aed351` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=1;consumer=1% | r_default_n1_c1 | valid | `757ac310ff55667c` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=1;consumer=100% | r_default_n1_c100 | valid | `f14adfed90d54f59` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=4;consumer=1% | r_default_n4_c1 | valid | `71926aea07aa689a` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=100% | r_default_n4_c100 | valid | `80859c63e975d4c1` | policy=default; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=1% | r_mat_n16_c1 | valid | `27276d7ea0694d34` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=100% | r_mat_n16_c100 | valid | `d33a595f67255747` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=1% | r_mat_n1_c1 | valid | `b006357d1cb08d14` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=100% | r_mat_n1_c100 | valid | `b2309a4b0dd4bc9b` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=1% | r_mat_n4_c1 | valid | `019ee96dff167386` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=100% | r_mat_n4_c100 | valid | `b1bab7382a051b0b` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=notmat;refs=16;consumer=1% | r_notmat_n16_c1 | valid | `4c3a79bcc7e31f31` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `add7113932f09219` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `69d471faf9b88e63` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `64d7122c959b1544` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `a8c9b77f5b8104d0` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `c8ecb3257ada41cb` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, DIST, DEBUG, FORMAT TEXT` supported on yugabyte
- distributed placement gate: three tablets confirmed by the engine, leader spread confirmed by the runner
- one host: this run used node budget 2 CPU 3g and client budget 1 CPU 2g, so budget comparisons are separate classes here, not independent hardware
- isolated phase: no concurrent writers in this run; the blended read/write phase is a separate labelled run
- cluster nodes are containers on the same host, so these are not network or independent-hardware measurements; container colocation is not data-colocation evidence
- readers round-robin across 3 endpoints; the writer's endpoint policy is recorded separately
- physical rank-index arm yb-range-rank; an engine whose planner ignores the index is a different physical design, not a faster one
- YugabyteDB DIST/DEBUG storage counters (rows and index rows scanned) are parsed per cell; RPC-level counters are not aggregated
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
