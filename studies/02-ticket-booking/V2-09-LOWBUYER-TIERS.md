# The low-buyer-count direction across the tier axis (study 02, `v2-09`)

**Question.** The 2026-09-28 buyer-count sweep isolated `v2-09`'s confound at **one tier** (1,000 seats) and
found the ledger slower in 55 of 55 trial pairs across 8–128 buyers in both orders. Its own first stated
weakness is that nothing exists above or below that tier: the primary run's 32-buyer arm read the ledger
*faster* at **every** tier (race-only `p3/x1` 0.377, 0.328, 0.585, 0.716, 0.781 at 10 / 100 / 1,000 /
10,000 / 100,000 seats), but with `p3` first and on the old harness, while the sweep's controlled 32-buyer
reading at 1,000 seats is 1.233 / 1.219. This tranche asks whether the primary's tier-wide "<1" reproduces
anywhere on the axis under control.

**What is held constant.** Designs `p3_precreated_cas` and `x1_cas_ledger`; topology `pg-single`; scale
`small`; phases `verify,race`; **32 buyers per event**; `-race-trials 5`; **all five tiers in every run**
(10, 100, 1,000, 10,000, 100,000 seats); dataset seed 42; the same organiser-edit cadence. The only factor
varied between the two runs is the **design order**.

**Design.** Two whole runs:

| # | buyers | tiers | order |
|---|---|---|---|
| 1 | 32 | all five | `p3`,`x1` |
| 2 | 32 | all five | `x1`,`p3` |

Two runs is enough for the direction question (the sweep already showed the sign is order-stable) and not
enough for a magnitude curve; the limits say so.

## Prediction

Recorded before measurement. The unit is the run-median `p3/x1` per tier, and the per-trial pair for the
direction.

| Id | Prediction | Falsifier |
|---|---|---|
| T1 | **The primary's tier-wide "<1" does not reproduce.** At 32 buyers the ledger is slower (ratio > 1) at **every** tier and in **both** orders — all 10 tier × order cells. | Any tier × order cell with ratio < 1, which would mean the primary's low-buyer reversal is real at some tier after all. |
| T2 | **The sign is order-stable at every tier**, as it was at 1,000 seats: the two orders agree on which design is faster at each tier. | A tier where the two orders disagree in sign. |
| T3 | **The 1,000-seat cell reproduces the sweep.** The 32-buyer tier-1,000 ratio falls within **±20%** of the sweep's 1.233 (p3 first) and 1.219 (x1 first). | Either order outside ±20%, which would make the sweep's own cell order-sensitive. |
| T4 | **The cost grows with the buyer count at every tier.** The 32-buyer ratio is **below** the 128-buyer pooled ratio measured earlier at the same tier (v2-08/tier axis: about 1.08, 1.24, 1.41, 1.85, 2.51). | A tier where the 32-buyer ratio meets or exceeds the 128-buyer one. |
| T5 | **Correctness stays green** in both runs: 0 over-booked, 0 under-booked, 0 timed-out, 0 errors. | Any failure, which voids that run's cells. |

**What a pass can and cannot decide.** It decides whether the direction is stable at 32 buyers across the
axis, and whether the cost keeps its growth with the buyer count at every tier. It does **not** give a
magnitude curve (one run per order), does not attribute anything to a mechanism, and says nothing about
YugabyteDB.

**Coverage gap, stated in advance.** `pg-single` only, one host, 32 buyers only, one run per order; the
2 GiB podman-machine ceiling applies to every number.

## Scope

| Field | Value |
|---|---|
| Runs | two, `pg-single`, two designs, `--race-buyers 32 -race-trials 5 -race-tiers 10,100,1000,10000,100000` |
| Compared with | primary `20260921T183722Z` (32 buyers, p3 first); sweep `20260928T01*` (32 buyers, tier 1,000); `v2-08` tier axis (128 buyers) |
| Unit of evidence | the run-median ratio per tier; the per-trial pair for direction |
| Boundary | One host, one engine, one buyer count, one run per order |
