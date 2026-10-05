# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 91 cells, 91 valid, engine `yugabyte`, topology `single`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-v3b-domain-yb-b3-20261005`
- run tag: `run/study-01/select-join-v3b-domain-yb-b3-20261005`
- environment: `host-zenbook-ux5406sa`
- repo commit: `7243fee074bfb16fc5ad8b6afc1a2d7549080935`
- engine version: `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit`
- database engine image: `docker.io/yugabytedb/yugabyte:2025.2.6.0-b111` @ `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`
- benchmark client image: `localhost/charitybench-select:select-join-v3b-domain-yb-b3-20261005` @ `f3f1e10c44cbcb4c7b5bc56ed95a4f071682b6f7d33b5e0a7b688c4271eeef88`
- inputs digest: `824b5579aeac034ee9c3afe73c771ad5f097531d215069a6a29649c67c1e30b7`
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
| yb-single | 934 | 2 |

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
| 20 (0.01%) | 2.1835 | 2.1824 | 2.1083 | 2.2096 | 2.2421 | 1.6458 | 2.3552 | 81.2771 | 1.9394 | 1.8361 |
| 200 (0.1%) | 2.9529 | 3.6614 | 3.7681 | 3.1969 | 3.0646 | 3.2804 | 3.1365 | 86.7706 | 2.8491 | 3.3524 |
| 2000 (1%) | 9.5336 | 9.6440 | 8.3933 | 9.7351 | 9.8200 | 9.8942 | 9.5881 | 82.8729 | 9.3775 | 8.7519 |
| 20000 (10%) | 63.6895 | 63.9069 | 62.9096 | 64.3350 | 65.6807 | 64.9550 | 64.1595 | 88.2650 | 63.4140 | 63.3152 |
| 100000 (50%) | 87.0115 | 88.8777 | 87.7463 | 87.2704 | 90.8677 | 86.0191 | 87.2393 | 98.4019 | 84.9575 | 87.6625 |
| 180000 (90%) | 112.8999 | 107.8579 | 112.0664 | 109.4673 | 114.6585 | 109.6103 | 112.5127 | 104.3876 | 109.0088 | 109.0250 |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | join_on vs join_where | 2.1835 | 2.1824 | 0.9995 | false |
| 20 | filtered_derived vs join_where | 2.1835 | 2.1083 | 0.9656 | false |
| 20 | cte_filtered vs join_where | 2.1835 | 2.2096 | 1.0120 | false |
| 20 | cte_filtered_mat vs join_where | 2.1835 | 2.2421 | 1.0268 | false |
| 20 | cte_filtered_notmat vs join_where | 2.1835 | 1.6458 | 0.7538 | false |
| 20 | cte_joined_outer vs join_where | 2.1835 | 2.3552 | 1.0786 | false |
| 20 | cte_joined_outer_mat vs join_where | 2.1835 | 81.2771 | 37.2235 | false |
| 20 | cte_joined_outer_notmat vs join_where | 2.1835 | 1.9394 | 0.8882 | false |
| 20 | comma_join vs join_where | 2.1835 | 1.8361 | 0.8409 | false |
| 200 | join_on vs join_where | 2.9529 | 3.6614 | 1.2399 | false |
| 200 | filtered_derived vs join_where | 2.9529 | 3.7681 | 1.2761 | false |
| 200 | cte_filtered vs join_where | 2.9529 | 3.1969 | 1.0826 | false |
| 200 | cte_filtered_mat vs join_where | 2.9529 | 3.0646 | 1.0378 | false |
| 200 | cte_filtered_notmat vs join_where | 2.9529 | 3.2804 | 1.1109 | false |
| 200 | cte_joined_outer vs join_where | 2.9529 | 3.1365 | 1.0622 | false |
| 200 | cte_joined_outer_mat vs join_where | 2.9529 | 86.7706 | 29.3845 | false |
| 200 | cte_joined_outer_notmat vs join_where | 2.9529 | 2.8491 | 0.9648 | false |
| 200 | comma_join vs join_where | 2.9529 | 3.3524 | 1.1353 | false |
| 2000 | join_on vs join_where | 9.5336 | 9.6440 | 1.0116 | false |
| 2000 | filtered_derived vs join_where | 9.5336 | 8.3933 | 0.8804 | false |
| 2000 | cte_filtered vs join_where | 9.5336 | 9.7351 | 1.0211 | false |
| 2000 | cte_filtered_mat vs join_where | 9.5336 | 9.8200 | 1.0300 | false |
| 2000 | cte_filtered_notmat vs join_where | 9.5336 | 9.8942 | 1.0378 | false |
| 2000 | cte_joined_outer vs join_where | 9.5336 | 9.5881 | 1.0057 | false |
| 2000 | cte_joined_outer_mat vs join_where | 9.5336 | 82.8729 | 8.6927 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 9.5336 | 9.3775 | 0.9836 | false |
| 2000 | comma_join vs join_where | 9.5336 | 8.7519 | 0.9180 | false |
| 20000 | join_on vs join_where | 63.6895 | 63.9069 | 1.0034 | false |
| 20000 | filtered_derived vs join_where | 63.6895 | 62.9096 | 0.9878 | false |
| 20000 | cte_filtered vs join_where | 63.6895 | 64.3350 | 1.0101 | false |
| 20000 | cte_filtered_mat vs join_where | 63.6895 | 65.6807 | 1.0313 | false |
| 20000 | cte_filtered_notmat vs join_where | 63.6895 | 64.9550 | 1.0199 | false |
| 20000 | cte_joined_outer vs join_where | 63.6895 | 64.1595 | 1.0074 | false |
| 20000 | cte_joined_outer_mat vs join_where | 63.6895 | 88.2650 | 1.3859 | false |
| 20000 | cte_joined_outer_notmat vs join_where | 63.6895 | 63.4140 | 0.9957 | false |
| 20000 | comma_join vs join_where | 63.6895 | 63.3152 | 0.9941 | false |
| 100000 | join_on vs join_where | 87.0115 | 88.8777 | 1.0214 | false |
| 100000 | filtered_derived vs join_where | 87.0115 | 87.7463 | 1.0084 | false |
| 100000 | cte_filtered vs join_where | 87.0115 | 87.2704 | 1.0030 | false |
| 100000 | cte_filtered_mat vs join_where | 87.0115 | 90.8677 | 1.0443 | false |
| 100000 | cte_filtered_notmat vs join_where | 87.0115 | 86.0191 | 0.9886 | false |
| 100000 | cte_joined_outer vs join_where | 87.0115 | 87.2393 | 1.0026 | false |
| 100000 | cte_joined_outer_mat vs join_where | 87.0115 | 98.4019 | 1.1309 | false |
| 100000 | cte_joined_outer_notmat vs join_where | 87.0115 | 84.9575 | 0.9764 | false |
| 100000 | comma_join vs join_where | 87.0115 | 87.6625 | 1.0075 | false |
| 180000 | join_on vs join_where | 112.8999 | 107.8579 | 0.9553 | false |
| 180000 | filtered_derived vs join_where | 112.8999 | 112.0664 | 0.9926 | false |
| 180000 | cte_filtered vs join_where | 112.8999 | 109.4673 | 0.9696 | false |
| 180000 | cte_filtered_mat vs join_where | 112.8999 | 114.6585 | 1.0156 | false |
| 180000 | cte_filtered_notmat vs join_where | 112.8999 | 109.6103 | 0.9709 | false |
| 180000 | cte_joined_outer vs join_where | 112.8999 | 112.5127 | 0.9966 | false |
| 180000 | cte_joined_outer_mat vs join_where | 112.8999 | 104.3876 | 0.9246 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 112.8999 | 109.0088 | 0.9655 | false |
| 180000 | comma_join vs join_where | 112.8999 | 109.0250 | 0.9657 | false |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_filtered_mat vs cte_filtered | 2.2096 | 2.2421 | 1.0147 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_filtered_notmat vs cte_filtered | 2.2096 | 1.6458 | 0.7448 | false | policy=not_materialized; observed_cte_scan=false |
| 20 | cte_joined_outer_mat vs cte_joined_outer | 2.3552 | 81.2771 | 34.5097 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 2.3552 | 1.9394 | 0.8234 | false | policy=not_materialized; observed_cte_scan=false |
| 200 | cte_filtered_mat vs cte_filtered | 3.1969 | 3.0646 | 0.9586 | false | policy=materialized; observed_cte_scan=true |
| 200 | cte_filtered_notmat vs cte_filtered | 3.1969 | 3.2804 | 1.0261 | false | policy=not_materialized; observed_cte_scan=false |
| 200 | cte_joined_outer_mat vs cte_joined_outer | 3.1365 | 86.7706 | 27.6646 | false | policy=materialized; observed_cte_scan=true |
| 200 | cte_joined_outer_notmat vs cte_joined_outer | 3.1365 | 2.8491 | 0.9084 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_filtered_mat vs cte_filtered | 9.7351 | 9.8200 | 1.0087 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_filtered_notmat vs cte_filtered | 9.7351 | 9.8942 | 1.0163 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 9.5881 | 82.8729 | 8.6433 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 9.5881 | 9.3775 | 0.9780 | false | policy=not_materialized; observed_cte_scan=false |
| 20000 | cte_filtered_mat vs cte_filtered | 64.3350 | 65.6807 | 1.0209 | false | policy=materialized; observed_cte_scan=true |
| 20000 | cte_filtered_notmat vs cte_filtered | 64.3350 | 64.9550 | 1.0096 | false | policy=not_materialized; observed_cte_scan=false |
| 20000 | cte_joined_outer_mat vs cte_joined_outer | 64.1595 | 88.2650 | 1.3757 | false | policy=materialized; observed_cte_scan=true |
| 20000 | cte_joined_outer_notmat vs cte_joined_outer | 64.1595 | 63.4140 | 0.9884 | false | policy=not_materialized; observed_cte_scan=false |
| 100000 | cte_filtered_mat vs cte_filtered | 87.2704 | 90.8677 | 1.0412 | false | policy=materialized; observed_cte_scan=true |
| 100000 | cte_filtered_notmat vs cte_filtered | 87.2704 | 86.0191 | 0.9857 | false | policy=not_materialized; observed_cte_scan=false |
| 100000 | cte_joined_outer_mat vs cte_joined_outer | 87.2393 | 98.4019 | 1.1280 | false | policy=materialized; observed_cte_scan=true |
| 100000 | cte_joined_outer_notmat vs cte_joined_outer | 87.2393 | 84.9575 | 0.9738 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_filtered_mat vs cte_filtered | 109.4673 | 114.6585 | 1.0474 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_filtered_notmat vs cte_filtered | 109.4673 | 109.6103 | 1.0013 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 112.5127 | 104.3876 | 0.9278 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 112.5127 | 109.0088 | 0.9689 | false | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 6.8390 | 1 | `c0f1461fc9af5c4a` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 8.4109 | 1 | `faf183ad391d43cf` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | 7.9880 | 2 | `ce7dd3ec58d1a4d4` | policy=default; observed_cte_scan=true | valid |
| default | 2 | 100% | 8.8538 | 2 | `3599c771c6b98e2d` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 1% | 8.5764 | 4 | `e2bceaddcb1249a4` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 8.2380 | 4 | `a0d8f9cea31e6a45` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | 8.2337 | 8 | `5a1468191cd70d6a` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 100% | 9.1456 | 8 | `eb1aeb3b45bacfa5` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 1% | 8.8372 | 16 | `95f12267e112124e` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 10.4310 | 16 | `9ed3a1dfbb8cc07b` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 8.0568 | 1 | `29a1a3e4e8eb501d` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 7.9181 | 1 | `929942fb68f5853c` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | 8.1190 | 2 | `1c39bdd6ea09ed3c` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 100% | 8.7531 | 2 | `12e954e597d6f056` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 1% | 9.0837 | 4 | `4b983f7b792bce7b` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 8.5006 | 4 | `f8e6091d1978b19b` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | 8.4037 | 8 | `3ce9832f813c60e3` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 100% | 9.3121 | 8 | `959fdd05644fcf63` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 1% | 9.2632 | 16 | `be7c8eea5de9a89a` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 10.4876 | 16 | `8fbbff9121da9ed6` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 7.6069 | 1 | `31287b2dd6c5bf42` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 8.1437 | 1 | `906939ec9aa0b481` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | 11.7632 | 2 | `e80965534823d481` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 100% | 14.7076 | 2 | `3b720d4f0253724f` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 1% | 23.0552 | 4 | `84f5132085e70848` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 28.8889 | 4 | `af5fe6e27ad7c0ff` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | 57.8272 | 8 | `af2002fdd774eb7c` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 100% | 56.2829 | 8 | `dd4e72d88884197e` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 1% | 83.9057 | 16 | `928ecdd499dc0738` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 99.4849 | 16 | `12baf4cfd000e6a2` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 1.1781 | false |
| notmat | 1 | 1% | 1.1123 | false |
| mat | 1 | 100% | 0.9414 | false |
| notmat | 1 | 100% | 0.9682 | false |
| mat | 2 | 1% | 1.0164 | false |
| notmat | 2 | 1% | 1.4726 | false |
| mat | 2 | 100% | 0.9886 | false |
| notmat | 2 | 100% | 1.6612 | false |
| mat | 4 | 1% | 1.0591 | false |
| notmat | 4 | 1% | 2.6882 | false |
| mat | 4 | 100% | 1.0319 | false |
| notmat | 4 | 100% | 3.5068 | false |
| mat | 8 | 1% | 1.0207 | false |
| notmat | 8 | 1% | 7.0233 | false |
| mat | 8 | 100% | 1.0182 | false |
| notmat | 8 | 100% | 6.1541 | false |
| mat | 16 | 1% | 1.0482 | false |
| notmat | 16 | 1% | 9.4946 | false |
| mat | 16 | 100% | 1.0054 | false |
| notmat | 16 | 100% | 9.5374 | false |

