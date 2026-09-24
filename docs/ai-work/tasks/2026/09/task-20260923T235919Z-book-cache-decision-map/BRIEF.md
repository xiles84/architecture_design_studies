# Execution brief — Rework the book cache material into the guarantee-first decision map

Read the cache brainstorm conclusion first
(`docs/ai-work/brainstorms/records/2026/09/brainstorm-20260923T172022Z-cache-strategy-coverage-by-schema-control/CONCLUSION.md`),
then `book/concepts/cache-consistency.typ`, `book/chapters/{major-families,minor-variants,scenarios,how-to-choose}.typ`,
the active evidence version, and `studies/05-cache-consistency/README.md`.

## Goal

The book explains the cache gain and the failure mode but offers no selection map, so a reader cannot decide
when the schema is frozen, when writes bypass the adapter, or what each tactic actually guarantees. Replace the
flat pattern discussion with the guarantee-first decision map, with every number on a registered claim.

## Decisions already made (do not re-litigate)

- Flow order: contract → writer-observation completeness → source capabilities (frozen/owned) → topology and
  publication protocol → performance controls separate. Schema control is a branch, not the root.
- The named mechanisms are distinct and must stay distinct: database mutation lock, fill lease, publication
  fence/CAS, source version token, outbox.
- Corrected claim: a cache-side fence is **sufficient**, not proven necessary; write "fence or a demonstrated
  equivalent barrier".
- Corrected claim: a source version not consulted at the hit/publication boundary does not make a relaxed shared
  cache strict (do not write "ownership did not buy consistency" as a general law).
- Cache-aside vs write-through are fill/publication paths, not guarantees, and the study has no head-to-head ratio.
- Consume the evidence version produced by the dependency task; do not edit `book/evidence/`.

## Acceptance criteria

See the task spec. Minimum: the guarantee-first flow is present; the protocol table separates the four evidence
statuses; every lock sentence names its race; the private-memory bypass trap is explicit; relaxed staleness is
defined against the strict comparator; every number resolves to a registered claim and the book rebuilds.

## Constraints

- No measurement; no benchmark lock.
- Never edit another analyst's analysis, report or digest; write your own signed analysis.
- Do not touch `book/build.sh`, `book/assets` or the figure pipeline (separate task).
- Commit explicit paths, tag `repo/book-cache-decision-map-v1`, merge through the queue lifecycle.

## Escalate only if

- a required sentence cannot be supported by any registered claim and the alternative changes a published
  conclusion (that is material, not a local edit).

NEXT MODEL: LOW
