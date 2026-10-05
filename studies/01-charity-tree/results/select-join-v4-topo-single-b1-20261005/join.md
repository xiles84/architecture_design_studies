# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 30 cells, 30 valid, engine `yugabyte`, topology `single`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-v4-topo-single-b1-20261005`
- run tag: `run/study-01/select-join-v4-topo-single-b1-20261005`
- environment: `host-zenbook-ux5406sa`
- repo commit: `ba821e7b686ac3adb1cda56dd45f15c8ae699ef8`
- engine version: `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit`
- database engine image: `docker.io/yugabytedb/yugabyte:2025.2.6.0-b111` @ `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`
- benchmark client image: `localhost/charitybench-select:select-join-v4-topo-single-b1-20261005` @ `b6ea3c007cf459112eaa8e9e85381a7dc25a7710ff07648fde3048f747df6ab0`
- inputs digest: `8321c3ddbaeb64c88aa2315bd66bd90c9cacdf42177e4b4540c1749b92205d27`
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
| 20 (0.01%) | 1.8427 | missing | missing | missing | missing | missing | 2.0173 | 82.3613 | 2.1840 | missing |
| 2000 (1%) | 8.9347 | missing | missing | missing | missing | missing | 8.3822 | 82.0311 | 8.7744 | missing |
| 180000 (90%) | 106.0982 | missing | missing | missing | missing | missing | 104.5941 | 99.6518 | 108.2737 | missing |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | cte_joined_outer vs join_where | 1.8427 | 2.0173 | 1.0947 | false |
| 20 | cte_joined_outer_mat vs join_where | 1.8427 | 82.3613 | 44.6949 | false |
| 20 | cte_joined_outer_notmat vs join_where | 1.8427 | 2.1840 | 1.1852 | false |
| 2000 | cte_joined_outer vs join_where | 8.9347 | 8.3822 | 0.9382 | false |
| 2000 | cte_joined_outer_mat vs join_where | 8.9347 | 82.0311 | 9.1812 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 8.9347 | 8.7744 | 0.9821 | false |
| 180000 | cte_joined_outer vs join_where | 106.0982 | 104.5941 | 0.9858 | false |
| 180000 | cte_joined_outer_mat vs join_where | 106.0982 | 99.6518 | 0.9392 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 106.0982 | 108.2737 | 1.0205 | false |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_joined_outer_mat vs cte_joined_outer | 2.0173 | 82.3613 | 40.8278 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 2.0173 | 2.1840 | 1.0827 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 8.3822 | 82.0311 | 9.7863 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 8.3822 | 8.7744 | 1.0468 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 104.5941 | 99.6518 | 0.9527 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 104.5941 | 108.2737 | 1.0352 | false | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 5.8544 | 1 | `bb9e4099fa1a2df6` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 9.0690 | 1 | `a47e1a3050df5d5c` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | missing | | | | |
| default | 2 | 100% | missing | | | | |
| default | 4 | 1% | 8.2388 | 4 | `67f7728bda61b3a6` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 8.8107 | 4 | `37159b149f0b8546` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | missing | | | | |
| default | 8 | 100% | missing | | | | |
| default | 16 | 1% | 8.8534 | 16 | `a7728d4e267f6ad3` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 9.8040 | 16 | `b80b8c5c05a08ab3` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 8.1215 | 1 | `743928056abdbcae` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 8.2521 | 1 | `a65dfeb6ea2169bc` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | missing | | | | |
| mat | 2 | 100% | missing | | | | |
| mat | 4 | 1% | 7.9730 | 4 | `ff628c2b5fe2bee0` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 9.1206 | 4 | `922d2b042a39eff6` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | missing | | | | |
| mat | 8 | 100% | missing | | | | |
| mat | 16 | 1% | 10.0662 | 16 | `274411b24a09586d` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 10.5017 | 16 | `78761dd0f2893b64` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 7.4205 | 1 | `fa132d157dd63aec` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 8.5940 | 1 | `78fa2ab05fde24f8` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | missing | | | | |
| notmat | 2 | 100% | missing | | | | |
| notmat | 4 | 1% | 22.8712 | 4 | `0a716c6b56cee0e6` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 29.2379 | 4 | `61283b84fd999c38` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | missing | | | | |
| notmat | 8 | 100% | missing | | | | |
| notmat | 16 | 1% | 87.5949 | 16 | `1bad92722e69f4f7` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 120.6135 | 16 | `839ce7846ef679b8` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 1.3872 | false |
| notmat | 1 | 1% | 1.2675 | false |
| mat | 1 | 100% | 0.9099 | false |
| notmat | 1 | 100% | 0.9476 | false |
| mat | 4 | 1% | 0.9677 | false |
| notmat | 4 | 1% | 2.7760 | false |
| mat | 4 | 100% | 1.0352 | false |
| notmat | 4 | 100% | 3.3184 | false |
| mat | 16 | 1% | 1.1370 | false |
| notmat | 16 | 1% | 9.8939 | false |
| mat | 16 | 100% | 1.0712 | false |
| notmat | 16 | 100% | 12.3025 | false |

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

