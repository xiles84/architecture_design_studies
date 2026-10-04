# Theory and applied cases — validation

## Scope and review decisions

The owner requested two clear reading paths. The implemented structure has six
theory groups and twenty strategy lessons, followed by five business cases.
Every strategy has a mechanism visual, short steps, fit/failure conditions,
cost rows, a case/coverage link and the eight-question reading checklist.
SELECT examples are illustrative SQL, not new measurements. Queue guidance
is conceptual; its delivery and throughput regimes are not measured here.

The applied-cases subagent's decision log and containerized source audit are
in `docs/book-applied-cases-20261004/`. HIGH reviewed the case definitions,
registered summary ranges, explicit uncertainty and all 31 claim citations.
Accepted its correction of the configuration document identifier and the
seating picture's acquisition/confirmation distinction. Parent changes keep
tables whole without altering the case results.

The active registry, reports, signed analyses and producing tags are unchanged.
Earlier book sources remain in the repository, but are excluded from the new
reading path. Full claim cards remain in the appendix. No database was started.

## Layout repairs found in development

- Removed empty conversion promises and duplicated older chapters.
- Kept fit tables and case summary tables together.
- Grouped costs, case link and checklist to avoid checklist-only continuations.
- Replaced the tall optimistic lesson diagram with a short conceptual path;
  retained its expanded original as optional appendix detail.
- Kept sink-side fencing separate from current source-generation validation.
- Corrected a template branch that initially rendered literal `else` text.

## Final canonical proof

- Source commit: `355db52837f7454b529329638a57da868141fbb8`; working tree clean at build.
- Edition 3 preview 2: **83 A4 pages**, 4/4 embedded fonts, 31/31 indexed claims.
- PDF SHA-256: `dd6aa64655e05b4d2e12fc4e2ab6b473d5814dd220690b0a114d8d64dc17e4a0`.
- All 17 registered figures reproduce through the pinned renderer.
- No compiler warnings, unresolved template tokens or unverified banner.
- Active evidence SHA-256 remains `6c5b2f80e98321c0c7cc90d4d6e42e8ebc91ef66ac6504bf4a732d0c5688e2db`.

The canonical `book/build.sh` ran through WSL with the existing jq and the shared
Windows Podman engine; no host toolchain was installed. The full build manifest
records the pinned Typst/renderer inputs and source hash.

Containerized Go checks passed:

1. Applied-case source audit: five cases, 31/31 active claim citations, no full
   cards in teaching, audited prose lines at most 60 words.
2. `check-reader.go`: six groups, twenty strategies, five cases, 31/31 claims
   in the rendered teaching portion, no legacy conversion promises or literal else branch.
3. Evidence resolver `validate-v6`: 31 claims, zero errors. Its first invocation
   mounted only the worktree and could not read the shared Git tags; the corrected
   read-only full-project mount passed. This was a mount/provenance access error,
   not missing producing tags.

All 83 final pages were rendered at 650 pixels and inspected in six numbered
contact sheets. High-resolution review covered pages 4, 13, 14, 21, 28, 39, 49,
50, 53 and 55: contents hierarchy, theory/cost/checklist, SQL, concurrency and
three applied-case layouts. Five representative proofs are preserved in `proof/`.
No clipping, overlaps, split summary/fit tables or orphan checklist rows remain.
Full claim cards stay in the appendix; their long audit text is not repeated in
the reader's teaching flow. The internal theory/case/claim links compile successfully.

The book source after the producing commit differs only in generated dist outputs.
No measurement, signed analysis, frozen package or producing tag changed.
Completion milestone: `repo/data-architecture-book-v3-preview-2`, local main only.
