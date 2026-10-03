# Cache publication and deletion protocol, v5

Task: `task-20260923T235925Z-cache-untested-protocols`. Planner/executor: GPT-6,
Codex desktop, declared HIGH; selected effort unavailable. Written before the
measurement implementation. This extends the existing study with a bounded protocol
experiment; no older design, result, analysis or claim is changed.

## Question and answer contract

Which of authoritative version validation on every hit, synchronous publication
fencing, and outbox-delivered invalidation can satisfy Study 05's contract, and what
work does each move to the database, acknowledgement path and recovery path?

The fixed answer is one keyed record `(key, generation, value, deleted)`. A read
starting after acknowledgement must return that generation or a newer committed one.
An overlapping read can return an older committed state, classified separately;
an uncommitted or fabricated state is forbidden. A separate Go operation ledger
records intended mutations, commit outcomes and acknowledgement order. It derives
expected values independently and retains every committed state, including absence.

This deliberately small projection isolates the protocols. It is **not** a re-run
of the existing donor portal, its aggregates or recent-donation payload, and supports
no cross-run comparison with its throughput. Payload size and network topology are
fixed limits, not silently generalized.

## Arm A: controlled protocol comparison

All candidates share SQL, indexes, data, Redis store, connection pool and workload.
Every mutation updates a durable generation, mutates the object and appends an outbox
event in one transaction, even where that event is not consumed. Only the protocol
changes:

- `version`: validate a hit against the authoritative generation. A miss reads a
  committed snapshot. Writes do not invalidate the cache.
- `sync`: fence before the database mutation and again after commit, then acknowledge.
- `outbox`: commit; a relay consumes committed unprocessed events, fences the cache,
  records delivery and only then releases the writer's acknowledgement barrier.
  Delivery scans unprocessed rows, **not a sequence cursor** that can skip an earlier
  transaction committing later. Relay replay is idempotent.

The source transaction is crossed with optimistic CAS and pessimistic lock-before-
update. Writers serialize by key for the independent ledger's commit order (the same
default as Study 05); cross-key work is concurrent. This is a protocol comparison,
not evidence of intra-key database lock scalability. A lost-update control is explicit.

Each cell gates correctness before timing and again afterward. Phases: warm read-only,
write-only, and concurrent eight-reader/two-writer workload. Three trials per phase;
at least three whole runs with recorded seeded order and a repeated first candidate
at the end as a position control. Report separate read/write counts, p50/p99, source
queries per logical read, transaction and commit-to-acknowledgement latency,
contention, outbox polling/delivery and recovery lag. No clean timing survives a
correctness failure. Closed-loop capacity is a floor; no server knee is inferred.

## Arm B: deletion, recreation and tombstone expiry

Compare a retained database soft-delete row, physical database deletion with a
durable generation record, and physical deletion with a cached absence tombstone.
Hold the source-generation record and synchronous fencing constant: these are
deletion packages, not evidence that soft deletion alone creates coherence.

Pause an old filler after its committed snapshot and before publication. Delete and
acknowledge, then attempt the paused publication; expire the short **absence marker**,
then repeat; recreate the same logical key at a higher generation and repeat again.
The 300-second value TTL is unchanged. The absence marker is a separate 100-ms test
parameter; correct fencing metadata has no TTL. A deliberately expiring fence
control must exhibit the ABA/resurrection hazard. Report observed expiry, not elapsed
time alone, and classify reads against the independent ledger at every step.

Report deletion/recreation mutation latency separately by package with three trials
and three whole runs. Do not pool Arm B with Arm A's workload or rank an invalid
control as a fast result.

## Prediction

