# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 30 cells, 30 valid, engine `yugabyte`, topology `cluster3`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-v4-topo-colocated-i2-20261005`
- run tag: `run/study-01/select-join-v4-topo-colocated-i2-20261005`
- environment: `host-zenbook-ux5406sa`
- repo commit: `ba821e7b686ac3adb1cda56dd45f15c8ae699ef8`
- engine version: `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit`
- database engine image: `docker.io/yugabytedb/yugabyte:2025.2.6.0-b111` @ `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`
- benchmark client image: `localhost/charitybench-select:select-join-v4-topo-colocated-i2-20261005` @ `b6ea3c007cf459112eaa8e9e85381a7dc25a7710ff07648fde3048f747df6ab0`
- inputs digest: `88e0feb9497ff35fc6cadf0e9abcbdc8d5d90c7655d3c6b6b9caf84aa32aad26`
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
- setting `placement`: colocated
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
| 20 (0.01%) | 2.3887 | missing | missing | missing | missing | missing | 1.8564 | 68.8441 | 1.6661 | missing |
| 2000 (1%) | 6.5681 | missing | missing | missing | missing | missing | 5.8743 | 105.7687 | 6.9282 | missing |
| 180000 (90%) | 126.0745 | missing | missing | missing | missing | missing | 97.4336 | 121.5728 | 92.6233 | missing |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | cte_joined_outer vs join_where | 2.3887 | 1.8564 | 0.7771 | false |
| 20 | cte_joined_outer_mat vs join_where | 2.3887 | 68.8441 | 28.8208 | false |
| 20 | cte_joined_outer_notmat vs join_where | 2.3887 | 1.6661 | 0.6975 | false |
| 2000 | cte_joined_outer vs join_where | 6.5681 | 5.8743 | 0.8944 | false |
| 2000 | cte_joined_outer_mat vs join_where | 6.5681 | 105.7687 | 16.1034 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 6.5681 | 6.9282 | 1.0548 | false |
| 180000 | cte_joined_outer vs join_where | 126.0745 | 97.4336 | 0.7728 | false |
| 180000 | cte_joined_outer_mat vs join_where | 126.0745 | 121.5728 | 0.9643 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 126.0745 | 92.6233 | 0.7347 | false |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_joined_outer_mat vs cte_joined_outer | 1.8564 | 68.8441 | 37.0853 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 1.8564 | 1.6661 | 0.8975 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 5.8743 | 105.7687 | 18.0055 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 5.8743 | 6.9282 | 1.1794 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 97.4336 | 121.5728 | 1.2477 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 97.4336 | 92.6233 | 0.9506 | false | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 1.5832 | 1 | `4ab4b656b153a01b` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 6.5227 | 1 | `9aaf44b3b957ee4b` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | missing | | | | |
| default | 2 | 100% | missing | | | | |
| default | 4 | 1% | 6.1805 | 4 | `d44bee39eb56a7e8` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 6.6765 | 4 | `e373c2afe9e8637a` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | missing | | | | |
| default | 8 | 100% | missing | | | | |
| default | 16 | 1% | 7.5374 | 16 | `b6d855f5a27587d3` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 8.3265 | 16 | `cfa39b857ab94119` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 6.3761 | 1 | `9598c15b4d6ba301` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 5.9511 | 1 | `4b44982a54d7cc77` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | missing | | | | |
| mat | 2 | 100% | missing | | | | |
| mat | 4 | 1% | 7.0613 | 4 | `95b1600adac11544` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 6.8154 | 4 | `525e14a6876a8db6` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | missing | | | | |
| mat | 8 | 100% | missing | | | | |
| mat | 16 | 1% | 7.1684 | 16 | `2ce259a585726938` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 7.8844 | 16 | `27567301536cf180` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 1.6063 | 1 | `e5e50bd74653336a` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 6.3530 | 1 | `3173137567f7149c` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | missing | | | | |
| notmat | 2 | 100% | missing | | | | |
| notmat | 4 | 1% | 3.7971 | 4 | `3ab1d7a9ec1a2a99` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 22.8080 | 4 | `ef2e757af5cc9898` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | missing | | | | |
| notmat | 8 | 100% | missing | | | | |
| notmat | 16 | 1% | 15.5538 | 16 | `3c80fa47cce8466e` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 73.4265 | 16 | `8a2eb701a66ef323` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 4.0273 | false |
| notmat | 1 | 1% | 1.0146 | false |
| mat | 1 | 100% | 0.9124 | false |
| notmat | 1 | 100% | 0.9740 | false |
| mat | 4 | 1% | 1.1425 | false |
| notmat | 4 | 1% | 0.6144 | false |
| mat | 4 | 100% | 1.0208 | false |
| notmat | 4 | 100% | 3.4162 | false |
| mat | 16 | 1% | 0.9510 | false |
| notmat | 16 | 1% | 2.0636 | false |
| mat | 16 | 100% | 0.9469 | false |
| notmat | 16 | 100% | 8.8184 | false |

