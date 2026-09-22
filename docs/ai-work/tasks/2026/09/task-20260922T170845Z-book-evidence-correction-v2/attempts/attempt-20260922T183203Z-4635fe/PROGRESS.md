# Progress — corrected book evidence registry v2

**Task:** `task-20260922T170845Z-book-evidence-correction-v2`
**Branch:** `repo/book-evidence-correction-v2`
**Attempt:** `attempt-20260922T183203Z-4635fe` (claim `claim-e4004f9358131d07`, epoch 1)
**Capability/role:** HIGH session executing a LOW task / executor (model `deepseek-flash`)
**Required tag:** `repo/book-evidence-registry-v2`

Artifact-only work over committed evidence. No database started, no benchmark run, no new design.

## Decision Log

1. **v2 is a new directory, v1 untouched.** `book/evidence/v2/` holds the active package;
   `book/evidence/claims.json` and `claims.schema.json` are byte-identical to tag
   `repo/book-evidence-registry-v1` (`git diff` empty; hashes in `CORRECTION_LEDGER.md`).
   The correction ledger is the single attributed change record, as the conclusion requires.
   Confidence: High.
2. **The structured support key resolves to the report's own label, not a curated alias.**
   `support[].resolves_to` must appear verbatim in the primary anchor's cited report; the validator
   enforces it. This is the semantic resolver v1 lacked. Confidence: High.
3. **Run level is a separate field from support cells.** `run_level` carries the run's own cell
   totals, failures and control scope; this fixes claim-04's `failed_cells: 0` → 42 cells / 2 failed
   and avoids the scope conflation critique-02 identified. Confidence: High.
4. **Split compound claims rather than qualify them loosely.** claim-01 → v2-01/02/03, claim-05 →
   v2-07 + a load-multiplier gap, claim-10 → v2-12/13. Confidence: High.
5. **Apply the three 2026-09-22 review corrections inside v2** (claim-07 range 1.7-5.7x, claim-05
   "8x load" removed, all-green three-instance attribution corrected in a new claim). The reviews
   are cited in the claims' `review.artifact`. Confidence: High.
6. **Gap claims carry a `gap_basis` and no support key.** "As applicable" is honoured: gaps name the
   artefact that establishes the absence; numeric claims must have ≥1 resolving support key.
   Confidence: High.
7. **Default selection moved to v2 by dispatching on `schema_version`.** `tools/evidence validate`
   now defaults to the v2 paths and prints which source it read; v1 is reachable only by explicit
   `--claims`/`--schema`. Confidence: High.
8. **The v1 43-name audit is a Go regression fixture, not a one-off script.** `v1-diagnostic` and
   `TestV1DiagnosticRegression` pin 12/11/43/36. Confidence: High.

## What was produced

- `book/evidence/v2/{claims.json,schema.json,confounds.json,coverage.json,CORRECTION_LEDGER.md,SYNTHESIS_HANDOFF.md}`
- `book/evidence/README.md` (active-source section), `tools/evidence/README.md`
- `tools/evidence/internal/evidence/v2.go`, `v2_test.go`; `main.go` `validate-v2`, `v1-diagnostic`,
  and `validate` re-pointed at v2
- `CONTEXT.md` active source, gate rule, and closed gap 2

## Residual risks / limitations

- v2 has 23 claims; it is a correction of v1's 12 numeric + 4 gap claims, not a new survey of all
  56 reports. A defect visible only inside an unread report body may remain.
- Support keys resolve to report tokens; they do not re-derive the measured value from result JSON.
  A future task could add that cross-check (as the evidence-registry result already noted).
- Five numeric claims (v2-01/02/03) share the one legacy Study 01 run and remain
  `legacy_provenance_incomplete`.
- The coverage matrix records `planned` work whose tasks are not yet released; it will drift if the
  protocols change.
