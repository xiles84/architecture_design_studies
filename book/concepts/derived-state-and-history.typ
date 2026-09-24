#import "../lib/config.typ": *
#import "../lib/evidence.typ": registry-card

#heading("Derived state and history")

#marker("concept", "derived state and history")
Two different things are stored beside the canonical row, and conflating them produces the wrong design.
*Derived state* — a rollup, a rolldown, an external read copy — optimizes an answer from current state
and owes its maintainer a correctness contract on every contributing write. *Retained history* — an
event or ledger table, an audit trail, a temporal history, snapshots — preserves information that would
otherwise disappear. Derived state can be rebuilt from the facts; retained history cannot be rebuilt
from a current-state-only representation.

#figure-evidence(
  "../assets/fig-derived-vs-history.svg",
  "structure",
  "conceptual illustration",
  "One canonical fact feeds a derived current-state copy or a retained history.",
  [A canonical fact has two stored branches. A *derived current-state copy* (rollup or rolldown) answers
   from current state, can be rebuilt from the facts and owes its maintainer a correctness contract on
   every contributing write. *Retained historical facts* (an event/ledger) preserve information that
   would otherwise disappear and cannot be rebuilt once a current-state-only representation has
   destroyed it. The figure asserts no rate.],
)

#heading(level: 2, "Maintenance is the design")
A rollup can be maintained by the application or by a database trigger. The trigger makes the
guarantee structural; the measured cost was 6.85x on whole-configuration replacement (83.7 vs 573
ops/s, p99 2.5 to 42 ms), while the application rollup kept writes near the reference for the same
read gain. Neither is universally right: choose who must never forget.

#heading(level: 2, "History is answerability")
Retained history is not faster; it is the difference between being able to answer a question after the
event was undone and not. A current-state-only representation that overwrites or deletes the earlier
fact cannot reconstruct it, whichever physical structure stores the answer — an event table, an audit
trail, a temporal history or a snapshot can each supply it. The refunds window was unanswerable in 14 of
15 ticketing designs and cost about 33–35% more storage with a ledger. That is a correctness purchase,
priced honestly. Retaining history does not itself fence concurrent writers or order events; that is a
separate, explicit decision.

#figure-evidence(
  "../assets/fig-history-vs-current.svg",
  "structure",
  "conceptual illustration",
  "Undoing the state is not the same as erasing the evidence.",
  [Left, a current-state-only representation overwrites a Donation's PAID status on refund, so
   "was this ever paid?" becomes unanswerable. Right, appended `DonationEvent` rows (PLEDGED, PAID,
   REFUNDED) answer the same question: undoing the state keeps the evidence. The figure is the
   difference between derived state and retained history, and it asserts no rate.],
)

#heading(level: 2, "Direct evidence")
#registry-card("v2-13-trigger-rollup-cost")
#registry-card("v2-09-postgres-ledger-ratio-unsettled")

#heading(level: 2, "Boundaries")
The ledger's write cost is engine- and order-sensitive at the small arm; quote the run's own range
and the disagreement, not a headline multiplier. The load-time multiplier is not published at all.
