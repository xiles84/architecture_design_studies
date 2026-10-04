# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 91 cells, 91 valid, engine `yugabyte`, topology `single`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-domain2-yb1-20261004`
- run tag: `run/study-01/select-join-domain2-yb1-20261004`
- environment: `host-zenbook-ux5406sa`
- repo commit: `445fbc1646d58074c45cbaea57a83234132adc34`
- engine version: `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit`
- database engine image: `docker.io/yugabytedb/yugabyte:2025.2.6.0-b111` @ `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`
- benchmark client image: `localhost/charitybench-select:select-join-domain2-yb1-20261004` @ `79bc753ff7e57e102982a323939c43e68e9b3664887683d3a5d5bf072650f4c1`
- inputs digest: `f03b1ad718dce4985f6849dd8f4cf33687d09daea675106ed2992eb3a1aecb2b`
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
| 20 (0.01%) | 55.8052 | 55.1772 | 55.6128 | 55.8798 | 56.3830 | 56.1749 | 56.6447 | 88.0918 | 54.5287 | 55.2845 |
| 200 (0.1%) | 56.8781 | 55.7799 | 54.1900 | 56.0599 | 54.8696 | 56.3623 | 55.8737 | 92.6668 | 57.2431 | 56.4192 |
| 2000 (1%) | 57.1350 | 56.8928 | 58.2196 | 58.4553 | 56.7935 | 55.3503 | 59.5656 | 83.0722 | 55.3730 | 54.6947 |
| 20000 (10%) | 66.8531 | 70.6190 | 68.1145 | 68.0566 | 73.1833 | 67.7933 | 69.9399 | 92.6200 | 72.6335 | 71.1162 |
| 100000 (50%) | 95.1290 | 87.1198 | 89.4769 | 88.0010 | 93.1772 | 88.1750 | 92.6738 | 93.7120 | 95.5195 | 92.2789 |
| 180000 (90%) | 109.7212 | 114.2475 | 107.2297 | 112.0193 | 118.5577 | 111.0952 | 108.6631 | 107.1316 | 110.6173 | 108.6497 |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | join_on vs join_where | 55.8052 | 55.1772 | 0.9887 | false |
| 20 | filtered_derived vs join_where | 55.8052 | 55.6128 | 0.9966 | false |
| 20 | cte_filtered vs join_where | 55.8052 | 55.8798 | 1.0013 | false |
| 20 | cte_filtered_mat vs join_where | 55.8052 | 56.3830 | 1.0104 | false |
| 20 | cte_filtered_notmat vs join_where | 55.8052 | 56.1749 | 1.0066 | false |
| 20 | cte_joined_outer vs join_where | 55.8052 | 56.6447 | 1.0150 | false |
| 20 | cte_joined_outer_mat vs join_where | 55.8052 | 88.0918 | 1.5786 | false |
| 20 | cte_joined_outer_notmat vs join_where | 55.8052 | 54.5287 | 0.9771 | false |
| 20 | comma_join vs join_where | 55.8052 | 55.2845 | 0.9907 | false |
| 200 | join_on vs join_where | 56.8781 | 55.7799 | 0.9807 | false |
| 200 | filtered_derived vs join_where | 56.8781 | 54.1900 | 0.9527 | false |
| 200 | cte_filtered vs join_where | 56.8781 | 56.0599 | 0.9856 | false |
| 200 | cte_filtered_mat vs join_where | 56.8781 | 54.8696 | 0.9647 | false |
| 200 | cte_filtered_notmat vs join_where | 56.8781 | 56.3623 | 0.9909 | false |
| 200 | cte_joined_outer vs join_where | 56.8781 | 55.8737 | 0.9823 | false |
| 200 | cte_joined_outer_mat vs join_where | 56.8781 | 92.6668 | 1.6292 | false |
| 200 | cte_joined_outer_notmat vs join_where | 56.8781 | 57.2431 | 1.0064 | false |
| 200 | comma_join vs join_where | 56.8781 | 56.4192 | 0.9919 | false |
| 2000 | join_on vs join_where | 57.1350 | 56.8928 | 0.9958 | false |
| 2000 | filtered_derived vs join_where | 57.1350 | 58.2196 | 1.0190 | false |
| 2000 | cte_filtered vs join_where | 57.1350 | 58.4553 | 1.0231 | false |
| 2000 | cte_filtered_mat vs join_where | 57.1350 | 56.7935 | 0.9940 | false |
| 2000 | cte_filtered_notmat vs join_where | 57.1350 | 55.3503 | 0.9688 | false |
| 2000 | cte_joined_outer vs join_where | 57.1350 | 59.5656 | 1.0425 | false |
| 2000 | cte_joined_outer_mat vs join_where | 57.1350 | 83.0722 | 1.4540 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 57.1350 | 55.3730 | 0.9692 | false |
| 2000 | comma_join vs join_where | 57.1350 | 54.6947 | 0.9573 | false |
| 20000 | join_on vs join_where | 66.8531 | 70.6190 | 1.0563 | false |
| 20000 | filtered_derived vs join_where | 66.8531 | 68.1145 | 1.0189 | false |
| 20000 | cte_filtered vs join_where | 66.8531 | 68.0566 | 1.0180 | false |
| 20000 | cte_filtered_mat vs join_where | 66.8531 | 73.1833 | 1.0947 | false |
| 20000 | cte_filtered_notmat vs join_where | 66.8531 | 67.7933 | 1.0141 | false |
| 20000 | cte_joined_outer vs join_where | 66.8531 | 69.9399 | 1.0462 | false |
| 20000 | cte_joined_outer_mat vs join_where | 66.8531 | 92.6200 | 1.3854 | false |
| 20000 | cte_joined_outer_notmat vs join_where | 66.8531 | 72.6335 | 1.0865 | false |
| 20000 | comma_join vs join_where | 66.8531 | 71.1162 | 1.0638 | false |
| 100000 | join_on vs join_where | 95.1290 | 87.1198 | 0.9158 | false |
| 100000 | filtered_derived vs join_where | 95.1290 | 89.4769 | 0.9406 | false |
| 100000 | cte_filtered vs join_where | 95.1290 | 88.0010 | 0.9251 | false |
| 100000 | cte_filtered_mat vs join_where | 95.1290 | 93.1772 | 0.9795 | false |
| 100000 | cte_filtered_notmat vs join_where | 95.1290 | 88.1750 | 0.9269 | false |
| 100000 | cte_joined_outer vs join_where | 95.1290 | 92.6738 | 0.9742 | false |
| 100000 | cte_joined_outer_mat vs join_where | 95.1290 | 93.7120 | 0.9851 | false |
| 100000 | cte_joined_outer_notmat vs join_where | 95.1290 | 95.5195 | 1.0041 | false |
| 100000 | comma_join vs join_where | 95.1290 | 92.2789 | 0.9700 | false |
| 180000 | join_on vs join_where | 109.7212 | 114.2475 | 1.0413 | false |
| 180000 | filtered_derived vs join_where | 109.7212 | 107.2297 | 0.9773 | false |
| 180000 | cte_filtered vs join_where | 109.7212 | 112.0193 | 1.0209 | false |
| 180000 | cte_filtered_mat vs join_where | 109.7212 | 118.5577 | 1.0805 | false |
| 180000 | cte_filtered_notmat vs join_where | 109.7212 | 111.0952 | 1.0125 | false |
| 180000 | cte_joined_outer vs join_where | 109.7212 | 108.6631 | 0.9904 | false |
| 180000 | cte_joined_outer_mat vs join_where | 109.7212 | 107.1316 | 0.9764 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 109.7212 | 110.6173 | 1.0082 | false |
| 180000 | comma_join vs join_where | 109.7212 | 108.6497 | 0.9902 | false |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_filtered_mat vs cte_filtered | 55.8798 | 56.3830 | 1.0090 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_filtered_notmat vs cte_filtered | 55.8798 | 56.1749 | 1.0053 | false | policy=not_materialized; observed_cte_scan=false |
| 20 | cte_joined_outer_mat vs cte_joined_outer | 56.6447 | 88.0918 | 1.5552 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 56.6447 | 54.5287 | 0.9626 | false | policy=not_materialized; observed_cte_scan=false |
| 200 | cte_filtered_mat vs cte_filtered | 56.0599 | 54.8696 | 0.9788 | false | policy=materialized; observed_cte_scan=true |
| 200 | cte_filtered_notmat vs cte_filtered | 56.0599 | 56.3623 | 1.0054 | false | policy=not_materialized; observed_cte_scan=false |
| 200 | cte_joined_outer_mat vs cte_joined_outer | 55.8737 | 92.6668 | 1.6585 | false | policy=materialized; observed_cte_scan=true |
| 200 | cte_joined_outer_notmat vs cte_joined_outer | 55.8737 | 57.2431 | 1.0245 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_filtered_mat vs cte_filtered | 58.4553 | 56.7935 | 0.9716 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_filtered_notmat vs cte_filtered | 58.4553 | 55.3503 | 0.9469 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 59.5656 | 83.0722 | 1.3946 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 59.5656 | 55.3730 | 0.9296 | false | policy=not_materialized; observed_cte_scan=false |
| 20000 | cte_filtered_mat vs cte_filtered | 68.0566 | 73.1833 | 1.0753 | false | policy=materialized; observed_cte_scan=true |
| 20000 | cte_filtered_notmat vs cte_filtered | 68.0566 | 67.7933 | 0.9961 | false | policy=not_materialized; observed_cte_scan=false |
| 20000 | cte_joined_outer_mat vs cte_joined_outer | 69.9399 | 92.6200 | 1.3243 | false | policy=materialized; observed_cte_scan=true |
| 20000 | cte_joined_outer_notmat vs cte_joined_outer | 69.9399 | 72.6335 | 1.0385 | false | policy=not_materialized; observed_cte_scan=false |
| 100000 | cte_filtered_mat vs cte_filtered | 88.0010 | 93.1772 | 1.0588 | false | policy=materialized; observed_cte_scan=true |
| 100000 | cte_filtered_notmat vs cte_filtered | 88.0010 | 88.1750 | 1.0020 | false | policy=not_materialized; observed_cte_scan=false |
| 100000 | cte_joined_outer_mat vs cte_joined_outer | 92.6738 | 93.7120 | 1.0112 | false | policy=materialized; observed_cte_scan=true |
| 100000 | cte_joined_outer_notmat vs cte_joined_outer | 92.6738 | 95.5195 | 1.0307 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_filtered_mat vs cte_filtered | 112.0193 | 118.5577 | 1.0584 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_filtered_notmat vs cte_filtered | 112.0193 | 111.0952 | 0.9918 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 108.6631 | 107.1316 | 0.9859 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 108.6631 | 110.6173 | 1.0180 | false | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 57.7764 | 1 | `ad5744afcbe7b6b6` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 62.0809 | 1 | `d8ebcb11b976ac36` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | 55.9737 | 2 | `254b1ba7d92b8ae0` | policy=default; observed_cte_scan=true | valid |
| default | 2 | 100% | 56.6995 | 2 | `8ed78cead8319934` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 1% | 56.7135 | 4 | `f2b9388c72edd0e1` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 58.8414 | 4 | `db4cb50a49872f43` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | 54.7866 | 8 | `16266d61c667863b` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 100% | 57.6598 | 8 | `88059c8f08680832` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 1% | 56.8527 | 16 | `59d871ffa0680a2f` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 59.4180 | 16 | `481da00ea7550bd2` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 57.1153 | 1 | `99d110cf51c37192` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 55.0072 | 1 | `d19a3c2e0b10bd83` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | 58.4316 | 2 | `d866cf131fd6f5f9` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 100% | 59.9335 | 2 | `95b46adef2209ea4` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 1% | 55.5403 | 4 | `889edb0ee90fce9c` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 59.0390 | 4 | `2bda37e1b03e45a2` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | 57.7580 | 8 | `8d36c125916eff8b` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 100% | 57.6682 | 8 | `1c33b2ce3213a7ac` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 1% | 58.5167 | 16 | `9e121a22d7f830b7` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 54.5271 | 16 | `587b4cb383830075` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 59.2344 | 1 | `ec810d1ab9da57d8` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 61.8831 | 1 | `599aae1e435bf4e4` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | 146.4694 | 2 | `c5c0f68bdd0b9b7a` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 100% | 112.6876 | 2 | `69fef1feaa7e3580` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 1% | 141.6230 | 4 | `e8da14cf5750a312` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 161.7322 | 4 | `12b5638e569a45ec` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | 305.7852 | 8 | `9fc14f17163c7d7a` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 100% | 316.3024 | 8 | `6aa0ea938505e779` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 1% | 554.5876 | 16 | `bafa437c0105c768` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 562.1632 | 16 | `4c882048f62cf8fd` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 0.9886 | false |
| notmat | 1 | 1% | 1.0252 | false |
| mat | 1 | 100% | 0.8861 | false |
| notmat | 1 | 100% | 0.9968 | false |
| mat | 2 | 1% | 1.0439 | false |
| notmat | 2 | 1% | 2.6168 | false |
| mat | 2 | 100% | 1.0570 | false |
| notmat | 2 | 100% | 1.9875 | false |
| mat | 4 | 1% | 0.9793 | false |
| notmat | 4 | 1% | 2.4972 | false |
| mat | 4 | 100% | 1.0034 | false |
| notmat | 4 | 100% | 2.7486 | false |
| mat | 8 | 1% | 1.0542 | false |
| notmat | 8 | 1% | 5.5814 | false |
| mat | 8 | 100% | 1.0001 | false |
| notmat | 8 | 100% | 5.4857 | false |
| mat | 16 | 1% | 1.0293 | false |
| notmat | 16 | 1% | 9.7548 | false |
| mat | 16 | 100% | 0.9177 | false |
| notmat | 16 | 100% | 9.4612 | false |

