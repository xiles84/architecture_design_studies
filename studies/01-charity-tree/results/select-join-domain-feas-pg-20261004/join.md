# Study 01 SELECT join/filter and CTE materialization — generated measurements

_Measurements only. Verdicts on P1–P4 belong in the signed analysis beside this report._

**TL;DR:** 0 cells, 0 valid, engine `postgres`, topology `single`, scale `cte-domain-small`, seed 42, 20000 fact rows.

## Provenance

- run id: `select-join-domain-feas-pg-20261004`
- run tag: `run/study-01/select-join-domain-feas-pg-20261004`
- environment: `host-zenbook-ux5406sa`
- repo commit: `1216bf0f0b8ca87f8aafb46b4e0605189a6d16a5`
- engine version: `PostgreSQL 17.11 (Debian 17.11-1.pgdg13+2) on x86_64-pc-linux-gnu, compiled by gcc (Debian 14.2.0-19) 14.2.0, 64-bit`
- database engine image: `docker.io/library/postgres:17.11` @ `7296f210ae81031ec955dbad9a67a84fe958572a2153b8d0826a647522904dc1`
- benchmark client image: `localhost/charitybench-select:select-join-domain-feas-pg-20261004` @ `6d86f9ae16cfca5f0b4ad1b752d44494c00ef6e71786f3a5bcaa20b4d5440b57`
- inputs digest: `9368238ae350a520970d28f21f73392c55956eef550a9535e9d51e3c5447ebcd`
- dataset fingerprint: `5c186ea19e202f5d`
- client CPU quota: `100000 100000`, memory: `2147483648`
- setting `client_workers`: 1 sequential reader; 0 writers
- setting `connection_pool`: 2 = 1 reader/setup + 1 reserve
- setting `customers`: 200
- setting `dataset`: study01-person-donation-charity
- setting `distributed_rpc`: N/A on PostgreSQL single-node
- setting `fact_count`: 20000
- setting `index_set`: child rank index and parent-key index, primary keys; no per-spelling index
- setting `prepared_plan_policy`: pgx simple protocol; no prepared statements
- setting `scale`: cte-domain-small
- setting `selectivity_mapping`: rank <= threshold; rank is a deterministic shuffle of 1..fact_count
- setting `statement_budget_ms`: 120000
- setting `statistics`: ANALYZE after deterministic load
- setting `thresholds`: [2 20 200 2000 10000 18000]
- setting `topology`: single
- setting `trials_per_cell`: 1
- setting `whole_run_count`: 1; no replicated winner in this file

## Fixture and selectivity mapping

- arm: `study01-person-donation-charity`
- fact rows: 20000
- parent rows: 200
- rank thresholds: [2 20 200 2000 10000 18000]

The predicate is `rank <= threshold`, where `rank` is a deterministic shuffle of 1..20000; `rank` is a labelled proxy for the teaching selective predicate, not the same question. The `cte`/`cte-small` arm is the synthetic customers/facts microfixture; the `cte-domain` arm is Study 01's person/donation fixture with the person's charity projected and the views binding it to this catalogue.

## Per-engine feasibility evidence

- fact rows 20000, customers 200, thresholds [2 20 200 2000 10000 18000]
- expensive reuse token supported by the engine: yes

| threshold | rows with rank <= threshold |
|---|---|
| 10000 | -1 |
| 18000 | -1 |
| 2 | -1 |
| 20 | -1 |
| 200 | -1 |
| 2000 | -1 |

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
- gap: threshold 2 count failed: ERROR: relation "selectjoin.facts" does not exist (SQLSTATE 42P01)
- gap: threshold 20 count failed: ERROR: relation "selectjoin.facts" does not exist (SQLSTATE 42P01)
- gap: threshold 200 count failed: ERROR: relation "selectjoin.facts" does not exist (SQLSTATE 42P01)
- gap: threshold 2000 count failed: ERROR: relation "selectjoin.facts" does not exist (SQLSTATE 42P01)
- gap: threshold 10000 count failed: ERROR: relation "selectjoin.facts" does not exist (SQLSTATE 42P01)
- gap: threshold 18000 count failed: ERROR: relation "selectjoin.facts" does not exist (SQLSTATE 42P01)

