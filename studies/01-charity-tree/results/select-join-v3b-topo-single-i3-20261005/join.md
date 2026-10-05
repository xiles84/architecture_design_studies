# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 30 cells, 30 valid, engine `yugabyte`, topology `single`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-v3b-topo-single-i3-20261005`
- run tag: `run/study-01/select-join-v3b-topo-single-i3-20261005`
- environment: `host-zenbook-ux5406sa`
- repo commit: `7243fee074bfb16fc5ad8b6afc1a2d7549080935`
- engine version: `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit`
- database engine image: `docker.io/yugabytedb/yugabyte:2025.2.6.0-b111` @ `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`
- benchmark client image: `localhost/charitybench-select:select-join-v3b-topo-single-i3-20261005` @ `f3f1e10c44cbcb4c7b5bc56ed95a4f071682b6f7d33b5e0a7b688c4271eeef88`
- inputs digest: `a7a8522c80f52a19e0a32e5bc1b01a5bc4b3431c482716d6385c1deeb6366229`
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
| 20 (0.01%) | 2.0646 | missing | missing | missing | missing | missing | 2.3719 | 117.2604 | 2.3966 | missing |
| 2000 (1%) | 12.4329 | missing | missing | missing | missing | missing | 12.2369 | 114.9470 | 12.0151 | missing |
| 180000 (90%) | 138.4676 | missing | missing | missing | missing | missing | 142.8837 | 133.1028 | 137.6820 | missing |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | cte_joined_outer vs join_where | 2.0646 | 2.3719 | 1.1488 | false |
| 20 | cte_joined_outer_mat vs join_where | 2.0646 | 117.2604 | 56.7947 | false |
| 20 | cte_joined_outer_notmat vs join_where | 2.0646 | 2.3966 | 1.1608 | false |
| 2000 | cte_joined_outer vs join_where | 12.4329 | 12.2369 | 0.9842 | false |
| 2000 | cte_joined_outer_mat vs join_where | 12.4329 | 114.9470 | 9.2454 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 12.4329 | 12.0151 | 0.9664 | false |
| 180000 | cte_joined_outer vs join_where | 138.4676 | 142.8837 | 1.0319 | false |
| 180000 | cte_joined_outer_mat vs join_where | 138.4676 | 133.1028 | 0.9613 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 138.4676 | 137.6820 | 0.9943 | false |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_joined_outer_mat vs cte_joined_outer | 2.3719 | 117.2604 | 49.4365 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 2.3719 | 2.3966 | 1.0104 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 12.2369 | 114.9470 | 9.3935 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 12.2369 | 12.0151 | 0.9819 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 142.8837 | 133.1028 | 0.9315 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 142.8837 | 137.6820 | 0.9636 | false | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 8.0703 | 1 | `5484f981760f8643` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 10.6639 | 1 | `6ddd2eee55c58480` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | missing | | | | |
| default | 2 | 100% | missing | | | | |
| default | 4 | 1% | 11.1384 | 4 | `01a150e13a4cffe8` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 10.5764 | 4 | `f095c41575860818` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | missing | | | | |
| default | 8 | 100% | missing | | | | |
| default | 16 | 1% | 11.6433 | 16 | `eec6a5cbc45a6d2c` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 12.0289 | 16 | `17a85bc26559e5c1` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 10.2226 | 1 | `9ed3c1f399783eaf` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 10.6396 | 1 | `499cb92442729431` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | missing | | | | |
| mat | 2 | 100% | missing | | | | |
| mat | 4 | 1% | 10.7332 | 4 | `60c0167c60ea22bf` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 11.5837 | 4 | `ea7d86eebae08c51` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | missing | | | | |
| mat | 8 | 100% | missing | | | | |
| mat | 16 | 1% | 12.0625 | 16 | `8872f51cf658d8b1` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 12.2842 | 16 | `c3724a627d7e90f4` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 8.6349 | 1 | `e50621f7bfccda66` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 10.8463 | 1 | `ad8abdc9e2a8ec6f` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | missing | | | | |
| notmat | 2 | 100% | missing | | | | |
| notmat | 4 | 1% | 29.9783 | 4 | `eefbbe6b352816b9` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 38.7751 | 4 | `de6fba4c1cef016d` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | missing | | | | |
| notmat | 8 | 100% | missing | | | | |
| notmat | 16 | 1% | 115.8743 | 16 | `93dc9c2b4eedd785` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 157.0835 | 16 | `e0480b27b2330ef3` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 1.2667 | false |
| notmat | 1 | 1% | 1.0700 | false |
| mat | 1 | 100% | 0.9977 | false |
| notmat | 1 | 100% | 1.0171 | false |
| mat | 4 | 1% | 0.9636 | false |
| notmat | 4 | 1% | 2.6914 | false |
| mat | 4 | 100% | 1.0952 | false |
| notmat | 4 | 100% | 3.6662 | false |
| mat | 16 | 1% | 1.0360 | false |
| notmat | 16 | 1% | 9.9520 | false |
| mat | 16 | 100% | 1.0212 | false |
| notmat | 16 | 100% | 13.0588 | false |

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