## Repeated-design position control

| first median ms | last median ms | ratio (last/first) | same plan |
|---|---|---|---|
| 59.565603 | 58.188388 | 0.9769 | false |

## Latency samples and derived throughput

_The reader is sequential, so `queries/s` is 1000/median_ms, a derived rate and not a concurrent-capacity test. min/max are the five sample extremes._

| axis | status | min ms | median ms | max ms | samples | derived queries/s |
|---|---|---|---|---|---|---|
| family=bag;stmt=comma_join;threshold=100000 | valid | 85.9662 | 92.2789 | 94.7647 | 5 | 10.8 |
| family=bag;stmt=comma_join;threshold=180000 | valid | 102.0204 | 108.6497 | 119.4531 | 5 | 9.2 |
| family=bag;stmt=comma_join;threshold=20 | valid | 53.2034 | 55.2845 | 60.6583 | 5 | 18.1 |
| family=bag;stmt=comma_join;threshold=200 | valid | 53.0999 | 56.4192 | 57.9696 | 5 | 17.7 |
| family=bag;stmt=comma_join;threshold=2000 | valid | 53.4220 | 54.6947 | 57.7838 | 5 | 18.3 |
| family=bag;stmt=comma_join;threshold=20000 | valid | 69.8752 | 71.1162 | 75.2847 | 5 | 14.1 |
| family=bag;stmt=cte_filtered;threshold=100000 | valid | 85.1389 | 88.0010 | 92.3975 | 5 | 11.4 |
| family=bag;stmt=cte_filtered;threshold=180000 | valid | 104.8037 | 112.0193 | 126.0665 | 5 | 8.9 |
| family=bag;stmt=cte_filtered;threshold=20 | valid | 53.2405 | 55.8798 | 56.5761 | 5 | 17.9 |
| family=bag;stmt=cte_filtered;threshold=200 | valid | 53.3387 | 56.0599 | 59.8090 | 5 | 17.8 |
| family=bag;stmt=cte_filtered;threshold=2000 | valid | 57.9902 | 58.4553 | 63.5213 | 5 | 17.1 |
| family=bag;stmt=cte_filtered;threshold=20000 | valid | 67.1458 | 68.0566 | 70.7480 | 5 | 14.7 |
| family=bag;stmt=cte_filtered_mat;threshold=100000 | valid | 92.1714 | 93.1772 | 94.9773 | 5 | 10.7 |
| family=bag;stmt=cte_filtered_mat;threshold=180000 | valid | 113.5827 | 118.5577 | 120.0571 | 5 | 8.4 |
| family=bag;stmt=cte_filtered_mat;threshold=20 | valid | 55.4122 | 56.3830 | 59.9619 | 5 | 17.7 |
| family=bag;stmt=cte_filtered_mat;threshold=200 | valid | 52.7480 | 54.8696 | 57.4732 | 5 | 18.2 |
| family=bag;stmt=cte_filtered_mat;threshold=2000 | valid | 55.6420 | 56.7935 | 57.6319 | 5 | 17.6 |
| family=bag;stmt=cte_filtered_mat;threshold=20000 | valid | 68.7326 | 73.1833 | 92.3822 | 5 | 13.7 |
| family=bag;stmt=cte_filtered_notmat;threshold=100000 | valid | 85.9273 | 88.1750 | 91.7281 | 5 | 11.3 |
| family=bag;stmt=cte_filtered_notmat;threshold=180000 | valid | 107.2232 | 111.0952 | 123.2243 | 5 | 9.0 |
| family=bag;stmt=cte_filtered_notmat;threshold=20 | valid | 52.8496 | 56.1749 | 61.1098 | 5 | 17.8 |
| family=bag;stmt=cte_filtered_notmat;threshold=200 | valid | 54.3285 | 56.3623 | 59.9419 | 5 | 17.7 |
| family=bag;stmt=cte_filtered_notmat;threshold=2000 | valid | 53.5982 | 55.3503 | 57.0987 | 5 | 18.1 |
| family=bag;stmt=cte_filtered_notmat;threshold=20000 | valid | 66.1108 | 67.7933 | 70.1547 | 5 | 14.8 |
| family=bag;stmt=cte_joined_outer;threshold=100000 | valid | 89.8678 | 92.6738 | 93.3582 | 5 | 10.8 |
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 107.7674 | 108.6631 | 110.5682 | 5 | 9.2 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 52.8696 | 56.6447 | 58.0285 | 5 | 17.7 |
| family=bag;stmt=cte_joined_outer;threshold=200 | valid | 55.3350 | 55.8737 | 59.0328 | 5 | 17.9 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 58.6344 | 59.5656 | 71.8006 | 5 | 16.8 |
| family=bag;stmt=cte_joined_outer;threshold=20000 | valid | 68.9636 | 69.9399 | 72.2647 | 5 | 14.3 |
| family=bag;stmt=cte_joined_outer;threshold=2000;position=last | valid | 58.1317 | 58.1884 | 62.0864 | 5 | 17.2 |
| family=bag;stmt=cte_joined_outer_mat;threshold=100000 | valid | 91.3984 | 93.7120 | 95.1624 | 5 | 10.7 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 105.9687 | 107.1316 | 108.5162 | 5 | 9.3 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 85.9469 | 88.0918 | 90.9667 | 5 | 11.4 |
| family=bag;stmt=cte_joined_outer_mat;threshold=200 | valid | 88.8494 | 92.6668 | 95.0705 | 5 | 10.8 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 79.9794 | 83.0722 | 89.2164 | 5 | 12.0 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20000 | valid | 87.3642 | 92.6200 | 94.6228 | 5 | 10.8 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=100000 | valid | 93.6253 | 95.5195 | 101.9880 | 5 | 10.5 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 105.1483 | 110.6173 | 113.2052 | 5 | 9.0 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 54.2881 | 54.5287 | 57.7546 | 5 | 18.3 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=200 | valid | 55.2436 | 57.2431 | 60.7324 | 5 | 17.5 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 55.1767 | 55.3730 | 59.0035 | 5 | 18.1 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20000 | valid | 71.7606 | 72.6335 | 74.1044 | 5 | 13.8 |
| family=bag;stmt=filtered_derived;threshold=100000 | valid | 86.7910 | 89.4769 | 90.3327 | 5 | 11.2 |
| family=bag;stmt=filtered_derived;threshold=180000 | valid | 104.2599 | 107.2297 | 119.0990 | 5 | 9.3 |
| family=bag;stmt=filtered_derived;threshold=20 | valid | 54.4031 | 55.6128 | 58.4491 | 5 | 18.0 |
| family=bag;stmt=filtered_derived;threshold=200 | valid | 51.9516 | 54.1900 | 58.8225 | 5 | 18.5 |
| family=bag;stmt=filtered_derived;threshold=2000 | valid | 54.7323 | 58.2196 | 60.0447 | 5 | 17.2 |
| family=bag;stmt=filtered_derived;threshold=20000 | valid | 67.6692 | 68.1145 | 71.1186 | 5 | 14.7 |
| family=bag;stmt=join_on;threshold=100000 | valid | 86.1533 | 87.1198 | 92.5386 | 5 | 11.5 |
| family=bag;stmt=join_on;threshold=180000 | valid | 110.7270 | 114.2475 | 117.6833 | 5 | 8.8 |
| family=bag;stmt=join_on;threshold=20 | valid | 54.9057 | 55.1772 | 57.0472 | 5 | 18.1 |
| family=bag;stmt=join_on;threshold=200 | valid | 53.2255 | 55.7799 | 57.4686 | 5 | 17.9 |
| family=bag;stmt=join_on;threshold=2000 | valid | 54.4505 | 56.8928 | 59.1472 | 5 | 17.6 |
| family=bag;stmt=join_on;threshold=20000 | valid | 69.3819 | 70.6190 | 71.8423 | 5 | 14.2 |
| family=bag;stmt=join_where;threshold=100000 | valid | 91.3726 | 95.1290 | 97.4432 | 5 | 10.5 |
| family=bag;stmt=join_where;threshold=180000 | valid | 103.1524 | 109.7212 | 112.9257 | 5 | 9.1 |
| family=bag;stmt=join_where;threshold=20 | valid | 53.9384 | 55.8052 | 58.9695 | 5 | 17.9 |
| family=bag;stmt=join_where;threshold=200 | valid | 55.0043 | 56.8781 | 59.4340 | 5 | 17.6 |
| family=bag;stmt=join_where;threshold=2000 | valid | 56.1826 | 57.1350 | 58.3650 | 5 | 17.5 |
| family=bag;stmt=join_where;threshold=20000 | valid | 64.9751 | 66.8531 | 71.2560 | 5 | 15.0 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 55.2413 | 56.8527 | 59.6047 | 5 | 17.6 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 58.2148 | 59.4180 | 63.3530 | 5 | 16.8 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 54.4914 | 57.7764 | 60.7214 | 5 | 17.3 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 60.6359 | 62.0809 | 78.8015 | 5 | 16.1 |
| family=reuse;policy=default;refs=2;consumer=1% | valid | 53.6642 | 55.9737 | 58.1410 | 5 | 17.9 |
| family=reuse;policy=default;refs=2;consumer=100% | valid | 53.1852 | 56.6995 | 59.7900 | 5 | 17.6 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 55.2874 | 56.7135 | 57.2598 | 5 | 17.6 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 57.0365 | 58.8414 | 59.8779 | 5 | 17.0 |
| family=reuse;policy=default;refs=8;consumer=1% | valid | 53.7241 | 54.7866 | 58.5666 | 5 | 18.3 |
| family=reuse;policy=default;refs=8;consumer=100% | valid | 54.4890 | 57.6598 | 58.2040 | 5 | 17.3 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 56.4055 | 58.5167 | 59.4087 | 5 | 17.1 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 54.2299 | 54.5271 | 61.5655 | 5 | 18.3 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 55.9680 | 57.1153 | 60.6204 | 5 | 17.5 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 53.1071 | 55.0072 | 59.7626 | 5 | 18.2 |
| family=reuse;policy=mat;refs=2;consumer=1% | valid | 57.5242 | 58.4316 | 62.8201 | 5 | 17.1 |
| family=reuse;policy=mat;refs=2;consumer=100% | valid | 57.2815 | 59.9335 | 62.3583 | 5 | 16.7 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 52.5335 | 55.5403 | 56.7424 | 5 | 18.0 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 58.4164 | 59.0390 | 63.5226 | 5 | 16.9 |
| family=reuse;policy=mat;refs=8;consumer=1% | valid | 55.6596 | 57.7580 | 59.1886 | 5 | 17.3 |
| family=reuse;policy=mat;refs=8;consumer=100% | valid | 56.0707 | 57.6682 | 58.8884 | 5 | 17.3 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 548.8117 | 554.5876 | 569.3471 | 5 | 1.8 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 557.2699 | 562.1632 | 579.5473 | 5 | 1.8 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 57.2188 | 59.2344 | 64.4785 | 5 | 16.9 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 59.8797 | 61.8831 | 64.9903 | 5 | 16.2 |
| family=reuse;policy=notmat;refs=2;consumer=1% | valid | 134.9220 | 146.4694 | 167.1140 | 5 | 6.8 |
| family=reuse;policy=notmat;refs=2;consumer=100% | valid | 106.0655 | 112.6876 | 120.5456 | 5 | 8.9 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 132.1269 | 141.6230 | 160.6559 | 5 | 7.1 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 138.5096 | 161.7322 | 173.9913 | 5 | 6.2 |
| family=reuse;policy=notmat;refs=8;consumer=1% | valid | 294.4376 | 305.7852 | 319.0154 | 5 | 3.3 |
| family=reuse;policy=notmat;refs=8;consumer=100% | valid | 274.8299 | 316.3024 | 338.1887 | 5 | 3.2 |

