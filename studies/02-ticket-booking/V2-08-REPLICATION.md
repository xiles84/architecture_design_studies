# Replication: the ledger's cost in the sell-out race (study 02, `v2-08`)

**Question.** `v2-08-ledger-race-cost-128-buyers` says that at 128 buyers an append-only ledger costs
**1.1–2.3x** of the sell-out race on PostgreSQL, *"with tight spreads"*, and gives ratios for YugabyteDB and
for the refund. It rests on **one run** (`20260921T205212Z`), one measurement per configuration. Its sibling
`v2-09` already records that the same pair of designs **swaps places** under different buyer counts, so the
ratio's stability is not a settled matter in this study — which makes "tight spreads" worth testing rather
than accepting.

**What is held constant.** Designs `p3_precreated_cas` and `x1_cas_ledger`; topology `pg-single`; scale
`small`; phases `verify,race`; **128 buyers per event**; three race trials per tier; dataset seed 42; the
same organiser-edit cadence (every 20 ms). Only the session differs.

**Baseline being replicated** (PostgreSQL, 1 node, 128 buyers, 3 trials — from `20260921T205212Z`):

| Seats per event | `p3` CAS | `x1` CAS+ledger | ratio | reported spreads |
|---|---:|---:|---:|---|
| 100 | 4.8k | 3.7k | 1.30x | 37% / 11% |
| 1k | 5.7k | 4.0k | 1.43x | 9% / 3% |
| 10k | 7.0k | 3.9k | 1.79x | 6% / 4% |

## Prediction

| Id | Prediction | Falsifier |
|---|---|---|
| R1 | **The ratio reproduces in size.** Each tier's `p3/x1` ratio falls within **±30%** of the baseline ratio for that tier (1.30x, 1.43x, 1.79x). | Any tier outside ±30%, which would move the claim's 1.1–2.3x band or add a third "unsettled" finding to this pair. |
| R2 | **"Tight spreads" holds.** Every cell's within-cell spread is **≤40%**, as the baseline's were (its widest was 37%). | A spread above 40% in any cell, which would make the ratio a band rather than the point the claim treats it as. |
| R3 | **The ordering is stable.** `x1_cas_ledger` is slower than `p3_precreated_cas` at **every** tier, as in the baseline. | A tier where the ledger is faster, which would connect this claim to `v2-09`'s sign reversal rather than leave them as separate findings. |
| R4 | **Correctness stays green**: all gates pass, no over-booked or under-booked seats, no timed-out events. | Any failure, which voids that cell. |

Recorded before measurement. Evidence form stays `isolate` (a two-design pair differing in whether a ledger
is written); this tranche raises the claim's replication tier if R1 and R3 hold.

**Coverage gap, stated in advance.** The claim's YugabyteDB figures, its refund ratio and its effect on the
organiser's own edit are **not** re-measured here, and the 10-seat and 100k-seat tiers are omitted to keep
the run short. A pass on R1–R4 says nothing about those.

## Scope

| Field | Value |
|---|---|
| Runs | one, `pg-single`, two designs, `--race-buyers 128 -race-trials 3 -race-tiers 100,1000,10000` |
| Compared with | `20260921T205212Z` (primary anchor of `v2-08`) |
| Boundary | One host, one engine, three of five tiers, one scale; the 2 GiB podman-machine ceiling applies to every number, as it does to the baseline. |
