# Environment: xiles-podman-select-access-20261003

Access-path measurement runs for Study 01's additive SELECT experiment, on the
documented [Zenbook UX5406SA host](host-zenbook-ux5406sa.md), Windows/WSL2
Podman (`podman-machine-default`). This is the same measurement slot described by
[`host-zenbook-ux5406sa.md`](host-zenbook-ux5406sa.md); it is named separately so
an access-path result is never confused with the earlier bounded harness
validation.

## Budgets and topology

- Single-node PostgreSQL 17.11 or YugabyteDB 2025.2.6.0-b111: **2 CPU / 3 GiB**.
- Benchmark client container: **2 CPU / 2 GiB**.
- Declared aggregate 4 CPU / 5 GiB on a guest that reports 8 online CPUs and
  about 15.36 GiB; the machine profile in Podman is 4 CPU / 2048 MiB. The
  declared memories (3 GiB + 2 GiB = 5 GiB) exceed the 2 GiB profile, so the
  guest ceiling, not the container limit, is the effective constraint. This is
  the standing sizing blocker; it is recorded per run and is not repaired by
  this study.
- One sequential reader, no concurrent writers, connection-pool ceiling 2 (one
  active reader/setup stream plus one reserved connection).
- Dataset scale `deep`: 200 customers × 60 orders = 12,000 rows, seed 42.

## What a run records

Fresh guest memory/CPU, effective container limits, before/after node and client
cgroup throttling counters, image digests, dataset fingerprint, engine settings,
per-cell EXPLAIN with normalized structure, and the fixed-snapshot page-equality
check. Planning/execution server time is separate from client round-trip/drain
time; decoded-row JSON length is not wire length.

## Coverage limits (explicit)

Single host and single node only. No 3-node RF=3 arm, no equal-total-budget arm,
no verified tablet/leader placement, no independent hardware and no real network
behaviour is established here. Those remain coverage gaps, not completed
comparisons. A single whole run establishes no replicated direction; each
measuring task runs three whole runs per engine.
