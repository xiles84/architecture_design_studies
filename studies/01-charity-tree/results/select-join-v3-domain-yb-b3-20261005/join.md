# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 91 cells, 91 valid, engine `yugabyte`, topology `single`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-v3-domain-yb-b3-20261005`
- run tag: `run/study-01/select-join-v3-domain-yb-b3-20261005`
- environment: `host-zenbook-ux5406sa`
- repo commit: `0f9d09212c6d52a80d8437d3b5a234f2f711ed36`
- engine version: `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit`
- database engine image: `docker.io/yugabytedb/yugabyte:2025.2.6.0-b111` @ `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`
- benchmark client image: `localhost/charitybench-select:select-join-v3-domain-yb-b3-20261005` @ `8436e36bb66caa8178559d54771b88c063e26848d4ef1d1a56eb524b7f0ecca6`
- inputs digest: `189b76d013d7777cceab041c2ecbea04499cba2d1c4b8641ca086b20f2c11189`
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
| 20 (0.01%) | 1.8276 | 2.3302 | 1.9853 | 1.8431 | 2.1152 | 1.8220 | 1.8464 | 80.8148 | 1.6937 | 1.7951 |
| 200 (0.1%) | 2.8353 | 2.9722 | 3.1216 | 3.4871 | 2.7394 | 3.1720 | 2.8459 | 83.1849 | 3.7779 | 3.0302 |
| 2000 (1%) | 9.6383 | 8.4901 | 9.4889 | 8.6093 | 9.8716 | 9.8899 | 9.7670 | 88.4376 | 10.0059 | 9.0935 |
| 20000 (10%) | 64.9114 | 66.6933 | 64.3321 | 62.6322 | 66.4493 | 64.2873 | 63.8291 | 87.7118 | 64.7567 | 62.3088 |
| 100000 (50%) | 90.8649 | 85.5925 | 87.3267 | 87.9869 | 95.0479 | 86.3328 | 85.4274 | 96.4324 | 86.6986 | 85.7539 |
| 180000 (90%) | 123.0986 | 110.6048 | 110.2073 | 107.5930 | 111.1159 | 109.2454 | 108.8021 | 102.2577 | 106.5249 | 108.2053 |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | join_on vs join_where | 1.8276 | 2.3302 | 1.2750 | false |
| 20 | filtered_derived vs join_where | 1.8276 | 1.9853 | 1.0863 | false |
| 20 | cte_filtered vs join_where | 1.8276 | 1.8431 | 1.0085 | false |
| 20 | cte_filtered_mat vs join_where | 1.8276 | 2.1152 | 1.1574 | false |
| 20 | cte_filtered_notmat vs join_where | 1.8276 | 1.8220 | 0.9969 | false |
| 20 | cte_joined_outer vs join_where | 1.8276 | 1.8464 | 1.0103 | false |
| 20 | cte_joined_outer_mat vs join_where | 1.8276 | 80.8148 | 44.2201 | false |
| 20 | cte_joined_outer_notmat vs join_where | 1.8276 | 1.6937 | 0.9267 | false |
| 20 | comma_join vs join_where | 1.8276 | 1.7951 | 0.9823 | false |
| 200 | join_on vs join_where | 2.8353 | 2.9722 | 1.0483 | false |
| 200 | filtered_derived vs join_where | 2.8353 | 3.1216 | 1.1010 | false |
| 200 | cte_filtered vs join_where | 2.8353 | 3.4871 | 1.2299 | false |
| 200 | cte_filtered_mat vs join_where | 2.8353 | 2.7394 | 0.9662 | false |
| 200 | cte_filtered_notmat vs join_where | 2.8353 | 3.1720 | 1.1188 | false |
| 200 | cte_joined_outer vs join_where | 2.8353 | 2.8459 | 1.0037 | false |
| 200 | cte_joined_outer_mat vs join_where | 2.8353 | 83.1849 | 29.3389 | false |
| 200 | cte_joined_outer_notmat vs join_where | 2.8353 | 3.7779 | 1.3324 | false |
| 200 | comma_join vs join_where | 2.8353 | 3.0302 | 1.0688 | false |
| 2000 | join_on vs join_where | 9.6383 | 8.4901 | 0.8809 | false |
| 2000 | filtered_derived vs join_where | 9.6383 | 9.4889 | 0.9845 | false |
| 2000 | cte_filtered vs join_where | 9.6383 | 8.6093 | 0.8932 | false |
| 2000 | cte_filtered_mat vs join_where | 9.6383 | 9.8716 | 1.0242 | false |
| 2000 | cte_filtered_notmat vs join_where | 9.6383 | 9.8899 | 1.0261 | false |
| 2000 | cte_joined_outer vs join_where | 9.6383 | 9.7670 | 1.0134 | false |
| 2000 | cte_joined_outer_mat vs join_where | 9.6383 | 88.4376 | 9.1756 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 9.6383 | 10.0059 | 1.0381 | false |
| 2000 | comma_join vs join_where | 9.6383 | 9.0935 | 0.9435 | false |
| 20000 | join_on vs join_where | 64.9114 | 66.6933 | 1.0275 | false |
| 20000 | filtered_derived vs join_where | 64.9114 | 64.3321 | 0.9911 | false |
| 20000 | cte_filtered vs join_where | 64.9114 | 62.6322 | 0.9649 | false |
| 20000 | cte_filtered_mat vs join_where | 64.9114 | 66.4493 | 1.0237 | false |
| 20000 | cte_filtered_notmat vs join_where | 64.9114 | 64.2873 | 0.9904 | false |
| 20000 | cte_joined_outer vs join_where | 64.9114 | 63.8291 | 0.9833 | false |
| 20000 | cte_joined_outer_mat vs join_where | 64.9114 | 87.7118 | 1.3513 | false |
| 20000 | cte_joined_outer_notmat vs join_where | 64.9114 | 64.7567 | 0.9976 | false |
| 20000 | comma_join vs join_where | 64.9114 | 62.3088 | 0.9599 | false |
| 100000 | join_on vs join_where | 90.8649 | 85.5925 | 0.9420 | false |
| 100000 | filtered_derived vs join_where | 90.8649 | 87.3267 | 0.9611 | false |
| 100000 | cte_filtered vs join_where | 90.8649 | 87.9869 | 0.9683 | false |
| 100000 | cte_filtered_mat vs join_where | 90.8649 | 95.0479 | 1.0460 | false |
| 100000 | cte_filtered_notmat vs join_where | 90.8649 | 86.3328 | 0.9501 | false |
| 100000 | cte_joined_outer vs join_where | 90.8649 | 85.4274 | 0.9402 | false |
| 100000 | cte_joined_outer_mat vs join_where | 90.8649 | 96.4324 | 1.0613 | false |
| 100000 | cte_joined_outer_notmat vs join_where | 90.8649 | 86.6986 | 0.9541 | false |
| 100000 | comma_join vs join_where | 90.8649 | 85.7539 | 0.9438 | false |
| 180000 | join_on vs join_where | 123.0986 | 110.6048 | 0.8985 | false |
| 180000 | filtered_derived vs join_where | 123.0986 | 110.2073 | 0.8953 | false |
| 180000 | cte_filtered vs join_where | 123.0986 | 107.5930 | 0.8740 | false |
| 180000 | cte_filtered_mat vs join_where | 123.0986 | 111.1159 | 0.9027 | false |
| 180000 | cte_filtered_notmat vs join_where | 123.0986 | 109.2454 | 0.8875 | false |
| 180000 | cte_joined_outer vs join_where | 123.0986 | 108.8021 | 0.8839 | false |
| 180000 | cte_joined_outer_mat vs join_where | 123.0986 | 102.2577 | 0.8307 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 123.0986 | 106.5249 | 0.8654 | false |
| 180000 | comma_join vs join_where | 123.0986 | 108.2053 | 0.8790 | false |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_filtered_mat vs cte_filtered | 1.8431 | 2.1152 | 1.1476 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_filtered_notmat vs cte_filtered | 1.8431 | 1.8220 | 0.9885 | false | policy=not_materialized; observed_cte_scan=false |
| 20 | cte_joined_outer_mat vs cte_joined_outer | 1.8464 | 80.8148 | 43.7696 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 1.8464 | 1.6937 | 0.9173 | false | policy=not_materialized; observed_cte_scan=false |
| 200 | cte_filtered_mat vs cte_filtered | 3.4871 | 2.7394 | 0.7856 | false | policy=materialized; observed_cte_scan=true |
| 200 | cte_filtered_notmat vs cte_filtered | 3.4871 | 3.1720 | 0.9096 | false | policy=not_materialized; observed_cte_scan=false |
| 200 | cte_joined_outer_mat vs cte_joined_outer | 2.8459 | 83.1849 | 29.2295 | false | policy=materialized; observed_cte_scan=true |
| 200 | cte_joined_outer_notmat vs cte_joined_outer | 2.8459 | 3.7779 | 1.3275 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_filtered_mat vs cte_filtered | 8.6093 | 9.8716 | 1.1466 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_filtered_notmat vs cte_filtered | 8.6093 | 9.8899 | 1.1488 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 9.7670 | 88.4376 | 9.0547 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 9.7670 | 10.0059 | 1.0245 | false | policy=not_materialized; observed_cte_scan=false |
| 20000 | cte_filtered_mat vs cte_filtered | 62.6322 | 66.4493 | 1.0609 | false | policy=materialized; observed_cte_scan=true |
| 20000 | cte_filtered_notmat vs cte_filtered | 62.6322 | 64.2873 | 1.0264 | false | policy=not_materialized; observed_cte_scan=false |
| 20000 | cte_joined_outer_mat vs cte_joined_outer | 63.8291 | 87.7118 | 1.3742 | false | policy=materialized; observed_cte_scan=true |
| 20000 | cte_joined_outer_notmat vs cte_joined_outer | 63.8291 | 64.7567 | 1.0145 | false | policy=not_materialized; observed_cte_scan=false |
| 100000 | cte_filtered_mat vs cte_filtered | 87.9869 | 95.0479 | 1.0803 | false | policy=materialized; observed_cte_scan=true |
| 100000 | cte_filtered_notmat vs cte_filtered | 87.9869 | 86.3328 | 0.9812 | false | policy=not_materialized; observed_cte_scan=false |
| 100000 | cte_joined_outer_mat vs cte_joined_outer | 85.4274 | 96.4324 | 1.1288 | false | policy=materialized; observed_cte_scan=true |
| 100000 | cte_joined_outer_notmat vs cte_joined_outer | 85.4274 | 86.6986 | 1.0149 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_filtered_mat vs cte_filtered | 107.5930 | 111.1159 | 1.0327 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_filtered_notmat vs cte_filtered | 107.5930 | 109.2454 | 1.0154 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 108.8021 | 102.2577 | 0.9399 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 108.8021 | 106.5249 | 0.9791 | false | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 7.6191 | 1 | `c6faa2b2a434a4ff` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 9.5224 | 1 | `1334f16b54340058` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | 8.7575 | 2 | `80c694ac69713bcf` | policy=default; observed_cte_scan=true | valid |
| default | 2 | 100% | 9.1384 | 2 | `cb5c42885c210034` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 1% | 9.5914 | 4 | `8d99023cbb8afdc7` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 9.0019 | 4 | `2496c33e3db1ea40` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | 8.8349 | 8 | `4eadbe18a4648146` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 100% | 9.7558 | 8 | `6f36ce56f2d4b0d5` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 1% | 8.7819 | 16 | `c0eff00634f64a4c` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 9.8725 | 16 | `9979f8fab58a1d21` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 7.8923 | 1 | `f4b8330c95f75e83` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 8.0605 | 1 | `e0b5d89ff8646f81` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | 8.2690 | 2 | `872459181a59e5c2` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 100% | 9.6435 | 2 | `92a6aaa309d0ce5c` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 1% | 8.0938 | 4 | `a0fd8e4360a67939` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 9.1732 | 4 | `7083ae8059b5fe80` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | 8.6092 | 8 | `30f7e07a2edc465b` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 100% | 9.1172 | 8 | `a5a60931cb48284a` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 1% | 9.6179 | 16 | `95e5790ee3291b8e` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 10.2686 | 16 | `07f1ca3b642ce096` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 7.4001 | 1 | `6b0896655055ba83` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 7.9445 | 1 | `bbbac4d17d7adc05` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | 12.0151 | 2 | `eaa463088ee68259` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 100% | 15.4806 | 2 | `f5ba60e75a457580` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 1% | 25.2586 | 4 | `acaec9b309b9ca63` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 30.5821 | 4 | `f9b3f5d54c85f804` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | 60.2314 | 8 | `e3680db6113c66ab` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 100% | 60.3749 | 8 | `a92b70780ebec35a` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 1% | 92.1133 | 16 | `52fd9c5cd35a3a9c` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 90.8527 | 16 | `edaee8bff2874abd` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 1.0359 | false |
| notmat | 1 | 1% | 0.9713 | false |
| mat | 1 | 100% | 0.8465 | false |
| notmat | 1 | 100% | 0.8343 | false |
| mat | 2 | 1% | 0.9442 | false |
| notmat | 2 | 1% | 1.3720 | false |
| mat | 2 | 100% | 1.0553 | false |
| notmat | 2 | 100% | 1.6940 | false |
| mat | 4 | 1% | 0.8439 | false |
| notmat | 4 | 1% | 2.6335 | false |
| mat | 4 | 100% | 1.0190 | false |
| notmat | 4 | 100% | 3.3973 | false |
| mat | 8 | 1% | 0.9745 | false |
| notmat | 8 | 1% | 6.8174 | false |
| mat | 8 | 100% | 0.9345 | false |
| notmat | 8 | 100% | 6.1886 | false |
| mat | 16 | 1% | 1.0952 | false |
| notmat | 16 | 1% | 10.4890 | false |
| mat | 16 | 100% | 1.0401 | false |
| notmat | 16 | 100% | 9.2026 | false |

