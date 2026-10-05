# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 30 cells, 30 valid, engine `yugabyte`, topology `cluster3`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-v4-topo-colocated-b1-20261005`
- run tag: `run/study-01/select-join-v4-topo-colocated-b1-20261005`
- environment: `host-zenbook-ux5406sa`
- repo commit: `ba821e7b686ac3adb1cda56dd45f15c8ae699ef8`
- engine version: `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit`
- database engine image: `docker.io/yugabytedb/yugabyte:2025.2.6.0-b111` @ `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`
- benchmark client image: `localhost/charitybench-select:select-join-v4-topo-colocated-b1-20261005` @ `b6ea3c007cf459112eaa8e9e85381a7dc25a7710ff07648fde3048f747df6ab0`
- inputs digest: `93eb5f6087c05d62860f236404ceb449cc9526da746005adc3a396506d030522`
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
| 20 (0.01%) | 2.4015 | missing | missing | missing | missing | missing | 2.2153 | 107.8110 | 1.8738 | missing |
| 2000 (1%) | 7.4860 | missing | missing | missing | missing | missing | 7.2387 | 66.9229 | 7.0996 | missing |
| 180000 (90%) | 97.3977 | missing | missing | missing | missing | missing | 131.4460 | 99.0774 | 137.6312 | missing |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | cte_joined_outer vs join_where | 2.4015 | 2.2153 | 0.9224 | false |
| 20 | cte_joined_outer_mat vs join_where | 2.4015 | 107.8110 | 44.8924 | false |
| 20 | cte_joined_outer_notmat vs join_where | 2.4015 | 1.8738 | 0.7803 | false |
| 2000 | cte_joined_outer vs join_where | 7.4860 | 7.2387 | 0.9670 | false |
| 2000 | cte_joined_outer_mat vs join_where | 7.4860 | 66.9229 | 8.9398 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 7.4860 | 7.0996 | 0.9484 | false |
| 180000 | cte_joined_outer vs join_where | 97.3977 | 131.4460 | 1.3496 | false |
| 180000 | cte_joined_outer_mat vs join_where | 97.3977 | 99.0774 | 1.0172 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 97.3977 | 137.6312 | 1.4131 | false |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_joined_outer_mat vs cte_joined_outer | 2.2153 | 107.8110 | 48.6666 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 2.2153 | 1.8738 | 0.8459 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 7.2387 | 66.9229 | 9.2451 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 7.2387 | 7.0996 | 0.9808 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 131.4460 | 99.0774 | 0.7538 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 131.4460 | 137.6312 | 1.0471 | false | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 1.4670 | 1 | `22fff1493a1d1287` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 6.5721 | 1 | `e624b4691fe9731f` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | missing | | | | |
| default | 2 | 100% | missing | | | | |
| default | 4 | 1% | 6.6974 | 4 | `b52aab85d2d026c4` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 7.3625 | 4 | `824ab44a469fc05d` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | missing | | | | |
| default | 8 | 100% | missing | | | | |
| default | 16 | 1% | 7.1613 | 16 | `08ed2a942185f705` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 8.0430 | 16 | `076912c3bb3d4805` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 6.6272 | 1 | `a9fbf01c4fc7fcf2` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 6.2434 | 1 | `2f437968a8488a8f` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | missing | | | | |
| mat | 2 | 100% | missing | | | | |
| mat | 4 | 1% | 6.8936 | 4 | `e0926d915f6bdd1d` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 6.2739 | 4 | `c354d93452367384` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | missing | | | | |
| mat | 8 | 100% | missing | | | | |
| mat | 16 | 1% | 7.5863 | 16 | `856d68ee593a75a9` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 7.7655 | 16 | `eebb278900e89711` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 1.6655 | 1 | `0c967d48610b70d7` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 5.7105 | 1 | `25b0d4710c33680c` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | missing | | | | |
| notmat | 2 | 100% | missing | | | | |
| notmat | 4 | 1% | 3.6780 | 4 | `c2bb60c75eee575c` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 23.8469 | 4 | `c5fe532a6af3b20c` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | missing | | | | |
| notmat | 8 | 100% | missing | | | | |
| notmat | 16 | 1% | 17.1661 | 16 | `cf90cb58b4e64f92` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 83.8295 | 16 | `e88ab93df96484fd` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 4.5175 | false |
| notmat | 1 | 1% | 1.1353 | false |
| mat | 1 | 100% | 0.9500 | false |
| notmat | 1 | 100% | 0.8689 | false |
| mat | 4 | 1% | 1.0293 | false |
| notmat | 4 | 1% | 0.5492 | false |
| mat | 4 | 100% | 0.8521 | false |
| notmat | 4 | 100% | 3.2390 | false |
| mat | 16 | 1% | 1.0593 | false |
| notmat | 16 | 1% | 2.3971 | false |
| mat | 16 | 100% | 0.9655 | false |
| notmat | 16 | 100% | 10.4227 | false |

## Physical arm

