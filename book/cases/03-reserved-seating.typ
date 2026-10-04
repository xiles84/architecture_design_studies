#import "../lib/config.typ": *
#import "../lib/evidence.typ": claim-badge

#pagebreak(weak: true)
== Reserved seats and checkout promises <case-seating>
#marker("case", "Study 03 · venue → event → seat")

Buyers choose physical seats and hold a block while paying. The system must honour a valid hold, refuse an expired confirmation and return abandoned seats to availability.

*Different from general admission:* another free seat is not a substitute for the one the buyer selected. A row's uniqueness protects sale identity; the hold protocol protects the time-limited promise.

#figure-evidence(
  "../assets/fig-hold-arbitration.svg", "sequence", "implemented design contract",
  "Two buyers attempt to hold the same seats.",
  [Separate checks let both buyers see free seats before either writes. A conditional update rechecks availability; lock-first waits and then checks. The illustration shows hold acquisition, not the confirmation refusal measurements below.],
  image-width: 100%,
)

=== Variation: check the promise inside the write

#block[
  #set text(size: 9.5pt)
  #table(columns: (1fr, 2fr, 1.3fr), inset: 6pt, stroke: 0.4pt + palette.rule,
    table.header([Variant], [What protects the hold], [Theory]),
    [s0 control], [Check, then change; the gap allows another buyer to intervene.], [#link(<theory-expiry>)[Expiry races]],
    [s1 conditional], [The write checks the current hold and state.], [#link(<theory-optimistic>)[Optimistic writes]],
    [s2 lock first], [Lock the relevant rows before checking and changing them.], [#link(<theory-pessimistic>)[Pessimistic writes]],
    [s1r retry once], [Retry the identical guarded confirmation after a transient refusal.], [#link(<theory-expiry>)[Promise boundary]],
  )
]

*Observed summary:* at the 10,000-seat tier, YugabyteDB s1 refused valid holds twice in one run and three times in reverse order; s2 refused once in each. The identical statement then matched every seat inside the refusing transaction.

*Observed correction:* s1r refused nothing; every retry succeeded. PostgreSQL's correct variants refused nothing across the measured tiers in both orders. This is a correctness result, not a throughput ranking.
#claim-badge("v2-10-guarded-confirm-refusals")

*Thrives — observed:* atomic guards or locks stop the check/change race. The engine-specific retry addressed the observed transient refusal.

*Dies — observed:* the unguarded negative control stole valid holds on both engines. A friendly pre-payment check cannot replace arbitration at confirmation.

=== Variation: read a section as one document

*Change:* L2 stores the section seat map as a document rather than assembling it from separate seat rows. L3 also changes the section's sharding package.

*Observed summary:* the replicated L2 comparison measures 24.4x and 27.8x on YugabyteDB single-node, versus 0.83x on PostgreSQL single-node. PostgreSQL shows no attributable advantage.
#claim-badge("v2-11-section-seat-map")

*Thrives — observed:* the YugabyteDB section-map package improves this broad read.

*Dies — observed boundary:* that gain does not transfer to PostgreSQL. The sharded L3 package also collapses point lookups in its primary single run; this has not been accepted as replicated evidence. A section read is not every read.

Theory: #link(<theory-embedding>)[Embedding] and #link(<theory-placement>)[Placement]. The family/package comparison changes more than one decision, so it cannot isolate embedding alone.

=== Variation: choose the hold window from history

*Business question:* how many holds convert to sales, expire or get released? That funnel would inform the checkout policy.

#gap[
  No measured design answers this hold-funnel report. Releasing an expired hold clears its row, making an expired hold indistinguishable from a seat never held.
  #claim-badge("v2-gap-02-hold-funnel-unanswerable")
]

*Thrives — inferred design need:* retained hold events would preserve that distinction. A history extension is a proposal here, not an executed design or measured result.

Theory: #link(<theory-history>)[History versus current state].

*Reproduce:* `studies/03-reserved-seating`. Badges identify the runs and limits. The configured hold policy and compressed benchmark time must not be confused with a measured optimal checkout duration.
