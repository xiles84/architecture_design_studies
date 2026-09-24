# Book evidence registry

Every number in the book must resolve to a claim here, and every claim must resolve to a run,
digest, report and signed analysis. The registry is validated by [`tools/evidence`](../../tools/evidence/README.md).

## One active source: `v3/`

| File | Purpose |
|---|---|
| `v3/claims.json` | **the active registry.** The still-active v2 claims carried forward, plus the six 2026-09-23 Study 05 claims and their successor gap claim; each claim names its family, study, strength, run-level status, structured support keys, confounds, limits and exact provenance |
| `v3/schema.json` | the committed JSON Schema the active registry must satisfy |
| `v3/confounds.json` | the six confounds that must be published beside the numbers they affect |
| `v3/coverage.json` | the coverage matrix (`measured` / `measured-but-confounded` / `planned` / `gap` / `not_applicable`) with review state |
| `v3/SUPERSESSION_LEDGER.md` | the attributed v2 → v3 change record and the clause-by-clause disposition of the retired `v2-gap-04` |
| `v3/ANALYSIS.md` | the signed ingest analysis naming, for each new claim, the run, digest and tag it resolves to |
| `v2/claims.json`, `v2/schema.json`, `v2/confounds.json`, `v2/coverage.json`, `v2/CORRECTION_LEDGER.md`, `v2/SYNTHESIS_HANDOFF.md` | **v2, frozen** at tag `repo/book-evidence-registry-v2`. Historical reference; v3 carries its still-active claims forward unchanged. |
| `claims.json`, `claims.schema.json` | **v1, frozen** at tag `repo/book-evidence-registry-v1`. Historical reference only; its claim-to-cell resolution failed (36 of 43 names resolve nowhere). Do not write book prose from it. |

Validation, v1 diagnostic and active-input listing:

```bash
cd tools/evidence
go run . validate    --repo ../..   # default: the v3 package, 0 errors required
go run . validate-v3 --repo ../..   # explicit v3 validation
go run . validate-v2 --repo ../..   # the frozen v2 predecessor (supersedes v1)
go run . v1-diagnostic --repo ../.. # the frozen v1 audit: 11/12 claims, 36/43 names
go run . validate --repo ../.. \
  --claims book/evidence/claims.json \
  --schema book/evidence/claims.schema.json   # explicit historical v1 inspection
```

`validate` dispatches on the file's own `schema_version`; the default resolves only v3, and v3
supersessions and retirements may name either frozen predecessor (v1 or v2).

## Adding a claim

A claim added to the active package may supersede a claim in either frozen predecessor (v1 or v2);
record the change in the version's ledger. Then:

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
