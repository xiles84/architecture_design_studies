# Corrective validation — reporting and evidence fixes (LOW)

Task `task-20261002T131729Z-select-join-cte`, attempt
`attempt-20261004T173630Z-d9091f`, LOW `deepseek-flash` / opencode CLI.
Responds to R3/R4 and the mandatory-rendering items of the
[CORRECTIVE-HANDOFF](../20261004-join-high/CORRECTIVE-HANDOFF.md). No database was
started for these checks; they are source, unit-test and file-digest checks.

## Fixed defects and their evidence

| # | Defect (review) | Fix | Evidence |
|---|---|---|---|
| 1 | `runSelectJoinReport` hashed `append(data, '\n')`, advertising a digest the run never had (R4) | one shared `joinInputsDigest` over the exact file bytes | corrected re-render of the original pg1 input advertises `070df793…` (the review's broken re-render showed `7c879564…`); see [rerender-pg1-fixed.md](rerender-pg1-fixed.md) |
| 2 | client image labelled `EngineImage` beside a database digest (R4) | `Provenance` gains `ClientImage`/`ClientImageDigest`; runner passes `ADS_ENGINE_IMAGE`; report prints both | `TestJoinProvenanceSeparatesClientAndEngine` |
| 3 | replicated report keyed only on engine, topology, axis (R4) | keyed on a compatibility class (environment, dataset, client/engine image digests, settings) plus engine/topology/axis; classes printed | `TestJoinReplicatedPartitionsByClass` |
| 4 | `EXPLAIN (ANALYZE, BUFFERS, FORMAT TEXT)` for both engines; RPC counters never requested (R3) | `probeExplainOptions` picks PG `ANALYZE, VERBOSE, BUFFERS, FORMAT TEXT` or YB `ANALYZE, DIST, DEBUG, FORMAT TEXT`, records the choice and any fallback | recorded per run in `explain_options` and the notes; the probe runs under the lock |
| 5 | token volatility asserted only for a constant expression (R3) | feasibility records distinct-id token evidence and that the reuse plan carries the expression | recorded in `feasibility.token_evidence` |
| 6 | no throughput/latency summary; capacity implied | min/median/max samples and derived `queries/s` added, labelled a derived rate from a sequential reader | report section "Latency samples and derived throughput" |
| 7 | live gate tested ON only; empty input only in Go (R2) | gate adds an empty-input fixture and the filtered-child spelling `j11_optional_derived`; the WHERE control must still differ | `runJoinGate`; new catalog statement |

## Failing-before / passing-after

`harness/select_join_report_test.go` contains three tests written so the
defective behaviour fails them:

- `TestJoinInputsDigestReproducible` asserts the report advertises the file's
  SHA-256 and that the old double-newline variant differs from it.
- `TestJoinProvenanceSeparatesClientAndEngine` asserts the client and engine
  images occupy different fields.
- `TestJoinReplicatedPartitionsByClass` asserts a compatibility-class section
  exists and two different dataset fingerprints are named separately.

`go vet ./... && go test ./...` pass on the corrected tree (root: no test files;
`charitytree/harness` and `charitytree/selectbench`: PASS).

## Limits of this validation

- It does not validate the absent domain arm, blended phase or cluster subset;
  those are C2, C5 and C4 and remain open.
- It does not re-measure anything. The six original runs and their tags are
  retained; new runs will carry new ids.
- The fixed re-render used the tracked original input; it is a rendering check,
  not a new measurement.
