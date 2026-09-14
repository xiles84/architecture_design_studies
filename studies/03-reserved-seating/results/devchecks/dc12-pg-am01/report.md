# Study 03 — results: reserved seating (venue → event → seat)

Generated from `results/dc12-pg-am01`. Every number comes from exactly one JSON file in that
directory. **This file contains measurements and no interpretation**; conclusions live in
signed analyses, indexed at the end.

## TL;DR — measured facts

*Selected from the tables below by fixed rules; no interpretation. Read a signed analysis for
what these facts mean.*

- **Cells:** 5 across 1 topologies; 0 failed.
- **Negative controls:** 0 of 0 fired as designed to.
- **Invariant violations in designs meant to be correct:** none observed.
- **Refused confirmations, S1 conditional, PostgreSQL, 1 node** (lifecycle, all tiers): 9 late, 0 boundary, 0 early (0 transient); 556 confirmed.
- **Refused confirmations, S1r retry, PostgreSQL, 1 node** (lifecycle, all tiers): 7 late, 0 boundary, 0 early (0 transient); 567 confirmed; confirmations retried 7, of which sold 0.
- **Refused confirmations, K1 window, PostgreSQL, 1 node** (lifecycle, all tiers): 0 late, 0 boundary, 0 early (0 transient); 568 confirmed.
- **S1r confirmation retries, PostgreSQL, 1 node** (race, all tiers and trials): 0 retried, 0 of them sold.
- **Race, PostgreSQL, 1 node** (seats held/s; designs meant to be correct, without violations):
  - 10 seats: highest S1 conditional 218/s, lowest E2 cart 75.4/s
  - 100 seats: highest S1 conditional 650/s, lowest E2 cart 368/s
  - 1k seats: highest S1 conditional 1.0k/s, lowest E2 cart 441/s
  - 10k seats: highest S1r retry 277/s, lowest E2 cart 164/s (1 of 5 timed out)
- **Publishing a 100k seats event, PostgreSQL, 1 node** (p50): 667–1287 ms with per-event inventory rows or documents, 1.45 ms without (L1).

## What was measured

|  |  |
|---|---|
| Run id | `devchecks/dc12-pg-am01` |
| Environment | [`host-zenbook-ux5406sa`](../../../docs/environments/host-zenbook-ux5406sa.md) |
| Scale | `tiny` — 5 venues, 84 events, 5741 tickets sold, 131 live and 34 expired holds at load |
| Events (kind_seats: count) | catalogue_10: 40 · catalogue_100: 8 · catalogue_1000: 2 · catalogue_10000: 1 · lifecycle_10: 4 · lifecycle_100: 2 · lifecycle_1000: 1 · race_10: 20 · race_100: 4 · race_1000: 1 · race_10000: 1 |
| Seed | 42 — identical data in every cell |
| Reads / isolated writes | 8 workers, 5s measured + 2s warmup, 1 trial(s) |
| Race | 32 buyers per event, real hold TTL 40m0s, 10% of holds deferred until the crowd has gone, timeout 1m0s, tier budget 3m0s, give up after 20 conflicts |
| Lifecycle — PostgreSQL, 1 node | one human minute = 50ms; TTL 40m0s, payment window (K1) 10m0s, guard G 2m0s, sweeper every 1m0s, outage 45m0s→1h25m0s in the first event of each tier, 25% abandoned (40% of those explicitly), E0 skew 10m0s; tiers 10,100,1000, 6 events per tier, 32 buyers |

**Engines**

| Topology | Version | READ COMMITTED runs as | Client connections to | DB clock − client clock (ms) |
|---|---|---|---|---:|
| PostgreSQL, 1 node | `PostgreSQL 17.11 (Debian 17.11-1.pgdg13+2) on x86_64-pc-linux-gnu, compiled by gcc (Debian 14.2.0-19) 14.2.0, 64-bit` | read committed | 1 node(s) | 0.01 |

> Limits of this environment that bound every number below: the client and the databases
> share one laptop's eight cores; the "3-node" cluster has no real network between nodes; the
> race and the lifecycle are closed loops, so tails are a floor; the lifecycle runs on
> compressed time; E0's clock skew is injected.

## Correctness gate

Six read questions checked against Go-computed truth on a load with sold seats, live holds and
expired holds, plus an audit of the loaded state. A failure aborts the cell.

