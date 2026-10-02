---
analysis_id: 20261002T-t1-refusal-replication--deepseek-flash--2026-10-02
run_id: 20261002T150903Z-t1yb (yb forward), 20261002T152604Z-t1yb (yb reversed), 20261002T154340Z-t1pg (pg forward), 20261002T155134Z-t1pg (pg reversed)
environment: host-zenbook-ux5406sa
analyst: deepseek-flash
analyst_kind: ai
analyst_version: "deepseek-flash (DeepSeek), HIGH session, Deep Code CLI; T1 replication of v2-10 (guarded confirmations refused on YugabyteDB)."
analyzed_at: 2026-10-02
inputs_digest: 5603bdaa1e6045cb (yb forward), f4e6b822c66767c2 (yb reversed), a6759a618001f494 (pg forward), aaa426d437c3c060 (pg reversed)
inputs: 20261002T150903Z-t1yb@5603bdaa1e6045cb, 20261002T152604Z-t1yb@f4e6b822c66767c2, 20261002T154340Z-t1pg@a6759a618001f494, 20261002T155134Z-t1pg@aaa426d437c3c060
repo_commit: 0d57e6ef
supersedes:
status: current
headline: v2-10 replicates with a sharper mechanism than it states - on YugabyteDB the correct guarded designs refused valid holds at the 10,000-seat tier (2 and 3 refusals for s1, 1 for s2) and every single refusal was classified transient, while the retry variant refused nothing and every retry it made succeeded; on PostgreSQL the same three designs refused nothing at any of five tiers in two whole runs.
---

# Analysis — T1 replication of `v2-10` — deepseek-flash

## TL;DR

- **The claim replicates, and the refusals are spurious rather than a design fault.** On `yb-single`,
  `s1_conditional_update` refused valid holds **2 times** in the forward run and **3 times** in the reversed
  run, and `s2_lock_then_update` once in each — with **every refusal classified `transient`**, which the
  harness defines as *the identical statement, re-issued inside the refusing transaction, matching every
  seat*. That is the claim's mechanism with a stronger warrant than the claim itself gives: it is not that the
  design is wrong, it is that the engine can refuse a confirmation that would have succeeded.
- **The retry fix holds exactly as claimed.** `s1r_confirm_retry` recorded **no refusals at any tier** in
  either run and made 5 retries in the forward run, **5 of them successful**.
- **PostgreSQL is the control that holds.** All three correct designs recorded **zero refusals, zero retries
  and zero correctness failures at all five tiers (10 to 100,000 seats) in both runs**.
- **The negative control fires on both engines**, which is what makes the silence meaningful: the study's
  declared control `s0_check_then_hold_rc` refused 242–351 holds per tier on YugabyteDB and 555–3,526 per tier
  on PostgreSQL, and **stole 474–599 holds per tier on YugabyteDB and 1,168–6,231 on PostgreSQL**. The audit
  therefore detects both the spurious refusal and the real theft, and PostgreSQL is not simply a quieter
  instrument.
- **The feasibility check is why this ran at all.** The task mandated one before any measurement because
  `CONTEXT.md` records a 2 GiB podman-machine blocker for YugabyteDB work. A single-node YB arm fits: the tiny
  pre-check started, ran both designs, and already showed `s1` refusing twice at 10,000 seats with `s1r`
  retrying successfully. **The blocker is about multi-node topologies**; `yb-cluster3` remains unmeasurable,
  but "no YugabyteDB work on this machine" was too broad a reading.
- **Two of my predictions were falsified**, both my error: `N3` because I forgot the negative control is
  *supposed* to refuse on PostgreSQL too, and `N4` because a reversed-run refusal landed at the 1,000-seat
  tier for `s2`. Both are recorded below rather than smoothed over.

## What I was looking at

