# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 30 cells, 30 valid, engine `yugabyte`, topology `cluster3`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-v4-topo-colocated-b2-20261005`
- run tag: `run/study-01/select-join-v4-topo-colocated-b2-20261005`
- environment: `host-zenbook-ux5406sa`
- repo commit: `ba821e7b686ac3adb1cda56dd45f15c8ae699ef8`
- engine version: `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit`
- database engine image: `docker.io/yugabytedb/yugabyte:2025.2.6.0-b111` @ `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`
- benchmark client image: `localhost/charitybench-select:select-join-v4-topo-colocated-b2-20261005` @ `b6ea3c007cf459112eaa8e9e85381a7dc25a7710ff07648fde3048f747df6ab0`
- inputs digest: `a3dc08acf943eb0f29692a94c144eee7b13abc3364ca7ecfd00f31ec0ad332c0`
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
| 20 (0.01%) | 2.2136 | missing | missing | missing | missing | missing | 2.7091 | 70.6817 | 2.3715 | missing |
| 2000 (1%) | 6.8613 | missing | missing | missing | missing | missing | 5.9729 | 70.1970 | 6.9862 | missing |
| 180000 (90%) | 128.0755 | missing | missing | missing | missing | missing | 93.6174 | 129.2045 | 127.7908 | missing |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | cte_joined_outer vs join_where | 2.2136 | 2.7091 | 1.2239 | false |
| 20 | cte_joined_outer_mat vs join_where | 2.2136 | 70.6817 | 31.9309 | false |
| 20 | cte_joined_outer_notmat vs join_where | 2.2136 | 2.3715 | 1.0713 | false |
| 2000 | cte_joined_outer vs join_where | 6.8613 | 5.9729 | 0.8705 | false |
| 2000 | cte_joined_outer_mat vs join_where | 6.8613 | 70.1970 | 10.2308 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 6.8613 | 6.9862 | 1.0182 | false |
| 180000 | cte_joined_outer vs join_where | 128.0755 | 93.6174 | 0.7310 | false |
| 180000 | cte_joined_outer_mat vs join_where | 128.0755 | 129.2045 | 1.0088 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 128.0755 | 127.7908 | 0.9978 | false |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_joined_outer_mat vs cte_joined_outer | 2.7091 | 70.6817 | 26.0901 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 2.7091 | 2.3715 | 0.8754 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 5.9729 | 70.1970 | 11.7526 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 5.9729 | 6.9862 | 1.1697 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 93.6174 | 129.2045 | 1.3801 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 93.6174 | 127.7908 | 1.3650 | false | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 1.6814 | 1 | `e6a2f24e2254d5c7` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 7.1381 | 1 | `cbf9159225e67882` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | missing | | | | |
| default | 2 | 100% | missing | | | | |
| default | 4 | 1% | 6.5491 | 4 | `6af3f9bd748caf76` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 6.2227 | 4 | `e6da29474c0a3e3d` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | missing | | | | |
| default | 8 | 100% | missing | | | | |
| default | 16 | 1% | 7.3384 | 16 | `d8026553338f6cca` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 7.8544 | 16 | `ef30912d367d8fa1` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 6.6599 | 1 | `9f2f501f7ddb9ea3` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 5.8177 | 1 | `1aa3e73f5157f3a9` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | missing | | | | |
| mat | 2 | 100% | missing | | | | |
| mat | 4 | 1% | 5.9674 | 4 | `a94d18d94ff3420a` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 6.8122 | 4 | `5d388f2faa3c874b` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | missing | | | | |
| mat | 8 | 100% | missing | | | | |
| mat | 16 | 1% | 7.6111 | 16 | `88733ef8b3d35b6c` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 8.0260 | 16 | `b672d7794d50ea65` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 1.4840 | 1 | `dca151555a146d0f` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 5.9803 | 1 | `55e4db44e8f0cdb6` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | missing | | | | |
| notmat | 2 | 100% | missing | | | | |
| notmat | 4 | 1% | 4.2791 | 4 | `3a1bf921b1d7569d` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 22.7672 | 4 | `4c5d853a00881969` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | missing | | | | |
| notmat | 8 | 100% | missing | | | | |
| notmat | 16 | 1% | 12.6650 | 16 | `3973964474471936` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 73.9947 | 16 | `6c3eb9addaa6b941` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 3.9610 | false |
| notmat | 1 | 1% | 0.8826 | false |
| mat | 1 | 100% | 0.8150 | false |
| notmat | 1 | 100% | 0.8378 | false |
| mat | 4 | 1% | 0.9112 | false |
| notmat | 4 | 1% | 0.6534 | false |
| mat | 4 | 100% | 1.0947 | false |
| notmat | 4 | 100% | 3.6587 | false |
| mat | 16 | 1% | 1.0372 | false |
| notmat | 16 | 1% | 1.7259 | false |
| mat | 16 | 100% | 1.0218 | false |
| notmat | 16 | 100% | 9.4208 | false |

## Physical arm

