# Execution brief — Run the Study 07 placement and routing arms

Implements part of `studies/07-topology/HANDOFF.md`. Read that handoff first (and the study README).

## Acceptance criteria

- Colocated vs non-colocated hold logical data, operations, replication and resources constant
- Actual tablet/leader distribution recorded for every cell; no record means gap
- Routing arms: single endpoint (control), round-robin, least-connections, with the skew shown
- A deliberately stale replica control is seen to serve stale data

## Constraints

- Take the benchmark lock; one matrix at a time; never measure in parallel.
- Correctness gates timing; record the environment, sizing and per-endpoint distribution.
- Produce a generated report (numbers only) and a signed analysis; never edit another analyst's file.
- Do not push or pull; commit and tag, and merge through the normal queue lifecycle.

NEXT MODEL: LOW
