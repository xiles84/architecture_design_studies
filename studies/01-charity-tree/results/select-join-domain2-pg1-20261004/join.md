# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 91 cells, 91 valid, engine `postgres`, topology `single`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-domain2-pg1-20261004`
- run tag: `run/study-01/select-join-domain2-pg1-20261004`
- environment: `host-zenbook-ux5406sa`
- repo commit: `445fbc1646d58074c45cbaea57a83234132adc34`
- engine version: `PostgreSQL 17.11 (Debian 17.11-1.pgdg13+2) on x86_64-pc-linux-gnu, compiled by gcc (Debian 14.2.0-19) 14.2.0, 64-bit`
- database engine image: `docker.io/library/postgres:17.11` @ `7296f210ae81031ec955dbad9a67a84fe958572a2153b8d0826a647522904dc1`
- benchmark client image: `localhost/charitybench-select:select-join-domain2-pg1-20261004` @ `79bc753ff7e57e102982a323939c43e68e9b3664887683d3a5d5bf072650f4c1`
- inputs digest: `347d12b60a3bffc23ae0193584db62e42e0b967778dd90d8c6ccd293ddf7ac1a`
- dataset fingerprint: `fcf26e08f377ec51`
- client CPU quota: `100000 100000`, memory: `2147483648`
- setting `client_workers`: 1 sequential reader; 0 writers
- setting `connection_pool`: 2 per endpoint = reader/setup + reserve
- setting `customers`: 2000
- setting `dataset`: study01-person-donation-charity
- setting `distributed_rpc`: N/A on PostgreSQL single-node
- setting `endpoint_mode`: single endpoint
- setting `endpoints`: pg-single
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
| pg-single | 914 | 1 |

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
| 20 (0.01%) | 0.4248 | 0.6485 | 0.4764 | 0.4150 | 0.5204 | 0.4540 | 0.6436 | 43.0594 | 0.6972 | 0.7149 |
| 200 (0.1%) | 0.7394 | 0.6980 | 0.8135 | 0.9459 | 0.7336 | 0.8541 | 1.0300 | 42.9176 | 0.7058 | 0.9324 |
| 2000 (1%) | 1.7501 | 1.9257 | 1.6899 | 2.3856 | 2.1761 | 1.7814 | 1.8746 | 48.9844 | 1.8568 | 1.5387 |
| 20000 (10%) | 6.9815 | 7.3674 | 7.0450 | 7.4676 | 9.6814 | 7.1539 | 6.7494 | 52.8927 | 6.8425 | 6.6844 |
| 100000 (50%) | 53.9293 | 57.7548 | 51.5489 | 48.3719 | 62.8347 | 46.4440 | 57.1731 | 67.7068 | 49.0801 | 56.6840 |
| 180000 (90%) | 50.3270 | 68.2629 | 68.8957 | 79.7254 | 73.8423 | 78.2505 | 76.3655 | 90.3074 | 76.7771 | 75.9659 |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | join_on vs join_where | 0.4248 | 0.6485 | 1.5266 | true |
| 20 | filtered_derived vs join_where | 0.4248 | 0.4764 | 1.1214 | true |
| 20 | cte_filtered vs join_where | 0.4248 | 0.4150 | 0.9769 | true |
| 20 | cte_filtered_mat vs join_where | 0.4248 | 0.5204 | 1.2249 | false |
| 20 | cte_filtered_notmat vs join_where | 0.4248 | 0.4540 | 1.0688 | true |
| 20 | cte_joined_outer vs join_where | 0.4248 | 0.6436 | 1.5149 | true |
| 20 | cte_joined_outer_mat vs join_where | 0.4248 | 43.0594 | 101.3585 | false |
| 20 | cte_joined_outer_notmat vs join_where | 0.4248 | 0.6972 | 1.6411 | true |
| 20 | comma_join vs join_where | 0.4248 | 0.7149 | 1.6828 | true |
| 200 | join_on vs join_where | 0.7394 | 0.6980 | 0.9440 | true |
| 200 | filtered_derived vs join_where | 0.7394 | 0.8135 | 1.1003 | true |
| 200 | cte_filtered vs join_where | 0.7394 | 0.9459 | 1.2794 | true |
| 200 | cte_filtered_mat vs join_where | 0.7394 | 0.7336 | 0.9921 | false |
| 200 | cte_filtered_notmat vs join_where | 0.7394 | 0.8541 | 1.1552 | true |
| 200 | cte_joined_outer vs join_where | 0.7394 | 1.0300 | 1.3931 | true |
| 200 | cte_joined_outer_mat vs join_where | 0.7394 | 42.9176 | 58.0448 | false |
| 200 | cte_joined_outer_notmat vs join_where | 0.7394 | 0.7058 | 0.9546 | true |
| 200 | comma_join vs join_where | 0.7394 | 0.9324 | 1.2611 | true |
| 2000 | join_on vs join_where | 1.7501 | 1.9257 | 1.1004 | true |
| 2000 | filtered_derived vs join_where | 1.7501 | 1.6899 | 0.9656 | true |
| 2000 | cte_filtered vs join_where | 1.7501 | 2.3856 | 1.3631 | true |
| 2000 | cte_filtered_mat vs join_where | 1.7501 | 2.1761 | 1.2434 | false |
| 2000 | cte_filtered_notmat vs join_where | 1.7501 | 1.7814 | 1.0179 | true |
| 2000 | cte_joined_outer vs join_where | 1.7501 | 1.8746 | 1.0712 | true |
| 2000 | cte_joined_outer_mat vs join_where | 1.7501 | 48.9844 | 27.9900 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 1.7501 | 1.8568 | 1.0610 | true |
| 2000 | comma_join vs join_where | 1.7501 | 1.5387 | 0.8792 | true |
| 20000 | join_on vs join_where | 6.9815 | 7.3674 | 1.0553 | true |
| 20000 | filtered_derived vs join_where | 6.9815 | 7.0450 | 1.0091 | true |
| 20000 | cte_filtered vs join_where | 6.9815 | 7.4676 | 1.0696 | true |
| 20000 | cte_filtered_mat vs join_where | 6.9815 | 9.6814 | 1.3867 | false |
| 20000 | cte_filtered_notmat vs join_where | 6.9815 | 7.1539 | 1.0247 | true |
| 20000 | cte_joined_outer vs join_where | 6.9815 | 6.7494 | 0.9668 | true |
| 20000 | cte_joined_outer_mat vs join_where | 6.9815 | 52.8927 | 7.5761 | false |
| 20000 | cte_joined_outer_notmat vs join_where | 6.9815 | 6.8425 | 0.9801 | true |
| 20000 | comma_join vs join_where | 6.9815 | 6.6844 | 0.9574 | true |
| 100000 | join_on vs join_where | 53.9293 | 57.7548 | 1.0709 | true |
| 100000 | filtered_derived vs join_where | 53.9293 | 51.5489 | 0.9559 | true |
| 100000 | cte_filtered vs join_where | 53.9293 | 48.3719 | 0.8969 | true |
| 100000 | cte_filtered_mat vs join_where | 53.9293 | 62.8347 | 1.1651 | false |
| 100000 | cte_filtered_notmat vs join_where | 53.9293 | 46.4440 | 0.8612 | true |
| 100000 | cte_joined_outer vs join_where | 53.9293 | 57.1731 | 1.0601 | true |
| 100000 | cte_joined_outer_mat vs join_where | 53.9293 | 67.7068 | 1.2555 | false |
| 100000 | cte_joined_outer_notmat vs join_where | 53.9293 | 49.0801 | 0.9101 | true |
| 100000 | comma_join vs join_where | 53.9293 | 56.6840 | 1.0511 | true |
| 180000 | join_on vs join_where | 50.3270 | 68.2629 | 1.3564 | true |
| 180000 | filtered_derived vs join_where | 50.3270 | 68.8957 | 1.3690 | true |
| 180000 | cte_filtered vs join_where | 50.3270 | 79.7254 | 1.5841 | true |
| 180000 | cte_filtered_mat vs join_where | 50.3270 | 73.8423 | 1.4673 | false |
| 180000 | cte_filtered_notmat vs join_where | 50.3270 | 78.2505 | 1.5548 | true |
| 180000 | cte_joined_outer vs join_where | 50.3270 | 76.3655 | 1.5174 | true |
| 180000 | cte_joined_outer_mat vs join_where | 50.3270 | 90.3074 | 1.7944 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 50.3270 | 76.7771 | 1.5256 | true |
| 180000 | comma_join vs join_where | 50.3270 | 75.9659 | 1.5094 | true |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_filtered_mat vs cte_filtered | 0.4150 | 0.5204 | 1.2539 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_filtered_notmat vs cte_filtered | 0.4150 | 0.4540 | 1.0941 | true | policy=not_materialized; observed_cte_scan=false |
| 20 | cte_joined_outer_mat vs cte_joined_outer | 0.6436 | 43.0594 | 66.9062 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 0.6436 | 0.6972 | 1.0832 | true | policy=not_materialized; observed_cte_scan=false |
| 200 | cte_filtered_mat vs cte_filtered | 0.9459 | 0.7336 | 0.7755 | false | policy=materialized; observed_cte_scan=true |
| 200 | cte_filtered_notmat vs cte_filtered | 0.9459 | 0.8541 | 0.9030 | true | policy=not_materialized; observed_cte_scan=false |
| 200 | cte_joined_outer_mat vs cte_joined_outer | 1.0300 | 42.9176 | 41.6661 | false | policy=materialized; observed_cte_scan=true |
| 200 | cte_joined_outer_notmat vs cte_joined_outer | 1.0300 | 0.7058 | 0.6853 | true | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_filtered_mat vs cte_filtered | 2.3856 | 2.1761 | 0.9122 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_filtered_notmat vs cte_filtered | 2.3856 | 1.7814 | 0.7467 | true | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 1.8746 | 48.9844 | 26.1299 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 1.8746 | 1.8568 | 0.9905 | true | policy=not_materialized; observed_cte_scan=false |
| 20000 | cte_filtered_mat vs cte_filtered | 7.4676 | 9.6814 | 1.2965 | false | policy=materialized; observed_cte_scan=true |
| 20000 | cte_filtered_notmat vs cte_filtered | 7.4676 | 7.1539 | 0.9580 | true | policy=not_materialized; observed_cte_scan=false |
| 20000 | cte_joined_outer_mat vs cte_joined_outer | 6.7494 | 52.8927 | 7.8366 | false | policy=materialized; observed_cte_scan=true |
| 20000 | cte_joined_outer_notmat vs cte_joined_outer | 6.7494 | 6.8425 | 1.0138 | true | policy=not_materialized; observed_cte_scan=false |
| 100000 | cte_filtered_mat vs cte_filtered | 48.3719 | 62.8347 | 1.2990 | false | policy=materialized; observed_cte_scan=true |
| 100000 | cte_filtered_notmat vs cte_filtered | 48.3719 | 46.4440 | 0.9601 | true | policy=not_materialized; observed_cte_scan=false |
| 100000 | cte_joined_outer_mat vs cte_joined_outer | 57.1731 | 67.7068 | 1.1842 | false | policy=materialized; observed_cte_scan=true |
| 100000 | cte_joined_outer_notmat vs cte_joined_outer | 57.1731 | 49.0801 | 0.8584 | true | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_filtered_mat vs cte_filtered | 79.7254 | 73.8423 | 0.9262 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_filtered_notmat vs cte_filtered | 79.7254 | 78.2505 | 0.9815 | true | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 76.3655 | 90.3074 | 1.1826 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 76.3655 | 76.7771 | 1.0054 | true | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 0.6730 | 1 | `549a07108ab01ba9` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 2.2777 | 1 | `2ce8b85a7b60c88e` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | 2.8907 | 2 | `b307c41eb54cf452` | policy=default; observed_cte_scan=true | valid |
| default | 2 | 100% | 2.3431 | 2 | `b307c41eb54cf452` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 1% | 2.5007 | 4 | `ccfe1af393deae98` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 2.4659 | 4 | `ccfe1af393deae98` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | 2.6444 | 8 | `c8c73778097931dc` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 100% | 3.2083 | 8 | `c8c73778097931dc` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 1% | 3.2078 | 16 | `b0c356323d30c425` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 3.7671 | 16 | `b0c356323d30c425` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 2.4710 | 1 | `eaa7d6058aaa5e03` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 2.7059 | 1 | `eaa7d6058aaa5e03` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | 2.7794 | 2 | `b307c41eb54cf452` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 100% | 2.6278 | 2 | `b307c41eb54cf452` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 1% | 2.5610 | 4 | `ccfe1af393deae98` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 2.8516 | 4 | `ccfe1af393deae98` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | 2.6540 | 8 | `c8c73778097931dc` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 100% | 3.0394 | 8 | `c8c73778097931dc` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 1% | 3.0843 | 16 | `b0c356323d30c425` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 3.3902 | 16 | `b0c356323d30c425` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 0.5032 | 1 | `549a07108ab01ba9` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 2.4461 | 1 | `2ce8b85a7b60c88e` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | 0.6465 | 2 | `28cc913d71d0d72b` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 100% | 4.5430 | 2 | `2043b765d5001655` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 1% | 1.0950 | 4 | `9f850e53560baafb` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 6.7414 | 4 | `cc53feafb616350b` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | 1.6683 | 8 | `477f7fbaf1c627cb` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 100% | 9.0577 | 8 | `714b97b0f1e43d36` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 1% | 2.5541 | 16 | `1c38b58d8cec08d3` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 15.0087 | 16 | `effce31a81556487` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 3.6715 | false |
| notmat | 1 | 1% | 0.7477 | true |
| mat | 1 | 100% | 1.1880 | false |
| notmat | 1 | 100% | 1.0739 | true |
| mat | 2 | 1% | 0.9615 | true |
| notmat | 2 | 1% | 0.2236 | false |
| mat | 2 | 100% | 1.1215 | true |
| notmat | 2 | 100% | 1.9389 | false |
| mat | 4 | 1% | 1.0241 | true |
| notmat | 4 | 1% | 0.4379 | false |
| mat | 4 | 100% | 1.1564 | true |
| notmat | 4 | 100% | 2.7339 | false |
| mat | 8 | 1% | 1.0036 | true |
| notmat | 8 | 1% | 0.6309 | false |
| mat | 8 | 100% | 0.9474 | true |
| notmat | 8 | 100% | 2.8232 | false |
| mat | 16 | 1% | 0.9615 | true |
| notmat | 16 | 1% | 0.7962 | false |
| mat | 16 | 100% | 0.9000 | true |
| notmat | 16 | 100% | 3.9842 | false |

