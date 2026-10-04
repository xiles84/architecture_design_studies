# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 30 cells, 30 valid, engine `yugabyte`, topology `cluster3`, scale `cte-domain`, seed 42, 200000 fact rows.

## Provenance

- run id: `select-join-topo-yb-noncolocated-20261004`
- run tag: `run/study-01/select-join-topo-yb-noncolocated-20261004`
- environment: `host-zenbook-ux5406sa`
- repo commit: `6dc48392c17cccee175f15f5a18497497c863a95`
- engine version: `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit`
- database engine image: `docker.io/yugabytedb/yugabyte:2025.2.6.0-b111` @ `561d5e6d2f20738778f4b8e909c541b0dc3a48ccaf8e6015d4d87b5a1af37c35`
- benchmark client image: `localhost/charitybench-select:select-join-topo-yb-noncolocated-20261004` @ `f5b4872c872894055c1e2e06c2575e59674c056208158c9eb02e1fa997367421`
- inputs digest: `48751521d673f8ba1437125942c21f30c71144fe8c97d848db3d44d66c09789f`
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
- setting `phase`: isolated
- setting `placement`: noncolocated
- setting `prepared_plan_policy`: pgx simple protocol; no prepared statements
- setting `scale`: cte-domain
- setting `selectivity_mapping`: rank <= threshold; rank is a deterministic shuffle of 1..fact_count
- setting `statement_budget_ms`: 120000
- setting `statistics`: ANALYZE after deterministic load
- setting `thresholds`: [20 2000 180000]
- setting `topology`: cluster3
- setting `trials_per_cell`: 1
- setting `whole_run_count`: 1; no replicated winner in this file
- setting `writer_keys`: 1000
- setting `writer_target_rate_per_sec`: 10

### Endpoint distribution

A cluster arm must rotate its reader across nodes; a single seed endpoint would not be balanced cluster access.

