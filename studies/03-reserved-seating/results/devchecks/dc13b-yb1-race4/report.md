# Study 03 — results: reserved seating (venue → event → seat)

Generated from `results/dc13b-yb1-race4`. Every number comes from exactly one JSON file in that
directory. **This file contains measurements and no interpretation**; conclusions live in
signed analyses, indexed at the end.

## TL;DR — measured facts

*Selected from the tables below by fixed rules; no interpretation. Read a signed analysis for
what these facts mean.*

- **Cells:** 5 across 1 topologies; 0 failed.
- **Negative controls:** 0 of 0 fired as designed to.
- ❌ **Invariant violations in designs meant to be correct:** S1 conditional on YugabyteDB, 1 node (RF=1) (2 early rejections (2 transient)); E1 sweeper on YugabyteDB, 1 node (RF=1) (5 early rejections (5 transient)); K1 window on YugabyteDB, 1 node (RF=1) (10 early rejections (10 transient)); L1 claims on YugabyteDB, 1 node (RF=1) (4 early rejections (4 transient)).
- **S1r confirmation retries, YugabyteDB, 1 node (RF=1)** (race, all tiers and trials): 6 retried, 6 of them sold.
- **Early rejections in the race, correct designs, YugabyteDB, 1 node (RF=1):** S1 conditional 2 (2 transient), E1 sweeper 5 (5 transient), K1 window 10 (10 transient), L1 claims 4 (4 transient).
- **Race, YugabyteDB, 1 node (RF=1)** (seats held/s; designs meant to be correct, without violations):
  - 10 seats: highest S1 conditional 14.2/s, lowest L1 claims 6.7/s
  - 100 seats: highest S1r retry 38.1/s, lowest L1 claims 24.9/s
  - 1k seats: highest L1 claims 38.1/s, lowest S1r retry 37.8/s
  - 10k seats: highest S1r retry 27.0/s, lowest S1r retry 27.0/s (1 of 1 timed out)

## What was measured

|  |  |
|---|---|
| Run id | `devchecks/dc13b-yb1-race4` |
| Environment | [`host-zenbook-ux5406sa`](../../../docs/environments/host-zenbook-ux5406sa.md) |
| Scale | `tiny` — 5 venues, 84 events, 5741 tickets sold, 131 live and 34 expired holds at load |
| Events (kind_seats: count) | catalogue_10: 40 · catalogue_100: 8 · catalogue_1000: 2 · catalogue_10000: 1 · lifecycle_10: 4 · lifecycle_100: 2 · lifecycle_1000: 1 · race_10: 20 · race_100: 4 · race_1000: 1 · race_10000: 1 |
| Seed | 42 — identical data in every cell |
| Reads / isolated writes | 8 workers, 5s measured + 2s warmup, 1 trial(s) |
| Race | 32 buyers per event, real hold TTL 40m0s, 10% of holds deferred until the crowd has gone, timeout 1m0s, tier budget 3m0s, give up after 20 conflicts |

**Engines**

| Topology | Version | READ COMMITTED runs as | Client connections to | DB clock − client clock (ms) |
|---|---|---|---|---:|
| YugabyteDB, 1 node (RF=1) | `PostgreSQL 15.12-YB-2025.2.6.0-b0 on x86_64-pc-linux-gnu, compiled by clang version 19.1.0 (https://github.com/yugabyte/llvm-project.git a2a6b655e14e7fa1fcf1011a6cb29cb8575249c0), 64-bit` | read committed | 1 node(s) | -0.01 |

> Limits of this environment that bound every number below: the client and the databases
> share one laptop's eight cores; the "3-node" cluster has no real network between nodes; the
> race and the lifecycle are closed loops, so tails are a floor; the lifecycle runs on
> compressed time; E0's clock skew is injected.

## Correctness gate

Six read questions checked against Go-computed truth on a load with sold seats, live holds and
expired holds, plus an audit of the loaded state. A failure aborts the cell.

