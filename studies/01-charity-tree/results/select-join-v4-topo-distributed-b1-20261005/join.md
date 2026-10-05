# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 30 cells, 30 valid, engine `yugabyte`, topology `cluster3`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-v4-topo-distributed-b1-20261005`
- run tag: `run/study-01/select-join-v4-topo-distributed-b1-20261005`
- environment: `host-zenbook-ux5406sa`
- repo commit: `ba821e7b686ac3adb1cda56dd45f15c8ae699ef8`
- engine version: `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit`
- database engine image: `docker.io/yugabytedb/yugabyte:2025.2.6.0-b111` @ `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`
- benchmark client image: `localhost/charitybench-select:select-join-v4-topo-distributed-b1-20261005` @ `b6ea3c007cf459112eaa8e9e85381a7dc25a7710ff07648fde3048f747df6ab0`
- inputs digest: `ca8e03004e7caf4a78d0441c47b3804975c22e3924cc1156f775731c60a23b24`
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
| yb-n1 | 169 | 2 |
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
| 20 (0.01%) | 3.2999 | missing | missing | missing | missing | missing | 3.6563 | 73.0584 | 3.5023 | missing |
| 2000 (1%) | 8.3302 | missing | missing | missing | missing | missing | 8.5606 | 78.3079 | 9.5904 | missing |
| 180000 (90%) | 87.7122 | missing | missing | missing | missing | missing | 87.2105 | 117.6337 | 83.8236 | missing |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | cte_joined_outer vs join_where | 3.2999 | 3.6563 | 1.1080 | false |
| 20 | cte_joined_outer_mat vs join_where | 3.2999 | 73.0584 | 22.1399 | false |
| 20 | cte_joined_outer_notmat vs join_where | 3.2999 | 3.5023 | 1.0613 | false |
| 2000 | cte_joined_outer vs join_where | 8.3302 | 8.5606 | 1.0277 | false |
| 2000 | cte_joined_outer_mat vs join_where | 8.3302 | 78.3079 | 9.4005 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 8.3302 | 9.5904 | 1.1513 | false |
| 180000 | cte_joined_outer vs join_where | 87.7122 | 87.2105 | 0.9943 | false |
| 180000 | cte_joined_outer_mat vs join_where | 87.7122 | 117.6337 | 1.3411 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 87.7122 | 83.8236 | 0.9557 | false |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_joined_outer_mat vs cte_joined_outer | 3.6563 | 73.0584 | 19.9816 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 3.6563 | 3.5023 | 0.9579 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 8.5606 | 78.3079 | 9.1475 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 8.5606 | 9.5904 | 1.1203 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 87.2105 | 117.6337 | 1.3488 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 87.2105 | 83.8236 | 0.9612 | false | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 5.6389 | 1 | `0f04a154dba0131c` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 7.1104 | 1 | `07360b73648d558d` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | missing | | | | |
| default | 2 | 100% | missing | | | | |
| default | 4 | 1% | 8.1952 | 4 | `43d9f404ac60395b` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 9.0057 | 4 | `7eeccf280d353ba3` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | missing | | | | |
| default | 8 | 100% | missing | | | | |
| default | 16 | 1% | 9.0413 | 16 | `2b406b2e17f63f76` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 10.1480 | 16 | `d1981973c37fcee2` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 8.0929 | 1 | `a8948b96fdd7b096` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 8.4941 | 1 | `19a077a677f12d06` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | missing | | | | |
| mat | 2 | 100% | missing | | | | |
| mat | 4 | 1% | 8.7506 | 4 | `5f478eb6c41e4ddf` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 8.2004 | 4 | `d3ce985f45147a10` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | missing | | | | |
| mat | 8 | 100% | missing | | | | |
| mat | 16 | 1% | 8.8969 | 16 | `a6f6385c1ebed495` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 10.9469 | 16 | `31bb8ff8bb94a13d` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 5.8764 | 1 | `d36323795b050331` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 7.5990 | 1 | `e0e08ec08e14c8ac` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | missing | | | | |
| notmat | 2 | 100% | missing | | | | |
| notmat | 4 | 1% | 20.0919 | 4 | `aaa030cf261db65e` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 27.2655 | 4 | `1664a36048ad180d` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | missing | | | | |
| notmat | 8 | 100% | missing | | | | |
| notmat | 16 | 1% | 93.5290 | 16 | `f31a1db0642c07b8` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 107.6959 | 16 | `98aafe9202286ce1` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 1.4352 | false |
| notmat | 1 | 1% | 1.0421 | false |
| mat | 1 | 100% | 1.1946 | false |
| notmat | 1 | 100% | 1.0687 | false |
| mat | 4 | 1% | 1.0678 | false |
| notmat | 4 | 1% | 2.4517 | false |
| mat | 4 | 100% | 0.9106 | false |
| notmat | 4 | 100% | 3.0276 | false |
| mat | 16 | 1% | 0.9840 | false |
| notmat | 16 | 1% | 10.3447 | false |
| mat | 16 | 100% | 1.0787 | false |
| notmat | 16 | 100% | 10.6125 | false |

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

