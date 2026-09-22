# Queue record schema

Schema version 1 uses human-readable Markdown for instructions and receipts and JSON for
machine-checked task metadata and immutable state transitions. UTC timestamps use RFC3339.

## Task metadata

Every `task.json` contains:

- `schema_version`, `task_id`, `goal_id`, `root_task_id`, and optional `parent_task_id`;
- `title`, `kind`, `created_at`, creator identity, priority `0..100`, and `not_before`;
- `minimum_capability`, `preferred_capability`, and `work_role`;
- dependencies and immutable source request/handoff references;
- canonical branch and worktree, owned paths, and forbidden paths;
- acceptance criteria, expected artifacts, validation, review policy, integration and
  tag requirements;
- `benchmark_required` and an optional `revalidate_after` timestamp.
- optional `originating_brainstorm_id` when the owner authorized tasks from a
  concluded brainstorm.

The task file is immutable after publication. Clarifications live in `amendments/` or in
a successor task linked by `supersedes`.

## Capability normalization

The persisted capability enum remains `HIGH | LOW`. Session-start input accepts the
case-insensitive aliases defined by `WORKFLOW.md`: `leader` and legacy `master` normalize
to `HIGH`; `worker`, `follower`, and legacy `slave` normalize to `LOW`. `primary` and
`replica` are deliberately excluded because they describe datastore topology here.

Actor and live-claim records use `session_capability` for the canonical value and may use
`session_capability_input` for the exact user declaration. Eligibility and state-machine
logic must never compare the raw value. Parsers accept legacy terms but help text,
generated prompts, events, statuses, and routing instructions emit only canonical
`HIGH`/`LOW` or the preferred human-facing `leader`/`worker` pair.

## Event metadata

Each file under `events/` contains:

- `schema_version`, `event_id`, `previous_event_id`, and monotonic `sequence`;
- task, goal, root, parent, attempt, and claim identifiers;
- `occurred_at`, full actor identity, event type, and resulting state;
- base, checkpoint, and result commits when applicable;
- summary, evidence paths, related run/digest/tag/escalation/review/amendment ids;
- optional related brainstorm and contribution ids;
- `next_capability` and `next_work_role`.

The validator rejects a missing predecessor, fork, sequence gap, duplicate id, malformed
timestamp or actor, impossible transition, changed task digest, future clock, or orphaned
live ref.

## Live claim blob

`refs/ads-queue/live/<task-id>` points to a JSON blob containing the current status,
claim and attempt ids, claim epoch, worker identity, capability and role, claim/heartbeat/
expiry timestamps, branch, worktree, base and checkpoint commits, task-spec digest, last
event, next capability/role, and optional benchmark run id.

Claims and heartbeats write a new blob, then use `git update-ref` with the previously
observed object id. One compare-and-swap wins. Reflogs are created for live refs. The
committed event archive, not the ref, is the permanent record.

## Evidence levels

Book claims use exactly one of:

- `repeated_controlled`
- `single_run_directional`
- `mechanism_supported`
- `analogy`
- `gap`
- `invalid`

Every numeric claim records a controlled comparison, run, inputs digest, producing tag,
environment, topology, generated report, signed analysis, and relevant valid cells.

## Brainstorm records

Multi-leader deliberations use their own schema under
[`brainstorms/SCHEMA.md`](brainstorms/SCHEMA.md). They share actor identity,
UTC timestamp, immutable-event, Git-CAS, and local-main-lock conventions with this
queue, but their slot lifecycle is not a task lifecycle.
