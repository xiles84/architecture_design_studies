# Study 03 — results: reserved seating (venue → event → seat)

Generated from `results/dc4-pg-e1-l2`. Every number comes from exactly one JSON file in that
directory. **This file contains measurements and no interpretation**; conclusions live in
signed analyses, indexed at the end.

## TL;DR — measured facts

*Selected from the tables below by fixed rules; no interpretation. Read a signed analysis for
what these facts mean.*

- **Cells:** 2 across 1 topologies; 0 failed.
- **Negative controls:** 1 of 1 fired as designed to.
- **Invariant violations in designs meant to be correct:** none observed.
- **Sweeper outage, PostgreSQL, 1 node:** E1 13 of 13 probed seats unavailable after expiry; 0 in every other design.
- **Race, PostgreSQL, 1 node** (seats held/s; designs meant to be correct, without violations):
  - 10 seats: highest E1 sweeper 322/s, lowest L2 document 206/s
  - 100 seats: highest E1 sweeper 1.2k/s, lowest L2 document 320/s
  - 1k seats: highest E1 sweeper 1.2k/s, lowest L2 document 458/s
  - 10k seats: highest E1 sweeper 663/s, lowest L2 document 379/s

## What was measured

|  |  |
|---|---|
| Run id | `devchecks/dc4-pg-e1-l2` |
| Environment | [`host-zenbook-ux5406sa`](../../../docs/environments/host-zenbook-ux5406sa.md) |
| Scale | `tiny` — 5 venues, 84 events, 5741 tickets sold, 131 live and 34 expired holds at load |
| Events (kind_seats: count) | catalogue_10: 40 · catalogue_100: 8 · catalogue_1000: 2 · catalogue_10000: 1 · lifecycle_10: 4 · lifecycle_100: 2 · lifecycle_1000: 1 · race_10: 20 · race_100: 4 · race_1000: 1 · race_10000: 1 |
| Seed | 42 — identical data in every cell |
| Reads / isolated writes | 8 workers, 2s measured + 1s warmup, 1 trial(s) |
| Race | 32 buyers per event, real hold TTL 40m0s, 10% of holds deferred until the crowd has gone, timeout 30s, tier budget 3m0s, give up after 20 conflicts |
| Lifecycle — PostgreSQL, 1 node | one human minute = 50ms; TTL 40m0s, payment window (K1) 10m0s, guard G 2m0s, sweeper every 1m0s, outage 45m0s→1h25m0s in the first event of each tier, 25% abandoned (40% of those explicitly), E0 skew 10m0s; tiers 10,100,1000, 2 events per tier, 32 buyers |

**Engines**

| Topology | Version | READ COMMITTED runs as | Client connections to | DB clock − client clock (ms) |
|---|---|---|---|---:|
| PostgreSQL, 1 node | `PostgreSQL 17.11 (Debian 17.11-1.pgdg13+2) on x86_64-pc-linux-gnu, compiled by gcc (Debian 14.2.0-19) 14.2.0, 64-bit` | read committed | 1 node(s) | -0.00 |

> Limits of this environment that bound every number below: the client and the databases
> share one laptop's eight cores; the "3-node" cluster has no real network between nodes; the
> race and the lifecycle are closed loops, so tails are a floor; the lifecycle runs on
> compressed time; E0's clock skew is injected.

## Correctness gate

Six read questions checked against Go-computed truth on a load with sold seats, live holds and
expired holds, plus an audit of the loaded state. A failure aborts the cell.

| Design | PostgreSQL, 1 node |
|---|---|
| E1 sweeper | ✅ 51/51 (after its sweeper released 65 expired seats) |
| L2 document | ✅ 51/51 |

## Negative controls

Deliberately wrong designs (and E1 with its sweeper stopped). Each must be seen to break the
invariant it exists to test; otherwise the absence of violations elsewhere in that experiment
is not evidence of correctness.

| Control | Experiment | Topology | Fired? | What was found |
|---|---|---|---|---|
| E1 sweeper stopped | lifecycle outage | PostgreSQL, 1 node | ✅ fired | 13 of 13 probed seats unavailable after expiry |

## The hold guarantee — lifecycle

Compressed time. *Rejected* confirmations are classified by how long before the hold's expiry
the confirming transaction started: **late** (after expiry), **boundary** (less than G before),
**early** (G or more before — a violation). Violations are counted by the ledger and the audit.
Release lag is in human minutes.

### PostgreSQL, 1 node

| Design | Tier | Confirmed seats/s | Holds | Abandoned | Expired at check | Refused late / boundary / early | Outage: unavailable / probed | Leaked | Release lag p50/p99 (min) | Idle held seat-min | Violations |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|---|
| E1 sweeper | 10 seats | 2.0 | 22 | 4+3 | 1 | 0 / 0 / 0 | 0 / 0 | 0 | 1.0 / 1.0 | 377 | none |
| E1 sweeper | 100 seats | 14.1 | 145 | 22+13 | 2 | 1 / 0 / 0 | 4 / 4 | 0 | 0.8 / 42.5 | 1759 | none |
| E1 sweeper | 1k seats | 71.2 | 595 | 89+45 | 31 | 10 / 0 / 0 | 9 / 9 | 0 | 0.7 / 59.6 | 9431 | none |
| L2 document | 10 seats | 1.5 | 28 | 9+1 | 0 | 0 / 0 / 0 | 0 / 0 | 0 | 44426.8 / 68989.7 | 522 | none |
| L2 document | 100 seats | 16.5 | 152 | 25+17 | 4 | 2 / 0 / 0 | 0 / 0 | 0 | 0.3 / 0.3 | 2418 | none |
| L2 document | 1k seats | 70.7 | 585 | 71+61 | 17 | 5 / 0 / 0 | 0 / 4 | 0 | 0.5 / 1.0 | 8512 | none |

