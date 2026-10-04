# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 30 cells, 30 valid, engine `yugabyte`, topology `cluster3`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-topo-yb-cluster-bo1-20261004`
- run tag: `run/study-01/select-join-topo-yb-cluster-bo1-20261004`
- environment: `host-zenbook-ux5406sa`
- repo commit: `0da058e8bbdc36acfdeec4082f56a5b0c8068a8e`
- engine version: `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit`
- database engine image: `docker.io/yugabytedb/yugabyte:2025.2.6.0-b111` @ `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`
- benchmark client image: `localhost/charitybench-select:select-join-topo-yb-cluster-bo1-20261004` @ `79bc753ff7e57e102982a323939c43e68e9b3664887683d3a5d5bf072650f4c1`
- inputs digest: `361859783b15304847cbd360653e6a106cfba0aa3754cfc04d40d936ee3032eb`
- dataset fingerprint: `fcf26e08f377ec51`
- client CPU quota: `100000 100000`, memory: `2147483648`
- setting `client_workers`: 1 sequential reader; 0 writers
- setting `connection_pool`: 2 per endpoint = reader/setup + reserve
- setting `customers`: 2000
- setting `dataset`: study01-person-donation-charity
- setting `distributed_rpc`: not parsed; raw EXPLAIN retained (coverage gap)
- setting `endpoint_mode`: round-robin across 3 endpoints
- setting `endpoints`: yb-n1,yb-n2,yb-n3
- setting `fact_count`: 200000
- setting `index_set`: child rank index and parent-key index, primary keys; no per-spelling index
- setting `phase`: blended
- setting `placement`: noncolocated
- setting `prepared_plan_policy`: pgx simple protocol; no prepared statements
- setting `scale`: cte-domain
- setting `selectivity_mapping`: rank <= threshold; rank is a deterministic shuffle of 1..fact_count
- setting `statement_budget_ms`: 120000
- setting `statistics`: ANALYZE after deterministic load
- setting `thresholds`: [20 2000 180000]
- setting `topology`: cluster3
- setting `trials_per_cell`: 5
- setting `whole_run_count`: 1; no replicated winner in this file
- setting `writer_keys`: 1000
- setting `writer_target_rate_per_sec`: 10

### Endpoint distribution

A cluster arm must rotate its reader across nodes; a single seed endpoint would not be balanced cluster access.

| endpoint | operations | connections at end |
|---|---|---|
| yb-n1 | 163 | 2 |
| yb-n2 | 163 | 1 |
| yb-n3 | 162 | 1 |

## Fixture and selectivity mapping

- arm: `study01-person-donation-charity`
- fact rows: 200000
- parent rows: 2000
- rank thresholds: [20 2000 180000]

The predicate is `rank <= threshold`, where `rank` is a deterministic shuffle of 1..200000; `rank` is a labelled proxy for the teaching selective predicate, not the same question. The `cte`/`cte-small` arm is the synthetic customers/facts microfixture; the `cte-domain` arm is Study 01's person/donation fixture with the person's charity projected and the views binding it to this catalogue.

## Per-engine feasibility evidence

- fact rows 200000, customers 2000, thresholds [20 2000 180000]
- expensive reuse token supported by the engine: yes

