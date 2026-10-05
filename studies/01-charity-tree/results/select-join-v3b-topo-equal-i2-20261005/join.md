# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 30 cells, 30 valid, engine `yugabyte`, topology `single`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-v3b-topo-equal-i2-20261005`
- run tag: `run/study-01/select-join-v3b-topo-equal-i2-20261005`
- environment: `host-zenbook-ux5406sa`
- repo commit: `7243fee074bfb16fc5ad8b6afc1a2d7549080935`
- engine version: `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit`
- database engine image: `docker.io/yugabytedb/yugabyte:2025.2.6.0-b111` @ `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`
- benchmark client image: `localhost/charitybench-select:select-join-v3b-topo-equal-i2-20261005` @ `f3f1e10c44cbcb4c7b5bc56ed95a4f071682b6f7d33b5e0a7b688c4271eeef88`
- inputs digest: `322de6312f5b799c2b27d5179639658b98a484b4db65467e918077f3f25676c4`
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
- setting `phase`: isolated
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
| yb-single | 504 | 1 |

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
| 20 (0.01%) | 2.2810 | missing | missing | missing | missing | missing | 2.1751 | 85.4364 | 1.9179 | missing |
| 2000 (1%) | 9.9750 | missing | missing | missing | missing | missing | 8.9026 | 85.3347 | 9.3560 | missing |
| 180000 (90%) | 107.9493 | missing | missing | missing | missing | missing | 109.0072 | 102.5118 | 106.1922 | missing |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | cte_joined_outer vs join_where | 2.2810 | 2.1751 | 0.9536 | false |
| 20 | cte_joined_outer_mat vs join_where | 2.2810 | 85.4364 | 37.4552 | false |
| 20 | cte_joined_outer_notmat vs join_where | 2.2810 | 1.9179 | 0.8408 | false |
| 2000 | cte_joined_outer vs join_where | 9.9750 | 8.9026 | 0.8925 | false |
| 2000 | cte_joined_outer_mat vs join_where | 9.9750 | 85.3347 | 8.5548 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 9.9750 | 9.3560 | 0.9379 | false |
| 180000 | cte_joined_outer vs join_where | 107.9493 | 109.0072 | 1.0098 | false |
| 180000 | cte_joined_outer_mat vs join_where | 107.9493 | 102.5118 | 0.9496 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 107.9493 | 106.1922 | 0.9837 | false |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_joined_outer_mat vs cte_joined_outer | 2.1751 | 85.4364 | 39.2787 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 2.1751 | 1.9179 | 0.8818 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 8.9026 | 85.3347 | 9.5853 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 8.9026 | 9.3560 | 1.0509 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 109.0072 | 102.5118 | 0.9404 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 109.0072 | 106.1922 | 0.9742 | false | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 7.0131 | 1 | `1a38f469c0434969` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 8.3700 | 1 | `6ac25249c9094cf0` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | missing | | | | |
| default | 2 | 100% | missing | | | | |
| default | 4 | 1% | 8.3882 | 4 | `f7d577fd77c91d0c` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 8.8771 | 4 | `7b466a96c02d928d` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | missing | | | | |
| default | 8 | 100% | missing | | | | |
| default | 16 | 1% | 9.5951 | 16 | `4c76a5ad1fc03158` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 10.1698 | 16 | `de0699428514f971` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 8.6235 | 1 | `97cc5d36b7c4e9a2` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 8.7043 | 1 | `54a75d15f34e6297` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | missing | | | | |
| mat | 2 | 100% | missing | | | | |
| mat | 4 | 1% | 9.3029 | 4 | `8144d4cff2916541` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 9.2324 | 4 | `90b3ec0ca439d541` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | missing | | | | |
| mat | 8 | 100% | missing | | | | |
| mat | 16 | 1% | 9.3735 | 16 | `63f8d80e6218ec9a` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 10.4579 | 16 | `f3e217fddc1cc0dc` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 6.6781 | 1 | `534ce4e64ddc5930` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 8.0905 | 1 | `03f955d02ae3bfc6` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | missing | | | | |
| notmat | 2 | 100% | missing | | | | |
| notmat | 4 | 1% | 24.7592 | 4 | `57bfe0f1c735a320` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 31.5258 | 4 | `0da478614c0ce1e9` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | missing | | | | |
| notmat | 8 | 100% | missing | | | | |
| notmat | 16 | 1% | 81.2836 | 16 | `593452cd8055b931` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 83.8727 | 16 | `d0d99a2f9e144fcb` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 1.2296 | false |
| notmat | 1 | 1% | 0.9522 | false |
| mat | 1 | 100% | 1.0399 | false |
| notmat | 1 | 100% | 0.9666 | false |
| mat | 4 | 1% | 1.1090 | false |
| notmat | 4 | 1% | 2.9517 | false |
| mat | 4 | 100% | 1.0400 | false |
| notmat | 4 | 100% | 3.5514 | false |
| mat | 16 | 1% | 0.9769 | false |
| notmat | 16 | 1% | 8.4713 | false |
| mat | 16 | 100% | 1.0283 | false |
| notmat | 16 | 100% | 8.2472 | false |

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

