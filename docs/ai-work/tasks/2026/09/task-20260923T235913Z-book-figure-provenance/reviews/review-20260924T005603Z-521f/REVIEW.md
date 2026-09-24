# Review — figure provenance for the book build

**Task:** `task-20260923T235913Z-book-figure-provenance`
**Reviewer:** HIGH session, model `deepseek-flash` (Deep Code CLI), capability HIGH, role reviewer
**Verdict:** **approved for integration** — no rework required
**Required tag:** `repo/book-figure-provenance-v1`

## Decision Log reviewed first

| # | LOW decision | Verdict |
|---|---|---|
| 1 | mechanism (b) generated export, chosen by pilot | accepted — pilot evidence shows (a) fails under `--root book`; (b) compiles under both |
| 2 | authored registry + generated manifest inside `book/assets/` | accepted — matches the acceptance ("figure manifest") and owned paths |
| 3 | canonical source stays a study `.puml`; no book copy | accepted — the non-duplication rule |
| 4 | gate re-renders with the pinned renderer and compares bytes | accepted — stronger than a stored hash; three negative controls fail |
| 5 | `build.sh` gates before compiling and widens the source-tree hash | accepted — manifests and hash both cover the figure layer |
| 6 | pilot figure `00_overview`, not the D3 source the audit must fix | accepted — avoids a self-inflicted staleness after the audit task |
| 7 | no full `book/build.sh` run | accepted — a full build would rewrite the released dist PDF; the clean-tree criterion is proven by the `--root book` compile of `main.typ` |

No decision materially degraded correctness, requirements or maintainability.

## Independent checks re-run by the reviewer

```text
book/assets/export.sh --check                       exit 0, "verified (1 figure(s))"
figure-manifest.json has source+asset+renderer keys verified for the pilot entry
book/build.sh contains the export.sh --check gate    yes, before the compile
mutation via book/build.sh (edited .puml)            exit 1, dist untouched
main.typ compile, book/ mounted only, --root book    exit 0
```

Acceptance criteria:

- Manifest records source path + revision, renderer digest, output hash, embedded bytes — met.
- Book source-tree hash covers figure sources and embedded assets — met.
- Build fails on a missing figure or stale asset — met (and demonstrated).
- Two-variant pilot executed and documented; chosen mechanism stated — met (`book/FIGURES.md`).
- Mutation test passes (edit `.puml`, leave SVG, build fails) — met, through `book/build.sh`.

## Findings (non-blocking)

1. A registered source that changes after the audit task must be re-exported; the build error names
   the fix. Recorded by the executor.
2. The check re-renders a whole diagrams directory per distinct source directory; fine for the pilot,
   worth caching when the proof set grows.
3. `book/dist` is not rebuilt, so the committed manifest lacks the new keys until the next build.

## Verdict

Approved for integration. The figure layer has real source/renderer/byte provenance and a build gate
that fails loudly on the specified defects.