| Design | PostgreSQL, 1 node |
|---|---|
| S1 conditional | ✅ 51/51 |
| E2 cart | ✅ 51/51 |
| K1 window | ✅ 51/51 |
| S1r retry | ✅ 51/51 |
| L1 claims | ✅ 51/51 |

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
| S1 conditional | 10 seats | 3.5 | 28 | 1+2 | 1 | 0 / 0 / 0 (0) | 0 / 0 | 0 | 45506.8 / 69006.1 | 128 | none |
| S1 conditional | 100 seats | 11.1 | 179 | 35+20 | 8 | 1 / 0 / 0 (0) | 0 / 0 | 0 | 0.3 / 1.0 | 2871 | none |
| S1 conditional | 1k seats | 68.4 | 587 | 93+53 | 16 | 8 / 0 / 0 (0) | 0 / 5 | 0 | 0.4 / 1.0 | 9579 | none |
| E2 cart | 10 seats | 3.1 | 37 | 4+3 | 1 | 0 / 0 / 0 (0) | 0 / 0 | 0 | 44521.2 / 69001.7 | 188 | none |
| E2 cart | 100 seats | 10.7 | 173 | 28+24 | 6 | 0 / 0 / 0 (0) | 0 / 0 | 0 | 0.5 / 0.7 | 2534 | none |
| E2 cart | 1k seats | 68.2 | 575 | 89+55 | 14 | 5 / 0 / 0 (0) | 0 / 5 | 0 | 0.4 / 2.1 | 9245 | none |
| K1 window | 10 seats | 3.1 | 37 | 3+5 | 1 | 0 / 0 / 0 (0) | 0 / 0 | 0 | 45502.7 / 69002.0 | 226 | none |
| K1 window | 100 seats | 9.1 | 155 | 30+15 | 4 | 0 / 0 / 0 (0) | 0 / 0 | 0 | 0.3 / 0.4 | 2657 | none |
| K1 window | 1k seats | 72.2 | 575 | 71+56 | 14 | 0 / 0 / 0 (0) | 0 / 13 | 0 | 0.4 / 1.0 | 8213 | none |
| S1r retry | 10 seats | 2.8 | 38 | 6+2 | 0 | 0 / 0 / 0 (0) | 0 / 0 | 0 | 44509.4 / 68989.8 | 488 | none |
| S1r retry | 100 seats | 19.7 | 162 | 26+18 | 1 | 2 / 0 / 0 (0) | 0 / 0 | 0 | 0.1 / 0.3 | 2600 | none |
| S1r retry | 1k seats | 72.1 | 581 | 90+50 | 14 | 5 / 0 / 0 (0) | 0 / 11 | 0 | 0.4 / 5.0 | 9113 | none |
| L1 claims | 10 seats | 2.7 | 43 | 10+5 | 1 | 1 / 0 / 0 (0) | 0 / 0 | 0 | 44523.8 / 69004.2 | 763 | none |
| L1 claims | 100 seats | 18.5 | 149 | 23+13 | 2 | 4 / 0 / 0 (0) | 0 / 0 | 0 | 0.6 / 0.6 | 1990 | none |
| L1 claims | 1k seats | 66.0 | 565 | 84+58 | 32 | 13 / 0 / 0 (0) | 0 / 4 | 0 | 0.8 / 53.6 | 9826 | none |

## The race — a hot drop with seat choice

Seats held per second, from the gate to the last hold granted. Conflicts and seat-map reads are
per granted hold. Deferred confirmations kept a real 40-minute hold through the crowd.

### PostgreSQL, 1 node

| Design | 10 seats | 100 seats | 1k seats | 10k seats |
|---|---:|---:|---:|---:|
| S1 conditional | 218 | 650 | 1.0k | 257 |
| E2 cart | 75.4 | 368 | 441 | 164 ⏱ 1 timed out at 99% sold |
| K1 window | 169 | 372 | 541 | 243 |
| S1r retry | 130 | 471 | 579 | 277 |
| L1 claims | 75.4 | 427 | 622 | 258 |

<details><summary>Race detail — PostgreSQL, 1 node</summary>

