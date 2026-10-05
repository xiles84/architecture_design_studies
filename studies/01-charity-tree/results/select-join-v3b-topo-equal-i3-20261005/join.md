# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 30 cells, 30 valid, engine `yugabyte`, topology `single`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-v3b-topo-equal-i3-20261005`
- run tag: `run/study-01/select-join-v3b-topo-equal-i3-20261005`
- environment: `host-zenbook-ux5406sa`
- repo commit: `7243fee074bfb16fc5ad8b6afc1a2d7549080935`
- engine version: `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit`
- database engine image: `docker.io/yugabytedb/yugabyte:2025.2.6.0-b111` @ `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`
- benchmark client image: `localhost/charitybench-select:select-join-v3b-topo-equal-i3-20261005` @ `f3f1e10c44cbcb4c7b5bc56ed95a4f071682b6f7d33b5e0a7b688c4271eeef88`
- inputs digest: `7d5d9b2325cefc32de2cc09f01557ab2a0ae7fefbda3eff894d5b0e99e2f9074`
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
| 20 (0.01%) | 1.8455 | missing | missing | missing | missing | missing | 1.8490 | 81.7325 | 1.5638 | missing |
| 2000 (1%) | 8.7323 | missing | missing | missing | missing | missing | 9.5234 | 78.1647 | 9.0075 | missing |
| 180000 (90%) | 101.6297 | missing | missing | missing | missing | missing | 100.2009 | 104.1211 | 104.1177 | missing |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | cte_joined_outer vs join_where | 1.8455 | 1.8490 | 1.0019 | false |
| 20 | cte_joined_outer_mat vs join_where | 1.8455 | 81.7325 | 44.2871 | false |
| 20 | cte_joined_outer_notmat vs join_where | 1.8455 | 1.5638 | 0.8473 | false |
| 2000 | cte_joined_outer vs join_where | 8.7323 | 9.5234 | 1.0906 | false |
| 2000 | cte_joined_outer_mat vs join_where | 8.7323 | 78.1647 | 8.9512 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 8.7323 | 9.0075 | 1.0315 | false |
| 180000 | cte_joined_outer vs join_where | 101.6297 | 100.2009 | 0.9859 | false |
| 180000 | cte_joined_outer_mat vs join_where | 101.6297 | 104.1211 | 1.0245 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 101.6297 | 104.1177 | 1.0245 | false |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_joined_outer_mat vs cte_joined_outer | 1.8490 | 81.7325 | 44.2032 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 1.8490 | 1.5638 | 0.8457 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 9.5234 | 78.1647 | 8.2076 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 9.5234 | 9.0075 | 0.9458 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 100.2009 | 104.1211 | 1.0391 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 100.2009 | 104.1177 | 1.0391 | false | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 6.8676 | 1 | `3dc80a1388cff719` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 8.9415 | 1 | `f68da5ef71a7cc11` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | missing | | | | |
| default | 2 | 100% | missing | | | | |
| default | 4 | 1% | 8.7095 | 4 | `d8e41b76f91da919` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 8.5290 | 4 | `5cdbf962b365a9f1` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | missing | | | | |
| default | 8 | 100% | missing | | | | |
| default | 16 | 1% | 9.9309 | 16 | `6816e14972da4dcf` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 9.5039 | 16 | `d0ffe682e83055c9` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 7.9369 | 1 | `a3d4159ec16201a7` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 7.8878 | 1 | `5e4d073109ad90cc` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | missing | | | | |
| mat | 2 | 100% | missing | | | | |
| mat | 4 | 1% | 8.3203 | 4 | `292bd75f5dbfa2b8` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 8.1596 | 4 | `634aeb10b03573f1` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | missing | | | | |
| mat | 8 | 100% | missing | | | | |
| mat | 16 | 1% | 9.0820 | 16 | `57e1fc876e25c0f2` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 9.5326 | 16 | `033fef38664dcbb2` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 7.5141 | 1 | `46583b5b4611965e` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 7.8245 | 1 | `dd5d98868ed4b380` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | missing | | | | |
| notmat | 2 | 100% | missing | | | | |
| notmat | 4 | 1% | 22.4908 | 4 | `c3697a3e0caf64d2` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 28.6019 | 4 | `2c242fdd9d0bce74` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | missing | | | | |
| notmat | 8 | 100% | missing | | | | |
| notmat | 16 | 1% | 70.9807 | 16 | `56390057e90ed23a` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 78.5057 | 16 | `a18b775d890dc28f` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 1.1557 | false |
| notmat | 1 | 1% | 1.0941 | false |
| mat | 1 | 100% | 0.8822 | false |
| notmat | 1 | 100% | 0.8751 | false |
| mat | 4 | 1% | 0.9553 | false |
| notmat | 4 | 1% | 2.5823 | false |
| mat | 4 | 100% | 0.9567 | false |
| notmat | 4 | 100% | 3.3535 | false |
| mat | 16 | 1% | 0.9145 | false |
| notmat | 16 | 1% | 7.1474 | false |
| mat | 16 | 100% | 1.0030 | false |
| notmat | 16 | 100% | 8.2604 | false |

