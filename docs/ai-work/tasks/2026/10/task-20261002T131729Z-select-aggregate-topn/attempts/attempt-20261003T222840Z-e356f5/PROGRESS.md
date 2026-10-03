# Progress — task-20261002T131729Z-select-aggregate-topn

Attempt `attempt-20261003T222840Z-e356f5`, branch `study-01/select-aggregate-topn`,
LOW executor (deepseek-flash / opencode CLI). Milestone tag:
`study-01/select-v4-aggregate-topn`.

## Plan

1. Recon and premeasurement plan (this file, the bounded-matrix plan and the
   environment page) committed before any run. **[done]**
2. Catalogue + oracles + controls + pure tests. **[done, 57c10e3a]**
3. Spec-based runner, P8/P9 matrix, measurement report. **[done, d13ce296]**
4. Harness mode, 10,000-parent sweep scale, runner family. **[done, 00f2aa3d]**
5. PostgreSQL single-node: 3 whole runs, 5 trials, one position control.
   **[done, 08a3c553 then 33b76901; a first attempt panicked in the harness]**
6. YugabyteDB single-node: 3 whole runs. **[done, 33b76901]**
7. Per-run and per-engine replicated reports. **[done]**
8. Signed analysis with P8/P9 verdicts; RESULT.md. **[done; BOOK.md integration
   belongs to the later select-book task]**
9. Submit for HIGH review. **[done]**

## Checkpoints

- `57c10e3a` catalogue/oracles/controls; `d13ce296` runner/matrix/report;
  `00f2aa3d` harness mode/scale/runner family; `ad048e0c` plan/env/Decision Log.
- `08a3c553` fixed a nil-predicate panic in the run-level control sweep found by
  the first (failed, exit 2) PostgreSQL attempt.
- `33b76901` added the opt-in position control.
- Six primary runs valid (PG 30/30/31, YB 30/31/30), both controls fired 6/6;
  replicated report and signed analysis written.


## Decision Log

1. **Parent selection by prefix, not per-level reloads.** `parent_max` is a
   valued parameter, so one 10,000-parent fixture serves every level. Confidence:
   High. Alternative: reload a smaller fixture per level (expensive, and it
   changes the data between cells). Affects every cell and the fingerprint.
2. **Added a 0.01% level (1 parent).** The protocol lists four percentages but
   its own falsifier requires ≥5 levels for a bounded crossing; measured levels
   are {1,10,100,1000,10000}. Confidence: High. Flagged for HIGH as a small
   expansion, like the access-paths regime expansion.
3. **Top-N is compared at N=1 only.** `DISTINCT ON` is N=1-only; the report says
   so on every row. Confidence: High. N>1 is a recorded gap (a separate arm).
4. **Two controls, not one.** The protocol names `X-tie`; the brief additionally
   requires the pre-aggregation wrong-grain counterexample. Added
   `x05_aggregate_double_count` and `x06_topn_rank_ties`. Confidence: High.
5. **GuaranteedTie ties the *top* of the smallest multi-child parent.** RANK only
   keeps an extra row when the tie is at the boundary rank, so a lower tie would
   leave the control silent. Confidence: High. Affects x06 only; it is in the
   dataset fingerprint.
6. **`RunCells` refactored onto `RunSpecs`/`runOne`.** A cell must bind its own
   `parent_max`, so the runner works from specs. Confidence: High. Note: the
   unmerged access-paths branch made the same refactor; the later merger must
   keep one copy. Affects the shared harness.
7. **Live teaching-fixture gate now uses N=1.** With N=3 the RANK-tie control
   cannot fire on the teaching fixture (no parent has more than 3 children).
   Confidence: High. Affects the six-control gate only.
8. **New SQL file for the family.** Statements live in
   `sql/select_aggregate_catalog.sql` (names `a01`–`a06`, controls `x05`/`x06`),
   so the merge with the access-paths catalogue is add-a-file rather than
   rewrite. Confidence: Medium (the shared loader sorts filenames; statement
   order changes but no contract depends on it).
9. **`topNArg` normalises N to ≥1 in both the binding and the oracle.** An unset
   N must never make SQL and oracle disagree about rows per parent. Confidence:
   High.
10. **`DISTINCT ON` is its own group (`sel_topn_one`), fixed at N=1.** It cannot
    accept N>1, so tying it to `Args.N` would make the generic validation bind
    N=3 and fail. Confidence: High.

## Limits already known

Single-node only; no 3-node/equal-total/placement arms; one sequential reader,
no writers; five samples cannot support tails; YB storage/RPC counters not
parsed. These are repeated in the environment page and the reports.
