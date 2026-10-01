# Result — T1: replicate Study 01's remaining read claims (`v2-04`, `v2-05`)

**Task:** `task-20261001T133759Z-t1-study01-replication` · **claim:** `claim-2a0ed23c8147525f` (epoch 1)
**Attempt:** `attempt-20261001T141449Z-2926d2` · **worker:** deepseek-flash, Deep Code CLI, LOW executor

**Outcome: delivered.** Both claims were replicated run-level on their PostgreSQL arms and promoted to
`run_replicated` with re-scoped statements; a runner defect that would have silently truncated every future
matrix was found and fixed on the way.

## What was delivered

| artifact | where |
|---|---|
| pre-registered protocol (L1–L6, required-comparison coverage) | `studies/01-charity-tree/T1-RECENCY-REPLICATION.md` (commit `49dde3b`) |
| two whole runs, 18 cells each, AB/BA by trial, tagged by the runner | `run/01-charity-tree/20261001T144451Z-t1`, `…145914Z-t1` |
| signed analysis with the verdicts | `studies/01-charity-tree/reports/analyses/20261001T-t1-recency-replication--deepseek-flash--2026-10-01.md` |
| registry: both claims re-scoped and promoted, two anchors each | `book/evidence/v6/claims.json` (`dec1c62`) |
| book prose checked (no recency figure quoted) and the study README re-scoped | `studies/01-charity-tree/README.md` |
| context, lessons and the supersession ledger | `CONTEXT.md`, `LESSONS_LEARNED.md`, `book/evidence/v6/SUPERSESSION_LEDGER.md` |
| merged and tagged | `main` at `88a9c04`, tag `study-01/t1-remaining` |

## Results

- **`v2-05` confirmed**: `d3_flattened_fk` beats `d2_normalized_indexed` on the core read score in all four
  fresh loads across two whole runs — 1.59x, 1.72x, 1.64x, 1.84x — against the primary's 2.47x.
- **`v2-04` replicated as a query-specific result**: the rollup index beats the flag on q13 (1.9–2.3x) and q15
  trailing (2.6–3.1x), ties it on q14 and is marginally slower on q16 (0.85–0.94x); the group geomean is
  1.43–1.44x here against 1.52x recomputed from the primary, whose quoted 2.9x is the q13-trailing comparison.
  The statement now says that, and the 7.9x YugabyteDB half is named as unreplicated.

## Decisions logged

1. **Scoped the runner instead of the study**: added optional `--designs`/`--topologies` filters to
   `run-enhancements.sh` (unset = the published matrix, unchanged; the manifest records the filter) so two
   claims could be replicated in 28 minutes rather than the full matrix's hours.
2. **Fixed the runner's loop stdin on fd 3** rather than only redirecting one command: proven with
   `podman --version` that *every* podman call consumes stdin, so `</dev/null` on the cell alone was not
   enough. Confidence: high.
3. **Kept the two one-cell diagnostic runs** (`20261001T142606Z-t1`, `…142837Z-t1`) as evidence of the defect
   instead of deleting them; they are labelled and not cited as replication evidence.
4. **Promoted `v2-04` with `verdict: falsified` on its prediction record**, because `L4`'s band was mine and
   wrong; the claim's ordering replicated, and the record says both things.

## Review notes for HIGH

- `L4` (my 1.5x floor) and `L5` (`d1` slowest, a design the scoped plan never contained) were failures of my
  protocol, not of the claims; the analysis says so, and a corrected `L5` needs `d1` added to a future plan.
- The branch is **already merged into `main`** (`88a9c04`) and tagged, per the standing completion rule; the
  claim's worktree `.worktrees/t1-study01` is left in place for the review rather than pruned.
- No other claim was touched, no existing analysis was edited, and the book needed no prose change.
