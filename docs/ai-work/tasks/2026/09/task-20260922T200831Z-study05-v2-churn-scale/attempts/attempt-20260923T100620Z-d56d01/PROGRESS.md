# Progress — Study 05 v2 churn and scale

**Task:** `task-20260922T200831Z-study05-v2-churn-scale`
**Attempt:** `attempt-20260923T100620Z-d56d01` (claim `claim-e15c19be0c0e29ab`, epoch 1)
**Capability/role:** HIGH executor (model `deepseek-flash`)
**Required tag:** `study-05/v2-churn-scale`

Measured: three churn cells (600 s write-through, 600 s write-aside, 1800 s write-through) on
PostgreSQL + Redis. One benchmark lock at a time.

## Decision Log

1. **Run the regimes the machine and harness can produce; record the rest as gaps.** The acceptance
   allows partial coverage when a regime cannot be produced; the burst multiplier and the large
   scales do not exist in the harness. Confidence: High.
2. **The 300 s TTL is never shortened**; the run simply lasts long enough. 600 s and 1800 s were
   chosen so the boundaries are 2.00 and 6.00. Confidence: High.
3. **The hard-expiry zero is reported as a finding, not smoothed over.** `expired_hard = 0` while
   probabilistic expiries fire means the boundary was crossed in time but the hard-expiry path did
   not remove an entry. Confidence: High.
4. **Designs: write-through and write-aside of the strict owned model** (`owned-opt-redis-through-strict-coord`,
   `owned-opt-redis-aside-strict-coord`), the cache write-path variants the amendment names. Confidence: High.
5. **The 8k/80k/800k donor tiers exceed the harness's largest scale** (3,000 people / ~120k
   donations, `scaleShape`); recorded as a gap rather than run at a different question. Confidence: High.

## Residual risks

- Hard expiry is unexercised; a cell that stops refreshing a subset or disables early expiry is needed.
- Single run per regime; no medium scale; PostgreSQL only.
