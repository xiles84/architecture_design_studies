# Progress — task-20261002T131729Z-select-join-cte

Attempt `attempt-20261004T152035Z-f82109`, branch `study-01/select-join-cte`,
LOW executor (model `deepseek-flash`, tool opencode/CLI, effort unknown).
Milestone tag: `study-01/select-v2-join-cte`. Claim `claim-a93dbde0c868f001`,
epoch 1.

Previous attempt `attempt-20261004T014019Z-c3a6b3` yielded a committed
recon/premeasurement plan with no run started; its `PROGRESS.md` and the
bounded-matrix file it authored are preserved unedited. This attempt resumes
from the dated HIGH handoff `docs/select-query-templates/JOIN-HANDOFF-20261004.md`
and the current local `main` (integrated before implementation).

## Status

- [x] Claim the task; integrate current local `main` into the retained worktree.
- [x] Implement the join/CTE family (SQL, fixture, oracles, runner, report,
      tests) as a self-contained addition; commit before any run.
- [x] Build the benchmark image and run per-engine feasibility checks.
- [x] Run >=3 whole runs per single-node engine (5 trials, seeded interleaving,
      position control): 6 runs, **91/91 valid each**.
- [x] Generate reports, RESULT.md, coverage matrix and exact input digests.
- [x] Establish the status of the 3-node YB subset and the placement pair as an
      explicit resource/harness coverage gap (not run).
- [ ] Submit for HIGH review.

### Outcome summary

Two harness defects were found by feasibility and fixed before the matrix
(inverted control detection; a MATERIALIZED cell was judged by a non-existent
`Materialize` node — PostgreSQL 17 prints a `CTE Scan` instead). The matrix:
six inner-join spellings are plan-identical within an engine; `MATERIALIZED` on
a joined-output CTE costs 6.8x–88x on PostgreSQL but 1.02x–1.61x on YugabyteDB;
the reuse crossing exists only on PostgreSQL, where `NOT MATERIALIZED` is faster
at a 1% consumer but slower at 100%, while on YugabyteDB it is uniformly ~9–11x
slower. Measured pairs and gaps are in `RESULT.md`; interpretation in the signed
analysis `reports/analyses/20261004T-select-join-cte--deepseek-flash--2026-10-04.md`.

### Post-run calibration outcome

The handoff's 200,000-row fixture and all six thresholds were **feasible on both
engines** (exact counts, 179,911-row widest drain, token supported), so the
calibration did not reject any axis. No axis was altered after measurement.

## Decision Log

1. **Separate schema and catalogue.** The join family lives in its own
   `selectjoin` schema and its own embedded catalogue (`sql/join/`), loaded by
   `LoadJoin()`; the teaching `Load()` does not see it. Reason: adding the join
   statements to the shared catalogue would have made every other family's
   validation gate run them against a fixture with no `rank` column. Confidence
   High. Affects: all join cells; no other family.
2. **`rank` permutation fixture at 200,000 rows** with thresholds
   20/200/2,000/20,000/100,000/180,000. Reason: exact selectivity, fixed index
   set, handoff-specified. Confidence High. Affects: every bag cell.
3. **Reuse ladder generated from committed code**, not hand-written 30 times:
   `buildReuseSQL` emits the CTE policy keyword and the `UNION ALL` references;
   the catalogue registers them programmatically. Reason: keeps the ladder
   consistent and reviewable; the SQL is still deterministic from committed
   code. Confidence High. Affects: all reuse cells.
4. **Materialization evidence is enforced.** A forced `MATERIALIZED` cell whose
   plan contains no `Materialize` node is recorded `unsupported`, never as a
   default fold. `NOT MATERIALIZED` / default record the observed flag.
   Confidence High. Affects: `j04m`/`j05m`/`r_mat_*`.
5. **Unsupported cells do not fail a run; invalid/timeout/error cells do.** The
   handoff requires unsupported behavior to be an explicit result, and an engine
   that cannot express a policy is not a wrong answer. Confidence High. Affects:
   run exit codes and report tables.
6. **Controls placed where they can fire.** `x01_optional_where` on the edge
   fixture (an unmatched parent exists); `x02_join_limit_before` on the large
   fixture (the fixed limit 7 truncates). Reason: a control that cannot fire
   proves nothing. Confidence High. Affects: the gate.
7. **Reuse token `ascii(substr(md5(repeat(id::text,32)),1,1))`**, support
   verified per engine at feasibility. If an engine refuses it, the reuse
   matrix is an explicit gap for that engine. Confidence Medium (engine
   support untested until feasibility). Affects: reuse cells.
8. **Cluster/placement arms are conditional.** The handoff's 3-node YB RF=3
   subset is attempted only after single-node work. If the fixed 2 GiB machine
   cannot carry it, it is recorded as a coverage gap with the failing resource
   evidence. A single seed endpoint in a cluster DSN is not balanced cluster
   access and will be labelled as such. Confidence Medium. Affects: cluster
   cells and the coverage matrix.

## Notes carried from the handoff

- One-use spellings and reuse contracts get separate results and predictions.
- Plans, valid-cell latency/throughput, variability, failure regimes and
  equal/no-resolved-difference outcomes are reported without a signed
  interpretation; a LOW signed interpretation is optional and separately
  attributed.
- The batch of commands started from a committed tree; producing runs get
  `run/study-01/<id>` tags via the runner.
