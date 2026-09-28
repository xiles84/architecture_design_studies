---
analysis_id: 20260928T-v2-09-buyer-sweep--deepseek-flash--2026-09-28
run_id: 20260928T014452Z (128, p3 first), 20260928T014541Z (128, x1 first), 20260928T014640Z (8, x1 first), 20260928T014723Z (8, p3 first), 20260928T014821Z (64, p3 first), 20260928T014904Z (64, x1 first), 20260928T014949Z (16, x1 first), 20260928T015030Z (16, p3 first), 20260928T015111Z (32, p3 first), 20260928T015214Z (32, x1 first), 20260928T015257Z (128, p3 first, reference)
environment: host-zenbook-ux5406sa
analyst: deepseek-flash
analyst_kind: ai
analyst_version: "deepseek-flash (DeepSeek), HIGH-capability session acting as leader and executor, Deep Code CLI; the buyer-count sweep that isolates v2-09's confound."
analyzed_at: 2026-09-28
inputs_digest: ae9d870a666ca0bc (128 A), e9dc893bb3091bb7 (128 B), 21e416d823582162 (8 B), 9facab46793a5ada (8 A), a1fd9ed7620140cf (64 A), 6e675d46efe7c395 (64 B), 710fe5846835eb28 (16 B), 7646fd4cf36481b9 (16 A), 982ade6072c6def0 (32 A), 084f96a7e44b9d35 (32 B), de441ceb812660dd (128 A, reference)
inputs: 20260928T014452Z@ae9d870a666ca0bc, 20260928T014541Z@e9dc893bb3091bb7, 20260928T014640Z@21e416d823582162, 20260928T014723Z@9facab46793a5ada, 20260928T014821Z@a1fd9ed7620140cf, 20260928T014904Z@6e675d46efe7c395, 20260928T014949Z@710fe5846835eb28, 20260928T015030Z@7646fd4cf36481b9, 20260928T015111Z@982ade6072c6def0, 20260928T015214Z@084f96a7e44b9d35, 20260928T015257Z@de441ceb812660dd
repo_commit: 0a153f8bffd3114c880d36fe83c026862523f3de
status: current
headline: v2-09's design swap does not reproduce under control - at 1,000 seats the ledger is slower at every buyer count from 8 to 128 in both design orders (55 of 55 trial pairs), the magnitude is non-monotone with no crossing of 1, and the primary's apparent swap is explained by its own buyer-count/order confound.
---

# Analysis — study 02 `v2-09` buyer-count sweep: the swap was the confound — deepseek-flash

## TL;DR

- **There is no crossover.** At 1,000 seats the ledger is slower than compare-and-set in **55 of 55 trial
  pairs**, at every buyer count from 8 to 128, and in **both** design orders. Nothing crosses ratio 1.
- **`v2-09`'s swap is not reproduced.** The primary compared 32 buyers `p3`-first against 128 buyers
  `x1`-first — two factors at once. Controlled, the 32-buyer cell reads the ledger **slower** (1.233 and
  1.219), not faster as the primary's 0.585 did.
- **The magnitude is not monotone in buyers** (pooled medians 3.08, 1.32, 1.23, 1.59, 1.57 at 8/16/32/64/128)
  — my pre-registered *monotone* and *crossing* predictions are both **falsified**, and the 8-buyer level is
  polluted by one cell whose `x1` collapsed from ~4,800 to ~1,000 ops/s in its last three trials.
- **The order moves the absolute, not the sign.** At 32 buyers `p3` read 2,630 in one order and 6,499 in the
  other — a 2.5x difference — while the ratio agreed (1.233 vs 1.219). "The ratio is more portable than the
  throughput" holds even when the orders disagree by more than the effect.
- **Correctness green**: 0 over-booked, 0 under-booked, 0 timed-out, 0 errors in all eleven runs; the
  reference cell repeated at the two ends read 1.712 and 1.859 (8.6% apart).
- **Consequence for `v2-09`**: its statement ("the two designs change places") is contradicted by the
  controlled evidence and should be re-scoped to "direction stable, magnitude non-monotone"; that is a HIGH
  decision recorded with this analysis.

## What I was looking at

Eleven whole runs on `pg-single`, scale `small`, phases `verify,race`, designs `p3_precreated_cas` and
`x1_cas_ledger`, **tier 1,000 only**, `-race-trials 5`, at buyer counts **8, 16, 32, 64, 128**, with two runs
per level (one per design order) plus a 128-buyer `p3`-first cell repeated at the start and end. The buyer
levels ran in the non-monotone order 128, 8, 64, 16, 32, so a level is not confounded with its position in
the session. Twelve runs were planned and eleven executed — the pre-registered table already counted the
repeated reference as the eleventh. All trees clean, all tagged `run/02-ticket-booking/<run-id>` at
`0a153f8`, no failed cells.

I excluded the primary runs' `churn` cells: `20260921T183722Z` records six tier-1,000 rows per design, of
which three are `mode: churn`; only the three `race` rows are the comparable evidence, and the `p3/x1`
mediation there is 0.585 rather than the 0.618 a mixed-mode median gives.

## Conclusions

