# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 30 cells, 30 valid, engine `yugabyte`, topology `single`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-v4-topo-single-b3-20261005`
- run tag: `run/study-01/select-join-v4-topo-single-b3-20261005`
- environment: `host-zenbook-ux5406sa`
- repo commit: `ba821e7b686ac3adb1cda56dd45f15c8ae699ef8`
- engine version: `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit`
- database engine image: `docker.io/yugabytedb/yugabyte:2025.2.6.0-b111` @ `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`
- benchmark client image: `localhost/charitybench-select:select-join-v4-topo-single-b3-20261005` @ `b6ea3c007cf459112eaa8e9e85381a7dc25a7710ff07648fde3048f747df6ab0`
- inputs digest: `339cf3f56a1f3e809ce65c22ec1e1227e135ef71f45fc7267205cd18f1a94d81`
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
- setting `phase`: blended
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
| yb-single | 504 | 2 |

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
| 20 (0.01%) | 2.2035 | missing | missing | missing | missing | missing | 2.0982 | 83.8518 | 2.1172 | missing |
| 2000 (1%) | 8.6908 | missing | missing | missing | missing | missing | 8.9308 | 85.2482 | 9.3874 | missing |
| 180000 (90%) | 111.6797 | missing | missing | missing | missing | missing | 106.1600 | 106.2035 | 108.6185 | missing |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | cte_joined_outer vs join_where | 2.2035 | 2.0982 | 0.9522 | false |
| 20 | cte_joined_outer_mat vs join_where | 2.2035 | 83.8518 | 38.0535 | false |
| 20 | cte_joined_outer_notmat vs join_where | 2.2035 | 2.1172 | 0.9608 | false |
| 2000 | cte_joined_outer vs join_where | 8.6908 | 8.9308 | 1.0276 | false |
| 2000 | cte_joined_outer_mat vs join_where | 8.6908 | 85.2482 | 9.8090 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 8.6908 | 9.3874 | 1.0802 | false |
| 180000 | cte_joined_outer vs join_where | 111.6797 | 106.1600 | 0.9506 | false |
| 180000 | cte_joined_outer_mat vs join_where | 111.6797 | 106.2035 | 0.9510 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 111.6797 | 108.6185 | 0.9726 | false |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_joined_outer_mat vs cte_joined_outer | 2.0982 | 83.8518 | 39.9628 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 2.0982 | 2.1172 | 1.0090 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 8.9308 | 85.2482 | 9.5454 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 8.9308 | 9.3874 | 1.0511 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 106.1600 | 106.2035 | 1.0004 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 106.1600 | 108.6185 | 1.0232 | false | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 6.4892 | 1 | `f4a9c7e473b4c14b` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 8.9336 | 1 | `9c8dde87aece8fd1` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | missing | | | | |
| default | 2 | 100% | missing | | | | |
| default | 4 | 1% | 9.2765 | 4 | `d469b7adde261b97` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 9.9674 | 4 | `250113c815fb905e` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | missing | | | | |
| default | 8 | 100% | missing | | | | |
| default | 16 | 1% | 9.3468 | 16 | `5ab8351b3c9f8689` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 11.0585 | 16 | `db74dae0a28ca251` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 8.8567 | 1 | `1644d35f1d92a152` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 8.7948 | 1 | `3c020b488c5ecc4e` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | missing | | | | |
| mat | 2 | 100% | missing | | | | |
| mat | 4 | 1% | 8.4168 | 4 | `b540e9396f00df80` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 9.9249 | 4 | `aa3a27bc81e17920` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | missing | | | | |
| mat | 8 | 100% | missing | | | | |
| mat | 16 | 1% | 9.0380 | 16 | `a628bcf06ded1be9` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 10.3431 | 16 | `174c0ba6e9f99764` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 7.5682 | 1 | `0771bac8c9e5cdfb` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 8.2308 | 1 | `6e4d58e6d12f3405` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | missing | | | | |
| notmat | 2 | 100% | missing | | | | |
| notmat | 4 | 1% | 22.4200 | 4 | `7274014f95d6f5c3` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 32.6549 | 4 | `e0538e7793dad451` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | missing | | | | |
| notmat | 8 | 100% | missing | | | | |
| notmat | 16 | 1% | 84.6282 | 16 | `ca013c0335a2faf8` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 118.6404 | 16 | `9b00cac17c98657b` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 1.3648 | false |
| notmat | 1 | 1% | 1.1663 | false |
| mat | 1 | 100% | 0.9845 | false |
| notmat | 1 | 100% | 0.9213 | false |
| mat | 4 | 1% | 0.9073 | false |
| notmat | 4 | 1% | 2.4168 | false |
| mat | 4 | 100% | 0.9957 | false |
| notmat | 4 | 100% | 3.2762 | false |
| mat | 16 | 1% | 0.9670 | false |
| notmat | 16 | 1% | 9.0543 | false |
| mat | 16 | 100% | 0.9353 | false |
| notmat | 16 | 100% | 10.7284 | false |

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

