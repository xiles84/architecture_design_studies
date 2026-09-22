# Result — book evidence registry and validators

**Task:** `task-20260922T025912Z-evidence-registry`
**Branch:** `repo/book-evidence-registry`, base `f82f5a3`
**Required tag:** `repo/book-evidence-registry-v1`

## Decision Log

1. **The registry is curated; the extractor only copies.** `tools/evidence extract` dumps
   analysis front-matter verbatim as candidates. A claim's statement is copied from a signed
   analysis; strength, cells and limits are curated and never invent a comparison. This is what
   "no new interpretation" means here.
2. **Strength is a contract, not a label.** `repeated_controlled` requires ≥ 2 trials,
   `single_run_directional` exactly 1, `mechanism_supported` ≥ 1, `gap`/`analogy` none. The
   validator rejects a mismatch, so the run's own trial count governs the level.
3. **Evidence levels are the WORKFLOW ladder** (`repeated_controlled`, `single_run_directional`,
   `mechanism_supported`, `analogy`, `gap`, `invalid`), reused rather than invented.
4. **A failed cell can never be a winner.** Every numeric claim lists its cells and winners; the
   validator rejects a winner whose status is `failed`/`invalid` and checks the declared
   `failed_cells` count. Negative fixture `failed-cell-winner.json` proves it.
5. **Stale evidence is a digest mismatch.** `inputs_digest` must appear verbatim in both the
   report and its signed analysis; a re-run changes the digest, so the mismatch is detected
   without trusting a date. Fixture `stale-digest.json`.
6. **Cross-study ratios are gated.** `comparability: cross-study` needs an existing
   `comparability_record`; no such claim is seeded yet, so no cross-study number can slip in.
   Fixture `cross-study-without-record.json`.
7. **Supersession is symmetric and explicit.** A superseding claim names the claim it replaces,
   which must name it back; `--list-inputs` excludes superseded claims. Fixture
   `superseded-not-marked.json`.
8. **The initial set is the current signed headlines plus the visible gaps** — 12 numeric claims
   (one per current analysis headline or headline sentence) and 4 gap records (Study 04 missing
   designs, Study 05 unmeasured regimes, Study 03 hold-window sizing, Study 02/03 validation
   open). No claim is superseded.
9. **No external dependency.** The schema subset is evaluated in-tree so validation is offline
   and reproducible.

## Evidence

- `book/evidence/claims.json` (16 claims), `book/evidence/claims.schema.json`,
  `book/evidence/README.md`.
- `tools/evidence/` — `main.go`, `internal/evidence/schema.go`, `internal/evidence/registry.go`,
  `internal/evidence/registry_test.go`, `testdata/*.json`, `README.md`.

## Validation performed

- `go test ./...` green (real registry validates cleanly; all five negative fixtures fail for
  the intended reason; the mini-schema rejects an out-of-enum value).
- `go run . validate --repo ../..` → `evidence: 16 claims, 0 errors, 16 book inputs`.
- Every referenced report, analysis, environment page, results directory and run tag resolves;
  every digest appears in both its report and analysis.

## Residual risks / limitations

- The seeded claim set is the headline claim of each current analysis, not every sentence in
  every report. Later book chapters must add claims (same process) rather than quoting numbers
  that are not yet in the registry.
- Cell names in `provenance.cells` are curator-recorded summaries of the run's design ids; the
  validator enforces their statuses and winner rule but does not re-derive them from result
  JSON. A future task can parse `results/*/*.json` to cross-check cell status.
- The registry covers Studies 01–04 at their current evidence level; when the Study 02/03 HIGH
  validation lands, those claims' `limits` should be updated (by a new claim, not an edit to a
  signed artifact).
