# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 30 cells, 30 valid, engine `yugabyte`, topology `single`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-v3b-topo-single-i1-20261005`
- run tag: `run/study-01/select-join-v3b-topo-single-i1-20261005`
- environment: `host-zenbook-ux5406sa`
- repo commit: `7243fee074bfb16fc5ad8b6afc1a2d7549080935`
- engine version: `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit`
- database engine image: `docker.io/yugabytedb/yugabyte:2025.2.6.0-b111` @ `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`
- benchmark client image: `localhost/charitybench-select:select-join-v3b-topo-single-i1-20261005` @ `f3f1e10c44cbcb4c7b5bc56ed95a4f071682b6f7d33b5e0a7b688c4271eeef88`
- inputs digest: `fba396164d3176debe7a9c1933075592133fa9ac0302b340ca8736ad1527c754`
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
| 20 (0.01%) | 2.0159 | missing | missing | missing | missing | missing | 1.8114 | 77.5173 | 2.4150 | missing |
| 2000 (1%) | 8.9062 | missing | missing | missing | missing | missing | 10.1694 | 88.2206 | 8.7812 | missing |
| 180000 (90%) | 105.4197 | missing | missing | missing | missing | missing | 105.3376 | 99.5245 | 100.0272 | missing |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | cte_joined_outer vs join_where | 2.0159 | 1.8114 | 0.8985 | false |
| 20 | cte_joined_outer_mat vs join_where | 2.0159 | 77.5173 | 38.4533 | false |
| 20 | cte_joined_outer_notmat vs join_where | 2.0159 | 2.4150 | 1.1980 | false |
| 2000 | cte_joined_outer vs join_where | 8.9062 | 10.1694 | 1.1418 | false |
| 2000 | cte_joined_outer_mat vs join_where | 8.9062 | 88.2206 | 9.9055 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 8.9062 | 8.7812 | 0.9860 | false |
| 180000 | cte_joined_outer vs join_where | 105.4197 | 105.3376 | 0.9992 | false |
| 180000 | cte_joined_outer_mat vs join_where | 105.4197 | 99.5245 | 0.9441 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 105.4197 | 100.0272 | 0.9488 | false |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_joined_outer_mat vs cte_joined_outer | 1.8114 | 77.5173 | 42.7949 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 1.8114 | 2.4150 | 1.3332 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 10.1694 | 88.2206 | 8.6751 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 10.1694 | 8.7812 | 0.8635 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 105.3376 | 99.5245 | 0.9448 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 105.3376 | 100.0272 | 0.9496 | false | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 6.9298 | 1 | `e97d0e8bc58f8529` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 8.8855 | 1 | `fbeb097620c760e6` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | missing | | | | |
| default | 2 | 100% | missing | | | | |
| default | 4 | 1% | 9.0704 | 4 | `da0b8ef3a319d9dc` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 9.6797 | 4 | `2c94f2d5af7d1981` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | missing | | | | |
| default | 8 | 100% | missing | | | | |
| default | 16 | 1% | 8.8007 | 16 | `100bc5ec645868e8` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 9.8736 | 16 | `dab643fc7fae9874` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 8.5555 | 1 | `ff4add51d7550a21` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 8.1036 | 1 | `a0d7a36464b17eac` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | missing | | | | |
| mat | 2 | 100% | missing | | | | |
| mat | 4 | 1% | 8.4936 | 4 | `023b89da07c73c91` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 8.1360 | 4 | `5601163629ba54e9` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | missing | | | | |
| mat | 8 | 100% | missing | | | | |
| mat | 16 | 1% | 9.6487 | 16 | `633521a2e119aca9` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 9.0270 | 16 | `7e3a422debb812c7` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 6.4678 | 1 | `da9b48012ae63e70` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 8.2653 | 1 | `8369a0adbf5f5585` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | missing | | | | |
| notmat | 2 | 100% | missing | | | | |
| notmat | 4 | 1% | 22.2474 | 4 | `cfb4d7bb10fc6e5c` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 28.0235 | 4 | `6773796e7bddaac9` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | missing | | | | |
| notmat | 8 | 100% | missing | | | | |
| notmat | 16 | 1% | 76.5929 | 16 | `4499332d50d664cc` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 86.7718 | 16 | `513f9163ff322fd7` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 1.2346 | false |
| notmat | 1 | 1% | 0.9333 | false |
| mat | 1 | 100% | 0.9120 | false |
| notmat | 1 | 100% | 0.9302 | false |
| mat | 4 | 1% | 0.9364 | false |
| notmat | 4 | 1% | 2.4527 | false |
| mat | 4 | 100% | 0.8405 | false |
| notmat | 4 | 100% | 2.8951 | false |
| mat | 16 | 1% | 1.0964 | false |
| notmat | 16 | 1% | 8.7031 | false |
| mat | 16 | 100% | 0.9143 | false |
| notmat | 16 | 100% | 8.7882 | false |

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

