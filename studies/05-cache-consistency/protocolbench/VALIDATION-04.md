# Final matrix validation

2026-10-03, GPT-6 via Codex desktop. The producing image passed vet and three
oracle tests under `go test -race ./protocolbench`. Source commit: `fed61266`;
client image id: `54d66b9c85d21b0d562c076f2c8cf2ddd7ce74755979c09461db416c926ed2f0`.

| Whole run | Input digest | Cells | Control gates | Admitted trials |
|---|---|---:|---:|---:|
| cache-v5-20261003-r1 | f2054e7244c3d2b9 | 10 | 57 | 72 |
| cache-v5-20261003-r2 | 52c1eb9b8c355d19 | 10 | 59 | 72 |
| cache-v5-20261003-r3 | 76c170e2f1a3e27e | 10 | 57 | 72 |

All three exact-image verification commands passed complete control gates,
phase/trial counts, warmed read-only source-query contrasts and zero candidate
wrong/impossible reads. Each run has fresh database/cache state, randomized
order and a same-candidate first/last position control. Runtime resource/image
records, source plans, faults and generated factual reports are retained.

The exact producing image generated the whole-run reduction successfully.
Large first-run position drift and substantial client throttling are explicit
analysis limits; no general throughput winner is assigned. The 500-ms control
run remains separate. No prior study result or signed artifact was edited.

The runner removed its own containers/data volume and released `ads-run-lock`.
This passes the bounded task's measurement gates, with the predeclared cluster,
placement, fixed-payload and contention limits. HIGH analysis/manifests review,
queue integration, final ancestry and milestone checks are next.

NEXT MODEL: HIGH
