# Environment: xiles-podman-select-order-correction-20261004

Order-corrected SELECT matrices (`task-20261004T100822Z-select-order-correction`)
on the documented [Zenbook UX5406SA host](host-zenbook-ux5406sa.md), Windows/WSL2
Podman.

**Live resources supersede the earlier pages' stale 2 GiB guest claim.** The
independent check's recorded `MemTotal` for the runs was 16,109,780 kB
(~15.36 GiB) with CPU ids 0-7 (8 online). The earlier SELECT pages described a
~2 GiB guest; that figure is stale and is corrected here.

- Database container: 2 CPU / 3 GiB (single-node PostgreSQL 17.11 or
  YugabyteDB 2025.2.6.0-b111).
- Client container: 2 CPU / 2 GiB. Aggregate 4 CPU / 5 GiB; relative to an
  8-CPU / ~15.36-GiB guest, the reserve is ~4 online processors / ~10 GiB.
- The runner records fresh guest CPU/memory, effective container limits and
  before/after CPU throttling per run.

**Execution order.** Each whole run permutes its cells with a seed derived from
the dataset seed and the run id, and takes one sample of every cell per round
(interleaved rounds) after each cell is bound, oracle-verified and planned once.
The seed, round count and the actual axis sequence are recorded in the run JSON
and report. The optional position control repeats the *actual* first cell of the
run's randomized order; the access expression-index arm is a separate labelled
phase with its own recorded order.

**Limits.** Single host and single node: no 3-node per-node/equal-total arm, no
verified placement, no real network, no concurrent demand. The recorded
resources make a cluster budget calculable but no cluster topology was measured,
so the cluster comparison remains an explicit gap. One sequential reader; five
samples per cell cannot support p95/p99; client CPU throttling being zero is not
a client-capacity result. YugabyteDB storage/RPC counters are not parsed.

NEXT MODEL: LOW
