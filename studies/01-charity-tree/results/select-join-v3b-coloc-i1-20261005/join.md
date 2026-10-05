# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 30 cells, 30 valid, engine `yugabyte`, topology `cluster3`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-v3b-coloc-i1-20261005`
- run tag: `run/study-01/select-join-v3b-coloc-i1-20261005`
- environment: `host-zenbook-ux5406sa`
- repo commit: `7243fee074bfb16fc5ad8b6afc1a2d7549080935`
- engine version: `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit`
- database engine image: `docker.io/yugabytedb/yugabyte:2025.2.6.0-b111` @ `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`
- benchmark client image: `localhost/charitybench-select:select-join-v3b-coloc-i1-20261005` @ `f3f1e10c44cbcb4c7b5bc56ed95a4f071682b6f7d33b5e0a7b688c4271eeef88`
- inputs digest: `933d5c83aaff4896b21f69cbb4b42b08dd9709b49a80e429b5966e1199d3766b`
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
| 20 (0.01%) | 2.4578 | missing | missing | missing | missing | missing | 1.8661 | 111.9835 | 2.4683 | missing |
| 2000 (1%) | 7.0962 | missing | missing | missing | missing | missing | 6.5196 | 112.1343 | 7.0956 | missing |
| 180000 (90%) | 135.1399 | missing | missing | missing | missing | missing | 131.2419 | 99.3995 | 134.7846 | missing |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | cte_joined_outer vs join_where | 2.4578 | 1.8661 | 0.7593 | false |
| 20 | cte_joined_outer_mat vs join_where | 2.4578 | 111.9835 | 45.5629 | false |
| 20 | cte_joined_outer_notmat vs join_where | 2.4578 | 2.4683 | 1.0043 | false |
| 2000 | cte_joined_outer vs join_where | 7.0962 | 6.5196 | 0.9187 | false |
| 2000 | cte_joined_outer_mat vs join_where | 7.0962 | 112.1343 | 15.8020 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 7.0962 | 7.0956 | 0.9999 | false |
| 180000 | cte_joined_outer vs join_where | 135.1399 | 131.2419 | 0.9712 | false |
| 180000 | cte_joined_outer_mat vs join_where | 135.1399 | 99.3995 | 0.7355 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 135.1399 | 134.7846 | 0.9974 | false |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_joined_outer_mat vs cte_joined_outer | 1.8661 | 111.9835 | 60.0079 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 1.8661 | 2.4683 | 1.3227 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 6.5196 | 112.1343 | 17.1996 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 6.5196 | 7.0956 | 1.0884 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 131.2419 | 99.3995 | 0.7574 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 131.2419 | 134.7846 | 1.0270 | false | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 1.9257 | 1 | `788b1b6d8a5ce7da` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 6.7059 | 1 | `ee0e486764294187` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | missing | | | | |
| default | 2 | 100% | missing | | | | |
| default | 4 | 1% | 6.3128 | 4 | `fa369cfd15f1a0e8` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 6.2539 | 4 | `d9195713485c8f72` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | missing | | | | |
| default | 8 | 100% | missing | | | | |
| default | 16 | 1% | 7.6562 | 16 | `aabdc202447abad1` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 7.6698 | 16 | `5a02a7390f5bb508` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 6.3418 | 1 | `09b38315e261e412` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 6.1495 | 1 | `0df1294ac5f0ea3e` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | missing | | | | |
| mat | 2 | 100% | missing | | | | |
| mat | 4 | 1% | 6.7352 | 4 | `64cc10fa28e08526` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 6.4506 | 4 | `6904dcb97810e7de` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | missing | | | | |
| mat | 8 | 100% | missing | | | | |
| mat | 16 | 1% | 7.9898 | 16 | `eccbd5b4d69cdffe` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 8.0604 | 16 | `cf97c1ae031a9f32` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 1.6589 | 1 | `130cbee7c13685f8` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 6.4259 | 1 | `d7aaa002d32b64ce` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | missing | | | | |
| notmat | 2 | 100% | missing | | | | |
| notmat | 4 | 1% | 5.0546 | 4 | `c77441084eae57dd` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 21.2756 | 4 | `32d43288fd343334` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | missing | | | | |
| notmat | 8 | 100% | missing | | | | |
| notmat | 16 | 1% | 12.9642 | 16 | `9ac73d1a2c0b55e3` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 82.3852 | 16 | `d97a98dd18f11059` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 3.2932 | false |
| notmat | 1 | 1% | 0.8614 | false |
| mat | 1 | 100% | 0.9170 | false |
| notmat | 1 | 100% | 0.9582 | false |
| mat | 4 | 1% | 1.0669 | false |
| notmat | 4 | 1% | 0.8007 | false |
| mat | 4 | 100% | 1.0315 | false |
| notmat | 4 | 100% | 3.4020 | false |
| mat | 16 | 1% | 1.0436 | false |
| notmat | 16 | 1% | 1.6933 | false |
| mat | 16 | 100% | 1.0509 | false |
| notmat | 16 | 100% | 10.7415 | false |

