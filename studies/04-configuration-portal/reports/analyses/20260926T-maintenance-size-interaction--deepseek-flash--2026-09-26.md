---
analysis_id: 20260926T-maintenance-size-interaction--deepseek-flash--2026-09-26
run_id: 20260926T152426Z, 20260926T152621Z, 20260926T152812Z
environment: host-zenbook-ux5406sa
analyst: deepseek-flash
analyst_kind: ai
analyst_version: "deepseek-flash (DeepSeek), HIGH-capability, Deep Code CLI; effort setting not exposed. Executor of task study-04/maintenance-contention."
analyzed_at: 2026-09-26
inputs_digest: per-run 69ef9738af6d2fd4 (tier 10), 5618b5ce5954adcf (tier 60), 62417ed21d652380 (tier 240)
repo_commit: 1e2d046
status: current
headline: The rollup-maintenance penalty is not a constant but a function of replacement size - 2.07x at 10 entries per installation, 4.38x at 60 and 24.20x at 240 - giving the corpus its first interaction claim and converting v2-13's stated gap into a measurement.
---

# Analysis — Study 04 interaction: maintenance strategy × replacement size — deepseek-flash

## TL;DR

- **The penalty grows with the size of a replacement, steeply and monotonically.** The application
  rollup (`n4_rollup_app`) runs 2.07x, 4.38x and 24.20x the trigger design (`n3_rollup_trigger`) on
  whole-configuration replacement at 10, 60 and 240 entries per installation.
- **`v2-13`'s 4.0–5.7x is one point on that curve, not a property of the design.** It was measured at the
  `small` scale's default 60 entries per installation. Quoting it without the size overstates the trigger
  design's cost at small configurations (2.07x at 10 entries) and badly understates it at large ones
  (24.20x at 240).
- **The trigger's throughput collapses with size while the application's merely declines**: 459 → 144 → 10
  ops/s for the trigger (a 46x fall), 952 → 631 → 242 for the application (3.9x). This is the mechanism
  `v2-13` asserted and listed as unmeasured: the trigger does work per row written, the application does
  one statement per replacement.
- **The tail grows faster than the throughput gap.** p99 is 3 / 44 / 156 ms for the trigger against
  1 / 3 / 8 ms for the application, a p99 ratio of 3.0, 14.7 and 19.5.
- **Correctness held at every tier for both designs**: every gate passed, and the audits recorded 0
  failures and 0 rollup mismatches across 600, 3,600 and 14,400 loaded entries. The trigger's structural
  guarantee is not traded away for size in this range.
- **This is the corpus's first `interaction` claim** — two factors crossed, maintenance strategy (2
  levels) × replacement size (3 levels), all other settings constant, measured in one session.

## The measured table

| Entries per installation | loaded | n3 ops/s | n4 ops/s | n4/n3 | p99 n3 ms | p99 n4 ms | p99 ratio |
|---|---|---|---|---|---|---|---|
| 10 | 600 | 459 | 952 | **2.07** | 3 | 1 | 3.0 |
| 60 | 3,600 | 144 | 631 | **4.38** | 44 | 3 | 14.7 |
| 240 | 14,400 | 10 | 242 | **24.20** | 156 | 8 | 19.5 |

`pg-single`, `small`, `verify,write`, three trials per cell (median), operation `w05_replace_all`.

**The tier-60 cell is a third observation of `v2-13`'s ratio.** It is the same operation at the same tier
as the replication's cells, measured on a different day, and it lands at **4.38x inside the replicated
range of 3.99–5.69x**. Given the ~30% day term the retry control measured, that is consistent — and it
means the interaction's middle cell corroborates the isolated claim rather than contradicting it.

## Prediction verdicts

| Id | Prediction | Verdict | Basis |
|---|---|---|---|
| R1 | the `n4/n3` ratio grows with entries per installation | **confirmed** | 2.07 → 4.38 → 24.20, a 12x spread across the axis, far outside the ±20% "no interaction" falsifier |
| R2 | the growth is monotone | **confirmed** | strictly increasing at all three points |
| R3 | correctness holds at every tier for both designs | **confirmed** | 0 audit failures and 0 rollup mismatches in all six cells |
| R4 | at tier 240 the trigger's p99 is at least 5x the application's | **confirmed** | 19.5x at tier 240, and already 14.7x at tier 60 |

## A crossing that could not be run

The first version of this experiment crossed maintenance strategy with concurrent writers. It cannot run:
`harness/main.go` guards the contention phase with `if !d.hasContention() { return nil }`, and the rollup
designs have no arbitration, so the phase **returns silently** with no result and no warning. Three runs
(`20260926T152250Z`, `20260926T152307Z`, `20260926T152326Z`) produced `concurrency: null` and are kept as
evidence of the feasibility failure. The protocol file records the revision and the reason. This is worth
knowing beyond this study: a phase that does not apply to a design produces no artifact, so "the phase
was skipped" and "the phase measured nothing" are indistinguishable in the results.

## Where this analysis is weak

1. **One run per tier, no repetition.** Three points on one axis, each measured once, so the trend has no
   spread at any level and the day term is unquantified within this tranche.
2. **No negative control in these runs.** `control_coverage` is recorded `absent`; the drift control
   belongs to the `v2` tranche and is not reused as evidence here.
3. **One operation.** `w05_replace_all` replaces the entire configuration. A workload that appends an
   entry instead would sit at the low end of this curve regardless of size, so the interaction is
   specific to whole-replacement writes.
4. **The `small` scale caps the axis.** 240 entries is the largest tier measured; the curve is still
   rising steeply there, so nothing here bounds where it flattens or whether it continues.
5. **Two designs, so no three-way structure.** The crossing says the penalty is size-dependent; it does
   not separate "the trigger fires per row" from "the trigger's per-row work is heavier than the
   application's single statement" — both predict the same shape.
6. **Two days are involved** if the tier-60 cell is read as corroborating `v2-13`; that comparison should
   carry the retry control's ~30% day term rather than being treated as exact.

## Next experiments

1. **Repeat the three tiers three times each**, with randomized order between the two designs, to put a
   spread on the curve and on the ratio at each level.
2. **Extend the axis past the `small` scale** — 480 and 960 entries per installation — to see whether the
   growth is superlinear, and where the trigger design stops being viable at all rather than merely slow.
3. **Cross the same decision with a second workload dimension** (single-entry publish rather than
   whole-replacement), which is where the two designs should converge rather than diverge.
4. **Give the harness a "phase not applicable" marker**, so a skipped phase is visibly skipped instead of
   absent — the defect this experiment stumbled into.
