# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 30 cells, 30 valid, engine `yugabyte`, topology `single`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-v4-topo-single-b2-20261005`
- run tag: `run/study-01/select-join-v4-topo-single-b2-20261005`
- environment: `host-zenbook-ux5406sa`
- repo commit: `ba821e7b686ac3adb1cda56dd45f15c8ae699ef8`
- engine version: `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit`
- database engine image: `docker.io/yugabytedb/yugabyte:2025.2.6.0-b111` @ `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`
- benchmark client image: `localhost/charitybench-select:select-join-v4-topo-single-b2-20261005` @ `b6ea3c007cf459112eaa8e9e85381a7dc25a7710ff07648fde3048f747df6ab0`
- inputs digest: `82d73027d7f16910d7cbf4e49cc60d32f8da47ef39a3a5c1104607ef6a11db63`
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
| 20 (0.01%) | 1.8952 | missing | missing | missing | missing | missing | 1.9335 | 93.7461 | 2.4198 | missing |
| 2000 (1%) | 9.8186 | missing | missing | missing | missing | missing | 9.4653 | 87.5076 | 10.8565 | missing |
| 180000 (90%) | 110.4733 | missing | missing | missing | missing | missing | 112.0943 | 125.9413 | 109.8973 | missing |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | cte_joined_outer vs join_where | 1.8952 | 1.9335 | 1.0202 | false |
| 20 | cte_joined_outer_mat vs join_where | 1.8952 | 93.7461 | 49.4653 | false |
| 20 | cte_joined_outer_notmat vs join_where | 1.8952 | 2.4198 | 1.2768 | false |
| 2000 | cte_joined_outer vs join_where | 9.8186 | 9.4653 | 0.9640 | false |
| 2000 | cte_joined_outer_mat vs join_where | 9.8186 | 87.5076 | 8.9124 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 9.8186 | 10.8565 | 1.1057 | false |
| 180000 | cte_joined_outer vs join_where | 110.4733 | 112.0943 | 1.0147 | false |
| 180000 | cte_joined_outer_mat vs join_where | 110.4733 | 125.9413 | 1.1400 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 110.4733 | 109.8973 | 0.9948 | false |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_joined_outer_mat vs cte_joined_outer | 1.9335 | 93.7461 | 48.4847 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 1.9335 | 2.4198 | 1.2515 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 9.4653 | 87.5076 | 9.2451 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 9.4653 | 10.8565 | 1.1470 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 112.0943 | 125.9413 | 1.1235 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 112.0943 | 109.8973 | 0.9804 | false | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 7.2793 | 1 | `9dce1f5dbe17be5c` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 8.4978 | 1 | `509e0017f52e40cb` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | missing | | | | |
| default | 2 | 100% | missing | | | | |
| default | 4 | 1% | 8.5868 | 4 | `66504facceaba213` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 8.6191 | 4 | `d581fa7f5e08579b` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | missing | | | | |
| default | 8 | 100% | missing | | | | |
| default | 16 | 1% | 9.5944 | 16 | `3992e61f6b808571` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 11.1998 | 16 | `67f45ccb2e3ba270` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 8.3394 | 1 | `cacba2210809d6a1` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 9.2849 | 1 | `86225bb09a69a5a0` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | missing | | | | |
| mat | 2 | 100% | missing | | | | |
| mat | 4 | 1% | 9.8650 | 4 | `2efca80386d5adb0` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 9.0289 | 4 | `dae3aebfbb510def` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | missing | | | | |
| mat | 8 | 100% | missing | | | | |
| mat | 16 | 1% | 10.6204 | 16 | `f2b237437475bc92` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 10.0090 | 16 | `153fd01083180b22` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 7.5097 | 1 | `952d97bc8dffa024` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 8.5587 | 1 | `d872aa899963cbd6` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | missing | | | | |
| notmat | 2 | 100% | missing | | | | |
| notmat | 4 | 1% | 24.4475 | 4 | `795c6aa7522efa68` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 31.0923 | 4 | `6ff281c7be5c01d7` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | missing | | | | |
| notmat | 8 | 100% | missing | | | | |
| notmat | 16 | 1% | 109.9561 | 16 | `53f819dc7dd8767d` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 114.4688 | 16 | `6e72a37535e7cad4` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 1.1456 | false |
| notmat | 1 | 1% | 1.0316 | false |
| mat | 1 | 100% | 1.0926 | false |
| notmat | 1 | 100% | 1.0072 | false |
| mat | 4 | 1% | 1.1489 | false |
| notmat | 4 | 1% | 2.8471 | false |
| mat | 4 | 100% | 1.0475 | false |
| notmat | 4 | 100% | 3.6074 | false |
| mat | 16 | 1% | 1.1069 | false |
| notmat | 16 | 1% | 11.4605 | false |
| mat | 16 | 100% | 0.8937 | false |
| notmat | 16 | 100% | 10.2206 | false |

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

