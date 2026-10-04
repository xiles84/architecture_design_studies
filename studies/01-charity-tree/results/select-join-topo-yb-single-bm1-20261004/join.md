# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 30 cells, 29 valid, engine `yugabyte`, topology `single`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-topo-yb-single-bm1-20261004`
- run tag: `run/study-01/select-join-topo-yb-single-bm1-20261004`
- environment: `host-zenbook-ux5406sa`
- repo commit: `fa5bbf57ddaa27d90fc5568018a21f3efa43e0ff`
- engine version: `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit`
- database engine image: `docker.io/yugabytedb/yugabyte:2025.2.6.0-b111` @ `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`
- benchmark client image: `localhost/charitybench-select:select-join-topo-yb-single-bm1-20261004` @ `741ebaefc42d21402a7c6cd2ebdd3f51d20588ca5655c76f348d18ae8cf85a62`
- inputs digest: `345c96f5dc35717c7c80bae2855d61e75687a03e10c563714b28daf7df109665`
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
- setting `phase`: blended
- setting `placement`: noncolocated
- setting `prepared_plan_policy`: pgx simple protocol; no prepared statements
- setting `scale`: cte-domain
- setting `selectivity_mapping`: rank <= threshold; rank is a deterministic shuffle of 1..fact_count
- setting `statement_budget_ms`: 120000
- setting `statistics`: ANALYZE after deterministic load
- setting `thresholds`: [20 2000 180000]
- setting `topology`: single
- setting `trials_per_cell`: 5
- setting `whole_run_count`: 1; no replicated winner in this file
- setting `writer_keys`: 1000
- setting `writer_target_rate_per_sec`: 10

### Endpoint distribution

A cluster arm must rotate its reader across nodes; a single seed endpoint would not be balanced cluster access.

| endpoint | operations | connections at end |
|---|---|---|
| yb-single | 486 | 2 |

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
| 20 (0.01%) | 57.0961 | missing | missing | missing | missing | missing | 53.6807 | 94.5046 | 53.3826 | missing |
| 2000 (1%) | 55.9449 | missing | missing | missing | missing | missing | 56.4526 | 88.9707 | 57.3047 | missing |
| 180000 (90%) | 110.7549 | missing | missing | missing | missing | missing | 117.2211 | 121.8806 | 117.8396 | missing |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | cte_joined_outer vs join_where | 57.0961 | 53.6807 | 0.9402 | false |
| 20 | cte_joined_outer_mat vs join_where | 57.0961 | 94.5046 | 1.6552 | false |
| 20 | cte_joined_outer_notmat vs join_where | 57.0961 | 53.3826 | 0.9350 | false |
| 2000 | cte_joined_outer vs join_where | 55.9449 | 56.4526 | 1.0091 | false |
| 2000 | cte_joined_outer_mat vs join_where | 55.9449 | 88.9707 | 1.5903 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 55.9449 | 57.3047 | 1.0243 | false |
| 180000 | cte_joined_outer vs join_where | 110.7549 | 117.2211 | 1.0584 | false |
| 180000 | cte_joined_outer_mat vs join_where | 110.7549 | 121.8806 | 1.1005 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 110.7549 | 117.8396 | 1.0640 | false |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_joined_outer_mat vs cte_joined_outer | 53.6807 | 94.5046 | 1.7605 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 53.6807 | 53.3826 | 0.9944 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 56.4526 | 88.9707 | 1.5760 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 56.4526 | 57.3047 | 1.0151 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 117.2211 | 121.8806 | 1.0397 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 117.2211 | 117.8396 | 1.0053 | false | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 55.5610 | 1 | `bb2fe1ff941c1eaa` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 56.9112 | 1 | `077292fc192eac2e` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | missing | | | | |
| default | 2 | 100% | missing | | | | |
| default | 4 | 1% | 53.6118 | 4 | `b1ee2091e6d36c3d` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 54.5469 | 4 | `995bd3be46b9aacd` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | missing | | | | |
| default | 8 | 100% | missing | | | | |
| default | 16 | 1% | 55.5421 | 16 | `cb446b89f78614ac` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 55.0792 | 16 | `e38c5e12fa1aea4d` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 56.1926 | 1 | `693e30efc0253b64` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 54.0045 | 1 | `067eafcd46c5531b` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | missing | | | | |
| mat | 2 | 100% | missing | | | | |
| mat | 4 | 1% | 56.4052 | 4 | `976e5a17b733db5a` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 53.8064 | 4 | `88c5ad0de1f73c2a` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | missing | | | | |
| mat | 8 | 100% | missing | | | | |
| mat | 16 | 1% | 53.7583 | 16 | `3cc6e07b49d806ba` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 59.4919 | 16 | `d61b1aaaecb87be7` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 58.8231 | 1 | `c1daf6372e4ff42e` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 58.1510 | 1 | `d2f5eadcc57a5b6e` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | missing | | | | |
| notmat | 2 | 100% | missing | | | | |
| notmat | 4 | 1% | error | 0 | `e3b0c44298fc1c14` |  | error |
| notmat | 4 | 100% | 167.2247 | 4 | `7a1407437f2688fb` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | missing | | | | |
| notmat | 8 | 100% | missing | | | | |
| notmat | 16 | 1% | 605.8273 | 16 | `35c82a5783e7c7bf` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 624.8906 | 16 | `5393c7bcac44bae9` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 1.0114 | false |
| notmat | 1 | 1% | 1.0587 | false |
| mat | 1 | 100% | 0.9489 | false |
| notmat | 1 | 100% | 1.0218 | false |
| mat | 4 | 1% | 1.0521 | false |
| notmat | 4 | 1% | N/A | N/A |
| mat | 4 | 100% | 0.9864 | false |
| notmat | 4 | 100% | 3.0657 | false |
| mat | 16 | 1% | 0.9679 | false |
| notmat | 16 | 1% | 10.9075 | false |
| mat | 16 | 100% | 1.0801 | false |
| notmat | 16 | 100% | 11.3453 | false |

