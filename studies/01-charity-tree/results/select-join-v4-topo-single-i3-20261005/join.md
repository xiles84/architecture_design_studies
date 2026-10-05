# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 30 cells, 30 valid, engine `yugabyte`, topology `single`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-v4-topo-single-i3-20261005`
- run tag: `run/study-01/select-join-v4-topo-single-i3-20261005`
- environment: `host-zenbook-ux5406sa`
- repo commit: `ba821e7b686ac3adb1cda56dd45f15c8ae699ef8`
- engine version: `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit`
- database engine image: `docker.io/yugabytedb/yugabyte:2025.2.6.0-b111` @ `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`
- benchmark client image: `localhost/charitybench-select:select-join-v4-topo-single-i3-20261005` @ `b6ea3c007cf459112eaa8e9e85381a7dc25a7710ff07648fde3048f747df6ab0`
- inputs digest: `a99d6402480ba00cece3985fbe62f122e1519382e6e02ed4cd47ed5b9218d42d`
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
- setting `placement`: noncolocated-single-tablet
- setting `prepared_plan_policy`: pgx simple protocol; no prepared statements
- setting `scale`: cte-domain
- setting `selectivity_mapping`: rank <= threshold; rank is a deterministic shuffle of 1..fact_count
- setting `shards_per_tserver`: ysql_num_shards_per_tserver=2
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
| 20 (0.01%) | 1.6959 | missing | missing | missing | missing | missing | 1.7732 | 77.6468 | 1.8929 | missing |
| 2000 (1%) | 8.3845 | missing | missing | missing | missing | missing | 8.2525 | 85.6812 | 9.5892 | missing |
| 180000 (90%) | 104.4094 | missing | missing | missing | missing | missing | 102.9684 | 105.1295 | 105.8212 | missing |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | cte_joined_outer vs join_where | 1.6959 | 1.7732 | 1.0456 | false |
| 20 | cte_joined_outer_mat vs join_where | 1.6959 | 77.6468 | 45.7852 | false |
| 20 | cte_joined_outer_notmat vs join_where | 1.6959 | 1.8929 | 1.1162 | false |
| 2000 | cte_joined_outer vs join_where | 8.3845 | 8.2525 | 0.9843 | false |
| 2000 | cte_joined_outer_mat vs join_where | 8.3845 | 85.6812 | 10.2190 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 8.3845 | 9.5892 | 1.1437 | false |
| 180000 | cte_joined_outer vs join_where | 104.4094 | 102.9684 | 0.9862 | false |
| 180000 | cte_joined_outer_mat vs join_where | 104.4094 | 105.1295 | 1.0069 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 104.4094 | 105.8212 | 1.0135 | false |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_joined_outer_mat vs cte_joined_outer | 1.7732 | 77.6468 | 43.7895 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 1.7732 | 1.8929 | 1.0675 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 8.2525 | 85.6812 | 10.3824 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 8.2525 | 9.5892 | 1.1620 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 102.9684 | 105.1295 | 1.0210 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 102.9684 | 105.8212 | 1.0277 | false | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 6.3985 | 1 | `0dab767a4e40745b` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 7.5035 | 1 | `5f362ab622513965` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | missing | | | | |
| default | 2 | 100% | missing | | | | |
| default | 4 | 1% | 8.7135 | 4 | `5dd1953a595e5213` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 7.8543 | 4 | `a3d22e9489c9d73d` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | missing | | | | |
| default | 8 | 100% | missing | | | | |
| default | 16 | 1% | 9.8091 | 16 | `1e0aa14f05020141` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 9.8147 | 16 | `14cec3d31e2fd91d` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 7.8159 | 1 | `4b44939ec5d04833` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 8.3748 | 1 | `ff637792f960a3a9` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | missing | | | | |
| mat | 2 | 100% | missing | | | | |
| mat | 4 | 1% | 7.6961 | 4 | `f5266582341d8730` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 7.5740 | 4 | `94030fe631b1672b` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | missing | | | | |
| mat | 8 | 100% | missing | | | | |
| mat | 16 | 1% | 8.4722 | 16 | `5b8a3d9fcfc2a1cd` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 9.3765 | 16 | `15257bfc8a462fff` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 6.0996 | 1 | `a29d9b18ae8c577d` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 7.6852 | 1 | `573e8a3328dd48dd` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | missing | | | | |
| notmat | 2 | 100% | missing | | | | |
| notmat | 4 | 1% | 19.5558 | 4 | `f6db36d9b4789649` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 28.4197 | 4 | `0f4e3436aa4d572d` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | missing | | | | |
| notmat | 8 | 100% | missing | | | | |
| notmat | 16 | 1% | 70.0928 | 16 | `8fc1ff2431a3c477` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 109.8145 | 16 | `c3596a35116abaa6` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 1.2215 | false |
| notmat | 1 | 1% | 0.9533 | false |
| mat | 1 | 100% | 1.1161 | false |
| notmat | 1 | 100% | 1.0242 | false |
| mat | 4 | 1% | 0.8832 | false |
| notmat | 4 | 1% | 2.2443 | false |
| mat | 4 | 100% | 0.9643 | false |
| notmat | 4 | 100% | 3.6184 | false |
| mat | 16 | 1% | 0.8637 | false |
| notmat | 16 | 1% | 7.1457 | false |
| mat | 16 | 100% | 0.9554 | false |
| notmat | 16 | 100% | 11.1888 | false |

