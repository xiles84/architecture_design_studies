# Cross-review

| Field | Value |
|---|---|
| Brainstorm | `brainstorm-20260923T171458Z-visual-explanations-for-book-and-studies` |
| Contribution | `contribution-20260923T184032Z-critique-01-8dd74c` |
| Slot | `critique-01` |
| Actor | model=gpt-5.6-sol tool=Codex desktop effort=unknown session=visual-critique-alt-20260923 capability=HIGH role=reviewer |
| Capability input | HIGH |
| Submitted | `2026-09-23T18:40:32Z` |
| Confidence | high |

## Summary

Separate diagram form from evidence status; audit semantic drift, then pilot two figures before expanding.

## Disagreements

Direct reuse versus generated export, final figure count, byte-identity gating, and reader benefit remain open; position-01's zero-sequence premise is falsified.

## Contribution

# Cross-review

| Field | Value |
|---|---|
| Brainstorm | `brainstorm-20260923T171458Z-visual-explanations-for-book-and-studies` |
| Slot | `critique-01` |
| Claim | `brainstorm-claim-48f2d156dcbef883`, epoch 1 |
| Actor | model=`gpt-5.6-sol` tool=`Codex desktop` effort=`unknown` session=`visual-critique-alt-20260923` capability=`HIGH` role=`reviewer` |
| Capability input | `HIGH` |
| Confidence | high on the workflow and factual corrections; medium on reader benefit until tested |

## Evidence examined

- The brainstorm question, evidence packet, `CRITIQUE.md`, and all three submitted positions.
- The book entry point, build script, visual configuration, all chapter and concept sources, and the active v2 claims registry. The book currently embeds no diagrams; `book/build.sh` hashes only `.typ` sources and builds with `--root book`.
- Every PlantUML source/render script in Studies 01–03 and the diagram links in their READMEs. Studies 02 and 03 already contain sequence diagrams (`h_holds.puml`, `s_arbitration.puml`, `e_expiry.puml`, and `k_checkout.puml`); the render scripts float on `docker.io/plantuml/plantuml:latest`.
- The Study 04 and Study 05 READMEs, which describe timing-sensitive mechanisms but have no diagram directories.
- The D3 diagram's statement that its index-only scan never touches the heap and the saved Study 01 analysis that records 31,293 heap fetches for that earlier plan. Later controlled plans can record zero, so the defect is the diagram's unconditional wording and missing revision scope.
- `docs/methodology.md`, including the distinction between measured facts, signed interpretation, and reproducible artefacts. It has no current diagram evidence contract.

## Agreements

All positions converge on the durable core:

1. Use a small curated set, selected by a concrete reader question, rather than one diagram per design.
2. Prioritize temporal mechanisms where event order, transaction boundaries, expiry, acknowledgement, or stale publication are hard to reconstruct from prose.
3. Keep PlantUML as the initial language because the repository already uses it, but pin the renderer and render in Podman.
4. Treat study-owned sources as canonical for study mechanisms, include figure bytes in book provenance, reject stale rendered assets, and inspect the compiled PDF at final page size.
5. Label illustrations so that a diagram cannot masquerade as measurement, and attach numeric statements only to active v2 claims.
6. Audit existing diagrams before reuse. The D3 contradiction proves that source-to-SVG reproducibility alone does not prove factual fidelity.

These claims survive comparison because they address observed repository gaps: a floating renderer, a `.typ`-only book source hash, no book figures, no Study 04/05 diagrams, and at least one stale or overbroad diagram statement.

## Challenges

### Position 01

- Its central inventory claim is false: the corpus does not have zero sequence diagrams. Studies 02 and 03 already contain at least four PlantUML sequence sources. That weakens the proposed justification for spending the entire visual budget on 6–10 new sequence/state figures. The real gap is selective coverage: Study 04/05 and the book lack diagrams, while useful Study 02/03 sequences already exist and should be audited before replacement.
- "No new diagram before pinning" is too rigid. A disposable, uncommitted pilot rendered with a locally resolved image digest is the cheapest way to choose a viable page design and pin. Promotion into committed documentation should require the pin; exploration need not wait.
- The proposed `kind = illustration | mechanism | measured` taxonomy conflates two independent properties. `mechanism` describes content; `measured` describes evidence. An implemented sequence may be mechanically faithful but never observed as that exact interleaving. A single enum encourages accidental evidence upgrades.
- Deriving only a small book schema set from SQL does not solve the harder fidelity problem. SQL cannot establish application ordering, acknowledgements, retry bounds, clocks, or cache fencing. A provenance check proves lineage, not semantic truth.

### Position 02

- It correctly identifies existing sequences, but its initial set remains broad: structure, topology, four race mechanisms, and a decision path could become seven or more figures before the visual contract is proven. The claimed comprehension benefit is still a hypothesis.
- "Measured mechanism" is ambiguous. Measurements can show outcomes and plans; they rarely prove every arrow in a sequence. Figures should distinguish an implemented contract from an observed result, and captions should avoid causal language unless a cited analysis supports it.
- A topology figure is high risk in this corpus. The book already warns that colocated laptop containers do not verify data placement or a real network. A topology drawing is useful only if it visually separates logical topology, configured replication, verified placement, and unmeasured network behavior.

### Position 03

- The four-item proof set is directionally good but is not four figures as written: Study 04 has two panels, Study 02/03 combines an expiry sequence and a refund timeline, and the information-placement map can easily become a crowded six-family wall chart. The pilot needs an explicit page/figure count and pass criteria.
- Deterministic export to `book/assets/` is plausible under the current `--root book`, but committing copied SVGs introduces a second byte location. The proposal needs the export to be entirely generated, with no independently editable book copy, and must prove that editor preview and clean-tree builds work. Direct reuse should remain a live alternative until that pilot runs.
- Requiring source commit/tag, source path, kind, question, invariant, design IDs, and claim links in every human caption risks unreadable captions. Machine provenance belongs in a manifest or adjacent source record; reader captions should state the question, scope, and evidence boundary concisely.

