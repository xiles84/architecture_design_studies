# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 30 cells, 30 valid, engine `yugabyte`, topology `cluster3`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-v4-smoke-distributed-20261005`
- run tag: `run/study-01/select-join-v4-smoke-distributed-20261005`
- environment: `host-zenbook-ux5406sa`
- repo commit: `8a1dd7680e0825069a27bfa3d1d60e98f10bd678`
- engine version: `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit`
- database engine image: `docker.io/yugabytedb/yugabyte:2025.2.6.0-b111` @ `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`
- benchmark client image: `localhost/charitybench-select:select-join-v4-smoke-distributed-20261005` @ `3bd8431c09af70be6c0028375e402c8929c8524cf1af65db89a6a36118b1718e`
- inputs digest: `bbea85e7f30117a004eeff66ba462be3e2f442b46d3220b8d0d3b7db4202b754`
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
- setting `placement`: distributed
- setting `prepared_plan_policy`: pgx simple protocol; no prepared statements
- setting `scale`: cte-domain
- setting `selectivity_mapping`: rank <= threshold; rank is a deterministic shuffle of 1..fact_count
- setting `shards_per_tserver`: ysql_num_shards_per_tserver= yb_num_shards_per_tserver=
- setting `statement_budget_ms`: 120000
- setting `statistics`: ANALYZE after deterministic load
- setting `thresholds`: [20 2000 180000]
- setting `topology`: cluster3
- setting `trials_per_cell`: 1
- setting `whole_run_count`: 1; no replicated winner in this file
- setting `writer_keys`: 1000
- setting `writer_target_rate_per_sec`: 10

### Endpoint distribution

A cluster arm must rotate its reader across nodes; a single seed endpoint would not be balanced cluster access.

| endpoint | operations | connections at end |
|---|---|---|
| yb-n1 | 128 | 1 |
| yb-n2 | 128 | 1 |
| yb-n3 | 128 | 1 |

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
| 20 (0.01%) | 2.7062 | missing | missing | missing | missing | missing | 2.6650 | 61.9694 | 2.8515 | missing |
| 2000 (1%) | 7.1708 | missing | missing | missing | missing | missing | 7.5343 | 62.4942 | 6.7450 | missing |
| 180000 (90%) | 73.4014 | missing | missing | missing | missing | missing | 76.5715 | 93.9451 | 74.6813 | missing |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | cte_joined_outer vs join_where | 2.7062 | 2.6650 | 0.9848 | false |
| 20 | cte_joined_outer_mat vs join_where | 2.7062 | 61.9694 | 22.8993 | false |
| 20 | cte_joined_outer_notmat vs join_where | 2.7062 | 2.8515 | 1.0537 | false |
| 2000 | cte_joined_outer vs join_where | 7.1708 | 7.5343 | 1.0507 | false |
| 2000 | cte_joined_outer_mat vs join_where | 7.1708 | 62.4942 | 8.7151 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 7.1708 | 6.7450 | 0.9406 | false |
| 180000 | cte_joined_outer vs join_where | 73.4014 | 76.5715 | 1.0432 | false |
| 180000 | cte_joined_outer_mat vs join_where | 73.4014 | 93.9451 | 1.2799 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 73.4014 | 74.6813 | 1.0174 | false |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_joined_outer_mat vs cte_joined_outer | 2.6650 | 61.9694 | 23.2529 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 2.6650 | 2.8515 | 1.0700 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 7.5343 | 62.4942 | 8.2946 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 7.5343 | 6.7450 | 0.8952 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 76.5715 | 93.9451 | 1.2269 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 76.5715 | 74.6813 | 0.9753 | false | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 4.1497 | 1 | `4c601ac25833abd0` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 5.4958 | 1 | `cd6a7f5990744c4e` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | missing | | | | |
| default | 2 | 100% | missing | | | | |
| default | 4 | 1% | 6.2653 | 4 | `9b3c60be999b0f5a` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 6.0358 | 4 | `47cebc9ec0f165fd` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | missing | | | | |
| default | 8 | 100% | missing | | | | |
| default | 16 | 1% | 8.2896 | 16 | `0758894389b9fb88` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 7.6516 | 16 | `aa7fc410fba75443` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 5.5467 | 1 | `8cb6072b8296a0d7` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 5.9079 | 1 | `5806600491dbc8dd` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | missing | | | | |
| mat | 2 | 100% | missing | | | | |
| mat | 4 | 1% | 6.3481 | 4 | `2804e1412a5e28ad` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 6.1049 | 4 | `2d329fa705582f7b` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | missing | | | | |
| mat | 8 | 100% | missing | | | | |
| mat | 16 | 1% | 7.5534 | 16 | `6a47c45f7b9f7dbe` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 7.3389 | 16 | `6a9ad2f84f171d7c` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 5.4496 | 1 | `93f1f20e0ff14ef2` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 6.2701 | 1 | `4d7817106a747cc6` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | missing | | | | |
| notmat | 2 | 100% | missing | | | | |
| notmat | 4 | 1% | 14.6676 | 4 | `126dc25c0a35a0ac` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 23.0565 | 4 | `48227a0a06dea92b` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | missing | | | | |
| notmat | 8 | 100% | missing | | | | |
| notmat | 16 | 1% | 85.1517 | 16 | `2f7d6c22560db8fb` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 82.5966 | 16 | `ceb415c9ceeb6763` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 1.3367 | false |
| notmat | 1 | 1% | 1.3133 | false |
| mat | 1 | 100% | 1.0750 | false |
| notmat | 1 | 100% | 1.1409 | false |
| mat | 4 | 1% | 1.0132 | false |
| notmat | 4 | 1% | 2.3411 | false |
| mat | 4 | 100% | 1.0114 | false |
| notmat | 4 | 100% | 3.8199 | false |
| mat | 16 | 1% | 0.9112 | false |
| notmat | 16 | 1% | 10.2721 | false |
| mat | 16 | 100% | 0.9591 | false |
| notmat | 16 | 100% | 10.7947 | false |

