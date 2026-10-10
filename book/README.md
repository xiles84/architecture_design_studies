# Data Architecture Reference

> **Current evidence:** the registry version is selected by the book build configuration. The active evidence package maintains the claim registry and coverage statistics. Preview-specific sections below document historical milestones.

The book has two reader parts:

- **Theory:** six topic groups and their strategy lessons with diagrams, selection boundaries and costs.
- **Study cases:** five business cases, theory-linked variations and compact evidence summaries.

The appendix holds result-reading guidance, reproduction, glossary and full claim cards.
Detailed evidence is kept out of the teaching flow.

## Source layout

| Path | Role |
|---|---|
| `main.typ` | Front matter, hierarchical contents, theory, cases and appendix |
| `theory/*.typ` | Schemas, SELECTs, concurrent writes/time, caches, queues and placement |
| `cases/*.typ` | Charity, ticket booking, reserved seating, configuration and cache consistency |
| `lib/config.typ`, `lib/teaching.typ` | Shared visual system and repeated short lesson template |
| `lib/evidence.typ` | Linked badges and full appendix cards |
| `STYLE.md` | Owner's structure, writing and visual contract |
| `assets/` | Canonical figures, generated SVGs and provenance |
| `evidence/` | Frozen predecessor packages, active package and source gates |
| `chapters/`, `concepts/` | Earlier source preserved; excluded from the new reading path |
| `dist/` | Canonical generated PDF and build manifest |

## Build

```bash
book/build.sh
book/evidence/check.sh
```

Everything builds in pinned Podman images. The manifest records the source commit,
dirty state, registry digest, renderer/toolchain and PDF hash. Inspect rendered pages
before accepting a layout change.

## Edition 3 preview 2

Owner clarification on 2026-10-04 replaces the mixed pilot/older-book layout with
Theory → topic → strategy and Study cases → business → variations → thrive/die situations.
The original preview remains under `repo/data-architecture-book-v3-preview-1`.
This revision's milestone is `repo/data-architecture-book-v3-preview-2`.

The active evidence package is unchanged. Queues are theory without a measured case;
new SELECT results are not promoted to book claims. Study gaps remain visible.
The full source/check/render receipt is in `docs/book-theory-cases-20261004/VALIDATION.md`.
