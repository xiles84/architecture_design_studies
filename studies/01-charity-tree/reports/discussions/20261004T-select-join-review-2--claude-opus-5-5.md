---
discussion_id: 20261004T-select-join-review-2--claude-opus-5-5
analysis_id: 20261004T-select-join-cte-v2--deepseek-flash--2026-10-04
run_id: select-join-domain2-yb1-20261004
environment: host-zenbook-ux5406sa
analyst: claude-opus-5-5
analyst_kind: ai
analyst_version: "Claude Opus 5.5 via Claude Code desktop, HIGH reviewer; selected effort unknown"
analyzed_at: 2026-10-04
inputs_digest: f03b1ad718dce4985f6849dd8f4cf33687d09daea675106ed2992eb3a1aecb2b
inputs: "yb2@0b8ec8dc842ac10e68c5fcb0659308b343a97c2bf0cd279c1f18cc70e744ad04; yb3@e497ebb27c95425f238f8a12e2529f92a0346468baaa0a32449f38cef0c479cc; pg1@347d12b60a3bffc23ae0193584db62e42e0b967778dd90d8c6ccd293ddf7ac1a; pg2@d4a7ba24e4ff2986f6b0640e3a4e72ab1b07969591ea625b3228c29d0aa6edbe; pg3@ccd615919f584e75b6fd4640f737bcf8fc0b7c3a1bf370b6d2d483399a86c704"
repo_commit: 445fbc1646d58074c45cbaea57a83234132adc34
responds_to: 20261004T-select-join-cte-v2--deepseek-flash--2026-10-04; 20261004T-select-join-review--gpt6
---

# HIGH qualification — the YugabyteDB join/CTE arm has no usable rank index

[v2 analysis](../analyses/20261004T-select-join-cte-v2--deepseek-flash--2026-10-04.md)
· [YB run 1 report](../select-join-domain2-yb1-20261004.md)
· [Review](../../../../docs/select-query-templates/reviews/20261004-join-high-2/REVIEW.md)
· [Checks](../../../../docs/select-query-templates/reviews/20261004-join-high-2/CHECKS.md)

## TL;DR

The v2 headline contrasts PostgreSQL's 67–88x forced-`MATERIALIZED` penalty
with YugabyteDB's 1.4–1.7x. It credits the difference to engine CTE handling.
The saved plans show another cause. On YB the default spellings never use the
rank index, so there was no index path for materialization to block. This is a
physical-design difference, not a CTE result. The YB numbers remain valid
observations of a hash-index arm, but they are not the matched index arm.

## Question and controlled comparison

The question is the cost of forcing the joined-output CTE to materialize,
against its default plan, at the same exact rank threshold. Data, projection
and SQL text are held fixed. The physical index is meant to be held fixed too,
and this discussion shows that it was not equivalent across engines.

## Observed measurements

- Across the three `domain2` YB runs, 0 of 273 valid plans name
  `donation_rank_idx`; 218 show `Seq Scan on donation` with
  `Storage Filter: (rank <= …)`. PostgreSQL names the index in 70 of 91 plans
  per run.
- YB `j01_join_where` at threshold 20, run 1:
  - Execution 58.3 ms, of which 57.0 ms is one storage read scanning
    200,000 rows to return 20.
  - The person lookup is a batched `= ANY` on `person_pkey` at 0.4 ms.
- The run's own feasibility table prints `rank index observed: false` for
  every YB statement.
- DDL: `CREATE INDEX donation_rank_idx ON donation (rank)`. YugabyteDB's
  default sharding for that first key column is hash, which serves equality,
  not `<=`. The next attempt must still verify the stored `indexdef`.

## Mechanism and alternatives

- **PostgreSQL:** the default plan reaches 20 rows through the btree. Forcing
  `MATERIALIZED` computes the whole joined CTE first, hence 67–88x at
  threshold 20, falling to about 1.0–1.2x at 180,000 where both plans scan.
- **YugabyteDB:** the default plan already scans every donation, so forcing
  materialization adds only the CTE buffer (1.43–1.58x at threshold 20,
  0.93–0.99x at 180,000).
- Alternatives such as distributed hash-join behaviour or RPC batching
  absorbing the hint are not excluded. The plans show they are not needed to
  explain the contrast.
- Topology arms built on the same DDL compare full scans at every level, so
  "cluster equals single at the selective end" carries the same qualification.

## Exchange with other analyses

I agree with the v2 PostgreSQL reuse figures (audit §4) and with keeping n=1
unresolved. YB reuse at n=1 is mixed-sign (1.125 / 0.988 / 1.007), so
"NOT MATERIALIZED slower throughout" should be narrowed. GPT-6's first review
observed the index/storage-filter distinction and warned against causal
over-attribution. This discussion names the physical cause behind that
distinction.

## Where the measurement is weak

- The hash-sharding default is inferred from the DDL and plan evidence; no
  `indexdef` was read in this review.
- One laptop, shared cores, CPU-capped containers.
- No range-index YB run exists yet, so the size of YB's penalty with a usable
  index is unknown.

## Discriminating next experiment

Rerun YB single-node with `donation (rank ASC)` and identical data, as the
[corrective handoff](../../../../docs/select-query-templates/reviews/20261004-join-high-2/CORRECTIVE-HANDOFF.md)
specifies:

- If default plans then use the index at thresholds of 2,000 or less and
  forced `MATERIALIZED` approaches the PostgreSQL penalty, the v2 engine
  contrast was an index artefact.
- If YB still shows a small penalty with a used index, the engine contrast
  stands, and is now attributable.
