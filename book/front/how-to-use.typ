// Front matter — "How to use this book" (STYLE.md, Edition 3 structure).
//
// One page: the anatomy of a concept chapter, the evidence badge legend, the short
// version of how to read a number, and the five studies at a glance. The legend and
// the badge meanings are generated from the active registry, so a strength or a gap
// kind the registry adds cannot appear on a page without a reader-facing meaning: the
// build stops and names it.
//
// The "how to read a number" material moved here from the provenance page, shortened:
// a reader who wants to know how to read the book should not have to first read the
// build's provenance to find it.

#import "../lib/config.typ": *
#import "../lib/evidence.typ": active-claims, strength-label, gap-kind-label, registry-version

// The kinds of gap the page explains are the ones the registry actually uses. A new
// kind therefore cannot appear without someone writing its definition: the build
// fails and names the missing one, instead of the page quietly listing three kinds
// while the cards use four.
#let gap-kind-description = (
  coverage: [the case simply has not been measured yet],
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

// The strengths the registry actually uses, in a fixed order, each with the plain
// wording STYLE.md requires. `strength-label` panics on a strength it does not know.
#let strength-order = (
  "run_replicated", "repeated_controlled", "single_run_directional",
  "mechanism_supported", "analogy", "gap",
)
#let strengths-in-use = {
  let s = ()
  for c in active-claims() {
    if c.strength not in s { s.push(c.strength) }
  }
  s.sorted(key: x => {
    let i = strength-order.position(y => y == x)
    if i == none { 99 } else { i }
  })
}

#heading(level: 1, numbering: none, "How to use this book")

This book teaches one concept per chapter and then shows the case studies that tested it. Read a
concept chapter to decide *what to do*; read the case boxes in it to see when that decision was
measured to pay off, and when it did not.

#heading(level: 2, "The anatomy of a concept chapter")

Every concept chapter answers the same ten things, in the same order. Once you know the shape, you can
read any chapter in a few minutes and stop where your decision is made.

#set text(size: 9.5pt)
#table(
  columns: (auto, 1fr),
  inset: 5pt,
  stroke: 0.4pt + palette.rule,
  align: (left + top, left + top),
  [Key idea], [The decision in one line, under the title.],
  [The problem], [Why the naive answer costs more than it looks.],
  [Picture], [One figure, with a text equivalent that says the same thing.],
  [How it works], [Three to five bullets: the mechanism, not the history.],
  [When it fits], [Two short lists: it fits when, and it does not fit when.],
  [What it costs], [Five fixed rows: reads, writes, storage, correctness, operations.],
  [In the studies], [The case boxes: what a study measured, and the verdict.],
  [Not measured yet], [The gaps, shown as gaps rather than answered by analogy.],
  [Remember], [At most three lines, for the reader who skims.],
  [The checklist], [The same eight questions every concept answers.],
)

#heading(level: 2, "How to read a badge")

A badge appears beside a sentence or a box and points at its evidence card in the index at the back of
the book. The badge names the claim, says how strong the evidence is in plain words, and gives the
card's page.

#set text(size: 9.5pt)
#table(
  columns: (auto, 1fr),
  inset: 5pt,
  stroke: 0.4pt + palette.rule,
  align: (left + top, left + top),
  ..strengths-in-use.map(s => (
    text(weight: "bold")[#strength-label(s)],
    if s == "run_replicated" [Measured more than once, in more than one whole run.]
    else if s == "repeated_controlled" [Repeated within one run, with a control.]
    else if s == "single_run_directional" [One run: the direction is what it supports.]
    else if s == "mechanism_supported" [The mechanism is implemented; no rate is claimed.]
    else if s == "analogy" [Transferred from another engine, family or workload.]
    else if s == "gap" [No evidence. Shown as a gap, never as a null result.]
    else [See the card.],
  )).flatten(),
)

A #text(weight: "bold")[gap] is not one thing. Each states its own kind, and the kinds close
differently: #gap-kinds-in-use.map(k => [*#gap-kind-label(k)* — #gap-kind-description.at(k)]).join(" · ").

#heading(level: 2, "Rules that hold everywhere")

#list(
  [A number appears only in a case box or a study profile, beside the badge of the claim that
    measured it.],
  [The number is copied from that claim's statement or limits. Nothing is recomputed for the page.],
  [A value a claim says did not reproduce is never quoted.],
  [Caveats live on the card. Each box carries the one limit that most changes how you read it.],
  [No cross-study numeric ranking appears anywhere, because no comparability record exists.],
  [A gap is never a null result: "not measured" must not be read as "would not change this".],
)

#heading(level: 2, "The five studies at a glance")

#set text(size: 9.5pt)
#table(
  columns: (auto, 1fr, 1fr),
  inset: 5pt,
  stroke: 0.4pt + palette.rule,
  align: (left + top, left + top, left + top),
  table.header(
    [Study], [The business], [Concepts it tests],
  ),
  [01 · charity donations], [Charities, donors and donations, with totals per charity and per donor.],
  [Where a fact lives · Indexes · Rolldown · Rollup · Embedding],
  [02 · ticket booking], [Events with seats, buyers racing for them, refunds and organiser edits.],
  [Hot rows · History vs current state],
  [03 · reserved seating], [Holds that must expire, and the clock that decides when.],
  [Expiry and the clock],
  [04 · configuration portal], [Products that read their configuration, and an overview that rolls it up.],
  [Rollup (trigger vs application) · Optimistic vs pessimistic],
  [05 · cache consistency], [Warm reads served from a cache that may lag its source.],
  [Caches and freshness · Where a cluster puts data],
)

The study paths, the designs each one compares and the runs behind every number are in chapter 12 and in
the evidence index. This book's registry is #raw("book/evidence/" + registry-version + "/claims.json").
