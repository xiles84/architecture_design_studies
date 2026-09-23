# Execution brief — Record derived-value staleness and rollup lag in the Study 04 cadence phase

Implements part of `studies/04-configuration-portal/HANDOFF-AMENDMENT-01-V2.md`. Read that handoff first (and the study README).

## Acceptance criteria

- The cadence phase stamps each derived value and each read so the age of the answer is recorded
- Staleness distribution and rollup lag are reported for trigger versus application maintenance at the 120-entry tier
- The correctness gate still passes and the change is re-verified before any number is published
- One matrix at a time under the benchmark lock

## Constraints

- Take the benchmark lock; one matrix at a time; never measure in parallel.
- Correctness gates timing; record the environment, sizing and per-endpoint distribution.
- Produce a generated report (numbers only) and a signed analysis; never edit another analyst's file.
- Do not push or pull; commit and tag, and merge through the normal queue lifecycle.

NEXT MODEL: LOW
