# Study 01 — results: tree structures under a charity → person → donation schema

Generated from `20260913-size-small`. Every number here comes from exactly one JSON file in that
directory; the JSON is the source of truth and this file only arranges it.

## What was measured

| | |
|---|---|
| Run id | `20260913-size-small` |
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

## Correctness gate

Every design must return the same answers, checked against values computed
independently in Go from the generated dataset. Timings are only reported for
cells that passed; a failing cell aborts before it is measured.

| Topology | Design | Checks passed |
|---|---|---|
| PostgreSQL, 1 node | D1 minimal | ✅ 12/12 |
| PostgreSQL, 1 node | D2 indexed | ✅ 12/12 |

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
| D1 minimal | 379 | 1.0x | 6.2k | 1.00x | 1 | 12.7 MiB |
| D2 indexed | 4.8k | 12.7x | 6.5k | 1.05x | 3 | 18.5 MiB |


## Read throughput by question (ops/s, higher is better)

### PostgreSQL, 1 node

| Question | D1 minimal | D2 indexed | 
|---|---:|---:|
| Information of the last donation (global) | 132 | 34k | 
| Information of the last donation (one charity) | 212 | 18k | 
| Who donates the most | 193 | 288 | 
| Top-10 donor leaderboard | 204 | 291 | 
| Who donated last | 212 | 19k | 
| First and last donation of a person | 392 | 32k | 
| Total donated (global) | 336 | 352 | 
| Total donated (one charity) | 236 | 356 | 
| A person's 20 most recent donations | 393 | 31k | 
| Donation by id (point lookup) | 39k | 36k | 
| How many donations a person made | 409 | 34k | 
| Charity activity feed (last 50, with names) | 196 | 1.1k | 

## Read tail latency (p99, ms, lower is better)

Means hide the cases a design handles badly; this is where design decisions
usually show up first.

### PostgreSQL, 1 node

| Question | D1 minimal | D2 indexed | 
|---|---:|---:|
| Information of the last donation (global) | 102 | 0.70 | 
| Information of the last donation (one charity) | 98 | 0.84 | 
| Who donates the most | 99 | 99 | 
| Top-10 donor leaderboard | 98 | 99 | 
| Who donated last | 97 | 0.82 | 
| First and last donation of a person | 89 | 0.60 | 
| Total donated (global) | 91 | 89 | 
| Total donated (one charity) | 94 | 96 | 
| A person's 20 most recent donations | 89 | 0.57 | 
| Donation by id (point lookup) | 0.50 | 0.55 | 
| How many donations a person made | 88 | 0.56 | 
| Charity activity feed (last 50, with names) | 98 | 82 | 

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

| Operation | D1 minimal | D2 indexed | 
|---|---:|---:|
| insert donation | 6.2k | 6.5k | 
| correct a donation amount | 6.5k | 6.4k | 
| remove one donation | 6.7k | 7.0k | 
| donor edits their profile | 6.6k | 6.6k | 
| erase a donor and all their donations | 212 | 6.2k | 
| insert p99 (ms) | 2.20 | 1.88 | 
| insert p99.9 (ms) | 3.11 | 5.31 | 
| insert p99.99 (ms) | — | — | 
| insert max (ms) | 21 | 41 | 

## Storage footprint and load time

Read speed bought with redundancy is paid for in bytes and in how long the
data takes to land. Sizes are PostgreSQL only — YugabyteDB stores data in
DocDB rather than PostgreSQL heap files, so `pg_total_relation_size()` does not
describe it.

### PostgreSQL, 1 node

| Design | Tables+indexes | Index bytes | Indexes on `donation` | Load total | COPY | Index build | Rollup |
|---|---:|---:|---:|---:|---:|---:|---:|
| D1 minimal | 12.7 MiB | 3.9 MiB | 1 | 0.4s | 0.3s | 0.0s | 0.0s |
| D2 indexed | 18.5 MiB | 9.5 MiB | 3 | 0.4s | 0.3s | 0.1s | 0.0s |

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

