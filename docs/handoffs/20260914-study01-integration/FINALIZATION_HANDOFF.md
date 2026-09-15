# Execution Handoff — merge the accepted HIGH review

| Field | Value |
|---|---|
| Handoff | `study01-integration/EH-02`, revision 1; supplements historical EH-01 |
| Planner/reviewer | GPT-6 via Codex desktop, HIGH, 2026-09-15; selected effort not exposed |
| Review checkpoint | `study-01/v3-high-review` (resolve the annotated tag to a commit) |
| Incoming-main baseline | `c7d4067e1c870602e9f794d9067c8163018523dd` |
| Already integrated implementation | `db77fe60b5bfe78b29ec04db9150f33cc1a0bb49`, `study-01/v3-integrated` |
| Acceptance | [HIGH_REVIEW.md](HIGH_REVIEW.md); scientific work and EH-01 integration accepted |
| Remaining work | Safe merge of review documentation, completion record and annotated final tag |
| Next setting | **LOW model / low effort**; record the actual executing model/tool |

## 1. Scope and ownership

Reuse `C:/extra/code/architecture_design_studies/.worktrees/study01-v3`, branch
`study-01/measurement-enhancements`. Main checkout:
`C:/extra/code/architecture_design_studies`. Read AGENTS.md, current CONTEXT.md,
methodology, this handoff and the review first. The task worktree must be clean or
contain only documented partial EH-02 work. Do not use the separate read-only audit
worktree as an executor. No remote operation, source change, measurement, runtime
repair, database teardown or report regeneration is needed or authorized by EH-02.

The HIGH acceptance applies to the exact review checkpoint and preservation conditions
below. All substantive review is complete. LOW may close the task after the mapped
merge without returning to HIGH. Unmapped changes use `SI-ER-NN` in
[ESCALATIONS.md](ESCALATIONS.md); pause only dependent work and return to HIGH.

## 2. Wait for a safe main checkout

Use PowerShell for Git commands. Check exit codes immediately after each native command;
empty output alone is not success. Record UTC time, revisions, model/tool and outcomes
in [PROGRESS.md](PROGRESS.md).

```powershell
git status --porcelain=v1
git branch --show-current
git rev-parse 'study-01/v3-high-review^{commit}'
git merge-base --is-ancestor study-01/v3-high-review HEAD
git -C C:/extra/code/architecture_design_studies status --porcelain=v1
git rev-parse main
podman volume ls --filter name=ads-run-lock --format '{{.Name}}'
podman ps --format '{{.Names}}'
```

As of HIGH review, Study 03 is measuring `20260915T002411Z` from main, with unfinished
result files and the benchmark lock. Wait while any other task owns the lock or main
has edits/untracked files. Inspect the current lock holder when useful; do not assume
the historical run ID remains current. Wait in intervals of at most 60 seconds and
give concise progress updates. Waiting remains LOW work, not Escalation Required.
Do not stage, stash, delete, commit or otherwise adopt the other task's files.
If containers remain without a lock, Podman inspection fails, or ownership is unclear,
record an escalation; no cleanup or lock removal is mapped here.
The owner confirmed another AI is active. Verify actual worktree/branch ownership;
do not move another task or change its checkout to enforce isolation during its run.

## 3. Incorporate and preserve later main work

Once main is clean and the runtime is idle, record its full SHA as `$incomingMain`.
Require it to descend from the incoming baseline. Inspect:

```powershell
git log --oneline c7d4067..main
git diff --name-only c7d4067 main
```

Permitted new main changes are inside `studies/02-ticket-booking/`,
`studies/03-reserved-seating/`, and updates to their context/lessons entries. These may
include their source, settings, results, reports and escalation decisions; preserve
their exact incoming trees. HIGH does not certify their experiments by integrating
them. Changes to Study 01, platform, infrastructure, AGENTS.md, methodology, this task's
handoffs or its own context/lessons meaning require HIGH review before any main update.
Unrelated paths also require escalation rather than silent scope expansion.

From a clean task worktree, integrate the recorded SHA with
`git merge --no-ff --no-commit $incomingMain`. An already-contained main needs no merge.
A clean merge is permitted; verify all conditions below and commit with actual model/tool
attribution, explicitly crediting the incoming session's work. On any conflict, record
paths, abort only this newly started clean-tree merge, and escalate; do not choose a side.

