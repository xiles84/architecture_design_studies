# Study 03 — results: reserved seating (venue → event → seat)

Generated from `results/dc5-yb1-all`. Every number comes from exactly one JSON file in that
directory. **This file contains measurements and no interpretation**; conclusions live in
signed analyses, indexed at the end.

## TL;DR — measured facts

*Selected from the tables below by fixed rules; no interpretation. Read a signed analysis for
what these facts mean.*

- **Cells:** 13 across 1 topologies; 0 failed.
- **Negative controls:** 5 of 5 fired as designed to.
- ❌ **Invariant violations in designs meant to be correct:** S1 conditional on YugabyteDB, 1 node (RF=1) (1 early rejections); S2 lock on YugabyteDB, 1 node (RF=1) (1 early rejections); S3 nowait on YugabyteDB, 1 node (RF=1) (1 early rejections); K1 window on YugabyteDB, 1 node (RF=1) (1 early rejections).
- **Refused confirmations, S1 conditional, YugabyteDB, 1 node (RF=1)** (lifecycle, all tiers): 1 late, 0 boundary, 0 early; 121 confirmed.
- **Refused confirmations, K1 window, YugabyteDB, 1 node (RF=1)** (lifecycle, all tiers): 0 late, 0 boundary, 0 early; 119 confirmed.
- **Sweeper outage, YugabyteDB, 1 node (RF=1):** E1 10 of 10 probed seats unavailable after expiry; 0 in every other design.
- **Race, YugabyteDB, 1 node (RF=1)** (seats held/s; designs meant to be correct, without violations):
  - 10 seats: highest S2 lock 61.7/s, lowest S4 check/SER 0.1/s (1 of 10 timed out)
  - 100 seats: highest S1 conditional 126/s, lowest S4 check/SER 0.1/s (1 of 10 timed out)
  - 1k seats: highest L1 claims 141/s, lowest S4 check/SER 0.1/s (1 of 9 timed out)
  - 10k seats: highest L1 claims 154/s, lowest S4 check/SER 0.3/s (7 of 7 timed out)
- **Publishing a 100k seats event, YugabyteDB, 1 node (RF=1)** (p50): 4.34–2948 ms with per-event inventory rows or documents, 2.54 ms without (L1).

## What was measured

|  |  |
|---|---|
| Run id | `devchecks/dc5-yb1-all` |
| Environment | [`host-zenbook-ux5406sa`](../../../docs/environments/host-zenbook-ux5406sa.md) |
| Scale | `tiny` — 5 venues, 84 events, 5741 tickets sold, 131 live and 34 expired holds at load |
| Events (kind_seats: count) | catalogue_10: 40 · catalogue_100: 8 · catalogue_1000: 2 · catalogue_10000: 1 · lifecycle_10: 4 · lifecycle_100: 2 · lifecycle_1000: 1 · race_10: 20 · race_100: 4 · race_1000: 1 · race_10000: 1 |
| Seed | 42 — identical data in every cell |
| Reads / isolated writes | 8 workers, 2s measured + 1s warmup, 1 trial(s) |
| Race | 32 buyers per event, real hold TTL 40m0s, 10% of holds deferred until the crowd has gone, timeout 20s, tier budget 1m0s, give up after 20 conflicts |
| Lifecycle — YugabyteDB, 1 node (RF=1) | one human minute = 500ms; TTL 40m0s, payment window (K1) 10m0s, guard G 2m0s, sweeper every 1m0s, outage 45m0s→1h25m0s in the first event of each tier, 25% abandoned (40% of those explicitly), E0 skew 10m0s; tiers 10,100, 2 events per tier, 32 buyers |

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
| S0 check/RC ✗ | ✅ 51/51 |
| S1 conditional | ✅ 51/51 |
| S2 lock | ✅ 51/51 |
| S3 nowait | ✅ 51/51 |
| S4 check/SER | ✅ 51/51 |
| E0 app clock ✗ | ✅ 51/51 |
| E1 sweeper | ✅ 51/51 (after its sweeper released 65 expired seats) |
| E2 cart | ✅ 51/51 |
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
| S0 check/RC ✗ | race | YugabyteDB, 1 node (RF=1) | ✅ fired | 9 deferred confirmations refused; 406 thefts, 219 early rejections |
| S0 check/RC ✗ | lifecycle | YugabyteDB, 1 node (RF=1) | ✅ fired | 172 thefts |
| E0 app clock ✗ | lifecycle | YugabyteDB, 1 node (RF=1) | ✅ fired | 56 thefts, 1 early rejections |
| K0 naive ✗ | lifecycle | YugabyteDB, 1 node (RF=1) | ✅ fired | 3 sales without the hold, 3 late sales |
| E1 sweeper stopped | lifecycle outage | YugabyteDB, 1 node (RF=1) | ✅ fired | 10 of 10 probed seats unavailable after expiry |

## The hold guarantee — lifecycle

Compressed time. *Rejected* confirmations are classified by how long before the hold's expiry
the confirming transaction started: **late** (after expiry), **boundary** (less than G before),
**early** (G or more before — a violation). Violations are counted by the ledger and the audit.
Release lag is in human minutes.

### YugabyteDB, 1 node (RF=1)

| Design | Tier | Confirmed seats/s | Holds | Abandoned | Expired at check | Refused late / boundary / early | Outage: unavailable / probed | Leaked | Release lag p50/p99 (min) | Idle held seat-min | Violations |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| S0 check/RC ✗ | 10 seats | ~~0.2~~ | 57 | 12+6 | 27 | 0 / 0 / 0 | 0 / 0 | 0 | 4470.0 / 6918.1 | 1019 | ❌ 63 thefts |
| S0 check/RC ✗ | 100 seats | ~~1.6~~ | 222 | 31+30 | 51 | 0 / 0 / 0 | 0 / 0 | 0 | 0.0 / 0.8 | 3202 | ❌ 109 thefts |
| S1 conditional | 10 seats | 0.3 | 21 | 3+2 | 1 | 0 / 0 / 0 | 0 / 0 | 0 | 4461.5 / 6909.6 | 212 | none |
| S1 conditional | 100 seats | 1.9 | 152 | 26+14 | 5 | 1 / 0 / 0 | 0 / 0 | 0 | 0.1 / 2.5 | 2151 | none |
| S2 lock | 10 seats | 0.2 | 23 | 5+3 | 1 | 0 / 0 / 0 | 0 / 0 | 0 | 4453.6 / 6909.9 | 429 | none |
| S2 lock | 100 seats | 1.7 | 163 | 29+17 | 4 | 2 / 0 / 0 | 0 / 0 | 0 | 0.2 / 0.7 | 2964 | none |
| S3 nowait | 10 seats | 0.1 | 23 | 2+1 | 4 | 0 / 0 / 0 | 0 / 0 | 0 | 4456.1 / 6912.4 | 188 | none |
| S3 nowait | 100 seats | 1.5 | 160 | 25+18 | 6 | 2 / 0 / 0 | 0 / 0 | 0 | 0.2 / 1.2 | 2330 | none |
| S4 check/SER | 10 seats | 0.0 | 4 | 0+0 | 1 | 2 / 0 / 0 | 0 / 3 | 0 | 4463.5 / 7219.2 | 0 | none |
| S4 check/SER | 100 seats | 0.0 | 16 | 1+1 | 2 | 5 / 0 / 0 | 0 / 8 | 0 | 57.0 / 188.9 | 175 | none |
| E0 app clock ✗ | 10 seats | ~~0.2~~ | 20 | 4+3 | 0 | 0 / 0 / 0 | 0 / 0 | 0 | 4469.2 / 6917.3 | 368 | ❌ 8 thefts |
| E0 app clock ✗ | 100 seats | ~~1.7~~ | 157 | 28+12 | 11 | 0 / 0 / 1 | 0 / 0 | 0 | -9.6 / -9.2 | 2244 | ❌ 48 thefts, 1 early rejections |
| E1 sweeper | 10 seats | 0.3 | 20 | 2+2 | 0 | 0 / 0 / 0 | 0 / 0 | 0 | 5.2 / 5.2 | 132 | none |
| E1 sweeper | 100 seats | 1.6 | 134 | 17+11 | 1 | 3 / 0 / 0 | 10 / 10 | 0 | 0.8 / 48.9 | 1913 | none |
| E2 cart | 10 seats | 0.1 | 32 | 8+3 | 2 | 1 / 0 / 0 | 0 / 0 | 0 | 4456.1 / 6912.4 | 520 | none |
| E2 cart | 100 seats | 1.7 | 153 | 27+11 | 4 | 3 / 0 / 0 | 0 / 0 | 0 | 0.3 / 0.5 | 2065 | none |
| K0 naive ✗ | 10 seats | ~~0.2~~ | 22 | 3+3 | 2 | 0 / 0 / 0 | 0 / 0 | 0 | 4559.7 / 6909.6 | 339 | ❌ 1 sales without the hold, 1 late sales |
| K0 naive ✗ | 100 seats | ~~1.3~~ | 168 | 29+18 | 12 | 0 / 0 / 0 | 0 / 0 | 0 | 0.1 / 0.2 | 2395 | ❌ 2 sales without the hold, 2 late sales |
| K1 window | 10 seats | 0.2 | 15 | 3+1 | 0 | 0 / 0 / 0 | 0 / 0 | 0 | 4461.0 / 6909.0 | 128 | none |
| K1 window | 100 seats | 1.6 | 156 | 25+19 | 4 | 0 / 0 / 0 | 0 / 0 | 0 | 0.0 / 0.0 | 2523 | none |
| L1 claims | 10 seats | 0.2 | 21 | 5+1 | 1 | 0 / 0 / 0 | 0 / 0 | 0 | 4454.1 / 6910.4 | 407 | none |
| L1 claims | 100 seats | 1.9 | 146 | 24+13 | 7 | 0 / 0 / 0 | 0 / 0 | 0 | 0.0 / 0.2 | 2129 | none |
| L2 document | 10 seats | 0.2 | 17 | 3+0 | 1 | 0 / 0 / 0 | 0 / 0 | 0 | 4460.5 / 6908.5 | 200 | none |
| L2 document | 100 seats | 1.9 | 142 | 21+12 | 7 | 2 / 0 / 0 | 0 / 0 | 0 | 0.1 / 0.1 | 1564 | none |
| L3 sharded | 10 seats | 0.2 | 24 | 6+3 | 1 | 0 / 0 / 0 | 0 / 0 | 0 | 4452.3 / 6908.6 | 462 | none |
| L3 sharded | 100 seats | 2.3 | 148 | 24+16 | 2 | 0 / 0 / 0 | 0 / 0 | 0 | 0.3 / 0.6 | 2215 | none |