| Design | YugabyteDB, 1 node (RF=1) |
|---|---|
| S1 conditional | ✅ 51/51 |
| E1 sweeper | ✅ 51/51 (after its sweeper released 65 expired seats) |
| K1 window | ✅ 51/51 |
| S1r retry | ✅ 51/51 |
| L1 claims | ✅ 51/51 |

## Negative controls

Deliberately wrong designs (and E1 with its sweeper stopped). Each must be seen to break the
invariant it exists to test; otherwise the absence of violations elsewhere in that experiment
is not evidence of correctness.

No control was measured in this run.

## The race — a hot drop with seat choice

Seats held per second, from the gate to the last hold granted. Conflicts and seat-map reads are
per granted hold. Deferred confirmations kept a real 40-minute hold through the crowd.

### YugabyteDB, 1 node (RF=1)

| Design | 10 seats | 100 seats | 1k seats | 10k seats |
|---|---:|---:|---:|---:|
| S1 conditional | 14.2 (4 trials, spread 30%) | 37.7 (4 trials, spread 36%) | ~~42.2 (4 trials, spread 57%)~~ ❌ 1 early rejections (1 transient) | ~~29.2 (4 trials, spread 15%)~~ ❌ 1 early rejections (1 transient); ⏱ 4 timed out at 18% sold |
| E1 sweeper | 14.1 (4 trials, spread 28%) | 34.4 (4 trials, spread 9%) | ~~36.1 (4 trials, spread 5%)~~ ❌ 2 early rejections (2 transient) | ~~28.1 (4 trials, spread 26%)~~ ❌ 3 early rejections (3 transient); ⏱ 4 timed out at 18% sold |
| K1 window | ~~14.0 (4 trials, spread 17%)~~ ❌ 1 early rejections (1 transient) | ~~35.2 (4 trials, spread 27%)~~ ❌ 1 early rejections (1 transient) | ~~41.0 (4 trials, spread 21%)~~ ❌ 3 early rejections (3 transient) | ~~26.9 (4 trials, spread 18%)~~ ❌ 5 early rejections (5 transient); ⏱ 4 timed out at 16% sold |
| S1r retry | 13.9 (4 trials, spread 40%) | 38.1 (4 trials, spread 19%) | 37.8 (4 trials, spread 22%) | 27.0 (4 trials, spread 6%) ⏱ 4 timed out at 16% sold |
| L1 claims | 6.7 (4 trials, spread 55%) | 24.9 (4 trials, spread 30%) | 38.1 (4 trials, spread 19%) | ~~48.7 (4 trials, spread 12%)~~ ❌ 4 early rejections (4 transient); ⏱ 4 timed out at 30% sold |

<details><summary>Race detail — YugabyteDB, 1 node (RF=1)</summary>