- client: nr_periods 347, nr_throttled 161, throttled_usec 6353847 | quota 100000 100000
, memory 2147483648

- yb-n1: nr_periods 697, nr_throttled 26, throttled_usec 1440554 | quota 200000 100000, memory 3221225472
- yb-n2: nr_periods 691, nr_throttled 6, throttled_usec 1142087 | quota 200000 100000, memory 3221225472
- yb-n3: nr_periods 684, nr_throttled 5, throttled_usec 126167 | quota 200000 100000, memory 3221225472

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
| family=bag;stmt=join_where;threshold=180000 | valid | false | true | 202050 | 0 | 0/0 | 0 | 1 |
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
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 76.1330 | 87.2105 | 200.0101 | 5 | 11.5 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 3.0684 | 3.6563 | 4.1631 | 5 | 273.5 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 8.0351 | 8.5606 | 10.8500 | 5 | 116.8 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 116.2316 | 117.6337 | 134.7595 | 5 | 8.5 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 70.0094 | 73.0584 | 76.3852 | 5 | 13.7 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 74.2400 | 78.3079 | 80.6776 | 5 | 12.8 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 81.2844 | 83.8236 | 101.4608 | 5 | 11.9 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 3.2034 | 3.5023 | 6.6423 | 5 | 285.5 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 8.5033 | 9.5904 | 14.2156 | 5 | 104.3 |
| family=bag;stmt=join_where;threshold=180000 | valid | 83.1314 | 87.7122 | 171.6445 | 5 | 11.4 |
| family=bag;stmt=join_where;threshold=20 | valid | 2.7816 | 3.2999 | 34.0332 | 5 | 303.0 |
| family=bag;stmt=join_where;threshold=2000 | valid | 7.8011 | 8.3302 | 10.3160 | 5 | 120.0 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 8.6409 | 9.0413 | 10.6447 | 5 | 110.6 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 9.5214 | 10.1480 | 11.0036 | 5 | 98.5 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 5.2718 | 5.6389 | 6.2454 | 5 | 177.3 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 6.7274 | 7.1104 | 7.3848 | 5 | 140.6 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 7.9903 | 8.1952 | 9.8222 | 5 | 122.0 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 7.9119 | 9.0057 | 294.9027 | 5 | 111.0 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 7.8849 | 8.8969 | 9.5080 | 5 | 112.4 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 9.3216 | 10.9469 | 11.2567 | 5 | 91.3 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 7.4821 | 8.0929 | 8.2583 | 5 | 123.6 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 8.1118 | 8.4941 | 9.0261 | 5 | 117.7 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 7.7688 | 8.7506 | 12.2413 | 5 | 114.3 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 7.9429 | 8.2004 | 10.0488 | 5 | 121.9 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 81.3690 | 93.5290 | 170.8736 | 5 | 10.7 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 92.2288 | 107.6959 | 113.9347 | 5 | 9.3 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 5.5865 | 5.8764 | 6.8755 | 5 | 170.2 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 6.8520 | 7.5990 | 8.9540 | 5 | 131.6 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 18.3958 | 20.0919 | 22.8884 | 5 | 49.8 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 26.3167 | 27.2655 | 28.4358 | 5 | 36.7 |

## Matched blended workload