## The race — a hot drop with seat choice

Seats held per second, from the gate to the last hold granted. Conflicts and seat-map reads are
per granted hold. Deferred confirmations kept a real 40-minute hold through the crowd.

### YugabyteDB, 1 node (RF=1)

| Design | 10 seats | 100 seats | 1k seats | 10k seats |
|---|---:|---:|---:|---:|
| S0 check/RC ✗ | ~~61.7~~ ❌ 124 thefts, 71 early rejections; ❌ 12 of 22 deferred confirmations refused; 277 errors | ~~151~~ ❌ 120 thefts, 58 early rejections; ❌ 17 of 26 deferred confirmations refused; 109 errors | ~~161~~ ❌ 100 thefts, 55 early rejections; ❌ 17 of 58 deferred confirmations refused; 121 errors | ~~96.2~~ ❌ 62 thefts, 35 early rejections; ❌ 9 of 97 deferred confirmations refused; ⏱ 1 timed out at 19% sold; 49 errors |
| S1 conditional | 50.3 | 126 | 133 | ~~90.6~~ ❌ 1 early rejections; ⏱ 1 timed out at 18% sold |
| S2 lock | 61.7 | 98.2 | 120 | ~~85.0~~ ❌ 1 early rejections; ⏱ 1 timed out at 17% sold |
| S3 nowait | 20.7 | 81.5 | ~~100~~ ❌ 1 early rejections | 81.2 ⏱ 1 timed out at 17% sold |
| S4 check/SER | 0.1 ⏱ 2 timed out at 45% sold; ⌛ 2 of 20 events raced within the tier budget | 0.1 ⏱ 3 timed out at 2% sold; ⌛ 3 of 4 events raced within the tier budget | 0.1 ⏱ 1 timed out at 0% sold | 0.3 ⏱ 1 timed out at 0% sold |
| E0 app clock ✗ | 50.2 | 113 | 140 | ~~79.9~~ ❌ 1 early rejections; ⏱ 1 timed out at 16% sold |
| E1 sweeper | 45.2 | 117 | 118 | 77.0 ⏱ 1 timed out at 16% sold |
| E2 cart | 15.3 | 50.5 | 75.5 | 63.6 ⏱ 1 timed out at 13% sold |
| K0 naive ✗ | 46.8 | 135 | 112 | 87.0 ⏱ 1 timed out at 18% sold |
| K1 window | 39.9 | 113 | 114 | ~~83.4~~ ❌ 1 early rejections; ⏱ 1 timed out at 17% sold |
| L1 claims | 26.5 | 87.1 | 141 | 154 ⏱ 1 timed out at 31% sold |
| L2 document | 33.8 | 60.3 | 82.4 | 143 ⏱ 1 timed out at 29% sold |
| L3 sharded | 15.5 | 46.8 | 68.8 | 105 ⏱ 1 timed out at 21% sold |

<details><summary>Race detail — YugabyteDB, 1 node (RF=1)</summary>

| Design | Tier | Conflicts/hold | Map reads/hold | Engine retries | Gave up | Hold p50/p99 ms | Confirm p99 ms | Deferred ok | Final buyer sold |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|
| S0 check/RC ✗ | 10 seats | 1.46 | 5.28 | 133 | 0 | 66 / 367 | 281 | 10/22 | 0 |
| S0 check/RC ✗ | 100 seats | 0.84 | 2.68 | 57 | 0 | 76 / 282 | 302 | 9/26 | 0 |
| S0 check/RC ✗ | 1k seats | 0.33 | 1.85 | 15 | 0 | 89 / 194 | 186 | 41/58 | 0 |
| S0 check/RC ✗ | 10k seats | 0.16 | 1.25 | 8 | 0 | 99 / 212 | 284 | 88/97 | 0 |
| S1 conditional | 10 seats | 6.46 | 13.68 | 0 | 0 | 12 / 204 | 309 | 15/15 | 0 |
| S1 conditional | 100 seats | 2.17 | 4.30 | 0 | 0 | 13 / 196 | 379 | 13/13 | 0 |
| S1 conditional | 1k seats | 0.92 | 3.17 | 0 | 0 | 77 / 181 | 189 | 47/47 | 0 |
| S1 conditional | 10k seats | 0.33 | 1.39 | 0 | 0 | 95 / 200 | 586 | 85/85 | 0 |
| S2 lock | 10 seats | 5.41 | 11.56 | 0 | 0 | 12 / 173 | 387 | 11/11 | 0 |
| S2 lock | 100 seats | 1.90 | 4.74 | 0 | 0 | 17 / 198 | 411 | 18/18 | 0 |
| S2 lock | 1k seats | 1.02 | 3.28 | 0 | 0 | 91 / 194 | 389 | 47/47 | 0 |
| S2 lock | 10k seats | 0.31 | 1.34 | 0 | 0 | 109 / 295 | 301 | 65/65 | 0 |
| S3 nowait | 10 seats | 21.47 | 29.46 | 0 | 0 | 83 / 722 | 500 | 15/15 | 0 |
| S3 nowait | 100 seats | 3.60 | 6.37 | 0 | 0 | 93 / 180 | 503 | 19/19 | 0 |
| S3 nowait | 1k seats | 1.36 | 3.76 | 0 | 0 | 97 / 203 | 288 | 45/45 | 0 |
| S3 nowait | 10k seats | 0.52 | 1.56 | 0 | 0 | 111 / 305 | 390 | 69/69 | 0 |
| S4 check/SER | 10 seats | 3.75 | 23.50 | 568 | 0 | 999 / 25058 | 78 | 2/2 | 0 |
| S4 check/SER | 100 seats | 0.33 | 32.00 | 1061 | 0 | 20915 / 20915 | 30 | 0/0 | 0 |
| S4 check/SER | 1k seats | 0.00 | 32.00 | 405 | 0 | 38307 / 38307 | 93 | 0/0 | 0 |
| S4 check/SER | 10k seats | 0.00 | 16.00 | 414 | 0 | 23008 / 23008 | 23 | 0/0 | 0 |
| E0 app clock ✗ | 10 seats | 6.25 | 13.20 | 0 | 0 | 10 / 280 | 379 | 17/17 | 0 |
| E0 app clock ✗ | 100 seats | 2.26 | 4.17 | 0 | 0 | 19 / 280 | 410 | 23/23 | 0 |
| E0 app clock ✗ | 1k seats | 1.00 | 2.56 | 0 | 0 | 84 / 181 | 193 | 46/46 | 0 |
| E0 app clock ✗ | 10k seats | 0.35 | 1.38 | 0 | 0 | 99 / 212 | 702 | 66/66 | 0 |
| E1 sweeper | 10 seats | 6.31 | 13.52 | 0 | 0 | 10 / 210 | 396 | 15/15 | 0 |
| E1 sweeper | 100 seats | 2.54 | 4.75 | 0 | 0 | 16 / 205 | 187 | 23/23 | 0 |
| E1 sweeper | 1k seats | 0.96 | 2.50 | 0 | 0 | 86 / 197 | 508 | 41/41 | 0 |
| E1 sweeper | 10k seats | 0.31 | 1.32 | 0 | 0 | 99 / 206 | 297 | 54/54 | 0 |
| E2 cart | 10 seats | 7.47 | 16.43 | 0 | 0 | 72 / 705 | 491 | 13/13 | 0 |
| E2 cart | 100 seats | 2.04 | 4.49 | 0 | 0 | 90 / 609 | 196 | 24/24 | 0 |
| E2 cart | 1k seats | 0.98 | 2.80 | 0 | 0 | 95 / 233 | 211 | 49/49 | 0 |
| E2 cart | 10k seats | 0.29 | 1.29 | 0 | 0 | 111 / 300 | 900 | 58/58 | 0 |
| K0 naive ✗ | 10 seats | 6.90 | 14.48 | 0 | 0 | 24 / 288 | 407 | 10/10 | 0 |
| K0 naive ✗ | 100 seats | 2.46 | 4.83 | 0 | 0 | 17 / 179 | 189 | 20/20 | 0 |
| K0 naive ✗ | 1k seats | 0.98 | 3.35 | 0 | 0 | 78 / 182 | 197 | 50/50 | 0 |
| K0 naive ✗ | 10k seats | 0.32 | 1.36 | 0 | 0 | 95 / 196 | 284 | 69/69 | 0 |
| K1 window | 10 seats | 6.77 | 14.55 | 0 | 0 | 13 / 379 | 406 | 11/11 | 0 |
| K1 window | 100 seats | 2.20 | 4.74 | 0 | 0 | 10 / 251 | 405 | 15/15 | 0 |
| K1 window | 1k seats | 0.89 | 2.70 | 0 | 0 | 87 / 190 | 197 | 36/36 | 0 |
| K1 window | 10k seats | 0.35 | 1.38 | 0 | 0 | 96 / 212 | 500 | 80/80 | 0 |
| L1 claims | 10 seats | 6.82 | 13.63 | 0 | 0 | 43 / 829 | 383 | 18/18 | 0 |
| L1 claims | 100 seats | 2.26 | 4.44 | 0 | 0 | 22 / 294 | 199 | 13/13 | 0 |
| L1 claims | 1k seats | 0.82 | 2.74 | 0 | 0 | 79 / 113 | 190 | 49/49 | 0 |
| L1 claims | 10k seats | 0.26 | 1.34 | 1 | 0 | 90 / 182 | 201 | 131/131 | 0 |
| L2 document | 10 seats | 13.47 | 26.05 | 2260 | 0 | 31 / 257 | 578 | 14/14 | 0 |
| L2 document | 100 seats | 4.68 | 7.05 | 3069 | 0 | 44 / 402 | 2108 | 17/17 | 0 |
| L2 document | 1k seats | 1.55 | 4.10 | 3989 | 0 | 64 / 878 | 5015 | 45/45 | 0 |
| L2 document | 10k seats | 0.48 | 1.53 | 5937 | 0 | 14 / 483 | 6907 | 117/117 | 0 |
| L3 sharded | 10 seats | 4.76 | 10.65 | 0 | 0 | 87 / 793 | 1180 | 14/14 | 0 |
| L3 sharded | 100 seats | 1.67 | 3.47 | 0 | 0 | 91 / 298 | 408 | 15/15 | 0 |
| L3 sharded | 1k seats | 0.55 | 2.17 | 0 | 0 | 99 / 300 | 316 | 43/43 | 0 |
| L3 sharded | 10k seats | 0.18 | 1.24 | 0 | 0 | 98 / 283 | 299 | 86/86 | 0 |

