# SELECT query templates: examples, measurements and selection guidance

Status: protocol/harness and three measurement families integrated; join/CTE remains ready and unmeasured. See the dated amendment below for the current execution route.
Owner request: [REQUEST.md](REQUEST.md). Existing goal: goal-20260922T025912Z-data-architecture-book.
Planner: GPT-6 / Codex desktop / effort unknown; user-declared HIGH.
Publication task: task-20261002T131729Z-select-query-task-publication.
Starting local main: 1a178ce; publication milestone: repo/select-query-task-program-v1.

## Current owner amendment — 2026-10-04

The theory now includes six join/filter spellings, outer-join semantics and CTE
policies. The [join handoff](JOIN-HANDOFF-20261004.md) resolves the yielded LOW
plan and controls its matrix: LOW measures in the existing Study 01 task; HIGH
reviews through the existing synthesis task; HIGH then publishes the case results.
New `task-20261004T142623Z-select-book-high` supersedes the original LOW book
assignment. Its predecessor amendment is mandatory even though queue v1 cannot
terminalize a superseded task through a command. No new measurement was run.

Current route: `select-join-cte` (LOW) -> `select-synthesis` (HIGH) ->
`select-book-high` (HIGH). The completed protocol/harness need no restart.

## Reader outcome

A reader starts with the required answer, chooses a readable SQL template and can see the measured conditions in which an alternative helps, hurts or produces no resolved difference. Correctness, output grain, NULL handling and ordering precede speed. No syntax is sold as an unconditional optimization. A crossover that never appears is a valid outcome.

Use a small customers/orders teaching fixture, then connect it explicitly to equivalent operations in Study 01's existing charity/person/donation domain. Toy data/output are examples. Only registered observed cells supply performance numbers.

## Work graph

| Task | Capability / role | Dependency | Deliverable |
|---|---|---|---|
| 1. [Plan SELECT query families, semantic contracts and crossover experiments](../ai-work/tasks/2026/10/task-20261002T131729Z-select-protocol/BRIEF.md) | HIGH / planner | publication complete | studies/01-charity-tree/protocols/select-query-templates.md |
| 2. [Build the containerized SELECT harness, oracle and readable plan reports](../ai-work/tasks/2026/10/task-20261002T131729Z-select-harness/BRIEF.md) | LOW / executor | 1 | additive SELECT mode and embedded SQL catalogue |
| 3. [Measure join/filter spellings and CTE materialization success and failure regimes](../ai-work/tasks/2026/10/task-20261002T131729Z-select-join-cte/BRIEF.md) | LOW / executor | 2 | join/CTE run manifests, results, plans and generated reports |
| 4. [Measure EXISTS, IN, JOIN DISTINCT and anti-join templates](../ai-work/tasks/2026/10/task-20261002T131729Z-select-existence/BRIEF.md) | LOW / executor | 2 | membership/absence manifests, results, plans and generated reports |
| 5. [Measure join aggregation, correlated subqueries and per-parent top-N](../ai-work/tasks/2026/10/task-20261002T131729Z-select-aggregate-topn/BRIEF.md) | LOW / executor | 2 | aggregate/top-N manifests, results, plans and generated reports |
| 6. [Measure index-friendly predicates, OR alternatives and pagination boundaries](../ai-work/tasks/2026/10/task-20261002T131729Z-select-access-paths/BRIEF.md) | LOW / executor | 2 | access-path/pagination manifests, results, plans and generated reports |
| 7. [Review SELECT evidence and publish measured query-template selection guidance](../ai-work/tasks/2026/10/task-20261002T131729Z-select-synthesis/BRIEF.md) | HIGH / analyst | 3, 4, 5, 6 | new signed SELECT analysis and mechanism companions |
| 8. [Publish reviewed SELECT cases using HIGH](../ai-work/tasks/2026/10/task-20261004T142623Z-select-book-high/BRIEF.md) | HIGH / executor | 7 + Edition 3 pilot + theory update | reviewed measured case variations, selection guide and rendered book |