## Physical arm

- `charity_pkey`: `CREATE UNIQUE INDEX charity_pkey ON selectjoin_domain.charity USING lsm (charity_id HASH)`
- `donation_person_idx`: `CREATE INDEX donation_person_idx ON selectjoin_domain.donation USING lsm (person_id HASH)`
- `donation_pkey`: `CREATE UNIQUE INDEX donation_pkey ON selectjoin_domain.donation USING lsm (donation_id HASH)`
- `donation_rank_idx`: `CREATE INDEX donation_rank_idx ON selectjoin_domain.donation USING lsm (rank ASC)`
- `person_pkey`: `CREATE UNIQUE INDEX person_pkey ON selectjoin_domain.person USING lsm (person_id HASH)`

- index arm `yb-range-rank`; placement `distributed`; node budget `2 CPU 3g`; client budget `1 CPU 2g`

| table | tablets | hash key columns | colocated | tablegroup | colocation id |
|---|---|---|---|---|---|
| selectjoin_domain.charity | 3 | 1 | false | unknown | unknown |
| selectjoin_domain.person | 3 | 1 | false | unknown | unknown |
| selectjoin_domain.donation | 3 | 1 | false | unknown | unknown |
| selectjoin_domain.donation_pkey | 0 | 0 | false | unknown | unknown |
| selectjoin_domain.donation_rank_idx | 3 | 0 | false | unknown | unknown |
| selectjoin_domain.donation_person_idx | 3 | 1 | false | unknown | unknown |

Physical-arm gate (default join must use the rank index at its three most selective thresholds):

- threshold 20: used
- threshold 2000: used

Tablet leaders per relation (placement evidence; the runner read these through yb-admin):

