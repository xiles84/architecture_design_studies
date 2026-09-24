# tools/evidence

Validates the book claim-to-evidence registry against its committed JSON Schema and the repository
it cites. No external dependencies; runs offline.

The **active** registry is the v3 package, `book/evidence/v3/claims.json`. `validate` defaults to it
and dispatches on the file's own `schema_version`; point `--claims`/`--schema` at the v2 or v1 files
to inspect an older registry explicitly.

```
go run . validate    --repo ../..                # active v3: schema + semantic resolver + rules
go run . validate    --repo ../.. --list-inputs  # v3 claim ids that may appear in the book
go run . validate-v3 --repo ../..                # explicit v3 validation
go run . validate-v2 --repo ../..                # explicit v2 validation (predecessor set = v1)
go run . v1-diagnostic --repo ../..              # the frozen v1 resolution audit (11/12, 36/43)
go run . validate --repo ../.. \
  --claims book/evidence/v2/claims.json --schema book/evidence/v2/schema.json \
  --confounds book/evidence/v2/confounds.json    # historical v2 inspection
go run . validate --repo ../.. \
  --claims book/evidence/claims.json --schema book/evidence/claims.schema.json   # historical v1
go run . extract     --repo ../..                # candidate claims from analysis front-matter
go test ./...
```

## What versioned `validate` enforces

1. **Schema.** `book/evidence/v3/schema.json` (the same JSON-Schema subset: type, const, enum,
   pattern, lengths, required, additionalProperties and local `$ref`).
2. **Every structured support key resolves.** For each `support[]` entry, `resolves_to` must appear
   **verbatim in the primary anchor's cited report**. This is the semantic resolver that v1 lacked:
   a green run now means the claim resolved to a measured cell, not that two internal lists agree.
3. **Provenance resolves.** Every anchor's report and analysis exist; `inputs_digest` appears in
   both; a numeric claim names a real `studies/<study>/results/<run_id>` directory; a `run_tag`, if
   given, exists and is a `run/<study>/...` tag.
4. **Run level is separate from support.** `run_level` carries the run's own cell totals, failures
   and control scope; `reported_failed` may not exceed `reported_cells`, and `controls_fired` may
   not exceed `controls_present`.
5. **Strength matches trials.** `repeated_controlled` needs ≥ 2 whole-run replications or
   fresh-load trials; `single_run_directional` may not declare more than one whole-run replication.
6. **Legacy provenance is explicit.** `legacy_provenance_incomplete` is allowed only when the
   primary run has no tag — and required when it has none. Tags and commits are never invented.
7. **Confounds are closed.** Every confound a claim names exists in `confounds.json` and lists that
   claim in `affected_claims`; every registered confound is referenced by at least one claim.
8. **Supersession and retirement are explicit.** `supersedes` must name a claim that exists in an
   allowed predecessor registry, and a retirement must too: v2 resolves against frozen v1, v3 against
   frozen v1 **and** v2. Tags and commits are never invented.
9. **Gaps carry a basis.** A `gap` claim must name the artefact that establishes the absence.

## Fixtures

`testdata/*.json` are the v1 negative fixtures, asserted in `internal/evidence/registry_test.go`.
`internal/evidence/v2_test.go` adds the v2 fixtures: the committed registry validates; an
unresolvable support key, a mismatched legacy flag, an unreferenced confound and an unknown
supersession each fail; and `TestV1DiagnosticRegression` pins the frozen v1 audit at 12/11/43/36 so
the exact failure mode the correction fixed stays reproducible. `internal/evidence/v3_test.go` adds
the v3 positive check (the committed package validates), a supersession-closure check (a v3 claim may
supersede a v2 claim) and an unresolved-support-key negative control.
