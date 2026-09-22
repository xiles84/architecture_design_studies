# Progress — Study 09 hierarchy protocol

**Task:** `task-20260922T025912Z-hierarchy-model-protocol`
**Attempt:** `attempt-20260922T204712Z-105ba6` (claim `claim-e23f78318d7f0af8`, epoch 1)
**Capability/role:** HIGH planner (model `deepseek-flash`)
**Required tag:** `repo/hierarchy-study-handoff-v1`

Planning only: no database started, no measurement.

## Decision Log

1. **Five representations, one operation set, cross-representation equality** as the correctness
   gate — not a within-design check. Confidence: High.
2. **Move cost is a function, not a number.** The move-subtree axis sweeps subtree size and depth;
   that is the axis that separates the representations. Confidence: High.
3. **Bounded embedding carries its bound and overflow rule** in the design. Confidence: High.
4. **Cycle, orphan and depth-bound controls must fire**; a concurrency arm adds a lost-update
   control. Confidence: High.
5. **Storage at three tree sizes**, so a representation is not judged at one cached size.
   Confidence: High.
6. **PostgreSQL v0, YugabyteDB coverage in the concurrency task**, engines reported separately.
   Confidence: High.
7. **Four execution tasks, harness first.** Confidence: High.

## Residual risks

- A 1M-node nested-sets tree may be impractical to renumber; the writes task reports the ceiling it
  reached rather than omitting the representation.
- Cross-representation equality may disagree on tie-breaking order; the checker must compare sets,
  not iteration order.
