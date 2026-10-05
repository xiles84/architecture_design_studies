# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 30 cells, 30 valid, engine `yugabyte`, topology `cluster3`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-v3b-topo-cluster-b1-20261005`
- run tag: `run/study-01/select-join-v3b-topo-cluster-b1-20261005`
- environment: `host-zenbook-ux5406sa`
- repo commit: `7243fee074bfb16fc5ad8b6afc1a2d7549080935`
- engine version: `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit`
- database engine image: `docker.io/yugabytedb/yugabyte:2025.2.6.0-b111` @ `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`
- benchmark client image: `localhost/charitybench-select:select-join-v3b-topo-cluster-b1-20261005` @ `f3f1e10c44cbcb4c7b5bc56ed95a4f071682b6f7d33b5e0a7b688c4271eeef88`
- inputs digest: `0bc9da3bd9efef9dd0aa95d6ca515c58b7c9b5d41344bd39beee9c1f8b9d571d`
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
- setting `placement`: noncolocated
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
| 20 (0.01%) | 2.1715 | missing | missing | missing | missing | missing | 2.8744 | 134.0262 | 2.2663 | missing |
| 2000 (1%) | 11.2935 | missing | missing | missing | missing | missing | 10.0018 | 136.6209 | 10.6172 | missing |
| 180000 (90%) | 162.2235 | missing | missing | missing | missing | missing | 122.2483 | 150.5998 | 159.4722 | missing |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | cte_joined_outer vs join_where | 2.1715 | 2.8744 | 1.3237 | false |
| 20 | cte_joined_outer_mat vs join_where | 2.1715 | 134.0262 | 61.7213 | false |
| 20 | cte_joined_outer_notmat vs join_where | 2.1715 | 2.2663 | 1.0437 | false |
| 2000 | cte_joined_outer vs join_where | 11.2935 | 10.0018 | 0.8856 | false |
| 2000 | cte_joined_outer_mat vs join_where | 11.2935 | 136.6209 | 12.0973 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 11.2935 | 10.6172 | 0.9401 | false |
| 180000 | cte_joined_outer vs join_where | 162.2235 | 122.2483 | 0.7536 | false |
| 180000 | cte_joined_outer_mat vs join_where | 162.2235 | 150.5998 | 0.9283 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 162.2235 | 159.4722 | 0.9830 | false |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_joined_outer_mat vs cte_joined_outer | 2.8744 | 134.0262 | 46.6281 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 2.8744 | 2.2663 | 0.7885 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 10.0018 | 136.6209 | 13.6596 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 10.0018 | 10.6172 | 1.0615 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 122.2483 | 150.5998 | 1.2319 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 122.2483 | 159.4722 | 1.3045 | false | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 7.6020 | 1 | `eee2a59193377453` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 9.8510 | 1 | `757482133c9f9d14` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | missing | | | | |
| default | 2 | 100% | missing | | | | |
| default | 4 | 1% | 10.9476 | 4 | `06a27cb307797bdf` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 10.8663 | 4 | `e77d216de178f03d` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | missing | | | | |
| default | 8 | 100% | missing | | | | |
| default | 16 | 1% | 11.3908 | 16 | `a6cf3a8c9a614446` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 11.4519 | 16 | `39d47be685a5e532` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 9.7062 | 1 | `52605aa03bcecb95` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 9.1000 | 1 | `d7c2d8c9042e9d84` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | missing | | | | |
| mat | 2 | 100% | missing | | | | |
| mat | 4 | 1% | 10.0172 | 4 | `29c416f56f2deb86` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 11.0697 | 4 | `8a8b1a64464640df` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | missing | | | | |
| mat | 8 | 100% | missing | | | | |
| mat | 16 | 1% | 11.2329 | 16 | `34322e34c4c2d2fd` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 11.3404 | 16 | `1e3f5534cc7665a5` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 8.9580 | 1 | `a3b56b2eabfa077f` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 10.1255 | 1 | `646113dcf4bbcfe8` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | missing | | | | |
| notmat | 2 | 100% | missing | | | | |
| notmat | 4 | 1% | 32.3579 | 4 | `f3b8fa5440018661` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 36.9767 | 4 | `dbda443d9ea6bcb9` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | missing | | | | |
| notmat | 8 | 100% | missing | | | | |
| notmat | 16 | 1% | 92.7883 | 16 | `d35ffbbe5265b38a` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 99.7729 | 16 | `98bf2654cb17ec9c` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 1.2768 | false |
| notmat | 1 | 1% | 1.1784 | false |
| mat | 1 | 100% | 0.9238 | false |
| notmat | 1 | 100% | 1.0279 | false |
| mat | 4 | 1% | 0.9150 | false |
| notmat | 4 | 1% | 2.9557 | false |
| mat | 4 | 100% | 1.0187 | false |
| notmat | 4 | 100% | 3.4029 | false |
| mat | 16 | 1% | 0.9861 | false |
| notmat | 16 | 1% | 8.1459 | false |
| mat | 16 | 100% | 0.9903 | false |
| notmat | 16 | 100% | 8.7123 | false |

