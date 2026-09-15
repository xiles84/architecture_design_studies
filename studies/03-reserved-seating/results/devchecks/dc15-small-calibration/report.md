# Study 03 — results: reserved seating (venue → event → seat)

Generated from `results/dc15-small-calibration`. Every number comes from exactly one JSON file in that
directory. **This file contains measurements and no interpretation**; conclusions live in
signed analyses, indexed at the end.

## TL;DR — measured facts

*Selected from the tables below by fixed rules; no interpretation. Read a signed analysis for
what these facts mean.*

- **Cells:** 3 across 3 topologies; 0 failed.
- **Negative controls:** 0 of 0 fired as designed to.
- ❌ **Invariant violations in designs meant to be correct:** S1 conditional on YugabyteDB, 1 node (RF=1) (4 early rejections (4 transient)); S1 conditional on YugabyteDB, 3 nodes (RF=3) (43 early rejections (43 transient)).
- **Refused confirmations, S1 conditional, PostgreSQL, 1 node** (lifecycle, all tiers): 40 late, 0 boundary, 0 early (0 transient); 2806 confirmed.
- **Refused confirmations, S1 conditional, YugabyteDB, 1 node (RF=1)** (lifecycle, all tiers): 2 late, 0 boundary, 0 early (0 transient); 183 confirmed.
- **Early rejections in the race, correct designs, YugabyteDB, 1 node (RF=1):** S1 conditional 4 (4 transient).
- **Refused confirmations, S1 conditional, YugabyteDB, 3 nodes (RF=3)** (lifecycle, all tiers): 3 late, 0 boundary, 0 early (0 transient); 181 confirmed.
- **Early rejections in the race, correct designs, YugabyteDB, 3 nodes (RF=3):** S1 conditional 43 (43 transient).
- **Race, PostgreSQL, 1 node** (seats held/s; designs meant to be correct, without violations):
  - 10 seats: highest S1 conditional 329/s, lowest S1 conditional 329/s
  - 100 seats: highest S1 conditional 765/s, lowest S1 conditional 765/s
  - 1k seats: highest S1 conditional 1.3k/s, lowest S1 conditional 1.3k/s
  - 10k seats: highest S1 conditional 647/s, lowest S1 conditional 647/s
  - 100k seats: highest S1 conditional 303/s, lowest S1 conditional 303/s (1 of 1 timed out)
- **Race, YugabyteDB, 1 node (RF=1)** (seats held/s; designs meant to be correct, without violations):
  - 10 seats: highest S1 conditional 48.1/s, lowest S1 conditional 48.1/s
- **Race, YugabyteDB, 3 nodes (RF=3)** (seats held/s; designs meant to be correct, without violations):
  - 100k seats: highest S1 conditional 45.1/s, lowest S1 conditional 45.1/s (1 of 1 timed out)

## What was measured

|  |  |
|---|---|
| Run id | `devchecks/dc15-small-calibration` |
| Environment | [`host-zenbook-ux5406sa`](../../../docs/environments/host-zenbook-ux5406sa.md) |
| Scale | `small` — 5 venues, 658 events, 56103 tickets sold, 1793 live and 424 expired holds at load |
| Events (kind_seats: count) | catalogue_10: 400 · catalogue_100: 80 · catalogue_1000: 16 · catalogue_10000: 3 · catalogue_100000: 1 · lifecycle_10: 10 · lifecycle_100: 10 · lifecycle_1000: 10 · race_10: 100 · race_100: 20 · race_1000: 5 · race_10000: 2 · race_100000: 1 |
| Seed | 42 — identical data in every cell |
| Reads / isolated writes | 8 workers, 5s measured + 2s warmup, 1 trial(s) |
| Race | 32 buyers per event, real hold TTL 40m0s, 10% of holds deferred until the crowd has gone, timeout 1m0s, tier budget 3m0s, give up after 20 conflicts |
| Lifecycle — PostgreSQL, 1 node | one human minute = 50ms; TTL 40m0s, payment window (K1) 10m0s, guard G 2m0s, sweeper every 1m0s, outage 45m0s→1h25m0s in the first event of each tier, 25% abandoned (40% of those explicitly), E0 skew 10m0s; tiers 10,100,1000, 6 events per tier, 32 buyers |
| Lifecycle — YugabyteDB, 1 node (RF=1) | one human minute = 500ms; TTL 40m0s, payment window (K1) 10m0s, guard G 2m0s, sweeper every 1m0s, outage 45m0s→1h25m0s in the first event of each tier, 25% abandoned (40% of those explicitly), E0 skew 10m0s; tiers 10,100, 3 events per tier, 32 buyers |
| Lifecycle — YugabyteDB, 3 nodes (RF=3) | one human minute = 500ms; TTL 40m0s, payment window (K1) 10m0s, guard G 2m0s, sweeper every 1m0s, outage 45m0s→1h25m0s in the first event of each tier, 25% abandoned (40% of those explicitly), E0 skew 10m0s; tiers 10,100, 3 events per tier, 32 buyers |

