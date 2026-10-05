# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 30 cells, 30 valid, engine `yugabyte`, topology `cluster3`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-v4-topo-tablet-i3-20261005`
- run tag: `run/study-01/select-join-v4-topo-tablet-i3-20261005`
- environment: `host-zenbook-ux5406sa`
- repo commit: `ba821e7b686ac3adb1cda56dd45f15c8ae699ef8`
- engine version: `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit`
- database engine image: `docker.io/yugabytedb/yugabyte:2025.2.6.0-b111` @ `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`
- benchmark client image: `localhost/charitybench-select:select-join-v4-topo-tablet-i3-20261005` @ `b6ea3c007cf459112eaa8e9e85381a7dc25a7710ff07648fde3048f747df6ab0`
- inputs digest: `34447d84de019b9cadd9dc7faf9ad95497d2c6212ce3eb4bfca4f79b5d799b9d`
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
| j01_join_where | valid | true |  |
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
| 20 (0.01%) | 2.6144 | missing | missing | missing | missing | missing | 2.7767 | 116.8543 | 1.9744 | missing |
| 2000 (1%) | 11.5129 | missing | missing | missing | missing | missing | 10.3747 | 124.2378 | 11.5250 | missing |
| 180000 (90%) | 107.7438 | missing | missing | missing | missing | missing | 143.3530 | 137.7133 | 141.4773 | missing |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | cte_joined_outer vs join_where | 2.6144 | 2.7767 | 1.0620 | false |
| 20 | cte_joined_outer_mat vs join_where | 2.6144 | 116.8543 | 44.6956 | false |
| 20 | cte_joined_outer_notmat vs join_where | 2.6144 | 1.9744 | 0.7552 | false |
| 2000 | cte_joined_outer vs join_where | 11.5129 | 10.3747 | 0.9011 | false |
| 2000 | cte_joined_outer_mat vs join_where | 11.5129 | 124.2378 | 10.7912 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 11.5129 | 11.5250 | 1.0011 | false |
| 180000 | cte_joined_outer vs join_where | 107.7438 | 143.3530 | 1.3305 | false |
| 180000 | cte_joined_outer_mat vs join_where | 107.7438 | 137.7133 | 1.2782 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 107.7438 | 141.4773 | 1.3131 | false |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_joined_outer_mat vs cte_joined_outer | 2.7767 | 116.8543 | 42.0845 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 2.7767 | 1.9744 | 0.7111 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 10.3747 | 124.2378 | 11.9750 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 10.3747 | 11.5250 | 1.1109 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 143.3530 | 137.7133 | 0.9607 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 143.3530 | 141.4773 | 0.9869 | false | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 7.7494 | 1 | `f7abe27cf7445c31` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 8.6454 | 1 | `995325a173d66144` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | missing | | | | |
| default | 2 | 100% | missing | | | | |
| default | 4 | 1% | 9.4074 | 4 | `0d1a5541b087fbde` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 9.7193 | 4 | `faa819fecb2b63f7` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | missing | | | | |
| default | 8 | 100% | missing | | | | |
| default | 16 | 1% | 9.4064 | 16 | `15defcb15ba2d19d` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 10.4564 | 16 | `2b89f17962a95c75` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 9.0683 | 1 | `f7350e5758e4c62d` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 8.1893 | 1 | `a6d6979bfd2f7db3` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | missing | | | | |
| mat | 2 | 100% | missing | | | | |
| mat | 4 | 1% | 8.3419 | 4 | `19d6e13c655bacda` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 9.7356 | 4 | `5acf99ce8aaac5ee` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | missing | | | | |
| mat | 8 | 100% | missing | | | | |
| mat | 16 | 1% | 9.3783 | 16 | `35d9c39722b750ad` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 9.7549 | 16 | `1802197ecb87d06a` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 7.4002 | 1 | `734a73d864f0b127` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 9.6020 | 1 | `25eb236d6cb31470` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | missing | | | | |
| notmat | 2 | 100% | missing | | | | |
| notmat | 4 | 1% | 22.0028 | 4 | `45de12fcc114fa0c` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 33.0502 | 4 | `124c7a18171c7b0e` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | missing | | | | |
| notmat | 8 | 100% | missing | | | | |
| notmat | 16 | 1% | 77.9223 | 16 | `a9070ae76a60c265` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 113.6569 | 16 | `6d7c569b5e35a402` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 1.1702 | false |
| notmat | 1 | 1% | 0.9549 | false |
| mat | 1 | 100% | 0.9472 | false |
| notmat | 1 | 100% | 1.1106 | false |
| mat | 4 | 1% | 0.8867 | false |
| notmat | 4 | 1% | 2.3389 | false |
| mat | 4 | 100% | 1.0017 | false |
| notmat | 4 | 100% | 3.4005 | false |
| mat | 16 | 1% | 0.9970 | false |
| notmat | 16 | 1% | 8.2840 | false |
| mat | 16 | 100% | 0.9329 | false |
| notmat | 16 | 100% | 10.8696 | false |

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