| ID | Expected result and mechanism | Falsifier |
|---|---|---|
| A1 | All three correctly implemented protocols produce zero stale-after-ack and impossible reads; the version check, two fences, or delivery barrier rejects superseded states. | Any such violation after the candidate's fired controls. |
| A2 | A warmed `version` hit adds one authoritative-generation query; sync/outbox hits need none. | A version hit requiring none, or a sync/outbox hit issuing that query. |
| A3 | Pausing the relay blocks outbox acknowledgement; resuming it releases the pending write with nonnegative recovery lag and no stale-after-ack read. | An acknowledgement before delivery or a stale read after it. |
| A4 | A relay without the acknowledgement barrier permits an acknowledged stale hit while delivery is paused. | The negative control cannot produce a stale-after-ack read. |
| A5 | Every comparison's direction is bounded by across-run and position-control spread, rather than assumed to rank the protocols. | A reported winner unsupported by those records. |
| B1 | Durable generations and non-expiring fencing metadata prevent old-fill resurrection across all three deletion packages. | Any observed old state after deletion/recreation acknowledgement. |
| B2 | Expiring the fence as well as the absence marker allows the old zero-fence filler to republish; soft deletion by itself does not repair that cache race. | The explicit ABA control does not publish and serve the superseded state. |
| B3 | Injecting a value the ledger never committed is detected as impossible in every relevant arm. | An injected uncommitted value passes the independent check. |

## Environment, sizing and required-comparison coverage

The existing Podman machine was stopped and has been restarted unchanged: configured
four CPUs, 2 GiB. Record machine inspect, guest memory, effective container budgets,
engine/image digests, Redis eviction/memory counters and per-container throttling.
Do not equate the configured VM memory with guest-observed WSL memory.

The cached alternatives use the same equal-total **service** budget: PostgreSQL
one CPU/512 MiB and Redis one CPU/256 MiB. Client one CPU/512 MiB; reserve one CPU/
768 MiB, totaling four CPUs/2 GiB. PostgreSQL uses shared_buffers=64 MiB, work_mem=4
MiB, maintenance_work_mem=64 MiB and max_connections=40. Redis maxmemory=64 MiB,
volatile-lru, persistence disabled; only value/absence entries expire or evict,
never fencing metadata. Keys=128, deterministic dataset seed 42. Eight readers,
two writers, one relay, plus three pool-reserve connections: pool ceiling 14.

Record actual worker supply, closed-loop rate, pool allowance, errors and client
throttling; a client-limited result remains labelled as such. All alternatives have
identical aggregate CPU/memory; this is not an add-cache-vs-no-cache measurement.

| Required comparison | Measured pair / applicability / gap |
|---|---|
| Single-node and three-node, sizing, equal-total | Three cache protocols on PostgreSQL single-node with equal service budgets. Three-node YugabyteDB is a resource gap: published multi-node failures and the configured 2-GiB envelope do not justify silently starting a larger cluster. No cluster claim follows. |
| Rollup/rolldown/embedding | Not applicable to this single-object publication question. Existing donor-portal evidence remains separately scoped; none of these storage decisions changes across a pair. |
| Colocated/non-colocated data | Single-node, so placement contrast is not applicable to measured cells. Physical multi-host and multi-node placement remain explicit unmeasured gaps. Record the actual single database endpoint. |
| Optimistic/pessimistic | CAS and lock-before-update crossed with all three protocols; same invariant and dataset. Per-key ledger serialization is a stated contention limitation. |

## Gates, provenance and review

Record a premeasurement HIGH self-review against this protocol. All runtime work is
containerized. Every database-starting script acquires `ads-run-lock`; never bypass
or remove a foreign lock. A unique task benchmark image protects another study's
existing image. Commit the source before a producing run; use annotated run tags,
immutable result files, generated factual reports and a new signed analysis carrying
the input digest and source commit. Preserve failed cells and source/placement gaps.

The owner renewed the instruction to execute the HIGH queue, including this task's
already specified comparative questions. That supplies the comparative need for its
conditional release. It does not approve the Edition 3 visual style or waive the
new-domain ordering gate. This phase remains entirely HIGH; no switch is needed.

NEXT MODEL: HIGH
