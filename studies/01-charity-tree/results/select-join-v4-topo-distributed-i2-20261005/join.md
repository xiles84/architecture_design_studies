# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 30 cells, 30 valid, engine `yugabyte`, topology `cluster3`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-v4-topo-distributed-i2-20261005`
- run tag: `run/study-01/select-join-v4-topo-distributed-i2-20261005`
- environment: `host-zenbook-ux5406sa`
- repo commit: `ba821e7b686ac3adb1cda56dd45f15c8ae699ef8`
- engine version: `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit`
- database engine image: `docker.io/yugabytedb/yugabyte:2025.2.6.0-b111` @ `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`
- benchmark client image: `localhost/charitybench-select:select-join-v4-topo-distributed-i2-20261005` @ `b6ea3c007cf459112eaa8e9e85381a7dc25a7710ff07648fde3048f747df6ab0`
- inputs digest: `c379b5634cc0c83d6c5dda73cbabb8240e851f6bdedc2c45499d4b6853647554`
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
| 20 (0.01%) | 2.8698 | missing | missing | missing | missing | missing | 2.7460 | 67.3520 | 2.7759 | missing |
| 2000 (1%) | 8.0216 | missing | missing | missing | missing | missing | 7.9676 | 67.9742 | 7.7897 | missing |
| 180000 (90%) | 79.1312 | missing | missing | missing | missing | missing | 80.5871 | 103.7744 | 81.4507 | missing |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | cte_joined_outer vs join_where | 2.8698 | 2.7460 | 0.9568 | false |
| 20 | cte_joined_outer_mat vs join_where | 2.8698 | 67.3520 | 23.4691 | false |
| 20 | cte_joined_outer_notmat vs join_where | 2.8698 | 2.7759 | 0.9673 | false |
| 2000 | cte_joined_outer vs join_where | 8.0216 | 7.9676 | 0.9933 | false |
| 2000 | cte_joined_outer_mat vs join_where | 8.0216 | 67.9742 | 8.4739 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 8.0216 | 7.7897 | 0.9711 | false |
| 180000 | cte_joined_outer vs join_where | 79.1312 | 80.5871 | 1.0184 | false |
| 180000 | cte_joined_outer_mat vs join_where | 79.1312 | 103.7744 | 1.3114 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 79.1312 | 81.4507 | 1.0293 | false |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_joined_outer_mat vs cte_joined_outer | 2.7460 | 67.3520 | 24.5277 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 2.7460 | 2.7759 | 1.0109 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 7.9676 | 67.9742 | 8.5314 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 7.9676 | 7.7897 | 0.9777 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 80.5871 | 103.7744 | 1.2877 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 80.5871 | 81.4507 | 1.0107 | false | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 5.2445 | 1 | `59626ab6c8c880bc` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 6.9410 | 1 | `fa6a8cea828eeb1a` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | missing | | | | |
| default | 2 | 100% | missing | | | | |
| default | 4 | 1% | 7.1358 | 4 | `20c2fdb9d2d4d853` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 7.3744 | 4 | `2b5c8f382e290ade` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | missing | | | | |
| default | 8 | 100% | missing | | | | |
| default | 16 | 1% | 7.9709 | 16 | `e4dab1a08c04f808` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 8.8635 | 16 | `2492b2441d9e77cf` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 7.2612 | 1 | `846348f0f7fdf0ab` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 7.2869 | 1 | `433678a7928ec317` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | missing | | | | |
| mat | 2 | 100% | missing | | | | |
| mat | 4 | 1% | 7.6119 | 4 | `f6e8be51a239476e` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 7.3757 | 4 | `0de57b0de695cf7e` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | missing | | | | |
| mat | 8 | 100% | missing | | | | |
| mat | 16 | 1% | 8.1524 | 16 | `7e5c3efa948c1a9c` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 8.9812 | 16 | `6320e7b258ce37b4` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 5.2130 | 1 | `d89bc8d4a590a84c` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 6.8842 | 1 | `c667a86455326416` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | missing | | | | |
| notmat | 2 | 100% | missing | | | | |
| notmat | 4 | 1% | 18.7986 | 4 | `f37a9be660af19b0` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 26.4841 | 4 | `718976007112c263` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | missing | | | | |
| notmat | 8 | 100% | missing | | | | |
| notmat | 16 | 1% | 91.2518 | 16 | `6361a7a9768e2e7d` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 93.4224 | 16 | `c7878f89545ef735` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 1.3845 | false |
| notmat | 1 | 1% | 0.9940 | false |
| mat | 1 | 100% | 1.0498 | false |
| notmat | 1 | 100% | 0.9918 | false |
| mat | 4 | 1% | 1.0667 | false |
| notmat | 4 | 1% | 2.6344 | false |
| mat | 4 | 100% | 1.0002 | false |
| notmat | 4 | 100% | 3.5913 | false |
| mat | 16 | 1% | 1.0228 | false |
| notmat | 16 | 1% | 11.4482 | false |
| mat | 16 | 100% | 1.0133 | false |
| notmat | 16 | 100% | 10.5401 | false |

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

