# Ongoing brainstorms

Generated from immutable brainstorm records; do not hand-edit.

## `brainstorm-20260923T171458Z-visual-explanations-for-book-and-studies` — Visual explanations for book and studies

- State/stage: `ongoing` / `collecting_positions`
- Summary: Use a small set of question-led structure and sequence figures, with detailed study sources and readable book adaptations. Pin rendering and include diagram assets in PDF provenance; validate visual legibility and factual alignment.
- Progress: positions 2/3; critiques 0/2; synthesis pending
- Disagreements: Figure count, direct SVG reuse versus book adaptations, and PlantUML versus an alternative renderer remain open until PDF and reader checks.
- Next: `CONTINUE BRAINSTORM brainstorm-20260923T171458Z-visual-explanations-for-book-and-studies`

## `brainstorm-20260923T172022Z-cache-strategy-coverage-by-schema-control` — Cache strategy coverage by schema control

- State/stage: `ongoing` / `collecting_positions`
- Summary: Map caches on two axes: who arbitrates the copy (cache-side fence, DB version token, authoritative read) and whether the schema may change (legacy vs owned), with aside/through, memory/redis and relaxed/strict as parameters. Ownership did not buy consistency (owned 87.81% vs legacy 87.24% wrong reads under relaxed write-through), and the strict arms came from a fence that fires before and after commit. Record soft deletes, outbox-driven invalidation, partitioned process-local caches, hard-TTL expiry and placement as gaps.
- Progress: positions 1/3; critiques 0/2; synthesis pending
- Disagreements: Anticipated: relaxed staleness is permitted so it should not be called wrong; one machine and one trial cannot be generalized; version columns, outboxes and soft deletes are standard and may be recommended.
- Next: `CONTINUE BRAINSTORM brainstorm-20260923T172022Z-cache-strategy-coverage-by-schema-control`