- client: nr_periods 314, nr_throttled 125, throttled_usec 5884815 | quota 100000 100000
, memory 2147483648

- yb-single: nr_periods 394, nr_throttled 92, throttled_usec 5185062 | quota 200000 100000, memory 3221225472

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
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 110.0574 | 112.0943 | 124.6551 | 5 | 8.9 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 1.8120 | 1.9335 | 2.1877 | 5 | 517.2 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 8.8349 | 9.4653 | 9.9153 | 5 | 105.6 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 102.0373 | 125.9413 | 141.1614 | 5 | 7.9 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 82.0237 | 93.7461 | 99.4540 | 5 | 10.7 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 81.6897 | 87.5076 | 91.0830 | 5 | 11.4 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 108.1593 | 109.8973 | 117.7630 | 5 | 9.1 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 2.3172 | 2.4198 | 2.5606 | 5 | 413.3 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 10.1286 | 10.8565 | 11.5676 | 5 | 92.1 |
| family=bag;stmt=join_where;threshold=180000 | valid | 104.9969 | 110.4733 | 113.1566 | 5 | 9.1 |
| family=bag;stmt=join_where;threshold=20 | valid | 1.8781 | 1.8952 | 18.5558 | 5 | 527.7 |
| family=bag;stmt=join_where;threshold=2000 | valid | 9.4070 | 9.8186 | 10.3406 | 5 | 101.8 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 8.9687 | 9.5944 | 11.0798 | 5 | 104.2 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 10.7464 | 11.1998 | 11.4473 | 5 | 89.3 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 6.6744 | 7.2793 | 7.5045 | 5 | 137.4 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 8.4322 | 8.4978 | 9.4473 | 5 | 117.7 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 8.3036 | 8.5868 | 8.9613 | 5 | 116.5 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 8.2770 | 8.6191 | 9.5853 | 5 | 116.0 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 10.2762 | 10.6204 | 39.7539 | 5 | 94.2 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 9.5936 | 10.0090 | 10.6434 | 5 | 99.9 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 8.1068 | 8.3394 | 8.5917 | 5 | 119.9 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 9.1434 | 9.2849 | 11.8638 | 5 | 107.7 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 9.4157 | 9.8650 | 35.4892 | 5 | 101.4 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 8.2152 | 9.0289 | 10.0519 | 5 | 110.8 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 88.4561 | 109.9561 | 113.9109 | 5 | 9.1 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 88.3893 | 114.4688 | 123.3557 | 5 | 8.7 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 7.1627 | 7.5097 | 7.6315 | 5 | 133.2 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 8.1988 | 8.5587 | 9.4709 | 5 | 116.8 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 22.5194 | 24.4475 | 32.9058 | 5 | 40.9 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 29.6696 | 31.0923 | 34.0533 | 5 | 32.2 |

## Matched blended workload

