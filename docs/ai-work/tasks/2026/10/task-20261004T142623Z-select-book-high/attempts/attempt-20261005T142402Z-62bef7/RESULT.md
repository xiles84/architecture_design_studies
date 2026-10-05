# Result — HIGH SELECT book publication

**Task:** `task-20261004T142623Z-select-book-high`. **Attempt:** `attempt-20261005T142402Z-62bef7`.
**Agent:** Claude Opus 5.5 / Claude Code desktop / HIGH / executor; effort xhigh.
**Required tag:** `repo/book-select-query-templates-high-v1`.

## Delivered

- **Charity donations case.** Measured SELECT variations replace the planned one:
  - join forms and the CTE fence;
  - CTE reuse;
  - membership and absence;
  - donor totals and the top donation;
  - pagination and disjunctions;
  - verified three-node placement;
  - a compact "Choosing a SELECT template" guide.

  All 18 `v6-sel-*` claims are cited with badges. The YugabyteDB scan-only gap is shown.
- **Theory.** Each SELECT strategy links to its measured variation. The pagination strategy gains a
  number-free access-path note. Placement theory names both placement cases.
- **Navigation, legend and other pages:**
  - the coverage map and its SELECT gap note;
  - the cache case's placement gaps;
  - the glossary (strength list and new terms);
  - the reading-guide legend;
  - the badge label for `repeated_controlled`.
- **Placement figure.** Regenerated from its source with the pinned renderer; it no longer claims
  placement is unverified everywhere.
- **Edition.** "Edition 3 preview 4", with a CHANGELOG entry.
- **Tracked PDF** `book/dist/data-architecture-reference.pdf` and its build manifest, from clean
  commit `bf03e408`.

## Validation

See `docs/book-select-high-20261005/VALIDATION.md`.

- Canonical build verified: 98 pages, 4/4 fonts, 49/49 claims indexed, evidence digest in the PDF,
  no unverified banner. PDF sha256 `be83c9422389a1b2e92265eae3d7291e7695b8a72119e6674a5cbe8a733b3799`.
- `check.sh` OK; 17 figures reproduce.
- Reader audit: 34/34 checks pass.
- Changed pages rendered and inspected; the build-1 findings were fixed in build 2.

## Limits

- The review is a labelled self-review by the same model and session.
- The SELECT figures carry their claims' limits: one laptop, one sequential reader, scan-only
  YugabyteDB access/existence plans, single-node families other than the join, and bracketed
  crossings.
- The old LOW `task-20261002T131729Z-select-book` still appears `ready` because queue v1 cannot
  retire it. Its mandatory amendment says not to execute it, and this task covers its deliverables.