| Design | Tier | Conflicts/hold | Map reads/hold | Engine retries | Gave up | Hold p50/p99 ms | Confirm p99 ms | Deferred ok | Early rej. (transient) | Confirm retries run / sold | Final buyer sold |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| S1 conditional | 10 seats | 8.84 | 17.40 | 0 | 0 | 5.39 / 60 | 68 | 14/14 | 0 (0) | 0 / 0 | 0 |
| S1 conditional | 100 seats | 2.90 | 6.39 | 0 | 0 | 5.19 / 53 | 81 | 16/16 | 0 (0) | 0 / 0 | 0 |
| S1 conditional | 1k seats | 1.18 | 3.09 | 0 | 0 | 6.95 / 88 | 99 | 48/48 | 0 (0) | 0 / 0 | 0 |
| S1 conditional | 10k seats | 0.41 | 3.65 | 0 | 0 | 31 / 304 | 320 | 415/415 | 0 (0) | 0 / 0 | 0 |
| E2 cart | 10 seats | 8.66 | 17.68 | 0 | 0 | 18 / 134 | 217 | 14/14 | 0 (0) | 0 / 0 | 0 |
| E2 cart | 100 seats | 3.25 | 5.92 | 0 | 0 | 9.39 / 100 | 91 | 15/15 | 0 (0) | 0 / 0 | 0 |
| E2 cart | 1k seats | 1.66 | 4.04 | 0 | 0 | 11 / 167 | 188 | 54/54 | 0 (0) | 0 / 0 | 0 |
| E2 cart | 10k seats | 0.41 | 3.22 | 4 | 0 | 79 / 391 | 397 | 365/365 | 0 (0) | 0 / 0 | 0 |
| K1 window | 10 seats | 9.03 | 17.89 | 0 | 0 | 6.90 / 84 | 76 | 14/14 | 0 (0) | 0 / 0 | 0 |
| K1 window | 100 seats | 3.09 | 6.51 | 0 | 0 | 9.04 / 87 | 103 | 17/17 | 0 (0) | 0 / 0 | 0 |
| K1 window | 1k seats | 1.30 | 3.74 | 0 | 0 | 15 / 102 | 105 | 42/42 | 0 (0) | 0 / 0 | 0 |
| K1 window | 10k seats | 0.43 | 3.84 | 1 | 0 | 30 / 314 | 326 | 412/412 | 0 (0) | 0 / 0 | 0 |
| S1r retry | 10 seats | 9.18 | 17.86 | 0 | 0 | 7.37 / 76 | 111 | 14/14 | 0 (0) | 0 / 0 | 0 |
| S1r retry | 100 seats | 3.01 | 5.52 | 0 | 0 | 9.87 / 87 | 80 | 13/13 | 0 (0) | 0 / 0 | 0 |
| S1r retry | 1k seats | 1.35 | 3.55 | 0 | 0 | 15 / 95 | 119 | 48/48 | 0 (0) | 0 / 0 | 0 |
| S1r retry | 10k seats | 0.42 | 3.58 | 0 | 0 | 28 / 305 | 306 | 407/407 | 0 (0) | 0 / 0 | 0 |
| L1 claims | 10 seats | 9.77 | 18.45 | 0 | 0 | 11 / 116 | 247 | 5/5 | 0 (0) | 0 / 0 | 0 |
| L1 claims | 100 seats | 2.92 | 5.41 | 0 | 0 | 9.54 / 81 | 101 | 16/16 | 0 (0) | 0 / 0 | 0 |
| L1 claims | 1k seats | 1.29 | 3.36 | 0 | 0 | 14 / 75 | 114 | 50/50 | 0 (0) | 0 / 0 | 0 |
| L1 claims | 10k seats | 0.41 | 4.21 | 0 | 0 | 23 / 300 | 386 | 433/433 | 0 (0) | 0 / 0 | 0 |

</details>

## Reads

Each question in isolation, fresh key per execution, ops/s.

### PostgreSQL, 1 node

