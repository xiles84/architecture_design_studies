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
#outline(title: "Contents", indent: 1em, depth: 3)

// -------------------------------------------------------- reader structure
#set heading(numbering: none)
#heading(level: 1, numbering: none, "Part I · Theory") <part-theory>
Theory explains the choice. Each topic groups strategies; each strategy shows its mechanism,
fit, cost and an experiment that could challenge the expectation.

#include "theory/schemas.typ"
#include "theory/selects.typ"
#include "theory/writers.typ"
#include "theory/caches.typ"
#include "theory/queues.typ"
#include "theory/placement.typ"

#heading(level: 1, numbering: none, "Part II · Study cases") <part-cases>
Each case starts with its business problem, then applies the strategies from Part I.
Summary numbers belong to the named comparison and environment, not to a universal ranking.

#heading(level: 2, numbering: none, "Study-case coverage map") <case-coverage>
#table(columns: (1fr, 1.7fr), inset: 6pt, stroke: 0.4pt + palette.rule,
  table.header([Case], [Theory put into practice]),
  [#link(<case-charity>)[Charity donations]], [Schemas, indexes, copied keys, rollups and embedding; measured SELECT templates from #link(<case-charity-join-plan>)[join and filter] to #link(<case-charity-pagination>)[pagination], #link(<case-charity-placement>)[verified cluster placement] and the #link(<case-charity-select-guide>)[template guide].],
  [#link(<case-ticket>)[Ticket booking]], [Hot rows, history and current state.],
  [#link(<case-seating>)[Reserved seating]], [Expiry, ownership and bounded embedding.],
  [#link(<case-config>)[Configuration portal]], [Rollup maintenance and concurrent editing.],
  [#link(<case-cache>)[Cache consistency]], [Freshness, invalidation, fencing and cluster access.],
)
#gap[The SELECT families are measured on one laptop with one sequential reader. Their YugabyteDB
membership, absence and pagination figures come from scan-only plans, and only the join family has a
cluster arm. Queues have no registered measured case. Their expected boundaries are teaching examples,
not benchmark results.]

#include "cases/01-charity-tree.typ"
#include "cases/02-ticket-booking.typ"
#include "cases/03-reserved-seating.typ"
#include "cases/04-configuration-portal.typ"
#include "cases/05-cache-consistency.typ"

#heading(level: 1, numbering: none, "Appendix · Evidence and reproduction")
#include "appendix/reading-results.typ"
#include "appendix/query-plans.typ"

// ------------------------------------------------------------------- glossary
#pagebreak()
#include "glossary.typ"

// -------------------------------------------------------------- evidence index
#pagebreak()
// The heading names the registry from the single derived value, so a version bump
// cannot leave one page claiming a different package from the cover.
#heading(level: 2, numbering: none, "Evidence registry (active " + registry-version + " claims)")
The active registry holds #active-claims().len() claims that may be cited as evidence. This index is
rendered from #raw(registry-label) at build time; the digest above identifies its content.

#evidence-index()
