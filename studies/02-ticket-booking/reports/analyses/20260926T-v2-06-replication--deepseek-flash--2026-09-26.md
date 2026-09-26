---
analysis_id: 20260926T-v2-06-replication--deepseek-flash--2026-09-26
run_id: 20260926T214452Z (designs in one order), 20260926T215213Z (the reverse order)
environment: host-zenbook-ux5406sa
analyst: deepseek-flash
analyst_kind: ai
analyst_version: "deepseek-flash (DeepSeek), HIGH-capability, Deep Code CLI; effort setting not exposed. Executor of task study-02/v2-06-replication."
analyzed_at: 2026-09-26
inputs_digest: 518f3bfc68e987df (run 20260926T214452Z), b1f88e0067f4db83 (run 20260926T215213Z, the reverse design order); the baseline for comparison is run/02-ticket-booking/20260913T021206Z
repo_commit: 9c608c89e8ff1e43a1b212916b2c1209468ec0c1
status: current
headline: The overbooking claim reproduces in direction and in its 4-8x range thirteen days later and under both design orders, but every absolute throughput is 26-72 percent below the original run and one design moves 58-76 percent with its position in the session.
---

# Analysis — study 02 `v2-06` replication: the range reproduces, the absolutes do not — deepseek-flash

## TL;DR

- **The headline reproduces, in both design orders.** `p3_precreated_cas` sold **4.30–7.06x** faster than the
  two counter designs in the first run and **4.64–6.99x** in the second, inside the 4–8x the claim states;
  `p2_precreated_skip_locked` beat both counters by **3.25–5.25x** and **5.78–8.71x**. T1 and T2 hold in
  both, from a session thirteen days after the original, with three fresh-load repeats per tier.
- **No absolute number reproduced.** Every cell is below its 2026-09-13 value, by **−26% to −72%**: `p3` at
  100 seats went 6,400 → 2,189 and `c3` 1,100 → 310. T3 is falsified. The claim's *range* survives; its
  figures are readings of their session.
- **One design moved with its position, and the counters did not.** Running the same five designs in the
  reverse order, `c3`, `c4` and `r1` stayed within 15% of their first-run values, `p3` within 7%, and
  `p2_precreated_skip_locked` moved **+58% and +76%** — it ran first in the first run and last in the
  second. So `p2` carries a position or warm-up effect of the same size as the whole comparison it is
  part of, and the claim's pairing of `p2` with `p3` hides that.
- **This is the study 04 lesson, on a different study.** Study 04 concluded that ratios are portable and
  absolute throughputs are not; here the same pattern appears in an independent study, on a different
  scenario, thirteen days apart.
- **Correctness is clean in every cell of both runs**: 0 errors, 0 overbooked seats, 0 underbooked, 0
  timed-out events across five designs, two tiers and two orders. T4 confirmed.
- **Coverage gap carried forward as promised**: the deliberately wrong controls were not re-run, so the
  registry claim gains a replication, not a re-fired negative control.

## The measurements

PostgreSQL 1 node, `small`, the sell-out race with 32 concurrent buyers, three fresh-load repeats per tier.
Medians of the three repeats; the first run's three readings are in the analysis of record below.

| Design | run 1, 100 / 1,000 seats | run 2 (reverse order), 100 / 1,000 | movement |
|---|---|---|---|
| `p3_precreated_cas` | 2,189 / 2,729 | 2,301 / 2,926 | +5% / +7% |
| `p2_precreated_skip_locked` | 1,627 / 2,100 | 2,866 / 3,318 | **+76% / +58%** |
| `c3_count_lock_event` | 310 / 637 | 329 / 539 | +6% / −15% |
| `c4_counter_guard` | 501 / 635 | 496 / 623 | −1% / −2% |
| `r1_inventory_row` | 468 / 559 | 477 / 617 | +2% / +10% |

Ratios against the two counters, of which the claim's range is made: run 1 `p3` 4.37x and 7.06x (100
seats), 4.30x and 4.28x (1,000); run 2 4.64x and 6.99x, 4.70x and 5.43x. `p2`: run 1 3.25x and 5.25x,
3.31x and 3.30x; run 2 5.78x and 8.71x, 5.32x and 6.16x.

## Prediction verdicts

| Id | Prediction | Verdict | Basis |
|---|---|---|---|
| T1 | `p3` is 4–8x faster than both counters at both tiers | **confirmed in both orders** | run 1: 4.37, 7.06, 4.30, 4.28; run 2: 4.64, 6.99, 4.70, 5.43 |
| T2 | `p2` beats both counters by ≥3x at both tiers | **confirmed in both orders** | run 1: 3.25, 5.25, 3.31, 3.30; run 2: 5.78, 8.71, 5.32, 6.16 |
| T3 | every cell within ±40% of the 2026-09-13 value | **falsified** | −26% to −72%, with only the 1,000-seat counters inside the band |
| T4 | the correctness gate stays green | **confirmed** | 0 errors, 0 overbooked, 0 underbooked, 0 timed-out in ten cells |

## What this changes

1. **`v2-06` becomes `run_replicated`** — three whole runs, the last two on the same day in opposite order
   (AB/BA), three repeats per tier — with the claim's figures re-scoped so the 4–8x range is the claim and
   the absolute numbers are labelled as one session's readings. The statement's range is unchanged because
   it reproduced under both orders.
2. **`p2`'s position effect is a new limitation on a design inside the claim.** The claim says the two
   pre-created designs are the fast side; it now needs to say that one of them moves by up to 76% with its
   place in the session while the other does not, and that the ≥3x result holds at both extremes.
3. **A different study reproduces the study 04 pattern.** "Carry the ratio, distrust the throughput" now has
   evidence from two studies, two workloads and thirteen days.
4. **The infrastructure finding that made this possible.** Study 02's runner passed `hostpath` where the
   image build needs `winpath`, which study 04 had already fixed in its own copy — so study 02 could not
   build at all from a WSL path, which is why none of its claims had ever been replicated. With that
   corrected each run took about five minutes.

## Where this analysis is weak

1. **One host, one engine, two of the original's five tiers.** 10 and 10,000 seats are not re-measured.
2. **The AB/BA pair is same-day**, so it separates design order from each other but not from the day; the
   comparison against 2026-09-13 mixes order, day and a change in repetition.
3. **The two runs are not like-for-like in repetition**: the original measured each cell once, the
   replications three times, so part of the observed difference is precision rather than session.
4. **The position effect is diagnosed, not explained.** `p2` moving with its slot is consistent with a
   warm-up, a container or a cache-state effect; nothing here distinguishes them.
5. **`r1` still shows one anomalous trial per tier** (in the first run, 122 and 117 against 468–653). It did
   not recur in the second run, which makes a harness artefact more likely than a design property, but it
   is unresolved and `r1` remains the weakest row.
6. **The controls did not fire**, so the correctness evidence for the wrong designs is unchanged from 2026-09-13.

## Next experiments

1. **Re-run the negative controls** in the same session as a replication, so the gate fires rather than passes.
2. **Repeat `p2` in several positions** within one session to size its position effect properly, since it is
   the design whose number is least portable.
3. **Sweep the missing tiers** (10 and 10,000 seats).
4. **Record host state beside every run** — the drift is now measured in four tranches and explained in none.
