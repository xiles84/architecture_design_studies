> **OUTDATED (moved 2026-09-12).** Interim preview generated from the PostgreSQL slice while the survey was still running. Superseded by `reports/20260912-small.md`. Kept unedited below this banner.

# Study 01 — results: tree structures under a charity → person → donation schema

Generated from `20260912-small`. Every number here comes from exactly one JSON file in that
directory; the JSON is the source of truth and this file only arranges it.

## What was measured

| | |
|---|---|
| Run id | `20260912-small` |
| Environment | [`host-zenbook-ux5406sa`](../../../docs/environments/host-zenbook-ux5406sa.md) |
| Dataset scale | `small` — 10 charities, 5000 people, 108081 donations |
| Donations per person | mean 21.6, max 496 |
| Dataset seed | 42 (identical data in every cell) |
| Client connections | 8 |
| Duration per query | 10s measured, 3s warmup discarded |

**Engines**

| Topology | Version |
|---|---|
| PostgreSQL, 1 node | `PostgreSQL 17.11 (Debian 17.11-1.pgdg13+2) on x86_64-pc-linux-gnu, compiled by gcc (Debian 14.2.0-19) 14.2.0, 64-bit` |
| YugabyteDB, 1 node (RF=1) | `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit` |

## Correctness gate

Every design must return the same answers, checked against values computed
independently in Go from the generated dataset. Timings are only reported for
cells that passed; a failing cell aborts before it is measured.

| Topology | Design | Checks passed |
|---|---|---|
| PostgreSQL, 1 node | D1 minimal | ✅ 12/12 |
| PostgreSQL, 1 node | D2 indexed | ✅ 12/12 |
| PostgreSQL, 1 node | D3 flat+FK | ✅ 12/12 |
| PostgreSQL, 1 node | D8 flat−FK | ✅ 12/12 |
| PostgreSQL, 1 node | D4 rollup/trg | ✅ 12/12 |
| PostgreSQL, 1 node | D5 rollup/app | ✅ 12/12 |
| PostgreSQL, 1 node | D6 embedded | ✅ 12/12 |
| YugabyteDB, 1 node (RF=1) | D1 minimal | ✅ 12/12 |
| YugabyteDB, 1 node (RF=1) | D2 indexed | ✅ 12/12 |

All cells passed.

## The trade-off at a glance

Every design here buys read speed with something — write throughput, bytes on
disk, or index maintenance on the hot path. This table puts the purchase and the
price in the same row.

*Read score* is the **geometric** mean across all twelve questions. These
throughputs span four orders of magnitude, so an arithmetic mean would be decided
entirely by the fastest query and say nothing about the other eleven.

### PostgreSQL, 1 node

| Design | Read score | vs D1 minimal | Insert/s | vs D1 minimal | Indexes on `donation` | Storage |
|---|---:|---:|---:|---:|---:|---:|
| D1 minimal | 396 | 1.0x | 6.2k | 1.00x | 1 | 12.5 MiB |
| D2 indexed | 5.2k | 13.3x | 6.7k | 1.07x | 3 | 18.1 MiB |
| D3 flat+FK | 8.1k | 20.5x | 5.5k | 0.89x | 5 | 26.3 MiB |
| D8 flat−FK | 8.1k | 20.4x | 5.6k | 0.89x | 5 | 26.3 MiB |
| D4 rollup/trg | 29k | 73.8x | 2.6k | 0.41x | 4 | 23.2 MiB |
| D5 rollup/app | 28k | 69.9x | 2.1k | 0.33x | 4 | 23.6 MiB |
| D6 embedded | 951 | 2.4x | 3.6k | 0.57x | 0 | 21.0 MiB |

### YugabyteDB, 1 node (RF=1)

| Design | Read score | vs D1 minimal | Insert/s | vs D1 minimal | Indexes on `donation` | Storage |
|---|---:|---:|---:|---:|---:|---:|
| D1 minimal | 40.8 | 1.0x | 614 | 1.00x | — | — |
| D2 indexed | 325 | 8.0x | 496 | 0.81x | — | — |


