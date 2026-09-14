# Study 03 — results: reserved seating (venue → event → seat)

Generated from `results/dc13a-yb1-s1r-all`. Every number comes from exactly one JSON file in that
directory. **This file contains measurements and no interpretation**; conclusions live in
signed analyses, indexed at the end.

## TL;DR — measured facts

*Selected from the tables below by fixed rules; no interpretation. Read a signed analysis for
what these facts mean.*

- **Cells:** 1 across 1 topologies; 0 failed.
- **Negative controls:** 0 of 0 fired as designed to.
- **Invariant violations in designs meant to be correct:** none observed.
- **Refused confirmations, S1r retry, YugabyteDB, 1 node (RF=1)** (lifecycle, all tiers): 1 late, 0 boundary, 0 early (0 transient); 127 confirmed; confirmations retried 1, of which sold 0.
- **S1r confirmation retries, YugabyteDB, 1 node (RF=1)** (race, all tiers and trials): 0 retried, 0 of them sold.
- **Race, YugabyteDB, 1 node (RF=1)** (seats held/s; designs meant to be correct, without violations):
  - 10 seats: highest S1r retry 14.2/s, lowest S1r retry 14.2/s
  - 100 seats: highest S1r retry 41.3/s, lowest S1r retry 41.3/s
  - 1k seats: highest S1r retry 46.5/s, lowest S1r retry 46.5/s
  - 10k seats: highest S1r retry 29.6/s, lowest S1r retry 29.6/s (1 of 1 timed out)

## What was measured

|  |  |
|---|---|
| Run id | `devchecks/dc13a-yb1-s1r-all` |
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
| YugabyteDB, 1 node (RF=1) | `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit` | read committed | 1 node(s) | 0.01 |

> Limits of this environment that bound every number below: the client and the databases
> share one laptop's eight cores; the "3-node" cluster has no real network between nodes; the
> race and the lifecycle are closed loops, so tails are a floor; the lifecycle runs on
> compressed time; E0's clock skew is injected.

## Correctness gate

Six read questions checked against Go-computed truth on a load with sold seats, live holds and
expired holds, plus an audit of the loaded state. A failure aborts the cell.

| Design | YugabyteDB, 1 node (RF=1) |
|---|---|
| S1r retry | ✅ 51/51 |

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

| Design | Tier | Confirmed seats/s | Holds | Abandoned | Expired at check | Refused late / boundary / early (transient) | Outage: unavailable / probed | Leaked | Release lag p50/p99 (min) | Idle held seat-min | Violations |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| S1r retry | 10 seats | 0.3 | 23 | 2+1 | 2 | 0 / 0 / 0 (0) | 0 / 0 | 0 | 0.0 / 0.0 | 125 | none |
| S1r retry | 100 seats | 1.7 | 152 | 21+15 | 6 | 1 / 0 / 0 (0) | 0 / 0 | 0 | 0.0 / 0.4 | 2129 | none |

## The race — a hot drop with seat choice

Seats held per second, from the gate to the last hold granted. Conflicts and seat-map reads are
per granted hold. Deferred confirmations kept a real 40-minute hold through the crowd.

### YugabyteDB, 1 node (RF=1)

| Design | 10 seats | 100 seats | 1k seats | 10k seats |
|---|---:|---:|---:|---:|
| S1r retry | 14.2 | 41.3 | 46.5 | 29.6 ⏱ 1 timed out at 18% sold |

<details><summary>Race detail — YugabyteDB, 1 node (RF=1)</summary>

| Design | Tier | Conflicts/hold | Map reads/hold | Engine retries | Gave up | Hold p50/p99 ms | Confirm p99 ms | Deferred ok | Early rej. (transient) | Confirm retries run / sold | Final buyer sold |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| S1r retry | 10 seats | 6.10 | 13.59 | 0 | 0 | 95 / 771 | 1313 | 9/9 | 0 (0) | 0 / 0 | 0 |
| S1r retry | 100 seats | 2.45 | 5.27 | 0 | 0 | 101 / 623 | 371 | 22/22 | 0 (0) | 0 / 0 | 0 |
| S1r retry | 1k seats | 0.92 | 2.86 | 0 | 0 | 180 / 490 | 602 | 44/44 | 0 (0) | 0 / 0 | 0 |
| S1r retry | 10k seats | 0.34 | 1.40 | 0 | 0 | 309 / 807 | 1100 | 73/73 | 0 (0) | 0 / 0 | 0 |

