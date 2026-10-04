# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 91 cells, 91 valid, engine `yugabyte`, topology `single`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-v3-domain-yb-b1-20261005`
- run tag: `run/study-01/select-join-v3-domain-yb-b1-20261005`
- environment: `host-zenbook-ux5406sa`
- repo commit: `0f9d09212c6d52a80d8437d3b5a234f2f711ed36`
- engine version: `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit`
- database engine image: `docker.io/yugabytedb/yugabyte:2025.2.6.0-b111` @ `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`
- benchmark client image: `localhost/charitybench-select:select-join-v3-domain-yb-b1-20261005` @ `8436e36bb66caa8178559d54771b88c063e26848d4ef1d1a56eb524b7f0ecca6`
- inputs digest: `29e7348dd181bc43a61ff5d94d96dbb93927e8494bc11ba23a6de48ac27f3401`
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
| 20 (0.01%) | 2.3872 | 1.9405 | 1.6542 | 1.8763 | 1.6770 | 1.9144 | 2.1505 | 101.2866 | 2.1674 | 1.7749 |
| 200 (0.1%) | 3.6163 | 3.5227 | 3.5052 | 3.8908 | 3.0212 | 2.9039 | 3.6775 | 84.2620 | 2.8620 | 3.7938 |
| 2000 (1%) | 9.3009 | 9.0142 | 10.0048 | 8.9453 | 9.9517 | 9.0050 | 9.1630 | 81.7983 | 9.3869 | 9.7697 |
| 20000 (10%) | 66.6481 | 63.5215 | 62.4875 | 65.3698 | 67.7109 | 64.2026 | 63.0163 | 93.1901 | 62.3922 | 65.6356 |
| 100000 (50%) | 86.6911 | 87.7967 | 93.9461 | 84.2825 | 91.0411 | 85.5384 | 86.2173 | 92.1862 | 87.3422 | 84.8986 |
| 180000 (90%) | 106.8836 | 109.7228 | 105.7867 | 107.0039 | 112.1658 | 104.4277 | 104.5115 | 103.8122 | 109.0223 | 106.4771 |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | join_on vs join_where | 2.3872 | 1.9405 | 0.8129 | false |
| 20 | filtered_derived vs join_where | 2.3872 | 1.6542 | 0.6929 | false |
| 20 | cte_filtered vs join_where | 2.3872 | 1.8763 | 0.7860 | false |
| 20 | cte_filtered_mat vs join_where | 2.3872 | 1.6770 | 0.7025 | false |
| 20 | cte_filtered_notmat vs join_where | 2.3872 | 1.9144 | 0.8020 | false |
| 20 | cte_joined_outer vs join_where | 2.3872 | 2.1505 | 0.9009 | false |
| 20 | cte_joined_outer_mat vs join_where | 2.3872 | 101.2866 | 42.4288 | false |
| 20 | cte_joined_outer_notmat vs join_where | 2.3872 | 2.1674 | 0.9079 | false |
| 20 | comma_join vs join_where | 2.3872 | 1.7749 | 0.7435 | false |
| 200 | join_on vs join_where | 3.6163 | 3.5227 | 0.9741 | false |
| 200 | filtered_derived vs join_where | 3.6163 | 3.5052 | 0.9693 | false |
| 200 | cte_filtered vs join_where | 3.6163 | 3.8908 | 1.0759 | false |
| 200 | cte_filtered_mat vs join_where | 3.6163 | 3.0212 | 0.8354 | false |
| 200 | cte_filtered_notmat vs join_where | 3.6163 | 2.9039 | 0.8030 | false |
| 200 | cte_joined_outer vs join_where | 3.6163 | 3.6775 | 1.0169 | false |
| 200 | cte_joined_outer_mat vs join_where | 3.6163 | 84.2620 | 23.3009 | false |
| 200 | cte_joined_outer_notmat vs join_where | 3.6163 | 2.8620 | 0.7914 | false |
| 200 | comma_join vs join_where | 3.6163 | 3.7938 | 1.0491 | false |
| 2000 | join_on vs join_where | 9.3009 | 9.0142 | 0.9692 | false |
| 2000 | filtered_derived vs join_where | 9.3009 | 10.0048 | 1.0757 | false |
| 2000 | cte_filtered vs join_where | 9.3009 | 8.9453 | 0.9618 | false |
| 2000 | cte_filtered_mat vs join_where | 9.3009 | 9.9517 | 1.0700 | false |
| 2000 | cte_filtered_notmat vs join_where | 9.3009 | 9.0050 | 0.9682 | false |
| 2000 | cte_joined_outer vs join_where | 9.3009 | 9.1630 | 0.9852 | false |
| 2000 | cte_joined_outer_mat vs join_where | 9.3009 | 81.7983 | 8.7947 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 9.3009 | 9.3869 | 1.0092 | false |
| 2000 | comma_join vs join_where | 9.3009 | 9.7697 | 1.0504 | false |
| 20000 | join_on vs join_where | 66.6481 | 63.5215 | 0.9531 | false |
| 20000 | filtered_derived vs join_where | 66.6481 | 62.4875 | 0.9376 | false |
| 20000 | cte_filtered vs join_where | 66.6481 | 65.3698 | 0.9808 | false |
| 20000 | cte_filtered_mat vs join_where | 66.6481 | 67.7109 | 1.0159 | false |
| 20000 | cte_filtered_notmat vs join_where | 66.6481 | 64.2026 | 0.9633 | false |
| 20000 | cte_joined_outer vs join_where | 66.6481 | 63.0163 | 0.9455 | false |
| 20000 | cte_joined_outer_mat vs join_where | 66.6481 | 93.1901 | 1.3982 | false |
| 20000 | cte_joined_outer_notmat vs join_where | 66.6481 | 62.3922 | 0.9361 | false |
| 20000 | comma_join vs join_where | 66.6481 | 65.6356 | 0.9848 | false |
| 100000 | join_on vs join_where | 86.6911 | 87.7967 | 1.0128 | false |
| 100000 | filtered_derived vs join_where | 86.6911 | 93.9461 | 1.0837 | false |
| 100000 | cte_filtered vs join_where | 86.6911 | 84.2825 | 0.9722 | false |
| 100000 | cte_filtered_mat vs join_where | 86.6911 | 91.0411 | 1.0502 | false |
| 100000 | cte_filtered_notmat vs join_where | 86.6911 | 85.5384 | 0.9867 | false |
| 100000 | cte_joined_outer vs join_where | 86.6911 | 86.2173 | 0.9945 | false |
| 100000 | cte_joined_outer_mat vs join_where | 86.6911 | 92.1862 | 1.0634 | false |
| 100000 | cte_joined_outer_notmat vs join_where | 86.6911 | 87.3422 | 1.0075 | false |
| 100000 | comma_join vs join_where | 86.6911 | 84.8986 | 0.9793 | false |
| 180000 | join_on vs join_where | 106.8836 | 109.7228 | 1.0266 | false |
| 180000 | filtered_derived vs join_where | 106.8836 | 105.7867 | 0.9897 | false |
| 180000 | cte_filtered vs join_where | 106.8836 | 107.0039 | 1.0011 | false |
| 180000 | cte_filtered_mat vs join_where | 106.8836 | 112.1658 | 1.0494 | false |
| 180000 | cte_filtered_notmat vs join_where | 106.8836 | 104.4277 | 0.9770 | false |
| 180000 | cte_joined_outer vs join_where | 106.8836 | 104.5115 | 0.9778 | false |
| 180000 | cte_joined_outer_mat vs join_where | 106.8836 | 103.8122 | 0.9713 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 106.8836 | 109.0223 | 1.0200 | false |
| 180000 | comma_join vs join_where | 106.8836 | 106.4771 | 0.9962 | false |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_filtered_mat vs cte_filtered | 1.8763 | 1.6770 | 0.8938 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_filtered_notmat vs cte_filtered | 1.8763 | 1.9144 | 1.0203 | false | policy=not_materialized; observed_cte_scan=false |
| 20 | cte_joined_outer_mat vs cte_joined_outer | 2.1505 | 101.2866 | 47.0981 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 2.1505 | 2.1674 | 1.0079 | false | policy=not_materialized; observed_cte_scan=false |
| 200 | cte_filtered_mat vs cte_filtered | 3.8908 | 3.0212 | 0.7765 | false | policy=materialized; observed_cte_scan=true |
| 200 | cte_filtered_notmat vs cte_filtered | 3.8908 | 2.9039 | 0.7464 | false | policy=not_materialized; observed_cte_scan=false |
| 200 | cte_joined_outer_mat vs cte_joined_outer | 3.6775 | 84.2620 | 22.9131 | false | policy=materialized; observed_cte_scan=true |
| 200 | cte_joined_outer_notmat vs cte_joined_outer | 3.6775 | 2.8620 | 0.7783 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_filtered_mat vs cte_filtered | 8.9453 | 9.9517 | 1.1125 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_filtered_notmat vs cte_filtered | 8.9453 | 9.0050 | 1.0067 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 9.1630 | 81.7983 | 8.9270 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 9.1630 | 9.3869 | 1.0244 | false | policy=not_materialized; observed_cte_scan=false |
| 20000 | cte_filtered_mat vs cte_filtered | 65.3698 | 67.7109 | 1.0358 | false | policy=materialized; observed_cte_scan=true |
| 20000 | cte_filtered_notmat vs cte_filtered | 65.3698 | 64.2026 | 0.9821 | false | policy=not_materialized; observed_cte_scan=false |
| 20000 | cte_joined_outer_mat vs cte_joined_outer | 63.0163 | 93.1901 | 1.4788 | false | policy=materialized; observed_cte_scan=true |
| 20000 | cte_joined_outer_notmat vs cte_joined_outer | 63.0163 | 62.3922 | 0.9901 | false | policy=not_materialized; observed_cte_scan=false |
| 100000 | cte_filtered_mat vs cte_filtered | 84.2825 | 91.0411 | 1.0802 | false | policy=materialized; observed_cte_scan=true |
| 100000 | cte_filtered_notmat vs cte_filtered | 84.2825 | 85.5384 | 1.0149 | false | policy=not_materialized; observed_cte_scan=false |
| 100000 | cte_joined_outer_mat vs cte_joined_outer | 86.2173 | 92.1862 | 1.0692 | false | policy=materialized; observed_cte_scan=true |
| 100000 | cte_joined_outer_notmat vs cte_joined_outer | 86.2173 | 87.3422 | 1.0130 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_filtered_mat vs cte_filtered | 107.0039 | 112.1658 | 1.0482 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_filtered_notmat vs cte_filtered | 107.0039 | 104.4277 | 0.9759 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 104.5115 | 103.8122 | 0.9933 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 104.5115 | 109.0223 | 1.0432 | false | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 6.7656 | 1 | `a9332eaf2ad42b18` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 8.2660 | 1 | `7f48d3c1957bfc58` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | 9.1250 | 2 | `a7fc49a7ae038dc3` | policy=default; observed_cte_scan=true | valid |
| default | 2 | 100% | 8.7609 | 2 | `b58272dc2862fd34` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 1% | 8.5936 | 4 | `33eec4fdd01346d2` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 8.2410 | 4 | `d747ea33f9a082e5` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | 9.7113 | 8 | `4800e33100cfb1f7` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 100% | 9.1730 | 8 | `52f91f5f0c289c07` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 1% | 9.6230 | 16 | `bca2affb9432ad22` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 9.6473 | 16 | `b6e601431458dfd1` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 8.3224 | 1 | `f85c9c1ae479571c` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 9.2281 | 1 | `6e5c68a34d79aebf` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | 8.4599 | 2 | `14202011ca94ebf0` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 100% | 8.6643 | 2 | `2265b0d841dfde6d` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 1% | 9.0723 | 4 | `927deff15ecc8e6e` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 9.4691 | 4 | `161b4fdb10a1e4be` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | 9.6097 | 8 | `ca6ffb091d642ca7` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 100% | 9.2089 | 8 | `8daa6f780fd423e8` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 1% | 10.2537 | 16 | `a4bdc272edeb1872` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 9.1744 | 16 | `8f734d447a5963e2` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 7.9593 | 1 | `ae5a488114d5aeb6` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 9.0763 | 1 | `8c217c33e94cef86` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | 12.2144 | 2 | `8c00cfac968093e7` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 100% | 15.3048 | 2 | `b2d472dbaf00dce5` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 1% | 21.7410 | 4 | `666711723d486851` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 28.7119 | 4 | `446fd5f2c3dae322` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | 98.9907 | 8 | `bea7a85ae0b95814` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 100% | 89.5510 | 8 | `1cc6cec0116ad179` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 1% | 101.2243 | 16 | `1a40bdb8fa700bbd` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 84.6702 | 16 | `084d621c25e41bc0` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 1.2301 | false |
| notmat | 1 | 1% | 1.1764 | false |
| mat | 1 | 100% | 1.1164 | false |
| notmat | 1 | 100% | 1.0980 | false |
| mat | 2 | 1% | 0.9271 | false |
| notmat | 2 | 1% | 1.3386 | false |
| mat | 2 | 100% | 0.9890 | false |
| notmat | 2 | 100% | 1.7469 | false |
| mat | 4 | 1% | 1.0557 | false |
| notmat | 4 | 1% | 2.5299 | false |
| mat | 4 | 100% | 1.1490 | false |
| notmat | 4 | 100% | 3.4840 | false |
| mat | 8 | 1% | 0.9895 | false |
| notmat | 8 | 1% | 10.1934 | false |
| mat | 8 | 100% | 1.0039 | false |
| notmat | 8 | 100% | 9.7625 | false |
| mat | 16 | 1% | 1.0655 | false |
| notmat | 16 | 1% | 10.5190 | false |
| mat | 16 | 100% | 0.9510 | false |
| notmat | 16 | 100% | 8.7766 | false |

