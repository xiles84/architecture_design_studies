#import "../lib/config.typ": *
#import "../lib/evidence.typ": registry-card

#heading("Major families")

#marker("chapter", "major families")
A family is a way of deciding where a fact lives. The six here are fixed by the goal; each section
follows the same progressive structure, and the mechanism chapters carry the reusable explanations.
Every cost below is a measured claim, labelled *direct* when it comes from the cited run or
*analogy* when it is a transfer.

// ---------------------------------------------------------------- 1
#heading(level: 2, "Normalized facts")
#heading(level: 3, "Quick choice in plain language")
Keep each fact in one place and let queries assemble the answer. Choose it unless a measured
question says otherwise.
#heading(level: 3, "Where it thrives and where it perishes")
Thrives when questions are answerable from the rows themselves and writes are spread across many
parents. Perishes on a hot row, or when the required question was never modelled.
#heading(level: 3, "Common scenarios")
Order lines, ticket sales, configuration entries — anything that must never overbook or double-sell.
#heading(level: 3, "Mechanism")
See the *indexes and hot rows* concept and the *concurrency control* concept. Indexes answer expressible questions; arbitration
answers contended rows.
#heading(level: 3, "Costs")
Reads pay per query; writes are cheap; storage is the baseline; correctness depends on the
constraint you actually declared; operations are ordinary.
#heading(level: 3, "Direct evidence")
#registry-card("v2-01-normalized-index-first")
#registry-card("v2-06-hot-drop-designs")
#registry-card("v2-12-lock-vs-version-upper-bound")
#heading(level: 3, "Boundaries and reproduction")
The ratios are single-run and single-host. Reproduce from the run tags in the cards; check the
`run_level` note for failed cells before quoting a percentage.

// ---------------------------------------------------------------- 2
#heading(level: 2, "Rolldown / duplicated data")
#heading(level: 3, "Quick choice in plain language")
Copy a parent value onto the child to remove a join. Do it when the copied value is small, stable
enough, and the join is on the hot path.
#heading(level: 3, "Where it thrives and where it perishes")
Thrives on read-heavy workloads with a bounded parent set. Perishes when the copied value changes
often, because every change becomes a fan-out write.
#heading(level: 3, "Common scenarios")
A donation row carrying its donor's display name; an order line carrying its product's category.
#heading(level: 3, "Mechanism")
See the *information placement* concept. The measured pair is a package — the copied key plus its SQL and
indexes — so attribute the gain to the package, not to the column alone.
#heading(level: 3, "Costs")
Reads improve; writes gain an index and a copy to maintain; storage grows; correctness needs a rule
for the stale copy; operations must backfill after a schema change.
#heading(level: 3, "Direct evidence")
#registry-card("v2-05-d3-read-advantage-repeated")
#heading(level: 3, "Boundaries and reproduction")
One family contrast is not a general law. The copied-key package is confounded (`conf-01`), and
placement is unverified (`v2-gap-05`).

// ---------------------------------------------------------------- 3
#heading(level: 2, "Rollups / materialized aggregates")
#heading(level: 3, "Quick choice in plain language")
Store the total on the parent. Use it when an aggregate is read far more often than its children
change, and be explicit about who maintains it.
#heading(level: 3, "Where it thrives and where it perishes")
Thrives on dashboards and parent-level reads. Perishes when child writes dominate, or when a write
path can forget the rollup.
#heading(level: 3, "Common scenarios")
Donor lifetime totals, event ticket counts, a configuration portal's overview.
#heading(level: 3, "Mechanism")
See the *derived state and history* concept. Trigger-maintained is a structural guarantee at a write cost;
application-maintained is cheaper but must be atomic with the publication.
#heading(level: 3, "Costs")
Reads win big; writes pay maintenance; storage grows modestly; correctness is the whole game;
operations must rebuild the rollup after a bulk load.
#heading(level: 3, "Direct evidence")
#registry-card("v2-02-rollup-family")
#registry-card("v2-13-trigger-rollup-cost")
#registry-card("v2-04-recency-maintained-index")
#heading(level: 3, "Boundaries and reproduction")
The recency pair's write cost is confounded (`conf-02`); quote its read result, not its write cost.