| endpoint | operations | connections at end |
|---|---|---|
| yb-n1 | 123 | 1 |
| yb-n2 | 122 | 1 |
| yb-n3 | 122 | 1 |

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
| 20 (0.01%) | 52.5475 | missing | missing | missing | missing | missing | 52.5154 | 141.0527 | 53.3922 | missing |
| 2000 (1%) | 52.4955 | missing | missing | missing | missing | missing | 54.0584 | 114.9939 | 58.3214 | missing |
| 180000 (90%) | 166.5453 | missing | missing | missing | missing | missing | 156.9741 | 100.7194 | 103.1469 | missing |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|
| 20 | cte_joined_outer vs join_where | 52.5475 | 52.5154 | 0.9994 | false |
| 20 | cte_joined_outer_mat vs join_where | 52.5475 | 141.0527 | 2.6843 | false |
| 20 | cte_joined_outer_notmat vs join_where | 52.5475 | 53.3922 | 1.0161 | false |
| 2000 | cte_joined_outer vs join_where | 52.4955 | 54.0584 | 1.0298 | false |
| 2000 | cte_joined_outer_mat vs join_where | 52.4955 | 114.9939 | 2.1905 | false |
| 2000 | cte_joined_outer_notmat vs join_where | 52.4955 | 58.3214 | 1.1110 | false |
| 180000 | cte_joined_outer vs join_where | 166.5453 | 156.9741 | 0.9425 | false |
| 180000 | cte_joined_outer_mat vs join_where | 166.5453 | 100.7194 | 0.6048 | false |
| 180000 | cte_joined_outer_notmat vs join_where | 166.5453 | 103.1469 | 0.6193 | false |

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|
| 20 | cte_joined_outer_mat vs cte_joined_outer | 52.5154 | 141.0527 | 2.6859 | false | policy=materialized; observed_cte_scan=true |
| 20 | cte_joined_outer_notmat vs cte_joined_outer | 52.5154 | 53.3922 | 1.0167 | false | policy=not_materialized; observed_cte_scan=false |
| 2000 | cte_joined_outer_mat vs cte_joined_outer | 54.0584 | 114.9939 | 2.1272 | false | policy=materialized; observed_cte_scan=true |
| 2000 | cte_joined_outer_notmat vs cte_joined_outer | 54.0584 | 58.3214 | 1.0789 | false | policy=not_materialized; observed_cte_scan=false |
| 180000 | cte_joined_outer_mat vs cte_joined_outer | 156.9741 | 100.7194 | 0.6416 | false | policy=materialized; observed_cte_scan=true |
| 180000 | cte_joined_outer_notmat vs cte_joined_outer | 156.9741 | 103.1469 | 0.6571 | false | policy=not_materialized; observed_cte_scan=false |

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | 52.8944 | 1 | `055861f35e6b3d7a` | policy=default; observed_cte_scan=false | valid |
| default | 1 | 100% | 52.6356 | 1 | `3718fa067a4d01a9` | policy=default; observed_cte_scan=false | valid |
| default | 2 | 1% | missing | | | | |
| default | 2 | 100% | missing | | | | |
| default | 4 | 1% | 49.9601 | 4 | `5906d86f9db94a85` | policy=default; observed_cte_scan=true | valid |
| default | 4 | 100% | 51.2172 | 4 | `8b1ad9eb17098a31` | policy=default; observed_cte_scan=true | valid |
| default | 8 | 1% | missing | | | | |
| default | 8 | 100% | missing | | | | |
| default | 16 | 1% | 51.2400 | 16 | `2088cbff38588c2b` | policy=default; observed_cte_scan=true | valid |
| default | 16 | 100% | 63.7907 | 16 | `76ccd9847bab8374` | policy=default; observed_cte_scan=true | valid |
| mat | 1 | 1% | 52.2369 | 1 | `99406df3d603f303` | policy=materialized; observed_cte_scan=true | valid |
| mat | 1 | 100% | 57.0412 | 1 | `bdd96ce612206afe` | policy=materialized; observed_cte_scan=true | valid |
| mat | 2 | 1% | missing | | | | |
| mat | 2 | 100% | missing | | | | |
| mat | 4 | 1% | 53.2084 | 4 | `0166f20f11b9b0ac` | policy=materialized; observed_cte_scan=true | valid |
| mat | 4 | 100% | 51.9478 | 4 | `890101c11aa758d6` | policy=materialized; observed_cte_scan=true | valid |
| mat | 8 | 1% | missing | | | | |
| mat | 8 | 100% | missing | | | | |
| mat | 16 | 1% | 58.8631 | 16 | `34d7facc63a4d2ea` | policy=materialized; observed_cte_scan=true | valid |
| mat | 16 | 100% | 51.7907 | 16 | `728c2ccee454d2c6` | policy=materialized; observed_cte_scan=true | valid |
| notmat | 1 | 1% | 55.1288 | 1 | `6255cd60d0d6f241` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 1 | 100% | 56.5764 | 1 | `4ea88db4cbd6ce3a` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 2 | 1% | missing | | | | |
| notmat | 2 | 100% | missing | | | | |
| notmat | 4 | 1% | 141.6509 | 4 | `331149b87df89d69` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 4 | 100% | 126.0177 | 4 | `239795c2ed6bd240` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 8 | 1% | missing | | | | |
| notmat | 8 | 100% | missing | | | | |
| notmat | 16 | 1% | 456.8403 | 16 | `af76738f46d19c91` | policy=not_materialized; observed_cte_scan=false | valid |
| notmat | 16 | 100% | 506.5164 | 16 | `47f79f83a3676aa7` | policy=not_materialized; observed_cte_scan=false | valid |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|
| mat | 1 | 1% | 0.9876 | false |
| notmat | 1 | 1% | 1.0422 | false |
| mat | 1 | 100% | 1.0837 | false |
| notmat | 1 | 100% | 1.0749 | false |
| mat | 4 | 1% | 1.0650 | false |
| notmat | 4 | 1% | 2.8353 | false |
| mat | 4 | 100% | 1.0143 | false |
| notmat | 4 | 100% | 2.4605 | false |
| mat | 16 | 1% | 1.1488 | false |
| notmat | 16 | 1% | 8.9157 | false |
| mat | 16 | 100% | 0.8119 | false |
| notmat | 16 | 100% | 7.9403 | false |

## Latency samples and derived throughput

_The reader is sequential, so `queries/s` is 1000/median_ms, a derived rate and not a concurrent-capacity test. min/max are the five sample extremes._

