# Review — Study 09 hierarchy protocol

**Task:** `task-20260922T025912Z-hierarchy-model-protocol`
**Reviewer:** HIGH session, model `deepseek-flash`, role `reviewer`
**Decision:** **approve.** Decision-complete; measured nothing.

## Decision Log reviewed

1. Five representations, one operation set, cross-representation equality — accept.
2. Move cost as a function of size/depth — accept; it is the discriminating axis.
3. Bounded embedding carries its bound and overflow rule — accept.
4. Cycle/orphan/depth-bound controls must fire; concurrency adds lost-update — accept.
5. Storage at three tree sizes — accept.
6. PostgreSQL v0, YugabyteDB in the concurrency task, engines separate — accept.
7. Four tasks, harness first — accept.

## Acceptance criteria

| Criterion | Verdict |
|---|---|
| Five representations under identical operations | met (§2–3) |
| Covers reads, moves, inserts, deletes, storage, concurrency, correctness | met (§4–5) |
| Creates execution tasks but does not measure | met (four published) |

No material defect. Integrate and tag `repo/hierarchy-study-handoff-v1`.
