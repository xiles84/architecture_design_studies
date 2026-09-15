# Diagnosis — yb-single / s4_check_then_hold_serializable (run 20260915T002411Z)

- diagnosed_by: claude-opus-5, setting `high` (LOW model/effort), Claude Code desktop
- diagnosed_at: 2026-09-15 03:30 UTC, while the matrix continued with other cells
- verdict: **cause established; the failure is the design's own load, not the environment.** Not
  re-run (§11 re-runs only environmental failures). Not an Escalation Required (§12 item 7 needs
  an unestablished cause).

## What failed

The cell ran verify (58/58), explain, reads, writes, the race (every tier) and the 10-seat lifecycle
tier, all without a violation. On the 100-seat lifecycle tier, the harness's sold-seat monitor query
failed and aborted the cell:

```
ERROR: lifecycle monitor event 640: ERROR: Timed out waiting kResponseSent, state: kProcessingRequest (SQLSTATE XX000)
```

The result JSON was written, with every phase that finished.

## Evidence

Source: the full PostgreSQL-layer logs, copied out of the running `yb-single` container with
`podman cp` (no process started inside the database container, so the next cell's measurement was
not disturbed). The last two log files are kept beside this file, gzipped.

- 57 080 `ERROR:  deadlock detected` lines during the cell (02:21:48 → 03:20:49 UTC). This is how
  YugabyteDB aborts conflicting SERIALIZABLE transactions (engine probe P7: 40P01); the harness
  retries them.
- The race shows the collapse: 0.1 seats/s on every tier, hold p50 22–62 s, every event timed out.
  The isolated hold benchmark ran at 1.0 ops/s with p50 8.2 s and 267 retries.
- 40 `Timed out waiting kResponseSent` statement RPC timeouts, the first at 03:16:39 on design
  `UPDATE`s. The sold-seat monitor's `SELECT COUNT(*) ... status = 'sold'` timed out at 03:19:11,
  which is the error that aborted the cell.
- tserver (captured `s4_check_then_hold_serializable.yb-single.log`): at 03:20:48, repeated
  `Write failed: Timed out (yb/docdb/wait_queue.cc:151): Waiter transaction timed out waiting in
  queue`.
- Database container `cpu.stat` at the end of the cell: throttled in 10 095 of 34 344 periods
  (29%), 3 142 s throttled.

## Reading (facts only)

Under S4's SERIALIZABLE read-then-write on a CPU-throttled node, conflicting transactions queue and
abort, and they saturate the tserver. Statement RPCs, including the harness's own monitor query,
then time out. Study 02 recorded the same pattern for its SERIALIZABLE design (LESSONS_LEARNED: "A
server-side failure needs the server's logs"). The harness treats a failing monitor query as fatal
for the cell. That is why the remaining lifecycle tier has no numbers; the design's statements
themselves were still being retried.

Expected on yb-cluster3 as well; that cell is diagnosed there if it fails.
