---
analysis_id: 20260926T-interaction-repeat--deepseek-flash--2026-09-26
run_id: 20260926T152426Z, 20260926T171715Z, 20260926T171907Z, 20260926T152621Z, 20260926T185214Z, 20260926T185408Z, 20260926T152812Z, 20260926T185614Z, 20260926T185816Z
environment: host-zenbook-ux5406sa
analyst: deepseek-flash
analyst_kind: ai
analyst_version: "deepseek-flash (DeepSeek), HIGH-capability, Deep Code CLI; effort setting not exposed. Executor of task study-04/interaction-repeat."
analyzed_at: 2026-09-26
inputs_digest: per-run 69ef9738af6d2fd4, 8163a927b72f07a7, 1547ba73ec39ee7a (tier 10); 5618b5ce5954adcf, 692ecfe37db5ca14, 7b0920167fe923ee (tier 60); 62417ed21d652380, 027f25ea3fd78f26, dd3008b72b5c6db9 (tier 240)
repo_commit: e6b0f42
status: current
headline: The interaction's trend reproduces in all three repetitions, but its magnitudes do not - the ratio spreads up to 55% within a tier, so v2-18's numbers are bands (about 2.1-2.4x, 2.5-4.4x, 13-24x) rather than the points it first reported.
---

# Analysis — Study 04 interaction repeat: the trend survives, the magnitudes do not — deepseek-flash

## TL;DR

- **The trend reproduces in every repetition.** Reading the three runs in any grouping, the ratio rises
  with replacement size: 2.07/2.42/2.32 at tier 10, 4.38/3.42/2.51 at tier 60, 24.20/20.88/13.36 at
  tier 240. The ordering by tier holds in all three repetitions, so R1 is confirmed.
- **The magnitudes do not, and R2 is falsified.** The three runs at a tier agreed within **15%** at tier
  10 but only within **55%** at tier 60 and **52%** at tier 240, against a falsifier of 40%. The level
  cannot carry a point, so `v2-18` is corrected from 2.07x / 4.38x / 24.20x to **bands**: roughly
  2.1–2.4x, 2.5–4.4x and 13–24x.
- **The scatter is not noise, it is drift.** Absolute throughput fell steadily as the session went on for
  both designs — n4 at tier 60 read 631 → 205 → 206 ops/s and n3 144 → 60 → 82 — and the *ratio* moved
  because the two designs drifted by different amounts. This is the third tranche in which a session-level
  term dominates the differences being measured.
- **R3 is confirmed, narrowly.** The tier-60 median is 3.42x, 22% below the 4.38x of the first run and
  inside the 30% band the prediction allowed, so the earlier corroboration of `v2-13` survives as
  approximate rather than exact.
- **R4 is confirmed.** Every gate passed in all nine runs, with 0 audit failures and 0 rollup mismatches
  at 10, 60 and 240 entries per installation.

## The nine runs

| Tier | Ratios (n4/n3) | Median | Spread | n3 ops/s | n4 ops/s |
|---|---|---|---|---|---|
| 10 | 2.07, 2.42, 2.32 | 2.32 | 15.1% | 459 → 397 → 238 | 952 → 960 → 552 |
| 60 | 4.38, 3.42, 2.51 | 3.42 | 54.7% | 144 → 60 → 82 | 631 → 205 → 206 |
| 240 | 24.20, 20.88, 13.36 | 20.88 | 51.9% | 10 → 8 → 11 | 242 → 167 → 147 |

`pg-single`, `small`, `verify,write`, three measurement repeats per configuration, operation
`w05_replace_all`. Within each tier the runs are listed in the order they were run.

## Prediction verdicts

| Id | Prediction | Verdict | Basis |
|---|---|---|---|
| R1 | the ratio rises with the tier in every repetition | **confirmed** | the ordering holds in all three groupings |
| R2 | the three runs at a tier agree within 40% of their median | **falsified** | 15.1% at tier 10, but 54.7% at 60 and 51.9% at 240 |
| R3 | the tier-60 median is within 30% of 4.38x | **confirmed** | median 3.42x, 22% below |
| R4 | correctness holds in every run | **confirmed** | 0 audit failures, 0 rollup mismatches in all nine |

## What this changes in `v2-18`

The interaction is real and reproducible in direction. What the first tranche could not see is that its
three numbers were three readings of one session. The claim now carries bands instead of points, a
`run_replicated` strength with a three-run record, and limits saying which part of the spread is
repetition and which part is the session.

**It also weakens one sentence I wrote in the first analysis**: the tier-60 cell was described as landing
"inside the replicated range of 3.99–5.69x" of `v2-13` and therefore corroborating it. With three runs the
tier-60 median is 3.42x, below that range. The corroboration is now approximate — same operation, same
size, same direction — not a match.

## Where this analysis is weak

1. **Repetition index and session time are confounded.** The six added runs went in wall-clock sequence
   (two per tier, tier by tier), and throughput fell as the session progressed. So the measured spread
   contains drift-by-time as well as run-to-run scatter, and this tranche cannot separate them. That is a
   defect in the design of the repeat, not in the data: randomizing the order of the *runs* across tiers,
   or interleaving them, would have separated the two.
2. **Three runs per level is still one machine and one day.** Nothing here is hardware-independent, and
   the drift measured is within-session at that.
3. **The design order varied but was not controlled.** Six of the nine runs realized one order and three
   realized the other; that variation is recorded, but it was not crossed with the tiers deliberately.
4. **Three points on the axis**, so no crossing can be found and none is claimed.
5. **The bands are wide enough to weaken, not to overturn.** At tier 60 the range 2.5–4.4x still excludes
   the tier-10 band and the tier-240 band, so the interaction stands; but a reader who wants a number
   should treat the median as the number and the range as the uncertainty.

## Next experiments

1. **Interleave the repeats** — run the three tiers in an order that alternates, or repeat the whole
   sweep several times in one session, so drift-by-time and repetition can be separated.
2. **Re-measure a fixed reference cell before and after each tranche**, so the session's drift is
   recorded as a number rather than inferred from the curve.
3. **Push the axis past the small scale** (480, 960 entries) to see whether the ratio keeps rising and to
   test whether the bands stay separated.
