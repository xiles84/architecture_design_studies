# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 30 cells, 30 valid, engine `yugabyte`, topology `cluster3`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-v3b-coloc-i2-20261005`
- run tag: `run/study-01/select-join-v3b-coloc-i2-20261005`
- environment: `host-zenbook-ux5406sa`
- repo commit: `7243fee074bfb16fc5ad8b6afc1a2d7549080935`
- engine version: `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit`
- database engine image: `docker.io/yugabytedb/yugabyte:2025.2.6.0-b111` @ `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`
- benchmark client image: `localhost/charitybench-select:select-join-v3b-coloc-i2-20261005` @ `f3f1e10c44cbcb4c7b5bc56ed95a4f071682b6f7d33b5e0a7b688c4271eeef88`
- inputs digest: `e05845d60bf9fc195b051f88301d65547fd0682fee93bba80fd97790cdcad153`
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
| 20 (0.01%) | 1.7407 | missing | missing | missing | missing | missing | 2.4730 | 109.5407 | 2.3242 | missing |
| 2000 (1%) | 7.1472 | missing | missing | missing | missing | missing | 7.0161 | 112.8755 | 6.9640 | missing |
| 180000 (90%) | 134.4102 | missing | missing | missing | missing | missing | 138.5225 | 134.5914 | 101.1228 | missing |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | cte_joined_outer vs join_where | 1.7407 | 2.4730 | 1.4206 | false |
| 20 | cte_joined_outer_mat vs join_where | 1.7407 | 109.5407 | 62.9277 | false |
| 20 | cte_joined_outer_notmat vs join_where | 1.7407 | 2.3242 | 1.3352 | false |
| 2000 | cte_joined_outer vs join_where | 7.1472 | 7.0161 | 0.9817 | false |
| 2000 | cte_joined_outer_mat vs join_where | 7.1472 | 112.8755 | 15.7929 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 7.1472 | 6.9640 | 0.9744 | false |
| 180000 | cte_joined_outer vs join_where | 134.4102 | 138.5225 | 1.0306 | false |
| 180000 | cte_joined_outer_mat vs join_where | 134.4102 | 134.5914 | 1.0013 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 134.4102 | 101.1228 | 0.7523 | false |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_joined_outer_mat vs cte_joined_outer | 2.4730 | 109.5407 | 44.2954 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 2.4730 | 2.3242 | 0.9398 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 7.0161 | 112.8755 | 16.0880 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 7.0161 | 6.9640 | 0.9926 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 138.5225 | 134.5914 | 0.9716 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 138.5225 | 101.1228 | 0.7300 | false | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 1.6259 | 1 | `c673b15f6c7f6249` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 6.6151 | 1 | `f78d6b55f8083108` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | missing | | | | |
| default | 2 | 100% | missing | | | | |
| default | 4 | 1% | 6.9971 | 4 | `1c47f0f191091ab6` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 6.2222 | 4 | `6be8559a51431b3f` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | missing | | | | |
| default | 8 | 100% | missing | | | | |
| default | 16 | 1% | 7.7729 | 16 | `b358d8ab35c8d912` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 8.9044 | 16 | `3d405d6865e24287` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 6.2895 | 1 | `31d56c927858d7d4` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 5.9706 | 1 | `90720934deca6059` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | missing | | | | |
| mat | 2 | 100% | missing | | | | |
| mat | 4 | 1% | 7.6848 | 4 | `503d9de3f492696f` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 6.9561 | 4 | `28206afb4cad9ff2` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | missing | | | | |
| mat | 8 | 100% | missing | | | | |
| mat | 16 | 1% | 7.5222 | 16 | `a96a5ba8f9fa31c6` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 8.2524 | 16 | `644cb18ea5a6abb8` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 1.3403 | 1 | `09826ed17aeead8a` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 6.5568 | 1 | `4319e2534fda0fa4` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | missing | | | | |
| notmat | 2 | 100% | missing | | | | |
| notmat | 4 | 1% | 4.6905 | 4 | `0083fff4d4f3de23` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 21.7213 | 4 | `fa8842dab3efefbf` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | missing | | | | |
| notmat | 8 | 100% | missing | | | | |
| notmat | 16 | 1% | 17.2355 | 16 | `eaa25f9bc2720bc7` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 80.5005 | 16 | `2cf20ed59b1d0df0` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 3.8682 | false |
| notmat | 1 | 1% | 0.8243 | false |
| mat | 1 | 100% | 0.9026 | false |
| notmat | 1 | 100% | 0.9912 | false |
| mat | 4 | 1% | 1.0983 | false |
| notmat | 4 | 1% | 0.6703 | false |
| mat | 4 | 100% | 1.1179 | false |
| notmat | 4 | 100% | 3.4909 | false |
| mat | 16 | 1% | 0.9677 | false |
| notmat | 16 | 1% | 2.2174 | false |
| mat | 16 | 100% | 0.9268 | false |
| notmat | 16 | 100% | 9.0405 | false |

