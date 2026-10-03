# HIGH SELECT harness correction check

**TL;DR:** final single-node checks on PostgreSQL and YugabyteDB each passed
17 templates, 53 live fixture cells and all four deliberately wrong SQL controls.
Every candidate recorded three verified drained-query samples. Exact saved-JSON
SHA-256 matches each generated report. This validates the harness; it establishes
no replicated performance winner, crossover or prediction verdict.

GPT-6 via Codex desktop, HIGH, selected effort unknown. Environment:
`xiles-podman-select-validation-20261003`. Resumed DeepSeek's canonical handed-off
task through epoch-4 CAS recovery after the ten-minute grace; prior worker
progress, results and reports were preserved. Supersedes the source-only checkpoint
now at `reports/outdated/20261003T-select-high-source-check.md`.

| Final check | PostgreSQL | YugabyteDB |
|---|---|---|
| Run id | select-pg-check-20261003-final | select-yb-check-20261003-final |
| Source commit | b4583c11 (full hash in result/manifest) | 4f54d848 (full hash in result/manifest) |
| Run tag | run/study-01/select-pg-check-20261003-final | run/study-01/select-yb-check-20261003-final |
| Candidate cells | 17 valid; 0 non-valid | 17 valid; 0 non-valid |
| Untimed live fixture cells | 53 valid | 53 valid |
| Wrong SQL controls | 4/4 detected; no query/decoding errors | 4/4 detected; no query/decoding errors |
| Verified client samples | 51 | 51 |
| Input SHA-256 | 3180cf4f671d5047281a5ab37373e14faf9da6eee3578d78d2d454fcf115ba22 | 6b2eabd76ca00c115ea0e6c148fb3e21e8021af9e86d8874361fff2582c7913d |
| Client cpu.max | 200000 100000 = 2 CPU | 200000 100000 = 2 CPU |
| Client memory.max | 2147483648 = 2 GiB | 2147483648 = 2 GiB |

Both final checks used client image
`18af3b106a1fcd39996e12ad5596eda4dd822c6c830995e63420dc879d9dfc83`.
Database images are pinned to PostgreSQL 17.11 and YugabyteDB 2025.2.6.0-b111;
their actual image ids and runtime banners are in each result/manifest. The runner
enforced 2 CPU/3 GiB per single database, 2 CPU/2 GiB client, one sequential reader
and a two-slot pool. It retained node/machine/image inspection, guest CPU/memory,
actual client cgroup limits, whole-check client/node CPU counters and raw logs.
Guest CPU observation: 0-7. Quotas are not core pins. Both runs finished with
exit code 0, their databases stopped and the benchmark lock released.

## What the checks establish

- Independent Go truth precedes EXPLAIN/client timing and checks every sample.
  A late mismatch discards earlier samples. Missing plans or engine errors admit
  no client timing. Unexpected errors are not called feature gaps or timeouts.
- Live teaching, empty, all-miss, keyset tie and NULL-cursor cases passed on both
  engines. Top-N emits no unmatched-parent row. Ordered contracts retain database
  row arrival order; verification cannot repair a wrong ORDER BY.
- All four wrong SQL controls executed successfully with matching contract column
  shape, then differed from independent truth. A decoding/query error cannot count
  as detection. No control is timed/ranked as a faster template.
- Samples measure client query/row-drain time. Server planning/execution times and
  complete EXPLAIN counters remain separate. Decoded-row JSON size is labelled as
  such, rather than called wire bytes. Missing fields render N/A.
- Structure signatures retain hierarchy/outer-inner roles and normalize aliases;
  plan fingerprints alone do not establish semantic equivalence. Reports retain
  raw plans and refuse a performance-winner claim from one whole run.
- `go vet ./... && go test ./... && go build ./...` passed in the pinned Go Podman
  image for both the historical harness and additive package. The image build
  repeats this gate. The new regressions cover wrong ordering, tie/NULL binding,
  lossy integers, false control detections, late mismatches, absent plans,
  alias-only renaming, real sample recording and N/A output. Runner syntax and
  source/docs diff checks passed; raw YB terminal spinners are retained unchanged
  as opaque logs, with byte conversion/diff disabled only for these new logs.

## Retained history and limits

Initial corrected checks `select-pg-check-20261003` and
`select-yb-check-20261003` also passed, using the earlier recorded source/image;
they are retained separately. Legacy worker checks `select-20261003T181045Z` and
`select-20261003T181716Z` remain immutable and are not performance evidence: their
missing image/environment/run-tag fields, zero timing placeholders and absent live
edge/control coverage prevent that use. The earlier worker report is preserved.

An accidentally replaced `.gitattributes` digest policy was restored in the next
commit; final policy retains every existing rule and adds only new SELECT paths.
No stored result bytes were changed. No source/report tag was moved.

Only bounded single-node validation was performed. Three-node/equal-total,
actual placement/balanced routing, mutation/storage costs, workload concurrency,
demand calibration, repeated whole-run position controls and crossover sweeps
remain measurement-task obligations. YB RPC metrics remain in raw plans where
provided, with parsed RPC coverage explicitly unmeasured. No numeric book claim or
signed performance analysis is added by this harness task.

NEXT MODEL: HIGH
