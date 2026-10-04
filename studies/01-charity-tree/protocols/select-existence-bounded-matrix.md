# Bounded measurement plan — membership and absence (`select-existence`)

LOW worker plan for the reviewed protocol's task-4 row
(`protocols/select-query-templates.md`, P6/P7). It fixes the matrix, controls,
levels and replication before any database run.

## 1. Question and predictions

- **P6** membership: `EXISTS` vs `IN` vs `JOIN + DISTINCT` on parent identity.
  Predicted: `EXISTS`/`IN` avoid row multiplication; `DISTINCT` pays a
  de-duplication cost that grows with child fanout.
  Falsifier: `DISTINCT` is never slower at high fanout, or `EXISTS`/`IN` differ
  in answer from the membership set.
- **P7** absence: `NOT EXISTS` vs `LEFT JOIN ... IS NULL` vs `NOT IN`.
  Predicted: `NOT EXISTS`/anti-join are safe; nullable `NOT IN` returns nothing.
  Falsifier: `NOT IN` over a nullable subquery returns the same set as
  `NOT EXISTS`.

## 2. Data and axes

- **Fanout sweep, constant total children.** Membership is compared at child
  fanout {1, 10, 100} with the total number of child rows held at 10,000, so
  the parent count is {10000, 1000, 100}. One deterministic fixture per level
  (seed 42); each level's fingerprint and row counts are recorded.
- **Membership spellings:** `EXISTS` (s06), `IN` over a nullable subquery (s07),
  `IN` over a NOT-NULL-filtered subquery (e12), `JOIN + DISTINCT` (s08).
- **Absence spellings:** `NOT EXISTS` (s09), `LEFT JOIN ... IS NULL` (s10),
  `NOT IN` over a NOT-NULL-filtered subquery (e13). The nullable-`NOT IN` case is
  the labelled semantic counterexample (`x02_not_in_nullable`), never ranked.
- The `{NULL, NOT NULL}` axis is the subquery contract: the valid `NOT IN` arm
  filters `customer_id IS NOT NULL`; the NULL arm is the control.
- Predicate (`amount > 100`) is identical across spellings in a pair, so a pair
  isolates the spelling. Qualifying proportion, hit fraction and skew are
  recorded as coverage gaps for this bounded run.

## 3. Controls (must fire)

- `x02_not_in_nullable` — nullable `NOT IN` returns the empty set; already live
  in the shared gate and detected on the nullable fixture.
- `x07_membership_multiplier` — plain `JOIN` without `DISTINCT` emits one row per
  qualifying child, so it differs from the membership set of parent ids. It is
  the ordinary child-row join the brief says is illustrated separately; it must
  be detected, never timed as a membership spelling.
- Both run on the measurement fixture each level; a non-firing control fails the
  run.

## 4. Arms, replication and evidence

- PostgreSQL 17.11 single-node: ≥3 whole runs (one with the position control).
- YugabyteDB 2025.2.6 single-node: 3 whole runs.
- 5 trials per cell after oracle verification; `EXPLAIN ANALYZE` kept per cell
  (raw plan text; structure signatures for the report).
- Reports: per run and replicated per engine, with membership trio ratios per
  fanout level, absence cells, controls, plans and limits. Verdicts in a signed
  analysis beside the report.

## 5. Explicit gaps

No 3-node/equal-total-budget arm, no verified placement, no concurrent demand,
no qualifying-proportion/hit-fraction/skew axis in this bounded run, no index
redesign arm, and YugabyteDB RPC/storage counters are not parsed.