## Repeated-design position control

| first median ms | last median ms | ratio (last/first) | same plan |
|---|---|---|---|
| 8.233669 | 8.253741 | 1.0024 | false |

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

- client: nr_periods 610, nr_throttled 242, throttled_usec 10359570 | quota 100000 100000
, memory 2147483648

- yb-single: nr_periods 672, nr_throttled 68, throttled_usec 2696404 | quota 200000 100000, memory 3221225472

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
| family=bag;stmt=cte_filtered;threshold=100000 | valid | false | true | 202050 | 0 | 0/0 | 0 | 1 |
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
| family=reuse;policy=default;refs=8;consumer=1%;position=last | valid | true | false | 4000 | 2000 | 0/0 | 1983 | 0 |
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
| family=bag;stmt=comma_join;threshold=100000 | valid | 85.3437 | 87.6625 | 90.6970 | 5 | 11.4 |
| family=bag;stmt=comma_join;threshold=180000 | valid | 104.8617 | 109.0250 | 113.9179 | 5 | 9.2 |
| family=bag;stmt=comma_join;threshold=20 | valid | 1.6130 | 1.8361 | 2.0064 | 5 | 544.6 |
| family=bag;stmt=comma_join;threshold=200 | valid | 3.0725 | 3.3524 | 3.8992 | 5 | 298.3 |
| family=bag;stmt=comma_join;threshold=2000 | valid | 8.5301 | 8.7519 | 9.5407 | 5 | 114.3 |
| family=bag;stmt=comma_join;threshold=20000 | valid | 61.5106 | 63.3152 | 65.4995 | 5 | 15.8 |
| family=bag;stmt=cte_filtered;threshold=100000 | valid | 85.4163 | 87.2704 | 182.3779 | 5 | 11.5 |
| family=bag;stmt=cte_filtered;threshold=180000 | valid | 108.2554 | 109.4673 | 111.7863 | 5 | 9.1 |
| family=bag;stmt=cte_filtered;threshold=20 | valid | 2.1026 | 2.2096 | 2.5276 | 5 | 452.6 |
| family=bag;stmt=cte_filtered;threshold=200 | valid | 2.9373 | 3.1969 | 3.2981 | 5 | 312.8 |
| family=bag;stmt=cte_filtered;threshold=2000 | valid | 8.7321 | 9.7351 | 10.2777 | 5 | 102.7 |
| family=bag;stmt=cte_filtered;threshold=20000 | valid | 62.0428 | 64.3350 | 68.5093 | 5 | 15.5 |
| family=bag;stmt=cte_filtered_mat;threshold=100000 | valid | 86.8724 | 90.8677 | 95.4986 | 5 | 11.0 |
| family=bag;stmt=cte_filtered_mat;threshold=180000 | valid | 108.0272 | 114.6585 | 117.8341 | 5 | 8.7 |
| family=bag;stmt=cte_filtered_mat;threshold=20 | valid | 2.1307 | 2.2421 | 2.6652 | 5 | 446.0 |
| family=bag;stmt=cte_filtered_mat;threshold=200 | valid | 2.8766 | 3.0646 | 3.2079 | 5 | 326.3 |
| family=bag;stmt=cte_filtered_mat;threshold=2000 | valid | 9.2264 | 9.8200 | 10.5024 | 5 | 101.8 |
| family=bag;stmt=cte_filtered_mat;threshold=20000 | valid | 65.3187 | 65.6807 | 67.4830 | 5 | 15.2 |
| family=bag;stmt=cte_filtered_notmat;threshold=100000 | valid | 81.0487 | 86.0191 | 89.3018 | 5 | 11.6 |
| family=bag;stmt=cte_filtered_notmat;threshold=180000 | valid | 103.8798 | 109.6103 | 112.6073 | 5 | 9.1 |
| family=bag;stmt=cte_filtered_notmat;threshold=20 | valid | 1.4371 | 1.6458 | 1.7668 | 5 | 607.6 |
| family=bag;stmt=cte_filtered_notmat;threshold=200 | valid | 2.9983 | 3.2804 | 3.4715 | 5 | 304.8 |
| family=bag;stmt=cte_filtered_notmat;threshold=2000 | valid | 9.4794 | 9.8942 | 11.5338 | 5 | 101.1 |
| family=bag;stmt=cte_filtered_notmat;threshold=20000 | valid | 61.6137 | 64.9550 | 71.4853 | 5 | 15.4 |
| family=bag;stmt=cte_joined_outer;threshold=100000 | valid | 85.0141 | 87.2393 | 89.7892 | 5 | 11.5 |
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 108.6206 | 112.5127 | 121.3721 | 5 | 8.9 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 2.2038 | 2.3552 | 2.6119 | 5 | 424.6 |
| family=bag;stmt=cte_joined_outer;threshold=200 | valid | 2.9853 | 3.1365 | 3.7108 | 5 | 318.8 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 9.1967 | 9.5881 | 10.5228 | 5 | 104.3 |
| family=bag;stmt=cte_joined_outer;threshold=20000 | valid | 62.1730 | 64.1595 | 67.1061 | 5 | 15.6 |
| family=bag;stmt=cte_joined_outer_mat;threshold=100000 | valid | 93.4801 | 98.4019 | 111.6866 | 5 | 10.2 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 97.7772 | 104.3876 | 105.4300 | 5 | 9.6 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 79.0060 | 81.2771 | 85.2179 | 5 | 12.3 |
| family=bag;stmt=cte_joined_outer_mat;threshold=200 | valid | 78.0759 | 86.7706 | 90.5733 | 5 | 11.5 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 79.9829 | 82.8729 | 106.7827 | 5 | 12.1 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20000 | valid | 85.5671 | 88.2650 | 106.9186 | 5 | 11.3 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=100000 | valid | 79.9871 | 84.9575 | 90.3854 | 5 | 11.8 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 105.7682 | 109.0088 | 115.0667 | 5 | 9.2 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 1.7392 | 1.9394 | 2.0805 | 5 | 515.6 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=200 | valid | 2.6763 | 2.8491 | 3.2002 | 5 | 351.0 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 9.2659 | 9.3775 | 10.3138 | 5 | 106.6 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20000 | valid | 60.5695 | 63.4140 | 67.0490 | 5 | 15.8 |
| family=bag;stmt=filtered_derived;threshold=100000 | valid | 82.5190 | 87.7463 | 91.8426 | 5 | 11.4 |
| family=bag;stmt=filtered_derived;threshold=180000 | valid | 106.8163 | 112.0664 | 143.3492 | 5 | 8.9 |
| family=bag;stmt=filtered_derived;threshold=20 | valid | 1.6336 | 2.1083 | 2.2434 | 5 | 474.3 |
| family=bag;stmt=filtered_derived;threshold=200 | valid | 3.3388 | 3.7681 | 4.4376 | 5 | 265.4 |
| family=bag;stmt=filtered_derived;threshold=2000 | valid | 8.0060 | 8.3933 | 8.9741 | 5 | 119.1 |
| family=bag;stmt=filtered_derived;threshold=20000 | valid | 61.2612 | 62.9096 | 64.3635 | 5 | 15.9 |
| family=bag;stmt=join_on;threshold=100000 | valid | 86.7476 | 88.8777 | 90.6817 | 5 | 11.3 |
| family=bag;stmt=join_on;threshold=180000 | valid | 104.1247 | 107.8579 | 112.5435 | 5 | 9.3 |
| family=bag;stmt=join_on;threshold=20 | valid | 1.9190 | 2.1824 | 2.4195 | 5 | 458.2 |
| family=bag;stmt=join_on;threshold=200 | valid | 3.4598 | 3.6614 | 3.9020 | 5 | 273.1 |
| family=bag;stmt=join_on;threshold=2000 | valid | 9.1390 | 9.6440 | 10.4166 | 5 | 103.7 |
| family=bag;stmt=join_on;threshold=20000 | valid | 61.5294 | 63.9069 | 70.2009 | 5 | 15.6 |
| family=bag;stmt=join_where;threshold=100000 | valid | 84.2123 | 87.0115 | 89.1125 | 5 | 11.5 |
| family=bag;stmt=join_where;threshold=180000 | valid | 106.8307 | 112.8999 | 116.9694 | 5 | 8.9 |
| family=bag;stmt=join_where;threshold=20 | valid | 2.0550 | 2.1835 | 2.6636 | 5 | 458.0 |
| family=bag;stmt=join_where;threshold=200 | valid | 2.8694 | 2.9529 | 3.2303 | 5 | 338.6 |
| family=bag;stmt=join_where;threshold=2000 | valid | 8.9714 | 9.5336 | 10.1533 | 5 | 104.9 |
| family=bag;stmt=join_where;threshold=20000 | valid | 62.4555 | 63.6895 | 66.9903 | 5 | 15.7 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 8.6374 | 8.8372 | 8.8804 | 5 | 113.2 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 9.2500 | 10.4310 | 10.9495 | 5 | 95.9 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 6.2537 | 6.8390 | 8.0046 | 5 | 146.2 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 7.9061 | 8.4109 | 8.7455 | 5 | 118.9 |
| family=reuse;policy=default;refs=2;consumer=1% | valid | 7.8144 | 7.9880 | 8.3307 | 5 | 125.2 |
| family=reuse;policy=default;refs=2;consumer=100% | valid | 8.1438 | 8.8538 | 9.3273 | 5 | 112.9 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 8.4032 | 8.5764 | 9.2899 | 5 | 116.6 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 8.1567 | 8.2380 | 8.8336 | 5 | 121.4 |
| family=reuse;policy=default;refs=8;consumer=1% | valid | 7.8892 | 8.2337 | 9.3087 | 5 | 121.5 |
| family=reuse;policy=default;refs=8;consumer=1%;position=last | valid | 7.8075 | 8.2537 | 8.6982 | 5 | 121.2 |
| family=reuse;policy=default;refs=8;consumer=100% | valid | 9.0671 | 9.1456 | 11.4284 | 5 | 109.3 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 9.0663 | 9.2632 | 9.3471 | 5 | 108.0 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 9.8983 | 10.4876 | 11.9328 | 5 | 95.4 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 7.6363 | 8.0568 | 8.5598 | 5 | 124.1 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 7.8090 | 7.9181 | 9.1881 | 5 | 126.3 |
| family=reuse;policy=mat;refs=2;consumer=1% | valid | 7.8625 | 8.1190 | 8.7421 | 5 | 123.2 |
| family=reuse;policy=mat;refs=2;consumer=100% | valid | 8.0461 | 8.7531 | 9.7351 | 5 | 114.2 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 8.5538 | 9.0837 | 9.4911 | 5 | 110.1 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 7.8893 | 8.5006 | 9.1138 | 5 | 117.6 |
| family=reuse;policy=mat;refs=8;consumer=1% | valid | 8.1711 | 8.4037 | 24.1696 | 5 | 119.0 |
| family=reuse;policy=mat;refs=8;consumer=100% | valid | 8.9349 | 9.3121 | 10.0831 | 5 | 107.4 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 78.0733 | 83.9057 | 99.3814 | 5 | 11.9 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 82.8487 | 99.4849 | 116.5034 | 5 | 10.1 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 6.7467 | 7.6069 | 8.6807 | 5 | 131.5 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 7.1937 | 8.1437 | 9.1239 | 5 | 122.8 |
| family=reuse;policy=notmat;refs=2;consumer=1% | valid | 11.4050 | 11.7632 | 14.2285 | 5 | 85.0 |
| family=reuse;policy=notmat;refs=2;consumer=100% | valid | 14.6618 | 14.7076 | 15.6753 | 5 | 68.0 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 21.8414 | 23.0552 | 23.8013 | 5 | 43.4 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 27.8149 | 28.8889 | 30.5304 | 5 | 34.6 |
| family=reuse;policy=notmat;refs=8;consumer=1% | valid | 54.4225 | 57.8272 | 61.3711 | 5 | 17.3 |
| family=reuse;policy=notmat;refs=8;consumer=100% | valid | 55.5281 | 56.2829 | 59.4610 | 5 | 17.8 |

