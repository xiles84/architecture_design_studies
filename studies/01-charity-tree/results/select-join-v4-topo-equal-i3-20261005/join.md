# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 30 cells, 30 valid, engine `yugabyte`, topology `single`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-v4-topo-equal-i3-20261005`
- run tag: `run/study-01/select-join-v4-topo-equal-i3-20261005`
- environment: `host-zenbook-ux5406sa`
- repo commit: `ba821e7b686ac3adb1cda56dd45f15c8ae699ef8`
- engine version: `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit`
- database engine image: `docker.io/yugabytedb/yugabyte:2025.2.6.0-b111` @ `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`
- benchmark client image: `localhost/charitybench-select:select-join-v4-topo-equal-i3-20261005` @ `b6ea3c007cf459112eaa8e9e85381a7dc25a7710ff07648fde3048f747df6ab0`
- inputs digest: `eb1eda2ea2f9f4df1f54c640b91c49ca45a903827d648ce2474d1535b5e99432`
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
| 20 (0.01%) | 1.8281 | missing | missing | missing | missing | missing | 1.8214 | 77.6508 | 1.7431 | missing |
| 2000 (1%) | 7.8514 | missing | missing | missing | missing | missing | 8.2564 | 78.4964 | 7.7667 | missing |
| 180000 (90%) | 100.9159 | missing | missing | missing | missing | missing | 102.2606 | 96.8530 | 101.7018 | missing |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | cte_joined_outer vs join_where | 1.8281 | 1.8214 | 0.9963 | false |
| 20 | cte_joined_outer_mat vs join_where | 1.8281 | 77.6508 | 42.4757 | false |
| 20 | cte_joined_outer_notmat vs join_where | 1.8281 | 1.7431 | 0.9535 | false |
| 2000 | cte_joined_outer vs join_where | 7.8514 | 8.2564 | 1.0516 | false |
| 2000 | cte_joined_outer_mat vs join_where | 7.8514 | 78.4964 | 9.9978 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 7.8514 | 7.7667 | 0.9892 | false |
| 180000 | cte_joined_outer vs join_where | 100.9159 | 102.2606 | 1.0133 | false |
| 180000 | cte_joined_outer_mat vs join_where | 100.9159 | 96.8530 | 0.9597 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 100.9159 | 101.7018 | 1.0078 | false |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_joined_outer_mat vs cte_joined_outer | 1.8214 | 77.6508 | 42.6319 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 1.8214 | 1.7431 | 0.9570 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 8.2564 | 78.4964 | 9.5074 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 8.2564 | 7.7667 | 0.9407 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 102.2606 | 96.8530 | 0.9471 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 102.2606 | 101.7018 | 0.9945 | false | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 6.4792 | 1 | `046c13305e59289e` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 7.8252 | 1 | `df8776e70f534bf5` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | missing | | | | |
| default | 2 | 100% | missing | | | | |
| default | 4 | 1% | 7.3689 | 4 | `df299e3420826222` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 7.6248 | 4 | `20486fb62ea425d9` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | missing | | | | |
| default | 8 | 100% | missing | | | | |
| default | 16 | 1% | 8.0490 | 16 | `3024c672e00c5d7f` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 9.4933 | 16 | `0848bb1e8cb9cc49` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 8.5634 | 1 | `4d57f70029bb92cb` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 7.0869 | 1 | `e0f311f1a15d9ee4` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | missing | | | | |
| mat | 2 | 100% | missing | | | | |
| mat | 4 | 1% | 8.4822 | 4 | `6459ad46492d024f` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 7.9357 | 4 | `fb24c48e8f537b6e` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | missing | | | | |
| mat | 8 | 100% | missing | | | | |
| mat | 16 | 1% | 8.8759 | 16 | `6dbab8b6d91547cf` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 8.7376 | 16 | `ff1436806155cf7b` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 7.1454 | 1 | `3a40ef719cb54998` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 8.3915 | 1 | `95fd2fd91dcf5492` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | missing | | | | |
| notmat | 2 | 100% | missing | | | | |
| notmat | 4 | 1% | 21.7533 | 4 | `572ef31553fa9279` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 25.7612 | 4 | `dcd1eaa1b3a72c22` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | missing | | | | |
| notmat | 8 | 100% | missing | | | | |
| notmat | 16 | 1% | 72.1884 | 16 | `b594f1e5d4cf0e14` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 75.7994 | 16 | `4eabd7751899dce2` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 1.3217 | false |
| notmat | 1 | 1% | 1.1028 | false |
| mat | 1 | 100% | 0.9056 | false |
| notmat | 1 | 100% | 1.0724 | false |
| mat | 4 | 1% | 1.1511 | false |
| notmat | 4 | 1% | 2.9520 | false |
| mat | 4 | 100% | 1.0408 | false |
| notmat | 4 | 100% | 3.3786 | false |
| mat | 16 | 1% | 1.1027 | false |
| notmat | 16 | 1% | 8.9686 | false |
| mat | 16 | 100% | 0.9204 | false |
| notmat | 16 | 100% | 7.9845 | false |

