# Book evidence registry

Every number in the book must resolve to a claim here, and every claim must resolve to a run,
digest, report and signed analysis. The registry is validated by [`tools/evidence`](../../tools/evidence/README.md),
and the Typst sources are checked by [`check.sh`](check.sh) before every build.

## One active source: `v6/`

| File | Purpose |
|---|---|
| `v6/claims.json` | **Active claim registry.** Contains current claims and their scientific method metadata; historical migration provenance is preserved. |
| `v6/schema.json` | the committed JSON Schema the active registry must satisfy, including the method enums |
| `v6/confounds.json` | the six confounds that must be published beside the numbers they affect |
| `v6/coverage.json` | Coverage dimensions and `active_claim_summary`, checked against the active registry. |
| `v6/method_migration.json` | the per-claim record of how each v6 field was assigned, so the classification is auditable rather than inferred |
| `v6/SUPERSESSION_LEDGER.md` | the attributed v5 → v6 change record (frozen v4 → v5 text below it) |
| `v6/ANALYSIS.md` | Historical signed analysis of the initial v6 method-metadata pass. Its original counts are a dated snapshot, not the live registry totals. |
| `v5/*` | **frozen** predecessor: the five claims re-scoped by the 2026-09-24 claim-versus-limit audit (v2-07, v2-10, v2-14, v2-16, v3-06) |
| `v4/*`, `v3/*`, `v2/*`, `claims.json`, `claims.schema.json` | **frozen** historical packages at tags `repo/book-evidence-registry-v4`, `-v3`, `-v2`, `-v1`. Never edited; a later package carries their still-active claims forward and records what it retires. |
| `check.sh`, `check-waivers.txt` | the structural source rules the build runs (including `concept-checklist` against `book/concepts/CHECKLIST.md`), and the standing exceptions to them |

## Claim lifecycle (v4 and later)

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
go run . validate    --repo ../..   # default: the v6 package, 0 errors required
go run . validate-v6 --repo ../..   # explicit v6 validation (lifecycle + repetition + method lints)
go run . validate-v5 --repo ../..   # explicit v5 validation (frozen predecessor)
go run . validate-v4 --repo ../..   # explicit v4 validation (frozen predecessor, unchanged)
go run . validate-v3 --repo ../..   # the frozen v3 predecessor
go run . validate-v2 --repo ../..   # the frozen v2 predecessor
go run . v1-diagnostic --repo ../.. # the frozen v1 audit: 11/12 claims, 36/43 names
go test -count=1 ./...              # the resolver and lifecycle rules, including negative controls
```

**Always pass `-count=1` when the registry may have changed.** These tests read the live registry from
disk at runtime, and Go's test cache only tracks the package's build inputs — so after a claim is edited,
`go test ./...` can report a cached `ok` from before the edit. A cache hit is not a passing gate. The same
applies to any wrapper that pipes the result: a pipeline's exit status is the *last* command's, so
`go test ./... | tail -1` will happily print `FAIL` and return success. The scripts in this repository set
`set -o pipefail` for that reason; ad-hoc command chains do not.

`validate` dispatches on the file's own `schema_version`; the default resolves v6, and a v6
supersession or retirement may name a claim in any frozen predecessor (v1, v2, v3, v4 or v5).

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
7. Run `go run . validate --repo ../..` and `go test -count=1 ./...` — both must be green before the
   claim is used, and the tests must not be served from the cache (see the note above the command list).

## Current coverage statistics

The active totals are derived from `v6/claims.json` and stored in `v6/coverage.json` under `active_claim_summary`. The build rejects stale totals through `check-active-coverage.sh`. These counts do not establish independent review or preregistration.

The signed `v6/ANALYSIS.md` preserves the original 2026-09-25 conclusions and must not be silently updated.