</details>

## Reads

Each question in isolation, fresh key per execution, ops/s.

### YugabyteDB, 1 node (RF=1)

| Question | S0 check/RC ✗ | S1 conditional | S2 lock | S3 nowait | S4 check/SER | E0 app clock ✗ | E1 sweeper | E2 cart | K0 naive ✗ | K1 window | L1 claims | L2 document | L3 sharded |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| Seat map of a section — 10 seats | 3.9k | 3.5k | 3.7k | 3.5k | 3.2k | 3.2k | 3.9k | 2.6k | 3.2k | 3.8k | 3.9k | 3.9k | 3.8k |
| Seat map of a section — 100 seats | 3.7k | 3.6k | 3.5k | 3.0k | 3.1k | 3.5k | 3.7k | 1.5k | 3.6k | 3.1k | 3.6k | 3.8k | 3.5k |
| Seat map of a section — 1k seats | 2.2k | 2.1k | 1.8k | 2.1k | 2.2k | 2.0k | 1.8k | 1.1k | 2.1k | 1.9k | 2.9k | 3.3k | 2.8k |
| Seat map of a section — 10k seats | 487 | 516 | 404 | 381 | 465 | 447 | 344 | 351 | 485 | 373 | 887 | 3.1k | 2.5k |
| Seats available per section — 10 seats | 3.1k | 2.8k | 2.9k | 3.2k | 3.2k | 2.8k | 2.8k | 1.7k | 3.2k | 3.1k | 1.5k | 4.0k | 179 |
| Seats available per section — 100 seats | 3.1k | 2.9k | 3.0k | 2.6k | 3.0k | 2.7k | 2.7k | 1.5k | 3.0k | 3.0k | 1.4k | 3.9k | 235 |
| Seats available per section — 1k seats | 1.9k | 1.9k | 1.6k | 1.5k | 1.5k | 1.5k | 2.0k | 1.1k | 1.9k | 1.6k | 1.2k | 2.9k | 179 |
| Seats available per section — 10k seats | 224 | 242 | 220 | 181 | 214 | 180 | 196 | 175 | 216 | 193 | 379 | 1.0k | 127 |
| Seats available for an event — 10 seats | 4.1k | 3.4k | 3.1k | 3.6k | 4.0k | 4.0k | 4.0k | 1.6k | 4.0k | 3.7k | 1.7k | 3.3k | 85.7 |
| Seats available for an event — 100 seats | 3.9k | 3.7k | 3.7k | 3.4k | 3.7k | 3.5k | 3.9k | 1.6k | 3.3k | 3.4k | 1.8k | 3.7k | 93.7 |
| Seats available for an event — 1k seats | 2.2k | 2.3k | 2.3k | 1.9k | 2.0k | 1.7k | 2.3k | 1.3k | 2.3k | 2.3k | 1.6k | 3.0k | 104 |
| Seats available for an event — 10k seats | 451 | 483 | 404 | 365 | 510 | 364 | 441 | 199 | 525 | 405 | 773 | 978 | 105 |
| A hold's seats (basket page) | 610 | 662 | 520 | 705 | 662 | 502 | 558 | 543 | 651 | 510 | 1.1k | 3.8k | 3.0k |
| A customer's tickets | 2.9k | 2.8k | 2.8k | 2.6k | 2.4k | 2.7k | 2.7k | 2.7k | 2.8k | 2.7k | 2.8k | 2.8k | 2.9k |
| Ticket by id (control) | 4.9k | 4.9k | 4.9k | 4.8k | 4.8k | 4.8k | 4.9k | 4.9k | 4.8k | 4.8k | 4.6k | 4.8k | 4.9k |

## Isolated writes, publishing and storage

### YugabyteDB, 1 node (RF=1)

