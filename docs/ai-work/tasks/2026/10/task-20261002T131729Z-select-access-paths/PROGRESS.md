# Progress — task-20261002T131729Z-select-access-paths

Owner: LOW executor, deepseek-flash (DeepSeek) via opencode CLI.
Branch `study-01/select-access-paths`, worktree `.worktrees/select-access-paths`.
Milestone tag on completion: `study-01/select-v5-access-paths`.

## Plan

1. Add the access-path matrix to the additive SELECT harness (done).
2. Commit the bounded measurement plan and environment before any run (done).
3. Run ≥3 whole runs on PostgreSQL single-node and ≥3 on YugabyteDB single-node.
4. Generate per-run and replicated reports; write the signed analysis with
   P10/P11 verdicts; record limits and gaps.
5. Submit for HIGH review, integrate into local main, tag, complete.

## Decision log

Each entry: decision — reason — confidence — alternative considered — affected work.

1. **Added a `deep` scale (200 × 60 = 12,000 rows) for the pagination axis.**
   The protocol's deepest page depth is 10,000, so the fixture must hold more
   than 10,000 rows. `large` (≈60,000 rows, one hot parent) would load far more
   slowly on the single-node 2 GiB YugabyteDB guest and is not needed for a page
   pointer; `deep` keeps the load bounded while making every depth a real page.
   Confidence High. Alternative: `large`. Affects: all access-path cells
   (dataset identity); recorded in every manifest/fingerprint.

2. **Added `s20_union_predicate` (plain `UNION`) as a third spelling.** The
   protocol's P10 pair is OR vs `UNION ALL`; because the arms are disjoint,
   `UNION` answers the same bag and shows the de-duplication cost the overlap-safe
   choice avoids. It is the answer-safe fallback the brief's "proven overlap
   handling" requires, adds no new scientific question, and never competes across
   contracts. Confidence High. Alternative: run only OR vs UNION ALL. Affects:
   P10 evidence.

3. **The predicate-alternative oracle defaults nil arms to the same SQL the
   harness binds by default.** The reviewed live teaching-fixture gate runs the
   whole catalogue with default args; if the oracle treated a nil arm as empty
   while the SQL bound `amount <= 100`, the gate would fail for a harness reason.
   The default is now explicit and shared (`defaultArms`). Confidence High.
   Affects: harness validation only; matrix cells always set both arms.

4. **The expression index is a labelled package arm, never a syntax arm.** The
   default arm runs with the index dropped; the index arm runs the identical
   query text with the index present and reports its size and insert cost in a
   separate section. This keeps the physical factor from being attributed to the
   spelling, per the protocol's attribution rule. Confidence High. Affects: P10
   interpretation (analysis).

5. **Fixed-snapshot page equality is checked inside one transaction.** A
   `REPEATABLE READ` transaction runs both pages at every depth and requires equal
   rows; if the engine rejects that isolation it falls back to the default and
   records the exact isolation string, so the check is never silently assumed.
   Confidence Medium (YB isolation behaviour is engine-reported, not recalled).
   Alternative: rely on the absence of writers alone. Affects: P11 evidence.

6. **`run-select.sh`'s committed-tree check now excludes `results/` and
   `reports/`.** Git's `:!` pathspec does not suppress *untracked* files, so the
   first run's output would make the tree look dirty and block every later run —
   the exact shared-infrastructure defect already in `LESSONS_LEARNED`. The
   filter matches the newer studies' runners. Confidence High. Affects: the
   runner only; no result bytes.

7. **The index-maintenance cost is a 500-row multi-row INSERT timed once with and
   once without the index.** This is a labelled write-cost indicator, not a
   write workload; it isolates the index's maintenance term and is reported with
   that limit. Confidence Medium. Alternative: a full write phase (out of scope
   for a syntax-isolate task). Affects: index-cost section.

## Checkpoints

- Matrix/oracle logic and unit tests committed (`2b0abbfd`).
- Harness `-cmd select-access`, runner `--family access`, reports and docs
  committed as the premeasurement checkpoint (hash in `RESULT.md`).
