# Execution brief — corrected book evidence registry v2

## Goal

Create one versioned, attributed correction package that becomes the sole evidence source for
numeric book prose. Preserve the accepted v1 registry unchanged at
`repo/book-evidence-registry-v1`. This is artifact-only work over committed evidence: do not run
a benchmark, start a database, add a design, or create a new scientific study.

This task implements the conclusion of
`brainstorm-20260922T155627Z-book-v1-evidence-defects-and-model-additions`. It does not redesign
the original registry task, repeat the three existing HIGH reviews, rebuild the Typst toolchain,
write the book chapters, or perform the release review.

## Current state and fixed decisions

- Registry v1 is schema-valid but not eligible as the book's direct numeric source. The
  conclusion mechanically found that 11 of 12 numeric claims have at least one declared name
  absent from their cited report and that 36 of 43 names do not resolve verbatim. Claim 01 also
  shifts D1-D3 and reverses D4/D5 relative to the cited corpus.
- V1 remains immutable and cited. Do not edit `book/evidence/claims.json`,
  `book/evidence/claims.schema.json`, any protected evidence artifact, or an existing tag.
- The correction is a new version under `book/evidence/v2/`. The book and validation tooling
  have exactly one active source; there is no silent fallback to v1 and no parallel narrative
  claim database.
- Keep the owner's fixed six major families. The brainstorm's model additions are decision axes
  and concept material beneath them, not new major families.
- Keep every evidence level conservative. A correction may narrow, split, qualify or invalidate
  a claim; it may not promote one.
- Run tags identify producing code state and need not contain the later committed result files.
  Document that provenance model rather than treating tag ancestry as a defect.

## Work

1. **Freeze and cite v1.** Add a signed/attributed correction ledger that identifies the v1 tag,
   records every changed/split claim id and reason, and proves the root v1 JSON/schema files are
   unchanged. Do not copy claims into an independent prose database.
2. **Create the sole active v2 registry.** Decide the exact internal schema within the required
   `book/evidence/v2/` location. Update `book/evidence/README.md` and the tool's default selection
   so book inputs and validation resolve only v2 unless a caller explicitly requests historical
   v1 inspection.
3. **Atomize support.** Split compound claims where one status/winner cannot honestly qualify all
   clauses. Give each atomic claim a structured support key: topology, design, scenario or
   operation, and trial/aggregate identity as applicable. Resolve it against committed report or
   result artifacts, not by a loose display alias or coincidental substring.
4. **Separate scopes.** Record complete run totals/status separately from selected support cells.
   Preserve partial, skipped, diagnostic, negative-control, expected-rejection, failed and invalid
   outcomes. A winner must resolve to the exact eligible support cell; do not require a selected
   subset's failed count to equal the whole run's failures.
5. **Make controls and repetitions explicit.** Record negative-control scope as same-run fired,
   present-but-not-exercised, inherited, absent, or justified not-applicable. Separate whole-run
   replication, fresh-load/cell trials and inner race iterations. Keep conservative strength when
   independence is ambiguous.
6. **Restore provenance.** Supply all nine extant producing tags omitted by numeric v1 claims,
   plus applicable gap-record tags; record commit/digest/environment/topology/report/analysis
   links; support multiple anchors for genuinely multi-run conclusions. Mark claim 01
   `legacy_provenance_incomplete` and never invent its commit or tag.
7. **Publish the reader-facing qualifications.** Put the brainstorm conclusion's six confounds
   beside the affected claim evidence. Add a coverage matrix with `measured`,
   `measured-but-confounded`, `planned`, `gap`, and justified `not_applicable`, including review
   state. Treat Study 01 D7/D24 and Study 03 L3 as layout/placement intent with measured effects,
   not verified physical colocation; retain Study 02's recorded placement gap. Add the missing
   native-datastore-family and real-network/balanced-endpoint gap records without replanning their
   existing protocol tasks.
8. **Produce a synthesis handoff, not chapters.** Map claims to the fixed six-family taxonomy and
   the selected analytic models: answerability-before-performance (explicitly a derived rule),
   multidimensional design vector, derive to index to materialize to copy/cache ladder,
   contention surface/arbitration, query-boundary portfolio, evidence-confidence cards, and
   failure-mode catalogue. Label every entry direct evidence, analogy, or gap and include the
   exact run/report/analysis/digest/tag references available. Do not quote cross-study rankings.
9. **Harden validation.** Add semantic resolver and negative tests for the failure modes above.
   Retain the 11-of-12 and 36-of-43 v1 diagnostic as a regression fixture or reproducible audit,
   but require zero unresolved structured support keys in the active registry.
10. **Record decisions and hand off cleanly.** Log meaningful schema/resolver choices and any
    narrowed or invalidated claim. Update `CONTEXT.md` with the active version, tag target, and
    explicit rule that the proposed book-synthesis task must not be released until this task is
    completed and integrated.

## Acceptance and validation

All task-spec acceptance criteria apply. In addition:

- Use Podman for all Go/tool validation, following repository methodology and the existing
  evidence tool's pinned environment.
- Validate from a clean committed tree; run queue audit; inspect generated tables or book inputs
  for identity swaps, truncated limits, duplicate sources and unreadable output.
- The HIGH reviewer must review the Decision Log first and independently rerun both the v1
  mechanical baseline and the v2 semantic resolver checks.
- If exact v2 package selection cannot be made compatible with the existing book toolchain,
  escalate before introducing a second active source. The allowed fallback is one signed audit
  as the sole book input with v1 excluded, exactly as the conclusion states.
- Finish through the normal queue lifecycle: commit explicit owned paths, submit, obtain HIGH
  approval, integrate into local `main`, create annotated tag `repo/book-evidence-registry-v2`,
  verify reachability and complete. Never push or pull.

## Genuine escalation triggers

Escalate only if committed artifacts cannot uniquely resolve an atomic claim, if satisfying one
existing HIGH review contradicts another, if the active-source switch requires altering a
protected v1 artifact or immutable existing task, or if a proposed correction would strengthen a
claim. Local schema shape, fixture organization and renderer wording are LOW decisions: decide,
log and continue.

NEXT MODEL: LOW
