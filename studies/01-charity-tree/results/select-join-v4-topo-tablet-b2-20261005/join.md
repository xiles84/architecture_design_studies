# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 30 cells, 30 valid, engine `yugabyte`, topology `cluster3`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-v4-topo-tablet-b2-20261005`
- run tag: `run/study-01/select-join-v4-topo-tablet-b2-20261005`
- environment: `host-zenbook-ux5406sa`
- repo commit: `ba821e7b686ac3adb1cda56dd45f15c8ae699ef8`
- engine version: `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit`
- database engine image: `docker.io/yugabytedb/yugabyte:2025.2.6.0-b111` @ `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`
- benchmark client image: `localhost/charitybench-select:select-join-v4-topo-tablet-b2-20261005` @ `b6ea3c007cf459112eaa8e9e85381a7dc25a7710ff07648fde3048f747df6ab0`
- inputs digest: `9bbb7a64c3dc758819f23782c3d16987a59c490cd657008ceaf5faf1bb6fc16a`
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
| 20 (0.01%) | 3.2479 | missing | missing | missing | missing | missing | 2.8920 | 133.5939 | 3.2901 | missing |
| 2000 (1%) | 11.3076 | missing | missing | missing | missing | missing | 13.1497 | 136.8571 | 9.7356 | missing |
| 180000 (90%) | 121.9118 | missing | missing | missing | missing | missing | 154.4661 | 150.2236 | 147.0753 | missing |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | cte_joined_outer vs join_where | 3.2479 | 2.8920 | 0.8904 | false |
| 20 | cte_joined_outer_mat vs join_where | 3.2479 | 133.5939 | 41.1328 | false |
| 20 | cte_joined_outer_notmat vs join_where | 3.2479 | 3.2901 | 1.0130 | false |
| 2000 | cte_joined_outer vs join_where | 11.3076 | 13.1497 | 1.1629 | false |
| 2000 | cte_joined_outer_mat vs join_where | 11.3076 | 136.8571 | 12.1031 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 11.3076 | 9.7356 | 0.8610 | false |
| 180000 | cte_joined_outer vs join_where | 121.9118 | 154.4661 | 1.2670 | false |
| 180000 | cte_joined_outer_mat vs join_where | 121.9118 | 150.2236 | 1.2322 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 121.9118 | 147.0753 | 1.2064 | false |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_joined_outer_mat vs cte_joined_outer | 2.8920 | 133.5939 | 46.1947 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 2.8920 | 3.2901 | 1.1377 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 13.1497 | 136.8571 | 10.4077 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 13.1497 | 9.7356 | 0.7404 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 154.4661 | 150.2236 | 0.9725 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 154.4661 | 147.0753 | 0.9522 | false | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 7.5779 | 1 | `1b8dd5b6cd16dffb` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 10.1817 | 1 | `2338fecfd41cad4e` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | missing | | | | |
| default | 2 | 100% | missing | | | | |
| default | 4 | 1% | 10.5773 | 4 | `6e63ec170f7a25b5` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 9.8072 | 4 | `a19f6f98b8d18899` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | missing | | | | |
| default | 8 | 100% | missing | | | | |
| default | 16 | 1% | 11.6859 | 16 | `31ccae67bde8eec6` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 11.7102 | 16 | `5f310c14ecc38884` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 8.4837 | 1 | `35e9281ad8dd9a84` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 10.4883 | 1 | `ddc80dede531ab3d` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | missing | | | | |
| mat | 2 | 100% | missing | | | | |
| mat | 4 | 1% | 10.0566 | 4 | `1cd1e600713acf3b` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 9.7623 | 4 | `5eda43197552e87d` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | missing | | | | |
| mat | 8 | 100% | missing | | | | |
| mat | 16 | 1% | 10.4245 | 16 | `b611bdc84bbb70a5` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 11.0817 | 16 | `3c85d9d5dafa3d7e` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 7.3729 | 1 | `0ea9e0cf102bccfa` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 8.6449 | 1 | `2622c212b196407e` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | missing | | | | |
| notmat | 2 | 100% | missing | | | | |
| notmat | 4 | 1% | 28.4491 | 4 | `afcb8cb6eac2e685` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 35.7623 | 4 | `09938a3e542f1fec` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | missing | | | | |
| notmat | 8 | 100% | missing | | | | |
| notmat | 16 | 1% | 94.5046 | 16 | `1f9ccf2cd6562a3c` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 117.8989 | 16 | `2a6ee79ff23e5d0c` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 1.1195 | false |
| notmat | 1 | 1% | 0.9729 | false |
| mat | 1 | 100% | 1.0301 | false |
| notmat | 1 | 100% | 0.8491 | false |
| mat | 4 | 1% | 0.9508 | false |
| notmat | 4 | 1% | 2.6896 | false |
| mat | 4 | 100% | 0.9954 | false |
| notmat | 4 | 100% | 3.6465 | false |
| mat | 16 | 1% | 0.8921 | false |
| notmat | 16 | 1% | 8.0871 | false |
| mat | 16 | 100% | 0.9463 | false |
| notmat | 16 | 100% | 10.0680 | false |

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

