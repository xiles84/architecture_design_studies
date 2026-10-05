# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 30 cells, 30 valid, engine `yugabyte`, topology `cluster3`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-v3b-topo-cluster-i2-20261005`
- run tag: `run/study-01/select-join-v3b-topo-cluster-i2-20261005`
- environment: `host-zenbook-ux5406sa`
- repo commit: `7243fee074bfb16fc5ad8b6afc1a2d7549080935`
- engine version: `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit`
- database engine image: `docker.io/yugabytedb/yugabyte:2025.2.6.0-b111` @ `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`
- benchmark client image: `localhost/charitybench-select:select-join-v3b-topo-cluster-i2-20261005` @ `f3f1e10c44cbcb4c7b5bc56ed95a4f071682b6f7d33b5e0a7b688c4271eeef88`
- inputs digest: `ebacdede2d5ecc2ec598e09048566cbfa8350ad279d33bca76fa29c4270fd9ad`
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
| 20 (0.01%) | 3.2214 | missing | missing | missing | missing | missing | 2.1342 | 146.1342 | 3.0783 | missing |
| 2000 (1%) | 11.2850 | missing | missing | missing | missing | missing | 10.0756 | 141.3146 | 12.0708 | missing |
| 180000 (90%) | 168.3196 | missing | missing | missing | missing | missing | 168.7502 | 152.8720 | 120.7980 | missing |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | cte_joined_outer vs join_where | 3.2214 | 2.1342 | 0.6625 | false |
| 20 | cte_joined_outer_mat vs join_where | 3.2214 | 146.1342 | 45.3631 | false |
| 20 | cte_joined_outer_notmat vs join_where | 3.2214 | 3.0783 | 0.9556 | false |
| 2000 | cte_joined_outer vs join_where | 11.2850 | 10.0756 | 0.8928 | false |
| 2000 | cte_joined_outer_mat vs join_where | 11.2850 | 141.3146 | 12.5223 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 11.2850 | 12.0708 | 1.0696 | false |
| 180000 | cte_joined_outer vs join_where | 168.3196 | 168.7502 | 1.0026 | false |
| 180000 | cte_joined_outer_mat vs join_where | 168.3196 | 152.8720 | 0.9082 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 168.3196 | 120.7980 | 0.7177 | false |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_joined_outer_mat vs cte_joined_outer | 2.1342 | 146.1342 | 68.4725 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 2.1342 | 3.0783 | 1.4424 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 10.0756 | 141.3146 | 14.0254 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 10.0756 | 12.0708 | 1.1980 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 168.7502 | 152.8720 | 0.9059 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 168.7502 | 120.7980 | 0.7158 | false | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 7.6610 | 1 | `0e33c030da3b9c9d` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 10.8564 | 1 | `d401822e9c9c82ca` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | missing | | | | |
| default | 2 | 100% | missing | | | | |
| default | 4 | 1% | 11.3636 | 4 | `5af43fc3681c9afc` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 9.7799 | 4 | `40caa23b129cba52` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | missing | | | | |
| default | 8 | 100% | missing | | | | |
| default | 16 | 1% | 11.4274 | 16 | `5d6e5df01bdc3adb` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 12.8053 | 16 | `3b3fde755d53868b` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 8.7797 | 1 | `800e08adddffb0ee` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 10.5959 | 1 | `a2880e4afc98d41d` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | missing | | | | |
| mat | 2 | 100% | missing | | | | |
| mat | 4 | 1% | 11.0561 | 4 | `28d2bc7359897193` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 10.8504 | 4 | `3b0ae1c80cfb6d91` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | missing | | | | |
| mat | 8 | 100% | missing | | | | |
| mat | 16 | 1% | 12.1679 | 16 | `12dd2f7e44e464d3` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 12.2825 | 16 | `0692eb44eb6c9de9` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 8.5162 | 1 | `2557ac391eece573` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 10.0453 | 1 | `3ef24b2c867afbf4` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | missing | | | | |
| notmat | 2 | 100% | missing | | | | |
| notmat | 4 | 1% | 26.6072 | 4 | `b38d20192dcf7277` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 38.5641 | 4 | `ee96e5744e02c459` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | missing | | | | |
| notmat | 8 | 100% | missing | | | | |
| notmat | 16 | 1% | 92.6179 | 16 | `dc0fb63d7cc06a37` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 103.2266 | 16 | `88497375a98ee09c` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 1.1460 | false |
| notmat | 1 | 1% | 1.1116 | false |
| mat | 1 | 100% | 0.9760 | false |
| notmat | 1 | 100% | 0.9253 | false |
| mat | 4 | 1% | 0.9729 | false |
| notmat | 4 | 1% | 2.3414 | false |
| mat | 4 | 100% | 1.1095 | false |
| notmat | 4 | 100% | 3.9432 | false |
| mat | 16 | 1% | 1.0648 | false |
| notmat | 16 | 1% | 8.1049 | false |
| mat | 16 | 100% | 0.9592 | false |
| notmat | 16 | 100% | 8.0613 | false |

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

