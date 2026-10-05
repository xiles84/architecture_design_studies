# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 30 cells, 30 valid, engine `yugabyte`, topology `cluster3`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-v4-topo-tablet-b1-20261005`
- run tag: `run/study-01/select-join-v4-topo-tablet-b1-20261005`
- environment: `host-zenbook-ux5406sa`
- repo commit: `ba821e7b686ac3adb1cda56dd45f15c8ae699ef8`
- engine version: `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit`
- database engine image: `docker.io/yugabytedb/yugabyte:2025.2.6.0-b111` @ `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`
- benchmark client image: `localhost/charitybench-select:select-join-v4-topo-tablet-b1-20261005` @ `b6ea3c007cf459112eaa8e9e85381a7dc25a7710ff07648fde3048f747df6ab0`
- inputs digest: `cc9d1cbf384650d419ca7fb3f16bcc9905806cda53b4f640ebe25024346aaaa2`
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
- setting `placement`: noncolocated-single-tablet
- setting `prepared_plan_policy`: pgx simple protocol; no prepared statements
- setting `scale`: cte-domain
- setting `selectivity_mapping`: rank <= threshold; rank is a deterministic shuffle of 1..fact_count
- setting `shards_per_tserver`: ysql_num_shards_per_tserver=2
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
| j01_join_where | valid | true |  |
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
| 20 (0.01%) | 2.0892 | missing | missing | missing | missing | missing | 2.0321 | 131.1921 | 3.1898 | missing |
| 2000 (1%) | 12.3004 | missing | missing | missing | missing | missing | 11.3222 | 91.4519 | 11.2379 | missing |
| 180000 (90%) | 117.4869 | missing | missing | missing | missing | missing | 124.9643 | 114.5862 | 153.9764 | missing |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | cte_joined_outer vs join_where | 2.0892 | 2.0321 | 0.9727 | false |
| 20 | cte_joined_outer_mat vs join_where | 2.0892 | 131.1921 | 62.7943 | false |
| 20 | cte_joined_outer_notmat vs join_where | 2.0892 | 3.1898 | 1.5268 | false |
| 2000 | cte_joined_outer vs join_where | 12.3004 | 11.3222 | 0.9205 | false |
| 2000 | cte_joined_outer_mat vs join_where | 12.3004 | 91.4519 | 7.4349 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 12.3004 | 11.2379 | 0.9136 | false |
| 180000 | cte_joined_outer vs join_where | 117.4869 | 124.9643 | 1.0636 | false |
| 180000 | cte_joined_outer_mat vs join_where | 117.4869 | 114.5862 | 0.9753 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 117.4869 | 153.9764 | 1.3106 | false |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_joined_outer_mat vs cte_joined_outer | 2.0321 | 131.1921 | 64.5593 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 2.0321 | 3.1898 | 1.5697 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 11.3222 | 91.4519 | 8.0772 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 11.3222 | 11.2379 | 0.9926 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 124.9643 | 114.5862 | 0.9170 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 124.9643 | 153.9764 | 1.2322 | false | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 8.0753 | 1 | `14a00f5f1ad937f3` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 10.4798 | 1 | `8198f03726a8f32e` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | missing | | | | |
| default | 2 | 100% | missing | | | | |
| default | 4 | 1% | 10.1721 | 4 | `2a91a83dfaaefb06` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 9.6506 | 4 | `5d90833085e644e4` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | missing | | | | |
| default | 8 | 100% | missing | | | | |
| default | 16 | 1% | 11.2682 | 16 | `6b862dcd379582d9` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 11.3531 | 16 | `08152107cec3664c` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 10.2271 | 1 | `edca44865121951b` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 11.0907 | 1 | `03f23851be87a1c0` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | missing | | | | |
| mat | 2 | 100% | missing | | | | |
| mat | 4 | 1% | 10.3487 | 4 | `5708f7bda5c4a66c` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 10.6480 | 4 | `804f2e621e0f45aa` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | missing | | | | |
| mat | 8 | 100% | missing | | | | |
| mat | 16 | 1% | 10.5595 | 16 | `0a61883d6471c751` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 11.7126 | 16 | `c8b9d3eb6d80745d` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 7.0769 | 1 | `16f0ac652a772348` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 10.6736 | 1 | `aaf655b47d00c5a0` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | missing | | | | |
| notmat | 2 | 100% | missing | | | | |
| notmat | 4 | 1% | 28.8182 | 4 | `23726b79b766a671` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 31.3932 | 4 | `c7222e297139fd64` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | missing | | | | |
| notmat | 8 | 100% | missing | | | | |
| notmat | 16 | 1% | 88.6976 | 16 | `b3471f19aeaeae86` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 99.4867 | 16 | `f3d05e6603228f2f` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 1.2665 | false |
| notmat | 1 | 1% | 0.8764 | false |
| mat | 1 | 100% | 1.0583 | false |
| notmat | 1 | 100% | 1.0185 | false |
| mat | 4 | 1% | 1.0174 | false |
| notmat | 4 | 1% | 2.8331 | false |
| mat | 4 | 100% | 1.1034 | false |
| notmat | 4 | 100% | 3.2530 | false |
| mat | 16 | 1% | 0.9371 | false |
| notmat | 16 | 1% | 7.8715 | false |
| mat | 16 | 100% | 1.0317 | false |
| notmat | 16 | 100% | 8.7630 | false |

