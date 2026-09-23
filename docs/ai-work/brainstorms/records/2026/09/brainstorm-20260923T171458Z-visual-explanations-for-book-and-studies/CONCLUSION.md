# Brainstorm conclusion

| Field | Value |
|---|---|
| Brainstorm | `brainstorm-20260923T171458Z-visual-explanations-for-book-and-studies` |
| Contribution | `contribution-20260923T233610Z-synthesis-01-ac85f2` |
| Slot | `synthesis-01` |
| Actor | model=deepseek-flash tool=unknown effort=unknown session=session-20260923T174300Z-deepseek-flash capability=HIGH role=integrator |
| Capability input | HIGH |
| Submitted | `2026-09-23T23:36:10Z` |
| Confidence | high |

## Summary

Adopt a small audited provenance-gated visual layer with one canonical PlantUML source per figure, pinned rendering, and asset bytes in the book build. Audit the 30 existing sources first (the D3 heap wording is the first correction), label form and evidence status separately, and pilot 3-4 figures before expanding.

## Disagreements

Preserved: whether one book-level structure figure earns its maintenance; byte-identity versus normalized drift gate; direct import versus generated book/assets export; how much caption detail.

## Contribution

## Decision

Adopt a small, audited, provenance-gated visual layer for the book and the studies, built from
**one canonical PlantUML source per figure**, with rendering pinned in Podman and every figure's
bytes covered by the book build.

Order of work:

1. **Audit before authoring.** Classify all 30 committed `.puml` sources — 26 structural, 4 temporal
   (`studies/02-ticket-booking/diagrams/h_holds.puml`,
   `studies/03-reserved-seating/diagrams/{s_arbitration,e_expiry,k_checkout}.puml`) — against the
   SQL, harness and saved plans at a named revision. The first concrete correction is
   `studies/01-charity-tree/diagrams/d3_flattened_fk.puml:44` ("index-only scan, never touching the
   heap").
2. **Pin and gate.** `PLANTUML_IMAGE` in `infra/versions.env`, one shared render helper, and asset
   bytes inside the book manifest with a build that fails on a missing or stale asset.
3. **Label two independent properties, never one enum:** *form* (structure | sequence | state) ×
   *evidence status* (conceptual illustration | implemented design contract | negative control |
   observed result). An `observed result` figure must cite an active v2 claim; an implemented
   contract cites SQL/harness plus a revision and uses no observed/performance wording.
4. **Pilot 3–4 figures**, then expand only on demonstrated reader benefit. No figure quota.

## Why

Against the four decision criteria:

- **Comprehension.** The corpus's hard material is temporal (sell-out arbitration, hold expiry and
  confirmation, the cache stale-fill race). Four sequence sources already exist for Studies 02–03;
  none exist for Studies 04–09, and the book (`book/`, Typst root `book/`) contains no figure at
  all. The gap is selective coverage, not absence.
- **Faithful and labelled.** The D3 case proves a source-faithful, byte-reproducible figure can still
  overstate. The record: under explicit vacuum preparation the covering plan visits 31,293 or 1,756
  index entries with **zero heap fetches** and D3 does too, while D15 trial 1 *starts* at 31,293 heap
  fetches and ends at 218 as the visibility map updates
  (`studies/01-charity-tree/reports/discussions/d2-d3-mechanisms--gpt-6--2026-09-13.md`). The defect
  is not a false measurement; it is an unqualified, preparation-dependent result presented as a
  timeless property. Hash and digest gates cannot catch that — only semantic review at a revision can.
- **Reproducible in Podman.** All three `render.sh` scripts use `docker.io/plantuml/plantuml:latest`,
  which is not in `infra/versions.env`, so two renders of the same source are not provably the same
  picture. The book build (`book/build.sh`) mounts only `book/`, compiles with `--root book`, and
  computes `source_tree_hash_sha256` over `*.typ` files only; `build-manifest.json` has no asset
  keys. A figure is therefore both unreachable and invisible to provenance today.
- **No duplication.** Study diagrams keep their per-design job beside the design tables that embed
  them; the book carries only cross-study concept/mechanism figures, citing the study source rather
  than re-drawing it.

**Falsified input, recorded.** Position-01 claimed the corpus has "zero sequence diagrams" from a
search for Mermaid's `sequenceDiagram` keyword; both critiques independently corrected it —
PlantUML sequence sources need only `participant`/`->`/`alt`. The conclusion therefore starts with
audit-and-reuse, not author-new.

