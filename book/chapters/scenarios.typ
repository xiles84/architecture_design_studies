#import "../lib/config.typ": *
#import "../lib/evidence.typ": registry-card

#heading("Scenarios")

#marker("chapter", "scenarios")
A scenario is a set of questions plus a workload shape. The same family can be right in one and
wrong in another, so each scenario below states its answerability question and its evidence status.

#heading(level: 2, "Ticket overbooking and refunds")
Never oversell, then report refunds. Pre-created seats with compare-and-set or `SKIP LOCKED` sold
4–8x faster than a counter row; a ledger made refunds answerable at +33–35% storage. The ratio
between the two arbitration designs at 128 buyers is 1.1–2.3x on PostgreSQL; the small arm is
unsettled.
#registry-card("v2-06-hot-drop-designs")
#registry-card("v2-08-ledger-race-cost-128-buyers")
#registry-card("v2-09-postgres-ledger-ratio-unsettled")

#heading(level: 2, "Reserved seating and holds")
Hold, expire, confirm, and re-sell. Guarantees held everywhere, but YugabyteDB refused valid
confirmations until the refused statement was retried; and the funnel that would size the hold
window is unanswerable because releasing an expired hold erases the evidence.
#registry-card("v2-10-guarded-confirm-refusals")
#registry-card("v2-gap-02-hold-funnel-unanswerable")

#heading(level: 2, "Configuration portal")
Read the overview often, replace configurations occasionally. A trigger-maintained rollup cost
6.85x on replacement; the application rollup kept writes near the reference. This is a small-scale
result at one cardinality.
#registry-card("v2-13-trigger-rollup-cost")

#heading(level: 2, "Donor portal with an external cache")
Serve warm profile reads from a cache. Throughput rose 2.2–3.5x and p99 fell an order of magnitude;
three instances sharing one relaxed cache recorded about 87% wrong reads unless the publication was
fenced.
#registry-card("v2-14-cache-throughput-gain")
#registry-card("v2-15-three-instance-staleness")

#heading(level: 2, "Scenario transfers are labelled")
#analogy[
  A supermarket stock counter, an auction bid ledger or an e-commerce inventory pool is structurally
  similar to the ticketing hot row, but none of them was measured here. Use those transfers as
  analogies from the hot-row mechanism, never as separate evidence.
]