## Physical arm

- `charity_pkey`: `CREATE UNIQUE INDEX charity_pkey ON selectjoin_domain.charity USING lsm (charity_id ASC)`
- `donation_person_idx`: `CREATE INDEX donation_person_idx ON selectjoin_domain.donation USING lsm (person_id ASC)`
- `donation_pkey`: `CREATE UNIQUE INDEX donation_pkey ON selectjoin_domain.donation USING lsm (donation_id ASC)`
- `donation_rank_idx`: `CREATE INDEX donation_rank_idx ON selectjoin_domain.donation USING lsm (rank ASC)`
- `person_pkey`: `CREATE UNIQUE INDEX person_pkey ON selectjoin_domain.person USING lsm (person_id ASC)`

- index arm `yb-range-rank`; placement `colocated`; node budget `2 CPU 3g`; client budget `1 CPU 2g`

| table | tablets | hash key columns | colocated | tablegroup | colocation id |
|---|---|---|---|---|---|
| selectjoin_domain.charity | 1 | 0 | true | 16388 | 865308968 |
| selectjoin_domain.person | 1 | 0 | true | 16388 | 2545373406 |
| selectjoin_domain.donation | 1 | 0 | true | 16388 | 38999154 |
| selectjoin_domain.donation_pkey | 0 | 0 | false | unknown | unknown |
| selectjoin_domain.donation_rank_idx | 1 | 0 | true | 16388 | 3055644992 |
| selectjoin_domain.donation_person_idx | 1 | 0 | true | 16388 | 2937279550 |

Physical-arm gate (default join must use the rank index at its three most selective thresholds):

- threshold 20: used
- threshold 2000: used

Tablet leaders per relation (placement evidence; the runner read these through yb-admin):

| relation | leader nodes |
|---|---|
| charity | yb-n3:9100 |
| person | yb-n3:9100 |
| donation | yb-n3:9100 |
| donation_rank_idx | yb-n3:9100 |
| donation_person_idx | yb-n3:9100 |

Effective sharding flag (from the engine): `ysql_num_shards_per_tserver=2`. A relation with fewer tablets than this setting was created without an explicit split; the per-table `num_tablets` above is the authority.

## Throttling (cgroup `cpu.stat` deltas)

- client: nr_periods 261, nr_throttled 91, throttled_usec 4489964 | quota 100000 100000
, memory 2147483648

- yb-n1: nr_periods 514, nr_throttled 7, throttled_usec 1167843 | quota 200000 100000, memory 3221225472
- yb-n2: nr_periods 503, nr_throttled 6, throttled_usec 828293 | quota 200000 100000, memory 3221225472
- yb-n3: nr_periods 502, nr_throttled 5, throttled_usec 138350 | quota 200000 100000, memory 3221225472

## Access path per cell

_Engine-reported from each cell's own plan. The two scan flags are **whole-plan** booleans: they are true when any node of the plan is an index scan / sequential scan, including small parent-table scans, so they are not the rank access path. `storage rows scanned` is the engine's own count (YugabyteDB); PostgreSQL does not print one and it is shown as n/a. A cell that scanned the whole child table at a selective threshold is a different physical regime, not a faster query._

| axis | status | whole-plan index scan | whole-plan seq scan | storage rows scanned | index rows scanned | buffers hit/read | rows removed | retries |
|---|---|---|---|---|---|---|---|---|
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | true | true | 202050 | 0 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | true | false | 110 | 20 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | true | false | 6050 | 2000 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | true | true | 202050 | 0 | 0/0 | 20000 | 0 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | true | true | 202050 | 0 | 0/0 | 199980 | 0 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | true | true | 202050 | 0 | 0/0 | 198000 | 0 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | true | true | 202050 | 0 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | true | false | 110 | 20 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | true | false | 6050 | 2000 | 0/0 | 0 | 0 |
| family=bag;stmt=join_where;threshold=180000 | valid | true | true | 202050 | 0 | 0/0 | 0 | 0 |
| family=bag;stmt=join_where;threshold=20 | valid | true | false | 110 | 20 | 0/0 | 0 | 0 |
| family=bag;stmt=join_where;threshold=2000 | valid | true | false | 6050 | 2000 | 0/0 | 0 | 0 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | true | false | 4000 | 2000 | 0/0 | 1983 | 0 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | true | false | 4000 | 2000 | 0/0 | 0 | 0 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | true | false | 2000 | 0 | 0/0 | 0 | 0 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | true | false | 4000 | 2000 | 0/0 | 0 | 0 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | true | false | 4000 | 2000 | 0/0 | 1983 | 0 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | true | false | 4000 | 2000 | 0/0 | 0 | 0 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | true | false | 4000 | 2000 | 0/0 | 1983 | 0 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | true | false | 4000 | 2000 | 0/0 | 0 | 0 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | true | false | 4000 | 2000 | 0/0 | 1983 | 0 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | true | false | 4000 | 2000 | 0/0 | 0 | 0 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | true | false | 4000 | 2000 | 0/0 | 1983 | 0 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | true | false | 4000 | 2000 | 0/0 | 0 | 0 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | true | false | 32000 | 0 | 0/0 | 0 | 0 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | true | false | 64000 | 2000 | 0/0 | 0 | 0 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | true | false | 2000 | 0 | 0/0 | 0 | 0 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | true | false | 4000 | 2000 | 0/0 | 0 | 0 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | true | false | 8000 | 0 | 0/0 | 0 | 0 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | true | false | 16000 | 2000 | 0/0 | 0 | 0 |

