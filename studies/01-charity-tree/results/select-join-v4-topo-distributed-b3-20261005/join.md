# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 30 cells, 30 valid, engine `yugabyte`, topology `cluster3`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-v4-topo-distributed-b3-20261005`
- run tag: `run/study-01/select-join-v4-topo-distributed-b3-20261005`
- environment: `host-zenbook-ux5406sa`
- repo commit: `ba821e7b686ac3adb1cda56dd45f15c8ae699ef8`
- engine version: `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit`
- database engine image: `docker.io/yugabytedb/yugabyte:2025.2.6.0-b111` @ `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`
- benchmark client image: `localhost/charitybench-select:select-join-v4-topo-distributed-b3-20261005` @ `b6ea3c007cf459112eaa8e9e85381a7dc25a7710ff07648fde3048f747df6ab0`
- inputs digest: `1e62a6d00eb0c539c464a0418f5fa923fbde0556ece7132df02b25d848c5a782`
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
- setting `placement`: distributed
- setting `prepared_plan_policy`: pgx simple protocol; no prepared statements
- setting `scale`: cte-domain
- setting `selectivity_mapping`: rank <= threshold; rank is a deterministic shuffle of 1..fact_count
- setting `shards_per_tserver`: ysql_num_shards_per_tserver=2
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
| yb-n1 | 169 | 2 |
| yb-n2 | 169 | 2 |
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
| 20 (0.01%) | 2.7574 | missing | missing | missing | missing | missing | 3.5478 | 69.2856 | 2.7102 | missing |
| 2000 (1%) | 8.2406 | missing | missing | missing | missing | missing | 8.9192 | 69.3664 | 8.7445 | missing |
| 180000 (90%) | 80.2491 | missing | missing | missing | missing | missing | 86.2732 | 114.3705 | 81.3291 | missing |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | cte_joined_outer vs join_where | 2.7574 | 3.5478 | 1.2867 | false |
| 20 | cte_joined_outer_mat vs join_where | 2.7574 | 69.2856 | 25.1273 | false |
| 20 | cte_joined_outer_notmat vs join_where | 2.7574 | 2.7102 | 0.9829 | false |
| 2000 | cte_joined_outer vs join_where | 8.2406 | 8.9192 | 1.0823 | false |
| 2000 | cte_joined_outer_mat vs join_where | 8.2406 | 69.3664 | 8.4176 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 8.2406 | 8.7445 | 1.0611 | false |
| 180000 | cte_joined_outer vs join_where | 80.2491 | 86.2732 | 1.0751 | false |
| 180000 | cte_joined_outer_mat vs join_where | 80.2491 | 114.3705 | 1.4252 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 80.2491 | 81.3291 | 1.0135 | false |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_joined_outer_mat vs cte_joined_outer | 3.5478 | 69.2856 | 19.5292 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 3.5478 | 2.7102 | 0.7639 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 8.9192 | 69.3664 | 7.7772 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 8.9192 | 8.7445 | 0.9804 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 86.2732 | 114.3705 | 1.3257 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 86.2732 | 81.3291 | 0.9427 | false | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 5.0454 | 1 | `b1005ebec4b05699` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 7.1589 | 1 | `9cc95d4b14140015` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | missing | | | | |
| default | 2 | 100% | missing | | | | |
| default | 4 | 1% | 7.3875 | 4 | `3ba19312469e2ea3` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 7.8845 | 4 | `bf3e1d192f8a0119` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | missing | | | | |
| default | 8 | 100% | missing | | | | |
| default | 16 | 1% | 8.7935 | 16 | `aded255a103591a8` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 8.4914 | 16 | `6fcaef612bd4d65b` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 6.9311 | 1 | `2cdede407390337f` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 7.7094 | 1 | `950b04dd42fb567d` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | missing | | | | |
| mat | 2 | 100% | missing | | | | |
| mat | 4 | 1% | 7.4220 | 4 | `b7d7d50aee2d7258` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 7.7442 | 4 | `8649e997a99db3f3` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | missing | | | | |
| mat | 8 | 100% | missing | | | | |
| mat | 16 | 1% | 8.5272 | 16 | `fd7133bc4f4745c1` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 8.8962 | 16 | `571b2d250569ab76` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 5.3236 | 1 | `7540129a8e5587a3` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 7.3270 | 1 | `33e59f8a8de2aafb` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | missing | | | | |
| notmat | 2 | 100% | missing | | | | |
| notmat | 4 | 1% | 19.7329 | 4 | `7ae9eecb093e5595` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 26.7053 | 4 | `409a81209a48cbd3` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | missing | | | | |
| notmat | 8 | 100% | missing | | | | |
| notmat | 16 | 1% | 87.9638 | 16 | `3ee9ba0c7a265493` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 88.9342 | 16 | `d9cbe9aef0eea144` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 1.3737 | false |
| notmat | 1 | 1% | 1.0551 | false |
| mat | 1 | 100% | 1.0769 | false |
| notmat | 1 | 100% | 1.0235 | false |
| mat | 4 | 1% | 1.0047 | false |
| notmat | 4 | 1% | 2.6711 | false |
| mat | 4 | 100% | 0.9822 | false |
| notmat | 4 | 100% | 3.3871 | false |
| mat | 16 | 1% | 0.9697 | false |
| notmat | 16 | 1% | 10.0033 | false |
| mat | 16 | 100% | 1.0477 | false |
| notmat | 16 | 100% | 10.4734 | false |

