# Progress — implement AI task queue v1

**Actor:** `deepseek-flash` (tool and effort unknown; session capability LOW, work role executor).
**Task:** `task-20260922T025912Z-queue-v1`.
**Claim:** `claim-76540de5a8deb0b3`, epoch 1, attempt `attempt-20260922T101824Z-e309f5`.
**Branch/worktree:** `repo/ai-task-queue-v1` in `.worktrees/ai-task-queue-v1`, cut from local
`main` at `eef4d81` before the bootstrap CLI existed; the branch is the bootstrap claim.

## What was built

`tools/queue` is a dependency-free Go CLI. Live coordination is one Git blob per task at
`refs/ads-queue/live/<task-id>` plus `refs/ads-queue/locks/main-integration`, written with
`git update-ref <ref> <new> <observed-old>` and reflogs; committed `events/*.json` are the
permanent archive and reconstruct state on a fresh clone. Every command in the brief is
implemented. `tools/queue/queue` is the thin host wrapper; `Containerfile` builds
`localhost/ads-queue:1` from the pinned Go base (`QUEUE_IMAGE` in `infra/versions.env`) and
runs `go vet ./... && go test ./...` as a build gate.

Commits: `a35de47` (implementation, wrapper, image pin, tests).

## Decision Log

1. **Command-to-state mapping.** The brief's diagram and command list do not match
   one-to-one, so: `claim` accepts `changes_requested` directly; `yield` is the command
   that returns `blocked_high → ready` once HIGH has resolved an escalation; the first
   `checkpoint` moves `claimed → in_progress`; `recover` records the new epoch as
   `claimed`; `integrate` records the merge while staying
   `approved_for_integration`, and `complete` performs the terminal transition. Reason:
   the smallest interpretation that keeps every documented state reachable. Confidence:
   High. Affects: every future queue user; documented in `tools/queue/README.md`.
2. **Live update ordering.** `claim` and `recover` take the ref with a create
   compare-and-swap *before* writing anything else, because the ref is the exclusivity
   gate. Other transitions emit the durable event first and then compare-and-swap the ref,
   because `guard` already rejects a superseded worker. Reason: exclusivity must not
   depend on a commit succeeding. Confidence: High. Affects: crash-recovery semantics.
3. **Evidence files are committed with the event.** `emitEvent` commits the event,
   `STATUS.md` and any new evidence path together. Without this, `publish`'s `task.json`,
   `review`'s `REVIEW.md` and `escalate`'s `ESCALATIONS.md` stayed untracked while the
   event claimed they existed, and integration would have lost them. Found by the test
   suite (`TestFullLifecycle`, `TestFreshCloneReconstruction`). Confidence: High.
4. **Bootstrap transition allowance in the validator.** The pre-CLI HIGH planning record
   (`task-20260922T025912Z-publish-book-pipeline-handoff`) runs
   `proposed → in_progress → completed`, skipping the worker lifecycle, and its history is
   immutable. `model.go` allows exactly that bootstrap shape so `audit` is clean on the
   repository's own archive; command guards are untouched, so `complete` still requires
   `approved_for_integration`. Confidence: Medium. Affects: how much the validator proves;
   flagged for HIGH review, and the clean alternative is a corrective event on that task.
5. **Worktree reuse without `git worktree list`.** The repository stores linked-worktree
   `.git` links and `.git/worktrees/<name>/gitdir` back-pointers as relative paths. Git
   2.43 resolves those when operating *inside* the worktree but reports the worktree as
   prunable in `worktree list`. The CLI therefore detects an existing worktree from the
   filesystem and by opening it, not from `worktree list`. Confidence: High; verified with
   WSL Git 2.43 and native Windows Git 2.53.
6. **Canonical paths resolve against the common repository root.** `canonical_worktree`
   is repository-root-relative, but the CLI may run inside a linked worktree whose root is
   different. `ensureWorktree` and claim-context storage resolve against
   `CommonRoot()` (the directory of `--git-common-dir`). Found by dogfooding `claim` from
   the task worktree. Confidence: High.
7. **Claim context is per worktree; flags are authoritative.** The worker's claim id and
   epoch are stored under its worktree's Git dir so `guard`/`heartbeat`/`submit` need no
   extra flags; a resumed old worker that must prove its epoch passes
   `--claim-id`/`--claim-epoch`, which is what the stale-epoch test uses. Recover overwrites
   the convenience copy, so the explicit form is the guarantee. Confidence: Medium.
8. **`integrate` requires tracked files only.** Another session's untracked files in the
   main checkout are normal and must not block a merge; modified tracked files still do.
   Confidence: High.
9. **The wrapper mounts the common repository root.** Mounting only the invoking worktree
   hid `<common>/.git` and every container git command failed. It also forwards the host's
   git author identity so commits inside the container attribute correctly. Confidence:
   High; found by running the wrapper against this repository.
10. **HIGH review runs in the task worktree.** Events emitted before integration are
    committed on the task branch, so `review`/`approve` must run where those events are
    visible (the worktree). `integrate` then operates on main from there. Confidence:
    Medium; a HIGH reviewer should know this and it is documented in the README.

## Validation performed

- `gofmt -l .` clean, `go vet ./...` clean, `go test ./...` green (local Go 1.26.3).
- `tools/queue/queue --build`: image `localhost/ads-queue:1` builds and runs
  `go vet ./... && go test ./...` inside Podman; all packages green.
- `tools/queue/test-crossplatform.sh`: the portability test passes against a
  Windows-visible repository with native `git.exe` reading the same relative worktree
  registration and the same `refs/ads-queue/live/...` ref as WSL git.
- `queue audit` on this repository's real archive: 16 tasks, 0 errors, 0 warnings.
- `queue next --capability LOW` returns this task as ready before the claim and "(no
  eligible task)" after it.
- Race detector not run: the host has no `gcc` and the pinned Go image is used without
  CGO; the only in-process concurrency is the integration-lock test, whose shared state is
  a mutex-protected counter.

## Next

Commit the attempt records, `checkpoint` (once, moving the task to `in_progress`), then
`submit` for HIGH review. No integration, tag or merge is performed here: the brief
requires HIGH approval first, and the required tag `repo/ai-task-queue-v1` belongs to the
integrated state.