## Latency samples and derived throughput

_The reader is sequential, so `queries/s` is 1000/median_ms, a derived rate and not a concurrent-capacity test. min/max are the five sample extremes. Each sample includes any retry and backoff time the statement needed after a transient conflict; per-cell retry counts appear in the access-path table._

| axis | status | min ms | median ms | max ms | samples | derived queries/s |
|---|---|---|---|---|---|---|
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 93.1043 | 97.4336 | 99.0320 | 5 | 10.3 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 1.6462 | 1.8564 | 22.7233 | 5 | 538.7 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 5.6338 | 5.8743 | 5.9358 | 5 | 170.2 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 120.0702 | 121.5728 | 146.6603 | 5 | 8.2 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 64.9440 | 68.8441 | 75.3333 | 5 | 14.5 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 101.4635 | 105.7687 | 115.9003 | 5 | 9.5 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 88.6938 | 92.6233 | 97.0141 | 5 | 10.8 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 1.6224 | 1.6661 | 1.8170 | 5 | 600.2 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 6.5447 | 6.9282 | 7.2746 | 5 | 144.3 |
| family=bag;stmt=join_where;threshold=180000 | valid | 120.4640 | 126.0745 | 133.9527 | 5 | 7.9 |
| family=bag;stmt=join_where;threshold=20 | valid | 2.2171 | 2.3887 | 2.6813 | 5 | 418.6 |
| family=bag;stmt=join_where;threshold=2000 | valid | 6.4272 | 6.5681 | 6.8863 | 5 | 152.3 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 7.3771 | 7.5374 | 8.0015 | 5 | 132.7 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 8.0768 | 8.3265 | 8.6144 | 5 | 120.1 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 1.5112 | 1.5832 | 1.7433 | 5 | 631.6 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 5.9627 | 6.5227 | 6.9237 | 5 | 153.3 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 5.9627 | 6.1805 | 7.1424 | 5 | 161.8 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 6.0000 | 6.6765 | 6.7400 | 5 | 149.8 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 7.0501 | 7.1684 | 7.7885 | 5 | 139.5 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 7.3963 | 7.8844 | 8.7976 | 5 | 126.8 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 6.2689 | 6.3761 | 6.5110 | 5 | 156.8 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 5.8581 | 5.9511 | 6.2406 | 5 | 168.0 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 6.7874 | 7.0613 | 7.6341 | 5 | 141.6 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 6.4883 | 6.8154 | 7.3939 | 5 | 146.7 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 14.8469 | 15.5538 | 25.5427 | 5 | 64.3 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 71.2000 | 73.4265 | 98.3862 | 5 | 13.6 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 1.5702 | 1.6063 | 1.7098 | 5 | 622.5 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 6.1616 | 6.3530 | 6.8086 | 5 | 157.4 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 3.7668 | 3.7971 | 3.8935 | 5 | 263.4 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 20.9348 | 22.8080 | 23.0382 | 5 | 43.8 |

## Matched blended workload