- phase: `blended`; target offered rate: 10 updates/s
- offered 169, committed 169, failed 0 over 16.881 s
- offered rate 10.011/s, completed rate 10.011/s
- writer operations per endpoint: [57 56 56]
- reader operations overlapping a committed write: 66 of 211
- reader retries of transient conflicts: 1
- per key: 169 keys, min 1, max 1 per key
- ledger: 169 keys reconciled, 0 mismatch
- update latency: avg 7.269 ms, min 2.717, max 193.137 over 169

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42890283`, 5 interleaved rounds
  1. family=bag;stmt=join_where;threshold=2000
  2. family=reuse;policy=notmat;refs=4;consumer=1%
  3. family=reuse;policy=notmat;refs=16;consumer=100%
  4. family=reuse;policy=default;refs=4;consumer=100%
  5. family=reuse;policy=default;refs=4;consumer=1%
  6. family=bag;stmt=cte_joined_outer;threshold=180000
  7. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  8. family=reuse;policy=default;refs=1;consumer=100%
  9. family=bag;stmt=cte_joined_outer_mat;threshold=20
  10. family=reuse;policy=default;refs=1;consumer=1%
  11. family=reuse;policy=notmat;refs=4;consumer=100%
  12. family=reuse;policy=notmat;refs=1;consumer=100%
  13. family=reuse;policy=notmat;refs=1;consumer=1%
  14. family=reuse;policy=default;refs=16;consumer=100%
  15. family=reuse;policy=notmat;refs=16;consumer=1%
  16. family=bag;stmt=cte_joined_outer;threshold=20
  17. family=reuse;policy=mat;refs=4;consumer=100%
  18. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  19. family=reuse;policy=mat;refs=1;consumer=100%
  20. family=bag;stmt=join_where;threshold=180000
  21. family=reuse;policy=mat;refs=16;consumer=100%
  22. family=reuse;policy=mat;refs=4;consumer=1%
  23. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  24. family=bag;stmt=cte_joined_outer;threshold=2000
  25. family=reuse;policy=default;refs=16;consumer=1%
  26. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  27. family=reuse;policy=mat;refs=1;consumer=1%
  28. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  29. family=reuse;policy=mat;refs=16;consumer=1%
  30. family=bag;stmt=join_where;threshold=20

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
| family=bag;stmt=cte_joined_outer;threshold=180000 | j05_cte_joined_outer | valid | `bd07152a21fe4b17` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20 | j05_cte_joined_outer | valid | `e64923b6b6772586` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=2000 | j05_cte_joined_outer | valid | `18b972b339ecca6d` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | j05m_cte_joined_outer_mat | valid | `9e0ad73535caffcf` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | j05m_cte_joined_outer_mat | valid | `cf9c0406585bec04` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | j05m_cte_joined_outer_mat | valid | `86927954729bedba` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | j05n_cte_joined_outer_notmat | valid | `c868adfe08409ff4` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | j05n_cte_joined_outer_notmat | valid | `ab4934cb707dbede` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | j05n_cte_joined_outer_notmat | valid | `769899bf3355c45e` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=join_where;threshold=180000 | j01_join_where | valid | `7428623f62da2eb3` |  |
| family=bag;stmt=join_where;threshold=20 | j01_join_where | valid | `e33da071104c620e` |  |
| family=bag;stmt=join_where;threshold=2000 | j01_join_where | valid | `90e6002d90577113` |  |
| family=reuse;policy=default;refs=16;consumer=1% | r_default_n16_c1 | valid | `2b406b2e17f63f76` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=16;consumer=100% | r_default_n16_c100 | valid | `d1981973c37fcee2` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=1;consumer=1% | r_default_n1_c1 | valid | `0f04a154dba0131c` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=1;consumer=100% | r_default_n1_c100 | valid | `07360b73648d558d` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=4;consumer=1% | r_default_n4_c1 | valid | `43d9f404ac60395b` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=100% | r_default_n4_c100 | valid | `7eeccf280d353ba3` | policy=default; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=1% | r_mat_n16_c1 | valid | `a6f6385c1ebed495` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=100% | r_mat_n16_c100 | valid | `31bb8ff8bb94a13d` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=1% | r_mat_n1_c1 | valid | `a8948b96fdd7b096` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=100% | r_mat_n1_c100 | valid | `19a077a677f12d06` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=1% | r_mat_n4_c1 | valid | `5f478eb6c41e4ddf` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=100% | r_mat_n4_c100 | valid | `d3ce985f45147a10` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=notmat;refs=16;consumer=1% | r_notmat_n16_c1 | valid | `f31a1db0642c07b8` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `98aafe9202286ce1` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `d36323795b050331` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `e0e08ec08e14c8ac` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `aaa030cf261db65e` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `1664a36048ad180d` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, DIST, DEBUG, FORMAT TEXT` supported on yugabyte
- distributed placement gate: three tablets confirmed by the engine, leader spread confirmed by the runner
- one host: this run used node budget 2 CPU 3g and client budget 1 CPU 2g, so budget comparisons are separate classes here, not independent hardware
- one client process with one fixed maintenance writer at 10 updates/s; readers are sequential; the reported sample latency includes retry and backoff time
- cluster nodes are containers on the same host, so these are not network or independent-hardware measurements; container colocation is not data-colocation evidence
- readers round-robin across 3 endpoints; the writer's endpoint policy is recorded separately
- physical rank-index arm yb-range-rank; an engine whose planner ignores the index is a different physical design, not a faster one
- YugabyteDB DIST/DEBUG storage counters (rows and index rows scanned) are parsed per cell; RPC-level counters are not aggregated
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