## Physical arm

- `charity_pkey`: `CREATE UNIQUE INDEX charity_pkey ON selectjoin_domain.charity USING lsm (charity_id HASH)`
- `donation_person_idx`: `CREATE INDEX donation_person_idx ON selectjoin_domain.donation USING lsm (person_id HASH)`
- `donation_pkey`: `CREATE UNIQUE INDEX donation_pkey ON selectjoin_domain.donation USING lsm (donation_id HASH)`
- `donation_rank_idx`: `CREATE INDEX donation_rank_idx ON selectjoin_domain.donation USING lsm (rank ASC)`
- `person_pkey`: `CREATE UNIQUE INDEX person_pkey ON selectjoin_domain.person USING lsm (person_id HASH)`

- index arm `yb-range-rank`; placement `distributed`; node budget `2 CPU 3g`; client budget `1 CPU 2g`

| table | tablets | hash key columns | colocated | tablegroup | colocation id |
|---|---|---|---|---|---|
| selectjoin_domain.charity | 3 | 1 | false | unknown | unknown |
| selectjoin_domain.person | 3 | 1 | false | unknown | unknown |
| selectjoin_domain.donation | 3 | 1 | false | unknown | unknown |
| selectjoin_domain.donation_pkey | 0 | 0 | false | unknown | unknown |
| selectjoin_domain.donation_rank_idx | 3 | 0 | false | unknown | unknown |
| selectjoin_domain.donation_person_idx | 3 | 1 | false | unknown | unknown |

Physical-arm gate (default join must use the rank index at its three most selective thresholds):

- threshold 20: used
- threshold 2000: used

Tablet leaders per relation (placement evidence; the runner read these through yb-admin):

| relation | leader nodes |
|---|---|
| charity | yb-n1:9100, yb-n2:9100, yb-n3:9100 |
| person | yb-n1:9100, yb-n2:9100, yb-n3:9100 |
| donation | yb-n1:9100, yb-n2:9100, yb-n3:9100 |
| donation_rank_idx | yb-n1:9100, yb-n2:9100, yb-n3:9100 |
| donation_person_idx | yb-n1:9100, yb-n2:9100, yb-n3:9100 |

Effective sharding flag (from the engine): `ysql_num_shards_per_tserver=2`. A relation with fewer tablets than this setting was created without an explicit split; the per-table `num_tablets` above is the authority.

## Throttling (cgroup `cpu.stat` deltas)

- client: nr_periods 313, nr_throttled 139, throttled_usec 5200375 | quota 100000 100000
, memory 2147483648

- yb-n1: nr_periods 652, nr_throttled 6, throttled_usec 789477 | quota 200000 100000, memory 3221225472
- yb-n2: nr_periods 644, nr_throttled 6, throttled_usec 898858 | quota 200000 100000, memory 3221225472
- yb-n3: nr_periods 640, nr_throttled 0, throttled_usec 0 | quota 200000 100000, memory 3221225472

## Access path per cell

_Engine-reported from each cell's own plan. The two scan flags are **whole-plan** booleans: they are true when any node of the plan is an index scan / sequential scan, including small parent-table scans, so they are not the rank access path. `storage rows scanned` is the engine's own count (YugabyteDB); PostgreSQL does not print one and it is shown as n/a. A cell that scanned the whole child table at a selective threshold is a different physical regime, not a faster query._

