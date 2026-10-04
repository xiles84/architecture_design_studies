# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 91 cells, 91 valid, engine `yugabyte`, topology `single`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-domain2-yb2-20261004`
- run tag: `run/study-01/select-join-domain2-yb2-20261004`
- environment: `host-zenbook-ux5406sa`
- repo commit: `445fbc1646d58074c45cbaea57a83234132adc34`
- engine version: `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit`
- database engine image: `docker.io/yugabytedb/yugabyte:2025.2.6.0-b111` @ `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`
- benchmark client image: `localhost/charitybench-select:select-join-domain2-yb2-20261004` @ `79bc753ff7e57e102982a323939c43e68e9b3664887683d3a5d5bf072650f4c1`
- inputs digest: `0b8ec8dc842ac10e68c5fcb0659308b343a97c2bf0cd279c1f18cc70e744ad04`
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
| 20 (0.01%) | 57.1368 | 55.1729 | 61.3341 | 55.7243 | 53.4782 | 54.9420 | 58.4982 | 83.6480 | 58.6603 | 58.9074 |
| 200 (0.1%) | 60.1059 | 55.6316 | 55.5974 | 58.4827 | 60.7451 | 53.9860 | 55.7608 | 82.5734 | 61.7535 | 57.6281 |
| 2000 (1%) | 59.4739 | 57.3489 | 59.3681 | 57.5083 | 56.1911 | 57.1676 | 61.7402 | 88.6591 | 73.0405 | 58.2441 |
| 20000 (10%) | 70.0067 | 69.2121 | 73.2025 | 66.5384 | 72.5271 | 67.2228 | 66.6795 | 97.2148 | 68.3166 | 66.4834 |
| 100000 (50%) | 95.6175 | 85.9017 | 94.2605 | 89.3132 | 92.0424 | 90.3157 | 89.0346 | 99.1983 | 91.2664 | 87.0013 |
| 180000 (90%) | 119.4605 | 110.5148 | 110.0782 | 108.8930 | 114.2940 | 110.0115 | 108.8648 | 105.7473 | 110.8330 | 113.8692 |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | join_on vs join_where | 57.1368 | 55.1729 | 0.9656 | false |
| 20 | filtered_derived vs join_where | 57.1368 | 61.3341 | 1.0735 | false |
| 20 | cte_filtered vs join_where | 57.1368 | 55.7243 | 0.9753 | false |
| 20 | cte_filtered_mat vs join_where | 57.1368 | 53.4782 | 0.9360 | false |
| 20 | cte_filtered_notmat vs join_where | 57.1368 | 54.9420 | 0.9616 | false |
| 20 | cte_joined_outer vs join_where | 57.1368 | 58.4982 | 1.0238 | false |
| 20 | cte_joined_outer_mat vs join_where | 57.1368 | 83.6480 | 1.4640 | false |
| 20 | cte_joined_outer_notmat vs join_where | 57.1368 | 58.6603 | 1.0267 | false |
| 20 | comma_join vs join_where | 57.1368 | 58.9074 | 1.0310 | false |
| 200 | join_on vs join_where | 60.1059 | 55.6316 | 0.9256 | false |
| 200 | filtered_derived vs join_where | 60.1059 | 55.5974 | 0.9250 | false |
| 200 | cte_filtered vs join_where | 60.1059 | 58.4827 | 0.9730 | false |
| 200 | cte_filtered_mat vs join_where | 60.1059 | 60.7451 | 1.0106 | false |
| 200 | cte_filtered_notmat vs join_where | 60.1059 | 53.9860 | 0.8982 | false |
| 200 | cte_joined_outer vs join_where | 60.1059 | 55.7608 | 0.9277 | false |
| 200 | cte_joined_outer_mat vs join_where | 60.1059 | 82.5734 | 1.3738 | false |
| 200 | cte_joined_outer_notmat vs join_where | 60.1059 | 61.7535 | 1.0274 | false |
| 200 | comma_join vs join_where | 60.1059 | 57.6281 | 0.9588 | false |
| 2000 | join_on vs join_where | 59.4739 | 57.3489 | 0.9643 | false |
| 2000 | filtered_derived vs join_where | 59.4739 | 59.3681 | 0.9982 | false |
| 2000 | cte_filtered vs join_where | 59.4739 | 57.5083 | 0.9670 | false |
| 2000 | cte_filtered_mat vs join_where | 59.4739 | 56.1911 | 0.9448 | false |
| 2000 | cte_filtered_notmat vs join_where | 59.4739 | 57.1676 | 0.9612 | false |
| 2000 | cte_joined_outer vs join_where | 59.4739 | 61.7402 | 1.0381 | false |
| 2000 | cte_joined_outer_mat vs join_where | 59.4739 | 88.6591 | 1.4907 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 59.4739 | 73.0405 | 1.2281 | false |
| 2000 | comma_join vs join_where | 59.4739 | 58.2441 | 0.9793 | false |
| 20000 | join_on vs join_where | 70.0067 | 69.2121 | 0.9887 | false |
| 20000 | filtered_derived vs join_where | 70.0067 | 73.2025 | 1.0456 | false |
| 20000 | cte_filtered vs join_where | 70.0067 | 66.5384 | 0.9505 | false |
| 20000 | cte_filtered_mat vs join_where | 70.0067 | 72.5271 | 1.0360 | false |
| 20000 | cte_filtered_notmat vs join_where | 70.0067 | 67.2228 | 0.9602 | false |
| 20000 | cte_joined_outer vs join_where | 70.0067 | 66.6795 | 0.9525 | false |
| 20000 | cte_joined_outer_mat vs join_where | 70.0067 | 97.2148 | 1.3887 | false |
| 20000 | cte_joined_outer_notmat vs join_where | 70.0067 | 68.3166 | 0.9759 | false |
| 20000 | comma_join vs join_where | 70.0067 | 66.4834 | 0.9497 | false |
| 100000 | join_on vs join_where | 95.6175 | 85.9017 | 0.8984 | false |
| 100000 | filtered_derived vs join_where | 95.6175 | 94.2605 | 0.9858 | false |
| 100000 | cte_filtered vs join_where | 95.6175 | 89.3132 | 0.9341 | false |
| 100000 | cte_filtered_mat vs join_where | 95.6175 | 92.0424 | 0.9626 | false |
| 100000 | cte_filtered_notmat vs join_where | 95.6175 | 90.3157 | 0.9446 | false |
| 100000 | cte_joined_outer vs join_where | 95.6175 | 89.0346 | 0.9312 | false |
| 100000 | cte_joined_outer_mat vs join_where | 95.6175 | 99.1983 | 1.0375 | false |
| 100000 | cte_joined_outer_notmat vs join_where | 95.6175 | 91.2664 | 0.9545 | false |
| 100000 | comma_join vs join_where | 95.6175 | 87.0013 | 0.9099 | false |
| 180000 | join_on vs join_where | 119.4605 | 110.5148 | 0.9251 | false |
| 180000 | filtered_derived vs join_where | 119.4605 | 110.0782 | 0.9215 | false |
| 180000 | cte_filtered vs join_where | 119.4605 | 108.8930 | 0.9115 | false |
| 180000 | cte_filtered_mat vs join_where | 119.4605 | 114.2940 | 0.9568 | false |
| 180000 | cte_filtered_notmat vs join_where | 119.4605 | 110.0115 | 0.9209 | false |
| 180000 | cte_joined_outer vs join_where | 119.4605 | 108.8648 | 0.9113 | false |
| 180000 | cte_joined_outer_mat vs join_where | 119.4605 | 105.7473 | 0.8852 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 119.4605 | 110.8330 | 0.9278 | false |
| 180000 | comma_join vs join_where | 119.4605 | 113.8692 | 0.9532 | false |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_filtered_mat vs cte_filtered | 55.7243 | 53.4782 | 0.9597 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_filtered_notmat vs cte_filtered | 55.7243 | 54.9420 | 0.9860 | false | policy=not_materialized; observed_cte_scan=false |
| 20 | cte_joined_outer_mat vs cte_joined_outer | 58.4982 | 83.6480 | 1.4299 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 58.4982 | 58.6603 | 1.0028 | false | policy=not_materialized; observed_cte_scan=false |
| 200 | cte_filtered_mat vs cte_filtered | 58.4827 | 60.7451 | 1.0387 | false | policy=materialized; observed_cte_scan=true |
| 200 | cte_filtered_notmat vs cte_filtered | 58.4827 | 53.9860 | 0.9231 | false | policy=not_materialized; observed_cte_scan=false |
| 200 | cte_joined_outer_mat vs cte_joined_outer | 55.7608 | 82.5734 | 1.4808 | false | policy=materialized; observed_cte_scan=true |
| 200 | cte_joined_outer_notmat vs cte_joined_outer | 55.7608 | 61.7535 | 1.1075 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_filtered_mat vs cte_filtered | 57.5083 | 56.1911 | 0.9771 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_filtered_notmat vs cte_filtered | 57.5083 | 57.1676 | 0.9941 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 61.7402 | 88.6591 | 1.4360 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 61.7402 | 73.0405 | 1.1830 | false | policy=not_materialized; observed_cte_scan=false |
| 20000 | cte_filtered_mat vs cte_filtered | 66.5384 | 72.5271 | 1.0900 | false | policy=materialized; observed_cte_scan=true |
| 20000 | cte_filtered_notmat vs cte_filtered | 66.5384 | 67.2228 | 1.0103 | false | policy=not_materialized; observed_cte_scan=false |
| 20000 | cte_joined_outer_mat vs cte_joined_outer | 66.6795 | 97.2148 | 1.4579 | false | policy=materialized; observed_cte_scan=true |
| 20000 | cte_joined_outer_notmat vs cte_joined_outer | 66.6795 | 68.3166 | 1.0246 | false | policy=not_materialized; observed_cte_scan=false |
| 100000 | cte_filtered_mat vs cte_filtered | 89.3132 | 92.0424 | 1.0306 | false | policy=materialized; observed_cte_scan=true |
| 100000 | cte_filtered_notmat vs cte_filtered | 89.3132 | 90.3157 | 1.0112 | false | policy=not_materialized; observed_cte_scan=false |
| 100000 | cte_joined_outer_mat vs cte_joined_outer | 89.0346 | 99.1983 | 1.1142 | false | policy=materialized; observed_cte_scan=true |
| 100000 | cte_joined_outer_notmat vs cte_joined_outer | 89.0346 | 91.2664 | 1.0251 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_filtered_mat vs cte_filtered | 108.8930 | 114.2940 | 1.0496 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_filtered_notmat vs cte_filtered | 108.8930 | 110.0115 | 1.0103 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 108.8648 | 105.7473 | 0.9714 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 108.8648 | 110.8330 | 1.0181 | false | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 59.9066 | 1 | `c2ca16497e0c7991` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 58.0273 | 1 | `0e188abe8e7b1cf4` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | 58.8795 | 2 | `6051fe5a29f0b0af` | policy=default; observed_cte_scan=true | valid |
| default | 2 | 100% | 58.8155 | 2 | `0f579ee7fbb3222a` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 1% | 58.2933 | 4 | `31ae360618a083b3` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 58.1791 | 4 | `c54456f2c407372e` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | 55.6342 | 8 | `7683aeaee94075fd` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 100% | 55.8027 | 8 | `7325b32b5fa2afcc` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 1% | 60.9444 | 16 | `68e270e48780529c` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 57.9738 | 16 | `86687320ebe5006c` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 56.5462 | 1 | `2c672d04070b9f42` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 59.0213 | 1 | `9936eccb1c9c5ade` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | 59.4605 | 2 | `2c15128f05f293b5` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 100% | 58.2268 | 2 | `a6682769406b66e9` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 1% | 59.7968 | 4 | `71ce98aee8dfc66b` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 57.7894 | 4 | `e8a7503518d8b24f` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | 56.4921 | 8 | `f447d7d6db160b5d` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 100% | 55.6503 | 8 | `4d49a128ab652533` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 1% | 59.9515 | 16 | `4ea9fbfb31f16d0f` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 60.8010 | 16 | `54c6bbb6369e521f` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 57.5455 | 1 | `0dcfa6f772c4eec0` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 58.3425 | 1 | `c2c78965659f3570` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | 131.7976 | 2 | `954e53963eb61fad` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 100% | 119.4502 | 2 | `298d7650a2c9b146` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 1% | 155.2522 | 4 | `d3d0d4c0b9b8171d` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 166.5635 | 4 | `7d297235505274d3` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | 313.4809 | 8 | `953f0f83e718cf0d` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 100% | 306.2287 | 8 | `783f969b5a64457d` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 1% | 556.6217 | 16 | `a70e4ca90cb07535` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 565.3148 | 16 | `92e51bbe7e13dfbd` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 0.9439 | false |
| notmat | 1 | 1% | 0.9606 | false |
| mat | 1 | 100% | 1.0171 | false |
| notmat | 1 | 100% | 1.0054 | false |
| mat | 2 | 1% | 1.0099 | false |
| notmat | 2 | 1% | 2.2384 | false |
| mat | 2 | 100% | 0.9900 | false |
| notmat | 2 | 100% | 2.0309 | false |
| mat | 4 | 1% | 1.0258 | false |
| notmat | 4 | 1% | 2.6633 | false |
| mat | 4 | 100% | 0.9933 | false |
| notmat | 4 | 100% | 2.8629 | false |
| mat | 8 | 1% | 1.0154 | false |
| notmat | 8 | 1% | 5.6347 | false |
| mat | 8 | 100% | 0.9973 | false |
| notmat | 8 | 100% | 5.4877 | false |
| mat | 16 | 1% | 0.9837 | false |
| notmat | 16 | 1% | 9.1333 | false |
| mat | 16 | 100% | 1.0488 | false |
| notmat | 16 | 100% | 9.7512 | false |

