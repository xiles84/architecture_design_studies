# Study 03 — results: reserved seating (venue → event → seat)

Generated from `results/dc3-pg-all`. Every number comes from exactly one JSON file in that
directory. **This file contains measurements and no interpretation**; conclusions live in
signed analyses, indexed at the end.

## TL;DR — measured facts

*Selected from the tables below by fixed rules; no interpretation. Read a signed analysis for
what these facts mean.*

- **Cells:** 12 across 1 topologies; 0 failed.
- **Negative controls:** 4 of 5 fired as designed to. ⚠️ **Did not fire:** E1 during the sweeper outage on PostgreSQL, 1 node — the absence of violations in the other designs of that experiment and topology is not evidence of correctness.
- **Invariant violations in designs meant to be correct:** none observed.
- **Refused confirmations, S1 conditional, PostgreSQL, 1 node** (lifecycle, all tiers): 12 late, 0 boundary, 0 early; 523 confirmed.
- **Refused confirmations, K1 window, PostgreSQL, 1 node** (lifecycle, all tiers): 0 late, 0 boundary, 0 early; 562 confirmed.
- **Sweeper outage, PostgreSQL, 1 node:** E1 0 of 0 probed seats unavailable after expiry; 0 in every other design.
- **Race, PostgreSQL, 1 node** (seats held/s; designs meant to be correct, without violations):
  - 10 seats: highest L1 claims 438/s, lowest L2 document 21.0/s
  - 100 seats: highest L1 claims 1.1k/s, lowest L2 document 19.8/s
  - 1k seats: highest L1 claims 1.9k/s, lowest L2 document 34.9/s
  - 10k seats: highest L1 claims 716/s, lowest L2 document 57.0/s (2 of 9 timed out)
- **Publishing a 100k seats event, PostgreSQL, 1 node** (p50): 6.10–583 ms with per-event inventory rows or documents, 0.68 ms without (L1).

## What was measured

|  |  |
|---|---|
| Run id | `devchecks/dc3-pg-all` |
| Environment | [`host-zenbook-ux5406sa`](../../../docs/environments/host-zenbook-ux5406sa.md) |
| Scale | `tiny` — 5 venues, 84 events, 5741 tickets sold, 131 live and 34 expired holds at load |
| Events | map[catalogue_10:40 catalogue_100:8 catalogue_1000:2 catalogue_10000:1 lifecycle_10:4 lifecycle_100:2 lifecycle_1000:1 race_10:20 race_100:4 race_1000:1 race_10000:1] |
| Seed | 42 — identical data in every cell |
| Reads / isolated writes | 8 workers, 2s measured + 1s warmup, 1 trial(s) |
| Race | 32 buyers per event, real hold TTL 40m0s, 10% of holds deferred until the crowd has gone, timeout 30s, tier budget 3m0s, give up after 20 conflicts |
| Lifecycle — PostgreSQL, 1 node | one human minute = 50ms; TTL 40m0s, payment window (K1) 10m0s, guard G 2m0s, sweeper every 1m0s, outage 45m0s→1h25m0s in the first event of each tier, 25% abandoned (40% of those explicitly), E0 skew 10m0s; tiers 10,100,1000, 2 events per tier, 32 buyers |

**Engines**

| Topology | Version | READ COMMITTED runs as | Client connections to | DB clock − client clock (ms) |
|---|---|---|---|---:|
| PostgreSQL, 1 node | `PostgreSQL 17.11 (Debian 17.11-1.pgdg13+2) on x86_64-pc-linux-gnu, compiled by gcc (Debian 14.2.0-19) 14.2.0, 64-bit` | read committed | 1 node(s) | 0.00 |

> Limits of this environment that bound every number below: the client and the databases
> share one laptop's eight cores; the "3-node" cluster has no real network between nodes; the
> race and the lifecycle are closed loops, so tails are a floor; the lifecycle runs on
> compressed time; E0's clock skew is injected.

## Correctness gate

Six read questions checked against Go-computed truth on a load with sold seats, live holds and
expired holds, plus an audit of the loaded state. A failure aborts the cell.

| Design | PostgreSQL, 1 node |
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

## Negative controls

Deliberately wrong designs (and E1 with its sweeper stopped). Each must be seen to break the
invariant it exists to test; otherwise the absence of violations elsewhere in that experiment
is not evidence of correctness.

| Control | Experiment | Topology | Fired? | What was found |
|---|---|---|---|---|
| S0 check/RC ✗ | race | PostgreSQL, 1 node | ✅ fired | 89 deferred confirmations refused; 3857 thefts, 2009 early rejections |
| S0 check/RC ✗ | lifecycle | PostgreSQL, 1 node | ✅ fired | 418 thefts |
| E0 app clock ✗ | lifecycle | PostgreSQL, 1 node | ✅ fired | 66 thefts, 13 early rejections |
| K0 naive ✗ | lifecycle | PostgreSQL, 1 node | ✅ fired | 10 sales without the hold, 14 late sales |
| E1 sweeper stopped | lifecycle outage | PostgreSQL, 1 node | ⚠️ **did not fire** | 0 of 0 probed seats unavailable after expiry |

## The hold guarantee — lifecycle

Compressed time. *Rejected* confirmations are classified by how long before the hold's expiry
the confirming transaction started: **late** (after expiry), **boundary** (less than G before),
**early** (G or more before — a violation). Violations are counted by the ledger and the audit.
Release lag is in human minutes.

### PostgreSQL, 1 node