## Physical arm

- `charity_pkey`: `CREATE UNIQUE INDEX charity_pkey ON selectjoin_domain.charity USING lsm (charity_id HASH)`
- `donation_person_idx`: `CREATE INDEX donation_person_idx ON selectjoin_domain.donation USING lsm (person_id HASH)`
- `donation_pkey`: `CREATE UNIQUE INDEX donation_pkey ON selectjoin_domain.donation USING lsm (donation_id HASH)`
- `donation_rank_idx`: `CREATE INDEX donation_rank_idx ON selectjoin_domain.donation USING lsm (rank ASC)`
- `person_pkey`: `CREATE UNIQUE INDEX person_pkey ON selectjoin_domain.person USING lsm (person_id HASH)`

- index arm `yb-range-rank`; placement `noncolocated-single-tablet`; node budget `2 CPU 3g`; client budget `1 CPU 2g`

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

Tablet leaders per relation (placement evidence; the runner read these through yb-admin):

| relation | leader nodes |
|---|---|
| charity | yb-n2:9100 |
| person | yb-n2:9100 |
| donation | yb-n2:9100 |
| donation_rank_idx | yb-n2:9100 |
| donation_person_idx | yb-n2:9100 |

Effective sharding flag (from the engine): `ysql_num_shards_per_tserver=2`. A relation with fewer tablets than this setting was created without an explicit split; the per-table `num_tablets` above is the authority.

## Throttling (cgroup `cpu.stat` deltas)

- client: nr_periods 329, nr_throttled 125, throttled_usec 4478417 | quota 100000 100000
, memory 2147483648

- yb-n1: nr_periods 602, nr_throttled 4, throttled_usec 529102 | quota 200000 100000, memory 3221225472
- yb-n2: nr_periods 593, nr_throttled 13, throttled_usec 2062689 | quota 200000 100000, memory 3221225472
- yb-n3: nr_periods 588, nr_throttled 0, throttled_usec 0 | quota 200000 100000, memory 3221225472

## Access path per cell

_Engine-reported from each cell's own plan. The two scan flags are **whole-plan** booleans: they are true when any node of the plan is an index scan / sequential scan, including small parent-table scans, so they are not the rank access path. `storage rows scanned` is the engine's own count (YugabyteDB); PostgreSQL does not print one and it is shown as n/a. A cell that scanned the whole child table at a selective threshold is a different physical regime, not a faster query._

