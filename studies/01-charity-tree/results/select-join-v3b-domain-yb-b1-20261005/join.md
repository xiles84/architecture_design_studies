# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 91 cells, 91 valid, engine `yugabyte`, topology `single`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-v3b-domain-yb-b1-20261005`
- run tag: `run/study-01/select-join-v3b-domain-yb-b1-20261005`
- environment: `host-zenbook-ux5406sa`
- repo commit: `7243fee074bfb16fc5ad8b6afc1a2d7549080935`
- engine version: `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit`
- database engine image: `docker.io/yugabytedb/yugabyte:2025.2.6.0-b111` @ `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`
- benchmark client image: `localhost/charitybench-select:select-join-v3b-domain-yb-b1-20261005` @ `f3f1e10c44cbcb4c7b5bc56ed95a4f071682b6f7d33b5e0a7b688c4271eeef88`
- inputs digest: `28deac3dcf98eeb5da0432bc8fb437001b425cc7b793e3e2c6f78a35de377235`
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
| yb-single | 933 | 2 |

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
| j01_join_where | valid | true |  |
| j02_join_on | valid | true |  |
| j03_filtered_derived | valid | true |  |
| j04_cte_filtered | valid | true | policy=default; observed_cte_scan=false |
| j04m_cte_filtered_mat | valid | true | policy=materialized; observed_cte_scan=true |
| j04n_cte_filtered_notmat | valid | true | policy=not_materialized; observed_cte_scan=false |
| j05_cte_joined_outer | valid | true | policy=default; observed_cte_scan=false |
| j05m_cte_joined_outer_mat | valid | false | policy=materialized; observed_cte_scan=true |
| j05n_cte_joined_outer_notmat | valid | true | policy=not_materialized; observed_cte_scan=false |
| j06_comma_join | valid | true |  |
| j11_optional_on | valid | true |  |
| r_default_n4_c100 | valid | false | policy=default; observed_cte_scan=true |
| r_mat_n4_c100 | valid | false | policy=materialized; observed_cte_scan=true |
| r_notmat_n4_c100 | valid | false | policy=not_materialized; observed_cte_scan=false |

## P1 — six inner-join spellings and their CTE policies (median ms)

