# Study 04 — configuration portal — run `devcheck-sub`

Generated 2026-09-21T11:49:55Z from 4 result file(s). This report contains measurements only.

| Field | Value |
|---|---|
| Study | 04-configuration-portal |
| Run id | `devcheck-sub` |
| Environment | `host-zenbook-ux5406sa` |
| Repository commit | `2e959908963509d26a8f6a9cabdded595b8fafbe` |
| `git describe` | `repo/open-loop-arrivals-dirty` |
| Working tree dirty at start | true |
| **Inputs digest** | `58e3f53cbfb137e0` |
| Cells | 4 |

## TL;DR — measured facts, selected by fixed rules

- Cells reported: 4; cells that failed: 2.
  - failed: `pg-single/x1_lost_update_control` (design x1_lost_update_control statement a_version_sum uses $n but declares no -- params)
  - failed: `pg-single/x2_rollup_drift_control` (phase contention: negative control x2_rollup_drift_control did not fire: 16 writers, 0 acknowledged, counter 0 (expected 0))
- Negative controls that fired: 1. Controls that did NOT fire: 1.
  - did not fire: `pg-single/x1_lost_update_control` — the associated correctness claim is **not demonstrated**
- Invariant violations recorded: 2.
  - `pg-single/x2_rollup_drift_control`: INV-8/INV-9: a_rollup_mismatches reports 6 mismatches
  - `pg-single/x2_rollup_drift_control`: INV-8/INV-9: a_rollup_mismatches reports 6 mismatches
- Highest and lowest valid read result, across all cells and reads: 50768.8 ops/s (`pg-single/n1/r03_revision_check`) and 17753.4 ops/s (`pg-single/n1/r04_list_installations`).

This list ranks by number only. Which design is *better* is not a measurement and is not stated here.

## Analyses of this data

No signed analysis exists yet for inputs digest `58e3f53cbfb137e0`.

## Cells

| Topology | Design | Family | Kind | Gate | Load s | Entries | Serialized B |
|---|---|---|---|---|---|---|---|
| pg-single | `d1_doc_row` | document | doc | pass | 0.00 | 48 | 5.3 KiB |
| pg-single | `n1_rows_indexed` | normalized | rows | pass | 0.00 | 48 | 5.3 KiB |
| pg-single | `x1_lost_update_control` | control | rows | aborted | 0.00 | 0 | — |
| pg-single | `x2_rollup_drift_control` | control | rows | aborted | 0.00 | 48 | 5.3 KiB |


## Dataset (identical for every design and engine)

| Property | Value |
|---|---|
| Seed | 42 |
| Product definitions | 10 |
| Environments | 3 |
| Business units | 4 |
| Installed products | 6 |
| Entries per installed product (tier) | 8 |
| Cardinality mode | exact |
| Value regime | small |
| Key bytes | 16 |
| Value bytes | 96 |
| Total entries | 48 |
| Total serialized bytes | 5.3 KiB |
| Mean serialized bytes per installation | 907 |

## Reads (isolated)

`ops/s` is the median across trials. A run with one trial has no spread and no error bar.

| Topology | Design | Read | ops/s | p50 ms | p99 ms | max ms | ops | errors | spread % |
|---|---|---|---|---|---|---|---|---|---|
| pg-single | `d1` | r01_effective_config | 33k | 0.13 | 0.43 | 41 | 13139 | 0 | 0.0 |
| pg-single | `d1` | r02_read_key | 44k | 0.11 | 0.38 | 36 | 17566 | 0 | 0.0 |
| pg-single | `d1` | r03_revision_check | 47k | 0.11 | 0.39 | 31 | 18921 | 0 | 0.0 |
| pg-single | `d1` | r04_list_installations | 19k | 0.18 | 0.62 | 57 | 7435 | 0 | 0.0 |
| pg-single | `d1` | r05_search_key | 32k | 0.13 | 0.43 | 43 | 12737 | 0 | 0.0 |
| pg-single | `n1` | r01_effective_config | 34k | 0.12 | 0.44 | 44 | 13561 | 0 | 0.0 |
| pg-single | `n1` | r02_read_key | 41k | 0.11 | 0.40 | 37 | 16239 | 0 | 0.0 |
| pg-single | `n1` | r03_revision_check | 51k | 0.10 | 0.35 | 27 | 20329 | 0 | 0.0 |
| pg-single | `n1` | r04_list_installations | 18k | 0.18 | 0.62 | 58 | 7107 | 0 | 0.0 |
| pg-single | `n1` | r05_search_key | 34k | 0.12 | 0.39 | 43 | 13778 | 0 | 0.0 |
| pg-single | `x2` | r01_effective_config | 33k | 0.13 | 0.42 | 43 | 13140 | 0 | 0.0 |
| pg-single | `x2` | r02_read_key | 41k | 0.11 | 0.37 | 38 | 16566 | 0 | 0.0 |
| pg-single | `x2` | r03_revision_check | 50k | 0.10 | 0.34 | 29 | 20057 | 0 | 0.0 |
| pg-single | `x2` | r04_list_installations | 25k | 0.14 | 0.47 | 48 | 11199 | 0 | 0.0 |
| pg-single | `x2` | r05_search_key | 34k | 0.13 | 0.40 | 42 | 13408 | 0 | 0.0 |


## Writes (isolated, one worker)

`ops/s` here is the cost of one update. Concurrency is a different question and is measured separately.

