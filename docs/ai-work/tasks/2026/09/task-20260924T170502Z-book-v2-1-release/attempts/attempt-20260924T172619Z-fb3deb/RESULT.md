# Result — Release the revised Edition 2 book (v2-1)

**Task:** `task-20260924T170502Z-book-v2-1-release`
**Attempt:** `attempt-20260924T172619Z-fb3deb`
**Worker:** model `deepseek-flash`, tool Deep Code CLI, effort unknown, capability HIGH, role executor.
**Commit:** `392e65c` (branch `book/v2-1-release`).
**Tag required:** `repo/data-architecture-book-v2-1`.

## What was delivered

The release artefact for the owner-directed revision, rebuilt from the fully merged tranches A and B
through the official path:

- **`book/dist/data-architecture-reference.pdf`** — 61 pages, PDF sha256
  `f736326150394a7574e94fac7cc4ea154cec94fddee097f64ab6883fe5b35390`.
- **`book/dist/build-manifest.json`** — 29/29 claims indexed, 6/6 fonts embedded, evidence digest in the
  PDF, no unverified banner, no unresolved template tokens, edition present on the page.
- **`book/CHANGELOG.md`** — grouped into technical corrections, clarified claims, added visuals, modified
  diagrams and unresolved questions/gaps.
- **`book/evidence/reviews/20260924-edition-2-revision--deepseek-flash.md`** — the signed self-review,
  item by item against A–G, 2.1–2.9, 3, 4 and 5.
- **`CONTEXT.md`** and **`LESSONS_LEARNED.md`** updated with the revision, its tags, and six reusable
  lessons about figures and the queue worktree.

No measured number, claim id, trial count, evidence classification, confound or gap changed;
`book/evidence/v5/` is untouched.

## Verification

- `book/build.sh` (release path, no override): compiled in the pinned image, wrote the manifest, and
  printed `verified`.
- Independently confirmed from the manifest: pages 61, claims 29/29, fonts 6/6, digest in PDF true,
  banner false, unresolved tokens false, edition on page true.
- `book/assets/export.sh --check` and `book/evidence/check.sh` both pass (run inside `book/build.sh`).

## Decision Log

| # | Decision | Reason | Confidence | Affects |
|---|---|---|---|---|
| D9 | Keep the edition string "Edition 2 (draft)" and distinguish this revision by tag `repo/data-architecture-book-v2-1` | The owner's plan said the revision is a new milestone of Edition 2; changing the edition string would touch the header and title page for no reader benefit | High | — |
| D10 | Tag `v2-1`, not `v2.1` | The queue task-id slug cannot contain a dot, so the milestone uses a hyphen everywhere | High | All future book tags |
| D11 | Commit the regenerated PDF even though it is binary | It is the release artefact the owner asked for, and the repository's convention is to track it | High | — |

## Where this is weak

- **One-model review** across all three tasks; the self-review is disclosed as such.
- **The book is 61 pages**, up from 54; the growth is recorded but a print/reader test was not run.
- **Two figures render below full width** (74% and 72%).
- **Item 2.9 was rejected with reason** and item 4 is partial; both are recorded in the changelog and
  self-review rather than smoothed over.
