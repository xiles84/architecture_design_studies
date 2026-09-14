# Study 03 — results: reserved seating (venue → event → seat)

Generated from `results/dc14b-yb3-race2`. Every number comes from exactly one JSON file in that
directory. **This file contains measurements and no interpretation**; conclusions live in
signed analyses, indexed at the end.

## TL;DR — measured facts

*Selected from the tables below by fixed rules; no interpretation. Read a signed analysis for
what these facts mean.*

- **Cells:** 8 across 1 topologies; 0 failed.
- **Negative controls:** 1 of 1 fired as designed to.
- ❌ **Invariant violations in designs meant to be correct:** S1 conditional on YugabyteDB, 3 nodes (RF=3) (23 early rejections (23 transient)); E1 sweeper on YugabyteDB, 3 nodes (RF=3) (15 early rejections (15 transient)); K1 window on YugabyteDB, 3 nodes (RF=3) (28 early rejections (28 transient)); L1 claims on YugabyteDB, 3 nodes (RF=3) (38 early rejections (38 transient)); L3 sharded on YugabyteDB, 3 nodes (RF=3) (33 early rejections (33 transient)).
- **S1r confirmation retries, YugabyteDB, 3 nodes (RF=3)** (race, all tiers and trials): 20 retried, 20 of them sold.
- **Early rejections in the race, correct designs, YugabyteDB, 3 nodes (RF=3):** S1 conditional 23 (23 transient), E1 sweeper 15 (15 transient), K1 window 28 (28 transient), L1 claims 38 (38 transient), L3 sharded 33 (33 transient).
- **Race, YugabyteDB, 3 nodes (RF=3)** (seats held/s; designs meant to be correct, without violations):
  - 10 seats: highest E1 sweeper 34.4/s, lowest L1 claims 13.0/s
  - 100 seats: highest S1r retry 95.0/s, lowest L1 claims 66.3/s
  - 1k seats: highest S1r retry 65.7/s, lowest S1r retry 65.7/s
  - 10k seats: highest S1r retry 37.1/s, lowest S1r retry 37.1/s (1 of 1 timed out)

## What was measured

|  |  |
|---|---|
| Run id | `devchecks/dc14b-yb3-race2` |
| Environment | [`host-zenbook-ux5406sa`](../../../docs/environments/host-zenbook-ux5406sa.md) |
| Scale | `tiny` — 5 venues, 84 events, 5741 tickets sold, 131 live and 34 expired holds at load |
| Events (kind_seats: count) | catalogue_10: 40 · catalogue_100: 8 · catalogue_1000: 2 · catalogue_10000: 1 · lifecycle_10: 4 · lifecycle_100: 2 · lifecycle_1000: 1 · race_10: 20 · race_100: 4 · race_1000: 1 · race_10000: 1 |
| Seed | 42 — identical data in every cell |
| Reads / isolated writes | 8 workers, 5s measured + 2s warmup, 1 trial(s) |
| Race | 32 buyers per event, real hold TTL 40m0s, 10% of holds deferred until the crowd has gone, timeout 1m0s, tier budget 3m0s, give up after 20 conflicts |

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
| S0 check/RC ✗ | ✅ 51/51 |
| S1 conditional | ✅ 51/51 |
| E0 app clock ✗ | ✅ 51/51 |
| E1 sweeper | ✅ 51/51 (after its sweeper released 65 expired seats) |
| K1 window | ✅ 51/51 |
| S1r retry | ✅ 51/51 |
| L1 claims | ✅ 51/51 |
| L3 sharded | ✅ 51/51 |

## Negative controls

Deliberately wrong designs (and E1 with its sweeper stopped). Each must be seen to break the
invariant it exists to test; otherwise the absence of violations elsewhere in that experiment
is not evidence of correctness.

| Control | Experiment | Topology | Fired? | What was found |
|---|---|---|---|---|
| S0 check/RC ✗ | race | YugabyteDB, 3 nodes (RF=3) | ✅ fired | 9 deferred confirmations refused; 646 thefts, 355 early rejections (19 transient) |

