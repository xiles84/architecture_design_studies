---
analysis_id: 20261003T-select-agg--deepseek-flash--2026-10-03
run_id: select-agg-20261003 (pg-1b/2/3, yb-1/2/3)
environment: xiles-podman-select-aggregate-20261003
analyst: deepseek-flash
analyst_kind: ai
analyst_version: "DeepSeek deepseek-flash via opencode CLI, effort not exposed"
analyzed_at: 2026-10-03
inputs_digest: 96fe2bd48af748a2add547115aa55d7f603ee9bbd709bd765b9fd71584044660
repo_commit: 33b769017da9648bef92cb1a804697e42c957e90
supersedes:
status: current
headline: "P8 reverses only on PostgreSQL (correlated→pre-agg→post-join as parents grow) and does not reverse on YugabyteDB, where correlated subqueries collapse to 8.5 s at 10k parents; P9's index probe wins only at one parent and DISTINCT ON is the fastest spelling at scale."
---

# Analysis — Study 01 SELECT aggregate & per-parent top-N (P8, P9) — deepseek-flash

Signed analysis of the aggregate/per-parent-top-N sweep for the additive SELECT
programme (task `task-20261002T131729Z-select-aggregate-topn`). Measurements are
in the per-run reports `results/select-agg-*/aggregate.md` and the generated
cross-run report [`20261003T-select-agg-replicated.md`](../20261003T-select-agg-replicated.md);
the bounded matrix is in
[`protocols/select-aggregate-topn-bounded-matrix.md`](../../protocols/select-aggregate-topn-bounded-matrix.md),
and predictions P8/P9 are in
[`protocols/select-query-templates.md`](../../protocols/select-query-templates.md).

## TL;DR

- **P8 reverses, but only on PostgreSQL.** At one parent, correlated beats
  post-join (0.70–0.98×) and pre-aggregation is 9–12× slower; at 10,000 parents
  the order flips: pre-aggregation becomes the fastest aggregate (0.81–0.91×
  post-join), correlated collapses to 2.7–3.0× post-join, and the correlated
  crossover sits between 10 and 100 parents in 3/3 runs. **The reversal is not a
  law of the query shape: on YugabyteDB pre-aggregation never wins at any level
  (1.13–3.26×) and the correlated crossover is between 1 and 10 parents.**
- **YugabyteDB's correlated per-parent aggregate is unusable at scale.** Its median
  grows 3.99 ms → 8.24 → 56.8 → 519 → **8,460 ms** across 1 → 10 → 100 → 1,000 →
  10,000 parents (489–507× post-join), while PostgreSQL's grows 0.24 → 20.45 ms.
  The YB cost is ~0.5–0.8 ms per parent, i.e. per-probe proxy/round-trip overhead.
- **P9's index probe almost never wins, and the winner is a spelling the protocol
  did not name as the favourite.** LATERAL never has a resolved advantage on
  PostgreSQL (0.88–1.05× ROW_NUMBER at 1–10 parents) and loses monotonically from
  100 parents (2.0–2.2× at 10,000); on YugabyteDB it wins only at one parent
  (0.57–0.63×) and then loses by 6.3–6.9× at 100 and **272–278×** at 10,000.
  `DISTINCT ON` is the fastest top-N spelling at every level from 100 parents up
  on both engines (PG 0.63–0.73× ROW_NUMBER; YB 0.83–1.02×), because it avoids
  both the global sort and the per-parent probe.
- **The controls fired in all six runs.** `x05_aggregate_double_count` returned
  10,000 rows that multiply the totals, and `x06_topn_rank_ties` returned 10,021
  rows (21 extra boundary-tie rows). The position control was 1.278 (PG) and
  0.951 (YB) on the same plan — below the level effects.
- **Biggest doubt:** the crossings are bracketed between adjacent levels, not
  located; all medians below ~0.3 ms are within a fixed-overhead cloud, and the
  YugabyteDB absolute numbers include single-node YSQL proxy cost, not a
  distributed 3-node penalty.

## What I was looking at

Thirty matrix cells per run plus one repeat, five trials each, on one
deterministic 10,000-parent fixture (2 orders per parent, seed 42) at
parent-prefix levels {1, 10, 100, 1000, 10000}: three parent-aggregate spellings
(post-join, pre-aggregate, correlated scalar) and three per-parent top-N
spellings (ROW_NUMBER, LATERAL, DISTINCT ON) compared at N=1. Six whole runs
(three PostgreSQL 17.11, three YugabyteDB 2025.2.6 single-node), all cells
oracle-verified, 0 invalid/unsupported/timed-out, and the client was not
throttled (0 CFS throttled periods). Every run passed the live six-control gate.
I am ignoring the aggregate-vs-rollup question: persistent rollups and their
freshness belong to Study 01/04/08, not here.

## Conclusions

