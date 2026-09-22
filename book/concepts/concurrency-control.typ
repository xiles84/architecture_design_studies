#import "../lib/config.typ": *
#import "../lib/evidence.typ": registry-card

#heading("Concurrency control")

#marker("concept", "concurrency control")
Optimistic (version-checked) and pessimistic (lock-before-update) strategies under the same business
invariant. Both local engine mechanisms and application retry loops belong here.
#registry-card("v2-12-lock-vs-version-upper-bound")