- client: nr_periods 404, nr_throttled 185, throttled_usec 7627490 | quota 100000 100000
, memory 2147483648

- yb-single: nr_periods 517, nr_throttled 114, throttled_usec 9189479 | quota 200000 100000, memory 3221225472

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
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 135.1126 | 142.8837 | 185.5889 | 5 | 7.0 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 2.0401 | 2.3719 | 2.5464 | 5 | 421.6 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 10.5747 | 12.2369 | 13.5023 | 5 | 81.7 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 117.8566 | 133.1028 | 172.4610 | 5 | 7.5 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 99.2089 | 117.2604 | 122.1728 | 5 | 8.5 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 88.7538 | 114.9470 | 121.4206 | 5 | 8.7 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 125.8222 | 137.6820 | 144.6085 | 5 | 7.3 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 2.2964 | 2.3966 | 3.0445 | 5 | 417.3 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 9.9002 | 12.0151 | 12.1880 | 5 | 83.2 |
| family=bag;stmt=join_where;threshold=180000 | valid | 115.1605 | 138.4676 | 164.2089 | 5 | 7.2 |
| family=bag;stmt=join_where;threshold=20 | valid | 2.0019 | 2.0646 | 2.6327 | 5 | 484.3 |
| family=bag;stmt=join_where;threshold=2000 | valid | 11.6724 | 12.4329 | 43.4805 | 5 | 80.4 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 10.0513 | 11.6433 | 17.5739 | 5 | 85.9 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 11.0921 | 12.0289 | 12.9154 | 5 | 83.1 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 7.5463 | 8.0703 | 27.6549 | 5 | 123.9 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 9.6190 | 10.6639 | 11.2142 | 5 | 93.8 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 9.1454 | 11.1384 | 12.7047 | 5 | 89.8 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 9.6228 | 10.5764 | 12.0355 | 5 | 94.6 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 11.2738 | 12.0625 | 12.7048 | 5 | 82.9 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 10.2694 | 12.2842 | 13.0615 | 5 | 81.4 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 8.8732 | 10.2226 | 10.5685 | 5 | 97.8 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 8.9135 | 10.6396 | 11.9595 | 5 | 94.0 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 8.5022 | 10.7332 | 11.8775 | 5 | 93.2 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 10.5670 | 11.5837 | 11.9642 | 5 | 86.3 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 92.6090 | 115.8743 | 153.4888 | 5 | 8.6 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 102.2758 | 157.0835 | 176.0871 | 5 | 6.4 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 7.4663 | 8.6349 | 10.0968 | 5 | 115.8 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 8.7754 | 10.8463 | 13.9957 | 5 | 92.2 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 24.4607 | 29.9783 | 33.9258 | 5 | 33.4 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 35.8462 | 38.7751 | 41.6286 | 5 | 25.8 |

## Matched blended workload

