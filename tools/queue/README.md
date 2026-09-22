# `queue` — the durable AI work queue (v1)

The operating contract is [`docs/ai-work/WORKFLOW.md`](../../docs/ai-work/WORKFLOW.md)
and [`docs/ai-work/SCHEMA.md`](../../docs/ai-work/SCHEMA.md). This directory is
the implementation: a Go CLI that owns live coordination in Git refs, and
committed event files under `docs/ai-work/tasks/` that remain the permanent
archive.

Nothing here starts a database or measures, so the queue never takes
`ads-run-lock`. `recover` *reads* it only to refuse recovering a task whose
worktree still holds the benchmark lock.

## Quick start

```bash
tools/queue/queue --build                 # build the pinned image (once)
tools/queue/queue list --capability LOW
tools/queue/queue next --capability LOW --claim-command
tools/queue/queue session-start --capability LOW --role executor
```

Everything runs in the pinned image `localhost/ads-queue:1` (built from
`Containerfile`, base pinned in `infra/versions.env`). The only host
requirement is Podman; `tools/queue/queue` resolves paths and starts the image,
and must contain no queue logic of its own.

## Commands

| Command | Effect |
|---|---|
| `session-start` | record this session's capability and role under `docs/ai-work/sessions/` |
| `list` | one row per task with its derived state; `--eligible` filters to what this capability may claim |
| `next` | the next eligible task in queue order (`--claim-command` prints the claim line) |
| `publish` | HIGH creates a task from `--spec task.json --brief BRIEF.md` (first event `published`), or releases a pre-created `proposed` task with `--release --task <id>` (`proposed → ready`) |
| `claim` | atomic claim; creates the canonical worktree and branch |
| `guard` | exits non-zero if this worker's claim id/epoch is no longer current |
| `heartbeat` | extends the lease with a compare-and-swap |
| `checkpoint` | records the current commit; first one moves `claimed → in_progress` |
| `yield` | releases a claim back to `ready` |
| `escalate` | writes `ESCALATION REQUIRED`, moves to `blocked_high`, releases the claim |
| `submit` | moves to `awaiting_review`; requires a committed result and a clean tree |
| `review` | HIGH records `reviews/<id>/REVIEW.md` |
| `approve` | `awaiting_review → approved_for_integration` |
| `request-changes` | `awaiting_review → changes_requested`, back to a LOW executor |
| `integrate` | takes the local-main lock, merges the branch, creates the required tag |
| `complete` | verifies result, review, reachability and tag; `→ completed` |
| `recover` | two-step recovery of an expired claim with a grace period |
| `audit` | validates every event chain and live ref; `--json` for machine output |
| `review-candidates` | `awaiting_review` tasks submitted within `--since` (default `7d`) |

### Deliberate brainstorm commands

| Command | Effect |
|---|---|
| `brainstorm-intent` | parse only the documented explicit owner phrases; ordinary questions return no intent |
| `brainstorm-start` | HIGH creates the question/evidence packet and independent-position round |
| `brainstorm-list` / `brainstorm-status` | show compact ongoing/concluded state and next action |
| `brainstorm-claim` | HIGH atomically claims the next ready contribution slot |
| `brainstorm-guard` / `brainstorm-heartbeat` | verify or extend a slot lease |
| `brainstorm-submit` | append a position, critique, or synthesis; synthesis emits `BRAINSTORM ENDED` |
| `brainstorm-cancel` / `brainstorm-supersede` | preserve an administrative ending |
| `brainstorm-link-tasks` | after explicit owner authorization, verify and link already-published tasks |
| `brainstorm-audit` | validate specs, event chains, actors, transitions, and contribution files |

The user phrases and content contract are in
[`docs/ai-work/brainstorms/WORKFLOW.md`](../../docs/ai-work/brainstorms/WORKFLOW.md).
Brainstorm slot refs use
`refs/ads-brainstorms/live/<brainstorm-id>/<slot-id>`. Archive writes reuse the
local-main integration lock. Expired slots are reclaimed by `brainstorm-claim` with a
higher epoch; there is no separate recover command. A fresh clone reconstructs submitted
state, but not an unsubmitted contributor draft.

Common flags: `--repo`, `--now` (freeze the clock; used by tests), `--json`,
`--capability`, `--role`, `--model`, `--tool`, `--effort`, `--session-id`.

## Live coordination

One blob per task lives at `refs/ads-queue/live/<task-id>`; the local-main
integration lock is `refs/ads-queue/locks/main-integration`. Writes are
`git update-ref <ref> <new> <observed-old>` with reflogs, so Git's own ref lock
is the mutual exclusion: twenty workers on the same common Git directory, or
two different git binaries, produce exactly one winner.

`task.json` is immutable; every state change is a new file under `events/`
linked by `previous_event_id`. `task_spec_digest` is a SHA-256 over the
LF-normalised `task.json` bytes, so a CRLF checkout on Windows cannot make a
`guard` reject a task that never changed.

## How commands map onto the state machine