**Engines**

| Topology | Version | READ COMMITTED runs as | Client connections to | DB clock − client clock (ms) |
|---|---|---|---|---:|
| PostgreSQL, 1 node | `PostgreSQL 17.11 (Debian 17.11-1.pgdg13+2) on x86_64-pc-linux-gnu, compiled by gcc (Debian 14.2.0-19) 14.2.0, 64-bit` | read committed | 1 node(s) | 0.00 |
| YugabyteDB, 1 node (RF=1) | `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit` | read committed | 1 node(s) | 0.00 |
| YugabyteDB, 3 nodes (RF=3) | `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit` | read committed | 3 node(s) | 0.01 |

> Limits of this environment that bound every number below: the client and the databases
> share one laptop's eight cores; the "3-node" cluster has no real network between nodes; the
> race and the lifecycle are closed loops, so tails are a floor; the lifecycle runs on
> compressed time; E0's clock skew is injected.

## Correctness gate

Six read questions checked against Go-computed truth on a load with sold seats, live holds and
expired holds, plus an audit of the loaded state. A failure aborts the cell.

| Design | PostgreSQL, 1 node | YugabyteDB, 1 node (RF=1) | YugabyteDB, 3 nodes (RF=3) |
|---|---|---|---|
| S1 conditional | ✅ 58/58 | ✅ 58/58 | ✅ 58/58 |

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

### PostgreSQL, 1 node

| Design | Tier | Confirmed seats/s | Holds | Abandoned | Expired at check | Refused late / boundary / early (transient) | Outage: unavailable / probed | Leaked | Release lag p50/p99 (min) | Idle held seat-min | Violations |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| S1 conditional | 10 seats | 3.0 | 57 | 7+5 | 9 | 0 / 0 / 0 (0) | 0 / 0 | 0 | 32.6 / 32.7 | 446 | none |
| S1 conditional | 100 seats | 18.1 | 443 | 69+50 | 8 | 5 / 0 / 0 (0) | 0 / 0 | 0 | 62.8 / 67.7 | 6477 | none |
| S1 conditional | 1k seats | 70.5 | 3464 | 514+350 | 106 | 35 / 0 / 0 (0) | 0 / 3 | 0 | 0.6 / 63.7 | 59237 | none |

### YugabyteDB, 1 node (RF=1)

| Design | Tier | Confirmed seats/s | Holds | Abandoned | Expired at check | Refused late / boundary / early (transient) | Outage: unavailable / probed | Leaked | Release lag p50/p99 (min) | Idle held seat-min | Violations |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| S1 conditional | 10 seats | 0.2 | 39 | 7+2 | 4 | 0 / 0 / 0 (0) | 0 / 0 | 0 | 0.1 / 0.2 | 426 | none |
| S1 conditional | 100 seats | 2.0 | 233 | 31+37 | 6 | 2 / 0 / 0 (0) | 0 / 0 | 0 | 0.7 / 0.9 | 2911 | none |

### YugabyteDB, 3 nodes (RF=3)

| Design | Tier | Confirmed seats/s | Holds | Abandoned | Expired at check | Refused late / boundary / early (transient) | Outage: unavailable / probed | Leaked | Release lag p50/p99 (min) | Idle held seat-min | Violations |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| S1 conditional | 10 seats | 0.3 | 21 | 2+0 | 1 | 0 / 0 / 0 (0) | 0 / 0 | 0 | 0.0 / 0.0 | 80 | none |
| S1 conditional | 100 seats | 1.6 | 244 | 32+33 | 13 | 3 / 0 / 0 (0) | 0 / 0 | 0 | 0.1 / 0.6 | 2880 | none |

## The race — a hot drop with seat choice