| relation | leader nodes |
|---|---|
| charity | [0x0000, [0x5555, [0xAAAA |
| person | [0x0000, [0x5555, [0xAAAA |
| donation | [0x0000, [0x5555, [0xAAAA |
| donation_rank_idx | [<start>, [DocKey([] |
| donation_person_idx | [0x0000, [0x5555, [0xAAAA |

Effective sharding flag (from the engine): `ysql_num_shards_per_tserver= yb_num_shards_per_tserver=`. A relation with fewer tablets than this setting was created without an explicit split; the per-table `num_tablets` above is the authority.

## Throttling (cgroup `cpu.stat` deltas)

- client: nr_periods 208, nr_throttled 73, throttled_usec 2674899 | quota 100000 100000
, memory 2147483648

- yb-n1: nr_periods 540, nr_throttled 12, throttled_usec 1083548 | quota 200000 100000, memory 3221225472
- yb-n2: nr_periods 530, nr_throttled 4, throttled_usec 561014 | quota 200000 100000, memory 3221225472
- yb-n3: nr_periods 528, nr_throttled 1, throttled_usec 26342 | quota 200000 100000, memory 3221225472

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
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 76.5715 | 76.5715 | 76.5715 | 1 | 13.1 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 2.6650 | 2.6650 | 2.6650 | 1 | 375.2 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 7.5343 | 7.5343 | 7.5343 | 1 | 132.7 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 93.9451 | 93.9451 | 93.9451 | 1 | 10.6 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 61.9694 | 61.9694 | 61.9694 | 1 | 16.1 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 62.4942 | 62.4942 | 62.4942 | 1 | 16.0 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 74.6813 | 74.6813 | 74.6813 | 1 | 13.4 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 2.8515 | 2.8515 | 2.8515 | 1 | 350.7 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 6.7450 | 6.7450 | 6.7450 | 1 | 148.3 |
| family=bag;stmt=join_where;threshold=180000 | valid | 73.4014 | 73.4014 | 73.4014 | 1 | 13.6 |
| family=bag;stmt=join_where;threshold=20 | valid | 2.7062 | 2.7062 | 2.7062 | 1 | 369.5 |
| family=bag;stmt=join_where;threshold=2000 | valid | 7.1708 | 7.1708 | 7.1708 | 1 | 139.5 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 8.2896 | 8.2896 | 8.2896 | 1 | 120.6 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 7.6516 | 7.6516 | 7.6516 | 1 | 130.7 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 4.1497 | 4.1497 | 4.1497 | 1 | 241.0 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 5.4958 | 5.4958 | 5.4958 | 1 | 182.0 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 6.2653 | 6.2653 | 6.2653 | 1 | 159.6 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 6.0358 | 6.0358 | 6.0358 | 1 | 165.7 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 7.5534 | 7.5534 | 7.5534 | 1 | 132.4 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 7.3389 | 7.3389 | 7.3389 | 1 | 136.3 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 5.5467 | 5.5467 | 5.5467 | 1 | 180.3 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 5.9079 | 5.9079 | 5.9079 | 1 | 169.3 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 6.3481 | 6.3481 | 6.3481 | 1 | 157.5 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 6.1049 | 6.1049 | 6.1049 | 1 | 163.8 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 85.1517 | 85.1517 | 85.1517 | 1 | 11.7 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 82.5966 | 82.5966 | 82.5966 | 1 | 12.1 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 5.4496 | 5.4496 | 5.4496 | 1 | 183.5 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 6.2701 | 6.2701 | 6.2701 | 1 | 159.5 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 14.6676 | 14.6676 | 14.6676 | 1 | 68.2 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 23.0565 | 23.0565 | 23.0565 | 1 | 43.4 |

## Matched blended workload

- phase: `isolated`; target offered rate: 10 updates/s
- offered 0, committed 0, failed 0 over 0.000 s
- offered rate 0.000/s, completed rate 0.000/s
- per key: no keys touched
- ledger: 0 keys reconciled, 0 mismatch

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42448382`, 1 interleaved rounds
  1. family=reuse;policy=default;refs=4;consumer=1%
  2. family=reuse;policy=default;refs=1;consumer=1%
  3. family=reuse;policy=default;refs=4;consumer=100%
  4. family=reuse;policy=default;refs=16;consumer=100%
  5. family=reuse;policy=mat;refs=1;consumer=100%
  6. family=reuse;policy=mat;refs=1;consumer=1%
  7. family=reuse;policy=notmat;refs=1;consumer=100%
  8. family=bag;stmt=join_where;threshold=180000
  9. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  10. family=reuse;policy=default;refs=16;consumer=1%
  11. family=reuse;policy=notmat;refs=1;consumer=1%
  12. family=bag;stmt=cte_joined_outer;threshold=20
  13. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  14. family=reuse;policy=mat;refs=16;consumer=1%
  15. family=bag;stmt=cte_joined_outer_mat;threshold=20
  16. family=reuse;policy=mat;refs=16;consumer=100%
  17. family=reuse;policy=mat;refs=4;consumer=100%
  18. family=reuse;policy=default;refs=1;consumer=100%
  19. family=reuse;policy=notmat;refs=4;consumer=1%
  20. family=reuse;policy=notmat;refs=4;consumer=100%
  21. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  22. family=bag;stmt=cte_joined_outer;threshold=2000
  23. family=reuse;policy=notmat;refs=16;consumer=1%
  24. family=bag;stmt=cte_joined_outer;threshold=180000
  25. family=bag;stmt=join_where;threshold=2000
  26. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  27. family=reuse;policy=mat;refs=4;consumer=1%
  28. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  29. family=bag;stmt=join_where;threshold=20
  30. family=reuse;policy=notmat;refs=16;consumer=100%

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
| family=bag;stmt=cte_joined_outer;threshold=180000 | j05_cte_joined_outer | valid | `100c6d79f0ab6ad6` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20 | j05_cte_joined_outer | valid | `5b4518f73bc7c399` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=2000 | j05_cte_joined_outer | valid | `c2e391c21fa0e62f` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | j05m_cte_joined_outer_mat | valid | `a715568e9e72ebaf` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | j05m_cte_joined_outer_mat | valid | `ebb85c235e99ae5a` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | j05m_cte_joined_outer_mat | valid | `c4cbf185d5bb326d` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | j05n_cte_joined_outer_notmat | valid | `f9eb82ceec6a6ca7` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | j05n_cte_joined_outer_notmat | valid | `a889c1cbcab9633b` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | j05n_cte_joined_outer_notmat | valid | `b10552bb9ec41589` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=join_where;threshold=180000 | j01_join_where | valid | `125e6fc85e2c120e` |  |
| family=bag;stmt=join_where;threshold=20 | j01_join_where | valid | `9f1dcb09e4230166` |  |
| family=bag;stmt=join_where;threshold=2000 | j01_join_where | valid | `df2b0f4779ab64e9` |  |
| family=reuse;policy=default;refs=16;consumer=1% | r_default_n16_c1 | valid | `0758894389b9fb88` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=16;consumer=100% | r_default_n16_c100 | valid | `aa7fc410fba75443` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=1;consumer=1% | r_default_n1_c1 | valid | `4c601ac25833abd0` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=1;consumer=100% | r_default_n1_c100 | valid | `cd6a7f5990744c4e` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=4;consumer=1% | r_default_n4_c1 | valid | `9b3c60be999b0f5a` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=100% | r_default_n4_c100 | valid | `47cebc9ec0f165fd` | policy=default; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=1% | r_mat_n16_c1 | valid | `6a47c45f7b9f7dbe` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=100% | r_mat_n16_c100 | valid | `6a9ad2f84f171d7c` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=1% | r_mat_n1_c1 | valid | `8cb6072b8296a0d7` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=100% | r_mat_n1_c100 | valid | `5806600491dbc8dd` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=1% | r_mat_n4_c1 | valid | `2804e1412a5e28ad` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=100% | r_mat_n4_c100 | valid | `2d329fa705582f7b` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=notmat;refs=16;consumer=1% | r_notmat_n16_c1 | valid | `2f7d6c22560db8fb` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `ceb415c9ceeb6763` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `93f1f20e0ff14ef2` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `4d7817106a747cc6` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `126dc25c0a35a0ac` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `48227a0a06dea92b` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, DIST, DEBUG, FORMAT TEXT` supported on yugabyte
- distributed placement gate: three tablets confirmed by the engine, leader spread confirmed by the runner
- one host: this run used node budget 2 CPU 3g and client budget 1 CPU 2g, so budget comparisons are separate classes here, not independent hardware
- isolated phase: no concurrent writers in this run; the blended read/write phase is a separate labelled run
- cluster nodes are containers on the same host, so these are not network or independent-hardware measurements; container colocation is not data-colocation evidence
- readers round-robin across 3 endpoints; the writer's endpoint policy is recorded separately
- physical rank-index arm yb-range-rank; an engine whose planner ignores the index is a different physical design, not a faster one
- YugabyteDB DIST/DEBUG storage counters (rows and index rows scanned) are parsed per cell; RPC-level counters are not aggregated
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
