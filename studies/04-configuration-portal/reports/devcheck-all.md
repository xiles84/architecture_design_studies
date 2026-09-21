# Study 04 — configuration portal — run `devcheck-all`

Generated 2026-09-21T11:48:29Z from 10 result file(s). This report contains measurements only.

| Field | Value |
|---|---|
| Study | 04-configuration-portal |
| Run id | `devcheck-all` |
| Environment | `host-zenbook-ux5406sa` |
| Repository commit | `2e959908963509d26a8f6a9cabdded595b8fafbe` |
| `git describe` | `repo/open-loop-arrivals-dirty` |
| Working tree dirty at start | true |
| **Inputs digest** | `1ec2c960be96097e` |
| Cells | 10 |

## TL;DR — measured facts, selected by fixed rules

- Cells reported: 10; cells that failed: 10.
  - failed: `pg-single/c1_optimistic_version` (phase contention: audit after contention failed: INV-3/INV-10: installation 1 is missing key "limits.pool_size_004"; INV-5: installation 1 is at revision 480, t…)
  - failed: `pg-single/c2_pessimistic_lock` (phase contention: audit after contention failed: INV-3/INV-10: installation 1 is missing key "db.ttl_006"; INV-5: installation 1 is at revision 510, the ledger …)
  - failed: `pg-single/d1_doc_row` (phase write: audit after writes failed: INV-8/INV-9: a_content_hash_mismatches reports 6 mismatches)
  - failed: `pg-single/n0_rows_unindexed` (phase contention: audit after contention failed: INV-3/INV-10: installation 1 is missing key "cache.pool_size_001"; INV-5: installation 1 is at revision 533, th…)
  - failed: `pg-single/n1_rows_indexed` (phase contention: audit after contention failed: INV-3/INV-10: installation 1 is missing key "db.pool_size_000"; INV-5: installation 1 is at revision 496, the l…)
  - failed: `pg-single/n2_rows_rolldown` (phase contention: audit after contention failed: INV-3/INV-10: installation 1 is missing key "limits.pool_size_004"; INV-5: installation 1 is at revision 483, t…)
  - failed: `pg-single/n3_rollup_trigger` (phase contention: audit after contention failed: INV-3/INV-10: installation 1 is missing key "cache.pool_size_001"; INV-5: installation 1 is at revision 356, th…)
  - failed: `pg-single/n4_rollup_app` (phase contention: audit after contention failed: INV-3/INV-10: installation 1 is missing key "db.pool_size_000"; INV-5: installation 1 is at revision 419, the l…)
  - failed: `pg-single/x1_lost_update_control` (design x1_lost_update_control statement a_version_sum uses $n but declares no -- params)
  - failed: `pg-single/x2_rollup_drift_control` (phase write: audit after writes failed: INV-8/INV-9: a_rollup_mismatches reports 6 mismatches)
- Negative controls that fired: 1. Controls that did NOT fire: 1.
  - did not fire: `pg-single/x1_lost_update_control` — the associated correctness claim is **not demonstrated**
- Invariant violations recorded: 9.
  - `pg-single/c1_optimistic_version`: INV-3/INV-10: installation 1 is missing key "limits.pool_size_004"
  - `pg-single/c2_pessimistic_lock`: INV-3/INV-10: installation 1 is missing key "db.ttl_006"
  - `pg-single/d1_doc_row`: INV-8/INV-9: a_content_hash_mismatches reports 6 mismatches
  - `pg-single/n0_rows_unindexed`: INV-3/INV-10: installation 1 is missing key "cache.pool_size_001"
  - `pg-single/n1_rows_indexed`: INV-3/INV-10: installation 1 is missing key "db.pool_size_000"
  - `pg-single/n2_rows_rolldown`: INV-3/INV-10: installation 1 is missing key "limits.pool_size_004"
  - `pg-single/n3_rollup_trigger`: INV-3/INV-10: installation 1 is missing key "cache.pool_size_001"
  - `pg-single/n4_rollup_app`: INV-3/INV-10: installation 1 is missing key "db.pool_size_000"
  - `pg-single/x2_rollup_drift_control`: INV-8/INV-9: a_rollup_mismatches reports 6 mismatches