## Repeated-design position control

| first median ms | last median ms | ratio (last/first) | same plan |
|---|---|---|---|
| 9.767017 | 9.236876 | 0.9457 | false |

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

- client: nr_periods 632, nr_throttled 253, throttled_usec 9865085 | quota 100000 100000
, memory 2147483648

- yb-single: nr_periods 696, nr_throttled 77, throttled_usec 2694500 | quota 200000 100000, memory 3221225472

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
| family=bag;stmt=cte_joined_outer;threshold=2000;position=last | valid | true | true | 6050 | 2000 | 0/0 | 0 | 0 |
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
| family=bag;stmt=comma_join;threshold=100000 | valid | 82.1933 | 85.7539 | 86.8291 | 5 | 11.7 |
| family=bag;stmt=comma_join;threshold=180000 | valid | 104.5065 | 108.2053 | 111.9823 | 5 | 9.2 |
| family=bag;stmt=comma_join;threshold=20 | valid | 1.4795 | 1.7951 | 2.2531 | 5 | 557.1 |
| family=bag;stmt=comma_join;threshold=200 | valid | 2.9319 | 3.0302 | 3.1234 | 5 | 330.0 |
| family=bag;stmt=comma_join;threshold=2000 | valid | 8.8783 | 9.0935 | 9.7160 | 5 | 110.0 |
| family=bag;stmt=comma_join;threshold=20000 | valid | 60.8984 | 62.3088 | 63.6823 | 5 | 16.0 |
| family=bag;stmt=cte_filtered;threshold=100000 | valid | 86.7699 | 87.9869 | 89.7488 | 5 | 11.4 |
| family=bag;stmt=cte_filtered;threshold=180000 | valid | 103.9212 | 107.5930 | 114.5849 | 5 | 9.3 |
| family=bag;stmt=cte_filtered;threshold=20 | valid | 1.7010 | 1.8431 | 2.0370 | 5 | 542.6 |
| family=bag;stmt=cte_filtered;threshold=200 | valid | 3.4601 | 3.4871 | 3.8352 | 5 | 286.8 |
| family=bag;stmt=cte_filtered;threshold=2000 | valid | 8.4165 | 8.6093 | 9.0252 | 5 | 116.2 |
| family=bag;stmt=cte_filtered;threshold=20000 | valid | 61.2502 | 62.6322 | 66.0703 | 5 | 16.0 |
| family=bag;stmt=cte_filtered_mat;threshold=100000 | valid | 90.5602 | 95.0479 | 120.1717 | 5 | 10.5 |
| family=bag;stmt=cte_filtered_mat;threshold=180000 | valid | 110.6700 | 111.1159 | 124.6182 | 5 | 9.0 |
| family=bag;stmt=cte_filtered_mat;threshold=20 | valid | 2.0697 | 2.1152 | 2.1866 | 5 | 472.8 |
| family=bag;stmt=cte_filtered_mat;threshold=200 | valid | 2.6341 | 2.7394 | 3.3474 | 5 | 365.0 |
| family=bag;stmt=cte_filtered_mat;threshold=2000 | valid | 9.3790 | 9.8716 | 11.7535 | 5 | 101.3 |
| family=bag;stmt=cte_filtered_mat;threshold=20000 | valid | 66.2709 | 66.4493 | 68.3884 | 5 | 15.0 |
| family=bag;stmt=cte_filtered_notmat;threshold=100000 | valid | 84.4707 | 86.3328 | 87.3343 | 5 | 11.6 |
| family=bag;stmt=cte_filtered_notmat;threshold=180000 | valid | 103.4964 | 109.2454 | 112.0706 | 5 | 9.2 |
| family=bag;stmt=cte_filtered_notmat;threshold=20 | valid | 1.7197 | 1.8220 | 2.1645 | 5 | 548.9 |
| family=bag;stmt=cte_filtered_notmat;threshold=200 | valid | 2.9090 | 3.1720 | 3.7811 | 5 | 315.3 |
| family=bag;stmt=cte_filtered_notmat;threshold=2000 | valid | 9.4283 | 9.8899 | 12.1084 | 5 | 101.1 |
| family=bag;stmt=cte_filtered_notmat;threshold=20000 | valid | 62.0949 | 64.2873 | 66.2359 | 5 | 15.6 |
| family=bag;stmt=cte_joined_outer;threshold=100000 | valid | 83.8290 | 85.4274 | 89.6541 | 5 | 11.7 |
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 104.1501 | 108.8021 | 110.7507 | 5 | 9.2 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 1.6790 | 1.8464 | 2.0393 | 5 | 541.6 |
| family=bag;stmt=cte_joined_outer;threshold=200 | valid | 2.7314 | 2.8459 | 30.0434 | 5 | 351.4 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 8.6217 | 9.7670 | 11.5506 | 5 | 102.4 |
| family=bag;stmt=cte_joined_outer;threshold=20000 | valid | 62.6236 | 63.8291 | 64.9600 | 5 | 15.7 |
| family=bag;stmt=cte_joined_outer;threshold=2000;position=last | valid | 9.1065 | 9.2369 | 9.6756 | 5 | 108.3 |
| family=bag;stmt=cte_joined_outer_mat;threshold=100000 | valid | 91.0417 | 96.4324 | 100.2840 | 5 | 10.4 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 101.6619 | 102.2577 | 110.5747 | 5 | 9.8 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 80.2148 | 80.8148 | 82.7569 | 5 | 12.4 |
| family=bag;stmt=cte_joined_outer_mat;threshold=200 | valid | 81.8302 | 83.1849 | 86.1652 | 5 | 12.0 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 80.5232 | 88.4376 | 91.6776 | 5 | 11.3 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20000 | valid | 86.5677 | 87.7118 | 89.8340 | 5 | 11.4 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=100000 | valid | 83.6461 | 86.6986 | 96.5284 | 5 | 11.5 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 104.2144 | 106.5249 | 110.0019 | 5 | 9.4 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 1.5082 | 1.6937 | 2.0035 | 5 | 590.4 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=200 | valid | 3.4431 | 3.7779 | 3.9689 | 5 | 264.7 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 9.6570 | 10.0059 | 10.1175 | 5 | 99.9 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20000 | valid | 62.6627 | 64.7567 | 69.7272 | 5 | 15.4 |
| family=bag;stmt=filtered_derived;threshold=100000 | valid | 84.6284 | 87.3267 | 92.5684 | 5 | 11.5 |
| family=bag;stmt=filtered_derived;threshold=180000 | valid | 105.4663 | 110.2073 | 115.0534 | 5 | 9.1 |
| family=bag;stmt=filtered_derived;threshold=20 | valid | 1.7835 | 1.9853 | 2.0589 | 5 | 503.7 |
| family=bag;stmt=filtered_derived;threshold=200 | valid | 2.9545 | 3.1216 | 3.8261 | 5 | 320.3 |
| family=bag;stmt=filtered_derived;threshold=2000 | valid | 8.4845 | 9.4889 | 11.0434 | 5 | 105.4 |
| family=bag;stmt=filtered_derived;threshold=20000 | valid | 62.8587 | 64.3321 | 68.5262 | 5 | 15.5 |
| family=bag;stmt=join_on;threshold=100000 | valid | 84.5734 | 85.5925 | 107.2712 | 5 | 11.7 |
| family=bag;stmt=join_on;threshold=180000 | valid | 105.9529 | 110.6048 | 114.1205 | 5 | 9.0 |
| family=bag;stmt=join_on;threshold=20 | valid | 2.0899 | 2.3302 | 2.5112 | 5 | 429.2 |
| family=bag;stmt=join_on;threshold=200 | valid | 2.5792 | 2.9722 | 3.0060 | 5 | 336.4 |
| family=bag;stmt=join_on;threshold=2000 | valid | 8.3844 | 8.4901 | 9.3964 | 5 | 117.8 |
| family=bag;stmt=join_on;threshold=20000 | valid | 64.2005 | 66.6933 | 72.6933 | 5 | 15.0 |
| family=bag;stmt=join_where;threshold=100000 | valid | 88.9423 | 90.8649 | 92.6718 | 5 | 11.0 |
| family=bag;stmt=join_where;threshold=180000 | valid | 110.1193 | 123.0986 | 128.5324 | 5 | 8.1 |
| family=bag;stmt=join_where;threshold=20 | valid | 1.7644 | 1.8276 | 2.0191 | 5 | 547.2 |
| family=bag;stmt=join_where;threshold=200 | valid | 2.6485 | 2.8353 | 3.1866 | 5 | 352.7 |
| family=bag;stmt=join_where;threshold=2000 | valid | 9.3598 | 9.6383 | 10.2071 | 5 | 103.8 |
| family=bag;stmt=join_where;threshold=20000 | valid | 62.7310 | 64.9114 | 66.9529 | 5 | 15.4 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 8.5469 | 8.7819 | 9.6559 | 5 | 113.9 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 9.5168 | 9.8725 | 25.6594 | 5 | 101.3 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 7.2218 | 7.6191 | 7.6764 | 5 | 131.2 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 8.5002 | 9.5224 | 9.8153 | 5 | 105.0 |
| family=reuse;policy=default;refs=2;consumer=1% | valid | 8.3726 | 8.7575 | 12.4110 | 5 | 114.2 |
| family=reuse;policy=default;refs=2;consumer=100% | valid | 8.6447 | 9.1384 | 9.6118 | 5 | 109.4 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 9.1212 | 9.5914 | 10.1767 | 5 | 104.3 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 8.5715 | 9.0019 | 10.2102 | 5 | 111.1 |
| family=reuse;policy=default;refs=8;consumer=1% | valid | 8.4867 | 8.8349 | 9.1985 | 5 | 113.2 |
| family=reuse;policy=default;refs=8;consumer=100% | valid | 8.9804 | 9.7558 | 10.0240 | 5 | 102.5 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 9.5658 | 9.6179 | 9.7267 | 5 | 104.0 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 10.1282 | 10.2686 | 10.5768 | 5 | 97.4 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 7.6440 | 7.8923 | 8.4727 | 5 | 126.7 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 7.4352 | 8.0605 | 9.0919 | 5 | 124.1 |
| family=reuse;policy=mat;refs=2;consumer=1% | valid | 8.0717 | 8.2690 | 8.5441 | 5 | 120.9 |
| family=reuse;policy=mat;refs=2;consumer=100% | valid | 9.4021 | 9.6435 | 10.4799 | 5 | 103.7 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 7.8538 | 8.0938 | 8.2682 | 5 | 123.6 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 8.8767 | 9.1732 | 9.5246 | 5 | 109.0 |
| family=reuse;policy=mat;refs=8;consumer=1% | valid | 8.5909 | 8.6092 | 9.0278 | 5 | 116.2 |
| family=reuse;policy=mat;refs=8;consumer=100% | valid | 8.8203 | 9.1172 | 9.6391 | 5 | 109.7 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 85.8591 | 92.1133 | 143.7843 | 5 | 10.9 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 82.4558 | 90.8527 | 110.6884 | 5 | 11.0 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 7.1191 | 7.4001 | 7.9734 | 5 | 135.1 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 7.3657 | 7.9445 | 8.0693 | 5 | 125.9 |
| family=reuse;policy=notmat;refs=2;consumer=1% | valid | 10.8348 | 12.0151 | 12.4197 | 5 | 83.2 |
| family=reuse;policy=notmat;refs=2;consumer=100% | valid | 14.3409 | 15.4806 | 15.8721 | 5 | 64.6 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 23.2268 | 25.2586 | 25.5324 | 5 | 39.6 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 27.9185 | 30.5821 | 31.4130 | 5 | 32.7 |
| family=reuse;policy=notmat;refs=8;consumer=1% | valid | 57.7223 | 60.2314 | 91.2037 | 5 | 16.6 |
| family=reuse;policy=notmat;refs=8;consumer=100% | valid | 59.3387 | 60.3749 | 73.2533 | 5 | 16.6 |

