# Result — Study 09 hierarchy protocol

**Task:** `task-20260922T025912Z-hierarchy-model-protocol`
**Capability/role:** HIGH planner (model `deepseek-flash`)
**Required tag:** `repo/hierarchy-study-handoff-v1`
**Measured:** nothing.

## Delivered

- `studies/09-hierarchy/HANDOFF.md` and `README.md`: a decision-complete protocol comparing
  adjacency list, materialized path, closure table, nested sets and bounded embedding under one
  operation set, with reads, moves, inserts, deletes, storage and concurrency, and cycle/orphan/
  depth-bound controls.
- Four separately claimable execution tasks published as `proposed`:
  `task-20260922T204845Z-hierarchy-harness`, `…-reads-storage`, `…-writes-moves`,
  `…-concurrency-engine`.
- `CONTEXT.md` open gap 9 records completion and the new tasks.

## Acceptance criteria check

| Criterion | Result |
|---|---|
| Compares adjacency list, materialized path, closure table, nested sets, bounded embedding under identical operations | met (handoff §2–3) |
| Covers reads, moves, inserts, deletes, storage, concurrency, correctness | met (§4–5) |
| Creates execution tasks but does not measure | met (four published) |

## Notes

No hierarchy claim exists in the active evidence registry; this study creates the evidence base for
one if a future edition needs it.
