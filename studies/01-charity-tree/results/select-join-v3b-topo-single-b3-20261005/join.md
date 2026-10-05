# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 30 cells, 30 valid, engine `yugabyte`, topology `single`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-v3b-topo-single-b3-20261005`
- run tag: `run/study-01/select-join-v3b-topo-single-b3-20261005`
- environment: `host-zenbook-ux5406sa`
- repo commit: `7243fee074bfb16fc5ad8b6afc1a2d7549080935`
- engine version: `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit`
- database engine image: `docker.io/yugabytedb/yugabyte:2025.2.6.0-b111` @ `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`
- benchmark client image: `localhost/charitybench-select:select-join-v3b-topo-single-b3-20261005` @ `f3f1e10c44cbcb4c7b5bc56ed95a4f071682b6f7d33b5e0a7b688c4271eeef88`
- inputs digest: `dbae495403d90aeafd9c88fef03a25f269cb2f38ceb8ade3a3417071dbbf5eed`
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
- setting `node_budget`: 2 CPU 3g
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
| 20 (0.01%) | 2.0190 | missing | missing | missing | missing | missing | 1.8656 | 103.4108 | 2.0867 | missing |
| 2000 (1%) | 9.8052 | missing | missing | missing | missing | missing | 9.2314 | 89.0833 | 10.1596 | missing |
| 180000 (90%) | 114.1303 | missing | missing | missing | missing | missing | 126.5837 | 104.8720 | 116.2756 | missing |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | cte_joined_outer vs join_where | 2.0190 | 1.8656 | 0.9240 | false |
| 20 | cte_joined_outer_mat vs join_where | 2.0190 | 103.4108 | 51.2181 | false |
| 20 | cte_joined_outer_notmat vs join_where | 2.0190 | 2.0867 | 1.0335 | false |
| 2000 | cte_joined_outer vs join_where | 9.8052 | 9.2314 | 0.9415 | false |
| 2000 | cte_joined_outer_mat vs join_where | 9.8052 | 89.0833 | 9.0853 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 9.8052 | 10.1596 | 1.0361 | false |
| 180000 | cte_joined_outer vs join_where | 114.1303 | 126.5837 | 1.1091 | false |
| 180000 | cte_joined_outer_mat vs join_where | 114.1303 | 104.8720 | 0.9189 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 114.1303 | 116.2756 | 1.0188 | false |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_joined_outer_mat vs cte_joined_outer | 1.8656 | 103.4108 | 55.4301 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 1.8656 | 2.0867 | 1.1185 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 9.2314 | 89.0833 | 9.6500 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 9.2314 | 10.1596 | 1.1005 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 126.5837 | 104.8720 | 0.8285 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 126.5837 | 116.2756 | 0.9186 | false | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 7.6182 | 1 | `4e664e6f83fb0988` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 8.3830 | 1 | `c0b133c19c912d8d` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | missing | | | | |
| default | 2 | 100% | missing | | | | |
| default | 4 | 1% | 9.6179 | 4 | `ac2f1f9388484b04` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 9.2042 | 4 | `44cfaff1538e8d2a` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | missing | | | | |
| default | 8 | 100% | missing | | | | |
| default | 16 | 1% | 9.5205 | 16 | `5bfacddf1af24cdc` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 10.0204 | 16 | `401f2193477159e0` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 9.7282 | 1 | `a20c1d4dbe61be4d` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 9.3078 | 1 | `35376211b4d5a600` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | missing | | | | |
| mat | 2 | 100% | missing | | | | |
| mat | 4 | 1% | 9.3609 | 4 | `01c006b5cda091f9` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 9.2607 | 4 | `c8e2f399606a4ff3` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | missing | | | | |
| mat | 8 | 100% | missing | | | | |
| mat | 16 | 1% | 9.6711 | 16 | `33824d4097439cea` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 10.0578 | 16 | `c10f6cf71b1a09f4` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 7.9971 | 1 | `8e93deb44b291c72` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 8.9194 | 1 | `f0e0164b770fad57` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | missing | | | | |
| notmat | 2 | 100% | missing | | | | |
| notmat | 4 | 1% | 26.0548 | 4 | `f1d27f2537de9704` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 30.7302 | 4 | `c58c3a23e9890f7f` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | missing | | | | |
| notmat | 8 | 100% | missing | | | | |
| notmat | 16 | 1% | 82.4492 | 16 | `b9f0e7c7acac0715` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 122.0000 | 16 | `af36d925c8e3383f` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 1.2770 | false |
| notmat | 1 | 1% | 1.0497 | false |
| mat | 1 | 100% | 1.1103 | false |
| notmat | 1 | 100% | 1.0640 | false |
| mat | 4 | 1% | 0.9733 | false |
| notmat | 4 | 1% | 2.7090 | false |
| mat | 4 | 100% | 1.0061 | false |
| notmat | 4 | 100% | 3.3387 | false |
| mat | 16 | 1% | 1.0158 | false |
| notmat | 16 | 1% | 8.6602 | false |
| mat | 16 | 100% | 1.0037 | false |
| notmat | 16 | 100% | 12.1752 | false |

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

