# Study 01 — results: tree structures under a charity → person → donation schema

Generated from `20260913-writes-isolated`. Every number here comes from exactly one JSON file in that
directory; the JSON is the source of truth and this file only arranges it.

## What was measured

| | |
|---|---|
| Run id | `20260913-writes-isolated` |
| Environment | [`host-zenbook-ux5406sa`](../../../docs/environments/host-zenbook-ux5406sa.md) |
| Dataset scale | `small` — 10 charities, 5000 people, 108081 donations |
| Donations per person | mean 21.6, max 496 |
| Dataset seed | 42 (identical data in every cell) |
| Client connections | 8 |
| Duration per query | 6s measured, 2s warmup discarded |

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
| PostgreSQL, 1 node | D8 flat−FK | ✅ 12/12 |
| PostgreSQL, 1 node | D4 rollup/trg | ✅ 12/12 |
| PostgreSQL, 1 node | D6 embedded | ✅ 12/12 |
| PostgreSQL, 1 node | D9 hybrid | ✅ 12/12 |
| YugabyteDB, 3 nodes (RF=3) | D3 flat+FK | ✅ 12/12 |
| YugabyteDB, 3 nodes (RF=3) | D8 flat−FK | ✅ 12/12 |
| YugabyteDB, 3 nodes (RF=3) | D4 rollup/trg | ✅ 12/12 |
| YugabyteDB, 3 nodes (RF=3) | D6 embedded | ✅ 12/12 |
| YugabyteDB, 3 nodes (RF=3) | D9 hybrid | ✅ 12/12 |

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
| D3 flat+FK | — | — | 6.7k | 1.00x | 5 | 26.0 MiB |
| D8 flat−FK | — | — | 6.3k | 0.95x | 5 | 25.6 MiB |
| D4 rollup/trg | — | — | 2.8k | 0.42x | 4 | 23.2 MiB |
| D6 embedded | — | — | 4.0k | 0.60x | 0 | 20.9 MiB |
| D9 hybrid | — | — | 5.2k | 0.78x | 5 | 29.7 MiB |

### YugabyteDB, 3 nodes (RF=3)

| Design | Read score | vs D3 flat+FK | Insert/s | vs D3 flat+FK | Indexes on `donation` | Storage |
|---|---:|---:|---:|---:|---:|---:|
| D3 flat+FK | — | — | 384 | 1.00x | — | — |
| D8 flat−FK | — | — | 321 | 0.84x | — | — |
| D4 rollup/trg | — | — | 103 | 0.27x | — | — |
| D6 embedded | — | — | 302 | 0.79x | — | — |
| D9 hybrid | — | — | 247 | 0.64x | — | — |


## Read throughput by question (ops/s, higher is better)

### PostgreSQL, 1 node

| Question | D3 flat+FK | D8 flat−FK | D4 rollup/trg | D6 embedded | D9 hybrid | 
|---|---:|---:|---:|---:|---:|

### YugabyteDB, 3 nodes (RF=3)

| Question | D3 flat+FK | D8 flat−FK | D4 rollup/trg | D6 embedded | D9 hybrid | 
|---|---:|---:|---:|---:|---:|

## Read tail latency (p99, ms, lower is better)

Means hide the cases a design handles badly; this is where design decisions
usually show up first.

### PostgreSQL, 1 node

| Question | D3 flat+FK | D8 flat−FK | D4 rollup/trg | D6 embedded | D9 hybrid | 
|---|---:|---:|---:|---:|---:|

### YugabyteDB, 3 nodes (RF=3)

| Question | D3 flat+FK | D8 flat−FK | D4 rollup/trg | D6 embedded | D9 hybrid | 
|---|---:|---:|---:|---:|---:|

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

| Operation | D3 flat+FK | D8 flat−FK | D4 rollup/trg | D6 embedded | D9 hybrid | 
|---|---:|---:|---:|---:|---:|
| insert donation | 6.7k | 6.3k | 2.8k | 4.0k | 5.2k | 
| correct a donation amount | 6.7k | 6.0k | 2.7k | 418 | 4.8k | 
| remove one donation | 7.3k | 6.8k | 1.3k | 457 | 5.3k | 
| donor edits their profile | 6.7k | 6.6k | 6.1k | 7.2k | 7.1k | 
| erase a donor and all their donations | 5.6k | 6.3k | 389 | 7.4k | 3.8k | 
| insert p99 (ms) | 1.97 | 2.07 | 15 | 18 | 2.63 | 
| insert p99.9 (ms) | 5.49 | 5.34 | 27 | 44 | 20 | 
| insert p99.99 (ms) | 21 | 8.85 | — | — | — | 
| insert max (ms) | 24 | 15 | 56 | 2113 | 36 | 

