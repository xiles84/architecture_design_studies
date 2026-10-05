# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 30 cells, 30 valid, engine `yugabyte`, topology `single`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-v4-topo-equal-i2-20261005`
- run tag: `run/study-01/select-join-v4-topo-equal-i2-20261005`
- environment: `host-zenbook-ux5406sa`
- repo commit: `ba821e7b686ac3adb1cda56dd45f15c8ae699ef8`
- engine version: `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit`
- database engine image: `docker.io/yugabytedb/yugabyte:2025.2.6.0-b111` @ `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`
- benchmark client image: `localhost/charitybench-select:select-join-v4-topo-equal-i2-20261005` @ `b6ea3c007cf459112eaa8e9e85381a7dc25a7710ff07648fde3048f747df6ab0`
- inputs digest: `bffe080fefcb13b42d5b1f89a1906d70060534d141e743c0fe13234b1449c26e`
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
- setting `phase`: isolated
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
| yb-single | 504 | 1 |

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
| 20 (0.01%) | 1.6249 | missing | missing | missing | missing | missing | 1.7659 | 76.2120 | 1.5243 | missing |
| 2000 (1%) | 8.5039 | missing | missing | missing | missing | missing | 8.4416 | 78.8539 | 8.0655 | missing |
| 180000 (90%) | 101.0593 | missing | missing | missing | missing | missing | 98.6336 | 94.3181 | 103.1642 | missing |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | cte_joined_outer vs join_where | 1.6249 | 1.7659 | 1.0868 | false |
| 20 | cte_joined_outer_mat vs join_where | 1.6249 | 76.2120 | 46.9017 | false |
| 20 | cte_joined_outer_notmat vs join_where | 1.6249 | 1.5243 | 0.9381 | false |
| 2000 | cte_joined_outer vs join_where | 8.5039 | 8.4416 | 0.9927 | false |
| 2000 | cte_joined_outer_mat vs join_where | 8.5039 | 78.8539 | 9.2727 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 8.5039 | 8.0655 | 0.9484 | false |
| 180000 | cte_joined_outer vs join_where | 101.0593 | 98.6336 | 0.9760 | false |
| 180000 | cte_joined_outer_mat vs join_where | 101.0593 | 94.3181 | 0.9333 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 101.0593 | 103.1642 | 1.0208 | false |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_joined_outer_mat vs cte_joined_outer | 1.7659 | 76.2120 | 43.1565 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 1.7659 | 1.5243 | 0.8632 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 8.4416 | 78.8539 | 9.3411 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 8.4416 | 8.0655 | 0.9554 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 98.6336 | 94.3181 | 0.9562 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 98.6336 | 103.1642 | 1.0459 | false | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 5.7841 | 1 | `0ce1a8604f6b26ff` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 7.7434 | 1 | `2e031e22d52067d2` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | missing | | | | |
| default | 2 | 100% | missing | | | | |
| default | 4 | 1% | 7.1551 | 4 | `1a385c04691dff22` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 7.4379 | 4 | `f81dfceca90d9347` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | missing | | | | |
| default | 8 | 100% | missing | | | | |
| default | 16 | 1% | 8.3040 | 16 | `eb067f8f8c4a2bc5` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 9.8119 | 16 | `4aa9eca98877a1d5` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 7.4250 | 1 | `f1c932a97757a387` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 7.8509 | 1 | `38ed824288ff0287` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | missing | | | | |
| mat | 2 | 100% | missing | | | | |
| mat | 4 | 1% | 8.1427 | 4 | `597eb155107a67ca` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 7.9761 | 4 | `9e20ed89f3f2f89d` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | missing | | | | |
| mat | 8 | 100% | missing | | | | |
| mat | 16 | 1% | 8.0252 | 16 | `779901e885ba5fce` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 8.7639 | 16 | `0a9511ce1ad0ca81` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 5.9633 | 1 | `7475c5bf3ac5dae9` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 7.5505 | 1 | `855371359f8f0747` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | missing | | | | |
| notmat | 2 | 100% | missing | | | | |
| notmat | 4 | 1% | 22.2108 | 4 | `65ae7bf2e8bf2c9c` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 28.0347 | 4 | `7e6b6dac8c7fd1e6` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | missing | | | | |
| notmat | 8 | 100% | missing | | | | |
| notmat | 16 | 1% | 73.1574 | 16 | `c1c163bac4a8fdaf` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 79.8006 | 16 | `d81a2c7189467026` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 1.2837 | false |
| notmat | 1 | 1% | 1.0310 | false |
| mat | 1 | 100% | 1.0139 | false |
| notmat | 1 | 100% | 0.9751 | false |
| mat | 4 | 1% | 1.1380 | false |
| notmat | 4 | 1% | 3.1042 | false |
| mat | 4 | 100% | 1.0724 | false |
| notmat | 4 | 100% | 3.7692 | false |
| mat | 16 | 1% | 0.9664 | false |
| notmat | 16 | 1% | 8.8099 | false |
| mat | 16 | 100% | 0.8932 | false |
| notmat | 16 | 100% | 8.1330 | false |

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