## Physical arm

- `charity_pkey`: `CREATE UNIQUE INDEX charity_pkey ON selectjoin_domain.charity USING lsm (charity_id HASH)`
- `donation_person_idx`: `CREATE INDEX donation_person_idx ON selectjoin_domain.donation USING lsm (person_id HASH)`
- `donation_pkey`: `CREATE UNIQUE INDEX donation_pkey ON selectjoin_domain.donation USING lsm (donation_id HASH)`
- `donation_rank_idx`: `CREATE INDEX donation_rank_idx ON selectjoin_domain.donation USING lsm (rank ASC)`
- `person_pkey`: `CREATE UNIQUE INDEX person_pkey ON selectjoin_domain.person USING lsm (person_id HASH)`

- index arm `yb-range-rank`; placement `noncolocated`; node budget `2 CPU 3g`; client budget `1 CPU 2g`

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

- client: nr_periods 348, nr_throttled 131, throttled_usec 5059249 | quota 100000 100000
, memory 2147483648

- yb-n1: nr_periods 643, nr_throttled 7, throttled_usec 702031 | quota 200000 100000, memory 3221225472
- yb-n2: nr_periods 633, nr_throttled 5, throttled_usec 656221 | quota 200000 100000, memory 3221225472
- yb-n3: nr_periods 629, nr_throttled 11, throttled_usec 314027 | quota 200000 100000, memory 3221225472

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
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 117.6292 | 122.2483 | 125.4995 | 5 | 8.2 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 2.7024 | 2.8744 | 3.9996 | 5 | 347.9 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 9.4350 | 10.0018 | 10.3150 | 5 | 100.0 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 149.5948 | 150.5998 | 167.7334 | 5 | 6.6 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 127.9079 | 134.0262 | 151.3818 | 5 | 7.5 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 131.1324 | 136.6209 | 155.7175 | 5 | 7.3 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 156.0638 | 159.4722 | 169.7512 | 5 | 6.3 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 2.1521 | 2.2663 | 2.3348 | 5 | 441.2 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 10.1923 | 10.6172 | 11.1515 | 5 | 94.2 |
| family=bag;stmt=join_where;threshold=180000 | valid | 154.4333 | 162.2235 | 165.8502 | 5 | 6.2 |
| family=bag;stmt=join_where;threshold=20 | valid | 1.9894 | 2.1715 | 2.2850 | 5 | 460.5 |
| family=bag;stmt=join_where;threshold=2000 | valid | 10.7658 | 11.2935 | 11.6051 | 5 | 88.5 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 10.7445 | 11.3908 | 12.2985 | 5 | 87.8 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 11.1387 | 11.4519 | 12.8011 | 5 | 87.3 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 7.1310 | 7.6020 | 8.1704 | 5 | 131.5 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 9.7618 | 9.8510 | 10.0400 | 5 | 101.5 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 10.5982 | 10.9476 | 11.0484 | 5 | 91.3 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 10.4161 | 10.8663 | 11.0095 | 5 | 92.0 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 10.6815 | 11.2329 | 12.0269 | 5 | 89.0 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 11.1457 | 11.3404 | 11.6824 | 5 | 88.2 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 9.2820 | 9.7062 | 10.3097 | 5 | 103.0 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 8.5996 | 9.1000 | 9.2918 | 5 | 109.9 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 9.7384 | 10.0172 | 10.8525 | 5 | 99.8 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 10.2453 | 11.0697 | 12.1915 | 5 | 90.3 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 85.0662 | 92.7883 | 114.4137 | 5 | 10.8 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 95.1152 | 99.7729 | 115.6970 | 5 | 10.0 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 8.7573 | 8.9580 | 9.1211 | 5 | 111.6 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 9.2928 | 10.1255 | 13.9937 | 5 | 98.8 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 31.0321 | 32.3579 | 39.0273 | 5 | 30.9 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 34.9567 | 36.9767 | 38.5434 | 5 | 27.0 |

## Matched blended workload

