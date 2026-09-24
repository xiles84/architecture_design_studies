# Execution brief — Add figure provenance to the book build and settle the figure import mechanism

Read the visual-explanations brainstorm conclusion first
(`docs/ai-work/brainstorms/records/2026/09/brainstorm-20260923T171458Z-visual-explanations-for-book-and-studies/CONCLUSION.md`),
then `book/build.sh`, `book/README.md` and `book/dist/build-manifest.json`.

## Goal

Today `book/build.sh` mounts only `book/`, compiles with `--root book`, and hashes `*.typ` files only; the
manifest has no asset keys. A figure would therefore be both unreachable and invisible to provenance. Make
figures first-class build inputs with a manifest and a stale/missing gate, and settle how a study-owned SVG
gets into the book.

## Decisions already made (do not re-litigate)

- One canonical editable source per figure; the book must never hold an independently editable copy.
- The open question is the mechanism, and it is settled by a **pilot, not preference**:
  (a) direct repo-root Typst import, or (b) a generated, read-only `book/assets/` export.
  Criteria: editor-preview portability, clean-tree build, manifest coverage.
- Machine provenance lives in the manifest; reader captions stay short (that is the later figures task).
- The renderer is already pinned by the dependency task; do not change the pin.

## Acceptance criteria

- Figure manifest records canonical source path + revision, renderer digest, output hash, embedded bytes.
- Book source-tree hash covers figure sources and assets.
- Build fails on a missing figure or a stale asset.
- Two-variant pilot executed and documented in `book/FIGURES.md`; the chosen mechanism is stated with its reason.
- Mutation test passes: edit a `.puml`, leave its SVG, build fails.

## Constraints

- No measurement; no benchmark lock. Do not start or stop databases.
- Do not edit chapter/concept prose or the figures themselves (later tasks own those).
- Do not edit `book/evidence/`; the active registry is owned by another task.
- Commit explicit paths, tag `repo/book-figure-provenance-v1`, merge through the queue lifecycle.

## Escalate only if

- neither import mechanism can keep the published book reproducible and editor preview working (then the book
  layout decision is material and belongs to HIGH).

NEXT MODEL: LOW