## Physical arm

- `charity_pkey`: `CREATE UNIQUE INDEX charity_pkey ON selectjoin_domain.charity USING lsm (charity_id HASH)`
- `donation_person_idx`: `CREATE INDEX donation_person_idx ON selectjoin_domain.donation USING lsm (person_id HASH)`
- `donation_pkey`: `CREATE UNIQUE INDEX donation_pkey ON selectjoin_domain.donation USING lsm (donation_id HASH)`
- `donation_rank_idx`: `CREATE INDEX donation_rank_idx ON selectjoin_domain.donation USING lsm (rank ASC)`
- `person_pkey`: `CREATE UNIQUE INDEX person_pkey ON selectjoin_domain.person USING lsm (person_id HASH)`

- index arm `yb-range-rank`; placement `noncolocated-single-tablet`; node budget `6 CPU 9g`; client budget `1 CPU 2g`

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

- client: nr_periods 269, nr_throttled 113, throttled_usec 3923859 | quota 100000 100000
, memory 2147483648

- yb-single: nr_periods 337, nr_throttled 0, throttled_usec 0 | quota 600000 100000, memory 9663676416

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
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 99.1573 | 102.2606 | 105.7127 | 5 | 9.8 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 1.6540 | 1.8214 | 2.3828 | 5 | 549.0 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 7.6568 | 8.2564 | 9.0400 | 5 | 121.1 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 91.7082 | 96.8530 | 102.0805 | 5 | 10.3 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 74.9899 | 77.6508 | 79.8357 | 5 | 12.9 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 76.3497 | 78.4964 | 80.2592 | 5 | 12.7 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 99.8584 | 101.7018 | 105.5037 | 5 | 9.8 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 1.6766 | 1.7431 | 1.9118 | 5 | 573.7 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 7.5482 | 7.7667 | 8.1581 | 5 | 128.8 |
| family=bag;stmt=join_where;threshold=180000 | valid | 97.7858 | 100.9159 | 104.2490 | 5 | 9.9 |
| family=bag;stmt=join_where;threshold=20 | valid | 1.6813 | 1.8281 | 2.1442 | 5 | 547.0 |
| family=bag;stmt=join_where;threshold=2000 | valid | 7.5831 | 7.8514 | 8.9492 | 5 | 127.4 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 7.7150 | 8.0490 | 8.5952 | 5 | 124.2 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 9.1889 | 9.4933 | 9.6808 | 5 | 105.3 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 6.1323 | 6.4792 | 7.4036 | 5 | 154.3 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 6.8698 | 7.8252 | 8.9294 | 5 | 127.8 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 7.0976 | 7.3689 | 7.7956 | 5 | 135.7 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 7.3148 | 7.6248 | 8.1105 | 5 | 131.2 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 8.4941 | 8.8759 | 9.9080 | 5 | 112.7 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 8.3673 | 8.7376 | 12.0692 | 5 | 114.4 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 8.2989 | 8.5634 | 8.9068 | 5 | 116.8 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 6.9308 | 7.0869 | 12.1442 | 5 | 141.1 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 7.8152 | 8.4822 | 8.7645 | 5 | 117.9 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 7.4464 | 7.9357 | 8.9092 | 5 | 126.0 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 67.7104 | 72.1884 | 88.7775 | 5 | 13.9 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 74.0516 | 75.7994 | 89.0932 | 5 | 13.2 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 6.8064 | 7.1454 | 7.9131 | 5 | 140.0 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 8.1112 | 8.3915 | 8.6347 | 5 | 119.2 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 20.2567 | 21.7533 | 22.3718 | 5 | 46.0 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 24.7827 | 25.7612 | 27.8725 | 5 | 38.8 |

## Matched blended workload

