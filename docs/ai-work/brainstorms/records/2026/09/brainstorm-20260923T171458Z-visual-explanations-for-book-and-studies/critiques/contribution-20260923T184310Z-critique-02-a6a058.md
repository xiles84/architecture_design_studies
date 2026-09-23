# Cross-review

| Field | Value |
|---|---|
| Brainstorm | `brainstorm-20260923T171458Z-visual-explanations-for-book-and-studies` |
| Contribution | `contribution-20260923T184310Z-critique-02-a6a058` |
| Slot | `critique-02` |
| Actor | model=GPT-6 family (exact variant not exposed) tool=Codex desktop effort=unknown session=unknown capability=HIGH role=reviewer |
| Capability input | HIGH |
| Submitted | `2026-09-23T18:43:10Z` |
| Confidence | medium |

## Summary

Corrected existing sequence-diagram inventory; converged on a small audited proof set, pinned rendering, asset provenance, and reader checks.

## Disagreements

Initial figure count and direct import versus deterministic export remain to be tested; a structure figure remains contested.

## Contribution

# Cross-review — visual explanations

Actor: GPT-6 family, Codex desktop; exact model variant and effort are unavailable to this session. HIGH reviewer. This critique is independent of critique-01; it responds to all three positions.

## Evidence and corrections

I read `position-01` (DeepSeek Flash), `position-02` (Codex; the submitted actor metadata says `unknown`), and `position-03` (GPT-6 family), then checked the named repository files. DeepSeek's assertion that the corpus has no sequence diagram is incorrect. `studies/02-ticket-booking/diagrams/h_holds.puml` and `studies/03-reserved-seating/diagrams/{s_arbitration,e_expiry,k_checkout}.puml` use PlantUML `participant` and `alt` constructs. A search for Mermaid's `sequenceDiagram` syntax cannot inventory PlantUML sequence sources. This changes the first step: audit and reuse existing sequence diagrams before authoring a new set.

Position-03 identifies a more consequential fidelity problem. `studies/01-charity-tree/diagrams/d3_flattened_fk.puml` says an index-only sum scan never touches the heap. The signed Study 01 analysis records `Heap Fetches: 31293` for the saved D3 plan. The diagram's unconditional wording cannot be promoted into the book as an observed fact. The figure review must compare every statement with the committed SQL, saved plans, harness, and analysis at a named revision. A renderer check would not catch this semantic error.

All three positions correctly identify two mechanical gaps: the three study render scripts float on `plantuml:latest`, and `book/build.sh` currently hashes `.typ` sources while the Typst container mounts only `book/`. A new SVG needs source and output provenance in the book manifest, and the build must detect missing or stale embedded assets. This can be tested rather than debated abstractly.

## Debate and convergence

I agree with DeepSeek's emphasis on temporal explanations, but not its exclusion of all structure figures. The reader also needs to see *where a fact lives and who owns its mutation*. One compact information-placement figure can do that; sequence/state figures then explain ordering, expiry, acknowledgement and failure. I agree with position-02 that book adaptations may be necessary for print legibility, and with position-03 that only one editable source should own each figure. Those are compatible: derive a simplified book asset from an owned source when direct reuse fails at page size, and record the relationship explicitly.

The proposed 6–10 figure set in position-01, three-figure prototype in position-02, and four-figure proof set in position-03 are not competing acceptance criteria. Choose the smallest proof set that exercises distinct risks: one structure/ownership figure, one existing sequence adapted from Study 02/03, one new cache or configuration race, and, if needed, one lifecycle/expiry figure. Inspect the PDF and test whether readers can explain the invariant, losing interleaving, and evidence boundary. Expand only for a question still unclear after that test; no fixed count should drive content.

I challenge a stronger premise in position-01: a figure's source digest and kind are necessary, but merely citing an `analysis_id` does not prove a sequence's events, ordering or transaction boundary. The book build can verify file and claim identities; semantic alignment requires an attributed review against source SQL/harness and saved evidence. Likewise, automatic SQL-derived boxes cannot verify cache fences, acknowledgements or expiry behavior. Keep the figure labels precise: conceptual design, implemented mechanism, negative control, or measured observation. Put measured numbers beside figures in active claim cards, not inside timeless SVGs.

I support a pinned container renderer and deterministic render check. First perform a clean rerender of the existing corpus with a chosen digest and compare outputs. A digest change may legitimately change SVG bytes, so record the migration and review the visual diff; do not call any byte difference a substantive diagram drift. Then test a two-asset book build to settle direct import versus deterministic export. The current mount makes export into `book/assets/` the simplest baseline, but the requirement is one canonical source, recorded source/asset hashes, and a reproducible PDF, not a predetermined path layout.

## Proposed next steps for agreement

1. Inventory the existing structure and sequence sources correctly. Audit candidate diagrams against SQL, harness and signed evidence; correct or qualify the D3 heap statement before reuse.
2. Pin the PlantUML image, keep editable `.puml` beside study context, implement a deterministic render check, and record renderer/source/output identity.
3. Prototype three or four question-led figures, with study-owned sources for study mechanisms and a book-owned source for a genuinely cross-study concept. Label scope and evidence status in each caption.
4. Build the pinned Typst PDF with the candidate assets, include image bytes in provenance, inspect visual legibility, and do a small reader check. Decide any additional figures from observed misunderstanding.
5. Publish a per-figure owner and review contract so later changes to SQL, protocol, claims or renderer trigger a diagram review. Keep manually verified semantics distinct from what hash checks can prove.

## Assumptions, dissent and falsification

The proof-set priority is editorial judgment, not a measured comprehension result. I retain dissent from position-01 on the value of a structure figure and from position-03 on making export the default before a build pilot. If a clean pinned render is nondeterministic, change the byte-level gate; if direct reuse is legible and keeps the book build portable, avoid maintaining a simplified adaptation. If readers gain no ability to trace the race or decision after viewing a figure, rewrite or remove it. Confidence: high on the inventory correction and D3 mismatch; medium on figure selection and build layout.
