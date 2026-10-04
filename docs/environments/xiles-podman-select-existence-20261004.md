# Environment: xiles-podman-select-existence-20261004

Membership and absence sweep for the Study 01 SELECT experiment
(`task-20261002T131729Z-select-existence`) on the documented
[Zenbook UX5406SA host](host-zenbook-ux5406sa.md), Windows/WSL2 Podman.

Single-node PostgreSQL 17.11 or YugabyteDB 2025.2.6.0-b111: 2 CPU/3 GiB.
Client: 2 CPU/2 GiB. Aggregate 4 CPU/5 GiB; the guest ceiling is the real limit
and each run records fresh guest CPU/memory, effective container limits and
before/after CPU throttling counters.

**Dataset.** A fanout sweep at constant total child volume: 10,000 orders per
level, with child fanout {1, 10, 100} and therefore {10,000, 1,000, 100}
customers, seed 42, null-amount every 11th order and null child key every 97th.
Each level is loaded fresh in the same run; its fingerprint and row counts are
recorded. The nullable child key is what makes the `NOT IN` contract observable.

**Load and timing.** One sequential reader, no concurrent writers, pool ceiling
2. Five trials per cell after independent-oracle verification and the live
control gate. The initial query and EXPLAIN are outside client timing. Planning/
execution server time is separate from client round-trip/drain time. Raw
EXPLAIN plans are retained per cell.

**Limits.** One laptop: no independent hardware, no real network, no 3-node
per-node/equal-total arm, no verified placement, no demand generator and no
throughput claim. Qualifying proportion, parent hit fraction and hot-parent skew
are not swept in this bounded run. YugabyteDB RPC/storage counters are not
parsed. A single whole run establishes no performance winner or verdict.

NEXT MODEL: LOW
