# Review — T1 Study 01 replication (`v2-04`, `v2-05`)

**Reviewer:** deepseek-flash, Deep Code CLI, **HIGH** (session capability changed from LOW after the task was
submitted — see the independence note below)
**Reviewed at:** 2026-10-01 · **Result under review:** `attempt-20261001T141449Z-2926d2/RESULT.md`,
branch commit `88a9c04` on `main`, tag `study-01/t1-remaining`

## Verdict: accept. Both promotions stand; one defect found in the *review* method, not the work.

## Independence note (stated up front)

This session authored the work as a LOW worker and reviews it as HIGH. That satisfies the queue's capability
gate but **not** the corpus's preference for a different vendor reading the same data. The mitigation is that
every headline number below was re-derived from the committed results with a **separate script** written for
this review, not by re-reading the worker's analysis.

## 1. The numbers, re-derived independently

`/tmp/high_t1_verify.py` (reviewer's own) recomputes both claims from
`studies/01-charity-tree/results/2026…-t1/pg-single-*/*.json`:

| comparison | worker's claim | my re-derivation | verdict |
|---|---|---|---|
| `d3`/`d2`, core read score, four fresh loads | 1.59x, 1.72x, 1.64x, 1.84x | **1.59x, 1.72x, 1.64x, 1.84x** | matches to the digit |
| `d22`/`d20`, recency group geomean per run | 1.44x, 1.43x | **1.44x, 1.43x** | matches |
| `d22`/`d20`, q13 trailing (the figure the claim quotes) | 1.92x, 2.33x | **1.92x, 2.33x** | matches |
| q16 (behind the flag) | 0.85–0.94x | **0.85x–0.94x** | matches |

**A defect in my first attempt at this check, recorded because it is instructive:** my first script matched
core queries by prefix (`q08…`) and therefore also swept in the cell's `q08_total_donated_charity@largest` and
`q12_charity_recent_feed@smallest` regime variants, producing 2.03–2.52x and an apparent mismatch. The worker's
exact-name match was correct and my shortcut was wrong. **The reviewer's first result was the wrong one; the
lesson is that a verification harness needs verifying too**, and the run's own query list (visible in
`cells.tsv`) is the authority for which queries form the score.

## 2. Decisions reviewed (from the attempt's Decision Log)

| decision | judgement |
|---|---|
| Scoped `--designs`/`--topologies` filters on `run-enhancements.sh` | **Accept.** Additive: both default empty, the unfiltered matrix is untouched, and the manifest records `design_filter`/`topology_filter` so a filtered run cannot masquerade as a full one. This is the right trade: a scoped replication is a labelled subset, not a new protocol. |
| The fd-3 loop fix (after `</dev/null` alone proved insufficient) | **Accept, and the most valuable output of the task.** The evidence is sound (`podman --version` ends a `while read … done < file` loop after one iteration) and the change is strictly corrective. |
| Keeping the two one-cell diagnostic runs in `results/` | **Accept.** They are labelled in the analysis and are the artefact that proves the defect existed. Deleting them would have destroyed the finding's evidence. |
| Promoting `v2-04` while its prediction record says `falsified` | **Accept with a note.** The falsified prediction is `L4` (my own 1.5x floor), and the claim's ordering replicated in every trial; the prediction text says so. Note for readers: a `run_replicated` claim with `verdict: falsified` looks contradictory at a glance, and the text has to be read. The alternative — recording `confirmed` because the claim survived — would have hidden that the band was wrong. |
| `L5` recorded as untestable (`d1` was never in the scoped plan) | **Accept.** Correctly reported as a protocol defect rather than quietly dropped. |

## 3. What I checked

- **Scope:** the merge touched exactly two claims (`v2-04`, `v2-05`), no claim added or removed; no other
  analyst's file, report or analysis edited.
- **Anchors:** all six anchor targets resolve — report file present, analysis path present (including the two
  moved to `outdated/` earlier and repointed), run tags `run/01-charity-tree/20261001T144451Z-t1` and
  `…145914Z-t1` present.
- **Gates:** `validate-v6` 31 claims / 0 errors; `go test ./tools/evidence/...` green; report regenerated so
  the analysis appears in its index; book rebuild clean with the evidence digest matching the registry.
- **Correctness gates inside the runs:** 18 cells per run, `failed processes 0`, no read or write error in any
  cell, both trees clean at `49dde3bd`, both runs tagged by the runner.

## 4. Where this evidence is still weak (carried into the claim limits, not just noted here)

1. PostgreSQL only: `v2-04`'s quoted 7.9x YugabyteDB half is unreplicated, and the machine is fixed at 2 GiB
   by the owner's instruction.
2. Five designs of twenty-five; `d1` absent, so `L5` remains open.
3. Two runs on one host on one day: absolute levels are day-confounded; only ratios and the per-query pattern
   are used, which the statements respect.
4. One measurement per order per run (`--trials 2` gives AB/BA, not repeats).

## 5. Required follow-ups created by this review

1. **Nothing blocking.** The two claims may stay promoted as written.
2. **For the next Study 01 tranche** (not this task): include `d1` so `L5` can be answered, and note the
   `q08`/`q12` regime variants in the protocol so the read score's query set is explicit for any future
   re-derivation.
3. **For the shared platform:** audit the other study runners for the same loop/stdin pattern. Only
   `run-enhancements.sh` reads a loop from a file today, but the fix belongs in the habit, not one script.

## 6. Integration

The task's own completion rule put the work on `main` before submission: `88a9c04` is reachable and
`study-01/t1-remaining` is tagged on the integrated state (the branch itself was renamed to
`study-01/t1b-work` to stop the branch/tag ref ambiguity that naming had created — a worktree-local rename, no
history rewrite, and I checked the tag resolves to the merge commit).
