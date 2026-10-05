# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 30 cells, 30 valid, engine `yugabyte`, topology `cluster3`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-v4-smoke2-distributed-20261005`
- run tag: `run/study-01/select-join-v4-smoke2-distributed-20261005`
- environment: `host-zenbook-ux5406sa`
- repo commit: `ba821e7b686ac3adb1cda56dd45f15c8ae699ef8`
- engine version: `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit`
- database engine image: `docker.io/yugabytedb/yugabyte:2025.2.6.0-b111` @ `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`
- benchmark client image: `localhost/charitybench-select:select-join-v4-smoke2-distributed-20261005` @ `b6ea3c007cf459112eaa8e9e85381a7dc25a7710ff07648fde3048f747df6ab0`
- inputs digest: `599ef356c632c0e9cfebc12c8ce2f799240c988aef78541366bb0caa2869a4e1`
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
- setting `trials_per_cell`: 1
- setting `whole_run_count`: 1; no replicated winner in this file
- setting `writer_keys`: 1000
- setting `writer_target_rate_per_sec`: 10

### Endpoint distribution

A cluster arm must rotate its reader across nodes; a single seed endpoint would not be balanced cluster access.

| endpoint | operations | connections at end |
|---|---|---|
| yb-n1 | 128 | 1 |
| yb-n2 | 128 | 1 |
| yb-n3 | 128 | 1 |

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
| 20 (0.01%) | 2.9762 | missing | missing | missing | missing | missing | 3.3129 | 58.7822 | 2.4407 | missing |
| 2000 (1%) | 7.8757 | missing | missing | missing | missing | missing | 7.3229 | 62.3983 | 7.4036 | missing |
| 180000 (90%) | 75.6077 | missing | missing | missing | missing | missing | 91.1791 | 99.3677 | 69.3173 | missing |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | cte_joined_outer vs join_where | 2.9762 | 3.3129 | 1.1131 | false |
| 20 | cte_joined_outer_mat vs join_where | 2.9762 | 58.7822 | 19.7510 | false |
| 20 | cte_joined_outer_notmat vs join_where | 2.9762 | 2.4407 | 0.8201 | false |
| 2000 | cte_joined_outer vs join_where | 7.8757 | 7.3229 | 0.9298 | false |
| 2000 | cte_joined_outer_mat vs join_where | 7.8757 | 62.3983 | 7.9229 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 7.8757 | 7.4036 | 0.9401 | false |
| 180000 | cte_joined_outer vs join_where | 75.6077 | 91.1791 | 1.2060 | false |
| 180000 | cte_joined_outer_mat vs join_where | 75.6077 | 99.3677 | 1.3143 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 75.6077 | 69.3173 | 0.9168 | false |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_joined_outer_mat vs cte_joined_outer | 3.3129 | 58.7822 | 17.7435 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 3.3129 | 2.4407 | 0.7367 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 7.3229 | 62.3983 | 8.5210 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 7.3229 | 7.4036 | 1.0110 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 91.1791 | 99.3677 | 1.0898 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 91.1791 | 69.3173 | 0.7602 | false | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 5.7106 | 1 | `da34e49f30e7004d` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 5.8149 | 1 | `fbd04bfe14ee640c` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | missing | | | | |
| default | 2 | 100% | missing | | | | |
| default | 4 | 1% | 7.2531 | 4 | `1574aaa6423ce120` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 7.7525 | 4 | `3d8b9625e95c4b2c` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | missing | | | | |
| default | 8 | 100% | missing | | | | |
| default | 16 | 1% | 8.0072 | 16 | `4711a776a0b8db3d` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 8.3762 | 16 | `1439875dbe8852ff` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 6.7960 | 1 | `4d1e7c0ca11413f5` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 7.0689 | 1 | `65b759676825ae5a` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | missing | | | | |
| mat | 2 | 100% | missing | | | | |
| mat | 4 | 1% | 6.7543 | 4 | `af1ff5d31848bad2` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 7.5129 | 4 | `16eba81b02a4a531` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | missing | | | | |
| mat | 8 | 100% | missing | | | | |
| mat | 16 | 1% | 6.5286 | 16 | `62bf322fc7765dd5` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 7.3455 | 16 | `868951bce3d78b3d` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 4.3423 | 1 | `2168ed0ad37dc2ca` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 6.0881 | 1 | `be8b1ea4e73a4cde` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | missing | | | | |
| notmat | 2 | 100% | missing | | | | |
| notmat | 4 | 1% | 17.9524 | 4 | `e04d62ac104e21e1` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 25.8989 | 4 | `ecc159095eaec64d` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | missing | | | | |
| notmat | 8 | 100% | missing | | | | |
| notmat | 16 | 1% | 87.6750 | 16 | `082001e7de704d62` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 79.5749 | 16 | `f92677b24a09635f` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 1.1901 | false |
| notmat | 1 | 1% | 0.7604 | false |
| mat | 1 | 100% | 1.2157 | false |
| notmat | 1 | 100% | 1.0470 | false |
| mat | 4 | 1% | 0.9312 | false |
| notmat | 4 | 1% | 2.4751 | false |
| mat | 4 | 100% | 0.9691 | false |
| notmat | 4 | 100% | 3.3407 | false |
| mat | 16 | 1% | 0.8153 | false |
| notmat | 16 | 1% | 10.9495 | false |
| mat | 16 | 100% | 0.8769 | false |
| notmat | 16 | 100% | 9.5001 | false |

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