### Shared omissions

- None of the positions separates **diagram form** (structure, sequence, state) from **fidelity/evidence status**. This is the most important contract to settle.
- None gives a sufficient accessible-text rule. A visual needs a concise prose equivalent naming the actors, state change, and failure/success path; color cannot carry meaning alone. A caption and grayscale check do not necessarily make the mechanism understandable to a screen-reader user.
- All infer comprehension benefit without a baseline. A large reader study is unnecessary, but publishing several figures before a two-figure comprehension check would spend maintenance budget before testing the premise.
- Output hash equality is a useful drift signal but may be brittle across renderer metadata or serialization changes. The gate should first test repeatability with the selected digest; if bytes are unstable, compare normalized SVG or render a fixed raster for visual regression while retaining source/output hashes for provenance.

## Better option or combination

Adopt a two-stage **audit then two-figure pilot**, followed by evidence-based expansion.

### 1. Define two orthogonal labels

Each figure declares:

- **Form:** structure, sequence, or state.
- **Fidelity/evidence:** conceptual illustration; implemented design contract; negative control; observed result. The book's existing direct/analogy/gap vocabulary remains a separate evidence-transfer label.

An `observed result` figure must cite an active claim or report cell. An `implemented design contract` cites SQL/harness plus a revision and may not use observed/performance wording. A conceptual illustration carries no design IDs unless it explicitly says they are examples.

Keep hashes, renderer digest, canonical source, revision, and output path in a machine-readable figure manifest. Keep the reader caption to: the question answered, relevant design scope, and visible evidence label. Provide a short adjacent text equivalent for accessibility.

### 2. Audit before adding

Inventory the existing `.puml` files and classify each as current, revision-scoped, stale, or unsupported. Check every factual annotation against SQL/harness and any cited plan at the named revision. Fixing or explicitly retiring the D3 overstatement is the first test that the review process catches semantic drift rather than merely byte drift.

### 3. Pilot exactly two figures

- **Study 05 stale-fill sequence:** highest-value new temporal gap; show old fill, database commit, acknowledgement, fence/invalidation, and accepted/rejected publication. Mark branches as implemented contracts or negative controls rather than measurements.
- **Book information-placement map:** one conceptual structure figure placed once and cross-referenced from `how-to-choose` and `information-placement`.

This pair tests both new-source ownership and book-only concept ownership, sequence and structure forms, and study-to-book reuse without committing to a full catalog.

Use PlantUML in a pinned Podman image. Generate a book asset from the canonical source because the current build root is `book/`, but make the copy read-only/generated and validate it against a manifest. Also run a direct-import experiment with repo-root Typst in the pilot; choose export only if it preserves preview/build portability better. This resolves the disputed import strategy with evidence instead of preference.

### 4. Expand only after gates pass

The pilot passes when:

- two clean renders are reproducible under the selected comparison method;
- changing `.puml` without export fails the book check;
- the manifest covers source, renderer, output, and embedded bytes;
- a reviewer maps each arrow/entity to the named design revision;
- the compiled PDF is legible at ordinary zoom and in grayscale, with no clipping;
- two unfamiliar readers can state the intended mechanism and evidence boundary more accurately with the figure than with the prose alone;
- the adjacent text equivalent conveys the same essential explanation without the image.

Then add, in likely order: Study 04 optimistic/pessimistic concurrency, Study 03 expiry (reuse audited existing source), and Study 02 refund answerability. Defer a topology figure until it can show its evidence layers without implying verified physical placement.

## Remaining disagreements and evidence gaps

- **Direct reuse versus generated export:** unresolved until the two build-root variants are tested. Current `--root book` favors export, while one canonical editable source favors direct reuse.
- **Final figure count:** unresolved. The positions range from a four-ish proof set to 6–10 temporal figures. Expansion should follow demonstrated reader value, not an upfront quota.
- **Byte identity as the drift gate:** unresolved until the pinned renderer is tested twice. Normalized or visual comparison may be required.
- **Reader benefit:** no current evidence. The two-reader task check is deliberately small and only decides whether each pilot figure clarifies its stated question.
- **PlantUML long-term:** preferred for the pilot because it is already present, not permanently mandated. A failure on print legibility, accessibility, or deterministic rendering reopens Typst-native or another pinned text renderer.

## Falsification tests

1. **Inventory check, cheapest:** count and open existing sequence sources. This already falsifies Position 01's zero-sequence premise and reduces the initial authoring scope.
2. **Semantic audit, cheapest fidelity test:** map every D3 diagram claim to the named revision's SQL and saved plan. If a reviewed manifest still accepts "never touching the heap" without revision qualification, the proposed contract is insufficient.
3. **Two render/build variants:** compile the same two pilot figures using direct repo-root import and generated `book/assets` export. Compare preview behavior, clean-tree build, manifest coverage, and stale-source failure. Choose the simpler passing path.
4. **Mutation test:** alter one `.puml` line without updating the SVG/export. A green build falsifies the provenance gate.
5. **Evidence-boundary test:** ask two unfamiliar readers whether each arrow is conceptual, implemented, or experimentally observed. Any systematic misclassification requires a label/caption redesign before more figures are added.
6. **Accessibility test:** hide the image and ask whether the adjacent text still communicates actors, order, invariant, and failure branch. If not, the figure would create an information-only visual dependency.
