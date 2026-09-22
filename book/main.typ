// Data Architecture Reference — assembled book.
//
// Build: book/build.sh (Podman only). The build passes provenance through
// --input so the title and provenance pages state actual values, never
// hand-written ones.

#import "lib/config.typ": *
#import "lib/evidence.typ": evidence-index, active-claims, registry-card

#let input(key, default: "unknown") = sys.inputs.at(key, default: default)

#let commit = input("commit")
#let describe = input("describe")
#let dirty = input("dirty")
#let built-at = input("built_at")
#let typst-version = input("typst_version")
#let typst-digest = input("typst_digest")
#let evidence-digest = input("evidence_digest")
#let pdf-source-hash = input("source_hash")

#set document(
  title: book-title,
  author: "Architecture Design Studies",
)

#show: page-setup

// ---------------------------------------------------------------- title page
#page(header: none, footer: none)[
  #v(30mm)
  #text(size: 32pt, weight: "bold", fill: palette.accent)[#book-title]
  #v(4mm)
  #text(size: 14pt, fill: palette.muted)[A measured reference for data architecture and state design]
  #v(10mm)
  #line(length: 60%, stroke: 1pt + palette.rule)
  #v(6mm)
  #text(size: 11pt)[#book-edition]
  #v(40mm)
  #align(left)[
    #set par(justify: false)
    #set text(size: 9pt, fill: palette.muted)
    Built from commit #raw(describe) (#raw(commit)) — working tree #dirty. \
    Typst #typst-version (image #raw(typst-digest)). \
    Evidence registry `book/evidence/v2/claims.json`, digest #raw(evidence-digest).
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
  Podman image. *Every claim in this book must resolve to a claim in the v2 registry*, and every
  number must resolve to a measured cell in a cited report. Nothing here is hand-entered at build
  time; the fields below are injected from the build script.

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
    [Evidence registry], [`book/evidence/v2/claims.json`],
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
    [Every family, variant, topology and scenario section follows the same progressive structure:
      quick choice, where it thrives and perishes, common scenarios, mechanism, evidence and
      reproduction.],
    [A #text(weight: "bold")[direct] callout is measured in the cited run. An
      #text(weight: "bold")[analogy] callout is a labelled transfer to a family that was *not*
      measured. A #text(weight: "bold")[gap] callout is missing evidence.],
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
#heading("Concepts")
#include "concepts/information-placement.typ"
#include "concepts/indexes-and-hot-rows.typ"
#include "concepts/concurrency-control.typ"
#include "concepts/derived-state-and-history.typ"
#include "concepts/expiry-and-clock-authority.typ"
#include "concepts/cache-consistency.typ"
#include "concepts/distributed-placement.typ"
#include "concepts/reading-benchmarks.typ"

// -------------------------------------------------------------- evidence index
#pagebreak()
#heading("Evidence registry (active v2 claims)")
The active registry holds #active-claims().len() claims. This index is rendered from the registry
file at build time; the digest above identifies its content.

#evidence-index()
