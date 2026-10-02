---
analysis_id: 20261002T-t1-answerability-replication--deepseek-flash--2026-10-02
run_id: 20261001T191704Z-t1 (read pass 1), 20261001T194638Z-t1 (read pass 2)
environment: host-zenbook-ux5406sa
analyst: deepseek-flash
analyst_kind: ai
analyst_version: "deepseek-flash (DeepSeek), LOW worker session, Deep Code CLI; T1 replication of v2-07 (answerability and the ledger's storage cost)."
analyzed_at: 2026-10-02
inputs_digest: 02a64483f9516a93 (read pass 1), f879f1cb732ebe49 (read pass 2)
inputs: 20261001T191704Z-t1@02a64483f9516a93, 20261001T194638Z-t1@f879f1cb732ebe49, 20261001T191132Z-t1@d5d39748f0846f1c, 20261001T191256Z-t1@dedd99c0ae0863ba
repo_commit: d0b4747a
supersedes:
status: current
headline: v2-07's schema half replicates exactly - r05 unanswerable in 14 of 15 designs and the ledger the only design that answers it, in two whole runs in both design orders - and its storage half replicates as a band rather than a point, because the relation sizes themselves move about a percentage point between runs.
---

# Analysis — T1 replication of `v2-07` — deepseek-flash

## TL;DR

- **The answerability half replicates exactly.** `r05` (the refunds report) is `unanswerable` in **14 of 15**
  designs and `x1_cas_ledger` is the **only** design that answers it, in both whole runs and both design
  orders; `r01` (sales/revenue) and `r03` (last purchase) are `partial` in the same 14, and `r02`/`r04` are
  `answerable` in all 15. The claim's *count* is exact and needs no re-scoping.
- **The storage half replicates as a band, and the band is the honest shape.** The ledger's delta over
  `p3_precreated_cas` is **+33.4% / +34.8% (tables only)**, **+35.0% / +35.9% (all relations)** and
  **+37.7% (indexes only)** in the two evidence runs, against the primary's **35.2% / 36.2%**. So "about
  33–35%" is right for the tables-only rule the primary used, and a reader who counts indexes too gets about
  38% — **the percentage depends on which relations are counted, and the claim should say which**.
- **`M4` was falsified, and what it shows is worth keeping.** I predicted the delta would be reproducible to
  within 1% between runs; it moved **2.6% relative** (35.0% → 35.9%, i.e. 0.9 percentage points) because the
  relation sizes themselves differ run to run (`p3` 76,251,136 against 76,054,528 bytes, 0.26%; `x1`
  102,907,904 against 103,333,888, 0.41%). PostgreSQL relation size includes free space left by the load, so
  **a single-run storage percentage is a reading of a page layout, not a constant of the schema** — which is
  exactly why the claim's own word is "about".
- **Correctness was green throughout**: 15/15 cells with a clean `verify` in each run, 0 over-booked, 0
  duplicate seats, 0 ledger mismatches, both trees clean and both runs tagged.
- **The claim is registered `run_replicated`** for its schema half, with its statement naming the rule the
  percentage is measured on and a limit keeping its remaining single-run halves visible.

## What I was looking at

Two whole runs of all fifteen designs on `pg-single` at scale `small`, the study's own runner and published
read footprint, in both design orders (`20261001T191704Z-t1` forward, `20261001T194638Z-t1` reversed), 15 cells
each. A **pre-check pair** (`20261001T191132Z-t1`, `20261001T191256Z-t1`) ran the same design set with
`--phases verify` only; it carries the storage and gate evidence (tables-only 34.1% / 33.9%, all-relations
35.3% / 35.2%) but no answerability data, because `report_coverage` is emitted inside the read phase — a
protocol defect I found and recorded by amending the protocol before using the evidence.

The two evidence runs sit at two different commits (`fcab38dc` and `d0b4747a`); `git diff` between them shows
**one file, the protocol markdown**, and no harness or SQL change.

| prediction | verdict |
|---|---|
| M1 `r05` unanswerable in exactly 14/15, `x1` the only answerable | **confirmed** in both runs |
| M2 `r01`/`r03` partial in the same 14, `r02`/`r04` answerable in all 15 | **confirmed** in both runs |
| M3 every storage rule inside +32…+38% | **confirmed** (33.4–37.7% across both runs) |
| M4 run-to-run difference under 1% relative | **falsified** — 2.6% relative (0.9 pp) |
| M5 the ledger largest, `p*` family within 1% | **confirmed** (`x1` largest; p-family spread 0.52% and 0.18%) |
| M6 gates green, clean tagged trees | **confirmed** (15/15 in each run) |

## Conclusions

**The claim's substance survives, and its count half is stronger than its strength said.** Answerability is a
schema fact and the schema has not changed, so this is a reproduction of a deterministic fact rather than a
new measurement of a noisy quantity — which is precisely what the claim's own limit says. Two whole runs in
opposite design orders make it `run_replicated` in the sense the registry means: the fact is the same fact in
more than one run.

**Its storage figure should name its rule.** The primary's analysis said "35% storage" — that is the
**tables-only** figure (35.2% in its own data); counting every relation gives 36.2%, and counting only indexes
37.7%. My two runs land at 33.4–34.8% (tables) and 35.0–35.9% (all). The re-scoped statement therefore quotes
the band per rule instead of one number, and states that the percentage moves about a percentage point
between runs.

**The falsified prediction is a finding about the measurement, not about the design.** `M4` assumed relation
size is a schema constant; it is a schema fact *plus* whatever free space the load left behind. A future
tranche that wants a tighter storage number would have to compare `pg_relation_size` on a vacuumed database,
and that is a different measurement from the one the claim is about.

## Where I think the measurement is weak

1. **The audit is the harness's own.** `report_coverage` is computed by the same code that defines the
   schemas, so this replication tests reproducibility, not the independence of the audit's rule. The
   September footprint-2 read work verified the reports' *answers*; this tranche does not re-open that.
2. **No timing here.** Nothing in this analysis speaks to the ledger's load or read cost; the read phase ran
   only because `report_coverage` is emitted inside it, and its figures are not used.
3. **Scale `small`, PostgreSQL only.** The dataset is 690 events; the sizes are what the study's primary
   measured, but a larger dataset could change the *percentage* (fixed per-table overheads matter less).
4. **One day, one host.** Two runs 29 minutes apart share the session's state; `M4`'s 0.9 pp movement is
   within-day movement and says nothing about week-to-week variation.
5. **The pre-check pair's phases differ** from the evidence pair's, so its storage figures corroborate by a
   slightly different path (same load, same `stats`); I report them as a pre-check and not as independent
   runs.

## What I would measure next

1. **`v2-08`'s PostgreSQL refund and organiser-edit arms** — the highest-value un-replicated figures left in
   this study; the survey in `T1-ANSWERABILITY-REPLICATION.md` ranks it first.
2. **`v2-09`'s buyer × tier surface** at the remaining tiers.
3. **A vacuumed storage comparison**, if anyone wants a storage figure tighter than a percentage point.
4. The `yb-cluster3` arms remain **unmeasurable on the owner's fixed 2 GiB machine** and are recorded as a
   coverage gap rather than a plan.

## Disagreements with other analyses

None. This reproduces the primary analysis's answerability finding exactly and corrects only the *shape* of
its storage figure — from "35%" to a band whose value depends on the rule — which its own text already hinted
at by writing "about".
