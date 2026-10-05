# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 30 cells, 30 valid, engine `yugabyte`, topology `cluster3`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-v3b-topo-cluster-i3-20261005`
- run tag: `run/study-01/select-join-v3b-topo-cluster-i3-20261005`
- environment: `host-zenbook-ux5406sa`
- repo commit: `7243fee074bfb16fc5ad8b6afc1a2d7549080935`
- engine version: `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit`
- database engine image: `docker.io/yugabytedb/yugabyte:2025.2.6.0-b111` @ `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`
- benchmark client image: `localhost/charitybench-select:select-join-v3b-topo-cluster-i3-20261005` @ `f3f1e10c44cbcb4c7b5bc56ed95a4f071682b6f7d33b5e0a7b688c4271eeef88`
- inputs digest: `7070cb55586f1625518eebf13b5817ca7852d1a870dff3ccc106fd8fcbfe87a2`
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
| 20 (0.01%) | 3.2375 | missing | missing | missing | missing | missing | 2.1286 | 133.5754 | 2.3107 | missing |
| 2000 (1%) | 12.5324 | missing | missing | missing | missing | missing | 12.3677 | 89.1868 | 11.6497 | missing |
| 180000 (90%) | 157.6872 | missing | missing | missing | missing | missing | 155.3366 | 155.5775 | 153.0432 | missing |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | cte_joined_outer vs join_where | 3.2375 | 2.1286 | 0.6575 | false |
| 20 | cte_joined_outer_mat vs join_where | 3.2375 | 133.5754 | 41.2586 | false |
| 20 | cte_joined_outer_notmat vs join_where | 3.2375 | 2.3107 | 0.7137 | false |
| 2000 | cte_joined_outer vs join_where | 12.5324 | 12.3677 | 0.9869 | false |
| 2000 | cte_joined_outer_mat vs join_where | 12.5324 | 89.1868 | 7.1165 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 12.5324 | 11.6497 | 0.9296 | false |
| 180000 | cte_joined_outer vs join_where | 157.6872 | 155.3366 | 0.9851 | false |
| 180000 | cte_joined_outer_mat vs join_where | 157.6872 | 155.5775 | 0.9866 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 157.6872 | 153.0432 | 0.9705 | false |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_joined_outer_mat vs cte_joined_outer | 2.1286 | 133.5754 | 62.7537 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 2.1286 | 2.3107 | 1.0856 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 12.3677 | 89.1868 | 7.2112 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 12.3677 | 11.6497 | 0.9419 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 155.3366 | 155.5775 | 1.0016 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 155.3366 | 153.0432 | 0.9852 | false | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 7.2377 | 1 | `9125ba2fad0e268f` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 9.6954 | 1 | `a256290fe813c6f6` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | missing | | | | |
| default | 2 | 100% | missing | | | | |
| default | 4 | 1% | 10.6321 | 4 | `4f232bfece2a62ab` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 11.0389 | 4 | `0fabfd5f89083d6e` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | missing | | | | |
| default | 8 | 100% | missing | | | | |
| default | 16 | 1% | 10.8586 | 16 | `3cfe0d6253cb1245` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 12.3066 | 16 | `636b7f8bdc557de9` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 10.5054 | 1 | `9893f9cba1f55c5c` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 8.6484 | 1 | `db49b88393f094f2` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | missing | | | | |
| mat | 2 | 100% | missing | | | | |
| mat | 4 | 1% | 9.9617 | 4 | `c3da7c34b898cd50` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 10.4243 | 4 | `461857ff39459b95` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | missing | | | | |
| mat | 8 | 100% | missing | | | | |
| mat | 16 | 1% | 9.5029 | 16 | `59ef3f630c6077c5` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 11.0011 | 16 | `034a9d4d8aef777f` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 7.9456 | 1 | `793bf5f42390908e` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 9.2810 | 1 | `3bab66a1477b1e1f` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | missing | | | | |
| notmat | 2 | 100% | missing | | | | |
| notmat | 4 | 1% | 30.2478 | 4 | `eacd91127f21835c` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 35.9026 | 4 | `310e7db5f2952c84` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | missing | | | | |
| notmat | 8 | 100% | missing | | | | |
| notmat | 16 | 1% | 83.4273 | 16 | `cad54b9949d435d7` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 99.7441 | 16 | `a58c422282e59c64` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 1.4515 | false |
| notmat | 1 | 1% | 1.0978 | false |
| mat | 1 | 100% | 0.8920 | false |
| notmat | 1 | 100% | 0.9573 | false |
| mat | 4 | 1% | 0.9370 | false |
| notmat | 4 | 1% | 2.8450 | false |
| mat | 4 | 100% | 0.9443 | false |
| notmat | 4 | 100% | 3.2524 | false |
| mat | 16 | 1% | 0.8751 | false |
| notmat | 16 | 1% | 7.6830 | false |
| mat | 16 | 100% | 0.8939 | false |
| notmat | 16 | 100% | 8.1049 | false |

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

