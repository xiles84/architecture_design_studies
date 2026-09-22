#import "../lib/config.typ": *
#import "../lib/evidence.typ": registry-card

#heading("Minor variants")

#marker("chapter", "minor variants")
A minor variant is one controlled change inside a family: index selection, full versus bounded
embedding, trigger versus application maintenance, optimistic versus pessimistic concurrency,
locks/CAS/isolation, pre-created versus on-demand inventory, cache-aside versus write-through, and
strict versus relaxed freshness. Only the first and the concurrency variants have measured
controlled pairs in the corpus; the rest are labelled with their gap or analogy status.
#registry-card("v2-01-normalized-index-first")
#registry-card("v2-12-lock-vs-version-upper-bound")
#registry-card("v2-13-trigger-rollup-cost")
#registry-card("v2-16-strict-freshness-read-cost")