## Repeated-design position control

| first median ms | last median ms | ratio (last/first) | same plan |
|---|---|---|---|
| 93.190127 | 94.355959 | 1.0125 | false |

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

- client: nr_periods 621, nr_throttled 263, throttled_usec 9384626 | quota 100000 100000
, memory 2147483648

- yb-single: nr_periods 684, nr_throttled 75, throttled_usec 3328962 | quota 200000 100000, memory 3221225472

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
| family=bag;stmt=cte_joined_outer_mat;threshold=100000 | valid | false | true | 202050 | 0 | 0/0 | 100000 | 0 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | false | true | 202050 | 0 | 0/0 | 20000 | 0 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | false | true | 202050 | 0 | 0/0 | 199980 | 0 |
| family=bag;stmt=cte_joined_outer_mat;threshold=200 | valid | false | true | 202050 | 0 | 0/0 | 199800 | 0 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | false | true | 202050 | 0 | 0/0 | 198000 | 0 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20000 | valid | false | true | 202050 | 0 | 0/0 | 180000 | 0 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20000;position=last | valid | false | true | 202050 | 0 | 0/0 | 180000 | 0 |
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
| family=bag;stmt=comma_join;threshold=100000 | valid | 84.4849 | 84.8986 | 87.9122 | 5 | 11.8 |
| family=bag;stmt=comma_join;threshold=180000 | valid | 95.8126 | 106.4771 | 108.0722 | 5 | 9.4 |
| family=bag;stmt=comma_join;threshold=20 | valid | 1.7042 | 1.7749 | 1.8207 | 5 | 563.4 |
| family=bag;stmt=comma_join;threshold=200 | valid | 3.2927 | 3.7938 | 4.6357 | 5 | 263.6 |
| family=bag;stmt=comma_join;threshold=2000 | valid | 9.5999 | 9.7697 | 10.3931 | 5 | 102.4 |
| family=bag;stmt=comma_join;threshold=20000 | valid | 63.3356 | 65.6356 | 74.9361 | 5 | 15.2 |
| family=bag;stmt=cte_filtered;threshold=100000 | valid | 83.0731 | 84.2825 | 91.6928 | 5 | 11.9 |
| family=bag;stmt=cte_filtered;threshold=180000 | valid | 103.9446 | 107.0039 | 111.6940 | 5 | 9.3 |
| family=bag;stmt=cte_filtered;threshold=20 | valid | 1.7866 | 1.8763 | 2.0950 | 5 | 533.0 |
| family=bag;stmt=cte_filtered;threshold=200 | valid | 3.4122 | 3.8908 | 3.9987 | 5 | 257.0 |
| family=bag;stmt=cte_filtered;threshold=2000 | valid | 8.3992 | 8.9453 | 10.1578 | 5 | 111.8 |
| family=bag;stmt=cte_filtered;threshold=20000 | valid | 62.5977 | 65.3698 | 66.8149 | 5 | 15.3 |
| family=bag;stmt=cte_filtered_mat;threshold=100000 | valid | 86.2174 | 91.0411 | 95.2083 | 5 | 11.0 |
| family=bag;stmt=cte_filtered_mat;threshold=180000 | valid | 109.7331 | 112.1658 | 114.8727 | 5 | 8.9 |
| family=bag;stmt=cte_filtered_mat;threshold=20 | valid | 1.4813 | 1.6770 | 1.9433 | 5 | 596.3 |
| family=bag;stmt=cte_filtered_mat;threshold=200 | valid | 2.5647 | 3.0212 | 3.2915 | 5 | 331.0 |
| family=bag;stmt=cte_filtered_mat;threshold=2000 | valid | 9.7770 | 9.9517 | 11.1673 | 5 | 100.5 |
| family=bag;stmt=cte_filtered_mat;threshold=20000 | valid | 65.4870 | 67.7109 | 68.3417 | 5 | 14.8 |
| family=bag;stmt=cte_filtered_notmat;threshold=100000 | valid | 82.4057 | 85.5384 | 88.7351 | 5 | 11.7 |
| family=bag;stmt=cte_filtered_notmat;threshold=180000 | valid | 98.2055 | 104.4277 | 109.7123 | 5 | 9.6 |
| family=bag;stmt=cte_filtered_notmat;threshold=20 | valid | 1.7042 | 1.9144 | 19.4939 | 5 | 522.3 |
| family=bag;stmt=cte_filtered_notmat;threshold=200 | valid | 2.8588 | 2.9039 | 3.3269 | 5 | 344.4 |
| family=bag;stmt=cte_filtered_notmat;threshold=2000 | valid | 8.7253 | 9.0050 | 9.7652 | 5 | 111.0 |
| family=bag;stmt=cte_filtered_notmat;threshold=20000 | valid | 63.5733 | 64.2026 | 65.8067 | 5 | 15.6 |
| family=bag;stmt=cte_joined_outer;threshold=100000 | valid | 81.8106 | 86.2173 | 89.9250 | 5 | 11.6 |
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 100.2028 | 104.5115 | 110.4194 | 5 | 9.6 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 1.9636 | 2.1505 | 2.2741 | 5 | 465.0 |
| family=bag;stmt=cte_joined_outer;threshold=200 | valid | 3.5273 | 3.6775 | 4.1781 | 5 | 271.9 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 8.6811 | 9.1630 | 9.7386 | 5 | 109.1 |
| family=bag;stmt=cte_joined_outer;threshold=20000 | valid | 60.8252 | 63.0163 | 65.5403 | 5 | 15.9 |
| family=bag;stmt=cte_joined_outer_mat;threshold=100000 | valid | 89.2071 | 92.1862 | 104.1814 | 5 | 10.8 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 98.4924 | 103.8122 | 106.0718 | 5 | 9.6 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 82.0178 | 101.2866 | 117.9713 | 5 | 9.9 |
| family=bag;stmt=cte_joined_outer_mat;threshold=200 | valid | 77.3271 | 84.2620 | 163.5425 | 5 | 11.9 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 81.4811 | 81.7983 | 85.8977 | 5 | 12.2 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20000 | valid | 86.9525 | 93.1901 | 96.0822 | 5 | 10.7 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20000;position=last | valid | 92.8851 | 94.3560 | 96.7122 | 5 | 10.6 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=100000 | valid | 86.7165 | 87.3422 | 88.9668 | 5 | 11.4 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 106.7026 | 109.0223 | 122.2672 | 5 | 9.2 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 2.0394 | 2.1674 | 2.8592 | 5 | 461.4 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=200 | valid | 2.7973 | 2.8620 | 2.8994 | 5 | 349.4 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 8.9212 | 9.3869 | 9.8605 | 5 | 106.5 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20000 | valid | 60.1322 | 62.3922 | 64.3778 | 5 | 16.0 |
| family=bag;stmt=filtered_derived;threshold=100000 | valid | 88.9351 | 93.9461 | 98.4415 | 5 | 10.6 |
| family=bag;stmt=filtered_derived;threshold=180000 | valid | 99.4345 | 105.7867 | 114.0406 | 5 | 9.5 |
| family=bag;stmt=filtered_derived;threshold=20 | valid | 1.5352 | 1.6542 | 1.8081 | 5 | 604.5 |
| family=bag;stmt=filtered_derived;threshold=200 | valid | 3.4204 | 3.5052 | 3.5687 | 5 | 285.3 |
| family=bag;stmt=filtered_derived;threshold=2000 | valid | 9.6597 | 10.0048 | 10.3133 | 5 | 100.0 |
| family=bag;stmt=filtered_derived;threshold=20000 | valid | 60.7534 | 62.4875 | 66.0819 | 5 | 16.0 |
| family=bag;stmt=join_on;threshold=100000 | valid | 85.0233 | 87.7967 | 90.5665 | 5 | 11.4 |
| family=bag;stmt=join_on;threshold=180000 | valid | 106.9055 | 109.7228 | 119.5205 | 5 | 9.1 |
| family=bag;stmt=join_on;threshold=20 | valid | 1.7195 | 1.9405 | 2.2164 | 5 | 515.3 |
| family=bag;stmt=join_on;threshold=200 | valid | 3.4562 | 3.5227 | 3.9625 | 5 | 283.9 |
| family=bag;stmt=join_on;threshold=2000 | valid | 8.3606 | 9.0142 | 9.3388 | 5 | 110.9 |
| family=bag;stmt=join_on;threshold=20000 | valid | 59.8061 | 63.5215 | 64.7454 | 5 | 15.7 |
| family=bag;stmt=join_where;threshold=100000 | valid | 81.4263 | 86.6911 | 89.5103 | 5 | 11.5 |
| family=bag;stmt=join_where;threshold=180000 | valid | 105.1510 | 106.8836 | 111.4541 | 5 | 9.4 |
| family=bag;stmt=join_where;threshold=20 | valid | 2.1065 | 2.3872 | 2.5401 | 5 | 418.9 |
| family=bag;stmt=join_where;threshold=200 | valid | 3.4761 | 3.6163 | 3.9491 | 5 | 276.5 |
| family=bag;stmt=join_where;threshold=2000 | valid | 9.1908 | 9.3009 | 9.6705 | 5 | 107.5 |
| family=bag;stmt=join_where;threshold=20000 | valid | 64.1990 | 66.6481 | 69.7674 | 5 | 15.0 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 9.2432 | 9.6230 | 10.6783 | 5 | 103.9 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 9.2470 | 9.6473 | 10.0533 | 5 | 103.7 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 6.2287 | 6.7656 | 6.8397 | 5 | 147.8 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 7.4911 | 8.2660 | 8.9467 | 5 | 121.0 |
| family=reuse;policy=default;refs=2;consumer=1% | valid | 8.8000 | 9.1250 | 9.9493 | 5 | 109.6 |
| family=reuse;policy=default;refs=2;consumer=100% | valid | 8.3824 | 8.7609 | 9.1547 | 5 | 114.1 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 8.2378 | 8.5936 | 8.8705 | 5 | 116.4 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 8.0482 | 8.2410 | 8.7212 | 5 | 121.3 |
| family=reuse;policy=default;refs=8;consumer=1% | valid | 9.0207 | 9.7113 | 10.5436 | 5 | 103.0 |
| family=reuse;policy=default;refs=8;consumer=100% | valid | 8.4092 | 9.1730 | 9.7205 | 5 | 109.0 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 9.9810 | 10.2537 | 10.3258 | 5 | 97.5 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 9.0018 | 9.1744 | 9.5716 | 5 | 109.0 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 7.6249 | 8.3224 | 8.8213 | 5 | 120.2 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 8.5492 | 9.2281 | 9.4944 | 5 | 108.4 |
| family=reuse;policy=mat;refs=2;consumer=1% | valid | 7.9951 | 8.4599 | 8.6483 | 5 | 118.2 |
| family=reuse;policy=mat;refs=2;consumer=100% | valid | 7.6743 | 8.6643 | 8.9317 | 5 | 115.4 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 8.3618 | 9.0723 | 10.2774 | 5 | 110.2 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 9.2088 | 9.4691 | 9.5790 | 5 | 105.6 |
| family=reuse;policy=mat;refs=8;consumer=1% | valid | 9.2864 | 9.6097 | 10.0615 | 5 | 104.1 |
| family=reuse;policy=mat;refs=8;consumer=100% | valid | 9.0698 | 9.2089 | 9.6451 | 5 | 108.6 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 91.7136 | 101.2243 | 117.7879 | 5 | 9.9 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 82.5856 | 84.6702 | 118.3631 | 5 | 11.8 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 7.3235 | 7.9593 | 47.0698 | 5 | 125.6 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 8.5036 | 9.0763 | 9.4355 | 5 | 110.2 |
| family=reuse;policy=notmat;refs=2;consumer=1% | valid | 11.3348 | 12.2144 | 12.4731 | 5 | 81.9 |
| family=reuse;policy=notmat;refs=2;consumer=100% | valid | 14.8552 | 15.3048 | 16.0221 | 5 | 65.3 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 21.3698 | 21.7410 | 24.5601 | 5 | 46.0 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 27.2850 | 28.7119 | 32.2504 | 5 | 34.8 |
| family=reuse;policy=notmat;refs=8;consumer=1% | valid | 85.4092 | 98.9907 | 103.0741 | 5 | 10.1 |
| family=reuse;policy=notmat;refs=8;consumer=100% | valid | 79.3601 | 89.5510 | 92.3490 | 5 | 11.2 |

