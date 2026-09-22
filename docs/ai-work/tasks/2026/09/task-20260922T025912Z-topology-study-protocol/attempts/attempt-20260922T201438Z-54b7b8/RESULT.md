# Result — Study 07 topology protocol

**Task:** `task-20260922T025912Z-topology-study-protocol`
**Capability/role:** HIGH planner (model `deepseek-flash`)
**Required tag:** `repo/topology-study-handoff-v1`
**Measured:** nothing.

## Delivered

- `studies/07-topology/HANDOFF.md` and `README.md`: a decision-complete protocol that separates node
  count, replication, placement, routing, resource budget, network and failure; requires independent
  hosts, a real network path and balanced endpoints with a recorded distribution; and defines the
  correctness and negative controls per axis.
- Four separately claimable execution tasks published as `proposed`:
  `task-20260922T202001Z-topology-env-harness`, `…-node-replication`, `…-placement-routing`,
  `…-network-failure`.
- `CONTEXT.md` open gap 6 records completion and the new tasks.

## Acceptance criteria check

| Criterion | Result |
|---|---|
| Separates node count, replication, placement, routing, resource budget, networking, failures | met (handoff §4) |
| Uses independent hosts and balanced endpoints | required by §2–3 and gated by execution task 1 |
| Defines correctness and negative controls | met (§5) |
| Creates separately claimable execution tasks | met (four published) |
| Does not measure | met |

## Notes

This protocol is the path to close `v2-gap-05` (colocation unverified) and `v2-gap-07` (real network
/ balanced endpoints). No runner exists yet; task 1 builds it.
