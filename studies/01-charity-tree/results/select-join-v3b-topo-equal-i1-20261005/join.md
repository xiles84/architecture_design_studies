# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 30 cells, 30 valid, engine `yugabyte`, topology `single`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-v3b-topo-equal-i1-20261005`
- run tag: `run/study-01/select-join-v3b-topo-equal-i1-20261005`
- environment: `host-zenbook-ux5406sa`
- repo commit: `7243fee074bfb16fc5ad8b6afc1a2d7549080935`
- engine version: `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit`
- database engine image: `docker.io/yugabytedb/yugabyte:2025.2.6.0-b111` @ `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`
- benchmark client image: `localhost/charitybench-select:select-join-v3b-topo-equal-i1-20261005` @ `f3f1e10c44cbcb4c7b5bc56ed95a4f071682b6f7d33b5e0a7b688c4271eeef88`
- inputs digest: `4d9d9c88b5d4abd291b5013a74b670c228e15a716a6129ccbeb1eaca813b26ac`
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
| 20 (0.01%) | 1.8066 | missing | missing | missing | missing | missing | 1.7763 | 86.1375 | 1.9820 | missing |
| 2000 (1%) | 9.3060 | missing | missing | missing | missing | missing | 9.6070 | 85.8868 | 9.8174 | missing |
| 180000 (90%) | 108.9713 | missing | missing | missing | missing | missing | 103.4669 | 103.9693 | 110.4309 | missing |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | cte_joined_outer vs join_where | 1.8066 | 1.7763 | 0.9832 | false |
| 20 | cte_joined_outer_mat vs join_where | 1.8066 | 86.1375 | 47.6794 | false |
| 20 | cte_joined_outer_notmat vs join_where | 1.8066 | 1.9820 | 1.0971 | false |
| 2000 | cte_joined_outer vs join_where | 9.3060 | 9.6070 | 1.0323 | false |
| 2000 | cte_joined_outer_mat vs join_where | 9.3060 | 85.8868 | 9.2291 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 9.3060 | 9.8174 | 1.0549 | false |
| 180000 | cte_joined_outer vs join_where | 108.9713 | 103.4669 | 0.9495 | false |
| 180000 | cte_joined_outer_mat vs join_where | 108.9713 | 103.9693 | 0.9541 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 108.9713 | 110.4309 | 1.0134 | false |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_joined_outer_mat vs cte_joined_outer | 1.7763 | 86.1375 | 48.4920 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 1.7763 | 1.9820 | 1.1158 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 9.6070 | 85.8868 | 8.9400 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 9.6070 | 9.8174 | 1.0219 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 103.4669 | 103.9693 | 1.0049 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 103.4669 | 110.4309 | 1.0673 | false | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 7.2165 | 1 | `087055ad80db1558` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 8.2153 | 1 | `ee2f347d5650689a` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | missing | | | | |
| default | 2 | 100% | missing | | | | |
| default | 4 | 1% | 8.6983 | 4 | `a7689a46aede15b2` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 8.8821 | 4 | `592340e8e9acacc9` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | missing | | | | |
| default | 8 | 100% | missing | | | | |
| default | 16 | 1% | 9.0918 | 16 | `40f69425ebdfc54f` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 10.0379 | 16 | `9a33b72d8eefb75b` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 8.7694 | 1 | `74e566257ed0c71f` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 8.9848 | 1 | `cbea172deb5cda23` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | missing | | | | |
| mat | 2 | 100% | missing | | | | |
| mat | 4 | 1% | 8.7947 | 4 | `4b22de01941003f2` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 8.6993 | 4 | `e5351a4cd0f065b9` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | missing | | | | |
| mat | 8 | 100% | missing | | | | |
| mat | 16 | 1% | 9.1317 | 16 | `baaf9c100ff64a20` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 10.2754 | 16 | `c2bef3223d85bf8e` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 7.4272 | 1 | `71bd3b64514a95df` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 9.4589 | 1 | `9d414757fe282160` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | missing | | | | |
| notmat | 2 | 100% | missing | | | | |
| notmat | 4 | 1% | 24.6731 | 4 | `8745952d55380d1c` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 30.2062 | 4 | `1f5e2172e1045471` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | missing | | | | |
| notmat | 8 | 100% | missing | | | | |
| notmat | 16 | 1% | 72.2798 | 16 | `5bb087e901dabc4b` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 83.9565 | 16 | `666a93176e2951b8` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 1.2152 | false |
| notmat | 1 | 1% | 1.0292 | false |
| mat | 1 | 100% | 1.0937 | false |
| notmat | 1 | 100% | 1.1514 | false |
| mat | 4 | 1% | 1.0111 | false |
| notmat | 4 | 1% | 2.8365 | false |
| mat | 4 | 100% | 0.9794 | false |
| notmat | 4 | 100% | 3.4008 | false |
| mat | 16 | 1% | 1.0044 | false |
| notmat | 16 | 1% | 7.9500 | false |
| mat | 16 | 100% | 1.0237 | false |
| notmat | 16 | 100% | 8.3640 | false |

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

