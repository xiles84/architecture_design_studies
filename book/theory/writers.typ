#import "../lib/config.typ": *
#import "../lib/teaching.typ": strategy, flow-picture
#set heading(numbering: none)
#pagebreak(weak: true)
#heading(level: 2, numbering: none, "Concurrent writes and time")
Start with the business invariant: no lost increment, no oversold seat, or no stale owner publishing a result. A faster invalid design is a failed design.

Strategies: #link(<theory-hot-rows>)[Hot rows] · #link(<theory-optimistic>)[Optimistic version checks] · #link(<theory-pessimistic>)[Pessimistic locking] · #link(<theory-history>)[History and current state] · #link(<theory-expiry>)[Expiry and clock authority].

#strategy("Hot rows", "theory-hot-rows", [A shared mutable row concentrates otherwise independent writers into one coordination point.],
  [Many donations may update one charity total. More workers do not remove the need to order updates to that row.],
  ([Identify which key every operation changes.], [Separate useful parallel work from serialized work.], [Measure correctness, conflicts and tail behavior together.],),
  ([Contention is low or bounded.], [The shared row is the intended arbitration point.], [Rejected work is handled explicitly.],),
  ([Writers overwhelmingly target one key.], [Retries amplify already overloaded work.], [Throughput is reported without invariant checks.],),
  ([May wait behind writers.], [Serialized updates or conflicting attempts.], [Little extra state by itself.], [Acknowledged changes must survive.], [Track queueing, rejection and conflict behavior.],),
  flow-picture(([Many writers], [One mutable key], [Ordered successful changes],), "Hot rows: the operation path.", [Many writers → One mutable key → Ordered successful changes. Arrows show logical order, not a measured duration.]),
  ([Acknowledged changes preserve the business invariant.], [Shared mutable rows.], [The popular key or row.], [No automatic read/write relocation.], [Counters or summaries when present.], [Cached copies and abandoned read decisions.], [Concentrated offered work exceeds serialized capacity.], [Sweep skew and writers with a violating control.],), "case-ticket", measured: true)

#strategy("Optimistic version checks", "theory-optimistic", [Read a version, then commit only if that version is still current.],
  [Two administrators edit the same configuration. A conditional update detects whether another writer changed it first.],
  ([Read the state and its version.], [Write using a version condition.], [On conflict, retry from fresh state within a defined budget.],),
  ([Conflicts are uncommon.], [Attempts are cheap to repeat.], [Bounded rejection is acceptable.],),
  ([A hot key creates frequent conflicts.], [Side effects occur before the conditional commit.], [Retries repeat expensive work without a bound.],),
  ([Read state and version.], [Conditional writes; conflicts may require another attempt.], [A version token.], [Never accept an update based on an obsolete version.], [Bound retries and define client-visible rejection.],),
  flow-picture(([Read state and version], [Conditional update], [Commit or reject conflict],), "Optimistic writes: validate before acceptance.", [A writer reads a version and attempts an update conditional on that version. An intervening change makes the update fail; the writer must re-read before retrying.]),
  ([No acknowledged update overwrites an intervening change.], [Source row plus version.], [The version-checked hot key.], [None inherently.], [Version token.], [The state read before a conflicting commit.], [Frequent conflicts or expensive repeated attempts.], [Match the invariant and sweep concurrent writers.],), "case-config", measured: true)

#strategy("Pessimistic locking", "theory-pessimistic", [Acquire the arbitration lock before reading state used to decide an update.],
  [A seat reservation reads availability and changes it. Holding the lock keeps another writer from making the same decision concurrently.],
  ([Lock the row before deciding.], [Keep protected work short.], [Commit or roll back to release the lock.],),
  ([Conflicts would otherwise be frequent.], [Protected operations are short.], [Waiting has a bounded, visible policy.],),
  ([External work happens while holding the lock.], [Lock order causes deadlocks.], [Clients wait indefinitely on a busy row.],),
  ([Protected reads may wait.], [Queued updates instead of speculative conflicts.], [Lock-management state.], [Read, decision and write must share protection.], [Set timeouts and consistent lock ordering.],),
  flow-picture(([Acquire row lock], [Read and decide], [Write and release],), "Pessimistic locking: the operation path.", [Acquire row lock → Read and decide → Write and release. Arrows show logical order, not a measured duration.]),
  ([The protected decision observes state that cannot change underneath it.], [Authoritative rows and lock state.], [The locked row and waiting writers.], [No inherent relocation; work is ordered.], [None inherently.], [Copies read outside the protected transaction.], [Long critical sections or overload.], [Compare version checks and locks under one invariant.],), "case-config", measured: true)

#strategy("History and current state", "theory-history", [An event history records what happened; current state records the answer needed now.],
  [A booking service must answer whether a seat is free and explain refunds. A log alone need not make the current answer cheap.],
  ([Keep immutable events for the historical contract.], [Maintain or derive a current-state projection.], [Decide how events and projections commit together.],),
  ([Audit or reconstruction is required.], [Current reads have a maintained projection.], [Projection maintenance has a clear ownership rule.],),
  ([Every hot read replays a growing log.], [An event append is mistaken for concurrency protection.], [History and projection can diverge silently.],),
  ([Read a projection or replay relevant events.], [Append events and maintain any current projection.], [Retained history plus current state.], [Events and acknowledged current state must reconcile.], [Version events and support replay or repair.],),
  figure-evidence("../assets/fig-history-vs-current.svg", "structure", "conceptual illustration", "History and current state: mechanism and trade-off.", [Current state answers availability now. Retained events preserve past sales and refunds that current-state deletion would destroy.]),
  ([Current state agrees with accepted historical changes.], [Event history and optional projection.], [Append coordination and shared current-state rows.], [Projection maintenance replaces replay on every read.], [Current-state projection.], [A projection that falls behind or misses an event.], [Growing replay or an incorrectly maintained projection.], [Match concurrency strategy; compare replay and projection cost.],), "case-ticket", measured: true)

#strategy("Expiry and clock authority", "theory-expiry", [A deadline has meaning only when one authority decides whether it has passed.],
  [A seat hold expires while a client attempts confirmation. The service needs one atomic decision about deadline and ownership.],
  ([Choose the authoritative clock.], [Check expiry inside the arbitration operation.], [Separate expiration semantics from cleanup timing.],),
  ([Deadline checks use a defined authority.], [Confirmation and expiration share arbitration.], [Cleanup delay does not extend authorization.],),
  ([Clients decide expiry using unrelated clocks.], [Cleanup is the only protection against late confirmation.], [A lease is assumed to prevent stale publication.],),
  ([Check ownership and authoritative time.], [Arbitrate confirmations and expiry transitions.], [Deadline and ownership token.], [No success after the authority rejects the deadline.], [Monitor delayed cleanup and clock assumptions.],),
  figure-evidence("../assets/fig-expiry-clock.svg", "sequence", "conceptual illustration", "Expiry and clock authority: mechanism and trade-off.", [A database-time guard can reclaim an expired hold without waiting for cleanup. A sweeper-only design delays availability; a skewed client clock can reclaim too early.]),
  ([A hold can be confirmed only while valid for its owner.], [Hold row, deadline and owner token.], [The hold's arbitration record.], [No inherent relocation.], [Availability projection or cache, if present.], [Client beliefs and cached availability.], [Undefined clock authority or non-atomic transitions.], [Race confirmation against expiry with an invalid control.],), "case-seating", measured: true)
