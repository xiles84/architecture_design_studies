# Result — SELECT synthesis

**Task:** `task-20261002T131729Z-select-synthesis`. **Attempt:**
`attempt-20261005T133549Z-3446d3`. **Agent:** Claude Opus 5.5 / Claude Code desktop / HIGH /
analyst; effort xhigh. **Required tag:** `study-01/select-v6-analysis`.

## Delivered

| Artifact | Path |
|---|---|
| Signed synthesis analysis: P1–P12 verdicts, selection matrix, limits, required-comparison coverage | `studies/01-charity-tree/reports/analyses/20261005T-select-synthesis--claude-opus-5-5--2026-10-05.md` |
| Generated evidence tables (numbers only, 60 runs) | `studies/01-charity-tree/reports/20261005T-select-synthesis-evidence.md` |
| Evidence generator, claim spec, registration helper, anchors | `docs/select-query-templates/synthesis-20261005/` |
| 18 registered claims; two gaps narrowed; ledger and coverage | `book/evidence/v6/claims.json`, `SUPERSESSION_LEDGER.md`, `coverage.json` |
| Validator run-tag namespace and tests | `tools/evidence/internal/evidence/v2.go`, `v6_test.go` |
| Book handoff for HIGH `select-book-high` | `docs/select-query-templates/BOOK-HANDOFF-20261005.md` |
| Decision Log | `PROGRESS.md` beside this file |

Inputs digest of the 60 runs: `baaaab3622c5170faf371f64a51bf4cf2c9a02aa644fcb8c82c143016211f0c6`.

## Verdicts

- **Confirmed:** P1, P2, P5, P7, P12.
- **Confirmed, refined:** P3. The cost depends on whether the filter is outside or inside the
  CTE.
- **Confirmed in 3 of 4 regimes:** P4. Falsified at the PostgreSQL 1% consumer.
- **Confirmed on PostgreSQL only:** P6 and P8. P6 is not observed on YugabyteDB; P8 is falsified
  there.
- **Confirmed on YugabyteDB:** P9. On PostgreSQL its small side is within drift.
- **Falsified in the tested shapes:** P10.
- **Confirmed in direction:** P11. Its flat-cost component was not tested.

## Validation

- `tools/evidence validate`: 49 claims, 0 errors; run tags resolved with git inside the container.
- `go test -count=1 ./...` (tools/evidence): ok.
- `book/evidence/check.sh` against v6: OK.
- Evidence report regenerated twice: byte-identical.
- No database or benchmark was started; `ads-run-lock` was never taken.

## Limits

- YugabyteDB access and existence numbers come from scan plans (`v6-sel-gap-01`).
- The three non-join families have no cluster arm.
- One laptop, one sequential reader, five samples per cell.
- Aggregate drift is 29–74%.
- Every crossing is bracketed, not located.
- The tracked book PDF still reflects the previous registry digest; `select-book-high` rebuilds it.

## Review

Labelled HIGH self-review (same model and session), allowed by the brief when no independent
reviewer is available. The checks re-derived every registered figure from the generated tables.
An independent analyst may disagree in a signed companion.

## Next

HIGH `task-20261004T142623Z-select-book-high`: publish the reviewed cases and the selection guide
from `BOOK-HANDOFF-20261005.md`.
