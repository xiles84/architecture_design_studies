# Build the containerized SELECT harness, oracle and readable plan reports

Capability: LOW; role: executor. Dependencies: `task-20261002T131729Z-select-protocol`.

## Goal and fixed scope

Implement the reviewed protocol in Go using the existing study harness and platform boundaries; Study 01's historical harness remains reproducible. Prefer an additive mode/subcommand and a distinct pinned select benchmark image, not a migration of old results or a shared image overwrite.

Embed the named SQL catalogue with explicit projections and stable semantic-group ids. Implement deterministic datasets, narrow/wide rows, controlled fanout/skew and independent Go bag/membership/aggregate/order oracles. Ordering checks require a documented total order; compare bags when order is unspecified. Nullable foreign keys/values, duplicate qualifying children, unmatched parents, all-miss sets, ties and empty input must be exercised.

Add read-path negative controls that must fail: outer-join WHERE mistakenly dropping a preserved parent; nullable NOT IN used as an unconditional NOT EXISTS substitute; wrong deduplication/grain; limit before a multiplicative join or wrong top-N ties. Do not display invalid controls as faster alternatives. Deliberate controls demonstrate detection; valid examples retain their own intended contracts.

Store raw results, plans, result row/byte counts, actual/estimated cardinalities, loops, buffers/temp spill, supported distributed RPCs, planning/execution/end-to-end latency, settings and environment. Plan signatures remove volatile costs/timings, retain access/join/materialization structure and do not assert semantic equivalence alone.

Make comparison reports show valid/invalid/unsupported/timed-out cells, same-plan/no-resolved-difference outcomes, trial and whole-run variability, generator caps and predictions. Generated measurements and interpretation remain separate.

Unit/fixture checks run in Podman. Any engine validation follows the benchmark lock and produces a check report; skip large matrices here. Freeze a committed/tagged harness for subsequent trials; future task branches integrate current main before runs.

## Shared constraints

Read the durable REQUEST.md/PLAN.md, AGENTS.md, current CONTEXT.md and methodology. Claim/guard atomically, work in your own task worktree, stage explicit owned paths, preserve prior signed artifacts and immutable evidence, identify actual model/tool/effort, integrate current main and complete with an annotated tag. No push/pull.
All runtime work is in Podman. Every check/calibration has a report and every database-starting operation holds ads-run-lock. Measurement starts from committed/tagged code. Drain result rows; randomize keys; gate valid timings on independent Go truth; make read-path controls fire per relevant topology/phase. Use >=3 independent whole runs and >=3 trials for argued deltas, seeded interleaving and repeated-design position controls, >=5 levels for a crossover, isolated and blended workloads, actual demand/client caps, per-node throttling and EXPLAIN ANALYZE cost counters.
Follow the reviewed bounded protocol: identical logical data/answer/indexes/settings in query-form isolate pairs; declared index variants include mutation/storage cost; single-node and 3-node per-node/equal-total arms separate; balanced endpoints and actual colocated/non-colocated placement evidence or an explicit gap. One laptop establishes neither independent hardware nor real network effects.
Decide/log/continue local reversible choices. Escalate changing semantics, question, predictions, accepted matrix, oracle/negative control, required evidence or book style that would invalidate downstream acceptance. Resource-infeasible/unsupported/timed-out cells remain visible gaps/results. Never weaken a gate or substitute a different answer to get a faster number. HIGH reviews Decision Log and material evidence before integration; label self-review honestly.

## Acceptance

- Independent Go checks and each deliberately wrong read-path fixture detect the intended defect; SQL semantic groups cannot be pooled accidentally.
- Go tooling/tests run inside pinned Podman images and report all checks; new DB-starting entrypoints enforce ads-run-lock.
- Old Study 01 execution remains reproducible; new output retains complete provenance and separates measurements from interpretations.

## Completion

Record RESULT.md, Decision Log, reports and actual limits. Submit for HIGH review; reconcile concurrent changes in your worktree, validate, fast-forward local main under its integration lock and verify reachability. Annotated milestone: `study-01/select-v1-harness`.

NEXT MODEL: LOW
