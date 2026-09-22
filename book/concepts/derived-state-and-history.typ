#import "../lib/config.typ": *
#import "../lib/evidence.typ": registry-card

#heading("Derived state and history")

#marker("concept", "derived state and history")
A rollup, a copied key and an append-only ledger are all "a fact stored twice". They differ in what
they can answer after the fact and what they cost on every write.

#heading(level: 2, "Maintenance is the design")
A rollup can be maintained by the application or by a database trigger. The trigger makes the
guarantee structural; the measured cost was 6.85x on whole-configuration replacement (83.7 vs 573
ops/s, p99 2.5 to 42 ms), while the application rollup kept writes near the reference for the same
read gain. Neither is universally right: choose who must never forget.

#heading(level: 2, "History is answerability")
Ledgers and snapshots are not faster; they are the only structures that can answer a question after
the event was undone. The refunds window was unanswerable in 14 of 15 ticketing designs and cost
about 33–35% more storage with a ledger. That is a correctness purchase, priced honestly.

#heading(level: 2, "Direct evidence")
#registry-card("v2-13-trigger-rollup-cost")
#registry-card("v2-09-postgres-ledger-ratio-unsettled")

#heading(level: 2, "Boundaries")
The ledger's write cost is engine- and order-sensitive at the small arm; quote the run's own range
and the disagreement, not a headline multiplier. The load-time multiplier is not published at all.
