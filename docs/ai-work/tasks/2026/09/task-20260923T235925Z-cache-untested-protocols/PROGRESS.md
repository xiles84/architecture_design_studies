# Cache protocol task progress

Owner declared HIGH and twice requested iteration through all HIGH tasks.
Session: `session-20261003-cache-protocols`, GPT-6 via Codex desktop; effort
unknown. Queue claim: `claim-f22601f822ed6d18`, epoch 1. Worktree:
`.worktrees/study05-v5-untested-protocols`; branch `study-05/v5-untested-protocols`.

## Decision log

- **Release the conditional comparative experiment.** The renewed execution
  request covers the already specified unanswered publication and deletion
  questions. Confidence: High. Alternative: leave it proposed awaiting a
  repeated request for this same comparison. Affects this task only; resource
  deferrals, T1/T2-before-T5 and Edition 3 approval remain intact.
- **Use a bounded keyed projection.** Isolate cache protocol mechanics with
  identical SQL and resources, preserving the strict acknowledgement contract
  and an independent ledger. Confidence: High. Alternative: change the existing
  donor-portal harness. Affects the new evidence's scope: no donor-portal or
  cross-run throughput claim follows.
- **Retain source identity across physical deletion.** The three deletion
  packages share durable generations and synchronous fencing. Confidence:
  High. Alternative: erase all identity metadata; that is the explicit ABA
  negative control, not a correct candidate.
- **Keep per-key client writer ordering.** This matches the study's default
  ledger contract and isolates publication costs. Confidence: High. Alternative:
  redesign the oracle for uncoordinated same-key writers. Affects contention
  interpretation; report client wait and do not claim database hot-key scaling.
- **Measure within the unchanged machine configuration.** Equal-total service
  budget 2 CPU/768 MiB; client 1 CPU/512 MiB; configured reserve 1 CPU/768 MiB.
  Confidence: High. Alternative: increase VM resources. Cluster/placement gaps
  remain explicit and do not become completed comparisons.

## Checkpoints

- Premeasurement protocol committed `d5ad17f3`, annotated milestone
  `study-05/v5-protocol-20261003` before implementation.
- New protocol implementation and HIGH self-review are in progress; no runtime
  measurement has started. Existing harnesses, signed analyses and run tags are
  untouched.

NEXT MODEL: HIGH