Seats held per second, from the gate to the last hold granted. Conflicts and seat-map reads are
per granted hold. Deferred confirmations kept a real 40-minute hold through the crowd.

### PostgreSQL, 1 node

| Design | 10 seats | 100 seats | 1k seats | 10k seats | 100k seats |
|---|---:|---:|---:|---:|---:|
| S1 conditional | 329 | 765 | 1.3k | 647 | 303 ⏱ 1 timed out at 18% sold |

<details><summary>Race detail — PostgreSQL, 1 node</summary>

| Design | Tier | Conflicts/hold | Map reads/hold | Engine retries | Gave up | Hold p50/p99 ms | Confirm p99 ms | Deferred ok | Early rej. (transient) | Confirm retries run / sold | Final buyer sold |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| S1 conditional | 10 seats | 8.76 | 16.56 | 0 | 0 | 3.43 / 50 | 49 | 70/70 | 0 (0) | 0 / 0 | 0 |
| S1 conditional | 100 seats | 3.05 | 5.70 | 0 | 0 | 3.53 / 67 | 71 | 87/87 | 0 (0) | 0 / 0 | 0 |
| S1 conditional | 1k seats | 1.44 | 4.16 | 0 | 1 | 3.43 / 79 | 79 | 245/245 | 0 (0) | 0 / 0 | 0 |
| S1 conditional | 10k seats | 0.40 | 3.81 | 0 | 0 | 5.85 / 93 | 96 | 852/852 | 0 (0) | 0 / 0 | 0 |
| S1 conditional | 100k seats | 0.11 | 1.27 | 0 | 0 | 22 / 179 | 187 | 729/729 | 0 (0) | 0 / 0 | 0 |

</details>

### YugabyteDB, 1 node (RF=1)

| Design | 10 seats | 100 seats | 1k seats | 10k seats | 100k seats |
|---|---:|---:|---:|---:|---:|
| S1 conditional | 48.1 | ~~98.8~~ ❌ 1 early rejections (1 transient) | ~~127~~ ❌ 1 early rejections (1 transient) | ~~114~~ ❌ 1 early rejections (1 transient); ⏱ 2 timed out at 69% sold | ~~36.8~~ ❌ 1 early rejections (1 transient); ⏱ 1 timed out at 2% sold |

<details><summary>Race detail — YugabyteDB, 1 node (RF=1)</summary>

| Design | Tier | Conflicts/hold | Map reads/hold | Engine retries | Gave up | Hold p50/p99 ms | Confirm p99 ms | Deferred ok | Early rej. (transient) | Confirm retries run / sold | Final buyer sold |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| S1 conditional | 10 seats | 7.25 | 15.20 | 0 | 0 | 15 / 210 | 301 | 63/63 | 0 (0) | 0 / 0 | 0 |
| S1 conditional | 100 seats | 2.74 | 5.13 | 0 | 0 | 15 / 208 | 186 | 100/100 | 1 (1) | 0 / 0 | 0 |
| S1 conditional | 1k seats | 0.96 | 3.20 | 0 | 0 | 84 / 177 | 308 | 206/206 | 1 (1) | 0 / 0 | 0 |
| S1 conditional | 10k seats | 0.18 | 1.55 | 0 | 0 | 94 / 196 | 285 | 531/531 | 1 (1) | 0 / 0 | 0 |
| S1 conditional | 100k seats | 0.04 | 1.04 | 0 | 0 | 94 / 197 | 390 | 78/78 | 1 (1) | 0 / 0 | 0 |

</details>

### YugabyteDB, 3 nodes (RF=3)

| Design | 10 seats | 100 seats | 1k seats | 10k seats | 100k seats |
|---|---:|---:|---:|---:|---:|
| S1 conditional | ~~46.2~~ ❌ 1 early rejections (1 transient) | ~~89.8~~ ❌ 3 early rejections (3 transient) | ~~132~~ ❌ 11 early rejections (11 transient) | ~~122~~ ❌ 28 early rejections (28 transient); ⏱ 2 timed out at 73% sold | 45.1 ⏱ 1 timed out at 3% sold |

<details><summary>Race detail — YugabyteDB, 3 nodes (RF=3)</summary>

