# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 30 cells, 30 valid, engine `yugabyte`, topology `single`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-v4-topo-equal-b3-20261005`
- run tag: `run/study-01/select-join-v4-topo-equal-b3-20261005`
- environment: `host-zenbook-ux5406sa`
- repo commit: `ba821e7b686ac3adb1cda56dd45f15c8ae699ef8`
- engine version: `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit`
- database engine image: `docker.io/yugabytedb/yugabyte:2025.2.6.0-b111` @ `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`
- benchmark client image: `localhost/charitybench-select:select-join-v4-topo-equal-b3-20261005` @ `b6ea3c007cf459112eaa8e9e85381a7dc25a7710ff07648fde3048f747df6ab0`
- inputs digest: `f6b9c08cf1e49b8c28b767a485a5bb4fbeee295d6c2605ab8743972334a3e1dd`
- dataset fingerprint: `fcf26e08f377ec51`
- client CPU quota: `100000 100000`, memory: `2147483648`
- setting `client_budget`: 1 CPU 2g
- setting `client_workers`: 1 sequential reader; 0 writers
- setting `connection_pool`: 2 per endpoint = reader/setup + reserve
- setting `customers`: 2000
- setting `dataset`: study01-person-donation-charity
- setting `distributed_rpc`: DIST/DEBUG storage counters parsed per cell (Storage Rows Scanned, Storage Index Rows Scanned); RPC-level counters not aggregated
- setting `endpoint_mode`: single endpoint
- setting `endpoints`: yb-single
- setting `fact_count`: 200000
- setting `index_arm`: yb-range-rank
- setting `index_set`: child rank index and parent-key index, primary keys; no per-spelling index
- setting `node_budget`: 6 CPU 9g
- setting `phase`: blended
- setting `placement`: noncolocated-single-tablet
- setting `prepared_plan_policy`: pgx simple protocol; no prepared statements
- setting `scale`: cte-domain
- setting `selectivity_mapping`: rank <= threshold; rank is a deterministic shuffle of 1..fact_count
- setting `shards_per_tserver`: ysql_num_shards_per_tserver=2
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
| yb-single | 504 | 2 |

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
| 20 (0.01%) | 1.8150 | missing | missing | missing | missing | missing | 1.8870 | 80.5408 | 1.8585 | missing |
| 2000 (1%) | 8.7390 | missing | missing | missing | missing | missing | 9.2359 | 80.0139 | 9.5321 | missing |
| 180000 (90%) | 109.3388 | missing | missing | missing | missing | missing | 106.6266 | 103.5151 | 105.5875 | missing |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | cte_joined_outer vs join_where | 1.8150 | 1.8870 | 1.0397 | false |
| 20 | cte_joined_outer_mat vs join_where | 1.8150 | 80.5408 | 44.3749 | false |
| 20 | cte_joined_outer_notmat vs join_where | 1.8150 | 1.8585 | 1.0239 | false |
| 2000 | cte_joined_outer vs join_where | 8.7390 | 9.2359 | 1.0569 | false |
| 2000 | cte_joined_outer_mat vs join_where | 8.7390 | 80.0139 | 9.1559 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 8.7390 | 9.5321 | 1.0908 | false |
| 180000 | cte_joined_outer vs join_where | 109.3388 | 106.6266 | 0.9752 | false |
| 180000 | cte_joined_outer_mat vs join_where | 109.3388 | 103.5151 | 0.9467 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 109.3388 | 105.5875 | 0.9657 | false |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_joined_outer_mat vs cte_joined_outer | 1.8870 | 80.5408 | 42.6821 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 1.8870 | 1.8585 | 0.9849 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 9.2359 | 80.0139 | 8.6633 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 9.2359 | 9.5321 | 1.0321 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 106.6266 | 103.5151 | 0.9708 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 106.6266 | 105.5875 | 0.9903 | false | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 7.3748 | 1 | `7dd1ece36a97410a` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 8.5258 | 1 | `019d12689e3d33d5` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | missing | | | | |
| default | 2 | 100% | missing | | | | |
| default | 4 | 1% | 8.6189 | 4 | `443725539e714bcc` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 8.0977 | 4 | `5cb9d9e0db2e7ceb` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | missing | | | | |
| default | 8 | 100% | missing | | | | |
| default | 16 | 1% | 9.0766 | 16 | `dbf26e4eab2a70bc` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 9.6206 | 16 | `11b30bdcec4ce8d2` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 7.3748 | 1 | `3911a9a67a0ec6a8` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 8.1132 | 1 | `3325749f5ffe159c` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | missing | | | | |
| mat | 2 | 100% | missing | | | | |
| mat | 4 | 1% | 8.2582 | 4 | `d9c5fd4016a9df91` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 8.4199 | 4 | `85e36f115b5d0b80` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | missing | | | | |
| mat | 8 | 100% | missing | | | | |
| mat | 16 | 1% | 8.6664 | 16 | `c4735e177f2f1eaa` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 9.3057 | 16 | `692563ea807296b5` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 7.4604 | 1 | `717a4ddb77d1487b` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 7.7323 | 1 | `3da904034315f440` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | missing | | | | |
| notmat | 2 | 100% | missing | | | | |
| notmat | 4 | 1% | 21.1959 | 4 | `2f519d8b3f051d20` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 27.6865 | 4 | `f7c47f59ba3c2640` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | missing | | | | |
| notmat | 8 | 100% | missing | | | | |
| notmat | 16 | 1% | 69.8838 | 16 | `2b2ff2dd0de51d55` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 83.0334 | 16 | `2d2f609ee1fabfbc` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 1.0000 | false |
| notmat | 1 | 1% | 1.0116 | false |
| mat | 1 | 100% | 0.9516 | false |
| notmat | 1 | 100% | 0.9069 | false |
| mat | 4 | 1% | 0.9581 | false |
| notmat | 4 | 1% | 2.4592 | false |
| mat | 4 | 100% | 1.0398 | false |
| notmat | 4 | 100% | 3.4190 | false |
| mat | 16 | 1% | 0.9548 | false |
| notmat | 16 | 1% | 7.6994 | false |
| mat | 16 | 100% | 0.9673 | false |
| notmat | 16 | 100% | 8.6308 | false |