## Physical arm

- `charity_pkey`: `CREATE UNIQUE INDEX charity_pkey ON selectjoin_domain.charity USING lsm (charity_id HASH)`
- `donation_person_idx`: `CREATE INDEX donation_person_idx ON selectjoin_domain.donation USING lsm (person_id HASH)`
- `donation_pkey`: `CREATE UNIQUE INDEX donation_pkey ON selectjoin_domain.donation USING lsm (donation_id HASH)`
- `donation_rank_idx`: `CREATE INDEX donation_rank_idx ON selectjoin_domain.donation USING lsm (rank ASC)`
- `person_pkey`: `CREATE UNIQUE INDEX person_pkey ON selectjoin_domain.person USING lsm (person_id HASH)`

- index arm `yb-range-rank`; placement `noncolocated-single-tablet`; node budget `2 CPU 3g`; client budget `1 CPU 2g`

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

Effective sharding flag (from the engine): `ysql_num_shards_per_tserver=2`. A relation with fewer tablets than this setting was created without an explicit split; the per-table `num_tablets` above is the authority.

## Throttling (cgroup `cpu.stat` deltas)

- client: nr_periods 295, nr_throttled 119, throttled_usec 3940006 | quota 100000 100000
, memory 2147483648

- yb-single: nr_periods 380, nr_throttled 79, throttled_usec 5674243 | quota 200000 100000, memory 3221225472

## Access path per cell

_Engine-reported from each cell's own plan. The two scan flags are **whole-plan** booleans: they are true when any node of the plan is an index scan / sequential scan, including small parent-table scans, so they are not the rank access path. `storage rows scanned` is the engine's own count (YugabyteDB); PostgreSQL does not print one and it is shown as n/a. A cell that scanned the whole child table at a selective threshold is a different physical regime, not a faster query._

| axis | status | whole-plan index scan | whole-plan seq scan | storage rows scanned | index rows scanned | buffers hit/read | rows removed | retries |
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
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 96.8375 | 102.9684 | 106.1327 | 5 | 9.7 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 1.5813 | 1.7732 | 2.1639 | 5 | 564.0 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 8.0597 | 8.2525 | 9.2184 | 5 | 121.2 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 95.9627 | 105.1295 | 113.7105 | 5 | 9.5 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 73.4000 | 77.6468 | 79.3768 | 5 | 12.9 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 75.7133 | 85.6812 | 117.5302 | 5 | 11.7 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 99.3706 | 105.8212 | 107.4326 | 5 | 9.4 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 1.6471 | 1.8929 | 2.1460 | 5 | 528.3 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 9.2275 | 9.5892 | 10.5340 | 5 | 104.3 |
| family=bag;stmt=join_where;threshold=180000 | valid | 101.6957 | 104.4094 | 113.1852 | 5 | 9.6 |
| family=bag;stmt=join_where;threshold=20 | valid | 1.5854 | 1.6959 | 2.0882 | 5 | 589.7 |
| family=bag;stmt=join_where;threshold=2000 | valid | 7.7263 | 8.3845 | 8.4290 | 5 | 119.3 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 8.6442 | 9.8091 | 10.8919 | 5 | 101.9 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 9.1012 | 9.8147 | 10.8652 | 5 | 101.9 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 6.1749 | 6.3985 | 7.0330 | 5 | 156.3 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 7.0019 | 7.5035 | 8.2784 | 5 | 133.3 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 8.2983 | 8.7135 | 9.0249 | 5 | 114.8 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 7.6660 | 7.8543 | 9.5987 | 5 | 127.3 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 8.1212 | 8.4722 | 9.3173 | 5 | 118.0 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 9.1343 | 9.3765 | 10.0192 | 5 | 106.6 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 6.9968 | 7.8159 | 8.4360 | 5 | 127.9 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 7.8006 | 8.3748 | 8.9677 | 5 | 119.4 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 7.5284 | 7.6961 | 26.6268 | 5 | 129.9 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 7.2462 | 7.5740 | 7.9305 | 5 | 132.0 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 65.8674 | 70.0928 | 71.6667 | 5 | 14.3 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 78.0708 | 109.8145 | 117.3601 | 5 | 9.1 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 5.4988 | 6.0996 | 7.3935 | 5 | 163.9 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 7.3475 | 7.6852 | 9.0906 | 5 | 130.1 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 19.0635 | 19.5558 | 29.2663 | 5 | 51.1 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 26.4169 | 28.4197 | 29.5713 | 5 | 35.2 |

## Matched blended workload