## Latency samples and derived throughput

_The reader is sequential, so `queries/s` is 1000/median_ms, a derived rate and not a concurrent-capacity test. min/max are the five sample extremes._

| axis | status | min ms | median ms | max ms | samples | derived queries/s |
|---|---|---|---|---|---|---|
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 114.7138 | 117.2211 | 155.8645 | 5 | 8.5 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 52.1748 | 53.6807 | 59.8353 | 5 | 18.6 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 54.9858 | 56.4526 | 97.7375 | 5 | 17.7 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 104.6433 | 121.8806 | 168.4263 | 5 | 8.2 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 81.5708 | 94.5046 | 101.7898 | 5 | 10.6 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 83.1625 | 88.9707 | 149.6489 | 5 | 11.2 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 102.9430 | 117.8396 | 137.1514 | 5 | 8.5 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 52.1307 | 53.3826 | 78.4183 | 5 | 18.7 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 52.8274 | 57.3047 | 78.0251 | 5 | 17.5 |
| family=bag;stmt=join_where;threshold=180000 | valid | 108.6359 | 110.7549 | 121.9139 | 5 | 9.0 |
| family=bag;stmt=join_where;threshold=20 | valid | 52.2001 | 57.0961 | 63.4345 | 5 | 17.5 |
| family=bag;stmt=join_where;threshold=2000 | valid | 53.6823 | 55.9449 | 63.2233 | 5 | 17.9 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 51.4050 | 55.5421 | 85.4321 | 5 | 18.0 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 53.1745 | 55.0792 | 57.6074 | 5 | 18.2 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 54.4575 | 55.5610 | 144.1732 | 5 | 18.0 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 55.0423 | 56.9112 | 79.7483 | 5 | 17.6 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 51.8284 | 53.6118 | 88.7719 | 5 | 18.7 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 51.5966 | 54.5469 | 56.5456 | 5 | 18.3 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 53.3031 | 53.7583 | 60.3848 | 5 | 18.6 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 56.5767 | 59.4919 | 62.4838 | 5 | 16.8 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 51.6239 | 56.1926 | 71.9166 | 5 | 17.8 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 51.0819 | 54.0045 | 85.4147 | 5 | 18.5 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 54.0015 | 56.4052 | 82.7538 | 5 | 17.7 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 50.6294 | 53.8064 | 63.1227 | 5 | 18.6 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 576.4201 | 605.8273 | 899.9466 | 5 | 1.7 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 590.1039 | 624.8906 | 628.1930 | 5 | 1.6 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 51.6655 | 58.8231 | 82.2801 | 5 | 17.0 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 49.6007 | 58.1510 | 62.6706 | 5 | 17.2 |
| family=reuse;policy=notmat;refs=4;consumer=1% | error | N/A | N/A | N/A | 0 | N/A |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 149.0915 | 167.2247 | 196.4165 | 5 | 6.0 |

## Matched blended workload

