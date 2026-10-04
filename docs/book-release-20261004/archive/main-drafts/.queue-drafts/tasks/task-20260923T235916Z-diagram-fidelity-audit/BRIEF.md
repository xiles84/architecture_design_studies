# Execution brief — Audit committed diagram sources for scope and fidelity; correct the D3 wording

Read the visual-explanations brainstorm conclusion first
(`docs/ai-work/brainstorms/records/2026/09/brainstorm-20260923T171458Z-visual-explanations-for-book-and-studies/CONCLUSION.md`),
then the three study `diagrams/` directories and the study 01 analyses/discussions they cite.

## Goal

A figure can be byte-reproducible and source-faithful and still overstate. The committed corpus has a
concrete example: `studies/01-charity-tree/diagrams/d3_flattened_fk.puml` says an index-only sum scan is
"never touching the heap", while the saved D3/D15 record shows the covering plan at zero heap fetches only
under preparation, with D15 trial 1 starting at 31,293 and falling to 218 as the visibility map updates
(`studies/01-charity-tree/reports/discussions/d2-d3-mechanisms--gpt-6--2026-09-13.md`). Audit every committed
diagram source at a named revision and correct or scope what cannot be supported.

## Decisions already made (do not re-litigate)

- The defect is an **unscoped conditional result**, not a false measurement. Do not "fix" it by deleting a
  correct measured outcome; scope it to preparation and revision.
- Do not re-draw diagrams for style. This task corrects facts and scope only.
- Counts for orientation: 30 sources total, 26 structural, 4 sequence
  (`studies/02-ticket-booking/diagrams/h_holds.puml`, `studies/03-reserved-seating/diagrams/{s_arbitration,e_expiry,k_checkout}.puml`).
- The renderer is pinned by the dependency task; render through the shared helper.

## Acceptance criteria

- Every factual annotation in the 30 sources is classified current / revision-scoped / stale / unsupported,
  with the SQL, harness or plan it was checked against and the revision.
- The D3 wording is corrected or scoped as described above.
- Corrected sources are re-rendered with the pinned helper and the SVGs committed.
- The audit record names, per changed diagram, the evidence used.
- No measured result, report or signed analysis is altered.

## Constraints

- No measurement; no benchmark lock. Do not start or stop databases.
- Never edit another analyst's report or analysis; if a diagram conflicts with one, qualify the diagram.
- Commit explicit paths, tag `repo/diagram-fidelity-audit-v1`, merge through the queue lifecycle.

## Escalate only if

- a diagram states something that contradicts a signed analysis in a way that cannot be resolved by scoping
  the diagram (e.g. the analysis itself is wrong) — that is a material finding, not a local edit.

NEXT MODEL: LOW