| Design | Hold+confirm /s | Release /s | Refund /s | Publish p50 ms by size | Storage | Audit |
|---|---:|---:|---:|---:|---:|---|
| S0 check/RC ✗ | ~~60.8~~ ❌ (1 errors) | 386 | 338 | 10 seats: 4.79 · 100 seats: 8.14 · 1k seats: 89 · 10k seats: 263 · 100k seats: 2847 | n/a on YugabyteDB | ❌ 6 thefts, 3 early rejections |
| S1 conditional | 53.2 | 378 | 345 | 10 seats: 4.84 · 100 seats: 8.67 · 1k seats: 98 · 10k seats: 395 · 100k seats: 2810 | n/a on YugabyteDB | consistent |
| S2 lock | 45.3 | 1.9k | 278 | 10 seats: 4.35 · 100 seats: 7.12 · 1k seats: 89 · 10k seats: 236 · 100k seats: 2916 | n/a on YugabyteDB | consistent |
| S3 nowait | 48.2 | 2.1k | 279 | 10 seats: 4.71 · 100 seats: 8.64 · 1k seats: 93 · 10k seats: 259 · 100k seats: 2780 | n/a on YugabyteDB | consistent |
| S4 check/SER | 1.7 | 1.9k | 329 | 10 seats: 5.02 · 100 seats: 7.05 · 1k seats: 87 · 10k seats: 213 · 100k seats: 2887 | n/a on YugabyteDB | consistent |
| E0 app clock ✗ | ~~46.8~~ ❌ | 370 | 261 | 10 seats: 4.75 · 100 seats: 7.04 · 1k seats: 96 · 10k seats: 234 · 100k seats: 2827 | n/a on YugabyteDB | ❌ 3 thefts |
| E1 sweeper | 48.2 | 2.2k | 280 | 10 seats: 4.52 · 100 seats: 7.93 · 1k seats: 102 · 10k seats: 296 · 100k seats: 2837 | n/a on YugabyteDB | consistent |
| E2 cart | 42.5 | 319 | 295 | 10 seats: 4.63 · 100 seats: 7.48 · 1k seats: 93 · 10k seats: 239 · 100k seats: 2879 | n/a on YugabyteDB | consistent |
| K0 naive ✗ | 51.9 | 368 | 359 | 10 seats: 4.63 · 100 seats: 8.88 · 1k seats: 89 · 10k seats: 216 · 100k seats: 2745 | n/a on YugabyteDB | consistent |
| K1 window | 49.4 | 1.8k | 338 | 10 seats: 4.48 · 100 seats: 8.35 · 1k seats: 105 · 10k seats: 392 · 100k seats: 2948 | n/a on YugabyteDB | consistent |
| L1 claims | 64.3 | 1.9k | 270 | 10 seats: 2.55 · 100 seats: 2.90 · 1k seats: 3.32 · 10k seats: 1.99 · 100k seats: 2.54 | n/a on YugabyteDB | consistent |
| L2 document | 109 | 352 | 177 | 10 seats: 4.55 · 100 seats: 4.26 · 1k seats: 3.57 · 10k seats: 4.02 · 100k seats: 4.34 | n/a on YugabyteDB | consistent |
| L3 sharded | 46.1 | 1.7k | 336 | 10 seats: 4.16 · 100 seats: 6.24 · 1k seats: 80 · 10k seats: 185 · 100k seats: 1482 | n/a on YugabyteDB | consistent |


## Controlled pairs — one decision at a time

**This run's error bar.** Ten designs read tickets through byte-identical SQL on identical data,
so their `q05`/`q06` differences are noise: 90 comparisons, median disagreement **1.02x**, worst **1.17x**.
Read rows inside the worst disagreement are omitted below. **Race, lifecycle and write rows are always
shown**, and so is every correctness difference.

### S0 check/RC ✗ → S4 check/SER

*Identical SQL: what does SERIALIZABLE cost, and does it prevent theft?*

| Topology | Measurement | S0 check/RC ✗ | S4 check/SER | Change |
|---|---|---:|---:|---|
| YugabyteDB, 1 node (RF=1) | race, 10 seats — seats/s | ~~61.7~~ ❌ 124 thefts, 71 early rejections; ❌ 12 of 22 deferred confirmations refused; 277 errors | 0.1 ⏱ 2 timed out at 45% sold; ⌛ 2 of 20 events raced within the tier budget | 502.9x slower — ❌ a design with violations has no valid speed |
| YugabyteDB, 1 node (RF=1) | race, 100 seats — seats/s | ~~151~~ ❌ 120 thefts, 58 early rejections; ❌ 17 of 26 deferred confirmations refused; 109 errors | 0.1 ⏱ 3 timed out at 2% sold; ⌛ 3 of 4 events raced within the tier budget | 1384.7x slower — ❌ a design with violations has no valid speed |
| YugabyteDB, 1 node (RF=1) | race, 1k seats — seats/s | ~~161~~ ❌ 100 thefts, 55 early rejections; ❌ 17 of 58 deferred confirmations refused; 121 errors | 0.1 ⏱ 1 timed out at 0% sold | 3083.8x slower — ❌ a design with violations has no valid speed |
| YugabyteDB, 1 node (RF=1) | race, 10k seats — seats/s | ~~96.2~~ ❌ 62 thefts, 35 early rejections; ❌ 9 of 97 deferred confirmations refused; ⏱ 1 timed out at 19% sold; 49 errors | 0.3 ⏱ 1 timed out at 0% sold | 296.5x slower — ❌ a design with violations has no valid speed |
| YugabyteDB, 1 node (RF=1) | lifecycle, 10 seats — refused late/boundary/early; violations | 0/0/0; 63 | 2/0/0; 0 |  |
| YugabyteDB, 1 node (RF=1) | lifecycle, 100 seats — refused late/boundary/early; violations | 0/0/0; 109 | 5/0/0; 0 |  |
| YugabyteDB, 1 node (RF=1) | hold /s | 60.8 | 1.7 | 34.8x slower |
| YugabyteDB, 1 node (RF=1) | release /s | 386 | 1.9k | 4.9x faster |
| YugabyteDB, 1 node (RF=1) | cancel /s | 338 | 329 | 1.0x slower |
| YugabyteDB, 1 node (RF=1) | publish 10 seats p50 ms | 4.79 | 5.02 | 1.0x slower |
| YugabyteDB, 1 node (RF=1) | publish 100 seats p50 ms | 8.14 | 7.05 | 1.2x faster |
| YugabyteDB, 1 node (RF=1) | publish 1k seats p50 ms | 89 | 87 | 1.0x faster |
| YugabyteDB, 1 node (RF=1) | publish 10k seats p50 ms | 263 | 213 | 1.2x faster |
| YugabyteDB, 1 node (RF=1) | publish 100k seats p50 ms | 2847 | 2887 | 1.0x slower |
| YugabyteDB, 1 node (RF=1) | Seat map of a section — 10 seats | 3.9k | 3.2k | 1.2x slower |
| YugabyteDB, 1 node (RF=1) | Seat map of a section — 100 seats | 3.7k | 3.1k | 1.2x slower |
| YugabyteDB, 1 node (RF=1) | Seats available per section — 1k seats | 1.9k | 1.5k | 1.2x slower |
| YugabyteDB, 1 node (RF=1) | A customer's tickets | 2.9k | 2.4k | 1.2x slower |

### S4 check/SER → S1 conditional

*Serializable read-then-write, or one conditional update?*

| Topology | Measurement | S4 check/SER | S1 conditional | Change |
|---|---|---:|---:|---|
| YugabyteDB, 1 node (RF=1) | race, 10 seats — seats/s | 0.1 ⏱ 2 timed out at 45% sold; ⌛ 2 of 20 events raced within the tier budget | 50.3 | 410.3x faster |
| YugabyteDB, 1 node (RF=1) | race, 100 seats — seats/s | 0.1 ⏱ 3 timed out at 2% sold; ⌛ 3 of 4 events raced within the tier budget | 126 | 1152.6x faster |
| YugabyteDB, 1 node (RF=1) | race, 1k seats — seats/s | 0.1 ⏱ 1 timed out at 0% sold | 133 | 2550.0x faster |
| YugabyteDB, 1 node (RF=1) | race, 10k seats — seats/s | 0.3 ⏱ 1 timed out at 0% sold | ~~90.6~~ ❌ 1 early rejections; ⏱ 1 timed out at 18% sold | 279.3x faster — ❌ a design with violations has no valid speed |
| YugabyteDB, 1 node (RF=1) | lifecycle, 10 seats — refused late/boundary/early; violations | 2/0/0; 0 | 0/0/0; 0 |  |
| YugabyteDB, 1 node (RF=1) | lifecycle, 100 seats — refused late/boundary/early; violations | 5/0/0; 0 | 1/0/0; 0 |  |
| YugabyteDB, 1 node (RF=1) | hold /s | 1.7 | 53.2 | 30.4x faster |
| YugabyteDB, 1 node (RF=1) | release /s | 1.9k | 378 | 5.0x slower |
| YugabyteDB, 1 node (RF=1) | cancel /s | 329 | 345 | 1.0x faster |
| YugabyteDB, 1 node (RF=1) | publish 10 seats p50 ms | 5.02 | 4.84 | 1.0x faster |
| YugabyteDB, 1 node (RF=1) | publish 100 seats p50 ms | 7.05 | 8.67 | 1.2x slower |
| YugabyteDB, 1 node (RF=1) | publish 1k seats p50 ms | 87 | 98 | 1.1x slower |
| YugabyteDB, 1 node (RF=1) | publish 10k seats p50 ms | 213 | 395 | 1.9x slower |
| YugabyteDB, 1 node (RF=1) | publish 100k seats p50 ms | 2887 | 2810 | 1.0x faster |
| YugabyteDB, 1 node (RF=1) | Seat map of a section — 100 seats | 3.1k | 3.6k | 1.2x faster |
| YugabyteDB, 1 node (RF=1) | Seats available per section — 1k seats | 1.5k | 1.9k | 1.3x faster |

### S1 conditional → S2 lock

*Optimistic conditional update, or lock the seats first?*

