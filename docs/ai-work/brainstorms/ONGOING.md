# Ongoing brainstorms

Generated from immutable brainstorm records; do not hand-edit.

## `brainstorm-20260923T171458Z-visual-explanations-for-book-and-studies` — Visual explanations for book and studies

- State/stage: `ongoing` / `collecting_positions`
- Summary: Pin the PlantUML image and label figures by provenance before adding any, and make the book build fail on an unlabelled or digest-stale figure. Spend the visual budget on 6-10 sequence/state diagrams for the temporal race mechanisms, which the corpus lacks entirely, not on more schema pictures. Gate drift with a --check re-render.
- Progress: positions 1/3; critiques 0/2; synthesis pending
- Disagreements: Anticipated: a floating :latest tag is acceptable because SVGs are committed; diagrams belong only in study docs; sequence diagrams cost more to maintain than they explain.
- Next: `CONTINUE BRAINSTORM brainstorm-20260923T171458Z-visual-explanations-for-book-and-studies`

## `brainstorm-20260923T172022Z-cache-strategy-coverage-by-schema-control` — Cache strategy coverage by schema control

- State/stage: `ongoing` / `collecting_positions`
- Summary: Map caches on two axes: who arbitrates the copy (cache-side fence, DB version token, authoritative read) and whether the schema may change (legacy vs owned), with aside/through, memory/redis and relaxed/strict as parameters. Ownership did not buy consistency (owned 87.81% vs legacy 87.24% wrong reads under relaxed write-through), and the strict arms came from a fence that fires before and after commit. Record soft deletes, outbox-driven invalidation, partitioned process-local caches, hard-TTL expiry and placement as gaps.
- Progress: positions 1/3; critiques 0/2; synthesis pending
- Disagreements: Anticipated: relaxed staleness is permitted so it should not be called wrong; one machine and one trial cannot be generalized; version columns, outboxes and soft deletes are standard and may be recommended.
- Next: `CONTINUE BRAINSTORM brainstorm-20260923T172022Z-cache-strategy-coverage-by-schema-control`
