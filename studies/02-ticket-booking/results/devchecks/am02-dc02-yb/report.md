# Study 02 — results: avoiding overbooking (band → event → ticket)

Generated from `results/results`. Every number comes from exactly one JSON file in that
directory. **This file contains measurements and no interpretation**; conclusions live in
signed analyses, indexed at the end.

## TL;DR — measured facts

*Selected from the tables below by fixed rules; no interpretation. Read a signed analysis for
what these facts mean.*

- **Cells:** 1 across 1 topology; 0 failed.
- **Negative controls:** 0 of 0 fired (overbooked, as designed to).
- **Overbooking in designs meant to be correct:** none observed.
- **Under-booking** (sold out with seats unsold): none observed.

## What was measured

|  |  |
|---|---|
| Run id | `am02-dc02-yb` |
| Environment | [`host-zenbook-ux5406sa`](../../../docs/environments/host-zenbook-ux5406sa.md) |
| Scale | `tiny` — 10 bands, 90 events, 26100 seats, 6880 tickets sold at load |
| Events (count × seats) | **catalogue**: 40 × 10 seats, 8 × 100 seats, 2 × 1k seats, 1 × 10k seats · **race**: 20 × 10 seats, 4 × 100 seats, 1 × 1k seats, 1 × 10k seats · **churn**: 10 × 10 seats, 2 × 100 seats, 1 × 1k seats |
| Seed | 42 — identical data in every cell |
| Reads / isolated writes | 4 workers, 2s measured + 1s warmup, 1 trial(s) |
| Sell-out race | 32 buyers per event, timeout 1m0s, organiser edit every 20ms; churn race cancels 10% of sales |

**Engines**

| Topology | Version | READ COMMITTED requests actually run as |
|---|---|---|
| YugabyteDB, 1 node (RF=1) | `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit` | read committed |

> Limits of this environment that bound every number below: the client and the
> databases share one laptop's eight cores, and the "3-node" cluster has no real network
> between nodes. The race is a closed loop (buyers wait for their answer), so its tails
> are a floor, not an SLO figure.

## Correctness gate

Every design must answer the read questions exactly as computed from the generated data and
pass the overbooking audit on the freshly loaded tables before anything is timed. The audit
then runs again after every phase that writes, on the state that phase left behind: events
over capacity, seats sold twice, counters/buckets/pools that disagree with the tickets, and a
reconciliation of the tickets that exist against the sales the harness saw commit.

| Topology | Design | Answers | Audit on load | After isolated writes | After races | After holds | Cell |
|---|---|---|---|---|---|---|---|
| YugabyteDB, 1 node (RF=1) | P3 CAS | ✅ 44/44 | ✅ | — | — | — | ✅ |

### Negative controls

Two designs are in the study **because they are wrong**. An audit that has never caught a
wrong design has not been shown to work, so each control must be seen to fail. If a control
did not fail, the absence of violations in the other designs of that experiment is not
evidence of their correctness — only of too little contention.

| Topology | Control | Expected failure | Observed |
|---|---|---|---|


## Pre-created or created: the cost outside the race

*Publish* creates a new event with its inventory in one transaction — one row, or one row
per seat — and is shown as the median time to publish one event. *Book (spread)* is demand
spread across the catalogue in proportion to seats left: little contention on any one event.
*Cancel* refunds a sold ticket. Each operation ran on its own fresh load.

### YugabyteDB, 1 node (RF=1)

| Design | publish 10 seats (ms) | publish 100 seats (ms) | publish 1k seats (ms) | publish 10k seats (ms) | publish 100k seats (ms) | book spread /s | cancel /s | load total | storage |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| P3 CAS | — | — | — | — | — | — | — | 6.1s | — |

Storage is PostgreSQL only: YugabyteDB keeps data in DocDB, where `pg_total_relation_size()`
does not describe it.


## Reads (ops/s, higher is better)

Each question benchmarked in isolation with a fresh key per execution. Availability is
measured separately for every event size, because that is where counting designs pay.

### YugabyteDB, 1 node (RF=1)

| Question | P3 CAS |
|---|---:|
| Seats left — 10 seats | 4.5k |
| Seats left — 100 seats | 4.3k |
| Seats left — 1k seats | 2.9k |
| Seats left — 10k seats | 1.7k |
| A band's events with availability | 603 |
| A customer's tickets | 2.7k |
| Ticket by id (control) | 4.7k |
| Tickets sold across a band | 208 |


## Operational reports (ops/s, higher is better)

The back office's questions (REPORTS.md), not the buyer's. Added to every design's
`queries.sql` with no change to any existing statement, schema, index or write path.
r03 and r05 are measured in both window regimes (`@trailing`, `@historical`) for the same
reason RECENCY.md section 2 gives: in the trailing regime the cheap wrong answer happens
to be right, and only the historical regime would catch a design that got it wrong.

### YugabyteDB, 1 node (RF=1)

**Answerability** — ✅ answerable, ⚠️ partial (caveat below), ✗ unanswerable:

| Report | P3 CAS |
|---|---:|
| Sales and revenue in a window | ⚠️ |
| An event's recent buyers | ✅ |
| Customers whose last purchase falls in a window | ✅ |
| Sell-through per event, by band | ✅ |
| Refunds in a window | ✗ |
| Outstanding holds | ✗ |


- **Sales and revenue in a window** (P3 CAS): a ticket sold in the window and later refunded is not counted, because cancellation erases the sale
- **Refunds in a window** (P3 CAS): cancellation deletes or resets the sale; nothing records that it happened
- **Outstanding holds** (P3 CAS): no reservation table in this design

| Report | P3 CAS |
|---|---:|
| Sales and revenue in a window (trailing) | 1.2k |
| Sales and revenue in a window (historical) | 1.2k |
| An event's recent buyers | 1.2k |
| Customers whose last purchase falls in a window (trailing) | 174 |
| Customers whose last purchase falls in a window (historical) | 176 |
| Sell-through per event, by band | 195 |



## Controlled pairs — one decision at a time

The C1/C2/C3 read control is not present in this run, so read differences are filtered at a
conventional 1.30x rather than a measured error bar.


## Measurement hazard: CPU-quota throttling

Every container runs under a CFS quota (`--cpus`). When a container exhausts a 100 ms period's
allowance it is paused until the next period, which puts tens of milliseconds into whatever
request it was serving or issuing — latency that belongs to no design. Cells show *throttled
periods / periods (throttled time)*: for the **database** container(s) over the whole cell
(summed across nodes), and for the **client** per phase.

### YugabyteDB, 1 node (RF=1)

| Design | database, whole cell | client read | client reports |
|---|---:|---:|---:|
| P3 CAS | — | — | — |


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
| **Inputs digest** | `87de4c86a8a9b82e` |

Quote the digest in any analysis: it identifies the data. The repository version
identifies the code — check it out (`git checkout <commit>`) to reproduce the run, or
`git diff <commit>` to see what has changed in the study since.

### Analyses of this run

**None yet.** This run has measurements but no interpretation.

### Adding an analysis

Copy `reports/analyses/TEMPLATE.md`, fill in the frontmatter, and regenerate this report.
If an analysis already exists for digest `87de4c86a8a9b82e`, read it first and write a new one only to
**disagree, extend, or bring a different perspective**. Never edit another analyst's file.
