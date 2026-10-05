# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 30 cells, 30 valid, engine `yugabyte`, topology `single`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-v4-topo-equal-i1-20261005`
- run tag: `run/study-01/select-join-v4-topo-equal-i1-20261005`
- environment: `host-zenbook-ux5406sa`
- repo commit: `ba821e7b686ac3adb1cda56dd45f15c8ae699ef8`
- engine version: `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit`
- database engine image: `docker.io/yugabytedb/yugabyte:2025.2.6.0-b111` @ `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`
- benchmark client image: `localhost/charitybench-select:select-join-v4-topo-equal-i1-20261005` @ `b6ea3c007cf459112eaa8e9e85381a7dc25a7710ff07648fde3048f747df6ab0`
- inputs digest: `64dcb81ee841e6d5e8e980097756dfe625516647b960dc1693ac2570685697d7`
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
| 20 (0.01%) | 1.6950 | missing | missing | missing | missing | missing | 1.9096 | 74.1659 | 1.9447 | missing |
| 2000 (1%) | 7.8345 | missing | missing | missing | missing | missing | 9.5134 | 77.2926 | 9.4829 | missing |
| 180000 (90%) | 99.4009 | missing | missing | missing | missing | missing | 96.1079 | 96.0912 | 97.3263 | missing |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | cte_joined_outer vs join_where | 1.6950 | 1.9096 | 1.1266 | false |
| 20 | cte_joined_outer_mat vs join_where | 1.6950 | 74.1659 | 43.7565 | false |
| 20 | cte_joined_outer_notmat vs join_where | 1.6950 | 1.9447 | 1.1473 | false |
| 2000 | cte_joined_outer vs join_where | 7.8345 | 9.5134 | 1.2143 | false |
| 2000 | cte_joined_outer_mat vs join_where | 7.8345 | 77.2926 | 9.8657 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 7.8345 | 9.4829 | 1.2104 | false |
| 180000 | cte_joined_outer vs join_where | 99.4009 | 96.1079 | 0.9669 | false |
| 180000 | cte_joined_outer_mat vs join_where | 99.4009 | 96.0912 | 0.9667 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 99.4009 | 97.3263 | 0.9791 | false |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_joined_outer_mat vs cte_joined_outer | 1.9096 | 74.1659 | 38.8386 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 1.9096 | 1.9447 | 1.0184 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 9.5134 | 77.2926 | 8.1246 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 9.5134 | 9.4829 | 0.9968 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 96.1079 | 96.0912 | 0.9998 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 96.1079 | 97.3263 | 1.0127 | false | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 5.9480 | 1 | `fe7c3013ba056c32` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 6.9311 | 1 | `b1879842aae341a2` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | missing | | | | |
| default | 2 | 100% | missing | | | | |
| default | 4 | 1% | 7.8210 | 4 | `c8ff11e72483fb42` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 8.0383 | 4 | `98875b983207fa11` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | missing | | | | |
| default | 8 | 100% | missing | | | | |
| default | 16 | 1% | 8.7344 | 16 | `f1d247ed70c554f8` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 9.3374 | 16 | `0d798719bf221d24` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 7.1422 | 1 | `26689b9a30f04643` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 7.4771 | 1 | `6121979ee37ba62a` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | missing | | | | |
| mat | 2 | 100% | missing | | | | |
| mat | 4 | 1% | 8.4432 | 4 | `bb6ff71ba77a76d5` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 8.4080 | 4 | `8230c7d6d590ba1e` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | missing | | | | |
| mat | 8 | 100% | missing | | | | |
| mat | 16 | 1% | 8.5172 | 16 | `686756c39c2621b1` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 9.1577 | 16 | `8ec50ae72bd2ae85` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 6.3419 | 1 | `de9bb2606c3b75d9` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 8.5447 | 1 | `77f75f759b564eb2` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | missing | | | | |
| notmat | 2 | 100% | missing | | | | |
| notmat | 4 | 1% | 20.5445 | 4 | `94fbc7dde5313dc4` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 25.8026 | 4 | `2a37e76abdc25203` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | missing | | | | |
| notmat | 8 | 100% | missing | | | | |
| notmat | 16 | 1% | 68.4239 | 16 | `ec549823bb3d0a92` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 77.3955 | 16 | `bacd7e0e7adb7b16` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 1.2008 | false |
| notmat | 1 | 1% | 1.0662 | false |
| mat | 1 | 100% | 1.0788 | false |
| notmat | 1 | 100% | 1.2328 | false |
| mat | 4 | 1% | 1.0796 | false |
| notmat | 4 | 1% | 2.6269 | false |
| mat | 4 | 100% | 1.0460 | false |
| notmat | 4 | 100% | 3.2099 | false |
| mat | 16 | 1% | 0.9751 | false |
| notmat | 16 | 1% | 7.8338 | false |
| mat | 16 | 100% | 0.9808 | false |
| notmat | 16 | 100% | 8.2887 | false |

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