## Matched blended workload

- phase: `isolated`; target offered rate: 10 updates/s
- offered 0, committed 0, failed 0 over 0.000 s
- offered rate 0.000/s, completed rate 0.000/s
- per key: no keys touched
- ledger: 0 keys reconciled, 0 mismatch

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42850290`, 5 interleaved rounds
  1. family=bag;stmt=cte_joined_outer;threshold=2000
  2. family=bag;stmt=comma_join;threshold=20000
  3. family=reuse;policy=notmat;refs=4;consumer=1%
  4. family=reuse;policy=notmat;refs=2;consumer=1%
  5. family=bag;stmt=cte_joined_outer_mat;threshold=200
  6. family=bag;stmt=cte_joined_outer_notmat;threshold=100000
  7. family=reuse;policy=mat;refs=2;consumer=100%
  8. family=reuse;policy=notmat;refs=8;consumer=1%
  9. family=bag;stmt=cte_joined_outer_notmat;threshold=20000
  10. family=reuse;policy=mat;refs=4;consumer=100%
  11. family=bag;stmt=cte_filtered;threshold=2000
  12. family=reuse;policy=default;refs=1;consumer=100%
  13. family=bag;stmt=cte_filtered_mat;threshold=20000
  14. family=bag;stmt=join_where;threshold=100000
  15. family=bag;stmt=cte_joined_outer;threshold=100000
  16. family=bag;stmt=comma_join;threshold=100000
  17. family=reuse;policy=notmat;refs=8;consumer=100%
  18. family=reuse;policy=notmat;refs=1;consumer=1%
  19. family=bag;stmt=cte_joined_outer;threshold=20000
  20. family=reuse;policy=default;refs=16;consumer=100%
  21. family=reuse;policy=mat;refs=2;consumer=1%
  22. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  23. family=bag;stmt=cte_joined_outer_notmat;threshold=200
  24. family=reuse;policy=default;refs=16;consumer=1%
  25. family=reuse;policy=default;refs=4;consumer=1%
  26. family=reuse;policy=mat;refs=8;consumer=100%
  27. family=bag;stmt=cte_filtered_mat;threshold=100000
  28. family=bag;stmt=cte_filtered_mat;threshold=2000
  29. family=bag;stmt=comma_join;threshold=20
  30. family=reuse;policy=notmat;refs=2;consumer=100%
  31. family=bag;stmt=cte_filtered_mat;threshold=20
  32. family=bag;stmt=cte_joined_outer;threshold=200
  33. family=bag;stmt=cte_filtered_mat;threshold=180000
  34. family=bag;stmt=cte_joined_outer_mat;threshold=20000
  35. family=bag;stmt=cte_filtered_notmat;threshold=180000
  36. family=bag;stmt=join_on;threshold=2000
  37. family=bag;stmt=cte_filtered_notmat;threshold=20000
  38. family=reuse;policy=default;refs=8;consumer=100%
  39. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  40. family=reuse;policy=default;refs=2;consumer=100%
  41. family=reuse;policy=default;refs=8;consumer=1%
  42. family=bag;stmt=filtered_derived;threshold=20
  43. family=bag;stmt=comma_join;threshold=2000
  44. family=bag;stmt=comma_join;threshold=200
  45. family=bag;stmt=cte_filtered;threshold=180000
  46. family=bag;stmt=filtered_derived;threshold=180000
  47. family=bag;stmt=join_where;threshold=20
  48. family=reuse;policy=mat;refs=16;consumer=100%
  49. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  50. family=bag;stmt=cte_filtered;threshold=100000
  51. family=bag;stmt=join_where;threshold=180000
  52. family=bag;stmt=filtered_derived;threshold=200
  53. family=bag;stmt=cte_filtered_notmat;threshold=2000
  54. family=reuse;policy=mat;refs=4;consumer=1%
  55. family=bag;stmt=filtered_derived;threshold=100000
  56. family=bag;stmt=join_on;threshold=200
  57. family=reuse;policy=default;refs=1;consumer=1%
  58. family=reuse;policy=notmat;refs=4;consumer=100%
  59. family=reuse;policy=notmat;refs=16;consumer=1%
  60. family=reuse;policy=notmat;refs=1;consumer=100%
  61. family=reuse;policy=mat;refs=1;consumer=1%
  62. family=reuse;policy=default;refs=4;consumer=100%
  63. family=bag;stmt=cte_joined_outer_mat;threshold=20
  64. family=bag;stmt=join_on;threshold=20000
  65. family=reuse;policy=mat;refs=16;consumer=1%
  66. family=bag;stmt=filtered_derived;threshold=2000
  67. family=bag;stmt=join_where;threshold=200
  68. family=bag;stmt=join_on;threshold=180000
  69. family=bag;stmt=filtered_derived;threshold=20000
  70. family=bag;stmt=cte_filtered;threshold=20
  71. family=reuse;policy=default;refs=2;consumer=1%
  72. family=bag;stmt=cte_joined_outer;threshold=180000
  73. family=bag;stmt=cte_joined_outer;threshold=20
  74. family=bag;stmt=join_on;threshold=20
  75. family=bag;stmt=cte_joined_outer_mat;threshold=100000
  76. family=bag;stmt=cte_filtered_notmat;threshold=20
  77. family=bag;stmt=cte_filtered_notmat;threshold=100000
  78. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  79. family=bag;stmt=cte_filtered_notmat;threshold=200
  80. family=bag;stmt=join_where;threshold=20000
  81. family=bag;stmt=cte_filtered_mat;threshold=200
  82. family=reuse;policy=mat;refs=8;consumer=1%
  83. family=bag;stmt=comma_join;threshold=180000
  84. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  85. family=reuse;policy=mat;refs=1;consumer=100%
  86. family=bag;stmt=cte_filtered;threshold=200
  87. family=bag;stmt=join_on;threshold=100000
  88. family=bag;stmt=join_where;threshold=2000
  89. family=bag;stmt=cte_filtered;threshold=20000
  90. family=reuse;policy=notmat;refs=16;consumer=100%
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
| family=bag;stmt=comma_join;threshold=100000 | j06_comma_join | valid | `34abfc88c4ec4a9a` |  |
| family=bag;stmt=comma_join;threshold=180000 | j06_comma_join | valid | `29b8810be25fa22f` |  |
| family=bag;stmt=comma_join;threshold=20 | j06_comma_join | valid | `bdf4635e7fd62f25` |  |
| family=bag;stmt=comma_join;threshold=200 | j06_comma_join | valid | `a5eabfce94df2fa7` |  |
| family=bag;stmt=comma_join;threshold=2000 | j06_comma_join | valid | `660d749c0584479b` |  |
| family=bag;stmt=comma_join;threshold=20000 | j06_comma_join | valid | `f021353d56fa2899` |  |
| family=bag;stmt=cte_filtered;threshold=100000 | j04_cte_filtered | valid | `56dd7a4c1606cddc` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=180000 | j04_cte_filtered | valid | `75341b4fcfacc1e1` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=20 | j04_cte_filtered | valid | `9052bdc28008704f` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=200 | j04_cte_filtered | valid | `75c17fd0222f6a0d` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=2000 | j04_cte_filtered | valid | `e9ccc49c35bcf451` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=20000 | j04_cte_filtered | valid | `00bf6e39e20ea596` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_mat;threshold=100000 | j04m_cte_filtered_mat | valid | `173d84c26175a0a4` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=180000 | j04m_cte_filtered_mat | valid | `dd3ebf5d283a785f` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=20 | j04m_cte_filtered_mat | valid | `6163ad9c56303d63` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=200 | j04m_cte_filtered_mat | valid | `8886922eb2fc04a3` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=2000 | j04m_cte_filtered_mat | valid | `66bac0921d6112a7` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=20000 | j04m_cte_filtered_mat | valid | `fa53a7bdf37978bf` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_notmat;threshold=100000 | j04n_cte_filtered_notmat | valid | `8a0a736c0f6a967d` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=180000 | j04n_cte_filtered_notmat | valid | `8b8a6cb9aad836eb` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=20 | j04n_cte_filtered_notmat | valid | `ce1b089c12d29807` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=200 | j04n_cte_filtered_notmat | valid | `a010976135a7af03` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=2000 | j04n_cte_filtered_notmat | valid | `9e3270fd223138a7` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=20000 | j04n_cte_filtered_notmat | valid | `42c25fb74286a04e` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=100000 | j05_cte_joined_outer | valid | `a5f1045c820ea815` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=180000 | j05_cte_joined_outer | valid | `cbfbc2718351efb9` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20 | j05_cte_joined_outer | valid | `4e3af29ab7a08032` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=200 | j05_cte_joined_outer | valid | `1d90a96597263856` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=2000 | j05_cte_joined_outer | valid | `7c0b3a1cac1a374d` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20000 | j05_cte_joined_outer | valid | `f00100f5c73af0a9` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=2000;position=last | j05_cte_joined_outer | valid | `bfe3cec5cd75121a` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_mat;threshold=100000 | j05m_cte_joined_outer_mat | valid | `ab55dea5efeb8dcd` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | j05m_cte_joined_outer_mat | valid | `1c1289ac0983d2f3` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | j05m_cte_joined_outer_mat | valid | `6638c64d7eef68a1` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=200 | j05m_cte_joined_outer_mat | valid | `49c475ca703687e5` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | j05m_cte_joined_outer_mat | valid | `a83142c128a8d4f5` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20000 | j05m_cte_joined_outer_mat | valid | `a49a4cb978183128` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_notmat;threshold=100000 | j05n_cte_joined_outer_notmat | valid | `d1acc8804ef0325c` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | j05n_cte_joined_outer_notmat | valid | `603613ae8090c543` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | j05n_cte_joined_outer_notmat | valid | `998191b852c2fdfd` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=200 | j05n_cte_joined_outer_notmat | valid | `99b7e4e74f7d4604` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | j05n_cte_joined_outer_notmat | valid | `fd99b5b214a8e4ec` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20000 | j05n_cte_joined_outer_notmat | valid | `3f851204349a0818` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=filtered_derived;threshold=100000 | j03_filtered_derived | valid | `16e4dab2e376922d` |  |
| family=bag;stmt=filtered_derived;threshold=180000 | j03_filtered_derived | valid | `cd637369a0bda1b2` |  |
| family=bag;stmt=filtered_derived;threshold=20 | j03_filtered_derived | valid | `d328633a128c5689` |  |
| family=bag;stmt=filtered_derived;threshold=200 | j03_filtered_derived | valid | `d7f1ca17b59092ec` |  |
| family=bag;stmt=filtered_derived;threshold=2000 | j03_filtered_derived | valid | `01bdaedb361ba303` |  |
| family=bag;stmt=filtered_derived;threshold=20000 | j03_filtered_derived | valid | `5a0bb116dae8ec50` |  |
| family=bag;stmt=join_on;threshold=100000 | j02_join_on | valid | `19a9db0e61b598ab` |  |
| family=bag;stmt=join_on;threshold=180000 | j02_join_on | valid | `b785f3981285438c` |  |
| family=bag;stmt=join_on;threshold=20 | j02_join_on | valid | `3366ad2ab305165a` |  |
| family=bag;stmt=join_on;threshold=200 | j02_join_on | valid | `6ac0664dc98155dd` |  |
| family=bag;stmt=join_on;threshold=2000 | j02_join_on | valid | `acee2e5372c6423d` |  |
| family=bag;stmt=join_on;threshold=20000 | j02_join_on | valid | `c47e8274f3816da9` |  |
| family=bag;stmt=join_where;threshold=100000 | j01_join_where | valid | `660a9224f088f6a9` |  |
| family=bag;stmt=join_where;threshold=180000 | j01_join_where | valid | `820ceee1f2defe2e` |  |
| family=bag;stmt=join_where;threshold=20 | j01_join_where | valid | `f2ef651bcc72b866` |  |
| family=bag;stmt=join_where;threshold=200 | j01_join_where | valid | `3827d5f9b8b48653` |  |
| family=bag;stmt=join_where;threshold=2000 | j01_join_where | valid | `7952d4a4cb6c1d07` |  |
| family=bag;stmt=join_where;threshold=20000 | j01_join_where | valid | `1265c46703788a29` |  |
| family=reuse;policy=default;refs=16;consumer=1% | r_default_n16_c1 | valid | `59d871ffa0680a2f` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=16;consumer=100% | r_default_n16_c100 | valid | `481da00ea7550bd2` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=1;consumer=1% | r_default_n1_c1 | valid | `ad5744afcbe7b6b6` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=1;consumer=100% | r_default_n1_c100 | valid | `d8ebcb11b976ac36` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=2;consumer=1% | r_default_n2_c1 | valid | `254b1ba7d92b8ae0` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=2;consumer=100% | r_default_n2_c100 | valid | `8ed78cead8319934` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=1% | r_default_n4_c1 | valid | `f2b9388c72edd0e1` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=100% | r_default_n4_c100 | valid | `db4cb50a49872f43` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=8;consumer=1% | r_default_n8_c1 | valid | `16266d61c667863b` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=8;consumer=100% | r_default_n8_c100 | valid | `88059c8f08680832` | policy=default; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=1% | r_mat_n16_c1 | valid | `9e121a22d7f830b7` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=100% | r_mat_n16_c100 | valid | `587b4cb383830075` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=1% | r_mat_n1_c1 | valid | `99d110cf51c37192` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=100% | r_mat_n1_c100 | valid | `d19a3c2e0b10bd83` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=2;consumer=1% | r_mat_n2_c1 | valid | `d866cf131fd6f5f9` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=2;consumer=100% | r_mat_n2_c100 | valid | `95b46adef2209ea4` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=1% | r_mat_n4_c1 | valid | `889edb0ee90fce9c` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=100% | r_mat_n4_c100 | valid | `2bda37e1b03e45a2` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=8;consumer=1% | r_mat_n8_c1 | valid | `8d36c125916eff8b` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=8;consumer=100% | r_mat_n8_c100 | valid | `1c33b2ce3213a7ac` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=notmat;refs=16;consumer=1% | r_notmat_n16_c1 | valid | `bafa437c0105c768` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `4c882048f62cf8fd` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `ec810d1ab9da57d8` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `599aae1e435bf4e4` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=2;consumer=1% | r_notmat_n2_c1 | valid | `c5c0f68bdd0b9b7a` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=2;consumer=100% | r_notmat_n2_c100 | valid | `69fef1feaa7e3580` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `e8da14cf5750a312` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `12b5638e569a45ec` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=8;consumer=1% | r_notmat_n8_c1 | valid | `9fc14f17163c7d7a` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=8;consumer=100% | r_notmat_n8_c100 | valid | `6aa0ea938505e779` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, DIST, DEBUG, FORMAT TEXT` supported on yugabyte
- one host, single client: no independent hardware, real-network or equal-total-budget comparison
- no concurrent writers: the isolated read contract only; a matched blended read/write phase is a recorded gap
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
