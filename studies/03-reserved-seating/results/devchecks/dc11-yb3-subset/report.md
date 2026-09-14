# Study 03 — results: reserved seating (venue → event → seat)

Generated from `results/dc11-yb3-subset`. Every number comes from exactly one JSON file in that
directory. **This file contains measurements and no interpretation**; conclusions live in
signed analyses, indexed at the end.

## TL;DR — measured facts

*Selected from the tables below by fixed rules; no interpretation. Read a signed analysis for
what these facts mean.*

- **Cells:** 9 across 1 topologies; 0 failed.
- **Negative controls:** 5 of 5 fired as designed to.
- ❌ **Invariant violations in designs meant to be correct:** S1 conditional on YugabyteDB, 3 nodes (RF=3) (3 early rejections); E1 sweeper on YugabyteDB, 3 nodes (RF=3) (12 early rejections); K1 window on YugabyteDB, 3 nodes (RF=3) (6 early rejections); L1 claims on YugabyteDB, 3 nodes (RF=3) (7 early rejections); L3 sharded on YugabyteDB, 3 nodes (RF=3) (7 early rejections).
- **Refused confirmations, S1 conditional, YugabyteDB, 3 nodes (RF=3)** (lifecycle, all tiers): 0 late, 0 boundary, 0 early; 123 confirmed.
- **Refused confirmations, K1 window, YugabyteDB, 3 nodes (RF=3)** (lifecycle, all tiers): 0 late, 0 boundary, 0 early; 117 confirmed.
- **Sweeper outage, YugabyteDB, 3 nodes (RF=3):** E1 6 of 6 probed seats unavailable after expiry; 0 in every other design.
- **Race, YugabyteDB, 3 nodes (RF=3)** (seats held/s; designs meant to be correct, without violations):
  - 10 seats: highest S1 conditional 61.5/s, lowest L1 claims 19.6/s
  - 100 seats: highest S1 conditional 161/s, lowest L2 document 41.3/s
  - 1k seats: highest L2 document 76.1/s, lowest L2 document 76.1/s
  - 10k seats: highest L2 document 106/s, lowest L2 document 106/s (1 of 1 timed out)
- **Publishing a 100k seats event, YugabyteDB, 3 nodes (RF=3)** (p50): 11–3704 ms with per-event inventory rows or documents, 3.88 ms without (L1).

## What was measured

|  |  |
|---|---|
| Run id | `devchecks/dc11-yb3-subset` |
| Environment | [`host-zenbook-ux5406sa`](../../../docs/environments/host-zenbook-ux5406sa.md) |
| Scale | `tiny` — 5 venues, 84 events, 5741 tickets sold, 131 live and 34 expired holds at load |
| Events (kind_seats: count) | catalogue_10: 40 · catalogue_100: 8 · catalogue_1000: 2 · catalogue_10000: 1 · lifecycle_10: 4 · lifecycle_100: 2 · lifecycle_1000: 1 · race_10: 20 · race_100: 4 · race_1000: 1 · race_10000: 1 |
| Seed | 42 — identical data in every cell |
| Reads / isolated writes | 8 workers, 2s measured + 1s warmup, 1 trial(s) |
| Race | 32 buyers per event, real hold TTL 40m0s, 10% of holds deferred until the crowd has gone, timeout 20s, tier budget 1m0s, give up after 20 conflicts |
| Lifecycle — YugabyteDB, 3 nodes (RF=3) | one human minute = 500ms; TTL 40m0s, payment window (K1) 10m0s, guard G 2m0s, sweeper every 1m0s, outage 45m0s→1h25m0s in the first event of each tier, 25% abandoned (40% of those explicitly), E0 skew 10m0s; tiers 10,100, 2 events per tier, 32 buyers |

**Engines**

| Topology | Version | READ COMMITTED runs as | Client connections to | DB clock − client clock (ms) |
|---|---|---|---|---:|
| YugabyteDB, 3 nodes (RF=3) | `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit` | read committed | 3 node(s) | 0.01 |

> Limits of this environment that bound every number below: the client and the databases
> share one laptop's eight cores; the "3-node" cluster has no real network between nodes; the
> race and the lifecycle are closed loops, so tails are a floor; the lifecycle runs on
> compressed time; E0's clock skew is injected.

## Correctness gate

Six read questions checked against Go-computed truth on a load with sold seats, live holds and
expired holds, plus an audit of the loaded state. A failure aborts the cell.

| Design | YugabyteDB, 3 nodes (RF=3) |
|---|---|
| S0 check/RC ✗ | ✅ 51/51 |
| S1 conditional | ✅ 51/51 |
| E0 app clock ✗ | ✅ 51/51 |
| E1 sweeper | ✅ 51/51 (after its sweeper released 65 expired seats) |
| K0 naive ✗ | ✅ 51/51 |
| K1 window | ✅ 51/51 |
| L1 claims | ✅ 51/51 |
| L2 document | ✅ 51/51 |
| L3 sharded | ✅ 51/51 |

## Negative controls

Deliberately wrong designs (and E1 with its sweeper stopped). Each must be seen to break the
invariant it exists to test; otherwise the absence of violations elsewhere in that experiment
is not evidence of correctness.

| Control | Experiment | Topology | Fired? | What was found |
|---|---|---|---|---|
| S0 check/RC ✗ | race | YugabyteDB, 3 nodes (RF=3) | ✅ fired | 12 deferred confirmations refused; 278 thefts, 152 early rejections |
| S0 check/RC ✗ | lifecycle | YugabyteDB, 3 nodes (RF=3) | ✅ fired | 164 thefts |
| E0 app clock ✗ | lifecycle | YugabyteDB, 3 nodes (RF=3) | ✅ fired | 73 thefts, 2 early rejections |
| K0 naive ✗ | lifecycle | YugabyteDB, 3 nodes (RF=3) | ✅ fired | 5 sales without the hold, 6 late sales |
| E1 sweeper stopped | lifecycle outage | YugabyteDB, 3 nodes (RF=3) | ✅ fired | 6 of 6 probed seats unavailable after expiry |

