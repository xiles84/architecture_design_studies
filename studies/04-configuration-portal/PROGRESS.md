# Progress — Study 04, configuration portal

Observed facts and routine decisions, in order. This log records what happened; it is not the
analysis. Decisions the handoff did not map live in `ESCALATIONS.md`; the specification lives in
`HANDOFF.md`.

**Agent:** DeepSeek HIGH, one agent for every phase (owner's instruction for this task). Model id
exposed by the session: `deepseek-flash`. Effort setting: not exposed. Tool identity: not exposed.

---

## Step 0 — WSL safety preflight (read-only), 2026-09-21

| Check | Result |
|---|---|
| `pwd -P` / `git rev-parse --show-toplevel` | `/mnt/c/extra/code/architecture_design_studies` — matches the required root |
| Shell / Git | `/usr/bin/bash`; `/usr/bin/git` 2.43.0 — native WSL Git throughout, no `git.exe` |
| `git status --short --branch` | `main` at `df13f2a`; **1970 modified files, 0 staged, 0 untracked** |
| Nature of that dirt | **100 % CRLF-vs-LF**: `git ls-files --eol` reports `i/lf w/crlf` for 1997 paths, `file(1)` confirms `with CRLF line terminators`, and `git diff --numstat` shows equal insert/delete counts. No content edit exists. |
| `core.autocrlf` / `core.filemode` | unset (nowhere) / `false` (`.git/config`) — no global git config was changed |
| `git worktree list --porcelain` | `main` + `.worktrees/recency-reports` (`repo/recency-and-reports`, marked `prunable`) |
| `.worktrees/recency-reports` | **Broken cross-platform worktree**: its `.git` holds `gitdir: C:/extra/...`, so native WSL Git reports `fatal: not a git repository`. Left untouched; no prune, no repair. Reported to the owner. |
| Benchmark lock | `podman.exe volume inspect ads-run-lock` → no such volume. **Free.** |
| Containers | one exited `pg-scratch`; no benchmark or database container running |
| Images present | `localhost/ticketbench:1` (study 02), `localhost/seatbench:1` (study 03), `golang:1.26-bookworm` |
| Podman client / server | `podman.exe` 5.8.1 (windows/amd64) → 5.8.5 (linux/amd64) on `podman-machine-default` |
| Environment | live 8 CPUs, 16 496 422 912 B, kernel `6.6.87.2-microsoft-standard-WSL2`, host ASUS Zenbook S 14 UX5406SA / Core Ultra 7 258V / 31.48 GB → `host-zenbook-ux5406sa` still applies |
| Tags | `git tag --list` read; `study-04/*` and `run/04-*` do not exist yet |

Forbidden operations were not run: no `git worktree prune`, `git gc --prune`, `git clean`,
`git reset --hard`, `git checkout -- <path>`; no `podman machine init`, `podman system reset`,
`podman volume rm ads-run-lock`, no `ADS_IGNORE_RUN_LOCK=1`.

## Step 1 — Task worktree created

```
git worktree add .worktrees/study04-configuration-portal -b study-04/configuration-portal
```

From local `main` = `df13f2a6baeded28bf5d1216b6ec6d622e259850e`. Verified afterwards: the new
worktree is **clean** (`git status --short --branch` shows only the branch line) and its files are
**LF** (`git ls-files --eol` → `i/lf w/lf`). This confirms the main checkout's CRLF state is local
to that checkout and does not propagate to a worktree — which is what makes the final integration
step tractable.

`.worktrees/` is excluded through `.git/info/exclude` (`/.worktrees/`), so the new worktree stays
invisible to the main checkout's `git status`.

## Step 2 — Planning checkpoint

Written and committed as `study-04/v0-handoff`:

- `HANDOFF.md` — `EH-04` revision 1, including AM-01 (ER-01 decision).
- `ESCALATIONS.md` — `04-ER-01` raised and decided.
- `PROGRESS.md` — this log.
- `.gitattributes` — the study's scoped LF policy.
- An initial `CONTEXT.md` entry.

`04-ER-01` (no native WSL podman) was raised before any container work and decided as AM-01:
drive the documented Windows engine through `podman.exe`, never initialise a WSL-local one.

---

## Step 3 — WSL → Podman bridge (`infra/lib.sh`), and the Phase B bind-mount probe

**Bridge (AM-01).** `infra/lib.sh` gained an additive engine resolver:

```
ADS_PODMAN = $ADS_PODMAN | $PODMAN | `podman` on PATH | /mnt/c/Program Files/RedHat/Podman/podman.exe
podman() { command "$ADS_PODMAN" "$@"; }
```

Verified live from the task worktree:

```
ADS_PODMAN=/mnt/c/Program Files/RedHat/Podman/podman.exe
command -v podman -> podman          # the function; every existing call site keeps working
need_podman OK
podman --version -> podman version 5.8.1
winpath <worktree>/infra -> C:/extra/.../.worktrees/study04-configuration-portal/infra
hostpath <worktree>/infra -> /mnt/c/extra/.../.worktrees/study04-configuration-portal/infra
```

`hostpath()` was **not** modified — studies 01–03 also pass its result to `git`, so widening its
meaning would have changed code behind published results. A separate `winpath()` was added as a
documented fallback. Study 01–03 runners are untouched.

**Phase B probe** (`probe-bind-mount.sh`, evidence in `results/devchecks/phase-b-bind-mount/`).
Took the benchmark lock, started one container (no database), and released the lock. Result:

| Form | Sentinel visible | Read-back identical | Host write survived | Verdict |
|---|---|---|---|---|
| `hostpath` → `/mnt/c/...` | yes | yes | yes | **WORKS** |
| `winpath` → `C:/...` | yes | yes | yes | **WORKS** |
| control → a path that cannot exist | no | n/a | no | **FAILS** (container exit 125) |

`RUNNER_DECISION=hostpath`. The runner therefore uses `hostpath()` exactly as studies 02 and 03 do.
The control matters: without it, a probe that always answered "works" would have looked identical.

Two bugs were found and fixed in the probe itself while building it, both recorded because they are
the kind that silently produce confident nonsense:

1. an `awk` expression inside the double-quoted container command was expanded by the *outer* shell
   (`$3`/`$4` unbound under `set -u`), so the probe died mid-run — the lock was released by its
   EXIT trap, and the failure was visible rather than silent;
2. the "empty mount" check grepped the **cumulative** report, so one failing form would have marked
   every later form as an empty mount. It now greps a per-form log, and the verdict window was
   widened past the diagnostics to the verdict line.

Probe runs are dated and kept (`report-*.txt`); the first is the crashed run above.

---

*(Steps 4 onward are appended as they complete.)*
