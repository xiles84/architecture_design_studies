# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 30 cells, 30 valid, engine `yugabyte`, topology `cluster3`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-v3b-coloc-b3-20261005`
- run tag: `run/study-01/select-join-v3b-coloc-b3-20261005`
- environment: `host-zenbook-ux5406sa`
- repo commit: `7243fee074bfb16fc5ad8b6afc1a2d7549080935`
- engine version: `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit`
- database engine image: `docker.io/yugabytedb/yugabyte:2025.2.6.0-b111` @ `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`
- benchmark client image: `localhost/charitybench-select:select-join-v3b-coloc-b3-20261005` @ `f3f1e10c44cbcb4c7b5bc56ed95a4f071682b6f7d33b5e0a7b688c4271eeef88`
- inputs digest: `6aafb9d2ead32bd8c148af1b31e67ee003ec26e675ccc0c85659737448df1acc`
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
| 20 (0.01%) | 2.4475 | missing | missing | missing | missing | missing | 2.0925 | 120.2742 | 2.1770 | missing |
| 2000 (1%) | 7.2744 | missing | missing | missing | missing | missing | 7.5938 | 80.3757 | 7.8713 | missing |
| 180000 (90%) | 141.1788 | missing | missing | missing | missing | missing | 140.0713 | 109.2234 | 145.1275 | missing |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | cte_joined_outer vs join_where | 2.4475 | 2.0925 | 0.8550 | false |
| 20 | cte_joined_outer_mat vs join_where | 2.4475 | 120.2742 | 49.1414 | false |
| 20 | cte_joined_outer_notmat vs join_where | 2.4475 | 2.1770 | 0.8895 | false |
| 2000 | cte_joined_outer vs join_where | 7.2744 | 7.5938 | 1.0439 | false |
| 2000 | cte_joined_outer_mat vs join_where | 7.2744 | 80.3757 | 11.0490 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 7.2744 | 7.8713 | 1.0820 | false |
| 180000 | cte_joined_outer vs join_where | 141.1788 | 140.0713 | 0.9922 | false |
| 180000 | cte_joined_outer_mat vs join_where | 141.1788 | 109.2234 | 0.7737 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 141.1788 | 145.1275 | 1.0280 | false |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_joined_outer_mat vs cte_joined_outer | 2.0925 | 120.2742 | 57.4783 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 2.0925 | 2.1770 | 1.0404 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 7.5938 | 80.3757 | 10.5844 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 7.5938 | 7.8713 | 1.0365 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 140.0713 | 109.2234 | 0.7798 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 140.0713 | 145.1275 | 1.0361 | false | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 1.6301 | 1 | `cf3aa50d4ec9e4d3` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 7.0792 | 1 | `e5c634fed9262062` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | missing | | | | |
| default | 2 | 100% | missing | | | | |
| default | 4 | 1% | 7.3308 | 4 | `6f8a3534869db956` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 6.7632 | 4 | `3737a7cf3adcd48e` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | missing | | | | |
| default | 8 | 100% | missing | | | | |
| default | 16 | 1% | 8.2549 | 16 | `93bb3cf9ff32021b` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 8.5219 | 16 | `e75b9fd014f67eda` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 6.2293 | 1 | `e212ae8224d9d514` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 7.5868 | 1 | `90a28fffa294f9a4` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | missing | | | | |
| mat | 2 | 100% | missing | | | | |
| mat | 4 | 1% | 7.0739 | 4 | `86114ed11c45defd` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 7.6239 | 4 | `a7740c6a9ed7b0a5` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | missing | | | | |
| mat | 8 | 100% | missing | | | | |
| mat | 16 | 1% | 8.2010 | 16 | `b15face56512108a` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 8.5150 | 16 | `8c2cc199f47dcd9e` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 1.7691 | 1 | `ec272556ffe6d1b3` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 6.4264 | 1 | `d0e0fdf0c7fc43fd` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | missing | | | | |
| notmat | 2 | 100% | missing | | | | |
| notmat | 4 | 1% | 4.8691 | 4 | `b089e8e8452d6a7b` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 23.8457 | 4 | `b7d67f6d400e0b84` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | missing | | | | |
| notmat | 8 | 100% | missing | | | | |
| notmat | 16 | 1% | 17.3321 | 16 | `21e4fbe4139008ec` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 89.5045 | 16 | `1f2db2dfc81c6642` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 3.8214 | false |
| notmat | 1 | 1% | 1.0852 | false |
| mat | 1 | 100% | 1.0717 | false |
| notmat | 1 | 100% | 0.9078 | false |
| mat | 4 | 1% | 0.9650 | false |
| notmat | 4 | 1% | 0.6642 | false |
| mat | 4 | 100% | 1.1272 | false |
| notmat | 4 | 100% | 3.5258 | false |
| mat | 16 | 1% | 0.9935 | false |
| notmat | 16 | 1% | 2.0996 | false |
| mat | 16 | 100% | 0.9992 | false |
| notmat | 16 | 100% | 10.5029 | false |

## Physical arm