| Topology | Measurement | S1 conditional | S2 lock | Change |
|---|---|---:|---:|---|
| YugabyteDB, 1 node (RF=1) | race, 10 seats — seats/s | 50.3 | 61.7 | 1.2x faster |
| YugabyteDB, 1 node (RF=1) | race, 100 seats — seats/s | 126 | 98.2 | 1.3x slower |
| YugabyteDB, 1 node (RF=1) | race, 1k seats — seats/s | 133 | 120 | 1.1x slower |
| YugabyteDB, 1 node (RF=1) | race, 10k seats — seats/s | ~~90.6~~ ❌ 1 early rejections; ⏱ 1 timed out at 18% sold | ~~85.0~~ ❌ 1 early rejections; ⏱ 1 timed out at 17% sold | 1.1x slower — ❌ a design with violations has no valid speed |
| YugabyteDB, 1 node (RF=1) | lifecycle, 10 seats — refused late/boundary/early; violations | 0/0/0; 0 | 0/0/0; 0 |  |
| YugabyteDB, 1 node (RF=1) | lifecycle, 100 seats — refused late/boundary/early; violations | 1/0/0; 0 | 2/0/0; 0 |  |
| YugabyteDB, 1 node (RF=1) | hold /s | 53.2 | 45.3 | 1.2x slower |
| YugabyteDB, 1 node (RF=1) | release /s | 378 | 1.9k | 5.0x faster |
| YugabyteDB, 1 node (RF=1) | cancel /s | 345 | 278 | 1.2x slower |
| YugabyteDB, 1 node (RF=1) | publish 10 seats p50 ms | 4.84 | 4.35 | 1.1x faster |
| YugabyteDB, 1 node (RF=1) | publish 100 seats p50 ms | 8.67 | 7.12 | 1.2x faster |
| YugabyteDB, 1 node (RF=1) | publish 1k seats p50 ms | 98 | 89 | 1.1x faster |
| YugabyteDB, 1 node (RF=1) | publish 10k seats p50 ms | 395 | 236 | 1.7x faster |
| YugabyteDB, 1 node (RF=1) | publish 100k seats p50 ms | 2810 | 2916 | 1.0x slower |
| YugabyteDB, 1 node (RF=1) | Seats available per section — 1k seats | 1.9k | 1.6k | 1.2x slower |
| YugabyteDB, 1 node (RF=1) | Seat map of a section — 10k seats | 516 | 404 | 1.3x slower |
| YugabyteDB, 1 node (RF=1) | Seats available for an event — 10k seats | 483 | 404 | 1.2x slower |
| YugabyteDB, 1 node (RF=1) | A hold's seats (basket page) | 662 | 520 | 1.3x slower |

### S2 lock → S3 nowait

*Wait for a seat someone is taking, or fail fast and choose again?*

| Topology | Measurement | S2 lock | S3 nowait | Change |
|---|---|---:|---:|---|
| YugabyteDB, 1 node (RF=1) | race, 10 seats — seats/s | 61.7 | 20.7 | 3.0x slower |
| YugabyteDB, 1 node (RF=1) | race, 100 seats — seats/s | 98.2 | 81.5 | 1.2x slower |
| YugabyteDB, 1 node (RF=1) | race, 1k seats — seats/s | 120 | ~~100~~ ❌ 1 early rejections | 1.2x slower — ❌ a design with violations has no valid speed |
| YugabyteDB, 1 node (RF=1) | race, 10k seats — seats/s | ~~85.0~~ ❌ 1 early rejections; ⏱ 1 timed out at 17% sold | 81.2 ⏱ 1 timed out at 17% sold | 1.0x slower — ❌ a design with violations has no valid speed |
| YugabyteDB, 1 node (RF=1) | lifecycle, 10 seats — refused late/boundary/early; violations | 0/0/0; 0 | 0/0/0; 0 |  |
| YugabyteDB, 1 node (RF=1) | lifecycle, 100 seats — refused late/boundary/early; violations | 2/0/0; 0 | 2/0/0; 0 |  |
| YugabyteDB, 1 node (RF=1) | hold /s | 45.3 | 48.2 | 1.1x faster |
| YugabyteDB, 1 node (RF=1) | release /s | 1.9k | 2.1k | 1.1x faster |
| YugabyteDB, 1 node (RF=1) | cancel /s | 278 | 279 | 1.0x faster |
| YugabyteDB, 1 node (RF=1) | publish 10 seats p50 ms | 4.35 | 4.71 | 1.1x slower |
| YugabyteDB, 1 node (RF=1) | publish 100 seats p50 ms | 7.12 | 8.64 | 1.2x slower |
| YugabyteDB, 1 node (RF=1) | publish 1k seats p50 ms | 89 | 93 | 1.1x slower |
| YugabyteDB, 1 node (RF=1) | publish 10k seats p50 ms | 236 | 259 | 1.1x slower |
| YugabyteDB, 1 node (RF=1) | publish 100k seats p50 ms | 2916 | 2780 | 1.0x faster |
| YugabyteDB, 1 node (RF=1) | Seats available for an event — 1k seats | 2.3k | 1.9k | 1.2x slower |
| YugabyteDB, 1 node (RF=1) | Seats available per section — 10k seats | 220 | 181 | 1.2x slower |
| YugabyteDB, 1 node (RF=1) | A hold's seats (basket page) | 520 | 705 | 1.4x faster |

### S1 conditional → E1 sweeper

*Expiry judged at use time, or made to take effect by a sweeper?*

| Topology | Measurement | S1 conditional | E1 sweeper | Change |
|---|---|---:|---:|---|
| YugabyteDB, 1 node (RF=1) | race, 10 seats — seats/s | 50.3 | 45.2 | 1.1x slower |
| YugabyteDB, 1 node (RF=1) | race, 100 seats — seats/s | 126 | 117 | 1.1x slower |
| YugabyteDB, 1 node (RF=1) | race, 1k seats — seats/s | 133 | 118 | 1.1x slower |
| YugabyteDB, 1 node (RF=1) | race, 10k seats — seats/s | ~~90.6~~ ❌ 1 early rejections; ⏱ 1 timed out at 18% sold | 77.0 ⏱ 1 timed out at 16% sold | 1.2x slower — ❌ a design with violations has no valid speed |
| YugabyteDB, 1 node (RF=1) | lifecycle, 10 seats — refused late/boundary/early; violations | 0/0/0; 0 | 0/0/0; 0 |  |
| YugabyteDB, 1 node (RF=1) | lifecycle, 100 seats — refused late/boundary/early; violations | 1/0/0; 0 | 3/0/0; 0 |  |
| YugabyteDB, 1 node (RF=1) | hold /s | 53.2 | 48.2 | 1.1x slower |
| YugabyteDB, 1 node (RF=1) | release /s | 378 | 2.2k | 5.7x faster |
| YugabyteDB, 1 node (RF=1) | cancel /s | 345 | 280 | 1.2x slower |
| YugabyteDB, 1 node (RF=1) | publish 10 seats p50 ms | 4.84 | 4.52 | 1.1x faster |
| YugabyteDB, 1 node (RF=1) | publish 100 seats p50 ms | 8.67 | 7.93 | 1.1x faster |
| YugabyteDB, 1 node (RF=1) | publish 1k seats p50 ms | 98 | 102 | 1.0x slower |
| YugabyteDB, 1 node (RF=1) | publish 10k seats p50 ms | 395 | 296 | 1.3x faster |
| YugabyteDB, 1 node (RF=1) | publish 100k seats p50 ms | 2810 | 2837 | 1.0x slower |
| YugabyteDB, 1 node (RF=1) | Seat map of a section — 10k seats | 516 | 344 | 1.5x slower |
| YugabyteDB, 1 node (RF=1) | Seats available per section — 10k seats | 242 | 196 | 1.2x slower |
| YugabyteDB, 1 node (RF=1) | A hold's seats (basket page) | 662 | 558 | 1.2x slower |

### S1 conditional → E2 cart

*Expiry copied onto every seat row, or stored once on the cart?*

