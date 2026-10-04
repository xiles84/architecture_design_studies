# Progress — task-20261002T131729Z-select-join-cte (corrective attempt)

Attempt `attempt-20261004T173630Z-d9091f`, claim `claim-03114d688bd5fe53`,
branch `study-01/select-join-cte`, LOW executor (model `deepseek-flash`, tool
opencode/CLI, effort unknown). Resumes after the independent HIGH review
`docs/select-query-templates/reviews/20261004-join-high/` moved the task to
`changes_requested`. The previous two attempts, the six original runs, their
reports, signed analysis and all ten producing tags are preserved unedited.

## Corrective checklist (from CORRECTIVE-HANDOFF.md)

- [x] **C1 reporting/evidence fixes**
  - inputs digest: one shared `joinInputsDigest`; the re-render no longer appends
    a second newline. Re-rendering the original pg1 input now advertises the
    original `070df793…` (the review's broken re-render showed `7c879564…`).
  - client vs engine provenance separated (`ClientImage`/`ClientImageDigest` vs
    `EngineImage`/`EngineImageDigest`); the runner passes `ADS_ENGINE_IMAGE`.
  - replicated report keyed by a compatibility class (environment, dataset,
    client and engine image digests, settings) in addition to engine/topology/
    axis; classes are printed and never pooled.
  - engine-specific EXPLAIN options probed and recorded (PG VERBOSE/BUFFERS, YB
    DIST/DEBUG) with fallback evidence.
  - token volatility evidence recorded (distinct ids → distinct tokens; reuse
    plan carries the expression).
  - latency min/median/max and derived queries/s added, labelled derived (not a
    capacity test).
  - live gate extended: empty-input fixture and the filtered-child LEFT JOIN
    spelling `j11_optional_derived`, plus the existing firing WHERE control.
  - failing-before/passing-after tests in `harness/select_join_report_test.go`.
- [~] **C2 Study 01 domain arm** (person/donation/charity) + re-run matrix.
  - [x] Domain fixture/DDL/views implemented (`selectbench/joindomain.go`); the
    catalogue is shared unchanged. Projection, generated keys and rank proxy
    documented in `docs/select-query-templates/reviews/20261004-join-corrective/DOMAIN-ARM.md`.
  - [x] PG live gate passes; exact threshold counts and rank-index evidence;
    `cte-domain` runs `select-join-domain2-feas-pg-20261004`.
  - [~] Six validation whole runs measured (3 PG, 3 YB), 91/91 valid each, 5
    samples and position control; report
    `reports/20261004T-select-join-domain-replicated.md`. They span two client
    builds (report-only changes landed between yb1 and yb2-3), so the completed
    matrix must be re-run from one frozen commit; that frozen 3+3 matrix is open.
- [~] **C3 effective WSL resource recheck** and budget calibration.
  - [x] Attestation: 8 CPU / 16,109,780 kB effective (the 4 CPU / 2048 MiB podman
    config is stale metadata). Budget plan recorded in
    `docs/select-query-templates/reviews/20261004-join-corrective/RESOURCES.md`;
    the D8 resource rejection is retracted.
  - [~] One-CPU client calibration: bounded read-only pass done
    (`select-join-calib-pg-20261004`, 0 throttled periods at 1 CPU); the
    10-updates/s writer and widest-threshold demand still to be calibrated.
- [~] **C4 bounded single/cluster topology subset**, round-robin endpoints, RF3,
      placement and colocation.
  - [x] Cluster reachable under the calibrated budget: 3 tservers, RF3, per-node
    2 CPU/3 GiB (`CLUSTER-BRINGUP.md`).
  - [x] Multi-endpoint client implemented (`harness/join_cluster.go`); first
    cluster3 subset run 30/30 valid with endpoint ops 123/122/122 and per-node
    throttling 20/82/11 (`select-join-topo-yb-smoke-20261004`).
  - [ ] Bounded subset matrix: 3 whole runs per arm (single, equal-total single,
    cluster3), colocated vs non-colocated placement.
- [~] **C5 matched blended read/write phase** (one fixed writer).
  - [x] Implemented: independent `maintenance_updates` counter outside the
    projection/rank/token; one fixed writer at 10 updates/s over 1000 keys;
    isolated and blended phases; offered/completed/failed, latency, per-key
    spread, Go ledger and reader-write overlap recorded.
  - [x] Validation pair on PG single-node: both 30/30 valid; blended offered 66,
    completed 66, failed 0 (10.002/s), ledger 0 mismatch, 23/90 reader ops
    overlapped a write (`BLENDED.md`).
  - [ ] Full 5-trial phase pairs on the domain matrix and the topology subset.
- [ ] **C6** new signed successor analysis with corrected verdicts.

## Decision Log (this attempt)

1. **Restore the owner-normalized CRLF copies of HIGH's review files to the
   committed LF bytes** rather than commit a line-ending-only change to another
   session's signed review packet. No content change was present.
2. **Do not overwrite the six original reports or the original signed analysis.**
   The fixes are validated by a corrected re-render under the new review outputs
   and by tests; new runs use new ids, as the handoff requires.
3. **Explain options are probed, not assumed**, and the fallback is recorded as
   evidence, so an absent RPC counter is reported as an absence rather than
   described as "not parsed".
4. **Compatibility class rather than rejection** for the replicated report:
   runs with differing environment/fixture/settings are partitioned explicitly,
   matching the handoff's "reject … or partition" and keeping one report usable
   across the intended single/cluster and blended arms.
5. **`selectbench.Provenance` gains two fields** (`ClientImage`,
   `ClientImageDigest`); the change is additive and only the join family sets
   them, so the teaching families' provenance is unchanged.

## Still open / carry forward

C2–C6 are not done in this checkpoint. They need new fixtures/runs (domain arm,
blended phase, cluster subset) and a new signed interpretation. The next session
resumes from this decision log with the same claim lifecycle.
