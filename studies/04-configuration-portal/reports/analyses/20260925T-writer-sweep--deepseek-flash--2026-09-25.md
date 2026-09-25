---
analysis_id: 20260925T-writer-sweep--deepseek-flash--2026-09-25
run_id: 20260923T0011Z-ctlC-w1, 20260925T224234Z, 20260923T0014Z-ctlC-w4, 20260925T224309Z, 20260923T00116Z-ctlC-w16, 20260925T224342Z, 20260923T00164Z-ctlC-w64
environment: host-zenbook-ux5406sa
analyst: deepseek-flash
analyst_kind: ai
analyst_version: "deepseek-flash (DeepSeek), HIGH-capability, Deep Code CLI; effort setting not exposed. Executor of task study-04/t1-followup."
analyzed_at: 2026-09-25
inputs_digest: per-run d897600aa583d03a (w1), b12d44532119fe3a (w2), 1e5cdb992e027ec8 (w4), 457a072982a9f7ed (w8), 7f3b085773289fb3 (w16), 88b7f53af11bb58b (w32), 624799b41ae48039 (w64)
repo_commit: ecd9638
status: current
headline: Optimistic version-checking and pessimistic locking cross over between one and two concurrent writers on one hot key — the lock leads by 1.68x at 32 writers and loses at one — giving the corpus its first bounded crossover.
---

# Analysis — Study 04 writer sweep: the corpus's first bounded crossover — deepseek-flash

## TL;DR

- **There is a crossing, and it is bracketed by two adjacent measured levels.** A lock-before-update
  (`c2_pessimistic_lock`) delivers **0.89x** the version check's throughput at **1 writer** and **1.10x** at
  **2 writers**. The optimistic version check wins below the crossing; the lock wins above it.
- **The sweep now has seven levels** (1, 2, 4, 8, 16, 32, 64 writers), so the recurrence is bounded by the
  interval **[1, 2] writers** rather than merely observed. Three levels (2, 8, 32) were added for this
  analysis; the earlier four (1, 4, 16, 64) are re-used unchanged.
- **The two sides behave differently, which is the mechanism.** The version check's throughput falls
  **2.9x** from 1 to 64 writers (1,329 → 466 ops/s) while the lock's falls far less (1,177 → 745). The
  version check's conflict count rises from **0** at one writer to **21,301** at 64.
- **The effect saturates rather than reversing again.** The multiplier peaks at 1.68x (32 writers) and is
  1.60x at 64, so the corpus has one crossing in this dimension, not two.
- **Boundary:** one measurement per level, one host, one engine, one key, single node. This is a curve,
  not a replicated one — every level is a single run.

## The curve

| Writers | c1 optimistic ops/s | c2 pessimistic ops/s | c2/c1 | c1 conflicts |
|---|---|---|---|---|
| 1 | 1,329 | 1,177 | **0.89** | 0 |
| **2** | **1,167** | **1,285** | **1.10** | 8,163 |
| 4 | 899 | 1,133 | 1.26 | 18,236 |
| 8 | 821 | 1,227 | 1.49 | 24,807 |
| 16 | 552 | 870 | 1.58 | 20,218 |
| 32 | 523 | 877 | 1.68 | 21,504 |
| 64 | 466 | 745 | 1.60 | 21,301 |

`pg-single`, `small`, designs `c1_optimistic_version,c2_pessimistic_lock`, `contention` phase, three
trials per level (median reported). Every cell passed its correctness gate. The three new levels ran
`verify,contention`; the earlier four ran `verify,read,write,contention`.

## Why this is a crossover and not a two-point comparison

The corpus's charter prefers a boundary to a bigger number, and the reason is visible here: at one writer
the two designs differ by 11% in the *opposite* direction from the difference at 32 writers (68%). A
report that measured only 16 writers would have concluded "the lock is 1.58x faster" and been wrong at
low contention — where a portal's real traffic mostly lives. The crossing interval is what makes the
recommendation conditional instead of unconditional.

## Relation to the registered claim

`v2-12-lock-vs-version-upper-bound` states the 16-writer result as an upper bound, because that run could
not separate the lock's advantage from per-retry harness work. This sweep neither strengthens nor weakens
that caveat: the same confound is present at every level, since the version check pays for retries that
the lock does not. The crossing is therefore a crossing in *this harness's* cost model, and the
mechanism — conflict handling being re-done per retry — is what makes it plausible beyond it.

## Where this analysis is weak

1. **One measurement per level, no repetition.** Seven levels is a curve in the workload dimension, not a
   replicated result; no level has a second run, so no level has a spread.
2. **Retries are not separated from the lock's advantage** (the v2-12 caveat, unchanged).
3. **Two different phase sets** across the seven runs. The contention metric should not depend on whether
   read and write phases also ran, but that is an assumption, not a measurement.
4. **No negative control in these runs.** The lost-update control (93.7% lost) lives in the v2 run, not
   here, so `control_coverage` is recorded as absent rather than implied.
5. **One key, one engine, single node, small scale.** The absolute numbers are a laptop comparison; the
   *shape* is the finding.
6. **The 1-writer point is at the edge of the sweep**, so the optimistic side is bracketed by only one
   level below the crossing. A 0.5-writer arm does not exist here.

## Next experiments

1. Repeat at least the three levels around the crossing (1, 2, 4) three times each, to put a spread on
   the interval and test whether the crossing stays between 1 and 2.
2. Disable retries at a level and re-measure, to separate the retry cost from the guard cost — the
   experiment `v2-12` has been waiting for.
3. Sweep a second dimension (keys, not writers) to see whether the crossing moves with contention per
   key or with client count.
