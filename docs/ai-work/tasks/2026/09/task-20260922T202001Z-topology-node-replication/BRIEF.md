# Execution brief — Run the Study 07 node-count, replication and equal-total arms

Implements part of `studies/07-topology/HANDOFF.md`. Read that handoff first (and the study README).

## Acceptance criteria

- Node count 1 vs 3 with a separately labelled equal-total arm; never pooled
- RF=1 vs RF=3, plus a lagging-follower variant; the ticketing invariant holds at every level
- Per-endpoint distribution recorded for every client run; a single-endpoint arm is labelled a control
- The unchecked read-modify-write control is seen to oversell

## Constraints

- Take the benchmark lock; one matrix at a time; never measure in parallel.
- Correctness gates timing; record the environment, sizing and per-endpoint distribution.
- Produce a generated report (numbers only) and a signed analysis; never edit another analyst's file.
- Do not push or pull; commit and tag, and merge through the normal queue lifecycle.

NEXT MODEL: LOW