- Highest and lowest valid read result, across all cells and reads: 52854.4 ops/s (`pg-single/n3/r03_revision_check`) and 15016.5 ops/s (`pg-single/c2/r04_list_installations`).

This list ranks by number only. Which design is *better* is not a measurement and is not stated here.

## Analyses of this data

No signed analysis exists yet for inputs digest `1ec2c960be96097e`.

## Cells

| Topology | Design | Family | Kind | Gate | Load s | Entries | Serialized B |
|---|---|---|---|---|---|---|---|
| pg-single | `c1_optimistic_version` | concurrency | rows | aborted | 0.00 | 48 | 5.3 KiB |
| pg-single | `c2_pessimistic_lock` | concurrency | rows | aborted | 0.00 | 48 | 5.3 KiB |
| pg-single | `d1_doc_row` | document | doc | aborted | 0.00 | 48 | 5.3 KiB |
| pg-single | `n0_rows_unindexed` | normalized | rows | aborted | 0.00 | 48 | 5.3 KiB |
| pg-single | `n1_rows_indexed` | normalized | rows | aborted | 0.00 | 48 | 5.3 KiB |
| pg-single | `n2_rows_rolldown` | normalized | rows | aborted | 0.00 | 48 | 5.3 KiB |
| pg-single | `n3_rollup_trigger` | rollup | rows | aborted | 0.00 | 48 | 5.3 KiB |
| pg-single | `n4_rollup_app` | rollup | rows | aborted | 0.00 | 48 | 5.3 KiB |
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
| pg-single | `c1` | r01_effective_config | 33k | 0.13 | 0.39 | 41 | 13395 | 0 | 0.0 |
| pg-single | `c1` | r02_read_key | 40k | 0.11 | 0.47 | 37 | 16085 | 0 | 0.0 |
| pg-single | `c1` | r03_revision_check | 49k | 0.11 | 0.47 | 27 | 19629 | 0 | 0.0 |
| pg-single | `c1` | r04_list_installations | 18k | 0.17 | 0.79 | 59 | 7228 | 0 | 0.0 |
| pg-single | `c1` | r05_search_key | 34k | 0.13 | 0.40 | 43 | 13695 | 0 | 0.0 |
| pg-single | `c2` | r01_effective_config | 33k | 0.13 | 0.45 | 41 | 13189 | 0 | 0.0 |
| pg-single | `c2` | r02_read_key | 39k | 0.11 | 0.38 | 41 | 15777 | 0 | 0.0 |
| pg-single | `c2` | r03_revision_check | 50k | 0.10 | 0.36 | 29 | 20120 | 0 | 0.0 |
| pg-single | `c2` | r04_list_installations | 15k | 0.21 | 0.70 | 59 | 6012 | 0 | 0.0 |
| pg-single | `c2` | r05_search_key | 33k | 0.13 | 0.42 | 41 | 13273 | 0 | 0.0 |
| pg-single | `d1` | r01_effective_config | 34k | 0.13 | 0.45 | 42 | 13444 | 0 | 0.0 |
| pg-single | `d1` | r02_read_key | 44k | 0.10 | 0.36 | 39 | 17752 | 0 | 0.0 |
| pg-single | `d1` | r03_revision_check | 47k | 0.11 | 0.41 | 30 | 19003 | 0 | 0.0 |
| pg-single | `d1` | r04_list_installations | 18k | 0.18 | 0.83 | 54 | 7360 | 0 | 0.0 |
| pg-single | `d1` | r05_search_key | 31k | 0.14 | 0.50 | 39 | 12455 | 0 | 0.0 |
| pg-single | `n0` | r01_effective_config | 37k | 0.12 | 0.40 | 39 | 14629 | 0 | 0.0 |
| pg-single | `n0` | r02_read_key | 42k | 0.12 | 0.44 | 32 | 16667 | 0 | 0.0 |
| pg-single | `n0` | r03_revision_check | 50k | 0.10 | 0.36 | 25 | 20161 | 0 | 0.0 |
| pg-single | `n0` | r04_list_installations | 17k | 0.19 | 0.74 | 55 | 6800 | 0 | 0.0 |
| pg-single | `n0` | r05_search_key | 31k | 0.13 | 0.47 | 45 | 12547 | 0 | 0.0 |
| pg-single | `n1` | r01_effective_config | 33k | 0.13 | 0.43 | 44 | 13183 | 0 | 0.0 |
| pg-single | `n1` | r02_read_key | 41k | 0.11 | 0.40 | 38 | 16384 | 0 | 0.0 |
| pg-single | `n1` | r03_revision_check | 48k | 0.11 | 0.43 | 29 | 19084 | 0 | 0.0 |
| pg-single | `n1` | r04_list_installations | 17k | 0.19 | 0.65 | 57 | 6994 | 0 | 0.0 |
| pg-single | `n1` | r05_search_key | 34k | 0.13 | 0.39 | 42 | 13577 | 0 | 0.0 |
| pg-single | `n2` | r01_effective_config | 38k | 0.12 | 0.41 | 37 | 15390 | 0 | 0.0 |
| pg-single | `n2` | r02_read_key | 44k | 0.11 | 0.35 | 36 | 17512 | 0 | 0.0 |
| pg-single | `n2` | r03_revision_check | 47k | 0.11 | 0.40 | 24 | 19323 | 0 | 0.0 |
| pg-single | `n2` | r04_list_installations | 15k | 0.20 | 0.72 | 60 | 6032 | 0 | 0.0 |
| pg-single | `n2` | r05_search_key | 41k | 0.12 | 0.37 | 35 | 16428 | 0 | 0.0 |
| pg-single | `n3` | r01_effective_config | 35k | 0.13 | 0.56 | 41 | 13843 | 0 | 0.0 |
| pg-single | `n3` | r02_read_key | 36k | 0.12 | 0.57 | 37 | 14554 | 0 | 0.0 |
| pg-single | `n3` | r03_revision_check | 53k | 0.10 | 0.33 | 33 | 21155 | 0 | 0.0 |
| pg-single | `n3` | r04_list_installations | 28k | 0.14 | 0.44 | 47 | 11065 | 0 | 0.0 |
| pg-single | `n3` | r05_search_key | 35k | 0.12 | 0.41 | 43 | 13853 | 0 | 0.0 |
| pg-single | `n4` | r01_effective_config | 31k | 0.13 | 0.55 | 43 | 12344 | 0 | 0.0 |
| pg-single | `n4` | r02_read_key | 42k | 0.11 | 0.38 | 37 | 16615 | 0 | 0.0 |
| pg-single | `n4` | r03_revision_check | 50k | 0.10 | 0.39 | 24 | 19840 | 0 | 0.0 |
| pg-single | `n4` | r04_list_installations | 28k | 0.14 | 0.41 | 47 | 11186 | 0 | 0.0 |
| pg-single | `n4` | r05_search_key | 33k | 0.13 | 0.45 | 42 | 13401 | 0 | 0.0 |
| pg-single | `x2` | r01_effective_config | 33k | 0.13 | 0.61 | 43 | 13332 | 0 | 0.0 |
| pg-single | `x2` | r02_read_key | 40k | 0.11 | 0.42 | 39 | 16126 | 0 | 0.0 |
| pg-single | `x2` | r03_revision_check | 43k | 0.11 | 0.51 | 32 | 17565 | 0 | 0.0 |
| pg-single | `x2` | r04_list_installations | 27k | 0.14 | 0.51 | 49 | 10854 | 0 | 0.0 |
| pg-single | `x2` | r05_search_key | 34k | 0.13 | 0.46 | 41 | 13529 | 0 | 0.0 |


