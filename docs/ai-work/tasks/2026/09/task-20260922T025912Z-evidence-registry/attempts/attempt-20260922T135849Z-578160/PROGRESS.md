# Progress — book evidence registry

Attempt `attempt-20260922T135849Z-578160`, branch `repo/book-evidence-registry`, base `f82f5a3`.

- [x] Claimed through the queue.
- [x] Built `tools/evidence` (schema validator + provenance/digest/strength/winner/supersession/
      comparability/transfer rules) with tests and five negative fixtures.
- [x] Wrote `book/evidence/claims.schema.json` and seeded `book/evidence/claims.json` with 12
      numeric claims drawn verbatim from current signed-analysis headlines plus 4 gap records.
- [x] `go test ./...` green; `validate --repo ../..` → 16 claims, 0 errors, 16 book inputs.
- [ ] HIGH review; then integrate and tag `repo/book-evidence-registry-v1`.

## Open items

- Every later chapter number must be added as a claim; nothing may cite a number absent here.
- Cell status cross-checking against raw result JSON is deferred (recorded in RESULT.md).
