# Measure join aggregation, correlated subqueries and per-parent top-N

Capability: LOW; role: executor. Dependencies: `task-20261002T131729Z-select-harness`.

## Goal and fixed scope

Execute equivalent parent-total pairs: aggregate after join, pre-aggregate children before join, correlated scalar aggregate and supported lateral aggregate. Fix empty-parent SUM/COUNT/COALESCE behavior and avoid multiplying totals by joining a second one-to-many relation. Required pre-aggregation counterexample deliberately shows the wrong grain.

Execute latest/top-N-per-parent pairs: correlated or LATERAL ORDER BY ... LIMIT N with a tie-breaker; window ROW_NUMBER filtering with the same tie rule; supported DISTINCT ON for N=1 with matching ordering. RANK/DENSE_RANK keep-all-ties answers are separate contracts. Global top-N and per-parent top-N cannot be mixed.

Sweep selected-parent count, fanout/skew, N and memory at fixed logical volume; exercise small indexed probes and broad scans. Capture nested-loop iterations, scans, sort/hash/aggregate cardinality, spill and RPCs. Hold indexes fixed in query-form isolate pairs; any index redesign is a separate experimental arm.

Persistent rollups/materialized read models, their freshness and maintenance are not equivalent to transient pre-aggregation by default. Reuse existing Study 01/04 evidence or link Study 08 coverage; don't duplicate those studies or claim their write cost is measured here.

Execute the reviewed replicated topology/resources and isolated/mixed matrix. Correctness checks use a stable snapshot appropriate to the answer; random concurrent state cannot be compared to an earlier unversioned oracle.

## Shared constraints

Read the durable REQUEST.md/PLAN.md, AGENTS.md, current CONTEXT.md and methodology. Claim/guard atomically, work in your own task worktree, stage explicit owned paths, preserve prior signed artifacts and immutable evidence, identify actual model/tool/effort, integrate current main and complete with an annotated tag. No push/pull.
All runtime work is in Podman. Every check/calibration has a report and every database-starting operation holds ads-run-lock. Measurement starts from committed/tagged code. Drain result rows; randomize keys; gate valid timings on independent Go truth; make read-path controls fire per relevant topology/phase. Use >=3 independent whole runs and >=3 trials for argued deltas, seeded interleaving and repeated-design position controls, >=5 levels for a crossover, isolated and blended workloads, actual demand/client caps, per-node throttling and EXPLAIN ANALYZE cost counters.
Follow the reviewed bounded protocol: identical logical data/answer/indexes/settings in query-form isolate pairs; declared index variants include mutation/storage cost; single-node and 3-node per-node/equal-total arms separate; balanced endpoints and actual colocated/non-colocated placement evidence or an explicit gap. One laptop establishes neither independent hardware nor real network effects.
Decide/log/continue local reversible choices. Escalate changing semantics, question, predictions, accepted matrix, oracle/negative control, required evidence or book style that would invalidate downstream acceptance. Resource-infeasible/unsupported/timed-out cells remain visible gaps/results. Never weaken a gate or substitute a different answer to get a faster number. HIGH reviews Decision Log and material evidence before integration; label self-review honestly.

## Acceptance

- Aggregates and top-N have explicit matched empty-set/tie/grain contracts and Go truth checks.
- Parent-count/fanout/N sweeps identify measured regimes or unresolved boundaries; supported feature limits are visible.
- Generated reports connect cost counters and timings without attributing simultaneous query/index changes to one cause.

## Completion

Record RESULT.md, Decision Log, reports and actual limits. Submit for HIGH review; reconcile concurrent changes in your worktree, validate, fast-forward local main under its integration lock and verify reachability. Annotated milestone: `study-01/select-v4-aggregate-topn`.

NEXT MODEL: LOW
