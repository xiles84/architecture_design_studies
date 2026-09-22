# Review — corrected book evidence registry v2

**Task:** `task-20260922T170845Z-book-evidence-correction-v2`
**Reviewer:** HIGH session, model `deepseek-flash`, role `reviewer`
**Reviewed:** `attempt-20260922T183203Z-4635fe` (Decision Log in `PROGRESS.md`, result in
`RESULT.md`), submitted at `ccd2698`.
**Decision:** **approve for integration.** No corrective handoff.

## Decision Log reviewed first

| # | Decision | Verdict |
|---|---|---|
| 1 | v2 in a new directory; v1 untouched, ledger as change record | Accept; independently confirmed `git diff repo/book-evidence-registry-v1` over the v1 files is empty. |
| 2 | support key resolves to the report's own label | Accept; this is the semantic resolver the conclusion required. |
| 3 | run level separate from support cells | Accept; it fixes claim-04's false `failed_cells: 0`. |
| 4 | split compound claims | Accept; claim-01/05/10 were each carrying two decisions. |
| 5 | apply the three 2026-09-22 review corrections | Accept; verified in `v2-07`/`v2-09`/`v2-15` and the ledger. |
| 6 | gap claims carry `gap_basis`, no support key | Accept; "as applicable" honoured without weakening numeric claims. |
| 7 | `validate` defaults to v2 by schema_version dispatch | Accept; the printed source makes the active registry unambiguous. |
| 8 | v1 audit as a Go regression fixture | Accept; it is what keeps the failure mode from recurring silently. |

All eight are within the brief and reversible. None changes a scientific claim's measured content;
the strength of every claim is the same or narrower than v1.

## Independent checks re-run (the brief requires both)

- `tools/evidence validate` (default v2): **23 claims, 0 errors**, active source
  `book/evidence/v2/claims.json`.
- `tools/evidence v1-diagnostic`: **12 numeric claims, 11 with a missing name, 36/43 names
  unresolved** — matching the brainstorm conclusion's computed baseline.
- `git diff repo/book-evidence-registry-v1 -- book/evidence/claims.json book/evidence/claims.schema.json`
  is empty; v1 is frozen exactly as required.
- The negative fixtures in `v2_test.go` (unresolvable support key, legacy-flag mismatch,
  unreferenced confound, unknown supersession) each fail for the intended reason, and the v1
  regression test pins 12/11/43/36.

## Result reviewed against the acceptance criteria

- **Sole active source with a cited predecessor**: yes — `book/evidence/v2/`, ledger cites
  `repo/book-evidence-registry-v1`, no second active claim database.
- **Atomized claims with source-resolving support keys**: yes; support keys use the reports' own
  labels (`D2 indexed`, `P3 CAS`, `l:na:mem:asd:str`, …), enforced by the validator.
- **Run-level vs support-cell status, controls, multi-run anchors, trial dimensions**: all present;
  claim-04 corrected to 42 cells / 2 failed; Study 01 untagged run marked
  `legacy_provenance_incomplete`.
- **Confounds and coverage published beside the numbers**: the six confounds are referenced by the
  affected claims and the coverage matrix carries review state plus the placement, native-datastore
  and real-network gaps.
- **Synthesis handoff, not chapters**: present and labelled direct/analogy/gap.
- **Hardened validation with the v1 diagnostic retained**: yes.
- **Context updated with the active version, tag target and the synthesis gate**: yes.
- **No measurement**: confirmed; the task touched registry, validator and prose only.

## Conditions on the outcome

The v2 package is the only source the book may read. The three corrections from the 2026-09-22
reviews must be kept (`v2-09` 1.7-5.7x, `v2-gap-01` load multiplier dropped, `v2-15` corrected
three-instance attribution). The evidence-validation residual risk — support keys resolve to report
tokens, not re-derived result values — is recorded and is acceptable for this task.

## Integration

Approve and integrate into local `main`; tag the integrated state `repo/book-evidence-registry-v2`.
