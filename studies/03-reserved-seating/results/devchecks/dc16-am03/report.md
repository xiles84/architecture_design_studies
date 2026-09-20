# Study 03 — results: reserved seating (venue → event → seat)

Generated from `results/dc16-am03`. Every number comes from exactly one JSON file in that
directory. **This file contains measurements and no interpretation**; conclusions live in
signed analyses, indexed at the end.

## TL;DR — measured facts

*Selected from the tables below by fixed rules; no interpretation. Read a signed analysis for
what these facts mean.*

- **Cells:** 2 across 1 topologies; 0 failed.
- **Negative controls:** 0 of 0 fired as designed to.
- **Invariant violations in designs meant to be correct:** none observed.
- **Race, YugabyteDB, 1 node (RF=1)** (seats held/s; designs meant to be correct, without violations):
  - 10 seats: highest L2 document 34.1/s, lowest E2 cart 19.4/s
  - 100 seats: highest L2 document 54.9/s, lowest E2 cart 54.0/s
  - 1k seats: highest L2 document 89.5/s, lowest E2 cart 80.9/s
  - 10k seats: highest L2 document 135/s, lowest E2 cart 63.7/s (2 of 2 timed out)

## What was measured

|  |  |
|---|---|
| Run id | `devchecks/dc16-am03` |
| Environment | [`host-zenbook-ux5406sa`](../../../docs/environments/host-zenbook-ux5406sa.md) |
| Scale | `tiny` — 5 venues, 84 events, 5741 tickets sold, 131 live and 34 expired holds at load |
| Events (kind_seats: count) | catalogue_10: 40 · catalogue_100: 8 · catalogue_1000: 2 · catalogue_10000: 1 · lifecycle_10: 4 · lifecycle_100: 2 · lifecycle_1000: 1 · race_10: 20 · race_100: 4 · race_1000: 1 · race_10000: 1 |
| Seed | 42 — identical data in every cell |
| Reads / isolated writes | 8 workers, 5s measured + 2s warmup, 1 trial(s) |
| Race | 32 buyers per event, real hold TTL 40m0s, 10% of holds deferred until the crowd has gone, timeout 1m0s, tier budget 3m0s, give up after 20 conflicts |
| Lifecycle — YugabyteDB, 1 node (RF=1) | one human minute = 500ms; TTL 40m0s, payment window (K1) 10m0s, guard G 2m0s, sweeper every 1m0s, outage 45m0s→1h25m0s in the first event of each tier, 25% abandoned (40% of those explicitly), E0 skew 10m0s; tiers 10,100, 3 events per tier, 32 buyers |

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
| E2 cart | ✅ 51/51 |
| L2 document | ✅ 51/51 |

## Negative controls

Deliberately wrong designs (and E1 with its sweeper stopped). Each must be seen to break the
invariant it exists to test; otherwise the absence of violations elsewhere in that experiment
is not evidence of correctness.

No control was measured in this run.

## The hold guarantee — lifecycle

Compressed time. *Rejected* confirmations are classified by how long before the hold's expiry
the confirming transaction started: **late** (after expiry), **boundary** (less than G before),
**early** (G or more before — a violation). Violations are counted by the ledger and the audit.
Release lag is in human minutes.

### YugabyteDB, 1 node (RF=1)

| Design | Tier | Confirmed seats/s | Holds | Abandoned | Expired at check | Refused late / boundary / early (transient) | Outage: unavailable / probed | Leaked | Release lag p50/p99 (min) | Idle held seat-min | Monitor retries / events ended early | Violations |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| E2 cart | 10 seats | 0.3 | 29 | 3+2 | 1 | 2 / 0 / 0 (0) | 0 / 0 | 0 | 0.0 / 0.0 | 236 | 0 / 0 | none |
| E2 cart | 100 seats | 1.7 | 160 | 33+15 | 2 | 3 / 0 / 0 (0) | 0 / 0 | 0 | 0.2 / 0.4 | 2893 | 0 / 0 | none |
| L2 document | 10 seats | 0.4 | 21 | 2+0 | 0 | 2 / 0 / 0 (0) | 0 / 0 | 0 | 0.0 / 0.0 | 120 | 0 / 0 | none |
| L2 document | 100 seats | 1.5 | 143 | 20+14 | 4 | 2 / 0 / 0 (0) | 0 / 0 | 0 | 0.0 / 0.0 | 1797 | 0 / 0 | none |

## The race — a hot drop with seat choice

Seats held per second, from the gate to the last hold granted. Conflicts and seat-map reads are
per granted hold. Deferred confirmations kept a real 40-minute hold through the crowd.

### YugabyteDB, 1 node (RF=1)

| Design | 10 seats | 100 seats | 1k seats | 10k seats |
|---|---:|---:|---:|---:|
| E2 cart | 19.4 | 54.0 | 80.9 | 63.7 ⏱ 1 timed out at 38% sold |
| L2 document | 34.1 | 54.9 | 89.5 | 135 ⏱ 1 timed out at 81% sold |

<details><summary>Race detail — YugabyteDB, 1 node (RF=1)</summary>