## P1 — six inner-join spellings and their CTE policies (median ms)

| threshold | join_where | join_on | filtered_derived | cte_filtered | cte_filtered_mat | cte_filtered_notmat | cte_joined_outer | cte_joined_outer_mat | cte_joined_outer_notmat | comma_join |
|---|---|---|---|---|---|---|---|---|---|---|
| 2 (0.01%) | missing | missing | missing | missing | missing | missing | missing | missing | missing | missing |
| 20 (0.1%) | missing | missing | missing | missing | missing | missing | missing | missing | missing | missing |
| 200 (1%) | missing | missing | missing | missing | missing | missing | missing | missing | missing | missing |
| 2000 (10%) | missing | missing | missing | missing | missing | missing | missing | missing | missing | missing |
| 10000 (50%) | missing | missing | missing | missing | missing | missing | missing | missing | missing | missing |
| 18000 (90%) | missing | missing | missing | missing | missing | missing | missing | missing | missing | missing |

## P1 — bag spelling pairs against the canonical join+WHERE spelling

| threshold | pair | A ms | B ms | B/A | same plan structure |
|---|---|---|---|---|---|

## P3 — forced CTE policy against the default spelling

| threshold | pair | default ms | forced ms | forced/default | same plan structure | observed policy |
|---|---|---|---|---|---|---|

## P4 — reuse ladder (median ms)

| policy | refs | consumer | median ms | checksum refs | plan structure | observed policy | status |
|---|---|---|---|---|---|---|---|
| default | 1 | 1% | missing | | | | |
| default | 1 | 100% | missing | | | | |
| default | 2 | 1% | missing | | | | |
| default | 2 | 100% | missing | | | | |
| default | 4 | 1% | missing | | | | |
| default | 4 | 100% | missing | | | | |
| default | 8 | 1% | missing | | | | |
| default | 8 | 100% | missing | | | | |
| default | 16 | 1% | missing | | | | |
| default | 16 | 100% | missing | | | | |
| mat | 1 | 1% | missing | | | | |
| mat | 1 | 100% | missing | | | | |
| mat | 2 | 1% | missing | | | | |
| mat | 2 | 100% | missing | | | | |
| mat | 4 | 1% | missing | | | | |
| mat | 4 | 100% | missing | | | | |
| mat | 8 | 1% | missing | | | | |
| mat | 8 | 100% | missing | | | | |
| mat | 16 | 1% | missing | | | | |
| mat | 16 | 100% | missing | | | | |
| notmat | 1 | 1% | missing | | | | |
| notmat | 1 | 100% | missing | | | | |
| notmat | 2 | 1% | missing | | | | |
| notmat | 2 | 100% | missing | | | | |
| notmat | 4 | 1% | missing | | | | |
| notmat | 4 | 100% | missing | | | | |
| notmat | 8 | 1% | missing | | | | |
| notmat | 8 | 100% | missing | | | | |
| notmat | 16 | 1% | missing | | | | |
| notmat | 16 | 100% | missing | | | | |

## P4 — reuse policy pairs (B/A)

| policy | refs | consumer | B/A | same plan structure |
|---|---|---|---|---|

## Latency samples and derived throughput

_The reader is sequential, so `queries/s` is 1000/median_ms, a derived rate and not a concurrent-capacity test. min/max are the five sample extremes._

| axis | status | min ms | median ms | max ms | samples | derived queries/s |
|---|---|---|---|---|---|---|

## Execution order

- order seed `0`, 0 interleaved rounds

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

## Recorded limits

- explain options `ANALYZE, VERBOSE, BUFFERS, FORMAT TEXT` supported on postgres
- one host, single client: no independent hardware, real-network or equal-total-budget comparison
- no concurrent writers: the isolated read contract only; a matched blended read/write phase is a recorded gap
- a flat cell reports a client or engine ceiling; it is not by itself a constant query cost
