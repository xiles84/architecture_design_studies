# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 30 cells, 30 valid, engine `yugabyte`, topology `single`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-v3b-topo-single-i2-20261005`
- run tag: `run/study-01/select-join-v3b-topo-single-i2-20261005`
- environment: `host-zenbook-ux5406sa`
- repo commit: `7243fee074bfb16fc5ad8b6afc1a2d7549080935`
- engine version: `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit`
- database engine image: `docker.io/yugabytedb/yugabyte:2025.2.6.0-b111` @ `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`
- benchmark client image: `localhost/charitybench-select:select-join-v3b-topo-single-i2-20261005` @ `f3f1e10c44cbcb4c7b5bc56ed95a4f071682b6f7d33b5e0a7b688c4271eeef88`
- inputs digest: `098b89d9c83de0dd890731025244234b10d1cb250c23118065fdf0f8878c11c8`
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
- setting `placement`: noncolocated
- setting `prepared_plan_policy`: pgx simple protocol; no prepared statements
- setting `scale`: cte-domain
- setting `selectivity_mapping`: rank <= threshold; rank is a deterministic shuffle of 1..fact_count
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
| 20 (0.01%) | 1.9586 | missing | missing | missing | missing | missing | 2.4953 | 77.8521 | 1.9253 | missing |
| 2000 (1%) | 8.5163 | missing | missing | missing | missing | missing | 8.2152 | 77.5481 | 9.0017 | missing |
| 180000 (90%) | 101.0970 | missing | missing | missing | missing | missing | 102.6497 | 97.6008 | 101.2703 | missing |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | cte_joined_outer vs join_where | 1.9586 | 2.4953 | 1.2740 | false |
| 20 | cte_joined_outer_mat vs join_where | 1.9586 | 77.8521 | 39.7486 | false |
| 20 | cte_joined_outer_notmat vs join_where | 1.9586 | 1.9253 | 0.9830 | false |
| 2000 | cte_joined_outer vs join_where | 8.5163 | 8.2152 | 0.9646 | false |
| 2000 | cte_joined_outer_mat vs join_where | 8.5163 | 77.5481 | 9.1059 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 8.5163 | 9.0017 | 1.0570 | false |
| 180000 | cte_joined_outer vs join_where | 101.0970 | 102.6497 | 1.0154 | false |
| 180000 | cte_joined_outer_mat vs join_where | 101.0970 | 97.6008 | 0.9654 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 101.0970 | 101.2703 | 1.0017 | false |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_joined_outer_mat vs cte_joined_outer | 2.4953 | 77.8521 | 31.1994 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 2.4953 | 1.9253 | 0.7716 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 8.2152 | 77.5481 | 9.4396 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 8.2152 | 9.0017 | 1.0957 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 102.6497 | 97.6008 | 0.9508 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 102.6497 | 101.2703 | 0.9866 | false | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 6.6818 | 1 | `10b47fef96fd6824` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 7.5329 | 1 | `4b0fc303a68a0b49` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | missing | | | | |
| default | 2 | 100% | missing | | | | |
| default | 4 | 1% | 7.9716 | 4 | `e77f07d5e6a4503a` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 8.4075 | 4 | `c3a41a9f40fa636c` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | missing | | | | |
| default | 8 | 100% | missing | | | | |
| default | 16 | 1% | 8.5079 | 16 | `657067656d8263c4` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 9.1553 | 16 | `5c636c5399e736d8` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 8.3475 | 1 | `ab078b5079848a27` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 8.5228 | 1 | `6241948476d1ed00` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | missing | | | | |
| mat | 2 | 100% | missing | | | | |
| mat | 4 | 1% | 8.8278 | 4 | `c64ac35e89850567` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 7.9760 | 4 | `f0ad8022db1b80b6` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | missing | | | | |
| mat | 8 | 100% | missing | | | | |
| mat | 16 | 1% | 8.7898 | 16 | `4dbe2dfcae740eae` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 9.3698 | 16 | `aca3bc92020de07c` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 6.4928 | 1 | `f66483d8cbc542be` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 7.4394 | 1 | `d42081e770e94ef2` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | missing | | | | |
| notmat | 2 | 100% | missing | | | | |
| notmat | 4 | 1% | 22.1093 | 4 | `9c214b41df0e91bd` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 30.1732 | 4 | `2df3fe6c777bc7f1` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | missing | | | | |
| notmat | 8 | 100% | missing | | | | |
| notmat | 16 | 1% | 99.2313 | 16 | `69ba42f9bfcbaf19` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 111.5181 | 16 | `39c538df37561609` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 1.2493 | false |
| notmat | 1 | 1% | 0.9717 | false |
| mat | 1 | 100% | 1.1314 | false |
| notmat | 1 | 100% | 0.9876 | false |
| mat | 4 | 1% | 1.1074 | false |
| notmat | 4 | 1% | 2.7735 | false |
| mat | 4 | 100% | 0.9487 | false |
| notmat | 4 | 100% | 3.5888 | false |
| mat | 16 | 1% | 1.0331 | false |
| notmat | 16 | 1% | 11.6634 | false |
| mat | 16 | 100% | 1.0234 | false |
| notmat | 16 | 100% | 12.1807 | false |

