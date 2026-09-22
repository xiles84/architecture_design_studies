# Progress — deliberate multi-leader brainstorm workflow v1

**Actor:** model/effort unknown through Codex desktop; session capability HIGH, work role
executor.

**Task:** `task-20260922T114507Z-brainstorm-workflow-v1`.

**Claim:** `claim-c402a2572067105e`, epoch 1,
attempt `attempt-20260922T115147Z-e0514a`.

**Branch/worktree:** `repo/brainstorm-workflow-v1` in
`.worktrees/brainstorm-workflow-v1`.

## What was built

The existing queue binary now provides explicit intent parsing, durable brainstorm
specs/contributions/events, generated ongoing and concluded views, compact status/end
output, per-slot Git CAS claims, leases/guards/heartbeats, automatic expired-slot
recovery, cancellation/supersession, task correlation, and audit. The repository rules,
workflow/schema, CLI reference, goal amendment, and content templates describe the same
contract.

## Decision Log

1. **Archive commits are a narrow queue-coordination exception.** Brainstorm start,
   submission, and administration commit their immutable records directly in the clean
   local `main` worktree while holding the existing main-integration lock. Reason: the
   next contributor may be in any task worktree and needs the state immediately; a
   separate integration task per contribution would prevent queue-like use. Only exact
   brainstorm paths are staged. Confidence: High. Affects: coordination and AGENTS rules.
2. **Claims are per contribution slot.** Each ready position or critique owns a distinct
   `refs/ads-brainstorms/live/<brainstorm>/<slot>` CAS ref. Reason: independent leaders
   can work concurrently while Git still guarantees one owner per slot. Confidence:
   High. Alternative: one brainstorm-wide claim, rejected because it serializes all
   thinking. Affects: claim, guard, heartbeat, submit, audit.
3. **Expired recovery is part of claim.** There is no separate brainstorm-recover
   command: the next claim CAS-replaces an expired slot and increments its epoch. Reason:
   no worktree or portable draft belongs to an unsubmitted contribution; recovery needs
   to reject the old writer, not preserve unknown local prose. Confidence: High. Affects:
   recovery and documentation.
4. **Record paths never change.** `STATE.md`, `ONGOING.md`, and `CONCLUDED.md` are
   generated views; the immutable spec, events, contributions, and conclusion remain at
   one stable path. Reason: moving a record at conclusion breaks citations and Git
   history. Confidence: High. Affects: navigation and fresh-clone reconstruction.
5. **Conclusion and task authorization remain separate.** Synthesis emits
   `BRAINSTORM ENDED` and “Tasks: NOT CREATED.” The later owner intent authorizes a HIGH
   session to publish tasks, each with `originating_brainstorm_id`, before the workflow
   can link them. Confidence: High. Affects: owner control and provenance.
6. **Blind positions are a workflow boundary, not access control.** A position claimant
   is shown only the question/evidence instruction and must not inspect `positions/`
   until the round closes. Git cannot make already committed paths secret from a local
   contributor. Confidence: High. Affects: bias reduction claims and documented limits.
7. **A task carries its brainstorm origin through every event.** Event emission fills
   `related.brainstorm_id` from immutable task metadata and chain validation rejects a
   mismatch. Reason: later reviews and corrections must trace back to the deliberation,
   not just a prose brief. Confidence: High. Affects: queue schema and audit.
8. **Submission refreshes ownership under the archive lock.** The first implementation
   guarded before waiting for the integration lock, leaving a window in which the slot
   could expire and a new leader could recover it before the old writer committed.
   Submission now revalidates the claim and CAS-extends its lease under that lock before
   writing. Administrative endings reject live claims and CAS-clear expired ones.
   Confidence: High. Affects: exactly-once contribution ownership and cancellation.

## Validation so far

- Containerized `gofmt`, `go vet ./...`, and five consecutive
  `go test ./...` runs pass.
- Tests cover deliberate trigger parsing, ordinary-question non-triggering, full default
  lifecycle and terminal output, retained dissent, fresh-clone reconstruction, twenty
  parallel claimers for two distinct slots, epoch rejection, HIGH-only changes,
  pre-conclusion task refusal, successful correlated task linking, and audit.
- `tools/queue/queue --build` passed and produced pinned runtime image digest
  `9481de7c07cb748ba90b059d3e139f1ac358ccd5d908daf3b840b120530878d4`.
- Real-archive `queue audit`: 18 tasks, 0 errors, 0 warnings. Real-archive
  `brainstorm-audit`: 0 records, 0 errors. Changed-document local links and
  `git diff --check` pass.
- No database or benchmark was started.

## Next

Build the pinned queue image, run repository audits and link checks, update current
context, write the final result receipt, then submit for HIGH review.