| Design | Tier | Confirmed seats/s | Holds | Abandoned | Expired at check | Refused late / boundary / early | Outage: unavailable / probed | Leaked | Release lag p50/p99 (min) | Idle held seat-min | Violations |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| S0 check/RC ✗ | 10 seats | ~~2.0~~ | 80 | 18+4 | 46 | 0 / 0 / 0 | 0 / 0 | 0 | 44509.2 / 68989.7 | 1734 | ❌ 125 thefts |
| S0 check/RC ✗ | 100 seats | ~~11.9~~ | 293 | 54+29 | 94 | 5 / 0 / 0 | 0 / 0 | 0 | 0.1 / 0.2 | 5808 | ❌ 222 thefts |
| S0 check/RC ✗ | 1k seats | ~~68.0~~ | 624 | 108+60 | 38 | 5 / 0 / 0 | 0 / 0 | 0 | 0.6 / 72.0 | 11396 | ❌ 71 thefts |
| S1 conditional | 10 seats | 2.3 | 21 | 4+4 | 0 | 0 / 0 / 0 | 0 / 0 | 0 | 45491.1 / 68990.4 | 290 | none |
| S1 conditional | 100 seats | 13.0 | 145 | 27+9 | 1 | 2 / 0 / 0 | 0 / 0 | 0 | 0.0 / 0.3 | 2446 | none |
| S1 conditional | 1k seats | 68.9 | 582 | 85+53 | 30 | 10 / 0 / 0 | 0 / 0 | 0 | 0.7 / 80.3 | 9435 | none |
| S2 lock | 10 seats | 1.8 | 21 | 3+2 | 2 | 0 / 0 / 0 | 0 / 0 | 0 | 45491.3 / 68990.6 | 303 | none |
| S2 lock | 100 seats | 14.3 | 156 | 27+18 | 1 | 2 / 0 / 0 | 0 / 0 | 0 | 0.2 / 33.2 | 2297 | none |
| S2 lock | 1k seats | 71.2 | 560 | 87+53 | 13 | 8 / 0 / 0 | 0 / 0 | 0 | 0.6 / 13.5 | 9390 | none |
| S3 nowait | 10 seats | 2.9 | 16 | 2+0 | 1 | 0 / 0 / 0 | 0 / 0 | 0 | 45491.5 / 68990.8 | 160 | none |
| S3 nowait | 100 seats | 14.3 | 137 | 17+10 | 2 | 1 / 0 / 0 | 0 / 0 | 0 | 0.0 / 0.0 | 1496 | none |
| S3 nowait | 1k seats | 71.8 | 587 | 96+48 | 8 | 13 / 0 / 0 | 0 / 0 | 0 | 0.6 / 36.0 | 8933 | none |
| S4 check/SER | 10 seats | 2.4 | 22 | 5+3 | 0 | 0 / 0 / 0 | 0 / 0 | 0 | 44510.5 / 68991.0 | 258 | none |
| S4 check/SER | 100 seats | 11.1 | 152 | 23+15 | 5 | 3 / 0 / 0 | 0 / 0 | 0 | 0.2 / 0.5 | 2100 | none |
| S4 check/SER | 1k seats | 71.8 | 591 | 82+53 | 12 | 3 / 0 / 0 | 0 / 0 | 0 | 0.7 / 71.0 | 8500 | none |
| E0 app clock ✗ | 10 seats | ~~2.7~~ | 19 | 1+2 | 1 | 0 / 0 / 0 | 0 / 0 | 0 | 45490.5 / 68989.8 | 143 | ❌ 4 thefts |
| E0 app clock ✗ | 100 seats | ~~13.4~~ | 177 | 20+19 | 26 | 3 / 0 / 1 | 0 / 0 | 0 | 33.5 / 43.5 | 2081 | ❌ 37 thefts, 1 early rejections |
| E0 app clock ✗ | 1k seats | ~~68.6~~ | 599 | 101+57 | 21 | 0 / 0 / 12 | 0 / 0 | 0 | -9.5 / 1.4 | 10793 | ❌ 25 thefts, 12 early rejections |
| E1 sweeper | 10 seats | 2.3 | 19 | 5+1 | 0 | 0 / 0 / 0 | 0 / 0 | 0 | 0.8 / 0.9 | 445 | none |
| E1 sweeper | 100 seats | 13.8 | 144 | 20+15 | 2 | 2 / 0 / 0 | 0 / 0 | 0 | 0.8 / 39.7 | 1823 | none |
| E1 sweeper | 1k seats | 71.5 | 613 | 97+53 | 25 | 8 / 0 / 0 | 0 / 0 | 0 | 0.8 / 69.0 | 9720 | none |
| E2 cart | 10 seats | 2.7 | 22 | 4+3 | 1 | 1 / 0 / 0 | 0 / 0 | 0 | 45490.4 / 68989.7 | 337 | none |
| E2 cart | 100 seats | 20.4 | 137 | 23+5 | 3 | 4 / 0 / 0 | 0 / 0 | 0 | 0.0 / 0.0 | 1953 | none |
| E2 cart | 1k seats | 70.5 | 597 | 106+54 | 19 | 9 / 0 / 0 | 0 / 0 | 0 | 0.7 / 63.5 | 10968 | none |
| K0 naive ✗ | 10 seats | 1.9 | 24 | 4+2 | 1 | 0 / 0 / 0 | 0 / 0 | 0 | 45491.9 / 68991.2 | 326 | none |
| K0 naive ✗ | 100 seats | ~~17.3~~ | 149 | 24+13 | 5 | 0 / 0 / 0 | 0 / 0 | 0 | 0.1 / 0.1 | 1967 | ❌ 3 sales without the hold, 3 late sales |
| K0 naive ✗ | 1k seats | ~~72.6~~ | 572 | 75+54 | 14 | 0 / 0 / 0 | 0 / 0 | 0 | 0.5 / 24.0 | 7638 | ❌ 7 sales without the hold, 11 late sales |
| K1 window | 10 seats | 1.8 | 22 | 5+2 | 1 | 0 / 0 / 0 | 0 / 0 | 0 | 45490.4 / 68989.7 | 330 | none |
| K1 window | 100 seats | 15.6 | 153 | 27+16 | 4 | 0 / 0 / 0 | 0 / 0 | 0 | 0.2 / 0.2 | 2529 | none |
| K1 window | 1k seats | 70.6 | 595 | 87+58 | 8 | 0 / 0 / 0 | 0 / 0 | 0 | 0.6 / 12.0 | 10320 | none |
| L1 claims | 10 seats | 3.7 | 18 | 0+1 | 1 | 0 / 0 / 0 | 0 / 0 | 0 | 45489.5 / 68988.8 | 5 | none |
| L1 claims | 100 seats | 15.7 | 195 | 28+33 | 19 | 8 / 0 / 0 | 0 / 0 | 0 | 41.3 / 41.6 | 2419 | none |
| L1 claims | 1k seats | 70.3 | 592 | 102+57 | 12 | 8 / 0 / 0 | 0 / 0 | 0 | 0.4 / 22.9 | 10694 | none |
| L2 document | 10 seats | 4.1 | 12 | 0+2 | 0 | 0 / 0 / 0 | 0 / 0 | 0 | 45490.7 / 68990.0 | 38 | none |
| L2 document | 100 seats | 10.4 | 160 | 30+11 | 5 | 2 / 0 / 0 | 0 / 0 | 0 | 0.1 / 0.2 | 2947 | none |
| L2 document | 1k seats | 70.2 | 603 | 89+54 | 29 | 7 / 0 / 0 | 0 / 0 | 0 | 0.9 / 68.2 | 10006 | none |

## The race — a hot drop with seat choice

Seats held per second, from the gate to the last hold granted. Conflicts and seat-map reads are
per granted hold. Deferred confirmations kept a real 40-minute hold through the crowd.

### PostgreSQL, 1 node

| Design | 10 seats | 100 seats | 1k seats | 10k seats |
|---|---:|---:|---:|---:|
| S0 check/RC ✗ | ~~21.3~~ ❌ 1212 thefts, 706 early rejections; ❌ 70 of 71 deferred confirmations refused; 183 errors | ~~108~~ ❌ 491 thefts, 232 early rejections; ❌ 29 of 45 deferred confirmations refused; 40 errors | ~~2.4k~~ ❌ 548 thefts, 275 early rejections; ❌ 37 of 81 deferred confirmations refused; 17 errors | ~~768~~ ❌ 1606 thefts, 796 early rejections; ❌ 89 of 498 deferred confirmations refused; 172 errors |
| S1 conditional | 413 | 1.1k | 716 | 595 |
| S2 lock | 380 | 1.0k | 1.3k | 538 |
| S3 nowait | 217 | 603 | 1.1k | 696 |
| S4 check/SER | 219 | 569 | 1.0k | 596 |
| E0 app clock ✗ | 293 | 966 | 1.5k | 629 |
| E1 sweeper | 301 | 911 | 1.3k | 700 |
| E2 cart | 212 | 638 | 672 | 330 ⏱ 1 timed out at 99% sold |
| K0 naive ✗ | 273 | 933 | 1.2k | 568 |
| K1 window | 414 | 1.0k | 1.3k | 538 |
| L1 claims | 438 | 1.1k | 1.9k | 716 |
| L2 document | 21.0 | 19.8 | 34.9 | 57.0 ⏱ 1 timed out at 17% sold |

