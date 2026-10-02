# T1: replicating `v2-10` — guarded confirmations refused on YugabyteDB

**Task:** `task-20261001T133759Z-t3-study03-replication` · **branch** `study-03/t1-replication`

**Goal.** `v2-10` is Study 03's only `single_run_directional` claim and it rests on one v1 matrix run: *"On
YugabyteDB a guarded confirmation is not enough: correct designs refused valid holds, and retrying the refused
statement once removed every refusal; on PostgreSQL the guarantee held everywhere."* This tranche re-runs the
pair the claim is about in **two whole runs in both design orders**, with PostgreSQL as the control.

**The mandated feasibility check, done first — and it passed.** The task requires a feasibility check before
any run because the claim's substance is on YugabyteDB and this machine's podman configuration allots it
2 GiB. Result on 2026-10-02 (`20261002T150337Z-feas`, `yb-single`, scale `tiny`, designs `s1_conditional_update`
and `s1r_confirm_retry`, phases `verify,race`): **the topology started, both cells completed with `errors=0`,
audits consistent, and the claim's own mechanism appeared** — at the 10,000-seat tier `s1` recorded
`immediate_confirms_refused=2` with `confirm_retries=0`, while `s1r` recorded `confirm_retries=1` and
`confirm_retry_successes=1`, i.e. the retry removed the refusal. So the arm is feasible and the run is the
replication, not a substitute for it. **The feasibility run is a pre-check and is labelled as one**: it is
`tiny`, one trial, one order, and it is not counted as a replication.

**Correction to the machine record, worth carrying forward.** The 2 GiB blocker documented in `CONTEXT.md` is
about **multi-node** topologies. A single-node YugabyteDB cell fits and completes cleanly here, which is why
`v2-10` can be replicated at all; `yb-cluster3` remains the part that does not fit. The claim's own limit and
`v2-11`'s YugabyteDB gap do not change, but "no YugabyteDB work on this machine" was too broad.

**What is held constant.** The study's own runner, harness, correctness gates, dataset at scale `small`, and
the published budgets and client footprint. **Varied: the run, the design order, and the engine** (`yb-single`
for the claim, `pg-single` for its control).

**Design.** Designs `s0_check_then_hold_rc`, `s1_conditional_update`, `s1r_confirm_retry` and
`s2_lock_then_update` — the guarded-confirmation pair the claim names plus the check-then-hold and lock-then-
update designs, so "correct designs refused valid holds" is tested on more than one design. Phases
`verify,race`: refusals are recorded by the race phase (`immediate_confirms_refused`, `confirm_retries`,
`confirm_retry_successes`, `deferred_refused`), and nothing here needs the read or write phases.
**Runs:** four whole runs — `yb-single` and `pg-single`, each in forward and reversed design order.

## Required-comparison coverage (owner rule, 2026-09-14 — recorded before running)

| requirement | coverage |
|---|---|
| 1. cluster and client sizing, equal-total-budget control | Published budgets and footprint; the same declared envelope for both engines. **Not applicable:** one node per engine, no cluster-size comparison and no equal-total arm — this claim is about refusals, not capacity. Note recorded: the *configured* machine limit is 2 GiB while containers see the 15.4 GiB WSL VM, so the declared budgets are not enforced (the documented sizing defect) and that applies to both arms equally. |
| 2. rollup, rolldown and embedding | **Not applicable** to `v2-10`: it is about hold confirmation and refusals. Study 03's `l2`/`l3` designs are its embedded variants and are not part of this claim. |
| 3. colocated and non-colocated data | **Not applicable** (single node per engine; no replication). Outstanding for `yb-cluster3`, which does not fit. |
| 4. optimistic and pessimistic concurrency | **Covered by design selection:** `s1`/`s1r` are the optimistic conditional-confirmation pair and `s2_lock_then_update` is the pessimistic lock-then-update arm; the refusal counts are reported for each, so the two strategies are compared on the same invariant. |

## Prediction

| Id | Prediction | Falsifier |
|---|---|---|
| N1 | **On YugabyteDB at least one design records a refused valid confirmation** (`immediate_confirms_refused > 0` or `deferred_refused > 0`) in **both** runs, as the claim asserts. | A run with zero refusals on every YB design. |
| N2 | **The retry variant removes every refusal it sees**: in `s1r_confirm_retry`, `confirm_retry_successes == confirm_retries` and `immediate_confirms_refused == 0`, in both runs. | Any unsuccessful retry, or a refusal in `s1r` that no retry covers. |
| N3 | **PostgreSQL is the control that holds**: no design records a refusal or a retry on `pg-single` in either run. | Any non-zero refusal or retry counter on PostgreSQL. |
| N4 | **Refusals concentrate at the largest tier measured** (10,000 seats), as the feasibility run showed, and do not appear at the two smallest tiers. | A refusal at the 10- or 100-seat tier, or the largest tier refusing nothing while a smaller one does. |
| N5 | **Correctness stays green** in every cell: the audit reports consistent, no undersold or leaked seats, no errors, and both runs are tagged from a clean tree. | A failed cell, an inconsistent audit, or a dirty tree. |
| N6 | **The refusal count is not an artefact of design order**: the same designs record refusals (or none) in both orders. | A design that refuses only in one order. |

**What a pass decides.** `v2-10` moves to `run_replicated` with its statement re-scoped to the tiers and
designs measured, and its limit on "a single measurement per configuration" replaced. A failure of N3 — a
PostgreSQL refusal — would be a **finding against the claim's control half** and would keep the claim's
strength with a limit naming it, because the claim's whole point is the engine difference.

**Known limits, stated in advance.** Scale `small`, one node per engine, one host; the refusals are the
harness's own classification of its own designs; and the claim is about refusals, so no throughput conclusion
follows from anything here.