## Read throughput by question (ops/s, higher is better)

### PostgreSQL, 1 node

| Question | D1 minimal | D2 indexed | D3 flat+FK | D8 flat−FK | D4 rollup/trg | D5 rollup/app | D6 embedded | 
|---|---:|---:|---:|---:|---:|---:|---:|
| Information of the last donation (global) | 140 | 39k | 40k | 38k | 27k | 28k | 35.9 | 
| Information of the last donation (one charity) | 231 | 19k | 39k | 39k | 29k | 30k | 310 | 
| Who donates the most | 227 | 305 | 563 | 559 | 34k | 35k | 444 | 
| Top-10 donor leaderboard | 219 | 306 | 480 | 532 | 29k | 30k | 422 | 
| Who donated last | 226 | 19k | 22k | 32k | 25k | 27k | 674 | 
| First and last donation of a person | 434 | 34k | 30k | 28k | 37k | 24k | 29k | 
| Total donated (global) | 373 | 378 | 286 | 277 | 34k | 35k | 50.5 | 
| Total donated (one charity) | 262 | 397 | 2.5k | 2.1k | 36k | 33k | 431 | 
| A person's 20 most recent donations | 409 | 34k | 28k | 25k | 28k | 20k | 19k | 
| Donation by id (point lookup) | 26k | 38k | 36k | 36k | 35k | 34k | 5.4k | 
| How many donations a person made | 404 | 38k | 33k | 32k | 37k | 36k | 35k | 
| Charity activity feed (last 50, with names) | 222 | 1.3k | 12k | 12k | 12k | 13k | 177 | 

### YugabyteDB, 1 node (RF=1)

| Question | D1 minimal | D2 indexed | 
|---|---:|---:|
| Information of the last donation (global) | 17.7 | 2.6k | 
| Information of the last donation (one charity) | 17.4 | 192 | 
| Who donates the most | 25.2 | 31.9 | 
| Top-10 donor leaderboard | 24.7 | 27.7 | 
| Who donated last | 26.1 | 215 | 
| First and last donation of a person | 45.1 | 3.5k | 
| Total donated (global) | 30.2 | 28.8 | 
| Total donated (one charity) | 25.3 | 29.4 | 
| A person's 20 most recent donations | 44.9 | 2.4k | 
| Donation by id (point lookup) | 4.1k | 4.0k | 
| How many donations a person made | 36.1 | 3.4k | 
| Charity activity feed (last 50, with names) | 18.7 | 148 | 

## Read tail latency (p99, ms, lower is better)

Means hide the cases a design handles badly; this is where design decisions
usually show up first.

### PostgreSQL, 1 node

| Question | D1 minimal | D2 indexed | D3 flat+FK | D8 flat−FK | D4 rollup/trg | D5 rollup/app | D6 embedded | 
|---|---:|---:|---:|---:|---:|---:|---:|
| Information of the last donation (global) | 101 | 0.47 | 0.41 | 0.56 | 1.11 | 0.82 | 296 | 
| Information of the last donation (one charity) | 94 | 0.66 | 0.40 | 0.41 | 1.16 | 0.81 | 97 | 
| Who donates the most | 94 | 98 | 85 | 86 | 0.93 | 0.69 | 90 | 
| Top-10 donor leaderboard | 96 | 99 | 86 | 86 | 1.08 | 0.78 | 90 | 
| Who donated last | 94 | 0.70 | 3.05 | 0.60 | 1.32 | 1.12 | 83 | 
| First and last donation of a person | 84 | 0.46 | 0.71 | 1.20 | 0.66 | 2.86 | 0.73 | 
| Total donated (global) | 87 | 87 | 87 | 87 | 0.68 | 0.65 | 278 | 
| Total donated (one charity) | 90 | 93 | 73 | 70 | 0.66 | 1.25 | 89 | 
| A person's 20 most recent donations | 86 | 0.45 | 0.78 | 1.73 | 0.94 | 3.03 | 0.96 | 
| Donation by id (point lookup) | 1.82 | 0.45 | 0.59 | 0.61 | 0.68 | 0.74 | 65 | 
| How many donations a person made | 85 | 0.41 | 0.64 | 0.77 | 0.71 | 0.75 | 0.65 | 
| Charity activity feed (last 50, with names) | 95 | 79 | 1.60 | 1.45 | 1.49 | 1.33 | 182 | 

