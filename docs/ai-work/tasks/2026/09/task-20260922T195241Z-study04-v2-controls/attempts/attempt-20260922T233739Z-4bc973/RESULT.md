# Result — Study 04 v2 controls

**Task:** `task-20260922T195241Z-study04-v2-controls`
**Capability/role:** HIGH executor (model `deepseek-flash`)
**Required tag:** `study-04/v2-controls`
**Measured:** eight valid PostgreSQL cells; commit `f1d9b3c0…`, environment `host-zenbook-ux5406sa`.

## Delivered

- **Control A (identical-SQL reproducibility):** `n1_rows_indexed` in three fresh-load phases;
  read throughput reproduces within 0.9%–6.7% (most reads ≤2%).
- **Control B (retries-disabled):** at 16 writers the optimistic design reached 731 ops/s with
  `-retries 1`, faster than its default 552; the lock's advantage is 1.19x against the
  retries-disabled variant and 1.57x against default.
- **Control C (writer sweep):** c2/c1 = 0.89 at 1 writer, 1.26 at 4, 1.58 at 16, 1.60 at 64, with
  `lost_updates = 0` and a passed gate everywhere.
- **Signed analysis:** `studies/04-configuration-portal/reports/analyses/20260923T00116Z-controls--deepseek-flash--2026-09-23.md`.
- Eight generated reports with inputs digests; eight run tags (A1 auto-tagged; the seven later runs
  tagged manually on the same commit because the runner's dirty check counted generated files).
- Five failed first-attempt contention runs kept and excluded, with the failure cause stated.

## Acceptance criteria check

| Criterion | Result |
|---|---|
| Control A: identical-SQL floor measured | **partly** — cross-phase reproducibility measured; within-run position effect recorded as the open half |
| Control B: `-retries 1` control run; 1.76x restated only with it | met — the claim is now 1.19x–1.57x, below the old upper bound |
| Control C: 1/4/16/64 writer sweep with throughput and abort/retry rate | met (single contention measurement per point) |
| Correctness gate passes; controls fired; one matrix at a time | met — all eight cells passed; no failed cell used |
| Generated report + signed analysis | met |

## Notes

All runs used the committed Study 04 code; no harness change was needed. The task leaves one named
gap (within-run identical-SQL position effect) and one provenance-tooling defect (runner dirty check
counts generated files), both recorded in `LESSONS_LEARNED.md`.
