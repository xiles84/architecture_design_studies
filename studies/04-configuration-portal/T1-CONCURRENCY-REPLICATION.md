# T1: replicating Study 04's concurrency claims — `v2-12` and `v2-17`

**Task:** `task-20261001T133759Z-t4-study04-replication` · **branch** `study-04/t1-replication`

**Goal.** Two claims about the same pair of designs under one hot key, both `single_run_directional`:

| claim | what it asserts |
|---|---|
| `v2-12-lock-vs-version-upper-bound` | at 16 writers a lock-before-update (`c2_pessimistic_lock`) delivered **1.76x** the throughput of a version check (`c1_optimistic_version`), called an **upper bound** because that run could not separate the lock's advantage from extra per-retry harness work; an unchecked read-modify-write (`x1_lost_update_control`) lost **93.7%** of acknowledged updates |
| `v2-17-optimistic-pessimistic-crossover` | the pair **crosses over between one and two writers**: the lock delivers 0.89x at 1 writer (1,177 against 1,329 ops/s), 1.10x at 2, 1.68x at 32 and 1.60x at 64, while the version check's own throughput falls 2.9x across the sweep and its conflict count rises from 0 to 21,301 |

**The tension this tranche tests.** `v2-12` still carries "an UPPER BOUND ... because the run could not
separate it from extra per-retry harness work", while `v2-17`'s limit records that the retry explanation was
**tested and failed** (with the retry budget at zero the 16-writer ratio was 1.60 against 1.61 with retries).
The same study therefore states the caveat in one claim and its refutation in another. Re-running the pair at
16 writers with the retry budget at its default and at zero is the cheapest way to settle which reading a
fresh run supports, and it is the reason this task runs the retry control rather than only the sweep.

**What is held constant.** The study's own runner, harness, correctness gates, dataset at scale `small`, the
published budgets, and phases `verify,contention` (the contention phase is what the claims measure; `explain`
and the read/write phases add cost and nothing to this question). **Varied: the writer count, the retry
budget, the run and the design order.**

**Design.** Designs `c1_optimistic_version`, `c2_pessimistic_lock` and `x1_lost_update_control` (the study's
declared negative control, which must lose updates). Levels **1, 2, 16 and 64 writers**, each level its own
whole run, in both design orders (the claim's crossover is between 1 and 2; 16 is the ratio both claims quote;
64 is the claim's other endpoint). **Retry control:** at 16 writers, one run with `-retries 0` and one with
the default, so the reading `v2-12` hedges and `v2-17` says was refuted is measured again.

## Required-comparison coverage (owner rule, 2026-09-14 — recorded before running)

| requirement | coverage |
|---|---|
| 1. cluster and client sizing, equal-total-budget control | Published budgets and footprint, unchanged. **Not applicable:** PostgreSQL single node, no cluster-size comparison, no equal-total arm; the claim is about contention on one key. |
| 2. rollup, rolldown and embedding | **Not applicable** to these two claims (concurrency strategies on one hot row). Study 04's `n3`/`n4` rollup and `n2` rolldown designs are other claims in the same study. |
| 3. colocated and non-colocated data | **Not applicable** (single node, no replication). |
| 4. **optimistic and pessimistic concurrency strategies** | **This is the owner's required comparison and this task is it**: `c1_optimistic_version` (conditional/version-checked, bounded retries) against `c2_pessimistic_lock` (lock-before-update) under the same business invariant (one hot configuration key), across a writer sweep, reporting throughput, the conflict count, retries and correctness — plus `x1_lost_update_control` as the negative control that must lose updates. |

## Prediction

| Id | Prediction | Falsifier |
|---|---|---|
| P1 | **The crossover exists at the same place**: `c2`/`c1` is **below 1 at 1 writer** and **above 1 at 2 writers**, in both design orders. | A ratio above 1 at 1 writer, or below 1 at 2. |
| P2 | **The lock's advantage at 16 writers is at least 1.4x** (the claim says 1.76x, one point on a curve), in both orders. | Any ratio below 1.4x at 16 writers. |
| P3 | **The retry budget does not explain the advantage**: at 16 writers the `c2`/`c1` ratio with `-retries 0` is within 10% of the ratio at the default budget, as `v2-17`'s limit asserts and `v2-12`'s upper-bound wording doubts. | A ratio difference of 10% or more between the two budgets. |
| P4 | **The version check degrades with the writer count** while the lock does not: `c1`'s throughput at 64 writers is below half its 1-writer level, and its conflict count rises by orders of magnitude across the sweep. | `c1` holding more than half its 1-writer throughput at 64, or a conflict count that does not rise. |
| P5 | **The negative control fires**: `x1_lost_update_control` loses a large fraction of acknowledged updates (the claim says 93.7%), and the two correct designs lose none. | A clean `x1`, or any lost update in `c1`/`c2`. |
| P6 | **Correctness and provenance**: `verify` clean in every cell, every run tagged from a clean tree, both orders per level. | A failed cell or a dirty tree. |

**What a pass decides.** Both claims move to `run_replicated`, with `v2-12`'s statement re-scoped so that its
"upper bound" wording reflects what the retry control shows: either the hedge stands (P3 fails) and the claim
keeps it, or the hedge is retired by two fresh budgets and the statement says so. `v2-17`'s crossover and
degradation are either reproduced at its own levels or the crossing point moves, in which case the statement
is re-scoped to the levels that do cross.

**Known limits, stated in advance.** PostgreSQL single node, scale `small`, one hot key, one measurement per
level per order, and no spread at any level unless the two orders disagree; absolute throughput is day-specific
and only within-run ratios travel; `x1`'s 93.7% is the primary's figure and this tranche reports its own.
