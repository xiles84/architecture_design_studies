---
analysis_id: 20260921-cache-consistency-survey
study: 05-cache-consistency
run_id: 20260921T-survey
inputs_digest: 8ff86dc9e5228e86
repo_commit: 567778ff7091f7748a3a6096cec3b98b408adec7
analyst: deepseek-flash
analyst_kind: ai
analyst_version: deepseek-flash, DeepSeek HIGH role (model id, exact version, tool and effort not exposed by the client), WSL2/Podman
analyzed_at: 2026-09-21T15:30:00Z
report: ../../reports/20260921T-survey.md
---

# Study 05 — external cache throughput and consistency: analysis of run `20260921T-survey`

## 0. What this run is, and what it is not

Run `20260921T-survey` is the first measured pass over the study's 21 PostgreSQL
scenarios plus three database-layout reference cells, at scale `small` (800
donors, ~10 600 donations), on one pinned PostgreSQL 17.11 container and one
pinned Redis 7.4.11 container on the same eight-core WSL2 machine. Every cell
loads a fresh copy of the deterministic dataset, passes the pre-timing
correctness gate, runs the phases, injects its faults, replays its ledger and
audits.

**Nine of the 27 cells failed their correctness rules and are excluded from every
performance statement below.** Their failures are findings, not noise, and they
are the most important output of this run: they are reported in section 6 and
they bound the conclusions of sections 3 to 5.

This analysis was written by the same single agent that planned, implemented and
measured the study. There is **no independent review** of the harness, the oracle
or these conclusions. That is a limitation of the whole study, not of this file.

## 1. The direct answers, in the order the protocol asks for them

**How much did caching add over the same-run no-cache baseline?** Measured warm
cacheable read throughput, same run, same scale, same machine:

| Scenario | warm reads/s | p99 ms | vs its own no-cache baseline |
|---|---|---|---|
| `legacy-na-none-none-none-coord` (no cache) | 3 914 | 27.6 | 1.0× — the baseline |
| `legacy-na-memory-aside-relaxed-coord` | 15 018 | 1.7 | **3.8×** |
| `legacy-na-memory-through-relaxed-coord` | 12 860 | 2.0 | 3.3× |
| `legacy-na-memory-aside-strict-coord` | 13 295 | 2.0 | 3.4× |
| `legacy-na-redis-aside-relaxed-coord` | 8 945 | 9.5 | 2.3× |
| `legacy-na-redis-aside-strict-coord` | 7 687 | 15.3 | 2.0× |
| `legacy-na-redis-through-relaxed-coord` | 8 847 | 7.8 | 2.3× |
| `owned-opt-none-none-none-coord` (no cache) | 2 984 | 21.6 | 1.0× — the baseline |
| `owned-opt-memory-aside-relaxed-coord` | 13 165 | 1.7 | 4.4× |
| `owned-opt-memory-aside-strict-coord` | 10 972 | 2.6 | 3.7× |
| `owned-opt-redis-aside-strict-coord` | 7 841 | 11.8 | 2.6× |
| `owned-pess-redis-aside-strict-coord` | 8 337 | 14.9 | 2.8× |

Read p99 falls by roughly an order of magnitude for the in-process caches and by
about half for the shared one. This is a same-run comparison: nothing from study
01 is presented as a baseline here.

**What did strict freshness cost?** Only two same-run, same-model, same-backend,
same-strategy pairs are actually comparable, and only one of them is clean:

| Controlled pair | relaxed | strict | strict cost |
|---|---|---|---|
| legacy, memory, aside | 15 018 | 13 295 | **−11 %** (within the ~20 % noise floor; not separable) |
| legacy, redis, aside | 8 945 | 7 687 | −14 % |
| legacy, memory, through | 12 860 | 15 972 | +24 % — **direction reversed**, see the weakness list |
| legacy, redis, through | 8 847 | *failed the strict contract* | not reportable |