- client: nr_periods 201, nr_throttled 72, throttled_usec 2452818 | quota 100000 100000
, memory 2147483648

- yb-n1: nr_periods 538, nr_throttled 6, throttled_usec 1116201 | quota 200000 100000, memory 3221225472
- yb-n2: nr_periods 529, nr_throttled 6, throttled_usec 651133 | quota 200000 100000, memory 3221225472
- yb-n3: nr_periods 522, nr_throttled 0, throttled_usec 0 | quota 200000 100000, memory 3221225472

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
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 91.1791 | 91.1791 | 91.1791 | 1 | 11.0 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 3.3129 | 3.3129 | 3.3129 | 1 | 301.9 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 7.3229 | 7.3229 | 7.3229 | 1 | 136.6 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 99.3677 | 99.3677 | 99.3677 | 1 | 10.1 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 58.7822 | 58.7822 | 58.7822 | 1 | 17.0 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 62.3983 | 62.3983 | 62.3983 | 1 | 16.0 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 69.3173 | 69.3173 | 69.3173 | 1 | 14.4 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 2.4407 | 2.4407 | 2.4407 | 1 | 409.7 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 7.4036 | 7.4036 | 7.4036 | 1 | 135.1 |
| family=bag;stmt=join_where;threshold=180000 | valid | 75.6077 | 75.6077 | 75.6077 | 1 | 13.2 |
| family=bag;stmt=join_where;threshold=20 | valid | 2.9762 | 2.9762 | 2.9762 | 1 | 336.0 |
| family=bag;stmt=join_where;threshold=2000 | valid | 7.8757 | 7.8757 | 7.8757 | 1 | 127.0 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 8.0072 | 8.0072 | 8.0072 | 1 | 124.9 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 8.3762 | 8.3762 | 8.3762 | 1 | 119.4 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 5.7106 | 5.7106 | 5.7106 | 1 | 175.1 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 5.8149 | 5.8149 | 5.8149 | 1 | 172.0 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 7.2531 | 7.2531 | 7.2531 | 1 | 137.9 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 7.7525 | 7.7525 | 7.7525 | 1 | 129.0 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 6.5286 | 6.5286 | 6.5286 | 1 | 153.2 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 7.3455 | 7.3455 | 7.3455 | 1 | 136.1 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 6.7960 | 6.7960 | 6.7960 | 1 | 147.1 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 7.0689 | 7.0689 | 7.0689 | 1 | 141.5 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 6.7543 | 6.7543 | 6.7543 | 1 | 148.1 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 7.5129 | 7.5129 | 7.5129 | 1 | 133.1 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 87.6750 | 87.6750 | 87.6750 | 1 | 11.4 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 79.5749 | 79.5749 | 79.5749 | 1 | 12.6 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 4.3423 | 4.3423 | 4.3423 | 1 | 230.3 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 6.0881 | 6.0881 | 6.0881 | 1 | 164.3 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 17.9524 | 17.9524 | 17.9524 | 1 | 55.7 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 25.8989 | 25.8989 | 25.8989 | 1 | 38.6 |

## Matched blended workload