- client: nr_periods 335, nr_throttled 141, throttled_usec 6213880 | quota 100000 100000
, memory 2147483648

- yb-single: nr_periods 411, nr_throttled 88, throttled_usec 6158308 | quota 200000 100000, memory 3221225472

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
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 113.1397 | 126.5837 | 149.1663 | 5 | 7.9 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 1.8538 | 1.8656 | 1.9478 | 5 | 536.0 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 9.0773 | 9.2314 | 9.9524 | 5 | 108.3 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 102.4441 | 104.8720 | 130.3536 | 5 | 9.5 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 97.7156 | 103.4108 | 124.8039 | 5 | 9.7 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 85.2712 | 89.0833 | 98.4628 | 5 | 11.2 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 110.9351 | 116.2756 | 121.7090 | 5 | 8.6 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 1.9741 | 2.0867 | 2.1670 | 5 | 479.2 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 9.7587 | 10.1596 | 10.8823 | 5 | 98.4 |
| family=bag;stmt=join_where;threshold=180000 | valid | 110.9732 | 114.1303 | 115.9111 | 5 | 8.8 |
| family=bag;stmt=join_where;threshold=20 | valid | 1.8366 | 2.0190 | 2.1553 | 5 | 495.3 |
| family=bag;stmt=join_where;threshold=2000 | valid | 8.9691 | 9.8052 | 10.1503 | 5 | 102.0 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 9.0457 | 9.5205 | 9.9749 | 5 | 105.0 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 9.8534 | 10.0204 | 11.1222 | 5 | 99.8 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 7.2407 | 7.6182 | 8.0250 | 5 | 131.3 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 8.1930 | 8.3830 | 9.3185 | 5 | 119.3 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 9.2320 | 9.6179 | 10.2849 | 5 | 104.0 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 9.0563 | 9.2042 | 9.6965 | 5 | 108.6 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 9.0696 | 9.6711 | 9.8812 | 5 | 103.4 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 9.7335 | 10.0578 | 10.9149 | 5 | 99.4 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 8.8356 | 9.7282 | 11.0660 | 5 | 102.8 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 8.7939 | 9.3078 | 10.0092 | 5 | 107.4 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 8.8204 | 9.3609 | 9.9448 | 5 | 106.8 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 8.9585 | 9.2607 | 9.3541 | 5 | 108.0 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 79.4164 | 82.4492 | 121.2426 | 5 | 12.1 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 94.5156 | 122.0000 | 126.6814 | 5 | 8.2 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 7.4363 | 7.9971 | 9.2726 | 5 | 125.0 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 8.7153 | 8.9194 | 9.3744 | 5 | 112.1 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 24.1282 | 26.0548 | 29.0487 | 5 | 38.4 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 29.3124 | 30.7302 | 32.0680 | 5 | 32.5 |

## Matched blended workload