## The hold guarantee — lifecycle

Compressed time. *Rejected* confirmations are classified by how long before the hold's expiry
the confirming transaction started: **late** (after expiry), **boundary** (less than G before),
**early** (G or more before — a violation). Violations are counted by the ledger and the audit.
Release lag is in human minutes.

### YugabyteDB, 3 nodes (RF=3)

| Design | Tier | Confirmed seats/s | Holds | Abandoned | Expired at check | Refused late / boundary / early | Outage: unavailable / probed | Leaked | Release lag p50/p99 (min) | Idle held seat-min | Violations |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| S0 check/RC ✗ | 10 seats | ~~0.2~~ | 59 | 11+5 | 28 | 1 / 0 / 0 | 0 / 0 | 0 | 4474.4 / 6930.7 | 987 | ❌ 70 thefts |
| S0 check/RC ✗ | 100 seats | ~~1.4~~ | 216 | 39+29 | 39 | 1 / 0 / 0 | 0 / 0 | 0 | 0.1 / 0.4 | 4128 | ❌ 94 thefts |
| S1 conditional | 10 seats | 0.3 | 18 | 1+1 | 3 | 0 / 0 / 0 | 0 / 0 | 0 | 4592.8 / 6942.8 | 88 | none |
| S1 conditional | 100 seats | 2.0 | 168 | 24+23 | 11 | 0 / 0 / 0 | 0 / 0 | 0 | 0.1 / 0.9 | 2502 | none |
| E0 app clock ✗ | 10 seats | ~~0.2~~ | 31 | 7+6 | 2 | 0 / 0 / 0 | 0 / 0 | 0 | 4481.1 / 6929.2 | 554 | ❌ 13 thefts |
| E0 app clock ✗ | 100 seats | ~~1.6~~ | 158 | 26+12 | 12 | 0 / 0 / 2 | 0 / 0 | 0 | -9.9 / -9.2 | 2360 | ❌ 60 thefts, 2 early rejections |
| E1 sweeper | 10 seats | 0.2 | 21 | 4+1 | 1 | 2 / 0 / 0 | 3 / 3 | 0 | 0.8 / 11.9 | 205 | none |
| E1 sweeper | 100 seats | 1.4 | 135 | 17+10 | 9 | 1 / 0 / 0 | 3 / 3 | 0 | 0.6 / 12.6 | 1683 | none |
| K0 naive ✗ | 10 seats | 0.3 | 20 | 2+1 | 2 | 0 / 0 / 0 | 0 / 0 | 0 | 4472.4 / 6928.7 | 165 | none |
| K0 naive ✗ | 100 seats | ~~1.7~~ | 154 | 29+8 | 12 | 0 / 0 / 0 | 0 / 0 | 0 | 0.0 / 0.2 | 2215 | ❌ 5 sales without the hold, 6 late sales |
| K1 window | 10 seats | 0.2 | 20 | 5+2 | 0 | 0 / 0 / 0 | 0 / 0 | 0 | 4469.9 / 6926.1 | 381 | none |
| K1 window | 100 seats | 1.2 | 147 | 32+8 | 3 | 0 / 0 / 0 | 0 / 0 | 0 | 0.2 / 0.3 | 2674 | none |
| L1 claims | 10 seats | 0.2 | 17 | 1+1 | 2 | 0 / 0 / 0 | 0 / 0 | 0 | 4571.6 / 6921.5 | 85 | none |
| L1 claims | 100 seats | 1.8 | 160 | 28+17 | 3 | 1 / 0 / 0 | 0 / 0 | 0 | 0.2 / 0.5 | 2422 | none |
| L2 document | 10 seats | 0.3 | 15 | 0+0 | 1 | 1 / 0 / 0 | 0 / 0 | 0 | 4573.2 / 6923.2 | 0 | none |
| L2 document | 100 seats | 1.7 | 164 | 32+11 | 5 | 3 / 0 / 0 | 0 / 0 | 0 | 0.2 / 0.4 | 2552 | none |
| L3 sharded | 10 seats | 0.3 | 21 | 3+2 | 0 | 0 / 0 / 0 | 0 / 0 | 0 | 4571.5 / 6921.4 | 210 | none |
| L3 sharded | 100 seats | 1.6 | 154 | 29+11 | 5 | 2 / 0 / 0 | 0 / 0 | 0 | 0.1 / 0.1 | 2327 | none |

## The race — a hot drop with seat choice

Seats held per second, from the gate to the last hold granted. Conflicts and seat-map reads are
per granted hold. Deferred confirmations kept a real 40-minute hold through the crowd.

### YugabyteDB, 3 nodes (RF=3)

| Design | 10 seats | 100 seats | 1k seats | 10k seats |
|---|---:|---:|---:|---:|
| S0 check/RC ✗ | ~~57.8~~ ❌ 87 thefts, 52 early rejections; ❌ 13 of 23 deferred confirmations refused; 329 errors | ~~179~~ ❌ 78 thefts, 38 early rejections; ❌ 8 of 17 deferred confirmations refused; 166 errors | ~~156~~ ❌ 63 thefts, 32 early rejections; ❌ 10 of 44 deferred confirmations refused; 161 errors | ~~113~~ ❌ 50 thefts, 30 early rejections; ❌ 12 of 99 deferred confirmations refused; ⏱ 1 timed out at 22% sold; 73 errors |
| S1 conditional | 61.5 | 161 | ~~159~~ ❌ 2 early rejections | ~~116~~ ❌ 1 early rejections; ⏱ 1 timed out at 23% sold |
| E0 app clock ✗ | 60.4 | 158 | ~~144~~ ❌ 2 early rejections | ~~114~~ ❌ 10 early rejections; ⏱ 1 timed out at 23% sold |
| E1 sweeper | ~~28.5~~ ❌ 1 early rejections | 76.6 | ~~59.6~~ ❌ 6 early rejections | ~~51.7~~ ❌ 5 early rejections; ⏱ 1 timed out at 10% sold |
| K0 naive ✗ | 68.5 | 129 | 131 | 105 ⏱ 1 timed out at 21% sold |
| K1 window | 42.7 | ~~97.2~~ ❌ 1 early rejections | ~~126~~ ❌ 1 early rejections | ~~82.6~~ ❌ 4 early rejections; ⏱ 1 timed out at 17% sold |
| L1 claims | 19.6 | 93.6 | ~~156~~ ❌ 1 early rejections | ~~166~~ ❌ 6 early rejections; ⏱ 1 timed out at 33% sold |
| L2 document | 30.2 | 41.3 | 76.1 | 106 ⏱ 1 timed out at 21% sold |
| L3 sharded | 23.8 | ~~66.7~~ ❌ 1 early rejections | ~~73.2~~ ❌ 2 early rejections | ~~134~~ ❌ 4 early rejections; ⏱ 1 timed out at 27% sold |

