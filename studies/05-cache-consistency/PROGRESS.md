# Progress — Study 05, external cache throughput and consistency

Observed facts and routine decisions, in order. This log records what happened; it is not the
analysis. Decisions the handoff did not map live in `ESCALATIONS.md`; the specification lives in
`HANDOFF.md`.

**Agent:** DeepSeek HIGH selected by user, one agent for every phase (owner's instruction for
this task). Model id exposed by the session: `deepseek-flash`. Effort setting: **not exposed**.
Tool identity: **not exposed**.

---

## Step 0 — WSL safety preflight (read-only), 2026-09-21

| Check | Result |
|---|---|
| `pwd -P` | `/mnt/c/extra/code/architecture_design_studies` — matches the required root |
| `git status --short --branch` in the main checkout | `main`, working tree clean; `main` = `e8e4009c664b5d1a7d2a99688acd6b3a46caa361` |
| `git worktree list --porcelain` | `main`; `.worktrees/study04-configuration-portal` (`study-04/configuration-portal`, clean, at the same commit); `.worktrees/recency-reports` and `C:/Users/xiles/.codex/worktrees/study05-cache-prompt/...` — both marked **prunable**, both **left untouched** (not repaired, not removed, not adopted) |
| Existing `study-05` branch | none — `git branch -a --list '*05*'` is empty |
| Existing `study-05` tag | none |
| Benchmark lock | `podman volume inspect ads-run-lock` → no such volume. **Free.** |
| Containers | one exited `pg-scratch`; no benchmark or database container running |
| Images present | `localhost/configbench:1` (study 04), `postgres:17.11`, many dangling layers |
| Podman client / server | `podman.exe` 5.8.1 (windows/amd64), reached from WSL through `infra/lib.sh`'s resolver exactly as study 04's AM-01 decided; **no WSL-local engine was initialised** |
| Network | `proxy.golang.org` reachable (HTTP 200); `registry-1.docker.io` reachable (HTTP 401, the expected unauthenticated answer) — the image build can run `go mod download` |
| Redis pin | `docker.io/library/redis:7.4.11-alpine` pulled; image id `f84b0c4678011602b9b98c227a4dcd5468bf8b088b02fdd4165cb7758bad8058` |

Nothing was pushed, pulled, rebased, amended or reset. No tag was moved or deleted. No foreign
worktree, branch, file or container was touched.

## Step 1 — Task worktree created

```
git worktree add .worktrees/study05-cache-consistency -b study-05/cache-consistency
```

From local `main` = `e8e4009c664b5d1a7d2a99688acd6b3a46caa361`. Verified afterwards: the worktree
is **clean** (`git status --short --branch` shows only the branch line) and at that commit.

## Step 2 — Planning checkpoint

Written and committed as `study-05/v0-handoff`:

- `HANDOFF.md` — `EH-05` revision 1: objective, terminology, the two database models, the
  scenario registry, the cache policy (LRU, 300 s TTL, the exact probabilistic-expiry formula,
  leases), the oracle and wrong-read accounting, controls and faults, workloads, engines and the
  two resource framings, the fixed reduction rules, the coverage mapping and the ordered
  execution table with commands, expected evidence and stop conditions.
- `ESCALATIONS.md` — escalation log for this task.
- `PROGRESS.md` — this log.
- `README.md` — the question and the designs, for a reader.
- `study.env` — study image name and engine flags.
- `.gitattributes` — the study's scoped LF policy, so a run's inputs digest survives a checkout.
- `infra/versions.env` — **appended** `REDIS_IMAGE` (pinned); no existing value changed.

*(Steps 3 onward are appended as they complete.)*

## Step 3 — Phase A: Redis capability and policy probe

`./probe-cache.sh` built the image and ran the harness's `-cmd probe` against a real
Redis. **Passed**, evidence in `results/devchecks/phase-a-redis/`:

| Primitive | Observed |
|---|---|
| `PING` | `PONG` |
| `SET NX PX` lease, first holder | acquired |
| `SET NX PX` second holder | refused |
| compare-and-delete release with the wrong token | value and lease untouched |
| compare-and-delete release with the right token | released, re-acquirable |
| older publication (lower seq) | refused |
| newer publication | accepted |
| round-trip content hash | identical |
| invalidation fence: publish from before the fence | refused |
| invalidation fence: publish from after the fence | accepted |
| 1 ms hard TTL | gone when read 30 ms later |
| `maxmemory` / policy / samples | 33554432 / `allkeys-lru` / 5 (read back from `INFO`) |

Redis image id `f84b0c4678011602b9b98c227a4dcd5468bf8b088b02fdd4165cb7758bad8058`.

## Step 4 — Dev checks: five iterations, each one a harness bug the gate caught

The `tiny` matrix on `pg-single` (all scenarios, all phases except churn) was run five
times. Every iteration failed in a way that would have produced a confidently wrong
record rather than an error, which is the failure mode this repository is built
around. In order:

1. **Forcing a mutation onto a hot donor after building it.** `runHotspot` and three
   fault phases built a mutation and then overwrote `m.PersonID`. For a
   correct/delete/reassign the donation belongs to the donor chosen inside the
   builder, so the database changed a row the ledger attributed to someone else. The
   gate caught it as an **impossible cache value** and failed a fault cell whose own
   logic was correct. Fixed by making the target a parameter of the builder
   (`buildMutationFor`).
2. **The cache under test was never started.** The runner called `need_redis` for a
   helper named `needs_redis`, so no Redis container existed — and because a cache
   outage degrades to authoritative reads by design, every redis scenario "passed"
   its phases while measuring nothing but the database. Only the fault that demands a
   working lease failed, and it looked like a lease bug. Fixed twice over: the typo,
   and a `PING` gate in `runCell` that refuses to measure against a cache that does
   not answer.
3. **The instances phase left the base cache stale.** It measures a different
   deployment through its own instances, whose writes cannot invalidate the base
   adapter's stores, so pre-phase entries survived it. It now starts and ends cold.
4. **The YB-only ids did not match the registry**, so the placement pair was
   scheduled on PostgreSQL and aborted.
5. **Cache-aside is not made safe by "publish only after a commit".** A reader whose
   fill began before a writer's invalidation holds a committed but superseded state
   and could republish it: ~14 000 stale-after-ack reads in a strict cell under
   coordinated writers, and thousands in an owned cell. Fixed with a cache-side
   invalidation FENCE (`Fence`/`FenceOf`, and `Put` refusing a publication whose
   fence has moved), captured before the database snapshot.
6. **The fence alone was not enough.** Fencing only *before* the commit leaves a
   window in which a reader captures the new fence and then reads the pre-mutation
   state, so it may publish a superseded value at a "current" fence. Strict writers
   now fence **before and after** the commit.
7. **An unrecorded state is not an impossible state.** The oracle learns of a
   committed state just after the commit returns, so a reader racing that commit can
   hold a state the ledger has not recorded. `finish` now settles an otherwise
   unknown hash with ONE conditional database read and classifies it `ahead`; it
   never calls an unrecorded committed state impossible, and the check runs only on
   that path, never on a hit.

Two smaller defects were fixed alongside: `readLog.Summary` returned its live maps
(so every phase printed the cell's final totals beside its own counters), and the
per-phase impossible count was summed once per phase, inflating the cell total.