### YugabyteDB, 1 node (RF=1)

| Question | D1 minimal | D2 indexed | 
|---|---:|---:|
| Information of the last donation (global) | 599 | 67 | 
| Information of the last donation (one charity) | 694 | 91 | 
| Who donates the most | 395 | 794 | 
| Top-10 donor leaderboard | 394 | 804 | 
| Who donated last | 384 | 89 | 
| First and last donation of a person | 289 | 65 | 
| Total donated (global) | 302 | 378 | 
| Total donated (one charity) | 395 | 799 | 
| A person's 20 most recent donations | 282 | 68 | 
| Donation by id (point lookup) | 60 | 63 | 
| How many donations a person made | 392 | 43 | 
| Charity activity feed (last 50, with names) | 691 | 106 | 

## Write throughput (ops/s, higher is better)

`insert` appends a donation, `update` corrects one amount, `delete` removes one
donation. Reads are only half the story: every design that made a read cheap
paid for it somewhere in this table.

### PostgreSQL, 1 node

| Operation | D1 minimal | D2 indexed | D3 flat+FK | D8 flat−FK | D4 rollup/trg | D5 rollup/app | D6 embedded | 
|---|---:|---:|---:|---:|---:|---:|---:|
| insert | 6.2k | 6.7k | 5.5k | 5.6k | 2.6k | 2.1k | 3.6k | 
| update | 6.8k | 6.8k | 5.9k | 5.0k | 3.1k | 6.5k | 305 | 
| delete | 7.2k | 7.1k | 6.2k | 3.0k | 2.5k | 6.0k | 28.5 | 
| insert p99 (ms) | 2.08 | 1.71 | 2.33 | 2.45 | 16 | 16 | 28 | 

### YugabyteDB, 1 node (RF=1)

| Operation | D1 minimal | D2 indexed | 
|---|---:|---:|
| insert | 614 | 496 | 
| update | 2.3k | 2.7k | 
| delete | 962 | 542 | 
| insert p99 (ms) | 74 | 78 | 

## Storage footprint and load time

Read speed bought with redundancy is paid for in bytes and in how long the
data takes to land. Sizes are PostgreSQL only — YugabyteDB stores data in
DocDB rather than PostgreSQL heap files, so `pg_total_relation_size()` does not
describe it.

### PostgreSQL, 1 node

| Design | Tables+indexes | Index bytes | Indexes on `donation` | Load total | COPY | Index build | Rollup |
|---|---:|---:|---:|---:|---:|---:|---:|
| D1 minimal | 12.5 MiB | 3.8 MiB | 1 | 0.4s | 0.4s | 0.0s | 0.0s |
| D2 indexed | 18.1 MiB | 9.5 MiB | 3 | 0.5s | 0.4s | 0.1s | 0.0s |
| D3 flat+FK | 26.3 MiB | 16.1 MiB | 5 | 0.7s | 0.5s | 0.2s | 0.0s |
| D8 flat−FK | 26.3 MiB | 16.1 MiB | 5 | 0.4s | 0.1s | 0.2s | 0.0s |
| D4 rollup/trg | 23.2 MiB | 13.3 MiB | 4 | 1.2s | 0.7s | 0.2s | 0.2s |
| D5 rollup/app | 23.6 MiB | 13.3 MiB | 4 | 0.9s | 0.6s | 0.1s | 0.1s |
| D6 embedded | 21.0 MiB | 14.3 MiB | 0 | 0.5s | 0.1s | 0.4s | 0.0s |