<details><summary>Race detail — YugabyteDB, 3 nodes (RF=3)</summary>

| Design | Tier | Conflicts/hold | Map reads/hold | Engine retries | Gave up | Hold p50/p99 ms | Confirm p99 ms | Deferred ok | Final buyer sold |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|
| S0 check/RC ✗ | 10 seats | 1.09 | 5.14 | 134 | 0 | 25 / 361 | 266 | 10/23 | 0 |
| S0 check/RC ✗ | 100 seats | 0.65 | 2.92 | 58 | 0 | 26 / 196 | 309 | 9/17 | 0 |
| S0 check/RC ✗ | 1k seats | 0.45 | 2.27 | 47 | 0 | 42 / 188 | 261 | 34/44 | 0 |
| S0 check/RC ✗ | 10k seats | 0.18 | 1.28 | 13 | 0 | 88 / 201 | 209 | 87/99 | 0 |
| S1 conditional | 10 seats | 6.64 | 13.54 | 0 | 0 | 21 / 145 | 238 | 17/17 | 0 |
| S1 conditional | 100 seats | 2.40 | 4.74 | 0 | 0 | 20 / 155 | 177 | 23/23 | 0 |
| S1 conditional | 1k seats | 0.97 | 2.75 | 0 | 0 | 25 / 106 | 173 | 53/53 | 0 |
| S1 conditional | 10k seats | 0.33 | 1.36 | 0 | 0 | 75 / 172 | 204 | 91/91 | 0 |
| E0 app clock ✗ | 10 seats | 6.59 | 13.72 | 0 | 0 | 21 / 141 | 324 | 13/13 | 0 |
| E0 app clock ✗ | 100 seats | 2.62 | 4.87 | 0 | 0 | 19 / 171 | 183 | 16/16 | 0 |
| E0 app clock ✗ | 1k seats | 1.07 | 3.23 | 0 | 0 | 28 / 110 | 142 | 45/45 | 0 |
| E0 app clock ✗ | 10k seats | 0.34 | 1.40 | 0 | 0 | 76 / 174 | 207 | 109/109 | 0 |
| E1 sweeper | 10 seats | 6.60 | 14.32 | 0 | 0 | 40 / 476 | 582 | 16/16 | 0 |
| E1 sweeper | 100 seats | 2.47 | 4.91 | 0 | 0 | 50 / 407 | 445 | 21/21 | 0 |
| E1 sweeper | 1k seats | 1.08 | 3.09 | 0 | 0 | 97 / 257 | 377 | 58/58 | 0 |
| E1 sweeper | 10k seats | 0.38 | 1.38 | 0 | 0 | 128 / 365 | 670 | 40/40 | 0 |
| K0 naive ✗ | 10 seats | 6.84 | 13.77 | 0 | 0 | 20 / 137 | 356 | 14/14 | 0 |
| K0 naive ✗ | 100 seats | 2.69 | 5.10 | 0 | 0 | 23 / 212 | 194 | 20/20 | 0 |
| K0 naive ✗ | 1k seats | 1.10 | 3.00 | 0 | 0 | 33 / 134 | 181 | 45/45 | 0 |
| K0 naive ✗ | 10k seats | 0.33 | 1.35 | 0 | 0 | 81 / 183 | 213 | 85/85 | 0 |
| K1 window | 10 seats | 6.96 | 14.60 | 0 | 0 | 29 / 235 | 382 | 15/15 | 0 |
| K1 window | 100 seats | 2.55 | 4.91 | 0 | 0 | 31 / 279 | 367 | 20/20 | 0 |
| K1 window | 1k seats | 1.00 | 3.15 | 0 | 0 | 41 / 156 | 181 | 45/45 | 0 |
| K1 window | 10k seats | 0.34 | 1.35 | 0 | 0 | 89 / 253 | 671 | 70/70 | 0 |
| L1 claims | 10 seats | 6.37 | 14.09 | 4 | 0 | 30 / 899 | 240 | 17/17 | 0 |
| L1 claims | 100 seats | 2.39 | 4.63 | 1 | 0 | 28 / 508 | 182 | 16/16 | 0 |
| L1 claims | 1k seats | 0.94 | 2.41 | 0 | 0 | 38 / 146 | 143 | 62/62 | 0 |
| L1 claims | 10k seats | 0.30 | 1.42 | 0 | 0 | 47 / 204 | 361 | 127/127 | 0 |
| L2 document | 10 seats | 14.52 | 25.64 | 2960 | 0 | 44 / 268 | 479 | 7/7 | 0 |
| L2 document | 100 seats | 6.08 | 9.37 | 4772 | 0 | 95 / 515 | 2895 | 12/12 | 0 |
| L2 document | 1k seats | 2.03 | 4.06 | 5284 | 0 | 56 / 1034 | 3325 | 52/52 | 0 |
| L2 document | 10k seats | 0.94 | 2.01 | 6762 | 0 | 41 / 1001 | 4669 | 80/80 | 0 |
| L3 sharded | 10 seats | 5.03 | 10.77 | 0 | 0 | 32 / 410 | 572 | 12/12 | 0 |
| L3 sharded | 100 seats | 1.96 | 3.94 | 0 | 0 | 28 / 296 | 234 | 15/15 | 0 |
| L3 sharded | 1k seats | 0.62 | 2.61 | 0 | 0 | 83 / 189 | 686 | 39/39 | 0 |
| L3 sharded | 10k seats | 0.16 | 1.27 | 0 | 0 | 76 / 174 | 289 | 114/114 | 0 |

