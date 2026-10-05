# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 30 cells, 30 valid, engine `yugabyte`, topology `cluster3`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-v4-topo-tablet-i1-20261005`
- run tag: `run/study-01/select-join-v4-topo-tablet-i1-20261005`
- environment: `host-zenbook-ux5406sa`
- repo commit: `ba821e7b686ac3adb1cda56dd45f15c8ae699ef8`
- engine version: `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit`
- database engine image: `docker.io/yugabytedb/yugabyte:2025.2.6.0-b111` @ `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`
- benchmark client image: `localhost/charitybench-select:select-join-v4-topo-tablet-i1-20261005` @ `b6ea3c007cf459112eaa8e9e85381a7dc25a7710ff07648fde3048f747df6ab0`
- inputs digest: `6122526384aa46083db9172dc68a117e37915dd37a9807571f85973be1ae4a3a`
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
| 20 (0.01%) | 1.8712 | missing | missing | missing | missing | missing | 2.7232 | 120.5125 | 2.7147 | missing |
| 2000 (1%) | 9.2925 | missing | missing | missing | missing | missing | 10.6796 | 125.4706 | 10.5949 | missing |
| 180000 (90%) | 142.8339 | missing | missing | missing | missing | missing | 143.9169 | 134.0807 | 141.0170 | missing |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | cte_joined_outer vs join_where | 1.8712 | 2.7232 | 1.4553 | false |
| 20 | cte_joined_outer_mat vs join_where | 1.8712 | 120.5125 | 64.4035 | false |
| 20 | cte_joined_outer_notmat vs join_where | 1.8712 | 2.7147 | 1.4508 | false |
| 2000 | cte_joined_outer vs join_where | 9.2925 | 10.6796 | 1.1493 | false |
| 2000 | cte_joined_outer_mat vs join_where | 9.2925 | 125.4706 | 13.5023 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 9.2925 | 10.5949 | 1.1402 | false |
| 180000 | cte_joined_outer vs join_where | 142.8339 | 143.9169 | 1.0076 | false |
| 180000 | cte_joined_outer_mat vs join_where | 142.8339 | 134.0807 | 0.9387 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 142.8339 | 141.0170 | 0.9873 | false |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_joined_outer_mat vs cte_joined_outer | 2.7232 | 120.5125 | 44.2539 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 2.7232 | 2.7147 | 0.9969 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 10.6796 | 125.4706 | 11.7487 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 10.6796 | 10.5949 | 0.9921 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 143.9169 | 134.0807 | 0.9317 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 143.9169 | 141.0170 | 0.9798 | false | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 8.4291 | 1 | `e9a23d3c3a2e453d` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 8.0878 | 1 | `19ecd6bba1f02bb3` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | missing | | | | |
| default | 2 | 100% | missing | | | | |
| default | 4 | 1% | 9.2660 | 4 | `f6232ffaa37a8bd1` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 9.6243 | 4 | `c169a7b934b607f5` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | missing | | | | |
| default | 8 | 100% | missing | | | | |
| default | 16 | 1% | 9.9788 | 16 | `9af5585d4f5d5586` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 9.9252 | 16 | `c62ad49ae769bada` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 9.1057 | 1 | `ffeea94550c863df` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 8.3138 | 1 | `ed43f53142efacad` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | missing | | | | |
| mat | 2 | 100% | missing | | | | |
| mat | 4 | 1% | 9.1150 | 4 | `cdb760de646c3e21` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 9.2572 | 4 | `afa27595d0ae2094` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | missing | | | | |
| mat | 8 | 100% | missing | | | | |
| mat | 16 | 1% | 9.9818 | 16 | `1c44e2e2bf80d083` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 11.3022 | 16 | `5ec367e4051d2df5` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 7.6431 | 1 | `8f048c3df89c4062` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 9.3544 | 1 | `ee69d3eefa0c17f1` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | missing | | | | |
| notmat | 2 | 100% | missing | | | | |
| notmat | 4 | 1% | 27.2331 | 4 | `25b722a28262b827` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 28.1329 | 4 | `95e41c8f4eae601b` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | missing | | | | |
| notmat | 8 | 100% | missing | | | | |
| notmat | 16 | 1% | 75.3324 | 16 | `6392bdc15e5d043c` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 89.6101 | 16 | `650923bc1c19fa37` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 1.0803 | false |
| notmat | 1 | 1% | 0.9067 | false |
| mat | 1 | 100% | 1.0279 | false |
| notmat | 1 | 100% | 1.1566 | false |
| mat | 4 | 1% | 0.9837 | false |
| notmat | 4 | 1% | 2.9390 | false |
| mat | 4 | 100% | 0.9619 | false |
| notmat | 4 | 100% | 2.9231 | false |
| mat | 16 | 1% | 1.0003 | false |
| notmat | 16 | 1% | 7.5493 | false |
| mat | 16 | 100% | 1.1387 | false |
| notmat | 16 | 100% | 9.0285 | false |

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

