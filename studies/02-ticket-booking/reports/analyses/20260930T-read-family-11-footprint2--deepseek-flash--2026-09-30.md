---
analysis_id: 20260930T-read-family-11-footprint2--deepseek-flash--2026-09-30
run_id: 20260930T201732Z (order A), 20260930T203913Z (order B)
environment: host-zenbook-ux5406sa
analyst: deepseek-flash
analyst_kind: ai
analyst_version: "deepseek-flash (DeepSeek), LOW worker session, Deep Code CLI; completes the read family's footprint-2 footing (eleven remaining designs)."
analyzed_at: 2026-09-30
inputs_digest: 17c1586f3e78956d (order A), 8ca41db9b81e969d (order B)
inputs: 20260930T201732Z@17c1586f3e78956d, 20260930T203913Z@8ca41db9b81e969d
repo_commit: c5c6452d29d66097d35edea7b9f0f7df152d9499
status: current
headline: The eleven remaining designs complete the read family's footprint-2 footing - none throttles (0 of 631 periods in all 22 read phases) and the ordering is order-stable (0.91-1.03), but the diagnostic refines - cells whose query returns many rows are bounded by the client's row decoding (in-flight 0.7-1.5) rather than the store (1.8-2.0), and the primary anchor's levels were inflated in proportion to how much its client was throttled per design (1.09x to 2.22x).
---

# Analysis — the read family complete at footprint 2 (study 02) — deepseek-flash

## TL;DR

- **`H1` confirmed: nothing throttles.** All **22 read phases** (eleven designs × two orders) record
  **`nr_throttled = 0`** of ~631 periods, including the five designs the primary anchor throttled hardest
  (`c4` 714, `h1` 719, `p4` 712, `h0` 672, `r2` 637 of 1085). With the four designs re-measured earlier, all
  **fifteen** of the study's designs now have a footprint-2 foot.
- **`H3` confirmed, and better than the last time**: the two orders agree on `q01`@1,000 to within
  **0.91–1.03** for every design, and no pair differing by ≥1.5x swaps — against a 33% order term in one
  earlier session. The ordering: `p4`/`c4` ≈ 18.4k > `h1`/`r2` 17.5k > `h0` 16.7k > `c2`/`c3` 14k >
  `c5`/`r3`/`p2`/`p1` 11.4–12.2k.
- **`H2` falsified as stated, and the diagnostic improves for it.** In-flight (`ops/s × p50`) is **1.8–2.0**
  for the point-lookup and single-row cells but **0.7–1.5** for `q02_band_events` and
  `q05_band_tickets_sold` — the two queries that **return 28 rows**. Those cells are bounded by the client's
  own row decoding, not by the store: the workers are idle between fetches. So `in-flight ≈ workers` means
  store-bound, and `in-flight < workers` means the client's per-fetch work bounds the loop — a client limit
  that is **not** throttling.
- **This corrects an over-broad claim of my previous analysis** ("every read cell is latency-bound at 1.9"),
  which was measured on `q01`@1,000 only. The correction is recorded in the registry limits.
- **The primary anchor's read levels were inflated per design, in proportion to its client's throttling.**
  Against `20260920T234953Z` (footprint 8, day-confounded), every design's `q01`@1,000 falls — from **1.09x**
  (`p1`, throttled 121/1085 then) to **2.22x** (`c4`, throttled 714/1085), with the five most-throttled designs
  falling 1.75–2.22x and the four least-throttled 1.09–1.47x. That is the footprint artefact quantified: the
  old table's *order* was mostly right, its *levels* were inflated by up to 2.2x, and its inflation tracked
  the client, not the design.

## What I was looking at

Two whole runs on `pg-single`, scale `small`, phases `verify,explain,read`, published budgets and duration,
the eleven designs this tranche covers in reversed order, at the read default of **2 workers**
(`client_read_conns: 2` in both manifests). Both trees clean, both tagged at `c5c6452`, no failed cells.
Reference: `20260920T234953Z` (the same eleven plus four others at **8** workers, a ten-day-older harness),
used only as a stated, day-confounded expectation and never scored.

## Conclusions

**The family now has one footing, and the footing is honest.** Read-phase client CPU at 2 workers ranges
**40.8–68.0 s** against a 126 s allowance, and no design throttles — including the five that lost 600–719
periods each in the primary anchor. The designs that use the most client CPU (`c4` 67.1/68.0 s, `p4`
66.8/67.8, `h1` 65.3/65.7, `h0` 65.5/65.3, `r2` 63.3/63.7) are exactly the ones the old run throttled
hardest, which is the mechanism the calibration predicted: client demand, not design cost, was setting those
cells.