- phase: `isolated`; target offered rate: 10 updates/s
- offered 0, committed 0, failed 0 over 0.000 s
- offered rate 0.000/s, completed rate 0.000/s
- per key: no keys touched
- ledger: 0 keys reconciled, 0 mismatch

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42323780`, 5 interleaved rounds
  1. family=reuse;policy=notmat;refs=4;consumer=100%
  2. family=reuse;policy=notmat;refs=16;consumer=1%
  3. family=bag;stmt=cte_joined_outer;threshold=20
  4. family=reuse;policy=default;refs=4;consumer=1%
  5. family=bag;stmt=join_where;threshold=20
  6. family=reuse;policy=default;refs=1;consumer=1%
  7. family=reuse;policy=mat;refs=1;consumer=1%
  8. family=reuse;policy=default;refs=16;consumer=1%
  9. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  10. family=reuse;policy=notmat;refs=1;consumer=1%
  11. family=reuse;policy=notmat;refs=1;consumer=100%
  12. family=reuse;policy=mat;refs=16;consumer=100%
  13. family=bag;stmt=join_where;threshold=180000
  14. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  15. family=reuse;policy=mat;refs=1;consumer=100%
  16. family=reuse;policy=notmat;refs=4;consumer=1%
  17. family=reuse;policy=mat;refs=4;consumer=1%
  18. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  19. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  20. family=bag;stmt=cte_joined_outer_mat;threshold=20
  21. family=reuse;policy=mat;refs=4;consumer=100%
  22. family=bag;stmt=cte_joined_outer;threshold=2000
  23. family=reuse;policy=mat;refs=16;consumer=1%
  24. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  25. family=reuse;policy=default;refs=1;consumer=100%
  26. family=reuse;policy=default;refs=16;consumer=100%
  27. family=reuse;policy=notmat;refs=16;consumer=100%
  28. family=bag;stmt=join_where;threshold=2000
  29. family=bag;stmt=cte_joined_outer;threshold=180000
  30. family=reuse;policy=default;refs=4;consumer=100%

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
| family=bag;stmt=cte_joined_outer;threshold=180000 | j05_cte_joined_outer | valid | `0bb184c9dfce4573` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20 | j05_cte_joined_outer | valid | `191891d1a902556b` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=2000 | j05_cte_joined_outer | valid | `1ecd5ba5f84307ad` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | j05m_cte_joined_outer_mat | valid | `b4a013e60ce8f460` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | j05m_cte_joined_outer_mat | valid | `5bd9fd682685ca10` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | j05m_cte_joined_outer_mat | valid | `3324d4a99248fd29` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | j05n_cte_joined_outer_notmat | valid | `e859ffa67c8ec8e2` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | j05n_cte_joined_outer_notmat | valid | `901a54faed1bae91` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | j05n_cte_joined_outer_notmat | valid | `3f27b16ad2fcd11a` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=join_where;threshold=180000 | j01_join_where | valid | `87bb32002cf3ff7f` |  |
| family=bag;stmt=join_where;threshold=20 | j01_join_where | valid | `b1e9a3ee984deec3` |  |
| family=bag;stmt=join_where;threshold=2000 | j01_join_where | valid | `46eb8216a6ee0c2b` |  |
| family=reuse;policy=default;refs=16;consumer=1% | r_default_n16_c1 | valid | `eec6a5cbc45a6d2c` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=16;consumer=100% | r_default_n16_c100 | valid | `17a85bc26559e5c1` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=1;consumer=1% | r_default_n1_c1 | valid | `5484f981760f8643` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=1;consumer=100% | r_default_n1_c100 | valid | `6ddd2eee55c58480` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=4;consumer=1% | r_default_n4_c1 | valid | `01a150e13a4cffe8` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=100% | r_default_n4_c100 | valid | `f095c41575860818` | policy=default; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=1% | r_mat_n16_c1 | valid | `8872f51cf658d8b1` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=100% | r_mat_n16_c100 | valid | `c3724a627d7e90f4` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=1% | r_mat_n1_c1 | valid | `9ed3c1f399783eaf` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=100% | r_mat_n1_c100 | valid | `499cb92442729431` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=1% | r_mat_n4_c1 | valid | `60c0167c60ea22bf` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=100% | r_mat_n4_c100 | valid | `ea7d86eebae08c51` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=notmat;refs=16;consumer=1% | r_notmat_n16_c1 | valid | `93dc9c2b4eedd785` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `e0480b27b2330ef3` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `e50621f7bfccda66` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `ad8abdc9e2a8ec6f` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `eefbbe6b352816b9` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `de6fba4c1cef016d` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, DIST, DEBUG, FORMAT TEXT` supported on yugabyte
- one host: this run used node budget 2 CPU 3g and client budget 1 CPU 2g, so budget comparisons are separate classes here, not independent hardware
- isolated phase: no concurrent writers in this run; the blended read/write phase is a separate labelled run
- physical rank-index arm yb-range-rank; an engine whose planner ignores the index is a different physical design, not a faster one
- YugabyteDB DIST/DEBUG storage counters (rows and index rows scanned) are parsed per cell; RPC-level counters are not aggregated
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