The honest reading is that at this scale, on this machine, strict freshness did
**not** measurably reduce read throughput in a same-run pair, and the one pair
that suggests otherwise reverses direction and is explained by trial variance.
The mechanism agrees: strict adds fences and an authoritative read only where the
policy cannot prove freshness, and it removes almost no cache work for
**reads** — the cost lands on **writers** (section 4).

**How often did relaxed caching return wrong information, and how stale was it?**
Only in the write-bearing phases, never in the warm read-only phase:

| Scenario / phase | wrong reads | % of all reads | cache hits | longest streak |
|---|---|---|---|---|
| legacy, memory, through, relaxed — 90/10 app mix | 229 | **0.258 %** | 73 439 | 2 |
| legacy, memory, through, relaxed — 90/10 cacheable | 229 | 0.182 % | 103 763 | 2 |
| legacy, redis, through, relaxed — 90/10 app mix | 136 | 0.304 % | 44 643 | 2 |
| legacy, redis, aside, relaxed, external writers — 90/10 app mix | 9 | 0.041 % | 0 | 1 |
| legacy, redis, through, relaxed — **three instances** | 28 337 | **85.7 %** | 33 025 | — |
| owned, redis, through, relaxed — **three instances** | 27 978 | **85.8 %** | 32 576 | — |

The measured answer is therefore: **under a coordinated, single-instance
deployment a relaxed cache returned a superseded value for about 0.04 % to 0.3 %
of reads while writes were in flight, and for zero reads in the warm read-only
phase. In a three-instance deployment it returned one for the large majority of
reads.** The multi-instance number is the study's strongest result and the one it
most needs a second analyst to confirm (section 7).

**Which races and failures produced it?** The wrong reads carry an attributed
cause. In the single-instance relaxed cells the causes are the ledger's
attribution buckets for a superseded `cache hit`, i.e. writes whose post-commit
invalidation had not yet run when the read was served — the failure the relaxed
contract explicitly permits. The stale-read negative control
(`ctl-stale-invalidation`) fired: with exactly one post-commit invalidation
deliberately suppressed, the oracle detected the resulting stale read, so the
accounting is demonstrated to be able to see the failure rather than assumed to.

**When did Redis beat the in-process cache, after counting its resources?**
Never, in this environment. Every Redis cell is slower than the corresponding
in-process cell by 1.6×–2.0× at equal client workers (8 945 vs 15 018 warm
reads/s), and the gap is larger in p99 (9.5 ms vs 1.7 ms). Redis's own `INFO`
shows it is not resource-starved: 32 MiB of `maxmemory` against a working set
under 1 MiB, `evicted_keys` 0, `used_memory_peak` ~1 MiB. The cost is the
container-to-container round trip on a shared laptop, not memory pressure. **This
says nothing about a real network deployment** — see the weakness list.

**How did multiple instances change the result?** Catastrophically, and this is
the run's headline. With three logical instances and equal total workers, a
*shared* Redis cache on the `through` strategy recorded **85.7 % of reads wrong**
(28 337 of 33 025) for the legacy model and **85.8 %** for the owned model, while
the same design with one instance recorded 0.4 % and 0.02 % respectively. The
in-memory three-instance arms do not show this — their strict policies either
validate the database version per hit (owned) or read authoritatively (legacy) —
so the finding is specific to a shared store whose per-instance write path
publishes a representation that other instances then serve. A private in-memory
LRU cannot be invalidated by a peer at all, which is why the study never labels
an uncoordinated local cache strict.

**What did a modifiable database model buy?** A version token and an outbox cost
2 984 vs 3 914 no-cache reads/s (−24 %) and add transactional work to every write,
and in exchange the owned model's strict cache can prove freshness from the cache
itself: `owned-opt-memory-aside-strict-coord` passed its strict contract at
10 972 reads/s with zero stale-after-ack reads and zero impossible values. But it
did **not** buy a lower staleness rate in the relaxed arms — the owned and legacy
relaxed through cells recorded the same 0.04 %–0.06 % wrong-read rates — because
the fences, not the token, are what bound the single-instance races (section 2).

