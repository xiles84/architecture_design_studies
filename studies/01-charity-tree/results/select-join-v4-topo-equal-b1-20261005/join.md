# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 30 cells, 30 valid, engine `yugabyte`, topology `single`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-v4-topo-equal-b1-20261005`
- run tag: `run/study-01/select-join-v4-topo-equal-b1-20261005`
- environment: `host-zenbook-ux5406sa`
- repo commit: `ba821e7b686ac3adb1cda56dd45f15c8ae699ef8`
- engine version: `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit`
- database engine image: `docker.io/yugabytedb/yugabyte:2025.2.6.0-b111` @ `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`
- benchmark client image: `localhost/charitybench-select:select-join-v4-topo-equal-b1-20261005` @ `b6ea3c007cf459112eaa8e9e85381a7dc25a7710ff07648fde3048f747df6ab0`
- inputs digest: `5b2b8c1a474427784366133e9a46b5d6b023eb5d2fdc72210a01c7cbea6181b0`
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
| 20 (0.01%) | 2.3514 | missing | missing | missing | missing | missing | 1.8522 | 90.5090 | 1.9874 | missing |
| 2000 (1%) | 10.5890 | missing | missing | missing | missing | missing | 9.8925 | 94.6921 | 9.2867 | missing |
| 180000 (90%) | 119.2955 | missing | missing | missing | missing | missing | 122.2914 | 119.3151 | 114.5344 | missing |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | cte_joined_outer vs join_where | 2.3514 | 1.8522 | 0.7877 | false |
| 20 | cte_joined_outer_mat vs join_where | 2.3514 | 90.5090 | 38.4919 | false |
| 20 | cte_joined_outer_notmat vs join_where | 2.3514 | 1.9874 | 0.8452 | false |
| 2000 | cte_joined_outer vs join_where | 10.5890 | 9.8925 | 0.9342 | false |
| 2000 | cte_joined_outer_mat vs join_where | 10.5890 | 94.6921 | 8.9425 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 10.5890 | 9.2867 | 0.8770 | false |
| 180000 | cte_joined_outer vs join_where | 119.2955 | 122.2914 | 1.0251 | false |
| 180000 | cte_joined_outer_mat vs join_where | 119.2955 | 119.3151 | 1.0002 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 119.2955 | 114.5344 | 0.9601 | false |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_joined_outer_mat vs cte_joined_outer | 1.8522 | 90.5090 | 48.8647 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 1.8522 | 1.9874 | 1.0730 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 9.8925 | 94.6921 | 9.5721 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 9.8925 | 9.2867 | 0.9388 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 122.2914 | 119.3151 | 0.9757 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 122.2914 | 114.5344 | 0.9366 | false | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 6.8576 | 1 | `acec6d3c962d9f8e` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 9.1768 | 1 | `fdf7bbae23810a40` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | missing | | | | |
| default | 2 | 100% | missing | | | | |
| default | 4 | 1% | 9.0018 | 4 | `2f159182999b2958` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 9.1655 | 4 | `3592bc453de02739` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | missing | | | | |
| default | 8 | 100% | missing | | | | |
| default | 16 | 1% | 11.5325 | 16 | `fc9811c3a72cc202` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 10.5784 | 16 | `a6af303aa40c6cd0` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 9.5822 | 1 | `142e264703902eb4` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 9.0049 | 1 | `3acd4f76ebc68ebd` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | missing | | | | |
| mat | 2 | 100% | missing | | | | |
| mat | 4 | 1% | 9.2065 | 4 | `fa80534a634a95b5` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 10.4165 | 4 | `9169f82364242c01` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | missing | | | | |
| mat | 8 | 100% | missing | | | | |
| mat | 16 | 1% | 10.4806 | 16 | `d75de8014a53d070` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 10.4087 | 16 | `b11bc53e86be050f` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 7.5263 | 1 | `acc51cfec5b421c0` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 8.9860 | 1 | `eab70d6ddfb893f2` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | missing | | | | |
| notmat | 2 | 100% | missing | | | | |
| notmat | 4 | 1% | 25.3471 | 4 | `c4167bc4d55203de` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 33.2972 | 4 | `f9c79fe3da7137e1` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | missing | | | | |
| notmat | 8 | 100% | missing | | | | |
| notmat | 16 | 1% | 79.7530 | 16 | `178354d188842a9c` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 89.1317 | 16 | `d4fd240d36563943` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 1.3973 | false |
| notmat | 1 | 1% | 1.0975 | false |
| mat | 1 | 100% | 0.9813 | false |
| notmat | 1 | 100% | 0.9792 | false |
| mat | 4 | 1% | 1.0227 | false |
| notmat | 4 | 1% | 2.8158 | false |
| mat | 4 | 100% | 1.1365 | false |
| notmat | 4 | 100% | 3.6329 | false |
| mat | 16 | 1% | 0.9088 | false |
| notmat | 16 | 1% | 6.9155 | false |
| mat | 16 | 100% | 0.9840 | false |
| notmat | 16 | 100% | 8.4258 | false |

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

