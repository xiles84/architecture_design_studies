# Independent position

| Field | Value |
|---|---|
| Brainstorm | `brainstorm-20260923T171458Z-visual-explanations-for-book-and-studies` |
| Contribution | `contribution-20260923T180208Z-position-02-749e83` |
| Slot | `position-02` |
| Actor | model=unknown tool=Codex desktop effort=unknown session=unknown capability=HIGH role=analyst |
| Capability input | HIGH |
| Submitted | `2026-09-23T18:02:08Z` |
| Confidence | medium |

## Summary

Use a small set of question-led structure and sequence figures, with detailed study sources and readable book adaptations. Pin rendering and include diagram assets in PDF provenance; validate visual legibility and factual alignment.

## Disagreements

Figure count, direct SVG reuse versus book adaptations, and PlantUML versus an alternative renderer remain open until PDF and reader checks.

## Contribution

## Evidence examined

- `book/main.typ`, `book/README.md`, `book/chapters/{how-to-choose,major-families,topologies,scenarios}.typ`, and `book/build.sh`: the released book has a common visual style and evidence cards, but the Typst source includes no diagrams. The build hashes `.typ` files and the evidence registry; its source hash does not currently include image assets.
- `studies/01-charity-tree/diagrams/`, `studies/02-ticket-booking/diagrams/`, and `studies/03-reserved-seating/diagrams/`: Study 01 has variant structure diagrams; Studies 02 and 03 have both structure and sequence diagrams. `s_arbitration.puml` shows why an interleaving matters, while `p_precreated.puml` places several variants in one dense schema figure. Their render scripts use a containerized PlantUML `:latest` image.
- `studies/04-configuration-portal/README.md` and `studies/05-cache-consistency/README.md`: both describe versioned state, concurrent actions and multiple designs in prose and tables, with no diagram directory in either study.
- `docs/methodology.md` and `book/evidence/v2/claims.json`: correctness and evidence scope govern any numeric or causal label on a figure.

## Recommendation

Add a small, deliberate visual layer to the book, backed by more detailed canonical diagrams in the relevant studies. Choose each figure by the reader question it answers:

1. **Structure:** a compact diagram for the derive/index/rollup/rolldown/embed/cache choices and one for single-node versus replicated placement. Distinguish *logical* relationships from verified *physical* placement.
2. **Time and concurrency:** sequence diagrams for the races that prose hides: two buyers and one seat; hold expiry versus confirmation; a cache fill racing a committed write; a configuration publication and revision check. Put the bad interleaving beside the fence or arbitration that changes the outcome, with the negative control clearly marked.
3. **Reader decisions:** one visual path from required answer and freshness contract to a candidate family and its evidence card. It should route to the text, not encode a performance ranking that the corpus does not support.

Reuse the existing Study 02/03 diagrams as technical sources, but make book versions legible at the printed page size. Create detailed Study 04/05 diagrams first where no source exists, then derive the few book figures worth carrying. A figure needs a caption naming its question, design IDs, source diagram, and whether it is a measured mechanism, an explanatory illustration, an analogy, or an open gap. Performance numbers must point to an active v2 claim ID; a sequence diagram is not itself evidence that a design was measured or correct.

Keep editable `.puml` sources and rendered SVG together. Pin the PlantUML image by immutable version or digest, use one shared containerized render command and style, and make the book build manifest hash every included diagram source and rendered asset. Add a check for stale SVGs and missing figure links, then visually inspect the assembled PDF for line weight, labels, contrast, clipping, and page flow. No host Java or design tool should be required.

## Reasoning and trade-offs

The repository already demonstrates that PlantUML sequence diagrams can explain a concurrency failure precisely (`s_arbitration.puml`); the book does not yet carry that visual language. A focused set improves explanations without making every design a second prose chapter. Study figures may be too dense for the book, so a simplified book view costs editorial maintenance. Source links and a shared design vocabulary make that cost explicit. Pinning the renderer and extending the manifest are prerequisites to calling the PDF visually reproducible, because the current scripts float on `:latest` and the book source hash covers only Typst files.

## Assumptions and evidence gaps

The inventory establishes where diagrams exist, not that a particular new figure improves comprehension. Study 04/05 mechanisms and newer results may exceed the v1 book's active claims; those figures must label evidence status at publication. A topology picture cannot turn colocated laptop containers into proof of data placement or a real network experiment.

## Alternatives

Diagram every design or auto-generate all figures from SQL. That gives apparent completeness but high upkeep and misses application order, retries, and acknowledgement boundaries. Reusing dense study SVGs unchanged in the book is cheaper, but likely harms readability at book scale. These remain reasonable for an online appendix, subject to visual review.

## Anticipated disagreements

Some readers may prefer more figures than this initial set; the figure budget should expand when a distinct question remains hard to answer. Others may prefer Mermaid in Markdown; PlantUML is already used here and can express the existing sequence and schema sources, but an equivalent pinned renderer could win if it better survives the book build and accessibility checks.

## Falsification tests

Prototype three representative figures (structure, race, cache) in a draft PDF. Ask readers to identify the changed design decision, the losing interleaving, and the evidence boundary from the figure and caption. If they cannot, revise or remove the figure. Confirm every line against the committed SQL/harness and evidence claim, render twice from a clean tree with identical output hashes, and inspect the PDF at ordinary viewing and print sizes. Failure on any of those checks weakens this recommendation.

**Confidence:** high on the coverage and reproducibility gaps; medium on the proposed figure selection until reader and PDF review.
