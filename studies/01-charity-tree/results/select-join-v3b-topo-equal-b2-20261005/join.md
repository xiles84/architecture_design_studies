# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 30 cells, 30 valid, engine `yugabyte`, topology `single`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-v3b-topo-equal-b2-20261005`
- run tag: `run/study-01/select-join-v3b-topo-equal-b2-20261005`
- environment: `host-zenbook-ux5406sa`
- repo commit: `7243fee074bfb16fc5ad8b6afc1a2d7549080935`
- engine version: `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit`
- database engine image: `docker.io/yugabytedb/yugabyte:2025.2.6.0-b111` @ `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`
- benchmark client image: `localhost/charitybench-select:select-join-v3b-topo-equal-b2-20261005` @ `f3f1e10c44cbcb4c7b5bc56ed95a4f071682b6f7d33b5e0a7b688c4271eeef88`
- inputs digest: `963c6e51422ac24f7c1dec336c15de7c519e9687538c04e31666262ab2bf1ce6`
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
| 20 (0.01%) | 1.8796 | missing | missing | missing | missing | missing | 2.1018 | 87.0739 | 1.7972 | missing |
| 2000 (1%) | 9.6159 | missing | missing | missing | missing | missing | 9.5416 | 90.7287 | 9.8266 | missing |
| 180000 (90%) | 117.6850 | missing | missing | missing | missing | missing | 117.7177 | 110.7488 | 114.5476 | missing |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | cte_joined_outer vs join_where | 1.8796 | 2.1018 | 1.1182 | false |
| 20 | cte_joined_outer_mat vs join_where | 1.8796 | 87.0739 | 46.3262 | false |
| 20 | cte_joined_outer_notmat vs join_where | 1.8796 | 1.7972 | 0.9562 | false |
| 2000 | cte_joined_outer vs join_where | 9.6159 | 9.5416 | 0.9923 | false |
| 2000 | cte_joined_outer_mat vs join_where | 9.6159 | 90.7287 | 9.4353 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 9.6159 | 9.8266 | 1.0219 | false |
| 180000 | cte_joined_outer vs join_where | 117.6850 | 117.7177 | 1.0003 | false |
| 180000 | cte_joined_outer_mat vs join_where | 117.6850 | 110.7488 | 0.9411 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 117.6850 | 114.5476 | 0.9733 | false |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_joined_outer_mat vs cte_joined_outer | 2.1018 | 87.0739 | 41.4281 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 2.1018 | 1.7972 | 0.8551 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 9.5416 | 90.7287 | 9.5087 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 9.5416 | 9.8266 | 1.0299 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 117.7177 | 110.7488 | 0.9408 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 117.7177 | 114.5476 | 0.9731 | false | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 6.9667 | 1 | `7a4399f5355e7177` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 8.8254 | 1 | `5b456880e2c342ac` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | missing | | | | |
| default | 2 | 100% | missing | | | | |
| default | 4 | 1% | 8.7022 | 4 | `e598a67e88b8334f` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 9.7623 | 4 | `06f6b29197146f62` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | missing | | | | |
| default | 8 | 100% | missing | | | | |
| default | 16 | 1% | 10.6348 | 16 | `da59001908595408` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 10.1686 | 16 | `767476c005658d4f` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 8.9680 | 1 | `a21fb338a2bcfd02` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 8.3939 | 1 | `fbd0a43dc15da09e` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | missing | | | | |
| mat | 2 | 100% | missing | | | | |
| mat | 4 | 1% | 9.1844 | 4 | `9923e3efdc9ad044` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 9.4250 | 4 | `9885f1eed72a4dbd` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | missing | | | | |
| mat | 8 | 100% | missing | | | | |
| mat | 16 | 1% | 10.7166 | 16 | `877faf8da8050553` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 10.9798 | 16 | `937bacf5d4033196` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 7.7151 | 1 | `fbb4de7d30a068f1` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 9.1769 | 1 | `235ba57e8ae4cd79` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | missing | | | | |
| notmat | 2 | 100% | missing | | | | |
| notmat | 4 | 1% | 24.0662 | 4 | `ae73eb4e47c2b4e5` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 33.3407 | 4 | `67f9a1592e1c482b` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | missing | | | | |
| notmat | 8 | 100% | missing | | | | |
| notmat | 16 | 1% | 75.7946 | 16 | `6f5054cd6858e270` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 95.7041 | 16 | `c007655aca761b5c` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 1.2873 | false |
| notmat | 1 | 1% | 1.1074 | false |
| mat | 1 | 100% | 0.9511 | false |
| notmat | 1 | 100% | 1.0398 | false |
| mat | 4 | 1% | 1.0554 | false |
| notmat | 4 | 1% | 2.7655 | false |
| mat | 4 | 100% | 0.9654 | false |
| notmat | 4 | 100% | 3.4152 | false |
| mat | 16 | 1% | 1.0077 | false |
| notmat | 16 | 1% | 7.1270 | false |
| mat | 16 | 100% | 1.0798 | false |
| notmat | 16 | 100% | 9.4117 | false |

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

