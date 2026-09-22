# Progress — Study 04 v2 completion protocol

**Task:** `task-20260922T025912Z-study04-v2-protocol`
**Attempt:** `attempt-20260922T194852Z-d22bf2` (claim `claim-74bcd738c331bd2e`, epoch 1)
**Capability/role:** HIGH planner (model `deepseek-flash`)
**Required tag:** `study-04/v2-handoff`

Planning only: no database started, no benchmark run, no measurement.

## Deliverables

- `studies/04-configuration-portal/HANDOFF-AMENDMENT-01-V2.md` — decision-complete protocol.
- A pointer from `HANDOFF.md` to the amendment.
- Three newly published execution tasks:
  `task-20260922T195241Z-study04-v2-controls`,
  `task-20260922T195241Z-study04-v2-cardinality-cadence`,
  `task-20260922T195241Z-study04-v2-remaining-designs`.
- `CONTEXT.md` open gap 3 updated.

## Decision Log

1. **Amend, do not rewrite.** The v1 handoff stays valid; the amendment adds controls and phases and
   is what a v2 runner reads. This keeps the v1 designs and dataset as the fixed baseline.
   Confidence: High.
2. **The three review controls are blocking, not advisory.** 1.76x may be restated only after the
   `-retries 1` control; the ~1.7x read floor must be measured before any design ratio is quoted.
   Confidence: High.
3. **Three execution tasks, split by measurement phase.** Controls first (they gate the others),
   then cardinality/cadence, then remaining designs + equal-total + YugabyteDB/placement. Each takes
   the benchmark lock and produces its own report and signed analysis. Confidence: High.
4. **Equal-total resources are a separately labelled arm**, never pooled with the per-node baseline;
   placement requires a recorded tablet/leader distribution or it is a gap. Confidence: High.
5. **Repeats default to ≥3 fresh-load trials per cell**, or a recorded randomized order with the
   limitation stated. Confidence: High.

## Residual risks

- The amendment fixes tier values (1…500) and writer counts (1/4/16/64); a pilot may show the
  crossover lies outside them, in which case the task extends rather than reports a flat line.
- The eight unbuilt designs may not all be buildable; each must be recorded as a gap with a reason.
