# Diagnosis — yb-single / e2_cart_expiry (run 20260915T002411Z)

- diagnosed_by: claude-opus-5, setting `high` (LOW model/effort), Claude Code desktop
- diagnosed_at: 2026-09-15 09:40 UTC, while the matrix continued with other cells
- verdict: **cause established: sustained CPU saturation of the database node made statement RPCs
  time out, and the harness treats one failed monitor query as fatal for the cell.** Not
  environmental in §11's sense, so not re-run. Escalated for a decision on re-running (ER-03).

## What failed

Verify (58/58), explain, reads, writes, the race (every tier) and the 10-seat lifecycle tier finished
without a violation. On the 100-seat lifecycle tier the sold-seat monitor query failed:

```
ERROR: lifecycle monitor event 641: ERROR: Timed out waiting kResponseSent, state: kProcessingRequest (SQLSTATE XX000)
```

`harness/lifecycle.go` (around line 423) returns that error, which ends the cell. The failure has
the same message and the same point in the phase as the S4 cell on this topology
(`s4_check_then_hold_serializable.DIAGNOSIS.md`, event 640), but E2 runs at READ COMMITTED.

## Evidence

The full PostgreSQL-layer log was copied out of the running container with `podman cp`, and the
file covering the failure is kept beside this one, gzipped. The tserver tail was captured by the
runner (`e2_cart_expiry.yb-single.log`).

- 24 `Timed out waiting kResponseSent` errors, 09:21:49 → 09:22:21 UTC, on several backends; the
  cell's monitor query is one of them.
- tserver, same minutes:
  - 200 × `Time spent Append to log took a long time` (real 58–120 ms, user and sys 0 — waiting,
    not working);
  - 64 × `Failed to request transaction ... status`;
  - an `OpenTable` RPC callback warning.
- The database container was throttled in **86%** of its CPU scheduling periods during this cell.
  Per cell, from the runner's before/after `cpu.stat` snapshots:

  | Cell | Throttled |
  |---|---:|
  | S4 | 30% |
  | L2 | 63% |
  | S0 | 67% |
  | S1r | 70% |
  | S3 | 72% |
  | E1 | 72% |
  | K0 | 73% |
  | K1 | 74% |
  | L1 | 74% |
  | S2 | 74% |
  | S1 | 75% |
  | E2 | 86% |
  | L3 | 96% |

  L3 passed at 96%, so throttling alone does not decide which cell fails.
- No deadlock storm in this cell: 323 `deadlock detected` lines in the whole log file, most from
  the preceding S0 cell.

## Reading (facts only)

The 2-CPU YugabyteDB container is saturated for most of every `small` cell on this topology. When
write-ahead-log appends and transaction-status RPCs stall long enough, statements hit the
YSQL RPC timeout. The design's own statements are retried or recorded as errors by the buyers.
The lifecycle's sold-seat monitor is not: its first failure aborts the cell, and the remaining
lifecycle tier has no numbers. Whether such cells should be re-run with a timeout-tolerant monitor is
not a mapped decision (ER-03).
