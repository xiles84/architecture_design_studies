// Data Architecture Reference — assembled book.
//
// Build: book/build.sh (Podman only). The build passes provenance through
// --input so the title and provenance pages state actual values, never
// hand-written ones. A build that does not pass them (an editor preview, or a
// hand-run typst compile) renders the UNVERIFIED banner instead of a provenance
// page full of placeholders.

#import "lib/config.typ": *
#import "lib/evidence.typ": evidence-index, active-claims, registry-card, registry-label, registry-version, gap-kind-label

// The kinds of gap the page explains are the ones the registry actually uses. A new
// kind therefore cannot appear without someone writing its definition: the build
// fails and names the missing one, instead of the page quietly listing three kinds
// while the cards use four.
#let gap-kind-description = (
  coverage: [a relevant regime has not been measured yet],
  schema_limitation: [the representation cannot answer the question, so no run would],
  instrument: [the property could not be observed with the tools the study had],
  unstable_measurement: [the number moved between runs with no code change],
)
#let gap-kind-order = ("schema_limitation", "coverage", "instrument", "unstable_measurement")
#let gap-kinds-in-use = {
  let kinds = ()
  for c in active-claims() {
    for k in c.at("gap_kind", default: ()) {
      if k not in kinds { kinds.push(k) }
    }
  }
  // A fixed reading order, with any kind the registry adds later at the end.
  kinds.sorted(key: k => {
    let i = gap-kind-order.position(x => x == k)
    if i == none { 99 } else { i }
  })
}
#for k in gap-kinds-in-use {
  if k not in gap-kind-description {
    panic("gap kind " + k + " has no definition in the legend that explains the evidence system")
  }
}
#let gap-kinds-legend() = {
  for k in gap-kinds-in-use [
    *#gap-kind-label(k)* — #gap-kind-description.at(k) · 
  ]
}

#let input(key, default: "unknown") = sys.inputs.at(key, default: default)

#let commit = input("commit")
#let describe = input("describe")
#let dirty = input("dirty")
#let built-at = input("built_at")
#let typst-version = input("typst_version")
#let typst-digest = input("typst_digest")
#let evidence-digest = input("evidence_digest")
#let pdf-source-hash = input("source_hash")
// build.sh sets release=1 only for a build that passed every input above.
#let release = input("release", default: "0") == "1"

#set document(
  title: book-title,
  author: "Architecture Design Studies",
)

#show: page-setup

// ---------------------------------------------------------------- title page
#page(header: none, footer: none)[
  #v(24mm)
  #text(size: 32pt, weight: "bold", fill: palette.accent)[#book-title]
  #v(4mm)
  #text(size: 14pt, fill: palette.muted)[A measured reference for data architecture and state design]
  #v(10mm)
  #line(length: 60%, stroke: 1pt + palette.rule)
  #v(6mm)
  #text(size: 11pt)[#book-edition]
  #v(28mm)
  #if release [
    #align(left)[
      #set par(justify: false)
      #set text(size: 9pt, fill: palette.muted)
      Built from commit #raw(describe) (#raw(commit)) — working tree #dirty. \
      Typst #typst-version (image #raw(typst-digest)). \
      Evidence registry #raw(registry-label), digest #raw(evidence-digest).
    ]
  ] else [
    #callout("Unverified development build", palette.gap)[
      Reproducibility metadata is incomplete: this PDF was compiled without the build script's
      provenance inputs. It is a preview. *Do not cite a measured value from this document.* \
      Rebuild with `book/build.sh` for a release-grade artefact.
    ]
  ]
]

