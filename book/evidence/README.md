# Book evidence registry

Every number in the book must resolve to a claim here, and every claim must resolve to a run,
digest, report and signed analysis. The registry is validated by [`tools/evidence`](../../tools/evidence/README.md),
and the Typst sources are checked by [`check.sh`](check.sh) before every build.

## One active source: `v4/`

| File | Purpose |
|---|---|
| `v4/claims.json` | **the active registry.** The still-active v3 claims carried forward unchanged, plus the two successor placement/endpoint gaps and the v4 claim lifecycle fields; each claim names its family, study, strength, status, run-level status, structured support keys, confounds, limits and exact provenance |
| `v4/schema.json` | the committed JSON Schema the active registry must satisfy |
| `v4/confounds.json` | the six confounds that must be published beside the numbers they affect |
| `v4/coverage.json` | the coverage matrix (`measured` / `measured-but-confounded` / `planned` / `gap` / `not_applicable`) with review state |
| `v4/SUPERSESSION_LEDGER.md` | the attributed v3 → v4 change record and the clause-by-clause disposition of the two retired gaps |
| `v4/ANALYSIS.md` | the signed ingest analysis naming, for each successor claim, the run, digest and tag it resolves to |
| `v3/*`, `v2/*`, `claims.json`, `claims.schema.json` | **frozen** historical packages at tags `repo/book-evidence-registry-v3`, `-v2`, `-v1`. Never edited; a later package carries their still-active claims forward and records what it retires. |
| `check.sh`, `check-waivers.txt` | the structural source rules the build runs, and the standing exceptions to them |

## Claim lifecycle (v4)

`claims[]` holds only claims that may appear as evidence. Retiring one means moving it out:

- `status` — `active`, `partially_superseded`, or (illegal in `claims[]`) `superseded` / `retired`.
- `superseded_by` — the successor claim ids. Required for `partially_superseded`.
- `remaining_dimensions` — what is still open after a partial supersession.
- `closed_dimensions` — on a retired entry, the specific clauses the new evidence closed. This is the
  field that would have caught the Edition 1 defect: v3 carried `v2-gap-05` and `v2-gap-07` forward
  unchanged while adding the v3-05/v3-06 runs that closed parts of both, and the book printed them
  side by side.
- `gap_kind` — `coverage`, `schema_limitation`, `instrument` or `unstable_measurement`. A gap must
  say which; the book badges them differently, because "no representation can answer this" is not
  the same finding as "nobody ran it".

## Validation

```bash
cd tools/evidence
go run . validate    --repo ../..   # default: the v4 package, 0 errors required
go run . validate-v4 --repo ../..   # explicit v4 validation (v4 lifecycle rules included)
go run . validate-v3 --repo ../..   # the frozen v3 predecessor
go run . validate-v2 --repo ../..   # the frozen v2 predecessor
go run . v1-diagnostic --repo ../.. # the frozen v1 audit: 11/12 claims, 36/43 names
go test ./...                       # the resolver and lifecycle rules, including negative controls
```

`validate` dispatches on the file's own `schema_version`; the default resolves v4, and a v4
supersession or retirement may name a claim in any frozen predecessor (v1, v2 or v3).

Source-level checks, run by `book/build.sh` before it compiles:

```bash
book/evidence/check.sh
```

It fails the build when a Typst source names a fixed registry version, when a chapter renders a
retired or unknown claim id, when a gap has no `gap_kind`, when an asset draws its own figure
number, or when a registered figure is never embedded.

## Adding a claim

A claim added to the active package may supersede a claim in any frozen predecessor; record the
change in the version's ledger, or — when the change is a retirement — in the ledger plus
`retired_predecessor_claims` with its `closed_dimensions`. Then:

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
6. **Re-read the gaps the new evidence touches, not only the claim you are adding.** If a new run
   closes a clause of an existing gap, retire or narrow that gap in the same change; that omission
   is what the v4 lifecycle fields exist to expose.
7. Run `go run . validate --repo ../..` and `go test ./...` — both must be green before the claim
   is used.
