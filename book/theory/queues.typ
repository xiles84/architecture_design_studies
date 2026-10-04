#import "../lib/config.typ": *
#import "../lib/teaching.typ": strategy, flow-picture
#set heading(numbering: none)
#pagebreak(weak: true)
#heading(level: 2, numbering: none, "Queues")
Queue theory extends the same ideas about ownership, retries and copies. This repository has no registered queue throughput study.

Strategies: #link(<theory-queues>)[Asynchronous handoff] · #link(<theory-queue-idempotency>)[Idempotency and duplicate delivery].

#strategy("Asynchronous handoff", "theory-queues", [A queue separates acceptance of work from completion of that work.],
  [A service accepts a notification job. The user-facing response and the worker's later delivery are different completion events.],
  ([Persist the handoff under a defined durability rule.], [Let a worker claim and process the job.], [Acknowledge according to the retry contract.],),
  ([Deferred completion is acceptable.], [Bursts can be absorbed within a bounded backlog.], [Failures have an explicit retry policy.],),
  ([The request requires an immediate final answer.], [An unbounded backlog hides overload.], [Acknowledgment occurs before durable required effects.],),
  ([Read status instead of assuming completion.], [Persist and process the handoff.], [Backlog, payloads and delivery state.], [Accepted, processed and completed are distinct states.], [Bound backlog; monitor age and failed jobs.],),
  flow-picture(([Accept durable job], [Worker claims work], [Effect and acknowledgment],), "Asynchronous handoff: the operation path.", [Accept durable job → Worker claims work → Effect and acknowledgment. Arrows show logical order, not a measured duration.]),
  ([An accepted job follows its stated completion contract.], [Producer state, queued job and consumer effect.], [Job claiming and constrained consumers.], [Request-path work moves to asynchronous processing.], [Job status and downstream projections.], [Completion status or downstream copies.], [Arrival demand exceeds sustainable service capacity.], [Measure offered load, backlog age, failures and completed work.],), "case-coverage", measured: false)

#strategy("Idempotency and duplicate delivery", "theory-queue-idempotency", [A retry must not repeat a business effect that was already accepted.],
  [A worker applies a job's effect, then crashes before acknowledging it. Redelivery can repeat the operation.],
  ([Give the logical operation a stable identity.], [Make effect and duplicate detection atomic where possible.], [Define behavior for external effects outside that transaction.],),
  ([Duplicate delivery is expected.], [Effects can be keyed or conditionally applied.], [Retry behavior is part of the business contract.],),
  ([Acknowledge-after-effect is assumed to prevent duplicates.], [Deduplication state expires before possible redelivery.], [External effects are repeated without their own protection.],),
  ([Check status or operation identity.], [Record effect and deduplication outcome.], [Operation identities and retained outcomes.], [Retries must preserve the logical operation's effect.], [Bound retention and handle poison jobs.],),
  flow-picture(([Same job redelivered], [Check operation identity], [Apply once or return prior outcome],), "Idempotency and duplicate delivery: the operation path.", [Same job redelivered → Check operation identity → Apply once or return prior outcome. Arrows show logical order, not a measured duration.]),
  ([Repeated delivery preserves one accepted logical effect.], [Job identity, effect and deduplication record.], [The same logical operation identity.], [Duplicate detection occurs at effect time.], [Delivery and completion status.], [A status record separated from its actual effect.], [Non-atomic effect tracking or unsafe external side effects.], [Crash between effect and acknowledgment, then redeliver.],), "case-coverage", measured: false)
