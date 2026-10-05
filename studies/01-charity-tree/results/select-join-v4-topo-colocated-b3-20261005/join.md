# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 30 cells, 30 valid, engine `yugabyte`, topology `cluster3`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-v4-topo-colocated-b3-20261005`
- run tag: `run/study-01/select-join-v4-topo-colocated-b3-20261005`
- environment: `host-zenbook-ux5406sa`
- repo commit: `ba821e7b686ac3adb1cda56dd45f15c8ae699ef8`
- engine version: `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit`
- database engine image: `docker.io/yugabytedb/yugabyte:2025.2.6.0-b111` @ `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`
- benchmark client image: `localhost/charitybench-select:select-join-v4-topo-colocated-b3-20261005` @ `b6ea3c007cf459112eaa8e9e85381a7dc25a7710ff07648fde3048f747df6ab0`
- inputs digest: `b4e20666adf6f00ecb4771f4d9c949f8d857571a0fd4e190953dbe1b4fafaa89`
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
| 20 (0.01%) | 2.3416 | missing | missing | missing | missing | missing | 2.6428 | 73.9177 | 1.8447 | missing |
| 2000 (1%) | 6.7446 | missing | missing | missing | missing | missing | 6.9049 | 111.1469 | 7.4350 | missing |
| 180000 (90%) | 140.2607 | missing | missing | missing | missing | missing | 102.6479 | 130.6704 | 130.7481 | missing |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | cte_joined_outer vs join_where | 2.3416 | 2.6428 | 1.1287 | false |
| 20 | cte_joined_outer_mat vs join_where | 2.3416 | 73.9177 | 31.5676 | false |
| 20 | cte_joined_outer_notmat vs join_where | 2.3416 | 1.8447 | 0.7878 | false |
| 2000 | cte_joined_outer vs join_where | 6.7446 | 6.9049 | 1.0238 | false |
| 2000 | cte_joined_outer_mat vs join_where | 6.7446 | 111.1469 | 16.4793 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 6.7446 | 7.4350 | 1.1024 | false |
| 180000 | cte_joined_outer vs join_where | 140.2607 | 102.6479 | 0.7318 | false |
| 180000 | cte_joined_outer_mat vs join_where | 140.2607 | 130.6704 | 0.9316 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 140.2607 | 130.7481 | 0.9322 | false |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_joined_outer_mat vs cte_joined_outer | 2.6428 | 73.9177 | 27.9692 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 2.6428 | 1.8447 | 0.6980 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 6.9049 | 111.1469 | 16.0967 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 6.9049 | 7.4350 | 1.0768 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 102.6479 | 130.6704 | 1.2730 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 102.6479 | 130.7481 | 1.2738 | false | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 1.5879 | 1 | `203db80a9fe3d9b2` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 6.4450 | 1 | `a80bb10f89f5bf18` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | missing | | | | |
| default | 2 | 100% | missing | | | | |
| default | 4 | 1% | 6.5818 | 4 | `94a53b6f6e197829` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 6.2500 | 4 | `c861ceb3bb4514fb` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | missing | | | | |
| default | 8 | 100% | missing | | | | |
| default | 16 | 1% | 6.9366 | 16 | `4f0d35fa05af02e3` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 8.3399 | 16 | `13dd79fa2a35b464` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 6.1899 | 1 | `ddeb83ab4917efde` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 6.5031 | 1 | `75b654dc2f76fa1b` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | missing | | | | |
| mat | 2 | 100% | missing | | | | |
| mat | 4 | 1% | 6.4149 | 4 | `0e6188864623898b` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 6.1680 | 4 | `f0116f4f88119aca` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | missing | | | | |
| mat | 8 | 100% | missing | | | | |
| mat | 16 | 1% | 7.7631 | 16 | `ecbc7b325bcd430f` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 8.1085 | 16 | `37066f64bc32d20f` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 1.6901 | 1 | `31da81d1ec923700` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 6.7069 | 1 | `ecc0a80c0559f976` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | missing | | | | |
| notmat | 2 | 100% | missing | | | | |
| notmat | 4 | 1% | 4.5916 | 4 | `453a0f2667298b2c` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 23.1546 | 4 | `51c23600c1115912` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | missing | | | | |
| notmat | 8 | 100% | missing | | | | |
| notmat | 16 | 1% | 16.4177 | 16 | `d6b9c6920a935ae1` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 84.0777 | 16 | `cb412ae9de5dd76b` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 3.8982 | false |
| notmat | 1 | 1% | 1.0644 | false |
| mat | 1 | 100% | 1.0090 | false |
| notmat | 1 | 100% | 1.0406 | false |
| mat | 4 | 1% | 0.9746 | false |
| notmat | 4 | 1% | 0.6976 | false |
| mat | 4 | 100% | 0.9869 | false |
| notmat | 4 | 100% | 3.7047 | false |
| mat | 16 | 1% | 1.1191 | false |
| notmat | 16 | 1% | 2.3668 | false |
| mat | 16 | 100% | 0.9722 | false |
| notmat | 16 | 100% | 10.0814 | false |