| axis | status | min ms | median ms | max ms | samples | derived queries/s |
|---|---|---|---|---|---|---|
| family=bag;stmt=cte_joined_outer;threshold=180000 | valid | 156.9741 | 156.9741 | 156.9741 | 1 | 6.4 |
| family=bag;stmt=cte_joined_outer;threshold=20 | valid | 52.5154 | 52.5154 | 52.5154 | 1 | 19.0 |
| family=bag;stmt=cte_joined_outer;threshold=2000 | valid | 54.0584 | 54.0584 | 54.0584 | 1 | 18.5 |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | valid | 100.7194 | 100.7194 | 100.7194 | 1 | 9.9 |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | valid | 141.0527 | 141.0527 | 141.0527 | 1 | 7.1 |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | valid | 114.9939 | 114.9939 | 114.9939 | 1 | 8.7 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | valid | 103.1469 | 103.1469 | 103.1469 | 1 | 9.7 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | valid | 53.3922 | 53.3922 | 53.3922 | 1 | 18.7 |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | valid | 58.3214 | 58.3214 | 58.3214 | 1 | 17.1 |
| family=bag;stmt=join_where;threshold=180000 | valid | 166.5453 | 166.5453 | 166.5453 | 1 | 6.0 |
| family=bag;stmt=join_where;threshold=20 | valid | 52.5475 | 52.5475 | 52.5475 | 1 | 19.0 |
| family=bag;stmt=join_where;threshold=2000 | valid | 52.4955 | 52.4955 | 52.4955 | 1 | 19.0 |
| family=reuse;policy=default;refs=16;consumer=1% | valid | 51.2400 | 51.2400 | 51.2400 | 1 | 19.5 |
| family=reuse;policy=default;refs=16;consumer=100% | valid | 63.7907 | 63.7907 | 63.7907 | 1 | 15.7 |
| family=reuse;policy=default;refs=1;consumer=1% | valid | 52.8944 | 52.8944 | 52.8944 | 1 | 18.9 |
| family=reuse;policy=default;refs=1;consumer=100% | valid | 52.6356 | 52.6356 | 52.6356 | 1 | 19.0 |
| family=reuse;policy=default;refs=4;consumer=1% | valid | 49.9601 | 49.9601 | 49.9601 | 1 | 20.0 |
| family=reuse;policy=default;refs=4;consumer=100% | valid | 51.2172 | 51.2172 | 51.2172 | 1 | 19.5 |
| family=reuse;policy=mat;refs=16;consumer=1% | valid | 58.8631 | 58.8631 | 58.8631 | 1 | 17.0 |
| family=reuse;policy=mat;refs=16;consumer=100% | valid | 51.7907 | 51.7907 | 51.7907 | 1 | 19.3 |
| family=reuse;policy=mat;refs=1;consumer=1% | valid | 52.2369 | 52.2369 | 52.2369 | 1 | 19.1 |
| family=reuse;policy=mat;refs=1;consumer=100% | valid | 57.0412 | 57.0412 | 57.0412 | 1 | 17.5 |
| family=reuse;policy=mat;refs=4;consumer=1% | valid | 53.2084 | 53.2084 | 53.2084 | 1 | 18.8 |
| family=reuse;policy=mat;refs=4;consumer=100% | valid | 51.9478 | 51.9478 | 51.9478 | 1 | 19.3 |
| family=reuse;policy=notmat;refs=16;consumer=1% | valid | 456.8403 | 456.8403 | 456.8403 | 1 | 2.2 |
| family=reuse;policy=notmat;refs=16;consumer=100% | valid | 506.5164 | 506.5164 | 506.5164 | 1 | 2.0 |
| family=reuse;policy=notmat;refs=1;consumer=1% | valid | 55.1288 | 55.1288 | 55.1288 | 1 | 18.1 |
| family=reuse;policy=notmat;refs=1;consumer=100% | valid | 56.5764 | 56.5764 | 56.5764 | 1 | 17.7 |
| family=reuse;policy=notmat;refs=4;consumer=1% | valid | 141.6509 | 141.6509 | 141.6509 | 1 | 7.1 |
| family=reuse;policy=notmat;refs=4;consumer=100% | valid | 126.0177 | 126.0177 | 126.0177 | 1 | 7.9 |

## Matched blended workload

- phase: `isolated`; target offered rate: 10 updates/s
- offered 0, committed 0, failed 0 over 0.000 s
- offered rate 0.000/s, completed rate 0.000/s
- per key: no keys touched
- ledger: 0 keys reconciled, 0 mismatch

A phase whose reader operations never overlapped a committed write is not a concurrent measurement. The counter is outside the projection, the rank/index set and the MD5 token, so read answers are unchanged; the ledger is the independent check that no update was lost or duplicated.

## Execution order

