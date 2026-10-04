---
analysis_id: 20261004T-select-exist--deepseek-flash--2026-10-04
run_id: select-exist-20261004 (pg-1b/2/3, yb-1/2/3)
environment: xiles-podman-select-existence-20261004
analyst: deepseek-flash
analyst_kind: ai
analyst_version: "DeepSeek deepseek-flash via opencode CLI, effort not exposed"
analyzed_at: 2026-10-04
inputs_digest: 21feb94c45582af40bb9d0a694e405df7cd9350f62e5923a923b4dad79ee6f42
repo_commit: 805b5b85828387ee6a74912e815dbbfe56105efa
supersedes:
status: current
headline: "P6's fanout-growing JOIN+DISTINCT penalty is real on PostgreSQL (1.4x at fanout 1 to 3.6-4.4x at fanout 100) and absent on YugabyteDB; nullable NOT IN returned the empty set at every level (x02 fired 18/18), and on PostgreSQL NOT EXISTS beats the LEFT-JOIN/NOT-IN forms by 2-3x at fanout 100."
---

# Analysis — Study 01 SELECT membership & absence (P6, P7) — deepseek-flash

Signed analysis of the membership/absence sweep for the additive SELECT
programme (task `task-20261002T131729Z-select-existence`). Measurements are in
the per-run reports `results/select-exist-*/existence.md` and the generated
cross-run report [`20261004T-select-exist-replicated.md`](../20261004T-select-exist-replicated.md);
the bounded matrix is in
[`protocols/select-existence-bounded-matrix.md`](../../protocols/select-existence-bounded-matrix.md),
and predictions P6/P7 are in
[`protocols/select-query-templates.md`](../../protocols/select-query-templates.md).

## TL;DR

- **P6 is confirmed on PostgreSQL and refuted on YugabyteDB.** On PostgreSQL the
  `JOIN + DISTINCT` membership spelling is 1.44–1.47× `EXISTS` at child fanout 1,
  1.30–1.57× at fanout 10, and **3.61–4.40×** at fanout 100, in 3/3 runs at every
  level. On YugabyteDB the same comparison is 0.92–1.06× at every level and run:
  there is no measurable de-duplication penalty to grow with fanout.
- **`IN` and `EXISTS` are interchangeable on both engines.** The `in`-vs-`exists`
  ratio is 0.86–1.13 on PostgreSQL and 0.99–1.04 on YugabyteDB with mixed signs;
  the NOT-NULL subquery contract (`in_notnull`) is also within noise of both. Use
  whichever reads better; the fanout penalty is not theirs.
- **P7's semantic trap is confirmed decisively.** The nullable-`NOT IN` control
  returned **0 rows** at every fanout level in all six runs (x02 fired 18/18),
  while the NOT-NULL `NOT IN` spelling matched the absence oracle. A single NULL
  child key makes `NOT IN` return nothing; `NOT EXISTS` and `LEFT JOIN ... IS
  NULL` are the safe anti-joins.
- **Absence has a second, performance finding.** On PostgreSQL at fanout 100,
  `NOT EXISTS` ran in 0.22–0.34 ms against 0.63–0.72 ms for `LEFT JOIN ... IS
  NULL` and 0.68–0.72 ms for `NOT IN` — a 2–3× anti-join advantage in 3/3 runs
  that is small or absent at fanout 1. On YugabyteDB the three absence spellings
  are within ~10% of each other at every level.
- **Biggest doubt:** fanout is a three-level fixture axis with 10,000 child rows
  held constant, so the parent count changes with it (10,000 / 1,000 / 100); the
  "penalty grows with fanout" and "the parent set shrinks" effects are not
  separated by this bounded run.

## What I was looking at

Seven cells per fanout level, three levels, five trials each, on three
deterministic fixtures that hold 10,000 child rows constant with 1, 10 and 100
children per parent (so 10,000, 1,000 and 100 parents). Membership spellings:
`EXISTS`, `IN` over the nullable child key, `IN` over a NOT-NULL subquery,
`JOIN + DISTINCT`. Absence spellings: `NOT EXISTS`, `LEFT JOIN ... IS NULL`,
`NOT IN` over a NOT-NULL subquery. Six whole runs (three PostgreSQL 17.11, three
YugabyteDB 2025.2.6 single-node), all cells oracle-verified, 0
invalid/unsupported/timed-out, client unthrottled. Each level loaded its own
fixture; the two controls ran at every level. The predicate (`amount > 100`) is
identical across a pair.

## Conclusions

