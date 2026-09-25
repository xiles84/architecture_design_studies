---
analysis_id: 20260925T1-replication--deepseek-flash--2026-09-25
run_id: 20260925T190831Z, 20260925T191610Z, 20260925T194447Z
environment: host-zenbook-ux5406sa
analyst: deepseek-flash
analyst_kind: ai
analyst_version: "deepseek-flash (DeepSeek), HIGH-capability, Deep Code CLI; effort setting not exposed. Executor of task study-04/t1-replication."
analyzed_at: 2026-09-25
inputs_digest: per-run 829381816473fbd8 (20260925T190831Z), dc76ad71636584b4 (20260925T191610Z), 536c9eb4dcd53f67 (20260925T194447Z)
repo_commit: f9f201e
status: current
headline: Three randomized-order runs reproduce the direction and the large write-cost gap of the trigger-vs-application rollup pair (median 4.94x, range 3.99-5.69x) but not the previously cited 6.85x, and the trigger cell's own trials disagree by 35-115%, so the pair supports a direction and a range, not a point ratio.
---

# Analysis — Study 04 T1 replication: the rollup-maintenance pair, three runs — deepseek-flash

## TL;DR

- **P4 confirmed.** The application-maintained rollup (`n4`) was faster than the trigger (`n3`) on
  `w05_replace_all` in **all three** runs, at every order seed.
- **P1 direction confirmed, magnitude NOT reproduced.** Ratios were 3.99x, 4.94x and 5.69x (median
  **4.94x**); the previously cited single-run figure was 6.85x. The claimed magnitude is not a
  reproducible quantity.
- **The trigger cell is imprecise by the corpus's own rule.** `n3`'s three trials inside a run spanned
  `[40,159,162]`, `[22,124,165]` and `[61,91,93]` — spreads of 77%, 115% and 35%, all above the ~20%
  threshold above which methodology 8 says a cell cannot support an argument about a difference. `n4`
  was flat (spreads 1.4%, 4.4%, 6.8%).
- **A position effect is visible inside every run.** `n3`'s first trial was the slowest every time
  (40 vs a 159 median; 22 vs 124; 61 vs 91), i.e. 1.5–5.6x slower than the cell's own median. The
  median, not the mean, is what makes this cell reportable.
- **P3 confirmed.** The `x2_rollup_drift_control` fired in all three runs: 60 rollup mismatches each,
  audit `passed: false` on INV-8/INV-9.
- **P2 unresolved.** The committed scope ran no read phase, so the prediction that `n3` and `n4` read
  identically was not tested. It needs its own run.
- **New, and not in the earlier claim: the tail ratio is far larger than the throughput ratio.** `n3`
  p99 was 61–109 ms against `n4`'s 2.2–2.7 ms — a **22–46x** gap, versus 4.0–5.7x on throughput.

## What was measured

Three whole runs of the same three cells at `small` scale on `pg-single`, each with its own order seed
and run tag, three trials per cell, 5 s per trial, 2 s discarded warmup.

| Run | Order seed | `n3` w05 ops/s | `n4` w05 ops/s | ratio | `n3` p50 / p99 ms | `n4` p50 / p99 ms | `n3` trials | `n4` trials | x2 drift |
|---|---|---|---|---|---|---|---|---|---|
| `20260925T190831Z` | 101 | 159 | 635 | 3.99x | 6.11 / 97.33 | 1.51 / 2.15 | 40, 159, 162 | 635, 644, 635 | 60 |
| `20260925T191610Z` | 202 | 124 | 612 | 4.94x | 6.28 / 109.48 | 1.58 / 2.39 | 22, 124, 165 | 587, 612, 614 | 60 |
| `20260925T194447Z` | 303 | 91 | 518 | 5.69x | 10.58 / 61.24 | 1.86 / 2.72 | 61, 91, 93 | 518, 538, 503 | 60 |

Across-run medians: `n3` 124 ops/s, `n4` 612 ops/s. Across-run spread
`(max − min) / median`: **54.8%** on `n3`, **19.1%** on `n4`, **34.4%** on the ratio.

Design order was shuffled per topology from each run's order seed and the order is recorded in every
run manifest. Every cell passed its correctness gate; every `x2` audit failed exactly as the control
is designed to.

## Prediction verdicts

| Id | Prediction | Verdict | Basis |
|---|---|---|---|
| P1 | n4 materially faster on `w05_replace_all`; the cited 6.85x reproduces | **direction confirmed, magnitude not reproduced** | 3.99/4.94/5.69x vs a cited 6.85x; no reversal in any run, so the stated falsifier did not fire — but the expected magnitude is outside the observed range |
| P2 | no attributable read difference between n3 and n4 | **unresolved** | the committed scope ran `verify,write` only; no read phase was executed |
| P3 | `x2` drift non-zero in every run | **confirmed** | 60 mismatches and `passed: false` in all three runs |
| P4 | the ordering n4 > n3 is stable across runs | **confirmed** | 3/3 runs |

## Mechanism consistent with the result

The trigger fires once per row touched by the whole-configuration replacement, so its cost scales with
the number of configuration entries in the replaced document and it serialises work inside the
statement; the application path recomputes the overview once per replacement. That is consistent with
(a) a ~4–6x throughput gap, (b) an order-of-magnitude worse tail, and (c) high variance in the trigger
cell — a per-row trigger's cost is sensitive to row order, cache state and autovacuum pressure in a way
a single bulk recompute is not. This is an observation consistent with the mechanism, not a
measurement of it; isolating *why* the first trial is slow needs its own experiment.

## Where this analysis is weak

1. **The ratio is not a point.** The `n3` cell violates the 20% rule in two of three runs; the honest
   output is a range (3.99–5.69x), not the median alone.
2. **P2 was not tested.** The prediction about reads was written beside a scope that excluded reads.
   That is a defect in the committed protocol, not a null result, and it is recorded as unresolved
   rather than quietly dropped.
3. **Three runs, one machine, one scale, one engine.** `run-level` replication only: it bounds
   run-to-run variance and does not remove the environment's hazards, and it is not
   hardware-independent.
4. **The `x2` control's speed is not evidence.** `x2` is a deliberately wrong design; its throughput is
   reported only so the audit's firing is visible, never as a result.
5. **The read-gain clause of the original claim is untouched.** The ~15x read gain came from the cited
   run and was not re-measured here.
6. **Wall-clock variance between runs was large** (6, 26 and 6 minutes for identical work), which is
   itself a warning about drawing magnitude conclusions from any single run on this host.

## What this changes in the registry

The active claim `v2-13-trigger-rollup-cost` cited a single run and a 6.85x point ratio. That paragraph
no longer describes what the corpus supports. The corrected claim should state the replicated
direction, the observed range, the imprecision of the trigger cell, the tail ratio, and the fact that
the read gain is from the earlier run and was not re-measured. Promoting it properly is a registry
change; the numbers it needs are in the table above.

## Next experiments

1. A read-phase replication of the same pair, to resolve P2.
2. A first-trial-versus-steady-state experiment on `n3`: warm the trigger path before the timed trials,
   or report the first trial separately, to establish whether the 2.5–6x first-trial penalty is plan
   cache, buffer cache or bloat.
3. A cardinality sweep of the same pair (entries per installation), which is where the trigger's per-row
   cost should cross the application's single recompute — the crossover T2 asks for.
