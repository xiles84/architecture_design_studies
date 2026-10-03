# Observed resource envelope

Recorded during the serial v5 matrix on 2026-10-03, GPT-6 via Codex desktop.
The earlier short control run first exposed a distinction between saved Podman
configuration and the WSL guest. The producing source and candidate resource
limits have not changed during the matrix.

Saved machine configuration: 4 CPU/2 GiB. Guest observation: 8 online processors
and `MemTotal=16109780 kB` (about 15.36 GiB); each run retains its own memory
snapshot. PostgreSQL is capped at 1 CPU/512 MiB, Redis at 1 CPU/256 MiB, and
the client at 1 CPU/512 MiB. Together they budget 3 CPU/1.25 GiB.

The protocol's planned configuration-based reserve is 1 CPU/768 MiB. Relative
to the observed guest total, the arithmetic remainder is instead 5 online
processors/about 14.11 GiB. Neither remainder promises exclusive host resources:
WSL, Windows and other applications share the host. Memory availability and
effective per-container quotas/counters are the actual retained observations.

The experiment deliberately remains a single-node protocol comparison.
Its unmeasured cluster gap must not be described as proof that this host cannot
run a cluster: a larger actual guest envelope would need a separately calculated,
validated cluster/client plan and new premeasurement protocol. No clustered
or hardware-independent conclusion follows from the current cells.

NEXT MODEL: HIGH