| threshold | rows with rank <= threshold |
|---|---|
| 180000 | 180000 |
| 20 | 20 |
| 2000 | 2000 |

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
| 20 (0.01%) | 64.2191 | missing | missing | missing | missing | missing | 61.9044 | 147.6074 | 66.5217 | missing |
| 2000 (1%) | 63.2172 | missing | missing | missing | missing | missing | 66.4139 | 147.2876 | 64.8700 | missing |
| 180000 (90%) | 171.7928 | missing | missing | missing | missing | missing | 166.0411 | 170.7248 | 131.4017 | missing |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | cte_joined_outer vs join_where | 64.2191 | 61.9044 | 0.9640 | false |
| 20 | cte_joined_outer_mat vs join_where | 64.2191 | 147.6074 | 2.2985 | false |
| 20 | cte_joined_outer_notmat vs join_where | 64.2191 | 66.5217 | 1.0359 | false |
| 2000 | cte_joined_outer vs join_where | 63.2172 | 66.4139 | 1.0506 | false |
| 2000 | cte_joined_outer_mat vs join_where | 63.2172 | 147.2876 | 2.3299 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 63.2172 | 64.8700 | 1.0261 | false |
| 180000 | cte_joined_outer vs join_where | 171.7928 | 166.0411 | 0.9665 | false |
| 180000 | cte_joined_outer_mat vs join_where | 171.7928 | 170.7248 | 0.9938 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 171.7928 | 131.4017 | 0.7649 | false |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_joined_outer_mat vs cte_joined_outer | 61.9044 | 147.6074 | 2.3844 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 61.9044 | 66.5217 | 1.0746 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 66.4139 | 147.2876 | 2.2177 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 66.4139 | 64.8700 | 0.9768 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 166.0411 | 170.7248 | 1.0282 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 166.0411 | 131.4017 | 0.7914 | false | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 66.1884 | 1 | `c54dd880faba877a` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 65.6087 | 1 | `cdee8e766ba6f835` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | missing | | | | |
| default | 2 | 100% | missing | | | | |
| default | 4 | 1% | 64.0627 | 4 | `454591938ad5ad60` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 68.8162 | 4 | `105a61a92d1220ad` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | missing | | | | |
| default | 8 | 100% | missing | | | | |
| default | 16 | 1% | 66.1929 | 16 | `3b6a2bef5c4b8c3b` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 65.1117 | 16 | `08c7792ad8e12faa` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 61.4513 | 1 | `3883734e89c731d9` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 62.6346 | 1 | `2674d8b64de7a1d8` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | missing | | | | |
| mat | 2 | 100% | missing | | | | |
| mat | 4 | 1% | 64.8848 | 4 | `33f3e68e1172ca78` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 60.8916 | 4 | `a6e788458d9df207` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | missing | | | | |
| mat | 8 | 100% | missing | | | | |
| mat | 16 | 1% | 62.6415 | 16 | `5864c87f7c85b2e6` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 64.2062 | 16 | `ad5fef7f4d3aca97` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 63.0788 | 1 | `d6b8b1f231391da8` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 67.4091 | 1 | `299b08c895622dcd` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | missing | | | | |
| notmat | 2 | 100% | missing | | | | |
| notmat | 4 | 1% | 284.0113 | 4 | `dba931842aeafeb7` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 174.6363 | 4 | `cecb6526f13b8b35` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | missing | | | | |
| notmat | 8 | 100% | missing | | | | |
| notmat | 16 | 1% | 892.2641 | 16 | `1d8f75d67b0a2116` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 637.0188 | 16 | `06a2a726c6e48fc0` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 0.9284 | false |
| notmat | 1 | 1% | 0.9530 | false |
| mat | 1 | 100% | 0.9547 | false |
| notmat | 1 | 100% | 1.0274 | false |
| mat | 4 | 1% | 1.0128 | false |
| notmat | 4 | 1% | 4.4333 | false |
| mat | 4 | 100% | 0.8848 | false |
| notmat | 4 | 100% | 2.5377 | false |
| mat | 16 | 1% | 0.9463 | false |
| notmat | 16 | 1% | 13.4798 | false |
| mat | 16 | 100% | 0.9861 | false |
| notmat | 16 | 100% | 9.7835 | false |

## Latency samples and derived throughput

_The reader is sequential, so `queries/s` is 1000/median_ms, a derived rate and not a concurrent-capacity test. min/max are the five sample extremes._

| axis | status | min ms | median ms | max ms | samples | derived queries/s |
|---|---|---|---|---|---|---|
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 160.0531 | 166.0411 | 171.5196 | 5 | 6.0 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 60.2569 | 61.9044 | 62.9021 | 5 | 16.2 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 64.1646 | 66.4139 | 70.2350 | 5 | 15.1 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 164.0012 | 170.7248 | 175.3141 | 5 | 5.9 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 140.0670 | 147.6074 | 162.7509 | 5 | 6.8 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 144.5180 | 147.2876 | 148.0057 | 5 | 6.8 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 127.7814 | 131.4017 | 138.5334 | 5 | 7.6 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 64.4909 | 66.5217 | 68.0618 | 5 | 15.0 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 62.4168 | 64.8700 | 138.3768 | 5 | 15.4 |
| family=bag;stmt=join_where;threshold=180000 | valid | 162.1672 | 171.7928 | 185.3226 | 5 | 5.8 |
| family=bag;stmt=join_where;threshold=20 | valid | 60.8097 | 64.2191 | 66.2106 | 5 | 15.6 |
| family=bag;stmt=join_where;threshold=2000 | valid | 61.0954 | 63.2172 | 67.5235 | 5 | 15.8 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 64.4057 | 66.1929 | 67.1320 | 5 | 15.1 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 63.1950 | 65.1117 | 68.3233 | 5 | 15.4 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 65.9415 | 66.1884 | 69.4012 | 5 | 15.1 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 63.0125 | 65.6087 | 66.6663 | 5 | 15.2 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 62.5058 | 64.0627 | 65.3607 | 5 | 15.6 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 66.2017 | 68.8162 | 72.9341 | 5 | 14.5 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 60.8047 | 62.6415 | 68.2421 | 5 | 16.0 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 61.3941 | 64.2062 | 68.7710 | 5 | 15.6 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 59.5900 | 61.4513 | 63.8385 | 5 | 16.3 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 58.5413 | 62.6346 | 68.4442 | 5 | 16.0 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 63.0185 | 64.8848 | 68.0038 | 5 | 15.4 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 58.5940 | 60.8916 | 64.9834 | 5 | 16.4 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 599.2572 | 892.2641 | 901.6146 | 5 | 1.1 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 623.4148 | 637.0188 | 682.5767 | 5 | 1.6 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 60.9016 | 63.0788 | 64.9308 | 5 | 15.9 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 64.8946 | 67.4091 | 106.6473 | 5 | 14.8 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 171.7181 | 284.0113 | 359.7373 | 5 | 3.5 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 172.3342 | 174.6363 | 198.0270 | 5 | 5.7 |

