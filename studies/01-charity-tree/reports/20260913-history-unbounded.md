# Study 01 — results: tree structures under a charity → person → donation schema

Generated from `20260913-history-unbounded`. Every number here comes from exactly one JSON file in that
directory; the JSON is the source of truth and this file only arranges it.

## What was measured

| | |
|---|---|
| Run id | `20260913-history-unbounded` |
| Environment | [`host-zenbook-ux5406sa`](../../../docs/environments/host-zenbook-ux5406sa.md) |
| Dataset scale | `small` — 10 charities, 5000 people, 108081 donations |
| Donations per person | mean 21.6, max 496 |
| Dataset seed | 42 (identical data in every cell) |
| Client connections | 8 |
| Duration per query | 8s measured, 3s warmup discarded |

**Engines**

| Topology | Version |
|---|---|
| PostgreSQL, 1 node | `PostgreSQL 17.11 (Debian 17.11-1.pgdg13+2) on x86_64-pc-linux-gnu, compiled by gcc (Debian 14.2.0-19) 14.2.0, 64-bit` |
| YugabyteDB, 3 nodes (RF=3) | `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit` |

## Correctness gate

Every design must return the same answers, checked against values computed
independently in Go from the generated dataset. Timings are only reported for
cells that passed; a failing cell aborts before it is measured.

| Topology | Design | Checks passed |
|---|---|---|
| PostgreSQL, 1 node | D3 flat+FK | ✅ 12/12 |
| PostgreSQL, 1 node | D4 rollup/trg | ✅ 12/12 |
| PostgreSQL, 1 node | D6 embedded | ✅ 12/12 |
| PostgreSQL, 1 node | D9 hybrid | ✅ 12/12 |
| YugabyteDB, 3 nodes (RF=3) | D3 flat+FK | ✅ 12/12 |
| YugabyteDB, 3 nodes (RF=3) | D6 embedded | ✅ 12/12 |
| YugabyteDB, 3 nodes (RF=3) | D7 yb-coloc | ✅ 12/12 |

All cells passed.

## The trade-off at a glance

Every design here buys read speed with something — write throughput, bytes on
disk, or index maintenance on the hot path. This table puts the purchase and the
price in the same row.

*Read score* is the **geometric** mean across all twelve questions. These
throughputs span four orders of magnitude, so an arithmetic mean would be decided
entirely by the fastest query and say nothing about the other eleven.

### PostgreSQL, 1 node

| Design | Read score | vs D3 flat+FK | Insert/s | vs D3 flat+FK | Indexes on `donation` | Storage |
|---|---:|---:|---:|---:|---:|---:|
| D3 flat+FK | 8.2k | 1.0x | 6.2k | 1.00x | 5 | 25.8 MiB |
| D4 rollup/trg | 32k | 3.9x | 2.9k | 0.46x | 4 | 23.6 MiB |
| D6 embedded | 866 | 0.1x | 3.9k | 0.62x | 0 | 20.7 MiB |
| D9 hybrid | 8.4k | 1.0x | 5.0k | 0.81x | 5 | 29.5 MiB |

### YugabyteDB, 3 nodes (RF=3)

| Design | Read score | vs D3 flat+FK | Insert/s | vs D3 flat+FK | Indexes on `donation` | Storage |
|---|---:|---:|---:|---:|---:|---:|
| D3 flat+FK | 696 | 1.0x | 264 | 1.00x | — | — |
| D6 embedded | 482 | 0.7x | 249 | 0.94x | — | — |
| D7 yb-coloc | 570 | 0.8x | 292 | 1.10x | — | — |


## Read throughput by question (ops/s, higher is better)

### PostgreSQL, 1 node