## Writes (isolated, one worker)

`ops/s` here is the cost of one update. Concurrency is a different question and is measured separately.

| Topology | Design | Write | ops/s | p50 ms | p99 ms | ops | errors | spread % |
|---|---|---|---|---|---|---|---|---|
| pg-single | `c1` | w01_modify_key | 1.4k | 0.68 | 1.26 | 546 | 0 | 0.0 |
| pg-single | `c1` | w02_add_key | 1.3k | 0.73 | 1.39 | 517 | 0 | 0.0 |
| pg-single | `c1` | w03_delete_key | 675 | 1.39 | 2.44 | 270 | 0 | 0.0 |
| pg-single | `c1` | w04_publish_batch | 1.2k | 0.82 | 1.33 | 463 | 0 | 0.0 |
| pg-single | `c1` | w05_replace_all | 935 | 0.96 | 1.79 | 375 | 0 | 0.0 |
| pg-single | `c1` | w06_update_metadata | 1.5k | 0.62 | 1.21 | 619 | 0 | 0.0 |
| pg-single | `c2` | w01_modify_key | 1.2k | 0.81 | 1.59 | 470 | 0 | 0.0 |
| pg-single | `c2` | w02_add_key | 1.3k | 0.72 | 1.21 | 527 | 0 | 0.0 |
| pg-single | `c2` | w03_delete_key | 672 | 1.41 | 2.41 | 269 | 0 | 0.0 |
| pg-single | `c2` | w04_publish_batch | 1.2k | 0.79 | 1.49 | 482 | 0 | 0.0 |
| pg-single | `c2` | w05_replace_all | 977 | 0.97 | 1.68 | 392 | 0 | 0.0 |
| pg-single | `c2` | w06_update_metadata | 1.6k | 0.60 | 1.07 | 640 | 0 | 0.0 |
| pg-single | `d1` | wd1_write_doc | 912 | 1.01 | 1.87 | 366 | 0 | 0.0 |
| pg-single | `d1` | w06_update_metadata | 1.6k | 0.60 | 1.11 | 631 | 0 | 0.0 |
| pg-single | `n0` | w01_modify_key | 1.5k | 0.64 | 1.18 | 582 | 0 | 0.0 |
| pg-single | `n0` | w02_add_key | 1.4k | 0.70 | 1.30 | 543 | 0 | 0.0 |
| pg-single | `n0` | w03_delete_key | 709 | 1.32 | 2.23 | 284 | 0 | 0.0 |
| pg-single | `n0` | w04_publish_batch | 1.2k | 0.81 | 1.51 | 466 | 0 | 0.0 |
| pg-single | `n0` | w05_replace_all | 1.1k | 0.89 | 1.60 | 426 | 0 | 0.0 |
| pg-single | `n0` | w06_update_metadata | 1.6k | 0.59 | 1.18 | 631 | 0 | 0.0 |
| pg-single | `n1` | w01_modify_key | 1.3k | 0.76 | 1.42 | 505 | 0 | 0.0 |
| pg-single | `n1` | w02_add_key | 1.2k | 0.76 | 1.63 | 497 | 0 | 0.0 |
| pg-single | `n1` | w03_delete_key | 650 | 1.46 | 2.42 | 261 | 0 | 0.0 |
| pg-single | `n1` | w04_publish_batch | 1.2k | 0.79 | 1.39 | 482 | 0 | 0.0 |
| pg-single | `n1` | w05_replace_all | 1.0k | 0.93 | 1.66 | 405 | 0 | 0.0 |
| pg-single | `n1` | w06_update_metadata | 1.6k | 0.60 | 1.05 | 642 | 0 | 0.0 |
| pg-single | `n2` | w01_modify_key | 1.2k | 0.81 | 1.40 | 470 | 0 | 0.0 |
| pg-single | `n2` | w02_add_key | 1.2k | 0.80 | 1.46 | 475 | 0 | 0.0 |
| pg-single | `n2` | w03_delete_key | 637 | 1.51 | 2.47 | 256 | 0 | 0.0 |
| pg-single | `n2` | w04_publish_batch | 1.1k | 0.91 | 1.52 | 423 | 0 | 0.0 |
| pg-single | `n2` | w05_replace_all | 866 | 1.04 | 1.71 | 347 | 0 | 0.0 |
| pg-single | `n2` | w06_update_metadata | 1.6k | 0.60 | 1.13 | 631 | 0 | 0.0 |
| pg-single | `n3` | w01_modify_key | 1.2k | 0.78 | 1.52 | 477 | 0 | 0.0 |
| pg-single | `n3` | w02_add_key | 1.2k | 0.78 | 1.44 | 481 | 0 | 0.0 |
| pg-single | `n3` | w03_delete_key | 593 | 1.59 | 2.54 | 238 | 0 | 0.0 |
| pg-single | `n3` | w04_publish_batch | 587 | 1.68 | 2.32 | 235 | 0 | 0.0 |
| pg-single | `n3` | w05_replace_all | 527 | 1.36 | 2.36 | 211 | 0 | 0.0 |
| pg-single | `n3` | w06_update_metadata | 1.5k | 0.61 | 1.10 | 615 | 0 | 0.0 |
| pg-single | `n4` | w01_modify_key | 1.1k | 0.89 | 1.56 | 425 | 0 | 0.0 |
| pg-single | `n4` | w02_add_key | 1.1k | 0.87 | 1.67 | 434 | 0 | 0.0 |
| pg-single | `n4` | w03_delete_key | 548 | 1.73 | 2.58 | 220 | 0 | 0.0 |
| pg-single | `n4` | w04_publish_batch | 866 | 1.10 | 1.84 | 347 | 0 | 0.0 |
| pg-single | `n4` | w05_replace_all | 945 | 1.00 | 1.79 | 378 | 0 | 0.0 |
| pg-single | `n4` | w06_update_metadata | 1.6k | 0.60 | 1.20 | 625 | 0 | 0.0 |
| pg-single | `x2` | w01_modify_key | 1.3k | 0.70 | 1.34 | 539 | 0 | 0.0 |
| pg-single | `x2` | w02_add_key | 1.3k | 0.70 | 1.41 | 537 | 0 | 0.0 |
| pg-single | `x2` | w03_delete_key | 665 | 1.42 | 2.41 | 267 | 0 | 0.0 |
| pg-single | `x2` | w04_publish_batch | 1.2k | 0.81 | 1.44 | 467 | 0 | 0.0 |
| pg-single | `x2` | w05_replace_all | 951 | 0.94 | 1.67 | 381 | 0 | 0.0 |
| pg-single | `x2` | w06_update_metadata | 1.6k | 0.57 | 1.25 | 657 | 0 | 0.0 |