Task ids share prefix task-20261002T131729Z. Suffixes: select-protocol, select-harness, select-join-cte, select-existence, select-aggregate-topn, select-access-paths, select-synthesis, select-book.
Publication is complete only after audit/review, local-main integration and the annotated milestone.
The first downstream task is HIGH protocol work. LOW workers may prepare/measure only after its committed decisions and the harness dependencies complete. Measurement tasks can be developed in independent worktrees, but the machine's measurement slot remains serial.

## Fixed scientific and editorial decisions

- Extend Study 01 with an additive SELECT experiment and distinct image/output namespaces. No new study number, migration of its historical harness, or mutation of old signed evidence. Final implementation details belong in the protocol.
- Query formulation, physical indexes, information placement and resource/topology changes are independent factors. Equivalent syntax pairs hold everything else fixed. Crossed factors are interaction evidence; realistic combined changes are package evidence.
- Partition answers into bags of joined rows, sets of parent identities, parent-preserving optional matches, aggregates, global/per-parent top-N, and ordered pagination. Comparisons across these contracts do not establish a faster implementation.
- SQL logical structure, optimizer transformation, observed plan structure and measured end-to-end behavior are four distinct facts. Query structure does not establish execution order. Plans accompany timing.
- Prefer ordinary explicit joins for the initial full-row example, EXISTS for the parent-membership example, and preserve NULL/duplicate/tie/snapshot assumptions. These are readable starting templates, not measured fastest choices.
- Reuse the concept-first Edition 3 layout: short chapters, runnable SQL, diagrams/plan explanations, fits/avoid rows, measured case boxes, glossary and a compact selection guide. Do not turn it into a long syntax catalogue or add reader exercises.
- Existing rollup/rolldown/embedding cases remain the source for information-placement costs. New query syntax never earns an attribution for a changed index or data model; persistent materialized-read-model expansion belongs to Study 08 and hierarchy recursion to Study 09.
- This program concerns an existing relational domain and does not release deferred Study 06–09 tasks or supersede the replication-first program. Priority 54/55 sits below existing T1 workers at 62; new domain/topology expansion still follows its gate.
- Build/review the new book material only after the approved Edition 3 pilot is integrated. Preserve its outstanding owner approval gate; this task request is not approval of another task's style preview.

## Query family coverage

| Family | Valid controlled comparisons | Important trap / rejection |
|---|---|---|
| Inner joins and filtering | all six supplied spellings, default CTE, supported MATERIALIZED/NOT MATERIALIZED, selective one-use vs expensive repeated reuse | written order is not physical order; CTE feature support is engine/version-specific |
| Optional child matches | LEFT JOIN with child filter in ON vs filtered child input | WHERE on child may drop an unmatched parent: a different answer, not a speed win |
| Membership | EXISTS, IN, JOIN plus parent-id DISTINCT | plain join emits one row per qualifying child; DISTINCT on names changes identity |
| Absence | NOT EXISTS, non-null child-id LEFT JOIN IS NULL, non-null-contract NOT IN | nullable NOT IN can produce UNKNOWN; NULL and empty inputs must be exercised |
| Parent aggregates | post-join aggregate, pre-aggregation, correlated scalar/lateral aggregate | two child collections multiply totals; COUNT/SUM/empty-parent semantics differ |
| Latest / top-N per parent | lateral/correlated indexed probe, matching ROW_NUMBER, supported DISTINCT ON for N=1 | global limit is another question; ROW_NUMBER vs RANK differs on ties |
| Predicates and set alternatives | matched range/function/cast predicates, OR vs overlap-safe UNION/UNION ALL | timezone/type/coercion or overlap changes the answer; new index is a separate factor |
| Pagination and physical access | same ordered page under fixed snapshot, OFFSET vs known-cursor keyset; matched index arms | arbitrary-page and next-page contracts differ; mutable offsets shift pages |

## What numbers must establish

