# REVIEW — task-20261002T131729Z-select-aggregate-topn

Reviewer: HIGH session `session-20261004T005525Z-c7a6bf` (deepseek-flash / opencode CLI).
Reviewed result commit: `1d153c58` (results/report/analysis) plus `b0b45d20`
(figure corrections made during this review). Branch `study-01/select-aggregate-topn`.

## What I checked

- **Evidence present and consistent.** Six primary runs
  (`select-agg-pg-1b/2/3`, `select-agg-yb-1/2/3`) each carry `aggregate.json`,
  `aggregate.md`, `manifest.env` and resources; all cells `valid`
  (PG 30/30/31, YB 30/31/30), 0 invalid/unsupported/timeout/error. Run tags
  exist and are annotated. The replicated report is generated from the run JSONs
  and keys every aggregate on (engine, topology, axis), so no engine pooling.
- **Correctness gate not weakened.** Every cell is oracle-verified before
  timing; the live six-control gate ran, and both family controls fired in all
  six runs (`x05_aggregate_double_count` 10,000 rows; `x06_topn_rank_ties`
  10,021 rows). The failed first PostgreSQL attempt (`select-agg-pg-1`, exit 2)
  is retained as a harness-defect record and is not counted as a measurement.
- **Figures checked against the report.** Every range in the analysis was
  re-derived from the replicated report. Four ranges were wrong and are corrected
  in `b0b45d20`: YB correlated/post-join 446–507× (not 489–507×); YB lateral at
  100 parents 6.3–7.1× and at 1,000 44.8–49.9× (not 6.3–6.9× / 45.6–49.9×); PG
  DISTINCT ON vs ROW_NUMBER 0.63–0.81× (not 0.63–0.73×) and YB 0.82–1.05× (not
  0.83–1.02×). Conclusions are unchanged.
- **Predictions judged, not restated.** P8 is reported as engine-specific (real
  reversal on PostgreSQL, none on YugabyteDB) and P9 as only partly as framed
  (probe wins only at one parent; DISTINCT ON is the fastest spelling at scale).
  Both are consistent with the measured signs across three runs per engine.

## Decision Log review

1. Parent prefix (`parent_max`) instead of per-level reloads — **accepted**;
   keeps the fixture identical across levels and is in the fingerprint.
2. Added a 0.01% level (1 parent) to reach the ≥5 levels the protocol's own
   falsifier requires — **accepted** as a coverage expansion, recorded.
3. Top-N compared at N=1 only, because DISTINCT ON cannot express N>1 —
   **accepted** and stated on every report row; N>1 recorded as a gap.
4. Two controls (wrong-grain `x05`, tie `x06`) instead of the protocol's single
   `X-tie` — **accepted**; the brief required the counterexample and both are
   live.
5. `GuaranteedTie` ties the top of the smallest multi-child parent — **accepted**;
   without it the tie control would be silent, and it is recorded as a fixture
   shape.
6. `RunCells` refactored onto `RunSpecs`/`runOne` — **accepted**; it is the same
   refactor the access-paths branch made. The later merger must keep one copy;
   flagged for integration reconciliation.
7. Teaching gate moved to N=1 — **accepted**; required for the RANK-tie control
   to fire on the teaching fixture.
8. New SQL file for the family — **accepted**; keeps the merge additive.
9. `topNArg` normalisation shared by binding and oracle — **accepted**; prevents
   an unset N making SQL and oracle disagree.
10. `DISTINCT ON` as its own group fixed at N=1 — **accepted**; the generic
    validation could not otherwise bind a consistent N.

## Findings and defects

- Correctness of the analysis's cited figures: four ranges adjusted (above).
- No gate weakening, no oracle substitution, no hidden failed cell.
- Known limits are recorded in the report and analysis: single host, single
  node, no YB 3-node/equal-total arm, no storage/RPC counters, fixed fanout,
  N=1 only, five samples.
- The 30 shared cells ran at two commits (`08a3c553`, `33b76901`); the only
  difference is the opt-in position cell, and this is recorded in the analysis.

## Decision

**Approved for integration.** The measurements are complete, oracle-gated,
replicated and honestly bounded; the predictions are judged rather than
restated; the one review finding (figures) is corrected before integration.
Integration must reconcile the shared `RunSpecs`/`select_validation.go` with the
access-paths branch.