- client: nr_periods 305, nr_throttled 134, throttled_usec 5710025 | quota 100000 100000
, memory 2147483648

- yb-single: nr_periods 384, nr_throttled 0, throttled_usec 0 | quota 600000 100000, memory 9663676416

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
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 105.8740 | 122.2914 | 128.3843 | 5 | 8.2 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 1.6953 | 1.8522 | 2.0553 | 5 | 539.9 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 9.0391 | 9.8925 | 11.3140 | 5 | 101.1 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 107.1621 | 119.3151 | 130.8336 | 5 | 8.4 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 87.7758 | 90.5090 | 92.5420 | 5 | 11.0 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 88.6017 | 94.6921 | 110.2673 | 5 | 10.6 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 111.6737 | 114.5344 | 130.6671 | 5 | 8.7 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 1.9009 | 1.9874 | 2.1339 | 5 | 503.2 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 8.9717 | 9.2867 | 11.6429 | 5 | 107.7 |
| family=bag;stmt=join_where;threshold=180000 | valid | 110.7766 | 119.2955 | 120.5198 | 5 | 8.4 |
| family=bag;stmt=join_where;threshold=20 | valid | 2.2388 | 2.3514 | 2.5529 | 5 | 425.3 |
| family=bag;stmt=join_where;threshold=2000 | valid | 9.9960 | 10.5890 | 10.7627 | 5 | 94.4 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 10.6210 | 11.5325 | 11.7534 | 5 | 86.7 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 9.8170 | 10.5784 | 10.7187 | 5 | 94.5 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 6.6866 | 6.8576 | 7.5663 | 5 | 145.8 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 9.0174 | 9.1768 | 9.6013 | 5 | 109.0 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 8.4761 | 9.0018 | 10.9841 | 5 | 111.1 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 8.8798 | 9.1655 | 9.6466 | 5 | 109.1 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 10.3288 | 10.4806 | 11.1359 | 5 | 95.4 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 9.8721 | 10.4087 | 12.6528 | 5 | 96.1 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 9.3130 | 9.5822 | 10.2451 | 5 | 104.4 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 8.3899 | 9.0049 | 18.5876 | 5 | 111.1 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 8.6863 | 9.2065 | 9.5809 | 5 | 108.6 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 10.1487 | 10.4165 | 10.6828 | 5 | 96.0 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 74.1872 | 79.7530 | 82.9619 | 5 | 12.5 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 85.8278 | 89.1317 | 95.6155 | 5 | 11.2 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 7.2224 | 7.5263 | 7.8556 | 5 | 132.9 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 8.6925 | 8.9860 | 9.9247 | 5 | 111.3 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 23.6459 | 25.3471 | 26.5290 | 5 | 39.5 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 28.9109 | 33.2972 | 33.6967 | 5 | 30.0 |

## Matched blended workload

