# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 30 cells, 30 valid, engine `yugabyte`, topology `cluster3`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-v4-topo-colocated-i1-20261005`
- run tag: `run/study-01/select-join-v4-topo-colocated-i1-20261005`
- environment: `host-zenbook-ux5406sa`
- repo commit: `ba821e7b686ac3adb1cda56dd45f15c8ae699ef8`
- engine version: `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit`
- database engine image: `docker.io/yugabytedb/yugabyte:2025.2.6.0-b111` @ `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`
- benchmark client image: `localhost/charitybench-select:select-join-v4-topo-colocated-i1-20261005` @ `b6ea3c007cf459112eaa8e9e85381a7dc25a7710ff07648fde3048f747df6ab0`
- inputs digest: `771f30fddb4e5880544603c9cc85e1424602ea23efe11e298b93cb35f0f11a38`
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
| 20 (0.01%) | 2.2624 | missing | missing | missing | missing | missing | 2.2971 | 66.3033 | 2.2220 | missing |
| 2000 (1%) | 6.2830 | missing | missing | missing | missing | missing | 6.0258 | 101.2475 | 6.6017 | missing |
| 180000 (90%) | 89.3831 | missing | missing | missing | missing | missing | 92.7433 | 118.3854 | 121.3268 | missing |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | cte_joined_outer vs join_where | 2.2624 | 2.2971 | 1.0153 | false |
| 20 | cte_joined_outer_mat vs join_where | 2.2624 | 66.3033 | 29.3071 | false |
| 20 | cte_joined_outer_notmat vs join_where | 2.2624 | 2.2220 | 0.9821 | false |
| 2000 | cte_joined_outer vs join_where | 6.2830 | 6.0258 | 0.9591 | false |
| 2000 | cte_joined_outer_mat vs join_where | 6.2830 | 101.2475 | 16.1144 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 6.2830 | 6.6017 | 1.0507 | false |
| 180000 | cte_joined_outer vs join_where | 89.3831 | 92.7433 | 1.0376 | false |
| 180000 | cte_joined_outer_mat vs join_where | 89.3831 | 118.3854 | 1.3245 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 89.3831 | 121.3268 | 1.3574 | false |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_joined_outer_mat vs cte_joined_outer | 2.2971 | 66.3033 | 28.8644 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 2.2971 | 2.2220 | 0.9673 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 6.0258 | 101.2475 | 16.8022 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 6.0258 | 6.6017 | 1.0956 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 92.7433 | 118.3854 | 1.2765 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 92.7433 | 121.3268 | 1.3082 | false | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 1.5258 | 1 | `1dd5eb2091ae388a` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 5.8739 | 1 | `150d58137e4fa691` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | missing | | | | |
| default | 2 | 100% | missing | | | | |
| default | 4 | 1% | 6.7233 | 4 | `3e6673c79c24524c` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 5.6241 | 4 | `bb93424b83cae805` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | missing | | | | |
| default | 8 | 100% | missing | | | | |
| default | 16 | 1% | 7.0203 | 16 | `7a11cf78d1b3c98b` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 7.3026 | 16 | `ca636720e5f5ef4f` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 5.6936 | 1 | `efe2be3a43aab680` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 6.0484 | 1 | `16606d5855393b57` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | missing | | | | |
| mat | 2 | 100% | missing | | | | |
| mat | 4 | 1% | 5.5072 | 4 | `4d8ea5f76138dc9c` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 5.7661 | 4 | `2b2e58cee632df25` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | missing | | | | |
| mat | 8 | 100% | missing | | | | |
| mat | 16 | 1% | 6.7819 | 16 | `347ae3c8e380bc3f` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 6.9184 | 16 | `b6c812fefc246896` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 1.5917 | 1 | `4533aa75e4d07f31` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 6.6999 | 1 | `9b8322013d285b8c` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | missing | | | | |
| notmat | 2 | 100% | missing | | | | |
| notmat | 4 | 1% | 3.6095 | 4 | `ee639e8f9ff0003d` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 20.4266 | 4 | `153ee79fea58b731` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | missing | | | | |
| notmat | 8 | 100% | missing | | | | |
| notmat | 16 | 1% | 14.2082 | 16 | `0f3fc32c5a134ddf` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 73.6152 | 16 | `ec24bb765d3a039c` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 3.7316 | false |
| notmat | 1 | 1% | 1.0432 | false |
| mat | 1 | 100% | 1.0297 | false |
| notmat | 1 | 100% | 1.1406 | false |
| mat | 4 | 1% | 0.8191 | false |
| notmat | 4 | 1% | 0.5369 | false |
| mat | 4 | 100% | 1.0252 | false |
| notmat | 4 | 100% | 3.6319 | false |
| mat | 16 | 1% | 0.9660 | false |
| notmat | 16 | 1% | 2.0239 | false |
| mat | 16 | 100% | 0.9474 | false |
| notmat | 16 | 100% | 10.0806 | false |

## Physical arm

