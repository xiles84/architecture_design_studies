#import "../lib/config.typ": *
#import "../lib/evidence.typ": strength-label, active-claims, gap-kind-label
#set heading(numbering: none)
#heading(level: 1, numbering: none, "How to use this book")

#callout("Two reading paths", palette.accent)[
  *Learn a strategy:* Theory → topic → strategy. \
  *Choose for a situation:* Study cases → business case → tested variation.
]

#table(columns: (1fr, 1.7fr), inset: 6pt, stroke: 0.4pt + palette.rule,
  table.header([Part], [What you will find]),
  [#link(<part-theory>)[Theory]], [The mechanism, a picture, fit/failure conditions, costs and a case link.],
  [#link(<part-cases>)[Study cases]], [The business, variations, thrive/die situations and compact measured summaries.],
  [Appendix], [How to judge a result, reproduce it and inspect its full evidence card.],
)

#heading(level: 3, numbering: none, "Reading a thrive-or-die example")
#list(
  [*Observed:* a named comparison in the cited evidence showed this effect.],
  [*Inferred:* the mechanism suggests this boundary, but that workload was not measured.],
  [*Not measured:* there is no result to resolve the choice.],
)

#heading(level: 3, numbering: none, "Reading a number")
A badge links a summary to its full claim card. Read the workload and the limit with the number.
Results come from one shared-core host with no real network between cluster containers.

#table(columns: (1fr, 1.7fr), inset: 6pt, stroke: 0.4pt + palette.rule,
  table.header([Badge], [Meaning]),
  [Replicated], [Observed across whole runs; the card still names spread and limits.],
  [Single run], [Direction in one run, not a stable universal magnitude.],
  [Mechanism], [An implemented mechanism; no measured rate is implied.],
  [Analogy], [Transferred intuition, not evidence for the new workload.],
  [Gap], [Absent coverage, a missing instrument, an unsupported schema or unstable measurement.],
)

#remember(
  [Learn the mechanism in Theory; check its tested boundary in the matching case.],
  [A read gain may carry a write, storage or correctness cost.],
  [The full evidence cards are reference material, not extra teaching chapters.],
)
