# Execution Handoff — <task>

| Field | Value |
|---|---|
| Handoff ID / revision | <stable ID and revision> |
| Planner | <actual model, known effort, tool, date; HIGH> |
| Starting source / main revision | <verified commits, worktree and branch> |
| Checkpoint tag | <new annotated tag; never move it> |
| Status | <ready / amended / superseded / executed-awaiting-HIGH-review> |
| Next setting | **LOW — <agreed model/effort or user-selected role>** |
| Progress / escalations | <task-owned paths> |

## Objective and boundaries

<Outcome, owned paths, dependencies, completed work, untouched other-task artifacts.>

## Decisions made by HIGH

<Complete choices, reasoning where useful, permitted routine choices and explicit
non-goals. Specify measurements, environment, invariants and negative controls if relevant.>

## Ordered execution

| Step | Exact action or command | Expected evidence | Failure action |
|---|---|---|---|
| 1 | <create/reuse owned worktree and recheck live state> | <expected revision/status> | <mapped wait or ER> |
| 2 | <implementation/measurement> | <result> | <mapped action or ER> |
| 3 | <predefined validation> | <acceptance criteria> | <ER; never weaken a check> |

<Include path-safe, shell-specific commands and idempotency/restart behavior. Identify
locks, resource limits, runtime bounds and which external writes are authorized.>

## Decision Log

<LOW appends one short entry per meaningful decision here or in PROGRESS.md, whichever
the task already uses: decision, short reason, confidence (High/Medium/Low), alternative
considered, tasks potentially affected. Log only what could plausibly affect correctness,
architecture, performance, maintainability, security, behavior, compatibility, acceptance
criteria or overall quality -- not every implementation choice. This is what HIGH reviews
at validation instead of re-deriving everything LOW touched.>

## Escalation Required

<Concrete triggers -- decisions that could invalidate substantial downstream work,
significantly alter the architecture, materially change requirements, carry real
security/correctness/data-loss risk, are expensive to reverse, contradict a critical
HIGH assumption, or genuinely cannot be resolved from this handoff. Several reasonable
options, imperfect confidence, or an unspecified minor detail are NOT triggers -- LOW
decides, logs, and continues instead. What work must stop; continue independent mapped
steps. Specify the task escalation log; the executor does not rewrite this handoff.>

## Acceptance, integration and next role

<Required evidence, HIGH review gate, explicit-path commits, immutable run/milestone
tags, local-main integration procedure and context/lessons updates. No push or pull.>

**Next after execution: HIGH — review the Decision Log first, then validate the
result.**

## Amendments

<HIGH appends dated, attributed AM-NN decisions with superseded instructions identified;
never erase previous handoffs or another planner's text.>