Effective sharding flag (from the engine): `ysql_num_shards_per_tserver=2`. A relation with fewer tablets than this setting was created without an explicit split; the per-table `num_tablets` above is the authority.

## Throttling (cgroup `cpu.stat` deltas)

- client: nr_periods 264, nr_throttled 111, throttled_usec 4428811 | quota 100000 100000
, memory 2147483648

- yb-single: nr_periods 337, nr_throttled 0, throttled_usec 0 | quota 600000 100000, memory 9663676416

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
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 96.7278 | 98.6336 | 100.0947 | 5 | 10.1 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 1.5747 | 1.7659 | 1.8104 | 5 | 566.3 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 7.8617 | 8.4416 | 9.7107 | 5 | 118.5 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 90.3502 | 94.3181 | 94.6380 | 5 | 10.6 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 74.3743 | 76.2120 | 80.5294 | 5 | 13.1 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 75.1503 | 78.8539 | 81.7149 | 5 | 12.7 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 94.1969 | 103.1642 | 119.0111 | 5 | 9.7 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 1.2622 | 1.5243 | 1.6624 | 5 | 656.0 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 7.6135 | 8.0655 | 9.7987 | 5 | 124.0 |
| family=bag;stmt=join_where;threshold=180000 | valid | 95.6899 | 101.0593 | 110.7379 | 5 | 9.9 |
| family=bag;stmt=join_where;threshold=20 | valid | 1.4489 | 1.6249 | 1.7806 | 5 | 615.4 |
| family=bag;stmt=join_where;threshold=2000 | valid | 8.1952 | 8.5039 | 8.7846 | 5 | 117.6 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 8.1331 | 8.3040 | 9.0353 | 5 | 120.4 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 9.6159 | 9.8119 | 10.2762 | 5 | 101.9 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 5.6235 | 5.7841 | 6.4960 | 5 | 172.9 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 6.8198 | 7.7434 | 8.4409 | 5 | 129.1 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 6.7493 | 7.1551 | 7.8894 | 5 | 139.8 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 6.8198 | 7.4379 | 7.6282 | 5 | 134.4 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 7.6150 | 8.0252 | 8.4549 | 5 | 124.6 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 8.7050 | 8.7639 | 8.9906 | 5 | 114.1 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 6.8357 | 7.4250 | 8.1375 | 5 | 134.7 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 7.4234 | 7.8509 | 8.2500 | 5 | 127.4 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 7.4117 | 8.1427 | 8.8128 | 5 | 122.8 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 7.3920 | 7.9761 | 9.2087 | 5 | 125.4 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 66.8125 | 73.1574 | 79.6365 | 5 | 13.7 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 77.9085 | 79.8006 | 88.3838 | 5 | 12.5 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 5.6427 | 5.9633 | 9.6945 | 5 | 167.7 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 7.3714 | 7.5505 | 8.4167 | 5 | 132.4 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 20.5461 | 22.2108 | 25.4652 | 5 | 45.0 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 25.4817 | 28.0347 | 29.2056 | 5 | 35.7 |

## Matched blended workload

