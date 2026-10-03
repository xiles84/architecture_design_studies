# Environment: xiles-podman-select-validation-20261003

Bounded SELECT harness validation on the documented
[Zenbook UX5406SA host](host-zenbook-ux5406sa.md), Windows/WSL2 Podman.
Current guest configuration/observations are recorded for each run: the saved
machine profile is 4 CPUs/2048 MiB; the previous live guest observation was
8 online CPUs/about 15.36 GiB. The saved profile is not the observed ceiling.

Single-node PostgreSQL 17.11 or YugabyteDB 2025.2.6.0-b111: 2 CPU/3 GiB.
Client: 2 CPU/2 GiB. Aggregate 4 CPU/5 GiB; relative to an 8-CPU/15.36-GiB
guest, reserve is 4 online processors/about 10.36 GiB for guest/system activity.
The runner records fresh guest CPU/memory plus effective container limits and
before/after CPU throttling counters. Changes in guest totals remain explicit.

One sequential reader, no concurrent writers, pool ceiling 2 = one active
reader/setup stream plus one reserved connection. Three drained-query samples
per template follow independent-oracle verification and live fixture/control
gates. Initial query and EXPLAIN are outside client timing. This is closed-loop
validation, without a demand-generator capacity or crossover claim. Planning/
execution server time is separate from client round-trip/drain time. Decoded-row
JSON length is not wire length. Index definitions and engine settings accompany
each run. Only fresh, tagged, committed source may run; exact images are retained.

No 3-node, equal-total-budget, placement, balanced-routing or mutation-cost
comparison is established by this validation. Those remain in the measuring
tasks. A single whole run establishes no performance winner or prediction verdict.

NEXT MODEL: HIGH