Tasks 3–6 follow a precommitted protocol, not a universal performance ranking:
- Valid-cell read throughput, end-to-end p50/p95/p99 with sample counts and applicable tail limits, per-trial and across-run spread; read/write throughput in a matched blended phase.
- Planning/execution and fully drained result row/byte counts, estimated/actual cardinalities, loops, buffers, supported RPC counters, memory/temp spills, index/storage and write-maintenance cost where a physical design changes.
- At least three independent runs and three trials for an argued difference; seeded randomized/interleaved order, repeated-design first/last control. Differences below measured variability remain unresolved.
- At least five levels for a claimed boundary sweep. Candidate axes include selectivity (0.01%, 0.1%, 1%, 10%, 50%, 90%), parent selection, child fanout/skew, row width, memory, reuse count, N and page depth. These are proposed input levels, never results; HIGH chooses/calibrates a feasible bounded matrix before measuring code.
- Each family reports best fit, tested adverse regime, equality/no-resolved-difference, invalid answer, timeout or unsupported gap as applicable. A reversal needs observations on both sides and a bounded crossing interval with spread; absence of a reversal is reported.
- Every cell retains run tag, producing commit, SQL/data/settings/image identifiers, environment and digest. Generated reports contain measurements; signed analyses contain recommendations.

## Correctness and workload contract

Use an independent Go oracle from deterministic logical data. Compare bag multiplicities as well as values, parent membership as identities, and ordered outputs under a declared total order. Include duplicates, NULLs, unmatched parents, ties, empty sets and overlapping predicates. Read-path deliberately wrong fixtures must actually fail their intended oracle per relevant engine/topology.

Concurrent reads use a defined snapshot/as-of/operation-ledger contract compatible with that oracle. An unversioned earlier expected answer cannot prove a changing database correct. Hold the writer implementation fixed; this is not a new concurrency strategy study. If strategy is varied, optimistic/pessimistic coverage and invariant controls become required.

## Resource and topology coverage

The protocol calculates current host and Podman VM resources, client/system reserves, node/pool/worker budgets and attainable load. It does not recall laptop sizing from old reports. Include PG single-node, YB single-node and YB 3-node RF3 with distinct per-node and equal-total arms when supported. Record actual endpoint distribution and per-node/client throttling.

A PG replication/read-routing experiment is not implied by multiple PostgreSQL containers. Unsupported engine/topology features remain explicit gaps. Representative cluster subsets may keep runtime bounded only if protocol-selected, labelled and all omitted combinations remain visible.

Physical colocated/non-colocated placement must be verified at the tablet/partition/leader level with identical logical data, operations, replication and resources; laptop container proximity proves nothing about data placement. Single-host network and hardware-independent conclusions remain gaps. All DB checks/calibration/runs use ads-run-lock; never clear another agent's lock or modify its inputs.

## Claims and publication

HIGH reviews worker decisions and adequacy, issues targeted corrections only for material defects, reads analyses already attached to each inputs digest, and publishes its own signed analysis. Every prediction gets a verdict; every book number must semantically resolve to an immutable valid source cell in the active evidence package. Follow its current successor/versioning policy without editing frozen versions.

The final HIGH book task uses approved claims and templates, runs the pinned Podman evidence checks/build, renders and visually inspects new pages, records figure text equivalents/two-reader strength and completes HIGH review. It updates current source/navigation, tracked PDF/manifest and changelog, merges into local main and tags completion. No benchmark number exists merely because these tasks were created.

## Deferred choices and escalation

HIGH protocol work decides feature support, exact SQL, bounded matrix/runtime, calibration, admissible memory/index/statistics/prepared-plan controls and snapshot implementation. This publication fixes coverage and contracts, not an untested runnable measurement plan. LOW may choose local reversible implementation details and log material choices; changes to semantics, hypotheses, scientific scope, accepted coverage or evidence gates require HIGH.

## Next action

Claim task-20261002T131729Z-select-join-cte with declared LOW capability and executor role; read the dated handoff and predecessor amendments, then resume its canonical worktree. No measurement task or external chat was started by this publication.

NEXT MODEL: LOW
