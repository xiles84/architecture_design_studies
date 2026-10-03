# Submission portability check, 2026-10-03

The first submission check from the pinned queue image reported thousands of modified
files in the Windows-created sweep worktree, while native Git reported it clean. The
Windows system Git configuration uses `core.autocrlf=true`; the Linux image does not
inherit that system configuration. Unpinned text paths were checked out with CRLF,
and the container interpreted them differently. No reported source change was staged.

A read-only container status check with `core.checkStat=minimal` and
`core.trustctime=false` returned clean. The subsequent pinned-image submission used
those settings plus `core.autocrlf=true`, supplied through per-command `GIT_CONFIG_*`
environment variables, and completed successfully at event
`event-20261003T124549Z-submitted`. This preserves checks of file size/mtime and content
normalization without changing the shared repository configuration, marking files
assume-unchanged, changing another worktree, or rewriting checked-out evidence.

These settings are operational for the audit's queue invocations only. The scientific
backlog and its acceptance gates are unaffected. Native Git clean-tree/diff checks,
containerized queue checks, and local-main ancestry/tag verification remain required.
