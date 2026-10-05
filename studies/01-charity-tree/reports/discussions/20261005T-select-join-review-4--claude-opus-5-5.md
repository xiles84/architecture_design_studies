---
discussion_id: 20261005T-select-join-review-4--claude-opus-5-5
analysis_id: 20261005T-select-join-cte-v4--deepseek-flash--2026-10-05
run_id: select-join-v4-topo-distributed-i1-20261005
environment: host-zenbook-ux5406sa
analyst: claude-opus-5-5
analyst_kind: ai
analyst_version: "Claude Opus 5.5 via Claude Code desktop, HIGH reviewer; selected effort xhigh"
analyzed_at: 2026-10-05
inputs_digest: "acc5d12126ae90644bbca88ef4323ff6dfc063f7400c1e70752aacc448ef3c18 (sha256 of docs/select-query-templates/reviews/20261005-join-high-4/V4-INPUTS.txt: one `run sha256(join.json)` line per v4 run)"
inputs: "select-join-v4-topo-{single,equal,tablet,distributed,colocated}-{i,b}{1,2,3}-20261005; P4 figures from select-join-v3b-domain-{pg,yb}-{i,b}{1,2,3}-20261005"
repo_commit: ba821e7b686ac3adb1cda56dd45f15c8ae699ef8
responds_to: 20261005T-select-join-cte-v4--deepseek-flash--2026-10-05
---

# HIGH qualification — what the v4 placement runs establish

[v4 analysis](../analyses/20261005T-select-join-cte-v4--deepseek-flash--2026-10-05.md)
· [Topology report](../select-join-v4-topology-20261005.md)
· [Review 4](../../../../docs/select-query-templates/reviews/20261005-join-high-4/REVIEW.md)
· [Checks](../../../../docs/select-query-templates/reviews/20261005-join-high-4/CHECKS.md)
· [Review 3 discussion](20261005T-select-join-review-3--claude-opus-5-5.md)

## TL;DR

I agree with v4's measurements, its body figures for the splitting effect and
the colocated/distributed pair, and its corrected P4 ladder. I disagree with
six statements in its TL;DR, headline and body. Read v4 with these
corrections:

- **Splitting** (single-tablet → distributed, same cluster and budget):
  **0.73x at 2,000 rows and 0.54x at 180,000 rows**, with the point read
  **1.22x slower** (isolated). v4's TL;DR gives 0.85x/0.71x/1.19x for this
  pair; the first two are the single-node ratios.
- **Colocation vs distribution:** colocated is faster at 20 and 2,000 rows and
  **1.21x–1.54x slower at 180,000 rows**. Its forced-MAT penalty is larger in
  every run. Its forced-MAT plan itself takes **1.14x** the single node's time,
  not 1.81x; 1.81x is a ratio of two penalties.
- **The blended phase changes some comparisons by far more than 10%**,
  including one direction: colocated vs single node at 180,000 rows is 0.85x
  isolated and 1.16x blended.

## Question and controlled comparison

On a 3-node RF3 YugabyteDB cluster on one host, what does data placement
change for a join/filter and a CTE, with logical data, index arm, budgets and
endpoints held fixed?

- **Splitting effect:** single-tablet (engine default, one tablet per
  relation, all led by one node) vs `distributed` (every relation pre-split into
  three tablets, leaders gated onto three nodes).
- **Required pair:** `colocated` (one shared tablet) vs `distributed`.
- **Controls:** a single node at the same per-node budget, and a single node
  at the cluster's total budget.

The arms differ only in placement DDL. Tablet count and leader spread change
together in the splitting pair. They cannot be separated without a
split-but-single-node arm, which was not run.

## Observed measurements

Per-run medians and separation are in
[CHECKS](../../../../docs/select-query-templates/reviews/20261005-join-high-4/CHECKS.md).
A pair is cited below only where every run of one arm lies beyond every run of
the other.

| pair, `join_where` | isolated (20 / 2,000 / 180,000 rows) | blended |
|---|---|---|
| distributed / single-tablet | **1.22 / 0.73 / 0.54** | 1.06 (overlap) / **0.71 / 0.71** |
| colocated / distributed | **0.71 / 0.78 / 1.21** | **0.75 / 0.82 / 1.54** |
| distributed / single node | **1.86 / 0.85 / 0.71** | **1.66 / 0.93 / 0.75** |
| distributed / equal-total single node | **1.88** / 1.02 (overlap) / **0.76** | **1.43 / 0.80 / 0.73** |
| single-tablet / single node | 1.53 / 1.15 / 1.31, all overlap | 1.57 (overlap) / **1.32 / 1.06** |
| colocated / single node | 1.33 (overlap) / **0.66** / 0.85 (overlap) | **1.24 / 0.77** / 1.16 (overlap) |