## Physical arm

- `charity_pkey`: `CREATE UNIQUE INDEX charity_pkey ON selectjoin_domain.charity USING lsm (charity_id ASC)`
- `donation_person_idx`: `CREATE INDEX donation_person_idx ON selectjoin_domain.donation USING lsm (person_id ASC)`
- `donation_pkey`: `CREATE UNIQUE INDEX donation_pkey ON selectjoin_domain.donation USING lsm (donation_id ASC)`
- `donation_rank_idx`: `CREATE INDEX donation_rank_idx ON selectjoin_domain.donation USING lsm (rank ASC)`
- `person_pkey`: `CREATE UNIQUE INDEX person_pkey ON selectjoin_domain.person USING lsm (person_id ASC)`

- index arm `yb-range-rank`; placement `colocated`; node budget `2 CPU 3g`; client budget `1 CPU 2g`

| table | tablets | hash key columns | colocated | tablegroup | colocation id |
|---|---|---|---|---|---|
| selectjoin_domain.charity | 1 | 0 | true | 16388 | 1684299457 |
| selectjoin_domain.person | 1 | 0 | true | 16388 | 3202345061 |
| selectjoin_domain.donation | 1 | 0 | true | 16388 | 2834703299 |
| selectjoin_domain.donation_pkey | 0 | 0 | false | unknown | unknown |
| selectjoin_domain.donation_rank_idx | 1 | 0 | true | 16388 | 4127594463 |
| selectjoin_domain.donation_person_idx | 1 | 0 | true | 16388 | 1999134316 |

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

- client: nr_periods 275, nr_throttled 94, throttled_usec 4919178 | quota 100000 100000
, memory 2147483648

- yb-n1: nr_periods 517, nr_throttled 6, throttled_usec 730253 | quota 200000 100000, memory 3221225472
- yb-n2: nr_periods 509, nr_throttled 5, throttled_usec 613130 | quota 200000 100000, memory 3221225472
- yb-n3: nr_periods 507, nr_throttled 1, throttled_usec 39795 | quota 200000 100000, memory 3221225472

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
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 94.6603 | 102.6479 | 113.0108 | 5 | 9.7 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 2.5506 | 2.6428 | 2.7845 | 5 | 378.4 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 6.6873 | 6.9049 | 7.6114 | 5 | 144.8 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 126.2181 | 130.6704 | 143.3696 | 5 | 7.7 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 68.7767 | 73.9177 | 79.0483 | 5 | 13.5 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 106.0016 | 111.1469 | 112.2696 | 5 | 9.0 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 127.7088 | 130.7481 | 134.8271 | 5 | 7.6 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 1.7032 | 1.8447 | 1.8995 | 5 | 542.1 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 7.1749 | 7.4350 | 7.5810 | 5 | 134.5 |
| family=bag;stmt=join_where;threshold=180000 | valid | 134.7653 | 140.2607 | 141.6872 | 5 | 7.1 |
| family=bag;stmt=join_where;threshold=20 | valid | 2.1929 | 2.3416 | 2.5712 | 5 | 427.1 |
| family=bag;stmt=join_where;threshold=2000 | valid | 6.5839 | 6.7446 | 7.0680 | 5 | 148.3 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 6.5444 | 6.9366 | 7.1606 | 5 | 144.2 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 8.1295 | 8.3399 | 9.1074 | 5 | 119.9 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 1.4231 | 1.5879 | 1.9517 | 5 | 629.8 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 6.2610 | 6.4450 | 6.6082 | 5 | 155.2 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 6.3654 | 6.5818 | 19.8791 | 5 | 151.9 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 6.0033 | 6.2500 | 6.3195 | 5 | 160.0 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 7.5131 | 7.7631 | 8.5189 | 5 | 128.8 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 7.9805 | 8.1085 | 8.4129 | 5 | 123.3 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 5.7706 | 6.1899 | 6.3405 | 5 | 161.6 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 6.1135 | 6.5031 | 6.5880 | 5 | 153.8 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 6.3121 | 6.4149 | 6.4776 | 5 | 155.9 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 6.0138 | 6.1680 | 6.5520 | 5 | 162.1 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 15.9143 | 16.4177 | 18.7685 | 5 | 60.9 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 75.8473 | 84.0777 | 94.7408 | 5 | 11.9 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 1.5602 | 1.6901 | 1.9187 | 5 | 591.7 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 6.5510 | 6.7069 | 7.4160 | 5 | 149.1 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 4.3083 | 4.5916 | 4.9314 | 5 | 217.8 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 22.6934 | 23.1546 | 26.4245 | 5 | 43.2 |

## Matched blended workload

