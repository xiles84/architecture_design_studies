# Execution brief — Release the revised Edition 2 book (v2-1)

Final task of the owner-directed revision. Tranches A and B are integrated; this task rebuilds the
tracked artefact from fully merged sources and records the change.

## Goal

Produce the release artefact `book/dist/data-architecture-reference.pdf` plus its manifest, a grouped
CHANGELOG, and a final self-review of the external review's items; update `CONTEXT.md` and
`LESSONS_LEARNED.md`; land the `repo/data-architecture-book-v2-1` tag.

## Work items

1. **Rebuild the release artefact.** Run `book/build.sh` (the release path, no `--out` override) so the
   PDF and `book/dist/build-manifest.json` are regenerated with provenance, then commit both.
2. **CHANGELOG.** Write `book/CHANGELOG.md`, grouped into: technical corrections, clarified claims,
   added visuals, modified diagrams, unresolved questions/gaps. Cite figures by id/caption, not by their
   old numbers (numbering shifts as figures are inserted). Record that the tag is `v2-1` (hyphen) rather
   than `v2.1` because the queue task-id slug cannot contain a dot.
3. **Final self-review.** Write `book/evidence/reviews/<YYYYMMDD>-edition-2-revision--deepseek-flash.md`
   stating, item by item, whether each external-review item (A, B, C, D, E, F, G; 2.1–2.9; 3; 4; 5) was
   completed, partially completed, or rejected with reason. Include the "where this is weak" section.
4. **Context and lessons.** Update `CONTEXT.md` (the revision, the merge outcome, the new tag) and
   `LESSONS_LEARNED.md` with any reusable lesson (for example: a sequence figure must give every
   participant named in the prose a lifeline; a fencing token and source-generation validation are
   different mechanisms and must be drawn differently).
5. **Verify and integrate.** Confirm `book/assets/export.sh --check` and `book/evidence/check.sh` pass,
   rasterise and inspect the final pages, merge the task branch into local `main`, and confirm the tag.

## Constraints

- Registry `v5` is read-only; no measured number, claim id, trial count or evidence classification
  may change.
- Podman only; `book/build.sh` takes `ads-book-build-lock`, never the benchmark lock.
- Stage explicit paths; never `git add -A`; never push or pull.
- A figure page that cannot be made legible is an escalation, not a reason to drop an accepted item.
