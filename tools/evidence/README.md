# tools/evidence

Validates the book claim-to-evidence registry (`book/evidence/claims.json`) against its
committed JSON Schema and the repository it cites. No external dependencies; runs offline.

```
go run . validate --repo ../..                 # schema + provenance + rules
go run . validate --repo ../.. --list-inputs   # claim ids that may appear in the book
go run . extract  --repo ../..                 # candidate claims from analysis front-matter
go test ./...
```

## What `validate` enforces

1. **Schema.** The registry must satisfy `book/evidence/claims.schema.json` (the validator
   evaluates the same JSON-Schema subset it uses: type, const, enum, pattern, lengths,
   required, additionalProperties and local `$ref`).
2. **Provenance resolves.** `report` and `analysis` exist; the environment page exists; numeric
   claims name a real `studies/<study>/results/<run_id>` directory; a `run_tag`, if given, exists.
3. **Digest is live.** `inputs_digest` must appear verbatim in both the report and the signed
   analysis. A re-run changes the digest, so a mismatch is a *stale digest*.
4. **Strength matches trials.** `repeated_controlled` needs ≥ 2 trials, `single_run_directional`
   exactly 1, `mechanism_supported` ≥ 1, `gap`/`analogy` none.
5. **Winners are valid cells.** A numeric claim's `winner_cells` must all appear in `cells` with
   status `valid`; a `failed`/`invalid` cell can never be a winner, and the declared
   `failed_cells` count must match the cells.
6. **Supersession is explicit.** If claim A supersedes B, B must be marked `superseded_by: A`;
   `--list-inputs` excludes anything superseded.
7. **Cross-study comparisons are gated.** `comparability: cross-study` requires a
   `comparability_record` path that exists.
8. **Transfers are labelled.** Analogy strength must carry `transfer: analogy`; a numeric claim
   cannot be labelled an analogy.

## Fixtures

`testdata/*.json` are negative fixtures, one violation each, asserted in
`internal/evidence/registry_test.go`: a failed cell as winner, a stale digest, a cross-study
ratio without a comparability record, a missing topology and an unmarked supersession.
