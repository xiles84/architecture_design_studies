#import "../lib/config.typ": *
#import "../lib/evidence.typ": registry-card

#heading("Major families")

#marker("chapter", "major families")
The six major families are fixed by the goal. Each section follows the same progressive structure;
prose is written by `book-synthesis-v1`, evidence is attached here.
#let families = (
  ("normalized facts", "v2-01-normalized-index-first", "v2-06-hot-drop-designs", "v2-12-lock-vs-version-upper-bound"),
  ("rolldown / duplicated data", "v2-05-d3-read-advantage-repeated", "v2-gap-05-colocation-unverified"),
  ("rollups / materialized aggregates", "v2-02-rollup-family", "v2-04-recency-maintained-index", "v2-13-trigger-rollup-cost"),
  ("embedded documents", "v2-03-embedding-not-ahead", "v2-11-section-seat-map"),
  ("append-only history / snapshots", "v2-07-refund-answerability-and-storage", "v2-08-ledger-race-cost-128-buyers", "v2-09-postgres-ledger-ratio-unsettled"),
  ("external read copies / caches", "v2-14-cache-throughput-gain", "v2-15-three-instance-staleness", "v2-16-strict-freshness-read-cost"),
)
#for (name, ..ids) in families [
  #heading(level: 2, name)
  #progressive-structure(depth: 3)
  #for id in ids { registry-card(id) }
]
