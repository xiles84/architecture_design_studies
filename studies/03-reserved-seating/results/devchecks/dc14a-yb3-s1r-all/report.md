# Study 03 — results: reserved seating (venue → event → seat)

Generated from `results/dc14a-yb3-s1r-all`. Every number comes from exactly one JSON file in that
directory. **This file contains measurements and no interpretation**; conclusions live in
signed analyses, indexed at the end.

## TL;DR — measured facts

*Selected from the tables below by fixed rules; no interpretation. Read a signed analysis for
what these facts mean.*

- **Cells:** 1 across 1 topologies; 0 failed.
- **Negative controls:** 0 of 0 fired as designed to.
- **Invariant violations in designs meant to be correct:** none observed.
- **Refused confirmations, S1r retry, YugabyteDB, 3 nodes (RF=3)** (lifecycle, all tiers): 2 late, 0 boundary, 0 early (0 transient); 139 confirmed; confirmations retried 2, of which sold 0.
- **S1r confirmation retries, YugabyteDB, 3 nodes (RF=3)** (race, all tiers and trials): 17 retried, 17 of them sold.
- **Race, YugabyteDB, 3 nodes (RF=3)** (seats held/s; designs meant to be correct, without violations):
  - 10 seats: highest S1r retry 32.3/s, lowest S1r retry 32.3/s
  - 100 seats: highest S1r retry 76.9/s, lowest S1r retry 76.9/s
  - 1k seats: highest S1r retry 65.9/s, lowest S1r retry 65.9/s
  - 10k seats: highest S1r retry 37.5/s, lowest S1r retry 37.5/s (1 of 1 timed out)

## What was measured

|  |  |
|---|---|
| Run id | `devchecks/dc14a-yb3-s1r-all` |
| Environment | [`host-zenbook-ux5406sa`](../../../docs/environments/host-zenbook-ux5406sa.md) |
| Scale | `tiny` — 5 venues, 84 events, 5741 tickets sold, 131 live and 34 expired holds at load |
| Events (kind_seats: count) | catalogue_10: 40 · catalogue_100: 8 · catalogue_1000: 2 · catalogue_10000: 1 · lifecycle_10: 4 · lifecycle_100: 2 · lifecycle_1000: 1 · race_10: 20 · race_100: 4 · race_1000: 1 · race_10000: 1 |
| Seed | 42 — identical data in every cell |
| Reads / isolated writes | 8 workers, 5s measured + 2s warmup, 1 trial(s) |
| Race | 32 buyers per event, real hold TTL 40m0s, 10% of holds deferred until the crowd has gone, timeout 1m0s, tier budget 3m0s, give up after 20 conflicts |
| Lifecycle — YugabyteDB, 3 nodes (RF=3) | one human minute = 500ms; TTL 40m0s, payment window (K1) 10m0s, guard G 2m0s, sweeper every 1m0s, outage 45m0s→1h25m0s in the first event of each tier, 25% abandoned (40% of those explicitly), E0 skew 10m0s; tiers 10,100, 3 events per tier, 32 buyers |

**Engines**

| Topology | Version | READ COMMITTED runs as | Client connections to | DB clock − client clock (ms) |
|---|---|---|---|---:|
| YugabyteDB, 3 nodes (RF=3) | `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit` | read committed | 3 node(s) | 0.03 |

> Limits of this environment that bound every number below: the client and the databases
> share one laptop's eight cores; the "3-node" cluster has no real network between nodes; the
> race and the lifecycle are closed loops, so tails are a floor; the lifecycle runs on
> compressed time; E0's clock skew is injected.

## Correctness gate

Six read questions checked against Go-computed truth on a load with sold seats, live holds and
expired holds, plus an audit of the loaded state. A failure aborts the cell.

| Design | YugabyteDB, 3 nodes (RF=3) |
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

### YugabyteDB, 3 nodes (RF=3)

| Design | Tier | Confirmed seats/s | Holds | Abandoned | Expired at check | Refused late / boundary / early (transient) | Outage: unavailable / probed | Leaked | Release lag p50/p99 (min) | Idle held seat-min | Violations |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| S1r retry | 10 seats | 0.2 | 34 | 6+5 | 2 | 0 / 0 / 0 (0) | 0 / 0 | 0 | 0.0 / 0.0 | 624 | none |
| S1r retry | 100 seats | 1.1 | 187 | 38+18 | 11 | 2 / 0 / 0 (0) | 0 / 0 | 0 | 0.2 / 0.7 | 3475 | none |

## The race — a hot drop with seat choice

Seats held per second, from the gate to the last hold granted. Conflicts and seat-map reads are
per granted hold. Deferred confirmations kept a real 40-minute hold through the crowd.

### YugabyteDB, 3 nodes (RF=3)

| Design | 10 seats | 100 seats | 1k seats | 10k seats |
|---|---:|---:|---:|---:|
| S1r retry | 32.3 | 76.9 | 65.9 | 37.5 ⏱ 1 timed out at 23% sold |