**P8: the prediction holds on PostgreSQL and fails on YugabyteDB.** On
PostgreSQL the three spellings swap places twice as the parent set grows. Below
~10 parents the correlated scalar subquery is the cheapest (0.70–0.98× post-join
across the three runs) and pre-aggregation is 9–12× more expensive because it
aggregates all 20,000 child rows no matter how few parents are requested. The
correlated/post-join ratio crosses between 10 and 100 parents in 3/3 runs
(0.70/0.85/1.02/1.28/2.87; 0.82/0.79/1.09/1.22/2.96; 0.98/0.96/1.23/1.35/2.65).
At the other end, pre-aggregation's cost is nearly flat across levels (2.80 →
6.21 ms) while post-join's grows to 7.53 ms, so pre-aggregation overtakes
post-join between 1,000 and 10,000 parents in 3/3 runs (1.82→0.82, 1.69→0.91,
1.95→0.81). That is the brief's intended reading: pre-aggregation is the right
answer when the selected parent set is large, and the wrong answer when it is
tiny. **But it is engine-specific.** On YugabyteDB pre-aggregation never wins:
its ratio to post-join is 2.97–3.26× at small sets and still 1.13–1.30× at
10,000, so post-join is the best aggregate at every measured level. And the
correlated shape does not merely lose at scale, it collapses: 3.99, 8.24, 56.8,
519 and 8,460 ms across the five levels, i.e. ~0.5–0.8 ms for every extra
parent. A correlated scalar aggregate over 10,000 parents on this YugabyteDB
build is three orders of magnitude slower than PostgreSQL's. A reader must not
carry "pre-aggregate before a large read" across engines without measuring; and
must not carry "a correlated subquery is fine for a small result set" to
YugabyteDB beyond a handful of parents.

**P9: the probe ranking is the opposite of the protocol's framing at scale, and
DISTINCT ON wins.** The protocol predicted that a per-parent indexed probe wins
at small parent sets and loses at large, with a bounded crossing. On PostgreSQL
LATERAL is never clearly better than ROW_NUMBER — 0.88–1.05× at 1–10 parents,
inside the noise — and then loses with a stable sign: 1.27–1.34× at 100,
1.68–1.79× at 1,000, 1.99–2.17× at 10,000. On YugabyteDB the probe does win at
exactly one parent (0.57–0.63×) but the penalty at scale is enormous: 6.3–6.9×
at 100, 45.6–49.9× at 1,000, 272–278× at 10,000. So the crossing exists, but its
useful side is one parent wide on YB and absent on PG. What the protocol did not
test is the spelling that actually wins: `DISTINCT ON` is fastest from 100
parents up on PostgreSQL (0.63–0.73× ROW_NUMBER) and is at least equal to
ROW_NUMBER on YugabyteDB (0.83–1.02×), because on this `(customer_id,
amount)`-indexed data it removes the global sort without paying one round trip
per parent. The practical rule is therefore *DISTINCT ON for N=1*, with
ROW_NUMBER as the portable fallback for N>1, and LATERAL only for a very small,
explicitly bounded parent set.

**Controls and the position term.** The wrong-grain control is not a formality:
`x05_aggregate_double_count` is accepted as an aggregate answer by the engine and
returns 10,000 rows whose totals are multiplied by each parent's child count, so
without the oracle gate it would look like a legitimate, fast aggregate. The tie
control is equally live: `x06_topn_rank_ties` returns 10,021 rows because RANK
keeps the boundary tie that ROW_NUMBER drops, which is exactly the "RANK/
DENSE_RANK keep-all-ties answers are separate contracts" case. Both fired in
6/6 runs. The repeated-design position control read 1.278 on PostgreSQL and
0.951 on YugabyteDB with the same plan, so the session position term is smaller
than every crossing effect quoted above.

## Where I think the measurement is weak

- **The crossings are brackets, not points.** I can say the PG correlated
  crossover is between 10 and 100 parents and the pre-aggregation crossover
  between 1,000 and 10,000; I did not measure the levels between them. The
  one-parent level is sub-millisecond on PostgreSQL and dominated by fixed
  overhead, which is why the small-parent half of P9 is `no resolved
  difference` rather than a win.
- **YugabyteDB absolute times include single-node proxy cost.** The per-parent
  correlated/lateral growth is a real property of this build on this host, but a
  3-node RF=3 arm would change the constant; no placement or equal-total-budget
  claim is made. The YB storage/RPC counters are not parsed.
- **Two commits produced the 30 shared cells.** `pg-1b` ran at `08a3c553` and the
  other five at `33b76901`; the only change between them is the opt-in
  position-control cell. The pair ratios agree run to run, but the provenance is
  recorded rather than identical.
- **Fanout is fixed at two children per parent and N is fixed at one.** The
  join-multiplier contract (the reason x05 exists) is not exercised by a real
  second one-to-many relation here; no skew axis; `DISTINCT ON` was not compared
  at N>1 (it cannot express it).
- Five samples per cell cannot support p95/p99, and one host with no concurrent
  writers supports latency only, not throughput.

## What I would measure next

1. **Bracket the crossings** with levels 30, 300 and 3,000 parents, and add
   1, 3, 10 to the fanout axis, to locate the PG correlated and pre-aggregation
   crossovers and test whether the YB pre-aggregation ratio ever falls below 1.
2. **A real join-multiplier arm:** add a second one-to-many child collection and
   the correct pre-aggregate-then-join spelling beside the wrong-grain control,
   so the brief's counterexample is a measured pair rather than a control alone.
3. **Per-parent pushed-down pre-aggregation** (`LATERAL` pre-agg restricted to
   the selected parents) to test whether PostgreSQL's constant pre-aggregate can
   be avoided and the 1,000→10,000 crossing moves.
4. **N sweep for the top-N family** (N = 1, 3, 10) for LATERAL and ROW_NUMBER,
   since DISTINCT ON cannot cover N>1 and the portable spelling is what a reader
   will actually use.

## Disagreements with other analyses

None: no earlier analysis exists for this inputs digest. The per-run reports are
generated measurements; the verdicts above are mine.
