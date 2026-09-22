# Result — deliberate multi-leader brainstorm workflow v1

## Outcome

Implemented the owner's opt-in workflow for several HIGH/leader sessions to form
independent positions, challenge them, synthesize a conclusion, and preserve the full
record without automatically authorizing implementation.

The owner-facing triggers are:

- `START BRAINSTORM: <question>`
- `CONTINUE BRAINSTORM <id>`
- `LIST BRAINSTORMS`
- `BRAINSTORM STATUS <id>`
- `CREATE TASKS FROM BRAINSTORM <id>`
- `CANCEL BRAINSTORM <id>`
- `SUPERSEDE BRAINSTORM <old-id> WITH <new-id>`

Ordinary questions do not trigger the workflow. Every state reports a compact summary,
progress, disagreements, and next action. Synthesis emits `BRAINSTORM ENDED`, preserves
dissent, and states `Tasks: NOT CREATED`. A later explicit owner instruction lets the
chosen HIGH session publish correlated queue tasks.

## Delivered

- Durable stable-path archive, immutable event chain, content templates, generated
  ongoing/concluded indexes, and fresh-clone reconstruction.
- Queue commands for intent parsing, start/list/status, per-slot claim/guard/heartbeat,
  submission, cancellation, supersession, task linking, and audit.
- Git CAS contribution leases with parallel distinct-slot ownership, epoch rejection,
  submit-time ownership refresh, and automatic expired-slot recovery.
- HIGH-only enforcement after capability alias normalization.
- `originating_brainstorm_id` on task metadata and automatic
  `related.brainstorm_id` propagation through task events.
- Repository rules, workflow/schema, CLI documentation, goal amendment, context, and
  one reusable ownership lesson.

Implementation commits:

- `bc85878` — core workflow, docs, schema, tests.
- `5c13d55` — submit ownership race fix, expired-claim administration, stronger audit.
- `19cb6ac` — current context, final validation evidence, lesson.

## Validation

- Pinned Podman build `tools/queue/queue --build` passed `go vet ./...` and
  `go test ./...`; runtime image digest:
  `9481de7c07cb748ba90b059d3e139f1ac358ccd5d908daf3b840b120530878d4`.
- Five consecutive containerized test runs passed.
- Tests cover the full default 3/2/1 lifecycle, exact terminal output, conclusion with
  dissent, trigger/non-trigger parsing, HIGH-only mutation, twenty parallel claimers,
  distinct slots, expired recovery, stale epoch rejection, active/expired cancellation,
  pre-conclusion task refusal, successful correlated task linking, malformed chains and
  identities, generated state, audit, and fresh-clone reconstruction.
- Repository queue audit: 18 tasks, 0 errors, 0 warnings.
- Repository brainstorm audit: 0 records, 0 errors.
- Changed-document local links and `git diff --check` pass.
- No database, benchmark, push, or pull was run.

## Known boundaries

- Independent positions are a documented workflow boundary, not filesystem secrecy;
  contributors with repository access could inspect committed positions.
- Live refs coordinate one shared Git common directory. Independent clones still need a
  future central coordinator.
- An unsubmitted draft remains local to its contributor; only submitted records are
  portable.
- Archive commands intentionally require a tracked-clean local `main` and fail while
  another integrator owns the shared lock.
- The implementation created no example brainstorm in the real archive. Its generated
  indexes correctly remain empty; full lifecycle evidence comes from disposable Git
  repositories in the test suite.

## Review focus

Review the submit-time lease refresh, direct-main coordination exception, explicit
conclusion/task boundary, and the honest limitation on blind positions. No scientific
evidence or protected analyst artifact changed.
