# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 30 cells, 30 valid, engine `yugabyte`, topology `cluster3`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-v3b-coloc-b1-20261005`
- run tag: `run/study-01/select-join-v3b-coloc-b1-20261005`
- environment: `host-zenbook-ux5406sa`
- repo commit: `7243fee074bfb16fc5ad8b6afc1a2d7549080935`
- engine version: `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit`
- database engine image: `docker.io/yugabytedb/yugabyte:2025.2.6.0-b111` @ `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`
- benchmark client image: `localhost/charitybench-select:select-join-v3b-coloc-b1-20261005` @ `f3f1e10c44cbcb4c7b5bc56ed95a4f071682b6f7d33b5e0a7b688c4271eeef88`
- inputs digest: `d33fd2a69e57be9a1a2ebfec6d37dd23a4af7ec621413a86d8849ad1e4ae0849`
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
| 20 (0.01%) | 2.4074 | missing | missing | missing | missing | missing | 2.1549 | 118.7996 | 2.7018 | missing |
| 2000 (1%) | 7.9846 | missing | missing | missing | missing | missing | 7.9813 | 121.9556 | 7.6128 | missing |
| 180000 (90%) | 115.9706 | missing | missing | missing | missing | missing | 140.6864 | 104.2239 | 105.9562 | missing |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | cte_joined_outer vs join_where | 2.4074 | 2.1549 | 0.8951 | false |
| 20 | cte_joined_outer_mat vs join_where | 2.4074 | 118.7996 | 49.3471 | false |
| 20 | cte_joined_outer_notmat vs join_where | 2.4074 | 2.7018 | 1.1223 | false |
| 2000 | cte_joined_outer vs join_where | 7.9846 | 7.9813 | 0.9996 | false |
| 2000 | cte_joined_outer_mat vs join_where | 7.9846 | 121.9556 | 15.2739 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 7.9846 | 7.6128 | 0.9534 | false |
| 180000 | cte_joined_outer vs join_where | 115.9706 | 140.6864 | 1.2131 | false |
| 180000 | cte_joined_outer_mat vs join_where | 115.9706 | 104.2239 | 0.8987 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 115.9706 | 105.9562 | 0.9136 | false |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_joined_outer_mat vs cte_joined_outer | 2.1549 | 118.7996 | 55.1304 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 2.1549 | 2.7018 | 1.2538 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 7.9813 | 121.9556 | 15.2802 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 7.9813 | 7.6128 | 0.9538 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 140.6864 | 104.2239 | 0.7408 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 140.6864 | 105.9562 | 0.7531 | false | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 1.4230 | 1 | `9a64a0613768050f` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 7.6359 | 1 | `8f90da054f9f2aca` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | missing | | | | |
| default | 2 | 100% | missing | | | | |
| default | 4 | 1% | 7.2355 | 4 | `f12d26b66a8684fe` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 6.3531 | 4 | `31855c33f8c7881b` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | missing | | | | |
| default | 8 | 100% | missing | | | | |
| default | 16 | 1% | 8.2966 | 16 | `4cf84397665b0d35` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 8.4698 | 16 | `107e20cb6bcfdbfb` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 6.8552 | 1 | `c1b22c44076b354f` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 6.9557 | 1 | `ff91c8711b7e3627` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | missing | | | | |
| mat | 2 | 100% | missing | | | | |
| mat | 4 | 1% | 6.4670 | 4 | `d191d915dc30144c` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 7.3732 | 4 | `907ab5476f8ad1f5` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | missing | | | | |
| mat | 8 | 100% | missing | | | | |
| mat | 16 | 1% | 7.8036 | 16 | `9f092f705fac93bb` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 9.2489 | 16 | `64ec315ef647ef2d` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 1.6953 | 1 | `8029c01188eb7447` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 7.0788 | 1 | `62f59d83157809d4` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | missing | | | | |
| notmat | 2 | 100% | missing | | | | |
| notmat | 4 | 1% | 3.8805 | 4 | `39f6874acd801497` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 24.6238 | 4 | `25e554fe591f4ea9` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | missing | | | | |
| notmat | 8 | 100% | missing | | | | |
| notmat | 16 | 1% | 16.9441 | 16 | `ebd81377094a97f2` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 86.9239 | 16 | `543d906da363acbe` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 4.8176 | false |
| notmat | 1 | 1% | 1.1914 | false |
| mat | 1 | 100% | 0.9109 | false |
| notmat | 1 | 100% | 0.9270 | false |
| mat | 4 | 1% | 0.8938 | false |
| notmat | 4 | 1% | 0.5363 | false |
| mat | 4 | 100% | 1.1606 | false |
| notmat | 4 | 100% | 3.8758 | false |
| mat | 16 | 1% | 0.9406 | false |
| notmat | 16 | 1% | 2.0423 | false |
| mat | 16 | 100% | 1.0920 | false |
| notmat | 16 | 100% | 10.2628 | false |

## Physical arm