| Question | D3 flat+FK | D4 rollup/trg | D6 embedded | D9 hybrid | 
|---|---:|---:|---:|---:|
| Information of the last donation (global) | 33k | 30k | 27.8 | 37k | 
| Information of the last donation (one charity) | 34k | 34k | 272 | 36k | 
| Who donates the most | 484 | 37k | 374 | 503 | 
| Top-10 donor leaderboard | 512 | 34k | 340 | 527 | 
| Who donated last | 34k | 34k | 539 | 36k | 
| First and last donation of a person | 31k | 39k | 36k | 33k | 
| Total donated (global) | 256 | 37k | 38.1 | 260 | 
| Total donated (one charity) | 2.3k | 41k | 350 | 2.4k | 
| A person's 20 most recent donations | 32k | 31k | 22k | 27k | 
| Donation by id (point lookup) | 36k | 32k | 6.2k | 36k | 
| How many donations a person made | 36k | 38k | 35k | 36k | 
| Charity activity feed (last 50, with names) | 14k | 13k | 149 | 13k | 

### YugabyteDB, 3 nodes (RF=3)

| Question | D3 flat+FK | D6 embedded | D7 yb-coloc | 
|---|---:|---:|---:|
| Information of the last donation (global) | 2.7k | 58.2 | 2.4k | 
| Information of the last donation (one charity) | 2.4k | 423 | 2.2k | 
| Who donates the most | 64.5 | 290 | 57.6 | 
| Top-10 donor leaderboard | 67.7 | 293 | 50.3 | 
| Who donated last | 1.8k | 431 | 1.5k | 
| First and last donation of a person | 3.9k | 3.9k | 3.8k | 
| Total donated (global) | 37.6 | 37.9 | 33.5 | 
| Total donated (one charity) | 324 | 303 | 61.5 | 
| A person's 20 most recent donations | 2.2k | 3.4k | 3.8k | 
| Donation by id (point lookup) | 3.7k | 1.7k | 2.5k | 
| How many donations a person made | 4.0k | 4.1k | 4.0k | 
| Charity activity feed (last 50, with names) | 172 | 169 | 172 | 

## Read tail latency (p99, ms, lower is better)

Means hide the cases a design handles badly; this is where design decisions
usually show up first.

### PostgreSQL, 1 node

| Question | D3 flat+FK | D4 rollup/trg | D6 embedded | D9 hybrid | 
|---|---:|---:|---:|---:|
| Information of the last donation (global) | 0.84 | 0.53 | 477 | 0.51 | 
| Information of the last donation (one charity) | 0.61 | 0.49 | 104 | 0.49 | 
| Who donates the most | 91 | 0.47 | 96 | 90 | 
| Top-10 donor leaderboard | 89 | 0.48 | 97 | 89 | 
| Who donated last | 0.52 | 0.45 | 90 | 0.43 | 
| First and last donation of a person | 0.70 | 0.48 | 0.45 | 0.47 | 
| Total donated (global) | 93 | 0.51 | 380 | 94 | 
| Total donated (one charity) | 75 | 0.44 | 98 | 75 | 
| A person's 20 most recent donations | 0.50 | 0.51 | 0.69 | 0.54 | 
| Donation by id (point lookup) | 0.49 | 0.79 | 65 | 0.53 | 
| How many donations a person made | 0.44 | 0.62 | 0.50 | 0.48 | 
| Charity activity feed (last 50, with names) | 0.89 | 0.93 | 197 | 0.94 | 

### YugabyteDB, 3 nodes (RF=3)

| Question | D3 flat+FK | D6 embedded | D7 yb-coloc | 
|---|---:|---:|---:|
| Information of the last donation (global) | 52 | 192 | 52 | 
| Information of the last donation (one charity) | 53 | 75 | 51 | 
| Who donates the most | 361 | 89 | 434 | 
| Top-10 donor leaderboard | 369 | 90 | 469 | 
| Who donated last | 54 | 74 | 53 | 
| First and last donation of a person | 52 | 53 | 51 | 
| Total donated (global) | 284 | 292 | 298 | 
| Total donated (one charity) | 86 | 90 | 467 | 
| A person's 20 most recent donations | 49 | 56 | 51 | 
| Donation by id (point lookup) | 51 | 58 | 59 | 
| How many donations a person made | 52 | 53 | 50 | 
| Charity activity feed (last 50, with names) | 92 | 173 | 85 | 