- client: nr_periods 357, nr_throttled 134, throttled_usec 6738599 | quota 100000 100000
, memory 2147483648

- yb-n1: nr_periods 644, nr_throttled 9, throttled_usec 756973 | quota 200000 100000, memory 3221225472
- yb-n2: nr_periods 637, nr_throttled 7, throttled_usec 791742 | quota 200000 100000, memory 3221225472
- yb-n3: nr_periods 633, nr_throttled 12, throttled_usec 228304 | quota 200000 100000, memory 3221225472

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
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 152.3184 | 168.7502 | 182.0982 | 5 | 5.9 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 1.8003 | 2.1342 | 2.3876 | 5 | 468.6 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 9.0461 | 10.0756 | 11.1541 | 5 | 99.2 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 141.1188 | 152.8720 | 180.9701 | 5 | 6.5 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 130.1188 | 146.1342 | 173.1848 | 5 | 6.8 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 121.9132 | 141.3146 | 154.8658 | 5 | 7.1 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 116.6124 | 120.7980 | 141.8977 | 5 | 8.3 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 2.7512 | 3.0783 | 3.4526 | 5 | 324.9 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 10.8747 | 12.0708 | 14.3678 | 5 | 82.8 |
| family=bag;stmt=join_where;threshold=180000 | valid | 149.5476 | 168.3196 | 191.2051 | 5 | 5.9 |
| family=bag;stmt=join_where;threshold=20 | valid | 3.1652 | 3.2214 | 3.6765 | 5 | 310.4 |
| family=bag;stmt=join_where;threshold=2000 | valid | 10.0596 | 11.2850 | 12.1312 | 5 | 88.6 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 11.1274 | 11.4274 | 13.1945 | 5 | 87.5 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 11.7371 | 12.8053 | 14.3150 | 5 | 78.1 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 6.5706 | 7.6610 | 8.0300 | 5 | 130.5 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 9.5014 | 10.8564 | 13.2416 | 5 | 92.1 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 9.8146 | 11.3636 | 12.6481 | 5 | 88.0 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 8.6211 | 9.7799 | 10.5261 | 5 | 102.3 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 10.8663 | 12.1679 | 14.7052 | 5 | 82.2 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 11.2827 | 12.2825 | 14.6059 | 5 | 81.4 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 8.5025 | 8.7797 | 10.6073 | 5 | 113.9 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 9.6020 | 10.5959 | 11.0744 | 5 | 94.4 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 9.9803 | 11.0561 | 11.9570 | 5 | 90.4 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 9.6745 | 10.8504 | 12.4851 | 5 | 92.2 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 84.5478 | 92.6179 | 110.7394 | 5 | 10.8 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 88.7020 | 103.2266 | 118.5358 | 5 | 9.7 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 7.4926 | 8.5162 | 9.5055 | 5 | 117.4 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 8.7679 | 10.0453 | 11.6821 | 5 | 99.5 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 21.8948 | 26.6072 | 31.2918 | 5 | 37.6 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 33.9235 | 38.5641 | 43.1657 | 5 | 25.9 |

## Matched blended workload

