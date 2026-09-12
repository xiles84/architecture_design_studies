# Partially invalid — do not use the `delete_person` numbers in this run

**Invalidated:** 2026-09-12 (UTC). **Replaced by:** `results/20260913-writes-isolated/`.

Same harness bug as `20260912-trials`: erasure consumes a pool of only 5 000 donors, the
pool drained during warmup, and every measured trial recorded tens of millions of instant
"no such row" failures and 0 ops/s.

`update_person` in this run is valid (0 errors, 3 trials, spread ≤ 11%).