- `charity_pkey`: `CREATE UNIQUE INDEX charity_pkey ON selectjoin_domain.charity USING lsm (charity_id ASC)`
- `donation_person_idx`: `CREATE INDEX donation_person_idx ON selectjoin_domain.donation USING lsm (person_id ASC)`
- `donation_pkey`: `CREATE UNIQUE INDEX donation_pkey ON selectjoin_domain.donation USING lsm (donation_id ASC)`
- `donation_rank_idx`: `CREATE INDEX donation_rank_idx ON selectjoin_domain.donation USING lsm (rank ASC)`
- `person_pkey`: `CREATE UNIQUE INDEX person_pkey ON selectjoin_domain.person USING lsm (person_id ASC)`

- index arm `yb-range-rank`; placement `colocated`; node budget `2 CPU 3g`; client budget `1 CPU 2g`

| table | tablets | hash key columns | colocated | tablegroup | colocation id |
|---|---|---|---|---|---|
| selectjoin_domain.charity | 1 | 0 | true | 16388 | 2752046158 |
| selectjoin_domain.person | 1 | 0 | true | 16388 | 47765700 |
| selectjoin_domain.donation | 1 | 0 | true | 16388 | 2484824743 |
| selectjoin_domain.donation_pkey | 0 | 0 | false | unknown | unknown |
| selectjoin_domain.donation_rank_idx | 1 | 0 | true | 16388 | 3545894751 |
| selectjoin_domain.donation_person_idx | 1 | 0 | true | 16388 | 1863026618 |

Physical-arm gate (default join must use the rank index at its three most selective thresholds):

- threshold 20: used
- threshold 2000: used

## Throttling (cgroup `cpu.stat` deltas)

- client: nr_periods 283, nr_throttled 99, throttled_usec 5328518 | quota 100000 100000
, memory 2147483648

- yb-n1: nr_periods 538, nr_throttled 5, throttled_usec 781983 | quota 200000 100000, memory 3221225472
- yb-n2: nr_periods 531, nr_throttled 9, throttled_usec 1719228 | quota 200000 100000, memory 3221225472
- yb-n3: nr_periods 523, nr_throttled 0, throttled_usec 0 | quota 200000 100000, memory 3221225472

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
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 132.9635 | 140.6864 | 143.4338 | 5 | 7.1 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 2.0690 | 2.1549 | 5.4466 | 5 | 464.1 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 7.7401 | 7.9813 | 8.7429 | 5 | 125.3 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 103.1348 | 104.2239 | 109.3476 | 5 | 9.6 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 111.9005 | 118.7996 | 138.5957 | 5 | 8.4 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 117.3512 | 121.9556 | 122.8076 | 5 | 8.2 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 103.9260 | 105.9562 | 108.7416 | 5 | 9.4 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 2.4103 | 2.7018 | 3.1405 | 5 | 370.1 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 7.2124 | 7.6128 | 9.1741 | 5 | 131.4 |
| family=bag;stmt=join_where;threshold=180000 | valid | 107.9316 | 115.9706 | 127.2457 | 5 | 8.6 |
| family=bag;stmt=join_where;threshold=20 | valid | 2.3195 | 2.4074 | 3.1638 | 5 | 415.4 |
| family=bag;stmt=join_where;threshold=2000 | valid | 7.8271 | 7.9846 | 8.7413 | 5 | 125.2 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 7.5842 | 8.2966 | 8.4766 | 5 | 120.5 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 8.2345 | 8.4698 | 8.9690 | 5 | 118.1 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 1.3216 | 1.4230 | 1.8601 | 5 | 702.8 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 7.3026 | 7.6359 | 8.0159 | 5 | 131.0 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 7.1513 | 7.2355 | 7.6658 | 5 | 138.2 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 6.1986 | 6.3531 | 6.8375 | 5 | 157.4 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 7.5824 | 7.8036 | 8.0344 | 5 | 128.1 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 8.3920 | 9.2489 | 11.8807 | 5 | 108.1 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 6.2166 | 6.8552 | 9.7034 | 5 | 145.9 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 6.8648 | 6.9557 | 9.4728 | 5 | 143.8 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 6.0160 | 6.4670 | 6.7588 | 5 | 154.6 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 7.2827 | 7.3732 | 7.6070 | 5 | 135.6 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 16.6247 | 16.9441 | 17.4138 | 5 | 59.0 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 83.1610 | 86.9239 | 87.6058 | 5 | 11.5 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 1.6298 | 1.6953 | 1.7356 | 5 | 589.9 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 6.6224 | 7.0788 | 9.5955 | 5 | 141.3 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 3.7981 | 3.8805 | 4.0067 | 5 | 257.7 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 24.3855 | 24.6238 | 28.3802 | 5 | 40.6 |

## Matched blended workload

