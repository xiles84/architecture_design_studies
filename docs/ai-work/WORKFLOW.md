# AI work queue

This directory is the durable coordination archive for AI-managed work. It preserves
what the owner asked, how HIGH interpreted it, which task a worker performed, every
handoff or escalation, and the evidence behind review and completion. A user should not
have to copy a prompt from one chat to another.

## Start every session

1. Read `AGENTS.md`, `CONTEXT.md`, this file, and the relevant task brief.
2. If the user has not already declared a session capability, ask once before claiming
   work. Normalize an accepted alias to `HIGH` or `LOW` and record the canonical value as
   `session_capability`; never infer it from a model.
3. Record the actual model, effort, tool, and session id when exposed. Write `unknown`
   when they are not exposed.
4. Inspect `git worktree list --porcelain`, the task branch, current queue state, and the
   `ads-run-lock` before any database or benchmark operation.

Capability is not a work role. Work roles are `planner`, `executor`, `reviewer`,
`analyst`, and `integrator`. HIGH may perform any role. LOW may claim only a task whose
`minimum_capability` is `LOW`.

## Capability vocabulary

Capability declarations are case-insensitive and normalize as follows:

| User declaration | Canonical capability |
|---|---|
| `high`, `leader`, `master` | `HIGH` |
| `low`, `worker`, `follower`, `slave` | `LOW` |

The declaration must be explicit or given in answer to the session-capability question;
the same words in a database or ordinary-language discussion do not change the session.
Store the raw declaration as `session_capability_input` when available, but use only the
canonical value for eligibility, events, live refs, task metadata, and routing. Accept
`master` and `slave` as legacy input, but never generate or recommend them. Human-facing
examples prefer `leader` and `worker`. `primary` and `replica` remain datastore-topology
terms and are not capability aliases.

## Selecting and claiming work

The queue order is:

1. dependencies complete;
2. `not_before` reached;
3. session capability eligible;
4. priority descending;
5. creation time ascending;
6. task id ascending.

The planned Go queue CLI stores live state in `refs/ads-queue/live/<task-id>` and claims
it with Git compare-and-swap. Until that CLI is implemented, use the task's exact branch
name from `task.json` as the bootstrap atomic claim:

```text
git worktree add <canonical-worktree> -b <canonical-branch> main
```

Git permits only one creation of that branch. If the branch or registered worktree
already exists, do not create an alternative name. Inspect `STATUS.md`, events, branch,
and worktree and either resume the same task or leave it to its current worker.

## States

```text
proposed -> ready -> claimed -> in_progress
                              -> yielded -> ready
                              -> blocked_high -> ready
                              -> awaiting_review
                                   -> changes_requested -> ready
                                   -> approved_for_integration -> completed
```

`cancelled` and `superseded` are terminal alternatives. A task is not `completed` until
its required checks pass, its result and review exist, its commits are reachable from
local `main`, and required annotated tags exist.

## Events, attempts, and recovery

- `task.json` is immutable. A changed requirement is an amendment or a superseding task.
- `STATUS.md` is a readable snapshot; immutable JSON event files are the history.
- Each event names its predecessor, task/goal/root/parent ids, attempt and claim ids,
  actor, UTC timestamp, resulting state, commits, evidence, and next capability/role.
- Each worker owns a new attempt directory. Never edit another worker's receipt.
- Default live-claim lease: two hours. Heartbeat at least every 15 minutes and before
  and after substantive mutations. Extend it before a mapped long run.
- Recovery requires an expired lease, a ten-minute `recovery_pending` grace period, a
  second state check, inspection of the canonical branch/worktree, and confirmation that
  the task does not own the benchmark lock. Reuse the existing worktree and preserve all
  dirty and untracked files. Never reset, clean, or discard them.
- Run the future `queue guard` before edits, commits, measurements, tags, and integration.
  A worker whose claim epoch is no longer current stops immediately.

Queue ownership never grants permission to measure. `ads-run-lock` remains the only
benchmark lock. `refs/ads-queue/locks/main-integration` separately serializes final
updates of local `main`.

## HIGH/LOW handoff

HIGH creates decision-complete briefs and acceptance criteria. LOW decides and logs
local, reversible implementation details, and escalates only under the material-risk
rules in `AGENTS.md`. Submission moves the task to `awaiting_review`; HIGH reviews the
Decision Log and result before integration. A HIGH session may also execute a LOW task.

All finished work is committed, tagged where required, merged into local `main`, and
never pushed or pulled by an AI.

## Portability boundary

Committed task records are the portable archive. Custom Git refs are live coordination
for agents sharing this one local Git common directory; they are not pushed by default.
Concurrent workers operating from independent clones would require a central coordinator
and are outside queue v1.
