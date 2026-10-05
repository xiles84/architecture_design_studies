# HIGH: review the T1/T2 results, decide the gate, and release Study 09 if closed

**Task:** `task-20261005T151500Z-t5-gate-review` · **capability:** HIGH · **role:** reviewer · **benchmark:** none
Depends on the six tasks of [`docs/gate-closure-20261005/PLAN.md`](../../../../../gate-closure-20261005/PLAN.md).

## Goal

Decide, from the integrated T1/T2 evidence, whether the charter's "T1 and T2 before T5" gate is closed, and if so release Study 09 (`task-20260922T204845Z-hierarchy-harness` first, then its three measurement tasks).

## Acceptance criteria

- Each of the six T1/T2 results was reviewed and integrated (or has a recorded reason it could not be)
- A signed gate decision states, claim by claim, whether T1 has no remaining single-run headline claim and whether each T2 sweep delivered a five-level verdict
- If the gate is closed, the Study 09 hierarchy tasks are released with publish --release in dependency order (harness first) and the program file records the decision; if not, the decision names the exact missing evidence
- CONTEXT.md and LESSONS_LEARNED.md updated; merged into local main and tagged

## Constraints

Do not release a deferred task (Study 06, 07, 08) by this decision; their blockers are resources, not this gate. A sweep that observed no reversal still delivers a verdict. A claim left `single_run_directional` keeps the gate open unless the charter's exception (the domain task is the only way to close a named hole) is argued explicitly.

NEXT MODEL: HIGH