**The two primary arms were not a comparison, and the control shows it.** `v2-09` rests on 32 buyers with
`p3` first against 128 buyers with `x1` first — buyer count and order moved together, which the claim's own
limit admits. Crossing them at one tier removes the confound: the ledger is slower at 32 buyers in both
orders (1.233 and 1.219) and at 128 buyers in both orders (1.430 and 1.712/1.859). The primary's 0.585 at
32 buyers is not reproduced on the current harness and image, and the most likely reason is that the primary
run is not the same experiment: it ran phases `verify,write,race,churn` against this tranche's
`verify,race`, on a five-day-old bench image, so its state entering the race differed. That makes the
primary's low-buyer arm a period artefact as far as the *sign* is concerned, and it is the cleanest
illustration yet in this study that a confounded pair cannot be repaired by adding runs — only by varying
the confounded factor.

**My crossing hypothesis was wrong, and the data says so plainly.** I predicted a monotone rise in the cost
with the buyer count and a bounded crossing of ratio 1 between adjacent levels. Neither happened: the
pooled medians are 3.08, 1.32, 1.23, 1.59 and 1.57, so the cost is *smallest in the middle* (16–32 buyers)
and larger at both ends, and every cell is on the ledger-slower side. The 8-buyer level's 3.08 is not even a
level effect — the `p3`-first cell's `x1` fell from ~4,800 to ~1,000 ops/s across its last three trials
(a 5x mid-cell collapse) while the `x1`-first cell sat flat at ~5,000, so the level's ratio is an average of
1.06 and 5.10. **A cell that changes state mid-run is not a datum**, and with five trials the median of that
cell is not protected (unlike the `x1` medians in the tier-axis tranche, which were).

**What the direction is worth.** 55 of 55 trial pairs favour `p3`, across five buyer counts and both orders,
with a minimum ratio of 1.06. That is a stronger statement about the *direction* at 1,000 seats than the
claim currently makes, and it survives an order-induced 2.5x shift in the absolutes. What it does not
support is a ratio: the same eleven runs span 1.06 to 5.10 (and 1.22–1.86 if the collapsed cell is set
aside), so the cost is a range whose width is set by which cell you keep, exactly the "range whose endpoints
are the story" pattern the tier axis showed one level up.

**Prediction verdicts** (recorded in `V2-09-BUYER-SWEEP.md` before the runs):

| Id | Prediction | Verdict |
|---|---|---|
| B1 | pooled ratio rises monotonically with buyers | **falsified** — 3.08, 1.32, 1.23, 1.59, 1.57 |
| B2 | a bounded crossing of ratio 1 between adjacent levels | **falsified** — every cell is above 1; no crossing |
| B3 | the sign follows the buyer count, not the order | **confirmed** — both orders agree at every level |
| B4 | correctness green in every run | **confirmed** — 0 bad counters in all eleven |
| B5 | the reference cell repeats within 25% | **confirmed** — 1.712 vs 1.859, 8.6% apart |

Two falsified predictions are the useful part of this tranche: they turn `v2-09` from "the designs change
places" into "the direction is stable and the magnitude is not a function of buyer count", and they leave
the actual shape of the cost-vs-concurrency curve unmeasured rather than guessed.

## Where I think the measurement is weak

1. **One tier.** Everything here is 1,000 seats. The primary's swap was claimed across five tiers, and no
   controlled low-buyer-count sweep exists for the other four; the 128-buyer direction at those tiers is
   settled by `v2-08`, the low-buyer-count direction is not.
2. **One primary run per cell.** The reference cell shows an 8.6% session term, so levels whose two orders
   differ by more than that (the 8-buyer level, 1.06 vs 5.10) cannot be placed from these runs alone.
3. **The primary evidence is from a different harness state.** `20260921T183722Z` used phases
   `verify,write,race,churn` and a five-day-old image; comparing it to this tranche mixes the confound with
   a period effect, which is why the re-scope is stated as a direction, not a ratio.
4. **One host, one engine, the 2 GiB machine ceiling** applies to every number, as to the rest of study 02.
5. **The `p3`-first 8-buyer cell changed state mid-run** and is reported, not discarded; a discard would be
   invisible in the artifact.

## What I would measure next

1. **The low-buyer-count direction at the other four tiers**, one tier per tranche if needed, to see whether
   the stable direction at 1,000 seats is a property of the pair rather than of the tier.
2. **The 8-buyer cell again**, several runs, to decide whether its collapse is a harness artefact (a client
   stall) or a reproducible low-concurrency regime.
3. **The shape of cost versus concurrence** with more levels between 8 and 128 now that the sign is known to
   be stable — the deliverable would be the curve's minimum, not a crossing.
4. **The YugabyteDB half of `v2-09`**, which this tranche did not touch and which the claim still asserts.

## Disagreements with other analyses

None. This extends `20260921T183722Z--deepseek-flash--2026-09-21`, whose own limit said the confound could
not be isolated, and it agrees with the 2026-09-27 `v2-08` tranches that the ledger is slower at 128 buyers.
It **disagrees with the primary analysis's reading** that the 32-buyer arm shows the ledger faster: on the
current harness that reading does not reproduce, and the disagreement is resolved in favour of the
controlled sweep, with the period difference named as the residual.
