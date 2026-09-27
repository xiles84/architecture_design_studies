# The ledger's cost across the tier axis (study 02, `v2-08`)

**Question.** `v2-08-ledger-race-cost-128-buyers` states the ledger costs **1.1–2.3x** of the sell-out race at
128 buyers. The primary run measured all five tiers, and its own ratios run **1.06 → 1.30 → 1.43 → 1.79 →
2.35** from 10 to 100,000 seats: the band's floor is the 100-seat tier and its top is the 10,000-seat tier,
so **both endpoint tiers are already outside the band** in the evidence the band was drawn from. The two
replication tranches since then covered tiers 100, 1,000 and 10,000 only, so the endpoints — the tiers where
the claim is most likely to be a boundary rather than a value — have never been re-measured.

This tranche measures **all five tiers in every run**, in both design orders, so the tier curve is a
within-run comparison rather than a set of ratios stitched across sessions. It is a coverage extension of
the 2026-09-27 settle (`V2-08-TIER100-SETTLE.md`), not a re-scoping: the signed band is untouched.

**What is held constant.** Designs `p3_precreated_cas` and `x1_cas_ledger`; topology `pg-single`; scale
`small`; phases `verify,race`; **128 buyers per event**; `-race-trials 5`; dataset seed 42; the same
organiser-edit cadence (every 20 ms). Only the tier list and the session differ from the settle.

**Baseline being extended** (PostgreSQL, 1 node, 128 buyers, primary `20260921T205212Z`, 3 trials):

| Seats per event | `p3` | `x1` | ratio | report spreads |
|---|---:|---:|---:|---|
| 10 | 1.9k | 1.8k | **1.06** | 21% / 27% |
| 100 | 4.8k | 3.7k | 1.30 | 37% / 11% |
| 1,000 | 5.7k | 4.0k | 1.43 | 9% / 3% |
| 10,000 | 7.0k | 3.9k | 1.79 | 6% / 4% |
| 100,000 | 8.0k | 3.4k | **2.35** | 4% / 4% |

The 2026-09-27 replication reproduced tiers 100/1,000/10,000 in both orders (1.53/1.71/2.00 and
0.99/1.30/1.74), and the settle then ran tier 100 in six more runs (six run medians 1.257–1.712; ledger
slower in 29 of 30 trial pairs). So the middle of the axis is well covered and the ends are not.

## Prediction

Recorded before measurement. `spread` is the report's definition, `(max - min) / median * 100`.
Per-trial pairs are the unit for the direction test; run medians are the unit for the axis.

| Id | Prediction | Falsifier |
|---|---|---|
| P1 | **The cost is monotone in tier.** The pooled run-median ratio increases from the 10-seat tier to the 100,000-seat tier without an adjacent inversion (10 < 100 < 1,000 < 10,000 < 100,000). | Any adjacent inversion in the pooled medians, which would make the cost a function of something other than the tier. |
| P2 | **The floor is a boundary, not a value.** The 10-seat ratio is **below the claim's 1.1x floor**, reproducing the primary's 1.06. | A 10-seat ratio **≥ 1.10** — the band would then cover the smallest tier after all, and the 1.06 would be a session reading. |
| P3 | **The ceiling is a boundary too.** The 100,000-seat ratio is **above the claim's 2.3x top**, reproducing the primary's 2.35. | A 100,000-seat ratio **≤ 2.30**, which would mean the primary's top-tier ratio was inflated and the band does cover the whole axis. |
| P4 | **The ledger is slower at every tier.** Across each tier's per-trial pairs (4 runs × 5 trials = 20 per tier), `p3` is faster in **at least 80%** of pairs. | Any tier where `p3` is faster in **fewer than 80%** of pairs — the smallest tier is the candidate, since that is where the ratio is near 1. |
| P5 | **Correctness stays green** in every run: `verify` passes and every race cell records 0 over-booked, 0 under-booked, 0 timed-out. | Any failure, which voids that run's cells. |

**What a pass can and cannot decide.** P1–P5 test whether the tier axis and its two endpoints reproduce.
They do **not** re-scope the claim: its `1.1–2.3x` band and "tight spreads" wording are signed, and changing
them is a HIGH decision (flagged by both 2026-09-27 analyses). This tranche supplies the boundary evidence
that decision needs.

**Coverage gap, stated in advance.** Only `pg-single`, one host; the 100,000-seat tier is a **single event**
per trial, so its ratio rests on one event and the corpus's `n=1` caveat applies to it more than to the
others. The YugabyteDB half, the refund ratio and the organiser-edit effect are **not** re-measured. The
2 GiB podman-machine ceiling applies to every number.

## Scope

| Field | Value |
|---|---|
| Runs | four, `pg-single`, two designs, `--race-buyers 128 -race-trials 5 -race-tiers 10,100,1000,10000,100000` |
| Orders | two `p3` then `x1`, two `x1` then `p3` (A B B A) |
| Unit of evidence | the per-trial pair (20 per tier), and the within-run tier curve |
| Compared with | primary `20260921T205212Z`; replication `20260927T213602Z`/`20260927T214133Z`; settle `20260927T22*` |
| Boundary | One host, one engine, one scale; the endpoint tiers have never had a repetition other than this; tier 100,000 is one event |
