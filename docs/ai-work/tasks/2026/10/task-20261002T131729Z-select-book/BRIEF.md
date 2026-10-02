# Publish concept-first SELECT examples, measured case boxes and a template selection guide

Capability: LOW; role: executor. Dependencies: `task-20261002T131729Z-select-synthesis`, `task-20261001T155900Z-book-ed3-pilot`.

## Goal and fixed scope

Use the approved synthesis/handoff and the CURRENT Edition 3 macros/style after the pilot is integrated. Add a short concept-first sequence about query answer shape, joins/filtering, existence, aggregation/top-N, and access paths/pagination; keep each concept chapter within the style's normal page/paragraph targets, and place longer SQL/plan details in a compact appendix/linked example atlas.

Start with the original customer/order fixture and show all six equivalent INNER JOIN forms plus EXISTS under its distinct parent-membership question. Use SQL highlighting, explicit columns and exact outputs; any numbers in the toy dataset are labelled example data, never measurements. Show unmatched-parent LEFT JOIN and duplicate/NULL examples before asking readers to compare timings.

Include runnable container-only SQL templates with requirements (grain, NULLs, tie ordering, cursor semantics, engine/version). For each family supply a recommended readable starting point, alternatives, When it fits/When to avoid rows, a verified observed plan excerpt or labelled conceptual diagram and measured case boxes from registered claims.

Separate syntax choices that produced the same plan/no-resolved-difference from plan-changing mechanisms, and correctness failures from slow valid queries. Numbers appear in evidence-backed case boxes under STYLE.md, with a named environment/engine/topology, varied axis, uncertainty, boundary and claim badge. No workload-independent speed ranking.

Wire the SELECT guide into how-to-choose and glossary/checklist, preserving all earlier claims and simultaneous book work. Modify STYLE.md only to record this owner's scope extension and necessary navigation; retain the approved style.

Run evidence validation/check, canonical pinned Podman book build with manifest and serial book-build lock. Render and visually inspect the new pages: legible SQL, unclipped plan snippets, compact case tables/curves, correct links and outputs. Figures follow FIGURES.md and the two-reader policy; label self-review honestly if independence is unavailable. Require HIGH review of content/evidence/visual result before integration; rebuild tracked dist from the approved combined source. No new owner style approval is needed unless the approved style materially changes.

## Shared constraints

Read the durable REQUEST.md/PLAN.md, AGENTS.md, current CONTEXT.md and methodology. Claim/guard atomically, work in your own task worktree, stage explicit owned paths, preserve prior signed artifacts and immutable evidence, identify actual model/tool/effort, integrate current main and complete with an annotated tag. No push/pull.
All runtime work is in Podman. Every check/calibration has a report and every database-starting operation holds ads-run-lock. Measurement starts from committed/tagged code. Drain result rows; randomize keys; gate valid timings on independent Go truth; make read-path controls fire per relevant topology/phase. Use >=3 independent whole runs and >=3 trials for argued deltas, seeded interleaving and repeated-design position controls, >=5 levels for a crossover, isolated and blended workloads, actual demand/client caps, per-node throttling and EXPLAIN ANALYZE cost counters.
Follow the reviewed bounded protocol: identical logical data/answer/indexes/settings in query-form isolate pairs; declared index variants include mutation/storage cost; single-node and 3-node per-node/equal-total arms separate; balanced endpoints and actual colocated/non-colocated placement evidence or an explicit gap. One laptop establishes neither independent hardware nor real network effects.
Decide/log/continue local reversible choices. Escalate changing semantics, question, predictions, accepted matrix, oracle/negative control, required evidence or book style that would invalidate downstream acceptance. Resource-infeasible/unsupported/timed-out cells remain visible gaps/results. Never weaken a gate or substitute a different answer to get a faster number. HIGH reviews Decision Log and material evidence before integration; label self-review honestly.

## Acceptance

- The book contains runnable examples, recommended templates and measured fit/failure boundaries for the approved SELECT families, including same-plan outcomes and semantic counterexamples.
- Every performance number comes from a validated claim; toy outputs and conceptual illustrations are unmistakably labelled.
- Containerized evidence checks, PDF build and page-by-page visual review pass; navigation, checklist, manifest and HIGH review are recorded before local-main completion.

## Completion

Record RESULT.md, Decision Log, reports and actual limits. Submit for HIGH review; reconcile concurrent changes in your worktree, validate, fast-forward local main under its integration lock and verify reachability. Annotated milestone: `repo/book-select-query-templates-v1`.

NEXT MODEL: LOW