<details><summary>Race detail — PostgreSQL, 1 node</summary>

| Design | Tier | Conflicts/hold | Map reads/hold | Engine retries | Gave up | Hold p50/p99 ms | Confirm p99 ms | Deferred ok | Final buyer sold |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|
| S0 check/RC ✗ | 10 seats | 0.17 | 1.57 | 90 | 0 | 17 / 12063 | 4019 | 1/71 | 0 |
| S0 check/RC ✗ | 100 seats | 0.14 | 1.47 | 16 | 0 | 3.94 / 4048 | 20 | 16/45 | 0 |
| S0 check/RC ✗ | 1k seats | 0.11 | 1.21 | 0 | 0 | 5.35 / 81 | 70 | 44/81 | 0 |
| S0 check/RC ✗ | 10k seats | 0.05 | 1.86 | 4 | 0 | 7.50 / 103 | 94 | 409/498 | 0 |
| S1 conditional | 10 seats | 8.34 | 16.37 | 0 | 0 | 2.50 / 54 | 50 | 19/19 | 0 |
| S1 conditional | 100 seats | 2.99 | 5.84 | 0 | 0 | 2.47 / 73 | 59 | 16/16 | 0 |
| S1 conditional | 1k seats | 1.17 | 3.18 | 1 | 0 | 2.71 / 66 | 66 | 54/54 | 0 |
| S1 conditional | 10k seats | 0.40 | 3.64 | 0 | 0 | 5.88 / 94 | 96 | 413/413 | 0 |
| S2 lock | 10 seats | 9.24 | 17.41 | 0 | 0 | 2.75 / 43 | 47 | 12/12 | 0 |
| S2 lock | 100 seats | 3.40 | 6.16 | 0 | 0 | 2.61 / 63 | 62 | 24/24 | 0 |
| S2 lock | 1k seats | 1.43 | 3.74 | 0 | 0 | 3.50 / 71 | 78 | 46/46 | 0 |
| S2 lock | 10k seats | 0.42 | 4.28 | 0 | 0 | 6.31 / 96 | 96 | 387/387 | 0 |
| S3 nowait | 10 seats | 26.00 | 34.91 | 0 | 0 | 6.95 / 65 | 80 | 13/13 | 0 |
| S3 nowait | 100 seats | 7.98 | 10.29 | 0 | 12 | 6.54 / 97 | 98 | 17/17 | 0 |
| S3 nowait | 1k seats | 2.75 | 4.54 | 0 | 5 | 5.73 / 76 | 76 | 39/39 | 0 |
| S3 nowait | 10k seats | 0.61 | 3.59 | 0 | 0 | 7.32 / 93 | 93 | 387/387 | 0 |
| S4 check/SER | 10 seats | 11.08 | 20.71 | 1698 | 0 | 4.64 / 64 | 88 | 10/10 | 0 |
| S4 check/SER | 100 seats | 3.99 | 6.54 | 1510 | 0 | 3.99 / 61 | 224 | 16/16 | 0 |
| S4 check/SER | 1k seats | 0.58 | 2.48 | 1157 | 0 | 2.23 / 60 | 663 | 54/54 | 0 |
| S4 check/SER | 10k seats | 0.44 | 3.61 | 6209 | 0 | 5.23 / 100 | 693 | 405/405 | 0 |
| E0 app clock ✗ | 10 seats | 9.33 | 17.24 | 0 | 0 | 3.25 / 55 | 49 | 21/21 | 0 |
| E0 app clock ✗ | 100 seats | 3.15 | 5.39 | 0 | 0 | 2.89 / 65 | 65 | 15/15 | 0 |
| E0 app clock ✗ | 1k seats | 1.30 | 3.59 | 0 | 0 | 3.08 / 69 | 70 | 51/51 | 0 |
| E0 app clock ✗ | 10k seats | 0.40 | 3.80 | 0 | 0 | 5.84 / 93 | 95 | 412/412 | 0 |
| E1 sweeper | 10 seats | 9.10 | 17.60 | 0 | 0 | 2.71 / 60 | 56 | 8/8 | 0 |
| E1 sweeper | 100 seats | 3.30 | 6.14 | 0 | 0 | 2.32 / 64 | 68 | 19/19 | 0 |
| E1 sweeper | 1k seats | 1.61 | 3.77 | 0 | 0 | 3.00 / 67 | 74 | 40/40 | 0 |
| E1 sweeper | 10k seats | 0.39 | 3.61 | 0 | 0 | 5.39 / 93 | 96 | 407/407 | 0 |
| E2 cart | 10 seats | 9.84 | 17.98 | 0 | 0 | 4.06 / 64 | 71 | 12/12 | 0 |
| E2 cart | 100 seats | 3.68 | 6.14 | 0 | 0 | 3.62 / 80 | 73 | 19/19 | 0 |
| E2 cart | 1k seats | 1.59 | 3.49 | 0 | 0 | 4.51 / 89 | 92 | 47/47 | 0 |
| E2 cart | 10k seats | 0.39 | 3.62 | 0 | 0 | 9.82 / 104 | 104 | 380/380 | 0 |
| K0 naive ✗ | 10 seats | 9.36 | 17.36 | 0 | 0 | 3.34 / 57 | 62 | 16/16 | 0 |
| K0 naive ✗ | 100 seats | 3.09 | 5.31 | 0 | 0 | 3.09 / 60 | 63 | 21/21 | 0 |
| K0 naive ✗ | 1k seats | 1.46 | 4.11 | 0 | 0 | 3.51 / 75 | 76 | 51/51 | 0 |
| K0 naive ✗ | 10k seats | 0.40 | 4.52 | 0 | 0 | 6.05 / 95 | 95 | 395/395 | 0 |
| K1 window | 10 seats | 8.63 | 15.92 | 0 | 0 | 2.41 / 48 | 43 | 20/20 | 0 |
| K1 window | 100 seats | 2.78 | 5.71 | 0 | 0 | 2.29 / 61 | 65 | 15/15 | 0 |
| K1 window | 1k seats | 1.51 | 4.05 | 0 | 0 | 2.94 / 73 | 78 | 46/46 | 0 |
| K1 window | 10k seats | 0.36 | 4.28 | 0 | 0 | 5.73 / 95 | 97 | 393/393 | 0 |
| L1 claims | 10 seats | 8.01 | 14.82 | 0 | 0 | 2.03 / 67 | 68 | 9/9 | 0 |
| L1 claims | 100 seats | 2.80 | 5.37 | 0 | 0 | 2.23 / 60 | 60 | 20/20 | 0 |
| L1 claims | 1k seats | 1.42 | 3.78 | 0 | 0 | 2.70 / 64 | 70 | 53/53 | 0 |
| L1 claims | 10k seats | 0.39 | 3.70 | 0 | 0 | 4.24 / 94 | 94 | 409/409 | 0 |
| L2 document | 10 seats | 19.60 | 35.52 | 3622 | 0 | 58 / 181 | 91 | 20/20 | 0 |
| L2 document | 100 seats | 8.53 | 13.27 | 6822 | 0 | 68 / 457 | 89 | 25/25 | 0 |
| L2 document | 1k seats | 3.68 | 6.41 | 9402 | 31 | 60 / 2106 | 137 | 60/60 | 0 |
| L2 document | 10k seats | 1.97 | 3.01 | 9527 | 20 | 4.89 / 2453 | 97 | 73/73 | 0 |