### YugabyteDB, 1 node (RF=1)

| Design | Tables+indexes | Index bytes | Indexes on `donation` | Load total | COPY | Index build | Rollup |
|---|---:|---:|---:|---:|---:|---:|---:|
| D1 minimal | — | — | — | 3.2s | 1.7s | 0.0s | 0.0s |
| D2 indexed | — | — | — | 10.5s | 2.3s | 6.4s | 0.0s |

**Indexes on `donation`** counts the index entries that must be maintained on every
donation INSERT — including the primary key. It is the structural driver of write
cost, and the number to look at when a write result surprises you.

## Controlled pairs — one decision at a time

Each pair below differs by exactly one design decision, so a measured
difference has exactly one explanation. This is the part of the report worth
reading if you read nothing else.

### This run's error bar

**D4 and D5 read identically by construction** — same columns, byte-identical
query SQL. Every read difference measured between them is therefore noise, which
makes that pair a free calibration of this run's own error bar:

| | |
|---|---|
| Control comparisons | 12 |
| Typical (median) disagreement | **1.04x** |
| Worst disagreement | **1.56x** |

Differences below the worst control disagreement (**1.56x**) are suppressed in the
tables that follow. A gap smaller than the largest error two *identical* designs
showed cannot be attributed to a design decision.

> ⚠️ A worst-case control disagreement of 1.56x is high. It means this run can
> only support order-of-magnitude claims, not tens-of-percent ones. Comparisons
> that matter should be re-measured with `-trials N`, which reports a median
> across repeats and a per-cell spread.

### D1 minimal → D2 indexed

*What do secondary indexes alone buy? (identical SQL)*

| Topology | Operation | D1 minimal | D2 indexed | Change |
|---|---|---:|---:|---|
| PostgreSQL, 1 node | Information of the last donation (global) | 140 | 39k | **279.6x faster** |
| PostgreSQL, 1 node | Information of the last donation (one charity) | 231 | 19k | **84.0x faster** |
| PostgreSQL, 1 node | Who donated last | 226 | 19k | **84.5x faster** |
| PostgreSQL, 1 node | First and last donation of a person | 434 | 34k | **79.2x faster** |
| PostgreSQL, 1 node | A person's 20 most recent donations | 409 | 34k | **82.6x faster** |
| PostgreSQL, 1 node | How many donations a person made | 404 | 38k | **92.9x faster** |
| PostgreSQL, 1 node | Charity activity feed (last 50, with names) | 222 | 1.3k | **5.8x faster** |
| PostgreSQL, 1 node | write: insert | 6.2k | 6.7k | no measurable change (1.07x, inside the 1.56x error bar) |
| PostgreSQL, 1 node | write: update | 6.8k | 6.8k | no measurable change (1.00x, inside the 1.56x error bar) |
| PostgreSQL, 1 node | write: delete | 7.2k | 7.1k | no measurable change (0.99x, inside the 1.56x error bar) |
| YugabyteDB, 1 node (RF=1) | Information of the last donation (global) | 17.7 | 2.6k | **148.0x faster** |
| YugabyteDB, 1 node (RF=1) | Information of the last donation (one charity) | 17.4 | 192 | **11.0x faster** |
| YugabyteDB, 1 node (RF=1) | Who donated last | 26.1 | 215 | **8.2x faster** |
| YugabyteDB, 1 node (RF=1) | First and last donation of a person | 45.1 | 3.5k | **78.5x faster** |
| YugabyteDB, 1 node (RF=1) | A person's 20 most recent donations | 44.9 | 2.4k | **53.1x faster** |
| YugabyteDB, 1 node (RF=1) | How many donations a person made | 36.1 | 3.4k | **94.3x faster** |
| YugabyteDB, 1 node (RF=1) | Charity activity feed (last 50, with names) | 18.7 | 148 | **7.9x faster** |
| YugabyteDB, 1 node (RF=1) | write: insert | 614 | 496 | no measurable change (0.81x, inside the 1.56x error bar) |
| YugabyteDB, 1 node (RF=1) | write: update | 2.3k | 2.7k | no measurable change (1.18x, inside the 1.56x error bar) |
| YugabyteDB, 1 node (RF=1) | write: delete | 962 | 542 | **1.8x slower** |

