# Environment: `xiles-podman-4cpu-2g-config-protocol-20261003`

Study 05's bounded v5 protocol condition, documented 2026-10-03 by GPT-6 via
Codex desktop. Do not pool these figures with earlier donor-portal throughput.

| Property | Observed value |
|---|---|
| Host | ASUSTeK ASUS Zenbook S 14 UX5406SA_UX5406SA |
| CPU | Intel Core Ultra 7 258V; 8 physical/8 logical cores |
| Host RAM | 33,804,754,944 bytes; approximately 31.48 GiB |
| Host OS | Microsoft Windows 11 Home, 10.0.26200, build 26200 |
| Podman machine | podman-machine-default; WSL2, rootless |
| Saved machine configuration | 4 CPUs, 2048 MiB, 100-GiB disk |
| Guest observation | 8 online CPUs; MemTotal 16,109,780 kB, approximately 15.36 GiB |
| Guest kernel | 6.6.87.2-microsoft-standard-WSL2, from Redis INFO |
| Source | PostgreSQL 17.11, endpoint cp-pg:5432 |
| Cache | Redis 7.4.11, endpoint cp-redis:6379; no persistence |
| Go build image | golang:1.26-bookworm |
| Client image id | 54d66b9c85d21b0d562c076f2c8cf2ddd7ce74755979c09461db416c926ed2f0 |

Host CPU/RAM/model/OS were re-read through Windows CIM after the matrix.
The baseline storage/core-topology description is in
[host-zenbook-ux5406sa](host-zenbook-ux5406sa.md); storage and clocks were not
re-measured here. Each run retains guest memory, machine/container inspect
and exact image identities under its `resources/` directory.

| Role | Effective budget/configuration |
|---|---|
| PostgreSQL | 1 CPU/512 MiB; shared_buffers 64 MiB, effective_cache_size 256 MiB, work_mem 4 MiB, maintenance_work_mem 64 MiB, max_connections 40 |
| Redis | 1 CPU/256 MiB; maxmemory 64 MiB, volatile-lru, samples 5; value TTL 300 s, absence-marker control TTL 100 ms; correct fences have no expiry |
| Client | 1 CPU/512 MiB; 8 readers, 2 writers, 1 relay, 3 spare connection slots; pool ceiling 14 |
| Service aggregate | 2 CPU/768 MiB, identical across candidates |
| Aggregate with client | 3 CPU/1.25 GiB |

Configuration-based reserve: 1 CPU/768 MiB. Relative to guest total: 5 online
processors/about 14.11 GiB. Both share the host with Windows/WSL applications;
saved configuration is not the observed guest ceiling. Quotas use CFS, not
pinning. The maximum whole-cell client throttled-period fraction was 68.86%;
the same-candidate first/last read ratio reached 2.2446. No isolated server
capacity or hardware-independent result is established.

Cluster and placement comparisons remain unmeasured, not proof of physical
resource impossibility. The 128-key projection does not test eviction pressure,
long-lived retention, payload scaling or source replication.

TASK COMPLETE