| Question | S1 conditional | E2 cart | K1 window | S1r retry | L1 claims |
|---|---:|---:|---:|---:|---:|
| Seat map of a section — 10 seats | 9.2k | 4.3k | 8.8k | 8.3k | 17k |
| Seat map of a section — 100 seats | 6.7k | 3.2k | 7.1k | 6.9k | 15k |
| Seat map of a section — 1k seats | 7.9k | 5.0k | 11k | 9.3k | 10k |
| Seat map of a section — 10k seats | 5.4k | 3.4k | 6.3k | 5.6k | 6.0k |
| Seats available per section — 10 seats | 7.8k | 4.0k | 7.9k | 8.4k | 11k |
| Seats available per section — 100 seats | 6.6k | 2.8k | 7.5k | 7.4k | 11k |
| Seats available per section — 1k seats | 3.3k | 1.6k | 3.3k | 2.8k | 5.2k |
| Seats available per section — 10k seats | 435 | 276 | 445 | 504 | 558 |
| Seats available for an event — 10 seats | 8.9k | 3.5k | 8.3k | 8.9k | 5.2k |
| Seats available for an event — 100 seats | 7.1k | 3.0k | 7.4k | 7.5k | 5.2k |
| Seats available for an event — 1k seats | 3.3k | 1.9k | 3.6k | 4.0k | 3.3k |
| Seats available for an event — 10k seats | 533 | 332 | 536 | 491 | 932 |
| A hold's seats (basket page) | 10k | 12k | 10k | 11k | 14k |
| A customer's tickets | 19k | 11k | 15k | 18k | 14k |
| Ticket by id (control) | 19k | 23k | 18k | 17k | 17k |

## Isolated writes, publishing and storage

### PostgreSQL, 1 node

| Design | Hold+confirm /s | Release /s | Refund /s | Publish p50 ms by size | Storage | Audit |
|---|---:|---:|---:|---:|---:|---|
| S1 conditional | 296 | 697 | 4.1k | 10 seats: 2.29 · 100 seats: 2.61 · 1k seats: 11 · 10k seats: 59 · 100k seats: 804 | 14.8 MiB | consistent |
| E2 cart | 179 | 1.5k | 3.7k | 10 seats: 2.29 · 100 seats: 2.88 · 1k seats: 12 · 10k seats: 72 · 100k seats: 748 | 15.6 MiB | consistent |
| K1 window | 369 | 1.7k | 3.8k | 10 seats: 2.48 · 100 seats: 3.31 · 1k seats: 13 · 10k seats: 99 · 100k seats: 1287 | 15.1 MiB | consistent |
| S1r retry | 267 | 1.4k | 3.6k | 10 seats: 2.57 · 100 seats: 3.12 · 1k seats: 12 · 10k seats: 60 · 100k seats: 667 | 15.4 MiB | consistent |
| L1 claims | 289 | 1.3k | 4.1k | 10 seats: 2.28 · 100 seats: 2.18 · 1k seats: 2.50 · 10k seats: 1.84 · 100k seats: 1.45 | 12.8 MiB | consistent |


## Controlled pairs — one decision at a time

**This run's error bar.** Ten designs read tickets through byte-identical SQL on identical data,
so their `q05`/`q06` differences are noise: 6 comparisons, median disagreement **1.11x**, worst **1.28x**.
Read rows inside the worst disagreement are omitted below. **Race, lifecycle and write rows are always
shown**, and so is every correctness difference.

### S1 conditional → E2 cart

*Expiry copied onto every seat row, or stored once on the cart?*

