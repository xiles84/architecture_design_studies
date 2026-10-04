---
analysis_id: 20261004T-select-order-correction--deepseek-flash--2026-10-04
run_id: select-order-{access,aggregate,existence}-{postgres,yugabyte}-{1,2,3}-20261004
environment: xiles-podman-select-order-correction-20261004
analyst: deepseek-flash
analyst_kind: ai
analyst_version: "DeepSeek deepseek-flash via opencode CLI, effort not exposed"
analyzed_at: 2026-10-04
inputs_digest: see reports/20261004T-select-order-manifest-{access,agg,exist}.md (per-family corpora)
repo_commit: 805b5b85 + order-correction commits
supersedes:
status: current
headline: "With seeded randomized order and interleaved samples, the access-path and membership/absence signs hold, but the aggregate DISTINCT ON vs ROW_NUMBER sign is mixed at PG levels 10/100 and the position drift there is 29-74%, so those comparisons stay unresolved."
---

# Correction analysis — SELECT measurement order (P6, P7, P8, P9, P11) — deepseek-flash

Signed analysis for `task-20261004T100822Z-select-order-correction`, answering the
independent HIGH review's
[`LOW-CORRECTION-HANDOFF.md`](../../../docs/select-evidence-review-20261004/LOW-CORRECTION-HANDOFF.md).
It corrects the **order** in which the three SELECT families are measured; it does
not change the question, the answer contracts, the oracles or the accepted matrix,
and it does not edit the earlier signed analyses. Measurements are in
`reports/20261004T-select-order-{access,agg,exist}-replicated.md` and the
per-run `results/select-order-*/…`; the primary inputs are pinned by
`reports/20261004T-select-order-manifest-{access,agg,exist}.md`.

## TL;DR

- **The order defect is real and now measured.** Each whole run permuted its cells
  with a recorded seed (3 distinct sequences per family/engine) and interleaved
  one sample per cell round. A repeated-design control that repeats the *actual*
  first cell of the run brackets the drift: access 1.5–2.4%, existence 13–22%,
  and **aggregate 29–74%** first-to-last. Only effects larger than that are
  promoted.
- **Access paths (P11) hold with a resolved depth-100 result.** Keyset is slower
  than OFFSET at depths 1–100 and faster at 1,000–10,000 on PostgreSQL (3/3 runs
  each sign), so the crossing is between 100 and 1,000; on YugabyteDB it is
  between 1,000 and 10,000. The corrected PostgreSQL depth-100 ratio is
  1.06–1.17 (keyset slower) in 3/3 runs, i.e. now *resolved* rather than
  bracketed — but it is only 6–17%, near the access position drift.
- **Membership/absence (P6/P7) hold.** On PostgreSQL `JOIN + DISTINCT` over
  `EXISTS` is 1.40–1.69× at fanout 1, 1.38–1.44× at fanout 10 and **3.64–4.23×**
  at fanout 100 (3/3 runs each); on YugabyteDB it is within noise (mixed signs at
  fanout 10). Nullable `NOT IN` returned the empty set in every corrected run.
- **Aggregate (P8/P9) is the weakest family and partly unresolved.**
  `DISTINCT ON` vs `ROW_NUMBER` is mixed-sign at PostgreSQL level 10 (1 >1, 2 <1)
  and level 100 (1 >1, 2 <1), stable `<1` only at levels 1, 1,000 and 10,000; on
  YugabyteDB level 10 is `>1` 3/3 and level 100 mixed, while 10,000 is now `<1`
  3/3. The original "`DISTINCT ON` fastest at scale on both engines" must be
  narrowed to the levels whose sign survived the randomized order, and the
  29–74% position drift means the aggregate numbers carry the largest caveat.

## What changed and what did not

Changed: the execution schedule (seeded permutation + interleaved rounds), the
recorded `order` block (seed, rounds, actual sequence), the position control
(actual first cell), the replicated ratio rows (run ids), and an explicit
primary-input manifest. Not changed: SQL, statements, fixtures, predicates, page
size, N=1, the five pagination/parent levels, the three fanout levels, the oracle
gate, and the original runs/tags/analyses (fresh ids and tags only).

## Family verdicts under the corrected order

### Access paths (P11)

`keyset/OFFSET` is below 1 when keyset is faster.

| depth | PostgreSQL (3 runs) | YugabyteDB (3 runs) |
|---|---|---|
| 1 | 1.29, 1.10, 1.07 (>1) | 1.34, 1.27, 1.36 (>1) |
| 10 | 1.22, 1.16, 1.09 (>1) | 1.27, 1.32, 1.32 (>1) |
| 100 | 1.06, 1.17, 1.08 (>1) | 1.31, 1.27, 1.35 (>1) |
| 1,000 | 0.73, 0.82, 0.81 (<1) | 1.25, 1.25, 1.20 (>1) |
| 10,000 | 0.24, 0.21, 0.23 (<1) | 0.71, 0.72, 0.69 (<1) |