</details>

## Reads

Each question in isolation, fresh key per execution, ops/s.

### PostgreSQL, 1 node

| Question | S0 check/RC ✗ | S1 conditional | S2 lock | S3 nowait | S4 check/SER | E0 app clock ✗ | E1 sweeper | E2 cart | K0 naive ✗ | K1 window | L1 claims | L2 document |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| Seat map of a section — 10 seats | 19k | 18k | 19k | 19k | 19k | 19k | 20k | 12k | 19k | 20k | 35k | 31k |
| Seat map of a section — 100 seats | 16k | 18k | 18k | 14k | 18k | 17k | 18k | 14k | 18k | 18k | 37k | 26k |
| Seat map of a section — 1k seats | 21k | 23k | 21k | 22k | 22k | 23k | 22k | 18k | 22k | 23k | 26k | 18k |
| Seat map of a section — 10k seats | 14k | 16k | 13k | 15k | 15k | 15k | 12k | 12k | 14k | 15k | 15k | 9.7k |
| Seats available per section — 10 seats | 20k | 20k | 21k | 21k | 21k | 20k | 21k | 11k | 20k | 20k | 25k | 31k |
| Seats available per section — 100 seats | 17k | 19k | 19k | 17k | 19k | 19k | 17k | 9.8k | 17k | 19k | 27k | 29k |
| Seats available per section — 1k seats | 8.0k | 7.7k | 8.6k | 8.8k | 8.6k | 7.7k | 8.5k | 4.2k | 6.6k | 8.6k | 10k | 12k |
| Seats available per section — 10k seats | 941 | 875 | 1.4k | 1.4k | 1.2k | 1.0k | 1.4k | 804 | 1.2k | 1.1k | 1.0k | 1.4k |
| Seats available for an event — 10 seats | 20k | 20k | 20k | 21k | 20k | 20k | 23k | 11k | 21k | 21k | 14k | 29k |
| Seats available for an event — 100 seats | 20k | 19k | 19k | 19k | 19k | 19k | 21k | 9.9k | 19k | 20k | 14k | 26k |
| Seats available for an event — 1k seats | 10k | 7.9k | 8.9k | 9.2k | 8.8k | 7.6k | 11k | 4.9k | 9.3k | 8.5k | 7.6k | 11k |
| Seats available for an event — 10k seats | 1.0k | 1.3k | 1.1k | 1.3k | 1.3k | 1.0k | 1.7k | 860 | 1.4k | 903 | 1.2k | 1.3k |
| A hold's seats (basket page) | 26k | 27k | 27k | 26k | 25k | 23k | 26k | 36k | 26k | 27k | 36k | 23k |
| A customer's tickets | 39k | 39k | 38k | 34k | 38k | 38k | 38k | 39k | 36k | 38k | 39k | 32k |
| Ticket by id (control) | 48k | 44k | 45k | 37k | 46k | 47k | 43k | 46k | 41k | 47k | 47k | 33k |

## Isolated writes, publishing and storage

### PostgreSQL, 1 node

| Design | Hold+confirm /s | Release /s | Refund /s | Publish p50 ms by size | Storage | Audit |
|---|---:|---:|---:|---:|---:|---|
| S0 check/RC ✗ | ~~650~~ ❌ (17 errors) | 5.3k | 4.8k | 10 seats: 1.44 · 100 seats: 1.86 · 1k seats: 7.40 · 10k seats: 42 · 100k seats: 491 | 14.8 MiB | ❌ 71 thefts, 39 early rejections |
| S1 conditional | 723 | 5.2k | 5.4k | 10 seats: 1.47 · 100 seats: 1.97 · 1k seats: 7.26 · 10k seats: 37 · 100k seats: 475 | 14.4 MiB | consistent |
| S2 lock | 551 | 4.6k | 4.5k | 10 seats: 1.69 · 100 seats: 2.01 · 1k seats: 8.75 · 10k seats: 47 · 100k seats: 491 | 15.4 MiB | consistent |
| S3 nowait | 559 | 4.4k | 4.2k | 10 seats: 1.73 · 100 seats: 2.38 · 1k seats: 11 · 10k seats: 57 · 100k seats: 583 | 14.7 MiB | consistent |
| S4 check/SER | 512 | 4.3k | 4.0k | 10 seats: 1.79 · 100 seats: 2.36 · 1k seats: 11 · 10k seats: 54 · 100k seats: 580 | 15.2 MiB | consistent |
| E0 app clock ✗ | ~~696~~ ❌ | 4.3k | 5.4k | 10 seats: 1.47 · 100 seats: 2.00 · 1k seats: 7.48 · 10k seats: 43 · 100k seats: 482 | 14.9 MiB | ❌ 8 thefts |
| E1 sweeper | 612 | 4.2k | 4.7k | 10 seats: 1.65 · 100 seats: 2.12 · 1k seats: 9.61 · 10k seats: 67 · 100k seats: 566 | 14.4 MiB | consistent |
| E2 cart | 442 | 5.2k | 5.6k | 10 seats: 1.41 · 100 seats: 1.92 · 1k seats: 8.65 · 10k seats: 54 · 100k seats: 493 | 15.4 MiB | consistent |
| K0 naive ✗ | 598 | 4.3k | 3.7k | 10 seats: 1.70 · 100 seats: 2.38 · 1k seats: 10 · 10k seats: 53 · 100k seats: 576 | 14.6 MiB | consistent |
| K1 window | 701 | 5.6k | 5.3k | 10 seats: 1.47 · 100 seats: 1.86 · 1k seats: 7.31 · 10k seats: 46 · 100k seats: 464 | 15.6 MiB | consistent |
| L1 claims | 713 | 6.1k | 5.7k | 10 seats: 1.00 · 100 seats: 1.10 · 1k seats: 1.10 · 10k seats: 1.23 · 100k seats: 0.68 | 12.3 MiB | consistent |
| L2 document | 367 | 2.9k | 1.2k | 10 seats: 5.53 · 100 seats: 5.42 · 1k seats: 5.30 · 10k seats: 3.40 · 100k seats: 6.10 | 11.4 MiB | consistent |


## Controlled pairs — one decision at a time

**This run's error bar.** Ten designs read tickets through byte-identical SQL on identical data,
so their `q05`/`q06` differences are noise: 90 comparisons, median disagreement **1.08x**, worst **1.45x**.
Read rows inside the worst disagreement are omitted below. **Race, lifecycle and write rows are always
shown**, and so is every correctness difference.

### S0 check/RC ✗ → S4 check/SER

*Identical SQL: what does SERIALIZABLE cost, and does it prevent theft?*

