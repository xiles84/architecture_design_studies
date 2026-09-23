# Result — Study 05 v2 equal-total, engines and placement

**Task:** `task-20260922T200831Z-study05-v2-resources-engines`
**Capability/role:** HIGH executor (model `deepseek-flash`)
**Required tag:** `study-05/v2-resources-engines`
**Measured:** eight cells; PostgreSQL 17.11 / YugabyteDB 2025.2.6 + Redis; scale small; commit `e2d3eb42…` (placement re-run `7ca60245…`).

## Delivered

- **Equal-total vs db-only** (same design and topology): warm 10,837 vs 8,027; cacheable-99 7,876 vs
  8,130; app-99 7,945 vs 7,873. `resource_framing` recorded in every result; the arms are labelled, not pooled.
- **Strict/relaxed per engine:** yb-single RF=1 strict 514 / relaxed 596 app-99 ops/s; yb-cluster3
  RF=3 strict 942 / relaxed 846 — with **all three endpoints spread** (`connection_nodes: 3`).
- **Colocated/non-colocated pair:** y1 882/852/558 vs y2 673/557/387 (warm/cacheable/app), ~1.3–1.5x.
- **Harness fix:** `pgxdb.OpenSpread` plus `splitList` in Study 05; `connection_nodes` in results.
- All eight cells: gate passed, `failed_cells: []`, strict violations 0, impossible cache values 0,
  ledger mismatches 0.
- **Signed analysis:** `reports/analyses/20260923T1215Z-resources-engines--deepseek-flash--2026-09-23.md`.

## Acceptance criteria check

| Criterion | Result |
|---|---|
| Equal-total arm separately labelled, never pooled | met |
| Strict/relaxed pair on YugabyteDB 1-node RF=1 and 3-node RF=3, reported per engine | met (endpoints spread over 3 nodes) |
| Colocated vs non-colocated pair with the tablet/leader distribution recorded | **placement distribution NOT recorded** — only the server list; the speed difference is measured, the physical placement is not evidenced |
| Correctness gate passes; controls fired; one matrix at a time | met |

## Notes

The harness defect was real and blocking: without it the cluster arm could not connect, and a
first-endpoint-only shortcut would have been a false cluster measurement. The placement evidence
gap is recorded rather than implied.
