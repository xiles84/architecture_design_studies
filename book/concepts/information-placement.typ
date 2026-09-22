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
  [*Copy it out (cache).* Serve the answer from beside the database. The fastest read and the only
   step that adds a consistency contract to the design.],
)

#heading(level: 2, "The rule that comes before performance")
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
