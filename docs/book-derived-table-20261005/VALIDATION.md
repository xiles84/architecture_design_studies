# Edition 3 preview 5 — derived-table clarity (validation and labelled self-review)

Claude Opus 5.5 / Claude Code desktop / HIGH; effort xhigh; 2026-10-05.
Task `task-20261005T160000Z-book-derived-table-clarity`. No database or benchmark was started.

## Owner question

"The book does not make clear to me when (or if there is none) situation when filtered derived
table queries outperform other forms."

## Answer now in the book

- **Filtered derived table** (case p. 56, guide p. 60): it ran the same plan as `JOIN … WHERE` at
  every selectivity, on both engines and in both phases. It never outperformed, and no measured
  situation favours it (`v6-sel-01`).
- **When a derived table is more than a spelling:**
  - theory, p. 21: number-free;
  - case, p. 56: labelled inferred.
- **Pre-aggregating derived table** (case p. 57, guide p. 60): `v6-sel-18`, registered from the
  existing synthesis evidence with no new measurement.
  - It lost at 1–1,000 donors on PostgreSQL.
  - It dipped below 1 only at 10,000, inside the measurement drift.
  - It never won on YugabyteDB.
  - The two-child-collection case it is meant for is marked not measured.

## Checks

- `tools/evidence validate`: 50 claims, 0 errors. `go test -count=1 ./...` (tools/evidence): ok.
- `book/evidence/check.sh`: OK.
- Canonical `book/build.sh`, run in WSL from clean commit `5326ef3d`, verified:
  - 99 pages, 4/4 fonts, 50/50 claims indexed;
  - edition "Edition 3 preview 5";
  - PDF sha256 `5648b27318e88932f8d3301567f7378f3e17a8be575c1437392748a1bf6bf623`
    ([log](build.log)).
- Pages 21, 56, 57 and 60 rendered and inspected (`proof/`): no clipping, split ids or stretched
  lines.

## Self-review

Same session as the author; labelled. Every number on the new lines is copied from the
`v6-sel-01` and `v6-sel-18` statements. Approve.