- client: nr_periods 285, nr_throttled 116, throttled_usec 4144397 | quota 100000 100000
, memory 2147483648

- yb-single: nr_periods 368, nr_throttled 71, throttled_usec 5308192 | quota 200000 100000, memory 3221225472

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
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 99.5786 | 104.5941 | 111.7884 | 5 | 9.6 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 1.7072 | 2.0173 | 2.3440 | 5 | 495.7 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 8.2111 | 8.3822 | 8.7519 | 5 | 119.3 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 97.7258 | 99.6518 | 120.8963 | 5 | 10.0 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 76.8035 | 82.3613 | 103.6257 | 5 | 12.1 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 80.8459 | 82.0311 | 92.2736 | 5 | 12.2 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 99.4517 | 108.2737 | 118.3648 | 5 | 9.2 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 2.0704 | 2.1840 | 2.9810 | 5 | 457.9 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 8.2720 | 8.7744 | 12.1425 | 5 | 114.0 |
| family=bag;stmt=join_where;threshold=180000 | valid | 103.8369 | 106.0982 | 115.8044 | 5 | 9.4 |
| family=bag;stmt=join_where;threshold=20 | valid | 1.6899 | 1.8427 | 2.0586 | 5 | 542.7 |
| family=bag;stmt=join_where;threshold=2000 | valid | 8.5776 | 8.9347 | 10.4428 | 5 | 111.9 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 7.9999 | 8.8534 | 11.1728 | 5 | 113.0 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 9.4441 | 9.8040 | 11.3400 | 5 | 102.0 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 5.6436 | 5.8544 | 6.7729 | 5 | 170.8 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 8.6259 | 9.0690 | 9.3519 | 5 | 110.3 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 7.8042 | 8.2388 | 8.7012 | 5 | 121.4 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 8.1935 | 8.8107 | 9.0676 | 5 | 113.5 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 9.7899 | 10.0662 | 11.0609 | 5 | 99.3 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 10.1642 | 10.5017 | 11.2763 | 5 | 95.2 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 7.5896 | 8.1215 | 8.3987 | 5 | 123.1 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 7.3643 | 8.2521 | 8.9766 | 5 | 121.2 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 7.7267 | 7.9730 | 9.4511 | 5 | 125.4 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 8.6248 | 9.1206 | 9.5161 | 5 | 109.6 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 73.5513 | 87.5949 | 95.1276 | 5 | 11.4 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 105.4560 | 120.6135 | 128.0777 | 5 | 8.3 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 7.3248 | 7.4205 | 7.5833 | 5 | 134.8 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 8.0381 | 8.5940 | 9.4551 | 5 | 116.4 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 21.3264 | 22.8712 | 23.8883 | 5 | 43.7 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 26.7741 | 29.2379 | 30.6393 | 5 | 34.2 |

## Matched blended workload