- phase: `blended`; target offered rate: 10 updates/s
- offered 155, committed 155, failed 0 over 15.442 s
- offered rate 10.037/s, completed rate 10.037/s
- writer operations per endpoint: [155]
- reader operations overlapping a committed write: 72 of 210
- reader retries of transient conflicts: 0
- per key: 155 keys, min 1, max 1 per key
- ledger: 155 keys reconciled, 0 mismatch
- update latency: avg 5.221 ms, min 1.643, max 67.144 over 155

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42961819`, 5 interleaved rounds
  1. family=reuse;policy=mat;refs=4;consumer=100%
  2. family=reuse;policy=notmat;refs=1;consumer=100%
  3. family=bag;stmt=join_where;threshold=20
  4. family=reuse;policy=default;refs=1;consumer=100%
  5. family=bag;stmt=cte_joined_outer;threshold=20
  6. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  7. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  8. family=bag;stmt=join_where;threshold=2000
  9. family=bag;stmt=cte_joined_outer;threshold=180000
  10. family=reuse;policy=default;refs=16;consumer=100%
  11. family=reuse;policy=default;refs=16;consumer=1%
  12. family=reuse;policy=default;refs=4;consumer=100%
  13. family=bag;stmt=join_where;threshold=180000
  14. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  15. family=reuse;policy=default;refs=1;consumer=1%
  16. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  17. family=reuse;policy=notmat;refs=1;consumer=1%
  18. family=reuse;policy=default;refs=4;consumer=1%
  19. family=reuse;policy=notmat;refs=4;consumer=100%
  20. family=reuse;policy=mat;refs=16;consumer=100%
  21. family=reuse;policy=mat;refs=1;consumer=1%
  22. family=bag;stmt=cte_joined_outer;threshold=2000
  23. family=reuse;policy=notmat;refs=4;consumer=1%
  24. family=reuse;policy=notmat;refs=16;consumer=100%
  25. family=reuse;policy=mat;refs=16;consumer=1%
  26. family=bag;stmt=cte_joined_outer_mat;threshold=20
  27. family=reuse;policy=notmat;refs=16;consumer=1%
  28. family=reuse;policy=mat;refs=1;consumer=100%
  29. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  30. family=reuse;policy=mat;refs=4;consumer=1%

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
| family=bag;stmt=cte_joined_outer;threshold=180000 | j05_cte_joined_outer | valid | `25ee840158f7c004` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20 | j05_cte_joined_outer | valid | `ce64b249bf8e67e3` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=2000 | j05_cte_joined_outer | valid | `f2c305426450532c` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | j05m_cte_joined_outer_mat | valid | `7582417914fc1237` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | j05m_cte_joined_outer_mat | valid | `52e7af1b1a454d2b` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | j05m_cte_joined_outer_mat | valid | `0d6685646a9fc1e4` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | j05n_cte_joined_outer_notmat | valid | `87d57c22e1ca51d2` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | j05n_cte_joined_outer_notmat | valid | `e4d16ea7ba5ab9c1` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | j05n_cte_joined_outer_notmat | valid | `9935b2adb78284b3` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=join_where;threshold=180000 | j01_join_where | valid | `67f0ff2ac4ad813d` |  |
| family=bag;stmt=join_where;threshold=20 | j01_join_where | valid | `fa739a6b7ef06856` |  |
| family=bag;stmt=join_where;threshold=2000 | j01_join_where | valid | `9ff2de9d95af1ff0` |  |
| family=reuse;policy=default;refs=16;consumer=1% | r_default_n16_c1 | valid | `3992e61f6b808571` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=16;consumer=100% | r_default_n16_c100 | valid | `67f45ccb2e3ba270` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=1;consumer=1% | r_default_n1_c1 | valid | `9dce1f5dbe17be5c` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=1;consumer=100% | r_default_n1_c100 | valid | `509e0017f52e40cb` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=4;consumer=1% | r_default_n4_c1 | valid | `66504facceaba213` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=100% | r_default_n4_c100 | valid | `d581fa7f5e08579b` | policy=default; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=1% | r_mat_n16_c1 | valid | `f2b237437475bc92` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=100% | r_mat_n16_c100 | valid | `153fd01083180b22` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=1% | r_mat_n1_c1 | valid | `cacba2210809d6a1` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=100% | r_mat_n1_c100 | valid | `86225bb09a69a5a0` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=1% | r_mat_n4_c1 | valid | `2efca80386d5adb0` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=100% | r_mat_n4_c100 | valid | `dae3aebfbb510def` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=notmat;refs=16;consumer=1% | r_notmat_n16_c1 | valid | `53f819dc7dd8767d` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `6e72a37535e7cad4` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `952d97bc8dffa024` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `d872aa899963cbd6` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `795c6aa7522efa68` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `6ff281c7be5c01d7` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, DIST, DEBUG, FORMAT TEXT` supported on yugabyte
- one host: this run used node budget 2 CPU 3g and client budget 1 CPU 2g, so budget comparisons are separate classes here, not independent hardware
- one client process with one fixed maintenance writer at 10 updates/s; readers are sequential; the reported sample latency includes retry and backoff time
- physical rank-index arm yb-range-rank; an engine whose planner ignores the index is a different physical design, not a faster one
- YugabyteDB DIST/DEBUG storage counters (rows and index rows scanned) are parsed per cell; RPC-level counters are not aggregated
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
