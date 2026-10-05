# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 30 cells, 30 valid, engine `yugabyte`, topology `cluster3`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-v4-topo-distributed-i3-20261005`
- run tag: `run/study-01/select-join-v4-topo-distributed-i3-20261005`
- environment: `host-zenbook-ux5406sa`
- repo commit: `ba821e7b686ac3adb1cda56dd45f15c8ae699ef8`
- engine version: `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit`
- database engine image: `docker.io/yugabytedb/yugabyte:2025.2.6.0-b111` @ `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`
- benchmark client image: `localhost/charitybench-select:select-join-v4-topo-distributed-i3-20261005` @ `b6ea3c007cf459112eaa8e9e85381a7dc25a7710ff07648fde3048f747df6ab0`
- inputs digest: `261a8263c1f684e8bb2ab17c74a9216b2f4c722b6c9cbc7dad47db99e13674dc`
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
| 20 (0.01%) | 3.1798 | missing | missing | missing | missing | missing | 2.6984 | 64.7019 | 3.1138 | missing |
| 2000 (1%) | 8.1008 | missing | missing | missing | missing | missing | 7.8671 | 65.3485 | 7.2819 | missing |
| 180000 (90%) | 77.0517 | missing | missing | missing | missing | missing | 77.9346 | 100.8829 | 74.6723 | missing |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | cte_joined_outer vs join_where | 3.1798 | 2.6984 | 0.8486 | false |
| 20 | cte_joined_outer_mat vs join_where | 3.1798 | 64.7019 | 20.3477 | false |
| 20 | cte_joined_outer_notmat vs join_where | 3.1798 | 3.1138 | 0.9792 | false |
| 2000 | cte_joined_outer vs join_where | 8.1008 | 7.8671 | 0.9712 | false |
| 2000 | cte_joined_outer_mat vs join_where | 8.1008 | 65.3485 | 8.0669 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 8.1008 | 7.2819 | 0.8989 | false |
| 180000 | cte_joined_outer vs join_where | 77.0517 | 77.9346 | 1.0115 | false |
| 180000 | cte_joined_outer_mat vs join_where | 77.0517 | 100.8829 | 1.3093 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 77.0517 | 74.6723 | 0.9691 | false |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_joined_outer_mat vs cte_joined_outer | 2.6984 | 64.7019 | 23.9778 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 2.6984 | 3.1138 | 1.1539 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 7.8671 | 65.3485 | 8.3065 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 7.8671 | 7.2819 | 0.9256 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 77.9346 | 100.8829 | 1.2945 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 77.9346 | 74.6723 | 0.9581 | false | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 5.4597 | 1 | `a0a6e7daddd2e810` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 7.2126 | 1 | `b44d12d4c3fd0b2c` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | missing | | | | |
| default | 2 | 100% | missing | | | | |
| default | 4 | 1% | 7.4319 | 4 | `09a86545e9766751` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 7.1322 | 4 | `467e0af039c24a4e` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | missing | | | | |
| default | 8 | 100% | missing | | | | |
| default | 16 | 1% | 8.1518 | 16 | `7f63f88d07a2aba6` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 8.8872 | 16 | `e0c4ebdd5f83afae` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 6.8174 | 1 | `635f0372eef5c74d` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 7.4503 | 1 | `c405b80b32ea414d` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | missing | | | | |
| mat | 2 | 100% | missing | | | | |
| mat | 4 | 1% | 6.9621 | 4 | `448055a2aebba1e2` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 6.7784 | 4 | `11e441065cda4994` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | missing | | | | |
| mat | 8 | 100% | missing | | | | |
| mat | 16 | 1% | 8.4047 | 16 | `f3bc6a1c6729d91a` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 8.2484 | 16 | `21032330647114e5` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 5.0769 | 1 | `67abfae2ec9dee8c` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 6.4237 | 1 | `29fd0d1fd25c4d8b` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | missing | | | | |
| notmat | 2 | 100% | missing | | | | |
| notmat | 4 | 1% | 15.3370 | 4 | `d1e226a2253d515d` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 25.3041 | 4 | `687ae5430e62407e` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | missing | | | | |
| notmat | 8 | 100% | missing | | | | |
| notmat | 16 | 1% | 76.1795 | 16 | `b951bc2bcf9c88da` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 102.1203 | 16 | `bcec7b0806dd9ba8` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 1.2487 | false |
| notmat | 1 | 1% | 0.9299 | false |
| mat | 1 | 100% | 1.0330 | false |
| notmat | 1 | 100% | 0.8906 | false |
| mat | 4 | 1% | 0.9368 | false |
| notmat | 4 | 1% | 2.0637 | false |
| mat | 4 | 100% | 0.9504 | false |
| notmat | 4 | 100% | 3.5479 | false |
| mat | 16 | 1% | 1.0310 | false |
| notmat | 16 | 1% | 9.3451 | false |
| mat | 16 | 100% | 0.9281 | false |
| notmat | 16 | 100% | 11.4907 | false |

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