| Design | Tier | Conflicts/hold | Map reads/hold | Engine retries | Gave up | Hold p50/p99 ms | Confirm p99 ms | Deferred ok | Early rej. (transient) | Confirm retries run / sold | Final buyer sold |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| S1 conditional | 10 seats | 7.53 | 15.41 | 0 | 0 | 29 / 252 | 270 | 59/59 | 1 (1) | 0 / 0 | 0 |
| S1 conditional | 100 seats | 2.87 | 5.32 | 0 | 0 | 37 / 273 | 193 | 99/99 | 3 (3) | 0 / 0 | 0 |
| S1 conditional | 1k seats | 1.04 | 3.17 | 0 | 0 | 44 / 176 | 184 | 206/206 | 11 (11) | 0 / 0 | 0 |
| S1 conditional | 10k seats | 0.17 | 1.48 | 0 | 0 | 76 / 176 | 277 | 613/613 | 28 (28) | 0 / 0 | 0 |
| S1 conditional | 100k seats | 0.03 | 1.03 | 0 | 0 | 71 / 164 | 214 | 106/106 | 0 (0) | 0 / 0 | 0 |

</details>

## Reads

Each question in isolation, fresh key per execution, ops/s.

### PostgreSQL, 1 node

| Question | S1 conditional |
|---|---:|
| Seat map of a section — 10 seats | 20k |
| Seat map of a section — 100 seats | 25k |
| Seat map of a section — 1k seats | 20k |
| Seat map of a section — 10k seats | 17k |
| Seat map of a section — 100k seats | 10k |
| Seats available per section — 10 seats | 20k |
| Seats available per section — 100 seats | 16k |
| Seats available per section — 1k seats | 8.7k |
| Seats available per section — 10k seats | 1.2k |
| Seats available per section — 100k seats | 144 |
| Seats available for an event — 10 seats | 18k |
| Seats available for an event — 100 seats | 17k |
| Seats available for an event — 1k seats | 9.9k |
| Seats available for an event — 10k seats | 1.3k |
| Seats available for an event — 100k seats | 153 |
| A hold's seats (basket page) | 9.0k |
| A customer's tickets | 34k |
| Ticket by id (control) | 39k |

### YugabyteDB, 1 node (RF=1)

| Question | S1 conditional |
|---|---:|
| Seat map of a section — 10 seats | 2.4k |
| Seat map of a section — 100 seats | 1.5k |
| Seat map of a section — 1k seats | 1.2k |
| Seat map of a section — 10k seats | 930 |
| Seat map of a section — 100k seats | 631 |
| Seats available per section — 10 seats | 3.3k |
| Seats available per section — 100 seats | 3.0k |
| Seats available per section — 1k seats | 1.6k |
| Seats available per section — 10k seats | 240 |
| Seats available per section — 100k seats | 25.8 |
| Seats available for an event — 10 seats | 3.9k |
| Seats available for an event — 100 seats | 3.7k |
| Seats available for an event — 1k seats | 2.2k |
| Seats available for an event — 10k seats | 472 |
| Seats available for an event — 100k seats | 52.0 |
| A hold's seats (basket page) | 764 |
| A customer's tickets | 2.7k |
| Ticket by id (control) | 4.7k |

### YugabyteDB, 3 nodes (RF=3)

| Question | S1 conditional |
|---|---:|
| Seat map of a section — 10 seats | 3.2k |
| Seat map of a section — 100 seats | 2.2k |
| Seat map of a section — 1k seats | 1.6k |
| Seat map of a section — 10k seats | 1.1k |
| Seat map of a section — 100k seats | 864 |
| Seats available per section — 10 seats | 4.9k |
| Seats available per section — 100 seats | 4.4k |
| Seats available per section — 1k seats | 2.5k |
| Seats available per section — 10k seats | 281 |
| Seats available per section — 100k seats | 31.4 |
| Seats available for an event — 10 seats | 5.5k |
| Seats available for an event — 100 seats | 4.7k |
| Seats available for an event — 1k seats | 2.6k |
| Seats available for an event — 10k seats | 514 |
| Seats available for an event — 100k seats | 58.0 |
| A hold's seats (basket page) | 1.0k |
| A customer's tickets | 2.7k |
| Ticket by id (control) | 5.4k |

## Isolated writes, publishing and storage

### PostgreSQL, 1 node

| Design | Hold+confirm /s | Release /s | Refund /s | Publish p50 ms by size | Storage | Audit |
|---|---:|---:|---:|---:|---:|---|
| S1 conditional | 125 | 3.7k | 3.6k | 10 seats: 5.61 · 100 seats: 9.72 · 1k seats: 15 · 10k seats: 67 · 100k seats: 798 | 55.5 MiB | consistent |