| Topology | Measurement | S1 conditional | E2 cart | Change |
|---|---|---:|---:|---|
| PostgreSQL, 1 node | race, 10 seats — seats/s | 218 | 75.4 | 2.9x slower |
| PostgreSQL, 1 node | race, 100 seats — seats/s | 650 | 368 | 1.8x slower |
| PostgreSQL, 1 node | race, 1k seats — seats/s | 1.0k | 441 | 2.3x slower |
| PostgreSQL, 1 node | race, 10k seats — seats/s | 257 | 164 ⏱ 1 timed out at 99% sold | 1.6x slower |
| PostgreSQL, 1 node | lifecycle, 10 seats — refused late/boundary/early (transient); violations | 0/0/0 (0); 0 | 0/0/0 (0); 0 |  |
| PostgreSQL, 1 node | lifecycle, 100 seats — refused late/boundary/early (transient); violations | 1/0/0 (0); 0 | 0/0/0 (0); 0 |  |
| PostgreSQL, 1 node | lifecycle, 1k seats — refused late/boundary/early (transient); violations | 8/0/0 (0); 0 | 5/0/0 (0); 0 |  |
| PostgreSQL, 1 node | hold /s | 296 | 179 | 1.6x slower |
| PostgreSQL, 1 node | release /s | 697 | 1.5k | 2.1x faster |
| PostgreSQL, 1 node | cancel /s | 4.1k | 3.7k | 1.1x slower |
| PostgreSQL, 1 node | publish 10 seats p50 ms | 2.29 | 2.29 | 1.0x faster |
| PostgreSQL, 1 node | publish 100 seats p50 ms | 2.61 | 2.88 | 1.1x slower |
| PostgreSQL, 1 node | publish 1k seats p50 ms | 11 | 12 | 1.1x slower |
| PostgreSQL, 1 node | publish 10k seats p50 ms | 59 | 72 | 1.2x slower |
| PostgreSQL, 1 node | publish 100k seats p50 ms | 804 | 748 | 1.1x faster |
| PostgreSQL, 1 node | Seat map of a section — 10 seats | 9.2k | 4.3k | 2.1x slower |
| PostgreSQL, 1 node | Seats available per section — 10 seats | 7.8k | 4.0k | 1.9x slower |
| PostgreSQL, 1 node | Seats available for an event — 10 seats | 8.9k | 3.5k | 2.5x slower |
| PostgreSQL, 1 node | Seat map of a section — 100 seats | 6.7k | 3.2k | 2.1x slower |
| PostgreSQL, 1 node | Seats available per section — 100 seats | 6.6k | 2.8k | 2.3x slower |
| PostgreSQL, 1 node | Seats available for an event — 100 seats | 7.1k | 3.0k | 2.4x slower |
| PostgreSQL, 1 node | Seat map of a section — 1k seats | 7.9k | 5.0k | 1.6x slower |
| PostgreSQL, 1 node | Seats available per section — 1k seats | 3.3k | 1.6k | 2.0x slower |
| PostgreSQL, 1 node | Seats available for an event — 1k seats | 3.3k | 1.9k | 1.8x slower |
| PostgreSQL, 1 node | Seat map of a section — 10k seats | 5.4k | 3.4k | 1.6x slower |
| PostgreSQL, 1 node | Seats available per section — 10k seats | 435 | 276 | 1.6x slower |
| PostgreSQL, 1 node | Seats available for an event — 10k seats | 533 | 332 | 1.6x slower |
| PostgreSQL, 1 node | A customer's tickets | 19k | 11k | 1.6x slower |

### S1 conditional → K1 window

*What does a guaranteed payment window cost, and what does it remove?*

| Topology | Measurement | S1 conditional | K1 window | Change |
|---|---|---:|---:|---|
| PostgreSQL, 1 node | race, 10 seats — seats/s | 218 | 169 | 1.3x slower |
| PostgreSQL, 1 node | race, 100 seats — seats/s | 650 | 372 | 1.7x slower |
| PostgreSQL, 1 node | race, 1k seats — seats/s | 1.0k | 541 | 1.9x slower |
| PostgreSQL, 1 node | race, 10k seats — seats/s | 257 | 243 | 1.1x slower |
| PostgreSQL, 1 node | lifecycle, 10 seats — refused late/boundary/early (transient); violations | 0/0/0 (0); 0 | 0/0/0 (0); 0 |  |
| PostgreSQL, 1 node | lifecycle, 100 seats — refused late/boundary/early (transient); violations | 1/0/0 (0); 0 | 0/0/0 (0); 0 |  |
| PostgreSQL, 1 node | lifecycle, 1k seats — refused late/boundary/early (transient); violations | 8/0/0 (0); 0 | 0/0/0 (0); 0 |  |
| PostgreSQL, 1 node | hold /s | 296 | 369 | 1.2x faster |
| PostgreSQL, 1 node | release /s | 697 | 1.7k | 2.4x faster |
| PostgreSQL, 1 node | cancel /s | 4.1k | 3.8k | 1.1x slower |
| PostgreSQL, 1 node | publish 10 seats p50 ms | 2.29 | 2.48 | 1.1x slower |
| PostgreSQL, 1 node | publish 100 seats p50 ms | 2.61 | 3.31 | 1.3x slower |
| PostgreSQL, 1 node | publish 1k seats p50 ms | 11 | 13 | 1.2x slower |
| PostgreSQL, 1 node | publish 10k seats p50 ms | 59 | 99 | 1.7x slower |
| PostgreSQL, 1 node | publish 100k seats p50 ms | 804 | 1287 | 1.6x slower |
| PostgreSQL, 1 node | Seat map of a section — 1k seats | 7.9k | 11k | 1.4x faster |
| PostgreSQL, 1 node | A customer's tickets | 19k | 15k | 1.3x slower |