| Topology | Measurement | S1 conditional | E2 cart | Change |
|---|---|---:|---:|---|
| YugabyteDB, 1 node (RF=1) | race, 10 seats — seats/s | 50.3 | 15.3 | 3.3x slower |
| YugabyteDB, 1 node (RF=1) | race, 100 seats — seats/s | 126 | 50.5 | 2.5x slower |
| YugabyteDB, 1 node (RF=1) | race, 1k seats — seats/s | 133 | 75.5 | 1.8x slower |
| YugabyteDB, 1 node (RF=1) | race, 10k seats — seats/s | ~~90.6~~ ❌ 1 early rejections; ⏱ 1 timed out at 18% sold | 63.6 ⏱ 1 timed out at 13% sold | 1.4x slower — ❌ a design with violations has no valid speed |
| YugabyteDB, 1 node (RF=1) | lifecycle, 10 seats — refused late/boundary/early; violations | 0/0/0; 0 | 1/0/0; 0 |  |
| YugabyteDB, 1 node (RF=1) | lifecycle, 100 seats — refused late/boundary/early; violations | 1/0/0; 0 | 3/0/0; 0 |  |
| YugabyteDB, 1 node (RF=1) | hold /s | 53.2 | 42.5 | 1.3x slower |
| YugabyteDB, 1 node (RF=1) | release /s | 378 | 319 | 1.2x slower |
| YugabyteDB, 1 node (RF=1) | cancel /s | 345 | 295 | 1.2x slower |
| YugabyteDB, 1 node (RF=1) | publish 10 seats p50 ms | 4.84 | 4.63 | 1.0x faster |
| YugabyteDB, 1 node (RF=1) | publish 100 seats p50 ms | 8.67 | 7.48 | 1.2x faster |
| YugabyteDB, 1 node (RF=1) | publish 1k seats p50 ms | 98 | 93 | 1.1x faster |
| YugabyteDB, 1 node (RF=1) | publish 10k seats p50 ms | 395 | 239 | 1.6x faster |
| YugabyteDB, 1 node (RF=1) | publish 100k seats p50 ms | 2810 | 2879 | 1.0x slower |
| YugabyteDB, 1 node (RF=1) | Seat map of a section — 10 seats | 3.5k | 2.6k | 1.3x slower |
| YugabyteDB, 1 node (RF=1) | Seats available per section — 10 seats | 2.8k | 1.7k | 1.7x slower |
| YugabyteDB, 1 node (RF=1) | Seats available for an event — 10 seats | 3.4k | 1.6k | 2.1x slower |
| YugabyteDB, 1 node (RF=1) | Seat map of a section — 100 seats | 3.6k | 1.5k | 2.4x slower |
| YugabyteDB, 1 node (RF=1) | Seats available per section — 100 seats | 2.9k | 1.5k | 1.9x slower |
| YugabyteDB, 1 node (RF=1) | Seats available for an event — 100 seats | 3.7k | 1.6k | 2.3x slower |
| YugabyteDB, 1 node (RF=1) | Seat map of a section — 1k seats | 2.1k | 1.1k | 1.9x slower |
| YugabyteDB, 1 node (RF=1) | Seats available per section — 1k seats | 1.9k | 1.1k | 1.8x slower |
| YugabyteDB, 1 node (RF=1) | Seats available for an event — 1k seats | 2.3k | 1.3k | 1.8x slower |
| YugabyteDB, 1 node (RF=1) | Seat map of a section — 10k seats | 516 | 351 | 1.5x slower |
| YugabyteDB, 1 node (RF=1) | Seats available per section — 10k seats | 242 | 175 | 1.4x slower |
| YugabyteDB, 1 node (RF=1) | Seats available for an event — 10k seats | 483 | 199 | 2.4x slower |
| YugabyteDB, 1 node (RF=1) | A hold's seats (basket page) | 662 | 543 | 1.2x slower |

### E0 app clock ✗ → S1 conditional

*The application's clock, or the database's?*

| Topology | Measurement | E0 app clock ✗ | S1 conditional | Change |
|---|---|---:|---:|---|
| YugabyteDB, 1 node (RF=1) | race, 10 seats — seats/s | 50.2 | 50.3 | 1.0x faster |
| YugabyteDB, 1 node (RF=1) | race, 100 seats — seats/s | 113 | 126 | 1.1x faster |
| YugabyteDB, 1 node (RF=1) | race, 1k seats — seats/s | 140 | 133 | 1.1x slower |
| YugabyteDB, 1 node (RF=1) | race, 10k seats — seats/s | ~~79.9~~ ❌ 1 early rejections; ⏱ 1 timed out at 16% sold | ~~90.6~~ ❌ 1 early rejections; ⏱ 1 timed out at 18% sold | 1.1x faster — ❌ a design with violations has no valid speed |
| YugabyteDB, 1 node (RF=1) | lifecycle, 10 seats — refused late/boundary/early; violations | 0/0/0; 8 | 0/0/0; 0 |  |
| YugabyteDB, 1 node (RF=1) | lifecycle, 100 seats — refused late/boundary/early; violations | 0/0/1; 49 | 1/0/0; 0 |  |
| YugabyteDB, 1 node (RF=1) | hold /s | 46.8 | 53.2 | 1.1x faster |
| YugabyteDB, 1 node (RF=1) | release /s | 370 | 378 | 1.0x faster |
| YugabyteDB, 1 node (RF=1) | cancel /s | 261 | 345 | 1.3x faster |
| YugabyteDB, 1 node (RF=1) | publish 10 seats p50 ms | 4.75 | 4.84 | 1.0x slower |
| YugabyteDB, 1 node (RF=1) | publish 100 seats p50 ms | 7.04 | 8.67 | 1.2x slower |
| YugabyteDB, 1 node (RF=1) | publish 1k seats p50 ms | 96 | 98 | 1.0x slower |
| YugabyteDB, 1 node (RF=1) | publish 10k seats p50 ms | 234 | 395 | 1.7x slower |
| YugabyteDB, 1 node (RF=1) | publish 100k seats p50 ms | 2827 | 2810 | 1.0x faster |
| YugabyteDB, 1 node (RF=1) | Seats available per section — 1k seats | 1.5k | 1.9k | 1.3x faster |
| YugabyteDB, 1 node (RF=1) | Seats available for an event — 1k seats | 1.7k | 2.3k | 1.4x faster |
| YugabyteDB, 1 node (RF=1) | Seats available per section — 10k seats | 180 | 242 | 1.3x faster |
| YugabyteDB, 1 node (RF=1) | Seats available for an event — 10k seats | 364 | 483 | 1.3x faster |
| YugabyteDB, 1 node (RF=1) | A hold's seats (basket page) | 502 | 662 | 1.3x faster |

### K0 naive ✗ → S1 conditional

*What does checking the hold at confirmation cost, and what does it prevent?*

| Topology | Measurement | K0 naive ✗ | S1 conditional | Change |
|---|---|---:|---:|---|
| YugabyteDB, 1 node (RF=1) | race, 10 seats — seats/s | 46.8 | 50.3 | 1.1x faster |
| YugabyteDB, 1 node (RF=1) | race, 100 seats — seats/s | 135 | 126 | 1.1x slower |
| YugabyteDB, 1 node (RF=1) | race, 1k seats — seats/s | 112 | 133 | 1.2x faster |
| YugabyteDB, 1 node (RF=1) | race, 10k seats — seats/s | 87.0 ⏱ 1 timed out at 18% sold | ~~90.6~~ ❌ 1 early rejections; ⏱ 1 timed out at 18% sold | 1.0x faster — ❌ a design with violations has no valid speed |
| YugabyteDB, 1 node (RF=1) | lifecycle, 10 seats — refused late/boundary/early; violations | 0/0/0; 2 | 0/0/0; 0 |  |
| YugabyteDB, 1 node (RF=1) | lifecycle, 100 seats — refused late/boundary/early; violations | 0/0/0; 4 | 1/0/0; 0 |  |
| YugabyteDB, 1 node (RF=1) | hold /s | 51.9 | 53.2 | 1.0x faster |
| YugabyteDB, 1 node (RF=1) | release /s | 368 | 378 | 1.0x faster |
| YugabyteDB, 1 node (RF=1) | cancel /s | 359 | 345 | 1.0x slower |
| YugabyteDB, 1 node (RF=1) | publish 10 seats p50 ms | 4.63 | 4.84 | 1.0x slower |
| YugabyteDB, 1 node (RF=1) | publish 100 seats p50 ms | 8.88 | 8.67 | 1.0x faster |
| YugabyteDB, 1 node (RF=1) | publish 1k seats p50 ms | 89 | 98 | 1.1x slower |
| YugabyteDB, 1 node (RF=1) | publish 10k seats p50 ms | 216 | 395 | 1.8x slower |
| YugabyteDB, 1 node (RF=1) | publish 100k seats p50 ms | 2745 | 2810 | 1.0x slower |

### S1 conditional → K1 window

*What does a guaranteed payment window cost, and what does it remove?*

