# Execution brief — Run the Study 08 search-copy arm

Implements part of `studies/08-analytics-read-models/HANDOFF.md`. Read that handoff first (and the study README).

## Acceptance criteria

- PostgreSQL full-text/trigram search copy kept in sync by trigger and by application, measured separately
- Recent-activity and top-N workloads measured with query gain, staleness, storage and correctness
- Search copy labelled as in-engine; any native search engine is a separate labelled extension
- Duplicate/lost-copy control fires

## Constraints

- Take the benchmark lock; one matrix at a time; never measure in parallel.
- Correctness gates timing; record the environment, sizing and per-endpoint distribution.
- Produce a generated report (numbers only) and a signed analysis; never edit another analyst's file.
- Do not push or pull; commit and tag, and merge through the normal queue lifecycle.

NEXT MODEL: LOW