## Write throughput (ops/s, higher is better)

`insert` appends a donation, `update` corrects one amount, `delete` removes one
donation. Reads are only half the story: every design that made a read cheap
paid for it somewhere in this table.

A dash in the p99.9 or p99.99 row means **too few samples to support that
percentile**, not a missing measurement: a quantile is only reported when at
least ten samples sit beyond it (10 000 samples for p99.9, 100 000 for p99.99).
Below that a "p99.9" would be the single worst request relabelled, which reads
as far more authoritative than it is. `max` is always shown, and for a
low-throughput cell it is the only honest tail figure available.

> **These tails are optimistic, by construction.** The harness is closed-loop:
> each worker waits for its own response before issuing the next request. When
> the server stalls, the load generator stalls with it, so the stall is recorded
> once instead of being charged to every request that would have arrived during
> it — the coordinated-omission effect. It barely touches p50, biases p99 a
> little, and biases p99.9 and beyond substantially. Treat deep tails here as a
> floor on what a real open-loop client would see, and compare them **between
> designs** (all of which pay the same bias) rather than as absolute SLO figures.

### PostgreSQL, 1 node

| Operation | D3 flat+FK | D4 rollup/trg | D6 embedded | D9 hybrid | 
|---|---:|---:|---:|---:|
| insert donation | 6.2k | 2.9k | 3.9k | 5.0k | 
| correct a donation amount | 6.5k | 3.3k | 434 | 4.4k | 
| remove one donation | 7.4k | 2.7k | 1.2k | 4.7k | 
| donor edits their profile | 7.1k | 7.1k | 6.6k | 7.2k | 
| erase a donor and all their donations | 6.7k | 519 | 7.5k | 3.9k | 
| insert p99 (ms) | 2.00 | 14 | 14 | 12 | 
| insert p99.9 (ms) | 5.12 | 24 | 47 | 18 | 
| insert p99.99 (ms) | — | — | — | — | 
| insert max (ms) | 18 | 63 | 3038 | 41 | 

### YugabyteDB, 3 nodes (RF=3)

| Operation | D3 flat+FK | D6 embedded | D7 yb-coloc | 
|---|---:|---:|---:|
| insert donation | 264 | 249 | 292 | 
| correct a donation amount | — | — | — | 
| remove one donation | — | — | — | 
| donor edits their profile | 2.3k | 2.1k | 2.2k | 
| erase a donor and all their donations | 114 | 273 | 125 | 
| insert p99 (ms) | 71 | 104 | 62 | 
| insert p99.9 (ms) | — | — | — | 
| insert p99.99 (ms) | — | — | — | 
| insert max (ms) | 94 | 192 | 74 | 

## Storage footprint and load time

Read speed bought with redundancy is paid for in bytes and in how long the
data takes to land. Sizes are PostgreSQL only — YugabyteDB stores data in
DocDB rather than PostgreSQL heap files, so `pg_total_relation_size()` does not
describe it.

### PostgreSQL, 1 node

| Design | Tables+indexes | Index bytes | Indexes on `donation` | Load total | COPY | Index build | Rollup |
|---|---:|---:|---:|---:|---:|---:|---:|
| D3 flat+FK | 25.8 MiB | 16.1 MiB | 5 | 0.7s | 0.5s | 0.1s | 0.0s |
| D4 rollup/trg | 23.6 MiB | 13.3 MiB | 4 | 0.8s | 0.4s | 0.1s | 0.1s |
| D6 embedded | 20.7 MiB | 14.3 MiB | 0 | 0.6s | 0.1s | 0.4s | 0.0s |
| D9 hybrid | 29.5 MiB | 16.1 MiB | 5 | 0.8s | 0.6s | 0.2s | 0.0s |

### YugabyteDB, 3 nodes (RF=3)

