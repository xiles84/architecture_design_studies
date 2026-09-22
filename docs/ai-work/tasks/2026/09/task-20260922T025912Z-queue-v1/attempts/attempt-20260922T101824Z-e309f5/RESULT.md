# Result — implement AI task queue v1

**Outcome:** implementation complete and submitted for HIGH review.
**Capability / role:** LOW / executor.
**Actor:** `deepseek-flash` (tool and effort not exposed by the session).
**Task:** `task-20260922T025912Z-queue-v1`.
**Claim:** `claim-76540de5a8deb0b3` (epoch 1), attempt
`attempt-20260922T101824Z-e309f5`.
**Branch:** `repo/ai-task-queue-v1` — **not integrated**, by the brief's instruction.

## Delivered

- `tools/queue/` — the Go CLI. All commands from the brief:
  `session-start`, `list`, `next`, `publish`, `claim`, `guard`, `heartbeat`,
  `checkpoint`, `yield`, `escalate`, `submit`, `review`, `approve`,
  `request-changes`, `integrate`, `complete`, `recover`, `audit`,
  `review-candidates --since`.
- Live coordination at `refs/ads-queue/live/<task-id>` and
  `refs/ads-queue/locks/main-integration`, each updated with
  `git update-ref <ref> <new> <observed-old>` and reflogs. No central mutable
  `QUEUE.md`.
- Two-hour lease, fifteen-minute heartbeat, two-step recovery with a ten-minute
  `recovery_pending` grace period, a second state check and a new claim epoch;
  stale-epoch heartbeat/guard rejection; refusal to recover while the task's
  worktree holds `ads-run-lock`; worktrees reused and never reset or cleaned.
- A pinned image: `localhost/ads-queue:1`, built by `tools/queue/Containerfile`
  from the Go base in `infra/versions.env`; the build gate runs
  `go vet ./... && go test ./...`, so a failing test produces no image.
- `tools/queue/queue` — the host wrapper: path resolution and `podman run` only,
  no queue logic. It mounts the common repository root at the shell's own path
  and forwards the host git identity.
- `tools/queue/test-crossplatform.sh` — optional host-side native-Windows-Git
  check against a Windows-visible repository.
- `tools/queue/README.md` — the command reference, state-machine mapping,
  lease/recovery rules and the v1 limits.

## Evidence

| Check | Command | Result |
|---|---|---|
| Format / vet / tests (host Go 1.26.3) | `go test ./...` | green; `gofmt -l .` clean, `go vet` clean |
| Canonical Podman gate | `tools/queue/queue --build` | image `localhost/ads-queue:1` built; image ran `go vet ./... && go test ./...`, all packages green |
| Twenty-process atomic claim | `TestParallelClaimersExactlyOneWinner` | exactly one winner (repeated 5x) |
| Integration lock | `TestParallelIntegrationLockExactlyOneHolder` | exactly one holder (repeated 5x) |
| Crash matrix | `TestCrashBeforeWorktreeIsRecoverable`, `TestRecoveryPreservesDirtyWorktree`, `TestFullLifecycle` | green; dirty files and commits preserved |
| Recovery / epoch | `TestGuardRejectsSupersededClaim` | old epoch rejected with exit code 3 |
| Benchmark lock | `TestRecoveryRefusedWhileBenchmarkLockHeld` | recovery refused |
| Event validation | `internal/archive` tests | forks, gaps, duplicates, malformed actor/time, future clock, impossible transition all rejected |
| Fresh clone | `TestFreshCloneReconstruction` | live refs absent, committed events reconstruct state |
| Cross-platform | `tools/queue/test-crossplatform.sh` | native `git.exe` 2.53 and WSL git 2.43 read the same relative worktree registration and the same live ref |
| Real archive | `queue audit` | 16 tasks, 0 errors, 0 warnings |

The working tree is clean and every artifact is committed on the task branch. The
images used are recorded by their config digests in this session's command output;
the requirement is satisfied by the pinned base rather than by a floating tag.

## What was not done, and why

- **No integration, merge or tag.** The brief says submit for HIGH review and do
  not integrate before approval. The required tag `repo/ai-task-queue-v1` marks
  the integrated state, so it is left for the integrator.
- **No push or pull.** The AI rule.
- **Race detector not run.** The host has no `gcc` and the pinned image is built
  with `CGO_ENABLED=0`; the one in-process concurrent path (integration lock) was
  otherwise exercised. The concurrency claims rest on multi-process tests, which
  is the stronger evidence for this design.

## Known limitations and residual risks for HIGH to weigh

- **Bootstrap transition allowance (Decision Log 4).** `model.go` accepts the
  pre-CLI planning record's `proposed → in_progress → completed` shape so `audit`
  is clean on the repository's own archive. Command guards still require
  `approved_for_integration` before `complete`. If HIGH prefers the validator to
  be strict, the corrective path is a new event on that task, superseding its
  history without rewriting it.
- **Claim-id convenience copy (Decision Log 7).** The per-worktree claim copy is
  overwritten by recovery; a resumed old worker must pass
  `--claim-id`/`--claim-epoch` to be rejected deterministically, which the test
  does. The live ref remains the authority either way.
- **Review location (Decision Log 10).** `review`/`approve` must run in the task
  worktree because pre-integration events are committed on the branch. This is
  documented but is the least obvious workflow constraint.
- **`integrate` conflict handling.** It fast-forwards when possible, otherwise
  merges and aborts on conflict, reporting that the main worktree needs manual
  resolution. It does not attempt automatic conflict resolution.
- **v1 scope.** Exclusivity is within one shared Git common directory; a
  multi-clone setup would need the central coordinator the workflow puts out of
  scope.
