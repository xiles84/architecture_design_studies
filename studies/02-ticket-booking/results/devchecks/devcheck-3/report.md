# Study 02 — results: avoiding overbooking (band → event → ticket)

Generated from `results/devcheck-3`. Every number comes from exactly one JSON file in that
directory. **This file contains measurements and no interpretation**; conclusions live in
signed analyses, indexed at the end.

## What was measured

|  |  |
|---|---|
| Run id | `devcheck-3` |
| Environment | [`host-zenbook-ux5406sa`](../../../docs/environments/host-zenbook-ux5406sa.md) |
| Scale | `tiny` — 10 bands, 90 events, 26100 seats, 6880 tickets sold at load |
| Events (count × seats) | **catalogue**: 40 × 10 seats, 8 × 100 seats, 2 × 1k seats, 1 × 10k seats · **race**: 20 × 10 seats, 4 × 100 seats, 1 × 1k seats, 1 × 10k seats · **churn**: 10 × 10 seats, 2 × 100 seats, 1 × 1k seats |
| Seed | 42 — identical data in every cell |
| Reads / isolated writes | 8 workers, 1s measured + 300ms warmup, 1 trial(s) |
| Sell-out race | 32 buyers per event, timeout 30s, organiser edit every 20ms; churn race cancels 10% of sales |
| Holds | 32 buyers, TTL 250ms, basket ≤ 200ms, payment 20ms-100ms, 25% abandoned, sweeper every 50ms |

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
| YugabyteDB, 1 node (RF=1) | C2 count (SER) | ✅ 38/38 | ✅ | ✅ | — | — | ❌ race race event 52: final count: failed to connect to `user=yugabyte database=yugabyte`: |
| YugabyteDB, 1 node (RF=1) | C5 seat-unique | ✅ 38/38 | ✅ | ✅ | ❌ 5 drift | — | ✅ |
| YugabyteDB, 1 node (RF=1) | R1 inventory | ✅ 38/38 | ✅ | ✅ | ✅ | — | ✅ |
| YugabyteDB, 1 node (RF=1) | R2 buckets | ✅ 38/38 | ✅ | ✅ | ✅ | — | ✅ |
| YugabyteDB, 1 node (RF=1) | R3 seat-pool | ✅ 38/38 | ✅ | ✅ | ✅ | — | ✅ |
| YugabyteDB, 1 node (RF=1) | H0 hold/naive ✗ | ✅ 38/38 | ✅ | ✅ | ✅ | ❌ 2 overbooked, 3 drift | ✅ |
| YugabyteDB, 1 node (RF=1) | H1 hold/checked | ✅ 38/38 | ✅ | ✅ | ✅ | ✅ | ✅ |

### Negative controls

Two designs are in the study **because they are wrong**. An audit that has never caught a
wrong design has not been shown to work, so each control must be seen to fail. If a control
did not fail, the absence of violations in the other designs of that experiment is not
evidence of their correctness — only of too little contention.

| Topology | Control | Expected failure | Observed |
|---|---|---|---|
| YugabyteDB, 1 node (RF=1) | H0 hold/naive ✗ | confirms a hold that expired during payment, after the seat may have been resold | ✅ fired: 2 events overbooked by 2 seats in total |


## The sell-out race — seats sold per second (higher is better)