| Design | Tables+indexes | Index bytes | Indexes on `donation` | Load total | COPY | Index build | Rollup |
|---|---:|---:|---:|---:|---:|---:|---:|
| D3 flat+FK | — | — | — | 27.1s | 1.3s | 24.2s | 0.0s |
| D6 embedded | — | — | — | 8.9s | 0.3s | 6.6s | 0.0s |
| D7 yb-coloc | — | — | — | 25.5s | 1.5s | 20.8s | 0.0s |

**Indexes on `donation`** counts the index entries that must be maintained on every
donation INSERT — including the primary key. It is the structural driver of write
cost, and the number to look at when a write result surprises you.

## Controlled pairs — one decision at a time

Each pair below differs by exactly one design decision, so a measured
difference has exactly one explanation. This is the part of the report worth
reading if you read nothing else.

### This run's error bar

The D4/D5 control pair is not present in this run, so the threshold below is a
conventional 1.30x rather than one measured from the data.

### D3 flat+FK → D4 rollup/trg

*What do consolidated aggregates buy and cost?*

| Topology | Operation | D3 flat+FK | D4 rollup/trg | Change |
|---|---|---:|---:|---|
| PostgreSQL, 1 node | Who donates the most | 484 | 37k | **77.2x faster** |
| PostgreSQL, 1 node | Top-10 donor leaderboard | 512 | 34k | **65.5x faster** |
| PostgreSQL, 1 node | Total donated (global) | 256 | 37k | **144.8x faster** |
| PostgreSQL, 1 node | Total donated (one charity) | 2.3k | 41k | **17.9x faster** |
| PostgreSQL, 1 node | insert donation | 6.2k | 2.9k | **2.2x slower** |
| PostgreSQL, 1 node | correct a donation amount | 6.5k | 3.3k | **2.0x slower** |
| PostgreSQL, 1 node | remove one donation | 7.4k | 2.7k | **2.7x slower** |
| PostgreSQL, 1 node | donor edits their profile | 7.1k | 7.1k | no measurable change (1.00x, inside the 1.30x error bar) |
| PostgreSQL, 1 node | erase a donor and all their donations | 6.7k | 519 | **12.9x slower** |

Read differences below this run's measured 1.30x error bar are omitted. **Write
rows are always shown**, even when the change is inside the error bar — writes are
the price every one of these designs pays for its read speed, and a hidden price
tag reads as no price at all.

### D3 flat+FK → D6 embedded

*What does embedding ALL children in the parent buy and cost?*

| Topology | Operation | D3 flat+FK | D6 embedded | Change |
|---|---|---:|---:|---|
| PostgreSQL, 1 node | Information of the last donation (global) | 33k | 27.8 | **1181.7x slower** |
| PostgreSQL, 1 node | Information of the last donation (one charity) | 34k | 272 | **125.0x slower** |
| PostgreSQL, 1 node | Top-10 donor leaderboard | 512 | 340 | **1.5x slower** |
| PostgreSQL, 1 node | Who donated last | 34k | 539 | **62.2x slower** |
| PostgreSQL, 1 node | Total donated (global) | 256 | 38.1 | **6.7x slower** |
| PostgreSQL, 1 node | Total donated (one charity) | 2.3k | 350 | **6.6x slower** |
| PostgreSQL, 1 node | A person's 20 most recent donations | 32k | 22k | **1.5x slower** |
| PostgreSQL, 1 node | Donation by id (point lookup) | 36k | 6.2k | **5.8x slower** |
| PostgreSQL, 1 node | Charity activity feed (last 50, with names) | 14k | 149 | **92.4x slower** |
| PostgreSQL, 1 node | insert donation | 6.2k | 3.9k | **1.6x slower** |
| PostgreSQL, 1 node | correct a donation amount | 6.5k | 434 | **15.1x slower** |
| PostgreSQL, 1 node | remove one donation | 7.4k | 1.2k | **6.2x slower** |
| PostgreSQL, 1 node | donor edits their profile | 7.1k | 6.6k | no measurable change (0.93x, inside the 1.30x error bar) |
| PostgreSQL, 1 node | erase a donor and all their donations | 6.7k | 7.5k | no measurable change (1.12x, inside the 1.30x error bar) |
| YugabyteDB, 3 nodes (RF=3) | Information of the last donation (global) | 2.7k | 58.2 | **46.1x slower** |
| YugabyteDB, 3 nodes (RF=3) | Information of the last donation (one charity) | 2.4k | 423 | **5.7x slower** |
| YugabyteDB, 3 nodes (RF=3) | Who donates the most | 64.5 | 290 | **4.5x faster** |
| YugabyteDB, 3 nodes (RF=3) | Top-10 donor leaderboard | 67.7 | 293 | **4.3x faster** |
| YugabyteDB, 3 nodes (RF=3) | Who donated last | 1.8k | 431 | **4.1x slower** |
| YugabyteDB, 3 nodes (RF=3) | A person's 20 most recent donations | 2.2k | 3.4k | **1.6x faster** |
| YugabyteDB, 3 nodes (RF=3) | Donation by id (point lookup) | 3.7k | 1.7k | **2.2x slower** |
| YugabyteDB, 3 nodes (RF=3) | insert donation | 264 | 249 | no measurable change (0.94x, inside the 1.30x error bar) |
| YugabyteDB, 3 nodes (RF=3) | donor edits their profile | 2.3k | 2.1k | no measurable change (0.95x, inside the 1.30x error bar) |
| YugabyteDB, 3 nodes (RF=3) | erase a donor and all their donations | 114 | 273 | **2.4x faster** |