// ------------------------------------------------------- provenance / edition
#page[
  // Not a level-1 heading: the chapter heading rule would force a break and
  // leave this page blank. This page is front matter, excluded from the outline.
  #text(size: 19pt, weight: "bold", fill: palette.accent)[Edition and provenance]
  #v(1mm)
  #line(length: 100%, stroke: 1.2pt + palette.accent)
  #v(3mm)
  This is a generated artefact: Typst source plus the evidence registry, compiled in a pinned
  Podman image. *Every claim in this book must resolve to a claim in the active registry*, and every
  number must resolve to a measured cell in a cited report. Nothing here is hand-entered at build
  time; the fields below are injected from the build script.

  #if not release [
    #v(3mm)
    #gap[
      *This is an unverified development build.* The fields below are placeholders because the
      document was not compiled by `book/build.sh`. Do not cite its numbers.
    ]
  ]

  #set text(size: 9.5pt)
  #table(
    columns: (auto, 1fr),
    stroke: 0.4pt + palette.rule,
    [Source commit], [#raw(commit)],
    [Source describe], [#raw(describe)],
    [Working tree at build], [#dirty],
    [Built at (UTC)], [#raw(built-at)],
    [Typst version], [#typst-version],
    [Typst image digest], [#raw(typst-digest)],
    [Evidence registry], [#raw(registry-label)],
    [Evidence digest], [#raw(evidence-digest)],
    [Source tree hash], [#raw(pdf-source-hash)],
  )

  #v(3mm)
  The book is assembled with #link("https://github.com/typst/typst")[Typst] inside a pinned Podman
  image; the evidence index at the end is rendered from the registry file at build time.

  #v(4mm)
  #gap[
    Environment and provenance limits travel with every number: all results come from one host
    (`host-zenbook-ux5406sa`) with shared cores and no real network between "cluster" nodes. Run
    tags mark the *producing code state*, not a tree containing the results. A reader reproducing
    a run must check out the run tag to get the code, then read the committed results directory.
  ]

  #heading("How to read a number here")
  #list(
    [Every *major-family* section follows the same progressive structure: quick choice, where it
      thrives and perishes, common scenarios, mechanism, costs, evidence and reproduction.
      Variant, scenario and concept chapters use the shape their role needs, because a case study
      and a reusable explanation are not the same thing.],
    [*Direct evidence* is measured in the cited run. *Mechanism evidence* demonstrates how
      something works without establishing a rate. *Analogy* is a labelled transfer to a family
      that was not measured. A #text(weight: "bold")[gap] callout is not one thing: each card
      names its own kind.],
    [The kinds of gap in use, read from the registry rather than restated here, because they close
      differently: #gap-kinds-legend()],
    [A claim can also be *partially superseded* — part of it has been replaced, the card names what
      moved to which successor, and the dimensions still open stay listed. A *retired* claim is no
      longer active and is not cited at all.],
    [Gaps are never null results: "not measured" must not be read as "would not change this".],
    [No cross-study numeric ranking appears anywhere; no comparability record exists.],
  )
]

// ------------------------------------------------------------------- contents
#outline(title: "Contents", indent: 1em, depth: 2)

// ------------------------------------------------------------------- chapters
#include "chapters/how-to-choose.typ"
#include "chapters/major-families.typ"
#include "chapters/minor-variants.typ"
#include "chapters/cross-family-comparisons.typ"
#include "chapters/topologies.typ"
#include "chapters/scenarios.typ"
#include "chapters/evidence-and-reproduction.typ"

// -------------------------------------------------------------------- concepts
#pagebreak()
#heading("Part II — Concepts")
The first part measured what each family does under this corpus's workloads. This part explains the
mechanisms those results came from, in the order the decisions bind: where information is placed,
what the indexes cost, who arbitrates a contended row, what derived state owes its maintainer, who
owns expiry, what a cache must prove, and where a cluster actually puts the data. The last chapter is
the checklist for reading any number in the first part.

#text(size: 9.5pt, fill: palette.muted)[Each concept names the claim cards that measured it, and the
glossary at the end of the book is the single definition site for every term used here.]

#include "concepts/information-placement.typ"
#include "concepts/indexes-and-hot-rows.typ"
#include "concepts/concurrency-control.typ"
#include "concepts/derived-state-and-history.typ"
#include "concepts/expiry-and-clock-authority.typ"
#include "concepts/cache-consistency.typ"
#include "concepts/distributed-placement.typ"
#include "concepts/reading-benchmarks.typ"

// ------------------------------------------------------------------- glossary
#pagebreak()
#include "glossary.typ"

// -------------------------------------------------------------- evidence index
#pagebreak()
// The heading names the registry from the single derived value, so a version bump
// cannot leave one page claiming a different package from the cover.
#heading("Evidence registry (active " + registry-version + " claims)")
The active registry holds #active-claims().len() claims that may be cited as evidence. This index is
rendered from #raw(registry-label) at build time; the digest above identifies its content.

#evidence-index()