- client: nr_periods 311, nr_throttled 115, throttled_usec 4670066 | quota 100000 100000
, memory 2147483648

- yb-n1: nr_periods 599, nr_throttled 4, throttled_usec 521500 | quota 200000 100000, memory 3221225472
- yb-n2: nr_periods 595, nr_throttled 17, throttled_usec 1672194 | quota 200000 100000, memory 3221225472
- yb-n3: nr_periods 588, nr_throttled 0, throttled_usec 0 | quota 200000 100000, memory 3221225472

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
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 141.7488 | 143.3530 | 151.9643 | 5 | 7.0 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 2.7181 | 2.7767 | 3.3864 | 5 | 360.1 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 9.8909 | 10.3747 | 10.6347 | 5 | 96.4 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 131.2495 | 137.7133 | 139.4444 | 5 | 7.3 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 113.3241 | 116.8543 | 118.6104 | 5 | 8.6 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 113.2593 | 124.2378 | 141.6260 | 5 | 8.0 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 134.1992 | 141.4773 | 145.3934 | 5 | 7.1 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 1.8706 | 1.9744 | 2.6902 | 5 | 506.5 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 10.8069 | 11.5250 | 11.8259 | 5 | 86.8 |
| family=bag;stmt=join_where;threshold=180000 | valid | 103.5134 | 107.7438 | 112.8000 | 5 | 9.3 |
| family=bag;stmt=join_where;threshold=20 | valid | 2.5002 | 2.6144 | 2.6919 | 5 | 382.5 |
| family=bag;stmt=join_where;threshold=2000 | valid | 11.1457 | 11.5129 | 12.4074 | 5 | 86.9 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 9.0777 | 9.4064 | 9.4720 | 5 | 106.3 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 10.3963 | 10.4564 | 10.6721 | 5 | 95.6 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 7.1014 | 7.7494 | 11.7202 | 5 | 129.0 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 8.4469 | 8.6454 | 9.8045 | 5 | 115.7 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 8.8928 | 9.4074 | 10.5469 | 5 | 106.3 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 9.3381 | 9.7193 | 9.8817 | 5 | 102.9 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 9.1695 | 9.3783 | 10.2220 | 5 | 106.6 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 9.4100 | 9.7549 | 10.3359 | 5 | 102.5 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 8.6157 | 9.0683 | 10.3396 | 5 | 110.3 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 7.9515 | 8.1893 | 8.2815 | 5 | 122.1 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 8.0548 | 8.3419 | 9.0895 | 5 | 119.9 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 9.4977 | 9.7356 | 10.3449 | 5 | 102.7 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 76.6052 | 77.9223 | 104.7403 | 5 | 12.8 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 80.9219 | 113.6569 | 124.6792 | 5 | 8.8 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 7.1918 | 7.4002 | 7.7883 | 5 | 135.1 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 9.0832 | 9.6020 | 11.0542 | 5 | 104.1 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 21.4963 | 22.0028 | 23.5023 | 5 | 45.4 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 32.3524 | 33.0502 | 33.6880 | 5 | 30.3 |

## Matched blended workload