- `charity_pkey`: `CREATE UNIQUE INDEX charity_pkey ON selectjoin_domain.charity USING lsm (charity_id ASC)`
- `donation_person_idx`: `CREATE INDEX donation_person_idx ON selectjoin_domain.donation USING lsm (person_id ASC)`
- `donation_pkey`: `CREATE UNIQUE INDEX donation_pkey ON selectjoin_domain.donation USING lsm (donation_id ASC)`
- `donation_rank_idx`: `CREATE INDEX donation_rank_idx ON selectjoin_domain.donation USING lsm (rank ASC)`
- `person_pkey`: `CREATE UNIQUE INDEX person_pkey ON selectjoin_domain.person USING lsm (person_id ASC)`

- index arm `yb-range-rank`; placement `colocated`; node budget `2 CPU 3g`; client budget `1 CPU 2g`

| table | tablets | hash key columns | colocated | tablegroup | colocation id |
|---|---|---|---|---|---|
| selectjoin_domain.charity | 1 | 0 | true | 16388 | 3653221518 |
| selectjoin_domain.person | 1 | 0 | true | 16388 | 23141175 |
| selectjoin_domain.donation | 1 | 0 | true | 16388 | 1862887405 |
| selectjoin_domain.donation_pkey | 0 | 0 | false | unknown | unknown |
| selectjoin_domain.donation_rank_idx | 1 | 0 | true | 16388 | 3225599719 |
| selectjoin_domain.donation_person_idx | 1 | 0 | true | 16388 | 3611304706 |

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

- client: nr_periods 251, nr_throttled 81, throttled_usec 3918876 | quota 100000 100000
, memory 2147483648

- yb-n1: nr_periods 495, nr_throttled 4, throttled_usec 544740 | quota 200000 100000, memory 3221225472
- yb-n2: nr_periods 484, nr_throttled 8, throttled_usec 1135954 | quota 200000 100000, memory 3221225472
- yb-n3: nr_periods 483, nr_throttled 0, throttled_usec 0 | quota 200000 100000, memory 3221225472

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
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 87.9696 | 92.7433 | 95.4299 | 5 | 10.8 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 2.0900 | 2.2971 | 2.8674 | 5 | 435.3 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 5.9100 | 6.0258 | 6.1310 | 5 | 166.0 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 113.8119 | 118.3854 | 123.8183 | 5 | 8.4 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 61.4691 | 66.3033 | 67.6194 | 5 | 15.1 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 94.1515 | 101.2475 | 104.8619 | 5 | 9.9 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 116.1943 | 121.3268 | 127.7459 | 5 | 8.2 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 2.0839 | 2.2220 | 2.4501 | 5 | 450.0 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 6.4699 | 6.6017 | 7.5651 | 5 | 151.5 |
| family=bag;stmt=join_where;threshold=180000 | valid | 88.0999 | 89.3831 | 100.7990 | 5 | 11.2 |
| family=bag;stmt=join_where;threshold=20 | valid | 2.1492 | 2.2624 | 2.4942 | 5 | 442.0 |
| family=bag;stmt=join_where;threshold=2000 | valid | 6.0962 | 6.2830 | 6.5825 | 5 | 159.2 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 6.7626 | 7.0203 | 7.1392 | 5 | 142.4 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 6.9292 | 7.3026 | 7.3360 | 5 | 136.9 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 1.4888 | 1.5258 | 1.5904 | 5 | 655.4 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 5.7107 | 5.8739 | 6.1514 | 5 | 170.2 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 6.4433 | 6.7233 | 6.8602 | 5 | 148.7 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 5.3520 | 5.6241 | 6.6601 | 5 | 177.8 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 6.5598 | 6.7819 | 7.7427 | 5 | 147.5 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 6.8510 | 6.9184 | 7.5059 | 5 | 144.5 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 5.4741 | 5.6936 | 5.7879 | 5 | 175.6 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 5.9200 | 6.0484 | 6.2757 | 5 | 165.3 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 5.4925 | 5.5072 | 5.6702 | 5 | 181.6 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 5.6786 | 5.7661 | 6.9847 | 5 | 173.4 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 14.0387 | 14.2082 | 15.1238 | 5 | 70.4 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 72.2974 | 73.6152 | 80.6499 | 5 | 13.6 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 1.5473 | 1.5917 | 1.6280 | 5 | 628.2 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 6.5302 | 6.6999 | 6.8952 | 5 | 149.3 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 3.5800 | 3.6095 | 3.8488 | 5 | 277.0 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 20.2956 | 20.4266 | 22.2695 | 5 | 49.0 |

## Matched blended workload