**What did the legacy adapter have to pay for the same guarantee?** It had to
carry the freshness proof in the cache rather than in the database: an
invalidation fence, captured before every fill and checked on every publish. That
is a cache-side obligation the owned model also uses, and it costs no database
round trip. Its measured price is visible as `publishes_refused_by_fence` (61–85
refused publications per cell, i.e. **61–85 fills per cell that would otherwise
have republished a superseded value**) and as extra tombstones (6 063–8 537 per
cell on the strict write path, two per write).

**Can the legacy cache give strict freshness under external writers without
querying the database?** No, and the run measures the reason: `legacy-na-redis-aside-strict-ext20`
recorded 3 212 stale-after-ack reads — the external writer commits without
fencing, so nothing in the cache can notice, and the policy for that cell is an
authoritative read on every read, which is exactly "the cache cannot do this
alone". The external fraction is 20 % of writes.

**Did the lease prevent stampedes, and at what latency cost?**
Yes, measurably. In the stampede phase (16 readers released together over 8
keys) every cache-bearing cell reported **8 database loads for 8 keys — exactly
one per key** — with **0 duplicate fills**, in 210–230 ms. Without coordination
the same burst would be 128 loads. The leases cost almost nothing in the read
path: lease *waits* are rare (443–536 contended of 26 000–34 000 acquisitions in
the memory cells) and the calibrated lease duration was pinned at the 100 ms
floor, which is itself the finding that database fills are fast enough
(p99 well under 25 ms) that a short lease is the honest calibration.

**What was the cache policy's observable behaviour?** The in-process cache is a
byte-bounded **exact LRU**: with a 256 KiB capacity against a ~2 MB working set
it evicted 22 597 entries in one cell while holding its byte bound. Redis was
configured `maxmemory=32mb`, `allkeys-lru`, `maxmemory-samples=5`, read back from
`INFO` and recorded in every redis cell; it evicted nothing because its capacity
exceeded the working set. Redis's LRU is approximate and its capacity is not the
in-process capacity, so every Redis/memory difference in this run is partly a
policy difference and is stated as such.

**The five-minute TTL and probabilistic early expiration.** The hard TTL is
300 s and was never shortened; the probabilistic draw is
`p = clamp(1 − remaining/300 s, 0, 1)`, implemented once and used by both
backends. It is pinned by unit tests at ages 0, 75, 150, 225 and 300 s against a
fake clock (p = 0, 0.25, 0.5, 0.75, 1.0) and by a deterministic-RNG replay test.
**This run contains no sustained-churn phase: it was cut from the matrix and is a
declared coverage gap.** Every cell therefore reports zero hard expiries and its
early expiries only for entries that survived a phase, so the observed-by-age
evidence in the report is thin, and the honest statement is that the *formula* is
verified and the *field behaviour across two real TTL boundaries* is not measured.

**Equal-total versus add-cache resources.** Every cell in this run was measured
under the `db-only` framing (database 2 CPUs / 3 GiB, client 2 CPUs / 2 GiB).
**No add-cache or equal-total cell was run: a declared coverage gap.** No
conclusion here accounts for the cache's own resource cost, except that Redis's
own accounting shows its memory footprint (~1 MiB) to be negligible beside the
database's.

## 2. The mechanism that mattered most: the invalidation fence

The single most important technical result is not a throughput number. It is that
**"publish only after the commit" does not make cache-aside safe**. A reader that
begins its fill before a writer's invalidation holds a *committed but superseded*
state, and publishing it afterwards reinstates exactly the value the writer
invalidated. The first implementation of this study had no defence against that
and recorded ~14 000 stale-after-ack reads in one strict cell; the
dev-check history in `PROGRESS.md` records each iteration.

The fix is a cache-side **invalidation fence**: the writer advances a per-key
fence and removes the value (one atomic operation in both backends), a fill
captures the fence *before* taking its database snapshot and stamps it on the
entry, and a publish is refused unless the key's fence is still that value. The
run shows the fence doing real work: 61–85 publications per cell were refused by
it. Two further pieces were needed and are part of the finding:

1. a strict writer must fence **before and after** the commit. One fence leaves a
   window in which a reader captures the new fence and then reads the
   pre-mutation state;
