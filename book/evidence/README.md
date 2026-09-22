# Book evidence registry

Every number in the book must resolve to a claim here, and every claim must resolve to a run,
digest, report and signed analysis. The registry is validated by [`tools/evidence`](../../tools/evidence/README.md).

| File | Purpose |
|---|---|
| `claims.json` | the claims; each one names its study, strength, limits and exact provenance |
| `claims.schema.json` | the committed JSON Schema the registry must satisfy |

## Adding a claim

1. Take the statement **verbatim** from a signed analysis headline (or an explicit sentence) —
   the extractor must not introduce interpretation.
2. Fill `provenance` from the analysis and report front-matter: `run_id`, `environment`,
   `topology`, `inputs_digest`, `repo_commit`, `report`, `analysis`, and `run_tag` when the run
   was tagged.
3. List the measured `cells` with their status and name the `winner_cells`. A failed cell can
   never be a winner.
4. Set `strength` and `trials` honestly: `repeated_controlled` only with ≥ 2 trials,
   `single_run_directional` for one, `mechanism_supported` when a mechanism is argued from at
   least one run, `gap` for a missing measurement, `analogy` for a labelled transfer.
5. Label scenario transfers `direct` or `analogy`; a cross-study numeric comparison needs a
   `comparability_record` path.
6. Run `go run ./tools/evidence validate --repo .` — it must report 0 errors.