- client: nr_periods 278, nr_throttled 115, throttled_usec 3564778 | quota 100000 100000
, memory 2147483648

- yb-single: nr_periods 357, nr_throttled 0, throttled_usec 0 | quota 600000 100000, memory 9663676416

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
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 100.8026 | 103.4669 | 118.0125 | 5 | 9.7 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 1.6925 | 1.7763 | 1.8530 | 5 | 563.0 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 8.2167 | 9.6070 | 10.7435 | 5 | 104.1 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 91.6097 | 103.9693 | 109.8428 | 5 | 9.6 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 82.0056 | 86.1375 | 93.1991 | 5 | 11.6 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 82.8117 | 85.8868 | 88.7481 | 5 | 11.6 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 103.3030 | 110.4309 | 119.0890 | 5 | 9.1 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 1.7918 | 1.9820 | 2.3429 | 5 | 504.5 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 9.1077 | 9.8174 | 9.9199 | 5 | 101.9 |
| family=bag;stmt=join_where;threshold=180000 | valid | 106.6567 | 108.9713 | 114.7080 | 5 | 9.2 |
| family=bag;stmt=join_where;threshold=20 | valid | 1.6641 | 1.8066 | 2.3002 | 5 | 553.5 |
| family=bag;stmt=join_where;threshold=2000 | valid | 8.4674 | 9.3060 | 10.3950 | 5 | 107.5 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 8.4980 | 9.0918 | 10.1620 | 5 | 110.0 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 9.4565 | 10.0379 | 10.9253 | 5 | 99.6 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 6.0169 | 7.2165 | 8.6656 | 5 | 138.6 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 7.8734 | 8.2153 | 30.4337 | 5 | 121.7 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 7.9631 | 8.6983 | 8.7806 | 5 | 115.0 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 7.9764 | 8.8821 | 10.1187 | 5 | 112.6 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 8.5894 | 9.1317 | 9.7346 | 5 | 109.5 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 9.5628 | 10.2754 | 10.7316 | 5 | 97.3 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 8.3908 | 8.7694 | 9.3969 | 5 | 114.0 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 8.6256 | 8.9848 | 9.3299 | 5 | 111.3 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 8.0848 | 8.7947 | 9.3405 | 5 | 113.7 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 8.2352 | 8.6993 | 10.4830 | 5 | 115.0 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 69.8136 | 72.2798 | 77.0265 | 5 | 13.8 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 78.1119 | 83.9565 | 109.4036 | 5 | 11.9 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 6.8453 | 7.4272 | 8.4648 | 5 | 134.6 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 8.3523 | 9.4589 | 10.9424 | 5 | 105.7 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 21.1569 | 24.6731 | 25.5084 | 5 | 40.5 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 28.7929 | 30.2062 | 33.8447 | 5 | 33.1 |

## Matched blended workload

