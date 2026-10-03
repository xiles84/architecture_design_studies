# Runner path validation correction

2026-10-03, GPT-6 via Codex desktop. The control runner acquired and released
`ads-run-lock`, then failed before creating a run tag or starting a database:
Podman interpreted the Git Bash `/c/extra/...` build context as `C:\c\extra\...`.
The task-owned runner now applies `hostpath` before `buildpath` for both its
Containerfile and build context. The shared helper and other runners are unchanged.

No result, measurement or completed runtime gate came from this failed invocation.
Initial image tests remain recorded in `VALIDATION.md`; the corrected producing
runner must rebuild and pass the final tests before its serial control matrix.

NEXT MODEL: HIGH