- client: nr_periods 308, nr_throttled 138, throttled_usec 5229216 | quota 100000 100000
, memory 2147483648

- yb-n1: nr_periods 638, nr_throttled 15, throttled_usec 1301458 | quota 200000 100000, memory 3221225472
- yb-n2: nr_periods 633, nr_throttled 6, throttled_usec 863774 | quota 200000 100000, memory 3221225472
- yb-n3: nr_periods 623, nr_throttled 0, throttled_usec 0 | quota 200000 100000, memory 3221225472

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
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 75.3855 | 80.5871 | 92.3537 | 5 | 12.4 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 2.6427 | 2.7460 | 3.6061 | 5 | 364.2 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 7.0619 | 7.9676 | 21.8666 | 5 | 125.5 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 99.3974 | 103.7744 | 114.7080 | 5 | 9.6 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 66.2736 | 67.3520 | 72.5752 | 5 | 14.8 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 65.8827 | 67.9742 | 68.6806 | 5 | 14.7 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 76.6762 | 81.4507 | 87.0833 | 5 | 12.3 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 2.5678 | 2.7759 | 3.1757 | 5 | 360.2 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 7.4872 | 7.7897 | 8.9564 | 5 | 128.4 |
| family=bag;stmt=join_where;threshold=180000 | valid | 72.8942 | 79.1312 | 103.1363 | 5 | 12.6 |
| family=bag;stmt=join_where;threshold=20 | valid | 2.6543 | 2.8698 | 3.3374 | 5 | 348.5 |
| family=bag;stmt=join_where;threshold=2000 | valid | 7.5418 | 8.0216 | 8.9395 | 5 | 124.7 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 7.3628 | 7.9709 | 8.2162 | 5 | 125.5 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 8.3198 | 8.8635 | 9.1887 | 5 | 112.8 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 4.9758 | 5.2445 | 5.9855 | 5 | 190.7 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 6.5114 | 6.9410 | 7.5111 | 5 | 144.1 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 6.9556 | 7.1358 | 7.5919 | 5 | 140.1 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 7.1802 | 7.3744 | 7.9844 | 5 | 135.6 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 8.0130 | 8.1524 | 8.1836 | 5 | 122.7 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 8.4581 | 8.9812 | 9.5105 | 5 | 111.3 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 6.5226 | 7.2612 | 7.6407 | 5 | 137.7 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 6.8503 | 7.2869 | 7.9342 | 5 | 137.2 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 7.2716 | 7.6119 | 7.7597 | 5 | 131.4 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 7.0371 | 7.3757 | 7.7529 | 5 | 135.6 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 74.7969 | 91.2518 | 98.1353 | 5 | 11.0 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 84.1295 | 93.4224 | 115.6557 | 5 | 10.7 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 4.8760 | 5.2130 | 5.4366 | 5 | 191.8 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 6.4728 | 6.8842 | 7.1851 | 5 | 145.3 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 17.8986 | 18.7986 | 19.6464 | 5 | 53.2 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 25.7004 | 26.4841 | 28.0573 | 5 | 37.8 |

## Matched blended workload