Read differences below this run's measured 1.56x error bar are omitted. **Write
rows are always shown**, even when the change is inside the error bar — writes are
the price every one of these designs pays for its read speed, and a hidden price
tag reads as no price at all.

### D2 indexed → D3 flat+FK

*What does denormalising the grandparent key buy?*

| Topology | Operation | D2 indexed | D3 flat+FK | Change |
|---|---|---:|---:|---|
| PostgreSQL, 1 node | Information of the last donation (one charity) | 19k | 39k | **2.0x faster** |
| PostgreSQL, 1 node | Who donates the most | 305 | 563 | **1.8x faster** |
| PostgreSQL, 1 node | Top-10 donor leaderboard | 306 | 480 | **1.6x faster** |
| PostgreSQL, 1 node | Total donated (one charity) | 397 | 2.5k | **6.4x faster** |
| PostgreSQL, 1 node | Charity activity feed (last 50, with names) | 1.3k | 12k | **9.2x faster** |
| PostgreSQL, 1 node | write: insert | 6.7k | 5.5k | no measurable change (0.83x, inside the 1.56x error bar) |
| PostgreSQL, 1 node | write: update | 6.8k | 5.9k | no measurable change (0.86x, inside the 1.56x error bar) |
| PostgreSQL, 1 node | write: delete | 7.1k | 6.2k | no measurable change (0.87x, inside the 1.56x error bar) |

Read differences below this run's measured 1.56x error bar are omitted. **Write
rows are always shown**, even when the change is inside the error bar — writes are
the price every one of these designs pays for its read speed, and a hidden price
tag reads as no price at all.

### D3 flat+FK → D8 flat−FK

*What does enforcing foreign keys cost?*

| Topology | Operation | D3 flat+FK | D8 flat−FK | Change |
|---|---|---:|---:|---|
| PostgreSQL, 1 node | write: insert | 5.5k | 5.6k | no measurable change (1.01x, inside the 1.56x error bar) |
| PostgreSQL, 1 node | write: update | 5.9k | 5.0k | no measurable change (0.85x, inside the 1.56x error bar) |
| PostgreSQL, 1 node | write: delete | 6.2k | 3.0k | **2.1x slower** |

Read differences below this run's measured 1.56x error bar are omitted. **Write
rows are always shown**, even when the change is inside the error bar — writes are
the price every one of these designs pays for its read speed, and a hidden price
tag reads as no price at all.

### D3 flat+FK → D4 rollup/trg

*What do consolidated aggregates buy and cost?*

| Topology | Operation | D3 flat+FK | D4 rollup/trg | Change |
|---|---|---:|---:|---|
| PostgreSQL, 1 node | Who donates the most | 563 | 34k | **60.0x faster** |
| PostgreSQL, 1 node | Top-10 donor leaderboard | 480 | 29k | **59.9x faster** |
| PostgreSQL, 1 node | Total donated (global) | 286 | 34k | **120.2x faster** |
| PostgreSQL, 1 node | Total donated (one charity) | 2.5k | 36k | **14.4x faster** |
| PostgreSQL, 1 node | write: insert | 5.5k | 2.6k | **2.2x slower** |
| PostgreSQL, 1 node | write: update | 5.9k | 3.1k | **1.9x slower** |
| PostgreSQL, 1 node | write: delete | 6.2k | 2.5k | **2.5x slower** |

