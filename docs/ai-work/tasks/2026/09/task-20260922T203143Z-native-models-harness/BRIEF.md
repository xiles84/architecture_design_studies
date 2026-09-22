# Execution brief — Build the Study 06 native-model harness, pins and correctness gate

Implements part of `studies/06-native-models/HANDOFF.md`. Read that handoff first (and the study README).

## Acceptance criteria

- Containerfiles and probe scripts for MongoDB 8.0, ScyllaDB 6.2 and Valkey 8.1; exact tags and image digests recorded in versions.env and manifests
- Deterministic dataset loader and per-technology native correctness gate implementing Study 04 INV-1..INV-13 plus the native check
- Every design has a negative control that is seen to fire
- A missing or moved upstream tag is a blocker, not a silent upgrade

## Constraints

- Take the benchmark lock; one matrix at a time; never measure in parallel.
- Correctness gates timing; record the environment, sizing and per-endpoint distribution.
- Produce a generated report (numbers only) and a signed analysis; never edit another analyst's file.
- Do not push or pull; commit and tag, and merge through the normal queue lifecycle.

NEXT MODEL: LOW
