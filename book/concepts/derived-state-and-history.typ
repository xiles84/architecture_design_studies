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
guarantee structural; the measured cost was 4.0–5.7x on whole-configuration replacement across three
runs (median 4.94x), with a 22–46x worse latency tail, while the application rollup kept writes near
the reference. That ratio is a reading at one size, not a property of the design: crossed with the size
of a replacement it runs about **1.1–3.2x** at 10 entries per installation, **2.5–9.0x** at 60 and
**10–24x** at 240, because the trigger does its work once per row written while the application does one
statement per replacement. Those are ranges and not single figures, and the ranges are wide: the runs at
each size disagreed by more than a factor of two, because the machine's speed wanders during a session —
sometimes up, sometimes down. What survives the wandering is the ordering: the same ranking appeared in
every arrangement that was measured, so treat the ranking as the finding and the size of the penalty as
uncertain by a factor of two or more.
Neither is universally right: choose who must never forget, and then ask how large a replacement is.
#registry-card("v2-13-trigger-rollup-cost")
#registry-card("v2-18-maintenance-size-interaction")

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
The ledger's write cost depends on the engine and on the order the designs were run at the small size, so
quote the run's own range and the disagreement rather than a single headline multiplier. The load-time
multiplier is not published at all.

#heading(level: 2, "The checklist")

The eight questions, answered for this concept; the status in brackets is the evidence behind the
answer.

- *Invariant.* Published derived state agrees with the entries it summarises, and retained history is
  never erased by the mutation that follows it. [direct]
- *Where the information lives.* In the base entries, in an aggregate on the parent, and — for
  answerability — in an append-only record that the mutation does not overwrite. [direct]
- *Where the contention is.* The parent row an aggregate is maintained on, which every child write
  touches. [direct]
- *What work moved from reads to writes.* Aggregate maintenance. History is different: it is
  write-only work bought for the questions it lets you answer later. [direct]
- *What state is derived.* The rollup columns. History is a fact, not a derived value, which is why it
  is not "stale" when it disagrees with current state. [direct]
- *What can become stale.* Any aggregate published after the commit that produced it — the drift
  control demonstrates it. [mechanism]
- *Where it stops working.* A per-row trigger's cost grows with the rows a replacement touches — now
  measured at three sizes, six runs each (about 1.3–3.2x, 2.5–9.0x and 10–24x at 10, 60 and 240
  entries, with the ranking holding in every sweep) rather than stated as a tendency; history costs
  storage (33–35% for the ledger). [direct]
- *The experiment that would settle it.* The n3/n4/x2 isolation pair, swept over entries per
  installation to find where the two maintenance strategies cross. [direct]
