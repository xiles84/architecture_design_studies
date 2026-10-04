#import "../lib/config.typ": *
#import "../lib/teaching.typ": strategy, flow-picture
#set heading(numbering: none)
#heading(level: 2, numbering: none, "Database schemas")
Choose where facts live before choosing how quickly to retrieve them. These strategies trade read work for storage, maintenance and coordination.

#figure-evidence("../assets/fig-families-er.svg", "structure", "conceptual illustration", "Database schemas: the choices at a glance.", [The picture maps the mechanisms. Layout and arrows do not represent measured throughput.])

Strategies: #link(<theory-normalized>)[Normalized rows] · #link(<theory-indexes>)[Indexes] · #link(<theory-rolldown>)[Rolldown] · #link(<theory-rollup>)[Rollup] · #link(<theory-embedding>)[Embedding].

#strategy("Normalized rows", "theory-normalized", [Keep each fact in its owning row; combine rows when an answer needs several owners.],
  [A donor belongs to a charity and has many donations. Repeating the same charity details on every donation creates extra update paths.],
  ([Give each entity a key.], [Store relationships as references.], [Join or aggregate when serving a combined answer.],),
  ([Facts change independently.], [Children grow without a small fixed bound.], [Integrity and selective updates matter.],),
  ([Nearly every read reconstructs an expensive aggregate.], [The same remote join dominates the workload.], [A proven read bottleneck justifies maintained copies.],),
  ([Joins and aggregation remain on the read path.], [Update the owning fact once.], [Separate rows plus relationship indexes.], [References must resolve to valid owners.], [Migrations preserve keys and relationships.],),
  figure-evidence("../assets/fig-placement-ownership.svg", "structure", "conceptual illustration", "Normalized rows: mechanism and trade-off.", [A read can derive from facts or use an index, aggregate, copied key or external copy. Every added copy needs maintenance.]),
  ([A fact has one authoritative owner.], [Entity rows linked by keys.], [Frequently changed rows and shared indexes.], [None; combined answers are computed when requested.], [Query results.], [Maintained external copies, if introduced.], [Repeated expensive reconstruction of the same answer.], [Compare matched indexes before adding a copied fact.],), "case-charity", measured: true)

#strategy("Indexes", "theory-indexes", [Store an access path so a query reaches the rows it needs without searching all rows.],
  [A donation feed filters by charity and orders by time. An index helps only when its keys support that question.],
  ([Match leading keys to filters and ordering.], [Read the actual plan, including heap visits.], [Maintain the index when indexed values change.],),
  ([Queries select a small useful subset.], [Filters and ordering match index keys.], [Read savings repay index maintenance.],),
  ([Most rows are needed anyway.], [Writes dominate and many indexes duplicate work.], [An index is blamed for a separate hot-row bottleneck.],),
  ([Fewer candidate rows for matching queries.], [Maintain each affected index.], [Extra key and row-location entries.], [An access path must preserve the answer.], [Review plans and retire unused indexes.],),
  figure-evidence("../assets/fig-index-vs-hotrow.svg", "structure", "conceptual illustration", "Indexes: mechanism and trade-off.", [An index shortens a read's search path. It does not remove the shared row on which concurrent writers collide.]),
  ([Return the same qualifying rows.], [Base table plus derived access paths.], [Index pages and hot base rows.], [Access-path maintenance.], [Index entries.], [A transactional index should track committed base rows.], [Low selectivity or excessive write maintenance.], [Match data and SQL, change one index, inspect plans.],), "case-charity", measured: true)

#strategy("Rolldown", "theory-rolldown", [Copy a parent key onto its children to shorten their read path.],
  [A donation lookup needs its charity, but the normalized route goes through the donor. A copied charity key removes that hop.],
  ([Copy the parent value when creating a child.], [Use it directly in child filters or indexes.], [Repair every affected child when the parent fact changes.],),
  ([Parent relationships rarely change.], [Child reads repeatedly need the parent key.], [Fan-out updates can be made correct.],),
  ([Parents change frequently.], [A parent has many children to repair.], [The gain actually comes from an unmatched index.],),
  ([May remove a relationship hop.], [Parent changes fan out to children.], [A copied field per child.], [Copied keys must agree with their authoritative parent.], [Backfill and repair copied information.],),
  figure-evidence("../assets/fig-rollup-vs-rolldown.svg", "structure", "conceptual illustration", "Rolldown: mechanism and trade-off.", [Rollup maintains a parent aggregate when children change. Rolldown maintains child copies when parent information changes.]),
  ([Each copied key matches its owner.], [Parent fact and child copies.], [Children touched by a parent change.], [Resolving the parent relationship.], [Copied parent information.], [Children skipped during maintenance.], [Large fan-out with frequent parent changes.], [Hold SQL and indexes fixed; vary the copied key.],), "case-charity", measured: true)

#strategy("Rollup", "theory-rollup", [Store an aggregate on its parent so reads skip the children and writes maintain the answer.],
  [A charity dashboard repeatedly sums donations. A stored total replaces that repeated scan with a parent lookup.],
  ([Keep child rows as the facts of record.], [Update the total with each contributing mutation.], [Choose trigger or application maintenance deliberately.],),
  ([Totals are read much more often than changed.], [The aggregate definition is stable.], [Maintenance can be atomic with the source change.],),
  ([Many writers collide on one parent total.], [Aggregate rules change frequently.], [Historical answers still require scanning the facts.],),
  ([Fetch a maintained total.], [Update shared parent state.], [Stored totals and supporting indexes.], [Every insert, correction and deletion must maintain totals.], [Rebuild or reconcile when rules change.],),
  figure-evidence("../assets/fig-rollup-vs-rolldown.svg", "structure", "conceptual illustration", "Rollup: mechanism and trade-off.", [Rollup maintains a parent aggregate when children change. Rolldown maintains child copies when parent information changes.]),
  ([The total equals its contributing facts.], [Child facts and parent aggregate.], [Shared aggregate row.], [Repeated aggregation.], [Totals and other maintained summaries.], [An aggregate with a missed mutation.], [Write-heavy hot parents or large replacement operations.], [Compare maintainers and sweep write size and read/write mix.],), "case-charity", measured: true)

#strategy("Embedding", "theory-embedding", [Store bounded child information with its parent when the application reads and changes that unit together.],
  [A product can read a configuration as one document. A growing donation history may require independent child access instead.],
  ([Choose the atomic document boundary.], [Read or rewrite the embedded unit.], [Bound growth and define concurrent update behavior.],),
  ([Parent and children are read together.], [The collection has a practical size bound.], [Whole-unit replacement matches the write pattern.],),
  ([Children grow without a useful bound.], [Small child edits rewrite a large unit.], [Writers contend on the same parent document.],),
  ([One unit can answer a combined read.], [May rewrite or lock a larger unit.], [Representation-dependent duplication and document overhead.], [Concurrent edits must not overwrite one another.], [Set size limits and migration rules.],),
  figure-evidence("../assets/fig-embedded-vs-normalized.svg", "structure", "conceptual illustration", "Embedding: mechanism and trade-off.", [Normalized children are independently addressable rows. Embedding places a bounded child collection inside the parent's update unit.]),
  ([The embedded unit represents its child facts correctly.], [Child information inside the parent.], [The parent document or row.], [Assembling the child collection.], [Only explicitly copied or summarized fields.], [Copies outside the authoritative document.], [Large units, independent child reads or concurrent small edits.], [Compare bounded sizes and matched read/write contracts.],), "case-charity", measured: true)
