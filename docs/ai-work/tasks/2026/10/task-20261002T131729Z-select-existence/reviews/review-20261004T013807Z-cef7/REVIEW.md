# REVIEW — task-20261002T131729Z-select-existence

Reviewer: HIGH session `session-20261004T005525Z-c7a6bf` (deepseek-flash / opencode CLI).
**Self-review: this session both executed and reviewed the task**, because the
owner declared it a LOW+HIGH session. The Decision Log, the evidence and the
figures below were checked against the generated reports; an independent reviewer
should still read the analysis.

Reviewed result commit: `89daacf3`. Branch `study-01/select-existence`.

## What I checked

- **Evidence present and consistent.** Six primary runs
  (`select-exist-pg-1b/2/3`, `select-exist-yb-1/2/3`) each carry `existence.json`,
  `existence.md`, `manifest.env` and resources; all cells valid (PG 21/22/21, YB
  21/22/21), 0 invalid/unsupported/timeout/error. Each level's parents, children
  and fingerprint are recorded. Run tags are annotated.
- **Correctness gate not weakened.** Every cell is oracle-verified before timing;
  the live eight-control teaching gate passed, and both family controls fired at
  all three fanout levels in all six runs: x02 nullable `NOT IN` with 0 rows
  (18/18) and x07 membership multiplier with 7,240–7,241 rows (18/18). The first
  attempt (`select-exist-pg-1`) failed the gate because the harness built the raw
  fixture; it is retained as a harness-defect record and not counted.
- **Figures checked against the report.** PG `join_distinct`/`exists`
  1.44–1.47 (fanout 1), 1.30–1.57 (10), 3.61–4.40 (100); YB 0.92–1.06 at every
  level; `in`/`exists` 0.86–1.13 (PG) and 0.99–1.04 (YB); PG fanout-100
  `not_exists` 0.22–0.34 ms vs `left_is_null` 0.63–0.72 and `not_in_notnull`
  0.68–0.72; position control 0.871 (PG) / 0.995 (YB), same plan. All match.
- **Predictions judged.** P6 is reported as engine-specific (confirmed on PG,
  refuted on YB) and P7 as the semantic trap plus a fanout-growing performance
  gap. Both are consistent across three runs per engine.

## Decision Log review

1. Fanout as a fixture axis with constant child volume — **accepted**; the report
   records each level's fingerprint. (The confound with parent count is stated
   as a limit.)
2. Reuse of the s06–s10 statements plus e12/e13 — **accepted**; no duplicate SQL
   for the same question.
3. `IN`/`NOT IN` at the NOT-NULL contract — **accepted**; the nullable case is
   the x02 control, never ranked.
4. `x07` membership multiplier — **accepted**; it is the plain join the brief
   says is illustrated separately.
5. Absence `{NULL, NOT NULL}` as contract, not spelling — **accepted**.
6. Reuse of the shared `CellSpec`/`RunSpecs` runner — **accepted**; no new runner.
7. Position control — **accepted**.
8. `ExistenceShape` — **accepted**; the fix is what makes both controls fire, and
   a permanent test now builds the real per-level fixtures and asserts it.
9. One fingerprint per level — **accepted**.

## Findings and defects

- The first whole run failed its gate on the raw fixture; fixed in `805b5b85`
  and guarded by `TestExistenceShapeMakesBothControlsFireAtEveryLevel`. This is
  the same class of defect the programme's controls exist to catch, and it was
  caught before any number was published.
- Two analysis wordings were corrected during review (the fanout-100 `exists`
  bound and the multiplier row-count description).
- Limits are recorded: fanout/parent-count confound, three levels, no
  selectivity/hit/skew axis, single host and single node, five samples, YB near
  its per-statement floor and no storage/RPC counters.

## Decision

**Approved for integration.** The measurements are complete, oracle-gated and
replicated; the control defect was found by the gate and fixed before
publication; the engine-specific verdict is stated as engine-specific. The
self-review is disclosed rather than hidden.
