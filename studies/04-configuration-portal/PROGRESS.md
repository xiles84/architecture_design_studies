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

## Step 4–7 — design catalogue, harness, image and dev checks

Added and committed:

- 10 SQL designs under `sql/`, five files each, in the repository catalogue format.
- The Go harness (`harness/`): deterministic dataset, load, gate, workload, contention,
  ledger, audits, report.
- `platform/core/measure/arrival.go` (already committed separately) plus two new tests.
- `Containerfile`, `study.env`, `sqlfs.go`, `run-study.sh`, `probe-engines.sh`, `README.md`,
  `sql/README.md`.

**Phase C dev checks — all 10 designs on `pg-single` at `tiny`, gate passed, 5 reads and 6 writes
measured per design, report generated, inputs digest `364ee7176b19777c`.**

Both negative controls fired:

| Design | Strategy | Acked | Counter | Expected | Lost updates |
|---|---|---|---|---|---|
| `c1_optimistic_version` | optimistic | 286 | 286 | 286 | **0** |
| `c2_pessimistic_lock` | pessimistic | 497 | 497 | 497 | **0** |
| `x1_lost_update_control` | none (control) | 653 | 46 | 653 | **607** |

`x2_rollup_drift_control` fired through its own `a_rollup_mismatches` audit.

## Bugs the gate caught before a number was reported

Five harness bugs, four of them found by the correctness gate rather than by reading the code.
All five are recorded because they are the failure mode this repository is built to avoid: each
one would have produced a confident wrong number.

1. **Positional comparison of an alphabetically sorted result against a section-ordered dataset.**
   The reference design was reported as returning the wrong values. Fixed by sorting the expected
   entries by key.
2. **Two statements used `$1` without declaring `-- params:`.** 2 263 write operations ran
   against a server answering "there is no parameter $1" and every one was counted as an error,
   so the read side looked clean while the write side measured nothing. Fixed in the SQL, and the
   harness now refuses to load a design whose statement uses `$n` without declaring parameters.
3. **One ledger per phase instead of per cell.** The contention audit then reconciled the
   database against a dataset the write phase had already legitimately changed, and reported every
   installation as missing keys.
4. **A negative control failing its audit was treated as a broken design.** x2 firing is the
   control working; failing the cell there would have made the control look like a defect.
5. **The lost-update detector compared the stored counter against the last value written** —
   a number with itself — and reported 583 acknowledged increments with a counter of 39 as a pass.
   It now compares the counter against the number of acknowledged increments.

Two more were found while building the probe and the runner:

6. **The podman build context is resolved client-side on Windows**, so `/mnt/c/...` becomes
   `C:\mnt\c\...` and the build fails with "context must be a directory". A bind-mount source
   tolerates `/mnt/c/...`; a build context does not. The runner now uses `winpath()` for the
   build and `hostpath()` for mounts. The phase B probe had proved mounts, not the build.
7. **The probes empty-mount check grepped the cumulative report**, so one failing form tainted
   every later form. It now greps a per-form log.

## Known gaps at this checkpoint

- **Eight designs are mapped and not implemented** (`d2`, `d3`, `d4`, `h1`, `s1`, `s2`, `y1`,
  `y2`). Recorded as a coverage gap; no conclusion may treat them as measured or as excluded.
- **Diagrams are not written yet.** `diagrams/` is empty; the ASCII tables in `README.md` and
  `sql/README.md` carry the same information for now.
- Only `pg-single` has been run. `yb-single` and `yb-cluster3` are untested at this commit, and
  the colocation pair that needs them is not implemented.
---

*(Steps 8 onward are appended as they complete.)*

## Step 8–9 — small matrix, report, and signed analysis

**Phase D (reduced), `20260921T1215Z-small`.** Ten cells, one per design, `pg-single`, scale
`small` (60 installed products × 60 entries), phases `verify,explain,read,write,contention`,
`--duration 3s --warmup 1s`. Run tag `run/04-configuration-portal/20260921T1215Z-small` on commit
`5e15abb`; inputs digest `ab0f6ef5e5500775`; **no cell failed**; both negative controls fired.

| Design | Strategy | Acked | Counter | Lost updates | ops/s |
|---|---|---|---|---|---|
| `c1_optimistic_version` | optimistic | 2 065 | 2 065 | **0** | 503 |
| `c2_pessimistic_lock` | pessimistic | 3 539 | 3 539 | **0** | 886 |
| `x1_lost_update_control` | none (control) | 3 416 | 215 | **3 201** (93.7 %) | 768 |
| `x2_rollup_drift_control` | drift (control) | — | — | fired: 60 rollup mismatches | 570 |

Signed analysis:
`reports/analyses/20260921T1215Z-small--deepseek-flash--2026-09-21.md`, with the mechanism
companion `reports/discussions/20260921-configuration-portal-mechanisms--deepseek-flash--2026-09-21.md`.
The report was regenerated afterwards so its analysis index picks the analysis up by digest.

**The analysis's principal finding is a weakness in this run.** `r01`/`r02`/`r03` have
byte-identical SQL in the eight row-per-key designs, yet their throughputs spread by ~65 %
(`r03` 25 114 to 40 010 ops/s). The likely cause is that every cell runs against the same
container in one pass and each drops and recreates its schema over the previous cell's dead
tuples, so the database is not the same instrument at cell 1 and cell 10 — but this run contains
no repeated design and cannot prove it. The consequence is stated wherever the numbers are used:
**no read difference below ~1.7x in this digest is attributable**, and only the ~15x rollup effect
is claimed.

Coverage gaps recorded rather than implied away: eight of the eighteen designs are unbuilt;
YugabyteDB topologies, cadence, cardinality crossover, churn, deployment controls and repeated
conclusion runs are mapped and unrun; no independent model has reviewed this work.

## Step 10 — integration into local main

- Pre-merge checks: benchmark lock free; no benchmark or database container running; `main` still
  at `df13f2a`, the commit the task branch was cut from; task worktree clean.
- **EOL-only proof over the main checkout:** `git status --porcelain` reported 1970 modified files,
  and `git diff --ignore-cr-at-eol --quiet` exited **0** over all of them — every one differed from
  the index only by a CR at end of line. 0 staged, 0 untracked.
- **Normalization (owner-authorized):** `git diff --name-only -z | xargs -0 git restore --worktree`
  in the main checkout. Dirty count 1970 → **0**; `main` still `df13f2a`; `git ls-files --eol`
  now reports `i/lf w/lf`. No index entry, commit, tag or history was touched.
- **Merge:** `git merge-base --is-ancestor main study-04/configuration-portal` succeeded, so the
  integration is a fast-forward over six commits with no merge commit and no concurrent-agent
  reconciliation to perform. `main` was fast-forwarded; the integrated state is tagged
  `study-04/v1-integrated`.
- Nothing was pushed or pulled.
