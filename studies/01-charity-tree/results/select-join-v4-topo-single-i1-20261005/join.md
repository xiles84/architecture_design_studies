# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 30 cells, 30 valid, engine `yugabyte`, topology `single`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-v4-topo-single-i1-20261005`
- run tag: `run/study-01/select-join-v4-topo-single-i1-20261005`
- environment: `host-zenbook-ux5406sa`
- repo commit: `ba821e7b686ac3adb1cda56dd45f15c8ae699ef8`
- engine version: `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit`
- database engine image: `docker.io/yugabytedb/yugabyte:2025.2.6.0-b111` @ `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`
- benchmark client image: `localhost/charitybench-select:select-join-v4-topo-single-i1-20261005` @ `b6ea3c007cf459112eaa8e9e85381a7dc25a7710ff07648fde3048f747df6ab0`
- inputs digest: `5ff0161edb2e32cb4c63da0123451f452b614e6c9ef4330c1c4724292cd9c92a`
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
| 20 (0.01%) | 2.2348 | missing | missing | missing | missing | missing | 2.0433 | 88.7571 | 2.0240 | missing |
| 2000 (1%) | 9.4894 | missing | missing | missing | missing | missing | 9.3226 | 89.0481 | 9.7642 | missing |
| 180000 (90%) | 108.8445 | missing | missing | missing | missing | missing | 118.7366 | 109.2512 | 121.2623 | missing |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | cte_joined_outer vs join_where | 2.2348 | 2.0433 | 0.9143 | false |
| 20 | cte_joined_outer_mat vs join_where | 2.2348 | 88.7571 | 39.7161 | false |
| 20 | cte_joined_outer_notmat vs join_where | 2.2348 | 2.0240 | 0.9057 | false |
| 2000 | cte_joined_outer vs join_where | 9.4894 | 9.3226 | 0.9824 | false |
| 2000 | cte_joined_outer_mat vs join_where | 9.4894 | 89.0481 | 9.3840 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 9.4894 | 9.7642 | 1.0290 | false |
| 180000 | cte_joined_outer vs join_where | 108.8445 | 118.7366 | 1.0909 | false |
| 180000 | cte_joined_outer_mat vs join_where | 108.8445 | 109.2512 | 1.0037 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 108.8445 | 121.2623 | 1.1141 | false |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_joined_outer_mat vs cte_joined_outer | 2.0433 | 88.7571 | 43.4371 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 2.0433 | 2.0240 | 0.9905 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 9.3226 | 89.0481 | 9.5519 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 9.3226 | 9.7642 | 1.0474 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 118.7366 | 109.2512 | 0.9201 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 118.7366 | 121.2623 | 1.0213 | false | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 7.0662 | 1 | `ce11dd712859024d` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 8.7268 | 1 | `99d207b163d67945` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | missing | | | | |
| default | 2 | 100% | missing | | | | |
| default | 4 | 1% | 8.6680 | 4 | `41c113d6bb795e65` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 8.8600 | 4 | `f34e257b77aba744` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | missing | | | | |
| default | 8 | 100% | missing | | | | |
| default | 16 | 1% | 10.0155 | 16 | `03f247f58da5647b` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 9.5782 | 16 | `935df13832e30cc7` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 8.4473 | 1 | `f079e3f770007f7b` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 8.8105 | 1 | `991020c6ee391f67` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | missing | | | | |
| mat | 2 | 100% | missing | | | | |
| mat | 4 | 1% | 8.6703 | 4 | `5eb3d880db012937` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 8.9971 | 4 | `bdd14bb04c0cdb4e` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | missing | | | | |
| mat | 8 | 100% | missing | | | | |
| mat | 16 | 1% | 9.1091 | 16 | `5b99d37ffdf83586` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 9.7011 | 16 | `db7479af6f1dc204` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 7.0045 | 1 | `2c3cbcf8c38ced0c` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 9.2390 | 1 | `8e8f588d48497ed4` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | missing | | | | |
| notmat | 2 | 100% | missing | | | | |
| notmat | 4 | 1% | 25.8475 | 4 | `1481e83aa8258fe6` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 31.5749 | 4 | `e815a6e3cbb0de08` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | missing | | | | |
| notmat | 8 | 100% | missing | | | | |
| notmat | 16 | 1% | 98.9009 | 16 | `9ebed617f1d65633` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 87.3387 | 16 | `2c47ffd21a7cf853` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 1.1954 | false |
| notmat | 1 | 1% | 0.9913 | false |
| mat | 1 | 100% | 1.0096 | false |
| notmat | 1 | 100% | 1.0587 | false |
| mat | 4 | 1% | 1.0003 | false |
| notmat | 4 | 1% | 2.9819 | false |
| mat | 4 | 100% | 1.0155 | false |
| notmat | 4 | 100% | 3.5638 | false |
| mat | 16 | 1% | 0.9095 | false |
| notmat | 16 | 1% | 9.8748 | false |
| mat | 16 | 100% | 1.0128 | false |
| notmat | 16 | 100% | 9.1185 | false |

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

