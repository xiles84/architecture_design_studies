# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 30 cells, 30 valid, engine `yugabyte`, topology `single`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-v3b-topo-equal-b3-20261005`
- run tag: `run/study-01/select-join-v3b-topo-equal-b3-20261005`
- environment: `host-zenbook-ux5406sa`
- repo commit: `7243fee074bfb16fc5ad8b6afc1a2d7549080935`
- engine version: `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit`
- database engine image: `docker.io/yugabytedb/yugabyte:2025.2.6.0-b111` @ `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`
- benchmark client image: `localhost/charitybench-select:select-join-v3b-topo-equal-b3-20261005` @ `f3f1e10c44cbcb4c7b5bc56ed95a4f071682b6f7d33b5e0a7b688c4271eeef88`
- inputs digest: `e785333b56256bce973f9296b59b43fd988f8b6871549fe60d3a5cf71a5900d2`
- dataset fingerprint: `fcf26e08f377ec51`
- client CPU quota: `100000 100000`, memory: `2147483648`
- setting `client_budget`: 1 CPU 2g
- setting `client_workers`: 1 sequential reader; 0 writers
- setting `connection_pool`: 2 per endpoint = reader/setup + reserve
- setting `customers`: 2000
- setting `dataset`: study01-person-donation-charity
- setting `distributed_rpc`: DIST/DEBUG storage counters parsed per cell (Storage Rows Scanned, Storage Index Rows Scanned); RPC-level counters not aggregated
- setting `endpoint_mode`: single endpoint
- setting `endpoints`: yb-single
- setting `fact_count`: 200000
- setting `index_arm`: yb-range-rank
- setting `index_set`: child rank index and parent-key index, primary keys; no per-spelling index
- setting `node_budget`: 6 CPU 9g
- setting `phase`: blended
- setting `placement`: noncolocated
- setting `prepared_plan_policy`: pgx simple protocol; no prepared statements
- setting `scale`: cte-domain
- setting `selectivity_mapping`: rank <= threshold; rank is a deterministic shuffle of 1..fact_count
- setting `statement_budget_ms`: 120000
- setting `statistics`: ANALYZE after deterministic load
- setting `thresholds`: [20 2000 180000]
- setting `topology`: single
- setting `trials_per_cell`: 5
- setting `whole_run_count`: 1; no replicated winner in this file
- setting `writer_keys`: 1000
- setting `writer_target_rate_per_sec`: 10

### Endpoint distribution

A cluster arm must rotate its reader across nodes; a single seed endpoint would not be balanced cluster access.