## The race — a hot drop with seat choice

Seats held per second, from the gate to the last hold granted. Conflicts and seat-map reads are
per granted hold. Deferred confirmations kept a real 40-minute hold through the crowd.

### YugabyteDB, 3 nodes (RF=3)

| Design | 10 seats | 100 seats | 1k seats | 10k seats |
|---|---:|---:|---:|---:|
| S0 check/RC ✗ | ~~41.6 (2 trials, spread 52%)~~ ❌ 159 thefts, 97 early rejections (0 transient); ❌ 26 of 41 deferred confirmations refused; 706 errors | ~~101 (2 trials, spread 11%)~~ ❌ 155 thefts, 80 early rejections (1 transient); ❌ 22 of 51 deferred confirmations refused; 285 errors | ~~79.3 (2 trials, spread 12%)~~ ❌ 195 thefts, 98 early rejections (2 transient); ❌ 39 of 104 deferred confirmations refused; 302 errors | ~~40.5 (2 trials, spread 18%)~~ ❌ 137 thefts, 80 early rejections (16 transient); ❌ 17 of 210 deferred confirmations refused; ⏱ 2 timed out at 25% sold; 137 errors |
| S1 conditional | 29.4 (2 trials, spread 32%) | ~~76.2 (2 trials, spread 7%)~~ ❌ 2 early rejections (2 transient) | ~~67.9 (2 trials, spread 24%)~~ ❌ 6 early rejections (6 transient) | ~~37.6 (2 trials, spread 2%)~~ ❌ 15 early rejections (15 transient); ⏱ 2 timed out at 23% sold |
| E0 app clock ✗ | ~~35.6 (2 trials, spread 9%)~~ ❌ 2 early rejections (2 transient) | ~~76.8 (2 trials, spread 1%)~~ ❌ 1 early rejections (1 transient) | ~~76.9 (2 trials, spread 13%)~~ ❌ 5 early rejections (5 transient) | ~~40.1 (2 trials, spread 6%)~~ ❌ 11 early rejections (11 transient); ⏱ 2 timed out at 25% sold |
| E1 sweeper | 34.4 (2 trials, spread 3%) | ~~79.6 (2 trials, spread 30%)~~ ❌ 3 early rejections (3 transient) | ~~65.5 (2 trials, spread 1%)~~ ❌ 1 early rejections (1 transient) | ~~37.0 (2 trials, spread 2%)~~ ❌ 11 early rejections (11 transient); ⏱ 2 timed out at 22% sold |
| K1 window | 34.0 (2 trials, spread 22%) | 72.9 (2 trials, spread 7%) | ~~61.8 (2 trials, spread 19%)~~ ❌ 10 early rejections (10 transient) | ~~37.3 (2 trials, spread 1%)~~ ❌ 18 early rejections (18 transient); ⏱ 2 timed out at 23% sold |
| S1r retry | 32.3 (2 trials, spread 6%) | 95.0 (2 trials, spread 10%) | 65.7 (2 trials, spread 6%) | 37.1 (2 trials, spread 1%) ⏱ 2 timed out at 23% sold |
| L1 claims | 13.0 (2 trials, spread 42%) | 66.3 (2 trials, spread 3%) | ~~86.0 (2 trials, spread 12%)~~ ❌ 4 early rejections (4 transient) | ~~77.1 (2 trials, spread 9%)~~ ❌ 34 early rejections (34 transient); ⏱ 2 timed out at 49% sold |
| L3 sharded | ~~10.1 (2 trials, spread 21%)~~ ❌ 1 early rejections (1 transient) | ~~27.1 (2 trials, spread 18%)~~ ❌ 3 early rejections (3 transient) | ~~29.5 (2 trials, spread 10%)~~ ❌ 3 early rejections (3 transient) | ~~48.6 (2 trials, spread 2%)~~ ❌ 26 early rejections (26 transient); ⏱ 2 timed out at 29% sold |