- phase: `isolated`; target offered rate: 10 updates/s
- offered 0, committed 0, failed 0 over 0.000 s
- offered rate 0.000/s, completed rate 0.000/s
- per key: no keys touched
- ledger: 0 keys reconciled, 0 mismatch

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42167567`, 5 interleaved rounds
  1. family=reuse;policy=notmat;refs=4;consumer=100%
  2. family=reuse;policy=default;refs=4;consumer=100%
  3. family=bag;stmt=join_where;threshold=20
  4. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  5. family=reuse;policy=notmat;refs=4;consumer=1%
  6. family=reuse;policy=mat;refs=1;consumer=100%
  7. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  8. family=reuse;policy=default;refs=16;consumer=1%
  9. family=reuse;policy=default;refs=1;consumer=1%
  10. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  11. family=reuse;policy=mat;refs=16;consumer=100%
  12. family=reuse;policy=mat;refs=16;consumer=1%
  13. family=reuse;policy=notmat;refs=16;consumer=1%
  14. family=bag;stmt=cte_joined_outer;threshold=2000
  15. family=bag;stmt=cte_joined_outer;threshold=20
  16. family=bag;stmt=join_where;threshold=2000
  17. family=reuse;policy=mat;refs=4;consumer=1%
  18. family=reuse;policy=mat;refs=1;consumer=1%
  19. family=reuse;policy=default;refs=1;consumer=100%
  20. family=bag;stmt=join_where;threshold=180000
  21. family=reuse;policy=notmat;refs=16;consumer=100%
  22. family=reuse;policy=default;refs=4;consumer=1%
  23. family=bag;stmt=cte_joined_outer_mat;threshold=20
  24. family=reuse;policy=notmat;refs=1;consumer=1%
  25. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  26. family=bag;stmt=cte_joined_outer;threshold=180000
  27. family=reuse;policy=notmat;refs=1;consumer=100%
  28. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  29. family=reuse;policy=mat;refs=4;consumer=100%
  30. family=reuse;policy=default;refs=16;consumer=100%

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
| family=bag;stmt=cte_joined_outer;threshold=180000 | j05_cte_joined_outer | valid | `15cc86614cfe0d76` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20 | j05_cte_joined_outer | valid | `e6bcba4083ede1af` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=2000 | j05_cte_joined_outer | valid | `fd262744b5216c4e` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | j05m_cte_joined_outer_mat | valid | `219386173d3f3b5d` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | j05m_cte_joined_outer_mat | valid | `1bb87225f709ac51` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | j05m_cte_joined_outer_mat | valid | `bf953139d1692f6b` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | j05n_cte_joined_outer_notmat | valid | `f3d8866fadc1e0b4` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | j05n_cte_joined_outer_notmat | valid | `3c7d2ca369ed57ef` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | j05n_cte_joined_outer_notmat | valid | `b786197f55f9acc1` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=join_where;threshold=180000 | j01_join_where | valid | `f2051edbd8227bf0` |  |
| family=bag;stmt=join_where;threshold=20 | j01_join_where | valid | `9d0655e3858ba862` |  |
| family=bag;stmt=join_where;threshold=2000 | j01_join_where | valid | `ecbc590a87caafa0` |  |
| family=reuse;policy=default;refs=16;consumer=1% | r_default_n16_c1 | valid | `7a11cf78d1b3c98b` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=16;consumer=100% | r_default_n16_c100 | valid | `ca636720e5f5ef4f` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=1;consumer=1% | r_default_n1_c1 | valid | `1dd5eb2091ae388a` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=1;consumer=100% | r_default_n1_c100 | valid | `150d58137e4fa691` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=4;consumer=1% | r_default_n4_c1 | valid | `3e6673c79c24524c` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=100% | r_default_n4_c100 | valid | `bb93424b83cae805` | policy=default; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=1% | r_mat_n16_c1 | valid | `347ae3c8e380bc3f` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=100% | r_mat_n16_c100 | valid | `b6c812fefc246896` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=1% | r_mat_n1_c1 | valid | `efe2be3a43aab680` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=100% | r_mat_n1_c100 | valid | `16606d5855393b57` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=1% | r_mat_n4_c1 | valid | `4d8ea5f76138dc9c` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=100% | r_mat_n4_c100 | valid | `2b2e58cee632df25` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=notmat;refs=16;consumer=1% | r_notmat_n16_c1 | valid | `0f3fc32c5a134ddf` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `ec24bb765d3a039c` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `4533aa75e4d07f31` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `9b8322013d285b8c` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `ee639e8f9ff0003d` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `153ee79fea58b731` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, DIST, DEBUG, FORMAT TEXT` supported on yugabyte
- one host: this run used node budget 2 CPU 3g and client budget 1 CPU 2g, so budget comparisons are separate classes here, not independent hardware
- isolated phase: no concurrent writers in this run; the blended read/write phase is a separate labelled run
- cluster nodes are containers on the same host, so these are not network or independent-hardware measurements; container colocation is not data-colocation evidence
- readers round-robin across 3 endpoints; the writer's endpoint policy is recorded separately
- physical rank-index arm yb-range-rank; an engine whose planner ignores the index is a different physical design, not a faster one
- YugabyteDB DIST/DEBUG storage counters (rows and index rows scanned) are parsed per cell; RPC-level counters are not aggregated
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
