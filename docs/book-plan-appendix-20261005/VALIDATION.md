# Edition 3 preview 6 — query-plan appendix (validation and labelled self-review)

Claude Opus 5.5 / Claude Code desktop / HIGH; effort xhigh; 2026-10-05.
Task `task-20261005T190500Z-book-plan-appendix`. No database or benchmark was started: every plan is
quoted from committed result files.

## Owner request

"I believe the tests we have performed analyze explains on the database queries. could we have an
appendix with them?"

## What the book now has

- **Appendix *Query plans behind the SELECT claims*** (pp. 75–87):
  - 25 exhibits holding 45 trimmed `EXPLAIN ANALYZE` plans;
  - a reading note per claim;
  - the trim rule and kept attributes, printed from the generated data;
  - per exhibit: run, run tag, result path and sha256.
- **Coverage:** every measured SELECT claim, `v6-sel-01` to `v6-sel-18`.
  - `v6-sel-07` has no exhibit. Its nullable `NOT IN` is a negative control that records rows, not a
    plan; the appendix says so on p. 80, and `TestPlanAppendixCoversEveryMeasuredSelectClaim`
    encodes the exception.
  - `v6-sel-gap-01` is a gap, and the generator refuses gap claims.
- **Links:** the Charity template guide (p. 60) and *Reproduce a measured run* (p. 74) link to the
  appendix.

## How the plans are made trustworthy

- **Generator:** `tools/evidence plans` reads the authored `book/appendix/query-plans.spec.json` and
  writes `book/appendix/query-plans.json`. It refuses an exhibit when:
  - the run is not an anchor of the claim;
  - sha256 of the result file differs from the anchor's `inputs_digest`;
  - the cell is missing, invalid or has no plan;
  - the claim is a gap or not active;
  - a declared node-tree statement is contradicted by the plans.
- **Tests** (`plans_test.go`):
  - trim and shape unit tests;
  - a staleness gate (committed data equals regenerated data);
  - a coverage test;
  - four negative controls: unanchored run, contradicted shape, missing axis, gap claim.
- **First generation caught an expectation error.** The spec said the colocated and distributed
  180,000-row joins had the same node tree. They do not: the colocated run read `person` and `charity`
  through primary-key index scans. The spec now expects "different", and the note says what differs.

## Checks

- `go vet ./...`, `gofmt -l .` (clean), `go test -count=1 ./...` (tools/evidence): ok.
- `go run . plans --repo ../..`: current (25 exhibits, 45 plans).
- `go run . validate --repo ../..`: 50 claims, 0 errors.
- `book/evidence/check.sh`: OK.
- Canonical `book/build.sh`, run in WSL from clean commit `9bcdb642`, verified ([log](build.log)):
  - 112 pages, 4/4 fonts, 50/50 claims indexed;
  - edition "Edition 3 preview 6";
  - PDF sha256 `c88b724770bea23c3e65cf4a26de22fcaaf2dceac1aaee8f77b569074917c6c2`.
- Pages 60 and 75–87 were inspected in draft and canonical renders; [`proof/`](proof/) keeps pages
  60, 75, 76, 80 and 87.
  - Fixed during inspection: a cross-run comparison printed the wrong plan label, five repetitive
    shape bullets were grouped into one line, and a heading and badge were orphaned at a page foot
    (now sticky).
  - Remaining: long plans continue across a page break, which is acceptable for quoted output.

## Self-review

Same session as the author; labelled.

- Every figure in a reading note was checked against the plan printed beside it (rows removed, loops,
  request counts, execution times).
- The notes add no conclusion beyond the cited claim. Where a plan suggests a mechanism the claim
  calls a hypothesis (`v6-sel-10`'s per-probe storage requests), the note says the plan cannot settle
  it.

Approve.