<details><summary>Race detail — YugabyteDB, 3 nodes (RF=3)</summary>

| Design | Tier | Conflicts/hold | Map reads/hold | Engine retries | Gave up | Hold p50/p99 ms | Confirm p99 ms | Deferred ok | Early rej. (transient) | Confirm retries run / sold | Final buyer sold |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| S0 check/RC ✗ | 10 seats | 0.93 | 5.12 | 105 | 0 | 40 / 463 | 351 | 15/41 | 97 (0) | 0 / 0 | 0 |
| S0 check/RC ✗ | 100 seats | 0.64 | 3.20 | 55 | 0 | 48 / 285 | 381 | 29/51 | 80 (1) | 0 / 0 | 0 |
| S0 check/RC ✗ | 1k seats | 0.45 | 2.15 | 45 | 0 | 99 / 406 | 389 | 65/104 | 98 (2) | 0 / 0 | 0 |
| S0 check/RC ✗ | 10k seats | 0.20 | 1.29 | 19 | 0 | 200 / 559 | 670 | 193/210 | 80 (16) | 0 / 0 | 0 |
| S1 conditional | 10 seats | 6.98 | 14.54 | 0 | 0 | 38 / 335 | 672 | 38/38 | 0 (0) | 0 / 0 | 0 |
| S1 conditional | 100 seats | 2.65 | 4.79 | 0 | 0 | 51 / 319 | 620 | 42/42 | 2 (2) | 0 / 0 | 0 |
| S1 conditional | 1k seats | 0.92 | 3.02 | 0 | 0 | 62 / 266 | 310 | 103/103 | 6 (6) | 0 / 0 | 0 |
| S1 conditional | 10k seats | 0.40 | 1.48 | 0 | 0 | 187 / 567 | 1072 | 174/174 | 15 (15) | 0 / 0 | 0 |
| E0 app clock ✗ | 10 seats | 6.01 | 13.49 | 0 | 0 | 35 / 276 | 726 | 24/24 | 2 (2) | 0 / 0 | 0 |
| E0 app clock ✗ | 100 seats | 2.79 | 5.32 | 0 | 0 | 47 / 570 | 258 | 37/37 | 1 (1) | 0 / 0 | 0 |
| E0 app clock ✗ | 1k seats | 0.99 | 2.81 | 0 | 0 | 64 / 223 | 328 | 89/89 | 5 (5) | 0 / 0 | 0 |
| E0 app clock ✗ | 10k seats | 0.27 | 1.37 | 0 | 0 | 144 / 503 | 1185 | 184/184 | 11 (11) | 0 / 0 | 0 |
| E1 sweeper | 10 seats | 7.11 | 14.81 | 0 | 0 | 31 / 371 | 632 | 26/26 | 0 (0) | 0 / 0 | 0 |
| E1 sweeper | 100 seats | 2.55 | 4.51 | 0 | 0 | 39 / 308 | 212 | 28/28 | 3 (3) | 0 / 0 | 0 |
| E1 sweeper | 1k seats | 0.94 | 3.30 | 0 | 0 | 84 / 312 | 407 | 80/80 | 1 (1) | 0 / 0 | 0 |
| E1 sweeper | 10k seats | 0.38 | 1.40 | 0 | 0 | 177 / 497 | 679 | 182/182 | 11 (11) | 0 / 0 | 0 |
| K1 window | 10 seats | 7.25 | 14.84 | 0 | 0 | 38 / 249 | 620 | 28/28 | 0 (0) | 0 / 0 | 0 |
| K1 window | 100 seats | 2.57 | 4.90 | 0 | 0 | 45 / 506 | 305 | 39/39 | 0 (0) | 0 / 0 | 0 |
| K1 window | 1k seats | 1.23 | 2.89 | 0 | 0 | 91 / 320 | 371 | 84/84 | 10 (10) | 0 / 0 | 0 |
| K1 window | 10k seats | 0.32 | 1.36 | 0 | 0 | 172 / 505 | 1366 | 178/178 | 18 (18) | 0 / 0 | 0 |
| S1r retry | 10 seats | 6.70 | 13.80 | 0 | 0 | 38 / 288 | 698 | 30/30 | 0 (0) | 0 / 0 | 0 |
| S1r retry | 100 seats | 2.40 | 4.74 | 0 | 0 | 38 / 222 | 286 | 40/40 | 0 (0) | 1 / 1 | 0 |
| S1r retry | 1k seats | 1.01 | 3.36 | 0 | 0 | 84 / 316 | 376 | 106/106 | 0 (0) | 8 / 8 | 0 |
| S1r retry | 10k seats | 0.35 | 1.38 | 0 | 0 | 179 / 504 | 1158 | 158/158 | 0 (0) | 11 / 11 | 0 |
| L1 claims | 10 seats | 7.31 | 15.36 | 4 | 0 | 52 / 771 | 794 | 29/29 | 0 (0) | 0 / 0 | 0 |
| L1 claims | 100 seats | 2.30 | 4.45 | 3 | 0 | 39 / 326 | 594 | 35/35 | 0 (0) | 0 / 0 | 0 |
| L1 claims | 1k seats | 0.89 | 2.84 | 0 | 0 | 65 / 343 | 444 | 103/103 | 4 (4) | 0 / 0 | 0 |
| L1 claims | 10k seats | 0.25 | 1.49 | 0 | 0 | 94 / 366 | 687 | 407/407 | 34 (34) | 0 / 0 | 0 |
| L3 sharded | 10 seats | 5.15 | 10.89 | 0 | 0 | 135 / 627 | 1928 | 22/22 | 1 (1) | 0 / 0 | 0 |
| L3 sharded | 100 seats | 2.12 | 3.77 | 0 | 0 | 125 / 612 | 569 | 38/38 | 3 (3) | 0 / 0 | 0 |
| L3 sharded | 1k seats | 0.71 | 2.72 | 0 | 0 | 132 / 526 | 659 | 86/86 | 3 (3) | 0 / 0 | 0 |
| L3 sharded | 10k seats | 0.16 | 1.27 | 0 | 0 | 159 / 492 | 1423 | 235/235 | 26 (26) | 0 / 0 | 0 |