## Physical arm

- `charity_pkey`: `CREATE UNIQUE INDEX charity_pkey ON selectjoin_domain.charity USING lsm (charity_id ASC)`
- `donation_person_idx`: `CREATE INDEX donation_person_idx ON selectjoin_domain.donation USING lsm (person_id ASC)`
- `donation_pkey`: `CREATE UNIQUE INDEX donation_pkey ON selectjoin_domain.donation USING lsm (donation_id ASC)`
- `donation_rank_idx`: `CREATE INDEX donation_rank_idx ON selectjoin_domain.donation USING lsm (rank ASC)`
- `person_pkey`: `CREATE UNIQUE INDEX person_pkey ON selectjoin_domain.person USING lsm (person_id ASC)`

- index arm `yb-range-rank`; placement `colocated`; node budget `2 CPU 3g`; client budget `1 CPU 2g`

| table | tablets | hash key columns | colocated | tablegroup | colocation id |
|---|---|---|---|---|---|
| selectjoin_domain.charity | 1 | 0 | true | 16388 | 4192460358 |
| selectjoin_domain.person | 1 | 0 | true | 16388 | 2029726056 |
| selectjoin_domain.donation | 1 | 0 | true | 16388 | 680993452 |
| selectjoin_domain.donation_pkey | 0 | 0 | false | unknown | unknown |
| selectjoin_domain.donation_rank_idx | 1 | 0 | true | 16388 | 3563693350 |
| selectjoin_domain.donation_person_idx | 1 | 0 | true | 16388 | 1530534549 |

Physical-arm gate (default join must use the rank index at its three most selective thresholds):

- threshold 20: used
- threshold 2000: used

## Throttling (cgroup `cpu.stat` deltas)

- client: nr_periods 280, nr_throttled 90, throttled_usec 4635477 | quota 100000 100000
, memory 2147483648

- yb-n1: nr_periods 535, nr_throttled 6, throttled_usec 987249 | quota 200000 100000, memory 3221225472
- yb-n2: nr_periods 526, nr_throttled 6, throttled_usec 915820 | quota 200000 100000, memory 3221225472
- yb-n3: nr_periods 518, nr_throttled 5, throttled_usec 297311 | quota 200000 100000, memory 3221225472

## Access path per cell

_Engine-reported from each cell's own plan: whether the rank index or a sequential scan was used, rows scanned, PostgreSQL buffers and filter removals, and reader retries of transient conflicts. A cell that scanned the whole table at a selective threshold is a different physical regime, not a faster query._