- `charity_pkey`: `CREATE UNIQUE INDEX charity_pkey ON selectjoin_domain.charity USING lsm (charity_id ASC)`
- `donation_person_idx`: `CREATE INDEX donation_person_idx ON selectjoin_domain.donation USING lsm (person_id ASC)`
- `donation_pkey`: `CREATE UNIQUE INDEX donation_pkey ON selectjoin_domain.donation USING lsm (donation_id ASC)`
- `donation_rank_idx`: `CREATE INDEX donation_rank_idx ON selectjoin_domain.donation USING lsm (rank ASC)`
- `person_pkey`: `CREATE UNIQUE INDEX person_pkey ON selectjoin_domain.person USING lsm (person_id ASC)`

- index arm `yb-range-rank`; placement `colocated`; node budget `2 CPU 3g`; client budget `1 CPU 2g`

| table | tablets | hash key columns | colocated | tablegroup | colocation id |
|---|---|---|---|---|---|
| selectjoin_domain.charity | 1 | 0 | true | 16388 | 1509739706 |
| selectjoin_domain.person | 1 | 0 | true | 16388 | 830521536 |
| selectjoin_domain.donation | 1 | 0 | true | 16388 | 1962647520 |
| selectjoin_domain.donation_pkey | 0 | 0 | false | unknown | unknown |
| selectjoin_domain.donation_rank_idx | 1 | 0 | true | 16388 | 2118599568 |
| selectjoin_domain.donation_person_idx | 1 | 0 | true | 16388 | 1760435484 |

Physical-arm gate (default join must use the rank index at its three most selective thresholds):

- threshold 20: used
- threshold 2000: used

## Throttling (cgroup `cpu.stat` deltas)

- client: nr_periods 293, nr_throttled 94, throttled_usec 4855464 | quota 100000 100000
, memory 2147483648

- yb-n1: nr_periods 545, nr_throttled 5, throttled_usec 827497 | quota 200000 100000, memory 3221225472
- yb-n2: nr_periods 536, nr_throttled 6, throttled_usec 771131 | quota 200000 100000, memory 3221225472
- yb-n3: nr_periods 532, nr_throttled 5, throttled_usec 95830 | quota 200000 100000, memory 3221225472

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
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 135.5069 | 140.0713 | 150.4988 | 5 | 7.1 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 1.8128 | 2.0925 | 2.4251 | 5 | 477.9 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 7.2093 | 7.5938 | 8.0733 | 5 | 131.7 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 102.0427 | 109.2234 | 111.6728 | 5 | 9.2 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 114.4467 | 120.2742 | 131.1581 | 5 | 8.3 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 77.3082 | 80.3757 | 83.4223 | 5 | 12.4 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 139.0303 | 145.1275 | 153.6073 | 5 | 6.9 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 2.1193 | 2.1770 | 2.2962 | 5 | 459.3 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 6.9435 | 7.8713 | 9.0814 | 5 | 127.0 |
| family=bag;stmt=join_where;threshold=180000 | valid | 136.3974 | 141.1788 | 145.9838 | 5 | 7.1 |
| family=bag;stmt=join_where;threshold=20 | valid | 2.3171 | 2.4475 | 2.6200 | 5 | 408.6 |
| family=bag;stmt=join_where;threshold=2000 | valid | 6.9921 | 7.2744 | 7.9526 | 5 | 137.5 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 7.7208 | 8.2549 | 8.3543 | 5 | 121.1 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 8.0930 | 8.5219 | 9.1922 | 5 | 117.3 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 1.5568 | 1.6301 | 1.7156 | 5 | 613.5 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 6.9000 | 7.0792 | 7.6793 | 5 | 141.3 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 6.8248 | 7.3308 | 10.0347 | 5 | 136.4 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 6.7050 | 6.7632 | 7.6251 | 5 | 147.9 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 7.5859 | 8.2010 | 29.8330 | 5 | 121.9 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 8.3574 | 8.5150 | 9.2464 | 5 | 117.4 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 6.0172 | 6.2293 | 9.3399 | 5 | 160.5 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 7.4920 | 7.5868 | 7.7449 | 5 | 131.8 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 6.8109 | 7.0739 | 31.8740 | 5 | 141.4 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 6.9934 | 7.6239 | 7.8047 | 5 | 131.2 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 16.9847 | 17.3321 | 19.0844 | 5 | 57.7 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 86.7978 | 89.5045 | 97.7097 | 5 | 11.2 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 1.6865 | 1.7691 | 2.1069 | 5 | 565.3 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 6.4083 | 6.4264 | 6.7137 | 5 | 155.6 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 4.7898 | 4.8691 | 5.1966 | 5 | 205.4 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 21.9625 | 23.8457 | 27.1367 | 5 | 41.9 |

## Matched blended workload

