# Measure EXISTS, IN, JOIN DISTINCT and anti-join templates

Capability: LOW; role: executor. Dependencies: `task-20261002T131729Z-select-harness`.

## Goal and fixed scope

Execute parent-membership pairs: EXISTS, IN subquery and JOIN plus the exact DISTINCT needed for the same parent identity. Ordinary child-row joins have a different grain and are illustrated separately. Execute absence pairs: NOT EXISTS, LEFT JOIN with a guaranteed non-null child identity tested for NULL, and NOT IN only under an explicit non-null contract; nullable NOT IN is a labelled semantic counterexample.

Sweep child fanout, qualifying proportion, parent hit/miss fraction and hot-parent skew while holding total logical volume constant for each shape sweep. Include empty right-side and NULL fixtures. Never deduplicate names when distinct parent ids are required.

Measure early-stop/semi/anti plans, row expansion/dedup work, sort/hash spill, index probes, actual rows and distributed RPCs, beside end-to-end throughput/latency with rows drained. Use the reviewed serial replicated matrix and isolated/mixed snapshots; don't infer one algorithm from SQL syntax.

Report where every form succeeds, hurts, has no measurable difference, is semantically unsuitable, unsupported or timed out. Keep invalid controls out of valid rankings.

## Shared constraints

Read the durable REQUEST.md/PLAN.md, AGENTS.md, current CONTEXT.md and methodology. Claim/guard atomically, work in your own task worktree, stage explicit owned paths, preserve prior signed artifacts and immutable evidence, identify actual model/tool/effort, integrate current main and complete with an annotated tag. No push/pull.
All runtime work is in Podman. Every check/calibration has a report and every database-starting operation holds ads-run-lock. Measurement starts from committed/tagged code. Drain result rows; randomize keys; gate valid timings on independent Go truth; make read-path controls fire per relevant topology/phase. Use >=3 independent whole runs and >=3 trials for argued deltas, seeded interleaving and repeated-design position controls, >=5 levels for a crossover, isolated and blended workloads, actual demand/client caps, per-node throttling and EXPLAIN ANALYZE cost counters.
Follow the reviewed bounded protocol: identical logical data/answer/indexes/settings in query-form isolate pairs; declared index variants include mutation/storage cost; single-node and 3-node per-node/equal-total arms separate; balanced endpoints and actual colocated/non-colocated placement evidence or an explicit gap. One laptop establishes neither independent hardware nor real network effects.
Decide/log/continue local reversible choices. Escalate changing semantics, question, predictions, accepted matrix, oracle/negative control, required evidence or book style that would invalidate downstream acceptance. Resource-infeasible/unsupported/timed-out cells remain visible gaps/results. Never weaken a gate or substitute a different answer to get a faster number. HIGH reviews Decision Log and material evidence before integration; label self-review honestly.

## Acceptance

- Membership and absence groups have identical output grain within each valid pair; NULL and duplicate traps are independently detected.
- Fanout/selectivity/hit-rate/skew sweeps include enough levels to bound any claimed crossover.
- Raw plans and generated reports retain valid/invalid/unsupported cells, replication, resources and topology coverage.

## Completion

Record RESULT.md, Decision Log, reports and actual limits. Submit for HIGH review; reconcile concurrent changes in your worktree, validate, fast-forward local main under its integration lock and verify reachability. Annotated milestone: `study-01/select-v3-existence`.

NEXT MODEL: LOW