- client: nr_periods 267, nr_throttled 102, throttled_usec 4339856 | quota 100000 100000
, memory 2147483648

- yb-single: nr_periods 335, nr_throttled 0, throttled_usec 0 | quota 600000 100000, memory 9663676416

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
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 93.6283 | 96.1079 | 102.2198 | 5 | 10.4 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 1.6199 | 1.9096 | 2.4780 | 5 | 523.7 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 9.2605 | 9.5134 | 12.7326 | 5 | 105.1 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 92.3197 | 96.0912 | 103.8036 | 5 | 10.4 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 70.6559 | 74.1659 | 96.3788 | 5 | 13.5 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 74.9540 | 77.2926 | 92.8054 | 5 | 12.9 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 95.4214 | 97.3263 | 111.6640 | 5 | 10.3 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 1.8510 | 1.9447 | 1.9923 | 5 | 514.2 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 9.0795 | 9.4829 | 9.5432 | 5 | 105.5 |
| family=bag;stmt=join_where;threshold=180000 | valid | 95.7502 | 99.4009 | 102.3971 | 5 | 10.1 |
| family=bag;stmt=join_where;threshold=20 | valid | 1.5977 | 1.6950 | 1.7720 | 5 | 590.0 |
| family=bag;stmt=join_where;threshold=2000 | valid | 7.5159 | 7.8345 | 8.3495 | 5 | 127.6 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 8.4778 | 8.7344 | 9.2346 | 5 | 114.5 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 8.3678 | 9.3374 | 10.2791 | 5 | 107.1 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 5.7529 | 5.9480 | 7.1979 | 5 | 168.1 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 6.8784 | 6.9311 | 8.1500 | 5 | 144.3 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 7.3484 | 7.8210 | 8.2507 | 5 | 127.9 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 7.8422 | 8.0383 | 9.9905 | 5 | 124.4 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 7.8408 | 8.5172 | 8.9852 | 5 | 117.4 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 8.2330 | 9.1577 | 11.0473 | 5 | 109.2 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 6.6656 | 7.1422 | 7.4659 | 5 | 140.0 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 7.1865 | 7.4771 | 8.7476 | 5 | 133.7 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 7.9871 | 8.4432 | 8.6449 | 5 | 118.4 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 8.1663 | 8.4080 | 8.8127 | 5 | 118.9 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 65.6390 | 68.4239 | 123.9784 | 5 | 14.6 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 76.2900 | 77.3955 | 93.8876 | 5 | 12.9 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 5.6453 | 6.3419 | 6.5772 | 5 | 157.7 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 8.4349 | 8.5447 | 8.8283 | 5 | 117.0 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 19.7035 | 20.5445 | 23.1099 | 5 | 48.7 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 24.6331 | 25.8026 | 32.9989 | 5 | 38.8 |

## Matched blended workload

