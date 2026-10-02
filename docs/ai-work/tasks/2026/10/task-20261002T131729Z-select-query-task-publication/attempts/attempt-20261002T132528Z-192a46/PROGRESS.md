# Publication Decision Log

Actor: GPT-6 / Codex desktop / effort unknown; session HIGH (owner answer HIGH (leader)).
Task: task-20261002T131729Z-select-query-task-publication.

## Decisions

- Extend Study 01 additively rather than create a new customer/orders study. Reason: reuse the current relational domain and avoid bypassing replication-first/new-domain gates; toy customer data remains a teaching analogue. Confidence: High. Alternative: a new study number. Affects: protocol, harness, all four measurements and book authorship.
- Publish eight dependent tasks instead of implementing now. Reason: the owner explicitly requested task creation; HIGH protocol freezes exact scientific choices before LOW measurement. Confidence: High. Affects: complete program.
- Partition semantically different outputs and include same-plan/no-resolved-difference results. Reason: performance cannot rank a different or invalid answer; an optimizer-collapsed rewrite still teaches useful intuition. Confidence: High. Affects: oracle, reports and synthesis.
- Retain Edition 3's approved concept-first style and pilot integration dependency. Reason: concurrent book work owns that preview/approval; this request adds SELECT scope without approving another task's pending result. Confidence: High. Affects: book integration.
- Use the pinned queue image through direct Podman orchestration for multi-command publication after WSL wrapper preflight. Reason: the native worktree had CRLF shell/environment files and a hidden absolute Git link; only this task's Git registration was made relative, temporary shell edits are restored and no runtime/harness change is committed. Confidence: High. Affects: publication operation only.
- Align core.autocrlf=true in this task's own Git worktree config with its Windows checkout before containerized submission. Reason: the container's default false setting treated unchanged CRLF files as edits and repeated a slow full-tree scan. Restore only the owned temporary script edits using the matching setting; final native diff shows just context/lessons plus new task documents. Confidence: High. Affects: cross-platform queue operation only.

## Validation baseline

Containerized queue audit before publication: 74 tasks, 0 errors, 2 pre-existing warnings
on task-20261001T133759Z-t2-study02-replication (expired lease and live_ref_drift).
No ownership recovery, event edit or ref change was made for that peer task.
Existing ads-run-lock was absent after starting the existing Podman machine; no database
or benchmark was started. Native/main checkout was clean at initial isolation.

All eight downstream tasks were published. Final pre-submission queue audit: 83 tasks,
0 errors, the two peer warnings plus this task's temporary live_ref_drift after claim.
Normal submit releases the owned live ref. Brainstorm audit: 3 records, 0 errors.
Dependency/uniqueness/role/link and scoped-diff checks passed. The read-only diagnostic
container that displayed line-ending noise was stopped; it was owned by this task.

## Routing

Task publication finishes with local-main integration; downstream execution remains queued.

TASK COMPLETE