- phase: `isolated`; target offered rate: 10 updates/s
- offered 0, committed 0, failed 0 over 0.000 s
- offered rate 0.000/s, completed rate 0.000/s
- per key: no keys touched
- ledger: 0 keys reconciled, 0 mismatch

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42207287`, 5 interleaved rounds
  1. family=bag;stmt=cte_joined_outer_mat;threshold=20
  2. family=bag;stmt=join_where;threshold=180000
  3. family=bag;stmt=join_where;threshold=2000
  4. family=reuse;policy=mat;refs=16;consumer=1%
  5. family=bag;stmt=cte_joined_outer;threshold=180000
  6. family=reuse;policy=notmat;refs=4;consumer=1%
  7. family=reuse;policy=mat;refs=1;consumer=100%
  8. family=reuse;policy=notmat;refs=1;consumer=1%
  9. family=reuse;policy=default;refs=1;consumer=1%
  10. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  11. family=reuse;policy=mat;refs=4;consumer=1%
  12. family=reuse;policy=default;refs=4;consumer=100%
  13. family=reuse;policy=notmat;refs=4;consumer=100%
  14. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  15. family=bag;stmt=cte_joined_outer;threshold=20
  16. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  17. family=reuse;policy=mat;refs=4;consumer=100%
  18. family=bag;stmt=cte_joined_outer;threshold=2000
  19. family=reuse;policy=mat;refs=16;consumer=100%
  20. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  21. family=reuse;policy=notmat;refs=16;consumer=100%
  22. family=reuse;policy=default;refs=16;consumer=100%
  23. family=reuse;policy=default;refs=16;consumer=1%
  24. family=reuse;policy=mat;refs=1;consumer=1%
  25. family=reuse;policy=notmat;refs=16;consumer=1%
  26. family=bag;stmt=join_where;threshold=20
  27. family=reuse;policy=notmat;refs=1;consumer=100%
  28. family=reuse;policy=default;refs=1;consumer=100%
  29. family=reuse;policy=default;refs=4;consumer=1%
  30. family=bag;stmt=cte_joined_outer_notmat;threshold=180000

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
| family=bag;stmt=cte_joined_outer;threshold=180000 | j05_cte_joined_outer | valid | `52802bb56204ba17` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20 | j05_cte_joined_outer | valid | `96ff1561b4661182` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=2000 | j05_cte_joined_outer | valid | `e814390362c960c5` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | j05m_cte_joined_outer_mat | valid | `e2bc857ef4ac78f3` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | j05m_cte_joined_outer_mat | valid | `6f2db14416bce749` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | j05m_cte_joined_outer_mat | valid | `e128a8b71aab7284` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | j05n_cte_joined_outer_notmat | valid | `a9d9a341f736b9fb` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | j05n_cte_joined_outer_notmat | valid | `9da01888b806c8ae` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | j05n_cte_joined_outer_notmat | valid | `0c76f1d38a6f427b` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=join_where;threshold=180000 | j01_join_where | valid | `2556560d56106bc8` |  |
| family=bag;stmt=join_where;threshold=20 | j01_join_where | valid | `155289556d1579ce` |  |
| family=bag;stmt=join_where;threshold=2000 | j01_join_where | valid | `a404577398d577da` |  |
| family=reuse;policy=default;refs=16;consumer=1% | r_default_n16_c1 | valid | `5d6e5df01bdc3adb` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=16;consumer=100% | r_default_n16_c100 | valid | `3b3fde755d53868b` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=1;consumer=1% | r_default_n1_c1 | valid | `0e33c030da3b9c9d` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=1;consumer=100% | r_default_n1_c100 | valid | `d401822e9c9c82ca` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=4;consumer=1% | r_default_n4_c1 | valid | `5af43fc3681c9afc` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=100% | r_default_n4_c100 | valid | `40caa23b129cba52` | policy=default; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=1% | r_mat_n16_c1 | valid | `12dd2f7e44e464d3` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=100% | r_mat_n16_c100 | valid | `0692eb44eb6c9de9` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=1% | r_mat_n1_c1 | valid | `800e08adddffb0ee` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=100% | r_mat_n1_c100 | valid | `a2880e4afc98d41d` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=1% | r_mat_n4_c1 | valid | `28d2bc7359897193` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=100% | r_mat_n4_c100 | valid | `3b0ae1c80cfb6d91` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=notmat;refs=16;consumer=1% | r_notmat_n16_c1 | valid | `dc0fb63d7cc06a37` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `88497375a98ee09c` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `2557ac391eece573` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `3ef24b2c867afbf4` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `b38d20192dcf7277` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `ee96e5744e02c459` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, DIST, DEBUG, FORMAT TEXT` supported on yugabyte
- one host: this run used node budget 2 CPU 3g and client budget 1 CPU 2g, so budget comparisons are separate classes here, not independent hardware
- isolated phase: no concurrent writers in this run; the blended read/write phase is a separate labelled run
- cluster nodes are containers on the same host, so these are not network or independent-hardware measurements; container colocation is not data-colocation evidence
- readers round-robin across 3 endpoints; the writer's endpoint policy is recorded separately
- physical rank-index arm yb-range-rank; an engine whose planner ignores the index is a different physical design, not a faster one
- YugabyteDB DIST/DEBUG storage counters (rows and index rows scanned) are parsed per cell; RPC-level counters are not aggregated
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