## Matched blended workload

- phase: `blended`; target offered rate: 10 updates/s
- offered 381, committed 381, failed 0 over 38.035 s
- offered rate 10.017/s, completed rate 10.017/s
- reader operations overlapping a committed write: 165 of 211
- reader retries of transient conflicts: 1
- per key: 381 keys, min 1, max 1 per key
- ledger: 381 keys reconciled, 0 mismatch
- update latency: avg 4.873 ms, min 2.663, max 89.624 over 381

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42774313`, 5 interleaved rounds
  1. family=bag;stmt=cte_joined_outer_mat;threshold=20
  2. family=reuse;policy=notmat;refs=1;consumer=100%
  3. family=reuse;policy=mat;refs=16;consumer=100%
  4. family=bag;stmt=cte_joined_outer;threshold=180000
  5. family=reuse;policy=mat;refs=1;consumer=100%
  6. family=reuse;policy=default;refs=1;consumer=100%
  7. family=reuse;policy=mat;refs=4;consumer=100%
  8. family=reuse;policy=notmat;refs=1;consumer=1%
  9. family=reuse;policy=notmat;refs=4;consumer=100%
  10. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  11. family=bag;stmt=join_where;threshold=20
  12. family=bag;stmt=cte_joined_outer;threshold=20
  13. family=reuse;policy=notmat;refs=16;consumer=100%
  14. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  15. family=reuse;policy=default;refs=16;consumer=1%
  16. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  17. family=bag;stmt=cte_joined_outer;threshold=2000
  18. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  19. family=reuse;policy=default;refs=1;consumer=1%
  20. family=reuse;policy=mat;refs=4;consumer=1%
  21. family=bag;stmt=join_where;threshold=2000
  22. family=reuse;policy=default;refs=4;consumer=1%
  23. family=reuse;policy=mat;refs=16;consumer=1%
  24. family=reuse;policy=mat;refs=1;consumer=1%
  25. family=reuse;policy=notmat;refs=4;consumer=1%
  26. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  27. family=reuse;policy=notmat;refs=16;consumer=1%
  28. family=reuse;policy=default;refs=4;consumer=100%
  29. family=bag;stmt=join_where;threshold=180000
  30. family=reuse;policy=default;refs=16;consumer=100%

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
| family=bag;stmt=cte_joined_outer;threshold=180000 | j05_cte_joined_outer | valid | `fa5dd3b3951b0b1a` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20 | j05_cte_joined_outer | valid | `a8c6c33eadd337b2` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=2000 | j05_cte_joined_outer | valid | `4a7536e86e43cb9c` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | j05m_cte_joined_outer_mat | valid | `b90e0bbd3958ea74` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | j05m_cte_joined_outer_mat | valid | `64b0f0a62e7aa081` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | j05m_cte_joined_outer_mat | valid | `b706264bd7c14ede` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | j05n_cte_joined_outer_notmat | valid | `4782831b0dde9347` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | j05n_cte_joined_outer_notmat | valid | `18f93e32d4c98c73` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | j05n_cte_joined_outer_notmat | valid | `906fc79ac4638534` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=join_where;threshold=180000 | j01_join_where | valid | `86d75c4d226858ef` |  |
| family=bag;stmt=join_where;threshold=20 | j01_join_where | valid | `b60ef9915d053748` |  |
| family=bag;stmt=join_where;threshold=2000 | j01_join_where | valid | `082d51da3b1c9a82` |  |
| family=reuse;policy=default;refs=16;consumer=1% | r_default_n16_c1 | valid | `3b6a2bef5c4b8c3b` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=16;consumer=100% | r_default_n16_c100 | valid | `08c7792ad8e12faa` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=1;consumer=1% | r_default_n1_c1 | valid | `c54dd880faba877a` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=1;consumer=100% | r_default_n1_c100 | valid | `cdee8e766ba6f835` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=4;consumer=1% | r_default_n4_c1 | valid | `454591938ad5ad60` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=100% | r_default_n4_c100 | valid | `105a61a92d1220ad` | policy=default; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=1% | r_mat_n16_c1 | valid | `5864c87f7c85b2e6` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=100% | r_mat_n16_c100 | valid | `ad5fef7f4d3aca97` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=1% | r_mat_n1_c1 | valid | `3883734e89c731d9` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=100% | r_mat_n1_c100 | valid | `2674d8b64de7a1d8` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=1% | r_mat_n4_c1 | valid | `33f3e68e1172ca78` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=100% | r_mat_n4_c100 | valid | `a6e788458d9df207` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=notmat;refs=16;consumer=1% | r_notmat_n16_c1 | valid | `1d8f75d67b0a2116` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `06a2a726c6e48fc0` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `d6b8b1f231391da8` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `299b08c895622dcd` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `dba931842aeafeb7` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `cecb6526f13b8b35` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, DIST, DEBUG, FORMAT TEXT` supported on yugabyte
- one host, single client: no independent hardware, real-network or equal-total-budget comparison
- no concurrent writers: the isolated read contract only; a matched blended read/write phase is a recorded gap
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