- client: nr_periods 292, nr_throttled 133, throttled_usec 5526846 | quota 100000 100000
, memory 2147483648

- yb-n1: nr_periods 618, nr_throttled 4, throttled_usec 646315 | quota 200000 100000, memory 3221225472
- yb-n2: nr_periods 610, nr_throttled 7, throttled_usec 934653 | quota 200000 100000, memory 3221225472
- yb-n3: nr_periods 608, nr_throttled 5, throttled_usec 264570 | quota 200000 100000, memory 3221225472

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
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 72.3113 | 77.9346 | 95.0828 | 5 | 12.8 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 2.5199 | 2.6984 | 3.0751 | 5 | 370.6 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 7.3149 | 7.8671 | 8.4979 | 5 | 127.1 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 99.6342 | 100.8829 | 102.8847 | 5 | 9.9 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 63.0568 | 64.7019 | 80.7975 | 5 | 15.5 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 64.1036 | 65.3485 | 75.6935 | 5 | 15.3 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 73.8257 | 74.6723 | 85.9201 | 5 | 13.4 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 2.6728 | 3.1138 | 3.5446 | 5 | 321.2 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 7.2511 | 7.2819 | 8.0791 | 5 | 137.3 |
| family=bag;stmt=join_where;threshold=180000 | valid | 74.5337 | 77.0517 | 83.0501 | 5 | 13.0 |
| family=bag;stmt=join_where;threshold=20 | valid | 2.8511 | 3.1798 | 3.4963 | 5 | 314.5 |
| family=bag;stmt=join_where;threshold=2000 | valid | 7.5346 | 8.1008 | 9.3505 | 5 | 123.4 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 7.1609 | 8.1518 | 8.7445 | 5 | 122.7 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 8.2690 | 8.8872 | 9.3260 | 5 | 112.5 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 4.9022 | 5.4597 | 5.8911 | 5 | 183.2 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 6.6276 | 7.2126 | 7.7869 | 5 | 138.6 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 6.6785 | 7.4319 | 8.4646 | 5 | 134.6 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 6.5571 | 7.1322 | 7.5097 | 5 | 140.2 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 7.6306 | 8.4047 | 8.9804 | 5 | 119.0 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 7.6815 | 8.2484 | 8.4814 | 5 | 121.2 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 6.5215 | 6.8174 | 8.3746 | 5 | 146.7 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 6.9379 | 7.4503 | 7.9233 | 5 | 134.2 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 6.2680 | 6.9621 | 8.7167 | 5 | 143.6 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 6.4655 | 6.7784 | 7.6156 | 5 | 147.5 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 73.6964 | 76.1795 | 94.5340 | 5 | 13.1 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 89.7541 | 102.1203 | 109.3751 | 5 | 9.8 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 4.8086 | 5.0769 | 5.5436 | 5 | 197.0 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 6.2161 | 6.4237 | 8.3156 | 5 | 155.7 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 14.6778 | 15.3370 | 17.8744 | 5 | 65.2 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 23.5915 | 25.3041 | 29.5622 | 5 | 39.5 |

## Matched blended workload

