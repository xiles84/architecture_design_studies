# Third HIGH review — join/CTE round 2: changes requested (narrow)

Reviewer: Claude Opus 5.5 (`claude-opus-5-5`) / Claude Code desktop / HIGH;
selected effort unknown. Independent of the LOW author (`deepseek-flash` /
opencode CLI). Reviewed submission `99233466`, result `a2f3fa9f`, attempt
`attempt-20261004T215518Z-975148`. I read the Decision Log first. No database
was started; [checks](CHECKS.md) read saved results through [audit.go](audit.go).

**Verdict: request changes, limited to the cluster placement arm and the
interpretation.** Everything else in round 2 is accepted and need not be
repeated.

## Accepted

- **Decision Log.** Retroactive entries 1–5 and round-2 decisions 1–6 are all
  acceptable. Decision 6 (gate on the named 20/200/2,000 thresholds, then
  re-run every arm under one build) is the right correction.
- **Runs.** All 37 `v3b` runs come from one client build (`f3f1e10c…`). All
  cells are valid, every cell has 5 samples, both controls fire in every run,
  and every per-run report digest equals `sha256(join.json)`.
- **M1 fixed.** `donation_rank_idx` is `lsm (rank ASC)`. YB domain runs have
  66–67 index-scan cells, and `j01` at threshold 20 scans 110 storage rows
  instead of 200,000. The corrected materialization figures are traceable to
  the generated report:
  - PostgreSQL 48.454 / 0.4948 = **97.9x**;
  - YugabyteDB 86.60 / 1.900 = **45.6x** (isolated, threshold 20).
- **M3/M4 fixed.** Classes separate budget, phase, placement and index arm.
  There is no stale limit text. Throttling, endpoint, writer and retry evidence
  is rendered. The topology report re-renders **byte-identically** twice with
  the producing image (`f973dc84…`).
- **M5 fixed.** The domain matrix has isolated and blended phases on both
  engines. The writer offers and completes ~10/s with 0 failures, rotates
  across cluster endpoints, and overlap is recorded.
- **M2 partly fixed.** The colocated database was created and measured, and
  placement evidence is saved per run. v1/v2 and the pooled report moved to
  `reports/outdated/` with a replacement index.
- **Process.** Events use the real clock and the claim was valid at
  submission. Queue audit: 0 errors.

## P1 — the "non-colocated" cluster arm is not distributed

Saved `yb_table_properties` and `yb-admin list_tablets` show the same layout
in **6 of 6** non-colocated runs:

- `charity`, `person`, `donation`, `donation_rank_idx` and
  `donation_person_idx` have **one tablet each**;
- **all five leaders are on the same node** (yb-n3 in five runs, yb-n2 in one).

The colocated runs show the same leader concentration with one shared tablet.
The tservers were started with `ysql_num_shards_per_tserver=2`, yet one tablet
was created. Both arms therefore serve every read from one leader node.

The measured pair is "one shared tablet vs five single-tablet Raft groups on
the same leader". It is not "data on one node vs data spread across nodes".
The required comparison (AGENTS.md, required comparisons §3) is not yet met:
the non-colocated scenario never placed data on more than one node. The v3
claims built on it are not supported by the placement evidence:

- "colocation keeps the point read local";
- "one tablet per node concentrates the materialising work";
- "a cluster is 1.4–1.6x slower than single node" as a cluster-distribution
  effect.

Those runs remain valid observations of a labelled **default single-tablet**
non-colocated arm.

## P2 — the v3 P4 verdict is inverted

Per-run `NOT MATERIALIZED / MATERIALIZED` ratios (audit §3 and CHECKS):

- **YugabyteDB, 1% consumer:** n=1 is 0.81–0.94 (**6/6 below 1**) and n=2 is
  1.43–1.52 (**6/6 above 1**). That is a replicated sign change between one
  and two references.
- **YugabyteDB, 100% consumer:** n=1 is 0.90–1.07 (mixed, unresolved); n≥2 is
  above 1.
- **PostgreSQL, 1% consumer:** 0.17–0.81 at every n up to 16; NOT MATERIALIZED
  stays faster and there is no crossing in 1–16.
- **PostgreSQL, 100% consumer:** n=1 is 0.87–1.05 (mixed); n=2 is 1.70–2.06.

v3 states the crossing "exists only on PostgreSQL and only at the unselective
consumer" and calls YB n=1 at 1% mixed-sign. Both statements contradict the
per-run signs. The YB n=1 effect is small (6–19%), so the successor should
state its size alongside its consistent sign.

## Minor

- The per-cell `plan_index_scan`/`plan_seq_scan` flags are whole-plan
  booleans. Seq scans of the small parent tables set them, which is why
  PostgreSQL shows 24–25 / 61. Do not read them as the rank access path.
- PostgreSQL `storage_rows_scanned` prints 0; render it as n/a.
- Record the tservers' effective sharding flag and explain the one-tablet
  outcome, or say that it is unexplained.

## Route

The [corrective handoff 3](CORRECTIVE-HANDOFF.md) asks for:

- one new build with a pre-split, leader-spread non-colocated arm, gated on
  placement;
- the 30-run topology subset re-run under that build;
- a v4 successor analysis.

The domain matrix is not re-run. The completion tag stays absent. The checkpoint
tag is `study-01/select-join-high-review-3-20261005`.

NEXT MODEL: LOW