- phase: `blended`; target offered rate: 10 updates/s
- offered 133, committed 133, failed 0 over 13.237 s
- offered rate 10.048/s, completed rate 10.048/s
- writer operations per endpoint: [45 44 44]
- reader operations overlapping a committed write: 61 of 210
- reader retries of transient conflicts: 0
- per key: 133 keys, min 1, max 1 per key
- ledger: 133 keys reconciled, 0 mismatch
- update latency: avg 7.728 ms, min 2.504, max 81.912 over 133

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42787714`, 5 interleaved rounds
  1. family=reuse;policy=default;refs=1;consumer=100%
  2. family=reuse;policy=notmat;refs=16;consumer=100%
  3. family=bag;stmt=join_where;threshold=2000
  4. family=reuse;policy=notmat;refs=4;consumer=1%
  5. family=reuse;policy=mat;refs=4;consumer=100%
  6. family=reuse;policy=mat;refs=1;consumer=1%
  7. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  8. family=reuse;policy=notmat;refs=1;consumer=1%
  9. family=reuse;policy=default;refs=4;consumer=100%
  10. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  11. family=reuse;policy=mat;refs=1;consumer=100%
  12. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  13. family=reuse;policy=notmat;refs=16;consumer=1%
  14. family=bag;stmt=cte_joined_outer;threshold=2000
  15. family=bag;stmt=cte_joined_outer;threshold=20
  16. family=reuse;policy=mat;refs=16;consumer=100%
  17. family=reuse;policy=default;refs=16;consumer=100%
  18. family=reuse;policy=notmat;refs=4;consumer=100%
  19. family=bag;stmt=cte_joined_outer_mat;threshold=20
  20. family=bag;stmt=join_where;threshold=180000
  21. family=reuse;policy=mat;refs=4;consumer=1%
  22. family=reuse;policy=default;refs=16;consumer=1%
  23. family=bag;stmt=cte_joined_outer;threshold=180000
  24. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  25. family=reuse;policy=default;refs=1;consumer=1%
  26. family=bag;stmt=join_where;threshold=20
  27. family=reuse;policy=notmat;refs=1;consumer=100%
  28. family=reuse;policy=mat;refs=16;consumer=1%
  29. family=reuse;policy=default;refs=4;consumer=1%
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
| family=bag;stmt=cte_joined_outer;threshold=180000 | j05_cte_joined_outer | valid | `5990b47e09ab9d58` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20 | j05_cte_joined_outer | valid | `5a72a9c889860d0c` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=2000 | j05_cte_joined_outer | valid | `c45ac642589f4862` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | j05m_cte_joined_outer_mat | valid | `8b6ef77e9442f070` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | j05m_cte_joined_outer_mat | valid | `d70e510b167d70a8` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | j05m_cte_joined_outer_mat | valid | `c957f757070ab3d6` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | j05n_cte_joined_outer_notmat | valid | `80c39cc33f17e0a8` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | j05n_cte_joined_outer_notmat | valid | `b498dd4b3240d849` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | j05n_cte_joined_outer_notmat | valid | `28b56e2b09b523bb` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=join_where;threshold=180000 | j01_join_where | valid | `89617d4abb4bc788` |  |
| family=bag;stmt=join_where;threshold=20 | j01_join_where | valid | `4fbadb3e57f74d47` |  |
| family=bag;stmt=join_where;threshold=2000 | j01_join_where | valid | `987664e98ae864c2` |  |
| family=reuse;policy=default;refs=16;consumer=1% | r_default_n16_c1 | valid | `93bb3cf9ff32021b` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=16;consumer=100% | r_default_n16_c100 | valid | `e75b9fd014f67eda` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=1;consumer=1% | r_default_n1_c1 | valid | `cf3aa50d4ec9e4d3` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=1;consumer=100% | r_default_n1_c100 | valid | `e5c634fed9262062` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=4;consumer=1% | r_default_n4_c1 | valid | `6f8a3534869db956` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=100% | r_default_n4_c100 | valid | `3737a7cf3adcd48e` | policy=default; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=1% | r_mat_n16_c1 | valid | `b15face56512108a` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=100% | r_mat_n16_c100 | valid | `8c2cc199f47dcd9e` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=1% | r_mat_n1_c1 | valid | `e212ae8224d9d514` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=100% | r_mat_n1_c100 | valid | `90a28fffa294f9a4` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=1% | r_mat_n4_c1 | valid | `86114ed11c45defd` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=100% | r_mat_n4_c100 | valid | `a7740c6a9ed7b0a5` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=notmat;refs=16;consumer=1% | r_notmat_n16_c1 | valid | `21e4fbe4139008ec` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `1f2db2dfc81c6642` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `ec272556ffe6d1b3` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `d0e0fdf0c7fc43fd` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `b089e8e8452d6a7b` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `b7d67f6d400e0b84` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, DIST, DEBUG, FORMAT TEXT` supported on yugabyte
- one host: this run used node budget 2 CPU 3g and client budget 1 CPU 2g, so budget comparisons are separate classes here, not independent hardware
- one client process with one fixed maintenance writer at 10 updates/s; readers are sequential; the reported sample latency includes retry and backoff time
- cluster nodes are containers on the same host, so these are not network or independent-hardware measurements; container colocation is not data-colocation evidence
- readers round-robin across 3 endpoints; the writer's endpoint policy is recorded separately
- physical rank-index arm yb-range-rank; an engine whose planner ignores the index is a different physical design, not a faster one
- YugabyteDB DIST/DEBUG storage counters (rows and index rows scanned) are parsed per cell; RPC-level counters are not aggregated
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
