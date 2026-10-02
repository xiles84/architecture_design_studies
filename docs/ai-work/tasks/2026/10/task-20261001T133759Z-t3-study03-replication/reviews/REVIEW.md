# Review — T1 Study 03 replication (`v2-10`)

**Reviewer:** deepseek-flash, Deep Code CLI, **HIGH**
**Reviewed at:** 2026-10-02 · **Result under review:** `attempt-20261002T145931Z-d97f36/RESULT.md`, merged on
`main` at `8aa5d6b`, tags `study-03-t1-replication-v1` and `study-03/v2-10-refusal-replication`

## Verdict: **accept**. No correction needed; the re-scoped statement matches the data to the event.

**Independence note.** Same session authored and reviewed (capability gate satisfied; vendor independence not).
Mitigation as before: every counter below was re-derived from the committed results by a separate script.

## 1. Re-derived independently

| quantity | claim | my re-derivation | verdict |
|---|---|---|---|
| YB forward, correct designs: refusals (audit) | 2 (`s1`) + 1 (`s2`), all transient | **3 total, 3 transient** (`s1` 2, `s2` 1) | matches |
| YB reversed: refusals | 3 (`s1`) + 1 (`s2`) | **4 total, 4 transient** (`s1` 3 at 10k, `s2` 1 at 1k) | matches |
| retry variant | 0 refusals, 5 retries all successful | **0 refusals, 5/5 successes** | matches |
| PG control, correct designs, both runs | zero refusals, retries, failures | **all counters 0 at all five tiers** | matches |
| negative control | fires on both engines | **1,120–1,144 refusals and 2,090–2,119 thefts (YB); 8,408–9,420 and 16,142–18,095 (PG)** | matches |
| provenance | clean, tagged | all four `repo_dirty=false`, runner-tagged | matches |

No numeric correction was needed — unlike `t2`, where the index-only figure was a point in a band's clothing.

## 2. Decisions reviewed

| decision | judgement |
|---|---|
| Ran the feasibility check as a **real attempt** rather than reasoning from the recorded blocker | **Accept, and it is the most valuable thing this tranche produced.** Five minutes of `yb-single` at `tiny` scale converted a probable "coverage gap" into a replication and falsified a standing assumption (see below). |
| **Refusals read from the audit, not the client's `errors`** | **Accept.** My first pass treated `s0`'s thousands of client errors as a defect; the audit showed `s0` is the declared negative control and its `thefts` are the control firing. The lesson is recorded. |
| Four designs (claim's pair + `s0` + `s2`) | **Accept.** Testing "correct designs refused" on one design would have been thin; including the control in the same cells is what makes PostgreSQL's zeros meaningful. |
| Promotion with `verdict: falsified` (`N3`, `N4`) | **Accept.** This is the third claim whose prediction record reads `falsified` while the claim is `run_replicated`; the review of `t2` already flagged that the registry cannot distinguish a falsified *author band* from a falsified *claim*. Here both failures are clearly the author's (`N3` forgot the negative control must refuse on PostgreSQL; `N4`'s tier rule broke in the reversed order), and the prediction text says so. |
| Recording the small counts as existence, not a rate | **Accept.** 2–4 events across four designs and four tiers is enough to establish the phenomenon and not enough to estimate it. |

## 3. The machine-record correction, reviewed

The finding that the **2 GiB podman blocker is a multi-node limit** is supported by this task's own evidence:
`yb-single` started, loaded, ran its race phase and completed with a consistent audit, twice in the pre-check
and four times in the evidence runs. That is direct evidence, not inference. The claim's own limit and
`v2-11`'s coverage gap are unchanged, and `CONTEXT.md` now states the refined scope next to the original
blocker rather than replacing it.

**Consequence for the wider programme:** the earlier survey results that ranked "YB arms do not fit" (in the
t2 protocol and in `CONTEXT.md`) were **too broad for single-node YB arms**. `yb-cluster3` genuinely does not
fit; `yb-single` does. The next planning session should re-rank the YB items with that distinction.

## 4. Where this evidence is weak (also in the claim's limits)

1. Counts are 2–4 events: existence, not a rate, and no tier law (`s2` refused at 10k in one order and 1k in
   the other).
2. The YB arm measured four tiers, the PostgreSQL control five — "held everywhere" means where the PG arm ran.
3. The classification is the harness's own (AM-01 re-issue), so this replicates the mechanism and its
   detection, not an independent definition of a valid hold.
4. Single node, single host, one day, scale `small`; the declared envelope still carries the sizing defect.
5. Two of the tranche's own predictions were mis-specified — a protocol defect of the author, recorded, and it
   does not touch the claim.

## 5. Follow-ups

1. **Re-rank the YB items** in the t2 survey now that `yb-single` is known feasible.
2. **`t4` and `t5`** remain ready; the new `select-*` programme (another session) is also queued for LOW.
3. **Registry ergonomics** (carried from the t2 review): a falsified prediction should say which side it
   falsified. Three claims now show the combination.

## 6. Integration

`main` at `8aa5d6b` contains the branch `study-03/t1-replication`; the required tag
`study-03-t1-replication-v1` and the milestone tag `study-03/v2-10-refusal-replication` both mark that state,
and neither collides with a branch name.
