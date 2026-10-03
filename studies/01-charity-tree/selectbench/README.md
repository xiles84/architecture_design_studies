# Additive SELECT harness

`run-select.sh` executes single-node validation in PostgreSQL or YugabyteDB.
Every run requires committed source, a new immutable run tag and a named
environment. It builds a unique client image with vet/tests gating it, takes
`ads-run-lock`, verifies live teaching/empty/all-miss/keyset tie/NULL fixtures
and all four wrong SQL controls, then records drained-query samples and plans.

```bash
studies/01-charity-tree/run-select.sh --engine postgres --scale tiny \
  --run-id select-pg-check-20261003 \
  --tag run/study-01/select-pg-check-20261003 \
  --environment xiles-podman-select-validation-20261003 --trials 3
```

Use a fresh id/tag for a new run. The client refuses missing runner lock/provenance
metadata. Results retain JSON, factual Markdown, manifest, image/node/machine
inspection, guest resources, logs and CPU throttling counters. Every timing sample
must match independent Go truth; any later mismatch discards the whole cell's
samples. Read-path controls need different answers under the same column shape;
query or decoding errors do not count as detected controls.

Ordering contracts preserve arrival order for pagination and top-N. Keyset uses
amount DESC NULLS LAST/id ASC and supports tie continuation and NULL cursors.
Plan structure preserves hierarchy and join roles. Unmeasured fields render N/A.
Report medians describe this closed-loop validation; they establish no replicated
winner. General performance guidance still requires the measuring tasks' full
three-run/three-trial/sweep/position/resource/placement controls and signed analysis.

No historical Study 01 mode or published result is migrated. Cluster/equal-total/
placement, balanced endpoints, mutation/storage costs and demand calibration remain
explicit measurement-stage gaps. Distributed RPC counters are retained only in
raw EXPLAIN where the engine supplies them; no parsed RPC metric is fabricated.