| Design | Tier | Conflicts/hold | Map reads/hold | Engine retries | Gave up | Hold p50/p99 ms | Confirm p99 ms | Deferred ok | Early rej. (transient) | Confirm retries run / sold | Final buyer sold |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| S1 conditional | 10 seats | 6.78 | 14.55 | 0 | 0 | 104 / 678 | 997 | 50/50 | 0 (0) | 0 / 0 | 0 |
| S1 conditional | 100 seats | 2.46 | 4.48 | 0 | 0 | 104 / 603 | 999 | 69/69 | 0 (0) | 0 / 0 | 0 |
| S1 conditional | 1k seats | 0.83 | 3.30 | 0 | 0 | 107 / 400 | 511 | 204/204 | 1 (1) | 0 / 0 | 0 |
| S1 conditional | 10k seats | 0.31 | 1.34 | 0 | 0 | 300 / 792 | 1703 | 293/293 | 1 (1) | 0 / 0 | 0 |
| E1 sweeper | 10 seats | 7.09 | 15.28 | 0 | 0 | 101 / 890 | 692 | 50/50 | 0 (0) | 0 / 0 | 0 |
| E1 sweeper | 100 seats | 2.63 | 5.01 | 0 | 0 | 113 / 603 | 811 | 80/80 | 0 (0) | 0 / 0 | 0 |
| E1 sweeper | 1k seats | 1.09 | 3.11 | 0 | 0 | 189 / 576 | 519 | 203/203 | 2 (2) | 0 / 0 | 0 |
| E1 sweeper | 10k seats | 0.32 | 1.35 | 0 | 0 | 294 / 791 | 1004 | 274/274 | 3 (3) | 0 / 0 | 0 |
| K1 window | 10 seats | 7.23 | 14.96 | 0 | 0 | 94 / 685 | 1171 | 57/57 | 1 (1) | 0 / 0 | 0 |
| K1 window | 100 seats | 2.69 | 5.05 | 0 | 0 | 109 / 585 | 477 | 77/77 | 1 (1) | 0 / 0 | 0 |
| K1 window | 1k seats | 1.05 | 2.95 | 0 | 0 | 125 / 402 | 498 | 209/209 | 3 (3) | 0 / 0 | 0 |
| K1 window | 10k seats | 0.38 | 1.41 | 0 | 0 | 291 / 885 | 1212 | 291/291 | 5 (5) | 0 / 0 | 0 |
| S1r retry | 10 seats | 6.04 | 13.44 | 0 | 0 | 104 / 583 | 912 | 48/48 | 0 (0) | 0 / 0 | 0 |
| S1r retry | 100 seats | 2.62 | 5.11 | 0 | 0 | 107 / 787 | 378 | 68/68 | 0 (0) | 1 / 1 | 0 |
| S1r retry | 1k seats | 0.96 | 2.57 | 0 | 0 | 192 / 490 | 589 | 186/186 | 0 (0) | 2 / 2 | 0 |
| S1r retry | 10k seats | 0.38 | 1.42 | 0 | 0 | 300 / 794 | 1890 | 292/292 | 0 (0) | 3 / 3 | 0 |
| L1 claims | 10 seats | 7.52 | 15.35 | 1 | 0 | 183 / 1824 | 1099 | 48/48 | 0 (0) | 0 / 0 | 0 |
| L1 claims | 100 seats | 2.28 | 4.32 | 1 | 0 | 114 / 1716 | 788 | 77/77 | 0 (0) | 0 / 0 | 0 |
| L1 claims | 1k seats | 0.79 | 2.91 | 0 | 0 | 193 / 408 | 500 | 211/211 | 0 (0) | 0 / 0 | 0 |
| L1 claims | 10k seats | 0.27 | 1.47 | 0 | 0 | 204 / 682 | 1525 | 499/499 | 4 (4) | 0 / 0 | 0 |

</details>

## Reads

Each question in isolation, fresh key per execution, ops/s.

### YugabyteDB, 1 node (RF=1)

| Question | S1 conditional | E1 sweeper | K1 window | S1r retry | L1 claims |
|---|---:|---:|---:|---:|---:|

## Isolated writes, publishing and storage

### YugabyteDB, 1 node (RF=1)

| Design | Hold+confirm /s | Release /s | Refund /s | Publish p50 ms by size | Storage | Audit |
|---|---:|---:|---:|---:|---:|---|
| S1 conditional | — | — | — |  | n/a on YugabyteDB | consistent |
| E1 sweeper | — | — | — |  | n/a on YugabyteDB | consistent |
| K1 window | — | — | — |  | n/a on YugabyteDB | consistent |
| S1r retry | — | — | — |  | n/a on YugabyteDB | consistent |
| L1 claims | — | — | — |  | n/a on YugabyteDB | consistent |


## Controlled pairs — one decision at a time

### S1 conditional → E1 sweeper

*Expiry judged at use time, or made to take effect by a sweeper?*

