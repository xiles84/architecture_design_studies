# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 30 cells, 30 valid, engine `yugabyte`, topology `cluster3`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-v4-topo-colocated-i3-20261005`
- run tag: `run/study-01/select-join-v4-topo-colocated-i3-20261005`
- environment: `host-zenbook-ux5406sa`
- repo commit: `ba821e7b686ac3adb1cda56dd45f15c8ae699ef8`
- engine version: `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit`
- database engine image: `docker.io/yugabytedb/yugabyte:2025.2.6.0-b111` @ `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`
- benchmark client image: `localhost/charitybench-select:select-join-v4-topo-colocated-i3-20261005` @ `b6ea3c007cf459112eaa8e9e85381a7dc25a7710ff07648fde3048f747df6ab0`
- inputs digest: `efa57cff77782a8c5bad2db94a32b3437bb84a5ca5d51709bff14055937684fa`
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
| 20 (0.01%) | 1.6140 | missing | missing | missing | missing | missing | 2.2459 | 98.4898 | 2.5284 | missing |
| 2000 (1%) | 5.9581 | missing | missing | missing | missing | missing | 5.6248 | 64.9106 | 6.9368 | missing |
| 180000 (90%) | 92.9268 | missing | missing | missing | missing | missing | 126.7825 | 122.0927 | 91.2534 | missing |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | cte_joined_outer vs join_where | 1.6140 | 2.2459 | 1.3915 | false |
| 20 | cte_joined_outer_mat vs join_where | 1.6140 | 98.4898 | 61.0218 | false |
| 20 | cte_joined_outer_notmat vs join_where | 1.6140 | 2.5284 | 1.5665 | false |
| 2000 | cte_joined_outer vs join_where | 5.9581 | 5.6248 | 0.9441 | false |
| 2000 | cte_joined_outer_mat vs join_where | 5.9581 | 64.9106 | 10.8945 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 5.9581 | 6.9368 | 1.1643 | false |
| 180000 | cte_joined_outer vs join_where | 92.9268 | 126.7825 | 1.3643 | false |
| 180000 | cte_joined_outer_mat vs join_where | 92.9268 | 122.0927 | 1.3139 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 92.9268 | 91.2534 | 0.9820 | false |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_joined_outer_mat vs cte_joined_outer | 2.2459 | 98.4898 | 43.8540 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 2.2459 | 2.5284 | 1.1258 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 5.6248 | 64.9106 | 11.5401 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 5.6248 | 6.9368 | 1.2333 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 126.7825 | 122.0927 | 0.9630 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 126.7825 | 91.2534 | 0.7198 | false | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 1.5382 | 1 | `52aacf1c5e1c9d61` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 6.2159 | 1 | `f8e2569b4e380e0c` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | missing | | | | |
| default | 2 | 100% | missing | | | | |
| default | 4 | 1% | 6.1913 | 4 | `1a4292d5be1a7f26` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 6.2755 | 4 | `cf06a7ac60e0fe55` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | missing | | | | |
| default | 8 | 100% | missing | | | | |
| default | 16 | 1% | 7.3896 | 16 | `cf76bd74b3099cea` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 7.6066 | 16 | `48aa80f3875fa54d` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 5.5235 | 1 | `350c70c48e950ede` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 6.0351 | 1 | `72df45eb6f58c2ea` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | missing | | | | |
| mat | 2 | 100% | missing | | | | |
| mat | 4 | 1% | 6.1713 | 4 | `5f5d138969f1d64b` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 6.8657 | 4 | `09d7eb607ec8b133` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | missing | | | | |
| mat | 8 | 100% | missing | | | | |
| mat | 16 | 1% | 7.0122 | 16 | `5bf78a399b1db022` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 7.7867 | 16 | `629f1d8d839a1d67` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 1.5234 | 1 | `692b6f2284779825` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 6.1468 | 1 | `5a73dafada25ea4b` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | missing | | | | |
| notmat | 2 | 100% | missing | | | | |
| notmat | 4 | 1% | 3.5399 | 4 | `7360ee3ca3216a2e` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 21.3547 | 4 | `ca7dc31e4f39227a` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | missing | | | | |
| notmat | 8 | 100% | missing | | | | |
| notmat | 16 | 1% | 14.5390 | 16 | `24701f898050bab5` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 74.0736 | 16 | `b5967e5b6a765227` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 3.5909 | false |
| notmat | 1 | 1% | 0.9904 | false |
| mat | 1 | 100% | 0.9709 | false |
| notmat | 1 | 100% | 0.9889 | false |
| mat | 4 | 1% | 0.9968 | false |
| notmat | 4 | 1% | 0.5717 | false |
| mat | 4 | 100% | 1.0941 | false |
| notmat | 4 | 100% | 3.4029 | false |
| mat | 16 | 1% | 0.9489 | false |
| notmat | 16 | 1% | 1.9675 | false |
| mat | 16 | 100% | 1.0237 | false |
| notmat | 16 | 100% | 9.7380 | false |

