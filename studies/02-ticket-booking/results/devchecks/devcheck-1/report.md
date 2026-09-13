# Study 02 — results: avoiding overbooking (band → event → ticket)

Generated from `results/devcheck-1`. Every number comes from exactly one JSON file in that
directory. **This file contains measurements and no interpretation**; conclusions live in
signed analyses, indexed at the end.

## What was measured

|  |  |
|---|---|
| Run id | `devcheck-1` |
| Environment | [`host-zenbook-ux5406sa`](../../../docs/environments/host-zenbook-ux5406sa.md) |
| Scale | `tiny` — 10 bands, 90 events, 26100 seats, 6880 tickets sold at load |
| Events by kind and size | `map[catalogue_10:40 catalogue_100:8 catalogue_1000:2 catalogue_10000:1 churn_10:10 churn_100:2 churn_1000:1 race_10:20 race_100:4 race_1000:1 race_10000:1]` |
| Seed | 42 — identical data in every cell |
| Reads / isolated writes | 8 workers, 1s measured + 300ms warmup, 1 trial(s) |
| Sell-out race | 32 buyers per event, timeout 1m0s, organiser edit every 20ms; churn race cancels 10% of sales |
| Holds | 32 buyers, TTL 250ms, basket ≤ 200ms, payment 20ms-100ms, 25% abandoned, sweeper every 50ms |

**Engines**

| Topology | Version | READ COMMITTED requests actually run as |
|---|---|---|
| PostgreSQL, 1 node | `PostgreSQL 17.11 (Debian 17.11-1.pgdg13+2) on x86_64-pc-linux-gnu, compiled by gcc (Debian 14.2.0-19) 14.2.0, 64-bit` | read committed |

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
| PostgreSQL, 1 node | P2 skip-locked | ✅ 38/38 | ✅ | ✅ | ✅ | — | ✅ |
| PostgreSQL, 1 node | C1 count (RC) ✗ | ✅ 38/38 | ✅ | ✅ | ❌ 26 overbooked | — | ✅ |
| PostgreSQL, 1 node | H1 hold/checked | ✅ 38/38 | ✅ | ✅ | ✅ | ✅ | ✅ |

### Negative controls

Two designs are in the study **because they are wrong**. An audit that has never caught a
wrong design has not been shown to work, so each control must be seen to fail. If a control
did not fail, the absence of violations in the other designs of that experiment is not
evidence of their correctness — only of too little contention.

| Topology | Control | Expected failure | Observed |
|---|---|---|---|
| PostgreSQL, 1 node | C1 count (RC) ✗ | check-then-insert at READ COMMITTED: two buyers can both see the last seat free | ✅ fired: 39 events overbooked by 866 seats in total |


## The sell-out race — seats sold per second (higher is better)

Every buyer arrives at an unsold event at the same instant and keeps buying until told
"sold out". Throughput is successful sales over the wall time of the races. Markers:
❌ overbooked (events/total, extra seats) and struck through — **not a valid speed**;
⚠️ underbooked (every buyer was sent away with seats unsold); ⏱ race hit its timeout.

### PostgreSQL, 1 node

| Design | 10 seats | 100 seats | 1k seats | 10k seats |
|---|---:|---:|---:|---:|
| P2 skip-locked | 705 | 4.1k | 3.8k | 4.0k |
| C1 count (RC) ✗ | ~~2.3k~~ ❌ overbooked 20/20 (+442) | ~~8.8k~~ ❌ overbooked 4/4 (+112) | ~~4.7k~~ ❌ overbooked 1/1 (+29) | ~~2.0k~~ ❌ overbooked 1/1 (+31) |
| H1 hold/checked | 487 | 783 | 791 | 655 |


## The race in detail

How each sale was achieved, not just how many. *Attempts per sale* is (sales + retries) ÷
sales: 1.00 means no buyer ever lost a race and started again. *Sold out answered in* is the
median latency of the final "no" each buyer received. *Organiser edit p99* is the latency of an
unrelated update to the event page made while the crowd was buying.

### PostgreSQL, 1 node

| Sale latency p99 (ms) | 10 seats | 100 seats | 1k seats | 10k seats |
|---|---:|---:|---:|---:|
| P2 skip-locked | 33 | 41 | 65 | 66 |
| C1 count (RC) ✗ | 29 ❌ | 11 ❌ | 62 ❌ | 83 ❌ |
| H1 hold/checked | 31 | 128 | 167 | 233 |

| Attempts per sale | 10 seats | 100 seats | 1k seats | 10k seats |
|---|---:|---:|---:|---:|
| P2 skip-locked | 20.72 | 2.13 | 1.09 | 1.00 |
| C1 count (RC) ✗ | 1.00 ❌ | 1.00 ❌ | 1.00 ❌ | 1.00 ❌ |
| H1 hold/checked | 1.00 | 1.00 | 1.00 | 1.00 |

| Sold out answered in, p50 (ms) | 10 seats | 100 seats | 1k seats | 10k seats |
|---|---:|---:|---:|---:|
| P2 skip-locked | 4.70 | 3.47 | 4.01 | 2.54 |
| C1 count (RC) ✗ | 0.60 ❌ | 0.47 ❌ | 1.00 ❌ | 2.25 ❌ |
| H1 hold/checked | 14 | 27 | 31 | 32 |

| Organiser edit p99 (ms) | 10 seats | 100 seats | 1k seats | 10k seats |
|---|---:|---:|---:|---:|
| P2 skip-locked | 24 | 4.70 | 4.62 | 50 |
| C1 count (RC) ✗ | 23 ❌ | 1.98 ❌ | 2.48 ❌ | 76 ❌ |
| H1 hold/checked | 17 | 52 | 150 | 214 |


## Churn race — selling out while buyers cancel