- phase: `blended`; target offered rate: 10 updates/s
- offered 125, committed 125, failed 0 over 12.417 s
- offered rate 10.067/s, completed rate 10.067/s
- writer operations per endpoint: [42 42 41]
- reader operations overlapping a committed write: 56 of 210
- reader retries of transient conflicts: 0
- per key: 125 keys, min 1, max 1 per key
- ledger: 125 keys reconciled, 0 mismatch
- update latency: avg 5.920 ms, min 2.335, max 165.694 over 125

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42828784`, 5 interleaved rounds
  1. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  2. family=reuse;policy=mat;refs=1;consumer=100%
  3. family=reuse;policy=default;refs=4;consumer=100%
  4. family=bag;stmt=join_where;threshold=20
  5. family=reuse;policy=default;refs=16;consumer=100%
  6. family=bag;stmt=cte_joined_outer;threshold=180000
  7. family=bag;stmt=cte_joined_outer;threshold=20
  8. family=reuse;policy=notmat;refs=4;consumer=100%
  9. family=reuse;policy=default;refs=16;consumer=1%
  10. family=reuse;policy=default;refs=1;consumer=100%
  11. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  12. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  13. family=reuse;policy=default;refs=1;consumer=1%
  14. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  15. family=reuse;policy=mat;refs=16;consumer=100%
  16. family=reuse;policy=notmat;refs=16;consumer=100%
  17. family=reuse;policy=notmat;refs=16;consumer=1%
  18. family=reuse;policy=default;refs=4;consumer=1%
  19. family=bag;stmt=cte_joined_outer;threshold=2000
  20. family=reuse;policy=notmat;refs=1;consumer=100%
  21. family=reuse;policy=mat;refs=4;consumer=100%
  22. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  23. family=bag;stmt=join_where;threshold=180000
  24. family=reuse;policy=mat;refs=16;consumer=1%
  25. family=reuse;policy=notmat;refs=1;consumer=1%
  26. family=reuse;policy=notmat;refs=4;consumer=1%
  27. family=reuse;policy=mat;refs=4;consumer=1%
  28. family=reuse;policy=mat;refs=1;consumer=1%
  29. family=bag;stmt=join_where;threshold=2000
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
| family=bag;stmt=cte_joined_outer;threshold=180000 | j05_cte_joined_outer | valid | `bb5c7d3dfae1eaaf` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20 | j05_cte_joined_outer | valid | `024cdb176f627309` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=2000 | j05_cte_joined_outer | valid | `e0b09170aaeaa238` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | j05m_cte_joined_outer_mat | valid | `a3a7ea7d9d988f89` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | j05m_cte_joined_outer_mat | valid | `bd0e0de4f592326b` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | j05m_cte_joined_outer_mat | valid | `382948d581b4b924` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | j05n_cte_joined_outer_notmat | valid | `44f8b4f8bfc9f07c` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | j05n_cte_joined_outer_notmat | valid | `347f11272d51a9fd` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | j05n_cte_joined_outer_notmat | valid | `fb855bf1c32008c2` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=join_where;threshold=180000 | j01_join_where | valid | `c5ad0b70cddf4d83` |  |
| family=bag;stmt=join_where;threshold=20 | j01_join_where | valid | `07950892bc2a778c` |  |
| family=bag;stmt=join_where;threshold=2000 | j01_join_where | valid | `c6ece4efb6df0e20` |  |
| family=reuse;policy=default;refs=16;consumer=1% | r_default_n16_c1 | valid | `4f0d35fa05af02e3` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=16;consumer=100% | r_default_n16_c100 | valid | `13dd79fa2a35b464` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=1;consumer=1% | r_default_n1_c1 | valid | `203db80a9fe3d9b2` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=1;consumer=100% | r_default_n1_c100 | valid | `a80bb10f89f5bf18` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=4;consumer=1% | r_default_n4_c1 | valid | `94a53b6f6e197829` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=100% | r_default_n4_c100 | valid | `c861ceb3bb4514fb` | policy=default; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=1% | r_mat_n16_c1 | valid | `ecbc7b325bcd430f` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=100% | r_mat_n16_c100 | valid | `37066f64bc32d20f` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=1% | r_mat_n1_c1 | valid | `ddeb83ab4917efde` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=100% | r_mat_n1_c100 | valid | `75b654dc2f76fa1b` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=1% | r_mat_n4_c1 | valid | `0e6188864623898b` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=100% | r_mat_n4_c100 | valid | `f0116f4f88119aca` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=notmat;refs=16;consumer=1% | r_notmat_n16_c1 | valid | `d6b9c6920a935ae1` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `cb412ae9de5dd76b` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `31da81d1ec923700` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `ecc0a80c0559f976` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `453a0f2667298b2c` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `51c23600c1115912` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, DIST, DEBUG, FORMAT TEXT` supported on yugabyte
- one host: this run used node budget 2 CPU 3g and client budget 1 CPU 2g, so budget comparisons are separate classes here, not independent hardware
- one client process with one fixed maintenance writer at 10 updates/s; readers are sequential; the reported sample latency includes retry and backoff time
- cluster nodes are containers on the same host, so these are not network or independent-hardware measurements; container colocation is not data-colocation evidence
- readers round-robin across 3 endpoints; the writer's endpoint policy is recorded separately
- physical rank-index arm yb-range-rank; an engine whose planner ignores the index is a different physical design, not a faster one
- YugabyteDB DIST/DEBUG storage counters (rows and index rows scanned) are parsed per cell; RPC-level counters are not aggregated
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
