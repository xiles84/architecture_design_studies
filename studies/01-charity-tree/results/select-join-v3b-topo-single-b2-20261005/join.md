# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 30 cells, 30 valid, engine `yugabyte`, topology `single`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-v3b-topo-single-b2-20261005`
- run tag: `run/study-01/select-join-v3b-topo-single-b2-20261005`
- environment: `host-zenbook-ux5406sa`
- repo commit: `7243fee074bfb16fc5ad8b6afc1a2d7549080935`
- engine version: `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit`
- database engine image: `docker.io/yugabytedb/yugabyte:2025.2.6.0-b111` @ `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`
- benchmark client image: `localhost/charitybench-select:select-join-v3b-topo-single-b2-20261005` @ `f3f1e10c44cbcb4c7b5bc56ed95a4f071682b6f7d33b5e0a7b688c4271eeef88`
- inputs digest: `9b08b86d7d22b914399dab929cb38fcbc3802b2af24828b5990d27dfc3e88704`
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
| 20 (0.01%) | 1.9582 | missing | missing | missing | missing | missing | 2.0619 | 90.6717 | 2.2333 | missing |
| 2000 (1%) | 9.6818 | missing | missing | missing | missing | missing | 10.3959 | 93.2029 | 9.3842 | missing |
| 180000 (90%) | 114.8487 | missing | missing | missing | missing | missing | 119.2145 | 110.1336 | 117.6864 | missing |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | cte_joined_outer vs join_where | 1.9582 | 2.0619 | 1.0530 | false |
| 20 | cte_joined_outer_mat vs join_where | 1.9582 | 90.6717 | 46.3035 | false |
| 20 | cte_joined_outer_notmat vs join_where | 1.9582 | 2.2333 | 1.1405 | false |
| 2000 | cte_joined_outer vs join_where | 9.6818 | 10.3959 | 1.0737 | false |
| 2000 | cte_joined_outer_mat vs join_where | 9.6818 | 93.2029 | 9.6266 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 9.6818 | 9.3842 | 0.9693 | false |
| 180000 | cte_joined_outer vs join_where | 114.8487 | 119.2145 | 1.0380 | false |
| 180000 | cte_joined_outer_mat vs join_where | 114.8487 | 110.1336 | 0.9589 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 114.8487 | 117.6864 | 1.0247 | false |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_joined_outer_mat vs cte_joined_outer | 2.0619 | 90.6717 | 43.9750 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 2.0619 | 2.2333 | 1.0831 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 10.3959 | 93.2029 | 8.9654 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 10.3959 | 9.3842 | 0.9027 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 119.2145 | 110.1336 | 0.9238 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 119.2145 | 117.6864 | 0.9872 | false | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 7.3297 | 1 | `69732cb2dddfe858` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 10.0718 | 1 | `857b743f5419be3f` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | missing | | | | |
| default | 2 | 100% | missing | | | | |
| default | 4 | 1% | 10.1148 | 4 | `60d69081d1866309` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 9.4066 | 4 | `5c09851aa60bf630` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | missing | | | | |
| default | 8 | 100% | missing | | | | |
| default | 16 | 1% | 10.3767 | 16 | `fa3b394e0edfb260` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 10.4461 | 16 | `329de7336a08da2f` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 9.0427 | 1 | `cae5d12b9da7ddac` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 10.6477 | 1 | `3ccd31c97fae05ef` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | missing | | | | |
| mat | 2 | 100% | missing | | | | |
| mat | 4 | 1% | 10.1524 | 4 | `2e2c82e5c9fccb74` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 9.7514 | 4 | `570527a5b99ce134` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | missing | | | | |
| mat | 8 | 100% | missing | | | | |
| mat | 16 | 1% | 11.1652 | 16 | `4df7bbe90d6dbf44` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 10.8862 | 16 | `6055dea32507a33b` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 7.6100 | 1 | `c765eb803cf89d76` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 8.7710 | 1 | `9fbc77819d74db4e` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | missing | | | | |
| notmat | 2 | 100% | missing | | | | |
| notmat | 4 | 1% | 39.2456 | 4 | `4ef337d18d9fc300` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 32.7979 | 4 | `0ed9b622c65617cc` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | missing | | | | |
| notmat | 8 | 100% | missing | | | | |
| notmat | 16 | 1% | 102.1546 | 16 | `2e76030f04e6fb4c` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 129.2604 | 16 | `cfdd6813e4eca1d8` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 1.2337 | false |
| notmat | 1 | 1% | 1.0382 | false |
| mat | 1 | 100% | 1.0572 | false |
| notmat | 1 | 100% | 0.8709 | false |
| mat | 4 | 1% | 1.0037 | false |
| notmat | 4 | 1% | 3.8800 | false |
| mat | 4 | 100% | 1.0367 | false |
| notmat | 4 | 100% | 3.4867 | false |
| mat | 16 | 1% | 1.0760 | false |
| notmat | 16 | 1% | 9.8446 | false |
| mat | 16 | 100% | 1.0421 | false |
| notmat | 16 | 100% | 12.3740 | false |

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