| endpoint | operations | connections at end |
|---|---|---|
| yb-single | 504 | 2 |

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
| 20 (0.01%) | 1.8267 | missing | missing | missing | missing | missing | 1.9828 | 87.9115 | 2.0601 | missing |
| 2000 (1%) | 9.6112 | missing | missing | missing | missing | missing | 9.7466 | 87.5445 | 11.2835 | missing |
| 180000 (90%) | 112.1092 | missing | missing | missing | missing | missing | 120.0761 | 108.7952 | 115.0683 | missing |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | cte_joined_outer vs join_where | 1.8267 | 1.9828 | 1.0854 | false |
| 20 | cte_joined_outer_mat vs join_where | 1.8267 | 87.9115 | 48.1255 | false |
| 20 | cte_joined_outer_notmat vs join_where | 1.8267 | 2.0601 | 1.1278 | false |
| 2000 | cte_joined_outer vs join_where | 9.6112 | 9.7466 | 1.0141 | false |
| 2000 | cte_joined_outer_mat vs join_where | 9.6112 | 87.5445 | 9.1086 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 9.6112 | 11.2835 | 1.1740 | false |
| 180000 | cte_joined_outer vs join_where | 112.1092 | 120.0761 | 1.0711 | false |
| 180000 | cte_joined_outer_mat vs join_where | 112.1092 | 108.7952 | 0.9704 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 112.1092 | 115.0683 | 1.0264 | false |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_joined_outer_mat vs cte_joined_outer | 1.9828 | 87.9115 | 44.3378 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 1.9828 | 2.0601 | 1.0390 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 9.7466 | 87.5445 | 8.9820 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 9.7466 | 11.2835 | 1.1577 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 120.0761 | 108.7952 | 0.9061 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 120.0761 | 115.0683 | 0.9583 | false | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 7.7286 | 1 | `47e0b35694db6da5` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 9.3265 | 1 | `038a7f7d48ffa5da` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | missing | | | | |
| default | 2 | 100% | missing | | | | |
| default | 4 | 1% | 9.4821 | 4 | `8f83aca0ae937564` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 9.2396 | 4 | `bae65d24dfbc1ad0` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | missing | | | | |
| default | 8 | 100% | missing | | | | |
| default | 16 | 1% | 9.8102 | 16 | `e33cdbe68fbf169a` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 10.2137 | 16 | `75a8c4dd8c0a4739` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 9.0351 | 1 | `29faa216d3dfa30d` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 9.6020 | 1 | `6020c01948d84f44` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | missing | | | | |
| mat | 2 | 100% | missing | | | | |
| mat | 4 | 1% | 8.7467 | 4 | `0a1e88a407f861da` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 10.0081 | 4 | `bbeaa1eb4016dbd3` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | missing | | | | |
| mat | 8 | 100% | missing | | | | |
| mat | 16 | 1% | 9.6120 | 16 | `cf1cb3928efbf66a` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 10.4970 | 16 | `6419e021c8506f26` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 6.9386 | 1 | `e5f9bc9fee9418a1` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 8.8536 | 1 | `1d77a6d535328dca` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | missing | | | | |
| notmat | 2 | 100% | missing | | | | |
| notmat | 4 | 1% | 22.8222 | 4 | `15aaa6e8e9632446` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 30.9524 | 4 | `2c051c650a265c67` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | missing | | | | |
| notmat | 8 | 100% | missing | | | | |
| notmat | 16 | 1% | 78.5458 | 16 | `369d5f7fc39874d3` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 87.7698 | 16 | `753a74a951ee286d` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 1.1691 | false |
| notmat | 1 | 1% | 0.8978 | false |
| mat | 1 | 100% | 1.0295 | false |
| notmat | 1 | 100% | 0.9493 | false |
| mat | 4 | 1% | 0.9224 | false |
| notmat | 4 | 1% | 2.4069 | false |
| mat | 4 | 100% | 1.0832 | false |
| notmat | 4 | 100% | 3.3500 | false |
| mat | 16 | 1% | 0.9798 | false |
| notmat | 16 | 1% | 8.0065 | false |
| mat | 16 | 100% | 1.0277 | false |
| notmat | 16 | 100% | 8.5933 | false |

## Physical arm

- `charity_pkey`: `CREATE UNIQUE INDEX charity_pkey ON selectjoin_domain.charity USING lsm (charity_id HASH)`
- `donation_person_idx`: `CREATE INDEX donation_person_idx ON selectjoin_domain.donation USING lsm (person_id HASH)`
- `donation_pkey`: `CREATE UNIQUE INDEX donation_pkey ON selectjoin_domain.donation USING lsm (donation_id HASH)`
- `donation_rank_idx`: `CREATE INDEX donation_rank_idx ON selectjoin_domain.donation USING lsm (rank ASC)`
- `person_pkey`: `CREATE UNIQUE INDEX person_pkey ON selectjoin_domain.person USING lsm (person_id HASH)`

- index arm `yb-range-rank`; placement `noncolocated`; node budget `6 CPU 9g`; client budget `1 CPU 2g`

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

## Throttling (cgroup `cpu.stat` deltas)

- client: nr_periods 313, nr_throttled 132, throttled_usec 5401798 | quota 100000 100000
, memory 2147483648

- yb-single: nr_periods 384, nr_throttled 0, throttled_usec 0 | quota 600000 100000, memory 9663676416

## Access path per cell

_Engine-reported from each cell's own plan: whether the rank index or a sequential scan was used, rows scanned, PostgreSQL buffers and filter removals, and reader retries of transient conflicts. A cell that scanned the whole table at a selective threshold is a different physical regime, not a faster query._