- `charity_pkey`: `CREATE UNIQUE INDEX charity_pkey ON selectjoin_domain.charity USING lsm (charity_id ASC)`
- `donation_person_idx`: `CREATE INDEX donation_person_idx ON selectjoin_domain.donation USING lsm (person_id ASC)`
- `donation_pkey`: `CREATE UNIQUE INDEX donation_pkey ON selectjoin_domain.donation USING lsm (donation_id ASC)`
- `donation_rank_idx`: `CREATE INDEX donation_rank_idx ON selectjoin_domain.donation USING lsm (rank ASC)`
- `person_pkey`: `CREATE UNIQUE INDEX person_pkey ON selectjoin_domain.person USING lsm (person_id ASC)`

- index arm `yb-range-rank`; placement `colocated`; node budget `2 CPU 3g`; client budget `1 CPU 2g`

| table | tablets | hash key columns | colocated | tablegroup | colocation id |
|---|---|---|---|---|---|
| selectjoin_domain.charity | 1 | 0 | true | 16388 | 2738386372 |
| selectjoin_domain.person | 1 | 0 | true | 16388 | 2526620291 |
| selectjoin_domain.donation | 1 | 0 | true | 16388 | 1460364601 |
| selectjoin_domain.donation_pkey | 0 | 0 | false | unknown | unknown |
| selectjoin_domain.donation_rank_idx | 1 | 0 | true | 16388 | 1907069609 |
| selectjoin_domain.donation_person_idx | 1 | 0 | true | 16388 | 751647271 |

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

- client: nr_periods 258, nr_throttled 83, throttled_usec 2948089 | quota 100000 100000
, memory 2147483648

- yb-n1: nr_periods 498, nr_throttled 5, throttled_usec 577296 | quota 200000 100000, memory 3221225472
- yb-n2: nr_periods 487, nr_throttled 5, throttled_usec 350507 | quota 200000 100000, memory 3221225472
- yb-n3: nr_periods 478, nr_throttled 2, throttled_usec 20643 | quota 200000 100000, memory 3221225472

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
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 88.1110 | 93.6174 | 97.9649 | 5 | 10.7 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 2.3855 | 2.7091 | 3.5368 | 5 | 369.1 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 5.9138 | 5.9729 | 8.5081 | 5 | 167.4 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 123.1895 | 129.2045 | 144.8239 | 5 | 7.7 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 69.2815 | 70.6817 | 71.4079 | 5 | 14.1 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 68.3183 | 70.1970 | 75.3525 | 5 | 14.2 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 122.3713 | 127.7908 | 130.9766 | 5 | 7.8 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 2.2006 | 2.3715 | 2.5401 | 5 | 421.7 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 6.9409 | 6.9862 | 42.5087 | 5 | 143.1 |
| family=bag;stmt=join_where;threshold=180000 | valid | 123.0662 | 128.0755 | 137.6946 | 5 | 7.8 |
| family=bag;stmt=join_where;threshold=20 | valid | 2.0933 | 2.2136 | 2.7540 | 5 | 451.8 |
| family=bag;stmt=join_where;threshold=2000 | valid | 6.6251 | 6.8613 | 7.1202 | 5 | 145.7 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 6.7889 | 7.3384 | 7.7450 | 5 | 136.3 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 7.4685 | 7.8544 | 8.7697 | 5 | 127.3 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 1.6242 | 1.6814 | 2.1561 | 5 | 594.8 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 6.7258 | 7.1381 | 7.8324 | 5 | 140.1 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 6.1456 | 6.5491 | 7.1713 | 5 | 152.7 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 5.9303 | 6.2227 | 6.3523 | 5 | 160.7 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 7.5058 | 7.6111 | 7.9875 | 5 | 131.4 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 7.7028 | 8.0260 | 8.5803 | 5 | 124.6 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 6.0566 | 6.6599 | 18.0602 | 5 | 150.2 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 5.3160 | 5.8177 | 6.3545 | 5 | 171.9 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 5.6044 | 5.9674 | 6.0444 | 5 | 167.6 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 6.4933 | 6.8122 | 9.2312 | 5 | 146.8 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 11.9538 | 12.6650 | 13.4016 | 5 | 79.0 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 68.5197 | 73.9947 | 85.1414 | 5 | 13.5 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 1.4706 | 1.4840 | 1.6215 | 5 | 673.8 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 5.7559 | 5.9803 | 6.3442 | 5 | 167.2 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 4.1563 | 4.2791 | 4.4716 | 5 | 233.7 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 22.0728 | 22.7672 | 23.1503 | 5 | 43.9 |

## Matched blended workload