### YugabyteDB, 3 nodes (RF=3)

| Operation | D3 flat+FK | D8 flat−FK | D4 rollup/trg | D6 embedded | D9 hybrid | 
|---|---:|---:|---:|---:|---:|
| insert donation | 384 | 321 | 103 | 302 | 247 | 
| correct a donation amount | 551 | 449 | 178 | 88.3 | 265 | 
| remove one donation | 257 | 280 | 62.1 | 103 | 204 | 
| donor edits their profile | 2.0k | 1.2k | 2.5k | 3.0k | 3.0k | 
| erase a donor and all their donations | 121 | 184 | 5.4 | 366 | 53.4 | 
| insert p99 (ms) | 51 | 58 | 416 | 97 | 58 | 
| insert p99.9 (ms) | — | — | — | — | — | 
| insert p99.99 (ms) | — | — | — | — | — | 
| insert max (ms) | 73 | 87 | 1018 | 209 | 98 | 

## Storage footprint and load time

Read speed bought with redundancy is paid for in bytes and in how long the
data takes to land. Sizes are PostgreSQL only — YugabyteDB stores data in
DocDB rather than PostgreSQL heap files, so `pg_total_relation_size()` does not
describe it.

### PostgreSQL, 1 node

| Design | Tables+indexes | Index bytes | Indexes on `donation` | Load total | COPY | Index build | Rollup |
|---|---:|---:|---:|---:|---:|---:|---:|
| D3 flat+FK | 26.0 MiB | 16.1 MiB | 5 | 0.6s | 0.4s | 0.1s | 0.0s |
| D8 flat−FK | 25.6 MiB | 16.1 MiB | 5 | 0.3s | 0.1s | 0.1s | 0.0s |
| D4 rollup/trg | 23.2 MiB | 13.3 MiB | 4 | 0.7s | 0.4s | 0.1s | 0.1s |
| D6 embedded | 20.9 MiB | 14.3 MiB | 0 | 0.5s | 0.1s | 0.4s | 0.0s |
| D9 hybrid | 29.7 MiB | 16.1 MiB | 5 | 0.7s | 0.5s | 0.1s | 0.0s |

### YugabyteDB, 3 nodes (RF=3)

| Design | Tables+indexes | Index bytes | Indexes on `donation` | Load total | COPY | Index build | Rollup |
|---|---:|---:|---:|---:|---:|---:|---:|
| D3 flat+FK | — | — | — | 24.3s | 1.4s | 21.0s | 0.0s |
| D8 flat−FK | — | — | — | 17.2s | 1.6s | 13.6s | 0.0s |
| D4 rollup/trg | — | — | — | 22.1s | 1.2s | 17.7s | 0.6s |
| D6 embedded | — | — | — | 8.5s | 0.3s | 6.2s | 0.0s |
| D9 hybrid | — | — | — | 17.0s | 1.3s | 12.9s | 0.0s |

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

### D3 flat+FK → D8 flat−FK

*What does enforcing foreign keys cost?*

| Topology | Operation | D3 flat+FK | D8 flat−FK | Change |
|---|---|---:|---:|---|
| PostgreSQL, 1 node | insert donation | 6.7k | 6.3k | no measurable change (0.95x, inside the 1.30x error bar) |
| PostgreSQL, 1 node | correct a donation amount | 6.7k | 6.0k | no measurable change (0.91x, inside the 1.30x error bar) |
| PostgreSQL, 1 node | remove one donation | 7.3k | 6.8k | no measurable change (0.93x, inside the 1.30x error bar) |
| PostgreSQL, 1 node | donor edits their profile | 6.7k | 6.6k | no measurable change (0.99x, inside the 1.30x error bar) |
| PostgreSQL, 1 node | erase a donor and all their donations | 5.6k | 6.3k | no measurable change (1.13x, inside the 1.30x error bar) |
| YugabyteDB, 3 nodes (RF=3) | insert donation | 384 | 321 | no measurable change (0.84x, inside the 1.30x error bar) |
| YugabyteDB, 3 nodes (RF=3) | correct a donation amount | 551 | 449 | no measurable change (0.81x, inside the 1.30x error bar) |
| YugabyteDB, 3 nodes (RF=3) | remove one donation | 257 | 280 | no measurable change (1.09x, inside the 1.30x error bar) |
| YugabyteDB, 3 nodes (RF=3) | donor edits their profile | 2.0k | 1.2k | **1.7x slower** |
| YugabyteDB, 3 nodes (RF=3) | erase a donor and all their donations | 121 | 184 | **1.5x faster** |