Read differences below this run's measured 1.30x error bar are omitted. **Write
rows are always shown**, even when the change is inside the error bar — writes are
the price every one of these designs pays for its read speed, and a hidden price
tag reads as no price at all.

### D6 embedded → D9 hybrid

*Does bounding the embedded slice fix what embedding broke?*

| Topology | Operation | D6 embedded | D9 hybrid | Change |
|---|---|---:|---:|---|
| PostgreSQL, 1 node | Information of the last donation (global) | 27.8 | 37k | **1336.1x faster** |
| PostgreSQL, 1 node | Information of the last donation (one charity) | 272 | 36k | **131.6x faster** |
| PostgreSQL, 1 node | Who donates the most | 374 | 503 | **1.3x faster** |
| PostgreSQL, 1 node | Top-10 donor leaderboard | 340 | 527 | **1.6x faster** |
| PostgreSQL, 1 node | Who donated last | 539 | 36k | **67.5x faster** |
| PostgreSQL, 1 node | Total donated (global) | 38.1 | 260 | **6.8x faster** |
| PostgreSQL, 1 node | Total donated (one charity) | 350 | 2.4k | **6.9x faster** |
| PostgreSQL, 1 node | Donation by id (point lookup) | 6.2k | 36k | **5.8x faster** |
| PostgreSQL, 1 node | Charity activity feed (last 50, with names) | 149 | 13k | **88.9x faster** |
| PostgreSQL, 1 node | insert donation | 3.9k | 5.0k | no measurable change (1.30x, inside the 1.30x error bar) |
| PostgreSQL, 1 node | correct a donation amount | 434 | 4.4k | **10.1x faster** |
| PostgreSQL, 1 node | remove one donation | 1.2k | 4.7k | **4.0x faster** |
| PostgreSQL, 1 node | donor edits their profile | 6.6k | 7.2k | no measurable change (1.10x, inside the 1.30x error bar) |
| PostgreSQL, 1 node | erase a donor and all their donations | 7.5k | 3.9k | **1.9x slower** |

Read differences below this run's measured 1.30x error bar are omitted. **Write
rows are always shown**, even when the change is inside the error bar — writes are
the price every one of these designs pays for its read speed, and a hidden price
tag reads as no price at all.

### D3 flat+FK → D9 hybrid

*What does a bounded embedded cache buy and cost?*