## The race — a hot drop with seat choice

Seats held per second, from the gate to the last hold granted. Conflicts and seat-map reads are
per granted hold. Deferred confirmations kept a real 40-minute hold through the crowd.

### PostgreSQL, 1 node

| Design | 10 seats | 100 seats | 1k seats | 10k seats |
|---|---:|---:|---:|---:|
| E1 sweeper | 322 | 1.2k | 1.2k | 663 |
| L2 document | 206 | 320 | 458 | 379 |

<details><summary>Race detail — PostgreSQL, 1 node</summary>

| Design | Tier | Conflicts/hold | Map reads/hold | Engine retries | Gave up | Hold p50/p99 ms | Confirm p99 ms | Deferred ok | Final buyer sold |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|
| E1 sweeper | 10 seats | 8.99 | 16.55 | 0 | 0 | 3.31 / 46 | 46 | 14/14 | 0 |
| E1 sweeper | 100 seats | 3.06 | 5.64 | 0 | 0 | 3.06 / 61 | 66 | 11/11 | 0 |
| E1 sweeper | 1k seats | 1.38 | 4.01 | 0 | 0 | 3.95 / 73 | 81 | 49/49 | 0 |
| E1 sweeper | 10k seats | 0.39 | 3.68 | 1 | 0 | 5.85 / 91 | 94 | 436/436 | 0 |
| L2 document | 10 seats | 11.87 | 22.87 | 1614 | 0 | 1.51 / 28 | 222 | 9/9 | 0 |
| L2 document | 100 seats | 4.34 | 6.75 | 2227 | 0 | 7.80 / 75 | 411 | 18/18 | 0 |
| L2 document | 1k seats | 1.84 | 4.24 | 2985 | 0 | 3.30 / 179 | 648 | 50/50 | 0 |
| L2 document | 10k seats | 0.64 | 4.53 | 12216 | 0 | 3.27 / 306 | 1211 | 433/433 | 0 |

</details>

## Reads

Each question in isolation, fresh key per execution, ops/s.

### PostgreSQL, 1 node

| Question | E1 sweeper | L2 document |
|---|---:|---:|

## Isolated writes, publishing and storage

### PostgreSQL, 1 node

| Design | Hold+confirm /s | Release /s | Refund /s | Publish p50 ms by size | Storage | Audit |
|---|---:|---:|---:|---:|---:|---|
| E1 sweeper | — | — | — |  | 15.2 MiB | consistent |
| L2 document | — | — | — |  | 11.9 MiB | consistent |


## Controlled pairs — one decision at a time


## Measurement hazard: CPU-quota throttling

Cells show *throttled periods / periods (throttled time)*: for the **database** container(s) over
the whole cell (summed across nodes), and for the **client** per phase.

### PostgreSQL, 1 node

| Design | database, whole cell | client lifecycle#1 | client race#1 |
|---|---:|---:|---:|
| E1 sweeper | 180/575 (80s) | 0/385 (0.0s) | 0/175 (0.0s) |
| L2 document | 266/732 (106s) | 0/391 (0.0s) | 2/330 (0.0s) |


---

Plans for every read and write statement are in `results/dc4-pg-e1-l2/<topology>/plans/<design>.txt`
(`EXPLAIN (ANALYZE, BUFFERS)` on PostgreSQL, `EXPLAIN (ANALYZE, DIST)` on YugabyteDB; write
statements inside rolled-back transactions).

## Provenance: code, data and analyses

This file is **generated from the measurements** and contains no interpretation.
Conclusions live in separate signed analyses, so several analysts — including
different AI models — can read the same numbers and each record what they make of
them. Where two analyses disagree, the disagreement is a finding and stays visible.

| | |
|---|---|
| Run id | `devchecks/dc4-pg-e1-l2` |
| Result files | 2 |
| **Inputs digest** | `3c99d4eb048cfcc6` |
| Repository version | `repo/platform-lock-not-available-3-g1b2801d-dirty` (`1b2801d6d10f`), ⚠️ **dirty** — uncommitted changes; not reproducible from the commit alone |

Quote the digest in any analysis: it identifies the data. The repository version
identifies the code — check it out (`git checkout <commit>`) to reproduce the run, or
`git diff <commit>` to see what has changed in the study since.

### Analyses of this run

**None yet.** This run has measurements but no interpretation.

### Adding an analysis

Copy `reports/analyses/TEMPLATE.md`, fill in the frontmatter, and regenerate this report.
If an analysis already exists for digest `3c99d4eb048cfcc6`, read it first and write a new one only to
**disagree, extend, or bring a different perspective**. Never edit another analyst's file.
