# Execution brief — Run the Study 06 wide-column-versus-relational comparison

Implements part of `studies/06-native-models/HANDOFF.md`. Read that handoff first (and the study README).

## Acceptance criteria

- ScyllaDB and the PostgreSQL baseline run identical logical data and the five operations
- Partition/clustering design documented; point and range access both measured
- Correctness gate passes and the negative controls fire; within-run comparison only

## Constraints

- Take the benchmark lock; one matrix at a time; never measure in parallel.
- Correctness gates timing; record the environment, sizing and per-endpoint distribution.
- Produce a generated report (numbers only) and a signed analysis; never edit another analyst's file.
- Do not push or pull; commit and tag, and merge through the normal queue lifecycle.

NEXT MODEL: LOW
