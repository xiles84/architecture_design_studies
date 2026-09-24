# Result — v4 evidence package: retire the stale placement/endpoint gaps

**Task:** `task-20260924T102917Z-book-evidence-v4`
**Branch:** `book/evidence-v4`
**Required tag:** `repo/book-evidence-registry-v4`
**Capability:** HIGH session executing a LOW task (model `deepseek-flash`, Deep Code CLI; effort not exposed)
**Benchmark:** none — artifact-only over committed evidence; no database, no measurement

## What was delivered

1. **The active v4 registry** at `book/evidence/v4/claims.json` — 29 claims. All 27 still-active v3
   claims are carried forward content-identically; `v2-gap-05-colocation-unverified` and
   `v2-gap-07-real-network-balanced-endpoints` are retired through `retired_predecessor_claims` with
   the clause each retirement closes; `v4-gap-01-physical-placement-unverified` and
   `v4-gap-02-endpoint-distribution-partial` carry their still-true clauses and claim no measurement
   of their own. `v3-gap-01-study05-remaining-regimes` is marked `partially_superseded` by the two
   successors, with hard-TTL expiry, medium scale and open-loop demand named as what remains open.
2. **The disclosure set**: `v4/schema.json` (adds the lifecycle fields inside the validator's
   supported JSON-Schema subset), `v4/confounds.json`, `v4/coverage.json` (the two placement/endpoint
   dimensions reconciled with the retirement), `v4/SUPERSESSION_LEDGER.md` and the signed
   `v4/ANALYSIS.md`.
3. **A versioned validator**: `tools/evidence validate` defaults to v4; `validate-v4` is explicit;
   `ValidateV4` closes supersessions and retirements over v1 ∪ v2 ∪ v3 and adds the lifecycle rules —
   a claim in `claims[]` is active or partially superseded, a retired or superseded status may not
   stay in `claims[]`, a partially superseded claim must name its successors and what remains open, a
   gap must declare `gap_kind`, every successor id must resolve to an active claim, and a retired
   predecessor may not also be active. `validate-v3` and `validate-v2` keep their semantics.
4. **Six new tests** in `internal/evidence/v4_test.go`: the positive check, the retirement
   regression (both stale gaps retired, successors active and superseding them), the carried-forward
   guarantee (27 of 29 statements byte-identical to v3), and four lifecycle negative controls.
5. **One source of truth for the registry path**: the version is a single build input;
   `book/lib/evidence.typ` derives both the source-relative and the display path. `book/evidence/check.sh`
   (run by `book/build.sh` before compiling) fails on a fixed registry version, a retired or unknown
   claim id, a gap with no `gap_kind`, an asset drawing its own figure number, or an orphan figure;
   `check-waivers.txt` records the two known figure-number defects as owned debt.
6. **Book prose corrected by reference**: `concepts/distributed-placement.typ` (successor gaps, and
   the instrument gap separated from the partial coverage gap), `chapters/topologies.typ`,
   `chapters/major-families.typ`, `chapters/evidence-and-reproduction.typ` (no fixed path), and
   `lib/config.typ` (the dead hand-written `evidence-card` macro removed; cards render from the
   registry).
7. **Build hardening**: `book/build.sh` takes `--out`/`--manifest`, takes its own
   `ads-book-build-lock` volume (never `ads-run-lock`), records the raw and cleaned Typst version
   separately, and fails if the unverified-build banner reaches a release PDF. An input-less compile
   renders that banner instead of a provenance page of placeholders.
8. **`CONTEXT.md` and `LESSONS_LEARNED.md`** updated with the v4 state, the review verdicts, and six
   lessons from this task.

## Validation (from committed files)

```text
go vet ./...                          clean
go test ./...                         ok (v1, v2, v3 and the six new v4 fixtures; ~25 s)
validate (default, dispatch v4)       evidence: 29 claims, 0 errors, 29 book inputs (active source book/evidence/v4/claims.json)
validate-v4                           evidence-v4: 29 claims, 0 errors
validate-v3 (frozen)                  evidence-v3: 29 claims, 0 errors
check.sh                              OK with 2 waived figure-number defects (owned by book/figures-v2)
book/build.sh --out dist/preview.pdf  53 pages, 29/29 claims indexed, 5/5 fonts embedded,
                                      evidence digest present, unverified banner absent — verified
v3 unchanged vs tag                   git diff repo/book-evidence-registry-v3 -- book/evidence/v3/ is empty
carried-forward claims                normalising the four lifecycle keys and diffing v3 vs v4 reports no change
```

### Negative controls (introduced, observed, reverted)

| Control | Expected | Observed |
|---|---|---|
| A retired claim id quoted in a source | source checker fails | `[FAIL] retired-claim-quoted: v2-gap-05-colocation-unverified is retired in v4; cite its successor instead` |
| The same id through the Typst loader | compile aborts | `panicked with: claim v2-gap-05-colocation-unverified is retired in book/evidence/v4/claims.json; cite its successor instead: v4-gap-01-physical-placement-unverified` |
| An unknown claim id | source checker fails | `[FAIL] unknown-claim-id: v9-does-not-exist is not a claim in v4` |
| A fixed registry version in a source | source checker fails | `[FAIL] registry-literal: … names a fixed registry version` |
| A compile without the build inputs | banner present | banner text present in the preview PDF (1), absent in the release build (0) |

## The defect this closes

| Retired clause | Why it was false | Successor |
|---|---|---|
| `v2-gap-05` "Study 05's placement pair was never executed" | `v3-06` records the key-locality pair executed on 2026-09-23 | closed; recorded in `closed_dimensions` |
| `v2-gap-05` "no study verifies actual node placement" | still true — no tablet/leader distribution was produced | `v4-gap-01` |
| `v2-gap-07` "balanced client access across cluster endpoints is not proven" | `v3-05`'s client spreads over all three endpoints | closed for the Study 05 cells; recorded |
| `v2-gap-07` "endpoint distribution is not measured" | still true for the studies 01–02 multi-node cells, whose `yb-cluster3` DSN names `yb-n1` only | `v4-gap-02` |
| `v2-gap-07` "no real network" | still true everywhere | `v4-gap-02` |

v3's own `coverage.json` had already been updated for the endpoint spread while `claims.json` still
said the opposite, and the book printed `claims.json`. That is the sharper form of finding R04: not
missing evidence, but two files in one package disagreeing.

## Conclusion

The active registry no longer asserts gaps its own evidence closes, and the lifecycle fields make the
next such omission visible to the validator instead of to a reader. No numeric claim, report, result
file or signed analysis was edited; v1, v2 and v3 are byte-unchanged and no tag moved. The tracked
PDF in `book/dist/` is unchanged and will be regenerated once at release from fully merged sources.
