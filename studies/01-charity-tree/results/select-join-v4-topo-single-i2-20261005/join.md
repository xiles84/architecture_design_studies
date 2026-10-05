# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 30 cells, 30 valid, engine `yugabyte`, topology `single`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-v4-topo-single-i2-20261005`
- run tag: `run/study-01/select-join-v4-topo-single-i2-20261005`
- environment: `host-zenbook-ux5406sa`
- repo commit: `ba821e7b686ac3adb1cda56dd45f15c8ae699ef8`
- engine version: `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit`
- database engine image: `docker.io/yugabytedb/yugabyte:2025.2.6.0-b111` @ `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`
- benchmark client image: `localhost/charitybench-select:select-join-v4-topo-single-i2-20261005` @ `b6ea3c007cf459112eaa8e9e85381a7dc25a7710ff07648fde3048f747df6ab0`
- inputs digest: `a63c107a532494fd6f8f9bbe8f49ca6e095cc14c0bead22fc848ac0e6a61693e`
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
- setting `node_budget`: 2 CPU 3g
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
| 20 (0.01%) | 1.7054 | missing | missing | missing | missing | missing | 1.9701 | 86.7975 | 1.8743 | missing |
| 2000 (1%) | 9.7444 | missing | missing | missing | missing | missing | 10.9920 | 91.3876 | 10.2569 | missing |
| 180000 (90%) | 126.7814 | missing | missing | missing | missing | missing | 131.4230 | 113.0238 | 108.0538 | missing |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | cte_joined_outer vs join_where | 1.7054 | 1.9701 | 1.1552 | false |
| 20 | cte_joined_outer_mat vs join_where | 1.7054 | 86.7975 | 50.8956 | false |
| 20 | cte_joined_outer_notmat vs join_where | 1.7054 | 1.8743 | 1.0990 | false |
| 2000 | cte_joined_outer vs join_where | 9.7444 | 10.9920 | 1.1280 | false |
| 2000 | cte_joined_outer_mat vs join_where | 9.7444 | 91.3876 | 9.3785 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 9.7444 | 10.2569 | 1.0526 | false |
| 180000 | cte_joined_outer vs join_where | 126.7814 | 131.4230 | 1.0366 | false |
| 180000 | cte_joined_outer_mat vs join_where | 126.7814 | 113.0238 | 0.8915 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 126.7814 | 108.0538 | 0.8523 | false |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_joined_outer_mat vs cte_joined_outer | 1.9701 | 86.7975 | 44.0578 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 1.9701 | 1.8743 | 0.9514 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 10.9920 | 91.3876 | 8.3140 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 10.9920 | 10.2569 | 0.9331 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 131.4230 | 113.0238 | 0.8600 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 131.4230 | 108.0538 | 0.8222 | false | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 7.4711 | 1 | `07c5b7241fdbac36` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 9.0043 | 1 | `6cc19b7ffe37dd7b` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | missing | | | | |
| default | 2 | 100% | missing | | | | |
| default | 4 | 1% | 10.2502 | 4 | `835a6bb9bba0714d` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 9.3801 | 4 | `1226eea54f78e5c4` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | missing | | | | |
| default | 8 | 100% | missing | | | | |
| default | 16 | 1% | 10.5321 | 16 | `ad40fe9ce4de55bf` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 9.9350 | 16 | `a32176d2253f9ef6` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 8.7035 | 1 | `9e689d6102b289af` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 8.6382 | 1 | `84436f44d02b0ba8` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | missing | | | | |
| mat | 2 | 100% | missing | | | | |
| mat | 4 | 1% | 9.9818 | 4 | `b2d7d291ebb7595c` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 9.3092 | 4 | `983a346bca42fa5c` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | missing | | | | |
| mat | 8 | 100% | missing | | | | |
| mat | 16 | 1% | 10.0384 | 16 | `852458c995675e5e` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 10.2148 | 16 | `356b3985aa9da7be` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 7.2312 | 1 | `38801e01ecb87a6e` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 9.5693 | 1 | `1cbbf26a53675d05` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | missing | | | | |
| notmat | 2 | 100% | missing | | | | |
| notmat | 4 | 1% | 26.7395 | 4 | `c7d85c7f9c9c6609` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 32.6548 | 4 | `64a865cddd01176f` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | missing | | | | |
| notmat | 8 | 100% | missing | | | | |
| notmat | 16 | 1% | 119.7908 | 16 | `6e45e4443c7c8fdd` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 94.6692 | 16 | `b60eaf55f8df19be` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 1.1650 | false |
| notmat | 1 | 1% | 0.9679 | false |
| mat | 1 | 100% | 0.9593 | false |
| notmat | 1 | 100% | 1.0627 | false |
| mat | 4 | 1% | 0.9738 | false |
| notmat | 4 | 1% | 2.6087 | false |
| mat | 4 | 100% | 0.9924 | false |
| notmat | 4 | 100% | 3.4813 | false |
| mat | 16 | 1% | 0.9531 | false |
| notmat | 16 | 1% | 11.3738 | false |
| mat | 16 | 100% | 1.0282 | false |
| notmat | 16 | 100% | 9.5289 | false |

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
| charity | yb-single:9100 |
| person | yb-single:9100 |
| donation | yb-single:9100 |
| donation_rank_idx | yb-single:9100 |
| donation_person_idx | yb-single:9100 |