| Topology | Measurement | S1 conditional | K1 window | Change |
|---|---|---:|---:|---|
| YugabyteDB, 1 node (RF=1) | race, 10 seats — seats/s | 50.3 | 39.9 | 1.3x slower |
| YugabyteDB, 1 node (RF=1) | race, 100 seats — seats/s | 126 | 113 | 1.1x slower |
| YugabyteDB, 1 node (RF=1) | race, 1k seats — seats/s | 133 | 114 | 1.2x slower |
| YugabyteDB, 1 node (RF=1) | race, 10k seats — seats/s | ~~90.6~~ ❌ 1 early rejections; ⏱ 1 timed out at 18% sold | ~~83.4~~ ❌ 1 early rejections; ⏱ 1 timed out at 17% sold | 1.1x slower — ❌ a design with violations has no valid speed |
| YugabyteDB, 1 node (RF=1) | lifecycle, 10 seats — refused late/boundary/early; violations | 0/0/0; 0 | 0/0/0; 0 |  |
| YugabyteDB, 1 node (RF=1) | lifecycle, 100 seats — refused late/boundary/early; violations | 1/0/0; 0 | 0/0/0; 0 |  |
| YugabyteDB, 1 node (RF=1) | hold /s | 53.2 | 49.4 | 1.1x slower |
| YugabyteDB, 1 node (RF=1) | release /s | 378 | 1.8k | 4.6x faster |
| YugabyteDB, 1 node (RF=1) | cancel /s | 345 | 338 | 1.0x slower |
| YugabyteDB, 1 node (RF=1) | publish 10 seats p50 ms | 4.84 | 4.48 | 1.1x faster |
| YugabyteDB, 1 node (RF=1) | publish 100 seats p50 ms | 8.67 | 8.35 | 1.0x faster |
| YugabyteDB, 1 node (RF=1) | publish 1k seats p50 ms | 98 | 105 | 1.1x slower |
| YugabyteDB, 1 node (RF=1) | publish 10k seats p50 ms | 395 | 392 | 1.0x faster |
| YugabyteDB, 1 node (RF=1) | publish 100k seats p50 ms | 2810 | 2948 | 1.0x slower |
| YugabyteDB, 1 node (RF=1) | Seats available per section — 1k seats | 1.9k | 1.6k | 1.2x slower |
| YugabyteDB, 1 node (RF=1) | Seat map of a section — 10k seats | 516 | 373 | 1.4x slower |
| YugabyteDB, 1 node (RF=1) | Seats available per section — 10k seats | 242 | 193 | 1.3x slower |
| YugabyteDB, 1 node (RF=1) | Seats available for an event — 10k seats | 483 | 405 | 1.2x slower |
| YugabyteDB, 1 node (RF=1) | A hold's seats (basket page) | 662 | 510 | 1.3x slower |

### S1 conditional → L1 claims

*Pre-created per-event seat rows, or claims created on hold?*

| Topology | Measurement | S1 conditional | L1 claims | Change |
|---|---|---:|---:|---|
| YugabyteDB, 1 node (RF=1) | race, 10 seats — seats/s | 50.3 | 26.5 | 1.9x slower |
| YugabyteDB, 1 node (RF=1) | race, 100 seats — seats/s | 126 | 87.1 | 1.4x slower |
| YugabyteDB, 1 node (RF=1) | race, 1k seats — seats/s | 133 | 141 | 1.1x faster |
| YugabyteDB, 1 node (RF=1) | race, 10k seats — seats/s | ~~90.6~~ ❌ 1 early rejections; ⏱ 1 timed out at 18% sold | 154 ⏱ 1 timed out at 31% sold | 1.7x faster — ❌ a design with violations has no valid speed |
| YugabyteDB, 1 node (RF=1) | lifecycle, 10 seats — refused late/boundary/early; violations | 0/0/0; 0 | 0/0/0; 0 |  |
| YugabyteDB, 1 node (RF=1) | lifecycle, 100 seats — refused late/boundary/early; violations | 1/0/0; 0 | 0/0/0; 0 |  |
| YugabyteDB, 1 node (RF=1) | hold /s | 53.2 | 64.3 | 1.2x faster |
| YugabyteDB, 1 node (RF=1) | release /s | 378 | 1.9k | 5.0x faster |
| YugabyteDB, 1 node (RF=1) | cancel /s | 345 | 270 | 1.3x slower |
| YugabyteDB, 1 node (RF=1) | publish 10 seats p50 ms | 4.84 | 2.55 | 1.9x faster |
| YugabyteDB, 1 node (RF=1) | publish 100 seats p50 ms | 8.67 | 2.90 | 3.0x faster |
| YugabyteDB, 1 node (RF=1) | publish 1k seats p50 ms | 98 | 3.32 | 29.5x faster |
| YugabyteDB, 1 node (RF=1) | publish 10k seats p50 ms | 395 | 1.99 | 197.9x faster |
| YugabyteDB, 1 node (RF=1) | publish 100k seats p50 ms | 2810 | 2.54 | 1108.3x faster |
| YugabyteDB, 1 node (RF=1) | Seats available per section — 10 seats | 2.8k | 1.5k | 1.9x slower |
| YugabyteDB, 1 node (RF=1) | Seats available for an event — 10 seats | 3.4k | 1.7k | 2.1x slower |
| YugabyteDB, 1 node (RF=1) | Seats available per section — 100 seats | 2.9k | 1.4k | 2.1x slower |
| YugabyteDB, 1 node (RF=1) | Seats available for an event — 100 seats | 3.7k | 1.8k | 2.0x slower |
| YugabyteDB, 1 node (RF=1) | Seat map of a section — 1k seats | 2.1k | 2.9k | 1.4x faster |
| YugabyteDB, 1 node (RF=1) | Seats available per section — 1k seats | 1.9k | 1.2k | 1.6x slower |
| YugabyteDB, 1 node (RF=1) | Seats available for an event — 1k seats | 2.3k | 1.6k | 1.4x slower |
| YugabyteDB, 1 node (RF=1) | Seat map of a section — 10k seats | 516 | 887 | 1.7x faster |
| YugabyteDB, 1 node (RF=1) | Seats available per section — 10k seats | 242 | 379 | 1.6x faster |
| YugabyteDB, 1 node (RF=1) | Seats available for an event — 10k seats | 483 | 773 | 1.6x faster |
| YugabyteDB, 1 node (RF=1) | A hold's seats (basket page) | 662 | 1.1k | 1.7x faster |

### S1 conditional → L2 document

*Seat rows, or an embedded section document?*

| Topology | Measurement | S1 conditional | L2 document | Change |
|---|---|---:|---:|---|
| YugabyteDB, 1 node (RF=1) | race, 10 seats — seats/s | 50.3 | 33.8 | 1.5x slower |
| YugabyteDB, 1 node (RF=1) | race, 100 seats — seats/s | 126 | 60.3 | 2.1x slower |
| YugabyteDB, 1 node (RF=1) | race, 1k seats — seats/s | 133 | 82.4 | 1.6x slower |
| YugabyteDB, 1 node (RF=1) | race, 10k seats — seats/s | ~~90.6~~ ❌ 1 early rejections; ⏱ 1 timed out at 18% sold | 143 ⏱ 1 timed out at 29% sold | 1.6x faster — ❌ a design with violations has no valid speed |
| YugabyteDB, 1 node (RF=1) | lifecycle, 10 seats — refused late/boundary/early; violations | 0/0/0; 0 | 0/0/0; 0 |  |
| YugabyteDB, 1 node (RF=1) | lifecycle, 100 seats — refused late/boundary/early; violations | 1/0/0; 0 | 2/0/0; 0 |  |
| YugabyteDB, 1 node (RF=1) | hold /s | 53.2 | 109 | 2.0x faster |
| YugabyteDB, 1 node (RF=1) | release /s | 378 | 352 | 1.1x slower |
| YugabyteDB, 1 node (RF=1) | cancel /s | 345 | 177 | 2.0x slower |
| YugabyteDB, 1 node (RF=1) | publish 10 seats p50 ms | 4.84 | 4.55 | 1.1x faster |
| YugabyteDB, 1 node (RF=1) | publish 100 seats p50 ms | 8.67 | 4.26 | 2.0x faster |
| YugabyteDB, 1 node (RF=1) | publish 1k seats p50 ms | 98 | 3.57 | 27.5x faster |
| YugabyteDB, 1 node (RF=1) | publish 10k seats p50 ms | 395 | 4.02 | 98.2x faster |
| YugabyteDB, 1 node (RF=1) | publish 100k seats p50 ms | 2810 | 4.34 | 648.1x faster |
| YugabyteDB, 1 node (RF=1) | Seats available per section — 10 seats | 2.8k | 4.0k | 1.4x faster |
| YugabyteDB, 1 node (RF=1) | Seats available per section — 100 seats | 2.9k | 3.9k | 1.3x faster |
| YugabyteDB, 1 node (RF=1) | Seat map of a section — 1k seats | 2.1k | 3.3k | 1.6x faster |
| YugabyteDB, 1 node (RF=1) | Seats available per section — 1k seats | 1.9k | 2.9k | 1.5x faster |
| YugabyteDB, 1 node (RF=1) | Seats available for an event — 1k seats | 2.3k | 3.0k | 1.3x faster |
| YugabyteDB, 1 node (RF=1) | Seat map of a section — 10k seats | 516 | 3.1k | 6.0x faster |
| YugabyteDB, 1 node (RF=1) | Seats available per section — 10k seats | 242 | 1.0k | 4.2x faster |
| YugabyteDB, 1 node (RF=1) | Seats available for an event — 10k seats | 483 | 978 | 2.0x faster |
| YugabyteDB, 1 node (RF=1) | A hold's seats (basket page) | 662 | 3.8k | 5.8x faster |

### S1 conditional → L3 sharded

*(YugabyteDB) all of an event's seats on one tablet, or spread by section?*

