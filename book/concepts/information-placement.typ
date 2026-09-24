#import "../lib/config.typ": *
#import "../lib/evidence.typ": registry-card

#heading("Information placement")

#marker("concept", "information placement")
Placement is the decision the rest of the book hangs on: which single fact owns which question,
and what is allowed to exist only because it answers a question faster. Every design in the corpus
is a different answer to "where does this fact live, and who maintains it".

#heading(level: 2, "The ladder")
#list(
  [*Derive it.* Answer the question from normalized facts at read time. Nothing to maintain,
   nothing to go stale; the cost is paid per read.],
  [*Index an existing fact.* Add the structure that turns a scan into a lookup. The cheapest
   derived object, but it taxes every write and only helps questions the schema can already
   express.],
  [*Materialize it.* Store the aggregate on the parent. Big read wins, a maintenance cost on every
   child write, and a correctness question whenever a write path forgets it.],
  [*Copy it down (rolldown).* Put parent information on the child to avoid a join. Measured as a
   package with its indexes, never as a bare column.],
  [*Copy it out (cache).* Serve the answer from beside the database, outside the authoritative
   transaction. It can be the lowest-latency hot read for the measured cacheable workload, and it is
   the step that crosses the transaction boundary — so the design acquires a separate freshness and
   failure-domain contract. Indexes, rollups and rolldowns carry their own maintenance correctness;
   what is special here is where the copy lives, not that it is copied.],
)

#figure-evidence(
  "../assets/fig-placement-ownership.svg",
  "structure",
  "conceptual illustration",
  "Placement and ownership: where an answer may live.",
  [The leaf row stays canonical for row questions. *Derive stores no additional copy.* The remaining
   rungs introduce some stored derived structure or copy, each owning its own maintenance cost and its
   own correctness contract: *index* adds a copy of an ordering, *materialize* copies the aggregate
   upward, *copy down* copies a parent key onto the child, and *copy out* serves the answer beside the
   database under a freshness contract. The ladder is ordered by how much read work each stored
   structure removes. The figure is a map of the corpus's decisions and carries no measured value.],
)

#direct[
  Ask what the schema must be able to *answer* before asking how fast it answers. Two measured
  answerability failures make this concrete: refunds were unanswerable in 14 of 15 ticketing
  designs, and no reserved-seating design could size a hold window because releasing an expired
  hold erases the evidence that it existed.
]
#registry-card("v2-07-refund-answerability-and-storage")
#registry-card("v2-gap-02-hold-funnel-unanswerable")

#heading(level: 2, "Boundaries")
#gap[
  Placement says nothing about *physical* colocation. The corpus does not verify where a cluster
  put its data, so a placement claim here is a design claim, not a topology result.
]