| axis | status | index scan | seq scan | storage rows scanned | index rows scanned | buffers hit/read | rows removed | retries |
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
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 111.9022 | 120.0761 | 121.2402 | 5 | 8.3 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 1.8060 | 1.9828 | 2.9389 | 5 | 504.3 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 9.3814 | 9.7466 | 10.4246 | 5 | 102.6 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 105.4617 | 108.7952 | 109.2912 | 5 | 9.2 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 86.7713 | 87.9115 | 90.2373 | 5 | 11.4 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 84.1625 | 87.5445 | 92.0709 | 5 | 11.4 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 109.8214 | 115.0683 | 122.2325 | 5 | 8.7 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 1.9428 | 2.0601 | 2.5244 | 5 | 485.4 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 9.7547 | 11.2835 | 24.1914 | 5 | 88.6 |
| family=bag;stmt=join_where;threshold=180000 | valid | 110.3630 | 112.1092 | 114.3441 | 5 | 8.9 |
| family=bag;stmt=join_where;threshold=20 | valid | 1.6797 | 1.8267 | 1.9711 | 5 | 547.4 |
| family=bag;stmt=join_where;threshold=2000 | valid | 8.7153 | 9.6112 | 9.6779 | 5 | 104.0 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 9.3365 | 9.8102 | 10.3136 | 5 | 101.9 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 9.2739 | 10.2137 | 11.6315 | 5 | 97.9 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 7.4533 | 7.7286 | 8.4786 | 5 | 129.4 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 8.8316 | 9.3265 | 9.7664 | 5 | 107.2 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 9.0428 | 9.4821 | 10.3831 | 5 | 105.5 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 8.7126 | 9.2396 | 9.9543 | 5 | 108.2 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 9.1095 | 9.6120 | 10.2198 | 5 | 104.0 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 9.9986 | 10.4970 | 11.1711 | 5 | 95.3 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 8.7678 | 9.0351 | 10.4560 | 5 | 110.7 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 9.2479 | 9.6020 | 10.1462 | 5 | 104.1 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 8.0059 | 8.7467 | 10.2858 | 5 | 114.3 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 9.4158 | 10.0081 | 10.9808 | 5 | 99.9 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 72.9822 | 78.5458 | 113.0245 | 5 | 12.7 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 84.6465 | 87.7698 | 92.5264 | 5 | 11.4 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 5.9953 | 6.9386 | 7.8492 | 5 | 144.1 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 8.2518 | 8.8536 | 9.4969 | 5 | 112.9 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 22.1227 | 22.8222 | 24.2977 | 5 | 43.8 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 29.6320 | 30.9524 | 35.7612 | 5 | 32.3 |

## Matched blended workload

