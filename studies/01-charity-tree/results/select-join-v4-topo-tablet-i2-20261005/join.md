# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 30 cells, 30 valid, engine `yugabyte`, topology `cluster3`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-v4-topo-tablet-i2-20261005`
- run tag: `run/study-01/select-join-v4-topo-tablet-i2-20261005`
- environment: `host-zenbook-ux5406sa`
- repo commit: `ba821e7b686ac3adb1cda56dd45f15c8ae699ef8`
- engine version: `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit`
- database engine image: `docker.io/yugabytedb/yugabyte:2025.2.6.0-b111` @ `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`
- benchmark client image: `localhost/charitybench-select:select-join-v4-topo-tablet-i2-20261005` @ `b6ea3c007cf459112eaa8e9e85381a7dc25a7710ff07648fde3048f747df6ab0`
- inputs digest: `020584f5fd06c11a3b471d4e25265c7ac45ecf5d310ed241e5e9ff4a7d1395b1`
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
- setting `placement`: noncolocated-single-tablet
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
| j01_join_where | valid | true |  |
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
| 20 (0.01%) | 2.7500 | missing | missing | missing | missing | missing | 2.6650 | 84.2837 | 15.2788 | missing |
| 2000 (1%) | 10.9394 | missing | missing | missing | missing | missing | 8.9196 | 122.3975 | 11.6093 | missing |
| 180000 (90%) | 146.7331 | missing | missing | missing | missing | missing | 140.2097 | 105.2482 | 107.6607 | missing |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | cte_joined_outer vs join_where | 2.7500 | 2.6650 | 0.9691 | false |
| 20 | cte_joined_outer_mat vs join_where | 2.7500 | 84.2837 | 30.6491 | false |
| 20 | cte_joined_outer_notmat vs join_where | 2.7500 | 15.2788 | 5.5560 | false |
| 2000 | cte_joined_outer vs join_where | 10.9394 | 8.9196 | 0.8154 | false |
| 2000 | cte_joined_outer_mat vs join_where | 10.9394 | 122.3975 | 11.1887 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 10.9394 | 11.6093 | 1.0612 | false |
| 180000 | cte_joined_outer vs join_where | 146.7331 | 140.2097 | 0.9555 | false |
| 180000 | cte_joined_outer_mat vs join_where | 146.7331 | 105.2482 | 0.7173 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 146.7331 | 107.6607 | 0.7337 | false |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_joined_outer_mat vs cte_joined_outer | 2.6650 | 84.2837 | 31.6261 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 2.6650 | 15.2788 | 5.7331 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 8.9196 | 122.3975 | 13.7224 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 8.9196 | 11.6093 | 1.3016 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 140.2097 | 105.2482 | 0.7506 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 140.2097 | 107.6607 | 0.7679 | false | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 6.9484 | 1 | `4bc03410a1f3e027` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 8.6243 | 1 | `71b781f7dde65b0d` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | missing | | | | |
| default | 2 | 100% | missing | | | | |
| default | 4 | 1% | 10.0216 | 4 | `4fb1644a8672c650` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 9.7735 | 4 | `719a4ff88ffe5e20` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | missing | | | | |
| default | 8 | 100% | missing | | | | |
| default | 16 | 1% | 10.0465 | 16 | `702064b628b13f38` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 10.8146 | 16 | `19863066ff84c3d6` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 9.5220 | 1 | `6e7d4087d880f8d0` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 8.4367 | 1 | `e91ef029eae871d2` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | missing | | | | |
| mat | 2 | 100% | missing | | | | |
| mat | 4 | 1% | 9.4803 | 4 | `38f44f8a55b58140` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 8.2991 | 4 | `4f7b8a7c4b4d8e64` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | missing | | | | |
| mat | 8 | 100% | missing | | | | |
| mat | 16 | 1% | 10.2975 | 16 | `224413dda96ca549` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 10.6065 | 16 | `429ed81a1b050ec9` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 8.3977 | 1 | `28f72e8792753e70` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 9.5104 | 1 | `dd2a0a1a189d043c` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | missing | | | | |
| notmat | 2 | 100% | missing | | | | |
| notmat | 4 | 1% | 22.9815 | 4 | `461c3cc88d991e82` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 33.3348 | 4 | `a4912a0b44c952d1` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | missing | | | | |
| notmat | 8 | 100% | missing | | | | |
| notmat | 16 | 1% | 71.3659 | 16 | `fd5d39603fef909e` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 95.7612 | 16 | `3b4dac0ad0783817` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 1.3704 | false |
| notmat | 1 | 1% | 1.2086 | false |
| mat | 1 | 100% | 0.9782 | false |
| notmat | 1 | 100% | 1.1027 | false |
| mat | 4 | 1% | 0.9460 | false |
| notmat | 4 | 1% | 2.2932 | false |
| mat | 4 | 100% | 0.8491 | false |
| notmat | 4 | 100% | 3.4107 | false |
| mat | 16 | 1% | 1.0250 | false |
| notmat | 16 | 1% | 7.1036 | false |
| mat | 16 | 100% | 0.9808 | false |
| notmat | 16 | 100% | 8.8548 | false |

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
| charity | yb-n3:9100 |
| person | yb-n3:9100 |
| donation | yb-n3:9100 |
| donation_rank_idx | yb-n3:9100 |
| donation_person_idx | yb-n3:9100 |

