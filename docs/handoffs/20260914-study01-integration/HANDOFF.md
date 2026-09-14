# Execution Handoff — finish Study 01 integration

| Field | Value |
|---|---|
| Handoff | `study01-integration/EH-01`, revision 1 |
| Planner | GPT-6 via Codex desktop, HIGH role, 2026-09-14; selected effort not exposed |
| Reviewed integration | `e3e31bcb95b38b7fa92e882be02b435d54d0e56a` |
| Reviewed local main | `025b72f184434a256c106f21bf78edefb5641c2d` |
| Planning checkpoint | `repo/study01-integration-handoff-v1` |
| Status | ready for execution; task not yet merged |
| Next setting | **LOW model / low effort — execute this handoff** |
| Logs | [PROGRESS.md](PROGRESS.md), [ESCALATIONS.md](ESCALATIONS.md) |

## 1. Outcome and scope

Finish the existing task on `study-01/measurement-enhancements`: put its completed
Study 01 v3 enhancements and repository workflow rules on local `main`, preserve the
other studies, record the outcome, and create annotated integrated milestones. No new
database experiment, result interpretation, source modification, push or pull is part
of this execution. Study 03's ER-02 belongs to its separate task and remains unresolved.

Reuse `C:/extra/code/architecture_design_studies/.worktrees/study01-v3`. The local main
checkout is `C:/extra/code/architecture_design_studies`. Inspect the branch/status and
read AGENTS.md, CONTEXT.md, methodology and this entire handoff before executing.
The named paths are verified Windows paths; commands below are for PowerShell, except
the explicitly labelled Git Bash lock wrapper. Use the actual model identity in commits.

This is a LOW execution checkpoint in the same unfinished task. The standing merge
authorization already covers local integration. HIGH has reviewed the exact source
combination and the acceptance rules below; LOW executes their mechanical checks.
After execution, **return to HIGH to judge the evidence and close the task**.

## 2. Decisions made by HIGH

- The original 261-cell measurement record and its signed analyses remain valid for
  their producing revisions. This integration does not require a new benchmark: Study 01
  keeps its own harness, and its infrastructure changes were previously validated with
  default budgets unchanged. Any source change that invalidates that premise is an ER.
- Integration `e3e31bc` incorporates all committed main work through `025b72f`. Both
  shared living documents merged without conflict. Incoming Study 02, Study 03 and
  platform trees matched main; Study 01 and infrastructure matched the tested task tree.
- The previous [integration review](../../../studies/01-charity-tree/reports/20260914-v3-integration-validation.md)
  concerns incoming `1b989d3`; keep it intact. Record this newer combination in the
  task's progress log and a new `EXECUTION_RESULT.md` beside this handoff.
- On this HIGH pass, both worktrees were clean, Podman answered, and neither a benchmark
  lock nor running containers were present. These are observations, not execution-time
  guarantees. The old ER-01 wait is over; do not resume waiting on its historical log.
- Repository-wide HIGH/LOW rules are canonical in AGENTS.md. Existing Study 03 handoffs
  and their particular model mapping are preserved, not rewritten by this task.

## 3. Preflight and mapped conditions

Run from the task worktree. Check exit codes immediately; `git diff --quiet` and
`git merge-base --is-ancestor` return 0 for success, 1 for a real difference, and larger
values for an execution error. A blank output alone is not a successful check.

```powershell
git status --porcelain=v1
git branch --show-current
git rev-parse main
git rev-parse 'repo/study01-integration-handoff-v1^{commit}'
git merge-base --is-ancestor repo/study01-integration-handoff-v1 HEAD
git -C C:/extra/code/architecture_design_studies status --porcelain=v1
podman volume ls --filter name=ads-run-lock --format '{{.Name}}'
podman ps --format '{{.Names}}'
```

