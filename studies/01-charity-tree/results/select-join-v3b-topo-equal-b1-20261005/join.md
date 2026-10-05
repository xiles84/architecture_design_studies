# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 30 cells, 30 valid, engine `yugabyte`, topology `single`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-v3b-topo-equal-b1-20261005`
- run tag: `run/study-01/select-join-v3b-topo-equal-b1-20261005`
- environment: `host-zenbook-ux5406sa`
- repo commit: `7243fee074bfb16fc5ad8b6afc1a2d7549080935`
- engine version: `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit`
- database engine image: `docker.io/yugabytedb/yugabyte:2025.2.6.0-b111` @ `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`
- benchmark client image: `localhost/charitybench-select:select-join-v3b-topo-equal-b1-20261005` @ `f3f1e10c44cbcb4c7b5bc56ed95a4f071682b6f7d33b5e0a7b688c4271eeef88`
- inputs digest: `6ba87ba204d27111ad77126a7daddb0de600ac59b4ee7279c5f634849672c2f9`
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
| yb-single | 505 | 2 |

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
| 20 (0.01%) | 2.0736 | missing | missing | missing | missing | missing | 1.8665 | 93.8055 | 2.2151 | missing |
| 2000 (1%) | 9.6234 | missing | missing | missing | missing | missing | 10.0710 | 88.3222 | 10.7140 | missing |
| 180000 (90%) | 119.1081 | missing | missing | missing | missing | missing | 118.5620 | 111.3526 | 115.5723 | missing |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | cte_joined_outer vs join_where | 2.0736 | 1.8665 | 0.9001 | false |
| 20 | cte_joined_outer_mat vs join_where | 2.0736 | 93.8055 | 45.2382 | false |
| 20 | cte_joined_outer_notmat vs join_where | 2.0736 | 2.2151 | 1.0682 | false |
| 2000 | cte_joined_outer vs join_where | 9.6234 | 10.0710 | 1.0465 | false |
| 2000 | cte_joined_outer_mat vs join_where | 9.6234 | 88.3222 | 9.1778 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 9.6234 | 10.7140 | 1.1133 | false |
| 180000 | cte_joined_outer vs join_where | 119.1081 | 118.5620 | 0.9954 | false |
| 180000 | cte_joined_outer_mat vs join_where | 119.1081 | 111.3526 | 0.9349 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 119.1081 | 115.5723 | 0.9703 | false |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_joined_outer_mat vs cte_joined_outer | 1.8665 | 93.8055 | 50.2584 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 1.8665 | 2.2151 | 1.1868 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 10.0710 | 88.3222 | 8.7700 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 10.0710 | 10.7140 | 1.0638 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 118.5620 | 111.3526 | 0.9392 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 118.5620 | 115.5723 | 0.9748 | false | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 7.8273 | 1 | `de53a79fec263782` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 8.6493 | 1 | `c631b1abc93dbd5a` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | missing | | | | |
| default | 2 | 100% | missing | | | | |
| default | 4 | 1% | 8.9034 | 4 | `aac5e9a8839ff8de` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 9.3293 | 4 | `83c30e955456d7f8` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | missing | | | | |
| default | 8 | 100% | missing | | | | |
| default | 16 | 1% | 10.2419 | 16 | `6e1c6bf161738607` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 10.6998 | 16 | `600d7711480b60c7` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 9.8492 | 1 | `e52f4a00fa8a7cf2` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 8.3108 | 1 | `5ce7f173123f67d4` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | missing | | | | |
| mat | 2 | 100% | missing | | | | |
| mat | 4 | 1% | 9.3678 | 4 | `ec3a5640b6b3f1c2` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 9.1888 | 4 | `4f445b8150b13ccc` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | missing | | | | |
| mat | 8 | 100% | missing | | | | |
| mat | 16 | 1% | 9.6948 | 16 | `0ccd1bead3009ffd` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 10.3947 | 16 | `f1eeff86d8a13462` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 6.9936 | 1 | `dbdc783d51d6e979` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 8.5113 | 1 | `cfb1906539a7b7f6` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | missing | | | | |
| notmat | 2 | 100% | missing | | | | |
| notmat | 4 | 1% | 24.8511 | 4 | `014df0072375c08d` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 33.2752 | 4 | `aeb1715cb9a4e2df` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | missing | | | | |
| notmat | 8 | 100% | missing | | | | |
| notmat | 16 | 1% | 81.4039 | 16 | `0fc5c1bac0a5f563` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 91.4780 | 16 | `d8fed7c6c64a9119` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 1.2583 | false |
| notmat | 1 | 1% | 0.8935 | false |
| mat | 1 | 100% | 0.9609 | false |
| notmat | 1 | 100% | 0.9841 | false |
| mat | 4 | 1% | 1.0522 | false |
| notmat | 4 | 1% | 2.7912 | false |
| mat | 4 | 100% | 0.9849 | false |
| notmat | 4 | 100% | 3.5667 | false |
| mat | 16 | 1% | 0.9466 | false |
| notmat | 16 | 1% | 7.9482 | false |
| mat | 16 | 100% | 0.9715 | false |
| notmat | 16 | 100% | 8.5495 | false |

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

