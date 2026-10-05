# Progress — Study 01 SELECT corrective round 3 (LOW `deepseek-flash`)

**Task:** `task-20261002T131729Z-select-join-cte`
**Attempt:** `attempt-20261005T111913Z-13271e`; claim `claim-901d6661ceb83a37`
(real clock, expires 14:30Z)
**Worktree:** `.worktrees/select-join-cte`; branch `study-01/select-join-cte`
**Handoff executed:** `docs/select-query-templates/reviews/20261005-join-high-3/CORRECTIVE-HANDOFF.md`

## Decision Log

1. **Podman machine was stopped.** `queue guard` and every container command
   failed with "podman is installed but not reachable"; `podman machine list`
   showed `LAST UP: Never`. Decision: start the existing machine
   (`podman machine start`) rather than work around it; it is required by the
   non-negotiables and is local and reversible. Confidence: High.
2. **Split DDL.** The `distributed` arm uses `SPLIT INTO 3 TABLETS` on the three
   tables and the person index, and `SPLIT AT VALUES ((66667),(133334))` on the
   range rank index, exactly as the handoff specifies. Feature-checked before
   any matrix: the pinned engine accepts it and all three relations report
   `num_tablets = 3` with donation's leaders on yb-n1, yb-n3 and yb-n2.
   Confidence: High.
3. **Leader gate by file handshake.** Tablet leaders are only visible through
   `yb-admin`, which the client image does not have. Decision: the client writes
   `resources/placement-prepared` after creating and splitting the schema, then
   waits (bounded) for the runner to write `resources/placement-gate.ok`; the
   runner polls `yb-admin list_tablets` for up to five minutes and only releases
   the client when three distinct leaders are seen. A failure aborts before any
   cell. Confidence: High. Alternative considered: run the client twice (once to
   prepare, once to measure) — rejected because the second prepare re-creates
   the tablets and invalidates the check.
4. **The runner backgrounds the client** so it can perform that handshake. The
   gate file logic is limited to `cluster3` + `distributed`; other arms are not
   gated on spread, as the handoff says.
5. **Backfilled `leader-summary.txt` for the 12 existing cluster runs** that were
   taken before the runner captured it, by parsing their saved raw
   `resources/placement.txt`. The raw listing is the primary evidence and is
   unchanged; the derived file only re-expresses it. Confidence: High.
6. **`leader_summary` is loaded at render time, not stored in `join.json`,** so
   `join.json` bytes (and the run digests the review checks) are unaffected by
   the reporting additions. Confidence: High.

## Deviations and defects found

- The first 30-run driver silently ran only the `single` arm because a
  sub-command consumed the driver's heredoc on stdin; the six runs were valid
  and are kept, and the remaining 24 were run by an array-driven script.
- The first leader parse read the `Range` column (which contains spaces) as the
  leader column; fixed to match the `node:port` token. The first smoke run is
  kept as evidence of that defect.
- `SHOW ysql_num_shards_per_tserver` is empty (it is a tserver flag, not a GUC);
  the runner now reads it from the node's inspect data and the client records it.

## Status

Complete against the round-3 acceptance list; submitted for review.

7. **Single-node leader evidence and report reproducibility.** The single-node
   arms have no `placement.txt`; their leader summary is written from the run's
   own recorded endpoint (`yb-single`), so the leader column is evidence from
   the run, not an inferred label. A render-time Go fallback for that case was
   written, then **reverted**: it landed after the 30 runs, so the frozen
   producing image would not have reproduced the report. The committed v4
   topology report is verified byte-identical to a re-render using the recorded
   producing image `b6ea3c007cf4` on a copy of the 30 runs.
