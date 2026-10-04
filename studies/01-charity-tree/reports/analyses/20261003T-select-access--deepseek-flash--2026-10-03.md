---
analysis_id: 20261003T-select-access--deepseek-flash--2026-10-03
run_id: select-access-20261003 (pg-1/2/3, pg-pos, yb-1/2/3)
environment: xiles-podman-select-access-20261003
analyst: deepseek-flash
analyst_kind: ai
analyst_version: "DeepSeek deepseek-flash via opencode CLI, effort not exposed"
analyzed_at: 2026-10-03
inputs_digest: 0efeae8e02bd0379cd50fa2a863cc4c1ab9086cc6cfc8fdafb1c96d0f55c8aa0
repo_commit: 524ce328afa3e738f095b71214051240bcc6f1f0
supersedes:
status: current
headline: "Keyset pagination overtakes OFFSET at a depth that is engine-specific (PG ~1k, YB ~10k); OR is never worse than UNION/UNION ALL, and a new expression index changes the expression plan only on PostgreSQL."
---

# Analysis — Study 01 SELECT access paths (P10, P11) — deepseek-flash

Signed analysis of the access-path measurement family for the additive SELECT
programme (task `task-20261002T131729Z-select-access-paths`). Measurements are in
the per-run reports `results/select-access-{pg,yb}-*/access.md` and the generated
cross-run report [`20261003T-select-access-replicated.md`](../20261003T-select-access-replicated.md).
Protocol predictions are P10 and P11 in
[`protocols/select-query-templates.md`](../../protocols/select-query-templates.md);
the concrete bounded matrix is in
[`protocols/select-access-paths-bounded-matrix.md`](../../protocols/select-access-paths-bounded-matrix.md).

## TL;DR

- **P11 is confirmed, with an engine-specific crossing.** Keyset pagination beats
  OFFSET only past a depth: on PostgreSQL the crossover is between depth 100 and
  1,000 (at depth 10,000 keyset is 4.3–5.1× faster, 4/4 runs), on YugabyteDB it is
  between 1,000 and 10,000 (at 10,000 keyset is 1.5–1.9× faster, 3/3 runs). Below
  the crossing the two are equal or OFFSET is slightly ahead. Telling a reader to
  "use keyset" without the depth and the engine is wrong in both directions.
- **P10 is falsified as stated.** No shape or regime showed `UNION ALL` using two
  index arms to beat `OR`; instead `OR` is never worse by a replicated margin.
  On YugabyteDB it is consistently better than both `UNION` and `UNION ALL`
  (1.13–2.12×, 36/36 shape×run ratios > 1). The predicted penalty needs an `OR`
  that *blocks* an index, and on these rows the selective indexable `OR` is
  served by one BitmapOr over `orders_amount_idx`, so the penalty never appears.
- **Adding an expression index is a physical decision, not a spelling decision,
  and it is engine-dependent.** With the query text unchanged, PostgreSQL's
  selective expression `OR` drops from 0.845 ms to 0.265 ms (3.2×) once
  `orders_amount_expr_idx` exists, while the same arm on YugabyteDB does not
  improve (4.05 ms → 4.15 ms) and `pg_relation_size` reports 0 for the index.
- **The overlap control fires everywhere it must** (`x05` detected in all six
  shapes in all seven runs; 228 duplicated rows at the selective arms, 4,546 at
  the broad arms), so the `UNION ALL`-is-safe claim is policed on the measured
  data, not assumed.
- **Biggest doubt:** the client is a single sequential reader with no writers, so
  these are latency measurements of tiny pages, not throughput; the across-run
  spread on the selective cells (up to 59%) is larger than several of the effects
  being compared, and the YugabyteDB index-maintenance probe is unusable.

## What I was looking at

