# Execution result — Study 01 integration EH-01

## TL;DR

- LOW executed the reviewed local integration through candidate `32aec56`.
- Both worktrees were clean; the task-owned integration lock was confirmed; no database
  containers were running.
- Existing Study 01 run provenance, result digests and report-analysis indexes match
  the HIGH-reviewed register. No new benchmark was needed.

Executor: GPT-5 via Codex desktop, LOW role, 2026-09-14; exact effort not exposed.
Handoff: [EH-01](HANDOFF.md). HIGH-reviewed incoming main:
`025b72f184434a256c106f21bf78edefb5641c2d`. Candidate:
`32aec5683a13cd0a220eb0a7e78e68d43bd3b983`.

## Evidence

Before updating `main`, these conditions passed:

```text
task worktree clean; branch study-01/measurement-enhancements
main worktree clean; main 025b72f184434a256c106f21bf78edefb5641c2d
repo/study01-integration-handoff-v1 resolves to 32aec56
Study 01/infrastructure equal the reviewed 0efd760 baseline
incoming Study 02/Study 03/platform equal 025b72f
git diff --check 025b72f -- AGENTS.md CONTEXT.md LESSONS_LEARNED.md docs infra
no unresolved conflict paths
```

The task acquired `ads-run-lock` with holder `Study 01 final local integration` and
worktree `C:/extra/code/architecture_design_studies/.worktrees/study01-v3`. The lock
was absent before acquisition; no containers were listed after acquisition.

The four `run/01-charity-tree/*-v3` tags resolve to their producing commits, and each
generated report contains the registered inputs digest and current final-analysis index:

| Run | Producing commit | Inputs digest |
|---|---|---|
| `20260913T124917Z-v3` | `41a1490de485f2fc905e36f2486842eee0aa2592` | `d2d1cc8380b381f2` |
| `20260913T125342Z-v3` | `c084c73b0cade873406c57fa65ce61f5a26c505a` | `0bc5900a6d7a5d3a` |
| `20260913T172624Z-v3` | `9b7aa48b329ca44e341fe2aedb6d39de8d4cb70d` | `bd95d304cb39e2de` |
| `20260914T104721Z-v3` | `d1b18bf437daf972860261e53b422dc87905cc09` | `9ead901cb74c4183` |

The first local-main update was a fast-forward from `025b72f` to `32aec56`. This closeout
revision is the target for the final fast-forward and the two annotated tags.
No remote was contacted, altered or synchronized.

**Next: HIGH model / high effort — review this evidence and the completed closeout.**