- phase: `blended`; target offered rate: 10 updates/s
- offered 113, committed 113, failed 0 over 11.205 s
- offered rate 10.085/s, completed rate 10.085/s
- writer operations per endpoint: [38 38 37]
- reader operations overlapping a committed write: 58 of 210
- reader retries of transient conflicts: 0
- per key: 113 keys, min 1, max 1 per key
- ledger: 113 keys reconciled, 0 mismatch
- update latency: avg 6.277 ms, min 2.632, max 152.270 over 113

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42131005`, 5 interleaved rounds
  1. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  2. family=reuse;policy=default;refs=1;consumer=100%
  3. family=bag;stmt=cte_joined_outer_mat;threshold=20
  4. family=reuse;policy=notmat;refs=4;consumer=100%
  5. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  6. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  7. family=reuse;policy=default;refs=1;consumer=1%
  8. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  9. family=reuse;policy=notmat;refs=1;consumer=100%
  10. family=bag;stmt=join_where;threshold=20
  11. family=bag;stmt=join_where;threshold=180000
  12. family=reuse;policy=mat;refs=4;consumer=100%
  13. family=bag;stmt=join_where;threshold=2000
  14. family=reuse;policy=mat;refs=1;consumer=1%
  15. family=reuse;policy=mat;refs=4;consumer=1%
  16. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  17. family=reuse;policy=default;refs=16;consumer=100%
  18. family=bag;stmt=cte_joined_outer;threshold=2000
  19. family=reuse;policy=default;refs=4;consumer=1%
  20. family=reuse;policy=notmat;refs=1;consumer=1%
  21. family=reuse;policy=default;refs=4;consumer=100%
  22. family=reuse;policy=notmat;refs=4;consumer=1%
  23. family=reuse;policy=mat;refs=16;consumer=100%
  24. family=bag;stmt=cte_joined_outer;threshold=180000
  25. family=reuse;policy=mat;refs=16;consumer=1%
  26. family=reuse;policy=default;refs=16;consumer=1%
  27. family=reuse;policy=mat;refs=1;consumer=100%
  28. family=reuse;policy=notmat;refs=16;consumer=100%
  29. family=bag;stmt=cte_joined_outer;threshold=20
  30. family=reuse;policy=notmat;refs=16;consumer=1%

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
| family=bag;stmt=cte_joined_outer;threshold=180000 | j05_cte_joined_outer | valid | `722b8be077b27cbc` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20 | j05_cte_joined_outer | valid | `d55eb40ff780235c` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=2000 | j05_cte_joined_outer | valid | `91499a6948878b0d` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | j05m_cte_joined_outer_mat | valid | `1269ef8759229b32` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | j05m_cte_joined_outer_mat | valid | `d15a515cab494dd2` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | j05m_cte_joined_outer_mat | valid | `a377f2e03913edbf` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | j05n_cte_joined_outer_notmat | valid | `73f329aaa4a33c14` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | j05n_cte_joined_outer_notmat | valid | `091109de47f55362` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | j05n_cte_joined_outer_notmat | valid | `c2d7f3dfd2182c6b` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=join_where;threshold=180000 | j01_join_where | valid | `13fc1ecb293ee7c0` |  |
| family=bag;stmt=join_where;threshold=20 | j01_join_where | valid | `89ab3b1b5215a75f` |  |
| family=bag;stmt=join_where;threshold=2000 | j01_join_where | valid | `d8c9ec48f381ff3e` |  |
| family=reuse;policy=default;refs=16;consumer=1% | r_default_n16_c1 | valid | `d8026553338f6cca` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=16;consumer=100% | r_default_n16_c100 | valid | `ef30912d367d8fa1` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=1;consumer=1% | r_default_n1_c1 | valid | `e6a2f24e2254d5c7` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=1;consumer=100% | r_default_n1_c100 | valid | `cbf9159225e67882` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=4;consumer=1% | r_default_n4_c1 | valid | `6af3f9bd748caf76` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=100% | r_default_n4_c100 | valid | `e6da29474c0a3e3d` | policy=default; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=1% | r_mat_n16_c1 | valid | `88733ef8b3d35b6c` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=100% | r_mat_n16_c100 | valid | `b672d7794d50ea65` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=1% | r_mat_n1_c1 | valid | `9f2f501f7ddb9ea3` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=100% | r_mat_n1_c100 | valid | `1aa3e73f5157f3a9` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=1% | r_mat_n4_c1 | valid | `a94d18d94ff3420a` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=100% | r_mat_n4_c100 | valid | `5d388f2faa3c874b` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=notmat;refs=16;consumer=1% | r_notmat_n16_c1 | valid | `3973964474471936` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `6c3eb9addaa6b941` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `dca151555a146d0f` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `55e4db44e8f0cdb6` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `3a1bf921b1d7569d` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `4c5d853a00881969` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, DIST, DEBUG, FORMAT TEXT` supported on yugabyte
- one host: this run used node budget 2 CPU 3g and client budget 1 CPU 2g, so budget comparisons are separate classes here, not independent hardware
- one client process with one fixed maintenance writer at 10 updates/s; readers are sequential; the reported sample latency includes retry and backoff time
- cluster nodes are containers on the same host, so these are not network or independent-hardware measurements; container colocation is not data-colocation evidence
- readers round-robin across 3 endpoints; the writer's endpoint policy is recorded separately
- physical rank-index arm yb-range-rank; an engine whose planner ignores the index is a different physical design, not a faster one
- YugabyteDB DIST/DEBUG storage counters (rows and index rows scanned) are parsed per cell; RPC-level counters are not aggregated
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
