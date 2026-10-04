#import "../lib/config.typ": *
#import "../lib/evidence.typ": claim-badge

#pagebreak(weak: true)
== General-admission ticket sales <case-ticket>
#marker("case", "Study 02 · band → event → ticket")

A show company sells interchangeable admissions. Buyers need any available unit of capacity, while organisers edit events and report sales and refunds.

*Business invariant:* committed tickets must never exceed an event's capacity. A capacity unit here is not a physical seat chosen on a map.

#block(inset: 8pt, fill: palette.surface)[
  *Conceptual sequence:* available capacity → successful booking → cancellation or refund.
  \
  Current inventory answers “what is available?” A retained sale/refund record answers “what happened?”
]

=== Variation: spread the hot drop across seat rows

#block[
  #set text(size: 9.5pt)
  #table(columns: (1fr, 2fr, 1.4fr), inset: 6pt, stroke: 0.4pt + palette.rule,
    table.header([Variant], [Arbitration point], [Theory]),
    [Counter designs], [All buyers compete on one event or inventory counter.], [#link(<theory-hot-rows>)[Hot rows]],
    [P3 CAS], [A conditional update wins one pre-created capacity row.], [#link(<theory-optimistic>)[Optimistic writes]],
    [P2 skip-locked], [Lock one free row and skip rows other buyers already lock.], [#link(<theory-pessimistic>)[Pessimistic writes]],
  )
]

*Observed summary:* P3 and P2 sell 4-8x faster than single-counter designs in the hot-drop comparison. The range reproduced; absolute throughput did not.
#claim-badge("v2-06-hot-drop-designs")

*Thrives — observed:* many buyers want any capacity unit. Independent rows allow useful progress around contention.

*Dies — inferred boundary:* skipping another buyer's row cannot honour a request for a specific physical seat. Use #link(<case-seating>)[Reserved seating] for that different contract.

*Watch:* retry backoff matters. A counter on the event row can also block an organiser's unrelated edit; moving inventory away changes who shares the contention point.

=== Variation: retain the sale after a refund

*Change:* X1 adds an append-only ledger to P3. Cancellation returns capacity without erasing the financial event.

*Observed summary:* refunds are unanswerable in 14 of 15 designs; X1 alone answers the refund report. The ledger adds about 33-35% by table size, or about 35-36% counting every relation.
#claim-badge("v2-07-refund-answerability-and-storage")

*Thrives — observed answerability:* a business that needs refund reporting must retain the erased event somewhere.

*Dies — observed boundary:* current inventory alone cannot reconstruct refunds after deletion. Fast booking does not repair a missing fact.

Theory: #link(<theory-history>)[History versus current state]. Storage percentages depend on the counting rule; they are not one universal overhead.

=== Variation: pay for history under contention

#block[
  #set text(size: 9.5pt)
  #table(columns: (1fr, 2.4fr), inset: 6pt, stroke: 0.4pt + palette.rule,
    table.header([Regime], [Observed ledger cost over CAS]),
    [128 buyers; 10 seats], [About 1.1x; direction unreliable at this smallest tier.],
    [128 buyers; 100-1000 seats], [About 1.2-1.4x.],
    [128 buyers; 10000 seats], [About 1.8x.],
    [128 buyers; 100000 seats], [About 2.5x.],
  )
]
#claim-badge("v2-08-ledger-race-cost-128-buyers")

*Observed:* at 1,000 seats, PostgreSQL CAS is faster at every tested buyer count, 8 to 128, in both design orders. The earlier apparent change of winner was a confounded comparison.
#claim-badge("v2-09-postgres-ledger-ratio-unsettled")

*Thrives — inferred business choice:* keep the ledger when historical answers justify the booking cost.

*Dies — observed shortcut:* choosing only the fastest current-state design leaves required refund answers unavailable. Buyer count and event size must travel with any performance summary.

#gap[
  A ledger load-time multiplier is not published: unchanged code produced incompatible readings. Loading overhead is unresolved, rather than a fixed number.
  #claim-badge("v2-gap-01-ledger-load-multiplier-unstable")
]

*Reproduce:* `studies/02-ticket-booking`. Follow badges for producing tags and reports. Sell-rate evidence does not license quoting the old client-throttled read levels.
