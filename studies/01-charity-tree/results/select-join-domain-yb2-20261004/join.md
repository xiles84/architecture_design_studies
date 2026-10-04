# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 91 cells, 91 valid, engine `yugabyte`, topology `single`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-domain-yb2-20261004`
- run tag: `run/study-01/select-join-domain-yb2-20261004`
- environment: `host-zenbook-ux5406sa`
- repo commit: `e8b0fb1e56da6ae4b7a64f0256c588ea08e983d3`
- engine version: `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit`
- database engine image: `docker.io/yugabytedb/yugabyte:2025.2.6.0-b111` @ `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`
- benchmark client image: `localhost/charitybench-select:select-join-domain-yb2-20261004` @ `bc7ce0737685b8d0125c7291b019ef64d8dcf34c681bcc75569f23b2d083414b`
- inputs digest: `865c4471b57f5ba38e11dceafa59d6b6638d5abe6699d504e7d8e303469aff47`
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
| 20 (0.01%) | 82.6690 | 109.0463 | 78.9719 | 91.8535 | 83.4065 | 100.7976 | 90.9569 | 135.0339 | 81.2951 | 71.3233 |
| 200 (0.1%) | 104.4598 | 93.6990 | 97.2641 | 83.3770 | 74.1115 | 75.9594 | 94.2939 | 194.2958 | 71.0185 | 95.7748 |
| 2000 (1%) | 102.6413 | 97.2891 | 82.5772 | 89.6260 | 96.6450 | 80.1999 | 85.2398 | 144.3756 | 99.7000 | 74.1813 |
| 20000 (10%) | 90.0696 | 97.7491 | 116.9710 | 127.3603 | 95.2844 | 120.9388 | 116.0596 | 217.0542 | 105.9995 | 125.7364 |
| 100000 (50%) | 175.1617 | 138.1147 | 192.4074 | 152.3758 | 136.0317 | 125.1805 | 167.3406 | 179.2701 | 166.4903 | 178.5353 |
| 180000 (90%) | 267.9862 | 247.3008 | 231.3123 | 166.7558 | 215.6168 | 235.0769 | 243.0776 | 167.8199 | 200.0032 | 256.3938 |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | join_on vs join_where | 82.6690 | 109.0463 | 1.3191 | false |
| 20 | filtered_derived vs join_where | 82.6690 | 78.9719 | 0.9553 | false |
| 20 | cte_filtered vs join_where | 82.6690 | 91.8535 | 1.1111 | false |
| 20 | cte_filtered_mat vs join_where | 82.6690 | 83.4065 | 1.0089 | false |
| 20 | cte_filtered_notmat vs join_where | 82.6690 | 100.7976 | 1.2193 | false |
| 20 | cte_joined_outer vs join_where | 82.6690 | 90.9569 | 1.1003 | false |
| 20 | cte_joined_outer_mat vs join_where | 82.6690 | 135.0339 | 1.6334 | false |
| 20 | cte_joined_outer_notmat vs join_where | 82.6690 | 81.2951 | 0.9834 | false |
| 20 | comma_join vs join_where | 82.6690 | 71.3233 | 0.8628 | false |
| 200 | join_on vs join_where | 104.4598 | 93.6990 | 0.8970 | false |
| 200 | filtered_derived vs join_where | 104.4598 | 97.2641 | 0.9311 | false |
| 200 | cte_filtered vs join_where | 104.4598 | 83.3770 | 0.7982 | false |
| 200 | cte_filtered_mat vs join_where | 104.4598 | 74.1115 | 0.7095 | false |
| 200 | cte_filtered_notmat vs join_where | 104.4598 | 75.9594 | 0.7272 | false |
| 200 | cte_joined_outer vs join_where | 104.4598 | 94.2939 | 0.9027 | false |
| 200 | cte_joined_outer_mat vs join_where | 104.4598 | 194.2958 | 1.8600 | false |
| 200 | cte_joined_outer_notmat vs join_where | 104.4598 | 71.0185 | 0.6799 | false |
| 200 | comma_join vs join_where | 104.4598 | 95.7748 | 0.9169 | false |
| 2000 | join_on vs join_where | 102.6413 | 97.2891 | 0.9479 | false |
| 2000 | filtered_derived vs join_where | 102.6413 | 82.5772 | 0.8045 | false |
| 2000 | cte_filtered vs join_where | 102.6413 | 89.6260 | 0.8732 | false |
| 2000 | cte_filtered_mat vs join_where | 102.6413 | 96.6450 | 0.9416 | false |
| 2000 | cte_filtered_notmat vs join_where | 102.6413 | 80.1999 | 0.7814 | false |
| 2000 | cte_joined_outer vs join_where | 102.6413 | 85.2398 | 0.8305 | false |
| 2000 | cte_joined_outer_mat vs join_where | 102.6413 | 144.3756 | 1.4066 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 102.6413 | 99.7000 | 0.9713 | false |
| 2000 | comma_join vs join_where | 102.6413 | 74.1813 | 0.7227 | false |
| 20000 | join_on vs join_where | 90.0696 | 97.7491 | 1.0853 | false |
| 20000 | filtered_derived vs join_where | 90.0696 | 116.9710 | 1.2987 | false |
| 20000 | cte_filtered vs join_where | 90.0696 | 127.3603 | 1.4140 | false |
| 20000 | cte_filtered_mat vs join_where | 90.0696 | 95.2844 | 1.0579 | false |
| 20000 | cte_filtered_notmat vs join_where | 90.0696 | 120.9388 | 1.3427 | false |
| 20000 | cte_joined_outer vs join_where | 90.0696 | 116.0596 | 1.2886 | false |
| 20000 | cte_joined_outer_mat vs join_where | 90.0696 | 217.0542 | 2.4098 | false |
| 20000 | cte_joined_outer_notmat vs join_where | 90.0696 | 105.9995 | 1.1769 | false |
| 20000 | comma_join vs join_where | 90.0696 | 125.7364 | 1.3960 | false |
| 100000 | join_on vs join_where | 175.1617 | 138.1147 | 0.7885 | false |
| 100000 | filtered_derived vs join_where | 175.1617 | 192.4074 | 1.0985 | false |
| 100000 | cte_filtered vs join_where | 175.1617 | 152.3758 | 0.8699 | false |
| 100000 | cte_filtered_mat vs join_where | 175.1617 | 136.0317 | 0.7766 | false |
| 100000 | cte_filtered_notmat vs join_where | 175.1617 | 125.1805 | 0.7147 | false |
| 100000 | cte_joined_outer vs join_where | 175.1617 | 167.3406 | 0.9553 | false |
| 100000 | cte_joined_outer_mat vs join_where | 175.1617 | 179.2701 | 1.0235 | false |
| 100000 | cte_joined_outer_notmat vs join_where | 175.1617 | 166.4903 | 0.9505 | false |
| 100000 | comma_join vs join_where | 175.1617 | 178.5353 | 1.0193 | false |
| 180000 | join_on vs join_where | 267.9862 | 247.3008 | 0.9228 | false |
| 180000 | filtered_derived vs join_where | 267.9862 | 231.3123 | 0.8631 | false |
| 180000 | cte_filtered vs join_where | 267.9862 | 166.7558 | 0.6223 | false |
| 180000 | cte_filtered_mat vs join_where | 267.9862 | 215.6168 | 0.8046 | false |
| 180000 | cte_filtered_notmat vs join_where | 267.9862 | 235.0769 | 0.8772 | false |
| 180000 | cte_joined_outer vs join_where | 267.9862 | 243.0776 | 0.9071 | false |
| 180000 | cte_joined_outer_mat vs join_where | 267.9862 | 167.8199 | 0.6262 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 267.9862 | 200.0032 | 0.7463 | false |
| 180000 | comma_join vs join_where | 267.9862 | 256.3938 | 0.9567 | false |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_filtered_mat vs cte_filtered | 91.8535 | 83.4065 | 0.9080 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_filtered_notmat vs cte_filtered | 91.8535 | 100.7976 | 1.0974 | false | policy=not_materialized; observed_cte_scan=false |
| 20 | cte_joined_outer_mat vs cte_joined_outer | 90.9569 | 135.0339 | 1.4846 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 90.9569 | 81.2951 | 0.8938 | false | policy=not_materialized; observed_cte_scan=false |
| 200 | cte_filtered_mat vs cte_filtered | 83.3770 | 74.1115 | 0.8889 | false | policy=materialized; observed_cte_scan=true |
| 200 | cte_filtered_notmat vs cte_filtered | 83.3770 | 75.9594 | 0.9110 | false | policy=not_materialized; observed_cte_scan=false |
| 200 | cte_joined_outer_mat vs cte_joined_outer | 94.2939 | 194.2958 | 2.0605 | false | policy=materialized; observed_cte_scan=true |
| 200 | cte_joined_outer_notmat vs cte_joined_outer | 94.2939 | 71.0185 | 0.7532 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_filtered_mat vs cte_filtered | 89.6260 | 96.6450 | 1.0783 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_filtered_notmat vs cte_filtered | 89.6260 | 80.1999 | 0.8948 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 85.2398 | 144.3756 | 1.6938 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 85.2398 | 99.7000 | 1.1696 | false | policy=not_materialized; observed_cte_scan=false |
| 20000 | cte_filtered_mat vs cte_filtered | 127.3603 | 95.2844 | 0.7481 | false | policy=materialized; observed_cte_scan=true |
| 20000 | cte_filtered_notmat vs cte_filtered | 127.3603 | 120.9388 | 0.9496 | false | policy=not_materialized; observed_cte_scan=false |
| 20000 | cte_joined_outer_mat vs cte_joined_outer | 116.0596 | 217.0542 | 1.8702 | false | policy=materialized; observed_cte_scan=true |
| 20000 | cte_joined_outer_notmat vs cte_joined_outer | 116.0596 | 105.9995 | 0.9133 | false | policy=not_materialized; observed_cte_scan=false |
| 100000 | cte_filtered_mat vs cte_filtered | 152.3758 | 136.0317 | 0.8927 | false | policy=materialized; observed_cte_scan=true |
| 100000 | cte_filtered_notmat vs cte_filtered | 152.3758 | 125.1805 | 0.8215 | false | policy=not_materialized; observed_cte_scan=false |
| 100000 | cte_joined_outer_mat vs cte_joined_outer | 167.3406 | 179.2701 | 1.0713 | false | policy=materialized; observed_cte_scan=true |
| 100000 | cte_joined_outer_notmat vs cte_joined_outer | 167.3406 | 166.4903 | 0.9949 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_filtered_mat vs cte_filtered | 166.7558 | 215.6168 | 1.2930 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_filtered_notmat vs cte_filtered | 166.7558 | 235.0769 | 1.4097 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 243.0776 | 167.8199 | 0.6904 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 243.0776 | 200.0032 | 0.8228 | false | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 97.8882 | 1 | `165362f15a6ea244` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 75.9922 | 1 | `b22677ec0ec02248` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | 80.2171 | 2 | `6ae27b58bbc4256a` | policy=default; observed_cte_scan=true | valid |
| default | 2 | 100% | 104.6142 | 2 | `ba6fab2e54b29a45` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 1% | 95.9475 | 4 | `1e8f6f01bcb192b4` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 84.1687 | 4 | `236c6aee93f6c37d` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | 82.8818 | 8 | `2c7c2dc5cb54295c` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 100% | 105.0926 | 8 | `a8db5c2d460bab96` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 1% | 95.6351 | 16 | `eaf10ea8374d564a` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 89.4784 | 16 | `fd84a7ab5371f4f5` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 87.8541 | 1 | `54567e5de7d3d133` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 87.4619 | 1 | `4690196800d0033f` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | 98.5504 | 2 | `62d2dd841d57ca7d` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 100% | 84.7355 | 2 | `49672e055da6c5b6` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 1% | 86.6428 | 4 | `5a9746bebcabb9df` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 94.4433 | 4 | `33911c082a983e51` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | 99.9630 | 8 | `4b139a2265bc6e1f` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 100% | 79.2663 | 8 | `1d1d31283489f2f1` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 1% | 73.6282 | 16 | `03b11860ecf0f1f0` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 83.1863 | 16 | `a4522afa62b8401f` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 98.4189 | 1 | `35244fa57bbc9aaf` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 110.1305 | 1 | `322c8d72f420230b` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | 216.0452 | 2 | `42e387d3debc8098` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 100% | 250.5153 | 2 | `e6ee3a8f59761d4a` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 1% | 344.3300 | 4 | `a6ebe1684ec4c4ec` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 337.3900 | 4 | `78d3779f6a8ce24c` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | 592.9320 | 8 | `e9743400383e8dbd` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 100% | 656.6843 | 8 | `35699f12b4a4b4f5` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 1% | 1033.6538 | 16 | `b0714f1432092995` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 1035.6494 | 16 | `82431be763abd93e` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 0.8975 | false |
| notmat | 1 | 1% | 1.0054 | false |
| mat | 1 | 100% | 1.1509 | false |
| notmat | 1 | 100% | 1.4492 | false |
| mat | 2 | 1% | 1.2285 | false |
| notmat | 2 | 1% | 2.6933 | false |
| mat | 2 | 100% | 0.8100 | false |
| notmat | 2 | 100% | 2.3947 | false |
| mat | 4 | 1% | 0.9030 | false |
| notmat | 4 | 1% | 3.5887 | false |
| mat | 4 | 100% | 1.1221 | false |
| notmat | 4 | 100% | 4.0085 | false |
| mat | 8 | 1% | 1.2061 | false |
| notmat | 8 | 1% | 7.1539 | false |
| mat | 8 | 100% | 0.7543 | false |
| notmat | 8 | 100% | 6.2486 | false |
| mat | 16 | 1% | 0.7699 | false |
| notmat | 16 | 1% | 10.8083 | false |
| mat | 16 | 100% | 0.9297 | false |
| notmat | 16 | 100% | 11.5743 | false |

