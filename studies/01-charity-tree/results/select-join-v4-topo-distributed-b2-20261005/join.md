# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 30 cells, 30 valid, engine `yugabyte`, topology `cluster3`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-v4-topo-distributed-b2-20261005`
- run tag: `run/study-01/select-join-v4-topo-distributed-b2-20261005`
- environment: `host-zenbook-ux5406sa`
- repo commit: `ba821e7b686ac3adb1cda56dd45f15c8ae699ef8`
- engine version: `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit`
- database engine image: `docker.io/yugabytedb/yugabyte:2025.2.6.0-b111` @ `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`
- benchmark client image: `localhost/charitybench-select:select-join-v4-topo-distributed-b2-20261005` @ `b6ea3c007cf459112eaa8e9e85381a7dc25a7710ff07648fde3048f747df6ab0`
- inputs digest: `233446b75e88c47c953204b767c1ef54ff47c91967edbe07448dceab512448a7`
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
| yb-n2 | 169 | 2 |
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
| 20 (0.01%) | 3.1432 | missing | missing | missing | missing | missing | 2.9347 | 74.8211 | 3.2771 | missing |
| 2000 (1%) | 8.3521 | missing | missing | missing | missing | missing | 8.3636 | 69.8051 | 8.3742 | missing |
| 180000 (90%) | 83.0955 | missing | missing | missing | missing | missing | 86.2752 | 111.8808 | 93.9067 | missing |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | cte_joined_outer vs join_where | 3.1432 | 2.9347 | 0.9337 | false |
| 20 | cte_joined_outer_mat vs join_where | 3.1432 | 74.8211 | 23.8041 | false |
| 20 | cte_joined_outer_notmat vs join_where | 3.1432 | 3.2771 | 1.0426 | false |
| 2000 | cte_joined_outer vs join_where | 8.3521 | 8.3636 | 1.0014 | false |
| 2000 | cte_joined_outer_mat vs join_where | 8.3521 | 69.8051 | 8.3578 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 8.3521 | 8.3742 | 1.0027 | false |
| 180000 | cte_joined_outer vs join_where | 83.0955 | 86.2752 | 1.0383 | false |
| 180000 | cte_joined_outer_mat vs join_where | 83.0955 | 111.8808 | 1.3464 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 83.0955 | 93.9067 | 1.1301 | false |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_joined_outer_mat vs cte_joined_outer | 2.9347 | 74.8211 | 25.4956 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 2.9347 | 3.2771 | 1.1167 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 8.3636 | 69.8051 | 8.3463 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 8.3636 | 8.3742 | 1.0013 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 86.2752 | 111.8808 | 1.2968 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 86.2752 | 93.9067 | 1.0885 | false | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 5.1566 | 1 | `910e828ce2d83a53` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 7.8833 | 1 | `0c36bfe442b498d2` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | missing | | | | |
| default | 2 | 100% | missing | | | | |
| default | 4 | 1% | 8.0836 | 4 | `4986f9dd128a75a1` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 8.0861 | 4 | `35e9442af87c9aa6` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | missing | | | | |
| default | 8 | 100% | missing | | | | |
| default | 16 | 1% | 8.5467 | 16 | `ef93d2dd142d454c` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 9.4018 | 16 | `1959a4dd61cdc254` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 7.6900 | 1 | `9e89513516726d63` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 7.6363 | 1 | `91815e6b6b077dab` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | missing | | | | |
| mat | 2 | 100% | missing | | | | |
| mat | 4 | 1% | 7.8430 | 4 | `1153e976b0202ecf` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 7.9377 | 4 | `54d7bc4a59043b78` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | missing | | | | |
| mat | 8 | 100% | missing | | | | |
| mat | 16 | 1% | 10.2841 | 16 | `98e0b277ea57ed30` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 9.7893 | 16 | `908cba4c119d5e56` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 5.6639 | 1 | `81f6169eaafee907` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 7.0357 | 1 | `8bda7581ea9d456a` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | missing | | | | |
| notmat | 2 | 100% | missing | | | | |
| notmat | 4 | 1% | 17.9784 | 4 | `ac22380ea6f4b720` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 27.2062 | 4 | `c0fb2dfe2892b01c` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | missing | | | | |
| notmat | 8 | 100% | missing | | | | |
| notmat | 16 | 1% | 88.6595 | 16 | `dc839babb00df3d9` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 108.7505 | 16 | `cedcd8c18284d938` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 1.4913 | false |
| notmat | 1 | 1% | 1.0984 | false |
| mat | 1 | 100% | 0.9687 | false |
| notmat | 1 | 100% | 0.8925 | false |
| mat | 4 | 1% | 0.9702 | false |
| notmat | 4 | 1% | 2.2241 | false |
| mat | 4 | 100% | 0.9816 | false |
| notmat | 4 | 100% | 3.3646 | false |
| mat | 16 | 1% | 1.2033 | false |
| notmat | 16 | 1% | 10.3736 | false |
| mat | 16 | 100% | 1.0412 | false |
| notmat | 16 | 100% | 11.5670 | false |

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

