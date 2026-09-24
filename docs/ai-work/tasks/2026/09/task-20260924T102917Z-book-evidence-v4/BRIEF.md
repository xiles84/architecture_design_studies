# Execution brief — v4 evidence package and the retirement of the stale placement gaps

Implement finding **R04** of the external review of the *Edition 1* draft, plus the two findings that
share its cause (R03, the fixed registry path in prose) and the provenance hardening that survives
verification (R02's residual, not its severity claim).

## Goal

The active evidence package must stop asserting gaps that its own newer evidence closes. v3 retired
`v2-gap-04` correctly but carried `v2-gap-05-colocation-unverified` and
`v2-gap-07-real-network-balanced-endpoints` forward byte-identical while adding `v3-06` (the Study 05
key-locality pair, executed) and `v3-05` (three-node cells whose client spreads over all three
endpoints). v3's own `coverage.json` was updated for the endpoint spread; `claims.json` was not, and
the book printed `claims.json`.

## Decisions already made (do not re-litigate)

- **A new versioned package, not an edit of v3.** v1, v2 and v3 stay frozen and byte-unchanged; no
  tag moves. v4 carries the 27 still-active claims forward content-identically and records the two
  retirements in a ledger.
- **Retire, do not rewrite.** `v2-gap-05` and `v2-gap-07` leave `claims[]` entirely; their still-true
  clauses become `v4-gap-01-physical-placement-unverified` and
  `v4-gap-02-endpoint-distribution-partial`. Their closed clauses are recorded in
  `retired_predecessor_claims[].closed_dimensions`.
- **`v3-gap-01` becomes `partially_superseded`**: its placement-distribution clause is owned by the
  two successors, and `remaining_dimensions` names only hard-TTL expiry, medium scale and open-loop
  demand.
- **No new measurement, no new number.** The successor gaps cite existing runs and claim nothing.
- **Gap kinds.** `gap_kind` is one or more of `coverage`, `schema_limitation`, `instrument`,
  `unstable_measurement`. `v2-gap-02` (no design can answer the hold funnel) is a schema limitation,
  not a coverage gap, and the book must be able to badge it differently.
- **R02 is not a defect of the committed artefact.** Pages 1–2 carry real commit, describe, dirty
  state, clock, image digest, evidence digest and tree hash; the committed manifest records them.
  The only literal `unknown` is inside Typst's own version string. Do not fail release builds for
  missing provenance that is present; do stop a *preview* build from looking like a release one.
- **R01, the cache-chapter rewrite and the figure fixes are other tasks** (`book/cache-chapter-v2`,
  `book/figures-v2`). This task touches `book/concepts/cache-consistency.typ` in exactly one place,
  the registry-label substitution the new checker requires.

## Constraints

- No measurement, no database: the benchmark lock is not involved. The book build takes its own
  `ads-book-build-lock` volume for the shared `localhost/ads-book:1` image tag and must never touch
  `ads-run-lock`.
- Stage explicit paths; never `git add -A`. Do not commit `book/dist/`; build to a scratch path and
  delete it before committing.
- Never edit another analyst's report, analysis or digest.
- Keep `CONTEXT.md` and `LESSONS_LEARNED.md` current as part of this task, not afterwards.

## Acceptance criteria

See `task.json`. The gate that matters most: `tools/evidence validate-v4` reports 0 errors, `go test
./...` is green, `book/build.sh` verifies 29/29 claims, and a retired claim id fails the build with a
message naming its successor.

## Escalate only if

- retiring the two gaps would require changing a numeric claim, a report, or a signed analysis —
  that is material, not a local implementation choice;
- the v3 package cannot be validated unchanged after the new Go lifecycle rules land.
