# Execution Handoff — Study 09, hierarchy representation

**Status:** decision-complete protocol; **does not measure**. Authorised by
`task-20260922T025912Z-hierarchy-model-protocol` (HIGH planning).
**Required tag on integration:** `repo/hierarchy-study-handoff-v1`.
**Numbering:** Studies 06–08 exist; this is Study 09.

## 1. The question

For one hierarchy — a product/category tree, an organisation chart, a bill of materials — how do the
five classic representations compare under identical operations? Adjacency list, materialized path,
closure table, nested sets, and bounded embedding (a bounded subtree stored inside its parent), all
with the same logical data, the same reads and the same correctness contract.

## 2. Representations

| Representation | Shape | The bet |
|---|---|---|
| Adjacency list | each row stores its parent id | writes are trivial, subtree reads recurse |
| Materialized path | each row stores its full path (`/a/b/c`) | subtree reads are a prefix scan, moves rewrite descendants |
| Closure table | a separate (ancestor, descendant, depth) table | arbitrary reads are joins, every write touches the table |
| Nested sets | left/right bounds | fast subtree reads, moves renumber |
| Bounded embedding | children embedded in the parent document, depth/size bounded | one read for the whole subtree, moves are document rewrites, overflow must have a rule |

Bounded embedding is measured with an explicit bound (e.g. ≤32 children, depth ≤4); the overflow
rule (re-split, or store a reference) is part of the design, not an implementation detail.

## 3. Operations (identical for every representation)

`root-ancestors`, `subtree`, `children`, `path-to-root`, `depth`, `insert-leaf`, `insert-subtree`,
`move-subtree` (including across parents), `delete-leaf`, `delete-subtree`, `reparent-and-recount`,
and a cardinality read (node counts per level). Every representation answers the same logical
questions; the correctness gate compares answers across representations, not just within one.

## 4. Measurement axes

- **Reads**: per-operation throughput and p50/p99 on a fixed tree (balanced) and a skewed tree.
- **Moves**: move-subtree cost as a function of subtree size and depth — the axis that separates the
  representations.
- **Inserts and deletes**: leaf and subtree, batched and single.
- **Storage**: row/column bytes and index bytes per representation, at three tree sizes.
- **Concurrency**: concurrent inserts under different parents, concurrent moves of the same subtree
  (the contended case), and a read during a move.
- **Correctness**: cross-representation answer equality, plus cycles (a move that would create one
  must be refused), orphan prevention, and depth-bound enforcement for embedding.

## 5. Correctness and negative controls

Every design carries: a cycle control (a move that creates a cycle must be refused), an
orphan control (a delete that would orphan children must be refused or cascade, as declared), and a
cross-representation answer check. A concurrency arm adds a lost-update control. A design whose
control did not fire is not evidence.

## 6. Resources and topology controls

Per-node 2 CPU / 3 GiB and client 2 CPU / 2 GiB as today; a separately labelled equal-total arm; v0 is
PostgreSQL single-node, with YugabyteDB coverage in a second task. Tree sizes 1k / 100k / 1M nodes
with a deterministic generator (seed 42) and a recorded skew.

## 7. Deliverables and separately claimable execution tasks

1. **Hierarchy harness and correctness gate** — generator, cross-representation checker,
   cycle/orphan controls, loaders for all five representations.
2. **Hierarchy reads and storage** — the read and storage axes on balanced and skewed trees.
3. **Hierarchy writes and moves** — insert/delete/move-subtree including the move-size sweep.
4. **Hierarchy concurrency and second engine** — contended moves, read-during-move, and the
   PostgreSQL/YugabyteDB comparison under the same correctness gate.

Each task takes the benchmark lock, runs one matrix at a time, and produces a generated report plus a
signed analysis under `studies/09-hierarchy/`. No cross-representation ranking across different runs.

## 8. Acceptance criteria for the study report

- All five representations implement the same operations and return equal answers.
- Move cost reported as a function of subtree size and depth, not one number.
- Storage reported per representation at three tree sizes.
- Cycle, orphan and (for embedding) depth-bound controls fire.
- Concurrency results include the contested-move and read-during-move cases with a lost-update
  control.
- Per-node and equal-total arms labelled and separate; engines reported separately.
- Every number resolved to a v2 evidence claim before it enters the book.
