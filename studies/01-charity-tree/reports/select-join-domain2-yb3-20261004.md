# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 91 cells, 91 valid, engine `yugabyte`, topology `single`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-domain2-yb3-20261004`
- run tag: `run/study-01/select-join-domain2-yb3-20261004`
- environment: `host-zenbook-ux5406sa`
- repo commit: `445fbc1646d58074c45cbaea57a83234132adc34`
- engine version: `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit`
- database engine image: `docker.io/yugabytedb/yugabyte:2025.2.6.0-b111` @ `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`
- benchmark client image: `localhost/charitybench-select:select-join-domain2-yb3-20261004` @ `79bc753ff7e57e102982a323939c43e68e9b3664887683d3a5d5bf072650f4c1`
- inputs digest: `e497ebb27c95425f238f8a12e2529f92a0346468baaa0a32449f38cef0c479cc`
- dataset fingerprint: `fcf26e08f377ec51`
- client CPU quota: `100000 100000`, memory: `2147483648`
- setting `client_workers`: 1 sequential reader; 0 writers
- setting `connection_pool`: 2 per endpoint = reader/setup + reserve
- setting `customers`: 2000
- setting `dataset`: study01-person-donation-charity
- setting `distributed_rpc`: not parsed; raw EXPLAIN retained (coverage gap)
- setting `endpoint_mode`: single endpoint
- setting `endpoints`: yb-single
- setting `fact_count`: 200000
- setting `index_set`: child rank index and parent-key index, primary keys; no per-spelling index
- setting `phase`: isolated
- setting `placement`: engine default
- setting `prepared_plan_policy`: pgx simple protocol; no prepared statements
- setting `scale`: cte-domain
- setting `selectivity_mapping`: rank <= threshold; rank is a deterministic shuffle of 1..fact_count
- setting `statement_budget_ms`: 120000
- setting `statistics`: ANALYZE after deterministic load
- setting `thresholds`: [20 200 2000 20000 100000 180000]
- setting `topology`: single
- setting `trials_per_cell`: 5
- setting `whole_run_count`: 1; no replicated winner in this file
- setting `writer_keys`: 1000
- setting `writer_target_rate_per_sec`: 10

### Endpoint distribution

A cluster arm must rotate its reader across nodes; a single seed endpoint would not be balanced cluster access.

| endpoint | operations | connections at end |
|---|---|---|
| yb-single | 914 | 1 |

## Fixture and selectivity mapping

- arm: `study01-person-donation-charity`
- fact rows: 200000
- parent rows: 2000
- rank thresholds: [20 200 2000 20000 100000 180000]

The predicate is `rank <= threshold`, where `rank` is a deterministic shuffle of 1..200000; `rank` is a labelled proxy for the teaching selective predicate, not the same question. The `cte`/`cte-small` arm is the synthetic customers/facts microfixture; the `cte-domain` arm is Study 01's person/donation fixture with the person's charity projected and the views binding it to this catalogue.

## Per-engine feasibility evidence

- fact rows 200000, customers 2000, thresholds [20 200 2000 20000 100000 180000]
- expensive reuse token supported by the engine: yes