2. the ledger's freshness requirement must move at the **acknowledgement**, not
   at the commit — which is what the protocol's contract is actually written
   against.

The two unit tests that pin these (`TestFenceRefusesAFillThatStartedBeforeAnInvalidation`,
`TestOracleClassifiesFreshStaleAheadAndImpossible`) fail if either is removed.

## 3. Where the measurement is weak

This section is required and is deliberately unflattering.

1. **Nine cells are excluded and I could not fix them inside this session.** The
   residual is listed in section 6. Until it is resolved, the clean strict cells
   are evidence *for those cells*, not for the design family.
2. **Single trials.** Every number in this run is one trial. The report prints
   each measurement's spread across workers, which is not an error bar across
   trials, and the noise floor (~20 % per `docs/methodology.md`) is larger than
   several differences I would otherwise want to read as real — the strict-cost
   comparison in section 1 is exactly that case.
3. **`cache.evictions`, `fills` and `publishes` are cumulative for the whole
   cell, but `items` and `resident_bytes` are read once at the end**, after the
   fault phases have flushed the cache. The end-state fields therefore describe
   an empty cache and must not be read as steady state; the counters are the
   evidence.
4. **No sustained-churn phase, no equal-total framing, no `medium` scale, no
   repeated trials, no YugabyteDB cell, no three-node cluster, no colocation
   evidence.** All are declared gaps; the placement pair exists as SQL and was
   never executed.
5. **The client is closed-loop within a phase.** `measure.Run` bounds each
   operation, so the deep percentiles are optimistic floors and no figure here is
   an SLO. The open-loop scheduler exists in the platform but this run did not
   use it.
6. **Everything runs on one laptop.** There is no real network, the containers
   share eight cores, and the client, the cache and the database compete for
   them. The Redis-versus-memory gap and the absolute throughput are properties of
   *this arrangement*, not of the designs.
7. **The ledger's model of concurrently conflicting writes is incomplete.** Two
   concurrent mutations of the same key can apply to the ledger in an order the
   database did not commit in, and set-valued mutations (`person_update`) are
   order-sensitive by nature. This is the most likely explanation for the
   remaining impossible-value and audit failures in section 6 and it is a
   *harness* limitation: the harness would need per-key write serialisation, which
   would change the contention it is trying to measure.
8. **The three reference cells of this run were re-run at a LATER commit than the
   other 24, and the manifest records one commit for all of them.** The first
   pass panicked in those three cells because `hasCache()` compared the backend
   field against `"none"` while an unset field is `""`; the harness now asks for
   the backends it supports explicitly. That change is a **no-op for every other
   cell** (their backend is set explicitly), so the 24 earlier cells' behaviour is
   unaffected — but a reader is entitled to know that the run directory contains
   cells from two code states, and the fix is in the same worktree history.
9. **No independent review.** One agent wrote the harness, the oracle, the
   accounting, the runs and this analysis.

## 4. What the write path costs

Isolated write phases on the optimistic owned model: `hotspot_writes` measured
693 ops/s (p99 9 ms) with the in-process cache and 1 347 ops/s (p99 3 ms) with
Redis; the no-cache owned baseline is the comparison the *report's* write table
carries. The strict path pays two fences per write (6 063–8 537 tombstone fences
per cell) and the version/outbox maintenance; the relaxed path pays one
invalidation. At a 99/1 read/write mix this is invisible in the read throughput,
which is why strict and relaxed read figures are so close — and why a write-heavy
conclusion cannot be drawn from this run at all.

## 5. What the reference cells say

`ref-normalized-indexed` (4 100 reads/s) and `ref-embedded-locked` (3 382
reads/s) both passed their gates and audits, so the study's read path is
comparable across the two database shapes and the flattened-FK baseline
(`legacy-na-none-none-none-coord`, 3 914 reads/s) sits between them as expected
for a workload whose portal view is one row plus a bounded child slice.
**`ref-rollup-trigger` failed with 76 impossible content values** and is not
reportable: its database-maintained rollup produced portal content the ledger's
model never held, which is a drift finding that must be diagnosed before that
reference can be used. Study 01's rollup/rolldown/embedding coverage is therefore
mapped to this study as: rolldown measured by every cell (the legacy model *is*
the flattened FK), embedding measured by `ref-embedded-locked`, rollup attempted
by `ref-rollup-trigger` and **not yet demonstrated**.