## Matched blended workload

- phase: `blended`; target offered rate: 10 updates/s
- offered 527, committed 527, failed 0 over 52.630 s
- offered rate 10.013/s, completed rate 10.013/s
- writer operations per endpoint: [527]
- reader operations overlapping a committed write: 219 of 638
- reader retries of transient conflicts: 1
- per key: 527 keys, min 1, max 1 per key
- ledger: 527 keys reconciled, 0 mismatch
- update latency: avg 3.070 ms, min 1.418, max 143.993 over 527

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42638838`, 5 interleaved rounds
  1. family=reuse;policy=default;refs=8;consumer=1%
  2. family=bag;stmt=join_where;threshold=200
  3. family=reuse;policy=notmat;refs=2;consumer=100%
  4. family=bag;stmt=join_where;threshold=100000
  5. family=reuse;policy=notmat;refs=8;consumer=100%
  6. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  7. family=reuse;policy=mat;refs=4;consumer=1%
  8. family=bag;stmt=cte_filtered_mat;threshold=200
  9. family=reuse;policy=mat;refs=8;consumer=1%
  10. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  11. family=bag;stmt=cte_joined_outer_notmat;threshold=100000
  12. family=bag;stmt=cte_filtered_notmat;threshold=100000
  13. family=bag;stmt=filtered_derived;threshold=100000
  14. family=bag;stmt=cte_joined_outer;threshold=2000
  15. family=bag;stmt=cte_filtered;threshold=180000
  16. family=bag;stmt=join_on;threshold=20
  17. family=reuse;policy=notmat;refs=2;consumer=1%
  18. family=reuse;policy=notmat;refs=16;consumer=1%
  19. family=bag;stmt=filtered_derived;threshold=20
  20. family=reuse;policy=mat;refs=2;consumer=1%
  21. family=reuse;policy=mat;refs=4;consumer=100%
  22. family=reuse;policy=notmat;refs=4;consumer=100%
  23. family=bag;stmt=comma_join;threshold=20
  24. family=bag;stmt=cte_filtered_notmat;threshold=20
  25. family=reuse;policy=default;refs=1;consumer=100%
  26. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  27. family=bag;stmt=cte_joined_outer_notmat;threshold=200
  28. family=reuse;policy=default;refs=1;consumer=1%
  29. family=bag;stmt=cte_filtered_notmat;threshold=20000
  30. family=bag;stmt=comma_join;threshold=20000
  31. family=bag;stmt=cte_joined_outer;threshold=180000
  32. family=bag;stmt=cte_filtered_notmat;threshold=2000
  33. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  34. family=reuse;policy=default;refs=2;consumer=1%
  35. family=bag;stmt=cte_filtered;threshold=200
  36. family=reuse;policy=default;refs=4;consumer=100%
  37. family=bag;stmt=cte_filtered_mat;threshold=100000
  38. family=reuse;policy=mat;refs=8;consumer=100%
  39. family=bag;stmt=cte_filtered_mat;threshold=180000
  40. family=reuse;policy=mat;refs=16;consumer=1%
  41. family=bag;stmt=cte_joined_outer_notmat;threshold=20000
  42. family=bag;stmt=comma_join;threshold=100000
  43. family=reuse;policy=mat;refs=16;consumer=100%
  44. family=reuse;policy=mat;refs=1;consumer=1%
  45. family=reuse;policy=notmat;refs=8;consumer=1%
  46. family=reuse;policy=default;refs=8;consumer=100%
  47. family=bag;stmt=comma_join;threshold=200
  48. family=bag;stmt=filtered_derived;threshold=180000
  49. family=bag;stmt=join_where;threshold=20
  50. family=bag;stmt=cte_filtered_notmat;threshold=200
  51. family=bag;stmt=cte_filtered;threshold=2000
  52. family=bag;stmt=comma_join;threshold=2000
  53. family=reuse;policy=mat;refs=2;consumer=100%
  54. family=reuse;policy=notmat;refs=1;consumer=100%
  55. family=bag;stmt=cte_joined_outer_mat;threshold=20000
  56. family=bag;stmt=join_on;threshold=2000
  57. family=reuse;policy=default;refs=16;consumer=100%
  58. family=bag;stmt=join_where;threshold=2000
  59. family=reuse;policy=mat;refs=1;consumer=100%
  60. family=bag;stmt=cte_joined_outer_mat;threshold=200
  61. family=bag;stmt=cte_filtered;threshold=100000
  62. family=bag;stmt=join_on;threshold=200
  63. family=bag;stmt=comma_join;threshold=180000
  64. family=bag;stmt=cte_joined_outer;threshold=20
  65. family=bag;stmt=cte_joined_outer;threshold=200
  66. family=bag;stmt=join_on;threshold=180000
  67. family=reuse;policy=notmat;refs=16;consumer=100%
  68. family=bag;stmt=join_on;threshold=20000
  69. family=bag;stmt=cte_joined_outer_mat;threshold=20
  70. family=bag;stmt=cte_joined_outer_mat;threshold=100000
  71. family=reuse;policy=notmat;refs=1;consumer=1%
  72. family=bag;stmt=join_on;threshold=100000
  73. family=bag;stmt=cte_filtered;threshold=20000
  74. family=reuse;policy=default;refs=2;consumer=100%
  75. family=reuse;policy=notmat;refs=4;consumer=1%
  76. family=bag;stmt=join_where;threshold=180000
  77. family=bag;stmt=cte_filtered_mat;threshold=2000
  78. family=bag;stmt=cte_joined_outer;threshold=100000
  79. family=bag;stmt=cte_filtered_notmat;threshold=180000
  80. family=bag;stmt=cte_filtered_mat;threshold=20
  81. family=reuse;policy=default;refs=16;consumer=1%
  82. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  83. family=bag;stmt=cte_filtered_mat;threshold=20000
  84. family=bag;stmt=filtered_derived;threshold=200
  85. family=bag;stmt=filtered_derived;threshold=20000
  86. family=bag;stmt=cte_joined_outer;threshold=20000
  87. family=bag;stmt=cte_filtered;threshold=20
  88. family=bag;stmt=join_where;threshold=20000
  89. family=reuse;policy=default;refs=4;consumer=1%
  90. family=bag;stmt=filtered_derived;threshold=2000
  91. family=reuse;policy=default;refs=8;consumer=1%;position=last

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
| family=bag;stmt=comma_join;threshold=100000 | j06_comma_join | valid | `206ee3d46391ee63` |  |
| family=bag;stmt=comma_join;threshold=180000 | j06_comma_join | valid | `f5f8f9dcc92b36a7` |  |
| family=bag;stmt=comma_join;threshold=20 | j06_comma_join | valid | `9422f965c25ec632` |  |
| family=bag;stmt=comma_join;threshold=200 | j06_comma_join | valid | `c34e55ee8c05149c` |  |
| family=bag;stmt=comma_join;threshold=2000 | j06_comma_join | valid | `a06cb2b11c9d27c8` |  |
| family=bag;stmt=comma_join;threshold=20000 | j06_comma_join | valid | `5ab0211aae199a2c` |  |
| family=bag;stmt=cte_filtered;threshold=100000 | j04_cte_filtered | valid | `365c0e258ee5270e` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=180000 | j04_cte_filtered | valid | `d5ea98f21c7ce63b` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=20 | j04_cte_filtered | valid | `f00187fdf2a547a4` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=200 | j04_cte_filtered | valid | `5ccdcad7412d22de` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=2000 | j04_cte_filtered | valid | `0649af41d90c6246` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=20000 | j04_cte_filtered | valid | `bcffe2d7eb127675` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_mat;threshold=100000 | j04m_cte_filtered_mat | valid | `89111a7dc4ec690b` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=180000 | j04m_cte_filtered_mat | valid | `574730ef30392230` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=20 | j04m_cte_filtered_mat | valid | `44d085fd86417566` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=200 | j04m_cte_filtered_mat | valid | `aaab078ac6f9c0c7` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=2000 | j04m_cte_filtered_mat | valid | `4936977499d21829` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=20000 | j04m_cte_filtered_mat | valid | `3806a3de72eed080` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_notmat;threshold=100000 | j04n_cte_filtered_notmat | valid | `ca1f8a7751e17cb4` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=180000 | j04n_cte_filtered_notmat | valid | `4d942c004def8ac9` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=20 | j04n_cte_filtered_notmat | valid | `f99fd3934d11079e` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=200 | j04n_cte_filtered_notmat | valid | `0270b7bf72d8aa1e` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=2000 | j04n_cte_filtered_notmat | valid | `d7f4c17ecb759192` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=20000 | j04n_cte_filtered_notmat | valid | `e51877174d2a0dc8` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=100000 | j05_cte_joined_outer | valid | `3e9bfc9b5236eaba` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=180000 | j05_cte_joined_outer | valid | `f2bba35cbdb79a3f` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20 | j05_cte_joined_outer | valid | `45ae60f13e1e35a0` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=200 | j05_cte_joined_outer | valid | `0332c622ed25e66f` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=2000 | j05_cte_joined_outer | valid | `186412bd3d1c2615` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20000 | j05_cte_joined_outer | valid | `e2b234b6214bbd25` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_mat;threshold=100000 | j05m_cte_joined_outer_mat | valid | `38997c9ca33993f6` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | j05m_cte_joined_outer_mat | valid | `d6b629ab300b0323` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | j05m_cte_joined_outer_mat | valid | `596c637892bbbc8b` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=200 | j05m_cte_joined_outer_mat | valid | `1ee6bd50d4843129` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | j05m_cte_joined_outer_mat | valid | `00ea332eb0d110db` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20000 | j05m_cte_joined_outer_mat | valid | `81b884b61ea2dce1` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_notmat;threshold=100000 | j05n_cte_joined_outer_notmat | valid | `9535752cef78f0a4` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | j05n_cte_joined_outer_notmat | valid | `c0d6c4159b3ddb0b` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | j05n_cte_joined_outer_notmat | valid | `71aa87a6510d1e5d` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=200 | j05n_cte_joined_outer_notmat | valid | `bc6db221dc0c611f` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | j05n_cte_joined_outer_notmat | valid | `5d5789d2c47351e3` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20000 | j05n_cte_joined_outer_notmat | valid | `7eb5d1026f6b9c62` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=filtered_derived;threshold=100000 | j03_filtered_derived | valid | `066be5ddb73e4cac` |  |
| family=bag;stmt=filtered_derived;threshold=180000 | j03_filtered_derived | valid | `8c4407829bdb5060` |  |
| family=bag;stmt=filtered_derived;threshold=20 | j03_filtered_derived | valid | `85ac3dc3cbd68bcd` |  |
| family=bag;stmt=filtered_derived;threshold=200 | j03_filtered_derived | valid | `64d7b87c4869b958` |  |
| family=bag;stmt=filtered_derived;threshold=2000 | j03_filtered_derived | valid | `c719359d8324c0f2` |  |
| family=bag;stmt=filtered_derived;threshold=20000 | j03_filtered_derived | valid | `56de841a51da6154` |  |
| family=bag;stmt=join_on;threshold=100000 | j02_join_on | valid | `924faa9e0cc5d90c` |  |
| family=bag;stmt=join_on;threshold=180000 | j02_join_on | valid | `8366702c145106e0` |  |
| family=bag;stmt=join_on;threshold=20 | j02_join_on | valid | `b1cf131473784268` |  |
| family=bag;stmt=join_on;threshold=200 | j02_join_on | valid | `4575ef838a1111ff` |  |
| family=bag;stmt=join_on;threshold=2000 | j02_join_on | valid | `e9faf8b070973ed0` |  |
| family=bag;stmt=join_on;threshold=20000 | j02_join_on | valid | `e01c12e98440b560` |  |
| family=bag;stmt=join_where;threshold=100000 | j01_join_where | valid | `506c795bc57f471f` |  |
| family=bag;stmt=join_where;threshold=180000 | j01_join_where | valid | `39138ac18b3042a3` |  |
| family=bag;stmt=join_where;threshold=20 | j01_join_where | valid | `6befa15877791f31` |  |
| family=bag;stmt=join_where;threshold=200 | j01_join_where | valid | `9d39fa09ac4cedf8` |  |
| family=bag;stmt=join_where;threshold=2000 | j01_join_where | valid | `fb44dd60352b7bfa` |  |
| family=bag;stmt=join_where;threshold=20000 | j01_join_where | valid | `7649594cc1f5890a` |  |
| family=reuse;policy=default;refs=16;consumer=1% | r_default_n16_c1 | valid | `95f12267e112124e` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=16;consumer=100% | r_default_n16_c100 | valid | `9ed3a1dfbb8cc07b` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=1;consumer=1% | r_default_n1_c1 | valid | `c0f1461fc9af5c4a` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=1;consumer=100% | r_default_n1_c100 | valid | `faf183ad391d43cf` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=2;consumer=1% | r_default_n2_c1 | valid | `ce7dd3ec58d1a4d4` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=2;consumer=100% | r_default_n2_c100 | valid | `3599c771c6b98e2d` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=1% | r_default_n4_c1 | valid | `e2bceaddcb1249a4` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=100% | r_default_n4_c100 | valid | `a0d8f9cea31e6a45` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=8;consumer=1% | r_default_n8_c1 | valid | `5a1468191cd70d6a` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=8;consumer=1%;position=last | r_default_n8_c1 | valid | `e71100f6d7bb8aec` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=8;consumer=100% | r_default_n8_c100 | valid | `eb1aeb3b45bacfa5` | policy=default; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=1% | r_mat_n16_c1 | valid | `be7c8eea5de9a89a` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=100% | r_mat_n16_c100 | valid | `8fbbff9121da9ed6` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=1% | r_mat_n1_c1 | valid | `29a1a3e4e8eb501d` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=100% | r_mat_n1_c100 | valid | `929942fb68f5853c` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=2;consumer=1% | r_mat_n2_c1 | valid | `1c39bdd6ea09ed3c` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=2;consumer=100% | r_mat_n2_c100 | valid | `12e954e597d6f056` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=1% | r_mat_n4_c1 | valid | `4b983f7b792bce7b` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=100% | r_mat_n4_c100 | valid | `f8e6091d1978b19b` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=8;consumer=1% | r_mat_n8_c1 | valid | `3ce9832f813c60e3` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=8;consumer=100% | r_mat_n8_c100 | valid | `959fdd05644fcf63` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=notmat;refs=16;consumer=1% | r_notmat_n16_c1 | valid | `928ecdd499dc0738` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `12baf4cfd000e6a2` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `31287b2dd6c5bf42` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `906939ec9aa0b481` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=2;consumer=1% | r_notmat_n2_c1 | valid | `e80965534823d481` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=2;consumer=100% | r_notmat_n2_c100 | valid | `3b720d4f0253724f` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `84f5132085e70848` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `af5fe6e27ad7c0ff` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=8;consumer=1% | r_notmat_n8_c1 | valid | `af2002fdd774eb7c` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=8;consumer=100% | r_notmat_n8_c100 | valid | `dd4e72d88884197e` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, DIST, DEBUG, FORMAT TEXT` supported on yugabyte
- one host: this run used node budget 2 CPU 3g and client budget 1 CPU 2g, so budget comparisons are separate classes here, not independent hardware
- one client process with one fixed maintenance writer at 10 updates/s; readers are sequential; the reported sample latency includes retry and backoff time
- physical rank-index arm yb-range-rank; an engine whose planner ignores the index is a different physical design, not a faster one
- YugabyteDB DIST/DEBUG storage counters (rows and index rows scanned) are parsed per cell; RPC-level counters are not aggregated
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
