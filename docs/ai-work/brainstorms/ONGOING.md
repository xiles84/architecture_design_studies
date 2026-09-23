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
- Summary: How should the book map cache strategies and consistency tactics across legacy databases whose schema cannot change and owned databases whose schema can change, including cache-aside, write-through, process-local and shared caches, versions, outbox, soft deletes or tombstones, leases and locks? Which tradeoffs and correctness guarantees are supported by Study 05, which claims about avoiding locks need qualification, and what coverage is still missing?
- Progress: positions 0/3; critiques 0/2; synthesis pending
- Disagreements: none recorded
- Next: `CONTINUE BRAINSTORM brainstorm-20260923T172022Z-cache-strategy-coverage-by-schema-control`