| Observation | LOW action |
|---|---|
| Own branch matches, own worktree clean, main is reviewed `025b72f`, main checkout clean, no lock/containers | Proceed. |
| Resuming with only commits/edits recorded in this handoff's progress log | Verify their origin and resume the first incomplete step; do not repeat a completed merge/tag. |
| Benchmark lock exists | Inspect its holder; wait in intervals no longer than 60 seconds, with concise updates. Do not remove it or change main. This is an expected wait; next model remains LOW. |
| Containers exist with no lock, stale-looking lock, or Podman cannot be inspected | Preserve state; record SI-ER-NN. No teardown, lock removal or runtime repair is authorized here. |
| Main checkout has tracked or untracked work | Leave it untouched and wait for its task to finish. Do independent read-only checks. If ownership/overlap needs a decision, record SI-ER-NN. Do not stash, stage, delete or commit those files. |
| Own worktree contains unexplained changes or another operation is in progress | Record SI-ER-NN; do not adopt the changes. |
| Main advanced beyond the reviewed commit | Read the new commit list/diff. If all changes are within Study 02/03 and additive context/lessons updates, LOW may integrate them as described below, then return to HIGH for renewed acceptance before changing main. Any other scope or non-additive shared-document change is SI-ER-NN. |
| Main equals the recorded successful candidate or closeout commit on a restart | Verify the preserved trees and log; continue the first incomplete step. Do not duplicate integration. |
| Main is a later descendant of a recorded target | Do not treat ancestry alone as acceptance of the new state. Classify the intervening commits under the drift rule before any further main update. Already-created tags on the recorded verified closeout commit remain valid. |

For the mapped main advance, capture its full incoming-main SHA, start from a clean task worktree and use
`git merge --no-ff --no-commit main`. If it merges cleanly, verify preservation, commit
with attribution to the incoming session, and set `Next: HIGH — review updated integration`.
In that drift path, compare incoming Study 02/03/platform against the newly captured
main SHA, not `025b72f`; retain `0efd760` for Study 01/infrastructure. Record both SHAs
and return to HIGH before updating main; HIGH supplies the amended acceptance baseline.
If a conflict appears, capture the paths/error, abort **only this just-started clean-tree
merge** with `git merge --abort`, then write/commit an escalation. Never choose a side.
No automatic merge of changed AGENTS.md, methodology, infrastructure, platform or Study 01
is authorized by this handoff.

## 4. Acceptance checks before updating main

For the exact reviewed main, run:

```powershell
git diff --quiet 0efd760 -- studies/01-charity-tree infra
git diff --quiet 025b72f -- studies/02-ticket-booking studies/03-reserved-seating platform
git diff --check 025b72f -- AGENTS.md CONTEXT.md LESSONS_LEARNED.md docs infra
git diff --name-only --diff-filter=U
git merge-base --is-ancestor main HEAD
```

The first two checks must pass; they protect all existing study source, results,
reports, analyst files and the incoming platform. The third checks this task's document
and infrastructure changes, not inherited raw-output whitespace. The fourth must list
no conflict, and the fifth must succeed. Resolve local links in the edited rules, this
handoff and its logs. Verify the four existing run tags against this register:

| Run tag suffix under `run/01-charity-tree/` | Producing commit | Inputs digest |
|---|---|---|
| `20260913T124917Z-v3` | `41a1490de485f2fc905e36f2486842eee0aa2592` | `d2d1cc8380b381f2` |
| `20260913T125342Z-v3` | `c084c73b0cade873406c57fa65ce61f5a26c505a` | `0bc5900a6d7a5d3a` |
| `20260913T172624Z-v3` | `9b7aa48b329ca44e341fe2aedb6d39de8d4cb70d` | `bd95d304cb39e2de` |
| `20260914T104721Z-v3` | `d1b18bf437daf972860261e53b422dc87905cc09` | `9ead901cb74c4183` |

Use `git rev-list -n 1 <tag>` and the generated reports' digest/index rows. Record each
pass in PROGRESS.md; create EXECUTION_RESULT.md with a TL;DR, actual model/effort if
known, UTC time, incoming main, candidate commit, commands and results. Commit these
explicit paths. No test rerun or report regeneration is needed if these checks pass.

## 5. Serialize and perform the local integration

Reserve the existing machine lock for the short final integration so another benchmark
cannot start between the idle check and the shared-infrastructure update. Use
`run_lock_acquire` from this task's `infra/lib.sh` in a Git Bash session; do not invent
or manually remove a volume. Invoke Git Bash explicitly from PowerShell:

```powershell
$integrationScript = @'
source infra/lib.sh
run_lock_acquire "Study 01 final local integration"
printf '%s\n' 'INTEGRATION_LOCK_READY'
read -r release_signal
'@
& 'C:/Program Files/Git/bin/bash.exe' -c $integrationScript
```

