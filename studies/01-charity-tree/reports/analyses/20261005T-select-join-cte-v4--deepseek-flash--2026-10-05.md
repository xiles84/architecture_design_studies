---
analysis_id: 20261005T-select-join-cte-v4--deepseek-flash--2026-10-05
run_id: select-join-v4-topo-{single,equal,distributed,tablet,colocated}-{i,b}{1,2,3}-20261005; domain figures carried from select-join-v3b-domain-{pg,yb}-{i,b}{1,2,3}-20261005
environment: host-zenbook-ux5406sa
analyst: deepseek-flash
analyst_kind: ai
analyst_version: "DeepSeek Flash via opencode CLI (LOW executor session)"
analyzed_at: 2026-10-05
inputs_digest: "the 30 v4 topology runs, one client build; each per-run report carries its own join.json digest"
repo_commit: 8a1dd768
supersedes: 20261004T-select-join-cte-v3--deepseek-flash--2026-10-05
status: current
headline: "With data actually split across three nodes, splitting turns the 90% scan from 1.31x the single node into 0.71x and the materializing CTE from 1.25x into 0.87x, while a point read pays ~1.2x; colocation is cheapest at a point read but its materializing plan costs 1.81x the single node."
---

# Analysis (v4 successor) — placement that is measured, and the corrected reuse signs — deepseek-flash

This supersedes v3 (now in `reports/outdated/`; replacement index extended).
Review 3 accepted the v3 domain conclusions and asked for two things: a cluster
arm whose data is genuinely split and gated on placement, and a P4 verdict that
follows the per-run signs. Domain figures are carried forward from the v3b
domain report, which was not re-run. Topology figures are across-run medians
from `reports/select-join-v4-topology-20261005.md`.

## TL;DR

- **The v3 "non-colocated" cluster arm was not distributed.** Its five relations
  each had one tablet, all led on one node. It is now relabelled
  `noncolocated-single-tablet`, and a new `distributed` arm pre-splits the three
  tables and the two indexes into three tablets, with a gate that checks
  `num_tablets = 3` and three distinct tablet leaders before any cell is timed
  (before and after listings saved). The 30-run subset is one build.
- **Splitting is a real, measured effect.** Against the single-tablet cluster at
  the same budget, `distributed` is **0.85x at 2,000 rows and 0.71x at 180,000
  rows** but **1.19x at the 20-row point read**; against the single node it is
  0.71x at 180,000 rows. Splitting pays for parallel range scans and costs a
  point read.
- **Colocation vs distributed** (the required pair): colocated is faster at the
  point read (2.26 vs 3.18 ms) and at 2,000 rows (6.28 vs 8.02), but slower at
  180,000 rows (92.93 vs 77.05), and its forced-materializing plan costs **17.24x
  its default** at 2,000 rows against **8.31x** for distributed.
- **P4, corrected:** on YugabyteDB at the 1% consumer there is a **replicated
  sign change between one reference (0.81-0.94, 6/6 below 1) and two
  (1.43-1.52, 6/6 above 1)** — a small effect (6-19%) but a consistent one. On
  PostgreSQL at the 1% consumer there is **no crossing in 1-16** (0.17-0.81).
  PostgreSQL at the full consumer crosses at or below two references (n=2
  1.70-2.06), with n=1 mixed (0.87-1.05); YB at the full consumer has n=1 mixed
  (0.90-1.07) and n>=2 above 1.

## What was measured

- **Placement-gated topology subset, 30 runs, one build:** single 2 CPU/3 GiB,
  equal-total single 6 CPU/9 GiB, and the three cluster placements `distributed`,
  `noncolocated-single-tablet` (relabelled) and `colocated`, each isolated and
  blended, three whole runs, 30 cells, five trials, position control. 30/30 valid
  in every run; both controls fire.
- **Gate evidence:** `donation` and `donation_rank_idx` report `num_tablets = 3`
  only in the distributed arm; in the distributed runs every relation's leaders
  are `yb-n1:9100, yb-n2:9100, yb-n3:9100`. The single-tablet and colocated runs
  record one leader node (yb-n3 and yb-n2 respectively). The engine's effective
  `ysql_num_shards_per_tserver` is recorded as 2, and the arithmetic is not the
  cause of the earlier one-tablet outcome: the tables were created without a
  split, **marked unexplained beyond that** (with the flag recorded, a one-tablet
  outcome for a relation that was not explicitly split is expected behaviour, and
  the per-relation `num_tablets` is the authority).
