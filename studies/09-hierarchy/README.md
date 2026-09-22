# Study 09 — hierarchy representation

**Question:** for one hierarchy, how do adjacency list, materialized path, closure table, nested sets
and bounded embedding compare under identical operations?

**Status:** protocol only; **no measurement yet**. The decision-complete handoff is
[`HANDOFF.md`](HANDOFF.md) (tag `repo/hierarchy-study-handoff-v1`).

**Operations:** ancestors, subtree, children, path, depth, insert (leaf/subtree), move-subtree,
delete (leaf/subtree), reparent-and-recount, per-level cardinality — the same for every
representation, with cross-representation answer equality as the correctness gate.

**Measures:** reads, moves (as a function of subtree size and depth), inserts/deletes, storage at
three tree sizes, and concurrency (contended move, read-during-move, lost-update control).

**Boundary:** v0 is PostgreSQL single-node with a YugabyteDB coverage task; equal-total and per-node
resource arms are separate; no cross-representation ranking across runs.