## Repeated-design position control

| first median ms | last median ms | ratio (last/first) | same plan |
|---|---|---|---|
| 2.851551 | 2.422852 | 0.8497 | true |

## Latency samples and derived throughput

_The reader is sequential, so `queries/s` is 1000/median_ms, a derived rate and not a concurrent-capacity test. min/max are the five sample extremes._

| axis | status | min ms | median ms | max ms | samples | derived queries/s |
|---|---|---|---|---|---|---|
| family=bag;stmt=comma_join;threshold=100000 | valid | 48.2640 | 56.6840 | 67.9431 | 5 | 17.6 |
| family=bag;stmt=comma_join;threshold=180000 | valid | 44.4336 | 75.9659 | 90.6343 | 5 | 13.2 |
| family=bag;stmt=comma_join;threshold=20 | valid | 0.6248 | 0.7149 | 0.9403 | 5 | 1398.8 |
| family=bag;stmt=comma_join;threshold=200 | valid | 0.8641 | 0.9324 | 0.9621 | 5 | 1072.5 |
| family=bag;stmt=comma_join;threshold=2000 | valid | 1.4675 | 1.5387 | 1.7496 | 5 | 649.9 |
| family=bag;stmt=comma_join;threshold=20000 | valid | 6.4698 | 6.6844 | 7.4745 | 5 | 149.6 |
| family=bag;stmt=cte_filtered;threshold=100000 | valid | 47.0135 | 48.3719 | 86.8807 | 5 | 20.7 |
| family=bag;stmt=cte_filtered;threshold=180000 | valid | 44.6886 | 79.7254 | 98.7604 | 5 | 12.5 |
| family=bag;stmt=cte_filtered;threshold=20 | valid | 0.3658 | 0.4150 | 0.6357 | 5 | 2409.7 |
| family=bag;stmt=cte_filtered;threshold=200 | valid | 0.7744 | 0.9459 | 1.0783 | 5 | 1057.2 |
| family=bag;stmt=cte_filtered;threshold=2000 | valid | 1.8808 | 2.3856 | 2.9825 | 5 | 419.2 |
| family=bag;stmt=cte_filtered;threshold=20000 | valid | 6.8424 | 7.4676 | 47.7550 | 5 | 133.9 |
| family=bag;stmt=cte_filtered_mat;threshold=100000 | valid | 53.9543 | 62.8347 | 73.1772 | 5 | 15.9 |
| family=bag;stmt=cte_filtered_mat;threshold=180000 | valid | 66.0116 | 73.8423 | 76.5814 | 5 | 13.5 |
| family=bag;stmt=cte_filtered_mat;threshold=20 | valid | 0.4883 | 0.5204 | 6.4247 | 5 | 1921.7 |
| family=bag;stmt=cte_filtered_mat;threshold=200 | valid | 0.6954 | 0.7336 | 0.8086 | 5 | 1363.2 |
| family=bag;stmt=cte_filtered_mat;threshold=2000 | valid | 2.0044 | 2.1761 | 2.6945 | 5 | 459.5 |
| family=bag;stmt=cte_filtered_mat;threshold=20000 | valid | 8.6410 | 9.6814 | 13.0391 | 5 | 103.3 |
| family=bag;stmt=cte_filtered_notmat;threshold=100000 | valid | 44.4873 | 46.4440 | 48.0420 | 5 | 21.5 |
| family=bag;stmt=cte_filtered_notmat;threshold=180000 | valid | 47.9443 | 78.2505 | 93.7206 | 5 | 12.8 |
| family=bag;stmt=cte_filtered_notmat;threshold=20 | valid | 0.4477 | 0.4540 | 0.5858 | 5 | 2202.5 |
| family=bag;stmt=cte_filtered_notmat;threshold=200 | valid | 0.7860 | 0.8541 | 11.4471 | 5 | 1170.8 |
| family=bag;stmt=cte_filtered_notmat;threshold=2000 | valid | 1.5860 | 1.7814 | 2.0202 | 5 | 561.4 |
| family=bag;stmt=cte_filtered_notmat;threshold=20000 | valid | 6.8636 | 7.1539 | 39.4158 | 5 | 139.8 |
| family=bag;stmt=cte_joined_outer;threshold=100000 | valid | 45.5083 | 57.1731 | 68.3337 | 5 | 17.5 |
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 46.4412 | 76.3655 | 93.1709 | 5 | 13.1 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 0.5399 | 0.6436 | 0.6680 | 5 | 1553.8 |
| family=bag;stmt=cte_joined_outer;threshold=200 | valid | 0.9043 | 1.0300 | 1.2895 | 5 | 970.8 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 1.7914 | 1.8746 | 2.2231 | 5 | 533.4 |
| family=bag;stmt=cte_joined_outer;threshold=20000 | valid | 6.2084 | 6.7494 | 7.5665 | 5 | 148.2 |
| family=bag;stmt=cte_joined_outer_mat;threshold=100000 | valid | 64.3378 | 67.7068 | 69.7718 | 5 | 14.8 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 77.3917 | 90.3074 | 110.8286 | 5 | 11.1 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 41.2032 | 43.0594 | 45.8176 | 5 | 23.2 |
| family=bag;stmt=cte_joined_outer_mat;threshold=200 | valid | 41.9541 | 42.9176 | 53.6686 | 5 | 23.3 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 43.8922 | 48.9844 | 57.5370 | 5 | 20.4 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20000 | valid | 48.5958 | 52.8927 | 59.6996 | 5 | 18.9 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=100000 | valid | 46.9735 | 49.0801 | 62.0432 | 5 | 20.4 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 44.2921 | 76.7771 | 92.4957 | 5 | 13.0 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 0.6672 | 0.6972 | 0.7570 | 5 | 1434.4 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=200 | valid | 0.6391 | 0.7058 | 0.9277 | 5 | 1416.7 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 1.8144 | 1.8568 | 2.6743 | 5 | 538.5 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20000 | valid | 6.4681 | 6.8425 | 7.3936 | 5 | 146.1 |
| family=bag;stmt=filtered_derived;threshold=100000 | valid | 44.0376 | 51.5489 | 72.6165 | 5 | 19.4 |
| family=bag;stmt=filtered_derived;threshold=180000 | valid | 46.1431 | 68.8957 | 92.2850 | 5 | 14.5 |
| family=bag;stmt=filtered_derived;threshold=20 | valid | 0.4547 | 0.4764 | 0.7762 | 5 | 2099.1 |
| family=bag;stmt=filtered_derived;threshold=200 | valid | 0.7692 | 0.8135 | 0.9104 | 5 | 1229.2 |
| family=bag;stmt=filtered_derived;threshold=2000 | valid | 1.5901 | 1.6899 | 1.7878 | 5 | 591.8 |
| family=bag;stmt=filtered_derived;threshold=20000 | valid | 6.8463 | 7.0450 | 8.8874 | 5 | 141.9 |
| family=bag;stmt=join_on;threshold=100000 | valid | 42.8406 | 57.7548 | 78.5830 | 5 | 17.3 |
| family=bag;stmt=join_on;threshold=180000 | valid | 50.0189 | 68.2629 | 75.1394 | 5 | 14.6 |
| family=bag;stmt=join_on;threshold=20 | valid | 0.5464 | 0.6485 | 0.6863 | 5 | 1542.0 |
| family=bag;stmt=join_on;threshold=200 | valid | 0.6662 | 0.6980 | 0.7942 | 5 | 1432.7 |
| family=bag;stmt=join_on;threshold=2000 | valid | 1.8404 | 1.9257 | 3.4931 | 5 | 519.3 |
| family=bag;stmt=join_on;threshold=20000 | valid | 6.9046 | 7.3674 | 8.8128 | 5 | 135.7 |
| family=bag;stmt=join_where;threshold=100000 | valid | 44.8585 | 53.9293 | 62.2804 | 5 | 18.5 |
| family=bag;stmt=join_where;threshold=180000 | valid | 50.1016 | 50.3270 | 78.5792 | 5 | 19.9 |
| family=bag;stmt=join_where;threshold=20 | valid | 0.4120 | 0.4248 | 0.5182 | 5 | 2353.9 |
| family=bag;stmt=join_where;threshold=200 | valid | 0.6820 | 0.7394 | 0.7651 | 5 | 1352.5 |
| family=bag;stmt=join_where;threshold=2000 | valid | 1.6706 | 1.7501 | 2.1947 | 5 | 571.4 |
| family=bag;stmt=join_where;threshold=20000 | valid | 6.5240 | 6.9815 | 7.3726 | 5 | 143.2 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 3.1703 | 3.2078 | 3.6748 | 5 | 311.7 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 3.4996 | 3.7671 | 4.0061 | 5 | 265.5 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 0.5841 | 0.6730 | 0.6994 | 5 | 1485.8 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 2.0666 | 2.2777 | 2.9402 | 5 | 439.0 |
| family=reuse;policy=default;refs=2;consumer=1% | valid | 2.6712 | 2.8907 | 3.5097 | 5 | 345.9 |
| family=reuse;policy=default;refs=2;consumer=100% | valid | 2.2562 | 2.3431 | 2.5121 | 5 | 426.8 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 2.4308 | 2.5007 | 3.1235 | 5 | 399.9 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 2.4138 | 2.4659 | 3.3035 | 5 | 405.5 |
| family=reuse;policy=default;refs=8;consumer=1% | valid | 2.5141 | 2.6444 | 2.8190 | 5 | 378.2 |
| family=reuse;policy=default;refs=8;consumer=100% | valid | 2.9895 | 3.2083 | 3.3291 | 5 | 311.7 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 2.8098 | 3.0843 | 3.5479 | 5 | 324.2 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 3.2291 | 3.3902 | 3.5696 | 5 | 295.0 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 2.4397 | 2.4710 | 2.6026 | 5 | 404.7 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 2.5694 | 2.7059 | 3.4865 | 5 | 369.6 |
| family=reuse;policy=mat;refs=2;consumer=1% | valid | 2.5920 | 2.7794 | 2.8045 | 5 | 359.8 |
| family=reuse;policy=mat;refs=2;consumer=100% | valid | 2.4994 | 2.6278 | 2.6969 | 5 | 380.5 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 2.4140 | 2.5610 | 3.1902 | 5 | 390.5 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 2.7022 | 2.8516 | 3.2126 | 5 | 350.7 |
| family=reuse;policy=mat;refs=4;consumer=100%;position=last | valid | 2.1871 | 2.4229 | 2.6596 | 5 | 412.7 |
| family=reuse;policy=mat;refs=8;consumer=1% | valid | 2.4998 | 2.6540 | 2.7205 | 5 | 376.8 |
| family=reuse;policy=mat;refs=8;consumer=100% | valid | 2.9272 | 3.0394 | 3.6932 | 5 | 329.0 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 2.4435 | 2.5541 | 3.0245 | 5 | 391.5 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 14.5233 | 15.0087 | 15.8481 | 5 | 66.6 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 0.4344 | 0.5032 | 0.5439 | 5 | 1987.2 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 2.3473 | 2.4461 | 2.5573 | 5 | 408.8 |
| family=reuse;policy=notmat;refs=2;consumer=1% | valid | 0.6232 | 0.6465 | 0.8103 | 5 | 1546.9 |
| family=reuse;policy=notmat;refs=2;consumer=100% | valid | 4.3482 | 4.5430 | 5.2364 | 5 | 220.1 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 1.0172 | 1.0950 | 1.1494 | 5 | 913.3 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 6.4126 | 6.7414 | 7.2979 | 5 | 148.3 |
| family=reuse;policy=notmat;refs=8;consumer=1% | valid | 1.4478 | 1.6683 | 1.8816 | 5 | 599.4 |
| family=reuse;policy=notmat;refs=8;consumer=100% | valid | 8.4599 | 9.0577 | 10.9626 | 5 | 110.4 |

