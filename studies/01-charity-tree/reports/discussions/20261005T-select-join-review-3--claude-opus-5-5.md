---
discussion_id: 20261005T-select-join-review-3--claude-opus-5-5
analysis_id: 20261004T-select-join-cte-v3--deepseek-flash--2026-10-05
run_id: select-join-v3b-domain-yb-i1-20261005
environment: host-zenbook-ux5406sa
analyst: claude-opus-5-5
analyst_kind: ai
analyst_version: "Claude Opus 5.5 via Claude Code desktop, HIGH reviewer; selected effort unknown"
analyzed_at: 2026-10-05
inputs_digest: "per-run join.json digests are printed in each v3b per-run report; all 37 verified equal (review 3 CHECKS.md)"
inputs: "select-join-v3b-domain-{pg,yb}-{i,b}{1,2,3}-20261005; select-join-v3b-topo-cluster-{i,b}{1,2,3}-20261005; select-join-v3b-coloc-{i,b}{1,2,3}-20261005"
repo_commit: 7243fee074bfb16fc5ad8b6afc1a2d7549080935
responds_to: 20261004T-select-join-cte-v3--deepseek-flash--2026-10-05
---

# HIGH qualification — reuse crossing and cluster placement in v3

[v3 analysis](../analyses/20261004T-select-join-cte-v3--deepseek-flash--2026-10-05.md)
· [Domain report](../select-join-v3b-domain-replicated-20261005.md)
· [Topology report](../select-join-v3b-topology-20261005.md)
· [Review 3](../../../../docs/select-query-templates/reviews/20261005-join-high-3/REVIEW.md)
· [Checks](../../../../docs/select-query-templates/reviews/20261005-join-high-3/CHECKS.md)

## TL;DR

I agree with v3's corrected materialization figures: 97.9x on PostgreSQL vs
45.6x on YugabyteDB at threshold 20, with a used range index. I disagree with
two of its readings:

- **The reuse crossing is not PostgreSQL-only.** YugabyteDB at the 1% consumer
  changes sign between one and two references in all six runs.
- **The cluster placement claims rest on an arm that was not distributed.**
  In every non-colocated run, every relation had one tablet and every leader
  sat on the same node.

## Question and controlled comparison

- **P4:** NOT MATERIALIZED vs MATERIALIZED for one CTE referenced n times,
  with the consumer selection and engine held fixed.
- **Placement:** the same data, RF3, budgets and endpoints, with colocated vs
  non-colocated storage on a 3-node cluster.

## Observed measurements

- **Reuse ladder** (n = 1, 2, 4, 8, 16; per-run medians; full table in
  CHECKS):
  - YB, 1% consumer: n=1 is 0.81–0.94 (6/6 below 1) and n=2 is 1.43–1.52
    (6/6 above 1).
  - YB, 100% consumer: n=1 is 0.90–1.07, mixed sign.
  - PG, 1% consumer: 0.17–0.81 at every n.
  - PG, 100% consumer: n=1 is 0.87–1.05, mixed; n=2 is 1.70–2.06.
- **Placement:** in 6 of 6 non-colocated cluster runs, `charity`, `person`,
  `donation`, `donation_rank_idx` and `donation_person_idx` each have one
  tablet, all led by one node. Colocated runs have one shared tablet led by
  one node. The tservers were started with `ysql_num_shards_per_tserver=2`.

## Mechanism and alternatives

- **Reuse.** On YB at a selective consumer, inlining one reference is slightly
  cheaper (6–19%) than building the CTE buffer. From two references, repeating
  the expensive token computation outweighs that saving. PostgreSQL at 1% never
  reaches that point within 16 references. This plausibly reflects relative
  executor costs; the plans would need loop-level comparison to establish it.
- **Placement.** Both cluster arms read from the same leader node. Their
  difference is shared tablet vs separate tablets, plus forwarding from the
  non-leader endpoints that the round-robin reader uses in both arms. It is
  not local vs cross-node data. Differences between the arms
  (e.g. 1.85 vs 3.22 ms at threshold 20; 16x vs 10x materialization penalty at
  2,000) are real observations of that narrower pair.

## Exchange with other analyses

This responds to `20261004T-select-join-cte-v3--deepseek-flash--2026-10-05`. Its
M1 correction, P1, P5 and P12 readings are consistent with the saved data. Its P4
verdict and placement mechanism are not, for the reasons above.

## Where the measurement is weak

- The YB n=1 effect at 1% is small, although its sign is consistent.
- Three runs per arm, on one host.
- The one-tablet outcome despite the requested shard flag is unexplained here.

## Discriminating next experiment

Pre-split every relation into three tablets with leaders on three nodes,
gate on that placement, and compare it with the colocated and single-tablet
arms under one build ([handoff 3](../../../../docs/select-query-templates/reviews/20261005-join-high-3/CORRECTIVE-HANDOFF.md)).
If the distributed arm is slower than single-tablet at the point read, data
spread has a measurable cost on this host. If it is not, the v3 cluster
penalty came from endpoint forwarding alone.
