# Study 04 — configuration portal — run `devchecks`

Generated 2026-09-21T12:00:31Z from 10 result file(s). This report contains measurements only.

| Field | Value |
|---|---|
| Study | 04-configuration-portal |
| Run id | `devchecks` |
| Environment | `host-zenbook-ux5406sa` |
| Repository commit | `2e959908963509d26a8f6a9cabdded595b8fafbe` |
| `git describe` | `repo/open-loop-arrivals-dirty` |
| Working tree dirty at start | true |
| **Inputs digest** | `364ee7176b19777c` |
| Cells | 10 |

## TL;DR — measured facts, selected by fixed rules

- Cells reported: 10; cells that failed: 0.
- Negative controls that fired: 2. Controls that did NOT fire: 0.
- No invariant violation was recorded by a design that is not a negative control.
- Highest and lowest valid read result, across all cells and reads: 51727.9 ops/s (`pg-single/x2/r03_revision_check`) and 3788.5 ops/s (`pg-single/c2/r05_search_key`).

This list ranks by number only. Which design is *better* is not a measurement and is not stated here.

## Analyses of this data

No signed analysis exists yet for inputs digest `364ee7176b19777c`.

## Cells

| Topology | Design | Family | Kind | Gate | Load s | Entries | Serialized B |
|---|---|---|---|---|---|---|---|
| pg-single | `c1_optimistic_version` | concurrency | rows | pass | 0.00 | 48 | 5.3 KiB |
| pg-single | `c2_pessimistic_lock` | concurrency | rows | pass | 0.00 | 48 | 5.3 KiB |
| pg-single | `d1_doc_row` | document | doc | pass | 0.00 | 48 | 5.3 KiB |
| pg-single | `n0_rows_unindexed` | normalized | rows | pass | 0.03 | 48 | 5.3 KiB |
| pg-single | `n1_rows_indexed` | normalized | rows | pass | 0.00 | 48 | 5.3 KiB |
| pg-single | `n2_rows_rolldown` | normalized | rows | pass | 0.00 | 48 | 5.3 KiB |
| pg-single | `n3_rollup_trigger` | rollup | rows | pass | 0.00 | 48 | 5.3 KiB |
| pg-single | `n4_rollup_app` | rollup | rows | pass | 0.00 | 48 | 5.3 KiB |
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
| pg-single | `c1` | r01_effective_config | 5.4k | 0.13 | 0.50 | 35 | 2148 | 0 | 0.0 |
| pg-single | `c1` | r02_read_key | 5.8k | 0.11 | 0.42 | 35 | 2323 | 0 | 0.0 |
| pg-single | `c1` | r03_revision_check | 49k | 0.11 | 0.35 | 30 | 19684 | 0 | 0.0 |
| pg-single | `c1` | r04_list_installations | 18k | 0.18 | 0.71 | 55 | 7269 | 0 | 0.0 |
| pg-single | `c1` | r05_search_key | 4.4k | 0.13 | 0.44 | 38 | 1771 | 0 | 0.0 |
| pg-single | `c2` | r01_effective_config | 5.0k | 0.12 | 0.38 | 36 | 2195 | 0 | 0.0 |
| pg-single | `c2` | r02_read_key | 5.6k | 0.12 | 0.44 | 30 | 2241 | 0 | 0.0 |
| pg-single | `c2` | r03_revision_check | 48k | 0.11 | 0.36 | 32 | 19233 | 0 | 0.0 |
| pg-single | `c2` | r04_list_installations | 18k | 0.19 | 0.62 | 56 | 7046 | 0 | 0.0 |
| pg-single | `c2` | r05_search_key | 3.8k | 0.16 | 0.59 | 41 | 1516 | 0 | 0.0 |
| pg-single | `d1` | r01_effective_config | 26k | 0.15 | 1.35 | 43 | 10426 | 0 | 0.0 |
| pg-single | `d1` | r02_read_key | 35k | 0.12 | 0.99 | 36 | 13973 | 0 | 0.0 |
| pg-single | `d1` | r03_revision_check | 46k | 0.10 | 0.59 | 33 | 18478 | 0 | 0.0 |
| pg-single | `d1` | r04_list_installations | 15k | 0.21 | 1.84 | 54 | 6056 | 0 | 0.0 |
| pg-single | `d1` | r05_search_key | 27k | 0.14 | 1.20 | 42 | 10680 | 0 | 0.0 |
| pg-single | `n0` | r01_effective_config | 5.3k | 0.12 | 0.35 | 36 | 2114 | 0 | 0.0 |
| pg-single | `n0` | r02_read_key | 5.6k | 0.11 | 0.46 | 33 | 2242 | 0 | 0.0 |
| pg-single | `n0` | r03_revision_check | 48k | 0.11 | 0.38 | 22 | 20204 | 0 | 0.0 |
| pg-single | `n0` | r04_list_installations | 17k | 0.19 | 0.66 | 58 | 6759 | 0 | 0.0 |
| pg-single | `n0` | r05_search_key | 3.8k | 0.14 | 0.51 | 47 | 1540 | 0 | 0.0 |
| pg-single | `n1` | r01_effective_config | 5.1k | 0.13 | 0.47 | 36 | 2034 | 0 | 0.0 |
| pg-single | `n1` | r02_read_key | 5.6k | 0.11 | 0.38 | 36 | 2240 | 0 | 0.0 |
| pg-single | `n1` | r03_revision_check | 51k | 0.10 | 0.36 | 24 | 21494 | 0 | 0.0 |
| pg-single | `n1` | r04_list_installations | 18k | 0.18 | 0.69 | 57 | 7303 | 0 | 0.0 |
| pg-single | `n1` | r05_search_key | 4.2k | 0.14 | 0.47 | 41 | 1691 | 0 | 0.0 |
| pg-single | `n2` | r01_effective_config | 6.2k | 0.12 | 0.33 | 29 | 2501 | 0 | 0.0 |
| pg-single | `n2` | r02_read_key | 6.0k | 0.11 | 0.33 | 36 | 2412 | 0 | 0.0 |
| pg-single | `n2` | r03_revision_check | 47k | 0.11 | 0.38 | 32 | 18971 | 0 | 0.0 |
| pg-single | `n2` | r04_list_installations | 16k | 0.20 | 0.82 | 57 | 6406 | 0 | 0.0 |
| pg-single | `n2` | r05_search_key | 5.4k | 0.12 | 0.35 | 35 | 2163 | 0 | 0.0 |
| pg-single | `n3` | r01_effective_config | 5.1k | 0.13 | 0.40 | 37 | 2038 | 0 | 0.0 |
| pg-single | `n3` | r02_read_key | 5.5k | 0.12 | 0.42 | 28 | 2207 | 0 | 0.0 |
| pg-single | `n3` | r03_revision_check | 47k | 0.11 | 0.40 | 30 | 18979 | 0 | 0.0 |
| pg-single | `n3` | r04_list_installations | 28k | 0.15 | 0.49 | 48 | 11151 | 0 | 0.0 |
| pg-single | `n3` | r05_search_key | 4.5k | 0.13 | 0.37 | 39 | 1796 | 0 | 0.0 |
| pg-single | `n4` | r01_effective_config | 5.5k | 0.12 | 0.41 | 38 | 2182 | 0 | 0.0 |
| pg-single | `n4` | r02_read_key | 5.7k | 0.11 | 0.43 | 34 | 2266 | 0 | 0.0 |
| pg-single | `n4` | r03_revision_check | 49k | 0.10 | 0.37 | 19 | 19848 | 0 | 0.0 |
| pg-single | `n4` | r04_list_installations | 27k | 0.14 | 0.51 | 50 | 10699 | 0 | 0.0 |
| pg-single | `n4` | r05_search_key | 4.5k | 0.12 | 0.41 | 42 | 1813 | 0 | 0.0 |
| pg-single | `x1` | r01_effective_config | 5.0k | 0.13 | 0.47 | 29 | 2154 | 0 | 0.0 |
| pg-single | `x1` | r02_read_key | 5.4k | 0.12 | 0.44 | 37 | 2175 | 0 | 0.0 |
| pg-single | `x1` | r03_revision_check | 51k | 0.10 | 0.34 | 29 | 20436 | 0 | 0.0 |
| pg-single | `x1` | r04_list_installations | 18k | 0.18 | 0.72 | 55 | 7232 | 0 | 0.0 |
| pg-single | `x1` | r05_search_key | 4.5k | 0.13 | 0.40 | 42 | 1814 | 0 | 0.0 |
| pg-single | `x2` | r01_effective_config | 5.4k | 0.12 | 0.61 | 35 | 2153 | 0 | 0.0 |
| pg-single | `x2` | r02_read_key | 5.6k | 0.11 | 0.45 | 33 | 2229 | 0 | 0.0 |
| pg-single | `x2` | r03_revision_check | 52k | 0.10 | 0.35 | 27 | 20720 | 0 | 0.0 |
| pg-single | `x2` | r04_list_installations | 26k | 0.15 | 0.49 | 47 | 10349 | 0 | 0.0 |
| pg-single | `x2` | r05_search_key | 4.1k | 0.14 | 0.49 | 41 | 1638 | 0 | 0.0 |