- phase: `blended`; target offered rate: 10 updates/s
- offered 176, committed 176, failed 0 over 17.562 s
- offered rate 10.022/s, completed rate 10.022/s
- writer operations per endpoint: [59 59 58]
- reader operations overlapping a committed write: 68 of 210
- reader retries of transient conflicts: 0
- per key: 176 keys, min 1, max 1 per key
- ledger: 176 keys reconciled, 0 mismatch
- update latency: avg 6.459 ms, min 2.532, max 196.417 over 176

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42431313`, 5 interleaved rounds
  1. family=bag;stmt=cte_joined_outer;threshold=20
  2. family=bag;stmt=cte_joined_outer_mat;threshold=20
  3. family=reuse;policy=default;refs=1;consumer=1%
  4. family=bag;stmt=join_where;threshold=2000
  5. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  6. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  7. family=reuse;policy=mat;refs=1;consumer=1%
  8. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  9. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  10. family=reuse;policy=notmat;refs=4;consumer=100%
  11. family=reuse;policy=notmat;refs=1;consumer=100%
  12. family=bag;stmt=cte_joined_outer;threshold=2000
  13. family=reuse;policy=notmat;refs=16;consumer=100%
  14. family=reuse;policy=notmat;refs=1;consumer=1%
  15. family=reuse;policy=notmat;refs=16;consumer=1%
  16. family=reuse;policy=default;refs=4;consumer=1%
  17. family=reuse;policy=notmat;refs=4;consumer=1%
  18. family=bag;stmt=join_where;threshold=20
  19. family=reuse;policy=mat;refs=4;consumer=100%
  20. family=bag;stmt=join_where;threshold=180000
  21. family=reuse;policy=mat;refs=4;consumer=1%
  22. family=reuse;policy=default;refs=16;consumer=1%
  23. family=reuse;policy=default;refs=1;consumer=100%
  24. family=bag;stmt=cte_joined_outer;threshold=180000
  25. family=reuse;policy=default;refs=4;consumer=100%
  26. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  27. family=reuse;policy=mat;refs=16;consumer=100%
  28. family=reuse;policy=mat;refs=16;consumer=1%
  29. family=reuse;policy=default;refs=16;consumer=100%
  30. family=reuse;policy=mat;refs=1;consumer=100%

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
| family=bag;stmt=cte_joined_outer;threshold=180000 | j05_cte_joined_outer | valid | `e2b0ea47f0487e3c` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20 | j05_cte_joined_outer | valid | `51db4eab80574f5e` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=2000 | j05_cte_joined_outer | valid | `64d4d52121d92d81` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | j05m_cte_joined_outer_mat | valid | `1b85638e8f6485e1` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | j05m_cte_joined_outer_mat | valid | `7b761f0da2af406b` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | j05m_cte_joined_outer_mat | valid | `84ad5deffafc0bd1` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | j05n_cte_joined_outer_notmat | valid | `1055f5f616ee4591` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | j05n_cte_joined_outer_notmat | valid | `db30426eb645b7b8` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | j05n_cte_joined_outer_notmat | valid | `c3e26239c5e1e725` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=join_where;threshold=180000 | j01_join_where | valid | `c084411c769f2c9d` |  |
| family=bag;stmt=join_where;threshold=20 | j01_join_where | valid | `ce1d7ec3d299a779` |  |
| family=bag;stmt=join_where;threshold=2000 | j01_join_where | valid | `96c20ea120e3dcfa` |  |
| family=reuse;policy=default;refs=16;consumer=1% | r_default_n16_c1 | valid | `a6cf3a8c9a614446` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=16;consumer=100% | r_default_n16_c100 | valid | `39d47be685a5e532` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=1;consumer=1% | r_default_n1_c1 | valid | `eee2a59193377453` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=1;consumer=100% | r_default_n1_c100 | valid | `757482133c9f9d14` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=4;consumer=1% | r_default_n4_c1 | valid | `06a27cb307797bdf` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=100% | r_default_n4_c100 | valid | `e77d216de178f03d` | policy=default; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=1% | r_mat_n16_c1 | valid | `34322e34c4c2d2fd` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=100% | r_mat_n16_c100 | valid | `1e3f5534cc7665a5` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=1% | r_mat_n1_c1 | valid | `52605aa03bcecb95` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=100% | r_mat_n1_c100 | valid | `d7c2d8c9042e9d84` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=1% | r_mat_n4_c1 | valid | `29c416f56f2deb86` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=100% | r_mat_n4_c100 | valid | `8a8b1a64464640df` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=notmat;refs=16;consumer=1% | r_notmat_n16_c1 | valid | `d35ffbbe5265b38a` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `98bf2654cb17ec9c` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `a3b56b2eabfa077f` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `646113dcf4bbcfe8` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `f3b8fa5440018661` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `dbda443d9ea6bcb9` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, DIST, DEBUG, FORMAT TEXT` supported on yugabyte
- one host: this run used node budget 2 CPU 3g and client budget 1 CPU 2g, so budget comparisons are separate classes here, not independent hardware
- one client process with one fixed maintenance writer at 10 updates/s; readers are sequential; the reported sample latency includes retry and backoff time
- cluster nodes are containers on the same host, so these are not network or independent-hardware measurements; container colocation is not data-colocation evidence
- readers round-robin across 3 endpoints; the writer's endpoint policy is recorded separately
- physical rank-index arm yb-range-rank; an engine whose planner ignores the index is a different physical design, not a faster one
- YugabyteDB DIST/DEBUG storage counters (rows and index rows scanned) are parsed per cell; RPC-level counters are not aggregated
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