Effective sharding flag (from the engine): `ysql_num_shards_per_tserver=2`. A relation with fewer tablets than this setting was created without an explicit split; the per-table `num_tablets` above is the authority.

## Throttling (cgroup `cpu.stat` deltas)

- client: nr_periods 298, nr_throttled 110, throttled_usec 3733757 | quota 100000 100000
, memory 2147483648

- yb-single: nr_periods 380, nr_throttled 82, throttled_usec 5214799 | quota 200000 100000, memory 3221225472

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
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 108.4198 | 118.7366 | 125.1355 | 5 | 8.4 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 1.7107 | 2.0433 | 2.3401 | 5 | 489.4 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 8.8545 | 9.3226 | 10.6467 | 5 | 107.3 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 99.0555 | 109.2512 | 150.5467 | 5 | 9.2 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 82.7491 | 88.7571 | 93.1349 | 5 | 11.3 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 87.0104 | 89.0481 | 94.7062 | 5 | 11.2 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 104.9820 | 121.2623 | 127.7793 | 5 | 8.2 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 1.7212 | 2.0240 | 2.5924 | 5 | 494.1 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 9.1424 | 9.7642 | 10.6138 | 5 | 102.4 |
| family=bag;stmt=join_where;threshold=180000 | valid | 106.9851 | 108.8445 | 110.3997 | 5 | 9.2 |
| family=bag;stmt=join_where;threshold=20 | valid | 1.7891 | 2.2348 | 2.4243 | 5 | 447.5 |
| family=bag;stmt=join_where;threshold=2000 | valid | 9.0042 | 9.4894 | 10.3569 | 5 | 105.4 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 9.5777 | 10.0155 | 10.5073 | 5 | 99.8 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 9.5145 | 9.5782 | 9.7812 | 5 | 104.4 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 6.7211 | 7.0662 | 7.3313 | 5 | 141.5 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 8.2755 | 8.7268 | 9.6705 | 5 | 114.6 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 8.4225 | 8.6680 | 10.8212 | 5 | 115.4 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 8.2961 | 8.8600 | 9.3584 | 5 | 112.9 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 9.0454 | 9.1091 | 9.8949 | 5 | 109.8 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 8.9741 | 9.7011 | 10.0289 | 5 | 103.1 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 8.1502 | 8.4473 | 9.2465 | 5 | 118.4 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 7.9313 | 8.8105 | 9.3917 | 5 | 113.5 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 8.5307 | 8.6703 | 11.1559 | 5 | 115.3 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 8.5771 | 8.9971 | 9.4764 | 5 | 111.1 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 95.6900 | 98.9009 | 140.1751 | 5 | 10.1 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 82.7070 | 87.3387 | 88.2464 | 5 | 11.4 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 6.7547 | 7.0045 | 7.7668 | 5 | 142.8 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 8.6683 | 9.2390 | 9.6108 | 5 | 108.2 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 24.4278 | 25.8475 | 26.6775 | 5 | 38.7 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 30.9437 | 31.5749 | 32.4271 | 5 | 31.7 |

## Matched blended workload

