#import "../lib/config.typ": *
#import "../lib/evidence.typ": registry-card

#heading("Reading benchmarks honestly")

#marker("concept", "reading benchmarks")
A benchmark answers the question its harness asks. Most wrong conclusions in this corpus came from
reading a number as if it answered a different question, not from a bad run.

#heading(level: 2, "Six questions to ask of any number")
#list(
  [*Closed or open loop?* Every race here is closed loop. Closed-loop load can understate
   fixed-rate behaviour, so its tail must not be read as an SLO prediction.],
  [*One trial or many?* Spreads within a measurement are not an across-trial error bar.],
  [*Per-node or equal-total framing?* Single-vs-cluster answers "add machines", not "same budget".],
  [*What was the design order?* In Study 02 the two arms disagree in sign on PostgreSQL, so no ratio
   is quotable until a repeated-design control is run.],
  [*What fired the control?* A correctness gate without a fired negative control proves less than it
   appears to.],
  [*What is the cell, and what is the run?* A run's totals and failures belong to the run; a claim's
   support cells are a selected subset.],
)

#heading(level: 2, "Worked examples of a number that cannot be quoted")
#registry-card("v2-09-postgres-ledger-ratio-unsettled")
#registry-card("v2-gap-01-ledger-load-multiplier-unstable")
#registry-card("v2-16-strict-freshness-read-cost")

#heading(level: 2, "Before quoting a number")

A table, not a figure: check it against the card beside the number.

#table(
  columns: (1.5fr, 1fr, 1fr),
  stroke: 0.4pt + palette.rule,
  inset: 4pt,
  align: (left, left, left),
  table.header([*Question*], [*Where the answer lives*], [*If it is missing*]),
  [Which claim does this number belong to?], [The claim id in the card], [Do not quote it.],
  [Run, replication or trial?], [`trials` on the card], [Say "one run, one trial" and stop.],
  [Closed or open loop?], [The report's workload section], [Do not read the tail as an SLO.],
  [Per-node or equal-total framing?], [`resource_framing` in the result], [Never pool the two arms.],
  [One decision or a package?], [The card's own label], [Attribute nothing to a single change.],
  [Did the negative control fire?], [The run's verification block], [Treat correctness as unproven.],
  [What confound moved with it?], [`confounds` on the card], [Do not attribute the effect.],
  [What gap does it leave?], [The gap's `gap_kind`], [Name the kind, never "unknown".],
)

#heading(level: 2, "Boundaries")
Laptop topologies with shared cores are a controlled comparison, not a capacity model. The
environment caveat belongs wherever a conclusion is drawn.