- client: nr_periods 346, nr_throttled 125, throttled_usec 5125823 | quota 100000 100000
, memory 2147483648

- yb-n1: nr_periods 621, nr_throttled 6, throttled_usec 562946 | quota 200000 100000, memory 3221225472
- yb-n2: nr_periods 616, nr_throttled 5, throttled_usec 709435 | quota 200000 100000, memory 3221225472
- yb-n3: nr_periods 608, nr_throttled 10, throttled_usec 331576 | quota 200000 100000, memory 3221225472

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
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 147.3953 | 154.4661 | 161.5391 | 5 | 6.5 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 2.6547 | 2.8920 | 3.3223 | 5 | 345.8 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 11.0823 | 13.1497 | 13.9561 | 5 | 76.0 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 138.2932 | 150.2236 | 193.6682 | 5 | 6.7 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 126.8480 | 133.5939 | 154.2017 | 5 | 7.5 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 124.6845 | 136.8571 | 147.9125 | 5 | 7.3 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 142.8765 | 147.0753 | 197.7163 | 5 | 6.8 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 3.0956 | 3.2901 | 3.7807 | 5 | 303.9 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 9.4110 | 9.7356 | 10.5338 | 5 | 102.7 |
| family=bag;stmt=join_where;threshold=180000 | valid | 118.2202 | 121.9118 | 151.8108 | 5 | 8.2 |
| family=bag;stmt=join_where;threshold=20 | valid | 2.9349 | 3.2479 | 3.6803 | 5 | 307.9 |
| family=bag;stmt=join_where;threshold=2000 | valid | 10.3686 | 11.3076 | 13.8827 | 5 | 88.4 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 10.9534 | 11.6859 | 12.4928 | 5 | 85.6 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 10.5169 | 11.7102 | 12.2267 | 5 | 85.4 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 7.0636 | 7.5779 | 8.6465 | 5 | 132.0 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 9.6008 | 10.1817 | 12.6163 | 5 | 98.2 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 9.9885 | 10.5773 | 12.9281 | 5 | 94.5 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 9.4162 | 9.8072 | 13.1803 | 5 | 102.0 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 10.1798 | 10.4245 | 12.5126 | 5 | 95.9 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 10.5278 | 11.0817 | 14.6789 | 5 | 90.2 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 8.4173 | 8.4837 | 10.8491 | 5 | 117.9 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 9.9192 | 10.4883 | 11.2197 | 5 | 95.3 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 9.2881 | 10.0566 | 10.3174 | 5 | 99.4 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 9.2542 | 9.7623 | 10.7750 | 5 | 102.4 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 86.1086 | 94.5046 | 113.4176 | 5 | 10.6 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 85.4796 | 117.8989 | 137.5708 | 5 | 8.5 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 7.0581 | 7.3729 | 8.6979 | 5 | 135.6 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 7.9914 | 8.6449 | 9.8424 | 5 | 115.7 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 27.9794 | 28.4491 | 35.9857 | 5 | 35.2 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 33.2735 | 35.7623 | 45.8358 | 5 | 28.0 |

## Matched blended workload

