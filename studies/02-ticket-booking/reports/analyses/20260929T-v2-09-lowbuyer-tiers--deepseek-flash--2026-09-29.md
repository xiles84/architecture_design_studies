---
analysis_id: 20260929T-v2-09-lowbuyer-tiers--deepseek-flash--2026-09-29
run_id: 20260929T002128Z (32 buyers, p3 first), 20260929T003126Z (32 buyers, x1 first)
environment: host-zenbook-ux5406sa
analyst: deepseek-flash
analyst_kind: ai
analyst_version: "deepseek-flash (DeepSeek), leader-and-executor session, Deep Code CLI; the 32-buyer tier-axis tranche that closes v2-09's first stated weakness."
analyzed_at: 2026-09-29
inputs_digest: b98961fb6b61b338 (p3 first), cda90300d6206a0f (x1 first)
inputs: 20260929T002128Z@b98961fb6b61b338, 20260929T003126Z@cda90300d6206a0f
repo_commit: adfdd12daa707589fd209ea2dca71776548da4a3
status: current
headline: At 32 buyers the ledger is slower at every tier from 10 to 100,000 seats in both design orders (all ten cells, 37 of 50 trial pairs), so the primary's tier-wide reversal does not reproduce anywhere on the axis; at 10 and 100 seats the direction sits at the floor and is unreliable pair by pair.
---

# Analysis — study 02 `v2-09`, the 32-buyer tier axis — deepseek-flash

## TL;DR

- **The primary's tier-wide reversal does not reproduce.** `20260921T183722Z` read the ledger *faster* at
  every tier at 32 buyers (0.377–0.781). Two controlled runs give ratios **above 1 at all ten tier × order
  cells**, from 1.021 (10 seats, `x1` first) to 1.606 (100,000 seats, `p3` first).
- **The cost grows with the buyer count at every tier.** The 32-buyer pooled ratio is below the 128-buyer
  pooled ratio at all five tiers: 1.06 vs 1.08, 1.06 vs 1.24, 1.23 vs 1.41, 1.56 vs 1.85, 1.56 vs 2.51.
- **At the two smallest tiers the direction is at the floor.** Only **3 of 10** trial pairs at 10 seats and
  6 of 10 at 100 seats favour `p3`, and those cells spread 86–174%; the medians are above 1 but a per-pair
  reader would not call it. At 1,000 seats and above it is 9/10 and 10/10.
- **The 1,000-seat cell reproduces the sweep exactly.** 1.245 and 1.219 against the sweep's 1.233 and 1.219
  — the same cell, a day apart, in both orders.
- **Correctness green** in both runs: 0 over-booked, 0 under-booked, 0 timed-out, 0 errors.
- **No re-scope is warranted**: the statement's core (direction stable, swap was the confound) stands; this
  tranche extends its tier coverage and adds the small-tier reliability caveat as a limit.

## What I was looking at

Two whole runs on `pg-single`, scale `small`, phases `verify,race`, designs `p3_precreated_cas` and
`x1_cas_ledger`, **32 buyers per event**, **5 race trials**, **all five tiers in every run** (10, 100,
1,000, 10,000, 100,000 seats), one run in each design order (`20260929T002128Z` `p3` first,
`20260929T003126Z` `x1` first). Both trees clean, both tagged `run/02-ticket-booking/<run-id>` at
`adfdd12`, no failed cells. I excluded no cells; the `race` filter is the one the sweep established, since
the file can also carry `churn` rows.

## Conclusions

**The primary's low-buyer reversal was the confound, not a regime.** `v2-09`'s original claim rested on
32 buyers with `p3` first, where the ledger looked faster at every tier. Controlled two ways — order
reversed, and the whole tier axis re-measured on the current harness — it does not appear: all ten cells
put the ledger behind, and the pooled ratio rises with tier exactly as it does at 128 buyers (1.06, 1.06,
1.23, 1.56, 1.56). The direction, which the 2026-09-28 sweep established at one tier, is now measured at
five.

