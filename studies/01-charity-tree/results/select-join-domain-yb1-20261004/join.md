# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 91 cells, 91 valid, engine `yugabyte`, topology `single`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-domain-yb1-20261004`
- run tag: `run/study-01/select-join-domain-yb1-20261004`
- environment: `host-zenbook-ux5406sa`
- repo commit: `4e249caed0279c6fe647214665b593c72b2e91d2`
- engine version: `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit`
- database engine image: `docker.io/yugabytedb/yugabyte:2025.2.6.0-b111` @ `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`
- benchmark client image: `localhost/charitybench-select:select-join-domain-yb1-20261004` @ `fa3606dc9ba1e4169c0124c33434b7aa827570eb10aa634f479901ca10b2f0be`
- inputs digest: `c3e4599c16117a8868ad19e9831a16a5b026c096bf511f6e115f0bec45f5902a`
- dataset fingerprint: `fcf26e08f377ec51`
- client CPU quota: `100000 100000`, memory: `2147483648`
- setting `client_workers`: 1 sequential reader; 0 writers
- setting `connection_pool`: 2 = 1 reader/setup + 1 reserve
- setting `customers`: 2000
- setting `dataset`: study01-person-donation-charity
- setting `distributed_rpc`: not parsed; raw EXPLAIN retained (coverage gap)
- setting `fact_count`: 200000
- setting `index_set`: child rank index and parent-key index, primary keys; no per-spelling index
- setting `prepared_plan_policy`: pgx simple protocol; no prepared statements
- setting `scale`: cte-domain
- setting `selectivity_mapping`: rank <= threshold; rank is a deterministic shuffle of 1..fact_count
- setting `statement_budget_ms`: 120000
- setting `statistics`: ANALYZE after deterministic load
- setting `thresholds`: [20 200 2000 20000 100000 180000]
- setting `topology`: single
- setting `trials_per_cell`: 5
- setting `whole_run_count`: 1; no replicated winner in this file

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
| 20 (0.01%) | 105.8785 | 76.8356 | 98.3247 | 78.4165 | 71.1056 | 79.3469 | 79.1075 | 188.3907 | 70.4831 | 80.1080 |
| 200 (0.1%) | 73.1880 | 70.6495 | 90.3503 | 92.6889 | 93.0677 | 84.4803 | 81.3453 | 200.7583 | 82.3018 | 74.1292 |
| 2000 (1%) | 76.7624 | 77.8327 | 91.5311 | 72.9088 | 84.5401 | 84.4173 | 73.2526 | 201.0103 | 100.7849 | 68.7936 |
| 20000 (10%) | 93.8910 | 87.3950 | 88.0785 | 98.1257 | 126.7666 | 88.4380 | 97.0167 | 148.2526 | 96.6108 | 126.1999 |
| 100000 (50%) | 120.5154 | 118.4268 | 143.6310 | 135.5453 | 155.3808 | 129.8132 | 193.0855 | 233.9574 | 118.9194 | 123.0144 |
| 180000 (90%) | 161.5619 | 175.2461 | 174.2613 | 170.7941 | 170.7817 | 166.3519 | 191.6006 | 220.1819 | 185.1344 | 220.4352 |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | join_on vs join_where | 105.8785 | 76.8356 | 0.7257 | false |
| 20 | filtered_derived vs join_where | 105.8785 | 98.3247 | 0.9287 | false |
| 20 | cte_filtered vs join_where | 105.8785 | 78.4165 | 0.7406 | false |
| 20 | cte_filtered_mat vs join_where | 105.8785 | 71.1056 | 0.6716 | false |
| 20 | cte_filtered_notmat vs join_where | 105.8785 | 79.3469 | 0.7494 | false |
| 20 | cte_joined_outer vs join_where | 105.8785 | 79.1075 | 0.7472 | false |
| 20 | cte_joined_outer_mat vs join_where | 105.8785 | 188.3907 | 1.7793 | false |
| 20 | cte_joined_outer_notmat vs join_where | 105.8785 | 70.4831 | 0.6657 | false |
| 20 | comma_join vs join_where | 105.8785 | 80.1080 | 0.7566 | false |
| 200 | join_on vs join_where | 73.1880 | 70.6495 | 0.9653 | false |
| 200 | filtered_derived vs join_where | 73.1880 | 90.3503 | 1.2345 | false |
| 200 | cte_filtered vs join_where | 73.1880 | 92.6889 | 1.2665 | false |
| 200 | cte_filtered_mat vs join_where | 73.1880 | 93.0677 | 1.2716 | false |
| 200 | cte_filtered_notmat vs join_where | 73.1880 | 84.4803 | 1.1543 | false |
| 200 | cte_joined_outer vs join_where | 73.1880 | 81.3453 | 1.1115 | false |
| 200 | cte_joined_outer_mat vs join_where | 73.1880 | 200.7583 | 2.7431 | false |
| 200 | cte_joined_outer_notmat vs join_where | 73.1880 | 82.3018 | 1.1245 | false |
| 200 | comma_join vs join_where | 73.1880 | 74.1292 | 1.0129 | false |
| 2000 | join_on vs join_where | 76.7624 | 77.8327 | 1.0139 | false |
| 2000 | filtered_derived vs join_where | 76.7624 | 91.5311 | 1.1924 | false |
| 2000 | cte_filtered vs join_where | 76.7624 | 72.9088 | 0.9498 | false |
| 2000 | cte_filtered_mat vs join_where | 76.7624 | 84.5401 | 1.1013 | false |
| 2000 | cte_filtered_notmat vs join_where | 76.7624 | 84.4173 | 1.0997 | false |
| 2000 | cte_joined_outer vs join_where | 76.7624 | 73.2526 | 0.9543 | false |
| 2000 | cte_joined_outer_mat vs join_where | 76.7624 | 201.0103 | 2.6186 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 76.7624 | 100.7849 | 1.3129 | false |
| 2000 | comma_join vs join_where | 76.7624 | 68.7936 | 0.8962 | false |
| 20000 | join_on vs join_where | 93.8910 | 87.3950 | 0.9308 | false |
| 20000 | filtered_derived vs join_where | 93.8910 | 88.0785 | 0.9381 | false |
| 20000 | cte_filtered vs join_where | 93.8910 | 98.1257 | 1.0451 | false |
| 20000 | cte_filtered_mat vs join_where | 93.8910 | 126.7666 | 1.3501 | false |
| 20000 | cte_filtered_notmat vs join_where | 93.8910 | 88.4380 | 0.9419 | false |
| 20000 | cte_joined_outer vs join_where | 93.8910 | 97.0167 | 1.0333 | false |
| 20000 | cte_joined_outer_mat vs join_where | 93.8910 | 148.2526 | 1.5790 | false |
| 20000 | cte_joined_outer_notmat vs join_where | 93.8910 | 96.6108 | 1.0290 | false |
| 20000 | comma_join vs join_where | 93.8910 | 126.1999 | 1.3441 | false |
| 100000 | join_on vs join_where | 120.5154 | 118.4268 | 0.9827 | false |
| 100000 | filtered_derived vs join_where | 120.5154 | 143.6310 | 1.1918 | false |
| 100000 | cte_filtered vs join_where | 120.5154 | 135.5453 | 1.1247 | false |
| 100000 | cte_filtered_mat vs join_where | 120.5154 | 155.3808 | 1.2893 | false |
| 100000 | cte_filtered_notmat vs join_where | 120.5154 | 129.8132 | 1.0772 | false |
| 100000 | cte_joined_outer vs join_where | 120.5154 | 193.0855 | 1.6022 | false |
| 100000 | cte_joined_outer_mat vs join_where | 120.5154 | 233.9574 | 1.9413 | false |
| 100000 | cte_joined_outer_notmat vs join_where | 120.5154 | 118.9194 | 0.9868 | false |
| 100000 | comma_join vs join_where | 120.5154 | 123.0144 | 1.0207 | false |
| 180000 | join_on vs join_where | 161.5619 | 175.2461 | 1.0847 | false |
| 180000 | filtered_derived vs join_where | 161.5619 | 174.2613 | 1.0786 | false |
| 180000 | cte_filtered vs join_where | 161.5619 | 170.7941 | 1.0571 | false |
| 180000 | cte_filtered_mat vs join_where | 161.5619 | 170.7817 | 1.0571 | false |
| 180000 | cte_filtered_notmat vs join_where | 161.5619 | 166.3519 | 1.0296 | false |
| 180000 | cte_joined_outer vs join_where | 161.5619 | 191.6006 | 1.1859 | false |
| 180000 | cte_joined_outer_mat vs join_where | 161.5619 | 220.1819 | 1.3628 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 161.5619 | 185.1344 | 1.1459 | false |
| 180000 | comma_join vs join_where | 161.5619 | 220.4352 | 1.3644 | false |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_filtered_mat vs cte_filtered | 78.4165 | 71.1056 | 0.9068 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_filtered_notmat vs cte_filtered | 78.4165 | 79.3469 | 1.0119 | false | policy=not_materialized; observed_cte_scan=false |
| 20 | cte_joined_outer_mat vs cte_joined_outer | 79.1075 | 188.3907 | 2.3815 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 79.1075 | 70.4831 | 0.8910 | false | policy=not_materialized; observed_cte_scan=false |
| 200 | cte_filtered_mat vs cte_filtered | 92.6889 | 93.0677 | 1.0041 | false | policy=materialized; observed_cte_scan=true |
| 200 | cte_filtered_notmat vs cte_filtered | 92.6889 | 84.4803 | 0.9114 | false | policy=not_materialized; observed_cte_scan=false |
| 200 | cte_joined_outer_mat vs cte_joined_outer | 81.3453 | 200.7583 | 2.4680 | false | policy=materialized; observed_cte_scan=true |
| 200 | cte_joined_outer_notmat vs cte_joined_outer | 81.3453 | 82.3018 | 1.0118 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_filtered_mat vs cte_filtered | 72.9088 | 84.5401 | 1.1595 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_filtered_notmat vs cte_filtered | 72.9088 | 84.4173 | 1.1578 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 73.2526 | 201.0103 | 2.7441 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 73.2526 | 100.7849 | 1.3759 | false | policy=not_materialized; observed_cte_scan=false |
| 20000 | cte_filtered_mat vs cte_filtered | 98.1257 | 126.7666 | 1.2919 | false | policy=materialized; observed_cte_scan=true |
| 20000 | cte_filtered_notmat vs cte_filtered | 98.1257 | 88.4380 | 0.9013 | false | policy=not_materialized; observed_cte_scan=false |
| 20000 | cte_joined_outer_mat vs cte_joined_outer | 97.0167 | 148.2526 | 1.5281 | false | policy=materialized; observed_cte_scan=true |
| 20000 | cte_joined_outer_notmat vs cte_joined_outer | 97.0167 | 96.6108 | 0.9958 | false | policy=not_materialized; observed_cte_scan=false |
| 100000 | cte_filtered_mat vs cte_filtered | 135.5453 | 155.3808 | 1.1463 | false | policy=materialized; observed_cte_scan=true |
| 100000 | cte_filtered_notmat vs cte_filtered | 135.5453 | 129.8132 | 0.9577 | false | policy=not_materialized; observed_cte_scan=false |
| 100000 | cte_joined_outer_mat vs cte_joined_outer | 193.0855 | 233.9574 | 1.2117 | false | policy=materialized; observed_cte_scan=true |
| 100000 | cte_joined_outer_notmat vs cte_joined_outer | 193.0855 | 118.9194 | 0.6159 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_filtered_mat vs cte_filtered | 170.7941 | 170.7817 | 0.9999 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_filtered_notmat vs cte_filtered | 170.7941 | 166.3519 | 0.9740 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 191.6006 | 220.1819 | 1.1492 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 191.6006 | 185.1344 | 0.9663 | false | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 74.3600 | 1 | `6ec8b6ef41670287` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 82.7491 | 1 | `ff79a352bb19d78b` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | 90.2165 | 2 | `8a59fa39b4065a59` | policy=default; observed_cte_scan=true | valid |
| default | 2 | 100% | 78.0111 | 2 | `767a1d0748588b41` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 1% | 91.3477 | 4 | `6688211ed8b87367` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 80.0318 | 4 | `99ea2be309be62f6` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | 86.2296 | 8 | `edcbbd3ca9e775bf` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 100% | 74.0249 | 8 | `3409d235d9dfb791` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 1% | 108.3874 | 16 | `4bda64629d563d8f` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 88.5839 | 16 | `bdf1228a11b92459` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 78.6449 | 1 | `c7153b3ff13fba9a` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 76.3532 | 1 | `af2b99761aee5d3e` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | 77.1046 | 2 | `441d44a095a9224d` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 100% | 95.3046 | 2 | `276f1d69e32c200e` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 1% | 84.3075 | 4 | `99ebebe3adc2c9e4` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 75.4162 | 4 | `9650472847f49c1f` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | 70.0125 | 8 | `7693d0c0a2d54734` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 100% | 76.7111 | 8 | `e1102e3f5214b125` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 1% | 75.6634 | 16 | `84e22830b01edd10` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 78.4442 | 16 | `ae8fdea8f950608c` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 81.0691 | 1 | `31c82f866e0681e3` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 79.1151 | 1 | `ce8bd2912759d0a0` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | 243.0256 | 2 | `055a50b82dc0dd8d` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 100% | 201.6763 | 2 | `845d07ddec0c6ebd` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 1% | 285.9821 | 4 | `5e2f75e98ca8e796` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 257.7474 | 4 | `084a1fdf707448ab` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | 566.6426 | 8 | `d8ae049c6fe2d047` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 100% | 538.3197 | 8 | `41e5eb2b074a6a63` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 1% | 979.1275 | 16 | `9e38d00b18b2bd84` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 1157.2493 | 16 | `1fa172cf77acbac1` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 1.0576 | false |
| notmat | 1 | 1% | 1.0902 | false |
| mat | 1 | 100% | 0.9227 | false |
| notmat | 1 | 100% | 0.9561 | false |
| mat | 2 | 1% | 0.8547 | false |
| notmat | 2 | 1% | 2.6938 | false |
| mat | 2 | 100% | 1.2217 | false |
| notmat | 2 | 100% | 2.5852 | false |
| mat | 4 | 1% | 0.9229 | false |
| notmat | 4 | 1% | 3.1307 | false |
| mat | 4 | 100% | 0.9423 | false |
| notmat | 4 | 100% | 3.2206 | false |
| mat | 8 | 1% | 0.8119 | false |
| notmat | 8 | 1% | 6.5713 | false |
| mat | 8 | 100% | 1.0363 | false |
| notmat | 8 | 100% | 7.2721 | false |
| mat | 16 | 1% | 0.6981 | false |
| notmat | 16 | 1% | 9.0336 | false |
| mat | 16 | 100% | 0.8855 | false |
| notmat | 16 | 100% | 13.0639 | false |