- client: nr_periods 313, nr_throttled 123, throttled_usec 5450145 | quota 100000 100000
, memory 2147483648

- yb-single: nr_periods 402, nr_throttled 88, throttled_usec 7444711 | quota 200000 100000, memory 3221225472

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
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 114.4907 | 119.2145 | 128.1636 | 5 | 8.4 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 1.7679 | 2.0619 | 2.3073 | 5 | 485.0 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 10.0478 | 10.3959 | 11.6564 | 5 | 96.2 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 104.8894 | 110.1336 | 127.3364 | 5 | 9.1 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 86.4363 | 90.6717 | 95.9746 | 5 | 11.0 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 89.3704 | 93.2029 | 107.3611 | 5 | 10.7 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 111.8021 | 117.6864 | 125.2274 | 5 | 8.5 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 2.0806 | 2.2333 | 2.6106 | 5 | 447.8 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 8.9128 | 9.3842 | 10.0300 | 5 | 106.6 |
| family=bag;stmt=join_where;threshold=180000 | valid | 111.4025 | 114.8487 | 130.7877 | 5 | 8.7 |
| family=bag;stmt=join_where;threshold=20 | valid | 1.9230 | 1.9582 | 2.0370 | 5 | 510.7 |
| family=bag;stmt=join_where;threshold=2000 | valid | 8.9176 | 9.6818 | 9.8002 | 5 | 103.3 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 9.8171 | 10.3767 | 12.4577 | 5 | 96.4 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 10.0846 | 10.4461 | 14.3819 | 5 | 95.7 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 6.5984 | 7.3297 | 8.7182 | 5 | 136.4 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 9.3968 | 10.0718 | 10.4195 | 5 | 99.3 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 9.6334 | 10.1148 | 10.5611 | 5 | 98.9 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 9.0847 | 9.4066 | 9.6445 | 5 | 106.3 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 10.3778 | 11.1652 | 11.4734 | 5 | 89.6 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 10.5437 | 10.8862 | 13.5174 | 5 | 91.9 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 8.6306 | 9.0427 | 10.9599 | 5 | 110.6 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 8.2358 | 10.6477 | 12.3075 | 5 | 93.9 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 9.7878 | 10.1524 | 11.0001 | 5 | 98.5 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 9.5039 | 9.7514 | 10.0036 | 5 | 102.5 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 78.4634 | 102.1546 | 109.7012 | 5 | 9.8 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 126.6263 | 129.2604 | 174.2080 | 5 | 7.7 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 7.3400 | 7.6100 | 8.0642 | 5 | 131.4 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 7.9635 | 8.7710 | 9.3553 | 5 | 114.0 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 24.4288 | 39.2456 | 53.2717 | 5 | 25.5 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 31.7241 | 32.7979 | 33.1163 | 5 | 30.5 |

## Matched blended workload