### S1 conditional → L1 claims

*Pre-created per-event seat rows, or claims created on hold?*

| Topology | Measurement | S1 conditional | L1 claims | Change |
|---|---|---:|---:|---|
| PostgreSQL, 1 node | race, 10 seats — seats/s | 218 | 75.4 | 2.9x slower |
| PostgreSQL, 1 node | race, 100 seats — seats/s | 650 | 427 | 1.5x slower |
| PostgreSQL, 1 node | race, 1k seats — seats/s | 1.0k | 622 | 1.6x slower |
| PostgreSQL, 1 node | race, 10k seats — seats/s | 257 | 258 | 1.0x faster |
| PostgreSQL, 1 node | lifecycle, 10 seats — refused late/boundary/early (transient); violations | 0/0/0 (0); 0 | 1/0/0 (0); 0 |  |
| PostgreSQL, 1 node | lifecycle, 100 seats — refused late/boundary/early (transient); violations | 1/0/0 (0); 0 | 4/0/0 (0); 0 |  |
| PostgreSQL, 1 node | lifecycle, 1k seats — refused late/boundary/early (transient); violations | 8/0/0 (0); 0 | 13/0/0 (0); 0 |  |
| PostgreSQL, 1 node | hold /s | 296 | 289 | 1.0x slower |
| PostgreSQL, 1 node | release /s | 697 | 1.3k | 1.9x faster |
| PostgreSQL, 1 node | cancel /s | 4.1k | 4.1k | 1.0x slower |
| PostgreSQL, 1 node | publish 10 seats p50 ms | 2.29 | 2.28 | 1.0x faster |
| PostgreSQL, 1 node | publish 100 seats p50 ms | 2.61 | 2.18 | 1.2x faster |
| PostgreSQL, 1 node | publish 1k seats p50 ms | 11 | 2.50 | 4.4x faster |
| PostgreSQL, 1 node | publish 10k seats p50 ms | 59 | 1.84 | 31.9x faster |
| PostgreSQL, 1 node | publish 100k seats p50 ms | 804 | 1.45 | 555.2x faster |
| PostgreSQL, 1 node | Seat map of a section — 10 seats | 9.2k | 17k | 1.9x faster |
| PostgreSQL, 1 node | Seats available per section — 10 seats | 7.8k | 11k | 1.4x faster |
| PostgreSQL, 1 node | Seats available for an event — 10 seats | 8.9k | 5.2k | 1.7x slower |
| PostgreSQL, 1 node | Seat map of a section — 100 seats | 6.7k | 15k | 2.2x faster |
| PostgreSQL, 1 node | Seats available per section — 100 seats | 6.6k | 11k | 1.7x faster |
| PostgreSQL, 1 node | Seats available for an event — 100 seats | 7.1k | 5.2k | 1.4x slower |
| PostgreSQL, 1 node | Seat map of a section — 1k seats | 7.9k | 10k | 1.3x faster |
| PostgreSQL, 1 node | Seats available per section — 1k seats | 3.3k | 5.2k | 1.6x faster |
| PostgreSQL, 1 node | Seats available for an event — 10k seats | 533 | 932 | 1.7x faster |
| PostgreSQL, 1 node | A hold's seats (basket page) | 10k | 14k | 1.4x faster |
| PostgreSQL, 1 node | A customer's tickets | 19k | 14k | 1.3x slower |

### S1 conditional → S1r retry

*Does retrying a short confirmation once remove transient refusals, and what does it cost?*