Four whole runs on 2026-10-02, each with four designs at scale `small`, phases `verify,race`, tagged by the
runner: two on `yb-single` (forward and reversed design order) and two on `pg-single` as the control (forward
and reversed). The refusals are read from the harness's own audit of each race cell
(`audit.ledger.rejected_early`, `rejected_early_transient`), not from the client's `errors` counter — a
distinction that mattered: `s0`'s client-level `errors` run to thousands while its audit is what classifies
the failure.

| prediction | verdict |
|---|---|
| N1 correct designs refuse valid confirms on YB in both runs | **confirmed** — `s1`: 2 (fwd) and 3 (rev) at 10k; `s2`: 1 each |
| N2 the retry variant removes every refusal | **confirmed** — `s1r`: 0 refusals, 5 retries / 5 successes (fwd), 0 retries needed (rev) |
| N3 PostgreSQL is the control that holds | **confirmed for the correct designs, falsified as written** — I forgot the negative control must refuse on PostgreSQL too; it does, by thousands |
| N4 refusals concentrate at the largest tier | **falsified as stated** — true in the forward run, but the reversed run put `s2`'s single refusal at 1,000 seats while its 10,000-seat cell refused nothing |
| N5 correctness green, control firing | **confirmed** — no correct-design failure in any cell; the control stole holds at every tier on both engines |
| N6 the refusal count is not a design-order artefact | **confirmed** — the phenomenon appears in both orders, though not at the same tier for `s2` |

## Conclusions

**`v2-10` can move to `run_replicated`, and its statement should say what the refusals are.** The claim as
written says *"correct designs refused valid holds, and retrying the refused statement once removed every
refusal"*. The replication confirms both halves and adds the classification the claim leaves implicit: every
refusal in a correct design was `transient`, meaning a single re-issue matched every seat. A reader should
take away that **a guarded confirmation can be refused spuriously on this engine and that one retry is the
whole fix** — not that the design was at fault.

**The PostgreSQL half is a real control, not an absence of evidence.** Zero refusals is only good news
because the same harness, on the same designs, with the same workload, *did* refuse and *did* record stolen
holds for the negative control on PostgreSQL — at a higher rate than on YugabyteDB. Without the control, "no
refusals on PostgreSQL" could have been an artefact of a quieter instrument; with it, the engine difference is
the explanation.

**The refusals are few, and the claim must not be read as a rate.** Two, three and one events across two runs
of four designs and four tiers each: enough to establish that the phenomenon exists and is transient, not
enough to estimate how often it happens. The re-scoped statement and the limits say so.

## Where I think the measurement is weak

1. **The counts are small** (2–3 events for `s1`, 1 for `s2`): existence, not a rate, and no tier law.
2. **The YB arm measured four tiers (10–10,000 seats) while the PostgreSQL control measured five
   (10–100,000)**, because that is how the study configures tiers per engine. "Held everywhere" means
   everywhere the PostgreSQL arm ran.
3. **The classification is the harness's own**: `transient` is decided by re-issuing the statement inside the
   refusing transaction (AM-01), so this replicates the mechanism and its detection, not an independent
   definition of a valid hold.
4. **Single node, one host, one day**, scale `small`; and the configured 2 GiB machine limit is not enforced on
   either arm (containers see the 15.4 GiB WSL VM), so the declared envelope is the documented sizing defect
   for every run here.
5. **My own `N3` was mis-specified**, which the review should treat as a protocol defect of mine and not as
   evidence about the claim.

## What I would measure next

1. **`v2-09`'s buyer × tier surface** and **`v2-08`'s PostgreSQL refund and organiser-edit arms** remain the
   queue's next replications (ranked in the t2 survey).
2. **The refusal phenomenon at more tiers and more buyers**, if anyone wants a rate rather than existence —
   the current tranche deliberately does not provide one.
3. **`yb-cluster3` stays a coverage gap** on this machine, but single-node YB arms are now demonstrably
   feasible and should not be excluded by the 2 GiB note alone.

## Disagreements with other analyses

None. It reproduces the primary's finding and adds the classification detail plus a control that fires on both
engines.
