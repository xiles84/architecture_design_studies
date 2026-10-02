#import "../lib/config.typ": *
#import "../lib/evidence.typ": registry-card

#heading("Expiry and clock authority")

#marker("concept", "expiry and clock authority")
Expiry is a write that no client issues. That makes it the place where "whose clock is right" stops
being philosophical: the sweeper, the application clock and the database all disagree eventually,
and the design must say who is allowed to decide.

#heading(level: 2, "What the corpus shows")
YugabyteDB refused valid hold confirmations under a guarded design — every refusal of that kind was
*transient*, meaning the identical statement re-issued inside the refusing transaction matched every
seat — and retrying the refused
statement once removed every refusal at no measurable cost. The practical lesson is smaller than it
sounds: an expiry design needs a retry story, because the mechanism can refuse a correct write.

#figure-evidence(
  "../assets/fig-expiry-clock.svg",
  "sequence",
  "implemented design contract",
  "When an abandoned hold gives its seat back, and whose clock says so.",
  [The figure lays out the expiry designs side by side — lazy expiry judged by the database clock, a
   sweeper that decides when a seat returns, and an application clock shown as a negative control that
   can release a seat early. It is the contract the reserved-seating designs implement and carries no
   measured value; the funnel that would size a hold window remains unanswerable (`v2-gap-02`). Source:
   `studies/03-reserved-seating/diagrams/e_expiry.puml`],
)

#heading(level: 2, "What it cannot yet show")
No design can answer the funnel that would size a hold window. Releasing an expired hold clears the
row, so an expired hold is indistinguishable from a seat that was never held. The missing evidence
is a schema property, not a missing measurement.

#registry-card("v2-10-guarded-confirm-refusals")
#registry-card("v2-gap-02-hold-funnel-unanswerable")

#heading(level: 2, "Boundaries")
Do not transfer a single-engine expiry result to a clock-synchronised cluster without a placement
and clock test; the corpus has neither.

#heading(level: 2, "The checklist")

The eight questions, answered for this concept; the status in brackets is the evidence behind the
answer.

- *Invariant.* A hold made while it is valid is honoured, nothing is granted or sold after expiry, and
  no seat leaks. [direct]
- *Where the information lives.* In `expires_at` on the held row — or, in one variant, once on a cart
  row rather than copied to every seat. [direct]
- *Where the contention is.* The seat rows of one chosen block, which cannot be skipped the way study
  02's interchangeable seats can. [mechanism]
- *What work moved from reads to writes.* Expiry checks move to the moment of use, and a sweeper does
  the work at rest instead of on the buyer's path. [mechanism]
- *What state is derived.* Availability is derived, not stored: with lazy expiry an expiry changes no
  row, so no exact availability counter can be maintained. [mechanism]
- *What can become stale.* Any judgement of expiry made from an application clock, and any cached
  availability figure. [direct]
- *Where it stops working.* When the clock is not the database's — the injected-skew control fired —
  and when the sweeper stalls, which leaks seats. [gap]
- *The experiment that would settle it.* An injected application-clock skew and a sweeper outage,
  beside the honoured-hold invariant. [direct]
