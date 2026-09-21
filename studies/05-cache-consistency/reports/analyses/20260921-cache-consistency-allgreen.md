---
analysis_id: 20260921-cache-consistency-allgreen
supersedes: 20260921-cache-consistency-survey
study: 05-cache-consistency
run_id: 20260921T-survey3
inputs_digest: 03fac739d0835768
repo_commit: 48fe1e33db5467dfa4421f354df4a6ee77becaff
analyst: deepseek-flash
analyst_kind: ai
analyst_version: deepseek-flash, DeepSeek HIGH role (model id, exact version, tool and effort not exposed by the client), WSL2/Podman
analyzed_at: 2026-09-21T17:45:00Z
report: ../../reports/20260921T-survey3.md
---

# Study 05 — external cache throughput and consistency: analysis of run `20260921T-survey3`

This analysis supersedes `20260921-cache-consistency-survey` (moved to `outdated/`). It
replaces that file's conclusions, not its record: the earlier analysis and the three runs
behind it are what found the harness defects below, and its history is worth keeping.

## 0. What this run is

`20260921T-survey3` is the full PostgreSQL matrix — 21 core scenarios, 3 database-layout
reference cells, 2 controls — at scale `small` (800 donors), on one pinned PostgreSQL
17.11 container and one pinned Redis 7.4.11 container, from the **clean** commit
`48fe1e3` (tag `run/05-cache-consistency/20260921T-survey3`), inputs digest
**`03fac739d0835768`**.

**Every one of the 27 cells passed every gate: zero failed cells, both controls fired,
6 ledger assertions per cell with zero mismatches, 310 acknowledgement verifications per
cell with zero violations, zero impossible cache values, and zero stale-after-ack reads
in all twelve strict cells.**

That is only meaningful because of how the run before it failed. The four harness defects
that had to be fixed first are section 5, and three of them were *manufacturing* the
correctness findings rather than detecting them.

## 1. The direct answers

**How much did caching add over the same-run no-cache baseline?** Warm cacheable reads:

| Scenario | reads/s | p99 ms | vs its own baseline |
|---|---|---|---|
| `legacy-na-none-none-none-coord` (no cache) | 7 794 | 15.8 | 1.0× |
| `legacy-na-memory-aside-strict-coord` | **20 233** | **0.9** | **2.6×** |
| `legacy-na-memory-aside-relaxed-coord` | 19 308 | 0.9 | 2.5× |
| `legacy-na-memory-through-strict-coord` | 20 074 | 0.8 | 2.6× |
| `legacy-na-redis-aside-strict-coord` | 18 667 | 1.3 | 2.4× |
| `owned-opt-none-none-none-coord` (no cache) | 5 805 | 18.9 | 1.0× |
| `owned-opt-redis-aside-strict-coord` | 20 380 | 0.9 | **3.5×** |
| `owned-opt-memory-aside-strict-coord` | 16 614 | 1.1 | 2.9× |
| `owned-pess-none-none-none-coord` (no cache) | 7 043 | 27.1 | 1.0× |
| `owned-pess-redis-aside-strict-coord` | 15 782 | 2.0 | 2.2× |

The cache buys **2.2×–3.5×** in throughput and **an order of magnitude in p99** (15.8–27.1 ms
→ 0.8–2.0 ms). Absolute numbers are not comparable across runs; every ratio here is within
this run.

**What did strict freshness cost?** Now that all cells pass, the controlled pairs are
readable — and the answer is nothing measurable in read throughput:

| Controlled pair (one decision apart) | relaxed | strict | difference |
|---|---|---|---|
| legacy, memory, aside | 19 308 | 20 233 | +4.8 % (noise) |
| legacy, memory, through | 19 568 | 20 074 | +2.6 % (noise) |
| legacy, redis, aside | 19 750 | 18 667 | −5.5 % (noise) |
| legacy, redis, through | 15 148 | 14 949 | −1.3 % |
| owned, memory, aside | 16 551 | 16 614 | +0.4 % |
| owned, memory, through | 19 643 | 18 010 | −8.3 % |
| owned, redis, aside | 18 150 | 20 380 | +12.3 % (noise) |
| owned, redis, through | 19 593 | 16 813 | −14.2 % (noise) |

