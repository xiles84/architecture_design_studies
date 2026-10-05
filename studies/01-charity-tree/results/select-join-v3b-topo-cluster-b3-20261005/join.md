# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 30 cells, 30 valid, engine `yugabyte`, topology `cluster3`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-v3b-topo-cluster-b3-20261005`
- run tag: `run/study-01/select-join-v3b-topo-cluster-b3-20261005`
- environment: `host-zenbook-ux5406sa`
- repo commit: `7243fee074bfb16fc5ad8b6afc1a2d7549080935`
- engine version: `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit`
- database engine image: `docker.io/yugabytedb/yugabyte:2025.2.6.0-b111` @ `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`
- benchmark client image: `localhost/charitybench-select:select-join-v3b-topo-cluster-b3-20261005` @ `f3f1e10c44cbcb4c7b5bc56ed95a4f071682b6f7d33b5e0a7b688c4271eeef88`
- inputs digest: `9ce338e8f28ebf95c57c51682b888e3ae6e3df49d58c24201804bb1d730d3c3a`
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
| 20 (0.01%) | 2.9576 | missing | missing | missing | missing | missing | 2.3839 | 138.5601 | 2.1410 | missing |
| 2000 (1%) | 10.1490 | missing | missing | missing | missing | missing | 11.7377 | 99.0413 | 10.0796 | missing |
| 180000 (90%) | 160.0020 | missing | missing | missing | missing | missing | 163.1814 | 152.1131 | 119.4367 | missing |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | cte_joined_outer vs join_where | 2.9576 | 2.3839 | 0.8060 | false |
| 20 | cte_joined_outer_mat vs join_where | 2.9576 | 138.5601 | 46.8484 | false |
| 20 | cte_joined_outer_notmat vs join_where | 2.9576 | 2.1410 | 0.7239 | false |
| 2000 | cte_joined_outer vs join_where | 10.1490 | 11.7377 | 1.1565 | false |
| 2000 | cte_joined_outer_mat vs join_where | 10.1490 | 99.0413 | 9.7587 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 10.1490 | 10.0796 | 0.9932 | false |
| 180000 | cte_joined_outer vs join_where | 160.0020 | 163.1814 | 1.0199 | false |
| 180000 | cte_joined_outer_mat vs join_where | 160.0020 | 152.1131 | 0.9507 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 160.0020 | 119.4367 | 0.7465 | false |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_joined_outer_mat vs cte_joined_outer | 2.3839 | 138.5601 | 58.1223 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 2.3839 | 2.1410 | 0.8981 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 11.7377 | 99.0413 | 8.4379 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 11.7377 | 10.0796 | 0.8587 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 163.1814 | 152.1131 | 0.9322 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 163.1814 | 119.4367 | 0.7319 | false | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 8.5783 | 1 | `86b440faddc84065` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 10.6686 | 1 | `01940907d5877551` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | missing | | | | |
| default | 2 | 100% | missing | | | | |
| default | 4 | 1% | 10.0524 | 4 | `d05f95aba30128c6` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 10.8159 | 4 | `0de57bff7d7cbc7c` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | missing | | | | |
| default | 8 | 100% | missing | | | | |
| default | 16 | 1% | 11.7796 | 16 | `cdd39a9e81f039ca` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 12.4767 | 16 | `7e77d38cffa10e69` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 10.6174 | 1 | `ec978509cdda97d7` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 10.2798 | 1 | `9d2ad1762f225a14` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | missing | | | | |
| mat | 2 | 100% | missing | | | | |
| mat | 4 | 1% | 10.9263 | 4 | `85331df27460366f` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 10.8191 | 4 | `a008ed754076d4a3` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | missing | | | | |
| mat | 8 | 100% | missing | | | | |
| mat | 16 | 1% | 11.2960 | 16 | `d3c82f8af711e30e` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 10.9831 | 16 | `af7f12349f1d5919` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 8.6992 | 1 | `1f425322838b069e` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 10.0981 | 1 | `c78ec7f625ad314a` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | missing | | | | |
| notmat | 2 | 100% | missing | | | | |
| notmat | 4 | 1% | 27.5075 | 4 | `2a9b59ea14f5043a` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 38.4589 | 4 | `a9e2ee7b851d82fc` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | missing | | | | |
| notmat | 8 | 100% | missing | | | | |
| notmat | 16 | 1% | 96.1861 | 16 | `621deaa879cebb91` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 116.9069 | 16 | `49064fe02ef3582e` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 1.2377 | false |
| notmat | 1 | 1% | 1.0141 | false |
| mat | 1 | 100% | 0.9636 | false |
| notmat | 1 | 100% | 0.9465 | false |
| mat | 4 | 1% | 1.0869 | false |
| notmat | 4 | 1% | 2.7364 | false |
| mat | 4 | 100% | 1.0003 | false |
| notmat | 4 | 100% | 3.5558 | false |
| mat | 16 | 1% | 0.9590 | false |
| notmat | 16 | 1% | 8.1655 | false |
| mat | 16 | 100% | 0.8803 | false |
| notmat | 16 | 100% | 9.3700 | false |

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