- client: nr_periods 301, nr_throttled 117, throttled_usec 4723409 | quota 100000 100000
, memory 2147483648

- yb-single: nr_periods 389, nr_throttled 79, throttled_usec 5163078 | quota 200000 100000, memory 3221225472

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
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 102.0876 | 105.3376 | 114.6068 | 5 | 9.5 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 1.5520 | 1.8114 | 2.3818 | 5 | 552.1 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 9.1970 | 10.1694 | 10.5085 | 5 | 98.3 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 98.6647 | 99.5245 | 117.8471 | 5 | 10.0 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 74.9103 | 77.5173 | 86.2943 | 5 | 12.9 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 81.2389 | 88.2206 | 110.7507 | 5 | 11.3 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 99.0036 | 100.0272 | 117.3641 | 5 | 10.0 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 2.1267 | 2.4150 | 2.6404 | 5 | 414.1 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 7.9077 | 8.7812 | 10.3246 | 5 | 113.9 |
| family=bag;stmt=join_where;threshold=180000 | valid | 99.8314 | 105.4197 | 110.5636 | 5 | 9.5 |
| family=bag;stmt=join_where;threshold=20 | valid | 1.7310 | 2.0159 | 2.1452 | 5 | 496.1 |
| family=bag;stmt=join_where;threshold=2000 | valid | 8.4229 | 8.9062 | 10.2308 | 5 | 112.3 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 8.6554 | 8.8007 | 9.7752 | 5 | 113.6 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 8.8193 | 9.8736 | 9.9269 | 5 | 101.3 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 6.2941 | 6.9298 | 8.1348 | 5 | 144.3 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 7.9662 | 8.8855 | 9.0451 | 5 | 112.5 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 8.5513 | 9.0704 | 9.4891 | 5 | 110.2 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 8.6011 | 9.6797 | 35.6594 | 5 | 103.3 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 9.4730 | 9.6487 | 10.4265 | 5 | 103.6 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 8.3955 | 9.0270 | 9.9423 | 5 | 110.8 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 7.2458 | 8.5555 | 22.2867 | 5 | 116.9 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 7.3019 | 8.1036 | 8.8776 | 5 | 123.4 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 7.4818 | 8.4936 | 9.4374 | 5 | 117.7 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 7.6040 | 8.1360 | 8.9787 | 5 | 122.9 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 69.4117 | 76.5929 | 82.2072 | 5 | 13.1 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 79.2734 | 86.7718 | 131.9684 | 5 | 11.5 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 6.0206 | 6.4678 | 7.4350 | 5 | 154.6 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 7.6651 | 8.2653 | 9.5051 | 5 | 121.0 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 19.8613 | 22.2474 | 27.7014 | 5 | 44.9 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 26.5426 | 28.0235 | 35.4052 | 5 | 35.7 |

## Matched blended workload