- phase: `isolated`; target offered rate: 10 updates/s
- offered 0, committed 0, failed 0 over 0.000 s
- offered rate 0.000/s, completed rate 0.000/s
- per key: no keys touched
- ledger: 0 keys reconciled, 0 mismatch

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42040399`, 5 interleaved rounds
  1. family=reuse;policy=default;refs=4;consumer=100%
  2. family=reuse;policy=mat;refs=16;consumer=1%
  3. family=reuse;policy=default;refs=16;consumer=100%
  4. family=reuse;policy=mat;refs=16;consumer=100%
  5. family=reuse;policy=mat;refs=1;consumer=100%
  6. family=reuse;policy=notmat;refs=4;consumer=1%
  7. family=reuse;policy=default;refs=1;consumer=100%
  8. family=reuse;policy=default;refs=4;consumer=1%
  9. family=bag;stmt=join_where;threshold=2000
  10. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  11. family=bag;stmt=cte_joined_outer;threshold=2000
  12. family=reuse;policy=notmat;refs=1;consumer=1%
  13. family=reuse;policy=mat;refs=1;consumer=1%
  14. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  15. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  16. family=reuse;policy=notmat;refs=16;consumer=1%
  17. family=bag;stmt=join_where;threshold=180000
  18. family=reuse;policy=notmat;refs=1;consumer=100%
  19. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  20. family=reuse;policy=default;refs=16;consumer=1%
  21. family=reuse;policy=notmat;refs=16;consumer=100%
  22. family=reuse;policy=mat;refs=4;consumer=100%
  23. family=reuse;policy=default;refs=1;consumer=1%
  24. family=bag;stmt=cte_joined_outer;threshold=20
  25. family=reuse;policy=notmat;refs=4;consumer=100%
  26. family=bag;stmt=cte_joined_outer;threshold=180000
  27. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  28. family=reuse;policy=mat;refs=4;consumer=1%
  29. family=bag;stmt=cte_joined_outer_mat;threshold=20
  30. family=bag;stmt=join_where;threshold=20

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
| family=bag;stmt=cte_joined_outer;threshold=180000 | j05_cte_joined_outer | valid | `e564960ed45c3a41` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20 | j05_cte_joined_outer | valid | `494b47a9cc1dadc0` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=2000 | j05_cte_joined_outer | valid | `dc4b3202157b7e7b` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | j05m_cte_joined_outer_mat | valid | `7e14ed87cc64bb48` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | j05m_cte_joined_outer_mat | valid | `5a33017c3270910c` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | j05m_cte_joined_outer_mat | valid | `65886ee9b679cf04` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | j05n_cte_joined_outer_notmat | valid | `87f98441bc5737fe` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | j05n_cte_joined_outer_notmat | valid | `968e86682f6c1644` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | j05n_cte_joined_outer_notmat | valid | `4d0390c1b6023be4` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=join_where;threshold=180000 | j01_join_where | valid | `4487cf08796c981c` |  |
| family=bag;stmt=join_where;threshold=20 | j01_join_where | valid | `3d97955fdd5fd664` |  |
| family=bag;stmt=join_where;threshold=2000 | j01_join_where | valid | `baee24645942042a` |  |
| family=reuse;policy=default;refs=16;consumer=1% | r_default_n16_c1 | valid | `f1d247ed70c554f8` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=16;consumer=100% | r_default_n16_c100 | valid | `0d798719bf221d24` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=1;consumer=1% | r_default_n1_c1 | valid | `fe7c3013ba056c32` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=1;consumer=100% | r_default_n1_c100 | valid | `b1879842aae341a2` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=4;consumer=1% | r_default_n4_c1 | valid | `c8ff11e72483fb42` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=100% | r_default_n4_c100 | valid | `98875b983207fa11` | policy=default; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=1% | r_mat_n16_c1 | valid | `686756c39c2621b1` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=100% | r_mat_n16_c100 | valid | `8ec50ae72bd2ae85` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=1% | r_mat_n1_c1 | valid | `26689b9a30f04643` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=100% | r_mat_n1_c100 | valid | `6121979ee37ba62a` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=1% | r_mat_n4_c1 | valid | `bb6ff71ba77a76d5` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=100% | r_mat_n4_c100 | valid | `8230c7d6d590ba1e` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=notmat;refs=16;consumer=1% | r_notmat_n16_c1 | valid | `ec549823bb3d0a92` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `bacd7e0e7adb7b16` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `de9bb2606c3b75d9` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `77f75f759b564eb2` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `94fbc7dde5313dc4` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `2a37e76abdc25203` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, DIST, DEBUG, FORMAT TEXT` supported on yugabyte
- one host: this run used node budget 6 CPU 9g and client budget 1 CPU 2g, so budget comparisons are separate classes here, not independent hardware
- isolated phase: no concurrent writers in this run; the blended read/write phase is a separate labelled run
- physical rank-index arm yb-range-rank; an engine whose planner ignores the index is a different physical design, not a faster one
- YugabyteDB DIST/DEBUG storage counters (rows and index rows scanned) are parsed per cell; RPC-level counters are not aggregated
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
