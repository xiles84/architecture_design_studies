# Book evidence registry

Every number in the book must resolve to a claim here, and every claim must resolve to a run,
digest, report and signed analysis. The registry is validated by [`tools/evidence`](../../tools/evidence/README.md).

## One active source: `v2/`

| File | Purpose |
|---|---|
| `v2/claims.json` | **the active registry.** Atomic claims; each names its family, study, strength, run-level status, structured support keys, confounds, limits and exact provenance |
| `v2/schema.json` | the committed JSON Schema the active registry must satisfy |
| `v2/confounds.json` | the six confounds that must be published beside the numbers they affect |
| `v2/coverage.json` | the coverage matrix (`measured` / `measured-but-confounded` / `planned` / `gap` / `not_applicable`) with review state |
| `v2/CORRECTION_LEDGER.md` | the attributed v1 → v2 change record, including the proof that v1 is unchanged |
| `v2/SYNTHESIS_HANDOFF.md` | claim → six-family-taxonomy mapping and the decision models for the book synthesis |
| `claims.json`, `claims.schema.json` | **v1, frozen** at tag `repo/book-evidence-registry-v1`. Historical reference only; its claim-to-cell resolution failed (36 of 43 names resolve nowhere). Do not write book prose from it. |

Validation, v1 diagnostic and active-input listing:

```bash
cd tools/evidence
go run . validate    --repo ../..   # default: the v2 correction package, 0 errors required
go run . v1-diagnostic --repo ../.. # the frozen v1 audit: 11/12 claims, 36/43 names
go run . validate --repo ../.. \
  --claims book/evidence/claims.json \
  --schema book/evidence/claims.schema.json   # explicit historical v1 inspection
```

`validate` dispatches on the file's own `schema_version`; the default resolves only v2.

## Adding a v2 claim

1. Write a statement a single status and a single winner can honestly qualify. If one clause needs
   a different strength or scope, split it into its own claim.
2. Give every numeric claim at least one **structured support key**
   (`topology` + `design` + `operation` + `metric`) whose `resolves_to` token appears **verbatim in
   the cited report**. This is the field whose absence made v1 unusable: a green run now means the
   claim resolved to a measured cell, not that two internal lists agree.
3. Keep `run_level` (the run's own cell totals, failures and control scope) separate from the
   claim's `support` cells. A selected subset's failures need not equal the whole run's.
4. Record `trials` as three dimensions — `whole_run_replications`, `fresh_load_trials_per_cell`,
   `inner_iterations` — and never raise `strength` above what they support.
5. Attach every applicable confound id from `confounds.json`; a registered confound nobody cites
   fails validation.
6. Run `go run . validate --repo ../..` — it must report 0 errors before the claim is used.