| axis | status | whole-plan index scan | whole-plan seq scan | storage rows scanned | index rows scanned | buffers hit/read | rows removed | retries |
|---|---|---|---|---|---|---|---|---|
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | false | true | 202050 | 0 | 0/0 | 0 | 1 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | true | true | 110 | 20 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | true | true | 6050 | 2000 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | false | true | 202050 | 0 | 0/0 | 20000 | 1 |
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
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 77.4613 | 86.2732 | 156.4101 | 5 | 11.6 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 3.1751 | 3.5478 | 3.7972 | 5 | 281.9 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 7.9085 | 8.9192 | 9.4513 | 5 | 112.1 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 110.0458 | 114.3705 | 200.7190 | 5 | 8.7 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 64.9221 | 69.2856 | 77.1282 | 5 | 14.4 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 67.3616 | 69.3664 | 79.7420 | 5 | 14.4 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 78.0846 | 81.3291 | 85.5875 | 5 | 12.3 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 2.5318 | 2.7102 | 42.0136 | 5 | 369.0 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 8.1425 | 8.7445 | 9.3023 | 5 | 114.4 |
| family=bag;stmt=join_where;threshold=180000 | valid | 74.9192 | 80.2491 | 89.9418 | 5 | 12.5 |
| family=bag;stmt=join_where;threshold=20 | valid | 2.6043 | 2.7574 | 2.8554 | 5 | 362.7 |
| family=bag;stmt=join_where;threshold=2000 | valid | 7.6095 | 8.2406 | 9.8397 | 5 | 121.3 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 8.3189 | 8.7935 | 9.0301 | 5 | 113.7 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 8.1309 | 8.4914 | 8.9416 | 5 | 117.8 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 4.7750 | 5.0454 | 5.4430 | 5 | 198.2 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 6.6261 | 7.1589 | 8.2955 | 5 | 139.7 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 6.7683 | 7.3875 | 9.2222 | 5 | 135.4 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 7.2037 | 7.8845 | 8.1930 | 5 | 126.8 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 7.7785 | 8.5272 | 8.6429 | 5 | 117.3 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 8.3423 | 8.8962 | 9.8250 | 5 | 112.4 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 6.4633 | 6.9311 | 6.9742 | 5 | 144.3 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 7.6518 | 7.7094 | 7.7733 | 5 | 129.7 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 7.2655 | 7.4220 | 8.2087 | 5 | 134.7 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 7.5840 | 7.7442 | 8.9427 | 5 | 129.1 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 80.4954 | 87.9638 | 96.8105 | 5 | 11.4 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 85.3783 | 88.9342 | 102.2716 | 5 | 11.2 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 4.9476 | 5.3236 | 6.3957 | 5 | 187.8 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 6.9023 | 7.3270 | 7.6391 | 5 | 136.5 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 19.6810 | 19.7329 | 20.2674 | 5 | 50.7 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 22.9977 | 26.7053 | 27.8521 | 5 | 37.4 |

## Matched blended workload