Required zero-exit comparisons after integration:

```powershell
git diff --quiet db77fe6 HEAD -- studies/01-charity-tree infra platform AGENTS.md docs/methodology.md docs/templates
git diff --quiet $incomingMain HEAD -- studies/02-ticket-booking studies/03-reserved-seating
git merge-base --is-ancestor $incomingMain HEAD
git diff --name-only --diff-filter=U
```

For checks before a pending merge commit, omit `HEAD` so Git compares the working tree.
The unresolved-path command must return an empty list. Compare the proposed tree against
`$incomingMain`: only CONTEXT.md, LESSONS_LEARNED.md and files in
`docs/handoffs/20260914-study01-integration/` may differ. Preserve all incoming context
and lessons, and verify this task's HIGH/LOW status and requirements were not lost.
Check local links in the edited files and `git diff --check $incomingMain --` those
explicit paths. Inherited raw Study 03 result whitespace is outside this check.

Recheck both EH-01 integrated tag targets equal `db77fe6`, all four run tags equal the
EH-01 register, and the original enhancement tag still equals `6bf6485`. No result bytes
or signed reports have changed; no benchmark/test rerun is required by this handoff.

## 4. Serialize the final update and close the task

After the mapped checks, acquire a short owned integration reservation through the
existing helper. In the task worktree, use the Git Bash wrapper from EH-01 section 5,
with holder **`Study 01 HIGH review final integration`**:

```powershell
$reviewIntegrationScript = @'
source infra/lib.sh
run_lock_acquire "Study 01 HIGH review final integration"
printf '%s\n' 'INTEGRATION_LOCK_READY'
read -r release_signal
'@
& 'C:/Program Files/Git/bin/bash.exe' -c $reviewIntegrationScript
```

Keep it alive in a PTY/session. Require the ready message and live session ID; record
the holder/worktree/start timestamp. Another task winning the lock returns to the wait.
While this session holds it, its own lock is expected. Recheck its ownership, main's
recorded SHA and clean status, no containers, and the task's clean committed state.
Changed conditions require releasing this owned wrapper and returning to the preflight.

Resolve HEAD to a full `$candidate` SHA, then:

```powershell
git -C C:/extra/code/architecture_design_studies merge --ff-only --no-edit $candidate
git merge-base --is-ancestor $candidate main
```

A rejected fast-forward is a safe stop: reclassify drift and repeat checks. Never force
main. After success, update only in the task worktree:

- CONTEXT.md: record HIGH acceptance plus the observed successful incoming/candidate
  merge, replace the pending-EH-02 paragraph with completed local integration, and list
  `study-01/v3-reviewed` as the final milestone. Preserve every other task entry.
- PROGRESS.md: append the actual commands, UTC observations, incoming/candidate SHAs,
  preservation passes and intended final tag. This is the completion receipt; do not
  overwrite the historical EH-01 receipt or HIGH review.
- LESSONS_LEARNED.md only if a new observed procedural lesson warrants it.

Check links/whitespace, stage explicit owned paths and commit with model/tool attribution.
Recheck main and the owned lock; fast-forward main to this closeout commit and verify
ancestry, clean status and all preservation conditions again. Record its full SHA in
the final response; a commit need not name its own hash inside itself.

Create a new annotated **`study-01/v3-reviewed`** tag pinned to that exact verified
closeout SHA. Check names first. On restart, an existing annotated tag at the recorded
accepted target with the intended meaning counts as completed; a different target or
lightweight tag is an escalation. Never move/delete/reuse any prior tag. Release only
this wrapper's lock by sending Enter to its live session and confirm release. If
another task then obtains a lock or advances main, leave it intact; verify the accepted
closeout remains reachable and describe the later state accurately.

Report the actual final commit/tag and local-main containment. **No further model switch
is required when these conditions pass: the task is complete.** If waiting continues,
state `Next: LOW — resume EH-02's mapped wait`; if an unmapped decision appears, state
`Next: HIGH — resolve SI-ER-NN`. No background follow-up is implied by this handoff.
