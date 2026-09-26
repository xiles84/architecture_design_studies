# Interaction: rollup maintenance × write contention (Study 04)

**Question.** `v2-13` measures the trigger-maintained versus application-maintained rollup pair under a
*write* workload at one concurrency level, and finds the trigger design 4.0–5.7x slower on
whole-configuration replacement. This experiment crosses that decision with the number of concurrent
writers on one hot key, to ask whether the maintenance penalty is **additive, constant, or absorbed**: 
does the trigger's extra work stay on the critical path as contention rises, or does it disappear into
queueing behind the same hot row?

The corpus has **no registered interaction claim** — every claim is `isolate`, `package` or
`not_applicable`. This is the smallest well-posed crossing available: both designs already exist, the
contention phase already exists, and the only new axis is the writer count.

## Prediction

| Id | Prediction | Falsifier |
|---|---|---|
| P1 | At **1 writer** the trigger design (`n3_rollup_trigger`) delivers lower contention throughput than the application design (`n4_rollup_app`), because the trigger does its work on the critical path of every write. | `n3` matches or beats `n4` at one writer, which would mean the contention-phase operation does not exercise the trigger. |
| P2 | **The interaction.** The `n4/n3` ratio **shrinks** as writers rise from 1 to 16: at high contention both designs queue on one hot row, so the trigger's extra work is absorbed rather than added. | The ratio stays within ±20% across the three writer counts (no interaction — the penalty is constant), or it grows (the penalty compounds). |
| P3 | The counter invariant holds for both designs at every writer count: `lost_updates` is 0, because neither design is the unsafe control. | Any `lost_updates` above 0, which would be a correctness failure and would void the cell. |
| P4 | Neither design's `retries` counts become large, since the rollup pair differs in who maintains the rollup, not in how a publication is guarded. | Large retry counts, which would mean the contention phase is exercising a guard decision rather than a maintenance decision, and the crossing would not be about maintenance. |

Recorded before measurement, from the committed tree. Evidence form: **`interaction`** — two factors
crossed, maintenance strategy (2 levels) × concurrent writers (3 levels), with every other setting held
constant across all six cells and measured in the same session.

## Scope

| Field | Value |
|---|---|
| Topology | `pg-single` |
| Scale | `small` |
| Designs | `n3_rollup_trigger`, `n4_rollup_app` |
| Phases | `verify,contention` |
| Writers | 1, 4, 16 (one run each) |
| Trials | 3 per cell |
| Extra flags | none |

**Boundaries.** The writer counts are three points, so a reversal is not testable and is not claimed; this
is a test of whether the *penalty changes*, not a crossover. The contention phase's operation is not the
`w05` whole-configuration replacement `v2-13` measures, so the absolute ratio here is not comparable to
`v2-13`'s, and the analysis must not present it as such. One host, one key, one engine, single node, one
run per writer count.

**What would make this an interaction claim rather than a package claim.** The two factors must be varied
independently, which they are: `n3`/`n4` differ only in who maintains the rollup, and the writer count is
the same setting in both arms at each level. The claim is only registrable as `interaction` if P2's effect
is visible; if the ratio is constant, the honest registration is "no interaction detected in this range",
which is itself a finding and does not get an `interaction` label.