- phase: `isolated`; target offered rate: 10 updates/s
- offered 0, committed 0, failed 0 over 0.000 s
- offered rate 0.000/s, completed rate 0.000/s
- per key: no keys touched
- ledger: 0 keys reconciled, 0 mismatch

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42231707`, 5 interleaved rounds
  1. family=reuse;policy=notmat;refs=4;consumer=1%
  2. family=reuse;policy=notmat;refs=1;consumer=1%
  3. family=bag;stmt=join_where;threshold=180000
  4. family=reuse;policy=notmat;refs=16;consumer=1%
  5. family=reuse;policy=mat;refs=1;consumer=100%
  6. family=reuse;policy=mat;refs=1;consumer=1%
  7. family=reuse;policy=mat;refs=4;consumer=1%
  8. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  9. family=reuse;policy=notmat;refs=4;consumer=100%
  10. family=bag;stmt=cte_joined_outer;threshold=20
  11. family=bag;stmt=cte_joined_outer_mat;threshold=20
  12. family=reuse;policy=default;refs=16;consumer=100%
  13. family=bag;stmt=join_where;threshold=20
  14. family=reuse;policy=mat;refs=16;consumer=100%
  15. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  16. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  17. family=reuse;policy=default;refs=4;consumer=1%
  18. family=reuse;policy=default;refs=1;consumer=100%
  19. family=bag;stmt=cte_joined_outer;threshold=180000
  20. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  21. family=reuse;policy=notmat;refs=16;consumer=100%
  22. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  23. family=reuse;policy=default;refs=16;consumer=1%
  24. family=bag;stmt=cte_joined_outer;threshold=2000
  25. family=reuse;policy=notmat;refs=1;consumer=100%
  26. family=reuse;policy=default;refs=4;consumer=100%
  27. family=reuse;policy=mat;refs=16;consumer=1%
  28. family=reuse;policy=mat;refs=4;consumer=100%
  29. family=bag;stmt=join_where;threshold=2000
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
| family=bag;stmt=cte_joined_outer;threshold=180000 | j05_cte_joined_outer | valid | `444488eb65bfb5bb` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20 | j05_cte_joined_outer | valid | `b99c62a6d6800cc2` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=2000 | j05_cte_joined_outer | valid | `842fd3e8293a3602` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | j05m_cte_joined_outer_mat | valid | `2f6a5fe45a3c5546` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | j05m_cte_joined_outer_mat | valid | `a72ec00ac2b9e54b` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | j05m_cte_joined_outer_mat | valid | `391b8930b9e2aeea` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | j05n_cte_joined_outer_notmat | valid | `287bb9081c0c0a20` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | j05n_cte_joined_outer_notmat | valid | `9f59ca08817e436c` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | j05n_cte_joined_outer_notmat | valid | `9c77e1c1f4ca4962` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=join_where;threshold=180000 | j01_join_where | valid | `81656d8510cdac36` |  |
| family=bag;stmt=join_where;threshold=20 | j01_join_where | valid | `388e3336eca98846` |  |
| family=bag;stmt=join_where;threshold=2000 | j01_join_where | valid | `3a042b1f64e53cab` |  |
| family=reuse;policy=default;refs=16;consumer=1% | r_default_n16_c1 | valid | `1e0aa14f05020141` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=16;consumer=100% | r_default_n16_c100 | valid | `14cec3d31e2fd91d` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=1;consumer=1% | r_default_n1_c1 | valid | `0dab767a4e40745b` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=1;consumer=100% | r_default_n1_c100 | valid | `5f362ab622513965` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=4;consumer=1% | r_default_n4_c1 | valid | `5dd1953a595e5213` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=100% | r_default_n4_c100 | valid | `a3d22e9489c9d73d` | policy=default; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=1% | r_mat_n16_c1 | valid | `5b8a3d9fcfc2a1cd` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=100% | r_mat_n16_c100 | valid | `15257bfc8a462fff` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=1% | r_mat_n1_c1 | valid | `4b44939ec5d04833` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=100% | r_mat_n1_c100 | valid | `ff637792f960a3a9` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=1% | r_mat_n4_c1 | valid | `f5266582341d8730` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=100% | r_mat_n4_c100 | valid | `94030fe631b1672b` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=notmat;refs=16;consumer=1% | r_notmat_n16_c1 | valid | `8fc1ff2431a3c477` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `c3596a35116abaa6` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `a29d9b18ae8c577d` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `573e8a3328dd48dd` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `f6db36d9b4789649` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `0f4e3436aa4d572d` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, DIST, DEBUG, FORMAT TEXT` supported on yugabyte
- one host: this run used node budget 2 CPU 3g and client budget 1 CPU 2g, so budget comparisons are separate classes here, not independent hardware
- isolated phase: no concurrent writers in this run; the blended read/write phase is a separate labelled run
- physical rank-index arm yb-range-rank; an engine whose planner ignores the index is a different physical design, not a faster one
- YugabyteDB DIST/DEBUG storage counters (rows and index rows scanned) are parsed per cell; RPC-level counters are not aggregated
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
