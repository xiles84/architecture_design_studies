---
analysis_id: 20260929T-study02-client-quota-audit--deepseek-flash--2026-09-29
supersedes: 20260929T-v2-09-8buyer-recurrence--deepseek-flash--2026-09-29
run_id: 20260929T010633Z, 20260929T010807Z, 20260929T010947Z, 20260929T011128Z, 20260929T011301Z, 20260929T011444Z (the 8-buyer tranche); 20260920T234953Z, 20260913T021206Z (the audited anchors)
environment: host-zenbook-ux5406sa
analyst: deepseek-flash
analyst_kind: ai
analyst_version: "deepseek-flash (DeepSeek), leader-and-executor session, Deep Code CLI; audits the client-quota accounting the platform has recorded since the study's first commit."
analyzed_at: 2026-09-29
inputs_digest: 76f431bf0eb136fc, 2e5d35a847eda9ee, a86da7ea90c90c9a, 2eeac99e2287cfa9, 4a03a51500ee3a3a, 0dbf1c3abb037c71 (8-buyer tranche); 160bd80c49bbd892, 71bcee71d725d43d (audited anchors)
inputs: 20260929T010633Z@76f431bf0eb136fc, 20260929T010807Z@2e5d35a847eda9ee, 20260929T010947Z@a86da7ea90c90c9a, 20260929T011128Z@2eeac99e2287cfa9, 20260929T011301Z@4a03a51500ee3a3a, 20260929T011444Z@0dbf1c3abb037c71, 20260920T234953Z@160bd80c49bbd892, 20260913T021206Z@71bcee71d725d43d
repo_commit: 5d9c102
status: current
headline: The client's CPU accounting already existed per trial and nobody had read it - it settles the 8-buyer low mode (no client throttling, the client waiting) and exposes 35 client-throttled read-phase cells in two early anchors, where design order within the run decided who was throttled.
---

# Analysis — the client's own quota, read at last (study 02) — deepseek-flash

*Supersedes `20260929T-v2-09-8buyer-recurrence--deepseek-flash--2026-09-29`, whose conclusion that the
8-buyer mechanism was unassignable was based on a false premise: the client's counters were in the artifact
all along. That file is moved to `reports/outdated/`; its measurements stand, its inference does not.*

## TL;DR