- phase: `isolated`; target offered rate: 10 updates/s
- offered 0, committed 0, failed 0 over 0.000 s
- offered rate 0.000/s, completed rate 0.000/s
- per key: no keys touched
- ledger: 0 keys reconciled, 0 mismatch

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42268374`, 1 interleaved rounds
  1. family=reuse;policy=notmat;refs=16;consumer=1%
  2. family=reuse;policy=notmat;refs=4;consumer=100%
  3. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  4. family=reuse;policy=notmat;refs=16;consumer=100%
  5. family=reuse;policy=default;refs=4;consumer=100%
  6. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  7. family=reuse;policy=notmat;refs=1;consumer=100%
  8. family=reuse;policy=mat;refs=16;consumer=1%
  9. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  10. family=reuse;policy=mat;refs=16;consumer=100%
  11. family=reuse;policy=notmat;refs=1;consumer=1%
  12. family=reuse;policy=default;refs=4;consumer=1%
  13. family=reuse;policy=notmat;refs=4;consumer=1%
  14. family=bag;stmt=cte_joined_outer;threshold=2000
  15. family=bag;stmt=join_where;threshold=180000
  16. family=bag;stmt=cte_joined_outer;threshold=180000
  17. family=reuse;policy=mat;refs=4;consumer=100%
  18. family=reuse;policy=default;refs=16;consumer=100%
  19. family=reuse;policy=default;refs=16;consumer=1%
  20. family=reuse;policy=mat;refs=1;consumer=1%
  21. family=bag;stmt=join_where;threshold=2000
  22. family=bag;stmt=cte_joined_outer;threshold=20
  23. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  24. family=reuse;policy=default;refs=1;consumer=1%
  25. family=reuse;policy=mat;refs=1;consumer=100%
  26. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  27. family=bag;stmt=join_where;threshold=20
  28. family=bag;stmt=cte_joined_outer_mat;threshold=20
  29. family=reuse;policy=default;refs=1;consumer=100%
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
| family=bag;stmt=cte_joined_outer;threshold=180000 | j05_cte_joined_outer | valid | `6c42006fe3c0d6e4` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20 | j05_cte_joined_outer | valid | `9433125bffd66064` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=2000 | j05_cte_joined_outer | valid | `a5e9bb7660ee37b0` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | j05m_cte_joined_outer_mat | valid | `3d40cf068966fd7a` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | j05m_cte_joined_outer_mat | valid | `a46d28794cde7601` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | j05m_cte_joined_outer_mat | valid | `32ca5fa794520b0e` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | j05n_cte_joined_outer_notmat | valid | `c7c7c2da0da7fa40` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | j05n_cte_joined_outer_notmat | valid | `a064705efbbbe878` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | j05n_cte_joined_outer_notmat | valid | `215dd4343fbcbb69` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=join_where;threshold=180000 | j01_join_where | valid | `727bb6a113f53ab8` |  |
| family=bag;stmt=join_where;threshold=20 | j01_join_where | valid | `62fc507a8c9fe590` |  |
| family=bag;stmt=join_where;threshold=2000 | j01_join_where | valid | `fc8ac2314e526066` |  |
| family=reuse;policy=default;refs=16;consumer=1% | r_default_n16_c1 | valid | `4711a776a0b8db3d` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=16;consumer=100% | r_default_n16_c100 | valid | `1439875dbe8852ff` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=1;consumer=1% | r_default_n1_c1 | valid | `da34e49f30e7004d` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=1;consumer=100% | r_default_n1_c100 | valid | `fbd04bfe14ee640c` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=4;consumer=1% | r_default_n4_c1 | valid | `1574aaa6423ce120` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=100% | r_default_n4_c100 | valid | `3d8b9625e95c4b2c` | policy=default; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=1% | r_mat_n16_c1 | valid | `62bf322fc7765dd5` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=100% | r_mat_n16_c100 | valid | `868951bce3d78b3d` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=1% | r_mat_n1_c1 | valid | `4d1e7c0ca11413f5` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=100% | r_mat_n1_c100 | valid | `65b759676825ae5a` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=1% | r_mat_n4_c1 | valid | `af1ff5d31848bad2` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=100% | r_mat_n4_c100 | valid | `16eba81b02a4a531` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=notmat;refs=16;consumer=1% | r_notmat_n16_c1 | valid | `082001e7de704d62` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `f92677b24a09635f` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `2168ed0ad37dc2ca` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `be8b1ea4e73a4cde` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `e04d62ac104e21e1` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `ecc159095eaec64d` | policy=not_materialized; observed_cte_scan=false |

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