- phase: `blended`; target offered rate: 10 updates/s
- offered 324, committed 324, failed 0 over 32.362 s
- offered rate 10.012/s, completed rate 10.012/s
- reader operations overlapping a committed write: 145 of 209
- reader retries of transient conflicts: 4
- per key: 324 keys, min 1, max 1 per key
- ledger: 324 keys reconciled, 0 mismatch
- update latency: avg 3.410 ms, min 1.695, max 56.631 over 324

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42495933`, 5 interleaved rounds
  1. family=reuse;policy=mat;refs=4;consumer=100%
  2. family=reuse;policy=default;refs=16;consumer=100%
  3. family=reuse;policy=mat;refs=16;consumer=1%
  4. family=reuse;policy=notmat;refs=16;consumer=100%
  5. family=reuse;policy=mat;refs=1;consumer=1%
  6. family=bag;stmt=cte_joined_outer;threshold=180000
  7. family=reuse;policy=notmat;refs=4;consumer=1%
  8. family=reuse;policy=notmat;refs=16;consumer=1%
  9. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  10. family=bag;stmt=cte_joined_outer_mat;threshold=20
  11. family=bag;stmt=join_where;threshold=20
  12. family=reuse;policy=mat;refs=16;consumer=100%
  13. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  14. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  15. family=reuse;policy=notmat;refs=1;consumer=1%
  16. family=reuse;policy=default;refs=16;consumer=1%
  17. family=reuse;policy=default;refs=4;consumer=1%
  18. family=reuse;policy=mat;refs=4;consumer=1%
  19. family=bag;stmt=cte_joined_outer;threshold=2000
  20. family=reuse;policy=default;refs=1;consumer=100%
  21. family=reuse;policy=mat;refs=1;consumer=100%
  22. family=reuse;policy=default;refs=1;consumer=1%
  23. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  24. family=bag;stmt=cte_joined_outer;threshold=20
  25. family=bag;stmt=join_where;threshold=180000
  26. family=reuse;policy=notmat;refs=4;consumer=100%
  27. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  28. family=reuse;policy=notmat;refs=1;consumer=100%
  29. family=reuse;policy=default;refs=4;consumer=100%
  30. family=bag;stmt=join_where;threshold=2000

## Live negative controls

| control | pretended contract | detected | rows | note |
|---|---|---|---|---|
| x01_optional_where | sel_join_optional | true | 1 |  |
| x02_join_limit_before | sel_join_bag | true | 7 |  |

## Excluded cells

- `family=reuse;policy=notmat;refs=4;consumer=1%` (r_notmat_n4_c1) — **error** — ERROR: Restart read required (SQLSTATE 40001)

A cell that is invalid, unsupported or timed out is never ranked.

## Plan signatures

| axis | statement | status | structure signature | policy note |
|---|---|---|---|---|
| family=bag;stmt=cte_joined_outer;threshold=180000 | j05_cte_joined_outer | valid | `8fc292d391d08938` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20 | j05_cte_joined_outer | valid | `3774da6c4eea961a` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=2000 | j05_cte_joined_outer | valid | `58c1844e2fe5437b` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | j05m_cte_joined_outer_mat | valid | `5adf2df601819d27` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | j05m_cte_joined_outer_mat | valid | `6746b6f6bba3d0d3` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | j05m_cte_joined_outer_mat | valid | `c06b4c9aba805b71` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | j05n_cte_joined_outer_notmat | valid | `f8e1b25ca599ce59` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | j05n_cte_joined_outer_notmat | valid | `3154d8205eb04fe5` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | j05n_cte_joined_outer_notmat | valid | `c2602cf06d21e693` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=join_where;threshold=180000 | j01_join_where | valid | `61b2ff6897427eec` |  |
| family=bag;stmt=join_where;threshold=20 | j01_join_where | valid | `0e515fb6a2d9d5bd` |  |
| family=bag;stmt=join_where;threshold=2000 | j01_join_where | valid | `7cbfbc13e10ca3e3` |  |
| family=reuse;policy=default;refs=16;consumer=1% | r_default_n16_c1 | valid | `cb446b89f78614ac` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=16;consumer=100% | r_default_n16_c100 | valid | `e38c5e12fa1aea4d` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=1;consumer=1% | r_default_n1_c1 | valid | `bb2fe1ff941c1eaa` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=1;consumer=100% | r_default_n1_c100 | valid | `077292fc192eac2e` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=4;consumer=1% | r_default_n4_c1 | valid | `b1ee2091e6d36c3d` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=100% | r_default_n4_c100 | valid | `995bd3be46b9aacd` | policy=default; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=1% | r_mat_n16_c1 | valid | `3cc6e07b49d806ba` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=100% | r_mat_n16_c100 | valid | `d61b1aaaecb87be7` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=1% | r_mat_n1_c1 | valid | `693e30efc0253b64` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=100% | r_mat_n1_c100 | valid | `067eafcd46c5531b` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=1% | r_mat_n4_c1 | valid | `976e5a17b733db5a` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=100% | r_mat_n4_c100 | valid | `88c5ad0de1f73c2a` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=notmat;refs=16;consumer=1% | r_notmat_n16_c1 | valid | `35c82a5783e7c7bf` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `5393c7bcac44bae9` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `c1daf6372e4ff42e` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `d2f5eadcc57a5b6e` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | error | `e3b0c44298fc1c14` |  |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `7a1407437f2688fb` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, DIST, DEBUG, FORMAT TEXT` supported on yugabyte
- one host, single client: no independent hardware, real-network or equal-total-budget comparison
- no concurrent writers: the isolated read contract only; a matched blended read/write phase is a recorded gap
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