- client: nr_periods 315, nr_throttled 136, throttled_usec 5587070 | quota 100000 100000
, memory 2147483648

- yb-single: nr_periods 383, nr_throttled 0, throttled_usec 0 | quota 600000 100000, memory 9663676416

## Access path per cell

_Engine-reported from each cell's own plan: whether the rank index or a sequential scan was used, rows scanned, PostgreSQL buffers and filter removals, and reader retries of transient conflicts. A cell that scanned the whole table at a selective threshold is a different physical regime, not a faster query._

| axis | status | index scan | seq scan | storage rows scanned | index rows scanned | buffers hit/read | rows removed | retries |
|---|---|---|---|---|---|---|---|---|
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | false | true | 202050 | 0 | 0/0 | 0 | 0 |
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
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 114.1029 | 118.5620 | 125.4959 | 5 | 8.4 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 1.8129 | 1.8665 | 2.1149 | 5 | 535.8 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 9.2682 | 10.0710 | 11.7552 | 5 | 99.3 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 106.2076 | 111.3526 | 196.8052 | 5 | 9.0 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 89.5245 | 93.8055 | 108.7289 | 5 | 10.7 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 84.8806 | 88.3222 | 90.2927 | 5 | 11.3 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 111.9749 | 115.5723 | 123.5693 | 5 | 8.7 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 1.9247 | 2.2151 | 3.0936 | 5 | 451.5 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 10.3821 | 10.7140 | 11.4297 | 5 | 93.3 |
| family=bag;stmt=join_where;threshold=180000 | valid | 117.1712 | 119.1081 | 127.4775 | 5 | 8.4 |
| family=bag;stmt=join_where;threshold=20 | valid | 1.9545 | 2.0736 | 2.3392 | 5 | 482.3 |
| family=bag;stmt=join_where;threshold=2000 | valid | 8.7648 | 9.6234 | 10.6479 | 5 | 103.9 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 9.8370 | 10.2419 | 10.4563 | 5 | 97.6 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 10.4119 | 10.6998 | 11.2853 | 5 | 93.5 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 7.7153 | 7.8273 | 8.1988 | 5 | 127.8 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 8.5115 | 8.6493 | 10.1891 | 5 | 115.6 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 8.2064 | 8.9034 | 9.4876 | 5 | 112.3 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 8.7234 | 9.3293 | 9.8209 | 5 | 107.2 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 9.1535 | 9.6948 | 10.2524 | 5 | 103.1 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 9.8751 | 10.3947 | 12.4595 | 5 | 96.2 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 9.5152 | 9.8492 | 10.4154 | 5 | 101.5 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 8.2104 | 8.3108 | 9.0705 | 5 | 120.3 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 8.4200 | 9.3678 | 28.9601 | 5 | 106.7 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 8.5514 | 9.1888 | 9.6042 | 5 | 108.8 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 75.4402 | 81.4039 | 108.8201 | 5 | 12.3 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 87.6783 | 91.4780 | 98.5966 | 5 | 10.9 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 6.5459 | 6.9936 | 8.8762 | 5 | 143.0 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 8.2425 | 8.5113 | 9.4304 | 5 | 117.5 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 23.0866 | 24.8511 | 42.0319 | 5 | 40.2 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 29.7644 | 33.2752 | 34.6984 | 5 | 30.1 |

## Matched blended workload