**P6: the de-duplication penalty is a PostgreSQL fact, not a property of
`JOIN + DISTINCT`.** On PostgreSQL the plain join emits one row per qualifying
child and the `DISTINCT` must collapse them back to parent identities; the cost
of doing so rises steeply as the qualifying fanout grows. The measured
`join_distinct`/`exists` ratio is stable across the three runs (1.459/1.444/1.471
at fanout 1, 1.350/1.297/1.566 at fanout 10, **4.397/3.989/3.612** at fanout
100). At fanout 100 the join-and-dedup path costs 1.06 ms against 0.24–0.30 ms
for the semi-join spellings. The protocol predicted exactly this shape; the
practical rule is that `EXISTS` (or `IN`) is the membership spelling, and a
`DISTINCT` join is a fanout-dependent tax. **But the rule does not transfer.**
On YugabyteDB the same ratio is 1.02–1.06 at fanout 1, 0.92–0.96 at fanout 10
and 0.98–1.06 at fanout 100 — within noise, with no trend — because every
membership spelling there sits at the engine's per-statement floor (3.5–9.4 ms)
and the de-duplication work is not the dominant term. A reader who measured only
YugabyteDB would correctly report "no difference"; a reader who measured only
PostgreSQL would correctly report a 4× penalty; both are true of their engine.

**`IN` and `EXISTS` did not differ.** Across all runs and levels the `in`-vs-
`exists` ratio stays between 0.86 and 1.13 with mixed signs, and the
`in_notnull`-vs-`exists` ratio (the same question under an explicit NOT-NULL
subquery) is likewise within noise. The falsifier "`EXISTS`/`IN` differ in answer
from the membership set" never fired: every cell passed the set oracle. So the
choice between `EXISTS` and `IN` on these engines is readability, not
performance — and the real decision is semi-join versus join-and-dedup.

**P7: the NULL trap is real, and `NOT EXISTS` is also the cheaper anti-join on
PostgreSQL.** The nullable-`NOT IN` control returned the empty set at every
level in every run: x02 fired 18/18 with 0 rows, against a non-empty absence set,
which is the protocol's "`NOT IN` over a NULL yields UNKNOWN for every row". The
NOT-NULL `NOT IN` spelling (`e13`, filtering `customer_id IS NOT NULL`) matched
the absence oracle and is a valid spelling — but only because the contract makes
the column non-null. On PostgreSQL there is also a performance reason to prefer
`NOT EXISTS`: at fanout 100 it ran 0.22–0.34 ms against 0.63–0.72 ms for
`LEFT JOIN ... IS NULL` and 0.68–0.72 ms for `NOT IN`, a 2–3× gap in 3/3 runs;
the gap shrinks at fanout 1 (the three spellings are within ~30%). On
YugabyteDB the three absence spellings are 3.5–8.3 ms and within ~10% of each
other at every level, so only the semantic reason remains there.

**Controls and the position term.** Both controls fired at all three fanout
levels in all six runs: x02 with 0 rows (the trap) and x07 with 7,240–7,241 rows
(the plain join's per-child multiplicity, against a membership set of at most
one row per parent). The repeated-design position control read 0.871 on PostgreSQL and 0.995 on
YugabyteDB with the same plan, below every fanout effect quoted above.

## Where I think the measurement is weak

- **Fanout and parent count move together.** Holding 10,000 child rows constant
  means fanout 100 has 100 parents and fanout 1 has 10,000, so the growing
  `JOIN + DISTINCT` penalty is confounded with the shrinking parent set. A
  follow-up that fixes the parent count and varies only the children per parent
  would separate them.
- **Three fanout levels.** The protocol's {1, 10, 100} bounds the direction but
  does not locate a crossing (no ≥5-level claim is made); the fanout-1 level is
  sub-millisecond on PostgreSQL and partly fixed-overhead.
- **YugabyteDB comparisons are close to its per-statement floor.** At 3.5–9 ms
  with 2–15% spreads, a 5–10% real difference would be invisible; "no measurable
  difference" on YB means exactly that, not "equal".
- **No selectivity, hit-fraction or skew axis**, no 3-node/equal-total/
  placement arm, no index redesign arm, five samples only, latency not
  throughput. The `ExistenceShape` fixture (one hot multi-child parent, one
  childless parent) is a control-shape construct and is recorded in the
  fingerprints.

## What I would measure next

1. **Fix the parent count and vary only fanout** (e.g. 1,000 parents with 1, 3,
   10, 30, 100 children each) to separate the de-duplication penalty from the
   shrinking parent set and to bound a crossing.
2. **A qualifying-proportion sweep** (1%, 10%, 50%, 90% of children matching) at
   fixed fanout, since the `DISTINCT` cost depends on qualifying rows, not total
   rows.
3. **The same matrix on a second child collection** with a non-nullable key, to
   show the NOT-NULL `NOT IN` contract is genuine rather than an artefact of the
   nullable study schema.
4. **A hot-parent skew axis**, which this bounded run records as a gap.

## Disagreements with other analyses

None: no earlier analysis exists for this inputs digest. The per-run reports are
generated measurements; the verdicts above are mine.