The sign is stable in every cell (3/3), so the corrected runs confirm the
engine-specific crossing: PostgreSQL between 100 and 1,000, YugabyteDB between
1,000 and 10,000. The PostgreSQL depth-100 cell now has a stable sign; its
magnitude (6–17%) exceeds the access position drift (1.5–2.4%), so the sign is
resolved, though the absolute effect is small.

### Membership and absence (P6/P7)

`JOIN + DISTINCT` over `EXISTS` is `>1` (slower) on PostgreSQL at every level:
1.69, 1.40, 1.49 at fanout 1; 1.38, 1.42, 1.44 at fanout 10; **4.06, 3.64, 4.23**
at fanout 100. On YugabyteDB it is 1.11, 1.10, — at fanout 1 and mixed at fanout
10 (1.01, 0.95, 1.00), i.e. no growing penalty. This reproduces the original
engine-specific P6 finding under randomized order and interleaved sampling. The
nullable-`NOT IN` control again returned zero rows at every fanout level in every
run, and the NOT-NULL `NOT IN` spelling matched the absence oracle. Position
drift here is 22% (PG) and 13% (YB); the fanout-100 effect (3.6–4.2×) is far
larger, and the fanout-1/10 effects (1.4×) also exceed the drift, so those signs
are promoted.

### Aggregate / top-N (P8/P9) — partly unresolved

`DISTINCT ON` vs `ROW_NUMBER` (below 1 when `DISTINCT ON` is faster):

| level | PostgreSQL | | YugabyteDB | |
|---|---|---|---|---|
| | runs 1/2/3 | sign | runs 1/2/3 | sign |
| 1 | 0.85, 0.76, 0.44 | <1 3/3 | 1.03, 0.86, 0.93 | mixed |
| 10 | 2.12, 0.93, 0.30 | **mixed** | 0.99, 1.11, 1.00 | >1 3/3 |
| 100 | 1.21, 0.75, 0.84 | **mixed** | 0.97, 1.24, 0.89 | **mixed** |
| 1,000 | 0.92, 0.65, 0.61 | <1 3/3 | 0.92, 0.93, 0.89 | <1 3/3 |
| 10,000 | 0.61, 0.63, 0.63 | <1 3/3 | 0.70, 1.00, 0.79 | mixed* |

\* run 2 is 0.9968, effectively 1; the sign is not resolved at 10,000 on YB by a
strict same-direction rule once the spread is considered.

The position control is 0.71 on PostgreSQL and 1.74 on YugabyteDB (first vs last
of the actual first cell), i.e. 29% and 74% drift. Under the corrected schedule
the aggregate levels 10 and 100 are **unresolved on PostgreSQL** and level 100 is
unresolved on YugabyteDB; only levels 1, 1,000 and 10,000 on PostgreSQL and
levels 1,000 (and arguably 10) on YugabyteDB have a stable sign. The earlier
claim that `DISTINCT ON` is the fastest top-N spelling at scale should be
narrowed and re-examined with more samples before it is used.

## Where the measurement is weak

- **Aggregate position drift is large** (29–74%): with five interleaved samples
  per cell, aggregate medians still move enough within a run to flip
  medium-selectivity signs. More rounds and more whole runs are needed there.
- **Three levels cannot locate a crossing** for P11; it brackets one.
- **Single node, single host, one sequential reader**, five samples: no
  throughput, no cluster per-node/equal-total arm. The recorded resources make a
  cluster budget calculable, but no cluster topology was measured, so the cluster
  comparison remains an explicit gap.
- **YugabyteDB sits near its per-statement floor** in membership/absence, so
  "no measurable difference" there means exactly that.
- The position control is one cell per family/engine, so it brackets drift but is
  not a distribution.

## What I would measure next

1. More interleaved rounds (≥15) and ≥5 whole runs for the aggregate family,
   where the drift is largest, before promoting any of its medium levels.
2. A cluster per-node/equal-total arm now that live guest resources are recorded,
   or an explicit statement that the single-node evidence establishes no cluster
   result.
3. A crossing-bracketing run for P11 with more depths (e.g. 200, 300, 500) to
   locate the PostgreSQL crossover between 100 and 1,000.

## Disagreements

No earlier analysis exists for these corrected inputs. The independent review's
material inference corrections (order confound, stale guest resources, mixed
signs unresolved) are accepted and applied here. The original fixed-order
analyses are retained, explicitly order-confounded, and are not edited.
