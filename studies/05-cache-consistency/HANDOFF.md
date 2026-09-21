# Execution Handoff — Study 05, external cache throughput and consistency

| Field | Value |
|---|---|
| Handoff ID / revision | `EH-05` revision 1 |
| Planner | DeepSeek HIGH selected by user. Model id exposed by the session: `deepseek-flash`. Effort setting: **not exposed**. Tool identity: **not exposed**. WSL, one agent for planning, implementation, execution, validation and analysis (owner's instruction for this task). |
| Starting source / main revision | local `main` = `e8e4009c664b5d1a7d2a99688acd6b3a46caa361`; task worktree `.worktrees/study05-cache-consistency`, branch `study-05/cache-consistency` |
| Checkpoint tag | `study-05/v0-handoff` (annotated; never moved) |
| Status | ready |
| Next setting | **HIGH — DeepSeek HIGH implements and executes this handoff itself.** No LOW executor; that is stated again as a limitation in the analysis. |
| Progress / escalations | `PROGRESS.md`, `ESCALATIONS.md` (this directory) |

---

## 1. Objective and boundaries

Measure, on the repository's own hardware and in Podman containers, **what an external cache
buys and what freshness costs** for a database-backed donor portal, and **what it costs to be
correct** when the database may or may not be changed to help the cache.

**Question, in the owner's words:** how much throughput and latency improve across no cache,
in-process memory cache and shared Redis; cache-aside and write-through; bounded-stale reads
accepted and stale reads forbidden; on a database model we own and may modify and on a legacy
database model that cannot be modified for caching.

**Hard correctness floor for every candidate scenario: a cache value may be published only if
it corresponds to a successfully committed database state.** No scenario may publish
speculative, rolled-back, ambiguous or otherwise uncommitted data.

### Terminology — used precisely, in every artefact

| Term | Meaning | Requirement |
|---|---|---|
| **Dirty cache write** | publishing speculative, rolled-back, ambiguous or otherwise uncommitted data into the cache | forbidden in every candidate scenario |
| **Stale / wrong read** | returning an older committed value than the latest write to that key acknowledged before the read began | permitted only by `relaxed`; must be counted in `strict` and fails the cell |
| **Concurrent / ambiguous read** | a read overlapping a write | counted separately; never silently classified as fresh or wrong |
| **Impossible cache value** | a payload or version corresponding to no committed database state | must be **zero** in every cell; non-zero invalidates the cell |

"Stale read" is the term used in artefacts. A database dirty read is a different thing and is
not what this study measures.

### Owned paths

Everything under `studies/05-cache-consistency/`. One genuinely study-independent capability
may be added under `platform/` only if implementation proves a second study would need it; that
change carries its own tests and a `repo/*` milestone tag (see §7).

### Untouched

Study 01's harness, SQL, results and reports; studies 02, 03 and 04 in full; `infra/*` except
the additive changes §7 allows; the root `AGENTS.md` rules. Study 01's D9/D10 are **read as
correctness precedents only** — its measurements are historical context and are never presented
as Study 05 baselines.

---

## 2. Decisions made by HIGH

### 2.1 Domain and workload

The study reuses Study 01's `charity → person → donation` domain and **D3 flattened-FK as the
principal source-of-truth baseline**. Primary cache key: `donor:<person_id>:portal`.

The cached donor-portal representation contains a coherent committed view including at least:
person identity and mutable metadata (`full_name`, `email`, `joined_at`, `charity_id`); charity
identity (`name`, `country`); donation count; donation total; the newest 20 donations in
deterministic order `(donated_at DESC, donation_id DESC)`; a source version or external cache
generation; a content hash; cache creation time; and the hard-expiry deadline.

The payload is assembled **in Go** from two design-supplied statements that return the logical
columns (`r_portal_person`, `r_portal_recent`) so that the content hash is computed by one
canonical encoder for both cache values and oracle expectations. No SQL `jsonb_build_object`
output is hashed: key order and numeric formatting would make that a test of PostgreSQL's JSON
renderer, not of the cache.

Mutations, all required: donation insert; donation amount correction; donation deletion;
person metadata update; donation reassignment between two people (both cache keys handled).

### 2.2 Two database models, identical logical data and operations

**Owned** (`sql/owned/d3_owned/`) — the database may be changed for cache correctness:
`person.cache_version BIGINT NOT NULL DEFAULT 0`, a `cache_outbox` change-log table, an atomic
version change for every mutation affecting the portal view, and ordered locking for
reassignment touching two people. The version bump and the outbox event commit **in the same
transaction as the authoritative mutation**.

**Legacy** (`sql/legacy/d3_legacy/`) — the unmodified logical D3 schema, assumed to have no
cache-specific version column, no outbox, no new trigger, no CDC, and no reliable
cache-specific `updated_at` token. The schema cannot be changed for this study. The harness may
keep **external metadata** (an oracle generation counter in measurement metadata); it must never
pretend cache-side metadata can detect an unknown external writer.

Two legacy write regimes:
1. **coordinated** — every write passes through the caching service;
2. **external** — a controlled fraction of harness writes bypasses the cache adapter while
   still committing correctly to the database.

For the external-writer regime a **strict** cache validates through the authoritative database
or bypasses the cache. Where that is the only way to keep the guarantee, the result records
`cache_bypassed` and the analysis says so plainly rather than weakening "strict".

No-cache baselines are run for **both** models, both version modes, so the version/outbox
overhead is visible before any cached comparison is made.

### 2.3 Concurrency strategies (owner requirement 4)

The owned model's write path is implemented twice and measured as a controlled pair on the
no-cache baseline and on one cached scenario:

- **optimistic** — read `cache_version` (or the donation's amount) and write back only if it has
  not moved; bounded retries within a deadline; conflicts counted;
- **pessimistic** — `SELECT … FOR UPDATE` on `person` before reading or writing; waiters block;
  lock wait time recorded.

Both are run under the same business invariant, dataset and offered work, with an invalid
negative control (`ctl-unchecked-rmw`) that must lose updates.

### 2.4 Scenario registry — explicit dimensions, not duplicated SQL

Scenario id: `<model>-<ver>-<backend>-<strategy>-<freshness>-<writers>`.

| Dimension | Values |
|---|---|
| model | `owned`, `legacy` |
| ver | `opt`, `pess` (owned); `na` (legacy) |
| backend | `none`, `memory`, `redis` |
| strategy | `none` (backend `none` only), `aside`, `through` |
| freshness | `none` (backend `none` only), `relaxed`, `strict` |
| writers | `coord`, `ext20` (legacy only, 20 % of harness writes bypass the adapter) |

**Core matrix on PostgreSQL** — 20 cells:

| # | Scenario |
|---|---|
| 1 | `owned-opt-none-none-none-coord` |
| 2 | `owned-pess-none-none-none-coord` |
| 3 | `legacy-na-none-none-none-coord` |
| 4–5 | `owned-{opt,pess}-none-none-none-coord` are 1–2; cached block below doubles over model × backend × strategy × freshness |
| 4–19 | `{owned,legacy}` × `{memory,redis}` × `{aside,through}` × `{relaxed,strict}` (16 cells; owned uses `opt`) |
| 20 | `legacy-na-redis-aside-relaxed-ext20` |

plus `legacy-na-redis-aside-strict-ext20` (21) as the external-writer strict cell that must
report either authoritative validation or an explicit bypass.

**Database-reference set** (no cache, remeasured in this run, four cells): `ref-normalized-indexed`
(D2-like), `ref-flattened-fk` (D3, the principal baseline), `ref-rollup-trigger`
(database-maintained count/total), `ref-embedded-locked` (D10-style bounded embedding at 20).

**Controls** (measured like a design; never reported as a fast result):

- `ctl-stale-invalidation` — commits a database mutation and suppresses **exactly one**
  post-commit invalidation; the oracle must detect ≥1 stale-after-ack read. Fires only for
  `relaxed`; a `strict` variant must detect zero and is a stop condition if it does not.
- `ctl-dirty-write-audit` — dev-check only: injects a cache record whose version corresponds to
  a rolled-back or nonexistent state; the audit must reject it. Not a benchmark strategy.
- `ctl-unchecked-rmw` — the lost-update control for the concurrency pair.

Every report retains the complete scenario id and every recorded dimension.

### 2.5 Cache policy — required in every actual cache scenario

**Capacity and LRU.**
- memory: byte-bounded, concurrency-safe **exact** LRU (container/list + map, one mutex).
- Redis: recorded `maxmemory`, `maxmemory-policy=allkeys-lru`, `maxmemory-samples=5`;
  artefacts state that Redis's LRU is **approximate**.
- cache memory and CPU are charged to the resource budget; an in-process cache is not free.
- primary capacity is **smaller than the logical working set** so eviction occurs, plus a
  fits-working-set control.
- recorded: logical payload bytes, metadata bytes, resident bytes (memory backend), eviction
  count, evicted keys, Redis `INFO` memory/CPU/net and `commandstats`.

**Hard expiration.** Every entry has an exact five-minute hard TTL, `TTL = 300 s`. The primary
TTL is never shortened to make the experiment finish faster.

**Probabilistic early expiration.** On every otherwise valid hit, before returning it:

```
time_remaining = max(0, hard_expiry - now)
p_expire       = clamp(1 - time_remaining / 300 s, 0, 1)
```

draw from the trial's recorded random source; if the draw fires, treat the entry as expired and
refresh through the lease path. Same semantics for Redis and memory; sufficient timing metadata
is stored with the value. Unit tests use a **fake clock and deterministic RNG** at ages 0, 75,
150, 225 and 300 s (p = 0, 0.25, 0.5, 0.75, 1.0), and the implementation additionally reports
observed expiration rates by age bucket.

**Lease and stampede prevention.** Every fill or refresh uses a per-key lease.
- Redis: unique token, `SET key token NX PX <ttl>`, release by compare-and-delete (Lua).
- memory: per-key coordination with bounded waiting and cleanup.
- An expired lease holder must not overwrite a newer committed version (version fencing on
  publish, which the owned model can prove from `cache_version`).
- Waiters use bounded jittered backoff and then fall back to an authoritative database read.
- Recorded: lease acquisitions, contention/waits, wait duration, lease timeouts/losses,
  duplicate fills, backend database loads per miss burst, fallback reads.
- The lease duration is **calibrated from measured database-fill latency** and the reasoning is
  recorded; no arbitrary large constant is chosen silently.

### 2.6 Exact wrong-read accounting

No timed cache hit is validated by a hidden synchronous full database read — that would destroy
the performance question. Instead:

- an independent **Go oracle** and **operation ledger**: every acknowledged authoritative write
  records key, committed logical version, expected payload hash, commit/ack event and timing;
- at read start, the latest write to that key acknowledged before the read began is captured;
- every response records key, required version/hash, returned version/hash, cache source, cache
  age, and whether a write overlapped it;
- after the phase the ledger is **replayed** against the final authoritative database state.

For legacy rows the harness assigns an oracle version that is not stored in the database; it is
measurement metadata, not a schema modification.

Reported at least: total reads; fresh reads; wrong/stale reads; wrong reads as a percentage of
all reads and of cache hits; unique keys affected; versions behind; staleness duration
p50/p95/p99/max; longest consecutive wrong-read streak; time to freshness after an acknowledged
write; concurrent/ambiguous reads; impossible/uncommitted cache values; and wrong reads grouped
by cause (invalidation race, refill race, cache update failure, local-instance incoherence,
external writer, cache outage, other evidenced cause).

**Strict contract:** a read that begins after a write to that key was acknowledged must not
return an older value. Strict cells with a wrong read are correctness failures and cannot
support performance conclusions. Relaxed cells may return stale committed data, and their
throughput is reportable **only with the wrong-read count and rate beside it** — never shown
alone. Zero impossible/uncommitted cache values is required of every candidate cell.

### 2.7 Negative controls and deterministic fault injection

Faults are reproducible from a recorded seed/schedule, driven by `-fault-seed` and a phase name:

| Fault | Injected where |
|---|---|
| database rollback after pre-write tombstone | strict write-through path, forced |
| cache update/invalidation failure after database commit | relaxed and strict write-through |
| lease-holder death / timeout | fill path, lease abandoned |
| Redis restart / unavailability | Redis backend, connection refused window |
| cache-aside read-fill racing a writer | fill path concurrent with a write |
| legacy external writer bypassing the adapter | `ext20` regime |

### 2.8 Workloads and deployment regimes

Deterministic, fixed-seed, skewed donor access with hot donors and a bounded long tail. Measured:

1. cold-cache fill; 2. warm read-only throughput; 3. mixed 99 % reads / 1 % writes; 4. mixed
90 % / 10 %; 5. controlled hotspot contention; 6. synchronized cold/soft-expiry stampede;
7. sustained churn crossing **at least two real five-minute TTL boundaries**;
8. one logical application instance; 9. three logical application instances with equal total
workers — separate in-memory LRUs vs shared Redis.

The multi-instance test must expose that a local in-memory lease and invalidation do not
coordinate other instances. For owned strict local caches, the database version is validated or
an authoritative fallback is used. For legacy strict local caches with external writers an
authoritative read is performed or the guarantee is **marked unsupported**. An uncoordinated
local cache is never labelled strict.

Both isolated cacheable reads and a blended application workload containing cacheable and
uncacheable operations are run; **cacheable-endpoint throughput and total application throughput
are reported separately**. The open-loop scheduler (`measure.RunOpenLoop`) is used for
offered-load/saturation work, recording offered, started, completed, rejected, queue depth,
scheduling lag and client saturation. Closed-loop measurements are retained for controlled
latency comparisons only and are never presented as production SLO evidence.

### 2.9 Engines, placement, resources

Engines: pinned PostgreSQL single-node, YugabyteDB single-node, YugabyteDB three-node RF=3
(`infra/*`, images in `infra/versions.env`). Full core matrix first on PostgreSQL; then
**selected conclusion-bearing controlled pairs** on YugabyteDB rather than every dimension.

For three-node YugabyteDB: connections balanced across all query endpoints, endpoint
distribution recorded, colocated vs non-colocated data scenarios measured with engine evidence,
`EXPLAIN (ANALYZE, DIST)` and RPC counts captured. Containers on one WSL machine are never
claimed to prove colocation or network behaviour.

Redis is pinned to `docker.io/library/redis:7.4.11-alpine` (image id recorded per run); never
`latest`, never a host-installed Redis.

**Resource framings — calculated before measurement, never pooled.** Live host resources are
re-verified in phase A (studies 01–04 measure 8 CPUs / ≈15.36 GiB on `host-zenbook-ux5406sa`).
The host has 8 CPUs; the container VM's live totals are recorded, not copied.

| Condition | database | Redis | client | Σ CPU |
|---|---|---|---|---|
| `db-only` (baseline) | 2 CPU / 3 GiB | — | 2 CPU / 2 GiB | 4 |
| `add-cache` | 2 CPU / 3 GiB | 1 CPU / 768 MiB | 2 CPU / 2 GiB | 5 |
| `equal-total` | 1 CPU / 2.25 GiB | 1 CPU / 0.75 GiB | 2 CPU / 2 GiB | 4 |

- `equal-total` keeps the **service** total at the `db-only` service total (2 CPU / 3 GiB), with
  Redis's share deducted from the database's.
- For the memory backend, `equal-total` charges the inclusive cache's resident bytes against the
  client's 2 GiB, and the harness records the measured resident bytes so the deduction is
  auditable.
- Client sizing: readers = `-conns` (default 8), writer goroutines from the blend, database pool
  allowance = readers + writers + 4 control connections, Redis connections = readers + writers +
  2 control. The generator's ability to supply the intended demand is checked from the open-loop
  result's `dropped`/`client_saturated`, never assumed.
- Per-container CFS throttling and memory are captured for the client, Redis and every database
  node.

### 2.10 Metrics and plans

Recorded: completed throughput and offered load; p50/p90/p95/p99/max (deeper percentiles only
when sample counts support them); cache hit/miss/fill/eviction/expiration counts; hard vs
probabilistic expiry; database queries and bytes avoided; write throughput and acknowledgement
latency; cache maintenance latency; Redis `INFO` stats, `commandstats`, keyspace hits/misses,
`expired_keys`, `evicted_keys`, memory, CPU and network bytes; in-memory allocations, GC,
resident bytes, lock contention and eviction counts; database CPU/memory/throttling and storage;
the lease and wrong-read metrics of §2.6.

Plans are captured for the no-cache source query, the cache-miss/fill queries, owned strict
version validation and legacy strict authoritative validation. A cache hit has no SQL plan, and
the report says the saved database operation rather than inventing one.

Every important comparison gets independent repeated trials, a median, the individual values and
a spread. No conclusion rests on a single survey trial or on a difference below the measured
noise floor.

### 2.11 Scale and runtime bounds

| Scale | Donors | Donations/donor | Payload | Purpose |
|---|---|---|---|---|
| `tiny` | 40 | 8 | small | dev checks only; numbers never reported |
| `small` | 800 | 25 | ~2–4 KiB | the reported primary survey |
| `medium` | 3000 | 40 | ~2–4 KiB | churn, instances, repeated trials |

Primary cache capacity 256 KiB (below the `small` working set of ≈2–3 MiB) so eviction occurs;
control capacity 8 MiB (fits). Redis primary `maxmemory 32mb`, control `256mb`.

Routine choices the implementer may take without escalation: flag names, internal package
structure inside the study, exact metric field names, the number of decimals in the report,
which of two equivalent SQL formulations to use, and the lease-calibration formula (`4 ×
measured p99 of the database fill`, clamped to `[100 ms, 2 s]`, recorded).

**Reduction rules, fixed before seeing results** (used only if the projected runtime is
excessive, in this order):

1. retain both database-ownership models; 2. retain no-cache baselines; 3. retain Redis and
memory; 4. retain cache-aside and write-through; 5. retain relaxed and strict in the primary
PostgreSQL matrix; 6. retain controls; 7. retain the five-minute TTL; 8. reduce topologies or
secondary regimes before dropping a scientific pair; 9. record every omitted cell as a coverage
gap; 10. never reduce trials for a number used in a conclusion.

---

## 3. Coverage mapping for the owner's prospective requirements

| Requirement | Mapped measurement | Status |
|---|---|---|
| Calculated node/client sizing, per-host accounting, equal-total control, measured endpoint/client bottlenecks | §2.9 table; live re-verification in phase A; `equal-total` vs `add-cache` conditions in phase H; client saturation from the open-loop result | planned |
| Rollup | `ref-rollup-trigger` (parent count/total) vs `ref-flattened-fk` | planned |
| Rolldown | `ref-flattened-fk` (`donation.charity_id` copied) vs `ref-normalized-indexed` | planned |
| Embedding | `ref-embedded-locked` (bounded newest-20 embedding, D10-style) vs `ref-flattened-fk`; mutation cost and drift audit | planned |
| Colocated vs non-colocated data, multi-node | `y1` colocated vs `y2` non-colocated on `yb-cluster3`, with placement evidence and RPC counts | planned; honest gap if the cluster cannot be started inside the run budget |
| Optimistic and pessimistic concurrency | `owned-opt-*` vs `owned-pess-*`, plus `ctl-unchecked-rmw` invalid control | planned |
| Every measured pair's independence | one decision per pair; pairs named in the report | planned |

---

## 4. Ordered execution

Commands are path-safe for WSL + `podman.exe`; every runner takes the shared benchmark lock
(`run_lock_acquire`). Working directory is the task worktree.

| Step | Exact action or command | Expected evidence | Failure action |
|---|---|---|---|
| A | `source infra/lib.sh; need_podman; podman volume inspect ads-run-lock`; free-run CPU/memory; `./probe-cache.sh` (Redis capability probe: `SET NX PX`, Lua compare-and-delete, `INFO`, `maxmemory-policy`, TTL) | engine reachable, lock free, host resources recorded, Redis policy + image id recorded | mapped wait for the lock; ER if Redis primitive behaves differently |
| B | `go vet ./... && go test ./...` inside the benchmark image build | LRU byte bound, TTL, probabilistic expiry at 0/75/150/225/300, lease, fencing, oracle, digest tests pass | fix; a failing test is never skipped |
| C | `./run-study.sh --scale tiny --topologies pg-single --phases verify,explain,warm,mixed,hotspot,faults,stampede,instances,audit --designs <all>` | every scenario loads, gate passes, both controls fire, dirty-write audit rejects, faults are reproducible | ER if a control cannot fire after the workload is made harsher |
| D | `./run-study.sh --tag --scale small --topologies pg-single` (core matrix + reference designs) | one JSON per cell, plans captured, report generated, digest recorded | failed cell recorded in the manifest; matrix continues |
| E | repeated relaxed-vs-strict wrong-read/churn comparison, 3 trials, `-scale medium` | median + individual trials + spread; wrong-read rates per cell | ER if a `strict` cell records a wrong read |
| F | stampede and lease-recovery runs; lease-death and Redis-outage faults | duplicate fills, fallback reads, recovery time | ER if the lease cannot be shown to prevent duplicate fills |
| G | one-instance vs three-instance runs (equal total workers) | local-LRU incoherence exposed for memory; Redis shared | ER if a strict local cache cannot be shown uncoordinated |
| H | selected `yb-single` pairs; `yb-cluster3` colocation pair with endpoint distribution and RPC counts; `add-cache` vs `equal-total` resource framings | engine evidence, per-node throttling, endpoint distribution | recorded coverage gap, never a silent omission |
| I | repeated conclusion runs (`--trials 3`) on the pairs the analysis cites | medians, individual values, spreads | ER if spread > 20 % on a cited number |
| J | report regeneration, signed analysis, discussion companion, artifact validation, integration into local `main` | digest stable; `main` contains the task; lock released; containers stopped | blocked integration reported as blocked |

**Estimated cells and duration.** Core matrix 21 + 4 reference + 3 controls ≈ 28 cells on
`pg-single` at `small` (≈40 s each) ≈ 20 min; dev checks ≈ 15 min; repeated trials ≈ 10 min;
churn (two real TTL boundaries) 2 cells × ≈11 min ≈ 22 min; stampede/lease/faults ≈ 10 min;
instances 4 cells ≈ 5 min; `yb-single` 4 cells + startup ≈ 12 min; `yb-cluster3` 2 cells +
startup ≈ 20 min; reporting and analysis ≈ 30 min. **Projected total ≈ 2.5 h** of machine time.
An escalation is raised if any single step is projected above 45 min or if the total exceeds
4.5 h — the reduction rules in §2.11 then apply in the stated order.

Stop conditions: a `strict` cell with a wrong read; any impossible/uncommitted cache value; a
correctness gate failure; a negative control that cannot be made to fire; a dirty-write audit
that fails to reject the injected record.

---

## 5. Escalation Required

Triggers (record in `ESCALATIONS.md`, stop only the dependent step, continue independent mapped
work):

1. A cache primitive behaves differently from §2.5 in a way that changes what a scenario means.
2. The probabilistic-expiry formula of §2.5 cannot be implemented exactly as written.
3. A `strict` scenario cannot be given a correctness proof for the legacy external-writer regime
   within the declared contract.
4. A negative control does not fire after the workload has been made harsher.
5. A requested topology cannot be started or a placement comparison cannot be evidenced.
6. Projected runtime exceeds the §4 bounds and the reduction rules would drop a scientific pair.
7. Anything that would require changing `platform/`, `infra/`, another study, or a shared
   document beyond an additive, study-scoped change.

The executor does not rewrite this handoff; HIGH appends a dated `AM-NN` amendment below and
tags it immutably.

---

## 6. Acceptance, integration and next role

**Implementation milestone:** unit tests and `go vet` pass in Podman; the complete scenario
registry exists; every candidate scenario uses LRU, the 300 s TTL, the specified probabilistic
expiry and leases; the stale-read negative control fires; the dirty-cache-write audit rejects
rolled-back/nonexistent versions; every candidate reports zero impossible/uncommitted values.

**A measured strict cell is valid only if** it has zero stale-after-ack reads, zero impossible
cache values, a passing authoritative database reconciliation, passing source-write invariants,
and passing lease and failure-recovery audits.

**A completed study requires:** clean committed producing code before the reported runs;
immutable annotated `run/05-cache-consistency/<run-id>` tags; environment, commit, `git
describe`, dirty flag, image identities, seeds and resource conditions in every result; a stable
inputs digest; a generated report, signed analysis and discussion companion; `CONTEXT.md` and
`LESSONS_LEARNED.md` updated; the latest local `main` integrated into the task branch; combined
checks rerun; local `main` fast-forwarded or otherwise integrated per `AGENTS.md`; the task
commit proved reachable from `main`; an integrated milestone tag; the benchmark lock released
and study containers stopped; and **no push and no pull**.

Tags, checked for collisions first, annotated only:
`study-05/v0-handoff`, `study-05/v0.N-handoff-amendment-NN`, `study-05/v1-harness`,
`run/05-cache-consistency/<run-id>`, `study-05/v1-measured`, `study-05/v1-analysis`,
`study-05/v1-integrated`.

Branch-back into `main`: commit explicit owned paths only (never `git add -A`), identify
DeepSeek HIGH in the commit messages, integrate the latest local `main`, reconcile the shared
living documents as one project, then fast-forward `main`.

**Next after execution: HIGH — DeepSeek HIGH reviews validation and results, then writes the
signed analysis.** This task has no separate model switch; the review and the analysis are the
same agent's later pass, and the analysis states that absence of independent model review as a
limitation.

---

## 7. Allowed additive changes outside the study directory

- `studies/05-cache-consistency/.gitattributes` — the study's scoped LF policy.
- `infra/versions.env` — **append** the pinned Redis image (`REDIS_IMAGE`). No existing value
  changes.
- `infra/redis.sh` — a new topology script following `pg-single.sh`'s shape, taking
  `run_lock_guard`. Nothing existing in `infra/` changes.
- `platform/` — **no change is expected.** If one becomes necessary, it carries its own tests and
  a `repo/*` tag, and an ER is raised first.
- `CONTEXT.md`, `LESSONS_LEARNED.md`, `README.md` (root) — updated at integration time only.

---

## Amendments

*(HIGH appends dated, attributed `AM-NN` decisions here; existing text is never erased.)*

### AM-01 — 2026-09-21 — DeepSeek HIGH (`deepseek-flash`)

**Status: decided by HIGH, implemented, measured and integrated. EH-05 rev 1 stands; this
amendment records three decisions the execution forced and one requirement the execution
CHANGED. The original text above is untouched.**

#### 1. Requirement change: an invalidation fence is required, not optional

EH-05 §"cache-aside sequence" and §"write-through sequence" required that a cache value
be published only if it corresponds to a successfully committed state, and that a strict
writer "fences/tombstones before the mutation". **That is not sufficient and the handoff
was wrong to treat it as sufficient.** Measured: a strict cell with only that rule
recorded ~14 000 stale-after-ack reads, because a reader whose fill began *before* the
invalidation holds a committed but SUPERSEDED state and can republish it afterwards.

Amendment to the design, now binding for this study and recommended for any later one:

* `EntryStore` carries a per-key **invalidation fence**. `Fence(key)` advances it and
  removes the value in ONE atomic operation (one mutex in the in-process store, one Lua
  script in Redis); `FenceOf(key)` reads it; `Put` is refused unless the stored fence
  equals the entry's.
* A fill captures the fence **before** it takes its database snapshot and stamps it on
  the entry.
* A strict writer fences **before AND after** the commit. One fence leaves a window in
  which a reader captures the new fence and then reads the pre-mutation state.
* A relaxed writer's post-commit invalidation is the same operation; what makes a
  scenario relaxed is WHEN it runs and that it may be skipped or fail, not a weaker fence.
* The ledger's freshness requirement moves at the **ACKNOWLEDGEMENT**, not at the commit,
  because the contract is written against the ack and a strict writer fences before it
  acks. A committed-but-unacknowledged state is classified `ahead`, never `impossible`.

Evidence: `PROGRESS.md` (the dev-check iterations), the signed analysis
`reports/analyses/20260921-cache-consistency-survey.md` §2, and the discussion companion
`reports/discussions/20260921-cache-fences-and-ledger.md` §1.

#### 2. Decision: Redis is a new shared topology script, pinned in `infra/versions.env`

`infra/redis.sh` starts `docker.io/library/redis:7.4.11-alpine`
(id `f84b0c4678011602b9b98c227a4dcd5468bf8b088b02fdd4165cb7758bad8058`) with an explicit
`maxmemory`, `maxmemory-policy=allkeys-lru`, `maxmemory-samples=5`, `--save ''` and
`--appendonly no`, under the benchmark lock. Append-only to `versions.env`; no existing
value changed. Studies 01–04 never start Redis, so the addition cannot affect them.
Escalation: `05-ER-02`, decided.

#### 3. Decision: no WSL-local Podman engine; the documented Windows engine is used

`podman` is not on `PATH` in WSL. `infra/lib.sh`'s resolver reaches the documented Windows
engine through `podman.exe`, which is exactly the bridge the owner added, so **no new
bridge and no new engine were created** (creating one would have violated the
"install nothing on the host" rule and risked a second, unrunnable engine).
Escalation: `05-ER-01`, decided.

#### 3a. Finding after the amendment: the ledger's requirement can be AHEAD of the database

A follow-up instrumented run (`results/diag-strict/`) recorded the first few stale reads of the
three strict `through` cells in full. Two things are now evidenced rather than suspected: no stale
read overlapped a write, and some of them came from **authoritative database reads** returning a
state 1–2 versions behind the requirement. A bypass read consults no cache, so it cannot be a cache
fault: the ledger's freshness requirement was ahead of the database's own state.

**Binding consequence for the next iteration:** before any strict cell is judged, assert that
`Required()` is never ahead of the content the database currently holds, and fail the cell if it is.
The strict `aside` results and the single-instance relaxed rates are unaffected (aside cells are
clean; the relaxed rates come from cells that passed every gate).

**AM-02 — outcome, 2026-09-21 (same session).** The assertion is implemented and the defect is
fixed. Cause: a mutation selected a child row from the ledger and then acted on it BY ID ONLY, so a
concurrent writer could move or delete that row in between and the database and ledger diverged.
Fixes: (i) every child-row statement enforces the observed owner in all five schema directories;
(ii) a statement matching no row is `errNoEffect` — the database did not change, so the ledger does
not change, the write is an acknowledged no-op and is counted; (iii) corrections are relative in
every schema, including the reference designs, where an absolute SET had been writing the delta as
the amount. The assertion now runs after every writing phase: **0 mismatches over 800 donors in
warm, mixed, hotspot and stampede**, and 1 donor in one cell after `instances`, which is the next
thing to chase. A second measurement defect was closed with it: violations recorded by the
instances-phase arms were invisible to the cell's acceptance check (an owned redis through cell had
passed with 11 hidden stale reads). With both fixed, the strict `through` failures are HONEST and
real — 7 to 32 stale reads with every ledger assertion passing and the fence refusing 260–625
publications per cell — so the next iteration should treat them as a cache-design question in the
through path, not as an accounting one.

#### 4. Decision: the ledger cannot model concurrently conflicting writes of one key

Recorded as a **known limitation, not fixed**: two concurrent mutations of the same key
can reach the ledger in an order the database did not commit in, which surfaces as an
impossible value or an audit mismatch. One instance of it *was* fixed properly
(amount corrections are relative in both the SQL and the ledger); the remaining instance
(set-valued mutations) is declared in the analysis' weakness list and in section 6 of the
analysis as a harness limitation, with the two honest routes out of it named. HIGH did not
choose between them because both change what the hotspot phase measures and the choice
belongs with the next planning iteration.
