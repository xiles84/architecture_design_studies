# Result — Study 04 v2 cardinality and cadence

**Task:** `task-20260922T195241Z-study04-v2-cardinality-cadence`
**Capability/role:** HIGH executor (model `deepseek-flash`)
**Required tag:** `study-04/v2-cardinality-cadence`
**Measured:** nine PostgreSQL cells at commit `f1d9b3c0…`; environment `host-zenbook-ux5406sa`.

## Delivered

- **Cardinality tiers 1/10/30/60/120/500** for n0_rows_unindexed, n1_rows_indexed,
  n3_rollup_trigger, n4_rollup_app. `r04_list_installations` collapses 62x for the row designs
  (11.8k → 188 ops/s) and is flat for the rollups (13.9k–14.7k); at 500 entries the rollup is ~74x.
- **Cadence regimes** quiet/steady/burst at 120 entries: offered rate delivered in every regime
  (1.15/s of 1/s, 60.1/s of 60/s, 66.3/s of 60/s), no drops, burst queue depth 22 of 64.
- **Trigger versus application maintenance**: `w05_replace_all` 10–11 ops/s (trigger) vs 413–448
  ops/s (application rollup) across the cadences — the v1 gap reproduced at a second cardinality.
- **Signed analysis:**
  `studies/04-configuration-portal/reports/analyses/20260923T0030Z-cardinality-cadence--deepseek-flash--2026-09-23.md`.
- Nine generated reports, nine run tags, all cells gate-passed.

## Acceptance criteria check

| Criterion | Result |
|---|---|
| Cardinality tiers 1…500 for the four designs | met |
| Unindexed scan cost grows with tier, else stop and report | grew (62x); phase continued |
| Cadence regimes at the 120-entry tier | met (delivery and lag) |
| Staleness and rollup lag for trigger vs application | **not met** — the harness records scheduling lag, not derived-state age; recorded as a gap |
| ≥3 trials per cell | met (all read cells) |
| Generated report + signed analysis | met |

## Notes

The one unmet clause is stated rather than approximated. Closing it needs a harness field that
stamps the derived value and its read; it belongs with the remaining Study 04 v2 work.