## Physical arm

- `charity_pkey`: `CREATE UNIQUE INDEX charity_pkey ON selectjoin_domain.charity USING lsm (charity_id HASH)`
- `donation_person_idx`: `CREATE INDEX donation_person_idx ON selectjoin_domain.donation USING lsm (person_id HASH)`
- `donation_pkey`: `CREATE UNIQUE INDEX donation_pkey ON selectjoin_domain.donation USING lsm (donation_id HASH)`
- `donation_rank_idx`: `CREATE INDEX donation_rank_idx ON selectjoin_domain.donation USING lsm (rank ASC)`
- `person_pkey`: `CREATE UNIQUE INDEX person_pkey ON selectjoin_domain.person USING lsm (person_id HASH)`

- index arm `yb-range-rank`; placement `noncolocated-single-tablet`; node budget `6 CPU 9g`; client budget `1 CPU 2g`

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
| charity | yb-single:9100 |
| person | yb-single:9100 |
| donation | yb-single:9100 |
| donation_rank_idx | yb-single:9100 |
| donation_person_idx | yb-single:9100 |

Effective sharding flag (from the engine): `ysql_num_shards_per_tserver=2`. A relation with fewer tablets than this setting was created without an explicit split; the per-table `num_tablets` above is the authority.

## Throttling (cgroup `cpu.stat` deltas)

- client: nr_periods 277, nr_throttled 113, throttled_usec 3388981 | quota 100000 100000
, memory 2147483648

