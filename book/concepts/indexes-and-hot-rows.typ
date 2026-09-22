#import "../lib/config.typ": *
#import "../lib/evidence.typ": registry-card

#heading("Indexes and hot rows")

#marker("concept", "indexes and hot rows")
An index is the cheapest derived structure, but it is not free: it charges every write and it does
not answer a question the schema cannot express. The hot-row case shows the arbitration bound.
#registry-card("v2-01-normalized-index-first")
#registry-card("v2-06-hot-drop-designs")
