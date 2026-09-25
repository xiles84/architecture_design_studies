#import "../lib/config.typ": *
#import "../lib/evidence.typ": registry-card

#heading("Indexes and hot rows")

#marker("concept", "indexes and hot rows")
Two different tools get confused because both look like "make the query fast". An index changes how
many rows the engine touches; a hot-row decision changes who may touch one row at the same time.

#figure-evidence(
  "../assets/fig-index-vs-hotrow.svg",
  "structure",
  "conceptual illustration",
  "An index narrows the scan; a hot row serialises the writers.",
  [Left, the index panel: a query follows a narrow index path to a small subset of rows, so the engine
   touches less — the index is a copy of an ordering and taxes every write to the indexed column. Right,
   the hot-row panel: several writers converge on one contended row, so they serialise and the design
   question is arbitration, not structure. The figure separates the two tools and asserts no rate.],
)

#heading(level: 2, "Indexes")
An index answers a question the schema can already express, at the cost of every write to the
indexed column. In the charity tree the first index was the first large read win, and flattening
the foreign key (a rolldown package) bought more still.

#heading(level: 2, "Hot rows")
When one row is the contention point — one event's ticket count, one product's stock — the design
question is arbitration, not structure. Pre-created seats taken by compare-and-set or `SKIP LOCKED`
sold 4–8x faster than any single counter row; the counter also blocks unrelated edits to the same
row, which is a correctness-adjacent cost the throughput number hides.

#heading(level: 2, "Direct evidence")
#registry-card("v2-01-normalized-index-first")
#registry-card("v2-06-hot-drop-designs")

#heading(level: 2, "Boundaries")
The hot-row result is one workload at small scale on one host. Read it as "the counter row
serialises everything", not as a capacity figure.

#heading(level: 2, "The checklist")

The eight questions, answered for this concept; the status in brackets is the evidence behind the
answer.

- *Invariant.* The answer must not change because an access path changed; an index may only change how
  the answer is found. [direct]
- *Where the information lives.* In the table. An index locates rows; it does not hold the answer. [direct]
- *Where the contention is.* The single hot row that a counter or an aggregate lives on — one event
  row, one parent. [direct]
- *What work moved from reads to writes.* Maintaining an index, or a counter, is write work bought to
  make reads cheap. [mechanism]
- *What state is derived.* Nothing, for an index. A counter or rollup column is derived, and belongs to
  the derived-state concept. [mechanism]
- *What can become stale.* A rollup column can disagree with its entries. An index cannot be stale, only
  unused. [mechanism]
- *Where it stops working.* On a table small enough that the whole relation is cached, an index stops
  paying — measured only as a control that could not fire at 3,600 rows. [mechanism]
- *The experiment that would settle it.* An identical-SQL index-versus-no-index pair, plus a writer
  sweep to show where the hot row takes over. [direct]