- **The instrumentation was already there.** Every result file carries `client_cpu` — one entry per
  **trial** (`race#1`…`race#5`) and per phase — written by `platform/adapters/cgroup` since the study's
  first commit, and rendered by the report's "Measurement hazard: CPU-quota throttling" table. No analysis
  had read it. My previous analysis ("the artifact records the database's counters but nothing from the
  client") was **wrong**; the correction below is the reason it matters.
- **The 8-buyer low mode is resolved.** In all twelve 8-buyer cells the client's throttling is **zero**, and
  in the slow trials the client used **more** CPU (2.6–2.8 s) while its throughput fell 4x. The client was
  not starved; it was **waiting**. With the database's own throttling negligible (~4 ms), the cost is in the
  server's response path, not in either container's quota.
- **The audit's real find:** across all 775 recorded client-CPU cells in study 02, **35 are throttled** — and
  they are concentrated in the **read phase** of two anchors, `20260920T234953Z` and `20260913T021206Z`.
- **Within one run, throttling followed the design**: in `20260920T234953Z` the read phase lost **2.3 s** for
  `p2_precreated_skip_locked` and **58.1 s** for `r1_inventory_row` — a 25x spread between designs measured
  in the same session. A difference that large is not a property of the designs' read cost.
- **No claim is retracted**, because none of the three that anchor those runs rests on read throughput
  (`v2-06` is about sell rate, `v2-07` about storage and answerability, `v2-gap-01` about load). What is
  affected is every *read* figure printed in those two reports, which should not be quoted; the limits of all
  three now say so.
- **Why it matters methodologically**: on 2026-09-14 the owner made client sizing part of every study's
  protocol. The numbers to check it were being recorded for thirteen days and read by nobody, and the
  analysis layer drew a mechanism conclusion (mine) that the recording already refuted.

## The 8-buyer resolution (supersedes the previous conclusion)

Per-trial client CPU, 8 buyers, tier 1,000 (`usage` = client CPU consumed in the trial, `thr` = throttled):

| Run | Cell | Per-trial `sold/s` / client CPU |
|---|---|---|
| 010807Z (B) | `p3` CAS | 2690/1870ms · 2176/2040ms · 2690/1780ms · **655/2684ms** · 2546/1885ms |
| 011301Z (A) | `p3` CAS | 2406/1979ms · 2717/1728ms · 2635/1817ms · **623/2782ms** · 2642/1750ms |
| 010947Z (B) | `x1` ledger | 2641/1855ms · 2421/1877ms · 2336/2024ms · 3095/1711ms · **619/2572ms** |
| 20260928T014723Z | `x1` ledger | 4840/879ms · 4703/907ms · **967/1080ms** · **972/1088ms** · **1010/1062ms** |

`nr_throttled` is **0 in every one of the sixty trials**. The slow trials are not starved: they consume
*more* client CPU than the fast ones (2.6–2.8 s against ~1.8–1.9 s), and in the original collapsed cell the
client went from 84% busy (0.88 s CPU in 1.05 s wall) to 20% busy (1.08 s CPU in 5.2 s wall). So the client
is idle-waiting on responses while working harder per completed sale, and the database container's own
`cpu.stat` for that cell shows ~4 ms of throttling in total. **The 8-buyer low mode is a latency excursion
in the round trip, not a quota effect on either side** — which is a different conclusion from "the cause is
unassignable", and it is one this data supported all along.

## The audit

Scanning every `client_cpu` entry in `studies/02-ticket-booking/results/**` (775 cells: per-trial for the
race phase, per-phase elsewhere), 35 cells record `nr_throttled > 0`. All 35 are **read-phase** cells (plus
one `reports` cell each in `h0`, `h1`, `r3`, `x1` of the 2026-09-20 run). They are confined to two runs.

**`20260920T234953Z`** — read-phase client throttling per design, out of ~1,085 CFS periods:

| Design | Throttled periods | Throttled time |
|---|---:|---:|
| `p2_precreated_skip_locked` | 73/1085 | 2.3 s |
| `p1_precreated_lock_first` | 121/1083 | 5.1 s |
| `p3_precreated_cas` | 120/1084 | 4.8 s |
| `c1_count_naive` | 100/1084 | 4.6 s |
| `c5_seat_unique` | 132/1084 | 8.3 s |
| `c2_count_serializable` | 141/1083 | 8.5 s |
| `c3_count_lock_event` | 127/1083 | 8.9 s |
| `x1_cas_ledger` | — | 3.0 s |
| `r2_inventory_buckets` | 637/1084 | 36.5 s |
| `c4_counter_guard` | 714/1086 | 53.9 s |
| `p4_precreated_counter` | 712/1086 | 56.3 s |
| `r1_inventory_row` | 716/1085 | 58.1 s |

**`20260913T021206Z`** is the same shape, milder: `p2` 0.2 s, `p1` 1.0 s, `p3` 1.5 s, `c1` 1.4 s, `c3`
2.6 s, `c5` 3.0 s, `c2` 3.4 s, `r3` 4.6 s, `h1` 22.2 s, `r1` 28.0 s, `p4` 28.1 s, `h0` 29.5 s.

Two things follow. First, the **absolute** read throughput of those runs is a floor set by the client's
quota, not a measurement of the database. Second, and worse for comparison, the throttling is
**design-dependent within a single run** — 2.3 s versus 58.1 s in one session — so a read-phase *ranking*
between designs from those runs can be an ordering of who was throttled. The four most-throttled designs
(`p4`, `c4`, `r1`, `r2`) are precisely the counter/inventory family the read comparisons were run to test.

**What this does not do.** None of the three claims anchored to those runs rests on read throughput:
`v2-06-hot-drop-designs` (primary anchor `20260913T021206Z`) is about hot-drop **sell** rate, and its race
phase shows one throttled cell in fourteen designs; `v2-07-refund-answerability-and-storage` and
`v2-gap-01-ledger-load-multiplier-unstable` (both primary-anchored on `20260920T234953Z`) rest on storage,
answerability and load. Their conclusions stand. Their reports' read tables do not, and a reader who quotes
them would be quoting the client's quota. All three claims now carry a limit recording this.

## Where I think the measurement is weak

1. **Two runs.** 34 of study 02's run directories record no client throttling outside `read`, and recent
   runs (including every 2026-09-28/29 tranche) show none at all; but the audit is a scan of what was
   recorded, and cells written before the cgroup adapter existed have no `client_cpu` key to read.
2. **The 8-buyer mechanism is narrowed, not identified.** "Waiting on the response path" excludes both
   containers' quotas; it does not distinguish a slow query plan, a lock wait, or a network/VM scheduling
   stall between the client and server containers.
3. **Attribution of the old throttling to phase design is not possible from the artifact**: the read phase's
   query mix per design is not in the result JSON, so why `r1`/`p4`/`c4`/`r2` cost the client more CPU than
   `p2` is not recoverable here.
4. **One host, one engine per cell, the 2 GiB machine ceiling** — and the observation that the database
   container was itself heavily throttled in those 2026-09-13/20 runs is *not* something this analysis can
   quantify, because that era's `dbcpu.txt` figures are not consistent with the period counts recorded.

## What I would measure next

1. **Re-run the affected read comparisons** on the current harness, which is the only way to know whether
   the read family's ordering was decided by the quota; until then the ordering is unquoted, not refuted.
2. **Make the hazard a gate, not a table**: a protocol could pre-register that a cell whose client throttled
   more than a stated fraction of its periods is void, decided before the run.
3. **Reduce the client's own footprint** for read-heavy phases (fewer, longer-lived connections), since the
   client now has a recorded history of being the bottleneck in exactly those phases.

## Disagreements with other analyses

It **corrects my own** `20260929T-v2-09-8buyer-recurrence` (superseded above) and, more quietly, every
analysis of `20260913T021206Z` and `20260920T234953Z`, which read those runs' read tables without the
throttling beside them. It agrees with the corpus's standing "ratio is more portable than absolute
throughput" lesson and gives it a second mechanism: when the client is quota-limited, the absolute moves and
the ratio survives only if the quota bit both designs equally, which in those two runs it did not.