- order seed `42843277`, 1 interleaved rounds
  1. family=reuse;policy=mat;refs=16;consumer=1%
  2. family=bag;stmt=join_where;threshold=180000
  3. family=bag;stmt=cte_joined_outer_mat;threshold=20
  4. family=reuse;policy=notmat;refs=1;consumer=1%
  5. family=reuse;policy=mat;refs=1;consumer=100%
  6. family=reuse;policy=default;refs=16;consumer=100%
  7. family=bag;stmt=cte_joined_outer_notmat;threshold=2000
  8. family=reuse;policy=notmat;refs=16;consumer=1%
  9. family=bag;stmt=cte_joined_outer;threshold=180000
  10. family=bag;stmt=cte_joined_outer_mat;threshold=180000
  11. family=reuse;policy=notmat;refs=16;consumer=100%
  12. family=bag;stmt=cte_joined_outer_notmat;threshold=20
  13. family=bag;stmt=cte_joined_outer;threshold=2000
  14. family=bag;stmt=cte_joined_outer_mat;threshold=2000
  15. family=reuse;policy=mat;refs=4;consumer=100%
  16. family=bag;stmt=join_where;threshold=2000
  17. family=bag;stmt=cte_joined_outer;threshold=20
  18. family=reuse;policy=notmat;refs=4;consumer=1%
  19. family=bag;stmt=join_where;threshold=20
  20. family=reuse;policy=default;refs=1;consumer=100%
  21. family=reuse;policy=mat;refs=4;consumer=1%
  22. family=reuse;policy=default;refs=4;consumer=100%
  23. family=reuse;policy=notmat;refs=4;consumer=100%
  24. family=reuse;policy=default;refs=16;consumer=1%
  25. family=bag;stmt=cte_joined_outer_notmat;threshold=180000
  26. family=reuse;policy=mat;refs=1;consumer=1%
  27. family=reuse;policy=default;refs=4;consumer=1%
  28. family=reuse;policy=mat;refs=16;consumer=100%
  29. family=reuse;policy=default;refs=1;consumer=1%
  30. family=reuse;policy=notmat;refs=1;consumer=100%

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
| family=bag;stmt=cte_joined_outer;threshold=180000 | j05_cte_joined_outer | valid | `fd6304cd49016dcf` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=20 | j05_cte_joined_outer | valid | `6b9b45e7ad508e90` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer;threshold=2000 | j05_cte_joined_outer | valid | `54b6ec73c2253c17` | policy=default; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_mat;threshold=180000 | j05m_cte_joined_outer_mat | valid | `fcebf6ba3d106c61` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=20 | j05m_cte_joined_outer_mat | valid | `dea06fe3d4212581` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_mat;threshold=2000 | j05m_cte_joined_outer_mat | valid | `d97e5141d3201934` | policy=materialized; observed_cte_scan=true |
| family=bag;stmt=cte_joined_outer_notmat;threshold=180000 | j05n_cte_joined_outer_notmat | valid | `9ebd8f516c3ee149` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=20 | j05n_cte_joined_outer_notmat | valid | `3af928a634556127` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=cte_joined_outer_notmat;threshold=2000 | j05n_cte_joined_outer_notmat | valid | `98781a77b58e8c5c` | policy=not_materialized; observed_cte_scan=false |
| family=bag;stmt=join_where;threshold=180000 | j01_join_where | valid | `c5fcc2046a066473` |  |
| family=bag;stmt=join_where;threshold=20 | j01_join_where | valid | `60c68429fffec056` |  |
| family=bag;stmt=join_where;threshold=2000 | j01_join_where | valid | `fa118f6a5ebf8316` |  |
| family=reuse;policy=default;refs=16;consumer=1% | r_default_n16_c1 | valid | `2088cbff38588c2b` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=16;consumer=100% | r_default_n16_c100 | valid | `76ccd9847bab8374` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=1;consumer=1% | r_default_n1_c1 | valid | `055861f35e6b3d7a` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=1;consumer=100% | r_default_n1_c100 | valid | `3718fa067a4d01a9` | policy=default; observed_cte_scan=false |
| family=reuse;policy=default;refs=4;consumer=1% | r_default_n4_c1 | valid | `5906d86f9db94a85` | policy=default; observed_cte_scan=true |
| family=reuse;policy=default;refs=4;consumer=100% | r_default_n4_c100 | valid | `8b1ad9eb17098a31` | policy=default; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=1% | r_mat_n16_c1 | valid | `34d7facc63a4d2ea` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=16;consumer=100% | r_mat_n16_c100 | valid | `728c2ccee454d2c6` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=1% | r_mat_n1_c1 | valid | `99406df3d603f303` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=1;consumer=100% | r_mat_n1_c100 | valid | `bdd96ce612206afe` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=1% | r_mat_n4_c1 | valid | `0166f20f11b9b0ac` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=mat;refs=4;consumer=100% | r_mat_n4_c100 | valid | `890101c11aa758d6` | policy=materialized; observed_cte_scan=true |
| family=reuse;policy=notmat;refs=16;consumer=1% | r_notmat_n16_c1 | valid | `af76738f46d19c91` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=16;consumer=100% | r_notmat_n16_c100 | valid | `47f79f83a3676aa7` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=1% | r_notmat_n1_c1 | valid | `6255cd60d0d6f241` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=1;consumer=100% | r_notmat_n1_c100 | valid | `4ea88db4cbd6ce3a` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=1% | r_notmat_n4_c1 | valid | `331149b87df89d69` | policy=not_materialized; observed_cte_scan=false |
| family=reuse;policy=notmat;refs=4;consumer=100% | r_notmat_n4_c100 | valid | `239795c2ed6bd240` | policy=not_materialized; observed_cte_scan=false |

## Recorded limits

- explain options `ANALYZE, DIST, DEBUG, FORMAT TEXT` supported on yugabyte
- one host, single client: no independent hardware, real-network or equal-total-budget comparison
- no concurrent writers: the isolated read contract only; a matched blended read/write phase is a recorded gap
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