- client: nr_periods 312, nr_throttled 130, throttled_usec 5468048 | quota 100000 100000
, memory 2147483648

- yb-single: nr_periods 380, nr_throttled 0, throttled_usec 0 | quota 600000 100000, memory 9663676416

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
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 110.6985 | 117.7177 | 124.5085 | 5 | 8.5 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 1.9693 | 2.1018 | 2.6867 | 5 | 475.8 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 9.0285 | 9.5416 | 9.8687 | 5 | 104.8 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 107.8044 | 110.7488 | 132.4987 | 5 | 9.0 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 86.4673 | 87.0739 | 144.2952 | 5 | 11.5 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 86.9920 | 90.7287 | 101.5776 | 5 | 11.0 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 112.0619 | 114.5476 | 119.6092 | 5 | 8.7 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 1.6424 | 1.7972 | 2.7211 | 5 | 556.4 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 9.1213 | 9.8266 | 10.7259 | 5 | 101.8 |
| family=bag;stmt=join_where;threshold=180000 | valid | 116.3129 | 117.6850 | 123.8079 | 5 | 8.5 |
| family=bag;stmt=join_where;threshold=20 | valid | 1.7930 | 1.8796 | 2.7256 | 5 | 532.0 |
| family=bag;stmt=join_where;threshold=2000 | valid | 9.4503 | 9.6159 | 9.8337 | 5 | 104.0 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 10.5780 | 10.6348 | 10.8841 | 5 | 94.0 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 10.0456 | 10.1686 | 10.5759 | 5 | 98.3 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 6.2153 | 6.9667 | 7.1624 | 5 | 143.5 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 8.4974 | 8.8254 | 9.3190 | 5 | 113.3 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 8.2182 | 8.7022 | 9.8694 | 5 | 114.9 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 9.5733 | 9.7623 | 10.4089 | 5 | 102.4 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 10.4110 | 10.7166 | 11.4097 | 5 | 93.3 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 10.7940 | 10.9798 | 11.1709 | 5 | 91.1 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 8.5666 | 8.9680 | 9.4054 | 5 | 111.5 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 8.1302 | 8.3939 | 9.1677 | 5 | 119.1 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 8.6344 | 9.1844 | 9.4466 | 5 | 108.9 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 9.0341 | 9.4250 | 9.8242 | 5 | 106.1 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 73.4942 | 75.7946 | 85.7104 | 5 | 13.2 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 86.1927 | 95.7041 | 96.1074 | 5 | 10.4 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 7.2448 | 7.7151 | 7.9784 | 5 | 129.6 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 8.7941 | 9.1769 | 9.5577 | 5 | 109.0 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 22.8911 | 24.0662 | 26.1037 | 5 | 41.6 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 32.8145 | 33.3407 | 35.1733 | 5 | 30.0 |

## Matched blended workload

