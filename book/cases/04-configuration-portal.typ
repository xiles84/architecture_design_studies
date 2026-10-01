// Chapter 12 — the Study 04 profile (STYLE.md, case-study template).
//
// One page: the business, the question, the designs compared, what it showed (each bullet
// with its badge), what is not measured, and how to reproduce it. Domain words are checked
// against studies/04-configuration-portal/README.md: an *installed product* owns its
// configuration, a *configuration entry* is one key and its value, the *effective
// configuration* is what a product reads back, and a *publication* makes a revision
// current.

#import "../lib/config.typ": *
#import "../lib/evidence.typ": claim-badge

#heading(level: 1, numbering: none, "Study 04 · configuration portal")
#marker("study", "04 · configuration portal")

#heading(level: 2, "The business")

A configuration portal is a service other products use: they create, publish and retrieve their
configuration through it. Configuration belongs to an installed product, so two installations of the same
product definition never share it.

#heading(level: 2, "The question")

What does each way of storing and serving that configuration cost to read, to write, to keep correct
under concurrency, and to hold?

#heading(level: 2, "Designs compared")

#set text(size: 9.5pt)
#table(
  columns: (auto, 1fr, auto),
  inset: 5pt,
  stroke: 0.4pt + palette.rule,
  align: (left + top, left + top, left + top),
  table.header([Design], [What it does], [Concept chapter]),
  [n1_rows_indexed], [One row per configuration entry, indexed: the reference design.], [2 Indexes],
  [n3_rollup_trigger], [Parent rollups maintained by a database trigger on every row written.], [4 Rollup],
  [n4_rollup_app], [The same rollups maintained by the application's own write path.], [4 Rollup],
  [n2_rows_rolldown], [Parent keys copied onto each child row to avoid a join.], [3 Rolldown],
  [n6_config_document], [The whole configuration held as one document per installation.], [5 Embedding],
)

#heading(level: 2, "What it showed")

#list(
  [Maintaining the overview with a trigger cost 4.0–5.7x the application-maintained rollup on a whole
    configuration replacement. #claim-badge("v2-13-trigger-rollup-cost")],
  [The penalty tracks the size of the write: roughly 1.1–3.2x at 10 entries, 2.5–9.0x at 60 and 10–24x
    at 240, with the ordering by size holding in every arrangement measured.
    #claim-badge("v2-18-maintenance-size-interaction")],
  [Repeated whole runs of the same pair disagreed by 35–115%, so the ordering is the finding and a single
    ratio is not. #claim-badge("v2-13-trigger-rollup-cost")],
  [No rollup mismatch was recorded at any size: the trigger's guarantee held while it cost more.],
)

#heading(level: 2, "Not measured")

#gap[
  Eight of the study's eighteen planned designs are not built. Cardinality, cadence and churn, repeated
  randomized trials, equal-total resources, placement and YugabyteDB are unanswered.
  #claim-badge("v2-gap-03-study04-missing-designs")
]

#heading(level: 2, "Reproduce")

Study path: `studies/04-configuration-portal`. The run tags behind the two claims above are named on
their evidence cards, which the badges link to; the designs, the runner and the environment are in the
study's README.
