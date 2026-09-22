# Study 04 — configuration portal — run `devcheck-ctrl`

Generated 2026-09-21T11:53:21Z from 4 result file(s). This report contains measurements only.

| Field | Value |
|---|---|
| Study | 04-configuration-portal |
| Run id | `devcheck-ctrl` |
| Environment | `host-zenbook-ux5406sa` |
| Repository commit | `2e959908963509d26a8f6a9cabdded595b8fafbe` |
| `git describe` | `repo/open-loop-arrivals-dirty` |
| Working tree dirty at start | true |
| **Inputs digest** | `525b32f62df9fb44` |
| Cells | 4 |

## TL;DR — measured facts, selected by fixed rules

- Cells reported: 4; cells that failed: 0.
- Negative controls that fired: 2. Controls that did NOT fire: 0.
- Invariant violations recorded: 2.
  - `pg-single/x2_rollup_drift_control`: INV-8/INV-9: a_rollup_mismatches reports 6 mismatches
  - `pg-single/x2_rollup_drift_control`: INV-8/INV-9: a_rollup_mismatches reports 6 mismatches

This list ranks by number only. Which design is *better* is not a measurement and is not stated here.

## Analyses of this data

No signed analysis exists yet for inputs digest `525b32f62df9fb44`.

## Cells

| Topology | Design | Family | Kind | Gate | Load s | Entries | Serialized B |
|---|---|---|---|---|---|---|---|
| pg-single | `c1_optimistic_version` | concurrency | rows | pass | 0.00 | 48 | 5.3 KiB |
| pg-single | `c2_pessimistic_lock` | concurrency | rows | pass | 0.00 | 48 | 5.3 KiB |
| pg-single | `x1_lost_update_control` | control | rows | pass | 0.00 | 48 | 5.3 KiB |
| pg-single | `x2_rollup_drift_control` | control | rows | pass | 0.00 | 48 | 5.3 KiB |


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


## Writes (isolated, one worker)

`ops/s` here is the cost of one update. Concurrency is a different question and is measured separately.

| Topology | Design | Write | ops/s | p50 ms | p99 ms | ops | errors | spread % |
|---|---|---|---|---|---|---|---|---|
| pg-single | `c1` | w01_modify_key | 1.4k | 0.68 | 1.25 | 554 | 0 | 0.0 |
| pg-single | `c1` | w02_add_key | 1.3k | 0.71 | 1.32 | 531 | 0 | 0.0 |
| pg-single | `c1` | w03_delete_key | 629 | 1.49 | 2.42 | 252 | 0 | 0.0 |
| pg-single | `c1` | w04_publish_batch | 1.1k | 0.80 | 1.72 | 433 | 0 | 0.0 |
| pg-single | `c1` | w05_replace_all | 988 | 0.97 | 1.68 | 396 | 0 | 0.0 |
| pg-single | `c1` | w06_update_metadata | 1.5k | 0.62 | 1.10 | 612 | 0 | 0.0 |
| pg-single | `c2` | w01_modify_key | 1.2k | 0.77 | 1.59 | 495 | 0 | 0.0 |
| pg-single | `c2` | w02_add_key | 1.3k | 0.71 | 1.34 | 527 | 0 | 0.0 |
| pg-single | `c2` | w03_delete_key | 672 | 1.40 | 2.33 | 269 | 0 | 0.0 |
| pg-single | `c2` | w04_publish_batch | 1.1k | 0.84 | 1.44 | 453 | 0 | 0.0 |
| pg-single | `c2` | w05_replace_all | 1.0k | 0.89 | 1.81 | 414 | 0 | 0.0 |
| pg-single | `c2` | w06_update_metadata | 1.6k | 0.59 | 1.31 | 628 | 0 | 0.0 |
| pg-single | `x1` | w01_modify_key | 1.2k | 0.79 | 1.49 | 478 | 0 | 0.0 |
| pg-single | `x1` | w02_add_key | 1.3k | 0.74 | 1.41 | 513 | 0 | 0.0 |
| pg-single | `x1` | w03_delete_key | 692 | 1.37 | 2.27 | 277 | 0 | 0.0 |
| pg-single | `x1` | w04_publish_batch | 1.2k | 0.81 | 1.32 | 470 | 0 | 0.0 |
| pg-single | `x1` | w05_replace_all | 986 | 0.93 | 1.80 | 395 | 0 | 0.0 |
| pg-single | `x1` | w06_update_metadata | 1.6k | 0.60 | 1.16 | 624 | 0 | 0.0 |
| pg-single | `x2` | w01_modify_key | 1.4k | 0.68 | 1.24 | 559 | 0 | 0.0 |
| pg-single | `x2` | w02_add_key | 1.3k | 0.72 | 1.41 | 523 | 0 | 0.0 |
| pg-single | `x2` | w03_delete_key | 677 | 1.40 | 2.26 | 272 | 0 | 0.0 |
| pg-single | `x2` | w04_publish_batch | 1.1k | 0.89 | 1.52 | 430 | 0 | 0.0 |
| pg-single | `x2` | w05_replace_all | 959 | 0.98 | 1.68 | 384 | 0 | 0.0 |
| pg-single | `x2` | w06_update_metadata | 1.6k | 0.60 | 1.17 | 629 | 0 | 0.0 |