- phase: `blended`; target offered rate: 10 updates/s
- offered 155, committed 155, failed 0 over 15.438 s
- offered rate 10.040/s, completed rate 10.040/s
- writer operations per endpoint: [155]
- reader operations overlapping a committed write: 72 of 210
- reader retries of transient conflicts: 0
- per key: 155 keys, min 1, max 1 per key
- ledger: 155 keys reconciled, 0 mismatch
- update latency: avg 3.667 ms, min 1.597, max 88.878 over 155

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42310428`, 5 interleaved rounds
  1. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  2. family=reuse;policy=default;refs=16;consumer=1%
  3. family=bag;stmt=cte_joined_outer;threshold=2000
  4. family=bag;stmt=cte_joined_outer;threshold=20
  5. family=reuse;policy=mat;refs=4;consumer=100%
  6. family=reuse;policy=notmat;refs=16;consumer=1%
  7. family=reuse;policy=mat;refs=16;consumer=1%
  8. family=reuse;policy=mat;refs=1;consumer=1%
  9. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  10. family=reuse;policy=mat;refs=16;consumer=100%
  11. family=reuse;policy=mat;refs=4;consumer=1%
  12. family=reuse;policy=mat;refs=1;consumer=100%
  13. family=bag;stmt=join_where;threshold=180000
  14. family=reuse;policy=notmat;refs=4;consumer=100%
  15. family=bag;stmt=cte_joined_outer;threshold=180000
  16. family=reuse;policy=default;refs=4;consumer=100%
  17. family=reuse;policy=notmat;refs=4;consumer=1%
  18. family=reuse;policy=default;refs=1;consumer=1%
  19. family=bag;stmt=join_where;threshold=20
  20. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  21. family=bag;stmt=cte_joined_outer_mat;threshold=20
  22. family=reuse;policy=notmat;refs=1;consumer=100%
  23. family=reuse;policy=notmat;refs=16;consumer=100%
  24. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  25. family=reuse;policy=notmat;refs=1;consumer=1%
  26. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  27. family=reuse;policy=default;refs=1;consumer=100%
  28. family=reuse;policy=default;refs=4;consumer=1%
  29. family=reuse;policy=default;refs=16;consumer=100%
  30. family=bag;stmt=join_where;threshold=2000

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
| family=bag;stmt=cte_joined_outer;threshold=180000 | j05_cte_joined_outer | valid | `754c664ee88d500d` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20 | j05_cte_joined_outer | valid | `35c46121749ea9dc` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=2000 | j05_cte_joined_outer | valid | `144b73e36e2d0395` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | j05m_cte_joined_outer_mat | valid | `3dace3b85179b355` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | j05m_cte_joined_outer_mat | valid | `4340ec6b0b48765c` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | j05m_cte_joined_outer_mat | valid | `e9d1f8a4f5721e70` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | j05n_cte_joined_outer_notmat | valid | `a74453d93079b01f` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | j05n_cte_joined_outer_notmat | valid | `30515c0f57076ff9` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | j05n_cte_joined_outer_notmat | valid | `ea8ce17457a7185a` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=join_where;threshold=180000 | j01_join_where | valid | `11aecd28fe804e37` |  |
| family=bag;stmt=join_where;threshold=20 | j01_join_where | valid | `5381c61805fcaee6` |  |
| family=bag;stmt=join_where;threshold=2000 | j01_join_where | valid | `d1a8f106b863af81` |  |
| family=reuse;policy=default;refs=16;consumer=1% | r_default_n16_c1 | valid | `da59001908595408` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=16;consumer=100% | r_default_n16_c100 | valid | `767476c005658d4f` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=1;consumer=1% | r_default_n1_c1 | valid | `7a4399f5355e7177` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=1;consumer=100% | r_default_n1_c100 | valid | `5b456880e2c342ac` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=4;consumer=1% | r_default_n4_c1 | valid | `e598a67e88b8334f` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=100% | r_default_n4_c100 | valid | `06f6b29197146f62` | policy=default; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=1% | r_mat_n16_c1 | valid | `877faf8da8050553` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=100% | r_mat_n16_c100 | valid | `937bacf5d4033196` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=1% | r_mat_n1_c1 | valid | `a21fb338a2bcfd02` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=100% | r_mat_n1_c100 | valid | `fbd0a43dc15da09e` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=1% | r_mat_n4_c1 | valid | `9923e3efdc9ad044` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=100% | r_mat_n4_c100 | valid | `9885f1eed72a4dbd` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=notmat;refs=16;consumer=1% | r_notmat_n16_c1 | valid | `6f5054cd6858e270` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `c007655aca761b5c` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `fbb4de7d30a068f1` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `235ba57e8ae4cd79` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `ae73eb4e47c2b4e5` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `67f9a1592e1c482b` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, DIST, DEBUG, FORMAT TEXT` supported on yugabyte
- one host: this run used node budget 6 CPU 9g and client budget 1 CPU 2g, so budget comparisons are separate classes here, not independent hardware
- one client process with one fixed maintenance writer at 10 updates/s; readers are sequential; the reported sample latency includes retry and backoff time
- physical rank-index arm yb-range-rank; an engine whose planner ignores the index is a different physical design, not a faster one
- YugabyteDB DIST/DEBUG storage counters (rows and index rows scanned) are parsed per cell; RPC-level counters are not aggregated
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
