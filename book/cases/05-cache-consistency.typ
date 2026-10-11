#import "../lib/config.typ": *
#import "../lib/evidence.typ": claim-badge

#pagebreak(weak: true)
== A cached donor portal <case-cache>
#marker("case", "Study 05 · source database → donor page copy")

A donor portal repeatedly serves a donor's identity, totals and latest donations. A cache avoids rebuilding that page; writes make the cached copy potentially obsolete.

*Business contract:* a cached value must correspond to committed database state. Relaxed freshness can return an older committed value; strict freshness forbids values older than writes acknowledged before the read began.

#figure-evidence(
  "../assets/fig-cache-read-copy.svg", "structure", "conceptual illustration",
  "A cache serves a copy of the donor page.",
  [The database remains authoritative. A cache hit serves a stored read copy; a miss reaches the database and may publish a replacement. Copies require a separate freshness protocol.],
  image-width: 76%,
)

=== Variation: cache a warm read

#block(breakable: false)[
  #set text(size: 9.5pt)
  #table(columns: (1fr, 2fr, 1.3fr), inset: 6pt, stroke: 0.4pt + palette.rule,
    table.header([Variant], [What changes], [Theory]),
    [No cache], [Build the donor page from the source database.], [#link(<theory-normalized>)[Source facts]],
    [Cache aside], [On a miss, read the source and fill a cache entry.], [#link(<theory-cache-aside>)[Cache aside]],
    [Write-through], [Maintain a cache copy through the coordinated write path.], [#link(<theory-cache-invalidation>)[Copies and writes]],
    [Strict protocol], [Coordinate or validate freshness against acknowledged writes.], [#link(<theory-cache-fencing>)[Freshness fencing]],
  )
]

*Observed summary:* warm cached reads buy 2.2-2.5x throughput in the repeated PostgreSQL single-node + Redis comparison. The primary tail-latency claim did not reproduce.
#claim-badge("v2-14-cache-throughput-gain")

*Thrives — observed:* warm, cacheable page reads avoid repeated source work.

*Dies — inferred boundary:* a miss-heavy or constantly changing workload may repay less of the copy's cost. No universal hit-rate crossover is established here.

=== Variation: coordinate copies across application instances

*Observed failure:* three logical instances sharing Redis under relaxed write-through record about 87% wrong reads. A lease private to one process cannot coordinate the other instances.
#claim-badge("v2-15-three-instance-staleness")

*Thrives — inferred requirement:* coordination must span every participant that can publish or mutate the shared value.

*Dies — observed:* adding a database version token without checking it does not remove stale reads. A single-instance result does not establish a multi-instance guarantee.

Theory: #link(<theory-cache-fencing>)[Lease and version fencing] and #link(<theory-cache-invalidation>)[Invalidation]. For an asynchronous extension, see #link(<theory-queue-idempotency>)[Idempotent queue consumers]; its full delivery behaviour is not measured here.

=== Variation: require strict freshness

*Observed summary:* on PostgreSQL, strict and relaxed read throughput differ by 7.3% and 5.7%, with the sign flipping inside trial spread. No read-throughput difference is resolved.
#claim-badge("v2-16-strict-freshness-read-cost")

*Observed on YugabyteDB:* the single-node strict/relaxed ratios, 1.043x and 1.105x, also lie inside spread. No winner is established.
#claim-badge("v3-04-yb-single-strict-relaxed")

*Single-run cluster reading:* three-node strict write-through records 942 ops/s versus relaxed 846. Each configuration is measured once; the opposite single-node direction does not justify an engine ranking.
#claim-badge("v3-05-yb-cluster3-strict-relaxed")

*Thrives — observed contract:* strictness can be maintained in the measured cells without a resolved read penalty.

*Dies — unsupported conclusion:* “strict freshness is free” overstates the result. The write-path cost of the protocol is not quantified by these read comparisons.

=== Variation: keep the cache correct through churn

*Observed duration:* a sustained 1801 s phase crosses 6.0 real TTL boundaries with no wrong reads. Hard expiry never fires; early expiration does.
#claim-badge("v3-01-churn-crosses-real-ttl")

*Observed comparison:* at 601 s, write-through and write-aside each cross 2.0 real TTL boundaries with no wrong reads; neither is declared faster. One cache-aside cell recorded an impossible value. Its preliminary verification passed, but final cell acceptance correctly aborted.
#claim-badge("v3-02-churn-600-through-vs-aside")

*Thrives — scoped evidence:* the tested coordination survives sustained churn for acknowledged-write freshness.

*Dies — observed correctness limit:* zero stale-after-ack reads does not guarantee every returned value was a committed state. The final gate rejected the anomalous cell, but the cause remains unexplained. End-of-phase ledger checks do not reconstruct every intermediate commit.

=== Variation: count cache resources and locality honestly

*Observed accounting:* the database-only and equal-total comparisons carry separate labels and budgets. Their repeated figures overlap trial spread; neither framing is pooled with the other.
#claim-badge("v3-03-equal-total-framing-labelled")

*Single-run locality result:* the donor-local layout reads faster than hashed distribution: warm reads 882 versus 673 ops/s, cacheable reads 852 versus 520. Both record no wrong reads.
#claim-badge("v3-06-colocated-vs-noncolocated-locality")

*Thrives — inferred fit:* fewer distributed lookups may help a donor page. The observed comparison supports engine key locality on one host; it does not prove physical placement.

Theory: #link(<theory-placement>)[Placement and resource controls].

#gap[
  Physical tablet/leader placement is unverified in this case. The layout labels describe intent, and all nodes share one host. Only the #link(<case-charity-placement>)[Charity placement variation] records tablets and leaders.
  #claim-badge("v4-gap-01-physical-placement-unverified")
]

#gap[
  Endpoint coverage is partial: earlier clusters use one endpoint, and this study records node count without a per-endpoint balance. The Charity SELECT cluster runs record a balanced per-endpoint count.
  #claim-badge("v4-gap-02-endpoint-distribution-partial")
]

#gap[
  Hard expiry, larger-than-harness scale, physical placement and open-loop arrivals remain unmeasured regimes. These are absent evidence, rather than null results.
  #claim-badge("v3-gap-01-study05-remaining-regimes")
]

*Reproduce:* `studies/05-cache-consistency`. Follow badges for committed inputs and run-specific budgets. Redis is a cache in this study, not an authoritative datastore comparison.