## Writes (isolated, one worker)

`ops/s` here is the cost of one update. Concurrency is a different question and is measured separately.

| Topology | Design | Write | ops/s | p50 ms | p99 ms | ops | errors | spread % |
|---|---|---|---|---|---|---|---|---|
| pg-single | `c1` | w01_modify_key | 1.4k | 0.69 | 1.39 | 544 | 0 | 0.0 |
| pg-single | `c1` | w02_add_key | 1.3k | 0.71 | 1.37 | 532 | 0 | 0.0 |
| pg-single | `c1` | w03_delete_key | 690 | 1.38 | 2.23 | 276 | 0 | 0.0 |
| pg-single | `c1` | w04_publish_batch | 1.2k | 0.79 | 1.48 | 481 | 0 | 0.0 |
| pg-single | `c1` | w05_replace_all | 937 | 1.00 | 1.92 | 375 | 0 | 0.0 |
| pg-single | `c1` | w06_update_metadata | 1.6k | 0.58 | 1.24 | 647 | 0 | 0.0 |
| pg-single | `c2` | w01_modify_key | 1.2k | 0.80 | 1.32 | 478 | 0 | 0.0 |
| pg-single | `c2` | w02_add_key | 1.2k | 0.81 | 1.59 | 468 | 0 | 0.0 |
| pg-single | `c2` | w03_delete_key | 698 | 1.36 | 2.12 | 280 | 0 | 0.0 |
| pg-single | `c2` | w04_publish_batch | 1.1k | 0.85 | 1.46 | 446 | 0 | 0.0 |
| pg-single | `c2` | w05_replace_all | 1.0k | 0.91 | 1.79 | 411 | 0 | 0.0 |
| pg-single | `c2` | w06_update_metadata | 1.6k | 0.57 | 1.10 | 656 | 0 | 0.0 |
| pg-single | `d1` | wd1_write_doc | 900 | 0.96 | 2.44 | 360 | 0 | 0.0 |
| pg-single | `d1` | w06_update_metadata | 1.4k | 0.61 | 1.36 | 575 | 0 | 0.0 |
| pg-single | `n0` | w01_modify_key | 1.3k | 0.72 | 1.32 | 518 | 0 | 0.0 |
| pg-single | `n0` | w02_add_key | 1.3k | 0.73 | 1.37 | 517 | 0 | 0.0 |
| pg-single | `n0` | w03_delete_key | 647 | 1.47 | 2.61 | 259 | 0 | 0.0 |
| pg-single | `n0` | w04_publish_batch | 1.2k | 0.84 | 1.54 | 461 | 0 | 0.0 |
| pg-single | `n0` | w05_replace_all | 962 | 0.93 | 1.90 | 386 | 0 | 0.0 |
| pg-single | `n0` | w06_update_metadata | 1.5k | 0.63 | 1.26 | 581 | 0 | 0.0 |
| pg-single | `n1` | w01_modify_key | 1.4k | 0.67 | 1.40 | 565 | 0 | 0.0 |
| pg-single | `n1` | w02_add_key | 1.3k | 0.72 | 1.32 | 524 | 0 | 0.0 |
| pg-single | `n1` | w03_delete_key | 666 | 1.40 | 2.33 | 267 | 0 | 0.0 |
| pg-single | `n1` | w04_publish_batch | 1.2k | 0.78 | 1.52 | 487 | 0 | 0.0 |
| pg-single | `n1` | w05_replace_all | 995 | 0.88 | 1.58 | 399 | 0 | 0.0 |
| pg-single | `n1` | w06_update_metadata | 1.3k | 0.72 | 1.21 | 535 | 0 | 0.0 |
| pg-single | `n2` | w01_modify_key | 1.3k | 0.72 | 1.37 | 524 | 0 | 0.0 |
| pg-single | `n2` | w02_add_key | 1.3k | 0.73 | 1.49 | 508 | 0 | 0.0 |
| pg-single | `n2` | w03_delete_key | 672 | 1.40 | 2.41 | 269 | 0 | 0.0 |
| pg-single | `n2` | w04_publish_batch | 1.0k | 0.90 | 1.57 | 420 | 0 | 0.0 |
| pg-single | `n2` | w05_replace_all | 897 | 1.00 | 1.88 | 359 | 0 | 0.0 |
| pg-single | `n2` | w06_update_metadata | 1.6k | 0.59 | 1.08 | 643 | 0 | 0.0 |
| pg-single | `n3` | w01_modify_key | 1.2k | 0.77 | 1.53 | 486 | 0 | 0.0 |
| pg-single | `n3` | w02_add_key | 1.1k | 0.87 | 1.46 | 438 | 0 | 0.0 |
| pg-single | `n3` | w03_delete_key | 510 | 1.89 | 2.92 | 204 | 0 | 0.0 |
| pg-single | `n3` | w04_publish_batch | 290 | 3.45 | 4.98 | 117 | 0 | 0.0 |
| pg-single | `n3` | w05_replace_all | 504 | 1.80 | 2.95 | 202 | 0 | 0.0 |
| pg-single | `n3` | w06_update_metadata | 1.5k | 0.61 | 1.12 | 620 | 0 | 0.0 |
| pg-single | `n4` | w01_modify_key | 1.2k | 0.82 | 1.52 | 465 | 0 | 0.0 |
| pg-single | `n4` | w02_add_key | 1.0k | 0.91 | 1.63 | 411 | 0 | 0.0 |
| pg-single | `n4` | w03_delete_key | 492 | 1.96 | 2.96 | 197 | 0 | 0.0 |
| pg-single | `n4` | w04_publish_batch | 654 | 1.52 | 2.26 | 262 | 0 | 0.0 |
| pg-single | `n4` | w05_replace_all | 894 | 1.07 | 1.88 | 358 | 0 | 0.0 |
| pg-single | `n4` | w06_update_metadata | 1.4k | 0.64 | 1.39 | 541 | 0 | 0.0 |
| pg-single | `x1` | w01_modify_key | 1.3k | 0.72 | 1.42 | 519 | 0 | 0.0 |
| pg-single | `x1` | w02_add_key | 1.3k | 0.73 | 1.41 | 516 | 0 | 0.0 |
| pg-single | `x1` | w03_delete_key | 706 | 1.34 | 2.34 | 283 | 0 | 0.0 |
| pg-single | `x1` | w04_publish_batch | 1.2k | 0.81 | 1.33 | 474 | 0 | 0.0 |
| pg-single | `x1` | w05_replace_all | 951 | 0.94 | 1.73 | 381 | 0 | 0.0 |
| pg-single | `x1` | w06_update_metadata | 1.5k | 0.62 | 1.15 | 608 | 0 | 0.0 |
| pg-single | `x2` | w01_modify_key | 1.2k | 0.78 | 1.34 | 489 | 0 | 0.0 |
| pg-single | `x2` | w02_add_key | 1.1k | 0.82 | 1.79 | 445 | 0 | 0.0 |
| pg-single | `x2` | w03_delete_key | 581 | 1.62 | 2.91 | 233 | 0 | 0.0 |
| pg-single | `x2` | w04_publish_batch | 981 | 0.92 | 2.58 | 397 | 0 | 0.0 |
| pg-single | `x2` | w05_replace_all | 857 | 1.00 | 2.55 | 343 | 0 | 0.0 |
| pg-single | `x2` | w06_update_metadata | 1.5k | 0.61 | 1.19 | 612 | 0 | 0.0 |


