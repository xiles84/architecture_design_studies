# Settling the tier-100 reading: is the ledger-race reversal real? (study 02, `v2-08`)

**Question.** The `v2-08` replication (2026-09-27, analysis
`20260927T-v2-08-replication--deepseek-flash--2026-09-27`) reproduced the ledger's cost ratios in both
design orders, but at **100 seats** the reversed pass read `p3/x1 = 0.99x` — the ledger marginally faster —
where the claim's band starts at 1.1x and the first pass read 1.53x. Its own next-experiment list asks for
a third pass at that tier. This tranche runs it: **is the 0.99x a real property of the 100-seat tier, or a
median taken from a cell whose trials disagree by ~90%?**

The design differs from the earlier passes in one way only: **more trials per cell** (5 instead of 3) and
**six runs instead of two**, in three AB and three BA orders. Deciding it from run medians alone is what
produced the ambiguous 0.99x; the unit here is the **per-trial pair**, so the sign of the difference is
tested 30 times rather than 6.

**What is held constant.** Designs `p3_precreated_cas` and `x1_cas_ledger`; topology `pg-single`; scale
`small`; phases `verify,race`; **128 buyers per event**; dataset seed 42; the same organiser-edit cadence
(every 20 ms). Only the session differs.

**Baseline being settled** (PostgreSQL, 1 node, 128 buyers — from `20260927T213602Z` and
`20260927T214133Z`):

| Tier | pass 1 `p3` / `x1` | ratio | pass 2 (reversed) `p3` / `x1` | ratio |
|---|---:|---:|---:|---:|
| 100 | 4,164 / 2,725 | 1.53x | 3,588 / 3,619 | 0.99x |
| 1,000 | 6,021 / 3,514 | 1.71x | 5,346 / 4,126 | 1.30x |

Per-trial tier-100 pairs already measured (6 pairs): 1.30, 1.65, 1.53, 1.28, **0.97**, 1.12 — so the single
`0.99x` reading is one pair of six, and the trial-to-trial spread of a single design at this tier is 50-92%
of its own median. The third pass tests whether that picture survives more pairs.

## Prediction

Recorded before measurement. `spread` is the report's definition, `(max - min) / median * 100`.

| Id | Prediction | Falsifier |
|---|---|---|
| P1 | **The reversal is noise, not the tier.** Across the six runs' tier-100 cells (6 × 5 = 30 trial pairs), `p3` is faster than `x1` in **at least 20 of 30 pairs (≥67%)**. | **≤15 of 30 (half)** — CAS slower or tied at least half the time — which would make `0.99x` a property of the 100-seat tier rather than a session reading. 16–19 of 30 is inconclusive and is reported as such. |
| P2 | **The pooled ratio still shows a cost at the tier, inside the claim's band.** The median of the 30 per-trial `p3/x1` ratios is **in [1.10, 1.60]** — at or above the claim's 1.1x floor, and no larger than the first pass's 1.53x. | Median **< 1.10**, which would mean the claim's floor is met by the higher tiers but not this one; or **> 1.60**, which would mean the first pass understated this tier and the claim's band describes the wrong range. |
| P3 | **Positive control at the neighbouring tier.** In **every one of the six runs**, the tier-1000 ratio is **> 1.10**, the cost the two earlier passes saw there (1.71x, 1.30x). | Any run ≤ 1.10 at 1,000 seats — the harness would have failed to detect a cost it detected twice before, which makes that run's tier-100 reading uninterpretable rather than informative. |
| P4 | **Correctness stays green in every run**: the `verify` phase passes and every race cell records 0 over-booked, 0 under-booked and 0 timed-out events. | Any failure, which voids that run's cells. |
| P5 | **The wide spread is the cell's, and persists.** The tier-100 cell spread exceeds **40%** in **at least 4 of the 6 runs** (for either design), as it did in both earlier passes (47–57%). | Spread ≤ 40% at tier 100 in **3 or more** runs, which would mean the earlier wide spreads were session effects and the reversal candidate does not live in an intrinsically noisy cell. |

**What a pass can and cannot decide.** P1–P5 decide whether the reversal is a cell property or a session
artefact. They do **not** re-scope the claim: its `1.1–2.3x` band and "tight spreads" wording are signed,
and changing them is a HIGH decision (recorded as such by the 2026-09-27 analysis). This tranche supplies
the evidence that decision needs.

**Coverage gap, stated in advance.** Only tiers **100 and 1,000**, only `pg-single`, one host. The
YugabyteDB half, the refund ratio, the effect on the organiser's own edit, and the 10-seat and 10k-seat
tiers are **not** re-measured. A pass on P1–P5 says nothing about them. The 2 GiB podman-machine ceiling
applies to every number, as it does to the baseline.

## Scope

| Field | Value |
|---|---|
| Runs | six, `pg-single`, two designs, `--race-buyers 128 -race-trials 5 -race-tiers 100,1000` |
| Orders | three `p3` then `x1`, three `x1` then `p3` (A B B A A B) |
| Unit of evidence | the per-trial pair (30 of them), not the run median |
| Compared with | `20260927T213602Z`, `20260927T214133Z`; primary anchor `20260921T205212Z` |
| Boundary | One host, one engine, two of five tiers, one scale; within-cell vs across-run spread are not distinguished by design, and the first-trial warm-up seen in both earlier passes is not controlled |
