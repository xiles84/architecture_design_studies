# Isolating `v2-09`'s confound: does the ledger's sign follow the buyer count or the design order?

**Question.** `v2-09-postgres-ledger-ratio-unsettled` records that on PostgreSQL "the two designs change
places depending on the workload (32 simultaneous buyers, P3 first: X1 faster; 128 buyers, X1 first: P3
faster)", and its own limit says **"buyer count and design order changed together, so neither can be
isolated"**. The two primary arms are exactly that pair: `20260921T183722Z` ran **32 buyers with `p3` first**
and read the ledger *faster* at every tier (race-only `p3/x1` 0.377, 0.328, 0.585, 0.716, 0.781), while
`20260921T205212Z` ran **128 buyers with `x1` first** and read the ledger *slower* at every tier (1.057,
1.288, 1.444, 1.797, 2.338). This tranche crosses the two factors at one fixed tier.

**What is held constant.** Designs `p3_precreated_cas` and `x1_cas_ledger`; topology `pg-single`; scale
`small`; phases `verify,race`; **tier 1,000 seats only**; `-race-trials 5`; dataset seed 42; the same
organiser-edit cadence. The factor varied is **buyer count** (8, 16, 32, 64, 128) and, separately, the
**design order** within each run.

**Design.** Eleven runs, each a whole run of the two designs at one buyer count — the ten
(level, order) cells plus one repeated reference:

| # | buyers | order | role |
|---|---|---|---|
| R0 | 128 | `p3`,`x1` | reference cell (start) |
| 1 | 128 | `x1`,`p3` | |
| 2 | 8 | `x1`,`p3` | |
| 3 | 8 | `p3`,`x1` | |
| 4 | 64 | `p3`,`x1` | |
| 5 | 64 | `x1`,`p3` | |
| 6 | 16 | `x1`,`p3` | |
| 7 | 16 | `p3`,`x1` | |
| 8 | 32 | `p3`,`x1` | |
| 9 | 32 | `x1`,`p3` | |
| R11 | 128 | `p3`,`x1` | reference cell (end) |

The buyer levels run in the order 128, 8, 64, 16, 32 — **deliberately not monotone**, so a level is not
confounded with its position in the session — and the two orders alternate. R0 and R11 are the same cell
repeated at the two ends, so the session term is measurable rather than inferred (the corpus's "the machine
wanders" lesson).

## Prediction

Recorded before measurement. The unit is the run-median `p3/x1` ratio, and the per-trial pairs are the
direction check.

| Id | Prediction | Falsifier |
|---|---|---|
| B1 | **The ratio rises with the buyer count.** Pooled run medians are ordered 8 < 16 < 32 < 64 < 128. | Any adjacent inversion in the pooled medians, which would make the buyer count not the axis. |
| B2 | **There is a bounded crossing of ratio 1.** The pooled medians cross 1.0 between two **adjacent** levels, are below 1 at every level under it and above 1 at every level over it, and do not cross back. This would be the corpus's second bounded crossover. | No adjacent pair straddling 1.0, or a second crossing (two reversals), either of which makes the sign non-monotone rather than a boundary. |
| B3 | **The sign follows the buyer count, not the order.** At every level the two orders agree about which design is faster (same side of 1.0); the order reverses positions, not the sign. | A level where the two orders disagree in sign — then the reversal is (also) an order effect, and no crossing may be registered from one order's data. |
| B4 | **Correctness stays green** in every run: `verify` passes and every race cell records 0 over-booked, 0 under-booked, 0 timed-out, 0 errors. | Any failure, which voids that run's cells. |
| B5 | **The reference cell is stable.** The two 128-buyer `p3`-first readings (R0, R11) differ by **less than 25%**, so a level effect larger than that is not the session. | A difference of 25% or more, which would put the measurement's session term at the same size as the effect and force the crossing to be reported as a band, not an interval. |

**What a pass can and cannot decide.** It decides whether the sign flip is a function of buyer count and
whether it is bounded. It does **not** attribute the flip to a mechanism: the pair changes CAS-plus-ledger
versus CAS alone, so a component-level explanation would need the isolation arms the `v2-08` attribution
already withholds. It also says nothing about other tiers, and nothing about YugabyteDB.

**Coverage gap, stated in advance.** `pg-single` only, one host, tier 1,000 only, and the 2 GiB
podman-machine ceiling applies to every number. The 32-buyer and 64-buyer cells are the noisiest
candidates, so their spreads are reported beside every ratio.

## Scope

| Field | Value |
|---|---|
| Runs | eleven, `pg-single`, two designs, `--race-buyers <level> -race-trials 5 -race-tiers 1000` |
| Compared with | `20260921T183722Z` (32 buyers, p3 first), `20260921T205212Z` (128 buyers, x1 first) |
| Unit of evidence | the run-median ratio per (buyers, order); the per-trial pair for direction |
| Boundary | One host, one engine, one tier; one reading per (buyers, order) plus a two-point reference |