Read differences below this run's measured 1.30x error bar are omitted. **Write
rows are always shown**, even when the change is inside the error bar — writes are
the price every one of these designs pays for its read speed, and a hidden price
tag reads as no price at all.

### D3 flat+FK → D4 rollup/trg

*What do consolidated aggregates buy and cost?*

| Topology | Operation | D3 flat+FK | D4 rollup/trg | Change |
|---|---|---:|---:|---|
| PostgreSQL, 1 node | insert donation | 6.7k | 2.8k | **2.4x slower** |
| PostgreSQL, 1 node | correct a donation amount | 6.7k | 2.7k | **2.5x slower** |
| PostgreSQL, 1 node | remove one donation | 7.3k | 1.3k | **5.5x slower** |
| PostgreSQL, 1 node | donor edits their profile | 6.7k | 6.1k | no measurable change (0.92x, inside the 1.30x error bar) |
| PostgreSQL, 1 node | erase a donor and all their donations | 5.6k | 389 | **14.4x slower** |
| YugabyteDB, 3 nodes (RF=3) | insert donation | 384 | 103 | **3.7x slower** |
| YugabyteDB, 3 nodes (RF=3) | correct a donation amount | 551 | 178 | **3.1x slower** |
| YugabyteDB, 3 nodes (RF=3) | remove one donation | 257 | 62.1 | **4.1x slower** |
| YugabyteDB, 3 nodes (RF=3) | donor edits their profile | 2.0k | 2.5k | no measurable change (1.24x, inside the 1.30x error bar) |
| YugabyteDB, 3 nodes (RF=3) | erase a donor and all their donations | 121 | 5.4 | **22.5x slower** |

Read differences below this run's measured 1.30x error bar are omitted. **Write
rows are always shown**, even when the change is inside the error bar — writes are
the price every one of these designs pays for its read speed, and a hidden price
tag reads as no price at all.

### D3 flat+FK → D6 embedded

*What does embedding ALL children in the parent buy and cost?*

| Topology | Operation | D3 flat+FK | D6 embedded | Change |
|---|---|---:|---:|---|
| PostgreSQL, 1 node | insert donation | 6.7k | 4.0k | **1.7x slower** |
| PostgreSQL, 1 node | correct a donation amount | 6.7k | 418 | **16.0x slower** |
| PostgreSQL, 1 node | remove one donation | 7.3k | 457 | **16.0x slower** |
| PostgreSQL, 1 node | donor edits their profile | 6.7k | 7.2k | no measurable change (1.08x, inside the 1.30x error bar) |
| PostgreSQL, 1 node | erase a donor and all their donations | 5.6k | 7.4k | **1.3x faster** |
| YugabyteDB, 3 nodes (RF=3) | insert donation | 384 | 302 | no measurable change (0.79x, inside the 1.30x error bar) |
| YugabyteDB, 3 nodes (RF=3) | correct a donation amount | 551 | 88.3 | **6.2x slower** |
| YugabyteDB, 3 nodes (RF=3) | remove one donation | 257 | 103 | **2.5x slower** |
| YugabyteDB, 3 nodes (RF=3) | donor edits their profile | 2.0k | 3.0k | **1.5x faster** |
| YugabyteDB, 3 nodes (RF=3) | erase a donor and all their donations | 121 | 366 | **3.0x faster** |

Read differences below this run's measured 1.30x error bar are omitted. **Write
rows are always shown**, even when the change is inside the error bar — writes are
the price every one of these designs pays for its read speed, and a hidden price
tag reads as no price at all.

### D6 embedded → D9 hybrid

*Does bounding the embedded slice fix what embedding broke?*