Effective sharding flag (from the engine): `ysql_num_shards_per_tserver=2`. A relation with fewer tablets than this setting was created without an explicit split; the per-table `num_tablets` above is the authority.

## Throttling (cgroup `cpu.stat` deltas)

- client: nr_periods 307, nr_throttled 116, throttled_usec 4947765 | quota 100000 100000
, memory 2147483648

- yb-n1: nr_periods 586, nr_throttled 4, throttled_usec 426573 | quota 200000 100000, memory 3221225472
- yb-n2: nr_periods 577, nr_throttled 5, throttled_usec 634129 | quota 200000 100000, memory 3221225472
- yb-n3: nr_periods 572, nr_throttled 10, throttled_usec 226197 | quota 200000 100000, memory 3221225472

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
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 135.9242 | 140.2097 | 156.5984 | 5 | 7.1 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 2.4290 | 2.6650 | 2.8394 | 5 | 375.2 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 8.6095 | 8.9196 | 9.6804 | 5 | 112.1 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 98.8099 | 105.2482 | 116.0795 | 5 | 9.5 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 80.8283 | 84.2837 | 87.2775 | 5 | 11.9 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 118.9668 | 122.3975 | 125.8827 | 5 | 8.2 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 103.6966 | 107.6607 | 113.1568 | 5 | 9.3 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 2.8007 | 15.2788 | 28.2454 | 5 | 65.5 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 10.6789 | 11.6093 | 12.2085 | 5 | 86.1 |
| family=bag;stmt=join_where;threshold=180000 | valid | 137.6945 | 146.7331 | 158.1263 | 5 | 6.8 |
| family=bag;stmt=join_where;threshold=20 | valid | 2.4474 | 2.7500 | 2.9673 | 5 | 363.6 |
| family=bag;stmt=join_where;threshold=2000 | valid | 10.8708 | 10.9394 | 11.4845 | 5 | 91.4 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 9.8783 | 10.0465 | 11.5475 | 5 | 99.5 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 10.2391 | 10.8146 | 11.3986 | 5 | 92.5 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 6.4899 | 6.9484 | 6.9648 | 5 | 143.9 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 8.0935 | 8.6243 | 8.7738 | 5 | 116.0 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 9.3665 | 10.0216 | 11.0831 | 5 | 99.8 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 8.9935 | 9.7735 | 10.5410 | 5 | 102.3 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 10.0154 | 10.2975 | 10.8100 | 5 | 97.1 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 10.2152 | 10.6065 | 11.1340 | 5 | 94.3 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 9.1480 | 9.5220 | 9.6813 | 5 | 105.0 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 7.9007 | 8.4367 | 8.9948 | 5 | 118.5 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 9.1215 | 9.4803 | 9.9428 | 5 | 105.5 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 8.0371 | 8.2991 | 8.8681 | 5 | 120.5 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 70.2010 | 71.3659 | 99.0338 | 5 | 14.0 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 89.2658 | 95.7612 | 103.9529 | 5 | 10.4 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 7.8359 | 8.3977 | 8.7264 | 5 | 119.1 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 9.3688 | 9.5104 | 9.5375 | 5 | 105.1 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 22.1155 | 22.9815 | 29.2365 | 5 | 43.5 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 32.3106 | 33.3348 | 35.0946 | 5 | 30.0 |

## Matched blended workload