| axis | status | index scan | seq scan | storage rows scanned | index rows scanned | buffers hit/read | rows removed | retries |
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
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 130.0434 | 138.5225 | 143.8841 | 5 | 7.2 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 2.3255 | 2.4730 | 2.7500 | 5 | 404.4 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 6.7649 | 7.0161 | 7.0644 | 5 | 142.5 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 129.9945 | 134.5914 | 147.3548 | 5 | 7.4 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 108.1771 | 109.5407 | 113.1244 | 5 | 9.1 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 106.6468 | 112.8755 | 116.0584 | 5 | 8.9 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 96.0426 | 101.1228 | 117.5845 | 5 | 9.9 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 2.2015 | 2.3242 | 2.4075 | 5 | 430.3 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 6.8402 | 6.9640 | 7.3303 | 5 | 143.6 |
| family=bag;stmt=join_where;threshold=180000 | valid | 132.7407 | 134.4102 | 138.4132 | 5 | 7.4 |
| family=bag;stmt=join_where;threshold=20 | valid | 1.5658 | 1.7407 | 2.2043 | 5 | 574.5 |
| family=bag;stmt=join_where;threshold=2000 | valid | 6.9687 | 7.1472 | 7.5268 | 5 | 139.9 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 7.4558 | 7.7729 | 8.6122 | 5 | 128.7 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 8.6680 | 8.9044 | 9.2009 | 5 | 112.3 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 1.6188 | 1.6259 | 1.7431 | 5 | 615.0 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 6.3070 | 6.6151 | 6.8575 | 5 | 151.2 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 6.5241 | 6.9971 | 7.2116 | 5 | 142.9 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 6.0225 | 6.2222 | 6.8689 | 5 | 160.7 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 7.1898 | 7.5222 | 7.6343 | 5 | 132.9 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 7.8799 | 8.2524 | 8.5751 | 5 | 121.2 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 6.0273 | 6.2895 | 6.5190 | 5 | 159.0 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 5.6796 | 5.9706 | 6.2658 | 5 | 167.5 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 7.2055 | 7.6848 | 7.7937 | 5 | 130.1 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 6.7246 | 6.9561 | 7.2985 | 5 | 143.8 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 16.3872 | 17.2355 | 17.3957 | 5 | 58.0 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 73.4356 | 80.5005 | 111.6951 | 5 | 12.4 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 1.2943 | 1.3403 | 1.3646 | 5 | 746.1 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 6.0985 | 6.5568 | 26.7453 | 5 | 152.5 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 4.5428 | 4.6905 | 4.7952 | 5 | 213.2 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 20.8195 | 21.7213 | 24.0272 | 5 | 46.0 |

## Matched blended workload