As the race, but a share of successful buyers cancel immediately. After the crowd leaves,
one last buyer sweeps up refunded seats; seats still unsold after that sweep are
**underbooked** — the design said "sold out" while it had seats.

### PostgreSQL, 1 node

| Design | 10 seats | 100 seats | 1k seats |
|---|---:|---:|---:|
| P2 skip-locked | 1.5k | 5.5k | 4.6k |
| C1 count (RC) ✗ | ~~5.2k~~ ❌ overbooked 10/10 (+189) | ~~9.2k~~ ❌ overbooked 2/2 (+44) | ~~5.3k~~ ❌ overbooked 1/1 (+19) |
| H1 hold/checked | 623 | 807 | 673 |

| Design | 10 seats — cancels / swept up / underbooked seats | 100 seats | 1k seats |
|---|---:|---:|---:|
| P2 skip-locked | 7 / 0 / 0 | 30 / 0 / 0 | 103 / 0 / 0 |
| C1 count (RC) ✗ | 33 / 0 / 0 | 30 / 0 / 0 | 103 / 0 / 0 |
| H1 hold/checked | 8 / 0 / 0 | 19 / 0 / 0 | 105 / 0 / 0 |


## Pre-created or created: the cost outside the race

*Publish* creates a new event with its inventory in one transaction — one row, or one row
per seat — and is shown as the median time to publish one event. *Book (spread)* is demand
spread across the catalogue in proportion to seats left: little contention on any one event.
*Cancel* refunds a sold ticket. Each operation ran on its own fresh load.

### PostgreSQL, 1 node

| Design | publish 10 seats (ms) | publish 100 seats (ms) | publish 1k seats (ms) | publish 10k seats (ms) | publish 100k seats (ms) | book spread /s | cancel /s | load total | storage |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| P2 skip-locked | 1.31 | 1.92 | 8.74 | 50 | 510 | 6.1k | 6.4k | 0.1s | 4.3 MiB |
| C1 count (RC) ✗ | 0.88 | 0.89 | 0.85 | 0.91 | 0.65 | 3.7k | 7.5k | 0.1s | 1.0 MiB |
| H1 hold/checked | 0.92 | 0.89 | 1.08 | 1.12 | 0.93 | 1.5k | 1.9k | 0.1s | 1.1 MiB |

Storage is PostgreSQL only: YugabyteDB keeps data in DocDB, where `pg_total_relation_size()`
does not describe it.


## Reads (ops/s, higher is better)

Each question benchmarked in isolation with a fresh key per execution. Availability is
measured separately for every event size, because that is where counting designs pay.

### PostgreSQL, 1 node

| Question | P2 skip-locked | C1 count (RC) ✗ | H1 hold/checked |
|---|---:|---:|---:|
| Seats left — 10 seats | 26k | 40k | 51k |
| Seats left — 100 seats | 23k | 38k | 51k |
| Seats left — 1k seats | 17k | 25k | 52k |
| Seats left — 10k seats | 5.8k | 8.8k | 51k |
| A band's events with availability | 10k | 22k | 39k |
| A customer's tickets | 43k | 41k | 40k |
| Ticket by id (control) | 51k | 50k | 50k |
| Tickets sold across a band | 5.1k | 22k | 23k |


## Holds with expiry (H0, H1)

Buyers hold a seat, spend time in the basket, abandon some holds, and pay for the rest; a
sweeper releases expired holds. Some payments finish after the hold expired. *Late confirms
rejected* are those the design refused (and would refund); a design that accepts them can sell a
released seat twice.

### PostgreSQL, 1 node

| Design | Event size | Events | Confirmed/s | Holds | Abandoned | Expired at pre-check | Late confirms rejected | Swept | Hold p99 (ms) | Confirm p99 (ms) | Overbooked | Errors |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|---:|
| H1 hold/checked | 10 seats | 10 | 11.1 | 154 | 42 | 1 | 11 | 54 | 18 | 5.73 | ✅ none | 0 |
| H1 hold/checked | 100 seats | 4 | 59.1 | 611 | 170 | 0 | 41 | 211 | 35 | 3.30 | ✅ none | 0 |
| H1 hold/checked | 1k seats | 1 | 124 | 1527 | 424 | 0 | 103 | 527 | 38 | 2.10 | ✅ none | 0 |


## Controlled pairs — one decision at a time

The C1/C2/C3 read control is not present in this run, so read differences are filtered at a
conventional 1.30x rather than a measured error bar.


---

Plans for every read and write statement are in `results/devcheck-1/<topology>/plans/<design>.txt`
(`EXPLAIN (ANALYZE, BUFFERS)` on PostgreSQL, `EXPLAIN (ANALYZE, DIST)` on YugabyteDB; write
statements inside rolled-back transactions).

## Provenance: code, data and analyses

This file is **generated from the measurements** and contains no interpretation.
Conclusions live in separate signed analyses, so several analysts — including
different AI models — can read the same numbers and each record what they make of
them. Where two analyses disagree, the disagreement is a finding and stays visible.

| | |
|---|---|
| Run id | `devcheck-1` |
| Result files | 3 |
| **Inputs digest** | `d93299a82c5fa358` |
| Repository version | `unknown` (`unknown`), clean |

Quote the digest in any analysis: it identifies the data. The repository version
identifies the code — check it out (`git checkout <commit>`) to reproduce the run, or
`git diff <commit>` to see what has changed in the study since.

### Analyses of this run

**None yet.** This run has measurements but no interpretation.

### Adding an analysis

Copy `reports/analyses/TEMPLATE.md`, fill in the frontmatter, and regenerate this report.
If an analysis already exists for digest `d93299a82c5fa358`, read it first and write a new one only to
**disagree, extend, or bring a different perspective**. Never edit another analyst's file.