## Contention (hot key, one installed product)

| Topology | Design | Strategy | Writers | Acked | Retries | Conflicts | ops/s | p99 ms | counter | expected | lost updates | control |
|---|---|---|---|---|---|---|---|---|---|---|---|---|
| pg-single | `c1` | optimistic | 16 | 321 | 88 | 1663 | 634 | 45 | 321 | 321 | 0 |  |
| pg-single | `c2` | pessimistic | 16 | 473 | 0 | 0 | 1.1k | 59 | 473 | 473 | 0 |  |
| pg-single | `x1` | unsafe | 16 | 552 | 0 | 0 | 1.0k | 72 | 36 | 552 | 516 | negative control |
| pg-single | `x2` | none | 16 | 371 | 0 | 0 | 650 | 81 | 0 | 0 | 0 | negative control |


## Audits

| Topology | Design | Phase | Passed | Checks | Key/value mismatches | Revision mismatches | Design checks | First failure |
|---|---|---|---|---|---|---|---|---|
| pg-single | `c1` | 1 | true | 12 | 0 | 0 |  |  |
| pg-single | `c1` | 2 | true | 12 | 0 | 0 |  |  |
| pg-single | `c2` | 1 | true | 12 | 0 | 0 |  |  |
| pg-single | `c2` | 2 | true | 12 | 0 | 0 |  |  |
| pg-single | `x1` | 1 | true | 12 | 0 | 0 |  |  |
| pg-single | `x1` | 2 | true | 12 | 0 | 0 |  |  |
| pg-single | `x2` | 1 | false | 13 | 0 | 0 | a_rollup_mismatches=6 | INV-8/INV-9: a_rollup_mismatches reports 6 mismatches |
| pg-single | `x2` | 2 | false | 13 | 0 | 0 | a_rollup_mismatches=6 | INV-8/INV-9: a_rollup_mismatches reports 6 mismatches |


## Storage

PostgreSQL reports `pg_total_relation_size` per relation. The route to a YugabyteDB
size figure is a different query, so none is reported here rather than a number that
would be read as bytes on disk.

| Topology | Design | Relation | Bytes |
|---|---|---|---|
| pg-single | `c1` | business_unit | 32.0 KiB |
| pg-single | `c1` | config_entry | 48.0 KiB |
| pg-single | `c1` | environment | 32.0 KiB |
| pg-single | `c1` | installed_product | 32.0 KiB |
| pg-single | `c1` | load_entry | 16.0 KiB |
| pg-single | `c1` | product_definition | 32.0 KiB |
| pg-single | `c2` | business_unit | 32.0 KiB |
| pg-single | `c2` | config_entry | 48.0 KiB |
| pg-single | `c2` | environment | 32.0 KiB |
| pg-single | `c2` | installed_product | 32.0 KiB |
| pg-single | `c2` | load_entry | 16.0 KiB |
| pg-single | `c2` | product_definition | 32.0 KiB |
| pg-single | `x1` | business_unit | 32.0 KiB |
| pg-single | `x1` | config_entry | 48.0 KiB |
| pg-single | `x1` | environment | 32.0 KiB |
| pg-single | `x1` | installed_product | 32.0 KiB |
| pg-single | `x1` | load_entry | 16.0 KiB |
| pg-single | `x1` | product_definition | 32.0 KiB |
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

