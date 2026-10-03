# SELECT metadata source check

**TL;DR:** both initial corrected single-node engine checks passed. Final source
adds exact raw-result SHA-256 to generated reports, actual client CPU/memory limits,
missing-plan rejection and alias-normalized physical structure signatures.

GPT-6 via Codex desktop, HIGH, effort unknown; named environment
`xiles-podman-select-validation-20261003`. Podman Go 1.26-bookworm:
`go vet ./... && go test ./... && go build ./...` passed for both Study 01 packages.
Runner syntax and diff checks passed. Added regressions show that alias spelling
does not change structure and missing EXPLAIN cannot admit timing. The late-answer
regression discards earlier samples when the final sample fails independent truth.

Initial checks (unique tags/source identities retained in their manifests):
`select-pg-check-20261003` and `select-yb-check-20261003` each passed 17 templates,
53 live fixture cells and four wrong SQL controls, with 51 drained client samples.
Their factual results remain immutable. Final live checks follow from the source
commit introducing this metadata checkpoint; no comparative matrix is claimed.

NEXT MODEL: HIGH
