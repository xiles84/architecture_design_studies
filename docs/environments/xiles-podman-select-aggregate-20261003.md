# Environment: xiles-podman-select-aggregate-20261003

Aggregate / per-parent top-N sweep for the Study 01 SELECT experiment
(`task-20261002T131729Z-select-aggregate-topn`) on the documented
[Zenbook UX5406SA host](host-zenbook-ux5406sa.md), Windows/WSL2 Podman.

Single-node PostgreSQL 17.11 or YugabyteDB 2025.2.6.0-b111: 2 CPU/3 GiB.
Client: 2 CPU/2 GiB. Aggregate 4 CPU/5 GiB; the guest ceiling is the real limit
and the runner records fresh guest CPU/memory plus effective container limits
and before/after CPU throttling counters.

**Dataset.** Scale `sweep`: 10,000 deterministic customers, 2 orders each
(20,000 orders), seed 42, null-amount every 29th order, null-FK every 313th.
One fixture serves both axes through a `parent_max` prefix filter, so the levels
1, 10, 100, 1,000 and 10,000 parents are the protocol's parent-selection
{0.01%, 0.1%, 1%, 10%, 100%} and the top-N parent-count axis at once. The
smallest multi-child parent is given a tie at the top of its order
(`GuaranteedTie`) so the x06 RANK-tie control always has a tie to keep; that is
a fixture shape, not a workload change, and it is in the dataset fingerprint.

**Load and timing.** One sequential reader, no concurrent writers, pool ceiling
2 = one active reader/setup stream plus one reserved connection. Five trials per
cell follow independent-oracle verification and the live six-control gate. The
initial query and EXPLAIN are outside client timing. Planning/execution server
time is separate from client round-trip/drain time. Decoded-row JSON length is
not wire length. Index definitions and engine settings accompany each run. Only
fresh, tagged, committed source may run; exact images are retained.

**Limits.** One laptop: no independent hardware, no real network, no 3-node
per-node/equal-total arm, no verified placement, no demand generator, and no
throughput claim. Latency of tiny per-parent answers; five samples cannot support
p95/p99. YugabyteDB storage/RPC counters are not parsed (raw EXPLAIN retained).
A single whole run establishes no performance winner or prediction verdict.

NEXT MODEL: LOW
