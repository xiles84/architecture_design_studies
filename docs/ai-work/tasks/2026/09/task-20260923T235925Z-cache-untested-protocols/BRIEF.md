# Execution brief — Measure the unmeasured cache protocol and deletion tactics (conditional)

This task is published **proposed**, not ready. The cache brainstorm conclusion named these experiments as the
candidates worth running, but also said comparative guidance may not be needed. Release it only if the book or a
reader needs a comparative answer the current evidence cannot give.

Read the cache brainstorm conclusion first
(`docs/ai-work/brainstorms/records/2026/09/brainstorm-20260923T172022Z-cache-strategy-coverage-by-schema-control/CONCLUSION.md`),
then `studies/05-cache-consistency/HANDOFF-AMENDMENT-01-V2.md`, `README.md` and `sql/README.md`.

## Goal

The corpus demonstrates one sufficient strict protocol for coordinated shared caches and leaves both the
alternatives and the deletion question unmeasured. Measure them under the same contract and gates as the
existing study.

## Decisions already made (do not re-litigate)

- Contract is Study 05's: a read beginning after acknowledgement of a write to the same key must not return an
  older committed value; overlapping reads are classified separately; impossible/uncommitted values are forbidden.
- The two arms are separable and may be split into successive runs; do not merge them into one matrix.
- A cache-side fence is not to be assumed necessary; the comparison must be fair to a version-validated hit and
  to an outbox-delivered invalidation, with the same correctness oracle.
- Negative controls must fire. A control that does not fire fails the run.

## Acceptance criteria

See the task spec. Minimum: Arm A protocol comparison under equal-total resources; Arm B deletion/resurrection
audit with a paused filler; every arm with fired controls; same-run wrong-read, eviction and placement records;
one matrix at a time under the benchmark lock.

## Constraints

- This is a measurement: take `ads-run-lock`, one matrix at a time, never measure in parallel; record the
  environment, sizing, per-endpoint distribution and seeds.
- Correctness gates timing. Do not weaken a gate to make a cell pass.
- Never edit another analyst's analysis, report or digest; produce new reports and your own signed analyses.
- Commit and tag; never push or pull; merge through the queue lifecycle.

## Escalate only if

- a required comparison cannot be built without changing an existing design's semantics or a published claim
  (that is material, not a local implementation choice).

NEXT MODEL: HIGH
