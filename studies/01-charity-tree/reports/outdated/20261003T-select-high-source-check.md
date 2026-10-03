# HIGH SELECT harness correction check

GPT-6 via Codex desktop; source baseline `9a5969951925ad9e844f5f47e2e007e9bdbb216b`,
integrated current main `460b66ab151bf970ac3f1de255a0675555eeb67f` in task worktree.
Environment: `xiles-podman-select-validation-20261003`.

**TL;DR:** baseline unit checks passed but left real semantic/provenance defects.
The corrective Go regressions and complete Study 01 vet/test/build checks pass
in `docker.io/library/golang:1.26-bookworm`. Live validation is recorded below.
This is a harness check, with no performance-winner or prediction verdict.

## Source checks

- Read-only baseline: `go vet ./selectbench/...`, 40 named package tests and
  `go build ./...` passed. This did not validate SQL control execution or ordering.
- Correction: `go vet ./... && go test ./... && go build ./...` passed for the
  historical harness and additive SELECT package.
- New meaningful regressions reject wrong pagination/top-N arrival order,
  fractional/non-finite/trailing-junk integer values, controls with correct
  answers or wrong shapes, and fabricated zero timing fields. Keyset parameters
  encode ascending-id tie continuation and NULL cursor state. Real drained samples
  and unexpected error classification are checked.
- `bash -n run-select.sh` and `git diff --check` passed.
- Legacy live results `select-20261003T181045Z` and `select-20261003T181716Z`
  are retained as worker checks. Their missing image/environment/run-tag fields,
  zero timing placeholders and unverified edge/control coverage prevent using
  them as performance evidence. The earlier worker's report is preserved.

## Live checks

Pending committed/tagged corrected source, with one database at a time.

NEXT MODEL: HIGH