- `charity_pkey`: `CREATE UNIQUE INDEX charity_pkey ON selectjoin_domain.charity USING lsm (charity_id ASC)`
- `donation_person_idx`: `CREATE INDEX donation_person_idx ON selectjoin_domain.donation USING lsm (person_id ASC)`
- `donation_pkey`: `CREATE UNIQUE INDEX donation_pkey ON selectjoin_domain.donation USING lsm (donation_id ASC)`
- `donation_rank_idx`: `CREATE INDEX donation_rank_idx ON selectjoin_domain.donation USING lsm (rank ASC)`
- `person_pkey`: `CREATE UNIQUE INDEX person_pkey ON selectjoin_domain.person USING lsm (person_id ASC)`

- index arm `yb-range-rank`; placement `colocated`; node budget `2 CPU 3g`; client budget `1 CPU 2g`

| table | tablets | hash key columns | colocated | tablegroup | colocation id |
|---|---|---|---|---|---|
| selectjoin_domain.charity | 1 | 0 | true | 16388 | 457281939 |
| selectjoin_domain.person | 1 | 0 | true | 16388 | 3018457220 |
| selectjoin_domain.donation | 1 | 0 | true | 16388 | 816215557 |
| selectjoin_domain.donation_pkey | 0 | 0 | false | unknown | unknown |
| selectjoin_domain.donation_rank_idx | 1 | 0 | true | 16388 | 2998113465 |
| selectjoin_domain.donation_person_idx | 1 | 0 | true | 16388 | 1707047207 |

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

- client: nr_periods 268, nr_throttled 92, throttled_usec 3535021 | quota 100000 100000
, memory 2147483648

- yb-n1: nr_periods 514, nr_throttled 5, throttled_usec 737945 | quota 200000 100000, memory 3221225472
- yb-n2: nr_periods 507, nr_throttled 5, throttled_usec 584832 | quota 200000 100000, memory 3221225472
- yb-n3: nr_periods 500, nr_throttled 5, throttled_usec 107882 | quota 200000 100000, memory 3221225472

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
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 129.3603 | 131.4460 | 138.6578 | 5 | 7.6 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 2.1933 | 2.2153 | 2.9440 | 5 | 451.4 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 6.9679 | 7.2387 | 9.2189 | 5 | 138.1 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 94.9601 | 99.0774 | 100.1070 | 5 | 10.1 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 106.8283 | 107.8110 | 123.1372 | 5 | 9.3 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 65.8050 | 66.9229 | 73.4714 | 5 | 14.9 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 131.9118 | 137.6312 | 147.7229 | 5 | 7.3 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 1.7740 | 1.8738 | 1.9145 | 5 | 533.7 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 6.5502 | 7.0996 | 24.6684 | 5 | 140.9 |
| family=bag;stmt=join_where;threshold=180000 | valid | 92.2739 | 97.3977 | 100.1301 | 5 | 10.3 |
| family=bag;stmt=join_where;threshold=20 | valid | 2.1967 | 2.4015 | 2.5382 | 5 | 416.4 |
| family=bag;stmt=join_where;threshold=2000 | valid | 7.0927 | 7.4860 | 8.3373 | 5 | 133.6 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 7.0432 | 7.1613 | 7.6244 | 5 | 139.6 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 7.8282 | 8.0430 | 8.2014 | 5 | 124.3 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 1.3742 | 1.4670 | 1.9701 | 5 | 681.7 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 6.3757 | 6.5721 | 6.9613 | 5 | 152.2 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 6.4540 | 6.6974 | 7.7320 | 5 | 149.3 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 7.1660 | 7.3625 | 8.2971 | 5 | 135.8 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 6.9513 | 7.5863 | 7.6708 | 5 | 131.8 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 7.4708 | 7.7655 | 7.9399 | 5 | 128.8 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 6.1754 | 6.6272 | 6.8885 | 5 | 150.9 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 6.0236 | 6.2434 | 12.6384 | 5 | 160.2 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 6.8085 | 6.8936 | 7.5498 | 5 | 145.1 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 6.0286 | 6.2739 | 6.8140 | 5 | 159.4 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 16.0936 | 17.1661 | 18.1227 | 5 | 58.3 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 75.1043 | 83.8295 | 98.1918 | 5 | 11.9 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 1.5660 | 1.6655 | 1.9935 | 5 | 600.4 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 5.6405 | 5.7105 | 6.1154 | 5 | 175.1 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 3.5502 | 3.6780 | 4.1093 | 5 | 271.9 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 22.7661 | 23.8469 | 27.4644 | 5 | 41.9 |

## Matched blended workload

