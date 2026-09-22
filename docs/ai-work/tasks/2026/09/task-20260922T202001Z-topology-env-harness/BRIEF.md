# Execution brief — Build the Study 07 two-host environment, balanced-endpoint and fault harness

Implements part of `studies/07-topology/HANDOFF.md`. Read that handoff first (and the study README).

## Acceptance criteria

- A second host environment page documents CPU, memory, disk, kernel, network path, clock source and failure modes
- Client endpoint round-robin/load-balancing instrumentation records per-endpoint connection and operation distribution
- netem latency/loss helper and a node-kill/pause helper exist and refuse to run without the benchmark lock
- The correctness gate is unchanged and runs before any timing

## Constraints

- Take the benchmark lock; one matrix at a time; never measure in parallel.
- Correctness gates timing; record the environment, sizing and per-endpoint distribution.
- Produce a generated report (numbers only) and a signed analysis; never edit another analyst's file.
- Do not push or pull; commit and tag, and merge through the normal queue lifecycle.

NEXT MODEL: LOW