- phase: `blended`; target offered rate: 10 updates/s
- offered 146, committed 146, failed 0 over 14.562 s
- offered rate 10.026/s, completed rate 10.026/s
- writer operations per endpoint: [146]
- reader operations overlapping a committed write: 69 of 210
- reader retries of transient conflicts: 0
- per key: 146 keys, min 1, max 1 per key
- ledger: 146 keys reconciled, 0 mismatch
- update latency: avg 4.004 ms, min 1.333, max 108.284 over 146

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42410088`, 5 interleaved rounds
  1. family=reuse;policy=default;refs=1;consumer=1%
  2. family=bag;stmt=cte_joined_outer;threshold=20
  3. family=reuse;policy=notmat;refs=16;consumer=100%
  4. family=reuse;policy=mat;refs=16;consumer=100%
  5. family=reuse;policy=default;refs=4;consumer=100%
  6. family=reuse;policy=default;refs=16;consumer=100%
  7. family=reuse;policy=mat;refs=4;consumer=1%
  8. family=bag;stmt=join_where;threshold=20
  9. family=bag;stmt=cte_joined_outer;threshold=180000
  10. family=reuse;policy=notmat;refs=1;consumer=1%
  11. family=bag;stmt=cte_joined_outer;threshold=2000
  12. family=reuse;policy=mat;refs=1;consumer=1%
  13. family=reuse;policy=default;refs=16;consumer=1%
  14. family=reuse;policy=default;refs=4;consumer=1%
  15. family=reuse;policy=notmat;refs=4;consumer=100%
  16. family=bag;stmt=cte_joined_outer_mat;threshold=20
  17. family=reuse;policy=mat;refs=4;consumer=100%
  18. family=reuse;policy=notmat;refs=16;consumer=1%
  19. family=reuse;policy=notmat;refs=4;consumer=1%
  20. family=reuse;policy=mat;refs=1;consumer=100%
  21. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  22. family=reuse;policy=mat;refs=16;consumer=1%
  23. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  24. family=bag;stmt=join_where;threshold=180000
  25. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  26. family=reuse;policy=notmat;refs=1;consumer=100%
  27. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  28. family=bag;stmt=join_where;threshold=2000
  29. family=bag;stmt=cte_joined_outer_mat;threshold=180000
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
| family=bag;stmt=cte_joined_outer;threshold=180000 | j05_cte_joined_outer | valid | `c179ec9811336150` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20 | j05_cte_joined_outer | valid | `ca27f40e07f48ea8` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=2000 | j05_cte_joined_outer | valid | `fb4cdc310976aef2` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | j05m_cte_joined_outer_mat | valid | `8f54115b9d1d4599` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | j05m_cte_joined_outer_mat | valid | `aaf5773dd80364ec` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | j05m_cte_joined_outer_mat | valid | `1edd3b72437e6ceb` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | j05n_cte_joined_outer_notmat | valid | `80065f07cae07aa1` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | j05n_cte_joined_outer_notmat | valid | `0a45b019d577c70d` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | j05n_cte_joined_outer_notmat | valid | `384a83a0c7c02daa` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=join_where;threshold=180000 | j01_join_where | valid | `4a375131d4c70ad5` |  |
| family=bag;stmt=join_where;threshold=20 | j01_join_where | valid | `4efb344393262b89` |  |
| family=bag;stmt=join_where;threshold=2000 | j01_join_where | valid | `5fb2b39d4a2f084a` |  |
| family=reuse;policy=default;refs=16;consumer=1% | r_default_n16_c1 | valid | `a7728d4e267f6ad3` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=16;consumer=100% | r_default_n16_c100 | valid | `b80b8c5c05a08ab3` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=1;consumer=1% | r_default_n1_c1 | valid | `bb9e4099fa1a2df6` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=1;consumer=100% | r_default_n1_c100 | valid | `a47e1a3050df5d5c` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=4;consumer=1% | r_default_n4_c1 | valid | `67f7728bda61b3a6` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=100% | r_default_n4_c100 | valid | `37159b149f0b8546` | policy=default; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=1% | r_mat_n16_c1 | valid | `274411b24a09586d` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=100% | r_mat_n16_c100 | valid | `78761dd0f2893b64` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=1% | r_mat_n1_c1 | valid | `743928056abdbcae` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=100% | r_mat_n1_c100 | valid | `a65dfeb6ea2169bc` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=1% | r_mat_n4_c1 | valid | `ff628c2b5fe2bee0` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=100% | r_mat_n4_c100 | valid | `922d2b042a39eff6` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=notmat;refs=16;consumer=1% | r_notmat_n16_c1 | valid | `1bad92722e69f4f7` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `839ce7846ef679b8` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `fa132d157dd63aec` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `78fa2ab05fde24f8` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `0a716c6b56cee0e6` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `61283b84fd999c38` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, DIST, DEBUG, FORMAT TEXT` supported on yugabyte
- one host: this run used node budget 2 CPU 3g and client budget 1 CPU 2g, so budget comparisons are separate classes here, not independent hardware
- one client process with one fixed maintenance writer at 10 updates/s; readers are sequential; the reported sample latency includes retry and backoff time
- physical rank-index arm yb-range-rank; an engine whose planner ignores the index is a different physical design, not a faster one
- YugabyteDB DIST/DEBUG storage counters (rows and index rows scanned) are parsed per cell; RPC-level counters are not aggregated
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
