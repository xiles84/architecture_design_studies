# Goal — evidence-backed data architecture book and AI work pipeline

| Field | Value |
|---|---|
| Goal id | `goal-20260922T025912Z-data-architecture-book` |
| Requested at | `2026-09-22T02:59:12Z` |
| Baseline | `f27cb14a17791af77822be77b38915fa39811ba7` |
| Planner | model id not exposed, Codex desktop, user-selected HIGH capability |
| Status | active |
| Root task | `task-20260922T025912Z-publish-book-pipeline-handoff` |

## Objective

Turn the repository's measured database-design corpus into a beginner-to-master Typst
reference book, and make future work autonomously claimable, resumable, reviewable, and
traceable across HIGH and LOW AI sessions with minimal human steering.

The first edition is a **measured data architecture and state-design reference**, not a
claim to cover all software architecture. It separates design family, design variant,
technology, topology, and workload. Unsupported transfers are labelled analogy or gap.

## Decisions already made

- Typst source and a Podman-built PDF are the canonical assembled book.
- Human operational records are Markdown; queue and evidence machine state are JSON.
- Major model means design family; minor model means a controlled variant.
- Current evidence is publishable with explicit confidence and limitations.
- Numeric cross-study rankings are forbidden unless semantic and workload comparability
  is demonstrated.
- Existing results, reports, analyses, discussions, handoffs, and tags remain immutable.
- Queue exclusivity is local to one shared Git common directory in v1.
- HIGH capability and work role are separate; HIGH may execute LOW-eligible work.

## Acceptance

1. `docs/ai-work/` records the request, tasks, attempts, events, results, escalations,
   reviews, and parent/root correlations without relying on chat transcripts.
2. Atomic task claims, leases, recovery, role filters, a main-integration lock, and a
   portable committed event archive are implemented and tested.
3. Root and study navigation accurately describe all five studies; current context is
   concise and historical detail is preserved in an archive.
4. A machine-validated evidence registry identifies every book recommendation's strength,
   provenance, scope, limitations, and supersession.
5. The Typst book covers major families, minor variants, controlled cross-family evidence,
   topologies, common scenarios, mechanisms, evidence, and reproduction.
6. The PDF is built only in Podman and carries a build manifest tying it to source and
   evidence digests.
7. An independent HIGH release review accepts the book or creates targeted correction
   tasks. The accepted state is merged into local `main` and tagged
   `repo/data-architecture-book-v1` without pushing or pulling.

## Work partition

Tasks 1–3 and 5 are LOW-preferred execution. Tasks 4a–4c, 6, and 7 require HIGH. Future
scientific gaps are proposed planning tasks and do not authorize measurements. The task
records under `docs/ai-work/tasks/2026/09/` are the execution handoffs.

## Material escalation triggers

Escalate for a change to the taxonomy or evidence contract, a need to rewrite protected
artifacts, a queue design that cannot make claims atomic across this Git common directory,
an inability to reproduce claim provenance, a Typst dependency whose licence prevents
redistribution, or a scientific result that contradicts a planned book recommendation.
Routine implementation details remain Decide → Log → Continue.

## Publication checkpoint

This goal and its initial task records are the HIGH planning checkpoint. The checkpoint
is complete only when committed, tagged `repo/data-architecture-book-handoff-v1`, merged
into local `main`, and reachable there. The first task after it is queue v1.

NEXT MODEL: LOW