## Physical arm

- `charity_pkey`: `CREATE UNIQUE INDEX charity_pkey ON selectjoin_domain.charity USING lsm (charity_id HASH)`
- `donation_person_idx`: `CREATE INDEX donation_person_idx ON selectjoin_domain.donation USING lsm (person_id HASH)`
- `donation_pkey`: `CREATE UNIQUE INDEX donation_pkey ON selectjoin_domain.donation USING lsm (donation_id HASH)`
- `donation_rank_idx`: `CREATE INDEX donation_rank_idx ON selectjoin_domain.donation USING lsm (rank ASC)`
- `person_pkey`: `CREATE UNIQUE INDEX person_pkey ON selectjoin_domain.person USING lsm (person_id HASH)`

- index arm `yb-range-rank`; placement `noncolocated`; node budget `2 CPU 3g`; client budget `1 CPU 2g`

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

## Throttling (cgroup `cpu.stat` deltas)

- client: nr_periods 280, nr_throttled 107, throttled_usec 4808165 | quota 100000 100000
, memory 2147483648

- yb-single: nr_periods 355, nr_throttled 80, throttled_usec 7407562 | quota 200000 100000, memory 3221225472

## Access path per cell

_Engine-reported from each cell's own plan: whether the rank index or a sequential scan was used, rows scanned, PostgreSQL buffers and filter removals, and reader retries of transient conflicts. A cell that scanned the whole table at a selective threshold is a different physical regime, not a faster query._

| axis | status | index scan | seq scan | storage rows scanned | index rows scanned | buffers hit/read | rows removed | retries |
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
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 96.9313 | 102.6497 | 104.6045 | 5 | 9.7 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 2.4112 | 2.4953 | 2.6085 | 5 | 400.8 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 8.1396 | 8.2152 | 9.3173 | 5 | 121.7 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 91.6873 | 97.6008 | 134.7627 | 5 | 10.2 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 76.2548 | 77.8521 | 80.3784 | 5 | 12.8 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 77.1998 | 77.5481 | 83.5538 | 5 | 12.9 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 98.0178 | 101.2703 | 104.6463 | 5 | 9.9 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 1.7993 | 1.9253 | 2.0694 | 5 | 519.4 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 8.1222 | 9.0017 | 10.2012 | 5 | 111.1 |
| family=bag;stmt=join_where;threshold=180000 | valid | 100.7052 | 101.0970 | 106.8338 | 5 | 9.9 |
| family=bag;stmt=join_where;threshold=20 | valid | 1.8709 | 1.9586 | 2.1050 | 5 | 510.6 |
| family=bag;stmt=join_where;threshold=2000 | valid | 8.3443 | 8.5163 | 8.9025 | 5 | 117.4 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 8.2838 | 8.5079 | 8.7962 | 5 | 117.5 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 8.4212 | 9.1553 | 9.8529 | 5 | 109.2 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 6.0413 | 6.6818 | 7.0441 | 5 | 149.7 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 6.9653 | 7.5329 | 8.4092 | 5 | 132.8 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 7.8680 | 7.9716 | 9.1102 | 5 | 125.4 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 7.7145 | 8.4075 | 33.5877 | 5 | 118.9 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 8.4243 | 8.7898 | 8.9743 | 5 | 113.8 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 9.1109 | 9.3698 | 10.7423 | 5 | 106.7 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 8.1332 | 8.3475 | 8.5078 | 5 | 119.8 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 7.9766 | 8.5228 | 8.8235 | 5 | 117.3 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 8.7495 | 8.8278 | 9.0960 | 5 | 113.3 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 7.7294 | 7.9760 | 9.6664 | 5 | 125.4 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 68.5120 | 99.2313 | 107.5435 | 5 | 10.1 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 83.7434 | 111.5181 | 117.1748 | 5 | 9.0 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 5.8598 | 6.4928 | 7.0065 | 5 | 154.0 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 7.2558 | 7.4394 | 10.1964 | 5 | 134.4 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 21.4529 | 22.1093 | 24.2939 | 5 | 45.2 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 26.1967 | 30.1732 | 41.2668 | 5 | 33.1 |

## Matched blended workload