</details>

## Reads

Each question in isolation, fresh key per execution, ops/s.

### YugabyteDB, 3 nodes (RF=3)

| Question | S0 check/RC ✗ | S1 conditional | E0 app clock ✗ | E1 sweeper | K1 window | S1r retry | L1 claims | L3 sharded |
|---|---:|---:|---:|---:|---:|---:|---:|---:|

## Isolated writes, publishing and storage

### YugabyteDB, 3 nodes (RF=3)

| Design | Hold+confirm /s | Release /s | Refund /s | Publish p50 ms by size | Storage | Audit |
|---|---:|---:|---:|---:|---:|---|
| S0 check/RC ✗ | — | — | — |  | n/a on YugabyteDB | consistent |
| S1 conditional | — | — | — |  | n/a on YugabyteDB | consistent |
| E0 app clock ✗ | — | — | — |  | n/a on YugabyteDB | consistent |
| E1 sweeper | — | — | — |  | n/a on YugabyteDB | consistent |
| K1 window | — | — | — |  | n/a on YugabyteDB | consistent |
| S1r retry | — | — | — |  | n/a on YugabyteDB | consistent |
| L1 claims | — | — | — |  | n/a on YugabyteDB | consistent |
| L3 sharded | — | — | — |  | n/a on YugabyteDB | consistent |


## Controlled pairs — one decision at a time

### S1 conditional → E1 sweeper

*Expiry judged at use time, or made to take effect by a sweeper?*