// ---------------------------------------------------------------- 4
#heading(level: 2, "Embedded documents")
#heading(level: 3, "Quick choice in plain language")
Put the children inside the parent when they are read together and bounded in size. It is a
read-locality choice, not a database-family choice.
#heading(level: 3, "Where it thrives and where it perishes")
Thrives when a bounded child set is always read with its parent. Perishes when children are queried
independently or grow without bound.
#heading(level: 3, "Common scenarios")
A reservation's seat map, a configuration file rendered whole.
#heading(level: 3, "Mechanism")
See the *information placement* concept. A section document answered a 100k-seat map 4.5–11x faster than
per-seat rows on YugabyteDB; a sharded variant collapsed point lookups to 26 ops/s.
#heading(level: 3, "Costs")
Reads of the whole child set improve; point access to one child can get much worse; storage is
similar; correctness needs a bounded-growth rule; operations can be awkward for partial updates.
#heading(level: 3, "Direct evidence")
#registry-card("v2-03-embedding-not-ahead")
#registry-card("v2-11-section-seat-map")
#heading(level: 3, "Boundaries and reproduction")
Embedding here means a JSON/array child column in PostgreSQL or a document row in YugabyteDB — a
labelled #text(weight: "bold")[analogy] to a native document database, which was not measured
(`v2-gap-06`). The row-versus-document contrast is a package comparison (`conf-03`).

// ---------------------------------------------------------------- 5
#heading(level: 2, "Append-only history / snapshots")
#heading(level: 3, "Quick choice in plain language")
Never overwrite the fact; append what happened. Choose it when a question must survive an undo.
#heading(level: 3, "Where it thrives and where it perishes")
Thrives wherever refunds, cancellations or corrections are in scope. Perishes under storage limits
with no retention policy.
#heading(level: 3, "Common scenarios")
Refund reporting, audit trails, last-purchase questions after a cancellation.
#heading(level: 3, "Mechanism")
See the *derived state and history* concept and the *expiry and clock authority* concept.
#heading(level: 3, "Costs")
Reads need a report over history; writes append and fence; storage grows (about 33–35% more in the
ticketing ledger); correctness *improves* — this family is the answerability fix; operations need
retention.
#heading(level: 3, "Direct evidence")
#registry-card("v2-07-refund-answerability-and-storage")
#registry-card("v2-08-ledger-race-cost-128-buyers")
#heading(level: 3, "Boundaries and reproduction")
The PostgreSQL ledger cost is order-sensitive and unsettled (`v2-09`); the load multiplier is not
published (`v2-gap-01`). Quote the 128-buyer scoped result, never a pooled ratio.

// ---------------------------------------------------------------- 6
#heading(level: 2, "External read copies / caches")
#heading(level: 3, "Quick choice in plain language")
Serve the answer from beside the database, with a stated freshness contract. Do it for warm,
cacheable, read-heavy questions.
#heading(level: 3, "Where it thrives and where it perishes")
Thrives on repeated reads of stable keys. Perishes when many instances share a cache under relaxed
freshness without a fence.
#heading(level: 3, "Common scenarios")
A donor portal's hot profile reads, a public dashboard over changing totals.
#heading(level: 3, "Mechanism")
See the *cache consistency* concept. Invalidation is a fence that refuses a stale publication, not a deletion.
#heading(level: 3, "Costs")
Reads gain 2.2–3.5x and p99 drops an order of magnitude; writes may gain a fence; storage moves to
the cache (whose budget is outside the comparison); correctness needs the fence; operations gain a
new failure domain.
#heading(level: 3, "Direct evidence")
#registry-card("v2-14-cache-throughput-gain")
#registry-card("v2-15-three-instance-staleness")
#registry-card("v2-16-strict-freshness-read-cost")
#heading(level: 3, "Boundaries and reproduction")
Redis was measured as a cache, never as an authoritative store; the result is `db-only` framed
(`conf-04`). The three-instance failure is a mechanism demonstration, not a rate to predict from.