## Repeated-design position control

| first median ms | last median ms | ratio (last/first) | same plan |
|---|---|---|---|
| 556.621749 | 617.038522 | 1.1085 | false |

## Latency samples and derived throughput

_The reader is sequential, so `queries/s` is 1000/median_ms, a derived rate and not a concurrent-capacity test. min/max are the five sample extremes._

| axis | status | min ms | median ms | max ms | samples | derived queries/s |
|---|---|---|---|---|---|---|
| family=bag;stmt=comma_join;threshold=100000 | valid | 86.3778 | 87.0013 | 96.6063 | 5 | 11.5 |
| family=bag;stmt=comma_join;threshold=180000 | valid | 107.9844 | 113.8692 | 115.1915 | 5 | 8.8 |
| family=bag;stmt=comma_join;threshold=20 | valid | 58.1513 | 58.9074 | 79.6202 | 5 | 17.0 |
| family=bag;stmt=comma_join;threshold=200 | valid | 55.3329 | 57.6281 | 61.0437 | 5 | 17.4 |
| family=bag;stmt=comma_join;threshold=2000 | valid | 56.9420 | 58.2441 | 65.6167 | 5 | 17.2 |
| family=bag;stmt=comma_join;threshold=20000 | valid | 65.7153 | 66.4834 | 68.3235 | 5 | 15.0 |
| family=bag;stmt=cte_filtered;threshold=100000 | valid | 86.5644 | 89.3132 | 92.6181 | 5 | 11.2 |
| family=bag;stmt=cte_filtered;threshold=180000 | valid | 106.6320 | 108.8930 | 113.6588 | 5 | 9.2 |
| family=bag;stmt=cte_filtered;threshold=20 | valid | 53.4133 | 55.7243 | 58.3098 | 5 | 17.9 |
| family=bag;stmt=cte_filtered;threshold=200 | valid | 56.1587 | 58.4827 | 60.3677 | 5 | 17.1 |
| family=bag;stmt=cte_filtered;threshold=2000 | valid | 53.9517 | 57.5083 | 58.0160 | 5 | 17.4 |
| family=bag;stmt=cte_filtered;threshold=20000 | valid | 65.0203 | 66.5384 | 69.5406 | 5 | 15.0 |
| family=bag;stmt=cte_filtered_mat;threshold=100000 | valid | 90.9918 | 92.0424 | 102.1980 | 5 | 10.9 |
| family=bag;stmt=cte_filtered_mat;threshold=180000 | valid | 113.1422 | 114.2940 | 132.7617 | 5 | 8.7 |
| family=bag;stmt=cte_filtered_mat;threshold=20 | valid | 52.5745 | 53.4782 | 54.5305 | 5 | 18.7 |
| family=bag;stmt=cte_filtered_mat;threshold=200 | valid | 60.3552 | 60.7451 | 67.9192 | 5 | 16.5 |
| family=bag;stmt=cte_filtered_mat;threshold=2000 | valid | 54.6394 | 56.1911 | 59.7654 | 5 | 17.8 |
| family=bag;stmt=cte_filtered_mat;threshold=20000 | valid | 71.0887 | 72.5271 | 74.9094 | 5 | 13.8 |
| family=bag;stmt=cte_filtered_notmat;threshold=100000 | valid | 89.2831 | 90.3157 | 100.9954 | 5 | 11.1 |
| family=bag;stmt=cte_filtered_notmat;threshold=180000 | valid | 107.4888 | 110.0115 | 120.3845 | 5 | 9.1 |
| family=bag;stmt=cte_filtered_notmat;threshold=20 | valid | 53.0638 | 54.9420 | 55.8199 | 5 | 18.2 |
| family=bag;stmt=cte_filtered_notmat;threshold=200 | valid | 52.5620 | 53.9860 | 59.4805 | 5 | 18.5 |
| family=bag;stmt=cte_filtered_notmat;threshold=2000 | valid | 54.8334 | 57.1676 | 62.1801 | 5 | 17.5 |
| family=bag;stmt=cte_filtered_notmat;threshold=20000 | valid | 66.1776 | 67.2228 | 67.2615 | 5 | 14.9 |
| family=bag;stmt=cte_joined_outer;threshold=100000 | valid | 87.3972 | 89.0346 | 90.5654 | 5 | 11.2 |
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 107.5861 | 108.8648 | 113.6735 | 5 | 9.2 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 56.4239 | 58.4982 | 64.6629 | 5 | 17.1 |
| family=bag;stmt=cte_joined_outer;threshold=200 | valid | 52.1845 | 55.7608 | 57.0306 | 5 | 17.9 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 58.2365 | 61.7402 | 66.1252 | 5 | 16.2 |
| family=bag;stmt=cte_joined_outer;threshold=20000 | valid | 66.0362 | 66.6795 | 71.0572 | 5 | 15.0 |
| family=bag;stmt=cte_joined_outer_mat;threshold=100000 | valid | 92.2242 | 99.1983 | 104.1240 | 5 | 10.1 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 99.3480 | 105.7473 | 111.7384 | 5 | 9.5 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 81.8967 | 83.6480 | 85.6258 | 5 | 12.0 |
| family=bag;stmt=cte_joined_outer_mat;threshold=200 | valid | 80.7979 | 82.5734 | 98.6535 | 5 | 12.1 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 84.3999 | 88.6591 | 96.9641 | 5 | 11.3 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20000 | valid | 95.1601 | 97.2148 | 101.1234 | 5 | 10.3 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=100000 | valid | 87.7002 | 91.2664 | 93.7315 | 5 | 11.0 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 104.2698 | 110.8330 | 114.8214 | 5 | 9.0 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 56.7233 | 58.6603 | 63.4232 | 5 | 17.0 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=200 | valid | 57.8229 | 61.7535 | 67.4667 | 5 | 16.2 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 56.6426 | 73.0405 | 85.3577 | 5 | 13.7 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20000 | valid | 68.0132 | 68.3166 | 75.6192 | 5 | 14.6 |
| family=bag;stmt=filtered_derived;threshold=100000 | valid | 93.3773 | 94.2605 | 100.3937 | 5 | 10.6 |
| family=bag;stmt=filtered_derived;threshold=180000 | valid | 104.4754 | 110.0782 | 113.8069 | 5 | 9.1 |
| family=bag;stmt=filtered_derived;threshold=20 | valid | 54.6898 | 61.3341 | 61.8204 | 5 | 16.3 |
| family=bag;stmt=filtered_derived;threshold=200 | valid | 53.0099 | 55.5974 | 58.0521 | 5 | 18.0 |
| family=bag;stmt=filtered_derived;threshold=2000 | valid | 58.0270 | 59.3681 | 61.0808 | 5 | 16.8 |
| family=bag;stmt=filtered_derived;threshold=20000 | valid | 69.8560 | 73.2025 | 83.0857 | 5 | 13.7 |
| family=bag;stmt=join_on;threshold=100000 | valid | 84.1644 | 85.9017 | 92.7619 | 5 | 11.6 |
| family=bag;stmt=join_on;threshold=180000 | valid | 106.0235 | 110.5148 | 124.0773 | 5 | 9.0 |
| family=bag;stmt=join_on;threshold=20 | valid | 53.2015 | 55.1729 | 57.2210 | 5 | 18.1 |
| family=bag;stmt=join_on;threshold=200 | valid | 54.9893 | 55.6316 | 59.4927 | 5 | 18.0 |
| family=bag;stmt=join_on;threshold=2000 | valid | 56.7379 | 57.3489 | 59.3582 | 5 | 17.4 |
| family=bag;stmt=join_on;threshold=20000 | valid | 66.5866 | 69.2121 | 71.7493 | 5 | 14.4 |
| family=bag;stmt=join_where;threshold=100000 | valid | 86.5853 | 95.6175 | 109.3758 | 5 | 10.5 |
| family=bag;stmt=join_where;threshold=180000 | valid | 115.0692 | 119.4605 | 148.2201 | 5 | 8.4 |
| family=bag;stmt=join_where;threshold=20 | valid | 56.6348 | 57.1368 | 60.4578 | 5 | 17.5 |
| family=bag;stmt=join_where;threshold=200 | valid | 57.2908 | 60.1059 | 66.3342 | 5 | 16.6 |
| family=bag;stmt=join_where;threshold=2000 | valid | 56.9705 | 59.4739 | 61.8621 | 5 | 16.8 |
| family=bag;stmt=join_where;threshold=20000 | valid | 69.2536 | 70.0067 | 73.2822 | 5 | 14.3 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 60.6303 | 60.9444 | 61.3860 | 5 | 16.4 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 56.3259 | 57.9738 | 62.9999 | 5 | 17.2 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 58.8156 | 59.9066 | 61.1706 | 5 | 16.7 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 55.4623 | 58.0273 | 58.5611 | 5 | 17.2 |
| family=reuse;policy=default;refs=2;consumer=1% | valid | 57.1439 | 58.8795 | 68.2020 | 5 | 17.0 |
| family=reuse;policy=default;refs=2;consumer=100% | valid | 56.5125 | 58.8155 | 62.3757 | 5 | 17.0 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 57.2187 | 58.2933 | 59.8976 | 5 | 17.2 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 55.0333 | 58.1791 | 58.7178 | 5 | 17.2 |
| family=reuse;policy=default;refs=8;consumer=1% | valid | 54.4256 | 55.6342 | 61.8605 | 5 | 18.0 |
| family=reuse;policy=default;refs=8;consumer=100% | valid | 54.3999 | 55.8027 | 58.1661 | 5 | 17.9 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 57.2313 | 59.9515 | 61.5703 | 5 | 16.7 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 58.5458 | 60.8010 | 63.9989 | 5 | 16.4 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 54.4028 | 56.5462 | 60.2040 | 5 | 17.7 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 55.7188 | 59.0213 | 70.6666 | 5 | 16.9 |
| family=reuse;policy=mat;refs=2;consumer=1% | valid | 57.8433 | 59.4605 | 62.1001 | 5 | 16.8 |
| family=reuse;policy=mat;refs=2;consumer=100% | valid | 56.4409 | 58.2268 | 64.8386 | 5 | 17.2 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 52.8391 | 59.7968 | 60.3651 | 5 | 16.7 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 56.6045 | 57.7894 | 58.4625 | 5 | 17.3 |
| family=reuse;policy=mat;refs=8;consumer=1% | valid | 56.0073 | 56.4921 | 59.4876 | 5 | 17.7 |
| family=reuse;policy=mat;refs=8;consumer=100% | valid | 54.1071 | 55.6503 | 57.8594 | 5 | 18.0 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 511.2845 | 556.6217 | 589.6828 | 5 | 1.8 |
| family=reuse;policy=notmat;refs=16;consumer=1%;position=last | valid | 604.2663 | 617.0385 | 650.5493 | 5 | 1.6 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 547.5045 | 565.3148 | 679.3618 | 5 | 1.8 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 55.8799 | 57.5455 | 59.2503 | 5 | 17.4 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 56.3819 | 58.3425 | 67.2190 | 5 | 17.1 |
| family=reuse;policy=notmat;refs=2;consumer=1% | valid | 120.1833 | 131.7976 | 167.5725 | 5 | 7.6 |
| family=reuse;policy=notmat;refs=2;consumer=100% | valid | 106.1597 | 119.4502 | 132.7100 | 5 | 8.4 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 138.8920 | 155.2522 | 175.8515 | 5 | 6.4 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 158.1108 | 166.5635 | 189.9686 | 5 | 6.0 |
| family=reuse;policy=notmat;refs=8;consumer=1% | valid | 301.4019 | 313.4809 | 360.5572 | 5 | 3.2 |
| family=reuse;policy=notmat;refs=8;consumer=100% | valid | 300.9967 | 306.2287 | 331.4483 | 5 | 3.3 |