Read differences below this run's measured 1.56x error bar are omitted. **Write
rows are always shown**, even when the change is inside the error bar — writes are
the price every one of these designs pays for its read speed, and a hidden price
tag reads as no price at all.

### D4 rollup/trg → D5 rollup/app

*Does it matter whether a trigger or the app maintains them?*

| Topology | Operation | D4 rollup/trg | D5 rollup/app | Change |
|---|---|---:|---:|---|
| PostgreSQL, 1 node | write: insert | 2.6k | 2.1k | no measurable change (0.80x, inside the 1.56x error bar) |
| PostgreSQL, 1 node | write: update | 3.1k | 6.5k | **2.1x faster** |
| PostgreSQL, 1 node | write: delete | 2.5k | 6.0k | **2.4x faster** |

Read differences below this run's measured 1.56x error bar are omitted. **Write
rows are always shown**, even when the change is inside the error bar — writes are
the price every one of these designs pays for its read speed, and a hidden price
tag reads as no price at all.

### D3 flat+FK → D6 embedded

*What does embedding ALL children in the parent buy and cost?*

| Topology | Operation | D3 flat+FK | D6 embedded | Change |
|---|---|---:|---:|---|
| PostgreSQL, 1 node | Information of the last donation (global) | 40k | 35.9 | **1117.4x slower** |
| PostgreSQL, 1 node | Information of the last donation (one charity) | 39k | 310 | **126.6x slower** |
| PostgreSQL, 1 node | Who donated last | 22k | 674 | **32.4x slower** |
| PostgreSQL, 1 node | Total donated (global) | 286 | 50.5 | **5.7x slower** |
| PostgreSQL, 1 node | Total donated (one charity) | 2.5k | 431 | **5.9x slower** |
| PostgreSQL, 1 node | Donation by id (point lookup) | 36k | 5.4k | **6.7x slower** |
| PostgreSQL, 1 node | Charity activity feed (last 50, with names) | 12k | 177 | **66.6x slower** |
| PostgreSQL, 1 node | write: insert | 5.5k | 3.6k | no measurable change (0.64x, inside the 1.56x error bar) |
| PostgreSQL, 1 node | write: update | 5.9k | 305 | **19.3x slower** |
| PostgreSQL, 1 node | write: delete | 6.2k | 28.5 | **216.8x slower** |

Read differences below this run's measured 1.56x error bar are omitted. **Write
rows are always shown**, even when the change is inside the error bar — writes are
the price every one of these designs pays for its read speed, and a hidden price
tag reads as no price at all.


---

Query plans for every cell are in `20260912-small/<topology>/plans/<design>.txt`,
captured with `EXPLAIN (ANALYZE, BUFFERS)` on PostgreSQL and
`EXPLAIN (ANALYZE, DIST)` on YugabyteDB.

## Conclusions and analysis provenance

This file is **generated from the measurements** and contains no interpretation.
Conclusions live in separate signed analyses, so that several analysts — including
different AI models — can read the same numbers and each record what they make of
them. Where two analyses disagree, that disagreement is itself a finding and is
left visible rather than resolved by editing one of them.

| | |
|---|---|
| Run id | `20260912-small` |
| Result files | 9 |
| **Inputs digest** | `d6b42ad83d857659` |

The digest is a content hash of every result file in the run. **Quote it in any
analysis.** A run id alone is not enough to identify what was analysed — a run can
be extended with extra cells afterwards — so the digest is what tells a later
analyst whether they are looking at the same data someone else already wrote about.

### Analyses of this run

**None yet.** This run has measurements but no interpretation.

### Adding an analysis

Copy `reports/analyses/TEMPLATE.md`, fill in the frontmatter, and write what you
conclude. Then regenerate this report so the table above picks it up.

Before writing one, check the table: if an analysis already exists for digest
`d6b42ad83d857659`, read it first. Add a new analysis to **disagree, extend, or bring a
different perspective** — not to restate what is already there. Never edit another
analyst's file; write your own and reference theirs.