| Topology | Measurement | S1 conditional | S1r retry | Change |
|---|---|---:|---:|---|
| PostgreSQL, 1 node | race, 10 seats — seats/s | 218 | 130 | 1.7x slower |
| PostgreSQL, 1 node | race, 100 seats — seats/s | 650 | 471 | 1.4x slower |
| PostgreSQL, 1 node | race, 1k seats — seats/s | 1.0k | 579 | 1.8x slower |
| PostgreSQL, 1 node | race, 10k seats — seats/s | 257 | 277 | 1.1x faster |
| PostgreSQL, 1 node | lifecycle, 10 seats — refused late/boundary/early (transient); violations | 0/0/0 (0); 0 | 0/0/0 (0); 0 |  |
| PostgreSQL, 1 node | lifecycle, 100 seats — refused late/boundary/early (transient); violations | 1/0/0 (0); 0 | 2/0/0 (0); 0 |  |
| PostgreSQL, 1 node | lifecycle, 1k seats — refused late/boundary/early (transient); violations | 8/0/0 (0); 0 | 5/0/0 (0); 0 |  |
| PostgreSQL, 1 node | hold /s | 296 | 267 | 1.1x slower |
| PostgreSQL, 1 node | release /s | 697 | 1.4k | 2.0x faster |
| PostgreSQL, 1 node | cancel /s | 4.1k | 3.6k | 1.1x slower |
| PostgreSQL, 1 node | publish 10 seats p50 ms | 2.29 | 2.57 | 1.1x slower |
| PostgreSQL, 1 node | publish 100 seats p50 ms | 2.61 | 3.12 | 1.2x slower |
| PostgreSQL, 1 node | publish 1k seats p50 ms | 11 | 12 | 1.1x slower |
| PostgreSQL, 1 node | publish 10k seats p50 ms | 59 | 60 | 1.0x slower |
| PostgreSQL, 1 node | publish 100k seats p50 ms | 804 | 667 | 1.2x faster |


## Measurement hazard: CPU-quota throttling

Cells show *throttled periods / periods (throttled time)*: for the **database** container(s) over
the whole cell (summed across nodes), and for the **client** per phase.

### PostgreSQL, 1 node

| Design | database, whole cell | client lifecycle#1 | client race#1 | client read | client write:cancel | client write:hold | client write:publish | client write:release |
|---|---:|---:|---:|---:|---:|---:|---:|---:|
| S1 conditional | 1529/2068 (447s) | 0/440 (0.0s) | 0/427 (0.0s) | 92/1054 (5.4s) | 0/7 (0.0s) | 0/41 (0.0s) | 0/20 (0.0s) | 0/1 (0.0s) |
| E2 cart | 1818/2390 (693s) | 0/463 (0.0s) | 1/680 (0.0s) | 89/1053 (4.6s) | 0/9 (0.0s) | 0/72 (0.0s) | 0/22 (0.0s) | — |
| K1 window | 1561/2167 (441s) | 0/488 (0.0s) | 0/467 (0.0s) | 74/1053 (3.5s) | 0/9 (0.0s) | 0/39 (0.0s) | 0/23 (0.0s) | 0/1 (0.0s) |
| S1r retry | 1511/2014 (414s) | 0/388 (0.0s) | 2/414 (0.5s) | 79/1053 (4.6s) | 0/10 (0.0s) | 0/50 (0.0s) | 0/22 (0.0s) | — |
| L1 claims | 1522/2053 (405s) | 0/403 (0.0s) | 0/453 (0.0s) | 121/1052 (7.6s) | 0/10 (0.0s) | 1/48 (0.0s) | 0/16 (0.0s) | 0/1 (0.0s) |


---

Plans for every read and write statement are in `results/dc12-pg-am01/<topology>/plans/<design>.txt`
(`EXPLAIN (ANALYZE, BUFFERS)` on PostgreSQL, `EXPLAIN (ANALYZE, DIST)` on YugabyteDB; write
statements inside rolled-back transactions).

## Provenance: code, data and analyses

This file is **generated from the measurements** and contains no interpretation.
Conclusions live in separate signed analyses, so several analysts — including
different AI models — can read the same numbers and each record what they make of
them. Where two analyses disagree, the disagreement is a finding and stays visible.

| | |
|---|---|
| Run id | `devchecks/dc12-pg-am01` |
| Result files | 5 |
| **Inputs digest** | `d1789a33950a0987` |
| Repository version | `study-03/v0.1-handoff-amendment-01-2-ge10dc6c` (`e10dc6c8a6af`), clean |

Quote the digest in any analysis: it identifies the data. The repository version
identifies the code — check it out (`git checkout <commit>`) to reproduce the run, or
`git diff <commit>` to see what has changed in the study since.

### Analyses of this run

**None yet.** This run has measurements but no interpretation.

### Adding an analysis

Copy `reports/analyses/TEMPLATE.md`, fill in the frontmatter, and regenerate this report.
If an analysis already exists for digest `d1789a33950a0987`, read it first and write a new one only to
**disagree, extend, or bring a different perspective**. Never edit another analyst's file.
