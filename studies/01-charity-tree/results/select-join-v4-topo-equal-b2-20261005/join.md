# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 30 cells, 30 valid, engine `yugabyte`, topology `single`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-v4-topo-equal-b2-20261005`
- run tag: `run/study-01/select-join-v4-topo-equal-b2-20261005`
- environment: `host-zenbook-ux5406sa`
- repo commit: `ba821e7b686ac3adb1cda56dd45f15c8ae699ef8`
- engine version: `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit`
- database engine image: `docker.io/yugabytedb/yugabyte:2025.2.6.0-b111` @ `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`
- benchmark client image: `localhost/charitybench-select:select-join-v4-topo-equal-b2-20261005` @ `b6ea3c007cf459112eaa8e9e85381a7dc25a7710ff07648fde3048f747df6ab0`
- inputs digest: `81698f69d454c8a835180e6bdc7c6f7a0c1c3ed6785717fbb80932bc01c8eacd`
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
| yb-single | 505 | 2 |

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
| 20 (0.01%) | 2.1925 | missing | missing | missing | missing | missing | 2.1198 | 95.6558 | 2.3058 | missing |
| 2000 (1%) | 10.3871 | missing | missing | missing | missing | missing | 10.5466 | 99.4964 | 9.3068 | missing |
| 180000 (90%) | 114.3228 | missing | missing | missing | missing | missing | 121.1096 | 111.5058 | 114.8506 | missing |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | cte_joined_outer vs join_where | 2.1925 | 2.1198 | 0.9668 | false |
| 20 | cte_joined_outer_mat vs join_where | 2.1925 | 95.6558 | 43.6288 | false |
| 20 | cte_joined_outer_notmat vs join_where | 2.1925 | 2.3058 | 1.0517 | false |
| 2000 | cte_joined_outer vs join_where | 10.3871 | 10.5466 | 1.0154 | false |
| 2000 | cte_joined_outer_mat vs join_where | 10.3871 | 99.4964 | 9.5788 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 10.3871 | 9.3068 | 0.8960 | false |
| 180000 | cte_joined_outer vs join_where | 114.3228 | 121.1096 | 1.0594 | false |
| 180000 | cte_joined_outer_mat vs join_where | 114.3228 | 111.5058 | 0.9754 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 114.3228 | 114.8506 | 1.0046 | false |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_joined_outer_mat vs cte_joined_outer | 2.1198 | 95.6558 | 45.1248 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 2.1198 | 2.3058 | 1.0877 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 10.5466 | 99.4964 | 9.4339 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 10.5466 | 9.3068 | 0.8824 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 121.1096 | 111.5058 | 0.9207 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 121.1096 | 114.8506 | 0.9483 | false | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 7.6845 | 1 | `165b8050a2fd6b4a` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 9.4879 | 1 | `824304a7b71302c7` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | missing | | | | |
| default | 2 | 100% | missing | | | | |
| default | 4 | 1% | 8.4036 | 4 | `234070f21cc3690d` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 9.6585 | 4 | `c81b9cf47a904465` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | missing | | | | |
| default | 8 | 100% | missing | | | | |
| default | 16 | 1% | 11.1651 | 16 | `5098e06aca9f341e` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 10.7672 | 16 | `1ee382e4209e0358` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 8.6920 | 1 | `191fb6a15b84a67a` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 10.2125 | 1 | `ca7721014b58f77d` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | missing | | | | |
| mat | 2 | 100% | missing | | | | |
| mat | 4 | 1% | 9.4010 | 4 | `58bf2a4564cba067` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 9.6774 | 4 | `bd4703468faea855` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | missing | | | | |
| mat | 8 | 100% | missing | | | | |
| mat | 16 | 1% | 9.9975 | 16 | `0f9b623c5daef9c0` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 10.2790 | 16 | `44d7d8e2a8fca978` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 7.8425 | 1 | `13b38c6f2e35bcea` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 8.3364 | 1 | `9974c7f3ebed730f` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | missing | | | | |
| notmat | 2 | 100% | missing | | | | |
| notmat | 4 | 1% | 24.8587 | 4 | `5b7e0d71c3704261` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 36.8032 | 4 | `83283082ec9958a0` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | missing | | | | |
| notmat | 8 | 100% | missing | | | | |
| notmat | 16 | 1% | 87.6861 | 16 | `fa4368024121ecd2` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 101.1018 | 16 | `1c04a91b0b3dc73b` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 1.1311 | false |
| notmat | 1 | 1% | 1.0206 | false |
| mat | 1 | 100% | 1.0764 | false |
| notmat | 1 | 100% | 0.8786 | false |
| mat | 4 | 1% | 1.1187 | false |
| notmat | 4 | 1% | 2.9581 | false |
| mat | 4 | 100% | 1.0020 | false |
| notmat | 4 | 100% | 3.8105 | false |
| mat | 16 | 1% | 0.8954 | false |
| notmat | 16 | 1% | 7.8536 | false |
| mat | 16 | 100% | 0.9547 | false |
| notmat | 16 | 100% | 9.3898 | false |

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