| Topology | Measurement | S1 conditional | E1 sweeper | Change |
|---|---|---:|---:|---|
| YugabyteDB, 1 node (RF=1) | race, 10 seats — seats/s | 14.2 (4 trials, spread 30%) | 14.1 (4 trials, spread 28%) | 1.0x slower |
| YugabyteDB, 1 node (RF=1) | race, 100 seats — seats/s | 37.7 (4 trials, spread 36%) | 34.4 (4 trials, spread 9%) | 1.1x slower |
| YugabyteDB, 1 node (RF=1) | race, 1k seats — seats/s | ~~42.2 (4 trials, spread 57%)~~ ❌ 1 early rejections (1 transient) | ~~36.1 (4 trials, spread 5%)~~ ❌ 2 early rejections (2 transient) | 1.2x slower — ❌ a design with violations has no valid speed |
| YugabyteDB, 1 node (RF=1) | race, 10k seats — seats/s | ~~29.2 (4 trials, spread 15%)~~ ❌ 1 early rejections (1 transient); ⏱ 4 timed out at 18% sold | ~~28.1 (4 trials, spread 26%)~~ ❌ 3 early rejections (3 transient); ⏱ 4 timed out at 18% sold | 1.0x slower — ❌ a design with violations has no valid speed |

### S1 conditional → K1 window

*What does a guaranteed payment window cost, and what does it remove?*

| Topology | Measurement | S1 conditional | K1 window | Change |
|---|---|---:|---:|---|
| YugabyteDB, 1 node (RF=1) | race, 10 seats — seats/s | 14.2 (4 trials, spread 30%) | ~~14.0 (4 trials, spread 17%)~~ ❌ 1 early rejections (1 transient) | 1.0x slower — ❌ a design with violations has no valid speed |
| YugabyteDB, 1 node (RF=1) | race, 100 seats — seats/s | 37.7 (4 trials, spread 36%) | ~~35.2 (4 trials, spread 27%)~~ ❌ 1 early rejections (1 transient) | 1.1x slower — ❌ a design with violations has no valid speed |
| YugabyteDB, 1 node (RF=1) | race, 1k seats — seats/s | ~~42.2 (4 trials, spread 57%)~~ ❌ 1 early rejections (1 transient) | ~~41.0 (4 trials, spread 21%)~~ ❌ 3 early rejections (3 transient) | 1.0x slower — ❌ a design with violations has no valid speed |
| YugabyteDB, 1 node (RF=1) | race, 10k seats — seats/s | ~~29.2 (4 trials, spread 15%)~~ ❌ 1 early rejections (1 transient); ⏱ 4 timed out at 18% sold | ~~26.9 (4 trials, spread 18%)~~ ❌ 5 early rejections (5 transient); ⏱ 4 timed out at 16% sold | 1.1x slower — ❌ a design with violations has no valid speed |

### S1 conditional → L1 claims

*Pre-created per-event seat rows, or claims created on hold?*

| Topology | Measurement | S1 conditional | L1 claims | Change |
|---|---|---:|---:|---|
| YugabyteDB, 1 node (RF=1) | race, 10 seats — seats/s | 14.2 (4 trials, spread 30%) | 6.7 (4 trials, spread 55%) | 2.1x slower |
| YugabyteDB, 1 node (RF=1) | race, 100 seats — seats/s | 37.7 (4 trials, spread 36%) | 24.9 (4 trials, spread 30%) | 1.5x slower |
| YugabyteDB, 1 node (RF=1) | race, 1k seats — seats/s | ~~42.2 (4 trials, spread 57%)~~ ❌ 1 early rejections (1 transient) | 38.1 (4 trials, spread 19%) | 1.1x slower — ❌ a design with violations has no valid speed |
| YugabyteDB, 1 node (RF=1) | race, 10k seats — seats/s | ~~29.2 (4 trials, spread 15%)~~ ❌ 1 early rejections (1 transient); ⏱ 4 timed out at 18% sold | ~~48.7 (4 trials, spread 12%)~~ ❌ 4 early rejections (4 transient); ⏱ 4 timed out at 30% sold | 1.7x faster — ❌ a design with violations has no valid speed |

### S1 conditional → S1r retry

*Does retrying a short confirmation once remove transient refusals, and what does it cost?*