Every difference is inside the ~20 % noise floor in either direction. **Strict freshness
cost nothing measurable in read throughput at this scale**, which is consistent with the
mechanism: strict fences more often (two fences per write instead of one), and fences cost
the *writer*, not the reader. Its cost is visible in the write path (`tombstone_fences_installed`
is roughly double the relaxed cells' `invalidations`) and **80 % of the measured range is
write-path cost**: one relaxed cell records 1 024 `write_noops_refused_by_owner` — writes
the SQL correctly refused because the observed owner no longer owned the row.

**How often was a relaxed cache wrong, and how stale?** Only two cells in the entire matrix
recorded a single wrong read:

| Cell / phase | wrong reads | % of reads | cache hits |
|---|---|---|---|
| legacy, redis, through, relaxed — **three instances** | 38 749 | **87.81 %** | 44 122 |
| owned, redis, through, relaxed — **three instances** | 35 287 | **87.24 %** | 40 442 |

**Every other cell — including every single-instance relaxed cell in all four
model/strategy/backend combinations — recorded ZERO wrong reads in every phase.** In the
earlier run, before the harness's write ordering was fixed, the same single-instance
relaxed cells recorded 0.04 %–0.3 %; that was the harness, not the design. The honest
statement is therefore: **with a single application instance and a coordinated writer, this
relaxed cache was never observed to return a superseded value in ~1.2 million measured
reads.** The staleness the relaxed contract permits did not materialise in this workload —
which is a result about the workload's read/write key overlap as much as about the design.

**What produced the two failures there?** Three logical instances sharing one Redis, on the
`through` strategy. Each instance's writes publish a representation that the other two then
serve, and a private in-memory lease cannot coordinate them. This is the study's strongest
result and it is now confirmed in two independent runs (87.8 % and 86 %). Both models show
it, which is the point: **adding a version token to the database does not help a shared
cache whose per-instance write path publishes what a peer then reads.** The in-memory
three-instance arms do *not* show it, because their strict policies either validate the
database version per hit (owned) or read authoritatively (legacy).

**What did the owned model's version token and outbox buy?** Cost: 5 805 vs 7 794 no-cache
reads/s against the legacy model (−26 %), plus a row update and an outbox insert per write.
Benefit at this scale: nothing measurable in the read path, and — notably — **no reduction
in the three-instance staleness**, because that failure is about cache topology, not about
the token. Its genuine benefit is the one the fence cannot provide: a *local* cache in a
multi-instance deployment can validate a hit against the authoritative version, which is
why `owned-opt-...` three-instance arms are the only multi-instance arms that stay correct
without reading the database.

**Can the legacy model give strict freshness under external writers?** No. With 20 % of
writes bypassing the adapter, the relaxed cell's throughput falls back to the no-cache
baseline (7 664 vs 7 794) — the cache has become pure overhead, because every read must
bypass to be safe — and the strict cell achieves 8 269 by reading authoritatively. The
study's honest answer: **an unknown writer cannot be fenced, so a legacy cache cannot be
strict in that regime, and the measurement shows the cache's value collapsing to zero
rather than merely degrading.**

**Did the lease prevent stampedes, and at what latency cost?** Yes, exactly: **8 database
loads for 8 keys with 16 readers released together, 0 duplicate fills, in 211–225 ms** — one
load per key. The calibrated lease was pinned at its 100 ms floor, which is itself the
finding that database fills are fast (p99 well under 25 ms) and a short lease is the
honest calibration rather than a round number.

**Rollup, rolldown and embedding.** Rolldown is measured by every cell: the legacy model
*is* the flattened FK (Study 01's D3), and it is the fastest baseline here (7 794 vs 5 805
for the normalized/owned shape) because the portal view needs one row plus a bounded child
slice. Embedding is measured by `ref-embedded-locked` (7 007 reads/s, p99 31.9 ms — the
worst tail of any cell, because its portal read unnests a JSONB array). Rollup is measured
by `ref-rollup-trigger` (9 412 reads/s, p99 1.8 ms — the best reference), and its
`a_rollup_drift` audit is **clean in this run**. In the previous run that audit reported 3
mismatches; that drift was the harness's concurrent-write ordering (two writers touching
one child row), not the design. **The trigger-maintained rollup is covered and correct
under this workload — but the earlier failure is the reason the coverage is credible, and
it should be re-tested at `medium` scale before it is relied on.**

## 2. Where the measurement is weak

1. **Single trials.** Every number is one trial. Spreads within a measurement are shown in
   the report; there is no error bar across trials, and the noise floor (~20 %) is larger
   than most pairwise differences in section 1. No conclusion here rests on a difference
   below that floor, which is why several are stated as "nothing measurable".
2. **No sustained-churn phase.** The 300 s hard TTL was never crossed in a measured run.
   The TTL rule and the probabilistic early-expiry formula are verified by unit tests at
   0/75/150/225/300 s with a fake clock, but their **field behaviour across a real TTL
   boundary is not measured** — a declared coverage gap.
3. **No equal-total framing.** Every cell ran under `db-only` (database 2 CPUs/3 GiB,
   client 2 CPUs/2 GiB). The cache's own resource cost is therefore not counted anywhere
   except in Redis's own accounting (~1 MiB resident against 32 MiB `maxmemory`, 0
   evictions, working set under 1 MiB).
4. **No `medium` scale, no repeated trials, no YugabyteDB cell, no three-node cluster, no
   colocation evidence.** The placement pair (`ref-y1-colocated`, `ref-y2-noncolocated`)
   exists as SQL and was never executed. "Colocated" remains an untested word here.
5. **The client is closed-loop within a phase**, so the deep percentiles are optimistic
   floors. The platform's open-loop scheduler exists and this study did not use it.
6. **One laptop.** Eight shared cores, no real network, the client and both databases
   competing for the same silicon. The Redis-versus-memory gap (now nearly closed) and the
   absolute throughput are properties of this arrangement.