- phase: `isolated`; target offered rate: 10 updates/s
- offered 0, committed 0, failed 0 over 0.000 s
- offered rate 0.000/s, completed rate 0.000/s
- per key: no keys touched
- ledger: 0 keys reconciled, 0 mismatch

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42466105`, 5 interleaved rounds
  1. family=bag;stmt=join_where;threshold=180000
  2. family=reuse;policy=notmat;refs=1;consumer=100%
  3. family=reuse;policy=default;refs=16;consumer=1%
  4. family=reuse;policy=mat;refs=1;consumer=100%
  5. family=reuse;policy=default;refs=4;consumer=1%
  6. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  7. family=bag;stmt=cte_joined_outer_mat;threshold=20
  8. family=reuse;policy=default;refs=16;consumer=100%
  9. family=bag;stmt=cte_joined_outer;threshold=2000
  10. family=reuse;policy=notmat;refs=16;consumer=100%
  11. family=reuse;policy=notmat;refs=4;consumer=1%
  12. family=bag;stmt=cte_joined_outer;threshold=20
  13. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  14. family=reuse;policy=mat;refs=16;consumer=1%
  15. family=reuse;policy=default;refs=1;consumer=100%
  16. family=reuse;policy=notmat;refs=16;consumer=1%
  17. family=reuse;policy=mat;refs=4;consumer=1%
  18. family=reuse;policy=mat;refs=16;consumer=100%
  19. family=bag;stmt=join_where;threshold=20
  20. family=reuse;policy=default;refs=1;consumer=1%
  21. family=bag;stmt=cte_joined_outer;threshold=180000
  22. family=reuse;policy=mat;refs=1;consumer=1%
  23. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  24. family=reuse;policy=notmat;refs=4;consumer=100%
  25. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  26. family=reuse;policy=notmat;refs=1;consumer=1%
  27. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  28. family=reuse;policy=mat;refs=4;consumer=100%
  29. family=reuse;policy=default;refs=4;consumer=100%
  30. family=bag;stmt=join_where;threshold=2000

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
| family=bag;stmt=cte_joined_outer;threshold=180000 | j05_cte_joined_outer | valid | `620b07f633aecc32` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20 | j05_cte_joined_outer | valid | `ed59a510d8e73555` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=2000 | j05_cte_joined_outer | valid | `d2a59e27c263f402` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | j05m_cte_joined_outer_mat | valid | `ee5dba0c676c03ee` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | j05m_cte_joined_outer_mat | valid | `0119aadad64f0f6c` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | j05m_cte_joined_outer_mat | valid | `ac4030ade2331251` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | j05n_cte_joined_outer_notmat | valid | `89ea86f494cdc224` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | j05n_cte_joined_outer_notmat | valid | `0a8515a058446406` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | j05n_cte_joined_outer_notmat | valid | `122302a966ee06fc` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=join_where;threshold=180000 | j01_join_where | valid | `c99c2b2bcdb06c16` |  |
| family=bag;stmt=join_where;threshold=20 | j01_join_where | valid | `1aa9c17410416935` |  |
| family=bag;stmt=join_where;threshold=2000 | j01_join_where | valid | `0fc90199366a6059` |  |
| family=reuse;policy=default;refs=16;consumer=1% | r_default_n16_c1 | valid | `3024c672e00c5d7f` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=16;consumer=100% | r_default_n16_c100 | valid | `0848bb1e8cb9cc49` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=1;consumer=1% | r_default_n1_c1 | valid | `046c13305e59289e` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=1;consumer=100% | r_default_n1_c100 | valid | `df8776e70f534bf5` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=4;consumer=1% | r_default_n4_c1 | valid | `df299e3420826222` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=100% | r_default_n4_c100 | valid | `20486fb62ea425d9` | policy=default; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=1% | r_mat_n16_c1 | valid | `6dbab8b6d91547cf` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=100% | r_mat_n16_c100 | valid | `ff1436806155cf7b` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=1% | r_mat_n1_c1 | valid | `4d57f70029bb92cb` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=100% | r_mat_n1_c100 | valid | `e0f311f1a15d9ee4` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=1% | r_mat_n4_c1 | valid | `6459ad46492d024f` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=100% | r_mat_n4_c100 | valid | `fb24c48e8f537b6e` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=notmat;refs=16;consumer=1% | r_notmat_n16_c1 | valid | `b594f1e5d4cf0e14` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `4eabd7751899dce2` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `3a40ef719cb54998` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `95fd2fd91dcf5492` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `572ef31553fa9279` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `dcd1eaa1b3a72c22` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, DIST, DEBUG, FORMAT TEXT` supported on yugabyte
- one host: this run used node budget 6 CPU 9g and client budget 1 CPU 2g, so budget comparisons are separate classes here, not independent hardware
- isolated phase: no concurrent writers in this run; the blended read/write phase is a separate labelled run
- physical rank-index arm yb-range-rank; an engine whose planner ignores the index is a different physical design, not a faster one
- YugabyteDB DIST/DEBUG storage counters (rows and index rows scanned) are parsed per cell; RPC-level counters are not aggregated
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