- phase: `isolated`; target offered rate: 10 updates/s
- offered 0, committed 0, failed 0 over 0.000 s
- offered rate 0.000/s, completed rate 0.000/s
- per key: no keys touched
- ledger: 0 keys reconciled, 0 mismatch

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42418261`, 5 interleaved rounds
  1. family=reuse;policy=default;refs=4;consumer=100%
  2. family=bag;stmt=cte_joined_outer;threshold=20
  3. family=reuse;policy=mat;refs=16;consumer=1%
  4. family=bag;stmt=join_where;threshold=180000
  5. family=reuse;policy=notmat;refs=16;consumer=100%
  6. family=reuse;policy=mat;refs=1;consumer=1%
  7. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  8. family=reuse;policy=notmat;refs=16;consumer=1%
  9. family=reuse;policy=default;refs=1;consumer=100%
  10. family=bag;stmt=cte_joined_outer;threshold=180000
  11. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  12. family=reuse;policy=mat;refs=1;consumer=100%
  13. family=reuse;policy=default;refs=16;consumer=100%
  14. family=reuse;policy=notmat;refs=4;consumer=1%
  15. family=reuse;policy=mat;refs=16;consumer=100%
  16. family=reuse;policy=default;refs=1;consumer=1%
  17. family=bag;stmt=join_where;threshold=20
  18. family=reuse;policy=mat;refs=4;consumer=100%
  19. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  20. family=bag;stmt=join_where;threshold=2000
  21. family=reuse;policy=default;refs=4;consumer=1%
  22. family=reuse;policy=mat;refs=4;consumer=1%
  23. family=bag;stmt=cte_joined_outer;threshold=2000
  24. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  25. family=reuse;policy=notmat;refs=4;consumer=100%
  26. family=reuse;policy=notmat;refs=1;consumer=1%
  27. family=bag;stmt=cte_joined_outer_mat;threshold=20
  28. family=reuse;policy=notmat;refs=1;consumer=100%
  29. family=bag;stmt=cte_joined_outer_mat;threshold=180000
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
| family=bag;stmt=cte_joined_outer;threshold=180000 | j05_cte_joined_outer | valid | `77906e27a921428c` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20 | j05_cte_joined_outer | valid | `166a1f5663ad0a98` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=2000 | j05_cte_joined_outer | valid | `c8115c42a4ee431c` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | j05m_cte_joined_outer_mat | valid | `b2e97e8085db775e` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | j05m_cte_joined_outer_mat | valid | `4f2e43d9ed6998ff` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | j05m_cte_joined_outer_mat | valid | `933a828128d0f44f` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | j05n_cte_joined_outer_notmat | valid | `f72aab5f77c90854` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | j05n_cte_joined_outer_notmat | valid | `f12fefe75240a126` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | j05n_cte_joined_outer_notmat | valid | `9e375f1495d7b298` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=join_where;threshold=180000 | j01_join_where | valid | `555433fcd68e3a02` |  |
| family=bag;stmt=join_where;threshold=20 | j01_join_where | valid | `05ead2fe0f4ad1a8` |  |
| family=bag;stmt=join_where;threshold=2000 | j01_join_where | valid | `e7d692f75c6f13f3` |  |
| family=reuse;policy=default;refs=16;consumer=1% | r_default_n16_c1 | valid | `03f247f58da5647b` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=16;consumer=100% | r_default_n16_c100 | valid | `935df13832e30cc7` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=1;consumer=1% | r_default_n1_c1 | valid | `ce11dd712859024d` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=1;consumer=100% | r_default_n1_c100 | valid | `99d207b163d67945` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=4;consumer=1% | r_default_n4_c1 | valid | `41c113d6bb795e65` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=100% | r_default_n4_c100 | valid | `f34e257b77aba744` | policy=default; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=1% | r_mat_n16_c1 | valid | `5b99d37ffdf83586` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=100% | r_mat_n16_c100 | valid | `db7479af6f1dc204` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=1% | r_mat_n1_c1 | valid | `f079e3f770007f7b` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=100% | r_mat_n1_c100 | valid | `991020c6ee391f67` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=1% | r_mat_n4_c1 | valid | `5eb3d880db012937` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=100% | r_mat_n4_c100 | valid | `bdd14bb04c0cdb4e` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=notmat;refs=16;consumer=1% | r_notmat_n16_c1 | valid | `9ebed617f1d65633` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `2c47ffd21a7cf853` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `2c3cbcf8c38ced0c` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `8e8f588d48497ed4` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `1481e83aa8258fe6` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `e815a6e3cbb0de08` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, DIST, DEBUG, FORMAT TEXT` supported on yugabyte
- one host: this run used node budget 2 CPU 3g and client budget 1 CPU 2g, so budget comparisons are separate classes here, not independent hardware
- isolated phase: no concurrent writers in this run; the blended read/write phase is a separate labelled run
- physical rank-index arm yb-range-rank; an engine whose planner ignores the index is a different physical design, not a faster one
- YugabyteDB DIST/DEBUG storage counters (rows and index rows scanned) are parsed per cell; RPC-level counters are not aggregated
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
