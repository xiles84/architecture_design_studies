# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 30 cells, 30 valid, engine `yugabyte`, topology `cluster3`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-v3b-topo-cluster-i1-20261005`
- run tag: `run/study-01/select-join-v3b-topo-cluster-i1-20261005`
- environment: `host-zenbook-ux5406sa`
- repo commit: `7243fee074bfb16fc5ad8b6afc1a2d7549080935`
- engine version: `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit`
- database engine image: `docker.io/yugabytedb/yugabyte:2025.2.6.0-b111` @ `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`
- benchmark client image: `localhost/charitybench-select:select-join-v3b-topo-cluster-i1-20261005` @ `f3f1e10c44cbcb4c7b5bc56ed95a4f071682b6f7d33b5e0a7b688c4271eeef88`
- inputs digest: `c050491b0ec871cd3829f7baba4d313f303db4ed175530f479b6744d17eb9ff7`
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
| 20 (0.01%) | 2.5730 | missing | missing | missing | missing | missing | 2.8159 | 119.9445 | 2.7984 | missing |
| 2000 (1%) | 12.2691 | missing | missing | missing | missing | missing | 12.4619 | 121.8307 | 12.2331 | missing |
| 180000 (90%) | 124.2879 | missing | missing | missing | missing | missing | 106.5822 | 156.2962 | 150.1061 | missing |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | cte_joined_outer vs join_where | 2.5730 | 2.8159 | 1.0944 | false |
| 20 | cte_joined_outer_mat vs join_where | 2.5730 | 119.9445 | 46.6166 | false |
| 20 | cte_joined_outer_notmat vs join_where | 2.5730 | 2.7984 | 1.0876 | false |
| 2000 | cte_joined_outer vs join_where | 12.2691 | 12.4619 | 1.0157 | false |
| 2000 | cte_joined_outer_mat vs join_where | 12.2691 | 121.8307 | 9.9299 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 12.2691 | 12.2331 | 0.9971 | false |
| 180000 | cte_joined_outer vs join_where | 124.2879 | 106.5822 | 0.8575 | false |
| 180000 | cte_joined_outer_mat vs join_where | 124.2879 | 156.2962 | 1.2575 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 124.2879 | 150.1061 | 1.2077 | false |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_joined_outer_mat vs cte_joined_outer | 2.8159 | 119.9445 | 42.5957 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 2.8159 | 2.7984 | 0.9938 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 12.4619 | 121.8307 | 9.7762 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 12.4619 | 12.2331 | 0.9816 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 106.5822 | 156.2962 | 1.4664 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 106.5822 | 150.1061 | 1.4084 | false | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 7.2216 | 1 | `2f17764496c2023b` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 8.2658 | 1 | `08d79acd9722304f` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | missing | | | | |
| default | 2 | 100% | missing | | | | |
| default | 4 | 1% | 9.7563 | 4 | `46ef1f0ee37b563f` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 8.9563 | 4 | `ad79148f426790bc` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | missing | | | | |
| default | 8 | 100% | missing | | | | |
| default | 16 | 1% | 11.4626 | 16 | `7f20479e849399a4` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 10.6587 | 16 | `3c4099ec4704221a` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 10.6110 | 1 | `d2ecf6d28226e256` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 9.9394 | 1 | `29ec3344b188afc9` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | missing | | | | |
| mat | 2 | 100% | missing | | | | |
| mat | 4 | 1% | 9.1470 | 4 | `65eb12a6bbf574a4` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 9.8541 | 4 | `68c22681333fa06e` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | missing | | | | |
| mat | 8 | 100% | missing | | | | |
| mat | 16 | 1% | 9.4616 | 16 | `004315f6f30e6a5d` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 11.2661 | 16 | `4e26a063e8ddffa4` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 8.6892 | 1 | `d95cd0e2a1d3d669` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 11.6640 | 1 | `99268aa99c5bf36b` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | missing | | | | |
| notmat | 2 | 100% | missing | | | | |
| notmat | 4 | 1% | 23.5047 | 4 | `e51993c16800bf40` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 28.4156 | 4 | `966c98473f2b7bfc` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | missing | | | | |
| notmat | 8 | 100% | missing | | | | |
| notmat | 16 | 1% | 80.5107 | 16 | `bd8b0fa4887a640e` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 127.2619 | 16 | `217d4163b8e0d7af` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 1.4693 | false |
| notmat | 1 | 1% | 1.2032 | false |
| mat | 1 | 100% | 1.2025 | false |
| notmat | 1 | 100% | 1.4111 | false |
| mat | 4 | 1% | 0.9376 | false |
| notmat | 4 | 1% | 2.4092 | false |
| mat | 4 | 100% | 1.1002 | false |
| notmat | 4 | 100% | 3.1727 | false |
| mat | 16 | 1% | 0.8254 | false |
| notmat | 16 | 1% | 7.0238 | false |
| mat | 16 | 100% | 1.0570 | false |
| notmat | 16 | 100% | 11.9398 | false |

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