- client: nr_periods 301, nr_throttled 120, throttled_usec 4535597 | quota 100000 100000
, memory 2147483648

- yb-single: nr_periods 379, nr_throttled 0, throttled_usec 0 | quota 600000 100000, memory 9663676416

## Access path per cell

_Engine-reported from each cell's own plan. The two scan flags are **whole-plan** booleans: they are true when any node of the plan is an index scan / sequential scan, including small parent-table scans, so they are not the rank access path. `storage rows scanned` is the engine's own count (YugabyteDB); PostgreSQL does not print one and it is shown as n/a. A cell that scanned the whole child table at a selective threshold is a different physical regime, not a faster query._

| axis | status | whole-plan index scan | whole-plan seq scan | storage rows scanned | index rows scanned | buffers hit/read | rows removed | retries |
|---|---|---|---|---|---|---|---|---|
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | false | true | 202050 | 0 | 0/0 | 0 | 1 |
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
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 110.6785 | 121.1096 | 221.6986 | 5 | 8.3 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 1.9213 | 2.1198 | 2.4483 | 5 | 471.7 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 8.9047 | 10.5466 | 21.5153 | 5 | 94.8 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 100.1868 | 111.5058 | 231.2937 | 5 | 9.0 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 83.5048 | 95.6558 | 128.0567 | 5 | 10.5 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 83.0383 | 99.4964 | 106.9724 | 5 | 10.1 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 112.8493 | 114.8506 | 124.6439 | 5 | 8.7 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 1.8128 | 2.3058 | 2.9170 | 5 | 433.7 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 8.6304 | 9.3068 | 10.3822 | 5 | 107.4 |
| family=bag;stmt=join_where;threshold=180000 | valid | 105.9588 | 114.3228 | 176.2042 | 5 | 8.7 |
| family=bag;stmt=join_where;threshold=20 | valid | 1.8592 | 2.1925 | 2.6709 | 5 | 456.1 |
| family=bag;stmt=join_where;threshold=2000 | valid | 8.5874 | 10.3871 | 12.7763 | 5 | 96.3 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 9.4064 | 11.1651 | 12.7921 | 5 | 89.6 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 9.3095 | 10.7672 | 11.9185 | 5 | 92.9 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 7.1495 | 7.6845 | 9.7127 | 5 | 130.1 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 9.3225 | 9.4879 | 15.1440 | 5 | 105.4 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 7.8946 | 8.4036 | 14.5986 | 5 | 119.0 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 8.9877 | 9.6585 | 13.1333 | 5 | 103.5 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 9.0735 | 9.9975 | 13.0611 | 5 | 100.0 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 9.6084 | 10.2790 | 12.8835 | 5 | 97.3 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 7.9389 | 8.6920 | 10.8671 | 5 | 115.0 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 7.6628 | 10.2125 | 11.7924 | 5 | 97.9 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 9.0223 | 9.4010 | 11.6476 | 5 | 106.4 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 9.4116 | 9.6774 | 11.5454 | 5 | 103.3 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 67.8028 | 87.6861 | 202.3943 | 5 | 11.4 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 80.7292 | 101.1018 | 104.6397 | 5 | 9.9 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 7.1246 | 7.8425 | 8.8864 | 5 | 127.5 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 8.1022 | 8.3364 | 9.6604 | 5 | 120.0 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 22.1108 | 24.8587 | 33.2092 | 5 | 40.2 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 29.4057 | 36.8032 | 54.5039 | 5 | 27.2 |

## Matched blended workload