Effective sharding flag (from the engine): `ysql_num_shards_per_tserver=2`. A relation with fewer tablets than this setting was created without an explicit split; the per-table `num_tablets` above is the authority.

## Throttling (cgroup `cpu.stat` deltas)

- client: nr_periods 301, nr_throttled 109, throttled_usec 3557702 | quota 100000 100000
, memory 2147483648

- yb-single: nr_periods 392, nr_throttled 79, throttled_usec 6376858 | quota 200000 100000, memory 3221225472

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
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 103.9393 | 131.4230 | 179.8862 | 5 | 7.6 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 1.9201 | 1.9701 | 2.3273 | 5 | 507.6 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 9.7021 | 10.9920 | 11.6903 | 5 | 91.0 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 101.2840 | 113.0238 | 147.9111 | 5 | 8.8 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 84.8263 | 86.7975 | 99.0724 | 5 | 11.5 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 84.8467 | 91.3876 | 99.3759 | 5 | 10.9 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 105.2178 | 108.0538 | 115.6056 | 5 | 9.3 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 1.7244 | 1.8743 | 2.2573 | 5 | 533.5 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 9.8302 | 10.2569 | 14.2086 | 5 | 97.5 |
| family=bag;stmt=join_where;threshold=180000 | valid | 108.4423 | 126.7814 | 134.5125 | 5 | 7.9 |
| family=bag;stmt=join_where;threshold=20 | valid | 1.6322 | 1.7054 | 2.0405 | 5 | 586.4 |
| family=bag;stmt=join_where;threshold=2000 | valid | 9.2544 | 9.7444 | 10.0624 | 5 | 102.6 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 9.1526 | 10.5321 | 10.8532 | 5 | 94.9 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 9.5693 | 9.9350 | 10.1967 | 5 | 100.7 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 6.4987 | 7.4711 | 9.0331 | 5 | 133.8 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 8.4338 | 9.0043 | 10.0911 | 5 | 111.1 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 8.4270 | 10.2502 | 11.2196 | 5 | 97.6 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 8.9971 | 9.3801 | 9.5736 | 5 | 106.6 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 9.2353 | 10.0384 | 10.3710 | 5 | 99.6 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 9.9611 | 10.2148 | 10.7792 | 5 | 97.9 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 7.5287 | 8.7035 | 9.1876 | 5 | 114.9 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 8.0859 | 8.6382 | 9.5669 | 5 | 115.8 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 8.5457 | 9.9818 | 12.3545 | 5 | 100.2 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 8.8813 | 9.3092 | 10.0974 | 5 | 107.4 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 105.8495 | 119.7908 | 158.7877 | 5 | 8.3 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 84.7913 | 94.6692 | 170.7610 | 5 | 10.6 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 6.6781 | 7.2312 | 7.6323 | 5 | 138.3 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 8.7932 | 9.5693 | 12.1604 | 5 | 104.5 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 23.7220 | 26.7395 | 30.7262 | 5 | 37.4 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 30.5458 | 32.6548 | 35.7305 | 5 | 30.6 |

## Matched blended workload

