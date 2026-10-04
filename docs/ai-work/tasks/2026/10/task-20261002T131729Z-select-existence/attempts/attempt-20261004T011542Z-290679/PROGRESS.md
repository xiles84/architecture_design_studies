# Progress — task-20261002T131729Z-select-existence

Attempt `attempt-20261004T011542Z-290679`, branch `study-01/select-existence`,
LOW executor (deepseek-flash / opencode CLI). Milestone tag:
`study-01/select-v3-existence`.

## Plan

1. Recon and premeasurement plan (this file, the bounded-matrix plan, the
   environment page) committed before any run.
2. Catalogue additions (`e12`, `e13`, `x07`), oracle reuse, controls, pure tests. **[done]**
3. Fanout-sweep runner (`BuildExistenceSpecs`), per-run and replicated reports. **[done]**
4. Harness `-cmd select-existence`, runner `--family existence`. **[done]**
5. PostgreSQL single-node: 3 whole runs (one position control); YugabyteDB
   single-node: 3 whole runs. **[done, all at 805b5b85; a first attempt failed
   the gate and is retained as a harness-defect record]**
6. Signed analysis with P6/P7 verdicts; RESULT.md; submit for HIGH review.
   **[done]**

## Checkpoints

- `78fdd9c9` plan/env/Decision Log; `f6aaecfd` catalogue/runner/report/controls.
- `805b5b85` applied `ExistenceShape` in the mode and added the per-level control
  test after the first run failed its gate on the raw fixture.
- Six primary runs valid (PG 21/22/21, YB 21/22/21), both controls fired 18/18,
  replicated report and signed analysis written.

## Decision Log

1. **Fanout is a fixture axis, not a query parameter.** Membership is compared at
   child fanout {1,10,100} with 10,000 child rows held constant, so one run loads
   three deterministic fixtures (10,000/1,000/100 parents). Confidence High. The
   alternative — a query parameter — cannot express fanout. Affects every cell
   and each level's fingerprint, which is recorded separately.
2. **Reuse the existing s06–s10 membership/absence statements instead of
   duplicating them.** They already implement EXISTS, IN, JOIN+DISTINCT,
   NOT EXISTS and LEFT-IS-NULL under `sel_membership`/`sel_absence` with oracles,
   and fanout is data. Confidence High. Alternative: new `e01…` copies (adds
   merge surface and a second SQL text for the same question).
3. **`IN`/`NOT IN` are paired at the NOT-NULL contract (`e12`, `e13`).** The
   valid `NOT IN` arm filters `customer_id IS NOT NULL`; the nullable case is
   `x02_not_in_nullable`, the labelled counterexample, and is never ranked.
   Confidence High. This is the protocol's `{NULL, NOT NULL}` axis.
4. **Added `x07_membership_multiplier`.** A plain child-row join returns one row
   per qualifying child, not one row per parent; it is the ordinary join the
   brief says is illustrated separately and must be detected against the
   membership set. Confidence High. Affects the control gate only.
5. **Absence `{NULL, NOT NULL}` is a contract axis, not a spelling axis.** The
   NULL arm cannot be a valid absence spelling (its answer is empty), so it is
   the control; the three ranked absence spellings all use the NOT-NULL
   subquery contract. Confidence High.
6. **Reuse the integrated shared runner (`CellSpec`/`RunSpecs`/`Axis`).** The
   two earlier families' reconciliation put one runner on `main`; this task is
   the first to build on it directly, with no new runner code. Confidence High.
7. **Position control repeats the first membership cell at the end of the phase**
   (opt-in), as in the other families. Confidence High.
8. **`ExistenceShape` moves the last parent's children to the first parent and
   gives the first parent two amount-400 orders.** Without it no parent has more
   than one qualifying child at fanout=1, so the x07 multiplier cannot fire, and
   no parent is childless, so the x02 nullable-`NOT IN` control cannot differ
   from the absence set. The total child count is unchanged and the shape is in
   the fingerprint. Confidence High. Affects the control gate at every level.
9. **One fingerprint per fanout level.** Each level is a different dataset, so
   the run's provenance table records the level's customer/order counts and
   fingerprint rather than a single global one. Confidence High.

## Limits already known

Single-node only; no 3-node/equal-total/placement arms; one sequential reader,
no writers; five samples; qualifying proportion, hit fraction and skew are not
swept here; YB RPC/storage counters not parsed.