- phase: `blended`; target offered rate: 10 updates/s
- offered 157, committed 157, failed 0 over 15.698 s
- offered rate 10.001/s, completed rate 10.001/s
- writer operations per endpoint: [157]
- reader operations overlapping a committed write: 71 of 211
- reader retries of transient conflicts: 1
- per key: 157 keys, min 1, max 1 per key
- ledger: 157 keys reconciled, 0 mismatch
- update latency: avg 3.594 ms, min 1.553, max 64.507 over 157

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42452221`, 5 interleaved rounds
  1. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  2. family=reuse;policy=notmat;refs=4;consumer=100%
  3. family=reuse;policy=default;refs=16;consumer=1%
  4. family=reuse;policy=default;refs=16;consumer=100%
  5. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  6. family=reuse;policy=mat;refs=4;consumer=1%
  7. family=reuse;policy=notmat;refs=1;consumer=1%
  8. family=bag;stmt=join_where;threshold=20
  9. family=bag;stmt=cte_joined_outer_mat;threshold=20
  10. family=reuse;policy=default;refs=4;consumer=100%
  11. family=bag;stmt=cte_joined_outer;threshold=2000
  12. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  13. family=bag;stmt=cte_joined_outer;threshold=180000
  14. family=reuse;policy=notmat;refs=16;consumer=1%
  15. family=reuse;policy=default;refs=1;consumer=1%
  16. family=reuse;policy=notmat;refs=4;consumer=1%
  17. family=reuse;policy=default;refs=4;consumer=1%
  18. family=bag;stmt=join_where;threshold=2000
  19. family=reuse;policy=mat;refs=1;consumer=1%
  20. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  21. family=bag;stmt=join_where;threshold=180000
  22. family=reuse;policy=mat;refs=4;consumer=100%
  23. family=reuse;policy=mat;refs=1;consumer=100%
  24. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  25. family=reuse;policy=notmat;refs=16;consumer=100%
  26. family=reuse;policy=default;refs=1;consumer=100%
  27. family=reuse;policy=mat;refs=16;consumer=100%
  28. family=reuse;policy=notmat;refs=1;consumer=100%
  29. family=reuse;policy=mat;refs=16;consumer=1%
  30. family=bag;stmt=cte_joined_outer;threshold=20

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
| family=bag;stmt=cte_joined_outer;threshold=180000 | j05_cte_joined_outer | valid | `72573f2574ab07db` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20 | j05_cte_joined_outer | valid | `3506d46afaa2982d` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=2000 | j05_cte_joined_outer | valid | `cb53b589edd26521` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | j05m_cte_joined_outer_mat | valid | `273b9b7bdc1f161d` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | j05m_cte_joined_outer_mat | valid | `cd6e1421342dbdba` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | j05m_cte_joined_outer_mat | valid | `7da0d7ed0e8ee1ac` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | j05n_cte_joined_outer_notmat | valid | `ea3a2ed8ecd32a4a` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | j05n_cte_joined_outer_notmat | valid | `2ca9f63dc484e7e6` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | j05n_cte_joined_outer_notmat | valid | `a689aebb83f44691` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=join_where;threshold=180000 | j01_join_where | valid | `ad1eb9c86da8bec9` |  |
| family=bag;stmt=join_where;threshold=20 | j01_join_where | valid | `6d34dc56f5227851` |  |
| family=bag;stmt=join_where;threshold=2000 | j01_join_where | valid | `47c6d3c85f32d892` |  |
| family=reuse;policy=default;refs=16;consumer=1% | r_default_n16_c1 | valid | `5098e06aca9f341e` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=16;consumer=100% | r_default_n16_c100 | valid | `1ee382e4209e0358` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=1;consumer=1% | r_default_n1_c1 | valid | `165b8050a2fd6b4a` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=1;consumer=100% | r_default_n1_c100 | valid | `824304a7b71302c7` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=4;consumer=1% | r_default_n4_c1 | valid | `234070f21cc3690d` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=100% | r_default_n4_c100 | valid | `c81b9cf47a904465` | policy=default; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=1% | r_mat_n16_c1 | valid | `0f9b623c5daef9c0` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=100% | r_mat_n16_c100 | valid | `44d7d8e2a8fca978` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=1% | r_mat_n1_c1 | valid | `191fb6a15b84a67a` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=100% | r_mat_n1_c100 | valid | `ca7721014b58f77d` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=1% | r_mat_n4_c1 | valid | `58bf2a4564cba067` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=100% | r_mat_n4_c100 | valid | `bd4703468faea855` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=notmat;refs=16;consumer=1% | r_notmat_n16_c1 | valid | `fa4368024121ecd2` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `1c04a91b0b3dc73b` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `13b38c6f2e35bcea` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `9974c7f3ebed730f` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `5b7e0d71c3704261` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `83283082ec9958a0` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, DIST, DEBUG, FORMAT TEXT` supported on yugabyte
- one host: this run used node budget 6 CPU 9g and client budget 1 CPU 2g, so budget comparisons are separate classes here, not independent hardware
- one client process with one fixed maintenance writer at 10 updates/s; readers are sequential; the reported sample latency includes retry and backoff time
- physical rank-index arm yb-range-rank; an engine whose planner ignores the index is a different physical design, not a faster one
- YugabyteDB DIST/DEBUG storage counters (rows and index rows scanned) are parsed per cell; RPC-level counters are not aggregated
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
