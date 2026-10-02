# Result — T1: replicate Study 02's `v2-07`, and survey what the machine can still measure

**Task:** `task-20261001T133759Z-t2-study02-replication` · **claim:** `claim-243943cc1e4478ab` (epoch 1)
**Attempt:** `attempt-20261001T190207Z-b02e4a` · **worker:** deepseek-flash, Deep Code CLI, LOW executor

**Outcome: delivered.** `v2-07` is `run_replicated` for its schema half, with its storage figure re-scoped to
the rule it is measured on; the survey ranks the remaining targets and records what the fixed 2 GiB machine
cannot carry.

## What was delivered

| artefact | where |
|---|---|
| pre-registered protocol (M1–M6, required-comparison coverage, the ranked survey) | `studies/02-ticket-booking/T1-ANSWERABILITY-REPLICATION.md` (`fcab38dc`, amended `d0b4747a`) |
| two evidence runs, 15 designs each, both design orders, tagged | `run/02-ticket-booking/20261001T191704Z-t1`, `…194638Z-t1` |
| two pre-check runs (storage and gates, no answerability) | `…191132Z-t1`, `…191256Z-t1` |
| signed analysis with the six verdicts | `studies/02-ticket-booking/reports/analyses/20261002T-t1-answerability-replication--deepseek-flash--2026-10-02.md` |
| registry: `v2-07` promoted, statement re-scoped, four anchors | `book/evidence/v6/claims.json` |
| book prose re-scoped | `book/chapters/how-to-choose.typ` (`+33–38% storage — about 35% by table size, 38% counting indexes too`) |
| ledger, context, lessons (two new lessons) | `book/evidence/v6/SUPERSESSION_LEDGER.md`, `CONTEXT.md`, `LESSONS_LEARNED.md` |
| merged and tagged | `main` at `35877d8`, tags `study-02-t1-remaining-v1` and `study-02/v2-07-answerability-replication` |

## Results

- **The answerability half reproduced exactly** in both whole runs and both design orders: `r05` unanswerable
  in **14 of 15** designs, `x1_cas_ledger` the only design that answers it; `r01`/`r03` `partial` in the same
  14; `r02`/`r04` `answerable` in all 15.
- **The storage half is a band that depends on the rule**: tables only **+33.4% / +34.8%** (primary 35.2%),
  every relation **+35.0% / +35.9%** (primary 36.2%), indexes alone **+37.7%**. The statement now names the
  rule, and the book's bullet no longer quotes one number.
- **Correctness green**: 15/15 clean `verify` cells per run, 0 over-booked, 0 duplicate seats, 0 ledger
  mismatches, both runs tagged from clean trees.
- **`M4` falsified**: the delta moved 2.6% relative (0.9 percentage points) between runs, because relation size
  carries the free space the load left behind.

## Decisions logged

1. **Four runs instead of two, and why.** The first pair used `--phases verify`; `report_coverage` turned out to
   be emitted inside the **read** phase (`harness/main.go`), so that pair carries storage and gates but no
   answerability. It is kept as a **labelled pre-check** rather than deleted, and the protocol was amended
   before the evidence was used. Confidence: high.
2. **The storage definition is pinned in the protocol before measuring** (all relations / tables only /
   indexes), because the primary's "35%" is the tables-only figure while every relation gives 36.2%.
   Confidence: high.
3. **`v2-07` promoted to `run_replicated` although `M4` was falsified**, because the falsified prediction was
   about *my* tolerance for the percentage, not about the answerability count or the existence of the cost —
   and the prediction record says both. Confidence: medium-high; the review should check this judgement.
4. **The book's `+33–35%` re-scoped rather than left**, since my runs show 38% by index size. Confidence: high.
5. **The two evidence runs sit at two commits** (`fcab38dc`, `d0b4747a`); the diff between them is
   `T1-ANSWERABILITY-REPLICATION.md` only, recorded in the analysis.

## Review notes for HIGH

- The survey in the protocol is a *plan*, not a measurement: it ranks `v2-08`'s PostgreSQL refund and
  organiser-edit arms first, then `v2-09`'s buyer×tier surface, and records the `yb-cluster3`/`v2-11`
  YugabyteDB arms as **not fitting** the owner's fixed 2 GiB machine.
- The answerability statuses are the harness's own audit of its own schema; this tranche tests reproducibility,
  not the audit's independence.
- Branch already merged into `main` (`35877d8`) and tagged; the claim's worktree is left in place.