| threshold | join_where | join_on | filtered_derived | cte_filtered | cte_filtered_mat | cte_filtered_notmat | cte_joined_outer | cte_joined_outer_mat | cte_joined_outer_notmat | comma_join |
|---|---|---|---|---|---|---|---|---|---|---|
| 20 (0.01%) | 2.1445 | 1.7592 | 1.6527 | 1.5390 | 1.9284 | 2.5929 | 2.1226 | 82.0061 | 1.7864 | 1.8820 |
| 200 (0.1%) | 2.7167 | 3.5448 | 2.9905 | 3.7898 | 3.0267 | 3.6592 | 3.5976 | 82.9406 | 2.9775 | 3.0425 |
| 2000 (1%) | 9.0819 | 9.1065 | 9.6755 | 9.1381 | 10.5616 | 9.1046 | 9.2196 | 83.4664 | 8.9502 | 8.5453 |
| 20000 (10%) | 63.0249 | 63.8258 | 63.8833 | 63.1623 | 64.9904 | 60.9284 | 62.7514 | 84.6818 | 65.8629 | 62.7880 |
| 100000 (50%) | 86.3689 | 87.6960 | 94.8574 | 89.5374 | 91.9753 | 86.6649 | 88.3008 | 94.9390 | 88.0015 | 86.7790 |
| 180000 (90%) | 111.4336 | 108.8605 | 110.4637 | 106.8858 | 116.0640 | 111.4675 | 108.3577 | 107.2638 | 107.1141 | 106.0419 |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | join_on vs join_where | 2.1445 | 1.7592 | 0.8203 | false |
| 20 | filtered_derived vs join_where | 2.1445 | 1.6527 | 0.7707 | false |
| 20 | cte_filtered vs join_where | 2.1445 | 1.5390 | 0.7177 | false |
| 20 | cte_filtered_mat vs join_where | 2.1445 | 1.9284 | 0.8992 | false |
| 20 | cte_filtered_notmat vs join_where | 2.1445 | 2.5929 | 1.2091 | false |
| 20 | cte_joined_outer vs join_where | 2.1445 | 2.1226 | 0.9898 | false |
| 20 | cte_joined_outer_mat vs join_where | 2.1445 | 82.0061 | 38.2394 | false |
| 20 | cte_joined_outer_notmat vs join_where | 2.1445 | 1.7864 | 0.8330 | false |
| 20 | comma_join vs join_where | 2.1445 | 1.8820 | 0.8776 | false |
| 200 | join_on vs join_where | 2.7167 | 3.5448 | 1.3048 | false |
| 200 | filtered_derived vs join_where | 2.7167 | 2.9905 | 1.1008 | false |
| 200 | cte_filtered vs join_where | 2.7167 | 3.7898 | 1.3950 | false |
| 200 | cte_filtered_mat vs join_where | 2.7167 | 3.0267 | 1.1141 | false |
| 200 | cte_filtered_notmat vs join_where | 2.7167 | 3.6592 | 1.3469 | false |
| 200 | cte_joined_outer vs join_where | 2.7167 | 3.5976 | 1.3243 | false |
| 200 | cte_joined_outer_mat vs join_where | 2.7167 | 82.9406 | 30.5303 | false |
| 200 | cte_joined_outer_notmat vs join_where | 2.7167 | 2.9775 | 1.0960 | false |
| 200 | comma_join vs join_where | 2.7167 | 3.0425 | 1.1199 | false |
| 2000 | join_on vs join_where | 9.0819 | 9.1065 | 1.0027 | false |
| 2000 | filtered_derived vs join_where | 9.0819 | 9.6755 | 1.0654 | false |
| 2000 | cte_filtered vs join_where | 9.0819 | 9.1381 | 1.0062 | false |
| 2000 | cte_filtered_mat vs join_where | 9.0819 | 10.5616 | 1.1629 | false |
| 2000 | cte_filtered_notmat vs join_where | 9.0819 | 9.1046 | 1.0025 | false |
| 2000 | cte_joined_outer vs join_where | 9.0819 | 9.2196 | 1.0152 | false |
| 2000 | cte_joined_outer_mat vs join_where | 9.0819 | 83.4664 | 9.1904 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 9.0819 | 8.9502 | 0.9855 | false |
| 2000 | comma_join vs join_where | 9.0819 | 8.5453 | 0.9409 | false |
| 20000 | join_on vs join_where | 63.0249 | 63.8258 | 1.0127 | false |
| 20000 | filtered_derived vs join_where | 63.0249 | 63.8833 | 1.0136 | false |
| 20000 | cte_filtered vs join_where | 63.0249 | 63.1623 | 1.0022 | false |
| 20000 | cte_filtered_mat vs join_where | 63.0249 | 64.9904 | 1.0312 | false |
| 20000 | cte_filtered_notmat vs join_where | 63.0249 | 60.9284 | 0.9667 | false |
| 20000 | cte_joined_outer vs join_where | 63.0249 | 62.7514 | 0.9957 | false |
| 20000 | cte_joined_outer_mat vs join_where | 63.0249 | 84.6818 | 1.3436 | false |
| 20000 | cte_joined_outer_notmat vs join_where | 63.0249 | 65.8629 | 1.0450 | false |
| 20000 | comma_join vs join_where | 63.0249 | 62.7880 | 0.9962 | false |
| 100000 | join_on vs join_where | 86.3689 | 87.6960 | 1.0154 | false |
| 100000 | filtered_derived vs join_where | 86.3689 | 94.8574 | 1.0983 | false |
| 100000 | cte_filtered vs join_where | 86.3689 | 89.5374 | 1.0367 | false |
| 100000 | cte_filtered_mat vs join_where | 86.3689 | 91.9753 | 1.0649 | false |
| 100000 | cte_filtered_notmat vs join_where | 86.3689 | 86.6649 | 1.0034 | false |
| 100000 | cte_joined_outer vs join_where | 86.3689 | 88.3008 | 1.0224 | false |
| 100000 | cte_joined_outer_mat vs join_where | 86.3689 | 94.9390 | 1.0992 | false |
| 100000 | cte_joined_outer_notmat vs join_where | 86.3689 | 88.0015 | 1.0189 | false |
| 100000 | comma_join vs join_where | 86.3689 | 86.7790 | 1.0047 | false |
| 180000 | join_on vs join_where | 111.4336 | 108.8605 | 0.9769 | false |
| 180000 | filtered_derived vs join_where | 111.4336 | 110.4637 | 0.9913 | false |
| 180000 | cte_filtered vs join_where | 111.4336 | 106.8858 | 0.9592 | false |
| 180000 | cte_filtered_mat vs join_where | 111.4336 | 116.0640 | 1.0416 | false |
| 180000 | cte_filtered_notmat vs join_where | 111.4336 | 111.4675 | 1.0003 | false |
| 180000 | cte_joined_outer vs join_where | 111.4336 | 108.3577 | 0.9724 | false |
| 180000 | cte_joined_outer_mat vs join_where | 111.4336 | 107.2638 | 0.9626 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 111.4336 | 107.1141 | 0.9612 | false |
| 180000 | comma_join vs join_where | 111.4336 | 106.0419 | 0.9516 | false |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_filtered_mat vs cte_filtered | 1.5390 | 1.9284 | 1.2530 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_filtered_notmat vs cte_filtered | 1.5390 | 2.5929 | 1.6848 | false | policy=not_materialized; observed_cte_scan=false |
| 20 | cte_joined_outer_mat vs cte_joined_outer | 2.1226 | 82.0061 | 38.6349 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 2.1226 | 1.7864 | 0.8416 | false | policy=not_materialized; observed_cte_scan=false |
| 200 | cte_filtered_mat vs cte_filtered | 3.7898 | 3.0267 | 0.7986 | false | policy=materialized; observed_cte_scan=true |
| 200 | cte_filtered_notmat vs cte_filtered | 3.7898 | 3.6592 | 0.9655 | false | policy=not_materialized; observed_cte_scan=false |
| 200 | cte_joined_outer_mat vs cte_joined_outer | 3.5976 | 82.9406 | 23.0545 | false | policy=materialized; observed_cte_scan=true |
| 200 | cte_joined_outer_notmat vs cte_joined_outer | 3.5976 | 2.9775 | 0.8276 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_filtered_mat vs cte_filtered | 9.1381 | 10.5616 | 1.1558 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_filtered_notmat vs cte_filtered | 9.1381 | 9.1046 | 0.9963 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 9.2196 | 83.4664 | 9.0532 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 9.2196 | 8.9502 | 0.9708 | false | policy=not_materialized; observed_cte_scan=false |
| 20000 | cte_filtered_mat vs cte_filtered | 63.1623 | 64.9904 | 1.0289 | false | policy=materialized; observed_cte_scan=true |
| 20000 | cte_filtered_notmat vs cte_filtered | 63.1623 | 60.9284 | 0.9646 | false | policy=not_materialized; observed_cte_scan=false |
| 20000 | cte_joined_outer_mat vs cte_joined_outer | 62.7514 | 84.6818 | 1.3495 | false | policy=materialized; observed_cte_scan=true |
| 20000 | cte_joined_outer_notmat vs cte_joined_outer | 62.7514 | 65.8629 | 1.0496 | false | policy=not_materialized; observed_cte_scan=false |
| 100000 | cte_filtered_mat vs cte_filtered | 89.5374 | 91.9753 | 1.0272 | false | policy=materialized; observed_cte_scan=true |
| 100000 | cte_filtered_notmat vs cte_filtered | 89.5374 | 86.6649 | 0.9679 | false | policy=not_materialized; observed_cte_scan=false |
| 100000 | cte_joined_outer_mat vs cte_joined_outer | 88.3008 | 94.9390 | 1.0752 | false | policy=materialized; observed_cte_scan=true |
| 100000 | cte_joined_outer_notmat vs cte_joined_outer | 88.3008 | 88.0015 | 0.9966 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_filtered_mat vs cte_filtered | 106.8858 | 116.0640 | 1.0859 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_filtered_notmat vs cte_filtered | 106.8858 | 111.4675 | 1.0429 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 108.3577 | 107.2638 | 0.9899 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 108.3577 | 107.1141 | 0.9885 | false | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 7.0894 | 1 | `2b13cc877c6bda0a` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 8.7224 | 1 | `e0694de84b598e72` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | 8.4098 | 2 | `a87de8c863a0d67a` | policy=default; observed_cte_scan=true | valid |
| default | 2 | 100% | 8.4711 | 2 | `1c527585183d2d2e` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 1% | 8.3626 | 4 | `e106c7b3d97f101d` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 8.6106 | 4 | `f59c9010b87620c2` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | 8.8291 | 8 | `e086fc43ae7af671` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 100% | 8.9518 | 8 | `85a8d384d1b80a31` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 1% | 9.3432 | 16 | `ce3db60b8252bd89` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 9.3055 | 16 | `a85b28ffd2c7ef3f` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 8.1718 | 1 | `86f73a0bc7e2f2b8` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 8.8633 | 1 | `b4d11a7446109d2b` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | 8.8764 | 2 | `6caee16a07cc236b` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 100% | 8.3026 | 2 | `29d0a4628a80ae67` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 1% | 8.9820 | 4 | `9d39564c77d9fd35` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 8.9492 | 4 | `df2d57d9e7136b35` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | 9.1722 | 8 | `b32f091ba65ecbed` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 100% | 9.8114 | 8 | `8153e78542b5b2ce` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 1% | 8.9113 | 16 | `9fff03bfad2adfbd` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 9.9545 | 16 | `a2bcc4581afd1802` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 7.4837 | 1 | `b9b78fbfeaf8be98` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 8.3393 | 1 | `9b325d87e73f2655` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | 12.7368 | 2 | `7f12da48414cb532` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 100% | 14.9117 | 2 | `fdc569290a0a5b20` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 1% | 21.9993 | 4 | `83183e20bd36a187` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 31.2281 | 4 | `0d9d53dadb4e295b` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | 89.8282 | 8 | `cde1d48c78542fc4` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 100% | 56.1452 | 8 | `f85c30e87bdb864a` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 1% | 89.5516 | 16 | `1ea309265e75ee47` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 109.1164 | 16 | `677ae2bcc8cd9190` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 1.1527 | false |
| notmat | 1 | 1% | 1.0556 | false |
| mat | 1 | 100% | 1.0162 | false |
| notmat | 1 | 100% | 0.9561 | false |
| mat | 2 | 1% | 1.0555 | false |
| notmat | 2 | 1% | 1.5145 | false |
| mat | 2 | 100% | 0.9801 | false |
| notmat | 2 | 100% | 1.7603 | false |
| mat | 4 | 1% | 1.0741 | false |
| notmat | 4 | 1% | 2.6307 | false |
| mat | 4 | 100% | 1.0393 | false |
| notmat | 4 | 100% | 3.6267 | false |
| mat | 8 | 1% | 1.0389 | false |
| notmat | 8 | 1% | 10.1742 | false |
| mat | 8 | 100% | 1.0960 | false |
| notmat | 8 | 100% | 6.2720 | false |
| mat | 16 | 1% | 0.9538 | false |
| notmat | 16 | 1% | 9.5847 | false |
| mat | 16 | 100% | 1.0697 | false |
| notmat | 16 | 100% | 11.7260 | false |