| axis | status | whole-plan index scan | whole-plan seq scan | storage rows scanned | index rows scanned | buffers hit/read | rows removed | retries |
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
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 119.3761 | 124.9643 | 149.2320 | 5 | 8.0 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 1.8881 | 2.0321 | 2.1871 | 5 | 492.1 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 10.8890 | 11.3222 | 11.4154 | 5 | 88.3 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 107.2664 | 114.5862 | 136.0120 | 5 | 8.7 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 130.3865 | 131.1921 | 153.8892 | 5 | 7.6 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 88.3908 | 91.4519 | 97.9947 | 5 | 10.9 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 150.1157 | 153.9764 | 163.5232 | 5 | 6.5 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 2.7060 | 3.1898 | 3.3450 | 5 | 313.5 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 11.1277 | 11.2379 | 13.7950 | 5 | 89.0 |
| family=bag;stmt=join_where;threshold=180000 | valid | 113.6584 | 117.4869 | 125.9169 | 5 | 8.5 |
| family=bag;stmt=join_where;threshold=20 | valid | 2.0444 | 2.0892 | 2.7153 | 5 | 478.6 |
| family=bag;stmt=join_where;threshold=2000 | valid | 11.4590 | 12.3004 | 12.3267 | 5 | 81.3 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 10.6898 | 11.2682 | 11.9995 | 5 | 88.7 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 11.1463 | 11.3531 | 11.6332 | 5 | 88.1 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 7.4201 | 8.0753 | 9.1815 | 5 | 123.8 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 10.2758 | 10.4798 | 10.6837 | 5 | 95.4 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 9.8613 | 10.1721 | 10.5061 | 5 | 98.3 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 9.1474 | 9.6506 | 10.1528 | 5 | 103.6 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 10.2312 | 10.5595 | 11.7835 | 5 | 94.7 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 11.4944 | 11.7126 | 13.0034 | 5 | 85.4 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 9.2865 | 10.2271 | 12.4157 | 5 | 97.8 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 10.4504 | 11.0907 | 11.8000 | 5 | 90.2 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 9.9009 | 10.3487 | 10.8602 | 5 | 96.6 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 10.2911 | 10.6480 | 11.2390 | 5 | 93.9 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 81.5366 | 88.6976 | 91.9235 | 5 | 11.3 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 91.0790 | 99.4867 | 103.1817 | 5 | 10.1 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 6.8713 | 7.0769 | 7.7369 | 5 | 141.3 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 10.2988 | 10.6736 | 12.2752 | 5 | 93.7 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 27.7474 | 28.8182 | 30.2510 | 5 | 34.7 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 31.0044 | 31.3932 | 33.1297 | 5 | 31.9 |

## Matched blended workload

