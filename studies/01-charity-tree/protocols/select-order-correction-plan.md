# Order-correction plan — randomized/interleaved reruns of the three SELECT families

LOW worker plan for `task-20261004T100822Z-select-order-correction`, implementing
the independent HIGH review's
[`LOW-CORRECTION-HANDOFF.md`](../../../docs/select-evidence-review-20261004/LOW-CORRECTION-HANDOFF.md).
It changes **when** cells are measured, never what is asked or how answers are
checked.

## 1. Defect being corrected

Every earlier matrix executed its cells in one fixed order and took each cell's
samples consecutively, so a session/position drift lands on whichever cells
happen to run late and is confounded with the spelling. The check recorded one
distinct sequence per family/engine.

## 2. Correction

- **Seeded per-run order.** Each whole run permutes its cells with a seed derived
  from the dataset seed and the run id (`OrderSeedFor`). The seed and the actual
  axis sequence are recorded in the run JSON and report.
- **Interleaved rounds.** Each cell is bound, oracle-verified and planned once
  (in the randomized order); then one sample of every valid cell is taken per
  round, `trials` rounds in total. Within-run drift is spread across cells.
- **Actual-first position control.** The optional repeated-design cell repeats the
  *first cell of that run's randomized order*, not a hard-coded cell. For the
  existence family it repeats within the fanout=1 phase, where the fixture is the
  same. The access expression-index arm is a separate labelled phase with its own
  recorded order.
- Correctness is unchanged: every timed query is drained and verified against the
  same oracle; a late mismatch still invalidates the cell; all fixture and
  control gates are preserved.

## 3. Matrix held fixed

Same families, statements, seeds, fixtures, predicates, page size, N=1, the five
paginations/parent levels and the three fanout levels as the originals. No
question, answer, oracle, prediction or accepted matrix is changed. Fresh run
ids and producing tags; the original runs and tags are untouched.

## 4. Replication, provenance and reports

- ≥3 whole runs per engine/family, 5 trials, serially under `ads-run-lock`, one
  run per engine/family carrying the position control.
- An explicit primary-input manifest lists the exact `*.json` files with
  per-file SHA-256 and a corpus digest, plus the producing tag; pilots, failed
  cells and old fixed-order runs are excluded.
- Per-run reports now carry the recorded order; a correction report exposes the
  run ids on the ratio rows and judges paired ratios with within/across-run
  variability and position drift together. An effect smaller than the spread is
  `unresolved`.

## 5. Interpretation rules

Mixed signs, missing ordering evidence or effects smaller than uncertainty stay
unresolved. No universal equality, no client-capacity inference from zero CFS
throttling, no unmeasured RPC attribution, and depth 100 on PostgreSQL stays
unresolved pending the corrected runs. A signed correction analysis is written
beside the reports; the original signed analyses are not edited.

## 6. Gaps

Single-node only; no cluster per-node/equal-total arm (the recorded resources make
the budget calculable, but no cluster topology was measured); no placement, no
concurrent demand; five samples.
