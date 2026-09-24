#import "lib/config.typ": *

#heading("Glossary")
#marker("appendix", "glossary")
The single definition site for the vocabulary the rest of the book uses. Where a term is also a
registry field, the registry's value is named in brackets.

#table(
  columns: (auto, 1fr),
  stroke: 0.4pt + palette.rule,
  inset: 4pt,
  align: (left, left),
  table.header([*Term*], [*Meaning*]),
  [claim], [One reviewed statement about what the corpus measured, registered in the active evidence package and identified by an id such as `v2-16`. A number may not appear in this book without one.],
  [strength], [How strongly a claim is supported: `single_run_directional` (one run, one trial per cell, a direction not a curve), `repeated_controlled` (repeated fresh loads with a control), `mechanism_supported` (the mechanism is shown, the rate is not), `gap` (an explicit absence of evidence).],
  [direct evidence], [A callout whose numbers come from the cited run of the claim beside it.],
  [mechanism evidence], [A card or callout that demonstrates how something works without establishing a rate.],
  [analogy], [A labelled transfer beyond the measured family or environment. Never a measurement.],
  [gap], [An explicit absence of evidence, of one of four kinds. A gap is never a null result.],
  [schema limitation], [The representation cannot answer the question, so no run would. The hold-funnel gap is one.],
  [coverage gap], [A relevant regime has not been measured yet, but could be.],
  [instrument gap], [The property could not be observed with the tools the study had. Physical tablet placement is one.],
  [unstable measurement], [The number moved between runs with no code change, so no value is quotable.],
  [partially superseded], [Part of the claim has been replaced by newer evidence; the card says what moved to which successor and what is still open.],
  [retired], [The claim is no longer active, is not cited, and names the successor that closes it.],
  [run], [One execution of the harness under the benchmark lock, identified by a run tag and an inputs digest. Results are comparable only within one environment.],
  [replication], [A whole-run repeat of the same cell. The cards report it as `whole_run_replications`; `Trials: 1` means exactly one.],
  [trial], [A measurement inside a run. A fresh-load trial repeats the load phase; an in-run repeat is another measurement window inside one process. Neither is an independent replication.],
  [read score], [A blended within-run composite of a family's read queries. It compares variants inside one run and nothing else: it is not the throughput of any single query and is not quotable across runs.],
  [wrong_reads], [The harness counter for reads that violated the strict comparator in force, not for corruption. A relaxed cell can record many without serving an impossible value.],
  [strict after acknowledgement], [`begin(R1) > ack(W1)` implies the read must not return a version older than `W1`, unless a later write to the key supersedes it. This is the comparator the harness counts against.],
  [relaxed freshness], [A family of named weakenings of the strict line: session/read-your-writes, bounded staleness with a declared delta, eventual freshness, and best effort. Bounded staleness is one member, not the definition.],
  [resource framing], [`db-only` counts the database's resources; `equal-total` gives the cache a share of an equal budget. The two are never pooled, and a single-node versus cluster comparison answers "add machines", not "same budget".],
  [colocation], [A design intent that a key's rows share a partition. Logical mapping can follow from the schema; verified physical placement requires an instrument the studies did not have.],
  [rollup], [A parent aggregate maintained from its children. Where it lives in the same database it is updated in the same transaction; otherwise it is published, with a named delivery and reconciliation path.],
  [rolldown], [Parent information copied onto the child to avoid a join.],
  [embedding], [A child collection stored inside its parent row, as JSON or an array. In this corpus that means a relational JSON column, not a native document database.],
  [source correctness], [A mechanism that makes the authoritative mutation itself correct, such as a lock plus invariant recheck or a conditional update.],
  [copy correctness], [A mechanism that orders what leaves the transaction: invalidation, a publication fence, version validation at the hit boundary, a fill lease, a TTL, or an outbox.],
  [invalidation], [Removing or marking a cache entry stale. It orders nothing about a read or a fill already in flight.],
  [publication fence], [The check that refuses a fill whose source moved after the reader took its snapshot. This, not invalidation, is what closes the stale-fill race.],
  [confound], [A second change that moved with the one being measured. Registered in the package's confound register and published beside the numbers it affects.],
  [family, variant, scenario], [A family is a way of deciding where a fact lives; a variant is one design inside a family; a scenario is a workload built around one question, such as a hot drop or a refund.],
)
