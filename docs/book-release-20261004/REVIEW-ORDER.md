# HIGH review — SELECT order correction

Reviewer: GPT-6 via Codex desktop, HIGH, selected effort unknown; 2026-10-04.
Reviewed LOW's Decision Log first, then runner, 18 result JSONs, controls,
producing tags, manifests and signed analysis. This is an independent review of
deepseek-flash's execution, plus disclosed self-review of the small exporter fix.

**Accept for integration with the qualifications in the signed companion.**
The nine logged choices are acceptable: deterministic per-run seeds, interleaving,
actual-first position control, fanout-local fixture boundaries, separate physical
index arm, unchanged contracts/matrix, live resources, exact input digests and
uncertainty qualification. Every timed sample is checked against the oracle and
a late mismatch discards earlier samples; a new regression test verifies this
and the executed round sequence. All 18 producing tags remain unchanged.

The exporter searched for `run_tag` but the actual SELECT provenance field is
`RunTag`, leaving blank tag cells in the original generated manifests. Fixed the
exporter with a compatibility test. The independently generated
`CORRECTION-INPUTS.md` supplies all exact input hashes, tags and family digests;
the original manifests remain as historical artifacts.

The signed LOW analysis has conflicting aggregate claims: YB level 10 is mixed,
and level 10,000 has a near-equality trial. The new signed companion corrects
these interpretations without editing the original. No SELECT timing claim is
promoted into the book in this release. The aggregate family remains qualified
by 29–74% position drift; three same-direction medians alone are insufficient
to establish a practical advantage. Cluster/throughput and join/CTE gaps remain.

Containerized `go vet ./...`, `go test ./...` and static benchmark build passed
in `localhost/charitybench-releasecheck:20261004` (image
`63b5e4ca3d1c6249f8ac8d304120e4732e478ccb9bc7354051b63142f91fb372`).
The input check runs in `docker.io/library/golang:1.26-bookworm`; no database was
started for this review.

TASK COMPLETE