## Repeated-design position control

| first median ms | last median ms | ratio (last/first) | same plan |
|---|---|---|---|
| 215.616795 | 198.234593 | 0.9194 | false |

## Latency samples and derived throughput

_The reader is sequential, so `queries/s` is 1000/median_ms, a derived rate and not a concurrent-capacity test. min/max are the five sample extremes._

| axis | status | min ms | median ms | max ms | samples | derived queries/s |
|---|---|---|---|---|---|---|
| family=bag;stmt=comma_join;threshold=100000 | valid | 125.4555 | 178.5353 | 243.2936 | 5 | 5.6 |
| family=bag;stmt=comma_join;threshold=180000 | valid | 149.5253 | 256.3938 | 266.6863 | 5 | 3.9 |
| family=bag;stmt=comma_join;threshold=20 | valid | 69.2740 | 71.3233 | 111.9297 | 5 | 14.0 |
| family=bag;stmt=comma_join;threshold=200 | valid | 68.8773 | 95.7748 | 110.2121 | 5 | 10.4 |
| family=bag;stmt=comma_join;threshold=2000 | valid | 70.9448 | 74.1813 | 109.9382 | 5 | 13.5 |
| family=bag;stmt=comma_join;threshold=20000 | valid | 87.1348 | 125.7364 | 138.1663 | 5 | 8.0 |
| family=bag;stmt=cte_filtered;threshold=100000 | valid | 131.9335 | 152.3758 | 190.1015 | 5 | 6.6 |
| family=bag;stmt=cte_filtered;threshold=180000 | valid | 154.2024 | 166.7558 | 254.0114 | 5 | 6.0 |
| family=bag;stmt=cte_filtered;threshold=20 | valid | 79.5040 | 91.8535 | 105.4548 | 5 | 10.9 |
| family=bag;stmt=cte_filtered;threshold=200 | valid | 71.3568 | 83.3770 | 106.0139 | 5 | 12.0 |
| family=bag;stmt=cte_filtered;threshold=2000 | valid | 69.2299 | 89.6260 | 108.9800 | 5 | 11.2 |
| family=bag;stmt=cte_filtered;threshold=20000 | valid | 100.7303 | 127.3603 | 142.6078 | 5 | 7.9 |
| family=bag;stmt=cte_filtered_mat;threshold=100000 | valid | 122.1924 | 136.0317 | 224.2656 | 5 | 7.4 |
| family=bag;stmt=cte_filtered_mat;threshold=180000 | valid | 195.0193 | 215.6168 | 299.9636 | 5 | 4.6 |
| family=bag;stmt=cte_filtered_mat;threshold=180000;position=last | valid | 183.3565 | 198.2346 | 211.8414 | 5 | 5.0 |
| family=bag;stmt=cte_filtered_mat;threshold=20 | valid | 69.1753 | 83.4065 | 113.6063 | 5 | 12.0 |
| family=bag;stmt=cte_filtered_mat;threshold=200 | valid | 67.6295 | 74.1115 | 110.1875 | 5 | 13.5 |
| family=bag;stmt=cte_filtered_mat;threshold=2000 | valid | 68.4883 | 96.6450 | 113.0272 | 5 | 10.3 |
| family=bag;stmt=cte_filtered_mat;threshold=20000 | valid | 89.1807 | 95.2844 | 112.6826 | 5 | 10.5 |
| family=bag;stmt=cte_filtered_notmat;threshold=100000 | valid | 115.8973 | 125.1805 | 189.2774 | 5 | 8.0 |
| family=bag;stmt=cte_filtered_notmat;threshold=180000 | valid | 165.5253 | 235.0769 | 249.9930 | 5 | 4.3 |
| family=bag;stmt=cte_filtered_notmat;threshold=20 | valid | 66.8448 | 100.7976 | 115.8586 | 5 | 9.9 |
| family=bag;stmt=cte_filtered_notmat;threshold=200 | valid | 71.1849 | 75.9594 | 88.2109 | 5 | 13.2 |
| family=bag;stmt=cte_filtered_notmat;threshold=2000 | valid | 71.3249 | 80.1999 | 113.4606 | 5 | 12.5 |
| family=bag;stmt=cte_filtered_notmat;threshold=20000 | valid | 89.2489 | 120.9388 | 134.7353 | 5 | 8.3 |
| family=bag;stmt=cte_joined_outer;threshold=100000 | valid | 125.6990 | 167.3406 | 212.1670 | 5 | 6.0 |
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 163.5420 | 243.0776 | 286.3225 | 5 | 4.1 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 74.9711 | 90.9569 | 106.0900 | 5 | 11.0 |
| family=bag;stmt=cte_joined_outer;threshold=200 | valid | 71.4590 | 94.2939 | 110.9628 | 5 | 10.6 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 70.7386 | 85.2398 | 120.7792 | 5 | 11.7 |
| family=bag;stmt=cte_joined_outer;threshold=20000 | valid | 110.6891 | 116.0596 | 132.3317 | 5 | 8.6 |
| family=bag;stmt=cte_joined_outer_mat;threshold=100000 | valid | 137.3149 | 179.2701 | 207.1787 | 5 | 5.6 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 131.3161 | 167.8199 | 237.3086 | 5 | 6.0 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 112.3852 | 135.0339 | 176.6703 | 5 | 7.4 |
| family=bag;stmt=cte_joined_outer_mat;threshold=200 | valid | 126.2995 | 194.2958 | 251.4050 | 5 | 5.1 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 126.9997 | 144.3756 | 236.2690 | 5 | 6.9 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20000 | valid | 167.2341 | 217.0542 | 225.8298 | 5 | 4.6 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=100000 | valid | 131.4471 | 166.4903 | 242.5769 | 5 | 6.0 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 168.6875 | 200.0032 | 247.3324 | 5 | 5.0 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 74.8944 | 81.2951 | 110.5743 | 5 | 12.3 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=200 | valid | 67.5023 | 71.0185 | 96.7726 | 5 | 14.1 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 75.2265 | 99.7000 | 116.3464 | 5 | 10.0 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20000 | valid | 84.6466 | 105.9995 | 129.9759 | 5 | 9.4 |
| family=bag;stmt=filtered_derived;threshold=100000 | valid | 129.4794 | 192.4074 | 217.5578 | 5 | 5.2 |
| family=bag;stmt=filtered_derived;threshold=180000 | valid | 189.8691 | 231.3123 | 278.8663 | 5 | 4.3 |
| family=bag;stmt=filtered_derived;threshold=20 | valid | 77.4304 | 78.9719 | 130.5345 | 5 | 12.7 |
| family=bag;stmt=filtered_derived;threshold=200 | valid | 69.2448 | 97.2641 | 110.3322 | 5 | 10.3 |
| family=bag;stmt=filtered_derived;threshold=2000 | valid | 71.9627 | 82.5772 | 100.8217 | 5 | 12.1 |
| family=bag;stmt=filtered_derived;threshold=20000 | valid | 89.0048 | 116.9710 | 125.1634 | 5 | 8.5 |
| family=bag;stmt=join_on;threshold=100000 | valid | 131.4940 | 138.1147 | 184.4842 | 5 | 7.2 |
| family=bag;stmt=join_on;threshold=180000 | valid | 174.8915 | 247.3008 | 284.2951 | 5 | 4.0 |
| family=bag;stmt=join_on;threshold=20 | valid | 82.7319 | 109.0463 | 132.8849 | 5 | 9.2 |
| family=bag;stmt=join_on;threshold=200 | valid | 76.5614 | 93.6990 | 127.9793 | 5 | 10.7 |
| family=bag;stmt=join_on;threshold=2000 | valid | 70.8475 | 97.2891 | 111.1697 | 5 | 10.3 |
| family=bag;stmt=join_on;threshold=20000 | valid | 89.5773 | 97.7491 | 142.9604 | 5 | 10.2 |
| family=bag;stmt=join_where;threshold=100000 | valid | 124.6993 | 175.1617 | 182.9949 | 5 | 5.7 |
| family=bag;stmt=join_where;threshold=180000 | valid | 159.9825 | 267.9862 | 290.0985 | 5 | 3.7 |
| family=bag;stmt=join_where;threshold=20 | valid | 68.7937 | 82.6690 | 104.7700 | 5 | 12.1 |
| family=bag;stmt=join_where;threshold=200 | valid | 87.9871 | 104.4598 | 117.1572 | 5 | 9.6 |
| family=bag;stmt=join_where;threshold=2000 | valid | 71.8711 | 102.6413 | 134.1941 | 5 | 9.7 |
| family=bag;stmt=join_where;threshold=20000 | valid | 84.8194 | 90.0696 | 111.0740 | 5 | 11.1 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 72.6364 | 95.6351 | 144.5873 | 5 | 10.5 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 78.3665 | 89.4784 | 128.8327 | 5 | 11.2 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 77.9772 | 97.8882 | 113.2269 | 5 | 10.2 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 73.6999 | 75.9922 | 110.2427 | 5 | 13.2 |
| family=reuse;policy=default;refs=2;consumer=1% | valid | 69.5058 | 80.2171 | 117.2276 | 5 | 12.5 |
| family=reuse;policy=default;refs=2;consumer=100% | valid | 80.1963 | 104.6142 | 139.9186 | 5 | 9.6 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 72.2385 | 95.9475 | 102.5995 | 5 | 10.4 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 72.3243 | 84.1687 | 118.2598 | 5 | 11.9 |
| family=reuse;policy=default;refs=8;consumer=1% | valid | 78.4750 | 82.8818 | 142.1289 | 5 | 12.1 |
| family=reuse;policy=default;refs=8;consumer=100% | valid | 70.7724 | 105.0926 | 127.7784 | 5 | 9.5 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 69.6924 | 73.6282 | 106.9794 | 5 | 13.6 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 69.4715 | 83.1863 | 100.8803 | 5 | 12.0 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 80.2481 | 87.8541 | 121.6974 | 5 | 11.4 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 67.9705 | 87.4619 | 121.2949 | 5 | 11.4 |
| family=reuse;policy=mat;refs=2;consumer=1% | valid | 80.6553 | 98.5504 | 138.5955 | 5 | 10.1 |
| family=reuse;policy=mat;refs=2;consumer=100% | valid | 72.8609 | 84.7355 | 114.1910 | 5 | 11.8 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 70.7382 | 86.6428 | 102.7144 | 5 | 11.5 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 71.3615 | 94.4433 | 115.8292 | 5 | 10.6 |
| family=reuse;policy=mat;refs=8;consumer=1% | valid | 69.9839 | 99.9630 | 107.5144 | 5 | 10.0 |
| family=reuse;policy=mat;refs=8;consumer=100% | valid | 74.1321 | 79.2663 | 138.5164 | 5 | 12.6 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 926.4951 | 1033.6538 | 1262.9827 | 5 | 1.0 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 760.9230 | 1035.6494 | 1316.6991 | 5 | 1.0 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 85.6912 | 98.4189 | 117.7920 | 5 | 10.2 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 79.9335 | 110.1305 | 147.0947 | 5 | 9.1 |
| family=reuse;policy=notmat;refs=2;consumer=1% | valid | 132.3455 | 216.0452 | 257.3783 | 5 | 4.6 |
| family=reuse;policy=notmat;refs=2;consumer=100% | valid | 213.3624 | 250.5153 | 300.6050 | 5 | 4.0 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 226.3762 | 344.3300 | 378.9052 | 5 | 2.9 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 287.6854 | 337.3900 | 354.7061 | 5 | 3.0 |
| family=reuse;policy=notmat;refs=8;consumer=1% | valid | 409.9215 | 592.9320 | 739.1843 | 5 | 1.7 |
| family=reuse;policy=notmat;refs=8;consumer=100% | valid | 483.4675 | 656.6843 | 741.0263 | 5 | 1.5 |