- phase: `isolated`; target offered rate: 10 updates/s
- offered 0, committed 0, failed 0 over 0.000 s
- offered rate 0.000/s, completed rate 0.000/s
- per key: no keys touched
- ledger: 0 keys reconciled, 0 mismatch

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42057602`, 5 interleaved rounds
  1. family=reuse;policy=default;refs=16;consumer=1%
  2. family=reuse;policy=default;refs=4;consumer=1%
  3. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  4. family=reuse;policy=mat;refs=4;consumer=1%
  5. family=bag;stmt=cte_joined_outer;threshold=20
  6. family=reuse;policy=notmat;refs=1;consumer=100%
  7. family=bag;stmt=join_where;threshold=180000
  8. family=reuse;policy=default;refs=16;consumer=100%
  9. family=reuse;policy=mat;refs=1;consumer=100%
  10. family=reuse;policy=mat;refs=16;consumer=1%
  11. family=bag;stmt=cte_joined_outer;threshold=180000
  12. family=reuse;policy=notmat;refs=16;consumer=100%
  13. family=reuse;policy=notmat;refs=16;consumer=1%
  14. family=reuse;policy=notmat;refs=4;consumer=1%
  15. family=reuse;policy=notmat;refs=4;consumer=100%
  16. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  17. family=bag;stmt=join_where;threshold=2000
  18. family=reuse;policy=default;refs=4;consumer=100%
  19. family=reuse;policy=default;refs=1;consumer=1%
  20. family=reuse;policy=default;refs=1;consumer=100%
  21. family=reuse;policy=notmat;refs=1;consumer=1%
  22. family=reuse;policy=mat;refs=4;consumer=100%
  23. family=reuse;policy=mat;refs=1;consumer=1%
  24. family=bag;stmt=join_where;threshold=20
  25. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  26. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  27. family=reuse;policy=mat;refs=16;consumer=100%
  28. family=bag;stmt=cte_joined_outer_mat;threshold=20
  29. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  30. family=bag;stmt=cte_joined_outer;threshold=2000

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
| family=bag;stmt=cte_joined_outer;threshold=180000 | j05_cte_joined_outer | valid | `6c8938dac90749d5` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20 | j05_cte_joined_outer | valid | `46a7321bf207d4a0` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=2000 | j05_cte_joined_outer | valid | `ccdf2574fc108b62` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | j05m_cte_joined_outer_mat | valid | `89261d009b0854cb` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | j05m_cte_joined_outer_mat | valid | `78faf2492f0f184f` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | j05m_cte_joined_outer_mat | valid | `c8a2c573754b6216` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | j05n_cte_joined_outer_notmat | valid | `569172c49c73e659` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | j05n_cte_joined_outer_notmat | valid | `cfaaa8cccf85edf0` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | j05n_cte_joined_outer_notmat | valid | `d6f739fed21e5ae4` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=join_where;threshold=180000 | j01_join_where | valid | `eee749cf8ab0b3fa` |  |
| family=bag;stmt=join_where;threshold=20 | j01_join_where | valid | `0967fb8eccdb649e` |  |
| family=bag;stmt=join_where;threshold=2000 | j01_join_where | valid | `bc16579aeeee313c` |  |
| family=reuse;policy=default;refs=16;consumer=1% | r_default_n16_c1 | valid | `b358d8ab35c8d912` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=16;consumer=100% | r_default_n16_c100 | valid | `3d405d6865e24287` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=1;consumer=1% | r_default_n1_c1 | valid | `c673b15f6c7f6249` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=1;consumer=100% | r_default_n1_c100 | valid | `f78d6b55f8083108` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=4;consumer=1% | r_default_n4_c1 | valid | `1c47f0f191091ab6` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=100% | r_default_n4_c100 | valid | `6be8559a51431b3f` | policy=default; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=1% | r_mat_n16_c1 | valid | `a96a5ba8f9fa31c6` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=100% | r_mat_n16_c100 | valid | `644cb18ea5a6abb8` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=1% | r_mat_n1_c1 | valid | `31d56c927858d7d4` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=100% | r_mat_n1_c100 | valid | `90720934deca6059` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=1% | r_mat_n4_c1 | valid | `503d9de3f492696f` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=100% | r_mat_n4_c100 | valid | `28206afb4cad9ff2` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=notmat;refs=16;consumer=1% | r_notmat_n16_c1 | valid | `eaa25f9bc2720bc7` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `2cf20ed59b1d0df0` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `09826ed17aeead8a` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `4319e2534fda0fa4` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `0083fff4d4f3de23` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `fa8842dab3efefbf` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, DIST, DEBUG, FORMAT TEXT` supported on yugabyte
- one host: this run used node budget 2 CPU 3g and client budget 1 CPU 2g, so budget comparisons are separate classes here, not independent hardware
- isolated phase: no concurrent writers in this run; the blended read/write phase is a separate labelled run
- cluster nodes are containers on the same host, so these are not network or independent-hardware measurements; container colocation is not data-colocation evidence
- readers round-robin across 3 endpoints; the writer's endpoint policy is recorded separately
- physical rank-index arm yb-range-rank; an engine whose planner ignores the index is a different physical design, not a faster one
- YugabyteDB DIST/DEBUG storage counters (rows and index rows scanned) are parsed per cell; RPC-level counters are not aggregated
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