| Topology | Measurement | S0 check/RC ✗ | S4 check/SER | Change |
|---|---|---:|---:|---|
| PostgreSQL, 1 node | race, 10 seats — seats/s | ~~21.3~~ ❌ 1212 thefts, 706 early rejections; ❌ 70 of 71 deferred confirmations refused; 183 errors | 219 | 10.2x faster — ❌ a design with violations has no valid speed |
| PostgreSQL, 1 node | race, 100 seats — seats/s | ~~108~~ ❌ 491 thefts, 232 early rejections; ❌ 29 of 45 deferred confirmations refused; 40 errors | 569 | 5.3x faster — ❌ a design with violations has no valid speed |
| PostgreSQL, 1 node | race, 1k seats — seats/s | ~~2.4k~~ ❌ 548 thefts, 275 early rejections; ❌ 37 of 81 deferred confirmations refused; 17 errors | 1.0k | 2.4x slower — ❌ a design with violations has no valid speed |
| PostgreSQL, 1 node | race, 10k seats — seats/s | ~~768~~ ❌ 1606 thefts, 796 early rejections; ❌ 89 of 498 deferred confirmations refused; 172 errors | 596 | 1.3x slower — ❌ a design with violations has no valid speed |
| PostgreSQL, 1 node | lifecycle, 10 seats — refused late/boundary/early; violations | 0/0/0; 125 | 0/0/0; 0 |  |
| PostgreSQL, 1 node | lifecycle, 100 seats — refused late/boundary/early; violations | 5/0/0; 222 | 3/0/0; 0 |  |
| PostgreSQL, 1 node | lifecycle, 1k seats — refused late/boundary/early; violations | 5/0/0; 71 | 3/0/0; 0 |  |
| PostgreSQL, 1 node | hold /s | 650 | 512 | 1.3x slower |
| PostgreSQL, 1 node | release /s | 5.3k | 4.3k | 1.2x slower |
| PostgreSQL, 1 node | cancel /s | 4.8k | 4.0k | 1.2x slower |
| PostgreSQL, 1 node | publish 10 seats p50 ms | 1.44 | 1.79 | 1.2x slower |
| PostgreSQL, 1 node | publish 100 seats p50 ms | 1.86 | 2.36 | 1.3x slower |
| PostgreSQL, 1 node | publish 1k seats p50 ms | 7.40 | 11 | 1.4x slower |
| PostgreSQL, 1 node | publish 10k seats p50 ms | 42 | 54 | 1.3x slower |
| PostgreSQL, 1 node | publish 100k seats p50 ms | 491 | 580 | 1.2x slower |

### S4 check/SER → S1 conditional

*Serializable read-then-write, or one conditional update?*

| Topology | Measurement | S4 check/SER | S1 conditional | Change |
|---|---|---:|---:|---|
| PostgreSQL, 1 node | race, 10 seats — seats/s | 219 | 413 | 1.9x faster |
| PostgreSQL, 1 node | race, 100 seats — seats/s | 569 | 1.1k | 1.9x faster |
| PostgreSQL, 1 node | race, 1k seats — seats/s | 1.0k | 716 | 1.4x slower |
| PostgreSQL, 1 node | race, 10k seats — seats/s | 596 | 595 | 1.0x slower |
| PostgreSQL, 1 node | lifecycle, 10 seats — refused late/boundary/early; violations | 0/0/0; 0 | 0/0/0; 0 |  |
| PostgreSQL, 1 node | lifecycle, 100 seats — refused late/boundary/early; violations | 3/0/0; 0 | 2/0/0; 0 |  |
| PostgreSQL, 1 node | lifecycle, 1k seats — refused late/boundary/early; violations | 3/0/0; 0 | 10/0/0; 0 |  |
| PostgreSQL, 1 node | hold /s | 512 | 723 | 1.4x faster |
| PostgreSQL, 1 node | release /s | 4.3k | 5.2k | 1.2x faster |
| PostgreSQL, 1 node | cancel /s | 4.0k | 5.4k | 1.4x faster |
| PostgreSQL, 1 node | publish 10 seats p50 ms | 1.79 | 1.47 | 1.2x faster |
| PostgreSQL, 1 node | publish 100 seats p50 ms | 2.36 | 1.97 | 1.2x faster |
| PostgreSQL, 1 node | publish 1k seats p50 ms | 11 | 7.26 | 1.4x faster |
| PostgreSQL, 1 node | publish 10k seats p50 ms | 54 | 37 | 1.5x faster |
| PostgreSQL, 1 node | publish 100k seats p50 ms | 580 | 475 | 1.2x faster |

### S1 conditional → S2 lock

*Optimistic conditional update, or lock the seats first?*

| Topology | Measurement | S1 conditional | S2 lock | Change |
|---|---|---:|---:|---|
| PostgreSQL, 1 node | race, 10 seats — seats/s | 413 | 380 | 1.1x slower |
| PostgreSQL, 1 node | race, 100 seats — seats/s | 1.1k | 1.0k | 1.0x slower |
| PostgreSQL, 1 node | race, 1k seats — seats/s | 716 | 1.3k | 1.8x faster |
| PostgreSQL, 1 node | race, 10k seats — seats/s | 595 | 538 | 1.1x slower |
| PostgreSQL, 1 node | lifecycle, 10 seats — refused late/boundary/early; violations | 0/0/0; 0 | 0/0/0; 0 |  |
| PostgreSQL, 1 node | lifecycle, 100 seats — refused late/boundary/early; violations | 2/0/0; 0 | 2/0/0; 0 |  |
| PostgreSQL, 1 node | lifecycle, 1k seats — refused late/boundary/early; violations | 10/0/0; 0 | 8/0/0; 0 |  |
| PostgreSQL, 1 node | hold /s | 723 | 551 | 1.3x slower |
| PostgreSQL, 1 node | release /s | 5.2k | 4.6k | 1.1x slower |
| PostgreSQL, 1 node | cancel /s | 5.4k | 4.5k | 1.2x slower |
| PostgreSQL, 1 node | publish 10 seats p50 ms | 1.47 | 1.69 | 1.1x slower |
| PostgreSQL, 1 node | publish 100 seats p50 ms | 1.97 | 2.01 | 1.0x slower |
| PostgreSQL, 1 node | publish 1k seats p50 ms | 7.26 | 8.75 | 1.2x slower |
| PostgreSQL, 1 node | publish 10k seats p50 ms | 37 | 47 | 1.3x slower |
| PostgreSQL, 1 node | publish 100k seats p50 ms | 475 | 491 | 1.0x slower |
| PostgreSQL, 1 node | Seats available per section — 10k seats | 875 | 1.4k | 1.6x faster |

### S2 lock → S3 nowait

*Wait for a seat someone is taking, or fail fast and choose again?*