**The cost is a function of both axes, and buyer count dominates.** At 32 buyers the tier curve spans
1.02–1.61; at 128 buyers it spans 1.08–2.51. So the tier raises the cost by about 1.5x and the buyer count
raises it by up to 1.6x at the top tier, and the two are not interchangeable: the same ratio (~1.06) can be
reached by 32 buyers at 10 seats or by 128 buyers at 10 seats, but at 100,000 seats only the buyer count
moves it (1.56 vs 2.51). That is a second axis of the `v2-08` finding rather than a new claim, and it is
recorded as a limit rather than in the statement.

**The small tiers are a floor, not a value.** At 10 seats only 3 of 10 pairs favour `p3` and the cells
spread 88–140%; at 100 seats, 6 of 10 and 27–174%. The medians still say ledger-slower, but the honest
reading is that at 32 buyers and 10–100 seats the ledger's cost is inside the measurement's own noise, and a
claim that quoted 1.02–1.10 there as a cost would be quoting the floor. This matches the 128-buyer
tier-axis finding that the 10-seat cell's direction flips in 40% of pairs; the two tranches now agree that
**the smallest tier is where this comparison stops being resolvable**, at both concurrency levels.

**The 1,000-seat cell is stable across a day.** 1.245 (`p3` first) and 1.219 (`x1` first) against the
sweep's 1.233 and 1.219 — within 1% and exactly equal. Given the corpus's repeated finding that absolute
throughput is day-specific and that this host wanders, a ratio reproducing that closely a day later, in
both orders, is a small positive: it is the kind of portability the `v2-08` work found for ratios over
absolutes, now shown across sessions rather than within one.

**Prediction verdicts** (recorded in `V2-09-LOWBUYER-TIERS.md` before the runs):

| Id | Prediction | Verdict |
|---|---|---|
| T1 | the primary's tier-wide "<1" does not reproduce: all 10 cells > 1 | **confirmed** — smallest 1.021 |
| T2 | the two orders agree in sign at every tier | **confirmed** — 5 of 5 tiers |
| T3 | the 1,000-seat ratio within ±20% of the sweep's 1.233 / 1.219 | **confirmed** — 1.245 / 1.219 |
| T4 | the 32-buyer ratio below the 128-buyer ratio at every tier | **confirmed** — 5 of 5 tiers |
| T5 | correctness green | **confirmed** — 0 bad counters |

## Where I think the measurement is weak

1. **One run per order.** The direction is settled; the magnitudes are two readings each and the small
   tiers are noisy, so the tier curve at 32 buyers is an indication, not a curve.
2. **32 buyers only.** Nothing here says what happens below 8 or between 32 and 128 at the outer tiers;
   the sweep covered the buyer axis at 1,000 seats only.
3. **The 10- and 100-seat cells spread 86–174%**, so their medians are weak summaries and are reported
   beside the per-pair counts rather than instead of them.
4. **One host, one engine, and the primary is five days and a harness change away** — this tranche agrees
   with the sweep that the primary's low-buyer arm is a period artefact, but it cannot say which of the two
   (phases, image, or the confound) removed the reversal.
5. **The 2 GiB podman-machine ceiling** applies to every number, as to the rest of study 02.

## What I would measure next

1. **The YugabyteDB half of `v2-09`** — the claim still asserts a 1.7–5.7x range there, never replicated,
   and it is the last un-re-measured arm of both ledger claims. It is also the one the machine's 2 GiB
   ceiling makes risky, so it needs the owner's environment decision first.
2. **A third and fourth reading of the 10-seat and 100-seat cells at 32 buyers**, where the direction is
   unresolvable, to decide whether the floor is a property of the tier or of the client.
3. **The buyer × tier surface at one more buyer count** (say 64) to separate the two axes' contributions at
   the top tier, where 32 vs 128 buyers moves the ratio most (1.56 vs 2.51).

## Disagreements with other analyses

None. This agrees with the 2026-09-28 sweep that the primary's 32-buyer reversal is a confound, and it
agrees with the 27 September tier-axis tranche that the 10-seat cell is the least resolvable. It extends
both; no other analysis covers these digests.