- phase: `isolated`; target offered rate: 10 updates/s
- offered 0, committed 0, failed 0 over 0.000 s
- offered rate 0.000/s, completed rate 0.000/s
- per key: no keys touched
- ledger: 0 keys reconciled, 0 mismatch

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42137041`, 5 interleaved rounds
  1. family=bag;stmt=join_where;threshold=2000
  2. family=reuse;policy=notmat;refs=16;consumer=1%
  3. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  4. family=bag;stmt=join_where;threshold=180000
  5. family=bag;stmt=cte_joined_outer;threshold=20
  6. family=bag;stmt=cte_joined_outer_mat;threshold=20
  7. family=bag;stmt=join_where;threshold=20
  8. family=reuse;policy=mat;refs=16;consumer=1%
  9. family=reuse;policy=mat;refs=4;consumer=100%
  10. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  11. family=reuse;policy=mat;refs=1;consumer=1%
  12. family=reuse;policy=mat;refs=16;consumer=100%
  13. family=reuse;policy=default;refs=1;consumer=1%
  14. family=reuse;policy=notmat;refs=1;consumer=100%
  15. family=reuse;policy=default;refs=16;consumer=100%
  16. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  17. family=reuse;policy=notmat;refs=4;consumer=1%
  18. family=reuse;policy=default;refs=1;consumer=100%
  19. family=reuse;policy=notmat;refs=16;consumer=100%
  20. family=reuse;policy=mat;refs=1;consumer=100%
  21. family=reuse;policy=default;refs=4;consumer=100%
  22. family=reuse;policy=notmat;refs=4;consumer=100%
  23. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  24. family=reuse;policy=default;refs=4;consumer=1%
  25. family=bag;stmt=cte_joined_outer;threshold=180000
  26. family=reuse;policy=mat;refs=4;consumer=1%
  27. family=reuse;policy=default;refs=16;consumer=1%
  28. family=reuse;policy=notmat;refs=1;consumer=1%
  29. family=bag;stmt=cte_joined_outer;threshold=2000
  30. family=bag;stmt=cte_joined_outer_notmat;threshold=2000

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
| family=bag;stmt=cte_joined_outer;threshold=180000 | j05_cte_joined_outer | valid | `ec3f969a76df87e0` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20 | j05_cte_joined_outer | valid | `b7fcd297e7c061e3` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=2000 | j05_cte_joined_outer | valid | `1f9864e29d2571fe` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | j05m_cte_joined_outer_mat | valid | `3bfb72330473d746` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | j05m_cte_joined_outer_mat | valid | `934ec996faaf5f09` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | j05m_cte_joined_outer_mat | valid | `dd4fa2b57c7fb001` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | j05n_cte_joined_outer_notmat | valid | `2ae5ccfa060df841` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | j05n_cte_joined_outer_notmat | valid | `81a6b16d7fa0bcd4` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | j05n_cte_joined_outer_notmat | valid | `cea4582ba7eaa751` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=join_where;threshold=180000 | j01_join_where | valid | `468b44015c488212` |  |
| family=bag;stmt=join_where;threshold=20 | j01_join_where | valid | `0d7914cc8eadf84d` |  |
| family=bag;stmt=join_where;threshold=2000 | j01_join_where | valid | `b2451cd275b28f72` |  |
| family=reuse;policy=default;refs=16;consumer=1% | r_default_n16_c1 | valid | `657067656d8263c4` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=16;consumer=100% | r_default_n16_c100 | valid | `5c636c5399e736d8` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=1;consumer=1% | r_default_n1_c1 | valid | `10b47fef96fd6824` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=1;consumer=100% | r_default_n1_c100 | valid | `4b0fc303a68a0b49` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=4;consumer=1% | r_default_n4_c1 | valid | `e77f07d5e6a4503a` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=100% | r_default_n4_c100 | valid | `c3a41a9f40fa636c` | policy=default; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=1% | r_mat_n16_c1 | valid | `4dbe2dfcae740eae` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=100% | r_mat_n16_c100 | valid | `aca3bc92020de07c` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=1% | r_mat_n1_c1 | valid | `ab078b5079848a27` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=100% | r_mat_n1_c100 | valid | `6241948476d1ed00` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=1% | r_mat_n4_c1 | valid | `c64ac35e89850567` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=100% | r_mat_n4_c100 | valid | `f0ad8022db1b80b6` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=notmat;refs=16;consumer=1% | r_notmat_n16_c1 | valid | `69ba42f9bfcbaf19` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `39c538df37561609` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `f66483d8cbc542be` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `d42081e770e94ef2` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `9c214b41df0e91bd` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `2df3fe6c777bc7f1` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, DIST, DEBUG, FORMAT TEXT` supported on yugabyte
- one host: this run used node budget 2 CPU 3g and client budget 1 CPU 2g, so budget comparisons are separate classes here, not independent hardware
- isolated phase: no concurrent writers in this run; the blended read/write phase is a separate labelled run
- physical rank-index arm yb-range-rank; an engine whose planner ignores the index is a different physical design, not a faster one
- YugabyteDB DIST/DEBUG storage counters (rows and index rows scanned) are parsed per cell; RPC-level counters are not aggregated
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
