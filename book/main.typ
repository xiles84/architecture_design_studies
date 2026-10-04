// Data Architecture Reference — assembled book.
//
// Build: book/build.sh (Podman only). The build passes provenance through
// --input so the title and provenance pages state actual values, never
// hand-written ones. A build that does not pass them (an editor preview, or a
// hand-run typst compile) renders the UNVERIFIED banner instead of a provenance
// page full of placeholders.

#import "lib/config.typ": *
#import "lib/evidence.typ": evidence-index, active-claims, registry-card, registry-label, registry-version

// The gap-kind legend that used to live here moved to the front page
// (front/how-to-use.typ) with the rest of "how to read a number": the reader who
// needs it is deciding how to read the book, not checking this build's provenance.
// Its panic for an undefined kind moved with it.

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
  number must trace back to something actually measured in a report this book cites. Nothing here is
  hand-entered at build time; the fields below are injected from the build script.

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

  How to read a number — the badge legend, the gap kinds and the rules that hold everywhere —
is in *How to use this book*, on the page after this one.
]

#include "front/how-to-use.typ"

// ------------------------------------------------------------------- contents
#outline(title: "Contents", indent: 1em, depth: 2)

// ------------------------------------------------- Edition 3 parts (book/STYLE.md)
//
// The concept-first edition is built part by part: a part appears when it exists, a part
// opener names the chapters it will hold, and a chapter appears only once it is converted.
// An unconverted chapter is named and numbered in its part opener, so a reader can see what
// is missing instead of inferring it. The Edition 2 chapters follow the divider below in
// full, so nothing measured leaves the book while the conversion is in progress.
#part-opener(1, "Placing data",
  [Where a fact lives, what an index costs, and the two directions of copying a value: down to the
   children, or up to the parent.])

#include "concepts/rollup.typ"

#part-opener(2, "Many writers",
  [What happens when several writers want the same row, and what each coordination choice guarantees.
   Chapters 6–8 are not converted yet.])

#part-opener(3, "Time and copies",
  [History against current state, caches against freshness, and where a cluster puts the data.
   Chapters 9–11 are not converted yet.])

#part-opener(4, "Putting it together",
  [One profile per study, then how to choose and how to read a benchmark. Chapter 12 holds the one
   converted profile; chapters 13 and 14 are not converted yet.])

#include "cases/04-configuration-portal.typ"

// ------------------------------------------- Edition 2 chapters (awaiting conversion)
#pagebreak()
#heading(level: 1, numbering: none, "Edition 2 chapters")
#marker("edition", "2 · awaiting conversion")

The chapters below were written before the concept-first edition. They are *kept in the build* so that no
measured result leaves the book, and they await conversion into the concept chapters that own them
(book/STYLE.md, "Old file → new home"). Where a chapter below overlaps a converted one, the converted
chapter is the one to read.

#include "chapters/how-to-choose.typ"
#include "chapters/major-families.typ"
#include "chapters/minor-variants.typ"
#include "chapters/cross-family-comparisons.typ"
#include "chapters/topologies.typ"
#include "chapters/scenarios.typ"
#include "chapters/evidence-and-reproduction.typ"

#pagebreak()
#heading(level: 1, numbering: none, "Edition 2 concepts")
#marker("edition", "2 · awaiting conversion")

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
