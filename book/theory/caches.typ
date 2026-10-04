#import "../lib/config.typ": *
#import "../lib/teaching.typ": strategy, flow-picture
#set heading(numbering: none)
#heading(level: 2, numbering: none, "Caches")
A cache is another copy of an answer. Its read benefit and its freshness contract are separate questions.

#figure-evidence("../assets/fig-cache-read-copy.svg", "structure", "conceptual illustration", "Caches: the choices at a glance.", [The picture maps the mechanisms. Layout and arrows do not represent measured throughput.])

Strategies: #link(<theory-cache-aside>)[Cache-aside and freshness] · #link(<theory-cache-invalidation>)[Invalidation and stale fills] · #link(<theory-cache-fencing>)[Leases and fencing].

#strategy("Cache-aside and freshness", "theory-cache-aside", [Read the copy when usable; otherwise load the source under an explicit freshness contract.],
  [A charity page reuses a donor profile. If a write completes, must the very next read see it, or may the copy lag?],
  ([Define what counts as an acceptable stale answer.], [Load misses from the authoritative source.], [Bound copy lifetime and record freshness violations.],),
  ([Repeated reads reuse the same answers.], [A defined staleness window is acceptable.], [Hit behavior is measured alongside mutation behavior.],),
  ([The required contract is immediate freshness.], [Keys rarely repeat.], [Hit throughput hides post-write stale reads.],),
  ([Hits avoid source work; misses still load it.], [Maintain the source and any copy protocol.], [Additional cached values.], [Freshness is a separate contract from source correctness.], [Bound memory, expiry and stampede behavior.],),
  figure-evidence("../assets/fig-freshness-timeline.svg", "structure", "conceptual illustration", "Cache-aside and freshness: mechanism and trade-off.", [The diagram shows ownership, work or coordination. It describes the mechanism, not a measured gain.]),
  ([Returned copies satisfy the declared freshness contract.], [Authoritative source plus cache.], [Popular cache keys and source misses.], [Copy population and maintenance.], [Cached answer.], [The cache after source changes.], [Low reuse or a stricter freshness requirement.], [Measure stale-after-ack reads as well as throughput.],), "case-cache", measured: true)

#strategy("Invalidation and stale fills", "theory-cache-invalidation", [Deleting a cache key does not order a reader that loaded old source data before the write.],
  [A reader starts a miss load. A writer commits and invalidates. The old reader then fills the cache with its earlier value.],
  ([Identify the fill that can outlive invalidation.], [Treat miss loading as part of the write/copy protocol.], [Validate the race with an adversarial schedule.],),
  ([The contract tolerates the uncovered race.], [A version or shared protocol rejects old fills.], [Every writer participates.],),
  ([Delete-after-write is claimed to guarantee freshness.], [A fill can publish after ownership changes.], [Negative controls never expose the intended race.],),
  ([Loads may be retried or rejected.], [Coordinate invalidation with fill publication.], [Version or coordination metadata where used.], [Source commit and copy publication must be ordered.], [Test slow fills and skipped writer paths.],),
  figure-evidence("../assets/fig-cache-stale-fill.svg", "structure", "conceptual illustration", "Invalidation and stale fills: mechanism and trade-off.", [The diagram shows ownership, work or coordination. It describes the mechanism, not a measured gain.]),
  ([An obsolete load cannot violate the promised freshness boundary.], [Source, cache and in-flight loader.], [Publication of the same cache key.], [Coordination replaces unchecked copy publication.], [The cached value.], [A late fill loaded before a source mutation.], [Uncoordinated loaders or writers.], [Pause an old load, commit a write, then release the fill.],), "case-cache", measured: true)

#strategy("Leases and fencing", "theory-cache-fencing", [A lease allows temporary work; a fencing token lets the receiver reject an obsolete worker.],
  [A slow cache loader loses its lease, but still finishes. The receiver must decide whether its result is authorized to publish.],
  ([Issue an ownership token with the lease.], [Attach it to attempted publication.], [Reject obsolete tokens at the receiver.],),
  ([The receiver enforces token order.], [Ownership can change while work is in flight.], [All publishers use the same protocol.],),
  ([Expiry alone is treated as publication protection.], [One path bypasses receiver checks.], [The token does not order competing owners.],),
  ([May wait for authorized loads.], [Issue ownership and validate publication tokens.], [Lease and fencing metadata.], [The receiver must reject obsolete owners atomically.], [Track stalled workers and protocol bypasses.],),
  figure-evidence("../assets/fig-lease-vs-fence.svg", "structure", "conceptual illustration", "Leases and fencing: mechanism and trade-off.", [The diagram shows ownership, work or coordination. It describes the mechanism, not a measured gain.]),
  ([An obsolete owner cannot publish an accepted result.], [Lease authority and fenced receiver.], [Ownership and publication record.], [Validation moves to publication.], [Ownership token and cached answer.], [Results produced after a worker loses authority.], [Unfenced or partial publication paths.], [Delay an owner past transfer and attempt stale publication.],), "case-cache", measured: true)
