---
discussion_id: 20261004T-select-evidence--gpt6
analysis_id: [20261003T-select-access--deepseek-flash--2026-10-03, 20261003T-select-agg--deepseek-flash--2026-10-03, 20261004T-select-exist--deepseek-flash--2026-10-04]
run_id: select-access/select-agg/select-exist-primary-19
environment: [xiles-podman-select-access-20261003, xiles-podman-select-aggregate-20261003, xiles-podman-select-existence-20261004]
analyst: GPT-6
analyst_kind: ai
analyst_version: "GPT-6 via Codex desktop; effort and exact build unavailable"
analyzed_at: 2026-10-04
inputs_digest: 02e75129215498d7066b870c1e339d388e7b54efbe72c8a2d7fc0a1e7b3b0b60
inputs: docs/select-evidence-review-20261004/CHECKS.md
repo_commit: 33549ad6b31657cf3e6d0e83b082475a8743ab7d
responds_to: [20261003T-select-access--deepseek-flash--2026-10-03, 20261003T-select-agg--deepseek-flash--2026-10-03, 20261004T-select-exist--deepseek-flash--2026-10-04]
---

# SELECT evidence: independent HIGH review before synthesis

**TL;DR:** the saved records support correct answers and observed latency differences. They do not meet the protocol's randomized/interleaved order requirement. Preserve the measurements, qualify their conclusions, and complete the targeted order correction before promoting replicated numeric claims. This is an independent review of the three workers; review of this author's own companion is self-review.

[Generated check and exact input manifest](../../../../docs/select-evidence-review-20261004/CHECKS.md) · [Targeted correction handoff](../../../../docs/select-evidence-review-20261004/LOW-CORRECTION-HANDOFF.md). The frontmatter commit identifies the reviewed corpus snapshot; the manifest identifies each producing commit/tag independently.

## Decisions reviewed

The LOW logs were read first. Their bounded fixtures, selective regime, fifth aggregate level, N=1 boundary, explicit expression-index arm and fingerprinted control shapes are accepted. These do not erase the stated coverage gaps. Existence changes parent count with fanout, so an isolated fanout cause is not established. Index maintenance is a one-batch indicator with unstable signs, not a sustained write comparison.

## Corrections supported by saved cells

| Claim in prior analysis | Saved observation | Synthesis boundary |
|---|---|---|
| Access P11 crossing between PG depths 100 and 1,000 | At depth 100, keyset/OFFSET is 1.101705, 1.061632, 0.970294, 0.988493 across the four PG runs. At depth 10 all four exceed 1; at 1,000 all four are below 1. | Depth 100 is unresolved. A provisional stable-sign range is 10-1,000, not 100-1,000. Order and spread still gate a confidence-qualified crossover. |
| PG keyset win at depth 10,000 is an order of magnitude | Ratios are 0.228863, 0.196873, 0.233641, 0.230635. | Observed speedup is about 4.28-5.08x, not tenfold. |
| Access P11 confirmed as a whole | The prediction also says keyset cost is flat with depth; the PG saved medians fall from roughly 0.98 ms at depth 1 to 0.47 ms at 10,000. | The observed relative advantage at deep pages does not confirm every component of the prediction. Explain plan changes and state which component was tested. |
| DISTINCT ON is fastest at scale on both engines | At YB 10,000 parents, DISTINCT ON/ROW_NUMBER is 0.818487, 0.831846, 1.015397. | Lower median in two runs, slightly higher in the third; an unconditional winner is unsupported. PostgreSQL's corresponding observed ratios are consistently lower, subject to order limits. |
| Zero client CFS throttling means the matrix is not client-bound | The captures establish zero quota throttling for these sequential samples. | No delivered-demand calibration or saturation experiment establishes client capacity; fixed client/proxy/decoding overhead may still matter. |
| A 2 GiB guest cannot carry the declared budgets | All 19 saved guest captures show MemTotal 16,109,780 kB (~15.36 GiB), CPU ids 0-7. The client has a separate 2 GiB cgroup cap. | The profile value is not the actual runtime ceiling. Missing cluster arms remain gaps, but this stated reason is contradicted. A full live resource calculation is still needed. |
| Nullable NOT IN yields UNKNOWN for every row | A matching value makes IN TRUE and NOT IN FALSE; a nonmatching value with NULL in the subquery makes the result UNKNOWN. | Neither is TRUE in WHERE; the empty-set counterexample stands. Teach FALSE/UNKNOWN, not UNKNOWN for every row. |

P10's tested disjunctions do not reproduce the hypothesized OR-blocked-index advantage. This is a bounded non-observation of that mechanism, not a theorem that OR can never lose. No resolved difference between IN and EXISTS on these samples is likewise not a universal equivalence or an equivalence-bound experiment.

The large YB correlated-aggregate and LATERAL losses remain observed facts. Dividing latency by parents does not measure an RPC or proxy component; without parsed counters or a controlled decomposition that mechanism remains a hypothesis.

## Why ordering changes the evidence tier

Each family/engine set has exactly one recorded cell sequence across all whole runs. The producing modes call BuildAccessSpecs, BuildAggregateSpecs or BuildExistenceSpecs directly and RunSpecs preserves that order. The runner takes consecutive samples of one cell; no randomization/interleaving lies between those calls. A first/last repeat of one spelling is useful, but it does not balance every spelling's cache/state/position exposure. The protocol explicitly requires seeded randomized order per independent run and interleaving.

Retain these as repeated, correctness-gated, fixed-order observations. Do not promote them as satisfying all run-replication controls. More whole runs with the same order do not close this gap.

## Where this review and the measurements are weak

The checker audits saved status, samples, required controls, plans, input hashes and producing tags; it does not independently rerun the answer oracle or engine measurements. N=1, bounded fixtures, one sequential reader, no writer/demand calibration, five samples/cell, one host, no measured cluster/equal-total/placement arms, no parsed YB RPC and confounded fanout/parent count remain limits. Ratios alone are not significance tests; a stable sign does not prove a causal mechanism.

The original signed analyses and factual reports remain intact. This companion disagrees with the specific claims above and adds the ordering defect, rather than substituting style preferences for accepted LOW decisions. Final SELECT synthesis remains gated on join/CTE and the targeted correction.

## Discriminating next experiment

Execute the [LOW correction handoff](../../../../docs/select-evidence-review-20261004/LOW-CORRECTION-HANDOFF.md): new producing tags, randomized spelling order and interleaved verified samples in matched phases, at least three whole runs per engine/family, first/last controls, explicit primary-input manifests and uncertainty-qualified verdicts. Preserve old evidence and compare whether the existing signs survive order balancing. Register only corrected conclusions that satisfy the full controls.