- client: nr_periods 327, nr_throttled 108, throttled_usec 4874043 | quota 100000 100000
, memory 2147483648

- yb-n1: nr_periods 602, nr_throttled 4, throttled_usec 575890 | quota 200000 100000, memory 3221225472
- yb-n2: nr_periods 593, nr_throttled 6, throttled_usec 764364 | quota 200000 100000, memory 3221225472
- yb-n3: nr_periods 590, nr_throttled 5, throttled_usec 49224 | quota 200000 100000, memory 3221225472

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
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 141.9285 | 143.9169 | 153.8211 | 5 | 6.9 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 2.5954 | 2.7232 | 2.9025 | 5 | 367.2 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 10.2361 | 10.6796 | 11.5831 | 5 | 93.6 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 131.4552 | 134.0807 | 142.5914 | 5 | 7.5 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 118.4166 | 120.5125 | 127.7312 | 5 | 8.3 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 118.9716 | 125.4706 | 130.8167 | 5 | 8.0 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 137.0123 | 141.0170 | 144.9920 | 5 | 7.1 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 2.5409 | 2.7147 | 2.8424 | 5 | 368.4 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 9.8819 | 10.5949 | 11.2910 | 5 | 94.4 |
| family=bag;stmt=join_where;threshold=180000 | valid | 136.4086 | 142.8339 | 148.3691 | 5 | 7.0 |
| family=bag;stmt=join_where;threshold=20 | valid | 1.7910 | 1.8712 | 2.0284 | 5 | 534.4 |
| family=bag;stmt=join_where;threshold=2000 | valid | 8.5744 | 9.2925 | 10.8383 | 5 | 107.6 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 9.5356 | 9.9788 | 10.1030 | 5 | 100.2 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 9.5315 | 9.9252 | 10.3401 | 5 | 100.8 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 8.0662 | 8.4291 | 8.4588 | 5 | 118.6 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 7.9456 | 8.0878 | 8.6136 | 5 | 123.6 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 9.0860 | 9.2660 | 10.8826 | 5 | 107.9 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 8.9801 | 9.6243 | 10.5647 | 5 | 103.9 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 9.4233 | 9.9818 | 10.0624 | 5 | 100.2 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 11.1635 | 11.3022 | 11.7996 | 5 | 88.5 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 8.8287 | 9.1057 | 9.3712 | 5 | 109.8 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 7.8676 | 8.3138 | 8.8049 | 5 | 120.3 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 8.8249 | 9.1150 | 9.7789 | 5 | 109.7 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 8.6335 | 9.2572 | 9.6919 | 5 | 108.0 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 66.4765 | 75.3324 | 111.3972 | 5 | 13.3 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 86.2329 | 89.6101 | 100.0275 | 5 | 11.2 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 7.2482 | 7.6431 | 8.0963 | 5 | 130.8 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 8.6885 | 9.3544 | 9.6838 | 5 | 106.9 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 27.0514 | 27.2331 | 30.8482 | 5 | 36.7 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 27.3616 | 28.1329 | 32.7901 | 5 | 35.5 |

## Matched blended workload