</details>

## Reads

Each question in isolation, fresh key per execution, ops/s.

### YugabyteDB, 3 nodes (RF=3)

| Question | S0 check/RC ✗ | S1 conditional | E0 app clock ✗ | E1 sweeper | K0 naive ✗ | K1 window | L1 claims | L2 document | L3 sharded |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| Seat map of a section — 10 seats | 5.3k | 5.6k | 5.5k | 5.3k | 4.6k | 4.6k | 5.5k | 5.6k | 4.9k |
| Seat map of a section — 100 seats | 5.4k | 5.3k | 5.3k | 5.4k | 4.5k | 3.0k | 5.9k | 6.4k | 5.3k |
| Seat map of a section — 1k seats | 2.9k | 3.0k | 3.0k | 2.9k | 2.5k | 2.6k | 4.2k | 6.2k | 4.6k |
| Seat map of a section — 10k seats | 586 | 580 | 590 | 544 | 549 | 537 | 1.1k | 5.5k | 3.6k |
| Seats available per section — 10 seats | 5.6k | 5.7k | 5.4k | 5.5k | 4.6k | 5.2k | 3.0k | 6.2k | 309 |
| Seats available per section — 100 seats | 5.1k | 5.2k | 5.0k | 4.9k | 4.4k | 4.2k | 3.0k | 5.9k | 317 |
| Seats available per section — 1k seats | 3.1k | 3.2k | 3.1k | 3.1k | 2.6k | 2.4k | 2.4k | 4.8k | 305 |
| Seats available per section — 10k seats | 367 | 378 | 375 | 364 | 311 | 319 | 555 | 2.2k | 208 |
| Seats available for an event — 10 seats | 6.1k | 6.2k | 6.1k | 6.0k | 5.4k | 3.8k | 3.2k | 6.5k | 153 |
| Seats available for an event — 100 seats | 5.4k | 5.7k | 5.0k | 5.7k | 4.8k | 4.5k | 3.2k | 6.2k | 147 |
| Seats available for an event — 1k seats | 3.0k | 3.1k | 3.0k | 3.2k | 2.4k | 2.5k | 2.9k | 4.7k | 147 |
| Seats available for an event — 10k seats | 582 | 565 | 578 | 649 | 520 | 521 | 986 | 1.9k | 153 |
| A hold's seats (basket page) | 816 | 820 | 838 | 790 | 724 | 725 | 1.4k | 5.4k | 4.3k |
| A customer's tickets | 3.8k | 3.9k | 5.1k | 3.7k | 3.5k | 3.2k | 4.9k | 4.4k | 5.0k |
| Ticket by id (control) | 7.2k | 7.5k | 7.2k | 7.2k | 5.5k | 5.8k | 6.2k | 6.7k | 7.1k |

## Isolated writes, publishing and storage

### YugabyteDB, 3 nodes (RF=3)

| Design | Hold+confirm /s | Release /s | Refund /s | Publish p50 ms by size | Storage | Audit |
|---|---:|---:|---:|---:|---:|---|
| S0 check/RC ✗ | ~~57.4~~ ❌ (4 errors) | 401 | 311 | 10 seats: 13 · 100 seats: 16 · 1k seats: 50 · 10k seats: 213 · 100k seats: 2505 | n/a on YugabyteDB | ❌ 3 thefts, 1 early rejections |
| S1 conditional | 61.6 | 460 | 347 | 10 seats: 12 · 100 seats: 16 · 1k seats: 45 · 10k seats: 224 · 100k seats: 2954 | n/a on YugabyteDB | consistent |
| E0 app clock ✗ | ~~54.5~~ ❌ | 498 | 350 | 10 seats: 13 · 100 seats: 17 · 1k seats: 40 · 10k seats: 222 · 100k seats: 3142 | n/a on YugabyteDB | ❌ 5 thefts, 1 early rejections |
| E1 sweeper | 60.2 | 747 | 342 | 10 seats: 13 · 100 seats: 24 · 1k seats: 59 · 10k seats: 297 · 100k seats: 3000 | n/a on YugabyteDB | consistent |
| K0 naive ✗ | 43.3 | 524 | 266 | 10 seats: 14 · 100 seats: 25 · 1k seats: 71 · 10k seats: 277 · 100k seats: 2948 | n/a on YugabyteDB | consistent |
| K1 window | 49.9 | 513 | 231 | 10 seats: 20 · 100 seats: 30 · 1k seats: 75 · 10k seats: 437 · 100k seats: 3704 | n/a on YugabyteDB | consistent |
| L1 claims | 67.2 | 526 | 274 | 10 seats: 6.95 · 100 seats: 6.50 · 1k seats: 6.06 · 10k seats: 3.87 · 100k seats: 3.88 | n/a on YugabyteDB | consistent |
| L2 document | 96.8 | 379 | 218 | 10 seats: 12 · 100 seats: 13 · 1k seats: 10 · 10k seats: 10 · 100k seats: 11 | n/a on YugabyteDB | consistent |
| L3 sharded | 59.4 | 549 | 357 | 10 seats: 11 · 100 seats: 13 · 1k seats: 35 · 10k seats: 134 · 100k seats: 1504 | n/a on YugabyteDB | consistent |


## Controlled pairs — one decision at a time

**This run's error bar.** Ten designs read tickets through byte-identical SQL on identical data,
so their `q05`/`q06` differences are noise: 42 comparisons, median disagreement **1.16x**, worst **1.60x**.
Read rows inside the worst disagreement are omitted below. **Race, lifecycle and write rows are always
shown**, and so is every correctness difference.

### S1 conditional → E1 sweeper

*Expiry judged at use time, or made to take effect by a sweeper?*