- client: nr_periods 325, nr_throttled 144, throttled_usec 5823602 | quota 100000 100000
, memory 2147483648

- yb-n1: nr_periods 675, nr_throttled 9, throttled_usec 1040129 | quota 200000 100000, memory 3221225472
- yb-n2: nr_periods 666, nr_throttled 5, throttled_usec 879690 | quota 200000 100000, memory 3221225472
- yb-n3: nr_periods 661, nr_throttled 3, throttled_usec 138115 | quota 200000 100000, memory 3221225472

## Access path per cell

_Engine-reported from each cell's own plan. The two scan flags are **whole-plan** booleans: they are true when any node of the plan is an index scan / sequential scan, including small parent-table scans, so they are not the rank access path. `storage rows scanned` is the engine's own count (YugabyteDB); PostgreSQL does not print one and it is shown as n/a. A cell that scanned the whole child table at a selective threshold is a different physical regime, not a faster query._

| axis | status | whole-plan index scan | whole-plan seq scan | storage rows scanned | index rows scanned | buffers hit/read | rows removed | retries |
|---|---|---|---|---|---|---|---|---|
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | false | true | 202050 | 0 | 0/0 | 0 | 1 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | true | true | 110 | 20 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | true | true | 6050 | 2000 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | false | true | 202050 | 0 | 0/0 | 20000 | 1 |
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
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 77.9997 | 86.2752 | 125.2234 | 5 | 11.6 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 2.7660 | 2.9347 | 3.9273 | 5 | 340.8 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 7.5324 | 8.3636 | 9.9204 | 5 | 119.6 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 104.3220 | 111.8808 | 156.8241 | 5 | 8.9 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 71.0936 | 74.8211 | 93.9668 | 5 | 13.4 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 67.5056 | 69.8051 | 86.6488 | 5 | 14.3 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 72.2954 | 93.9067 | 99.3655 | 5 | 10.6 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 2.6293 | 3.2771 | 23.5418 | 5 | 305.2 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 7.8998 | 8.3742 | 9.2526 | 5 | 119.4 |
| family=bag;stmt=join_where;threshold=180000 | valid | 75.9689 | 83.0955 | 96.7101 | 5 | 12.0 |
| family=bag;stmt=join_where;threshold=20 | valid | 3.0060 | 3.1432 | 3.8461 | 5 | 318.1 |
| family=bag;stmt=join_where;threshold=2000 | valid | 7.9851 | 8.3521 | 9.5082 | 5 | 119.7 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 7.8187 | 8.5467 | 9.3202 | 5 | 117.0 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 9.0251 | 9.4018 | 14.2277 | 5 | 106.4 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 4.9117 | 5.1566 | 6.5719 | 5 | 193.9 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 6.8704 | 7.8833 | 9.0700 | 5 | 126.9 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 7.3254 | 8.0836 | 9.5091 | 5 | 123.7 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 7.0277 | 8.0861 | 8.6178 | 5 | 123.7 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 7.8866 | 10.2841 | 13.6141 | 5 | 97.2 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 8.6973 | 9.7893 | 10.9669 | 5 | 102.2 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 7.1716 | 7.6900 | 9.8050 | 5 | 130.0 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 7.5982 | 7.6363 | 9.4050 | 5 | 131.0 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 7.4639 | 7.8430 | 9.6935 | 5 | 127.5 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 7.1890 | 7.9377 | 8.9543 | 5 | 126.0 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 75.6920 | 88.6595 | 124.9491 | 5 | 11.3 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 83.9541 | 108.7505 | 141.9182 | 5 | 9.2 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 5.2667 | 5.6639 | 5.9210 | 5 | 176.6 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 6.6175 | 7.0357 | 9.2698 | 5 | 142.1 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 17.6992 | 17.9784 | 21.1313 | 5 | 55.6 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 24.8691 | 27.2062 | 31.3583 | 5 | 36.8 |

## Matched blended workload