- client: nr_periods 284, nr_throttled 117, throttled_usec 3892542 | quota 100000 100000
, memory 2147483648

- yb-single: nr_periods 358, nr_throttled 0, throttled_usec 0 | quota 600000 100000, memory 9663676416

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
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 102.5824 | 109.0072 | 110.4081 | 5 | 9.2 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 2.0482 | 2.1751 | 2.6259 | 5 | 459.7 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 8.5530 | 8.9026 | 10.6159 | 5 | 112.3 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 97.1445 | 102.5118 | 108.5881 | 5 | 9.8 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 84.1804 | 85.4364 | 97.0015 | 5 | 11.7 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 78.9961 | 85.3347 | 100.4581 | 5 | 11.7 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 100.4095 | 106.1922 | 118.6023 | 5 | 9.4 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 1.7170 | 1.9179 | 2.4139 | 5 | 521.4 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 8.9757 | 9.3560 | 9.6093 | 5 | 106.9 |
| family=bag;stmt=join_where;threshold=180000 | valid | 102.3429 | 107.9493 | 110.4221 | 5 | 9.3 |
| family=bag;stmt=join_where;threshold=20 | valid | 2.1263 | 2.2810 | 2.7276 | 5 | 438.4 |
| family=bag;stmt=join_where;threshold=2000 | valid | 8.5922 | 9.9750 | 14.4760 | 5 | 100.3 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 8.9001 | 9.5951 | 11.7535 | 5 | 104.2 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 9.4603 | 10.1698 | 10.3447 | 5 | 98.3 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 6.5569 | 7.0131 | 8.0128 | 5 | 142.6 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 7.8981 | 8.3700 | 8.7717 | 5 | 119.5 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 7.9477 | 8.3882 | 8.6440 | 5 | 119.2 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 7.9715 | 8.8771 | 10.1833 | 5 | 112.6 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 9.0371 | 9.3735 | 9.9450 | 5 | 106.7 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 9.7946 | 10.4579 | 10.8964 | 5 | 95.6 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 7.4999 | 8.6235 | 9.5832 | 5 | 116.0 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 8.2586 | 8.7043 | 10.4347 | 5 | 114.9 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 8.7604 | 9.3029 | 9.4209 | 5 | 107.5 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 8.2957 | 9.2324 | 9.9065 | 5 | 108.3 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 70.8277 | 81.2836 | 83.8718 | 5 | 12.3 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 82.9922 | 83.8727 | 95.4241 | 5 | 11.9 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 6.2754 | 6.6781 | 6.9665 | 5 | 149.7 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 7.9555 | 8.0905 | 8.3166 | 5 | 123.6 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 23.2482 | 24.7592 | 25.5707 | 5 | 40.4 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 29.0141 | 31.5258 | 42.0294 | 5 | 31.7 |

## Matched blended workload

