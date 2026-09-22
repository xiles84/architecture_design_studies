# HIGH review — AI task queue v1

**Review id:** `review-20260922T113300Z-high`
**Reviewer:** HIGH session; model id and effort not exposed by the session. Work role `reviewer`.
**Task:** `task-20260922T025912Z-queue-v1` (root `task-20260922T025912Z-publish-book-pipeline-handoff`).
**Reviewed submission:** attempt `attempt-20260922T101824Z-e309f5`; branch `repo/ai-task-queue-v1`,
submitted event `event-20260922T102109Z-submitted`, LOW actor `deepseek-flash`.
**Integration base:** local `main` advanced to `e6bc387` (concurrent capability-aliases task) and was
merged into the branch at `4b093fc`; the only conflict (`CONTEXT.md`) was reconciled.
**Verdict:** **accepted for integration into local `main`.**

## Scope reviewed

- The brief, acceptance criteria, and amendment `20260922T101113Z-capability-aliases` (which landed on
  `main` after the LOW submission and explicitly requires alias parsing and tests in queue v1).
- The LOW Decision Log (10 entries), `PROGRESS.md` and `RESULT.md`.
- `tools/queue` (packages `gitx`, `model`, `archive`, `app`), the host wrapper, `Containerfile`, the
  `QUEUE_IMAGE` pin, and the test suite against every required test in the brief.

## Decision Log assessment

| # | Decision | Verdict |
|---|---|---|
| 1 | Command-to-state mapping (`yield` resolves `blocked_high`; `claim` accepts `changes_requested`; first `checkpoint` → `in_progress`) | Accepted; the smallest mapping that keeps every documented state reachable, and it is documented in the README. |
| 2 | `claim`/`recover` take the ref before writing; other transitions commit the event first and CAS after | Accepted; `guard` protects the second case. |
| 3 | Evidence files committed with the event | Accepted; this fixed a real durability bug and is now covered by the lifecycle and fresh-clone tests. |
| 4 | Bootstrap transition allowance in the validator | Accepted as a recorded exception; command guards are untouched, so `complete` still requires `approved_for_integration`. Kept visible for a later corrective event on the planning record if the owner prefers a strict validator. |
| 5 | Worktree reuse detected from the filesystem, not `git worktree list` | Accepted; verified with WSL Git 2.43 and Windows Git 2.53. |
| 6 | Canonical paths resolve against the common repository root | Accepted; fixed a real defect found by dogfooding. |
| 7 | Per-worktree claim-context copy; explicit `--claim-id/--claim-epoch` authoritative | Accepted with the stated limitation; the stale-epoch test uses the explicit form and the live ref remains the authority. |
| 8 | `integrate` requires only tracked files to be clean | Accepted. |
| 9 | Wrapper mounts the common root and forwards host git identity | Accepted; both were found by running the wrapper for real. |
| 10 | `review`/`approve` run in the task worktree | Accepted and documented; pre-integration events live on the branch. |

## Defects found by this review, and their corrections

Both were fixed in this integration pass before the merge, with tests. The LOW attempt's own
`PROGRESS.md`/`RESULT.md` were **not** edited; this review is the record of the corrections.

1. **`STATUS.md` lagged the archive by one event.** `emitEvent` rendered the snapshot from the record
   *before* the new event, so a submitted task's `STATUS.md` still named its checkpoint and
   `in_progress`. Fixed by rendering from a post-event view (`7226e0e`), and the lifecycle test now
   asserts `STATUS.md` matches the newest event's state, sequence and id
   (`assertStatusMatchesLastEvent`).
2. **The capability-alias amendment was unimplemented.** Added `model.ParseCapability` (case-insensitive
   `high`/`leader`/`master` → `HIGH`, `low`/`worker`/`follower`/`slave` → `LOW`, whitespace tolerated,
   `primary`/`replica` rejected), canonical-only storage and eligibility, raw input preserved as
   `session_capability_input` on session records, actors and live claims, and canonical-only help and
   generated claim commands (`1acdba1`). Tests cover every alias, mixed case, whitespace, the excluded
   terms, role independence, and the end-to-end session/claim records.

## Validation performed on the merged state

- `gofmt -l .` clean, `go vet ./...` clean, `go test ./...` green on the host Go 1.26.3.
- `tools/queue/queue --build` rebuilt the pinned image (`120e410ef2e5`) and ran
  `go vet ./... && go test ./...` inside Podman: all packages green.
- `queue audit` on the merged archive: **17 tasks, 0 errors, 0 warnings**.
- `queue review-candidates --since 7d` lists this task; `queue list` shows it `awaiting_review`.
- Earlier, unchanged evidence still stands: the twenty-process claim test, the integration-lock test,
  the crash/recovery matrix, fresh-clone reconstruction, and the native-Windows-Git/WSL-Git
  portability check (`tools/queue/test-crossplatform.sh`).

## Residual risks accepted

The limitations in the LOW `RESULT.md` are accepted: the documented bootstrap transition allowance; the
per-worktree claim-context convenience copy; `review`/`approve` needing to run in the task worktree;
`integrate` failing cleanly on a real conflict rather than resolving it; v1 exclusivity bounded to one
shared Git common directory; and the race detector being unavailable (no host `gcc`), with
multi-process tests standing in as the stronger concurrency evidence.

## Integration instruction

Fast-forward local `main` to this branch and create the annotated tag `repo/ai-task-queue-v1`. Push and
pull remain the owner's; no other tag is moved, deleted or reused.