| Topology | Measurement | S1 conditional | E1 sweeper | Change |
|---|---|---:|---:|---|
| YugabyteDB, 3 nodes (RF=3) | race, 10 seats — seats/s | 61.5 | ~~28.5~~ ❌ 1 early rejections | 2.2x slower — ❌ a design with violations has no valid speed |
| YugabyteDB, 3 nodes (RF=3) | race, 100 seats — seats/s | 161 | 76.6 | 2.1x slower |
| YugabyteDB, 3 nodes (RF=3) | race, 1k seats — seats/s | ~~159~~ ❌ 2 early rejections | ~~59.6~~ ❌ 6 early rejections | 2.7x slower — ❌ a design with violations has no valid speed |
| YugabyteDB, 3 nodes (RF=3) | race, 10k seats — seats/s | ~~116~~ ❌ 1 early rejections; ⏱ 1 timed out at 23% sold | ~~51.7~~ ❌ 5 early rejections; ⏱ 1 timed out at 10% sold | 2.3x slower — ❌ a design with violations has no valid speed |
| YugabyteDB, 3 nodes (RF=3) | lifecycle, 10 seats — refused late/boundary/early; violations | 0/0/0; 0 | 2/0/0; 0 |  |
| YugabyteDB, 3 nodes (RF=3) | lifecycle, 100 seats — refused late/boundary/early; violations | 0/0/0; 0 | 1/0/0; 0 |  |
| YugabyteDB, 3 nodes (RF=3) | hold /s | 61.6 | 60.2 | 1.0x slower |
| YugabyteDB, 3 nodes (RF=3) | release /s | 460 | 747 | 1.6x faster |
| YugabyteDB, 3 nodes (RF=3) | cancel /s | 347 | 342 | 1.0x slower |
| YugabyteDB, 3 nodes (RF=3) | publish 10 seats p50 ms | 12 | 13 | 1.0x slower |
| YugabyteDB, 3 nodes (RF=3) | publish 100 seats p50 ms | 16 | 24 | 1.5x slower |
| YugabyteDB, 3 nodes (RF=3) | publish 1k seats p50 ms | 45 | 59 | 1.3x slower |
| YugabyteDB, 3 nodes (RF=3) | publish 10k seats p50 ms | 224 | 297 | 1.3x slower |
| YugabyteDB, 3 nodes (RF=3) | publish 100k seats p50 ms | 2954 | 3000 | 1.0x slower |

### E0 app clock ✗ → S1 conditional

*The application's clock, or the database's?*

| Topology | Measurement | E0 app clock ✗ | S1 conditional | Change |
|---|---|---:|---:|---|
| YugabyteDB, 3 nodes (RF=3) | race, 10 seats — seats/s | 60.4 | 61.5 | 1.0x faster |
| YugabyteDB, 3 nodes (RF=3) | race, 100 seats — seats/s | 158 | 161 | 1.0x faster |
| YugabyteDB, 3 nodes (RF=3) | race, 1k seats — seats/s | ~~144~~ ❌ 2 early rejections | ~~159~~ ❌ 2 early rejections | 1.1x faster — ❌ a design with violations has no valid speed |
| YugabyteDB, 3 nodes (RF=3) | race, 10k seats — seats/s | ~~114~~ ❌ 10 early rejections; ⏱ 1 timed out at 23% sold | ~~116~~ ❌ 1 early rejections; ⏱ 1 timed out at 23% sold | 1.0x faster — ❌ a design with violations has no valid speed |
| YugabyteDB, 3 nodes (RF=3) | lifecycle, 10 seats — refused late/boundary/early; violations | 0/0/0; 13 | 0/0/0; 0 |  |
| YugabyteDB, 3 nodes (RF=3) | lifecycle, 100 seats — refused late/boundary/early; violations | 0/0/2; 62 | 0/0/0; 0 |  |
| YugabyteDB, 3 nodes (RF=3) | hold /s | 54.5 | 61.6 | 1.1x faster |
| YugabyteDB, 3 nodes (RF=3) | release /s | 498 | 460 | 1.1x slower |
| YugabyteDB, 3 nodes (RF=3) | cancel /s | 350 | 347 | 1.0x slower |
| YugabyteDB, 3 nodes (RF=3) | publish 10 seats p50 ms | 13 | 12 | 1.0x faster |
| YugabyteDB, 3 nodes (RF=3) | publish 100 seats p50 ms | 17 | 16 | 1.0x faster |
| YugabyteDB, 3 nodes (RF=3) | publish 1k seats p50 ms | 40 | 45 | 1.1x slower |
| YugabyteDB, 3 nodes (RF=3) | publish 10k seats p50 ms | 222 | 224 | 1.0x slower |
| YugabyteDB, 3 nodes (RF=3) | publish 100k seats p50 ms | 3142 | 2954 | 1.1x faster |

### K0 naive ✗ → S1 conditional

*What does checking the hold at confirmation cost, and what does it prevent?*

| Topology | Measurement | K0 naive ✗ | S1 conditional | Change |
|---|---|---:|---:|---|
| YugabyteDB, 3 nodes (RF=3) | race, 10 seats — seats/s | 68.5 | 61.5 | 1.1x slower |
| YugabyteDB, 3 nodes (RF=3) | race, 100 seats — seats/s | 129 | 161 | 1.2x faster |
| YugabyteDB, 3 nodes (RF=3) | race, 1k seats — seats/s | 131 | ~~159~~ ❌ 2 early rejections | 1.2x faster — ❌ a design with violations has no valid speed |
| YugabyteDB, 3 nodes (RF=3) | race, 10k seats — seats/s | 105 ⏱ 1 timed out at 21% sold | ~~116~~ ❌ 1 early rejections; ⏱ 1 timed out at 23% sold | 1.1x faster — ❌ a design with violations has no valid speed |
| YugabyteDB, 3 nodes (RF=3) | lifecycle, 10 seats — refused late/boundary/early; violations | 0/0/0; 0 | 0/0/0; 0 |  |
| YugabyteDB, 3 nodes (RF=3) | lifecycle, 100 seats — refused late/boundary/early; violations | 0/0/0; 11 | 0/0/0; 0 |  |
| YugabyteDB, 3 nodes (RF=3) | hold /s | 43.3 | 61.6 | 1.4x faster |
| YugabyteDB, 3 nodes (RF=3) | release /s | 524 | 460 | 1.1x slower |
| YugabyteDB, 3 nodes (RF=3) | cancel /s | 266 | 347 | 1.3x faster |
| YugabyteDB, 3 nodes (RF=3) | publish 10 seats p50 ms | 14 | 12 | 1.1x faster |
| YugabyteDB, 3 nodes (RF=3) | publish 100 seats p50 ms | 25 | 16 | 1.5x faster |
| YugabyteDB, 3 nodes (RF=3) | publish 1k seats p50 ms | 71 | 45 | 1.6x faster |
| YugabyteDB, 3 nodes (RF=3) | publish 10k seats p50 ms | 277 | 224 | 1.2x faster |
| YugabyteDB, 3 nodes (RF=3) | publish 100k seats p50 ms | 2948 | 2954 | 1.0x slower |