| threshold | rows with rank <= threshold |
|---|---|
| 100000 | 100000 |
| 180000 | 180000 |
| 20 | 20 |
| 200 | 200 |
| 2000 | 2000 |
| 20000 | 20000 |

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
| 20 (0.01%) | 56.3768 | 57.8984 | 56.7067 | 60.8675 | 53.8213 | 55.9078 | 56.2751 | 88.7820 | 58.0128 | 55.0820 |
| 200 (0.1%) | 60.7198 | 55.7144 | 54.8072 | 56.4916 | 57.1962 | 58.7585 | 56.5626 | 85.1143 | 57.7281 | 54.6742 |
| 2000 (1%) | 57.4667 | 58.9139 | 57.8433 | 57.2131 | 57.6344 | 55.1731 | 57.4451 | 88.5428 | 56.9881 | 59.8192 |
| 20000 (10%) | 67.5504 | 69.1130 | 66.6618 | 66.9959 | 73.1248 | 70.2809 | 70.0080 | 93.1909 | 69.2583 | 70.3956 |
| 100000 (50%) | 88.6960 | 93.1732 | 91.7815 | 89.2507 | 96.0385 | 92.0216 | 89.9376 | 98.7561 | 93.4837 | 91.9948 |
| 180000 (90%) | 114.4575 | 112.2459 | 110.3726 | 109.8643 | 112.7105 | 108.3784 | 114.9823 | 107.1639 | 114.0428 | 113.7601 |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | join_on vs join_where | 56.3768 | 57.8984 | 1.0270 | false |
| 20 | filtered_derived vs join_where | 56.3768 | 56.7067 | 1.0059 | false |
| 20 | cte_filtered vs join_where | 56.3768 | 60.8675 | 1.0797 | false |
| 20 | cte_filtered_mat vs join_where | 56.3768 | 53.8213 | 0.9547 | false |
| 20 | cte_filtered_notmat vs join_where | 56.3768 | 55.9078 | 0.9917 | false |
| 20 | cte_joined_outer vs join_where | 56.3768 | 56.2751 | 0.9982 | false |
| 20 | cte_joined_outer_mat vs join_where | 56.3768 | 88.7820 | 1.5748 | false |
| 20 | cte_joined_outer_notmat vs join_where | 56.3768 | 58.0128 | 1.0290 | false |
| 20 | comma_join vs join_where | 56.3768 | 55.0820 | 0.9770 | false |
| 200 | join_on vs join_where | 60.7198 | 55.7144 | 0.9176 | false |
| 200 | filtered_derived vs join_where | 60.7198 | 54.8072 | 0.9026 | false |
| 200 | cte_filtered vs join_where | 60.7198 | 56.4916 | 0.9304 | false |
| 200 | cte_filtered_mat vs join_where | 60.7198 | 57.1962 | 0.9420 | false |
| 200 | cte_filtered_notmat vs join_where | 60.7198 | 58.7585 | 0.9677 | false |
| 200 | cte_joined_outer vs join_where | 60.7198 | 56.5626 | 0.9315 | false |
| 200 | cte_joined_outer_mat vs join_where | 60.7198 | 85.1143 | 1.4018 | false |
| 200 | cte_joined_outer_notmat vs join_where | 60.7198 | 57.7281 | 0.9507 | false |
| 200 | comma_join vs join_where | 60.7198 | 54.6742 | 0.9004 | false |
| 2000 | join_on vs join_where | 57.4667 | 58.9139 | 1.0252 | false |
| 2000 | filtered_derived vs join_where | 57.4667 | 57.8433 | 1.0066 | false |
| 2000 | cte_filtered vs join_where | 57.4667 | 57.2131 | 0.9956 | false |
| 2000 | cte_filtered_mat vs join_where | 57.4667 | 57.6344 | 1.0029 | false |
| 2000 | cte_filtered_notmat vs join_where | 57.4667 | 55.1731 | 0.9601 | false |
| 2000 | cte_joined_outer vs join_where | 57.4667 | 57.4451 | 0.9996 | false |
| 2000 | cte_joined_outer_mat vs join_where | 57.4667 | 88.5428 | 1.5408 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 57.4667 | 56.9881 | 0.9917 | false |
| 2000 | comma_join vs join_where | 57.4667 | 59.8192 | 1.0409 | false |
| 20000 | join_on vs join_where | 67.5504 | 69.1130 | 1.0231 | false |
| 20000 | filtered_derived vs join_where | 67.5504 | 66.6618 | 0.9868 | false |
| 20000 | cte_filtered vs join_where | 67.5504 | 66.9959 | 0.9918 | false |
| 20000 | cte_filtered_mat vs join_where | 67.5504 | 73.1248 | 1.0825 | false |
| 20000 | cte_filtered_notmat vs join_where | 67.5504 | 70.2809 | 1.0404 | false |
| 20000 | cte_joined_outer vs join_where | 67.5504 | 70.0080 | 1.0364 | false |
| 20000 | cte_joined_outer_mat vs join_where | 67.5504 | 93.1909 | 1.3796 | false |
| 20000 | cte_joined_outer_notmat vs join_where | 67.5504 | 69.2583 | 1.0253 | false |
| 20000 | comma_join vs join_where | 67.5504 | 70.3956 | 1.0421 | false |
| 100000 | join_on vs join_where | 88.6960 | 93.1732 | 1.0505 | false |
| 100000 | filtered_derived vs join_where | 88.6960 | 91.7815 | 1.0348 | false |
| 100000 | cte_filtered vs join_where | 88.6960 | 89.2507 | 1.0063 | false |
| 100000 | cte_filtered_mat vs join_where | 88.6960 | 96.0385 | 1.0828 | false |
| 100000 | cte_filtered_notmat vs join_where | 88.6960 | 92.0216 | 1.0375 | false |
| 100000 | cte_joined_outer vs join_where | 88.6960 | 89.9376 | 1.0140 | false |
| 100000 | cte_joined_outer_mat vs join_where | 88.6960 | 98.7561 | 1.1134 | false |
| 100000 | cte_joined_outer_notmat vs join_where | 88.6960 | 93.4837 | 1.0540 | false |
| 100000 | comma_join vs join_where | 88.6960 | 91.9948 | 1.0372 | false |
| 180000 | join_on vs join_where | 114.4575 | 112.2459 | 0.9807 | false |
| 180000 | filtered_derived vs join_where | 114.4575 | 110.3726 | 0.9643 | false |
| 180000 | cte_filtered vs join_where | 114.4575 | 109.8643 | 0.9599 | false |
| 180000 | cte_filtered_mat vs join_where | 114.4575 | 112.7105 | 0.9847 | false |
| 180000 | cte_filtered_notmat vs join_where | 114.4575 | 108.3784 | 0.9469 | false |
| 180000 | cte_joined_outer vs join_where | 114.4575 | 114.9823 | 1.0046 | false |
| 180000 | cte_joined_outer_mat vs join_where | 114.4575 | 107.1639 | 0.9363 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 114.4575 | 114.0428 | 0.9964 | false |
| 180000 | comma_join vs join_where | 114.4575 | 113.7601 | 0.9939 | false |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_filtered_mat vs cte_filtered | 60.8675 | 53.8213 | 0.8842 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_filtered_notmat vs cte_filtered | 60.8675 | 55.9078 | 0.9185 | false | policy=not_materialized; observed_cte_scan=false |
| 20 | cte_joined_outer_mat vs cte_joined_outer | 56.2751 | 88.7820 | 1.5776 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 56.2751 | 58.0128 | 1.0309 | false | policy=not_materialized; observed_cte_scan=false |
| 200 | cte_filtered_mat vs cte_filtered | 56.4916 | 57.1962 | 1.0125 | false | policy=materialized; observed_cte_scan=true |
| 200 | cte_filtered_notmat vs cte_filtered | 56.4916 | 58.7585 | 1.0401 | false | policy=not_materialized; observed_cte_scan=false |
| 200 | cte_joined_outer_mat vs cte_joined_outer | 56.5626 | 85.1143 | 1.5048 | false | policy=materialized; observed_cte_scan=true |
| 200 | cte_joined_outer_notmat vs cte_joined_outer | 56.5626 | 57.7281 | 1.0206 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_filtered_mat vs cte_filtered | 57.2131 | 57.6344 | 1.0074 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_filtered_notmat vs cte_filtered | 57.2131 | 55.1731 | 0.9643 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 57.4451 | 88.5428 | 1.5413 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 57.4451 | 56.9881 | 0.9920 | false | policy=not_materialized; observed_cte_scan=false |
| 20000 | cte_filtered_mat vs cte_filtered | 66.9959 | 73.1248 | 1.0915 | false | policy=materialized; observed_cte_scan=true |
| 20000 | cte_filtered_notmat vs cte_filtered | 66.9959 | 70.2809 | 1.0490 | false | policy=not_materialized; observed_cte_scan=false |
| 20000 | cte_joined_outer_mat vs cte_joined_outer | 70.0080 | 93.1909 | 1.3311 | false | policy=materialized; observed_cte_scan=true |
| 20000 | cte_joined_outer_notmat vs cte_joined_outer | 70.0080 | 69.2583 | 0.9893 | false | policy=not_materialized; observed_cte_scan=false |
| 100000 | cte_filtered_mat vs cte_filtered | 89.2507 | 96.0385 | 1.0761 | false | policy=materialized; observed_cte_scan=true |
| 100000 | cte_filtered_notmat vs cte_filtered | 89.2507 | 92.0216 | 1.0310 | false | policy=not_materialized; observed_cte_scan=false |
| 100000 | cte_joined_outer_mat vs cte_joined_outer | 89.9376 | 98.7561 | 1.0981 | false | policy=materialized; observed_cte_scan=true |
| 100000 | cte_joined_outer_notmat vs cte_joined_outer | 89.9376 | 93.4837 | 1.0394 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_filtered_mat vs cte_filtered | 109.8643 | 112.7105 | 1.0259 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_filtered_notmat vs cte_filtered | 109.8643 | 108.3784 | 0.9865 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 114.9823 | 107.1639 | 0.9320 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 114.9823 | 114.0428 | 0.9918 | false | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 55.8683 | 1 | `88832648c1133a9b` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 60.0464 | 1 | `1a133a13b7051a2b` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | 58.3900 | 2 | `2b6e36da3d1af17f` | policy=default; observed_cte_scan=true | valid |
| default | 2 | 100% | 54.9913 | 2 | `179dbfd898a87226` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 1% | 55.6060 | 4 | `e300101f501895e9` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 58.6463 | 4 | `ba9898a9027d81dc` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | 55.4290 | 8 | `fa27563ccff22d0e` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 100% | 54.4486 | 8 | `dece197c8f297db0` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 1% | 56.7183 | 16 | `72e2200f58a6b85e` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 56.2788 | 16 | `3d6e2be0b00ff813` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 55.2001 | 1 | `dfa7dd7ef8545674` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 56.7918 | 1 | `c4c72f423291626f` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | 56.6383 | 2 | `f5847cdbdbefacbb` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 100% | 57.4015 | 2 | `9af650e3aaf44181` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 1% | 56.9978 | 4 | `798f8430f2e62ba9` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 56.3408 | 4 | `7774211ba86ffb34` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | 58.6201 | 8 | `4e687ec3ffcd47c3` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 100% | 59.9441 | 8 | `e410ab4c138560ba` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 1% | 56.0490 | 16 | `440945c23ba87593` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 59.4047 | 16 | `369c79d5e33489ff` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 59.5878 | 1 | `58444e2881d8ace1` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 57.1860 | 1 | `f375352b14c53a5d` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | 106.7516 | 2 | `d2d60d2f690d2458` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 100% | 103.0167 | 2 | `08916c0b1967f3af` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 1% | 159.3048 | 4 | `1d1bcace08ebe9c8` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 168.3432 | 4 | `08d6fe70c4691cd5` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | 282.2786 | 8 | `69ea9d1da1b4ee7b` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 100% | 303.9960 | 8 | `bdabb7031fac52ff` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 1% | 547.2481 | 16 | `4599983814752c8d` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 550.2252 | 16 | `884a64cfe91ffab7` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 0.9880 | false |
| notmat | 1 | 1% | 1.0666 | false |
| mat | 1 | 100% | 0.9458 | false |
| notmat | 1 | 100% | 0.9524 | false |
| mat | 2 | 1% | 0.9700 | false |
| notmat | 2 | 1% | 1.8282 | false |
| mat | 2 | 100% | 1.0438 | false |
| notmat | 2 | 100% | 1.8733 | false |
| mat | 4 | 1% | 1.0250 | false |
| notmat | 4 | 1% | 2.8649 | false |
| mat | 4 | 100% | 0.9607 | false |
| notmat | 4 | 100% | 2.8705 | false |
| mat | 8 | 1% | 1.0576 | false |
| notmat | 8 | 1% | 5.0926 | false |
| mat | 8 | 100% | 1.1009 | false |
| notmat | 8 | 100% | 5.5832 | false |
| mat | 16 | 1% | 0.9882 | false |
| notmat | 16 | 1% | 9.6485 | false |
| mat | 16 | 100% | 1.0555 | false |
| notmat | 16 | 100% | 9.7768 | false |