## Physical arm

- `charity_pkey`: `CREATE UNIQUE INDEX charity_pkey ON selectjoin_domain.charity USING lsm (charity_id ASC)`
- `donation_person_idx`: `CREATE INDEX donation_person_idx ON selectjoin_domain.donation USING lsm (person_id ASC)`
- `donation_pkey`: `CREATE UNIQUE INDEX donation_pkey ON selectjoin_domain.donation USING lsm (donation_id ASC)`
- `donation_rank_idx`: `CREATE INDEX donation_rank_idx ON selectjoin_domain.donation USING lsm (rank ASC)`
- `person_pkey`: `CREATE UNIQUE INDEX person_pkey ON selectjoin_domain.person USING lsm (person_id ASC)`

- index arm `yb-range-rank`; placement `colocated`; node budget `2 CPU 3g`; client budget `1 CPU 2g`

| table | tablets | hash key columns | colocated | tablegroup | colocation id |
|---|---|---|---|---|---|
| selectjoin_domain.charity | 1 | 0 | true | 16388 | 1890171544 |
| selectjoin_domain.person | 1 | 0 | true | 16388 | 4051657324 |
| selectjoin_domain.donation | 1 | 0 | true | 16388 | 1300651145 |
| selectjoin_domain.donation_pkey | 0 | 0 | false | unknown | unknown |
| selectjoin_domain.donation_rank_idx | 1 | 0 | true | 16388 | 2570119830 |
| selectjoin_domain.donation_person_idx | 1 | 0 | true | 16388 | 1383227021 |

Physical-arm gate (default join must use the rank index at its three most selective thresholds):

- threshold 20: used
- threshold 2000: used

## Throttling (cgroup `cpu.stat` deltas)

- client: nr_periods 278, nr_throttled 95, throttled_usec 4547990 | quota 100000 100000
, memory 2147483648

- yb-n1: nr_periods 535, nr_throttled 9, throttled_usec 827369 | quota 200000 100000, memory 3221225472
- yb-n2: nr_periods 526, nr_throttled 5, throttled_usec 684797 | quota 200000 100000, memory 3221225472
- yb-n3: nr_periods 520, nr_throttled 9, throttled_usec 388943 | quota 200000 100000, memory 3221225472

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
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 129.7357 | 131.2419 | 139.9777 | 5 | 7.6 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 1.6179 | 1.8661 | 1.9332 | 5 | 535.9 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 6.4020 | 6.5196 | 7.2506 | 5 | 153.4 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 97.1453 | 99.3995 | 103.5379 | 5 | 10.1 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 108.2512 | 111.9835 | 118.3680 | 5 | 8.9 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 103.9161 | 112.1343 | 118.7371 | 5 | 8.9 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 130.7176 | 134.7846 | 137.0664 | 5 | 7.4 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 2.3665 | 2.4683 | 2.8354 | 5 | 405.1 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 6.9774 | 7.0956 | 7.4120 | 5 | 140.9 |
| family=bag;stmt=join_where;threshold=180000 | valid | 131.6713 | 135.1399 | 146.1459 | 5 | 7.4 |
| family=bag;stmt=join_where;threshold=20 | valid | 2.3690 | 2.4578 | 2.6875 | 5 | 406.9 |
| family=bag;stmt=join_where;threshold=2000 | valid | 6.6634 | 7.0962 | 8.2111 | 5 | 140.9 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 7.4054 | 7.6562 | 7.8807 | 5 | 130.6 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 7.3934 | 7.6698 | 7.9537 | 5 | 130.4 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 1.8039 | 1.9257 | 2.0632 | 5 | 519.3 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 6.4181 | 6.7059 | 6.8460 | 5 | 149.1 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 6.2043 | 6.3128 | 6.3931 | 5 | 158.4 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 5.8933 | 6.2539 | 6.5473 | 5 | 159.9 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 7.7640 | 7.9898 | 8.6295 | 5 | 125.2 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 7.8918 | 8.0604 | 8.2151 | 5 | 124.1 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 6.0695 | 6.3418 | 6.9205 | 5 | 157.7 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 6.0241 | 6.1495 | 6.5946 | 5 | 162.6 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 6.3742 | 6.7352 | 6.8619 | 5 | 148.5 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 6.3390 | 6.4506 | 6.8490 | 5 | 155.0 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 12.7300 | 12.9642 | 13.4284 | 5 | 77.1 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 73.1484 | 82.3852 | 94.1720 | 5 | 12.1 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 1.5435 | 1.6589 | 2.0003 | 5 | 602.8 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 6.1739 | 6.4259 | 6.6261 | 5 | 155.6 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 4.9372 | 5.0546 | 5.6270 | 5 | 197.8 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 20.1209 | 21.2756 | 23.2701 | 5 | 47.0 |

## Matched blended workload

