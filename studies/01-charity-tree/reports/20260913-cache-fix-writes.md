# Study 01 — results: tree structures under a charity → person → donation schema

Generated from `20260913-cache-fix-writes`. Every number here comes from exactly one JSON file in that
directory; the JSON is the source of truth and this file only arranges it.

## What was measured

| | |
|---|---|
| Run id | `20260913-cache-fix-writes` |
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

## Correctness gate

Every design must return the same answers, checked against values computed
independently in Go from the generated dataset. Timings are only reported for
cells that passed; a failing cell aborts before it is measured.

| Topology | Design | Checks passed |
|---|---|---|
| PostgreSQL, 1 node | D9 hybrid | ✅ 12/12 |
| PostgreSQL, 1 node | D10 hybrid/locked | ✅ 12/12 |

All cells passed.

## The trade-off at a glance

Every design here buys read speed with something — write throughput, bytes on
disk, or index maintenance on the hot path. This table puts the purchase and the
price in the same row.

*Read score* is the **geometric** mean across all twelve questions. These
throughputs span four orders of magnitude, so an arithmetic mean would be decided
entirely by the fastest query and say nothing about the other eleven.

### PostgreSQL, 1 node

| Design | Read score | vs D9 hybrid | Insert/s | vs D9 hybrid | Indexes on `donation` | Storage |
|---|---:|---:|---:|---:|---:|---:|
| D9 hybrid | — | — | 4.8k | 1.00x | 5 | 29.9 MiB |
| D10 hybrid/locked | — | — | 4.6k | 0.96x | 5 | 29.3 MiB |


## Read throughput by question (ops/s, higher is better)

### PostgreSQL, 1 node

| Question | D9 hybrid | D10 hybrid/locked | 
|---|---:|---:|

## Read tail latency (p99, ms, lower is better)

Means hide the cases a design handles badly; this is where design decisions
usually show up first.

### PostgreSQL, 1 node

| Question | D9 hybrid | D10 hybrid/locked | 
|---|---:|---:|

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

| Operation | D9 hybrid | D10 hybrid/locked | 
|---|---:|---:|
| insert donation | 4.8k | 4.6k | 
| correct a donation amount | 4.5k | 4.3k | 
| remove one donation | 4.9k | 4.8k | 
| donor edits their profile | 6.9k | 6.9k | 
| erase a donor and all their donations | 3.4k | 3.0k | 
| insert p99 (ms) | 14 | 16 | 
| insert p99.9 (ms) | 19 | 22 | 
| insert p99.99 (ms) | — | — | 
| insert max (ms) | 43 | 44 | 

## Storage footprint and load time

Read speed bought with redundancy is paid for in bytes and in how long the
data takes to land. Sizes are PostgreSQL only — YugabyteDB stores data in
DocDB rather than PostgreSQL heap files, so `pg_total_relation_size()` does not
describe it.

### PostgreSQL, 1 node

| Design | Tables+indexes | Index bytes | Indexes on `donation` | Load total | COPY | Index build | Rollup |
|---|---:|---:|---:|---:|---:|---:|---:|
| D9 hybrid | 29.9 MiB | 16.1 MiB | 5 | 0.7s | 0.5s | 0.1s | 0.0s |
| D10 hybrid/locked | 29.3 MiB | 16.1 MiB | 5 | 0.7s | 0.5s | 0.1s | 0.0s |

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

### D9 hybrid → D10 hybrid/locked

*What does making the cache trigger concurrency-correct cost?*

| Topology | Operation | D9 hybrid | D10 hybrid/locked | Change |
|---|---|---:|---:|---|
| PostgreSQL, 1 node | insert donation | 4.8k | 4.6k | no measurable change (0.96x, inside the 1.30x error bar) |
| PostgreSQL, 1 node | correct a donation amount | 4.5k | 4.3k | no measurable change (0.97x, inside the 1.30x error bar) |
| PostgreSQL, 1 node | remove one donation | 4.9k | 4.8k | no measurable change (0.98x, inside the 1.30x error bar) |
| PostgreSQL, 1 node | donor edits their profile | 6.9k | 6.9k | no measurable change (1.01x, inside the 1.30x error bar) |
| PostgreSQL, 1 node | erase a donor and all their donations | 3.4k | 3.0k | no measurable change (0.89x, inside the 1.30x error bar) |

Read differences below this run's measured 1.30x error bar are omitted. **Write
rows are always shown**, even when the change is inside the error bar — writes are
the price every one of these designs pays for its read speed, and a hidden price
tag reads as no price at all.


---

Query plans for every cell are in `20260913-cache-fix-writes/<topology>/plans/<design>.txt`,
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
| Run id | `20260913-cache-fix-writes` |
| Result files | 2 |
| **Inputs digest** | `3e7824b31a92942f` |

The digest is a content hash of every result file in the run. **Quote it in any
analysis.** A run id alone is not enough to identify what was analysed — a run can
be extended with extra cells afterwards — so the digest is what tells a later
analyst whether they are looking at the same data someone else already wrote about.

### Analyses of this run

| Analyst | Kind | Date | Digest analysed | Status | Headline |
|---|---|---|---|---|---|
| [claude-opus-5](analyses/20260912-study01--claude-opus-5--2026-09-12.md) | ai | 2026-09-12 | `3e7824b31a92942f` ✅ | current | Rollups win aggregate questions by 60-130x but cost 2.4x on inserts and collapse under many writers unless parents are many and small; plain indexes answer everything else; embedding never wins outright on PostgreSQL and a naive cache trigger silently corrupts data. |

### Adding an analysis

Copy `reports/analyses/TEMPLATE.md`, fill in the frontmatter, and write what you
conclude. Then regenerate this report so the table above picks it up.

Before writing one, check the table: if an analysis already exists for digest
`3e7824b31a92942f`, read it first. Add a new analysis to **disagree, extend, or bring a
different perspective** — not to restate what is already there. Never edit another
analyst's file; write your own and reference theirs.