### D1 minimal → D2 indexed

*What do secondary indexes alone buy? (identical SQL)*

| Topology | Operation | D1 minimal | D2 indexed | Change |
|---|---|---:|---:|---|
| PostgreSQL, 1 node | Information of the last donation (global) | 132 | 34k | **254.7x faster** |
| PostgreSQL, 1 node | Information of the last donation (one charity) | 212 | 18k | **84.1x faster** |
| PostgreSQL, 1 node | Who donates the most | 193 | 288 | **1.5x faster** |
| PostgreSQL, 1 node | Top-10 donor leaderboard | 204 | 291 | **1.4x faster** |
| PostgreSQL, 1 node | Who donated last | 212 | 19k | **89.1x faster** |
| PostgreSQL, 1 node | First and last donation of a person | 392 | 32k | **81.2x faster** |
| PostgreSQL, 1 node | Total donated (one charity) | 236 | 356 | **1.5x faster** |
| PostgreSQL, 1 node | A person's 20 most recent donations | 393 | 31k | **79.5x faster** |
| PostgreSQL, 1 node | How many donations a person made | 409 | 34k | **83.2x faster** |
| PostgreSQL, 1 node | Charity activity feed (last 50, with names) | 196 | 1.1k | **5.8x faster** |
| PostgreSQL, 1 node | insert donation | 6.2k | 6.5k | no measurable change (1.05x, inside the 1.30x error bar) |
| PostgreSQL, 1 node | correct a donation amount | 6.5k | 6.4k | no measurable change (0.98x, inside the 1.30x error bar) |
| PostgreSQL, 1 node | remove one donation | 6.7k | 7.0k | no measurable change (1.05x, inside the 1.30x error bar) |
| PostgreSQL, 1 node | donor edits their profile | 6.6k | 6.6k | no measurable change (0.99x, inside the 1.30x error bar) |
| PostgreSQL, 1 node | erase a donor and all their donations | 212 | 6.2k | **29.3x faster** |

Read differences below this run's measured 1.30x error bar are omitted. **Write
rows are always shown**, even when the change is inside the error bar — writes are
the price every one of these designs pays for its read speed, and a hidden price
tag reads as no price at all.


---

Query plans for every cell are in `20260913-size-small/<topology>/plans/<design>.txt`,
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
| Run id | `20260913-size-small` |
| Result files | 2 |
| **Inputs digest** | `d556eb4496efc31a` |

The digest is a content hash of every result file in the run. **Quote it in any
analysis.** A run id alone is not enough to identify what was analysed — a run can
be extended with extra cells afterwards — so the digest is what tells a later
analyst whether they are looking at the same data someone else already wrote about.

### Analyses of this run

| Analyst | Kind | Date | Digest analysed | Status | Headline |
|---|---|---|---|---|---|
| [claude-opus-5](analyses/20260912-study01--claude-opus-5--2026-09-12.md) | ai | 2026-09-12 | `d556eb4496efc31a` ✅ | current | Store totals on the parent only for aggregate questions and only when writes spread across many parents; use indexes for everything else; embedding donations in the donor row never came out ahead on reads. |
| [gpt-6](analyses/20260912-study01--gpt-6--2026-09-12.md) | ai | 2026-09-12 | `d556eb4496efc31a` ✅ | current | Use ordinary indexes for donor reads, copy the charity key when it enables selective charity indexes, and store totals only when their read benefit justifies contention and maintenance. |

### Adding an analysis

Copy `reports/analyses/TEMPLATE.md`, fill in the frontmatter, and write what you
conclude. Then regenerate this report so the table above picks it up.

Before writing one, check the table: if an analysis already exists for digest
`d556eb4496efc31a`, read it first. Add a new analysis to **disagree, extend, or bring a
different perspective** — not to restate what is already there. Never edit another
analyst's file; write your own and reference theirs.