| Topology | Measurement | S2 lock | S3 nowait | Change |
|---|---|---:|---:|---|
| PostgreSQL, 1 node | race, 10 seats — seats/s | 380 | 217 | 1.8x slower |
| PostgreSQL, 1 node | race, 100 seats — seats/s | 1.0k | 603 | 1.7x slower |
| PostgreSQL, 1 node | race, 1k seats — seats/s | 1.3k | 1.1k | 1.2x slower |
| PostgreSQL, 1 node | race, 10k seats — seats/s | 538 | 696 | 1.3x faster |
| PostgreSQL, 1 node | lifecycle, 10 seats — refused late/boundary/early; violations | 0/0/0; 0 | 0/0/0; 0 |  |
| PostgreSQL, 1 node | lifecycle, 100 seats — refused late/boundary/early; violations | 2/0/0; 0 | 1/0/0; 0 |  |
| PostgreSQL, 1 node | lifecycle, 1k seats — refused late/boundary/early; violations | 8/0/0; 0 | 13/0/0; 0 |  |
| PostgreSQL, 1 node | hold /s | 551 | 559 | 1.0x faster |
| PostgreSQL, 1 node | release /s | 4.6k | 4.4k | 1.0x slower |
| PostgreSQL, 1 node | cancel /s | 4.5k | 4.2k | 1.1x slower |
| PostgreSQL, 1 node | publish 10 seats p50 ms | 1.69 | 1.73 | 1.0x slower |
| PostgreSQL, 1 node | publish 100 seats p50 ms | 2.01 | 2.38 | 1.2x slower |
| PostgreSQL, 1 node | publish 1k seats p50 ms | 8.75 | 11 | 1.2x slower |
| PostgreSQL, 1 node | publish 10k seats p50 ms | 47 | 57 | 1.2x slower |
| PostgreSQL, 1 node | publish 100k seats p50 ms | 491 | 583 | 1.2x slower |

### S1 conditional → E1 sweeper

*Expiry judged at use time, or made to take effect by a sweeper?*

| Topology | Measurement | S1 conditional | E1 sweeper | Change |
|---|---|---:|---:|---|
| PostgreSQL, 1 node | race, 10 seats — seats/s | 413 | 301 | 1.4x slower |
| PostgreSQL, 1 node | race, 100 seats — seats/s | 1.1k | 911 | 1.2x slower |
| PostgreSQL, 1 node | race, 1k seats — seats/s | 716 | 1.3k | 1.9x faster |
| PostgreSQL, 1 node | race, 10k seats — seats/s | 595 | 700 | 1.2x faster |
| PostgreSQL, 1 node | lifecycle, 10 seats — refused late/boundary/early; violations | 0/0/0; 0 | 0/0/0; 0 |  |
| PostgreSQL, 1 node | lifecycle, 100 seats — refused late/boundary/early; violations | 2/0/0; 0 | 2/0/0; 0 |  |
| PostgreSQL, 1 node | lifecycle, 1k seats — refused late/boundary/early; violations | 10/0/0; 0 | 8/0/0; 0 |  |
| PostgreSQL, 1 node | hold /s | 723 | 612 | 1.2x slower |
| PostgreSQL, 1 node | release /s | 5.2k | 4.2k | 1.2x slower |
| PostgreSQL, 1 node | cancel /s | 5.4k | 4.7k | 1.2x slower |
| PostgreSQL, 1 node | publish 10 seats p50 ms | 1.47 | 1.65 | 1.1x slower |
| PostgreSQL, 1 node | publish 100 seats p50 ms | 1.97 | 2.12 | 1.1x slower |
| PostgreSQL, 1 node | publish 1k seats p50 ms | 7.26 | 9.61 | 1.3x slower |
| PostgreSQL, 1 node | publish 10k seats p50 ms | 37 | 67 | 1.8x slower |
| PostgreSQL, 1 node | publish 100k seats p50 ms | 475 | 566 | 1.2x slower |
| PostgreSQL, 1 node | Seats available per section — 10k seats | 875 | 1.4k | 1.6x faster |

### S1 conditional → E2 cart

*Expiry copied onto every seat row, or stored once on the cart?*

| Topology | Measurement | S1 conditional | E2 cart | Change |
|---|---|---:|---:|---|
| PostgreSQL, 1 node | race, 10 seats — seats/s | 413 | 212 | 1.9x slower |
| PostgreSQL, 1 node | race, 100 seats — seats/s | 1.1k | 638 | 1.7x slower |
| PostgreSQL, 1 node | race, 1k seats — seats/s | 716 | 672 | 1.1x slower |
| PostgreSQL, 1 node | race, 10k seats — seats/s | 595 | 330 ⏱ 1 timed out at 99% sold | 1.8x slower |
| PostgreSQL, 1 node | lifecycle, 10 seats — refused late/boundary/early; violations | 0/0/0; 0 | 1/0/0; 0 |  |
| PostgreSQL, 1 node | lifecycle, 100 seats — refused late/boundary/early; violations | 2/0/0; 0 | 4/0/0; 0 |  |
| PostgreSQL, 1 node | lifecycle, 1k seats — refused late/boundary/early; violations | 10/0/0; 0 | 9/0/0; 0 |  |
| PostgreSQL, 1 node | hold /s | 723 | 442 | 1.6x slower |
| PostgreSQL, 1 node | release /s | 5.2k | 5.2k | 1.0x faster |
| PostgreSQL, 1 node | cancel /s | 5.4k | 5.6k | 1.0x faster |
| PostgreSQL, 1 node | publish 10 seats p50 ms | 1.47 | 1.41 | 1.0x faster |
| PostgreSQL, 1 node | publish 100 seats p50 ms | 1.97 | 1.92 | 1.0x faster |
| PostgreSQL, 1 node | publish 1k seats p50 ms | 7.26 | 8.65 | 1.2x slower |
| PostgreSQL, 1 node | publish 10k seats p50 ms | 37 | 54 | 1.5x slower |
| PostgreSQL, 1 node | publish 100k seats p50 ms | 475 | 493 | 1.0x slower |
| PostgreSQL, 1 node | Seat map of a section — 10 seats | 18k | 12k | 1.4x slower |
| PostgreSQL, 1 node | Seats available per section — 10 seats | 20k | 11k | 1.9x slower |
| PostgreSQL, 1 node | Seats available for an event — 10 seats | 20k | 11k | 1.8x slower |
| PostgreSQL, 1 node | Seats available per section — 100 seats | 19k | 9.8k | 2.0x slower |
| PostgreSQL, 1 node | Seats available for an event — 100 seats | 19k | 9.9k | 1.9x slower |
| PostgreSQL, 1 node | Seats available per section — 1k seats | 7.7k | 4.2k | 1.8x slower |
| PostgreSQL, 1 node | Seats available for an event — 1k seats | 7.9k | 4.9k | 1.6x slower |
| PostgreSQL, 1 node | Seats available for an event — 10k seats | 1.3k | 860 | 1.5x slower |

### E0 app clock ✗ → S1 conditional

*The application's clock, or the database's?*

| Topology | Measurement | E0 app clock ✗ | S1 conditional | Change |
|---|---|---:|---:|---|
| PostgreSQL, 1 node | race, 10 seats — seats/s | 293 | 413 | 1.4x faster |
| PostgreSQL, 1 node | race, 100 seats — seats/s | 966 | 1.1k | 1.1x faster |
| PostgreSQL, 1 node | race, 1k seats — seats/s | 1.5k | 716 | 2.1x slower |
| PostgreSQL, 1 node | race, 10k seats — seats/s | 629 | 595 | 1.1x slower |
| PostgreSQL, 1 node | lifecycle, 10 seats — refused late/boundary/early; violations | 0/0/0; 4 | 0/0/0; 0 |  |
| PostgreSQL, 1 node | lifecycle, 100 seats — refused late/boundary/early; violations | 3/0/1; 38 | 2/0/0; 0 |  |
| PostgreSQL, 1 node | lifecycle, 1k seats — refused late/boundary/early; violations | 0/0/12; 37 | 10/0/0; 0 |  |
| PostgreSQL, 1 node | hold /s | 696 | 723 | 1.0x faster |
| PostgreSQL, 1 node | release /s | 4.3k | 5.2k | 1.2x faster |
| PostgreSQL, 1 node | cancel /s | 5.4k | 5.4k | 1.0x faster |
| PostgreSQL, 1 node | publish 10 seats p50 ms | 1.47 | 1.47 | 1.0x faster |
| PostgreSQL, 1 node | publish 100 seats p50 ms | 2.00 | 1.97 | 1.0x faster |
| PostgreSQL, 1 node | publish 1k seats p50 ms | 7.48 | 7.26 | 1.0x faster |
| PostgreSQL, 1 node | publish 10k seats p50 ms | 43 | 37 | 1.2x faster |
| PostgreSQL, 1 node | publish 100k seats p50 ms | 482 | 475 | 1.0x faster |