## Execution order

- order seed `42174513`, 5 interleaved rounds
  1. family=bag;stmt=cte_filtered_mat;threshold=180000
  2. family=bag;stmt=cte_filtered;threshold=200
  3. family=bag;stmt=comma_join;threshold=200
  4. family=reuse;policy=notmat;refs=2;consumer=1%
  5. family=bag;stmt=cte_filtered_notmat;threshold=2000
  6. family=bag;stmt=cte_filtered_mat;threshold=20
  7. family=reuse;policy=notmat;refs=8;consumer=1%
  8. family=reuse;policy=mat;refs=4;consumer=1%
  9. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  10. family=bag;stmt=cte_filtered;threshold=100000
  11. family=bag;stmt=join_on;threshold=20000
  12. family=bag;stmt=comma_join;threshold=20
  13. family=bag;stmt=comma_join;threshold=2000
  14. family=bag;stmt=cte_joined_outer_mat;threshold=20
  15. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  16. family=bag;stmt=comma_join;threshold=180000
  17. family=bag;stmt=cte_joined_outer;threshold=100000
  18. family=bag;stmt=comma_join;threshold=20000
  19. family=reuse;policy=mat;refs=8;consumer=1%
  20. family=bag;stmt=filtered_derived;threshold=180000
  21. family=bag;stmt=cte_joined_outer;threshold=200
  22. family=bag;stmt=cte_joined_outer_mat;threshold=20000
  23. family=bag;stmt=cte_filtered;threshold=2000
  24. family=bag;stmt=filtered_derived;threshold=100000
  25. family=bag;stmt=join_on;threshold=180000
  26. family=bag;stmt=cte_joined_outer_notmat;threshold=100000
  27. family=bag;stmt=cte_filtered_notmat;threshold=180000
  28. family=reuse;policy=default;refs=2;consumer=100%
  29. family=reuse;policy=notmat;refs=4;consumer=100%
  30. family=bag;stmt=join_where;threshold=200
  31. family=bag;stmt=join_on;threshold=20
  32. family=bag;stmt=join_on;threshold=2000
  33. family=reuse;policy=default;refs=4;consumer=1%
  34. family=bag;stmt=filtered_derived;threshold=20000
  35. family=bag;stmt=cte_filtered_mat;threshold=2000
  36. family=bag;stmt=cte_filtered_notmat;threshold=20000
  37. family=bag;stmt=join_where;threshold=180000
  38. family=bag;stmt=join_on;threshold=100000
  39. family=bag;stmt=cte_joined_outer_mat;threshold=200
  40. family=bag;stmt=cte_filtered_notmat;threshold=200
  41. family=reuse;policy=notmat;refs=16;consumer=1%
  42. family=reuse;policy=mat;refs=2;consumer=1%
  43. family=bag;stmt=filtered_derived;threshold=200
  44. family=bag;stmt=join_where;threshold=100000
  45. family=reuse;policy=default;refs=1;consumer=100%
  46. family=bag;stmt=cte_filtered;threshold=180000
  47. family=bag;stmt=comma_join;threshold=100000
  48. family=bag;stmt=cte_filtered_mat;threshold=200
  49. family=reuse;policy=notmat;refs=2;consumer=100%
  50. family=bag;stmt=cte_joined_outer;threshold=20000
  51. family=reuse;policy=mat;refs=1;consumer=100%
  52. family=reuse;policy=notmat;refs=1;consumer=1%
  53. family=bag;stmt=cte_filtered_notmat;threshold=20
  54. family=reuse;policy=mat;refs=16;consumer=1%
  55. family=bag;stmt=cte_joined_outer_notmat;threshold=200
  56. family=reuse;policy=notmat;refs=16;consumer=100%
  57. family=reuse;policy=default;refs=8;consumer=1%
  58. family=reuse;policy=notmat;refs=8;consumer=100%
  59. family=bag;stmt=cte_filtered;threshold=20000
  60. family=reuse;policy=mat;refs=4;consumer=100%
  61. family=bag;stmt=cte_joined_outer;threshold=180000
  62. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  63. family=reuse;policy=default;refs=16;consumer=1%
  64. family=reuse;policy=notmat;refs=1;consumer=100%
  65. family=bag;stmt=filtered_derived;threshold=20
  66. family=bag;stmt=cte_joined_outer;threshold=20
  67. family=bag;stmt=cte_joined_outer;threshold=2000
  68. family=reuse;policy=mat;refs=8;consumer=100%
  69. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  70. family=reuse;policy=mat;refs=2;consumer=100%
  71. family=reuse;policy=mat;refs=1;consumer=1%
  72. family=reuse;policy=notmat;refs=4;consumer=1%
  73. family=bag;stmt=join_where;threshold=2000
  74. family=reuse;policy=default;refs=2;consumer=1%
  75. family=bag;stmt=cte_filtered_notmat;threshold=100000
  76. family=reuse;policy=default;refs=4;consumer=100%
  77. family=bag;stmt=cte_filtered_mat;threshold=100000
  78. family=reuse;policy=default;refs=16;consumer=100%
  79. family=bag;stmt=join_on;threshold=200
  80. family=reuse;policy=default;refs=8;consumer=100%
  81. family=bag;stmt=cte_filtered;threshold=20
  82. family=reuse;policy=default;refs=1;consumer=1%
  83. family=bag;stmt=cte_joined_outer_notmat;threshold=20000
  84. family=bag;stmt=join_where;threshold=20
  85. family=reuse;policy=mat;refs=16;consumer=100%
  86. family=bag;stmt=cte_filtered_mat;threshold=20000
  87. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  88. family=bag;stmt=cte_joined_outer_mat;threshold=100000
  89. family=bag;stmt=filtered_derived;threshold=2000
  90. family=bag;stmt=join_where;threshold=20000
  91. family=bag;stmt=cte_filtered_mat;threshold=180000;position=last

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
| family=bag;stmt=comma_join;threshold=100000 | j06_comma_join | valid | `c2a661839eacb796` |  |
| family=bag;stmt=comma_join;threshold=180000 | j06_comma_join | valid | `d4c35ca36120991a` |  |
| family=bag;stmt=comma_join;threshold=20 | j06_comma_join | valid | `42acd6114b609f28` |  |
| family=bag;stmt=comma_join;threshold=200 | j06_comma_join | valid | `e193b6f4824facb1` |  |
| family=bag;stmt=comma_join;threshold=2000 | j06_comma_join | valid | `98e84910e04e5046` |  |
| family=bag;stmt=comma_join;threshold=20000 | j06_comma_join | valid | `9bfa5ad09db526f9` |  |
| family=bag;stmt=cte_filtered;threshold=100000 | j04_cte_filtered | valid | `aa40809113242a91` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=180000 | j04_cte_filtered | valid | `a23118328616942e` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=20 | j04_cte_filtered | valid | `8c2e5519794b1e83` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=200 | j04_cte_filtered | valid | `bcc09d8ddd38d1d1` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=2000 | j04_cte_filtered | valid | `01db388421217671` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered;threshold=20000 | j04_cte_filtered | valid | `c229a0c2c8991f4f` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_mat;threshold=100000 | j04m_cte_filtered_mat | valid | `55850bdeefdbbe61` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=180000 | j04m_cte_filtered_mat | valid | `6e23b22dabf82283` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=180000;position=last | j04m_cte_filtered_mat | valid | `45eb97fd151c5df0` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=20 | j04m_cte_filtered_mat | valid | `5c0c5cc40890af47` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=200 | j04m_cte_filtered_mat | valid | `347acea6982a55cd` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=2000 | j04m_cte_filtered_mat | valid | `3165f14d724973c5` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_mat;threshold=20000 | j04m_cte_filtered_mat | valid | `f18769cc5163ffdd` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_filtered_notmat;threshold=100000 | j04n_cte_filtered_notmat | valid | `b52054f70a006bfd` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=180000 | j04n_cte_filtered_notmat | valid | `bd9a6f42b87d6ff3` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=20 | j04n_cte_filtered_notmat | valid | `ef8cfec717ac25b1` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=200 | j04n_cte_filtered_notmat | valid | `7986c06d8e5a4816` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=2000 | j04n_cte_filtered_notmat | valid | `30cb524e8b0b971d` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_filtered_notmat;threshold=20000 | j04n_cte_filtered_notmat | valid | `554cce784f77bc98` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=100000 | j05_cte_joined_outer | valid | `dc2914d6e2d61b4f` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=180000 | j05_cte_joined_outer | valid | `aa6dbd1b34c3f460` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20 | j05_cte_joined_outer | valid | `dafbce2fae25f666` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=200 | j05_cte_joined_outer | valid | `4969790424b092a5` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=2000 | j05_cte_joined_outer | valid | `e1be1caa82ae0be7` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20000 | j05_cte_joined_outer | valid | `93912b3d5a7dc757` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_mat;threshold=100000 | j05m_cte_joined_outer_mat | valid | `98e4d4923ff2f48b` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | j05m_cte_joined_outer_mat | valid | `f9cdec98b0b0a2a9` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | j05m_cte_joined_outer_mat | valid | `1dfaeb486799d28b` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=200 | j05m_cte_joined_outer_mat | valid | `5641c579fb0dc257` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | j05m_cte_joined_outer_mat | valid | `32017806fd2f1156` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20000 | j05m_cte_joined_outer_mat | valid | `809d3a01ae63adb3` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_notmat;threshold=100000 | j05n_cte_joined_outer_notmat | valid | `8801a7eecde5b9c5` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | j05n_cte_joined_outer_notmat | valid | `f6fea04b97a27181` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | j05n_cte_joined_outer_notmat | valid | `66823415312bc514` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=200 | j05n_cte_joined_outer_notmat | valid | `47b312a16a3fca14` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | j05n_cte_joined_outer_notmat | valid | `525164020e6351f6` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20000 | j05n_cte_joined_outer_notmat | valid | `dd307fa909f504d6` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=filtered_derived;threshold=100000 | j03_filtered_derived | valid | `8fc0aec4fe1ac31c` |  |
| family=bag;stmt=filtered_derived;threshold=180000 | j03_filtered_derived | valid | `9b1a17e9f3e6ed15` |  |
| family=bag;stmt=filtered_derived;threshold=20 | j03_filtered_derived | valid | `b60eb6062c6745af` |  |
| family=bag;stmt=filtered_derived;threshold=200 | j03_filtered_derived | valid | `fa4ac6eac952ac01` |  |
| family=bag;stmt=filtered_derived;threshold=2000 | j03_filtered_derived | valid | `75d20aa995960159` |  |
| family=bag;stmt=filtered_derived;threshold=20000 | j03_filtered_derived | valid | `532577cd97acdc2c` |  |
| family=bag;stmt=join_on;threshold=100000 | j02_join_on | valid | `0411509a806d61e6` |  |
| family=bag;stmt=join_on;threshold=180000 | j02_join_on | valid | `fa1ccfa8a7f18ea7` |  |
| family=bag;stmt=join_on;threshold=20 | j02_join_on | valid | `400966ef3de9d97f` |  |
| family=bag;stmt=join_on;threshold=200 | j02_join_on | valid | `40742602d35b607e` |  |
| family=bag;stmt=join_on;threshold=2000 | j02_join_on | valid | `7d07143727ebb4fa` |  |
| family=bag;stmt=join_on;threshold=20000 | j02_join_on | valid | `2168721cc48dda67` |  |
| family=bag;stmt=join_where;threshold=100000 | j01_join_where | valid | `8345fb7e73392bce` |  |
| family=bag;stmt=join_where;threshold=180000 | j01_join_where | valid | `19780546faf938ae` |  |
| family=bag;stmt=join_where;threshold=20 | j01_join_where | valid | `d162f9f8495cb2d2` |  |
| family=bag;stmt=join_where;threshold=200 | j01_join_where | valid | `c0d8ca6ce7f1324c` |  |
| family=bag;stmt=join_where;threshold=2000 | j01_join_where | valid | `4cbe1e2ab771ef8c` |  |
| family=bag;stmt=join_where;threshold=20000 | j01_join_where | valid | `c8d554a59983029b` |  |
| family=reuse;policy=default;refs=16;consumer=1% | r_default_n16_c1 | valid | `eaf10ea8374d564a` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=16;consumer=100% | r_default_n16_c100 | valid | `fd84a7ab5371f4f5` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=1;consumer=1% | r_default_n1_c1 | valid | `165362f15a6ea244` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=1;consumer=100% | r_default_n1_c100 | valid | `b22677ec0ec02248` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=2;consumer=1% | r_default_n2_c1 | valid | `6ae27b58bbc4256a` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=2;consumer=100% | r_default_n2_c100 | valid | `ba6fab2e54b29a45` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=1% | r_default_n4_c1 | valid | `1e8f6f01bcb192b4` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=100% | r_default_n4_c100 | valid | `236c6aee93f6c37d` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=8;consumer=1% | r_default_n8_c1 | valid | `2c7c2dc5cb54295c` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=8;consumer=100% | r_default_n8_c100 | valid | `a8db5c2d460bab96` | policy=default; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=1% | r_mat_n16_c1 | valid | `03b11860ecf0f1f0` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=100% | r_mat_n16_c100 | valid | `a4522afa62b8401f` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=1% | r_mat_n1_c1 | valid | `54567e5de7d3d133` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=100% | r_mat_n1_c100 | valid | `4690196800d0033f` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=2;consumer=1% | r_mat_n2_c1 | valid | `62d2dd841d57ca7d` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=2;consumer=100% | r_mat_n2_c100 | valid | `49672e055da6c5b6` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=1% | r_mat_n4_c1 | valid | `5a9746bebcabb9df` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=100% | r_mat_n4_c100 | valid | `33911c082a983e51` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=8;consumer=1% | r_mat_n8_c1 | valid | `4b139a2265bc6e1f` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=8;consumer=100% | r_mat_n8_c100 | valid | `1d1d31283489f2f1` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=notmat;refs=16;consumer=1% | r_notmat_n16_c1 | valid | `b0714f1432092995` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `82431be763abd93e` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `35244fa57bbc9aaf` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `322c8d72f420230b` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=2;consumer=1% | r_notmat_n2_c1 | valid | `42e387d3debc8098` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=2;consumer=100% | r_notmat_n2_c100 | valid | `e6ee3a8f59761d4a` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `a6ebe1684ec4c4ec` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `78d3779f6a8ce24c` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=8;consumer=1% | r_notmat_n8_c1 | valid | `e9743400383e8dbd` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=8;consumer=100% | r_notmat_n8_c100 | valid | `35699f12b4a4b4f5` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, DIST, DEBUG, FORMAT TEXT` supported on yugabyte
- one host, single client: no independent hardware, real-network or equal-total-budget comparison
- no concurrent writers: the isolated read contract only; a matched blended read/write phase is a recorded gap
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
