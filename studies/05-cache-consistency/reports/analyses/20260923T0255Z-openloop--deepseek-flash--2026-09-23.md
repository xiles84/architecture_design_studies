---
analysis_id: 20260923T0255Z-openloop--deepseek-flash--2026-09-23
run_id: 20260923T0255Z-openloop
environment: host-zenbook-ux5406sa
analyst: deepseek-flash
analyst_kind: ai
analyst_version: "deepseek-flash (DeepSeek), HIGH-capability, Deep Code CLI; effort setting not exposed. Author of this run and of the open-loop phase it uses."
analyzed_at: 2026-09-23
inputs_digest: 50d5a426c89e96cf
repo_commit: ee92bc154aeccb54f5c351aeb78da32f8702a251
supersedes:
status: current
headline: Open-loop demand is now measured: offered 5,000/s was delivered at 4,969/s with nothing dropped, while offered 30,000/s was NOT sustained — 47.2% of arrivals were never attempted and the client's bounded buffer, not the server, was the first limit.
---

# Analysis — Study 05 v2 open-loop demand — deepseek-flash

> Deliverable of `task-20260922T200831Z-study05-v2-open-loop`, implementing section 6 of
> `studies/05-cache-consistency/HANDOFF-AMENDMENT-01-V2.md` (tag `study-05/v4-handoff`).
> One cell: `owned-opt-redis-aside-relaxed-coord` on `pg-single`, `resource_framing: db-only`,
> `small` scale, one trial. Digest `50d5a426c89e96cf`, commit `ee92bc1`.

## TL;DR

- **Below capacity, the offered rate is delivered.** Offered **5,000 ops/s**, the cell completed
  **50,000 / 50,000** reads at **4,969 ops/s** (99.4%), `dropped = 0`. p50 0.14 ms, p99 0.58 ms.
- **Above capacity, the offered stream is not sustained.** Offered **30,000 ops/s**, the client
  handed over **158,294 / 300,004** arrivals and **141,710 (47.2%)** were never attempted; the
  delivered rate was **15,599 ops/s**. p50 0.14 ms, p99 0.74 ms, but the scheduling lag p99 was
  **48.8 ms** — the generator was far behind its own schedule.
- **The server did not saturate first.** The same cell's closed-loop floor is **21,175 ops/s**
  (`warm_read_only`), above every delivered open-loop rate here, so the first limit observed was the
  client's bounded arrival buffer. Because the client is flagged saturated, **15,599/s is a
  client-limited lower bound, not a measured server saturation point**.
- **Correctness held in both regimes:** 0 wrong reads of 50,000 and 0 of 158,294, 0 impossible values,
  gate 14/14. The design is relaxed, but this phase has no writers, so this is not a freshness result.
- **Closed-loop numbers are floors and are not reused as open-loop.** The closed-loop floor (21.2k)
  is *higher* than the open-loop delivered peak (15.6k): on this shared laptop the open-loop generator
  itself becomes the bottleneck before the server does.

## What was run

The harness gains an `openloop` phase driven by the platform's `measure.RunOpenLoop`: arrivals are
scheduled on a clock and the offered stream does not slow down when the server does. One cell:

| | |
|---|---|
| Scenario | `owned-opt-redis-aside-relaxed-coord` |
| Topology / framing | `pg-single` / `db-only` |
| Scale / trials | `small` / 1 |
| Client workers | 8 (scale default) |
| Arrival buffer | 32 (4 × workers) |
| Offered window | 10 s per rate |
| Rates | 5,000 and 30,000 ops/s |

## The measurements

| Offered ops/s | Delivered ops/s | Delivered % | Offered | Started/completed | Not attempted | Rejected/errors | p50 ms | p99 ms | max ms | sched lag p99 ms | wrong reads |
|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| 5,000 | 4,969 | 99.4% | 50,000 | 50,000 | 0 | 0 / 0 | 0.14 | 0.58 | 221 | 1.12 | 0 of 50,000 |
| 30,000 | 15,599 | 52.0% | 300,004 | 158,294 | 141,710 (47.2%) | 0 / 0 | 0.14 | 0.74 | 224 | 48.8 | 0 of 158,294 |

Closed-loop reference from the same cell: `warm_read_only` **21,175 ops/s** (2 s, 1 trial).
Cache state: 820 fills, capacity 256 KiB, 0 evictions; 11,554 fill-lease acquisitions.

The low rate is served: delivered ≈ offered, nothing dropped, and the tail is short (p99 0.58 ms).
The high rate is not: nearly half the offered operations were never attempted, and the ones that were
served took the same p50/p99 as the low rate — the shortfall is entirely the arrivals the generator
could not hand to a worker, exactly the coordinated-omission gap the open-loop driver exists to expose.

## Whether the server saturated

No. The amendment asks for the offered rate and the server-side saturation point **if reached**. At
5,000/s the server was nowhere near its floor. At 30,000/s the delivered rate stopped at 15,599/s,
which is *below* the closed-loop floor of 21,175/s for the same cell, and the driver's own rules mark
the client as saturated (`arrivals dropped: the buffer filled before a worker could take them`). A
saturated client invalidates a claim about the server, so the honest statement is: **the offered
30,000/s is not sustainable end to end, at least 47.2% of it was never attempted, and the server's
own saturation point is not isolated by this run.**

## Where the measurement is weak

1. **One trial per rate, no error bar.** The closed-loop floor for the *same design and topology*
   varied across three runs of this task — 7,787, 19,921 and 21,175 ops/s — so no conclusion may rest
   on a difference smaller than that spread, and the three open-loop runs are not comparable to each
   other rate-for-rate.
2. **The client saturated in both regimes**, by the lag rule at 5,000/s (p99 lag 1.12 ms > one
   200 µs slot) and by drops at 30,000/s. The delivered 15,599/s is therefore a lower bound; a run
   with more workers or a deeper buffer might deliver more and find the server's knee, at the cost of
   the buffer absorbing the very saturation the phase is meant to show.
3. **Read-only overload.** There are no writers in this phase, so the relaxed contract is never
   exercised: 0 wrong reads here says nothing about staleness under concurrent writes.
4. **One design, one topology, one engine.** Only the owned/Redis/cache-aside/relaxed cell on
   PostgreSQL single-node was offered; strict cells, the memory backend and the cluster topologies
   were not.
5. **Shared laptop cores, no network, cache budget outside the comparison** (`db-only`). The client,
   the Redis container and PostgreSQL share the same eight throttled cores, which is precisely why
   the generator became the limit.

## Two superseded runs

Two earlier open-loop runs of the same task used different rate pairs (4,000/16,000 and 6,000/12,000)
and were superseded by the rate choice, not by a defect; their reports were moved to
`reports/outdated/` and their results remain under `results/`. They are the source of the closed-loop
floor spread quoted above.
