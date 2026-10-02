# Review — T1 Study 02 replication (`v2-07`)

**Reviewer:** deepseek-flash, Deep Code CLI, **HIGH**
**Reviewed at:** 2026-10-02 · **Result under review:** `attempt-20261002T002231Z-83ca64/RESULT.md` (mirrored from
`attempt-20261001T190207Z-b02e4a`), merged on `main` at `35877d8`, tags `study-02-t1-remaining-v1` and
`study-02/v2-07-answerability-replication`

## Verdict: **accept**. `v2-07` stays `run_replicated` with one precision correction made during this review.

**Independence note.** This session executed the task as a LOW worker and reviews it as HIGH: the capability
gate is satisfied, the corpus's preference for a different vendor reading the same data is not. Mitigation, as
with t1: every number below was re-derived from the committed results with a **separate script** written for
this review.

## 1. Re-derived independently

| quantity | worker's claim | my re-derivation | verdict |
|---|---|---|---|
| `r05` unanswerable, forward run | 14/15, `x1` the only answerable | **14/15**, `x1` | matches |
| `r05` unanswerable, reversed run | 14/15, `x1` | **14/15**, `x1` | matches |
| `r01`+`r03` `partial` | 14/15 both runs | **14/15** both runs | matches |
| storage, tables only | 33.4% / 34.8% (primary 35.2%) | **33.4% / 34.8%** (primary 35.2%) | matches |
| storage, all relations | 35.0% / 35.9% (primary 36.2%) | **35.0% / 35.9%** (primary 36.2%) | matches |
| storage, indexes only | 37.7% | **37.7% / 37.7%**, but 37.4% / 37.4% in the pre-check and 37.9% in the primary | **corrected below** |
| gates | 15/15 clean, 0 over/dup/mismatch, tagged, clean | **15/15** both runs, all four counters 0, both tagged | matches |

**The one correction I made.** The statement and one limit quoted the index-only delta as `37.7%`, a point,
when the five runs span **37.4–37.9%** (the two evidence runs happen to coincide at 37.7%). The claim's own
thesis is that a storage percentage is a band, so quoting a point in the same sentence undercut it. Both now
carry the band; validator and tests pass afterwards.

## 2. Decisions reviewed (the attempt's Decision Log)

| decision | judgement |
|---|---|
| Four runs instead of two; the `verify`-only pair kept as a **labelled pre-check** after `report_coverage` turned out to live in the read phase | **Accept, and the transparency is the right call.** The protocol was amended *before* the evidence was used, the pre-check is labelled everywhere it appears (protocol, analysis, registry note), and no claim rests on it for answerability. Deleting the pair would have hidden a real trap for the next worker. |
| Storage definition pinned in the protocol (all / tables / indexes) because the primary's "35%" is the tables-only figure | **Accept.** This is the second claim in a row where the primary's single number turned out to be one rule's answer (t1's `2.9x` was one query); the pattern is now explicit in the protocol, which is where it belongs. |
| **Promotion to `run_replicated` while the prediction record says `falsified`** (`M4`) | **Accept, with a note that matters beyond this task.** The falsified prediction was *my tolerance for the percentage*, and the review confirms the claim's substance (the count, and the existence and rule-dependence of the cost) reproduced in both orders. But the registry now has **two** `run_replicated` claims whose prediction verdict is `falsified` (this one and `v2-04`), and a reader scanning strengths against verdicts will read that as a contradiction. The registry has no field for "the prediction was wrong, the claim was right"; the prediction's `text` carries it. **Suggested, not required:** the next registry revision could label a falsified prediction with which side it falsified — the claim's thesis or the author's band. |
| Book prose `+33–35%` re-scoped to `+33–38% storage — about 35% by table size, 38% counting indexes too` | **Accept.** Verified line by line against the statement's rules: the tables band 33.4–35.2% and the index band 37.4–37.9% are both inside it, and no value the claim says did not reproduce is quoted. |
| Two evidence runs at two commits (`fcab38dc`, `d0b4747a`) | **Accept.** I re-checked the diff: one file, the protocol markdown. No harness or SQL difference, so the runs are comparable. |

## 3. What I checked beyond the numbers

- **Scope:** the merge touched `v2-07` only — no claim added, removed, or otherwise re-scoped; no other
  analyst's file, report or analysis edited; the survey is a plan inside the protocol, not a claim.
- **Anchors:** all seven targets resolve (report present, analysis present, run tag present for the four new
  runs and the three earlier ones).
- **Gates:** `validate-v6` 31 claims / 0 errors; `go test ./tools/evidence/...` green; the four reports
  regenerated so the analysis is indexed; book rebuild clean at 83 pages with the evidence digest matching the
  registry; queue audit 0 errors.
- **The four new runs' own gates:** 15/15 clean `verify` per run, 0 over-booked, 0 duplicate seats, 0 ledger
  mismatches, clean trees, runner-created tags.

## 4. Where this evidence is still weak (the claim's limits, not just my notes)

1. **The audit is the harness's own**: this is a reproducibility result, not an independent audit of the rule.
2. **No timing**: the read phase ran only because `report_coverage` is emitted inside it; its figures are not
   used, and nothing here re-opens the ledger's load or read cost.
3. **Scale `small`, PostgreSQL only, one host, one day**: `M4`'s 0.9 pp movement is within-day movement.
4. **The pre-check's phases differ from the evidence pair's**, so its storage figures corroborate by a
   slightly different path.

## 5. Follow-ups

1. **The survey's own ranking is the queue's next work**: `v2-08`'s PostgreSQL refund and organiser-edit arms
   are the highest-value un-replicated figures left in this study, then `v2-09`'s buyer × tier surface. The
   `yb-cluster3` and `v2-11` YugabyteDB arms do not fit the owner's fixed machine and stay coverage gaps.
2. **Registry ergonomics, optional**: distinguish a falsified *claim* prediction from a falsified *author band*
   (D-note above); two claims now show the confusing combination.
3. **Queue hygiene**: the submit on this task raced a compare-and-swap and left `live_ref_drift` for
   `queue audit`; it clears on completion. Worth a `recover`-adjacent repair command if it recurs.

## 6. Integration

`main` at `35877d8` already contains the branch (`study-02/t1-remaining`), the required tag
`study-02-t1-remaining-v1` exists on that state, and the milestone tag
`study-02/v2-07-answerability-replication` is collision-free by construction — the naming rule recorded after
t1 was applied.