Tablet leaders per relation (placement evidence; the runner read these through yb-admin):

| relation | leader nodes |
|---|---|
| charity | yb-single:9100 |
| person | yb-single:9100 |
| donation | yb-single:9100 |
| donation_rank_idx | yb-single:9100 |
| donation_person_idx | yb-single:9100 |

Effective sharding flag (from the engine): `ysql_num_shards_per_tserver=2`. A relation with fewer tablets than this setting was created without an explicit split; the per-table `num_tablets` above is the authority.

## Throttling (cgroup `cpu.stat` deltas)

- client: nr_periods 301, nr_throttled 125, throttled_usec 4978502 | quota 100000 100000
, memory 2147483648

- yb-single: nr_periods 383, nr_throttled 83, throttled_usec 5824482 | quota 200000 100000, memory 3221225472

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
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 103.3150 | 106.1600 | 130.3901 | 5 | 9.4 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 1.7999 | 2.0982 | 2.2603 | 5 | 476.6 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 8.5863 | 8.9308 | 9.5802 | 5 | 112.0 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 103.1093 | 106.2035 | 127.7605 | 5 | 9.4 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 77.2949 | 83.8518 | 97.2539 | 5 | 11.9 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 78.7612 | 85.2482 | 91.9467 | 5 | 11.7 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 103.4631 | 108.6185 | 114.0382 | 5 | 9.2 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 1.8392 | 2.1172 | 27.6064 | 5 | 472.3 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 8.7756 | 9.3874 | 10.1608 | 5 | 106.5 |
| family=bag;stmt=join_where;threshold=180000 | valid | 106.0701 | 111.6797 | 117.7851 | 5 | 9.0 |
| family=bag;stmt=join_where;threshold=20 | valid | 1.9528 | 2.2035 | 2.5667 | 5 | 453.8 |
| family=bag;stmt=join_where;threshold=2000 | valid | 8.2491 | 8.6908 | 9.5121 | 5 | 115.1 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 9.1744 | 9.3468 | 10.8707 | 5 | 107.0 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 10.3209 | 11.0585 | 42.7378 | 5 | 90.4 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 6.0306 | 6.4892 | 8.4760 | 5 | 154.1 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 8.5529 | 8.9336 | 31.9518 | 5 | 111.9 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 8.8306 | 9.2765 | 9.6654 | 5 | 107.8 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 9.0010 | 9.9674 | 10.1333 | 5 | 100.3 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 8.4985 | 9.0380 | 10.1225 | 5 | 110.6 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 9.5697 | 10.3431 | 10.9659 | 5 | 96.7 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 8.6700 | 8.8567 | 9.2732 | 5 | 112.9 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 8.0480 | 8.7948 | 8.9270 | 5 | 113.7 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 8.1333 | 8.4168 | 9.1654 | 5 | 118.8 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 9.3879 | 9.9249 | 10.1384 | 5 | 100.8 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 70.4549 | 84.6282 | 113.1493 | 5 | 11.8 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 100.7733 | 118.6404 | 126.5532 | 5 | 8.4 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 7.4506 | 7.5682 | 7.9641 | 5 | 132.1 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 7.9218 | 8.2308 | 9.3447 | 5 | 121.5 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 20.6870 | 22.4200 | 23.2741 | 5 | 44.6 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 28.6013 | 32.6549 | 33.9185 | 5 | 30.6 |

## Matched blended workload

