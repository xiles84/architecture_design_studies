# Execution brief — AI task queue v1

## Goal

Implement the Go queue CLI and pinned Podman wrapper described by `docs/ai-work/WORKFLOW.md`
and `SCHEMA.md`. The CLI owns live coordination; committed files remain the durable archive.

## Required commands

`session-start`, `list`, `next`, `publish`, `claim`, `guard`, `heartbeat`, `checkpoint`,
`yield`, `escalate`, `submit`, `review`, `approve`, `request-changes`, `integrate`,
`complete`, `recover`, `audit`, and `review-candidates --since <duration>`.

Use `refs/ads-queue/live/<task-id>` and `refs/ads-queue/locks/main-integration`.
Serialize live JSON as Git blobs and update refs with the observed old object id. Create
reflogs. Do not implement a central mutable `QUEUE.md`.

Default lease is two hours; heartbeat interval is fifteen minutes. Recovery uses an
expired lease, `recovery_pending`, a ten-minute grace period, a second compare-and-swap,
benchmark-lock and worktree inspection, and a new claim epoch. It preserves existing
files and rejects a late heartbeat from the old epoch.

The thin host wrapper may only resolve paths and start the pinned Podman image. Queue
semantics and validation stay in Go. Do not create worktrees from paths meaningful only
inside the container; the CLI should emit or execute host-safe canonical commands using
the repository's existing path helpers.

## Required tests

- Twenty parallel claimers: exactly one winner.
- LOW cannot claim HIGH-only; HIGH can claim either.
- Dependencies, priority, FIFO, `not_before`, and seven-day UTC boundaries.
- Crashes before worktree creation, with dirty files, after commit, and before submission.
- Expired recovery and rejection of the previous epoch.
- No recovery while the task owns `ads-run-lock`.
- Exactly one integration-lock holder.
- Event forks/gaps/duplicates, malformed actor/time, changed task digest, and orphan refs.
- Parent/root correlation through escalation, review, correction, and completion.
- Fresh-clone reconstruction from committed events.
- Native Windows Git and WSL Git against the same common Git directory.

Use Decide → Log → Continue for local implementation choices. Escalate only if Git cannot
provide atomicity on the supported shared-repository setup, the Podman-only rule cannot be
met, or the archive cannot reconstruct live state.

On implementation completion submit for HIGH review; do not integrate before approval.

NEXT MODEL: LOW