- phase: `blended`; target offered rate: 10 updates/s
- offered 156, committed 156, failed 0 over 15.549 s
- offered rate 10.033/s, completed rate 10.033/s
- writer operations per endpoint: [156]
- reader operations overlapping a committed write: 68 of 211
- reader retries of transient conflicts: 1
- per key: 156 keys, min 1, max 1 per key
- ledger: 156 keys reconciled, 0 mismatch
- update latency: avg 3.435 ms, min 1.770, max 71.066 over 156

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42107455`, 5 interleaved rounds
  1. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  2. family=reuse;policy=notmat;refs=4;consumer=100%
  3. family=reuse;policy=notmat;refs=1;consumer=1%
  4. family=bag;stmt=cte_joined_outer;threshold=20
  5. family=reuse;policy=default;refs=16;consumer=1%
  6. family=reuse;policy=default;refs=4;consumer=1%
  7. family=reuse;policy=mat;refs=4;consumer=100%
  8. family=reuse;policy=mat;refs=16;consumer=1%
  9. family=reuse;policy=default;refs=1;consumer=100%
  10. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  11. family=bag;stmt=join_where;threshold=20
  12. family=reuse;policy=default;refs=16;consumer=100%
  13. family=reuse;policy=mat;refs=1;consumer=100%
  14. family=reuse;policy=notmat;refs=16;consumer=1%
  15. family=reuse;policy=default;refs=1;consumer=1%
  16. family=bag;stmt=cte_joined_outer_mat;threshold=20
  17. family=bag;stmt=cte_joined_outer;threshold=180000
  18. family=reuse;policy=notmat;refs=4;consumer=1%
  19. family=bag;stmt=join_where;threshold=2000
  20. family=reuse;policy=notmat;refs=1;consumer=100%
  21. family=reuse;policy=notmat;refs=16;consumer=100%
  22. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  23. family=bag;stmt=join_where;threshold=180000
  24. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  25. family=reuse;policy=mat;refs=1;consumer=1%
  26. family=reuse;policy=mat;refs=4;consumer=1%
  27. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  28. family=reuse;policy=default;refs=4;consumer=100%
  29. family=reuse;policy=mat;refs=16;consumer=100%
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
| family=bag;stmt=cte_joined_outer;threshold=180000 | j05_cte_joined_outer | valid | `b258538f17f9fea4` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20 | j05_cte_joined_outer | valid | `b6913c06b8fc4abd` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=2000 | j05_cte_joined_outer | valid | `8fadd6f6ca26b9a0` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | j05m_cte_joined_outer_mat | valid | `f616b42366f65dff` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | j05m_cte_joined_outer_mat | valid | `4e8f768fb4797bb5` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | j05m_cte_joined_outer_mat | valid | `b7f308c7ef298fa3` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | j05n_cte_joined_outer_notmat | valid | `8564342e80079155` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | j05n_cte_joined_outer_notmat | valid | `ef3dd067c8a91050` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | j05n_cte_joined_outer_notmat | valid | `26cd3734c12711c0` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=join_where;threshold=180000 | j01_join_where | valid | `fa33152803b5a79a` |  |
| family=bag;stmt=join_where;threshold=20 | j01_join_where | valid | `90722a8a4be78434` |  |
| family=bag;stmt=join_where;threshold=2000 | j01_join_where | valid | `acc7f0d20d5174d5` |  |
| family=reuse;policy=default;refs=16;consumer=1% | r_default_n16_c1 | valid | `6e1c6bf161738607` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=16;consumer=100% | r_default_n16_c100 | valid | `600d7711480b60c7` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=1;consumer=1% | r_default_n1_c1 | valid | `de53a79fec263782` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=1;consumer=100% | r_default_n1_c100 | valid | `c631b1abc93dbd5a` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=4;consumer=1% | r_default_n4_c1 | valid | `aac5e9a8839ff8de` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=100% | r_default_n4_c100 | valid | `83c30e955456d7f8` | policy=default; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=1% | r_mat_n16_c1 | valid | `0ccd1bead3009ffd` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=100% | r_mat_n16_c100 | valid | `f1eeff86d8a13462` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=1% | r_mat_n1_c1 | valid | `e52f4a00fa8a7cf2` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=100% | r_mat_n1_c100 | valid | `5ce7f173123f67d4` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=1% | r_mat_n4_c1 | valid | `ec3a5640b6b3f1c2` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=100% | r_mat_n4_c100 | valid | `4f445b8150b13ccc` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=notmat;refs=16;consumer=1% | r_notmat_n16_c1 | valid | `0fc5c1bac0a5f563` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `d8fed7c6c64a9119` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `dbdc783d51d6e979` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `cfb1906539a7b7f6` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `014df0072375c08d` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `aeb1715cb9a4e2df` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, DIST, DEBUG, FORMAT TEXT` supported on yugabyte
- one host: this run used node budget 6 CPU 9g and client budget 1 CPU 2g, so budget comparisons are separate classes here, not independent hardware
- one client process with one fixed maintenance writer at 10 updates/s; readers are sequential; the reported sample latency includes retry and backoff time
- physical rank-index arm yb-range-rank; an engine whose planner ignores the index is a different physical design, not a faster one
- YugabyteDB DIST/DEBUG storage counters (rows and index rows scanned) are parsed per cell; RPC-level counters are not aggregated
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
