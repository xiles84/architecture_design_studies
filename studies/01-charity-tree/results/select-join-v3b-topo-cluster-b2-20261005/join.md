# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 30 cells, 30 valid, engine `yugabyte`, topology `cluster3`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-v3b-topo-cluster-b2-20261005`
- run tag: `run/study-01/select-join-v3b-topo-cluster-b2-20261005`
- environment: `host-zenbook-ux5406sa`
- repo commit: `7243fee074bfb16fc5ad8b6afc1a2d7549080935`
- engine version: `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit`
- database engine image: `docker.io/yugabytedb/yugabyte:2025.2.6.0-b111` @ `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`
- benchmark client image: `localhost/charitybench-select:select-join-v3b-topo-cluster-b2-20261005` @ `f3f1e10c44cbcb4c7b5bc56ed95a4f071682b6f7d33b5e0a7b688c4271eeef88`
- inputs digest: `9eb4ca553a7a1f7d9e82b64c1fa959aa3c633ac8b770d57fd5a4a82f0f60a7d6`
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
- setting `phase`: blended
- setting `placement`: noncolocated
- setting `prepared_plan_policy`: pgx simple protocol; no prepared statements
- setting `scale`: cte-domain
- setting `selectivity_mapping`: rank <= threshold; rank is a deterministic shuffle of 1..fact_count
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
| yb-n1 | 168 | 2 |
| yb-n2 | 168 | 2 |
| yb-n3 | 168 | 2 |

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
| 20 (0.01%) | 3.1531 | missing | missing | missing | missing | missing | 2.0955 | 132.4396 | 2.0223 | missing |
| 2000 (1%) | 11.4720 | missing | missing | missing | missing | missing | 10.7292 | 139.7980 | 10.5862 | missing |
| 180000 (90%) | 126.9917 | missing | missing | missing | missing | missing | 163.2515 | 159.5452 | 128.5588 | missing |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | cte_joined_outer vs join_where | 3.1531 | 2.0955 | 0.6646 | false |
| 20 | cte_joined_outer_mat vs join_where | 3.1531 | 132.4396 | 42.0027 | false |
| 20 | cte_joined_outer_notmat vs join_where | 3.1531 | 2.0223 | 0.6414 | false |
| 2000 | cte_joined_outer vs join_where | 11.4720 | 10.7292 | 0.9353 | false |
| 2000 | cte_joined_outer_mat vs join_where | 11.4720 | 139.7980 | 12.1860 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 11.4720 | 10.5862 | 0.9228 | false |
| 180000 | cte_joined_outer vs join_where | 126.9917 | 163.2515 | 1.2855 | false |
| 180000 | cte_joined_outer_mat vs join_where | 126.9917 | 159.5452 | 1.2563 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 126.9917 | 128.5588 | 1.0123 | false |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_joined_outer_mat vs cte_joined_outer | 2.0955 | 132.4396 | 63.2008 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 2.0955 | 2.0223 | 0.9651 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 10.7292 | 139.7980 | 13.0297 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 10.7292 | 10.5862 | 0.9867 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 163.2515 | 159.5452 | 0.9773 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 163.2515 | 128.5588 | 0.7875 | false | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 8.4139 | 1 | `34b6b025d25b90f8` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 9.9445 | 1 | `4f3d943dcc24e22c` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | missing | | | | |
| default | 2 | 100% | missing | | | | |
| default | 4 | 1% | 9.2930 | 4 | `ba3a329f08b399d8` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 11.2436 | 4 | `4dac6aa63fc6dcce` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | missing | | | | |
| default | 8 | 100% | missing | | | | |
| default | 16 | 1% | 10.5318 | 16 | `339958a953ee694f` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 11.8938 | 16 | `b81c9dd524de9c9a` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 10.2495 | 1 | `5a39baaf28f01405` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 10.2450 | 1 | `18d068775cf2cb58` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | missing | | | | |
| mat | 2 | 100% | missing | | | | |
| mat | 4 | 1% | 9.9866 | 4 | `778f4e0ce461debb` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 9.4559 | 4 | `3c5d2f7f3eeaea4e` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | missing | | | | |
| mat | 8 | 100% | missing | | | | |
| mat | 16 | 1% | 11.2256 | 16 | `8e630d466de7ba60` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 11.0195 | 16 | `4127af03c54c7930` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 7.9312 | 1 | `c3412cd35361ff66` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 9.7422 | 1 | `02112ce9050b806e` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | missing | | | | |
| notmat | 2 | 100% | missing | | | | |
| notmat | 4 | 1% | 31.2613 | 4 | `8f9b58cbf8294499` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 38.2764 | 4 | `39b412bda7f67090` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | missing | | | | |
| notmat | 8 | 100% | missing | | | | |
| notmat | 16 | 1% | 94.9560 | 16 | `29c35b8670d6aea6` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 105.4313 | 16 | `dccb6645ba4912dc` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 1.2182 | false |
| notmat | 1 | 1% | 0.9426 | false |
| mat | 1 | 100% | 1.0302 | false |
| notmat | 1 | 100% | 0.9797 | false |
| mat | 4 | 1% | 1.0746 | false |
| notmat | 4 | 1% | 3.3640 | false |
| mat | 4 | 100% | 0.8410 | false |
| notmat | 4 | 100% | 3.4043 | false |
| mat | 16 | 1% | 1.0659 | false |
| notmat | 16 | 1% | 9.0161 | false |
| mat | 16 | 100% | 0.9265 | false |
| notmat | 16 | 100% | 8.8644 | false |

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