## Physical arm

- `charity_pkey`: `CREATE UNIQUE INDEX charity_pkey ON selectjoin_domain.charity USING lsm (charity_id HASH)`
- `donation_person_idx`: `CREATE INDEX donation_person_idx ON selectjoin_domain.donation USING lsm (person_id HASH)`
- `donation_pkey`: `CREATE UNIQUE INDEX donation_pkey ON selectjoin_domain.donation USING lsm (donation_id HASH)`
- `donation_rank_idx`: `CREATE INDEX donation_rank_idx ON selectjoin_domain.donation USING lsm (rank ASC)`
- `person_pkey`: `CREATE UNIQUE INDEX person_pkey ON selectjoin_domain.person USING lsm (person_id HASH)`

- index arm `yb-range-rank`; placement `noncolocated`; node budget `6 CPU 9g`; client budget `1 CPU 2g`

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

- client: nr_periods 278, nr_throttled 123, throttled_usec 5104909 | quota 100000 100000
, memory 2147483648

- yb-single: nr_periods 348, nr_throttled 0, throttled_usec 0 | quota 600000 100000, memory 9663676416

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
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 99.5479 | 100.2009 | 107.3408 | 5 | 10.0 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 1.7595 | 1.8490 | 1.8950 | 5 | 540.8 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 8.9848 | 9.5234 | 9.8002 | 5 | 105.0 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 93.4102 | 104.1211 | 113.3163 | 5 | 9.6 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 79.0957 | 81.7325 | 83.1762 | 5 | 12.2 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 74.2336 | 78.1647 | 82.9827 | 5 | 12.8 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 100.8454 | 104.1177 | 105.3098 | 5 | 9.6 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 1.5238 | 1.5638 | 1.8650 | 5 | 639.5 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 8.5852 | 9.0075 | 9.2508 | 5 | 111.0 |
| family=bag;stmt=join_where;threshold=180000 | valid | 98.2645 | 101.6297 | 109.6620 | 5 | 9.8 |
| family=bag;stmt=join_where;threshold=20 | valid | 1.6642 | 1.8455 | 2.1644 | 5 | 541.9 |
| family=bag;stmt=join_where;threshold=2000 | valid | 8.3045 | 8.7323 | 29.3358 | 5 | 114.5 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 9.1664 | 9.9309 | 10.5549 | 5 | 100.7 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 8.8801 | 9.5039 | 10.0600 | 5 | 105.2 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 6.6638 | 6.8676 | 7.5429 | 5 | 145.6 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 8.6252 | 8.9415 | 9.4714 | 5 | 111.8 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 8.3211 | 8.7095 | 9.2595 | 5 | 114.8 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 7.7812 | 8.5290 | 10.4841 | 5 | 117.2 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 8.7860 | 9.0820 | 9.5056 | 5 | 110.1 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 9.3831 | 9.5326 | 9.9220 | 5 | 104.9 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 7.7835 | 7.9369 | 8.7889 | 5 | 126.0 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 7.4076 | 7.8878 | 8.5244 | 5 | 126.8 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 7.7491 | 8.3203 | 8.8923 | 5 | 120.2 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 8.0161 | 8.1596 | 8.7897 | 5 | 122.6 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 70.1633 | 70.9807 | 76.7722 | 5 | 14.1 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 76.6398 | 78.5057 | 85.3234 | 5 | 12.7 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 6.2444 | 7.5141 | 8.1333 | 5 | 133.1 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 7.4360 | 7.8245 | 10.9016 | 5 | 127.8 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 20.7341 | 22.4908 | 23.8225 | 5 | 44.5 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 27.6202 | 28.6019 | 31.6555 | 5 | 35.0 |

## Matched blended workload