## Physical arm

- `charity_pkey`: `CREATE UNIQUE INDEX charity_pkey ON selectjoin_domain.charity USING lsm (charity_id ASC)`
- `donation_person_idx`: `CREATE INDEX donation_person_idx ON selectjoin_domain.donation USING lsm (person_id ASC)`
- `donation_pkey`: `CREATE UNIQUE INDEX donation_pkey ON selectjoin_domain.donation USING lsm (donation_id ASC)`
- `donation_rank_idx`: `CREATE INDEX donation_rank_idx ON selectjoin_domain.donation USING lsm (rank ASC)`
- `person_pkey`: `CREATE UNIQUE INDEX person_pkey ON selectjoin_domain.person USING lsm (person_id ASC)`

- index arm `yb-range-rank`; placement `colocated`; node budget `2 CPU 3g`; client budget `1 CPU 2g`

| table | tablets | hash key columns | colocated | tablegroup | colocation id |
|---|---|---|---|---|---|
| selectjoin_domain.charity | 1 | 0 | true | 16388 | 1425957926 |
| selectjoin_domain.person | 1 | 0 | true | 16388 | 3204409570 |
| selectjoin_domain.donation | 1 | 0 | true | 16388 | 2457828944 |
| selectjoin_domain.donation_pkey | 0 | 0 | false | unknown | unknown |
| selectjoin_domain.donation_rank_idx | 1 | 0 | true | 16388 | 4256055230 |
| selectjoin_domain.donation_person_idx | 1 | 0 | true | 16388 | 2642029028 |

Physical-arm gate (default join must use the rank index at its three most selective thresholds):

- threshold 20: used
- threshold 2000: used

Tablet leaders per relation (placement evidence; the runner read these through yb-admin):

| relation | leader nodes |
|---|---|
| charity | yb-n2:9100 |
| person | yb-n2:9100 |
| donation | yb-n2:9100 |
| donation_rank_idx | yb-n2:9100 |
| donation_person_idx | yb-n2:9100 |

Effective sharding flag (from the engine): `ysql_num_shards_per_tserver=2`. A relation with fewer tablets than this setting was created without an explicit split; the per-table `num_tablets` above is the authority.

## Throttling (cgroup `cpu.stat` deltas)

- client: nr_periods 247, nr_throttled 81, throttled_usec 3498020 | quota 100000 100000
, memory 2147483648

