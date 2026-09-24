# Result — version the Study 05 v2 evidence (v3 package)

**Task:** `task-20260923T235907Z-book-evidence-study05-v3`
**Branch:** `repo/book-evidence-study05-v3`
**Required tag:** `repo/book-evidence-registry-v3`
**Capability:** HIGH session executing a LOW task (model `deepseek-flash`)
**Benchmark:** none — artifact-only over committed evidence

## What was delivered

1. **The active v3 registry** at `book/evidence/v3/claims.json` — 29 claims: the 22 still-active v2
   claims carried forward byte-for-byte in content, six new numeric Study 05 claims (`v3-01`…`v3-06`)
   and the successor gap claim `v3-gap-01`. The retired `v2-gap-04-study05-unmeasured-regimes` is not
   reproduced.
2. **The supersession ledger** `book/evidence/v3/SUPERSESSION_LEDGER.md`: the frozen-predecessor
   hashes, the carried-forward list, the clause-by-clause disposition of `v2-gap-04`, the
   claim → report cell/digest/tag table, and the book-wiring record.
3. **The disclosure set**: `v3/schema.json`, `v3/confounds.json` (conf-04 no longer lists the retired
   id), `v3/coverage.json` (cadence/churn, equal-total, repeated trials and YugabyteDB coverage
   updated; placement/native/real-network/open-loop gaps retained), and the signed ingest analysis
   `v3/ANALYSIS.md`.
4. **A versioned validator**: `tools/evidence validate` now defaults to v3; `validate-v3` is explicit;
   `ValidateV3` closes supersessions and retirements over v1 ∪ v2 while `ValidateV2`/`v1-diagnostic`
   keep their old semantics. `internal/evidence/v3_test.go` adds positive, supersession and
   unresolved-support-key fixtures.
5. **Book sources corrected by reference**: `book/chapters/scenarios.typ` (gap paragraph + cards),
   `book/README.md`, `book/lib/evidence.typ`, `book/build.sh`, `book/main.typ`, and `CONTEXT.md`.

## Validation (from committed files)

```text
go vet ./...                    clean
go test ./...                   ok (v1 fixtures, v2 fixtures + regression, v3 fixtures)
validate (default v3)           evidence: 29 claims, 0 errors, 29 book inputs (active source book/evidence/v3/claims.json)
validate-v2 (frozen v2)         evidence-v2: 23 claims, 0 errors, active source book/evidence/v2/claims.json
v1-diagnostic                   12 numeric claims, 11 with a missing name, 36/43 names unresolved
v2 unchanged vs tag             git diff repo/book-evidence-registry-v2 -- book/evidence/v2/ is empty
carried-forward claims          jq deep-compare v2 vs v3 (except v2-gap-04) returns []
book compile (pinned image)     typst compile --root book main.typ (throwaway output) exits 0
```

## New claims (each resolves to a report cell, digest and run tag)

| Claim | Run | Digest | Tag |
|---|---|---|---|
| `v3-01-churn-crosses-real-ttl` | `20260923T1100Z-churn1800-through` | `8ad59fb36a65c17f` | `run/05-cache-consistency/20260923T1100Z-churn1800-through` |
| `v3-02-churn-600-through-vs-aside` | `20260923T1020Z-churn600-through` (+ aside) | `75af94d4c846e528` (+ `cf16e598ad9b32f8`) | `run/05-cache-consistency/20260923T1020Z-churn600-through` |
| `v3-03-equal-total-framing-labelled` | `20260923T1150Z-res2-equaltotal-pg` (+ db-only) | `e0dfc55f3c363944` (+ `b28d8dd331c4c968`) | `run/05-cache-consistency/20260923T1150Z-res2-equaltotal-pg` |
| `v3-04-yb-single-strict-relaxed` | `20260923T1152Z-res2-yb1` | `6f04a8faf4a428bf` | `run/05-cache-consistency/20260923T1152Z-res2-yb1` |
| `v3-05-yb-cluster3-strict-relaxed` | `20260923T1157Z-res2-yb3` | `f2133069805cf82d` | `run/05-cache-consistency/20260923T1157Z-res2-yb3` |
| `v3-06-colocated-vs-noncolocated-locality` | `20260923T1215Z-res2-placement-yb3` | `a18817d592f767c5` | `run/05-cache-consistency/20260923T1215Z-res2-placement-yb3` |
| `v3-gap-01-study05-remaining-regimes` | `20260923T1100Z-churn1800-through` | `8ad59fb36a65c17f` | `run/05-cache-consistency/20260923T1100Z-churn1800-through` |

## Supersession v2 → v3 (summary)

- `v2-gap-04` real-TTL churn → **superseded** by `v3-01`/`v3-02` (2.0 and 6.0 boundaries, 0 wrong
  reads; hard expiries 0).
- `v2-gap-04` cache resource accounting → **superseded** by `v3-03` (labelled equal-total arm).
- `v2-gap-04` YugabyteDB/cluster → **superseded** by `v3-04`/`v3-05`.
- `v2-gap-04` verified placement → **partly superseded** by `v3-06`; the missing tablet/leader
  distribution is retained in `v3-gap-01`.
- `v2-gap-04` hard TTL expiry, medium scale, open-loop demand → **retained** in `v3-gap-01`.

## Conclusion

The book's active numeric source is now the v3 package; v1 and v2 stay frozen and readable, and no
earlier analysis, report or digest was edited. No measurement was performed or authorized.