- phase: `isolated`; target offered rate: 10 updates/s
- offered 0, committed 0, failed 0 over 0.000 s
- offered rate 0.000/s, completed rate 0.000/s
- per key: no keys touched
- ledger: 0 keys reconciled, 0 mismatch

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42459358`, 5 interleaved rounds
  1. family=reuse;policy=notmat;refs=16;consumer=100%
  2. family=reuse;policy=notmat;refs=1;consumer=100%
  3. family=bag;stmt=cte_joined_outer;threshold=180000
  4. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  5. family=reuse;policy=default;refs=1;consumer=1%
  6. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  7. family=reuse;policy=mat;refs=4;consumer=100%
  8. family=reuse;policy=default;refs=16;consumer=1%
  9. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  10. family=bag;stmt=cte_joined_outer_mat;threshold=20
  11. family=bag;stmt=cte_joined_outer;threshold=20
  12. family=reuse;policy=notmat;refs=1;consumer=1%
  13. family=bag;stmt=join_where;threshold=20
  14. family=reuse;policy=notmat;refs=16;consumer=1%
  15. family=reuse;policy=default;refs=1;consumer=100%
  16. family=reuse;policy=default;refs=4;consumer=1%
  17. family=bag;stmt=cte_joined_outer;threshold=2000
  18. family=reuse;policy=mat;refs=4;consumer=1%
  19. family=bag;stmt=join_where;threshold=180000
  20. family=reuse;policy=default;refs=4;consumer=100%
  21. family=reuse;policy=mat;refs=1;consumer=1%
  22. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  23. family=reuse;policy=default;refs=16;consumer=100%
  24. family=bag;stmt=join_where;threshold=2000
  25. family=reuse;policy=mat;refs=16;consumer=100%
  26. family=reuse;policy=notmat;refs=4;consumer=1%
  27. family=reuse;policy=mat;refs=16;consumer=1%
  28. family=reuse;policy=notmat;refs=4;consumer=100%
  29. family=reuse;policy=mat;refs=1;consumer=100%
  30. family=bag;stmt=cte_joined_outer_notmat;threshold=180000

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
| family=bag;stmt=cte_joined_outer;threshold=180000 | j05_cte_joined_outer | valid | `7c6fbaa1cfb6f9ca` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20 | j05_cte_joined_outer | valid | `b14f0c61fd3b060b` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=2000 | j05_cte_joined_outer | valid | `99718efbcf134cf0` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | j05m_cte_joined_outer_mat | valid | `fc822e95e2cb10ac` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | j05m_cte_joined_outer_mat | valid | `e9b1b5232ab896ee` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | j05m_cte_joined_outer_mat | valid | `d625afa3108fe18a` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | j05n_cte_joined_outer_notmat | valid | `ec1dfe8a26afbf21` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | j05n_cte_joined_outer_notmat | valid | `f854f3757cd04282` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | j05n_cte_joined_outer_notmat | valid | `3357f1b06fcada84` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=join_where;threshold=180000 | j01_join_where | valid | `cce0d41e9575dd9c` |  |
| family=bag;stmt=join_where;threshold=20 | j01_join_where | valid | `df6123457081adfd` |  |
| family=bag;stmt=join_where;threshold=2000 | j01_join_where | valid | `d148f27184e1673c` |  |
| family=reuse;policy=default;refs=16;consumer=1% | r_default_n16_c1 | valid | `ad40fe9ce4de55bf` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=16;consumer=100% | r_default_n16_c100 | valid | `a32176d2253f9ef6` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=1;consumer=1% | r_default_n1_c1 | valid | `07c5b7241fdbac36` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=1;consumer=100% | r_default_n1_c100 | valid | `6cc19b7ffe37dd7b` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=4;consumer=1% | r_default_n4_c1 | valid | `835a6bb9bba0714d` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=100% | r_default_n4_c100 | valid | `1226eea54f78e5c4` | policy=default; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=1% | r_mat_n16_c1 | valid | `852458c995675e5e` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=100% | r_mat_n16_c100 | valid | `356b3985aa9da7be` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=1% | r_mat_n1_c1 | valid | `9e689d6102b289af` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=100% | r_mat_n1_c100 | valid | `84436f44d02b0ba8` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=1% | r_mat_n4_c1 | valid | `b2d7d291ebb7595c` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=100% | r_mat_n4_c100 | valid | `983a346bca42fa5c` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=notmat;refs=16;consumer=1% | r_notmat_n16_c1 | valid | `6e45e4443c7c8fdd` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `b60eaf55f8df19be` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `38801e01ecb87a6e` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `1cbbf26a53675d05` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `c7d85c7f9c9c6609` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `64a865cddd01176f` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, DIST, DEBUG, FORMAT TEXT` supported on yugabyte
- one host: this run used node budget 2 CPU 3g and client budget 1 CPU 2g, so budget comparisons are separate classes here, not independent hardware
- isolated phase: no concurrent writers in this run; the blended read/write phase is a separate labelled run
- physical rank-index arm yb-range-rank; an engine whose planner ignores the index is a different physical design, not a faster one
- YugabyteDB DIST/DEBUG storage counters (rows and index rows scanned) are parsed per cell; RPC-level counters are not aggregated
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