| Topology | Measurement | S1 conditional | L3 sharded | Change |
|---|---|---:|---:|---|
| YugabyteDB, 1 node (RF=1) | race, 10 seats — seats/s | 50.3 | 15.5 | 3.2x slower |
| YugabyteDB, 1 node (RF=1) | race, 100 seats — seats/s | 126 | 46.8 | 2.7x slower |
| YugabyteDB, 1 node (RF=1) | race, 1k seats — seats/s | 133 | 68.8 | 1.9x slower |
| YugabyteDB, 1 node (RF=1) | race, 10k seats — seats/s | ~~90.6~~ ❌ 1 early rejections; ⏱ 1 timed out at 18% sold | 105 ⏱ 1 timed out at 21% sold | 1.2x faster — ❌ a design with violations has no valid speed |
| YugabyteDB, 1 node (RF=1) | lifecycle, 10 seats — refused late/boundary/early; violations | 0/0/0; 0 | 0/0/0; 0 |  |
| YugabyteDB, 1 node (RF=1) | lifecycle, 100 seats — refused late/boundary/early; violations | 1/0/0; 0 | 0/0/0; 0 |  |
| YugabyteDB, 1 node (RF=1) | hold /s | 53.2 | 46.1 | 1.2x slower |
| YugabyteDB, 1 node (RF=1) | release /s | 378 | 1.7k | 4.5x faster |
| YugabyteDB, 1 node (RF=1) | cancel /s | 345 | 336 | 1.0x slower |
| YugabyteDB, 1 node (RF=1) | publish 10 seats p50 ms | 4.84 | 4.16 | 1.2x faster |
| YugabyteDB, 1 node (RF=1) | publish 100 seats p50 ms | 8.67 | 6.24 | 1.4x faster |
| YugabyteDB, 1 node (RF=1) | publish 1k seats p50 ms | 98 | 80 | 1.2x faster |
| YugabyteDB, 1 node (RF=1) | publish 10k seats p50 ms | 395 | 185 | 2.1x faster |
| YugabyteDB, 1 node (RF=1) | publish 100k seats p50 ms | 2810 | 1482 | 1.9x faster |
| YugabyteDB, 1 node (RF=1) | Seats available per section — 10 seats | 2.8k | 179 | 15.6x slower |
| YugabyteDB, 1 node (RF=1) | Seats available for an event — 10 seats | 3.4k | 85.7 | 40.1x slower |
| YugabyteDB, 1 node (RF=1) | Seats available per section — 100 seats | 2.9k | 235 | 12.2x slower |
| YugabyteDB, 1 node (RF=1) | Seats available for an event — 100 seats | 3.7k | 93.7 | 39.6x slower |
| YugabyteDB, 1 node (RF=1) | Seat map of a section — 1k seats | 2.1k | 2.8k | 1.3x faster |
| YugabyteDB, 1 node (RF=1) | Seats available per section — 1k seats | 1.9k | 179 | 10.8x slower |
| YugabyteDB, 1 node (RF=1) | Seats available for an event — 1k seats | 2.3k | 104 | 22.5x slower |
| YugabyteDB, 1 node (RF=1) | Seat map of a section — 10k seats | 516 | 2.5k | 4.9x faster |
| YugabyteDB, 1 node (RF=1) | Seats available per section — 10k seats | 242 | 127 | 1.9x slower |
| YugabyteDB, 1 node (RF=1) | Seats available for an event — 10k seats | 483 | 105 | 4.6x slower |
| YugabyteDB, 1 node (RF=1) | A hold's seats (basket page) | 662 | 3.0k | 4.6x faster |


## Measurement hazard: CPU-quota throttling

Cells show *throttled periods / periods (throttled time)*: for the **database** container(s) over
the whole cell (summed across nodes), and for the **client** per phase.

### YugabyteDB, 1 node (RF=1)

| Design | database, whole cell | client lifecycle#1 | client race#1 | client read | client write:cancel | client write:hold | client write:publish | client write:release |
|---|---:|---:|---:|---:|---:|---:|---:|---:|
| S0 check/RC ✗ | 1285/4125 (386s) | 0/2052 (0.0s) | 0/406 (0.0s) | 0/453 (0.0s) | 0/59 (0.0s) | 0/133 (0.0s) | 0/93 (0.0s) | 0/2 (0.0s) |
| S1 conditional | 1284/3617 (381s) | 0/1811 (0.0s) | 0/392 (0.0s) | 0/453 (0.0s) | 0/60 (0.0s) | 0/134 (0.0s) | 0/101 (0.0s) | 0/1 (0.0s) |
| S2 lock | 1333/3838 (376s) | 0/2071 (0.0s) | 0/410 (0.0s) | 0/453 (0.0s) | 0/66 (0.0s) | 0/151 (0.0s) | 0/89 (0.0s) | 0/1 (0.0s) |
| S3 nowait | 1436/4583 (409s) | 0/2621 (0.0s) | 0/490 (0.0s) | 0/452 (0.0s) | 0/66 (0.0s) | 0/155 (0.0s) | 0/101 (0.0s) | 0/1 (0.0s) |
| S4 check/SER | 4913/19233 (1783s) | 0/2412 (0.0s) | 0/671 (0.0s) | 0/452 (0.0s) | 0/59 (0.0s) | 0/5869 (0.0s) | 0/95 (0.0s) | 0/1 (0.0s) |
| E0 app clock ✗ | 1328/3838 (398s) | 0/2061 (0.0s) | 0/390 (0.0s) | 0/453 (0.0s) | 0/67 (0.0s) | 0/154 (0.0s) | 0/101 (0.0s) | 0/2 (0.0s) |
| E1 sweeper | 1337/3714 (403s) | 0/1861 (0.0s) | 0/403 (0.0s) | 0/453 (0.0s) | 0/63 (0.0s) | 0/144 (0.0s) | 0/98 (0.0s) | 0/1 (0.0s) |
| E2 cart | 1574/5164 (448s) | 0/2957 (0.0s) | 0/597 (0.0s) | 0/452 (0.0s) | 0/64 (0.0s) | 0/190 (0.0s) | 0/97 (0.0s) | 0/3 (0.0s) |
| K0 naive ✗ | 1298/4316 (382s) | 0/2405 (0.0s) | 0/408 (0.0s) | 0/453 (0.0s) | 0/62 (0.0s) | 0/137 (0.0s) | 0/92 (0.0s) | 0/2 (0.0s) |
| K1 window | 1364/4071 (386s) | 0/2231 (0.0s) | 0/423 (0.0s) | 0/452 (0.0s) | 0/69 (0.0s) | 0/151 (0.0s) | 0/97 (0.0s) | 0/1 (0.0s) |
| L1 claims | 1235/3679 (337s) | 0/1872 (0.0s) | 0/438 (0.0s) | 0/451 (0.0s) | 0/68 (0.0s) | 0/124 (0.0s) | 0/44 (0.0s) | 0/1 (0.0s) |
| L2 document | 1292/3698 (334s) | 0/1876 (0.0s) | 0/524 (0.0s) | 0/451 (0.0s) | 0/102 (0.0s) | 0/84 (0.0s) | 0/76 (0.0s) | 0/2 (0.0s) |
| L3 sharded | 1716/3780 (586s) | 0/1603 (0.0s) | 0/762 (0.0s) | 0/458 (0.0s) | 0/63 (0.0s) | 0/164 (0.0s) | 0/80 (0.0s) | 0/1 (0.0s) |


---

Plans for every read and write statement are in `results/dc5-yb1-all/<topology>/plans/<design>.txt`
(`EXPLAIN (ANALYZE, BUFFERS)` on PostgreSQL, `EXPLAIN (ANALYZE, DIST)` on YugabyteDB; write
statements inside rolled-back transactions).

## Provenance: code, data and analyses

This file is **generated from the measurements** and contains no interpretation.
Conclusions live in separate signed analyses, so several analysts — including
different AI models — can read the same numbers and each record what they make of
them. Where two analyses disagree, the disagreement is a finding and stays visible.

| | |
|---|---|
| Run id | `devchecks/dc5-yb1-all` |
| Result files | 13 |
| **Inputs digest** | `d3c93141b3825a81` |
| Repository version | `repo/platform-lock-not-available-4-g1024aa7-dirty` (`1024aa73d282`), ⚠️ **dirty** — uncommitted changes; not reproducible from the commit alone |

Quote the digest in any analysis: it identifies the data. The repository version
identifies the code — check it out (`git checkout <commit>`) to reproduce the run, or
`git diff <commit>` to see what has changed in the study since.

### Analyses of this run

**None yet.** This run has measurements but no interpretation.

### Adding an analysis

Copy `reports/analyses/TEMPLATE.md`, fill in the frontmatter, and regenerate this report.
If an analysis already exists for digest `d3c93141b3825a81`, read it first and write a new one only to
**disagree, extend, or bring a different perspective**. Never edit another analyst's file.
