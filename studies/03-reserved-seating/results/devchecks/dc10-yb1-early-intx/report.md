# Study 03 — results: reserved seating (venue → event → seat)

Generated from `results/dc10-yb1-early-intx`. Every number comes from exactly one JSON file in that
directory. **This file contains measurements and no interpretation**; conclusions live in
signed analyses, indexed at the end.

## TL;DR — measured facts

*Selected from the tables below by fixed rules; no interpretation. Read a signed analysis for
what these facts mean.*

- **Cells:** 1 across 1 topologies; 0 failed.
- **Negative controls:** 0 of 0 fired as designed to.
- ❌ **Invariant violations in designs meant to be correct:** S1 conditional on YugabyteDB, 1 node (RF=1) (2 early rejections).
- **Race, YugabyteDB, 1 node (RF=1)** (seats held/s; designs meant to be correct, without violations):

## What was measured

|  |  |
|---|---|
| Run id | `devchecks/dc10-yb1-early-intx` |
| Environment | [`host-zenbook-ux5406sa`](../../../docs/environments/host-zenbook-ux5406sa.md) |
| Scale | `tiny` — 5 venues, 84 events, 5741 tickets sold, 131 live and 34 expired holds at load |
| Events (kind_seats: count) | catalogue_10: 40 · catalogue_100: 8 · catalogue_1000: 2 · catalogue_10000: 1 · lifecycle_10: 4 · lifecycle_100: 2 · lifecycle_1000: 1 · race_10: 20 · race_100: 4 · race_1000: 1 · race_10000: 1 |
| Seed | 42 — identical data in every cell |
| Reads / isolated writes | 8 workers, 5s measured + 2s warmup, 1 trial(s) |
| Race | 32 buyers per event, real hold TTL 40m0s, 10% of holds deferred until the crowd has gone, timeout 20s, tier budget 3m0s, give up after 20 conflicts |

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
| S1 conditional | ✅ 51/51 |

## Negative controls

Deliberately wrong designs (and E1 with its sweeper stopped). Each must be seen to break the
invariant it exists to test; otherwise the absence of violations elsewhere in that experiment
is not evidence of correctness.

No control was measured in this run.

## The race — a hot drop with seat choice

Seats held per second, from the gate to the last hold granted. Conflicts and seat-map reads are
per granted hold. Deferred confirmations kept a real 40-minute hold through the crowd.

### YugabyteDB, 1 node (RF=1)

| Design | 10k seats |
|---|---:|
| S1 conditional | ~~81.2 (6 trials, spread 15%)~~ ❌ 2 early rejections; ⏱ 6 timed out at 17% sold |

<details><summary>Race detail — YugabyteDB, 1 node (RF=1)</summary>

| Design | Tier | Conflicts/hold | Map reads/hold | Engine retries | Gave up | Hold p50/p99 ms | Confirm p99 ms | Deferred ok | Final buyer sold |
|---|---|---:|---:|---:|---:|---:|---:|---:|---:|
| S1 conditional | 10k seats | 0.27 | 1.27 | 0 | 0 | 97 / 279 | 700 | 414/414 | 0 |

</details>

## Reads

Each question in isolation, fresh key per execution, ops/s.

### YugabyteDB, 1 node (RF=1)

| Question | S1 conditional |
|---|---:|

## Isolated writes, publishing and storage

### YugabyteDB, 1 node (RF=1)

| Design | Hold+confirm /s | Release /s | Refund /s | Publish p50 ms by size | Storage | Audit |
|---|---:|---:|---:|---:|---:|---|
| S1 conditional | — | — | — |  | n/a on YugabyteDB | consistent |


## Controlled pairs — one decision at a time


## Measurement hazard: CPU-quota throttling

Cells show *throttled periods / periods (throttled time)*: for the **database** container(s) over
the whole cell (summed across nodes), and for the **client** per phase.

### YugabyteDB, 1 node (RF=1)

| Design | database, whole cell | client race#1 | client race#2 | client race#3 | client race#4 | client race#5 | client race#6 |
|---|---:|---:|---:|---:|---:|---:|---:|
| S1 conditional | 1340/1989 (592s) | 0/205 (0.0s) | 0/205 (0.0s) | 0/206 (0.0s) | 0/206 (0.0s) | 0/206 (0.0s) | 0/207 (0.0s) |


---

Plans for every read and write statement are in `results/dc10-yb1-early-intx/<topology>/plans/<design>.txt`
(`EXPLAIN (ANALYZE, BUFFERS)` on PostgreSQL, `EXPLAIN (ANALYZE, DIST)` on YugabyteDB; write
statements inside rolled-back transactions).

## Provenance: code, data and analyses

This file is **generated from the measurements** and contains no interpretation.
Conclusions live in separate signed analyses, so several analysts — including
different AI models — can read the same numbers and each record what they make of
them. Where two analyses disagree, the disagreement is a finding and stays visible.

| | |
|---|---|
| Run id | `devchecks/dc10-yb1-early-intx` |
| Result files | 1 |
| **Inputs digest** | `a9ca693d13c06822` |
| Repository version | `repo/platform-lock-not-available-4-g1024aa7-dirty` (`1024aa73d282`), ⚠️ **dirty** — uncommitted changes; not reproducible from the commit alone |

Quote the digest in any analysis: it identifies the data. The repository version
identifies the code — check it out (`git checkout <commit>`) to reproduce the run, or
`git diff <commit>` to see what has changed in the study since.

### Analyses of this run

**None yet.** This run has measurements but no interpretation.

### Adding an analysis

Copy `reports/analyses/TEMPLATE.md`, fill in the frontmatter, and regenerate this report.
If an analysis already exists for digest `a9ca693d13c06822`, read it first and write a new one only to
**disagree, extend, or bring a different perspective**. Never edit another analyst's file.
