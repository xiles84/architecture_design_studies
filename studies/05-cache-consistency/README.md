# Study 05 — external cache throughput and consistency

**Status:** in progress (see `HANDOFF.md`, `PROGRESS.md`, `ESCALATIONS.md`)
**Environment:** [`host-zenbook-ux5406sa`](../../docs/environments/host-zenbook-ux5406sa.md)
**Engines:** PostgreSQL 17.11 · YugabyteDB 2025.2.6.0 (1 node and 3 nodes, RF=3) · Redis 7.4.11-alpine

## The question

A donor portal reads the same donor's page over and over, and a small share of requests mutate
data under it. An external cache is the obvious answer; every cache is also a copy that can be
wrong. This study measures **what the cache buys** and **what freshness costs**, on a database
model we own and may modify and on a legacy model that cannot be modified for caching.

It reuses Study 01's `charity → person → donation` domain and its D3 flattened-FK schema as the
source of truth. The primary cache key is `donor:<person_id>:portal`; the cached value is a
coherent committed view: person identity and mutable metadata, charity identity, donation count,
donation total, the newest 20 donations in deterministic order, a source version or external
generation, a content hash, a creation time and a hard-expiry deadline.

**A cache value may be published only if it corresponds to a successfully committed database
state.** That is the floor every scenario is built on, including the ones that are allowed to
return a stale value.

### Terminology, used precisely

- **Dirty cache write** — publishing speculative, rolled-back, ambiguous or otherwise
  uncommitted data into the cache. Forbidden in every candidate scenario.
- **Stale / wrong read** — returning an older committed value than the latest write to that key
  acknowledged before the read began. Permitted only by `relaxed`; counted in `strict`, where it
  fails the cell.
- **Concurrent / ambiguous read** — a read overlapping a write. Counted separately; never
  silently classified as fresh or wrong.
- **Impossible cache value** — a payload or version corresponding to no committed database
  state. Must always be zero.

"The cache returned a stale read" is what this study says. A database dirty read is a different
thing.

## The two database models

| Model | What the database provides | Where its SQL lives |
|---|---|---|
| **owned** | `person.cache_version` (monotonic), a transactional `cache_outbox`, an atomic version change with every mutation, ordered locking for reassignment | `sql/owned/d3_owned/` |
| **legacy** | unmodified D3: no cache-specific version, no outbox, no trigger, no CDC, no reliable `updated_at` token, schema frozen | `sql/legacy/d3_legacy/` |

Legacy is measured in two write regimes: **coordinated** (all writes pass through the caching
service) and **external** (20 % of harness writes bypass the adapter while committing correctly
to the database).

## The scenarios

Scenario id: `<model>-<ver>-<backend>-<strategy>-<freshness>-<writers>`.
Backends `none|memory|redis`; strategies `none|aside|through`; freshness `none|relaxed|strict`;
`ver` is `opt|pess` for the owned model and `na` for legacy; `writers` is `coord` or `ext20`.

The cache scenarios are **dimensions over one SQL catalogue per model**, not 18 copies of the
same SQL. On PostgreSQL the core matrix is 21 cells: two owned no-cache baselines (optimistic and
pessimistic), one legacy no-cache baseline, sixteen cached combinations
(`{owned,legacy} × {memory,redis} × {aside,through} × {relaxed,strict}`), and two legacy external
writer cells. Four database-reference designs — `ref-normalized-indexed`, `ref-flattened-fk`
(the principal baseline), `ref-rollup-trigger`, `ref-embedded-locked` — are measured without a
cache in the same run, so the rollup, rolldown and embedding comparisons are attributable to the
database layout and not to the cache.

### Cache policy, identical in every cached scenario

- byte-bounded exact LRU in memory; Redis with `maxmemory`, `allkeys-lru` and
  `maxmemory-samples` recorded (Redis's LRU is approximate, and the report says so);
- an exact **300-second** hard TTL on every entry;
- probabilistic early expiration on every valid hit,
  `p = clamp(1 - time_remaining/300 s, 0, 1)`, drawn from the trial's recorded random source;
- a per-key fill lease — Redis `SET NX PX` with compare-and-delete release, in-memory per-key
  coordination with bounded waiting — with waits, losses, duplicate fills and fallback reads
  recorded, and version fencing so an expired lease holder cannot overwrite a newer committed
  version;
- a primary capacity smaller than the working set so eviction actually happens, plus a
  fits-working-set control.

### Correctness, measured rather than asserted

Timed cache hits are never validated by a hidden synchronous full database read. An independent
Go oracle and operation ledger record every acknowledged write's key, committed version and
expected payload hash, and capture the latest acknowledged write at each read start. The ledger is
replayed against the final database state after every phase, and wrong reads are reported with
their count, rate, staleness distribution, affected keys and evidenced cause.

Both controls are measured like designs and never reported as fast results:
`ctl-stale-invalidation` (one suppressed post-commit invalidation must produce a stale-after-ack
read) and `ctl-unchecked-rmw` (an unguarded read-modify-write must lose updates). A third,
`ctl-dirty-write-audit`, injects a rolled-back or nonexistent version in a dev check and requires
the audit to reject it.

## How to run it

```bash
# dev checks only; numbers are never reported
./run-study.sh --scale tiny --topologies pg-single \
  --phases verify,explain,warm,mixed,hotspot,stampede,instances,faults,audit

# the reported primary survey (core matrix + reference designs), PostgreSQL single node
./run-study.sh --tag --scale small --topologies pg-single

# the two real five-minute TTL boundaries, one instance, sustained churn
./run-study.sh --tag --scale medium --topologies pg-single --designs <scenario> \
  --phases verify,warm,mixed,churn --churn-duration 11m
```

Every runner takes the machine-wide benchmark lock, pins the benchmark and Redis image ids,
records the committed source, the environment, the resource framing and the seeds, and writes one
JSON result per cell plus a readable plan file. Results live in `results/<run-id>/`; the generated
report is `reports/<run-id>.md` and contains measurements only. Conclusions are separate signed
files under `reports/analyses/`, with mechanism detail in `reports/discussions/`.

## Layout

```
HANDOFF.md  PROGRESS.md  ESCALATIONS.md  README.md  study.env  Containerfile
run-study.sh  probe-cache.sh
sql/legacy/d3_legacy/   sql/owned/d3_owned/   sql/reference/{normalized_indexed,flattened_fk,rollup_trigger,embedded_locked}/
harness/    the Go benchmark client (main.go is the entry adapter)
results/    reports/{analyses,discussions,outdated}/
```

Shared, study-independent measurement code lives in `platform/`; study-specific cache code lives
here. Only `harness/main.go` wires in a concrete database adapter.