- phase: `blended`; target offered rate: 10 updates/s
- offered 128, committed 128, failed 0 over 12.793 s
- offered rate 10.005/s, completed rate 10.005/s
- writer operations per endpoint: [43 43 42]
- reader operations overlapping a committed write: 61 of 210
- reader retries of transient conflicts: 0
- per key: 128 keys, min 1, max 1 per key
- ledger: 128 keys reconciled, 0 mismatch
- update latency: avg 8.735 ms, min 2.545, max 186.452 over 128

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42047964`, 5 interleaved rounds
  1. family=bag;stmt=cte_joined_outer;threshold=180000
  2. family=bag;stmt=cte_joined_outer;threshold=20
  3. family=reuse;policy=default;refs=4;consumer=1%
  4. family=reuse;policy=notmat;refs=1;consumer=1%
  5. family=reuse;policy=default;refs=1;consumer=1%
  6. family=bag;stmt=cte_joined_outer_mat;threshold=20
  7. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  8. family=reuse;policy=mat;refs=1;consumer=1%
  9. family=reuse;policy=mat;refs=16;consumer=100%
  10. family=reuse;policy=default;refs=16;consumer=100%
  11. family=reuse;policy=default;refs=4;consumer=100%
  12. family=reuse;policy=notmat;refs=16;consumer=100%
  13. family=bag;stmt=join_where;threshold=2000
  14. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  15. family=reuse;policy=notmat;refs=16;consumer=1%
  16. family=reuse;policy=mat;refs=4;consumer=100%
  17. family=reuse;policy=mat;refs=4;consumer=1%
  18. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  19. family=reuse;policy=notmat;refs=4;consumer=100%
  20. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  21. family=reuse;policy=default;refs=1;consumer=100%
  22. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  23. family=reuse;policy=mat;refs=1;consumer=100%
  24. family=reuse;policy=notmat;refs=1;consumer=100%
  25. family=reuse;policy=mat;refs=16;consumer=1%
  26. family=bag;stmt=join_where;threshold=180000
  27. family=bag;stmt=cte_joined_outer;threshold=2000
  28. family=reuse;policy=default;refs=16;consumer=1%
  29. family=reuse;policy=notmat;refs=4;consumer=1%
  30. family=bag;stmt=join_where;threshold=20

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
| family=bag;stmt=cte_joined_outer;threshold=180000 | j05_cte_joined_outer | valid | `d08ecce88af722b6` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20 | j05_cte_joined_outer | valid | `60d011356e039542` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=2000 | j05_cte_joined_outer | valid | `94774d6c2967f38d` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | j05m_cte_joined_outer_mat | valid | `4c14e0f134f802e7` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | j05m_cte_joined_outer_mat | valid | `d7af4e17205f4f22` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | j05m_cte_joined_outer_mat | valid | `62688ee4ef7cd1fa` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | j05n_cte_joined_outer_notmat | valid | `8de9f7bdd88bc2d2` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | j05n_cte_joined_outer_notmat | valid | `b687a0f411d7b2ab` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | j05n_cte_joined_outer_notmat | valid | `a0de9d5bcc30f638` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=join_where;threshold=180000 | j01_join_where | valid | `a542fca6125630a3` |  |
| family=bag;stmt=join_where;threshold=20 | j01_join_where | valid | `ba19056e4dc77f9c` |  |
| family=bag;stmt=join_where;threshold=2000 | j01_join_where | valid | `cab3a9fb46e687f5` |  |
| family=reuse;policy=default;refs=16;consumer=1% | r_default_n16_c1 | valid | `4cf84397665b0d35` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=16;consumer=100% | r_default_n16_c100 | valid | `107e20cb6bcfdbfb` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=1;consumer=1% | r_default_n1_c1 | valid | `9a64a0613768050f` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=1;consumer=100% | r_default_n1_c100 | valid | `8f90da054f9f2aca` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=4;consumer=1% | r_default_n4_c1 | valid | `f12d26b66a8684fe` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=100% | r_default_n4_c100 | valid | `31855c33f8c7881b` | policy=default; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=1% | r_mat_n16_c1 | valid | `9f092f705fac93bb` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=100% | r_mat_n16_c100 | valid | `64ec315ef647ef2d` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=1% | r_mat_n1_c1 | valid | `c1b22c44076b354f` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=100% | r_mat_n1_c100 | valid | `ff91c8711b7e3627` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=1% | r_mat_n4_c1 | valid | `d191d915dc30144c` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=100% | r_mat_n4_c100 | valid | `907ab5476f8ad1f5` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=notmat;refs=16;consumer=1% | r_notmat_n16_c1 | valid | `ebd81377094a97f2` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `543d906da363acbe` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `8029c01188eb7447` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `62f59d83157809d4` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `39f6874acd801497` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `25e554fe591f4ea9` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, DIST, DEBUG, FORMAT TEXT` supported on yugabyte
- one host: this run used node budget 2 CPU 3g and client budget 1 CPU 2g, so budget comparisons are separate classes here, not independent hardware
- one client process with one fixed maintenance writer at 10 updates/s; readers are sequential; the reported sample latency includes retry and backoff time
- cluster nodes are containers on the same host, so these are not network or independent-hardware measurements; container colocation is not data-colocation evidence
- readers round-robin across 3 endpoints; the writer's endpoint policy is recorded separately
- physical rank-index arm yb-range-rank; an engine whose planner ignores the index is a different physical design, not a faster one
- YugabyteDB DIST/DEBUG storage counters (rows and index rows scanned) are parsed per cell; RPC-level counters are not aggregated
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