- yb-single: nr_periods 345, nr_throttled 0, throttled_usec 0 | quota 600000 100000, memory 9663676416

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
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 103.9652 | 106.6266 | 115.7923 | 5 | 9.4 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 1.6976 | 1.8870 | 1.9239 | 5 | 529.9 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 8.7087 | 9.2359 | 12.1279 | 5 | 108.3 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 101.1656 | 103.5151 | 104.9719 | 5 | 9.7 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 76.4148 | 80.5408 | 91.1068 | 5 | 12.4 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 75.7421 | 80.0139 | 82.6962 | 5 | 12.5 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 99.5254 | 105.5875 | 110.3651 | 5 | 9.5 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 1.7467 | 1.8585 | 31.0938 | 5 | 538.1 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 9.0857 | 9.5321 | 10.5776 | 5 | 104.9 |
| family=bag;stmt=join_where;threshold=180000 | valid | 102.5839 | 109.3388 | 113.7941 | 5 | 9.1 |
| family=bag;stmt=join_where;threshold=20 | valid | 1.6597 | 1.8150 | 1.8798 | 5 | 551.0 |
| family=bag;stmt=join_where;threshold=2000 | valid | 8.6815 | 8.7390 | 9.7147 | 5 | 114.4 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 8.3716 | 9.0766 | 9.1832 | 5 | 110.2 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 9.4881 | 9.6206 | 20.3704 | 5 | 103.9 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 6.2891 | 7.3748 | 8.2579 | 5 | 135.6 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 8.1233 | 8.5258 | 8.7493 | 5 | 117.3 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 8.2823 | 8.6189 | 10.5861 | 5 | 116.0 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 7.6613 | 8.0977 | 9.7205 | 5 | 123.5 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 8.3574 | 8.6664 | 8.7959 | 5 | 115.4 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 8.8835 | 9.3057 | 10.0128 | 5 | 107.5 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 7.0322 | 7.3748 | 7.4896 | 5 | 135.6 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 7.9639 | 8.1132 | 9.3148 | 5 | 123.3 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 7.4098 | 8.2582 | 8.4083 | 5 | 121.1 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 8.3840 | 8.4199 | 8.7262 | 5 | 118.8 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 69.3570 | 69.8838 | 82.1654 | 5 | 14.3 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 77.5516 | 83.0334 | 90.0001 | 5 | 12.0 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 6.8883 | 7.4604 | 7.9616 | 5 | 134.0 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 7.2838 | 7.7323 | 9.0767 | 5 | 129.3 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 19.5781 | 21.1959 | 23.3729 | 5 | 47.2 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 26.6699 | 27.6865 | 31.1805 | 5 | 36.1 |

## Matched blended workload