## Repeated-design position control

| first median ms | last median ms | ratio (last/first) | same plan |
|---|---|---|---|
| 1.882033 | 1.758401 | 0.9343 | false |

## Physical arm

- `charity_pkey`: `CREATE UNIQUE INDEX charity_pkey ON selectjoin_domain.charity USING lsm (charity_id HASH)`
- `donation_person_idx`: `CREATE INDEX donation_person_idx ON selectjoin_domain.donation USING lsm (person_id HASH)`
- `donation_pkey`: `CREATE UNIQUE INDEX donation_pkey ON selectjoin_domain.donation USING lsm (donation_id HASH)`
- `donation_rank_idx`: `CREATE INDEX donation_rank_idx ON selectjoin_domain.donation USING lsm (rank ASC)`
- `person_pkey`: `CREATE UNIQUE INDEX person_pkey ON selectjoin_domain.person USING lsm (person_id HASH)`

- index arm `yb-range-rank`; placement `engine default`; node budget `2 CPU 3g`; client budget `1 CPU 2g`

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
- threshold 200: used
- threshold 2000: used

## Throttling (cgroup `cpu.stat` deltas)

- client: nr_periods 604, nr_throttled 247, throttled_usec 8107224 | quota 100000 100000
, memory 2147483648

- yb-single: nr_periods 671, nr_throttled 76, throttled_usec 2994430 | quota 200000 100000, memory 3221225472

## Access path per cell

_Engine-reported from each cell's own plan: whether the rank index or a sequential scan was used, rows scanned, PostgreSQL buffers and filter removals, and reader retries of transient conflicts. A cell that scanned the whole table at a selective threshold is a different physical regime, not a faster query._