### K0 naive ✗ → S1 conditional

*What does checking the hold at confirmation cost, and what does it prevent?*

| Topology | Measurement | K0 naive ✗ | S1 conditional | Change |
|---|---|---:|---:|---|
| PostgreSQL, 1 node | race, 10 seats — seats/s | 273 | 413 | 1.5x faster |
| PostgreSQL, 1 node | race, 100 seats — seats/s | 933 | 1.1k | 1.2x faster |
| PostgreSQL, 1 node | race, 1k seats — seats/s | 1.2k | 716 | 1.7x slower |
| PostgreSQL, 1 node | race, 10k seats — seats/s | 568 | 595 | 1.0x faster |
| PostgreSQL, 1 node | lifecycle, 10 seats — refused late/boundary/early; violations | 0/0/0; 0 | 0/0/0; 0 |  |
| PostgreSQL, 1 node | lifecycle, 100 seats — refused late/boundary/early; violations | 0/0/0; 6 | 2/0/0; 0 |  |
| PostgreSQL, 1 node | lifecycle, 1k seats — refused late/boundary/early; violations | 0/0/0; 18 | 10/0/0; 0 |  |
| PostgreSQL, 1 node | hold /s | 598 | 723 | 1.2x faster |
| PostgreSQL, 1 node | release /s | 4.3k | 5.2k | 1.2x faster |
| PostgreSQL, 1 node | cancel /s | 3.7k | 5.4k | 1.5x faster |
| PostgreSQL, 1 node | publish 10 seats p50 ms | 1.70 | 1.47 | 1.2x faster |
| PostgreSQL, 1 node | publish 100 seats p50 ms | 2.38 | 1.97 | 1.2x faster |
| PostgreSQL, 1 node | publish 1k seats p50 ms | 10 | 7.26 | 1.4x faster |
| PostgreSQL, 1 node | publish 10k seats p50 ms | 53 | 37 | 1.5x faster |
| PostgreSQL, 1 node | publish 100k seats p50 ms | 576 | 475 | 1.2x faster |

### S1 conditional → K1 window

*What does a guaranteed payment window cost, and what does it remove?*

| Topology | Measurement | S1 conditional | K1 window | Change |
|---|---|---:|---:|---|
| PostgreSQL, 1 node | race, 10 seats — seats/s | 413 | 414 | 1.0x faster |
| PostgreSQL, 1 node | race, 100 seats — seats/s | 1.1k | 1.0k | 1.0x slower |
| PostgreSQL, 1 node | race, 1k seats — seats/s | 716 | 1.3k | 1.8x faster |
| PostgreSQL, 1 node | race, 10k seats — seats/s | 595 | 538 | 1.1x slower |
| PostgreSQL, 1 node | lifecycle, 10 seats — refused late/boundary/early; violations | 0/0/0; 0 | 0/0/0; 0 |  |
| PostgreSQL, 1 node | lifecycle, 100 seats — refused late/boundary/early; violations | 2/0/0; 0 | 0/0/0; 0 |  |
| PostgreSQL, 1 node | lifecycle, 1k seats — refused late/boundary/early; violations | 10/0/0; 0 | 0/0/0; 0 |  |
| PostgreSQL, 1 node | hold /s | 723 | 701 | 1.0x slower |
| PostgreSQL, 1 node | release /s | 5.2k | 5.6k | 1.1x faster |
| PostgreSQL, 1 node | cancel /s | 5.4k | 5.3k | 1.0x slower |
| PostgreSQL, 1 node | publish 10 seats p50 ms | 1.47 | 1.47 | 1.0x faster |
| PostgreSQL, 1 node | publish 100 seats p50 ms | 1.97 | 1.86 | 1.1x faster |
| PostgreSQL, 1 node | publish 1k seats p50 ms | 7.26 | 7.31 | 1.0x slower |
| PostgreSQL, 1 node | publish 10k seats p50 ms | 37 | 46 | 1.2x slower |
| PostgreSQL, 1 node | publish 100k seats p50 ms | 475 | 464 | 1.0x faster |

### S1 conditional → L1 claims

*Pre-created per-event seat rows, or claims created on hold?*

| Topology | Measurement | S1 conditional | L1 claims | Change |
|---|---|---:|---:|---|
| PostgreSQL, 1 node | race, 10 seats — seats/s | 413 | 438 | 1.1x faster |
| PostgreSQL, 1 node | race, 100 seats — seats/s | 1.1k | 1.1k | 1.0x faster |
| PostgreSQL, 1 node | race, 1k seats — seats/s | 716 | 1.9k | 2.6x faster |
| PostgreSQL, 1 node | race, 10k seats — seats/s | 595 | 716 | 1.2x faster |
| PostgreSQL, 1 node | lifecycle, 10 seats — refused late/boundary/early; violations | 0/0/0; 0 | 0/0/0; 0 |  |
| PostgreSQL, 1 node | lifecycle, 100 seats — refused late/boundary/early; violations | 2/0/0; 0 | 8/0/0; 0 |  |
| PostgreSQL, 1 node | lifecycle, 1k seats — refused late/boundary/early; violations | 10/0/0; 0 | 8/0/0; 0 |  |
| PostgreSQL, 1 node | hold /s | 723 | 713 | 1.0x slower |
| PostgreSQL, 1 node | release /s | 5.2k | 6.1k | 1.2x faster |
| PostgreSQL, 1 node | cancel /s | 5.4k | 5.7k | 1.1x faster |
| PostgreSQL, 1 node | publish 10 seats p50 ms | 1.47 | 1.00 | 1.5x faster |
| PostgreSQL, 1 node | publish 100 seats p50 ms | 1.97 | 1.10 | 1.8x faster |
| PostgreSQL, 1 node | publish 1k seats p50 ms | 7.26 | 1.10 | 6.6x faster |
| PostgreSQL, 1 node | publish 10k seats p50 ms | 37 | 1.23 | 29.8x faster |
| PostgreSQL, 1 node | publish 100k seats p50 ms | 475 | 0.68 | 700.9x faster |
| PostgreSQL, 1 node | Seat map of a section — 10 seats | 18k | 35k | 2.0x faster |
| PostgreSQL, 1 node | Seat map of a section — 100 seats | 18k | 37k | 2.1x faster |

### S1 conditional → L2 document

*Seat rows, or an embedded section document?*