- phase: `isolated`; target offered rate: 10 updates/s
- offered 0, committed 0, failed 0 over 0.000 s
- offered rate 0.000/s, completed rate 0.000/s
- per key: no keys touched
- ledger: 0 keys reconciled, 0 mismatch

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42325590`, 5 interleaved rounds
  1. family=reuse;policy=notmat;refs=4;consumer=100%
  2. family=reuse;policy=default;refs=16;consumer=1%
  3. family=reuse;policy=default;refs=4;consumer=1%
  4. family=bag;stmt=cte_joined_outer;threshold=180000
  5. family=reuse;policy=default;refs=16;consumer=100%
  6. family=reuse;policy=default;refs=1;consumer=100%
  7. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  8. family=reuse;policy=notmat;refs=1;consumer=1%
  9. family=reuse;policy=notmat;refs=16;consumer=100%
  10. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  11. family=reuse;policy=notmat;refs=1;consumer=100%
  12. family=reuse;policy=default;refs=1;consumer=1%
  13. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  14. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  15. family=reuse;policy=mat;refs=1;consumer=1%
  16. family=reuse;policy=mat;refs=16;consumer=100%
  17. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  18. family=bag;stmt=cte_joined_outer_mat;threshold=20
  19. family=bag;stmt=join_where;threshold=180000
  20. family=reuse;policy=mat;refs=1;consumer=100%
  21. family=reuse;policy=mat;refs=16;consumer=1%
  22. family=bag;stmt=join_where;threshold=20
  23. family=reuse;policy=mat;refs=4;consumer=1%
  24. family=reuse;policy=mat;refs=4;consumer=100%
  25. family=reuse;policy=default;refs=4;consumer=100%
  26. family=bag;stmt=cte_joined_outer;threshold=2000
  27. family=bag;stmt=join_where;threshold=2000
  28. family=reuse;policy=notmat;refs=4;consumer=1%
  29. family=bag;stmt=cte_joined_outer;threshold=20
  30. family=reuse;policy=notmat;refs=16;consumer=1%

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
| family=bag;stmt=cte_joined_outer;threshold=180000 | j05_cte_joined_outer | valid | `baf6f9e9c2d01899` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20 | j05_cte_joined_outer | valid | `9c19ec294201fe93` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=2000 | j05_cte_joined_outer | valid | `a3752438acc68895` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | j05m_cte_joined_outer_mat | valid | `f1495295463e744d` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | j05m_cte_joined_outer_mat | valid | `b42d6c12a480d77b` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | j05m_cte_joined_outer_mat | valid | `b2a8318e491ea892` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | j05n_cte_joined_outer_notmat | valid | `8d541cb106807d76` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | j05n_cte_joined_outer_notmat | valid | `fa075c7111d0f010` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | j05n_cte_joined_outer_notmat | valid | `356872c6be2812b5` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=join_where;threshold=180000 | j01_join_where | valid | `c5513623782fb472` |  |
| family=bag;stmt=join_where;threshold=20 | j01_join_where | valid | `89c13832100f755b` |  |
| family=bag;stmt=join_where;threshold=2000 | j01_join_where | valid | `69a874de93e261ff` |  |
| family=reuse;policy=default;refs=16;consumer=1% | r_default_n16_c1 | valid | `40f69425ebdfc54f` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=16;consumer=100% | r_default_n16_c100 | valid | `9a33b72d8eefb75b` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=1;consumer=1% | r_default_n1_c1 | valid | `087055ad80db1558` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=1;consumer=100% | r_default_n1_c100 | valid | `ee2f347d5650689a` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=4;consumer=1% | r_default_n4_c1 | valid | `a7689a46aede15b2` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=100% | r_default_n4_c100 | valid | `592340e8e9acacc9` | policy=default; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=1% | r_mat_n16_c1 | valid | `baaf9c100ff64a20` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=100% | r_mat_n16_c100 | valid | `c2bef3223d85bf8e` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=1% | r_mat_n1_c1 | valid | `74e566257ed0c71f` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=100% | r_mat_n1_c100 | valid | `cbea172deb5cda23` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=1% | r_mat_n4_c1 | valid | `4b22de01941003f2` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=100% | r_mat_n4_c100 | valid | `e5351a4cd0f065b9` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=notmat;refs=16;consumer=1% | r_notmat_n16_c1 | valid | `5bb087e901dabc4b` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `666a93176e2951b8` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `71bd3b64514a95df` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `9d414757fe282160` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `8745952d55380d1c` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `1f5e2172e1045471` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, DIST, DEBUG, FORMAT TEXT` supported on yugabyte
- one host: this run used node budget 6 CPU 9g and client budget 1 CPU 2g, so budget comparisons are separate classes here, not independent hardware
- isolated phase: no concurrent writers in this run; the blended read/write phase is a separate labelled run
- physical rank-index arm yb-range-rank; an engine whose planner ignores the index is a different physical design, not a faster one
- YugabyteDB DIST/DEBUG storage counters (rows and index rows scanned) are parsed per cell; RPC-level counters are not aggregated
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
