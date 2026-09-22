#import "../lib/config.typ": *
#import "../lib/evidence.typ": registry-card

#heading("Expiry and clock authority")

#marker("concept", "expiry and clock authority")
Expiry is a write that no client issues. That makes it the place where "whose clock is right" stops
being philosophical: the sweeper, the application clock and the database all disagree eventually,
and the design must say who is allowed to decide.

#heading(level: 2, "What the corpus shows")
YugabyteDB refused valid hold confirmations under a guarded design, and retrying the refused
statement once removed every refusal at no measurable cost. The practical lesson is smaller than it
sounds: an expiry design needs a retry story, because the mechanism can refuse a correct write.

#heading(level: 2, "What it cannot yet show")
No design can answer the funnel that would size a hold window. Releasing an expired hold clears the
row, so an expired hold is indistinguishable from a seat that was never held. The missing evidence
is a schema property, not a missing measurement.

#registry-card("v2-10-guarded-confirm-refusals")
#registry-card("v2-gap-02-hold-funnel-unanswerable")

#heading(level: 2, "Boundaries")
Do not transfer a single-engine expiry result to a clock-synchronised cluster without a placement
and clock test; the corpus has neither.