**Two kinds of cell, and only one is store-bound.** For `q01` (availability), `q03` (a customer's tickets) and
`q04` (a ticket by id) the in-flight count sits at 1.8–2.0 against 2 workers: the workers are waiting on the
store, and the level is a store measurement. For `q02` and `q05` — both of which return 28 rows — it sits at
0.7–1.5: the workers spend most of the loop decoding rows, so the level is a **client** measurement even at
the reduced footprint. The distinction matters for what the corpus quotes: `q05`'s `c1`-versus-`p3` result is
an **ordering** (a design with no per-seat rows has no per-seat aggregate to compute) and its *levels* are not
store numbers at any footprint measured so far.

**What the primary anchor's read table was.** Comparing per design at `q01`@1,000, the fall from footprint 8
to footprint 2 ranks exactly with how many periods the old run throttled that design: `p1` 121 → 1.09x, `p2`
73 → 1.37x, `r3` 91 → 1.25x, `c5` 132 → 1.47x, `c2` 141 → 1.69x, `h0` 672 → 1.75x, `r2` 637 → 1.81x, `c3`
127 → 1.89x, `h1` 719 → 2.09x, `p4` 712 → 2.14x, `c4` 714 → 2.22x. The 2026-09-20 ordering of designs is
therefore broadly reproducible, and its levels carry a design-dependent inflation of up to 2.2x — which is
why they were marked unquotable rather than merely caveated.

**Prediction verdicts** (recorded in `V2-READ-FAMILY-11-FOOTPRINT2.md` before the runs):

| Id | Prediction | Verdict |
|---|---|---|
| H1 | no design throttles at the new footprint | **confirmed** — 0 of ~631 periods in all 22 read phases |
| H2 | every cell latency-bound (`ops/s × p50` in 1.6–2.4) | **falsified** — true for `q01`/`q03`/`q04` (1.8–2.0), false for the two many-row queries (0.7–1.5), which are client-decode-bound |
| H3 | no ≥1.5x pair swaps between orders | **confirmed** — order agreement 0.91–1.03, no swaps |
| H4 | client read CPU ≤ 63 s (half the quota) for every design | **falsified as stated** — 10 of 22 phases are 63.7–68.0 s (1.01–1.08 cores of 2); no design throttles, so the substance of the prediction holds and the threshold was mine, not a property |
| H5 | correctness green, clean tagged trees | **confirmed** — no failed cells, both runs tagged at `c5c6452` |

## Where I think the measurement is weak

1. **One run per order** (as for the four-design tranche): `p4` 18.4k / 19.3k and `h1` 17.5k / 19.2k show a
   session term of up to 10%, so single levels should not be quoted from one run.
2. **The client-decode finding is inferred**, not directly measured: the harness records no per-operation
   client time, so "the workers are decoding rows" is the explanation that fits `in-flight < workers` for the
   two 28-row queries, not an instrument reading.
3. **The footprint-8 comparison is day- and harness-confounded**, as stated in the protocol; only its rank
   pattern (which tracks the recorded per-design throttling) is used here.
4. **`q02`/`q05` levels remain client-bound at the only footprint measured**, so those two queries have no
   store-bound level from this program yet; a single-worker run would test whether they become store-bound,
   at the cost of halving the in-flight count everywhere.
5. **One topology, one engine, one scale, one host**; the machine is a stated project constraint, and the YB
   read half is untouched.

## What I would measure next

1. **A single-worker read run** for `q02` and `q05` to see whether their in-flight count reaches 1 and their
   levels become store-bound; if it does not, those queries should be quoted as orderings only.
2. **The primary anchors' own read cells at footprint 2** (the four extra designs' cells included), so the
   old tables can be replaced rather than caveated.
3. **A second run per order** wherever a level rather than an ordering is to be quoted.

## Disagreements with other analyses

It completes `20260930T-read-family-footprint2` and **corrects one over-broad sentence in it** ("every cell
is latency-bound at 1.9"), which came from measuring only `q01`@1,000; the corrected statement is in the
registry limits and in this file. It agrees with `20260930T-read-q01-plan-check` (the `q01` flatness is O(1)
by design — visible here: `p4`, `c4` and `r2` read 18.4k at every tier) and with the calibration's mechanism
(client demand, not design cost, decided the old throttling).