- phase: `blended`; target offered rate: 10 updates/s
- offered 160, committed 160, failed 0 over 15.975 s
- offered rate 10.016/s, completed rate 10.016/s
- writer operations per endpoint: [160]
- reader operations overlapping a committed write: 65 of 210
- reader retries of transient conflicts: 0
- per key: 160 keys, min 1, max 1 per key
- ledger: 160 keys reconciled, 0 mismatch
- update latency: avg 5.676 ms, min 1.577, max 77.579 over 160

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42302301`, 5 interleaved rounds
  1. family=bag;stmt=join_where;threshold=180000
  2. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  3. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  4. family=reuse;policy=default;refs=4;consumer=1%
  5. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  6. family=reuse;policy=default;refs=1;consumer=1%
  7. family=reuse;policy=notmat;refs=16;consumer=100%
  8. family=reuse;policy=notmat;refs=1;consumer=1%
  9. family=bag;stmt=cte_joined_outer;threshold=180000
  10. family=reuse;policy=notmat;refs=16;consumer=1%
  11. family=bag;stmt=cte_joined_outer_mat;threshold=20
  12. family=reuse;policy=mat;refs=1;consumer=100%
  13. family=reuse;policy=notmat;refs=4;consumer=1%
  14. family=reuse;policy=default;refs=16;consumer=100%
  15. family=bag;stmt=join_where;threshold=2000
  16. family=reuse;policy=mat;refs=4;consumer=1%
  17. family=bag;stmt=join_where;threshold=20
  18. family=reuse;policy=notmat;refs=1;consumer=100%
  19. family=reuse;policy=default;refs=16;consumer=1%
  20. family=reuse;policy=default;refs=1;consumer=100%
  21. family=bag;stmt=cte_joined_outer;threshold=20
  22. family=bag;stmt=cte_joined_outer;threshold=2000
  23. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  24. family=reuse;policy=mat;refs=1;consumer=1%
  25. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  26. family=reuse;policy=default;refs=4;consumer=100%
  27. family=reuse;policy=notmat;refs=4;consumer=100%
  28. family=reuse;policy=mat;refs=16;consumer=1%
  29. family=reuse;policy=mat;refs=16;consumer=100%
  30. family=reuse;policy=mat;refs=4;consumer=100%

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
| family=bag;stmt=cte_joined_outer;threshold=180000 | j05_cte_joined_outer | valid | `1c6bd98300f06f70` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20 | j05_cte_joined_outer | valid | `fd21775c7b6a27ec` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=2000 | j05_cte_joined_outer | valid | `15e64a2430383a99` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | j05m_cte_joined_outer_mat | valid | `ec45482ef77498cd` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | j05m_cte_joined_outer_mat | valid | `0b13f13ca9d6939b` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | j05m_cte_joined_outer_mat | valid | `de1dede1970a6577` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | j05n_cte_joined_outer_notmat | valid | `937be48e53b6cf32` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | j05n_cte_joined_outer_notmat | valid | `f6e75fe070eea5a4` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | j05n_cte_joined_outer_notmat | valid | `c825a5b08d64149d` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=join_where;threshold=180000 | j01_join_where | valid | `3217548a1920c8a2` |  |
| family=bag;stmt=join_where;threshold=20 | j01_join_where | valid | `a06919591f9da599` |  |
| family=bag;stmt=join_where;threshold=2000 | j01_join_where | valid | `6981e0f16a824281` |  |
| family=reuse;policy=default;refs=16;consumer=1% | r_default_n16_c1 | valid | `5bfacddf1af24cdc` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=16;consumer=100% | r_default_n16_c100 | valid | `401f2193477159e0` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=1;consumer=1% | r_default_n1_c1 | valid | `4e664e6f83fb0988` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=1;consumer=100% | r_default_n1_c100 | valid | `c0b133c19c912d8d` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=4;consumer=1% | r_default_n4_c1 | valid | `ac2f1f9388484b04` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=100% | r_default_n4_c100 | valid | `44cfaff1538e8d2a` | policy=default; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=1% | r_mat_n16_c1 | valid | `33824d4097439cea` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=100% | r_mat_n16_c100 | valid | `c10f6cf71b1a09f4` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=1% | r_mat_n1_c1 | valid | `a20c1d4dbe61be4d` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=100% | r_mat_n1_c100 | valid | `35376211b4d5a600` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=1% | r_mat_n4_c1 | valid | `01c006b5cda091f9` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=100% | r_mat_n4_c100 | valid | `c8e2f399606a4ff3` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=notmat;refs=16;consumer=1% | r_notmat_n16_c1 | valid | `b9f0e7c7acac0715` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `af36d925c8e3383f` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `8e93deb44b291c72` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `f0e0164b770fad57` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `f1d27f2537de9704` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `c58c3a23e9890f7f` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, DIST, DEBUG, FORMAT TEXT` supported on yugabyte
- one host: this run used node budget 2 CPU 3g and client budget 1 CPU 2g, so budget comparisons are separate classes here, not independent hardware
- one client process with one fixed maintenance writer at 10 updates/s; readers are sequential; the reported sample latency includes retry and backoff time
- physical rank-index arm yb-range-rank; an engine whose planner ignores the index is a different physical design, not a faster one
- YugabyteDB DIST/DEBUG storage counters (rows and index rows scanned) are parsed per cell; RPC-level counters are not aggregated
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