| Topology | Measurement | S1 conditional | E1 sweeper | Change |
|---|---|---:|---:|---|
| YugabyteDB, 3 nodes (RF=3) | race, 10 seats — seats/s | 29.4 (2 trials, spread 32%) | 34.4 (2 trials, spread 3%) | 1.2x faster |
| YugabyteDB, 3 nodes (RF=3) | race, 100 seats — seats/s | ~~76.2 (2 trials, spread 7%)~~ ❌ 2 early rejections (2 transient) | ~~79.6 (2 trials, spread 30%)~~ ❌ 3 early rejections (3 transient) | 1.0x faster — ❌ a design with violations has no valid speed |
| YugabyteDB, 3 nodes (RF=3) | race, 1k seats — seats/s | ~~67.9 (2 trials, spread 24%)~~ ❌ 6 early rejections (6 transient) | ~~65.5 (2 trials, spread 1%)~~ ❌ 1 early rejections (1 transient) | 1.0x slower — ❌ a design with violations has no valid speed |
| YugabyteDB, 3 nodes (RF=3) | race, 10k seats — seats/s | ~~37.6 (2 trials, spread 2%)~~ ❌ 15 early rejections (15 transient); ⏱ 2 timed out at 23% sold | ~~37.0 (2 trials, spread 2%)~~ ❌ 11 early rejections (11 transient); ⏱ 2 timed out at 22% sold | 1.0x slower — ❌ a design with violations has no valid speed |

### E0 app clock ✗ → S1 conditional

*The application's clock, or the database's?*

| Topology | Measurement | E0 app clock ✗ | S1 conditional | Change |
|---|---|---:|---:|---|
| YugabyteDB, 3 nodes (RF=3) | race, 10 seats — seats/s | ~~35.6 (2 trials, spread 9%)~~ ❌ 2 early rejections (2 transient) | 29.4 (2 trials, spread 32%) | 1.2x slower — ❌ a design with violations has no valid speed |
| YugabyteDB, 3 nodes (RF=3) | race, 100 seats — seats/s | ~~76.8 (2 trials, spread 1%)~~ ❌ 1 early rejections (1 transient) | ~~76.2 (2 trials, spread 7%)~~ ❌ 2 early rejections (2 transient) | 1.0x slower — ❌ a design with violations has no valid speed |
| YugabyteDB, 3 nodes (RF=3) | race, 1k seats — seats/s | ~~76.9 (2 trials, spread 13%)~~ ❌ 5 early rejections (5 transient) | ~~67.9 (2 trials, spread 24%)~~ ❌ 6 early rejections (6 transient) | 1.1x slower — ❌ a design with violations has no valid speed |
| YugabyteDB, 3 nodes (RF=3) | race, 10k seats — seats/s | ~~40.1 (2 trials, spread 6%)~~ ❌ 11 early rejections (11 transient); ⏱ 2 timed out at 25% sold | ~~37.6 (2 trials, spread 2%)~~ ❌ 15 early rejections (15 transient); ⏱ 2 timed out at 23% sold | 1.1x slower — ❌ a design with violations has no valid speed |

### S1 conditional → K1 window

*What does a guaranteed payment window cost, and what does it remove?*

| Topology | Measurement | S1 conditional | K1 window | Change |
|---|---|---:|---:|---|
| YugabyteDB, 3 nodes (RF=3) | race, 10 seats — seats/s | 29.4 (2 trials, spread 32%) | 34.0 (2 trials, spread 22%) | 1.2x faster |
| YugabyteDB, 3 nodes (RF=3) | race, 100 seats — seats/s | ~~76.2 (2 trials, spread 7%)~~ ❌ 2 early rejections (2 transient) | 72.9 (2 trials, spread 7%) | 1.0x slower — ❌ a design with violations has no valid speed |
| YugabyteDB, 3 nodes (RF=3) | race, 1k seats — seats/s | ~~67.9 (2 trials, spread 24%)~~ ❌ 6 early rejections (6 transient) | ~~61.8 (2 trials, spread 19%)~~ ❌ 10 early rejections (10 transient) | 1.1x slower — ❌ a design with violations has no valid speed |
| YugabyteDB, 3 nodes (RF=3) | race, 10k seats — seats/s | ~~37.6 (2 trials, spread 2%)~~ ❌ 15 early rejections (15 transient); ⏱ 2 timed out at 23% sold | ~~37.3 (2 trials, spread 1%)~~ ❌ 18 early rejections (18 transient); ⏱ 2 timed out at 23% sold | 1.0x slower — ❌ a design with violations has no valid speed |