- phase: `isolated`; target offered rate: 10 updates/s
- offered 0, committed 0, failed 0 over 0.000 s
- offered rate 0.000/s, completed rate 0.000/s
- per key: no keys touched
- ledger: 0 keys reconciled, 0 mismatch

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42274093`, 5 interleaved rounds
  1. family=reuse;policy=mat;refs=1;consumer=1%
  2. family=reuse;policy=default;refs=1;consumer=1%
  3. family=reuse;policy=default;refs=16;consumer=100%
  4. family=reuse;policy=mat;refs=16;consumer=100%
  5. family=reuse;policy=mat;refs=16;consumer=1%
  6. family=reuse;policy=notmat;refs=16;consumer=1%
  7. family=reuse;policy=notmat;refs=16;consumer=100%
  8. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  9. family=reuse;policy=default;refs=16;consumer=1%
  10. family=bag;stmt=join_where;threshold=2000
  11. family=reuse;policy=mat;refs=1;consumer=100%
  12. family=reuse;policy=notmat;refs=4;consumer=100%
  13. family=bag;stmt=cte_joined_outer;threshold=2000
  14. family=reuse;policy=default;refs=4;consumer=100%
  15. family=reuse;policy=mat;refs=4;consumer=100%
  16. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  17. family=reuse;policy=mat;refs=4;consumer=1%
  18. family=reuse;policy=notmat;refs=1;consumer=100%
  19. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  20. family=bag;stmt=cte_joined_outer_mat;threshold=20
  21. family=reuse;policy=notmat;refs=4;consumer=1%
  22. family=reuse;policy=notmat;refs=1;consumer=1%
  23. family=reuse;policy=default;refs=1;consumer=100%
  24. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  25. family=bag;stmt=join_where;threshold=180000
  26. family=bag;stmt=cte_joined_outer;threshold=20
  27. family=bag;stmt=cte_joined_outer;threshold=180000
  28. family=bag;stmt=join_where;threshold=20
  29. family=reuse;policy=default;refs=4;consumer=1%
  30. family=bag;stmt=cte_joined_outer_notmat;threshold=20

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
| family=bag;stmt=cte_joined_outer;threshold=180000 | j05_cte_joined_outer | valid | `361e24d0786f3586` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20 | j05_cte_joined_outer | valid | `0a4a2145c6af23bb` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=2000 | j05_cte_joined_outer | valid | `94ee73cf02bb9e31` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | j05m_cte_joined_outer_mat | valid | `7740a34ea3708421` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | j05m_cte_joined_outer_mat | valid | `c599ad7cc731e751` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | j05m_cte_joined_outer_mat | valid | `c0a7a1d352634333` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | j05n_cte_joined_outer_notmat | valid | `1124a2baf4cfe663` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | j05n_cte_joined_outer_notmat | valid | `1112c1ab966dad8b` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | j05n_cte_joined_outer_notmat | valid | `e7eba9ce41a19d00` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=join_where;threshold=180000 | j01_join_where | valid | `4a02da675d1461df` |  |
| family=bag;stmt=join_where;threshold=20 | j01_join_where | valid | `14e5e9b018ac8460` |  |
| family=bag;stmt=join_where;threshold=2000 | j01_join_where | valid | `6f4a8664f0627c47` |  |
| family=reuse;policy=default;refs=16;consumer=1% | r_default_n16_c1 | valid | `4c76a5ad1fc03158` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=16;consumer=100% | r_default_n16_c100 | valid | `de0699428514f971` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=1;consumer=1% | r_default_n1_c1 | valid | `1a38f469c0434969` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=1;consumer=100% | r_default_n1_c100 | valid | `6ac25249c9094cf0` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=4;consumer=1% | r_default_n4_c1 | valid | `f7d577fd77c91d0c` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=100% | r_default_n4_c100 | valid | `7b466a96c02d928d` | policy=default; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=1% | r_mat_n16_c1 | valid | `63f8d80e6218ec9a` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=100% | r_mat_n16_c100 | valid | `f3e217fddc1cc0dc` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=1% | r_mat_n1_c1 | valid | `97cc5d36b7c4e9a2` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=100% | r_mat_n1_c100 | valid | `54a75d15f34e6297` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=1% | r_mat_n4_c1 | valid | `8144d4cff2916541` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=100% | r_mat_n4_c100 | valid | `90b3ec0ca439d541` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=notmat;refs=16;consumer=1% | r_notmat_n16_c1 | valid | `593452cd8055b931` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `d0d99a2f9e144fcb` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `534ce4e64ddc5930` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `03f955d02ae3bfc6` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `57bfe0f1c735a320` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `0da478614c0ce1e9` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, DIST, DEBUG, FORMAT TEXT` supported on yugabyte
- one host: this run used node budget 6 CPU 9g and client budget 1 CPU 2g, so budget comparisons are separate classes here, not independent hardware
- isolated phase: no concurrent writers in this run; the blended read/write phase is a separate labelled run
- physical rank-index arm yb-range-rank; an engine whose planner ignores the index is a different physical design, not a faster one
- YugabyteDB DIST/DEBUG storage counters (rows and index rows scanned) are parsed per cell; RPC-level counters are not aggregated
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