The diagram in `WORKFLOW.md` has no separate command for every arrow, so the
mapping is:

```text
proposed --publish--> ready --claim--> claimed --checkpoint--> in_progress
  claimed|in_progress --yield--> ready              (released)
  claimed|in_progress --escalate--> blocked_high --yield--> ready   (HIGH resolved)
  claimed|in_progress --submit--> awaiting_review
  awaiting_review --review--> awaiting_review        (review recorded)
  awaiting_review --approve--> approved_for_integration --integrate--> (same state, commit recorded) --complete--> completed
  awaiting_review --request-changes--> changes_requested --claim--> claimed
```

`claim` accepts a `changes_requested` task directly rather than detouring
through `ready`, and `recover` records a new epoch as `claimed` even from
`in_progress`; both widenings are in the transition table with comments. This is
the only place the command list forced an interpretation. It is logged in the
task's Decision Log.

## Leases and recovery

- Default lease 2 h; heartbeat interval 15 min; a heartbeat refreshes both
  `heartbeat_at` and `expires_at`.
- `recover` first compare-and-swaps the live blob to `recovery_pending` with
  `recovery_ready_at = now + 10 min`, then, after the grace period and a second
  state check, compare-and-swaps a fresh claim with `claim_epoch + 1`. A
  heartbeat or guard that still presents the old epoch is rejected with exit
  code 3.
- Recovery reuses the existing worktree and never resets or cleans it: a crashed
  worker's uncommitted files are exactly what recovery must preserve. It refuses
  while this task's worktree holds `ads-run-lock`.
- The worker's own claim id and epoch are saved inside its worktree's Git dir
  (`<git-dir>/ads-queue/claims/<task>.json`), so `guard`/`heartbeat`/`submit`
  work without restating them; passing `--claim-id`/`--claim-epoch` explicitly
  is the authoritative form and is what a resumed old worker must use.

## Capability and roles

Session capability (`HIGH`/`LOW`) is asked once per session and recorded, never
inferred from a model name. `LOW` may claim only tasks whose
`minimum_capability` is `LOW`; `HIGH` may claim either. Capability is separate
from the work role (`planner`, `executor`, `reviewer`, `analyst`, `integrator`),
which is carried in every event's actor.

`--capability` and `ADS_QUEUE_CAPABILITY` also accept the case-insensitive
aliases in `docs/ai-work/WORKFLOW.md`'s capability vocabulary (for example
`leader` for `HIGH`, `worker` for `LOW`, with legacy compatibility terms
accepted only as input). Every declaration is normalized immediately; only the
canonical `HIGH`/`LOW` is stored, compared or emitted, and the raw wording is
preserved separately as `session_capability_input`. `primary`/`replica` are
datastore-topology terms and are rejected as capability declarations.

## Portability

- `git worktree add` writes absolute paths on Git 2.43; the CLI rewrites the
  linked worktree's `.git` link and the `.git/worktrees/<name>/gitdir`
  back-pointer to **relative** paths, so the same registration resolves under
  WSL Git, native Windows Git and a Podman bind mount. Git reads relative links
  even where it does not write them.
- Because Git 2.43 cannot resolve a relative back-pointer in `git worktree
  list` (it reports the worktree as prunable), the CLI detects an existing
  worktree from the filesystem and by opening it, not from `worktree list`.
- The wrapper mounts the repository at the same absolute path the shell uses, so
  no container-internal path ever reaches a worktree registration or a committed
  record.
- `tools/queue/test-crossplatform.sh` (optional, needs a host Go) runs the
  portability test against a Windows-visible repository with `git.exe` as the
  second binary.

## Tests

Run in Podman through the image build, or directly:

```bash
tools/queue/queue --build          # runs go vet ./... && go test ./...
```

or, with a host Go toolchain:

```bash
(cd tools/queue && go vet ./... && go test ./...)
```

The suite covers: twenty real processes racing one claim (exactly one winner),
HIGH/LOW filtering, dependency/priority/FIFO/`not_before`/seven-day UTC
boundaries, the crash matrix (before the worktree, with dirty files, after
commits, before submission), expired-lease recovery and rejection of the old
epoch, no recovery while the task owns the benchmark lock, exactly one
integration-lock holder, event forks/gaps/duplicates/malformed actors and
times/future clocks/impossible transitions, changed task digests, orphan live
refs, parent/root correlation through the full lifecycle, and fresh-clone
reconstruction from committed events alone. Brainstorm coverage includes explicit
intent parsing, ordinary-question non-triggering, staged position/critique/synthesis
transitions, parallel distinct-slot claims, expired-epoch rejection, HIGH-only changes,
exact terminal output, task-boundary enforcement, audit, and fresh-clone reconstruction.

## Limits of v1

- Exclusivity holds within one shared Git common directory. Independent clones
  would need a central coordinator; that is out of scope.
- `integrate` requires the main worktree to be free of *tracked* modifications;
  untracked files from another session do not block it, and it never discards
  them.
- Push and pull are the owner's; the queue never touches remotes.