</details>

## Reads

Each question in isolation, fresh key per execution, ops/s.

### YugabyteDB, 1 node (RF=1)

| Question | S1r retry |
|---|---:|
| Seat map of a section — 10 seats | 1.4k |
| Seat map of a section — 100 seats | 1.1k |
| Seat map of a section — 1k seats | 790 |
| Seat map of a section — 10k seats | 144 |
| Seats available per section — 10 seats | 1.1k |
| Seats available per section — 100 seats | 900 |
| Seats available per section — 1k seats | 668 |
| Seats available per section — 10k seats | 76.7 |
| Seats available for an event — 10 seats | 1.2k |
| Seats available for an event — 100 seats | 1.1k |
| Seats available for an event — 1k seats | 850 |
| Seats available for an event — 10k seats | 140 |
| A hold's seats (basket page) | 179 |
| A customer's tickets | 880 |
| Ticket by id (control) | 1.4k |

## Isolated writes, publishing and storage

### YugabyteDB, 1 node (RF=1)

| Design | Hold+confirm /s | Release /s | Refund /s | Publish p50 ms by size | Storage | Audit |
|---|---:|---:|---:|---:|---:|---|
| S1r retry | 22.7 | 107 | 103 | 10 seats: 31 · 100 seats: 80 · 1k seats: 204 · 10k seats: 660 · 100k seats: 9742 | n/a on YugabyteDB | consistent |


## Controlled pairs — one decision at a time


## Measurement hazard: CPU-quota throttling

Cells show *throttled periods / periods (throttled time)*: for the **database** container(s) over
the whole cell (summed across nodes), and for the **client** per phase.

### YugabyteDB, 1 node (RF=1)

| Design | database, whole cell | client lifecycle#1 | client race#1 | client read | client write:cancel | client write:hold | client write:publish | client write:release |
|---|---:|---:|---:|---:|---:|---:|---:|---:|
| S1r retry | 3763/6985 (1441s) | 0/2176 (0.0s) | 0/1153 (0.0s) | 0/1056 (0.0s) | 0/199 (0.0s) | 0/422 (0.0s) | 0/280 (0.0s) | 0/7 (0.0s) |


---

Plans for every read and write statement are in `results/dc13a-yb1-s1r-all/<topology>/plans/<design>.txt`
(`EXPLAIN (ANALYZE, BUFFERS)` on PostgreSQL, `EXPLAIN (ANALYZE, DIST)` on YugabyteDB; write
statements inside rolled-back transactions).

## Provenance: code, data and analyses

This file is **generated from the measurements** and contains no interpretation.
Conclusions live in separate signed analyses, so several analysts — including
different AI models — can read the same numbers and each record what they make of
them. Where two analyses disagree, the disagreement is a finding and stays visible.

| | |
|---|---|
| Run id | `devchecks/dc13a-yb1-s1r-all` |
| Result files | 1 |
| **Inputs digest** | `726948824a224669` |
| Repository version | `study-03/v0.1-handoff-amendment-01-4-gf4bf9f8` (`f4bf9f82c972`), clean |

Quote the digest in any analysis: it identifies the data. The repository version
identifies the code — check it out (`git checkout <commit>`) to reproduce the run, or
`git diff <commit>` to see what has changed in the study since.

### Analyses of this run

**None yet.** This run has measurements but no interpretation.

### Adding an analysis

Copy `reports/analyses/TEMPLATE.md`, fill in the frontmatter, and regenerate this report.
If an analysis already exists for digest `726948824a224669`, read it first and write a new one only to
**disagree, extend, or bring a different perspective**. Never edit another analyst's file.
