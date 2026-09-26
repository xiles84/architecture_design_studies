#import "../lib/config.typ": *
#import "../lib/evidence.typ": registry-card

#heading("Cross-family comparisons")

#marker("chapter", "cross-family comparisons")
Comparisons across families are where over-claiming happens. The rule here: a pair is quoted only
when one decision changed, and every comparison carries its confounds on the page.

#heading(level: 2, "What a controlled pair looks like")
#list(
  [*One decision changed.* The charity tree's index and rolldown pairs each change one placement
   decision; their confounds are named (`conf-01`).],
  [*Same workload and resources.* Logical data, operations and per-node budget held constant.],
  [*Controls present and fired.* A pair whose negative control did not fire is weaker than it
   looks.],
)

#heading(level: 2, "Comparisons this book deliberately does not make")
#list(
  [*No cross-study numeric ranking.* No comparability record exists, so ratios from different runs
   are not pooled.],
  [*No equal-total-resource ranking.* Single-versus-cluster answers "add machines", not "same
   budget".],
  [*No native-family transfer.* JSONB/state is not a native document store and Redis was measured
   only as a cache; such a transfer is analogy, labelled as such.],
)

#heading(level: 2, "Direct evidence for the two family contrasts")
#registry-card("v2-03-embedding-not-ahead")
#registry-card("v2-11-section-seat-map")

#heading(level: 2, "What these six comparisons cannot decide")

Each item below is a comparison the book shows, followed by the reason it cannot settle the question it
looks like it settles. The reason comes first because it is the useful part; the short label in brackets
is only how the rest of the book refers back to it.

#list(
  [*A copied key carries more than the key.* The D2→D3 rolldown compares two whole designs, so the read
   gain it measures includes the copied column, the index that supports it and the changed SQL — not the
   copied column on its own. [*`conf-01`*]],
  [*The recency pair's write cost is not only about recency.* The design with the recency index also does
   all of the rollup maintenance, so the experiment cannot say how much of the write cost comes from the
   index itself. Quote its read comparison, not its write cost. [*`conf-02`*]],
  [*A table of rows against a stored document is a whole-family comparison.* The seat map differs both in
   storage shape and in how the data is organised, so any difference belongs to the two designs together,
   not to one decision inside them. [*`conf-03`*]],
  [*The cache gain was measured against one baseline.* The reference point was the database alone, so the
   measured gain cannot be attributed to the cache by itself. [*`conf-04`*]],
  [*The two buyer models swap places.* On PostgreSQL the two designs change order depending on the
   workload, so neither can be reported as the faster one. [*`conf-05`*]],
  [*The strict-freshness numbers run in both directions.* The range goes from −14.2% to +12.3%, so
   quoting either end alone would misrepresent the result; the whole range travels with any summary.
   [*`conf-06`*]],
)
