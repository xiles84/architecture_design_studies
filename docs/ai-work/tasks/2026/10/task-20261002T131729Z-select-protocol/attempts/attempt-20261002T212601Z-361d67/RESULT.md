# Result — Plan SELECT query families, semantic contracts and crossover experiments

**Task:** `task-20261002T131729Z-select-protocol` · **claim:** `claim-3521080f3466b3b9` (epoch 1)
**Attempt:** `attempt-20261002T212601Z-361d67` · **worker:** deepseek-flash, Deep Code CLI, HIGH session, planner
**Starting main:** `e4059822` · **Result commit:** `ac7d89ce` on `study-01/select-protocol`

## Outcome

The premeasurement protocol and SQL/contract atlas are written and committed at
`studies/01-charity-tree/protocols/select-query-templates.md`. **No measurement code was written and
no benchmark was run** — the brief forbids measuring, and the four LOW measurement tasks stay
blocked on `select-harness` until this protocol is reviewed and integrated.

## Deliverable against the acceptance criteria

- **Predictions, answer groups, SQL contracts and falsifiers with revision/tag provenance** — §2
  records P1–P12 with expected direction, mechanism, expected crossover and a falsifier each, marked
  as written before any measurement code; §3 partitions the six result contracts (full rows, parent
  membership, optional match, aggregates, top-N, pagination); the header carries revision 1, the
  required tag `study-01/select-v0-protocol` and the starting main revision `e4059822`.
- **All eight task scopes bounded** — §12 gives executable matrices, controls, acceptance and
  escalation for `select-harness`, `join-cte`, `existence`, `aggregate-topn`, `access-paths`,
  `synthesis` and `book`; §8 covers sizing/placement (per-node and equal-total arms, balanced
  endpoints, verified colocation or gap).
- **Teaching examples with exact outputs** — §11 gives fixture A (the REQUEST example, six forms →
  exactly one row) and fixture B with exact outputs for NULLs, duplicate child matches, unmatched
  parents, ties and empty sets; it states explicitly that no number there is a measurement.

## Decision Log

1. **Deliverable path `studies/01-charity-tree/protocols/select-query-templates.md`; no new study
   number.** The SELECT work is an additive experiment inside Study 01 per PLAN.md §32–34.
   Confidence high.
2. **Recorded the REQUEST's six inner-join forms and the seventh EXISTS form separately, and forbade
   ranking EXISTS against full-row output.** The two answer *different* questions. Confidence high.
3. **Six result contracts are first-class; no comparison crosses contracts.** This is the main
   scientific decision that keeps a faster-but-different answer from being read as a win.
   Confidence high.
4. **Kept the REQUEST fixture A byte-faithful and added a separate fixture B for the edge cases**
   (NULL amount, duplicate 200 tie, unmatched parent Cara, empty sets) rather than mutating A.
   Confidence high.
5. **Marked engine feature support `unverified` and gated the dependent matrix rows.** CTE inlining
   defaults, `MATERIALIZED/NOT MATERIALIZED`, `DISTINCT ON`, `LATERAL`, `NOT IN` NULL semantics and
   snapshot support must be verified against the pinned PG 17 and YugabyteDB images in tasks 2–3, not
   recalled. Confidence high (PG 17 changed the one-use CTE default; YB is independent).
6. **Rollup/rolldown/embedding: justified non-applicability for the syntax-only isolate pairs, mapped
   to existing designs** (`d4/d5/d6/d9/d10/d11`), with a rule that a physical change must be a
   separate pair that the syntax arm may not claim. Confidence high.
7. **Optimistic/pessimistic concurrency recorded non-applicable** (arbitration unchanged); a varied
   writer would require both strategies and is a different study. Confidence high.
8. **Topology: PG single + YB single required for every measuring task; YB 3-node RF3 only where
   topology-relevant, otherwise an explicit `coverage_gap`;** a PG cluster/read-routing result is not
   implied by multiple containers. Confidence medium-high (bounded by the fixed 2 GiB host).
9. **≥5 levels for every claimed crossover; ≥3 runs and ≥3 trials for any argued difference; seeded
   order and a first/last position control.** Confidence high.
10. **Prediction record-at-commit:** the introducing protocol commit is the pre-registration point;
    its hash is recorded here rather than re-edited into the table, to avoid rewriting a
    pre-registration record. Confidence high.

## Validation

- No measurement, no database start, no run lock needed (planner task).
- Acceptance self-check against the three criteria above: met.
- Files changed: `studies/01-charity-tree/protocols/select-query-templates.md` (new), `CONTEXT.md`
  (select-protocol entry). Explicit paths staged; no other analyst's files touched.

## Limits and gaps recorded

- Feature support remains `unverified` until tasks 2–3 check the pinned images.
- The PG-cluster/read-routing and any unsupported YB 3-node combination are explicit gaps.
- All fixture outputs are logical, not measured; the atlas invents no performance number.

## Remaining for HIGH review

Review the Decision Log and the acceptance evidence (self-review is required since no independent
reviewer is available), then integrate `study-01/select-protocol` into `main`, create the annotated
tag `study-01/select-v0-protocol`, complete the task and remove `.worktrees/select-protocol`.
