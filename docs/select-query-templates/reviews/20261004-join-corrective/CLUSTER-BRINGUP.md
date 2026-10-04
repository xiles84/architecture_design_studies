# Cluster bring-up under the calibrated budget (review R1)

LOW `deepseek-flash` / opencode CLI, 2026-10-04. Read under `ads-run-lock`, no
other benchmark running. This is the check the first submission skipped: the
3-node arm is **attainable**, not infeasible.

## Bring-up

```bash
source infra/lib.sh && run_lock_acquire "C4 yb-cluster3 bring-up ..."
export ADS_DB_CPUS=2 ADS_DB_MEMORY=3g
bash infra/yb-cluster3.sh up        # seed, then join n2, then join n3
```

The cluster reached ready in 47 s with:

| observation | value |
|---|---|
| tservers (`yb_servers()`) | 3 — `yb-n1`, `yb-n2`, `yb-n3` on 5433 |
| replication factor (`get_universe_config`) | `"numReplicas":3` |
| per-node cgroup `cpu.max` | `200000 100000` = 2 CPU |
| per-node cgroup `memory.max` | `3221225472` = 3 GiB |
| per-node `MemTotal` | `16,109,780 kB` (container view) |

`yb-cluster3.sh up` asserts 3 tservers and RF3 before returning, so an RF1
universe cannot be reported as replicated. The cluster was brought down and the
volumes removed, and the lock released.

## Calibrated topology budget (ceilings/reserves)

| Arm | Database | Client | System | Total |
|---|---|---|---|---|
| single-node per-node baseline | 2 CPU / 3 GiB | 1 CPU / 2 GiB | 1 CPU / 2 GiB | 4 CPU / 7 GiB |
| equal-total single-node control | 6 CPU / 9 GiB | 1 CPU / 2 GiB | 1 CPU / 2 GiB | 8 CPU / 13 GiB |
| 3-node cluster, per node (RF3) | 2 CPU / 3 GiB × 3 | 1 CPU / 2 GiB | 1 CPU / 2 GiB | 8 CPU / 13 GiB aggregate |

Master/coordinator memory is inside the 3 GiB per node (the node script pins
`MASTER_FLAGS=memory_limit_hard_bytes=536870912` and the tserver limit
separately). These calibrated arms are kept out of the original 2-CPU-client
observations by the provenance class.

## Balanced endpoints

`yb-cluster3.sh endpoints` now prints all three SQL endpoints comma-separated.
The join client takes a comma-separated DSN list, builds one pool per endpoint
and rotates successive statements across them (`harness/join_cluster.go`),
recording per-endpoint operation and connection counts in the result and the
report's "Endpoint distribution" table. A cluster run can therefore show it was
not served by the seed node alone.

## Still open

- The bounded single/cluster subset runs themselves (3 whole runs per arm), with
  the cluster arm using these endpoints.
- Colocated vs non-colocated domain schemas, with a feature check and error
  evidence if the pinned image cannot express the placement.
