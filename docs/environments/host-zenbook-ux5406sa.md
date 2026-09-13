# Environment: `host-zenbook-ux5406sa`

Every result file in this repository names the environment it was produced on. This is
that environment's specification. Hardware materially changes benchmark numbers, so
results are only comparable to each other when the environment id matches.

**Environment id:** `host-zenbook-ux5406sa`
**First documented:** 2026-09-12

## Host machine

| Property | Value |
| --- | --- |
| Model | ASUSTeK ASUS Zenbook S 14 UX5406SA |
| CPU | Intel Core Ultra 7 258V ("Lunar Lake") |
| Cores / threads | 8 physical / 8 logical (no SMT) |
| Core topology | 4 P-cores (Lion Cove) + 4 LP E-cores (Skymont) — **heterogeneous** |
| Base / max clock | 2.20 GHz nominal, ~3.30 GHz observed in-VM |
| RAM | 31.48 GiB LPDDR5X-8533 (8 × 4 GiB, on-package, Micron) |
| Storage | WD PC SN560 SDDPNQE-1T00 — 954 GB NVMe SSD |
| OS | Windows 11 Home 10.0.26200 (build 26200.9445) |

## Container runtime

| Property | Value |
| --- | --- |
| Podman | 5.8.1 (client, Windows) |
| Machine provider | WSL2 |
| Machine name | `podman-machine-default` |
| Rootless | yes |
| Guest kernel | 6.6.87.2-microsoft-standard-WSL2 |
| WSL | 2.6.2.0 |
| CPUs visible in VM | 8 |
| Memory visible in VM | 16 496 414 720 B (≈15.4 GiB) — WSL2 default of 50% of host RAM |
| Compose provider | none installed — orchestration uses the plain `podman` CLI (see [replication.md](../replication.md)) |

## Images under test

| Role | Image | Digest pinned in | Notes |
| --- | --- | --- | --- |
| PostgreSQL | `docker.io/library/postgres:17.11` | `infra/versions.env` | PostgreSQL 17 |
| YugabyteDB | `docker.io/yugabytedb/yugabyte:2025.2.6.0-b111` | `infra/versions.env` | Stable track; PostgreSQL-15-derived query layer |
| Benchmark client | `docker.io/library/golang:1.26-bookworm` | `infra/versions.env` | Build stage for the Go harness |

## Known measurement hazards on this environment

### Study 01 v3 resource conditions

These are separately named experimental conditions on this same hardware. They must not
be pooled as though the database configuration were identical.

| Condition | Database budget | Additional configuration |
|---|---|---|
| standard | 2 CPU / 3 GiB per node | Original settings |
| constrained (PostgreSQL growth) | 2 CPU / 256 MiB | shared_buffers 64 MiB, effective_cache_size 192 MiB, work_mem 4 MiB, maintenance_work_mem 32 MiB |
| large (YB equal-total single node) | 6 CPU / 9 GiB | tserver hard limit 4.5 GiB, master hard limit 1.5 GiB |
| equal-total YB cluster | 3 × standard nodes = 6 CPU / 9 GiB total | tserver limits sum to 4.5 GiB, master limits sum to 1.5 GiB |

Client budget stays 2 CPU / 2 GiB. Each cell captures live container limits and cgroup
CPU/memory/I/O statistics before and after it; the result records the database memory
limit, preparation and engine's effective isolation. Relation bytes exceeding a memory
limit demonstrate the stored footprint exceeds that limit; neither that fact nor a
`pg_stat_database.blks_read` count alone proves every read reached the physical SSD.
The cgroup `io.stat` and actual plans supply additional evidence about storage work.

The node-stop diagnostic stops `yb-n3` while the client continues through `yb-n1`, then
restarts it and records the recovered donation count. It is a local non-query-node
availability test, not a network partition, loss of the query endpoint, or host failure.

These are **not** incidental caveats — they bound what the numbers in this repo mean.

1. **Heterogeneous cores.** Lunar Lake mixes 4 performance cores with 4 low-power
   efficiency cores. The guest kernel sees 8 undifferentiated CPUs and cannot tell
   Podman which is which. A container pinned with `--cpuset-cpus` may land wholly on
   E-cores and run materially slower than the same container on P-cores, with no
   signal in the output. **Consequence:** this repo does *not* pin cores. It uses
   `--cpus=<quota>` (CFS bandwidth) so every container gets an equal *share* of
   whatever cores it lands on, and it reports the spread across repeated trials so
   that scheduling luck shows up as variance instead of hiding inside a single number.

2. **Client and server share the machine.** The benchmark client runs in a container on
   the same host as the database. There is no isolated network, and client CPU time is
   taken from the same 8 cores. Absolute throughput is therefore lower than a
   two-machine setup would give. Relative comparisons between designs — which is what
   this repo is for — remain valid because every design pays the same tax.

3. **No real network latency.** All traffic is loopback inside one WSL2 VM. Distributed
   designs (YugabyteDB 3-node) look *better* here than they would across real availability
   zones, because cross-node consensus round-trips cost microseconds rather than
   milliseconds. Read any "sharding is cheap" conclusion with that in mind; the
   RPC *counts* reported by `EXPLAIN (ANALYZE, DIST)` are the portable signal, not the
   wall-clock deltas.

4. **Laptop thermals.** A 28 W-class mobile CPU in a fanless-ish chassis throttles under
   sustained load. Long runs drift slower than short ones. Mitigation: fixed-duration
   runs with a warmup phase, randomised design order within a run, and repeated trials.

5. **WSL2 filesystem.** Container volumes live inside the WSL2 VM's ext4 disk, not on a
   Windows bind mount. Bind-mounting from `C:\` would add a 9p/drvfs translation layer
   and destroy write numbers. All database data directories are podman named volumes.