### S1 conditional → K1 window

*What does a guaranteed payment window cost, and what does it remove?*

| Topology | Measurement | S1 conditional | K1 window | Change |
|---|---|---:|---:|---|
| YugabyteDB, 3 nodes (RF=3) | race, 10 seats — seats/s | 61.5 | 42.7 | 1.4x slower |
| YugabyteDB, 3 nodes (RF=3) | race, 100 seats — seats/s | 161 | ~~97.2~~ ❌ 1 early rejections | 1.7x slower — ❌ a design with violations has no valid speed |
| YugabyteDB, 3 nodes (RF=3) | race, 1k seats — seats/s | ~~159~~ ❌ 2 early rejections | ~~126~~ ❌ 1 early rejections | 1.3x slower — ❌ a design with violations has no valid speed |
| YugabyteDB, 3 nodes (RF=3) | race, 10k seats — seats/s | ~~116~~ ❌ 1 early rejections; ⏱ 1 timed out at 23% sold | ~~82.6~~ ❌ 4 early rejections; ⏱ 1 timed out at 17% sold | 1.4x slower — ❌ a design with violations has no valid speed |
| YugabyteDB, 3 nodes (RF=3) | lifecycle, 10 seats — refused late/boundary/early; violations | 0/0/0; 0 | 0/0/0; 0 |  |
| YugabyteDB, 3 nodes (RF=3) | lifecycle, 100 seats — refused late/boundary/early; violations | 0/0/0; 0 | 0/0/0; 0 |  |
| YugabyteDB, 3 nodes (RF=3) | hold /s | 61.6 | 49.9 | 1.2x slower |
| YugabyteDB, 3 nodes (RF=3) | release /s | 460 | 513 | 1.1x faster |
| YugabyteDB, 3 nodes (RF=3) | cancel /s | 347 | 231 | 1.5x slower |
| YugabyteDB, 3 nodes (RF=3) | publish 10 seats p50 ms | 12 | 20 | 1.6x slower |
| YugabyteDB, 3 nodes (RF=3) | publish 100 seats p50 ms | 16 | 30 | 1.8x slower |
| YugabyteDB, 3 nodes (RF=3) | publish 1k seats p50 ms | 45 | 75 | 1.7x slower |
| YugabyteDB, 3 nodes (RF=3) | publish 10k seats p50 ms | 224 | 437 | 2.0x slower |
| YugabyteDB, 3 nodes (RF=3) | publish 100k seats p50 ms | 2954 | 3704 | 1.3x slower |
| YugabyteDB, 3 nodes (RF=3) | Seats available for an event — 10 seats | 6.2k | 3.8k | 1.6x slower |
| YugabyteDB, 3 nodes (RF=3) | Seat map of a section — 100 seats | 5.3k | 3.0k | 1.8x slower |

### S1 conditional → L1 claims

*Pre-created per-event seat rows, or claims created on hold?*

| Topology | Measurement | S1 conditional | L1 claims | Change |
|---|---|---:|---:|---|
| YugabyteDB, 3 nodes (RF=3) | race, 10 seats — seats/s | 61.5 | 19.6 | 3.1x slower |
| YugabyteDB, 3 nodes (RF=3) | race, 100 seats — seats/s | 161 | 93.6 | 1.7x slower |
| YugabyteDB, 3 nodes (RF=3) | race, 1k seats — seats/s | ~~159~~ ❌ 2 early rejections | ~~156~~ ❌ 1 early rejections | 1.0x slower — ❌ a design with violations has no valid speed |
| YugabyteDB, 3 nodes (RF=3) | race, 10k seats — seats/s | ~~116~~ ❌ 1 early rejections; ⏱ 1 timed out at 23% sold | ~~166~~ ❌ 6 early rejections; ⏱ 1 timed out at 33% sold | 1.4x faster — ❌ a design with violations has no valid speed |
| YugabyteDB, 3 nodes (RF=3) | lifecycle, 10 seats — refused late/boundary/early; violations | 0/0/0; 0 | 0/0/0; 0 |  |
| YugabyteDB, 3 nodes (RF=3) | lifecycle, 100 seats — refused late/boundary/early; violations | 0/0/0; 0 | 1/0/0; 0 |  |
| YugabyteDB, 3 nodes (RF=3) | hold /s | 61.6 | 67.2 | 1.1x faster |
| YugabyteDB, 3 nodes (RF=3) | release /s | 460 | 526 | 1.1x faster |
| YugabyteDB, 3 nodes (RF=3) | cancel /s | 347 | 274 | 1.3x slower |
| YugabyteDB, 3 nodes (RF=3) | publish 10 seats p50 ms | 12 | 6.95 | 1.8x faster |
| YugabyteDB, 3 nodes (RF=3) | publish 100 seats p50 ms | 16 | 6.50 | 2.5x faster |
| YugabyteDB, 3 nodes (RF=3) | publish 1k seats p50 ms | 45 | 6.06 | 7.4x faster |
| YugabyteDB, 3 nodes (RF=3) | publish 10k seats p50 ms | 224 | 3.87 | 57.8x faster |
| YugabyteDB, 3 nodes (RF=3) | publish 100k seats p50 ms | 2954 | 3.88 | 760.9x faster |
| YugabyteDB, 3 nodes (RF=3) | Seats available per section — 10 seats | 5.7k | 3.0k | 1.9x slower |
| YugabyteDB, 3 nodes (RF=3) | Seats available for an event — 10 seats | 6.2k | 3.2k | 1.9x slower |
| YugabyteDB, 3 nodes (RF=3) | Seats available per section — 100 seats | 5.2k | 3.0k | 1.7x slower |
| YugabyteDB, 3 nodes (RF=3) | Seats available for an event — 100 seats | 5.7k | 3.2k | 1.8x slower |
| YugabyteDB, 3 nodes (RF=3) | Seat map of a section — 10k seats | 580 | 1.1k | 1.9x faster |
| YugabyteDB, 3 nodes (RF=3) | Seats available for an event — 10k seats | 565 | 986 | 1.7x faster |
| YugabyteDB, 3 nodes (RF=3) | A hold's seats (basket page) | 820 | 1.4k | 1.7x faster |

