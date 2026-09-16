# Study 03 — results: reserved seating (venue → event → seat)

Generated from `results/results`. Every number comes from exactly one JSON file in that
directory. **This file contains measurements and no interpretation**; conclusions live in
signed analyses, indexed at the end.

## TL;DR — measured facts

*Selected from the tables below by fixed rules; no interpretation. Read a signed analysis for
what these facts mean.*

- **Cells:** 1 across 1 topologies; 0 failed.
- **Negative controls:** 0 of 0 fired as designed to.
- **Invariant violations in designs meant to be correct:** none observed.
- **Race, YugabyteDB, 1 node (RF=1)** (seats held/s; designs meant to be correct, without violations):

## What was measured

|  |  |
|---|---|
| Run id | `am02-dc02-yb` |
| Environment | [`host-zenbook-ux5406sa`](../../../docs/environments/host-zenbook-ux5406sa.md) |
| Scale | `tiny` — 5 venues, 84 events, 5741 tickets sold, 131 live and 34 expired holds at load |
| Events (kind_seats: count) | catalogue_10: 40 · catalogue_100: 8 · catalogue_1000: 2 · catalogue_10000: 1 · lifecycle_10: 4 · lifecycle_100: 2 · lifecycle_1000: 1 · race_10: 20 · race_100: 4 · race_1000: 1 · race_10000: 1 |
| Seed | 42 — identical data in every cell |
| Reads / isolated writes | 4 workers, 2s measured + 1s warmup, 1 trial(s) |
| Race | 32 buyers per event, real hold TTL 40m0s, 10% of holds deferred until the crowd has gone, timeout 1m0s, tier budget 3m0s, give up after 20 conflicts |

**Engines**

| Topology | Version | READ COMMITTED runs as | Client connections to | DB clock − client clock (ms) |
|---|---|---|---|---:|
| YugabyteDB, 1 node (RF=1) | `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit` | read committed | 1 node(s) | 0.00 |

> Limits of this environment that bound every number below: the client and the databases
> share one laptop's eight cores; the "3-node" cluster has no real network between nodes; the
> race and the lifecycle are closed loops, so tails are a floor; the lifecycle runs on
> compressed time; E0's clock skew is injected.

## Correctness gate

Six read questions checked against Go-computed truth on a load with sold seats, live holds and
expired holds, plus an audit of the loaded state. A failure aborts the cell.

| Design | YugabyteDB, 1 node (RF=1) |
|---|---|
| S1 conditional | ✅ 58/58 |

## Negative controls

Deliberately wrong designs (and E1 with its sweeper stopped). Each must be seen to break the
invariant it exists to test; otherwise the absence of violations elsewhere in that experiment
is not evidence of correctness.

No control was measured in this run.

## Operational reports: answerability (REPORTS.md v2)

The back office's and operations desk's questions, not the buyer's. Added to every
design's `queries.sql` with no change to any existing statement, schema, index or write
path. ✅ answerable, ⚠️ partial (caveat below), ✗ unanswerable. Throughput for whatever
each design can actually answer is in [Reads](#reads) below, under the same report id.

### YugabyteDB, 1 node (RF=1)

| Report | S1 conditional |
|---|---:|
| Holds expiring soon | ✅ |
| Seat status lookup | ✅ |
| Confirmed sales per section in a window | ✅ |
| An event's recent confirmations | ✅ |
| Customers whose last purchase falls in a window | ✅ |
| Hold funnel in a window | ✗ |


- **Hold funnel in a window** (S1 conditional): releasing an expired hold clears the row; no design in this study keeps a hold-event history

## Reads

Each question in isolation, fresh key per execution, ops/s.

### YugabyteDB, 1 node (RF=1)

| Question | S1 conditional |
|---|---:|
| Seat map of a section — 10 seats | 4.1k |
| Seat map of a section — 100 seats | 3.5k |
| Seat map of a section — 1k seats | 2.4k |
| Seat map of a section — 10k seats | 508 |
| Seats available per section — 10 seats | 2.4k |
| Seats available per section — 100 seats | 3.3k |
| Seats available per section — 1k seats | 2.0k |
| Seats available per section — 10k seats | 253 |
| Seats available for an event — 10 seats | 3.9k |
| Seats available for an event — 100 seats | 3.1k |
| Seats available for an event — 1k seats | 2.0k |
| Seats available for an event — 10k seats | 511 |
| A hold's seats (basket page) | 623 |
| A customer's tickets | 2.6k |
| Ticket by id (control) | 4.1k |
| Holds expiring soon | 3.0k |
| Seat status lookup | 4.1k |
| Section sales (historical week) | 1.6k |
| Section sales (trailing week) | 1.8k |
| Recent confirmations | 1.9k |
| Customers' last purchase (historical) | 455 |
| Customers' last purchase (trailing) | 511 |

## Isolated writes, publishing and storage

### YugabyteDB, 1 node (RF=1)

| Design | Hold+confirm /s | Release /s | Refund /s | Publish p50 ms by size | Storage | Audit |
|---|---:|---:|---:|---:|---:|---|
| S1 conditional | — | — | — |  | n/a on YugabyteDB | consistent |


## Controlled pairs — one decision at a time


## Measurement hazard: CPU-quota throttling

Cells show *throttled periods / periods (throttled time)*: for the **database** container(s) over
the whole cell (summed across nodes), and for the **client** per phase.

### YugabyteDB, 1 node (RF=1)

| Design | database, whole cell | client read |
|---|---:|---:|
| S1 conditional | — | — |


---

Plans for every read and write statement are in `results/results/<topology>/plans/<design>.txt`
(`EXPLAIN (ANALYZE, BUFFERS)` on PostgreSQL, `EXPLAIN (ANALYZE, DIST)` on YugabyteDB; write
statements inside rolled-back transactions).

## Provenance: code, data and analyses

This file is **generated from the measurements** and contains no interpretation.
Conclusions live in separate signed analyses, so several analysts — including
different AI models — can read the same numbers and each record what they make of
them. Where two analyses disagree, the disagreement is a finding and stays visible.

| | |
|---|---|
| Run id | `am02-dc02-yb` |
| Result files | 1 |
| **Inputs digest** | `f4f7481c13e744c5` |

Quote the digest in any analysis: it identifies the data. The repository version
identifies the code — check it out (`git checkout <commit>`) to reproduce the run, or
`git diff <commit>` to see what has changed in the study since.

### Analyses of this run

**None yet.** This run has measurements but no interpretation.

### Adding an analysis

Copy `reports/analyses/TEMPLATE.md`, fill in the frontmatter, and regenerate this report.
If an analysis already exists for digest `f4f7481c13e744c5`, read it first and write a new one only to
**disagree, extend, or bring a different perspective**. Never edit another analyst's file.
