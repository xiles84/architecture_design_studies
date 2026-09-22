# Execution brief — Run the Study 07 network and failure arms

Implements part of `studies/07-topology/HANDOFF.md`. Read that handoff first (and the study README).

## Acceptance criteria

- netem latency at 1, 10 and 50 ms RTT plus one packet-loss level, with a no-op injection control
- Failure arms: kill one node, pause one node, drop one network link; each records the failover event
- A no-op injection or a kill with no recorded failover is a harness/coverage gap, not a finding
- Per-node and equal-total budgets reported for every arm

## Constraints

- Take the benchmark lock; one matrix at a time; never measure in parallel.
- Correctness gates timing; record the environment, sizing and per-endpoint distribution.
- Produce a generated report (numbers only) and a signed analysis; never edit another analyst's file.
- Do not push or pull; commit and tag, and merge through the normal queue lifecycle.

NEXT MODEL: LOW
