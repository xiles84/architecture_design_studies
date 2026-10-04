# Matched blended read/write phase (review R2 / handoff C5)

LOW `deepseek-flash` / opencode CLI, 2026-10-04. Two phases of the same bounded
subset on PostgreSQL single-node, domain arm, 1 trial each — a validation pair,
not the final matrix.

| run | tag | phase |
|---|---|---|
| `select-join-blend-pg-isolated-20261004` | `run/study-01/select-join-blend-pg-isolated-20261004` | isolated, writer disabled |
| `select-join-blend-pg-blended-20261004` | `run/study-01/select-join-blend-pg-blended-20261004` | blended, one writer |

Both: 30/30 valid, 0 invalid/timeout.

## Writer contract

- The counter `maintenance_updates` is added to the child table with
  `DEFAULT 0`. It is **not** in the catalogue projection, not in the rank or
  index set and not in the MD5 reuse token, so a concurrent update cannot change
  a measured answer; the independent oracles ran in both phases and both are
  30/30 valid.
- One writer updates a fixed set of the lowest 1,000 child ids, rotating
  deterministically, at an offered **10 committed updates/s**. It is a fixed
  maintenance load, not an optimistic-versus-pessimistic comparison.
- The isolated phase loads exactly the same column and data with the writer
  disabled, so the pair differs only in the writer.

## Blended result

| observation | value |
|---|---|
| offered / completed / failed | 66 / 66 / 0 |
| wall time | 6.599 s |
| offered rate | 10.002 /s |
| completed rate | 10.002 /s |
| per-update latency samples | 66 |
| keys touched | 66 (rotation over the fixed set) |
| Go ledger | 66 keys reconciled, **0 mismatch** |
| reader operations overlapping a committed write | 23 of 90 |

The writer sustained the offered demand exactly and the ledger proves no update
was lost or duplicated. 23 of 90 reader operations ran while a write committed,
so the phase is a genuine concurrent measurement; the metric is recorded so a
short cell with no overlap cannot be reported as one.

## Isolated control

The isolated run records `phase=isolated`, writer offered/completed zero, and the
same 30/30 valid cells. Because `phase` and the writer target are provenance
settings, the compatibility class separates the two phases: they can never be
pooled in a replicated report.

## Still open

- The full phase pair at 5 trials on the domain matrix and on the bounded
  topology subset (single, equal-total single, cluster3), and the
  colocated/non-colocated arms.
- A writer on the cluster arm is bound to the first endpoint (one SQL endpoint);
  reads stay round-robin. Recorded, not hidden.