- client: nr_periods 348, nr_throttled 133, throttled_usec 5278383 | quota 100000 100000
, memory 2147483648

- yb-n1: nr_periods 630, nr_throttled 4, throttled_usec 689694 | quota 200000 100000, memory 3221225472
- yb-n2: nr_periods 624, nr_throttled 6, throttled_usec 551110 | quota 200000 100000, memory 3221225472
- yb-n3: nr_periods 619, nr_throttled 6, throttled_usec 65130 | quota 200000 100000, memory 3221225472

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
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 155.5660 | 163.2515 | 166.3501 | 5 | 6.1 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 1.9589 | 2.0955 | 2.5485 | 5 | 477.2 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 10.2946 | 10.7292 | 10.9154 | 5 | 93.2 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 150.0146 | 159.5452 | 165.4730 | 5 | 6.3 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 131.3888 | 132.4396 | 137.8584 | 5 | 7.6 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 133.1418 | 139.7980 | 158.5585 | 5 | 7.2 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 128.3931 | 128.5588 | 132.5427 | 5 | 7.8 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 1.9304 | 2.0223 | 2.3864 | 5 | 494.5 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 9.8956 | 10.5862 | 11.3565 | 5 | 94.5 |
| family=bag;stmt=join_where;threshold=180000 | valid | 118.6147 | 126.9917 | 147.8834 | 5 | 7.9 |
| family=bag;stmt=join_where;threshold=20 | valid | 2.7139 | 3.1531 | 3.4714 | 5 | 317.1 |
| family=bag;stmt=join_where;threshold=2000 | valid | 11.1133 | 11.4720 | 11.9169 | 5 | 87.2 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 10.3746 | 10.5318 | 10.8929 | 5 | 95.0 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 11.3062 | 11.8938 | 13.2020 | 5 | 84.1 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 8.2223 | 8.4139 | 9.5298 | 5 | 118.9 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 9.8914 | 9.9445 | 11.0862 | 5 | 100.6 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 8.7882 | 9.2930 | 10.1861 | 5 | 107.6 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 10.7552 | 11.2436 | 11.8668 | 5 | 88.9 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 10.8679 | 11.2256 | 13.6081 | 5 | 89.1 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 10.5669 | 11.0195 | 11.3564 | 5 | 90.7 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 10.1544 | 10.2495 | 10.5282 | 5 | 97.6 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 9.6887 | 10.2450 | 11.1627 | 5 | 97.6 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 9.7667 | 9.9866 | 11.5859 | 5 | 100.1 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 8.9134 | 9.4559 | 9.6119 | 5 | 105.8 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 84.0482 | 94.9560 | 98.1479 | 5 | 10.5 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 98.7947 | 105.4313 | 107.7036 | 5 | 9.5 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 7.5764 | 7.9312 | 8.1699 | 5 | 126.1 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 9.2916 | 9.7422 | 10.5265 | 5 | 102.6 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 28.8380 | 31.2613 | 34.0165 | 5 | 32.0 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 36.3776 | 38.2764 | 39.8076 | 5 | 26.1 |

## Matched blended workload