- phase: `isolated`; target offered rate: 10 updates/s
- offered 0, committed 0, failed 0 over 0.000 s
- offered rate 0.000/s, completed rate 0.000/s
- per key: no keys touched
- ledger: 0 keys reconciled, 0 mismatch

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42902477`, 5 interleaved rounds
  1. family=reuse;policy=notmat;refs=4;consumer=1%
  2. family=reuse;policy=mat;refs=1;consumer=1%
  3. family=bag;stmt=join_where;threshold=20
  4. family=reuse;policy=notmat;refs=1;consumer=100%
  5. family=reuse;policy=default;refs=4;consumer=100%
  6. family=reuse;policy=default;refs=16;consumer=100%
  7. family=reuse;policy=mat;refs=4;consumer=100%
  8. family=bag;stmt=cte_joined_outer_mat;threshold=20
  9. family=reuse;policy=default;refs=16;consumer=1%
  10. family=bag;stmt=cte_joined_outer;threshold=20
  11. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  12. family=reuse;policy=mat;refs=16;consumer=1%
  13. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  14. family=bag;stmt=join_where;threshold=180000
  15. family=reuse;policy=notmat;refs=16;consumer=1%
  16. family=reuse;policy=default;refs=1;consumer=1%
  17. family=bag;stmt=cte_joined_outer;threshold=2000
  18. family=reuse;policy=default;refs=1;consumer=100%
  19. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  20. family=reuse;policy=notmat;refs=1;consumer=1%
  21. family=reuse;policy=notmat;refs=4;consumer=100%
  22. family=reuse;policy=default;refs=4;consumer=1%
  23. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  24. family=bag;stmt=join_where;threshold=2000
  25. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  26. family=reuse;policy=mat;refs=16;consumer=100%
  27. family=reuse;policy=mat;refs=1;consumer=100%
  28. family=reuse;policy=notmat;refs=16;consumer=100%
  29. family=bag;stmt=cte_joined_outer;threshold=180000
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
| family=bag;stmt=cte_joined_outer;threshold=180000 | j05_cte_joined_outer | valid | `41283d395bc5537f` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20 | j05_cte_joined_outer | valid | `4fbceda2aa1cba81` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=2000 | j05_cte_joined_outer | valid | `3aec9709475828f8` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | j05m_cte_joined_outer_mat | valid | `c9c25ee6dc54dbac` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | j05m_cte_joined_outer_mat | valid | `b8a51dcec1d79380` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | j05m_cte_joined_outer_mat | valid | `a336155972f95ef8` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | j05n_cte_joined_outer_notmat | valid | `c1919aee4b4a1db8` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | j05n_cte_joined_outer_notmat | valid | `c7d3c2af0f846792` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | j05n_cte_joined_outer_notmat | valid | `85109f2fe750fb40` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=join_where;threshold=180000 | j01_join_where | valid | `255b3cfef9e387b3` |  |
| family=bag;stmt=join_where;threshold=20 | j01_join_where | valid | `e22a07ff76f8cf34` |  |
| family=bag;stmt=join_where;threshold=2000 | j01_join_where | valid | `6f86f2ab617bacc7` |  |
| family=reuse;policy=default;refs=16;consumer=1% | r_default_n16_c1 | valid | `9af5585d4f5d5586` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=16;consumer=100% | r_default_n16_c100 | valid | `c62ad49ae769bada` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=1;consumer=1% | r_default_n1_c1 | valid | `e9a23d3c3a2e453d` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=1;consumer=100% | r_default_n1_c100 | valid | `19ecd6bba1f02bb3` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=4;consumer=1% | r_default_n4_c1 | valid | `f6232ffaa37a8bd1` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=100% | r_default_n4_c100 | valid | `c169a7b934b607f5` | policy=default; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=1% | r_mat_n16_c1 | valid | `1c44e2e2bf80d083` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=100% | r_mat_n16_c100 | valid | `5ec367e4051d2df5` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=1% | r_mat_n1_c1 | valid | `ffeea94550c863df` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=100% | r_mat_n1_c100 | valid | `ed43f53142efacad` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=1% | r_mat_n4_c1 | valid | `cdb760de646c3e21` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=100% | r_mat_n4_c100 | valid | `afa27595d0ae2094` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=notmat;refs=16;consumer=1% | r_notmat_n16_c1 | valid | `6392bdc15e5d043c` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `650923bc1c19fa37` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `8f048c3df89c4062` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `ee69d3eefa0c17f1` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `25b722a28262b827` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `95e41c8f4eae601b` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, DIST, DEBUG, FORMAT TEXT` supported on yugabyte
- one host: this run used node budget 2 CPU 3g and client budget 1 CPU 2g, so budget comparisons are separate classes here, not independent hardware
- isolated phase: no concurrent writers in this run; the blended read/write phase is a separate labelled run
- cluster nodes are containers on the same host, so these are not network or independent-hardware measurements; container colocation is not data-colocation evidence
- readers round-robin across 3 endpoints; the writer's endpoint policy is recorded separately
- physical rank-index arm yb-range-rank; an engine whose planner ignores the index is a different physical design, not a faster one
- YugabyteDB DIST/DEBUG storage counters (rows and index rows scanned) are parsed per cell; RPC-level counters are not aggregated
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
