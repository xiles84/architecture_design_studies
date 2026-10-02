# Result — T1: replicate Study 03's `v2-10` (guarded confirmations refused on YugabyteDB)

**Task:** `task-20261001T133759Z-t3-study03-replication` · **claim:** `claim-ee93f03b2c2d6416`
**Attempt:** `attempt-20261002T145931Z-d97f36` · **worker:** deepseek-flash, Deep Code CLI, HIGH executor

**Outcome: delivered.** The mandated feasibility check passed, the claim replicated on both engines, and the
check also corrected a standing assumption about what this machine can run.

## What was delivered

| artefact | where |
|---|---|
| pre-registered protocol with the feasibility result and N1–N6 | `studies/03-reserved-seating/T1-REFUSAL-REPLICATION.md` (`6dbeac27`) |
| feasibility pre-check (`tiny`, `yb-single`) | `run/03-reserved-seating/20261002T150337Z-feas` |
| four whole runs: YB and PG, each in both design orders | `…20261002T150903Z-t1yb`, `…152604Z-t1yb`, `…154340Z-t1pg`, `…155134Z-t1pg` |
| signed analysis with the six verdicts | `studies/03-reserved-seating/reports/analyses/20261002T-t1-refusal-replication--deepseek-flash--2026-10-02.md` |
| registry: `v2-10` promoted, statement re-scoped, four anchors | `book/evidence/v6/claims.json` |
| book prose sharpened with the transient finding; machine blocker scope refined | `book/concepts/expiry-and-clock-authority.typ`, `CONTEXT.md` |
| ledger and two lessons | `book/evidence/v6/SUPERSESSION_LEDGER.md`, `LESSONS_LEARNED.md` |
| merged and tagged | `main` at `8aa5d6b`, tags `study-03-t1-replication-v1` and `study-03/v2-10-refusal-replication` |

## Results

- **The claim replicates on YugabyteDB**: `s1_conditional_update` refused valid holds **2 times** in the
  forward run and **3 times** in the reversed run, `s2_lock_then_update` **once in each**, and **every refusal
  was classified `transient`** — the identical statement re-issued inside the refusing transaction matched
  every seat.
- **The retry fix holds**: `s1r_confirm_retry` refused nothing at any tier in either run, and its 5 retries in
  the forward run all succeeded.
- **PostgreSQL is a real control**: all three correct designs recorded zero refusals, zero retries and zero
  correctness failures at all five tiers (10–100,000 seats) in both runs.
- **The negative control fires on both engines**: `s0_check_then_hold_rc` refused 242–351 per tier on
  YugabyteDB and 555–3,526 on PostgreSQL, stealing 474–599 and 1,168–6,231 holds respectively — so PostgreSQL
  refused *more* than YugabyteDB for a wrong design, which is what makes its silence for the correct ones
  evidence rather than a quieter instrument.
- **The feasibility check changed a standing assumption**: the 2 GiB podman blocker is a **multi-node** limit.
  A single-node YB arm starts and completes cleanly; `yb-cluster3` remains unmeasurable.

## Decisions logged

1. **Ran the feasibility check as a real attempt** (`tiny`, `yb-single`, phases `verify,race`) rather than
   reasoning from the recorded blocker — it took five minutes and reproduced the claim's mechanism, which
   turned a probable coverage-gap note into a real replication. Confidence: high.
2. **Selected four designs** (the claim's `s1`/`s1r` pair plus `s0` and `s2`) so "correct designs refused"
   would be tested on more than one design and the study's negative control would run in the same cells.
   Confidence: high.
3. **Read refusals from the audit, not the client's `errors`.** My first pass treated `s0`'s thousands of
   client errors as a defect; the audit showed `s0` is the **declared negative control** and that its
   `thefts` are the control firing. Confidence: high — and it changed the analysis.
4. **Promoted with `verdict: falsified`** because `N3` and `N4` both failed as written: `N3` forgot the
   negative control must refuse on PostgreSQL too, and `N4`'s "largest tier" rule failed in the reversed run
   (`s2` refused once at the 1,000-seat tier while its 10,000-seat cell refused nothing). The claim's
   substance replicated; my two bands did not.
5. **Recorded the small counts honestly**: 2, 3 and 1 events are existence, not a rate, and the statement and
   limits say so.

## Review notes for HIGH

- Two protocol wording defects of mine are recorded in the analysis and in the prediction record (`N3`, `N4`),
  so the review can see them rather than infer them.
- The YB arm measured four tiers and the PostgreSQL control five, because that is how the study configures
  tiers per engine — stated as a limit.
- Branch already merged into `main` (`8aa5d6b`) and tagged; the claim's worktree is left in place.