- One feature check on the pinned engine confirmed `SPLIT INTO 3 TABLETS` and
  `SPLIT AT VALUES ((66667),(133334))` are accepted and produce three tablets.

## Topology and placement results (across-run medians, `join_where`)

| arm (2 CPU/3 GiB unless noted) | 20 rows | 2,000 rows | 180,000 rows | forced materialization / default at 2,000 |
|---|---|---|---|---|
| single node | 1.71 | 9.49 | 108.84 | 9.55x |
| single node, equal total (6 CPU/9 GiB) | 1.70 | 7.85 | 100.92 | 9.30x |
| cluster3, noncolocated-single-tablet | 2.61 | 10.94 | 142.83 | 11.98x |
| cluster3, distributed | 3.18 | 8.02 | 77.05 | 8.31x |
| cluster3, colocated | 2.26 | 6.28 | 92.93 | 17.24x |

Read against the **single node**: the single-tablet cluster is 1.53x at a point
read and 1.31x at 90% (its one leader carries every relation); the distributed
cluster is 1.86x at a point read but **0.71x at 90%** — splitting more than pays
for the broad scan; the colocated cluster is 1.32x at a point read and 0.85x at
90%, but its materializing plan is 1.81x the single node's. Read as the
**splitting effect** (single-tablet -> distributed, same budget, same data, same
index arm): 1.22x slower at the point read, 0.73x at 2,000 rows, 0.54x at
180,000 rows. The equal-total budget still buys nothing measurable.

The blended phase moves these comparisons by under ~10% and the writer is
balanced across all three endpoints (e.g. [168/168/168] reader operations); it
is not a capacity test.

## P4 reuse, per run

Carried from the v3b domain report and read per run (no re-run), per the review's
table:

- **YugabyteDB, 1% consumer:** n=1 is 0.81-0.94 in all six runs and n=2 is
  1.43-1.52 in all six. This is a **replicated sign change between one and two
  references**. Its size at one reference is 6-19%, which is small but consistent
  across every run.
- **YugabyteDB, 100% consumer:** n=1 is 0.90-1.07 across six runs — **mixed,
  unresolved**; n>=2 is above 1 in every run.
- **PostgreSQL, 1% consumer:** 0.17-0.81 at every n up to 16; NOT MATERIALIZED
  stays faster and there is **no crossing in 1-16**.
- **PostgreSQL, 100% consumer:** n=1 is 0.87-1.05 (**mixed, unresolved**); n=2 is
  1.70-2.06, so the crossing is at or below two references.

**Action:** the earlier blanket "the crossing exists only on PostgreSQL" was
wrong; the correct statement is the four cases above, and a YB crossing at two
references exists at the selective consumer.

## Unchanged from v3 (accepted by review 3)

The domain conclusions stand: with the matched range arm, forcing `MATERIALIZED`
on the joined-output CTE costs PostgreSQL ~97x and YugabyteDB ~46x at 20 rows;
the default spellings are oracle-identical and converge to ~1x at 90%; the `x01`
and `x02` negative controls are detected in every run. The physical-arm gate,
per-cell access evidence, compatibility classes (now including placement), and
the PG `storage_rows_scanned = n/a` and whole-plan scan-flag naming are in force.

## Where the measurement is weak

- Three whole runs per arm; the placement claims rest on engine metadata plus
  `yb-admin` listings read before the cells, not on independent hardware — the
  three "nodes" are containers on one host, so the numbers show the effect of
  *data placement*, not of a real network.
- The distributed arm's point-read penalty (1.22x) is within the spread of a
  three-run median; treat it as "no better, possibly slightly worse".
- The writer is a fixed 10/s maintenance load; this is not an optimistic-vs-
  pessimistic concurrency comparison.
- Rows-scanned and leader data are parsed from engine text; an unparsed value is
  reported as `n/a`/`not recorded`, not as zero.

## Coverage

- Measured pairs: single-tablet vs distributed (the splitting effect), colocated
  vs distributed (the required pair), single vs equal-total budget, isolated vs
  blended, and the reuse pair over five reference counts and two consumers.
- Non-applicable here: rollup/rolldown/embedding — this experiment holds the
  physical model fixed and varies query spelling and placement; the study's other
  arms measure those decisions.
- Outstanding: an independent placement mechanism that spreads leaders without
  an explicit split; more runs per arm; a real multi-host cluster.
