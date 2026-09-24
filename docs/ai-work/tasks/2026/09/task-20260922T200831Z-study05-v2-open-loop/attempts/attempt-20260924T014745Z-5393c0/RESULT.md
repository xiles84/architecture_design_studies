# Result — Study 05 v2 open-loop demand phase

**Task:** `task-20260922T200831Z-study05-v2-open-loop`
**Branch:** `study-05/v2-open-loop`
**Required tag:** `study-05/v2-open-loop`
**Capability:** HIGH session executing the task / executor (model `deepseek-flash`)
**Benchmark:** taken for each run; one cell at a time; nothing measured in parallel

## What was delivered

1. **An open-loop demand phase** in the Study 05 harness, driven by the platform's
   `measure.RunOpenLoop` (no closed-loop backpressure). It records, per offered rate: offered, started,
   completed, rejected, errors, dropped (never attempted), delivered and offered rates, latency p50/p99/
   max, scheduling-lag p99, the client-saturation verdict and reason, and the phase's own wrong-read
   account. Flags: `-openloop-rates`, `-openloop-duration`, `-openloop-queue`, passed through by the
   runner and recorded in the manifest. The generated report has an open-loop section that states
   closed-loop results are floors and are never reused as open-loop.
2. **Three runs** (the first two superseded by rate choice, their reports in `reports/outdated/`,
   their results retained):
   - `20260923T0220Z-openloop` — 4,000 / 16,000 ops/s
   - `20260923T0240Z-openloop` — 6,000 / 12,000 ops/s
   - **`20260923T0255Z-openloop` — 5,000 / 30,000 ops/s (canonical)**
3. **A generated report** `reports/20260923T0255Z-openloop.md` (inputs digest `50d5a426c89e96cf`)
   and a **signed analysis**
   `reports/analyses/20260923T0255Z-openloop--deepseek-flash--2026-09-23.md`.
4. `CONTEXT.md` updated (open-loop landed; the remaining-gap sentence corrected).

## Canonical result (run `20260923T0255Z-openloop`)

```text
design/topology    owned-opt-redis-aside-relaxed-coord / pg-single, db-only, small, 1 trial
closed-loop floor  21,175 ops/s (warm_read_only, same cell)

offered 5,000/s    delivered 4,969/s (99.4%) · 50,000/50,000 completed · 0 dropped
                   p50 0.14 ms · p99 0.58 ms · max 221 ms · sched lag p99 1.12 ms
offered 30,000/s   delivered 15,599/s (52.0%) · 158,294 started of 300,004 offered
                   141,710 (47.2%) never attempted · p50 0.14 ms · p99 0.74 ms · max 224 ms
                   sched lag p99 48.8 ms · client saturated (buffer full)
correctness        gate 14/14 · 0 wrong reads of 50,000 and 0 of 158,294 · 0 impossible values
```

**Answer to "did the server saturate?"** No server saturation was isolated: the delivered rate at
30,000/s offered (15,599/s) is *below* the same cell's closed-loop floor (21,175/s), and the driver
flags the client as saturated, so the shortfall is a lower bound caused by the client's bounded
arrival buffer on the shared laptop. The offered 30,000/s is not sustainable end-to-end.

## Acceptance checks

```text
open-loop arrival at two fixed rates, no closed-loop backpressure   yes (measure.RunOpenLoop)
achieved throughput per offered rate                                yes
non-attempted (dropped) operations per offered rate                 yes
latency percentiles per offered rate                                yes
offered rate stated; server saturation stated if reached            yes; not isolated (client saturated first)
closed-loop numbers not reused as open-loop                         yes (stated in report + analysis)
correctness gates timing; wrong-read rate with the throughput       yes (same cell, same phase)
generated report + signed analysis                                  yes
benchmark lock taken, one matrix at a time                          yes
```

## Limitations (also in the analysis)

- One trial per rate; the closed-loop floor varied 7,787 → 19,921 → 21,175 ops/s across the three runs
  of this same design, so small differences are not interpretable.
- The client saturated in both regimes; no server-side saturation point was measured.
- Read-only overload: relaxed staleness is not exercised (no writers), so 0 wrong reads is not a
  freshness result.
- One design, one topology, one engine, shared laptop cores, `db-only` framing.