| Topology | Measurement | S1 conditional | S1r retry | Change |
|---|---|---:|---:|---|
| YugabyteDB, 1 node (RF=1) | race, 10 seats — seats/s | 14.2 (4 trials, spread 30%) | 13.9 (4 trials, spread 40%) | 1.0x slower |
| YugabyteDB, 1 node (RF=1) | race, 100 seats — seats/s | 37.7 (4 trials, spread 36%) | 38.1 (4 trials, spread 19%) | 1.0x faster |
| YugabyteDB, 1 node (RF=1) | race, 1k seats — seats/s | ~~42.2 (4 trials, spread 57%)~~ ❌ 1 early rejections (1 transient) | 37.8 (4 trials, spread 22%) | 1.1x slower — ❌ a design with violations has no valid speed |
| YugabyteDB, 1 node (RF=1) | race, 10k seats — seats/s | ~~29.2 (4 trials, spread 15%)~~ ❌ 1 early rejections (1 transient); ⏱ 4 timed out at 18% sold | 27.0 (4 trials, spread 6%) ⏱ 4 timed out at 16% sold | 1.1x slower — ❌ a design with violations has no valid speed |


## Measurement hazard: CPU-quota throttling

Cells show *throttled periods / periods (throttled time)*: for the **database** container(s) over
the whole cell (summed across nodes), and for the **client** per phase.

### YugabyteDB, 1 node (RF=1)

| Design | database, whole cell | client race#1 | client race#2 | client race#3 | client race#4 |
|---|---:|---:|---:|---:|---:|
| S1 conditional | 5117/5748 (2558s) | 0/1301 (0.0s) | 0/1096 (0.0s) | 0/1232 (0.0s) | 0/1163 (0.0s) |
| E1 sweeper | 5215/5870 (2587s) | 0/1279 (0.0s) | 0/1246 (0.0s) | 0/1221 (0.0s) | 0/1212 (0.0s) |
| K1 window | 5129/5846 (2552s) | 0/1236 (0.0s) | 0/1251 (0.0s) | 0/1155 (0.0s) | 0/1168 (0.0s) |
| S1r retry | 5192/5868 (2557s) | 0/1261 (0.0s) | 0/1175 (0.0s) | 0/1203 (0.0s) | 0/1237 (0.0s) |
| L1 claims | 6051/6875 (2694s) | 0/1474 (0.0s) | 0/1409 (0.0s) | 0/1456 (0.0s) | 0/1420 (0.0s) |


---

Plans for every read and write statement are in `results/dc13b-yb1-race4/<topology>/plans/<design>.txt`
(`EXPLAIN (ANALYZE, BUFFERS)` on PostgreSQL, `EXPLAIN (ANALYZE, DIST)` on YugabyteDB; write
statements inside rolled-back transactions).

## Provenance: code, data and analyses

This file is **generated from the measurements** and contains no interpretation.
Conclusions live in separate signed analyses, so several analysts — including
different AI models — can read the same numbers and each record what they make of
them. Where two analyses disagree, the disagreement is a finding and stays visible.

| | |
|---|---|
| Run id | `devchecks/dc13b-yb1-race4` |
| Result files | 5 |
| **Inputs digest** | `689ff733a0feb2e6` |
| Repository version | `study-03/v0.1-handoff-amendment-01-6-g425966e-dirty` (`425966e9452b`), ⚠️ **dirty** — uncommitted changes; not reproducible from the commit alone |

Quote the digest in any analysis: it identifies the data. The repository version
identifies the code — check it out (`git checkout <commit>`) to reproduce the run, or
`git diff <commit>` to see what has changed in the study since.

### Analyses of this run

**None yet.** This run has measurements but no interpretation.

### Adding an analysis

Copy `reports/analyses/TEMPLATE.md`, fill in the frontmatter, and regenerate this report.
If an analysis already exists for digest `689ff733a0feb2e6`, read it first and write a new one only to
**disagree, extend, or bring a different perspective**. Never edit another analyst's file.