- client: nr_periods 315, nr_throttled 109, throttled_usec 4231433 | quota 100000 100000
, memory 2147483648

- yb-n1: nr_periods 605, nr_throttled 6, throttled_usec 564235 | quota 200000 100000, memory 3221225472
- yb-n2: nr_periods 594, nr_throttled 4, throttled_usec 533493 | quota 200000 100000, memory 3221225472
- yb-n3: nr_periods 590, nr_throttled 15, throttled_usec 650060 | quota 200000 100000, memory 3221225472

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
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 100.7262 | 106.5822 | 231.2283 | 5 | 9.4 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 2.5781 | 2.8159 | 4.6528 | 5 | 355.1 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 10.7402 | 12.4619 | 17.1734 | 5 | 80.2 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 127.3344 | 156.2962 | 190.8477 | 5 | 6.4 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 113.4320 | 119.9445 | 162.2924 | 5 | 8.3 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 114.0092 | 121.8307 | 168.3970 | 5 | 8.2 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 133.0160 | 150.1061 | 262.6174 | 5 | 6.7 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 2.6678 | 2.7984 | 3.5223 | 5 | 357.3 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 10.7138 | 12.2331 | 12.9378 | 5 | 81.7 |
| family=bag;stmt=join_where;threshold=180000 | valid | 104.3501 | 124.2879 | 132.3471 | 5 | 8.0 |
| family=bag;stmt=join_where;threshold=20 | valid | 2.2697 | 2.5730 | 4.1174 | 5 | 388.7 |
| family=bag;stmt=join_where;threshold=2000 | valid | 9.3872 | 12.2691 | 21.2872 | 5 | 81.5 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 10.1221 | 11.4626 | 13.1985 | 5 | 87.2 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 10.3289 | 10.6587 | 17.7106 | 5 | 93.8 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 6.5406 | 7.2216 | 9.6472 | 5 | 138.5 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 7.7005 | 8.2658 | 12.4574 | 5 | 121.0 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 7.7261 | 9.7563 | 10.4585 | 5 | 102.5 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 7.9622 | 8.9563 | 11.9979 | 5 | 111.7 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 8.2958 | 9.4616 | 13.5371 | 5 | 105.7 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 10.4913 | 11.2661 | 17.4451 | 5 | 88.8 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 8.6274 | 10.6110 | 12.1707 | 5 | 94.2 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 9.3525 | 9.9394 | 11.1809 | 5 | 100.6 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 8.9089 | 9.1470 | 12.1357 | 5 | 109.3 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 9.2811 | 9.8541 | 15.2968 | 5 | 101.5 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 75.5990 | 80.5107 | 201.9489 | 5 | 12.4 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 107.9790 | 127.2619 | 138.0552 | 5 | 7.9 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 7.7660 | 8.6892 | 10.1080 | 5 | 115.1 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 9.1027 | 11.6640 | 18.3117 | 5 | 85.7 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 21.0560 | 23.5047 | 39.5069 | 5 | 42.5 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 28.2163 | 28.4156 | 44.9131 | 5 | 35.2 |

## Matched blended workload

