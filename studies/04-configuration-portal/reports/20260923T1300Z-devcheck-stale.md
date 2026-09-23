# Study 04 — configuration portal — run `20260923T1300Z-devcheck-stale`

Generated 2026-09-23T12:47:46Z from 2 result file(s). This report contains measurements only.

| Field | Value |
|---|---|
| Study | 04-configuration-portal |
| Run id | `20260923T1300Z-devcheck-stale` |
| Environment | `host-zenbook-ux5406sa` |
| Repository commit | `5e0906d09e5592c0e9068d55addff2577ab43c75` |
| `git describe` | `study-05/v2-resources-engines-4-g5e0906d-dirty` |
| Working tree dirty at start | true |
| **Inputs digest** | `bb727608b5369f7d` |
| Cells | 2 |

## TL;DR — measured facts, selected by fixed rules

- Cells reported: 2; cells that failed: 0.
- Negative controls that fired: 0. Controls that did NOT fire: 0.
- No invariant violation was recorded by a design that is not a negative control.

This list ranks by number only. Which design is *better* is not a measurement and is not stated here.

## Analyses of this data

No signed analysis exists yet for inputs digest `bb727608b5369f7d`.

## Cells

| Topology | Design | Family | Kind | Gate | Load s | Entries | Serialized B |
|---|---|---|---|---|---|---|---|
| pg-single | `n3_rollup_trigger` | rollup | rows | pass | 0.02 | 48 | 5.3 KiB |
| pg-single | `n4_rollup_app` | rollup | rows | pass | 0.01 | 48 | 5.3 KiB |


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
| pg-single | `n3` | w01_modify_key | 319 | 3.04 | 4.16 | 128 | 0 | 0.0 |
| pg-single | `n3` | w02_add_key | 325 | 3.07 | 3.66 | 131 | 0 | 0.0 |
| pg-single | `n3` | w03_delete_key | 167 | 5.87 | 7.84 | 67 | 0 | 0.0 |
| pg-single | `n3` | w04_publish_batch | 197 | 4.01 | 17 | 79 | 0 | 0.0 |
| pg-single | `n3` | w05_replace_all | 223 | 4.38 | 6.42 | 90 | 0 | 0.0 |
| pg-single | `n3` | w06_update_metadata | 341 | 2.59 | 6.32 | 137 | 0 | 0.0 |
| pg-single | `n4` | w01_modify_key | 282 | 3.49 | 4.53 | 113 | 0 | 0.0 |
| pg-single | `n4` | w02_add_key | 299 | 3.20 | 4.70 | 120 | 0 | 0.0 |
| pg-single | `n4` | w03_delete_key | 150 | 6.59 | 8.76 | 61 | 0 | 0.0 |
| pg-single | `n4` | w04_publish_batch | 265 | 3.77 | 4.83 | 106 | 0 | 0.0 |
| pg-single | `n4` | w05_replace_all | 255 | 3.76 | 5.20 | 103 | 0 | 0.0 |
| pg-single | `n4` | w06_update_metadata | 342 | 2.76 | 4.90 | 138 | 0 | 0.0 |


## Contention (hot key, one installed product)

| Topology | Design | Strategy | Writers | Acked | Retries | Conflicts | ops/s | p99 ms | counter | expected | lost updates | control |
|---|---|---|---|---|---|---|---|---|---|---|---|---|


## Cadence

The offered rate is *calculated*: installed products / period. A compressed rate is a
compressed-time validation, never the temporal result it stands in for.

| Topology | Design | Regime | Fleet | Calculated rate /s | Offered | Started | Completed | Dropped | Queue max | Client saturated | lag p99 ms | Compressed |
|---|---|---|---|---|---|---|---|---|---|---|---|---|
| pg-single | `n3` | secondly | 6 | 6.0000 | 3 | 3 | 3 | 0 | 0 | false | 0.75 | false |
| pg-single | `n4` | secondly | 6 | 6.0000 | 3 | 3 | 3 | 0 | 0 | false | 0.89 | false |


## Audits

| Topology | Design | Phase | Passed | Checks | Key/value mismatches | Revision mismatches | Design checks | First failure |
|---|---|---|---|---|---|---|---|---|
| pg-single | `n3` | 1 | true | 13 | 0 | 0 | a_rollup_mismatches=0 |  |
| pg-single | `n4` | 1 | true | 13 | 0 | 0 | a_rollup_mismatches=0 |  |


## Storage

PostgreSQL reports `pg_total_relation_size` per relation. The route to a YugabyteDB
size figure is a different query, so none is reported here rather than a number that
would be read as bytes on disk.

| Topology | Design | Relation | Bytes |
|---|---|---|---|
| pg-single | `n3` | business_unit | 32.0 KiB |
| pg-single | `n3` | config_entry | 48.0 KiB |
| pg-single | `n3` | environment | 32.0 KiB |
| pg-single | `n3` | installed_product | 32.0 KiB |
| pg-single | `n3` | load_entry | 16.0 KiB |
| pg-single | `n3` | product_definition | 32.0 KiB |
| pg-single | `n4` | business_unit | 32.0 KiB |
| pg-single | `n4` | config_entry | 48.0 KiB |
| pg-single | `n4` | environment | 32.0 KiB |
| pg-single | `n4` | installed_product | 32.0 KiB |
| pg-single | `n4` | load_entry | 16.0 KiB |
| pg-single | `n4` | product_definition | 32.0 KiB |


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