- client: nr_periods 352, nr_throttled 128, throttled_usec 5343505 | quota 100000 100000
, memory 2147483648

- yb-n1: nr_periods 628, nr_throttled 5, throttled_usec 553282 | quota 200000 100000, memory 3221225472
- yb-n2: nr_periods 621, nr_throttled 6, throttled_usec 969611 | quota 200000 100000, memory 3221225472
- yb-n3: nr_periods 617, nr_throttled 13, throttled_usec 354926 | quota 200000 100000, memory 3221225472

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
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 158.0604 | 163.1814 | 170.4584 | 5 | 6.1 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 2.3202 | 2.3839 | 3.3443 | 5 | 419.5 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 11.2657 | 11.7377 | 12.6479 | 5 | 85.2 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 150.1687 | 152.1131 | 159.7173 | 5 | 6.6 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 133.5978 | 138.5601 | 150.8648 | 5 | 7.2 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 97.6332 | 99.0413 | 102.7265 | 5 | 10.1 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 117.4776 | 119.4367 | 127.7585 | 5 | 8.4 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 1.8705 | 2.1410 | 2.3882 | 5 | 467.1 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 9.9320 | 10.0796 | 10.4856 | 5 | 99.2 |
| family=bag;stmt=join_where;threshold=180000 | valid | 156.0450 | 160.0020 | 164.6139 | 5 | 6.2 |
| family=bag;stmt=join_where;threshold=20 | valid | 2.7561 | 2.9576 | 3.6217 | 5 | 338.1 |
| family=bag;stmt=join_where;threshold=2000 | valid | 9.3466 | 10.1490 | 11.5404 | 5 | 98.5 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 11.4731 | 11.7796 | 12.0622 | 5 | 84.9 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 12.1719 | 12.4767 | 12.8012 | 5 | 80.1 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 8.2188 | 8.5783 | 10.7410 | 5 | 116.6 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 9.9754 | 10.6686 | 11.1698 | 5 | 93.7 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 9.7458 | 10.0524 | 10.8930 | 5 | 99.5 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 10.2899 | 10.8159 | 12.0530 | 5 | 92.5 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 10.3803 | 11.2960 | 11.4402 | 5 | 88.5 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 10.4178 | 10.9831 | 11.3209 | 5 | 91.0 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 9.6511 | 10.6174 | 11.3667 | 5 | 94.2 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 9.5698 | 10.2798 | 10.9884 | 5 | 97.3 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 10.1379 | 10.9263 | 11.0365 | 5 | 91.5 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 10.6237 | 10.8191 | 11.4432 | 5 | 92.4 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 87.6633 | 96.1861 | 102.5416 | 5 | 10.4 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 100.4618 | 116.9069 | 140.3688 | 5 | 8.6 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 8.0221 | 8.6992 | 10.3543 | 5 | 115.0 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 9.4526 | 10.0981 | 10.5911 | 5 | 99.0 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 26.3309 | 27.5075 | 32.2623 | 5 | 36.4 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 37.3026 | 38.4589 | 41.7909 | 5 | 26.0 |

## Matched blended workload