- phase: `isolated`; target offered rate: 10 updates/s
- offered 0, committed 0, failed 0 over 0.000 s
- offered rate 0.000/s, completed rate 0.000/s
- per key: no keys touched
- ledger: 0 keys reconciled, 0 mismatch

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42286121`, 5 interleaved rounds
  1. family=reuse;policy=notmat;refs=16;consumer=100%
  2. family=bag;stmt=join_where;threshold=180000
  3. family=reuse;policy=default;refs=4;consumer=100%
  4. family=bag;stmt=cte_joined_outer_mat;threshold=20
  5. family=reuse;policy=notmat;refs=1;consumer=1%
  6. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  7. family=reuse;policy=notmat;refs=16;consumer=1%
  8. family=reuse;policy=notmat;refs=4;consumer=100%
  9. family=reuse;policy=default;refs=4;consumer=1%
  10. family=reuse;policy=mat;refs=4;consumer=1%
  11. family=reuse;policy=default;refs=1;consumer=100%
  12. family=reuse;policy=notmat;refs=1;consumer=100%
  13. family=reuse;policy=mat;refs=16;consumer=100%
  14. family=reuse;policy=mat;refs=1;consumer=100%
  15. family=bag;stmt=cte_joined_outer;threshold=20
  16. family=bag;stmt=join_where;threshold=20
  17. family=reuse;policy=mat;refs=16;consumer=1%
  18. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  19. family=reuse;policy=notmat;refs=4;consumer=1%
  20. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  21. family=bag;stmt=cte_joined_outer;threshold=180000
  22. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  23. family=reuse;policy=mat;refs=4;consumer=100%
  24. family=reuse;policy=default;refs=16;consumer=1%
  25. family=reuse;policy=default;refs=16;consumer=100%
  26. family=reuse;policy=mat;refs=1;consumer=1%
  27. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  28. family=bag;stmt=join_where;threshold=2000
  29. family=reuse;policy=default;refs=1;consumer=1%
  30. family=bag;stmt=cte_joined_outer;threshold=2000

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
| family=bag;stmt=cte_joined_outer;threshold=180000 | j05_cte_joined_outer | valid | `1ab3867325f61fb6` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20 | j05_cte_joined_outer | valid | `7bb2b5f97de69d64` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=2000 | j05_cte_joined_outer | valid | `1efc753a0c0c2cd0` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | j05m_cte_joined_outer_mat | valid | `f43df07d76e2f0bf` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | j05m_cte_joined_outer_mat | valid | `61f952366a46e708` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | j05m_cte_joined_outer_mat | valid | `52549f97adfcfe41` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | j05n_cte_joined_outer_notmat | valid | `cdec658f5d1f18a3` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | j05n_cte_joined_outer_notmat | valid | `25a2be83ba376dbe` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | j05n_cte_joined_outer_notmat | valid | `e2b14c34cd9601d5` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=join_where;threshold=180000 | j01_join_where | valid | `ae10de850472344c` |  |
| family=bag;stmt=join_where;threshold=20 | j01_join_where | valid | `5fb07c7e114e4b21` |  |
| family=bag;stmt=join_where;threshold=2000 | j01_join_where | valid | `4b0a17287778c8f5` |  |
| family=reuse;policy=default;refs=16;consumer=1% | r_default_n16_c1 | valid | `e4dab1a08c04f808` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=16;consumer=100% | r_default_n16_c100 | valid | `2492b2441d9e77cf` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=1;consumer=1% | r_default_n1_c1 | valid | `59626ab6c8c880bc` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=1;consumer=100% | r_default_n1_c100 | valid | `fa6a8cea828eeb1a` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=4;consumer=1% | r_default_n4_c1 | valid | `20c2fdb9d2d4d853` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=100% | r_default_n4_c100 | valid | `2b5c8f382e290ade` | policy=default; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=1% | r_mat_n16_c1 | valid | `7e5c3efa948c1a9c` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=100% | r_mat_n16_c100 | valid | `6320e7b258ce37b4` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=1% | r_mat_n1_c1 | valid | `846348f0f7fdf0ab` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=100% | r_mat_n1_c100 | valid | `433678a7928ec317` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=1% | r_mat_n4_c1 | valid | `f6e8be51a239476e` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=100% | r_mat_n4_c100 | valid | `0de57b0de695cf7e` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=notmat;refs=16;consumer=1% | r_notmat_n16_c1 | valid | `6361a7a9768e2e7d` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `c7878f89545ef735` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `d89bc8d4a590a84c` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `c667a86455326416` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `f37a9be660af19b0` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `718976007112c263` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, DIST, DEBUG, FORMAT TEXT` supported on yugabyte
- distributed placement gate: three tablets confirmed by the engine, leader spread confirmed by the runner
- one host: this run used node budget 2 CPU 3g and client budget 1 CPU 2g, so budget comparisons are separate classes here, not independent hardware
- isolated phase: no concurrent writers in this run; the blended read/write phase is a separate labelled run
- cluster nodes are containers on the same host, so these are not network or independent-hardware measurements; container colocation is not data-colocation evidence
- readers round-robin across 3 endpoints; the writer's endpoint policy is recorded separately
- physical rank-index arm yb-range-rank; an engine whose planner ignores the index is a different physical design, not a faster one
- YugabyteDB DIST/DEBUG storage counters (rows and index rows scanned) are parsed per cell; RPC-level counters are not aggregated
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