| axis | status | index scan | seq scan | storage rows scanned | index rows scanned | buffers hit/read | rows removed | retries |
|---|---|---|---|---|---|---|---|---|
| family=bag;stmt=comma_join;threshold=100000 | valid | false | true | 202050 | 0 | 0/0 | 0 | 0 |
| family=bag;stmt=comma_join;threshold=180000 | valid | false | true | 202050 | 0 | 0/0 | 0 | 0 |
| family=bag;stmt=comma_join;threshold=20 | valid | true | true | 110 | 20 | 0/0 | 0 | 0 |
| family=bag;stmt=comma_join;threshold=200 | valid | true | true | 2450 | 200 | 0/0 | 0 | 0 |
| family=bag;stmt=comma_join;threshold=2000 | valid | true | true | 6050 | 2000 | 0/0 | 0 | 0 |
| family=bag;stmt=comma_join;threshold=20000 | valid | true | true | 42050 | 20000 | 0/0 | 0 | 0 |
| family=bag;stmt=comma_join;threshold=20;position=last | valid | true | true | 110 | 20 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_filtered;threshold=100000 | valid | false | true | 202050 | 0 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_filtered;threshold=180000 | valid | false | true | 202050 | 0 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_filtered;threshold=20 | valid | true | true | 110 | 20 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_filtered;threshold=200 | valid | true | true | 2450 | 200 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_filtered;threshold=2000 | valid | true | true | 6050 | 2000 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_filtered;threshold=20000 | valid | true | true | 42050 | 20000 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_filtered_mat;threshold=100000 | valid | false | true | 202050 | 0 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_filtered_mat;threshold=180000 | valid | false | true | 202050 | 0 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_filtered_mat;threshold=20 | valid | true | true | 110 | 20 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_filtered_mat;threshold=200 | valid | true | true | 2450 | 200 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_filtered_mat;threshold=2000 | valid | true | true | 6050 | 2000 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_filtered_mat;threshold=20000 | valid | true | true | 42050 | 20000 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_filtered_notmat;threshold=100000 | valid | false | true | 202050 | 0 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_filtered_notmat;threshold=180000 | valid | false | true | 202050 | 0 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_filtered_notmat;threshold=20 | valid | true | true | 110 | 20 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_filtered_notmat;threshold=200 | valid | true | true | 2450 | 200 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_filtered_notmat;threshold=2000 | valid | true | true | 6050 | 2000 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_filtered_notmat;threshold=20000 | valid | true | true | 42050 | 20000 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_joined_outer;threshold=100000 | valid | false | true | 202050 | 0 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | false | true | 202050 | 0 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | true | true | 110 | 20 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_joined_outer;threshold=200 | valid | true | true | 2450 | 200 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | true | true | 6050 | 2000 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_joined_outer;threshold=20000 | valid | true | true | 42050 | 20000 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_joined_outer_mat;threshold=100000 | valid | false | true | 202050 | 0 | 0/0 | 100000 | 0 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | false | true | 202050 | 0 | 0/0 | 20000 | 0 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | false | true | 202050 | 0 | 0/0 | 199980 | 0 |
| family=bag;stmt=cte_joined_outer_mat;threshold=200 | valid | false | true | 202050 | 0 | 0/0 | 199800 | 0 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | false | true | 202050 | 0 | 0/0 | 198000 | 0 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20000 | valid | false | true | 202050 | 0 | 0/0 | 180000 | 0 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=100000 | valid | false | true | 202050 | 0 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | false | true | 202050 | 0 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | true | true | 110 | 20 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=200 | valid | true | true | 2450 | 200 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | true | true | 6050 | 2000 | 0/0 | 0 | 0 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20000 | valid | true | true | 42050 | 20000 | 0/0 | 0 | 0 |
| family=bag;stmt=filtered_derived;threshold=100000 | valid | false | true | 202050 | 0 | 0/0 | 0 | 0 |
| family=bag;stmt=filtered_derived;threshold=180000 | valid | false | true | 202050 | 0 | 0/0 | 0 | 0 |
| family=bag;stmt=filtered_derived;threshold=20 | valid | true | true | 110 | 20 | 0/0 | 0 | 0 |
| family=bag;stmt=filtered_derived;threshold=200 | valid | true | true | 2450 | 200 | 0/0 | 0 | 0 |
| family=bag;stmt=filtered_derived;threshold=2000 | valid | true | true | 6050 | 2000 | 0/0 | 0 | 0 |
| family=bag;stmt=filtered_derived;threshold=20000 | valid | true | true | 42050 | 20000 | 0/0 | 0 | 0 |
| family=bag;stmt=join_on;threshold=100000 | valid | false | true | 202050 | 0 | 0/0 | 0 | 0 |
| family=bag;stmt=join_on;threshold=180000 | valid | false | true | 202050 | 0 | 0/0 | 0 | 0 |
| family=bag;stmt=join_on;threshold=20 | valid | true | true | 110 | 20 | 0/0 | 0 | 0 |
| family=bag;stmt=join_on;threshold=200 | valid | true | true | 2450 | 200 | 0/0 | 0 | 0 |
| family=bag;stmt=join_on;threshold=2000 | valid | true | true | 6050 | 2000 | 0/0 | 0 | 0 |
| family=bag;stmt=join_on;threshold=20000 | valid | true | true | 42050 | 20000 | 0/0 | 0 | 0 |
| family=bag;stmt=join_where;threshold=100000 | valid | false | true | 202050 | 0 | 0/0 | 0 | 0 |
| family=bag;stmt=join_where;threshold=180000 | valid | false | true | 202050 | 0 | 0/0 | 0 | 0 |
| family=bag;stmt=join_where;threshold=20 | valid | true | true | 110 | 20 | 0/0 | 0 | 0 |
| family=bag;stmt=join_where;threshold=200 | valid | true | true | 2450 | 200 | 0/0 | 0 | 0 |
| family=bag;stmt=join_where;threshold=2000 | valid | true | true | 6050 | 2000 | 0/0 | 0 | 0 |
| family=bag;stmt=join_where;threshold=20000 | valid | true | true | 42050 | 20000 | 0/0 | 0 | 0 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | true | false | 4000 | 2000 | 0/0 | 1983 | 0 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | true | false | 4000 | 2000 | 0/0 | 0 | 0 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | true | false | 4000 | 2000 | 0/0 | 0 | 0 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | true | false | 4000 | 2000 | 0/0 | 0 | 0 |
| family=reuse;policy=default;refs=2;consumer=1% | valid | true | false | 4000 | 2000 | 0/0 | 1983 | 0 |
| family=reuse;policy=default;refs=2;consumer=100% | valid | true | false | 4000 | 2000 | 0/0 | 0 | 0 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | true | false | 4000 | 2000 | 0/0 | 1983 | 0 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | true | false | 4000 | 2000 | 0/0 | 0 | 0 |
| family=reuse;policy=default;refs=8;consumer=1% | valid | true | false | 4000 | 2000 | 0/0 | 1983 | 0 |
| family=reuse;policy=default;refs=8;consumer=100% | valid | true | false | 4000 | 2000 | 0/0 | 0 | 0 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | true | false | 4000 | 2000 | 0/0 | 1983 | 0 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | true | false | 4000 | 2000 | 0/0 | 0 | 0 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | true | false | 4000 | 2000 | 0/0 | 1983 | 0 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | true | false | 4000 | 2000 | 0/0 | 0 | 0 |
| family=reuse;policy=mat;refs=2;consumer=1% | valid | true | false | 4000 | 2000 | 0/0 | 1983 | 0 |
| family=reuse;policy=mat;refs=2;consumer=100% | valid | true | false | 4000 | 2000 | 0/0 | 0 | 0 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | true | false | 4000 | 2000 | 0/0 | 1983 | 0 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | true | false | 4000 | 2000 | 0/0 | 0 | 0 |
| family=reuse;policy=mat;refs=8;consumer=1% | valid | true | false | 4000 | 2000 | 0/0 | 1983 | 0 |
| family=reuse;policy=mat;refs=8;consumer=100% | valid | true | false | 4000 | 2000 | 0/0 | 0 | 0 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | true | false | 64000 | 2000 | 0/0 | 0 | 0 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | true | false | 64000 | 2000 | 0/0 | 0 | 0 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | true | false | 4000 | 2000 | 0/0 | 0 | 0 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | true | false | 4000 | 2000 | 0/0 | 0 | 0 |
| family=reuse;policy=notmat;refs=2;consumer=1% | valid | true | false | 8000 | 2000 | 0/0 | 0 | 0 |
| family=reuse;policy=notmat;refs=2;consumer=100% | valid | true | false | 8000 | 2000 | 0/0 | 0 | 0 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | true | false | 16000 | 2000 | 0/0 | 0 | 0 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | true | false | 16000 | 2000 | 0/0 | 0 | 0 |
| family=reuse;policy=notmat;refs=8;consumer=1% | valid | true | false | 32000 | 2000 | 0/0 | 0 | 0 |
| family=reuse;policy=notmat;refs=8;consumer=100% | valid | true | false | 32000 | 2000 | 0/0 | 0 | 0 |

## Latency samples and derived throughput

_The reader is sequential, so `queries/s` is 1000/median_ms, a derived rate and not a concurrent-capacity test. min/max are the five sample extremes. Each sample includes any retry and backoff time the statement needed after a transient conflict; per-cell retry counts appear in the access-path table._

