---
analysis_id: 20260926T-interleaved-sweep--deepseek-flash--2026-09-26
run_id: 20260926T193412Z, 20260926T193504Z, 20260926T193657Z, 20260926T193850Z, 20260926T194051Z, 20260926T194253Z, 20260926T194500Z, 20260926T194704Z, 20260926T194909Z, 20260926T195111Z, 20260926T195315Z
environment: host-zenbook-ux5406sa
analyst: deepseek-flash
analyst_kind: ai
analyst_version: "deepseek-flash (DeepSeek), HIGH-capability, Deep Code CLI; effort setting not exposed. Executor of task study-04/interleaved-sweep."
analyzed_at: 2026-09-26
inputs_digest: per-run f3daa3168f6fa4a1, 5156884200884bd8, c3b5ea462bdcd1e5, 95bbd5df836d276b, 4ca2c77002b2f1e2, 4122fe5de9bb5bf6, 3870fd766da741f2, 352076a00d795148, dd9b98282c0e5ab5 (interaction runs); reference runs 20260926T193412Z and 20260926T195315Z
repo_commit: 0f4a76f386fe6f2bbee7edd3a7e99e23077f5b99
status: current
headline: Interleaving the sweeps confirms the interaction's ordering in all three sweeps but roughly doubles the measured spread - so the sequential repeat's narrow spread was temporal proximity, not stability, and the levels' uncertainty is about a factor of two, not 15-55%.
---

# Analysis — Study 04 interleaved sweep: proximity was hiding the drift — deepseek-flash

## TL;DR

- **The ordering by size holds in every sweep.** Sweep A (10, 60, 240) gave 1.80 / 4.78 / 22.00; sweep B
  (60, 240, 10) gave 8.95 / 11.33 / 3.15; sweep C (240, 10, 60) gave 10.20 / 1.34 / 4.28. Read by size,
  every sweep rises. L1 is confirmed, and the interaction does not depend on the sequential design.
- **The apparent spread roughly doubled, and that is the finding.** Within a size the interleaved runs
  spread **101%, 98% and 104%** at tiers 10, 60 and 240 — against 15%, 55% and 52% in the sequential
  tranche. L3 is falsified, in the opposite direction to the prediction.
- **Why: repetition close in time is not repetition.** In the sequential tranche each size's three runs
  were adjacent in wall-clock order, so they sampled almost the same session state and agreed with each
  other. Interleaving spreads them across the session, so each size samples the whole drift range. The
  sequential design therefore **understated** the uncertainty rather than measuring it.
- **The drift is now a number, not an inference.** The reference cell — `c2_pessimistic_lock` in the
  contention phase at 16 writers, measured immediately before and after the sweeps — read **445 ops/s at
  the start and 640 at the end, +44%**. L2 is falsified: the session was not stable, and the instrument
  caught it. (It rose this time; in the previous session the same kind of drift ran downwards.)
- **The bands are wide but the ordering is not fragile.** Pooling all six runs per size, the ratio ranges
  are 1.34–3.15x at 10 entries, 2.51–8.95x at 60 and 10.20–24.20x at 240. Within the interleaved tranche
  the three ranges are **disjoint** (3.15 < 4.28 < 10.20); pooled with the sequential runs, the two lower
  ranges touch (3.15 against 2.51).
- **L4 confirmed**: 0 audit failures and 0 rollup mismatches in all nine interaction runs.

## The tranche

| Position | Run | Size | n3 ops/s | n4 ops/s | n4/n3 |
|---|---|---|---|---|---|
| reference (before) | 20260926T193412Z | — | — | — | c2 = 445 ops/s |
| sweep A | 20260926T193504Z | 10 | 292 | 527 | 1.80 |
| sweep A | 20260926T193657Z | 60 | 65 | 311 | 4.78 |
| sweep A | 20260926T193850Z | 240 | 7 | 154 | 22.00 |
| sweep B | 20260926T194051Z | 60 | 43 | 385 | 8.95 |
| sweep B | 20260926T194253Z | 240 | 9 | 102 | 11.33 |
| sweep B | 20260926T194500Z | 10 | 181 | 570 | 3.15 |
| sweep C | 20260926T194704Z | 240 | 10 | 102 | 10.20 |
| sweep C | 20260926T194909Z | 10 | 175 | 234 | 1.34 |
| sweep C | 20260926T195111Z | 60 | 43 | 184 | 4.28 |
| reference (after) | 20260926T195315Z | — | — | — | c2 = 640 ops/s |

## Prediction verdicts

| Id | Prediction | Verdict | Basis |
|---|---|---|---|
| L1 | the ratio rises with the size within every sweep | **confirmed** | holds in A, B and C |
| L2 | the reference cell differs by <20% across the session | **falsified** | 445 → 640 ops/s, +44% |
| L3 | interleaving reduces the spread below the sequential tranche's | **falsified, opposite sign** | 101% / 98% / 104% against 15% / 55% / 52% |
| L4 | correctness holds in every run | **confirmed** | 0 audit failures, 0 rollup mismatches in nine runs |

## What this changes

1. **`v2-18`'s bands widen to what the data actually spans** — about 1.3–3.2x, 2.5–9.0x and 10–24x over six
   runs per size — and the claim records that the two lower bands touch when the sequential and
   interleaved runs are pooled.
2. **The previous lesson needs correcting.** It said a sequential repeat "measures the session too", which
   was right, and implied that interleaving would give a truer picture, which was also right — but it
   assumed the sequential spread was an *upper* bound contaminated by drift. The opposite is true: the
   sequential spread was a **lower** bound, because adjacent runs share session state. A repeat design is
   not merely a way to add runs; its *arrangement in time* decides what it can see.
3. **The reference cell works and should be standard.** Two extra runs costing about two minutes turned
   "the session drifted, probably" into "+44% across this tranche". It is the cheapest instrument in this
   corpus and it is now used by a claim rather than recommended in prose.

## Where this analysis is weak

1. **The reference cell is a proxy, not the same measurement.** It is the contention phase on a different
   design pair, so its +44% bounds the session's movement without being the drift of `w05_replace_all`.
   A reference cell drawn from the *same* operation would be stronger.
2. **Six runs per size is still small**, and the ranges are min–max of six points, not an interval with a
   stated coverage.
3. **Drift is measured, not explained.** Nothing here says whether it is thermal, cache state, container
   reuse or something else; the reference run is a symptom recorder.
4. **One host, one day, three sizes, one scale.** The ordering's robustness is established here and
   nowhere else.
5. **The Latin square removes the size/position confound but not the drift**, which it distributes evenly
   rather than removing — so the spread it reports is a fair summary of this session, not a floor.

## Next experiments

1. **Take the reference cell from the same operation** (`w05_replace_all` at one size) so the drift is
   measured on the quantity being compared.
2. **Run the reference cell between every run, not only at the ends**, which turns the drift into a
   per-run covariate that the ratios can be read against.
3. **Repeat the whole design on another day** to see whether the disjointness of the three ranges holds
   when the session's drift runs in the other direction, as it did in the previous tranche.