## Repeated-design position control

| first median ms | last median ms | ratio (last/first) | same plan |
|---|---|---|---|
| 175.246114 | 188.569504 | 1.0760 | false |

## Latency samples and derived throughput

_The reader is sequential, so `queries/s` is 1000/median_ms, a derived rate and not a concurrent-capacity test. min/max are the five sample extremes._

| axis | status | min ms | median ms | max ms | samples | derived queries/s |
|---|---|---|---|---|---|---|
| family=bag;stmt=comma_join;threshold=100000 | valid | 108.4621 | 123.0144 | 211.8221 | 5 | 8.1 |
| family=bag;stmt=comma_join;threshold=180000 | valid | 144.2953 | 220.4352 | 252.5918 | 5 | 4.5 |
| family=bag;stmt=comma_join;threshold=20 | valid | 67.4194 | 80.1080 | 109.1537 | 5 | 12.5 |
| family=bag;stmt=comma_join;threshold=200 | valid | 71.0271 | 74.1292 | 79.4192 | 5 | 13.5 |
| family=bag;stmt=comma_join;threshold=2000 | valid | 65.0174 | 68.7936 | 117.6297 | 5 | 14.5 |
| family=bag;stmt=comma_join;threshold=20000 | valid | 93.8390 | 126.1999 | 192.4069 | 5 | 7.9 |
| family=bag;stmt=cte_filtered;threshold=100000 | valid | 115.6487 | 135.5453 | 184.2503 | 5 | 7.4 |
| family=bag;stmt=cte_filtered;threshold=180000 | valid | 143.2106 | 170.7941 | 210.3800 | 5 | 5.9 |
| family=bag;stmt=cte_filtered;threshold=20 | valid | 69.0592 | 78.4165 | 122.3579 | 5 | 12.8 |
| family=bag;stmt=cte_filtered;threshold=200 | valid | 66.2782 | 92.6889 | 142.5819 | 5 | 10.8 |
| family=bag;stmt=cte_filtered;threshold=2000 | valid | 70.5059 | 72.9088 | 82.9565 | 5 | 13.7 |
| family=bag;stmt=cte_filtered;threshold=20000 | valid | 81.9952 | 98.1257 | 130.5954 | 5 | 10.2 |
| family=bag;stmt=cte_filtered_mat;threshold=100000 | valid | 121.5704 | 155.3808 | 218.3339 | 5 | 6.4 |
| family=bag;stmt=cte_filtered_mat;threshold=180000 | valid | 161.6522 | 170.7817 | 208.3747 | 5 | 5.9 |
| family=bag;stmt=cte_filtered_mat;threshold=20 | valid | 67.8623 | 71.1056 | 102.2652 | 5 | 14.1 |
| family=bag;stmt=cte_filtered_mat;threshold=200 | valid | 73.0385 | 93.0677 | 97.6678 | 5 | 10.7 |
| family=bag;stmt=cte_filtered_mat;threshold=2000 | valid | 73.2292 | 84.5401 | 150.6191 | 5 | 11.8 |
| family=bag;stmt=cte_filtered_mat;threshold=20000 | valid | 85.9660 | 126.7666 | 131.6050 | 5 | 7.9 |
| family=bag;stmt=cte_filtered_notmat;threshold=100000 | valid | 122.5818 | 129.8132 | 176.3935 | 5 | 7.7 |
| family=bag;stmt=cte_filtered_notmat;threshold=180000 | valid | 161.4104 | 166.3519 | 291.0726 | 5 | 6.0 |
| family=bag;stmt=cte_filtered_notmat;threshold=20 | valid | 69.7267 | 79.3469 | 106.5919 | 5 | 12.6 |
| family=bag;stmt=cte_filtered_notmat;threshold=200 | valid | 69.6562 | 84.4803 | 105.1984 | 5 | 11.8 |
| family=bag;stmt=cte_filtered_notmat;threshold=2000 | valid | 72.8228 | 84.4173 | 121.1704 | 5 | 11.8 |
| family=bag;stmt=cte_filtered_notmat;threshold=20000 | valid | 78.7122 | 88.4380 | 119.8487 | 5 | 11.3 |
| family=bag;stmt=cte_joined_outer;threshold=100000 | valid | 116.8227 | 193.0855 | 214.0105 | 5 | 5.2 |
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 181.3059 | 191.6006 | 306.0953 | 5 | 5.2 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 75.3415 | 79.1075 | 94.2056 | 5 | 12.6 |
| family=bag;stmt=cte_joined_outer;threshold=200 | valid | 70.3148 | 81.3453 | 124.9729 | 5 | 12.3 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 70.6699 | 73.2526 | 145.2482 | 5 | 13.7 |
| family=bag;stmt=cte_joined_outer;threshold=20000 | valid | 79.2594 | 97.0167 | 125.4073 | 5 | 10.3 |
| family=bag;stmt=cte_joined_outer_mat;threshold=100000 | valid | 143.3495 | 233.9574 | 263.0308 | 5 | 4.3 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 169.1780 | 220.1819 | 267.9369 | 5 | 4.5 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 122.6661 | 188.3907 | 238.3668 | 5 | 5.3 |
| family=bag;stmt=cte_joined_outer_mat;threshold=200 | valid | 114.2947 | 200.7583 | 271.2229 | 5 | 5.0 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 112.8199 | 201.0103 | 217.7086 | 5 | 5.0 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20000 | valid | 136.0929 | 148.2526 | 233.2037 | 5 | 6.7 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=100000 | valid | 114.3497 | 118.9194 | 140.0751 | 5 | 8.4 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 159.1564 | 185.1344 | 229.7439 | 5 | 5.4 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 67.3184 | 70.4831 | 96.2605 | 5 | 14.2 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=200 | valid | 73.9942 | 82.3018 | 95.8087 | 5 | 12.2 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 72.2145 | 100.7849 | 120.7038 | 5 | 9.9 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20000 | valid | 78.1484 | 96.6108 | 174.2411 | 5 | 10.4 |
| family=bag;stmt=filtered_derived;threshold=100000 | valid | 119.5888 | 143.6310 | 200.2466 | 5 | 7.0 |
| family=bag;stmt=filtered_derived;threshold=180000 | valid | 163.4923 | 174.2613 | 253.0807 | 5 | 5.7 |
| family=bag;stmt=filtered_derived;threshold=20 | valid | 69.7077 | 98.3247 | 109.9986 | 5 | 10.2 |
| family=bag;stmt=filtered_derived;threshold=200 | valid | 73.6603 | 90.3503 | 99.0857 | 5 | 11.1 |
| family=bag;stmt=filtered_derived;threshold=2000 | valid | 70.1549 | 91.5311 | 120.5044 | 5 | 10.9 |
| family=bag;stmt=filtered_derived;threshold=20000 | valid | 81.0200 | 88.0785 | 124.3264 | 5 | 11.4 |
| family=bag;stmt=join_on;threshold=100000 | valid | 115.5503 | 118.4268 | 199.5247 | 5 | 8.4 |
| family=bag;stmt=join_on;threshold=180000 | valid | 141.1656 | 175.2461 | 253.0165 | 5 | 5.7 |
| family=bag;stmt=join_on;threshold=180000;position=last | valid | 166.3044 | 188.5695 | 262.5184 | 5 | 5.3 |
| family=bag;stmt=join_on;threshold=20 | valid | 68.7260 | 76.8356 | 103.2077 | 5 | 13.0 |
| family=bag;stmt=join_on;threshold=200 | valid | 62.6058 | 70.6495 | 99.8900 | 5 | 14.2 |
| family=bag;stmt=join_on;threshold=2000 | valid | 74.7077 | 77.8327 | 122.2634 | 5 | 12.8 |
| family=bag;stmt=join_on;threshold=20000 | valid | 82.0913 | 87.3950 | 129.3258 | 5 | 11.4 |
| family=bag;stmt=join_where;threshold=100000 | valid | 110.0645 | 120.5154 | 186.8895 | 5 | 8.3 |
| family=bag;stmt=join_where;threshold=180000 | valid | 153.3947 | 161.5619 | 240.9407 | 5 | 6.2 |
| family=bag;stmt=join_where;threshold=20 | valid | 90.2168 | 105.8785 | 110.5053 | 5 | 9.4 |
| family=bag;stmt=join_where;threshold=200 | valid | 67.7113 | 73.1880 | 98.8265 | 5 | 13.7 |
| family=bag;stmt=join_where;threshold=2000 | valid | 67.6810 | 76.7624 | 105.9817 | 5 | 13.0 |
| family=bag;stmt=join_where;threshold=20000 | valid | 82.4200 | 93.8910 | 136.0736 | 5 | 10.7 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 69.7530 | 108.3874 | 112.3135 | 5 | 9.2 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 77.8201 | 88.5839 | 99.8085 | 5 | 11.3 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 71.0893 | 74.3600 | 104.3355 | 5 | 13.4 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 74.5282 | 82.7491 | 135.7788 | 5 | 12.1 |
| family=reuse;policy=default;refs=2;consumer=1% | valid | 71.3378 | 90.2165 | 137.3403 | 5 | 11.1 |
| family=reuse;policy=default;refs=2;consumer=100% | valid | 67.8377 | 78.0111 | 105.7442 | 5 | 12.8 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 69.0811 | 91.3477 | 117.3120 | 5 | 10.9 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 67.5916 | 80.0318 | 102.1877 | 5 | 12.5 |
| family=reuse;policy=default;refs=8;consumer=1% | valid | 69.5340 | 86.2296 | 152.2154 | 5 | 11.6 |
| family=reuse;policy=default;refs=8;consumer=100% | valid | 68.9520 | 74.0249 | 111.9430 | 5 | 13.5 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 70.1091 | 75.6634 | 120.0494 | 5 | 13.2 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 74.1561 | 78.4442 | 134.8297 | 5 | 12.7 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 74.0199 | 78.6449 | 91.0353 | 5 | 12.7 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 74.7808 | 76.3532 | 82.9662 | 5 | 13.1 |
| family=reuse;policy=mat;refs=2;consumer=1% | valid | 72.4744 | 77.1046 | 134.1430 | 5 | 13.0 |
| family=reuse;policy=mat;refs=2;consumer=100% | valid | 74.5961 | 95.3046 | 115.6546 | 5 | 10.5 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 69.5870 | 84.3075 | 116.9175 | 5 | 11.9 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 72.8812 | 75.4162 | 110.7067 | 5 | 13.3 |
| family=reuse;policy=mat;refs=8;consumer=1% | valid | 61.0549 | 70.0125 | 95.9639 | 5 | 14.3 |
| family=reuse;policy=mat;refs=8;consumer=100% | valid | 71.0829 | 76.7111 | 122.3271 | 5 | 13.0 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 961.7414 | 979.1275 | 1180.1948 | 5 | 1.0 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 983.5691 | 1157.2493 | 1352.6502 | 5 | 0.9 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 67.6220 | 81.0691 | 106.5409 | 5 | 12.3 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 72.1381 | 79.1151 | 107.0771 | 5 | 12.6 |
| family=reuse;policy=notmat;refs=2;consumer=1% | valid | 168.4377 | 243.0256 | 291.5500 | 5 | 4.1 |
| family=reuse;policy=notmat;refs=2;consumer=100% | valid | 149.4775 | 201.6763 | 247.9705 | 5 | 5.0 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 243.4778 | 285.9821 | 466.4656 | 5 | 3.5 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 164.9903 | 257.7474 | 427.5675 | 5 | 3.9 |
| family=reuse;policy=notmat;refs=8;consumer=1% | valid | 433.8209 | 566.6426 | 674.7388 | 5 | 1.8 |
| family=reuse;policy=notmat;refs=8;consumer=100% | valid | 458.1956 | 538.3197 | 784.2502 | 5 | 1.9 |

