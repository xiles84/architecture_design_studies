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