- phase: `blended`; target offered rate: 10 updates/s
- offered 153, committed 153, failed 0 over 15.288 s
- offered rate 10.008/s, completed rate 10.008/s
- writer operations per endpoint: [153]
- reader operations overlapping a committed write: 69 of 210
- reader retries of transient conflicts: 0
- per key: 153 keys, min 1, max 1 per key
- ledger: 153 keys reconciled, 0 mismatch
- update latency: avg 3.305 ms, min 1.621, max 68.299 over 153

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42421353`, 5 interleaved rounds
  1. family=reuse;policy=mat;refs=4;consumer=100%
  2. family=bag;stmt=cte_joined_outer;threshold=2000
  3. family=reuse;policy=notmat;refs=16;consumer=1%
  4. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  5. family=reuse;policy=default;refs=4;consumer=100%
  6. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  7. family=reuse;policy=notmat;refs=4;consumer=100%
  8. family=reuse;policy=notmat;refs=1;consumer=1%
  9. family=bag;stmt=cte_joined_outer;threshold=20
  10. family=reuse;policy=notmat;refs=1;consumer=100%
  11. family=bag;stmt=cte_joined_outer;threshold=180000
  12. family=reuse;policy=mat;refs=1;consumer=100%
  13. family=reuse;policy=mat;refs=4;consumer=1%
  14. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  15. family=reuse;policy=mat;refs=1;consumer=1%
  16. family=reuse;policy=default;refs=16;consumer=100%
  17. family=reuse;policy=mat;refs=16;consumer=100%
  18. family=reuse;policy=mat;refs=16;consumer=1%
  19. family=reuse;policy=notmat;refs=4;consumer=1%
  20. family=bag;stmt=cte_joined_outer_mat;threshold=20
  21. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  22. family=reuse;policy=default;refs=1;consumer=1%
  23. family=bag;stmt=join_where;threshold=20
  24. family=reuse;policy=default;refs=16;consumer=1%
  25. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  26. family=reuse;policy=default;refs=1;consumer=100%
  27. family=bag;stmt=join_where;threshold=180000
  28. family=reuse;policy=default;refs=4;consumer=1%
  29. family=bag;stmt=join_where;threshold=2000
  30. family=reuse;policy=notmat;refs=16;consumer=100%

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
| family=bag;stmt=cte_joined_outer;threshold=180000 | j05_cte_joined_outer | valid | `ec6d177eeaa94273` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20 | j05_cte_joined_outer | valid | `d6a5b9fa31c93afd` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=2000 | j05_cte_joined_outer | valid | `a04e1e4704ecf67d` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | j05m_cte_joined_outer_mat | valid | `2e43e8ec9e687f52` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | j05m_cte_joined_outer_mat | valid | `413cec3a773816df` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | j05m_cte_joined_outer_mat | valid | `66154ec428cf9904` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | j05n_cte_joined_outer_notmat | valid | `9627bad998abf360` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | j05n_cte_joined_outer_notmat | valid | `f5f8f427cde84c38` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | j05n_cte_joined_outer_notmat | valid | `a22308c4df06e575` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=join_where;threshold=180000 | j01_join_where | valid | `8cb942fb3cc73ba6` |  |
| family=bag;stmt=join_where;threshold=20 | j01_join_where | valid | `3401575f7feb2252` |  |
| family=bag;stmt=join_where;threshold=2000 | j01_join_where | valid | `a4144e25011c5056` |  |
| family=reuse;policy=default;refs=16;consumer=1% | r_default_n16_c1 | valid | `e33cdbe68fbf169a` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=16;consumer=100% | r_default_n16_c100 | valid | `75a8c4dd8c0a4739` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=1;consumer=1% | r_default_n1_c1 | valid | `47e0b35694db6da5` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=1;consumer=100% | r_default_n1_c100 | valid | `038a7f7d48ffa5da` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=4;consumer=1% | r_default_n4_c1 | valid | `8f83aca0ae937564` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=100% | r_default_n4_c100 | valid | `bae65d24dfbc1ad0` | policy=default; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=1% | r_mat_n16_c1 | valid | `cf1cb3928efbf66a` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=100% | r_mat_n16_c100 | valid | `6419e021c8506f26` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=1% | r_mat_n1_c1 | valid | `29faa216d3dfa30d` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=100% | r_mat_n1_c100 | valid | `6020c01948d84f44` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=1% | r_mat_n4_c1 | valid | `0a1e88a407f861da` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=100% | r_mat_n4_c100 | valid | `bbeaa1eb4016dbd3` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=notmat;refs=16;consumer=1% | r_notmat_n16_c1 | valid | `369d5f7fc39874d3` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `753a74a951ee286d` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `e5f9bc9fee9418a1` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `1d77a6d535328dca` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `15aaa6e8e9632446` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `2c051c650a265c67` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, DIST, DEBUG, FORMAT TEXT` supported on yugabyte
- one host: this run used node budget 6 CPU 9g and client budget 1 CPU 2g, so budget comparisons are separate classes here, not independent hardware
- one client process with one fixed maintenance writer at 10 updates/s; readers are sequential; the reported sample latency includes retry and backoff time
- physical rank-index arm yb-range-rank; an engine whose planner ignores the index is a different physical design, not a faster one
- YugabyteDB DIST/DEBUG storage counters (rows and index rows scanned) are parsed per cell; RPC-level counters are not aggregated
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