- phase: `blended`; target offered rate: 10 updates/s
- offered 152, committed 152, failed 0 over 15.150 s
- offered rate 10.033/s, completed rate 10.033/s
- writer operations per endpoint: [51 51 50]
- reader operations overlapping a committed write: 62 of 212
- reader retries of transient conflicts: 2
- per key: 152 keys, min 1, max 1 per key
- ledger: 152 keys reconciled, 0 mismatch
- update latency: avg 7.131 ms, min 2.648, max 182.290 over 152

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42757925`, 5 interleaved rounds
  1. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  2. family=reuse;policy=notmat;refs=4;consumer=100%
  3. family=reuse;policy=default;refs=1;consumer=1%
  4. family=reuse;policy=mat;refs=1;consumer=1%
  5. family=reuse;policy=mat;refs=16;consumer=1%
  6. family=bag;stmt=cte_joined_outer;threshold=180000
  7. family=reuse;policy=mat;refs=16;consumer=100%
  8. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  9. family=reuse;policy=mat;refs=1;consumer=100%
  10. family=reuse;policy=default;refs=4;consumer=1%
  11. family=bag;stmt=cte_joined_outer;threshold=2000
  12. family=reuse;policy=default;refs=16;consumer=1%
  13. family=reuse;policy=default;refs=16;consumer=100%
  14. family=bag;stmt=join_where;threshold=180000
  15. family=bag;stmt=cte_joined_outer;threshold=20
  16. family=bag;stmt=cte_joined_outer_mat;threshold=20
  17. family=reuse;policy=default;refs=1;consumer=100%
  18. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  19. family=reuse;policy=notmat;refs=1;consumer=1%
  20. family=reuse;policy=default;refs=4;consumer=100%
  21. family=reuse;policy=notmat;refs=4;consumer=1%
  22. family=reuse;policy=notmat;refs=16;consumer=100%
  23. family=bag;stmt=join_where;threshold=2000
  24. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  25. family=reuse;policy=mat;refs=4;consumer=1%
  26. family=reuse;policy=notmat;refs=16;consumer=1%
  27. family=reuse;policy=mat;refs=4;consumer=100%
  28. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  29. family=reuse;policy=notmat;refs=1;consumer=100%
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
| family=bag;stmt=cte_joined_outer;threshold=180000 | j05_cte_joined_outer | valid | `a4676205c8216a07` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20 | j05_cte_joined_outer | valid | `5bd441c7a211910b` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=2000 | j05_cte_joined_outer | valid | `50fc12753c8ff9a7` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | j05m_cte_joined_outer_mat | valid | `914ccca957069dec` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | j05m_cte_joined_outer_mat | valid | `ee472b66590d46bc` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | j05m_cte_joined_outer_mat | valid | `d26430e7c812535c` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | j05n_cte_joined_outer_notmat | valid | `2aa4c7d6ab1e8fac` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | j05n_cte_joined_outer_notmat | valid | `7986684b82a3edfa` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | j05n_cte_joined_outer_notmat | valid | `7a7cde36bf75d390` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=join_where;threshold=180000 | j01_join_where | valid | `0eefb39f8538965e` |  |
| family=bag;stmt=join_where;threshold=20 | j01_join_where | valid | `a507929658272232` |  |
| family=bag;stmt=join_where;threshold=2000 | j01_join_where | valid | `0c8d73775ad84aa0` |  |
| family=reuse;policy=default;refs=16;consumer=1% | r_default_n16_c1 | valid | `aded255a103591a8` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=16;consumer=100% | r_default_n16_c100 | valid | `6fcaef612bd4d65b` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=1;consumer=1% | r_default_n1_c1 | valid | `b1005ebec4b05699` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=1;consumer=100% | r_default_n1_c100 | valid | `9cc95d4b14140015` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=4;consumer=1% | r_default_n4_c1 | valid | `3ba19312469e2ea3` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=100% | r_default_n4_c100 | valid | `bf3e1d192f8a0119` | policy=default; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=1% | r_mat_n16_c1 | valid | `fd7133bc4f4745c1` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=100% | r_mat_n16_c100 | valid | `571b2d250569ab76` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=1% | r_mat_n1_c1 | valid | `2cdede407390337f` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=100% | r_mat_n1_c100 | valid | `950b04dd42fb567d` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=1% | r_mat_n4_c1 | valid | `b7d7d50aee2d7258` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=100% | r_mat_n4_c100 | valid | `8649e997a99db3f3` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=notmat;refs=16;consumer=1% | r_notmat_n16_c1 | valid | `3ee9ba0c7a265493` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `d9cbe9aef0eea144` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `7540129a8e5587a3` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `33e59f8a8de2aafb` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `7ae9eecb093e5595` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `409a81209a48cbd3` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, DIST, DEBUG, FORMAT TEXT` supported on yugabyte
- distributed placement gate: three tablets confirmed by the engine, leader spread confirmed by the runner
- one host: this run used node budget 2 CPU 3g and client budget 1 CPU 2g, so budget comparisons are separate classes here, not independent hardware
- one client process with one fixed maintenance writer at 10 updates/s; readers are sequential; the reported sample latency includes retry and backoff time
- cluster nodes are containers on the same host, so these are not network or independent-hardware measurements; container colocation is not data-colocation evidence
- readers round-robin across 3 endpoints; the writer's endpoint policy is recorded separately
- physical rank-index arm yb-range-rank; an engine whose planner ignores the index is a different physical design, not a faster one
- YugabyteDB DIST/DEBUG storage counters (rows and index rows scanned) are parsed per cell; RPC-level counters are not aggregated
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