- phase: `blended`; target offered rate: 10 updates/s
- offered 137, committed 137, failed 0 over 13.632 s
- offered rate 10.050/s, completed rate 10.050/s
- writer operations per endpoint: [137]
- reader operations overlapping a committed write: 67 of 210
- reader retries of transient conflicts: 0
- per key: 137 keys, min 1, max 1 per key
- ledger: 137 keys reconciled, 0 mismatch
- update latency: avg 3.319 ms, min 1.512, max 72.091 over 137

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42701616`, 5 interleaved rounds
  1. family=reuse;policy=notmat;refs=1;consumer=100%
  2. family=bag;stmt=join_where;threshold=180000
  3. family=reuse;policy=notmat;refs=4;consumer=100%
  4. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  5. family=bag;stmt=join_where;threshold=2000
  6. family=reuse;policy=default;refs=16;consumer=100%
  7. family=reuse;policy=mat;refs=16;consumer=1%
  8. family=reuse;policy=default;refs=16;consumer=1%
  9. family=bag;stmt=cte_joined_outer;threshold=180000
  10. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  11. family=reuse;policy=notmat;refs=4;consumer=1%
  12. family=reuse;policy=notmat;refs=16;consumer=100%
  13. family=reuse;policy=notmat;refs=1;consumer=1%
  14. family=bag;stmt=cte_joined_outer;threshold=20
  15. family=bag;stmt=cte_joined_outer_mat;threshold=20
  16. family=reuse;policy=mat;refs=1;consumer=100%
  17. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  18. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  19. family=reuse;policy=default;refs=1;consumer=100%
  20. family=reuse;policy=mat;refs=1;consumer=1%
  21. family=reuse;policy=mat;refs=16;consumer=100%
  22. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  23. family=reuse;policy=mat;refs=4;consumer=100%
  24. family=reuse;policy=notmat;refs=16;consumer=1%
  25. family=reuse;policy=default;refs=4;consumer=1%
  26. family=reuse;policy=mat;refs=4;consumer=1%
  27. family=bag;stmt=join_where;threshold=20
  28. family=bag;stmt=cte_joined_outer;threshold=2000
  29. family=reuse;policy=default;refs=1;consumer=1%
  30. family=reuse;policy=default;refs=4;consumer=100%

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
| family=bag;stmt=cte_joined_outer;threshold=180000 | j05_cte_joined_outer | valid | `3de07e000ef97606` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20 | j05_cte_joined_outer | valid | `4621999ec9fd6919` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=2000 | j05_cte_joined_outer | valid | `7e579dfdf988d4cd` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | j05m_cte_joined_outer_mat | valid | `2e484f7c9b3615e9` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | j05m_cte_joined_outer_mat | valid | `2fc591c9d2017d55` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | j05m_cte_joined_outer_mat | valid | `94cdb3de24e2e751` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | j05n_cte_joined_outer_notmat | valid | `7a9e35bf28c736ac` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | j05n_cte_joined_outer_notmat | valid | `7749490460f1ff09` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | j05n_cte_joined_outer_notmat | valid | `197c671a91569642` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=join_where;threshold=180000 | j01_join_where | valid | `bb13a8e387d82e0b` |  |
| family=bag;stmt=join_where;threshold=20 | j01_join_where | valid | `9df18e7a3c863198` |  |
| family=bag;stmt=join_where;threshold=2000 | j01_join_where | valid | `87bd1b0481f45b33` |  |
| family=reuse;policy=default;refs=16;consumer=1% | r_default_n16_c1 | valid | `dbf26e4eab2a70bc` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=16;consumer=100% | r_default_n16_c100 | valid | `11b30bdcec4ce8d2` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=1;consumer=1% | r_default_n1_c1 | valid | `7dd1ece36a97410a` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=1;consumer=100% | r_default_n1_c100 | valid | `019d12689e3d33d5` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=4;consumer=1% | r_default_n4_c1 | valid | `443725539e714bcc` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=100% | r_default_n4_c100 | valid | `5cb9d9e0db2e7ceb` | policy=default; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=1% | r_mat_n16_c1 | valid | `c4735e177f2f1eaa` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=100% | r_mat_n16_c100 | valid | `692563ea807296b5` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=1% | r_mat_n1_c1 | valid | `3911a9a67a0ec6a8` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=100% | r_mat_n1_c100 | valid | `3325749f5ffe159c` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=1% | r_mat_n4_c1 | valid | `d9c5fd4016a9df91` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=100% | r_mat_n4_c100 | valid | `85e36f115b5d0b80` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=notmat;refs=16;consumer=1% | r_notmat_n16_c1 | valid | `2b2ff2dd0de51d55` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `2d2f609ee1fabfbc` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `717a4ddb77d1487b` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `3da904034315f440` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `2f519d8b3f051d20` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `f7c47f59ba3c2640` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, DIST, DEBUG, FORMAT TEXT` supported on yugabyte
- one host: this run used node budget 6 CPU 9g and client budget 1 CPU 2g, so budget comparisons are separate classes here, not independent hardware
- one client process with one fixed maintenance writer at 10 updates/s; readers are sequential; the reported sample latency includes retry and backoff time
- physical rank-index arm yb-range-rank; an engine whose planner ignores the index is a different physical design, not a faster one
- YugabyteDB DIST/DEBUG storage counters (rows and index rows scanned) are parsed per cell; RPC-level counters are not aggregated
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