- phase: `blended`; target offered rate: 10 updates/s
- offered 175, committed 175, failed 0 over 17.462 s
- offered rate 10.022/s, completed rate 10.022/s
- writer operations per endpoint: [59 58 58]
- reader operations overlapping a committed write: 76 of 210
- reader retries of transient conflicts: 0
- per key: 175 keys, min 1, max 1 per key
- ledger: 175 keys reconciled, 0 mismatch
- update latency: avg 6.649 ms, min 2.241, max 86.542 over 175

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42385475`, 5 interleaved rounds
  1. family=reuse;policy=notmat;refs=4;consumer=100%
  2. family=reuse;policy=notmat;refs=4;consumer=1%
  3. family=reuse;policy=mat;refs=1;consumer=1%
  4. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  5. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  6. family=reuse;policy=default;refs=1;consumer=1%
  7. family=bag;stmt=join_where;threshold=2000
  8. family=reuse;policy=notmat;refs=16;consumer=1%
  9. family=bag;stmt=join_where;threshold=180000
  10. family=reuse;policy=default;refs=16;consumer=1%
  11. family=reuse;policy=mat;refs=16;consumer=1%
  12. family=reuse;policy=notmat;refs=1;consumer=100%
  13. family=reuse;policy=default;refs=16;consumer=100%
  14. family=bag;stmt=cte_joined_outer;threshold=2000
  15. family=reuse;policy=mat;refs=4;consumer=1%
  16. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  17. family=bag;stmt=cte_joined_outer;threshold=20
  18. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  19. family=bag;stmt=cte_joined_outer;threshold=180000
  20. family=reuse;policy=mat;refs=1;consumer=100%
  21. family=reuse;policy=notmat;refs=16;consumer=100%
  22. family=bag;stmt=join_where;threshold=20
  23. family=reuse;policy=default;refs=4;consumer=1%
  24. family=reuse;policy=mat;refs=4;consumer=100%
  25. family=reuse;policy=default;refs=1;consumer=100%
  26. family=bag;stmt=cte_joined_outer_mat;threshold=20
  27. family=reuse;policy=notmat;refs=1;consumer=1%
  28. family=reuse;policy=default;refs=4;consumer=100%
  29. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  30. family=reuse;policy=mat;refs=16;consumer=100%

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
| family=bag;stmt=cte_joined_outer;threshold=180000 | j05_cte_joined_outer | valid | `b287293dea88a5ca` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20 | j05_cte_joined_outer | valid | `a9df78ca56f502ae` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=2000 | j05_cte_joined_outer | valid | `b6b348b0c4c946b1` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | j05m_cte_joined_outer_mat | valid | `2c89f4abcd5fab43` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | j05m_cte_joined_outer_mat | valid | `bec8e43fe5e73744` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | j05m_cte_joined_outer_mat | valid | `7d6c5d368c7d08c9` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | j05n_cte_joined_outer_notmat | valid | `eb77b21866bd3651` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | j05n_cte_joined_outer_notmat | valid | `eeabc9d1abc43991` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | j05n_cte_joined_outer_notmat | valid | `1000dcdbc4d5a60f` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=join_where;threshold=180000 | j01_join_where | valid | `49808a0d635858e9` |  |
| family=bag;stmt=join_where;threshold=20 | j01_join_where | valid | `df32f5c17fc4aea2` |  |
| family=bag;stmt=join_where;threshold=2000 | j01_join_where | valid | `01716465efc9d9ca` |  |
| family=reuse;policy=default;refs=16;consumer=1% | r_default_n16_c1 | valid | `31ccae67bde8eec6` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=16;consumer=100% | r_default_n16_c100 | valid | `5f310c14ecc38884` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=1;consumer=1% | r_default_n1_c1 | valid | `1b8dd5b6cd16dffb` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=1;consumer=100% | r_default_n1_c100 | valid | `2338fecfd41cad4e` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=4;consumer=1% | r_default_n4_c1 | valid | `6e63ec170f7a25b5` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=100% | r_default_n4_c100 | valid | `a19f6f98b8d18899` | policy=default; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=1% | r_mat_n16_c1 | valid | `b611bdc84bbb70a5` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=100% | r_mat_n16_c100 | valid | `3c85d9d5dafa3d7e` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=1% | r_mat_n1_c1 | valid | `35e9281ad8dd9a84` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=100% | r_mat_n1_c100 | valid | `ddc80dede531ab3d` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=1% | r_mat_n4_c1 | valid | `1cd1e600713acf3b` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=100% | r_mat_n4_c100 | valid | `5eda43197552e87d` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=notmat;refs=16;consumer=1% | r_notmat_n16_c1 | valid | `1f9ccf2cd6562a3c` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `2a6ee79ff23e5d0c` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `0ea9e0cf102bccfa` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `2622c212b196407e` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `afcb8cb6eac2e685` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `09938a3e542f1fec` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, DIST, DEBUG, FORMAT TEXT` supported on yugabyte
- one host: this run used node budget 2 CPU 3g and client budget 1 CPU 2g, so budget comparisons are separate classes here, not independent hardware
- one client process with one fixed maintenance writer at 10 updates/s; readers are sequential; the reported sample latency includes retry and backoff time
- cluster nodes are containers on the same host, so these are not network or independent-hardware measurements; container colocation is not data-colocation evidence
- readers round-robin across 3 endpoints; the writer's endpoint policy is recorded separately
- physical rank-index arm yb-range-rank; an engine whose planner ignores the index is a different physical design, not a faster one
- YugabyteDB DIST/DEBUG storage counters (rows and index rows scanned) are parsed per cell; RPC-level counters are not aggregated
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