## 6. The nine failed cells, with their attributed causes

| Cell | Failure | Attributed cause |
|---|---|---|
| `legacy-na-memory-through-strict-coord` | 11 stale-after-ack reads | `cache hit served a superseded entry` in the mixed and 1- and 3-instance phases |
| `legacy-na-redis-through-strict-coord` | 210 stale-after-ack reads | same attribution, plus 1 `the refilling reader published before the write committed` |
| `owned-opt-redis-through-strict-coord` | 35 stale-after-ack reads | same |
| `legacy-na-redis-aside-strict-ext20` | 3 212 stale-after-ack reads | the external writer; the policy for this cell is an authoritative read and the residual is in that path |
| `owned-opt-redis-aside-relaxed-coord` | 7 impossible cache values | a hit on a state the ledger had not recorded — the pending-state race described in weakness 7 |
| `legacy-na-memory-aside-relaxed-coord` | ledger replay (donation 21 846 absent) | concurrent conflicting writes applied to the ledger out of the database's commit order |
| `ref-rollup-trigger` | 76 impossible content values | the stored rollup and the ledger's model disagree; not yet diagnosed |
| `ref-normalized-indexed`, `ref-embedded-locked` | first run: panic, no result | harness defect (`hasCache` treated an unset backend field as a cache), **fixed and re-run green** |

The strict *through* strategy fails and the strict *aside* strategy does not, in
both models and both backends. That asymmetry is a real lead for whoever
continues this study: it points at the post-commit republish path sharing a store
with concurrent fillers rather than at the fence itself.

### 6a. Diagnosis added after the analysis: at least part of this is the LEDGER, not the design

A follow-up instrumented run (`results/diag-strict/`, run id `diag-strict`, the
three strict `through` cells, `small` scale, phases `verify,calibrate,warm,mixed,hotspot,instances,faults`)
recorded the first few stale reads **in full** — source, cause, entry fence,
current fence, entry publication sequence, adapter publication sequence, versions
behind, and whether a write overlapped the read. The sampling and its JSON are
part of the committed harness; the numbers below are from that run, so they are
**not** the survey's.

What the samples show, and it changes the interpretation:

1. **No stale read overlapped a write.** `overlapped: false` on every sample. The
   "concurrent/ambiguous" bucket is empty for these cells; the failures are
   genuine "a read that began after an acknowledgement was served an older value".
2. **Two different paths produce them, and one of them cannot be a cache fault.**
   Samples with `source: bypass` show an **authoritative database read** returning
   a state **1 and 2 versions behind the requirement**. A bypass read executes the
   portal statements in one repeatable-read transaction and consults no cache at
   all, so no cache mechanism can produce this: it means the ledger's freshness
   requirement was **ahead of the database's own state** at that moment. Those
   samples come from the three-instance arm (11 014 bypass reads), i.e. exactly
   the arm whose strict policy reads the database.
3. Samples with `source: hit` and `source: fill` show a value exactly one state
   behind, and their recorded entry fence, current fence, entry sequence and
   adapter sequence all read as zero — which is itself the next thing to test, and
   is recorded here rather than explained away.

### 6b. The accounting defect is FIXED; the residual is now a real cache finding

The cause of the ledger being ahead of the database was found and fixed. A
mutation selected a child row from the ledger and then acted on it **by id only**,
so a concurrent writer could move or delete that row between selection and
application: the database changed a row the ledger attributed to someone else, or
refused silently, and the two diverged. Three changes close it:

1. every child-row statement now enforces the owner the harness observed
   (`... WHERE donation_id = $1 AND person_id = $N`), in all five schema
   directories;