| Topology | Measurement | S1 conditional | L2 document | Change |
|---|---|---:|---:|---|
| PostgreSQL, 1 node | race, 10 seats — seats/s | 413 | 21.0 | 19.7x slower |
| PostgreSQL, 1 node | race, 100 seats — seats/s | 1.1k | 19.8 | 54.6x slower |
| PostgreSQL, 1 node | race, 1k seats — seats/s | 716 | 34.9 | 20.5x slower |
| PostgreSQL, 1 node | race, 10k seats — seats/s | 595 | 57.0 ⏱ 1 timed out at 17% sold | 10.4x slower |
| PostgreSQL, 1 node | lifecycle, 10 seats — refused late/boundary/early; violations | 0/0/0; 0 | 0/0/0; 0 |  |
| PostgreSQL, 1 node | lifecycle, 100 seats — refused late/boundary/early; violations | 2/0/0; 0 | 2/0/0; 0 |  |
| PostgreSQL, 1 node | lifecycle, 1k seats — refused late/boundary/early; violations | 10/0/0; 0 | 7/0/0; 0 |  |
| PostgreSQL, 1 node | hold /s | 723 | 367 | 2.0x slower |
| PostgreSQL, 1 node | release /s | 5.2k | 2.9k | 1.8x slower |
| PostgreSQL, 1 node | cancel /s | 5.4k | 1.2k | 4.6x slower |
| PostgreSQL, 1 node | publish 10 seats p50 ms | 1.47 | 5.53 | 3.8x slower |
| PostgreSQL, 1 node | publish 100 seats p50 ms | 1.97 | 5.42 | 2.8x slower |
| PostgreSQL, 1 node | publish 1k seats p50 ms | 7.26 | 5.30 | 1.4x faster |
| PostgreSQL, 1 node | publish 10k seats p50 ms | 37 | 3.40 | 10.8x faster |
| PostgreSQL, 1 node | publish 100k seats p50 ms | 475 | 6.10 | 77.9x faster |
| PostgreSQL, 1 node | Seat map of a section — 10 seats | 18k | 31k | 1.8x faster |
| PostgreSQL, 1 node | Seats available per section — 10 seats | 20k | 31k | 1.6x faster |
| PostgreSQL, 1 node | Seats available for an event — 10 seats | 20k | 29k | 1.5x faster |
| PostgreSQL, 1 node | Seats available per section — 100 seats | 19k | 29k | 1.5x faster |
| PostgreSQL, 1 node | Seats available per section — 1k seats | 7.7k | 12k | 1.5x faster |
| PostgreSQL, 1 node | Seat map of a section — 10k seats | 16k | 9.7k | 1.6x slower |
| PostgreSQL, 1 node | Seats available per section — 10k seats | 875 | 1.4k | 1.6x faster |


## Measurement hazard: CPU-quota throttling

Cells show *throttled periods / periods (throttled time)*: for the **database** container(s) over
the whole cell (summed across nodes), and for the **client** per phase.

### PostgreSQL, 1 node

| Design | database, whole cell | client lifecycle#1 | client race#1 | client read | client write:cancel | client write:hold | client write:publish | client write:release |
|---|---:|---:|---:|---:|---:|---:|---:|---:|
| S0 check/RC ✗ | 659/1463 (224s) | 0/412 (0.0s) | 0/341 (0.0s) | 1/452 (0.0s) | 0/6 (0.0s) | 0/20 (0.0s) | 0/14 (0.0s) | 0/1 (0.0s) |
| S1 conditional | 685/1103 (241s) | 0/382 (0.0s) | 0/191 (0.0s) | 14/452 (0.3s) | 0/5 (0.0s) | 0/19 (0.0s) | 0/14 (0.0s) | — |
| S2 lock | 709/1129 (254s) | 0/386 (0.0s) | 0/209 (0.0s) | 5/452 (0.0s) | 0/7 (0.0s) | 0/23 (0.0s) | 0/15 (0.0s) | — |
| S3 nowait | 673/1061 (232s) | 0/345 (0.0s) | 0/174 (0.0s) | 30/452 (2.2s) | 0/7 (0.0s) | 0/24 (0.0s) | 0/14 (0.0s) | — |
| S4 check/SER | 705/1145 (243s) | 0/399 (0.0s) | 0/206 (0.0s) | 7/452 (0.0s) | 0/7 (0.0s) | 0/26 (0.0s) | 0/14 (0.0s) | — |
| E0 app clock ✗ | 680/1076 (238s) | 0/366 (0.0s) | 0/184 (0.0s) | 4/452 (0.2s) | 0/6 (0.0s) | 0/19 (0.0s) | 0/13 (0.0s) | — |
| E1 sweeper | 662/1072 (218s) | 0/373 (0.0s) | 0/166 (0.0s) | 13/451 (0.6s) | 0/6 (0.0s) | 0/21 (0.0s) | 0/16 (0.0s) | 0/1 (0.0s) |
| E2 cart | 845/1180 (355s) | 0/306 (0.0s) | 0/338 (0.0s) | 23/451 (0.9s) | 0/5 (0.0s) | 1/29 (0.0s) | 0/14 (0.0s) | — |
| K0 naive ✗ | 701/1101 (250s) | 0/359 (0.0s) | 0/203 (0.0s) | 15/452 (1.1s) | 0/8 (0.0s) | 0/22 (0.0s) | 0/16 (0.0s) | — |
| K1 window | 705/1113 (254s) | 0/377 (0.0s) | 0/207 (0.0s) | 1/452 (0.0s) | 0/6 (0.0s) | 0/19 (0.0s) | 0/13 (0.0s) | — |
| L1 claims | 640/988 (217s) | 0/323 (0.0s) | 0/159 (0.0s) | 31/452 (0.6s) | 0/5 (0.0s) | 0/18 (0.0s) | 0/6 (0.0s) | — |
| L2 document | 506/1869 (142s) | 0/380 (0.0s) | 0/914 (0.0s) | 27/452 (0.9s) | 24/25 (1.7s) | 0/34 (0.0s) | 0/32 (0.0s) | — |


---

Plans for every read and write statement are in `results/dc3-pg-all/<topology>/plans/<design>.txt`
(`EXPLAIN (ANALYZE, BUFFERS)` on PostgreSQL, `EXPLAIN (ANALYZE, DIST)` on YugabyteDB; write
statements inside rolled-back transactions).

## Provenance: code, data and analyses

This file is **generated from the measurements** and contains no interpretation.
Conclusions live in separate signed analyses, so several analysts — including
different AI models — can read the same numbers and each record what they make of
them. Where two analyses disagree, the disagreement is a finding and stays visible.

| | |
|---|---|
| Run id | `devchecks/dc3-pg-all` |
| Result files | 12 |
| **Inputs digest** | `917778be094b098a` |
| Repository version | `repo/platform-lock-not-available-3-g1b2801d-dirty` (`1b2801d6d10f`), ⚠️ **dirty** — uncommitted changes; not reproducible from the commit alone |

Quote the digest in any analysis: it identifies the data. The repository version
identifies the code — check it out (`git checkout <commit>`) to reproduce the run, or
`git diff <commit>` to see what has changed in the study since.

### Analyses of this run

**None yet.** This run has measurements but no interpretation.

### Adding an analysis

Copy `reports/analyses/TEMPLATE.md`, fill in the frontmatter, and regenerate this report.
If an analysis already exists for digest `917778be094b098a`, read it first and write a new one only to
**disagree, extend, or bring a different perspective**. Never edit another analyst's file.