- phase: `blended`; target offered rate: 10 updates/s
- offered 160, committed 160, failed 0 over 15.963 s
- offered rate 10.023/s, completed rate 10.023/s
- writer operations per endpoint: [160]
- reader operations overlapping a committed write: 74 of 210
- reader retries of transient conflicts: 0
- per key: 160 keys, min 1, max 1 per key
- ledger: 160 keys reconciled, 0 mismatch
- update latency: avg 4.238 ms, min 1.636, max 66.537 over 160

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42249168`, 5 interleaved rounds
  1. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  2. family=reuse;policy=notmat;refs=1;consumer=1%
  3. family=reuse;policy=notmat;refs=1;consumer=100%
  4. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  5. family=reuse;policy=default;refs=4;consumer=1%
  6. family=reuse;policy=notmat;refs=16;consumer=1%
  7. family=reuse;policy=notmat;refs=16;consumer=100%
  8. family=reuse;policy=mat;refs=16;consumer=1%
  9. family=reuse;policy=notmat;refs=4;consumer=1%
  10. family=reuse;policy=default;refs=16;consumer=100%
  11. family=reuse;policy=mat;refs=1;consumer=1%
  12. family=bag;stmt=join_where;threshold=2000
  13. family=bag;stmt=cte_joined_outer;threshold=2000
  14. family=reuse;policy=mat;refs=16;consumer=100%
  15. family=reuse;policy=default;refs=1;consumer=1%
  16. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  17. family=reuse;policy=default;refs=4;consumer=100%
  18. family=bag;stmt=join_where;threshold=180000
  19. family=reuse;policy=notmat;refs=4;consumer=100%
  20. family=bag;stmt=cte_joined_outer_mat;threshold=20
  21. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  22. family=reuse;policy=mat;refs=4;consumer=100%
  23. family=bag;stmt=join_where;threshold=20
  24. family=bag;stmt=cte_joined_outer;threshold=180000
  25. family=reuse;policy=mat;refs=4;consumer=1%
  26. family=reuse;policy=default;refs=16;consumer=1%
  27. family=reuse;policy=mat;refs=1;consumer=100%
  28. family=bag;stmt=cte_joined_outer;threshold=20
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
| family=bag;stmt=cte_joined_outer;threshold=180000 | j05_cte_joined_outer | valid | `12abdc86386153a2` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20 | j05_cte_joined_outer | valid | `a17c0fec4ea013a6` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=2000 | j05_cte_joined_outer | valid | `e3d0a68a40a7a7fd` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | j05m_cte_joined_outer_mat | valid | `01826b92638ae7a6` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | j05m_cte_joined_outer_mat | valid | `86bf473a2ed264dd` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | j05m_cte_joined_outer_mat | valid | `7b49a0d3479ef9ed` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | j05n_cte_joined_outer_notmat | valid | `dff3ac56be554a8f` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | j05n_cte_joined_outer_notmat | valid | `f301df51fec2b5e2` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | j05n_cte_joined_outer_notmat | valid | `cf151c5827019b5f` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=join_where;threshold=180000 | j01_join_where | valid | `e287033daa355865` |  |
| family=bag;stmt=join_where;threshold=20 | j01_join_where | valid | `9ea3310dd61a6920` |  |
| family=bag;stmt=join_where;threshold=2000 | j01_join_where | valid | `9c2dc5e4c53bae14` |  |
| family=reuse;policy=default;refs=16;consumer=1% | r_default_n16_c1 | valid | `fa3b394e0edfb260` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=16;consumer=100% | r_default_n16_c100 | valid | `329de7336a08da2f` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=1;consumer=1% | r_default_n1_c1 | valid | `69732cb2dddfe858` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=1;consumer=100% | r_default_n1_c100 | valid | `857b743f5419be3f` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=4;consumer=1% | r_default_n4_c1 | valid | `60d69081d1866309` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=100% | r_default_n4_c100 | valid | `5c09851aa60bf630` | policy=default; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=1% | r_mat_n16_c1 | valid | `4df7bbe90d6dbf44` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=100% | r_mat_n16_c100 | valid | `6055dea32507a33b` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=1% | r_mat_n1_c1 | valid | `cae5d12b9da7ddac` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=100% | r_mat_n1_c100 | valid | `3ccd31c97fae05ef` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=1% | r_mat_n4_c1 | valid | `2e2c82e5c9fccb74` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=100% | r_mat_n4_c100 | valid | `570527a5b99ce134` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=notmat;refs=16;consumer=1% | r_notmat_n16_c1 | valid | `2e76030f04e6fb4c` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `cfdd6813e4eca1d8` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `c765eb803cf89d76` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `9fbc77819d74db4e` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `4ef337d18d9fc300` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `0ed9b622c65617cc` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, DIST, DEBUG, FORMAT TEXT` supported on yugabyte
- one host: this run used node budget 2 CPU 3g and client budget 1 CPU 2g, so budget comparisons are separate classes here, not independent hardware
- one client process with one fixed maintenance writer at 10 updates/s; readers are sequential; the reported sample latency includes retry and backoff time
- physical rank-index arm yb-range-rank; an engine whose planner ignores the index is a different physical design, not a faster one
- YugabyteDB DIST/DEBUG storage counters (rows and index rows scanned) are parsed per cell; RPC-level counters are not aggregated
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
