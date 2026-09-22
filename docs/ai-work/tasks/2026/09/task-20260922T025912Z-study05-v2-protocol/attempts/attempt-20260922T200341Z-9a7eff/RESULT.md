# Result — Study 05 v2 completion protocol

**Task:** `task-20260922T025912Z-study05-v2-protocol`
**Capability/role:** HIGH planner (model `deepseek-flash`)
**Required tag:** `study-05/v4-handoff`
**Measured:** nothing.

## Delivered

- `studies/05-cache-consistency/HANDOFF-AMENDMENT-01-V2.md`: real-TTL churn, medium scale and a
  tested ceiling, ≥3 repeated randomized trials, separately labelled equal-total cache resource
  accounting, YugabyteDB/cluster/verified placement, open-loop demand, and the v2 report acceptance
  criteria.
- `HANDOFF.md` points to the amendment.
- Three separately claimable execution tasks published as `proposed`:
  - `task-20260922T200831Z-study05-v2-churn-scale`
  - `task-20260922T200831Z-study05-v2-resources-engines`
  - `task-20260922T200831Z-study05-v2-open-loop`
- `CONTEXT.md` gaps 4 and 5 updated.

## Acceptance criteria check

| Criterion | Result |
|---|---|
| Covers real-TTL churn, medium scale, repeats, cache resource accounting, YugabyteDB/cluster/placement, open-loop demand | met (amendment sections 1–6) |
| Creates separately claimable execution tasks | met (three published) |
| Does not measure | met |

## Notes

The book's cache claims (`v2-14`, `v2-15`, `v2-16`) remain scoped to `db-only`, single-trial,
800-donor evidence until these tasks run; the amendment is the path to widening them.
