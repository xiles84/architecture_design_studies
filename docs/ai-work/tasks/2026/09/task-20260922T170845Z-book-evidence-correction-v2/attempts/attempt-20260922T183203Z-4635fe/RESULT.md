# Result — corrected book evidence registry v2

**Task:** `task-20260922T170845Z-book-evidence-correction-v2`
**Branch:** `repo/book-evidence-correction-v2`
**Required tag:** `repo/book-evidence-registry-v2`
**Capability:** HIGH session executing a LOW task (model `deepseek-flash`)
**Benchmark:** none — artifact-only over committed evidence

## What was delivered

1. **The sole active v2 registry** at `book/evidence/v2/claims.json` — 23 atomic claims with a
   structured support key per numeric claim, run-level status separate from support cells, control
   scope, multi-run anchors, three-dimensional trial counts and the `legacy_provenance_incomplete`
   mark on Study 01's untagged run.
2. **The correction ledger** `book/evidence/v2/CORRECTION_LEDGER.md`: every changed, split or retired
   v1 claim with its reason, and proof that v1 is byte-identical to `repo/book-evidence-registry-v1`.
3. **The disclosure set**: `confounds.json` (the six required confounds), `coverage.json` (coverage
   matrix with review state, including the new placement, native-datastore and real-network gaps) and
   `SYNTHESIS_HANDOFF.md` (claims → fixed six-family taxonomy + decision models).
4. **A hardened validator**: `tools/evidence validate` defaults to v2 and dispatches on
   `schema_version`; `validate-v2` enforces the semantic resolver; `v1-diagnostic` recomputes the
   frozen v1 audit; `internal/evidence/v2_test.go` adds positive and four negative fixtures plus the
   `12/11/43/36` regression.
5. `book/evidence/README.md` and `tools/evidence/README.md` document one active source; `CONTEXT.md`
   records the active version, the tag target and the gate that `book-synthesis-v1` may not be
   released until this task is integrated.

## Validation (all in the pinned Go image, from committed files)

```text
go vet ./...                    clean
go test ./...                   ok (positive + 4 negative fixtures + v1 regression)
validate (default v2)           evidence: 23 claims, 0 errors, 23 book inputs (active source book/evidence/v2/claims.json)
validate (historical v1)        evidence: 16 claims, 0 errors, 16 book inputs
v1-diagnostic                   12 numeric claims, 11 with a missing name, 36/43 names unresolved
v1 unchanged vs tag             git diff repo/book-evidence-registry-v1 -- <v1 files> is empty
```

## Corrections v1 → v2 (summary)

- claim-01 split into three, cell names replaced with the report's own labels, legacy provenance
  marked; conf-01 attached.
- claim-04 `failed_cells` corrected 0 → 2 of 42; control scope recorded.
- claim-05 split: answerability + 33-35% storage kept, "8x load" removed to a gap.
- claim-07 range corrected to the run's own 1.7-5.7x; "PostgreSQL unsettled" kept.
- claim-10 split and narrowed: 1.76x is an upper bound; trigger cost kept.
- claim-12 range widened to the full -14.2%..+12.3%, still labelled noise.
- claim-gap-04 retired.
- New: `v2-15-three-instance-staleness` (corrected model attribution), and three gap claims.

## Conclusion

Registry v1 is immutable and cited; the v2 package is the book's single active numeric source. No
measurement was performed or authorized.