<details><summary>Race detail — YugabyteDB, 3 nodes (RF=3)</summary>

| Design | Tier | Conflicts/hold | Map reads/hold | Engine retries | Gave up | Hold p50/p99 ms | Confirm p99 ms | Deferred ok | Early rej. (transient) | Confirm retries run / sold | Final buyer sold |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| S1r retry | 10 seats | 6.65 | 14.45 | 0 | 0 | 34 / 299 | 575 | 11/11 | 0 (0) | 0 / 0 | 0 |
| S1r retry | 100 seats | 2.32 | 4.47 | 0 | 0 | 47 / 288 | 411 | 12/12 | 0 (0) | 2 / 2 | 0 |
| S1r retry | 1k seats | 1.04 | 3.09 | 0 | 0 | 87 / 291 | 326 | 46/46 | 0 (0) | 6 / 6 | 0 |
| S1r retry | 10k seats | 0.38 | 1.45 | 0 | 0 | 194 / 487 | 687 | 118/118 | 0 (0) | 9 / 9 | 0 |

</details>

## Reads

Each question in isolation, fresh key per execution, ops/s.

### YugabyteDB, 3 nodes (RF=3)

| Question | S1r retry |
|---|---:|
| Seat map of a section — 10 seats | 2.8k |
| Seat map of a section — 100 seats | 2.6k |
| Seat map of a section — 1k seats | 1.2k |
| Seat map of a section — 10k seats | 181 |
| Seats available per section — 10 seats | 2.6k |
| Seats available per section — 100 seats | 2.2k |
| Seats available per section — 1k seats | 1.2k |
| Seats available per section — 10k seats | 134 |
| Seats available for an event — 10 seats | 2.9k |
| Seats available for an event — 100 seats | 2.3k |
| Seats available for an event — 1k seats | 1.1k |
| Seats available for an event — 10k seats | 170 |
| A hold's seats (basket page) | 214 |
| A customer's tickets | 1.6k |
| Ticket by id (control) | 3.6k |

## Isolated writes, publishing and storage

### YugabyteDB, 3 nodes (RF=3)

| Design | Hold+confirm /s | Release /s | Refund /s | Publish p50 ms by size | Storage | Audit |
|---|---:|---:|---:|---:|---:|---|
| S1r retry | 19.4 | 419 | 240 | 10 seats: 20 · 100 seats: 23 · 1k seats: 56 · 10k seats: 360 · 100k seats: 3267 | n/a on YugabyteDB | consistent |


## Controlled pairs — one decision at a time


## Measurement hazard: CPU-quota throttling

Cells show *throttled periods / periods (throttled time)*: for the **database** container(s) over
the whole cell (summed across nodes), and for the **client** per phase.

### YugabyteDB, 3 nodes (RF=3)

| Design | database, whole cell | client lifecycle#1 | client race#1 | client read | client write:cancel | client write:hold | client write:publish | client write:release |
|---|---:|---:|---:|---:|---:|---:|---:|---:|
| S1r retry | 2900/23157 (577s) | 0/3198 (0.0s) | 0/944 (0.0s) | 0/1053 (0.0s) | 0/108 (0.0s) | 0/346 (0.0s) | 0/148 (0.0s) | 0/4 (0.0s) |


---

Plans for every read and write statement are in `results/dc14a-yb3-s1r-all/<topology>/plans/<design>.txt`
(`EXPLAIN (ANALYZE, BUFFERS)` on PostgreSQL, `EXPLAIN (ANALYZE, DIST)` on YugabyteDB; write
statements inside rolled-back transactions).

## Provenance: code, data and analyses

This file is **generated from the measurements** and contains no interpretation.
Conclusions live in separate signed analyses, so several analysts — including
different AI models — can read the same numbers and each record what they make of
them. Where two analyses disagree, the disagreement is a finding and stays visible.

| | |
|---|---|
| Run id | `devchecks/dc14a-yb3-s1r-all` |
| Result files | 1 |
| **Inputs digest** | `5da9a95bf6f62c21` |
| Repository version | `study-03/v0.1-handoff-amendment-01-6-g425966e-dirty` (`425966e9452b`), ⚠️ **dirty** — uncommitted changes; not reproducible from the commit alone |

Quote the digest in any analysis: it identifies the data. The repository version
identifies the code — check it out (`git checkout <commit>`) to reproduce the run, or
`git diff <commit>` to see what has changed in the study since.

### Analyses of this run

**None yet.** This run has measurements but no interpretation.

### Adding an analysis

Copy `reports/analyses/TEMPLATE.md`, fill in the frontmatter, and regenerate this report.
If an analysis already exists for digest `5da9a95bf6f62c21`, read it first and write a new one only to
**disagree, extend, or bring a different perspective**. Never edit another analyst's file.