| Topology | Operation | D6 embedded | D9 hybrid | Change |
|---|---|---:|---:|---|
| PostgreSQL, 1 node | insert donation | 4.0k | 5.2k | **1.3x faster** |
| PostgreSQL, 1 node | correct a donation amount | 418 | 4.8k | **11.5x faster** |
| PostgreSQL, 1 node | remove one donation | 457 | 5.3k | **11.7x faster** |
| PostgreSQL, 1 node | donor edits their profile | 7.2k | 7.1k | no measurable change (0.99x, inside the 1.30x error bar) |
| PostgreSQL, 1 node | erase a donor and all their donations | 7.4k | 3.8k | **2.0x slower** |
| YugabyteDB, 3 nodes (RF=3) | insert donation | 302 | 247 | no measurable change (0.82x, inside the 1.30x error bar) |
| YugabyteDB, 3 nodes (RF=3) | correct a donation amount | 88.3 | 265 | **3.0x faster** |
| YugabyteDB, 3 nodes (RF=3) | remove one donation | 103 | 204 | **2.0x faster** |
| YugabyteDB, 3 nodes (RF=3) | donor edits their profile | 3.0k | 3.0k | no measurable change (0.99x, inside the 1.30x error bar) |
| YugabyteDB, 3 nodes (RF=3) | erase a donor and all their donations | 366 | 53.4 | **6.9x slower** |

Read differences below this run's measured 1.30x error bar are omitted. **Write
rows are always shown**, even when the change is inside the error bar — writes are
the price every one of these designs pays for its read speed, and a hidden price
tag reads as no price at all.

### D3 flat+FK → D9 hybrid

*What does a bounded embedded cache buy and cost?*

| Topology | Operation | D3 flat+FK | D9 hybrid | Change |
|---|---|---:|---:|---|
| PostgreSQL, 1 node | insert donation | 6.7k | 5.2k | no measurable change (0.78x, inside the 1.30x error bar) |
| PostgreSQL, 1 node | correct a donation amount | 6.7k | 4.8k | **1.4x slower** |
| PostgreSQL, 1 node | remove one donation | 7.3k | 5.3k | **1.4x slower** |
| PostgreSQL, 1 node | donor edits their profile | 6.7k | 7.1k | no measurable change (1.07x, inside the 1.30x error bar) |
| PostgreSQL, 1 node | erase a donor and all their donations | 5.6k | 3.8k | **1.5x slower** |
| YugabyteDB, 3 nodes (RF=3) | insert donation | 384 | 247 | **1.6x slower** |
| YugabyteDB, 3 nodes (RF=3) | correct a donation amount | 551 | 265 | **2.1x slower** |
| YugabyteDB, 3 nodes (RF=3) | remove one donation | 257 | 204 | no measurable change (0.79x, inside the 1.30x error bar) |
| YugabyteDB, 3 nodes (RF=3) | donor edits their profile | 2.0k | 3.0k | **1.5x faster** |
| YugabyteDB, 3 nodes (RF=3) | erase a donor and all their donations | 121 | 53.4 | **2.3x slower** |

Read differences below this run's measured 1.30x error bar are omitted. **Write
rows are always shown**, even when the change is inside the error bar — writes are
the price every one of these designs pays for its read speed, and a hidden price
tag reads as no price at all.


---

Query plans for every cell are in `20260913-writes-isolated/<topology>/plans/<design>.txt`,
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
| Run id | `20260913-writes-isolated` |
| Result files | 10 |
| **Inputs digest** | `de0426905b95e2b0` |

The digest is a content hash of every result file in the run. **Quote it in any
analysis.** A run id alone is not enough to identify what was analysed — a run can
be extended with extra cells afterwards — so the digest is what tells a later
analyst whether they are looking at the same data someone else already wrote about.

### Analyses of this run

| Analyst | Kind | Date | Digest analysed | Status | Headline |
|---|---|---|---|---|---|
| [claude-opus-5](analyses/20260912-study01--claude-opus-5--2026-09-12.md) | ai | 2026-09-12 | `de0426905b95e2b0` ✅ | current | Rollups win aggregate questions by 60-130x but cost 2.4x on inserts and collapse under many writers unless parents are many and small; plain indexes answer everything else; embedding never wins outright on PostgreSQL and a naive cache trigger silently corrupts data. |

### Adding an analysis

Copy `reports/analyses/TEMPLATE.md`, fill in the frontmatter, and write what you
conclude. Then regenerate this report so the table above picks it up.

Before writing one, check the table: if an analysis already exists for digest
`de0426905b95e2b0`, read it first. Add a new analysis to **disagree, extend, or bring a
different perspective** — not to restate what is already there. Never edit another
analyst's file; write your own and reference theirs.