## Rejected alternatives

- **Diagram every design, or auto-generate all figures from SQL.** SQL cannot express application
  order, acknowledgement, lease expiry, retries or cache fencing, and derivation proves lineage, not
  semantic truth (the D3 test). Retained only as an optional schema-drift check for a small set.
- **Reuse every existing study SVG verbatim in the book.** Cheapest, but the D3 wording shows
  unrevised reuse can import an overstatement, and dense study figures are untested at page size.
  Remains a live option for an online appendix after audit and print review.
- **Switch to Mermaid or Typst-native figures now.** Rejected for the pilot because PlantUML already
  expresses the existing corpus and one visual language is worth more than a marginally simpler
  single build; reopened only if a pinned pilot fails on legibility, accessibility or determinism.
- **One crowded six-family wall chart.** Rejected; the placement/ownership view is one figure,
  placed once and cross-referenced.

## Remaining dissent

Preserved material objections:

1. **Whether one book-level structure/placement figure earns its maintenance.** Position-01 excluded
   it; positions 02–03 and both critiques want it. Recorded resolution: the pilot includes exactly
   one, and it is dropped if the reader check shows no gain. This is an empirical question, not a
   settled argument.
2. **Byte-identity as the drift gate.** If the pinned renderer is not byte-stable, switch to
   normalized/visual comparison while keeping source and output hashes for provenance.
3. **Direct repo-root import versus generated `book/assets/` export.** Both plausible; settle by a
   two-variant pilot judged on preview portability, clean-tree build and manifest coverage — not by
   preference.
4. **Caption detail.** Machine provenance belongs in the figure manifest; the human caption should
   carry only the question, the design scope and the evidence label, or captions become unreadable.

## Assumptions, gaps, and transfer risks

- Reader comprehension benefit is a **hypothesis with no baseline**; the reader check is deliberately
  small and only decides whether each figure clarifies its stated question.
- No compiled-PDF legibility, grayscale or accessibility evidence exists yet; the adjacent text
  equivalent is required so the explanation is not image-only.
- The D3 wording must be corrected or scoped before any reuse; the same audit applies to the other 29
  sources, and one may retire rather than be fixed.
- `docs/methodology.md` currently has **no rule for figures or diagrams**; until one is added, nothing
  in the methodology prevents an illustration from reading as evidence.
- Nothing here changes a measured result or any book number.

## Falsification and review triggers

Reopen or supersede this conclusion if:

1. two renders under the chosen pinned image are materially nondeterministic after a fair test (change
   the gate, not the pin);
2. the mutation test — edit a `.puml`, leave the SVG stale — still produces a green, unchanged
   manifest (the provenance gate is insufficient);
3. readers systematically misclassify a figure's evidence status (label/caption redesign required);
4. both direct import and generated export fail the build/preview pilot; or
5. the D3 audit shows the wording cannot be scoped without rewriting the diagram's source.

## Task-ready next steps

Not authorized by this conclusion; the owner must separately issue
`CREATE TASKS FROM BRAINSTORM brainstorm-20260923T171458Z-visual-explanations-for-book-and-studies`.

1. Pin PlantUML in `infra/versions.env`, move one render helper into `infra/`, and have the three
   study scripts call it; record the resolved image digest beside the rendered output.
2. Audit all 30 `.puml` sources against SQL/harness/saved plans at a named revision and classify each
   as current, revision-scoped, stale or unsupported; correct or qualify the D3 heap wording first.
3. Extend book provenance with a figure manifest (canonical source, revision, renderer digest, output
   hash, embedded bytes) and a build failure on missing or stale assets; run the two-variant
   import/export pilot.
4. Author the 3–4 figure proof set with form × evidence-status labels and adjacent text equivalents,
   then run the two-reader check.
5. Add the missing figure-evidence rule to `docs/methodology.md`.

**Synthesizer disclosure.** This synthesis was written by the session that authored `position-01` in
this record (`deepseek-flash`, HIGH, `session-20260923T174300Z-deepseek-flash`). It read all three
positions and both critiques before writing. Both critiques independently falsified a claim in
position-01, and that correction is carried in the decision; the owner should nonetheless weigh the
reduced independence when reviewing.
