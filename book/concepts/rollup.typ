// Chapter 4 — Rollup (STYLE.md, Edition 3 pilot).
//
// One concept, the template's ten sections in order, and every number inside a case box
// beside the badge of the claim that measured it. The three claims cited here are the
// rollup family (v2-02), the trigger-versus-application maintenance cost and its
// size dependence (v2-13, v2-18), and the Study 04 coverage gap (v2-gap-03).
//
// The sibling chapter — rolldown, copying a parent value down to children — is chapter 3
// and is not repeated here: a rollup aggregates a child change up, and the two directions
// are opposite, which is what the figure shows.

#import "../lib/config.typ": *
#import "../lib/evidence.typ": claim-badge, registry-card

// Edition 3 chapters carry their number in the opener, not in a heading counter, so
// their sections are unnumbered: without this they print as 0.1, 0.2 ... because the
// level-1 heading that would own the counter is unnumbered.
#set heading(numbering: none)

#concept-opener(4, "Rollup",
  [Store a total on the parent so reads skip the children — and pay for it on every child write.])

A charity page shows totals per charity and per donor. Calculating them on every view scans the same
donations repeatedly.

A rollup stores each total on its parent row. Reads fetch that row; every contributing donation
change must maintain it.

#figure-evidence(
  "../assets/fig-rollup-vs-rolldown.svg",
  "structure",
  "conceptual illustration",
  "Rolldown fans a parent change out; rollup aggregates a child change up.",
  [*Rolldown* copies parent information to children: a parent change causes fan-out writes.
   *Rollup* aggregates child facts on the parent: a child change updates its total.
   The arrows show maintenance direction, not a measured rate.],
)

#heading(level: 2, "How it works")

#list(
  [The parent row holds one stored value per total — a column, not a computed view.],
  [Every write that changes a contributing child updates that parent value in the same transaction.],
  [The maintainer is either a database trigger or the application's own write path.],
  [Readers then answer the aggregate from the parent row and never scan the children.],
  [Nothing else is derived: the child rows stay the facts of record.],
)

#heading(level: 2, "When it fits")

#fit-table(
  use: (
    [Totals are read far more often than the contributing rows are written.],
    [One parent row is an acceptable contention point for its writers.],
    [The write that changes the total can maintain it atomically.],
  ),
  avoid: (
    [Children are written far more often than their parents are read.],
    [A parent row is already hot for the application's own writes.],
    [The total must be recomputed from history at read time anyway.],
    [The aggregate's definition changes often, so stored values would be rebuilt often.],
  ),
)

#heading(level: 2, "What it costs")

#cost-table(
  reads: [One row per total, no child scan.],
  writes: [One parent write per contributing child change.],
  storage: [A column per maintained total.],
  correctness: [Silently wrong if one maintenance path is missed.],
  operations: [A backfill whenever the rule changes.],
)

#heading(level: 2, "In the studies")

#case-box(
  study: "Study 01 · charity donations",
  verdict: "fits",
  situation: [A charity site shows totals per charity and per donor.],
  compared: [Totals stored on the parent by a trigger or by the application, against totals computed
    from the donations on every read.],
  result: [Stored totals scored about 61–85x the minimal normalized design on aggregate reads. Inserts
    ran at about a third to two-fifths of its rate.],
  lesson: [A rollup fits when totals are read much more often than donations are written.],
  watch-out: [The read gain and the insert penalty reproduced; the absolute insert rates did not.],
  evidence: [#claim-badge("v2-02-rollup-family")],
)

#case-box(
  study: "Study 04 · configuration portal",
  verdict: "depends",
  situation: [Products read their configuration from a portal, and an overview rolls each one up.],
  compared: [The overview kept by a database trigger against kept by the application.],
  result: [Replacing a whole configuration was 4.0–5.7x slower with the trigger. The gap grew with
    size: about 1.1–3.2x at 10 entries, 2.5–9.0x at 60 and 10–24x at 240.],
  lesson: [A trigger buys a guarantee no code path can forget, paid on every row written. Check how
    large your writes are.],
  watch-out: [Small scale, one machine, and repeats disagreed widely: trust the ordering, not the
    ratio.],
  evidence: [
    #claim-badge("v2-13-trigger-rollup-cost")
    #claim-badge("v2-18-maintenance-size-interaction")
  ],
)

#gap[
  #text(weight: "bold")[Not measured yet.] Eight of Study 04's eighteen planned designs are not built,
  so cardinality, churn, placement, equal-total resources and YugabyteDB are unanswered for this pair.
  #claim-badge("v2-gap-03-study04-missing-designs")
]

#remember(
  [Read the ratio, not the rate: the rate moves between runs.],

  [Maintenance is the design: choose the maintainer, then the total.],

  [The penalty is a function of write size, so measure your write.],
)

// Kept in one unbreakable group: a checklist that splits after seven rows leaves the eighth
// alone on a page, which reads as an afterthought rather than the chapter's frame.
#block(breakable: false)[
#heading(level: 2, "The checklist")

The eight questions every concept chapter answers, with the evidence behind each answer in brackets.

#checklist-table((
  ("What is the invariant?",
   [A stored total equals the sum of the entries it summarises.], "direct"),
  ("Where does the information live?",
   [In the child rows, and in one stored value on the parent.], "direct"),
  ("Where is the contention?",
   [The parent row that every child write touches.], "direct"),
  ("What work moved from reads to writes?",
   [Recomputing the total on each read.], "direct"),
  ("What state is derived?",
   [The stored total; the child rows stay the facts.], "direct"),
  ("What can become stale?",
   [A total whose maintenance path was skipped.], "mechanism"),
  ("Under which workload does it stop working?",
   [Write-heavy children, or a size the ratio has not been measured at.], "gap"),
  ("Which experiment would prove it?",
   [One controlled pair maintained by trigger against application.], "direct"),
))
]