- client: nr_periods 354, nr_throttled 128, throttled_usec 4594089 | quota 100000 100000
, memory 2147483648

- yb-n1: nr_periods 642, nr_throttled 6, throttled_usec 737504 | quota 200000 100000, memory 3221225472
- yb-n2: nr_periods 634, nr_throttled 23, throttled_usec 1814258 | quota 200000 100000, memory 3221225472
- yb-n3: nr_periods 625, nr_throttled 0, throttled_usec 0 | quota 200000 100000, memory 3221225472

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
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 147.8495 | 155.3366 | 159.0533 | 5 | 6.4 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 2.0908 | 2.1286 | 2.3598 | 5 | 469.8 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 11.6576 | 12.3677 | 13.4049 | 5 | 80.9 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 151.0314 | 155.5775 | 158.2231 | 5 | 6.4 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 130.8333 | 133.5754 | 138.6756 | 5 | 7.5 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 87.9020 | 89.1868 | 96.9910 | 5 | 11.2 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 147.5813 | 153.0432 | 169.7648 | 5 | 6.5 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 2.0144 | 2.3107 | 2.5075 | 5 | 432.8 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 10.6742 | 11.6497 | 16.3471 | 5 | 85.8 |
| family=bag;stmt=join_where;threshold=180000 | valid | 150.5547 | 157.6872 | 189.5270 | 5 | 6.3 |
| family=bag;stmt=join_where;threshold=20 | valid | 3.0947 | 3.2375 | 3.5172 | 5 | 308.9 |
| family=bag;stmt=join_where;threshold=2000 | valid | 11.9860 | 12.5324 | 13.0364 | 5 | 79.8 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 10.6908 | 10.8586 | 11.5028 | 5 | 92.1 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 11.1982 | 12.3066 | 14.2832 | 5 | 81.3 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 7.1305 | 7.2377 | 8.0676 | 5 | 138.2 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 9.2765 | 9.6954 | 10.3523 | 5 | 103.1 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 10.0017 | 10.6321 | 10.9543 | 5 | 94.1 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 10.5625 | 11.0389 | 11.3647 | 5 | 90.6 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 9.4124 | 9.5029 | 10.8267 | 5 | 105.2 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 10.4681 | 11.0011 | 11.4724 | 5 | 90.9 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 10.1469 | 10.5054 | 11.1251 | 5 | 95.2 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 8.3099 | 8.6484 | 10.1402 | 5 | 115.6 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 9.8176 | 9.9617 | 11.3061 | 5 | 100.4 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 9.7783 | 10.4243 | 11.8265 | 5 | 95.9 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 73.3380 | 83.4273 | 87.6451 | 5 | 12.0 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 96.8963 | 99.7441 | 104.2412 | 5 | 10.0 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 7.8184 | 7.9456 | 8.8514 | 5 | 125.9 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 8.8194 | 9.2810 | 9.4633 | 5 | 107.7 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 27.0754 | 30.2478 | 30.8293 | 5 | 33.1 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 34.8094 | 35.9026 | 38.0430 | 5 | 27.9 |

## Matched blended workload

