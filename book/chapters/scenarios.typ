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

#figure-evidence(
  "../assets/fig-hold-arbitration.svg",
  "sequence",
  "implemented design contract",
  "Two buyers choose the same marked seats at the same moment.",
  [Because a reserved-seating buyer may not be moved to another seat, Study 02's `SKIP LOCKED` scan
   does not apply here: the arbitration is a conditional update on the chosen seats (or a lock first,
   with the loser retrying) and one unguarded design is shown as a negative control. The figure is the
   contract the reserved-seating designs implement; it carries no measured value. Source:
   `studies/03-reserved-seating/diagrams/s_arbitration.puml`],
)

#heading(level: 2, "Configuration portal")
Read the overview often, replace configurations occasionally. A trigger-maintained rollup cost
6.85x on replacement; the application rollup kept writes near the reference. This is a small-scale
result at one cardinality.
#registry-card("v2-13-trigger-rollup-cost")

#heading(level: 2, "Donor portal with an external cache")
Serve warm profile reads from a cache. Throughput rose 2.2–3.5x and p99 fell an order of magnitude.
Three instances sharing one relaxed cache recorded 87.81% wrong reads on the owned model and 87.24% on
the legacy model: the publication must be fenced against the stale-fill race — a reader republishing
S0 after a writer commits S1 and invalidates — not merely ordered after the commit. The single-instance
relaxed cells recorded zero wrong reads, and the legacy process-local arms' zero was a bypass result
(`v2-15`), so neither is evidence that an unfenced relaxed cache is safe.
#registry-card("v2-14-cache-throughput-gain")
#registry-card("v2-15-three-instance-staleness")
#registry-card("v3-01-churn-crosses-real-ttl")
#registry-card("v3-03-equal-total-framing-labelled")

#heading(level: 2, "Scenario transfers are labelled")
#analogy[
  A supermarket stock counter, an auction bid ledger or an e-commerce inventory pool is structurally
  similar to the ticketing hot row, but none of them was measured here. Use those transfers as
  analogies from the hot-row mechanism, never as separate evidence.
]

#heading(level: 2, "Where a scenario stops")
#gap[
  Two scenario families are explicitly incomplete, and the book does not fill them with words: the
  configuration-portal study did not build eight of its eighteen designs and did not measure
  cardinality or cadence, and the cache study still owes its hard TTL expiry, medium scale,
  tablet/leader placement and open-loop demand. The cache study *did* measure real TTL churn,
  in-run repeats and equal-total cache accounting on 2026-09-23, so those are no longer gaps. The
  vocabulary is worth fixing here: an *in-run repeat* is another measurement window inside one
  benchmark process, which reduces within-run noise but is not an independent run replication —
  where the cards say `Trials: 1`, they mean one whole-run replication and no more.
  These are coverage gaps, not null results.
]
#registry-card("v2-gap-03-study04-missing-designs")
#registry-card("v3-gap-01-study05-remaining-regimes")