- phase: `isolated`; target offered rate: 10 updates/s
- offered 0, committed 0, failed 0 over 0.000 s
- offered rate 0.000/s, completed rate 0.000/s
- per key: no keys touched
- ledger: 0 keys reconciled, 0 mismatch

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42214560`, 5 interleaved rounds
  1. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  2. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  3. family=reuse;policy=default;refs=4;consumer=100%
  4. family=bag;stmt=cte_joined_outer;threshold=180000
  5. family=reuse;policy=default;refs=16;consumer=1%
  6. family=bag;stmt=join_where;threshold=2000
  7. family=reuse;policy=mat;refs=16;consumer=1%
  8. family=reuse;policy=mat;refs=1;consumer=100%
  9. family=reuse;policy=mat;refs=4;consumer=1%
  10. family=bag;stmt=join_where;threshold=180000
  11. family=reuse;policy=notmat;refs=16;consumer=1%
  12. family=reuse;policy=default;refs=4;consumer=1%
  13. family=reuse;policy=mat;refs=16;consumer=100%
  14. family=reuse;policy=notmat;refs=4;consumer=100%
  15. family=bag;stmt=cte_joined_outer;threshold=20
  16. family=reuse;policy=notmat;refs=16;consumer=100%
  17. family=reuse;policy=default;refs=1;consumer=1%
  18. family=bag;stmt=join_where;threshold=20
  19. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  20. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  21. family=bag;stmt=cte_joined_outer;threshold=2000
  22. family=reuse;policy=mat;refs=4;consumer=100%
  23. family=reuse;policy=notmat;refs=4;consumer=1%
  24. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  25. family=reuse;policy=default;refs=1;consumer=100%
  26. family=reuse;policy=notmat;refs=1;consumer=1%
  27. family=reuse;policy=notmat;refs=1;consumer=100%
  28. family=reuse;policy=default;refs=16;consumer=100%
  29. family=bag;stmt=cte_joined_outer_mat;threshold=20
  30. family=reuse;policy=mat;refs=1;consumer=1%

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
| family=bag;stmt=cte_joined_outer;threshold=180000 | j05_cte_joined_outer | valid | `3eb94342c066642d` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20 | j05_cte_joined_outer | valid | `67ba6156f0e1a040` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=2000 | j05_cte_joined_outer | valid | `1f57d4d1c89bcc67` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | j05m_cte_joined_outer_mat | valid | `2e94f9010c3c55c5` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | j05m_cte_joined_outer_mat | valid | `c13d1781990942be` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | j05m_cte_joined_outer_mat | valid | `28e7bf3b55a38f66` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | j05n_cte_joined_outer_notmat | valid | `592a86cf631707dc` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | j05n_cte_joined_outer_notmat | valid | `76b71bc24115e5c8` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | j05n_cte_joined_outer_notmat | valid | `bb5bb7567a50e15c` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=join_where;threshold=180000 | j01_join_where | valid | `feb5338dd797ca5d` |  |
| family=bag;stmt=join_where;threshold=20 | j01_join_where | valid | `1429bff8ee51f247` |  |
| family=bag;stmt=join_where;threshold=2000 | j01_join_where | valid | `0a8c1c60c2e83d32` |  |
| family=reuse;policy=default;refs=16;consumer=1% | r_default_n16_c1 | valid | `6816e14972da4dcf` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=16;consumer=100% | r_default_n16_c100 | valid | `d0ffe682e83055c9` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=1;consumer=1% | r_default_n1_c1 | valid | `3dc80a1388cff719` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=1;consumer=100% | r_default_n1_c100 | valid | `f68da5ef71a7cc11` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=4;consumer=1% | r_default_n4_c1 | valid | `d8e41b76f91da919` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=100% | r_default_n4_c100 | valid | `5cdbf962b365a9f1` | policy=default; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=1% | r_mat_n16_c1 | valid | `57e1fc876e25c0f2` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=100% | r_mat_n16_c100 | valid | `033fef38664dcbb2` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=1% | r_mat_n1_c1 | valid | `a3d4159ec16201a7` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=100% | r_mat_n1_c100 | valid | `5e4d073109ad90cc` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=1% | r_mat_n4_c1 | valid | `292bd75f5dbfa2b8` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=100% | r_mat_n4_c100 | valid | `634aeb10b03573f1` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=notmat;refs=16;consumer=1% | r_notmat_n16_c1 | valid | `56390057e90ed23a` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `a18b775d890dc28f` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `46583b5b4611965e` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `dd5d98868ed4b380` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `c3697a3e0caf64d2` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `2c242fdd9d0bce74` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, DIST, DEBUG, FORMAT TEXT` supported on yugabyte
- one host: this run used node budget 6 CPU 9g and client budget 1 CPU 2g, so budget comparisons are separate classes here, not independent hardware
- isolated phase: no concurrent writers in this run; the blended read/write phase is a separate labelled run
- physical rank-index arm yb-range-rank; an engine whose planner ignores the index is a different physical design, not a faster one
- YugabyteDB DIST/DEBUG storage counters (rows and index rows scanned) are parsed per cell; RPC-level counters are not aggregated
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