## Repeated-design position control

| first median ms | last median ms | ratio (last/first) | same plan |
|---|---|---|---|
| 56.275085 | 55.662172 | 0.9891 | false |

## Latency samples and derived throughput

_The reader is sequential, so `queries/s` is 1000/median_ms, a derived rate and not a concurrent-capacity test. min/max are the five sample extremes._

| axis | status | min ms | median ms | max ms | samples | derived queries/s |
|---|---|---|---|---|---|---|
| family=bag;stmt=comma_join;threshold=100000 | valid | 89.5743 | 91.9948 | 96.9116 | 5 | 10.9 |
| family=bag;stmt=comma_join;threshold=180000 | valid | 109.3937 | 113.7601 | 126.9213 | 5 | 8.8 |
| family=bag;stmt=comma_join;threshold=20 | valid | 53.2145 | 55.0820 | 56.4422 | 5 | 18.2 |
| family=bag;stmt=comma_join;threshold=200 | valid | 52.8766 | 54.6742 | 58.3135 | 5 | 18.3 |
| family=bag;stmt=comma_join;threshold=2000 | valid | 59.0946 | 59.8192 | 76.3266 | 5 | 16.7 |
| family=bag;stmt=comma_join;threshold=20000 | valid | 68.0175 | 70.3956 | 76.4342 | 5 | 14.2 |
| family=bag;stmt=cte_filtered;threshold=100000 | valid | 87.2616 | 89.2507 | 91.3997 | 5 | 11.2 |
| family=bag;stmt=cte_filtered;threshold=180000 | valid | 107.6889 | 109.8643 | 112.3091 | 5 | 9.1 |
| family=bag;stmt=cte_filtered;threshold=20 | valid | 58.0684 | 60.8675 | 60.9368 | 5 | 16.4 |
| family=bag;stmt=cte_filtered;threshold=200 | valid | 52.8814 | 56.4916 | 58.5563 | 5 | 17.7 |
| family=bag;stmt=cte_filtered;threshold=2000 | valid | 56.7932 | 57.2131 | 58.9069 | 5 | 17.5 |
| family=bag;stmt=cte_filtered;threshold=20000 | valid | 65.8306 | 66.9959 | 70.7919 | 5 | 14.9 |
| family=bag;stmt=cte_filtered_mat;threshold=100000 | valid | 94.0182 | 96.0385 | 99.0065 | 5 | 10.4 |
| family=bag;stmt=cte_filtered_mat;threshold=180000 | valid | 106.7712 | 112.7105 | 121.0301 | 5 | 8.9 |
| family=bag;stmt=cte_filtered_mat;threshold=20 | valid | 52.5542 | 53.8213 | 61.3174 | 5 | 18.6 |
| family=bag;stmt=cte_filtered_mat;threshold=200 | valid | 56.1536 | 57.1962 | 57.5547 | 5 | 17.5 |
| family=bag;stmt=cte_filtered_mat;threshold=2000 | valid | 56.6042 | 57.6344 | 61.3802 | 5 | 17.4 |
| family=bag;stmt=cte_filtered_mat;threshold=20000 | valid | 69.7025 | 73.1248 | 79.7572 | 5 | 13.7 |
| family=bag;stmt=cte_filtered_notmat;threshold=100000 | valid | 90.7091 | 92.0216 | 97.9692 | 5 | 10.9 |
| family=bag;stmt=cte_filtered_notmat;threshold=180000 | valid | 105.9979 | 108.3784 | 111.8268 | 5 | 9.2 |
| family=bag;stmt=cte_filtered_notmat;threshold=20 | valid | 55.8215 | 55.9078 | 59.7105 | 5 | 17.9 |
| family=bag;stmt=cte_filtered_notmat;threshold=200 | valid | 56.0892 | 58.7585 | 62.9243 | 5 | 17.0 |
| family=bag;stmt=cte_filtered_notmat;threshold=2000 | valid | 52.6448 | 55.1731 | 59.3341 | 5 | 18.1 |
| family=bag;stmt=cte_filtered_notmat;threshold=20000 | valid | 65.4987 | 70.2809 | 71.9927 | 5 | 14.2 |
| family=bag;stmt=cte_joined_outer;threshold=100000 | valid | 87.1075 | 89.9376 | 92.7876 | 5 | 11.1 |
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 107.3994 | 114.9823 | 115.8450 | 5 | 8.7 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 55.2432 | 56.2751 | 60.4404 | 5 | 17.8 |
| family=bag;stmt=cte_joined_outer;threshold=200 | valid | 55.8071 | 56.5626 | 60.3743 | 5 | 17.7 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 55.8816 | 57.4451 | 57.9640 | 5 | 17.4 |
| family=bag;stmt=cte_joined_outer;threshold=20000 | valid | 68.3207 | 70.0080 | 75.9994 | 5 | 14.3 |
| family=bag;stmt=cte_joined_outer;threshold=20;position=last | valid | 53.1773 | 55.6622 | 58.2298 | 5 | 18.0 |
| family=bag;stmt=cte_joined_outer_mat;threshold=100000 | valid | 97.3164 | 98.7561 | 102.0419 | 5 | 10.1 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 106.2465 | 107.1639 | 111.3738 | 5 | 9.3 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 84.8267 | 88.7820 | 90.3483 | 5 | 11.3 |
| family=bag;stmt=cte_joined_outer_mat;threshold=200 | valid | 84.4105 | 85.1143 | 89.6326 | 5 | 11.7 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 85.2195 | 88.5428 | 95.8191 | 5 | 11.3 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20000 | valid | 89.5812 | 93.1909 | 99.4853 | 5 | 10.7 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=100000 | valid | 85.8075 | 93.4837 | 100.4824 | 5 | 10.7 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 111.1394 | 114.0428 | 116.4458 | 5 | 8.8 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 55.3357 | 58.0128 | 61.2653 | 5 | 17.2 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=200 | valid | 56.1986 | 57.7281 | 61.0285 | 5 | 17.3 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 56.5631 | 56.9881 | 59.9923 | 5 | 17.5 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20000 | valid | 65.4845 | 69.2583 | 70.2180 | 5 | 14.4 |
| family=bag;stmt=filtered_derived;threshold=100000 | valid | 89.6302 | 91.7815 | 105.1221 | 5 | 10.9 |
| family=bag;stmt=filtered_derived;threshold=180000 | valid | 105.4511 | 110.3726 | 117.9270 | 5 | 9.1 |
| family=bag;stmt=filtered_derived;threshold=20 | valid | 55.4712 | 56.7067 | 57.9140 | 5 | 17.6 |
| family=bag;stmt=filtered_derived;threshold=200 | valid | 53.4956 | 54.8072 | 56.6012 | 5 | 18.2 |
| family=bag;stmt=filtered_derived;threshold=2000 | valid | 57.0031 | 57.8433 | 72.1579 | 5 | 17.3 |
| family=bag;stmt=filtered_derived;threshold=20000 | valid | 65.2134 | 66.6618 | 67.2865 | 5 | 15.0 |
| family=bag;stmt=join_on;threshold=100000 | valid | 92.2633 | 93.1732 | 97.2295 | 5 | 10.7 |
| family=bag;stmt=join_on;threshold=180000 | valid | 106.4955 | 112.2459 | 113.8421 | 5 | 8.9 |
| family=bag;stmt=join_on;threshold=20 | valid | 54.6056 | 57.8984 | 58.7487 | 5 | 17.3 |
| family=bag;stmt=join_on;threshold=200 | valid | 54.7100 | 55.7144 | 56.9069 | 5 | 17.9 |
| family=bag;stmt=join_on;threshold=2000 | valid | 57.7082 | 58.9139 | 60.7816 | 5 | 17.0 |
| family=bag;stmt=join_on;threshold=20000 | valid | 68.4410 | 69.1130 | 72.8066 | 5 | 14.5 |
| family=bag;stmt=join_where;threshold=100000 | valid | 84.0468 | 88.6960 | 92.9751 | 5 | 11.3 |
| family=bag;stmt=join_where;threshold=180000 | valid | 112.8066 | 114.4575 | 127.4657 | 5 | 8.7 |
| family=bag;stmt=join_where;threshold=20 | valid | 55.2062 | 56.3768 | 59.3275 | 5 | 17.7 |
| family=bag;stmt=join_where;threshold=200 | valid | 56.1983 | 60.7198 | 83.2787 | 5 | 16.5 |
| family=bag;stmt=join_where;threshold=2000 | valid | 54.9878 | 57.4667 | 58.4795 | 5 | 17.4 |
| family=bag;stmt=join_where;threshold=20000 | valid | 66.3947 | 67.5504 | 68.6264 | 5 | 14.8 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 55.9626 | 56.7183 | 57.5826 | 5 | 17.6 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 56.1775 | 56.2788 | 57.2890 | 5 | 17.8 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 54.1609 | 55.8683 | 59.8053 | 5 | 17.9 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 58.0839 | 60.0464 | 60.8551 | 5 | 16.7 |
| family=reuse;policy=default;refs=2;consumer=1% | valid | 56.8999 | 58.3900 | 59.3456 | 5 | 17.1 |
| family=reuse;policy=default;refs=2;consumer=100% | valid | 52.3631 | 54.9913 | 56.1706 | 5 | 18.2 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 54.6743 | 55.6060 | 56.3201 | 5 | 18.0 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 54.7316 | 58.6463 | 64.2791 | 5 | 17.1 |
| family=reuse;policy=default;refs=8;consumer=1% | valid | 54.0885 | 55.4290 | 57.5875 | 5 | 18.0 |
| family=reuse;policy=default;refs=8;consumer=100% | valid | 54.0173 | 54.4486 | 56.7163 | 5 | 18.4 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 54.7454 | 56.0490 | 57.8045 | 5 | 17.8 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 57.7667 | 59.4047 | 61.8184 | 5 | 16.8 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 51.9352 | 55.2001 | 56.4109 | 5 | 18.1 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 54.8576 | 56.7918 | 59.1254 | 5 | 17.6 |
| family=reuse;policy=mat;refs=2;consumer=1% | valid | 54.5965 | 56.6383 | 59.6803 | 5 | 17.7 |
| family=reuse;policy=mat;refs=2;consumer=100% | valid | 56.7742 | 57.4015 | 60.0183 | 5 | 17.4 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 56.0719 | 56.9978 | 59.8260 | 5 | 17.5 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 53.8817 | 56.3408 | 58.2055 | 5 | 17.7 |
| family=reuse;policy=mat;refs=8;consumer=1% | valid | 56.8025 | 58.6201 | 60.0923 | 5 | 17.1 |
| family=reuse;policy=mat;refs=8;consumer=100% | valid | 59.1012 | 59.9441 | 62.3977 | 5 | 16.7 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 520.4830 | 547.2481 | 574.9995 | 5 | 1.8 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 537.7348 | 550.2252 | 589.6911 | 5 | 1.8 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 57.0271 | 59.5878 | 64.9599 | 5 | 16.8 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 55.6332 | 57.1860 | 58.8092 | 5 | 17.5 |
| family=reuse;policy=notmat;refs=2;consumer=1% | valid | 103.8959 | 106.7516 | 118.7071 | 5 | 9.4 |
| family=reuse;policy=notmat;refs=2;consumer=100% | valid | 101.0680 | 103.0167 | 105.4813 | 5 | 9.7 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 151.7531 | 159.3048 | 209.5189 | 5 | 6.3 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 160.2440 | 168.3432 | 170.3516 | 5 | 5.9 |
| family=reuse;policy=notmat;refs=8;consumer=1% | valid | 270.1323 | 282.2786 | 324.2557 | 5 | 3.5 |
| family=reuse;policy=notmat;refs=8;consumer=100% | valid | 303.5826 | 303.9960 | 343.2992 | 5 | 3.3 |