- phase: `isolated`; target offered rate: 10 updates/s
- offered 0, committed 0, failed 0 over 0.000 s
- offered rate 0.000/s, completed rate 0.000/s
- per key: no keys touched
- ledger: 0 keys reconciled, 0 mismatch

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42162010`, 5 interleaved rounds
  1. family=reuse;policy=default;refs=16;consumer=100%
  2. family=reuse;policy=notmat;refs=16;consumer=1%
  3. family=bag;stmt=join_where;threshold=20
  4. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  5. family=reuse;policy=mat;refs=4;consumer=1%
  6. family=reuse;policy=notmat;refs=1;consumer=1%
  7. family=bag;stmt=cte_joined_outer_mat;threshold=20
  8. family=reuse;policy=notmat;refs=1;consumer=100%
  9. family=reuse;policy=default;refs=16;consumer=1%
  10. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  11. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  12. family=reuse;policy=notmat;refs=4;consumer=100%
  13. family=reuse;policy=mat;refs=4;consumer=100%
  14. family=bag;stmt=cte_joined_outer;threshold=20
  15. family=reuse;policy=notmat;refs=16;consumer=100%
  16. family=bag;stmt=join_where;threshold=2000
  17. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  18. family=bag;stmt=join_where;threshold=180000
  19. family=reuse;policy=mat;refs=1;consumer=1%
  20. family=reuse;policy=mat;refs=1;consumer=100%
  21. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  22. family=reuse;policy=default;refs=4;consumer=100%
  23. family=reuse;policy=default;refs=1;consumer=1%
  24. family=bag;stmt=cte_joined_outer;threshold=180000
  25. family=bag;stmt=cte_joined_outer;threshold=2000
  26. family=reuse;policy=mat;refs=16;consumer=100%
  27. family=reuse;policy=default;refs=4;consumer=1%
  28. family=reuse;policy=notmat;refs=4;consumer=1%
  29. family=reuse;policy=mat;refs=16;consumer=1%
  30. family=reuse;policy=default;refs=1;consumer=100%

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
| family=bag;stmt=cte_joined_outer;threshold=180000 | j05_cte_joined_outer | valid | `2c9f25d58be90fa7` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20 | j05_cte_joined_outer | valid | `7a995bc55891cf3b` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=2000 | j05_cte_joined_outer | valid | `6bb14333e0e49eb8` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | j05m_cte_joined_outer_mat | valid | `de67a24f2d85a726` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | j05m_cte_joined_outer_mat | valid | `4857f6cfc96f3b05` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | j05m_cte_joined_outer_mat | valid | `668ae02577c1dd53` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | j05n_cte_joined_outer_notmat | valid | `e46670c38fe41904` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | j05n_cte_joined_outer_notmat | valid | `984e189acddc25f7` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | j05n_cte_joined_outer_notmat | valid | `8d75790b0c253c03` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=join_where;threshold=180000 | j01_join_where | valid | `995baba046c6eb9c` |  |
| family=bag;stmt=join_where;threshold=20 | j01_join_where | valid | `11d7493eacda4324` |  |
| family=bag;stmt=join_where;threshold=2000 | j01_join_where | valid | `a629d103ed61bbb2` |  |
| family=reuse;policy=default;refs=16;consumer=1% | r_default_n16_c1 | valid | `3cfe0d6253cb1245` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=16;consumer=100% | r_default_n16_c100 | valid | `636b7f8bdc557de9` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=1;consumer=1% | r_default_n1_c1 | valid | `9125ba2fad0e268f` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=1;consumer=100% | r_default_n1_c100 | valid | `a256290fe813c6f6` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=4;consumer=1% | r_default_n4_c1 | valid | `4f232bfece2a62ab` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=100% | r_default_n4_c100 | valid | `0fabfd5f89083d6e` | policy=default; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=1% | r_mat_n16_c1 | valid | `59ef3f630c6077c5` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=100% | r_mat_n16_c100 | valid | `034a9d4d8aef777f` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=1% | r_mat_n1_c1 | valid | `9893f9cba1f55c5c` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=100% | r_mat_n1_c100 | valid | `db49b88393f094f2` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=1% | r_mat_n4_c1 | valid | `c3da7c34b898cd50` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=100% | r_mat_n4_c100 | valid | `461857ff39459b95` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=notmat;refs=16;consumer=1% | r_notmat_n16_c1 | valid | `cad54b9949d435d7` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `a58c422282e59c64` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `793bf5f42390908e` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `3bab66a1477b1e1f` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `eacd91127f21835c` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `310e7db5f2952c84` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, DIST, DEBUG, FORMAT TEXT` supported on yugabyte
- one host: this run used node budget 2 CPU 3g and client budget 1 CPU 2g, so budget comparisons are separate classes here, not independent hardware
- isolated phase: no concurrent writers in this run; the blended read/write phase is a separate labelled run
- cluster nodes are containers on the same host, so these are not network or independent-hardware measurements; container colocation is not data-colocation evidence
- readers round-robin across 3 endpoints; the writer's endpoint policy is recorded separately
- physical rank-index arm yb-range-rank; an engine whose planner ignores the index is a different physical design, not a faster one
- YugabyteDB DIST/DEBUG storage counters (rows and index rows scanned) are parsed per cell; RPC-level counters are not aggregated
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