- phase: `isolated`; target offered rate: 10 updates/s
- offered 0, committed 0, failed 0 over 0.000 s
- offered rate 0.000/s, completed rate 0.000/s
- per key: no keys touched
- ledger: 0 keys reconciled, 0 mismatch

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42428777`, 5 interleaved rounds
  1. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  2. family=reuse;policy=mat;refs=16;consumer=1%
  3. family=reuse;policy=notmat;refs=4;consumer=100%
  4. family=reuse;policy=notmat;refs=1;consumer=100%
  5. family=reuse;policy=mat;refs=1;consumer=1%
  6. family=reuse;policy=notmat;refs=16;consumer=100%
  7. family=reuse;policy=default;refs=1;consumer=1%
  8. family=bag;stmt=join_where;threshold=20
  9. family=bag;stmt=cte_joined_outer;threshold=2000
  10. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  11. family=reuse;policy=default;refs=1;consumer=100%
  12. family=reuse;policy=default;refs=4;consumer=100%
  13. family=reuse;policy=notmat;refs=1;consumer=1%
  14. family=reuse;policy=mat;refs=4;consumer=1%
  15. family=bag;stmt=cte_joined_outer;threshold=20
  16. family=reuse;policy=default;refs=16;consumer=1%
  17. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  18. family=reuse;policy=mat;refs=16;consumer=100%
  19. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  20. family=bag;stmt=join_where;threshold=2000
  21. family=reuse;policy=default;refs=16;consumer=100%
  22. family=bag;stmt=join_where;threshold=180000
  23. family=reuse;policy=notmat;refs=4;consumer=1%
  24. family=reuse;policy=default;refs=4;consumer=1%
  25. family=reuse;policy=mat;refs=4;consumer=100%
  26. family=reuse;policy=mat;refs=1;consumer=100%
  27. family=reuse;policy=notmat;refs=16;consumer=1%
  28. family=bag;stmt=cte_joined_outer_mat;threshold=20
  29. family=bag;stmt=cte_joined_outer;threshold=180000
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
| family=bag;stmt=cte_joined_outer;threshold=180000 | j05_cte_joined_outer | valid | `500de8fcc0030611` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20 | j05_cte_joined_outer | valid | `2838e42e1b15c936` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=2000 | j05_cte_joined_outer | valid | `fa2000b4635f5a1f` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | j05m_cte_joined_outer_mat | valid | `dd6e04821ec74748` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | j05m_cte_joined_outer_mat | valid | `2e8bf96922493fd8` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | j05m_cte_joined_outer_mat | valid | `0c427cb337e00159` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | j05n_cte_joined_outer_notmat | valid | `65b9cc0532728685` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | j05n_cte_joined_outer_notmat | valid | `b6813d74256ad9d9` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | j05n_cte_joined_outer_notmat | valid | `97376e2cf1db51c3` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=join_where;threshold=180000 | j01_join_where | valid | `7cc4bff13189079d` |  |
| family=bag;stmt=join_where;threshold=20 | j01_join_where | valid | `1c82357bb3be47b3` |  |
| family=bag;stmt=join_where;threshold=2000 | j01_join_where | valid | `e7648363636c407c` |  |
| family=reuse;policy=default;refs=16;consumer=1% | r_default_n16_c1 | valid | `aabdc202447abad1` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=16;consumer=100% | r_default_n16_c100 | valid | `5a02a7390f5bb508` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=1;consumer=1% | r_default_n1_c1 | valid | `788b1b6d8a5ce7da` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=1;consumer=100% | r_default_n1_c100 | valid | `ee0e486764294187` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=4;consumer=1% | r_default_n4_c1 | valid | `fa369cfd15f1a0e8` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=100% | r_default_n4_c100 | valid | `d9195713485c8f72` | policy=default; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=1% | r_mat_n16_c1 | valid | `eccbd5b4d69cdffe` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=100% | r_mat_n16_c100 | valid | `cf97c1ae031a9f32` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=1% | r_mat_n1_c1 | valid | `09b38315e261e412` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=100% | r_mat_n1_c100 | valid | `0df1294ac5f0ea3e` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=1% | r_mat_n4_c1 | valid | `64cc10fa28e08526` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=100% | r_mat_n4_c100 | valid | `6904dcb97810e7de` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=notmat;refs=16;consumer=1% | r_notmat_n16_c1 | valid | `9ac73d1a2c0b55e3` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `d97a98dd18f11059` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `130cbee7c13685f8` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `d7aaa002d32b64ce` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `c77441084eae57dd` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `32d43288fd343334` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, DIST, DEBUG, FORMAT TEXT` supported on yugabyte
- one host: this run used node budget 2 CPU 3g and client budget 1 CPU 2g, so budget comparisons are separate classes here, not independent hardware
- isolated phase: no concurrent writers in this run; the blended read/write phase is a separate labelled run
- cluster nodes are containers on the same host, so these are not network or independent-hardware measurements; container colocation is not data-colocation evidence
- readers round-robin across 3 endpoints; the writer's endpoint policy is recorded separately
- physical rank-index arm yb-range-rank; an engine whose planner ignores the index is a different physical design, not a faster one
- YugabyteDB DIST/DEBUG storage counters (rows and index rows scanned) are parsed per cell; RPC-level counters are not aggregated
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