- phase: `blended`; target offered rate: 10 updates/s
- offered 162, committed 162, failed 0 over 16.163 s
- offered rate 10.023/s, completed rate 10.023/s
- writer operations per endpoint: [54 54 54]
- reader operations overlapping a committed write: 61 of 212
- reader retries of transient conflicts: 2
- per key: 162 keys, min 1, max 1 per key
- ledger: 162 keys reconciled, 0 mismatch
- update latency: avg 9.572 ms, min 2.582, max 220.012 over 162

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42216824`, 5 interleaved rounds
  1. family=reuse;policy=notmat;refs=1;consumer=1%
  2. family=reuse;policy=notmat;refs=4;consumer=100%
  3. family=reuse;policy=default;refs=16;consumer=1%
  4. family=reuse;policy=mat;refs=4;consumer=100%
  5. family=bag;stmt=cte_joined_outer;threshold=2000
  6. family=reuse;policy=default;refs=4;consumer=100%
  7. family=bag;stmt=join_where;threshold=180000
  8. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  9. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  10. family=reuse;policy=default;refs=4;consumer=1%
  11. family=reuse;policy=default;refs=1;consumer=1%
  12. family=reuse;policy=notmat;refs=16;consumer=1%
  13. family=reuse;policy=default;refs=16;consumer=100%
  14. family=reuse;policy=notmat;refs=16;consumer=100%
  15. family=reuse;policy=default;refs=1;consumer=100%
  16. family=bag;stmt=join_where;threshold=20
  17. family=bag;stmt=cte_joined_outer_mat;threshold=20
  18. family=reuse;policy=mat;refs=16;consumer=100%
  19. family=reuse;policy=mat;refs=1;consumer=100%
  20. family=reuse;policy=mat;refs=4;consumer=1%
  21. family=bag;stmt=cte_joined_outer;threshold=180000
  22. family=reuse;policy=mat;refs=1;consumer=1%
  23. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  24. family=bag;stmt=cte_joined_outer;threshold=20
  25. family=reuse;policy=notmat;refs=1;consumer=100%
  26. family=reuse;policy=mat;refs=16;consumer=1%
  27. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  28. family=bag;stmt=join_where;threshold=2000
  29. family=reuse;policy=notmat;refs=4;consumer=1%
  30. family=bag;stmt=cte_joined_outer_notmat;threshold=2000

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
| family=bag;stmt=cte_joined_outer;threshold=180000 | j05_cte_joined_outer | valid | `b7aef1f040fea0da` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20 | j05_cte_joined_outer | valid | `78663a5626a5f416` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=2000 | j05_cte_joined_outer | valid | `c80b3fc708f97183` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | j05m_cte_joined_outer_mat | valid | `83f3f7fc14e95a76` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | j05m_cte_joined_outer_mat | valid | `d4d17a39d0aa4d25` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | j05m_cte_joined_outer_mat | valid | `89dae894beaf6b49` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | j05n_cte_joined_outer_notmat | valid | `2422fbdcb3acb3a1` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | j05n_cte_joined_outer_notmat | valid | `164958a6a30a4566` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | j05n_cte_joined_outer_notmat | valid | `8be11cf0fe3a210d` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=join_where;threshold=180000 | j01_join_where | valid | `b6102eaf39b9b85f` |  |
| family=bag;stmt=join_where;threshold=20 | j01_join_where | valid | `d9acc0b4ffbcbb54` |  |
| family=bag;stmt=join_where;threshold=2000 | j01_join_where | valid | `d0e29d17ef9d93b8` |  |
| family=reuse;policy=default;refs=16;consumer=1% | r_default_n16_c1 | valid | `ef93d2dd142d454c` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=16;consumer=100% | r_default_n16_c100 | valid | `1959a4dd61cdc254` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=1;consumer=1% | r_default_n1_c1 | valid | `910e828ce2d83a53` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=1;consumer=100% | r_default_n1_c100 | valid | `0c36bfe442b498d2` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=4;consumer=1% | r_default_n4_c1 | valid | `4986f9dd128a75a1` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=100% | r_default_n4_c100 | valid | `35e9442af87c9aa6` | policy=default; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=1% | r_mat_n16_c1 | valid | `98e0b277ea57ed30` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=100% | r_mat_n16_c100 | valid | `908cba4c119d5e56` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=1% | r_mat_n1_c1 | valid | `9e89513516726d63` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=100% | r_mat_n1_c100 | valid | `91815e6b6b077dab` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=1% | r_mat_n4_c1 | valid | `1153e976b0202ecf` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=100% | r_mat_n4_c100 | valid | `54d7bc4a59043b78` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=notmat;refs=16;consumer=1% | r_notmat_n16_c1 | valid | `dc839babb00df3d9` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `cedcd8c18284d938` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `81f6169eaafee907` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `8bda7581ea9d456a` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `ac22380ea6f4b720` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `c0fb2dfe2892b01c` | policy=not_materialized; observed_cte_scan=false |

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