- phase: `blended`; target offered rate: 10 updates/s
- offered 178, committed 178, failed 0 over 17.756 s
- offered rate 10.025/s, completed rate 10.025/s
- writer operations per endpoint: [60 59 59]
- reader operations overlapping a committed write: 78 of 210
- reader retries of transient conflicts: 0
- per key: 178 keys, min 1, max 1 per key
- ledger: 178 keys reconciled, 0 mismatch
- update latency: avg 6.412 ms, min 2.494, max 173.690 over 178

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42233159`, 5 interleaved rounds
  1. family=reuse;policy=default;refs=16;consumer=1%
  2. family=bag;stmt=cte_joined_outer;threshold=2000
  3. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  4. family=reuse;policy=default;refs=1;consumer=100%
  5. family=bag;stmt=join_where;threshold=180000
  6. family=reuse;policy=mat;refs=16;consumer=1%
  7. family=bag;stmt=cte_joined_outer;threshold=180000
  8. family=reuse;policy=mat;refs=4;consumer=100%
  9. family=reuse;policy=notmat;refs=4;consumer=1%
  10. family=reuse;policy=mat;refs=1;consumer=100%
  11. family=bag;stmt=join_where;threshold=20
  12. family=reuse;policy=mat;refs=16;consumer=100%
  13. family=reuse;policy=notmat;refs=1;consumer=100%
  14. family=reuse;policy=default;refs=4;consumer=1%
  15. family=reuse;policy=notmat;refs=16;consumer=100%
  16. family=reuse;policy=notmat;refs=4;consumer=100%
  17. family=reuse;policy=mat;refs=4;consumer=1%
  18. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  19. family=reuse;policy=mat;refs=1;consumer=1%
  20. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  21. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  22. family=bag;stmt=cte_joined_outer_mat;threshold=20
  23. family=reuse;policy=default;refs=1;consumer=1%
  24. family=bag;stmt=join_where;threshold=2000
  25. family=reuse;policy=notmat;refs=1;consumer=1%
  26. family=reuse;policy=notmat;refs=16;consumer=1%
  27. family=bag;stmt=cte_joined_outer;threshold=20
  28. family=reuse;policy=default;refs=16;consumer=100%
  29. family=reuse;policy=default;refs=4;consumer=100%
  30. family=bag;stmt=cte_joined_outer_mat;threshold=2000

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
| family=bag;stmt=cte_joined_outer;threshold=180000 | j05_cte_joined_outer | valid | `329e1eeb22a68f55` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20 | j05_cte_joined_outer | valid | `404c7ff24dd6256e` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=2000 | j05_cte_joined_outer | valid | `c3b7923e59318dd2` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | j05m_cte_joined_outer_mat | valid | `b006cc456715fc8f` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | j05m_cte_joined_outer_mat | valid | `ad9cfcd17b3fbb4b` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | j05m_cte_joined_outer_mat | valid | `09c006feae106623` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | j05n_cte_joined_outer_notmat | valid | `41cc3f0995c99624` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | j05n_cte_joined_outer_notmat | valid | `ed114d92556a5ab4` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | j05n_cte_joined_outer_notmat | valid | `4bdcc4c1ef96e99c` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=join_where;threshold=180000 | j01_join_where | valid | `85750eba71be7793` |  |
| family=bag;stmt=join_where;threshold=20 | j01_join_where | valid | `f29f91a706cf54f5` |  |
| family=bag;stmt=join_where;threshold=2000 | j01_join_where | valid | `9a224d5d13c8115c` |  |
| family=reuse;policy=default;refs=16;consumer=1% | r_default_n16_c1 | valid | `cdd39a9e81f039ca` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=16;consumer=100% | r_default_n16_c100 | valid | `7e77d38cffa10e69` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=1;consumer=1% | r_default_n1_c1 | valid | `86b440faddc84065` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=1;consumer=100% | r_default_n1_c100 | valid | `01940907d5877551` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=4;consumer=1% | r_default_n4_c1 | valid | `d05f95aba30128c6` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=100% | r_default_n4_c100 | valid | `0de57bff7d7cbc7c` | policy=default; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=1% | r_mat_n16_c1 | valid | `d3c82f8af711e30e` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=100% | r_mat_n16_c100 | valid | `af7f12349f1d5919` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=1% | r_mat_n1_c1 | valid | `ec978509cdda97d7` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=100% | r_mat_n1_c100 | valid | `9d2ad1762f225a14` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=1% | r_mat_n4_c1 | valid | `85331df27460366f` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=100% | r_mat_n4_c100 | valid | `a008ed754076d4a3` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=notmat;refs=16;consumer=1% | r_notmat_n16_c1 | valid | `621deaa879cebb91` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `49064fe02ef3582e` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `1f425322838b069e` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `c78ec7f625ad314a` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `2a9b59ea14f5043a` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `a9e2ee7b851d82fc` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, DIST, DEBUG, FORMAT TEXT` supported on yugabyte
- one host: this run used node budget 2 CPU 3g and client budget 1 CPU 2g, so budget comparisons are separate classes here, not independent hardware
- one client process with one fixed maintenance writer at 10 updates/s; readers are sequential; the reported sample latency includes retry and backoff time
- cluster nodes are containers on the same host, so these are not network or independent-hardware measurements; container colocation is not data-colocation evidence
- readers round-robin across 3 endpoints; the writer's endpoint policy is recorded separately
- physical rank-index arm yb-range-rank; an engine whose planner ignores the index is a different physical design, not a faster one
- YugabyteDB DIST/DEBUG storage counters (rows and index rows scanned) are parsed per cell; RPC-level counters are not aggregated
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
