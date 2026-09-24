# Progress — Study 05 v2 open-loop demand phase

**Task:** `task-20260922T200831Z-study05-v2-open-loop`
**Branch:** `study-05/v2-open-loop`
**Attempt:** `attempt-20260924T014745Z-5393c0` (claim `claim-e08dd51705b86b14`, epoch 1)
**Capability/role:** HIGH session executing the task / executor (model `deepseek-flash`)
**Required tag:** `study-05/v2-open-loop`

Measurement task: benchmark lock held by the runner, one cell at a time, nothing else measured.

## Decision Log

1. **The phase uses the platform driver, not a local loop.** `measure.RunOpenLoop` already separates
   offered, started, completed, rejected, errors, dropped, scheduling lag and a fixed saturation
   verdict; re-implementing it would have produced a second, unvalidated scheduler. Confidence: High.
2. **Rates chosen to bracket the same cell's closed-loop floor.** Canonical run offers 5,000/s (below
   the 21,175/s floor) and 30,000/s (above it). Two earlier rate pairs (4,000/16,000 and 6,000/12,000)
   were superseded by the rate choice, not by a defect; their reports are in `reports/outdated/` and
   their results remain. Confidence: High.
3. **Per-rate wrong-read accounting.** The study's wrong-read log is cumulative within a phase, so the
   second rate initially inherited the first rate's reads; the loop now resets per rate. Confidence:
   High.
4. **The server-saturation question is answered honestly, not forced.** At 30,000/s the driver marks
   the *client* saturated (143k dropped arrivals), which by the driver's own rule invalidates a claim
   about the server. The analysis therefore states the offered rate, the delivered lower bound and
   that the server's knee is not isolated. Confidence: High.
5. **Read-only overload is labelled read-only.** The cell is relaxed but has no writers, so 0 wrong
   reads is not a freshness result and the analysis says so. Confidence: High.
6. **The canonical run is the third.** Run 1 closed-loop floor 7,787/s, run 2 19,921/s, run 3
   21,175/s — the run-to-run spread is recorded as a limitation rather than averaged away.
   Confidence: Medium-High.

## What was produced

- `harness/{main,workload,result,report}.go` + `run-study.sh` — the `openloop` phase and its report
- `results/20260923T0220Z-openloop/`, `results/20260923T0240Z-openloop/`,
  `results/20260923T0255Z-openloop/` — the three runs (results immutable)
- `reports/20260923T0255Z-openloop.md` — canonical generated report (digest `50d5a426c89e96cf`)
- `reports/outdated/20260923T0220Z-openloop.md`, `reports/outdated/20260923T0240Z-openloop.md`
- `reports/analyses/20260923T0255Z-openloop--deepseek-flash--2026-09-23.md` — signed analysis
- `CONTEXT.md` — open-loop status and the remaining-gap update

## Residual risks / limitations

- One trial per rate; the closed-loop floor moved 2.7x across three runs of the same design.
- The client saturated in both regimes, so no server-side saturation point is measured.
- One design/topology/engine, read-only, shared laptop cores, `db-only` framing.