- yb-n1: nr_periods 494, nr_throttled 5, throttled_usec 643974 | quota 200000 100000, memory 3221225472
- yb-n2: nr_periods 486, nr_throttled 10, throttled_usec 1569323 | quota 200000 100000, memory 3221225472
- yb-n3: nr_periods 480, nr_throttled 0, throttled_usec 0 | quota 200000 100000, memory 3221225472

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
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 118.5327 | 126.7825 | 132.2492 | 5 | 7.9 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 2.0020 | 2.2459 | 2.4351 | 5 | 445.3 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 5.1758 | 5.6248 | 6.0149 | 5 | 177.8 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 114.0115 | 122.0927 | 133.0484 | 5 | 8.2 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 96.6961 | 98.4898 | 105.3204 | 5 | 10.2 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 62.4458 | 64.9106 | 66.8538 | 5 | 15.4 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 88.0254 | 91.2534 | 92.9997 | 5 | 11.0 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 2.3888 | 2.5284 | 2.7401 | 5 | 395.5 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 6.5672 | 6.9368 | 44.7661 | 5 | 144.2 |
| family=bag;stmt=join_where;threshold=180000 | valid | 86.7005 | 92.9268 | 96.3484 | 5 | 10.8 |
| family=bag;stmt=join_where;threshold=20 | valid | 1.4228 | 1.6140 | 1.8063 | 5 | 619.6 |
| family=bag;stmt=join_where;threshold=2000 | valid | 5.4872 | 5.9581 | 6.0852 | 5 | 167.8 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 7.1659 | 7.3896 | 7.6872 | 5 | 135.3 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 7.1134 | 7.6066 | 8.7696 | 5 | 131.5 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 1.4128 | 1.5382 | 1.7175 | 5 | 650.1 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 6.1270 | 6.2159 | 6.4888 | 5 | 160.9 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 5.9071 | 6.1913 | 6.4313 | 5 | 161.5 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 6.0260 | 6.2755 | 9.5401 | 5 | 159.4 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 6.6372 | 7.0122 | 7.4339 | 5 | 142.6 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 7.7154 | 7.7867 | 7.8257 | 5 | 128.4 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 5.4839 | 5.5235 | 5.7259 | 5 | 181.0 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 5.6838 | 6.0351 | 7.0115 | 5 | 165.7 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 5.9893 | 6.1713 | 6.3268 | 5 | 162.0 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 6.8029 | 6.8657 | 7.5070 | 5 | 145.7 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 13.9999 | 14.5390 | 15.2206 | 5 | 68.8 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 66.7862 | 74.0736 | 80.4662 | 5 | 13.5 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 1.5124 | 1.5234 | 1.6115 | 5 | 656.4 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 5.9267 | 6.1468 | 6.3032 | 5 | 162.7 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 3.4462 | 3.5399 | 3.9552 | 5 | 282.5 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 20.4998 | 21.3547 | 22.9467 | 5 | 46.8 |

## Matched blended workload

