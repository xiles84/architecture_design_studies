# Progress — audit committed diagram sources for scope and fidelity

**Task:** `task-20260923T235916Z-diagram-fidelity-audit`
**Branch:** `repo/diagram-fidelity-audit`
**Attempt:** `attempt-20260924T005828Z-1414e9` (claim `claim-eb5ce8e7bbc3195b`, epoch 1)
**Capability/role:** HIGH session executing the task / analyst (model `deepseek-flash`)
**Required tag:** `repo/diagram-fidelity-audit-v1`

Documentation-and-migration work over committed sources. No database started, no benchmark lock taken,
no measurement. No measured result, report or signed analysis was edited.

## Decision Log

1. **D3's "never touching the heap" is scoped, not deleted.** It is a true prepared-read result
   (0 heap fetches for 1,756 / 31,293 entries under explicit `VACUUM`) that was stated as a timeless
   property; the ANALYZE-only diagnostic still starts at 31,293 fetches. The corrected note keeps the
   zero and names its preparation, per `d2-d3-mechanisms--gpt-6--2026-09-13.md`. Confidence: High.
2. **D3's "six of twelve queries lose a join" is stale and is corrected.** At the current sixteen-query
   catalogue D3 drops a join in three queries (q02, q08, q15) versus D2. Confidence: High.
3. **D10's corruption count is kept; its "per 30 s" unit is not supported and is removed.** The count
   (3 and 5 wrong rows) is in `reports/20260913-d9-cache-race.md`; no cited artefact names a 30 s
   window. The note now names the report and digest. Confidence: High.
4. **D4's "seven of twelve" is left as written and classified revision-scoped.** The count is exactly
   right for the twelve-query catalogue the SQL header itself names; the current catalogue makes it
   11 of 16. Editing it would remove a correct historical statement. Confidence: High.
5. **D7's "one person's 40 donations" is classified illustrative, not measured** (the dataset averages
   ~22 donations/donor). It is a hypothetical mechanism sketch, not a result; left unchanged.
   Confidence: Medium.
6. **Only study 01 sources were edited, so only study 01 was re-rendered.** The two changed sources
   (`d3`, `d10`) produced changed SVGs; the other 18 are byte-identical, confirming the corrections are
   the only delta. Confidence: High.
7. **`00_overview.puml` was not touched.** It is registered in the book figure manifest
   (`book/assets/figures.json`, task `book-figure-provenance`) and `book/` is a forbidden path here;
   it also needed no correction. Confidence: High.

## What was produced

- `studies/01-charity-tree/diagrams/{d3_flattened_fk,d10_embedded_hybrid_locked}.puml` — corrected
- `studies/01-charity-tree/diagrams/rendered/{d3_flattened_fk,d10_embedded_hybrid_locked}.svg` — re-rendered
- `studies/0{1,2,3}-*/diagrams/FIDELITY.md` — the per-study audit record
- `studies/01-charity-tree/diagrams/RENDERER.md` — refreshed by the renderer
- `CONTEXT.md` — the audit note

## Residual risks / limitations

- The audit is an AI self-audit with no independent model review; classification of mechanism prose
  (e.g. D6's O(history) claim, D9's "highest-QPS" remark) is judgement against the design SQL, not a
  measurement.
- The audit reads the sources at one revision; a later SQL change can make a "current" annotation
  stale, which is why the `FIDELITY.md` records name the revision.
- Sequence-diagram timing labels (e.g. "minutes 45–85") are scenario parameters of the harness phases,
  not measured outcomes.
