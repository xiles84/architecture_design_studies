---
analysis_id: 20260927T-v2-08-replication--deepseek-flash--2026-09-27
run_id: 20260927T213602Z (p3 then x1), 20260927T214133Z (x1 then p3)
environment: host-zenbook-ux5406sa
analyst: deepseek-flash
analyst_kind: ai
analyst_version: "deepseek-flash (DeepSeek), LOW/worker-capability session, Deep Code CLI; executing the study-02/v2-08-replication task."
analyzed_at: 2026-09-27
inputs_digest: c67336c6dcee4a99 (pass 1, pg-single), 0d64607d9c460c14 (pass 2, pg-single); baseline 6c88f6463c80c6ac. Corrected before commit: an earlier draft of this frontmatter carried an invented digest for pass 1; the value above is the one the run's own report states.
repo_commit: 677801ad73efbf512c9e65f5bafad8c27b1dbd4e
status: current
headline: The ledger's cost ratios reproduce in size and in both design orders, but the ordering reverses at the smallest tier (0.99x, ledger marginally faster) and the "tight spreads" claim fails there at 47-57%, so the claim's 1.1-2.3x band and its spread wording need re-scoping rather than silent acceptance.
---

# Analysis — study 02 `v2-08` replication: the ratio holds, the ordering and the spreads do not — deepseek-flash

## TL;DR

- **The ratios reproduce in size.** `p3/x1` came out **1.53x, 1.71x, 2.00x** at 100/1,000/10,000 seats in the
  first order and **0.99x, 1.30x, 1.74x** in the reversed order — every tier within ±30% of the baseline's
  1.30x / 1.43x / 1.79x, and all inside the claim's 1.1–2.3x band except one.
- **One ratio falls below the claim's floor and the ordering reverses with it.** At 100 seats the ledger read
  **3,618 ops/s against CAS's 3,588** — the ledger *marginally faster*, 0.99x, where the claim's band starts
  at 1.1x and the first pass measured 1.53x. R5's ordering clause is **falsified**; its ratio clause held.
- **"Tight spreads" fails at that same tier, and it is a property of the cell.** The 100-seat spread was
  **57%** for CAS in the first order and **47%** in the reversed one (and 51% for the ledger reversed), while
  every other cell sat at 4–17%. R6 is **confirmed**: the wide spread follows the tier, not the position in
  the session.
- **This connects `v2-08` to `v2-09` rather than contradicting it.** `v2-09` already records that this pair
  swaps places depending on the workload; the reversal here shows the swap needs no change of buyer count —
  the smallest tier is enough.
- **Correctness green in both orders**: 46/46 and 48/48, 0 errors, 0 over-booked, 0 timed-out.
- **What I did not do**: rewrite the statement. Its 1.1–2.3x band and its "tight spreads" wording are
  contradicted in part, and re-scoping a signed claim's numbers is a decision for HIGH, not a mechanical
  update for a worker. The evidence is recorded as limits and the decision is flagged.

## Both passes

| Tier | pass 1 `p3` / `x1` | ratio | pass 2 (reversed) `p3` / `x1` | ratio | baseline ratio |
|---|---:|---:|---:|---:|---:|
| 100 | 4,163 / 2,724 | 1.53x | 3,588 / 3,618 | **0.99x** | 1.30x |
| 1,000 | 6,021 / 3,514 | 1.71x | 5,346 / 4,126 | 1.30x | 1.43x |
| 10,000 | 7,123 / 3,553 | 2.00x | 6,694 / 3,855 | 1.74x | 1.79x |

Medians of three race trials per cell at 128 buyers, `pg-single`, `small`.

## Prediction verdicts

| Id | Prediction | Verdict | Basis |
|---|---|---|---|
| R1 | each ratio within ±30% of the baseline | **confirmed** | 1.53/1.71/2.00 against 1.30/1.43/1.79 |
| R2 | every cell's spread ≤40% | **falsified** | CAS at 100 seats spread 57% |
| R3 | the ledger is slower at every tier | **confirmed in pass 1** | slower at all three |
| R4 | correctness stays green | **confirmed** | 46/46, 48/48, no over-booking |
| R5 | ratio and ordering survive the reversal | **ratio clause confirmed, ordering clause falsified** | 0.99x at 100 seats |
| R6 | the small-tier spread is a cell property | **confirmed** | 47% in the reversed order, same tier |

## What this changes

1. **`v2-08` gains the replication** with `order: abba` and both passes as anchors, and four limits recording
   the reproduced ratios, the tier-100 reversal, the spread finding and its cell-specificity.
2. **A needed re-scoping is flagged, not taken.** The statement's band (1.1–2.3x) excludes the measured
   0.99x, and "tight spreads" is false for one of six cells. Both are signed wording; HIGH decides.
3. **The claim gets a boundary it lacked.** The ledger's cost is a function of tier: clearly positive at
   1,000 and above, indistinguishable from zero at 100 — which is a more useful statement than a single band.

## Where this analysis is weak

1. **Two runs, one session each, three trials per cell** — enough to see a reversal at the small tier, not
   enough to characterise it. The 0.99x is one median against another and could be noise; a third pass would
   settle it.
2. **The 100-seat cell is intrinsically noisy** (spreads 47–57%), so the tier where the claim's floor breaks
   is also the tier where the measurement is weakest — the least comfortable place to hang a conclusion.
3. **Two of five tiers and the whole YugabyteDB half were not re-measured**, as the protocol said in advance.
4. **The environment ceiling** (2 GiB machine against declared budgets) applies to both passes and to the
   baseline, so the comparison is like-for-like in that respect and unlike any other machine's.

## Next experiments

1. **A third pass at the 100-seat tier only**, several repeats, to decide whether 0.99x is a real reversal or
   the noise that tier's 50% spread would predict.
2. **The re-scoping decision** (HIGH): widen the band downward, or state the cost as tier-dependent.
3. **The YugabyteDB half**, untouched here, which the claim also asserts.