- phase: `blended`; target offered rate: 10 updates/s
- offered 122, committed 122, failed 0 over 12.127 s
- offered rate 10.060/s, completed rate 10.060/s
- writer operations per endpoint: [41 41 40]
- reader operations overlapping a committed write: 59 of 210
- reader retries of transient conflicts: 0
- per key: 122 keys, min 1, max 1 per key
- ledger: 122 keys reconciled, 0 mismatch
- update latency: avg 5.143 ms, min 2.316, max 64.894 over 122

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42940262`, 5 interleaved rounds
  1. family=reuse;policy=default;refs=4;consumer=100%
  2. family=reuse;policy=notmat;refs=1;consumer=1%
  3. family=reuse;policy=mat;refs=4;consumer=100%
  4. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  5. family=reuse;policy=default;refs=16;consumer=100%
  6. family=reuse;policy=notmat;refs=1;consumer=100%
  7. family=reuse;policy=mat;refs=16;consumer=100%
  8. family=reuse;policy=mat;refs=1;consumer=100%
  9. family=reuse;policy=default;refs=1;consumer=1%
  10. family=reuse;policy=default;refs=4;consumer=1%
  11. family=bag;stmt=cte_joined_outer_mat;threshold=20
  12. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  13. family=reuse;policy=notmat;refs=4;consumer=100%
  14. family=reuse;policy=mat;refs=16;consumer=1%
  15. family=reuse;policy=notmat;refs=4;consumer=1%
  16. family=reuse;policy=default;refs=1;consumer=100%
  17. family=reuse;policy=default;refs=16;consumer=1%
  18. family=reuse;policy=notmat;refs=16;consumer=100%
  19. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  20. family=reuse;policy=notmat;refs=16;consumer=1%
  21. family=bag;stmt=join_where;threshold=180000
  22. family=bag;stmt=join_where;threshold=2000
  23. family=bag;stmt=cte_joined_outer;threshold=2000
  24. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  25. family=bag;stmt=cte_joined_outer;threshold=20
  26. family=bag;stmt=cte_joined_outer;threshold=180000
  27. family=reuse;policy=mat;refs=4;consumer=1%
  28. family=reuse;policy=mat;refs=1;consumer=1%
  29. family=bag;stmt=join_where;threshold=20
  30. family=bag;stmt=cte_joined_outer_mat;threshold=180000

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
| family=bag;stmt=cte_joined_outer;threshold=180000 | j05_cte_joined_outer | valid | `f401317a4f995a97` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20 | j05_cte_joined_outer | valid | `78079c8fe7767f21` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=2000 | j05_cte_joined_outer | valid | `954c01152e48401f` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | j05m_cte_joined_outer_mat | valid | `2b284cb3658d7fa1` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | j05m_cte_joined_outer_mat | valid | `52dd23009635063f` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | j05m_cte_joined_outer_mat | valid | `b017c115a32dc7db` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | j05n_cte_joined_outer_notmat | valid | `d03b7e2ef2ce6bf9` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | j05n_cte_joined_outer_notmat | valid | `6524b3c9c7856a18` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | j05n_cte_joined_outer_notmat | valid | `cfc2f439850d61ec` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=join_where;threshold=180000 | j01_join_where | valid | `c868bb65b4cca916` |  |
| family=bag;stmt=join_where;threshold=20 | j01_join_where | valid | `ec3ac2b31152f1f9` |  |
| family=bag;stmt=join_where;threshold=2000 | j01_join_where | valid | `dd1831e90910d5ca` |  |
| family=reuse;policy=default;refs=16;consumer=1% | r_default_n16_c1 | valid | `08ed2a942185f705` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=16;consumer=100% | r_default_n16_c100 | valid | `076912c3bb3d4805` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=1;consumer=1% | r_default_n1_c1 | valid | `22fff1493a1d1287` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=1;consumer=100% | r_default_n1_c100 | valid | `e624b4691fe9731f` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=4;consumer=1% | r_default_n4_c1 | valid | `b52aab85d2d026c4` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=100% | r_default_n4_c100 | valid | `824ab44a469fc05d` | policy=default; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=1% | r_mat_n16_c1 | valid | `856d68ee593a75a9` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=100% | r_mat_n16_c100 | valid | `eebb278900e89711` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=1% | r_mat_n1_c1 | valid | `a9fbf01c4fc7fcf2` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=100% | r_mat_n1_c100 | valid | `2f437968a8488a8f` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=1% | r_mat_n4_c1 | valid | `e0926d915f6bdd1d` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=100% | r_mat_n4_c100 | valid | `c354d93452367384` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=notmat;refs=16;consumer=1% | r_notmat_n16_c1 | valid | `cf90cb58b4e64f92` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `e88ab93df96484fd` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `0c967d48610b70d7` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `25b0d4710c33680c` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `c2bb60c75eee575c` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `c5fe532a6af3b20c` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, DIST, DEBUG, FORMAT TEXT` supported on yugabyte
- one host: this run used node budget 2 CPU 3g and client budget 1 CPU 2g, so budget comparisons are separate classes here, not independent hardware
- one client process with one fixed maintenance writer at 10 updates/s; readers are sequential; the reported sample latency includes retry and backoff time
- cluster nodes are containers on the same host, so these are not network or independent-hardware measurements; container colocation is not data-colocation evidence
- readers round-robin across 3 endpoints; the writer's endpoint policy is recorded separately
- physical rank-index arm yb-range-rank; an engine whose planner ignores the index is a different physical design, not a faster one
- YugabyteDB DIST/DEBUG storage counters (rows and index rows scanned) are parsed per cell; RPC-level counters are not aggregated
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
