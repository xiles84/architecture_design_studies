---
analysis_id: 20260926T-retry-control--deepseek-flash--2026-09-26
run_id: 20260926T142926Z, 20260926T143014Z
environment: host-zenbook-ux5406sa
analyst: deepseek-flash
analyst_kind: ai
analyst_version: "deepseek-flash (DeepSeek), HIGH-capability, Deep Code CLI; effort setting not exposed. Executor of task study-04/retry-control."
analyzed_at: 2026-09-26
inputs_digest: per-run 9f170ae830a8c736 (20260926T142926Z, retries=0), 3f18ecc053c1c0e4 (20260926T143014Z, retries=8)
repo_commit: 14cafaa
status: current
headline: Removing the retry budget does not change the pessimistic/optimistic ratio at 16 writers (1.61 -> 1.60), and the design that never retries moves by the same amount as the one that does - so the lock's advantage is not explained by retry work, and the mechanism sentence in v2-17 must be withdrawn.
---

# Analysis — Study 04 retry control: the retry budget is not the mechanism — deepseek-flash

## TL;DR

- **The prediction that the version check's throughput depends on retrying is falsified.** With the retry
  budget set to zero, `c1_optimistic_version` did not collapse — it went from 657 to **725** ops/s. The
  claim it was meant to support ("the version check pays for retries the lock does not") is **refuted, and
  `v2-17`'s mechanism sentence has been withdrawn** rather than left standing.
- **The ratio did not move: `c2/c1` is 1.61 with retries and 1.60 without.** If retry work were the
  mechanism, removing it would have narrowed the gap. It did not narrow it at all.
- **The decisive detail is the design that never retries.** `c2_pessimistic_lock` has `retries = 0` in
  every run, at any setting of the flag, and it still moved 1,056 → 1,163 ops/s (+10%) between the two
  arms. So the arm-to-arm difference is a **run-level effect shared by both designs**, not a retry
  effect. Comparing 552 → 725 for c1 and 869 → 1,163 for c2 across days would have looked like a 30%
  retry effect on both sides at once, which is not a mechanism.
- **Cross-day comparisons carry a ~30% whole-run term.** `c2` at the same settings, same writer count,
  same phases measured 869 ops/s on 2026-09-23 and 1,056 on 2026-09-26. Absolute ops/s figures in this
  corpus are day-specific; only within-run ratios are safe to carry across days.
- One thing did hold: conflicts are a property of the workload, not of the budget (26,014 with the budget
  removed, up from 24,176 with it).

> **Correction, 2026-09-26, same session.** This analysis first reported the contention-phase counter
> invariant as *not recorded* and P4 as only weakly confirmed. That was wrong: the invariant **is** in every
> result file under `final_counter_value`, `expected_counter_value` and `lost_updates`, and the report's
> contention table renders it. The error came from querying the fields I guessed (`final_value`,
> `expected_value`), getting `null` for a key that does not exist, and reading my own miss as an absence in
> the artifact. P4 is now **confirmed**: with the retry budget removed, `c1` ended at 4,948 against 4,948
> expected and `c2` at 7,917 against 7,917, both with **0 lost updates** — so rejecting publications does
> not break the counter invariant, which is a positive result the first version of this file threw away.
> The verdict table and the weaknesses list below are corrected in place, and the lesson is recorded.

## What was measured

Same day, same topology, scale, phases (`verify,contention`), writer count (16) and trials; only the
retry budget differs.

| Arm | c1 acked | c1 retries | c1 conflicts | c1 ops/s | c2 acked | c2 ops/s | c2/c1 |
|---|---|---|---|---|---|---|---|
| retries = 0 (`9f170ae…`) | 4,948 | 26,014 | 26,014 | 725 | 7,917 | 1,163 | **1.60** |
| retries = 8 (`3f18ecc…`) | 4,615 | 1,325 | 24,176 | 657 | 7,445 | 1,056 | **1.61** |

Three days earlier, retries = 8, `verify,read,write,contention`: c1 552 ops/s, c2 869 ops/s (1.58).

## Prediction verdicts

| Id | Prediction | Verdict | Basis |
|---|---|---|---|
| P1 | c2 within 20% of 870 ops/s | **falsified as written, confirmed same-day** | 1,163 is +34% against the cross-day baseline but +10% against the same-day one. The falsifier named a baseline two days old; that was a defect in the prediction, not in the design. |
| P2 | c1 collapses by >50% or the cell fails | **falsified** | c1 rose to 725 ops/s; the cell passed. The mechanism sentence is withdrawn. |
| P3 | conflicts still > 10,000 | **confirmed** | 26,014 conflicts with no retry budget. |
| P4 | the correctness gate still passes | **confirmed** (corrected — see the note below) | the gate passed (16 checks), and the counter invariant holds under rejected publications: `c1` ended at 4,948 with 4,948 expected and `c2` at 7,917 with 7,917 expected, both **0 lost updates**. |

## What this changes

1. **`v2-17`** no longer asserts "the version check pays for retries the lock does not". It now records
   that the retry budget was set to zero, the ratio did not narrow, and the mechanism is **unseparated**.
   The crossover stands; the explanation for it does not.
2. **`v2-12`'s caveat is now stronger, not weaker.** The 16-writer upper bound was flagged because the
   run could not separate the lock's advantage from per-retry harness work. The control has now tried to
   separate it and failed to find an effect, so the caveat should read "not separated, and the attempt to
   separate it found no retry effect" rather than "could not be separated".
3. **Absolute throughput must not be compared across days.** A ~30% same-settings drift was measured
   between 2026-09-23 and 2026-09-26 on a design that does not retry.

## Where this analysis is weak

1. **One run per arm.** The ratio difference between arms (1.61 vs 1.60) is far below any resolvable
   effect at this precision, which is the point — but it also means a real retry effect smaller than
   ~10% would not be visible.
2. **The two arms realized different design orders** (`c2` first with retries = 0, `c1` first with
   retries = 8), and the position effect measured in this study's read phase was ~1.55x. The orders are
   recorded, but with one run per arm order and budget are confounded here. This does not rescue the
   retry hypothesis (both arms cannot be explained by retries plus order without order doing all the
   work), but it does mean the *magnitudes* in the table above are not clean.
3. **`-retries 0` is not "no conflict handling".** It makes the optimistic path abandon on first
   conflict, so the two arms differ in *what the design does*, not only in an internal counter. A cleaner
   arm would keep retries but exclude retry work from the measured interval.
4. **A rejected publication's cost is not measured.** The counter invariant confirms nothing was lost, but
   nothing here says how much work the 26,014 abandoned attempts consumed — the run measures outcomes, not
   wasted effort.
5. **One writer count.** This tests the mechanism sentence at 16 writers; it says nothing about whether
   the crossing *moves* when the budget changes.

## Next experiments

1. **A proper retry-cost arm**: keep the retry budget and measure the retry path's own latency, rather
   than removing retries. That separates cost from abandonment.
2. **Repeat both arms three times with randomized order**, so order and budget are not confounded — the
   same discipline the read tranche needed.
3. **Re-measure a fixed reference cell on every measurement day** (one design, one writer count) so the
   day term is recorded rather than discovered.
