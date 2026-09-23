# Result — Study 05 v2 churn and scale

**Task:** `task-20260922T200831Z-study05-v2-churn-scale`
**Capability/role:** HIGH executor (model `deepseek-flash`)
**Required tag:** `study-05/v2-churn-scale`
**Measured:** three cells, PostgreSQL 17.11 + Redis, scale small, commit `7ae82b71…`.

## Delivered

- **600 s churn, write-through:** crossed **2.00** real 300 s boundaries; 1,467,721 reads /
  197,196 writes; 12,517 probabilistic expiries; **0 wrong reads**, 0 impossible values.
- **600 s churn, write-aside:** 2.00 boundaries; 1,364,997 reads / 182,282 writes; 11,321
  probabilistic expiries; 0 wrong reads.
- **1800 s churn, write-through:** **6.00** boundaries; 2,879,077 reads / 380,824 writes; 37,063
  probabilistic expiries; 0 wrong reads.
- All three cells passed the correctness gate (14 checks); no cell failed; each has a run tag.
- **Signed analysis:** `reports/analyses/20260923T1100Z-churn--deepseek-flash--2026-09-23.md`.

## Acceptance criteria check

| Criterion | Result |
|---|---|
| Churn cells cross the 300 s TTL (600 s, 1800 s, 5x burst 120 s) | 600 s and 1800 s **met** (2.00 and 6.00 boundaries); the 5x burst is a **gap** (no write-rate knob) |
| Record refreshes and expiries actually fired | expiries recorded (probabilistic; hard 0 — reported); refreshes **not recorded** by the harness |
| Medium scale at 8k/80k/800k donors with working-set:RAM and eviction | **gap** — largest harness scale is 3,000 people / ~120k donations |
| ≥3 trials per cell, or recorded randomized order with the limitation stated | met for the measured phases; churn is one sustained measurement, stated |
| Wrong-read rate and throughput from the same cell | met — wrong reads 0 alongside the throughput |
| A regime that crossed no TTL is a gap | no regime needed this; all crossed boundaries |

## Notes

The task's central claim is now measurable: the TTL boundary is crossed under load and correctness
holds. The hard-expiry path remains unexercised, and that is stated rather than implied.