## Matched blended workload

- phase: `isolated`; target offered rate: 10 updates/s
- offered 0, committed 0, failed 0 over 0.000 s
- offered rate 0.000/s, completed rate 0.000/s
- per key: no keys touched
- ledger: 0 keys reconciled, 0 mismatch

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42982648`, 5 interleaved rounds
  1. family=bag;stmt=cte_joined_outer;threshold=20
  2. family=reuse;policy=notmat;refs=1;consumer=1%
  3. family=reuse;policy=default;refs=4;consumer=100%
  4. family=bag;stmt=cte_filtered;threshold=20
  5. family=bag;stmt=join_where;threshold=20
  6. family=reuse;policy=notmat;refs=2;consumer=1%
  7. family=reuse;policy=notmat;refs=4;consumer=1%
  8. family=bag;stmt=comma_join;threshold=2000
  9. family=bag;stmt=join_where;threshold=180000
  10. family=bag;stmt=cte_filtered_mat;threshold=100000
  11. family=reuse;policy=mat;refs=1;consumer=100%
  12. family=bag;stmt=cte_joined_outer_mat;threshold=20000
  13. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  14. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  15. family=bag;stmt=cte_filtered_notmat;threshold=100000
  16. family=bag;stmt=cte_filtered;threshold=2000
  17. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  18. family=bag;stmt=join_on;threshold=20000
  19. family=reuse;policy=notmat;refs=8;consumer=100%
  20. family=reuse;policy=mat;refs=8;consumer=100%
  21. family=bag;stmt=cte_joined_outer_mat;threshold=20
  22. family=reuse;policy=mat;refs=8;consumer=1%
  23. family=bag;stmt=cte_filtered_mat;threshold=200
  24. family=bag;stmt=cte_joined_outer_notmat;threshold=200
  25. family=bag;stmt=join_on;threshold=100000
  26. family=reuse;policy=mat;refs=16;consumer=100%
  27. family=bag;stmt=cte_joined_outer_mat;threshold=100000
  28. family=bag;stmt=cte_filtered_notmat;threshold=20
  29. family=bag;stmt=cte_filtered_notmat;threshold=20000
  30. family=bag;stmt=cte_joined_outer_mat;threshold=200
  31. family=reuse;policy=default;refs=16;consumer=1%
  32. family=bag;stmt=cte_joined_outer;threshold=200
  33. family=bag;stmt=comma_join;threshold=100000
  34. family=reuse;policy=default;refs=1;consumer=100%
  35. family=reuse;policy=notmat;refs=2;consumer=100%
  36. family=bag;stmt=comma_join;threshold=20000
  37. family=reuse;policy=mat;refs=2;consumer=100%
  38. family=bag;stmt=cte_joined_outer;threshold=20000
  39. family=bag;stmt=comma_join;threshold=180000
  40. family=bag;stmt=cte_joined_outer_notmat;threshold=20000
  41. family=bag;stmt=filtered_derived;threshold=100000
  42. family=bag;stmt=cte_filtered;threshold=20000
  43. family=bag;stmt=cte_joined_outer;threshold=180000
  44. family=reuse;policy=notmat;refs=8;consumer=1%
  45. family=bag;stmt=cte_filtered_mat;threshold=20000
  46. family=bag;stmt=cte_filtered_notmat;threshold=200
  47. family=reuse;policy=mat;refs=4;consumer=1%
  48. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  49. family=bag;stmt=comma_join;threshold=20
  50. family=reuse;policy=default;refs=4;consumer=1%
  51. family=bag;stmt=cte_filtered_mat;threshold=20
  52. family=bag;stmt=cte_joined_outer;threshold=100000
  53. family=bag;stmt=cte_filtered;threshold=200
  54. family=bag;stmt=cte_filtered;threshold=100000
  55. family=bag;stmt=cte_filtered;threshold=180000
  56. family=bag;stmt=comma_join;threshold=200
  57. family=bag;stmt=cte_joined_outer_notmat;threshold=100000
  58. family=bag;stmt=join_on;threshold=180000
  59. family=bag;stmt=filtered_derived;threshold=20000
  60. family=reuse;policy=mat;refs=16;consumer=1%
  61. family=bag;stmt=join_where;threshold=100000
  62. family=bag;stmt=join_on;threshold=200
  63. family=reuse;policy=default;refs=8;consumer=100%
  64. family=reuse;policy=default;refs=16;consumer=100%
  65. family=bag;stmt=cte_filtered_notmat;threshold=180000
  66. family=reuse;policy=mat;refs=4;consumer=100%
  67. family=bag;stmt=cte_filtered_mat;threshold=180000
  68. family=bag;stmt=filtered_derived;threshold=180000
  69. family=bag;stmt=join_where;threshold=20000
  70. family=reuse;policy=default;refs=8;consumer=1%
  71. family=bag;stmt=cte_filtered_notmat;threshold=2000
  72. family=bag;stmt=filtered_derived;threshold=200
  73. family=reuse;policy=mat;refs=1;consumer=1%
  74. family=reuse;policy=notmat;refs=1;consumer=100%
  75. family=reuse;policy=default;refs=2;consumer=100%
  76. family=reuse;policy=notmat;refs=16;consumer=100%
  77. family=bag;stmt=join_where;threshold=200
  78. family=reuse;policy=notmat;refs=4;consumer=100%
  79. family=reuse;policy=default;refs=2;consumer=1%
  80. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  81. family=bag;stmt=join_on;threshold=2000
  82. family=bag;stmt=cte_joined_outer;threshold=2000
  83. family=bag;stmt=join_where;threshold=2000
  84. family=bag;stmt=filtered_derived;threshold=20
  85. family=reuse;policy=default;refs=1;consumer=1%
  86. family=reuse;policy=mat;refs=2;consumer=1%
  87. family=reuse;policy=notmat;refs=16;consumer=1%
  88. family=bag;stmt=filtered_derived;threshold=2000
  89. family=bag;stmt=cte_filtered_mat;threshold=2000
  90. family=bag;stmt=join_on;threshold=20
  91. family=bag;stmt=cte_joined_outer;threshold=20;position=last

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
| family=bag;stmt=comma_join;threshold=100000 | j06_comma_join | valid | `297e281c684bb215` |  |
| family=bag;stmt=comma_join;threshold=180000 | j06_comma_join | valid | `b7b6d59405460fe0` |  |
| family=bag;stmt=comma_join;threshold=20 | j06_comma_join | valid | `718b53f317684a59` |  |
| family=bag;stmt=comma_join;threshold=200 | j06_comma_join | valid | `ed6b0c1f59f1d950` |  |
| family=bag;stmt=comma_join;threshold=2000 | j06_comma_join | valid | `6d9efa91e9850b6e` |  |
| family=bag;stmt=comma_join;threshold=20000 | j06_comma_join | valid | `ec26e2972ef0b295` |  |
| family=bag;stmt=cte_filtered;threshold=100000 | j04_cte_filtered | valid | `8a2681983b70f275` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=180000 | j04_cte_filtered | valid | `fefc22d0b3ab39a9` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=20 | j04_cte_filtered | valid | `fb04b8cfeec0c948` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=200 | j04_cte_filtered | valid | `ddcdb097053903e6` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=2000 | j04_cte_filtered | valid | `b3947c0b2832a526` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=20000 | j04_cte_filtered | valid | `e8a59337e46d4dc9` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_mat;threshold=100000 | j04m_cte_filtered_mat | valid | `de4cdb0468a76bc0` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=180000 | j04m_cte_filtered_mat | valid | `ad61b87c064e2951` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=20 | j04m_cte_filtered_mat | valid | `fb405b602f2b3f78` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=200 | j04m_cte_filtered_mat | valid | `95c900d2338fe891` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=2000 | j04m_cte_filtered_mat | valid | `0dd64dcf6b4d8c01` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=20000 | j04m_cte_filtered_mat | valid | `61895553885c2fe3` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_notmat;threshold=100000 | j04n_cte_filtered_notmat | valid | `3c9e8f43de54ece4` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=180000 | j04n_cte_filtered_notmat | valid | `893a7d55ce3d8e8e` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=20 | j04n_cte_filtered_notmat | valid | `5404e29385bcb96a` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=200 | j04n_cte_filtered_notmat | valid | `7860e22edd887a48` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=2000 | j04n_cte_filtered_notmat | valid | `f096d35a9dc2cb05` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=20000 | j04n_cte_filtered_notmat | valid | `7bad6f76ed4c5e13` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=100000 | j05_cte_joined_outer | valid | `fcfffe5c695192ec` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=180000 | j05_cte_joined_outer | valid | `d0b4eb61b8479c1c` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20 | j05_cte_joined_outer | valid | `1f2f6a607133a92f` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=200 | j05_cte_joined_outer | valid | `71dee3b73d1d9816` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=2000 | j05_cte_joined_outer | valid | `d8b1b767baf7e55a` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20000 | j05_cte_joined_outer | valid | `5ed7fd236835c425` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20;position=last | j05_cte_joined_outer | valid | `a531ad8996f5694f` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_mat;threshold=100000 | j05m_cte_joined_outer_mat | valid | `0d68de0c5caae90f` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | j05m_cte_joined_outer_mat | valid | `869e306b5fce9634` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | j05m_cte_joined_outer_mat | valid | `4e1a26cab3bac0b6` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=200 | j05m_cte_joined_outer_mat | valid | `6d93c8d86a553364` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | j05m_cte_joined_outer_mat | valid | `29ac36a5e22f41aa` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20000 | j05m_cte_joined_outer_mat | valid | `d224088dbcfbe1f4` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_notmat;threshold=100000 | j05n_cte_joined_outer_notmat | valid | `666c6f2cc0aea91c` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | j05n_cte_joined_outer_notmat | valid | `1ce140d2d7adfdde` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | j05n_cte_joined_outer_notmat | valid | `0eb19391f26601df` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=200 | j05n_cte_joined_outer_notmat | valid | `76e6836dec653d7b` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | j05n_cte_joined_outer_notmat | valid | `925be91a207cc8c8` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20000 | j05n_cte_joined_outer_notmat | valid | `37a2147dbc8ad176` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=filtered_derived;threshold=100000 | j03_filtered_derived | valid | `b93ae517f69bc198` |  |
| family=bag;stmt=filtered_derived;threshold=180000 | j03_filtered_derived | valid | `d540c1159ba295d5` |  |
| family=bag;stmt=filtered_derived;threshold=20 | j03_filtered_derived | valid | `c8fdabded402e861` |  |
| family=bag;stmt=filtered_derived;threshold=200 | j03_filtered_derived | valid | `fa269c1175d92bcb` |  |
| family=bag;stmt=filtered_derived;threshold=2000 | j03_filtered_derived | valid | `7678c688507fe884` |  |
| family=bag;stmt=filtered_derived;threshold=20000 | j03_filtered_derived | valid | `032a3a91bd28e166` |  |
| family=bag;stmt=join_on;threshold=100000 | j02_join_on | valid | `e747b8fa9c183710` |  |
| family=bag;stmt=join_on;threshold=180000 | j02_join_on | valid | `798297d4ca3cd69d` |  |
| family=bag;stmt=join_on;threshold=20 | j02_join_on | valid | `c0833b6a22f6ab49` |  |
| family=bag;stmt=join_on;threshold=200 | j02_join_on | valid | `a166abc6a69228b1` |  |
| family=bag;stmt=join_on;threshold=2000 | j02_join_on | valid | `f507a09f168758f1` |  |
| family=bag;stmt=join_on;threshold=20000 | j02_join_on | valid | `eef2eba30bbf7ee1` |  |
| family=bag;stmt=join_where;threshold=100000 | j01_join_where | valid | `4c2277e93832036a` |  |
| family=bag;stmt=join_where;threshold=180000 | j01_join_where | valid | `2faea02a0d8f860e` |  |
| family=bag;stmt=join_where;threshold=20 | j01_join_where | valid | `71040d1f5cd2cec5` |  |
| family=bag;stmt=join_where;threshold=200 | j01_join_where | valid | `23d669c031e6ab31` |  |
| family=bag;stmt=join_where;threshold=2000 | j01_join_where | valid | `7879e60eab1cf7e7` |  |
| family=bag;stmt=join_where;threshold=20000 | j01_join_where | valid | `1c1cd0cb9b0ab32b` |  |
| family=reuse;policy=default;refs=16;consumer=1% | r_default_n16_c1 | valid | `72e2200f58a6b85e` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=16;consumer=100% | r_default_n16_c100 | valid | `3d6e2be0b00ff813` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=1;consumer=1% | r_default_n1_c1 | valid | `88832648c1133a9b` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=1;consumer=100% | r_default_n1_c100 | valid | `1a133a13b7051a2b` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=2;consumer=1% | r_default_n2_c1 | valid | `2b6e36da3d1af17f` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=2;consumer=100% | r_default_n2_c100 | valid | `179dbfd898a87226` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=1% | r_default_n4_c1 | valid | `e300101f501895e9` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=100% | r_default_n4_c100 | valid | `ba9898a9027d81dc` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=8;consumer=1% | r_default_n8_c1 | valid | `fa27563ccff22d0e` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=8;consumer=100% | r_default_n8_c100 | valid | `dece197c8f297db0` | policy=default; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=1% | r_mat_n16_c1 | valid | `440945c23ba87593` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=100% | r_mat_n16_c100 | valid | `369c79d5e33489ff` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=1% | r_mat_n1_c1 | valid | `dfa7dd7ef8545674` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=100% | r_mat_n1_c100 | valid | `c4c72f423291626f` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=2;consumer=1% | r_mat_n2_c1 | valid | `f5847cdbdbefacbb` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=2;consumer=100% | r_mat_n2_c100 | valid | `9af650e3aaf44181` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=1% | r_mat_n4_c1 | valid | `798f8430f2e62ba9` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=100% | r_mat_n4_c100 | valid | `7774211ba86ffb34` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=8;consumer=1% | r_mat_n8_c1 | valid | `4e687ec3ffcd47c3` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=8;consumer=100% | r_mat_n8_c100 | valid | `e410ab4c138560ba` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=notmat;refs=16;consumer=1% | r_notmat_n16_c1 | valid | `4599983814752c8d` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `884a64cfe91ffab7` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `58444e2881d8ace1` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `f375352b14c53a5d` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=2;consumer=1% | r_notmat_n2_c1 | valid | `d2d60d2f690d2458` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=2;consumer=100% | r_notmat_n2_c100 | valid | `08916c0b1967f3af` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `1d1bcace08ebe9c8` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `08d6fe70c4691cd5` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=8;consumer=1% | r_notmat_n8_c1 | valid | `69ea9d1da1b4ee7b` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=8;consumer=100% | r_notmat_n8_c100 | valid | `bdabb7031fac52ff` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, DIST, DEBUG, FORMAT TEXT` supported on yugabyte
- one host, single client: no independent hardware, real-network or equal-total-budget comparison
- no concurrent writers: the isolated read contract only; a matched blended read/write phase is a recorded gap
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