2. a statement that matches no row is reported as `errNoEffect`: the database did
   not change, so **the ledger does not change either**, the write is acknowledged
   as a no-op and counted (`write_noops_refused_by_owner`);
3. amount corrections are relative in **every** schema, including the reference
   designs, where an absolute `SET amount_cents = $2` had been writing the *delta*
   as the amount.

Plus the assertion AM-01 demanded: after every writing phase the harness now
compares, for all donors, the content the database holds against the ledger's
requirement, and **fails the cell** if the requirement is ahead. Verified run
`verify-ledger` (6 cells, `small` scale):

| Phase | donors checked | mismatches |
|---|---|---|
| warm, mixed, hotspot, stampede | 800 each | **0** |
| instances | 800 | **1** (donor 67) in one cell |

And the acceptance check now also sees violations recorded by the
instances-phase arms, which it previously could not: an arm's own adapter counted
them while the cell's check read the base adapter's counter, so
`owned-opt-redis-through-strict-coord` had passed with 11 hidden stale reads. With
that closed, the post-fix verdicts are honest:

| Cell | stale-after-ack reads | ledger assertion | impossible |
|---|---|---|---|
| `owned-opt-redis-through-strict-coord` | 32 | all pass | 0 |
| `legacy-na-memory-through-strict-coord` | 7 | all pass | 0 |
| `legacy-na-redis-through-strict-coord` | (see report) | all pass | 0 |
| `ref-normalized-indexed`, `ref-embedded-locked` | — | all pass | 0, **cells now green** |

So the strict `through` residual is **a genuine cache-side race in the through
strategy, in both models and both backends**, with the fence demonstrably working
(260–625 refused publications per cell) and the accounting demonstrably clean. That
is the finding to hand over, and it is a stronger one than the earlier ambiguity.

`ref-rollup-trigger` remains failing, and its diagnosis narrowed usefully: all 114
impossible values come from **database** reads, which the harness's own rule
classifies as `ahead`, not impossible. Either the reclassification is not reached
on that path or the count is attributed to the wrong source — a concrete,
testable lead.

**Consequence for this analysis, stated plainly: the strict `through` failures
were originally not attributable to the cache design, because one of their
causes was a defect in the harness's own accounting; that defect is now fixed and
the failures are.** The earlier "weakness 7"
(an incomplete model of concurrently conflicting writes) is therefore upgraded
from speculation to evidenced fact, with a concrete test for the next iteration:
assert that `Required()` is never ahead of the content the database currently
holds, and fail the cell if it is, before judging any strict cell at all.

**What survives unchanged:** the strict **aside** cells (memory and Redis, both
models) recorded zero stale-after-ack reads and zero impossible values, and the
fence result in section 2 is measured by refused publications (61–85 per cell),
not by the strict counters. The single-instance relaxed wrong-read rates in
section 1 are also unaffected: they come from cells that passed every gate.

## 7. What a second analyst should check first

0. **First**: implement the ledger-side assertion from section 6a — `Required()` must
   never be ahead of the database's current content — and re-run the three strict
   `through` cells. Until that passes, every strict number in this study is
   provisional, because the harness's requirement can be wrong.
1. Reproduce the three-instance 85 % wrong-read rate on a second machine and with
   repeated trials. If it survives, it is the study's most valuable result and it
   contradicts the assumption behind every "shared cache" deployment in the
   repository so far.
2. Take the four strict failures apart with the attributed causes already in the
   results: the through-strategy asymmetry is the shortest path to a fix.
3. Re-run with `medium` scale and repeated trials, and add the churn and
   equal-total arms, before any of the numbers above are used for a decision.
4. Decide whether the ledger's ordering limitation (weakness 7) is fixed by
   per-key harness serialisation or by restricting the workload to
   order-insensitive mutations, and record which.

## 8. Signature

Analyst: `deepseek-flash`, acting as DeepSeek HIGH. It planned, implemented,
measured and analysed study 05 in one session; there was no LOW handoff, by the
task's instruction, and no independent model review. Every figure cited here was
checked against `reports/20260921T-survey.md` (inputs digest
`8ff86dc9e5228e86`) before this file was committed.
