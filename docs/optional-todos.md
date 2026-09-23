# Optional future todos — work that needs resources this machine does not have

**Owner policy, 2026-09-23:** all work in this project is executed on the owner's current machine.
Anything that needs a resource the machine does not have is **not** expected work; it is parked here
so it is not forgotten if the resources appear later. Nothing in this file is authorized work.

The queue already holds each of these as a published `proposed` task. A task listed here must not be
released for execution until the resource named in its row exists **and** the owner says so. The
corresponding task directory carries a `DEFERRED.md` repeating the reason.

## What this machine currently has

- One host (Windows/WSL2, 8 heterogeneous cores) running one container at a time under the benchmark
  lock; PostgreSQL 17.11 and YugabyteDB 2025.2.6 images are local.
- Podman only; one measurement at a time; no second physical or virtual host; no real network path
  between nodes (containers on one host share the kernel and a loopback bridge).

## Deferred items

| Queue task | Needs | What it would answer |
|---|---|---|
| `task-20260922T202001Z-topology-env-harness` | a second host | Study 07 environment page, balanced-endpoint instrumentation and fault helpers |
| `task-20260922T202001Z-topology-placement-routing` | two hosts (or a cluster that can place data on distinct hosts) | real data placement and routing, not two containers on one laptop |
| `task-20260922T202001Z-topology-node-replication` | two or three hosts | replication and durability across node failure between hosts |
| `task-20260922T202001Z-topology-network-failure` | a real inter-host path with `netem` | latency, loss and partition behaviour between hosts |
| `task-20260922T203143Z-native-models-harness` | MongoDB, ScyllaDB, Valkey images | a second implementation of the harness for non-relational engines |
| `task-20260922T203143Z-native-document-vs-relational` | MongoDB | document store versus PostgreSQL for the same workload |
| `task-20260922T203143Z-native-widecolumn-vs-relational` | ScyllaDB (or Cassandra) | wide-column versus relational under the same operations |
| `task-20260922T203143Z-native-keyvalue-vs-relational` | Valkey (or Redis) | key/value cache-store versus relational |
| `task-20260922T204130Z-analytics-harness` | ClickHouse and a search engine image | the analytical harness for Study 08 |
| `task-20260922T204130Z-analytics-columnar-copy` | ClickHouse | columnar read copy versus the row store |
| `task-20260922T204130Z-analytics-search-copy` | a search engine (OpenSearch/Elasticsearch) | inverted-index copy versus the row store |
| `task-20260922T204130Z-analytics-rollup-arms` | ClickHouse | rollup arms against the analytical copies |

## How to re-activate

1. The resource exists and is documented in `docs/environments/`.
2. The owner approves the work (a task is released only by HIGH, and this file records that it is
   deferred until the resource exists).
3. The task's `DEFERRED.md` is removed and the normal queue lifecycle applies: release, claim,
   execute, review, integrate, tag.

## Not deferred

Studies 04, 05 and 09 run on this machine with the current images and are not listed above. A task
that turns out to need a resource this machine lacks during execution is added here rather than
half-executed.