## Contention (hot key, one installed product)

| Topology | Design | Strategy | Writers | Acked | Retries | Conflicts | ops/s | p99 ms | counter | expected | lost updates | control |
|---|---|---|---|---|---|---|---|---|---|---|---|---|
| pg-single | `c1` | optimistic | 16 | 277 | 74 | 1474 | 606 | 52 | 277 | 277 | 0 |  |
| pg-single | `c2` | pessimistic | 16 | 492 | 0 | 0 | 1.0k | 56 | 492 | 492 | 0 |  |
| pg-single | `x1` | unsafe | 16 | 467 | 0 | 0 | 1.1k | 59 | 23 | 467 | 444 | negative control |
| pg-single | `x2` | none | 16 | 373 | 0 | 0 | 709 | 65 | 0 | 0 | 0 | negative control |


## Audits

| Topology | Design | Phase | Passed | Checks | Key/value mismatches | Revision mismatches | Design checks | First failure |
|---|---|---|---|---|---|---|---|---|
| pg-single | `c1` | 1 | true | 12 | 0 | 0 |  |  |
| pg-single | `c1` | 2 | true | 12 | 0 | 0 |  |  |
| pg-single | `c2` | 1 | true | 12 | 0 | 0 |  |  |
| pg-single | `c2` | 2 | true | 12 | 0 | 0 |  |  |
| pg-single | `d1` | 1 | true | 14 | 0 | 0 | a_content_hash_mismatches=0 a_document_revision_mismatches=0 |  |
| pg-single | `n0` | 1 | true | 12 | 0 | 0 |  |  |
| pg-single | `n1` | 1 | true | 12 | 0 | 0 |  |  |
| pg-single | `n2` | 1 | true | 13 | 0 | 0 | a_rolldown_mismatches=0 |  |
| pg-single | `n3` | 1 | true | 13 | 0 | 0 | a_rollup_mismatches=0 |  |
| pg-single | `n4` | 1 | true | 13 | 0 | 0 | a_rollup_mismatches=0 |  |
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

