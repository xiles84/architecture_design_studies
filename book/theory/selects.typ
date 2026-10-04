#import "../lib/config.typ": *
#import "../lib/teaching.typ": strategy, flow-picture
#set heading(numbering: none)
#heading(level: 2, numbering: none, "Database SELECTs")
Two SQL expressions are comparable only if they return the same logical answer. Establish grain, duplicates, NULLs, ties and ordering first.

Strategies: #link(<theory-select-contract>)[The answer contract] · #link(<theory-select-membership>)[Membership and absence] · #link(<theory-select-aggregate>)[Aggregation and per-parent top-N] · #link(<theory-select-pagination>)[Pagination and access paths].

#strategy("The answer contract", "theory-select-contract", [Define the answer before choosing the SQL spelling.],
  [A charity needs either one row per donor or one row per donation. Joining children changes row multiplicity unless the query restores the intended grain.],
  ([Name the output grain.], [Specify duplicates, NULLs, ties and order.], [Check answers against an independent oracle.],),
  ([A query has an explicit business question.], [Edge cases are part of verification.], [Alternative spellings preserve the contract.],),
  ([An inner join silently drops required parents.], [Deduplication removes meaningful repeated facts.], [A faster query answers a narrower question.],),
  ([Work required by the answer contract.], [None unless a maintained answer is introduced.], [Intermediate rows or maintained results.], [Validate all alternatives against the same oracle.], [Keep wrong-answer controls.],),
  flow-picture(([Business question], [Answer contract], [Oracle-checked result],), "The answer contract: the operation path.", [Business question → Answer contract → Oracle-checked result. Arrows show logical order, not a measured duration.]),
  ([Same logical answer across alternatives.], [Base facts and query results.], [Execution resources; not a new storage invariant.], [None for plain query spelling.], [The returned relation.], [Source snapshots or external copies.], [Unspecified grain, multiplicity, NULLs, ties or order.], [Compare edge-case outputs before timing.],), "case-coverage", measured: false)

#strategy("Membership and absence", "theory-select-membership", [Use a membership test when the question asks whether a related row exists.],
  [A charity asks for donors with a donation, not every donor-donation combination. Absence queries must also define NULL behavior.],
  ([Use existence semantics for a yes/no relationship.], [Check join duplication before adding DISTINCT.], [Test absence with nullable inputs.],),
  ([Only membership or absence is needed.], [Child fan-out is high.], [NULL handling is explicit.],),
  ([The output requires child details.], [NOT IN receives NULL without an explicit contract.], [An equivalent spelling is assumed faster everywhere.],),
  ([Plan-dependent probes or scans.], [No extra maintained state.], [Transient plan structures.], [Membership must preserve NULL and duplicate semantics.], [Inspect plans under measured fan-out regimes.],),
  flow-picture(([Parent candidate], [Related-row test], [Keep or reject parent],), "Membership and absence: the operation path.", [Parent candidate → Related-row test → Keep or reject parent. Arrows show logical order, not a measured duration.]),
  ([Return exactly the required parent membership.], [Parent and child rows.], [Read execution resources.], [None for equivalent SQL spellings.], [Membership answer.], [The selected snapshot can differ from later writes.], [Wrong NULL behavior or changed output grain.], [Sweep fan-out with matched answers and negative controls.],), "case-coverage", measured: false)

#strategy("Aggregation and per-parent top-N", "theory-select-aggregate", [Choose between repeated probes and shared processing according to the requested parent set.],
  [A page requests a total or latest donation for each donor. Processing all donations may waste work when the page asks for few donors.],
  ([Define empty-parent and tie behavior.], [Compare a per-parent probe with a shared aggregate.], [Vary the requested parent set while holding facts fixed.],),
  ([A narrow request has a useful lookup path.], [Shared processing amortizes across many requested parents.], [The selected plan matches the request shape.],),
  ([Every parent causes an expensive repeated scan.], [Pre-aggregation processes mostly unrequested data.], [A top-one query silently discards tied answers.],),
  ([Per-parent probes or shared processing.], [No extra writes for read-time aggregation.], [Intermediate aggregates, sorting or ranking.], [Empty groups and ties need an explicit rule.], [Measure several regimes and query plans.],),
  flow-picture(([Requested parents], [Probe or process together], [Totals or top-N],), "Aggregation and per-parent top-N: the operation path.", [Requested parents → Probe or process together → Totals or top-N. Arrows show logical order, not a measured duration.]),
  ([Aggregate and top-N answers preserve empty groups and ties.], [Child facts plus transient processing.], [Query execution resources.], [None until a stored rollup is introduced.], [Read-time summaries and ranks.], [Maintained aggregates, if added separately.], [An unfavorable requested-parent count or access path.], [Sweep parent count with oracle checks and shuffled measurement order.],), "case-coverage", measured: false)

#strategy("Pagination and access paths", "theory-select-pagination", [A cursor can continue from an ordered key; OFFSET skips through an ordered prefix.],
  [A donation feed needs the next page. Deep page numbers and a stable continuation cursor impose different navigation contracts.],
  ([Define a total order with a tie-breaker.], [Use a cursor carrying the last ordering keys.], [Compare equivalent pages at several depths.],),
  ([Navigation continues from the previous page.], [Ordering keys have a matching access path.], [Cursor semantics match concurrent change behavior.],),
  ([Users need arbitrary page-number jumps.], [Order lacks a stable tie-breaker.], [A shallow-page result is extrapolated to all depths.],),
  ([Continue after a key or scan and discard a prefix.], [Normal index maintenance.], [Cursor state plus access-path entries.], [Define behavior when rows appear, move or disappear.], [Version cursor formats and bound page sizes.],),
  flow-picture(([Ordered rows], [Last-key cursor], [Next matching page],), "Pagination and access paths: the operation path.", [Ordered rows → Last-key cursor → Next matching page. Arrows show logical order, not a measured duration.]),
  ([Return the correct next page under the chosen order.], [Rows, index and cursor keys.], [Query execution resources.], [Only ordinary index maintenance.], [Continuation cursor.], [A cursor can refer to rows changed since its creation.], [Unsupported jumps or an unfavorable engine/depth regime.], [Verify identical pages and sweep depth by engine.],), "case-coverage", measured: false)