| Design | Tier | Conflicts/hold | Map reads/hold | Engine retries | Gave up | Hold p50/p99 ms | Confirm p99 ms | Deferred ok | Early rej. (transient) | Confirm retries run / sold | Final buyer sold |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| E2 cart | 10 seats | 7.46 | 16.12 | 0 | 0 | 80 / 599 | 378 | 7/7 | 0 (0) | 0 / 0 | 0 |
| E2 cart | 100 seats | 2.21 | 4.24 | 0 | 0 | 88 / 400 | 605 | 16/16 | 0 (0) | 0 / 0 | 0 |
| E2 cart | 1k seats | 0.74 | 2.67 | 0 | 0 | 94 / 216 | 201 | 52/52 | 0 (0) | 0 / 0 | 0 |
| E2 cart | 10k seats | 0.22 | 1.40 | 0 | 0 | 109 / 299 | 604 | 151/151 | 0 (0) | 0 / 0 | 0 |
| L2 document | 10 seats | 14.05 | 24.30 | 2424 | 0 | 47 / 203 | 1064 | 7/7 | 0 (0) | 0 / 0 | 0 |
| L2 document | 100 seats | 4.91 | 7.94 | 3472 | 0 | 58 / 375 | 1893 | 21/21 | 0 (0) | 0 / 0 | 0 |
| L2 document | 1k seats | 1.41 | 3.36 | 4166 | 1 | 63 / 461 | 5210 | 37/37 | 0 (0) | 0 / 0 | 0 |
| L2 document | 10k seats | 0.60 | 1.95 | 16237 | 0 | 75 / 810 | 5382 | 330/330 | 0 (0) | 0 / 0 | 0 |

</details>

## Reads

Each question in isolation, fresh key per execution, ops/s.

### YugabyteDB, 1 node (RF=1)

| Question | E2 cart | L2 document |
|---|---:|---:|
| Seat map of a section — 10 seats | 2.7k | 4.1k |
| Seat map of a section — 100 seats | 1.8k | 3.8k |
| Seat map of a section — 1k seats | 1.2k | 3.6k |
| Seat map of a section — 10k seats | 328 | 3.3k |
| Seats available per section — 10 seats | 1.8k | 3.6k |
| Seats available per section — 100 seats | 1.9k | 4.0k |
| Seats available per section — 1k seats | 1.4k | 3.2k |
| Seats available per section — 10k seats | 205 | 1.0k |
| Seats available for an event — 10 seats | 1.8k | 4.0k |
| Seats available for an event — 100 seats | 1.8k | 3.9k |
| Seats available for an event — 1k seats | 1.5k | 3.1k |
| Seats available for an event — 10k seats | 260 | 1.0k |
| A hold's seats (basket page) | 586 | 3.9k |
| A customer's tickets | 2.9k | 2.9k |
| Ticket by id (control) | 5.1k | 5.0k |

## Isolated writes, publishing and storage

### YugabyteDB, 1 node (RF=1)

| Design | Hold+confirm /s | Release /s | Refund /s | Publish p50 ms by size | Storage | Audit |
|---|---:|---:|---:|---:|---:|---|
| E2 cart | 47.3 | 346 | 369 | 10 seats: 5.36 · 100 seats: 9.57 · 1k seats: 95 · 10k seats: 263 · 100k seats: 2843 | n/a on YugabyteDB | consistent |
| L2 document | 110 | 409 | 185 | 10 seats: 4.60 · 100 seats: 4.54 · 1k seats: 4.32 · 10k seats: 3.68 · 100k seats: 5.48 | n/a on YugabyteDB | consistent |


## Controlled pairs — one decision at a time


## Measurement hazard: CPU-quota throttling

Cells show *throttled periods / periods (throttled time)*: for the **database** container(s) over
the whole cell (summed across nodes), and for the **client** per phase.

### YugabyteDB, 1 node (RF=1)

| Design | database, whole cell | client lifecycle#1 | client race#1 | client read | client write:cancel | client write:hold | client write:publish | client write:release |
|---|---:|---:|---:|---:|---:|---:|---:|---:|
| E2 cart | 2548/5380 (693s) | 0/2098 (0.0s) | 0/964 (0.0s) | 0/1053 (0.0s) | 0/78 (0.0s) | 0/186 (0.0s) | 0/99 (0.0s) | 0/2 (0.0s) |
| L2 document | 2324/4912 (559s) | 0/2043 (0.0s) | 0/933 (0.0s) | 0/1051 (0.0s) | 0/120 (0.0s) | 0/111 (0.0s) | 0/67 (0.0s) | 0/1 (0.0s) |


---

Plans for every read and write statement are in `results/dc16-am03/<topology>/plans/<design>.txt`
(`EXPLAIN (ANALYZE, BUFFERS)` on PostgreSQL, `EXPLAIN (ANALYZE, DIST)` on YugabyteDB; write
statements inside rolled-back transactions).

## Provenance: code, data and analyses

This file is **generated from the measurements** and contains no interpretation.
Conclusions live in separate signed analyses, so several analysts — including
different AI models — can read the same numbers and each record what they make of
them. Where two analyses disagree, the disagreement is a finding and stays visible.

| | |
|---|---|
| Run id | `devchecks/dc16-am03` |
| Result files | 2 |
| **Inputs digest** | `5a35fec44c5024d1` |
| Repository version | `study-03/v0.3-handoff-amendment-03-1-g6d17f3f` (`6d17f3f48247`), clean |

Quote the digest in any analysis: it identifies the data. The repository version
identifies the code — check it out (`git checkout <commit>`) to reproduce the run, or
`git diff <commit>` to see what has changed in the study since.

### Analyses of this run

**None yet.** This run has measurements but no interpretation.

### Adding an analysis

Copy `reports/analyses/TEMPLATE.md`, fill in the frontmatter, and regenerate this report.
If an analysis already exists for digest `5a35fec44c5024d1`, read it first and write a new one only to
**disagree, extend, or bring a different perspective**. Never edit another analyst's file.