### S1 conditional → L1 claims

*Pre-created per-event seat rows, or claims created on hold?*

| Topology | Measurement | S1 conditional | L1 claims | Change |
|---|---|---:|---:|---|
| YugabyteDB, 3 nodes (RF=3) | race, 10 seats — seats/s | 29.4 (2 trials, spread 32%) | 13.0 (2 trials, spread 42%) | 2.3x slower |
| YugabyteDB, 3 nodes (RF=3) | race, 100 seats — seats/s | ~~76.2 (2 trials, spread 7%)~~ ❌ 2 early rejections (2 transient) | 66.3 (2 trials, spread 3%) | 1.1x slower — ❌ a design with violations has no valid speed |
| YugabyteDB, 3 nodes (RF=3) | race, 1k seats — seats/s | ~~67.9 (2 trials, spread 24%)~~ ❌ 6 early rejections (6 transient) | ~~86.0 (2 trials, spread 12%)~~ ❌ 4 early rejections (4 transient) | 1.3x faster — ❌ a design with violations has no valid speed |
| YugabyteDB, 3 nodes (RF=3) | race, 10k seats — seats/s | ~~37.6 (2 trials, spread 2%)~~ ❌ 15 early rejections (15 transient); ⏱ 2 timed out at 23% sold | ~~77.1 (2 trials, spread 9%)~~ ❌ 34 early rejections (34 transient); ⏱ 2 timed out at 49% sold | 2.1x faster — ❌ a design with violations has no valid speed |

### S1 conditional → L3 sharded

*(YugabyteDB) all of an event's seats on one tablet, or spread by section?*

| Topology | Measurement | S1 conditional | L3 sharded | Change |
|---|---|---:|---:|---|
| YugabyteDB, 3 nodes (RF=3) | race, 10 seats — seats/s | 29.4 (2 trials, spread 32%) | ~~10.1 (2 trials, spread 21%)~~ ❌ 1 early rejections (1 transient) | 2.9x slower — ❌ a design with violations has no valid speed |
| YugabyteDB, 3 nodes (RF=3) | race, 100 seats — seats/s | ~~76.2 (2 trials, spread 7%)~~ ❌ 2 early rejections (2 transient) | ~~27.1 (2 trials, spread 18%)~~ ❌ 3 early rejections (3 transient) | 2.8x slower — ❌ a design with violations has no valid speed |
| YugabyteDB, 3 nodes (RF=3) | race, 1k seats — seats/s | ~~67.9 (2 trials, spread 24%)~~ ❌ 6 early rejections (6 transient) | ~~29.5 (2 trials, spread 10%)~~ ❌ 3 early rejections (3 transient) | 2.3x slower — ❌ a design with violations has no valid speed |
| YugabyteDB, 3 nodes (RF=3) | race, 10k seats — seats/s | ~~37.6 (2 trials, spread 2%)~~ ❌ 15 early rejections (15 transient); ⏱ 2 timed out at 23% sold | ~~48.6 (2 trials, spread 2%)~~ ❌ 26 early rejections (26 transient); ⏱ 2 timed out at 29% sold | 1.3x faster — ❌ a design with violations has no valid speed |

### S1 conditional → S1r retry

*Does retrying a short confirmation once remove transient refusals, and what does it cost?*