Twenty-eight matrix cells plus two expression-index package cells per run
(35 cells when the position control is included): `OR` / `UNION ALL` / `UNION`
for one disjunction at {selective, broad} × {indexable, expression, cast}, and
`OFFSET` / keyset at page depths 1, 10, 100, 1,000, 10,000, on PostgreSQL 17.11
and YugabyteDB 2025.2.6 single-node, scale `deep` (12,000 orders, seed 42), five
trials per cell. Seven whole runs: three PG, three YB, and one PG position
control (plus a YB position control on run 2). All 34/35 cells were `valid`; no
cell was invalid, unsupported, timed out or errored, and every cell was verified
against the independent Go oracle before timing. The client was **not** throttled
(0 CFS throttled periods at the 2-CPU / 2 GiB quota in every run), so the read
matrix is not client-bound. I ignored the teaching-fixture validation cells and
the customers/orders teaching analogue, which carry no performance claim.

## Conclusions

**Pagination (P11): confirmed, and the crossing is engine-specific.** The
keyset/OFFSET median ratio is the number to read. On PostgreSQL, four runs give
1.15/1.19/1.20/1.15 at depth 1 and 1.09/1.16/1.17/1.14 at depth 10 — OFFSET is
*not* worse, it is slightly ahead. At depth 1,000 the ratio is 0.83/0.71/0.73/0.75
and at 10,000 it is 0.23/0.20/0.23/0.23: keyset wins, and the win is an order of
magnitude at the deepest level. On YugabyteDB the ratio stays above 1 through
depth 1,000 (1.43/1.17/1.33) and only drops to 0.54/0.67/0.69 at 10,000. The
plan explains it: at depth 10,000 PostgreSQL's OFFSET statement runs the full
ordered scan-and-discard and reads 2.11 ms against the keyset page's 0.47 ms
(different plan structures, `ac4509…` vs `dd3164…`); the keyset plan reaches the
same rows through a BitmapOr plus a top-N heapsort. So the actionable rule is
*bounded*: keyset pays off past a depth that is roughly 1,000 on PostgreSQL and
10,000 on YugabyteDB for this data shape, and below it OFFSET is acceptable.
Page equality under a fixed snapshot held at every depth in every run (all five
depths equal, `REPEATABLE READ` accepted by both engines), so the two spellings
really return the same page, not merely similar times.

**Predicate alternatives (P10): falsified.** The prediction was that a
suitably-aligned `UNION ALL` could beat `OR` where `OR` blocks an index. It does
not. Across all 36 (shape × run) pairs on YugabyteDB, `OR` is faster than both
alternatives, by 1.74–1.84× at `selective/indexable` and 1.48–1.86× at
`selective/expression`; the signs are unanimous. On PostgreSQL the broad shapes
mostly favour `OR` (yes) and the selective shapes are unresolved: the selective
ratios sit between 0.88 and 1.38 with mixed signs across the four runs, and the
only selective pair with a stable sign, `selective/indexable OR vs UNION`
(1.09/1.18/1.09/0.96), is still marginal. The mechanism the prediction needed is
absent: PostgreSQL plans the selective indexable disjunction as a single
BitmapOr over `orders_amount_idx` with a Bitmap Heap Scan, so `OR` is not
penalised, while `UNION`/`UNION ALL` pay a second scan and an append/de-duplicate
step. `UNION` and `UNION ALL` are indistinguishable here because the arms are
provably disjoint (`OverlapCount = 0` for all six shapes); their only difference,
the de-duplication, costs nothing at these sizes.