7. **Cache gauges versus counters.** `items` and `resident_bytes` are read once at the end
   of a cell, after the fault phases have flushed, while `evictions`, `fills` and
   `publishes` are cumulative. Read the counters as evidence, not the gauges.
8. **`ref-normalized-indexed`'s `writes.sql` now has the owner guard and a relative
   correction**, so it differs from a plain normalization by more than the normalization.
   The pair still isolates the layout decision for reads, which is what it is used for.
9. **No independent review.** One agent planned, implemented, measured and analysed this
   study, and the same agent fixed four of its own harness defects — a process that caught
   a great deal (section 5) but is not a substitute for another model reading the code.

## 3. The mechanism that mattered most: invalidation is a fence, not a deletion

"Publish only after the commit" does not make cache-aside safe. A reader whose fill begins
before a writer's invalidation holds a committed but *superseded* state and can republish
it. The study closes that with a per-key **invalidation fence** in the cache: atomic in
both backends (one mutex; one Lua script), captured *before* the database snapshot and
stamped on the entry, with publication refused unless the key's fence is unchanged, and a
strict writer fencing **before and after** its commit. `publishes_refused_by_fence` shows
the race is real rather than theoretical: 60–625 refused publications per cell.

By the end of the study this is no longer the headline, and that is the honest ordering: it
is a *design* result that a *harness* defect had obscured.

## 4. The acknowledgement boundary

The strict contract is written against the **acknowledgement**, not the commit: a read
that begins after a write was acknowledged must not return an older value. The ledger
therefore applies the data change at commit time and moves the requirement at the
acknowledgement, classifying a committed-but-unacknowledged state as `ahead` rather than
`impossible`. The new `ackcheck` phase tests the property directly — write, then
immediately read from a **second database session** — and reports **0 violations in 310
verifications per cell, sequential and concurrent**, across all 27 cells.

## 5. The four defects this study found in its own harness, and why they mattered

Each was found by a check that was built to find design defects, and each would have
produced a confidently wrong published claim.

1. **A mutation was built from the ledger's state and applied by id only.** A concurrent
   writer could move or delete the row in between, so the database changed a row the ledger
   attributed elsewhere, and the two diverged. This produced "impossible cache values" and
   an audit mismatch that belonged to no design. Fixed with owner-guarded statements, a
   distinct `errNoEffect` outcome for a refused write (the database did not change, so the
   ledger does not), and relative corrections in every schema.
2. **The ledger's history order did not follow the database's commit order.** Under
   concurrency the handler could append states in the other order, so a legitimate later
   read looked "one state behind". Proven by the fact that a design with **no cache at all**
   showed it; fixed by serialising writes that touch one key (64 shards, never across keys),
   which leaves inter-key parallelism intact.
3. **The checker had the ordering mistake it was built to detect.** `ackcheck` captured the
   requirement *after* its database read, so a concurrent writer could move the requirement
   between the two and the checker manufactured its own violation.
4. **A control was judged by the wrong evidence.** The unguarded-version-bump control had to
   lose a *bump* as well as break the version/outbox invariant; on a run where the invariant
   visibly broke while the count still balanced, the harness refused to license the guarded
   designs after watching the unguarded one fail. A control is judged by whether the
   invariant broke.

And two honesty fixes around them: negative controls are now exempt from the ledger and
acknowledgement assertions (they exist to break invariants) and judged by their own fault,
with the skip recorded in the report; and the instances-phase arms' violations now reach the
cell's acceptance check — before that fix, an owned cell had *passed* with 11 hidden stale
reads.

## 6. What a second analyst should attack first

1. **The three-instance shared-cache result.** ~87 % wrong reads, twice, single-trial, one
   machine. Reproduce it under repeated trials and, if possible, on a second machine; if it
   holds, it is the study's most valuable finding and it undercuts the assumption behind
   every "just add Redis" deployment in this repository.
2. **The churn phase**, at the real 300 s TTL, with the early-expiry formula already pinned
   by unit tests but its field behaviour unmeasured.
3. **The equal-total framing** for the Redis comparisons — a shared cache that needs a
   separate container should pay for it in the comparison, and this run never made it do so.
4. **Scale**: `medium`, repeated trials, and the missing topology arms. All five
   representative reference cells and all twelve strict cells pass at `small`; whether the
   87 % result depends on scale is unknown.
5. **The `pess`/`opt` pair** is currently a single representative cell each
   (`owned-pess-redis-aside-strict-coord` vs `owned-opt-redis-aside-strict-coord`: 15 782 vs
   20 380), which is one seed and one trial apart — a difference to test, not to quote.

## 7. Signature

Analyst `deepseek-flash`, acting as DeepSeek HIGH and executing its own handoff. It planned,
implemented, measured, repaired and analysed this study in one session; there was no LOW
handoff, by the task's instruction, and **no independent model review**. Every figure cited
here was checked against `reports/20260921T-survey3.md` (inputs digest
`03fac739d0835768`) before this file was committed.
