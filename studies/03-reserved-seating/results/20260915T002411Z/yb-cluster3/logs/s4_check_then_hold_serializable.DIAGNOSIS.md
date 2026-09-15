# Diagnosis — yb-cluster3 / s4_check_then_hold_serializable (run 20260915T002411Z)

- diagnosed_by: claude-opus-5, setting `high` (LOW model/effort), Claude Code desktop
- diagnosed_at: 2026-09-15 14:15 UTC, while the matrix continued
- verdict: **cause established: the design's SERIALIZABLE conflict storm, as on yb-single.** Not
  environmental, so not re-run; the cluster was healthy afterwards. The fatal-instrumentation part
  is covered by ER-03. This cell failed in a load, not in the monitor.

## What failed

Verify (58/58), explain, reads, writes and the race ran. The race collapsed, and the reload before
the lifecycle then failed:

```
ERROR: load before lifecycle#1: analyze venue_seat: ERROR: Timed out waiting kResponseSent, state: kRequestSent (SQLSTATE XX000)
```

A load is harness setup, and its failure ends the cell. The lifecycle phase has no numbers.

## Evidence

From the cell log and the node log tails the runner captured before teardown:

- Isolated hold benchmark: 1.0 ops/s, p50 10.4 s, 404 retries.
- Race, by tier:

  | Tier | Seats/s | Errors | Deferred confirmations | Sold |
  |---|---:|---:|---:|---:|
  | 10 seats | 0.1 | — | — | 43% |
  | 100 seats | 0.0 | 29 | — | 4% |
  | 1 000 seats | 2.8 | 7 585 | 6 of 22 | 0% |
  | 10 000 seats | 28.4 | 22 696 | 23 of 128 | 9% |

  Every event timed out. The first error in the three largest tiers was
  `Timed out waiting kResponseSent`. Most transactions were aborted, so the ledger and audit found
  no violation.
- Node logs, per node:

  | Node | Deadlock-aborted transaction heartbeats | `kResponseSent` timeouts |
  |---|---:|---:|
  | yb-n1 | 64 | 51 |
  | yb-n2 | 61 | 51 |
  | yb-n3 | 62 | 0 |

  The heartbeat lines read `Send heartbeat failed: ... aborted due to a deadlock`.
- CPU throttling during the cell: yb-n1 15%, yb-n3 10%, yb-n2 7% of periods.
- Afterwards the cluster was healthy: all three containers up; the next cell (L3) loaded in 30 s and
  passed its gate 58/58.

## Reading (facts only)

At SERIALIZABLE on YugabyteDB, S4's read-then-write aborts almost every concurrent transaction
through the deadlock detector. The race's 32 buyers per event, retrying, keep the conflict going.
Statement RPCs time out, and the ANALYZE issued by the next load is caught in the same backlog.
It is the same mechanism as the yb-single S4 cell and study 02's SERIALIZABLE design.