| Topology | Operation | D3 flat+FK | D9 hybrid | Change |
|---|---|---:|---:|---|
| PostgreSQL, 1 node | insert donation | 6.2k | 5.0k | no measurable change (0.81x, inside the 1.30x error bar) |
| PostgreSQL, 1 node | correct a donation amount | 6.5k | 4.4k | **1.5x slower** |
| PostgreSQL, 1 node | remove one donation | 7.4k | 4.7k | **1.6x slower** |
| PostgreSQL, 1 node | donor edits their profile | 7.1k | 7.2k | no measurable change (1.02x, inside the 1.30x error bar) |
| PostgreSQL, 1 node | erase a donor and all their donations | 6.7k | 3.9k | **1.7x slower** |

Read differences below this run's measured 1.30x error bar are omitted. **Write
rows are always shown**, even when the change is inside the error bar — writes are
the price every one of these designs pays for its read speed, and a hidden price
tag reads as no price at all.

### D3 flat+FK → D7 yb-coloc

*What does sharding children next to their parent buy?*

| Topology | Operation | D3 flat+FK | D7 yb-coloc | Change |
|---|---|---:|---:|---|
| YugabyteDB, 3 nodes (RF=3) | Top-10 donor leaderboard | 67.7 | 50.3 | **1.3x slower** |
| YugabyteDB, 3 nodes (RF=3) | Total donated (one charity) | 324 | 61.5 | **5.3x slower** |
| YugabyteDB, 3 nodes (RF=3) | A person's 20 most recent donations | 2.2k | 3.8k | **1.8x faster** |
| YugabyteDB, 3 nodes (RF=3) | Donation by id (point lookup) | 3.7k | 2.5k | **1.5x slower** |
| YugabyteDB, 3 nodes (RF=3) | insert donation | 264 | 292 | no measurable change (1.10x, inside the 1.30x error bar) |
| YugabyteDB, 3 nodes (RF=3) | donor edits their profile | 2.3k | 2.2k | no measurable change (0.97x, inside the 1.30x error bar) |
| YugabyteDB, 3 nodes (RF=3) | erase a donor and all their donations | 114 | 125 | no measurable change (1.10x, inside the 1.30x error bar) |

Read differences below this run's measured 1.30x error bar are omitted. **Write
rows are always shown**, even when the change is inside the error bar — writes are
the price every one of these designs pays for its read speed, and a hidden price
tag reads as no price at all.


---

Query plans for every cell are in `20260913-history-unbounded/<topology>/plans/<design>.txt`,
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
| Run id | `20260913-history-unbounded` |
| Result files | 7 |
| **Inputs digest** | `bbfb475a6d879410` |

The digest is a content hash of every result file in the run. **Quote it in any
analysis.** A run id alone is not enough to identify what was analysed — a run can
be extended with extra cells afterwards — so the digest is what tells a later
analyst whether they are looking at the same data someone else already wrote about.

### Analyses of this run

| Analyst | Kind | Date | Digest analysed | Status | Headline |
|---|---|---|---|---|---|
| [claude-opus-5](analyses/20260912-study01--claude-opus-5--2026-09-12.md) | ai | 2026-09-12 | `bbfb475a6d879410` ✅ | current | Store totals on the parent only for aggregate questions and only when writes spread across many parents; use indexes for everything else; embedding donations in the donor row never came out ahead on reads. |
| [gpt-6](analyses/20260912-study01--gpt-6--2026-09-12.md) | ai | 2026-09-12 | `bbfb475a6d879410` ✅ | current | Use ordinary indexes for donor reads, copy the charity key when it enables selective charity indexes, and store totals only when their read benefit justifies contention and maintenance. |

### Adding an analysis

Copy `reports/analyses/TEMPLATE.md`, fill in the frontmatter, and write what you
conclude. Then regenerate this report so the table above picks it up.

Before writing one, check the table: if an analysis already exists for digest
`bbfb475a6d879410`, read it first. Add a new analysis to **disagree, extend, or bring a
different perspective** — not to restate what is already there. Never edit another
analyst's file; write your own and reference theirs.
