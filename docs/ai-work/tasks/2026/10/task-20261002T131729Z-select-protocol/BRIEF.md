# Plan SELECT query families, semantic contracts and crossover experiments

Capability: HIGH; role: planner. Dependencies: `task-20261002T131729Z-select-query-task-publication`.

## Goal and fixed scope

Read the owner's request and this program, current Study 01 SQL/harness, the empirical-method program, book/STYLE.md and the active evidence registry. Extend Study 01's existing charity/person/donation domain rather than allocate a new study or bypass the replication-first ordering. The customers/orders fixture is a teaching analogue, not existing measured evidence.

Produce a protocol from docs/templates/STUDY_PROTOCOL.md plus an SQL example/contract atlas. Include all seven forms in REQUEST.md and the further query families in PLAN.md. Partition equivalent full joined rows, parent membership, parent aggregates, parent-preserving optional matches, per-parent top-N and pagination into separate result contracts.

Write falsifiable predictions BEFORE measuring code: simple inner-join spelling will often share a plan; a materialization boundary can inhibit selective access but reusable expensive work can benefit; membership methods can avoid row multiplication; per-parent indexed probes can win at small parent sets and lose at large sets; broad matches can reverse index lookup benefits. These are hypotheses to test, never assumed findings.

Decide a bounded matrix and runtime budget, at least five levels for a proposed crossover, three independent runs and at least three within-run trials, seeded interleaving and repeated-design position controls. Specify metrics, timeout rules, snapshots, row drainage, random keys, result sizes and plan-signature normalization.

Calculate LIVE host/VM, client and system reserves, per-node and equal-total arms, readers/writers/pools, balanced endpoints and measured throttling. Cover PG single-node, YB single-node and YB 3-node RF3; unsupported PG cluster execution is an explicit gap, not RF3 evidence. Require verified colocated/non-colocated pairs where supported.

Map required rollup/rolldown/embedding to existing designs or a justified non-applicability for syntax-only isolate pairs; if varied, use separate pairs with maintenance/storage/correctness. Mixed writers use a fixed correct write contract. Optimistic/pessimistic concurrency is non-applicable if arbitration is unchanged; varying it requires both.

Resolve feature support against PINNED engines, statistics/prepared plan policy, documentation and existing code. Freeze a compact LOW handoff for each measuring task; material narrowing or semantic changes need a dated amendment and HIGH review. Do not measure.

## Shared constraints

Read the durable REQUEST.md/PLAN.md, AGENTS.md, current CONTEXT.md and methodology. Claim/guard atomically, work in your own task worktree, stage explicit owned paths, preserve prior signed artifacts and immutable evidence, identify actual model/tool/effort, integrate current main and complete with an annotated tag. No push/pull.
All runtime work is in Podman. Every check/calibration has a report and every database-starting operation holds ads-run-lock. Measurement starts from committed/tagged code. Drain result rows; randomize keys; gate valid timings on independent Go truth; make read-path controls fire per relevant topology/phase. Use >=3 independent whole runs and >=3 trials for argued deltas, seeded interleaving and repeated-design position controls, >=5 levels for a crossover, isolated and blended workloads, actual demand/client caps, per-node throttling and EXPLAIN ANALYZE cost counters.
Follow the reviewed bounded protocol: identical logical data/answer/indexes/settings in query-form isolate pairs; declared index variants include mutation/storage cost; single-node and 3-node per-node/equal-total arms separate; balanced endpoints and actual colocated/non-colocated placement evidence or an explicit gap. One laptop establishes neither independent hardware nor real network effects.
Decide/log/continue local reversible choices. Escalate changing semantics, question, predictions, accepted matrix, oracle/negative control, required evidence or book style that would invalidate downstream acceptance. Resource-infeasible/unsupported/timed-out cells remain visible gaps/results. Never weaken a gate or substitute a different answer to get a faster number. HIGH reviews Decision Log and material evidence before integration; label self-review honestly.

## Acceptance

- Premeasurement protocol records predictions, equivalent answer groups, SQL contracts and falsifiers with revision/tag provenance.
- All eight task scopes receive bounded, executable matrices, sizing/placement coverage, controls, acceptance and escalation criteria.
- Teaching example SQL includes exact outputs for NULLs, duplicate child matches, unmatched parents, ties and empty sets; no performance claims are invented.

## Completion

Record RESULT.md, Decision Log, reports and actual limits. Submit for HIGH review; reconcile concurrent changes in your worktree, validate, fast-forward local main under its integration lock and verify reachability. Annotated milestone: `study-01/select-v0-protocol`.

NEXT MODEL: HIGH