## Matched blended workload

- phase: `blended`; target offered rate: 10 updates/s
- offered 546, committed 546, failed 0 over 54.595 s
- offered rate 10.001/s, completed rate 10.001/s
- writer operations per endpoint: [546]
- reader operations overlapping a committed write: 240 of 637
- reader retries of transient conflicts: 0
- per key: 546 keys, min 1, max 1 per key
- ledger: 546 keys reconciled, 0 mismatch
- update latency: avg 3.043 ms, min 1.416, max 221.559 over 546

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42891112`, 5 interleaved rounds
  1. family=bag;stmt=cte_joined_outer;threshold=2000
  2. family=bag;stmt=cte_filtered;threshold=180000
  3. family=reuse;policy=default;refs=1;consumer=100%
  4. family=bag;stmt=join_where;threshold=20000
  5. family=bag;stmt=cte_filtered_notmat;threshold=180000
  6. family=bag;stmt=cte_joined_outer;threshold=20000
  7. family=reuse;policy=default;refs=2;consumer=100%
  8. family=reuse;policy=notmat;refs=16;consumer=1%
  9. family=reuse;policy=mat;refs=16;consumer=1%
  10. family=reuse;policy=mat;refs=2;consumer=1%
  11. family=bag;stmt=cte_joined_outer_mat;threshold=20
  12. family=reuse;policy=mat;refs=8;consumer=100%
  13. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  14. family=reuse;policy=notmat;refs=4;consumer=1%
  15. family=bag;stmt=join_where;threshold=20
  16. family=bag;stmt=cte_filtered_notmat;threshold=100000
  17. family=reuse;policy=notmat;refs=16;consumer=100%
  18. family=bag;stmt=join_where;threshold=180000
  19. family=bag;stmt=cte_filtered_mat;threshold=180000
  20. family=bag;stmt=cte_filtered_notmat;threshold=2000
  21. family=reuse;policy=default;refs=8;consumer=1%
  22. family=reuse;policy=notmat;refs=2;consumer=1%
  23. family=bag;stmt=comma_join;threshold=20000
  24. family=reuse;policy=default;refs=8;consumer=100%
  25. family=bag;stmt=comma_join;threshold=200
  26. family=bag;stmt=comma_join;threshold=2000
  27. family=reuse;policy=mat;refs=8;consumer=1%
  28. family=reuse;policy=notmat;refs=2;consumer=100%
  29. family=reuse;policy=mat;refs=1;consumer=1%
  30. family=bag;stmt=join_on;threshold=2000
  31. family=bag;stmt=cte_filtered;threshold=2000
  32. family=bag;stmt=comma_join;threshold=180000
  33. family=bag;stmt=comma_join;threshold=100000
  34. family=reuse;policy=mat;refs=16;consumer=100%
  35. family=reuse;policy=mat;refs=4;consumer=1%
  36. family=reuse;policy=mat;refs=1;consumer=100%
  37. family=bag;stmt=cte_filtered;threshold=20000
  38. family=bag;stmt=cte_joined_outer;threshold=180000
  39. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  40. family=bag;stmt=cte_joined_outer;threshold=100000
  41. family=bag;stmt=filtered_derived;threshold=20000
  42. family=bag;stmt=cte_filtered_mat;threshold=20000
  43. family=reuse;policy=mat;refs=4;consumer=100%
  44. family=bag;stmt=cte_filtered;threshold=20
  45. family=bag;stmt=join_on;threshold=100000
  46. family=reuse;policy=notmat;refs=8;consumer=100%
  47. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  48. family=bag;stmt=join_where;threshold=2000
  49. family=bag;stmt=cte_filtered_notmat;threshold=20
  50. family=bag;stmt=cte_joined_outer_mat;threshold=200
  51. family=bag;stmt=filtered_derived;threshold=20
  52. family=reuse;policy=notmat;refs=8;consumer=1%
  53. family=bag;stmt=cte_filtered;threshold=200
  54. family=bag;stmt=cte_joined_outer;threshold=200
  55. family=bag;stmt=cte_filtered_mat;threshold=100000
  56. family=reuse;policy=mat;refs=2;consumer=100%
  57. family=reuse;policy=default;refs=4;consumer=100%
  58. family=bag;stmt=filtered_derived;threshold=200
  59. family=bag;stmt=join_where;threshold=200
  60. family=bag;stmt=cte_filtered_mat;threshold=200
  61. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  62. family=reuse;policy=notmat;refs=4;consumer=100%
  63. family=bag;stmt=cte_joined_outer_mat;threshold=100000
  64. family=reuse;policy=default;refs=4;consumer=1%
  65. family=bag;stmt=cte_filtered;threshold=100000
  66. family=bag;stmt=cte_joined_outer_mat;threshold=20000
  67. family=bag;stmt=cte_filtered_mat;threshold=20
  68. family=bag;stmt=cte_filtered_mat;threshold=2000
  69. family=reuse;policy=default;refs=2;consumer=1%
  70. family=bag;stmt=join_where;threshold=100000
  71. family=bag;stmt=join_on;threshold=20
  72. family=bag;stmt=join_on;threshold=20000
  73. family=bag;stmt=filtered_derived;threshold=180000
  74. family=bag;stmt=join_on;threshold=180000
  75. family=bag;stmt=cte_joined_outer_notmat;threshold=200
  76. family=bag;stmt=cte_joined_outer_notmat;threshold=100000
  77. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  78. family=reuse;policy=default;refs=16;consumer=100%
  79. family=bag;stmt=filtered_derived;threshold=2000
  80. family=bag;stmt=cte_joined_outer_notmat;threshold=20000
  81. family=reuse;policy=notmat;refs=1;consumer=1%
  82. family=reuse;policy=default;refs=16;consumer=1%
  83. family=reuse;policy=notmat;refs=1;consumer=100%
  84. family=bag;stmt=filtered_derived;threshold=100000
  85. family=reuse;policy=default;refs=1;consumer=1%
  86. family=bag;stmt=cte_joined_outer;threshold=20
  87. family=bag;stmt=comma_join;threshold=20
  88. family=bag;stmt=cte_filtered_notmat;threshold=200
  89. family=bag;stmt=join_on;threshold=200
  90. family=bag;stmt=cte_filtered_notmat;threshold=20000
  91. family=bag;stmt=cte_joined_outer;threshold=2000;position=last

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
| family=bag;stmt=comma_join;threshold=100000 | j06_comma_join | valid | `69e5c5a8d7c8fa8e` |  |
| family=bag;stmt=comma_join;threshold=180000 | j06_comma_join | valid | `5076949c3aba6e6b` |  |
| family=bag;stmt=comma_join;threshold=20 | j06_comma_join | valid | `cdaa71c504b52d02` |  |
| family=bag;stmt=comma_join;threshold=200 | j06_comma_join | valid | `49a6920c807cebf5` |  |
| family=bag;stmt=comma_join;threshold=2000 | j06_comma_join | valid | `944dc89980ac61c0` |  |
| family=bag;stmt=comma_join;threshold=20000 | j06_comma_join | valid | `751fc42e2c0aff0b` |  |
| family=bag;stmt=cte_filtered;threshold=100000 | j04_cte_filtered | valid | `39c6ba116bae94c2` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=180000 | j04_cte_filtered | valid | `1196a7f18f021d99` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=20 | j04_cte_filtered | valid | `f62fd432fc246d37` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=200 | j04_cte_filtered | valid | `97ef9a27130c97d7` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=2000 | j04_cte_filtered | valid | `a8c0e4875e577d9d` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=20000 | j04_cte_filtered | valid | `532ae0b1ec12b6e8` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_mat;threshold=100000 | j04m_cte_filtered_mat | valid | `8cf4d9a7ceee90eb` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=180000 | j04m_cte_filtered_mat | valid | `fbc473b5e69cbf09` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=20 | j04m_cte_filtered_mat | valid | `4c060efb65471e3a` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=200 | j04m_cte_filtered_mat | valid | `b2423d865acb29f8` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=2000 | j04m_cte_filtered_mat | valid | `6612dcd9334a9df0` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=20000 | j04m_cte_filtered_mat | valid | `125d073f7c84f586` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_notmat;threshold=100000 | j04n_cte_filtered_notmat | valid | `c2e7d79eff631a12` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=180000 | j04n_cte_filtered_notmat | valid | `7642f709884b3470` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=20 | j04n_cte_filtered_notmat | valid | `ffa81c69c186dbb1` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=200 | j04n_cte_filtered_notmat | valid | `2294e10ea085900b` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=2000 | j04n_cte_filtered_notmat | valid | `f68f52eea6fb16a3` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=20000 | j04n_cte_filtered_notmat | valid | `cbc654f86225630e` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=100000 | j05_cte_joined_outer | valid | `10f4009e569caede` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=180000 | j05_cte_joined_outer | valid | `096401580690616f` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20 | j05_cte_joined_outer | valid | `7938c0984573af05` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=200 | j05_cte_joined_outer | valid | `a248d836d59d7cb5` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=2000 | j05_cte_joined_outer | valid | `f7eca020581a23e5` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20000 | j05_cte_joined_outer | valid | `7f416ebb39af6df8` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=2000;position=last | j05_cte_joined_outer | valid | `bc98fed7223e0be9` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_mat;threshold=100000 | j05m_cte_joined_outer_mat | valid | `718d6c094fbe4b04` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | j05m_cte_joined_outer_mat | valid | `128b8c1447bb758a` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | j05m_cte_joined_outer_mat | valid | `2d7890663945ff03` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=200 | j05m_cte_joined_outer_mat | valid | `9a788dcac9b07abe` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | j05m_cte_joined_outer_mat | valid | `1a70f88c28f6ce80` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20000 | j05m_cte_joined_outer_mat | valid | `8bbb703fa614948a` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_notmat;threshold=100000 | j05n_cte_joined_outer_notmat | valid | `229f18f6ebbfb53d` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | j05n_cte_joined_outer_notmat | valid | `2fa2c7a8e62e5dce` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | j05n_cte_joined_outer_notmat | valid | `8f47da145246773f` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=200 | j05n_cte_joined_outer_notmat | valid | `d85f7bf8868e0774` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | j05n_cte_joined_outer_notmat | valid | `2c502df8c990624e` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20000 | j05n_cte_joined_outer_notmat | valid | `e25f3b0c97d03a59` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=filtered_derived;threshold=100000 | j03_filtered_derived | valid | `015aaaa8ce376518` |  |
| family=bag;stmt=filtered_derived;threshold=180000 | j03_filtered_derived | valid | `1e712c06fb2c4ccb` |  |
| family=bag;stmt=filtered_derived;threshold=20 | j03_filtered_derived | valid | `a28a11f79b1b994e` |  |
| family=bag;stmt=filtered_derived;threshold=200 | j03_filtered_derived | valid | `2be20529ba2bff9b` |  |
| family=bag;stmt=filtered_derived;threshold=2000 | j03_filtered_derived | valid | `394584213a52f2b6` |  |
| family=bag;stmt=filtered_derived;threshold=20000 | j03_filtered_derived | valid | `b0c579d4760da41d` |  |
| family=bag;stmt=join_on;threshold=100000 | j02_join_on | valid | `426f7ec64a790925` |  |
| family=bag;stmt=join_on;threshold=180000 | j02_join_on | valid | `0f1963b9b63bc908` |  |
| family=bag;stmt=join_on;threshold=20 | j02_join_on | valid | `6f83d657c174e765` |  |
| family=bag;stmt=join_on;threshold=200 | j02_join_on | valid | `b3d0917943db8dfd` |  |
| family=bag;stmt=join_on;threshold=2000 | j02_join_on | valid | `4ebd5f1e9321fb58` |  |
| family=bag;stmt=join_on;threshold=20000 | j02_join_on | valid | `d6968347e6195bb7` |  |
| family=bag;stmt=join_where;threshold=100000 | j01_join_where | valid | `382e91420ef50bc0` |  |
| family=bag;stmt=join_where;threshold=180000 | j01_join_where | valid | `e126cd1e6d7acea8` |  |
| family=bag;stmt=join_where;threshold=20 | j01_join_where | valid | `eebbab527ebc1219` |  |
| family=bag;stmt=join_where;threshold=200 | j01_join_where | valid | `85efa96a14ec09e2` |  |
| family=bag;stmt=join_where;threshold=2000 | j01_join_where | valid | `9ffbf9b0cefbb566` |  |
| family=bag;stmt=join_where;threshold=20000 | j01_join_where | valid | `76f3119d597c03c7` |  |
| family=reuse;policy=default;refs=16;consumer=1% | r_default_n16_c1 | valid | `c0eff00634f64a4c` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=16;consumer=100% | r_default_n16_c100 | valid | `9979f8fab58a1d21` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=1;consumer=1% | r_default_n1_c1 | valid | `c6faa2b2a434a4ff` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=1;consumer=100% | r_default_n1_c100 | valid | `1334f16b54340058` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=2;consumer=1% | r_default_n2_c1 | valid | `80c694ac69713bcf` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=2;consumer=100% | r_default_n2_c100 | valid | `cb5c42885c210034` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=1% | r_default_n4_c1 | valid | `8d99023cbb8afdc7` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=100% | r_default_n4_c100 | valid | `2496c33e3db1ea40` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=8;consumer=1% | r_default_n8_c1 | valid | `4eadbe18a4648146` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=8;consumer=100% | r_default_n8_c100 | valid | `6f36ce56f2d4b0d5` | policy=default; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=1% | r_mat_n16_c1 | valid | `95e5790ee3291b8e` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=100% | r_mat_n16_c100 | valid | `07f1ca3b642ce096` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=1% | r_mat_n1_c1 | valid | `f4b8330c95f75e83` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=100% | r_mat_n1_c100 | valid | `e0b5d89ff8646f81` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=2;consumer=1% | r_mat_n2_c1 | valid | `872459181a59e5c2` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=2;consumer=100% | r_mat_n2_c100 | valid | `92a6aaa309d0ce5c` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=1% | r_mat_n4_c1 | valid | `a0fd8e4360a67939` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=100% | r_mat_n4_c100 | valid | `7083ae8059b5fe80` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=8;consumer=1% | r_mat_n8_c1 | valid | `30f7e07a2edc465b` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=8;consumer=100% | r_mat_n8_c100 | valid | `a5a60931cb48284a` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=notmat;refs=16;consumer=1% | r_notmat_n16_c1 | valid | `52fd9c5cd35a3a9c` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `edaee8bff2874abd` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `6b0896655055ba83` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `bbbac4d17d7adc05` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=2;consumer=1% | r_notmat_n2_c1 | valid | `eaa463088ee68259` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=2;consumer=100% | r_notmat_n2_c100 | valid | `f5ba60e75a457580` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `acaec9b309b9ca63` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `f9b3f5d54c85f804` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=8;consumer=1% | r_notmat_n8_c1 | valid | `e3680db6113c66ab` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=8;consumer=100% | r_notmat_n8_c100 | valid | `a92b70780ebec35a` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, DIST, DEBUG, FORMAT TEXT` supported on yugabyte
- one host: this run used node budget 2 CPU 3g and client budget 1 CPU 2g, so budget comparisons are separate classes here, not independent hardware
- one client process with one fixed maintenance writer at 10 updates/s; readers are sequential; the reported sample latency includes retry and backoff time
- physical rank-index arm yb-range-rank; an engine whose planner ignores the index is a different physical design, not a faster one
- YugabyteDB DIST/DEBUG storage counters (rows and index rows scanned) are parsed per cell; RPC-level counters are not aggregated
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
