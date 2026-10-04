# Progress — task-20261002T131729Z-select-join-cte (corrective attempt)

Attempt `attempt-20261004T173630Z-d9091f`, claim `claim-03114d688bd5fe53`,
branch `study-01/select-join-cte`, LOW executor (model `deepseek-flash`, tool
opencode/CLI, effort unknown). Resumes after the independent HIGH review
`docs/select-query-templates/reviews/20261004-join-high/` moved the task to
`changes_requested`. The previous two attempts, the six original runs, their
reports, signed analysis and all ten producing tags are preserved unedited.

## Corrective checklist (from CORRECTIVE-HANDOFF.md)

- [x] **C1 reporting/evidence fixes**
  - inputs digest: one shared `joinInputsDigest`; the re-render no longer appends
    a second newline. Re-rendering the original pg1 input now advertises the
    original `070df793…` (the review's broken re-render showed `7c879564…`).
  - client vs engine provenance separated (`ClientImage`/`ClientImageDigest` vs
    `EngineImage`/`EngineImageDigest`); the runner passes `ADS_ENGINE_IMAGE`.
  - replicated report keyed by a compatibility class (environment, dataset,
    client and engine image digests, settings) in addition to engine/topology/
    axis; classes are printed and never pooled.
  - engine-specific EXPLAIN options probed and recorded (PG VERBOSE/BUFFERS, YB
    DIST/DEBUG) with fallback evidence.
  - token volatility evidence recorded (distinct ids → distinct tokens; reuse
    plan carries the expression).
  - latency min/median/max and derived queries/s added, labelled derived (not a
    capacity test).
  - live gate extended: empty-input fixture and the filtered-child LEFT JOIN
    spelling `j11_optional_derived`, plus the existing firing WHERE control.
  - failing-before/passing-after tests in `harness/select_join_report_test.go`.
- [ ] **C2 Study 01 domain arm** (person/donation/charity) + re-run matrix.
- [ ] **C3 effective WSL resource recheck** and budget calibration.
- [ ] **C4 bounded single/cluster topology subset**, round-robin endpoints, RF3,
      placement and colocation.
- [ ] **C5 matched blended read/write phase** (one fixed writer).
- [ ] **C6** new signed successor analysis with corrected verdicts.

## Decision Log (this attempt)

1. **Restore the owner-normalized CRLF copies of HIGH's review files to the
   committed LF bytes** rather than commit a line-ending-only change to another
   session's signed review packet. No content change was present.
2. **Do not overwrite the six original reports or the original signed analysis.**
   The fixes are validated by a corrected re-render under the new review outputs
   and by tests; new runs use new ids, as the handoff requires.
3. **Explain options are probed, not assumed**, and the fallback is recorded as
   evidence, so an absent RPC counter is reported as an absence rather than
   described as "not parsed".
4. **Compatibility class rather than rejection** for the replicated report:
   runs with differing environment/fixture/settings are partitioned explicitly,
   matching the handoff's "reject … or partition" and keeping one report usable
   across the intended single/cluster and blended arms.
5. **`selectbench.Provenance` gains two fields** (`ClientImage`,
   `ClientImageDigest`); the change is additive and only the join family sets
   them, so the teaching families' provenance is unchanged.

## Still open / carry forward

C2–C6 are not done in this checkpoint. They need new fixtures/runs (domain arm,
blended phase, cluster subset) and a new signed interpretation. The next session
resumes from this decision log with the same claim lifecycle.
