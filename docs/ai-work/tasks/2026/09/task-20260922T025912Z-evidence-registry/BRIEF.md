# Execution brief — evidence registry

Implement a Go validator/extractor and `book/evidence/claims.json`. Extraction may copy
facts and provenance from current reports/analyses but may not invent or strengthen an
interpretation. Use exactly these evidence levels: `repeated_controlled`,
`single_run_directional`, `mechanism_supported`, `analogy`, `gap`, `invalid`.

Every numeric claim needs its controlled pair, valid cells, run, digest, producing tag,
environment, topology, report, and signed analysis. Record scope and limitations. Reject
cross-study numeric rankings without an explicit comparability record. Reject a failed or
invalid cell used as winning performance evidence. Record supersession and stale digests.

Seed the registry across all five studies, including gap records for native datastore
families, real-network topology, and incomplete Studies 04/05 coverage. HIGH reviewers,
not the extractor, decide whether provisional claims are acceptable.

Run entirely in Podman, commit tests and fixtures, and submit for HIGH review.

NEXT MODEL: LOW