- phase: `isolated`; target offered rate: 10 updates/s
- offered 0, committed 0, failed 0 over 0.000 s
- offered rate 0.000/s, completed rate 0.000/s
- per key: no keys touched
- ledger: 0 keys reconciled, 0 mismatch

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42623580`, 5 interleaved rounds
  1. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  2. family=reuse;policy=default;refs=4;consumer=1%
  3. family=reuse;policy=default;refs=16;consumer=1%
  4. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  5. family=reuse;policy=mat;refs=16;consumer=1%
  6. family=reuse;policy=mat;refs=16;consumer=100%
  7. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  8. family=reuse;policy=notmat;refs=4;consumer=100%
  9. family=reuse;policy=mat;refs=4;consumer=100%
  10. family=reuse;policy=notmat;refs=1;consumer=1%
  11. family=bag;stmt=cte_joined_outer;threshold=20
  12. family=bag;stmt=join_where;threshold=180000
  13. family=bag;stmt=join_where;threshold=20
  14. family=reuse;policy=default;refs=4;consumer=100%
  15. family=reuse;policy=notmat;refs=16;consumer=1%
  16. family=reuse;policy=mat;refs=1;consumer=100%
  17. family=bag;stmt=cte_joined_outer;threshold=2000
  18. family=reuse;policy=notmat;refs=1;consumer=100%
  19. family=reuse;policy=default;refs=16;consumer=100%
  20. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  21. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  22. family=reuse;policy=default;refs=1;consumer=1%
  23. family=reuse;policy=mat;refs=4;consumer=1%
  24. family=reuse;policy=mat;refs=1;consumer=1%
  25. family=bag;stmt=cte_joined_outer_mat;threshold=20
  26. family=bag;stmt=join_where;threshold=2000
  27. family=reuse;policy=notmat;refs=16;consumer=100%
  28. family=bag;stmt=cte_joined_outer;threshold=180000
  29. family=reuse;policy=default;refs=1;consumer=100%
  30. family=reuse;policy=notmat;refs=4;consumer=1%

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
| family=bag;stmt=cte_joined_outer;threshold=180000 | j05_cte_joined_outer | valid | `d64ca9a77511d765` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20 | j05_cte_joined_outer | valid | `8782bdaa3a17a7d0` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=2000 | j05_cte_joined_outer | valid | `48bb7853e7936163` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | j05m_cte_joined_outer_mat | valid | `3765efe1991ee355` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | j05m_cte_joined_outer_mat | valid | `ab282fb1a02adb2e` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | j05m_cte_joined_outer_mat | valid | `c91636729318f132` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | j05n_cte_joined_outer_notmat | valid | `535c599218d3ad01` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | j05n_cte_joined_outer_notmat | valid | `0351d7f59bd588a2` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | j05n_cte_joined_outer_notmat | valid | `819b664a61ec27f5` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=join_where;threshold=180000 | j01_join_where | valid | `d631854eb91b329d` |  |
| family=bag;stmt=join_where;threshold=20 | j01_join_where | valid | `f02dc867164bd370` |  |
| family=bag;stmt=join_where;threshold=2000 | j01_join_where | valid | `bf83918ab74da44f` |  |
| family=reuse;policy=default;refs=16;consumer=1% | r_default_n16_c1 | valid | `7f63f88d07a2aba6` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=16;consumer=100% | r_default_n16_c100 | valid | `e0c4ebdd5f83afae` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=1;consumer=1% | r_default_n1_c1 | valid | `a0a6e7daddd2e810` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=1;consumer=100% | r_default_n1_c100 | valid | `b44d12d4c3fd0b2c` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=4;consumer=1% | r_default_n4_c1 | valid | `09a86545e9766751` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=100% | r_default_n4_c100 | valid | `467e0af039c24a4e` | policy=default; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=1% | r_mat_n16_c1 | valid | `f3bc6a1c6729d91a` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=100% | r_mat_n16_c100 | valid | `21032330647114e5` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=1% | r_mat_n1_c1 | valid | `635f0372eef5c74d` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=100% | r_mat_n1_c100 | valid | `c405b80b32ea414d` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=1% | r_mat_n4_c1 | valid | `448055a2aebba1e2` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=100% | r_mat_n4_c100 | valid | `11e441065cda4994` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=notmat;refs=16;consumer=1% | r_notmat_n16_c1 | valid | `b951bc2bcf9c88da` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `bcec7b0806dd9ba8` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `67abfae2ec9dee8c` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `29fd0d1fd25c4d8b` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `d1e226a2253d515d` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `687ae5430e62407e` | policy=not_materialized; observed_cte_scan=false |

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
