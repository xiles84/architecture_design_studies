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

#heading(level: 2, "The confounds, published here rather than in a footnote")
#list(
  [*`conf-01`* — the D2→D3 rolldown prices a copied key plus its SQL and indexes, not the column.],
  [*`conf-02`* — the recency pair's write axis is confounded by the whole rollup package.],
  [*`conf-03`* — the row-versus-document seat map is a family/package comparison.],
  [*`conf-04`* — the cache gain is `db-only` framed.],
  [*`conf-05`* — the ticketing buyer arms disagree in sign on PostgreSQL.],
  [*`conf-06`* — the strict-freshness range must be shown in full (−14.2% to +12.3%).],
)
