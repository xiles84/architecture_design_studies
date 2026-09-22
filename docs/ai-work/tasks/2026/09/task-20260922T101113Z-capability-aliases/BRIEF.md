# Execution brief — AI capability aliases

## Goal

Let the owner start an AI session with familiar leader/worker terminology without
changing the queue's canonical capability values or confusing AI capability with
database topology.

## Decided mapping

- `leader` and legacy `master` mean `HIGH`.
- `worker`, `follower`, and legacy `slave` mean `LOW`.
- `primary` and `replica` are not aliases.
- Persist `HIGH`/`LOW`; preserve the raw declaration separately when available.
- Accept but never generate or recommend the legacy pair.

Update the live rules and queue contracts, publish amendments without rewriting the
original goal or queue-v1 task, record the decision in current context and lessons, and
merge the tagged result into local `main`. No benchmark is authorized.

TASK COMPLETE