### S1 conditional → L2 document

*Seat rows, or an embedded section document?*

| Topology | Measurement | S1 conditional | L2 document | Change |
|---|---|---:|---:|---|
| YugabyteDB, 3 nodes (RF=3) | race, 10 seats — seats/s | 61.5 | 30.2 | 2.0x slower |
| YugabyteDB, 3 nodes (RF=3) | race, 100 seats — seats/s | 161 | 41.3 | 3.9x slower |
| YugabyteDB, 3 nodes (RF=3) | race, 1k seats — seats/s | ~~159~~ ❌ 2 early rejections | 76.1 | 2.1x slower — ❌ a design with violations has no valid speed |
| YugabyteDB, 3 nodes (RF=3) | race, 10k seats — seats/s | ~~116~~ ❌ 1 early rejections; ⏱ 1 timed out at 23% sold | 106 ⏱ 1 timed out at 21% sold | 1.1x slower — ❌ a design with violations has no valid speed |
| YugabyteDB, 3 nodes (RF=3) | lifecycle, 10 seats — refused late/boundary/early; violations | 0/0/0; 0 | 1/0/0; 0 |  |
| YugabyteDB, 3 nodes (RF=3) | lifecycle, 100 seats — refused late/boundary/early; violations | 0/0/0; 0 | 3/0/0; 0 |  |
| YugabyteDB, 3 nodes (RF=3) | hold /s | 61.6 | 96.8 | 1.6x faster |
| YugabyteDB, 3 nodes (RF=3) | release /s | 460 | 379 | 1.2x slower |
| YugabyteDB, 3 nodes (RF=3) | cancel /s | 347 | 218 | 1.6x slower |
| YugabyteDB, 3 nodes (RF=3) | publish 10 seats p50 ms | 12 | 12 | 1.0x faster |
| YugabyteDB, 3 nodes (RF=3) | publish 100 seats p50 ms | 16 | 13 | 1.3x faster |
| YugabyteDB, 3 nodes (RF=3) | publish 1k seats p50 ms | 45 | 10 | 4.3x faster |
| YugabyteDB, 3 nodes (RF=3) | publish 10k seats p50 ms | 224 | 10 | 22.0x faster |
| YugabyteDB, 3 nodes (RF=3) | publish 100k seats p50 ms | 2954 | 11 | 278.9x faster |
| YugabyteDB, 3 nodes (RF=3) | Seat map of a section — 1k seats | 3.0k | 6.2k | 2.1x faster |
| YugabyteDB, 3 nodes (RF=3) | Seat map of a section — 10k seats | 580 | 5.5k | 9.4x faster |
| YugabyteDB, 3 nodes (RF=3) | Seats available per section — 10k seats | 378 | 2.2k | 5.7x faster |
| YugabyteDB, 3 nodes (RF=3) | Seats available for an event — 10k seats | 565 | 1.9k | 3.4x faster |
| YugabyteDB, 3 nodes (RF=3) | A hold's seats (basket page) | 820 | 5.4k | 6.6x faster |

### S1 conditional → L3 sharded

*(YugabyteDB) all of an event's seats on one tablet, or spread by section?*

| Topology | Measurement | S1 conditional | L3 sharded | Change |
|---|---|---:|---:|---|
| YugabyteDB, 3 nodes (RF=3) | race, 10 seats — seats/s | 61.5 | 23.8 | 2.6x slower |
| YugabyteDB, 3 nodes (RF=3) | race, 100 seats — seats/s | 161 | ~~66.7~~ ❌ 1 early rejections | 2.4x slower — ❌ a design with violations has no valid speed |
| YugabyteDB, 3 nodes (RF=3) | race, 1k seats — seats/s | ~~159~~ ❌ 2 early rejections | ~~73.2~~ ❌ 2 early rejections | 2.2x slower — ❌ a design with violations has no valid speed |
| YugabyteDB, 3 nodes (RF=3) | race, 10k seats — seats/s | ~~116~~ ❌ 1 early rejections; ⏱ 1 timed out at 23% sold | ~~134~~ ❌ 4 early rejections; ⏱ 1 timed out at 27% sold | 1.2x faster — ❌ a design with violations has no valid speed |
| YugabyteDB, 3 nodes (RF=3) | lifecycle, 10 seats — refused late/boundary/early; violations | 0/0/0; 0 | 0/0/0; 0 |  |
| YugabyteDB, 3 nodes (RF=3) | lifecycle, 100 seats — refused late/boundary/early; violations | 0/0/0; 0 | 2/0/0; 0 |  |
| YugabyteDB, 3 nodes (RF=3) | hold /s | 61.6 | 59.4 | 1.0x slower |
| YugabyteDB, 3 nodes (RF=3) | release /s | 460 | 549 | 1.2x faster |
| YugabyteDB, 3 nodes (RF=3) | cancel /s | 347 | 357 | 1.0x faster |
| YugabyteDB, 3 nodes (RF=3) | publish 10 seats p50 ms | 12 | 11 | 1.2x faster |
| YugabyteDB, 3 nodes (RF=3) | publish 100 seats p50 ms | 16 | 13 | 1.2x faster |
| YugabyteDB, 3 nodes (RF=3) | publish 1k seats p50 ms | 45 | 35 | 1.3x faster |
| YugabyteDB, 3 nodes (RF=3) | publish 10k seats p50 ms | 224 | 134 | 1.7x faster |
| YugabyteDB, 3 nodes (RF=3) | publish 100k seats p50 ms | 2954 | 1504 | 2.0x faster |
| YugabyteDB, 3 nodes (RF=3) | Seats available per section — 10 seats | 5.7k | 309 | 18.6x slower |
| YugabyteDB, 3 nodes (RF=3) | Seats available for an event — 10 seats | 6.2k | 153 | 40.4x slower |
| YugabyteDB, 3 nodes (RF=3) | Seats available per section — 100 seats | 5.2k | 317 | 16.2x slower |
| YugabyteDB, 3 nodes (RF=3) | Seats available for an event — 100 seats | 5.7k | 147 | 38.7x slower |
| YugabyteDB, 3 nodes (RF=3) | Seats available per section — 1k seats | 3.2k | 305 | 10.5x slower |
| YugabyteDB, 3 nodes (RF=3) | Seats available for an event — 1k seats | 3.1k | 147 | 21.0x slower |
| YugabyteDB, 3 nodes (RF=3) | Seat map of a section — 10k seats | 580 | 3.6k | 6.3x faster |
| YugabyteDB, 3 nodes (RF=3) | Seats available per section — 10k seats | 378 | 208 | 1.8x slower |
| YugabyteDB, 3 nodes (RF=3) | Seats available for an event — 10k seats | 565 | 153 | 3.7x slower |
| YugabyteDB, 3 nodes (RF=3) | A hold's seats (basket page) | 820 | 4.3k | 5.2x faster |


