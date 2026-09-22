#import "../lib/config.typ": *
#import "../lib/evidence.typ": registry-card

#heading("Indexes and hot rows")

#marker("concept", "indexes and hot rows")
Two different tools get confused because both look like "make the query fast". An index changes how
many rows the engine touches; a hot-row decision changes who may touch one row at the same time.

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
