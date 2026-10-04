# Effective-resource attestation and calibrated budget (review R1)

LOW `deepseek-flash` / opencode CLI, 2026-10-04, task
`task-20261002T131729Z-select-join-cte`, attempt
`attempt-20261004T173630Z-d9091f`. Read under the benchmark lock with no database
running; the lock was released immediately after.

The first submission deferred the 3-node and equal-total arms as a resource
rejection. The review showed the saved runtime evidence contradicts that, and it
does. This attempt retracts the rejection: the arm is **unattempted**, not
infeasible.

## Attestation

| Source | Value |
|---|---|
| host `/proc/meminfo` `MemTotal` | `16,109,780 kB` (~15.36 GiB) |
| host `MemAvailable` at check | `14,058,712 kB` |
| host `SwapTotal` | `4,194,304 kB` |
| host `/sys/devices/system/cpu/online` / `nproc` | `0-7` / 8 |
| `podman machine inspect` `Resources` | 4 cpus, 2048 MiB — **configuration metadata only** |
| saved container `guest-memory.txt`, run `select-join-pg1-20261004` | `MemTotal: 16,109,780 kB` |
| saved container `guest-cpus.txt`, run `select-join-pg1-20261004` | `0-7` |

The container sees the same 8 CPU / 15.36 GiB as the host; the 4 CPU / 2 GiB
machine configuration is stale and must not be used to size budgets or to
declare an arm infeasible. All corrected runs re-read these files.

## Calibrated budget for the corrected arms

Ceilings/reserves, not proof of attainable load:

| Arm | Database | Client | System reserve | Total |
|---|---|---|---|---|
| Single-node, per-node baseline | 2 CPU / 3 GiB | 1 CPU / 2 GiB | 1 CPU / 2 GiB | 4 CPU / 7 GiB |
| Single-node, equal-total control | 6 CPU / 9 GiB | 1 CPU / 2 GiB | 1 CPU / 2 GiB | 8 CPU / 13 GiB |
| 3-node cluster, per node (RF3) | 2 CPU / 3 GiB each | 1 CPU / 2 GiB | 1 CPU / 2 GiB | 8 CPU / 13 GiB aggregate |

Master/coordinator memory is inside the database total (`yb-cluster3.sh` pins the
master limit separately; the aggregate 3 × 3 GiB already contains it). These
corrected budgets are **separate from the original 2-CPU-client observations**,
which are retained unchanged.

## Exact invocation (no defaults changed, so other studies are unaffected)

`infra/lib.sh` already honours the environment overrides:

```bash
# equal-total single-node control: 6 CPU / 9 GiB database, 1 CPU / 2 GiB client
ADS_DB_CPUS=6 ADS_DB_MEMORY=9g ADS_CLIENT_CPUS=1 ADS_CLIENT_MEMORY=2g \
  bash studies/01-charity-tree/run-select.sh --family join --engine postgres \
  --scale <domain-scale> --trials 5 --position-control ...

# 3-node per-node baseline: defaults already give 2 CPU / 3 GiB per node;
# only the client moves to 1 CPU
ADS_CLIENT_CPUS=1 ADS_CLIENT_MEMORY=2g \
  bash studies/01-charity-tree/run-select.sh --family join --engine yugabyte \
  --topology cluster3 --scale <domain-scale> --trials 5 --position-control ...
```

## Open calibration items (must precede the corrected runs)

1. Confirm the one-CPU client can supply the offered demand and record its
   throttling (`client.cpu.stat`), including the 10-updates/s blended writer.
2. Confirm the 3-node cluster starts under the aggregate budget, RF3 is effective
   and per-node throttling is recorded.
3. Round-robin endpoint selection and per-endpoint operation/connection counts
   must be implemented in the client before a cluster arm is measured; a single
   seed endpoint is not balanced cluster access.
4. Colocated vs non-colocated schema placement, with a feature check and error
   evidence if the pinned image cannot express it.

Item 1 is a C3 follow-up; items 2–4 are C4. Neither the deferral nor an invented
placement success is acceptable; if a feature genuinely fails, the failure is the
recorded gap.

## One-CPU client calibration (read-only, first pass)

Run `select-join-calib-pg-20261004` (tag
`run/study-01/select-join-calib-pg-20261004`), PostgreSQL 17.11 single-node,
`cte-small` (20,000 facts), one sequential reader, no writer, feasibility gate
only:

| Field | Value |
|---|---|
| client budget (manifest) | 1 CPU / 2 GiB (`100000 100000` cgroup quota) |
| client `nr_throttled`, before → after | 0 → 0 |
| database budget (manifest) | 2 CPU / 3 GiB |
| explain options probed | `ANALYZE, VERBOSE, BUFFERS, FORMAT TEXT` |
| token evidence | `token(id=1)=49 token(id=2)=97 distinct=true`; default, MATERIALIZED and NOT MATERIALIZED reuse plans all carry the expression |
| client vs engine image | client `localhost/charitybench-select:select-join-calib-pg-20261004`; engine `docker.io/library/postgres:17.11` (separate fields) |
| controls | `x01_optional_where` and `x02_join_limit_before` both fired |

This is a bounded read-only pass: it does not exercise the 10-updates/s writer,
the widest threshold on the 200,000-row fixture, or the cluster. Those remain
items 1–4 above.