| Topology | Measurement | S1 conditional | S1r retry | Change |
|---|---|---:|---:|---|
| YugabyteDB, 3 nodes (RF=3) | race, 10 seats — seats/s | 29.4 (2 trials, spread 32%) | 32.3 (2 trials, spread 6%) | 1.1x faster |
| YugabyteDB, 3 nodes (RF=3) | race, 100 seats — seats/s | ~~76.2 (2 trials, spread 7%)~~ ❌ 2 early rejections (2 transient) | 95.0 (2 trials, spread 10%) | 1.2x faster — ❌ a design with violations has no valid speed |
| YugabyteDB, 3 nodes (RF=3) | race, 1k seats — seats/s | ~~67.9 (2 trials, spread 24%)~~ ❌ 6 early rejections (6 transient) | 65.7 (2 trials, spread 6%) | 1.0x slower — ❌ a design with violations has no valid speed |
| YugabyteDB, 3 nodes (RF=3) | race, 10k seats — seats/s | ~~37.6 (2 trials, spread 2%)~~ ❌ 15 early rejections (15 transient); ⏱ 2 timed out at 23% sold | 37.1 (2 trials, spread 1%) ⏱ 2 timed out at 23% sold | 1.0x slower — ❌ a design with violations has no valid speed |


## Measurement hazard: CPU-quota throttling

Cells show *throttled periods / periods (throttled time)*: for the **database** container(s) over
the whole cell (summed across nodes), and for the **client** per phase.

### YugabyteDB, 3 nodes (RF=3)

| Design | database, whole cell | client race#1 | client race#2 |
|---|---:|---:|---:|
| S0 check/RC ✗ | 1923/7575 (550s) | 0/950 (0.0s) | 0/914 (0.0s) |
| S1 conditional | 1967/7597 (572s) | 0/951 (0.0s) | 0/942 (0.0s) |
| E0 app clock ✗ | 1966/7390 (583s) | 0/922 (0.0s) | 0/915 (0.0s) |
| E1 sweeper | 1980/7676 (590s) | 0/976 (0.0s) | 0/931 (0.0s) |
| K1 window | 2005/7815 (591s) | 0/950 (0.0s) | 0/964 (0.0s) |
| S1r retry | 1919/7589 (566s) | 0/936 (0.0s) | 0/931 (0.0s) |
| L1 claims | 2291/7975 (425s) | 0/1017 (0.0s) | 0/1028 (0.0s) |
| L3 sharded | 3465/11986 (1236s) | 0/1677 (0.0s) | 0/1709 (0.0s) |


---

Plans for every read and write statement are in `results/dc14b-yb3-race2/<topology>/plans/<design>.txt`
(`EXPLAIN (ANALYZE, BUFFERS)` on PostgreSQL, `EXPLAIN (ANALYZE, DIST)` on YugabyteDB; write
statements inside rolled-back transactions).

## Provenance: code, data and analyses

This file is **generated from the measurements** and contains no interpretation.
Conclusions live in separate signed analyses, so several analysts — including
different AI models — can read the same numbers and each record what they make of
them. Where two analyses disagree, the disagreement is a finding and stays visible.

| | |
|---|---|
| Run id | `devchecks/dc14b-yb3-race2` |
| Result files | 8 |
| **Inputs digest** | `d7dce151a5fc574c` |
| Repository version | `study-03/v0.1-handoff-amendment-01-7-gf43fc0f-dirty` (`f43fc0f2c1e1`), ⚠️ **dirty** — uncommitted changes; not reproducible from the commit alone |

Quote the digest in any analysis: it identifies the data. The repository version
identifies the code — check it out (`git checkout <commit>`) to reproduce the run, or
`git diff <commit>` to see what has changed in the study since.

### Analyses of this run

**None yet.** This run has measurements but no interpretation.

### Adding an analysis

Copy `reports/analyses/TEMPLATE.md`, fill in the frontmatter, and regenerate this report.
If an analysis already exists for digest `d7dce151a5fc574c`, read it first and write a new one only to
**disagree, extend, or bring a different perspective**. Never edit another analyst's file.