## Measurement hazard: CPU-quota throttling

Cells show *throttled periods / periods (throttled time)*: for the **database** container(s) over
the whole cell (summed across nodes), and for the **client** per phase.

### YugabyteDB, 3 nodes (RF=3)

| Design | database, whole cell | client lifecycle#1 | client race#1 | client read | client write:cancel | client write:hold | client write:publish | client write:release |
|---|---:|---:|---:|---:|---:|---:|---:|---:|
| S0 check/RC ✗ | 1203/14609 (213s) | 0/2224 (0.0s) | 0/387 (0.0s) | 0/453 (0.0s) | 0/65 (0.0s) | 0/127 (0.0s) | 0/99 (0.0s) | 0/2 (0.0s) |
| S1 conditional | 1213/13047 (220s) | 0/1687 (0.0s) | 0/365 (0.0s) | 0/452 (0.0s) | 0/62 (0.0s) | 0/126 (0.0s) | 0/98 (0.0s) | 0/2 (0.0s) |
| E0 app clock ✗ | 1247/14519 (224s) | 0/2244 (0.0s) | 0/371 (0.0s) | 0/451 (0.0s) | 0/63 (0.0s) | 0/137 (0.0s) | 0/102 (0.0s) | 0/1 (0.0s) |
| E1 sweeper | 1538/16062 (272s) | 0/2406 (0.0s) | 0/575 (0.0s) | 0/452 (0.0s) | 0/61 (0.0s) | 0/137 (0.0s) | 0/108 (0.0s) | 0/1 (0.0s) |
| K0 naive ✗ | 1309/13712 (233s) | 0/1857 (0.0s) | 0/383 (0.0s) | 0/453 (0.0s) | 0/69 (0.0s) | 0/153 (0.0s) | 0/111 (0.0s) | 0/2 (0.0s) |
| K1 window | 1425/16721 (223s) | 0/2726 (0.0s) | 0/421 (0.0s) | 0/452 (0.0s) | 0/72 (0.0s) | 0/157 (0.0s) | 0/143 (0.0s) | 0/3 (0.0s) |
| L1 claims | 1309/12600 (156s) | 0/1810 (0.0s) | 0/459 (0.0s) | 0/452 (0.0s) | 0/72 (0.0s) | 0/124 (0.0s) | 0/58 (0.0s) | 0/1 (0.0s) |
| L2 document | 1268/13084 (112s) | 0/1833 (0.0s) | 0/568 (0.0s) | 0/451 (0.0s) | 0/89 (0.0s) | 0/90 (0.0s) | 0/85 (0.0s) | 0/2 (0.0s) |
| L3 sharded | 1434/13737 (358s) | 0/1911 (0.0s) | 0/618 (0.0s) | 0/455 (0.0s) | 0/64 (0.0s) | 0/144 (0.0s) | 0/87 (0.0s) | 0/2 (0.0s) |


---

Plans for every read and write statement are in `results/dc11-yb3-subset/<topology>/plans/<design>.txt`
(`EXPLAIN (ANALYZE, BUFFERS)` on PostgreSQL, `EXPLAIN (ANALYZE, DIST)` on YugabyteDB; write
statements inside rolled-back transactions).

## Provenance: code, data and analyses

This file is **generated from the measurements** and contains no interpretation.
Conclusions live in separate signed analyses, so several analysts — including
different AI models — can read the same numbers and each record what they make of
them. Where two analyses disagree, the disagreement is a finding and stays visible.

| | |
|---|---|
| Run id | `devchecks/dc11-yb3-subset` |
| Result files | 9 |
| **Inputs digest** | `0b563d22f8118d0f` |
| Repository version | `repo/platform-lock-not-available-7-g31c467b-dirty` (`31c467b5d172`), ⚠️ **dirty** — uncommitted changes; not reproducible from the commit alone |

Quote the digest in any analysis: it identifies the data. The repository version
identifies the code — check it out (`git checkout <commit>`) to reproduce the run, or
`git diff <commit>` to see what has changed in the study since.

### Analyses of this run

**None yet.** This run has measurements but no interpretation.

### Adding an analysis

Copy `reports/analyses/TEMPLATE.md`, fill in the frontmatter, and regenerate this report.
If an analysis already exists for digest `0b563d22f8118d0f`, read it first and write a new one only to
**disagree, extend, or bring a different perspective**. Never edit another analyst's file.