- phase: `isolated`; target offered rate: 10 updates/s
- offered 0, committed 0, failed 0 over 0.000 s
- offered rate 0.000/s, completed rate 0.000/s
- per key: no keys touched
- ledger: 0 keys reconciled, 0 mismatch

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42797484`, 5 interleaved rounds
  1. family=reuse;policy=notmat;refs=16;consumer=1%
  2. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  3. family=bag;stmt=join_where;threshold=2000
  4. family=reuse;policy=default;refs=16;consumer=1%
  5. family=bag;stmt=join_where;threshold=20
  6. family=bag;stmt=cte_joined_outer;threshold=2000
  7. family=reuse;policy=mat;refs=4;consumer=1%
  8. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  9. family=bag;stmt=join_where;threshold=180000
  10. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  11. family=reuse;policy=notmat;refs=4;consumer=100%
  12. family=reuse;policy=mat;refs=4;consumer=100%
  13. family=reuse;policy=default;refs=1;consumer=100%
  14. family=reuse;policy=mat;refs=1;consumer=1%
  15. family=reuse;policy=mat;refs=1;consumer=100%
  16. family=reuse;policy=notmat;refs=1;consumer=1%
  17. family=bag;stmt=cte_joined_outer;threshold=180000
  18. family=reuse;policy=default;refs=16;consumer=100%
  19. family=reuse;policy=mat;refs=16;consumer=1%
  20. family=reuse;policy=default;refs=4;consumer=100%
  21. family=bag;stmt=cte_joined_outer;threshold=20
  22. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  23. family=reuse;policy=notmat;refs=1;consumer=100%
  24. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  25. family=reuse;policy=mat;refs=16;consumer=100%
  26. family=reuse;policy=default;refs=1;consumer=1%
  27. family=reuse;policy=default;refs=4;consumer=1%
  28. family=reuse;policy=notmat;refs=16;consumer=100%
  29. family=reuse;policy=notmat;refs=4;consumer=1%
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
| family=bag;stmt=cte_joined_outer;threshold=180000 | j05_cte_joined_outer | valid | `5215d2081d0b01f3` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20 | j05_cte_joined_outer | valid | `fd163d7570b5ead0` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=2000 | j05_cte_joined_outer | valid | `f7a8694c07b95012` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | j05m_cte_joined_outer_mat | valid | `c97fdc034ed12f26` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | j05m_cte_joined_outer_mat | valid | `e91b3f3119e9bd03` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | j05m_cte_joined_outer_mat | valid | `83a06cbe33642728` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | j05n_cte_joined_outer_notmat | valid | `5d4436c2ca48ee97` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | j05n_cte_joined_outer_notmat | valid | `c6ca063f00bc15bc` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | j05n_cte_joined_outer_notmat | valid | `f48c14eecb1aa479` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=join_where;threshold=180000 | j01_join_where | valid | `f81d9299e5f308d6` |  |
| family=bag;stmt=join_where;threshold=20 | j01_join_where | valid | `fb35ae6862f674e0` |  |
| family=bag;stmt=join_where;threshold=2000 | j01_join_where | valid | `084a33fcc421e26e` |  |
| family=reuse;policy=default;refs=16;consumer=1% | r_default_n16_c1 | valid | `eb067f8f8c4a2bc5` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=16;consumer=100% | r_default_n16_c100 | valid | `4aa9eca98877a1d5` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=1;consumer=1% | r_default_n1_c1 | valid | `0ce1a8604f6b26ff` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=1;consumer=100% | r_default_n1_c100 | valid | `2e031e22d52067d2` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=4;consumer=1% | r_default_n4_c1 | valid | `1a385c04691dff22` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=100% | r_default_n4_c100 | valid | `f81dfceca90d9347` | policy=default; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=1% | r_mat_n16_c1 | valid | `779901e885ba5fce` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=100% | r_mat_n16_c100 | valid | `0a9511ce1ad0ca81` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=1% | r_mat_n1_c1 | valid | `f1c932a97757a387` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=100% | r_mat_n1_c100 | valid | `38ed824288ff0287` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=1% | r_mat_n4_c1 | valid | `597eb155107a67ca` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=100% | r_mat_n4_c100 | valid | `9e20ed89f3f2f89d` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=notmat;refs=16;consumer=1% | r_notmat_n16_c1 | valid | `c1c163bac4a8fdaf` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `d81a2c7189467026` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `7475c5bf3ac5dae9` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `855371359f8f0747` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `65ae7bf2e8bf2c9c` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `7e6b6dac8c7fd1e6` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, DIST, DEBUG, FORMAT TEXT` supported on yugabyte
- one host: this run used node budget 6 CPU 9g and client budget 1 CPU 2g, so budget comparisons are separate classes here, not independent hardware
- isolated phase: no concurrent writers in this run; the blended read/write phase is a separate labelled run
- physical rank-index arm yb-range-rank; an engine whose planner ignores the index is a different physical design, not a faster one
- YugabyteDB DIST/DEBUG storage counters (rows and index rows scanned) are parsed per cell; RPC-level counters are not aggregated
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