| axis | status | min ms | median ms | max ms | samples | derived queries/s |
|---|---|---|---|---|---|---|
| family=bag;stmt=comma_join;threshold=100000 | valid | 84.1858 | 86.7790 | 90.1206 | 5 | 11.5 |
| family=bag;stmt=comma_join;threshold=180000 | valid | 105.2917 | 106.0419 | 109.6491 | 5 | 9.4 |
| family=bag;stmt=comma_join;threshold=20 | valid | 1.7158 | 1.8820 | 2.1738 | 5 | 531.3 |
| family=bag;stmt=comma_join;threshold=200 | valid | 2.8425 | 3.0425 | 3.4771 | 5 | 328.7 |
| family=bag;stmt=comma_join;threshold=2000 | valid | 7.9574 | 8.5453 | 9.3442 | 5 | 117.0 |
| family=bag;stmt=comma_join;threshold=20000 | valid | 59.8787 | 62.7880 | 64.9048 | 5 | 15.9 |
| family=bag;stmt=comma_join;threshold=20;position=last | valid | 1.1408 | 1.7584 | 2.3250 | 5 | 568.7 |
| family=bag;stmt=cte_filtered;threshold=100000 | valid | 87.9186 | 89.5374 | 90.8439 | 5 | 11.2 |
| family=bag;stmt=cte_filtered;threshold=180000 | valid | 104.6377 | 106.8858 | 110.7107 | 5 | 9.4 |
| family=bag;stmt=cte_filtered;threshold=20 | valid | 1.4383 | 1.5390 | 1.6235 | 5 | 649.8 |
| family=bag;stmt=cte_filtered;threshold=200 | valid | 3.6402 | 3.7898 | 3.8793 | 5 | 263.9 |
| family=bag;stmt=cte_filtered;threshold=2000 | valid | 8.1209 | 9.1381 | 9.2205 | 5 | 109.4 |
| family=bag;stmt=cte_filtered;threshold=20000 | valid | 60.2840 | 63.1623 | 68.2080 | 5 | 15.8 |
| family=bag;stmt=cte_filtered_mat;threshold=100000 | valid | 89.9841 | 91.9753 | 97.2099 | 5 | 10.9 |
| family=bag;stmt=cte_filtered_mat;threshold=180000 | valid | 112.3535 | 116.0640 | 141.1652 | 5 | 8.6 |
| family=bag;stmt=cte_filtered_mat;threshold=20 | valid | 1.7684 | 1.9284 | 2.2199 | 5 | 518.6 |
| family=bag;stmt=cte_filtered_mat;threshold=200 | valid | 2.7765 | 3.0267 | 3.7234 | 5 | 330.4 |
| family=bag;stmt=cte_filtered_mat;threshold=2000 | valid | 9.7652 | 10.5616 | 10.7069 | 5 | 94.7 |
| family=bag;stmt=cte_filtered_mat;threshold=20000 | valid | 63.6720 | 64.9904 | 89.4164 | 5 | 15.4 |
| family=bag;stmt=cte_filtered_notmat;threshold=100000 | valid | 85.7799 | 86.6649 | 88.3146 | 5 | 11.5 |
| family=bag;stmt=cte_filtered_notmat;threshold=180000 | valid | 105.6213 | 111.4675 | 113.1549 | 5 | 9.0 |
| family=bag;stmt=cte_filtered_notmat;threshold=20 | valid | 2.0203 | 2.5929 | 3.4268 | 5 | 385.7 |
| family=bag;stmt=cte_filtered_notmat;threshold=200 | valid | 3.4391 | 3.6592 | 3.9390 | 5 | 273.3 |
| family=bag;stmt=cte_filtered_notmat;threshold=2000 | valid | 8.8017 | 9.1046 | 9.6889 | 5 | 109.8 |
| family=bag;stmt=cte_filtered_notmat;threshold=20000 | valid | 60.2743 | 60.9284 | 68.3546 | 5 | 16.4 |
| family=bag;stmt=cte_joined_outer;threshold=100000 | valid | 83.0484 | 88.3008 | 93.1463 | 5 | 11.3 |
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 103.4996 | 108.3577 | 127.3550 | 5 | 9.2 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 1.7953 | 2.1226 | 2.3185 | 5 | 471.1 |
| family=bag;stmt=cte_joined_outer;threshold=200 | valid | 3.2871 | 3.5976 | 4.1386 | 5 | 278.0 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 8.5679 | 9.2196 | 9.6911 | 5 | 108.5 |
| family=bag;stmt=cte_joined_outer;threshold=20000 | valid | 61.6547 | 62.7514 | 65.0261 | 5 | 15.9 |
| family=bag;stmt=cte_joined_outer_mat;threshold=100000 | valid | 90.5287 | 94.9390 | 102.5862 | 5 | 10.5 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 98.4770 | 107.2638 | 123.5414 | 5 | 9.3 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 76.9829 | 82.0061 | 82.1176 | 5 | 12.2 |
| family=bag;stmt=cte_joined_outer_mat;threshold=200 | valid | 80.4918 | 82.9406 | 85.5006 | 5 | 12.1 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 79.8437 | 83.4664 | 89.3615 | 5 | 12.0 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20000 | valid | 81.7015 | 84.6818 | 89.4280 | 5 | 11.8 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=100000 | valid | 82.2846 | 88.0015 | 89.6453 | 5 | 11.4 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 104.3111 | 107.1141 | 117.6392 | 5 | 9.3 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 1.6662 | 1.7864 | 2.3296 | 5 | 559.8 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=200 | valid | 2.6502 | 2.9775 | 3.0768 | 5 | 335.8 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 8.2961 | 8.9502 | 9.6437 | 5 | 111.7 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20000 | valid | 61.3536 | 65.8629 | 71.9439 | 5 | 15.2 |
| family=bag;stmt=filtered_derived;threshold=100000 | valid | 86.0358 | 94.8574 | 117.8501 | 5 | 10.5 |
| family=bag;stmt=filtered_derived;threshold=180000 | valid | 107.1168 | 110.4637 | 112.1668 | 5 | 9.1 |
| family=bag;stmt=filtered_derived;threshold=20 | valid | 1.2937 | 1.6527 | 2.0782 | 5 | 605.1 |
| family=bag;stmt=filtered_derived;threshold=200 | valid | 2.6530 | 2.9905 | 3.3191 | 5 | 334.4 |
| family=bag;stmt=filtered_derived;threshold=2000 | valid | 9.3876 | 9.6755 | 10.4966 | 5 | 103.4 |
| family=bag;stmt=filtered_derived;threshold=20000 | valid | 61.1357 | 63.8833 | 67.4193 | 5 | 15.7 |
| family=bag;stmt=join_on;threshold=100000 | valid | 86.8859 | 87.6960 | 88.8250 | 5 | 11.4 |
| family=bag;stmt=join_on;threshold=180000 | valid | 105.7406 | 108.8605 | 115.2779 | 5 | 9.2 |
| family=bag;stmt=join_on;threshold=20 | valid | 1.6867 | 1.7592 | 1.8153 | 5 | 568.4 |
| family=bag;stmt=join_on;threshold=200 | valid | 3.4927 | 3.5448 | 4.3923 | 5 | 282.1 |
| family=bag;stmt=join_on;threshold=2000 | valid | 8.0061 | 9.1065 | 9.4320 | 5 | 109.8 |
| family=bag;stmt=join_on;threshold=20000 | valid | 61.9221 | 63.8258 | 74.6866 | 5 | 15.7 |
| family=bag;stmt=join_where;threshold=100000 | valid | 84.5346 | 86.3689 | 91.4910 | 5 | 11.6 |
| family=bag;stmt=join_where;threshold=180000 | valid | 103.2358 | 111.4336 | 124.6433 | 5 | 9.0 |
| family=bag;stmt=join_where;threshold=20 | valid | 1.9957 | 2.1445 | 2.3660 | 5 | 466.3 |
| family=bag;stmt=join_where;threshold=200 | valid | 2.5771 | 2.7167 | 3.3416 | 5 | 368.1 |
| family=bag;stmt=join_where;threshold=2000 | valid | 8.7644 | 9.0819 | 9.7664 | 5 | 110.1 |
| family=bag;stmt=join_where;threshold=20000 | valid | 60.6013 | 63.0249 | 65.2133 | 5 | 15.9 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 8.8713 | 9.3432 | 9.7269 | 5 | 107.0 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 9.0935 | 9.3055 | 10.3753 | 5 | 107.5 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 6.8698 | 7.0894 | 7.7193 | 5 | 141.1 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 8.2939 | 8.7224 | 9.6080 | 5 | 114.6 |
| family=reuse;policy=default;refs=2;consumer=1% | valid | 8.1941 | 8.4098 | 9.8882 | 5 | 118.9 |
| family=reuse;policy=default;refs=2;consumer=100% | valid | 7.9191 | 8.4711 | 9.4702 | 5 | 118.0 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 8.1359 | 8.3626 | 8.7411 | 5 | 119.6 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 8.3465 | 8.6106 | 10.7249 | 5 | 116.1 |
| family=reuse;policy=default;refs=8;consumer=1% | valid | 8.7576 | 8.8291 | 9.3775 | 5 | 113.3 |
| family=reuse;policy=default;refs=8;consumer=100% | valid | 8.7050 | 8.9518 | 9.4269 | 5 | 111.7 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 8.4303 | 8.9113 | 9.7728 | 5 | 112.2 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 9.2600 | 9.9545 | 18.7418 | 5 | 100.5 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 7.5424 | 8.1718 | 9.8552 | 5 | 122.4 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 8.2322 | 8.8633 | 9.8914 | 5 | 112.8 |
| family=reuse;policy=mat;refs=2;consumer=1% | valid | 8.1814 | 8.8764 | 9.2818 | 5 | 112.7 |
| family=reuse;policy=mat;refs=2;consumer=100% | valid | 7.9018 | 8.3026 | 8.8845 | 5 | 120.4 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 8.8125 | 8.9820 | 9.6651 | 5 | 111.3 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 8.7834 | 8.9492 | 9.3558 | 5 | 111.7 |
| family=reuse;policy=mat;refs=8;consumer=1% | valid | 8.6843 | 9.1722 | 9.7600 | 5 | 109.0 |
| family=reuse;policy=mat;refs=8;consumer=100% | valid | 9.0101 | 9.8114 | 10.0397 | 5 | 101.9 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 83.9701 | 89.5516 | 110.7257 | 5 | 11.2 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 85.1362 | 109.1164 | 125.7337 | 5 | 9.2 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 6.9919 | 7.4837 | 8.4423 | 5 | 133.6 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 7.7565 | 8.3393 | 8.9152 | 5 | 119.9 |
| family=reuse;policy=notmat;refs=2;consumer=1% | valid | 11.3808 | 12.7368 | 14.1709 | 5 | 78.5 |
| family=reuse;policy=notmat;refs=2;consumer=100% | valid | 14.1424 | 14.9117 | 15.1724 | 5 | 67.1 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 21.1643 | 21.9993 | 23.5059 | 5 | 45.5 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 29.9078 | 31.2281 | 32.1737 | 5 | 32.0 |
| family=reuse;policy=notmat;refs=8;consumer=1% | valid | 53.1935 | 89.8282 | 100.2613 | 5 | 11.1 |
| family=reuse;policy=notmat;refs=8;consumer=100% | valid | 55.7259 | 56.1452 | 64.3705 | 5 | 17.8 |