**The expression index is the factor, not the syntax — and it only helps
PostgreSQL.** The matrix deliberately runs the same expression SQL with the index
absent (default arm) and present (`index=expression_index`). On PostgreSQL the
selective expression `OR` falls from 0.845 ms to 0.265 ms across the four runs
(3.2×, and the same pattern for `UNION`/`UNION ALL`); the broad expression shape
does not change (1.32 ms both, a sequence scan either way, 39% of rows). On
YugabyteDB the selective expression `OR` is 4.05 ms default against 4.15 ms with
the index — no gain — and the index reports size 0. A reader must not conclude
"write `(amount+0)` instead of `amount`"; the correct instruction is that the
expression spelling needs a matching expression index, and that this is worth
3× on PostgreSQL selective predicates but was not usable on this YugabyteDB
build. The storage cost of that index on PostgreSQL is 114,688 bytes; the
500-row insert batch showed no stable maintenance penalty (with/without
4.07/3.42, 3.77/4.21, 3.98/4.23, 4.38/3.13 ms — the sign flips, so the batch is
inside noise).

**Controls and the position term.** The required `x05` overlap control was
detected in all six shapes in every run: `UNION ALL` over an arm with itself
returns twice the arm-A matches (228 rows selective, 4,546 broad), so the
answer-safe claim for disjoint `UNION ALL` is enforced by a control that fires,
and the `x01`–`x03` controls that the measurement fixture cannot exercise are
recorded as `control_absent` rather than counted as passes. The repeated-design
position control put the first cell last: PostgreSQL read 0.319 ms first and
0.282 ms last (ratio 0.885), YugabyteDB 3.548/3.945 (ratio 1.112), same plan
both times. That is a ~11% session position term, comparable to the per-cell
spreads, and far smaller than the engine gaps and the depth effect.

## Where I think the measurement is weak

- **Latency, not throughput.** Every cell is one drained query returning five
  rows (the disjunction cells return a few thousand rows) measured as
  client round-trip time by one sequential reader. There is no concurrency, no
  demand generator and no throughput claim; "faster" here means lower latency
  for the same page/answer.
- **Spread dominates the selective cells.** The selective cells have sub-ms
  medians and across-run spreads up to 59% (`selective/indexable union_all`
  0.281–0.465 ms); several P10 differences are inside that spread and are
  reported as unresolved rather than as effects. The five-trial sample cannot
  support p95/p99 tails.
- **The YugabyteDB index probe is unusable.** `pg_relation_size` returns 0 and
  the insert "with index" is *faster* than "without" in all three runs
  (15.5/21.4, 22.9/31.2, 23.5/62.1 ms), which cannot be a real index effect and
  is most likely the order of the two batches and cold state. The YB storage and
  maintenance numbers should be treated as absent, not as a small cost.
- **One host, single node.** No independent hardware, no real network, no YB
  3-node RF=3 arm, no equal-total-budget arm and no verified tablet/leader
  placement. The YB 3-node comparison required by the owner's rules is an
  explicit coverage gap on this 2 GiB machine, not a completed comparison. The
  declared 3 GiB + 2 GiB budgets also still exceed the 2 GiB guest.
- **The selective threshold is one choice.** `amount <= 5 / amount >= 496`
  selects ~0.95% per arm; a different threshold could move where the index plan
  is chosen, and the sweep does not test the boundary between they do and do not.

## What I would measure next

1. **A selectivity sweep for the disjunction** (0.01%, 0.1%, 1%, 10%, 50%, 90%)
   crossing `OR`/`UNION ALL` with the expression-index arm, to find the exact
   selectivity above which PostgreSQL stops using `orders_amount_idx` and the
   `OR` penalty the prediction expects actually appears.
2. **The same pagination pair under concurrent writes**, so "same page under a
   fixed snapshot" is tested where the snapshot is doing work, with the writer
   held fixed and the page re-read under its cursor.
3. **A matched index-cost phase** that inserts into a dedicated table with and
   without each index *in a randomized order*, several batches, so the
   maintenance term is separable from cold state (and a YB-specific size query).
4. **Rows-per-page and row-width levels** at the keyset/OFFSET crossing, to test
   whether the PG ~1,000 / YB ~10,000 crossing moves.

## Disagreements with other analyses

None: no earlier analysis exists for this inputs digest. The per-run reports this
analysis reads are generated measurements; the verdicts above are mine.