- phase: `blended`; target offered rate: 10 updates/s
- offered 152, committed 152, failed 0 over 15.144 s
- offered rate 10.037/s, completed rate 10.037/s
- writer operations per endpoint: [152]
- reader operations overlapping a committed write: 72 of 210
- reader retries of transient conflicts: 0
- per key: 152 keys, min 1, max 1 per key
- ledger: 152 keys reconciled, 0 mismatch
- update latency: avg 4.281 ms, min 1.641, max 60.488 over 152

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42637854`, 5 interleaved rounds
  1. family=bag;stmt=cte_joined_outer;threshold=20
  2. family=reuse;policy=mat;refs=1;consumer=100%
  3. family=bag;stmt=cte_joined_outer;threshold=180000
  4. family=reuse;policy=notmat;refs=1;consumer=1%
  5. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  6. family=reuse;policy=default;refs=4;consumer=1%
  7. family=reuse;policy=mat;refs=16;consumer=1%
  8. family=bag;stmt=cte_joined_outer;threshold=2000
  9. family=reuse;policy=notmat;refs=4;consumer=1%
  10. family=bag;stmt=join_where;threshold=2000
  11. family=reuse;policy=notmat;refs=16;consumer=100%
  12. family=reuse;policy=notmat;refs=4;consumer=100%
  13. family=reuse;policy=mat;refs=16;consumer=100%
  14. family=reuse;policy=mat;refs=4;consumer=1%
  15. family=reuse;policy=default;refs=1;consumer=1%
  16. family=bag;stmt=cte_joined_outer_mat;threshold=20
  17. family=reuse;policy=notmat;refs=16;consumer=1%
  18. family=bag;stmt=join_where;threshold=20
  19. family=reuse;policy=default;refs=16;consumer=100%
  20. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  21. family=reuse;policy=default;refs=1;consumer=100%
  22. family=bag;stmt=join_where;threshold=180000
  23. family=reuse;policy=default;refs=4;consumer=100%
  24. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  25. family=reuse;policy=mat;refs=1;consumer=1%
  26. family=reuse;policy=notmat;refs=1;consumer=100%
  27. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  28. family=reuse;policy=mat;refs=4;consumer=100%
  29. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  30. family=reuse;policy=default;refs=16;consumer=1%

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
| family=bag;stmt=cte_joined_outer;threshold=180000 | j05_cte_joined_outer | valid | `a300cbdcfe0570f3` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20 | j05_cte_joined_outer | valid | `d87b25d58c289ca6` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=2000 | j05_cte_joined_outer | valid | `1d48f5de508e85cc` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | j05m_cte_joined_outer_mat | valid | `7f56394cb2988741` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | j05m_cte_joined_outer_mat | valid | `7086b9117b72d565` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | j05m_cte_joined_outer_mat | valid | `7fd48dc9136b873f` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | j05n_cte_joined_outer_notmat | valid | `51e86723a959f9d6` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | j05n_cte_joined_outer_notmat | valid | `a46592f172369bc0` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | j05n_cte_joined_outer_notmat | valid | `d20f19ab63e35ad9` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=join_where;threshold=180000 | j01_join_where | valid | `5f1d683e1919c2b4` |  |
| family=bag;stmt=join_where;threshold=20 | j01_join_where | valid | `b44a7304ae413140` |  |
| family=bag;stmt=join_where;threshold=2000 | j01_join_where | valid | `5de97c8647d13aef` |  |
| family=reuse;policy=default;refs=16;consumer=1% | r_default_n16_c1 | valid | `5ab8351b3c9f8689` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=16;consumer=100% | r_default_n16_c100 | valid | `db74dae0a28ca251` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=1;consumer=1% | r_default_n1_c1 | valid | `f4a9c7e473b4c14b` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=1;consumer=100% | r_default_n1_c100 | valid | `9c8dde87aece8fd1` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=4;consumer=1% | r_default_n4_c1 | valid | `d469b7adde261b97` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=100% | r_default_n4_c100 | valid | `250113c815fb905e` | policy=default; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=1% | r_mat_n16_c1 | valid | `a628bcf06ded1be9` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=100% | r_mat_n16_c100 | valid | `174c0ba6e9f99764` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=1% | r_mat_n1_c1 | valid | `1644d35f1d92a152` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=100% | r_mat_n1_c100 | valid | `3c020b488c5ecc4e` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=1% | r_mat_n4_c1 | valid | `b540e9396f00df80` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=100% | r_mat_n4_c100 | valid | `aa3a27bc81e17920` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=notmat;refs=16;consumer=1% | r_notmat_n16_c1 | valid | `ca013c0335a2faf8` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `9b00cac17c98657b` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `0771bac8c9e5cdfb` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `6e4d58e6d12f3405` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `7274014f95d6f5c3` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `e0538e7793dad451` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, DIST, DEBUG, FORMAT TEXT` supported on yugabyte
- one host: this run used node budget 2 CPU 3g and client budget 1 CPU 2g, so budget comparisons are separate classes here, not independent hardware
- one client process with one fixed maintenance writer at 10 updates/s; readers are sequential; the reported sample latency includes retry and backoff time
- physical rank-index arm yb-range-rank; an engine whose planner ignores the index is a different physical design, not a faster one
- YugabyteDB DIST/DEBUG storage counters (rows and index rows scanned) are parsed per cell; RPC-level counters are not aggregated
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