- phase: `isolated`; target offered rate: 10 updates/s
- offered 0, committed 0, failed 0 over 0.000 s
- offered rate 0.000/s, completed rate 0.000/s
- per key: no keys touched
- ledger: 0 keys reconciled, 0 mismatch

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42476268`, 5 interleaved rounds
  1. family=reuse;policy=default;refs=16;consumer=1%
  2. family=reuse;policy=default;refs=4;consumer=100%
  3. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  4. family=reuse;policy=notmat;refs=1;consumer=100%
  5. family=reuse;policy=notmat;refs=4;consumer=100%
  6. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  7. family=reuse;policy=default;refs=16;consumer=100%
  8. family=reuse;policy=mat;refs=16;consumer=100%
  9. family=bag;stmt=cte_joined_outer;threshold=20
  10. family=reuse;policy=mat;refs=16;consumer=1%
  11. family=reuse;policy=default;refs=1;consumer=1%
  12. family=reuse;policy=notmat;refs=16;consumer=100%
  13. family=reuse;policy=mat;refs=4;consumer=1%
  14. family=reuse;policy=mat;refs=1;consumer=1%
  15. family=reuse;policy=notmat;refs=4;consumer=1%
  16. family=reuse;policy=notmat;refs=16;consumer=1%
  17. family=bag;stmt=join_where;threshold=20
  18. family=reuse;policy=mat;refs=1;consumer=100%
  19. family=bag;stmt=join_where;threshold=2000
  20. family=reuse;policy=default;refs=1;consumer=100%
  21. family=bag;stmt=cte_joined_outer;threshold=2000
  22. family=bag;stmt=join_where;threshold=180000
  23. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  24. family=reuse;policy=notmat;refs=1;consumer=1%
  25. family=reuse;policy=mat;refs=4;consumer=100%
  26. family=reuse;policy=default;refs=4;consumer=1%
  27. family=bag;stmt=cte_joined_outer;threshold=180000
  28. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  29. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  30. family=bag;stmt=cte_joined_outer_mat;threshold=20

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
| family=bag;stmt=cte_joined_outer;threshold=180000 | j05_cte_joined_outer | valid | `ec00edea6d9ac865` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20 | j05_cte_joined_outer | valid | `43f556d758bb3c59` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=2000 | j05_cte_joined_outer | valid | `0880481b0349931d` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | j05m_cte_joined_outer_mat | valid | `7d4a858e3be96f0e` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | j05m_cte_joined_outer_mat | valid | `38854c693f0746b0` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | j05m_cte_joined_outer_mat | valid | `613d8946a8e128cc` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | j05n_cte_joined_outer_notmat | valid | `3a479514617dc6b0` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | j05n_cte_joined_outer_notmat | valid | `877ed9a27b198a2e` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | j05n_cte_joined_outer_notmat | valid | `a7c45702afb9af9c` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=join_where;threshold=180000 | j01_join_where | valid | `8ba84bcb3cbe4173` |  |
| family=bag;stmt=join_where;threshold=20 | j01_join_where | valid | `da240d2b52609ba6` |  |
| family=bag;stmt=join_where;threshold=2000 | j01_join_where | valid | `7c1b49574cb8d768` |  |
| family=reuse;policy=default;refs=16;consumer=1% | r_default_n16_c1 | valid | `b6d855f5a27587d3` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=16;consumer=100% | r_default_n16_c100 | valid | `cfa39b857ab94119` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=1;consumer=1% | r_default_n1_c1 | valid | `4ab4b656b153a01b` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=1;consumer=100% | r_default_n1_c100 | valid | `9aaf44b3b957ee4b` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=4;consumer=1% | r_default_n4_c1 | valid | `d44bee39eb56a7e8` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=100% | r_default_n4_c100 | valid | `e373c2afe9e8637a` | policy=default; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=1% | r_mat_n16_c1 | valid | `2ce259a585726938` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=100% | r_mat_n16_c100 | valid | `27567301536cf180` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=1% | r_mat_n1_c1 | valid | `9598c15b4d6ba301` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=100% | r_mat_n1_c100 | valid | `4b44982a54d7cc77` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=1% | r_mat_n4_c1 | valid | `95b1600adac11544` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=100% | r_mat_n4_c100 | valid | `525e14a6876a8db6` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=notmat;refs=16;consumer=1% | r_notmat_n16_c1 | valid | `3c80fa47cce8466e` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `8a2eb701a66ef323` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `e5e50bd74653336a` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `3173137567f7149c` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `3ab1d7a9ec1a2a99` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `ef2e757af5cc9898` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, DIST, DEBUG, FORMAT TEXT` supported on yugabyte
- one host: this run used node budget 2 CPU 3g and client budget 1 CPU 2g, so budget comparisons are separate classes here, not independent hardware
- isolated phase: no concurrent writers in this run; the blended read/write phase is a separate labelled run
- cluster nodes are containers on the same host, so these are not network or independent-hardware measurements; container colocation is not data-colocation evidence
- readers round-robin across 3 endpoints; the writer's endpoint policy is recorded separately
- physical rank-index arm yb-range-rank; an engine whose planner ignores the index is a different physical design, not a faster one
- YugabyteDB DIST/DEBUG storage counters (rows and index rows scanned) are parsed per cell; RPC-level counters are not aggregated
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
