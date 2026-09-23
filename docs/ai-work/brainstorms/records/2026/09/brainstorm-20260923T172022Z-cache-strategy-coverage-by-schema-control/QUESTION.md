# Brainstorm question

**Brainstorm:** `brainstorm-20260923T172022Z-cache-strategy-coverage-by-schema-control`  
**Title:** Cache strategy coverage by schema control  
**Created:** `2026-09-23T17:20:22Z`

## Question

How should the book map cache strategies and consistency tactics across legacy databases whose schema cannot change and owned databases whose schema can change, including cache-aside, write-through, process-local and shared caches, versions, outbox, soft deletes or tombstones, leases and locks? Which tradeoffs and correctness guarantees are supported by Study 05, which claims about avoiding locks need qualification, and what coverage is still missing?

## Decision criteria

- Map current book coverage to Study 05 designs and evidence
- Separate measured findings from untested tactics and proposed gaps
- Compare correctness, database reads, write cost, coordination, failures and operational complexity
- Test when versioning and soft deletes reduce or do not remove the need for locks and fences
- Create no tasks or measurements during the brainstorm

## Evidence packet

- `book/concepts/cache-consistency.typ`
- `book/chapters/major-families.typ`
- `book/chapters/minor-variants.typ`
- `book/chapters/scenarios.typ`
- `book/evidence/v2/claims.json`
- `studies/05-cache-consistency/README.md`
- `studies/05-cache-consistency/HANDOFF-AMENDMENT-01-V2.md`
- `studies/05-cache-consistency/reports/analyses`
- `studies/05-cache-consistency/reports/discussions`
- `studies/05-cache-consistency/sql/README.md`
- `docs/methodology.md`
- `CONTEXT.md`