- phase: `isolated`; target offered rate: 10 updates/s
- offered 0, committed 0, failed 0 over 0.000 s
- offered rate 0.000/s, completed rate 0.000/s
- per key: no keys touched
- ledger: 0 keys reconciled, 0 mismatch

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42237843`, 5 interleaved rounds
  1. family=reuse;policy=notmat;refs=4;consumer=100%
  2. family=reuse;policy=mat;refs=1;consumer=100%
  3. family=bag;stmt=cte_joined_outer_mat;threshold=20
  4. family=bag;stmt=join_where;threshold=20
  5. family=reuse;policy=mat;refs=16;consumer=1%
  6. family=reuse;policy=default;refs=1;consumer=1%
  7. family=reuse;policy=default;refs=4;consumer=1%
  8. family=bag;stmt=join_where;threshold=180000
  9. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  10. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  11. family=reuse;policy=default;refs=16;consumer=1%
  12. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  13. family=bag;stmt=cte_joined_outer;threshold=20
  14. family=reuse;policy=notmat;refs=16;consumer=1%
  15. family=bag;stmt=cte_joined_outer;threshold=180000
  16. family=bag;stmt=join_where;threshold=2000
  17. family=reuse;policy=notmat;refs=16;consumer=100%
  18. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  19. family=reuse;policy=notmat;refs=1;consumer=100%
  20. family=reuse;policy=mat;refs=16;consumer=100%
  21. family=reuse;policy=notmat;refs=1;consumer=1%
  22. family=reuse;policy=default;refs=16;consumer=100%
  23. family=reuse;policy=mat;refs=4;consumer=1%
  24. family=bag;stmt=cte_joined_outer;threshold=2000
  25. family=reuse;policy=default;refs=4;consumer=100%
  26. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  27. family=reuse;policy=mat;refs=4;consumer=100%
  28. family=reuse;policy=mat;refs=1;consumer=1%
  29. family=reuse;policy=notmat;refs=4;consumer=1%
  30. family=reuse;policy=default;refs=1;consumer=100%

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
| family=bag;stmt=cte_joined_outer;threshold=180000 | j05_cte_joined_outer | valid | `92c20a3c1df684a2` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20 | j05_cte_joined_outer | valid | `ed20085e0a1cccc0` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=2000 | j05_cte_joined_outer | valid | `b7ef65e0f4f6b483` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | j05m_cte_joined_outer_mat | valid | `6df097a949fced5a` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | j05m_cte_joined_outer_mat | valid | `69835cfb3305646a` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | j05m_cte_joined_outer_mat | valid | `06ae9586e790dd40` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | j05n_cte_joined_outer_notmat | valid | `537d7c4e9fb71a87` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | j05n_cte_joined_outer_notmat | valid | `c19335cb1ba7a4f5` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | j05n_cte_joined_outer_notmat | valid | `18d88c01695cf3c5` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=join_where;threshold=180000 | j01_join_where | valid | `56e43e45c18f28c6` |  |
| family=bag;stmt=join_where;threshold=20 | j01_join_where | valid | `cbd9bc8f93c7ec80` |  |
| family=bag;stmt=join_where;threshold=2000 | j01_join_where | valid | `b68ba13511db2118` |  |
| family=reuse;policy=default;refs=16;consumer=1% | r_default_n16_c1 | valid | `15defcb15ba2d19d` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=16;consumer=100% | r_default_n16_c100 | valid | `2b89f17962a95c75` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=1;consumer=1% | r_default_n1_c1 | valid | `f7abe27cf7445c31` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=1;consumer=100% | r_default_n1_c100 | valid | `995325a173d66144` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=4;consumer=1% | r_default_n4_c1 | valid | `0d1a5541b087fbde` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=100% | r_default_n4_c100 | valid | `faa819fecb2b63f7` | policy=default; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=1% | r_mat_n16_c1 | valid | `35d9c39722b750ad` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=100% | r_mat_n16_c100 | valid | `1802197ecb87d06a` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=1% | r_mat_n1_c1 | valid | `f7350e5758e4c62d` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=100% | r_mat_n1_c100 | valid | `a6d6979bfd2f7db3` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=1% | r_mat_n4_c1 | valid | `19d6e13c655bacda` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=100% | r_mat_n4_c100 | valid | `5acf99ce8aaac5ee` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=notmat;refs=16;consumer=1% | r_notmat_n16_c1 | valid | `a9070ae76a60c265` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `6d7c569b5e35a402` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `734a73d864f0b127` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `25eb236d6cb31470` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `45de12fcc114fa0c` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `124c7a18171c7b0e` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, DIST, DEBUG, FORMAT TEXT` supported on yugabyte
- one host: this run used node budget 2 CPU 3g and client budget 1 CPU 2g, so budget comparisons are separate classes here, not independent hardware
- isolated phase: no concurrent writers in this run; the blended read/write phase is a separate labelled run
- cluster nodes are containers on the same host, so these are not network or independent-hardware measurements; container colocation is not data-colocation evidence
- readers round-robin across 3 endpoints; the writer's endpoint policy is recorded separately
- physical rank-index arm yb-range-rank; an engine whose planner ignores the index is a different physical design, not a faster one
- YugabyteDB DIST/DEBUG storage counters (rows and index rows scanned) are parsed per cell; RPC-level counters are not aggregated
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
