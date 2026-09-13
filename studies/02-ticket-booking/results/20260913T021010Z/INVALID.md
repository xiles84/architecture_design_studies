# INVALID — aborted at the first cell, provenance not captured

This was the first launch of the `small` matrix (2026-09-13 02:10 UTC). It was stopped by
hand within minutes, during PostgreSQL cell `p1_precreated_lock_first`, because the runner
had failed to read the repository version:

- `git -C /c/extra/...` could not resolve the POSIX path once `infra/lib.sh` disabled Git
  Bash's path conversion, so the manifest records `repo_commit: unknown` and
  `repo_dirty: false` — both wrong (the code was commit `774d258`, clean).
- `run_tag` names `run/02-ticket-booking/20260913T021010Z`, but **no such tag was created**:
  the same git failure prevented it.

The single result file is partial and says so (`"error": "... interrupted by signal: results
are partial"` — the harness's signal handling, working as intended). Nothing here may be
reported or analysed. It is kept rather than deleted, as the project keeps every invalid
run with the reason beside it. The runner now passes git a host-style path and refuses to
start when it cannot read the commit (LESSONS_LEARNED).