## Matched blended workload

- phase: `blended`; target offered rate: 10 updates/s
- offered 541, committed 541, failed 0 over 54.056 s
- offered rate 10.008/s, completed rate 10.008/s
- writer operations per endpoint: [541]
- reader operations overlapping a committed write: 249 of 637
- reader retries of transient conflicts: 0
- per key: 541 keys, min 1, max 1 per key
- ledger: 541 keys reconciled, 0 mismatch
- update latency: avg 3.487 ms, min 1.506, max 161.747 over 541

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42118878`, 5 interleaved rounds
  1. family=bag;stmt=cte_joined_outer_mat;threshold=20000
  2. family=bag;stmt=comma_join;threshold=20000
  3. family=bag;stmt=join_where;threshold=20000
  4. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  5. family=bag;stmt=cte_filtered_mat;threshold=2000
  6. family=bag;stmt=cte_filtered_notmat;threshold=200
  7. family=bag;stmt=join_on;threshold=2000
  8. family=bag;stmt=cte_filtered_mat;threshold=20000
  9. family=bag;stmt=cte_filtered;threshold=20000
  10. family=bag;stmt=comma_join;threshold=200
  11. family=bag;stmt=cte_filtered_notmat;threshold=2000
  12. family=bag;stmt=cte_filtered_notmat;threshold=20000
  13. family=reuse;policy=mat;refs=8;consumer=1%
  14. family=bag;stmt=cte_filtered_mat;threshold=180000
  15. family=reuse;policy=default;refs=2;consumer=1%
  16. family=bag;stmt=filtered_derived;threshold=20000
  17. family=bag;stmt=filtered_derived;threshold=200
  18. family=reuse;policy=default;refs=4;consumer=1%
  19. family=bag;stmt=cte_filtered;threshold=2000
  20. family=reuse;policy=notmat;refs=4;consumer=100%
  21. family=bag;stmt=cte_joined_outer_notmat;threshold=20000
  22. family=reuse;policy=notmat;refs=1;consumer=100%
  23. family=bag;stmt=cte_joined_outer;threshold=180000
  24. family=bag;stmt=cte_joined_outer_mat;threshold=100000
  25. family=bag;stmt=cte_joined_outer;threshold=20
  26. family=reuse;policy=mat;refs=8;consumer=100%
  27. family=reuse;policy=notmat;refs=16;consumer=100%
  28. family=bag;stmt=cte_joined_outer_mat;threshold=20
  29. family=bag;stmt=join_where;threshold=2000
  30. family=bag;stmt=cte_filtered;threshold=20
  31. family=reuse;policy=default;refs=2;consumer=100%
  32. family=bag;stmt=join_on;threshold=20
  33. family=bag;stmt=cte_filtered_mat;threshold=20
  34. family=bag;stmt=filtered_derived;threshold=20
  35. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  36. family=bag;stmt=filtered_derived;threshold=180000
  37. family=bag;stmt=join_where;threshold=20
  38. family=bag;stmt=cte_joined_outer;threshold=20000
  39. family=bag;stmt=join_on;threshold=200
  40. family=bag;stmt=join_on;threshold=20000
  41. family=bag;stmt=comma_join;threshold=180000
  42. family=reuse;policy=notmat;refs=1;consumer=1%
  43. family=reuse;policy=default;refs=1;consumer=100%
  44. family=reuse;policy=default;refs=1;consumer=1%
  45. family=reuse;policy=mat;refs=16;consumer=100%
  46. family=bag;stmt=cte_joined_outer_notmat;threshold=200
  47. family=bag;stmt=cte_filtered_notmat;threshold=100000
  48. family=bag;stmt=filtered_derived;threshold=2000
  49. family=bag;stmt=cte_joined_outer_notmat;threshold=100000
  50. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  51. family=bag;stmt=join_on;threshold=100000
  52. family=bag;stmt=cte_joined_outer;threshold=200
  53. family=bag;stmt=cte_joined_outer;threshold=2000
  54. family=reuse;policy=default;refs=4;consumer=100%
  55. family=reuse;policy=notmat;refs=2;consumer=100%
  56. family=bag;stmt=join_where;threshold=100000
  57. family=reuse;policy=mat;refs=1;consumer=100%
  58. family=reuse;policy=mat;refs=1;consumer=1%
  59. family=reuse;policy=default;refs=16;consumer=100%
  60. family=reuse;policy=notmat;refs=4;consumer=1%
  61. family=bag;stmt=cte_joined_outer;threshold=100000
  62. family=bag;stmt=cte_filtered_notmat;threshold=180000
  63. family=bag;stmt=cte_filtered;threshold=100000
  64. family=bag;stmt=join_on;threshold=180000
  65. family=bag;stmt=cte_filtered;threshold=200
  66. family=reuse;policy=default;refs=8;consumer=100%
  67. family=bag;stmt=cte_filtered_mat;threshold=100000
  68. family=bag;stmt=join_where;threshold=180000
  69. family=bag;stmt=cte_filtered;threshold=180000
  70. family=reuse;policy=mat;refs=4;consumer=100%
  71. family=bag;stmt=comma_join;threshold=100000
  72. family=bag;stmt=comma_join;threshold=2000
  73. family=bag;stmt=cte_filtered_notmat;threshold=20
  74. family=bag;stmt=cte_joined_outer_mat;threshold=200
  75. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  76. family=reuse;policy=mat;refs=2;consumer=1%
  77. family=bag;stmt=comma_join;threshold=20
  78. family=reuse;policy=notmat;refs=2;consumer=1%
  79. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  80. family=reuse;policy=mat;refs=4;consumer=1%
  81. family=reuse;policy=notmat;refs=16;consumer=1%
  82. family=reuse;policy=notmat;refs=8;consumer=1%
  83. family=reuse;policy=default;refs=8;consumer=1%
  84. family=reuse;policy=notmat;refs=8;consumer=100%
  85. family=bag;stmt=join_where;threshold=200
  86. family=bag;stmt=cte_filtered_mat;threshold=200
  87. family=reuse;policy=default;refs=16;consumer=1%
  88. family=bag;stmt=filtered_derived;threshold=100000
  89. family=reuse;policy=mat;refs=16;consumer=1%
  90. family=reuse;policy=mat;refs=2;consumer=100%
  91. family=bag;stmt=cte_joined_outer_mat;threshold=20000;position=last

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
| family=bag;stmt=comma_join;threshold=100000 | j06_comma_join | valid | `46c86791970a4095` |  |
| family=bag;stmt=comma_join;threshold=180000 | j06_comma_join | valid | `ad0ad470ca1fde61` |  |
| family=bag;stmt=comma_join;threshold=20 | j06_comma_join | valid | `04d22c80b1fe8978` |  |
| family=bag;stmt=comma_join;threshold=200 | j06_comma_join | valid | `458f574c6550734d` |  |
| family=bag;stmt=comma_join;threshold=2000 | j06_comma_join | valid | `02edda33b5684aaa` |  |
| family=bag;stmt=comma_join;threshold=20000 | j06_comma_join | valid | `98cd990949ccca76` |  |
| family=bag;stmt=cte_filtered;threshold=100000 | j04_cte_filtered | valid | `f0f4e2aeb7b19133` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=180000 | j04_cte_filtered | valid | `b4caf5580957b513` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=20 | j04_cte_filtered | valid | `e47f6bfe0eeeb3cc` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=200 | j04_cte_filtered | valid | `47520d38d6599ad8` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=2000 | j04_cte_filtered | valid | `cfd96ebbb1832811` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=20000 | j04_cte_filtered | valid | `bd6d4a4e7af9ea2b` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_mat;threshold=100000 | j04m_cte_filtered_mat | valid | `b43f9d12a5bae7d7` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=180000 | j04m_cte_filtered_mat | valid | `4e0995f28898d6a9` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=20 | j04m_cte_filtered_mat | valid | `1daa6f841590809d` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=200 | j04m_cte_filtered_mat | valid | `e897ddd9d67dbe8a` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=2000 | j04m_cte_filtered_mat | valid | `42f43def07e69d82` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=20000 | j04m_cte_filtered_mat | valid | `207d78ca6822f4cd` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_notmat;threshold=100000 | j04n_cte_filtered_notmat | valid | `299ec71557277551` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=180000 | j04n_cte_filtered_notmat | valid | `e35f3f865922f5dc` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=20 | j04n_cte_filtered_notmat | valid | `c69e498841332281` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=200 | j04n_cte_filtered_notmat | valid | `a027bb24599a63ce` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=2000 | j04n_cte_filtered_notmat | valid | `fb2a3c88b1cb4994` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=20000 | j04n_cte_filtered_notmat | valid | `73c83da62309edce` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=100000 | j05_cte_joined_outer | valid | `e57684ad350517a1` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=180000 | j05_cte_joined_outer | valid | `481e948e5dcb146a` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20 | j05_cte_joined_outer | valid | `9263c8269279aefb` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=200 | j05_cte_joined_outer | valid | `7ce2198e9d722814` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=2000 | j05_cte_joined_outer | valid | `956bccad7d8ee91f` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20000 | j05_cte_joined_outer | valid | `18cce2c952841c7f` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_mat;threshold=100000 | j05m_cte_joined_outer_mat | valid | `a2666cb6566d2100` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | j05m_cte_joined_outer_mat | valid | `59b4ee78868efe5f` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | j05m_cte_joined_outer_mat | valid | `a8045103d37203e6` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=200 | j05m_cte_joined_outer_mat | valid | `47789c44b6129065` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | j05m_cte_joined_outer_mat | valid | `5a1e57a7a8feabbb` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20000 | j05m_cte_joined_outer_mat | valid | `0a575b0eb42be716` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20000;position=last | j05m_cte_joined_outer_mat | valid | `c364f4e7bc4969a6` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_notmat;threshold=100000 | j05n_cte_joined_outer_notmat | valid | `e0bc669b6656a48f` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | j05n_cte_joined_outer_notmat | valid | `6f39964e6ab373cf` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | j05n_cte_joined_outer_notmat | valid | `a0ab792dab145a9c` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=200 | j05n_cte_joined_outer_notmat | valid | `b7f6cdc134396e18` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | j05n_cte_joined_outer_notmat | valid | `646368eb545b00f2` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20000 | j05n_cte_joined_outer_notmat | valid | `d67189eea561f1f3` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=filtered_derived;threshold=100000 | j03_filtered_derived | valid | `e994f984bcbdfadb` |  |
| family=bag;stmt=filtered_derived;threshold=180000 | j03_filtered_derived | valid | `4232b9d70527d36c` |  |
| family=bag;stmt=filtered_derived;threshold=20 | j03_filtered_derived | valid | `297848ff6e20b44f` |  |
| family=bag;stmt=filtered_derived;threshold=200 | j03_filtered_derived | valid | `7fcc9cf75eae861a` |  |
| family=bag;stmt=filtered_derived;threshold=2000 | j03_filtered_derived | valid | `cebf38cc456fe918` |  |
| family=bag;stmt=filtered_derived;threshold=20000 | j03_filtered_derived | valid | `d89492d4a58fdd6a` |  |
| family=bag;stmt=join_on;threshold=100000 | j02_join_on | valid | `1d2b5eab17cb5dcd` |  |
| family=bag;stmt=join_on;threshold=180000 | j02_join_on | valid | `75ea04dd5aa8212f` |  |
| family=bag;stmt=join_on;threshold=20 | j02_join_on | valid | `155dd5ce5235924b` |  |
| family=bag;stmt=join_on;threshold=200 | j02_join_on | valid | `4949ebde05ff7346` |  |
| family=bag;stmt=join_on;threshold=2000 | j02_join_on | valid | `b1a7f10d8b3b87b8` |  |
| family=bag;stmt=join_on;threshold=20000 | j02_join_on | valid | `e92a537130ba57e0` |  |
| family=bag;stmt=join_where;threshold=100000 | j01_join_where | valid | `b6d1c929e3e50529` |  |
| family=bag;stmt=join_where;threshold=180000 | j01_join_where | valid | `a0d2401daf294379` |  |
| family=bag;stmt=join_where;threshold=20 | j01_join_where | valid | `669d3a25b15bb296` |  |
| family=bag;stmt=join_where;threshold=200 | j01_join_where | valid | `a27878a378b74fe7` |  |
| family=bag;stmt=join_where;threshold=2000 | j01_join_where | valid | `efd5f0b7d7490617` |  |
| family=bag;stmt=join_where;threshold=20000 | j01_join_where | valid | `14ee786b0a468065` |  |
| family=reuse;policy=default;refs=16;consumer=1% | r_default_n16_c1 | valid | `bca2affb9432ad22` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=16;consumer=100% | r_default_n16_c100 | valid | `b6e601431458dfd1` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=1;consumer=1% | r_default_n1_c1 | valid | `a9332eaf2ad42b18` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=1;consumer=100% | r_default_n1_c100 | valid | `7f48d3c1957bfc58` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=2;consumer=1% | r_default_n2_c1 | valid | `a7fc49a7ae038dc3` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=2;consumer=100% | r_default_n2_c100 | valid | `b58272dc2862fd34` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=1% | r_default_n4_c1 | valid | `33eec4fdd01346d2` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=100% | r_default_n4_c100 | valid | `d747ea33f9a082e5` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=8;consumer=1% | r_default_n8_c1 | valid | `4800e33100cfb1f7` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=8;consumer=100% | r_default_n8_c100 | valid | `52f91f5f0c289c07` | policy=default; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=1% | r_mat_n16_c1 | valid | `a4bdc272edeb1872` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=100% | r_mat_n16_c100 | valid | `8f734d447a5963e2` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=1% | r_mat_n1_c1 | valid | `f85c9c1ae479571c` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=100% | r_mat_n1_c100 | valid | `6e5c68a34d79aebf` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=2;consumer=1% | r_mat_n2_c1 | valid | `14202011ca94ebf0` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=2;consumer=100% | r_mat_n2_c100 | valid | `2265b0d841dfde6d` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=1% | r_mat_n4_c1 | valid | `927deff15ecc8e6e` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=100% | r_mat_n4_c100 | valid | `161b4fdb10a1e4be` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=8;consumer=1% | r_mat_n8_c1 | valid | `ca6ffb091d642ca7` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=8;consumer=100% | r_mat_n8_c100 | valid | `8daa6f780fd423e8` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=notmat;refs=16;consumer=1% | r_notmat_n16_c1 | valid | `1a40bdb8fa700bbd` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `084d621c25e41bc0` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `ae5a488114d5aeb6` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `8c217c33e94cef86` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=2;consumer=1% | r_notmat_n2_c1 | valid | `8c00cfac968093e7` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=2;consumer=100% | r_notmat_n2_c100 | valid | `b2d472dbaf00dce5` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `666711723d486851` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `446fd5f2c3dae322` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=8;consumer=1% | r_notmat_n8_c1 | valid | `bea7a85ae0b95814` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=8;consumer=100% | r_notmat_n8_c100 | valid | `1cc6cec0116ad179` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, DIST, DEBUG, FORMAT TEXT` supported on yugabyte
- one host: this run used node budget 2 CPU 3g and client budget 1 CPU 2g, so budget comparisons are separate classes here, not independent hardware
- one client process with one fixed maintenance writer at 10 updates/s; readers are sequential; the reported sample latency includes retry and backoff time
- physical rank-index arm yb-range-rank; an engine whose planner ignores the index is a different physical design, not a faster one
- YugabyteDB DIST/DEBUG storage counters (rows and index rows scanned) are parsed per cell; RPC-level counters are not aggregated
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