## Matched blended workload

- phase: `isolated`; target offered rate: 10 updates/s
- offered 0, committed 0, failed 0 over 0.000 s
- offered rate 0.000/s, completed rate 0.000/s
- per key: no keys touched
- ledger: 0 keys reconciled, 0 mismatch

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42047777`, 5 interleaved rounds
  1. family=reuse;policy=notmat;refs=16;consumer=1%
  2. family=bag;stmt=comma_join;threshold=20
  3. family=bag;stmt=cte_joined_outer_mat;threshold=20000
  4. family=reuse;policy=default;refs=16;consumer=1%
  5. family=reuse;policy=default;refs=1;consumer=1%
  6. family=bag;stmt=cte_filtered_mat;threshold=20000
  7. family=bag;stmt=join_where;threshold=2000
  8. family=bag;stmt=comma_join;threshold=200
  9. family=bag;stmt=filtered_derived;threshold=2000
  10. family=bag;stmt=join_where;threshold=100000
  11. family=reuse;policy=notmat;refs=4;consumer=1%
  12. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  13. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  14. family=reuse;policy=default;refs=2;consumer=1%
  15. family=bag;stmt=cte_joined_outer;threshold=20
  16. family=reuse;policy=mat;refs=2;consumer=100%
  17. family=bag;stmt=cte_filtered_notmat;threshold=100000
  18. family=bag;stmt=comma_join;threshold=2000
  19. family=bag;stmt=cte_joined_outer_notmat;threshold=20000
  20. family=reuse;policy=default;refs=4;consumer=1%
  21. family=bag;stmt=cte_joined_outer_mat;threshold=100000
  22. family=reuse;policy=mat;refs=16;consumer=100%
  23. family=bag;stmt=join_where;threshold=20
  24. family=reuse;policy=notmat;refs=1;consumer=1%
  25. family=bag;stmt=cte_filtered;threshold=200
  26. family=reuse;policy=default;refs=4;consumer=100%
  27. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  28. family=bag;stmt=cte_filtered_notmat;threshold=180000
  29. family=reuse;policy=default;refs=8;consumer=1%
  30. family=bag;stmt=cte_filtered_mat;threshold=100000
  31. family=reuse;policy=mat;refs=4;consumer=1%
  32. family=bag;stmt=filtered_derived;threshold=180000
  33. family=bag;stmt=cte_filtered_mat;threshold=20
  34. family=reuse;policy=default;refs=1;consumer=100%
  35. family=bag;stmt=join_on;threshold=100000
  36. family=bag;stmt=comma_join;threshold=100000
  37. family=bag;stmt=cte_filtered_mat;threshold=2000
  38. family=bag;stmt=cte_filtered;threshold=20
  39. family=bag;stmt=join_on;threshold=200
  40. family=bag;stmt=join_on;threshold=2000
  41. family=bag;stmt=cte_joined_outer;threshold=20000
  42. family=bag;stmt=cte_filtered_notmat;threshold=2000
  43. family=reuse;policy=default;refs=16;consumer=100%
  44. family=bag;stmt=cte_filtered_mat;threshold=180000
  45. family=reuse;policy=notmat;refs=1;consumer=100%
  46. family=reuse;policy=notmat;refs=16;consumer=100%
  47. family=reuse;policy=notmat;refs=2;consumer=1%
  48. family=bag;stmt=filtered_derived;threshold=20000
  49. family=reuse;policy=mat;refs=1;consumer=100%
  50. family=reuse;policy=notmat;refs=8;consumer=1%
  51. family=bag;stmt=cte_joined_outer_notmat;threshold=200
  52. family=reuse;policy=mat;refs=2;consumer=1%
  53. family=reuse;policy=notmat;refs=2;consumer=100%
  54. family=bag;stmt=join_where;threshold=200
  55. family=bag;stmt=cte_joined_outer;threshold=2000
  56. family=reuse;policy=notmat;refs=4;consumer=100%
  57. family=bag;stmt=join_where;threshold=180000
  58. family=reuse;policy=notmat;refs=8;consumer=100%
  59. family=bag;stmt=cte_filtered_mat;threshold=200
  60. family=reuse;policy=mat;refs=16;consumer=1%
  61. family=bag;stmt=filtered_derived;threshold=20
  62. family=bag;stmt=join_where;threshold=20000
  63. family=bag;stmt=filtered_derived;threshold=100000
  64. family=reuse;policy=mat;refs=4;consumer=100%
  65. family=bag;stmt=cte_joined_outer_notmat;threshold=100000
  66. family=reuse;policy=default;refs=2;consumer=100%
  67. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  68. family=bag;stmt=join_on;threshold=180000
  69. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  70. family=bag;stmt=filtered_derived;threshold=200
  71. family=reuse;policy=mat;refs=8;consumer=100%
  72. family=bag;stmt=cte_filtered_notmat;threshold=200
  73. family=bag;stmt=cte_filtered;threshold=100000
  74. family=bag;stmt=cte_joined_outer;threshold=180000
  75. family=bag;stmt=cte_filtered_notmat;threshold=20000
  76. family=bag;stmt=cte_joined_outer;threshold=100000
  77. family=bag;stmt=join_on;threshold=20
  78. family=bag;stmt=join_on;threshold=20000
  79. family=reuse;policy=mat;refs=1;consumer=1%
  80. family=bag;stmt=cte_filtered;threshold=180000
  81. family=bag;stmt=cte_filtered_notmat;threshold=20
  82. family=bag;stmt=cte_joined_outer;threshold=200
  83. family=reuse;policy=mat;refs=8;consumer=1%
  84. family=bag;stmt=cte_filtered;threshold=20000
  85. family=bag;stmt=comma_join;threshold=180000
  86. family=bag;stmt=comma_join;threshold=20000
  87. family=bag;stmt=cte_joined_outer_mat;threshold=200
  88. family=bag;stmt=cte_filtered;threshold=2000
  89. family=bag;stmt=cte_joined_outer_mat;threshold=20
  90. family=reuse;policy=default;refs=8;consumer=100%
  91. family=reuse;policy=notmat;refs=16;consumer=1%;position=last

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
| family=bag;stmt=comma_join;threshold=100000 | j06_comma_join | valid | `60968cfa64c0c0ae` |  |
| family=bag;stmt=comma_join;threshold=180000 | j06_comma_join | valid | `0ad768a9f15a8a73` |  |
| family=bag;stmt=comma_join;threshold=20 | j06_comma_join | valid | `623b1faaf525e6fe` |  |
| family=bag;stmt=comma_join;threshold=200 | j06_comma_join | valid | `d952be70c4683e2c` |  |
| family=bag;stmt=comma_join;threshold=2000 | j06_comma_join | valid | `d227675832d4d635` |  |
| family=bag;stmt=comma_join;threshold=20000 | j06_comma_join | valid | `ad0422888d410181` |  |
| family=bag;stmt=cte_filtered;threshold=100000 | j04_cte_filtered | valid | `cd5c6a5f625cd787` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=180000 | j04_cte_filtered | valid | `b557a50e043c2309` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=20 | j04_cte_filtered | valid | `b15e047309c5ab0f` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=200 | j04_cte_filtered | valid | `f79293c51a53df56` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=2000 | j04_cte_filtered | valid | `f2b89f4cbc7e0d8f` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=20000 | j04_cte_filtered | valid | `c6dd1949286fd636` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_mat;threshold=100000 | j04m_cte_filtered_mat | valid | `06b63298b05c4971` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=180000 | j04m_cte_filtered_mat | valid | `a41e7edc827536b8` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=20 | j04m_cte_filtered_mat | valid | `121497c3d959c9d4` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=200 | j04m_cte_filtered_mat | valid | `6b10042d56b00f32` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=2000 | j04m_cte_filtered_mat | valid | `26ec53acb537c31e` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=20000 | j04m_cte_filtered_mat | valid | `e70106084389c83c` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_notmat;threshold=100000 | j04n_cte_filtered_notmat | valid | `2c57cde7d3bff4ee` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=180000 | j04n_cte_filtered_notmat | valid | `672525adbfbe8260` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=20 | j04n_cte_filtered_notmat | valid | `db8aba44326193e9` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=200 | j04n_cte_filtered_notmat | valid | `1c3a28a39b8e48f9` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=2000 | j04n_cte_filtered_notmat | valid | `e8a41359da638991` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=20000 | j04n_cte_filtered_notmat | valid | `3ec7dcb16083a31a` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=100000 | j05_cte_joined_outer | valid | `a6c7af8f7a09902b` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=180000 | j05_cte_joined_outer | valid | `4a2078ba81ebc522` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20 | j05_cte_joined_outer | valid | `121200c950a0dab6` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=200 | j05_cte_joined_outer | valid | `7e9d5fa5c0680996` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=2000 | j05_cte_joined_outer | valid | `0d0f3fdbe5da8ef1` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20000 | j05_cte_joined_outer | valid | `75d684ad113c415e` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_mat;threshold=100000 | j05m_cte_joined_outer_mat | valid | `a0e9c049b4d0156f` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | j05m_cte_joined_outer_mat | valid | `e36ad29bccce0cd6` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | j05m_cte_joined_outer_mat | valid | `45b7c327a67a2860` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=200 | j05m_cte_joined_outer_mat | valid | `38fc8858d5608aca` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | j05m_cte_joined_outer_mat | valid | `5ff48837519fa322` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20000 | j05m_cte_joined_outer_mat | valid | `8fe79639cdb6b01d` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_notmat;threshold=100000 | j05n_cte_joined_outer_notmat | valid | `c13365af0df25122` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | j05n_cte_joined_outer_notmat | valid | `786c126dd6a47195` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | j05n_cte_joined_outer_notmat | valid | `a606be9db88b56b8` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=200 | j05n_cte_joined_outer_notmat | valid | `ec0dea9099c26268` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | j05n_cte_joined_outer_notmat | valid | `d8c092d4b4dea464` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20000 | j05n_cte_joined_outer_notmat | valid | `e3e4fd69ccf7fb68` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=filtered_derived;threshold=100000 | j03_filtered_derived | valid | `4ba8de5cbe647f8f` |  |
| family=bag;stmt=filtered_derived;threshold=180000 | j03_filtered_derived | valid | `9323aebed6963e70` |  |
| family=bag;stmt=filtered_derived;threshold=20 | j03_filtered_derived | valid | `ced1216912c79527` |  |
| family=bag;stmt=filtered_derived;threshold=200 | j03_filtered_derived | valid | `66ada109f9335997` |  |
| family=bag;stmt=filtered_derived;threshold=2000 | j03_filtered_derived | valid | `0bd651374271b126` |  |
| family=bag;stmt=filtered_derived;threshold=20000 | j03_filtered_derived | valid | `6a7cd560d296f098` |  |
| family=bag;stmt=join_on;threshold=100000 | j02_join_on | valid | `2895535286804f38` |  |
| family=bag;stmt=join_on;threshold=180000 | j02_join_on | valid | `6fd5b66b6b2d62c6` |  |
| family=bag;stmt=join_on;threshold=20 | j02_join_on | valid | `bb653cfb099be57b` |  |
| family=bag;stmt=join_on;threshold=200 | j02_join_on | valid | `e5fb2d073474c377` |  |
| family=bag;stmt=join_on;threshold=2000 | j02_join_on | valid | `038a931d857e5fa9` |  |
| family=bag;stmt=join_on;threshold=20000 | j02_join_on | valid | `e5d58e22658c449b` |  |
| family=bag;stmt=join_where;threshold=100000 | j01_join_where | valid | `625f6104200c5bde` |  |
| family=bag;stmt=join_where;threshold=180000 | j01_join_where | valid | `1e59cf28a85348c8` |  |
| family=bag;stmt=join_where;threshold=20 | j01_join_where | valid | `fdd1f3430e2d00d2` |  |
| family=bag;stmt=join_where;threshold=200 | j01_join_where | valid | `d6f683d83eaa6573` |  |
| family=bag;stmt=join_where;threshold=2000 | j01_join_where | valid | `a3a2a0ee626c6f2a` |  |
| family=bag;stmt=join_where;threshold=20000 | j01_join_where | valid | `311216aa0e56b7e5` |  |
| family=reuse;policy=default;refs=16;consumer=1% | r_default_n16_c1 | valid | `68e270e48780529c` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=16;consumer=100% | r_default_n16_c100 | valid | `86687320ebe5006c` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=1;consumer=1% | r_default_n1_c1 | valid | `c2ca16497e0c7991` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=1;consumer=100% | r_default_n1_c100 | valid | `0e188abe8e7b1cf4` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=2;consumer=1% | r_default_n2_c1 | valid | `6051fe5a29f0b0af` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=2;consumer=100% | r_default_n2_c100 | valid | `0f579ee7fbb3222a` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=1% | r_default_n4_c1 | valid | `31ae360618a083b3` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=100% | r_default_n4_c100 | valid | `c54456f2c407372e` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=8;consumer=1% | r_default_n8_c1 | valid | `7683aeaee94075fd` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=8;consumer=100% | r_default_n8_c100 | valid | `7325b32b5fa2afcc` | policy=default; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=1% | r_mat_n16_c1 | valid | `4ea9fbfb31f16d0f` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=100% | r_mat_n16_c100 | valid | `54c6bbb6369e521f` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=1% | r_mat_n1_c1 | valid | `2c672d04070b9f42` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=100% | r_mat_n1_c100 | valid | `9936eccb1c9c5ade` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=2;consumer=1% | r_mat_n2_c1 | valid | `2c15128f05f293b5` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=2;consumer=100% | r_mat_n2_c100 | valid | `a6682769406b66e9` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=1% | r_mat_n4_c1 | valid | `71ce98aee8dfc66b` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=100% | r_mat_n4_c100 | valid | `e8a7503518d8b24f` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=8;consumer=1% | r_mat_n8_c1 | valid | `f447d7d6db160b5d` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=8;consumer=100% | r_mat_n8_c100 | valid | `4d49a128ab652533` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=notmat;refs=16;consumer=1% | r_notmat_n16_c1 | valid | `a70e4ca90cb07535` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=16;consumer=1%;position=last | r_notmat_n16_c1 | valid | `976bd38a5a6cd2bf` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `92e51bbe7e13dfbd` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `0dcfa6f772c4eec0` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `c2c78965659f3570` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=2;consumer=1% | r_notmat_n2_c1 | valid | `954e53963eb61fad` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=2;consumer=100% | r_notmat_n2_c100 | valid | `298d7650a2c9b146` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `d3d0d4c0b9b8171d` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `7d297235505274d3` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=8;consumer=1% | r_notmat_n8_c1 | valid | `953f0f83e718cf0d` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=8;consumer=100% | r_notmat_n8_c100 | valid | `783f969b5a64457d` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, DIST, DEBUG, FORMAT TEXT` supported on yugabyte
- one host, single client: no independent hardware, real-network or equal-total-budget comparison
- no concurrent writers: the isolated read contract only; a matched blended read/write phase is a recorded gap
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
