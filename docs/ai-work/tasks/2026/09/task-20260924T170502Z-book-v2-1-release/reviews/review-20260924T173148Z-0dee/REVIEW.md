# Review — task-20260924T170502Z-book-v2-1-release

**Reviewer:** HIGH session, model `deepseek-flash`, Deep Code CLI (effort not exposed)
**Verdict:** approve
**Disclosure:** self-review by the session that executed all three tasks; weaker than an independent
review.

## What was checked

| Check | Method | Result |
|---|---|---|
| Artefact rebuilt through the official path | the build's own output | `book/build.sh` wrote dist and printed `verified` (no `--out` override) |
| Pages | `dist/build-manifest.json` | 61 |
| Claims | manifest | 29/29 indexed |
| Fonts | manifest | 6/6 embedded |
| Provenance | manifest | evidence digest in PDF true; unverified banner false; unresolved tokens false; edition on page true |
| Figure layer | `export.sh --check` (inside the build) | 17 registered figures, pinned renderer, no stale asset |
| Source rules | `evidence/check.sh` (inside the build) | OK |
| CHANGELOG present and grouped | read `book/CHANGELOG.md` | technical corrections, clarified claims, added visuals, modified diagrams, unresolved gaps |
| Self-review covers every review item | read the self-review table | A–G, 2.1–2.9, 3, 4 and 5 each marked completed / partial / rejected |
| Context and lessons current | read the diffs | CONTEXT records the revision and tags; LESSONS adds six entries |
| No measured number changed | `git diff` over `book/evidence/v5` and prose numbers | registry untouched |

## Judgement on the Decision Log

1. **D9 (keep the edition string) is accepted.** The tag and provenance distinguish the revision; changing
   the edition text would touch the header for no reader benefit.
2. **D10 (`v2-1`) is accepted** as a mechanical consequence of the task-id slug.
3. **D11 (track the PDF) is accepted** and is the repository's convention.

## Where this task is weak

- **Self-review across the whole revision.** No independent reviewer or reader test exists.
- **61 pages** (from 54); the growth is recorded but not reader-tested.
- **Two figures below full width**, and **item 2.9 rejected**, are recorded in the changelog and
  self-review rather than hidden.
- **The `Source: ...` line in a couple of text equivalents is long**, but it now fits within the page.

## Required tag

`repo/data-architecture-book-v2-1`, created by `queue integrate`.