- phase: `isolated`; target offered rate: 10 updates/s
- offered 0, committed 0, failed 0 over 0.000 s
- offered rate 0.000/s, completed rate 0.000/s
- per key: no keys touched
- ledger: 0 keys reconciled, 0 mismatch

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42376054`, 5 interleaved rounds
  1. family=bag;stmt=cte_joined_outer;threshold=20
  2. family=bag;stmt=cte_joined_outer;threshold=180000
  3. family=reuse;policy=notmat;refs=16;consumer=1%
  4. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  5. family=reuse;policy=notmat;refs=16;consumer=100%
  6. family=reuse;policy=notmat;refs=4;consumer=1%
  7. family=reuse;policy=default;refs=16;consumer=100%
  8. family=reuse;policy=default;refs=4;consumer=1%
  9. family=reuse;policy=default;refs=1;consumer=100%
  10. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  11. family=bag;stmt=join_where;threshold=2000
  12. family=bag;stmt=cte_joined_outer_mat;threshold=20
  13. family=reuse;policy=notmat;refs=1;consumer=1%
  14. family=bag;stmt=join_where;threshold=180000
  15. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  16. family=reuse;policy=notmat;refs=4;consumer=100%
  17. family=reuse;policy=default;refs=4;consumer=100%
  18. family=reuse;policy=default;refs=1;consumer=1%
  19. family=bag;stmt=join_where;threshold=20
  20. family=reuse;policy=notmat;refs=1;consumer=100%
  21. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  22. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  23. family=reuse;policy=mat;refs=1;consumer=1%
  24. family=reuse;policy=mat;refs=1;consumer=100%
  25. family=reuse;policy=default;refs=16;consumer=1%
  26. family=reuse;policy=mat;refs=16;consumer=1%
  27. family=bag;stmt=cte_joined_outer;threshold=2000
  28. family=reuse;policy=mat;refs=16;consumer=100%
  29. family=reuse;policy=mat;refs=4;consumer=1%
  30. family=reuse;policy=mat;refs=4;consumer=100%

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
| family=bag;stmt=cte_joined_outer;threshold=180000 | j05_cte_joined_outer | valid | `6214d9b0af0b6e32` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20 | j05_cte_joined_outer | valid | `4bb1f8e9391d097f` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=2000 | j05_cte_joined_outer | valid | `4ab2ecbedc9acbcd` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | j05m_cte_joined_outer_mat | valid | `234ecbef89f75896` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | j05m_cte_joined_outer_mat | valid | `155c73aa9082e7a7` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | j05m_cte_joined_outer_mat | valid | `bb7f239e2a9b298f` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | j05n_cte_joined_outer_notmat | valid | `98d80c0e8cd671b3` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | j05n_cte_joined_outer_notmat | valid | `54d6ff9e68fe2de3` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | j05n_cte_joined_outer_notmat | valid | `deee7f14565100d7` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=join_where;threshold=180000 | j01_join_where | valid | `5fd85201611c2ad3` |  |
| family=bag;stmt=join_where;threshold=20 | j01_join_where | valid | `a3f50f96abf3f88f` |  |
| family=bag;stmt=join_where;threshold=2000 | j01_join_where | valid | `43bcc244a35f54f6` |  |
| family=reuse;policy=default;refs=16;consumer=1% | r_default_n16_c1 | valid | `702064b628b13f38` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=16;consumer=100% | r_default_n16_c100 | valid | `19863066ff84c3d6` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=1;consumer=1% | r_default_n1_c1 | valid | `4bc03410a1f3e027` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=1;consumer=100% | r_default_n1_c100 | valid | `71b781f7dde65b0d` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=4;consumer=1% | r_default_n4_c1 | valid | `4fb1644a8672c650` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=100% | r_default_n4_c100 | valid | `719a4ff88ffe5e20` | policy=default; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=1% | r_mat_n16_c1 | valid | `224413dda96ca549` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=100% | r_mat_n16_c100 | valid | `429ed81a1b050ec9` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=1% | r_mat_n1_c1 | valid | `6e7d4087d880f8d0` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=100% | r_mat_n1_c100 | valid | `e91ef029eae871d2` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=1% | r_mat_n4_c1 | valid | `38f44f8a55b58140` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=100% | r_mat_n4_c100 | valid | `4f7b8a7c4b4d8e64` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=notmat;refs=16;consumer=1% | r_notmat_n16_c1 | valid | `fd5d39603fef909e` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `3b4dac0ad0783817` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `28f72e8792753e70` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `dd2a0a1a189d043c` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `461c3cc88d991e82` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `a4912a0b44c952d1` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, DIST, DEBUG, FORMAT TEXT` supported on yugabyte
- one host: this run used node budget 2 CPU 3g and client budget 1 CPU 2g, so budget comparisons are separate classes here, not independent hardware
- isolated phase: no concurrent writers in this run; the blended read/write phase is a separate labelled run
- cluster nodes are containers on the same host, so these are not network or independent-hardware measurements; container colocation is not data-colocation evidence
- readers round-robin across 3 endpoints; the writer's endpoint policy is recorded separately
- physical rank-index arm yb-range-rank; an engine whose planner ignores the index is a different physical design, not a faster one
- YugabyteDB DIST/DEBUG storage counters (rows and index rows scanned) are parsed per cell; RPC-level counters are not aggregated
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