### YugabyteDB, 1 node (RF=1)

| Design | Hold+confirm /s | Release /s | Refund /s | Publish p50 ms by size | Storage | Audit |
|---|---:|---:|---:|---:|---:|---|
| S1 conditional | 16.4 | 557 | 335 | 10 seats: 4.73 · 100 seats: 7.38 · 1k seats: 96 · 10k seats: 333 · 100k seats: 3098 | n/a on YugabyteDB | consistent |

### YugabyteDB, 3 nodes (RF=3)

| Design | Hold+confirm /s | Release /s | Refund /s | Publish p50 ms by size | Storage | Audit |
|---|---:|---:|---:|---:|---:|---|
| S1 conditional | 17.6 | 527 | 263 | 10 seats: 17 · 100 seats: 27 · 1k seats: 77 · 10k seats: 316 · 100k seats: 3465 | n/a on YugabyteDB | consistent |


## Controlled pairs — one decision at a time


## Measurement hazard: CPU-quota throttling

Cells show *throttled periods / periods (throttled time)*: for the **database** container(s) over
the whole cell (summed across nodes), and for the **client** per phase.

### PostgreSQL, 1 node

| Design | database, whole cell | client lifecycle#1 | client race#1 | client read | client write:cancel | client write:hold | client write:publish | client write:release |
|---|---:|---:|---:|---:|---:|---:|---:|---:|
| S1 conditional | 3105/4583 (1321s) | 0/1365 (0.0s) | 0/1042 (0.0s) | 72/1266 (6.2s) | 0/74 (0.0s) | 0/671 (0.0s) | 0/41 (0.0s) | 0/2 (0.0s) |

### YugabyteDB, 1 node (RF=1)

| Design | database, whole cell | client lifecycle#1 | client race#1 | client read | client write:cancel | client write:hold | client write:publish | client write:release |
|---|---:|---:|---:|---:|---:|---:|---:|---:|
| S1 conditional | 10112/14266 (3194s) | 0/3019 (0.0s) | 0/2750 (0.0s) | 0/1246 (0.0s) | 0/457 (0.0s) | 0/5079 (0.0s) | 0/92 (0.0s) | 0/16 (0.0s) |

### YugabyteDB, 3 nodes (RF=3)

| Design | database, whole cell | client lifecycle#1 | client race#1 | client read | client write:cancel | client write:hold | client write:publish | client write:release |
|---|---:|---:|---:|---:|---:|---:|---:|---:|
| S1 conditional | 10834/44181 (2016s) | 0/2640 (0.0s) | 0/2795 (0.0s) | 0/1254 (0.0s) | 0/603 (0.0s) | 0/5046 (0.0s) | 0/136 (0.0s) | 0/19 (0.0s) |


---

Plans for every read and write statement are in `results/dc15-small-calibration/<topology>/plans/<design>.txt`
(`EXPLAIN (ANALYZE, BUFFERS)` on PostgreSQL, `EXPLAIN (ANALYZE, DIST)` on YugabyteDB; write
statements inside rolled-back transactions).

## Provenance: code, data and analyses

This file is **generated from the measurements** and contains no interpretation.
Conclusions live in separate signed analyses, so several analysts — including
different AI models — can read the same numbers and each record what they make of
them. Where two analyses disagree, the disagreement is a finding and stays visible.

| | |
|---|---|
| Run id | `devchecks/dc15-small-calibration` |
| Result files | 3 |
| **Inputs digest** | `1e55f71d66964578` |
| Repository version | `study-03/v0.2-handoff-amendment-02-2-g8b5b0ef` (`8b5b0ef362a7`), clean |

Quote the digest in any analysis: it identifies the data. The repository version
identifies the code — check it out (`git checkout <commit>`) to reproduce the run, or
`git diff <commit>` to see what has changed in the study since.

### Analyses of this run

**None yet.** This run has measurements but no interpretation.

### Adding an analysis

Copy `reports/analyses/TEMPLATE.md`, fill in the frontmatter, and regenerate this report.
If an analysis already exists for digest `1e55f71d66964578`, read it first and write a new one only to
**disagree, extend, or bring a different perspective**. Never edit another analyst's file.
