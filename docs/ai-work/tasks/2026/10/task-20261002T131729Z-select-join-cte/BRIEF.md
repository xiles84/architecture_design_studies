# Measure join/filter spellings and CTE materialization success and failure regimes

Capability: LOW; role: executor. Dependencies: `task-20261002T131729Z-select-harness`.

## Goal and fixed scope

Execute the reviewed inner-join/filter matrix: WHERE after join, filtered derived table, CTE of filtered input, CTE of joined rows, predicate in ON and legacy comma join. Keep data, projection, indexes, semantics and parameters matched. Compare supported default CTE folding with MATERIALIZED and NOT MATERIALIZED, including one-use selective access and repeated expensive input reuse as separate contracts.

Sweep filter selectivity (proposed 0.01%, 0.1%, 1%, 10%, 50%, 90%), fanout, row width and memory/reuse conditions according to the committed protocol; final levels must be feasibility-calibrated before measurements. Do not force join algorithm switches and attribute that result to spelling.

Keep LEFT JOIN predicate-placement examples as distinct semantic contracts and detected invalid substitutions; an incorrect faster answer is not a performance competitor.

Capture normalized plan signatures and actual plans for each group/condition. Report plan-equivalent variants and unresolved deltas plainly; a no-difference result is useful teaching evidence. Document observed materialization/folding on each pinned engine instead of transferring PostgreSQL behavior by assumption.

Use the program's replicated, serial, topology/resource/placement and isolated/mixed workload requirements. If the protocol budgets a representative topology subset, it must be labelled and every omitted dimension retained as a gap.

## Shared constraints

Read the durable REQUEST.md/PLAN.md, AGENTS.md, current CONTEXT.md and methodology. Claim/guard atomically, work in your own task worktree, stage explicit owned paths, preserve prior signed artifacts and immutable evidence, identify actual model/tool/effort, integrate current main and complete with an annotated tag. No push/pull.
All runtime work is in Podman. Every check/calibration has a report and every database-starting operation holds ads-run-lock. Measurement starts from committed/tagged code. Drain result rows; randomize keys; gate valid timings on independent Go truth; make read-path controls fire per relevant topology/phase. Use >=3 independent whole runs and >=3 trials for argued deltas, seeded interleaving and repeated-design position controls, >=5 levels for a crossover, isolated and blended workloads, actual demand/client caps, per-node throttling and EXPLAIN ANALYZE cost counters.
Follow the reviewed bounded protocol: identical logical data/answer/indexes/settings in query-form isolate pairs; declared index variants include mutation/storage cost; single-node and 3-node per-node/equal-total arms separate; balanced endpoints and actual colocated/non-colocated placement evidence or an explicit gap. One laptop establishes neither independent hardware nor real network effects.
Decide/log/continue local reversible choices. Escalate changing semantics, question, predictions, accepted matrix, oracle/negative control, required evidence or book style that would invalidate downstream acceptance. Resource-infeasible/unsupported/timed-out cells remain visible gaps/results. Never weaken a gate or substitute a different answer to get a faster number. HIGH reviews Decision Log and material evidence before integration; label self-review honestly.

## Acceptance

- All six equivalent inner-join spellings are verified against the same independent answer and measured on the reviewed regimes.
- Supported CTE controls distinguish folding/fencing and reuse costs; unsupported behavior is reported explicitly.
- Reports include plans, valid-cell latency/throughput, variability, failure regimes and equal/no-resolved-difference outcomes without a signed interpretation.

## Completion

Record RESULT.md, Decision Log, reports and actual limits. Submit for HIGH review; reconcile concurrent changes in your worktree, validate, fast-forward local main under its integration lock and verify reachability. Annotated milestone: `study-01/select-v2-join-cte`.

NEXT MODEL: LOW