- phase: `isolated`; target offered rate: 10 updates/s
- offered 0, committed 0, failed 0 over 0.000 s
- offered rate 0.000/s, completed rate 0.000/s
- per key: no keys touched
- ledger: 0 keys reconciled, 0 mismatch

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42593273`, 5 interleaved rounds
  1. family=reuse;policy=notmat;refs=16;consumer=1%
  2. family=bag;stmt=join_where;threshold=2000
  3. family=bag;stmt=cte_joined_outer;threshold=20
  4. family=reuse;policy=mat;refs=4;consumer=1%
  5. family=bag;stmt=join_where;threshold=20
  6. family=reuse;policy=default;refs=1;consumer=100%
  7. family=reuse;policy=default;refs=16;consumer=100%
  8. family=bag;stmt=cte_joined_outer;threshold=2000
  9. family=reuse;policy=mat;refs=16;consumer=1%
  10. family=reuse;policy=notmat;refs=4;consumer=100%
  11. family=reuse;policy=notmat;refs=4;consumer=1%
  12. family=reuse;policy=default;refs=4;consumer=1%
  13. family=reuse;policy=mat;refs=1;consumer=100%
  14. family=reuse;policy=mat;refs=1;consumer=1%
  15. family=reuse;policy=default;refs=4;consumer=100%
  16. family=bag;stmt=cte_joined_outer_mat;threshold=20
  17. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  18. family=reuse;policy=mat;refs=16;consumer=100%
  19. family=bag;stmt=cte_joined_outer;threshold=180000
  20. family=reuse;policy=default;refs=1;consumer=1%
  21. family=reuse;policy=default;refs=16;consumer=1%
  22. family=reuse;policy=notmat;refs=1;consumer=100%
  23. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  24. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  25. family=reuse;policy=notmat;refs=16;consumer=100%
  26. family=reuse;policy=notmat;refs=1;consumer=1%
  27. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  28. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  29. family=bag;stmt=join_where;threshold=180000
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
| family=bag;stmt=cte_joined_outer;threshold=180000 | j05_cte_joined_outer | valid | `e70a619cd1280d2c` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20 | j05_cte_joined_outer | valid | `0cde81a8c10bfe47` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=2000 | j05_cte_joined_outer | valid | `abd677c6ec45c985` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | j05m_cte_joined_outer_mat | valid | `6eeae5043ea31372` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | j05m_cte_joined_outer_mat | valid | `0864f5a134898925` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | j05m_cte_joined_outer_mat | valid | `b90cfb40cdfe9428` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | j05n_cte_joined_outer_notmat | valid | `2f086a8d6e750f7c` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | j05n_cte_joined_outer_notmat | valid | `00e2c44604f36147` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | j05n_cte_joined_outer_notmat | valid | `f6fd79f069cd8c1f` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=join_where;threshold=180000 | j01_join_where | valid | `c3fc96c733566cac` |  |
| family=bag;stmt=join_where;threshold=20 | j01_join_where | valid | `68503deedc3f7b4d` |  |
| family=bag;stmt=join_where;threshold=2000 | j01_join_where | valid | `2d7e8eeec22680da` |  |
| family=reuse;policy=default;refs=16;consumer=1% | r_default_n16_c1 | valid | `cf76bd74b3099cea` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=16;consumer=100% | r_default_n16_c100 | valid | `48aa80f3875fa54d` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=1;consumer=1% | r_default_n1_c1 | valid | `52aacf1c5e1c9d61` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=1;consumer=100% | r_default_n1_c100 | valid | `f8e2569b4e380e0c` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=4;consumer=1% | r_default_n4_c1 | valid | `1a4292d5be1a7f26` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=100% | r_default_n4_c100 | valid | `cf06a7ac60e0fe55` | policy=default; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=1% | r_mat_n16_c1 | valid | `5bf78a399b1db022` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=100% | r_mat_n16_c100 | valid | `629f1d8d839a1d67` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=1% | r_mat_n1_c1 | valid | `350c70c48e950ede` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=100% | r_mat_n1_c100 | valid | `72df45eb6f58c2ea` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=1% | r_mat_n4_c1 | valid | `5f5d138969f1d64b` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=100% | r_mat_n4_c100 | valid | `09d7eb607ec8b133` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=notmat;refs=16;consumer=1% | r_notmat_n16_c1 | valid | `24701f898050bab5` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `b5967e5b6a765227` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `692b6f2284779825` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `5a73dafada25ea4b` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `7360ee3ca3216a2e` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `ca7dc31e4f39227a` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, DIST, DEBUG, FORMAT TEXT` supported on yugabyte
- one host: this run used node budget 2 CPU 3g and client budget 1 CPU 2g, so budget comparisons are separate classes here, not independent hardware
- isolated phase: no concurrent writers in this run; the blended read/write phase is a separate labelled run
- cluster nodes are containers on the same host, so these are not network or independent-hardware measurements; container colocation is not data-colocation evidence
- readers round-robin across 3 endpoints; the writer's endpoint policy is recorded separately
- physical rank-index arm yb-range-rank; an engine whose planner ignores the index is a different physical design, not a faster one
- YugabyteDB DIST/DEBUG storage counters (rows and index rows scanned) are parsed per cell; RPC-level counters are not aggregated
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
