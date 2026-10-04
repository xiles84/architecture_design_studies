# Execution brief — Book Edition 2 revision B

Tranche B of the owner-directed revision of *Data Architecture Reference*: the four remaining accepted
visuals from the external review (items 2.2, 2.4, 2.5, 2.6). Tranche A (correctness A–G, Figure 6/7
redesign, glossary terms, 2.1/2.3/2.7/2.8) is already integrated and this task starts from that state.

## Goal

Add four explanatory figures that reduce the state a reader must hold mentally, using the same visual
language and the same Donor / Donation domain as tranche A.

## Decisions already made (do not re-litigate)

- **Registry `v5` is read-only.** No measured number, claim id, trial count, evidence classification,
  confound or gap changes.
- **Every figure is a `conceptual illustration`** with form `structure` or `sequence`, a short caption, a
  text equivalent, no embedded figure number and no embedded rate.
- **2.9 stays a table** (do not add a flowchart), and 2.1/2.3/2.7/2.8 belong to tranche A (already done).
- **One recurring Donor / Donation domain.**
- **Do not touch `book/dist/`.** Build to `book/build.sh --out dist/preview.pdf --manifest
  dist/preview-manifest.json` and delete both before committing.

## New figures

All sources live in `book/assets/sources/` (shared `_style.puml`; solid arrows = synchronous, dashed =
derived/asynchronous), get a `figures.json` entry, are embedded once, and carry a text equivalent.

| id / source | form | content | placed in |
|---|---|---|---|
| `fig-embedded-vs-normalized` / `fig_embedded_vs_normalized.puml` | structure | left: normalized `Donor 1—* Donation` independent rows; right: `Donor { … donations: [ Donation, … ] }` inside the parent. Annotate: whole bounded aggregate reads become local; querying/updating one independent child becomes less convenient; an unbounded child collection is dangerous | `book/chapters/major-families.typ`, Embedded documents family |
| `fig-cache-read-copy` / `fig_cache_read_copy.puml` | structure | `Client → App`, branching to `source DB` (authoritative) and `cache` (non-authoritative read copy). Annotate: authoritative mutation owns source correctness; propagation/publication/hit validation owns copy correctness; the freshness contract defines what a hit may mean | `book/concepts/cache-consistency.typ`, before the race sequences |
| `fig-index-vs-hotrow` / `fig_index_vs_hotrow.puml` | structure | panel 1: query → narrow index path → small subset of rows; panel 2: Writer 1..4 → one row. Make the contention explicit | `book/concepts/indexes-and-hot-rows.typ` |
| `fig-optimistic-vs-pessimistic` / `fig_optimistic_vs_pessimistic.puml` | sequence | optimistic: R1/R2 read version 5; R1 `UPDATE … WHERE version=5` succeeds → version 6; R2 `UPDATE … WHERE version=5` affects 0 rows → re-read/retry. pessimistic: R1 `SELECT … FOR UPDATE` acquires the lock; R2 waits; R1 rechecks the invariant, updates, commits; R2 acquires the lock and rechecks the invariant before deciding. Label that the **invariant recheck**, not the lock, is what rejects an invalidated operation | `book/concepts/concurrency-control.typ` |

## Page-count checkpoint

Record the rendered page count before and after this tranche in `RESULT.md`. If a figure cannot be made
legible at page size without dropping an accepted item, escalate rather than shrink the text.

## Validation

- `book/assets/export.sh --check`, `book/evidence/check.sh`.
- Scratch `book/build.sh --out dist/preview.pdf --manifest dist/preview-manifest.json`; verify 29/29
  claims, 6/6 fonts, digest present, no banner, no unresolved tokens; delete both preview files.
- Rasterise the figure pages and inspect them for overlap, clipping and unreadable labels.