- phase: `blended`; target offered rate: 10 updates/s
- offered 167, committed 167, failed 0 over 16.629 s
- offered rate 10.043/s, completed rate 10.043/s
- writer operations per endpoint: [56 56 55]
- reader operations overlapping a committed write: 68 of 210
- reader retries of transient conflicts: 0
- per key: 167 keys, min 1, max 1 per key
- ledger: 167 keys reconciled, 0 mismatch
- update latency: avg 5.634 ms, min 2.672, max 183.814 over 167

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42933744`, 5 interleaved rounds
  1. family=reuse;policy=notmat;refs=16;consumer=100%
  2. family=bag;stmt=cte_joined_outer;threshold=180000
  3. family=reuse;policy=mat;refs=1;consumer=100%
  4. family=reuse;policy=mat;refs=4;consumer=100%
  5. family=reuse;policy=notmat;refs=4;consumer=100%
  6. family=bag;stmt=cte_joined_outer_mat;threshold=20
  7. family=reuse;policy=mat;refs=4;consumer=1%
  8. family=bag;stmt=join_where;threshold=180000
  9. family=reuse;policy=notmat;refs=4;consumer=1%
  10. family=reuse;policy=mat;refs=16;consumer=1%
  11. family=reuse;policy=notmat;refs=1;consumer=1%
  12. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  13. family=reuse;policy=default;refs=1;consumer=100%
  14. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  15. family=reuse;policy=default;refs=16;consumer=1%
  16. family=bag;stmt=cte_joined_outer;threshold=2000
  17. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  18. family=reuse;policy=notmat;refs=1;consumer=100%
  19. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  20. family=reuse;policy=default;refs=4;consumer=100%
  21. family=reuse;policy=mat;refs=1;consumer=1%
  22. family=reuse;policy=default;refs=16;consumer=100%
  23. family=bag;stmt=join_where;threshold=20
  24. family=reuse;policy=mat;refs=16;consumer=100%
  25. family=reuse;policy=default;refs=4;consumer=1%
  26. family=bag;stmt=cte_joined_outer;threshold=20
  27. family=reuse;policy=notmat;refs=16;consumer=1%
  28. family=bag;stmt=join_where;threshold=2000
  29. family=reuse;policy=default;refs=1;consumer=1%
  30. family=bag;stmt=cte_joined_outer_notmat;threshold=2000

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
| family=bag;stmt=cte_joined_outer;threshold=180000 | j05_cte_joined_outer | valid | `a107eeecfd356361` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20 | j05_cte_joined_outer | valid | `1dc67e148a51ddd1` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=2000 | j05_cte_joined_outer | valid | `561ef45209b79da5` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | j05m_cte_joined_outer_mat | valid | `06a8b0d2d795a4ab` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | j05m_cte_joined_outer_mat | valid | `fc5ef5b754e8eb8f` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | j05m_cte_joined_outer_mat | valid | `54f83c340caced0a` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | j05n_cte_joined_outer_notmat | valid | `dfcdc067ffb10699` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | j05n_cte_joined_outer_notmat | valid | `eaed1abe395d93f4` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | j05n_cte_joined_outer_notmat | valid | `d6860b3843afbd09` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=join_where;threshold=180000 | j01_join_where | valid | `fc0211587dde73e1` |  |
| family=bag;stmt=join_where;threshold=20 | j01_join_where | valid | `cd20d10d235c49fd` |  |
| family=bag;stmt=join_where;threshold=2000 | j01_join_where | valid | `b14543b6d4a0c947` |  |
| family=reuse;policy=default;refs=16;consumer=1% | r_default_n16_c1 | valid | `6b862dcd379582d9` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=16;consumer=100% | r_default_n16_c100 | valid | `08152107cec3664c` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=1;consumer=1% | r_default_n1_c1 | valid | `14a00f5f1ad937f3` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=1;consumer=100% | r_default_n1_c100 | valid | `8198f03726a8f32e` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=4;consumer=1% | r_default_n4_c1 | valid | `2a91a83dfaaefb06` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=100% | r_default_n4_c100 | valid | `5d90833085e644e4` | policy=default; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=1% | r_mat_n16_c1 | valid | `0a61883d6471c751` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=100% | r_mat_n16_c100 | valid | `c8b9d3eb6d80745d` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=1% | r_mat_n1_c1 | valid | `edca44865121951b` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=100% | r_mat_n1_c100 | valid | `03f23851be87a1c0` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=1% | r_mat_n4_c1 | valid | `5708f7bda5c4a66c` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=100% | r_mat_n4_c100 | valid | `804f2e621e0f45aa` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=notmat;refs=16;consumer=1% | r_notmat_n16_c1 | valid | `b3471f19aeaeae86` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `f3d05e6603228f2f` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `16f0ac652a772348` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `aaf655b47d00c5a0` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `23726b79b766a671` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `c7222e297139fd64` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, DIST, DEBUG, FORMAT TEXT` supported on yugabyte
- one host: this run used node budget 2 CPU 3g and client budget 1 CPU 2g, so budget comparisons are separate classes here, not independent hardware
- one client process with one fixed maintenance writer at 10 updates/s; readers are sequential; the reported sample latency includes retry and backoff time
- cluster nodes are containers on the same host, so these are not network or independent-hardware measurements; container colocation is not data-colocation evidence
- readers round-robin across 3 endpoints; the writer's endpoint policy is recorded separately
- physical rank-index arm yb-range-rank; an engine whose planner ignores the index is a different physical design, not a faster one
- YugabyteDB DIST/DEBUG storage counters (rows and index rows scanned) are parsed per cell; RPC-level counters are not aggregated
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