## Execution order

- order seed `42620866`, 5 interleaved rounds
  1. family=bag;stmt=join_on;threshold=180000
  2. family=bag;stmt=join_where;threshold=20000
  3. family=bag;stmt=cte_filtered_mat;threshold=100000
  4. family=reuse;policy=default;refs=8;consumer=1%
  5. family=bag;stmt=filtered_derived;threshold=100000
  6. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  7. family=bag;stmt=join_on;threshold=100000
  8. family=bag;stmt=comma_join;threshold=20
  9. family=reuse;policy=notmat;refs=1;consumer=100%
  10. family=bag;stmt=join_on;threshold=200
  11. family=bag;stmt=comma_join;threshold=2000
  12. family=reuse;policy=default;refs=8;consumer=100%
  13. family=bag;stmt=join_where;threshold=2000
  14. family=reuse;policy=default;refs=2;consumer=100%
  15. family=bag;stmt=cte_filtered_notmat;threshold=180000
  16. family=bag;stmt=join_where;threshold=100000
  17. family=bag;stmt=cte_filtered;threshold=2000
  18. family=bag;stmt=cte_joined_outer_notmat;threshold=100000
  19. family=reuse;policy=notmat;refs=2;consumer=1%
  20. family=bag;stmt=cte_filtered_notmat;threshold=20
  21. family=reuse;policy=mat;refs=8;consumer=1%
  22. family=bag;stmt=join_where;threshold=200
  23. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  24. family=bag;stmt=cte_filtered_mat;threshold=20
  25. family=reuse;policy=default;refs=1;consumer=1%
  26. family=bag;stmt=filtered_derived;threshold=20000
  27. family=bag;stmt=join_on;threshold=20
  28. family=reuse;policy=default;refs=16;consumer=100%
  29. family=reuse;policy=default;refs=4;consumer=100%
  30. family=bag;stmt=cte_joined_outer;threshold=20
  31. family=reuse;policy=mat;refs=1;consumer=1%
  32. family=bag;stmt=filtered_derived;threshold=200
  33. family=bag;stmt=cte_filtered_notmat;threshold=20000
  34. family=bag;stmt=join_where;threshold=180000
  35. family=reuse;policy=mat;refs=4;consumer=1%
  36. family=reuse;policy=notmat;refs=4;consumer=1%
  37. family=bag;stmt=cte_joined_outer;threshold=200
  38. family=bag;stmt=cte_joined_outer_mat;threshold=20000
  39. family=bag;stmt=cte_filtered_mat;threshold=200
  40. family=bag;stmt=cte_filtered_mat;threshold=2000
  41. family=reuse;policy=notmat;refs=8;consumer=100%
  42. family=reuse;policy=mat;refs=16;consumer=1%
  43. family=reuse;policy=mat;refs=8;consumer=100%
  44. family=reuse;policy=mat;refs=16;consumer=100%
  45. family=reuse;policy=mat;refs=4;consumer=100%
  46. family=reuse;policy=mat;refs=1;consumer=100%
  47. family=bag;stmt=comma_join;threshold=200
  48. family=bag;stmt=cte_filtered_mat;threshold=180000
  49. family=bag;stmt=comma_join;threshold=100000
  50. family=bag;stmt=join_on;threshold=20000
  51. family=bag;stmt=cte_joined_outer_notmat;threshold=200
  52. family=bag;stmt=filtered_derived;threshold=180000
  53. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  54. family=bag;stmt=cte_filtered;threshold=180000
  55. family=bag;stmt=cte_filtered;threshold=100000
  56. family=reuse;policy=notmat;refs=16;consumer=1%
  57. family=reuse;policy=notmat;refs=16;consumer=100%
  58. family=bag;stmt=cte_filtered_notmat;threshold=200
  59. family=bag;stmt=join_on;threshold=2000
  60. family=bag;stmt=cte_filtered;threshold=20000
  61. family=reuse;policy=notmat;refs=4;consumer=100%
  62. family=bag;stmt=comma_join;threshold=20000
  63. family=reuse;policy=notmat;refs=8;consumer=1%
  64. family=bag;stmt=cte_joined_outer_mat;threshold=20
  65. family=bag;stmt=comma_join;threshold=180000
  66. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  67. family=reuse;policy=notmat;refs=2;consumer=100%
  68. family=bag;stmt=cte_filtered;threshold=200
  69. family=bag;stmt=cte_joined_outer_mat;threshold=100000
  70. family=bag;stmt=join_where;threshold=20
  71. family=bag;stmt=cte_filtered_mat;threshold=20000
  72. family=bag;stmt=filtered_derived;threshold=20
  73. family=reuse;policy=default;refs=16;consumer=1%
  74. family=bag;stmt=cte_joined_outer;threshold=100000
  75. family=bag;stmt=cte_joined_outer_mat;threshold=200
  76. family=bag;stmt=cte_filtered;threshold=20
  77. family=reuse;policy=mat;refs=2;consumer=1%
  78. family=bag;stmt=cte_joined_outer;threshold=180000
  79. family=bag;stmt=cte_joined_outer;threshold=20000
  80. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  81. family=reuse;policy=default;refs=2;consumer=1%
  82. family=bag;stmt=cte_filtered_notmat;threshold=100000
  83. family=bag;stmt=cte_joined_outer;threshold=2000
  84. family=reuse;policy=default;refs=1;consumer=100%
  85. family=bag;stmt=filtered_derived;threshold=2000
  86. family=reuse;policy=mat;refs=2;consumer=100%
  87. family=bag;stmt=cte_joined_outer_notmat;threshold=20000
  88. family=reuse;policy=default;refs=4;consumer=1%
  89. family=reuse;policy=notmat;refs=1;consumer=1%
  90. family=bag;stmt=cte_filtered_notmat;threshold=2000
  91. family=bag;stmt=join_on;threshold=180000;position=last

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
| family=bag;stmt=comma_join;threshold=100000 | j06_comma_join | valid | `ba9e8a39d5a8d3f3` |  |
| family=bag;stmt=comma_join;threshold=180000 | j06_comma_join | valid | `1007b090fb1e9b09` |  |
| family=bag;stmt=comma_join;threshold=20 | j06_comma_join | valid | `4c766e3ea312feda` |  |
| family=bag;stmt=comma_join;threshold=200 | j06_comma_join | valid | `8ea2b71b868a3db8` |  |
| family=bag;stmt=comma_join;threshold=2000 | j06_comma_join | valid | `7da53698194b3c1a` |  |
| family=bag;stmt=comma_join;threshold=20000 | j06_comma_join | valid | `b97b072864d6ab79` |  |
| family=bag;stmt=cte_filtered;threshold=100000 | j04_cte_filtered | valid | `e90ee7e47792ddb1` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=180000 | j04_cte_filtered | valid | `69d7da82bbe4e05e` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=20 | j04_cte_filtered | valid | `80a32a38e312bdb2` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=200 | j04_cte_filtered | valid | `19664f1520ee06ca` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=2000 | j04_cte_filtered | valid | `8d328bd583fe655c` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=20000 | j04_cte_filtered | valid | `83870b831b53ff49` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_mat;threshold=100000 | j04m_cte_filtered_mat | valid | `79564af5a3b13050` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=180000 | j04m_cte_filtered_mat | valid | `c179588db8781255` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=20 | j04m_cte_filtered_mat | valid | `bd94917a6913f338` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=200 | j04m_cte_filtered_mat | valid | `f4da85b34e86bf3a` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=2000 | j04m_cte_filtered_mat | valid | `692b34b63968d9b3` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=20000 | j04m_cte_filtered_mat | valid | `fd7976b0392b4a0c` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_notmat;threshold=100000 | j04n_cte_filtered_notmat | valid | `4451871da3064865` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=180000 | j04n_cte_filtered_notmat | valid | `a5ef76b28f36dc20` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=20 | j04n_cte_filtered_notmat | valid | `426ee5b629774c36` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=200 | j04n_cte_filtered_notmat | valid | `2d71f6d0010e2f0a` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=2000 | j04n_cte_filtered_notmat | valid | `4f112ad380b8a3d5` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=20000 | j04n_cte_filtered_notmat | valid | `32d19d25845b5726` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=100000 | j05_cte_joined_outer | valid | `9e21ccc5034b34bc` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=180000 | j05_cte_joined_outer | valid | `fb45012fb07b5880` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20 | j05_cte_joined_outer | valid | `93d599727c9bc072` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=200 | j05_cte_joined_outer | valid | `5f79d32b637c4fb5` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=2000 | j05_cte_joined_outer | valid | `1e57d1b3e6c01142` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20000 | j05_cte_joined_outer | valid | `0c23e46e79182e35` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_mat;threshold=100000 | j05m_cte_joined_outer_mat | valid | `79063a9a543ce1db` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | j05m_cte_joined_outer_mat | valid | `f16ddab6a79b454c` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | j05m_cte_joined_outer_mat | valid | `fdebc827bfbad417` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=200 | j05m_cte_joined_outer_mat | valid | `94d26c50f75f006e` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | j05m_cte_joined_outer_mat | valid | `5b61435b2a352184` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20000 | j05m_cte_joined_outer_mat | valid | `c1c9d807ca16d685` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_notmat;threshold=100000 | j05n_cte_joined_outer_notmat | valid | `4034085f03b7766e` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | j05n_cte_joined_outer_notmat | valid | `58d88ccb9e1a6fad` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | j05n_cte_joined_outer_notmat | valid | `21199a28660c3aa9` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=200 | j05n_cte_joined_outer_notmat | valid | `8d4f4f851b149315` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | j05n_cte_joined_outer_notmat | valid | `1cbfcabded1c2b57` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20000 | j05n_cte_joined_outer_notmat | valid | `e22e28cde46317e8` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=filtered_derived;threshold=100000 | j03_filtered_derived | valid | `56b5d1c9b9e5b3f7` |  |
| family=bag;stmt=filtered_derived;threshold=180000 | j03_filtered_derived | valid | `6d0af5828a2aece7` |  |
| family=bag;stmt=filtered_derived;threshold=20 | j03_filtered_derived | valid | `3495920bbe9fef68` |  |
| family=bag;stmt=filtered_derived;threshold=200 | j03_filtered_derived | valid | `794a61eb6023fbad` |  |
| family=bag;stmt=filtered_derived;threshold=2000 | j03_filtered_derived | valid | `9fed90339f5f7f44` |  |
| family=bag;stmt=filtered_derived;threshold=20000 | j03_filtered_derived | valid | `dec1a2e09cdc244a` |  |
| family=bag;stmt=join_on;threshold=100000 | j02_join_on | valid | `943a2db28e4e493e` |  |
| family=bag;stmt=join_on;threshold=180000 | j02_join_on | valid | `f6a1f88fc5ece574` |  |
| family=bag;stmt=join_on;threshold=180000;position=last | j02_join_on | valid | `a9bb7e42d80dc899` |  |
| family=bag;stmt=join_on;threshold=20 | j02_join_on | valid | `4ec97963449d86aa` |  |
| family=bag;stmt=join_on;threshold=200 | j02_join_on | valid | `70d78727db5c9c6d` |  |
| family=bag;stmt=join_on;threshold=2000 | j02_join_on | valid | `6082d40086d0da96` |  |
| family=bag;stmt=join_on;threshold=20000 | j02_join_on | valid | `0e7d2d4d444f4669` |  |
| family=bag;stmt=join_where;threshold=100000 | j01_join_where | valid | `b84ddebb8fb66942` |  |
| family=bag;stmt=join_where;threshold=180000 | j01_join_where | valid | `86ca6c3b7e5cf09b` |  |
| family=bag;stmt=join_where;threshold=20 | j01_join_where | valid | `45f44c72083fbfb4` |  |
| family=bag;stmt=join_where;threshold=200 | j01_join_where | valid | `5fa44f807afaf131` |  |
| family=bag;stmt=join_where;threshold=2000 | j01_join_where | valid | `c9b6ad352d00aa28` |  |
| family=bag;stmt=join_where;threshold=20000 | j01_join_where | valid | `5949e518446410ce` |  |
| family=reuse;policy=default;refs=16;consumer=1% | r_default_n16_c1 | valid | `4bda64629d563d8f` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=16;consumer=100% | r_default_n16_c100 | valid | `bdf1228a11b92459` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=1;consumer=1% | r_default_n1_c1 | valid | `6ec8b6ef41670287` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=1;consumer=100% | r_default_n1_c100 | valid | `ff79a352bb19d78b` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=2;consumer=1% | r_default_n2_c1 | valid | `8a59fa39b4065a59` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=2;consumer=100% | r_default_n2_c100 | valid | `767a1d0748588b41` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=1% | r_default_n4_c1 | valid | `6688211ed8b87367` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=100% | r_default_n4_c100 | valid | `99ea2be309be62f6` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=8;consumer=1% | r_default_n8_c1 | valid | `edcbbd3ca9e775bf` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=8;consumer=100% | r_default_n8_c100 | valid | `3409d235d9dfb791` | policy=default; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=1% | r_mat_n16_c1 | valid | `84e22830b01edd10` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=100% | r_mat_n16_c100 | valid | `ae8fdea8f950608c` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=1% | r_mat_n1_c1 | valid | `c7153b3ff13fba9a` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=100% | r_mat_n1_c100 | valid | `af2b99761aee5d3e` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=2;consumer=1% | r_mat_n2_c1 | valid | `441d44a095a9224d` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=2;consumer=100% | r_mat_n2_c100 | valid | `276f1d69e32c200e` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=1% | r_mat_n4_c1 | valid | `99ebebe3adc2c9e4` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=100% | r_mat_n4_c100 | valid | `9650472847f49c1f` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=8;consumer=1% | r_mat_n8_c1 | valid | `7693d0c0a2d54734` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=8;consumer=100% | r_mat_n8_c100 | valid | `e1102e3f5214b125` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=notmat;refs=16;consumer=1% | r_notmat_n16_c1 | valid | `9e38d00b18b2bd84` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `1fa172cf77acbac1` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `31c82f866e0681e3` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `ce8bd2912759d0a0` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=2;consumer=1% | r_notmat_n2_c1 | valid | `055a50b82dc0dd8d` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=2;consumer=100% | r_notmat_n2_c100 | valid | `845d07ddec0c6ebd` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `5e2f75e98ca8e796` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `084a1fdf707448ab` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=8;consumer=1% | r_notmat_n8_c1 | valid | `d8ae049c6fe2d047` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=8;consumer=100% | r_notmat_n8_c100 | valid | `41e5eb2b074a6a63` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, DIST, DEBUG, FORMAT TEXT` supported on yugabyte
- one host, single client: no independent hardware, real-network or equal-total-budget comparison
- no concurrent writers: the isolated read contract only; a matched blended read/write phase is a recorded gap
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
