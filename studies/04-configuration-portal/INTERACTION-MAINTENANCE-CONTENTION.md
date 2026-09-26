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

---

## Revision, same session, before measuring the revised design

**The crossing above cannot be run, and the reason is in the harness.** Three runs were started
(`20260926T152250Z`, `20260926T152307Z`, `20260926T152326Z`) and each produced `concurrency: null`: the
contention phase returned without recording anything. `harness/main.go`:

```go
case "contention":
    // A design with no arbitration has nothing to contend with -- except the
    // drift control, whose race is on the derived rollup rather than on a key
    // and which therefore must still be run here.
    if !d.hasContention() {
        return nil
    }
```

`n3_rollup_trigger` and `n4_rollup_app` have no arbitration, so `hasContention()` is false and the phase
returns `nil` — **silently**, with no result and no warning in the output. P1's falsifier anticipated the
weak form of this ("the contention-phase operation does not exercise the trigger"); the actual outcome is
that it does not run at all. The three runs are kept as evidence of the feasibility failure and are not
analysed as an interaction.

**The revised crossing is the one `v2-13` already names as a gap.** That claim's limits say "a per-row
trigger's cost grows with the rows a replacement touches", and record it as **unmeasured**. The write
phase does run for these designs, and the tier (entries per installed product) is settable, so the
maintenance decision can be crossed with the size of a replacement:

| Id | Prediction | Falsifier |
|---|---|---|
| R1 | **The interaction.** The `n4/n3` ratio on `w05_replace_all` **grows** with entries per installation across tiers 10, 60 and 240, because the trigger fires once per row written while the application's rollup update is one statement per replacement. | The ratio stays within ±20% across the three tiers (the penalty is size-independent), or it shrinks. |
| R2 | The growth is monotone: the ratio at tier 10 is below the ratio at tier 60, which is below tier 240. | A non-monotone sequence, which would point at a threshold rather than a trend. |
| R3 | Correctness holds at every tier for both designs: the audits pass with no rollup mismatch (INV-8/INV-9), so the trigger's structural guarantee is not traded for size. | Any audit failure or rollup mismatch, which would void the cell regardless of the timing. |
| R4 | At tier 240 the trigger cell's p99 latency is at least 5x the application cell's at the same tier. | A p99 ratio below 5x, which would mean the tail does not scale with replacement size even when the throughput does. |

Evidence form: **`interaction`** — maintenance strategy (2 levels) × entries per installation (3 levels),
all other settings constant, one run per tier. Recorded before measuring the revised design.

**Boundary of the revision.** This is three points on one axis in one workload on one host; it tests
whether the *penalty changes with size*, and cannot find a crossover, because no second design is
available at a larger size band. Scale remains `small`, so tier 240 (14,400 entries) is the largest cell
measured here.