| Topology | Design | Write | ops/s | p50 ms | p99 ms | ops | errors | spread % |
|---|---|---|---|---|---|---|---|---|
| pg-single | `d1` | wd1_write_doc | 1.1k | 0.88 | 1.59 | 430 | 0 | 0.0 |
| pg-single | `d1` | w06_update_metadata | 1.6k | 0.59 | 1.20 | 642 | 0 | 0.0 |
| pg-single | `n1` | w01_modify_key | 1.3k | 0.72 | 1.38 | 531 | 0 | 0.0 |
| pg-single | `n1` | w02_add_key | 1.3k | 0.71 | 1.25 | 535 | 0 | 0.0 |
| pg-single | `n1` | w03_delete_key | 647 | 1.47 | 2.36 | 259 | 0 | 0.0 |
| pg-single | `n1` | w04_publish_batch | 1.1k | 0.85 | 1.51 | 439 | 0 | 0.0 |
| pg-single | `n1` | w05_replace_all | 945 | 0.96 | 1.60 | 378 | 0 | 0.0 |
| pg-single | `n1` | w06_update_metadata | 1.6k | 0.58 | 1.10 | 649 | 0 | 0.0 |
| pg-single | `x2` | w01_modify_key | 1.4k | 0.67 | 1.16 | 559 | 0 | 0.0 |
| pg-single | `x2` | w02_add_key | 1.3k | 0.74 | 1.32 | 509 | 0 | 0.0 |
| pg-single | `x2` | w03_delete_key | 685 | 1.39 | 2.10 | 274 | 0 | 0.0 |
| pg-single | `x2` | w04_publish_batch | 1.1k | 0.86 | 1.53 | 440 | 0 | 0.0 |
| pg-single | `x2` | w05_replace_all | 890 | 1.07 | 1.82 | 356 | 0 | 0.0 |
| pg-single | `x2` | w06_update_metadata | 1.3k | 0.67 | 1.42 | 531 | 0 | 0.0 |


## Contention (hot key, one installed product)

| Topology | Design | Strategy | Writers | Acked | Retries | Conflicts | ops/s | p99 ms | counter | expected | lost updates | control |
|---|---|---|---|---|---|---|---|---|---|---|---|---|
| pg-single | `d1` |  | 16 | 0 | 0 | 0 | — | — | 0 | 0 | 0 |  |
| pg-single | `n1` |  | 16 | 0 | 0 | 0 | — | — | 0 | 0 | 0 |  |
| pg-single | `x2` |  | 16 | 0 | 0 | 0 | — | — | 0 | 0 | 0 | negative control |


## Audits

| Topology | Design | Phase | Passed | Checks | Key/value mismatches | Revision mismatches | Design checks | First failure |
|---|---|---|---|---|---|---|---|---|
| pg-single | `d1` | 1 | true | 14 | 0 | 0 | a_content_hash_mismatches=0 a_document_revision_mismatches=0 |  |
| pg-single | `d1` | 2 | true | 14 | 0 | 0 | a_content_hash_mismatches=0 a_document_revision_mismatches=0 |  |
| pg-single | `n1` | 1 | true | 12 | 0 | 0 |  |  |
| pg-single | `n1` | 2 | true | 12 | 0 | 0 |  |  |
| pg-single | `x2` | 1 | false | 13 | 0 | 0 | a_rollup_mismatches=6 | INV-8/INV-9: a_rollup_mismatches reports 6 mismatches |
| pg-single | `x2` | 2 | false | 13 | 0 | 0 | a_rollup_mismatches=6 | INV-8/INV-9: a_rollup_mismatches reports 6 mismatches |


## Storage

PostgreSQL reports `pg_total_relation_size` per relation. The route to a YugabyteDB
size figure is a different query, so none is reported here rather than a number that
would be read as bytes on disk.

| Topology | Design | Relation | Bytes |
|---|---|---|---|
| pg-single | `d1` | business_unit | 32.0 KiB |
| pg-single | `d1` | config_document | 32.0 KiB |
| pg-single | `d1` | environment | 32.0 KiB |
| pg-single | `d1` | installed_product | 32.0 KiB |
| pg-single | `d1` | load_entry | 16.0 KiB |
| pg-single | `d1` | product_definition | 32.0 KiB |
| pg-single | `n1` | business_unit | 32.0 KiB |
| pg-single | `n1` | config_entry | 48.0 KiB |
| pg-single | `n1` | environment | 32.0 KiB |
| pg-single | `n1` | installed_product | 32.0 KiB |
| pg-single | `n1` | load_entry | 16.0 KiB |
| pg-single | `n1` | product_definition | 32.0 KiB |
| pg-single | `x2` | business_unit | 32.0 KiB |
| pg-single | `x2` | config_entry | 48.0 KiB |
| pg-single | `x2` | environment | 32.0 KiB |
| pg-single | `x2` | installed_product | 32.0 KiB |
| pg-single | `x2` | load_entry | 16.0 KiB |
| pg-single | `x2` | product_definition | 32.0 KiB |


## Limitations stated with the numbers

- The client and the database share the same eight heterogeneous cores; absolute
  throughput is lower than a two-machine setup would give, and relative comparisons are
  the reason the numbers are comparable at all.
- There is no real network. Distributed designs look better here than they would across
  availability zones, so `EXPLAIN (ANALYZE, DIST)` RPC counts are the portable signal.
- The read and write loops are closed-loop per operation; deep tails are optimistic
  floors (coordinated omission) and are compared between designs, never quoted as SLOs.
- A single trial has no error bar. Cells whose own trials disagree by more than about a
  fifth are flagged in the tables above.
- Plan capture was not part of this run's phases.

