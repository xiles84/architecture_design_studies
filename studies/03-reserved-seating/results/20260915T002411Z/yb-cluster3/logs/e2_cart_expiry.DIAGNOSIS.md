# Diagnosis — yb-cluster3 / e2_cart_expiry (run 20260915T002411Z)

- diagnosed_by: claude-opus-5, setting `high` (LOW model/effort), Claude Code desktop
- diagnosed_at: 2026-09-15 16:10 UTC, while the matrix continued
- verdict: **same cause as yb-single E2** (`../../yb-single/logs/e2_cart_expiry.DIAGNOSIS.md`):
  saturated nodes stalled WAL appends and transaction-status RPCs; statements timed out; the
  lifecycle's sold-seat monitor query was one of them, and the harness treats that as fatal. Not
  re-run; covered by ER-03.

## What failed

Verify (58/58), explain, reads, writes, the race (every tier) and the 10-seat lifecycle tier
finished. On the 100-seat lifecycle tier:

```
ERROR: lifecycle monitor event 639: ERROR: Timed out waiting kResponseSent, state: kRequestSent (SQLSTATE XX000)
```

The race recorded early rejections, all transient: 4 in the 1 000-seat tier and 3 in the
10 000-seat tier. They are allowed on YugabyteDB by AM-01.5. Deferred confirmations were 100%.

## Evidence (node log tails captured by the runner, `cpu.stat` before and after)

| Node | `kResponseSent` timeouts | `Append to log took a long time` | `Failed to request transaction ... status` | CPU throttled during cell |
|---|---:|---:|---:|---:|
| yb-n1 | 51 | 200 | 140 | 56% |
| yb-n2 | 50 | 304 | 4 | 44% |
| yb-n3 | 11 | 304 | 0 | 32% |

The `lease` lines in the tails are raft elections for tablets created by the phase's reload
(`pre-election: Vote granted`, `Leader lease expiration was not set`), not expired leases.

After the failure all three containers were up, and the next cell (E1) started normally.
