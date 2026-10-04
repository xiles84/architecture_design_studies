# Progress — task-20261004T100822Z-select-order-correction

Attempt `attempt-20261004T102305Z-cd5604`, branch `study-01/select-order-correction`,
LOW executor (deepseek-flash / opencode CLI). Milestone tag:
`study-01/select-v7-order-controls`.

Source: independent HIGH review
`docs/select-evidence-review-20261004/LOW-CORRECTION-HANDOFF.md`.

## Plan

1. Recon and premeasurement plan, corrected environment page, Decision Log.
2. Seeded per-run order + recorded sequence in `selectbench` (`PlanOrder`,
   `OrderSeedFor`, `RunOrder`, `RunSpecsInterleaved`).
3. Interleaved sample rounds; actual-first position control; unit tests.
4. Wire access, aggregate and existence modes; access index arm as its own phase.
5. >=3 whole runs per engine/family (18 runs), 5 trials, under `ads-run-lock`,
   one run per engine/family with the position control.
6. Primary-input manifest (per-file + corpus digests) and correction reports
   with run ids on ratio rows.
7. Signed correction analysis; RESULT.md; CONTEXT/LESSONS; submit for HIGH.

## Decision Log

1. **Order seed derived from (dataset seed, run id).** Deterministic and
   reproducible per whole run while differing between runs. Confidence High.
   Recorded in the run's `order` block with the actual axis sequence.
2. **Interleaved rounds, not consecutive samples.** Each cell is verified and
   planned once, then one sample of every cell is taken per round. This changes
   only the sampling schedule; a late mismatch still invalidates the cell.
   Confidence High. Alternative: re-verify each round (more queries, no benefit).
3. **Position control repeats the actual first cell.** The randomized order's
   first cell is cloned with a `;position=last` axis, so no hard-coded cell is
   ever called the randomized first. Confidence High.
4. **Existence position control stays inside the fanout=1 phase**, because each
   fanout level is a different fixture and the repeat must share the first cell's
   fixture. Recorded per level. Confidence High.
5. **Access expression-index arm is a separate phase with its own seed and
   recorded order.** It is a labelled physical arm and must not be pooled with
   the default matrix. Confidence High.
6. **Matrix, fixtures, predicates, levels and N are unchanged.** Only run ids,
   schedule and provenance reporting change; original runs/tags/analyses are not
   touched. Confidence High. No question/oracle/prediction change, so no
   escalation.
7. **Corrected environment page.** The check's recorded guest was 8 CPUs /
   ~15.36 GiB; the earlier ~2 GiB-guest claim is stale and is superseded. The
   cluster per-node/equal-total arm stays an explicit gap. Confidence High.
8. **Input manifest with per-file and corpus digests**, excluding pilots, failed
   cells and old fixed-order runs. Confidence High.
9. **Uncertainty-qualified judgment.** Paired ratios are read with within-run
   spread, across-run variability and the position drift together; mixed signs
   or effects smaller than the spread stay unresolved. Confidence High.

## Limits

Single-node only; no cluster arm; one sequential reader; five samples; YB
RPC/storage counters not parsed.
