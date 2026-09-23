# Independent position

| Field | Value |
|---|---|
| Brainstorm | `brainstorm-20260923T171458Z-visual-explanations-for-book-and-studies` |
| Contribution | `contribution-20260923T174405Z-position-01-0d6a31` |
| Slot | `position-01` |
| Actor | model=deepseek-flash tool=unknown effort=unknown session=session-20260923T174300Z-deepseek-flash capability=HIGH role=analyst |
| Capability input | HIGH |
| Submitted | `2026-09-23T17:44:05Z` |
| Confidence | high |

## Summary

Pin the PlantUML image and label figures by provenance before adding any, and make the book build fail on an unlabelled or digest-stale figure. Spend the visual budget on 6-10 sequence/state diagrams for the temporal race mechanisms, which the corpus lacks entirely, not on more schema pictures. Gate drift with a --check re-render.

## Disagreements

Anticipated: a floating :latest tag is acceptable because SVGs are committed; diagrams belong only in study docs; sequence diagrams cost more to maintain than they explain.

## Contribution

## Position

Pin the renderer, label every figure by provenance, gate drift mechanically — and then spend the
whole visual budget on **sequence/state diagrams for the race mechanisms**, not on more schema
pictures. The book today has no figures at all, and the corpus has no sequence diagram anywhere;
that gap is exactly where comprehension is weakest.

### 1. Pin the renderer before drawing anything new

All three renderers (`studies/01-charity-tree/diagrams/render.sh`,
`studies/02-ticket-booking/diagrams/render.sh`, `studies/03-reserved-seating/diagrams/render.sh`)
use `IMAGE="docker.io/plantuml/plantuml:latest"`, and `PLANTUML_IMAGE` is absent from
`infra/versions.env`. That violates the repository's own rule 3 (image versions are pinned; a
floating tag makes a result unreproducible) for the visual layer while the measurement layer is
pinned to tags (`postgres:17.11`, `yugabyte:2025.2.6.0-b111`, `redis:7.4.11-alpine`) and the book
toolchain is pinned by manifest digest. Consequences today:

- two `render.sh` invocations at different times can produce different SVGs from identical sources;
- a book figure compiled into a released PDF has no recorded renderer provenance, unlike the
  Typst-image digest the book already prints on its provenance page;
- the same script is copy-pasted three times, so a fix drifts.

Do: add `PLANTUML_IMAGE="docker.io/plantuml/plantuml:<exact tag or digest>"` to
`infra/versions.env`, move one render helper into `infra/` and have the three study scripts call it,
and have the helper print/write the resolved image id next to the rendered file list, mirroring the
`book/dist/build-manifest.json` idea (source list, image digest, output set). Until this exists, no
new diagram should be promoted into the book: an unpinned renderer makes "the picture matches the
SQL" unprovable.

### 2. Add a figure-provenance vocabulary and make the build refuse an unlabelled figure

The book's visual system (`book/lib/config.typ`) has a precise vocabulary — `direct`, `analogy`,
`gap` callouts and `evidence-card` — and the README states the governing rule: *no number without a
claim*. There is no equivalent rule for a picture, so a diagram could silently read as evidence.
Add to `lib/config.typ`:

```text
figure(path, caption, source, digest, kind)
  kind = illustration   drawing of committed design/SQL; explicitly NOT a measurement
       | mechanism      timeline reconstructed from a named analysis/discussion id
       | measured       rendered from run results (must also carry a claim id)
```

and extend `book/build.sh` to abort when a referenced figure is missing or its content digest
differs from the manifest input — the same failure mode that already aborts on an unknown claim id.
`book/README.md` should gain the one-line rule ("no figure without a source digest and a kind"), and
`docs/methodology.md`, which currently contains **no** rule for figures or diagrams at all, should
state that an illustration is never evidence and a mechanism figure must cite the analysis it was
reconstructed from.

### 3. Spend the budget where comprehension is actually missing: sequence and state diagrams

Evidence of the gap:

- `studies/0{1,2,3}-*/diagrams/` hold 30 committed SVGs (20 / 5 / 5) rendered from 30 `.puml`
  sources; `_style.puml` defines color roles for base/denormalised/rollup/embedded.
- All are structural, class-style schema pictures. A repo-wide search returns **zero**
  `sequenceDiagram` occurrences — no sequence, state, activity or component diagram exists.
- `studies/0{4,5,6,7,8,9}-*/` have **no `diagrams/` directory** at all.
- `book/` contains **no `#image`, no figure and no diagram reference**; the layout is entirely
  prose + callouts + evidence cards.

The corpus's hardest findings are all *temporal*, and today they are explained only in paragraphs:

- the strict cache fence must fire **before and after** commit, because a fill that captures the new
  fence and reads the pre-commit snapshot can still republish it
  (`studies/05-cache-consistency/reports/discussions/20260921-cache-fences-and-ledger.md` §1);
- the acknowledgement boundary is the ack, not the commit (same discussion);
- arbitration/convoy under a hot row (study 02 sell-out, study 03 hold-checkout);
- expiry and confirmation (study 03 E/K families), where "expiry is a write no client issues".