Keep that terminal/session open while doing the following steps in separate tool calls.
With the Codex execution tool, allocate a PTY (`tty: true`) and use a short initial
yield. Require both `INTEGRATION_LOCK_READY` and a live session ID before continuing;
an exited command has released the lock. Send Enter through that same session to end it.
The wrapper's existing EXIT trap releases only its owned lock when LOW sends Enter at
the end. If the lock was acquired by another task first, the helper refuses; go back to
the mapped wait. If the wrapper exits or ownership is uncertain, recheck before proceeding.
Hold it only for integration, not while waiting for model changes or another task's edits.

Immediately recheck main's commit/status, no unexpected containers, the lock holder,
and own clean worktree. If anything differs, release this session's lock by ending its
wrapper and return to section 3. Resolve the candidate as a full commit hash and record it.

During these integration/closeout preflights, **this wrapper's own lock is expected**;
section 3's ordinary "lock exists" wait applies to another holder. Record the acquired
lock's holder, worktree and start timestamp plus the live wrapper session ID, and require
those same values while holding it. A dead wrapper or changed ownership invalidates the
reservation. Do not wait on or remove your own lock through a separate volume command.

```powershell
git rev-parse HEAD
git -C C:/extra/code/architecture_design_studies merge --ff-only --no-edit <candidate-sha>
git merge-base --is-ancestor <candidate-sha> main
```

The merge command is the sole initial write to the main checkout. A rejected fast-forward
is a safe stop, not a reason to force the reference. Re-evaluate drift under section 3.

After success, edit **in the task worktree only**:

- CONTEXT.md: replace the handoff-ready paragraph with the actual merged incoming/candidate
  revisions, note validation evidence collected and `Next: HIGH — final evidence review`.
  Mark repository integration complete, while final task review is pending. Add the
  integrated tags below and link EXECUTION_RESULT.md. Preserve every other task entry.
- PROGRESS.md and EXECUTION_RESULT.md: record the successful first fast-forward and
  checks. State the intended final milestone tags and that HIGH review is next.
- LESSONS_LEARNED.md only if execution exposes a new procedural lesson; use observations
  and do not add speculative performance conclusions.

Commit those explicit owned paths, then repeat the preflight and fast-forward main to
that closeout commit. If main advanced in between, follow section 3. Verify containment
of both task commits and the unchanged Study 01/incoming-study trees. The final commit
need not be inserted into itself: the tag identifies the final revision; report its hash
in the response and add it to progress at a later checkpoint if needed.

Record the full closeout SHA after its successful fast-forward and verification. Pin
both tag targets to **that recorded closeout SHA**, not a later value of `main`.
Create these **new annotated tags on that verified commit**, with annotations
describing the scope and identifying the executing model/tool:

- `repo/worktree-to-main-workflow` — standing worktree, local merge and HIGH/LOW workflow.
- `study-01/v3-integrated` — the integrated Study 01 v3 measurements, concise analysis
  and discussion companions, preserving all other studies.

Check names before creation. If a tag already exists with the same recorded verified target and
intended meaning from this execution, treat it as a completed step. A different target,
lightweight tag or uncertain meaning is SI-ER-NN; never move, delete or rename around it.
Do not alter `study-01/v3-enhancements` or any producing-run tag. Release the owned lock,
confirm its release, and leave worktrees available for review. Never push, pull or change
remotes. If a new task acquires the lock immediately afterward, leave its lock intact.

## 6. Escalation and model transition

Use [ESCALATIONS.md](ESCALATIONS.md), IDs `SI-ER-01`, `SI-ER-02`, etc. Record evidence and
blocked steps for any unmapped conflict, preservation failure, tag conflict, source
change, uncertain lock ownership or scope change. Release any owned integration lock
before a model-switch pause. Keep completed work committed and attributable.

On success, report actual main revision, both tag targets, clean/dirty state accurately,
and the execution-result link. End:

**Next: HIGH model / high effort — review the execution evidence and close the task.**

This HIGH review may accept the result or publish another Execution Handoff if it finds
a deficiency. LOW must not certify a new scientific interpretation or quietly relax a
failed acceptance check. On ordinary waiting, retain `Next: LOW — resume the mapped wait`.