## Matched blended workload

- phase: `blended`; target offered rate: 10 updates/s
- offered 526, committed 526, failed 0 over 52.566 s
- offered rate 10.007/s, completed rate 10.007/s
- writer operations per endpoint: [526]
- reader operations overlapping a committed write: 261 of 637
- reader retries of transient conflicts: 0
- per key: 526 keys, min 1, max 1 per key
- ledger: 526 keys reconciled, 0 mismatch
- update latency: avg 3.262 ms, min 1.353, max 146.136 over 526

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42777216`, 5 interleaved rounds
  1. family=bag;stmt=comma_join;threshold=20
  2. family=bag;stmt=filtered_derived;threshold=180000
  3. family=reuse;policy=notmat;refs=4;consumer=1%
  4. family=bag;stmt=comma_join;threshold=2000
  5. family=bag;stmt=join_where;threshold=2000
  6. family=bag;stmt=cte_filtered_mat;threshold=100000
  7. family=bag;stmt=cte_filtered;threshold=180000
  8. family=bag;stmt=join_on;threshold=20000
  9. family=reuse;policy=notmat;refs=1;consumer=1%
  10. family=reuse;policy=default;refs=16;consumer=100%
  11. family=bag;stmt=cte_filtered_mat;threshold=20000
  12. family=reuse;policy=mat;refs=8;consumer=100%
  13. family=bag;stmt=join_on;threshold=20
  14. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  15. family=bag;stmt=cte_filtered_notmat;threshold=20
  16. family=reuse;policy=mat;refs=1;consumer=1%
  17. family=bag;stmt=cte_joined_outer;threshold=20000
  18. family=bag;stmt=cte_joined_outer_mat;threshold=200
  19. family=reuse;policy=default;refs=4;consumer=1%
  20. family=bag;stmt=comma_join;threshold=200
  21. family=bag;stmt=filtered_derived;threshold=20000
  22. family=bag;stmt=join_where;threshold=20
  23. family=bag;stmt=join_where;threshold=100000
  24. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  25. family=reuse;policy=default;refs=1;consumer=1%
  26. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  27. family=bag;stmt=cte_filtered;threshold=20
  28. family=reuse;policy=mat;refs=16;consumer=1%
  29. family=reuse;policy=notmat;refs=2;consumer=100%
  30. family=bag;stmt=comma_join;threshold=20000
  31. family=bag;stmt=cte_filtered_notmat;threshold=200
  32. family=reuse;policy=notmat;refs=16;consumer=1%
  33. family=reuse;policy=notmat;refs=8;consumer=1%
  34. family=reuse;policy=default;refs=16;consumer=1%
  35. family=bag;stmt=filtered_derived;threshold=100000
  36. family=reuse;policy=notmat;refs=8;consumer=100%
  37. family=bag;stmt=cte_filtered_mat;threshold=2000
  38. family=bag;stmt=filtered_derived;threshold=2000
  39. family=bag;stmt=cte_joined_outer;threshold=20
  40. family=reuse;policy=mat;refs=16;consumer=100%
  41. family=bag;stmt=cte_filtered;threshold=100000
  42. family=bag;stmt=cte_joined_outer_mat;threshold=20000
  43. family=reuse;policy=mat;refs=8;consumer=1%
  44. family=bag;stmt=cte_joined_outer;threshold=2000
  45. family=bag;stmt=cte_filtered_mat;threshold=20
  46. family=bag;stmt=cte_filtered_notmat;threshold=20000
  47. family=reuse;policy=mat;refs=2;consumer=1%
  48. family=bag;stmt=cte_filtered;threshold=2000
  49. family=bag;stmt=join_where;threshold=20000
  50. family=reuse;policy=notmat;refs=4;consumer=100%
  51. family=bag;stmt=cte_joined_outer_notmat;threshold=100000
  52. family=bag;stmt=cte_joined_outer;threshold=200
  53. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  54. family=bag;stmt=cte_filtered_notmat;threshold=100000
  55. family=bag;stmt=join_on;threshold=200
  56. family=bag;stmt=cte_filtered_notmat;threshold=2000
  57. family=reuse;policy=mat;refs=2;consumer=100%
  58. family=reuse;policy=notmat;refs=2;consumer=1%
  59. family=bag;stmt=cte_joined_outer_mat;threshold=20
  60. family=reuse;policy=default;refs=8;consumer=100%
  61. family=bag;stmt=cte_filtered_mat;threshold=200
  62. family=bag;stmt=join_where;threshold=180000
  63. family=reuse;policy=mat;refs=1;consumer=100%
  64. family=bag;stmt=cte_joined_outer_mat;threshold=100000
  65. family=bag;stmt=cte_filtered_notmat;threshold=180000
  66. family=bag;stmt=comma_join;threshold=100000
  67. family=bag;stmt=cte_filtered;threshold=200
  68. family=bag;stmt=cte_joined_outer_notmat;threshold=200
  69. family=reuse;policy=default;refs=2;consumer=100%
  70. family=bag;stmt=cte_filtered;threshold=20000
  71. family=reuse;policy=mat;refs=4;consumer=1%
  72. family=bag;stmt=join_on;threshold=2000
  73. family=bag;stmt=cte_joined_outer;threshold=100000
  74. family=bag;stmt=cte_joined_outer;threshold=180000
  75. family=reuse;policy=default;refs=8;consumer=1%
  76. family=bag;stmt=join_where;threshold=200
  77. family=bag;stmt=comma_join;threshold=180000
  78. family=reuse;policy=default;refs=4;consumer=100%
  79. family=bag;stmt=filtered_derived;threshold=200
  80. family=bag;stmt=filtered_derived;threshold=20
  81. family=reuse;policy=notmat;refs=16;consumer=100%
  82. family=reuse;policy=default;refs=2;consumer=1%
  83. family=bag;stmt=cte_joined_outer_notmat;threshold=20000
  84. family=bag;stmt=cte_filtered_mat;threshold=180000
  85. family=reuse;policy=mat;refs=4;consumer=100%
  86. family=reuse;policy=notmat;refs=1;consumer=100%
  87. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  88. family=bag;stmt=join_on;threshold=100000
  89. family=bag;stmt=join_on;threshold=180000
  90. family=reuse;policy=default;refs=1;consumer=100%
  91. family=bag;stmt=comma_join;threshold=20;position=last

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
| family=bag;stmt=comma_join;threshold=100000 | j06_comma_join | valid | `36f94721e8dbe105` |  |
| family=bag;stmt=comma_join;threshold=180000 | j06_comma_join | valid | `3d1485409c32ccab` |  |
| family=bag;stmt=comma_join;threshold=20 | j06_comma_join | valid | `e51dcf6b01b5c55d` |  |
| family=bag;stmt=comma_join;threshold=200 | j06_comma_join | valid | `081431db202ffed8` |  |
| family=bag;stmt=comma_join;threshold=2000 | j06_comma_join | valid | `7479bfe419df6f4e` |  |
| family=bag;stmt=comma_join;threshold=20000 | j06_comma_join | valid | `012d27efd0c5e7ff` |  |
| family=bag;stmt=comma_join;threshold=20;position=last | j06_comma_join | valid | `41c171fbad00e873` |  |
| family=bag;stmt=cte_filtered;threshold=100000 | j04_cte_filtered | valid | `0b2ffc06493a64b9` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=180000 | j04_cte_filtered | valid | `39f7c2cf4e3a312c` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=20 | j04_cte_filtered | valid | `a5cd1ef269cef698` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=200 | j04_cte_filtered | valid | `122bf52f3e1fde3d` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=2000 | j04_cte_filtered | valid | `e5b6ff0ed2c9984b` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=20000 | j04_cte_filtered | valid | `0fceea1374ef14d2` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_mat;threshold=100000 | j04m_cte_filtered_mat | valid | `f5441a3f0cb8902e` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=180000 | j04m_cte_filtered_mat | valid | `fbb6b8ddc813623d` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=20 | j04m_cte_filtered_mat | valid | `be37109f29f6f507` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=200 | j04m_cte_filtered_mat | valid | `e58e6118b1bbe504` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=2000 | j04m_cte_filtered_mat | valid | `04f7b12841e14bad` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=20000 | j04m_cte_filtered_mat | valid | `12c2ca7aefda9bfa` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_notmat;threshold=100000 | j04n_cte_filtered_notmat | valid | `b0c6cefa3c4bd3de` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=180000 | j04n_cte_filtered_notmat | valid | `d8906be3f0ea629f` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=20 | j04n_cte_filtered_notmat | valid | `b26ab71ac1a7bab5` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=200 | j04n_cte_filtered_notmat | valid | `07c5c6c6426a4cc9` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=2000 | j04n_cte_filtered_notmat | valid | `31bdedf8b944bfe3` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=20000 | j04n_cte_filtered_notmat | valid | `3e5af0cae48db3df` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=100000 | j05_cte_joined_outer | valid | `2f42cf9baaed3457` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=180000 | j05_cte_joined_outer | valid | `1919f58cdbb92faa` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20 | j05_cte_joined_outer | valid | `278ed466781d962c` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=200 | j05_cte_joined_outer | valid | `d6b6b726e390f1e8` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=2000 | j05_cte_joined_outer | valid | `f804f15dca69974e` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20000 | j05_cte_joined_outer | valid | `bdd9fa3c77cd7ba4` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_mat;threshold=100000 | j05m_cte_joined_outer_mat | valid | `3aa58cd6c8a4ea0b` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | j05m_cte_joined_outer_mat | valid | `20e172cd892669d8` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | j05m_cte_joined_outer_mat | valid | `cb2676fe32402e9e` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=200 | j05m_cte_joined_outer_mat | valid | `07fab5e25e917a52` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | j05m_cte_joined_outer_mat | valid | `b659d22dce4bf47c` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20000 | j05m_cte_joined_outer_mat | valid | `42d565b597a711a1` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_notmat;threshold=100000 | j05n_cte_joined_outer_notmat | valid | `0540b01b4813046a` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | j05n_cte_joined_outer_notmat | valid | `d81edfbb4a2238b5` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | j05n_cte_joined_outer_notmat | valid | `97679c953efa0038` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=200 | j05n_cte_joined_outer_notmat | valid | `a37bca9bfbb12fce` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | j05n_cte_joined_outer_notmat | valid | `1795f4d9580c37d3` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20000 | j05n_cte_joined_outer_notmat | valid | `03e2e0d56d47f86f` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=filtered_derived;threshold=100000 | j03_filtered_derived | valid | `0adf66044183260f` |  |
| family=bag;stmt=filtered_derived;threshold=180000 | j03_filtered_derived | valid | `9cc5f7dd6f61c716` |  |
| family=bag;stmt=filtered_derived;threshold=20 | j03_filtered_derived | valid | `c3185dfdfdd042a1` |  |
| family=bag;stmt=filtered_derived;threshold=200 | j03_filtered_derived | valid | `17d6440b9fffed5a` |  |
| family=bag;stmt=filtered_derived;threshold=2000 | j03_filtered_derived | valid | `6114139499ac9c02` |  |
| family=bag;stmt=filtered_derived;threshold=20000 | j03_filtered_derived | valid | `64e17026486178e2` |  |
| family=bag;stmt=join_on;threshold=100000 | j02_join_on | valid | `cb114663d96eabde` |  |
| family=bag;stmt=join_on;threshold=180000 | j02_join_on | valid | `2fa37292310e6962` |  |
| family=bag;stmt=join_on;threshold=20 | j02_join_on | valid | `ae4ee8e815e30be2` |  |
| family=bag;stmt=join_on;threshold=200 | j02_join_on | valid | `cd981fc452c8a8d6` |  |
| family=bag;stmt=join_on;threshold=2000 | j02_join_on | valid | `f2863e65b7ed1021` |  |
| family=bag;stmt=join_on;threshold=20000 | j02_join_on | valid | `d9b8ea2a084c04b0` |  |
| family=bag;stmt=join_where;threshold=100000 | j01_join_where | valid | `3aa6630d4c5b48ac` |  |
| family=bag;stmt=join_where;threshold=180000 | j01_join_where | valid | `a80213055018d6ec` |  |
| family=bag;stmt=join_where;threshold=20 | j01_join_where | valid | `febc673d6277009d` |  |
| family=bag;stmt=join_where;threshold=200 | j01_join_where | valid | `ef27ebb75b1eeb93` |  |
| family=bag;stmt=join_where;threshold=2000 | j01_join_where | valid | `3d48c3319e4e9861` |  |
| family=bag;stmt=join_where;threshold=20000 | j01_join_where | valid | `116222d665dd4beb` |  |
| family=reuse;policy=default;refs=16;consumer=1% | r_default_n16_c1 | valid | `ce3db60b8252bd89` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=16;consumer=100% | r_default_n16_c100 | valid | `a85b28ffd2c7ef3f` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=1;consumer=1% | r_default_n1_c1 | valid | `2b13cc877c6bda0a` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=1;consumer=100% | r_default_n1_c100 | valid | `e0694de84b598e72` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=2;consumer=1% | r_default_n2_c1 | valid | `a87de8c863a0d67a` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=2;consumer=100% | r_default_n2_c100 | valid | `1c527585183d2d2e` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=1% | r_default_n4_c1 | valid | `e106c7b3d97f101d` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=100% | r_default_n4_c100 | valid | `f59c9010b87620c2` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=8;consumer=1% | r_default_n8_c1 | valid | `e086fc43ae7af671` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=8;consumer=100% | r_default_n8_c100 | valid | `85a8d384d1b80a31` | policy=default; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=1% | r_mat_n16_c1 | valid | `9fff03bfad2adfbd` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=100% | r_mat_n16_c100 | valid | `a2bcc4581afd1802` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=1% | r_mat_n1_c1 | valid | `86f73a0bc7e2f2b8` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=100% | r_mat_n1_c100 | valid | `b4d11a7446109d2b` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=2;consumer=1% | r_mat_n2_c1 | valid | `6caee16a07cc236b` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=2;consumer=100% | r_mat_n2_c100 | valid | `29d0a4628a80ae67` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=1% | r_mat_n4_c1 | valid | `9d39564c77d9fd35` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=100% | r_mat_n4_c100 | valid | `df2d57d9e7136b35` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=8;consumer=1% | r_mat_n8_c1 | valid | `b32f091ba65ecbed` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=8;consumer=100% | r_mat_n8_c100 | valid | `8153e78542b5b2ce` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=notmat;refs=16;consumer=1% | r_notmat_n16_c1 | valid | `1ea309265e75ee47` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `677ae2bcc8cd9190` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `b9b78fbfeaf8be98` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `9b325d87e73f2655` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=2;consumer=1% | r_notmat_n2_c1 | valid | `7f12da48414cb532` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=2;consumer=100% | r_notmat_n2_c100 | valid | `fdc569290a0a5b20` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `83183e20bd36a187` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `0d9d53dadb4e295b` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=8;consumer=1% | r_notmat_n8_c1 | valid | `cde1d48c78542fc4` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=8;consumer=100% | r_notmat_n8_c100 | valid | `f85c30e87bdb864a` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, DIST, DEBUG, FORMAT TEXT` supported on yugabyte
- one host: this run used node budget 2 CPU 3g and client budget 1 CPU 2g, so budget comparisons are separate classes here, not independent hardware
- one client process with one fixed maintenance writer at 10 updates/s; readers are sequential; the reported sample latency includes retry and backoff time
- physical rank-index arm yb-range-rank; an engine whose planner ignores the index is a different physical design, not a faster one
- YugabyteDB DIST/DEBUG storage counters (rows and index rows scanned) are parsed per cell; RPC-level counters are not aggregated
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