A class diagram cannot show any of these. Recommend a **closed set of 6–10 sequence/state figures**
authored for the book, each anchored to a run or analysis id (e.g. the fence timeline above), plus a
ticket/seat **state diagram** for the lifecycle the append-only and expiry concepts already discuss.
Cap the vocabulary at class (structure), sequence (interaction), state (lifecycle); reject
use-case/component/deployment UML as not serving this book — deployment belongs to the topology
study's prose.

### 4. Split the two audiences so the book does not duplicate the studies

Study `diagrams/` answer "what does D10 change versus D9" and belong beside the design tables that
already embed them (study READMEs 01–03 use `![…](diagrams/rendered/….svg)`). Book figures should
answer "why does this family behave this way" and be authored once at concept/mechanism level,
citing the study diagram or SQL instead of re-drawing it. A book figure that merely re-renders a
study schema diagram is duplication and should be a link. This keeps the figure set small, which is
also the maintenance budget.

### 5. Gate drift cheaply; generate from SQL only where it pays

The real maintenance risk is silent drift between `sql/<design>/schema.sql` and its picture. Two
options, in order:

1. **v1 (cheap, do now):** add a `--check` mode to the shared render helper that re-renders to a
   temporary directory and exits non-zero if the committed SVGs differ. Deterministic re-render plus
   the pinned image is enough to make drift visible in the study's existing checks.
2. **v2 (stronger, only for the book set):** a small Go generator that emits the class boxes and
   keys from the committed `schema.sql`. This makes "the picture provably matches the SQL" a build
   property rather than a convention, but it is new tooling and only worth it for the small,
   stable book figure set.

Sequence/state figures have no schema to derive from, so their accuracy gate is the cited
`analysis_id`/`run_id` in the figure's `source` field, checked by the book build.

## Reasoning and trade-offs

- The decision criteria are met in order: *faithful + labelled* (step 2), *reproducible in Podman*
  (step 1 + step 5), *comprehension* (step 3), *no duplication* (step 4).
- Cost: pinning and the `--check` mode are small; the figure-provenance function is a few lines of
  Typst; the sequence figures are real authoring work and must be cited, which is the point.
- The pinned renderer is a prerequisite, not a nice-to-have: without it, a figure's digest can only
  prove the SVG file did not change, not that a reader can reproduce the same SVG from the source.
- Sequence figures are deliberately few. Their value is that each one replaces a paragraph a reader
  currently has to reconstruct mentally; adding twenty would recreate the maintenance problem the
  schema diagrams already struggle with.

## Assumptions and evidence gaps

Observed: the three `render.sh` files and their `:latest` tag; the absence of `PLANTUML_IMAGE` in
`infra/versions.env`; the 30 committed SVGs and 30 sources; zero sequence/state diagrams; no
`diagrams/` in studies 04–09; no figures in `book/`; no figure rule in `docs/methodology.md`; the
book build's existing claim-id abort and manifest fields.

Inferred (not measured): that sequence diagrams *improve comprehension* — no reader study exists
and none is proposed; this is a documentation judgement, not a claim for the registry. That
PlantUML is pinnable by digest in the local Podman store is assumed. That Typst 0.15.1 embeds SVG
figures cleanly at the pinned image digest is assumed (only fonts and text are verified today).

## Alternatives

- **Mermaid instead of PlantUML.** Renders inline on Git hosts and is friendlier to review, but adds
  a second toolchain, discards the existing `_style.puml` visual language, and still needs a pinned
  container in this repository. Not selected: consistency and the existing style win.
- **Hand-authored SVG (diagrams.net etc.).** Best-looking, but not reproducible, no digest gate, and
  the source is not reviewable as text. Rejected outright under rule 3.
- **Prose only (status quo).** Cheapest, but leaves the temporal mechanisms — the corpus's strongest
  material — unshown, which is the question's whole complaint.
- **Generate everything from SQL.** Too much tooling for the 30 existing schema pictures; reserved
  for the book set.

## Disagreements

Anticipated objections, recorded without reading other independent positions:

- "PlantUML `:latest` is fine because the SVGs are committed." The committed SVG is frozen, but the
  source-to-output relation is not reproducible; the next render of an edited source is untraceable.
- "Diagrams belong inside the study SQL/README, not the book." Partly right for schema pictures;
  the book needs mechanism figures the studies do not have.
- "Sequence diagrams are expensive to keep accurate." True, which is why the recommendation caps the
  set and requires a citation, rather than banning them.

## Falsification tests

This recommendation changes if any of the following is observed:

1. `render.sh --check` shows the committed SVGs are non-deterministic even at one pinned PlantUML
   version (then the gate must be the SQL generator, and the digest may be meaningless);
2. PlantUML cannot be pulled/pinned by digest on this host, or the pinned image cannot render the
   existing 30 sources identically (then the pin must be documented as impossible and the
   provenance claim weakened honestly);
3. Typst at the pinned digest cannot embed SVG figures, or the PDF budget/page count changes
   materially (then figures must be raster artifacts with their own provenance);
4. the owner reviews the mechanism figures and judges an expanded textual walkthrough more useful
   than the figures (then the budget should shift back to prose + the labelling rule alone).
