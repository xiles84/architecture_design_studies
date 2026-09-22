#import "../lib/config.typ": *
#import "../lib/evidence.typ": registry-card

#heading("Evidence and reproduction")

#marker("chapter", "evidence and reproduction")
Every number in this book resolves to a claim in `book/evidence/v2/claims.json`, and every claim to
a measured cell in a cited report. This chapter explains the confidence cards, the run-tag
provenance model, and how to reproduce a run in Podman.
#heading(level: 2, "The provenance model")
#direct[
  A run tag marks the producing code state, not a tree that contains the run's results. Check out
  the tag for the code, then read the committed results directory for the data.
]
#registry-card("v2-gap-06-native-datastore-families")
#registry-card("v2-gap-01-ledger-load-multiplier-unstable")