Every buyer arrives at an unsold event at the same instant and keeps buying until told
"sold out". Throughput is successful sales over the time from the start to the last sale
(summed over the tier's events); the time spent turning the rest of the crowd away is
reported separately below as *sold out answered in*. Markers:
❌ overbooked (events/total, extra seats) and struck through — **not a valid speed**;
⚠️ underbooked (every buyer was sent away with seats unsold); ⏱ race hit its timeout.

### YugabyteDB, 1 node (RF=1)

| Design | 10 seats | 100 seats |
|---|---:|---:|
| C2 count (SER) | — | — |
| C5 seat-unique | 24.3 | 24.2 |
| R1 inventory | 20.0 | 18.4 |
| R2 buckets | 25.0 | 67.0 |
| R3 seat-pool | 32.3 | 87.0 |
| H0 hold/naive ✗ | 16.7 | 15.7 |
| H1 hold/checked | 16.8 | 15.5 |


## The race in detail

How each sale was achieved, not just how many. *Attempts per sale* is (sales + retries) ÷
sales: 1.00 means no buyer ever lost a race and started again. *Sold out answered in* is the
median latency of the final "no" each buyer received. *Organiser edit p99* is the latency of an
unrelated update to the event page made while the crowd was buying.

### YugabyteDB, 1 node (RF=1)

| Sale latency p99 (ms) | 10 seats | 100 seats |
|---|---:|---:|
| C2 count (SER) | — | — |
| C5 seat-unique | 550 | 3885 |
| R1 inventory | 595 | 2817 |
| R2 buckets | 869 | 1088 |
| R3 seat-pool | 624 | 768 |
| H0 hold/naive ✗ | 1286 | 2828 |
| H1 hold/checked | 815 | 3154 |

| Attempts per sale | 10 seats | 100 seats |
|---|---:|---:|
| C2 count (SER) | — | — |
| C5 seat-unique | 20.82 | 20.18 |
| R1 inventory | 1.00 | 1.00 |
| R2 buckets | 10.30 | 3.11 |
| R3 seat-pool | 7.08 | 1.01 |
| H0 hold/naive ✗ | 1.00 | 1.00 |
| H1 hold/checked | 1.00 | 1.00 |

| Sold out answered in, p50 (ms) | 10 seats | 100 seats |
|---|---:|---:|
| C2 count (SER) | — | — |
| C5 seat-unique | 400 | 911 |
| R1 inventory | 486 | 1067 |
| R2 buckets | 198 | 410 |
| R3 seat-pool | 178 | 206 |
| H0 hold/naive ✗ | 519 | 1102 |
| H1 hold/checked | 512 | 1101 |

| Organiser edit p99 (ms) | 10 seats | 100 seats |
|---|---:|---:|
| C2 count (SER) | — | — |
| C5 seat-unique | 175 | 102 |
| R1 inventory | 77 | 68 |
| R2 buckets | 195 | 94 |
| R3 seat-pool | 273 | 92 |
| H0 hold/naive ✗ | 633 | 1615 |
| H1 hold/checked | 563 | 1162 |


## Churn race — selling out while buyers cancel

As the race, but a share of successful buyers cancel immediately. After the crowd leaves,
one last buyer sweeps up refunded seats; seats still unsold after that sweep are
**underbooked** — the design said "sold out" while it had seats.

### YugabyteDB, 1 node (RF=1)

| Design | 10 seats | 100 seats |
|---|---:|---:|
| C2 count (SER) | — | — |
| C5 seat-unique | 24.2 ⚠️ underbooked 3/10 (−3), ⚠️ drift 3 | 23.0 ⚠️ underbooked 2/2 (−10), ⚠️ drift 5 |
| R1 inventory | 17.5 | 19.7 |
| R2 buckets | 18.0 | 59.1 |
| R3 seat-pool | 23.7 | 94.6 |
| H0 hold/naive ✗ | 20.6 | 14.9 |
| H1 hold/checked | 16.7 | 15.4 |

| Design | 10 seats — cancels / swept up / underbooked seats | 100 seats |
|---|---:|---:|
| C2 count (SER) | — | — |
| C5 seat-unique | 9 / 0 / 3 | 22 / 0 / 10 |
| R1 inventory | 11 / 0 / 0 | 18 / 0 / 0 |
| R2 buckets | 12 / 0 / 0 | 27 / 0 / 0 |
| R3 seat-pool | 10 / 0 / 0 | 21 / 0 / 0 |
| H0 hold/naive ✗ | 14 / 0 / 0 | 11 / 0 / 0 |
| H1 hold/checked | 10 / 0 / 0 | 17 / 0 / 0 |


## Pre-created or created: the cost outside the race

*Publish* creates a new event with its inventory in one transaction — one row, or one row
per seat — and is shown as the median time to publish one event. *Book (spread)* is demand
spread across the catalogue in proportion to seats left: little contention on any one event.
*Cancel* refunds a sold ticket. Each operation ran on its own fresh load.

### YugabyteDB, 1 node (RF=1)

| Design | publish 10 seats (ms) | publish 100 seats (ms) | publish 1k seats (ms) | publish 10k seats (ms) | publish 100k seats (ms) | book spread /s | cancel /s | load total | storage |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| C2 count (SER) | 2.94 | 2.86 | 2.85 | 2.34 | 1.95 | 4.6 | 530 | 4.6s | — |
| C5 seat-unique | 2.46 | 3.05 | 2.51 | 2.15 | 1.97 | 104 | 590 | 23.9s | — |
| R1 inventory | 4.02 | 4.66 | 5.32 | 4.06 | 3.72 | 79.7 | 53.8 | 4.7s | — |
| R2 buckets | 4.64 | 5.38 | 5.97 | 51 | 3.87 | 190 | 185 | 5.6s | — |
| R3 seat-pool | 4.10 | 5.83 | 16 | 95 | 1917 | 207 | 358 | 6.3s | — |
| H0 hold/naive ✗ | 2.29 | 2.45 | 2.13 | 1.84 | 1.73 | 96.7 | 48.3 | 4.8s | — |
| H1 hold/checked | 2.34 | 2.34 | 2.10 | 1.56 | 1.76 | 86.5 | 47.6 | 5.0s | — |

Storage is PostgreSQL only: YugabyteDB keeps data in DocDB, where `pg_total_relation_size()`
does not describe it.


## Reads (ops/s, higher is better)

Each question benchmarked in isolation with a fresh key per execution. Availability is
measured separately for every event size, because that is where counting designs pay.

### YugabyteDB, 1 node (RF=1)

| Question | C2 count (SER) | C5 seat-unique | R1 inventory | R2 buckets | R3 seat-pool | H0 hold/naive ✗ | H1 hold/checked |
|---|---:|---:|---:|---:|---:|---:|---:|
| A band's events with availability | — | — | — | — | — | — | — |
| A customer's tickets | — | — | — | — | — | — | — |
| Ticket by id (control) | — | — | — | — | — | — | — |
| Tickets sold across a band | — | — | — | — | — | — | — |


## Holds with expiry (H0, H1)

Buyers hold a seat, spend time in the basket, abandon some holds, and pay for the rest; a
sweeper releases expired holds. Some payments finish after the hold expired. *Late confirms
rejected* are those the design refused (and would refund); a design that accepts them can sell a
released seat twice.

### YugabyteDB, 1 node (RF=1)

| Design | Event size | Events | Confirmed/s | Holds | Abandoned | Expired at pre-check | Late confirms rejected | Swept | Hold p99 (ms) | Confirm p99 (ms) | Overbooked | Errors |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---:|
| H0 hold/naive ✗ | 10 seats | 4 | 6.8 | 84 | 19 | 23 | 0 | 44 | 559 | 117 | ❌ 2 events (+2 seats) | 0 |
| H0 hold/naive ✗ | 100 seats | 4 | 0.7 | 2478 | 600 | 1787 | 0 | 2388 | 2439 | 105 | ✅ none, ⏱ 4 timed out | 0 |
| H1 hold/checked | 10 seats | 4 | 4.5 | 147 | 46 | 44 | 17 | 107 | 576 | 106 | ✅ none | 0 |
| H1 hold/checked | 100 seats | 4 | 0.4 | 2485 | 629 | 1773 | 35 | 2437 | 2481 | 82 | ✅ none, ⏱ 4 timed out | 0 |


## Controlled pairs — one decision at a time

The C1/C2/C3 read control is not present in this run, so read differences are filtered at a
conventional 1.30x rather than a measured error bar.

### R1 inventory → R2 buckets

*What does sharding the counter into buckets buy?*

| Topology | Measurement | R1 inventory | R2 buckets | Change |
|---|---|---:|---:|---|
| YugabyteDB, 1 node (RF=1) | race, 10 seats — sold/s | 20.0 | 25.0 | 1.2x faster |
| YugabyteDB, 1 node (RF=1) | race, 100 seats — sold/s | 18.4 | 67.0 | 3.6x faster |
| YugabyteDB, 1 node (RF=1) | churn, 10 seats — sold/s | 17.5 | 18.0 | 1.0x faster |
| YugabyteDB, 1 node (RF=1) | churn, 100 seats — sold/s | 19.7 | 59.1 | 3.0x faster |
| YugabyteDB, 1 node (RF=1) | book /s | 79.7 | 190 | 2.4x faster |
| YugabyteDB, 1 node (RF=1) | cancel /s | 53.8 | 185 | 3.4x faster |
| YugabyteDB, 1 node (RF=1) | publish 10 seats p50 ms | 4.02 | 4.64 | 1.2x slower |
| YugabyteDB, 1 node (RF=1) | publish 100 seats p50 ms | 4.66 | 5.38 | 1.2x slower |
| YugabyteDB, 1 node (RF=1) | publish 1k seats p50 ms | 5.32 | 5.97 | 1.1x slower |
| YugabyteDB, 1 node (RF=1) | publish 10k seats p50 ms | 4.06 | 51 | 12.6x slower |
| YugabyteDB, 1 node (RF=1) | publish 100k seats p50 ms | 3.72 | 3.87 | 1.0x slower |

### H0 hold/naive ✗ → H1 hold/checked

*What does validating the hold at confirmation cost?*

| Topology | Measurement | H0 hold/naive ✗ | H1 hold/checked | Change |
|---|---|---:|---:|---|
| YugabyteDB, 1 node (RF=1) | race, 10 seats — sold/s | 16.7 | 16.8 | 1.0x faster |
| YugabyteDB, 1 node (RF=1) | race, 100 seats — sold/s | 15.7 | 15.5 | 1.0x slower |
| YugabyteDB, 1 node (RF=1) | churn, 10 seats — sold/s | 20.6 | 16.7 | 1.2x slower |
| YugabyteDB, 1 node (RF=1) | churn, 100 seats — sold/s | 14.9 | 15.4 | 1.0x faster |
| YugabyteDB, 1 node (RF=1) | book /s | 96.7 | 86.5 | 1.1x slower |
| YugabyteDB, 1 node (RF=1) | cancel /s | 48.3 | 47.6 | 1.0x slower |
| YugabyteDB, 1 node (RF=1) | publish 10 seats p50 ms | 2.29 | 2.34 | 1.0x slower |
| YugabyteDB, 1 node (RF=1) | publish 100 seats p50 ms | 2.45 | 2.34 | 1.0x faster |
| YugabyteDB, 1 node (RF=1) | publish 1k seats p50 ms | 2.13 | 2.10 | 1.0x faster |
| YugabyteDB, 1 node (RF=1) | publish 10k seats p50 ms | 1.84 | 1.56 | 1.2x faster |
| YugabyteDB, 1 node (RF=1) | publish 100k seats p50 ms | 1.73 | 1.76 | 1.0x slower |
| YugabyteDB, 1 node (RF=1) | holds, 10 seats — confirmed/s (overbooked events) | 6.8 (2) | 4.5 (0) | 1.5x slower |
| YugabyteDB, 1 node (RF=1) | holds, 100 seats — confirmed/s (overbooked events) | 0.7 (0) | 0.4 (0) | 1.9x slower |


## Measurement hazard: client CPU throttling

The benchmark client runs under a CFS quota (`--cpus`). When it exhausts a 100 ms period's
allowance it is paused until the next period, which puts tens of milliseconds into whatever
request it was issuing. Cells show *throttled periods / periods (throttled time)* for the
client container during each phase. Database-side counters are in `logs/<design>.dbcpu.txt`.

### YugabyteDB, 1 node (RF=1)

| Design | churn#1 | holds | race#1 | write:book | write:cancel | write:publish |
|---|---:|---:|---:|---:|---:|---:|
| C2 count (SER) | — | — | 0/82 (0.0s) | 0/43 (0.0s) | 0/14 (0.0s) | 0/49 (0.0s) |
| C5 seat-unique | 0/144 (0.0s) | — | 0/273 (0.0s) | 0/16 (0.0s) | 0/13 (0.0s) | 0/41 (0.0s) |
| R1 inventory | 0/181 (0.0s) | — | 0/341 (0.0s) | 0/15 (0.0s) | 0/14 (0.0s) | 0/63 (0.0s) |
| R2 buckets | 0/109 (0.0s) | — | 0/165 (0.0s) | 0/16 (0.0s) | 0/14 (0.0s) | 0/76 (0.0s) |
| R3 seat-pool | 0/77 (0.0s) | — | 0/135 (0.0s) | 0/16 (0.0s) | 0/14 (0.0s) | 0/78 (0.0s) |
| H0 hold/naive ✗ | 0/200 (0.0s) | 0/1306 (0.0s) | 0/387 (0.0s) | 0/16 (0.0s) | 0/15 (0.0s) | 0/38 (0.0s) |
| H1 hold/checked | 0/207 (0.0s) | 0/1333 (0.0s) | 0/390 (0.0s) | 0/17 (0.0s) | 0/15 (0.0s) | 0/39 (0.0s) |


---

Plans for every read and write statement are in `results/devcheck-3/<topology>/plans/<design>.txt`
(`EXPLAIN (ANALYZE, BUFFERS)` on PostgreSQL, `EXPLAIN (ANALYZE, DIST)` on YugabyteDB; write
statements inside rolled-back transactions).

## Provenance: code, data and analyses

This file is **generated from the measurements** and contains no interpretation.
Conclusions live in separate signed analyses, so several analysts — including
different AI models — can read the same numbers and each record what they make of
them. Where two analyses disagree, the disagreement is a finding and stays visible.

| | |
|---|---|
| Run id | `devcheck-3` |
| Result files | 7 |
| **Inputs digest** | `4115c06111f49340` |
| Repository version | `unknown` (`unknown`), clean |

Quote the digest in any analysis: it identifies the data. The repository version
identifies the code — check it out (`git checkout <commit>`) to reproduce the run, or
`git diff <commit>` to see what has changed in the study since.

### Analyses of this run

**None yet.** This run has measurements but no interpretation.

### Adding an analysis

Copy `reports/analyses/TEMPLATE.md`, fill in the frontmatter, and regenerate this report.
If an analysis already exists for digest `4115c06111f49340`, read it first and write a new one only to
**disagree, extend, or bring a different perspective**. Never edit another analyst's file.
