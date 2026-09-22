#import "../lib/config.typ": *
#import "../lib/evidence.typ": registry-card

#heading("Reading benchmarks honestly")

#marker("concept", "reading benchmarks")
A benchmark answers the question its harness asks. Most wrong conclusions in this corpus came from
reading a number as if it answered a different question, not from a bad run.

#heading(level: 2, "Six questions to ask of any number")
#list(
  [*Closed or open loop?* Every race here is closed loop; its tail is a floor, not an SLO.],
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

#heading(level: 2, "Boundaries")
Laptop topologies with shared cores are a controlled comparison, not a capacity model. The
environment caveat belongs wherever a conclusion is drawn.