## Matched blended workload

- phase: `isolated`; target offered rate: 10 updates/s
- offered 0, committed 0, failed 0 over 0.000 s
- offered rate 0.000/s, completed rate 0.000/s
- per key: no keys touched
- ledger: 0 keys reconciled, 0 mismatch

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42874128`, 5 interleaved rounds
  1. family=reuse;policy=mat;refs=4;consumer=100%
  2. family=reuse;policy=notmat;refs=1;consumer=1%
  3. family=reuse;policy=mat;refs=8;consumer=1%
  4. family=bag;stmt=cte_joined_outer_notmat;threshold=200
  5. family=reuse;policy=default;refs=2;consumer=100%
  6. family=bag;stmt=cte_filtered_mat;threshold=180000
  7. family=bag;stmt=cte_joined_outer;threshold=100000
  8. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  9. family=bag;stmt=join_where;threshold=2000
  10. family=bag;stmt=cte_joined_outer;threshold=180000
  11. family=reuse;policy=notmat;refs=4;consumer=1%
  12. family=bag;stmt=filtered_derived;threshold=100000
  13. family=bag;stmt=cte_filtered_mat;threshold=2000
  14. family=bag;stmt=cte_joined_outer;threshold=20000
  15. family=bag;stmt=cte_filtered;threshold=200
  16. family=bag;stmt=comma_join;threshold=100000
  17. family=bag;stmt=cte_joined_outer;threshold=200
  18. family=reuse;policy=mat;refs=4;consumer=1%
  19. family=reuse;policy=notmat;refs=2;consumer=1%
  20. family=reuse;policy=notmat;refs=2;consumer=100%
  21. family=bag;stmt=filtered_derived;threshold=2000
  22. family=bag;stmt=cte_filtered_mat;threshold=200
  23. family=reuse;policy=mat;refs=16;consumer=1%
  24. family=reuse;policy=notmat;refs=4;consumer=100%
  25. family=bag;stmt=cte_joined_outer_mat;threshold=20
  26. family=bag;stmt=cte_filtered_notmat;threshold=100000
  27. family=reuse;policy=mat;refs=1;consumer=1%
  28. family=bag;stmt=cte_filtered_notmat;threshold=180000
  29. family=bag;stmt=cte_filtered_notmat;threshold=20000
  30. family=bag;stmt=cte_joined_outer_mat;threshold=100000
  31. family=bag;stmt=comma_join;threshold=20
  32. family=bag;stmt=cte_filtered_notmat;threshold=20
  33. family=bag;stmt=cte_filtered_notmat;threshold=2000
  34. family=bag;stmt=cte_joined_outer_mat;threshold=200
  35. family=bag;stmt=join_on;threshold=20000
  36. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  37. family=reuse;policy=mat;refs=2;consumer=100%
  38. family=bag;stmt=cte_joined_outer_notmat;threshold=20000
  39. family=bag;stmt=join_on;threshold=2000
  40. family=bag;stmt=join_where;threshold=180000
  41. family=reuse;policy=mat;refs=2;consumer=1%
  42. family=bag;stmt=comma_join;threshold=180000
  43. family=reuse;policy=notmat;refs=8;consumer=1%
  44. family=bag;stmt=cte_filtered;threshold=20000
  45. family=bag;stmt=cte_joined_outer;threshold=20
  46. family=bag;stmt=cte_filtered;threshold=20
  47. family=bag;stmt=join_where;threshold=20000
  48. family=reuse;policy=notmat;refs=1;consumer=100%
  49. family=reuse;policy=notmat;refs=8;consumer=100%
  50. family=bag;stmt=cte_filtered_mat;threshold=20
  51. family=bag;stmt=cte_joined_outer_mat;threshold=20000
  52. family=reuse;policy=default;refs=16;consumer=100%
  53. family=reuse;policy=default;refs=1;consumer=100%
  54. family=bag;stmt=join_on;threshold=200
  55. family=reuse;policy=default;refs=4;consumer=1%
  56. family=bag;stmt=cte_filtered_mat;threshold=20000
  57. family=reuse;policy=default;refs=8;consumer=100%
  58. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  59. family=reuse;policy=default;refs=1;consumer=1%
  60. family=bag;stmt=cte_filtered;threshold=180000
  61. family=reuse;policy=default;refs=16;consumer=1%
  62. family=bag;stmt=filtered_derived;threshold=180000
  63. family=bag;stmt=join_on;threshold=20
  64. family=bag;stmt=filtered_derived;threshold=200
  65. family=bag;stmt=comma_join;threshold=20000
  66. family=reuse;policy=notmat;refs=16;consumer=1%
  67. family=reuse;policy=mat;refs=8;consumer=100%
  68. family=bag;stmt=join_where;threshold=200
  69. family=bag;stmt=filtered_derived;threshold=20
  70. family=bag;stmt=join_where;threshold=20
  71. family=reuse;policy=notmat;refs=16;consumer=100%
  72. family=bag;stmt=cte_filtered_notmat;threshold=200
  73. family=bag;stmt=filtered_derived;threshold=20000
  74. family=bag;stmt=cte_filtered;threshold=100000
  75. family=reuse;policy=default;refs=2;consumer=1%
  76. family=bag;stmt=cte_joined_outer_notmat;threshold=100000
  77. family=reuse;policy=mat;refs=1;consumer=100%
  78. family=reuse;policy=default;refs=8;consumer=1%
  79. family=reuse;policy=default;refs=4;consumer=100%
  80. family=bag;stmt=join_on;threshold=180000
  81. family=bag;stmt=cte_joined_outer;threshold=2000
  82. family=reuse;policy=mat;refs=16;consumer=100%
  83. family=bag;stmt=comma_join;threshold=2000
  84. family=bag;stmt=join_where;threshold=100000
  85. family=bag;stmt=comma_join;threshold=200
  86. family=bag;stmt=cte_filtered_mat;threshold=100000
  87. family=bag;stmt=cte_filtered;threshold=2000
  88. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  89. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  90. family=bag;stmt=join_on;threshold=100000
  91. family=reuse;policy=mat;refs=4;consumer=100%;position=last

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
| family=bag;stmt=comma_join;threshold=100000 | j06_comma_join | valid | `4bb1561e58681728` |  |
| family=bag;stmt=comma_join;threshold=180000 | j06_comma_join | valid | `18e11c5d2e685630` |  |
| family=bag;stmt=comma_join;threshold=20 | j06_comma_join | valid | `61819ba7d0f5f952` |  |
| family=bag;stmt=comma_join;threshold=200 | j06_comma_join | valid | `41f23de3437828d5` |  |
| family=bag;stmt=comma_join;threshold=2000 | j06_comma_join | valid | `5573434208edb047` |  |
| family=bag;stmt=comma_join;threshold=20000 | j06_comma_join | valid | `09a028d269057a3e` |  |
| family=bag;stmt=cte_filtered;threshold=100000 | j04_cte_filtered | valid | `4bb1561e58681728` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=180000 | j04_cte_filtered | valid | `18e11c5d2e685630` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=20 | j04_cte_filtered | valid | `61819ba7d0f5f952` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=200 | j04_cte_filtered | valid | `41f23de3437828d5` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=2000 | j04_cte_filtered | valid | `5573434208edb047` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=20000 | j04_cte_filtered | valid | `09a028d269057a3e` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_mat;threshold=100000 | j04m_cte_filtered_mat | valid | `4b4ce4c83ae7fcc8` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=180000 | j04m_cte_filtered_mat | valid | `6a10a6fd7e521978` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=20 | j04m_cte_filtered_mat | valid | `d4ca4a9d1c0fd3e3` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=200 | j04m_cte_filtered_mat | valid | `a398c2b3383a151d` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=2000 | j04m_cte_filtered_mat | valid | `d4df6a0ef2f05082` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=20000 | j04m_cte_filtered_mat | valid | `5db776896b6f593f` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_notmat;threshold=100000 | j04n_cte_filtered_notmat | valid | `4bb1561e58681728` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=180000 | j04n_cte_filtered_notmat | valid | `18e11c5d2e685630` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=20 | j04n_cte_filtered_notmat | valid | `61819ba7d0f5f952` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=200 | j04n_cte_filtered_notmat | valid | `41f23de3437828d5` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=2000 | j04n_cte_filtered_notmat | valid | `5573434208edb047` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=20000 | j04n_cte_filtered_notmat | valid | `09a028d269057a3e` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=100000 | j05_cte_joined_outer | valid | `4bb1561e58681728` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=180000 | j05_cte_joined_outer | valid | `18e11c5d2e685630` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20 | j05_cte_joined_outer | valid | `61819ba7d0f5f952` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=200 | j05_cte_joined_outer | valid | `41f23de3437828d5` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=2000 | j05_cte_joined_outer | valid | `5573434208edb047` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20000 | j05_cte_joined_outer | valid | `09a028d269057a3e` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_mat;threshold=100000 | j05m_cte_joined_outer_mat | valid | `7f87bace4cf1da4c` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | j05m_cte_joined_outer_mat | valid | `7f87bace4cf1da4c` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | j05m_cte_joined_outer_mat | valid | `7f87bace4cf1da4c` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=200 | j05m_cte_joined_outer_mat | valid | `7f87bace4cf1da4c` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | j05m_cte_joined_outer_mat | valid | `7f87bace4cf1da4c` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20000 | j05m_cte_joined_outer_mat | valid | `7f87bace4cf1da4c` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_notmat;threshold=100000 | j05n_cte_joined_outer_notmat | valid | `4bb1561e58681728` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | j05n_cte_joined_outer_notmat | valid | `18e11c5d2e685630` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | j05n_cte_joined_outer_notmat | valid | `61819ba7d0f5f952` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=200 | j05n_cte_joined_outer_notmat | valid | `41f23de3437828d5` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | j05n_cte_joined_outer_notmat | valid | `5573434208edb047` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20000 | j05n_cte_joined_outer_notmat | valid | `09a028d269057a3e` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=filtered_derived;threshold=100000 | j03_filtered_derived | valid | `4bb1561e58681728` |  |
| family=bag;stmt=filtered_derived;threshold=180000 | j03_filtered_derived | valid | `18e11c5d2e685630` |  |
| family=bag;stmt=filtered_derived;threshold=20 | j03_filtered_derived | valid | `61819ba7d0f5f952` |  |
| family=bag;stmt=filtered_derived;threshold=200 | j03_filtered_derived | valid | `41f23de3437828d5` |  |
| family=bag;stmt=filtered_derived;threshold=2000 | j03_filtered_derived | valid | `5573434208edb047` |  |
| family=bag;stmt=filtered_derived;threshold=20000 | j03_filtered_derived | valid | `09a028d269057a3e` |  |
| family=bag;stmt=join_on;threshold=100000 | j02_join_on | valid | `4bb1561e58681728` |  |
| family=bag;stmt=join_on;threshold=180000 | j02_join_on | valid | `18e11c5d2e685630` |  |
| family=bag;stmt=join_on;threshold=20 | j02_join_on | valid | `61819ba7d0f5f952` |  |
| family=bag;stmt=join_on;threshold=200 | j02_join_on | valid | `41f23de3437828d5` |  |
| family=bag;stmt=join_on;threshold=2000 | j02_join_on | valid | `5573434208edb047` |  |
| family=bag;stmt=join_on;threshold=20000 | j02_join_on | valid | `09a028d269057a3e` |  |
| family=bag;stmt=join_where;threshold=100000 | j01_join_where | valid | `4bb1561e58681728` |  |
| family=bag;stmt=join_where;threshold=180000 | j01_join_where | valid | `18e11c5d2e685630` |  |
| family=bag;stmt=join_where;threshold=20 | j01_join_where | valid | `61819ba7d0f5f952` |  |
| family=bag;stmt=join_where;threshold=200 | j01_join_where | valid | `41f23de3437828d5` |  |
| family=bag;stmt=join_where;threshold=2000 | j01_join_where | valid | `5573434208edb047` |  |
| family=bag;stmt=join_where;threshold=20000 | j01_join_where | valid | `09a028d269057a3e` |  |
| family=reuse;policy=default;refs=16;consumer=1% | r_default_n16_c1 | valid | `b0c356323d30c425` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=16;consumer=100% | r_default_n16_c100 | valid | `b0c356323d30c425` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=1;consumer=1% | r_default_n1_c1 | valid | `549a07108ab01ba9` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=1;consumer=100% | r_default_n1_c100 | valid | `2ce8b85a7b60c88e` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=2;consumer=1% | r_default_n2_c1 | valid | `b307c41eb54cf452` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=2;consumer=100% | r_default_n2_c100 | valid | `b307c41eb54cf452` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=1% | r_default_n4_c1 | valid | `ccfe1af393deae98` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=100% | r_default_n4_c100 | valid | `ccfe1af393deae98` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=8;consumer=1% | r_default_n8_c1 | valid | `c8c73778097931dc` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=8;consumer=100% | r_default_n8_c100 | valid | `c8c73778097931dc` | policy=default; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=1% | r_mat_n16_c1 | valid | `b0c356323d30c425` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=100% | r_mat_n16_c100 | valid | `b0c356323d30c425` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=1% | r_mat_n1_c1 | valid | `eaa7d6058aaa5e03` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=100% | r_mat_n1_c100 | valid | `eaa7d6058aaa5e03` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=2;consumer=1% | r_mat_n2_c1 | valid | `b307c41eb54cf452` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=2;consumer=100% | r_mat_n2_c100 | valid | `b307c41eb54cf452` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=1% | r_mat_n4_c1 | valid | `ccfe1af393deae98` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=100% | r_mat_n4_c100 | valid | `ccfe1af393deae98` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=100%;position=last | r_mat_n4_c100 | valid | `ccfe1af393deae98` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=8;consumer=1% | r_mat_n8_c1 | valid | `c8c73778097931dc` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=8;consumer=100% | r_mat_n8_c100 | valid | `c8c73778097931dc` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=notmat;refs=16;consumer=1% | r_notmat_n16_c1 | valid | `1c38b58d8cec08d3` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `effce31a81556487` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `549a07108ab01ba9` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `2ce8b85a7b60c88e` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=2;consumer=1% | r_notmat_n2_c1 | valid | `28cc913d71d0d72b` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=2;consumer=100% | r_notmat_n2_c100 | valid | `2043b765d5001655` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `9f850e53560baafb` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `cc53feafb616350b` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=8;consumer=1% | r_notmat_n8_c1 | valid | `477f7fbaf1c627cb` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=8;consumer=100% | r_notmat_n8_c100 | valid | `714b97b0f1e43d36` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, VERBOSE, BUFFERS, FORMAT TEXT` supported on postgres
- one host, single client: no independent hardware, real-network or equal-total-budget comparison
- no concurrent writers: the isolated read contract only; a matched blended read/write phase is a recorded gap
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
