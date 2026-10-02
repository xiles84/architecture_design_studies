# Measure index-friendly predicates, OR alternatives and pagination boundaries

Capability: LOW; role: executor. Dependencies: `task-20261002T131729Z-select-harness`.

## Goal and fixed scope

Execute sargable ranges versus equivalent function/cast predicates with exact type/time-zone semantics, matched indexed/unindexed pairs and supported composite/covering alternatives. Separate query spelling from adding an expression/index, so index maintenance is not attributed to syntax.

Compare OR with UNION/UNION ALL only under matched bag semantics and proven overlap handling; include overlapping branches as a negative control for naive UNION ALL. Document why UNION's duplicate removal can also change an answer.

Compare OFFSET with keyset pagination using the same stable total order and output page under a fixed snapshot and known cursor; keyset's sequential cursor contract differs from arbitrary-page access. Sweep depth and explain that concurrency can shift pages under a different consistency contract.

Use the reviewed selectivity/cardinality/skew/row-width and pagination-depth sweeps. Prepared generic/custom plans or stale-statistics controls must be separately labelled, capability-verified on pinned engines and must not contaminate the default arm.

Measure read throughput/tails, estimated versus actual rows, buffer/RPC work, spill, rows/bytes drained, and inserts/corrections/deletes plus index size/load/maintenance for changed indexes. Include a representative three-node placement/routing arm with actual placement and balanced endpoints; separate equal-total budgets.

No hard-coded forced planner settings underpin the headline recommendation. Record memory limits and timeout boundaries; budget/resource-infeasible cells remain gaps.

## Shared constraints

Read the durable REQUEST.md/PLAN.md, AGENTS.md, current CONTEXT.md and methodology. Claim/guard atomically, work in your own task worktree, stage explicit owned paths, preserve prior signed artifacts and immutable evidence, identify actual model/tool/effort, integrate current main and complete with an annotated tag. No push/pull.
All runtime work is in Podman. Every check/calibration has a report and every database-starting operation holds ads-run-lock. Measurement starts from committed/tagged code. Drain result rows; randomize keys; gate valid timings on independent Go truth; make read-path controls fire per relevant topology/phase. Use >=3 independent whole runs and >=3 trials for argued deltas, seeded interleaving and repeated-design position controls, >=5 levels for a crossover, isolated and blended workloads, actual demand/client caps, per-node throttling and EXPLAIN ANALYZE cost counters.
Follow the reviewed bounded protocol: identical logical data/answer/indexes/settings in query-form isolate pairs; declared index variants include mutation/storage cost; single-node and 3-node per-node/equal-total arms separate; balanced endpoints and actual colocated/non-colocated placement evidence or an explicit gap. One laptop establishes neither independent hardware nor real network effects.
Decide/log/continue local reversible choices. Escalate changing semantics, question, predictions, accepted matrix, oracle/negative control, required evidence or book style that would invalidate downstream acceptance. Resource-infeasible/unsupported/timed-out cells remain visible gaps/results. Never weaken a gate or substitute a different answer to get a faster number. HIGH reviews Decision Log and material evidence before integration; label self-review honestly.

## Acceptance

- Predicate, OR and pagination pairs preserve their declared bag/order/page contracts; overlap/cursor/time-zone controls fire.
- Index changes include write/storage costs, with query-form and physical-design effects isolated.
- Replicated reports expose selective/broad and shallow/deep regimes, plan evidence, client limits and verified distributed coverage.

## Completion

Record RESULT.md, Decision Log, reports and actual limits. Submit for HIGH review; reconcile concurrent changes in your worktree, validate, fast-forward local main under its integration lock and verify reachability. Annotated milestone: `study-01/select-v5-access-paths`.

NEXT MODEL: LOW