- phase: `isolated`; target offered rate: 10 updates/s
- offered 0, committed 0, failed 0 over 0.000 s
- offered rate 0.000/s, completed rate 0.000/s
- per key: no keys touched
- ledger: 0 keys reconciled, 0 mismatch

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42147252`, 5 interleaved rounds
  1. family=reuse;policy=mat;refs=4;consumer=1%
  2. family=reuse;policy=default;refs=16;consumer=100%
  3. family=reuse;policy=default;refs=1;consumer=1%
  4. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  5. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  6. family=bag;stmt=cte_joined_outer;threshold=180000
  7. family=bag;stmt=cte_joined_outer;threshold=2000
  8. family=reuse;policy=mat;refs=4;consumer=100%
  9. family=reuse;policy=default;refs=1;consumer=100%
  10. family=reuse;policy=notmat;refs=16;consumer=1%
  11. family=reuse;policy=mat;refs=16;consumer=100%
  12. family=reuse;policy=notmat;refs=4;consumer=100%
  13. family=bag;stmt=cte_joined_outer;threshold=20
  14. family=bag;stmt=join_where;threshold=20
  15. family=reuse;policy=default;refs=4;consumer=100%
  16. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  17. family=reuse;policy=notmat;refs=1;consumer=1%
  18. family=reuse;policy=mat;refs=16;consumer=1%
  19. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  20. family=bag;stmt=cte_joined_outer_mat;threshold=20
  21. family=reuse;policy=notmat;refs=16;consumer=100%
  22. family=reuse;policy=default;refs=16;consumer=1%
  23. family=reuse;policy=mat;refs=1;consumer=1%
  24. family=reuse;policy=default;refs=4;consumer=1%
  25. family=bag;stmt=join_where;threshold=2000
  26. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  27. family=bag;stmt=join_where;threshold=180000
  28. family=reuse;policy=mat;refs=1;consumer=100%
  29. family=reuse;policy=notmat;refs=1;consumer=100%
  30. family=reuse;policy=notmat;refs=4;consumer=1%

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
| family=bag;stmt=cte_joined_outer;threshold=180000 | j05_cte_joined_outer | valid | `9b53a12d21d71e42` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20 | j05_cte_joined_outer | valid | `bc357258025e6432` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=2000 | j05_cte_joined_outer | valid | `b2fd38e3fa1ba637` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | j05m_cte_joined_outer_mat | valid | `ce4e55b9e7d421d5` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | j05m_cte_joined_outer_mat | valid | `f89103d54a3af992` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | j05m_cte_joined_outer_mat | valid | `bfdc0eae52c2064f` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | j05n_cte_joined_outer_notmat | valid | `a515a24c484328f9` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | j05n_cte_joined_outer_notmat | valid | `1311e1ca09bd75d4` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | j05n_cte_joined_outer_notmat | valid | `dd0817d1af6ebcaf` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=join_where;threshold=180000 | j01_join_where | valid | `2a77e5b73fcd1e02` |  |
| family=bag;stmt=join_where;threshold=20 | j01_join_where | valid | `a6de84e1dcfd6335` |  |
| family=bag;stmt=join_where;threshold=2000 | j01_join_where | valid | `4dff90e8e2372778` |  |
| family=reuse;policy=default;refs=16;consumer=1% | r_default_n16_c1 | valid | `7f20479e849399a4` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=16;consumer=100% | r_default_n16_c100 | valid | `3c4099ec4704221a` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=1;consumer=1% | r_default_n1_c1 | valid | `2f17764496c2023b` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=1;consumer=100% | r_default_n1_c100 | valid | `08d79acd9722304f` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=4;consumer=1% | r_default_n4_c1 | valid | `46ef1f0ee37b563f` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=100% | r_default_n4_c100 | valid | `ad79148f426790bc` | policy=default; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=1% | r_mat_n16_c1 | valid | `004315f6f30e6a5d` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=100% | r_mat_n16_c100 | valid | `4e26a063e8ddffa4` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=1% | r_mat_n1_c1 | valid | `d2ecf6d28226e256` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=100% | r_mat_n1_c100 | valid | `29ec3344b188afc9` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=1% | r_mat_n4_c1 | valid | `65eb12a6bbf574a4` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=100% | r_mat_n4_c100 | valid | `68c22681333fa06e` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=notmat;refs=16;consumer=1% | r_notmat_n16_c1 | valid | `bd8b0fa4887a640e` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `217d4163b8e0d7af` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `d95cd0e2a1d3d669` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `99268aa99c5bf36b` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `e51993c16800bf40` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `966c98473f2b7bfc` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, DIST, DEBUG, FORMAT TEXT` supported on yugabyte
- one host: this run used node budget 2 CPU 3g and client budget 1 CPU 2g, so budget comparisons are separate classes here, not independent hardware
- isolated phase: no concurrent writers in this run; the blended read/write phase is a separate labelled run
- cluster nodes are containers on the same host, so these are not network or independent-hardware measurements; container colocation is not data-colocation evidence
- readers round-robin across 3 endpoints; the writer's endpoint policy is recorded separately
- physical rank-index arm yb-range-rank; an engine whose planner ignores the index is a different physical design, not a faster one
- YugabyteDB DIST/DEBUG storage counters (rows and index rows scanned) are parsed per cell; RPC-level counters are not aggregated
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