## Contention (hot key, one installed product)

| Topology | Design | Strategy | Writers | Acked | Retries | Conflicts | ops/s | p99 ms | counter | expected | lost updates | control |
|---|---|---|---|---|---|---|---|---|---|---|---|---|
| pg-single | `c1` | optimistic | 16 | 0 | 0 | 0 | — | — | 0 | 0 | 0 |  |
| pg-single | `c2` | pessimistic | 16 | 0 | 0 | 0 | — | — | 0 | 0 | 0 |  |
| pg-single | `n0` |  | 16 | 0 | 0 | 0 | — | — | 0 | 0 | 0 |  |
| pg-single | `n1` |  | 16 | 0 | 0 | 0 | — | — | 0 | 0 | 0 |  |
| pg-single | `n2` |  | 16 | 0 | 0 | 0 | — | — | 0 | 0 | 0 |  |
| pg-single | `n3` |  | 16 | 0 | 0 | 0 | — | — | 0 | 0 | 0 |  |
| pg-single | `n4` |  | 16 | 0 | 0 | 0 | — | — | 0 | 0 | 0 |  |


## Audits

| Topology | Design | Phase | Passed | Checks | Key/value mismatches | Revision mismatches | Design checks | First failure |
|---|---|---|---|---|---|---|---|---|
| pg-single | `c1` | 1 | true | 12 | 0 | 0 |  |  |
| pg-single | `c1` | 2 | false | 12 | 6 | 6 |  | INV-3/INV-10: installation 1 is missing key "limits.pool_size_004" |
| pg-single | `c2` | 1 | true | 12 | 0 | 0 |  |  |
| pg-single | `c2` | 2 | false | 12 | 6 | 6 |  | INV-3/INV-10: installation 1 is missing key "db.ttl_006" |
| pg-single | `d1` | 1 | false | 14 | 0 | 0 | a_content_hash_mismatches=6 a_document_revision_mismatches=0 | INV-8/INV-9: a_content_hash_mismatches reports 6 mismatches |
| pg-single | `n0` | 1 | true | 12 | 0 | 0 |  |  |
| pg-single | `n0` | 2 | false | 12 | 6 | 6 |  | INV-3/INV-10: installation 1 is missing key "cache.pool_size_001" |
| pg-single | `n1` | 1 | true | 12 | 0 | 0 |  |  |
| pg-single | `n1` | 2 | false | 12 | 6 | 6 |  | INV-3/INV-10: installation 1 is missing key "db.pool_size_000" |
| pg-single | `n2` | 1 | true | 13 | 0 | 0 | a_rolldown_mismatches=0 |  |
| pg-single | `n2` | 2 | false | 13 | 6 | 6 | a_rolldown_mismatches=0 | INV-3/INV-10: installation 1 is missing key "limits.pool_size_004" |
| pg-single | `n3` | 1 | true | 13 | 0 | 0 | a_rollup_mismatches=0 |  |
| pg-single | `n3` | 2 | false | 13 | 6 | 6 | a_rollup_mismatches=0 | INV-3/INV-10: installation 1 is missing key "cache.pool_size_001" |
| pg-single | `n4` | 1 | true | 13 | 0 | 0 | a_rollup_mismatches=0 |  |
| pg-single | `n4` | 2 | false | 13 | 6 | 6 | a_rollup_mismatches=0 | INV-3/INV-10: installation 1 is missing key "db.pool_size_000" |
| pg-single | `x2` | 1 | false | 13 | 0 | 0 | a_rollup_mismatches=6 | INV-8/INV-9: a_rollup_mismatches reports 6 mismatches |


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
| pg-single | `d1` | business_unit | 32.0 KiB |
| pg-single | `d1` | config_document | 32.0 KiB |
| pg-single | `d1` | environment | 32.0 KiB |
| pg-single | `d1` | installed_product | 32.0 KiB |
| pg-single | `d1` | load_entry | 16.0 KiB |
| pg-single | `d1` | product_definition | 32.0 KiB |
| pg-single | `n0` | business_unit | 32.0 KiB |
| pg-single | `n0` | config_entry | 32.0 KiB |
| pg-single | `n0` | environment | 32.0 KiB |
| pg-single | `n0` | installed_product | 32.0 KiB |
| pg-single | `n0` | load_entry | 16.0 KiB |
| pg-single | `n0` | product_definition | 32.0 KiB |
| pg-single | `n1` | business_unit | 32.0 KiB |
| pg-single | `n1` | config_entry | 48.0 KiB |
| pg-single | `n1` | environment | 32.0 KiB |
| pg-single | `n1` | installed_product | 32.0 KiB |
| pg-single | `n1` | load_entry | 16.0 KiB |
| pg-single | `n1` | product_definition | 32.0 KiB |
| pg-single | `n2` | business_unit | 32.0 KiB |
| pg-single | `n2` | config_entry | 112.0 KiB |
| pg-single | `n2` | environment | 32.0 KiB |
| pg-single | `n2` | installed_product | 32.0 KiB |
| pg-single | `n2` | load_entry | 16.0 KiB |
| pg-single | `n2` | product_definition | 32.0 KiB |
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