- phase: `blended`; target offered rate: 10 updates/s
- offered 175, committed 175, failed 0 over 17.455 s
- offered rate 10.026/s, completed rate 10.026/s
- writer operations per endpoint: [59 58 58]
- reader operations overlapping a committed write: 75 of 210
- reader retries of transient conflicts: 0
- per key: 175 keys, min 1, max 1 per key
- ledger: 175 keys reconciled, 0 mismatch
- update latency: avg 7.077 ms, min 2.456, max 188.129 over 175

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42529834`, 5 interleaved rounds
  1. family=reuse;policy=mat;refs=1;consumer=1%
  2. family=bag;stmt=cte_joined_outer;threshold=180000
  3. family=bag;stmt=cte_joined_outer;threshold=2000
  4. family=reuse;policy=default;refs=16;consumer=100%
  5. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  6. family=bag;stmt=join_where;threshold=180000
  7. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  8. family=bag;stmt=join_where;threshold=20
  9. family=reuse;policy=default;refs=16;consumer=1%
  10. family=reuse;policy=notmat;refs=1;consumer=100%
  11. family=reuse;policy=notmat;refs=1;consumer=1%
  12. family=bag;stmt=cte_joined_outer;threshold=20
  13. family=bag;stmt=cte_joined_outer_mat;threshold=20
  14. family=reuse;policy=default;refs=1;consumer=100%
  15. family=reuse;policy=default;refs=4;consumer=1%
  16. family=reuse;policy=notmat;refs=4;consumer=1%
  17. family=reuse;policy=mat;refs=4;consumer=1%
  18. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  19. family=bag;stmt=join_where;threshold=2000
  20. family=reuse;policy=notmat;refs=16;consumer=100%
  21. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  22. family=reuse;policy=notmat;refs=16;consumer=1%
  23. family=reuse;policy=default;refs=4;consumer=100%
  24. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  25. family=reuse;policy=default;refs=1;consumer=1%
  26. family=reuse;policy=mat;refs=16;consumer=1%
  27. family=reuse;policy=mat;refs=4;consumer=100%
  28. family=reuse;policy=notmat;refs=4;consumer=100%
  29. family=reuse;policy=mat;refs=1;consumer=100%
  30. family=reuse;policy=mat;refs=16;consumer=100%

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
| family=bag;stmt=cte_joined_outer;threshold=180000 | j05_cte_joined_outer | valid | `30d4f22a78aac9db` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20 | j05_cte_joined_outer | valid | `d1a270b8c36c3320` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=2000 | j05_cte_joined_outer | valid | `ca7d855e379c7a2c` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | j05m_cte_joined_outer_mat | valid | `d82c34d0b0b083c5` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | j05m_cte_joined_outer_mat | valid | `834b852ead573a19` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | j05m_cte_joined_outer_mat | valid | `ceb1eb4739610a34` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | j05n_cte_joined_outer_notmat | valid | `ec29c290f2d9ad35` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | j05n_cte_joined_outer_notmat | valid | `133c20a05028be17` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | j05n_cte_joined_outer_notmat | valid | `18071a49db3e24a8` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=join_where;threshold=180000 | j01_join_where | valid | `0815ebf86b6eaa35` |  |
| family=bag;stmt=join_where;threshold=20 | j01_join_where | valid | `d6b3b9d6df2a2c9b` |  |
| family=bag;stmt=join_where;threshold=2000 | j01_join_where | valid | `2abfb45434a8e005` |  |
| family=reuse;policy=default;refs=16;consumer=1% | r_default_n16_c1 | valid | `339958a953ee694f` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=16;consumer=100% | r_default_n16_c100 | valid | `b81c9dd524de9c9a` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=1;consumer=1% | r_default_n1_c1 | valid | `34b6b025d25b90f8` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=1;consumer=100% | r_default_n1_c100 | valid | `4f3d943dcc24e22c` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=4;consumer=1% | r_default_n4_c1 | valid | `ba3a329f08b399d8` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=100% | r_default_n4_c100 | valid | `4dac6aa63fc6dcce` | policy=default; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=1% | r_mat_n16_c1 | valid | `8e630d466de7ba60` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=100% | r_mat_n16_c100 | valid | `4127af03c54c7930` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=1% | r_mat_n1_c1 | valid | `5a39baaf28f01405` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=100% | r_mat_n1_c100 | valid | `18d068775cf2cb58` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=1% | r_mat_n4_c1 | valid | `778f4e0ce461debb` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=100% | r_mat_n4_c100 | valid | `3c5d2f7f3eeaea4e` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=notmat;refs=16;consumer=1% | r_notmat_n16_c1 | valid | `29c35b8670d6aea6` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `dccb6645ba4912dc` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `c3412cd35361ff66` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `02112ce9050b806e` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `8f9b58cbf8294499` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `39b412bda7f67090` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, DIST, DEBUG, FORMAT TEXT` supported on yugabyte
- one host: this run used node budget 2 CPU 3g and client budget 1 CPU 2g, so budget comparisons are separate classes here, not independent hardware
- one client process with one fixed maintenance writer at 10 updates/s; readers are sequential; the reported sample latency includes retry and backoff time
- cluster nodes are containers on the same host, so these are not network or independent-hardware measurements; container colocation is not data-colocation evidence
- readers round-robin across 3 endpoints; the writer's endpoint policy is recorded separately
- physical rank-index arm yb-range-rank; an engine whose planner ignores the index is a different physical design, not a faster one
- YugabyteDB DIST/DEBUG storage counters (rows and index rows scanned) are parsed per cell; RPC-level counters are not aggregated
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