Forced-MAT penalty at 2,000 rows (`cte_joined_outer_mat` /
`cte_joined_outer`, median of per-run ratios): single node 9.55x,
single-tablet 11.98x, distributed 8.50x, colocated 16.80x (isolated).
Colocated exceeds distributed in every run of both phases. The forced-MAT
plan's own time relative to the single node is 1.40x (single-tablet), 0.73x
(distributed) and 1.14x (colocated).

Placement is measured, not inferred:

- the distributed arm has three tablets per relation and leaders on yb-n1,
  yb-n2 and yb-n3 before and after every run;
- the single-tablet and colocated arms have every relation led on one node per
  run, yb-n2 or yb-n3 depending on the run.

## Mechanism and alternatives

- **Scan gain from splitting.** With three tablets on three leaders, a
  180,000-row filter is served by three tservers' scans in parallel. The
  equal-total single node has the same aggregate CPU but one tablet, and it is
  slower than the distributed arm at 180,000 rows in every run. Budget alone is
  therefore not the cause. The data do not separate "three tablets" from
  "three nodes", because both change together.
- **Point-read cost of splitting.** A 20-row join touching hash-spread parent
  and child rows needs more RPC fan-out when the relations are spread. That
  matches a 1.22x isolated cost, but the plans' RPC counters were not
  aggregated, so the mechanism is a hypothesis.
- **Colocation.** One shared tablet keeps all relations on one leader: the
  point and medium reads win and the broad scan loses. The large forced-MAT
  penalty is observed. Why the shared tablet amplifies materialization is not
  established by these runs.
- **Blended phase.** A 10/s writer competes with the readers for the same
  single leader in the colocated arm. That is consistent with its 180,000-row
  read growing from 92.9 to 128.1 ms, but per-tablet contention was not
  measured.

## Exchange with other analyses

This responds to `20261005T-select-join-cte-v4--deepseek-flash--2026-10-05`,
item by item:

1. **TL;DR baseline:** single-node ratios are labelled as single-tablet ones,
   and 1.19x matches no pair.
2. **Headline "1.25x into 0.87x" and body "1.81x the single node's":** these
   are ratios of MAT/default penalties (11.98/9.55, 8.31/9.55, 17.24/9.55),
   not plan costs.
3. **"The blended phase moves these comparisons by under ~10%":** false (see
   TL;DR).
4. **"One leader node (yb-n3 and yb-n2 respectively)":** the leader varies per
   run in both arms.
5. **"Single-tablet cluster 1.53x at a point read and 1.31x at 90%":**
   unresolved in the isolated phase.
6. **Shard flag "expected behaviour":** the configured
   `ysql_num_shards_per_tserver=2` on three tservers predicts six tablets for an
   unsplit hash table, and one was observed. The recorded value was not read
   back from a running tserver. The cause is unexplained.

Its P4 ladder agrees with review 3's per-run table, and its accepted domain
conclusions (97.9x/45.6x forced-MAT at threshold 20) are unchanged.

## Where the measurement is weak

- Three whole runs per arm and phase on one laptop; the "nodes" are containers
  sharing eight cores, so this is data placement without a network.
- The 1 CPU client is throttled in 32–46% of its CFS periods in every run.
  Point reads of 2–3 ms are the most exposed; medians of five samples damp
  but do not remove this.
- The single 2 CPU node is throttled in 19–23% of its periods. The
  single-node baseline is partly CPU-bound, which inflates every "/single
  node" ratio's denominator.
- Tablet count and leader spread are confounded in the splitting pair.
- The one-tablet default under a two-shards-per-tserver flag is unexplained.

## Discriminating next experiment

Run the `distributed` schema (`SPLIT INTO 3 TABLETS`) on the equal-total
single node, 6 CPU/9 GiB. If it matches the cluster's 0.71–0.76x at 180,000
rows, the gain is tablet parallelism and node spread adds nothing on this
host. If it stays near the single-tablet single node, spreading leaders
across nodes is the cause.
