# Study 04 — configuration portal — run `devcheck-n1`

Generated 2026-09-21T11:46:55Z from 1 result file(s). This report contains measurements only.

| Field | Value |
|---|---|
| Study | 04-configuration-portal |
| Run id | `devcheck-n1` |
| Environment | `host-zenbook-ux5406sa` |
| Repository commit | `2e959908963509d26a8f6a9cabdded595b8fafbe` |
| `git describe` | `repo/open-loop-arrivals-dirty` |
| Working tree dirty at start | true |
| **Inputs digest** | `76db9ac5535dedc6` |
| Cells | 1 |

## TL;DR — measured facts, selected by fixed rules

- Cells reported: 1; cells that failed: 0.
- Negative controls that fired: 0. Controls that did NOT fire: 0.
- No invariant violation was recorded by a design that is not a negative control.
- Highest and lowest valid read result, across all cells and reads: 49571.6 ops/s (`pg-single/n1/r03_revision_check`) and 4412.0 ops/s (`pg-single/n1/r05_search_key`).

This list ranks by number only. Which design is *better* is not a measurement and is not stated here.

## Analyses of this data

No signed analysis exists yet for inputs digest `76db9ac5535dedc6`.

## Cells

| Topology | Design | Family | Kind | Gate | Load s | Entries | Serialized B |
|---|---|---|---|---|---|---|---|
| pg-single | `n1_rows_indexed` | normalized | rows | pass | 0.00 | 48 | 5.3 KiB |


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
| pg-single | `n1` | r01_effective_config | 5.6k | 0.12 | 0.39 | 39 | 2226 | 0 | 0.0 |
| pg-single | `n1` | r02_read_key | 6.0k | 0.11 | 0.35 | 35 | 2390 | 0 | 0.0 |
| pg-single | `n1` | r03_revision_check | 50k | 0.11 | 0.44 | 19 | 19838 | 0 | 0.0 |
| pg-single | `n1` | r04_list_installations | 19k | 0.18 | 0.62 | 55 | 7663 | 0 | 0.0 |
| pg-single | `n1` | r05_search_key | 4.4k | 0.13 | 0.37 | 43 | 1766 | 0 | 0.0 |


## Writes (isolated, one worker)

`ops/s` here is the cost of one update. Concurrency is a different question and is measured separately.

| Topology | Design | Write | ops/s | p50 ms | p99 ms | ops | errors | spread % |
|---|---|---|---|---|---|---|---|---|
| pg-single | `n1` | w01_modify_key | 1.3k | 0.74 | 1.41 | 506 | 0 | 0.0 |
| pg-single | `n1` | w02_add_key | 1.2k | 0.77 | 1.41 | 484 | 0 | 0.0 |
| pg-single | `n1` | w03_delete_key | 631 | 1.50 | 2.51 | 253 | 0 | 0.0 |
| pg-single | `n1` | w04_publish_batch | 1.2k | 0.81 | 1.52 | 465 | 0 | 0.0 |
| pg-single | `n1` | w05_replace_all | 1.0k | 0.86 | 1.63 | 412 | 0 | 0.0 |
| pg-single | `n1` | w06_update_metadata | 1.6k | 0.58 | 1.02 | 646 | 0 | 0.0 |


## Contention (hot key, one installed product)

| Topology | Design | Strategy | Writers | Acked | Retries | Conflicts | ops/s | p99 ms | counter | expected | lost updates | control |
|---|---|---|---|---|---|---|---|---|---|---|---|---|


## Audits

| Topology | Design | Phase | Passed | Checks | Key/value mismatches | Revision mismatches | Design checks | First failure |
|---|---|---|---|---|---|---|---|---|
| pg-single | `n1` | 1 | true | 12 | 0 | 0 |  |  |


## Storage

PostgreSQL reports `pg_total_relation_size` per relation. The route to a YugabyteDB
size figure is a different query, so none is reported here rather than a number that
would be read as bytes on disk.

| Topology | Design | Relation | Bytes |
|---|---|---|---|
| pg-single | `n1` | business_unit | 32.0 KiB |
| pg-single | `n1` | config_entry | 48.0 KiB |
| pg-single | `n1` | environment | 32.0 KiB |
| pg-single | `n1` | installed_product | 32.0 KiB |
| pg-single | `n1` | load_entry | 16.0 KiB |
| pg-single | `n1` | product_definition | 32.0 KiB |


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