- phase: `isolated`; target offered rate: 10 updates/s
- offered 0, committed 0, failed 0 over 0.000 s
- offered rate 0.000/s, completed rate 0.000/s
- per key: no keys touched
- ledger: 0 keys reconciled, 0 mismatch

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42235562`, 5 interleaved rounds
  1. family=reuse;policy=mat;refs=16;consumer=100%
  2. family=bag;stmt=join_where;threshold=20
  3. family=bag;stmt=cte_joined_outer;threshold=180000
  4. family=reuse;policy=default;refs=1;consumer=100%
  5. family=reuse;policy=notmat;refs=4;consumer=100%
  6. family=reuse;policy=notmat;refs=4;consumer=1%
  7. family=bag;stmt=cte_joined_outer_mat;threshold=20
  8. family=reuse;policy=default;refs=16;consumer=1%
  9. family=reuse;policy=notmat;refs=1;consumer=1%
  10. family=reuse;policy=default;refs=16;consumer=100%
  11. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  12. family=reuse;policy=mat;refs=16;consumer=1%
  13. family=bag;stmt=join_where;threshold=2000
  14. family=reuse;policy=mat;refs=1;consumer=100%
  15. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  16. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  17. family=bag;stmt=cte_joined_outer;threshold=20
  18. family=reuse;policy=notmat;refs=16;consumer=100%
  19. family=reuse;policy=default;refs=4;consumer=100%
  20. family=reuse;policy=mat;refs=4;consumer=1%
  21. family=reuse;policy=mat;refs=4;consumer=100%
  22. family=reuse;policy=mat;refs=1;consumer=1%
  23. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  24. family=bag;stmt=join_where;threshold=180000
  25. family=bag;stmt=cte_joined_outer;threshold=2000
  26. family=reuse;policy=notmat;refs=16;consumer=1%
  27. family=reuse;policy=default;refs=4;consumer=1%
  28. family=reuse;policy=notmat;refs=1;consumer=100%
  29. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  30. family=reuse;policy=default;refs=1;consumer=1%

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
| family=bag;stmt=cte_joined_outer;threshold=180000 | j05_cte_joined_outer | valid | `4aa6072877a04dbd` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20 | j05_cte_joined_outer | valid | `22e615c2bd5cb052` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=2000 | j05_cte_joined_outer | valid | `d9c0e72798812cea` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | j05m_cte_joined_outer_mat | valid | `82445963dd36c45f` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | j05m_cte_joined_outer_mat | valid | `cd898b2786945383` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | j05m_cte_joined_outer_mat | valid | `acca4b1776de5538` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | j05n_cte_joined_outer_notmat | valid | `81feb9c7eddf6e3f` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | j05n_cte_joined_outer_notmat | valid | `c6988e1a2bb4c20a` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | j05n_cte_joined_outer_notmat | valid | `64a91609764fcace` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=join_where;threshold=180000 | j01_join_where | valid | `58fe8e461351e79e` |  |
| family=bag;stmt=join_where;threshold=20 | j01_join_where | valid | `67b36c969ff616c3` |  |
| family=bag;stmt=join_where;threshold=2000 | j01_join_where | valid | `e1208bc5f505b2a0` |  |
| family=reuse;policy=default;refs=16;consumer=1% | r_default_n16_c1 | valid | `100bc5ec645868e8` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=16;consumer=100% | r_default_n16_c100 | valid | `dab643fc7fae9874` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=1;consumer=1% | r_default_n1_c1 | valid | `e97d0e8bc58f8529` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=1;consumer=100% | r_default_n1_c100 | valid | `fbeb097620c760e6` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=4;consumer=1% | r_default_n4_c1 | valid | `da0b8ef3a319d9dc` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=100% | r_default_n4_c100 | valid | `2c94f2d5af7d1981` | policy=default; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=1% | r_mat_n16_c1 | valid | `633521a2e119aca9` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=100% | r_mat_n16_c100 | valid | `7e3a422debb812c7` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=1% | r_mat_n1_c1 | valid | `ff4add51d7550a21` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=100% | r_mat_n1_c100 | valid | `a0d7a36464b17eac` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=1% | r_mat_n4_c1 | valid | `023b89da07c73c91` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=100% | r_mat_n4_c100 | valid | `5601163629ba54e9` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=notmat;refs=16;consumer=1% | r_notmat_n16_c1 | valid | `4499332d50d664cc` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `513f9163ff322fd7` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `da9b48012ae63e70` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `8369a0adbf5f5585` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `cfb4d7bb10fc6e5c` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `6773796e7bddaac9` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, DIST, DEBUG, FORMAT TEXT` supported on yugabyte
- one host: this run used node budget 2 CPU 3g and client budget 1 CPU 2g, so budget comparisons are separate classes here, not independent hardware
- isolated phase: no concurrent writers in this run; the blended read/write phase is a separate labelled run
- physical rank-index arm yb-range-rank; an engine whose planner ignores the index is a different physical design, not a faster one
- YugabyteDB DIST/DEBUG storage counters (rows and index rows scanned) are parsed per cell; RPC-level counters are not aggregated
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
