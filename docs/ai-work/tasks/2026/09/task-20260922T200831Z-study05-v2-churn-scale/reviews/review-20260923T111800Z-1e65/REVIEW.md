# Review — Study 05 v2 churn and scale

**Task:** `task-20260922T200831Z-study05-v2-churn-scale`
**Reviewer:** HIGH session, model `deepseek-flash`, role `reviewer`
**Decision:** **approve for integration.** The regime the machine can produce was run and is
gate-clean; the regimes the harness cannot produce are named as gaps, and the hard-expiry zero is
reported rather than glossed.

## Decision Log reviewed

1. Run what the harness can produce, record the rest as gaps — accept; the acceptance authorises the
   gap branch for unproducible regimes.
2. 600 s and 1800 s chosen so the boundaries are 2.00 and 6.00 at a fixed 300 s TTL — accept.
3. `expired_hard = 0` reported as a finding — accept; this is the honest core of the result.
4. Write-through and write-aside of the strict owned model — accept; these are the write-path variants.
5. The 8k/80k/800k tiers exceed the harness maximum — accept as a gap.

## Independent checks

- Recomputed from result JSON: boundaries 2.00 / 2.00 / 6.00; reads 1,467,721 / 1,364,997 / 2,879,077;
  hard expiries 0 / 0 / 0; probabilistic 12,517 / 11,321 / 37,063; wrong reads and impossible values
  all 0. Matches the analysis.
- All three manifests: `failed_cells: []`; all cells `gate.passed = true` (14 checks).
- Three run tags exist on the recorded commit `7ae82b71…`.

## Acceptance criteria

| Criterion | Verdict |
|---|---|
| 600 s regime crosses the TTL and records expiries | met in duration (2.00; expiries recorded) |
| 1800 s regime | met in duration (6.00) |
| 5× burst for 120 s | **gap** — no write-rate multiplier in the harness |
| Refreshes recorded | **gap** — no refresh counter |
| Medium scale 8k/80k/800k donors | **gap** — harness maximum is 3,000 people |
| ≥3 trials, or recorded order with the limitation | met (measured phases); churn is single, stated |
| Wrong-read rate and throughput from the same cell | met — 0 wrong reads alongside throughput |

## Conditions

The book and registry may now say the 300 s TTL boundary was crossed under load with zero wrong reads
in the measured cells. They must **not** say the hard TTL was observed to expire an entry: that path
fired zero times. The burst, refresh counter and large-scale tiers stay gaps.

## Integration

Approve and integrate; tag `study-05/v2-churn-scale`.
