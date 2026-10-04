# REVIEW — task-20261002T131729Z-select-access-paths

Reviewer: HIGH session `session-20261004T005525Z-c7a6bf` (deepseek-flash / opencode CLI).
Reviewed result commit: `f4e84bec`. Branch `study-01/select-access-paths`.

## What I checked

- **Evidence present and consistent.** Seven whole runs (PG 1/2/3 + a PG position
  control, YB 1/2/3), each with `access.json`, `access.md`, `manifest.env` and
  resources; all cells valid (34/34, 35/35), 0 invalid/unsupported/timeout/error.
  A labelled broad-only pilot tranche (four runs) is committed as corroboration
  and is excluded from the primary replicated set. Run tags are annotated.
- **Correctness gate not weakened.** Every cell is oracle-verified before timing;
  the required `x05_union_all_overlap` control fired in all six shapes in all
  seven runs (228 duplicated rows selective, 4,546 broad); `x01`–`x03` are
  recorded `control_absent` where the measurement fixture cannot exercise them,
  not counted as passes.
- **Figures checked against the report.** Keyset/OFFSET 4.3–5.1× (PG) and
  1.5–1.9× (YB) at depth 10,000; OR vs UNION/UNION ALL on YB 1.13–2.12× across
  36 shape×run ratios; expression index 0.845 → 0.265 ms on PG and no gain on YB;
  position control 0.885 (PG) / 1.112 (YB); index size 114,688 bytes. All match
  the generated reports. Page equality under `REPEATABLE READ` held at all five
  depths in every run.
- **Predictions judged, not restated.** P11 confirmed with an engine-specific
  crossing (PG 100→1,000; YB 1,000→10,000). P10 reported as falsified as stated:
  no shape/regime shows UNION ALL with a stable advantage, and on YB OR is
  consistently faster. The analysis names where the measurement is weak.

## Decision Log review

1. `deep` scale (12,000 rows) for a real 10,000-depth page — **accepted**;
   recorded in every manifest/fingerprint.
2. Added plain `UNION` beside `UNION ALL` on provably disjoint arms —
   **accepted**; shows the de-duplication cost with no new scientific question.
3. Default arms shared between oracle and bound SQL — **accepted**; a harness
   correctness fix, not a measurement change.
4. Expression index as a labelled package arm with the query text unchanged —
   **accepted**; this is exactly the protocol's attribution rule and is why the
   3.2× PG effect is attributed to the index, not the spelling.
5. Fixed-snapshot page equality with an engine-reported isolation fallback —
   **accepted**; the exact isolation is recorded (`repeatable read` on both).
6. Runner clean-tree check excludes untracked `results/`/`reports/` —
   **accepted**; mirrors the existing shared-infrastructure lesson.
7. 500-row INSERT as a labelled index-maintenance indicator — **accepted**; the
   report states it is not a write workload, and the sign was unstable, so no
   maintenance cost is claimed.
8. Added the `selective` regime, making the P10 matrix {selective, broad} ×
   {indexable, expression, cast} — **accepted**. The acceptance itself requires
   the reports to expose selective/broad regimes, the broad-only pilot proved the
   mechanism was untestable at 39% selectivity, and the expansion changes no
   contract, prediction, control or gate. Correctly flagged for review; this
   review confirms it.

## Findings and defects

- No figure errors found; no gate weakening; no hidden failed cell.
- Limits are honest and recorded: one host, single node, no YB 3-node
  per-node/equal-total arm, no verified placement, no concurrent demand, five
  trials, latency not throughput, YB storage/RPC counters not parsed, and the
  YB index-maintenance probe is unusable (with-index faster, 3/3) and reported as
  absent rather than as a cost.

## Decision

**Approved for integration.** The expansion in Decision 8 is required by the
acceptance and was surfaced rather than smuggled in; the falsified prediction is
reported as falsified; the evidence is complete and oracle-gated. Integration
must reconcile the shared `runner.go`/`select_validation.go` with the
aggregate-topn branch.
