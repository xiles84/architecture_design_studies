# HIGH review of the SELECT synthesis — labelled self-review

Reviewer: Claude Opus 5.5 (`claude-opus-5-5`) / Claude Code desktop / HIGH; effort xhigh;
2026-10-05. **This is a self-review.** The same model and session authored the synthesis. The
brief allows a labelled self-review when no independent reviewer is available. An independent
analyst may disagree in a signed companion.

**Verdict: approve for integration.**

## Checks

- **Decision Log** (`attempts/attempt-20261005T133549Z-3446d3/PROGRESS.md`): eleven entries, each
  with reason, confidence and affected scope. Two are lowest-confidence:
  - anchor environment naming (Medium): each affected claim states it in a limit;
  - strength assignment (Medium): `repeated_controlled` wherever no direction was recorded before
    the runs.
- **Evidence tables:** generated twice from the 60 tagged runs, byte-identical. Inputs digest
  `baaaab36…`.
- **Registry:**
  - `tools/evidence validate`: 49 claims, 0 errors, including run-tag existence, report and
    analysis digests, results directories and support tokens;
  - `go test -count=1 ./...`: ok;
  - `book/evidence/check.sh`: OK.
- **Figures:** all 27 numeric sequences in the claims' support fields and every `resolves_to`
  token appear verbatim in the evidence tables. Statement ranges were re-derived from the same
  tables.
- **Scope:**
  - No database was started.
  - No earlier analysis, discussion, report or result was edited.
  - Frozen registries v1–v5 are untouched.
  - The two narrowed gaps keep their anchors, trials and confounds and are recorded in the ledger
    and the carry-forward test.
- **Amendment (2026-10-04 publication route):** every requested item is covered, and the narrowed
  aggregate verdicts are not reinstated. The requested items are:
  - P1–P5/P12;
  - fixed versus blocked predicate placement (filter inside versus outside the CTE);
  - the verified default CTE policy (plan-equal to `MATERIALIZED` from two references);
  - seeded interleaving;
  - first-cell drift;
  - gaps.

## Residual risks

- Self-review cannot catch an error the author does not see. The highest-value independent check
  is whether the scan-only YugabyteDB label (`v6-sel-gap-01`) and the `repeated_controlled`
  strengths are applied consistently in the book.
- The tracked PDF predates the 18 claims until `select-book-high` rebuilds it.