- phase: `blended`; target offered rate: 10 updates/s
- offered 155, committed 155, failed 0 over 15.452 s
- offered rate 10.031/s, completed rate 10.031/s
- writer operations per endpoint: [155]
- reader operations overlapping a committed write: 71 of 210
- reader retries of transient conflicts: 0
- per key: 155 keys, min 1, max 1 per key
- ledger: 155 keys reconciled, 0 mismatch
- update latency: avg 3.562 ms, min 1.569, max 67.615 over 155

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42813094`, 5 interleaved rounds
  1. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  2. family=reuse;policy=mat;refs=1;consumer=100%
  3. family=reuse;policy=default;refs=4;consumer=1%
  4. family=bag;stmt=cte_joined_outer;threshold=2000
  5. family=reuse;policy=mat;refs=4;consumer=1%
  6. family=reuse;policy=notmat;refs=4;consumer=100%
  7. family=reuse;policy=mat;refs=16;consumer=100%
  8. family=bag;stmt=cte_joined_outer;threshold=180000
  9. family=bag;stmt=join_where;threshold=20
  10. family=reuse;policy=default;refs=1;consumer=100%
  11. family=reuse;policy=notmat;refs=16;consumer=1%
  12. family=bag;stmt=join_where;threshold=2000
  13. family=reuse;policy=notmat;refs=4;consumer=1%
  14. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  15. family=reuse;policy=default;refs=16;consumer=1%
  16. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  17. family=reuse;policy=notmat;refs=16;consumer=100%
  18. family=reuse;policy=mat;refs=4;consumer=100%
  19. family=bag;stmt=cte_joined_outer_mat;threshold=20
  20. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  21. family=reuse;policy=notmat;refs=1;consumer=1%
  22. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  23. family=reuse;policy=mat;refs=1;consumer=1%
  24. family=reuse;policy=default;refs=16;consumer=100%
  25. family=reuse;policy=notmat;refs=1;consumer=100%
  26. family=reuse;policy=default;refs=1;consumer=1%
  27. family=bag;stmt=cte_joined_outer;threshold=20
  28. family=reuse;policy=default;refs=4;consumer=100%
  29. family=bag;stmt=join_where;threshold=180000
  30. family=reuse;policy=mat;refs=16;consumer=1%

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
| family=bag;stmt=cte_joined_outer;threshold=180000 | j05_cte_joined_outer | valid | `ce0b3c2ebd1a540a` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20 | j05_cte_joined_outer | valid | `4eaa737118162b65` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=2000 | j05_cte_joined_outer | valid | `d51534e87314ba28` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | j05m_cte_joined_outer_mat | valid | `235a42e62331d8c6` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | j05m_cte_joined_outer_mat | valid | `49e44f4189accce1` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | j05m_cte_joined_outer_mat | valid | `bcdc6c08247c752d` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | j05n_cte_joined_outer_notmat | valid | `e0481d0d6cca1e19` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | j05n_cte_joined_outer_notmat | valid | `491c3d138caf8405` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | j05n_cte_joined_outer_notmat | valid | `16d43c5505c5cebe` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=join_where;threshold=180000 | j01_join_where | valid | `bf2d9dc4abd7ed6e` |  |
| family=bag;stmt=join_where;threshold=20 | j01_join_where | valid | `cfe88920045d9f42` |  |
| family=bag;stmt=join_where;threshold=2000 | j01_join_where | valid | `4d14166799c327dd` |  |
| family=reuse;policy=default;refs=16;consumer=1% | r_default_n16_c1 | valid | `fc9811c3a72cc202` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=16;consumer=100% | r_default_n16_c100 | valid | `a6af303aa40c6cd0` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=1;consumer=1% | r_default_n1_c1 | valid | `acec6d3c962d9f8e` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=1;consumer=100% | r_default_n1_c100 | valid | `fdf7bbae23810a40` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=4;consumer=1% | r_default_n4_c1 | valid | `2f159182999b2958` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=100% | r_default_n4_c100 | valid | `3592bc453de02739` | policy=default; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=1% | r_mat_n16_c1 | valid | `d75de8014a53d070` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=100% | r_mat_n16_c100 | valid | `b11bc53e86be050f` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=1% | r_mat_n1_c1 | valid | `142e264703902eb4` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=100% | r_mat_n1_c100 | valid | `3acd4f76ebc68ebd` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=1% | r_mat_n4_c1 | valid | `fa80534a634a95b5` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=100% | r_mat_n4_c100 | valid | `9169f82364242c01` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=notmat;refs=16;consumer=1% | r_notmat_n16_c1 | valid | `178354d188842a9c` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `d4fd240d36563943` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `acc51cfec5b421c0` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `eab70d6ddfb893f2` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `c4167bc4d55203de` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `f9c79fe3da7137e1` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, DIST, DEBUG, FORMAT TEXT` supported on yugabyte
- one host: this run used node budget 6 CPU 9g and client budget 1 CPU 2g, so budget comparisons are separate classes here, not independent hardware
- one client process with one fixed maintenance writer at 10 updates/s; readers are sequential; the reported sample latency includes retry and backoff time
- physical rank-index arm yb-range-rank; an engine whose planner ignores the index is a different physical design, not a faster one
- YugabyteDB DIST/DEBUG storage counters (rows and index rows scanned) are parsed per cell; RPC-level counters are not aggregated
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
